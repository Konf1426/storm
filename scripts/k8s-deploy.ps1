$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $PSScriptRoot
$K8sDir = Join-Path $RootDir "infra\k8s"
$Namespace = if ($env:NAMESPACE) { $env:NAMESPACE } else { "storm" }
$BuildImages = if ($env:BUILD_IMAGES) { $env:BUILD_IMAGES } else { "1" }
$GatewayImage = if ($env:GATEWAY_IMAGE) { $env:GATEWAY_IMAGE } else { "storm-gateway:latest" }
$MessagesImage = if ($env:MESSAGES_IMAGE) { $env:MESSAGES_IMAGE } else { "storm-messages:latest" }

if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
  throw "kubectl is not installed or not in PATH"
}

if ($BuildImages -eq "1") {
  if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    throw "docker is not installed or not in PATH"
  }
  Write-Host "`n==> Build local images"
  docker build -t $GatewayImage (Join-Path $RootDir "services\gateway")
  docker build -t $MessagesImage (Join-Path $RootDir "services\messages")
}

Write-Host "`n==> Apply manifests"
kubectl apply -k $K8sDir

Write-Host "`n==> Wait for core deployments"
kubectl -n $Namespace rollout status deploy/nats --timeout=180s
kubectl -n $Namespace rollout status deploy/postgres --timeout=180s
kubectl -n $Namespace rollout status deploy/redis --timeout=180s
kubectl -n $Namespace rollout status deploy/gateway --timeout=180s
kubectl -n $Namespace rollout status deploy/messages --timeout=180s

Write-Host "`n==> Pods"
kubectl -n $Namespace get pods -o wide

Write-Host "`n==> Services"
kubectl -n $Namespace get svc
