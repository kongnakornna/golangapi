param([string]$Action = "up")

$ErrorActionPreference = "Stop"
$projectCompose = Join-Path $PSScriptRoot "..\docker-compose.yml"
$services = "logstash kibana"

if (-not (Test-Path $projectCompose)) {
    throw "docker-compose.yml not found at $projectCompose"
}

if ($Action -eq "up") {
    Write-Host "Raising vm.max_map_count in docker-desktop WSL (needed by Elasticsearch)..."
    wsl -d docker-desktop -u root sysctl -w vm.max_map_count=262144

    Write-Host "Starting Logstash and Kibana (reusing the project's Elasticsearch)..."
    docker compose -f $projectCompose up -d $services
    if ($LASTEXITCODE -ne 0) { throw "docker compose up failed" }

    Write-Host ""
    Write-Host "Kibana:   http://localhost:5601"
    Write-Host "Logstash: tcp input on localhost:5001, beats on localhost:5044"
} elseif ($Action -eq "down") {
    docker compose -f $projectCompose stop $services
} elseif ($Action -eq "logs") {
    docker compose -f $projectCompose logs -f $services
} else {
    throw "Unknown action '$Action'. Use up, down, or logs."
}
