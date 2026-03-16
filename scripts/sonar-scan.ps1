$ErrorActionPreference = "Stop"

$RootDir = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$ComposeFile = Join-Path $RootDir "infra\docker\docker-compose.yml"
$SonarHostUrl = if ($env:SONAR_HOST_URL) { $env:SONAR_HOST_URL } else { "http://host.docker.internal:9000" }
$SonarToken = $env:SONAR_TOKEN

if (-not $SonarToken) {
  throw "SONAR_TOKEN is required. Create a token in SonarQube (My Account > Security) and set `$env:SONAR_TOKEN."
}

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
  throw "docker is not installed or not in PATH"
}

Write-Host "`n==> Start SonarQube"
docker compose -f $ComposeFile up -d sonarqube | Out-Null

Write-Host "`n==> Wait for SonarQube API"
$maxAttempts = 60
for ($i = 1; $i -le $maxAttempts; $i++) {
  try {
    $res = Invoke-RestMethod -Uri "http://localhost:9000/api/system/status" -TimeoutSec 5
    if ($res.status -eq "UP") {
      break
    }
  } catch {
    Start-Sleep -Seconds 2
  }
  if ($i -eq $maxAttempts) {
    throw "SonarQube is not ready on http://localhost:9000"
  }
  Start-Sleep -Seconds 2
}

$SonarDir = Join-Path $RootDir ".sonar"
New-Item -ItemType Directory -Path $SonarDir -Force | Out-Null

Write-Host "`n==> Generate Go coverage reports"
docker run --rm -v "${RootDir}:/repo" -w /repo/services/gateway golang:1.24-alpine sh -lc "go test ./... -coverprofile /repo/.sonar/coverage-gateway.out"
docker run --rm -v "${RootDir}:/repo" -w /repo/services/messages golang:1.24-alpine sh -lc "go test ./... -coverprofile /repo/.sonar/coverage-messages.out"

Write-Host "`n==> Run Sonar Scanner"
docker run --rm `
  -e SONAR_HOST_URL=$SonarHostUrl `
  -e SONAR_TOKEN=$SonarToken `
  -v "${RootDir}:/usr/src" `
  sonarsource/sonar-scanner-cli

Write-Host "`nSonarQube scan complete. Dashboard: http://localhost:9000"
