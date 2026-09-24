<#
.SYNOPSIS
    Dev run script for icmongolang with per-step options.

.DESCRIPTION
    Runs the maintenance pipeline step-by-step and can start the dev server
    (air) at the end. Each Go / Swagger step is optional; combine any set of
    switches or use -All for the full pipeline.

.PARAMETER Clean
    go clean -cache && go clean -modcache

.PARAMETER Tidy
    go mod tidy

.PARAMETER Download
    go mod download

.PARAMETER Verify
    go mod verify

.PARAMETER Migrate
    go run cmd/api/main.go migrate

.PARAMETER Swag
    swag init -g cmd/api/main.go

.PARAMETER Vendor
    go mod vendor

.PARAMETER Test
    go test ./...

.PARAMETER Run
    Start the dev server with air after the selected build steps.

.PARAMETER All
    Shortcut for: Clean Tidy Download Verify Migrate Swag Vendor Test Run.

.PARAMETER AirArgs
    Extra arguments forwarded to air (e.g. "-config .air.toml").

.EXAMPLE
    .\run.ps1 -All

.EXAMPLE
    .\run.ps1 -Migrate -Swag -Run

.EXAMPLE
    .\run.ps1 -Test

.NOTES
    Requires: Go, swag (go install github.com/swaggo/swag/cmd/swag@latest),
    air (go install github.com/air-verse/air@latest).
#>
[CmdletBinding()]
param(
    [switch]$Clean,
    [switch]$Tidy,
    [switch]$Download,
    [switch]$Verify,
    [switch]$Migrate,
    [switch]$Swag,
    [switch]$Vendor,
    [switch]$Test,
    [switch]$Run,
    [switch]$All,
    [string]$AirArgs = ""
)

$ErrorActionPreference = 'Stop'

if ($All) {
    $Clean = $Tidy = $Download = $Verify = $Migrate = $Swag = $Vendor = $Test = $Run = $true
}

$any = $Clean -or $Tidy -or $Download -or $Verify -or $Migrate -or $Swag -or $Vendor -or $Test -or $Run
if (-not $any) {
    Write-Host "Usage: .\run.ps1 [-Clean] [-Tidy] [-Download] [-Verify] [-Migrate] [-Swag] [-Vendor] [-Test] [-Run] [-All] [-AirArgs <args>]"
    Write-Host "  -All  runs everything then starts air."
    exit 0
}

function Step {
    param(
        [string]$Name,
        [scriptblock]$Body
    )
    Write-Host "== $Name" -ForegroundColor Cyan
    & $Body
    if ($LASTEXITCODE -ne 0) {
        Write-Error "$Name failed (exit code $LASTEXITCODE)"
        exit 1
    }
}

function Assert-Tool {
    param([string]$Name, [string]$InstallHint)
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        Write-Error "Tool '$Name' not found on PATH. Install it with: $InstallHint"
        exit 1
    }
}

if ($Clean)    { Step "go clean -cache / -modcache"        { go clean -cache; go clean -modcache } }
if ($Tidy)     { Step "go mod tidy"                        { go mod tidy } }
if ($Download) { Step "go mod download"                    { go mod download } }
if ($Verify)   { Step "go mod verify"                      { go mod verify } }
if ($Migrate)  { Step "go run cmd/api/main.go migrate"     { go run cmd/api/main.go migrate } }
if ($Swag) {
    Assert-Tool swag "go install github.com/swaggo/swag/cmd/swag@latest"
    Step "swag init -g cmd/api/main.go" { swag init -g cmd/api/main.go }
}
if ($Vendor)   { Step "go mod vendor"                      { go mod vendor } }
if ($Test)     { Step "go test ./..."                      { go test ./... } }

if ($Run) {
    Assert-Tool air "go install github.com/air-verse/air@latest"
    Write-Host "== starting air" -ForegroundColor Green
    if ($AirArgs) {
        air $AirArgs
    } else {
        air
    }
}