# 📚 PART 5, 6, 7 — CORE LAYERS DEEP DIVE

> **ขนาด**: ใหญ่มาก — แยกเป็น 3 PART ใหญ่
> **PART 5**: Domain Layer Deep Dive
> **PART 6**: Application Layer Deep Dive
> **PART 7**: Infrastructure Layer Deep Dive
> **เป้าหมาย**: ให้ dev เข้าใจและ implement ทั้ง 3 layer ได้ทันที ตาม Clean Architecture + DDD

---
---

# 🏛️ PART 5 — DOMAIN LAYER DEEP DIVE

> **ขนาด**: ใหญ่ — แยก 9 ตอนย่อย
> **Part 5A**: Domain Layer Architecture & Principles
> **Part 5B**: Aggregate Design Patterns
> **Part 5C**: All Aggregates by Context (7 contexts)
> **Part 5D**: Value Objects (complete catalog)
> **Part 5E**: Domain Events
> **Part 5F**: Domain Services & Specifications
> **Part 5G**: Repository Interfaces
> **Part 5H**: Factories, Invariants & Guards
> **Part 5I**: Domain Errors (complete catalog)

---

## 🅰️ PART 5A — DOMAIN LAYER ARCHITECTURE & PRINCIPLES

### A.1 Domain Layer Position

```
┌─────────────────────────────────────────────────────────────┐
│                    Interface Layer                          │
│  HTTP · WebSocket · MQTT · CLI · Scheduler                  │
├─────────────────────────────────────────────────────────────┤
│                   Application Layer                         │
│  Use Cases · DTOs · Orchestration                           │
├─────────────────────────────────────────────────────────────┤
│  ★★★ DOMAIN LAYER ★★★  ← PART 5 นี้                        │
│  Entities · VOs · Events · Services · Repos (interface)     │
│  ⚠️  ห้าม import: gorm, gin, sarama, redis, mqtt, influxdb  │
├─────────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                       │
│  Postgres · Redis · Kafka · InfluxDB · ES · MQTT · LLM      │
└─────────────────────────────────────────────────────────────┘
```

### A.2 Domain Layer Rules

| # | กฎ | เหตุผล |
|:-:|:---|:---|
| 1 | Entity มี constructor `New{{Entity}}` | บังคับ invariant ตอนสร้าง |
| 2 | Entity ไม่มี public setter | เปลี่ยน state ผ่าน behavior method |
| 3 | Value Object ทุกตัวมี `IsValid()` | Validation สม่ำเสมอ |
| 4 | Repository เป็น interface เท่านั้น | ไม่ผูกกับ infrastructure |
| 5 | Domain error เป็น sentinel error | เทียบด้วย `errors.Is()` |
| 6 | ❌ ห้าม import ORM/HTTP/MQTT | Domain ต้อง pure |
| 7 | Invariants บังคับใน constructor + behavior | ป้องกัน invalid state |
| 8 | 1 Aggregate Root / 1 transaction | Consistency boundary |
| 9 | Reference ข้าม aggregate ใช้ ID | ไม่ใช้ object reference |
| 10 | Domain event เกิดจาก behavior method | Track state change |

### A.3 Allowed Imports

```go
// ✅ อนุญาตใน Domain Layer
import (
    "context"                    // ✅ stdlib
    "errors"                     // ✅ stdlib
    "time"                       // ✅ stdlib
    "strings"                    // ✅ stdlib
    "regexp"                     // ✅ stdlib
    "github.com/google/uuid"    // ✅ uuid
    "github.com/shopspring/decimal" // ✅ decimal (สำหรับ Money)
)

// ❌ ห้าม
import (
    "gorm.io/gorm"              // ❌ ORM
    "github.com/gin-gonic/gin"  // ❌ HTTP
    "github.com/IBM/sarama"     // ❌ Kafka
    "github.com/redis/go-redis" // ❌ Redis
    "github.com/eclipse/paho.mqtt.golang" // ❌ MQTT
    "github.com/influxdata/influxdb-client-go" // ❌ Influx
)
```

### A.4 Domain Layer Folder Structure

```
internal/modules/{{module}}/domain/
├── entity/                    # Aggregate Roots + Child Entities
│   ├── {{aggregate}}.go
│   ├── {{child_entity}}.go
│   └── {{aggregate}}_test.go
├── value_object/              # Immutable VOs
│   ├── status.go
│   ├── {{vo_name}}.go
│   ├── id.go
│   └── *_test.go
├── event/                     # Domain Events
│   └── {{aggregate}}_events.go
├── repository/                # Interfaces เท่านั้น
│   ├── {{aggregate}}_repository.go
│   ├── read_model_repository.go
│   └── *_test.go
├── service/                   # Stateless domain logic
│   ├── {{domain}}_service.go
│   ├── {{port}}_port.go       # Outbound ports
│   └── *_test.go
├── specification/             # Business rules
│   └── {{name}}_spec.go
└── errors/                    # Sentinel errors
    └── errors.go
```

---

## 🅱️ PART 5B — AGGREGATE DESIGN PATTERNS

### B.1 Anatomy of an Aggregate

```
┌─────────────────────────────────────────────────────────┐
│              AGGREGATE ROOT (Device)                    │
│  ┌──────────────────────────────────────────────────┐   │
│  │  ID: uuid.UUID                                   │   │
│  │  TenantID, SiteID, ZoneID  ← Reference by ID    │   │
│  │  SerialNumber (VO)                               │   │
│  │  Status (VO)                                     │   │
│  │                                                  │   │
│  │  ── Behavior Methods ──                          │   │
│  │  MarkOnline()                                    │   │
│  │  MarkOffline()                                   │   │
│  │  SetError(msg)                                   │   │
│  │  UpdateFirmware(v) error                         │   │
│  │  MoveTo(siteID, zoneID) error                    │   │
│  │                                                  │   │
│  │  ── Domain Events (internal) ──                  │   │
│  │  events: []DomainEvent                           │   │
│  │  PullEvents() []DomainEvent                      │   │
│  └──────────────────────────────────────────────────┘   │
│           │                                             │
│           │ owns (1-to-many)                            │
│           ▼                                             │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  Command (child) │  │  Sensor (child)  │            │
│  │  - ID            │  │  - ID            │            │
│  │  - Type          │  │  - Metric        │            │
│  │  - Status        │  │  - Unit          │            │
│  └──────────────────┘  └──────────────────┘            │
└─────────────────────────────────────────────────────────┘
```

### B.2 Aggregate Boundary Rules

| กฎ | ตัวอย่าง |
|:---|:---|
| **1 Aggregate Root / Transaction** | `Device` + `Command` = 1 transaction |
| **Reference ข้าม aggregate ใช้ ID** | `Device.SiteID uuid.UUID` (ไม่ใช่ `*Site`) |
| **Consistency ภายใน aggregate** | Invariant ภายใน aggregate รับประกัน |
| **Consistency ข้าม aggregate** | ใช้ domain event + eventual consistency |
| **Aggregate เล็ก ดีกว่าใหญ่** | แยก `Telemetry` ออกจาก `Device` |
| **โหลดทั้ง aggregate** | Repository โหลด root + children |

### B.3 Aggregate Size Guidelines

| ขนาด | ตัวอย่าง | จำนวน children |
|:---|:---|:---|
| **เล็ก** | `Package`, `Lead` | 0-5 |
| **กลาง** | `Device`, `Order` | 5-50 |
| **ใหญ่** | `Customer` (มี sites, contacts) | 50-500 |
| **ใหญ่มาก** ❌ | ควรแยก | > 500 |

> **กฎ:** ถ้า aggregate มี children > 500 → ควรแยกเป็น aggregate ใหม่

### B.4 Aggregate Root Template

```go
// internal/modules/{{module}}/domain/entity/{{aggregate}}.go
package entity

import (
    "time"
    "github.com/google/uuid"
    domainerrors "icmongolang/internal/modules/{{module}}/domain/errors"
    domainevent "icmongolang/internal/modules/{{module}}/domain/event"
    vo "icmongolang/internal/modules/{{module}}/domain/value_object"
)

// {{Aggregate}} – Aggregate Root
type {{Aggregate}} struct {
    // Identity
    ID       uuid.UUID
    TenantID uuid.UUID

    // Attributes (Value Objects)
    Code   vo.{{Aggregate}}Code
    Name   string
    Status vo.{{Aggregate}}Status

    // References to other aggregates (by ID only)
    ParentID uuid.UUID

    // Child entities (owned)
    Children []{{ChildEntity}}

    // Audit
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy uuid.UUID

    // Domain events (unexported)
    events []domainevent.DomainEvent
}

// ============================================================
// CONSTRUCTOR — บังคับ invariant
// ============================================================
func New{{Aggregate}}(
    tenantID uuid.UUID,
    code vo.{{Aggregate}}Code,
    name string,
    actorID uuid.UUID,
) (*{{Aggregate}}, error) {
    // Validate
    if tenantID == uuid.Nil {
        return nil, domainerrors.ErrInvalidTenant
    }
    if !code.IsValid() {
        return nil, domainerrors.ErrInvalidCode
    }
    if strings.TrimSpace(name) == "" {
        return nil, domainerrors.ErrInvalidName
    }

    now := time.Now().UTC()
    a := &{{Aggregate}}{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Code:      code,
        Name:      strings.TrimSpace(name),
        Status:    vo.StatusDraft,
        CreatedAt: now,
        UpdatedAt: now,
        CreatedBy: actorID,
    }
    a.addEvent(domainevent.New{{Aggregate}}Created(a.ID, tenantID, actorID))
    return a, nil
}

// ============================================================
// BEHAVIOR METHODS — เปลี่ยน state ผ่าน method เท่านั้น
// ============================================================
func (a *{{Aggregate}}) Activate(actorID uuid.UUID) error {
    if a.Status == vo.StatusActive {
        return domainerrors.ErrAlreadyActive
    }
    if a.Status == vo.StatusCancelled {
        return domainerrors.ErrCannotActivateCancelled
    }
    a.Status = vo.StatusActive
    a.UpdatedAt = time.Now().UTC()
    a.addEvent(domainevent.New{{Aggregate}}Activated(a.ID, actorID))
    return nil
}

func (a *{{Aggregate}}) Update(name string, actorID uuid.UUID) error {
    if strings.TrimSpace(name) == "" {
        return domainerrors.ErrInvalidName
    }
    a.Name = strings.TrimSpace(name)
    a.UpdatedAt = time.Now().UTC()
    a.addEvent(domainevent.New{{Aggregate}}Updated(a.ID, actorID))
    return nil
}

// ============================================================
// DOMAIN EVENTS
// ============================================================
func (a *{{Aggregate}}) addEvent(e domainevent.DomainEvent) {
    a.events = append(a.events, e)
}

func (a *{{Aggregate}}) PullEvents() []domainevent.DomainEvent {
    evts := a.events
    a.events = nil
    return evts
}

func (a *{{Aggregate}}) HasEvents() bool {
    return len(a.events) > 0
}

// ============================================================
// RECONSTITUTION — สำหรับ repository โหลดจาก DB
// ============================================================
type {{Aggregate}}Snapshot struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    Code      string
    Name      string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func Reconstitute{{Aggregate}}(s {{Aggregate}}Snapshot) *{{Aggregate}} {
    return &{{Aggregate}}{
        ID:        s.ID,
        TenantID:  s.TenantID,
        Code:      vo.{{Aggregate}}Code(s.Code),
        Name:      s.Name,
        Status:    vo.{{Aggregate}}Status(s.Status),
        CreatedAt: s.CreatedAt,
        UpdatedAt: s.UpdatedAt,
    }
}
```

---

## 🅲 PART 5C — ALL AGGREGATES BY CONTEXT

### C.1 Aggregate Inventory

| Context | Aggregate Roots | Child Entities | VOs |
|:---|:---|:---|:---|
| **Customer** | Customer | Contact | CustomerCode, CustomerType, Segment, Lifecycle, Address, TaxID |
| **Package** | Package, Subscription | Feature, Quota | PackageCode, Tier, Cycle, Money, QuotaLimits |
| **Device** | Device, AlertRule | Command, Shadow | DeviceSerial, DeviceType, Protocol, DeviceStatus, CommandStatus |
| **ERP** | Order, Invoice, Payment, Product, Warehouse | OrderLine, InvoiceLine, Allocation | DocumentNumber, SKU, Money, TaxType, OrderStatus |
| **CRM** | Lead, Opportunity, Ticket | Activity | LeadSource, LeadStatus, TicketPriority, SLA |
| **Logistics** | Shipment, InstallationJob, Technician | ShipmentItem, ChecklistItem | TrackingNo, JobStatus, GPS, SLA |
| **Report** | ReportDefinition, Dashboard | Widget, Parameter | ReportCategory, KPI code |

### C.2 Customer Context — `Customer` Aggregate

```go
// internal/modules/customer/domain/entity/customer.go
package entity

type Customer struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    Code      vo.CustomerCode
    Name      string
    LegalName string
    Type      vo.CustomerType  // CORPORATE, INDIVIDUAL, GOVERNMENT
    Segment   vo.Segment       // ENTERPRISE, SME, STARTUP, INDIVIDUAL, GOVERNMENT
    Lifecycle vo.Lifecycle     // LEAD, PROSPECT, CUSTOMER, CHURNED
    Status    vo.CustomerStatus // PENDING, ACTIVE, SUSPENDED, CANCELLED

    // Business
    TaxID       *vo.TaxID
    Email       string
    Phone       string
    Address     vo.Address
    CreditLimit vo.Money

    // References
    PackageID uuid.UUID
    Sites     []uuid.UUID  // reference to Device context

    // Children
    Contacts []Contact

    // Audit
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy uuid.UUID

    events []domainevent.DomainEvent
}

// ---- Behavior Methods ----
func (c *Customer) Qualify() error {
    if c.Lifecycle != vo.LifecycleLead {
        return domainerrors.ErrNotALead
    }
    c.Lifecycle = vo.LifecycleProspect
    c.UpdatedAt = time.Now().UTC()
    c.addEvent(domainevent.NewCustomerQualified(c.ID, c.TenantID))
    return nil
}

func (c *Customer) Activate(actorID uuid.UUID) error {
    if c.Status == vo.CustomerStatusActive {
        return domainerrors.ErrCustomerAlreadyActive
    }
    if c.Status == vo.CustomerStatusCancelled {
        return domainerrors.ErrCannotReactivateCancelled
    }
    // Invariant: ต้องเป็น PROSPECT หรือ PENDING
    if c.Lifecycle == vo.LifecycleLead {
        return domainerrors.ErrMustQualifyFirst
    }
    c.Status = vo.CustomerStatusActive
    c.Lifecycle = vo.LifecycleCustomer
    c.UpdatedAt = time.Now().UTC()
    c.addEvent(domainevent.NewCustomerActivated(c.ID, c.TenantID, actorID))
    return nil
}

func (c *Customer) Suspend(reason string, actorID uuid.UUID) error {
    if c.Status == vo.CustomerStatusSuspended {
        return domainerrors.ErrAlreadySuspended
    }
    if c.Status == vo.CustomerStatusCancelled {
        return domainerrors.ErrCannotSuspendCancelled
    }
    c.Status = vo.CustomerStatusSuspended
    c.UpdatedAt = time.Now().UTC()
    c.addEvent(domainevent.NewCustomerSuspended(c.ID, c.TenantID, reason, actorID))
    return nil
}

func (c *Customer) SetTaxID(taxID vo.TaxID) error {
    // Invariant: CORPORATE + GOVERNMENT ต้องมี TaxID
    if c.Type == vo.CustomerTypeIndividual {
        return domainerrors.ErrIndividualNoTaxID
    }
    if !taxID.IsValid() {
        return domainerrors.ErrInvalidTaxID
    }
    c.TaxID = &taxID
    c.UpdatedAt = time.Now().UTC()
    return nil
}

func (c *Customer) AddContact(contact Contact) error {
    if len(c.Contacts) >= 20 {
        return domainerrors.ErrTooManyContacts
    }
    if contact.IsPrimary {
        // ปลด primary คนเก่า
        for i := range c.Contacts {
            c.Contacts[i].IsPrimary = false
        }
    }
    c.Contacts = append(c.Contacts, contact)
    c.UpdatedAt = time.Now().UTC()
    return nil
}

func (c *Customer) UpdateCreditLimit(amount vo.Money) error {
    if amount.Amount < 0 {
        return domainerrors.ErrNegativeCreditLimit
    }
    if amount.Currency != c.CreditLimit.Currency {
        return domainerrors.ErrCurrencyMismatch
    }
    c.CreditLimit = amount
    c.UpdatedAt = time.Now().UTC()
    return nil
}
```

### C.3 Package Context — `Package` + `Subscription`

```go
// internal/modules/package/domain/entity/package.go
package entity

type Package struct {
    ID           uuid.UUID
    TenantID     *uuid.UUID  // nil = platform-level
    Code         vo.PackageCode
    Name         string
    Description  string
    Tier         vo.PackageTier      // FREE, BASIC, PRO, ENTERPRISE
    Cycle        vo.BillingCycle     // MONTHLY, YEARLY
    Price        vo.Money
    PriceHistory []vo.PricePoint
    Quotas       vo.QuotaLimits
    Features     []Feature
    TrialPolicy  vo.TrialPolicy

    // Status
    IsActive     bool
    IsPublic     bool
    IsDeprecated bool
    DeprecatedAt *time.Time
    DeprecReason string

    Version  int
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy uuid.UUID

    events []domainevent.DomainEvent
}

func NewPackage(
    tenantID *uuid.UUID,
    code vo.PackageCode,
    name string,
    tier vo.PackageTier,
    cycle vo.BillingCycle,
    price vo.Money,
    quotas vo.QuotaLimits,
    actorID uuid.UUID,
) (*Package, error) {
    if !code.IsValid() {
        return nil, domainerrors.ErrInvalidPackageCode
    }
    if !tier.IsValid() {
        return nil, domainerrors.ErrInvalidTier
    }
    if !cycle.IsValid() {
        return nil, domainerrors.ErrInvalidCycle
    }
    if price.Amount < 0 {
        return nil, domainerrors.ErrNegativePrice
    }
    if err := quotas.Validate(); err != nil {
        return nil, err
    }

    now := time.Now().UTC()
    p := &Package{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Code:      code,
        Name:      name,
        Tier:      tier,
        Cycle:     cycle,
        Price:     price,
        PriceHistory: []vo.PricePoint{
            {Version: 1, Price: price, EffectiveAt: now, ChangedBy: actorID, ChangedAt: now},
        },
        Quotas:    quotas,
        IsActive:  true,
        IsPublic:  true,
        Version:   1,
        CreatedAt: now,
        UpdatedAt: now,
        CreatedBy: actorID,
    }
    p.addEvent(domainevent.NewPackageCreated(p.ID, code.String()))
    return p, nil
}

func (p *Package) UpdatePrice(newPrice vo.Money, actorID uuid.UUID) error {
    if newPrice.Amount < 0 {
        return domainerrors.ErrNegativePrice
    }
    if newPrice.Currency != p.Price.Currency {
        return domainerrors.ErrCurrencyMismatch
    }
    if newPrice.Amount == p.Price.Amount {
        return nil // no-op
    }
    now := time.Now().UTC()
    p.Version++
    p.Price = newPrice
    p.PriceHistory = append(p.PriceHistory, vo.PricePoint{
        Version:     p.Version,
        Price:       newPrice,
        EffectiveAt: now,
        ChangedBy:   actorID,
        ChangedAt:   now,
    })
    p.UpdatedAt = now
    p.addEvent(domainevent.NewPackagePriceChanged(p.ID, p.Price, newPrice))
    return nil
}

func (p *Package) Deprecate(reason string, actorID uuid.UUID) error {
    if p.IsDeprecated {
        return domainerrors.ErrAlreadyDeprecated
    }
    now := time.Now().UTC()
    p.IsDeprecated = true
    p.DeprecatedAt = &now
    p.DeprecReason = reason
    p.IsActive = false
    p.UpdatedAt = now
    p.addEvent(domainevent.NewPackageDeprecated(p.ID, reason, actorID))
    return nil
}

func (p *Package) CanBeSubscribed() error {
    if p.IsDeprecated {
        return domainerrors.ErrPackageDeprecated
    }
    if !p.IsActive {
        return domainerrors.ErrPackageInactive
    }
    return nil
}
```

### C.4 Device Context — `Device` (กับ child entities + domain events)

```go
// internal/modules/device/domain/entity/device.go
package entity

type Device struct {
    ID         uuid.UUID
    TenantID   uuid.UUID
    CustomerID uuid.UUID
    SiteID     uuid.UUID
    ZoneID     uuid.UUID
    ModelID    uuid.UUID

    // Identity
    SerialNo vo.DeviceSerial
    Name     string
    Type     vo.DeviceType
    Protocol vo.Protocol

    // Status
    Status         vo.DeviceStatus
    LastSeenAt     *time.Time
    OfflineReason  string
    FaultCode      string
    FaultMessage   string

    // Provisioning
    MQTTClientID     string
    DeviceTokenHash  string
    FirmwareVersion  string
    FirmwareChannel  string
    ProvisionedAt    *time.Time

    // Config
    Capabilities []vo.DeviceCapability
    Config       map[string]interface{}
    Tags         []string
    Metadata     map[string]interface{}

    // Children
    Commands []Command

    CreatedBy uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time

    events []domainevent.DomainEvent
}

// ---- Heartbeat ----
func (d *Device) Heartbeat() {
    now := time.Now().UTC()
    wasOffline := d.Status == vo.DeviceStatusOffline
    d.Status = vo.DeviceStatusOnline
    d.LastSeenAt = &now
    d.UpdatedAt = now

    if wasOffline {
        d.addEvent(domainevent.NewDeviceOnline(d.ID, d.TenantID, now))
    }
}

// ---- Detect offline (scheduled job) ----
func (d *Device) MarkOfflineIfStale(threshold time.Duration) bool {
    if d.LastSeenAt == nil {
        return false
    }
    if time.Since(*d.LastSeenAt) < threshold {
        return false
    }
    if d.Status == vo.DeviceStatusOffline {
        return false
    }
    d.Status = vo.DeviceStatusOffline
    d.OfflineReason = "heartbeat_timeout"
    d.UpdatedAt = time.Now().UTC()
    d.addEvent(domainevent.NewDeviceOffline(d.ID, d.TenantID, d.LastSeenAt))
    return true
}

// ---- Fault ----
func (d *Device) ReportFault(code, message string) error {
    if code == "" {
        return domainerrors.ErrInvalidFaultCode
    }
    d.Status = vo.DeviceStatusFault
    d.FaultCode = code
    d.FaultMessage = message
    d.UpdatedAt = time.Now().UTC()
    d.addEvent(domainevent.NewDeviceFault(d.ID, d.TenantID, code, message))
    return nil
}

func (d *Device) ClearFault(actorID uuid.UUID) error {
    if d.Status != vo.DeviceStatusFault {
        return domainerrors.ErrNotInFault
    }
    d.Status = vo.DeviceStatusOffline // ต้อง heartbeat ใหม่
    d.FaultCode = ""
    d.FaultMessage = ""
    d.UpdatedAt = time.Now().UTC()
    d.addEvent(domainevent.NewDeviceFaultCleared(d.ID, d.TenantID, actorID))
    return nil
}

// ---- Firmware ----
func (d *Device) UpdateFirmware(version, channel string, actorID uuid.UUID) error {
    if !isSemver(version) {
        return domainerrors.ErrInvalidFirmwareVersion
    }
    if channel != "stable" && channel != "beta" && channel != "canary" {
        return domainerrors.ErrInvalidFirmwareChannel
    }
    old := d.FirmwareVersion
    d.FirmwareVersion = version
    d.FirmwareChannel = channel
    d.UpdatedAt = time.Now().UTC()
    d.addEvent(domainevent.NewDeviceFirmwareUpdated(d.ID, d.TenantID, old, version, actorID))
    return nil
}

// ---- Command ----
func (d *Device) IssueCommand(
    cmdType string,
    payload map[string]interface{},
    priority int,
    issuedBy uuid.UUID,
    ttl time.Duration,
) (*Command, error) {
    if !d.Status.IsOperational() {
        return nil, domainerrors.ErrDeviceOffline
    }
    if cmdType == "" {
        return nil, domainerrors.ErrInvalidCommandType
    }
    if priority < 1 || priority > 5 {
        return nil, domainerrors.ErrInvalidPriority
    }
    if !d.IsCommandSupported(cmdType) {
        return nil, domainerrors.ErrCommandNotSupported
    }

    now := time.Now().UTC()
    cmd := Command{
        ID:        uuid.New(),
        DeviceID:  d.ID,
        TenantID:  d.TenantID,
        Type:      cmdType,
        Payload:   payload,
        Status:    vo.CommandStatusPending,
        Priority:  priority,
        IssuedBy:  issuedBy,
        IssuedAt:  now,
        ExpiresAt: now.Add(ttl),
    }
    d.Commands = append(d.Commands, cmd)
    d.UpdatedAt = now
    d.addEvent(domainevent.NewCommandIssued(d.ID, d.TenantID, cmd.ID, cmdType, issuedBy))
    return &cmd, nil
}

func (d *Device) IsCommandSupported(cmd string) bool {
    for _, c := range d.Capabilities {
        if c.Type == "ACTUATOR" && c.Command == cmd {
            return true
        }
    }
    return false
}

// ---- Provisioning ----
func (d *Device) Provision(tokenHash, mqttClientID string) error {
    if d.ProvisionedAt != nil {
        return domainerrors.ErrAlreadyProvisioned
    }
    if tokenHash == "" || mqttClientID == "" {
        return domainerrors.ErrInvalidProvisioning
    }
    now := time.Now().UTC()
    d.DeviceTokenHash = tokenHash
    d.MQTTClientID = mqttClientID
    d.ProvisionedAt = &now
    d.UpdatedAt = now
    d.addEvent(domainevent.NewDeviceProvisioned(d.ID, d.TenantID, mqttClientID))
    return nil
}

func (d *Device) RotateToken(newTokenHash string, actorID uuid.UUID) error {
    if newTokenHash == "" {
        return domainerrors.ErrInvalidToken
    }
    d.DeviceTokenHash = newTokenHash
    d.UpdatedAt = time.Now().UTC()
    d.addEvent(domainevent.NewDeviceTokenRotated(d.ID, d.TenantID, actorID))
    return nil
}
```

### C.5 ERP Context — `Order` Aggregate

```go
// internal/modules/erp/domain/entity/order.go
package entity

type Order struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    OrderNo   vo.DocumentNumber
    Type      vo.OrderType      // SALES, PURCHASE
    Status    vo.OrderStatus    // DRAFT, CONFIRMED, SHIPPED, RECEIVED, CLOSED, CANCELLED

    CustomerID *uuid.UUID
    SupplierID *uuid.UUID

    OrderDate   time.Time
    ConfirmedAt *time.Time
    ShippedAt   *time.Time
    ReceivedAt  *time.Time
    CancelledAt *time.Time
    CancelReason string

    ShipTo vo.Address
    BillTo vo.Address

    Lines []OrderLine

    Subtotal  vo.Money
    TaxAmount vo.Money
    TotalAmount vo.Money

    Currency     string
    PaymentTerms vo.PaymentTerms

    CreatedBy  uuid.UUID
    ApprovedBy *uuid.UUID
    CreatedAt  time.Time
    UpdatedAt  time.Time

    events []domainevent.DomainEvent
}

// ---- Line Management ----
func (o *Order) AddLine(line OrderLine) error {
    if o.Status != vo.OrderStatusDraft {
        return domainerrors.ErrCannotModifyNonDraft
    }
    if line.Quantity <= 0 {
        return domainerrors.ErrInvalidQuantity
    }
    if line.UnitPrice.Amount < 0 {
        return domainerrors.ErrNegativePrice
    }
    line.LineNo = len(o.Lines) + 1
    line.OrderID = o.ID
    o.Lines = append(o.Lines, line)
    o.recalculate()
    o.UpdatedAt = time.Now().UTC()
    return nil
}

func (o *Order) RemoveLine(lineID uuid.UUID) error {
    if o.Status != vo.OrderStatusDraft {
        return domainerrors.ErrCannotModifyNonDraft
    }
    idx := -1
    for i, l := range o.Lines {
        if l.ID == lineID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return domainerrors.ErrLineNotFound
    }
    o.Lines = append(o.Lines[:idx], o.Lines[idx+1:]...)
    // re-number
    for i := range o.Lines {
        o.Lines[i].LineNo = i + 1
    }
    o.recalculate()
    o.UpdatedAt = time.Now().UTC()
    return nil
}

func (o *Order) recalculate() {
    subtotal := decimal.Zero
    tax := decimal.Zero
    for i := range o.Lines {
        l := &o.Lines[i]
        l.Subtotal = l.UnitPrice.Mul(l.Quantity)
        l.TaxAmount = l.Subtotal.Mul(l.TaxRate)
        l.LineTotal = l.Subtotal.Add(l.TaxAmount)
        subtotal = subtotal.Add(l.Subtotal.Amount)
        tax = tax.Add(l.TaxAmount.Amount)
    }
    o.Subtotal = vo.Money{Amount: subtotal, Currency: o.Currency}
    o.TaxAmount = vo.Money{Amount: tax, Currency: o.Currency}
    o.TotalAmount = vo.Money{Amount: subtotal.Add(tax), Currency: o.Currency}
}

// ---- State Transitions ----
func (o *Order) SubmitForApproval() error {
    if o.Status != vo.OrderStatusDraft {
        return domainerrors.ErrInvalidTransition
    }
    if len(o.Lines) == 0 {
        return domainerrors.ErrNoLines
    }
    if o.Type == vo.OrderTypeSales && o.CustomerID == nil {
        return domainerrors.ErrCustomerRequired
    }
    if o.Type == vo.OrderTypePurchase && o.SupplierID == nil {
        return domainerrors.ErrSupplierRequired
    }
    o.Status = vo.OrderStatusPendingApproval
    o.UpdatedAt = time.Now().UTC()
    return nil
}

func (o *Order) Confirm(approverID uuid.UUID) error {
    if o.Status != vo.OrderStatusPendingApproval && o.Status != vo.OrderStatusDraft {
        return domainerrors.ErrInvalidTransition
    }
    now := time.Now().UTC()
    o.Status = vo.OrderStatusConfirmed
    o.ConfirmedAt = &now
    o.ApprovedBy = &approverID
    o.UpdatedAt = now
    o.addEvent(domainevent.NewOrderConfirmed(o.ID, o.TenantID, approverID))
    return nil
}

func (o *Order) Ship(shippedBy uuid.UUID) error {
    if o.Status != vo.OrderStatusConfirmed {
        return domainerrors.ErrMustBeConfirmed
    }
    if o.Type != vo.OrderTypeSales {
        return domainerrors.ErrOnlySalesCanShip
    }
    now := time.Now().UTC()
    o.Status = vo.OrderStatusShipped
    o.ShippedAt = &now
    o.UpdatedAt = now
    for i := range o.Lines {
        o.Lines[i].QtyShipped = o.Lines[i].Quantity
    }
    o.addEvent(domainevent.NewOrderShipped(o.ID, o.TenantID, shippedBy))
    return nil
}

func (o *Order) Receive(receivedBy uuid.UUID) error {
    if o.Status != vo.OrderStatusConfirmed {
        return domainerrors.ErrMustBeConfirmed
    }
    if o.Type != vo.OrderTypePurchase {
        return domainerrors.ErrOnlyPurchaseCanReceive
    }
    now := time.Now().UTC()
    o.Status = vo.OrderStatusReceived
    o.ReceivedAt = &now
    o.UpdatedAt = now
    for i := range o.Lines {
        o.Lines[i].QtyReceived = o.Lines[i].Quantity
    }
    o.addEvent(domainevent.NewOrderReceived(o.ID, o.TenantID, receivedBy))
    return nil
}

func (o *Order) Close() error {
    if o.Status != vo.OrderStatusShipped && o.Status != vo.OrderStatusReceived {
        return domainerrors.ErrCannotClose
    }
    o.Status = vo.OrderStatusClosed
    o.UpdatedAt = time.Now().UTC()
    o.addEvent(domainevent.NewOrderClosed(o.ID, o.TenantID))
    return nil
}

func (o *Order) Cancel(reason string, actorID uuid.UUID) error {
    if o.Status == vo.OrderStatusShipped || o.Status == vo.OrderStatusReceived || o.Status == vo.OrderStatusClosed {
        return domainerrors.ErrCannotCancelAfterShip
    }
    if reason == "" {
        return domainerrors.ErrCancelReasonRequired
    }
    now := time.Now().UTC()
    o.Status = vo.OrderStatusCancelled
    o.CancelledAt = &now
    o.CancelReason = reason
    o.UpdatedAt = now
    o.addEvent(domainevent.NewOrderCancelled(o.ID, o.TenantID, reason, actorID))
    return nil
}
```

### C.6 Aggregate — Others (ย่อ)

```go
// CRM Context
type Lead struct {
    ID, TenantID uuid.UUID
    Name, Email, Phone, Company string
    Source   vo.LeadSource
    Status   vo.LeadStatus
    Score    int
    AssignedTo, ConvertedTo *uuid.UUID
    CreatedAt, UpdatedAt time.Time
    events []domainevent.DomainEvent
}
func (l *Lead) Qualify(score int) error { /* 0-100 check */ }
func (l *Lead) Convert(customerID uuid.UUID) error { /* once only */ }
func (l *Lead) AssignTo(userID uuid.UUID) error { /* ... */ }

// Logistics Context
type Shipment struct {
    ID, TenantID uuid.UUID
    TrackingNo vo.TrackingNumber
    OrderID, CustomerID uuid.UUID
    Status vo.ShipmentStatus
    Origin, Destination vo.Location
    Items []ShipmentItem
    Temperature *vo.TemperatureLog
    // ...
}
func (s *Shipment) Dispatch(driverID uuid.UUID) error { /* ... */ }
func (s *Shipment) RecordTemperature(reading vo.TemperatureReading) error {
    // cold chain 2-8°C
    if reading.Value < 2 || reading.Value > 8 {
        s.addEvent(domainevent.NewColdChainBreached(s.ID, reading))
    }
    // ...
}
func (s *Shipment) MarkDelivered(receiver, signature string) error { /* ... */ }

// Report Context
type ReportDefinition struct {
    ID, TenantID uuid.UUID
    Code vo.ReportCode
    Name string
    Category vo.ReportCategory
    Query string
    Parameters []Parameter
    Schedule *Schedule
    CreatedBy uuid.UUID
    events []domainevent.DomainEvent
}
func (r *ReportDefinition) ValidateQuery() error { /* SQL check */ }
func (r *ReportDefinition) SetSchedule(s Schedule) error { /* cron check */ }
```

---

## 🅳 PART 5D — VALUE OBJECTS (Complete Catalog)

### D.1 VO Inventory

| VO | Module | หน้าที่ | Immutable |
|:---|:---|:---|:-:|
| `CustomerCode` | customer | รหัสลูกค้า `CUS-YYYY-NNNN` | ✅ |
| `CustomerType` | customer | ประเภทลูกค้า | ✅ |
| `Segment` | customer | ระดับลูกค้า | ✅ |
| `Lifecycle` | customer | lifecycle | ✅ |
| `TaxID` | customer | เลขผู้เสียภาษี (13 หลัก) | ✅ |
| `Address` | shared | ที่อยู่ | ✅ |
| `Email` | shared | อีเมล | ✅ |
| `Phone` | shared | เบอร์โทร | ✅ |
| `Money` | shared | เงิน + currency | ✅ |
| `PackageCode` | package | รหัสแพ็กเกจ | ✅ |
| `Tier` | package | ระดับแพ็กเกจ | ✅ |
| `BillingCycle` | package | รอบบิล | ✅ |
| `QuotaLimits` | package | โควตา | ✅ |
| `DeviceSerial` | device | serial number | ✅ |
| `DeviceType` | device | ประเภทอุปกรณ์ | ✅ |
| `Protocol` | device | protocol | ✅ |
| `DeviceStatus` | device | สถานะอุปกรณ์ | ✅ |
| `CommandStatus` | device | สถานะคำสั่ง | ✅ |
| `DocumentNumber` | erp | เลขเอกสาร | ✅ |
| `SKU` | erp | รหัสสินค้า | ✅ |
| `TaxType` | erp | ประเภทภาษี | ✅ |
| `OrderStatus` | erp | สถานะออเดอร์ | ✅ |
| `InvoiceStatus` | erp | สถานะใบแจ้งหนี้ | ✅ |
| `PaymentTerms` | erp | เงื่อนไขชำระ | ✅ |
| `LeadStatus` | crm | สถานะ lead | ✅ |
| `TicketPriority` | crm | ความสำคัญ ticket | ✅ |
| `TrackingNumber` | logistics | เลขติดตาม | ✅ |
| `ShipmentStatus` | logistics | สถานะ shipment | ✅ |
| `JobStatus` | logistics | สถานะงานติดตั้ง | ✅ |
| `GPS` | logistics | พิกัด | ✅ |
| `ReportCategory` | report | หมวดหมู่ report | ✅ |

### D.2 Shared VOs — `Money`

```go
// internal/shared/domain/value_object/money.go
package valueobject

import (
    "errors"
    "github.com/shopspring/decimal"
)

type Money struct {
    Amount   decimal.Decimal
    Currency string // THB, USD
}

func NewMoney(amount float64, currency string) (Money, error) {
    if currency == "" {
        return Money{}, errors.New("currency required")
    }
    if len(currency) != 3 {
        return Money{}, errors.New("currency must be ISO 4217 (3 chars)")
    }
    return Money{
        Amount:   decimal.NewFromFloat(amount),
        Currency: currency,
    }, nil
}

func ZeroMoney(currency string) Money {
    return Money{Amount: decimal.Zero, Currency: currency}
}

func (m Money) IsZero() bool {
    return m.Amount.IsZero()
}

func (m Money) IsNegative() bool {
    return m.Amount.IsNegative()
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}, nil
}

func (m Money) Sub(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{Amount: m.Amount.Sub(other.Amount), Currency: m.Currency}, nil
}

func (m Money) Mul(qty decimal.Decimal) Money {
    return Money{Amount: m.Amount.Mul(qty), Currency: m.Currency}
}

func (m Money) Equals(other Money) bool {
    return m.Currency == other.Currency && m.Amount.Equal(other.Amount)
}

func (m Money) String() string {
    return m.Amount.StringFixed(2) + " " + m.Currency
}
```

### D.3 Device VOs

```go
// internal/modules/device/domain/value_object/device_serial.go
package valueobject

import (
    "errors"
    "regexp"
    "strings"
)

var serialRegex = regexp.MustCompile(`^[A-Z0-9][A-Z0-9\-]{2,49}$`)

type DeviceSerial string

func NewDeviceSerial(s string) (DeviceSerial, error) {
    s = strings.ToUpper(strings.TrimSpace(s))
    if !serialRegex.MatchString(s) {
        return "", errors.New("invalid serial format")
    }
    return DeviceSerial(s), nil
}

func (s DeviceSerial) String() string { return string(s) }
func (s DeviceSerial) IsValid() bool  { return serialRegex.MatchString(string(s)) }

// ---- Device Status ----
type DeviceStatus string

const (
    DeviceStatusOnline      DeviceStatus = "ONLINE"
    DeviceStatusOffline     DeviceStatus = "OFFLINE"
    DeviceStatusFault       DeviceStatus = "FAULT"
    DeviceStatusMaintenance DeviceStatus = "MAINTENANCE"
    DeviceStatusProvisioned DeviceStatus = "PROVISIONED"
)

func (s DeviceStatus) IsValid() bool {
    switch s {
    case DeviceStatusOnline, DeviceStatusOffline, DeviceStatusFault,
        DeviceStatusMaintenance, DeviceStatusProvisioned:
        return true
    }
    return false
}

func (s DeviceStatus) IsOperational() bool { return s == DeviceStatusOnline }
func (s DeviceStatus) CanReceiveCommands() bool {
    return s == DeviceStatusOnline
}

// ---- Device Type ----
type DeviceType string

const (
    DeviceTypeSensor   DeviceType = "SENSOR"
    DeviceTypeActuator DeviceType = "ACTUATOR"
    DeviceTypeGateway  DeviceType = "GATEWAY"
    DeviceTypeCamera   DeviceType = "CAMERA"
    DeviceTypeMeter    DeviceType = "METER"
    DeviceTypeTracker  DeviceType = "TRACKER"
)

func (t DeviceType) IsValid() bool {
    switch t {
    case DeviceTypeSensor, DeviceTypeActuator, DeviceTypeGateway,
        DeviceTypeCamera, DeviceTypeMeter, DeviceTypeTracker:
        return true
    }
    return false
}

func (t DeviceType) IsControllable() bool {
    return t == DeviceTypeActuator
}

// ---- Protocol ----
type Protocol string

const (
    ProtocolMQTT    Protocol = "MQTT"
    ProtocolModbus  Protocol = "MODBUS"
    ProtocolLoRaWAN Protocol = "LORAWAN"
    ProtocolZigbee  Protocol = "ZIGBEE"
    ProtocolHTTP    Protocol = "HTTP"
    ProtocolCoAP    Protocol = "COAP"
)

func (p Protocol) IsValid() bool {
    switch p {
    case ProtocolMQTT, ProtocolModbus, ProtocolLoRaWAN,
        ProtocolZigbee, ProtocolHTTP, ProtocolCoAP:
        return true
    }
    return false
}

// ---- Command Status ----
type CommandStatus string

const (
    CommandStatusPending  CommandStatus = "PENDING"
    CommandStatusSent     CommandStatus = "SENT"
    CommandStatusAcked    CommandStatus = "ACKED"
    CommandStatusFailed   CommandStatus = "FAILED"
    CommandStatusTimeout  CommandStatus = "TIMEOUT"
    CommandStatusExpired  CommandStatus = "EXPIRED"
)

func (s CommandStatus) IsTerminal() bool {
    switch s {
    case CommandStatusAcked, CommandStatusFailed,
        CommandStatusTimeout, CommandStatusExpired:
        return true
    }
    return false
}
```

### D.4 ERP VOs

```go
// internal/modules/erp/domain/value_object/order_status.go
package valueobject

type OrderStatus string

const (
    OrderStatusDraft            OrderStatus = "DRAFT"
    OrderStatusPendingApproval  OrderStatus = "PENDING_APPROVAL"
    OrderStatusConfirmed        OrderStatus = "CONFIRMED"
    OrderStatusShipped          OrderStatus = "SHIPPED"
    OrderStatusReceived         OrderStatus = "RECEIVED"
    OrderStatusInvoiced         OrderStatus = "INVOICED"
    OrderStatusClosed           OrderStatus = "CLOSED"
    OrderStatusCancelled        OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsValid() bool {
    switch s {
    case OrderStatusDraft, OrderStatusPendingApproval, OrderStatusConfirmed,
        OrderStatusShipped, OrderStatusReceived, OrderStatusInvoiced,
        OrderStatusClosed, OrderStatusCancelled:
        return true
    }
    return false
}

func (s OrderStatus) IsTerminal() bool {
    return s == OrderStatusClosed || s == OrderStatusCancelled
}

func (s OrderStatus) CanModifyLines() bool {
    return s == OrderStatusDraft
}

// ---- DocumentNumber ----
type DocumentNumber string

// Format: <PREFIX>-<YYYY>-<NNNNNN>
//   SO  → Sales Order
//   PO  → Purchase Order
//   INV → Invoice
//   BILL→ Bill
//   PAY → Payment
//   JOB → Installation Job
//   SHP → Shipment

type DocType string

const (
    DocTypeSO   DocType = "SO"
    DocTypePO   DocType = "PO"
    DocTypeINV  DocType = "INV"
    DocTypeBILL DocType = "BILL"
    DocTypePAY  DocType = "PAY"
    DocTypeJOB  DocType = "JOB"
    DocTypeSHP  DocType = "SHP"
)

func (t DocType) IsValid() bool {
    switch t {
    case DocTypeSO, DocTypePO, DocTypeINV, DocTypeBILL, DocTypePAY, DocTypeJOB, DocTypeSHP:
        return true
    }
    return false
}

func NewDocumentNumber(dt DocType, year, seq int) (DocumentNumber, error) {
    if !dt.IsValid() {
        return "", errors.New("invalid doc type")
    }
    if year < 2020 || year > 2100 {
        return "", errors.New("invalid year")
    }
    if seq < 1 || seq > 999999 {
        return "", errors.New("invalid sequence")
    }
    return DocumentNumber(fmt.Sprintf("%s-%d-%06d", dt, year, seq)), nil
}

func (n DocumentNumber) String() string { return string(n) }
func (n DocumentNumber) DocType() DocType {
    parts := strings.SplitN(string(n), "-", 2)
    if len(parts) > 0 {
        return DocType(parts[0])
    }
    return ""
}

// ---- TaxType ----
type TaxType string

const (
    TaxTypeVAT      TaxType = "VAT"
    TaxTypeZero     TaxType = "ZERO"
    TaxTypeExempt   TaxType = "EXEMPT"
    TaxTypeNone     TaxType = "NONE"
)

func (t TaxType) Rate() float64 {
    switch t {
    case TaxTypeVAT:    return 0.07
    case TaxTypeZero:   return 0.00
    case TaxTypeExempt: return 0.00
    default:            return 0.00
    }
}
```

### D.5 VO Template + Test

```go
// internal/modules/{{module}}/domain/value_object/{{vo}}.go
package valueobject

import (
    "errors"
    "regexp"
    "strings"
)

var {{vo}}Regex = regexp.MustCompile(`^{{pattern}}$`)

type {{VO}} string

func New{{VO}}(s string) ({{VO}}, error) {
    s = strings.TrimSpace(s)
    if !{{vo}}Regex.MatchString(s) {
        return "", errors.New("invalid {{vo}} format")
    }
    return {{VO}}(s), nil
}

func (v {{VO}}) String() string { return string(v) }
func (v {{VO}}) IsValid() bool  { return {{vo}}Regex.MatchString(string(v)) }
func (v {{VO}}) Equals(o {{VO}}) bool { return v == o }
```

```go
// internal/modules/{{module}}/domain/value_object/{{vo}}_test.go
package valueobject_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "icmongolang/internal/modules/{{module}}/domain/value_object"
)

func TestNew{{VO}}(t *testing.T) {
    t.Run("valid", func(t *testing.T) {
        v, err := valueobject.New{{VO}}("VALID-001")
        assert.NoError(t, err)
        assert.True(t, v.IsValid())
    })
    t.Run("invalid", func(t *testing.T) {
        _, err := valueobject.New{{VO}}("bad format")
        assert.Error(t, err)
    })
    t.Run("empty", func(t *testing.T) {
        _, err := valueobject.New{{VO}}("")
        assert.Error(t, err)
    })
}
```

---

## 🅴 PART 5E — DOMAIN EVENTS

### E.1 Event Naming Convention

```
<context>.<aggregate>.<past-tense-verb>

device.device.created
device.device.online
device.device.offline
device.command.issued
device.command.acked
device.telemetry.ingested

customer.customer.created
customer.customer.activated
customer.customer.suspended

package.package.created
package.package.price_changed
package.subscription.created
package.subscription.upgraded

erp.order.created
erp.order.confirmed
erp.order.shipped
erp.invoice.issued
erp.payment.received

logistics.shipment.dispatched
logistics.shipment.delivered
logistics.coldchain.breached

crm.lead.created
crm.lead.converted
crm.ticket.opened
```

### E.2 Domain Event Interface

```go
// internal/shared/domain/event/event.go
package event

import (
    "time"
    "github.com/google/uuid"
)

type DomainEvent interface {
    EventID() uuid.UUID
    EventName() string
    OccurredAt() time.Time
    AggregateID() uuid.UUID
    TenantID() uuid.UUID
}

// BaseEvent — embed สำหรับ event ทุกตัว
type BaseEvent struct {
    ID        uuid.UUID `json:"event_id"`
    Name      string    `json:"event_name"`
    At        time.Time `json:"occurred_at"`
    AggID     uuid.UUID `json:"aggregate_id"`
    TID       uuid.UUID `json:"tenant_id"`
    Version   int       `json:"version"`
}

func (b BaseEvent) EventID() uuid.UUID     { return b.ID }
func (b BaseEvent) EventName() string      { return b.Name }
func (b BaseEvent) OccurredAt() time.Time  { return b.At }
func (b BaseEvent) AggregateID() uuid.UUID { return b.AggID }
func (b BaseEvent) TenantID() uuid.UUID    { return b.TID }

func NewBase(name string, aggID, tenantID uuid.UUID) BaseEvent {
    return BaseEvent{
        ID:      uuid.New(),
        Name:    name,
        At:      time.Now().UTC(),
        AggID:   aggID,
        TID:     tenantID,
        Version: 1,
    }
}
```

### E.3 Device Events (complete)

```go
// internal/modules/device/domain/event/device_events.go
package event

import (
    "time"
    "github.com/google/uuid"
    sharedevent "icmongolang/internal/shared/domain/event"
)

// ---- Device Lifecycle ----
type DeviceCreated struct {
    sharedevent.BaseEvent
    SerialNo string `json:"serial_no"`
    Type     string `json:"type"`
    SiteID   uuid.UUID `json:"site_id"`
}

func NewDeviceCreated(id, tenantID uuid.UUID, serial, dtype string, siteID uuid.UUID) DeviceCreated {
    return DeviceCreated{
        BaseEvent: sharedevent.NewBase("device.device.created", id, tenantID),
        SerialNo:  serial,
        Type:      dtype,
        SiteID:    siteID,
    }
}

type DeviceOnline struct {
    sharedevent.BaseEvent
    LastSeenAt time.Time `json:"last_seen_at"`
}

func NewDeviceOnline(id, tenantID uuid.UUID, at time.Time) DeviceOnline {
    return DeviceOnline{
        BaseEvent:  sharedevent.NewBase("device.device.online", id, tenantID),
        LastSeenAt: at,
    }
}

type DeviceOffline struct {
    sharedevent.BaseEvent
    LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
    Reason     string     `json:"reason"`
}

func NewDeviceOffline(id, tenantID uuid.UUID, lastSeen *time.Time) DeviceOffline {
    return DeviceOffline{
        BaseEvent:  sharedevent.NewBase("device.device.offline", id, tenantID),
        LastSeenAt: lastSeen,
        Reason:     "heartbeat_timeout",
    }
}

type DeviceFault struct {
    sharedevent.BaseEvent
    FaultCode    string `json:"fault_code"`
    FaultMessage string `json:"fault_message"`
}

func NewDeviceFault(id, tenantID uuid.UUID, code, msg string) DeviceFault {
    return DeviceFault{
        BaseEvent:    sharedevent.NewBase("device.device.fault", id, tenantID),
        FaultCode:    code,
        FaultMessage: msg,
    }
}

type DeviceFirmwareUpdated struct {
    sharedevent.BaseEvent
    FromVersion string `json:"from_version"`
    ToVersion   string `json:"to_version"`
    UpdatedBy   uuid.UUID `json:"updated_by"`
}

// ---- Command ----
type CommandIssued struct {
    sharedevent.BaseEvent
    CommandID uuid.UUID              `json:"command_id"`
    Type      string                 `json:"type"`
    Payload   map[string]interface{} `json:"payload"`
    IssuedBy  uuid.UUID              `json:"issued_by"`
}

func NewCommandIssued(deviceID, tenantID, cmdID uuid.UUID, cmdType string, by uuid.UUID) CommandIssued {
    return CommandIssued{
        BaseEvent: sharedevent.NewBase("device.command.issued", deviceID, tenantID),
        CommandID: cmdID,
        Type:      cmdType,
        IssuedBy:  by,
    }
}

type CommandAcked struct {
    sharedevent.BaseEvent
    CommandID uuid.UUID              `json:"command_id"`
    Status    string                 `json:"status"`
    Result    map[string]interface{} `json:"result,omitempty"`
}

type CommandFailed struct {
    sharedevent.BaseEvent
    CommandID uuid.UUID `json:"command_id"`
    Error     string    `json:"error"`
}

// ---- Telemetry ----
type TelemetryIngested struct {
    sharedevent.BaseEvent
    Metrics   []MetricReading `json:"metrics"`
    Timestamp time.Time       `json:"timestamp"`
}

type MetricReading struct {
    Metric string  `json:"metric"`
    Value  float64 `json:"value"`
    Unit   string  `json:"unit"`
}

// ---- Automation ----
type AutomationTriggered struct {
    sharedevent.BaseEvent
    RuleID    uuid.UUID `json:"rule_id"`
    RuleName  string    `json:"rule_name"`
    Triggered time.Time `json:"triggered"`
}
```

### E.4 Event Registry & Dispatcher Port

```go
// internal/shared/domain/event/registry.go
package event

import "sync"

var registry = struct {
    sync.RWMutex
    names map[string]bool
}{
    names: make(map[string]bool),
}

func Register(name string) {
    registry.Lock()
    defer registry.Unlock()
    registry.names[name] = true
}

func IsRegistered(name string) bool {
    registry.RLock()
    defer registry.RUnlock()
    return registry.names[name]
}

func AllRegistered() []string {
    registry.RLock()
    defer registry.RUnlock()
    out := make([]string, 0, len(registry.names))
    for n := range registry.names {
        out = append(out, n)
    }
    return out
}

// ---- Dispatcher Port (implemented in Infrastructure) ----
type Dispatcher interface {
    Dispatch(ctx context.Context, events ...DomainEvent) error
    DispatchAsync(ctx context.Context, events ...DomainEvent) error
}
```

---

## 🅵 PART 5F — DOMAIN SERVICES & SPECIFICATIONS

### F.1 Domain Service Pattern

> **Domain Service** = stateless logic ที่ไม่ผูกกับ entity ตัวใดตัวหนึ่ง

```go
// internal/modules/device/domain/service/automation_service.go
package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/modules/device/domain/event"
)

// AutomationService – evaluate rules + issue commands
type AutomationService struct{}

type EvaluationResult struct {
    Matched  []MatchedRule
    Actions  []Action
    Duration time.Duration
}

type MatchedRule struct {
    RuleID   string
    RuleName string
    Score    int
}

type Action struct {
    DeviceID uuid.UUID
    Command  string
    Payload  map[string]interface{}
    Priority int
}

// EvaluateAll – ประเมินทุก rule กับ telemetry ที่เข้ามา
func (s *AutomationService) EvaluateAll(
    rules []entity.AutomationRule,
    device *entity.Device,
    metrics []event.MetricReading,
) (*EvaluationResult, error) {
    start := time.Now()
    res := &EvaluationResult{}

    for _, rule := range rules {
        if !rule.Enabled {
            continue
        }
        matched, err := s.evaluateRule(rule, device, metrics)
        if err != nil {
            continue
        }
        if matched {
            res.Matched = append(res.Matched, MatchedRule{
                RuleID:   rule.ID.String(),
                RuleName: rule.Name,
            })
            res.Actions = append(res.Actions, rule.Actions...)
        }
    }

    res.Duration = time.Since(start)
    return res, nil
}

func (s *AutomationService) evaluateRule(
    rule entity.AutomationRule,
    device *entity.Device,
    metrics []event.MetricReading,
) (bool, error) {
    // 1. Check device scope
    if rule.DeviceID != uuid.Nil && rule.DeviceID != device.ID {
        return false, nil
    }
    if rule.SiteID != uuid.Nil && rule.SiteID != device.SiteID {
        return false, nil
    }

    // 2. Evaluate conditions (AND)
    metricMap := make(map[string]float64)
    for _, m := range metrics {
        metricMap[m.Metric] = m.Value
    }

    for _, cond := range rule.Conditions {
        val, ok := metricMap[cond.Metric]
        if !ok {
            return false, nil // metric not present
        }
        if !compare(val, cond.Operator, cond.Value) {
            return false, nil
        }
    }

    return true, nil
}

func compare(a float64, op string, b float64) bool {
    switch op {
    case ">":  return a > b
    case ">=": return a >= b
    case "<":  return a < b
    case "<=": return a <= b
    case "==": return a == b
    case "!=": return a != b
    }
    return false
}

// ---- Threshold Service ----
type ThresholdService struct{}

func (s *ThresholdService) Check(rule entity.AlertRule, reading event.MetricReading) (*entity.AlertEvent, error) {
    if !compare(reading.Value, rule.Operator, rule.Threshold) {
        return nil, nil
    }
    return &entity.AlertEvent{
        RuleID:    rule.ID,
        Metric:    reading.Metric,
        Value:     reading.Value,
        Threshold: rule.Threshold,
        Severity:  rule.Severity,
        TriggeredAt: time.Now().UTC(),
    }, nil
}

// ---- Quota Service ----
type QuotaService struct{}

type QuotaCheck struct {
    Allowed bool
    Current int
    Limit   int
    Message string
}

func (s *QuotaService) CheckDeviceQuota(pkg entity.Package, currentDevices int) QuotaCheck {
    limit := pkg.Quotas.MaxDevices
    if limit == -1 {
        return QuotaCheck{Allowed: true, Limit: -1}
    }
    if currentDevices >= limit {
        return QuotaCheck{
            Allowed: false,
            Current: currentDevices,
            Limit:   limit,
            Message: fmt.Sprintf("device quota exceeded: %d/%d", currentDevices, limit),
        }
    }
    return QuotaCheck{Allowed: true, Current: currentDevices, Limit: limit}
}

// ---- Pricing Service ----
type PricingService struct{}

func (s *PricingService) CalculateSubscriptionPrice(pkg entity.Package, cycle vo.BillingCycle, months int) (vo.Money, error) {
    base := pkg.Price
    if cycle == vo.CycleYearly {
        // yearly discount 2 months
        multiplier := decimal.NewFromInt(int64(months - 2))
        return vo.Money{
            Amount:   base.Amount.Mul(multiplier),
            Currency: base.Currency,
        }, nil
    }
    multiplier := decimal.NewFromInt(int64(months))
    return vo.Money{
        Amount:   base.Amount.Mul(multiplier),
        Currency: base.Currency,
    }, nil
}
```

### F.2 Specification Pattern

```go
// internal/shared/domain/specification/specification.go
package specification

import "context"

type Specification[T any] interface {
    IsSatisfiedBy(ctx context.Context, candidate T) (bool, error)
    And(other Specification[T]) Specification[T]
    Or(other Specification[T]) Specification[T]
    Not() Specification[T]
}

// ---- Base ----
type base[T any] struct {
    fn func(context.Context, T) (bool, error)
}

func New[T any](fn func(context.Context, T) (bool, error)) Specification[T] {
    return &base[T]{fn: fn}
}

func (s *base[T]) IsSatisfiedBy(ctx context.Context, c T) (bool, error) {
    return s.fn(ctx, c)
}

func (s *base[T]) And(other Specification[T]) Specification[T] {
    return New[T](func(ctx context.Context, c T) (bool, error) {
        ok, err := s.fn(ctx, c)
        if err != nil || !ok {
            return false, err
        }
        return other.IsSatisfiedBy(ctx, c)
    })
}

func (s *base[T]) Or(other Specification[T]) Specification[T] {
    return New[T](func(ctx context.Context, c T) (bool, error) {
        ok, err := s.fn(ctx, c)
        if err != nil {
            return false, err
        }
        if ok {
            return true, nil
        }
        return other.IsSatisfiedBy(ctx, c)
    })
}

func (s *base[T]) Not() Specification[T] {
    return New[T](func(ctx context.Context, c T) (bool, error) {
        ok, err := s.fn(ctx, c)
        return !ok, err
    })
}
```

```go
// internal/modules/device/domain/specification/device_spec.go
package specification

import (
    "context"
    "icmongolang/internal/modules/device/domain/entity"
    "icmongolang/internal/shared/domain/specification"
)

type DeviceSpec = specification.Specification[*entity.Device]

func IsOnline() DeviceSpec {
    return specification.New(func(_ context.Context, d *entity.Device) (bool, error) {
        return d.Status == vo.DeviceStatusOnline, nil
    })
}

func IsActuator() DeviceSpec {
    return specification.New(func(_ context.Context, d *entity.Device) (bool, error) {
        return d.Type == vo.DeviceTypeActuator, nil
    })
}

func InSite(siteID uuid.UUID) DeviceSpec {
    return specification.New(func(_ context.Context, d *entity.Device) (bool, error) {
        return d.SiteID == siteID, nil
    })
}

// Usage
func ExampleSpec() {
    spec := IsOnline().And(IsActuator()).And(InSite(siteID))
    ok, _ := spec.IsSatisfiedBy(ctx, device)
}
```

---

## 🅶 PART 5G — REPOSITORY INTERFACES

### G.1 Repository Interface Pattern

```go
// internal/modules/device/domain/repository/device_repository.go
package repository

import (
    "context"
    "time"
    "github.com/google/uuid"
    "icmongolang/internal/modules/device/domain/entity"
    vo "icmongolang/internal/modules/device/domain/value_object"
)

// DeviceRepository – command side (write model)
type DeviceRepository interface {
    // Create
    Save(ctx context.Context, d *entity.Device) error

    // Read (single)
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Device, error)
    FindBySerial(ctx context.Context, tenantID uuid.UUID, serial vo.DeviceSerial) (*entity.Device, error)

    // Read (list)
    FindBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*entity.Device, error)
    FindByStatus(ctx context.Context, tenantID uuid.UUID, status vo.DeviceStatus) ([]*entity.Device, error)
    FindByCustomer(ctx context.Context, tenantID, customerID uuid.UUID) ([]*entity.Device, error)

    // Update
    Update(ctx context.Context, d *entity.Device) error

    // Delete
    Delete(ctx context.Context, tenantID, id uuid.UUID) error

    // Existence
    ExistsBySerial(ctx context.Context, tenantID uuid.UUID, serial vo.DeviceSerial) (bool, error)

    // Bulk
    FindStale(ctx context.Context, tenantID uuid.UUID, threshold time.Duration, limit int) ([]*entity.Device, error)
    UpdateBatch(ctx context.Context, devices []*entity.Device) error
}

// TelemetryRepository – time-series (InfluxDB)
type TelemetryRepository interface {
    Write(ctx context.Context, t *entity.Telemetry) error
    WriteBatch(ctx context.Context, batch []*entity.Telemetry) error

    Query(ctx context.Context, q TelemetryQuery) (*TelemetryResult, error)
    Latest(ctx context.Context, tenantID, deviceID uuid.UUID, metric string) (*entity.Telemetry, error)

    // Aggregation
    Aggregate(ctx context.Context, q AggregateQuery) ([]AggregatePoint, error)
    DeleteOlderThan(ctx context.Context, before time.Time) (int64, error)
}

type TelemetryQuery struct {
    TenantID uuid.UUID
    DeviceID uuid.UUID
    Metrics  []string
    From     time.Time
    To       time.Time
    Interval string // 1m, 5m, 1h
    Limit    int
}

type TelemetryResult struct {
    DeviceID uuid.UUID
    Series   map[string][]TelemetryPoint
}

type TelemetryPoint struct {
    Timestamp time.Time
    Value     float64
    Unit      string
    Quality   string
}

type AggregateQuery struct {
    TenantID uuid.UUID
    DeviceID uuid.UUID
    Metric   string
    From     time.Time
    To       time.Time
    Window   string
    Function string // mean, max, min, sum, count
}

type AggregatePoint struct {
    Timestamp time.Time
    Value     float64
}

// CommandRepository
type CommandRepository interface {
    Save(ctx context.Context, c *entity.Command) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Command, error)
    FindPendingByDevice(ctx context.Context, tenantID, deviceID uuid.UUID) ([]*entity.Command, error)
    FindByDevice(ctx context.Context, tenantID, deviceID uuid.UUID, limit int) ([]*entity.Command, error)
    Update(ctx context.Context, c *entity.Command) error
    MarkExpired(ctx context.Context, before time.Time) (int64, error)
}

// ReadModelRepository – CQRS query side
type DeviceReadModelRepository interface {
    List(ctx context.Context, q DeviceListQuery) (*DeviceListResult, error)
    Get(ctx context.Context, tenantID, id uuid.UUID) (*DeviceReadModel, error)
    StatsBySite(ctx context.Context, tenantID uuid.UUID) ([]SiteStats, error)
}

type DeviceListQuery struct {
    TenantID   uuid.UUID
    SiteID     *uuid.UUID
    CustomerID *uuid.UUID
    Status     *vo.DeviceStatus
    Type       *vo.DeviceType
    Search     string
    Page       int
    PageSize   int
    SortBy     string // created_at, name, status
    SortOrder  string // asc, desc
}

type DeviceListResult struct {
    Items      []*DeviceReadModel
    Pagination Pagination
}

type Pagination struct {
    Page       int
    PageSize   int
    Total      int64
    TotalPages int
}

type DeviceReadModel struct {
    ID           uuid.UUID
    SerialNo     string
    Name         string
    Type         string
    Status       string
    SiteID       uuid.UUID
    SiteName     string
    LastSeenAt   *time.Time
    CreatedAt    time.Time
}

type SiteStats struct {
    SiteID       uuid.UUID
    SiteName     string
    TotalDevices int
    OnlineCount  int
    OfflineCount int
    FaultCount   int
}
```

### G.2 Repository Interface Inventory

| Module | Repository | Methods |
|:---|:---|:-:|
| customer | CustomerRepository | 12 |
| customer | CustomerReadModelRepository | 4 |
| package | PackageRepository | 8 |
| package | SubscriptionRepository | 10 |
| device | DeviceRepository | 10 |
| device | TelemetryRepository | 7 |
| device | CommandRepository | 6 |
| device | AlertRuleRepository | 6 |
| device | ShadowRepository | 5 |
| erp | OrderRepository | 10 |
| erp | InvoiceRepository | 9 |
| erp | PaymentRepository | 8 |
| erp | ProductRepository | 8 |
| erp | InventoryRepository | 10 |
| erp | JournalRepository | 6 |
| crm | LeadRepository | 8 |
| crm | TicketRepository | 8 |
| logistics | ShipmentRepository | 9 |
| logistics | JobRepository | 10 |
| logistics | TechnicianRepository | 7 |
| report | ReportDefinitionRepository | 6 |
| report | DashboardRepository | 6 |
| report | KPIRepository | 5 |

---

## 🅷 PART 5H — FACTORIES, INVARIANTS & GUARDS

### H.1 Factory Pattern

```go
// internal/modules/erp/domain/factory/order_factory.go
package factory

import (
    "time"
    "github.com/google/uuid"
    "icmongolang/internal/modules/erp/domain/entity"
    vo "icmongolang/internal/modules/erp/domain/value_object"
)

type OrderFactory struct {
    sequence SequenceGenerator  // port
}

type SequenceGenerator interface {
    Next(ctx context.Context, tenantID uuid.UUID, docType vo.DocType) (int, error)
}

func NewOrderFactory(seq SequenceGenerator) *OrderFactory {
    return &OrderFactory{sequence: seq}
}

type CreateOrderInput struct {
    TenantID    uuid.UUID
    Type        vo.OrderType
    CustomerID  *uuid.UUID
    SupplierID  *uuid.UUID
    Currency    string
    ActorID     uuid.UUID
    ShipTo      vo.Address
    BillTo      vo.Address
}

func (f *OrderFactory) CreateSalesOrder(ctx context.Context, in CreateOrderInput) (*entity.Order, error) {
    seq, err := f.sequence.Next(ctx, in.TenantID, vo.DocTypeSO)
    if err != nil {
        return nil, err
    }
    docNo, err := vo.NewDocumentNumber(vo.DocTypeSO, time.Now().Year(), seq)
    if err != nil {
        return nil, err
    }
    order, err := entity.NewOrder(in.TenantID, docNo, vo.OrderTypeSales, time.Now().UTC(), in.Currency, in.ActorID)
    if err != nil {
        return nil, err
    }
    if in.CustomerID != nil {
        if err := order.SetCustomer(*in.CustomerID, in.ShipTo); err != nil {
            return nil, err
        }
    }
    return order, nil
}
```

### H.2 Invariant Enforcement

```go
// internal/modules/package/domain/entity/invariants.go
package entity

// Invariants ที่ต้องบังคับ:
// 1. Price ≥ 0
// 2. Quotas ทุกตัว ≥ 0 (หรือ -1 = unlimited)
// 3. Tier ต้อง valid
// 4. Cycle ต้อง valid
// 5. Deprecated → IsActive = false

func (p *Package) checkInvariants() error {
    if p.Price.Amount.IsNegative() {
        return domainerrors.ErrNegativePrice
    }
    if err := p.Quotas.Validate(); err != nil {
        return err
    }
    if !p.Tier.IsValid() {
        return domainerrors.ErrInvalidTier
    }
    if !p.Cycle.IsValid() {
        return domainerrors.ErrInvalidCycle
    }
    if p.IsDeprecated && p.IsActive {
        return domainerrors.ErrDeprecatedMustBeInactive
    }
    return nil
}

// เรียกใน constructor + ทุก behavior method ที่เปลี่ยน state
func (p *Package) Validate() error {
    return p.checkInvariants()
}
```

### H.3 Guard Clauses

```go
// internal/shared/domain/guard/guard.go
package guard

import (
    "errors"
    "regexp"
    "strings"
)

func NotEmpty(field, value string) error {
    if strings.TrimSpace(value) == "" {
        return errors.New(field + " must not be empty")
    }
    return nil
}

func Length(field, value string, min, max int) error {
    l := len(value)
    if l < min || l > max {
        return errors.New(field + " length out of range")
    }
    return nil
}

func Positive(field string, v int) error {
    if v <= 0 {
        return errors.New(field + " must be positive")
    }
    return nil
}

func NonNegative(field string, v int) error {
    if v < 0 {
        return errors.New(field + " must be non-negative")
    }
    return nil
}

func UUID(field string, id uuid.UUID) error {
    if id == uuid.Nil {
        return errors.New(field + " must be a valid UUID")
    }
    return nil
}

func Email(field, value string) error {
    if !emailRegex.MatchString(value) {
        return errors.New(field + " invalid email format")
    }
    return nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
```

---

## 🅸 PART 5I — DOMAIN ERRORS (Complete Catalog)

### I.1 Error Convention

```go
// internal/modules/device/domain/errors/errors.go
package domainerrors

import "errors"

// ---- Generic ----
var (
    ErrNotFound       = errors.New("not found")
    ErrAlreadyExists  = errors.New("already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrConflict       = errors.New("conflict")
    ErrInvalidState   = errors.New("invalid state")
    ErrInvalidTransition = errors.New("invalid status transition")
)

// ---- Tenant ----
var (
    ErrInvalidTenant = errors.New("invalid tenant")
    ErrTenantMismatch = errors.New("tenant mismatch")
)

// ---- Device ----
var (
    ErrDeviceNotFound        = errors.New("device not found")
    ErrDeviceAlreadyExists   = errors.New("device already exists")
    ErrInvalidSerial         = errors.New("invalid serial number")
    ErrSerialAlreadyExists   = errors.New("serial number already exists")
    ErrInvalidDeviceType     = errors.New("invalid device type")
    ErrInvalidProtocol       = errors.New("invalid protocol")
    ErrInvalidFirmware       = errors.New("invalid firmware version")
    ErrInvalidFirmwareVersion = errors.New("invalid firmware version (semver required)")
    ErrInvalidFirmwareChannel = errors.New("invalid firmware channel")
    ErrInvalidSite           = errors.New("invalid site")
    ErrInvalidZone           = errors.New("invalid zone")
    ErrDeviceOffline         = errors.New("device is offline")
    ErrDeviceFault           = errors.New("device is in fault state")
    ErrNotInFault            = errors.New("device is not in fault state")
    ErrAlreadyProvisioned    = errors.New("device already provisioned")
    ErrInvalidProvisioning   = errors.New("invalid provisioning data")
    ErrInvalidToken          = errors.New("invalid token")
    ErrCommandNotAllowed     = errors.New("command not allowed for this device")
    ErrCommandNotSupported   = errors.New("command not supported by device capabilities")
    ErrInvalidCommandType    = errors.New("invalid command type")
    ErrInvalidFaultCode      = errors.New("invalid fault code")
    ErrInvalidPriority       = errors.New("priority must be 1-5")
    ErrCommandExpired        = errors.New("command expired")
)

// ---- Customer ----
var (
    ErrCustomerNotFound          = errors.New("customer not found")
    ErrCustomerAlreadyExists     = errors.New("customer already exists")
    ErrInvalidCustomer           = errors.New("invalid customer")
    ErrInvalidCustomerCode       = errors.New("invalid customer code")
    ErrCustomerAlreadyActive     = errors.New("customer already active")
    ErrCustomerAlreadySuspended  = errors.New("customer already suspended")
    ErrCannotReactivateCancelled = errors.New("cannot reactivate cancelled customer")
    ErrCannotSuspendCancelled    = errors.New("cannot suspend cancelled customer")
    ErrNotALead                  = errors.New("customer is not a lead")
    ErrMustQualifyFirst          = errors.New("must qualify lead before activation")
    ErrInvalidTaxID              = errors.New("invalid tax ID")
    ErrIndividualNoTaxID         = errors.New("individual customer cannot have tax ID")
    ErrTooManyContacts           = errors.New("too many contacts (max 20)")
    ErrNegativeCreditLimit       = errors.New("credit limit cannot be negative")
)

// ---- Package ----
var (
    ErrPackageNotFound          = errors.New("package not found")
    ErrPackageAlreadyExists     = errors.New("package already exists")
    ErrInvalidPackage           = errors.New("invalid package")
    ErrInvalidPackageCode       = errors.New("invalid package code")
    ErrInvalidTier              = errors.New("invalid package tier")
    ErrInvalidCycle             = errors.New("invalid billing cycle")
    ErrNegativePrice            = errors.New("price cannot be negative")
    ErrPackageInactive          = errors.New("package is inactive")
    ErrPackageDeprecated        = errors.New("package is deprecated")
    ErrAlreadyDeprecated        = errors.New("package already deprecated")
    ErrDeprecatedMustBeInactive = errors.New("deprecated package must be inactive")
    ErrCurrencyMismatch         = errors.New("currency mismatch")
    ErrInvalidQuota             = errors.New("invalid quota")
)

// ---- Subscription ----
var (
    ErrSubscriptionNotFound   = errors.New("subscription not found")
    ErrSubscriptionCancelled  = errors.New("subscription already cancelled")
    ErrSubscriptionExpired    = errors.New("subscription expired")
    ErrSubscriptionPastDue    = errors.New("subscription past due")
    ErrCannotUpgradeCancelled = errors.New("cannot upgrade cancelled subscription")
    ErrQuotaExceeded          = errors.New("quota exceeded")
)

// ---- ERP ----
var (
    ErrOrderNotFound           = errors.New("order not found")
    ErrOrderAlreadyExists      = errors.New("order already exists")
    ErrInvalidOrderNumber      = errors.New("invalid order number")
    ErrCannotModifyNonDraft    = errors.New("cannot modify non-draft order")
    ErrCannotCancelAfterShip   = errors.New("cannot cancel after shipment")
    ErrNoLines                 = errors.New("order must have at least one line")
    ErrLineNotFound            = errors.New("order line not found")
    ErrCustomerRequired        = errors.New("customer required for sales order")
    ErrSupplierRequired        = errors.New("supplier required for purchase order")
    ErrMustBeConfirmed         = errors.New("order must be confirmed first")
    ErrOnlySalesCanShip        = errors.New("only sales orders can be shipped")
    ErrOnlyPurchaseCanReceive  = errors.New("only purchase orders can be received")
    ErrCannotClose             = errors.New("cannot close order in current state")
    ErrCancelReasonRequired    = errors.New("cancel reason required")
    ErrInvalidQuantity         = errors.New("quantity must be positive")
    ErrInvalidSKU              = errors.New("invalid SKU")
    ErrInvalidTaxType          = errors.New("invalid tax type")

    ErrInvoiceNotFound         = errors.New("invoice not found")
    ErrInvoiceAlreadyPaid      = errors.New("invoice already paid")
    ErrInvoiceCannotVoid       = errors.New("cannot void paid invoice")
    ErrInvoiceOverdue          = errors.New("invoice overdue")
    ErrPaymentAmountMismatch   = errors.New("payment amount mismatch")
    ErrPaymentAlreadyAllocated = errors.New("payment already fully allocated")
    ErrInsufficientStock       = errors.New("insufficient stock")
    ErrInventoryNegative       = errors.New("inventory cannot be negative")
)

// ---- CRM ----
var (
    ErrLeadNotFound           = errors.New("lead not found")
    ErrLeadAlreadyConverted   = errors.New("lead already converted")
    ErrInvalidScore           = errors.New("score must be 0-100")
    ErrInvalidLeadSource      = errors.New("invalid lead source")
    ErrTicketNotFound         = errors.New("ticket not found")
    ErrTicketAlreadyClosed    = errors.New("ticket already closed")
    ErrTicketAlreadyResolved  = errors.New("ticket already resolved")
    ErrInvalidPriority        = errors.New("invalid priority")
    ErrSLABreach              = errors.New("SLA breach")
)

// ---- Logistics ----
var (
    ErrShipmentNotFound       = errors.New("shipment not found")
    ErrTrackingAlreadyExists  = errors.New("tracking number already exists")
    ErrAlreadyDelivered       = errors.New("shipment already delivered")
    ErrCannotDispatch         = errors.New("cannot dispatch shipment in current state")
    ErrInvalidGPS             = errors.New("invalid GPS coordinates")
    ErrColdChainBreached      = errors.New("cold chain breached")
    ErrJobNotFound            = errors.New("job not found")
    ErrJobNotAssigned         = errors.New("job not assigned")
    ErrJobAlreadyStarted      = errors.New("job already started")
    ErrJobAlreadyCompleted    = errors.New("job already completed")
    ErrTechnicianUnavailable  = errors.New("technician unavailable")
    ErrChecklistIncomplete    = errors.New("checklist incomplete")
)

// ---- Report ----
var (
    ErrReportNotFound      = errors.New("report not found")
    ErrInvalidQuery        = errors.New("invalid SQL query")
    ErrInvalidCron         = errors.New("invalid cron expression")
    ErrReportRunning       = errors.New("report already running")
    ErrReportFailed        = errors.New("report execution failed")
)

// ---- Auth ----
var (
    ErrInvalidCredentials    = errors.New("invalid credentials")
    ErrTokenExpired          = errors.New("token expired")
    ErrTokenInvalid          = errors.New("token invalid")
    ErrRefreshTokenRevoked   = errors.New("refresh token revoked")
    ErrUserDisabled          = errors.New("user disabled")
    ErrEmailAlreadyExists    = errors.New("email already exists")
    ErrWeakPassword          = errors.New("password too weak")
    ErrUnauthorized          = errors.New("unauthorized")
    ErrForbidden             = errors.New("forbidden")
    ErrRateLimitExceeded     = errors.New("rate limit exceeded")
)
```

### I.2 Error Mapping (Domain → HTTP)

```go
// internal/shared/errors/mapping.go
package errors

import (
    "errors"
    "net/http"
    domainerrors "icmongolang/internal/modules/device/domain/errors"
)

func HTTPStatus(err error) int {
    switch {
    case errors.Is(err, domainerrors.ErrNotFound),
        errors.Is(err, domainerrors.ErrDeviceNotFound),
        errors.Is(err, domainerrors.ErrCustomerNotFound),
        errors.Is(err, domainerrors.ErrPackageNotFound),
        errors.Is(err, domainerrors.ErrOrderNotFound),
        errors.Is(err, domainerrors.ErrInvoiceNotFound),
        errors.Is(err, domainerrors.ErrLeadNotFound),
        errors.Is(err, domainerrors.ErrShipmentNotFound),
        errors.Is(err, domainerrors.ErrJobNotFound),
        errors.Is(err, domainerrors.ErrReportNotFound):
        return http.StatusNotFound

    case errors.Is(err, domainerrors.ErrAlreadyExists),
        errors.Is(err, domainerrors.ErrDeviceAlreadyExists),
        errors.Is(err, domainerrors.ErrCustomerAlreadyExists),
        errors.Is(err, domainerrors.ErrPackageAlreadyExists),
        errors.Is(err, domainerrors.ErrOrderAlreadyExists),
        errors.Is(err, domainerrors.ErrTrackingAlreadyExists),
        errors.Is(err, domainerrors.ErrSerialAlreadyExists),
        errors.Is(err, domainerrors.ErrEmailAlreadyExists):
        return http.StatusConflict

    case errors.Is(err, domainerrors.ErrUnauthorized),
        errors.Is(err, domainerrors.ErrInvalidCredentials),
        errors.Is(err, domainerrors.ErrTokenExpired),
        errors.Is(err, domainerrors.ErrTokenInvalid),
        errors.Is(err, domainerrors.ErrRefreshTokenRevoked):
        return http.StatusUnauthorized

    case errors.Is(err, domainerrors.ErrForbidden):
        return http.StatusForbidden

    case errors.Is(err, domainerrors.ErrRateLimitExceeded):
        return http.StatusTooManyRequests

    case errors.Is(err, domainerrors.ErrDeviceOffline),
        errors.Is(err, domainerrors.ErrDeviceFault):
        return http.StatusServiceUnavailable

    case errors.Is(err, domainerrors.ErrQuotaExceeded):
        return http.StatusPaymentRequired

    default:
        return http.StatusBadRequest // 400
    }
}
```

### I.3 Error Wrapping Pattern

```go
// ใช้ fmt.Errorf("%w", err) เพื่อ preserve sentinel
func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, in RegisterDeviceInput) (*Output, error) {
    existing, err := uc.repo.FindBySerial(ctx, in.TenantID, in.Serial)
    if err != nil && !errors.Is(err, domainerrors.ErrDeviceNotFound) {
        return nil, fmt.Errorf("find device by serial: %w", err)
    }
    if existing != nil {
        return nil, fmt.Errorf("serial %s: %w", in.Serial, domainerrors.ErrSerialAlreadyExists)
    }
    // ...
}
```

---

# 🎯 PART 5 — SUMMARY

### Domain Layer Deliverables

| Category | Count | Location |
|---|:-:|---|
| **Aggregates** | 18 | `domain/entity/` |
| **Value Objects** | 80+ | `domain/value_object/` |
| **Domain Events** | 60+ | `domain/event/` |
| **Domain Services** | 15+ | `domain/service/` |
| **Specifications** | 20+ | `domain/specification/` |
| **Repository Interfaces** | 25+ | `domain/repository/` |
| **Domain Errors** | 100+ | `domain/errors/` |
| **Factories** | 7 | `domain/factory/` |

### Layer Dependency Check

```
✅ Domain Layer imports:
   - stdlib (context, errors, time, strings, regexp)
   - github.com/google/uuid
   - github.com/shopspring/decimal

❌ Domain Layer NEVER imports:
   - gorm.io/gorm
   - github.com/gin-gonic/gin
   - github.com/IBM/sarama
   - github.com/redis/go-redis/v9
   - github.com/eclipse/paho.mqtt.golang
   - github.com/influxdata/influxdb-client-go/v2
   - github.com/elastic/go-elasticsearch/v8
```

---
