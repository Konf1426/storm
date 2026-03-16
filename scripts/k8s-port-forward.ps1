$ErrorActionPreference = "Stop"

$Namespace = if ($env:NAMESPACE) { $env:NAMESPACE } else { "storm" }

if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
  throw "kubectl is not installed or not in PATH"
}

$gw = Start-Process kubectl -ArgumentList "-n", $Namespace, "port-forward", "svc/gateway", "8080:8080" -PassThru
$gf = Start-Process kubectl -ArgumentList "-n", $Namespace, "port-forward", "svc/grafana", "3000:3000" -PassThru
$pr = Start-Process kubectl -ArgumentList "-n", $Namespace, "port-forward", "svc/prometheus", "9090:9090" -PassThru

Write-Host "Port-forward active:"
Write-Host "  gateway    http://localhost:8080"
Write-Host "  grafana    http://localhost:3000"
Write-Host "  prometheus http://localhost:9090"
Write-Host "Press Enter to stop..."
[void][System.Console]::ReadLine()

foreach ($p in @($gw, $gf, $pr)) {
  if ($p -and -not $p.HasExited) {
    Stop-Process -Id $p.Id -Force
  }
}
