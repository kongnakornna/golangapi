# 🚀 Module-by-Module — AI-Ready Generation Kit

> **วัตถุประสงค์:** เอกสารนี้ทำให้คุณสามารถป้อนให้ AI เพื่อ generate code ได้ทั้ง **Go (Gin + GORM)** และ **Python (FastAPI + SQLAlchemy)** ตาม Clean Architecture + DDD
>
> **กลยุทธ์:** ทำ **Module 1 (customer) แบบสมบูรณ์ 100%** ทั้ง 2 ภาษา → เป็น **แม่แบบ (Gold Standard)** ให้ทุกโมดูลที่เหลือทำตาม

---

## 📐 กติกากลาง (ใช้ทั้ง Go และ Python)

### โครงสร้าง Layer (บังคับ)

```
domain/         ← ไม่ import อะไรภายนอกเลย (ยกเว้น stdlib + uuid)
application/    ← import domain เท่านั้น
infrastructure/ ← import domain + application
interfaces/     ← import ทุก layer
```

### Naming Convention (เหมือนกันทั้ง 2 ภาษา)

| สิ่ง | รูปแบบ | ตัวอย่าง |
| :--- | :--- | :--- |
| Aggregate | PascalCase | `Customer`, `Package`, `Device` |
| Table | `{module}_{plural}` | `customer_customers` |
| Use Case | `{Verb}{Aggregate}` | `CreateCustomer` |
| Endpoint | `/api/v1/{module}/{plural}` | `/api/v1/customers` |
| Kafka Topic | `{module}.{aggregate}.{event}` | `customer.customer.created` |
| Error | `Err{Reason}` / `{Reason}Error` | `ErrCustomerNotFound` |

### Output ที่ AI ต้องผลิต (ต่อ module)

1. **Domain layer** — entity, VO, event, error, repo interface
2. **Application layer** — use case + DTO + port
3. **Infrastructure layer** — repo impl, cache, producer, indexer
4. **Interface layer** — HTTP handler + routes + middleware
5. **Migration SQL**
6. **Module composition root** (wire-up)
7. **Test suite**
8. **.env.example**

---

# 📦 MODULE 01: `customer`

> **สถานะ:** 🟢 Greenfield · **Priority:** P0 · **Owner:** Core Team

---

## PART A — Business Specification (Framework-Agnostic)

> ส่วนนี้ใช้เหมือนกันทั้ง Go และ Python — AI จะอ่านส่วนนี้แล้ว generate code

### A.1 Business Context

**Purpose:** จัดการข้อมูลลูกค้า (บุคคล/นิติบุคคล/ราชการ) รองรับ multi-tenant, lifecycle, contact, KYC

**Stakeholders:**
- Sales — สร้าง lead → customer
- Support — ดูประวัติลูกค้า
- Billing — ผูก subscription
- Compliance — PDPA, audit trail
- Tenant Admin — จัดการลูกค้าของตัวเอง

**Out of Scope:**
- Subscription (อยู่ที่ `package`)
- Payment (อยู่ที่ `payment`)
- Sales pipeline (อยู่ที่ `crm`)
- Site/Device (อยู่ที่ `device`)

### A.2 Ubiquitous Language

| Term | ความหมาย |
| :--- | :--- |
| **Customer** | ลูกค้าที่ซื้อบริการ/สินค้า |
| **Tenant** | องค์กรที่แยกข้อมูล (multi-tenant boundary) |
| **Contact** | ผู้ติดต่อของลูกค้า |
| **Segment** | การแบ่งกลุ่ม (SME, Enterprise, Government, Startup, Individual) |
| **Lifecycle** | LEAD → PROSPECT → CUSTOMER → CHURNED |
| **Status** | PENDING → ACTIVE → SUSPENDED → CHURNED |
| **KYC** | Know Your Customer (เฉพาะ BUSINESS/GOVERNMENT) |
| **Primary Contact** | ผู้ติดต่อหลัก (มีได้ 1 คน) |

### A.3 Domain Model

#### Aggregate Root: `Customer`

| Field | Type | Constraint | หมายเหตุ |
| :--- | :--- | :--- | :--- |
| id | UUID | PK | |
| tenant_id | UUID | NOT NULL, indexed | multi-tenant |
| code | string(50) | NOT NULL, unique per tenant | `^[A-Z0-9-]{3,50}$` |
| name | string(255) | NOT NULL | |
| type | enum | PERSONAL / BUSINESS / GOVERNMENT / NGO | |
| segment | enum | SME / ENTERPRISE / GOVERNMENT / STARTUP / INDIVIDUAL | |
| lifecycle | enum | LEAD / PROSPECT / CUSTOMER / CHURNED | |
| status | enum | PENDING / ACTIVE / SUSPENDED / CHURNED | |
| tax_id | string(20) | required if type ∈ {BUSINESS, GOVERNMENT} | |
| address | Address (VO) | | |
| contacts | []Contact | max 20 | |
| package_id | *UUID | nullable | |
| kyc | *KYCInfo | nullable, only BUSINESS | |
| metadata | JSONB | | ข้อมูลเพิ่มเติม |
| created_at | timestamp | | |
| updated_at | timestamp | | |

#### Entity: `Contact` (in Customer aggregate)

| Field | Type | Constraint |
| :--- | :--- | :--- |
| id | UUID | PK |
| customer_id | UUID | FK |
| name | string(255) | NOT NULL |
| phone | string(30) | valid phone |
| email | string(255) | valid email |
| position | string(100) | |
| is_primary | bool | exactly 1 per customer |
| created_at | timestamp | |

#### Value Object: `KYCInfo`

| Field | Type |
| :--- | :--- |
| verified | bool |
| verified_at | *timestamp |
| verified_by | *UUID |
| document_type | string |
| document_no | string |
| expires_at | *timestamp |

### A.4 Value Objects (with validation)

| VO | รูปแบบ | Validation |
| :--- | :--- | :--- |
| `CustomerType` | enum | ∈ {PERSONAL, BUSINESS, GOVERNMENT, NGO} |
| `CustomerStatus` | enum | ∈ {PENDING, ACTIVE, SUSPENDED, CHURNED} |
| `Segment` | enum | ∈ {SME, ENTERPRISE, GOVERNMENT, STARTUP, INDIVIDUAL} |
| `LifecycleStage` | enum | ∈ {LEAD, PROSPECT, CUSTOMER, CHURNED} |
| `CustomerCode` | string | `^[A-Z0-9-]{3,50}$` |
| `Email` | string | RFC 5322, lowercase |
| `PhoneNumber` | string | E.164 หรือ Thai format |
| `TaxID` | string | 13 หลัก (Thai) |
| `Address` | struct | postcode = 5 หลัก |

### A.5 Invariants

| # | Invariant | บังคับที่ |
| :-: | :--- | :--- |
| 1 | `code` ไม่ซ้ำใน `tenant_id` | Repository + unique index |
| 2 | `tax_id` required สำหรับ BUSINESS/GOVERNMENT | Constructor |
| 3 | `contacts.length ≤ 20` | `AddContact` |
| 4 | Primary contact มีได้ 1 | `AddContact`, `SetPrimaryContact` |
| 5 | Status transitions เท่านั้น: PENDING→ACTIVE, ACTIVE→SUSPENDED, ACTIVE→CHURNED, SUSPENDED→ACTIVE, SUSPENDED→CHURNED | Behavior methods |
| 6 | `KYC` verify เฉพาะ BUSINESS | `VerifyKYC` use case |
| 7 | `postcode` = 5 หลัก | `Address.validate` |

### A.6 Domain Events

| Event | Trigger | Payload |
| :--- | :--- | :--- |
| `CustomerCreated` | CreateCustomer | `{customer_id, tenant_id, code, name, type}` |
| `CustomerUpdated` | UpdateCustomer | `{customer_id, changes}` |
| `CustomerActivated` | ActivateCustomer | `{customer_id}` |
| `CustomerSuspended` | SuspendCustomer | `{customer_id, reason}` |
| `CustomerChurned` | ChurnCustomer | `{customer_id, reason}` |
| `ContactAdded` | AddContact | `{customer_id, contact_id}` |
| `ContactRemoved` | RemoveContact | `{customer_id, contact_id}` |
| `PackageAssigned` | AssignPackage | `{customer_id, package_id}` |
| `KYCVerified` | VerifyKYC | `{customer_id, verified_at}` |

### A.7 Use Cases

| # | Use Case | Input | Output | Side Effects |
| :-: | :--- | :--- | :--- | :--- |
| 1 | `CreateCustomer` | tenant_id, code, name, type, segment, tax_id, address | customer_id | DB, Audit, Kafka |
| 2 | `GetCustomer` | customer_id | CustomerDetail | Cache read |
| 3 | `ListCustomers` | tenant_id, filter, page | []CustomerSummary | |
| 4 | `UpdateCustomer` | customer_id, patch | | DB, Audit, Kafka, Cache invalidate |
| 5 | `ActivateCustomer` | customer_id, user_id | | DB, Audit, Kafka, WS |
| 6 | `SuspendCustomer` | customer_id, reason, user_id | | DB, Audit, Kafka, WS |
| 7 | `ChurnCustomer` | customer_id, reason, user_id | | DB, Audit, Kafka |
| 8 | `AddContact` | customer_id, contact_data | contact_id | DB, Audit |
| 9 | `RemoveContact` | customer_id, contact_id | | DB, Audit |
| 10 | `SetPrimaryContact` | customer_id, contact_id | | DB, Audit |
| 11 | `AssignPackage` | customer_id, package_id | | DB, Kafka, Audit |
| 12 | `VerifyKYC` | customer_id, doc_type, doc_no | | DB, External, Audit |
| 13 | `SearchCustomers` | tenant_id, query | []CustomerSummary | ES read |

### A.8 API Contract

| Method | Path | Auth | Role | Description |
| :--- | :--- | :-: | :--- | :--- |
| GET | `/api/v1/customers` | ✅ | viewer | List + filter + pagination |
| GET | `/api/v1/customers/:id` | ✅ | viewer | Get by ID |
| POST | `/api/v1/customers` | ✅ | manager | Create |
| PUT | `/api/v1/customers/:id` | ✅ | manager | Full update |
| PATCH | `/api/v1/customers/:id` | ✅ | manager | Partial update |
| DELETE | `/api/v1/customers/:id` | ✅ | admin | Soft delete |
| POST | `/api/v1/customers/:id/activate` | ✅ | manager | Activate |
| POST | `/api/v1/customers/:id/suspend` | ✅ | manager | Suspend (body: reason) |
| POST | `/api/v1/customers/:id/churn` | ✅ | manager | Churn (body: reason) |
| POST | `/api/v1/customers/:id/contacts` | ✅ | manager | Add contact |
| DELETE | `/api/v1/customers/:id/contacts/:contactID` | ✅ | manager | Remove |
| POST | `/api/v1/customers/:id/contacts/:contactID/primary` | ✅ | manager | Set primary |
| POST | `/api/v1/customers/:id/package` | ✅ | manager | Assign package |
| POST | `/api/v1/customers/:id/kyc` | ✅ | admin | Verify KYC |
| GET | `/api/v1/customers/search?q=` | ✅ | viewer | Search (ES) |

**Request/Response envelope (ทั้ง 2 ภาษา):**
```json
// Success
{ "data": {...}, "meta": {"request_id": "...", "timestamp": "..."} }
// Error
{ "error": {"code": "CUSTOMER_NOT_FOUND", "message": "...", "details": {...}}, "meta": {...} }
```

### A.9 Database Schema (PostgreSQL)

```sql
-- customer_customers
CREATE TABLE customer_customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOT NULL,
    segment VARCHAR(30) NOT NULL,
    lifecycle VARCHAR(30) NOT NULL DEFAULT 'LEAD',
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    tax_id VARCHAR(20),
    address JSONB,
    package_id UUID,
    kyc JSONB,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT uq_customer_tenant_code UNIQUE (tenant_id, code)
);
CREATE INDEX idx_customer_customers_tenant ON customer_customers (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_customers_status ON customer_customers (tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_customers_segment ON customer_customers (tenant_id, segment) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_customers_lifecycle ON customer_customers (tenant_id, lifecycle) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_customers_package ON customer_customers (package_id) WHERE package_id IS NOT NULL;

-- customer_contacts
CREATE TABLE customer_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(255),
    position VARCHAR(100),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_customer_contacts_customer ON customer_contacts (customer_id);
CREATE UNIQUE INDEX uq_customer_contacts_primary ON customer_contacts (customer_id) WHERE is_primary = TRUE;

-- customer_sites (link to device context)
CREATE TABLE customer_sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_customers(id) ON DELETE CASCADE,
    site_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_customer_sites_customer ON customer_sites (customer_id);
CREATE INDEX idx_customer_sites_site ON customer_sites (site_id);
```

### A.10 Events (Kafka)

**Published:**
| Topic | Trigger |
| :--- | :--- |
| `customer.customer.created` | CreateCustomer |
| `customer.customer.updated` | UpdateCustomer |
| `customer.customer.activated` | ActivateCustomer |
| `customer.customer.suspended` | SuspendCustomer |
| `customer.customer.churned` | ChurnCustomer |
| `customer.contact.added` | AddContact |
| `customer.package.assigned` | AssignPackage |

**Consumed:**
| Topic | Handler |
| :--- | :--- |
| `package.subscription.created` | Update customer.package_id |
| `package.subscription.cancelled` | Clear customer.package_id |
| `payment.payment.succeeded` | Update lifecycle → CUSTOMER |

---

## PART B — Go Implementation (icmongolang)

### B.1 Directory Tree

```
internal/modules/customer/
├── domain/
│   ├── entity/
│   │   ├── customer.go
│   │   ├── contact.go
│   │   └── kyc.go
│   ├── value_object/
│   │   ├── customer_type.go
│   │   ├── customer_status.go
│   │   ├── segment.go
│   │   ├── lifecycle.go
│   │   ├── address.go
│   │   ├── email.go
│   │   ├── phone.go
│   │   └── tax_id.go
│   ├── repository/
│   │   ├── customer_repository.go
│   │   └── read_model_repository.go
│   ├── service/
│   │   ├── customer_service.go
│   │   └── kyc_verifier.go
│   ├── event/
│   │   └── events.go
│   └── errors/
│       └── errors.go
├── application/
│   ├── create_customer.go
│   ├── get_customer.go
│   ├── list_customers.go
│   ├── update_customer.go
│   ├── activate_customer.go
│   ├── suspend_customer.go
│   ├── churn_customer.go
│   ├── add_contact.go
│   ├── remove_contact.go
│   ├── set_primary_contact.go
│   ├── assign_package.go
│   ├── verify_kyc.go
│   ├── search_customers.go
│   ├── dto.go
│   └── ports.go
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── customer_repo_impl.go
│   │   │   ├── read_model_repo.go
│   │   │   └── models.go
│   │   └── redis/
│   │       └── customer_cache.go
│   ├── messaging/
│   │   └── producer.go
│   ├── search/
│   │   └── elasticsearch/
│   │       └── customer_indexer.go
│   └── services/
│       └── kyc/
│           └── dbd_provider.go
├── interfaces/
│   ├── http/
│   │   ├── customer_handler.go
│   │   ├── routes.go
│   │   └── dto.go
│   └── middleware/
│       └── tenant.go
└── module.go
```

### B.2 Domain — Value Objects

```go
// internal/modules/customer/domain/value_object/customer_type.go
package valueobject

import "errors"

type CustomerType string

const (
	CustomerTypePersonal   CustomerType = "PERSONAL"
	CustomerTypeBusiness   CustomerType = "BUSINESS"
	CustomerTypeGovernment CustomerType = "GOVERNMENT"
	CustomerTypeNGO        CustomerType = "NGO"
)

var ErrInvalidCustomerType = errors.New("invalid customer type")

func (t CustomerType) IsValid() bool {
	switch t {
	case CustomerTypePersonal, CustomerTypeBusiness, CustomerTypeGovernment, CustomerTypeNGO:
		return true
	}
	return false
}

func (t CustomerType) RequiresTaxID() bool {
	return t == CustomerTypeBusiness || t == CustomerTypeGovernment
}

func (t CustomerType) CanVerifyKYC() bool {
	return t == CustomerTypeBusiness || t == CustomerTypeGovernment
}
```

```go
// internal/modules/customer/domain/value_object/customer_status.go
package valueobject

type CustomerStatus string

const (
	CustomerStatusPending   CustomerStatus = "PENDING"
	CustomerStatusActive    CustomerStatus = "ACTIVE"
	CustomerStatusSuspended CustomerStatus = "SUSPENDED"
	CustomerStatusChurned   CustomerStatus = "CHURNED"
)

func (s CustomerStatus) IsValid() bool {
	switch s {
	case CustomerStatusPending, CustomerStatusActive, CustomerStatusSuspended, CustomerStatusChurned:
		return true
	}
	return false
}

func (s CustomerStatus) CanTransitionTo(next CustomerStatus) bool {
	switch s {
	case CustomerStatusPending:
		return next == CustomerStatusActive
	case CustomerStatusActive:
		return next == CustomerStatusSuspended || next == CustomerStatusChurned
	case CustomerStatusSuspended:
		return next == CustomerStatusActive || next == CustomerStatusChurned
	case CustomerStatusChurned:
		return false
	}
	return false
}
```

```go
// internal/modules/customer/domain/value_object/segment.go
package valueobject

type Segment string

const (
	SegmentSME        Segment = "SME"
	SegmentEnterprise Segment = "ENTERPRISE"
	SegmentGovernment Segment = "GOVERNMENT"
	SegmentStartup    Segment = "STARTUP"
	SegmentIndividual Segment = "INDIVIDUAL"
)

func (s Segment) IsValid() bool {
	switch s {
	case SegmentSME, SegmentEnterprise, SegmentGovernment, SegmentStartup, SegmentIndividual:
		return true
	}
	return false
}
```

```go
// internal/modules/customer/domain/value_object/lifecycle.go
package valueobject

type LifecycleStage string

const (
	LifecycleLead     LifecycleStage = "LEAD"
	LifecycleProspect LifecycleStage = "PROSPECT"
	LifecycleCustomer LifecycleStage = "CUSTOMER"
	LifecycleChurned  LifecycleStage = "CHURNED"
)
```

```go
// internal/modules/customer/domain/value_object/address.go
package valueobject

import (
	"errors"
	"regexp"
)

var postcodePattern = regexp.MustCompile(`^\d{5}$`)

type Address struct {
	Line1    string   `json:"line1"`
	Line2    string   `json:"line2,omitempty"`
	District string   `json:"district"`
	Amphoe   string   `json:"amphoe"`
	Province string   `json:"province"`
	Postcode string   `json:"postcode"`
	Country  string   `json:"country"`
	Lat      *float64 `json:"lat,omitempty"`
	Lng      *float64 `json:"lng,omitempty"`
}

func (a Address) Validate() error {
	if len(a.Line1) > 255 || len(a.Line2) > 255 {
		return errors.New("address line too long")
	}
	if a.Postcode != "" && !postcodePattern.MatchString(a.Postcode) {
		return errors.New("invalid postcode (must be 5 digits)")
	}
	if a.Country == "" {
		return errors.New("country required")
	}
	return nil
}

func (a Address) IsEmpty() bool {
	return a.Line1 == "" && a.Province == "" && a.Country == ""
}
```

```go
// internal/modules/customer/domain/value_object/email.go
package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email string

func NewEmail(s string) (Email, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return "", nil
	}
	if !emailPattern.MatchString(s) {
		return "", errors.New("invalid email format")
	}
	return Email(s), nil
}

func (e Email) String() string { return string(e) }
func (e Email) IsEmpty() bool  { return e == "" }
```

```go
// internal/modules/customer/domain/value_object/tax_id.go
package valueobject

import (
	"errors"
	"regexp"
)

var taxIDPattern = regexp.MustCompile(`^\d{13}$`)

type TaxID string

func NewTaxID(s string) (TaxID, error) {
	if s == "" {
		return "", nil
	}
	if !taxIDPattern.MatchString(s) {
		return "", errors.New("invalid tax id (must be 13 digits)")
	}
	return TaxID(s), nil
}

func (t TaxID) String() string { return string(t) }
func (t TaxID) IsEmpty() bool  { return t == "" }
```

### B.3 Domain — Entities

```go
// internal/modules/customer/domain/entity/contact.go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

type Contact struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

func NewContact(name, phone, email, position string, isPrimary bool) (Contact, error) {
	if name == "" || len(name) > 255 {
		return Contact{}, domainerrors.ErrInvalidContactName
	}

	if email != "" {
		if _, err := valueobject.NewEmail(email); err != nil {
			return Contact{}, err
		}
	}

	return Contact{
		ID:        uuid.New(),
		Name:      name,
		Phone:     phone,
		Email:     email,
		Position:  position,
		IsPrimary: isPrimary,
		CreatedAt: time.Now(),
	}, nil
}
```

```go
// internal/modules/customer/domain/entity/kyc.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type KYCInfo struct {
	Verified     bool       `json:"verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	VerifiedBy   *uuid.UUID `json:"verified_by,omitempty"`
	DocumentType string     `json:"document_type,omitempty"`
	DocumentNo   string     `json:"document_no,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

func (k KYCInfo) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

func (k KYCInfo) IsValid() bool {
	return k.Verified && !k.IsExpired()
}
```

```go
// internal/modules/customer/domain/entity/customer.go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

const MaxContacts = 20

// Customer – Aggregate Root
type Customer struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Code      string
	Name      string
	Type      valueobject.CustomerType
	Segment   valueobject.Segment
	Lifecycle valueobject.LifecycleStage
	Status    valueobject.CustomerStatus
	TaxID     string
	Address   valueobject.Address
	Contacts  []Contact
	PackageID *uuid.UUID
	KYC       *KYCInfo
	Metadata  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCustomer – Constructor บังคับ invariants
func NewCustomer(
	tenantID uuid.UUID,
	code, name string,
	cType valueobject.CustomerType,
	segment valueobject.Segment,
) (*Customer, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if !isValidCode(code) {
		return nil, domainerrors.ErrInvalidCode
	}
	if name == "" || len(name) > 255 {
		return nil, domainerrors.ErrInvalidName
	}
	if !cType.IsValid() {
		return nil, domainerrors.ErrInvalidType
	}
	if !segment.IsValid() {
		return nil, domainerrors.ErrInvalidSegment
	}

	now := time.Now()
	return &Customer{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		Type:      cType,
		Segment:   segment,
		Lifecycle: valueobject.LifecycleLead,
		Status:    valueobject.CustomerStatusPending,
		Contacts:  []Contact{},
		Metadata:  map[string]interface{}{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ---------- Behavior methods ----------

func (c *Customer) Rename(newName string) error {
	if newName == "" || len(newName) > 255 {
		return domainerrors.ErrInvalidName
	}
	c.Name = newName
	c.touch()
	return nil
}

func (c *Customer) SetTaxID(taxID string) error {
	if c.Type.RequiresTaxID() && taxID == "" {
		return domainerrors.ErrTaxIDRequired
	}
	if taxID != "" {
		if _, err := valueobject.NewTaxID(taxID); err != nil {
			return err
		}
	}
	c.TaxID = taxID
	c.touch()
	return nil
}

func (c *Customer) UpdateAddress(addr valueobject.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}
	c.Address = addr
	c.touch()
	return nil
}

func (c *Customer) ChangeSegment(segment valueobject.Segment) error {
	if !segment.IsValid() {
		return domainerrors.ErrInvalidSegment
	}
	c.Segment = segment
	c.touch()
	return nil
}

// Activate – เปลี่ยน status → ACTIVE
func (c *Customer) Activate() error {
	if !c.Status.CanTransitionTo(valueobject.CustomerStatusActive) {
		return domainerrors.ErrInvalidStatusTransition
	}
	// invariant: BUSINESS/GOVERNMENT ต้องมี tax_id
	if c.Type.RequiresTaxID() && c.TaxID == "" {
		return domainerrors.ErrTaxIDRequired
	}
	c.Status = valueobject.CustomerStatusActive
	c.Lifecycle = valueobject.LifecycleCustomer
	c.touch()
	return nil
}

func (c *Customer) Suspend(reason string) error {
	if !c.Status.CanTransitionTo(valueobject.CustomerStatusSuspended) {
		return domainerrors.ErrInvalidStatusTransition
	}
	if reason == "" {
		return domainerrors.ErrReasonRequired
	}
	c.Status = valueobject.CustomerStatusSuspended
	c.Metadata["suspend_reason"] = reason
	c.Metadata["suspended_at"] = time.Now()
	c.touch()
	return nil
}

func (c *Customer) Churn(reason string) error {
	if !c.Status.CanTransitionTo(valueobject.CustomerStatusChurned) {
		return domainerrors.ErrInvalidStatusTransition
	}
	c.Status = valueobject.CustomerStatusChurned
	c.Lifecycle = valueobject.LifecycleChurned
	c.Metadata["churn_reason"] = reason
	c.Metadata["churned_at"] = time.Now()
	c.touch()
	return nil
}

func (c *Customer) AddContact(contact Contact) error {
	if len(c.Contacts) >= MaxContacts {
		return domainerrors.ErrTooManyContacts
	}
	if contact.IsPrimary {
		c.clearPrimary()
	}
	c.Contacts = append(c.Contacts, contact)
	c.touch()
	return nil
}

func (c *Customer) RemoveContact(contactID uuid.UUID) error {
	for i, ct := range c.Contacts {
		if ct.ID == contactID {
			c.Contacts = append(c.Contacts[:i], c.Contacts[i+1:]...)
			c.touch()
			return nil
		}
	}
	return domainerrors.ErrContactNotFound
}

func (c *Customer) SetPrimaryContact(contactID uuid.UUID) error {
	found := false
	for i := range c.Contacts {
		if c.Contacts[i].ID == contactID {
			c.Contacts[i].IsPrimary = true
			found = true
		} else {
			c.Contacts[i].IsPrimary = false
		}
	}
	if !found {
		return domainerrors.ErrContactNotFound
	}
	c.touch()
	return nil
}

func (c *Customer) AssignPackage(pkgID uuid.UUID) error {
	if pkgID == uuid.Nil {
		return domainerrors.ErrInvalidPackage
	}
	c.PackageID = &pkgID
	c.touch()
	return nil
}

func (c *Customer) AttachKYC(kyc KYCInfo) error {
	if !c.Type.CanVerifyKYC() {
		return domainerrors.ErrKYCNotAllowed
	}
	c.KYC = &kyc
	c.touch()
	return nil
}

// ---------- Query methods ----------

func (c *Customer) IsActive() bool {
	return c.Status == valueobject.CustomerStatusActive
}

func (c *Customer) IsBusiness() bool {
	return c.Type == valueobject.CustomerTypeBusiness
}

func (c *Customer) HasKYC() bool {
	return c.KYC != nil && c.KYC.IsValid()
}

func (c *Customer) PrimaryContact() *Contact {
	for i := range c.Contacts {
		if c.Contacts[i].IsPrimary {
			return &c.Contacts[i]
		}
	}
	return nil
}

// ---------- Private ----------

func (c *Customer) touch() { c.UpdatedAt = time.Now() }

func (c *Customer) clearPrimary() {
	for i := range c.Contacts {
		c.Contacts[i].IsPrimary = false
	}
}

func isValidCode(code string) bool {
	if len(code) < 3 || len(code) > 50 {
		return false
	}
	for _, r := range code {
		if !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}
```

### B.4 Domain — Repository Interface

```go
// internal/modules/customer/domain/repository/customer_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/customer/domain/entity"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

type CustomerRepository interface {
	Save(ctx context.Context, c *entity.Customer) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Customer, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Customer, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.CustomerStatus) ([]entity.Customer, error)
	ExistsByCode(ctx context.Context, tenantID uuid.UUID, code string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Read Model (CQRS)
type CustomerReadModel interface {
	GetDetail(ctx context.Context, id uuid.UUID) (*CustomerDetail, error)
	ListSummaries(ctx context.Context, filter CustomerFilter) ([]CustomerSummary, int64, error)
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit int) ([]CustomerSummary, error)
}

type CustomerDetail struct {
	Customer     *entity.Customer
	DeviceCount  int
	SiteCount    int
	OpenTickets  int
	TotalRevenue float64
}

type CustomerFilter struct {
	TenantID uuid.UUID
	Status   *valueobject.CustomerStatus
	Segment  *valueobject.Segment
	Search   string
	Page     int
	Size     int
	SortBy   string
	SortDesc bool
}

type CustomerSummary struct {
	ID           uuid.UUID
	Code         string
	Name         string
	Type         valueobject.CustomerType
	Segment      valueobject.Segment
	Status       valueobject.CustomerStatus
	DeviceCount  int
	TotalRevenue float64
	CreatedAt    time.Time
}
```

### B.5 Domain — Errors

```go
// internal/modules/customer/domain/errors/errors.go
package domainerrors

import "errors"

var (
	// Validation
	ErrInvalidTenant   = errors.New("invalid tenant")
	ErrInvalidCode     = errors.New("invalid customer code")
	ErrInvalidName     = errors.New("invalid customer name")
	ErrInvalidType     = errors.New("invalid customer type")
	ErrInvalidSegment  = errors.New("invalid segment")
	ErrInvalidPackage  = errors.New("invalid package")
	ErrInvalidContactName = errors.New("invalid contact name")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidTaxID    = errors.New("invalid tax id")
	ErrTaxIDRequired   = errors.New("tax id required for this customer type")
	ErrReasonRequired  = errors.New("reason is required")

	// Not found
	ErrCustomerNotFound = errors.New("customer not found")
	ErrContactNotFound  = errors.New("contact not found")

	// Conflict
	ErrCodeDuplicate = errors.New("customer code already exists")

	// State machine
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// Business rules
	ErrTooManyContacts = errors.New("too many contacts (max 20)")
	ErrKYCNotAllowed   = errors.New("KYC only allowed for BUSINESS/GOVERNMENT")
	ErrKYCVerificationFailed = errors.New("KYC verification failed")
	ErrCustomerNotActive = errors.New("customer is not active")
)
```

### B.6 Application — Ports

```go
// internal/modules/customer/application/ports.go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/customer/domain/event"
)

// EventProducer – outbound port สำหรับ publish event
type EventProducer interface {
	PublishCustomerCreated(ctx context.Context, evt event.CustomerCreated) error
	PublishCustomerUpdated(ctx context.Context, evt event.CustomerUpdated) error
	PublishCustomerActivated(ctx context.Context, evt event.CustomerActivated) error
	PublishCustomerSuspended(ctx context.Context, evt event.CustomerSuspended) error
	PublishCustomerChurned(ctx context.Context, evt event.CustomerChurned) error
	PublishPackageAssigned(ctx context.Context, evt event.PackageAssigned) error
}

// AuditRepository – audit trail
type AuditRepository interface {
	Save(ctx context.Context, trail *AuditTrail) error
}

type AuditTrail struct {
	UserID     uuid.UUID
	Action     string
	EntityType string
	EntityID   uuid.UUID
	Payload    map[string]interface{}
	IPAddress  string
	UserAgent  string
	OccurredAt time.Time
}

// WSHub – realtime broadcast
type WSHub interface {
	BroadcastToTenant(tenantID string, event interface{})
}

// PackageClient – query package (outbound)
type PackageClient interface {
	Exists(ctx context.Context, packageID uuid.UUID) (bool, error)
}

// Cache – cache port
type Cache interface {
	SetCustomer(ctx context.Context, c *entity.Customer) error
	GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
}

// Indexer – ES port
type Indexer interface {
	Index(ctx context.Context, c *entity.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}
```

### B.7 Application — Use Case: CreateCustomer

```go
// internal/modules/customer/application/create_customer.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/customer/domain/entity"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	"icmongolang/internal/modules/customer/domain/event"
	"icmongolang/internal/modules/customer/domain/repository"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

type CreateCustomerUseCase struct {
	repo      repository.CustomerRepository
	auditRepo AuditRepository
	producer  EventProducer
	indexer   Indexer
}

func NewCreateCustomerUseCase(
	repo repository.CustomerRepository,
	auditRepo AuditRepository,
	producer EventProducer,
	indexer Indexer,
) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		repo:      repo,
		auditRepo: auditRepo,
		producer:  producer,
		indexer:   indexer,
	}
}

type CreateCustomerInput struct {
	TenantID  uuid.UUID
	Code      string
	Name      string
	Type      valueobject.CustomerType
	Segment   valueobject.Segment
	TaxID     string
	Address   valueobject.Address
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
}

type CreateCustomerOutput struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Code       string    `json:"code"`
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, in CreateCustomerInput) (*CreateCustomerOutput, error) {
	// 1. Check duplicate
	exists, err := uc.repo.ExistsByCode(ctx, in.TenantID, in.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrCodeDuplicate
	}

	// 2. Create aggregate
	c, err := entity.NewCustomer(in.TenantID, in.Code, in.Name, in.Type, in.Segment)
	if err != nil {
		return nil, err
	}
	if err := c.SetTaxID(in.TaxID); err != nil {
		return nil, err
	}
	if err := c.UpdateAddress(in.Address); err != nil {
		return nil, err
	}

	// 3. Persist
	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, err
	}

	// 4. Audit
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID:     in.UserID,
		Action:     "CUSTOMER_CREATED",
		EntityType: "customer",
		EntityID:   c.ID,
		Payload: map[string]interface{}{
			"code": c.Code,
			"name": c.Name,
			"type": string(c.Type),
		},
		IPAddress:  in.IPAddress,
		UserAgent:  in.UserAgent,
		OccurredAt: time.Now(),
	})

	// 5. Index for search
	_ = uc.indexer.Index(ctx, c)

	// 6. Publish event
	_ = uc.producer.PublishCustomerCreated(ctx, event.CustomerCreated{
		CustomerID: c.ID,
		TenantID:   c.TenantID,
		Code:       c.Code,
		Name:       c.Name,
		Type:       string(c.Type),
		OccurredAt: time.Now(),
	})

	return &CreateCustomerOutput{CustomerID: c.ID, Code: c.Code}, nil
}
```

### B.8 Application — Use Case: SuspendCustomer

```go
// internal/modules/customer/application/suspend_customer.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	"icmongolang/internal/modules/customer/domain/event"
	"icmongolang/internal/modules/customer/domain/repository"
)

type SuspendCustomerUseCase struct {
	repo      repository.CustomerRepository
	auditRepo AuditRepository
	producer  EventProducer
	wsHub     WSHub
	cache     Cache
}

func NewSuspendCustomerUseCase(
	repo repository.CustomerRepository,
	auditRepo AuditRepository,
	producer EventProducer,
	wsHub WSHub,
	cache Cache,
) *SuspendCustomerUseCase {
	return &SuspendCustomerUseCase{repo: repo, auditRepo: auditRepo, producer: producer, wsHub: wsHub, cache: cache}
}

type SuspendCustomerInput struct {
	CustomerID uuid.UUID
	Reason     string
	UserID     uuid.UUID
	IPAddress  string
}

func (uc *SuspendCustomerUseCase) Execute(ctx context.Context, in SuspendCustomerInput) error {
	// 1. Load aggregate
	c, err := uc.repo.FindByID(ctx, in.CustomerID)
	if err != nil {
		return err
	}

	// 2. Call behavior
	if err := c.Suspend(in.Reason); err != nil {
		return err
	}

	// 3. Persist
	if err := uc.repo.Save(ctx, c); err != nil {
		return err
	}

	// 4. Invalidate cache
	_ = uc.cache.DeleteCustomer(ctx, c.ID)

	// 5. Audit
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID:     in.UserID,
		Action:     "CUSTOMER_SUSPENDED",
		EntityType: "customer",
		EntityID:   c.ID,
		Payload:    map[string]interface{}{"reason": in.Reason},
		IPAddress:  in.IPAddress,
		OccurredAt: time.Now(),
	})

	// 6. Publish
	_ = uc.producer.PublishCustomerSuspended(ctx, event.CustomerSuspended{
		CustomerID: c.ID,
		TenantID:   c.TenantID,
		Reason:     in.Reason,
		OccurredAt: time.Now(),
	})

	// 7. WebSocket
	uc.wsHub.BroadcastToTenant(c.TenantID.String(), map[string]interface{}{
		"type": "customer.suspended",
		"data": map[string]interface{}{"customer_id": c.ID.String(), "reason": in.Reason},
	})

	// 8. Check ErrCustomerNotActive guard
	_ = domainerrors.ErrCustomerNotActive
	return nil
}
```

### B.9 Infrastructure — Postgres

```go
// internal/modules/customer/infrastructure/persistence/postgres/models.go
package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CustomerModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index"`
	Code      string         `gorm:"type:varchar(50);not null"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Type      string         `gorm:"type:varchar(30);not null"`
	Segment   string         `gorm:"type:varchar(30);not null"`
	Lifecycle string         `gorm:"type:varchar(30);not null"`
	Status    string         `gorm:"type:varchar(30);not null;index"`
	TaxID     string         `gorm:"type:varchar(20)"`
	Address   datatypes.JSON `gorm:"type:jsonb"`
	PackageID *uuid.UUID
	KYC       datatypes.JSON `gorm:"type:jsonb"`
	Metadata  datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	Contacts []ContactModel `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

func (CustomerModel) TableName() string { return "customer_customers" }

type ContactModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CustomerID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Phone      string    `gorm:"type:varchar(30)"`
	Email      string    `gorm:"type:varchar(255)"`
	Position   string    `gorm:"type:varchar(100)"`
	IsPrimary  bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time
}

func (ContactModel) TableName() string { return "customer_contacts" }
```

```go
// internal/modules/customer/infrastructure/persistence/postgres/customer_repo_impl.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/customer/domain/entity"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

type customerRepoImpl struct{ db *gorm.DB }

func NewCustomerRepository(db *gorm.DB) *customerRepoImpl {
	return &customerRepoImpl{db: db}
}

func (r *customerRepoImpl) Save(ctx context.Context, c *entity.Customer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := toModel(c)
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		// Sync contacts: delete all, re-insert (simple strategy)
		if err := tx.Where("customer_id = ?", c.ID).Delete(&ContactModel{}).Error; err != nil {
			return err
		}
		for _, ct := range c.Contacts {
			cm := ContactModel{
				ID:         ct.ID,
				CustomerID: c.ID,
				Name:       ct.Name,
				Phone:      ct.Phone,
				Email:      ct.Email,
				Position:   ct.Position,
				IsPrimary:  ct.IsPrimary,
				CreatedAt:  ct.CreatedAt,
			}
			if err := tx.Create(&cm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *customerRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	var m CustomerModel
	err := r.db.WithContext(ctx).
		Preload("Contacts").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *customerRepoImpl) FindByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Customer, error) {
	var m CustomerModel
	err := r.db.WithContext(ctx).
		Preload("Contacts").
		Where("tenant_id = ? AND code = ? AND deleted_at IS NULL", tenantID, code).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *customerRepoImpl) FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Customer, error) {
	var models []CustomerModel
	if err := r.db.WithContext(ctx).
		Preload("Contacts").
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toEntities(models), nil
}

func (r *customerRepoImpl) FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.CustomerStatus) ([]entity.Customer, error) {
	var models []CustomerModel
	if err := r.db.WithContext(ctx).
		Preload("Contacts").
		Where("tenant_id = ? AND status = ? AND deleted_at IS NULL", tenantID, status).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return toEntities(models), nil
}

func (r *customerRepoImpl) ExistsByCode(ctx context.Context, tenantID uuid.UUID, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&CustomerModel{}).
		Where("tenant_id = ? AND code = ? AND deleted_at IS NULL", tenantID, code).
		Count(&count).Error
	return count > 0, err
}

func (r *customerRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ---------- Mappers ----------

func toModel(c *entity.Customer) *CustomerModel {
	addrJSON, _ := json.Marshal(c.Address)
	kycJSON, _ := json.Marshal(c.KYC)
	metaJSON, _ := json.Marshal(c.Metadata)

	return &CustomerModel{
		ID:        c.ID,
		TenantID:  c.TenantID,
		Code:      c.Code,
		Name:      c.Name,
		Type:      string(c.Type),
		Segment:   string(c.Segment),
		Lifecycle: string(c.Lifecycle),
		Status:    string(c.Status),
		TaxID:     c.TaxID,
		Address:   datatypes.JSON(addrJSON),
		PackageID: c.PackageID,
		KYC:       datatypes.JSON(kycJSON),
		Metadata:  datatypes.JSON(metaJSON),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toEntity(m *CustomerModel) *entity.Customer {
	var addr valueobject.Address
	var kyc entity.KYCInfo
	var meta map[string]interface{}

	_ = json.Unmarshal(m.Address, &addr)
	if len(m.KYC) > 0 {
		_ = json.Unmarshal(m.KYC, &kyc)
	}
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = map[string]interface{}{}
	}

	c := &entity.Customer{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Code:      m.Code,
		Name:      m.Name,
		Type:      valueobject.CustomerType(m.Type),
		Segment:   valueobject.Segment(m.Segment),
		Lifecycle: valueobject.LifecycleStage(m.Lifecycle),
		Status:    valueobject.CustomerStatus(m.Status),
		TaxID:     m.TaxID,
		Address:   addr,
		PackageID: m.PackageID,
		Metadata:  meta,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if kyc.Verified {
		c.KYC = &kyc
	}
	for _, ct := range m.Contacts {
		c.Contacts = append(c.Contacts, entity.Contact{
			ID:        ct.ID,
			Name:      ct.Name,
			Phone:     ct.Phone,
			Email:     ct.Email,
			Position:  ct.Position,
			IsPrimary: ct.IsPrimary,
			CreatedAt: ct.CreatedAt,
		})
	}
	return c
}

func toEntities(models []CustomerModel) []entity.Customer {
	out := make([]entity.Customer, 0, len(models))
	for i := range models {
		out = append(out, *toEntity(&models[i]))
	}
	return out
}
```

### B.10 Infrastructure — Kafka Producer

```go
// internal/modules/customer/infrastructure/messaging/producer.go
package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/customer/domain/event"
)

const (
	TopicCustomerCreated   = "customer.customer.created"
	TopicCustomerUpdated   = "customer.customer.updated"
	TopicCustomerActivated = "customer.customer.activated"
	TopicCustomerSuspended = "customer.customer.suspended"
	TopicCustomerChurned   = "customer.customer.churned"
	TopicPackageAssigned   = "customer.package.assigned"
)

type kafkaProducer struct {
	sync sarama.SyncProducer
}

func NewKafkaProducer(brokers []string, clientID string) (*kafkaProducer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{sync: p}, nil
}

func (k *kafkaProducer) publish(ctx context.Context, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, _, err = k.sync.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (k *kafkaProducer) PublishCustomerCreated(ctx context.Context, evt event.CustomerCreated) error {
	return k.publish(ctx, TopicCustomerCreated, evt.CustomerID.String(), evt)
}

func (k *kafkaProducer) PublishCustomerUpdated(ctx context.Context, evt event.CustomerUpdated) error {
	return k.publish(ctx, TopicCustomerUpdated, evt.CustomerID.String(), evt)
}

func (k *kafkaProducer) PublishCustomerActivated(ctx context.Context, evt event.CustomerActivated) error {
	return k.publish(ctx, TopicCustomerActivated, evt.CustomerID.String(), evt)
}

func (k *kafkaProducer) PublishCustomerSuspended(ctx context.Context, evt event.CustomerSuspended) error {
	return k.publish(ctx, TopicCustomerSuspended, evt.CustomerID.String(), evt)
}

func (k *kafkaProducer) PublishCustomerChurned(ctx context.Context, evt event.CustomerChurned) error {
	return k.publish(ctx, TopicCustomerChurned, evt.CustomerID.String(), evt)
}

func (k *kafkaProducer) PublishPackageAssigned(ctx context.Context, evt event.PackageAssigned) error {
	return k.publish(ctx, TopicPackageAssigned, evt.CustomerID.String(), evt)
}

var _ = fmt.Sprintf
```

### B.11 Interface — HTTP Handler

```go
// internal/modules/customer/interfaces/http/customer_handler.go
package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/customer/application"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
	"icmongolang/pkg/responses"
)

type CustomerHandler struct {
	createUC  *application.CreateCustomerUseCase
	getUC     *application.GetCustomerUseCase
	listUC    *application.ListCustomersUseCase
	updateUC  *application.UpdateCustomerUseCase
	activateUC *application.ActivateCustomerUseCase
	suspendUC *application.SuspendCustomerUseCase
	churnUC   *application.ChurnCustomerUseCase
	addCtUC   *application.AddContactUseCase
	removeCtUC *application.RemoveContactUseCase
	searchUC  *application.SearchCustomersUseCase
}

func NewCustomerHandler(
	createUC *application.CreateCustomerUseCase,
	getUC *application.GetCustomerUseCase,
	listUC *application.ListCustomersUseCase,
	updateUC *application.UpdateCustomerUseCase,
	activateUC *application.ActivateCustomerUseCase,
	suspendUC *application.SuspendCustomerUseCase,
	churnUC *application.ChurnCustomerUseCase,
	addCtUC *application.AddContactUseCase,
	removeCtUC *application.RemoveContactUseCase,
	searchUC *application.SearchCustomersUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createUC: createUC, getUC: getUC, listUC: listUC, updateUC: updateUC,
		activateUC: activateUC, suspendUC: suspendUC, churnUC: churnUC,
		addCtUC: addCtUC, removeCtUC: removeCtUC, searchUC: searchUC,
	}
}

// POST /api/v1/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	out, err := h.createUC.Execute(c.Request.Context(), application.CreateCustomerInput{
		TenantID:  tenantID,
		Code:      req.Code,
		Name:      req.Name,
		Type:      valueobject.CustomerType(req.Type),
		Segment:   valueobject.Segment(req.Segment),
		TaxID:     req.TaxID,
		Address:   req.Address.ToVO(),
		UserID:    userID,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

// GET /api/v1/customers/:id
func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid customer id", nil))
		return
	}
	out, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(out))
}

// POST /api/v1/customers/:id/suspend
func (h *CustomerHandler) Suspend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid customer id", nil))
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	if err := h.suspendUC.Execute(c.Request.Context(), application.SuspendCustomerInput{
		CustomerID: id,
		Reason:     req.Reason,
		UserID:     userID,
		IPAddress:  c.ClientIP(),
	}); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "suspended"}))
}

// GET /api/v1/customers
func (h *CustomerHandler) List(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }
	if size < 1 || size > 100 { size = 20 }

	var status *valueobject.CustomerStatus
	if s := c.Query("status"); s != "" {
		v := valueobject.CustomerStatus(s)
		status = &v
	}
	var segment *valueobject.Segment
	if s := c.Query("segment"); s != "" {
		v := valueobject.Segment(s)
		segment = &v
	}

	out, total, err := h.listUC.Execute(c.Request.Context(), application.ListCustomersInput{
		TenantID: tenantID,
		Status:   status,
		Segment:  segment,
		Search:   c.Query("q"),
		Page:     page,
		Size:     size,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Paginated(out, total, page, size))
}
```

```go
// internal/modules/customer/interfaces/http/routes.go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
	Customer *CustomerHandler
}

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handlers,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	g := r.Group("/customers")
	g.Use(auth, tenant)

	g.GET("", rbac("viewer"), h.Customer.List)
	g.GET("/search", rbac("viewer"), h.Customer.Search)
	g.GET("/:id", rbac("viewer"), h.Customer.Get)

	g.POST("", rbac("manager"), h.Customer.Create)
	g.PUT("/:id", rbac("manager"), h.Customer.Update)
	g.PATCH("/:id", rbac("manager"), h.Customer.Update)
	g.DELETE("/:id", rbac("admin"), h.Customer.Delete)

	g.POST("/:id/activate", rbac("manager"), h.Customer.Activate)
	g.POST("/:id/suspend", rbac("manager"), h.Customer.Suspend)
	g.POST("/:id/churn", rbac("manager"), h.Customer.Churn)

	g.POST("/:id/contacts", rbac("manager"), h.Customer.AddContact)
	g.DELETE("/:id/contacts/:contactID", rbac("manager"), h.Customer.RemoveContact)
	g.POST("/:id/contacts/:contactID/primary", rbac("manager"), h.Customer.SetPrimaryContact)

	g.POST("/:id/package", rbac("manager"), h.Customer.AssignPackage)
	g.POST("/:id/kyc", rbac("admin"), h.Customer.VerifyKYC)
}
```

### B.12 Composition Root — `module.go`

```go
// internal/modules/customer/module.go
package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/customer/application"
	"icmongolang/internal/modules/customer/infrastructure/messaging"
	pgrepo "icmongolang/internal/modules/customer/infrastructure/persistence/postgres"
	redisrepo "icmongolang/internal/modules/customer/infrastructure/persistence/redis"
	esindexer "icmongolang/internal/modules/customer/infrastructure/search/elasticsearch"
	httpiface "icmongolang/internal/modules/customer/interfaces/http"
)

type Dependencies struct {
	DB           *gorm.DB
	Redis        *redis.Client
	Producer     *messaging.KafkaProducer
	ESClient     ESClient
	AuditRepo    application.AuditRepository
	WSHub        application.WSHub
	PackageClient application.PackageClient
}

// Init – Wire-up ทั้งหมด
func Init(
	router *gin.RouterGroup,
	deps Dependencies,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	// Repositories
	customerRepo := pgrepo.NewCustomerRepository(deps.DB)
	readModel := pgrepo.NewCustomerReadModel(deps.DB)
	cache := redisrepo.NewCustomerCache(deps.Redis)
	indexer := esindexer.NewCustomerIndexer(deps.ESClient)

	// Use cases
	createUC := application.NewCreateCustomerUseCase(customerRepo, deps.AuditRepo, deps.Producer, indexer)
	getUC := application.NewGetCustomerUseCase(readModel, cache)
	listUC := application.NewListCustomersUseCase(readModel)
	updateUC := application.NewUpdateCustomerUseCase(customerRepo, deps.AuditRepo, deps.Producer, cache)
	activateUC := application.NewActivateCustomerUseCase(customerRepo, deps.AuditRepo, deps.Producer, deps.WSHub, cache)
	suspendUC := application.NewSuspendCustomerUseCase(customerRepo, deps.AuditRepo, deps.Producer, deps.WSHub, cache)
	churnUC := application.NewChurnCustomerUseCase(customerRepo, deps.AuditRepo, deps.Producer, cache)
	addCtUC := application.NewAddContactUseCase(customerRepo, deps.AuditRepo)
	removeCtUC := application.NewRemoveContactUseCase(customerRepo, deps.AuditRepo)
	searchUC := application.NewSearchCustomersUseCase(readModel)

	// Handlers
	handler := httpiface.NewCustomerHandler(
		createUC, getUC, listUC, updateUC, activateUC, suspendUC, churnUC,
		addCtUC, removeCtUC, searchUC,
	)

	// Routes
	httpiface.RegisterRoutes(router, &httpiface.Handlers{Customer: handler}, auth, tenant, rbac)
}
```

### B.13 Go Testing

```go
// internal/modules/customer/domain/entity/customer_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"icmongolang/internal/modules/customer/domain/entity"
	domainerrors "icmongolang/internal/modules/customer/domain/errors"
	valueobject "icmongolang/internal/modules/customer/domain/value_object"
)

func TestCustomer_StatusTransitions(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *entity.Customer
		action  func(c *entity.Customer) error
		wantErr error
	}{
		{
			name: "PENDING → ACTIVE (personal, ok)",
			setup: func() *entity.Customer {
				c, _ := entity.NewCustomer(uuid.New(), "C001", "Alice",
					valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)
				return c
			},
			action:  func(c *entity.Customer) error { return c.Activate() },
			wantErr: nil,
		},
		{
			name: "PENDING → ACTIVE (business without tax_id, error)",
			setup: func() *entity.Customer {
				c, _ := entity.NewCustomer(uuid.New(), "C002", "ACME",
					valueobject.CustomerTypeBusiness, valueobject.SegmentEnterprise)
				return c
			},
			action:  func(c *entity.Customer) error { return c.Activate() },
			wantErr: domainerrors.ErrTaxIDRequired,
		},
		{
			name: "ACTIVE → SUSPENDED (ok)",
			setup: func() *entity.Customer {
				c, _ := entity.NewCustomer(uuid.New(), "C003", "Bob",
					valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)
				_ = c.Activate()
				return c
			},
			action:  func(c *entity.Customer) error { return c.Suspend("payment failed") },
			wantErr: nil,
		},
		{
			name: "PENDING → SUSPENDED (invalid transition)",
			setup: func() *entity.Customer {
				c, _ := entity.NewCustomer(uuid.New(), "C004", "Eve",
					valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)
				return c
			},
			action:  func(c *entity.Customer) error { return c.Suspend("no reason") },
			wantErr: domainerrors.ErrInvalidStatusTransition,
		},
		{
			name: "CHURNED → ACTIVE (terminal, error)",
			setup: func() *entity.Customer {
				c, _ := entity.NewCustomer(uuid.New(), "C005", "Dan",
					valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)
				_ = c.Activate()
				_ = c.Churn("moved away")
				return c
			},
			action:  func(c *entity.Customer) error { return c.Activate() },
			wantErr: domainerrors.ErrInvalidStatusTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup()
			err := tt.action(c)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomer_AddContact_MaxLimit(t *testing.T) {
	c, _ := entity.NewCustomer(uuid.New(), "C010", "Test",
		valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)

	for i := 0; i < entity.MaxContacts; i++ {
		ct, _ := entity.NewContact("Contact", "", "", "", false)
		if err := c.AddContact(ct); err != nil {
			t.Fatalf("add #%d failed: %v", i, err)
		}
	}

	extra, _ := entity.NewContact("Extra", "", "", "", false)
	err := c.AddContact(extra)
	if !errors.Is(err, domainerrors.ErrTooManyContacts) {
		t.Errorf("expected ErrTooManyContacts, got %v", err)
	}
}

func TestCustomer_SetPrimaryContact_Uniqueness(t *testing.T) {
	c, _ := entity.NewCustomer(uuid.New(), "C011", "Test",
		valueobject.CustomerTypePersonal, valueobject.SegmentIndividual)

	ct1, _ := entity.NewContact("Alice", "", "", "", true)
	ct2, _ := entity.NewContact("Bob", "", "", "", false)
	_ = c.AddContact(ct1)
	_ = c.AddContact(ct2)

	_ = c.SetPrimaryContact(ct2.ID)

	primaryCount := 0
	for _, ct := range c.Contacts {
		if ct.IsPrimary {
			primaryCount++
		}
	}
	if primaryCount != 1 {
		t.Errorf("expected exactly 1 primary contact, got %d", primaryCount)
	}
	if c.PrimaryContact().ID != ct2.ID {
		t.Errorf("expected Bob to be primary")
	}
}
```

---

## PART C — Python FastAPI Implementation

> **Stack:** FastAPI + SQLAlchemy 2.0 (async) + Pydantic v2 + aiokafka + redis-py + elasticsearch-py

### C.1 Directory Tree

```
app/modules/customer/
├── domain/
│   ├── entities/
│   │   ├── customer.py
│   │   ├── contact.py
│   │   └── kyc.py
│   ├── value_objects/
│   │   ├── customer_type.py
│   │   ├── customer_status.py
│   │   ├── segment.py
│   │   ├── lifecycle.py
│   │   ├── address.py
│   │   ├── email.py
│   │   └── tax_id.py
│   ├── repositories/
│   │   └── customer_repository.py   # Protocol
│   ├── services/
│   │   └── kyc_verifier.py
│   ├── events/
│   │   └── events.py
│   └── errors/
│       └── errors.py
├── application/
│   ├── use_cases/
│   │   ├── create_customer.py
│   │   ├── get_customer.py
│   │   ├── list_customers.py
│   │   ├── suspend_customer.py
│   │   └── ...
│   ├── dtos/
│   │   └── customer_dto.py
│   └── ports/
│       ├── event_producer.py
│       ├── audit_repository.py
│       ├── cache.py
│       └── indexer.py
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.py
│   │   │   ├── customer_repository_impl.py
│   │   │   └── read_model_repository.py
│   │   └── redis/
│   │       └── customer_cache.py
│   ├── messaging/
│   │   └── kafka_producer.py
│   ├── search/
│   │   └── elasticsearch/
│   │       └── customer_indexer.py
│   └── services/
│       └── kyc/
│           └── dbd_provider.py
├── interfaces/
│   ├── http/
│   │   ├── customer_router.py
│   │   ├── schemas.py
│   │   └── dependencies.py
│   └── middleware/
│       └── tenant.py
└── module.py   # Dependency injection container
```

### C.2 Domain — Value Objects (Python)

```python
# app/modules/customer/domain/value_objects/customer_type.py
from enum import Enum


class CustomerType(str, Enum):
    PERSONAL = "PERSONAL"
    BUSINESS = "BUSINESS"
    GOVERNMENT = "GOVERNMENT"
    NGO = "NGO"

    def requires_tax_id(self) -> bool:
        return self in (CustomerType.BUSINESS, CustomerType.GOVERNMENT)

    def can_verify_kyc(self) -> bool:
        return self in (CustomerType.BUSINESS, CustomerType.GOVERNMENT)
```

```python
# app/modules/customer/domain/value_objects/customer_status.py
from enum import Enum


class CustomerStatus(str, Enum):
    PENDING = "PENDING"
    ACTIVE = "ACTIVE"
    SUSPENDED = "SUSPENDED"
    CHURNED = "CHURNED"

    def can_transition_to(self, next_status: "CustomerStatus") -> bool:
        transitions = {
            CustomerStatus.PENDING: {CustomerStatus.ACTIVE},
            CustomerStatus.ACTIVE: {CustomerStatus.SUSPENDED, CustomerStatus.CHURNED},
            CustomerStatus.SUSPENDED: {CustomerStatus.ACTIVE, CustomerStatus.CHURNED},
            CustomerStatus.CHURNED: set(),
        }
        return next_status in transitions.get(self, set())
```

```python
# app/modules/customer/domain/value_objects/address.py
import re
from pydantic import BaseModel, Field, field_validator

POSTCODE_PATTERN = re.compile(r"^\d{5}$")


class Address(BaseModel):
    line1: str = Field(..., max_length=255)
    line2: str = Field("", max_length=255)
    district: str = ""
    amphoe: str = ""
    province: str = ""
    postcode: str = ""
    country: str = "TH"
    lat: float | None = None
    lng: float | None = None

    @field_validator("postcode")
    @classmethod
    def validate_postcode(cls, v: str) -> str:
        if v and not POSTCODE_PATTERN.match(v):
            raise ValueError("invalid postcode (must be 5 digits)")
        return v

    def is_empty(self) -> bool:
        return not (self.line1 or self.province or self.country)
```

```python
# app/modules/customer/domain/value_objects/email.py
import re
from pydantic import BaseModel, field_validator

EMAIL_PATTERN = re.compile(r"^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$")


class Email(str):
    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        if not isinstance(v, str):
            raise TypeError("string required")
        v = v.strip().lower()
        if v and not EMAIL_PATTERN.match(v):
            raise ValueError("invalid email format")
        return cls(v)
```

### C.3 Domain — Entities (Python)

```python
# app/modules/customer/domain/entities/contact.py
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field


class Contact(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    name: str = Field(..., min_length=1, max_length=255)
    phone: str = ""
    email: str = ""
    position: str = ""
    is_primary: bool = False
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(
        cls,
        name: str,
        phone: str = "",
        email: str = "",
        position: str = "",
        is_primary: bool = False,
    ) -> "Contact":
        if not name or len(name) > 255:
            raise ValueError("invalid contact name")
        return cls(
            name=name,
            phone=phone,
            email=email,
            position=position,
            is_primary=is_primary,
        )
```

```python
# app/modules/customer/domain/entities/customer.py
from datetime import datetime, timezone
from uuid import UUID, uuid4
import re

from pydantic import BaseModel, Field

from app.modules.customer.domain.entities.contact import Contact
from app.modules.customer.domain.entities.kyc import KYCInfo
from app.modules.customer.domain.value_objects.address import Address
from app.modules.customer.domain.value_objects.customer_status import CustomerStatus
from app.modules.customer.domain.value_objects.customer_type import CustomerType
from app.modules.customer.domain.value_objects.lifecycle import LifecycleStage
from app.modules.customer.domain.value_objects.segment import Segment
from app.modules.customer.domain.errors.errors import (
    InvalidStatusTransitionError,
    TaxIDRequiredError,
    InvalidCodeError,
    TooManyContactsError,
    ContactNotFoundError,
)

MAX_CONTACTS = 20
CODE_PATTERN = re.compile(r"^[A-Z0-9-]{3,50}$")


class Customer(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    code: str
    name: str
    type: CustomerType
    segment: Segment
    lifecycle: LifecycleStage = LifecycleStage.LEAD
    status: CustomerStatus = CustomerStatus.PENDING
    tax_id: str = ""
    address: Address = Field(default_factory=Address)
    contacts: list[Contact] = Field(default_factory=list)
    package_id: UUID | None = None
    kyc: KYCInfo | None = None
    metadata: dict = Field(default_factory=dict)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    # ---------- Constructor ----------
    @classmethod
    def create(
        cls,
        tenant_id: UUID,
        code: str,
        name: str,
        customer_type: CustomerType,
        segment: Segment,
    ) -> "Customer":
        if not CODE_PATTERN.match(code):
            raise InvalidCodeError(f"invalid code: {code}")
        if not name or len(name) > 255:
            raise InvalidCodeError("invalid name")
        return cls(
            tenant_id=tenant_id,
            code=code,
            name=name,
            type=customer_type,
            segment=segment,
        )

    # ---------- Behavior methods ----------
    def rename(self, new_name: str) -> None:
        if not new_name or len(new_name) > 255:
            raise ValueError("invalid name")
        self.name = new_name
        self._touch()

    def set_tax_id(self, tax_id: str) -> None:
        if self.type.requires_tax_id() and not tax_id:
            raise TaxIDRequiredError()
        if tax_id and not re.match(r"^\d{13}$", tax_id):
            raise ValueError("invalid tax id (must be 13 digits)")
        self.tax_id = tax_id
        self._touch()

    def update_address(self, addr: Address) -> None:
        self.address = addr
        self._touch()

    def activate(self) -> None:
        if not self.status.can_transition_to(CustomerStatus.ACTIVE):
            raise InvalidStatusTransitionError(
                f"cannot activate from {self.status.value}"
            )
        if self.type.requires_tax_id() and not self.tax_id:
            raise TaxIDRequiredError()
        self.status = CustomerStatus.ACTIVE
        self.lifecycle = LifecycleStage.CUSTOMER
        self._touch()

    def suspend(self, reason: str) -> None:
        if not self.status.can_transition_to(CustomerStatus.SUSPENDED):
            raise InvalidStatusTransitionError(
                f"cannot suspend from {self.status.value}"
            )
        if not reason:
            raise ValueError("reason required")
        self.status = CustomerStatus.SUSPENDED
        self.metadata["suspend_reason"] = reason
        self.metadata["suspended_at"] = datetime.now(timezone.utc).isoformat()
        self._touch()

    def churn(self, reason: str) -> None:
        if not self.status.can_transition_to(CustomerStatus.CHURNED):
            raise InvalidStatusTransitionError(
                f"cannot churn from {self.status.value}"
            )
        self.status = CustomerStatus.CHURNED
        self.lifecycle = LifecycleStage.CHURNED
        self.metadata["churn_reason"] = reason
        self.metadata["churned_at"] = datetime.now(timezone.utc).isoformat()
        self._touch()

    def add_contact(self, contact: Contact) -> None:
        if len(self.contacts) >= MAX_CONTACTS:
            raise TooManyContactsError(f"max {MAX_CONTACTS} contacts")
        if contact.is_primary:
            for c in self.contacts:
                c.is_primary = False
        self.contacts.append(contact)
        self._touch()

    def remove_contact(self, contact_id: UUID) -> None:
        for i, c in enumerate(self.contacts):
            if c.id == contact_id:
                self.contacts.pop(i)
                self._touch()
                return
        raise ContactNotFoundError()

    def set_primary_contact(self, contact_id: UUID) -> None:
        found = False
        for c in self.contacts:
            if c.id == contact_id:
                c.is_primary = True
                found = True
            else:
                c.is_primary = False
        if not found:
            raise ContactNotFoundError()
        self._touch()

    def assign_package(self, package_id: UUID) -> None:
        self.package_id = package_id
        self._touch()

    def attach_kyc(self, kyc: KYCInfo) -> None:
        if not self.type.can_verify_kyc():
            raise ValueError("KYC only for BUSINESS/GOVERNMENT")
        self.kyc = kyc
        self._touch()

    # ---------- Query methods ----------
    def is_active(self) -> bool:
        return self.status == CustomerStatus.ACTIVE

    def is_business(self) -> bool:
        return self.type == CustomerType.BUSINESS

    def has_kyc(self) -> bool:
        return self.kyc is not None and self.kyc.is_valid()

    def primary_contact(self) -> Contact | None:
        for c in self.contacts:
            if c.is_primary:
                return c
        return None

    # ---------- Private ----------
    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)

    class Config:
        arbitrary_types_allowed = True
```

### C.4 Domain — Repository Protocol (Python)

```python
# app/modules/customer/domain/repositories/customer_repository.py
from typing import Protocol
from uuid import UUID

from app.modules.customer.domain.entities.customer import Customer
from app.modules.customer.domain.value_objects.customer_status import CustomerStatus


class CustomerRepository(Protocol):
    async def save(self, customer: Customer) -> None: ...
    async def find_by_id(self, customer_id: UUID) -> Customer | None: ...
    async def find_by_code(self, tenant_id: UUID, code: str) -> Customer | None: ...
    async def find_by_tenant(self, tenant_id: UUID) -> list[Customer]: ...
    async def find_by_status(
        self, tenant_id: UUID, status: CustomerStatus
    ) -> list[Customer]: ...
    async def exists_by_code(self, tenant_id: UUID, code: str) -> bool: ...
    async def delete(self, customer_id: UUID) -> None: ...


class CustomerReadModel(Protocol):
    async def get_detail(self, customer_id: UUID) -> "CustomerDetail | None": ...
    async def list_summaries(
        self, filter: "CustomerFilter"
    ) -> tuple[list["CustomerSummary"], int]: ...
    async def search(
        self, tenant_id: UUID, query: str, limit: int
    ) -> list["CustomerSummary"]: ...
```

### C.5 Domain — Errors (Python)

```python
# app/modules/customer/domain/errors/errors.py


class DomainError(Exception):
    """Base for all domain errors."""
    code: str = "DOMAIN_ERROR"

    def __init__(self, message: str = ""):
        self.message = message or self.__class__.__name__
        super().__init__(self.message)


class InvalidCodeError(DomainError):
    code = "INVALID_CODE"


class InvalidNameError(DomainError):
    code = "INVALID_NAME"


class InvalidStatusTransitionError(DomainError):
    code = "INVALID_STATUS_TRANSITION"


class TaxIDRequiredError(DomainError):
    code = "TAX_ID_REQUIRED"


class TooManyContactsError(DomainError):
    code = "TOO_MANY_CONTACTS"


class ContactNotFoundError(DomainError):
    code = "CONTACT_NOT_FOUND"


class CustomerNotFoundError(DomainError):
    code = "CUSTOMER_NOT_FOUND"


class CodeDuplicateError(DomainError):
    code = "CODE_DUPLICATE"


class KYCNotAllowedError(DomainError):
    code = "KYC_NOT_ALLOWED"
```

### C.6 Application — Use Case (Python)

```python
# app/modules/customer/application/use_cases/create_customer.py
from dataclasses import dataclass
from datetime import datetime, timezone
from uuid import UUID

from app.modules.customer.application.ports.event_producer import EventProducer
from app.modules.customer.application.ports.audit_repository import (
    AuditRepository,
    AuditTrail,
)
from app.modules.customer.application.ports.indexer import Indexer
from app.modules.customer.domain.entities.customer import Customer
from app.modules.customer.domain.errors.errors import CodeDuplicateError
from app.modules.customer.domain.events.events import CustomerCreatedEvent
from app.modules.customer.domain.repositories.customer_repository import (
    CustomerRepository,
)
from app.modules.customer.domain.value_objects.address import Address
from app.modules.customer.domain.value_objects.customer_type import CustomerType
from app.modules.customer.domain.value_objects.segment import Segment


@dataclass
class CreateCustomerInput:
    tenant_id: UUID
    code: str
    name: str
    type: CustomerType
    segment: Segment
    tax_id: str
    address: Address
    user_id: UUID
    ip_address: str = ""
    user_agent: str = ""


@dataclass
class CreateCustomerOutput:
    customer_id: UUID
    code: str


class CreateCustomerUseCase:
    def __init__(
        self,
        repo: CustomerRepository,
        audit_repo: AuditRepository,
        producer: EventProducer,
        indexer: Indexer,
    ):
        self._repo = repo
        self._audit = audit_repo
        self._producer = producer
        self._indexer = indexer

    async def execute(self, inp: CreateCustomerInput) -> CreateCustomerOutput:
        # 1. Check duplicate
        if await self._repo.exists_by_code(inp.tenant_id, inp.code):
            raise CodeDuplicateError(f"code {inp.code} already exists")

        # 2. Create aggregate
        customer = Customer.create(
            tenant_id=inp.tenant_id,
            code=inp.code,
            name=inp.name,
            customer_type=inp.type,
            segment=inp.segment,
        )
        customer.set_tax_id(inp.tax_id)
        customer.update_address(inp.address)

        # 3. Persist
        await self._repo.save(customer)

        # 4. Audit
        await self._audit.save(
            AuditTrail(
                user_id=inp.user_id,
                action="CUSTOMER_CREATED",
                entity_type="customer",
                entity_id=customer.id,
                payload={
                    "code": customer.code,
                    "name": customer.name,
                    "type": customer.type.value,
                },
                ip_address=inp.ip_address,
                user_agent=inp.user_agent,
                occurred_at=datetime.now(timezone.utc),
            )
        )

        # 5. Index
        await self._indexer.index(customer)

        # 6. Publish
        await self._producer.publish_customer_created(
            CustomerCreatedEvent(
                customer_id=customer.id,
                tenant_id=customer.tenant_id,
                code=customer.code,
                name=customer.name,
                type=customer.type.value,
                occurred_at=datetime.now(timezone.utc),
            )
        )

        return CreateCustomerOutput(customer_id=customer.id, code=customer.code)
```

### C.7 Infrastructure — SQLAlchemy Models (Python)

```python
# app/modules/customer/infrastructure/persistence/postgres/models.py
from datetime import datetime
from uuid import UUID, uuid4

from sqlalchemy import (
    Boolean,
    DateTime,
    ForeignKey,
    Index,
    String,
    UniqueConstraint,
    func,
)
from sqlalchemy.dialects.postgresql import JSONB, UUID as PGUUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class CustomerModel(Base):
    __tablename__ = "customer_customers"
    __table_args__ = (
        UniqueConstraint("tenant_id", "code", name="uq_customer_tenant_code"),
        Index("idx_customer_tenant", "tenant_id"),
        Index("idx_customer_status", "tenant_id", "status"),
        Index("idx_customer_segment", "tenant_id", "segment"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    code: Mapped[str] = mapped_column(String(50), nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    type: Mapped[str] = mapped_column(String(30), nullable=False)
    segment: Mapped[str] = mapped_column(String(30), nullable=False)
    lifecycle: Mapped[str] = mapped_column(String(30), nullable=False, default="LEAD")
    status: Mapped[str] = mapped_column(String(30), nullable=False, default="PENDING")
    tax_id: Mapped[str] = mapped_column(String(20), default="")
    address: Mapped[dict] = mapped_column(JSONB, default=dict)
    package_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    kyc: Mapped[dict | None] = mapped_column(JSONB, nullable=True)
    metadata_: Mapped[dict] = mapped_column("metadata", JSONB, default=dict)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(
        DateTime, server_default=func.now(), onupdate=func.now()
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)

    contacts: Mapped[list["ContactModel"]] = relationship(
        back_populates="customer", cascade="all, delete-orphan"
    )


class ContactModel(Base):
    __tablename__ = "customer_contacts"
    __table_args__ = (Index("idx_contact_customer", "customer_id"),)

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    customer_id: Mapped[UUID] = mapped_column(
        PGUUID(as_uuid=True), ForeignKey("customer_customers.id", ondelete="CASCADE")
    )
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    phone: Mapped[str] = mapped_column(String(30), default="")
    email: Mapped[str] = mapped_column(String(255), default="")
    position: Mapped[str] = mapped_column(String(100), default="")
    is_primary: Mapped[bool] = mapped_column(Boolean, default=False)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())

    customer: Mapped["CustomerModel"] = relationship(back_populates="contacts")
```

### C.8 Infrastructure — Repository Impl (Python)

```python
# app/modules/customer/infrastructure/persistence/postgres/customer_repository_impl.py
from uuid import UUID

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.modules.customer.domain.entities.contact import Contact
from app.modules.customer.domain.entities.customer import Customer
from app.modules.customer.domain.entities.kyc import KYCInfo
from app.modules.customer.domain.value_objects.address import Address
from app.modules.customer.domain.value_objects.customer_status import CustomerStatus
from app.modules.customer.domain.value_objects.customer_type import CustomerType
from app.modules.customer.domain.value_objects.lifecycle import LifecycleStage
from app.modules.customer.domain.value_objects.segment import Segment
from app.modules.customer.infrastructure.persistence.postgres.models import (
    ContactModel,
    CustomerModel,
)


class PostgresCustomerRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, customer: Customer) -> None:
        model = await self._session.get(CustomerModel, customer.id)
        if model is None:
            model = CustomerModel(id=customer.id)
            self._session.add(model)

        model.tenant_id = customer.tenant_id
        model.code = customer.code
        model.name = customer.name
        model.type = customer.type.value
        model.segment = customer.segment.value
        model.lifecycle = customer.lifecycle.value
        model.status = customer.status.value
        model.tax_id = customer.tax_id
        model.address = customer.address.model_dump()
        model.package_id = customer.package_id
        model.kyc = customer.kyc.model_dump() if customer.kyc else None
        model.metadata_ = customer.metadata

        # replace contacts
        model.contacts.clear()
        for ct in customer.contacts:
            model.contacts.append(
                ContactModel(
                    id=ct.id,
                    customer_id=customer.id,
                    name=ct.name,
                    phone=ct.phone,
                    email=ct.email,
                    position=ct.position,
                    is_primary=ct.is_primary,
                )
            )

        await self._session.flush()

    async def find_by_id(self, customer_id: UUID) -> Customer | None:
        stmt = (
            select(CustomerModel)
            .options(selectinload(CustomerModel.contacts))
            .where(CustomerModel.id == customer_id, CustomerModel.deleted_at.is_(None))
        )
        model = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(model) if model else None

    async def find_by_code(self, tenant_id: UUID, code: str) -> Customer | None:
        stmt = (
            select(CustomerModel)
            .options(selectinload(CustomerModel.contacts))
            .where(
                CustomerModel.tenant_id == tenant_id,
                CustomerModel.code == code,
                CustomerModel.deleted_at.is_(None),
            )
        )
        model = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(model) if model else None

    async def find_by_tenant(self, tenant_id: UUID) -> list[Customer]:
        stmt = (
            select(CustomerModel)
            .options(selectinload(CustomerModel.contacts))
            .where(CustomerModel.tenant_id == tenant_id, CustomerModel.deleted_at.is_(None))
            .order_by(CustomerModel.created_at.desc())
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(r) for r in rows]

    async def exists_by_code(self, tenant_id: UUID, code: str) -> bool:
        stmt = select(func.count()).select_from(CustomerModel).where(
            CustomerModel.tenant_id == tenant_id,
            CustomerModel.code == code,
            CustomerModel.deleted_at.is_(None),
        )
        return (await self._session.execute(stmt)).scalar() > 0

    async def delete(self, customer_id: UUID) -> None:
        from datetime import datetime, timezone
        model = await self._session.get(CustomerModel, customer_id)
        if model:
            model.deleted_at = datetime.now(timezone.utc)

    # ---------- Mappers ----------
    def _to_entity(self, m: CustomerModel) -> Customer:
        return Customer(
            id=m.id,
            tenant_id=m.tenant_id,
            code=m.code,
            name=m.name,
            type=CustomerType(m.type),
            segment=Segment(m.segment),
            lifecycle=LifecycleStage(m.lifecycle),
            status=CustomerStatus(m.status),
            tax_id=m.tax_id or "",
            address=Address(**m.address) if m.address else Address(line1=""),
            contacts=[
                Contact(
                    id=c.id,
                    name=c.name,
                    phone=c.phone,
                    email=c.email,
                    position=c.position,
                    is_primary=c.is_primary,
                    created_at=c.created_at,
                )
                for c in m.contacts
            ],
            package_id=m.package_id,
            kyc=KYCInfo(**m.kyc) if m.kyc else None,
            metadata=m.metadata_ or {},
            created_at=m.created_at,
            updated_at=m.updated_at,
        )
```

### C.9 Infrastructure — Kafka Producer (Python)

```python
# app/modules/customer/infrastructure/messaging/kafka_producer.py
import json
from dataclasses import asdict
from datetime import datetime

from aiokafka import AIOKafkaProducer

from app.modules.customer.domain.events.events import (
    CustomerCreatedEvent,
    CustomerSuspendedEvent,
    PackageAssignedEvent,
)


def _default_encoder(o):
    if isinstance(o, datetime):
        return o.isoformat()
    if hasattr(o, "hex"):  # UUID
        return str(o)
    raise TypeError(f"not serializable: {type(o)}")


class KafkaEventProducer:
    TOPIC_CREATED = "customer.customer.created"
    TOPIC_SUSPENDED = "customer.customer.suspended"
    TOPIC_PACKAGE_ASSIGNED = "customer.package.assigned"

    def __init__(self, producer: AIOKafkaProducer):
        self._producer = producer

    async def _send(self, topic: str, key: str, payload: dict) -> None:
        data = json.dumps(payload, default=_default_encoder).encode("utf-8")
        await self._producer.send_and_wait(topic, key=key.encode("utf-8"), value=data)

    async def publish_customer_created(self, evt: CustomerCreatedEvent) -> None:
        await self._send(self.TOPIC_CREATED, str(evt.customer_id), asdict(evt))

    async def publish_customer_suspended(self, evt: CustomerSuspendedEvent) -> None:
        await self._send(self.TOPIC_SUSPENDED, str(evt.customer_id), asdict(evt))

    async def publish_package_assigned(self, evt: PackageAssignedEvent) -> None:
        await self._send(self.TOPIC_PACKAGE_ASSIGNED, str(evt.customer_id), asdict(evt))
```

### C.10 Interface — FastAPI Router (Python)

```python
# app/modules/customer/interfaces/http/schemas.py
from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, Field


class AddressSchema(BaseModel):
    line1: str = Field(..., max_length=255)
    line2: str = ""
    district: str = ""
    amphoe: str = ""
    province: str = ""
    postcode: str = ""
    country: str = "TH"
    lat: float | None = None
    lng: float | None = None


class CreateCustomerRequest(BaseModel):
    code: str = Field(..., min_length=3, max_length=50, pattern=r"^[A-Z0-9-]+$")
    name: str = Field(..., min_length=1, max_length=255)
    type: str = Field(..., pattern=r"^(PERSONAL|BUSINESS|GOVERNMENT|NGO)$")
    segment: str = Field(..., pattern=r"^(SME|ENTERPRISE|GOVERNMENT|STARTUP|INDIVIDUAL)$")
    tax_id: str = ""
    address: AddressSchema


class CustomerResponse(BaseModel):
    id: UUID
    tenant_id: UUID
    code: str
    name: str
    type: str
    segment: str
    lifecycle: str
    status: str
    tax_id: str
    created_at: datetime
    updated_at: datetime


class SuspendRequest(BaseModel):
    reason: str = Field(..., min_length=1)
```

```python
# app/modules/customer/interfaces/http/customer_router.py
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, status

from app.modules.customer.application.use_cases.create_customer import (
    CreateCustomerInput,
    CreateCustomerUseCase,
)
from app.modules.customer.application.use_cases.suspend_customer import (
    SuspendCustomerInput,
    SuspendCustomerUseCase,
)
from app.modules.customer.domain.errors.errors import (
    CodeDuplicateError,
    CustomerNotFoundError,
    DomainError,
    InvalidStatusTransitionError,
)
from app.modules.customer.domain.value_objects.address import Address
from app.modules.customer.domain.value_objects.customer_type import CustomerType
from app.modules.customer.domain.value_objects.segment import Segment
from app.modules.customer.interfaces.http.dependencies import (
    get_create_customer_uc,
    get_current_user,
    get_suspend_customer_uc,
    get_tenant_id,
)
from app.modules.customer.interfaces.http.schemas import (
    CreateCustomerRequest,
    CustomerResponse,
    SuspendRequest,
)

router = APIRouter(prefix="/api/v1/customers", tags=["customers"])


@router.post("", response_model=CustomerResponse, status_code=status.HTTP_201_CREATED)
async def create_customer(
    req: CreateCustomerRequest,
    uc: CreateCustomerUseCase = Depends(get_create_customer_uc),
    user_id: UUID = Depends(get_current_user),
    tenant_id: UUID = Depends(get_tenant_id),
):
    try:
        out = await uc.execute(
            CreateCustomerInput(
                tenant_id=tenant_id,
                code=req.code,
                name=req.name,
                type=CustomerType(req.type),
                segment=Segment(req.segment),
                tax_id=req.tax_id,
                address=Address(**req.address.model_dump()),
                user_id=user_id,
            )
        )
        return {"id": out.customer_id, "code": out.code, "status": "created"}
    except DomainError as e:
        raise _domain_error_to_http(e)


@router.post("/{customer_id}/suspend", status_code=status.HTTP_200_OK)
async def suspend_customer(
    customer_id: UUID,
    req: SuspendRequest,
    uc: SuspendCustomerUseCase = Depends(get_suspend_customer_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(
            SuspendCustomerInput(
                customer_id=customer_id, reason=req.reason, user_id=user_id
            )
        )
        return {"message": "suspended"}
    except DomainError as e:
        raise _domain_error_to_http(e)


def _domain_error_to_http(e: DomainError) -> HTTPException:
    status_map = {
        CodeDuplicateError: 409,
        CustomerNotFoundError: 404,
        InvalidStatusTransitionError: 409,
    }
    code = status_map.get(type(e), 400)
    return HTTPException(status_code=code, detail={"code": e.code, "message": str(e)})
```

### C.11 Composition Root — `module.py` (Python)

```python
# app/modules/customer/module.py
from dataclasses import dataclass

from aiokafka import AIOKafkaProducer
from elasticsearch import AsyncElasticsearch
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.customer.application.use_cases.create_customer import (
    CreateCustomerUseCase,
)
from app.modules.customer.application.use_cases.suspend_customer import (
    SuspendCustomerUseCase,
)
from app.modules.customer.infrastructure.messaging.kafka_producer import KafkaEventProducer
from app.modules.customer.infrastructure.persistence.postgres.customer_repository_impl import (
    PostgresCustomerRepository,
)
from app.modules.customer.infrastructure.persistence.redis.customer_cache import (
    RedisCustomerCache,
)
from app.modules.customer.infrastructure.search.elasticsearch.customer_indexer import (
    ESCustomerIndexer,
)


@dataclass
class CustomerModuleDeps:
    session: AsyncSession
    redis: Redis
    kafka: AIOKafkaProducer
    es: AsyncElasticsearch
    audit_repo: "AuditRepository"


def build_module(deps: CustomerModuleDeps) -> dict:
    repo = PostgresCustomerRepository(deps.session)
    producer = KafkaEventProducer(deps.kafka)
    cache = RedisCustomerCache(deps.redis)
    indexer = ESCustomerIndexer(deps.es)

    return {
        "create_customer": CreateCustomerUseCase(repo, deps.audit_repo, producer, indexer),
        "suspend_customer": SuspendCustomerUseCase(repo, deps.audit_repo, producer, cache),
    }
```

---

## PART D — Unified Test Suite (ทั้ง 2 ภาษา)

### Go Test
```go
func TestCustomer_StatusTransitions(t *testing.T) { ... }  // ตาม B.13
```

### Python Test
```python
# tests/modules/customer/test_customer_entity.py
import pytest
from uuid import uuid4

from app.modules.customer.domain.entities.customer import Customer
from app.modules.customer.domain.errors.errors import (
    InvalidStatusTransitionError,
    TaxIDRequiredError,
    TooManyContactsError,
)
from app.modules.customer.domain.value_objects.customer_type import CustomerType
from app.modules.customer.domain.value_objects.segment import Segment


def test_activate_personal_ok():
    c = Customer.create(uuid4(), "C001", "Alice", CustomerType.PERSONAL, Segment.INDIVIDUAL)
    c.activate()
    assert c.is_active()


def test_activate_business_without_tax_id_fails():
    c = Customer.create(uuid4(), "C002", "ACME", CustomerType.BUSINESS, Segment.ENTERPRISE)
    with pytest.raises(TaxIDRequiredError):
        c.activate()


def test_pending_to_suspended_fails():
    c = Customer.create(uuid4(), "C003", "Eve", CustomerType.PERSONAL, Segment.INDIVIDUAL)
    with pytest.raises(InvalidStatusTransitionError):
        c.suspend("test")


def test_max_contacts():
    c = Customer.create(uuid4(), "C004", "Test", CustomerType.PERSONAL, Segment.INDIVIDUAL)
    from app.modules.customer.domain.entities.contact import Contact
    for i in range(20):
        c.add_contact(Contact.create(f"Contact{i}"))
    with pytest.raises(TooManyContactsError):
        c.add_contact(Contact.create("Extra"))
```

---

## PART E — AI Generation Prompts

> ใช้ 2 prompt นี้ป้อนให้ AI (ChatGPT/Claude/Gemini) เพื่อ generate code เต็มชุด

### E.1 Prompt สำหรับ Go

```
Generate a complete Go module named "customer" for the project "icmongolang"
following Clean Architecture + DDD. Use the specification provided below.

## Tech Stack
- Go 1.21+
- Gin for HTTP
- GORM for PostgreSQL
- IBM Sarama for Kafka
- go-redis/v8
- go-elasticsearch/v8
- google/uuid

## Project Structure (bắt buộc)
internal/modules/customer/
├── domain/{entity,value_object,repository,service,event,errors}/
├── application/
├── infrastructure/{persistence/postgres,persistence/redis,messaging,search/elasticsearch,services/kyc}/
├── interfaces/{http,middleware}/
└── module.go

## Rules
1. Domain layer MUST NOT import gorm, gin, sarama, redis, es
2. Entity state changes ONLY via behavior methods (no setters)
3. All errors are sentinel errors in domain/errors
4. Repository is an interface in domain, impl in infrastructure
5. Use cases: 1 file per use case, Execute(ctx, input) (output, error)
6. All side effects (audit, kafka, cache, es) are via ports (interfaces)
7. Include unit tests for domain entity
8. Use `icmongolang/pkg/{responses,httpErrors,logger,transaction}` for shared concerns

## Output
1. Full directory tree
2. All .go files (complete, no truncation)
3. migration SQL (YYYYMMDD_customer_init.sql)
4. .env.example additions
5. Wire-up in cmd/api/main.go

## Specification
[PASTE PART A — A.1 to A.10 here]
```

### E.2 Prompt สำหรับ Python FastAPI

```
Generate a complete Python module named "customer" for a FastAPI project
following Clean Architecture + DDD. Use the specification provided below.

## Tech Stack
- Python 3.11+
- FastAPI
- Pydantic v2
- SQLAlchemy 2.0 (async) + asyncpg
- aiokafka
- redis-py (async)
- elasticsearch-py (async)
- pytest + pytest-asyncio

## Project Structure (bắt buộc)
app/modules/customer/
├── domain/{entities,value_objects,repositories,services,events,errors}/
├── application/{use_cases,dtos,ports}/
├── infrastructure/{persistence/postgres,persistence/redis,messaging,search/elasticsearch,services/kyc}/
├── interfaces/{http,middleware}/
└── module.py

## Rules
1. Domain layer MUST NOT import sqlalchemy, fastapi, aiokafka
2. Entities: use Pydantic BaseModel with behavior methods
3. Repositories: use typing.Protocol for interfaces
4. Errors: custom exception classes with `code` attribute
5. Use cases: async class with `execute(input) -> output`, 1 file per use case
6. All side effects via Protocol-based ports
7. Include pytest tests for domain entity
8. Use dependency injection via FastAPI Depends

## Output
1. Full directory tree
2. All .py files (complete, no truncation)
3. Alembic migration
4. .env.example additions
5. App factory wire-up (main.py)

## Specification
[PASTE PART A — A.1 to A.10 here]
```

---

# 📋 MODULE 02-08 — Compact Pattern

> ใช้ **โครงเดียวกับ Module 01** ทุกประการ — AI สามารถ generate ได้จาก spec สั้นๆ

## Module 02: `package`

**Business:** Package + Subscription + Quota + Feature gating

**Aggregates:**
- `Package` — code, name, tier, version, price, billing_cycle, quotas, features
- `Subscription` — customer_id, package_id, status, started_at, expires_at, auto_renew, usage

**Key Use Cases:** `CreatePackage`, `Subscribe`, `UpgradeSubscription`, `CancelSubscription`, `RenewSubscription`, `CheckQuota`, `AutoExpireSubscriptions`

**Events:** `package.package.created`, `package.subscription.created|upgraded|cancelled|expired`, `package.quota.exceeded`

**Tables:** `package_packages`, `package_subscriptions`

**Migration:** `20260102_package_init.sql`

---

## Module 03: `device`

**Business:** Site + Zone + Device + Telemetry + Command + Automation Rule

**Aggregates:**
- `Site` (root) → zones
- `Device` (root) — serial, type, protocol, status, shadow
- `AutomationRule` (root) — conditions, actions, cooldown
- `Telemetry` (InfluxDB)
- `Command` (root)

**Key Use Cases:** `CreateSite`, `RegisterDevice`, `IngestTelemetry`, `IssueCommand`, `AckCommand`, `CreateAutomationRule`, `EvaluateAutomation`, `UpdateFirmware`, `DecommissionDevice`

**Events:** `device.device.registered|online|offline`, `device.telemetry.ingested`, `device.command.issued|acked`, `device.rule.fired`

**Tables:** `device_sites`, `device_zones`, `device_devices`, `device_commands`, `device_automation_rules`

**InfluxDB:** measurement `telemetry`, tags `device_id,site_id,metric`, field `value`

**MQTT Topics:** `iot/{serial}/telemetry`, `/heartbeat`, `/status`, `/ack`, `/command`, `/config`, `/firmware`

**Migration:** `20260103_device_init.sql`

---

## Module 04: `erp`

**Aggregates:** `Order`, `Invoice`, `StockItem`, `Warehouse`

**Key Use Cases:** `CreateOrder`, `ConfirmOrder`, `ShipOrder`, `DeliverOrder`, `CreateInvoice`, `IssueInvoice`, `RecordPayment`, `RestockItem`

**Events:** `erp.order.created|confirmed|shipped|delivered|cancelled`, `erp.invoice.issued|paid`, `erp.stock.low`

**Tables:** `erp_orders`, `erp_order_items`, `erp_invoices`, `erp_warehouses`, `erp_stock_items`

**Migration:** `20260104_erp_init.sql`

---

## Module 05: `crm`

**Aggregates:** `Lead`, `Opportunity`, `Ticket`, `Activity`, `Campaign`

**Key Use Cases:** `CaptureLead`, `QualifyLead`, `ConvertLead`, `OpenTicket`, `AssignTicket`, `ResolveTicket`, `EscalateTicket`

**Events:** `crm.lead.created|converted`, `crm.ticket.opened|resolved`

**Tables:** `crm_leads`, `crm_opportunities`, `crm_tickets`, `crm_activities`, `crm_campaigns`

**Migration:** `20260105_crm_init.sql`

---

## Module 06: `logistics`

**Aggregates:** `Shipment`, `Route`, `Vehicle`, `Driver`

**Key Use Cases:** `CreateShipment`, `AssignVehicle`, `UpdateLocation`, `RecordTemperature`, `DeliverShipment`

**Events:** `logistics.shipment.created|delivered`, `logistics.temperature.breached`

**Tables:** `logistics_shipments`, `logistics_routes`, `logistics_vehicles`, `logistics_drivers`

**Migration:** `20260106_logistics_init.sql`

---

## Module 07: `report` (Upgrade)

**Business:** Refactor ของเดิม + เพิ่ม IoT analytics

**Aggregates:** `ReportDefinition`, `ReportRun`, `Dashboard`, `Widget`

**Key Use Cases:** `CreateReport`, `RunReport`, `ScheduleReport`, `ExportReport`, `GetDashboard`, `GenerateInsight`

**Data Sources:** Postgres, InfluxDB, Elasticsearch

**Tables:** `report_definitions`, `report_runs`, `report_dashboards`

**Migration:** `20260107_report_upgrade.sql`

**⚠️ Breaking Changes:** alias old endpoints + 2-release deprecation

---

## Module 08: `auth` (Upgrade)

**Business:** เพิ่ม tenant-aware JWT + SSO + MFA + session management

**Key Use Cases:** `Login`, `LoginWithSSO`, `EnableMFA`, `VerifyMFA`, `RefreshToken`, `RevokeAllSessions`

**Breaking:** JWT payload เพิ่ม `tenant_id`, `session_id`, `mfa_verified`

**Migration:** `20260108_auth_upgrade.sql` (add columns only, no drops)

---

# 🗺️ Roadmap — Module Production Order

| Phase | Modules | เดือน | เหตุผล |
| :--- | :--- | :-: | :--- |
| **P0.1** | `customer`, `package` | 1-2 | รากฐาน: ลูกค้า + subscription |
| **P0.2** | `device`, `mqtt`, `influxdb` | 3-4 | IoT core |
| **P0.3** | `alarm`, `notifier`, `websocket` | 5 | Realtime + alert |
| **P0.4** | `payment`, `auth` (upgrade), `users` (upgrade) | 6 | Billing + security |
| **P1.1** | `erp`, `items`, `purchaseorder` | 7-9 | Business core |
| **P1.2** | `crm`, `quotation` | 10-11 | Sales |
| **P1.3** | `report` (upgrade), `dashboard` (upgrade) | 12 | Analytics |
| **P2.1** | `logistics`, `wos` | 13-15 | Supply chain |
| **P2.2** | `settings` (upgrade), `auditlog` (upgrade), `pdpa` (upgrade) | 16-18 | Governance |
| **P3** | ที่เหลือ + ของเดิม upgrade ทยอย | 19-24 | ครบระบบ |

---

# ✅ AI-Ready Checklist (ต่อ module)

ก่อนป้อนให้ AI generate ให้ตรวจ:

- [ ] **Business spec** ครบ (Purpose, Scope, Stakeholders)
- [ ] **Domain model** ครบ (Aggregate, Entity, VO, Event, Error)
- [ ] **Invariants** ระบุชัด
- [ ] **Use cases** ระบุ input/output/side effects
- [ ] **API contract** ระบุ endpoint + role
- [ ] **DB schema** มี DDL
- [ ] **Kafka topics** ระบุ publish/consume
- [ ] **Env vars** ระบุ
- [ ] **Migration filename** ระบุ
- [ ] **Breaking changes** ระบุ + expand-contract plan

---

## สรุป

| สิ่งที่ได้ | จำนวน |
| :--- | :-: |
| **โมดูลสมบูรณ์** | 1 (`customer`) |
| **ภาษา** | 2 (Go + Python FastAPI) |
| **Layer** | 4 (domain, application, infra, interface) |
| **Use cases** | 13 |
| **API endpoints** | 15 |
| **Tables** | 3 |
| **Kafka topics** | 7 |
| **Test cases** | 5+ |
| **AI prompts** | 2 (Go + Python) |
| **Modules ที่มี compact spec** | 7 (package, device, erp, crm, logistics, report, auth) |
| **Modules ที่ต้องทำต่อ** | 32 |

**ขั้นตอนถัดไป:**
1. ป้อน **Module 01 (customer)** + **Prompt E.1** ให้ AI → ได้ Go code
2. ป้อน **Module 01** + **Prompt E.2** ให้ AI → ได้ Python code
3. ทำ Module 02-08 ตาม pattern เดียวกัน
4. Modules ที่เหลือ — ใช้โครงเดียวกัน ระบุ spec ในเอกสารแยก
 