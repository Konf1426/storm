param(
  [string]$SubscriptionId,
  [string]$ResourceGroup = "rg-storm-loadtest",
  [string]$Location = "westeurope",
  [string]$LoadTestResource = "storm-loadtest",
  [string]$TestId = "storm-ws-100k",
  [string]$TestRunId = "",
  [string]$ConfigFile = "",
  [switch]$CreateResource,
  [switch]$CreateTestOnly
)

$ErrorActionPreference = "Stop"

function Invoke-AzCli {
  param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$Args
  )

  & az @Args --only-show-errors
  if ($LASTEXITCODE -ne 0) {
    throw "Azure CLI command failed: az $($Args -join ' ')"
  }
}

$RootDir = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
if (-not $ConfigFile) {
  $ConfigFile = Join-Path $PSScriptRoot "loadtest-100k.yaml"
}

if (-not (Get-Command az -ErrorAction SilentlyContinue)) {
  throw "Azure CLI is not installed or not in PATH. Install it first: https://learn.microsoft.com/cli/azure/install-azure-cli"
}

& az extension add --name load --only-show-errors 2>$null
if ($LASTEXITCODE -ne 0) {
  throw "Azure CLI extension install failed: az extension add --name load --only-show-errors"
}

if ($SubscriptionId) {
  Invoke-AzCli account set --subscription $SubscriptionId
}

if ($CreateResource) {
  Write-Host "==> Create resource group"
  Invoke-AzCli group create `
    --name $ResourceGroup `
    --location $Location `
    --output table

  Write-Host "==> Create Azure Load Testing resource"
  Invoke-AzCli load create `
    --name $LoadTestResource `
    --resource-group $ResourceGroup `
    --location $Location `
    --output table
}

Write-Host "==> Create or update the load test"
$testExists = $false
& az load test show `
  --load-test-resource $LoadTestResource `
  --resource-group $ResourceGroup `
  --test-id $TestId `
  --output none `
  --only-show-errors 2>$null
if ($LASTEXITCODE -eq 0) {
  $testExists = $true
}

if ($testExists) {
  Write-Host "Test '$TestId' already exists, updating it"
  Invoke-AzCli load test update `
    --load-test-resource $LoadTestResource `
    --resource-group $ResourceGroup `
    --test-id $TestId `
    --load-test-config-file $ConfigFile `
    --output table
}
else {
  Invoke-AzCli load test create `
    --load-test-resource $LoadTestResource `
    --resource-group $ResourceGroup `
    --test-id $TestId `
    --load-test-config-file $ConfigFile `
    --output table
}

if ($CreateTestOnly) {
  exit 0
}

if (-not $TestRunId) {
  $TestRunId = "run-{0}" -f (Get-Date -Format "yyyyMMdd-HHmmss")
}

Write-Host "==> Start test run"
Invoke-AzCli load test-run create `
  --load-test-resource $LoadTestResource `
  --resource-group $ResourceGroup `
  --test-id $TestId `
  --test-run-id $TestRunId `
  --display-name $TestRunId `
  --description "STORM 100k WebSocket test via Azure CLI" `
  --output table

Write-Host "==> Show test run"
Invoke-AzCli load test-run show `
  --load-test-resource $LoadTestResource `
  --resource-group $ResourceGroup `
  --test-run-id $TestRunId `
  --output table
