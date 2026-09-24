## 📚 เล่มที่ 10: เทมเพลต กระบวนการพัฒนา และตัวอย่างโค้ด + ภาคผนวก

เล่มนี้เป็นเล่มสุดท้ายที่จะรวบรวมทุกอย่างเข้าด้วยกัน เริ่มจากตัวอย่างโปรเจกต์ที่สมบูรณ์ที่คุณสามารถนำไปใช้เป็นต้นแบบได้จริง, เทมเพลตและ Checklist ที่ช่วยให้การพัฒนาของคุณเป็นระบบมากขึ้น, แผนภาพที่ช่วยให้เข้าใจสถาปัตยกรรม, การจัดการ Configuration ที่ยืดหยุ่น, และภาคผนวกที่เต็มไปด้วยความรู้เสริมที่ใช้งานได้จริง

---

### บทที่ 59: ตัวอย่างโค้ดครบวงจร (Full-stack Example)

ในบทนี้เราจะสร้าง **REST API สำหรับระบบจัดการสินค้า (Product Management System)** ที่สมบูรณ์ พร้อมเชื่อมต่อฐานข้อมูล, Authentication, Validation, การทดสอบ, และ Docker

#### โครงสร้างโปรเจกต์:
```
product-api/
├── cmd/
│   └── api/
│       └── main.go                 # จุดเริ่มต้น
├── internal/
│   ├── domain/
│   │   ├── product.go              # Entity Product
│   │   └── repository.go           # Repository Interface
│   ├── usecase/
│   │   └── product/
│   │       ├── interface.go        # Input/Output Ports
│   │       ├── create.go
│   │       ├── get.go
│   │       ├── update.go
│   │       └── delete.go
│   ├── handler/
│   │   ├── product_handler.go
│   │   └── dto/
│   │       └── product_dto.go
│   ├── repository/
│   │   ├── product_repo.go
│   │   └── db.go
│   └── config/
│       └── config.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── middleware/
│       ├── auth.go
│       ├── logging.go
│       └── recovery.go
├── test/
│   └── integration/
│       └── product_test.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

#### 1. domain/product.go
```go
package domain

import (
    "errors"
    "time"
)

type Product struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(name, description string, price float64, stock int) (*Product, error) {
    if name == "" {
        return nil, errors.New("ชื่อสินค้า不能为空")
    }
    if price < 0 {
        return nil, errors.New("ราคาต้องมากกว่า 0")
    }
    if stock < 0 {
        return nil, errors.New("สต็อกต้องมากกว่า 0")
    }
    return &Product{
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}

func (p *Product) UpdatePrice(newPrice float64) error {
    if newPrice < 0 {
        return errors.New("ราคาต้องมากกว่า 0")
    }
    p.Price = newPrice
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) UpdateStock(quantity int) error {
    if p.Stock+quantity < 0 {
        return errors.New("สต็อกไม่พอ")
    }
    p.Stock += quantity
    p.UpdatedAt = time.Now()
    return nil
}
```

#### 2. domain/repository.go
```go
package domain

type ProductRepository interface {
    Save(product *Product) error
    FindByID(id int64) (*Product, error)
    FindAll(page, limit int) ([]Product, int64, error)
    Update(product *Product) error
    Delete(id int64) error
}
```

#### 3. usecase/product/interface.go
```go
package product

type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type CreateProductResponse struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type GetProductResponse struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type ListProductResponse struct {
    Products []GetProductResponse `json:"products"`
    Total    int64                `json:"total"`
    Page     int                  `json:"page"`
    Limit    int                  `json:"limit"`
}

// Use Case Interfaces
type CreateProductUseCase interface {
    Execute(req CreateProductRequest) (*CreateProductResponse, error)
}

type GetProductUseCase interface {
    Execute(id int64) (*GetProductResponse, error)
}

type ListProductsUseCase interface {
    Execute(page, limit int) (*ListProductResponse, error)
}

type UpdateProductUseCase interface {
    Execute(id int64, req CreateProductRequest) (*CreateProductResponse, error)
}

type DeleteProductUseCase interface {
    Execute(id int64) error
}
```

#### 4. usecase/product/create.go
```go
package product

import (
    "myproject/internal/domain"
)

type createProductUseCase struct {
    repo domain.ProductRepository
}

func NewCreateProductUseCase(repo domain.ProductRepository) CreateProductUseCase {
    return &createProductUseCase{repo: repo}
}

func (uc *createProductUseCase) Execute(req CreateProductRequest) (*CreateProductResponse, error) {
    product, err := domain.NewProduct(req.Name, req.Description, req.Price, req.Stock)
    if err != nil {
        return nil, err
    }

    if err := uc.repo.Save(product); err != nil {
        return nil, err
    }

    return &CreateProductResponse{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
    }, nil
}
```

#### 5. handler/product_handler.go
```go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "myproject/internal/usecase/product"
)

type ProductHandler struct {
    createUC product.CreateProductUseCase
    getUC    product.GetProductUseCase
    listUC   product.ListProductsUseCase
    updateUC product.UpdateProductUseCase
    deleteUC product.DeleteProductUseCase
}

func NewProductHandler(
    createUC product.CreateProductUseCase,
    getUC product.GetProductUseCase,
    listUC product.ListProductsUseCase,
    updateUC product.UpdateProductUseCase,
    deleteUC product.DeleteProductUseCase,
) *ProductHandler {
    return &ProductHandler{
        createUC: createUC,
        getUC:    getUC,
        listUC:   listUC,
        updateUC: updateUC,
        deleteUC: deleteUC,
    }
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var req product.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    resp, err := h.createUC.Execute(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    resp, err := h.getUC.Execute(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 10
    }

    resp, err := h.listUC.Execute(page, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

#### 6. main.go (Entrypoint)
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "myproject/internal/config"
    "myproject/internal/handler"
    "myproject/internal/repository"
    "myproject/internal/usecase/product"
    "myproject/pkg/logger"
    "myproject/pkg/middleware"
)

func main() {
    // โหลด Config
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Load config error:", err)
    }

    // ตั้งค่า Logger
    logger.Init(cfg.LogLevel)

    // เชื่อมต่อฐานข้อมูล
    db, err := sql.Open("mysql", cfg.Database.DSN())
    if err != nil {
        logger.Fatal("Connect database error:", err)
    }
    defer db.Close()
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // สร้าง Repositories
    productRepo := repository.NewProductRepository(db)

    // สร้าง Use Cases
    createUC := product.NewCreateProductUseCase(productRepo)
    getUC := product.NewGetProductUseCase(productRepo)
    listUC := product.NewListProductsUseCase(productRepo)
    updateUC := product.NewUpdateProductUseCase(productRepo)
    deleteUC := product.NewDeleteProductUseCase(productRepo)

    // สร้าง Handlers
    productHandler := handler.NewProductHandler(createUC, getUC, listUC, updateUC, deleteUC)

    // ตั้งค่า Router
    mux := http.NewServeMux()
    mux.HandleFunc("POST /api/products", productHandler.CreateProduct)
    mux.HandleFunc("GET /api/products", productHandler.ListProducts)
    mux.HandleFunc("GET /api/products/{id}", productHandler.GetProduct)
    mux.HandleFunc("PUT /api/products/{id}", productHandler.UpdateProduct)
    mux.HandleFunc("DELETE /api/products/{id}", productHandler.DeleteProduct)

    // ใช้ Middleware
    handler := middleware.Recovery(
        middleware.Logging(
            middleware.Auth(mux),
        ),
    )

    // สร้าง HTTP Server
    srv := &http.Server{
        Addr:         cfg.Server.Addr,
        Handler:      handler,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // เริ่ม Server
    go func() {
        logger.Infof("Server เริ่มที่ %s", cfg.Server.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server error:", err)
        }
    }()

    // Graceful Shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown:", err)
    }
    logger.Info("Server shutdown complete")
}
```

---

### บทที่ 60: Task List Template (เทมเพลตรายการงาน)

เทมเพลตนี้ช่วยให้ทีมพัฒนาจัดการงานและติดตามความคืบหน้าอย่างเป็นระบบ

```markdown
# 📋 Task List Template

## Project: [ชื่อโปรเจกต์]
## Sprint: [Sprint X]
## Date: [วันที่เริ่ม - สิ้นสุด]

---

### ✅ Completed Tasks
- [x] Task 1: [คำอธิบาย] - Assignee: [ชื่อ] - Done: [วันที่]
- [x] Task 2: [คำอธิบาย] - Assignee: [ชื่อ] - Done: [วันที่]

---

### 🚧 In Progress
- [ ] Task 3: [คำอธิบาย] - Assignee: [ชื่อ] - Est: [เวลา] - Start: [วันที่]
- [ ] Task 4: [คำอธิบาย] - Assignee: [ชื่อ] - Est: [เวลา] - Start: [วันที่]

---

### 📝 To Do
- [ ] Task 5: [คำอธิบาย] - Assignee: [ชื่อ] - Priority: [High/Medium/Low]
- [ ] Task 6: [คำอธิบาย] - Assignee: [ชื่อ] - Priority: [High/Medium/Low]

---

### 🐛 Bugs
- [ ] Bug 1: [คำอธิบาย] - Severity: [Critical/High/Medium/Low] - Assignee: [ชื่อ]
- [ ] Bug 2: [คำอธิบาย] - Severity: [High/Medium/Low] - Assignee: [ชื่อ]

---

### 📊 Blockers
- [ ] Blocker 1: [คำอธิบาย] - Affects: [Task ID] - Resolution: [Pending/In Progress]

---

### 📈 Notes / Comments
- ความเสี่ยง: ...
- การตัดสินใจ: ...
- สิ่งที่ต้องเรียนรู้เพิ่มเติม: ...
```

**Go Implementation สำหรับจัดการ Tasks:**
```go
package task

type TaskStatus string

const (
    StatusTodo       TaskStatus = "todo"
    StatusInProgress TaskStatus = "in_progress"
    StatusDone       TaskStatus = "done"
)

type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
)

type Task struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    Priority    Priority   `json:"priority"`
    Assignee    string     `json:"assignee"`
    Estimated   time.Duration `json:"estimated"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DoneAt      *time.Time `json:"done_at,omitempty"`
}

type TaskList struct {
    Tasks []Task `json:"tasks"`
}

func (tl *TaskList) AddTask(task Task) {
    tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) GetTasksByStatus(status TaskStatus) []Task {
    result := []Task{}
    for _, t := range tl.Tasks {
        if t.Status == status {
            result = append(result, t)
        }
    }
    return result
}

func (tl *TaskList) GetTasksByAssignee(assignee string) []Task {
    result := []Task{}
    for _, t := range tl.Tasks {
        if t.Assignee == assignee {
            result = append(result, t)
        }
    }
    return result
}
```

---

### บทที่ 61: Checklist Template

Checklist สำหรับการเตรียมตัวก่อน Deploy, Code Review, และการเริ่มต้นโปรเจกต์ใหม่

#### 61.1 Pre-Deployment Checklist
```markdown
# 🚀 Pre-Deployment Checklist

## Before Deploy
- [ ] รัน Test ทั้งหมด: `go test ./...`
- [ ] รัน Benchmark: `go test -bench=.`
- [ ] ตรวจสอบ Code Coverage: `go test -cover`
- [ ] รัน Linter: `golangci-lint run`
- [ ] ตรวจสอบ Vulnerability: `go vet ./...`
- [ ] อัปเดต Documentation (API Docs, README)
- [ ] ทดสอบ Integration/E2E
- [ ] ตรวจสอบ Configuration (Production vs Staging)
- [ ] ตรวจสอบ Environment Variables
- [ ] ตรวจสอบ Migrations (ขึ้น-ลงได้)

## During Deploy
- [ ] Backup Database
- [ ] ตรวจสอบ Health Check Endpoint
- [ ] ตรวจสอบ Graceful Shutdown
- [ ] ตรวจสอบ Logging (ระดับที่เหมาะสม)
- [ ] ตรวจสอบ Metrics/Monitoring
- [ ] ตรวจสอบ Alert Configuration

## After Deploy
- [ ] ทดสอบ Smoke Test (API หลัก)
- [ ] ตรวจสอบ Logs (ไม่มี Error)
- [ ] ตรวจสอบ Performance (Response Time, Memory)
- [ ] แจ้งทีมผ่าน Discord/Slack
- [ ] อัปเดต Status Page
- [ ] ตั้งค่า Rollback Plan (ถ้าจำเป็น)
```

#### 61.2 Code Review Checklist
```markdown
# 📝 Code Review Checklist

## Code Quality
- [ ] โค้ดอ่านง่ายและเข้าใจได้
- [ ] ตั้งชื่อตัวแปร/ฟังก์ชันมีความหมาย
- [ ] ไม่มีโค้ดที่ซ้ำกัน (DRY)
- [ ] ใช้ Standard Library อย่างเหมาะสม
- [ ] หลีกเลี่ยง Global Variables

## Error Handling
- [ ] จัดการ Error ทุกกรณี
- [ ] ใช้ `%w` สำหรับ Wrap Error
- [ ] ไม่กลืน Error (Don't swallow errors)
- [ ] ใช้ Custom Error Types อย่างเหมาะสม

## Concurrency
- [ ] ใช้ Goroutine อย่างเหมาะสม
- [ ] ป้องกัน Data Race (ใช้ Mutex/Channels)
- [ ] ใช้ WaitGroup เมื่อจำเป็น
- [ ] หลีกเลี่ยง Goroutine Leak

## Performance
- [ ] หลีกเลี่ยงการจัดสรรหน่วยความจำโดยไม่จำเป็น
- [ ] ใช้ Pool (sync.Pool) สำหรับวัตถุที่ใช้บ่อย
- [ ] ใช้ Context สำหรับ Timeout/Cancel

## Testing
- [ ] มี Unit Test ครอบคลุม
- [ ] มี Table-Driven Tests
- [ ] ทดสอบ Edge Cases
- [ ] ใช้ Mock สำหรับ Dependencies

## Security
- [ ] ตรวจสอบ SQL Injection (ใช้ Prepared Statements)
- [ ] ตรวจสอบ XSS (ใช้ template escaping)
- [ ] ใช้ HTTPS ใน Production
- [ ] ตั้งค่า CORS อย่างเหมาะสม
```

#### 61.3 New Project Checklist
```markdown
# 🆕 New Project Checklist

## Initial Setup
- [ ] สร้าง Repository (GitHub/GitLab)
- [ ] ตั้งค่า README
- [ ] ตั้งค่า .gitignore
- [ ] ตั้งค่า Go Modules: `go mod init`
- [ ] ตั้งค่า Directory Structure (cmd/, internal/, pkg/)
- [ ] ตั้งค่า Makefile

## Development Setup
- [ ] ติดตั้ง VS Code + Go Extension
- [ ] ติดตั้ง Golangci-lint
- [ ] ติดตั้ง Air (Live Reload)
- [ ] ตั้งค่า Pre-commit Hooks

## Dependencies
- [ ] เลือก HTTP Router (chi, gin, echo)
- [ ] เลือก ORM (GORM, sqlx)
- [ ] เลือก Logger (zap, logrus)
- [ ] เลือก Config Manager (viper)
- [ ] เลือก Validator (go-playground/validator)

## CI/CD
- [ ] ตั้งค่า GitHub Actions / GitLab CI
- [ ] ตั้งค่า Testing Pipeline
- [ ] ตั้งค่า Linting Pipeline
- [ ] ตั้งค่า Build Pipeline
- [ ] ตั้งค่า Deploy Pipeline

## Monitoring
- [ ] ตั้งค่า Health Check
- [ ] ตั้งค่า Metrics (Prometheus)
- [ ] ตั้งค่า Logging
- [ ] ตั้งค่า Alert (Discord/Slack)
```

---

### บทที่ 62: แผนภาพการทำงาน (Workflow Diagram)

แผนภาพนี้แสดงให้เห็นภาพรวมของสถาปัตยกรรมและกระบวนการทำงานของระบบ

#### 62.1 Clean Architecture Diagram
```
┌─────────────────────────────────────────────────┐
│              Presentation Layer                  │
│          (HTTP Handlers / gRPC)                  │
│  ┌─────────────────────────────────────────┐    │
│  │          DTO / Request/Response         │    │
│  └─────────────────────────────────────────┘    │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│             Application Layer                    │
│            (Use Cases / Services)                │
│  ┌─────────────────────────────────────────┐    │
│  │       Business Logic / Workflow         │    │
│  └─────────────────────────────────────────┘    │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│               Domain Layer                       │
│        (Entities / Value Objects)                │
│  ┌─────────────────────────────────────────┐    │
│  │  Interfaces (Repository, Service)       │    │
│  └─────────────────────────────────────────┘    │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│           Infrastructure Layer                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │Database  │  │  Redis   │  │  RabbitMQ│      │
│  └──────────┘  └──────────┘  └──────────┘      │
└─────────────────────────────────────────────────┘
```

#### 62.2 Request Flow Diagram
```
HTTP Request
    ↓
┌────────────────────────┐
│   Middleware Chain     │
│  - Logging             │
│  - Authentication      │
│  - Rate Limiting       │
│  - Recovery            │
└────────────────────────┘
    ↓
┌────────────────────────┐
│   Router (chi/mux)     │
│  - Route Matching      │
│  - URL Parameter       │
└────────────────────────┘
    ↓
┌────────────────────────┐
│   Handler              │
│  - Parse Request       │
│  - Validate Input      │
└────────────────────────┘
    ↓
┌────────────────────────┐
│   Use Case             │
│  - Business Logic      │
│  - Use Repository      │
└────────────────────────┘
    ↓
┌────────────────────────┐
│   Repository           │
│  - Database Query      │
│  - Cache Operation     │
└────────────────────────┘
    ↓
┌────────────────────────┐
│   Response             │
│  - Serialize JSON      │
│  - HTTP Status Code    │
└────────────────────────┘
    ↓
HTTP Response
```

#### 62.3 Deployment Workflow Diagram
```
Code Push → GitHub
     ↓
GitHub Actions (CI)
     ↓
├─ Lint (golangci-lint)
├─ Test (go test)
├─ Build (go build)
└─ Security Scan (govulncheck)
     ↓
Docker Build
     ↓
Push to Docker Registry
     ↓
Deploy to Staging
     ↓
├─ Run Smoke Tests
└─ Manual Verification
     ↓
Deploy to Production
     ↓
├─ Health Check
├─ Monitoring
└─ Rollback Ready
```

---

### บทที่ 63: mop Config – การจัดการ Configuration

`mop` เป็นเครื่องมือที่ช่วยให้การจัดการ Configuration ง่ายและปลอดภัยขึ้น รวมถึงการสร้างไฟล์ config ที่ปลอดภัยด้วยการเข้ารหัส

#### 63.1 การสร้าง Config Structure
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"ssl_mode"`
}

func (d DatabaseConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        d.User, d.Password, d.Host, d.Port, d.DBName)
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    SecretKey string        `mapstructure:"secret_key"`
    Expiry    time.Duration `mapstructure:"expiry"`
}

type LogConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"` // json or console
}
```

#### 63.2 การโหลด Config จากหลายแหล่ง
```go
func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/myapp/")
    viper.AddConfigPath("$HOME/.myapp/")

    // ตั้งค่า Environment Variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // ค่า Default
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.read_timeout", "10s")
    viper.SetDefault("server.write_timeout", "10s")
    viper.SetDefault("log.level", "info")
    viper.SetDefault("log.format", "json")

    // อ่านไฟล์ config
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // ไฟล์ไม่พบ ใช้ ENV อย่างเดียว
        } else {
            return nil, err
        }
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

#### 63.3 การเข้ารหัสและถอดรหัส Config (Security)
```go
package config

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

var encryptionKey = []byte("your-32-byte-key-for-aes-256") // ต้องเก็บใน Environment Variable

func EncryptConfig(data []byte) (string, error) {
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptConfig(encryptedData string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    return plaintext, nil
}

// ใช้ใน Load Function
func LoadEncrypted(encryptedFile string) (*Config, error) {
    encryptedData, err := os.ReadFile(encryptedFile)
    if err != nil {
        return nil, err
    }

    decrypted, err := DecryptConfig(string(encryptedData))
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(decrypted, &config); err != nil {
        return nil, err
    }
    return &config, nil
}
```

---

### ภาคผนวก

#### ภาคผนวก A: GORM CRUD กับฐานข้อมูลหลายประเภท

GORM รองรับฐานข้อมูลหลายประเภท การใช้ CRUD พื้นฐานเหมือนกันเพียงเปลี่ยน Driver

**SQLite:**
```go
import "gorm.io/driver/sqlite"
db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
```

**MySQL:**
```go
import "gorm.io/driver/mysql"
dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

**PostgreSQL:**
```go
import "gorm.io/driver/postgres"
dsn := "host=localhost user=postgres password=pass dbname=test port=9920 sslmode=disable"
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

**CRUD Operation (ใช้ได้ทุกฐานข้อมูล):**
```go
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"index"`
    Age  int
}

// Create
user := User{Name: "Somchai", Age: 30}
db.Create(&user)

// Read
var user User
db.First(&user, 1) // WHERE id = 1
db.Where("name = ?", "Somchai").First(&user)

// Update
db.Model(&user).Update("Age", 31)
db.Model(&user).Updates(User{Name: "New Name", Age: 32})

// Delete
db.Delete(&user, 1) // Soft Delete
db.Unscoped().Delete(&user, 1) // Hard Delete
```

---

#### ภาคผนวก B: คู่มือภาษา Go ฉบับนำไปทำงาน

| หมวดหมู่ | คำสั่ง/เครื่องมือ | คำอธิบาย |
|---------|-----------------|----------|
| **项目管理** | `go mod init` | สร้าง Go Module ใหม่ |
| | `go mod tidy` | จัดการ Dependencies |
| | `go mod vendor` | คัดลอก Dependencies ไป vendor/ |
| **编译** | `go build` | Compile โปรแกรม |
| | `go build -o myapp` | ตั้งชื่อไฟล์ Output |
| | `GOOS=linux GOARCH=amd64 go build` | Cross-compile |
| **测试** | `go test ./...` | รัน Test ทั้งหมด |
| | `go test -v` | แสดงรายละเอียด |
| | `go test -cover` | แสดง Code Coverage |
| | `go test -bench=.` | รัน Benchmark |
| **工具** | `go fmt ./...` | จัดรูปแบบโค้ด |
| | `go vet ./...` | ตรวจสอบโค้ด |
| | `golangci-lint run` | Linter อย่างละเอียด |
| **文档** | `go doc` | ดู Documentation |
| | `go doc fmt.Printf` | ดูเอกสารฟังก์ชัน |
| **性能** | `go test -cpuprofile=cpu.prof` | CPU Profiling |
| | `go test -memprofile=mem.prof` | Memory Profiling |
| | `go tool pprof -http=:8080 cpu.prof` | ดู Profile ด้วย Web UI |

**Best Practices ติดตัว:**
1. ใช้ `gofmt` ทุกครั้งก่อน Commit
2. ใช้ `go vet` เพื่อหา Bug ที่อาจเกิดขึ้น
3. ใช้ `-race` flag เมื่อ Test เพื่อหา Data Race: `go test -race ./...`
4. ใช้ `defer` สำหรับ Cleanup (Close file, Unlock Mutex)
5. ใช้ Context สำหรับ Timeout และ Cancellation
6. หลีกเลี่ยง `panic` ใน Library (ใช้ Error แทน)
7. ใช้ Interface ที่เล็กที่สุดเท่าที่จำเป็น

---

#### ภาคผนวก C: Golang Live-reload (Air)

**Air** เป็นเครื่องมือที่ช่วยให้โค้ดของคุณ reload อัตโนมัติเมื่อมีการเปลี่ยนแปลง ช่วยเพิ่มความเร็วในการพัฒนา

**การติดตั้ง:**
```bash
# ติดตั้ง Air
go install github.com/cosmtrek/air@latest

# หรือใช้ go get
go get -u github.com/cosmtrek/air
```

**การสร้างไฟล์ Configuration (.air.toml):**
```bash
air init
```

**ตัวอย่าง .air.toml:**
```toml
root = "."
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_ansi = false
```

**การใช้งาน:**
```bash
# รัน Air
air

# ใช้ไฟล์ config ที่กำหนดเอง
air -c .air.toml
```

**การใช้ Air ใน Docker:**
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# ติดตั้ง Air
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ใช้ Air สำหรับ Development
CMD ["air", "-c", ".air.toml"]
```

---

### 📌 สรุปเล่มที่ 10
ในเล่มสุดท้ายนี้คุณได้เรียนรู้:
✅ **ตัวอย่างโค้ดครบวงจร** - REST API ที่สมบูรณ์สำหรับระบบจัดการสินค้า
✅ **Task List Template** - เทมเพลตสำหรับจัดการงานและติดตามความคืบหน้า
✅ **Checklist Template** - Pre-Deployment, Code Review, และ New Project Checklist
✅ **Workflow Diagram** - แผนภาพสถาปัตยกรรมและการทำงาน
✅ **mop Config** - การจัดการ Configuration ที่ปลอดภัยและยืดหยุ่น
✅ **ภาคผนวก** - GORM CRUD กับฐานข้อมูลหลายประเภท, คู่มือภาษา Go ฉบับนำไปทำงาน, Golang Live-reload

---

## 🎉 จบบริบูรณ์
 