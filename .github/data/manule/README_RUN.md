ใน Windows PowerShell ไม่มีคำสั่ง `make` (GNU make) โดยค่าเริ่มต้น แต่คุณสามารถรันคำสั่ง Go และเครื่องมืออื่น ๆ ได้โดยตรง

มี **3 วิธี** ที่แนะนำ:

---
# API

- http://localhost:5000/api/health
- http://localhost:5000/api/ping
- http://localhost:5000/swagger/index.html
- http://localhost:5000/apimetric
- http://localhost:5000/metrics



## วิธีที่ 1: รันทีละคำสั่ง (ง่ายที่สุด)

คัดลอกและวางทีละบรรทัดลงใน PowerShell:

```powershell

# 1. ล้าง cache
go clean -cache
go clean -modcache

# 2. จัดการ dependencies
go mod tidy
go mod download
go mod verify

# 3. Migrate database
go run cmd/api/main.go migrate

# 4. สร้าง Swagger docs (ต้องติดตั้ง swag ก่อน)
swag init -g cmd/api/main.go

# 5. Vendor (ถ้าต้องการ)
go mod vendor

# 6. ทดสอบ
go test ./...

# 7. รัน server (ใช้ air หรือ go run)
air




# หรือถ้าไม่มี air:
go run cmd/api/main.go serve

```

---

## วิธีที่ 2: สร้าง PowerShell script (`build.ps1`)

สร้างไฟล์ `build.ps1` ใน root ของโปรเจกต์ ด้วยเนื้อหา:

```powershell
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
```

จากนั้นรันด้วย:

```powershell
.\build.ps1
```

> **หมายเหตุ:** ถ้า Windows ปิดกั้นการรันสคริปต์ ให้ใช้คำสั่งนี้ก่อน:
> ```powershell
> Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
> ```

---

## วิธีที่ 3: ติดตั้ง GNU make (ทางเลือก)

ถ้าคุณยังต้องการใช้ `make` จริง ๆ สามารถติดตั้งผ่าน Chocolatey:

```powershell
choco install make
```

(ต้องติดตั้ง Chocolatey ก่อน: https://chocolatey.org/install)

จากนั้นก็ใช้ `make run` ได้ตามที่ต้องการ

---

## สรุปคำสั่งที่ใช้บ่อย

| คำสั่ง | ความหมาย |
|-------|----------|
| `go clean -cache` | ล้าง build cache |
| `go clean -modcache` | ล้าง module cache |
| `go mod tidy` | จัดระเบียบ go.mod / go.sum |
| `go mod download` | ดาวน์โหลด dependencies |
| `go mod verify` | ตรวจสอบความถูกต้องของ dependencies |
| `go run cmd/api/main.go migrate` | รัน migration |
| `swag init -g cmd/api/main.go` | สร้าง Swagger docs |
| `go mod vendor` | คัดลอก dependencies ไปที่ vendor/ |
| `go test ./...` | วิ่งทดสอบทุก package |
| `air` | รัน hot-reload dev server |
| `go run cmd/api/main.go serve` | รัน server แบบธรรมดา |

เลือกวิธีที่สะดวกที่สุดสำหรับคุณครับ 🚀

✅ Starting Kafka + WebSocket service...
✅ PostgreSQL connected
✅ Kafka producer ready
✅ WebSocket hub started
✅ Kafka consumer started
🚀 HTTP server listening on port 5000
📚 Swagger UI at http://localhost:5000/swagger/index.html
🔌 WebSocket endpoint at ws://localhost:5000/ws


/*
	icmongolang/
		├── cmd/
		│   └── websocket/
		│       └── main.go
		├── internal/
		│   ├── queue/
		│   │   ├── manager.go
		│   │   ├── noop_queue.go
		│   │   └── queue_test.go
		│   └── websocket/
		│       ├── models/
		│       │   └── ws_models.go
		│       ├── repository/
		│       │   ├── ws_repo.go
		│       │   └── postgres/
		│       │       └── ws_repo_pg.go
		│       ├── usecase/
		│       │   ├── ws_usecase.go
		│       │   └── ws_usecase_test.go
		│       └── delivery/
		│           ├── ws/
		│           │   └── handler.go
		│           └── http/
		│               └── handler.go
		├── pkg/
		│   ├── logger/
		│   │   └── zap_logger.go    (มีอยู่แล้ว)
		│   └── websocket/
		│       ├── hub.go
		│       ├── client.go
		│       ├── message.go
		│       └── interface.go
		└── migrations/
			└── 20250619_websocket_tables.sql
*/


http://172.17.0.1:9087


NmaR7cCkQirwnLvKZakZXyl12NTUyLNV0-gjAdBtIiN0VYpwA6QHfh0-pjKfJCOZZJsHdgZJEEbLKeITEyNE1w==


http://localhost:9087








