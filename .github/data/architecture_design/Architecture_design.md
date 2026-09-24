### SRS & Roadmap — แพลตฟอร์ม IoT (Smart Farm / Smart Building)
> **โปรเจกต์:** `icmongolang` — IoT Platform as a Service  
> **สถาปัตยกรรม:** Clean Architecture + Domain-Driven Design  
> **รูปแบบเอกสาร:** Software Requirements Specification (SRS) + Roadmap + Module Design  
> **อ้างอิง:** `Template_Module.md`, `Template_Module_MA.md`, `structure-map.md`

---
## สารบัญ

1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [โครงสร้าง Module](#2-โครงสร้าง-module)
3. [Bounded Contexts และ Context Map](#3-bounded-contexts-และ-context-map)
4. [Ubiquitous Language](#4-ubiquitous-language)
5. [Domain Layer](#5-domain-layer)
6. [Application Layer](#6-application-layer)
7. [Infrastructure Layer](#7-infrastructure-layer)
8. [Interface Layer](#8-interface-layer)
9. [Database Migrations](#9-database-migrations)
10. [System Flow](#10-system-flow)
11. [Workflow Diagram](#11-workflow-diagram)
12. [การติดตั้งและใช้งาน](#12-การติดตั้งและใช้งาน)
13. [Business Model](#13-business-model)
14. [Roadmap](#14-roadmap)
15. [ภาคผนวก](#15-ภาคผนวก)
16. [Prompt สำหรับการขยายระบบในอนาคต](#16-prompt-สำหรับการขยายระบบในอนาคต)
17. [การตั้งชื่อตาราง Database (Prefix)](#17-การตั้งชื่อตาราง-database-prefix)
18. [DDD Validation Checklist](#18-ddd-validation-checklist)

---

## 1. ภาพรวมระบบ

### 1.1 Vision

> **"แพลตฟอร์ม IoT แบบครบวงจรสำหรับ Smart Farm และ Smart Building ที่รวมการจัดการอุปกรณ์, AI/Automation, ERP, CRM และ Logistics ไว้ในระบบเดียว"**

### 1.2 กลุ่มเป้าหมาย

| กลุ่ม | ความต้องการหลัก |
| :--- | :--- |
| **เกษตรกร / Smart Farm** | ควบคุมโรงเรือน, sensor ดิน/น้ำ/อากาศ, AI วิเคราะห์ผลผลิต |
| **เจ้าของอาคาร / Smart Building** | จัดการ HVAC, ไฟฟ้า, ความปลอดภัย, พลังงาน |
| **System Integrator** | ติดตั้งและดูแลอุปกรณ์ให้ลูกค้าหลายราย |
| **Enterprise / OEM** | ต้องการ white-label platform |
| **ผู้ให้บริการ Logistics** | ติดตามสินค้าเกษตร/อาหารแช่เย็น |

### 1.3 Core Capabilities

```
┌─────────────────────────────────────────────────────────────┐
│                    IoT Platform (icmongolang)               │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ Customer │  │ Package  │  │   ERP    │  │   CRM    │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│  │  Device  │  │Logistics │  │  Report  │                   │
│  │ IoT/AI   │  │   IoT    │  │Analytics │                   │
│  └──────────┘  └──────────┘  └──────────┘                   │
├─────────────────────────────────────────────────────────────┤
│  Shared: auth · users · payment · notifier · mqtt · kafka   │
│          influxdb · elasticsearch · websocket · llm         │
└─────────────────────────────────────────────────────────────┘
```

### 1.4 Tech Stack

| Layer | Technology |
| :--- | :--- |
| Language | Go 1.21+ |
| Web Framework | Gin |
| ORM | GORM (PostgreSQL) |
| Cache | Redis |
| Message Broker | Kafka (IBM Sarama) |
| Time-series DB | InfluxDB |
| Search | Elasticsearch |
| IoT Protocol | MQTT |
| AI/LLM | Ollama / OpenAI-compatible |
| Realtime | WebSocket (Gorilla) |
| Observability | Prometheus + Grafana + Loki |

### 1.5 Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                     Interface Layer                          │
│  HTTP (Gin) │ WebSocket │ MQTT Gateway │ CLI │ Scheduler     │
├──────────────────────────────────────────────────────────────┤
│                    Application Layer                         │
│  Use Cases │ DTOs │ Orchestration │ Idempotency              │
├──────────────────────────────────────────────────────────────┤
│                      Domain Layer                            │
│  Entities │ Value Objects │ Repos (interface) │ Services     │
├──────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                       │
│  Postgres │ Redis │ Kafka │ InfluxDB │ ES │ MQTT │ LLM       │
└──────────────────────────────────────────────────────────────┘
```

---

## 2. โครงสร้าง Module

### 2.1 Module ทั้งหมด (33 เดิม + 7 ใหม่ = 40)

#### โมดูลที่มีอยู่แล้ว (33)
```
alarm        apimanager   auditlog     auth         batch
control      customer     dashboard    document     elasticsearch
email        flowengine   fullschedule i18n         influxdb
iot          items        job          kafka        mqtt
notifier     payment      pdpa         purchaseorder queue
quotation    realtime     report       settings     users
vectordata   websocket    wos
```

#### โมดูลใหม่สำหรับ IoT Platform (7)
```
package      erp          crm          device       logistics
farm         building
```

> **หมายเหตุ:** `device` และ `farm`/`building` อาจรวมเป็น module เดียวชื่อ `device` โดยแยก sub-domain ภายใน

### 2.2 โครงสร้างมาตรฐาน (Clean Architecture + DDD)

```
internal/modules/{{module_name}}/
│
├── domain/                                    # 🏛️ DOMAIN LAYER
│   ├── entity/
│   │   ├── {{aggregate}}.go                  # Aggregate Root
│   │   └── {{child_entity}}.go               # Child Entities
│   ├── value_object/
│   │   ├── status.go                         # Status VO
│   │   ├── {{vo_name}}.go                    # Immutable VOs
│   │   └── id.go                             # Typed IDs
│   ├── repository/
│   │   ├── {{aggregate}}_repository.go       # Interface เท่านั้น
│   │   └── read_model_repository.go          # CQRS read side
│   ├── service/
│   │   ├── {{domain}}_service.go             # Stateless domain logic
│   │   └── {{port}}_port.go                  # Outbound ports
│   ├── event/
│   │   └── {{aggregate}}_events.go           # Domain events
│   └── errors/
│       └── errors.go                         # Sentinel errors
│
├── application/                                # 🎯 APPLICATION LAYER
│   ├── {{verb}}_{{aggregate}}.go             # 1 use case / 1 file
│   ├── dto.go                                # Request/Response DTOs
│   ├── mappers.go                            # Entity ↔ DTO
│   └── query/                                # CQRS read models
│       └── {{query}}_handler.go
│
├── infrastructure/                             # 🔧 INFRASTRUCTURE LAYER
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── {{aggregate}}_repo_impl.go
│   │   │   ├── read_model_repo_impl.go
│   │   │   └── models.go                     # GORM models (prefix)
│   │   ├── redis/
│   │   │   └── {{aggregate}}_cache.go
│   │   └── influxdb/
│   │       └── {{metric}}_repo.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── consumers/
│   │       ├── {{topic}}_consumer.go
│   │       └── consumer_group.go
│   ├── search/elasticsearch/
│   │   └── {{aggregate}}_indexer.go
│   ├── iot/
│   │   ├── mqtt_subscriber.go
│   │   ├── mqtt_publisher.go
│   │   └── protocol/
│   │       ├── modbus.go
│   │       ├── lorawan.go
│   │       └── zigbee.go
│   ├── services/
│   │   ├── llm/ollama.go
│   │   ├── email/smtp.go
│   │   ├── sms/twilio.go
│   │   ├── payment/stripe.go
│   │   └── blockchain/ethereum.go
│   └── scheduler/
│       └── {{job}}_job.go
│
├── interfaces/                                 # 🌐 INTERFACE LAYER
│   ├── http/
│   │   ├── {{aggregate}}_handler.go
│   │   ├── routes.go
│   │   └── dto.go
│   ├── websocket/
│   │   └── hub.go
│   ├── mqtt/
│   │   └── gateway.go
│   └── middleware/
│       └── {{module}}_middleware.go
│
└── module.go                                   # Composition Root (Wire-up)
```

### 2.3 Entry Points

```
cmd/
├── api/main.go                # REST + WebSocket
├── scheduler/main.go          # Cron jobs
├── migrate/main.go            # DB migrations
├── mqtt-gateway/main.go       # MQTT bridge
└── workers/
    ├── device-telemetry/main.go
    ├── ai-inference/main.go
    ├── notification/main.go
    └── erp-sync/main.go
```

---

## 3. Bounded Contexts และ Context Map

### 3.1 Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────┐
│                     IoT Platform Contexts                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌────────────┐         ┌────────────┐         ┌────────────┐  │
│  │  Customer  │◄────────│  Package   │────────►│  Payment   │  │
│  │  Context   │  U/D    │  Context   │   PL    │  Context   │  │
│  └─────┬──────┘         └─────┬──────┘         └────────────┘  │
│        │                      │                                 │
│        │ U/D                  │ U/D                             │
│        ▼                      ▼                                 │
│  ┌────────────┐         ┌────────────┐         ┌────────────┐  │
│  │    CRM     │◄────────│   Device   │────────►│    ERP     │  │
│  │  Context   │   CF    │  Context   │   ACL   │  Context   │  │
│  └─────┬──────┘         └─────┬──────┘         └─────┬──────┘  │
│        │                      │                      │         │
│        │                      │                      │         │
│        ▼                      ▼                      ▼         │
│  ┌────────────┐         ┌────────────┐         ┌────────────┐  │
│  │ Logistics  │◄────────│  Report    │────────►│   Alarm    │  │
│  │  Context   │  U/D    │  Context   │   PL    │  Context   │  │
│  └────────────┘         └────────────┘         └────────────┘  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 Context Map Relationships

| Upstream (U) | Downstream (D) | Pattern | หมายเหตุ |
| :--- | :--- | :--- | :--- |
| Customer | CRM | Customer-Supplier | CRM ใช้ข้อมูลลูกค้า |
| Customer | Package | Customer-Supplier | ลูกค้าเลือกแพ็กเกจ |
| Package | Payment | Partnership | แพ็กเกจกำหนดราคา |
| Device | ERP | ACL (Anti-Corruption) | ERP ดึงข้อมูลอุปกรณ์ |
| Device | Logistics | Conformist | Logistics ใช้ device tracking |
| Device | Report | Published Language | Report อ่าน telemetry |
| CRM | Report | Shared Kernel | ใช้ DTO ร่วม |
| Logistics | Report | Published Language | ข้อมูลการขนส่ง |
| Alarm | Notifier | Partnership | แจ้งเตือนผ่าน notifier |

### 3.3 Context Integration Patterns

```go
// ตัวอย่าง: ERP ดึงข้อมูล Device ผ่าน ACL
// infrastructure/erp/client/device_acl.go
package client

type DeviceACL struct {
    deviceAPI DeviceAPIClient
}

func (a *DeviceACL) GetDeviceForERP(ctx context.Context, id uuid.UUID) (*ERPDevice, error) {
    d, err := a.deviceAPI.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    // แปลง Device → ERP-specific model
    return &ERPDevice{
        Code:     d.SerialNumber,
        Name:     d.Name,
        Cost:     d.PurchasePrice,
        Location: d.SiteID,
    }, nil
}
```

---

## 4. Ubiquitous Language

| คำศัพท์ | ความหมาย | Context |
| :--- | :--- | :--- |
| **Site** | พื้นที่ติดตั้ง (ฟาร์ม/อาคาร/โรงงาน) | Device |
| **Zone** | พื้นที่ย่อยภายใน Site (โรงเรือน/ชั้น) | Device |
| **Device** | อุปกรณ์ IoT (Sensor, Actuator, Gateway) | Device |
| **Gateway** | อุปกรณ์รับส่งข้อมูลจาก Device ไป Cloud | Device |
| **Telemetry** | ข้อมูลที่วัดได้จาก Sensor (อุณหภูมิ, ความชื้น) | Device |
| **Command** | คำสั่งควบคุม Actuator (เปิด/ปิด, setpoint) | Device |
| **Automation Rule** | กฎอัตโนมัติ (if temp > 30 then turn on fan) | Device |
| **Customer** | ลูกค้าที่ซื้อบริการ | Customer |
| **Tenant** | องค์กร/บัญชีที่แยกข้อมูล | Customer |
| **Package** | แพ็กเกจบริการ (Free, Pro, Enterprise) | Package |
| **Subscription** | การสมัครใช้แพ็กเกจ | Package |
| **Quota** | โควตาการใช้งาน (จำนวน device, storage) | Package |
| **Order** | คำสั่งซื้อสินค้า/บริการ | ERP |
| **Invoice** | ใบแจ้งหนี้ | ERP |
| **Stock** | สินค้าคงคลัง | ERP |
| **Lead** | ลูกค้าเป้าหมาย | CRM |
| **Opportunity** | โอกาสในการขาย | CRM |
| **Ticket** | ใบแจ้งปัญหาบริการ | CRM |
| **Shipment** | การจัดส่ง | Logistics |
| **Route** | เส้นทางขนส่ง | Logistics |
| **Cold Chain** | ห่วงโซ่ความเย็น | Logistics |
| **Report** | รายงาน/แดชบอร์ด | Report |
| **Metric** | ตัวชี้วัด (KPI) | Report |
| **Alert** | การแจ้งเตือนเมื่อเข้าเงื่อนไข | Alarm |

---

## 5. Domain Layer

### 5.1 Entities (Aggregates)

#### 5.1.1 Customer Context — Aggregate `Customer`

```go
// internal/modules/customer/domain/entity/customer.go
package entity

import (
	"time"
	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

// Customer – Aggregate Root
type Customer struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Code        string
	Name        string
	Type        valueobject.CustomerType // PERSONAL, BUSINESS, GOVERNMENT
	Status      valueobject.CustomerStatus
	TaxID       string
	Address     Address // Value Object
	Contacts    []Contact
	Sites       []uuid.UUID // reference to Device context
	PackageID   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCustomer(tenantID uuid.UUID, code, name string, ctype valueobject.CustomerType) (*Customer, error) {
	if code == "" || name == "" {
		return nil, domainerrors.ErrInvalidCustomer
	}
	now := time.Now()
	return &Customer{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		Type:      ctype,
		Status:    valueobject.CustomerStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Customer) Activate() error {
	if c.Status == valueobject.CustomerStatusActive {
		return domainerrors.ErrCustomerAlreadyActive
	}
	c.Status = valueobject.CustomerStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Customer) Suspend(reason string) error {
	if c.Status == valueobject.CustomerStatusSuspended {
		return domainerrors.ErrCustomerAlreadySuspended
	}
	c.Status = valueobject.CustomerStatusSuspended
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Customer) AssignPackage(pkgID uuid.UUID) error {
	if pkgID == uuid.Nil {
		return domainerrors.ErrInvalidPackage
	}
	c.PackageID = pkgID
	c.UpdatedAt = time.Now()
	return nil
}

// Address – Value Object
type Address struct {
	Line1    string
	Line2    string
	District string
	Amphoe   string
	Province string
	Postcode string
	Country  string
}

// Contact – Entity within Customer aggregate
type Contact struct {
	ID        uuid.UUID
	Name      string
	Phone     string
	Email     string
	Position  string
	IsPrimary bool
}
```

#### 5.1.2 Device Context — Aggregate `Device`

```go
// internal/modules/device/domain/entity/device.go
package entity

import (
	"time"
	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/device/domain/errors"
	valueobject "icmongolang/internal/modules/device/domain/value_object"
)

// Device – Aggregate Root
type Device struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	SiteID        uuid.UUID
	ZoneID        uuid.UUID
	SerialNumber  string
	Name          string
	Type          valueobject.DeviceType   // SENSOR, ACTUATOR, GATEWAY, CAMERA
	Protocol      valueobject.Protocol      // MQTT, MODBUS, LORAWAN, ZIGBEE
	Status        valueobject.DeviceStatus  // ONLINE, OFFLINE, ERROR, MAINTENANCE
	Firmware      string
	Config        map[string]interface{}
	LastSeenAt    *time.Time
	InstalledAt   *time.Time
	Metadata      DeviceMetadata
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewDevice(tenantID, siteID uuid.UUID, serial, name string, dtype valueobject.DeviceType, proto valueobject.Protocol) (*Device, error) {
	if serial == "" {
		return nil, domainerrors.ErrInvalidSerial
	}
	if !dtype.IsValid() {
		return nil, domainerrors.ErrInvalidDeviceType
	}
	now := time.Now()
	return &Device{
		ID:           uuid.New(),
		TenantID:     tenantID,
		SiteID:       siteID,
		SerialNumber: serial,
		Name:         name,
		Type:         dtype,
		Protocol:     proto,
		Status:       valueobject.DeviceStatusOffline,
		Config:       make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (d *Device) MarkOnline() {
	now := time.Now()
	d.Status = valueobject.DeviceStatusOnline
	d.LastSeenAt = &now
	d.UpdatedAt = now
}

func (d *Device) MarkOffline() {
	d.Status = valueobject.DeviceStatusOffline
	d.UpdatedAt = time.Now()
}

func (d *Device) SetError(msg string) {
	d.Status = valueobject.DeviceStatusError
	d.Config["last_error"] = msg
	d.UpdatedAt = time.Now()
}

func (d *Device) UpdateFirmware(version string) error {
	if version == "" {
		return domainerrors.ErrInvalidFirmware
	}
	d.Firmware = version
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Device) MoveTo(siteID, zoneID uuid.UUID) error {
	if siteID == uuid.Nil {
		return domainerrors.ErrInvalidSite
	}
	d.SiteID = siteID
	d.ZoneID = zoneID
	d.UpdatedAt = time.Now()
	return nil
}

type DeviceMetadata struct {
	Manufacturer string
	Model        string
	PurchaseDate *time.Time
	WarrantyEnd  *time.Time
	PurchasePrice float64
}

// Telemetry – Time-series entity (persisted to InfluxDB)
type Telemetry struct {
	DeviceID   uuid.UUID
	Metric     string
	Value      float64
	Unit       string
	Tags       map[string]string
	Timestamp  time.Time
}

// Command – Entity for device control
type Command struct {
	ID         uuid.UUID
	DeviceID   uuid.UUID
	Type       string
	Payload    map[string]interface{}
	Status     valueobject.CommandStatus // PENDING, SENT, ACKED, FAILED
	IssuedBy   uuid.UUID
	IssuedAt   time.Time
	ExecutedAt *time.Time
}
```

#### 5.1.3 Package Context — Aggregate `Package` + `Subscription`

```go
// internal/modules/package/domain/entity/package.go
package entity

import (
	"time"
	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type Package struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
	Tier        valueobject.PackageTier // FREE, BASIC, PRO, ENTERPRISE
	Price       Money
	BillingCycle valueobject.BillingCycle // MONTHLY, YEARLY
	Quotas      Quotas
	Features    []Feature
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Quotas struct {
	MaxDevices      int
	MaxSites        int
	MaxUsers        int
	StorageGB       int
	RetentionDays   int
	APICallsPerDay  int
}

type Feature struct {
	Code        string
	Name        string
	Description string
	Enabled     bool
}

type Money struct {
	Amount   float64
	Currency string // THB, USD
}

func NewPackage(code, name string, tier valueobject.PackageTier, price Money, cycle valueobject.BillingCycle) (*Package, error) {
	if code == "" || name == "" {
		return nil, domainerrors.ErrInvalidPackage
	}
	if !tier.IsValid() {
		return nil, domainerrors.ErrInvalidTier
	}
	if price.Amount < 0 {
		return nil, domainerrors.ErrInvalidPrice
	}
	now := time.Now()
	return &Package{
		ID:           uuid.New(),
		Code:         code,
		Name:         name,
		Tier:         tier,
		Price:        price,
		BillingCycle: cycle,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (p *Package) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// Subscription – Aggregate Root
type Subscription struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	PackageID   uuid.UUID
	Status      valueobject.SubscriptionStatus // ACTIVE, PAST_DUE, CANCELLED, EXPIRED
	StartedAt   time.Time
	ExpiresAt   time.Time
	AutoRenew   bool
	Usage       Quotas
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSubscription(customerID, packageID uuid.UUID, cycle valueobject.BillingCycle) *Subscription {
	now := time.Now()
	expires := now.AddDate(0, 1, 0)
	if cycle == valueobject.BillingYearly {
		expires = now.AddDate(1, 0, 0)
	}
	return &Subscription{
		ID:         uuid.New(),
		CustomerID: customerID,
		PackageID:  packageID,
		Status:     valueobject.SubscriptionStatusActive,
		StartedAt:  now,
		ExpiresAt:  expires,
		AutoRenew:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (s *Subscription) Upgrade(pkgID uuid.UUID, newExpiry time.Time) error {
	if s.Status == valueobject.SubscriptionStatusCancelled {
		return domainerrors.ErrSubscriptionCancelled
	}
	s.PackageID = pkgID
	s.ExpiresAt = newExpiry
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Subscription) Cancel() {
	s.Status = valueobject.SubscriptionStatusCancelled
	s.AutoRenew = false
	s.UpdatedAt = time.Now()
}
```

#### 5.1.4 ERP Context — Aggregate `Order`, `Invoice`, `StockItem`

```go
// internal/modules/erp/domain/entity/order.go
package entity

type Order struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	CustomerID  uuid.UUID
	OrderNumber string
	Items       []OrderItem
	Total       Money
	Status      valueobject.OrderStatus // DRAFT, CONFIRMED, SHIPPED, DELIVERED, CANCELLED
	OrderedAt   time.Time
	DeliveredAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	ProductID  uuid.UUID
	SKU        string
	Name       string
	Quantity   int
	UnitPrice  Money
	Subtotal   Money
}

type Invoice struct {
	ID          uuid.UUID
	OrderID     uuid.UUID
	CustomerID  uuid.UUID
	Number      string
	Amount      Money
	TaxAmount   Money
	Total       Money
	Status      valueobject.InvoiceStatus // DRAFT, ISSUED, PAID, OVERDUE, CANCELLED
	IssuedAt    time.Time
	DueDate     time.Time
	PaidAt      *time.Time
}

type StockItem struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	WarehouseID uuid.UUID
	SKU         string
	Quantity    int
	Reserved    int
	Available   int
	ReorderPoint int
	UpdatedAt   time.Time
}
```

#### 5.1.5 CRM Context — Aggregate `Lead`, `Opportunity`, `Ticket`

```go
// internal/modules/crm/domain/entity/lead.go
package entity

type Lead struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Email       string
	Phone       string
	Company     string
	Source      valueobject.LeadSource // WEB, REFERRAL, CAMPAIGN, COLD_CALL
	Status      valueobject.LeadStatus // NEW, CONTACTED, QUALIFIED, LOST, CONVERTED
	Score       int
	AssignedTo  *uuid.UUID // user ID
	ConvertedTo *uuid.UUID // customer ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (l *Lead) Qualify(score int) error {
	if score < 0 || score > 100 {
		return domainerrors.ErrInvalidScore
	}
	l.Score = score
	l.Status = valueobject.LeadStatusQualified
	l.UpdatedAt = time.Now()
	return nil
}

func (l *Lead) Convert(customerID uuid.UUID) error {
	if l.Status == valueobject.LeadStatusConverted {
		return domainerrors.ErrLeadAlreadyConverted
	}
	l.Status = valueobject.LeadStatusConverted
	l.ConvertedTo = &customerID
	l.UpdatedAt = time.Now()
	return nil
}

type Opportunity struct {
	ID          uuid.UUID
	LeadID      uuid.UUID
	CustomerID  *uuid.UUID
	Title       string
	Amount      Money
	Probability int // 0-100
	Stage       valueobject.OpportunityStage // PROSPECTING, PROPOSAL, NEGOTIATION, WON, LOST
	ExpectedCloseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Ticket struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Subject     string
	Description string
	Priority    valueobject.TicketPriority // LOW, MEDIUM, HIGH, URGENT
	Status      valueobject.TicketStatus    // OPEN, IN_PROGRESS, RESOLVED, CLOSED
	AssignedTo  *uuid.UUID
	SLA_DueAt   time.Time
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

#### 5.1.6 Logistics Context — Aggregate `Shipment`, `Route`

```go
// internal/modules/logistics/domain/entity/shipment.go
package entity

type Shipment struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	TrackingNumber string
	OrderID        uuid.UUID
	CustomerID     uuid.UUID
	Status         valueobject.ShipmentStatus // PENDING, PICKED, IN_TRANSIT, DELIVERED, FAILED
	Items          []ShipmentItem
	Origin         Location
	Destination    Location
	Carrier        string
	VehicleID      *uuid.UUID
	DriverID       *uuid.UUID
	EstimatedAt    time.Time
	DeliveredAt    *time.Time
	Temperature    *TemperatureLog // สำหรับ cold chain
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Location struct {
	Address  string
	Lat      float64
	Lng      float64
	Postcode string
}

type TemperatureLog struct {
	Min        float64
	Max        float64
	Current    float64
	Readings   []TemperatureReading
	IsBreached bool
}

type TemperatureReading struct {
	Value     float64
	Timestamp time.Time
	DeviceID  uuid.UUID
}

func (s *Shipment) UpdateLocation(lat, lng float64) error {
	// validate
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Shipment) MarkDelivered() error {
	if s.Status == valueobject.ShipmentStatusDelivered {
		return domainerrors.ErrAlreadyDelivered
	}
	now := time.Now()
	s.Status = valueobject.ShipmentStatusDelivered
	s.DeliveredAt = &now
	s.UpdatedAt = now
	return nil
}

func (s *Shipment) RecordTemperature(reading TemperatureReading) {
	if s.Temperature == nil {
		s.Temperature = &TemperatureLog{Min: reading.Value, Max: reading.Value}
	}
	if reading.Value < s.Temperature.Min {
		s.Temperature.Min = reading.Value
	}
	if reading.Value > s.Temperature.Max {
		s.Temperature.Max = reading.Value
	}
	s.Temperature.Current = reading.Value
	s.Temperature.Readings = append(s.Temperature.Readings, reading)
	if reading.Value < 2 || reading.Value > 8 { // cold chain 2-8°C
		s.Temperature.IsBreached = true
	}
}

type Route struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Waypoints   []Waypoint
	DistanceKM  float64
	EstimatedMin int
	CreatedAt   time.Time
}

type Waypoint struct {
	Seq     int
	Lat     float64
	Lng     float64
	Address string
	ETA     time.Time
}
```

#### 5.1.7 Report Context — Aggregate `ReportDefinition`

```go
// internal/modules/report/domain/entity/report.go
package entity

type ReportDefinition struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Code        string
	Name        string
	Category    valueobject.ReportCategory // OPERATION, FINANCE, ENERGY, AGRICULTURE
	Type        valueobject.ReportType     // TABLE, CHART, DASHBOARD, EXPORT
	Query       string
	Parameters  []Parameter
	Schedule    *Schedule
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Parameter struct {
	Name     string
	Type     string // string, int, date, uuid
	Required bool
	Default  interface{}
}

type Schedule struct {
	Cron     string
	Format   string // PDF, EXCEL, CSV
	EmailTo  []string
	Enabled  bool
}

type ReportRun struct {
	ID         uuid.UUID
	ReportID   uuid.UUID
	RunBy      uuid.UUID
	Status     valueobject.ReportRunStatus // PENDING, RUNNING, SUCCESS, FAILED
	Parameters map[string]interface{}
	ResultURL  string
	ErrorMsg   string
	StartedAt  time.Time
	FinishedAt *time.Time
}
```

### 5.2 Value Objects

```go
// internal/modules/device/domain/value_object/device_status.go
package valueobject

type DeviceStatus string

const (
	DeviceStatusOnline      DeviceStatus = "ONLINE"
	DeviceStatusOffline     DeviceStatus = "OFFLINE"
	DeviceStatusError       DeviceStatus = "ERROR"
	DeviceStatusMaintenance DeviceStatus = "MAINTENANCE"
)

func (s DeviceStatus) IsValid() bool {
	switch s {
	case DeviceStatusOnline, DeviceStatusOffline, DeviceStatusError, DeviceStatusMaintenance:
		return true
	}
	return false
}

func (s DeviceStatus) IsOperational() bool {
	return s == DeviceStatusOnline
}

// DeviceType
type DeviceType string

const (
	DeviceTypeSensor   DeviceType = "SENSOR"
	DeviceTypeActuator DeviceType = "ACTUATOR"
	DeviceTypeGateway  DeviceType = "GATEWAY"
	DeviceTypeCamera   DeviceType = "CAMERA"
	DeviceTypeMeter    DeviceType = "METER"
)

func (t DeviceType) IsValid() bool {
	switch t {
	case DeviceTypeSensor, DeviceTypeActuator, DeviceTypeGateway, DeviceTypeCamera, DeviceTypeMeter:
		return true
	}
	return false
}

// Protocol
type Protocol string

const (
	ProtocolMQTT    Protocol = "MQTT"
	ProtocolModbus  Protocol = "MODBUS"
	ProtocolLoRaWAN Protocol = "LORAWAN"
	ProtocolZigbee  Protocol = "ZIGBEE"
	ProtocolHTTP    Protocol = "HTTP"
)

// CommandStatus
type CommandStatus string

const (
	CommandStatusPending CommandStatus = "PENDING"
	CommandStatusSent    CommandStatus = "SENT"
	CommandStatusAcked   CommandStatus = "ACKED"
	CommandStatusFailed  CommandStatus = "FAILED"
)
```

```go
// internal/modules/package/domain/value_object/tier.go
package valueobject

type PackageTier string

const (
	TierFree       PackageTier = "FREE"
	TierBasic      PackageTier = "BASIC"
	TierPro        PackageTier = "PRO"
	TierEnterprise PackageTier = "ENTERPRISE"
)

func (t PackageTier) IsValid() bool {
	switch t {
	case TierFree, TierBasic, TierPro, TierEnterprise:
		return true
	}
	return false
}

type BillingCycle string

const (
	BillingMonthly BillingCycle = "MONTHLY"
	BillingYearly  BillingCycle = "YEARLY"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusPastDue   SubscriptionStatus = "PAST_DUE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
)
```

### 5.3 Repository Interfaces

```go
// internal/modules/device/domain/repository/device_repository.go
package repository

import (
	"context"
	"github.com/google/uuid"
	"icmongolang/internal/modules/device/domain/entity"
	"icmongolang/internal/modules/device/domain/value_object"
)

type DeviceRepository interface {
	Save(ctx context.Context, d *entity.Device) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	FindBySerial(ctx context.Context, serial string) (*entity.Device, error)
	FindBySite(ctx context.Context, siteID uuid.UUID) ([]entity.Device, error)
	FindByStatus(ctx context.Context, status valueobject.DeviceStatus) ([]entity.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TelemetryRepository interface {
	Write(ctx context.Context, t *entity.Telemetry) error
	WriteBatch(ctx context.Context, batch []entity.Telemetry) error
	Query(ctx context.Context, deviceID uuid.UUID, metric string, from, to time.Time) ([]entity.Telemetry, error)
	Latest(ctx context.Context, deviceID uuid.UUID, metric string) (*entity.Telemetry, error)
}

type CommandRepository interface {
	Save(ctx context.Context, c *entity.Command) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Command, error)
	FindPendingByDevice(ctx context.Context, deviceID uuid.UUID) ([]entity.Command, error)
}
```

### 5.4 Domain Services

```go
// internal/modules/device/domain/service/automation_service.go
package service

import (
	"context"
	"icmongolang/internal/modules/device/domain/entity"
)

// AutomationService – stateless logic ที่ไม่ผูกกับ entity ตัวใดตัวหนึ่ง
type AutomationService struct{}

// EvaluateRule – ประเมินกฎอัตโนมัติ
func (s *AutomationService) EvaluateRule(rule AutomationRule, t entity.Telemetry) (Action, bool, error) {
	// logic
	return Action{}, false, nil
}

type AutomationRule struct {
	ID         string
	Conditions []Condition
	Actions    []Action
	Enabled    bool
}

type Condition struct {
	Metric   string
	Operator string // >, <, ==, >=, <=
	Value    float64
}

type Action struct {
	DeviceID uuid.UUID
	Command  string
	Payload  map[string]interface{}
}
```

### 5.5 Domain Errors

```go
// internal/modules/device/domain/errors/errors.go
package domainerrors

import "errors"

var (
	ErrDeviceNotFound        = errors.New("device not found")
	ErrDeviceAlreadyExists   = errors.New("device already exists")
	ErrInvalidSerial         = errors.New("invalid serial number")
	ErrInvalidDeviceType     = errors.New("invalid device type")
	ErrInvalidFirmware       = errors.New("invalid firmware version")
	ErrInvalidSite           = errors.New("invalid site")
	ErrDeviceOffline         = errors.New("device is offline")
	ErrCommandNotAllowed     = errors.New("command not allowed for this device")
)
```

### 5.6 Invariants

| Aggregate | Invariant |
| :--- | :--- |
| Customer | ต้องมี Code ไม่ซ้ำ, TaxID ต้องถูกต้องถ้าเป็น BUSINESS |
| Device | SerialNumber ไม่ซ้ำ, ต้องมี SiteID, Firmware version ต้อง semver |
| Package | Price ≥ 0, Quotas ทุกตัว ≥ 0 |
| Subscription | ExpiresAt > StartedAt, ไม่สามารถ upgrade ถ้า Cancelled |
| Order | Total = sum(Items.Subtotal), ต้องมีอย่างน้อย 1 item |
| Invoice | Total = Amount + TaxAmount, DueDate > IssuedAt |
| Lead | Score 0-100, Convert ได้ครั้งเดียว |
| Shipment | DeliveredAt ≥ CreatedAt, Cold chain temp ต้อง 2-8°C |
| ReportDefinition | Query ต้อง valid SQL, Cron ต้อง valid expression |

---

## 6. Application Layer

### 6.1 Use Cases

#### 6.1.1 Device — Register Device

```go
// internal/modules/device/application/register_device.go
package application

import (
	"context"
	"github.com/google/uuid"
	"icmongolang/internal/modules/device/domain/entity"
	domainerrors "icmongolang/internal/modules/device/domain/errors"
	"icmongolang/internal/modules/device/domain/repository"
	valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type RegisterDeviceUseCase struct {
	repo        repository.DeviceRepository
	auditRepo   repository.AuditRepository
	producer    EventProducer
}

func NewRegisterDeviceUseCase(
	repo repository.DeviceRepository,
	auditRepo repository.AuditRepository,
	producer EventProducer,
) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{repo: repo, auditRepo: auditRepo, producer: producer}
}

type RegisterDeviceInput struct {
	TenantID     uuid.UUID
	SiteID       uuid.UUID
	ZoneID       uuid.UUID
	SerialNumber string
	Name         string
	Type         valueobject.DeviceType
	Protocol     valueobject.Protocol
	UserID       uuid.UUID
	IPAddress    string
}

type RegisterDeviceOutput struct {
	DeviceID uuid.UUID
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, input RegisterDeviceInput) (*RegisterDeviceOutput, error) {
	// 1. Validate
	if input.TenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if input.SerialNumber == "" {
		return nil, domainerrors.ErrInvalidSerial
	}

	// 2. Check duplicate
	if existing, _ := uc.repo.FindBySerial(ctx, input.SerialNumber); existing != nil {
		return nil, domainerrors.ErrDeviceAlreadyExists
	}

	// 3. Create aggregate
	device, err := entity.NewDevice(input.TenantID, input.SiteID, input.SerialNumber, input.Name, input.Type, input.Protocol)
	if err != nil {
		return nil, err
	}
	device.ZoneID = input.ZoneID

	// 4. Persist
	if err := uc.repo.Save(ctx, device); err != nil {
		return nil, err
	}

	// 5. Side effects
	_ = uc.auditRepo.Save(ctx, entity.NewAuditTrail(&input.UserID, "DEVICE_REGISTERED", map[string]interface{}{
		"device_id": device.ID.String(),
		"serial":    device.SerialNumber,
		"ip":        input.IPAddress,
	}))

	_ = uc.producer.PublishDeviceRegistered(ctx, device)

	return &RegisterDeviceOutput{DeviceID: device.ID}, nil
}
```

#### 6.1.2 Device — Ingest Telemetry

```go
// internal/modules/device/application/ingest_telemetry.go
package application

import (
	"context"
	"time"
	"github.com/google/uuid"
	"icmongolang/internal/modules/device/domain/entity"
	"icmongolang/internal/modules/device/domain/repository"
	"icmongolang/internal/modules/device/domain/service"
)

type IngestTelemetryUseCase struct {
	deviceRepo    repository.DeviceRepository
	telemetryRepo repository.TelemetryRepository
	automation    *service.AutomationService
	wsHub         WSHub
	producer      EventProducer
}

type IngestTelemetryInput struct {
	DeviceID  uuid.UUID
	Metrics   []MetricReading
	Timestamp time.Time
}

type MetricReading struct {
	Metric string
	Value  float64
	Unit   string
}

func (uc *IngestTelemetryUseCase) Execute(ctx context.Context, input IngestTelemetryInput) error {
	// 1. Verify device exists
	device, err := uc.deviceRepo.FindByID(ctx, input.DeviceID)
	if err != nil {
		return err
	}

	// 2. Update device status if needed
	if !device.Status.IsOperational() {
		device.MarkOnline()
		_ = uc.deviceRepo.Save(ctx, device)
	}

	// 3. Write telemetry to InfluxDB
	batch := make([]entity.Telemetry, 0, len(input.Metrics))
	for _, m := range input.Metrics {
		batch = append(batch, entity.Telemetry{
			DeviceID:  input.DeviceID,
			Metric:    m.Metric,
			Value:     m.Value,
			Unit:      m.Unit,
			Timestamp: input.Timestamp,
		})
	}
	if err := uc.telemetryRepo.WriteBatch(ctx, batch); err != nil {
		return err
	}

	// 4. Evaluate automation rules
	actions, _ := uc.automation.EvaluateAll(device, batch)

	// 5. Broadcast via WebSocket
	for _, t := range batch {
		uc.wsHub.BroadcastToTenant(device.TenantID.String(), WSEvent{
			Type: "telemetry",
			Data: t,
		})
	}

	// 6. Publish event for async processing (AI, alarm)
	_ = uc.producer.PublishTelemetryIngested(ctx, device, batch)

	// 7. Execute automation actions
	for _, action := range actions {
		_ = uc.producer.PublishCommandIssued(ctx, action)
	}

	return nil
}
```

#### 6.1.3 Package — Subscribe

```go
// internal/modules/package/application/subscribe.go
package application

type SubscribeUseCase struct {
	subRepo     repository.SubscriptionRepository
	pkgRepo     repository.PackageRepository
	paymentRepo repository.PaymentRepository
	producer    EventProducer
}

type SubscribeInput struct {
	CustomerID uuid.UUID
	PackageID  uuid.UUID
	Cycle      valueobject.BillingCycle
	PaymentRef string
}

type SubscribeOutput struct {
	SubscriptionID uuid.UUID
	PaymentURL     string
}

func (uc *SubscribeUseCase) Execute(ctx context.Context, input SubscribeInput) (*SubscribeOutput, error) {
	// 1. Load package
	pkg, err := uc.pkgRepo.FindByID(ctx, input.PackageID)
	if err != nil {
		return nil, err
	}
	if !pkg.IsActive {
		return nil, domainerrors.ErrPackageInactive
	}

	// 2. Create subscription
	sub := entity.NewSubscription(input.CustomerID, pkg.ID, input.Cycle)
	if err := uc.subRepo.Save(ctx, sub); err != nil {
		return nil, err
	}

	// 3. Create payment
	payment := entity.NewPayment(input.CustomerID, pkg.Price, "SUBSCRIPTION", sub.ID)
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return nil, err
	}

	// 4. Publish event
	_ = uc.producer.PublishSubscriptionCreated(ctx, sub)

	return &SubscribeOutput{
		SubscriptionID: sub.ID,
		PaymentURL:     payment.CheckoutURL,
	}, nil
}
```

### 6.2 DTOs

```go
// internal/modules/device/application/dto.go
package application

import "time"

type DeviceResponse struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	SiteID       string                 `json:"site_id"`
	ZoneID       string                 `json:"zone_id,omitempty"`
	SerialNumber string                 `json:"serial_number"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Protocol     string                 `json:"protocol"`
	Status       string                 `json:"status"`
	Firmware     string                 `json:"firmware"`
	LastSeenAt   *time.Time             `json:"last_seen_at,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type TelemetryQuery struct {
	DeviceID string    `form:"device_id" binding:"required"`
	Metric   string    `form:"metric"`
	From     time.Time `form:"from" binding:"required"`
	To       time.Time `form:"to" binding:"required"`
	Interval string    `form:"interval"` // 1m, 5m, 1h
}
```

---

## 7. Infrastructure Layer

### 7.1 Repository Implementations

```go
// internal/modules/device/infrastructure/persistence/postgres/device_repo_impl.go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/device/domain/entity"
	domainerrors "icmongolang/internal/modules/device/domain/errors"
	valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type DeviceModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index"`
	SiteID       uuid.UUID `gorm:"type:uuid;not null;index"`
	ZoneID       uuid.UUID `gorm:"type:uuid;index"`
	SerialNumber string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Type         string    `gorm:"type:varchar(50);not null;index"`
	Protocol     string    `gorm:"type:varchar(50);not null"`
	Status       string    `gorm:"type:varchar(30);not null;index"`
	Firmware     string    `gorm:"type:varchar(50)"`
	Config       string    `gorm:"type:jsonb"`
	LastSeenAt   *time.Time
	InstalledAt  *time.Time
	Metadata     string `gorm:"type:jsonb"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (DeviceModel) TableName() string { return "device_devices" }

type deviceRepoImpl struct{ db *gorm.DB }

func NewDeviceRepository(db *gorm.DB) *deviceRepoImpl {
	return &deviceRepoImpl{db: db}
}

func (r *deviceRepoImpl) Save(ctx context.Context, d *entity.Device) error {
	m := toDeviceModel(d)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *deviceRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	var m DeviceModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDeviceEntity(&m), nil
}

func (r *deviceRepoImpl) FindBySerial(ctx context.Context, serial string) (*entity.Device, error) {
	var m DeviceModel
	err := r.db.WithContext(ctx).First(&m, "serial_number = ?", serial).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrDeviceNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDeviceEntity(&m), nil
}

func (r *deviceRepoImpl) FindBySite(ctx context.Context, siteID uuid.UUID) ([]entity.Device, error) {
	var models []DeviceModel
	if err := r.db.WithContext(ctx).Where("site_id = ?", siteID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]entity.Device, len(models))
	for i, m := range models {
		out[i] = *toDeviceEntity(&m)
	}
	return out, nil
}

func (r *deviceRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&DeviceModel{}, "id = ?", id).Error
}

func toDeviceModel(d *entity.Device) *DeviceModel {
	return &DeviceModel{
		ID:           d.ID,
		TenantID:     d.TenantID,
		SiteID:       d.SiteID,
		ZoneID:       d.ZoneID,
		SerialNumber: d.SerialNumber,
		Name:         d.Name,
		Type:         string(d.Type),
		Protocol:     string(d.Protocol),
		Status:       string(d.Status),
		Firmware:     d.Firmware,
		LastSeenAt:   d.LastSeenAt,
		InstalledAt:  d.InstalledAt,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func toDeviceEntity(m *DeviceModel) *entity.Device {
	return &entity.Device{
		ID:           m.ID,
		TenantID:     m.TenantID,
		SiteID:       m.SiteID,
		ZoneID:       m.ZoneID,
		SerialNumber: m.SerialNumber,
		Name:         m.Name,
		Type:         valueobject.DeviceType(m.Type),
		Protocol:     valueobject.Protocol(m.Protocol),
		Status:       valueobject.DeviceStatus(m.Status),
		Firmware:     m.Firmware,
		LastSeenAt:   m.LastSeenAt,
		InstalledAt:  m.InstalledAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
```

### 7.2 Kafka Consumers

```go
// internal/modules/device/infrastructure/messaging/consumers/telemetry_consumer.go
package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/device/application"
)

type TelemetryConsumer struct {
	ingestUC *application.IngestTelemetryUseCase
}

func NewTelemetryConsumer(uc *application.IngestTelemetryUseCase) *TelemetryConsumer {
	return &TelemetryConsumer{ingestUC: uc}
}

func (c *TelemetryConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *TelemetryConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *TelemetryConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload application.IngestTelemetryInput
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("[telemetry] unmarshal error: %v", err)
			s.MarkMessage(msg, "")
			continue
		}
		if err := c.ingestUC.Execute(context.Background(), payload); err != nil {
			log.Printf("[telemetry] ingest error: %v", err)
		}
		s.MarkMessage(msg, "")
	}
	return nil
}
```

### 7.3 WebSocket Broadcaster

```go
// internal/modules/device/interfaces/websocket/hub.go
package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool // tenantID → clients
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]bool)}
}

func (h *Hub) Register(tenantID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[tenantID] == nil {
		h.clients[tenantID] = make(map[*Client]bool)
	}
	h.clients[tenantID][c] = true
}

func (h *Hub) Unregister(tenantID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m, ok := h.clients[tenantID]; ok {
		delete(m, c)
	}
}

func (h *Hub) BroadcastToTenant(tenantID string, event interface{}) {
	data, _ := json.Marshal(event)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[tenantID] {
		select {
		case c.Send <- data:
		default:
		}
	}
}

type Client struct {
	TenantID string
	UserID   string
	Conn     *websocket.Conn
	Send     chan []byte
}
```

### 7.4 MQTT Gateway

```go
// internal/modules/device/infrastructure/iot/mqtt_subscriber.go
package iot

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"icmongolang/internal/modules/device/application"
)

type MQTTSubscriber struct {
	client   mqtt.Client
	ingestUC *application.IngestTelemetryUseCase
}

func NewMQTTSubscriber(broker, clientID string, uc *application.IngestTelemetryUseCase) (*MQTTSubscriber, error) {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	s := &MQTTSubscriber{client: c, ingestUC: uc}
	if token := c.Subscribe("iot/+/telemetry", 1, s.handle); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return s, nil
}

func (s *MQTTSubscriber) handle(_ mqtt.Client, msg mqtt.Message) {
	var payload application.IngestTelemetryInput
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[mqtt] unmarshal: %v", err)
		return
	}
	if err := s.ingestUC.Execute(context.Background(), payload); err != nil {
		log.Printf("[mqtt] ingest: %v", err)
	}
}
```

### 7.5 JWT, Bcrypt, Rate Limit

ใช้ `pkg/jwt`, `pkg/cryptpass` ที่มีอยู่แล้ว และ middleware ใน `internal/middleware/rate_limit.go`

---

## 8. Interface Layer

### 8.1 HTTP Handlers

```go
// internal/modules/device/interfaces/http/device_handler.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icmongolang/internal/modules/device/application"
	valueobject "icmongolang/internal/modules/device/domain/value_object"
)

type DeviceHandler struct {
	registerUC *application.RegisterDeviceUseCase
	getUC      *application.GetDeviceUseCase
	telemetryUC *application.IngestTelemetryUseCase
}

func NewDeviceHandler(
	registerUC *application.RegisterDeviceUseCase,
	getUC *application.GetDeviceUseCase,
	telemetryUC *application.IngestTelemetryUseCase,
) *DeviceHandler {
	return &DeviceHandler{registerUC: registerUC, getUC: getUC, telemetryUC: telemetryUC}
}

func (h *DeviceHandler) Register(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var req struct {
		SiteID       string `json:"site_id" binding:"required"`
		ZoneID       string `json:"zone_id"`
		SerialNumber string `json:"serial_number" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Type         string `json:"type" binding:"required"`
		Protocol     string `json:"protocol" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	siteID, _ := uuid.Parse(req.SiteID)
	zoneID, _ := uuid.Parse(req.ZoneID)

	out, err := h.registerUC.Execute(c.Request.Context(), application.RegisterDeviceInput{
		TenantID:     tenantID,
		SiteID:       siteID,
		ZoneID:       zoneID,
		SerialNumber: req.SerialNumber,
		Name:         req.Name,
		Type:         valueobject.DeviceType(req.Type),
		Protocol:     valueobject.Protocol(req.Protocol),
		UserID:       uid,
		IPAddress:    c.ClientIP(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": out})
}

func (h *DeviceHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	d, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}
```

### 8.2 Routes

```go
// internal/modules/device/interfaces/http/routes.go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
	Device *DeviceHandler
}

func RegisterRoutes(r *gin.RouterGroup, h *Handlers, auth gin.HandlerFunc, tenant gin.HandlerFunc) {
	g := r.Group("/devices")
	g.Use(auth, tenant)

	g.POST("", h.Device.Register)
	g.GET("/:id", h.Device.Get)
	g.GET("", h.Device.List)
	g.PUT("/:id", h.Device.Update)
	g.DELETE("/:id", h.Device.Delete)
	g.POST("/:id/command", h.Device.SendCommand)
	g.GET("/:id/telemetry", h.Device.QueryTelemetry)
}
```

### 8.3 Middleware

```go
// internal/modules/device/interfaces/middleware/tenant.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.GetHeader("X-Tenant-ID")
		if tid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tenant required"})
			return
		}
		id, err := uuid.Parse(tid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid tenant"})
			return
		}
		c.Set("tenant_id", id)
		c.Next()
	}
}
```

---

## 9. Database Migrations

### 9.1 Customer

```sql
-- migrations/20260101_customer_init.sql
CREATE TABLE IF NOT EXISTS customer_customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOT NULL DEFAULT 'PERSONAL',
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    tax_id VARCHAR(20),
    address JSONB,
    package_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (tenant_id, code)
);
CREATE INDEX idx_customer_customers_tenant ON customer_customers (tenant_id);
CREATE INDEX idx_customer_customers_status ON customer_customers (status);

CREATE TABLE IF NOT EXISTS customer_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(255),
    position VARCHAR(100),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_customer_contacts_customer ON customer_contacts (customer_id);
```

### 9.2 Device

```sql
-- migrations/20260101_device_init.sql
CREATE TABLE IF NOT EXISTS device_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    site_id UUID NOT NULL,
    zone_id UUID,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    protocol VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'OFFLINE',
    firmware VARCHAR(50),
    config JSONB,
    last_seen_at TIMESTAMP,
    installed_at TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_device_devices_tenant ON device_devices (tenant_id);
CREATE INDEX idx_device_devices_site ON device_devices (site_id);
CREATE INDEX idx_device_devices_status ON device_devices (status);
CREATE INDEX idx_device_devices_type ON device_devices (type);

CREATE TABLE IF NOT EXISTS device_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES device_devices(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    payload JSONB,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    issued_by UUID NOT NULL,
    issued_at TIMESTAMP DEFAULT NOW(),
    executed_at TIMESTAMP
);
CREATE INDEX idx_device_commands_device ON device_commands (device_id, status);

-- Telemetry เก็บใน InfluxDB (ไม่สร้าง table ใน Postgres)
```

### 9.3 Package

```sql
-- migrations/20260101_package_init.sql
CREATE TABLE IF NOT EXISTS package_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tier VARCHAR(30) NOT NULL,
    price_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    price_currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    billing_cycle VARCHAR(20) NOT NULL,
    quotas JSONB,
    features JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS package_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    package_id UUID NOT NULL REFERENCES package_packages(id),
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    auto_renew BOOLEAN DEFAULT TRUE,
    usage JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_package_subscriptions_customer ON package_subscriptions (customer_id, status);
```

### 9.4 ERP

```sql
-- migrations/20260101_erp_init.sql
CREATE TABLE IF NOT EXISTS erp_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount NUMERIC(12,2) NOT NULL,
    total_currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    ordered_at TIMESTAMP DEFAULT NOW(),
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_erp_orders_customer ON erp_orders (customer_id);
CREATE INDEX idx_erp_orders_status ON erp_orders (status);

CREATE TABLE IF NOT EXISTS erp_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES erp_orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(12,2) NOT NULL,
    subtotal NUMERIC(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS erp_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES erp_orders(id),
    customer_id UUID NOT NULL,
    number VARCHAR(50) UNIQUE NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    tax_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    total NUMERIC(12,2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    issued_at TIMESTAMP,
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS erp_stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    sku VARCHAR(100) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved INT NOT NULL DEFAULT 0,
    reorder_point INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (product_id, warehouse_id)
);
```

### 9.5 CRM

```sql
-- migrations/20260101_crm_init.sql
CREATE TABLE IF NOT EXISTS crm_leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(30),
    company VARCHAR(255),
    source VARCHAR(50),
    status VARCHAR(30) NOT NULL DEFAULT 'NEW',
    score INT DEFAULT 0,
    assigned_to UUID,
    converted_to UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_crm_leads_tenant ON crm_leads (tenant_id);
CREATE INDEX idx_crm_leads_status ON crm_leads (status);

CREATE TABLE IF NOT EXISTS crm_opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id UUID REFERENCES crm_leads(id),
    customer_id UUID,
    title VARCHAR(255) NOT NULL,
    amount NUMERIC(12,2),
    probability INT,
    stage VARCHAR(30) NOT NULL DEFAULT 'PROSPECTING',
    expected_close_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS crm_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    subject VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN',
    assigned_to UUID,
    sla_due_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_crm_tickets_customer ON crm_tickets (customer_id, status);
```

### 9.6 Logistics

```sql
-- migrations/20260101_logistics_init.sql
CREATE TABLE IF NOT EXISTS logistics_shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    tracking_number VARCHAR(50) UNIQUE NOT NULL,
    order_id UUID,
    customer_id UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    origin JSONB,
    destination JSONB,
    carrier VARCHAR(100),
    vehicle_id UUID,
    driver_id UUID,
    estimated_at TIMESTAMP,
    delivered_at TIMESTAMP,
    temperature JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_logistics_shipments_tenant ON logistics_shipments (tenant_id);
CREATE INDEX idx_logistics_shipments_status ON logistics_shipments (status);
CREATE INDEX idx_logistics_shipments_tracking ON logistics_shipments (tracking_number);

CREATE TABLE IF NOT EXISTS logistics_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    waypoints JSONB,
    distance_km NUMERIC(10,2),
    estimated_min INT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 9.7 Report

```sql
-- migrations/20260101_report_init.sql
CREATE TABLE IF NOT EXISTS report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50),
    type VARCHAR(30),
    query TEXT NOT NULL,
    parameters JSONB,
    schedule JSONB,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (tenant_id, code)
);

CREATE TABLE IF NOT EXISTS report_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES report_definitions(id),
    run_by UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    parameters JSONB,
    result_url TEXT,
    error_msg TEXT,
    started_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP
);
CREATE INDEX idx_report_runs_report ON report_runs (report_id, status);
```

---

## 10. System Flow

### 10.1 Device Telemetry Flow

```
[Device/Sensor]
     │ (MQTT publish)
     ▼
[MQTT Broker] ──► [MQTT Subscriber] ──► [IngestTelemetryUseCase]
                                              │
                    ┌─────────────────────────┼──────────────────────┐
                    ▼                         ▼                      ▼
              [InfluxDB]              [AutomationService]     [Kafka Producer]
              (time-series)            (evaluate rules)              │
                    │                         │                      │
                    │                         ▼                      ▼
                    │                  [CommandIssuer]         [Kafka Topics]
                    │                         │                      │
                    │                         ▼                      ├─► [AI Consumer]
                    │                  [MQTT Publish]               ├─► [Alarm Consumer]
                    │                  (to Actuator)                ├─► [ES Indexer]
                    │                                              └─► [WS Broadcaster]
                    ▼
              [Grafana Dashboard]
```

### 10.2 Customer Onboarding Flow

```
1. [Visitor] ──► ลงทะเบียนในเว็บ
2. [System] ──► สร้าง User + Customer (status = PENDING)
3. [System] ──► ส่งอีเมลยืนยัน
4. [User] ──► ยืนยันอีเมล
5. [User] ──► เลือก Package
6. [System] ──► สร้าง Subscription + Payment
7. [User] ──► ชำระเงิน
8. [System] ──► เปิดใช้งาน (status = ACTIVE)
9. [System] ──► Provision tenant + ส่ง welcome email
10. [User] ──► เพิ่ม Site, Zone, Device
```

### 10.3 AI/Automation Flow

```
[Telemetry] ──► [Kafka: device.telemetry.ingested]
                        │
                        ▼
              [AI Consumer Group]
                        │
         ┌──────────────┼──────────────┐
         ▼              ▼              ▼
   [LLM Analyze]  [Anomaly]    [Predict]
   (Ollama)       (ML model)   (forecast)
         │              │              │
         └──────────────┼──────────────┘
                        ▼
              [Publish: device.ai.insight]
                        │
         ┌──────────────┼──────────────┐
         ▼              ▼              ▼
   [Alarm]        [Automation]    [Report]
   (notify)       (auto-action)   (dashboard)
```

### 10.4 Logistics Cold Chain Flow

```
[Shipment Created] ──► [Temp Sensor attached]
                              │
                              ▼
                    [MQTT: logistics/{id}/temp]
                              │
                              ▼
                    [Logistics Consumer]
                              │
                    ┌─────────┼─────────┐
                    ▼         ▼         ▼
              [Record]   [Check]   [Broadcast]
              InfluxDB   Breach?   WebSocket
                         │
                    ┌────┴────┐
                    ▼         ▼
                 [Alert]   [Continue]
                 (SMS/Email)
```

---

## 11. Workflow Diagram

### 11.1 Module Interaction (Sequence)

```
Client ──► API Gateway ──► Device Handler ──► UseCase ──► Repository ──► DB
                                                  │
                                                  ├──► Kafka ──► Consumers
                                                  ├──► Redis (cache)
                                                  ├──► InfluxDB (telemetry)
                                                  └──► WebSocket Hub
```

### 11.2 Bounded Context Integration

```
┌──────────┐    events    ┌──────────┐
│  Device  │─────────────►│  Alarm   │
└────┬─────┘              └────┬─────┘
     │                         │
     │ events                  │ events
     ▼                         ▼
┌──────────┐              ┌──────────┐
│  Report  │◄─────────────│ Notifier │
└──────────┘              └──────────┘
     ▲
     │ events
     │
┌────┴─────┐    events    ┌──────────┐
│   ERP    │◄─────────────│Logistics │
└──────────┘              └──────────┘
```

### 11.3 State Machines

**Device Status:**
```
OFFLINE ──(heartbeat)──► ONLINE ──(error)──► ERROR
   ▲                        │                 │
   │                        │ (maintenance)   │ (fixed)
   │                        ▼                 │
   └────(timeout)──── MAINTENANCE ◄───────────┘
```

**Subscription Status:**
```
ACTIVE ──(payment fail)──► PAST_DUE ──(pay)──► ACTIVE
   │                          │
   │ (cancel)                 │ (30 days)
   ▼                          ▼
CANCELLED                  EXPIRED
```

**Shipment Status:**
```
PENDING ──► PICKED ──► IN_TRANSIT ──► DELIVERED
                            │
                            └──► FAILED ──► (retry) ──► IN_TRANSIT
```

---

## 12. การติดตั้งและใช้งาน

### 12.1 Requirements

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Kafka 3.x
- InfluxDB 2.x
- Elasticsearch 8.x
- MQTT Broker (Mosquitto/EMQX)
- Docker + Docker Compose (recommended)

### 12.2 Quick Start

```bash
# 1. Clone
git clone <repo> icmongolang
cd icmongolang

# 2. Setup environment
cp .env.example .env
# แก้ไข DB_DSN, REDIS_ADDR, KAFKA_BROKERS, MQTT_BROKER

# 3. Start infrastructure
docker-compose -f docker-compose.dev.yml up -d

# 4. Run migrations
go run cmd/migrate.go up

# 5. Seed data
go run cmd/initdata.go

# 6. Run API
go run cmd/api/main.go

# 7. Run MQTT Gateway (อีก terminal)
go run cmd/mqtt-gateway/main.go

# 8. Run Workers
go run cmd/workers/device-telemetry/main.go
```

### 12.3 Environment Variables

```env
# Core
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug

# Database
DB_DSN=postgres://user:pass@localhost:5432/icmongolang?sslmode=disable
DB_MAX_OPEN=25
DB_MAX_IDLE=10

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=icmongolang-group

# MQTT
MQTT_BROKER=tcp://localhost:1883
MQTT_CLIENT_ID=icmongolang-gateway
MQTT_USERNAME=
MQTT_PASSWORD=

# InfluxDB
INFLUX_URL=http://localhost:8086
INFLUX_TOKEN=
INFLUX_ORG=icmongolang
INFLUX_BUCKET=telemetry

# Elasticsearch
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_INDEX=icmongolang_*

# JWT
JWT_SECRET=change-me-in-production
JWT_TTL=24h

# LLM
LLM_PROVIDER=ollama
LLM_API_URL=http://localhost:11434
LLM_MODEL=llama3

# Payment
PAYMENT_PROVIDER=stripe
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=

# Module-specific
DEVICE_TELEMETRY_RETENTION_DAYS=90
PACKAGE_DEFAULT_TIER=FREE
```

### 12.4 Docker Compose (ตัวอย่าง)

```yaml
version: '3.9'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: icmongolang
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    ports: ["5432:5432"]
    volumes: [pgdata:/var/lib/postgresql/data]

  redis:
    image: redis:7
    ports: ["6379:6379"]

  kafka:
    image: confluentinc/cp-kafka:latest
    ports: ["9092:9092"]
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  influxdb:
    image: influxdb:2.7
    ports: ["8086:8086"]
    environment:
      DOCKER_INFLUXDB_INIT_MODE: setup
      DOCKER_INFLUXDB_INIT_USERNAME: admin
      DOCKER_INFLUXDB_INIT_PASSWORD: adminpass
      DOCKER_INFLUXDB_INIT_ORG: icmongolang
      DOCKER_INFLUXDB_INIT_BUCKET: telemetry

  mosquitto:
    image: eclipse-mosquitto:2
    ports: ["1883:1883", "9001:9001"]
    volumes: [./mqtt/mosquitto.conf:/mosquitto/config/mosquitto.conf]

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      discovery.type: single-node
      xpack.security.enabled: "false"
    ports: ["9200:9200"]

volumes:
  pgdata:
```

---

## 13. Business Model

### 13.1 Revenue Streams

| Stream | Model | ตัวอย่าง |
| :--- | :--- | :--- |
| **SaaS Subscription** | รายเดือน/ปี | Free, Basic ฿990/เดือน, Pro ฿4,900/เดือน, Enterprise custom |
| **Device Management Fee** | ต่อ device/เดือน | ฿20/device/เดือน |
| **Hardware Sales** | ขายขาด | Sensor, Gateway, Actuator |
| **Installation Service** | ครั้งเดียว | ค่าติดตั้งตามไซต์ |
| **AI/Analytics Add-on** | รายเดือน | ฿2,900/เดือน |
| **API/Integration Fee** | ตามการใช้งาน | ฿0.10/API call เกินโควตา |
| **White-label** | License | สัญญารายปี |
| **Support & SLA** | รายเดือน | 24/7 support |

### 13.2 Package Tiers

| Feature | Free | Basic | Pro | Enterprise |
| :--- | :-: | :-: | :-: | :-: |
| Devices | 5 | 50 | 500 | Unlimited |
| Sites | 1 | 5 | 50 | Unlimited |
| Users | 2 | 10 | 100 | Unlimited |
| Storage | 1 GB | 10 GB | 100 GB | Custom |
| Data retention | 7 วัน | 30 วัน | 90 วัน | 365 วัน |
| Automation rules | 5 | 50 | 500 | Unlimited |
| AI Analytics | ❌ | ✅ (basic) | ✅ | ✅ |
| API access | ❌ | ✅ | ✅ | ✅ |
| WebSocket | ❌ | ✅ | ✅ | ✅ |
| SLA | – | 99% | 99.9% | 99.95% |
| Support | Community | Email | Email+Chat | 24/7 Phone |
| ราคา/เดือน | ฟรี | ฿990 | ฿4,900 | ติดต่อ |

### 13.3 Target Market (TAM/SAM/SOM)

| ระดับ | ขนาด | หมายเหตุ |
| :--- | :--- | :--- |
| **TAM** | ตลาด IoT ไทย ~฿50,000 ล้าน | ทั้งหมด |
| **SAM** | Smart Farm + Smart Building ~฿8,000 ล้าน | กลุ่มเป้าหมาย |
| **SOM** | 3 ปีแรก ~฿240 ล้าน | ส่วนแบ่งที่ตั้งเป้า |

### 13.4 Unit Economics (ตัวอย่าง)

```
ARPU (Basic): ฿990/เดือน
CAC: ฿3,500
LTV: ฿35,640 (36 เดือน × ฿990)
LTV/CAC: 10.2x
Gross Margin: 75%
Payback Period: 3.5 เดือน
```

### 13.5 Go-to-Market

1. **Phase 1 (เดือน 1-6):** Pilot กับฟาร์ม/อาคาร 10 ราย ฟรี
2. **Phase 2 (เดือน 7-12):** เปิด Free + Basic, หา partner ติดตั้ง
3. **Phase 3 (ปี 2):** ขยาย Pro/Enterprise, เพิ่ม AI
4. **Phase 4 (ปี 3):** White-label, ขยาย ASEAN

---

## 14. Roadmap

### 14.1 Phase Overview

```
Phase 0 ─ Foundation     (เดือน 1-2)
Phase 1 ─ MVP            (เดือน 3-6)
Phase 2 ─ Growth         (เดือน 7-12)
Phase 3 ─ Enterprise     (เดือน 13-18)
Phase 4 ─ Scale          (เดือน 19-24)
```

### 14.2 Phase 0 — Foundation (เดือน 1-2)

| สัปดาห์ | งาน | Deliverable |
| :--- | :--- | :--- |
| 1-2 | Setup repo, CI/CD, docker-compose | Repo + CI ผ่าน |
| 3-4 | Auth, Users, Customer module | Login ได้, CRUD customer |
| 5-6 | Package module | สร้างแพ็กเกจ + subscription |
| 7-8 | Device module (basic) | Register device, MQTT ingest |

**Milestone:** ระบบพื้นฐานทำงานได้ end-to-end

### 14.3 Phase 1 — MVP (เดือน 3-6)

| เดือน | งาน | Deliverable |
| :--- | :--- | :--- |
| 3 | Telemetry pipeline (MQTT→Kafka→InfluxDB) | เก็บข้อมูลได้ |
| 4 | WebSocket realtime + Dashboard | ดู realtime ได้ |
| 5 | Automation rules + Alarm | กฎอัตโนมัติทำงาน |
| 6 | Payment + Billing | ชำระเงินได้ |

**Milestone:** MVP พร้อม pilot 10 ราย

### 14.4 Phase 2 — Growth (เดือน 7-12)

| เดือน | งาน | Deliverable |
| :--- | :--- | :--- |
| 7-8 | ERP module (Order, Invoice, Stock) | ออกใบเสร็จได้ |
| 9-10 | CRM module (Lead, Opportunity, Ticket) | จัดการลูกค้าได้ |
| 11 | Report module + Analytics | Dashboard ขั้นสูง |
| 12 | AI/Automation (LLM integration) | วิเคราะห์อัตโนมัติ |

**Milestone:** เปิด sprice ขายจริง

### 14.5 Phase 3 — Enterprise (เดือน 13-18)

| เดือน | งาน | Deliverable |
| :--- | :--- | :--- |
| 13-14 | Logistics + Cold Chain | ติดตามการขนส่ง |
| 15-16 | Multi-tenancy + RLS | แยกข้อมูล tenant |
| 17 | Mobile App (React Native) | แอปมือถือ |
| 18 | White-label + API marketplace | ขาย OEM |

**Milestone:** Enterprise-ready

### 14.6 Phase 4 — Scale (เดือน 19-24)

| เดือน | งาน | Deliverable |
| :--- | :--- | :--- |
| 19-20 | Performance tuning (10k devices) | รองรับ 10k devices |
| 21-22 | Global expansion (multi-region) | Deploy 2 region |
| 23-24 | Advanced AI (predictive, computer vision) | AI ขั้นสูง |

**Milestone:** Scale ได้ 100k devices

### 14.7 Feature Priority Matrix

| Feature | Impact | Effort | Priority |
| :--- | :-: | :-: | :-: |
| Device management | High | Med | P0 |
| Telemetry pipeline | High | Med | P0 |
| Realtime dashboard | High | Med | P0 |
| Automation | High | High | P1 |
| Payment | High | Med | P1 |
| ERP | Med | High | P2 |
| CRM | Med | Med | P2 |
| AI Analytics | High | High | P2 |
| Logistics | Med | High | P3 |
| Mobile App | High | High | P3 |
| White-label | Med | Med | P4 |

---

## 15. ภาคผนวก

### A. Dependency Matrix

| Layer | stdlib | uuid | gorm | gin | sarama | redis | influx | es | mqtt |
| :--- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| Domain | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Application | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Infrastructure | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Interface | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ |

### B. Kafka Topic Convention

```
<context>.<aggregate>.<event>

device.device.registered
device.device.online
device.device.offline
device.telemetry.ingested
device.command.issued
device.command.acked
device.ai.insight

customer.customer.created
customer.customer.activated
customer.customer.suspended

package.subscription.created
package.subscription.upgraded
package.subscription.cancelled

erp.order.created
erp.order.shipped
erp.invoice.issued
erp.invoice.paid

crm.lead.created
crm.lead.converted
crm.ticket.opened
crm.ticket.resolved

logistics.shipment.created
logistics.shipment.delivered
logistics.temperature.breached

report.run.started
report.run.completed
```

### C. Error Code Convention (HTTP)

| Domain Error | HTTP |
| :--- | :-: |
| `ErrNotFound` | 404 |
| `ErrAlreadyExists` | 409 |
| `ErrInvalid*` | 400 |
| `ErrUnauthorized` | 401 |
| `ErrForbidden` | 403 |
| `ErrRateLimit` | 429 |
| `ErrDeviceOffline` | 503 |

### D. Table Prefix Summary

| Module | Prefix |
| :--- | :--- |
| customer | `customer_` |
| device | `device_` |
| package | `package_` |
| erp | `erp_` |
| crm | `crm_` |
| logistics | `logistics_` |
| report | `report_` |
| alarm | `alarm_` |
| auth | `auth_` |
| payment | `payment_` |

---

## 16. Prompt สำหรับการขยายระบบในอนาคต

### 16.1 Prompt สร้าง Module ใหม่

```
สร้างโมดูลใหม่ในโปรเจกต์ icmongolang ตาม Template_Module.md

## ข้อมูลนำเข้า
- module_name: <ชื่อ เช่น irrigation, energy, security>
- aggregate: <ชื่อ aggregate root เช่น IrrigationZone>
- actions: Create, Get, Update, Delete, List
- cross-cutting: <Kafka/ES/Redis/InfluxDB/WebSocket/MQTT/LLM>
- database prefix: <module_name>_
- bounded context: <context ที่สังกัด>

## กฎการสร้าง
1. โครงสร้างตาม Clean Architecture + DDD
2. Domain Layer:
   - Entity มี constructor + behavior methods
   - Value Object มี IsValid()
   - Repository เป็น interface
   - Domain error เป็น sentinel error
3. Application Layer:
   - 1 use case ต่อ 1 ไฟล์, Execute() return error
4. Infrastructure Layer:
   - GORM model มี TableName() + prefix
   - Repository impl map model ↔ entity
5. Interface Layer:
   - Gin handler ดึง user_id/tenant_id จาก context
   - Route ลงทะเบียนใน routes.go
6. Import จาก pkg/ เท่านั้น
7. สร้าง migration YYYYMMDD_<module>_init.sql
8. Wire-up ใน module.go + cmd/api/main.go
9. รัน go build ./... && go test ./... ให้ผ่าน

## Output
1. Tree ของไฟล์
2. โค้ดทุกไฟล์
3. Migration SQL
4. ตัวอย่าง .env
5. วิธีรัน
```

### 16.2 Prompt เพิ่ม Use Case

```
เพิ่ม use case {{Verb}}{{Aggregate}} ใน module {{module_name}}
- Input: {{fields}}
- Business rules: {{อธิบาย}}
- Side effects: {{audit/kafka/cache/influx/ws}}
- สร้างที่ application/{{verb}}_{{aggregate}}.go
- อัปเดต module.go และ handler
- เขียน test
```

### 16.3 Prompt เพิ่ม AI Feature

```
เพิ่ม AI feature ใน module {{module_name}}
- Use case: Analyze{{Aggregate}}
- Input: {{data}}
- LLM: ใช้ pkg/llm
- Prompt: {{prompt template}}
- Output: structured JSON
- Side effects: publish insight ไป Kafka, broadcast WS
- สร้างที่ application/analyze_{{aggregate}}.go
- สร้าง consumer ที่ infrastructure/messaging/consumers/
```

### 16.4 Prompt เพิ่ม Device Protocol

```
เพิ่ม protocol {{protocol}} ใน module device
- สร้าง parser ที่ infrastructure/iot/protocol/{{protocol}}.go
- Implement interface ProtocolParser
- Map payload → entity.Telemetry
- ลงทะเบียนใน protocol registry
- เขียน test
```

### 16.5 Prompt เพิ่ม Report

```
เพิ่ม report {{name}} ใน module report
- Category: {{category}}
- Query: {{SQL}}
- Parameters: {{list}}
- Schedule: {{cron}}
- Output: PDF/Excel/CSV
- สร้าง ReportDefinition seed
- เขียน test
```

---

## 17. การตั้งชื่อตาราง Database (Prefix)

| Module | Prefix | ตัวอย่าง |
| :--- | :--- | :--- |
| customer | `customer_` | `customer_customers`, `customer_contacts` |
| device | `device_` | `device_devices`, `device_commands` |
| package | `package_` | `package_packages`, `package_subscriptions` |
| erp | `erp_` | `erp_orders`, `erp_invoices`, `erp_stock_items` |
| crm | `crm_` | `crm_leads`, `crm_opportunities`, `crm_tickets` |
| logistics | `logistics_` | `logistics_shipments`, `logistics_routes` |
| report | `report_` | `report_definitions`, `report_runs` |
| alarm | `alarm_` | `alarm_rules`, `alarm_events` |
| auth | `auth_` | `auth_sessions`, `auth_refresh_tokens` |
| payment | `payment_` | `payment_payments`, `payment_transactions` |
| notifier | `notifier_` | `notifier_templates`, `notifier_logs` |
| auditlog | `audit_` | `audit_trails` |
| pdpa | `pdpa_` | `pdpa_consents`, `pdpa_dsar` |

**กฎ:**
- ✅ Table: `<module>_<plural>` (snake_case)
- ✅ Column: `snake_case`
- ✅ PK: `id UUID DEFAULT gen_random_uuid()`
- ✅ FK: `<entity>_id`
- ✅ Index: `idx_<table>_<column>`
- ✅ Unique: `uq_<table>_<column>`
- ✅ Migration: `YYYYMMDD_<module>_<desc>.sql`

---

## 18. DDD Validation Checklist

### 18.1 Domain Layer
- [ ] Entity ทุกตัวมี constructor (`New{{Entity}}`)
- [ ] Entity ไม่มี setter ตรง – เปลี่ยน state ผ่าน behavior method
- [ ] Value Object ทุกตัวมี `IsValid()`
- [ ] Repository เป็น interface เท่านั้น
- [ ] Domain error เป็น sentinel error
- [ ] Domain **ไม่ import** gorm, gin, sarama, redis, mqtt
- [ ] Invariants ถูกบังคับใน constructor และ behavior methods
- [ ] Aggregate Root เดียวต่อ transaction

### 18.2 Application Layer
- [ ] Use case ละ 1 ไฟล์, ชื่อ `{{verb}}_{{aggregate}}.go`
- [ ] Input/Output DTO แยกชัดเจน
- [ ] Execute() return error, ไม่ panic
- [ ] ไม่มี SQL/HTTP ใน use case
- [ ] Side effects ที่ไม่ critical → log warn แต่ไม่ fail
- [ ] Idempotency key สำหรับ operation ที่สำคัญ

### 18.3 Infrastructure Layer
- [ ] GORM model มี `TableName()` + prefix
- [ ] Repository impl map model ↔ entity ถูกต้อง
- [ ] Kafka message ใช้ JSON + schema version
- [ ] Redis key มี namespace (`{{module}}:{{entity}}:<id>`)
- [ ] ES index มีชื่อสอดคล้อง (`{{module}}_{{entity}}`)
- [ ] MQTT topic ตาม convention
- [ ] InfluxDB measurement + tag ตาม convention

### 18.4 Interface Layer
- [ ] Handler ดึง `user_id` + `tenant_id` จาก context
- [ ] Route ลงทะเบียนครบ
- [ ] Error response เป็น JSON สอดคล้อง
- [ ] Rate limit / auth / tenant middleware applied
- [ ] Input validation ที่ handler (defense in depth)

### 18.5 Build & Test
- [ ] `go build ./...` ผ่าน
- [ ] `go vet ./...` ผ่าน
- [ ] `go test ./...` ผ่าน
- [ ] Coverage: Domain ≥ 90%, App ≥ 80%
- [ ] ไม่มี import นอก whitelist

### 18.6 Migration
- [ ] ไฟล์ชื่อ `YYYYMMDD_<module>_<desc>.sql`
- [ ] ทุกตารางมี prefix
- [ ] มี index สำหรับ query ที่ใช้บ่อย
- [ ] มี FK constraint
- [ ] มี down script (ถ้าใช้ golang-migrate)

### 18.7 Observability
- [ ] ทุก use case มี log + metric
- [ ] Health check endpoint (`/healthz`, `/readyz`)
- [ ] Prometheus metrics
- [ ] Distributed tracing (ถ้ามี)

### 18.8 Security
- [ ] Input validation ที่ handler + domain
- [ ] Parameterized query (GORM)
- [ ] JWT + refresh token rotation
- [ ] Rate limiting
- [ ] Security headers
- [ ] Secret ไม่อยู่ใน code

### 18.9 Documentation
- [ ] อัปเดต README ของ module
- [ ] เพิ่มตัวอย่าง .env ถ้ามี env ใหม่
- [ ] อัปเดต docker-compose ถ้ามี service ใหม่
- [ ] Swagger annotations

---

## สรุป

เอกสารนี้เป็น **SRS + Roadmap + Module Design** สำหรับแพลตฟอร์ม IoT (Smart Farm / Smart Building) บน `icmongolang` โดย:

1. **ใช้ Clean Architecture + DDD** ตาม `Template_Module.md`
2. **7 โมดูลใหม่:** customer, package, erp, crm, device, logistics, report
3. **ทำงานร่วมกับ 33 โมดูลเดิม** ผ่าน Bounded Contexts
4. **Roadmap 24 เดือน** แบ่ง 5 phase
5. **Business Model** ชัดเจน: SaaS + Device fee + Hardware + Service
6. **DDD Validation Checklist** ครบ 9 หมวด
7. **Prompt Templates** สำหรับขยายระบบในอนาคต

> **หมายเหตุ:** โค้ดในเอกสารนี้เป็นตัวอย่าง (illustrative) — ควรปรับให้เข้ากับ business logic จริง และรัน `go build ./... && go test ./...` ก่อน merge ทุกครั้ง