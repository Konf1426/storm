# ──────────────────────────────────────────────────────
# azure-load-test.ps1 – Run distributed k6 load tests
# ──────────────────────────────────────────────────────
# Usage:
#   ./scripts/azure-load-test.ps1                   # default 25k VUs per pod (100k total)
#   ./scripts/azure-load-test.ps1 -TargetVUs 10000  # 10k per pod (40k total)
# ──────────────────────────────────────────────────────

param (
    [int]$TargetVUs  = 25000,
    [int]$Pods       = 4,
    [string]$RampUp  = "3m",
    [string]$Hold    = "5m",
    [string]$RampDown = "1m"
)

$ErrorActionPreference = "Stop"
$JobFile = "$PSScriptRoot\..\infra\k8s-azure\k6-runner\k6-job.yaml"

Write-Host "`n═══ STORM Load Test ($Pods pods × $TargetVUs VUs = $($Pods * $TargetVUs) total) ═══`n" -ForegroundColor Cyan

# ─── 1. Update k6 script ConfigMap ───────────────────
Write-Host "📦  Updating k6 script ConfigMap..." -ForegroundColor Yellow
kubectl create configmap k6-script `
    --from-file=storm-azure.js="$PSScriptRoot\k6\storm-azure.js" `
    -n storm --dry-run=client -o yaml | kubectl apply -f -

# ─── 2. Clean previous job if exists ─────────────────
Write-Host "🧹  Cleaning previous job..." -ForegroundColor Yellow
kubectl delete job k6-load-test -n storm --ignore-not-found

# ─── 3. Patch and apply job ──────────────────────────
Write-Host "🚀  Launching k6 runners..." -ForegroundColor Yellow

# Create a temp copy with the right values
$JobContent = Get-Content $JobFile -Raw
$JobContent = $JobContent -replace 'parallelism: \d+', "parallelism: $Pods"
$JobContent = $JobContent -replace 'completions: \d+', "completions: $Pods"
$JobContent = $JobContent -replace 'value: "\d+"', "value: `"$TargetVUs`""  # TARGET_VUS
$TempFile = [System.IO.Path]::GetTempFileName() + ".yaml"
$JobContent | Set-Content $TempFile
kubectl apply -f $TempFile
Remove-Item $TempFile

# ─── 4. Follow logs ─────────────────────────────────
Write-Host "`n📊  Streaming logs (Ctrl+C to stop watching)..." -ForegroundColor Yellow
Write-Host "    Monitor HPA:  kubectl get hpa -n storm -w" -ForegroundColor DarkGray
Write-Host "    Monitor pods: kubectl get pods -n storm -w" -ForegroundColor DarkGray
Write-Host ""

Start-Sleep -Seconds 5
kubectl logs -f -l app=k6-runner -n storm --max-log-requests=$Pods

# ─── 5. Summary ─────────────────────────────────────
Write-Host "`n═══ Load test complete ═══" -ForegroundColor Cyan
kubectl get jobs -n storm
kubectl get hpa -n storm
