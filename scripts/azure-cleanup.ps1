# ──────────────────────────────────────────────
# azure-cleanup.ps1 – Destroy all Azure resources
# ──────────────────────────────────────────────
# Usage: powershell -ExecutionPolicy Bypass -File scripts/azure-cleanup.ps1
# ──────────────────────────────────────────────

param (
    [string]$DbPassword = "",
    [switch]$Force
)

$ErrorActionPreference = "Stop"
$TfDir = "$PSScriptRoot\..\infra\azure\terraform"

Write-Host "`n═══ STORM Azure Cleanup ═══`n" -ForegroundColor Red

if (-not $Force) {
    $confirm = Read-Host "⚠️  This will DESTROY all Azure resources. Type 'yes' to confirm"
    if ($confirm -ne "yes") {
        Write-Host "Aborted." -ForegroundColor Yellow
        exit 0
    }
}

if ($DbPassword -eq "") {
    $DbPassword = Read-Host "Enter the DB password used during deployment"
}

Write-Host "💥  Destroying infrastructure..." -ForegroundColor Red
Push-Location $TfDir
terraform destroy -auto-approve -var "db_admin_password=$DbPassword"
Pop-Location

Write-Host "`n✅  All Azure resources destroyed." -ForegroundColor Green
