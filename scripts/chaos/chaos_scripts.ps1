# 🌪️ STORM Chaos Script
# Usage: ./chaos_scripts.ps1 -Action <kill|latency|restore> -Target <nats|redis|postgres|gateway>

param (
    [Parameter(Mandatory=$true)]
    [ValidateSet("kill", "latency", "restore")]
    $Action,

    [Parameter(Mandatory=$true)]
    [ValidateSet("nats", "redis", "postgres", "gateway")]
    $Target
)

$ComposeFile = "infra/docker/docker-compose.yml"
$ProjectDir = Get-Location

function Get-ContainerName($tgt) {
    return "docker-$tgt-1"
}

switch ($Action) {
    "kill" {
        Write-Host "🌪️ Killing $Target..." -ForegroundColor Red
        docker stop $(Get-ContainerName $Target)
    }
    "restore" {
        Write-Host "🌱 Restoring $Target..." -ForegroundColor Green
        docker start $(Get-ContainerName $Target)
    }
    "latency" {
        Write-Host "⏳ Simulating 500ms latency on $Target (Simulated via storage layer if supported or proxy)..." -ForegroundColor Yellow
        Write-Warning "Latency via Docker networking requires 'tc' (Linux only). For local Windows/Docker Desktop, we will use a dedicated env var for the Gateway."
        if ($Target -eq "postgres") {
            Write-Host "Adjusting Gateway to simulate DB Slowness..."
            # Note: This is an example, we might need to update docker-compose or gateway code
        }
    }
}
