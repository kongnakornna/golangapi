# 📚 ชุดคู่มือ Golang Module Master Series

> **สำหรับสร้าง แก้ไข และขาย Module ด้วยภาษา Go**
> ยึดตาม Clean Architecture + DDD + Modular Monolith
> อ้างอิงจากโปรเจกต์ `icmongolang`

---

##  เล่ม 1: สร้าง Module ใหม่ (Creating New Modules)
##  เล่ม 2: แก้ไข Module เดิม (Modifying Existing Modules)
##  เล่ม 3: ขาย Module / ต่อยอดเชิงพาณิชย์ (Selling/Commercializing Modules)
##  เล่ม 4: ทดสอบและ Deployment
##  เล่ม 5: บำรุงรักษาและขยายระบบ


# 📕 เล่ม 1: คู่มือสร้าง Module ใหม่ (Module Creation Guide)

## สารบัญ
1. [ปรัชญาและหลักการ](#1-ปรัชญาและหลักการ)
2. [การเตรียมความพร้อม](#2-การเตรียมความพร้อม)
3. [ขั้นตอนการสร้าง Module 12 ขั้น](#3-ขั้นตอนการสร้าง-module-12-ขั้น)
4. [Code Template ทุก Layer](#4-code-template-ทุก-layer)
5. [Migration & Schema](#5-migration--schema)
6. [Wire-up & Bootstrap](#6-wire-up--bootstrap)
7. [Checklist ก่อน Commit](#7-checklist-ก่อน-commit)

---

## 1. ปรัชญาและหลักการ

### 1.1 Dependency Rule (กฎเหล็ก)

```
┌─────────────────────────────────────────┐
│  Interfaces (HTTP/WS/gRPC)              │  ← รู้จักทุก Layer
├─────────────────────────────────────────┤
│  Infrastructure (Postgres/Redis/Kafka)  │  ← รู้จัก Domain, Application
├─────────────────────────────────────────┤
│  Application (Use Cases)                │  ← รู้จัก Domain เท่านั้น
├─────────────────────────────────────────┤
│  Domain (Entity, VO, Repository IF)     │  ← ไม่รู้จักใครเลย ❤️
└─────────────────────────────────────────┘
```

| Layer | import ได้ | ห้าม import |
|---|---|---|
| Domain | `stdlib`, `uuid`, `decimal` | gorm, gin, chi, sarama, redis |
| Application | Domain | Infrastructure, Interfaces |
| Infrastructure | Domain, Application | Interfaces |
| Interfaces | ทุก Layer | – |

### 1.2 เมื่อไหร่ควรใช้ Clean + DDD?

✅ **ควรใช้:**
- Business logic ซับซ้อน มี state transitions หลายแบบ
- ต้องสลับ technology (Postgres ↔ MongoDB)
- ทีม > 2 คน, โมดูล > 5
- ต้อง test โดยไม่ต้อง spin infra

❌ **ไม่ควรใช้:**
- CRUD ง่ายๆ ไร้ business rules
- Prototype / PoC < 1 สัปดาห์
- Script ขนาดเล็ก

---

## 2. การเตรียมความพร้อม

### 2.1 ตรวจสอบ Environment

```bash
# ตรวจสอบ Go version (ต้อง 1.22+)
go version

# ตรวจสอบ tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# ตรวจสอบโครงสร้างเดิม
ls -la internal/modules/          # ดูว่า module อะไรมีแล้ว
ls -la pkg/                       # ดูว่า shared package อะไรมีแล้ว
cat docs/Template_Structure.md    # อ่านโครงสร้างจริง
```

### 2.2 ตั้งชื่อ Module (Naming Convention)

| ประเภท | รูปแบบ | ตัวอย่าง |
|---|---|---|
| Module folder | lowercase, snake | `payment`, `purchase_order` |
| Entity | PascalCase | `Payment`, `PurchaseOrder` |
| Table prefix | module_name + `_` | `payment_transactions` |
| Kafka topic | `<module>.<entity>.<action>` | `payment.txn.created` |
| Redis key | `<module>:<entity>:<id>` | `payment:txn:abc-123` |
| Migration | `YYYYMMDD_<module>_<desc>.sql` | `20260401_payment_init.sql` |

### 2.3 เตรียม Folder Structure

```bash
export MODULE=payment
mkdir -p internal/modules/$MODULE/{domain/{entity,value_object,repository,service,event,errors},application,infrastructure/{persistence/{postgres,redis},messaging},interfaces/http}
```

---

## 3. ขั้นตอนการสร้าง Module 12 ขั้น

### ขั้นตอนที่ 1: ออกแบบ Domain Model

**คำถามที่ต้องตอบ:**
1. Aggregate Root คืออะไร?
2. มี Value Object อะไรบ้าง?
3. State transitions มีอะไรบ้าง?
4. Invariants (กฎที่ต้องเป็นจริงเสมอ) คืออะไร?
5. Domain Events ที่เกิดขึ้นคืออะไร?

**ตัวอย่าง: Payment Module**

```
Aggregate Root: Payment
├── id: UUID
├── amount: Money (VO)
├── status: PaymentStatus (VO)
├── method: PaymentMethod (VO)
└── transactions: []Transaction (Entity ย่อย)

State Transitions:
PENDING → PROCESSING → SUCCESS → REFUNDED
                    ↘ FAILED

Domain Events:
- PaymentCreated
- PaymentSucceeded
- PaymentFailed
- PaymentRefunded
```

### ขั้นตอนที่ 2: สร้าง Value Objects

```go
// internal/modules/payment/domain/value_object/status.go
package valueobject

type PaymentStatus string

const (
    PaymentStatusPending    PaymentStatus = "PENDING"
    PaymentStatusProcessing PaymentStatus = "PROCESSING"
    PaymentStatusSuccess    PaymentStatus = "SUCCESS"
    PaymentStatusFailed     PaymentStatus = "FAILED"
    PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

// IsValid ตรวจสอบว่าสถานะเป็นค่าที่ถูกต้อง
// IsValid checks if the status is valid
func (s PaymentStatus) IsValid() bool {
    switch s {
    case PaymentStatusPending, PaymentStatusProcessing,
         PaymentStatusSuccess, PaymentStatusFailed, PaymentStatusRefunded:
        return true
    }
    return false
}

// CanTransitionTo ตรวจสอบว่าสามารถเปลี่ยนสถานะได้หรือไม่
// CanTransitionTo checks state transition rules
func (s PaymentStatus) CanTransitionTo(next PaymentStatus) bool {
    transitions := map[PaymentStatus][]PaymentStatus{
        PaymentStatusPending:    {PaymentStatusProcessing, PaymentStatusFailed},
        PaymentStatusProcessing: {PaymentStatusSuccess, PaymentStatusFailed},
        PaymentStatusSuccess:    {PaymentStatusRefunded},
        PaymentStatusFailed:     {}, // terminal
        PaymentStatusRefunded:   {}, // terminal
    }
    for _, allowed := range transitions[s] {
        if allowed == next {
            return true
        }
    }
    return false
}
```

### ขั้นตอนที่ 3: สร้าง Entity (Aggregate Root)

```go
// internal/modules/payment/domain/entity/payment.go
package entity

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

// Payment is the aggregate root for payment transactions
// Payment เป็น aggregate root สำหรับรายการชำระเงิน
type Payment struct {
    ID          uuid.UUID                     `json:"id"`
    UserID      uuid.UUID                     `json:"user_id"`
    OrderID     uuid.UUID                     `json:"order_id"`
    Amount      decimal.Decimal               `json:"amount"`
    Currency    string                        `json:"currency"`
    Method      valueobject.PaymentMethod     `json:"method"`
    Status      valueobject.PaymentStatus     `json:"status"`
    RefundedAt  *time.Time                    `json:"refunded_at,omitempty"`
    CreatedAt   time.Time                     `json:"created_at"`
    UpdatedAt   time.Time                     `json:"updated_at"`
}

// NewPayment creates a new payment with validation
// NewPayment สร้าง payment ใหม่พร้อม validation
func NewPayment(userID, orderID uuid.UUID, amount decimal.Decimal, currency string, method valueobject.PaymentMethod) (*Payment, error) {
    // Validate user
    // ตรวจสอบ user ID
    if userID == uuid.Nil {
        return nil, domainerrors.ErrInvalidUserID
    }
    // Validate order
    // ตรวจสอบ order ID
    if orderID == uuid.Nil {
        return nil, domainerrors.ErrInvalidOrderID
    }
    // Amount must be positive
    // ยอดเงินต้องมากกว่าศูนย์
    if amount.LessThanOrEqual(decimal.Zero) {
        return nil, domainerrors.ErrInvalidAmount
    }
    // Method must be valid
    // วิธีชำระต้องถูกต้อง
    if !method.IsValid() {
        return nil, domainerrors.ErrInvalidPaymentMethod
    }

    now := time.Now()
    return &Payment{
        ID:        uuid.New(),
        UserID:    userID,
        OrderID:   orderID,
        Amount:    amount,
        Currency:  currency,
        Method:    method,
        Status:    valueobject.PaymentStatusPending,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// StartProcessing changes status to PROCESSING
// StartProcessing เปลี่ยนสถานะเป็น PROCESSING
func (p *Payment) StartProcessing() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusProcessing) {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusProcessing
    p.UpdatedAt = time.Now()
    return nil
}

// MarkSuccess marks the payment as successful
// MarkSuccess ทำเครื่องหมายว่าชำระสำเร็จ
func (p *Payment) MarkSuccess() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusSuccess) {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusSuccess
    p.UpdatedAt = time.Now()
    return nil
}

// Refund refunds a successful payment
// Refund คืนเงินสำหรับ payment ที่สำเร็จ
func (p *Payment) Refund() error {
    if !p.Status.CanTransitionTo(valueobject.PaymentStatusRefunded) {
        return domainerrors.ErrCannotRefund
    }
    now := time.Now()
    p.Status = valueobject.PaymentStatusRefunded
    p.RefundedAt = &now
    p.UpdatedAt = now
    return nil
}

// IsRefundable checks if refund is possible
// IsRefundable ตรวจสอบว่าสามารถคืนเงินได้หรือไม่
func (p *Payment) IsRefundable() bool {
    return p.Status == valueobject.PaymentStatusSuccess && p.RefundedAt == nil
}
```

### ขั้นตอนที่ 4: สร้าง Domain Errors

```go
// internal/modules/payment/domain/errors/errors.go
package domainerrors

import "errors"

var (
    ErrPaymentNotFound         = errors.New("payment not found")
    ErrInvalidUserID           = errors.New("invalid user id")
    ErrInvalidOrderID          = errors.New("invalid order id")
    ErrInvalidAmount           = errors.New("amount must be positive")
    ErrInvalidPaymentMethod    = errors.New("invalid payment method")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrCannotRefund            = errors.New("payment cannot be refunded")
    ErrPaymentAlreadyExists    = errors.New("payment already exists")
    ErrUnauthorized            = errors.New("unauthorized")
)
```

### ขั้นตอนที่ 5: สร้าง Repository Interface

```go
// internal/modules/payment/domain/repository/payment_repository.go
package repository

import (
    "context"
    "github.com/google/uuid"
    "icmongolang/internal/modules/payment/domain/entity"
)

// PaymentRepository defines the contract for payment persistence
// PaymentRepository กำหนดสัญญาสำหรับการจัดเก็บ payment
type PaymentRepository interface {
    Save(ctx context.Context, p *entity.Payment) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
    FindByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Payment, error)
    Update(ctx context.Context, p *entity.Payment) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### ขั้นตอนที่ 6: สร้าง Application Use Case

```go
// internal/modules/payment/application/create_payment.go
package application

import (
    "context"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    "icmongolang/internal/modules/payment/domain/repository"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
    "icmongolang/pkg/logger"
)

// CreatePaymentUseCase handles creating a new payment
// CreatePaymentUseCase จัดการการสร้าง payment ใหม่
type CreatePaymentUseCase struct {
    repo   repository.PaymentRepository
    logger logger.Logger
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository, log logger.Logger) *CreatePaymentUseCase {
    return &CreatePaymentUseCase{repo: repo, logger: log}
}

// CreatePaymentInput is the input for creating a payment
// CreatePaymentInput เป็น input สำหรับการสร้าง payment
type CreatePaymentInput struct {
    UserID   uuid.UUID
    OrderID  uuid.UUID
    Amount   decimal.Decimal
    Currency string
    Method   string
}

// CreatePaymentOutput is the output of payment creation
// CreatePaymentOutput เป็น output ของการสร้าง payment
type CreatePaymentOutput struct {
    ID       uuid.UUID `json:"id"`
    Status   string    `json:"status"`
    Amount   string    `json:"amount"`
    Currency string    `json:"currency"`
}

// Execute runs the use case
// Execute รัน use case
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // Step 1: Check duplicate
    // ขั้นที่ 1: ตรวจสอบซ้ำ
    existing, _ := uc.repo.FindByOrderID(ctx, input.OrderID)
    if existing != nil {
        return nil, domainerrors.ErrPaymentAlreadyExists
    }

    // Step 2: Create entity
    // ขั้นที่ 2: สร้าง entity
    method := valueobject.PaymentMethod(input.Method)
    payment, err := entity.NewPayment(input.UserID, input.OrderID, input.Amount, input.Currency, method)
    if err != nil {
        return nil, err
    }

    // Step 3: Persist
    // ขั้นที่ 3: บันทึก
    if err := uc.repo.Save(ctx, payment); err != nil {
        uc.logger.Error("save payment failed", "error", err, "order_id", input.OrderID)
        return nil, err
    }

    uc.logger.Info("payment created", "id", payment.ID, "order_id", input.OrderID)

    // Step 4: Return output
    // ขั้นที่ 4: คืนค่า output
    return &CreatePaymentOutput{
        ID:       payment.ID,
        Status:   string(payment.Status),
        Amount:   payment.Amount.String(),
        Currency: payment.Currency,
    }, nil
}
```

### ขั้นตอนที่ 7: สร้าง Infrastructure - GORM Model

```go
// internal/modules/payment/infrastructure/persistence/postgres/models.go
package postgres

import (
    "time"
    "github.com/google/uuid"
)

// PaymentModel is the GORM model for payments table
// PaymentModel เป็น GORM model สำหรับตาราง payments
type PaymentModel struct {
    ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID     uuid.UUID  `gorm:"type:uuid;index;not null"`
    OrderID    uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
    Amount     string     `gorm:"type:numeric(15,2);not null"`
    Currency   string     `gorm:"type:varchar(3);not null;default:'THB'"`
    Method     string     `gorm:"type:varchar(30);not null"`
    Status     string     `gorm:"type:varchar(20);not null;index"`
    RefundedAt *time.Time
    CreatedAt  time.Time  `gorm:"default:now()"`
    UpdatedAt  time.Time  `gorm:"default:now()"`
}

// TableName returns the table name with module prefix
// TableName คืนชื่อตารางพร้อม prefix ของ module
func (PaymentModel) TableName() string {
    return "payment_transactions"
}
```

### ขั้นตอนที่ 8: สร้าง Repository Implementation

```go
// internal/modules/payment/infrastructure/persistence/postgres/payment_repo_impl.go
package postgres

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
    "icmongolang/internal/modules/payment/domain/entity"
    domainerrors "icmongolang/internal/modules/payment/domain/errors"
    valueobject "icmongolang/internal/modules/payment/domain/value_object"
)

type paymentRepoImpl struct {
    db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepoImpl {
    return &paymentRepoImpl{db: db}
}

func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    m := r.toModel(p)
    return r.db.WithContext(ctx).Create(m).Error
}

func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrPaymentNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

func (r *paymentRepoImpl) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.db.WithContext(ctx).First(&m, "order_id = ?", orderID).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrPaymentNotFound
    }
    if err != nil {
        return nil, err
    }
    return r.toEntity(&m), nil
}

func (r *paymentRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Payment, error) {
    var models []PaymentModel
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Limit(limit).Offset(offset).
        Find(&models).Error
    if err != nil {
        return nil, err
    }
    payments := make([]*entity.Payment, len(models))
    for i := range models {
        payments[i] = r.toEntity(&models[i])
    }
    return payments, nil
}

func (r *paymentRepoImpl) Update(ctx context.Context, p *entity.Payment) error {
    m := r.toModel(p)
    return r.db.WithContext(ctx).Save(m).Error
}

func (r *paymentRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&PaymentModel{}, "id = ?", id).Error
}

// toModel converts entity to GORM model
// toModel แปลง entity เป็น GORM model
func (r *paymentRepoImpl) toModel(e *entity.Payment) *PaymentModel {
    return &PaymentModel{
        ID:         e.ID,
        UserID:     e.UserID,
        OrderID:    e.OrderID,
        Amount:     e.Amount.String(),
        Currency:   e.Currency,
        Method:     string(e.Method),
        Status:     string(e.Status),
        RefundedAt: e.RefundedAt,
        CreatedAt:  e.CreatedAt,
        UpdatedAt:  e.UpdatedAt,
    }
}

// toEntity converts GORM model to domain entity
// toEntity แปลง GORM model เป็น domain entity
func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    amount, _ := decimal.NewFromString(m.Amount)
    return &entity.Payment{
        ID:         m.ID,
        UserID:     m.UserID,
        OrderID:    m.OrderID,
        Amount:     amount,
        Currency:   m.Currency,
        Method:     valueobject.PaymentMethod(m.Method),
        Status:     valueobject.PaymentStatus(m.Status),
        RefundedAt: m.RefundedAt,
        CreatedAt:  m.CreatedAt,
        UpdatedAt:  m.UpdatedAt,
    }
}
```

### ขั้นตอนที่ 9: สร้าง HTTP Handler

```go
// internal/modules/payment/interfaces/http/payment_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "icmongolang/internal/modules/payment/application"
    "icmongolang/pkg/httputil"
    "icmongolang/pkg/validator"
)

type PaymentHandler struct {
    createUC  *application.CreatePaymentUseCase
    validator *validator.Validator
}

func NewPaymentHandler(createUC *application.CreatePaymentUseCase, v *validator.Validator) *PaymentHandler {
    return &PaymentHandler{createUC: createUC, validator: v}
}

// CreatePaymentRequest is the HTTP request for creating a payment
// CreatePaymentRequest เป็น HTTP request สำหรับสร้าง payment
type CreatePaymentRequest struct {
    OrderID  string `json:"order_id" validate:"required,uuid"`
    Amount   string `json:"amount" validate:"required"`
    Currency string `json:"currency" validate:"required,len=3"`
    Method   string `json:"method" validate:"required,oneof=CARD QR_PROMPT_PAY BANK_TRANSFER"`
}

// Create handles POST /api/v1/payments
// Create จัดการ POST /api/v1/payments
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Step 1: Extract user from context
    // ขั้นที่ 1: ดึง user จาก context
    userID, ok := r.Context().Value("user_id").(uuid.UUID)
    if !ok {
        httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }

    // Step 2: Parse and validate request
    // ขั้นที่ 2: parse และ validate request
    var req CreatePaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }
    if err := h.validator.Validate(req); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    // Step 3: Parse IDs and amount
    // ขั้นที่ 3: parse IDs และ amount
    orderID, err := uuid.Parse(req.OrderID)
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid order_id"})
        return
    }
    amount, err := decimal.NewFromString(req.Amount)
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid amount"})
        return
    }

    // Step 4: Execute use case
    // ขั้นที่ 4: รัน use case
    output, err := h.createUC.Execute(r.Context(), application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   amount,
        Currency: req.Currency,
        Method:   req.Method,
    })
    if err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    // Step 5: Return response
    // ขั้นที่ 5: คืนค่า response
    httputil.JSON(w, http.StatusCreated, output)
}

// Get handles GET /api/v1/payments/{id}
// Get จัดการ GET /api/v1/payments/{id}
func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    if _, err := uuid.Parse(idStr); err != nil {
        httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
        return
    }
    // TODO: implement get logic
    httputil.JSON(w, http.StatusOK, map[string]string{"id": idStr})
}
```

### ขั้นตอนที่ 10: สร้าง Routes

```go
// internal/modules/payment/interfaces/http/routes.go
package http

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

// Handlers aggregates all handlers of the payment module
// Handlers รวม handler ทั้งหมดของโมดูล payment
type Handlers struct {
    Payment *PaymentHandler
}

// RegisterRoutes registers all routes for the payment module
// RegisterRoutes ลงทะเบียน route ทั้งหมดสำหรับโมดูล payment
func RegisterRoutes(r chi.Router, h *Handlers, authMiddleware func(http.Handler) http.Handler) {
    r.Route("/api/v1/payments", func(r chi.Router) {
        r.Use(authMiddleware)
        r.Post("/", h.Payment.Create)
        r.Get("/{id}", h.Payment.Get)
    })
}
```

### ขั้นตอนที่ 11: สร้าง Migration

```sql
-- migrations/20260401_payment_init.sql
-- ============================================================
-- Payment Module — Initial Schema
-- Prefix: payment_
-- ============================================================

CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    order_id UUID NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    method VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    refunded_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_payment_transactions_order_id
    ON payment_transactions(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_user_id
    ON payment_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_status
    ON payment_transactions(status);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION payment_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS payment_transactions_updated_at ON payment_transactions;
CREATE TRIGGER payment_transactions_updated_at
    BEFORE UPDATE ON payment_transactions
    FOR EACH ROW EXECUTE FUNCTION payment_update_updated_at();
```

### ขั้นตอนที่ 12: Wire-up ใน main.go

```go
// cmd/api/main.go (ตัดตอนเฉพาะ payment)
package main

import (
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"
    
    paymentApp "icmongolang/internal/modules/payment/application"
    paymentHTTP "icmongolang/internal/modules/payment/interfaces/http"
    paymentPG "icmongolang/internal/modules/payment/infrastructure/persistence/postgres"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

func setupPaymentModule(r chi.Router, db *gorm.DB, log logger.Logger, authMW func(http.Handler) http.Handler) {
    // 1. AutoMigrate (dev only)
    if err := db.AutoMigrate(&paymentPG.PaymentModel{}); err != nil {
        log.Fatal("payment migration failed", "error", err)
    }

    // 2. Repository
    paymentRepo := paymentPG.NewPaymentRepository(db)

    // 3. Use Cases
    createPaymentUC := paymentApp.NewCreatePaymentUseCase(paymentRepo, log)

    // 4. Validator
    v := validator.New()

    // 5. Handlers
    paymentHandler := paymentHTTP.NewPaymentHandler(createPaymentUC, v)

    // 6. Routes
    paymentHTTP.RegisterRoutes(r, &paymentHTTP.Handlers{
        Payment: paymentHandler,
    }, authMW)
}
```

---

## 4. Code Template ทุก Layer

### 4.1 Template แบบย่อ (คัดลอกใช้ได้ทันที)

**Domain Entity Template:**
```go
// internal/modules/{MODULE}/domain/entity/{ENTITY}.go
package entity

import (
    "time"
    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/{MODULE}/domain/errors"
)

type {ENTITY} struct {
    ID        uuid.UUID `json:"id"`
    // fields...
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func New{ENTITY}(/* required args */) (*{ENTITY}, error) {
    // Validate invariants
    now := time.Now()
    return &{ENTITY}{
        ID:        uuid.New(),
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

func (e *{ENTITY}) /* BehaviorMethod */() error {
    // Change state through behavior only
    e.UpdatedAt = time.Now()
    return nil
}
```

**Repository Interface Template:**
```go
// internal/modules/{MODULE}/domain/repository/{ENTITY}_repository.go
package repository

import (
    "context"
    "github.com/google/uuid"
    "icmongolang/internal/modules/{MODULE}/domain/entity"
)

type {ENTITY}Repository interface {
    Save(ctx context.Context, e *entity.{ENTITY}) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.{ENTITY}, error)
    Update(ctx context.Context, e *entity.{ENTITY}) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**Use Case Template:**
```go
// internal/modules/{MODULE}/application/{VERB}_{ENTITY}.go
package application

type {VERB}{ENTITY}UseCase struct {
    repo   repository.{ENTITY}Repository
    logger logger.Logger
}

func New{VERB}{ENTITY}UseCase(repo repository.{ENTITY}Repository, log logger.Logger) *{VERB}{ENTITY}UseCase {
    return &{VERB}{ENTITY}UseCase{repo: repo, logger: log}
}

type {VERB}{ENTITY}Input struct { /* fields */ }
type {VERB}{ENTITY}Output struct { /* fields */ }

func (uc *{VERB}{ENTITY}UseCase) Execute(ctx context.Context, input {VERB}{ENTITY}Input) (*{VERB}{ENTITY}Output, error) {
    // 1. Validate
    // 2. Load aggregate
    // 3. Call domain behavior
    // 4. Persist
    // 5. Side effects (audit/event/cache)
    return &{VERB}{ENTITY}Output{}, nil
}
```

---

## 5. Migration & Schema

### 5.1 Naming Convention

| สิ่ง | รูปแบบ | ตัวอย่าง |
|---|---|---|
| Table | `{module}_{entity_plural}` | `payment_transactions` |
| Column | `snake_case` | `user_id`, `created_at` |
| PK | `id` (UUID) | – |
| FK | `{entity}_id` | `order_id` |
| Index | `idx_{table}_{col}` | `idx_payment_transactions_user_id` |
| Unique | `uq_{table}_{col}` | `uq_payment_transactions_order_id` |

### 5.2 Migration Checklist

- [ ] ตั้งชื่อไฟล์ `YYYYMMDD_{module}_{desc}.sql`
- [ ] ทุกตารางมี prefix `{module}_`
- [ ] มี `CREATE TABLE IF NOT EXISTS`
- [ ] มี index สำหรับ column ที่ query บ่อย
- [ ] มี FK constraint (ถ้าอ้างอิงตารางอื่น)
- [ ] มี CHECK constraint (เช่น `amount > 0`)
- [ ] มี trigger สำหรับ `updated_at` (ถ้าจำเป็น)
- [ ] มี rollback script

### 5.3 Rollback Template

```sql
-- migrations/20260401_payment_init.down.sql
DROP TRIGGER IF EXISTS payment_transactions_updated_at ON payment_transactions;
DROP FUNCTION IF EXISTS payment_update_updated_at();
DROP TABLE IF EXISTS payment_transactions;
```

---

## 6. Wire-up & Bootstrap

### 6.1 Module Composition Root

```go
// internal/modules/payment/module.go
package payment

import (
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"
    
    "icmongolang/internal/modules/payment/application"
    "icmongolang/internal/modules/payment/infrastructure/persistence/postgres"
    paymentHTTP "icmongolang/internal/modules/payment/interfaces/http"
    "icmongolang/pkg/logger"
    "icmongolang/pkg/validator"
)

type Module struct {
    db  *gorm.DB
    log logger.Logger
}

func NewModule(db *gorm.DB, log logger.Logger) *Module {
    return &Module{db: db, log: log}
}

func (m *Module) Init(r chi.Router, authMW func(http.Handler) http.Handler) {
    // Repositories
    paymentRepo := postgres.NewPaymentRepository(m.db)

    // Use Cases
    createPaymentUC := application.NewCreatePaymentUseCase(paymentRepo, m.log)

    // Handlers
    v := validator.New()
    paymentHandler := paymentHTTP.NewPaymentHandler(createPaymentUC, v)

    // Routes
    paymentHTTP.RegisterRoutes(r, &paymentHTTP.Handlers{
        Payment: paymentHandler,
    }, authMW)
}
```

### 6.2 การลงทะเบียนใน main.go

```go
// cmd/api/main.go
func main() {
    // ... setup db, logger, router, authMW

    // Register all modules
    payment.NewModule(db, log).Init(r, authMW)
    // order.NewModule(db, log).Init(r, authMW)
    // inventory.NewModule(db, log).Init(r, authMW)
}
```

---

## 7. Checklist ก่อน Commit

### 7.1 Domain Layer
- [ ] Entity มี constructor `New{Entity}`
- [ ] ไม่มี setter ตรง เปลี่ยน state ผ่าน behavior method เท่านั้น
- [ ] Value Object มี `IsValid()` และ logic methods
- [ ] Repository เป็น interface เท่านั้น
- [ ] Domain error เป็น sentinel error
- [ ] Domain **ไม่ import** gorm, gin, chi, sarama, redis

### 7.2 Application Layer
- [ ] Use case ละ 1 ไฟล์ `{verb}_{entity}.go`
- [ ] มี Input/Output DTO แยกชัดเจน
- [ ] `Execute()` return `error` ไม่ panic
- [ ] ไม่มี SQL/HTTP ใน use case

### 7.3 Infrastructure Layer
- [ ] GORM model มี `TableName()` + prefix
- [ ] Repository impl map model ↔ entity ถูกต้อง
- [ ] Kafka message ใช้ JSON
- [ ] Redis key มี namespace `{module}:{entity}:{id}`

### 7.4 Interface Layer
- [ ] Handler ดึง `user_id` จาก context เท่านั้น
- [ ] Route ลงทะเบียนใน `routes.go`
- [ ] Error response เป็น JSON consistent
- [ ] Auth middleware ถูก apply

### 7.5 Build & Test
```bash
go build ./...     # ต้องผ่าน
go vet ./...       # ต้องผ่าน
go test ./...      # ต้องผ่าน
gosec ./...        # security scan
```

### 7.6 Documentation
- [ ] อัปเดต README ของ module
- [ ] เพิ่มตัวอย่าง .env ถ้ามี env ใหม่
- [ ] อัปเดต docker-compose ถ้ามี service ใหม่

---

# 📗 เล่ม 2: คู่มือแก้ไข Module เดิม (Module Modification Guide)

## สารบัญ
1. [Principle: การเปลี่ยนแปลงอย่างปลอดภัย](#1-principle-การเปลี่ยนแปลงอย่างปลอดภัย)
2. [ประเภทของการแก้ไข](#2-ประเภทของการแก้ไข)
3. [Scenario ต่างๆ](#3-scenario-ต่างๆ)
4. [การทำ Migration แบบปลอดภัย](#4-การทำ-migration-แบบปลอดภัย)
5. [การทดสอบก่อน Merge](#5-การทดสอบก่อน-merge)

---

## 1. Principle: การเปลี่ยนแปลงอย่างปลอดภัย

### 1.1 กฎ 5 ข้อก่อนแก้ Module เดิม

```
1. อ่านโค้ดเดิมให้เข้าใจ 100% ก่อนแก้ 1 บรรทัด
2. เขียน test ครอบ behavior เดิมก่อน refactor
3. แก้ทีละเล็กๆ ทีละ commit
4. อย่าแก้ Domain เพื่อให้ Infrastructure ทำงาน
5. ถ้าต้อง breaking change → version API
```

### 1.2 การประเมินผลกระทบ (Impact Analysis)

| คำถาม | ถ้าใช่ → ต้องระวัง |
|---|---|
| มีใครใช้ module นี้อยู่? | ตรวจสอบ call sites ทั้งหมด |
| เปลี่ยน public API หรือไม่? | Version bump |
| เปลี่ยน Schema หรือไม่? | Migration + backward compatible |
| เปลี่ยน business rule หรือไม่? | แจ้ง stakeholder |
| เปลี่ยน event contract หรือไม่? | แจ้ง consumer ทุกตัว |

---

## 2. ประเภทของการแก้ไข

### 2.1 เพิ่ม Field ใหม่ (Add Field)

**Level 1: เพิ่ม field ที่ optional (ปลอดภัย)**

```go
// Before
type Payment struct {
    ID     uuid.UUID
    Amount decimal.Decimal
}

// After
type Payment struct {
    ID       uuid.UUID
    Amount   decimal.Decimal
    Notes    string  // ← เพิ่มใหม่ optional
    Metadata map[string]interface{}  // ← เพิ่มใหม่
}
```

**Migration (backward compatible):**
```sql
-- migrations/20260405_payment_add_notes.sql
ALTER TABLE payment_transactions 
    ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT '',
    ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';
```

**Update GORM Model:**
```go
type PaymentModel struct {
    // ... existing
    Notes    string         `gorm:"type:text;default:''"`
    Metadata datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
}
```

**Update DTO (ทั้ง request/response):**
```go
type CreatePaymentRequest struct {
    // ... existing
    Notes    string                 `json:"notes,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}
```

### 2.2 เพิ่ม Behavior Method ใหม่

```go
// เพิ่ม behavior โดยไม่แก้ของเดิม
func (p *Payment) MarkAsProcessing(processorID string) error {
    if p.Status != valueobject.PaymentStatusPending {
        return domainerrors.ErrInvalidStatusTransition
    }
    p.Status = valueobject.PaymentStatusProcessing
    p.ProcessorID = processorID
    p.UpdatedAt = time.Now()
    return nil
}
```

**เขียน Test ก่อน:**
```go
func TestPayment_MarkAsProcessing(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    
    err := p.MarkAsProcessing("stripe")
    assert.NoError(t, err)
    assert.Equal(t, valueobject.PaymentStatusProcessing, p.Status)
    assert.Equal(t, "stripe", p.ProcessorID)
    
    // Try again → should fail
    err = p.MarkAsProcessing("stripe")
    assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
}
```

### 2.3 เพิ่ม Use Case ใหม่

**ไม่แก้ของเดิม แค่เพิ่ม:**

```go
// internal/modules/payment/application/refund_payment.go
package application

type RefundPaymentUseCase struct {
    repo   repository.PaymentRepository
    logger logger.Logger
}

func NewRefundPaymentUseCase(repo repository.PaymentRepository, log logger.Logger) *RefundPaymentUseCase {
    return &RefundPaymentUseCase{repo: repo, logger: log}
}

type RefundPaymentInput struct {
    PaymentID uuid.UUID
    UserID    uuid.UUID
    Reason    string
}

func (uc *RefundPaymentUseCase) Execute(ctx context.Context, input RefundPaymentInput) error {
    payment, err := uc.repo.FindByID(ctx, input.PaymentID)
    if err != nil {
        return err
    }
    
    // Authorization check
    if payment.UserID != input.UserID {
        return domainerrors.ErrUnauthorized
    }
    
    // Business logic
    if err := payment.Refund(); err != nil {
        return err
    }
    
    return uc.repo.Update(ctx, payment)
}
```

### 2.4 เปลี่ยน Business Rule

**⚠️ ต้องระวังมาก - อาจกระทบผู้ใช้**

**วิธีที่ปลอดภัย - ใช้ Feature Flag:**

```go
// internal/modules/payment/application/create_payment.go
type CreatePaymentUseCase struct {
    repo      repository.PaymentRepository
    logger    logger.Logger
    featureFlag FeatureFlagService  // ← เพิ่มใหม่
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*CreatePaymentOutput, error) {
    // ... existing code
    
    // New business rule behind flag
    // กฎธุรกิจใหม่ภายใต้ feature flag
    if uc.featureFlag.IsEnabled(ctx, "payment.min_amount_check") {
        if input.Amount.LessThan(decimal.NewFromInt(20)) {
            return nil, domainerrors.ErrAmountTooLow
        }
    }
    
    // ... continue
}
```

---

## 3. Scenario ต่างๆ

### 3.1 Scenario: เปลี่ยน Field Name (Rename)

**⚠️ Breaking change - ต้องทำ 3 ขั้น**

**ขั้นที่ 1: เพิ่ม field ใหม่ (ยังไม่ลบของเก่า)**
```go
type Payment struct {
    ID          uuid.UUID
    UserID      uuid.UUID  // เก่า
    CustomerID  uuid.UUID  // ใหม่ ← เพิ่ม
    // ...
}
```

**ขั้นที่ 2: Dual write + backward compat read**
```go
func (r *paymentRepoImpl) toEntity(m *PaymentModel) *entity.Payment {
    // Fallback: read from old if new is empty
    customerID := m.CustomerID
    if customerID == uuid.Nil {
        customerID = m.UserID
    }
    return &entity.Payment{
        CustomerID: customerID,
        // ...
    }
}
```

**ขั้นที่ 3: หลัง deploy 2-3 sprints → drop column เก่า**
```sql
-- migrations/20260601_payment_drop_user_id.sql
ALTER TABLE payment_transactions DROP COLUMN IF EXISTS user_id;
```

### 3.2 Scenario: เปลี่ยน Type ของ Field

**ตัวอย่าง: `amount float64` → `amount decimal.Decimal`**

**ขั้นที่ 1: เพิ่ม column ใหม่ type ถูกต้อง**
```sql
ALTER TABLE payment_transactions 
    ADD COLUMN amount_decimal NUMERIC(15,2);
```

**ขั้นที่ 2: Migrate data**
```sql
UPDATE payment_transactions 
    SET amount_decimal = amount::numeric(15,2)
    WHERE amount_decimal IS NULL;
```

**ขั้นที่ 3: Switch code ใช้ column ใหม่ → drop column เก่า**

### 3.3 Scenario: เพิ่ม Relationship

**ตัวอย่าง: Payment → มีได้หลาย Refund**

```go
// เพิ่ม Entity ใหม่
type Refund struct {
    ID        uuid.UUID
    PaymentID uuid.UUID
    Amount    decimal.Decimal
    Reason    string
    CreatedAt time.Time
}

// เพิ่มความสัมพันธ์
type Payment struct {
    // ... existing
    Refunds []Refund  // ← เพิ่มใหม่
}
```

**เพิ่ม Repository:**
```go
type RefundRepository interface {
    Save(ctx context.Context, r *entity.Refund) error
    FindByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*entity.Refund, error)
}
```

### 3.4 Scenario: เปลี่ยนจาก Sync → Async (Event-Driven)

**Before:**
```go
func (uc *CreatePaymentUseCase) Execute(ctx, input) (*Output, error) {
    // ... create payment
    // ... send email directly (blocking!)
    uc.emailService.Send(ctx, payment.UserID, "payment.created")
    return output, nil
}
```

**After:**
```go
func (uc *CreatePaymentUseCase) Execute(ctx, input) (*Output, error) {
    // ... create payment
    
    // Publish event instead
    event := events.PaymentCreatedEvent{
        PaymentID: payment.ID,
        UserID:    payment.UserID,
        Amount:    payment.Amount,
    }
    if err := uc.eventBus.Publish(ctx, "payment.created", event); err != nil {
        uc.logger.Warn("failed to publish event", "error", err)
        // ไม่ fail operation หลัก
    }
    return output, nil
}
```

**เพิ่ม Consumer:**
```go
// internal/modules/payment/infrastructure/messaging/consumers/payment_created_consumer.go
type PaymentCreatedHandler struct {
    emailService EmailService
}

func (h *PaymentCreatedHandler) Handle(ctx context.Context, payload map[string]interface{}) error {
    userID := payload["user_id"].(string)
    return h.emailService.Send(ctx, userID, "payment.created")
}
```

---

## 4. การทำ Migration แบบปลอดภัย

### 4.1 หลักการ Expand-Contract

```
Phase 1: EXPAND   → เพิ่มของใหม่ (ยัง backward compat)
Phase 2: MIGRATE  → ย้ายข้อมูล + switch code
Phase 3: CONTRACT → ลบของเก่า
```

### 4.2 Zero-Downtime Migration Checklist

- [ ] ไม่ lock table นาน (หลีกเลี่ยง `ALTER TABLE ... REWRITE`)
- [ ] ใช้ `CREATE INDEX CONCURRENTLY`
- [ ] Backfill ด้วย batch (ไม่ UPDATE ทั้งตารางทีเดียว)
- [ ] Test บน staging ที่มี data ≈ production
- [ ] มี rollback plan ที่ verified

### 4.3 ตัวอย่าง: เพิ่ม NOT NULL Column ปลอดภัย

```sql
-- Step 1: Add nullable + default
ALTER TABLE payment_transactions
    ADD COLUMN processor VARCHAR(50) DEFAULT 'legacy';

-- Step 2: Backfill batch
DO $$
DECLARE
    batch_size INT := 1000;
BEGIN
    LOOP
        UPDATE payment_transactions
        SET processor = 'unknown'
        WHERE processor IS NULL
          AND id IN (
              SELECT id FROM payment_transactions
              WHERE processor IS NULL LIMIT batch_size
          );
        EXIT WHEN NOT FOUND;
        COMMIT;
        PERFORM pg_sleep(0.1);  -- ให้ DB พัก
    END LOOP;
END $$;

-- Step 3: Add NOT NULL constraint
ALTER TABLE payment_transactions
    ALTER COLUMN processor SET NOT NULL,
    ALTER COLUMN processor DROP DEFAULT;
```

---

## 5. การทดสอบก่อน Merge

### 5.1 Test Pyramid

```
         /\
        /  \  E2E (น้อย)
       /────\
      /      \  Integration (กลาง)
     /────────\
    /          \  Unit (เยอะ)
   /────────────\
```

### 5.2 Unit Test Template

```go
// internal/modules/payment/domain/entity/payment_test.go
package entity_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "icmongolang/internal/modules/payment/domain/entity"
    "icmongolang/internal/modules/payment/domain/value_object"
)

func TestNewPayment_Valid(t *testing.T) {
    userID := uuid.New()
    orderID := uuid.New()
    
    p, err := entity.NewPayment(userID, orderID, decimal.NewFromInt(100), "THB", "CARD")
    
    require.NoError(t, err)
    assert.Equal(t, userID, p.UserID)
    assert.Equal(t, orderID, p.OrderID)
    assert.Equal(t, valueobject.PaymentStatusPending, p.Status)
}

func TestNewPayment_InvalidAmount(t *testing.T) {
    _, err := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(-10), "THB", "CARD")
    assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount)
}

func TestPayment_StateTransition_Valid(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    
    require.NoError(t, p.StartProcessing())
    assert.Equal(t, valueobject.PaymentStatusProcessing, p.Status)
    
    require.NoError(t, p.MarkSuccess())
    assert.Equal(t, valueobject.PaymentStatusSuccess, p.Status)
}

func TestPayment_StateTransition_Invalid(t *testing.T) {
    p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    
    // Cannot go from PENDING to SUCCESS directly
    err := p.MarkSuccess()
    assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
}
```

### 5.3 Integration Test Template

```go
//go:build integration

package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func TestPaymentRepository_SaveAndFind(t *testing.T) {
    // Setup test container
    ctx := context.Background()
    pgContainer, err := postgres.RunContainer(ctx, /* ... */)
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    // Connect
    dsn, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    db.AutoMigrate(&PaymentModel{})

    // Test
    repo := NewPaymentRepository(db)
    payment, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    
    require.NoError(t, repo.Save(ctx, payment))
    
    found, err := repo.FindByID(ctx, payment.ID)
    require.NoError(t, err)
    assert.Equal(t, payment.ID, found.ID)
}
```

---

# 📘 เล่ม 3: คู่มือขาย Module (Module Commercialization Guide)

## สารบัญ
1. [โมเดลธุรกิจ](#1-โมเดลธุรกิจ)
2. [การเตรียม Module ให้พร้อมขาย](#2-การเตรียม-module-ให้พร้อมขาย)
3. [Package & Distribution](#3-package--distribution)
4. [Licensing & Pricing](#4-licensing--pricing)
5. [การตลาดและ Sales Kit](#5-การตลาดและ-sales-kit)
6. [Support & SLA](#6-support--sla)

---

## 1. โมเดลธุรกิจ

### 1.1 รูปแบบการขาย Module Go

| รูปแบบ | ตัวอย่าง | รายได้ |
|---|---|---|
| **Library/Module** | ขาย Go package | License fee ต่อปี |
| **SaaS Product** | ใช้ module เป็น backend | Subscription |
| **White-label** | ให้ partner ใช้แบรนด์ตัวเอง | Revenue share |
| **Source Code** | ขาย source ทั้งชุด | One-time + support |
| **Consulting** | ให้คำปรึกษา | Hourly/Daily rate |

### 1.2 Pricing Tiers

| Tier | ราคา/เดือน | เหมาะกับ |
|---|---|---|
| **Starter** | ฿990 | 1 tenant, 10 users |
| **Pro** | ฿2,990 | 5 tenants, 50 users |
| **Business** | ฿9,990 | 20 tenants, 200 users |
| **Enterprise** | Custom | Unlimited + SLA 99.9% |

### 1.3 Value Proposition Canvas

```
         PAINS                    GAINS
    ┌──────────────┐         ┌──────────────┐
    │ เขียนเอง 6 เดือน│         │ ใช้ได้เลย     │
    │ บั๊กเยอะ      │         │ Test ครอบคลุม │
    │ ไม่มี docs     │         │ มี docs       │
    │ Maintain ยาก   │         │ Update ให้    │
    └──────────────┘         └──────────────┘
              ↓                     ↑
         ┌────────────────────────────┐
         │     YOUR MODULE            │
         │  • Clean Architecture      │
         │  • Production-ready        │
         │  • Full documentation      │
         │  • Active maintenance      │
         └────────────────────────────┘
```

---

## 2. การเตรียม Module ให้พร้อมขาย

### 2.1 Code Quality Checklist

| # | รายการ | เกณฑ์ |
|---|---|---|
| 1 | Test Coverage | ≥ 80% |
| 2 | Lint | ผ่าน golangci-lint |
| 3 | Security | ผ่าน gosec |
| 4 | Documentation | ทุก public symbol |
| 5 | Examples | มี runnable examples |
| 6 | Benchmarks | มี performance data |
| 7 | Changelog | ตาม Semantic Versioning |
| 8 | License | ระบุชัดเจน |

### 2.2 โครงสร้าง Package ที่ขายได้

```
payment-module/
├── README.md                    # ← อธิบายคุณค่า ไม่ใช่แค่ API
├── LICENSE                      # ← MIT / Apache / Commercial
├── CHANGELOG.md                 # ← Keep a Changelog format
├── go.mod                       # ← module path
├── domain/                      # ← public API
│   ├── entity/
│   ├── value_object/
│   └── repository/
├── application/                 # ← use cases
├── infrastructure/              # ← adapters (postgres, redis, kafka)
├── interfaces/                  # ← HTTP/gRPC handlers
├── examples/                    # ← รันได้จริง
│   ├── basic/
│   └── with-kafka/
├── docs/                        # ← detailed docs
│   ├── getting-started.md
│   ├── architecture.md
│   └── api-reference.md
├── migrations/                  # ← SQL migrations
├── Makefile
└── docker-compose.yml           # ← demo environment
```

### 2.3 README Template สำหรับขาย

```markdown
# 💳 Payment Module for Go

> Production-ready payment module with Clean Architecture + DDD

## ✨ Features
- ✅ Multi-provider support (Stripe, Omise, 2C2P)
- ✅ Full refund workflow
- ✅ Webhook handling with idempotency
- ✅ Event-driven with Kafka
- ✅ Multi-currency support
- ✅ 90% test coverage
- ✅ Production-tested (10M+ transactions)

## 🚀 Quick Start (5 minutes)

\`\`\`bash
go get github.com/yourorg/payment-module
\`\`\`

\`\`\`go
import "github.com/yourorg/payment-module"

paymentModule := payment.NewModule(db, eventBus, logger)
paymentModule.Init(router, authMiddleware)
\`\`\`

## 📊 Performance

| Metric | Value |
|---|---|
| Throughput | 5,000 TPS |
| P95 Latency | 45ms |
| Memory | 120 MB |

## 💰 Pricing

| Tier | Price | Users |
|---|---|---|
| Starter | $29/mo | 10 |
| Pro | $99/mo | 100 |
| Enterprise | Contact | Unlimited |

## 📚 Documentation
- [Getting Started](docs/getting-started.md)
- [Architecture](docs/architecture.md)
- [API Reference](docs/api-reference.md)

## 📞 Contact
- Email: sales@yourorg.com
- Website: https://yourorg.com
```

---

## 3. Package & Distribution

### 3.1 Private Go Module

```bash
# ตั้ง GOPRIVATE
export GOPRIVATE=github.com/yourorg/*

# Consumer ต้องมี access token
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

### 3.2 Versioning Strategy

```
v1.0.0    → Initial release
v1.1.0    → New feature (backward compatible)
v1.1.1    → Bug fix
v2.0.0    → Breaking change
```

**Semantic Import Versioning:**
```go
// v1
import "github.com/yourorg/payment-module"

// v2 (breaking change)
import "github.com/yourorg/payment-module/v2"
```

### 3.3 Docker Distribution

```dockerfile
# Dockerfile for payment module as a service
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /payment-service ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /payment-service /payment-service
EXPOSE 8080
ENTRYPOINT ["/payment-service"]
```

---

## 4. Licensing & Pricing

### 4.1 License Options

| License | ใช้เมื่อ | ข้อดี | ข้อเสีย |
|---|---|---|---|
| **MIT** | Open source, ให้ใช้ฟรี | ง่าย, adopt เยอะ | ไม่มีรายได้ |
| **Apache 2.0** | Open source + patent protection | Legal ชัดเจน | ไม่มีรายได้ |
| **AGPL** | Open source + copyleft | ดึงดูด contributor | องค์กรหลีกเลี่ยง |
| **Commercial** | ขาย | มีรายได้ | ยาก adopt |
| **Dual** | ทั้ง 2 แบบ | ยืดหยุ่น | ซับซ้อน |

### 4.2 Commercial License Template

```
PAYMENT MODULE COMMERCIAL LICENSE AGREEMENT

1. LICENSE GRANT
   Licensor grants Licensee a non-exclusive, non-transferable license
   to use the Payment Module for commercial purposes.

2. RESTRICTIONS
   - No redistribution of source code
   - No sublicensing without permission
   - One production deployment per license

3. SUPPORT
   - Business hours: 9:00-18:00 ICT
   - Response time: 4 hours (Business tier)
   - Response time: 1 hour (Enterprise tier)

4. FEES
   - Annual license fee: $X,XXX
   - Renewal: 20% discount

5. TERM
   12 months from purchase date.

6. TERMINATION
   Either party may terminate with 30 days notice.
```

### 4.3 ROI Calculator Template (สำหรับลูกค้า)

```markdown
# ROI Analysis: Payment Module vs Build In-House

## Build In-House
| Item | Cost |
|---|---|
| Developer (3 months) | $45,000 |
| QA (1 month) | $8,000 |
| Infrastructure setup | $5,000 |
| **Total** | **$58,000** |
| Time to market | 4 months |

## Buy Payment Module
| Item | Cost |
|---|---|
| License (1 year) | $1,200 |
| Integration (2 weeks) | $6,000 |
| **Total** | **$7,200** |
| Time to market | 2 weeks |

## ROI
- **Savings: $50,800/year**
- **Payback: 0.15 months**
- **ROI: 705%**
```

---

## 5. การตลาดและ Sales Kit

### 5.1 Sales Kit Checklist

- [ ] Product one-pager (PDF)
- [ ] Demo video (5 minutes)
- [ ] Interactive demo (Docker)
- [ ] Case study (2-3 ตัวอย่าง)
- [ ] ROI calculator (Excel/Web)
- [ ] Pricing sheet
- [ ] FAQ document
- [ ] Technical whitepaper

### 5.2 Product One-Pager Template

```
┌────────────────────────────────────────────┐
│  💳 PAYMENT MODULE FOR GO                  │
│                                            │
│  Production-ready payment processing       │
│  for modern Go applications                │
├────────────────────────────────────────────┤
│  ✅ 6 months saved on development          │
│  ✅ 90% test coverage                      │
│  ✅ Multi-provider (Stripe/Omise/2C2P)     │
│  ✅ Event-driven architecture              │
│  ✅ Multi-currency support                 │
├────────────────────────────────────────────┤
│  📊 Performance                            │
│  • 5,000 TPS                               │
│  • P95 < 50ms                              │
│  • 99.99% uptime                           │
├────────────────────────────────────────────┤
│  💰 Pricing                                │
│  • Starter: $29/mo                         │
│  • Pro: $99/mo                             │
│  • Enterprise: Contact                     │
├────────────────────────────────────────────┤
│  📞 Contact                                │
│  sales@yourorg.com                        │
│  https://yourorg.com                      │
└────────────────────────────────────────────┘
```

### 5.3 Demo Environment

```yaml
# docker-compose.demo.yml
version: '3.9'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: payment_demo
      POSTGRES_USER: demo
      POSTGRES_PASSWORD: demo
    ports:
      - "5432:5432"
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  payment-api:
    image: yourorg/payment-module:latest
    ports:
      - "8080:8080"
    environment:
      DB_DSN: postgres://demo:demo@postgres:5432/payment_demo
      REDIS_ADDR: redis:6379
      STRIPE_API_KEY: sk_test_xxx
    depends_on:
      - postgres
      - redis
  
  payment-ui:
    image: yourorg/payment-demo-ui:latest
    ports:
      - "3000:3000"
    depends_on:
      - payment-api
```

---

## 6. Support & SLA

### 6.1 Support Tiers

| Tier | Response Time | Channels | Hours |
|---|---|---|---|
| **Community** | Best effort | GitHub Issues | – |
| **Business** | 4 hours | Email, Slack | 9-18 ICT |
| **Enterprise** | 1 hour | Email, Slack, Phone | 24/7 |

### 6.2 SLA Template

```
SERVICE LEVEL AGREEMENT

1. UPTIME COMMITMENT
   • 99.9% monthly uptime (Business)
   • 99.99% monthly uptime (Enterprise)

2. RESPONSE TIME
   Severity 1 (Critical): 1 hour
   Severity 2 (High): 4 hours
   Severity 3 (Normal): 1 business day
   Severity 4 (Low): 3 business days

3. RESOLUTION TARGETS
   Severity 1: 8 hours
   Severity 2: 2 business days
   Severity 3: 5 business days
   Severity 4: Next release

4. CREDITS
   Downtime 99.0-99.9%: 10% credit
   Downtime 95.0-99.0%: 25% credit
   Downtime < 95.0%: 50% credit
```

---

# 📙 เล่ม 4: คู่มือทดสอบและ Deployment

## สารบัญ
1. [Testing Strategy](#1-testing-strategy)
2. [CI/CD Pipeline](#2-cicd-pipeline)
3. [Deployment Patterns](#3-deployment-patterns)
4. [Monitoring & Observability](#4-monitoring--observability)

---

## 1. Testing Strategy

### 1.1 Test Pyramid สำหรับ Go

```
         /\
        /E2E\        5%   - Full stack
       /──────\
      / Integr \    15%   - Repository, DB, Kafka
     /──────────\
    /   Unit     \ 80%   - Domain, Use Case
   /──────────────\
```

### 1.2 Unit Test - Domain Layer

```go
package entity_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPayment_StateMachine(t *testing.T) {
    tests := []struct {
        name          string
        setup         func() *entity.Payment
        action        func(*entity.Payment) error
        wantErr       error
        wantStatus    valueobject.PaymentStatus
    }{
        {
            name: "pending to processing",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
                return p
            },
            action:     func(p *entity.Payment) error { return p.StartProcessing() },
            wantErr:    nil,
            wantStatus: valueobject.PaymentStatusProcessing,
        },
        {
            name: "cannot success from pending",
            setup: func() *entity.Payment {
                p, _ := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
                return p
            },
            action:     func(p *entity.Payment) error { return p.MarkSuccess() },
            wantErr:    domainerrors.ErrInvalidStatusTransition,
            wantStatus: valueobject.PaymentStatusPending,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := tt.setup()
            err := tt.action(p)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.wantStatus, p.Status)
        })
    }
}
```

### 1.3 Mock Repository

```go
// internal/modules/payment/application/mocks/payment_repo_mock.go
package mocks

import (
    "context"
    "github.com/google/uuid"
    "github.com/stretchr/testify/mock"
    "icmongolang/internal/modules/payment/domain/entity"
)

type PaymentRepositoryMock struct {
    mock.Mock
}

func (m *PaymentRepositoryMock) Save(ctx context.Context, p *entity.Payment) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func (m *PaymentRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Payment), args.Error(1)
}

// ... etc
```

### 1.4 Use Case Test with Mock

```go
func TestCreatePaymentUseCase_Success(t *testing.T) {
    repo := new(mocks.PaymentRepositoryMock)
    log := logger.NewNoop()
    uc := application.NewCreatePaymentUseCase(repo, log)
    
    userID := uuid.New()
    orderID := uuid.New()
    
    repo.On("FindByOrderID", mock.Anything, orderID).
        Return(nil, domainerrors.ErrPaymentNotFound)
    repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Payment")).
        Return(nil)
    
    output, err := uc.Execute(context.Background(), application.CreatePaymentInput{
        UserID:   userID,
        OrderID:  orderID,
        Amount:   decimal.NewFromInt(100),
        Currency: "THB",
        Method:   "CARD",
    })
    
    require.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, output.ID)
    assert.Equal(t, "PENDING", output.Status)
    repo.AssertExpectations(t)
}
```

### 1.5 Integration Test with Testcontainers

```go
//go:build integration

func TestPaymentRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    ctx := context.Background()
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).WithStartupTimeout(30*time.Second)),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)
    
    dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    require.NoError(t, err)
    
    require.NoError(t, db.AutoMigrate(&PaymentModel{}))
    
    repo := NewPaymentRepository(db)
    
    // Test create
    payment, err := entity.NewPayment(uuid.New(), uuid.New(), decimal.NewFromInt(100), "THB", "CARD")
    require.NoError(t, err)
    require.NoError(t, repo.Save(ctx, payment))
    
    // Test find
    found, err := repo.FindByID(ctx, payment.ID)
    require.NoError(t, err)
    assert.Equal(t, payment.ID, found.ID)
    assert.Equal(t, payment.Amount.String(), found.Amount.String())
}
```

---

## 2. CI/CD Pipeline

### 2.1 GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

env:
  GO_VERSION: '1.22'

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          args: --timeout=5m

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports: ['5432:5432']
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7-alpine
        ports: ['6379:6379']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Download deps
        run: go mod download
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
        env:
          DB_DSN: postgres://postgres:test@localhost:5432/testdb?sslmode=disable
          REDIS_ADDR: localhost:6379
      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage below 80%"
            exit 1
          fi
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec.sarif ./...'
      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec.sarif
      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'HIGH,CRITICAL'
          exit-code: '1'

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, test, security]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=linux go build \
            -ldflags="-w -s -X main.version=${{ github.sha }}" \
            -o bin/api ./cmd/api
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: api-binary
          path: bin/api
```

### 2.2 Deployment Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]
    tags: ['v*']

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: docker/setup-buildx-action@v3
      
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - uses: docker/metadata-action@v5
        id: meta
        with:
          images: ghcr.io/${{ github.repository }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=sha,prefix=
      
      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
  
  deploy-staging:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: staging
    steps:
      - name: Deploy to staging
        run: |
          kubectl set image deployment/api \
            api=ghcr.io/${{ github.repository }}:${{ github.sha }} \
            -n staging
          kubectl rollout status deployment/api -n staging
  
  deploy-production:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/v')
    environment: production
    steps:
      - name: Deploy to production (canary)
        run: |
          kubectl set image deployment/api-canary \
            api=ghcr.io/${{ github.repository }}:${{ github.sha }} \
            -n production
          kubectl rollout status deployment/api-canary -n production
      
      - name: Wait for canary validation
        run: sleep 300
      
      - name: Promote to stable
        run: |
          kubectl set image deployment/api \
            api=ghcr.io/${{ github.repository }}:${{ github.sha }} \
            -n production
          kubectl rollout status deployment/api -n production
```

---

## 3. Deployment Patterns

### 3.1 Deployment Strategies

| Strategy | Risk | Downtime | Rollback |
|---|---|---|---|
| **Recreate** | High | Yes | Slow |
| **Rolling Update** | Medium | No | Medium |
| **Blue-Green** | Low | No | Fast |
| **Canary** | Very Low | No | Very Fast |

### 3.2 Kubernetes Deployment

```yaml
# deployments/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-api
  namespace: production
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: payment-api
  template:
    metadata:
      labels:
        app: payment-api
    spec:
      containers:
      - name: api
        image: ghcr.io/yourorg/payment-module:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_DSN
          valueFrom:
            secretKeyRef:
              name: payment-secrets
              key: db-dsn
        - name: REDIS_ADDR
          value: "redis-service:6379"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
---
apiVersion: v1
kind: Service
metadata:
  name: payment-api
  namespace: production
spec:
  selector:
    app: payment-api
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: payment-api-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: payment-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## 4. Monitoring & Observability

### 4.1 Three Pillars

| Pillar | Tool | Purpose |
|---|---|---|
| **Metrics** | Prometheus + Grafana | Numbers over time |
| **Logs** | Loki / ELK | Events |
| **Traces** | Jaeger / Tempo | Request flow |

### 4.2 Metrics Instrumentation

```go
// pkg/metrics/payment.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    PaymentCreatedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "payment_created_total",
            Help: "Total number of payments created",
        },
        []string{"method", "currency", "status"},
    )
    
    PaymentDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "payment_duration_seconds",
            Help:    "Payment processing duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method"},
    )
)

// Usage in use case
func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*Output, error) {
    timer := prometheus.NewTimer(PaymentDuration.WithLabelValues(input.Method))
    defer timer.ObserveDuration()
    
    // ... business logic
    
    PaymentCreatedTotal.WithLabelValues(
        input.Method, input.Currency, string(payment.Status),
    ).Inc()
    
    return output, nil
}
```

### 4.3 Structured Logging

```go
// pkg/logger/logger.go
package logger

import (
    "log/slog"
    "os"
)

type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
}

type slogLogger struct {
    l *slog.Logger
}

func New() Logger {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
    return &slogLogger{l: slog.New(handler)}
}

func (s *slogLogger) Info(msg string, args ...any) {
    s.l.Info(msg, args...)
}
// ...
```

**Usage:**
```go
uc.logger.Info("payment created",
    "payment_id", payment.ID,
    "user_id", payment.UserID,
    "amount", payment.Amount.String(),
    "currency", payment.Currency,
)
```

---

# 📕 เล่ม 5: คู่มือบำรุงรักษาและ Scale

## สารบัญ
1. [Maintenance Schedule](#1-maintenance-schedule)
2. [Performance Tuning](#2-performance-tuning)
3. [Scaling Strategies](#3-scaling-strategies)
4. [Incident Response (RCA)](#4-incident-response-rca)

---

## 1. Maintenance Schedule

### 1.1 ตารางการบำรุงรักษา

| รอบ | งาน | เครื่องมือ |
|---|---|---|
| **Daily** | ตรวจ health check, alert log | Grafana, PagerDuty |
| **Weekly** | Update dependencies, review security | Dependabot |
| **Monthly** | Backup restore test, cert rotation | pg_dump, Vault |
| **Quarterly** | Load test, capacity review | k6, Grafana |
| **Yearly** | Major version upgrade, DR drill | – |

### 1.2 Daily Checklist

```bash
#!/bin/bash
# scripts/daily-check.sh

echo "=== Health Check ==="
curl -s http://localhost:8080/health/detailed | jq

echo "=== Error Rate (last 24h) ==="
curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~\"5..\"}[24h])" | jq

echo "=== DB Connection Pool ==="
psql -c "SELECT count(*) FROM pg_stat_activity WHERE state = 'active';"

echo "=== Disk Usage ==="
df -h | grep -E "(/|/var)"

echo "=== Backup Status ==="
ls -la /backups/ | tail -5
```

---

## 2. Performance Tuning

### 2.1 Database Optimization

**Slow Query Detection:**
```sql
-- Enable pg_stat_statements
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Top 10 slow queries
SELECT 
    query,
    calls,
    mean_exec_time,
    total_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Index Analysis:**
```sql
-- Missing indexes (sequential scans)
SELECT 
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    idx_tup_fetch
FROM pg_stat_user_tables
WHERE seq_scan > 1000
ORDER BY seq_scan DESC;

-- Unused indexes
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
  AND indexname NOT LIKE '%_pkey';
```

### 2.2 Application Tuning

```go
// Connection pool tuning
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)      // ปรับตาม DB max_connections
sqlDB.SetMaxIdleConns(10)      // ~40% ของ MaxOpen
sqlDB.SetConnMaxLifetime(30 * time.Minute)
sqlDB.SetConnMaxIdleTime(5 * time.Minute)
```

### 2.3 Cache Strategy

```go
// Cache-Aside Pattern
func (uc *GetPaymentUseCase) Execute(ctx context.Context, id uuid.UUID) (*Payment, error) {
    cacheKey := fmt.Sprintf("payment:txn:%s", id)
    
    // 1. Try cache
    var payment entity.Payment
    if err := uc.cache.Get(ctx, cacheKey, &payment); err == nil {
        return &payment, nil
    }
    
    // 2. Cache miss → DB
    p, err := uc.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. Save to cache (TTL 10 min)
    _ = uc.cache.Set(ctx, cacheKey, p, 10*time.Minute)
    
    return p, nil
}
```

---

## 3. Scaling Strategies

### 3.1 Scaling Progression

```
Phase 1: Vertical      Phase 2: Read Replica    Phase 3: Sharding
   ↓                      ↓                        ↓
[Bigger DB]           [Master + Replicas]      [Shard by tenant_id]
```

### 3.2 Read Replica Pattern

```go
// Repository ที่แยก read/write
type paymentRepoImpl struct {
    writeDB *gorm.DB  // master
    readDB  *gorm.DB  // replica
}

func (r *paymentRepoImpl) Save(ctx context.Context, p *entity.Payment) error {
    return r.writeDB.WithContext(ctx).Create(r.toModel(p)).Error
}

func (r *paymentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
    var m PaymentModel
    err := r.readDB.WithContext(ctx).First(&m, "id = ?", id).Error
    // ...
}
```

### 3.3 Sharding by Tenant

```go
// Shard router
type ShardRouter struct {
    shards map[string]*gorm.DB  // tenantID → DB
}

func (r *ShardRouter) GetDB(tenantID string) *gorm.DB {
    // Hash tenant ID to shard
    hash := fnv.New32a()
    hash.Write([]byte(tenantID))
    shardIdx := hash.Sum32() % uint32(len(r.shards))
    return r.shards[fmt.Sprintf("shard-%d", shardIdx)]
}
```

---

## 4. Incident Response (RCA)

### 4.1 Incident Severity Levels

| Level | Impact | Response | Example |
|---|---|---|---|
| **P0** | System down | Immediate, all hands | Payment API 500 |
| **P1** | Major degradation | 15 min response | Latency > 5s |
| **P2** | Minor issue | 1 hour response | One feature slow |
| **P3** | Cosmetic | Next sprint | Wrong color |

### 4.2 RCA Process

```
1. Detect      → Alert / User report
2. Triage      → Assess severity
3. Mitigate    → Stop the bleeding
4. Investigate → Find root cause
5. Fix         → Permanent solution
6. Prevent     → Add safeguards
7. Document    → Write RCA report
```

### 4.3 RCA Report Template

```markdown
# RCA Report: [INCIDENT-ID]

## Summary
- **Date**: 2026-04-01 14:30 ICT
- **Duration**: 45 minutes
- **Severity**: P1
- **Impact**: 500 users, ฿150,000 revenue loss

## Timeline
| Time | Event |
|---|---|
| 14:30 | Alert: P95 latency > 5s |
| 14:32 | On-call engineer acknowledged |
| 14:35 | Identified slow query |
| 14:45 | Added missing index |
| 15:15 | System recovered |

## Root Cause
Missing index on `payment_transactions.user_id` caused full table scan.

## 5 Whys
1. Why slow? → Query takes 4s
2. Why slow query? → Full table scan
3. Why full scan? → No index
4. Why no index? → Migration forgot
5. Why forgot? → No checklist

## Action Items
| # | Action | Owner | Due |
|---|---|---|---|
| 1 | Add missing index | DBA | 2026-04-01 |
| 2 | Update migration | Dev | 2026-04-02 |
| 3 | Add migration checklist | Lead | 2026-04-03 |
| 4 | Add slow query alert | DevOps | 2026-04-05 |

## Lessons Learned
- Migration ต้องมี index review
- ต้องมี query performance test ใน CI
- ต้องมี alert สำหรับ slow query

## Prevention
- Automated index detection
- Query plan analysis ใน CI
- Performance budget ใน test
```

### 4.4 Post-Mortem Template (Blameless)

```markdown
# Post-Mortem: [INCIDENT-ID]

## What happened?
[อธิบายเหตุการณ์]

## Impact
[ผลกระทบต่อผู้ใช้ ธุรกิจ]

## Root Cause
[สาเหตุที่แท้จริง]

## What went well?
- Alert ทำงานถูกต้อง
- ทีม respond เร็ว

## What went wrong?
- ไม่มี index
- ไม่มี monitoring ที่ดีพอ

## Action Items
[รายการ + owner + due date]

## Lessons Learned
[บทเรียน]
```

---

# 🎯 สรุป - แผนที่การใช้งานคู่มือทั้ง 5 เล่ม

| สถานการณ์ | เล่มที่ต้องอ่าน |
|---|---|
| เริ่มโปรเจกต์ใหม่ | เล่ม 1 → 2 → 4 |
| รับงาน module เดิม | เล่ม 2 → 4 → 5 |
| ต้องการขาย module | เล่ม 3 → 1 → 2 |
| ระบบมีปัญหา | เล่ม 5 |
| Onboarding ทีมใหม่ | เล่ม 1 → 2 |

---

## 🎁 Bonus: Quick Reference Card

```
┌──────────────────────────────────────────────────┐
│         GOLANG MODULE CHEAT SHEET                │
├──────────────────────────────────────────────────┤
│ 📁 STRUCTURE                                     │
│   domain/        ← pure business logic           │
│   application/   ← use cases                     │
│   infrastructure/← adapters                      │
│   interfaces/    ← HTTP/gRPC                     │
├──────────────────────────────────────────────────┤
│ 🔨 CREATE MODULE                                 │
│   mkdir internal/modules/{name}/...              │
│   domain/entity, value_object, repository        │
│   application/{verb}_{entity}.go                 │
│   infrastructure/persistence/postgres/           │
│   interfaces/http/handler.go, routes.go          │
│   migrations/YYYYMMDD_{module}_init.sql          │
│   wire-up ใน main.go                             │
├──────────────────────────────────────────────────┤
│ ✅ BEFORE COMMIT                                 │
│   go build ./...                                 │
│   go vet ./...                                   │
│   go test -race -cover ./...                     │
│   gosec ./...                                    │
├──────────────────────────────────────────────────┤
│ 🚀 DEPLOY                                        │
│   git tag v1.0.0                                 │
│   git push --tags                                │
│   → CI/CD auto-deploy                            │
└──────────────────────────────────────────────────┘
```

---

**จัดทำโดย:** ทีมสถาปัตยกรรมซอฟต์แวร์
**เวอร์ชัน:** 1.0 (เมษายน 2026)
**อ้างอิง:** โปรเจกต์ `icmongolang` + Clean Architecture + DDD