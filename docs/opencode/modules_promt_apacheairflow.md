# ROLE
คุณเป็น Senior Go Developer ที่เชี่ยวชาญ Clean Architecture + DDD
และ Modular Monolith Pattern

# CONTEXT
โปรเจกต์: icmongolang
Go version: 1.22+
Architecture: Clean Architecture + DDD
Structure: Modular Monolith (33 modules)

โครงสร้าง Module มาตรฐาน:
internal/modules/apacheairflow/
├── domain/.
│   ├── entity/{entity}.go
│   ├── value_object/{vo}.go
│   ├── repository/{entity}_repository.go
│   └── errors/errors.go
├── application/
│   ├── {verb}_{entity}.go
│   └── mocks/{entity}_repository_mock.go
├── infrastructure/
│   └── persistence/postgres/
│       ├── models.go
│       └── {entity}_repo_impl.go
├── interfaces/
│   └── http/
│       ├── {entity}_handler.go
│       ├── dto.go
│       └── routes.go
└── module.go

Shared packages ที่ใช้ได้:
- icmongolang/pkg/logger (structured logging)
- icmongolang/pkg/validator (input validation)
- icmongolang/pkg/httputil (JSON responses)
- github.com/google/uuid
- github.com/shopspring/decimal
- gorm.io/gorm
- github.com/go-chi/chi/v5

# TASK  apacheairflow
สร้าง Module ใหม่ชื่อ "apacheairflow" ที่มี:

## Business Requirements
{business_requirements}

## Entities
{entity_list}

## Use Cases
{use_case_list}

## API Endpoints
{api_endpoints}

# CONSTRAINTS
1. Domain Layer ห้าม import: gorm, gin, chi, sarama, redis
2. Entity ต้องมี constructor `New{Entity}` + validation
3. เปลี่ยน state ผ่าน behavior methods เท่านั้น (ไม่มี setter)
4. Repository เป็น interface เท่านั้น
5. Errors เป็น sentinel errors
6. Comment 2 ภาษา (ไทย/อังกฤษ)
7. Test coverage ≥ 80%
8. Table prefix: apacheairflow_
9. Redis key: apacheairflow:{entity}:{id}

# FORMAT
ตอบในรูปแบบ:

## 1. Architecture Diagram
[แสดงโครงสร้าง module]

## 2. File Tree
internal/modules/apacheairflow/ ├── ... (ระบุทุกไฟล์)


## 3. Code ทั้งหมด
### 3.1 Domain Layer
- entity/{entity}.go (พร้อม test)
- value_object/{vo}.go (พร้อม test)
- repository/{entity}_repository.go
- errors/errors.go

### 3.2 Application Layer
- application/{verb}_{entity}.go (ทุก use case)
- application/mocks/{entity}_repository_mock.go

### 3.3 Infrastructure Layer
- infrastructure/persistence/postgres/models.go
- infrastructure/persistence/postgres/{entity}_repo_impl.go

### 3.4 Interface Layer
- interfaces/http/{entity}_handler.go
- interfaces/http/dto.go
- interfaces/http/routes.go

### 3.5 Module Bootstrap
- module.go
- Wire-up ใน main.go

## 4. Migration SQL
migrations/YYYYMMDD_apacheairflow_init.sql

## 5. Tests
- Unit tests ทุก layer
- Integration test
- Handler test

## 6. Documentation
- README.md ของ module
- API examples

# VERIFICATION
ก่อนส่งคำตอบ ให้ตรวจสอบ:
- [ ] Domain ไม่มี GORM tag
- [ ] Entity มี constructor + behavior methods
- [ ] Repository interface อยู่ใน domain
- [ ] Use case ละ 1 ไฟล์
- [ ] Handler ดึง user_id จาก context
- [ ] ทุก error เป็น sentinel
- [ ] Migration มี prefix
- [ ] Test ครอบ happy + error paths
- [ ] Comment ครบ 2 ภาษา