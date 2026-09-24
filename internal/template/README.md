# Template Module

โครงสร้าง template สำหรับสร้าง module ใหม่ ตามรูปแบบ DDD (Domain-Driven Design) ที่ใช้ในโปรเจกต์

## โครงสร้าง

```
internal/template/
├── handler.go              # Handler interface (HTTP methods)
├── usecase.go              # UseCase interface (business logic)
├── pg_repository.go        # PgRepository interface (data access)
├── redis_repository.go     # RedisRepository interface (cache)
├── worker.go               # Task distributor/processor interfaces (async)
├── models/
│   └── model.go            # GORM model
├── presenter/
│   └── presenter.go        # Request/Response DTOs
├── delivery/http/
│   ├── handler.go          # HTTP handler implementation
│   └── routes.go           # Route mapping
├── usecase/
│   └── usecase.go          # UseCase implementation
├── repository/
│   ├── pg_repository.go    # PostgreSQL repository implementation
│   └── redis_repository.go # Redis repository implementation
├── processor/
│   └── processor.go        # Async task processor (optional)
├── distributor/
│   └── distributor.go      # Async task distributor (optional)
└── migration.sql           # SQL migration script
```

## วิธีใช้งาน

### 1. คัดลอก template

```bash
cp -r internal/template internal/<module-name>
```

### 2. เปลี่ยนชื่อ package และ struct

ใช้ replace all ใน IDE:
- `template` → `<module-name>` (package name)
- `Template` → `<ModuleName>` (struct name)
- `templates` → `<table-name>` (table name)

### 3. แก้ไข model

`internal/<module-name>/models/model.go`
- แก้ไข field ให้ตรงกับตารางฐานข้อมูล
- เปลี่ยน `TableName()` return value

### 4. แก้ไข presenter

`internal/<module-name>/presenter/presenter.go`
- เพิ่ม/แก้ไข request/response structs

### 5. แก้ไข repository

`internal/<module-name>/repository/pg_repository.go`
- เพิ่ม custom query methods
- เปลี่ยนเงื่อนไขตาม business logic

### 6. แก้ไข usecase

`internal/<module-name>/usecase/usecase.go`
- เพิ่ม business logic methods
- เพิ่ม Redis cache pattern (ถ้าต้องการ)

### 7. แก้ไข handler

`internal/<module-name>/delivery/http/handler.go`
- เพิ่ม/แก้ไข endpoint handlers
- แก้ไข Swagger annotations

### 8. แก้ไข routes

`internal/<module-name>/delivery/http/routes.go`
- เปลี่ยน route path
- ปรับ middleware ตามต้องการ

### 9. ลงทะเบียนใน router

แก้ไข `internal/delivery/rest/router.go`:
```go
import (
    templateHttp "icmongolang/internal/<module-name>/delivery/http"
    templateRepo "icmongolang/internal/<module-name>/repository"
    templateUC "icmongolang/internal/<module-name>/usecase"
)

// สร้าง dependencies
pgRepo := templateRepo.NewTemplatePgRepository(db)
redisRepo := templateRepo.NewTemplateRedisRepository(redisClient)
uc := templateUC.NewTemplateUseCase(pgRepo, redisRepo, cfg, log)
handler := templateHttp.NewTemplateHandler(uc, cfg, log)

// Mount routes
templateHttp.MapTemplateRoutes(apiRouter, handler, mw)
```

### 10. สร้าง migration

แก้ไข `internal/<module-name>/migration.sql` และรันกับฐานข้อมูล

## กรณีที่ไม่ต้องใช้บาง layer

| Layer | เมื่อไม่ต้องใช้ |
|-------|----------------|
| `redis_repository.go` | ไม่ต้องใช้ Redis cache |
| `processor/`, `distributor/` | ไม่มี async task (email, notification) |
| `models/` | ถ้าใช้ model จาก `internal/models/` ร่วมกัน |

## ดูตัวอย่าง module จริง

| Module | Features |
|--------|----------|
| `internal/modules/items/` | CRUD พื้นฐาน, ไม่มี Redis |
| `internal/modules/alarm/` | CRUD + processor |
| `internal/modules/users/` | CRUD + Redis cache + async task (email) + distributor/processor |
| `internal/modules/iot/` | CRUD + MQTT + InfluxDB + Redis |
