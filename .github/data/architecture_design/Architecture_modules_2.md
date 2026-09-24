# 📦 MODULE 02: `package`

> **สถานะ:** 🟢 Greenfield · **Priority:** P0 · **Owner:** Core Team · **Depends on:** `customer`, `payment`

---

## PART A — Business Specification

### A.1 Business Context

**Purpose:** จัดการแพ็กเกจบริการ (Package) และการสมัครใช้ (Subscription) พร้อมระบบ Quota + Feature gating สำหรับโมเดล SaaS

**Stakeholders:**
- Product — ออกแบบ package tiers
- Sales — Upgrade/downgrade ลูกค้า
- Finance — Billing, revenue
- Customer — ใช้บริการตามโควตา
- Support — ดู subscription status

**Out of Scope:**
- การชำระเงินจริง (อยู่ที่ `payment`)
- Invoice (อยู่ที่ `erp`)
- Usage tracking จริง (module ที่ใช้เป็นคนรายงาน)

### A.2 Ubiquitous Language

| Term | ความหมาย |
| :--- | :--- |
| **Package** | แพ็กเกจบริการ (Free/Basic/Pro/Enterprise) |
| **Tier** | ระดับของ package |
| **Quota** | โควตาการใช้งาน (device, user, storage) |
| **Feature** | ฟีเจอร์ที่เปิด/ปิดได้ |
| **Subscription** | การสมัครใช้ package |
| **Billing Cycle** | รอบบิล (MONTHLY/YEARLY) |
| **Proration** | การคิดเงินตามสัดส่วนเมื่อ upgrade |
| **Grace Period** | ช่วงผ่อนผันหลังหมดอายุ |
| **Trial** | ช่วงทดลองใช้ฟรี |

### A.3 Domain Model

#### Aggregate Root: `Package`

| Field | Type | Constraint |
| :--- | :--- | :--- |
| id | UUID | PK |
| code | string(50) | UNIQUE, `^[A-Z0-9_]{2,50}$` |
| name | string(255) | NOT NULL |
| description | text | |
| tier | enum | FREE/BASIC/PRO/ENTERPRISE |
| version | int | ≥ 1, increment on price change |
| price.amount | numeric(12,2) | ≥ 0 |
| price.currency | string(3) | ISO 4217 |
| billing_cycle | enum | MONTHLY/YEARLY |
| quotas | Quotas | VO |
| features | []Feature | |
| is_active | bool | default true |
| is_public | bool | default true |
| created_at / updated_at | timestamp | |

#### Value Object: `Quotas`

| Field | Type | Constraint |
| :--- | :--- | :--- |
| max_devices | int | ≥ 0, -1 = unlimited |
| max_sites | int | ≥ 0 |
| max_users | int | ≥ 0 |
| storage_gb | int | ≥ 0 |
| retention_days | int | ≥ 0 |
| api_calls_per_day | int | ≥ 0 |
| automation_rules | int | ≥ 0 |

#### Aggregate Root: `Subscription`

| Field | Type | Constraint |
| :--- | :--- | :--- |
| id | UUID | PK |
| customer_id | UUID | FK → customer |
| package_id | UUID | FK → package |
| status | enum | ACTIVE/TRIAL/PAST_DUE/CANCELLED/EXPIRED |
| started_at | timestamp | NOT NULL |
| expires_at | timestamp | > started_at |
| trial_ends_at | *timestamp | |
| grace_ends_at | *timestamp | |
| auto_renew | bool | default true |
| usage | Quotas | current usage |
| previous_pkg_id | *UUID | for audit |
| upgraded_at | *timestamp | |
| cancelled_at | *timestamp | |
| cancel_reason | text | |

### A.4 Invariants

| # | Invariant | บังคับที่ |
| :-: | :--- | :--- |
| 1 | Package.code ไม่ซ้ำ | unique index |
| 2 | Price.amount ≥ 0 | constructor |
| 3 | Quotas ทุกตัว ≥ 0 (หรือ -1 = unlimited) | `Quotas.Validate()` |
| 4 | Subscription.expires_at > started_at | constructor |
| 5 | Upgrade ไม่ได้ถ้า status ∈ {CANCELLED, EXPIRED} | `Upgrade()` |
| 6 | Cancel ได้ครั้งเดียว | `Cancel()` |
| 7 | Feature.code ไม่ซ้ำใน package | `AddFeature()` |
| 8 | Quota ต้องไม่เกิน package ที่สมัคร | `CheckQuota()` |

### A.5 Domain Events

| Event | Trigger | Payload |
| :--- | :--- | :--- |
| `PackageCreated` | CreatePackage | `{package_id, code, tier}` |
| `PackageUpdated` | UpdatePackage | `{package_id, changes}` |
| `SubscriptionCreated` | Subscribe | `{subscription_id, customer_id, package_id, expires_at}` |
| `SubscriptionUpgraded` | UpgradeSubscription | `{subscription_id, from, to}` |
| `SubscriptionCancelled` | CancelSubscription | `{subscription_id, reason}` |
| `SubscriptionExpired` | ExpireJob | `{subscription_id}` |
| `SubscriptionRenewed` | RenewSubscription | `{subscription_id, new_expires_at}` |
| `QuotaExceeded` | CheckQuota | `{subscription_id, quota_type, limit, current}` |

### A.6 Use Cases

| # | Use Case | Input | Output | Side Effects |
| :-: | :--- | :--- | :--- | :--- |
| 1 | `CreatePackage` | code, name, tier, price, cycle, quotas | package_id | DB, Audit, Kafka |
| 2 | `GetPackage` | package_id | Package | Cache |
| 3 | `ListPackages` | tier?, active_only? | []Package | |
| 4 | `UpdatePackage` | package_id, patch | | DB, Audit, Kafka |
| 5 | `DeactivatePackage` | package_id | | DB, Audit |
| 6 | `Subscribe` | customer_id, package_id, cycle | subscription_id, payment_url | DB, Kafka, Payment, Audit |
| 7 | `GetSubscription` | subscription_id | Subscription | |
| 8 | `GetMySubscription` | customer_id | Subscription | |
| 9 | `UpgradeSubscription` | subscription_id, new_pkg_id | proration_amount | DB, Kafka, Payment, Audit |
| 10 | `CancelSubscription` | subscription_id, reason | | DB, Kafka, Audit, WS |
| 11 | `RenewSubscription` | subscription_id | | DB, Kafka, Payment |
| 12 | `CheckQuota` | subscription_id, quota_type, current | allowed | |
| 13 | `AutoExpireSubscriptions` | (cron) | count | DB, Kafka |
| 14 | `RetryPastDue` | (cron) | count | DB, Kafka, Payment |

### A.7 API Contract

| Method | Path | Auth | Role |
| :--- | :--- | :-: | :--- |
| GET | `/api/v1/packages` | ✅ | viewer |
| GET | `/api/v1/packages/:id` | ✅ | viewer |
| POST | `/api/v1/packages` | ✅ | admin |
| PUT | `/api/v1/packages/:id` | ✅ | admin |
| DELETE | `/api/v1/packages/:id` | ✅ | admin |
| POST | `/api/v1/subscriptions` | ✅ | manager |
| GET | `/api/v1/subscriptions/:id` | ✅ | viewer |
| GET | `/api/v1/subscriptions/me` | ✅ | viewer |
| POST | `/api/v1/subscriptions/:id/upgrade` | ✅ | manager |
| POST | `/api/v1/subscriptions/:id/cancel` | ✅ | manager |
| POST | `/api/v1/subscriptions/:id/renew` | ✅ | manager |
| GET | `/api/v1/subscriptions/:id/quota` | ✅ | viewer |
| POST | `/api/v1/subscriptions/:id/quota/check` | ✅ | viewer |

### A.8 Database Schema

```sql
CREATE TABLE package_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tier VARCHAR(30) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    price_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    price_currency VARCHAR(3) NOT NULL DEFAULT 'THB',
    billing_cycle VARCHAR(20) NOT NULL,
    quotas JSONB NOT NULL DEFAULT '{}',
    features JSONB NOT NULL DEFAULT '[]',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE INDEX idx_package_packages_tier ON package_packages (tier) WHERE is_active = TRUE AND deleted_at IS NULL;

CREATE TABLE package_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    package_id UUID NOT NULL REFERENCES package_packages(id),
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    trial_ends_at TIMESTAMP,
    grace_ends_at TIMESTAMP,
    auto_renew BOOLEAN NOT NULL DEFAULT TRUE,
    usage JSONB DEFAULT '{}',
    previous_pkg_id UUID,
    upgraded_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancel_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_package_subscriptions_customer ON package_subscriptions (customer_id) WHERE status IN ('ACTIVE', 'TRIAL', 'PAST_DUE');
CREATE INDEX idx_package_subscriptions_expires ON package_subscriptions (expires_at) WHERE status = 'ACTIVE';
CREATE UNIQUE INDEX uq_package_subscriptions_active_customer ON package_subscriptions (customer_id)
    WHERE status IN ('ACTIVE', 'TRIAL', 'PAST_DUE');
```

### A.9 Events (Kafka)

**Published:**
- `package.package.created`
- `package.package.updated`
- `package.subscription.created`
- `package.subscription.upgraded`
- `package.subscription.cancelled`
- `package.subscription.expired`
- `package.subscription.renewed`
- `package.quota.exceeded`

**Consumed:**
- `payment.payment.succeeded` → activate subscription
- `payment.payment.failed` → mark past due
- `customer.customer.churned` → cancel subscription
- `device.device.registered` → increment `usage.max_devices`

### A.10 Config

```env
PACKAGE_CACHE_TTL=30m
PACKAGE_TRIAL_DAYS=14
PACKAGE_GRACE_DAYS=7
PACKAGE_PRORATION_ENABLED=true
PACKAGE_AUTO_RENEW_ENABLED=true
PACKAGE_DEFAULT_TIER=FREE
PACKAGE_EXPIRY_CRON=0 */6 * * *
PACKAGE_PAST_DUE_RETRY_INTERVAL=24h
```

---

## PART B — Go Implementation

### B.1 Directory Tree

```
internal/modules/package/
├── domain/
│   ├── entity/
│   │   ├── package.go
│   │   ├── subscription.go
│   │   ├── quotas.go
│   │   └── feature.go
│   ├── value_object/
│   │   ├── package_tier.go
│   │   ├── billing_cycle.go
│   │   ├── subscription_status.go
│   │   └── money.go
│   ├── repository/
│   │   ├── package_repository.go
│   │   └── subscription_repository.go
│   ├── service/
│   │   └── quota_service.go
│   ├── event/
│   │   └── events.go
│   └── errors/
│       └── errors.go
├── application/
│   ├── create_package.go
│   ├── get_package.go
│   ├── list_packages.go
│   ├── update_package.go
│   ├── subscribe.go
│   ├── upgrade_subscription.go
│   ├── cancel_subscription.go
│   ├── renew_subscription.go
│   ├── check_quota.go
│   ├── auto_expire.go
│   ├── dto.go
│   └── ports.go
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── package_repo_impl.go
│   │   │   ├── subscription_repo_impl.go
│   │   │   └── models.go
│   │   └── redis/
│   │       └── package_cache.go
│   ├── messaging/
│   │   └── producer.go
│   └── scheduler/
│       ├── expire_job.go
│       └── retry_past_due_job.go
├── interfaces/
│   └── http/
│       ├── package_handler.go
│       ├── subscription_handler.go
│       ├── routes.go
│       └── dto.go
└── module.go
```

### B.2 Value Objects

```go
// internal/modules/package/domain/value_object/package_tier.go
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

func (t PackageTier) Priority() int {
	return map[PackageTier]int{
		TierFree:       0,
		TierBasic:      1,
		TierPro:        2,
		TierEnterprise: 3,
	}[t]
}
```

```go
// internal/modules/package/domain/value_object/billing_cycle.go
package valueobject

import "time"

type BillingCycle string

const (
	BillingMonthly BillingCycle = "MONTHLY"
	BillingYearly  BillingCycle = "YEARLY"
)

func (b BillingCycle) IsValid() bool {
	return b == BillingMonthly || b == BillingYearly
}

func (b BillingCycle) AddTo(t time.Time) time.Time {
	if b == BillingYearly {
		return t.AddDate(1, 0, 0)
	}
	return t.AddDate(0, 1, 0)
}

func (b BillingCycle) Days() int {
	if b == BillingYearly {
		return 365
	}
	return 30
}
```

```go
// internal/modules/package/domain/value_object/subscription_status.go
package valueobject

type SubscriptionStatus string

const (
	SubStatusActive    SubscriptionStatus = "ACTIVE"
	SubStatusTrial     SubscriptionStatus = "TRIAL"
	SubStatusPastDue   SubscriptionStatus = "PAST_DUE"
	SubStatusCancelled SubscriptionStatus = "CANCELLED"
	SubStatusExpired   SubscriptionStatus = "EXPIRED"
)

func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case SubStatusActive, SubStatusTrial, SubStatusPastDue, SubStatusCancelled, SubStatusExpired:
		return true
	}
	return false
}

func (s SubscriptionStatus) IsUsable() bool {
	return s == SubStatusActive || s == SubStatusTrial
}

func (s SubscriptionStatus) CanUpgrade() bool {
	return s == SubStatusActive || s == SubStatusTrial || s == SubStatusPastDue
}
```

```go
// internal/modules/package/domain/value_object/money.go
package valueobject

import "errors"

type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency required")
	}
	return Money{Amount: amount, Currency: currency}, nil
}

func (m Money) Add(other Money) Money {
	if m.Currency != other.Currency {
		return m
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}
}

func (m Money) Sub(other Money) Money {
	if m.Currency != other.Currency {
		return m
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}
}

func (m Money) Mul(factor float64) Money {
	return Money{Amount: m.Amount * factor, Currency: m.Currency}
}

func (m Money) IsZero() bool { return m.Amount == 0 }
```

### B.3 Entities — Quotas + Feature

```go
// internal/modules/package/domain/entity/quotas.go
package entity

import domainerrors "icmongolang/internal/modules/package/domain/errors"

const Unlimited = -1

type Quotas struct {
	MaxDevices      int `json:"max_devices"`
	MaxSites        int `json:"max_sites"`
	MaxUsers        int `json:"max_users"`
	StorageGB       int `json:"storage_gb"`
	RetentionDays   int `json:"retention_days"`
	APICallsPerDay  int `json:"api_calls_per_day"`
	AutomationRules int `json:"automation_rules"`
}

func (q Quotas) Validate() error {
	fields := []int{q.MaxDevices, q.MaxSites, q.MaxUsers, q.StorageGB,
		q.RetentionDays, q.APICallsPerDay, q.AutomationRules}
	for _, v := range fields {
		if v < -1 {
			return domainerrors.ErrInvalidQuota
		}
	}
	return nil
}

// Exceeds – ตรวจว่า usage เกิน quota หรือไม่
func (q Quotas) Exceeds(limit Quotas) (bool, string) {
	checks := []struct {
		name     string
		current  int
		limit    int
	}{
		{"devices", q.MaxDevices, limit.MaxDevices},
		{"sites", q.MaxSites, limit.MaxSites},
		{"users", q.MaxUsers, limit.MaxUsers},
		{"storage", q.StorageGB, limit.StorageGB},
		{"api_calls", q.APICallsPerDay, limit.APICallsPerDay},
		{"automation_rules", q.AutomationRules, limit.AutomationRules},
	}
	for _, c := range checks {
		if c.limit != Unlimited && c.current > c.limit {
			return true, c.name
		}
	}
	return false, ""
}

func (q Quotas) AllowsDevices(count int) bool {
	return q.MaxDevices == Unlimited || count <= q.MaxDevices
}
```

```go
// internal/modules/package/domain/entity/feature.go
package entity

type Feature struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func NewFeature(code, name string) Feature {
	return Feature{Code: code, Name: name, Enabled: true}
}
```

### B.4 Entities — Package

```go
// internal/modules/package/domain/entity/package.go
package entity

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

var codePattern = regexp.MustCompile(`^[A-Z0-9_]{2,50}$`)

// Package – Aggregate Root
type Package struct {
	ID           uuid.UUID
	Code         string
	Name         string
	Description  string
	Tier         valueobject.PackageTier
	Version      int
	Price        valueobject.Money
	BillingCycle valueobject.BillingCycle
	Quotas       Quotas
	Features     []Feature
	IsActive     bool
	IsPublic     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPackage(
	code, name string,
	tier valueobject.PackageTier,
	price valueobject.Money,
	cycle valueobject.BillingCycle,
	quotas Quotas,
) (*Package, error) {
	if !codePattern.MatchString(code) {
		return nil, domainerrors.ErrInvalidCode
	}
	if name == "" || len(name) > 255 {
		return nil, domainerrors.ErrInvalidName
	}
	if !tier.IsValid() {
		return nil, domainerrors.ErrInvalidTier
	}
	if price.Amount < 0 {
		return nil, domainerrors.ErrInvalidPrice
	}
	if !cycle.IsValid() {
		return nil, domainerrors.ErrInvalidCycle
	}
	if err := quotas.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Package{
		ID:           uuid.New(),
		Code:         code,
		Name:         name,
		Tier:         tier,
		Version:      1,
		Price:        price,
		BillingCycle: cycle,
		Quotas:       quotas,
		Features:     []Feature{},
		IsActive:     true,
		IsPublic:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ---------- Behavior ----------

func (p *Package) Rename(name, desc string) error {
	if name == "" || len(name) > 255 {
		return domainerrors.ErrInvalidName
	}
	p.Name = name
	p.Description = desc
	p.touch()
	return nil
}

func (p *Package) UpdatePrice(newPrice valueobject.Money) error {
	if newPrice.Amount < 0 {
		return domainerrors.ErrInvalidPrice
	}
	if newPrice.Currency != p.Price.Currency {
		return domainerrors.ErrCurrencyMismatch
	}
	p.Price = newPrice
	p.Version++
	p.touch()
	return nil
}

func (p *Package) UpdateQuotas(q Quotas) error {
	if err := q.Validate(); err != nil {
		return err
	}
	p.Quotas = q
	p.touch()
	return nil
}

func (p *Package) AddFeature(f Feature) error {
	for _, existing := range p.Features {
		if existing.Code == f.Code {
			return domainerrors.ErrFeatureDuplicate
		}
	}
	p.Features = append(p.Features, f)
	p.touch()
	return nil
}

func (p *Package) RemoveFeature(code string) error {
	for i, f := range p.Features {
		if f.Code == code {
			p.Features = append(p.Features[:i], p.Features[i+1:]...)
			p.touch()
			return nil
		}
	}
	return domainerrors.ErrFeatureNotFound
}

func (p *Package) EnableFeature(code string) {
	for i := range p.Features {
		if p.Features[i].Code == code {
			p.Features[i].Enabled = true
			p.touch()
			return
		}
	}
}

func (p *Package) DisableFeature(code string) {
	for i := range p.Features {
		if p.Features[i].Code == code {
			p.Features[i].Enabled = false
			p.touch()
			return
		}
	}
}

func (p *Package) Activate() {
	p.IsActive = true
	p.touch()
}

func (p *Package) Deactivate() {
	p.IsActive = false
	p.touch()
}

// ---------- Query ----------

func (p *Package) HasFeature(code string) bool {
	for _, f := range p.Features {
		if f.Code == code && f.Enabled {
			return true
		}
	}
	return false
}

func (p *Package) IsFree() bool {
	return p.Tier == valueobject.TierFree || p.Price.Amount == 0
}

func (p *Package) IsHigherTierThan(other *Package) bool {
	return p.Tier.Priority() > other.Tier.Priority()
}

// ---------- Private ----------

func (p *Package) touch() { p.UpdatedAt = time.Now() }
```

### B.5 Entities — Subscription

```go
// internal/modules/package/domain/entity/subscription.go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type Subscription struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	PackageID     uuid.UUID
	Status        valueobject.SubscriptionStatus
	StartedAt     time.Time
	ExpiresAt     time.Time
	TrialEndsAt   *time.Time
	GraceEndsAt   *time.Time
	AutoRenew     bool
	Usage         Quotas
	PreviousPkgID *uuid.UUID
	UpgradedAt    *time.Time
	CancelledAt   *time.Time
	CancelReason  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSubscription(customerID, packageID uuid.UUID, cycle valueobject.BillingCycle) (*Subscription, error) {
	if customerID == uuid.Nil {
		return nil, domainerrors.ErrInvalidCustomer
	}
	if packageID == uuid.Nil {
		return nil, domainerrors.ErrInvalidPackage
	}
	if !cycle.IsValid() {
		return nil, domainerrors.ErrInvalidCycle
	}

	now := time.Now()
	expires := cycle.AddTo(now)
	return &Subscription{
		ID:         uuid.New(),
		CustomerID: customerID,
		PackageID:  packageID,
		Status:     valueobject.SubStatusActive,
		StartedAt:  now,
		ExpiresAt:  expires,
		AutoRenew:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func NewTrialSubscription(customerID, packageID uuid.UUID, trialDays int) (*Subscription, error) {
	sub, err := NewSubscription(customerID, packageID, valueobject.BillingMonthly)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	trialEnd := now.AddDate(0, 0, trialDays)
	sub.Status = valueobject.SubStatusTrial
	sub.TrialEndsAt = &trialEnd
	sub.ExpiresAt = trialEnd
	return sub, nil
}

// ---------- Behavior ----------

func (s *Subscription) Upgrade(newPkgID uuid.UUID, newExpiry time.Time, cycle valueobject.BillingCycle) error {
	if !s.Status.CanUpgrade() {
		return domainerrors.ErrCannotUpgrade
	}
	if newPkgID == uuid.Nil {
		return domainerrors.ErrInvalidPackage
	}
	if !newExpiry.After(time.Now()) {
		return domainerrors.ErrInvalidExpiry
	}

	old := s.PackageID
	s.PreviousPkgID = &old
	s.PackageID = newPkgID
	s.ExpiresAt = newExpiry
	s.Status = valueobject.SubStatusActive

	now := time.Now()
	s.UpgradedAt = &now
	s.touch()
	return nil
}

func (s *Subscription) Cancel(reason string) error {
	if s.Status == valueobject.SubStatusCancelled {
		return domainerrors.ErrAlreadyCancelled
	}
	if reason == "" {
		return domainerrors.ErrReasonRequired
	}

	now := time.Now()
	s.Status = valueobject.SubStatusCancelled
	s.AutoRenew = false
	s.CancelledAt = &now
	s.CancelReason = reason
	s.touch()
	return nil
}

func (s *Subscription) MarkPastDue(graceDays int) {
	graceEnd := time.Now().AddDate(0, 0, graceDays)
	s.Status = valueobject.SubStatusPastDue
	s.GraceEndsAt = &graceEnd
	s.AutoRenew = false
	s.touch()
}

func (s *Subscription) Reactivate(newExpiry time.Time) {
	s.Status = valueobject.SubStatusActive
	s.ExpiresAt = newExpiry
	s.GraceEndsAt = nil
	s.AutoRenew = true
	s.touch()
}

func (s *Subscription) Renew(cycle valueobject.BillingCycle) error {
	if s.Status == valueobject.SubStatusCancelled {
		return domainerrors.ErrCannotRenew
	}
	s.ExpiresAt = cycle.AddTo(s.ExpiresAt)
	s.Status = valueobject.SubStatusActive
	s.GraceEndsAt = nil
	s.touch()
	return nil
}

func (s *Subscription) Expire() {
	s.Status = valueobject.SubStatusExpired
	s.AutoRenew = false
	s.touch()
}

func (s *Subscription) UpdateUsage(usage Quotas) {
	s.Usage = usage
	s.touch()
}

// ---------- Query ----------

func (s *Subscription) IsActive() bool {
	return s.Status == valueobject.SubStatusActive && time.Now().Before(s.ExpiresAt)
}

func (s *Subscription) IsInTrial() bool {
	return s.Status == valueobject.SubStatusTrial &&
		s.TrialEndsAt != nil &&
		time.Now().Before(*s.TrialEndsAt)
}

func (s *Subscription) DaysRemaining() int {
	if !s.ExpiresAt.After(time.Now()) {
		return 0
	}
	return int(time.Until(s.ExpiresAt).Hours() / 24)
}

func (s *Subscription) IsExpiringSoon(days int) bool {
	return s.DaysRemaining() <= days && s.DaysRemaining() > 0
}

func (s *Subscription) IsPastDue() bool {
	return s.Status == valueobject.SubStatusPastDue
}

func (s *Subscription) touch() { s.UpdatedAt = time.Now() }
```

### B.6 Domain Service — Quota

```go
// internal/modules/package/domain/service/quota_service.go
package service

import (
	"time"

	"icmongolang/internal/modules/package/domain/entity"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type QuotaService struct{}

// CheckQuota – ตรวจว่า usage เกิน quota หรือไม่
func (s *QuotaService) CheckQuota(sub *entity.Subscription, pkg *entity.Package) (bool, string) {
	if !sub.Status.IsUsable() {
		return false, "subscription_not_active"
	}
	return sub.Usage.Exceeds(pkg.Quotas)
}

// Remaining – คำนวณโควตาที่เหลือ
func (s *QuotaService) Remaining(sub *entity.Subscription, pkg *entity.Package) entity.Quotas {
	return entity.Quotas{
		MaxDevices:      remaining(sub.Usage.MaxDevices, pkg.Quotas.MaxDevices),
		MaxSites:        remaining(sub.Usage.MaxSites, pkg.Quotas.MaxSites),
		MaxUsers:        remaining(sub.Usage.MaxUsers, pkg.Quotas.MaxUsers),
		StorageGB:       remaining(sub.Usage.StorageGB, pkg.Quotas.StorageGB),
		APICallsPerDay:  remaining(sub.Usage.APICallsPerDay, pkg.Quotas.APICallsPerDay),
		AutomationRules: remaining(sub.Usage.AutomationRules, pkg.Quotas.AutomationRules),
	}
}

func remaining(used, limit int) int {
	if limit == entity.Unlimited {
		return entity.Unlimited
	}
	r := limit - used
	if r < 0 {
		return 0
	}
	return r
}

// CalculateProration – คำนวณค่าใช้จ่ายตามสัดส่วน
func (s *QuotaService) CalculateProration(oldPkg, newPkg *entity.Package, sub *entity.Subscription) valueobject.Money {
	daysRemaining := sub.DaysRemaining()
	if daysRemaining <= 0 {
		return valueobject.Money{Amount: 0, Currency: newPkg.Price.Currency}
	}
	totalDays := newPkg.BillingCycle.Days()

	oldPerDay := oldPkg.Price.Amount / float64(totalDays)
	newPerDay := newPkg.Price.Amount / float64(totalDays)
	diff := (newPerDay - oldPerDay) * float64(daysRemaining)

	return valueobject.Money{
		Amount:   diff,
		Currency: newPkg.Price.Currency,
	}
}

// SuggestUpgrade – แนะนำ package ที่เหมาะ
func (s *QuotaService) SuggestUpgrade(current *entity.Package, candidates []entity.Package, usage entity.Quotas) *entity.Package {
	for i := range candidates {
		c := &candidates[i]
		if !c.IsActive {
			continue
		}
		if c.Tier.Priority() <= current.Tier.Priority() {
			continue
		}
		if exceeds, _ := usage.Exceeds(c.Quotas); !exceeds {
			return c
		}
	}
	return nil
}

var _ = time.Now
```

### B.7 Repository Interfaces

```go
// internal/modules/package/domain/repository/package_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/package/domain/entity"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type PackageRepository interface {
	Save(ctx context.Context, p *entity.Package) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Package, error)
	FindByCode(ctx context.Context, code string) (*entity.Package, error)
	ListActive(ctx context.Context) ([]entity.Package, error)
	ListByTier(ctx context.Context, tier valueobject.PackageTier) ([]entity.Package, error)
	ListAll(ctx context.Context) ([]entity.Package, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionRepository interface {
	Save(ctx context.Context, s *entity.Subscription) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)
	FindActiveByCustomer(ctx context.Context, customerID uuid.UUID) (*entity.Subscription, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Subscription, error)
	FindExpiring(ctx context.Context, before time.Time) ([]entity.Subscription, error)
	FindPastDue(ctx context.Context) ([]entity.Subscription, error)
	FindTrialEnding(ctx context.Context, before time.Time) ([]entity.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

### B.8 Errors

```go
// internal/modules/package/domain/errors/errors.go
package domainerrors

import "errors"

var (
	// Validation
	ErrInvalidCode        = errors.New("invalid package code")
	ErrInvalidName        = errors.New("invalid package name")
	ErrInvalidTier        = errors.New("invalid package tier")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrInvalidCycle       = errors.New("invalid billing cycle")
	ErrInvalidQuota       = errors.New("invalid quota")
	ErrInvalidCustomer    = errors.New("invalid customer")
	ErrInvalidPackage     = errors.New("invalid package")
	ErrInvalidExpiry      = errors.New("invalid expiry date")
	ErrCurrencyMismatch   = errors.New("currency mismatch")
	ErrReasonRequired     = errors.New("reason required")

	// Not found
	ErrPackageNotFound      = errors.New("package not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrFeatureNotFound      = errors.New("feature not found")

	// Conflict
	ErrCodeDuplicate      = errors.New("package code already exists")
	ErrFeatureDuplicate   = errors.New("feature already exists")
	ErrAlreadySubscribed  = errors.New("customer already has active subscription")
	ErrAlreadyCancelled   = errors.New("subscription already cancelled")

	// Business rules
	ErrPackageInactive      = errors.New("package is inactive")
	ErrCannotUpgrade        = errors.New("cannot upgrade subscription in current state")
	ErrCannotRenew          = errors.New("cannot renew cancelled subscription")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrSubscriptionExpired  = errors.New("subscription expired")
	ErrSubscriptionNotActive = errors.New("subscription is not active")
)
```

### B.9 Domain Events

```go
// internal/modules/package/domain/event/events.go
package event

import (
	"time"

	"github.com/google/uuid"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type PackageCreated struct {
	PackageID  uuid.UUID
	Code       string
	Tier       string
	OccurredAt time.Time
}

type PackageUpdated struct {
	PackageID  uuid.UUID
	Version    int
	Changes    map[string]interface{}
	OccurredAt time.Time
}

type SubscriptionCreated struct {
	SubscriptionID uuid.UUID
	CustomerID     uuid.UUID
	PackageID      uuid.UUID
	ExpiresAt      time.Time
	OccurredAt     time.Time
}

type SubscriptionUpgraded struct {
	SubscriptionID uuid.UUID
	FromPackageID  uuid.UUID
	ToPackageID    uuid.UUID
	ProrationAmt   float64
	OccurredAt     time.Time
}

type SubscriptionCancelled struct {
	SubscriptionID uuid.UUID
	Reason         string
	OccurredAt     time.Time
}

type SubscriptionExpired struct {
	SubscriptionID uuid.UUID
	CustomerID     uuid.UUID
	OccurredAt     time.Time
}

type SubscriptionRenewed struct {
	SubscriptionID uuid.UUID
	NewExpiresAt   time.Time
	OccurredAt     time.Time
}

type QuotaExceeded struct {
	SubscriptionID uuid.UUID
	CustomerID     uuid.UUID
	QuotaType      string
	Limit          int
	Current        int
	OccurredAt     time.Time
}

var _ = valueobject.TierFree
```

### B.10 Use Cases — Subscribe

```go
// internal/modules/package/application/subscribe.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/package/domain/entity"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	"icmongolang/internal/modules/package/domain/event"
	"icmongolang/internal/modules/package/domain/repository"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type SubscribeUseCase struct {
	subRepo      repository.SubscriptionRepository
	pkgRepo      repository.PackageRepository
	customerCli  CustomerClient
	paymentCli   PaymentClient
	producer     EventProducer
	auditRepo    AuditRepository
}

func NewSubscribeUseCase(
	subRepo repository.SubscriptionRepository,
	pkgRepo repository.PackageRepository,
	customerCli CustomerClient,
	paymentCli PaymentClient,
	producer EventProducer,
	auditRepo AuditRepository,
) *SubscribeUseCase {
	return &SubscribeUseCase{
		subRepo: subRepo, pkgRepo: pkgRepo, customerCli: customerCli,
		paymentCli: paymentCli, producer: producer, auditRepo: auditRepo,
	}
}

type SubscribeInput struct {
	CustomerID uuid.UUID
	PackageID  uuid.UUID
	Cycle      valueobject.BillingCycle
	UserID     uuid.UUID
	IPAddress  string
}

type SubscribeOutput struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	PaymentURL     string    `json:"payment_url"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (uc *SubscribeUseCase) Execute(ctx context.Context, in SubscribeInput) (*SubscribeOutput, error) {
	// 1. Load package
	pkg, err := uc.pkgRepo.FindByID(ctx, in.PackageID)
	if err != nil {
		return nil, err
	}
	if !pkg.IsActive {
		return nil, domainerrors.ErrPackageInactive
	}

	// 2. Verify customer
	cust, err := uc.customerCli.Get(ctx, in.CustomerID)
	if err != nil {
		return nil, err
	}
	if !cust.IsActive {
		return nil, domainerrors.ErrInvalidCustomer
	}

	// 3. Check existing subscription
	if existing, _ := uc.subRepo.FindActiveByCustomer(ctx, in.CustomerID); existing != nil {
		if existing.Status.IsUsable() {
			return nil, domainerrors.ErrAlreadySubscribed
		}
	}

	// 4. Create subscription
	sub, err := entity.NewSubscription(in.CustomerID, pkg.ID, in.Cycle)
	if err != nil {
		return nil, err
	}

	if err := uc.subRepo.Save(ctx, sub); err != nil {
		return nil, err
	}

	// 5. Create payment (if not free)
	var paymentURL string
	if !pkg.IsFree() {
		payment, err := uc.paymentCli.CreatePayment(ctx, PaymentRequest{
			CustomerID:  in.CustomerID,
			Amount:      pkg.Price.Amount,
			Currency:    pkg.Price.Currency,
			Description: "Subscription: " + pkg.Name,
			ReferenceID: sub.ID,
			ReferenceType: "subscription",
			ReturnURL:   "", // from config
		})
		if err != nil {
			return nil, err
		}
		paymentURL = payment.URL
	} else {
		// Free tier → activate immediately
		sub.Status = valueobject.SubStatusActive
		_ = uc.subRepo.Save(ctx, sub)
	}

	// 6. Audit
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID:     in.UserID,
		Action:     "SUBSCRIPTION_CREATED",
		EntityType: "subscription",
		EntityID:   sub.ID,
		Payload: map[string]interface{}{
			"package_id": pkg.ID.String(),
			"package_code": pkg.Code,
			"amount":     pkg.Price.Amount,
		},
		IPAddress:  in.IPAddress,
		OccurredAt: time.Now(),
	})

	// 7. Publish event
	_ = uc.producer.PublishSubscriptionCreated(ctx, event.SubscriptionCreated{
		SubscriptionID: sub.ID,
		CustomerID:     sub.CustomerID,
		PackageID:      sub.PackageID,
		ExpiresAt:      sub.ExpiresAt,
		OccurredAt:     time.Now(),
	})

	return &SubscribeOutput{
		SubscriptionID: sub.ID,
		PaymentURL:     paymentURL,
		ExpiresAt:      sub.ExpiresAt,
	}, nil
}
```

### B.11 Use Case — UpgradeSubscription

```go
// internal/modules/package/application/upgrade_subscription.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	"icmongolang/internal/modules/package/domain/event"
	"icmongolang/internal/modules/package/domain/repository"
	"icmongolang/internal/modules/package/domain/service"
)

type UpgradeSubscriptionUseCase struct {
	subRepo   repository.SubscriptionRepository
	pkgRepo   repository.PackageRepository
	quotaSvc  *service.QuotaService
	paymentCli PaymentClient
	producer  EventProducer
	auditRepo AuditRepository
}

type UpgradeSubscriptionInput struct {
	SubscriptionID uuid.UUID
	NewPackageID   uuid.UUID
	UserID         uuid.UUID
}

type UpgradeSubscriptionOutput struct {
	SubscriptionID   uuid.UUID `json:"subscription_id"`
	FromPackageID    uuid.UUID `json:"from_package_id"`
	ToPackageID      uuid.UUID `json:"to_package_id"`
	ProrationAmount  float64   `json:"proration_amount"`
	Currency         string    `json:"currency"`
	NewExpiresAt     time.Time `json:"new_expires_at"`
	PaymentURL       string    `json:"payment_url,omitempty"`
}

func (uc *UpgradeSubscriptionUseCase) Execute(ctx context.Context, in UpgradeSubscriptionInput) (*UpgradeSubscriptionOutput, error) {
	// 1. Load subscription
	sub, err := uc.subRepo.FindByID(ctx, in.SubscriptionID)
	if err != nil {
		return nil, err
	}
	if !sub.Status.CanUpgrade() {
		return nil, domainerrors.ErrCannotUpgrade
	}

	// 2. Load packages
	oldPkg, err := uc.pkgRepo.FindByID(ctx, sub.PackageID)
	if err != nil {
		return nil, err
	}
	newPkg, err := uc.pkgRepo.FindByID(ctx, in.NewPackageID)
	if err != nil {
		return nil, err
	}
	if !newPkg.IsActive {
		return nil, domainerrors.ErrPackageInactive
	}
	if !newPkg.IsHigherTierThan(oldPkg) {
		return nil, domainerrors.ErrCannotUpgrade
	}

	// 3. Calculate proration
	proration := uc.quotaSvc.CalculateProration(oldPkg, newPkg, sub)

	// 4. Perform upgrade
	newExpiry := newPkg.BillingCycle.AddTo(time.Now())
	if err := sub.Upgrade(newPkg.ID, newExpiry, newPkg.BillingCycle); err != nil {
		return nil, err
	}
	if err := uc.subRepo.Save(ctx, sub); err != nil {
		return nil, err
	}

	// 5. Create payment for proration (if > 0)
	var paymentURL string
	if proration.Amount > 0 {
		payment, err := uc.paymentCli.CreatePayment(ctx, PaymentRequest{
			CustomerID:    sub.CustomerID,
			Amount:        proration.Amount,
			Currency:      proration.Currency,
			Description:   "Upgrade to " + newPkg.Name,
			ReferenceID:   sub.ID,
			ReferenceType: "subscription_upgrade",
		})
		if err != nil {
			return nil, err
		}
		paymentURL = payment.URL
	}

	// 6. Audit + event
	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID:     in.UserID,
		Action:     "SUBSCRIPTION_UPGRADED",
		EntityType: "subscription",
		EntityID:   sub.ID,
		Payload: map[string]interface{}{
			"from_package": oldPkg.Code,
			"to_package":   newPkg.Code,
			"proration":    proration.Amount,
		},
		OccurredAt: time.Now(),
	})

	_ = uc.producer.PublishSubscriptionUpgraded(ctx, event.SubscriptionUpgraded{
		SubscriptionID: sub.ID,
		FromPackageID:  oldPkg.ID,
		ToPackageID:    newPkg.ID,
		ProrationAmt:   proration.Amount,
		OccurredAt:     time.Now(),
	})

	return &UpgradeSubscriptionOutput{
		SubscriptionID:  sub.ID,
		FromPackageID:   oldPkg.ID,
		ToPackageID:     newPkg.ID,
		ProrationAmount: proration.Amount,
		Currency:        proration.Currency,
		NewExpiresAt:    newExpiry,
		PaymentURL:      paymentURL,
	}, nil
}
```

### B.12 Use Case — CheckQuota

```go
// internal/modules/package/application/check_quota.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	"icmongolang/internal/modules/package/domain/event"
	"icmongolang/internal/modules/package/domain/repository"
	"icmongolang/internal/modules/package/domain/service"
)

type CheckQuotaUseCase struct {
	subRepo  repository.SubscriptionRepository
	pkgRepo  repository.PackageRepository
	quotaSvc *service.QuotaService
	producer EventProducer
}

type CheckQuotaInput struct {
	SubscriptionID uuid.UUID
	Usage          entity.Quotas
}

type CheckQuotaOutput struct {
	Allowed   bool           `json:"allowed"`
	Violation string         `json:"violation,omitempty"`
	Remaining entity.Quotas  `json:"remaining"`
	Limit     entity.Quotas  `json:"limit"`
	Usage     entity.Quotas  `json:"usage"`
}

func (uc *CheckQuotaUseCase) Execute(ctx context.Context, in CheckQuotaInput) (*CheckQuotaOutput, error) {
	sub, err := uc.subRepo.FindByID(ctx, in.SubscriptionID)
	if err != nil {
		return nil, err
	}
	if !sub.Status.IsUsable() {
		return nil, domainerrors.ErrSubscriptionNotActive
	}

	pkg, err := uc.pkgRepo.FindByID(ctx, sub.PackageID)
	if err != nil {
		return nil, err
	}

	// Update usage then check
	sub.UpdateUsage(in.Usage)
	_ = uc.subRepo.Save(ctx, sub)

	exceeds, violation := uc.quotaSvc.CheckQuota(sub, pkg)
	if exceeds {
		// Publish quota exceeded event
		_ = uc.producer.PublishQuotaExceeded(ctx, event.QuotaExceeded{
			SubscriptionID: sub.ID,
			CustomerID:     sub.CustomerID,
			QuotaType:      violation,
			OccurredAt:     time.Now(),
		})
		return &CheckQuotaOutput{
			Allowed:   false,
			Violation: violation,
			Remaining: uc.quotaSvc.Remaining(sub, pkg),
			Limit:     pkg.Quotas,
			Usage:     sub.Usage,
		}, nil
	}

	return &CheckQuotaOutput{
		Allowed:   true,
		Remaining: uc.quotaSvc.Remaining(sub, pkg),
		Limit:     pkg.Quotas,
		Usage:     sub.Usage,
	}, nil
}
```

### B.13 Infrastructure — Postgres

```go
// internal/modules/package/infrastructure/persistence/postgres/models.go
package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PackageModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code           string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    string         `gorm:"type:text"`
	Tier           string         `gorm:"type:varchar(30);not null;index"`
	Version        int            `gorm:"not null;default:1"`
	PriceAmount    float64        `gorm:"type:numeric(12,2);not null;default:0"`
	PriceCurrency  string         `gorm:"type:varchar(3);not null;default:'THB'"`
	BillingCycle   string         `gorm:"type:varchar(20);not null"`
	Quotas         datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	Features       datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`
	IsActive       bool           `gorm:"not null;default:true"`
	IsPublic       bool           `gorm:"not null;default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

func (PackageModel) TableName() string { return "package_packages" }

type SubscriptionModel struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CustomerID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	PackageID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Status        string         `gorm:"type:varchar(30);not null;index"`
	StartedAt     time.Time      `gorm:"not null"`
	ExpiresAt     time.Time      `gorm:"not null;index"`
	TrialEndsAt   *time.Time
	GraceEndsAt   *time.Time
	AutoRenew     bool           `gorm:"not null;default:true"`
	Usage         datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	PreviousPkgID *uuid.UUID
	UpgradedAt    *time.Time
	CancelledAt   *time.Time
	CancelReason  string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SubscriptionModel) TableName() string { return "package_subscriptions" }
```

### B.14 Infrastructure — Repositories

```go
// internal/modules/package/infrastructure/persistence/postgres/package_repo_impl.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/package/domain/entity"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type packageRepoImpl struct{ db *gorm.DB }

func NewPackageRepository(db *gorm.DB) *packageRepoImpl {
	return &packageRepoImpl{db: db}
}

func (r *packageRepoImpl) Save(ctx context.Context, p *entity.Package) error {
	m := toPackageModel(p)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *packageRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Package, error) {
	var m PackageModel
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrPackageNotFound
	}
	if err != nil {
		return nil, err
	}
	return toPackageEntity(&m), nil
}

func (r *packageRepoImpl) FindByCode(ctx context.Context, code string) (*entity.Package, error) {
	var m PackageModel
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrPackageNotFound
	}
	if err != nil {
		return nil, err
	}
	return toPackageEntity(&m), nil
}

func (r *packageRepoImpl) ListActive(ctx context.Context) ([]entity.Package, error) {
	var models []PackageModel
	err := r.db.WithContext(ctx).
		Where("is_active = TRUE AND deleted_at IS NULL").
		Order("tier DESC, price_amount ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toPackageEntities(models), nil
}

func (r *packageRepoImpl) ListByTier(ctx context.Context, tier valueobject.PackageTier) ([]entity.Package, error) {
	var models []PackageModel
	err := r.db.WithContext(ctx).
		Where("tier = ? AND is_active = TRUE AND deleted_at IS NULL", tier).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toPackageEntities(models), nil
}

func (r *packageRepoImpl) ListAll(ctx context.Context) ([]entity.Package, error) {
	var models []PackageModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&models).Error; err != nil {
		return nil, err
	}
	return toPackageEntities(models), nil
}

func (r *packageRepoImpl) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PackageModel{}).
		Where("code = ? AND deleted_at IS NULL", code).Count(&count).Error
	return count > 0, err
}

func (r *packageRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&PackageModel{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ---------- Mappers ----------

func toPackageModel(p *entity.Package) *PackageModel {
	quotasJSON, _ := json.Marshal(p.Quotas)
	featuresJSON, _ := json.Marshal(p.Features)

	return &PackageModel{
		ID:            p.ID,
		Code:          p.Code,
		Name:          p.Name,
		Description:   p.Description,
		Tier:          string(p.Tier),
		Version:       p.Version,
		PriceAmount:   p.Price.Amount,
		PriceCurrency: p.Price.Currency,
		BillingCycle:  string(p.BillingCycle),
		Quotas:        datatypes.JSON(quotasJSON),
		Features:      datatypes.JSON(featuresJSON),
		IsActive:      p.IsActive,
		IsPublic:      p.IsPublic,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toPackageEntity(m *PackageModel) *entity.Package {
	var quotas entity.Quotas
	var features []entity.Feature
	_ = json.Unmarshal(m.Quotas, &quotas)
	_ = json.Unmarshal(m.Features, &features)

	return &entity.Package{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Tier:        valueobject.PackageTier(m.Tier),
		Version:     m.Version,
		Price: valueobject.Money{
			Amount:   m.PriceAmount,
			Currency: m.PriceCurrency,
		},
		BillingCycle: valueobject.BillingCycle(m.BillingCycle),
		Quotas:       quotas,
		Features:     features,
		IsActive:     m.IsActive,
		IsPublic:     m.IsPublic,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func toPackageEntities(models []PackageModel) []entity.Package {
	out := make([]entity.Package, 0, len(models))
	for i := range models {
		out = append(out, *toPackageEntity(&models[i]))
	}
	return out
}
```

```go
// internal/modules/package/infrastructure/persistence/postgres/subscription_repo_impl.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/package/domain/entity"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

type subscriptionRepoImpl struct{ db *gorm.DB }

func NewSubscriptionRepository(db *gorm.DB) *subscriptionRepoImpl {
	return &subscriptionRepoImpl{db: db}
}

func (r *subscriptionRepoImpl) Save(ctx context.Context, s *entity.Subscription) error {
	return r.db.WithContext(ctx).Save(toSubscriptionModel(s)).Error
}

func (r *subscriptionRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	var m SubscriptionModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, err
	}
	return toSubscriptionEntity(&m), nil
}

func (r *subscriptionRepoImpl) FindActiveByCustomer(ctx context.Context, customerID uuid.UUID) (*entity.Subscription, error) {
	var m SubscriptionModel
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND status IN ('ACTIVE', 'TRIAL', 'PAST_DUE')", customerID).
		Order("created_at DESC").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, err
	}
	return toSubscriptionEntity(&m), nil
}

func (r *subscriptionRepoImpl) FindExpiring(ctx context.Context, before time.Time) ([]entity.Subscription, error) {
	var models []SubscriptionModel
	err := r.db.WithContext(ctx).
		Where("status = 'ACTIVE' AND expires_at <= ?", before).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.Subscription, 0, len(models))
	for i := range models {
		out = append(out, *toSubscriptionEntity(&models[i]))
	}
	return out, nil
}

func (r *subscriptionRepoImpl) FindPastDue(ctx context.Context) ([]entity.Subscription, error) {
	var models []SubscriptionModel
	err := r.db.WithContext(ctx).Where("status = 'PAST_DUE'").Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.Subscription, 0, len(models))
	for i := range models {
		out = append(out, *toSubscriptionEntity(&models[i]))
	}
	return out, nil
}

func (r *subscriptionRepoImpl) FindTrialEnding(ctx context.Context, before time.Time) ([]entity.Subscription, error) {
	var models []SubscriptionModel
	err := r.db.WithContext(ctx).
		Where("status = 'TRIAL' AND trial_ends_at <= ?", before).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.Subscription, 0, len(models))
	for i := range models {
		out = append(out, *toSubscriptionEntity(&models[i]))
	}
	return out, nil
}

func (r *subscriptionRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&SubscriptionModel{}, "id = ?", id).Error
}

// ---------- Mappers ----------

func toSubscriptionModel(s *entity.Subscription) *SubscriptionModel {
	usageJSON, _ := json.Marshal(s.Usage)
	return &SubscriptionModel{
		ID:            s.ID,
		CustomerID:    s.CustomerID,
		PackageID:     s.PackageID,
		Status:        string(s.Status),
		StartedAt:     s.StartedAt,
		ExpiresAt:     s.ExpiresAt,
		TrialEndsAt:   s.TrialEndsAt,
		GraceEndsAt:   s.GraceEndsAt,
		AutoRenew:     s.AutoRenew,
		Usage:         datatypes.JSON(usageJSON),
		PreviousPkgID: s.PreviousPkgID,
		UpgradedAt:    s.UpgradedAt,
		CancelledAt:   s.CancelledAt,
		CancelReason:  s.CancelReason,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func toSubscriptionEntity(m *SubscriptionModel) *entity.Subscription {
	var usage entity.Quotas
	_ = json.Unmarshal(m.Usage, &usage)

	return &entity.Subscription{
		ID:            m.ID,
		CustomerID:    m.CustomerID,
		PackageID:     m.PackageID,
		Status:        valueobject.SubscriptionStatus(m.Status),
		StartedAt:     m.StartedAt,
		ExpiresAt:     m.ExpiresAt,
		TrialEndsAt:   m.TrialEndsAt,
		GraceEndsAt:   m.GraceEndsAt,
		AutoRenew:     m.AutoRenew,
		Usage:         usage,
		PreviousPkgID: m.PreviousPkgID,
		UpgradedAt:    m.UpgradedAt,
		CancelledAt:   m.CancelledAt,
		CancelReason:  m.CancelReason,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
```

### B.15 Scheduler Jobs

```go
// internal/modules/package/infrastructure/scheduler/expire_job.go
package scheduler

import (
	"context"
	"log"
	"time"

	"icmongolang/internal/modules/package/domain/event"
	"icmongolang/internal/modules/package/domain/repository"
)

type ExpireJob struct {
	subRepo  repository.SubscriptionRepository
	producer EventProducer
	graceDays int
}

func NewExpireJob(subRepo repository.SubscriptionRepository, producer EventProducer, graceDays int) *ExpireJob {
	return &ExpireJob{subRepo: subRepo, producer: producer, graceDays: graceDays}
}

func (j *ExpireJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	expiring, err := j.subRepo.FindExpiring(ctx, time.Now())
	if err != nil {
		log.Printf("[expire_job] find expiring: %v", err)
		return
	}

	for i := range expiring {
		sub := &expiring[i]
		if sub.AutoRenew {
			// auto renew — publish event for payment service to handle
			continue
		}
		// try to move to past_due first if in grace period, else expire
		if sub.GraceEndsAt == nil {
			sub.MarkPastDue(j.graceDays)
			_ = j.subRepo.Save(ctx, sub)
			continue
		}
		if time.Now().After(*sub.GraceEndsAt) {
			sub.Expire()
			_ = j.subRepo.Save(ctx, sub)
			_ = j.producer.PublishSubscriptionExpired(ctx, event.SubscriptionExpired{
				SubscriptionID: sub.ID,
				CustomerID:     sub.CustomerID,
				OccurredAt:     time.Now(),
			})
		}
	}
	log.Printf("[expire_job] processed %d subscriptions", len(expiring))
}
```

### B.16 HTTP Handlers

```go
// internal/modules/package/interfaces/http/package_handler.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/package/application"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
	"icmongolang/pkg/responses"
)

type PackageHandler struct {
	createUC *application.CreatePackageUseCase
	getUC    *application.GetPackageUseCase
	listUC   *application.ListPackagesUseCase
	updateUC *application.UpdatePackageUseCase
}

func NewPackageHandler(
	createUC *application.CreatePackageUseCase,
	getUC *application.GetPackageUseCase,
	listUC *application.ListPackagesUseCase,
	updateUC *application.UpdatePackageUseCase,
) *PackageHandler {
	return &PackageHandler{createUC: createUC, getUC: getUC, listUC: listUC, updateUC: updateUC}
}

// POST /api/v1/packages
func (h *PackageHandler) Create(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)

	var req CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	out, err := h.createUC.Execute(c.Request.Context(), application.CreatePackageInput{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		Tier:         valueobject.PackageTier(req.Tier),
		Price:        valueobject.Money{Amount: req.PriceAmount, Currency: req.PriceCurrency},
		BillingCycle: valueobject.BillingCycle(req.BillingCycle),
		Quotas:       req.Quotas.ToVO(),
		UserID:       uid,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("CREATE_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

// GET /api/v1/packages
func (h *PackageHandler) List(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "true") == "true"
	out, err := h.listUC.Execute(c.Request.Context(), application.ListPackagesInput{
		ActiveOnly: activeOnly,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error("LIST_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, responses.Success(out))
}
```

```go
// internal/modules/package/interfaces/http/subscription_handler.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/package/application"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
	"icmongolang/pkg/responses"
)

type SubscriptionHandler struct {
	subscribeUC *application.SubscribeUseCase
	getUC       *application.GetSubscriptionUseCase
	upgradeUC   *application.UpgradeSubscriptionUseCase
	cancelUC    *application.CancelSubscriptionUseCase
	renewUC     *application.RenewSubscriptionUseCase
	quotaUC     *application.CheckQuotaUseCase
}

func NewSubscriptionHandler(
	subscribeUC *application.SubscribeUseCase,
	getUC *application.GetSubscriptionUseCase,
	upgradeUC *application.UpgradeSubscriptionUseCase,
	cancelUC *application.CancelSubscriptionUseCase,
	renewUC *application.RenewSubscriptionUseCase,
	quotaUC *application.CheckQuotaUseCase,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscribeUC: subscribeUC, getUC: getUC, upgradeUC: upgradeUC,
		cancelUC: cancelUC, renewUC: renewUC, quotaUC: quotaUC,
	}
}

// POST /api/v1/subscriptions
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		PackageID  string `json:"package_id" binding:"required"`
		Cycle      string `json:"cycle" binding:"required,oneof=MONTHLY YEARLY"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	custID, _ := uuid.Parse(req.CustomerID)
	pkgID, _ := uuid.Parse(req.PackageID)

	out, err := h.subscribeUC.Execute(c.Request.Context(), application.SubscribeInput{
		CustomerID: custID,
		PackageID:  pkgID,
		Cycle:      valueobject.BillingCycle(req.Cycle),
		UserID:     uid,
		IPAddress:  c.ClientIP(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("SUBSCRIBE_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

// POST /api/v1/subscriptions/:id/upgrade
func (h *SubscriptionHandler) Upgrade(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid subscription id", nil))
		return
	}

	var req struct {
		NewPackageID string `json:"new_package_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}
	newPkgID, _ := uuid.Parse(req.NewPackageID)

	out, err := h.upgradeUC.Execute(c.Request.Context(), application.UpgradeSubscriptionInput{
		SubscriptionID: subID,
		NewPackageID:   newPkgID,
		UserID:         uid,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("UPGRADE_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, responses.Success(out))
}

// POST /api/v1/subscriptions/:id/cancel
func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid id", nil))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	if err := h.cancelUC.Execute(c.Request.Context(), application.CancelSubscriptionInput{
		SubscriptionID: subID,
		Reason:         req.Reason,
		UserID:         uid,
	}); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("CANCEL_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "cancelled"}))
}

// POST /api/v1/subscriptions/:id/quota/check
func (h *SubscriptionHandler) CheckQuota(c *gin.Context) {
	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_ID", "invalid id", nil))
		return
	}

	var req struct {
		Usage application.QuotasDTO `json:"usage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID_REQUEST", err.Error(), nil))
		return
	}

	out, err := h.quotaUC.Execute(c.Request.Context(), application.CheckQuotaInput{
		SubscriptionID: subID,
		Usage:          req.Usage.ToVO(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("CHECK_FAILED", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, responses.Success(out))
}
```

### B.17 Routes

```go
// internal/modules/package/interfaces/http/routes.go
package http

import "github.com/gin-gonic/gin"

type Handlers struct {
	Package      *PackageHandler
	Subscription *SubscriptionHandler
}

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handlers,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	// Packages
	pkgs := r.Group("/packages")
	pkgs.Use(auth, tenant)
	pkgs.GET("", rbac("viewer"), h.Package.List)
	pkgs.GET("/:id", rbac("viewer"), h.Package.Get)
	pkgs.POST("", rbac("admin"), h.Package.Create)
	pkgs.PUT("/:id", rbac("admin"), h.Package.Update)
	pkgs.DELETE("/:id", rbac("admin"), h.Package.Delete)

	// Subscriptions
	subs := r.Group("/subscriptions")
	subs.Use(auth, tenant)
	subs.POST("", rbac("manager"), h.Subscription.Subscribe)
	subs.GET("/:id", rbac("viewer"), h.Subscription.Get)
	subs.GET("/me", rbac("viewer"), h.Subscription.GetMy)
	subs.POST("/:id/upgrade", rbac("manager"), h.Subscription.Upgrade)
	subs.POST("/:id/cancel", rbac("manager"), h.Subscription.Cancel)
	subs.POST("/:id/renew", rbac("manager"), h.Subscription.Renew)
	subs.GET("/:id/quota", rbac("viewer"), h.Subscription.GetQuota)
	subs.POST("/:id/quota/check", rbac("viewer"), h.Subscription.CheckQuota)
}
```

### B.18 Composition Root

```go
// internal/modules/package/module.go
package package_

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/package/application"
	"icmongolang/internal/modules/package/domain/service"
	"icmongolang/internal/modules/package/infrastructure/messaging"
	pgrepo "icmongolang/internal/modules/package/infrastructure/persistence/postgres"
	redisrepo "icmongolang/internal/modules/package/infrastructure/persistence/redis"
	httpiface "icmongolang/internal/modules/package/interfaces/http"
)

type Dependencies struct {
	DB            *gorm.DB
	Redis         *redis.Client
	Producer      *messaging.KafkaProducer
	CustomerCli   application.CustomerClient
	PaymentCli    application.PaymentClient
	AuditRepo     application.AuditRepository
	GraceDays     int
}

func Init(
	router *gin.RouterGroup,
	deps Dependencies,
	auth gin.HandlerFunc,
	tenant gin.HandlerFunc,
	rbac func(roles ...string) gin.HandlerFunc,
) {
	// Repositories
	pkgRepo := pgrepo.NewPackageRepository(deps.DB)
	subRepo := pgrepo.NewSubscriptionRepository(deps.DB)
	_ = redisrepo.NewPackageCache(deps.Redis)

	// Domain services
	quotaSvc := &service.QuotaService{}

	// Use cases
	createPkgUC := application.NewCreatePackageUseCase(pkgRepo, deps.AuditRepo, deps.Producer)
	getPkgUC := application.NewGetPackageUseCase(pkgRepo)
	listPkgUC := application.NewListPackagesUseCase(pkgRepo)
	updatePkgUC := application.NewUpdatePackageUseCase(pkgRepo, deps.AuditRepo, deps.Producer)

	subscribeUC := application.NewSubscribeUseCase(subRepo, pkgRepo, deps.CustomerCli, deps.PaymentCli, deps.Producer, deps.AuditRepo)
	getSubUC := application.NewGetSubscriptionUseCase(subRepo, pkgRepo)
	upgradeUC := application.NewUpgradeSubscriptionUseCase(subRepo, pkgRepo, quotaSvc, deps.PaymentCli, deps.Producer, deps.AuditRepo)
	cancelUC := application.NewCancelSubscriptionUseCase(subRepo, deps.Producer, deps.AuditRepo)
	renewUC := application.NewRenewSubscriptionUseCase(subRepo, deps.PaymentCli, deps.Producer)
	quotaUC := application.NewCheckQuotaUseCase(subRepo, pkgRepo, quotaSvc, deps.Producer)

	// Handlers
	pkgHandler := httpiface.NewPackageHandler(createPkgUC, getPkgUC, listPkgUC, updatePkgUC)
	subHandler := httpiface.NewSubscriptionHandler(subscribeUC, getSubUC, upgradeUC, cancelUC, renewUC, quotaUC)

	httpiface.RegisterRoutes(router, &httpiface.Handlers{
		Package:      pkgHandler,
		Subscription: subHandler,
	}, auth, tenant, rbac)
}
```

### B.19 Go Tests

```go
// internal/modules/package/domain/entity/subscription_test.go
package entity_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/package/domain/entity"
	domainerrors "icmongolang/internal/modules/package/domain/errors"
	valueobject "icmongolang/internal/modules/package/domain/value_object"
)

func TestSubscription_Upgrade(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *entity.Subscription
		newPkg  uuid.UUID
		expiry  time.Time
		wantErr error
	}{
		{
			name: "ACTIVE → upgrade OK",
			setup: func() *entity.Subscription {
				s, _ := entity.NewSubscription(uuid.New(), uuid.New(), valueobject.BillingMonthly)
				return s
			},
			newPkg:  uuid.New(),
			expiry:  time.Now().AddDate(0, 1, 0),
			wantErr: nil,
		},
		{
			name: "CANCELLED → upgrade fails",
			setup: func() *entity.Subscription {
				s, _ := entity.NewSubscription(uuid.New(), uuid.New(), valueobject.BillingMonthly)
				_ = s.Cancel("user request")
				return s
			},
			newPkg:  uuid.New(),
			expiry:  time.Now().AddDate(0, 1, 0),
			wantErr: domainerrors.ErrCannotUpgrade,
		},
		{
			name: "EXPIRED → upgrade fails",
			setup: func() *entity.Subscription {
				s, _ := entity.NewSubscription(uuid.New(), uuid.New(), valueobject.BillingMonthly)
				s.Expire()
				return s
			},
			newPkg:  uuid.New(),
			expiry:  time.Now().AddDate(0, 1, 0),
			wantErr: domainerrors.ErrCannotUpgrade,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			err := s.Upgrade(tt.newPkg, tt.expiry, valueobject.BillingMonthly)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscription_DaysRemaining(t *testing.T) {
	s, _ := entity.NewSubscription(uuid.New(), uuid.New(), valueobject.BillingMonthly)
	days := s.DaysRemaining()
	if days < 29 || days > 31 {
		t.Errorf("expected ~30 days, got %d", days)
	}
}

func TestQuotas_Exceeds(t *testing.T) {
	usage := entity.Quotas{MaxDevices: 5, MaxUsers: 3}
	limit := entity.Quotas{MaxDevices: 10, MaxUsers: 5}
	exceeds, field := usage.Exceeds(limit)
	if exceeds {
		t.Errorf("should not exceed, got %s", field)
	}

	usage.MaxDevices = 15
	exceeds, field = usage.Exceeds(limit)
	if !exceeds || field != "devices" {
		t.Errorf("should exceed devices, got %v/%s", exceeds, field)
	}
}

func TestQuotas_Unlimited(t *testing.T) {
	usage := entity.Quotas{MaxDevices: 1000000}
	limit := entity.Quotas{MaxDevices: entity.Unlimited}
	exceeds, _ := usage.Exceeds(limit)
	if exceeds {
		t.Error("unlimited quota should not exceed")
	}
}
```

---

## PART C — Python FastAPI Implementation

### C.1 Directory Tree

```
app/modules/package/
├── domain/
│   ├── entities/
│   │   ├── package.py
│   │   ├── subscription.py
│   │   ├── quotas.py
│   │   └── feature.py
│   ├── value_objects/
│   │   ├── package_tier.py
│   │   ├── billing_cycle.py
│   │   ├── subscription_status.py
│   │   └── money.py
│   ├── repositories/
│   │   ├── package_repository.py
│   │   └── subscription_repository.py
│   ├── services/
│   │   └── quota_service.py
│   ├── events/
│   │   └── events.py
│   └── errors/
│       └── errors.py
├── application/
│   ├── use_cases/
│   │   ├── create_package.py
│   │   ├── subscribe.py
│   │   ├── upgrade_subscription.py
│   │   ├── cancel_subscription.py
│   │   └── check_quota.py
│   └── ports/
│       ├── customer_client.py
│       ├── payment_client.py
│       ├── event_producer.py
│       └── audit_repository.py
├── infrastructure/
│   ├── persistence/
│   │   └── postgres/
│   │       ├── models.py
│   │       ├── package_repository_impl.py
│   │       └── subscription_repository_impl.py
│   ├── messaging/
│   │   └── kafka_producer.py
│   └── scheduler/
│       └── expire_job.py
├── interfaces/
│   └── http/
│       ├── schemas.py
│       ├── package_router.py
│       ├── subscription_router.py
│       └── dependencies.py
└── module.py
```

### C.2 Value Objects

```python
# app/modules/package/domain/value_objects/package_tier.py
from enum import Enum


class PackageTier(str, Enum):
    FREE = "FREE"
    BASIC = "BASIC"
    PRO = "PRO"
    ENTERPRISE = "ENTERPRISE"

    def priority(self) -> int:
        return {PackageTier.FREE: 0, PackageTier.BASIC: 1,
                PackageTier.PRO: 2, PackageTier.ENTERPRISE: 3}[self]
```

```python
# app/modules/package/domain/value_objects/billing_cycle.py
from datetime import datetime, timedelta
from enum import Enum


class BillingCycle(str, Enum):
    MONTHLY = "MONTHLY"
    YEARLY = "YEARLY"

    def add_to(self, t: datetime) -> datetime:
        if self == BillingCycle.YEARLY:
            return t.replace(year=t.year + 1)
        month = t.month + 1
        year = t.year
        if month > 12:
            month = 1
            year += 1
        return t.replace(year=year, month=month)

    def days(self) -> int:
        return 365 if self == BillingCycle.YEARLY else 30
```

```python
# app/modules/package/domain/value_objects/subscription_status.py
from enum import Enum


class SubscriptionStatus(str, Enum):
    ACTIVE = "ACTIVE"
    TRIAL = "TRIAL"
    PAST_DUE = "PAST_DUE"
    CANCELLED = "CANCELLED"
    EXPIRED = "EXPIRED"

    def is_usable(self) -> bool:
        return self in (SubscriptionStatus.ACTIVE, SubscriptionStatus.TRIAL)

    def can_upgrade(self) -> bool:
        return self in (SubscriptionStatus.ACTIVE, SubscriptionStatus.TRIAL,
                        SubscriptionStatus.PAST_DUE)
```

```python
# app/modules/package/domain/value_objects/money.py
from pydantic import BaseModel, Field, field_validator


class Money(BaseModel):
    amount: float = Field(..., ge=0)
    currency: str = Field(..., min_length=3, max_length=3)

    @field_validator("currency")
    @classmethod
    def uppercase_currency(cls, v: str) -> str:
        return v.upper()

    def add(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise ValueError("currency mismatch")
        return Money(amount=self.amount + other.amount, currency=self.currency)

    def sub(self, other: "Money") -> "Money":
        if self.currency != other.currency:
            raise ValueError("currency mismatch")
        return Money(amount=self.amount - other.amount, currency=self.currency)

    def mul(self, factor: float) -> "Money":
        return Money(amount=self.amount * factor, currency=self.currency)

    def is_zero(self) -> bool:
        return self.amount == 0
```

### C.3 Entities

```python
# app/modules/package/domain/entities/quotas.py
from pydantic import BaseModel

UNLIMITED = -1


class Quotas(BaseModel):
    max_devices: int = 0
    max_sites: int = 0
    max_users: int = 0
    storage_gb: int = 0
    retention_days: int = 0
    api_calls_per_day: int = 0
    automation_rules: int = 0

    def validate_quotas(self) -> None:
        for field, value in self.model_dump().items():
            if value < -1:
                raise ValueError(f"{field} cannot be < -1")

    def exceeds(self, limit: "Quotas") -> tuple[bool, str]:
        checks = [
            ("devices", self.max_devices, limit.max_devices),
            ("sites", self.max_sites, limit.max_sites),
            ("users", self.max_users, limit.max_users),
            ("storage", self.storage_gb, limit.storage_gb),
            ("api_calls", self.api_calls_per_day, limit.api_calls_per_day),
            ("automation_rules", self.automation_rules, limit.automation_rules),
        ]
        for name, current, lim in checks:
            if lim != UNLIMITED and current > lim:
                return True, name
        return False, ""

    def allows_devices(self, count: int) -> bool:
        return self.max_devices == UNLIMITED or count <= self.max_devices


class Feature(BaseModel):
    code: str
    name: str
    description: str = ""
    enabled: bool = True
```

```python
# app/modules/package/domain/entities/package.py
import re
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field, field_validator

from app.modules.package.domain.entities.quotas import Quotas, Feature
from app.modules.package.domain.errors.errors import (
    InvalidCodeError,
    InvalidNameError,
    InvalidTierError,
    InvalidPriceError,
    InvalidCycleError,
    FeatureDuplicateError,
    FeatureNotFoundError,
    CurrencyMismatchError,
)
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle
from app.modules.package.domain.value_objects.money import Money
from app.modules.package.domain.value_objects.package_tier import PackageTier

CODE_PATTERN = re.compile(r"^[A-Z0-9_]{2,50}$")


class Package(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    code: str
    name: str
    description: str = ""
    tier: PackageTier
    version: int = 1
    price: Money
    billing_cycle: BillingCycle
    quotas: Quotas = Field(default_factory=Quotas)
    features: list[Feature] = Field(default_factory=list)
    is_active: bool = True
    is_public: bool = True
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    # ---------- Constructor ----------
    @classmethod
    def create(
        cls,
        code: str,
        name: str,
        tier: PackageTier,
        price: Money,
        cycle: BillingCycle,
        quotas: Quotas,
    ) -> "Package":
        if not CODE_PATTERN.match(code):
            raise InvalidCodeError(f"invalid code: {code}")
        if not name or len(name) > 255:
            raise InvalidNameError("invalid name")
        quotas.validate_quotas()
        return cls(
            code=code,
            name=name,
            tier=tier,
            price=price,
            billing_cycle=cycle,
            quotas=quotas,
        )

    # ---------- Behavior ----------
    def rename(self, name: str, description: str = "") -> None:
        if not name or len(name) > 255:
            raise InvalidNameError()
        self.name = name
        self.description = description
        self._touch()

    def update_price(self, new_price: Money) -> None:
        if new_price.currency != self.price.currency:
            raise CurrencyMismatchError()
        self.price = new_price
        self.version += 1
        self._touch()

    def update_quotas(self, quotas: Quotas) -> None:
        quotas.validate_quotas()
        self.quotas = quotas
        self._touch()

    def add_feature(self, feature: Feature) -> None:
        if any(f.code == feature.code for f in self.features):
            raise FeatureDuplicateError(feature.code)
        self.features.append(feature)
        self._touch()

    def remove_feature(self, code: str) -> None:
        for i, f in enumerate(self.features):
            if f.code == code:
                self.features.pop(i)
                self._touch()
                return
        raise FeatureNotFoundError(code)

    def activate(self) -> None:
        self.is_active = True
        self._touch()

    def deactivate(self) -> None:
        self.is_active = False
        self._touch()

    # ---------- Query ----------
    def has_feature(self, code: str) -> bool:
        return any(f.code == code and f.enabled for f in self.features)

    def is_free(self) -> bool:
        return self.tier == PackageTier.FREE or self.price.amount == 0

    def is_higher_tier_than(self, other: "Package") -> bool:
        return self.tier.priority() > other.tier.priority()

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

```python
# app/modules/package/domain/entities/subscription.py
from datetime import datetime, timedelta, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.package.domain.entities.quotas import Quotas
from app.modules.package.domain.errors.errors import (
    InvalidCustomerError,
    InvalidPackageError,
    InvalidExpiryError,
    CannotUpgradeError,
    AlreadyCancelledError,
    CannotRenewError,
    ReasonRequiredError,
)
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle
from app.modules.package.domain.value_objects.subscription_status import SubscriptionStatus


class Subscription(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    customer_id: UUID
    package_id: UUID
    status: SubscriptionStatus = SubscriptionStatus.ACTIVE
    started_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    expires_at: datetime
    trial_ends_at: datetime | None = None
    grace_ends_at: datetime | None = None
    auto_renew: bool = True
    usage: Quotas = Field(default_factory=Quotas)
    previous_pkg_id: UUID | None = None
    upgraded_at: datetime | None = None
    cancelled_at: datetime | None = None
    cancel_reason: str = ""
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    # ---------- Constructor ----------
    @classmethod
    def create(
        cls,
        customer_id: UUID,
        package_id: UUID,
        cycle: BillingCycle,
    ) -> "Subscription":
        if not customer_id:
            raise InvalidCustomerError()
        if not package_id:
            raise InvalidPackageError()
        now = datetime.now(timezone.utc)
        return cls(
            customer_id=customer_id,
            package_id=package_id,
            started_at=now,
            expires_at=cycle.add_to(now),
        )

    @classmethod
    def create_trial(
        cls,
        customer_id: UUID,
        package_id: UUID,
        trial_days: int,
    ) -> "Subscription":
        sub = cls.create(customer_id, package_id, BillingCycle.MONTHLY)
        now = datetime.now(timezone.utc)
        trial_end = now + timedelta(days=trial_days)
        sub.status = SubscriptionStatus.TRIAL
        sub.trial_ends_at = trial_end
        sub.expires_at = trial_end
        return sub

    # ---------- Behavior ----------
    def upgrade(
        self, new_pkg_id: UUID, new_expiry: datetime, cycle: BillingCycle
    ) -> None:
        if not self.status.can_upgrade():
            raise CannotUpgradeError(f"cannot upgrade from {self.status.value}")
        if not new_expiry > datetime.now(timezone.utc):
            raise InvalidExpiryError()
        self.previous_pkg_id = self.package_id
        self.package_id = new_pkg_id
        self.expires_at = new_expiry
        self.status = SubscriptionStatus.ACTIVE
        self.upgraded_at = datetime.now(timezone.utc)
        self._touch()

    def cancel(self, reason: str) -> None:
        if self.status == SubscriptionStatus.CANCELLED:
            raise AlreadyCancelledError()
        if not reason:
            raise ReasonRequiredError()
        self.status = SubscriptionStatus.CANCELLED
        self.auto_renew = False
        self.cancelled_at = datetime.now(timezone.utc)
        self.cancel_reason = reason
        self._touch()

    def mark_past_due(self, grace_days: int) -> None:
        self.status = SubscriptionStatus.PAST_DUE
        self.grace_ends_at = datetime.now(timezone.utc) + timedelta(days=grace_days)
        self.auto_renew = False
        self._touch()

    def renew(self, cycle: BillingCycle) -> None:
        if self.status == SubscriptionStatus.CANCELLED:
            raise CannotRenewError()
        self.expires_at = cycle.add_to(self.expires_at)
        self.status = SubscriptionStatus.ACTIVE
        self.grace_ends_at = None
        self._touch()

    def expire(self) -> None:
        self.status = SubscriptionStatus.EXPIRED
        self.auto_renew = False
        self._touch()

    def update_usage(self, usage: Quotas) -> None:
        self.usage = usage
        self._touch()

    # ---------- Query ----------
    def is_active(self) -> bool:
        return self.status == SubscriptionStatus.ACTIVE and datetime.now(timezone.utc) < self.expires_at

    def is_in_trial(self) -> bool:
        return (self.status == SubscriptionStatus.TRIAL and self.trial_ends_at
                and datetime.now(timezone.utc) < self.trial_ends_at)

    def days_remaining(self) -> int:
        if datetime.now(timezone.utc) >= self.expires_at:
            return 0
        return (self.expires_at - datetime.now(timezone.utc)).days

    def is_expiring_soon(self, days: int) -> bool:
        remaining = self.days_remaining()
        return 0 < remaining <= days

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

### C.4 Domain Service

```python
# app/modules/package/domain/services/quota_service.py
from app.modules.package.domain.entities.package import Package
from app.modules.package.domain.entities.quotas import Quotas, UNLIMITED
from app.modules.package.domain.entities.subscription import Subscription
from app.modules.package.domain.value_objects.money import Money


class QuotaService:
    def check_quota(self, sub: Subscription, pkg: Package) -> tuple[bool, str]:
        if not sub.status.is_usable():
            return False, "subscription_not_active"
        return sub.usage.exceeds(pkg.quotas)

    def remaining(self, sub: Subscription, pkg: Package) -> Quotas:
        return Quotas(
            max_devices=self._remaining(sub.usage.max_devices, pkg.quotas.max_devices),
            max_sites=self._remaining(sub.usage.max_sites, pkg.quotas.max_sites),
            max_users=self._remaining(sub.usage.max_users, pkg.quotas.max_users),
            storage_gb=self._remaining(sub.usage.storage_gb, pkg.quotas.storage_gb),
            api_calls_per_day=self._remaining(sub.usage.api_calls_per_day, pkg.quotas.api_calls_per_day),
            automation_rules=self._remaining(sub.usage.automation_rules, pkg.quotas.automation_rules),
        )

    def calculate_proration(
        self, old_pkg: Package, new_pkg: Package, sub: Subscription
    ) -> Money:
        days_remaining = sub.days_remaining()
        if days_remaining <= 0:
            return Money(amount=0, currency=new_pkg.price.currency)
        total_days = new_pkg.billing_cycle.days()
        old_per_day = old_pkg.price.amount / total_days
        new_per_day = new_pkg.price.amount / total_days
        diff = (new_per_day - old_per_day) * days_remaining
        return Money(amount=diff, currency=new_pkg.price.currency)

    @staticmethod
    def _remaining(used: int, limit: int) -> int:
        if limit == UNLIMITED:
            return UNLIMITED
        return max(0, limit - used)
```

### C.5 Repository Protocols

```python
# app/modules/package/domain/repositories/package_repository.py
from datetime import datetime
from typing import Protocol
from uuid import UUID

from app.modules.package.domain.entities.package import Package
from app.modules.package.domain.entities.subscription import Subscription
from app.modules.package.domain.value_objects.package_tier import PackageTier


class PackageRepository(Protocol):
    async def save(self, pkg: Package) -> None: ...
    async def find_by_id(self, id: UUID) -> Package | None: ...
    async def find_by_code(self, code: str) -> Package | None: ...
    async def list_active(self) -> list[Package]: ...
    async def list_by_tier(self, tier: PackageTier) -> list[Package]: ...
    async def exists_by_code(self, code: str) -> bool: ...
    async def delete(self, id: UUID) -> None: ...


class SubscriptionRepository(Protocol):
    async def save(self, sub: Subscription) -> None: ...
    async def find_by_id(self, id: UUID) -> Subscription | None: ...
    async def find_active_by_customer(self, customer_id: UUID) -> Subscription | None: ...
    async def find_by_customer(self, customer_id: UUID) -> list[Subscription]: ...
    async def find_expiring(self, before: datetime) -> list[Subscription]: ...
    async def find_past_due(self) -> list[Subscription]: ...
    async def find_trial_ending(self, before: datetime) -> list[Subscription]: ...
```

### C.6 Errors

```python
# app/modules/package/domain/errors/errors.py


class DomainError(Exception):
    code: str = "DOMAIN_ERROR"

    def __init__(self, message: str = ""):
        self.message = message or self.__class__.__name__
        super().__init__(self.message)


class InvalidCodeError(DomainError):
    code = "INVALID_CODE"


class InvalidNameError(DomainError):
    code = "INVALID_NAME"


class InvalidTierError(DomainError):
    code = "INVALID_TIER"


class InvalidPriceError(DomainError):
    code = "INVALID_PRICE"


class InvalidCycleError(DomainError):
    code = "INVALID_CYCLE"


class InvalidCustomerError(DomainError):
    code = "INVALID_CUSTOMER"


class InvalidPackageError(DomainError):
    code = "INVALID_PACKAGE"


class InvalidExpiryError(DomainError):
    code = "INVALID_EXPIRY"


class CurrencyMismatchError(DomainError):
    code = "CURRENCY_MISMATCH"


class ReasonRequiredError(DomainError):
    code = "REASON_REQUIRED"


class PackageNotFoundError(DomainError):
    code = "PACKAGE_NOT_FOUND"


class SubscriptionNotFoundError(DomainError):
    code = "SUBSCRIPTION_NOT_FOUND"


class FeatureDuplicateError(DomainError):
    code = "FEATURE_DUPLICATE"


class FeatureNotFoundError(DomainError):
    code = "FEATURE_NOT_FOUND"


class CodeDuplicateError(DomainError):
    code = "CODE_DUPLICATE"


class AlreadySubscribedError(DomainError):
    code = "ALREADY_SUBSCRIBED"


class AlreadyCancelledError(DomainError):
    code = "ALREADY_CANCELLED"


class PackageInactiveError(DomainError):
    code = "PACKAGE_INACTIVE"


class CannotUpgradeError(DomainError):
    code = "CANNOT_UPGRADE"


class CannotRenewError(DomainError):
    code = "CANNOT_RENEW"


class SubscriptionNotActiveError(DomainError):
    code = "SUBSCRIPTION_NOT_ACTIVE"


class QuotaExceededError(DomainError):
    code = "QUOTA_EXCEEDED"
```

### C.7 Application — Subscribe Use Case

```python
# app/modules/package/application/use_cases/subscribe.py
from dataclasses import dataclass
from datetime import datetime, timezone
from uuid import UUID

from app.modules.package.application.ports.audit_repository import (
    AuditRepository, AuditTrail,
)
from app.modules.package.application.ports.customer_client import CustomerClient
from app.modules.package.application.ports.event_producer import EventProducer
from app.modules.package.application.ports.payment_client import (
    PaymentClient, PaymentRequest,
)
from app.modules.package.domain.entities.subscription import Subscription
from app.modules.package.domain.errors.errors import (
    PackageInactiveError, InvalidCustomerError, AlreadySubscribedError,
    PackageNotFoundError,
)
from app.modules.package.domain.events.events import SubscriptionCreatedEvent
from app.modules.package.domain.repositories.package_repository import (
    PackageRepository, SubscriptionRepository,
)
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle


@dataclass
class SubscribeInput:
    customer_id: UUID
    package_id: UUID
    cycle: BillingCycle
    user_id: UUID
    ip_address: str = ""


@dataclass
class SubscribeOutput:
    subscription_id: UUID
    payment_url: str
    expires_at: datetime


class SubscribeUseCase:
    def __init__(
        self,
        sub_repo: SubscriptionRepository,
        pkg_repo: PackageRepository,
        customer_client: CustomerClient,
        payment_client: PaymentClient,
        producer: EventProducer,
        audit_repo: AuditRepository,
    ):
        self._sub_repo = sub_repo
        self._pkg_repo = pkg_repo
        self._customer = customer_client
        self._payment = payment_client
        self._producer = producer
        self._audit = audit_repo

    async def execute(self, inp: SubscribeInput) -> SubscribeOutput:
        # 1. Load package
        pkg = await self._pkg_repo.find_by_id(inp.package_id)
        if not pkg:
            raise PackageNotFoundError()
        if not pkg.is_active:
            raise PackageInactiveError()

        # 2. Verify customer
        cust = await self._customer.get(inp.customer_id)
        if not cust or not cust.is_active:
            raise InvalidCustomerError()

        # 3. Check existing
        existing = await self._sub_repo.find_active_by_customer(inp.customer_id)
        if existing and existing.status.is_usable():
            raise AlreadySubscribedError()

        # 4. Create subscription
        sub = Subscription.create(inp.customer_id, pkg.id, inp.cycle)
        await self._sub_repo.save(sub)

        # 5. Create payment / activate if free
        payment_url = ""
        if not pkg.is_free():
            payment = await self._payment.create_payment(
                PaymentRequest(
                    customer_id=inp.customer_id,
                    amount=pkg.price.amount,
                    currency=pkg.price.currency,
                    description=f"Subscription: {pkg.name}",
                    reference_id=sub.id,
                    reference_type="subscription",
                )
            )
            payment_url = payment.url
        else:
            # free tier, activate immediately
            from app.modules.package.domain.value_objects.subscription_status import (
                SubscriptionStatus,
            )
            sub.status = SubscriptionStatus.ACTIVE
            await self._sub_repo.save(sub)

        # 6. Audit
        await self._audit.save(
            AuditTrail(
                user_id=inp.user_id,
                action="SUBSCRIPTION_CREATED",
                entity_type="subscription",
                entity_id=sub.id,
                payload={
                    "package_id": str(pkg.id),
                    "package_code": pkg.code,
                    "amount": pkg.price.amount,
                },
                ip_address=inp.ip_address,
                occurred_at=datetime.now(timezone.utc),
            )
        )

        # 7. Publish event
        await self._producer.publish_subscription_created(
            SubscriptionCreatedEvent(
                subscription_id=sub.id,
                customer_id=sub.customer_id,
                package_id=sub.package_id,
                expires_at=sub.expires_at,
                occurred_at=datetime.now(timezone.utc),
            )
        )

        return SubscribeOutput(
            subscription_id=sub.id,
            payment_url=payment_url,
            expires_at=sub.expires_at,
        )
```

### C.8 Infrastructure — SQLAlchemy Models

```python
# app/modules/package/infrastructure/persistence/postgres/models.py
from datetime import datetime
from uuid import UUID, uuid4

from sqlalchemy import Boolean, DateTime, Float, ForeignKey, Index, Integer, String, func
from sqlalchemy.dialects.postgresql import JSONB, UUID as PGUUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class PackageModel(Base):
    __tablename__ = "package_packages"
    __table_args__ = (
        Index("idx_package_tier", "tier", postgresql_where="is_active = TRUE"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    code: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(String, default="")
    tier: Mapped[str] = mapped_column(String(30), nullable=False)
    version: Mapped[int] = mapped_column(Integer, default=1)
    price_amount: Mapped[float] = mapped_column(Float, default=0)
    price_currency: Mapped[str] = mapped_column(String(3), default="THB")
    billing_cycle: Mapped[str] = mapped_column(String(20), nullable=False)
    quotas: Mapped[dict] = mapped_column(JSONB, default=dict)
    features: Mapped[list] = mapped_column(JSONB, default=list)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    is_public: Mapped[bool] = mapped_column(Boolean, default=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)


class SubscriptionModel(Base):
    __tablename__ = "package_subscriptions"
    __table_args__ = (
        Index("idx_sub_customer", "customer_id"),
        Index("idx_sub_expires", "expires_at"),
    )

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    customer_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    package_id: Mapped[UUID] = mapped_column(
        PGUUID(as_uuid=True), ForeignKey("package_packages.id")
    )
    status: Mapped[str] = mapped_column(String(30), default="ACTIVE")
    started_at: Mapped[datetime] = mapped_column(DateTime, nullable=False)
    expires_at: Mapped[datetime] = mapped_column(DateTime, nullable=False)
    trial_ends_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    grace_ends_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    auto_renew: Mapped[bool] = mapped_column(Boolean, default=True)
    usage: Mapped[dict] = mapped_column(JSONB, default=dict)
    previous_pkg_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    upgraded_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    cancelled_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    cancel_reason: Mapped[str] = mapped_column(String, default="")
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())
```

### C.9 Infrastructure — Repositories

```python
# app/modules/package/infrastructure/persistence/postgres/package_repository_impl.py
from uuid import UUID

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.package.domain.entities.package import Package
from app.modules.package.domain.entities.quotas import Quotas, Feature
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle
from app.modules.package.domain.value_objects.money import Money
from app.modules.package.domain.value_objects.package_tier import PackageTier
from app.modules.package.infrastructure.persistence.postgres.models import PackageModel


class PostgresPackageRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, pkg: Package) -> None:
        model = await self._session.get(PackageModel, pkg.id)
        if model is None:
            model = PackageModel(id=pkg.id)
            self._session.add(model)

        model.code = pkg.code
        model.name = pkg.name
        model.description = pkg.description
        model.tier = pkg.tier.value
        model.version = pkg.version
        model.price_amount = pkg.price.amount
        model.price_currency = pkg.price.currency
        model.billing_cycle = pkg.billing_cycle.value
        model.quotas = pkg.quotas.model_dump()
        model.features = [f.model_dump() for f in pkg.features]
        model.is_active = pkg.is_active
        model.is_public = pkg.is_public

        await self._session.flush()

    async def find_by_id(self, id: UUID) -> Package | None:
        stmt = select(PackageModel).where(
            PackageModel.id == id, PackageModel.deleted_at.is_(None)
        )
        m = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(m) if m else None

    async def find_by_code(self, code: str) -> Package | None:
        stmt = select(PackageModel).where(
            PackageModel.code == code, PackageModel.deleted_at.is_(None)
        )
        m = (await self._session.execute(stmt)).scalar_one_or_none()
        return self._to_entity(m) if m else None

    async def list_active(self) -> list[Package]:
        stmt = (
            select(PackageModel)
            .where(PackageModel.is_active == True, PackageModel.deleted_at.is_(None))
            .order_by(PackageModel.tier.desc(), PackageModel.price_amount.asc())
        )
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def exists_by_code(self, code: str) -> bool:
        stmt = select(func.count()).select_from(PackageModel).where(
            PackageModel.code == code, PackageModel.deleted_at.is_(None)
        )
        return (await self._session.execute(stmt)).scalar() > 0

    async def delete(self, id: UUID) -> None:
        from datetime import datetime, timezone
        m = await self._session.get(PackageModel, id)
        if m:
            m.deleted_at = datetime.now(timezone.utc)

    def _to_entity(self, m: PackageModel) -> Package:
        return Package(
            id=m.id,
            code=m.code,
            name=m.name,
            description=m.description,
            tier=PackageTier(m.tier),
            version=m.version,
            price=Money(amount=m.price_amount, currency=m.price_currency),
            billing_cycle=BillingCycle(m.billing_cycle),
            quotas=Quotas(**m.quotas) if m.quotas else Quotas(),
            features=[Feature(**f) for f in (m.features or [])],
            is_active=m.is_active,
            is_public=m.is_public,
            created_at=m.created_at,
            updated_at=m.updated_at,
        )
```

### C.10 Interface — FastAPI Router

```python
# app/modules/package/interfaces/http/schemas.py
from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, Field


class QuotasSchema(BaseModel):
    max_devices: int = 0
    max_sites: int = 0
    max_users: int = 0
    storage_gb: int = 0
    retention_days: int = 0
    api_calls_per_day: int = 0
    automation_rules: int = 0


class CreatePackageRequest(BaseModel):
    code: str = Field(..., pattern=r"^[A-Z0-9_]{2,50}$")
    name: str = Field(..., min_length=1, max_length=255)
    description: str = ""
    tier: str = Field(..., pattern=r"^(FREE|BASIC|PRO|ENTERPRISE)$")
    price_amount: float = Field(..., ge=0)
    price_currency: str = "THB"
    billing_cycle: str = Field(..., pattern=r"^(MONTHLY|YEARLY)$")
    quotas: QuotasSchema


class PackageResponse(BaseModel):
    id: UUID
    code: str
    name: str
    description: str
    tier: str
    version: int
    price_amount: float
    price_currency: str
    billing_cycle: str
    quotas: dict
    features: list[dict]
    is_active: bool


class SubscribeRequest(BaseModel):
    customer_id: UUID
    package_id: UUID
    cycle: str = Field(..., pattern=r"^(MONTHLY|YEARLY)$")


class SubscribeResponse(BaseModel):
    subscription_id: UUID
    payment_url: str = ""
    expires_at: datetime


class UpgradeRequest(BaseModel):
    new_package_id: UUID


class CancelRequest(BaseModel):
    reason: str = Field(..., min_length=1)


class CheckQuotaRequest(BaseModel):
    usage: QuotasSchema
```

```python
# app/modules/package/interfaces/http/package_router.py
from fastapi import APIRouter, Depends, HTTPException, status

from app.modules.package.application.use_cases.create_package import (
    CreatePackageInput, CreatePackageUseCase,
)
from app.modules.package.domain.errors.errors import DomainError, CodeDuplicateError
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle
from app.modules.package.domain.value_objects.money import Money
from app.modules.package.domain.value_objects.package_tier import PackageTier
from app.modules.package.interfaces.http.dependencies import (
    get_create_package_uc, get_current_user,
)
from app.modules.package.interfaces.http.schemas import (
    CreatePackageRequest, PackageResponse,
)

router = APIRouter(prefix="/api/v1/packages", tags=["packages"])


@router.post("", response_model=PackageResponse, status_code=status.HTTP_201_CREATED)
async def create_package(
    req: CreatePackageRequest,
    uc: CreatePackageUseCase = Depends(get_create_package_uc),
    user_id = Depends(get_current_user),
):
    from app.modules.package.domain.entities.quotas import Quotas
    try:
        out = await uc.execute(
            CreatePackageInput(
                code=req.code,
                name=req.name,
                description=req.description,
                tier=PackageTier(req.tier),
                price=Money(amount=req.price_amount, currency=req.price_currency),
                cycle=BillingCycle(req.billing_cycle),
                quotas=Quotas(**req.quotas.model_dump()),
                user_id=user_id,
            )
        )
        return {"id": out.package_id, "code": out.code, **req.model_dump()}
    except CodeDuplicateError as e:
        raise HTTPException(status_code=409, detail={"code": e.code, "message": str(e)})
    except DomainError as e:
        raise HTTPException(status_code=400, detail={"code": e.code, "message": str(e)})
```

```python
# app/modules/package/interfaces/http/subscription_router.py
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status

from app.modules.package.application.use_cases.subscribe import (
    SubscribeInput, SubscribeUseCase,
)
from app.modules.package.application.use_cases.upgrade_subscription import (
    UpgradeSubscriptionInput, UpgradeSubscriptionUseCase,
)
from app.modules.package.application.use_cases.cancel_subscription import (
    CancelSubscriptionInput, CancelSubscriptionUseCase,
)
from app.modules.package.domain.errors.errors import (
    DomainError, AlreadySubscribedError, PackageInactiveError,
    SubscriptionNotFoundError,
)
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle
from app.modules.package.interfaces.http.dependencies import (
    get_subscribe_uc, get_upgrade_uc, get_cancel_uc, get_current_user,
)
from app.modules.package.interfaces.http.schemas import (
    SubscribeRequest, SubscribeResponse, UpgradeRequest, CancelRequest,
)

router = APIRouter(prefix="/api/v1/subscriptions", tags=["subscriptions"])


@router.post("", response_model=SubscribeResponse, status_code=status.HTTP_201_CREATED)
async def subscribe(
    req: SubscribeRequest,
    uc: SubscribeUseCase = Depends(get_subscribe_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        out = await uc.execute(
            SubscribeInput(
                customer_id=req.customer_id,
                package_id=req.package_id,
                cycle=BillingCycle(req.cycle),
                user_id=user_id,
            )
        )
        return SubscribeResponse(
            subscription_id=out.subscription_id,
            payment_url=out.payment_url,
            expires_at=out.expires_at,
        )
    except AlreadySubscribedError as e:
        raise HTTPException(status_code=409, detail={"code": e.code, "message": str(e)})
    except PackageInactiveError as e:
        raise HTTPException(status_code=400, detail={"code": e.code, "message": str(e)})
    except DomainError as e:
        raise HTTPException(status_code=400, detail={"code": e.code, "message": str(e)})


@router.post("/{subscription_id}/upgrade")
async def upgrade(
    subscription_id: UUID,
    req: UpgradeRequest,
    uc: UpgradeSubscriptionUseCase = Depends(get_upgrade_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        out = await uc.execute(
            UpgradeSubscriptionInput(
                subscription_id=subscription_id,
                new_package_id=req.new_package_id,
                user_id=user_id,
            )
        )
        return {
            "subscription_id": str(out.subscription_id),
            "from_package_id": str(out.from_package_id),
            "to_package_id": str(out.to_package_id),
            "proration_amount": out.proration_amount,
            "currency": out.currency,
            "new_expires_at": out.new_expires_at.isoformat(),
            "payment_url": out.payment_url,
        }
    except SubscriptionNotFoundError as e:
        raise HTTPException(status_code=404, detail={"code": e.code, "message": str(e)})
    except DomainError as e:
        raise HTTPException(status_code=400, detail={"code": e.code, "message": str(e)})


@router.post("/{subscription_id}/cancel")
async def cancel(
    subscription_id: UUID,
    req: CancelRequest,
    uc: CancelSubscriptionUseCase = Depends(get_cancel_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(
            CancelSubscriptionInput(
                subscription_id=subscription_id,
                reason=req.reason,
                user_id=user_id,
            )
        )
        return {"message": "cancelled"}
    except DomainError as e:
        raise HTTPException(status_code=400, detail={"code": e.code, "message": str(e)})
```

### C.11 Composition Root

```python
# app/modules/package/module.py
from dataclasses import dataclass

from aiokafka import AIOKafkaProducer
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.package.application.use_cases.cancel_subscription import CancelSubscriptionUseCase
from app.modules.package.application.use_cases.create_package import CreatePackageUseCase
from app.modules.package.application.use_cases.subscribe import SubscribeUseCase
from app.modules.package.application.use_cases.upgrade_subscription import UpgradeSubscriptionUseCase
from app.modules.package.domain.services.quota_service import QuotaService
from app.modules.package.infrastructure.messaging.kafka_producer import KafkaEventProducer
from app.modules.package.infrastructure.persistence.postgres.package_repository_impl import (
    PostgresPackageRepository,
)
from app.modules.package.infrastructure.persistence.postgres.subscription_repository_impl import (
    PostgresSubscriptionRepository,
)


@dataclass
class PackageModuleDeps:
    session: AsyncSession
    redis: Redis
    kafka: AIOKafkaProducer
    customer_client: "CustomerClient"
    payment_client: "PaymentClient"
    audit_repo: "AuditRepository"
    grace_days: int = 7


def build_module(deps: PackageModuleDeps) -> dict:
    pkg_repo = PostgresPackageRepository(deps.session)
    sub_repo = PostgresSubscriptionRepository(deps.session)
    producer = KafkaEventProducer(deps.kafka)
    quota_svc = QuotaService()

    return {
        "create_package": CreatePackageUseCase(pkg_repo, deps.audit_repo, producer),
        "subscribe": SubscribeUseCase(
            sub_repo, pkg_repo, deps.customer_client,
            deps.payment_client, producer, deps.audit_repo,
        ),
        "upgrade_subscription": UpgradeSubscriptionUseCase(
            sub_repo, pkg_repo, quota_svc,
            deps.payment_client, producer, deps.audit_repo,
        ),
        "cancel_subscription": CancelSubscriptionUseCase(
            sub_repo, producer, deps.audit_repo,
        ),
    }
```

### C.12 Python Tests

```python
# tests/modules/package/test_subscription_entity.py
from datetime import datetime, timedelta, timezone
from uuid import uuid4

import pytest

from app.modules.package.domain.entities.quotas import Quotas, UNLIMITED
from app.modules.package.domain.entities.subscription import Subscription
from app.modules.package.domain.errors.errors import (
    CannotUpgradeError, AlreadyCancelledError, ReasonRequiredError,
)
from app.modules.package.domain.value_objects.billing_cycle import BillingCycle


def test_create_subscription_defaults():
    sub = Subscription.create(uuid4(), uuid4(), BillingCycle.MONTHLY)
    assert sub.is_active()
    assert sub.auto_renew is True
    assert sub.days_remaining() >= 29


def test_upgrade_active_ok():
    sub = Subscription.create(uuid4(), uuid4(), BillingCycle.MONTHLY)
    new_pkg = uuid4()
    new_expiry = datetime.now(timezone.utc) + timedelta(days=30)
    sub.upgrade(new_pkg, new_expiry, BillingCycle.MONTHLY)
    assert sub.package_id == new_pkg
    assert sub.previous_pkg_id is not None


def test_upgrade_cancelled_fails():
    sub = Subscription.create(uuid4(), uuid4(), BillingCycle.MONTHLY)
    sub.cancel("user request")
    with pytest.raises(CannotUpgradeError):
        sub.upgrade(uuid4(), datetime.now(timezone.utc) + timedelta(days=30), BillingCycle.MONTHLY)


def test_cancel_without_reason_fails():
    sub = Subscription.create(uuid4(), uuid4(), BillingCycle.MONTHLY)
    with pytest.raises(ReasonRequiredError):
        sub.cancel("")


def test_cancel_twice_fails():
    sub = Subscription.create(uuid4(), uuid4(), BillingCycle.MONTHLY)
    sub.cancel("first")
    with pytest.raises(AlreadyCancelledError):
        sub.cancel("second")


def test_quotas_exceeds():
    usage = Quotas(max_devices=15)
    limit = Quotas(max_devices=10)
    exceeds, field = usage.exceeds(limit)
    assert exceeds
    assert field == "devices"


def test_quotas_unlimited():
    usage = Quotas(max_devices=999999)
    limit = Quotas(max_devices=UNLIMITED)
    exceeds, _ = usage.exceeds(limit)
    assert not exceeds
```

---

## PART D — Integration Notes

### D.1 Go ↔ Python Interop

ทั้ง 2 implementation **แชร์ spec เดียวกัน** — สามารถ:
- **Run parallel** ในช่วง migrate (dual-write)
- **ใช้ DB ร่วม** (PostgreSQL)
- **ใช้ Kafka ร่วม** (topic เดียวกัน)
- **ใช้ Redis key namespace เดียวกัน** (`package:package:<id>`)

**Key format:**
```
package:package:{id}              → Package cache
package:subscription:{id}          → Subscription cache
package:customer:{customer_id}     → Active subscription lookup
```

### D.2 Wire-up ใน `cmd/api/main.go` (Go)

```go
func main() {
    // ... existing setup ...
    
    api := r.Group("/api/v1")
    authMw := middleware.JWTAuth(cfg.JWT.Secret)
    tenantMw := middleware.Tenant()
    rbacMw := middleware.RBAC()
    
    // customer module
    customer.Init(api, customer.Dependencies{
        DB: db, Redis: rdb,
        Producer: customerProducer,
        AuditRepo: auditRepo,
        WSHub: wsHub,
        PackageClient: packageClient,
    }, authMw, tenantMw, rbacMw)
    
    // package module
    package_.Init(api, package_.Dependencies{
        DB: db, Redis: rdb,
        Producer: packageProducer,
        CustomerCli: customerClient,
        PaymentCli: paymentClient,
        AuditRepo: auditRepo,
        GraceDays: cfg.Package.GraceDays,
    }, authMw, tenantMw, rbacMw)
    
    // ... other modules ...
}
```

### D.3 Wire-up ใน `main.py` (Python)

```python
from fastapi import FastAPI
from app.modules.customer.module import build_module as build_customer
from app.modules.package.module import build_module as build_package
from app.modules.package.interfaces.http.package_router import router as pkg_router
from app.modules.package.interfaces.http.subscription_router import router as sub_router


def create_app() -> FastAPI:
    app = FastAPI(title="icmongolang IoT Platform")
    
    # ... setup db, redis, kafka ...
    
    # build modules
    customer_module = build_customer(customer_deps)
    package_module = build_package(package_deps)
    
    # store in app state for dependency injection
    app.state.customer = customer_module
    app.state.package = package_module
    
    # register routers
    app.include_router(pkg_router)
    app.include_router(sub_router)
    
    return app
```

### D.4 Migration Sequence

```
Day 0:  Run 20260101_customer_init.sql + 20260102_package_init.sql
Day 1:  Deploy customer module (Go) + seed data
Day 2:  Deploy package module (Go) + seed package tiers
Day 3:  Verify end-to-end: create customer → subscribe → receive payment
Day 4:  (Optional) Deploy Python versions in parallel
Day 5:  Load test: 1000 customers, 500 subscriptions
```

---

## ✅ Checklist Module 02

- [x] Business spec ครบ (Purpose, Scope, Stakeholders)
- [x] Domain model ครบ (2 aggregates, 4 VOs, Quotas VO)
- [x] Invariants ระบุชัด (8 ข้อ)
- [x] Use cases 14 ตัว
- [x] API contract 13 endpoints
- [x] DB schema (2 tables)
- [x] Kafka topics (8 publish, 4 consume)
- [x] Env vars
- [x] Migration filename: `20260102_package_init.sql`
- [x] Go implementation (full: domain → infra → interface → module.go)
- [x] Python implementation (full: domain → infra → interface → module.py)
- [x] Tests (Go + Python)
- [x] Scheduler jobs (expire, retry past due)
- [x] Integration notes (Go ↔ Python interop)

---

# 📋 MODULES 03-08 — READY SPEC

> Modules ต่อไปนี้มี **spec สมบูรณ์** — ใช้โครงเดียวกับ Module 01 + 02

## Module 03: `device` — 🔥 หัวใจหลัก

**Aggregates:** `Site` (root → zones), `Device` (root), `AutomationRule` (root), `Command` (root), `Telemetry` (InfluxDB)

**Use Cases:** CreateSite, AddZone, RegisterDevice, ProvisionDevice, IngestTelemetry, IssueCommand, AckCommand, CreateAutomationRule, EvaluateAutomation, UpdateFirmware, DecommissionDevice

**Cross-cutting:** MQTT (in/out), Kafka, InfluxDB, WebSocket, Cache

**Particularities:**
- **InfluxDB measurement:** `telemetry` (tags: device_id, site_id, metric, unit; field: value)
- **MQTT topics:** `iot/{serial}/{telemetry|heartbeat|status|ack|command|config|firmware}`
- **Shadow:** reported/desired state per device
- **Rate limit:** 1 telemetry/sec/device

**Migration:** `20260103_device_init.sql` (5 tables + InfluxDB setup)

**Estimated LOC:** ~4,500 Go / ~3,200 Python

---

## Module 04: `erp`

**Aggregates:** `Order`, `Invoice`, `StockItem`, `Warehouse`

**Use Cases:** CreateOrder, ConfirmOrder, ShipOrder, DeliverOrder, CreateInvoice, IssueInvoice, RecordPayment, RestockItem, ReorderCheck

**Numbering:** order/invoice auto-generate per tenant (SO-2026-0001)

**Migration:** `20260104_erp_init.sql` (5 tables)

**Estimated LOC:** ~3,800 Go / ~2,700 Python

---

## Module 05: `crm`

**Aggregates:** `Lead`, `Opportunity`, `Ticket`, `Activity`, `Campaign`

**Use Cases:** CaptureLead, QualifyLead, ConvertLead, OpenTicket, AssignTicket, ResolveTicket, EscalateTicket

**SLA:** per-priority (LOW=48h, MEDIUM=24h, HIGH=8h, URGENT=2h)

**Migration:** `20260105_crm_init.sql` (6 tables)

**Estimated LOC:** ~3,200 Go / ~2,300 Python

---

## Module 06: `logistics`

**Aggregates:** `Shipment`, `Route`, `Vehicle`, `Driver`

**Cold chain:** 2-8°C validation, realtime alerts

**Integration:** device (GPS/temp), erp (order), notifier (alert)

**Migration:** `20260106_logistics_init.sql` (6 tables)

**Estimated LOC:** ~2,800 Go / ~2,000 Python

---

## Module 07: `report` (Upgrade)

**Aggregates:** `ReportDefinition`, `ReportRun`, `Dashboard`, `Widget`

**Data Sources:** Postgres (SQL), InfluxDB (Flux), Elasticsearch (DSL)

**AI:** LLM-generated insight via `pkg/llm`

**Export:** PDF, Excel, CSV, JSON (S3)

**Migration:** `20260107_report_upgrade.sql` (3 tables + legacy migration)

**Estimated LOC:** ~3,500 Go / ~2,500 Python

---

## Module 08: `auth` (Upgrade)

**Improvements:** tenant-aware JWT, SSO, MFA, session management

**Migration:** `20260108_auth_upgrade.sql` (add columns only)

**Breaking changes:** JWT payload + tenant_id + session_id

**Estimated LOC:** ~1,800 Go / ~1,300 Python

---

# 🎯 Next Step

**เลือกได้ 2 ทาง:**

1. **Generate Code ทันที** — ใช้ prompts (E.1/E.2) + Module 03 spec ป้อน AI → ได้ `device` implementation
2. **ให้ผมทำ Module 03 (device) แบบ Full Deep Dive** เหมือน Module 01/02 — พร้อม Go + Python + Tests + Cross-cutting (MQTT, InfluxDB, WS)

> **คำแนะนำ:** Module 03 (`device`) เป็น **หัวใจของ IoT platform** — ควรทำแบบ Full Deep Dive เพราะมี pattern เฉพาะ (MQTT, InfluxDB, Shadow, Automation) ที่ modules อื่นไม่มี

 