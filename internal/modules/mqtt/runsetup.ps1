# 1. format
gofmt -w internal\modules\mqtt\

# 2. build
go build ./internal/modules/mqtt/...

# 3. unit test
go test ./internal/modules/mqtt/delivery/http/...

# 4. clean + regenerate swagger
Remove-Item docs\swagger.json -Force -ErrorAction SilentlyContinue
Remove-Item docs\swagger.yaml -Force -ErrorAction SilentlyContinue
Remove-Item docs\docs.go     -Force -ErrorAction SilentlyContinue

swag init `
  -g cmd/api/main.go `
  -o docs `
  --parseDependency `
  --parseInternal `
  --parseDepth 3

# 5. ตรวจสอบ schema
Select-String -Path docs\swagger.json -Pattern "presenter.StatusResponse"
Select-String -Path docs\swagger.json -Pattern "presenter.ErrResponse"
Select-String -Path docs\swagger.json -Pattern "swaggertype"
Select-String -Path docs\swagger.json -Pattern "BAACTW05"

# 6. run
go run cmd/api/main.go serve