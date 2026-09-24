## 📚 เล่มที่ 8: การออกแบบสถาปัตยกรรมและ DDD

เล่มนี้จะพาคุณเข้าสู่โลกของการออกแบบซอฟต์แวร์ระดับองค์กร ตั้งแต่การจัดโครงสร้างโปรเจกต์ด้วย Clean Architecture, การวางแผน Workflow, จนถึงหลักการ DDD, Aggregates, Event Storming, CQRS และการออกแบบบริการด้วย Go-DDD ซึ่งเป็นพื้นฐานสำคัญสำหรับการสร้างระบบที่ซับซ้อนและมีอายุการใช้งานยาวนาน

---

### บทที่ 46: Clean Architecture และโครงสร้างโปรเจกต์

Clean Architecture (หรือ Onion Architecture) เป็นแนวคิดที่เสนอโดย Robert C. Martin (Uncle Bob) ที่แยกชั้นของซอฟต์แวร์ออกเป็นวงกลมซ้อนกัน โดยมี **โดเมน (Domain)** อยู่ตรงกลาง และชั้นนอกสุดคือ Infrastructure ซึ่งเป็นการพึ่งพาแบบมุ่งเข้าสู่ศูนย์กลาง (Dependency Inversion)

**หลักการสำคัญ:**
1. **ไม่ขึ้นอยู่กับเฟรมเวิร์ก (Framework Independent)** : โค้ดไม่ผูกติดกับไลบรารีหรือเฟรมเวิร์กใดๆ
2. **ทดสอบได้ง่าย (Testable)** : แต่ละชั้นสามารถทดสอบแยกกันได้
3. **ไม่ขึ้นอยู่กับ UI / Database / External** : ตรรกะธุรกิจไม่รู้ว่าข้อมูลมาจากไหนหรือจะแสดงผลที่ไหน
4. **Dependency Rule** : การพึ่งพาทิศทางต้องจากนอกเข้าในเท่านั้น (ชั้นในไม่รู้จักชั้นนอก)

**โครงสร้างโปรเจกต์ตาม Clean Architecture:**
```
myproject/
├── cmd/
│   └── api/
│       └── main.go                 # จุดเริ่มต้นของแอป
├── internal/
│   ├── domain/                     # ชั้นโดเมน (แกนกลาง)
│   │   ├── user.go                 # Entity, Value Object
│   │   └── repository.go           # Interface ของ Repository
│   ├── usecase/                    # ชั้น Application (Use Cases)
│   │   ├── user/
│   │   │   ├── create.go
│   │   │   ├── get.go
│   │   │   └── interface.go        # Input/Output Ports
│   │   └── ...
│   ├── handler/                    # ชั้น Delivery (HTTP, gRPC)
│   │   ├── user_handler.go
│   │   └── dto/                    # Data Transfer Objects
│   ├── repository/                 # ชั้น Infrastructure (Implementation)
│   │   ├── user_repo.go            # implements domain.Repository
│   │   └── db.go
│   └── config/                     # การตั้งค่า
│       └── config.go
├── pkg/                            # ไลบรารีที่ใช้ร่วมกัน (ไม่ใช่ business logic)
│   ├── logger/
│   └── validator/
├── go.mod
└── go.sum
```

**ตัวอย่างการใช้งานจริง:**

**1. domain/user.go (Entity)**
```go
package domain

import (
    "errors"
    "time"
)

type User struct {
    ID        int64
    Name      string
    Email     string
    Password  string // hashed
    CreatedAt time.Time
    UpdatedAt time.Time
}

// ฟังก์ชันสร้าง User (Factory)
func NewUser(name, email, password string) (*User, error) {
    if name == "" {
        return nil, errors.New("ชื่อผู้ใช้ไม่สามารถว่างได้")
    }
    if !isValidEmail(email) {
        return nil, errors.New("อีเมลไม่ถูกต้อง")
    }
    // hash password ฯลฯ
    return &User{
        Name:      name,
        Email:     email,
        Password:  hashPassword(password),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}

func isValidEmail(email string) bool {
    // ตรรกะตรวจสอบอีเมล
    return true
}

func hashPassword(password string) string {
    // เข้ารหัส
    return "hashed_" + password
}
```

**2. domain/repository.go (Interface)**
```go
package domain

type UserRepository interface {
    Save(user *User) error
    FindByID(id int64) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id int64) error
}
```

**3. usecase/user/interface.go (Ports)**
```go
package user

import "myproject/internal/domain"

// Input Port (คำขอจากภายนอก)
type CreateUserRequest struct {
    Name     string
    Email    string
    Password string
}

type CreateUserResponse struct {
    ID    int64
    Name  string
    Email string
}

// Output Port (Interface สำหรับ Presenter)
type CreateUserPresenter interface {
    Present(response *CreateUserResponse) interface{}
}

// Use Case Interface
type CreateUserUseCase interface {
    Execute(req CreateUserRequest, presenter CreateUserPresenter) (interface{}, error)
}
```

**4. usecase/user/create.go (Implementation)**
```go
package user

import (
    "errors"
    "myproject/internal/domain"
)

type createUserUseCase struct {
    repo domain.UserRepository
}

func NewCreateUserUseCase(repo domain.UserRepository) CreateUserUseCase {
    return &createUserUseCase{repo: repo}
}

func (uc *createUserUseCase) Execute(req CreateUserRequest, presenter CreateUserPresenter) (interface{}, error) {
    // ตรวจสอบอีเมลซ้ำ
    existing, _ := uc.repo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.New("อีเมลนี้ถูกใช้งานแล้ว")
    }

    // สร้าง Entity
    user, err := domain.NewUser(req.Name, req.Email, req.Password)
    if err != nil {
        return nil, err
    }

    // บันทึก
    if err := uc.repo.Save(user); err != nil {
        return nil, err
    }

    // สร้าง Response
    resp := &CreateUserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }

    // ส่งให้ Presenter แปลงเป็นรูปแบบที่ต้องการ (JSON, HTML ฯลฯ)
    return presenter.Present(resp), nil
}
```

**5. repository/user_repo.go (Infrastructure)**
```go
package repository

import (
    "database/sql"
    "myproject/internal/domain"
)

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Save(user *domain.User) error {
    result, err := r.db.Exec(
        "INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
    )
    if err != nil {
        return err
    }
    id, _ := result.LastInsertId()
    user.ID = id
    return nil
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
    row := r.db.QueryRow("SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ?", id)
    var user domain.User
    err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &user, err
}
// ... implement method อื่นๆ
```

**6. handler/user_handler.go (Delivery)**
```go
package handler

import (
    "encoding/json"
    "net/http"
    "myproject/internal/usecase/user"
)

type UserHandler struct {
    createUseCase user.CreateUserUseCase
}

func NewUserHandler(createUC user.CreateUserUseCase) *UserHandler {
    return &UserHandler{createUseCase: createUC}
}

// Presenter ที่แปลง Response เป็น JSON
type jsonPresenter struct{}

func (p *jsonPresenter) Present(resp *user.CreateUserResponse) interface{} {
    return resp // จะถูก Marshal เป็น JSON โดยตรง
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req user.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    presenter := &jsonPresenter{}
    result, err := h.createUseCase.Execute(req, presenter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(result)
}
```

**7. cmd/api/main.go (Entrypoint)**
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    "myproject/internal/handler"
    "myproject/internal/repository"
    "myproject/internal/usecase/user"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "user:pass@/dbname")
    defer db.Close()

    // สร้าง Repo
    userRepo := repository.NewUserRepository(db)

    // สร้าง Use Case
    createUserUC := user.NewCreateUserUseCase(userRepo)

    // สร้าง Handler
    userHandler := handler.NewUserHandler(createUserUC)

    // ตั้งค่า Routes
    http.HandleFunc("/users", userHandler.CreateUser)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

### บทที่ 47: Blueprint สำหรับโปรเจกต์ Go ระดับ Production

นอกจากการแบ่งชั้นตาม Clean Architecture แล้ว ยังมีแนวทางปฏิบัติอื่นๆ ที่จะทำให้โปรเจกต์ของคุณพร้อมสำหรับ Production:

**1. การจัดระเบียบไฟล์และแพคเกจ:**
- `cmd/` : มี subdirectory สำหรับแต่ละ executable (เช่น `api`, `worker`, `cli`)
- `internal/` : โค้ดส่วนตัวที่ไม่ต้องการให้แพคเกจอื่น import (ใช้ได้เฉพาะในโปรเจกต์นี้)
- `pkg/` : โค้ดที่สามารถนำไปใช้กับโปรเจกต์อื่นๆ ได้ (public library)
- `api/` : ไฟล์ proto/gRPC หรือ OpenAPI spec
- `scripts/` : สคริปต์สำหรับ build, migration, deployment
- `test/` : การทดสอบระดับ Integration / E2E

**2. การจัดการ Config ใน Production:**
ใช้ Environment Variables เป็นหลัก ร่วมกับไฟล์ Config เฉพาะ Environment
```go
package config

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

**3. Graceful Shutdown:**
```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: router}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s", err)
        }
    }()

    // รอ Signal (Ctrl+C)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    log.Println("Server exiting")
}
```

**4. Health Check Endpoint:**
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
})

http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
    // ตรวจสอบ DB, Redis, ฯลฯ
    if db.Ping() == nil && redis.Ping() == nil {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
})
```

**5. การใช้ Middleware สำหรับ Logging, Recovery, CORS, Rate Limiting:**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

---

### บทที่ 48: การออกแบบ Workflow และ Task Management

ในระบบที่ซับซ้อน การจัดการ workflow (ลำดับขั้นตอนการทำงาน) เป็นสิ่งสำคัญ สามารถออกแบบได้หลายรูปแบบ:

**1. Saga Pattern (สำหรับ Distributed Transactions):**
```go
type SagaStep interface {
    Execute(ctx context.Context, data map[string]interface{}) error
    Compensate(ctx context.Context, data map[string]interface{}) error
}

type Saga struct {
    steps []SagaStep
}

func (s *Saga) Run(ctx context.Context, data map[string]interface{}) error {
    executed := []int{}
    for i, step := range s.steps {
        if err := step.Execute(ctx, data); err != nil {
            // ถ้าผิดพลาด ให้ย้อนกลับ (Compensate)
            for j := len(executed) - 1; j >= 0; j-- {
                s.steps[executed[j]].Compensate(ctx, data)
            }
            return err
        }
        executed = append(executed, i)
    }
    return nil
}
```

**2. Workflow Engine (เช่น Temporal, Cadence):**
สามารถใช้ไลบรารีเพื่อจัดการ workflow ที่ซับซ้อน แต่ในบทนี้เราจะแสดงการออกแบบง่ายๆ ด้วย Channel และ Goroutine:

```go
type Task struct {
    ID   string
    Name string
    Data map[string]interface{}
}

type Worker struct {
    tasks chan Task
    wg    sync.WaitGroup
}

func NewWorker(numWorkers int) *Worker {
    w := &Worker{tasks: make(chan Task, 100)}
    for i := 0; i < numWorkers; i++ {
        w.wg.Add(1)
        go w.worker()
    }
    return w
}

func (w *Worker) worker() {
    defer w.wg.Done()
    for task := range w.tasks {
        fmt.Printf("Processing task %s: %s\n", task.ID, task.Name)
        // ประมวลผล task
        time.Sleep(1 * time.Second)
    }
}

func (w *Worker) AddTask(task Task) {
    w.tasks <- task
}

func (w *Worker) Shutdown() {
    close(w.tasks)
    w.wg.Wait()
}
```

---

### บทที่ 49: หลักการ DDD และการนำไปใช้ใน Go

Domain-Driven Design (DDD) เป็นแนวทางในการออกแบบซอฟต์แวร์ที่เน้นที่ **โดเมนธุรกิจ (Business Domain)** และการทำงานร่วมกันระหว่างนักพัฒนาและผู้เชี่ยวชาญทางธุรกิจ

**แนวคิดหลักของ DDD:**

1. **Ubiquitous Language (ภาษาร่วม)** : ใช้คำศัพท์เดียวกันระหว่างทีมธุรกิจและทีมพัฒนา เช่น "ลูกค้า" "คำสั่งซื้อ" "การชำระเงิน"
2. **Bounded Context (บริบทที่มีขอบเขต)** : แบ่งระบบออกเป็นบริบทย่อยๆ แต่ละบริบทมีโมเดลของตัวเอง
3. **Entity** : วัตถุที่มีเอกลักษณ์ (Identity) และเปลี่ยนแปลงได้ตลอดอายุ
4. **Value Object** : วัตถุที่ไม่มีเอกลักษณ์ ระบุด้วยค่า (เช่น ที่อยู่, เงิน)
5. **Aggregate** : กลุ่มของ Entity และ Value Object ที่ถือเป็นหน่วยเดียวกัน มี Aggregate Root เป็นตัวควบคุม
6. **Repository** : ให้วิธีการเข้าถึง Aggregate
7. **Domain Service** : บริการที่มีตรรกะที่ไม่เกี่ยวข้องกับ Entity เดียว
8. **Domain Event** : เหตุการณ์ที่เกิดขึ้นในโดเมน

**ตัวอย่างการนำ DDD ไปใช้ใน Go:**

**1. Entity และ Value Object:**
```go
// Value Object: Money
type Money struct {
    Amount   float64
    Currency string
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("สกุลเงินไม่ตรงกัน")
    }
    return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Entity: Order
type Order struct {
    ID         string
    CustomerID string
    Items      []OrderItem
    Total      Money
    Status     OrderStatus
    CreatedAt  time.Time
}

// Entity มี Identity (ID) และสามารถเปลี่ยนแปลงได้
func (o *Order) AddItem(item OrderItem) {
    o.Items = append(o.Items, item)
    o.Total = o.Total.Add(item.Price) // สมมติว่าคำนวณใหม่
}

// Aggregate Root: Order (ควบคุมการเข้าถึง OrderItem)
type OrderItem struct {
    ProductID string
    Quantity  int
    Price     Money
}
```

**2. Repository Interface ตาม DDD:**
```go
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id string) (*Order, error)
    FindByCustomer(customerID string) ([]*Order, error)
}
```

**3. Domain Service:**
```go
type PricingService interface {
    CalculateTotal(items []OrderItem) Money
}

type DefaultPricingService struct {
    // dependencies เช่น DiscountService
}

func (s *DefaultPricingService) CalculateTotal(items []OrderItem) Money {
    total := Money{Amount: 0, Currency: "THB"}
    for _, item := range items {
        total, _ = total.Add(item.Price)
    }
    // apply discounts
    return total
}
```

**4. Domain Event:**
```go
type DomainEvent interface {
    GetAggregateID() string
    GetOccurredAt() time.Time
}

type OrderCreatedEvent struct {
    AggregateID string
    CustomerID  string
    Total       Money
    OccurredAt  time.Time
}

func (e OrderCreatedEvent) GetAggregateID() string { return e.AggregateID }
func (e OrderCreatedEvent) GetOccurredAt() time.Time { return e.OccurredAt }

// Event Publisher
type EventPublisher interface {
    Publish(event DomainEvent) error
}
```

---

### บทที่ 50: Aggregates, Event Storming และ CQRS

**1. Aggregates (กลุ่มวัตถุ):**
Aggregate คือกลุ่มของ Entity และ Value Object ที่ถือเป็นหน่วยเดียวกัน โดยมี **Aggregate Root** เป็นตัวแทนในการเข้าถึงจากภายนอก

**ตัวอย่าง:**
```go
// Aggregate Root: Order
type Order struct {
    ID         string
    CustomerID string
    Items      []OrderItem
    Status     OrderStatus
    Events     []DomainEvent // เก็บเหตุการณ์ที่เกิดขึ้น
}

// ฟังก์ชันใน Aggregate Root ที่เปลี่ยนแปลงสถานะ
func (o *Order) Cancel() error {
    if o.Status == Shipped {
        return errors.New("ไม่สามารถยกเลิกคำสั่งซื้อที่จัดส่งแล้ว")
    }
    o.Status = Cancelled
    o.Events = append(o.Events, OrderCancelledEvent{
        AggregateID: o.ID,
        OccurredAt:  time.Now(),
    })
    return nil
}

// OrderItem ไม่สามารถเข้าถึงโดยตรงจากภายนอกได้ ต้องผ่าน Order เท่านั้น
```

**2. Event Storming (การระดมสมองด้วยเหตุการณ์):**
Event Storming คือ workshop ที่ทีมธุรกิจและพัฒนาช่วยกันระบุ **Domain Events** ทั้งหมดที่เกิดขึ้นในระบบ แล้วจัดลำดับเหตุการณ์เหล่านั้นเพื่อให้เข้าใจกระบวนการทางธุรกิจ

**ตัวอย่าง Domain Events ที่พบบ่อย:**
- `OrderPlaced`
- `PaymentReceived`
- `OrderShipped`
- `OrderDelivered`
- `OrderCancelled`

**การนำ Event Storming มาใช้ในโค้ด:**
```go
// กำหนด Event types
const (
    EventOrderPlaced   = "order.placed"
    EventPaymentReceived = "payment.received"
    EventOrderShipped  = "order.shipped"
)

type Event struct {
    Type      string
    Aggregate string
    Data      map[string]interface{}
    Timestamp time.Time
}

// Event Handler
type EventHandler interface {
    Handle(event Event) error
}

// Event Dispatcher
type EventDispatcher struct {
    handlers map[string][]EventHandler
}

func (d *EventDispatcher) Register(eventType string, handler EventHandler) {
    d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *EventDispatcher) Dispatch(event Event) error {
    for _, handler := range d.handlers[event.Type] {
        if err := handler.Handle(event); err != nil {
            return err
        }
    }
    return nil
}
```

**3. CQRS (Command Query Responsibility Segregation):**
CQRS คือการแยก **Command** (การเปลี่ยนแปลงข้อมูล) ออกจาก **Query** (การอ่านข้อมูล) โดยใช้โมเดลและฐานข้อมูลที่ต่างกัน

- **Command Side (Write)**: ใช้โมเดลที่เน้นตรรกะธุรกิจ ใช้ Event Sourcing
- **Query Side (Read)**: ใช้โมเดลที่ปรับให้เหมาะกับการอ่านและแสดงผล (Denormalized)

**ตัวอย่าง CQRS ใน Go:**
```go
// Command
type CreateOrderCommand struct {
    CustomerID string
    Items      []OrderItem
}

type CommandHandler interface {
    Handle(cmd interface{}) error
}

// Query
type GetOrderQuery struct {
    OrderID string
}

type OrderReadModel struct {
    ID         string
    Customer   string
    Items      []OrderItemDTO
    Total      float64
    Status     string
}

type QueryHandler interface {
    Handle(query interface{}) (interface{}, error)
}

// ใช้ Repository ที่แยกจากกัน
type OrderCommandRepository interface {
    Save(order *Order) error
}

type OrderQueryRepository interface {
    FindByID(id string) (*OrderReadModel, error)
}
```

**ประโยชน์ของ CQRS:**
- ปรับประสิทธิภาพการอ่านและเขียนแยกกัน
- ลดความซับซ้อนของโมเดล
- รองรับการขยายระบบในอนาคต

---

### บทที่ 51: การออกแบบบริการด้วย Go-DDD

Go-DDD เป็นแนวทางการนำ DDD มาประยุกต์ใช้ใน Go โดยใช้แพคเกจเสริมและเครื่องมือต่างๆ

**1. ใช้ Event Sourcing กับ Event Store:**
```go
// Event Store Interface
type EventStore interface {
    Save(event Event) error
    Load(aggregateID string) ([]Event, error)
}

// Aggregate ที่ใช้ Event Sourcing
type OrderAggregate struct {
    ID      string
    Version int
    Events  []Event
    // สถานะ
    Status string
}

func (a *OrderAggregate) ApplyEvent(event Event) {
    // เปลี่ยนสถานะตาม event
    switch e := event.(type) {
    case OrderPlacedEvent:
        a.Status = "placed"
    case OrderShippedEvent:
        a.Status = "shipped"
    }
}

func (a *OrderAggregate) Rebuild(events []Event) {
    for _, e := range events {
        a.ApplyEvent(e)
    }
}
```

**2. ใช้ Dependency Injection Container (เช่น wire, fx):**
```go
// wire.go (ใช้ Google Wire)
// +build wireinject

func InitializeApp() (*App, error) {
    wire.Build(
        config.Load,
        repository.NewUserRepository,
        usecase.NewCreateUserUseCase,
        handler.NewUserHandler,
        NewApp,
    )
    return nil, nil
}
```

**3. การจัดการ Unit of Work และ Transaction:**
```go
type UnitOfWork interface {
    Begin(ctx context.Context) error
    Commit() error
    Rollback() error
    UserRepository() domain.UserRepository
    OrderRepository() domain.OrderRepository
}

func (uc *createOrderUseCase) Execute(req CreateOrderRequest) error {
    uow := NewUnitOfWork()
    uow.Begin(ctx)
    defer uow.Rollback() // ถ้ายังไม่ commit จะ rollback

    order := domain.NewOrder(req.CustomerID)
    if err := uow.OrderRepository().Save(order); err != nil {
        return err
    }

    // ทำงานอื่นๆ
    return uow.Commit()
}
```

**4. ใช้ Middleware สำหรับ Authentication / Authorization:**
```go
type AuthMiddleware struct {
    jwtSecret string
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        // ตรวจสอบ token, ดึง user ID
        userID, err := m.validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        // ส่ง userID ผ่าน Context
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**5. การทำ Unit Test สำหรับ DDD:**
ใช้ Mock Repository เพื่อทดสอบ Use Case:
```go
type mockUserRepo struct {
    users map[int64]*domain.User
}

func (m *mockUserRepo) Save(user *domain.User) error {
    user.ID = int64(len(m.users) + 1)
    m.users[user.ID] = user
    return nil
}

func TestCreateUserUseCase(t *testing.T) {
    repo := &mockUserRepo{users: make(map[int64]*domain.User)}
    uc := NewCreateUserUseCase(repo)
    req := CreateUserRequest{Name: "สมชาย", Email: "test@test.com", Password: "1234"}
    presenter := &mockPresenter{}
    _, err := uc.Execute(req, presenter)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if len(repo.users) != 1 {
        t.Errorf("expected 1 user, got %d", len(repo.users))
    }
}
```

---

### 📌 สรุปเล่มที่ 8
ในเล่มนี้คุณได้เรียนรู้การออกแบบซอฟต์แวร์ระดับองค์กร:
✅ **Clean Architecture** และการจัดโครงสร้างโปรเจกต์ที่แยกชั้นชัดเจน
✅ **Blueprint สำหรับ Production** (Graceful Shutdown, Health Check, Middleware)
✅ การออกแบบ **Workflow** และ Task Management
✅ หลักการ **Domain-Driven Design (DDD)** และการนำไปใช้ใน Go
✅ **Aggregates, Event Storming** และ **CQRS** เพื่อจัดการความซับซ้อน
✅ การออกแบบบริการด้วย Go-DDD พร้อมตัวอย่างและแนวทางปฏิบัติ

--- 