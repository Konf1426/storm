$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $PSScriptRoot
$K8sDir = Join-Path $RootDir "infra\k8s"

if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
  throw "kubectl is not installed or not in PATH"
}

kubectl delete -k $K8sDir --ignore-not-found=true
