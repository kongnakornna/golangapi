# 🏗️ ระบบให้บริการ IoT Platform (Smart Farm / Smart Building)
## Master Design Document — ตามโครงสร้าง `icmongolang` + Template_Module.md (DDD + Clean Architecture)

---

# สารบัญภาพรวม

| # | Module | โฟลเดอร์ | หน้าที่ |
|---|---|---|---|
| 1 | ระบบลูกค้า | `internal/modules/customer` | Customer 360, Site, Contract |
| 2 | Package & Service | `internal/modules/packagecatalog` | Catalog, Subscription, Quota |
| 3 | ERP IoT Solution | `internal/modules/erp` | Product, Inventory, PO/SO, Invoice |
| 4 | CRM IoT Solution | `internal/modules/crm` | Lead, Opportunity, Ticket |
| 5 | จัดการอุปกรณ์ IoT/AI/Automation | `internal/modules/iotdevice` | Device, Telemetry, Command, AI Rules, Automation |
| 6 | บริหารจัดการ Logistic IoT | `internal/modules/iotlogistics` | Shipment, Installation, Maintenance, RMA |
| 7 | ระบบรายงาน | `internal/modules/report` | Report, KPI, Dashboard |

> ทุกโมดูลใช้ **โครงสร้าง 17 หัวข้อ** เดียวกัน ตามที่กำหนด

---

# PART 0 — SYSTEM OVERVIEW (ภาพรวมทั้งแพลตฟอร์ม)

## 0.1 สถาปัตยกรรมโดยรวม

```
┌─────────────────────────────────────────────────────────────────┐
│                    IoT Service Platform (SaaS)                    │
│                                                                   │
│  ┌──────────┐ ┌──────────────┐ ┌──────┐ ┌──────┐                │
│  │ Customer │ │PackageService│ │ ERP  │ │ CRM  │                │
│  └────┬─────┘ └──────┬───────┘ └──┬───┘ └──┬───┘                │
│       │              │            │        │                     │
│  ┌────┴──────────────┴────────────┴────────┴────┐               │
│  │       Shared Kernel (tenant, audit, iam)      │               │
│  └────┬──────────────┬────────────┬────────────┬─┘               │
│       │              │            │            │                 │
│  ┌────┴────┐ ┌───────┴──────┐ ┌───┴──────┐ ┌───┴─────┐          │
│  │IoTDevice│ │IoTLogistics  │ │  Report  │ │Realtime │          │
│  │+AI+Auto │ │              │ │          │ │  (WS)   │          │
│  └────┬────┘ └──────┬───────┘ └──────────┘ └─────────┘          │
└───────┼─────────────┼───────────────────────────────────────────┘
        │             │
   MQTT/Kafka    Kafka/REST
        │             │
   [Devices]    [Mobile/Web/Technician]
```

## 0.2 Integration Matrix (Kafka Topics ทั้งระบบ)

| Producer | Topic | Consumer |
|---|---|---|
| `crm` | `crm.lead.converted` | customer, package |
| `customer` | `customer.created` | crm, package, erp |
| `packagecatalog` | `package.subscription.created` | payment, iotdevice, notifier |
| `erp` | `erp.invoice.issued` | payment, report |
| `iotdevice` | `iot.alert.triggered` | notifier, crm |
| `iotdevice` | `iot.telemetry.aggregated` | report, package |
| `iotlogistics` | `iotlogistics.installation.completed` | iotdevice, package |
| `iotlogistics` | `iotlogistics.maintenance.due` | notifier, crm |
| `report` | `report.generated` | notifier |

## 0.3 Shared Infrastructure

```
pkg/                                  internal/
├─ cryptpass  ├─ jwt      ├─ kafka   ├─ middleware/   (tenant, auth)
├─ db         ├─ mqtt     ├─ influxdb├─ delivery/rest/(router)
├─ helpers    ├─ logger   ├─ llm     ├─ models/       (base)
├─ httpErrors ├─ responses├─ report  └─ repository/   (pg, redis)
├─ secureRandom├─sendEmail├─emailTemplates
├─ transaction├─ utils    ├─ vectordb
└─ elasticsearch └─ websocket  ├─ http-swagger
```

## 0.4 Multi-tenant Strategy

- ทุกตารางธุรกิจมีคอลัมน์ `tenant_id UUID NOT NULL`
- ทุก Repository query กรองด้วย `WHERE tenant_id = ?`
- Middleware `tenant.Inject()` ดึงค่าจาก header `X-Tenant-ID` (หรือ JWT claim)
- Redis key namespace: `{tenant_id}:{module}:{entity}:{id}`
- ES index: `{module}_{tenant_id}_*`

---

# PART 1 — ระบบลูกค้า (`customer`)

## 1. ภาพรวมระบบ

ระบบลูกค้าเป็น **Core Domain** ที่เก็บข้อมูลลูกค้า (Customer 360) ทั้งบุคคลและนิติบุคคล, สถานที่ติดตั้ง (Site) แยกตาม Smart Farm / Smart Building, และสัญญาบริการ (Contract) รองรับ KYC/Onboarding, เครดิต, และการระงับบริการ เป็นต้นทางของ events ที่ drive โมดูลอื่นทั้งหมด

**Bounded Context Role:** Upstream (Customer-Supplier)

## 2. โครงสร้าง Module

```
internal/modules/customer/
├── domain/
│   ├── entity/
│   │   ├── customer.go              # Aggregate Root
│   │   ├── contact.go               # Entity ย่อย
│   │   ├── site.go                  # Entity ย่อย
│   │   └── contract.go              # Aggregate Root แยก
│   ├── value_object/
│   │   ├── customer_type.go
│   │   ├── customer_status.go
│   │   ├── site_type.go
│   │   ├── tax_id.go
│   │   └── credit_limit.go
│   ├── repository/
│   │   ├── customer_repository.go
│   │   ├── site_repository.go
│   │   └── contract_repository.go
│   ├── service/
│   │   ├── kyc_service.go
│   │   └── credit_check_port.go
│   ├── event/
│   │   └── customer_events.go
│   └── errors/errors.go
├── application/
│   ├── create_customer.go
│   ├── update_customer.go
│   ├── onboard_customer.go
│   ├── suspend_customer.go
│   ├── create_site.go
│   ├── create_contract.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   │   ├── customer_repo_impl.go
│   │   ├── site_repo_impl.go
│   │   ├── contract_repo_impl.go
│   │   └── models.go
│   ├── persistence/redis/customer_cache.go
│   ├── search/elasticsearch/customer_indexer.go
│   └── messaging/kafka_producer.go
├── interfaces/
│   ├── http/
│   │   ├── customer_handler.go
│   │   ├── site_handler.go
│   │   ├── contract_handler.go
│   │   └── routes.go
│   └── middleware/tenant.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌──────────────────┐
│  CRM Context     │──(lead.converted)──┐
└──────────────────┘                     │
                                         ▼
┌──────────────────┐  (customer.created) ┌──────────────────┐
│  Customer Context│────────────────────▶│ Package Context  │
│  [CORE DOMAIN]   │                     └──────────────────┘
└────────┬─────────┘
         │ (customer.created/status.changed)
         ├───────────────▶ ERP Context
         ├───────────────▶ IoTDevice Context
         ├───────────────▶ IoTLogistics Context
         └───────────────▶ Report Context (read model)
```

**Relationship Types:**
- Customer ↔ CRM: **Customer-Supplier** (CRM ส่ง lead ที่ convert แล้ว)
- Customer ↔ Package: **Customer-Supplier** (Customer เป็น upstream)
- Customer ↔ ERP/IoT: **Published Language** ผ่าน Kafka events

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Customer | ลูกค้าที่มี identity (บุคคล/นิติบุคคล) |
| Contact | ผู้ติดต่อของลูกค้า (1 Customer : N Contacts) |
| Site | สถานที่ติดตั้ง IoT (ฟาร์ม, อาคาร, โรงงาน) |
| Contract | สัญญาบริการระหว่างลูกค้ากับแพ็กเกจ |
| KYC | Know Your Customer — การยืนยันตัวตน |
| Credit Limit | วงเงินเครดิตที่ให้ลูกค้า |
| Onboarding | กระบวนการรับลูกค้าใหม่ (KYC → Contract → Site) |

## 5. Domain Layer

### 5.1 Entities (Aggregates)

**Customer — Aggregate Root**
```go
// domain/entity/customer.go
package entity

type Customer struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Code        string                       // CUS-2026-0001
    Type        valueobject.CustomerType     // INDIVIDUAL/CORPORATE/GOVERNMENT
    Name        string
    TaxID       valueobject.TaxID
    Email       string
    Phone       string
    Address     Address                      // VO
    Status      valueobject.CustomerStatus
    CreditLimit valueobject.CreditLimit
    Contacts    []*Contact                   // Entity ย่อยใน aggregate
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewCustomer(tenantID uuid.UUID, typ valueobject.CustomerType, name string) (*Customer, error) {
    if name == "" {
        return nil, domainerrors.ErrInvalidName
    }
    if !typ.IsValid() {
        return nil, domainerrors.ErrInvalidCustomerType
    }
    now := time.Now()
    return &Customer{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Type:      typ,
        Name:      name,
        Status:    valueobject.CustomerStatusLead,
        Contacts:  []*Contact{},
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// Behavior methods
func (c *Customer) Activate() error {
    if c.Status == valueobject.CustomerStatusActive {
        return domainerrors.ErrCustomerAlreadyActive
    }
    if c.Status == valueobject.CustomerStatusChurned {
        return domainerrors.ErrCustomerChurned
    }
    c.Status = valueobject.CustomerStatusActive
    c.UpdatedAt = time.Now()
    return nil
}

func (c *Customer) Suspend(reason string) error {
    if c.Status != valueobject.CustomerStatusActive {
        return domainerrors.ErrCustomerNotActive
    }
    c.Status = valueobject.CustomerStatusSuspended
    c.UpdatedAt = time.Now()
    return nil
}

func (c *Customer) AddContact(ct *Contact) error {
    if len(c.Contacts) >= 20 {
        return domainerrors.ErrTooManyContacts
    }
    ct.CustomerID = c.ID
    c.Contacts = append(c.Contacts, ct)
    return nil
}

func (c *Customer) IsActive() bool { return c.Status == valueobject.CustomerStatusActive }
```

**Site — Entity (own aggregate via Customer)**
```go
type Site struct {
    ID         uuid.UUID
    CustomerID uuid.UUID
    Type       valueobject.SiteType   // FARM/BUILDING/FACTORY/WAREHOUSE
    Name       string
    Location   valueobject.GeoPoint
    Address    Address
    AreaSize   float64
    Metadata   map[string]any         // farm: crop; building: floors
    CreatedAt  time.Time
}
```

**Contract — Aggregate Root (แยก)**
```go
type Contract struct {
    ID         uuid.UUID
    CustomerID uuid.UUID
    PackageID  uuid.UUID
    ContractNo string
    Period     valueobject.DateRange
    Status     valueobject.ContractStatus
    SignedAt   *time.Time
    DocumentURL string
}
func (c *Contract) Sign(url string) error { /* ... */ }
func (c *Contract) Terminate(reason string) error { /* ... */ }
```

### 5.2 Value Objects

```go
// domain/value_object/customer_status.go
package valueobject

type CustomerStatus string
const (
    CustomerStatusLead      CustomerStatus = "LEAD"
    CustomerStatusActive    CustomerStatus = "ACTIVE"
    CustomerStatusSuspended CustomerStatus = "SUSPENDED"
    CustomerStatusChurned   CustomerStatus = "CHURNED"
)
func (s CustomerStatus) IsValid() bool {
    switch s {
    case CustomerStatusLead, CustomerStatusActive,
         CustomerStatusSuspended, CustomerStatusChurned:
        return true
    }
    return false
}

// domain/value_object/tax_id.go
type TaxID string
func (t TaxID) IsValid() bool {
    if len(string(t)) != 13 { return false }
    for _, r := range t { if r < '0' || r > '9' { return false } }
    return true
}

// domain/value_object/credit_limit.go
type CreditLimit struct{ Amount float64; Currency string }
func NewCreditLimit(a float64, c string) (CreditLimit, error) {
    if a < 0 { return CreditLimit{}, ErrNegativeCredit }
    return CreditLimit{a, c}, nil
}

// domain/value_object/geo_point.go
type GeoPoint struct{ Lat, Lng float64 }
func (g GeoPoint) IsValid() bool {
    return g.Lat >= -90 && g.Lat <= 90 && g.Lng >= -180 && g.Lng <= 180
}

// domain/value_object/date_range.go
type DateRange struct{ Start, End time.Time }
func (d DateRange) Contains(t time.Time) bool { /* ... */ }
```

### 5.3 Repository Interfaces

```go
// domain/repository/customer_repository.go
package repository

type CustomerRepository interface {
    Save(ctx context.Context, c *entity.Customer) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Customer, error)
    FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Customer, error)
    FindByTaxID(ctx context.Context, tenantID uuid.UUID, taxID string) (*entity.Customer, error)
    List(ctx context.Context, tenantID uuid.UUID, f ListFilter) ([]*entity.Customer, int64, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}

type SiteRepository interface {
    Save(ctx context.Context, s *entity.Site) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Site, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Site, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
}

type ContractRepository interface {
    Save(ctx context.Context, c *entity.Contract) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Contract, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Contract, error)
}
```

### 5.4 Domain Services

```go
// domain/service/kyc_service.go
package service

type KYCService struct{}
func (s *KYCService) Validate(c *entity.Customer) error {
    if c.Type == valueobject.CustomerTypeCorporate {
        if !c.TaxID.IsValid() {
            return domainerrors.ErrInvalidTaxID
        }
    }
    return nil
}

// domain/service/credit_check_port.go — Outbound Port
type CreditCheckPort interface {
    CheckCredit(ctx context.Context, customerID uuid.UUID, amount float64) (bool, error)
}
```

### 5.5 Domain Errors

```go
// domain/errors/errors.go
package domainerrors

var (
    ErrCustomerNotFound        = errors.New("customer not found")
    ErrCustomerAlreadyActive   = errors.New("customer already active")
    ErrCustomerNotActive       = errors.New("customer not active")
    ErrCustomerChurned         = errors.New("customer churned")
    ErrInvalidCustomerType     = errors.New("invalid customer type")
    ErrInvalidTaxID            = errors.New("invalid tax id")
    ErrInvalidName             = errors.New("invalid customer name")
    ErrNegativeCredit          = errors.New("credit limit cannot be negative")
    ErrTooManyContacts         = errors.New("too many contacts")
    ErrSiteNotFound            = errors.New("site not found")
    ErrContractNotFound        = errors.New("contract not found")
    ErrContractAlreadySigned   = errors.New("contract already signed")
)
```

### 5.6 Invariants

1. Customer ต้องมี `tenant_id` เสมอ
2. `TaxID` ต้อง valid ถ้าเป็น `CORPORATE` หรือ `GOVERNMENT`
3. `Status` ต้องเป็นหนึ่งใน `LEAD/ACTIVE/SUSPENDED/CHURNED`
4. `Contacts` ไม่เกิน 20 รายการ
5. Contract ที่ sign แล้วห้ามแก้ period
6. `CreditLimit.Amount ≥ 0`
7. Site ต้องมี `CustomerID` ที่อ้างถึง Customer จริง

## 6. Application Layer

### 6.1 Use Cases

```go
// application/create_customer.go
package application

type CreateCustomerUseCase struct {
    repo    repository.CustomerRepository
    kyc     *service.KYCService
    producer messaging.Producer
    audit   repository.AuditRepository
}

type CreateCustomerInput struct {
    TenantID  uuid.UUID
    Type      valueobject.CustomerType
    Name      string
    TaxID     string
    Email     string
    Phone     string
    ActorID   uuid.UUID
    IPAddress string
    UserAgent string
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, in CreateCustomerInput) (*CustomerResponse, error) {
    // 1. Validate
    if in.Name == "" { return nil, domainerrors.ErrInvalidName }

    // 2. Factory
    c, err := entity.NewCustomer(in.TenantID, in.Type, in.Name)
    if err != nil { return nil, err }

    c.TaxID = valueobject.TaxID(in.TaxID)
    c.Email = in.Email
    c.Phone = in.Phone

    // 3. Domain Service
    if err := uc.kyc.Validate(c); err != nil { return nil, err }

    // 4. Persist
    if err := uc.repo.Save(ctx, c); err != nil { return nil, err }

    // 5. Side effects
    _ = uc.producer.Publish(ctx, "customer.created", c.ID.String(), map[string]any{
        "customer_id": c.ID.String(),
        "tenant_id":   c.TenantID.String(),
        "type":        string(c.Type),
    })

    return toCustomerResponse(c), nil
}
```

**OnboardCustomer** (multi-step orchestration):
```go
// application/onboard_customer.go
type OnboardCustomerUseCase struct {
    customerRepo repository.CustomerRepository
    siteRepo     repository.SiteRepository
    contractRepo repository.ContractRepository
    kyc          *service.KYCService
    producer     messaging.Producer
}

type OnboardCustomerInput struct {
    TenantID   uuid.UUID
    CustomerID uuid.UUID
    PackageID  uuid.UUID
    Sites      []SiteInput
    ContractNo string
    StartDate  time.Time
}

func (uc *OnboardCustomerUseCase) Execute(ctx context.Context, in OnboardCustomerInput) error {
    c, err := uc.customerRepo.FindByID(ctx, in.TenantID, in.CustomerID)
    if err != nil { return err }

    if err := c.Activate(); err != nil { return err }
    if err := uc.customerRepo.Save(ctx, c); err != nil { return err }

    for _, si := range in.Sites {
        s := entity.NewSite(c.ID, si.Type, si.Name)
        if err := uc.siteRepo.Save(ctx, s); err != nil { return err }
    }

    contract := entity.NewContract(c.ID, in.PackageID, in.ContractNo)
    if err := uc.contractRepo.Save(ctx, contract); err != nil { return err }

    _ = uc.producer.Publish(ctx, "customer.onboarded", c.ID.String(), map[string]any{
        "customer_id": c.ID.String(),
        "package_id":  in.PackageID.String(),
    })
    return nil
}
```

### 6.2 DTOs

```go
// application/dto.go
type CustomerResponse struct {
    ID          string    `json:"id"`
    Code        string    `json:"code"`
    Type        string    `json:"type"`
    Name        string    `json:"name"`
    Status      string    `json:"status"`
    CreditLimit float64   `json:"credit_limit"`
    Sites       []SiteResponse `json:"sites,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type SiteResponse struct {
    ID       string  `json:"id"`
    Type     string  `json:"type"`
    Name     string  `json:"name"`
    Lat      float64 `json:"lat"`
    Lng      float64 `json:"lng"`
    AreaSize float64 `json:"area_size"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
// infrastructure/persistence/postgres/models.go
type CustomerModel struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    TenantID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    Code        string         `gorm:"size:30;not null"`
    Type        string         `gorm:"size:20;not null"`
    Name        string         `gorm:"size:255;not null"`
    TaxID       string         `gorm:"size:20;index"`
    Email       string         `gorm:"size:255"`
    Phone       string         `gorm:"size:30"`
    Address     datatypes.JSON `gorm:"type:jsonb"`
    Status      string         `gorm:"size:20;index"`
    CreditLimit float64        `gorm:"type:numeric(15,2);default:0"`
    Currency    string         `gorm:"size:3;default:'THB'"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Contacts    []ContactModel `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}
func (CustomerModel) TableName() string { return "customer_customers" }

// infrastructure/persistence/postgres/customer_repo_impl.go
type customerRepoImpl struct{ db *gorm.DB }

func (r *customerRepoImpl) Save(ctx context.Context, c *entity.Customer) error {
    m := toCustomerModel(c)
    return r.db.WithContext(ctx).
        Where("tenant_id = ?", c.TenantID).
        Save(m).Error
}

func (r *customerRepoImpl) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Customer, error) {
    var m CustomerModel
    err := r.db.WithContext(ctx).
        Preload("Contacts").
        Where("tenant_id = ? AND id = ?", tenantID, id).
        First(&m).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrCustomerNotFound
    }
    if err != nil { return nil, err }
    return toCustomerEntity(&m), nil
}
```

### 7.2 Kafka Consumers

```go
// infrastructure/messaging/consumers/crm_lead_converted_consumer.go
type LeadConvertedConsumer struct {
    createUC *application.CreateCustomerUseCase
}
func (c *LeadConvertedConsumer) Handle(ctx context.Context, payload map[string]any) error {
    // map CRM lead → CreateCustomerInput
    return c.createUC.Execute(ctx, application.CreateCustomerInput{ /* ... */ })
}
```

### 7.3 WebSocket Broadcaster
Customer module ส่ง realtime event เมื่อ customer ถูกสร้าง/ระงับ ผ่าน `pkg/websocket.Hub.SendToUser(...)` — Frontend dashboard รับทันที

### 7.4 JWT, Bcrypt, Rate Limit
- JWT: ใช้ `pkg/jwt` — claims: `sub`, `tenant_id`, `role`
- Bcrypt: ไม่ใช้ใน module นี้ (ใช้ `pkg/cryptpass` สำหรับ device token ใน `iotdevice`)
- Rate Limit: `internal/middleware/rate_limit.go` — 100 req/min ต่อ tenant

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
// interfaces/http/customer_handler.go
type CustomerHandler struct {
    createUC *application.CreateCustomerUseCase
    getUC    *application.GetCustomerUseCase
    listUC   *application.ListCustomersUseCase
    onboardUC *application.OnboardCustomerUseCase
}

func (h *CustomerHandler) Create(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    userID := c.MustGet("user_id").(uuid.UUID)

    var req struct {
        Type  string `json:"type" binding:"required"`
        Name  string `json:"name" binding:"required"`
        TaxID string `json:"tax_id"`
        Email string `json:"email"`
        Phone string `json:"phone"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }

    res, err := h.createUC.Execute(c.Request.Context(), application.CreateCustomerInput{
        TenantID: tenantID, Type: valueobject.CustomerType(req.Type),
        Name: req.Name, TaxID: req.TaxID, Email: req.Email, Phone: req.Phone,
        ActorID: userID, IPAddress: c.ClientIP(), UserAgent: c.Request.UserAgent(),
    })
    if err != nil { writeDomainError(c, err); return }
    c.JSON(201, res)
}
```

### 8.2 Routes

```go
// interfaces/http/routes.go
func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth, tenant gin.HandlerFunc) {
    g := r.Group("/customers")
    g.Use(auth, tenant)

    g.POST   ("",              h.Customer.Create)
    g.GET    ("",              h.Customer.List)
    g.GET    ("/:id",          h.Customer.Get)
    g.PUT    ("/:id",          h.Customer.Update)
    g.POST   ("/:id/onboard",  h.Customer.Onboard)
    g.POST   ("/:id/suspend",  h.Customer.Suspend)
    g.GET    ("/:id/sites",    h.Site.List)
    g.POST   ("/:id/sites",    h.Site.Create)
    g.GET    ("/:id/contracts",h.Contract.List)
}
```

### 8.3 Middleware
- `auth` — JWT
- `tenant` — ดึง `X-Tenant-ID` (หรือ JWT claim)
- `rate_limit` — ต่อ tenant

## 9. Database Migrations

```sql
-- migrations/20260101_customer_init.sql
CREATE TABLE IF NOT EXISTS customer_customers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    code          VARCHAR(30) NOT NULL,
    type          VARCHAR(20) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    tax_id        VARCHAR(20),
    email         VARCHAR(255),
    phone         VARCHAR(30),
    address       JSONB,
    status        VARCHAR(20) NOT NULL DEFAULT 'LEAD',
    credit_limit  NUMERIC(15,2) DEFAULT 0,
    currency      VARCHAR(3) DEFAULT 'THB',
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);
CREATE INDEX idx_customer_customers_tenant_status ON customer_customers(tenant_id, status);
CREATE INDEX idx_customer_customers_tenant_tax    ON customer_customers(tenant_id, tax_id);

CREATE TABLE IF NOT EXISTS customer_contacts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    name        VARCHAR(255),
    email       VARCHAR(255),
    phone       VARCHAR(30),
    position    VARCHAR(100),
    is_primary  BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customer_sites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    site_type   VARCHAR(30) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    geo_lat     NUMERIC(10,7),
    geo_lng     NUMERIC(10,7),
    address     JSONB,
    area_size   NUMERIC(12,2),
    metadata    JSONB,
    created_at  TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_customer_sites_customer ON customer_sites(customer_id);

CREATE TABLE IF NOT EXISTS customer_contracts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id   UUID NOT NULL REFERENCES customer_customers(id),
    package_id    UUID NOT NULL,
    contract_no   VARCHAR(50) NOT NULL,
    start_date    DATE NOT NULL,
    end_date      DATE,
    status        VARCHAR(20) DEFAULT 'DRAFT',
    signed_at     TIMESTAMP,
    document_url  TEXT,
    created_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(customer_id, contract_no)
);
```

## 10. System Flow

```
[Web Form] → POST /customers
    ↓
[CustomerHandler.Create]
    ↓ auth → tenant → rate limit
[CreateCustomerUseCase.Execute]
    ├─ Validate (KYCService)
    ├─ NewCustomer() factory
    ├─ CustomerRepository.Save()
    └─ Kafka.publish("customer.created")
         ↓
    [crm, package, erp consumers react]
```

## 11. Workflow Diagram (Onboarding)

```
Lead (CRM) ──convert──▶ Customer(LEAD)
                            │
                    POST /customers/:id/onboard
                            │
                  ┌─────────┴─────────┐
                  ▼                   ▼
              KYC Check          Credit Check
                  │                   │
                  └─────────┬─────────┘
                            ▼
                     Customer.ACTIVE
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
        Create Sites   Create Contract  Publish Event
              │             │             │
              └─────────────┴─────────────┘
                            ▼
                  Kafka: customer.onboarded
                            ▼
              package.subscription.start
```

## 12. การติดตั้งและใช้งาน

```bash
# Migration
make migrate-up MODULE=customer

# Env
CUSTOMER_DEFAULT_CURRENCY=THB
CUSTOMER_MAX_CONTACTS=20
CUSTOMER_KYC_ENABLED=true

# Run API
go run ./cmd/api

# Test
go test ./internal/modules/customer/...
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue Stream | ค่าบริการรายเดือน, ค่าติดตั้ง, ค่า overage |
| Cost Driver | ค่า device, ค่า bandwidth, ค่าซ่อมบำรุง |
| Key Metric | ARPU, Churn Rate, LTV, CAC |
| Customer Segment | Smart Farm (S/M/L), Smart Building (Office/Mall/Factory) |
| Pricing Model | Subscription + Usage-based + Setup fee |

## 14. ภาคผนวก

- **Env vars**: `CUSTOMER_DEFAULT_CURRENCY`, `CUSTOMER_KYC_ENABLED`
- **Indexes**: `(tenant_id, status)`, `(tenant_id, tax_id)`, `(customer_id)` on sites
- **Retention**: ลูกค้าที่ churn 5 ปี → anonymize
- **Security**: PII (tax_id, email, phone) → encrypted at rest (TDE)

## 15. Prompt สำหรับการขยายระบบ

```
เพิ่ม use case UpdateCustomerAddress ใน module customer
- Input: customer_id, address (JSON), actor_id
- Business rules: ต้องเป็น customer ACTIVE เท่านั้น
- Side effects: Kafka customer.updated, audit
- สร้าง application/update_customer_address.go
- อัปเดต customer_handler.go + routes.go
```

## 16. การตั้งชื่อตาราง Database (Prefix)

| Prefix | ตาราง |
|---|---|
| `customer_` | customers, contacts, sites, contracts |

## 17. DDD Validation Checklist

- [x] Customer เป็น Aggregate Root ที่มี invariant ครบ
- [x] Entity ไม่ expose setter ตรง — ใช้ behavior methods
- [x] Value Object immutable (TaxID, GeoPoint, DateRange)
- [x] Repository เป็น interface อยู่ที่ domain layer
- [x] Domain errors เป็น sentinel errors
- [x] Domain ไม่ import gorm/gin/sarama
- [x] Use case 1 ไฟล์ 1 responsibility
- [x] Cross-module ผ่าน Kafka เท่านั้น
- [x] Multi-tenant ผ่าน `tenant_id`
- [x] มี unit test ทุก use case

---

# PART 2 — Package & Service (`packagecatalog`)

> **หมายเหตุ**: ใช้ Go package name `packagecatalog` แต่โฟลเดอร์ `internal/modules/packagecatalog`

## 1. ภาพรวมระบบ

Catalog ของแพ็กเกจบริการ (Smart Farm / Smart Building) + Subscription Lifecycle (สมัคร/อัปเกรด/ยกเลิก) + Usage Metering + Quota Enforcement เป็น **Core Domain** ที่กำหนดขอบเขตบริการให้ลูกค้า

## 2. โครงสร้าง Module

```
internal/modules/packagecatalog/
├── domain/
│   ├── entity/
│   │   ├── package.go              # Aggregate Root
│   │   ├── service_item.go
│   │   ├── subscription.go         # Aggregate Root
│   │   ├── usage_record.go
│   │   └── entitlement.go
│   ├── value_object/
│   │   ├── billing_cycle.go
│   │   ├── package_category.go
│   │   ├── subscription_status.go
│   │   ├── service_code.go
│   │   └── quota.go
│   ├── repository/
│   │   ├── package_repository.go
│   │   ├── subscription_repository.go
│   │   └── usage_repository.go
│   ├── service/
│   │   ├── entitlement_service.go
│   │   └── proration_service.go
│   └── errors/errors.go
├── application/
│   ├── create_package.go
│   ├── publish_package.go
│   ├── subscribe.go
│   ├── upgrade_subscription.go
│   ├── cancel_subscription.go
│   ├── record_usage.go
│   ├── check_entitlement.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/quota_counter.go
│   ├── messaging/kafka/
│   │   ├── usage_producer.go
│   │   └── consumers/telemetry_usage_consumer.go
│   └── scheduler/quota_reset_job.go
├── interfaces/http/
│   ├── package_handler.go
│   ├── subscription_handler.go
│   └── routes.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌────────────────┐   ┌──────────────────────┐
│Customer Context│──▶│ Package Catalog      │
└────────────────┘   │ [CORE DOMAIN]        │
                     └────┬──────────┬──────┘
                          │          │
             (subscription.created)  (usage.recorded)◀── IoTDevice
                          │          │
                          ▼          ▼
                     ┌────────┐  ┌────────┐
                     │Payment │  │Report  │
                     └────────┘  └────────┘
```

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Package | แพ็กเกจบริการ (Smart Farm Pro, Building Basic) |
| Service Item | บริการย่อยในแพ็กเกจ (DEVICE_QUOTA, STORAGE) |
| Subscription | การสมัครใช้ของลูกค้า |
| Entitlement | สิทธิ์ที่ลูกค้าได้รับตาม subscription |
| Quota | โควต้าที่กำหนดในแพ็กเกจ |
| Usage Record | การบันทึกการใช้ (event-sourced) |
| Overage | การใช้เกินโควต้า → คิดเงินเพิ่ม |
| Proration | คำนวณค่าบริการตามสัดส่วนวัน |

## 5. Domain Layer

### 5.1 Entities

```go
// Package — Aggregate Root
type Package struct {
    ID          uuid.UUID
    TenantID    *uuid.UUID   // nil = platform-level
    Code        string
    Name        string
    Category    valueobject.PackageCategory
    BillingCycle valueobject.BillingCycle
    BasePrice   valueobject.Money
    Status      valueobject.PackageStatus
    Items       []*ServiceItem
    Metadata    map[string]any
    CreatedAt   time.Time
}

func NewPackage(code, name string, cat valueobject.PackageCategory, cycle valueobject.BillingCycle) (*Package, error) {
    if code == "" || name == "" { return nil, domainerrors.ErrInvalidPackage }
    if !cat.IsValid() || !cycle.IsValid() { return nil, domainerrors.ErrInvalidCategory }
    return &Package{
        ID: uuid.New(), Code: code, Name: name,
        Category: cat, BillingCycle: cycle,
        Status: valueobject.PackageStatusDraft,
        Items: []*ServiceItem{}, CreatedAt: time.Now(),
    }, nil
}

func (p *Package) AddItem(item *ServiceItem) error {
    if p.Status != valueobject.PackageStatusDraft {
        return domainerrors.ErrPackagePublished
    }
    for _, it := range p.Items {
        if it.ServiceCode == item.ServiceCode {
            return domainerrors.ErrDuplicateServiceCode
        }
    }
    p.Items = append(p.Items, item)
    return nil
}

func (p *Package) Publish() error {
    if p.Status == valueobject.PackageStatusPublished {
        return domainerrors.ErrAlreadyPublished
    }
    if len(p.Items) == 0 { return domainerrors.ErrEmptyPackage }
    p.Status = valueobject.PackageStatusPublished
    return nil
}

// Subscription — Aggregate Root
type Subscription struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    CustomerID      uuid.UUID
    PackageID       uuid.UUID
    ContractID      *uuid.UUID
    Status          valueobject.SubscriptionStatus
    Period          valueobject.DateRange
    NextBillingDate time.Time
    AutoRenew       bool
    CreatedAt       time.Time
}

func (s *Subscription) Activate() error {
    if s.Status != valueobject.SubscriptionStatusPending {
        return domainerrors.ErrInvalidStatusTransition
    }
    s.Status = valueobject.SubscriptionStatusActive
    return nil
}

func (s *Subscription) Suspend(reason string) error {
    if s.Status != valueobject.SubscriptionStatusActive {
        return domainerrors.ErrInvalidStatusTransition
    }
    s.Status = valueobject.SubscriptionStatusSuspended
    return nil
}

func (s *Subscription) Cancel() error {
    if s.Status == valueobject.SubscriptionStatusCancelled {
        return domainerrors.ErrAlreadyCancelled
    }
    s.Status = valueobject.SubscriptionStatusCancelled
    return nil
}

// UsageRecord — Event entity
type UsageRecord struct {
    ID              int64
    SubscriptionID  uuid.UUID
    ServiceCode     valueobject.ServiceCode
    Quantity        float64
    PeriodStart     time.Time
    IdempotencyKey  string
    RecordedAt      time.Time
}
```

### 5.2 Value Objects

```go
type BillingCycle string
const (
    BillingCycleMonthly    BillingCycle = "MONTHLY"
    BillingCycleQuarterly  BillingCycle = "QUARTERLY"
    BillingCycleYearly     BillingCycle = "YEARLY"
    BillingCycleUsageBased BillingCycle = "USAGE_BASED"
)

type PackageCategory string
const (
    PackageCategorySmartFarm     PackageCategory = "SMART_FARM"
    PackageCategorySmartBuilding PackageCategory = "SMART_BUILDING"
    PackageCategoryMixed         PackageCategory = "MIXED"
)

type SubscriptionStatus string
const (
    SubscriptionStatusPending   SubscriptionStatus = "PENDING"
    SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
    SubscriptionStatusSuspended SubscriptionStatus = "SUSPENDED"
    SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
)

type ServiceCode string
const (
    ServiceCodeDeviceQuota ServiceCode = "DEVICE_QUOTA"
    ServiceCodeStorage     ServiceCode = "STORAGE_GB"
    ServiceCodeAPICall     ServiceCode = "API_CALL"
    ServiceCodeAIInference ServiceCode = "AI_INFERENCE"
    ServiceCodeAutomation  ServiceCode = "AUTOMATION_EXEC"
)

type Quota struct {
    Limit    float64
    Unit     string
    Overage  float64   // ราคาต่อหน่วยที่เกิน
}
```

### 5.3 Repository Interfaces

```go
type PackageRepository interface {
    Save(ctx context.Context, p *entity.Package) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Package, error)
    FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Package, error)
    List(ctx context.Context, tenantID uuid.UUID, f ListFilter) ([]*entity.Package, int64, error)
}

type SubscriptionRepository interface {
    Save(ctx context.Context, s *entity.Subscription) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Subscription, error)
    FindByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Subscription, error)
    FindActiveByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) (*entity.Subscription, error)
}

type UsageRepository interface {
    Record(ctx context.Context, u *entity.UsageRecord) error
    SumBySubscriptionPeriod(ctx context.Context, subID uuid.UUID, serviceCode valueobject.ServiceCode, period time.Time) (float64, error)
}
```

### 5.4 Domain Services

```go
// EntitlementService ตรวจสิทธิ์ตาม subscription
type EntitlementService struct{ pkgRepo repository.PackageRepository }
func (s *EntitlementService) Check(sub *entity.Subscription, pkg *entity.Package, code valueobject.ServiceCode, requested float64, currentUsage float64) error {
    item := findItem(pkg, code)
    if item == nil { return domainerrors.ErrServiceNotInPackage }
    if currentUsage+requested > item.Quota.Limit {
        return domainerrors.ErrQuotaExceeded
    }
    return nil
}

// ProrationService คำนวณค่าบริการตามสัดส่วนวัน
type ProrationService struct{}
func (s *ProrationService) Calculate(oldPrice, newPrice valueobject.Money, daysRemaining, daysInCycle int) valueobject.Money {
    // ...
}
```

### 5.5 Domain Errors

```go
var (
    ErrInvalidPackage        = errors.New("invalid package")
    ErrInvalidCategory       = errors.New("invalid category")
    ErrInvalidBillingCycle   = errors.New("invalid billing cycle")
    ErrPackagePublished      = errors.New("package already published")
    ErrAlreadyPublished      = errors.New("already published")
    ErrEmptyPackage          = errors.New("package has no items")
    ErrDuplicateServiceCode  = errors.New("duplicate service code")
    ErrPackageNotFound       = errors.New("package not found")
    ErrSubscriptionNotFound  = errors.New("subscription not found")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrAlreadyCancelled      = errors.New("already cancelled")
    ErrQuotaExceeded         = errors.New("quota exceeded")
    ErrServiceNotInPackage   = errors.New("service not in package")
)
```

### 5.6 Invariants

1. `Package.Items` ต้อง unique ตาม `ServiceCode`
2. `Package` published แล้ว ห้าม add/remove items
3. `Subscription` ต้องมี `Package` ที่ published
4. Usage ต้องมี `IdempotencyKey` unique — ป้องกันนับซ้ำ
5. Quota reset ตาม `BillingCycle` ทุกต้นรอบ
6. Overage คิดราคาตาม `ServiceItem.OveragePrice`
7. Subscription cancelled → entitlement ถูก revoke ทันที

## 6. Application Layer

### 6.1 Use Cases

```go
// application/subscribe.go
type SubscribeUseCase struct {
    pkgRepo   repository.PackageRepository
    subRepo   repository.SubscriptionRepository
    custRepo  interface{ FindByID(...) }
    producer  messaging.Producer
    entitler  *service.EntitlementService
}

type SubscribeInput struct {
    TenantID   uuid.UUID
    CustomerID uuid.UUID
    PackageID  uuid.UUID
    ContractID *uuid.UUID
    StartDate  time.Time
    AutoRenew  bool
}

func (uc *SubscribeUseCase) Execute(ctx context.Context, in SubscribeInput) (*SubscriptionResponse, error) {
    pkg, err := uc.pkgRepo.FindByID(ctx, in.TenantID, in.PackageID)
    if err != nil { return nil, err }
    if pkg.Status != valueobject.PackageStatusPublished {
        return nil, domainerrors.ErrPackageNotPublished
    }

    sub := entity.NewSubscription(in.TenantID, in.CustomerID, in.PackageID, in.StartDate)
    sub.AutoRenew = in.AutoRenew

    if err := uc.subRepo.Save(ctx, sub); err != nil { return nil, err }

    _ = uc.producer.Publish(ctx, "package.subscription.created", sub.ID.String(), map[string]any{
        "subscription_id": sub.ID.String(),
        "customer_id":     in.CustomerID.String(),
        "package_id":      in.PackageID.String(),
        "start_date":      in.StartDate,
    })
    return toSubscriptionResponse(sub), nil
}
```

```go
// application/record_usage.go — hot path, idempotent
type RecordUsageUseCase struct {
    usageRepo repository.UsageRepository
    subRepo   repository.SubscriptionRepository
    redis     *redis.Client
    producer  messaging.Producer
}

type RecordUsageInput struct {
    SubscriptionID uuid.UUID
    ServiceCode    valueobject.ServiceCode
    Quantity       float64
    IdempotencyKey string
}

func (uc *RecordUsageUseCase) Execute(ctx context.Context, in RecordUsageInput) error {
    // Redis idempotency lock
    key := "usage:idem:" + in.IdempotencyKey
    ok, _ := uc.redis.SetNX(ctx, key, "1", 24*time.Hour).Result()
    if !ok { return nil } // ซ้ำ — ignore

    period := time.Now().Truncate(24 * time.Hour)
    return uc.usageRepo.Record(ctx, &entity.UsageRecord{
        SubscriptionID: in.SubscriptionID,
        ServiceCode:    in.ServiceCode,
        Quantity:       in.Quantity,
        PeriodStart:    period,
        IdempotencyKey: in.IdempotencyKey,
        RecordedAt:     time.Now(),
    })
}
```

### 6.2 DTOs

```go
type PackageResponse struct {
    ID           string                 `json:"id"`
    Code         string                 `json:"code"`
    Name         string                 `json:"name"`
    Category     string                 `json:"category"`
    BillingCycle string                 `json:"billing_cycle"`
    BasePrice    float64                `json:"base_price"`
    Currency     string                 `json:"currency"`
    Status       string                 `json:"status"`
    Items        []ServiceItemResponse  `json:"items"`
}

type SubscriptionResponse struct {
    ID              string    `json:"id"`
    CustomerID      string    `json:"customer_id"`
    PackageID       string    `json:"package_id"`
    Status          string    `json:"status"`
    StartDate       time.Time `json:"start_date"`
    NextBillingDate time.Time `json:"next_billing_date"`
    AutoRenew       bool      `json:"auto_renew"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type PackageModel struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
    TenantID     *uuid.UUID     `gorm:"type:uuid;index"`
    Code         string         `gorm:"size:50;not null"`
    Name         string         `gorm:"size:255;not null"`
    Category     string         `gorm:"size:30;index"`
    BillingCycle string         `gorm:"size:20"`
    BasePrice    float64        `gorm:"type:numeric(15,2)"`
    Currency     string         `gorm:"size:3;default:'THB'"`
    Status       string         `gorm:"size:20;index"`
    Metadata     datatypes.JSON `gorm:"type:jsonb"`
    Items        []ServiceItemModel `gorm:"foreignKey:PackageID;constraint:OnDelete:CASCADE"`
    CreatedAt    time.Time
}
func (PackageModel) TableName() string { return "package_packages" }
```

### 7.2 Kafka Consumers

```go
// consumers/telemetry_usage_consumer.go
type TelemetryUsageConsumer struct{ recordUC *application.RecordUsageUseCase }

func (c *TelemetryUsageConsumer) Handle(ctx context.Context, payload map[string]any) error {
    subID, _ := uuid.Parse(payload["subscription_id"].(string))
    qty, _ := payload["quantity"].(float64)
    return c.recordUC.Execute(ctx, application.RecordUsageInput{
        SubscriptionID: subID,
        ServiceCode:    valueobject.ServiceCode(payload["service_code"].(string)),
        Quantity:       qty,
        IdempotencyKey: payload["event_id"].(string),
    })
}
```

### 7.3 WebSocket Broadcaster
- publish `quota.exceeded` → WS broadcast ไปยัง tenant dashboard

### 7.4 JWT, Bcrypt, Rate Limit
- Rate limit เข้มงวดกับ `/packages/:id/publish` (10 req/min)

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *PackageHandler) Publish(c *gin.Context) {
    tenantID := c.MustGet("tenant_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    if err := h.publishUC.Execute(c.Request.Context(), application.PublishPackageInput{
        TenantID: tenantID, PackageID: id,
    }); err != nil { writeDomainError(c, err); return }
    c.JSON(200, gin.H{"message": "published"})
}
```

### 8.2 Routes

```go
g := r.Group("/packages"); g.Use(auth, tenant)
g.GET   ("",             h.Pkg.List)
g.POST  ("",             h.Pkg.Create)
g.GET   ("/:id",         h.Pkg.Get)
g.POST  ("/:id/publish", h.Pkg.Publish)
g.POST  ("/:id/items",   h.Pkg.AddItem)

sub := r.Group("/subscriptions"); sub.Use(auth, tenant)
sub.GET   ("",                    h.Sub.List)
sub.POST  ("",                    h.Sub.Subscribe)
sub.GET   ("/:id",                h.Sub.Get)
sub.POST  ("/:id/upgrade",        h.Sub.Upgrade)
sub.POST  ("/:id/cancel",         h.Sub.Cancel)
sub.GET   ("/:id/usage",          h.Sub.Usage)
sub.POST  ("/:id/entitlement",    h.Sub.CheckEntitlement)
```

### 8.3 Middleware
เพิ่ม `entitlement` middleware ที่เรียก `CheckEntitlementUseCase` ก่อนเข้า endpoint ที่ต้องมีโควต้า

## 9. Database Migrations

```sql
-- migrations/20260102_package_init.sql
CREATE TABLE IF NOT EXISTS package_packages (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID,
    code          VARCHAR(50) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    category      VARCHAR(30) NOT NULL,
    description   TEXT,
    billing_cycle VARCHAR(20) NOT NULL,
    base_price    NUMERIC(15,2) NOT NULL,
    currency      VARCHAR(3) DEFAULT 'THB',
    status        VARCHAR(20) DEFAULT 'DRAFT',
    metadata      JSONB,
    created_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);
CREATE INDEX idx_pkg_packages_tenant_status ON package_packages(tenant_id, status);

CREATE TABLE IF NOT EXISTS package_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    package_id    UUID NOT NULL REFERENCES package_packages(id) ON DELETE CASCADE,
    service_code  VARCHAR(50) NOT NULL,
    unit          VARCHAR(20),
    quota         NUMERIC(15,2),
    overage_price NUMERIC(15,4),
    UNIQUE(package_id, service_code)
);

CREATE TABLE IF NOT EXISTS package_subscriptions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    customer_id       UUID NOT NULL,
    package_id        UUID NOT NULL,
    contract_id       UUID,
    status            VARCHAR(20) DEFAULT 'PENDING',
    start_date        DATE NOT NULL,
    end_date          DATE,
    next_billing_date DATE,
    auto_renew        BOOLEAN DEFAULT TRUE,
    created_at        TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_pkg_subs_tenant_customer ON package_subscriptions(tenant_id, customer_id);
CREATE INDEX idx_pkg_subs_status_nextbill  ON package_subscriptions(status, next_billing_date);

CREATE TABLE IF NOT EXISTS package_usage_records (
    id              BIGSERIAL PRIMARY KEY,
    subscription_id UUID NOT NULL REFERENCES package_subscriptions(id),
    service_code    VARCHAR(50) NOT NULL,
    quantity        NUMERIC(15,4) NOT NULL,
    period_start    DATE NOT NULL,
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,
    recorded_at     TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_pkg_usage_sub_period ON package_usage_records(subscription_id, period_start);
```

## 10. System Flow

```
Customer subscribes:
  POST /subscriptions → SubscribeUseCase
    ├─ FindPackage (validate published)
    ├─ NewSubscription
    ├─ Save
    └─ Kafka: package.subscription.created
         → payment (issue invoice)
         → iotdevice (provision quota)

Usage metering:
  MQTT telemetry → iotdevice → Kafka: iot.telemetry.raw
    → TelemetryUsageConsumer → RecordUsageUseCase
      ├─ Redis idempotency
      ├─ usageRepo.Record
      └─ (if over quota) Kafka: package.quota.exceeded
```

## 11. Workflow Diagram

```
[Draft Package] ──publish──▶ [Published Package]
                                    │
                          Customer Subscribe
                                    │
                                    ▼
                          [Subscription: PENDING]
                                    │
                             Activate (via payment)
                                    │
                                    ▼
                          [Subscription: ACTIVE]
                                    │
              ┌─────────────────────┼────────────────────┐
              ▼                     ▼                    ▼
        Quota Counter         Usage Records        Entitlement
        (Redis)               (PG)                 Check
              │
        Reset ทุก cycle (scheduler)
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=packagecatalog
CURRENCY_DEFAULT=THB
QUOTA_RESET_CRON="0 0 1 * *"
go run ./cmd/api
```

## 13. Business Model

| Package | ราคา/เดือน | Devices | Storage | AI Calls |
|---|---|---|---|---|
| Smart Farm Basic | ฿1,500 | 10 | 5 GB | 100 |
| Smart Farm Pro | ฿4,500 | 50 | 50 GB | 1,000 |
| Smart Building Basic | ฿3,500 | 30 | 20 GB | 300 |
| Smart Building Pro | ฿9,900 | 200 | 200 GB | 5,000 |

## 14. ภาคผนวก
- Overage pricing ตาม `ServiceItem.OveragePrice`
- Upgrade: prorate ค่าบริการตามวันคงเหลือ
- Cancel: ไม่คืนเงิน, สิทธิ์หมดทันที

## 15. Prompt ขยายระบบ

```
เพิ่ม use case BulkSubscribe ใน module packagecatalog
- รับ list ของ (customer_id, package_id)
- ตรวจ entitlement ก่อน
- Atomic transaction ผ่าน pkg/transaction
- Publish Kafka ต่อ subscription
```

## 16. Prefix ตาราง

| Prefix | ตาราง |
|---|---|
| `package_` | packages, items, subscriptions, usage_records |

## 17. DDD Validation Checklist

- [x] Package/Subscription เป็น Aggregate Root แยกกัน
- [x] Usage เป็น event-sourced (idempotent)
- [x] Quota เก็บใน Redis + DB (dual write)
- [x] Proration คำนวณใน domain service
- [x] Entitlement ตรวจผ่าน domain service
- [x] Snapshot ราคา ณ เวลา subscribe (ห้าม retrive จาก package ปัจจุบัน)
- [x] Kafka events สำหรับ integration

---

# PART 3 — ERP IoT Solution (`erp`)

## 1. ภาพรวมระบบ

ระบบ ERP สำหรับธุรกิจ IoT: จัดการ Product/Device Catalog, Inventory (Warehouse + Stock Movements), Purchase Orders, Sales Orders, Invoicing, Payment, Vendor Management รองรับ multi-warehouse และ multi-currency

## 2. โครงสร้าง Module

```
internal/modules/erp/
├── domain/
│   ├── entity/
│   │   ├── product.go
│   │   ├── warehouse.go
│   │   ├── inventory.go            # Aggregate Root
│   │   ├── purchase_order.go       # Aggregate Root
│   │   ├── sales_order.go          # Aggregate Root
│   │   ├── invoice.go
│   │   ├── payment.go
│   │   ├── vendor.go
│   │   └── stock_movement.go
│   ├── value_object/
│   │   ├── money.go
│   │   ├── sku.go
│   │   ├── uom.go
│   │   ├── document_status.go
│   │   └── movement_type.go
│   ├── repository/
│   │   ├── product_repository.go
│   │   ├── inventory_repository.go
│   │   ├── po_repository.go
│   │   ├── so_repository.go
│   │   └── invoice_repository.go
│   ├── service/
│   │   ├── inventory_service.go
│   │   ├── pricing_service.go
│   │   └── tax_service.go
│   └── errors/errors.go
├── application/
│   ├── create_product.go
│   ├── stock_adjust.go
│   ├── create_po.go
│   ├── approve_po.go
│   ├── receive_goods.go
│   ├── create_so.go
│   ├── issue_invoice.go
│   ├── post_payment.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── services/pdf/invoice_pdf.go
│   ├── messaging/kafka/accounting_producer.go
│   └── scheduler/low_stock_alert_job.go
├── interfaces/http/
│   ├── product_handler.go
│   ├── inventory_handler.go
│   ├── po_handler.go
│   ├── so_handler.go
│   ├── invoice_handler.go
│   └── routes.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌──────────────────┐  (customer.created)  ┌──────────────┐
│Customer Context  │─────────────────────▶│ ERP Context  │
└──────────────────┘                      │              │
                                          │  Inventory   │
┌──────────────────┐  (installation.done) │  AR / AP     │
│IoTLogistics      │─────────────────────▶│  Product     │
└──────────────────┘                      └──────┬───────┘
                                                 │
                                                 │ (invoice.issued)
                                                 ▼
                                          ┌──────────────┐
                                          │Payment/Report│
                                          └──────────────┘
```

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Product | สินค้า/อุปกรณ์ (SKU) |
| SKU | Stock Keeping Unit |
| Inventory | สต็อกคงเหลือ ณ คลังหนึ่ง |
| Stock Movement | การเคลื่อนไหวของสต็อก (IN/OUT/ADJUST) |
| PO | Purchase Order |
| GRN | Goods Receipt Note |
| SO | Sales Order |
| Invoice | ใบแจ้งหนี้ |
| Vendor | ผู้ขาย/ซัพพลายเออร์ |
| UOM | Unit of Measure |

## 5. Domain Layer

### 5.1 Entities

```go
// Product — Aggregate Root
type Product struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    SKU        valueobject.SKU
    Name       string
    Category   string
    UOM        valueobject.UOM
    CostPrice  valueobject.Money
    SellPrice  valueobject.Money
    IsIoTDevice bool
    DeviceModelID *uuid.UUID   // link iotdevice.DeviceModel
    Status     valueobject.ProductStatus
    Metadata   map[string]any
    CreatedAt  time.Time
}

// Inventory — Aggregate Root (per product+warehouse)
type Inventory struct {
    ID            uuid.UUID
    TenantID      uuid.UUID
    ProductID     uuid.UUID
    WarehouseID   uuid.UUID
    QtyOnHand     float64
    QtyReserved   float64
    UpdatedAt     time.Time
}

func (i *Inventory) Available() float64 { return i.QtyOnHand - i.QtyReserved }

func (i *Inventory) Reserve(qty float64) error {
    if qty <= 0 { return domainerrors.ErrInvalidQuantity }
    if i.Available() < qty { return domainerrors.ErrInsufficientStock }
    i.QtyReserved += qty
    i.UpdatedAt = time.Now()
    return nil
}

func (i *Inventory) Release(qty float64) error {
    if i.QtyReserved < qty { return domainerrors.ErrInvalidQuantity }
    i.QtyReserved -= qty
    return nil
}

func (i *Inventory) Adjust(delta float64, reason string) error {
    if i.QtyOnHand+delta < 0 { return domainerrors.ErrNegativeStock }
    i.QtyOnHand += delta
    return nil
}

// PurchaseOrder — Aggregate Root
type PurchaseOrder struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    PONo        string
    VendorID    uuid.UUID
    Status      valueobject.DocumentStatus
    Lines       []*POLine
    TotalAmount valueobject.Money
    ExpectedAt  time.Time
    CreatedBy   uuid.UUID
    CreatedAt   time.Time
    ApprovedAt  *time.Time
}

func (po *PurchaseOrder) Approve(actor uuid.UUID) error {
    if po.Status != valueobject.DocStatusDraft {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(po.Lines) == 0 { return domainerrors.ErrEmptyDocument }
    po.Status = valueobject.DocStatusApproved
    now := time.Now()
    po.ApprovedAt = &now
    return nil
}

func (po *PurchaseOrder) Receive(warehouseID uuid.UUID) ([]*StockMovement, error) {
    if po.Status != valueobject.DocStatusApproved {
        return nil, domainerrors.ErrInvalidStatusTransition
    }
    movements := make([]*StockMovement, 0, len(po.Lines))
    for _, l := range po.Lines {
        movements = append(movements, &StockMovement{
            ProductID:    l.ProductID,
            WarehouseID:  warehouseID,
            Type:         valueobject.MovementTypeIn,
            Quantity:     l.Quantity,
            RefType:      "PO",
            RefID:        po.ID,
        })
    }
    po.Status = valueobject.DocStatusReceived
    return movements, nil
}

// Invoice — Aggregate Root
type Invoice struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    InvoiceNo   string
    CustomerID  uuid.UUID
    OrderID     *uuid.UUID
    Lines       []*InvoiceLine
    TotalAmount valueobject.Money
    TaxAmount   valueobject.Money
    Status      valueobject.InvoiceStatus
    DueDate     time.Time
    IssuedAt    time.Time
    PaidAt      *time.Time
}
```

### 5.2 Value Objects

```go
type Money struct {
    Amount   float64
    Currency string
}
func (m Money) Add(o Money) (Money, error) {
    if m.Currency != o.Currency { return Money{}, domainerrors.ErrCurrencyMismatch }
    return Money{m.Amount + o.Amount, m.Currency}, nil
}
func (m Money) IsZero() bool { return m.Amount == 0 }

type SKU string
func (s SKU) IsValid() bool { return len(s) >= 3 && len(s) <= 50 }

type UOM string
const (
    UOMPiece  UOM = "PCS"
    UOMBox    UOM = "BOX"
    UOMKg     UOM = "KG"
    UOMMeter  UOM = "M"
)

type DocumentStatus string
const (
    DocStatusDraft    DocumentStatus = "DRAFT"
    DocStatusApproved DocumentStatus = "APPROVED"
    DocStatusPosted   DocumentStatus = "POSTED"
    DocStatusReceived DocumentStatus = "RECEIVED"
    DocStatusVoid     DocumentStatus = "VOID"
)

type MovementType string
const (
    MovementTypeIn       MovementType = "IN"
    MovementTypeOut      MovementType = "OUT"
    MovementTypeAdjust   MovementType = "ADJUST"
    MovementTypeTransfer MovementType = "TRANSFER"
)
```

### 5.3 Repository Interfaces

```go
type ProductRepository interface {
    Save(ctx context.Context, p *entity.Product) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Product, error)
    FindBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*entity.Product, error)
    List(ctx context.Context, tenantID uuid.UUID, f ListFilter) ([]*entity.Product, int64, error)
}

type InventoryRepository interface {
    Save(ctx context.Context, i *entity.Inventory) error
    FindByProductWarehouse(ctx context.Context, tenantID, productID, warehouseID uuid.UUID) (*entity.Inventory, error)
    ApplyMovement(ctx context.Context, m *entity.StockMovement) error  // transaction
    ListLowStock(ctx context.Context, tenantID uuid.UUID, threshold float64) ([]*entity.Inventory, error)
}

type PurchaseOrderRepository interface {
    Save(ctx context.Context, po *entity.PurchaseOrder) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.PurchaseOrder, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.PurchaseOrder, error)
}

type InvoiceRepository interface {
    Save(ctx context.Context, inv *entity.Invoice) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Invoice, error)
    FindOverdue(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.Invoice, error)
}
```

### 5.4 Domain Services

```go
// InventoryService — logic ข้ามหลาย inventory
type InventoryService struct{ repo repository.InventoryRepository }
func (s *InventoryService) TransferStock(ctx context.Context, tenantID, productID, from, to uuid.UUID, qty float64) error {
    // load both inventories → out from 'from', in to 'to'
}

// PricingService — คำนวณราคาตาม tier + ส่วนลด + ภาษี
type PricingService struct{}
func (s *PricingService) Calculate(subtotal valueobject.Money, discountPct float64, taxPct float64) (total, tax valueobject.Money, err error)
```

### 5.5 Domain Errors

```go
var (
    ErrProductNotFound       = errors.New("product not found")
    ErrSKUExists             = errors.New("sku already exists")
    ErrInvalidSKU            = errors.New("invalid sku")
    ErrInventoryNotFound     = errors.New("inventory not found")
    ErrInsufficientStock     = errors.New("insufficient stock")
    ErrNegativeStock         = errors.New("negative stock not allowed")
    ErrInvalidQuantity       = errors.New("invalid quantity")
    ErrPONotFound            = errors.New("purchase order not found")
    ErrSONotFound            = errors.New("sales order not found")
    ErrInvoiceNotFound       = errors.New("invoice not found")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrEmptyDocument         = errors.New("document has no lines")
    ErrCurrencyMismatch      = errors.New("currency mismatch")
)
```

### 5.6 Invariants

1. SKU unique ต่อ tenant
2. `Inventory.QtyOnHand ≥ 0` เสมอ
3. `Inventory.QtyReserved ≤ QtyOnHand`
4. PO ที่ `APPROVED` แล้ว แก้ lines ไม่ได้
5. Invoice ต้องมี lines ≥ 1
6. Payment ต้อง ≤ ยอดค้างชำระ
7. StockMovement ต้อง update inventory ใน transaction เดียวกัน
8. Currency ใน line ต้องตรงกับ header

## 6. Application Layer

### 6.1 Use Cases

```go
// application/stock_adjust.go
type StockAdjustUseCase struct {
    invRepo repository.InventoryRepository
    tx      *transaction.Manager
    producer messaging.Producer
}

type StockAdjustInput struct {
    TenantID    uuid.UUID
    ProductID   uuid.UUID
    WarehouseID uuid.UUID
    Delta       float64
    Reason      string
    ActorID     uuid.UUID
}

func (uc *StockAdjustUseCase) Execute(ctx context.Context, in StockAdjustInput) error {
    return uc.tx.Do(ctx, func(tx *gorm.DB) error {
        inv, err := uc.invRepo.FindByProductWarehouse(ctx, in.TenantID, in.ProductID, in.WarehouseID)
        if err != nil { return err }

        if err := inv.Adjust(in.Delta, in.Reason); err != nil { return err }
        if err := uc.invRepo.Save(ctx, inv); err != nil { return err }

        m := &entity.StockMovement{
            ProductID: in.ProductID, WarehouseID: in.WarehouseID,
            Type: valueobject.MovementTypeAdjust, Quantity: in.Delta,
            RefType: "ADJUST", Reason: in.Reason,
        }
        if err := uc.invRepo.ApplyMovement(ctx, m); err != nil { return err }

        if inv.Available() < 10 {
            _ = uc.producer.Publish(ctx, "erp.inventory.low_stock", in.ProductID.String(), map[string]any{
                "product_id": in.ProductID.String(), "available": inv.Available(),
            })
        }
        return nil
    })
}
```

```go
// application/receive_goods.go
type ReceiveGoodsUseCase struct {
    poRepo  repository.PurchaseOrderRepository
    invRepo repository.InventoryRepository
    tx      *transaction.Manager
    producer messaging.Producer
}

func (uc *ReceiveGoodsUseCase) Execute(ctx context.Context, in ReceiveGoodsInput) error {
    return uc.tx.Do(ctx, func(tx *gorm.DB) error {
        po, err := uc.poRepo.FindByID(ctx, in.TenantID, in.POID)
        if err != nil { return err }

        movements, err := po.Receive(in.WarehouseID)
        if err != nil { return err }

        for _, m := range movements {
            if err := uc.invRepo.ApplyMovement(ctx, m); err != nil { return err }
        }
        if err := uc.poRepo.Save(ctx, po); err != nil { return err }

        _ = uc.producer.Publish(ctx, "erp.po.received", po.ID.String(), map[string]any{
            "po_id": po.ID.String(), "received_at": time.Now(),
        })
        return nil
    })
}
```

### 6.2 DTOs

```go
type ProductResponse struct {
    ID          string  `json:"id"`
    SKU         string  `json:"sku"`
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    CostPrice   float64 `json:"cost_price"`
    SellPrice   float64 `json:"sell_price"`
    IsIoTDevice bool    `json:"is_iot_device"`
    Available   float64 `json:"available,omitempty"`
}

type InvoiceResponse struct {
    ID          string    `json:"id"`
    InvoiceNo   string    `json:"invoice_no"`
    CustomerID  string    `json:"customer_id"`
    TotalAmount float64   `json:"total_amount"`
    Status      string    `json:"status"`
    DueDate     time.Time `json:"due_date"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type InventoryModel struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    TenantID    uuid.UUID `gorm:"type:uuid;not null;index:idx_inv_unique,unique"`
    ProductID   uuid.UUID `gorm:"type:uuid;not null;index:idx_inv_unique,unique"`
    WarehouseID uuid.UUID `gorm:"type:uuid;not null;index:idx_inv_unique,unique"`
    QtyOnHand   float64   `gorm:"type:numeric(15,2);default:0"`
    QtyReserved float64   `gorm:"type:numeric(15,2);default:0"`
    UpdatedAt   time.Time
}
func (InventoryModel) TableName() string { return "erp_inventory" }
```

### 7.2 Kafka Consumers

```go
// consumers/customer_created_consumer.go — provision default AR account
// consumers/installation_completed_consumer.go — สร้าง invoice ค่าติดตั้ง
```

### 7.3 WebSocket Broadcaster
- publish low stock → WS ทันที
- publish invoice overdue → WS

### 7.4 JWT, Bcrypt, Rate Limit
- PO approve ต้องมี role `manager` (ผ่าน JWT claim)

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *POHandler) Approve(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    if err := h.approveUC.Execute(c.Request.Context(), application.ApprovePOInput{
        TenantID: tid, POID: id, ActorID: uid,
    }); err != nil { writeDomainError(c, err); return }
    c.JSON(200, gin.H{"message": "approved"})
}
```

### 8.2 Routes

```go
p := r.Group("/products"); p.Use(auth, tenant)
p.GET("", h.Prod.List); p.POST("", h.Prod.Create); p.GET("/:id", h.Prod.Get)

i := r.Group("/inventory"); i.Use(auth, tenant)
i.GET("", h.Inv.List); i.POST("/adjust", h.Inv.Adjust); i.GET("/low-stock", h.Inv.LowStock)

po := r.Group("/purchase-orders"); po.Use(auth, tenant)
po.POST("", h.PO.Create); po.POST("/:id/approve", h.PO.Approve); po.POST("/:id/receive", h.PO.Receive)

so := r.Group("/sales-orders"); so.Use(auth, tenant)
so.POST("", h.SO.Create); so.POST("/:id/invoice", h.SO.Invoice)

inv := r.Group("/invoices"); inv.Use(auth, tenant)
inv.GET("", h.Inv2.List); inv.POST("/:id/pay", h.Inv2.Pay)
```

### 8.3 Middleware
- `role.Require("manager")` สำหรับ approve

## 9. Database Migrations

```sql
-- migrations/20260103_erp_init.sql
CREATE TABLE IF NOT EXISTS erp_products (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL,
    sku             VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    category        VARCHAR(50),
    uom             VARCHAR(20) DEFAULT 'PCS',
    cost_price      NUMERIC(15,2),
    sell_price      NUMERIC(15,2),
    is_iot_device   BOOLEAN DEFAULT FALSE,
    device_model_id UUID,
    status          VARCHAR(20) DEFAULT 'ACTIVE',
    metadata        JSONB,
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, sku)
);

CREATE TABLE IF NOT EXISTS erp_warehouses (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    code       VARCHAR(30) NOT NULL,
    name       VARCHAR(255),
    address    JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE IF NOT EXISTS erp_inventory (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    product_id   UUID NOT NULL REFERENCES erp_products(id),
    warehouse_id UUID NOT NULL REFERENCES erp_warehouses(id),
    qty_on_hand  NUMERIC(15,2) DEFAULT 0,
    qty_reserved NUMERIC(15,2) DEFAULT 0,
    updated_at   TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, product_id, warehouse_id)
);

CREATE TABLE IF NOT EXISTS erp_stock_movements (
    id           BIGSERIAL PRIMARY KEY,
    tenant_id    UUID NOT NULL,
    product_id   UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    movement_type VARCHAR(20) NOT NULL,
    quantity     NUMERIC(15,2) NOT NULL,
    ref_type     VARCHAR(30),
    ref_id       UUID,
    reason       TEXT,
    created_by   UUID,
    created_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_stock_mv_product ON erp_stock_movements(product_id, created_at DESC);

CREATE TABLE IF NOT EXISTS erp_purchase_orders (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    po_no         VARCHAR(30) NOT NULL,
    vendor_id     UUID,
    status        VARCHAR(20) DEFAULT 'DRAFT',
    total_amount  NUMERIC(15,2),
    currency      VARCHAR(3) DEFAULT 'THB',
    expected_at   DATE,
    created_by    UUID,
    approved_at   TIMESTAMP,
    created_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, po_no)
);

CREATE TABLE IF NOT EXISTS erp_po_lines (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_id       UUID NOT NULL REFERENCES erp_purchase_orders(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL,
    quantity    NUMERIC(15,2) NOT NULL,
    unit_price  NUMERIC(15,2),
    total       NUMERIC(15,2)
);

CREATE TABLE IF NOT EXISTS erp_invoices (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    invoice_no    VARCHAR(30) NOT NULL,
    customer_id   UUID NOT NULL,
    order_id      UUID,
    subtotal      NUMERIC(15,2),
    tax_amount    NUMERIC(15,2),
    total_amount  NUMERIC(15,2),
    currency      VARCHAR(3) DEFAULT 'THB',
    status        VARCHAR(20) DEFAULT 'DRAFT',
    due_date      DATE,
    issued_at     TIMESTAMP,
    paid_at       TIMESTAMP,
    UNIQUE(tenant_id, invoice_no)
);
CREATE INDEX idx_erp_invoices_status_due ON erp_invoices(status, due_date);
```

## 10. System Flow

```
Procurement:
  Create PO → Approve → Send to vendor → Receive Goods (GRN)
    → StockMovements IN → Inventory++ 
    → Kafka: erp.po.received

Sales:
  Sales Order (จาก CRM/Portal)
    → Reserve Inventory
    → Issue Invoice
    → Kafka: erp.invoice.issued
    → Payment.posted → Inventory OUT (delivery)
```

## 11. Workflow Diagram

```
[SO Created]──▶ [Reserve Stock]──▶ [Invoice Issued]
                                        │
                                        ▼
                                  [Customer Pays]
                                        │
                                        ▼
                             [Payment Posted]
                                        │
                        ┌───────────────┼───────────────┐
                        ▼               ▼               ▼
                   Stock OUT        AR Updated      Report
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=erp
LOW_STOCK_THRESHOLD=10
INVOICE_DUE_DAYS=30
go run ./cmd/api
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue | ขายอุปกรณ์ (margin 15-30%) + ค่าติดตั้ง |
| Cost | COGS, logistics, warranty |
| KPI | Inventory turnover, DSO, Gross margin |

## 14. ภาคผนวก
- Multi-currency: เก็บ currency ที่ line + header
- Tax: VAT 7% (default), tax-exempt สำหรับบางสินค้า
- Stock opname: ใช้ `ADJUST` movement

## 15. Prompt ขยายระบบ

```
เพิ่ม use case CreateSalesOrder ใน module erp
- Input: customer_id, lines[], warehouse_id
- Business rules: ต้อง reserve stock ให้ครบก่อน save
- Side effects: Kafka erp.so.created, low_stock check
```

## 16. Prefix ตาราง

| Prefix | ตาราง |
|---|---|
| `erp_` | products, warehouses, inventory, stock_movements, purchase_orders, po_lines, sales_orders, so_lines, invoices, invoice_lines, payments, vendors |

## 17. DDD Validation Checklist

- [x] Inventory/PO/SO/Invoice เป็น Aggregate Root แยกกัน
- [x] Money เป็น Value Object พร้อม currency safety
- [x] Stock movement update inventory ใน transaction เดียว
- [x] SKU unique ต่อ tenant
- [x] Approve PO ต้องมี role
- [x] Kafka events สำหรับ payment/report
- [x] Low stock → event + WS

---

# PART 4 — CRM IoT Solution (`crm`)

## 1. ภาพรวมระบบ

ระบบ CRM สำหรับธุรกิจ IoT: Lead Management (จาก website/event/ads), Sales Pipeline (Opportunity), Activity Tracking (call/email/meeting), Support Ticket (จาก device alert), Campaign, Lead Scoring ด้วย AI (LLM)

## 2. โครงสร้าง Module

```
internal/modules/crm/
├── domain/
│   ├── entity/
│   │   ├── lead.go                 # Aggregate Root
│   │   ├── opportunity.go          # Aggregate Root
│   │   ├── pipeline_stage.go
│   │   ├── activity.go
│   │   ├── ticket.go               # Aggregate Root
│   │   └── campaign.go
│   ├── value_object/
│   │   ├── lead_source.go
│   │   ├── lead_status.go
│   │   ├── opportunity_stage.go
│   │   ├── ticket_priority.go
│   │   └── ticket_status.go
│   ├── repository/
│   │   ├── lead_repository.go
│   │   ├── opportunity_repository.go
│   │   ├── activity_repository.go
│   │   └── ticket_repository.go
│   ├── service/
│   │   ├── lead_scoring_service.go
│   │   ├── sla_service.go
│   │   └── assignment_service.go
│   └── errors/errors.go
├── application/
│   ├── create_lead.go
│   ├── convert_lead.go
│   ├── move_opportunity.go
│   ├── close_opportunity.go
│   ├── log_activity.go
│   ├── create_ticket.go
│   ├── assign_ticket.go
│   ├── resolve_ticket.go
│   ├── classify_ticket.go          # AI
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/redis/lead_assignment.go
│   ├── services/llm/ticket_classifier.go
│   ├── messaging/kafka/
│   └── scheduler/sla_check_job.go
├── interfaces/http/
│   ├── lead_handler.go
│   ├── opportunity_handler.go
│   ├── ticket_handler.go
│   └── routes.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌───────────┐  (lead.created)   ┌──────────────┐
│ Web/Ads   │──────────────────▶│ CRM Context  │
└───────────┘                   │ [SUPPORTING] │
                                └───┬──────┬───┘
                                    │      │
                        (lead.converted)  (ticket.created)
                                    │      │
                                    ▼      ▼
                          ┌──────────┐  ┌───────────┐
                          │ Customer │  │ Notifier  │
                          └──────────┘  └───────────┘
                                ▲
                                │ (iot.alert.triggered)
                          ┌─────┴────┐
                          │ IoTDevice│
                          └──────────┘
```

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Lead | ผู้สนใจ (ยังไม่เป็นลูกค้า) |
| Opportunity | โอกาสขาย (มีมูลค่า, มี stage) |
| Pipeline | ลำดับ stage ของการขาย |
| Activity | การติดต่อ (call/email/meeting) |
| Ticket | ใบแจ้งปัญหา/คำขอรับบริการ |
| SLA | Service Level Agreement |
| Lead Score | คะแนนความสนใจ (0-100) |
| Campaign | แคมเปญการตลาด |

## 5. Domain Layer

### 5.1 Entities

```go
// Lead — Aggregate Root
type Lead struct {
    ID                  uuid.UUID
    TenantID            uuid.UUID
    Source              valueobject.LeadSource
    Name                string
    Email               string
    Phone               string
    Company             string
    InterestCategory    valueobject.PackageCategory   // SMART_FARM | SMART_BUILDING
    Score               int
    Status              valueobject.LeadStatus
    AssignedTo          *uuid.UUID
    ConvertedCustomerID *uuid.UUID
    Metadata            map[string]any
    CreatedAt           time.Time
}

func NewLead(tenantID uuid.UUID, source valueobject.LeadSource, name string) (*Lead, error) {
    if name == "" { return nil, domainerrors.ErrInvalidName }
    if !source.IsValid() { return nil, domainerrors.ErrInvalidSource }
    return &Lead{
        ID: uuid.New(), TenantID: tenantID, Source: source, Name: name,
        Status: valueobject.LeadStatusNew, Score: 0, CreatedAt: time.Now(),
    }, nil
}

func (l *Lead) UpdateScore(score int) error {
    if score < 0 || score > 100 { return domainerrors.ErrInvalidScore }
    l.Score = score
    return nil
}

func (l *Lead) Qualify() error {
    if l.Status != valueobject.LeadStatusNew && l.Status != valueobject.LeadStatusContacted {
        return domainerrors.ErrInvalidStatusTransition
    }
    l.Status = valueobject.LeadStatusQualified
    return nil
}

func (l *Lead) Convert(customerID uuid.UUID) error {
    if l.Status != valueobject.LeadStatusQualified {
        return domainerrors.ErrLeadNotQualified
    }
    l.Status = valueobject.LeadStatusConverted
    l.ConvertedCustomerID = &customerID
    return nil
}

// Opportunity — Aggregate Root
type Opportunity struct {
    ID           uuid.UUID
    TenantID     uuid.UUID
    CustomerID   *uuid.UUID
    LeadID       *uuid.UUID
    Name         string
    Stage        valueobject.OpportunityStage
    Amount       valueobject.Money
    Probability  int
    ExpectedCloseDate time.Time
    OwnerID      uuid.UUID
    Activities   []*Activity
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (o *Opportunity) MoveToStage(stage valueobject.OpportunityStage) error {
    if !isValidTransition(o.Stage, stage) {
        return domainerrors.ErrInvalidStageTransition
    }
    o.Stage = stage
    o.UpdatedAt = time.Now()
    return nil
}

func (o *Opportunity) Win() error {
    if o.Stage == valueobject.OpportunityStageWon { return domainerrors.ErrAlreadyWon }
    o.Stage = valueobject.OpportunityStageWon
    o.Probability = 100
    return nil
}

func (o *Opportunity) Lose(reason string) error {
    o.Stage = valueobject.OpportunityStageLost
    o.Probability = 0
    return nil
}

// Ticket — Aggregate Root
type Ticket struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    CustomerID  uuid.UUID
    DeviceID    *uuid.UUID
    Subject     string
    Description string
    Category    string           // AI-classified
    Priority    valueobject.TicketPriority
    Status      valueobject.TicketStatus
    SLADueAt    time.Time
    AssignedTo  *uuid.UUID
    ResolvedAt  *time.Time
    CreatedAt   time.Time
}

func (t *Ticket) Assign(actor uuid.UUID) error {
    if t.Status == valueobject.TicketStatusClosed {
        return domainerrors.ErrTicketClosed
    }
    t.AssignedTo = &actor
    t.Status = valueobject.TicketStatusInProgress
    return nil
}

func (t *Ticket) Resolve(note string) error {
    now := time.Now()
    t.Status = valueobject.TicketStatusResolved
    t.ResolvedAt = &now
    return nil
}

func (t *Ticket) IsSLABreached() bool {
    return t.SLADueAt.Before(time.Now()) && t.Status != valueobject.TicketStatusResolved
}
```

### 5.2 Value Objects

```go
type LeadSource string
const (
    LeadSourceWeb      LeadSource = "WEB"
    LeadSourceReferral LeadSource = "REFERRAL"
    LeadSourceEvent    LeadSource = "EVENT"
    LeadSourceAds      LeadSource = "ADS"
    LeadSourceWalkin   LeadSource = "WALKIN"
)

type LeadStatus string
const (
    LeadStatusNew       LeadStatus = "NEW"
    LeadStatusContacted LeadStatus = "CONTACTED"
    LeadStatusQualified LeadStatus = "QUALIFIED"
    LeadStatusConverted LeadStatus = "CONVERTED"
    LeadStatusLost      LeadStatus = "LOST"
)

type OpportunityStage string
const (
    OpportunityStageNew       OpportunityStage = "NEW"
    OpportunityStageQualified OpportunityStage = "QUALIFIED"
    OpportunityStageProposal  OpportunityStage = "PROPOSAL"
    OpportunityStageNegotiation OpportunityStage = "NEGOTIATION"
    OpportunityStageWon       OpportunityStage = "WON"
    OpportunityStageLost      OpportunityStage = "LOST"
)

type TicketPriority string
const (
    TicketPriorityLow      TicketPriority = "LOW"
    TicketPriorityNormal   TicketPriority = "NORMAL"
    TicketPriorityHigh     TicketPriority = "HIGH"
    TicketPriorityCritical TicketPriority = "CRITICAL"
)

type TicketStatus string
const (
    TicketStatusOpen       TicketStatus = "OPEN"
    TicketStatusInProgress TicketStatus = "IN_PROGRESS"
    TicketStatusResolved   TicketStatus = "RESOLVED"
    TicketStatusClosed     TicketStatus = "CLOSED"
)
```

### 5.3 Repository Interfaces

```go
type LeadRepository interface {
    Save(ctx context.Context, l *entity.Lead) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Lead, error)
    FindByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.Lead, error)
    List(ctx context.Context, tenantID uuid.UUID, f ListFilter) ([]*entity.Lead, int64, error)
}

type OpportunityRepository interface {
    Save(ctx context.Context, o *entity.Opportunity) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Opportunity, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Opportunity, error)
    ListByStage(ctx context.Context, tenantID uuid.UUID, stage valueobject.OpportunityStage) ([]*entity.Opportunity, error)
}

type TicketRepository interface {
    Save(ctx context.Context, t *entity.Ticket) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Ticket, error)
    FindSLABreached(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.Ticket, error)
}
```

### 5.4 Domain Services

```go
// LeadScoringService
type LeadScoringService struct{}
func (s *LeadScoringService) Score(l *entity.Lead) int {
    score := 0
    switch l.Source {
    case valueobject.LeadSourceReferral: score += 30
    case valueobject.LeadSourceWeb:      score += 20
    case valueobject.LeadSourceEvent:    score += 15
    }
    if l.Company != "" { score += 20 }
    if l.Email != "" && l.Phone != "" { score += 20 }
    if score > 100 { score = 100 }
    return score
}

// SLAService — คำนวณ SLADueAt ตาม priority
type SLAService struct{}
func (s *SLAService) Calculate(priority valueobject.TicketPriority, createdAt time.Time) time.Time {
    switch priority {
    case valueobject.TicketPriorityCritical: return createdAt.Add(1 * time.Hour)
    case valueobject.TicketPriorityHigh:     return createdAt.Add(4 * time.Hour)
    case valueobject.TicketPriorityNormal:   return createdAt.Add(24 * time.Hour)
    default:                                  return createdAt.Add(72 * time.Hour)
    }
}
```

### 5.5 Domain Errors

```go
var (
    ErrLeadNotFound           = errors.New("lead not found")
    ErrLeadNotQualified       = errors.New("lead not qualified")
    ErrInvalidSource          = errors.New("invalid lead source")
    ErrInvalidScore           = errors.New("invalid score")
    ErrOpportunityNotFound    = errors.New("opportunity not found")
    ErrInvalidStageTransition = errors.New("invalid stage transition")
    ErrAlreadyWon             = errors.New("already won")
    ErrTicketNotFound         = errors.New("ticket not found")
    ErrTicketClosed           = errors.New("ticket closed")
    ErrInvalidName            = errors.New("invalid name")
)
```

### 5.6 Invariants

1. Lead email unique ต่อ tenant
2. Lead score 0-100
3. Convert ได้เฉพาะ lead ที่ `QUALIFIED`
4. Opportunity stage transition ต้องถูกต้องตาม state machine
5. Ticket SLA due คำนวณตาม priority
6. Ticket closed แก้ไม่ได้
7. Activity ต้องมี `ref_type` + `ref_id`

## 6. Application Layer

### 6.1 Use Cases

```go
// application/convert_lead.go — orchestration ข้ามโมดูลผ่าน Kafka
type ConvertLeadUseCase struct {
    leadRepo  repository.LeadRepository
    oppRepo   repository.OpportunityRepository
    producer  messaging.Producer
    tx        *transaction.Manager
}

type ConvertLeadInput struct {
    TenantID uuid.UUID
    LeadID   uuid.UUID
    ActorID  uuid.UUID
    InitialOpportunity struct {
        Name   string
        Amount float64
        Currency string
    }
}

func (uc *ConvertLeadUseCase) Execute(ctx context.Context, in ConvertLeadInput) error {
    return uc.tx.Do(ctx, func(tx *gorm.DB) error {
        lead, err := uc.leadRepo.FindByID(ctx, in.TenantID, in.LeadID)
        if err != nil { return err }

        // Note: ไม่ได้สร้าง customer โดยตรง — ส่ง event
        // customer module จะ consume แล้ว create + callback
        if err := lead.Qualify(); err != nil { return err }
        if err := uc.leadRepo.Save(ctx, lead); err != nil { return err }

        _ = uc.producer.Publish(ctx, "crm.lead.converted", lead.ID.String(), map[string]any{
            "lead_id":   lead.ID.String(),
            "tenant_id": in.TenantID.String(),
            "name":      lead.Name,
            "email":     lead.Email,
            "phone":     lead.Phone,
            "company":   lead.Company,
        })
        return nil
    })
}
```

```go
// application/classify_ticket.go — ใช้ LLM
type ClassifyTicketUseCase struct {
    ticketRepo repository.TicketRepository
    llm        llm.Client
}

func (uc *ClassifyTicketUseCase) Execute(ctx context.Context, in ClassifyTicketInput) error {
    t, err := uc.ticketRepo.FindByID(ctx, in.TenantID, in.TicketID)
    if err != nil { return err }

    prompt := fmt.Sprintf(`Classify this IoT support ticket into one category:
DEVICE_OFFLINE, SENSOR_FAULT, NETWORK, POWER, CONFIG, BILLING, OTHER.
Title: %s
Description: %s
Reply with just the category name.`, t.Subject, t.Description)

    category, err := uc.llm.Generate(ctx, prompt)
    if err != nil { return err }
    t.Category = strings.TrimSpace(category)
    return uc.ticketRepo.Save(ctx, t)
}
```

### 6.2 DTOs

```go
type LeadResponse struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Source   string `json:"source"`
    Score    int    `json:"score"`
    Status   string `json:"status"`
}

type TicketResponse struct {
    ID         string     `json:"id"`
    Subject    string     `json:"subject"`
    Priority   string     `json:"priority"`
    Status     string     `json:"status"`
    SLADueAt   time.Time  `json:"sla_due_at"`
    AssignedTo *string    `json:"assigned_to,omitempty"`
    ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type TicketModel struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
    TenantID    uuid.UUID  `gorm:"type:uuid;not null;index"`
    CustomerID  uuid.UUID  `gorm:"type:uuid;not null;index"`
    DeviceID    *uuid.UUID `gorm:"type:uuid;index"`
    Subject     string     `gorm:"size:255"`
    Description string     `gorm:"type:text"`
    Category    string     `gorm:"size:50"`
    Priority    string     `gorm:"size:20"`
    Status      string     `gorm:"size:20;index"`
    SLADueAt    time.Time  `gorm:"index"`
    AssignedTo  *uuid.UUID
    ResolvedAt  *time.Time
    CreatedAt   time.Time
}
func (TicketModel) TableName() string { return "crm_tickets" }
```

### 7.2 Kafka Consumers

```go
// consumers/iot_alert_consumer.go — สร้าง ticket จาก device alert
type IoTAlertConsumer struct{ createUC *application.CreateTicketUseCase }

func (c *IoTAlertConsumer) Handle(ctx context.Context, payload map[string]any) error {
    // map alert → ticket
    return c.createUC.Execute(ctx, application.CreateTicketInput{
        TenantID:    uuid.MustParse(payload["tenant_id"].(string)),
        CustomerID:  uuid.MustParse(payload["customer_id"].(string)),
        DeviceID:    parseUUIDPtr(payload["device_id"]),
        Subject:     "Device Alert: " + payload["metric"].(string),
        Description: payload["message"].(string),
        Priority:    mapSeverityToPriority(payload["severity"].(string)),
    })
}
```

### 7.3 WebSocket Broadcaster
- Lead ใหม่ → WS ไปยัง sales team
- Ticket SLA breach → WS ไปยัง manager
- Ticket resolved → WS ไปยัง customer

### 7.4 JWT, Bcrypt, Rate Limit
- Role-based: `sales`, `cs`, `manager`
- Webhook endpoint (lead intake) — rate limit 1000 req/min

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *TicketHandler) Assign(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))
    if err := h.assignUC.Execute(c.Request.Context(), application.AssignTicketInput{
        TenantID: tid, TicketID: id, ActorID: uid,
    }); err != nil { writeDomainError(c, err); return }
    c.JSON(200, gin.H{"message": "assigned"})
}
```

### 8.2 Routes

```go
l := r.Group("/leads"); l.Use(auth, tenant)
l.GET("", h.Lead.List); l.POST("", h.Lead.Create)
l.POST("/:id/convert", h.Lead.Convert)
l.POST("/:id/score", h.Lead.Score)

o := r.Group("/opportunities"); o.Use(auth, tenant)
o.POST("", h.Opp.Create); o.POST("/:id/move", h.Opp.MoveStage)
o.POST("/:id/win", h.Opp.Win); o.POST("/:id/lose", h.Opp.Lose)

t := r.Group("/tickets"); t.Use(auth, tenant)
t.GET("", h.Ticket.List); t.POST("", h.Ticket.Create)
t.POST("/:id/assign", h.Ticket.Assign)
t.POST("/:id/resolve", h.Ticket.Resolve)
t.POST("/:id/classify", h.Ticket.Classify)  // AI
```

### 8.3 Middleware
- `role.Require("sales")` สำหรับ lead/opp
- `role.Require("cs", "manager")` สำหรับ ticket

## 9. Database Migrations

```sql
-- migrations/20260104_crm_init.sql
CREATE TABLE IF NOT EXISTS crm_leads (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL,
    source                VARCHAR(30),
    name                  VARCHAR(255),
    email                 VARCHAR(255),
    phone                 VARCHAR(30),
    company               VARCHAR(255),
    interest_category     VARCHAR(30),
    score                 INT DEFAULT 0,
    status                VARCHAR(20) DEFAULT 'NEW',
    assigned_to           UUID,
    converted_customer_id UUID,
    metadata              JSONB,
    created_at            TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);
CREATE INDEX idx_crm_leads_tenant_status ON crm_leads(tenant_id, status);

CREATE TABLE IF NOT EXISTS crm_opportunities (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL,
    customer_id         UUID,
    lead_id             UUID,
    name                VARCHAR(255),
    stage               VARCHAR(30),
    amount              NUMERIC(15,2),
    currency            VARCHAR(3) DEFAULT 'THB',
    probability         INT DEFAULT 0,
    expected_close_date DATE,
    owner_id            UUID,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_crm_opps_tenant_stage ON crm_opportunities(tenant_id, stage);

CREATE TABLE IF NOT EXISTS crm_activities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    ref_type    VARCHAR(30),
    ref_id      UUID,
    activity_type VARCHAR(30),
    subject     VARCHAR(255),
    notes       TEXT,
    performed_by UUID,
    performed_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_crm_activities_ref ON crm_activities(ref_type, ref_id);

CREATE TABLE IF NOT EXISTS crm_tickets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    customer_id UUID NOT NULL,
    device_id   UUID,
    subject     VARCHAR(255),
    description TEXT,
    category    VARCHAR(50),
    priority    VARCHAR(20) DEFAULT 'NORMAL',
    status      VARCHAR(20) DEFAULT 'OPEN',
    sla_due_at  TIMESTAMP,
    assigned_to UUID,
    resolved_at TIMESTAMP,
    created_at  TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_crm_tickets_tenant_status ON crm_tickets(tenant_id, status);
CREATE INDEX idx_crm_tickets_sla ON crm_tickets(status, sla_due_at);

CREATE TABLE IF NOT EXISTS crm_campaigns (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    name        VARCHAR(255),
    channel     VARCHAR(50),
    start_date  DATE,
    end_date    DATE,
    budget      NUMERIC(15,2),
    status      VARCHAR(20),
    created_at  TIMESTAMP DEFAULT NOW()
);
```

## 10. System Flow

```
Lead Intake:
  Web form → POST /leads → CreateLeadUseCase
    → LeadScoringService.Score()
    → Kafka: crm.lead.created
    → Assignment (round-robin/zone)

Conversion:
  Sales convert → POST /leads/:id/convert
    → lead.Qualify() + lead.Convert()
    → Kafka: crm.lead.converted
    → customer module: CreateCustomerUseCase
    
IoT Alert → Ticket:
  Kafka: iot.alert.triggered
    → IoTAlertConsumer → CreateTicketUseCase
    → SLAService.Calculate()
    → Kafka: crm.ticket.created
    → notifier
```

## 11. Workflow Diagram

```
[Lead]──score──▶[Qualified]──convert──▶[Customer+Opportunity]
                                              │
                            ┌─────────────────┼─────────────────┐
                            ▼                 ▼                 ▼
                         Proposal        Negotiation         Won
                                                                  │
                                                                  ▼
                                                            Kafka: opp.won
                                                                  │
                                                                  ▼
                                                          package.subscribe
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=crm
SLA_CRITICAL_HOURS=1
SLA_HIGH_HOURS=4
LLM_ENABLED=true
go run ./cmd/api
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue | ไม่มีตรง — เป็น support function |
| Contribution | เพิ่ม conversion rate, ลด churn, ticket deflection ด้วย AI |
| KPI | Lead-to-customer %, Win rate, Avg SLA, CSAT |

## 14. ภาคผนวก
- Lead assignment: round-robin + zone-based
- SLA breach: escalate ไป manager
- AI classify: ใช้ `pkg/llm` กับ Ollama/OpenAI

## 15. Prompt ขยายระบบ

```
เพิ่ม use case AILeadScoring ใน module crm
- ใช้ LLM วิเคราะห์ lead (จาก email, company, notes)
- ให้ score 0-100 + reasoning
- เก็บ reasoning ใน lead.metadata
```

## 16. Prefix ตาราง

| Prefix | ตาราง |
|---|---|
| `crm_` | leads, opportunities, activities, tickets, campaigns |

## 17. DDD Validation Checklist

- [x] Lead/Opportunity/Ticket เป็น Aggregate Root แยกกัน
- [x] Stage transition ถูกต้องตาม state machine
- [x] SLA คำนวณใน domain service
- [x] AI classify ผ่าน outbound port (llm.Client)
- [x] Convert lead → Kafka (ไม่ import customer ตรง)
- [x] Activity ref ผ่าน polymorphic ref_type/ref_id

---

# PART 5 — จัดการอุปกรณ์ IoT, AI, Automation (`iotdevice`)

## 1. ภาพรวมระบบ

**Core Domain** ของแพลตฟอร์ม: จัดการ lifecycle ของอุปกรณ์ IoT (register → provision → online → offline), รับ telemetry ผ่าน MQTT (hot path), ส่งคำสั่งควบคุม, ตั้ง alert rules, รัน AI inference, และ Automation Engine (event → condition → action)

## 2. โครงสร้าง Module

```
internal/modules/iotdevice/
├── domain/
│   ├── entity/
│   │   ├── device.go                  # Aggregate Root
│   │   ├── device_model.go
│   │   ├── device_group.go
│   │   ├── firmware.go
│   │   ├── telemetry.go
│   │   ├── command.go
│   │   ├── alert_rule.go
│   │   ├── ai_model.go
│   │   └── automation.go              # Automation rules
│   ├── value_object/
│   │   ├── device_status.go
│   │   ├── protocol.go
│   │   ├── device_type.go
│   │   ├── metric.go
│   │   ├── capability.go
│   │   └── automation_trigger.go
│   ├── repository/
│   │   ├── device_repository.go
│   │   ├── telemetry_repository.go
│   │   ├── command_repository.go
│   │   ├── alert_repository.go
│   │   └── automation_repository.go
│   ├── service/
│   │   ├── provisioning_service.go
│   │   ├── health_monitor_service.go
│   │   ├── alert_evaluator.go
│   │   └── automation_engine.go
│   └── errors/errors.go
├── application/
│   ├── register_device.go
│   ├── provision_device.go
│   ├── send_command.go
│   ├── ingest_telemetry.go            # hot path
│   ├── evaluate_alert.go
│   ├── ota_update_firmware.go
│   ├── create_alert_rule.go
│   ├── create_automation.go
│   ├── run_ai_inference.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/influxdb/telemetry_repo.go
│   ├── persistence/redis/device_state_cache.go
│   ├── messaging/mqtt/
│   │   ├── broker.go
│   │   ├── topic_router.go
│   │   └── publisher.go
│   ├── messaging/kafka/
│   │   ├── telemetry_producer.go
│   │   └── consumers/
│   │       ├── telemetry_consumer.go
│   │       └── command_ack_consumer.go
│   ├── search/elasticsearch/device_indexer.go
│   ├── services/ai/
│   │   ├── inference_client.go
│   │   └── ollama_client.go
│   └── scheduler/offline_detector_job.go
├── interfaces/
│   ├── http/device_handler.go
│   ├── http/telemetry_handler.go
│   ├── http/automation_handler.go
│   ├── websocket/live_telemetry_hub.go
│   └── routes.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
                ┌───────────────────────────┐
                │   IoTDevice Context       │
                │   [CORE DOMAIN]           │
                │                            │
                │  ┌──────────┐  ┌────────┐ │
                │  │ Telemetry│  │Command │ │
                │  └──────────┘  └────────┘ │
                │  ┌──────────┐  ┌────────┐ │
                │  │ AI Engine│  │Automate│ │
                │  └──────────┘  └────────┘ │
                └───┬────────┬──────────┬───┘
                    │        │          │
       (device.registered) (telemetry) (alert.triggered)
                    │        │          │
                    ▼        ▼          ▼
              ┌─────────┐ ┌────────┐ ┌──────────┐
              │Package  │ │Report  │ │CRM/Notify│
              └─────────┘ └────────┘ └──────────┘
                    ▲
                    │ (subscription.created)
              ┌─────┴─────┐
              │ Package   │
              └───────────┘
```

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Device | อุปกรณ์ IoT |
| Device Model | รุ่นอุปกรณ์ |
| Device Group | กลุ่มอุปกรณ์ (zone/floor) |
| Provision | ตั้งค่าอุปกรณ์ให้พร้อมใช้งาน |
| Telemetry | ข้อมูลที่อุปกรณ์ส่งเข้ามา |
| Command | คำสั่งควบคุม |
| Alert Rule | กฎแจ้งเตือน |
| Automation | กฎ event → condition → action |
| AI Model | โมเดล AI (anomaly, forecast) |
| OTA | Over-the-Air firmware update |

## 5. Domain Layer

### 5.1 Entities

```go
// Device — Aggregate Root
type Device struct {
    ID               uuid.UUID
    TenantID         uuid.UUID
    CustomerID       uuid.UUID
    SiteID           uuid.UUID
    ModelID          *uuid.UUID
    SerialNo         string
    Name             string
    Status           valueobject.DeviceStatus
    Protocol         valueobject.Protocol
    MQTTClientID     string
    DeviceTokenHash  string
    FirmwareVersion  string
    LastSeenAt       *time.Time
    InstalledAt      *time.Time
    GroupID          *uuid.UUID
    Metadata         map[string]any
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func NewDevice(tenantID, customerID, siteID uuid.UUID, serial string, protocol valueobject.Protocol) (*Device, error) {
    if serial == "" { return nil, domainerrors.ErrInvalidSerial }
    if !protocol.IsValid() { return nil, domainerrors.ErrInvalidProtocol }
    return &Device{
        ID: uuid.New(), TenantID: tenantID, CustomerID: customerID,
        SiteID: siteID, SerialNo: serial, Protocol: protocol,
        Status: valueobject.DeviceStatusProvisioned,
        CreatedAt: time.Now(), UpdatedAt: time.Now(),
    }, nil
}

func (d *Device) Provision(token string, mqttClientID string) error {
    if d.Status != valueobject.DeviceStatusRegistered {
        return domainerrors.ErrDeviceAlreadyProvisioned
    }
    d.DeviceTokenHash = token
    d.MQTTClientID = mqttClientID
    d.Status = valueobject.DeviceStatusProvisioned
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) MarkOnline() error {
    now := time.Now()
    d.LastSeenAt = &now
    d.Status = valueobject.DeviceStatusOnline
    d.UpdatedAt = now
    return nil
}

func (d *Device) MarkOffline() error {
    if d.Status == valueobject.DeviceStatusOffline { return nil }
    d.Status = valueobject.DeviceStatusOffline
    d.UpdatedAt = time.Now()
    return nil
}

func (d *Device) ReportFault(reason string) error {
    d.Status = valueobject.DeviceStatusFault
    if d.Metadata == nil { d.Metadata = map[string]any{} }
    d.Metadata["fault_reason"] = reason
    return nil
}

// Telemetry — Entity (event)
type Telemetry struct {
    DeviceID  uuid.UUID
    TenantID  uuid.UUID
    Metric    valueobject.Metric
    Value     float64
    Quality   string
    Timestamp time.Time
}

// Command — Aggregate Root
type Command struct {
    ID        uuid.UUID
    DeviceID  uuid.UUID
    Command   string
    Payload   map[string]any
    Status    valueobject.CommandStatus
    IssuedBy  uuid.UUID
    IssuedAt  time.Time
    AckedAt   *time.Time
    ErrorMsg  string
}

func (c *Command) Ack(errMsg string) error {
    if c.Status != valueobject.CommandStatusSent { return domainerrors.ErrInvalidCommandState }
    now := time.Now()
    if errMsg != "" {
        c.Status = valueobject.CommandStatusFailed
        c.ErrorMsg = errMsg
    } else {
        c.Status = valueobject.CommandStatusAcked
    }
    c.AckedAt = &now
    return nil
}

// AlertRule — Aggregate Root
type AlertRule struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    SiteID    *uuid.UUID
    DeviceID  *uuid.UUID
    Metric    valueobject.Metric
    Operator  string   // ">", "<", ">=", "<=", "=="
    Threshold float64
    Severity  string   // INFO|WARN|CRITICAL
    Actions   []AlertAction
    IsActive  bool
}

func (r *AlertRule) Evaluate(v float64) bool {
    switch r.Operator {
    case ">":  return v > r.Threshold
    case "<":  return v < r.Threshold
    case ">=": return v >= r.Threshold
    case "<=": return v <= r.Threshold
    case "==": return v == r.Threshold
    }
    return false
}

// Automation — Aggregate Root (event-condition-action)
type Automation struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    Name       string
    Trigger    valueobject.AutomationTrigger
    Conditions []Condition
    Actions    []Action
    IsActive   bool
    CreatedAt  time.Time
}

type AutomationTrigger struct {
    Type     string   // TELEMETRY | SCHEDULE | COMMAND_ACK | ALERT
    DeviceID *uuid.UUID
    Metric   *valueobject.Metric
}

type Condition struct {
    Field    string
    Operator string
    Value    any
}

type Action struct {
    Type     string   // SEND_COMMAND | SEND_EMAIL | SEND_WS | WEBHOOK
    DeviceID *uuid.UUID
    Payload  map[string]any
}

// AIModel — Entity
type AIModel struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Name        string
    Type        string   // ANOMALY | FORECAST | CLASSIFY
    Provider    string   // OLLAMA | OPENAI | CUSTOM
    Endpoint    string
    Config      map[string]any
    IsActive    bool
    CreatedAt   time.Time
}
```

### 5.2 Value Objects

```go
type DeviceStatus string
const (
    DeviceStatusRegistered  DeviceStatus = "REGISTERED"
    DeviceStatusProvisioned DeviceStatus = "PROVISIONED"
    DeviceStatusOnline      DeviceStatus = "ONLINE"
    DeviceStatusOffline     DeviceStatus = "OFFLINE"
    DeviceStatusFault       DeviceStatus = "FAULT"
    DeviceStatusDecommissioned DeviceStatus = "DECOMMISSIONED"
)

type Protocol string
const (
    ProtocolMQTT   Protocol = "MQTT"
    ProtocolHTTP   Protocol = "HTTP"
    ProtocolLoRaWAN Protocol = "LORAWAN"
    ProtocolModbus Protocol = "MODBUS"
    ProtocolZigbee Protocol = "ZIGBEE"
)

type DeviceType string
const (
    DeviceTypeSensor   DeviceType = "SENSOR"
    DeviceTypeGateway  DeviceType = "GATEWAY"
    DeviceTypeActuator DeviceType = "ACTUATOR"
    DeviceTypeCamera   DeviceType = "CAMERA"
)

type Metric string
const (
    MetricTemperature Metric = "temperature"
    MetricHumidity    Metric = "humidity"
    MetricSoilMoisture Metric = "soil_moisture"
    MetricPower       Metric = "power"
    MetricEnergy      Metric = "energy"
    MetricCO2         Metric = "co2"
    MetricPressure    Metric = "pressure"
    MetricWaterFlow   Metric = "water_flow"
)

type CommandStatus string
const (
    CommandStatusPending CommandStatus = "PENDING"
    CommandStatusSent    CommandStatus = "SENT"
    CommandStatusAcked   CommandStatus = "ACKED"
    CommandStatusFailed  CommandStatus = "FAILED"
)
```

### 5.3 Repository Interfaces

```go
type DeviceRepository interface {
    Save(ctx context.Context, d *entity.Device) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Device, error)
    FindBySerial(ctx context.Context, tenantID uuid.UUID, serial string) (*entity.Device, error)
    ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Device, error)
    ListOfflineCandidates(ctx context.Context, threshold time.Time) ([]*entity.Device, error)
}

type TelemetryRepository interface {  // InfluxDB
    Write(ctx context.Context, t *entity.Telemetry) error
    Query(ctx context.Context, f TelemetryFilter) ([]*entity.Telemetry, error)
}

type CommandRepository interface {
    Save(ctx context.Context, c *entity.Command) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Command, error)
    ListPending(ctx context.Context, deviceID uuid.UUID) ([]*entity.Command, error)
}

type AlertRuleRepository interface {
    Save(ctx context.Context, r *entity.AlertRule) error
    ListActive(ctx context.Context, tenantID uuid.UUID, deviceID uuid.UUID) ([]*entity.AlertRule, error)
}

type AutomationRepository interface {
    Save(ctx context.Context, a *entity.Automation) error
    ListActiveByTrigger(ctx context.Context, tenantID uuid.UUID, triggerType string) ([]*entity.Automation, error)
}
```

### 5.4 Domain Services

```go
// AlertEvaluator
type AlertEvaluator struct{ repo repository.AlertRuleRepository }
func (s *AlertEvaluator) Evaluate(ctx context.Context, t *entity.Telemetry) ([]*entity.AlertRule, error) {
    rules, err := s.repo.ListActive(ctx, t.TenantID, t.DeviceID)
    if err != nil { return nil, err }
    triggered := make([]*entity.AlertRule, 0)
    for _, r := range rules {
        if r.Metric == t.Metric && r.Evaluate(t.Value) {
            triggered = append(triggered, r)
        }
    }
    return triggered, nil
}

// AutomationEngine — event → condition → action
type AutomationEngine struct {
    repo    repository.AutomationRepository
    cmdUC   *application.SendCommandUseCase
    notify  Notifier
}

func (e *AutomationEngine) OnTelemetry(ctx context.Context, t *entity.Telemetry) error {
    automations, err := e.repo.ListActiveByTrigger(ctx, t.TenantID, "TELEMETRY")
    if err != nil { return err }
    for _, a := range automations {
        if matchConditions(a.Conditions, t) {
            for _, act := range a.Actions {
                if err := e.execute(ctx, act, t); err != nil {
                    log.Printf("automation %s action failed: %v", a.ID, err)
                }
            }
        }
    }
    return nil
}
```

### 5.5 Domain Errors

```go
var (
    ErrDeviceNotFound          = errors.New("device not found")
    ErrInvalidSerial           = errors.New("invalid serial")
    ErrInvalidProtocol         = errors.New("invalid protocol")
    ErrDeviceAlreadyProvisioned = errors.New("device already provisioned")
    ErrDeviceNotProvisioned    = errors.New("device not provisioned")
    ErrInvalidCommandState     = errors.New("invalid command state")
    ErrAlertRuleNotFound       = errors.New("alert rule not found")
    ErrAutomationNotFound      = errors.New("automation not found")
    ErrTelemetryWriteFailed    = errors.New("telemetry write failed")
)
```

### 5.6 Invariants

1. Serial No unique ต่อ tenant
2. Device ต้องมี `customer_id` + `site_id`
3. Provisioning ต้องสร้าง token hash (bcrypt)
4. `LastSeenAt` update ทุกครั้งที่รับ telemetry/status
5. Command ต้อง ACK ภายใน timeout (default 30s)
6. Alert rule threshold + operator ต้อง valid
7. Automation ต้องมีอย่างน้อย 1 action
8. Telemetry ต้องมี timestamp จาก device (fallback: server time)

## 6. Application Layer

### 6.1 Use Cases

```go
// application/register_device.go
type RegisterDeviceUseCase struct {
    repo     repository.DeviceRepository
    producer messaging.Producer
    entitler *entitlement.Checker   // ตรวจ quota
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*DeviceResponse, error) {
    if err := uc.entitler.CheckDeviceQuota(ctx, in.TenantID, in.CustomerID); err != nil {
        return nil, err
    }

    d, err := entity.NewDevice(in.TenantID, in.CustomerID, in.SiteID, in.SerialNo, in.Protocol)
    if err != nil { return nil, err }

    if err := uc.repo.Save(ctx, d); err != nil { return nil, err }

    _ = uc.producer.Publish(ctx, "iot.device.registered", d.ID.String(), map[string]any{
        "device_id":   d.ID.String(),
        "customer_id": d.CustomerID.String(),
        "site_id":     d.SiteID.String(),
        "serial_no":   d.SerialNo,
    })
    return toDeviceResponse(d), nil
}
```

```go
// application/provision_device.go
type ProvisionDeviceUseCase struct {
    repo    repository.DeviceRepository
    crypter cryptpass.Crypter
    producer messaging.Producer
}

func (uc *ProvisionDeviceUseCase) Execute(ctx context.Context, in ProvisionDeviceInput) (*ProvisionResponse, error) {
    d, err := uc.repo.FindByID(ctx, in.TenantID, in.DeviceID)
    if err != nil { return nil, err }

    rawToken, _ := secureRandom.String(32)
    hash, _ := uc.crypter.Hash(rawToken)
    mqttClientID := fmt.Sprintf("dev-%s", d.ID.String()[:8])

    if err := d.Provision(hash, mqttClientID); err != nil { return nil, err }
    if err := uc.repo.Save(ctx, d); err != nil { return nil, err }

    // ส่ง token กลับเป็นครั้งเดียว — ไม่เก็บ plaintext
    return &ProvisionResponse{
        DeviceID:     d.ID.String(),
        MQTTClientID: mqttClientID,
        DeviceToken:  rawToken,
        MQTTBroker:   in.MQTTBroker,
    }, nil
}
```

```go
// application/ingest_telemetry.go — hot path, ห้าม block นาน
type IngestTelemetryUseCase struct {
    telemetryRepo repository.TelemetryRepository   // InfluxDB
    deviceRepo    repository.DeviceRepository
    stateCache    DeviceStateCache                 // Redis
    evaluator     *service.AlertEvaluator
    automation    *service.AutomationEngine
    producer      messaging.Producer
}

type IngestTelemetryInput struct {
    TenantID  uuid.UUID
    DeviceID  uuid.UUID
    Metrics   []MetricValue
    Timestamp time.Time
}

func (uc *IngestTelemetryUseCase) Execute(ctx context.Context, in IngestTelemetryInput) error {
    // 1. Fast state update (Redis)
    _ = uc.stateCache.SetLastSeen(ctx, in.DeviceID, in.Timestamp)

    // 2. Write to InfluxDB (batch)
    points := make([]*entity.Telemetry, 0, len(in.Metrics))
    for _, m := range in.Metrics {
        points = append(points, &entity.Telemetry{
            DeviceID: in.DeviceID, TenantID: in.TenantID,
            Metric: m.Metric, Value: m.Value, Timestamp: in.Timestamp,
        })
    }
    if err := uc.telemetryRepo.Write(ctx, points...); err != nil {
        return domainerrors.ErrTelemetryWriteFailed
    }

    // 3. Alert evaluation (fast, cached rules)
    for _, p := range points {
        triggered, _ := uc.evaluator.Evaluate(ctx, p)
        for _, r := range triggered {
            _ = uc.producer.Publish(ctx, "iot.alert.triggered", in.DeviceID.String(), map[string]any{
                "device_id":   in.DeviceID.String(),
                "metric":      string(r.Metric),
                "value":       p.Value,
                "threshold":   r.Threshold,
                "severity":    r.Severity,
                "customer_id": in.CustomerID.String(),
            })
        }
        // 4. Automation
        _ = uc.automation.OnTelemetry(ctx, p)
    }

    // 5. Aggregate → Kafka for report/usage
    _ = uc.producer.Publish(ctx, "iot.telemetry.aggregated", in.DeviceID.String(), map[string]any{
        "device_id":   in.DeviceID.String(),
        "tenant_id":   in.TenantID.String(),
        "customer_id": in.CustomerID.String(),
        "count":       len(points),
        "window":      in.Timestamp,
    })
    return nil
}
```

```go
// application/send_command.go
type SendCommandUseCase struct {
    deviceRepo  repository.DeviceRepository
    cmdRepo     repository.CommandRepository
    mqtt        mqtt.Publisher
    producer    messaging.Producer
}

func (uc *SendCommandUseCase) Execute(ctx context.Context, in SendCommandInput) (*CommandResponse, error) {
    d, err := uc.deviceRepo.FindByID(ctx, in.TenantID, in.DeviceID)
    if err != nil { return nil, err }
    if d.Status == valueobject.DeviceStatusOffline {
        return nil, domainerrors.ErrDeviceOffline
    }

    cmd := entity.NewCommand(d.ID, in.Command, in.Payload, in.ActorID)
    if err := uc.cmdRepo.Save(ctx, cmd); err != nil { return nil, err }

    topic := fmt.Sprintf("iot/%s/%s/cmd", in.TenantID.String(), d.ID.String())
    if err := uc.mqtt.Publish(ctx, topic, cmd); err != nil {
        cmd.Status = valueobject.CommandStatusFailed
        cmd.ErrorMsg = err.Error()
        _ = uc.cmdRepo.Save(ctx, cmd)
        return nil, err
    }
    cmd.Status = valueobject.CommandStatusSent
    _ = uc.cmdRepo.Save(ctx, cmd)

    _ = uc.producer.Publish(ctx, "iot.command.sent", cmd.ID.String(), map[string]any{
        "command_id": cmd.ID.String(),
        "device_id":  d.ID.String(),
    })
    return toCommandResponse(cmd), nil
}
```

### 6.2 DTOs

```go
type DeviceResponse struct {
    ID          string     `json:"id"`
    SerialNo    string     `json:"serial_no"`
    Name        string     `json:"name"`
    Status      string     `json:"status"`
    Protocol    string     `json:"protocol"`
    SiteID      string     `json:"site_id"`
    CustomerID  string     `json:"customer_id"`
    LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
    FirmwareVersion string `json:"firmware_version,omitempty"`
}

type ProvisionResponse struct {
    DeviceID     string `json:"device_id"`
    MQTTClientID string `json:"mqtt_client_id"`
    DeviceToken  string `json:"device_token"`
    MQTTBroker   string `json:"mqtt_broker"`
}

type TelemetryResponse struct {
    DeviceID  string    `json:"device_id"`
    Metric    string    `json:"metric"`
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
}

type CommandResponse struct {
    ID       string     `json:"id"`
    DeviceID string     `json:"device_id"`
    Command  string     `json:"command"`
    Status   string     `json:"status"`
    IssuedAt time.Time  `json:"issued_at"`
    AckedAt  *time.Time `json:"acked_at,omitempty"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type DeviceModel struct {
    ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
    TenantID        uuid.UUID  `gorm:"type:uuid;not null;index"`
    CustomerID      uuid.UUID  `gorm:"type:uuid;not null;index"`
    SiteID          uuid.UUID  `gorm:"type:uuid;not null;index"`
    ModelID         *uuid.UUID `gorm:"type:uuid"`
    SerialNo        string     `gorm:"size:100;not null;index:idx_dev_serial,unique"`
    Name            string     `gorm:"size:255"`
    Status          string     `gorm:"size:20;index"`
    Protocol        string     `gorm:"size:30"`
    MQTTClientID    string     `gorm:"size:100"`
    DeviceTokenHash string     `gorm:"size:255"`
    FirmwareVersion string     `gorm:"size:50"`
    LastSeenAt      *time.Time
    InstalledAt     *time.Time
    GroupID         *uuid.UUID
    Metadata        datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
func (DeviceModel) TableName() string { return "iotdevice_devices" }
```

**InfluxDB Telemetry Repository:**
```go
// infrastructure/persistence/influxdb/telemetry_repo.go
type telemetryRepo struct{ client influxdb2.Client; org, bucket string }

func (r *telemetryRepo) Write(ctx context.Context, points ...*entity.Telemetry) error {
    writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)
    for _, t := range points {
        p := influxdb2.NewPointWithMeasurement("device_telemetry").
            AddTag("tenant_id", t.TenantID.String()).
            AddTag("device_id", t.DeviceID.String()).
            AddTag("metric", string(t.Metric)).
            AddField("value", t.Value).
            SetTime(t.Timestamp)
        if err := writeAPI.WritePoint(ctx, p); err != nil { return err }
    }
    return nil
}
```

### 7.2 Kafka Consumers

```go
// consumers/telemetry_consumer.go — hot path
type TelemetryConsumer struct{ uc *application.IngestTelemetryUseCase }

func (c *TelemetryConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    batch := make([]*application.IngestTelemetryInput, 0, 500)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case msg, ok := <-claim.Messages():
            if !ok { flush(c.uc, batch); return nil }
            var in application.IngestTelemetryInput
            if err := json.Unmarshal(msg.Value, &in); err != nil { continue }
            batch = append(batch, &in)
            if len(batch) >= 500 { flush(c.uc, batch); batch = batch[:0] }
            s.MarkMessage(msg, "")
        case <-ticker.C:
            if len(batch) > 0 { flush(c.uc, batch); batch = batch[:0] }
        }
    }
}
```

### 7.3 WebSocket Broadcaster

```go
// interfaces/websocket/live_telemetry_hub.go
type LiveTelemetryHub struct{ hub *pkgws.Hub }

func (h *LiveTelemetryHub) Broadcast(tenantID, deviceID string, t *entity.Telemetry) {
    h.hub.Broadcast(tenantID, pkgws.Event{
        Type: "telemetry",
        Data: map[string]any{
            "device_id": deviceID, "metric": t.Metric,
            "value": t.Value, "ts": t.Timestamp,
        },
    })
}
```

### 7.4 JWT, Bcrypt, Rate Limit
- **JWT**: สำหรับ user → API
- **Bcrypt** (`pkg/cryptpass`): hash device token
- **Rate limit**: MQTT ingest แยก, ไม่ผ่าน HTTP

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *DeviceHandler) SendCommand(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req struct {
        Command string         `json:"command" binding:"required"`
        Payload map[string]any `json:"payload"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }

    res, err := h.cmdUC.Execute(c.Request.Context(), application.SendCommandInput{
        TenantID: tid, DeviceID: id, Command: req.Command,
        Payload: req.Payload, ActorID: uid,
    })
    if err != nil { writeDomainError(c, err); return }
    c.JSON(202, res)
}
```

### 8.2 Routes

```go
d := r.Group("/devices"); d.Use(auth, tenant)
d.GET("",                      h.Device.List)
d.POST("",                     h.Device.Register)
d.GET("/:id",                  h.Device.Get)
d.POST("/:id/provision",       h.Device.Provision)
d.POST("/:id/commands",        h.Device.SendCommand)
d.GET("/:id/telemetry",        h.Telemetry.Query)
d.GET("/:id/commands",         h.Command.List)
d.DELETE("/:id",               h.Device.Decommission)

s := r.Group("/sites"); s.Use(auth, tenant)
s.GET("/:site_id/devices",     h.Device.ListBySite)

a := r.Group("/alert-rules"); a.Use(auth, tenant)
a.POST("", h.Alert.Create); a.GET("", h.Alert.List); a.PUT("/:id", h.Alert.Update)

au := r.Group("/automations"); au.Use(auth, tenant)
au.POST("", h.Automation.Create); au.GET("", h.Automation.List)

ws := r.Group("/ws"); ws.Use(auth, tenant)
ws.GET("/devices/:id/live", h.WS.LiveTelemetry)
```

### 8.3 Middleware
- `auth` + `tenant`
- Device authentication (แยก): token ผ่าน MQTT CONNECT

## 9. Database Migrations

```sql
-- migrations/20260105_iotdevice_init.sql
CREATE TABLE IF NOT EXISTS iotdevice_models (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor       VARCHAR(100),
    model_no     VARCHAR(100) UNIQUE,
    device_type  VARCHAR(30),
    protocol     VARCHAR(30),
    capabilities JSONB,
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS iotdevice_devices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    customer_id       UUID NOT NULL,
    site_id           UUID NOT NULL,
    model_id          UUID REFERENCES iotdevice_models(id),
    serial_no         VARCHAR(100) NOT NULL,
    name              VARCHAR(255),
    status            VARCHAR(20) DEFAULT 'REGISTERED',
    protocol          VARCHAR(30),
    mqtt_client_id    VARCHAR(100),
    device_token_hash VARCHAR(255),
    firmware_version  VARCHAR(50),
    last_seen_at      TIMESTAMP,
    installed_at      TIMESTAMP,
    group_id          UUID,
    metadata          JSONB,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, serial_no)
);
CREATE INDEX idx_iotdevice_devices_tenant_customer ON iotdevice_devices(tenant_id, customer_id);
CREATE INDEX idx_iotdevice_devices_site           ON iotdevice_devices(site_id);
CREATE INDEX idx_iotdevice_devices_status_seen    ON iotdevice_devices(status, last_seen_at);

CREATE TABLE IF NOT EXISTS iotdevice_groups (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    site_id    UUID NOT NULL,
    name       VARCHAR(100),
    group_type VARCHAR(30),
    metadata   JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS iotdevice_commands (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    device_id  UUID NOT NULL,
    command    VARCHAR(100) NOT NULL,
    payload    JSONB,
    status     VARCHAR(20) DEFAULT 'PENDING',
    issued_by  UUID,
    issued_at  TIMESTAMP DEFAULT NOW(),
    acked_at   TIMESTAMP,
    error_msg  TEXT
);
CREATE INDEX idx_iotdevice_commands_device ON iotdevice_commands(device_id, issued_at DESC);

CREATE TABLE IF NOT EXISTS iotdevice_alert_rules (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    site_id    UUID,
    device_id  UUID,
    metric     VARCHAR(50) NOT NULL,
    operator   VARCHAR(10) NOT NULL,
    threshold  NUMERIC(15,4),
    severity   VARCHAR(20),
    actions    JSONB,
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_iotdevice_alerts_tenant_active ON iotdevice_alert_rules(tenant_id, is_active);

CREATE TABLE IF NOT EXISTS iotdevice_automations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    name       VARCHAR(255) NOT NULL,
    trigger    JSONB NOT NULL,
    conditions JSONB,
    actions    JSONB NOT NULL,
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS iotdevice_ai_models (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    name       VARCHAR(100),
    type       VARCHAR(30),   -- ANOMALY|FORECAST|CLASSIFY
    provider   VARCHAR(30),   -- OLLAMA|OPENAI|CUSTOM
    endpoint   VARCHAR(255),
    config     JSONB,
    is_active  BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Telemetry เก็บใน InfluxDB ไม่ใช่ PostgreSQL
-- measurement: device_telemetry
-- tags: tenant_id, device_id, metric
-- fields: value
```

## 10. System Flow

```
Device onboarding:
  Register → Provision (gen token) → Device connects MQTT → Online

Telemetry hot path:
  Device → MQTT publish iot/{tenant}/{device}/telemetry
    → TopicRouter → Kafka producer (iot.telemetry.raw)
      → TelemetryConsumer → IngestTelemetryUseCase
        ├─ Redis: LastSeen
        ├─ InfluxDB: write points
        ├─ AlertEvaluator (cached rules)
        ├─ AutomationEngine
        ├─ WS: live broadcast
        └─ Kafka: iot.telemetry.aggregated

Command:
  User → POST /devices/:id/commands
    → SendCommandUseCase → Save + MQTT publish
    → Device ACK → MQTT → Kafka iot.command.acked
```

## 11. Workflow Diagram (Automation)

```
[Telemetry: temp=45°C]
        │
        ▼
[AutomationEngine.OnTelemetry]
        │
        ▼
[Match Trigger: metric=temperature]
        │
        ▼
[Match Conditions: temp > 40]
        │
        ├─▶ [Action: Send Command fan=ON]
        │       │
        │       ▼
        │   [MQTT publish iot/.../cmd]
        │
        ├─▶ [Action: Send WS to dashboard]
        │
        └─▶ [Action: Create CRM ticket]
                │
                ▼
          Kafka: crm.ticket.create
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=iotdevice

# Env
MQTT_BROKER=tcp://mosquitto:1883
MQTT_USERNAME=iot-platform
MQTT_PASSWORD=***
INFLUX_URL=http://influxdb:8086
INFLUX_TOKEN=***
INFLUX_ORG=icmon
INFLUX_BUCKET=iot_telemetry
OFFLINE_THRESHOLD_MINUTES=5
COMMAND_TIMEOUT_SECONDS=30

# Run
go run ./cmd/api
go run ./cmd/workers/telemetry       # Kafka consumer
go run ./cmd/workers/mqtt-ingest     # MQTT → Kafka bridge
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue | ค่าบริการต่อ device/เดือน + AI inference + Automation |
| Usage-based | จำนวน telemetry points, API calls, AI calls |
| Tier | Basic (no AI), Pro (AI anomaly), Enterprise (custom) |
| KPI | Device uptime, MTBF, Alert accuracy, AI precision |

## 14. ภาคผนวก

- **MQTT Topic**: `iot/{tenant_id}/{device_id}/{action}`
- **QoS**: telemetry=0, command=1, firmware=1
- **Retention**: InfluxDB 90 วัน (raw), 5 ปี (aggregated)
- **AI Provider**: `pkg/llm` → Ollama (self-host) หรือ OpenAI
- **Automation**: จำกัด 100 rules/tenant, 5 actions/rule

## 15. Prompt ขยายระบบ

```
เพิ่ม use case RunAIAutomation ใน module iotdevice
- Input: automation_id, telemetry payload
- Business: เรียก AI model (anomaly detection) → ถ้า anomaly > threshold → execute actions
- Side effects: Kafka iot.ai.anomaly.detected, WS broadcast
```

## 16. Prefix ตาราง

| Prefix | ตาราง | Storage |
|---|---|---|
| `iotdevice_` | models, devices, groups, commands, alert_rules, automations, ai_models | PostgreSQL |
| — | `device_telemetry` (measurement) | InfluxDB |
| — | Device state (LastSeen, online) | Redis |

## 17. DDD Validation Checklist

- [x] Device เป็น Aggregate Root ที่ encapsulate lifecycle
- [x] Telemetry เป็น event (ไม่ mutate aggregate)
- [x] Command/Automation/AlertRule เป็น Aggregate Root แยก
- [x] Telemetry เก็บใน InfluxDB (ไม่ใช่ PG)
- [x] Hot path ใช้ batch + Redis cache
- [x] AI เป็น outbound port (`llm.Client`, `inference_client`)
- [x] Cross-module ผ่าน Kafka
- [x] Device auth แยกจาก user auth
- [x] Multi-tenant MQTT topic namespace

---

# PART 6 — บริหารจัดการ Logistic IoT (`iotlogistics`)

## 1. ภาพรวมระบบ

บริหาร **Logistics ของอุปกรณ์ IoT** ตลอด lifecycle: Shipment (ขนส่งไปติดตั้ง), Installation Job (ติดตั้งที่ site), Maintenance Schedule (PM), Work Order (ซ่อม), RMA (เคลม), Spare Parts, Technician Management + Route Optimization (AI)

## 2. โครงสร้าง Module

```
internal/modules/iotlogistics/
├── domain/
│   ├── entity/
│   │   ├── shipment.go               # Aggregate Root
│   │   ├── shipment_item.go
│   │   ├── installation_job.go       # Aggregate Root
│   │   ├── technician.go
│   │   ├── maintenance_schedule.go
│   │   ├── work_order.go
│   │   ├── rma.go                    # Aggregate Root
│   │   └── spare_part.go
│   ├── value_object/
│   │   ├── shipment_status.go
│   │   ├── job_status.go
│   │   ├── job_type.go
│   │   ├── gps_location.go
│   │   └── route.go
│   ├── repository/
│   │   ├── shipment_repository.go
│   │   ├── installation_repository.go
│   │   ├── maintenance_repository.go
│   │   ├── rma_repository.go
│   │   └── technician_repository.go
│   ├── service/
│   │   ├── route_planning_service.go
│   │   ├── geo_fence_service.go
│   │   ├── sla_calculator.go
│   │   └── technician_assignment_service.go
│   └── errors/errors.go
├── application/
│   ├── create_shipment.go
│   ├── dispatch_shipment.go
│   ├── confirm_delivery.go
│   ├── schedule_installation.go
│   ├── assign_technician.go
│   ├── start_installation.go
│   ├── complete_installation.go       # → iotdevice.provision
│   ├── schedule_maintenance.go
│   ├── create_work_order.go
│   ├── complete_work_order.go
│   ├── create_rma.go
│   ├── close_rma.go
│   ├── optimize_route.go              # AI
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── services/maps/google_maps.go
│   ├── services/llm/route_advisor.go
│   ├── messaging/kafka/
│   └── scheduler/pm_due_job.go
├── interfaces/
│   ├── http/
│   │   ├── shipment_handler.go
│   │   ├── installation_handler.go
│   │   ├── maintenance_handler.go
│   │   ├── rma_handler.go
│   │   └── routes.go
│   └── websocket/technician_tracking.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌──────────┐  (customer.onboarded)  ┌─────────────────────┐
│ Customer │───────────────────────▶│ IoTLogistics        │
└──────────┘                        │ [SUPPORTING]        │
                                     │                     │
┌──────────┐  (erp.product.created) │  Shipment           │
│   ERP    │────────────────────────▶│  Installation       │
└──────────┘                         │  Maintenance        │
                                     │  RMA                │
┌──────────┐                         └──┬──────────┬───────┘
│IoTDevice │◀──(installation.done)──────┘          │
└──────────┘                                        │
                                                    │ (maintenance.due)
                                                    ▼
                                            ┌──────────────┐
                                            │CRM/Notifier  │
                                            └──────────────┘
```

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Shipment | การขนส่งอุปกรณ์ไปยัง site |
| Shipment Item | รายการอุปกรณ์ในการขนส่ง |
| Installation Job | งานติดตั้งอุปกรณ์ |
| Technician | ช่างที่ติดตั้ง/ซ่อม |
| Work Order | ใบสั่งงาน (ซ่อม/บำรุง) |
| PM | Preventive Maintenance |
| RMA | Return Merchandise Authorization |
| Spare Part | อะไหล่ |
| Route | เส้นทางของช่าง |
| Geo-fence | พื้นที่ทำงาน |

## 5. Domain Layer

### 5.1 Entities

```go
// Shipment — Aggregate Root
type Shipment struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    ShipmentNo     string
    CustomerID     uuid.UUID
    SiteID         uuid.UUID
    FromWarehouseID uuid.UUID
    Carrier        string
    TrackingNo     string
    Status         valueobject.ShipmentStatus
    Items          []*ShipmentItem
    ScheduledAt    time.Time
    DeliveredAt    *time.Time
    GPSLast        *valueobject.GPSLocation
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func NewShipment(tenantID, customerID, siteID, warehouseID uuid.UUID, no string) (*Shipment, error) {
    if no == "" { return nil, domainerrors.ErrInvalidShipmentNo }
    return &Shipment{
        ID: uuid.New(), TenantID: tenantID, CustomerID: customerID,
        SiteID: siteID, FromWarehouseID: warehouseID, ShipmentNo: no,
        Status: valueobject.ShipmentStatusPending,
        Items: []*ShipmentItem{}, CreatedAt: time.Now(), UpdatedAt: time.Now(),
    }, nil
}

func (s *Shipment) AddItem(item *ShipmentItem) error {
    if s.Status != valueobject.ShipmentStatusPending {
        return domainerrors.ErrShipmentLocked
    }
    s.Items = append(s.Items, item)
    return nil
}

func (s *Shipment) Dispatch(carrier, trackingNo string) error {
    if s.Status != valueobject.ShipmentStatusPending {
        return domainerrors.ErrInvalidStatusTransition
    }
    if len(s.Items) == 0 { return domainerrors.ErrEmptyShipment }
    s.Carrier = carrier
    s.TrackingNo = trackingNo
    s.Status = valueobject.ShipmentStatusInTransit
    s.UpdatedAt = time.Now()
    return nil
}

func (s *Shipment) UpdateGPS(loc valueobject.GPSLocation) error {
    if !loc.IsValid() { return domainerrors.ErrInvalidGPS }
    s.GPSLast = &loc
    return nil
}

func (s *Shipment) ConfirmDelivery(signature string) error {
    if s.Status != valueobject.ShipmentStatusInTransit {
        return domainerrors.ErrInvalidStatusTransition
    }
    now := time.Now()
    s.DeliveredAt = &now
    s.Status = valueobject.ShipmentStatusDelivered
    return nil
}

// InstallationJob — Aggregate Root
type InstallationJob struct {
    ID           uuid.UUID
    TenantID     uuid.UUID
    ShipmentID   *uuid.UUID
    SiteID       uuid.UUID
    TechnicianID *uuid.UUID
    JobType      valueobject.JobType
    Status       valueobject.JobStatus
    ScheduledAt  time.Time
    StartedAt    *time.Time
    CompletedAt  *time.Time
    Checklist    []ChecklistItem
    Photos       []string
    SignatureURL string
    Notes        string
    CreatedAt    time.Time
}

func (j *InstallationJob) AssignTechnician(techID uuid.UUID) error {
    if j.Status != valueobject.JobStatusScheduled {
        return domainerrors.ErrInvalidStatusTransition
    }
    j.TechnicianID = &techID
    return nil
}

func (j *InstallationJob) Start() error {
    if j.Status != valueobject.JobStatusScheduled { return domainerrors.ErrInvalidStatusTransition }
    if j.TechnicianID == nil { return domainerrors.ErrNoTechnician }
    now := time.Now()
    j.StartedAt = &now
    j.Status = valueobject.JobStatusInProgress
    return nil
}

func (j *InstallationJob) Complete(checklist []ChecklistItem, photos []string, signature string) error {
    if j.Status != valueobject.JobStatusInProgress { return domainerrors.ErrInvalidStatusTransition }
    if !allChecked(checklist) { return domainerrors.ErrChecklistIncomplete }
    now := time.Now()
    j.CompletedAt = &now
    j.Status = valueobject.JobStatusDone
    j.Checklist = checklist
    j.Photos = photos
    j.SignatureURL = signature
    return nil
}

// WorkOrder — Aggregate Root
type WorkOrder struct {
    ID           uuid.UUID
    TenantID     uuid.UUID
    WONo         string
    DeviceID     *uuid.UUID
    SiteID       uuid.UUID
    Type         valueobject.JobType   // MAINTENANCE | REPAIR | REMOVAL
    Priority     string
    Status       valueobject.JobStatus
    AssignedTo   *uuid.UUID
    Description  string
    ScheduledAt  time.Time
    CompletedAt  *time.Time
    PartsUsed    []*SparePartUsage
    CreatedAt    time.Time
}

// RMA — Aggregate Root
type RMA struct {
    ID                  uuid.UUID
    TenantID            uuid.UUID
    RMANo               string
    DeviceID            uuid.UUID
    CustomerID          uuid.UUID
    Reason              string
    Status              valueobject.RMAStatus
    ReplacementDeviceID *uuid.UUID
    ShippedAt           *time.Time
    ReceivedAt          *time.Time
    ClosedAt            *time.Time
    CreatedAt           time.Time
}

func (r *RMA) Ship() error {
    if r.Status != valueobject.RMAStatusOpen { return domainerrors.ErrInvalidStatusTransition }
    now := time.Now()
    r.ShippedAt = &now
    r.Status = valueobject.RMAStatusShipped
    return nil
}

func (r *RMA) Receive() error { /* ... */ }
func (r *RMA) Close(notes string) error { /* ... */ }
```

### 5.2 Value Objects

```go
type ShipmentStatus string
const (
    ShipmentStatusPending   ShipmentStatus = "PENDING"
    ShipmentStatusInTransit ShipmentStatus = "IN_TRANSIT"
    ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
    ShipmentStatusInstalled ShipmentStatus = "INSTALLED"
    ShipmentStatusCancelled ShipmentStatus = "CANCELLED"
)

type JobStatus string
const (
    JobStatusScheduled  JobStatus = "SCHEDULED"
    JobStatusInProgress JobStatus = "IN_PROGRESS"
    JobStatusDone       JobStatus = "DONE"
    JobStatusFailed     JobStatus = "FAILED"
    JobStatusCancelled  JobStatus = "CANCELLED"
)

type JobType string
const (
    JobTypeInstall     JobType = "INSTALL"
    JobTypeMaintenance JobType = "MAINTENANCE"
    JobTypeRepair      JobType = "REPAIR"
    JobTypeRemoval     JobType = "REMOVAL"
    JobTypeInspection  JobType = "INSPECTION"
)

type RMAStatus string
const (
    RMAStatusOpen      RMAStatus = "OPEN"
    RMAStatusShipped   RMAStatus = "SHIPPED"
    RMAStatusReceived  RMAStatus = "RECEIVED"
    RMAStatusRepaired  RMAStatus = "REPAIRED"
    RMAStatusReplaced  RMAStatus = "REPLACED"
    RMAStatusRejected  RMAStatus = "REJECTED"
    RMAStatusClosed    RMAStatus = "CLOSED"
)

type GPSLocation struct {
    Lat, Lng float64
    Accuracy float64
    Timestamp time.Time
}
func (g GPSLocation) IsValid() bool {
    return g.Lat >= -90 && g.Lat <= 90 && g.Lng >= -180 && g.Lng <= 180
}

type Route struct {
    Waypoints []GPSLocation
    Distance  float64   // km
    Duration  int       // minutes
}
```

### 5.3 Repository Interfaces

```go
type ShipmentRepository interface {
    Save(ctx context.Context, s *entity.Shipment) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Shipment, error)
    FindByNo(ctx context.Context, tenantID uuid.UUID, no string) (*entity.Shipment, error)
    ListByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.ShipmentStatus) ([]*entity.Shipment, error)
}

type InstallationRepository interface {
    Save(ctx context.Context, j *entity.InstallationJob) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.InstallationJob, error)
    ListByTechnician(ctx context.Context, tenantID, techID uuid.UUID, date time.Time) ([]*entity.InstallationJob, error)
    ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.InstallationJob, error)
}

type MaintenanceRepository interface {
    Save(ctx context.Context, m *entity.MaintenanceSchedule) error
    FindDue(ctx context.Context, tenantID uuid.UUID, asOf time.Time) ([]*entity.MaintenanceSchedule, error)
}

type RMARepository interface {
    Save(ctx context.Context, r *entity.RMA) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.RMA, error)
    ListByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.RMA, error)
}

type TechnicianRepository interface {
    Save(ctx context.Context, t *entity.Technician) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Technician, error)
    FindAvailable(ctx context.Context, tenantID uuid.UUID, zone string, skills []string) ([]*entity.Technician, error)
    UpdateGPS(ctx context.Context, id uuid.UUID, loc valueobject.GPSLocation) error
}
```

### 5.4 Domain Services

```go
// RoutePlanningService — ใช้ AI
type RoutePlanningService struct{ maps MapsClient }

func (s *RoutePlanningService) Plan(ctx context.Context, jobs []*entity.InstallationJob, startGPS valueobject.GPSLocation) (*valueobject.Route, error) {
    // Nearest-neighbor + traffic + working hours
    // ...
    return &valueobject.Route{}, nil
}

// TechnicianAssignmentService — มอบหมายงาน
type TechnicianAssignmentService struct{ repo repository.TechnicianRepository }

func (s *TechnicianAssignmentService) Assign(ctx context.Context, job *entity.InstallationJob, siteGPS valueobject.GPSLocation) (*entity.Technician, error) {
    techs, _ := s.repo.FindAvailable(ctx, job.TenantID, "", nil)
    // เลือกช่างที่ใกล้สุด + skill ตรง + workload น้อยสุด
    return techs[0], nil
}

// SLACalculator
type SLACalculator struct{}
func (s *SLACalculator) InstallationSLA(jobType valueobject.JobType, priority string) time.Duration {
    // INSTALL + CRITICAL = 24h, etc.
    return 48 * time.Hour
}
```

### 5.5 Domain Errors

```go
var (
    ErrShipmentNotFound       = errors.New("shipment not found")
    ErrInvalidShipmentNo      = errors.New("invalid shipment no")
    ErrShipmentLocked         = errors.New("shipment locked")
    ErrEmptyShipment          = errors.New("empty shipment")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrInstallationNotFound   = errors.New("installation not found")
    ErrNoTechnician           = errors.New("no technician assigned")
    ErrChecklistIncomplete    = errors.New("checklist incomplete")
    ErrMaintenanceNotFound    = errors.New("maintenance not found")
    ErrRMANotFound            = errors.New("rma not found")
    ErrInvalidGPS             = errors.New("invalid gps")
    ErrTechnicianNotFound     = errors.New("technician not found")
)
```

### 5.6 Invariants

1. Shipment number unique ต่อ tenant
2. Shipment ที่ dispatch แล้วแก้ items ไม่ได้
3. Installation job ต้อง assign technician ก่อน start
4. Complete ต้องผ่าน checklist ทั้งหมด
5. RMA ต้องมี device_id + customer_id
6. Maintenance schedule ต้องมี `next_due_at` > `last_done_at`
7. Technician GPS update ต้อง valid
8. Route optimization ต้องไม่เกิน 8 ชั่วโมง/วัน

## 6. Application Layer

### 6.1 Use Cases

```go
// application/complete_installation.go — orchestrate ข้ามโมดูลผ่าน Kafka
type CompleteInstallationUseCase struct {
    jobRepo   repository.InstallationRepository
    shipRepo  repository.ShipmentRepository
    producer  messaging.Producer
    tx        *transaction.Manager
}

type CompleteInstallationInput struct {
    TenantID     uuid.UUID
    JobID        uuid.UUID
    Checklist    []application.ChecklistItemDTO
    Photos       []string
    SignatureURL string
    ActorID      uuid.UUID
}

func (uc *CompleteInstallationUseCase) Execute(ctx context.Context, in CompleteInstallationInput) error {
    return uc.tx.Do(ctx, func(tx *gorm.DB) error {
        job, err := uc.jobRepo.FindByID(ctx, in.TenantID, in.JobID)
        if err != nil { return err }

        if err := job.Complete(toChecklist(in.Checklist), in.Photos, in.SignatureURL); err != nil {
            return err
        }
        if err := uc.jobRepo.Save(ctx, job); err != nil { return err }

        // Update shipment → INSTALLED
        if job.ShipmentID != nil {
            ship, err := uc.shipRepo.FindByID(ctx, in.TenantID, *job.ShipmentID)
            if err == nil {
                ship.Status = valueobject.ShipmentStatusInstalled
                _ = uc.shipRepo.Save(ctx, ship)
            }
        }

        // Publish event → iotdevice จะ provision devices
        _ = uc.producer.Publish(ctx, "iotlogistics.installation.completed", job.ID.String(), map[string]any{
            "job_id":        job.ID.String(),
            "site_id":       job.SiteID.String(),
            "shipment_id":   safeUUID(job.ShipmentID),
            "completed_at":  time.Now(),
        })

        return nil
    })
}
```

```go
// application/optimize_route.go — AI route planning
type OptimizeRouteUseCase struct {
    techRepo    repository.TechnicianRepository
    jobRepo     repository.InstallationRepository
    maps        service.MapsClient
    llm         llm.Client
}

type OptimizeRouteInput struct {
    TenantID     uuid.UUID
    TechnicianID uuid.UUID
    Date         time.Time
}

func (uc *OptimizeRouteUseCase) Execute(ctx context.Context, in OptimizeRouteInput) (*RouteResponse, error) {
    tech, err := uc.techRepo.FindByID(ctx, in.TenantID, in.TechnicianID)
    if err != nil { return nil, err }

    jobs, err := uc.jobRepo.ListByTechnician(ctx, in.TenantID, in.TechnicianID, in.Date)
    if err != nil { return nil, err }

    startGPS := tech.GPSLast

    // ใช้ Maps API คำนวณ distance matrix
    matrix, err := uc.maps.DistanceMatrix(ctx, startGPS, extractGPS(jobs))
    if err != nil { return nil, err }

    // ส่งให้ LLM ช่วยเลือกเส้นทาง (หรือใช้ TSP algorithm ใน Go)
    route, err := solveTSP(matrix, startGPS, jobs)
    if err != nil { return nil, err }

    return toRouteResponse(route), nil
}
```

### 6.2 DTOs

```go
type ShipmentResponse struct {
    ID          string     `json:"id"`
    ShipmentNo  string     `json:"shipment_no"`
    CustomerID  string     `json:"customer_id"`
    SiteID      string     `json:"site_id"`
    Carrier     string     `json:"carrier"`
    TrackingNo  string     `json:"tracking_no"`
    Status      string     `json:"status"`
    ScheduledAt time.Time  `json:"scheduled_at"`
    DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

type InstallationJobResponse struct {
    ID           string     `json:"id"`
    SiteID       string     `json:"site_id"`
    TechnicianID *string    `json:"technician_id,omitempty"`
    JobType      string     `json:"job_type"`
    Status       string     `json:"status"`
    ScheduledAt  time.Time  `json:"scheduled_at"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

type RouteResponse struct {
    Waypoints []GPSPoint `json:"waypoints"`
    Distance  float64    `json:"distance_km"`
    Duration  int        `json:"duration_min"`
    JobOrder  []string   `json:"job_order"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type ShipmentModel struct {
    ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
    TenantID        uuid.UUID  `gorm:"type:uuid;not null;index"`
    ShipmentNo      string     `gorm:"size:30;not null"`
    CustomerID      uuid.UUID  `gorm:"type:uuid;not null;index"`
    SiteID          uuid.UUID  `gorm:"type:uuid;not null;index"`
    FromWarehouseID uuid.UUID  `gorm:"type:uuid"`
    Carrier         string     `gorm:"size:100"`
    TrackingNo      string     `gorm:"size:100"`
    Status          string     `gorm:"size:30;index"`
    ScheduledAt     time.Time
    DeliveredAt     *time.Time
    GPSLast         datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
func (ShipmentModel) TableName() string { return "iotlogistics_shipments" }
```

### 7.2 Kafka Consumers

```go
// consumers/customer_onboarded_consumer.go — สร้าง shipment draft
type CustomerOnboardedConsumer struct{ createShipmentUC *application.CreateShipmentUseCase }

// consumers/erp_po_received_consumer.go — สร้าง shipment จาก PO
// consumers/device_offline_consumer.go — สร้าง work order ถ้า offline นาน
```

### 7.3 WebSocket Broadcaster
- Technician GPS tracking (realtime)
- Job status update → dashboard
- Shipment status change → customer

### 7.4 JWT, Bcrypt, Rate Limit
- Role: `technician`, `dispatcher`, `manager`
- Technician mobile app — JWT + refresh token

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *InstallationHandler) Complete(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid := c.MustGet("user_id").(uuid.UUID)
    id, _ := uuid.Parse(c.Param("id"))

    var req application.CompleteInstallationInput
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }
    req.TenantID = tid
    req.JobID = id
    req.ActorID = uid

    if err := h.completeUC.Execute(c.Request.Context(), req); err != nil {
        writeDomainError(c, err); return
    }
    c.JSON(200, gin.H{"message": "completed"})
}
```

### 8.2 Routes

```go
sh := r.Group("/shipments"); sh.Use(auth, tenant)
sh.POST("",             h.Shipment.Create)
sh.GET("",              h.Shipment.List)
sh.GET("/:id",          h.Shipment.Get)
sh.POST("/:id/dispatch",h.Shipment.Dispatch)
sh.POST("/:id/deliver", h.Shipment.Deliver)
sh.POST("/:id/gps",     h.Shipment.UpdateGPS)

ij := r.Group("/installation-jobs"); ij.Use(auth, tenant)
ij.POST("",                h.Installation.Schedule)
ij.GET("",                 h.Installation.List)
ij.POST("/:id/assign",     h.Installation.Assign)
ij.POST("/:id/start",      h.Installation.Start)
ij.POST("/:id/complete",   h.Installation.Complete)

mt := r.Group("/maintenance"); mt.Use(auth, tenant)
mt.POST("/schedule", h.Maintenance.Schedule)
mt.GET("/due",       h.Maintenance.Due)

rma := r.Group("/rma"); rma.Use(auth, tenant)
rma.POST("",             h.RMA.Create)
rma.POST("/:id/ship",    h.RMA.Ship)
rma.POST("/:id/receive", h.RMA.Receive)
rma.POST("/:id/close",   h.RMA.Close)

rt := r.Group("/routes"); rt.Use(auth, tenant)
rt.POST("/optimize", h.Route.Optimize)
```

### 8.3 Middleware
- `role.Require("technician")` สำหรับ start/complete
- Geo-fence check: ตรวจว่าช่างอยู่ในพื้นที่ก่อน start

## 9. Database Migrations

```sql
-- migrations/20260106_iotlogistics_init.sql
CREATE TABLE IF NOT EXISTS iotlogistics_shipments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    shipment_no       VARCHAR(30) NOT NULL,
    customer_id       UUID NOT NULL,
    site_id           UUID NOT NULL,
    from_warehouse_id UUID,
    carrier           VARCHAR(100),
    tracking_no       VARCHAR(100),
    status            VARCHAR(30) DEFAULT 'PENDING',
    scheduled_at      TIMESTAMP,
    delivered_at      TIMESTAMP,
    gps_last          JSONB,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, shipment_no)
);
CREATE INDEX idx_iotlogistics_shipments_status ON iotlogistics_shipments(tenant_id, status);

CREATE TABLE IF NOT EXISTS iotlogistics_shipment_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL REFERENCES iotlogistics_shipments(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL,
    device_id   UUID,
    qty         NUMERIC(15,2),
    serials     TEXT[]
);

CREATE TABLE IF NOT EXISTS iotlogistics_technicians (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    user_id      UUID NOT NULL,
    skill_set    TEXT[],
    zone         VARCHAR(100),
    gps_last     JSONB,
    is_available BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_iotlogistics_tech_tenant_zone ON iotlogistics_technicians(tenant_id, zone);

CREATE TABLE IF NOT EXISTS iotlogistics_installation_jobs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    shipment_id   UUID,
    site_id       UUID NOT NULL,
    technician_id UUID,
    job_type      VARCHAR(30),
    status        VARCHAR(20) DEFAULT 'SCHEDULED',
    scheduled_at  TIMESTAMP,
    started_at    TIMESTAMP,
    completed_at  TIMESTAMP,
    checklist     JSONB,
    photos        TEXT[],
    signature_url TEXT,
    notes         TEXT,
    created_at    TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_iotlogistics_jobs_tech_date ON iotlogistics_installation_jobs(technician_id, scheduled_at);
CREATE INDEX idx_iotlogistics_jobs_status    ON iotlogistics_installation_jobs(tenant_id, status);

CREATE TABLE IF NOT EXISTS iotlogistics_maintenance_schedules (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL,
    device_id     UUID NOT NULL,
    pm_type       VARCHAR(30),
    next_due_at   TIMESTAMP,
    last_done_at  TIMESTAMP,
    assigned_to   UUID,
    created_at    TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_iotlogistics_pm_due ON iotlogistics_maintenance_schedules(next_due_at);

CREATE TABLE IF NOT EXISTS iotlogistics_work_orders (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    wo_no        VARCHAR(30) NOT NULL,
    device_id    UUID,
    site_id      UUID NOT NULL,
    job_type     VARCHAR(30),
    priority     VARCHAR(20),
    status       VARCHAR(20) DEFAULT 'SCHEDULED',
    assigned_to  UUID,
    description  TEXT,
    scheduled_at TIMESTAMP,
    completed_at TIMESTAMP,
    parts_used   JSONB,
    created_at   TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, wo_no)
);

CREATE TABLE IF NOT EXISTS iotlogistics_rma (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL,
    rma_no                 VARCHAR(30) NOT NULL,
    device_id              UUID NOT NULL,
    customer_id            UUID NOT NULL,
    reason                 TEXT,
    status                 VARCHAR(20) DEFAULT 'OPEN',
    replacement_device_id  UUID,
    shipped_at             TIMESTAMP,
    received_at            TIMESTAMP,
    closed_at              TIMESTAMP,
    created_at             TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, rma_no)
);
```

## 10. System Flow

```
Installation Flow:
  Customer onboarded → Kafka: customer.onboarded
    → CreateShipmentUseCase (draft)
    → Dispatch (carrier, tracking)
    → Delivered at site
    → ScheduleInstallationJob
    → AssignTechnician (nearest + skill)
    → Technician starts (geo-fence check)
    → Complete (checklist + photos + signature)
    → Kafka: iotlogistics.installation.completed
      → iotdevice.provision(serial_no list)
      → package.start_billing

Maintenance Flow:
  Scheduler: pm_due_job (daily)
    → FindDue
    → CreateWorkOrder + AssignTechnician
    → Kafka: iotlogistics.maintenance.due
      → notifier: email/LINE
      → crm: create ticket
```

## 11. Workflow Diagram

```
[Customer Onboarded]
        │
        ▼
[Create Shipment Draft]──add items──▶[Dispatch]──▶[In Transit]
                                                        │
                                              GPS tracking (WS)
                                                        │
                                                        ▼
                                                  [Delivered]
                                                        │
                                                        ▼
                                            [Schedule Installation]
                                                        │
                                                        ▼
                                            [Assign Technician]
                                                        │
                                                        ▼
                                                  [Start Job]
                                                        │
                                                        ▼
                                                [Complete Job]
                                                        │
                                    ┌───────────────────┼───────────────┐
                                    ▼                   ▼               ▼
                            [iotdevice.provision]  [start billing]  [close shipment]
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=iotlogistics

# Env
GOOGLE_MAPS_API_KEY=***
TECHNICIAN_MAX_JOBS_PER_DAY=6
INSTALLATION_SLA_HOURS=48
PM_REMINDER_DAYS_BEFORE=7
GEOFENCE_RADIUS_METERS=500

# Run
go run ./cmd/api
go run ./cmd/scheduler   # pm_due_job, sla_check
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue | ค่าติดตั้ง, ค่าบำรุงรักษารายปี, ค่า RMA |
| Cost | ค่าช่าง, ค่าเดินทาง, ค่าอะไหล่ |
| KPI | Jobs/day/tech, First-time-fix rate, MTTR, Route efficiency |

## 14. ภาคผนวก

- **Route optimization**: TSP + Google Maps Distance Matrix
- **Geo-fence**: ตรวจ GPS ช่างก่อน start job
- **Photos**: อัปโหลดไป S3/MinIO, เก็บ URL
- **Signature**: เก็บเป็น base64 หรือ URL

## 15. Prompt ขยายระบบ

```
เพิ่ม use case AutoScheduleMaintenance ใน module iotlogistics
- Input: tenant_id
- Business: ดู device ทั้งหมดที่ติดตั้ง ≥ 6 เดือน → สร้าง PM schedule
- Side effects: Kafka iotlogistics.maintenance.scheduled, notifier
```

## 16. Prefix ตาราง

| Prefix | ตาราง |
|---|---|
| `iotlogistics_` | shipments, shipment_items, technicians, installation_jobs, maintenance_schedules, work_orders, rma |

## 17. DDD Validation Checklist

- [x] Shipment/InstallationJob/RMA เป็น Aggregate Root
- [x] Technician เป็น Entity แยก
- [x] ตรวจ geo-fence ก่อน start job
- [x] Complete installation → publish event (ไม่ import iotdevice ตรง)
- [x] Route optimization ผ่าน outbound port (maps + llm)
- [x] PM schedule → scheduler job
- [x] Multi-tenant + role-based access

---

# PART 7 — ระบบรายงาน (`report`)

## 1. ภาพรวมระบบ

ระบบรายงานและ Analytics: Report Definition (SQL/template), Scheduled Report (cron → email/LINE), KPI Snapshot (รายวัน), Dashboard API, Export (PDF/XLSX/CSV), Ad-hoc Query Builder + AI Insights (LLM สรุปแนวโน้ม)

## 2. โครงสร้าง Module

```
internal/modules/report/
├── domain/
│   ├── entity/
│   │   ├── report_definition.go
│   │   ├── report_schedule.go
│   │   ├── report_execution.go
│   │   ├── kpi_snapshot.go
│   │   └── dashboard.go
│   ├── value_object/
│   │   ├── report_type.go
│   │   ├── export_format.go
│   │   ├── kpi_code.go
│   │   └── date_range.go
│   ├── repository/
│   │   ├── definition_repository.go
│   │   ├── schedule_repository.go
│   │   ├── execution_repository.go
│   │   └── kpi_repository.go
│   ├── service/
│   │   ├── report_engine.go
│   │   ├── kpi_calculator.go
│   │   └── insight_generator.go     # AI
│   └── errors/errors.go
├── application/
│   ├── create_definition.go
│   ├── generate_report.go
│   ├── schedule_report.go
│   ├── export_report.go
│   ├── get_dashboard_kpi.go
│   ├── snapshot_kpi.go
│   ├── generate_ai_insight.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/
│   ├── persistence/clickhouse/       # OLAP (optional)
│   ├── renderer/pdf_renderer.go
│   ├── renderer/xlsx_renderer.go
│   ├── renderer/csv_renderer.go
│   ├── templates/                    # Go templates
│   │   ├── customer_360.html
│   │   ├── revenue_by_package.html
│   │   ├── device_uptime.html
│   │   └── ...
│   ├── messaging/kafka/report_producer.go
│   └── scheduler/report_scheduler_job.go
├── interfaces/http/
│   ├── definition_handler.go
│   ├── dashboard_handler.go
│   ├── export_handler.go
│   └── routes.go
└── module.go
```

## 3. Bounded Contexts และ Context Map

```
┌──────────┐  ┌──────────┐  ┌──────┐  ┌──────┐  ┌──────────┐
│ Customer │  │ Package  │  │ ERP  │  │ IoT  │  │Logistics │
└────┬─────┘  └────┬─────┘  └──┬───┘  └──┬───┘  └────┬─────┘
     │             │            │         │           │
     │ (events via Kafka)       │         │           │
     └─────────────┴────────────┴─────────┴───────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Report Context  │
                    │  [GENERIC]       │
                    │                  │
                    │  KPI Snapshot    │
                    │  Report Engine   │
                    │  Dashboard API   │
                    │  AI Insights     │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ Email / Line / WS│
                    └──────────────────┘
```

**Pattern**: Report Context เป็น **Conformist** ต่อทุก upstream — เก็บ read model ของตัวเองจาก Kafka events

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย |
|---|---|
| Report Definition | นิยามรายงาน (query + template) |
| Report Schedule | ตารางเวลาสร้างรายงาน |
| Report Execution | การรันครั้งหนึ่ง |
| KPI Snapshot | ค่า KPI ณ วันหนึ่ง |
| Dashboard | หน้าจอสรุป KPI |
| Export | ส่งออกเป็นไฟล์ |
| Insight | ข้อสรุปจาก AI |

## 5. Domain Layer

### 5.1 Entities

```go
// ReportDefinition — Aggregate Root
type ReportDefinition struct {
    ID           uuid.UUID
    TenantID     *uuid.UUID   // nil = shared
    Code         string
    Name         string
    Type         valueobject.ReportType
    QueryConfig  QueryConfig
    TemplatePath string
    DefaultFormat valueobject.ExportFormat
    IsActive     bool
    CreatedAt    time.Time
}

type QueryConfig struct {
    Source     string            // CUSTOMER | PACKAGE | ERP | IOT | LOGISTICS
    SQL        string            // parameterized
    Params     map[string]any
    GroupBy    []string
    OrderBy    []string
}

func NewReportDefinition(tenantID *uuid.UUID, code, name string, typ valueobject.ReportType) (*ReportDefinition, error) {
    if code == "" || name == "" { return nil, domainerrors.ErrInvalidDefinition }
    if !typ.IsValid() { return nil, domainerrors.ErrInvalidType }
    return &ReportDefinition{
        ID: uuid.New(), TenantID: tenantID, Code: code, Name: name,
        Type: typ, IsActive: true, CreatedAt: time.Now(),
    }, nil
}

// ReportSchedule — Aggregate Root
type ReportSchedule struct {
    ID           uuid.UUID
    TenantID     uuid.UUID
    ReportID     uuid.UUID
    CronExpr     string
    Recipients   []string
    Format       valueobject.ExportFormat
    Filters      map[string]any
    LastRunAt    *time.Time
    NextRunAt    time.Time
    CreatedAt    time.Time
}

func (s *ReportSchedule) MarkRun(at time.Time, next time.Time) {
    s.LastRunAt = &at
    s.NextRunAt = next
}

// ReportExecution — Entity
type ReportExecution struct {
    ID           uuid.UUID
    ScheduleID   *uuid.UUID
    ReportID     uuid.UUID
    TenantID     uuid.UUID
    Status       string   // PENDING|RUNNING|DONE|FAILED
    Format       valueobject.ExportFormat
    FileURL      string
    RowCount     int
    DurationMs   int
    ErrorMsg     string
    StartedAt    time.Time
    CompletedAt  *time.Time
}

// KPISnapshot — Entity
type KPISnapshot struct {
    ID         int64
    TenantID   uuid.UUID
    KPICode    valueobject.KPICode
    Value      float64
    Dimensions map[string]any
    SnapshotDate time.Time
}

// Dashboard — Aggregate Root
type Dashboard struct {
    ID       uuid.UUID
    TenantID uuid.UUID
    Name     string
    Layout   []Widget
    IsDefault bool
    OwnerID  uuid.UUID
    CreatedAt time.Time
}

type Widget struct {
    Type     string   // KPI_CARD | LINE_CHART | BAR_CHART | TABLE | PIE
    KPICodes []valueobject.KPICode
    Filters  map[string]any
    Position struct{ X, Y, W, H int }
}
```

### 5.2 Value Objects

```go
type ReportType string
const (
    ReportTypeFinancial   ReportType = "FINANCIAL"
    ReportTypeOperational ReportType = "OPERATIONAL"
    ReportTypeIoT         ReportType = "IOT"
    ReportTypeSales       ReportType = "SALES"
    ReportTypeCustomer    ReportType = "CUSTOMER"
)

type ExportFormat string
const (
    ExportFormatPDF  ExportFormat = "PDF"
    ExportFormatXLSX ExportFormat = "XLSX"
    ExportFormatCSV  ExportFormat = "CSV"
    ExportFormatJSON ExportFormat = "JSON"
)

type KPICode string
const (
    KPICodeMRR              KPICode = "MRR"                  // Monthly Recurring Revenue
    KPICodeARR              KPICode = "ARR"
    KPICodeActiveCustomers  KPICode = "ACTIVE_CUSTOMERS"
    KPICodeChurnRate        KPICode = "CHURN_RATE"
    KPICodeDeviceUptime     KPICode = "DEVICE_UPTIME_PCT"
    KPICodeAlertCount       KPICode = "ALERT_COUNT"
    KPICodeTicketSLA        KPICode = "TICKET_SLA_PCT"
    KPICodeInstallationSLA  KPICode = "INSTALLATION_SLA_PCT"
    KPICodeRevenueByPackage KPICode = "REVENUE_BY_PACKAGE"
    KPICodeUsageVariance    KPICode = "USAGE_VARIANCE_PCT"
)
```

### 5.3 Repository Interfaces

```go
type DefinitionRepository interface {
    Save(ctx context.Context, d *entity.ReportDefinition) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.ReportDefinition, error)
    FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.ReportDefinition, error)
    List(ctx context.Context, tenantID uuid.UUID) ([]*entity.ReportDefinition, error)
}

type ScheduleRepository interface {
    Save(ctx context.Context, s *entity.ReportSchedule) error
    FindDue(ctx context.Context, asOf time.Time) ([]*entity.ReportSchedule, error)
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.ReportSchedule, error)
}

type ExecutionRepository interface {
    Save(ctx context.Context, e *entity.ReportExecution) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.ReportExecution, error)
    ListByReport(ctx context.Context, tenantID, reportID uuid.UUID, limit int) ([]*entity.ReportExecution, error)
}

type KPIRepository interface {
    Upsert(ctx context.Context, s *entity.KPISnapshot) error
    Query(ctx context.Context, tenantID uuid.UUID, codes []valueobject.KPICode, from, to time.Time) ([]*entity.KPISnapshot, error)
    Latest(ctx context.Context, tenantID uuid.UUID, codes []valueobject.KPICode) ([]*entity.KPISnapshot, error)
}
```

### 5.4 Domain Services

```go
// ReportEngine — execute query + render
type ReportEngine struct {
    db       *gorm.DB
    renderer RendererRegistry
}

func (e *ReportEngine) Execute(ctx context.Context, def *entity.ReportDefinition, params map[string]any) (*ReportResult, error) {
    rows, err := e.db.WithContext(ctx).Raw(def.QueryConfig.SQL, toArgs(params)...).Rows()
    if err != nil { return nil, err }
    defer rows.Close()
    // ... scan rows
    return &ReportResult{Rows: /* ... */}, nil
}

// KPICalculator — คำนวณ KPI จาก raw data
type KPICalculator struct {
    db *gorm.DB
}

func (c *KPICalculator) CalculateMRR(ctx context.Context, tenantID uuid.UUID, asOf time.Time) (float64, error) {
    var total float64
    err := c.db.WithContext(ctx).Raw(`
        SELECT COALESCE(SUM(p.base_price), 0)
        FROM package_subscriptions s
        JOIN package_packages p ON p.id = s.package_id
        WHERE s.tenant_id = ? AND s.status = 'ACTIVE'
          AND s.start_date <= ? AND (s.end_date IS NULL OR s.end_date >= ?)
    `, tenantID, asOf, asOf).Scan(&total).Error
    return total, err
}

// InsightGenerator — AI สรุป
type InsightGenerator struct{ llm llm.Client }

func (g *InsightGenerator) Generate(ctx context.Context, kpis []*entity.KPISnapshot, period string) (string, error) {
    prompt := fmt.Sprintf(`You are a business analyst. Summarize these IoT platform KPIs for %s:
%s
Provide 3 bullet insights and 1 recommendation in Thai.`,
        period, formatKPIs(kpis))
    return g.llm.Generate(ctx, prompt)
}
```

### 5.5 Domain Errors

```go
var (
    ErrDefinitionNotFound = errors.New("report definition not found")
    ErrScheduleNotFound   = errors.New("report schedule not found")
    ErrExecutionNotFound  = errors.New("report execution not found")
    ErrInvalidDefinition  = errors.New("invalid definition")
    ErrInvalidType        = errors.New("invalid report type")
    ErrInvalidCron        = errors.New("invalid cron expression")
    ErrUnsupportedFormat  = errors.New("unsupported format")
    ErrQueryFailed        = errors.New("query failed")
)
```

### 5.6 Invariants

1. Report code unique ต่อ tenant
2. `QueryConfig.SQL` ต้อง parameterized (ห้าม concat string)
3. Cron expression ต้อง valid
4. Recipients ต้อง valid email/LINE ID
5. KPI snapshot unique ต่อ `(tenant, kpi_code, dimensions, date)`
6. Export file เก็บ 30 วัน
7. Dashboard ต้องมี widgets ≥ 1

## 6. Application Layer

### 6.1 Use Cases

```go
// application/generate_report.go
type GenerateReportUseCase struct {
    defRepo    repository.DefinitionRepository
    execRepo   repository.ExecutionRepository
    engine     *service.ReportEngine
    renderer   RendererRegistry
    producer   messaging.Producer
    storage    StoragePort
}

type GenerateReportInput struct {
    TenantID   uuid.UUID
    ReportCode string
    Format     valueobject.ExportFormat
    Params     map[string]any
    ActorID    uuid.UUID
}

func (uc *GenerateReportUseCase) Execute(ctx context.Context, in GenerateReportInput) (*ExecutionResponse, error) {
    def, err := uc.defRepo.FindByCode(ctx, in.TenantID, in.ReportCode)
    if err != nil { return nil, err }

    exec := &entity.ReportExecution{
        ID: uuid.New(), ReportID: def.ID, TenantID: in.TenantID,
        Status: "RUNNING", Format: in.Format, StartedAt: time.Now(),
    }
    _ = uc.execRepo.Save(ctx, exec)

    start := time.Now()
    result, err := uc.engine.Execute(ctx, def, in.Params)
    if err != nil {
        exec.Status = "FAILED"
        exec.ErrorMsg = err.Error()
        _ = uc.execRepo.Save(ctx, exec)
        return nil, domainerrors.ErrQueryFailed
    }

    // render
    var buf bytes.Buffer
    if err := uc.renderer.Render(in.Format, def.TemplatePath, result, &buf); err != nil {
        return nil, err
    }

    // upload
    url, err := uc.storage.Upload(ctx, fmt.Sprintf("reports/%s/%s.%s",
        in.TenantID, exec.ID, strings.ToLower(string(in.Format))), &buf)
    if err != nil { return nil, err }

    now := time.Now()
    exec.Status = "DONE"
    exec.FileURL = url
    exec.RowCount = len(result.Rows)
    exec.DurationMs = int(now.Sub(start).Milliseconds())
    exec.CompletedAt = &now
    _ = uc.execRepo.Save(ctx, exec)

    _ = uc.producer.Publish(ctx, "report.generated", exec.ID.String(), map[string]any{
        "execution_id": exec.ID.String(),
        "report_code":  def.Code,
        "file_url":     url,
        "tenant_id":    in.TenantID.String(),
    })
    return toExecutionResponse(exec), nil
}
```

```go
// application/snapshot_kpi.go — called by scheduler daily
type SnapshotKPIUseCase struct {
    calc    *service.KPICalculator
    kpiRepo repository.KPIRepository
}

func (uc *SnapshotKPIUseCase) Execute(ctx context.Context, tenantID uuid.UUID, asOf time.Time) error {
    codes := []valueobject.KPICode{
        valueobject.KPICodeMRR,
        valueobject.KPICodeActiveCustomers,
        valueobject.KPICodeDeviceUptime,
        valueobject.KPICodeAlertCount,
        valueobject.KPICodeTicketSLA,
        valueobject.KPICodeInstallationSLA,
    }
    for _, code := range codes {
        v, err := uc.calc.Calculate(ctx, tenantID, code, asOf)
        if err != nil { continue }
        _ = uc.kpiRepo.Upsert(ctx, &entity.KPISnapshot{
            TenantID: tenantID, KPICode: code,
            Value: v, SnapshotDate: asOf,
        })
    }
    return nil
}
```

### 6.2 DTOs

```go
type ExecutionResponse struct {
    ID          string     `json:"id"`
    ReportCode  string     `json:"report_code"`
    Status      string     `json:"status"`
    FileURL     string     `json:"file_url,omitempty"`
    RowCount    int        `json:"row_count"`
    DurationMs  int        `json:"duration_ms"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type DashboardResponse struct {
    KPIs      []KPIValue              `json:"kpis"`
    Charts    map[string]ChartData    `json:"charts"`
    GeneratedAt time.Time             `json:"generated_at"`
}

type KPIValue struct {
    Code  string  `json:"code"`
    Value float64 `json:"value"`
    Trend float64 `json:"trend_pct"`
    Unit  string  `json:"unit,omitempty"`
}

type InsightResponse struct {
    Period     string    `json:"period"`
    Summary    string    `json:"summary"`
    Insights   []string  `json:"insights"`
    Recommendation string `json:"recommendation"`
    GeneratedAt time.Time `json:"generated_at"`
}
```

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
type ReportDefinitionModel struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
    TenantID      *uuid.UUID     `gorm:"type:uuid;index"`
    Code          string         `gorm:"size:50;not null"`
    Name          string         `gorm:"size:255"`
    ReportType    string         `gorm:"size:30"`
    QueryConfig   datatypes.JSON `gorm:"type:jsonb"`
    TemplatePath  string         `gorm:"size:255"`
    DefaultFormat string         `gorm:"size:10"`
    IsActive      bool           `gorm:"default:true"`
    CreatedAt     time.Time
}
func (ReportDefinitionModel) TableName() string { return "report_definitions" }

type KPISnapshotModel struct {
    ID           int64          `gorm:"primaryKey;autoIncrement"`
    TenantID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    KPICode      string         `gorm:"size:50;not null;index"`
    Value        float64        `gorm:"type:numeric(20,4)"`
    Dimensions   datatypes.JSON `gorm:"type:jsonb"`
    SnapshotDate time.Time      `gorm:"type:date;index"`
}
func (KPISnapshotModel) TableName() string { return "report_kpi_snapshots" }
```

### 7.2 Kafka Consumers

```go
// consumers/kpi_trigger_consumer.go — trigger snapshot เมื่อมี event สำคัญ
type KPITriggerConsumer struct{ snapshotUC *application.SnapshotKPIUseCase }

func (c *KPITriggerConsumer) Handle(ctx context.Context, payload map[string]any) error {
    tid := uuid.MustParse(payload["tenant_id"].(string))
    return c.snapshotUC.Execute(ctx, tid, time.Now())
}
```

Subscribe: `package.subscription.created`, `erp.invoice.issued`, `iot.alert.triggered`, `crm.ticket.resolved`, `iotlogistics.installation.completed`

### 7.3 WebSocket Broadcaster
- KPI update realtime → dashboard
- Report generated → notify user

### 7.4 JWT, Bcrypt, Rate Limit
- Role: `executive`, `manager`, `analyst`
- Export — rate limit 10/min (heavy operation)

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
func (h *DashboardHandler) GetKPIs(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    codes := c.QueryArray("code")

    res, err := h.getKPIsUC.Execute(c.Request.Context(), application.GetKPIsInput{
        TenantID: tid,
        Codes:    toKPICodes(codes),
        From:     parseTimeOrDefault(c.Query("from"), time.Now().AddDate(0, -1, 0)),
        To:       parseTimeOrDefault(c.Query("to"), time.Now()),
    })
    if err != nil { writeDomainError(c, err); return }
    c.JSON(200, res)
}

func (h *ExportHandler) Generate(c *gin.Context) {
    tid := c.MustGet("tenant_id").(uuid.UUID)
    uid := c.MustGet("user_id").(uuid.UUID)

    var req struct {
        ReportCode string         `json:"report_code" binding:"required"`
        Format     string         `json:"format" binding:"required"`
        Params     map[string]any `json:"params"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()}); return
    }

    res, err := h.genUC.Execute(c.Request.Context(), application.GenerateReportInput{
        TenantID: tid, ReportCode: req.ReportCode,
        Format: valueobject.ExportFormat(req.Format),
        Params: req.Params, ActorID: uid,
    })
    if err != nil { writeDomainError(c, err); return }
    c.JSON(202, res)
}
```

### 8.2 Routes

```go
rep := r.Group("/reports"); rep.Use(auth, tenant, role.Require("analyst", "manager", "executive"))
rep.GET   ("",                    h.Def.List)
rep.POST  ("",                    h.Def.Create)
rep.GET   ("/:code",              h.Def.Get)
rep.POST  ("/:code/generate",     h.Export.Generate)
rep.GET   ("/executions/:id",     h.Export.Get)

sch := r.Group("/report-schedules"); sch.Use(auth, tenant)
sch.POST  ("",          h.Schedule.Create)
sch.GET   ("",          h.Schedule.List)
sch.POST  ("/:id/run",  h.Schedule.RunNow)

db := r.Group("/dashboard"); db.Use(auth, tenant)
db.GET("/kpis",       h.Dash.GetKPIs)
db.GET("/charts",     h.Dash.GetCharts)
db.GET("/insights",   h.Dash.GetInsights)   // AI
db.GET("/layout",     h.Dash.GetLayout)
db.PUT("/layout",     h.Dash.UpdateLayout)

ex := r.Group("/exports"); ex.Use(auth, tenant, rate_limit.Require(10, time.Minute))
ex.POST("",            h.Export.Generate)
ex.GET ("/:id",        h.Export.Get)
ex.GET ("/:id/download", h.Export.Download)
```

### 8.3 Middleware
- `role.Require` สำหรับ dashboard (executive)
- `rate_limit` เข้มงวดกับ export

## 9. Database Migrations

```sql
-- migrations/20260107_report_init.sql
CREATE TABLE IF NOT EXISTS report_definitions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID,
    code           VARCHAR(50) NOT NULL,
    name           VARCHAR(255) NOT NULL,
    report_type    VARCHAR(30),
    query_config   JSONB NOT NULL,
    template_path  VARCHAR(255),
    default_format VARCHAR(10) DEFAULT 'PDF',
    is_active      BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE IF NOT EXISTS report_schedules (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL,
    report_id    UUID NOT NULL REFERENCES report_definitions(id),
    cron_expr    VARCHAR(50) NOT NULL,
    recipients   TEXT[],
    format       VARCHAR(10),
    filters      JSONB,
    last_run_at  TIMESTAMP,
    next_run_at  TIMESTAMP,
    created_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_report_sched_next_run ON report_schedules(next_run_at);

CREATE TABLE IF NOT EXISTS report_executions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id  UUID,
    report_id    UUID NOT NULL,
    tenant_id    UUID NOT NULL,
    status       VARCHAR(20),
    format       VARCHAR(10),
    file_url     TEXT,
    row_count    INT,
    duration_ms  INT,
    error_msg    TEXT,
    started_at   TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);
CREATE INDEX idx_report_exec_tenant_started ON report_executions(tenant_id, started_at DESC);

CREATE TABLE IF NOT EXISTS report_kpi_snapshots (
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     UUID NOT NULL,
    kpi_code      VARCHAR(50) NOT NULL,
    value         NUMERIC(20,4) NOT NULL,
    dimensions    JSONB,
    snapshot_date DATE NOT NULL,
    UNIQUE(tenant_id, kpi_code, snapshot_date)
);
CREATE INDEX idx_report_kpi_tenant_code_date ON report_kpi_snapshots(tenant_id, kpi_code, snapshot_date DESC);

CREATE TABLE IF NOT EXISTS report_dashboards (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL,
    name       VARCHAR(255),
    layout     JSONB NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    owner_id   UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## 10. System Flow

```
Scheduled Report:
  Scheduler (cron) → FindDue
    → GenerateReportUseCase
      ├─ ReportEngine.Execute (SQL query)
      ├─ Render (PDF/XLSX)
      ├─ Upload to storage
      ├─ Save execution
      └─ Kafka: report.generated
        → notifier.sendEmail

Dashboard:
  GET /dashboard/kpis
    → GetKPIsUseCase → KPIRepository.Latest()
    → Cache in Redis (TTL 5m)

Daily Snapshot:
  Scheduler (daily 2 AM)
    → SnapshotKPIUseCase
      → KPICalculator.Calculate(...)
      → KPIRepository.Upsert()

AI Insights:
  GET /dashboard/insights?period=last_month
    → GetInsightsUseCase
      → InsightGenerator.Generate(kpis)
        → LLM → สรุปเป็น bullet + recommendation
```

## 11. Workflow Diagram

```
[Events from other modules]
        │
        ▼
[Kafka → KPI Trigger Consumer]
        │
        ▼
[SnapshotKPIUseCase]
        │
        ▼
[KPICalculator queries PG + InfluxDB]
        │
        ▼
[report_kpi_snapshots upsert]
        │
        ▼
[Dashboard API reads snapshots + Redis cache]
        │
        ▼
[Frontend renders charts + AI insights]

Report generation flow:
[Cron] → [GenerateReportUseCase]
            │
            ├─ SQL query (PG)
            ├─ Render (PDF/XLSX/CSV)
            ├─ Upload (S3/MinIO)
            ├─ Save execution
            └─ Kafka: report.generated
                → notifier → email/LINE
```

## 12. การติดตั้งและใช้งาน

```bash
make migrate-up MODULE=report

# Env
REPORT_STORAGE_URL=s3://reports
REPORT_RETENTION_DAYS=30
KPI_SNAPSHOT_CRON="0 2 * * *"        # daily 2 AM
REPORT_SCHEDULER_INTERVAL=60s
CLICKHOUSE_ENABLED=false
LLM_ENABLED=true

# Run
go run ./cmd/api
go run ./cmd/scheduler   # report_scheduler_job, kpi_snapshot_job
```

## 13. Business Model

| รายการ | รายละเอียด |
|---|---|
| Revenue | Free (Basic dashboard), Pro (custom reports), Enterprise (AI insights) |
| Value | Data-driven decisions, ลดเวลา report, insight ที่ actionable |
| KPI | Report adoption rate, Dashboard DAU, Time-to-insight |

## 14. ภาคผนวก

- **Reports ที่ต้องมี**: Customer 360, Revenue by Package, Device Uptime, Telemetry Analytics (Smart Farm/Building), Energy Analytics, Logistics SLA, Maintenance Due, Ticket SLA, Inventory Aging, Usage vs Quota
- **Storage**: S3/MinIO
- **OLAP**: ClickHouse (optional) สำหรับ query หนักๆ
- **AI**: `pkg/llm` — summary, anomaly, forecast

## 15. Prompt ขยายระบบ

```
เพิ่ม use case ScheduleRecurringReport ใน module report
- Input: report_code, cron_expr, recipients, format
- Business: ตรวจ cron + recipients → save schedule
- Side effects: scheduler job จะ pickup ตาม next_run_at
```

## 16. Prefix ตาราง

| Prefix | ตาราง |
|---|---|
| `report_` | definitions, schedules, executions, kpi_snapshots, dashboards |

## 17. DDD Validation Checklist

- [x] ReportDefinition/Schedule/Dashboard เป็น Aggregate Root แยก
- [x] KPISnapshot เป็น Entity (ไม่ mutate)
- [x] Query parameterized (ห้าม SQL injection)
- [x] Export file เก็บใน storage (ไม่ใช่ DB)
- [x] AI insight ผ่าน outbound port (llm.Client)
- [x] Caching KPI ด้วย Redis (TTL)
- [x] Kafka trigger สำหรับ snapshot
- [x] Multi-tenant isolation

---

# PART 8 — SYSTEM-WIDE CHECKLIST

## 8.1 Cross-Module Integration Checklist

| # | Integration | Producer | Consumer | Event | Status |
|---|---|---|---|---|---|
| 1 | CRM → Customer | crm | customer | `crm.lead.converted` | ✅ |
| 2 | Customer → Package | customer | packagecatalog | `customer.onboarded` | ✅ |
| 3 | Package → IoTDevice | packagecatalog | iotdevice | `package.subscription.created` | ✅ |
| 4 | Package → Payment | packagecatalog | payment | `package.subscription.created` | ✅ |
| 5 | IoTDevice → CRM | iotdevice | crm | `iot.alert.triggered` | ✅ |
| 6 | IoTDevice → Package | iotdevice | packagecatalog | `iot.telemetry.aggregated` | ✅ |
| 7 | IoTDevice → Report | iotdevice | report | `iot.alert.triggered` | ✅ |
| 8 | Logistics → IoTDevice | iotlogistics | iotdevice | `iotlogistics.installation.completed` | ✅ |
| 9 | Logistics → Package | iotlogistics | packagecatalog | `iotlogistics.installation.completed` | ✅ |
| 10 | ERP → Payment | erp | payment | `erp.invoice.issued` | ✅ |
| 11 | ERP → Report | erp | report | `erp.payment.posted` | ✅ |
| 12 | Report → Notifier | report | notifier | `report.generated` | ✅ |

## 8.2 Multi-Tenant Rules (ทุกโมดูล)

- [x] ทุกตารางมี `tenant_id UUID NOT NULL`
- [x] Index `(tenant_id, ...)` ทุกตาราง
- [x] Repository query กรองด้วย `tenant_id`
- [x] Middleware `tenant.Inject()`
- [x] Redis key: `{tenant}:{module}:{entity}:{id}`
- [x] ES index: `{module}_{tenant_id}_*`
- [x] MQTT topic: `iot/{tenant_id}/{device_id}/{action}`
- [x] S3 path: `{module}/{tenant_id}/{...}`
- [x] Kafka message มี `tenant_id`

## 8.3 Clean Architecture Compliance

| Layer | ห้าม import | ตรวจสอบ |
|---|---|---|
| Domain | gorm, gin, sarama, redis, mqtt | ✅ |
| Application | gorm, gin | ✅ |
| Infrastructure | gin | ✅ |
| Interface | — | ✅ |

## 8.4 Package Whitelist (ทุกโมดูล)

ใช้ได้เฉพาะ: `pkg/{cryptpass, db, elasticsearch, emailTemplates, helpers, http-swagger, httpErrors, influxdb, jwt, kafka, llm, logger, mqtt, report, responses, secureRandom, sendEmail, transaction, utils, vectordb, websocket}`

## 8.5 Environment Variables (รวม)

```env
# Core
DB_DSN=postgres://...
REDIS_ADDR=redis:6379
KAFKA_BROKERS=kafka:9092
JWT_SECRET=***
JWT_TTL=24h

# MQTT
MQTT_BROKER=tcp://mosquitto:1883
MQTT_USERNAME=iot
MQTT_PASSWORD=***

# InfluxDB
INFLUX_URL=http://influxdb:8086
INFLUX_TOKEN=***
INFLUX_ORG=icmon
INFLUX_BUCKET=iot_telemetry

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200

# Storage
S3_ENDPOINT=minio:9000
S3_BUCKET=iot-platform

# LLM
LLM_ENABLED=true
LLM_PROVIDER=ollama
LLM_API_URL=http://ollama:11434
LLM_MODEL=llama3

# Maps
GOOGLE_MAPS_API_KEY=***

# Module-specific
CUSTOMER_DEFAULT_CURRENCY=THB
PACKAGE_QUOTA_RESET_CRON="0 0 1 * *"
IOT_OFFLINE_THRESHOLD_MINUTES=5
IOT_COMMAND_TIMEOUT_SECONDS=30
LOGISTICS_INSTALLATION_SLA_HOURS=48
REPORT_KPI_SNAPSHOT_CRON="0 2 * * *"
```

## 8.6 Bootstrap — `cmd/api/main.go`

```go
func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DB)
    rdb := redis.Connect(cfg.Redis)
    kp := kafka.NewProducer(cfg.Kafka.Brokers)
    mqttClient := mqtt.NewClient(cfg.MQTT)
    influx := influxdb.NewClient(cfg.Influx)
    es := elasticsearch.NewClient(cfg.ES)
    llmClient := llm.New(cfg.LLM)
    wsHub := websocket.NewHub()

    go wsHub.Run()

    r := gin.New()
    r.Use(gin.Recovery(), middleware.Logging(), middleware.CORS())

    api := r.Group("/api/v1")
    auth := middleware.Auth(cfg.JWT.Secret)
    tenant := middleware.Inject()

    // Wire all modules
    customer.Init(api, customer.Deps{DB: db, Producer: kp}, auth, tenant)
    packagecatalog.Init(api, packagecatalog.Deps{DB: db, Redis: rdb, Producer: kp}, auth, tenant)
    erp.Init(api, erp.Deps{DB: db, Producer: kp}, auth, tenant)
    crm.Init(api, crm.Deps{DB: db, Producer: kp, LLM: llmClient}, auth, tenant)
    iotdevice.Init(api, iotdevice.Deps{
        DB: db, Redis: rdb, Influx: influx, ES: es,
        MQTT: mqttClient, WS: wsHub, LLM: llmClient, Producer: kp,
    }, auth, tenant)
    iotlogistics.Init(api, iotlogistics.Deps{DB: db, Producer: kp}, auth, tenant)
    report.Init(api, report.Deps{DB: db, Redis: rdb, LLM: llmClient, Producer: kp}, auth, tenant)

    r.Run(":" + cfg.Port)
}
```

---

# PART 9 — FINAL SUMMARY

| # | Module | Tables | Aggregates | Kafka Topics | Storage |
|---|---|---|---|---|---|
| 1 | `customer` | 4 | Customer, Contract | 5 | PG, Redis, ES |
| 2 | `packagecatalog` | 4 | Package, Subscription | 5 | PG, Redis |
| 3 | `erp` | 12 | Inventory, PO, SO, Invoice | 8 | PG, S3 (PDF) |
| 4 | `crm` | 5 | Lead, Opportunity, Ticket | 6 | PG, ES, LLM |
| 5 | `iotdevice` | 7 | Device, Command, AlertRule, Automation | 10 | PG, Influx, Redis, MQTT |
| 6 | `iotlogistics` | 7 | Shipment, Installation, RMA | 6 | PG, S3 (photos) |
| 7 | `report` | 5 | ReportDef, Schedule, Dashboard | 3 | PG, S3, Redis |

**รวม:** 44 ตาราง PostgreSQL, 43 Kafka Topics, 7 Bounded Contexts, 100% Clean Architecture + DDD ตาม `Template_Module.md`

---

 