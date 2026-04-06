param (
    [string]$DbPassword = "",
    [switch]$SkipTerraform,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"
$TfDir = "$PSScriptRoot\..\infra\azure\terraform"
$K8sDir = "$PSScriptRoot\..\infra\k8s-azure"
$ScriptDir = "$PSScriptRoot\..\scripts\k6"

Write-Host "Starting Azure Deployment..." -ForegroundColor Cyan

function Run-Command {
    param($Command, $Arguments)
    & $Command @Arguments
    if ($LastExitCode -ne 0) {
        throw "Command '$Command' failed with exit code $LastExitCode"
    }
}

if (-not $SkipTerraform) {
    Write-Host "Step 1: Terraform provisioning..." -ForegroundColor Yellow
    if ($DbPassword -eq "") {
        $DbPassword = Read-Host "Enter PostgreSQL admin password" -AsSecureString | ForEach-Object { [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($_)) }
    }
    Push-Location $TfDir
    Run-Command "terraform" @("init", "-input=false")
    Run-Command "terraform" @("apply", "-auto-approve", "-var", "db_admin_password=$DbPassword")

    $ACR_SERVER = terraform output -raw acr_login_server
    $AKS_NAME = terraform output -raw aks_cluster_name
    $RG_NAME = terraform output -raw resource_group_name
    $PG_DSN = terraform output -raw postgres_dsn
    $REDIS_ADDR = terraform output -raw redis_connection_string
    $JWT = terraform output -raw jwt_secret
    $JWT_REFRESH = terraform output -raw jwt_refresh_secret
    Pop-Location
} else {
    Write-Host "Skipping Terraform..." -ForegroundColor DarkGray
    Push-Location $TfDir
    $ACR_SERVER = terraform output -raw acr_login_server
    $AKS_NAME = terraform output -raw aks_cluster_name
    $RG_NAME = terraform output -raw resource_group_name
    $PG_DSN = terraform output -raw postgres_dsn
    $REDIS_ADDR = terraform output -raw redis_connection_string
    $JWT = terraform output -raw jwt_secret
    $JWT_REFRESH = terraform output -raw jwt_refresh_secret
    Pop-Location
}

if (-not $SkipBuild) {
    Write-Host "Step 2: Building Docker images..." -ForegroundColor Yellow
    $AcrName = ($ACR_SERVER -split "\.")[0]
    Run-Command "az" @("acr", "login", "--name", $AcrName)
    
    docker build -t "$($ACR_SERVER)/storm-gateway:latest" -f services/gateway/Dockerfile services/gateway
    docker push "$($ACR_SERVER)/storm-gateway:latest"
    
    docker build -t "$($ACR_SERVER)/storm-messages:latest" -f services/messages/Dockerfile services/messages
    docker push "$($ACR_SERVER)/storm-messages:latest"
}

Write-Host "Step 3: Fetching AKS credentials..." -ForegroundColor Yellow
Run-Command "az" @("aks", "get-credentials", "--resource-group", $RG_NAME, "--name", $AKS_NAME, "--overwrite-existing")

Write-Host "Step 4: Deploying to AKS..." -ForegroundColor Yellow
$KustomFile = "$K8sDir\kustomization.yaml"
$KContent = Get-Content $KustomFile -Raw
$KContent = $KContent.Replace("REPLACE_WITH_ACR.azurecr.io", $ACR_SERVER)
Set-Content -Path $KustomFile -Value $KContent

$SecretFile = "$K8sDir\secret-azure.yaml"
$SContent = Get-Content $SecretFile -Raw
$SContent = $SContent.Replace("REPLACE_PASSWORD", $DbPassword)

$PgHost = ($PG_DSN -split "@")[1].Split(":")[0]
$SContent = $SContent.Replace("REPLACE_PG_FQDN", $PgHost)

$RedisHost = ($REDIS_ADDR -split ":")[0]
$SContent = $SContent.Replace("REPLACE_REDIS_HOSTNAME", $RedisHost)

$SContent = $SContent.Replace("REPLACE_JWT_SECRET", $JWT)
$SContent = $SContent.Replace("REPLACE_JWT_REFRESH_SECRET", $JWT_REFRESH)
Set-Content -Path $SecretFile -Value $SContent

kubectl create configmap k6-script --from-file=storm-azure.js="$ScriptDir\storm-azure.js" -n storm --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -k $K8sDir

Write-Host "Step 5: Waiting for pods..." -ForegroundColor Yellow
kubectl rollout status deployment/gateway -n storm --timeout=120s
kubectl rollout status deployment/nats -n storm --timeout=60s

Write-Host "Deployment complete." -ForegroundColor Green
