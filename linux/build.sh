#!/bin/bash
set -e

echo -e "\033[36m🧹 Cleaning cache...\033[0m"
go clean -cache
go clean -modcache

echo -e "\033[36m📦 Tidying modules...\033[0m"
go mod tidy

echo -e "\033[36m⬇️ Downloading modules...\033[0m"
go mod download

echo -e "\033[36m✅ Verifying modules...\033[0m"
go mod verify

echo -e "\033[36m🗄️ Running migration...\033[0m"
go run cmd/api/main.go migrate

echo -e "\033[36m📄 Generating Swagger docs...\033[0m"
swag init -g cmd/api/main.go

echo -e "\033[36m📁 Vendoring...\033[0m"
go mod vendor

echo -e "\033[36m🧪 Running tests...\033[0m"
go test ./...

echo -e "\033[32m🚀 Starting server with Air...\033[0m"
air