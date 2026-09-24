
### โฟลเดอร์หลัก  icmongolang

swag init -g cmd/api/main.go -o docs
swag init -g cmd/api/main.go 
go mod tidy
go mod download
go mod verify
go run cmd/api/main.go migrate
swag init -g cmd/api/main.go 
go mod vendor
air


```
icmongolang/
├── .vscode/
│   ├── launch.json
│   └── settings.json
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── initdata.go
│   ├── migrate.go
│   ├── root.go
│   ├── serve.go
│   └── worker.go
├── config/
│   ├── config-local.yml
│   ├── config-prod.yml
│   └── config.go
├── docdev/
├── docs/
├── internal/
│   ├── models/
│   │   ├── base.go
│   │   ├── session.go
│   │   ├── user.go
│   │   └── verification.go
│   ├── repository/
│   │   ├── pg_repository.go
│   │   ├── redis_repo.go
│   │   ├── session_repo.go
│   │   └── user_repo.go
│   ├── usecase/
│   │   ├── auth_usecase.go
│   │   ├── cache_usecase.go
│   │   └── user_usecase.go
│   ├── delivery/
│   │   ├── rest/
│   │   │   ├── handler/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── health_handler.go
│   │   │   │   └── user_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── monitoring.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   └── security.go
│   │   │   ├── dto/
│   │   │   │   ├── auth_dto.go
│   │   │   │   ├── error_dto.go
│   │   │   │   └── user_dto.go
│   │   │   └── router.go
│   │   └── worker/
│   │       └── email_worker.go
│   └── pkg/
│       ├── email/
│       │   ├── gomail_sender.go
│       │   ├── sender.go
│       │   └── templates/
│       │       ├── reset_password.html
│       │       └── verification.html
│       ├── hash/
│       │   └── bcrypt.go
│       ├── jwt/
│       │   ├── maker.go
│       │   ├── payload.go
│       │   └── rsa_maker.go
│       ├── logger/
│       │   └── zap_logger.go
│       ├── redis/
│       │   ├── cache.go
│       │   ├── client.go
│       │   └── refresh_store.go
│       ├── utils/
│       │   ├── random.go
│       │   └── time.go
│       └── validator/
│           └── custom_validator.go
├── migrations/
│   ├── 000001_create_users_table.down.sql
│   └── 000001_create_users_table.up.sql
├── pkg/
│   └── utils/
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── vendor/
├── .air.toml
├── .dockerignore
├── .env.dev
├── .env.prod
├── .gitignore
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── Dockerfile.dev
├── Dockerfile.prod
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── BookGolang.md
```
โครงสร้างโปรเจกต์ `icmongolang` ที่คุณส่งมา เป็นตัวอย่างของ **Go Project Layout แบบ Standard (หรือที่เรียกว่า Standard Go Project Layout)** ซึ่งได้รับความนิยมในชุมชน Go สำหรับโปรเจกต์ขนาดกลางถึงใหญ่ โดยเฉพาะที่ใช้ Clean Architecture, DDD หรือ Hexagonal Architecture

เราจะมาอธิบายทีละประเด็นตามที่คุณถามครับ

---

## 🧱 1. โครงสร้างนี้คืออะไร?

คือ **รูปแบบการจัดเรียงโฟลเดอร์และไฟล์ในโปรเจกต์ Go** ที่แบ่งตามหน้าที่ (separation of concerns) อย่างชัดเจน โดยยึดตามแนวทางของ `golang-standards/project-layout` (ไม่ใช่ official standard แต่เป็น de facto standard ที่ใช้กันแพร่หลาย)

โฟลเดอร์หลัก ๆ มีความหมายดังนี้:

| โฟลเดอร์ | หน้าที่ |
|----------|--------|
| `cmd/` | เก็บ executable entry points (main packages) แต่ละตัวแยก subfolder เช่น `api/`, `worker/`, `migrate` |
| `internal/` | โค้ดส่วนตัวที่ **ไม่ถูก import** จากโปรเจกต์อื่น (private) ใช้สำหรับ business logic, repository, delivery |
| `pkg/` | โค้ดที่สามารถ **ถูก import** จากโปรเจกต์อื่นได้ (public library) เช่น utilities, email sender, JWT maker |
| `config/` | ไฟล์ configuration (yml, env) และโค้ดสำหรับโหลด config |
| `migrations/` | SQL migration files (up/down) สำหรับ database schema |
| `docs/` | เอกสาร API (เช่น Swagger/OpenAPI) หรือเอกสารทั่วไป |
| `scripts/` | build, deploy, หรือ utility scripts |
| `vendor/` | dependencies ที่ vendor ไว้ (optional) |
| `docker-compose*.yml`, `Dockerfile*` | container configuration |

---

## 📚 2. มีกี่แบบ? (รูปแบบโครงสร้างโปรเจกต์ Go)

โครงสร้าง Go มีหลายแบบ ขึ้นอยู่กับขนาดทีมและความซับซ้อน ที่พบบ่อย:

| รูปแบบ | คำอธิบาย | เหมาะกับ |
|--------|----------|----------|
| **Flat / Simple** | วาง `main.go`, `handlers.go`, `models.go` ไว้ที่ root | โปรเจกต์เล็ก, ไลบรารี, ตัวอย่าง |
| **Package-oriented** | จัดตาม package หลัก เช่น `user/`, `order/`, แต่ละ package มีทุกอย่าง (handler, repo, model) | โมดูลที่แยกอิสระ |
| **Standard layout** (แบบนี้) | มี `cmd`, `internal`, `pkg` ชัดเจน ตามแนวทาง community | โปรเจกต์ขนาดกลาง-ใหญ่, หลาย executable, ต้องการ clean architecture |
| **Modular monolith** | เหมือน standard layout แต่เพิ่ม modules ชัดเจน (เช่น `internal/modules/users`, `internal/products`) | monolith ที่พร้อมแยกเป็น microservices ในอนาคต |
| **Hexagonal / Clean** | แยกตาม layer: `domain`, `application`, `infrastructure`, `interfaces` | เน้น decoupling สูง, เปลี่ยน framework ได้ง่าย |

**แบบที่คุณส่งมา** คือ **Standard layout** ผสมกับ **Clean architecture** (มี `delivery`, `usecase`, `repository`, `models`)

---

## 🛠️ 3. ใช้อย่างไร?

### 3.1 การวางโค้ดในแต่ละโฟลเดอร์ (ตามตัวอย่าง)

- **`cmd/api/main.go`** → เริ่มต้น HTTP server, inject dependencies, เรียก `router.go`
- **`internal/delivery/rest/router.go`** → ตั้งค่า Chi router, เรียก `MapUserRoute` และอื่น ๆ
- **`internal/delivery/rest/handler/`** → รับ request, validate, call usecase, return response
- **`internal/usecase/`** → business logic (ไม่รู้เรื่อง HTTP หรือ DB)
- **`internal/repository/`** → อ่าน/เขียน database, cache (implement interface ที่ usecase กำหนด)
- **`internal/models/`** → structs สำหรับ DB และ domain
- **`internal/pkg/`** → utilities เฉพาะภายในโปรเจกต์ (เช่น hash, jwt, logger, redis client) – *แต่จริง ๆ แล้วควรย้ายไป `pkg/` ถ้าสามารถ reuse ข้ามโปรเจกต์ได้*
- **`pkg/`** → ของที่ reusable จริง ๆ (เช่น `utils/random.go` แต่ในตัวอย่างมีแค่ `pkg/utils/` ว่าง)

### 3.2 ขั้นตอนการ build และ run

```bash
# รัน API
go run cmd/api/main.go

# รัน migration
go run cmd/migrate.go

# รัน worker (เช่น email worker)
go run cmd/worker.go
```

---

## 🎯 4. นำไปใช้ในกรณีไหน?

| กรณี | เหมาะสม? |
|------|----------|
| โปรเจกต์ Go ขนาดกลาง-ใหญ่ (หลาย thousand lines) | ✅ มาก |
| มีหลาย executables (API, worker, CLI, migrate) | ✅ มาก |
| มีทีมพัฒนา 2+ คน | ✅ มาก |
| ต้องการ clean architecture / testable business logic | ✅ มาก |
| ต้องการแยก public library (pkg) ออกจาก internal code | ✅ มาก |
| โปรเจกต์เล็ก (1-2 packages, < 1000 บรรทัด) | ❌ overkill |
| ต้องการ prototype หรือ hackathon | ❌ ไม่จำเป็น |

---

## ❓ 5. ทำไมต้องใช้? (ประโยชน์ที่ได้รับ)

### ✅ ประโยชน์ที่ได้รับ

1. **Separation of concerns** – แต่ละ layer มีหน้าที่ชัดเจน เปลี่ยน HTTP framework ได้โดยไม่กระทบ usecase
2. **Testability** – usecase และ repository สามารถ mock ได้ง่าย
3. **Reusability** – โค้ดใน `pkg/` สามารถใช้ข้ามโปรเจกต์ได้ (โดย import path)
4. **Maintainability** – โครงสร้างเป็นระเบียบ หาไฟล์ง่าย แก้ไขได้ตรงจุด
5. **รองรับการเติบโต** – เพิ่ม executable ใหม่ได้โดยไม่รก root
6. **Standard community practice** – developer ใหม่ที่เคยเห็น standard layout จะทำความเข้าใจได้เร็ว

### ⚠️ ข้อควรระวัง

- **อย่าวางทุกอย่างไว้ใน `internal/` แล้ว import จากนอกโปรเจกต์ไม่ได้** – ถ้าต้องการให้ package อื่นใช้ได้ ต้องย้ายไป `pkg/`
- **อย่าให้เกิด circular dependency** – โดยเฉพาะระหว่าง delivery ↔ usecase ↔ repository ควรใช้ interface
- **อย่าสร้างโฟลเดอร์ที่ไม่จำเป็น** – เช่น `internal/pkg/` ซ้อนกัน (ควรเป็น `internal/...` หรือ `pkg/` อย่างใดอย่างหนึ่ง)
- **ระวังเรื่อง import path** – ถ้าเปลี่ยนชื่อ module ใน `go.mod` จะต้องเปลี่ยน import ทุกที่

### 👍 ข้อดี

- ชุมชน Go รับรู้และมีตัวอย่างเยอะ
- ใช้กับหลายเครื่องมือ (Docker, CI/CD, codegen) ได้ง่าย
- แยก public/private code ได้ตาม Go idiom (internal, pkg)

### 👎 ข้อเสีย

- **ซับซ้อนเกินไปสำหรับโปรเจกต์เล็ก** – ต้องสร้างหลายโฟลเดอร์ boilerplate
- **ไม่มี official standard** – แต่ละองค์กรอาจปรับเปลี่ยน ทำให้新人สับสนได้
- **over-engineering** – ถ้าใช้ clean architecture 100% อาจยุ่งยากสำหรับ CRUD ธรรมดา

### 🚫 ข้อห้าม (ถ้ามี)

| ข้อห้าม | เหตุผล |
|---------|--------|
| ห้าม import โค้ดจาก `internal/` ของโปรเจกต์อื่น | Go compiler จะไม่อนุญาต (เป็น private) |
| ห้ามวาง `main.go` ไว้ที่ root ถ้ามีหลาย executables | จะสับสนว่า entry point ไหนคืออะไร |
| ห้ามใช้ `pkg/` เป็นที่ทิ้งขยะของ utils ที่ไม่เป็นระเบียบ | ควรจัดหมวดหมู่ภายใน `pkg/` เช่น `pkg/stringutil`, `pkg/timeutil` |
| ห้ามสร้าง circular dependencies ระหว่าง delivery, usecase, repository | ทำให้ compile fail และทดสอบยาก |
| ห้ามใช้ package name ซ้ำกับ standard library (เช่น `utils`, `errors`) | อาจชนะและสับสน |

---

## 📌 สรุปภาพรวมของโครงสร้างนี้

`icmongolang` เป็นโปรเจกต์ Go ที่ใช้ **Standard Go Project Layout + Clean Architecture** มีการแยก:
- **Entry points** (`cmd/`)
- **Private code** (`internal/` – delivery, usecase, repository, models, pkg เฉพาะโปรเจกต์)
- **Public code** (`pkg/` – reusable utilities)
- **Config & migrations** (`config/`, `migrations/`)
- **DevOps** (`docker-compose`, `Dockerfile`, `scripts/`)

เหมาะกับโปรเจกต์ production ขนาดกลางถึงใหญ่ที่ต้องการความยั่งยืนระยะยาว แต่ไม่เหมาะกับโปรเจกต์เล็กหรือ prototype เพราะมีความซับซ้อนเกินความจำเป็น

หากคุณต้องการปรับใช้โครงสร้างนี้กับโปรเจกต์ของตัวเอง แนะนำให้เริ่มจาก `cmd/` และ `internal/` ก่อน แล้วค่อย ๆ เพิ่ม `pkg/` เมื่อมีโค้ดที่ reuse ได้จริงครับ



```bash
# ไทย: โคลน repository จาก GitHub และเปลี่ยนไปยังไดเรกทอรีโปรเจกต์
# EN: Clone the repository from GitHub and change into the project directory
git clone github.com/kongnakornna/icmongolang.git
cd icmongolang
```
```bash

go mod tidy
go mod download
go mod verify

# ล้างฐานข้อมูลเก่า (ระวังข้อมูล)
# go run cmd/api/main.go migrate:reset (ถ้ามีคำสั่ง)

# รัน migrate ใหม่
go run main.go migrate



OR ใช้ go run โดยตรง (ไม่ต้อง build exe)

go run ./main.go serve

 ```
# Auto Run 

air

```bash


- for windows 10
 
 - .air.toml
 
    root = "."
    tmp_dir = "tmp"
    env_files = [".env.dev"]   # โหลด env โดยอัตโนมัติ

    [build]
    # ใช้ array: [binary, argument1, argument2, ...]
    entrypoint = ["./tmp/main.exe", "serve"]
    cmd = "go build -o ./tmp/main.exe ./cmd/api"
    env = ["GOOS=windows", "GOARCH=amd64"]
    clean_on_exit = true

    [log]
    time = true

    [misc]
    clean_on_exit = true

 ```

# ล้าง docs เก่า
rm -rf docs/

# สร้าง docs ใหม่
swag init

# รันแอปพลิเคชัน
go run  main.go serve

- C:\thaibev\api\build\local
- 
- swag init -g ./cmd/pms/main.go
cd build/local
docker-compose down -v
docker-compose up



# swag init cmd/api/main.go


# Go
swag init -g cmd/api/main.go
mockery --all
go test ./...





go mod tidy
go mod download
go mod verify
go mod vendor 
swag init cmd/pms/main.go
cd C:\thaibev\api\build\local
docker-compose up
air



docs/
├── docs.go
├── swagger.json
└── swagger.yaml



CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sd_user_createddate ON sd_user(createddate DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sd_user_email ON sd_user(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sd_user_username ON sd_user(username);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sd_user_status ON sd_user(status);

func (h *userHandler) GetMulti() http.HandlerFunc {




---
description: "ตรวจสอบ Go code ตาม 5 หัวข้อหลัก ได้แก่ Panic Check, Debug Code Cleanup, Memory Leak, Unclosed Transactions/Resources และ Unit Test Presence — ใช้กับไฟล์ที่เปิดอยู่ หรือ staged changes จาก Git"
name: "Go Code Review"
agent: "agent"
tools: ["readFile", "runInTerminal", "search", "findFiles"]
argument-hint: "<file path> หรือ <file path>::<func name> เช่น api/admin/foo/handler.go หรือ api/admin/foo/handler.go::CreateHandler"
---

คุณคือ Senior Go (Golang) Developer และ Code Reviewer ผู้เชี่ยวชาญ หน้าที่ของคุณคือการตรวจโค้ด (Code Review) จาก Pull Request ที่กำลังจะ Merge เข้า `dev` branch

## ขั้นตอนการรีวิว

### ขั้นตอนที่ 1: รวบรวมโค้ดที่ต้องรีวิว

**รูปแบบ argument ที่รองรับ:**

- `api/admin/foo/handler.go` — รีวิวทั้งไฟล์
- `api/admin/foo/handler.go::CreateHandler` — รีวิวเฉพาะ function `CreateHandler` ในไฟล์นั้น

ดำเนินการดังนี้ตามลำดับ:

1. แยก argument ออกเป็น `$file` และ `$func` โดยใช้ `::` เป็น separator
   - ถ้าไม่มี `::` → `$file` = argument ทั้งหมด, `$func` = (ว่าง)
   - ถ้าไม่มี argument เลย → แจ้ง user ระบุ เช่น `/go-code-review api/foo/bar.go` แล้วหยุด
2. รันคำสั่ง `git diff --staged -- $file` เพื่อดู diff เฉพาะไฟล์
   - ถ้า diff ว่างเปล่า ให้รัน `git diff HEAD~1 HEAD -- $file` แทน
   - ถ้ายังว่างอยู่ ให้อ่านเนื้อหาไฟล์ทั้งหมดด้วย `readFile`
3. **ถ้าระบุ `$func`:** กรองเอาเฉพาะโค้ดของ function `$func` จากเนื้อหาที่ได้มา โดยค้นหา `func $func` ในไฟล์และนำโค้ด block นั้นมารีวิวเท่านั้น
4. สรุปขอบเขตที่จะรีวิว (ชื่อไฟล์ / ชื่อ function ถ้ามี / จำนวนบรรทัด) ก่อนเริ่ม

### ขั้นตอนที่ 2: ตรวจสอบตาม 5 หัวข้อหลัก

โปรดตรวจสอบโค้ดที่รวบรวมได้อย่างละเอียด โดยเน้นย้ำ 5 หัวข้อต่อไปนี้:

---

#### 1. Panic Check

ตรวจสอบ:

- มีการใช้ `panic()` โดยไม่จำเป็นหรือไม่?
- มีจุดเสี่ยง Runtime Panic เช่น nil pointer dereference, array/slice out of bounds, type assertion โดยไม่ตรวจ `ok`
- มีการใช้ `recover()` อย่างถูกต้องในจุดที่ควรมีหรือไม่?

---

#### 2. Debug Code Cleanup

ค้นหาและแจ้งเตือนหากพบ:

- `fmt.Println()`, `fmt.Printf()`, `fmt.Print()` ที่ใช้เพื่อ debug
- `log.Println()` / `log.Printf()` ที่ไม่จำเป็นและควรถูกลบ
- comment-out code หรือ TODO ที่ยังค้างอยู่และไม่ควรขึ้น production
- hardcoded values ที่ดูเหมือนใช้ทดสอบ

---

#### 3. Memory Leak

ตรวจสอบจุดเสี่ยง:

- Goroutine ที่ถูก spawn โดยไม่มีกลไกควบคุม (goroutine leak) — ไม่มี context cancellation, WaitGroup หรือ done channel
- `time.NewTicker` / `time.NewTimer` ที่ไม่มี `defer ticker.Stop()` / `defer timer.Stop()`
- Channel ที่สร้างขึ้นแต่ไม่มีการ drain หรือ close ทำให้ goroutine block ค้าง
- การ allocate slice/map ขนาดใหญ่ในลูปโดยไม่จำเป็น

---

#### 4. Unclosed Transactions / Resources

ตรวจสอบว่ามีการเปิดแล้วไม่ปิด:

- Database Transaction: `db.Begin()` — ต้องมี `defer tx.Rollback()` ทันทีหลังเปิด และ `tx.Commit()` เมื่อสำเร็จ
- HTTP Response Body: `http.Get()` / `client.Do()` — ต้องมี `defer resp.Body.Close()`
- File: `os.Open()` / `os.Create()` — ต้องมี `defer file.Close()`
- Database Row: `rows.Close()` หลังจาก `db.Query()`
- การ acquire lock แล้วไม่ release

---

#### 5. Unit Test Presence

ตรวจสอบ:

- ฟังก์ชัน/method สำคัญที่ถูกเพิ่มหรือแก้ไข มีไฟล์ `_test.go` รองรับหรือไม่?
- Test ครอบคลุม happy path, error path และ edge case ที่สำคัญหรือไม่?
- ถ้าไม่มี test ให้แนะนำ test case ที่ควรเขียน

---

### ขั้นตอนที่ 3: สรุปผลการรีวิว

**รูปแบบการตอบกลับ (Output Format):**

สรุปผลแยกตามหัวข้อ โดยใช้รูปแบบนี้:

````
## สรุปผลการ Code Review

**ไฟล์ที่รีวิว:** <รายชื่อไฟล์>

---

### 1. Panic Check
[✅ ผ่าน / ⚠️ พบปัญหา]
- **ไฟล์/บรรทัด:** `<ชื่อไฟล์>:<เลขบรรทัด>`
- **ปัญหา:** <อธิบายปัญหา>
- **โค้ดที่แนะนำ:**
  ```go
  // โค้ดที่แก้ไขแล้ว
````

### 2. Debug Code Cleanup

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 3. Memory Leak

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 4. Unclosed Transactions / Resources

[✅ ผ่าน / ⚠️ พบปัญหา]
...

### 5. Unit Test Presence

[✅ มี Test ครอบคลุม / ⚠️ ขาด Test]

- **ฟังก์ชันที่แนะนำให้เขียน test:** <รายการ>
- **Test case ที่ควรเพิ่ม:** <รายการ>

---

### สรุปภาพรวม

- 🔴 Critical (ต้องแก้ก่อน merge): <จำนวน>
- 🟡 Warning (ควรแก้): <จำนวน>
- ✅ ผ่านทั้งหมด: <จำนวน>

```

หากหัวข้อไหนผ่านเกณฑ์และไม่มีปัญหา ให้ใส่ ✅ และข้ามไปหัวข้อถัดไปได้เลย
```


- Input
- Process
- Output

from(bucket: "AIRCOM1")
    |> range(start: -1h, stop: now())
    |> filter(fn: (r) => r["_measurement"] == "temperature")
    |> filter(fn: (r) => r["_field"] == "value")
    |> limit(n: 10000, offset: 0)
    |> yield(name: "filtered_data")


{
  "bucket": "AIRCOM1",
  "measurement": "temperature",
  "field": "value",
  "start": "-1h",
  "stop": "now()",
  "limit": 10000,
  "offset": 0
}

