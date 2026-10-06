# 1. format
gofmt -w internal\modules\keys\

# 2. build
go build ./internal/modules/keys/...

# 3. ปิด editor ก่อน, ลบ swagger เก่า
Remove-Item docs\swagger.json -Force -ErrorAction SilentlyContinue
Remove-Item docs\swagger.yaml -Force -ErrorAction SilentlyContinue
Remove-Item docs\docs.go     -Force -ErrorAction SilentlyContinue

# 4. regenerate
swag init `
  -g cmd/api/main.go `
  -o docs `
  --parseDependency `
  --parseInternal `
  --parseDepth 3

# 5. ตรวจสอบ
Select-String -Path docs\swagger.json -Pattern "api-keys"
Select-String -Path docs\swagger.json -Pattern "APIKeyWithSecretResponse"

# 6. เปิด UI
go run cmd/api/main.go serve
# → http://localhost:5000/swagger/index.html → tag "api-keys"