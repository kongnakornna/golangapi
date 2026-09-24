# ระบบจัดการสิทธิ์ส่วนบุคคล (PDPA Module) - ฉบับสมบูรณ์

## 📦 โครงสร้างโปรเจกต์

```
icmongolang/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── migrate/
│   │   └── main.go
│   ├── scheduler/
│   │   └── main.go
│   └── workers/
│       ├── dsar/
│       │   └── main.go
│       ├── email/
│       │   └── main.go
│       ├── llm/
│       │   └── main.go
│       ├── revoke/
│       │   └── main.go
│       └── account/
│           └── main.go
├── internal/
│   └── modules/
│       └── pdpa/
│           ├── domain/
│           │   ├── entity/
│           │   │   ├── consent_log.go
│           │   │   ├── dsar_request.go
│           │   │   ├── audit_trail.go
│           │   │   └── user_account_status.go
│           │   ├── value_object/
│           │   │   ├── consent_purpose.go
│           │   │   ├── consent_status.go
│           │   │   ├── dsar_type.go
│           │   │   ├── dsar_status.go
│           │   │   └── account_status.go
│           │   ├── repository/
│           │   │   ├── consent_repository.go
│           │   │   ├── dsar_repository.go
│           │   │   ├── audit_repository.go
│           │   │   └── user_account_status_repository.go
│           │   ├── service/
│           │   │   ├── consent_validator.go
│           │   │   ├── anonymization_service.go
│           │   │   ├── localization_service.go
│           │   │   └── deletion_policy_service.go
│           │   └── errors/
│           │       └── errors.go
│           ├── application/
│           │   ├── record_consent.go
│           │   ├── revoke_consent.go
│           │   ├── immediate_deletion.go
│           │   ├── auto_delete_expired_consents.go
│           │   ├── submit_dsar.go
│           │   ├── process_dsar.go
│           │   ├── get_consent_history.go
│           │   ├── get_dsar_status.go
│           │   ├── get_admin_report.go
│           │   ├── handle_account_event.go
│           │   └── dto.go
│           ├── infrastructure/
│           │   ├── persistence/
│           │   │   ├── postgres/
│           │   │   │   ├── consent_repo_impl.go
│           │   │   │   ├── dsar_repo_impl.go
│           │   │   │   ├── audit_repo_impl.go
│           │   │   │   ├── user_account_status_repo_impl.go
│           │   │   │   └── models.go
│           │   │   └── redis/
│           │   │       └── consent_cache.go
│           │   ├── messaging/
│           │   │   ├── kafka_producer.go
│           │   │   └── consumers/
│           │   │       ├── dsar_worker.go
│           │   │       ├── llm_worker.go
│           │   │       ├── email_worker.go
│           │   │       ├── blockchain_worker.go
│           │   │       ├── revoke_consumer.go
│           │   │       └── account_event_consumer.go
│           │   ├── search/
│           │   │   └── elasticsearch/
│           │   │       └── consent_indexer.go
│           │   └── scheduler/
│           │       └── consent_cleanup_job.go
│           └── interfaces/
│               ├── http/
│               │   ├── consent_handler.go
│               │   ├── dsar_handler.go
│               │   ├── admin_handler.go
│               │   ├── deletion_handler.go
│               │   ├── routes.go
│               │   └── dto.go
│               ├── websocket/
│               │   └── dsar_broadcaster.go
│               └── localization/
│                   ├── locales/
│                   │   ├── en.toml
│                   │   └── th.toml
│                   └── i18n.go
├── migrations/
│   └── 001_initial_pdpa_schema.sql
├── go.mod
├── go.sum
├── .env
└── docker-compose.yml
```

---

## 🗄️ Database Migrations (พร้อม Prefix `pdpa_`)

#### `migrations/001_initial_pdpa_schema.sql`

```sql
-- ============================================================
-- ระบบ PDPA (Personal Data Protection Act)
-- Prefix: pdpa_ เพื่อแยกจากโมดูลอื่น
-- ============================================================

-- 1. ตารางเก็บเวอร์ชันนโยบายความเป็นส่วนตัว
CREATE TABLE IF NOT EXISTS pdpa_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(20) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    effective_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_policies_version ON pdpa_policies (version);
CREATE INDEX IF NOT EXISTS idx_pdpa_policies_active ON pdpa_policies (is_active);

-- 2. ตารางเก็บวัตถุประสงค์ในการขอข้อมูล
CREATE TABLE IF NOT EXISTS pdpa_purposes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_purposes_code ON pdpa_purposes (code);

-- ข้อมูลเริ่มต้น
INSERT INTO pdpa_purposes (id, code, name, description, is_required) VALUES
    (gen_random_uuid(), 'NECESSARY', 'ข้อมูลจำเป็นเพื่อการให้บริการ', 'จำเป็นต่อการดำเนินงานตามสัญญา', TRUE),
    (gen_random_uuid(), 'ANALYTICS', 'การวิเคราะห์ข้อมูล', 'เพื่อพัฒนาปรับปรุงบริการ', FALSE),
    (gen_random_uuid(), 'MARKETING', 'การตลาดและการโฆษณา', 'เพื่อนำเสนอสินค้าและบริการ', FALSE)
ON CONFLICT (code) DO NOTHING;

-- 3. ตารางบันทึกการให้ความยินยอม (Consent Log)
CREATE TABLE IF NOT EXISTS pdpa_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    session_id VARCHAR(255),
    purpose_code VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'GRANTED',
    ip_address VARCHAR(45),
    user_agent TEXT,
    granted_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP,
    auto_deleted_at TIMESTAMP,
    CONSTRAINT fk_pdpa_consents_purpose 
        FOREIGN KEY (purpose_code) REFERENCES pdpa_purposes(code) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_pdpa_consents_user_id ON pdpa_consents (user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_purpose ON pdpa_consents (purpose_code);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_status ON pdpa_consents (status);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_granted_at ON pdpa_consents (granted_at);
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_revoked_at ON pdpa_consents (revoked_at);

-- 4. ตารางคำร้องขอใช้สิทธิ์ (DSAR)
CREATE TABLE IF NOT EXISTS pdpa_user_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    request_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    requested_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    data_payload JSONB,
    rejection_reason TEXT,
    otp_code VARCHAR(255),
    otp_expired_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX IF NOT EXISTS idx_pdpa_user_requests_user_id ON pdpa_user_requests (user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_user_requests_status ON pdpa_user_requests (status);
CREATE INDEX IF NOT EXISTS idx_pdpa_user_requests_requested_at ON pdpa_user_requests (requested_at);

-- 5. ตารางสถานะบัญชีผู้ใช้
CREATE TABLE IF NOT EXISTS pdpa_user_account_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    suspended_at TIMESTAMP,
    terminated_at TIMESTAMP,
    deletion_confirmed_at TIMESTAMP,
    retention_deadline TIMESTAMP,
    auto_deleted_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_status_user_id ON pdpa_user_account_statuses (user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_status_retention ON pdpa_user_account_statuses (retention_deadline) 
    WHERE status = 'SUSPENDED';

-- 6. ตารางบันทึกการดำเนินการ (Audit Trail)
CREATE TABLE IF NOT EXISTS pdpa_audit_trails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    action VARCHAR(50) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pdpa_audit_user_id ON pdpa_audit_trails (user_id);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_action ON pdpa_audit_trails (action);
CREATE INDEX IF NOT EXISTS idx_pdpa_audit_created_at ON pdpa_audit_trails (created_at);

-- 7. ตารางเก็บการตอบกลับคำร้อง (ไฟล์)
CREATE TABLE IF NOT EXISTS pdpa_request_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL,
    file_name VARCHAR(255),
    file_path VARCHAR(500),
    file_size BIGINT,
    mime_type VARCHAR(100),
    data_content TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    CONSTRAINT fk_pdpa_response_request 
        FOREIGN KEY (request_id) REFERENCES pdpa_user_requests(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pdpa_responses_request_id ON pdpa_request_responses (request_id);
```

---

## 🏛️ Domain Layer

### 1. Entities

#### `domain/entity/consent_log.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/value_object"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
)

type ConsentLog struct {
	ID            uuid.UUID                   `json:"id"`
	UserID        uuid.UUID                   `json:"user_id"`
	SessionID     string                      `json:"session_id"`
	Purpose       valueobject.ConsentPurpose  `json:"purpose"`
	Status        valueobject.ConsentStatus   `json:"status"`
	IPAddress     string                      `json:"ip_address"`
	UserAgent     string                      `json:"user_agent"`
	GrantedAt     time.Time                   `json:"granted_at"`
	ExpiresAt     time.Time                   `json:"expires_at"`
	RevokedAt     *time.Time                  `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time                  `json:"auto_deleted_at,omitempty"`
}

func NewConsentLog(userID uuid.UUID, sessionID string, purpose valueobject.ConsentPurpose, ip, userAgent string) *ConsentLog {
	now := time.Now()
	return &ConsentLog{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
		Purpose:   purpose,
		Status:    valueobject.ConsentGranted,
		IPAddress: ip,
		UserAgent: userAgent,
		GrantedAt: now,
		ExpiresAt: now.AddDate(1, 0, 0), // 1 year
	}
}

func (c *ConsentLog) Revoke() error {
	if c.Status == valueobject.ConsentRevoked {
		return domainerrors.ErrConsentAlreadyRevoked
	}
	if c.Status == valueobject.ConsentExpired {
		return domainerrors.ErrConsentExpired
	}
	now := time.Now()
	c.Status = valueobject.ConsentRevoked
	c.RevokedAt = &now
	return nil
}

func (c *ConsentLog) MarkDeleted() {
	now := time.Now()
	c.Status = valueobject.ConsentDeleted
	c.AutoDeletedAt = &now
}

func (c *ConsentLog) IsActive() bool {
	return c.Status == valueobject.ConsentGranted && c.ExpiresAt.After(time.Now())
}
```

#### `domain/entity/dsar_request.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/value_object"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
)

type DSARRequest struct {
	ID              uuid.UUID                 `json:"id"`
	UserID          uuid.UUID                 `json:"user_id"`
	RequestType     valueobject.DSARType      `json:"request_type"`
	Status          valueobject.DSARStatus    `json:"status"`
	RequestedAt     time.Time                 `json:"requested_at"`
	CompletedAt     *time.Time                `json:"completed_at,omitempty"`
	DataPayload     []byte                    `json:"data_payload,omitempty"`
	RejectionReason string                    `json:"rejection_reason,omitempty"`
	OTPCode         string                    `json:"otp_code,omitempty"`
	OTPExpiredAt    time.Time                 `json:"otp_expired_at,omitempty"`
	IPAddress       string                    `json:"ip_address,omitempty"`
	UserAgent       string                    `json:"user_agent,omitempty"`
}

func NewDSARRequest(userID uuid.UUID, requestType valueobject.DSARType, otpCode string, otpExpiry time.Time, ip, userAgent string) *DSARRequest {
	return &DSARRequest{
		ID:           uuid.New(),
		UserID:       userID,
		RequestType:  requestType,
		Status:       valueobject.DSARStatusPending,
		RequestedAt:  time.Now(),
		OTPCode:      otpCode,
		OTPExpiredAt: otpExpiry,
		IPAddress:    ip,
		UserAgent:    userAgent,
	}
}

func (d *DSARRequest) VerifyOTP(otp string) bool {
	if d.OTPCode != otp {
		return false
	}
	if time.Now().After(d.OTPExpiredAt) {
		return false
	}
	return true
}

func (d *DSARRequest) MarkProcessing() error {
	if d.Status != valueobject.DSARStatusPending {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusProcessing
	return nil
}

func (d *DSARRequest) MarkCompleted(payload []byte) error {
	if d.Status == valueobject.DSARStatusCompleted {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	now := time.Now()
	d.Status = valueobject.DSARStatusCompleted
	d.CompletedAt = &now
	d.DataPayload = payload
	return nil
}

func (d *DSARRequest) MarkRejected(reason string) error {
	if d.Status == valueobject.DSARStatusCompleted {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusRejected
	d.RejectionReason = reason
	return nil
}
```

#### `domain/entity/audit_trail.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditTrail struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	Action    string      `json:"action"`
	Details   interface{} `json:"details,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewAuditTrail(userID *uuid.UUID, action string, details interface{}) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}
}

func NewAuditTrailWithMeta(userID *uuid.UUID, action string, details interface{}, ip, userAgent string) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
}
```

#### `domain/entity/user_account_status.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type UserAccountStatus struct {
	ID                   uuid.UUID                    `json:"id"`
	UserID               uuid.UUID                    `json:"user_id"`
	Status               valueobject.AccountStatus    `json:"status"`
	SuspendedAt          *time.Time                   `json:"suspended_at,omitempty"`
	TerminatedAt         *time.Time                   `json:"terminated_at,omitempty"`
	DeletionConfirmedAt  *time.Time                   `json:"deletion_confirmed_at,omitempty"`
	RetentionDeadline    *time.Time                   `json:"retention_deadline,omitempty"`
	AutoDeletedAt        *time.Time                   `json:"auto_deleted_at,omitempty"`
	UpdatedAt            time.Time                    `json:"updated_at"`
}

func NewUserAccountStatus(userID uuid.UUID, status valueobject.AccountStatus) *UserAccountStatus {
	return &UserAccountStatus{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    status,
		UpdatedAt: time.Now(),
	}
}

func (u *UserAccountStatus) SetSuspended(suspendedAt time.Time, retentionYears int) {
	now := time.Now()
	u.Status = valueobject.AccountSuspended
	u.SuspendedAt = &suspendedAt
	deadline := suspendedAt.AddDate(retentionYears, 0, 0)
	u.RetentionDeadline = &deadline
	u.UpdatedAt = now
}

func (u *UserAccountStatus) SetTerminated(terminatedAt time.Time) {
	now := time.Now()
	u.Status = valueobject.AccountTerminated
	u.TerminatedAt = &terminatedAt
	u.UpdatedAt = now
}

func (u *UserAccountStatus) ConfirmDeletion() {
	now := time.Now()
	u.DeletionConfirmedAt = &now
	u.UpdatedAt = now
}

func (u *UserAccountStatus) MarkDeleted() {
	now := time.Now()
	u.Status = "DELETED"
	u.AutoDeletedAt = &now
	u.UpdatedAt = now
}

func (u *UserAccountStatus) IsActive() bool {
	return u.Status == valueobject.AccountActive
}

func (u *UserAccountStatus) IsSuspended() bool {
	return u.Status == valueobject.AccountSuspended
}

func (u *UserAccountStatus) IsTerminated() bool {
	return u.Status == valueobject.AccountTerminated
}

func (u *UserAccountStatus) IsDeleted() bool {
	return u.Status == "DELETED"
}
```

---

### 2. Value Objects

#### `domain/value_object/consent_purpose.go`

```go
package valueobject

type ConsentPurpose string

const (
	PurposeNecessary ConsentPurpose = "NECESSARY"
	PurposeAnalytics ConsentPurpose = "ANALYTICS"
	PurposeMarketing ConsentPurpose = "MARKETING"
)

func (c ConsentPurpose) IsValid() bool {
	switch c {
	case PurposeNecessary, PurposeAnalytics, PurposeMarketing:
		return true
	}
	return false
}

func (c ConsentPurpose) String() string {
	return string(c)
}
```

#### `domain/value_object/consent_status.go`

```go
package valueobject

type ConsentStatus string

const (
	ConsentGranted ConsentStatus = "GRANTED"
	ConsentRevoked ConsentStatus = "REVOKED"
	ConsentExpired ConsentStatus = "EXPIRED"
	ConsentDeleted ConsentStatus = "DELETED"
)
```

#### `domain/value_object/dsar_type.go`

```go
package valueobject

type DSARType string

const (
	DSARTypeAccess          DSARType = "ACCESS"
	DSARTypeErasure         DSARType = "ERASURE"
	DSARTypeWithdrawConsent DSARType = "WITHDRAW_CONSENT"
)

func (d DSARType) IsValid() bool {
	switch d {
	case DSARTypeAccess, DSARTypeErasure, DSARTypeWithdrawConsent:
		return true
	}
	return false
}

func (d DSARType) String() string {
	return string(d)
}
```

#### `domain/value_object/dsar_status.go`

```go
package valueobject

type DSARStatus string

const (
	DSARStatusPending    DSARStatus = "PENDING"
	DSARStatusProcessing DSARStatus = "PROCESSING"
	DSARStatusCompleted  DSARStatus = "COMPLETED"
	DSARStatusRejected   DSARStatus = "REJECTED"
)
```

#### `domain/value_object/account_status.go`

```go
package valueobject

type AccountStatus string

const (
	AccountActive    AccountStatus = "ACTIVE"
	AccountSuspended AccountStatus = "SUSPENDED"
	AccountTerminated AccountStatus = "TERMINATED"
)

func (a AccountStatus) IsValid() bool {
	switch a {
	case AccountActive, AccountSuspended, AccountTerminated:
		return true
	}
	return false
}

func (a AccountStatus) String() string {
	return string(a)
}
```

---

### 3. Repository Interfaces

#### `domain/repository/consent_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type ConsentRepository interface {
	Save(ctx context.Context, log *entity.ConsentLog) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error)
	FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error)
	Revoke(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error
	FindExpiredRevoked(ctx context.Context, cutoffDate time.Time) ([]entity.ConsentLog, error)
	DeletePermanently(ctx context.Context, id uuid.UUID) error
	IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

#### `domain/repository/dsar_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
)

type DSARRepository interface {
	Save(ctx context.Context, request *entity.DSARRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

#### `domain/repository/audit_repository.go`

```go
package repository

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type AuditRepository interface {
	Save(ctx context.Context, audit *entity.AuditTrail) error
}
```

#### `domain/repository/user_account_status_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type UserAccountStatusRepository interface {
	Save(ctx context.Context, status *entity.UserAccountStatus) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error)
	FindSuspendedWithoutDeletion(ctx context.Context, cutoffDate time.Time) ([]entity.UserAccountStatus, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, status valueobject.AccountStatus) error
	UpdateDeletionConfirmed(ctx context.Context, userID uuid.UUID, confirmedAt time.Time) error
	UpdateAutoDeleted(ctx context.Context, userID uuid.UUID, deletedAt time.Time) error
}
```

---

### 4. Domain Services

#### `domain/service/deletion_policy_service.go`

```go
package service

import (
	"time"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type DeletionPolicyService struct{}

func NewDeletionPolicyService() *DeletionPolicyService {
	return &DeletionPolicyService{}
}

func (s *DeletionPolicyService) CanImmediateDeletion(status *entity.UserAccountStatus) bool {
	return status.Status == valueobject.AccountTerminated && status.DeletionConfirmedAt != nil
}

func (s *DeletionPolicyService) IsReadyForAutoDeletion(status *entity.UserAccountStatus, now time.Time) bool {
	if status.Status != valueobject.AccountSuspended {
		return false
	}
	if status.RetentionDeadline == nil {
		return false
	}
	return !status.RetentionDeadline.After(now)
}

func (s *DeletionPolicyService) CalculateRetentionDeadline(suspendedAt time.Time, retentionYears int) time.Time {
	return suspendedAt.AddDate(retentionYears, 0, 0)
}
```

#### `domain/service/localization_service.go`

```go
package service

import "context"

type LocalizationService interface {
	Translate(ctx context.Context, key, lang string) string
}
```

#### `domain/service/consent_validator.go`

```go
package service

import (
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type ConsentValidator struct{}

func NewConsentValidator() *ConsentValidator {
	return &ConsentValidator{}
}

func (v *ConsentValidator) ValidatePurpose(purpose string) bool {
	return valueobject.ConsentPurpose(purpose).IsValid()
}
```

---

### 5. Domain Errors

#### `domain/errors/errors.go`

```go
package domainerrors

import "errors"

var (
	ErrConsentNotFound          = errors.New("consent record not found")
	ErrConsentAlreadyRevoked    = errors.New("consent already revoked")
	ErrConsentExpired           = errors.New("consent expired")
	ErrDSARNotFound             = errors.New("DSAR request not found")
	ErrDSARAlreadyProcessed     = errors.New("DSAR request already processed")
	ErrInvalidPurpose           = errors.New("invalid consent purpose")
	ErrOTPExpired               = errors.New("OTP expired or invalid")
	ErrRateLimitExceeded        = errors.New("too many DSAR requests, please try again later")
	ErrRevokeNotAllowed         = errors.New("cannot revoke non-granted consent")
	ErrAccountNotFound          = errors.New("user account status not found")
	ErrAccountNotSuspended      = errors.New("account is not suspended")
	ErrAccountNotTerminated     = errors.New("account is not terminated")
	ErrDeletionNotConfirmed     = errors.New("deletion not confirmed for this account")
	ErrImmediateDeletionNotAllowed = errors.New("immediate deletion is not allowed for this account")
	ErrAccountNotActive         = errors.New("account is not active, cannot perform this action")
)
```

---

## ⚙️ Application Layer

### 1. Use Cases

#### `application/record_consent.go`

```go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
)

type RecordConsentUseCase struct {
	consentRepo       repository.ConsentRepository
	accountStatusRepo repository.UserAccountStatusRepository
	auditRepo         repository.AuditRepository
	cache             redis.ConsentCache
	producer          messaging.KafkaProducer
}

func NewRecordConsentUseCase(
	consentRepo repository.ConsentRepository,
	accountStatusRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	cache redis.ConsentCache,
	producer messaging.KafkaProducer,
) *RecordConsentUseCase {
	return &RecordConsentUseCase{
		consentRepo:       consentRepo,
		accountStatusRepo: accountStatusRepo,
		auditRepo:         auditRepo,
		cache:             cache,
		producer:          producer,
	}
}

type RecordConsentInput struct {
	UserID    uuid.UUID
	SessionID string
	Purposes  map[string]bool
	IPAddress string
	UserAgent string
}

func (uc *RecordConsentUseCase) Execute(ctx context.Context, input RecordConsentInput) error {
	// Check account status
	status, err := uc.accountStatusRepo.FindByUserID(ctx, input.UserID)
	if err == nil && (status.Status == valueobject.AccountSuspended || status.Status == valueobject.AccountTerminated) {
		return domainerrors.ErrAccountNotActive
	}

	for purposeStr, granted := range input.Purposes {
		if !granted {
			continue
		}
		purpose := valueobject.ConsentPurpose(purposeStr)
		if !purpose.IsValid() {
			return domainerrors.ErrInvalidPurpose
		}

		// Create new consent log
		log := entity.NewConsentLog(input.UserID, input.SessionID, purpose, input.IPAddress, input.UserAgent)

		// Save to DB
		if err := uc.consentRepo.Save(ctx, log); err != nil {
			return err
		}

		// Cache latest consent
		_ = uc.cache.Set(ctx, input.UserID, purpose, log.Status)

		// Publish event
		_ = uc.producer.PublishConsentEvent(ctx, log)

		// Audit trail
		audit := entity.NewAuditTrailWithMeta(
			&input.UserID,
			"CONSENT_GRANTED",
			map[string]interface{}{
				"purpose": purpose.String(),
				"status":  log.Status,
			},
			input.IPAddress,
			input.UserAgent,
		)
		_ = uc.auditRepo.Save(ctx, audit)
	}

	return nil
}
```

#### `application/revoke_consent.go`

```go
package application

import (
	"context"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
)

type RevokeConsentUseCase struct {
	consentRepo repository.ConsentRepository
	auditRepo   repository.AuditRepository
	cache       redis.ConsentCache
	producer    messaging.KafkaProducer
}

func NewRevokeConsentUseCase(
	consentRepo repository.ConsentRepository,
	auditRepo repository.AuditRepository,
	cache redis.ConsentCache,
	producer messaging.KafkaProducer,
) *RevokeConsentUseCase {
	return &RevokeConsentUseCase{
		consentRepo: consentRepo,
		auditRepo:   auditRepo,
		cache:       cache,
		producer:    producer,
	}
}

type RevokeConsentInput struct {
	UserID    uuid.UUID
	Purpose   valueobject.ConsentPurpose
	IPAddress string
	UserAgent string
}

func (uc *RevokeConsentUseCase) Execute(ctx context.Context, input RevokeConsentInput) error {
	// Check if consent exists and is active
	active, err := uc.consentRepo.IsConsentActive(ctx, input.UserID, input.Purpose)
	if err != nil {
		return err
	}
	if !active {
		return domainerrors.ErrConsentNotFound
	}

	// Revoke in DB
	if err := uc.consentRepo.Revoke(ctx, input.UserID, input.Purpose); err != nil {
		return err
	}

	// Update cache
	_ = uc.cache.Set(ctx, input.UserID, input.Purpose, valueobject.ConsentRevoked)

	// Publish revoke event
	_ = uc.producer.PublishRevokeEvent(ctx, input.UserID, input.Purpose)

	// Audit trail
	audit := entity.NewAuditTrailWithMeta(
		&input.UserID,
		"CONSENT_REVOKED",
		map[string]interface{}{
			"purpose": input.Purpose.String(),
		},
		input.IPAddress,
		input.UserAgent,
	)
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}
```

#### `application/immediate_deletion.go`

```go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/search/elasticsearch"
)

type ImmediateDeletionUseCase struct {
	consentRepo       repository.ConsentRepository
	dsarRepo          repository.DSARRepository
	accountStatusRepo repository.UserAccountStatusRepository
	auditRepo         repository.AuditRepository
	cache             redis.ConsentCache
	indexer           elasticsearch.ConsentIndexer
	producer          messaging.KafkaProducer
	policy            *service.DeletionPolicyService
}

func NewImmediateDeletionUseCase(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountStatusRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	cache redis.ConsentCache,
	indexer elasticsearch.ConsentIndexer,
	producer messaging.KafkaProducer,
	policy *service.DeletionPolicyService,
) *ImmediateDeletionUseCase {
	return &ImmediateDeletionUseCase{
		consentRepo:       consentRepo,
		dsarRepo:          dsarRepo,
		accountStatusRepo: accountStatusRepo,
		auditRepo:         auditRepo,
		cache:             cache,
		indexer:           indexer,
		producer:          producer,
		policy:            policy,
	}
}

func (uc *ImmediateDeletionUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	// Get account status
	status, err := uc.accountStatusRepo.FindByUserID(ctx, userID)
	if err != nil {
		return domainerrors.ErrAccountNotFound
	}

	// Check if immediate deletion is allowed
	if !uc.policy.CanImmediateDeletion(status) {
		return domainerrors.ErrImmediateDeletionNotAllowed
	}

	// Delete consent logs
	if err := uc.consentRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	// Delete DSAR requests
	if err := uc.dsarRepo.DeleteByUserID(ctx, userID); err != nil {
		// Log error but continue
	}

	// Delete Redis cache
	_ = uc.cache.DeleteAllByUser(ctx, userID)

	// Delete Elasticsearch index
	_ = uc.indexer.DeleteByUserID(ctx, userID)

	// Mark account as deleted
	status.MarkDeleted()
	if err := uc.accountStatusRepo.Save(ctx, status); err != nil {
		// Log error
	}

	// Publish immediate deletion event
	_ = uc.producer.PublishImmediateDeletionEvent(ctx, userID)

	// Audit trail
	audit := entity.NewAuditTrail(&userID, "IMMEDIATE_DELETION", map[string]interface{}{
		"user_id": userID.String(),
		"reason":  "termination_confirmed",
	})
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}
```

#### `application/auto_delete_expired_consents.go`

```go
package application

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/search/elasticsearch"
)

type AutoDeleteExpiredConsentsUseCase struct {
	consentRepo       repository.ConsentRepository
	dsarRepo          repository.DSARRepository
	accountStatusRepo repository.UserAccountStatusRepository
	auditRepo         repository.AuditRepository
	cache             redis.ConsentCache
	indexer           elasticsearch.ConsentIndexer
	producer          messaging.KafkaProducer
	policy            *service.DeletionPolicyService
}

func NewAutoDeleteExpiredConsentsUseCase(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountStatusRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	cache redis.ConsentCache,
	indexer elasticsearch.ConsentIndexer,
	producer messaging.KafkaProducer,
	policy *service.DeletionPolicyService,
) *AutoDeleteExpiredConsentsUseCase {
	return &AutoDeleteExpiredConsentsUseCase{
		consentRepo:       consentRepo,
		dsarRepo:          dsarRepo,
		accountStatusRepo: accountStatusRepo,
		auditRepo:         auditRepo,
		cache:             cache,
		indexer:           indexer,
		producer:          producer,
		policy:            policy,
	}
}

func (uc *AutoDeleteExpiredConsentsUseCase) Execute(ctx context.Context, retentionYears int) error {
	now := time.Now()

	// Find accounts that are suspended and past retention deadline
	expiredAccounts, err := uc.accountStatusRepo.FindSuspendedWithoutDeletion(ctx, now)
	if err != nil {
		return err
	}

	for _, status := range expiredAccounts {
		// Double-check policy
		if !uc.policy.IsReadyForAutoDeletion(&status, now) {
			continue
		}

		userID := status.UserID

		// Delete consent logs
		if err := uc.consentRepo.DeleteByUserID(ctx, userID); err != nil {
			// Log error and continue with next user
			continue
		}

		// Delete DSAR requests
		if err := uc.dsarRepo.DeleteByUserID(ctx, userID); err != nil {
			// Log error
		}

		// Delete Redis cache
		_ = uc.cache.DeleteAllByUser(ctx, userID)

		// Delete Elasticsearch index
		_ = uc.indexer.DeleteByUserID(ctx, userID)

		// Mark account as deleted
		status.MarkDeleted()
		if err := uc.accountStatusRepo.Save(ctx, &status); err != nil {
			// Log error
		}

		// Publish auto deletion event
		_ = uc.producer.PublishAutoDeletionEvent(ctx, userID, status.SuspendedAt)

		// Audit trail
		audit := entity.NewAuditTrail(&userID, "AUTO_DELETION", map[string]interface{}{
			"user_id":         userID.String(),
			"suspended_at":    status.SuspendedAt,
			"retention_period": retentionYears,
		})
		_ = uc.auditRepo.Save(ctx, audit)
	}

	return nil
}
```

#### `application/submit_dsar.go`

```go
package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
)

type SubmitDSARUseCase struct {
	dsarRepo          repository.DSARRepository
	accountStatusRepo repository.UserAccountStatusRepository
	auditRepo         repository.AuditRepository
	producer          messaging.KafkaProducer
}

func NewSubmitDSARUseCase(
	dsarRepo repository.DSARRepository,
	accountStatusRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	producer messaging.KafkaProducer,
) *SubmitDSARUseCase {
	return &SubmitDSARUseCase{
		dsarRepo:          dsarRepo,
		accountStatusRepo: accountStatusRepo,
		auditRepo:         auditRepo,
		producer:          producer,
	}
}

type SubmitDSARInput struct {
	UserID      uuid.UUID
	RequestType valueobject.DSARType
	IPAddress   string
	UserAgent   string
}

func (uc *SubmitDSARUseCase) Execute(ctx context.Context, input SubmitDSARInput) (*entity.DSARRequest, error) {
	// Check account status
	status, err := uc.accountStatusRepo.FindByUserID(ctx, input.UserID)
	if err == nil && (status.Status == valueobject.AccountSuspended || status.Status == valueobject.AccountTerminated) {
		return nil, domainerrors.ErrAccountNotActive
	}

	// Generate OTP
	otpCode, err := generateOTP()
	if err != nil {
		return nil, err
	}
	otpExpiry := time.Now().Add(15 * time.Minute)

	// Create DSAR request
	request := entity.NewDSARRequest(input.UserID, input.RequestType, otpCode, otpExpiry, input.IPAddress, input.UserAgent)

	if err := uc.dsarRepo.Save(ctx, request); err != nil {
		return nil, err
	}

	// Publish DSAR event for async processing
	_ = uc.producer.PublishDSAREvent(ctx, request)

	// Audit trail
	audit := entity.NewAuditTrailWithMeta(
		&input.UserID,
		"DSAR_SUBMITTED",
		map[string]interface{}{
			"request_id":   request.ID.String(),
			"request_type": input.RequestType.String(),
		},
		input.IPAddress,
		input.UserAgent,
	)
	_ = uc.auditRepo.Save(ctx, audit)

	return request, nil
}

func generateOTP() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:8], nil
}
```

#### `application/handle_account_event.go`

```go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type HandleAccountEventUseCase struct {
	accountStatusRepo repository.UserAccountStatusRepository
	auditRepo         repository.AuditRepository
	policy            *service.DeletionPolicyService
}

func NewHandleAccountEventUseCase(
	accountStatusRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	policy *service.DeletionPolicyService,
) *HandleAccountEventUseCase {
	return &HandleAccountEventUseCase{
		accountStatusRepo: accountStatusRepo,
		auditRepo:         auditRepo,
		policy:            policy,
	}
}

func (uc *HandleAccountEventUseCase) HandleSuspended(ctx context.Context, userID uuid.UUID, suspendedAt time.Time, retentionYears int) error {
	existing, err := uc.accountStatusRepo.FindByUserID(ctx, userID)
	if err != nil && err != domainerrors.ErrAccountNotFound {
		return err
	}

	var status *entity.UserAccountStatus
	if existing != nil {
		status = existing
	} else {
		status = entity.NewUserAccountStatus(userID, valueobject.AccountSuspended)
	}
	status.SetSuspended(suspendedAt, retentionYears)

	if err := uc.accountStatusRepo.Save(ctx, status); err != nil {
		return err
	}

	// Audit trail
	audit := entity.NewAuditTrail(&userID, "ACCOUNT_SUSPENDED", map[string]interface{}{
		"suspended_at":      suspendedAt,
		"retention_deadline": status.RetentionDeadline,
	})
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}

func (uc *HandleAccountEventUseCase) HandleTerminated(ctx context.Context, userID uuid.UUID, terminatedAt time.Time) error {
	existing, err := uc.accountStatusRepo.FindByUserID(ctx, userID)
	if err != nil && err != domainerrors.ErrAccountNotFound {
		return err
	}

	var status *entity.UserAccountStatus
	if existing != nil {
		status = existing
	} else {
		status = entity.NewUserAccountStatus(userID, valueobject.AccountTerminated)
	}
	status.SetTerminated(terminatedAt)

	if err := uc.accountStatusRepo.Save(ctx, status); err != nil {
		return err
	}

	// Audit trail
	audit := entity.NewAuditTrail(&userID, "ACCOUNT_TERMINATED", map[string]interface{}{
		"terminated_at": terminatedAt,
	})
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}

func (uc *HandleAccountEventUseCase) ConfirmDeletion(ctx context.Context, userID uuid.UUID) error {
	status, err := uc.accountStatusRepo.FindByUserID(ctx, userID)
	if err != nil {
		return domainerrors.ErrAccountNotFound
	}
	if status.Status != valueobject.AccountTerminated {
		return domainerrors.ErrAccountNotTerminated
	}

	status.ConfirmDeletion()
	if err := uc.accountStatusRepo.Save(ctx, status); err != nil {
		return err
	}

	// Audit trail
	audit := entity.NewAuditTrail(&userID, "DELETION_CONFIRMED", map[string]interface{}{
		"confirmed_at": status.DeletionConfirmedAt,
	})
	_ = uc.auditRepo.Save(ctx, audit)

	return nil
}
```

---

## 🗄️ Infrastructure Layer

### 1. GORM Models (พร้อม Prefix `pdpa_`)

#### `infrastructure/persistence/postgres/models.go`

```go
package postgres

import (
	"time"

	"github.com/google/uuid"
)

type PdpAPolicyModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Version       string    `gorm:"type:varchar(20);not null;index"`
	Title         string    `gorm:"type:text;not null"`
	Content       string    `gorm:"type:text;not null"`
	EffectiveDate time.Time `gorm:"type:date;not null"`
	IsActive      bool      `gorm:"default:false;index"`
	CreatedAt     time.Time `gorm:"default:now()"`
	UpdatedAt     time.Time `gorm:"default:now()"`
}

func (PdpAPolicyModel) TableName() string { return "pdpa_policies" }

type PdpAPurposeModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	IsRequired  bool      `gorm:"default:false"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
}

func (PdpAPurposeModel) TableName() string { return "pdpa_purposes" }

type PdpAConsentModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index:idx_pdpa_consents_user_id"`
	SessionID     string     `gorm:"type:varchar(255);index"`
	PurposeCode   string     `gorm:"type:varchar(50);not null;index:idx_pdpa_consents_purpose"`
	Status        string     `gorm:"type:varchar(20);not null;default:GRANTED;index:idx_pdpa_consents_status"`
	IPAddress     string     `gorm:"type:varchar(45)"`
	UserAgent     string     `gorm:"type:text"`
	GrantedAt     time.Time  `gorm:"default:now();index:idx_pdpa_consents_granted_at"`
	ExpiresAt     time.Time  `gorm:"type:timestamp"`
	RevokedAt     *time.Time `gorm:"type:timestamp;index:idx_pdpa_consents_revoked_at"`
	AutoDeletedAt *time.Time `gorm:"type:timestamp"`
	Purpose       *PdpAPurposeModel `gorm:"foreignKey:PurposeCode;references:Code"`
}

func (PdpAConsentModel) TableName() string { return "pdpa_consents" }

type PdpAUserRequestModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index:idx_pdpa_user_requests_user_id"`
	RequestType     string     `gorm:"type:varchar(20);not null"`
	Status          string     `gorm:"type:varchar(20);default:PENDING;index:idx_pdpa_user_requests_status"`
	RequestedAt     time.Time  `gorm:"default:now();index:idx_pdpa_user_requests_requested_at"`
	CompletedAt     *time.Time `gorm:"type:timestamp"`
	DataPayload     []byte     `gorm:"type:jsonb"`
	RejectionReason string     `gorm:"type:text"`
	OTPCode         string     `gorm:"type:varchar(255)"`
	OTPExpiredAt    time.Time  `gorm:"type:timestamp"`
	IPAddress       string     `gorm:"type:varchar(45)"`
	UserAgent       string     `gorm:"type:text"`
}

func (PdpAUserRequestModel) TableName() string { return "pdpa_user_requests" }

type PdpAUserAccountStatusModel struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	Status               string     `gorm:"type:varchar(20);not null;default:ACTIVE"`
	SuspendedAt          *time.Time `gorm:"type:timestamp"`
	TerminatedAt         *time.Time `gorm:"type:timestamp"`
	DeletionConfirmedAt  *time.Time `gorm:"type:timestamp"`
	RetentionDeadline    *time.Time `gorm:"type:timestamp;index:idx_pdpa_status_retention,where:status='SUSPENDED'"`
	AutoDeletedAt        *time.Time `gorm:"type:timestamp"`
	UpdatedAt            time.Time  `gorm:"default:now()"`
}

func (PdpAUserAccountStatusModel) TableName() string { return "pdpa_user_account_statuses" }

type PdpAAuditTrailModel struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    *uuid.UUID   `gorm:"type:uuid;index:idx_pdpa_audit_user_id"`
	Action    string       `gorm:"type:varchar(50);not null;index:idx_pdpa_audit_action"`
	Details   string       `gorm:"type:jsonb"`
	IPAddress string       `gorm:"type:varchar(45)"`
	UserAgent string       `gorm:"type:text"`
	CreatedAt time.Time    `gorm:"default:now();index:idx_pdpa_audit_created_at"`
}

func (PdpAAuditTrailModel) TableName() string { return "pdpa_audit_trails" }

type PdpARequestResponseModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID   uuid.UUID  `gorm:"type:uuid;not null;index:idx_pdpa_responses_request_id"`
	FileName    string     `gorm:"type:varchar(255)"`
	FilePath    string     `gorm:"type:varchar(500)"`
	FileSize    int64      `gorm:"type:bigint"`
	MimeType    string     `gorm:"type:varchar(100)"`
	DataContent string     `gorm:"type:text"`
	CreatedAt   time.Time  `gorm:"default:now()"`
	ExpiresAt   time.Time  `gorm:"type:timestamp"`
	Request     *PdpAUserRequestModel `gorm:"foreignKey:RequestID;references:ID"`
}

func (PdpARequestResponseModel) TableName() string { return "pdpa_request_responses" }
```

---

### 2. Repository Implementations

#### `infrastructure/persistence/postgres/consent_repo_impl.go`

```go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type consentRepoImpl struct {
	db *gorm.DB
}

func NewConsentRepository(db *gorm.DB) *consentRepoImpl {
	return &consentRepoImpl{db: db}
}

func (r *consentRepoImpl) Save(ctx context.Context, log *entity.ConsentLog) error {
	model := &PdpAConsentModel{
		ID:            log.ID,
		UserID:        log.UserID,
		SessionID:     log.SessionID,
		PurposeCode:   log.Purpose.String(),
		Status:        string(log.Status),
		IPAddress:     log.IPAddress,
		UserAgent:     log.UserAgent,
		GrantedAt:     log.GrantedAt,
		ExpiresAt:     log.ExpiresAt,
		RevokedAt:     log.RevokedAt,
		AutoDeletedAt: log.AutoDeletedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *consentRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error) {
	var models []PdpAConsentModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("granted_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *consentRepoImpl) FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error) {
	var model PdpAConsentModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose_code = ?", userID, purpose.String()).
		Order("granted_at DESC").
		First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *consentRepoImpl) Revoke(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&PdpAConsentModel{}).
		Where("user_id = ? AND purpose_code = ? AND status = ?", userID, purpose.String(), valueobject.ConsentGranted).
		Updates(map[string]interface{}{
			"status":     valueobject.ConsentRevoked,
			"revoked_at": now,
		}).Error
}

func (r *consentRepoImpl) FindExpiredRevoked(ctx context.Context, cutoffDate time.Time) ([]entity.ConsentLog, error) {
	var models []PdpAConsentModel
	err := r.db.WithContext(ctx).
		Where("status IN (?) AND revoked_at <= ?", []string{string(valueobject.ConsentRevoked), string(valueobject.ConsentExpired)}, cutoffDate).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *consentRepoImpl) DeletePermanently(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&PdpAConsentModel{}, "id = ?", id).Error
}

func (r *consentRepoImpl) IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&PdpAConsentModel{}).
		Where("user_id = ? AND purpose_code = ? AND status = ? AND expires_at > ?", userID, purpose.String(), valueobject.ConsentGranted, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *consentRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PdpAConsentModel{}).Error
}

func (r *consentRepoImpl) modelToEntity(m *PdpAConsentModel) *entity.ConsentLog {
	return &entity.ConsentLog{
		ID:            m.ID,
		UserID:        m.UserID,
		SessionID:     m.SessionID,
		Purpose:       valueobject.ConsentPurpose(m.PurposeCode),
		Status:        valueobject.ConsentStatus(m.Status),
		IPAddress:     m.IPAddress,
		UserAgent:     m.UserAgent,
		GrantedAt:     m.GrantedAt,
		ExpiresAt:     m.ExpiresAt,
		RevokedAt:     m.RevokedAt,
		AutoDeletedAt: m.AutoDeletedAt,
	}
}

func (r *consentRepoImpl) modelsToEntities(models []PdpAConsentModel) []entity.ConsentLog {
	entities := make([]entity.ConsentLog, len(models))
	for i, m := range models {
		entities[i] = *r.modelToEntity(&m)
	}
	return entities
}
```

#### `infrastructure/persistence/postgres/dsar_repo_impl.go`

```go
package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type dsarRepoImpl struct {
	db *gorm.DB
}

func NewDSARRepository(db *gorm.DB) *dsarRepoImpl {
	return &dsarRepoImpl{db: db}
}

func (r *dsarRepoImpl) Save(ctx context.Context, request *entity.DSARRequest) error {
	model := &PdpAUserRequestModel{
		ID:              request.ID,
		UserID:          request.UserID,
		RequestType:     string(request.RequestType),
		Status:          string(request.Status),
		RequestedAt:     request.RequestedAt,
		CompletedAt:     request.CompletedAt,
		DataPayload:     request.DataPayload,
		RejectionReason: request.RejectionReason,
		OTPCode:         request.OTPCode,
		OTPExpiredAt:    request.OTPExpiredAt,
		IPAddress:       request.IPAddress,
		UserAgent:       request.UserAgent,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *dsarRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error) {
	var model PdpAUserRequestModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domainerrors.ErrDSARNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *dsarRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error) {
	var models []PdpAUserRequestModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("requested_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *dsarRepoImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&PdpAUserRequestModel{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *dsarRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PdpAUserRequestModel{}).Error
}

func (r *dsarRepoImpl) modelToEntity(m *PdpAUserRequestModel) *entity.DSARRequest {
	return &entity.DSARRequest{
		ID:              m.ID,
		UserID:          m.UserID,
		RequestType:     valueobject.DSARType(m.RequestType),
		Status:          valueobject.DSARStatus(m.Status),
		RequestedAt:     m.RequestedAt,
		CompletedAt:     m.CompletedAt,
		DataPayload:     m.DataPayload,
		RejectionReason: m.RejectionReason,
		OTPCode:         m.OTPCode,
		OTPExpiredAt:    m.OTPExpiredAt,
		IPAddress:       m.IPAddress,
		UserAgent:       m.UserAgent,
	}
}

func (r *dsarRepoImpl) modelsToEntities(models []PdpAUserRequestModel) []entity.DSARRequest {
	entities := make([]entity.DSARRequest, len(models))
	for i, m := range models {
		entities[i] = *r.modelToEntity(&m)
	}
	return entities
}
```

#### `infrastructure/persistence/postgres/user_account_status_repo_impl.go`

```go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type userAccountStatusRepoImpl struct {
	db *gorm.DB
}

func NewUserAccountStatusRepository(db *gorm.DB) *userAccountStatusRepoImpl {
	return &userAccountStatusRepoImpl{db: db}
}

func (r *userAccountStatusRepoImpl) Save(ctx context.Context, status *entity.UserAccountStatus) error {
	model := &PdpAUserAccountStatusModel{
		ID:                   status.ID,
		UserID:               status.UserID,
		Status:               string(status.Status),
		SuspendedAt:          status.SuspendedAt,
		TerminatedAt:         status.TerminatedAt,
		DeletionConfirmedAt:  status.DeletionConfirmedAt,
		RetentionDeadline:    status.RetentionDeadline,
		AutoDeletedAt:        status.AutoDeletedAt,
		UpdatedAt:            status.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userAccountStatusRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error) {
	var model PdpAUserAccountStatusModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domainerrors.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userAccountStatusRepoImpl) FindSuspendedWithoutDeletion(ctx context.Context, cutoffDate time.Time) ([]entity.UserAccountStatus, error) {
	var models []PdpAUserAccountStatusModel
	err := r.db.WithContext(ctx).
		Where("status = ? AND retention_deadline <= ?", valueobject.AccountSuspended, cutoffDate).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *userAccountStatusRepoImpl) UpdateStatus(ctx context.Context, userID uuid.UUID, status valueobject.AccountStatus) error {
	return r.db.WithContext(ctx).
		Model(&PdpAUserAccountStatusModel{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *userAccountStatusRepoImpl) UpdateDeletionConfirmed(ctx context.Context, userID uuid.UUID, confirmedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&PdpAUserAccountStatusModel{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"deletion_confirmed_at": confirmedAt,
			"updated_at":            time.Now(),
		}).Error
}

func (r *userAccountStatusRepoImpl) UpdateAutoDeleted(ctx context.Context, userID uuid.UUID, deletedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&PdpAUserAccountStatusModel{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":          "DELETED",
			"auto_deleted_at": deletedAt,
			"updated_at":      time.Now(),
		}).Error
}

func (r *userAccountStatusRepoImpl) modelToEntity(m *PdpAUserAccountStatusModel) *entity.UserAccountStatus {
	return &entity.UserAccountStatus{
		ID:                   m.ID,
		UserID:               m.UserID,
		Status:               valueobject.AccountStatus(m.Status),
		SuspendedAt:          m.SuspendedAt,
		TerminatedAt:         m.TerminatedAt,
		DeletionConfirmedAt:  m.DeletionConfirmedAt,
		RetentionDeadline:    m.RetentionDeadline,
		AutoDeletedAt:        m.AutoDeletedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func (r *userAccountStatusRepoImpl) modelsToEntities(models []PdpAUserAccountStatusModel) []entity.UserAccountStatus {
	entities := make([]entity.UserAccountStatus, len(models))
	for i, m := range models {
		entities[i] = *r.modelToEntity(&m)
	}
	return entities
}
```

#### `infrastructure/persistence/postgres/audit_repo_impl.go`

```go
package postgres

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/entity"
)

type auditRepoImpl struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *auditRepoImpl {
	return &auditRepoImpl{db: db}
}

func (r *auditRepoImpl) Save(ctx context.Context, audit *entity.AuditTrail) error {
	detailsJSON, err := json.Marshal(audit.Details)
	if err != nil {
		return err
	}
	model := &PdpAAuditTrailModel{
		ID:        audit.ID,
		UserID:    audit.UserID,
		Action:    audit.Action,
		Details:   string(detailsJSON),
		IPAddress: audit.IPAddress,
		UserAgent: audit.UserAgent,
		CreatedAt: audit.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}
```

---

### 3. Redis Cache

#### `infrastructure/persistence/redis/consent_cache.go`

```go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type ConsentCache interface {
	Set(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose, status valueobject.ConsentStatus) error
	Get(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (valueobject.ConsentStatus, error)
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}

type redisConsentCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewConsentCache(client *redis.Client, ttl time.Duration) *redisConsentCache {
	return &redisConsentCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *redisConsentCache) Set(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose, status valueobject.ConsentStatus) error {
	key := c.key(userID, purpose)
	return c.client.Set(ctx, key, string(status), c.ttl).Err()
}

func (c *redisConsentCache) Get(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (valueobject.ConsentStatus, error) {
	key := c.key(userID, purpose)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return valueobject.ConsentStatus(val), nil
}

func (c *redisConsentCache) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("consent:%s:*", userID.String())
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (c *redisConsentCache) key(userID uuid.UUID, purpose valueobject.ConsentPurpose) string {
	return fmt.Sprintf("consent:%s:%s", userID.String(), purpose.String())
}
```

---

### 4. Kafka Producer

#### `infrastructure/messaging/kafka_producer.go`

```go
package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type KafkaProducer interface {
	PublishConsentEvent(ctx context.Context, log *entity.ConsentLog) error
	PublishRevokeEvent(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error
	PublishDSAREvent(ctx context.Context, request *entity.DSARRequest) error
	PublishImmediateDeletionEvent(ctx context.Context, userID uuid.UUID) error
	PublishAutoDeletionEvent(ctx context.Context, userID uuid.UUID, suspendedAt *time.Time) error
	PublishAccountEvent(ctx context.Context, eventType string, userID uuid.UUID, data map[string]interface{}) error
}

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{producer: producer}, nil
}

func (p *kafkaProducer) publish(topic string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (p *kafkaProducer) PublishConsentEvent(ctx context.Context, log *entity.ConsentLog) error {
	event := map[string]interface{}{
		"event_type": "consent.granted",
		"user_id":    log.UserID.String(),
		"purpose":    log.Purpose.String(),
		"granted_at": log.GrantedAt,
	}
	return p.publish("pdpa.consent.log", event)
}

func (p *kafkaProducer) PublishRevokeEvent(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error {
	event := map[string]interface{}{
		"event_type": "consent.revoked",
		"user_id":    userID.String(),
		"purpose":    purpose.String(),
		"revoked_at": time.Now(),
	}
	return p.publish("pdpa.consent.revoked", event)
}

func (p *kafkaProducer) PublishDSAREvent(ctx context.Context, request *entity.DSARRequest) error {
	event := map[string]interface{}{
		"event_type":   "dsar.request",
		"dsar_id":      request.ID.String(),
		"user_id":      request.UserID.String(),
		"request_type": request.RequestType.String(),
		"requested_at": request.RequestedAt,
	}
	return p.publish("pdpa.dsar.request", event)
}

func (p *kafkaProducer) PublishImmediateDeletionEvent(ctx context.Context, userID uuid.UUID) error {
	event := map[string]interface{}{
		"event_type": "data.immediate_deletion",
		"user_id":    userID.String(),
		"deleted_at": time.Now(),
	}
	return p.publish("pdpa.data.deleted", event)
}

func (p *kafkaProducer) PublishAutoDeletionEvent(ctx context.Context, userID uuid.UUID, suspendedAt *time.Time) error {
	event := map[string]interface{}{
		"event_type":  "data.auto_deletion",
		"user_id":     userID.String(),
		"suspended_at": suspendedAt,
		"deleted_at":  time.Now(),
	}
	return p.publish("pdpa.data.deleted", event)
}

func (p *kafkaProducer) PublishAccountEvent(ctx context.Context, eventType string, userID uuid.UUID, data map[string]interface{}) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"user_id":    userID.String(),
		"data":       data,
		"timestamp":  time.Now(),
	}
	return p.publish("pdpa.account.event", event)
}
```

---

### 5. Scheduler

#### `infrastructure/scheduler/consent_cleanup_job.go`

```go
package scheduler

import (
	"context"
	"log"
	"time"

	"icmongolang/internal/modules/pdpa/application"
)

type ConsentCleanupJob struct {
	useCase *application.AutoDeleteExpiredConsentsUseCase
}

func NewConsentCleanupJob(uc *application.AutoDeleteExpiredConsentsUseCase) *ConsentCleanupJob {
	return &ConsentCleanupJob{useCase: uc}
}

func (j *ConsentCleanupJob) Run() {
	log.Println("Running consent cleanup job...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := j.useCase.Execute(ctx, 1); err != nil {
		log.Printf("Error cleaning up expired consents: %v", err)
	} else {
		log.Println("Consent cleanup completed successfully")
	}
}
```

---

### 6. Elasticsearch Indexer (ตัวอย่าง)

#### `infrastructure/search/elasticsearch/consent_indexer.go`

```go
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
)

type ConsentIndexer interface {
	IndexConsent(ctx context.Context, consent interface{}) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type consentIndexer struct {
	client *elasticsearch.Client
	index  string
}

func NewConsentIndexer(addresses []string, index string) (ConsentIndexer, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consentIndexer{
		client: client,
		index:  index,
	}, nil
}

func (i *consentIndexer) IndexConsent(ctx context.Context, consent interface{}) error {
	data, err := json.Marshal(consent)
	if err != nil {
		return err
	}
	_, err = i.client.Index(i.index, bytes.NewReader(data))
	return err
}

func (i *consentIndexer) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`{"query":{"term":{"user_id":"%s"}}}`, userID.String())
	_, err := i.client.DeleteByQuery(
		[]string{i.index},
		bytes.NewReader([]byte(query)),
	)
	return err
}
```

---

## 🌐 Interface Layer

### 1. HTTP Handlers

#### `interfaces/http/consent_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type ConsentHandler struct {
	recordUC *application.RecordConsentUseCase
	revokeUC *application.RevokeConsentUseCase
}

func NewConsentHandler(recordUC *application.RecordConsentUseCase, revokeUC *application.RevokeConsentUseCase) *ConsentHandler {
	return &ConsentHandler{
		recordUC: recordUC,
		revokeUC: revokeUC,
	}
}

func (h *ConsentHandler) RecordConsent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var req struct {
		Purposes map[string]bool `json:"purposes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := application.RecordConsentInput{
		UserID:    uid,
		SessionID: c.GetString("session_id"),
		Purposes:  req.Purposes,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if err := h.recordUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "consent recorded successfully"})
}

func (h *ConsentHandler) RevokeConsent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var req struct {
		Purpose string `json:"purpose"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purpose := valueobject.ConsentPurpose(req.Purpose)
	if !purpose.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purpose"})
		return
	}

	input := application.RevokeConsentInput{
		UserID:    uid,
		Purpose:   purpose,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if err := h.revokeUC.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "consent revoked successfully"})
}
```

#### `interfaces/http/dsar_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application"
	"icmongolang/internal/modules/pdpa/domain/value_object"
)

type DSARHandler struct {
	submitUC *application.SubmitDSARUseCase
}

func NewDSARHandler(submitUC *application.SubmitDSARUseCase) *DSARHandler {
	return &DSARHandler{submitUC: submitUC}
}

func (h *DSARHandler) Submit(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var req struct {
		RequestType string `json:"request_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestType := valueobject.DSARType(req.RequestType)
	if !requestType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request_type"})
		return
	}

	input := application.SubmitDSARInput{
		UserID:      uid,
		RequestType: requestType,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	request, err := h.submitUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         request.ID.String(),
		"status":     string(request.Status),
		"created_at": request.RequestedAt,
	})
}
```

#### `interfaces/http/deletion_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
)

type DeletionHandler struct {
	immediateUC *application.ImmediateDeletionUseCase
	handleUC    *application.HandleAccountEventUseCase
}

func NewDeletionHandler(immediateUC *application.ImmediateDeletionUseCase, handleUC *application.HandleAccountEventUseCase) *DeletionHandler {
	return &DeletionHandler{
		immediateUC: immediateUC,
		handleUC:    handleUC,
	}
}

func (h *DeletionHandler) ConfirmDeletion(c *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check permission: only admin or DPO
	// ...

	if err := h.handleUC.ConfirmDeletion(c.Request.Context(), req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.immediateUC.Execute(c.Request.Context(), req.UserID); err != nil {
		if err == domainerrors.ErrImmediateDeletionNotAllowed {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data deleted immediately"})
}

func (h *DeletionHandler) AccountEvent(c *gin.Context) {
	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Process or forward to Kafka
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
```

---

### 2. Routes

#### `interfaces/http/routes.go`

```go
package http

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Consent  *ConsentHandler
	DSAR     *DSARHandler
	Deletion *DeletionHandler
}

func RegisterRoutes(r *gin.RouterGroup, handlers *Handlers, authMiddleware gin.HandlerFunc) {
	pdpa := r.Group("/pdpa")
	pdpa.Use(authMiddleware)

	pdpa.POST("/consent", handlers.Consent.RecordConsent)
	pdpa.DELETE("/consent", handlers.Consent.RevokeConsent)

	pdpa.POST("/dsar", handlers.DSAR.Submit)

	pdpa.POST("/deletion/confirm", handlers.Deletion.ConfirmDeletion)
	pdpa.POST("/events/account", handlers.Deletion.AccountEvent)
}
```

---

### 3. Localization (i18n)

#### `interfaces/localization/i18n.go`

```go
package localization

import (
	"context"
	"embed"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed locales/*.toml
var localesFS embed.FS

type I18nManager struct {
	translations map[string]map[string]string
	defaultLang  string
}

func NewI18nManager(defaultLang string) (*I18nManager, error) {
	m := &I18nManager{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
	}
	for _, lang := range []string{"en", "th"} {
		if err := m.loadLocale(lang); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *I18nManager) loadLocale(lang string) error {
	path := "locales/" + lang + ".toml"
	data, err := localesFS.ReadFile(path)
	if err != nil {
		return err
	}
	var locale map[string]interface{}
	if err := toml.Unmarshal(data, &locale); err != nil {
		return err
	}
	translations := make(map[string]string)
	for key, val := range locale {
		if msg, ok := val.(map[string]interface{}); ok {
			if other, ok := msg["other"].(string); ok {
				translations[key] = other
			}
		}
	}
	m.translations[lang] = translations
	return nil
}

func (m *I18nManager) Translate(ctx context.Context, key, lang string) string {
	if lang == "" {
		lang = m.defaultLang
	}
	if trans, ok := m.translations[lang][key]; ok {
		return trans
	}
	if trans, ok := m.translations[m.defaultLang][key]; ok {
		return trans
	}
	return key
}

func (m *I18nManager) GetLangFromHeader(header string) string {
	if header == "" {
		return m.defaultLang
	}
	parts := strings.Split(header, ",")
	if len(parts) > 0 {
		lang := strings.TrimSpace(parts[0])
		if _, ok := m.translations[lang]; ok {
			return lang
		}
	}
	return m.defaultLang
}
```

#### `interfaces/localization/locales/en.toml`

```toml
[ImmediateDeletionSuccess]
other = "Data deleted immediately as per termination confirmation"

[AutoDeletionExecuted]
other = "Expired data (suspended + 1 year) has been deleted"

[AccountSuspended]
other = "Account suspended, data will be retained for 1 year"

[AccountTerminated]
other = "Account terminated, pending deletion confirmation"

[DeletionConfirmed]
other = "Deletion confirmed, data will be removed immediately"

[ErrAccountNotActive]
other = "Account is not active, cannot perform this action"
```

#### `interfaces/localization/locales/th.toml`

```toml
[ImmediateDeletionSuccess]
other = "ลบข้อมูลทันทีตามการยืนยันการยกเลิกใช้งาน"

[AutoDeletionExecuted]
other = "ลบข้อมูลที่หมดอายุ (ระงับ + 1 ปี) เรียบร้อยแล้ว"

[AccountSuspended]
other = "ระงับบัญชีผู้ใช้ ข้อมูลจะถูกเก็บไว้ 1 ปี"

[AccountTerminated]
other = "ยกเลิกบัญชีผู้ใช้ รอการยืนยันการลบ"

[DeletionConfirmed]
other = "ยืนยันการลบแล้ว ข้อมูลจะถูกลบทันที"

[ErrAccountNotActive]
other = "บัญชีผู้ใช้ไม่สามารถใช้งานได้ ไม่สามารถดำเนินการนี้"
```

---

## 🚀 Main Entry Points

### `cmd/api/main.go`

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/application"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/postgres"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/search/elasticsearch"
	"icmongolang/internal/modules/pdpa/interfaces/http"
)

func main() {
	// Database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=admin password=secret dbname=pdpa_db port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	}

	// Kafka
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	producer, err := messaging.NewKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}

	// Repositories
	consentRepo := postgres.NewConsentRepository(db)
	dsarRepo := postgres.NewDSARRepository(db)
	auditRepo := postgres.NewAuditRepository(db)
	accountStatusRepo := postgres.NewUserAccountStatusRepository(db)

	// Cache
	cache := redis.NewConsentCache(redisClient, 30*24*time.Hour)

	// Indexer (Elasticsearch)
	indexer, _ := elasticsearch.NewConsentIndexer(
		[]string{os.Getenv("ELASTICSEARCH_URL")},
		"pdpa_consents",
	)

	// Domain services
	policy := service.NewDeletionPolicyService()

	// Use cases
	recordConsentUC := application.NewRecordConsentUseCase(consentRepo, accountStatusRepo, auditRepo, cache, producer)
	revokeConsentUC := application.NewRevokeConsentUseCase(consentRepo, auditRepo, cache, producer)
	immediateDeletionUC := application.NewImmediateDeletionUseCase(consentRepo, dsarRepo, accountStatusRepo, auditRepo, cache, indexer, producer, policy)
	handleAccountUC := application.NewHandleAccountEventUseCase(accountStatusRepo, auditRepo, policy)
	submitDSARUC := application.NewSubmitDSARUseCase(dsarRepo, accountStatusRepo, auditRepo, producer)

	// Handlers
	consentHandler := http.NewConsentHandler(recordConsentUC, revokeConsentUC)
	dsarHandler := http.NewDSARHandler(submitDSARUC)
	deletionHandler := http.NewDeletionHandler(immediateDeletionUC, handleAccountUC)

	handlers := &http.Handlers{
		Consent:  consentHandler,
		DSAR:     dsarHandler,
		Deletion: deletionHandler,
	}

	// Router
	r := gin.Default()
	authMiddleware := func(c *gin.Context) {
		// Placeholder JWT validation
		c.Set("user_id", uuid.New())
		c.Next()
	}
	api := r.Group("/api")
	http.RegisterRoutes(api, handlers, authMiddleware)

	log.Println("API server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

### `cmd/scheduler/main.go`

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/application"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/postgres"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/scheduler"
	"icmongolang/internal/modules/pdpa/infrastructure/search/elasticsearch"
)

func main() {
	// Database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=admin password=secret dbname=pdpa_db port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	}

	// Kafka
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	producer, err := messaging.NewKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}

	// Repositories
	consentRepo := postgres.NewConsentRepository(db)
	dsarRepo := postgres.NewDSARRepository(db)
	auditRepo := postgres.NewAuditRepository(db)
	accountStatusRepo := postgres.NewUserAccountStatusRepository(db)

	// Cache & Indexer
	cache := redis.NewConsentCache(redisClient, 30*24*time.Hour)
	indexer, _ := elasticsearch.NewConsentIndexer(
		[]string{os.Getenv("ELASTICSEARCH_URL")},
		"pdpa_consents",
	)

	policy := service.NewDeletionPolicyService()

	autoDeleteUC := application.NewAutoDeleteExpiredConsentsUseCase(
		consentRepo, dsarRepo, accountStatusRepo, auditRepo, cache, indexer, producer, policy,
	)

	job := scheduler.NewConsentCleanupJob(autoDeleteUC)

	c := cron.New()
	_, err = c.AddFunc("0 2 * * *", job.Run)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	log.Println("Scheduler started. Running at 2:00 AM daily.")
	select {}
}
```

---
 
---

## 📄 Environment Variables (`.env`)

```env
# PostgreSQL
DB_DSN=host=postgres user=admin password=secret dbname=pdpa_db port=5432 sslmode=disable

# Redis
REDIS_ADDR=redis:6379

# Kafka
KAFKA_BROKERS=kafka:9092

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200

# JWT
JWT_SECRET=your_secret_key

# Retention years for suspended accounts
RETENTION_YEARS=1

# Default language
DEFAULT_LANG=en
```

---

## 📦 Go Module (`go.mod`)

```go
module icmongolang

go 1.21

require (
    github.com/BurntSushi/toml v1.3.2
    github.com/IBM/sarama v1.41.3
    github.com/elastic/go-elasticsearch/v8 v8.11.0
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/google/uuid v1.5.0
    github.com/robfig/cron/v3 v3.0.1
    gorm.io/driver/postgres v1.5.4
    gorm.io/gorm v1.25.5
)
```

---

## สรุป

ระบบนี้ถูกออกแบบตามหลัก **Clean Architecture + DDD** พร้อม:

- ✅ **Domain Layer** - Entities, Value Objects, Repository Interfaces, Domain Services, Domain Errors
- ✅ **Application Layer** - Use Cases ที่เป็น business logic
- ✅ **Infrastructure Layer** - PostgreSQL, Redis, Kafka, Elasticsearch, Scheduler
- ✅ **Interface Layer** - HTTP Handlers, Routes, Localization
- ✅ **Prefix `pdpa_`** ในทุกตารางเพื่อแยกโมดูล

 