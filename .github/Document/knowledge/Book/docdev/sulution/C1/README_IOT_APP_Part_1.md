ขอรับคำขอของท่านแล้ว ซึ่งเป็นโครงการที่มีรายละเอียดสูงมาก ครอบคลุมทั้งการออกแบบระบบ API ด้วย GoLang, สถาปัตยกรรม 3 ชั้น, การทำงานของระบบ IoT และระบบแจ้งเตือนแบบ Real-time รวมถึงการจัดทำเอกสารเป็นเล่มต่างๆ

เนื่องจากข้อจำกัดด้านความยาวในการตอบในรอบเดียว กระผมขอเสนอแผนการดำเนินการดังนี้:

1. **เริ่มต้นที่ "เล่ม 2: สถาปัตยกรรมโครงสร้างระบบ"** ซึ่งเป็นหัวใจหลักของโครงสร้างโฟลเดอร์และการออกแบบ Data Flow ก่อน
2. ในแต่ละบท จะทำให้สมบูรณ์ตามข้อกำหนดทั้งหมด (Concept, Code, Table, Exercise, References)
3. จะแทรก Mermaid Flowchart, ตาราง Checklist, และตัวอย่าง Code ที่รันได้จริง
4. ใช้ภาษาไทยเป็นหลัก และอังกฤษในคำศัพท์เทคนิคหรือ comment ใน code

**เริ่มด้วยบทที่ 1 ของเล่ม 2: โครงสร้างพื้นฐานของโปรเจกต์ (Project Structure Foundation)**

---

# เล่ม 2: สถาปัตยกรรมโครงสร้างระบบ (System Architecture)
## บทที่ 1: โครงสร้างพื้นฐานของโปรเจกต์และการออกแบบ Data Flow

### สรุปสั้นก่อนเริ่ม
โครงสร้างโฟลเดอร์ `gobackend` ที่กำหนดไว้เป็นรูปแบบ **Clean Architecture 3 ชั้น** (Repository – Usecase – Delivery) เหมาะสำหรับระบบผลิตจริง (Production-ready) รองรับการแยก concerns, ทดสอบง่าย และขยายฟังก์ชันได้โดยไม่แก้ไขโค้ดเดิมมากนัก บทนี้จะอธิบายที่มา, หน้าที่ของแต่ละโฟลเดอร์, และการไหลของข้อมูล (Data Flow) จาก HTTP request จนถึงฐานข้อมูล พร้อมตัวอย่างการนำไปใช้งานจริง

---

### คำอธิบายแนวคิด (Concept Explanation)

#### Clean Architecture คืออะไร?
Clean Architecture (Robert C. Martin) เป็นแนวคิดที่แบ่งซอฟต์แวร์ออกเป็นวงกลมซ้อนกัน โดยชั้นในสุดคือ **Entities (Models)** และชั้นนอกสุดคือ **Delivery (HTTP handlers, workers)** กฎสำคัญ: dependencies ชี้เข้าหาศูนย์กลางเท่านั้น (ชั้นนอกขึ้นอยู่กับชั้นในได้ แต่ชั้นในไม่รู้จักชั้นนอก)

#### ทำไมต้องใช้三层架构 (3-layer) แบบ Repository-Usecase-Delivery?
| Layer | หน้าที่ | ตัวอย่างใน Go |
|-------|--------|----------------|
| **Repository** | ติดต่อฐานข้อมูล / แหล่งข้อมูลภายนอก (PostgreSQL, Redis, API อื่น) | `user_repo.go` มี method `Create()`, `FindByID()` |
| **Usecase** | ธุรกิจลอจิก (business logic) ตรวจสอบ rules, ประสาน repo ต่างๆ | `auth_usecase.go` ตรวจสอบ password, สร้าง JWT |
| **Delivery** | รับ input จาก user (HTTP body, query string) และส่ง output กลับ | `auth_handler.go` ผูก JSON, เรียก usecase |

**ประโยชน์ที่ได้รับ**
- แยกความรับผิดชอบ (Separation of Concerns) – แต่ละ layer มีหน้าที่เดียว
- ทดสอบง่าย (Testable) – mock repository ทดสอบ usecase โดยไม่ต้องมี DB จริง
- เปลี่ยนเทคโนโลยีได้โดยไม่กระทบ layer อื่น เช่น เปลี่ยน GORM เป็น SQLx แก้แค่ repository

**ข้อควรระวัง**
- อย่าให้ usecase ส่ง `*gorm.DB` หรือ HTTP status code กลับมา – usecase ต้องคืนค่า business error เท่านั้น
- อย่าให้ handler เรียก repository โดยตรง – ต้องผ่าน usecase เสมอ

**ข้อดี**
- โค้ดเป็นระเบียบ อ่านง่าย
- ขยายฟีเจอร์ใหม่โดยเพิ่ม usecase/handler ใหม่ ไม่ต้องแก้ของเก่า

**ข้อเสีย**
- จำนวนไฟล์เยอะกว่า monolithic
- มือใหม่อาจงงกับการส่งต่อ dependencies

**ข้อห้าม**
- ห้ามใช้ `context.Context` ใน repository เพื่อ logging? ควรใช้ context เพื่อ tracing, timeout เสมอ
- ห้าม import delivery ใน usecase

---

### โครงสร้างการทำงาน (Workflow) ของแต่ละโฟลเดอร์

#### โฟลเดอร์หลักใน `gobackend` และหน้าที่

| โฟลเดอร์/ไฟล์ | หน้าที่ |
|----------------|---------|
| `cmd/` | จุดเริ่มต้นของแอปพลิเคชัน (entry points) – แต่ละไฟล์สร้าง CLI command ด้วย cobra |
| `config/` | โหลด config จาก YAML/env, struct `Config` |
| `internal/` | โค้ด private ของโปรเจกต์ (ไม่อนุญาตให้ import โดยโปรเจกต์อื่น) |
| `internal/models/` | Entity (struct) ที่สอดคล้องกับตาราง DB และ DTO บางตัว |
| `internal/repository/` | อินเทอร์เฟซและอิมพลีเมนต์การติดต่อ DB/Cache |
| `internal/usecase/` | อินเทอร์เฟซและอิมพลีเมนต์ business logic |
| `internal/delivery/rest/` | HTTP handlers, middleware, DTO, router |
| `internal/pkg/` | ส่วนประกอบภายในที่ reuse ได้ (email, jwt, logger, redis) |
| `migrations/` | SQL migration files (up/down) |
| `scripts/` | build, deploy scripts |
| `docker-compose.*.yml` | สำหรับ dev และ prod environment |

#### Data Flow แบบละเอียด (จาก HTTP Request → DB → Response)

```
Client Request (POST /api/login)
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. router.go (chi.Router)                                   │
│    - Route: /api/login -> auth_handler.Login               │
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Middleware chain                                         │
│    - CORS, Logger, RateLimit, Recovery                     │
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. auth_handler.Login (Delivery)                           │
│    - Bind request body to LoginRequest DTO                 │
│    - Validate (email, password required)                   │
│    - Call authUsecase.Login(ctx, dto)                      │
│    - Map result to HTTP status (200, 401, 500)             │
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. auth_usecase.Login (Usecase)                            │
│    - Validate business rules: user exists? password match? │
│    - Call userRepo.FindByEmail(email)                      │
│    - Call jwtMaker.CreateToken(userID, role)               │
│    - Call sessionRepo.Create(session) (Redis)              │
│    - Return token response or business error               │
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. userRepo.FindByEmail (Repository)                       │
│    - Execute SQL: SELECT * FROM users WHERE email=$1      │
│    - Scan row into models.User                             │
│    - Return user or sql.ErrNoRows                          │
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Database (PostgreSQL)                                   │
└─────────────────────────────────────────────────────────────┘
   │
   ▼ (กลับขึ้นไป)
   Response 200 JSON { access_token, refresh_token }
```

#### Flowchart แสดง Data Flow (ใช้ Mermaid)

```mermaid
flowchart TB
    Client[Client HTTP Request] --> Router[chi Router]
    Router --> Middleware[Middleware: Logger, RateLimit, CORS]
    Middleware --> Handler[Delivery: auth_handler.Login]
    Handler --> Bind[Bind & Validate DTO]
    Bind --> Usecase[Usecase: auth_usecase.Login]
    Usecase --> RepoCall[Repository: userRepo.FindByEmail]
    RepoCall --> DB[(PostgreSQL)]
    DB --> RepoReturn[Return models.User]
    RepoReturn --> BusinessCheck{Password valid?}
    BusinessCheck -- No --> ErrorResp[Return 401 Unauthorized]
    BusinessCheck -- Yes --> CreateJWT[jwtMaker.CreateToken]
    CreateJWT --> StoreSession[Redis: sessionRepo.Create]
    StoreSession --> SuccessResp[Return 200 + tokens]
    ErrorResp --> Client
    SuccessResp --> Client
```

**รูปที่ 1:** แผนภาพการไหลของข้อมูลสำหรับ endpoint `/api/login` แสดงการทำงานตั้งแต่ request เข้าสู่ router, middleware, handler, usecase, repository, database และตอบกลับ

---

### ตัวอย่างโค้ดที่รันได้จริง (Runnable Code Example)

เราจะสร้างตัวอย่างเล็กๆ ที่สอดคล้องกับโครงสร้างข้างต้น – **ระบบสมัครสมาชิก (Register)** แบบไม่มี JWT ซับซ้อน เพื่อให้เห็นภาพการทำงานของทั้ง 3 layer

#### ขั้นตอนที่ 1: โครงสร้างโปรเจกต์ตัวอย่าง (ย่อส่วน)
```bash
gobackend-demo/
├── cmd/api/main.go
├── config/config.go
├── internal/
│   ├── models/user.go
│   ├── repository/user_repo.go
│   ├── usecase/user_usecase.go
│   └── delivery/rest/
│       ├── handler/user_handler.go
│       └── router.go
├── go.mod
└── .env
```

#### ไฟล์ `go.mod`
```go
module gobackend-demo

go 1.21

require (
    github.com/go-chi/chi/v5 v5.0.10
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/postgres v1.5.4
    gorm.io/gorm v1.25.5
)
```

#### ไฟล์ `config/config.go`
```go
package config

import (
    "os"
    "log"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
}

func Load() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file, using system env")
    }
    return &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBPort:     os.Getenv("DB_PORT"),
    }
}
```

#### ไฟล์ `internal/models/user.go`
```go
package models

import "time"

// User represents the user entity in database
// ใช้ GORM tags สำหรับ mapping กับ PostgreSQL
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"` // hashed
    FullName  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### ไฟล์ `internal/repository/user_repo.go`
```go
package repository

import (
    "gobackend-demo/internal/models"
    "gorm.io/gorm"
    "context"
)

// UserRepository interface defines database operations for users
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct {
    db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepo{db: db}
}

// Create inserts a new user into database
// ฟังก์ชันนี้จะบันทึกผู้ใช้ใหม่ลงฐานข้อมูล
func (r *userRepo) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail retrieves a user by email
// ค้นหาผู้ใช้จากอีเมล
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

#### ไฟล์ `internal/usecase/user_usecase.go`
```go
package usecase

import (
    "context"
    "errors"
    "gobackend-demo/internal/models"
    "gobackend-demo/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

// UserUsecase handles business logic for user operations
type UserUsecase interface {
    Register(ctx context.Context, email, password, fullName string) error
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

// Register validates and creates a new user
// ตรวจสอบว่า email ซ้ำหรือไม่, hash password, แล้วบันทึก
func (u *userUsecase) Register(ctx context.Context, email, password, fullName string) error {
    // Check if user already exists
    existing, _ := u.userRepo.FindByEmail(ctx, email)
    if existing != nil {
        return errors.New("email already registered")
    }

    // Hash password with bcrypt
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return errors.New("failed to hash password")
    }

    user := &models.User{
        Email:    email,
        Password: string(hashed),
        FullName: fullName,
    }

    return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return u.userRepo.FindByEmail(ctx, email)
}
```

#### ไฟล์ `internal/delivery/rest/handler/user_handler.go`
```go
package handler

import (
    "encoding/json"
    "net/http"
    "gobackend-demo/internal/usecase"
)

type UserHandler struct {
    userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: userUsecase}
}

// RegisterRequest DTO (Data Transfer Object)
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    FullName string `json:"full_name"`
}

// Register handles POST /register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Basic validation (in real project use validator)
    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password required", http.StatusBadRequest)
        return
    }

    err := h.userUsecase.Register(r.Context(), req.Email, req.Password, req.FullName)
    if err != nil {
        if err.Error() == "email already registered" {
            http.Error(w, err.Error(), http.StatusConflict)
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
}
```

#### ไฟล์ `internal/delivery/rest/router.go`
```go
package rest

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "gobackend-demo/internal/delivery/rest/handler"
)

func SetupRouter(userHandler *handler.UserHandler) *chi.Mux {
    r := chi.NewRouter()
    r.Use(middleware.Logger)   // log requests
    r.Use(middleware.Recoverer) // recover from panics

    r.Post("/register", userHandler.Register)

    return r
}
```

#### ไฟล์ `cmd/api/main.go`
```go
package main

import (
    "log"
    "gobackend-demo/config"
    "gobackend-demo/internal/delivery/rest"
    "gobackend-demo/internal/delivery/rest/handler"
    "gobackend-demo/internal/repository"
    "gobackend-demo/internal/usecase"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    cfg := config.Load()

    // Connect to PostgreSQL
    dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword +
           " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Auto migrate schema (for demo)
    db.AutoMigrate(&models.User{})

    // Initialize layers
    userRepo := repository.NewUserRepository(db)
    userUsecase := usecase.NewUserUsecase(userRepo)
    userHandler := handler.NewUserHandler(userUsecase)

    router := rest.SetupRouter(userHandler)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

#### วิธีการรันจริง
1. ติดตั้ง Go 1.21+, PostgreSQL
2. สร้าง database `demo_db`
3. ตั้ง environment variables: `DB_HOST=localhost DB_USER=postgres DB_PASSWORD=xxx DB_NAME=demo_db DB_PORT=5432`
4. `go mod tidy`
5. `go run cmd/api/main.go`
6. ทดสอบ `curl -X POST http://localhost:8080/register -d '{"email":"test@example.com","password":"123456","full_name":"Test User"}' -H "Content-Type: application/json"`

---

### ตารางสรุปการทำงานของแต่ละ Layer

| Layer | รับผิดชอบ | ตัวอย่าง Method | สิ่งที่ควรคืน | สิ่งที่ห้ามทำ |
|-------|-----------|----------------|--------------|----------------|
| **Delivery** | แปลง HTTP request/response, validation ระดับ input | `Register(w,r)` | HTTP status, JSON | ไม่ควรมี SQL หรือ business logic |
| **Usecase** | Business rules, orchestrating repos | `Register(ctx,email,pass)` | business error หรือ result | ไม่ควรรู้ว่าเป็น HTTP |
| **Repository** | Database/Cache operations | `Create(ctx, user)` | DB model หรือ error | ไม่ควรมี business logic |

---

### แบบฝึกหัดท้ายบท (3 ข้อ)

1. **จงเพิ่ม method `UpdateUser` ใน repository และ usecase** ที่อนุญาตให้แก้ไข `full_name` เท่านั้น โดยต้องตรวจสอบว่าผู้ใช้มีอยู่ในระบบก่อน
2. **เขียน middleware ตัวหนึ่ง** (ชื่อ `AuthMiddleware`) ที่อ่าน header `Authorization: Bearer <token>` แล้วตรวจสอบ token (สมมติว่ามี JWT service) ถ้าไม่ถูกต้องให้ตอบ 401 โดยไม่ถึง handler
3. **จากตัวอย่าง Data Flow Diagram** จงปรับปรุงให้รองรับการแคชผู้ใช้ใน Redis โดยเพิ่ม cache layer ระหว่าง usecase และ repository (อธิบายเป็นลำดับขั้นตอน)

---

### แหล่งอ้างอิง (References)

- Clean Architecture by Robert C. Martin – [https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- Go Chi router examples – [https://github.com/go-chi/chi](https://github.com/go-chi/chi)
- GORM documentation – [https://gorm.io/docs/](https://gorm.io/docs/)
- การจัดโครงสร้างโปรเจกต์ Go – [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

---

**หมายเหตุ:** เนื่องจากข้อความตอบต้องไม่ยาวเกินไป กระผมได้ส่งเฉพาะ **บทที่ 1 ของเล่ม 2** เท่านั้น หากท่านต้องการให้ทำ **เล่ม 3 (การพัฒนาโค้ด)** หรือ **เล่ม 4 (การปรับใช้และการบำรุงรักษา)** หรือต้องการให้ขยายความในส่วนของ **ระบบ IoT และ MQTT** ตาม requirement ที่ให้ไว้ (CMON IoT Solution) โปรดแจ้ง และผมจะส่งเนื้อหาบทถัดไปให้ทันที พร้อมทั้งแนบเทมเพลต Task List, Checklist และ Timeline Project ให้สมบูรณ์ครบถ้วน

**ท่านต้องการให้ดำเนินการต่อในส่วนใดก่อน?**  
1. เล่ม 3: การพัฒนา Authentication (JWT, Refresh Token)  
2. เล่ม 3: การเชื่อมต่อ MQTT และ Real-time Alert  
3. เล่ม 4: การ Deploy ด้วย Docker Compose และ Monitoring  
4. เทมเพลต Task List / Checklist สำหรับโครงการนี้โดยเฉพาะ