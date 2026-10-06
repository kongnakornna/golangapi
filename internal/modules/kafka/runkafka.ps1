# 1. ลบ docs เก่าให้สะอาด (ปิด VS Code / editor ก่อน)
Remove-Item docs\* -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path docs -Force | Out-Null

# 2. format
gofmt -w internal\modules\kafka\

# 3. build (ต้องผ่านก่อน swag init)
go build ./internal/modules/kafka/...

# 4. regenerate swagger
swag init `
  -g cmd/api/main.go `
  -o docs `
  --parseDependency `
  --parseInternal `
  --parseDepth 3

# 5. ตรวจสอบว่ามี docs.go
Get-ChildItem docs\
# → ต้องเห็น docs.go, swagger.json, swagger.yaml

# 6. ตรวจ schema ใหม่
Select-String -Path docs\swagger.json -Pattern "CreateOrderAcceptedResponse"
Select-String -Path docs\swagger.json -Pattern "models.ErrorResponse"
Select-String -Path docs\swagger.json -Pattern "PROD-001"

# 7. build ทั้งหมด
go build ./...

# 8. run
go run cmd/api/main.go serve