# ──────────────────────────────────────────────────────
# azure-deploy.ps1 – Deploy STORM on Azure AKS
# ──────────────────────────────────────────────────────
# Usage: powershell -ExecutionPolicy Bypass -File scripts/azure-deploy.ps1
# Prerequisites: az CLI, terraform, docker, kubectl
# ──────────────────────────────────────────────────────

param (
    [string]$DbPassword = "",
    [switch]$SkipTerraform,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"
$TfDir    = "$PSScriptRoot\..\infra\azure\terraform"
$K8sDir   = "$PSScriptRoot\..\infra\k8s-azure"
$ScriptDir = "$PSScriptRoot\..\scripts\k6"

Write-Host "`n═══ STORM Azure Deployment ═══`n" -ForegroundColor Cyan

# ─── 1. Terraform ─────────────────────────────────────
if (-not $SkipTerraform) {
    Write-Host "🏗️  Step 1/5: Provisioning Azure infrastructure..." -ForegroundColor Yellow

    if ($DbPassword -eq "") {
        $DbPassword = Read-Host "Enter PostgreSQL admin password" -AsSecureString |
            ForEach-Object { [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($_)) }
    }

    Push-Location $TfDir
    terraform init -input=false
    terraform apply -auto-approve -var "db_admin_password=$DbPassword"

    $ACR_SERVER  = terraform output -raw acr_login_server
    $AKS_NAME    = terraform output -raw aks_cluster_name
    $RG_NAME     = terraform output -raw resource_group_name
    $PG_DSN      = terraform output -raw postgres_dsn
    $REDIS_ADDR  = terraform output -raw redis_connection_string
    $JWT         = terraform output -raw jwt_secret
    $JWT_REFRESH = terraform output -raw jwt_refresh_secret
    Pop-Location

    Write-Host "✅  Infrastructure provisioned" -ForegroundColor Green
} else {
    Write-Host "⏭️  Skipping Terraform (--SkipTerraform)" -ForegroundColor DarkGray
    Push-Location $TfDir
    $ACR_SERVER  = terraform output -raw acr_login_server
    $AKS_NAME    = terraform output -raw aks_cluster_name
    $RG_NAME     = terraform output -raw resource_group_name
    $PG_DSN      = terraform output -raw postgres_dsn
    $REDIS_ADDR  = terraform output -raw redis_connection_string
    $JWT         = terraform output -raw jwt_secret
    $JWT_REFRESH = terraform output -raw jwt_refresh_secret
    Pop-Location
}

# ─── 2. Build & Push images ──────────────────────────
if (-not $SkipBuild) {
    Write-Host "`n🐳  Step 2/5: Building and pushing Docker images..." -ForegroundColor Yellow

    az acr login --name ($ACR_SERVER -split '\.')[0]

    docker build -t "${ACR_SERVER}/storm-gateway:latest" -f services/gateway/Dockerfile services/gateway
    docker push "${ACR_SERVER}/storm-gateway:latest"

    docker build -t "${ACR_SERVER}/storm-messages:latest" -f services/messages/Dockerfile services/messages
    docker push "${ACR_SERVER}/storm-messages:latest"

    Write-Host "✅  Images pushed to ACR" -ForegroundColor Green
} else {
    Write-Host "`n⏭️  Skipping Docker build (--SkipBuild)" -ForegroundColor DarkGray
}

# ─── 3. Get AKS credentials ─────────────────────────
Write-Host "`n🔑  Step 3/5: Fetching AKS credentials..." -ForegroundColor Yellow
az aks get-credentials --resource-group $RG_NAME --name $AKS_NAME --overwrite-existing
Write-Host "✅  kubectl configured" -ForegroundColor Green

# ─── 4. Update K8s manifests & deploy ────────────────
Write-Host "`n☸️  Step 4/5: Deploying to AKS..." -ForegroundColor Yellow

# Update kustomization with real ACR name
$KustomFile = "$K8sDir\kustomization.yaml"
(Get-Content $KustomFile) -replace 'REPLACE_WITH_ACR\.azurecr\.io', $ACR_SERVER | Set-Content $KustomFile

# Update secrets with real values
$SecretFile = "$K8sDir\secret-azure.yaml"
(Get-Content $SecretFile) `
    -replace 'REPLACE_PASSWORD', $DbPassword `
    -replace 'REPLACE_PG_FQDN', ($PG_DSN -split '@')[1].Split(':')[0] `
    -replace 'REPLACE_REDIS_HOSTNAME', ($REDIS_ADDR -split ':')[0] `
    -replace 'REPLACE_JWT_SECRET', $JWT `
    -replace 'REPLACE_JWT_REFRESH_SECRET', $JWT_REFRESH |
    Set-Content $SecretFile

# Create k6 ConfigMap from script
kubectl create configmap k6-script `
    --from-file=storm-azure.js="$ScriptDir\storm-azure.js" `
    -n storm --dry-run=client -o yaml | kubectl apply -f -

# Apply all manifests
kubectl apply -k $K8sDir

Write-Host "✅  Manifests applied" -ForegroundColor Green

# ─── 5. Health check ─────────────────────────────────
Write-Host "`n🩺  Step 5/5: Waiting for pods..." -ForegroundColor Yellow
kubectl rollout status deployment/gateway -n storm --timeout=120s
kubectl rollout status deployment/nats -n storm --timeout=60s

Write-Host "`n✅  All pods ready!" -ForegroundColor Green
kubectl get pods -n storm
kubectl get hpa -n storm

Write-Host "`n═══ Deployment complete ═══" -ForegroundColor Cyan
Write-Host "Gateway: kubectl port-forward svc/gateway 8080:8080 -n storm"
Write-Host "Grafana: kubectl port-forward svc/grafana 3000:3000 -n storm"
