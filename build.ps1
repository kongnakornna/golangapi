Write-Host "🧹 Cleaning cache..." -ForegroundColor Cyan
go clean -cache
go clean -modcache

Write-Host "📦 Tidying modules..." -ForegroundColor Cyan
go mod tidy

Write-Host "⬇️ Downloading modules..." -ForegroundColor Cyan
go mod download

Write-Host "✅ Verifying modules..." -ForegroundColor Cyan
go mod verify

Write-Host "🗄️ Running migration..." -ForegroundColor Cyan
go run cmd/api/main.go migrate

Write-Host "📄 Generating Swagger docs..." -ForegroundColor Cyan
swag init -g cmd/api/main.go

Write-Host "📁 Vendoring..." -ForegroundColor Cyan
go mod vendor

Write-Host "🧪 Running tests..." -ForegroundColor Cyan
go test ./...

Write-Host "🚀 Starting server with Air..." -ForegroundColor Green
air