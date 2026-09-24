# ระบบ PDPA Module - สร้างทีละไฟล์ (Clean Architecture + DDD + EDA)

ผมจะสร้างระบบทีละไฟล์ตามลำดับ dependency จากชั้นในสุด (Domain) ออกมาชั้นนอก (Infrastructure/Interface)

---

## 1. Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                     INTERFACE LAYER (HTTP/WS)                        │
│  ConsentHandler │ DSARHandler │ AccountHandler │ PolicyHandler       │
│  AdminHandler   │ WebSocket Broadcaster │ i18n                       │
└────────────────────────────┬────────────────────────────────────────┘
                             │ (Commands/Queries)
┌────────────────────────────▼────────────────────────────────────────┐
│                      APPLICATION LAYER                               │
│  Commands: RecordConsent, RevokeConsent, SubmitDSAR, ProcessDSAR,   │
│            ConfirmDeletion, ImmediateDeletion, AutoDeleteExpired,   │
│            HandleAccountSuspended, HandleAccountTerminated          │
│  Queries:  GetConsentHistory, GetDSARStatus, GetAdminReport         │
│  EventHandlers: ConsentGranted, ConsentRevoked, DSARSubmitted,      │
│                 DataDeletion, AccountSuspended, LLMAnalysis,        │
│                 EmailNotification, BlockchainRecord                 │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────────┐
│                        DOMAIN LAYER                                  │
│  Aggregates: ConsentLog, DSARRequest, UserAccountStatus             │
│  Entities:   AuditTrail, PrivacyPolicy                              │
│  ValueObjects: ConsentPurpose, ConsentStatus, DSARType, DSARStatus, │
│                AccountStatus                                        │
│  DomainEvents: ConsentGranted, ConsentRevoked, DSARSubmitted, ...  │
│  Services: DeletionPolicyService, AnonymizationService              │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────────┐
│                     INFRASTRUCTURE LAYER                             │
│  PostgreSQL (GORM) │ Redis Cache │ Kafka Producer/Consumers         │
│  Outbox Publisher  │ Elasticsearch │ Scheduler │ External Clients   │
└─────────────────────────────────────────────────────────────────────┘

Kafka Topics Flow:
  pdpa.consent.granted    ──► ConsentGrantedHandler ──► Email + ES + Cache
  pdpa.consent.revoked    ──► ConsentRevokedHandler ──► Cache + Email
  pdpa.dsar.submitted     ──► DSARSubmittedHandler  ──► OTP Email + LLM
  pdpa.dsar.completed     ──► DSARCompletedHandler  ──► Email + WS
  pdpa.data.deletion_req  ──► DataDeletionHandler   ──► Multi-service
  pdpa.data.deleted       ──► BlockchainRecord      ──► Blockchain
  pdpa.account.suspended  ──► AccountSuspended      ──► Schedule + Email
  pdpa.account.terminated ──► AccountTerminated     ──► Immediate del.
  pdpa.llm.analysis.req   ──► LLMAnalysisHandler    ──► LLM API
```

---

## 2. File Tree

```
	internal/modules/pdpa/
	├── domain/
	│   ├── entity/
	│   │   ├── consent_log.go
	│   │   ├── dsar_request.go
	│   │   ├── user_account_status.go
	│   │   ├── audit_trail.go
	│   │   └── privacy_policy.go
	│   ├── value_object/
	│   │   ├── consent_purpose.go
	│   │   ├── consent_status.go
	│   │   ├── dsar_type.go
	│   │   ├── dsar_status.go
	│   │   └── account_status.go
	│   ├── event/
	│   │   ├── event_metadata.go
	│   │   ├── consent_events.go
	│   │   ├── dsar_events.go
	│   │   ├── account_events.go
	│   │   └── data_events.go
	│   ├── repository/
	│   │   ├── consent_repository.go
	│   │   ├── dsar_repository.go
	│   │   ├── account_status_repository.go
	│   │   ├── audit_repository.go
	│   │   ├── policy_repository.go
	│   │   └── outbox_repository.go
	│   ├── service/
	│   │   ├── consent_validator.go
	│   │   ├── deletion_policy_service.go
	│   │   └── anonymization_service.go
	│   └── errors/
	│       └── errors.go
	├── application/
	│   ├── command/
	│   │   ├── record_consent.go
	│   │   ├── revoke_consent.go
	│   │   ├── submit_dsar.go
	│   │   ├── process_dsar.go
	│   │   ├── confirm_deletion.go
	│   │   ├── immediate_deletion.go
	│   │   ├── auto_delete_expired.go
	│   │   ├── handle_account_suspended.go
	│   │   ├── handle_account_terminated.go
	│   │   └── publish_privacy_policy.go
	│   ├── query/
	│   │   ├── get_consent_history.go
	│   │   ├── get_dsar_status.go
	│   │   ├── get_admin_report.go
	│   │   └── get_active_privacy_policy.go
	│   ├── event_handler/
	│   │   ├── consent_granted_handler.go
	│   │   ├── consent_revoked_handler.go
	│   │   ├── dsar_submitted_handler.go
	│   │   ├── data_deletion_handler.go
	│   │   ├── account_suspended_handler.go
	│   │   ├── llm_analysis_handler.go
	│   │   ├── email_notification_handler.go
	│   │   └── blockchain_record_handler.go
	│   └── dto/
	│       └── dto.go
	├── infrastructure/
	│   ├── persistence/
	│   │   └── postgres/
	│   │       ├── models.go
	│   │       ├── consent_repo_impl.go
	│   │       ├── dsar_repo_impl.go
	│   │       ├── account_status_repo_impl.go
	│   │       ├── audit_repo_impl.go
	│   │       ├── policy_repo_impl.go
	│   │       └── outbox_repo_impl.go
	│   ├── cache/redis/
	│   │   └── consent_cache.go
	│   ├── messaging/
	│   │   ├── kafka/
	│   │   │   ├── producer.go
	│   │   │   └── event_serializer.go
	│   │   └── outbox/
	│   │       └── outbox_publisher.go
	│   └── external/
	│       ├── llm_client.go
	│       ├── email_client.go
	│       └── blockchain_client.go
	├── interfaces/
	│   └── http/
	│       ├── consent_handler.go
	│       ├── dsar_handler.go
	│       ├── account_handler.go
	│       ├── policy_handler.go
	│       ├── admin_handler.go
	│       ├── dto.go
	│       └── routes.go
	└── module.go
```
---

## 3. Code Implementation

### 3.1.1 Value Object: ConsentPurpose

**`internal/modules/pdpa/domain/value_object/consent_purpose.go`**

```go
package valueobject

import (
	"fmt"
	"strings"
)

// ConsentPurpose วัตถุประสงค์ในการขอความยินยอม
// ConsentPurpose represents the purpose for which consent is requested
type ConsentPurpose string

const (
	// PurposeNecessary - ข้อมูลจำเป็นเพื่อการให้บริการ (บังคับ)
	// Necessary data for service provision (mandatory)
	PurposeNecessary ConsentPurpose = "NECESSARY"

	// PurposeAnalytics - การวิเคราะห์ข้อมูลเพื่อพัฒนาบริการ
	// Analytics for service improvement
	PurposeAnalytics ConsentPurpose = "ANALYTICS"

	// PurposeMarketing - การตลาดและการโฆษณา
	// Marketing and advertising
	PurposeMarketing ConsentPurpose = "MARKETING"

	// PurposeAccountSystem - ระบบบัญชีผู้ใช้งาน
	// User account system
	PurposeAccountSystem ConsentPurpose = "ACCOUNT_SYSTEM"

	// PurposeUsageLogs - ประวัติการใช้งานระบบ
	// System usage logs
	PurposeUsageLogs ConsentPurpose = "USAGE_LOGS"

	// PurposeTransactionHistory - ประวัติการทำธุรกรรม
	// Transaction history
	PurposeTransactionHistory ConsentPurpose = "TRANSACTION_HISTORY"
)

// allPurposes รายการ purpose ทั้งหมดที่รองรับ
// allPurposes holds all supported purposes
var allPurposes = []ConsentPurpose{
	PurposeNecessary,
	PurposeAnalytics,
	PurposeMarketing,
	PurposeAccountSystem,
	PurposeUsageLogs,
	PurposeTransactionHistory,
}

// NewConsentPurpose สร้าง ConsentPurpose พร้อม validation
// NewConsentPurpose creates ConsentPurpose with validation
func NewConsentPurpose(s string) (ConsentPurpose, error) {
	p := ConsentPurpose(strings.ToUpper(strings.TrimSpace(s)))
	if !p.IsValid() {
		return "", fmt.Errorf("invalid consent purpose: %s", s)
	}
	return p, nil
}

// IsValid ตรวจสอบว่า purpose ถูกต้องหรือไม่
// IsValid checks if the purpose is valid
func (p ConsentPurpose) IsValid() bool {
	for _, valid := range allPurposes {
		if p == valid {
			return true
		}
	}
	return false
}

// IsRequired ตรวจสอบว่าเป็น purpose ที่บังคับหรือไม่
// IsRequired checks if this purpose is mandatory
func (p ConsentPurpose) IsRequired() bool {
	return p == PurposeNecessary
}

// String คืนค่า string representation
// String returns the string representation
func (p ConsentPurpose) String() string {
	return string(p)
}

// AllPurposes คืนค่ารายการ purpose ทั้งหมด
// AllPurposes returns all supported purposes
func AllPurposes() []ConsentPurpose {
	result := make([]ConsentPurpose, len(allPurposes))
	copy(result, allPurposes)
	return result
}
```

---

### 3.1.2 Value Object: ConsentStatus

**`internal/modules/pdpa/domain/value_object/consent_status.go`**

```go
package valueobject

// ConsentStatus สถานะของความยินยอม
// ConsentStatus represents the status of consent
type ConsentStatus string

const (
	// ConsentGranted - ผู้ใช้ให้ความยินยอม
	// User has granted consent
	ConsentGranted ConsentStatus = "GRANTED"

	// ConsentRevoked - ผู้ใช้ถอนความยินยอม
	// User has revoked consent
	ConsentRevoked ConsentStatus = "REVOKED"

	// ConsentExpired - ความยินยอมหมดอายุ
	// Consent has expired
	ConsentExpired ConsentStatus = "EXPIRED"

	// ConsentDeleted - ความยินยอมถูกลบ (auto-delete)
	// Consent has been deleted (auto-delete)
	ConsentDeleted ConsentStatus = "DELETED"
)

// IsValid ตรวจสอบว่าสถานะถูกต้องหรือไม่
// IsValid checks if the status is valid
func (s ConsentStatus) IsValid() bool {
	switch s {
	case ConsentGranted, ConsentRevoked, ConsentExpired, ConsentDeleted:
		return true
	}
	return false
}

// IsActive ตรวจสอบว่าความยินยอมยังมีผลอยู่หรือไม่
// IsActive checks if consent is still active
func (s ConsentStatus) IsActive() bool {
	return s == ConsentGranted
}

// String คืนค่า string representation
// String returns the string representation
func (s ConsentStatus) String() string {
	return string(s)
}
```

---

### 3.1.3 Value Object: DSARType

**`internal/modules/pdpa/domain/value_object/dsar_type.go`**

```go
package valueobject

import "fmt"

// DSARType ประเภทของคำร้องขอใช้สิทธิ์
// DSARType represents the type of Data Subject Access Request
type DSARType string

const (
	// DSARTypeAccess - ขอเข้าถึงข้อมูล
	// Request to access data
	DSARTypeAccess DSARType = "ACCESS"

	// DSARTypeErasure - ขอลบข้อมูล (Right to be Forgotten)
	// Request to erase data (Right to be Forgotten)
	DSARTypeErasure DSARType = "ERASURE"

	// DSARTypeWithdrawConsent - ขอถอนความยินยอม
	// Request to withdraw consent
	DSARTypeWithdrawConsent DSARType = "WITHDRAW_CONSENT"
)

// NewDSARType สร้าง DSARType พร้อม validation
// NewDSARType creates DSARType with validation
func NewDSARType(s string) (DSARType, error) {
	t := DSARType(s)
	if !t.IsValid() {
		return "", fmt.Errorf("invalid DSAR type: %s", s)
	}
	return t, nil
}

// IsValid ตรวจสอบว่าประเภทถูกต้องหรือไม่
// IsValid checks if the type is valid
func (t DSARType) IsValid() bool {
	switch t {
	case DSARTypeAccess, DSARTypeErasure, DSARTypeWithdrawConsent:
		return true
	}
	return false
}

// String คืนค่า string representation
// String returns the string representation
func (t DSARType) String() string {
	return string(t)
}
```

---

### 3.1.4 Value Object: DSARStatus

**`internal/modules/pdpa/domain/value_object/dsar_status.go`**

```go
package valueobject

// DSARStatus สถานะของคำร้อง
// DSARStatus represents the status of a DSAR request
type DSARStatus string

const (
	// DSARStatusPending - รอดำเนินการ
	// Waiting for processing
	DSARStatusPending DSARStatus = "PENDING"

	// DSARStatusProcessing - กำลังดำเนินการ
	// Being processed
	DSARStatusProcessing DSARStatus = "PROCESSING"

	// DSARStatusCompleted - ดำเนินการเสร็จสิ้น
	// Completed
	DSARStatusCompleted DSARStatus = "COMPLETED"

	// DSARStatusRejected - ถูกปฏิเสธ
	// Rejected
	DSARStatusRejected DSARStatus = "REJECTED"
)

// IsValid ตรวจสอบว่าสถานะถูกต้องหรือไม่
// IsValid checks if the status is valid
func (s DSARStatus) IsValid() bool {
	switch s {
	case DSARStatusPending, DSARStatusProcessing, DSARStatusCompleted, DSARStatusRejected:
		return true
	}
	return false
}

// IsTerminal ตรวจสอบว่าเป็นสถานะสุดท้ายหรือไม่
// IsTerminal checks if this is a terminal status
func (s DSARStatus) IsTerminal() bool {
	return s == DSARStatusCompleted || s == DSARStatusRejected
}

// String คืนค่า string representation
// String returns the string representation
func (s DSARStatus) String() string {
	return string(s)
}
```

---

### 3.1.5 Value Object: AccountStatus

**`internal/modules/pdpa/domain/value_object/account_status.go`**

```go
package valueobject

// AccountStatus สถานะบัญชีผู้ใช้
// AccountStatus represents the status of a user account
type AccountStatus string

const (
	// AccountActive - บัญชีปกติ
	// Account is active
	AccountActive AccountStatus = "ACTIVE"

	// AccountSuspended - ระงับบัญชี (เก็บข้อมูล 1 ปี)
	// Account suspended (retain data for 1 year)
	AccountSuspended AccountStatus = "SUSPENDED"

	// AccountTerminated - ยกเลิกบัญชี (ลบทันทีเมื่อยืนยัน)
	// Account terminated (delete immediately upon confirmation)
	AccountTerminated AccountStatus = "TERMINATED"

	// AccountDeleted - ลบข้อมูลแล้ว
	// Data has been deleted
	AccountDeleted AccountStatus = "DELETED"
)

// IsValid ตรวจสอบว่าสถานะถูกต้องหรือไม่
// IsValid checks if the status is valid
func (s AccountStatus) IsValid() bool {
	switch s {
	case AccountActive, AccountSuspended, AccountTerminated, AccountDeleted:
		return true
	}
	return false
}

// IsActive ตรวจสอบว่าบัญชีใช้งานได้หรือไม่
// IsActive checks if account is active
func (s AccountStatus) IsActive() bool {
	return s == AccountActive
}

// IsDeleted ตรวจสอบว่าถูกลบแล้วหรือไม่
// IsDeleted checks if account has been deleted
func (s AccountStatus) IsDeleted() bool {
	return s == AccountDeleted
}

// String คืนค่า string representation
// String returns the string representation
func (s AccountStatus) String() string {
	return string(s)
}
```

---

### 3.2 Domain Errors

**`internal/modules/pdpa/domain/errors/errors.go`**

```go
package domainerrors

import "errors"

// Sentinel errors สำหรับ PDPA domain
// Sentinel errors for PDPA domain
var (
	// Consent errors
	ErrConsentNotFound       = errors.New("consent record not found")
	ErrConsentAlreadyRevoked = errors.New("consent already revoked")
	ErrConsentExpired        = errors.New("consent expired")
	ErrInvalidPurpose        = errors.New("invalid consent purpose")
	ErrRevokeNotAllowed      = errors.New("cannot revoke non-granted consent")

	// DSAR errors
	ErrDSARNotFound          = errors.New("DSAR request not found")
	ErrDSARAlreadyProcessed  = errors.New("DSAR request already processed")
	ErrDSARNotPending        = errors.New("DSAR request is not in pending state")
	ErrOTPInvalid            = errors.New("OTP is invalid")
	ErrOTPExpired            = errors.New("OTP expired")
	ErrRateLimitExceeded     = errors.New("too many DSAR requests, please try again later")

	// Account errors
	ErrAccountNotFound             = errors.New("user account status not found")
	ErrAccountNotActive            = errors.New("account is not active")
	ErrAccountNotSuspended         = errors.New("account is not suspended")
	ErrAccountNotTerminated        = errors.New("account is not terminated")
	ErrDeletionNotConfirmed        = errors.New("deletion not confirmed for this account")
	ErrImmediateDeletionNotAllowed = errors.New("immediate deletion is not allowed for this account")
	ErrAlreadyDeleted              = errors.New("account data already deleted")

	// Policy errors
	ErrPolicyNotFound      = errors.New("privacy policy not found")
	ErrPolicyVersionExists = errors.New("privacy policy version already exists")

	// General
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInternal     = errors.New("internal server error")
)
```

---

### 3.3 Domain Event Metadata

**`internal/modules/pdpa/domain/event/event_metadata.go`**

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// EventMetadata ข้อมูล metadata สำหรับทุก event
// EventMetadata holds metadata for every event (tracing, correlation)
type EventMetadata struct {
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID string    `json:"correlation_id"` // trace ทั้ง flow
	CausationID   string    `json:"causation_id"`   // event ที่ทำให้เกิด event นี้
	TraceID       string    `json:"trace_id"`       // distributed tracing
	Timestamp     time.Time `json:"timestamp"`
	Version       int       `json:"version"`
	Source        string    `json:"source"`         // module name
}

// NewEventMetadata สร้าง EventMetadata ใหม่
// NewEventMetadata creates a new EventMetadata
func NewEventMetadata(correlationID, causationID, traceID string) EventMetadata {
	return EventMetadata{
		EventID:       uuid.New(),
		CorrelationID: correlationID,
		CausationID:   causationID,
		TraceID:       traceID,
		Timestamp:     time.Now().UTC(),
		Version:       1,
		Source:        "pdpa",
	}
}

// DomainEvent interface ที่ทุก event ต้อง implement
// DomainEvent interface that all events must implement
type DomainEvent interface {
	EventName() string
	GetMetadata() EventMetadata
	AggregateID() uuid.UUID
}

// Kafka topic constants
const (
	TopicConsentGranted     = "pdpa.consent.granted"
	TopicConsentRevoked     = "pdpa.consent.revoked"
	TopicDSARSubmitted      = "pdpa.dsar.submitted"
	TopicDSARCompleted      = "pdpa.dsar.completed"
	TopicAccountSuspended   = "pdpa.account.suspended"
	TopicAccountTerminated  = "pdpa.account.terminated"
	TopicDataDeletionReq    = "pdpa.data.deletion_requested"
	TopicDataDeleted        = "pdpa.data.deleted"
	TopicAuditTrail         = "pdpa.audit.trail"
	TopicLLMAnalysisRequest = "pdpa.llm.analysis.requested"
	TopicEmailNotification  = "pdpa.email.notification"
	TopicBlockchainRecord   = "pdpa.blockchain.record"
	TopicDLQPrefix          = "pdpa.dlq."
)
```

---

### 3.4 Consent Events

**`internal/modules/pdpa/domain/event/consent_events.go`**

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// ConsentGrantedEvent เหตุการณ์เมื่อผู้ใช้ให้ความยินยอม
// ConsentGrantedEvent is fired when a user grants consent
type ConsentGrantedEvent struct {
	Metadata  EventMetadata `json:"metadata"`
	ConsentID uuid.UUID     `json:"consent_id"`
	UserID    uuid.UUID     `json:"user_id"`
	SessionID string        `json:"session_id"`
	Purpose   string        `json:"purpose"`
	GrantedAt time.Time     `json:"granted_at"`
	ExpiresAt time.Time     `json:"expires_at"`
	IPAddress string        `json:"ip_address"`
	UserAgent string        `json:"user_agent"`
}

func (e ConsentGrantedEvent) EventName() string           { return TopicConsentGranted }
func (e ConsentGrantedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e ConsentGrantedEvent) AggregateID() uuid.UUID      { return e.ConsentID }

// ConsentRevokedEvent เหตุการณ์เมื่อผู้ใช้ถอนความยินยอม
// ConsentRevokedEvent is fired when a user revokes consent
type ConsentRevokedEvent struct {
	Metadata  EventMetadata `json:"metadata"`
	ConsentID uuid.UUID     `json:"consent_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Purpose   string        `json:"purpose"`
	RevokedAt time.Time     `json:"revoked_at"`
	IPAddress string        `json:"ip_address"`
	UserAgent string        `json:"user_agent"`
}

func (e ConsentRevokedEvent) EventName() string           { return TopicConsentRevoked }
func (e ConsentRevokedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e ConsentRevokedEvent) AggregateID() uuid.UUID      { return e.ConsentID }
```

---

### 3.5 DSAR Events

**`internal/modules/pdpa/domain/event/dsar_events.go`**

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// DSARSubmittedEvent เหตุการณ์เมื่อมีการส่งคำร้อง DSAR
// DSARSubmittedEvent is fired when a DSAR is submitted
type DSARSubmittedEvent struct {
	Metadata    EventMetadata `json:"metadata"`
	DSARID      uuid.UUID     `json:"dsar_id"`
	UserID      uuid.UUID     `json:"user_id"`
	RequestType string        `json:"request_type"`
	RequestedAt time.Time     `json:"requested_at"`
	IPAddress   string        `json:"ip_address"`
	UserAgent   string        `json:"user_agent"`
}

func (e DSARSubmittedEvent) EventName() string           { return TopicDSARSubmitted }
func (e DSARSubmittedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e DSARSubmittedEvent) AggregateID() uuid.UUID      { return e.DSARID }

// DSARCompletedEvent เหตุการณ์เมื่อคำร้อง DSAR เสร็จสิ้น
// DSARCompletedEvent is fired when a DSAR is completed
type DSARCompletedEvent struct {
	Metadata    EventMetadata `json:"metadata"`
	DSARID      uuid.UUID     `json:"dsar_id"`
	UserID      uuid.UUID     `json:"user_id"`
	RequestType string        `json:"request_type"`
	Status      string        `json:"status"`
	CompletedAt time.Time     `json:"completed_at"`
	DataPayload []byte        `json:"data_payload,omitempty"`
	Reason      string        `json:"reason,omitempty"`
}

func (e DSARCompletedEvent) EventName() string           { return TopicDSARCompleted }
func (e DSARCompletedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e DSARCompletedEvent) AggregateID() uuid.UUID      { return e.DSARID }
```

---

### 3.6 Account Events

**`internal/modules/pdpa/domain/event/account_events.go`**

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// AccountSuspendedEvent เหตุการณ์เมื่อบัญชีถูกระงับ
// AccountSuspendedEvent is fired when an account is suspended
type AccountSuspendedEvent struct {
	Metadata       EventMetadata `json:"metadata"`
	UserID         uuid.UUID     `json:"user_id"`
	SuspendedAt    time.Time     `json:"suspended_at"`
	RetentionYears int           `json:"retention_years"`
	Deadline       time.Time     `json:"retention_deadline"`
}

func (e AccountSuspendedEvent) EventName() string           { return TopicAccountSuspended }
func (e AccountSuspendedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e AccountSuspendedEvent) AggregateID() uuid.UUID      { return e.UserID }

// AccountTerminatedEvent เหตุการณ์เมื่อบัญชีถูกยกเลิก
// AccountTerminatedEvent is fired when an account is terminated
type AccountTerminatedEvent struct {
	Metadata     EventMetadata `json:"metadata"`
	UserID       uuid.UUID     `json:"user_id"`
	TerminatedAt time.Time     `json:"terminated_at"`
}

func (e AccountTerminatedEvent) EventName() string           { return TopicAccountTerminated }
func (e AccountTerminatedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e AccountTerminatedEvent) AggregateID() uuid.UUID      { return e.UserID }
```

---

### 3.7 Data Events

**`internal/modules/pdpa/domain/event/data_events.go`**

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// DataDeletionRequestedEvent เหตุการณ์เมื่อมีการขอให้ลบข้อมูล
// DataDeletionRequestedEvent is fired when data deletion is requested
type DataDeletionRequestedEvent struct {
	Metadata    EventMetadata `json:"metadata"`
	UserID      uuid.UUID     `json:"user_id"`
	Reason      string        `json:"reason"`
	RequestedBy string        `json:"requested_by"` // ADMIN or USER or SYSTEM
	RequestedAt time.Time     `json:"requested_at"`
}

func (e DataDeletionRequestedEvent) EventName() string           { return TopicDataDeletionReq }
func (e DataDeletionRequestedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e DataDeletionRequestedEvent) AggregateID() uuid.UUID      { return e.UserID }

// DataDeletedEvent เหตุการณ์เมื่อข้อมูลถูกลบสำเร็จ
// DataDeletedEvent is fired when data has been deleted
type DataDeletedEvent struct {
	Metadata   EventMetadata `json:"metadata"`
	UserID     uuid.UUID     `json:"user_id"`
	DeletedAt  time.Time     `json:"deleted_at"`
	DeleteType string        `json:"delete_type"` // IMMEDIATE, AUTO, DSAR
	Reason     string        `json:"reason,omitempty"`
}

func (e DataDeletedEvent) EventName() string           { return TopicDataDeleted }
func (e DataDeletedEvent) GetMetadata() EventMetadata  { return e.Metadata }
func (e DataDeletedEvent) AggregateID() uuid.UUID      { return e.UserID }

// PrivacyPolicyPublishedEvent เหตุการณ์เมื่อมีการเผยแพร่นโยบายใหม่
// PrivacyPolicyPublishedEvent is fired when a new policy is published
type PrivacyPolicyPublishedEvent struct {
	Metadata      EventMetadata `json:"metadata"`
	PolicyID      uuid.UUID     `json:"policy_id"`
	Version       string        `json:"version"`
	Title         string        `json:"title"`
	EffectiveDate time.Time     `json:"effective_date"`
	PublishedAt   time.Time     `json:"published_at"`
}

func (e PrivacyPolicyPublishedEvent) EventName() string          { return "pdpa.policy.published" }
func (e PrivacyPolicyPublishedEvent) GetMetadata() EventMetadata { return e.Metadata }
func (e PrivacyPolicyPublishedEvent) AggregateID() uuid.UUID     { return e.PolicyID }

// LLMAnalysisRequestedEvent เหตุการณ์ขอให้ LLM วิเคราะห์
// LLMAnalysisRequestedEvent is fired to request LLM analysis
type LLMAnalysisRequestedEvent struct {
	Metadata     EventMetadata `json:"metadata"`
	DSARID       uuid.UUID     `json:"dsar_id"`
	UserID       uuid.UUID     `json:"user_id"`
	AnalysisType string        `json:"analysis_type"`
	RequestedAt  time.Time     `json:"requested_at"`
}

func (e LLMAnalysisRequestedEvent) EventName() string          { return TopicLLMAnalysisRequest }
func (e LLMAnalysisRequestedEvent) GetMetadata() EventMetadata { return e.Metadata }
func (e LLMAnalysisRequestedEvent) AggregateID() uuid.UUID     { return e.DSARID }

// EmailNotificationEvent เหตุการณ์สำหรับส่ง email
// EmailNotificationEvent is fired to trigger an email
type EmailNotificationEvent struct {
	Metadata   EventMetadata          `json:"metadata"`
	To         string                 `json:"to"`
	Template   string                 `json:"template"`
	Subject    string                 `json:"subject"`
	Locale     string                 `json:"locale"`
	Variables  map[string]interface{} `json:"variables"`
	RefType    string                 `json:"ref_type"`
	RefID      string                 `json:"ref_id"`
	CreatedAt  time.Time              `json:"created_at"`
}

func (e EmailNotificationEvent) EventName() string          { return TopicEmailNotification }
func (e EmailNotificationEvent) GetMetadata() EventMetadata { return e.Metadata }
func (e EmailNotificationEvent) AggregateID() uuid.UUID     { return uuid.Nil }

// BlockchainRecordEvent เหตุการณ์สำหรับบันทึก audit log ลง blockchain
// BlockchainRecordEvent is fired to record an audit log on blockchain
type BlockchainRecordEvent struct {
	Metadata  EventMetadata `json:"metadata"`
	UserID    uuid.UUID     `json:"user_id"`
	Action    string        `json:"action"`
	DataHash  string        `json:"data_hash"` // SHA-256 hash ของข้อมูล
	CreatedAt time.Time     `json:"created_at"`
}

func (e BlockchainRecordEvent) EventName() string          { return TopicBlockchainRecord }
func (e BlockchainRecordEvent) GetMetadata() EventMetadata { return e.Metadata }
func (e BlockchainRecordEvent) AggregateID() uuid.UUID     { return e.UserID }
```

---

### 3.8 Entity: ConsentLog

**`internal/modules/pdpa/domain/entity/consent_log.go`**

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentLog Aggregate Root สำหรับบันทึกประวัติความยินยอม
// ConsentLog is the Aggregate Root for recording consent history
type ConsentLog struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	SessionID     string            `json:"session_id"`
	Purpose       vo.ConsentPurpose `json:"purpose"`
	Status        vo.ConsentStatus  `json:"status"`
	IPAddress     string            `json:"ip_address"`
	UserAgent     string            `json:"user_agent"`
	GrantedAt     time.Time         `json:"granted_at"`
	ExpiresAt     time.Time         `json:"expires_at"`
	RevokedAt     *time.Time        `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time        `json:"auto_deleted_at,omitempty"`

	// domain events ที่จะถูก publish (ไม่บันทึกใน DB)
	// domain events to be published (not persisted in DB)
	events []interface{}
}

// NewConsentLog สร้าง ConsentLog ใหม่พร้อม validation
// NewConsentLog creates a new ConsentLog with validation
func NewConsentLog(
	userID uuid.UUID,
	sessionID string,
	purpose vo.ConsentPurpose,
	ip, userAgent string,
) (*ConsentLog, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidInput
	}
	if !purpose.IsValid() {
		return nil, domainerrors.ErrInvalidPurpose
	}

	now := time.Now().UTC()
	return &ConsentLog{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
		Purpose:   purpose,
		Status:    vo.ConsentGranted,
		IPAddress: ip,
		UserAgent: userAgent,
		GrantedAt: now,
		ExpiresAt: now.AddDate(1, 0, 0), // อายุ 1 ปี / 1 year expiration
	}, nil
}

// Revoke ถอนความยินยอม
// Revoke revokes the consent
func (c *ConsentLog) Revoke(ip, userAgent string) error {
	if c.Status == vo.ConsentRevoked {
		return domainerrors.ErrConsentAlreadyRevoked
	}
	if c.Status == vo.ConsentExpired {
		return domainerrors.ErrConsentExpired
	}
	if c.Status != vo.ConsentGranted {
		return domainerrors.ErrRevokeNotAllowed
	}
	now := time.Now().UTC()
	c.Status = vo.ConsentRevoked
	c.RevokedAt = &now
	c.IPAddress = ip
	c.UserAgent = userAgent
	return nil
}

// MarkDeleted ทำเครื่องหมายว่าถูกลบ (auto-delete)
// MarkDeleted marks the consent as deleted
func (c *ConsentLog) MarkDeleted() {
	now := time.Now().UTC()
	c.Status = vo.ConsentDeleted
	c.AutoDeletedAt = &now
}

// MarkExpired ทำเครื่องหมายว่าหมดอายุ
// MarkExpired marks the consent as expired
func (c *ConsentLog) MarkExpired() {
	c.Status = vo.ConsentExpired
}

// IsActive ตรวจสอบว่ายังมีผลอยู่หรือไม่
// IsActive checks if the consent is currently active
func (c *ConsentLog) IsActive() bool {
	return c.Status == vo.ConsentGranted && c.ExpiresAt.After(time.Now().UTC())
}

// CanBeAutoDeleted ตรวจสอบว่าสามารถ auto-delete ได้หรือไม่
// CanBeAutoDeleted checks if the consent can be auto-deleted
func (c *ConsentLog) CanBeAutoDeleted(cutoff time.Time) bool {
	if c.Status != vo.ConsentRevoked && c.Status != vo.ConsentExpired {
		return false
	}
	return c.RevokedAt != nil && c.RevokedAt.Before(cutoff)
}

// AddEvent เพิ่ม domain event ที่จะถูก publish
// AddEvent adds a domain event to be published
func (c *ConsentLog) AddEvent(e interface{}) {
	c.events = append(c.events, e)
}

// Events คืนค่ารายการ domain events และ clear ข้อมูล
// Events returns and clears the domain events
func (c *ConsentLog) Events() []interface{} {
	events := c.events
	c.events = nil
	return events
}
```

---

### 3.9 Entity: DSARRequest

**`internal/modules/pdpa/domain/entity/dsar_request.go`**

```go
package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// DSARRequest Aggregate Root สำหรับคำร้องขอใช้สิทธิ์
// DSARRequest is the Aggregate Root for Data Subject Access Requests
type DSARRequest struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	RequestType     vo.DSARType     `json:"request_type"`
	Status          vo.DSARStatus   `json:"status"`
	RequestedAt     time.Time       `json:"requested_at"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	DataPayload     []byte          `json:"data_payload,omitempty"`
	RejectionReason string          `json:"rejection_reason,omitempty"`
	OTPHash         string          `json:"-"` // ไม่ส่งออก, เก็บแบบ hash
	OTPExpiredAt    time.Time       `json:"otp_expired_at,omitempty"`
	OTPVerified     bool            `json:"otp_verified"`
	IPAddress       string          `json:"ip_address,omitempty"`
	UserAgent       string          `json:"user_agent,omitempty"`

	events []interface{}
}

// otpTTL อายุของ OTP
// otpTTL is the OTP time-to-live
const otpTTL = 15 * time.Minute

// NewDSARRequest สร้าง DSARRequest ใหม่พร้อม generate OTP
// NewDSARRequest creates a new DSARRequest with generated OTP
// Returns: request, plainOTP (สำหรับส่ง email), error
func NewDSARRequest(
	userID uuid.UUID,
	requestType vo.DSARType,
	ip, userAgent string,
) (*DSARRequest, string, error) {
	if userID == uuid.Nil {
		return nil, "", domainerrors.ErrInvalidInput
	}
	if !requestType.IsValid() {
		return nil, "", domainerrors.ErrInvalidInput
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	req := &DSARRequest{
		ID:           uuid.New(),
		UserID:       userID,
		RequestType:  requestType,
		Status:       vo.DSARStatusPending,
		RequestedAt:  now,
		OTPHash:      hashOTP(otp),
		OTPExpiredAt: now.Add(otpTTL),
		IPAddress:    ip,
		UserAgent:    userAgent,
	}
	return req, otp, nil
}

// VerifyOTP ตรวจสอบ OTP
// VerifyOTP verifies the OTP
func (d *DSARRequest) VerifyOTP(otp string) error {
	if d.OTPVerified {
		return nil // already verified
	}
	if time.Now().UTC().After(d.OTPExpiredAt) {
		return domainerrors.ErrOTPExpired
	}
	if hashOTP(otp) != d.OTPHash {
		return domainerrors.ErrOTPInvalid
	}
	d.OTPVerified = true
	return nil
}

// MarkProcessing เปลี่ยนสถานะเป็นกำลังดำเนินการ
// MarkProcessing transitions to processing state
func (d *DSARRequest) MarkProcessing() error {
	if d.Status != vo.DSARStatusPending {
		return domainerrors.ErrDSARNotPending
	}
	if !d.OTPVerified {
		return domainerrors.ErrOTPInvalid
	}
	d.Status = vo.DSARStatusProcessing
	return nil
}

// MarkCompleted เปลี่ยนสถานะเป็นเสร็จสิ้น
// MarkCompleted transitions to completed state
func (d *DSARRequest) MarkCompleted(payload []byte) error {
	if d.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	now := time.Now().UTC()
	d.Status = vo.DSARStatusCompleted
	d.CompletedAt = &now
	d.DataPayload = payload
	return nil
}

// MarkRejected ปฏิเสธคำร้อง
// MarkRejected rejects the request
func (d *DSARRequest) MarkRejected(reason string) error {
	if d.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = vo.DSARStatusRejected
	d.RejectionReason = reason
	return nil
}

// AddEvent เพิ่ม domain event
// AddEvent adds a domain event
func (d *DSARRequest) AddEvent(e interface{}) {
	d.events = append(d.events, e)
}

// Events คืนค่าและ clear domain events
// Events returns and clears domain events
func (d *DSARRequest) Events() []interface{} {
	events := d.events
	d.events = nil
	return events
}

// generateOTP สร้าง OTP 6 หลัก
// generateOTP generates a 6-digit OTP
func generateOTP() (string, error) {
	const digits = "0123456789"
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = digits[int(b[i])%10]
	}
	return string(b), nil
}

// hashOTP hash OTP ด้วย SHA-256
// hashOTP hashes the OTP with SHA-256
func hashOTP(otp string) string {
	h := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(h[:])
}
```

---

### 3.10 Entity: UserAccountStatus

**`internal/modules/pdpa/domain/entity/user_account_status.go`**

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// UserAccountStatus Aggregate Root สำหรับสถานะบัญชีผู้ใช้
// UserAccountStatus is the Aggregate Root for account lifecycle status
type UserAccountStatus struct {
	ID                  uuid.UUID         `json:"id"`
	UserID              uuid.UUID         `json:"user_id"`
	Status              vo.AccountStatus  `json:"status"`
	SuspendedAt         *time.Time        `json:"suspended_at,omitempty"`
	TerminatedAt        *time.Time        `json:"terminated_at,omitempty"`
	DeletionConfirmedAt *time.Time        `json:"deletion_confirmed_at,omitempty"`
	RetentionDeadline   *time.Time        `json:"retention_deadline,omitempty"`
	AutoDeletedAt       *time.Time        `json:"auto_deleted_at,omitempty"`
	UpdatedAt           time.Time         `json:"updated_at"`

	events []interface{}
}

// NewUserAccountStatus สร้าง UserAccountStatus ใหม่
// NewUserAccountStatus creates a new UserAccountStatus
func NewUserAccountStatus(userID uuid.UUID) (*UserAccountStatus, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidInput
	}
	return &UserAccountStatus{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    vo.AccountActive,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

// Suspend ระงับบัญชี (Retain for retentionYears)
// Suspend suspends the account and sets retention deadline
func (u *UserAccountStatus) Suspend(suspendedAt time.Time, retentionYears int) error {
	if u.Status == vo.AccountDeleted {
		return domainerrors.ErrAlreadyDeleted
	}
	if u.Status == vo.AccountTerminated {
		return domainerrors.ErrAccountNotActive
	}
	if retentionYears < 1 {
		retentionYears = 1
	}
	deadline := suspendedAt.AddDate(retentionYears, 0, 0)
	u.Status = vo.AccountSuspended
	u.SuspendedAt = &suspendedAt
	u.RetentionDeadline = &deadline
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// Terminate ยกเลิกบัญชี (ต้องยืนยันการลบ)
// Terminate terminates the account (requires deletion confirmation)
func (u *UserAccountStatus) Terminate(terminatedAt time.Time) error {
	if u.Status == vo.AccountDeleted {
		return domainerrors.ErrAlreadyDeleted
	}
	u.Status = vo.AccountTerminated
	u.TerminatedAt = &terminatedAt
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// ConfirmDeletion ยืนยันการลบข้อมูล
// ConfirmDeletion confirms the deletion
func (u *UserAccountStatus) ConfirmDeletion() error {
	if u.Status != vo.AccountTerminated {
		return domainerrors.ErrAccountNotTerminated
	}
	now := time.Now().UTC()
	u.DeletionConfirmedAt = &now
	u.UpdatedAt = now
	return nil
}

// MarkDeleted ทำเครื่องหมายว่าลบข้อมูลแล้ว
// MarkDeleted marks the account as deleted
func (u *UserAccountStatus) MarkDeleted() {
	now := time.Now().UTC()
	u.Status = vo.AccountDeleted
	u.AutoDeletedAt = &now
	u.UpdatedAt = now
}

// IsReadyForAutoDeletion ตรวจสอบว่าพร้อม auto-delete หรือไม่
// IsReadyForAutoDeletion checks if ready for auto-deletion
func (u *UserAccountStatus) IsReadyForAutoDeletion(now time.Time) bool {
	if u.Status != vo.AccountSuspended {
		return false
	}
	if u.RetentionDeadline == nil {
		return false
	}
	return !u.RetentionDeadline.After(now)
}

// CanImmediateDelete ตรวจสอบว่าสามารถลบทันทีได้หรือไม่
// CanImmediateDelete checks if immediate deletion is allowed
func (u *UserAccountStatus) CanImmediateDelete() bool {
	return u.Status == vo.AccountTerminated && u.DeletionConfirmedAt != nil
}

// AddEvent เพิ่ม domain event
// AddEvent adds a domain event
func (u *UserAccountStatus) AddEvent(e interface{}) {
	u.events = append(u.events, e)
}

// Events คืนค่าและ clear events
// Events returns and clears events
func (u *UserAccountStatus) Events() []interface{} {
	events := u.events
	u.events = nil
	return events
}
```

---

### 3.11 Entity: AuditTrail

**`internal/modules/pdpa/domain/entity/audit_trail.go`**

```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditTrail Entity สำหรับบันทึกการดำเนินการ (ไม่ใช่ aggregate root)
// AuditTrail is an Entity for recording actions (not an aggregate root)
type AuditTrail struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	Action    string      `json:"action"`
	Details   interface{} `json:"details,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// NewAuditTrail สร้าง AuditTrail ใหม่ (ไม่มี metadata)
// NewAuditTrail creates a new AuditTrail without request metadata
func NewAuditTrail(userID *uuid.UUID, action string, details interface{}) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now().UTC(),
	}
}

// NewAuditTrailWithMeta สร้าง AuditTrail พร้อม request metadata
// NewAuditTrailWithMeta creates a new AuditTrail with request metadata
func NewAuditTrailWithMeta(
	userID *uuid.UUID,
	action string,
	details interface{},
	ip, userAgent string,
) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now().UTC(),
	}
}

// Audit action constants
const (
	ActionConsentGranted       = "CONSENT_GRANTED"
	ActionConsentRevoked       = "CONSENT_REVOKED"
	ActionDSARSubmitted        = "DSAR_SUBMITTED"
	ActionDSARCompleted        = "DSAR_COMPLETED"
	ActionDSARRejected         = "DSAR_REJECTED"
	ActionOTPVerified          = "OTP_VERIFIED"
	ActionAccountSuspended     = "ACCOUNT_SUSPENDED"
	ActionAccountTerminated    = "ACCOUNT_TERMINATED"
	ActionDeletionConfirmed    = "DELETION_CONFIRMED"
	ActionImmediateDeletion    = "IMMEDIATE_DELETION"
	ActionAutoDeletion         = "AUTO_DELETION"
	ActionPrivacyPolicyPublish = "PRIVACY_POLICY_PUBLISHED"
)
```

---

### 3.12 Entity: PrivacyPolicy

**`internal/modules/pdpa/domain/entity/privacy_policy.go`**

```go
package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
)

// PrivacyPolicy Entity สำหรับนโยบายความเป็นส่วนตัว
// PrivacyPolicy is an Entity for privacy policy versions
type PrivacyPolicy struct {
	ID            uuid.UUID `json:"id"`
	Version       string    `json:"version"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	EffectiveDate time.Time `json:"effective_date"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewPrivacyPolicy สร้าง PrivacyPolicy ใหม่พร้อม validation
// NewPrivacyPolicy creates a new PrivacyPolicy with validation
func NewPrivacyPolicy(version, title, content string, effectiveDate time.Time) (*PrivacyPolicy, error) {
	version = strings.TrimSpace(version)
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if version == "" || title == "" || content == "" {
		return nil, domainerrors.ErrInvalidInput
	}
	if effectiveDate.IsZero() {
		return nil, domainerrors.ErrInvalidInput
	}

	now := time.Now().UTC()
	return &PrivacyPolicy{
		ID:            uuid.New(),
		Version:       version,
		Title:         title,
		Content:       content,
		EffectiveDate: effectiveDate,
		IsActive:      false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Activate เปิดใช้งานนโยบายนี้
// Activate activates this policy version
func (p *PrivacyPolicy) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now().UTC()
}

// Deactivate ปิดการใช้งานนโยบายนี้
// Deactivate deactivates this policy version
func (p *PrivacyPolicy) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now().UTC()
}
```

---

## สรุปส่วนที่ 1 (Domain Layer)

สร้างเสร็จแล้ว 12 ไฟล์ใน Domain Layer:

| # | ไฟล์ | หน้าที่ |
|---|------|--------|
| 1 | `value_object/consent_purpose.go` | Purpose VO + validation |
| 2 | `value_object/consent_status.go` | Consent status VO |
| 3 | `value_object/dsar_type.go` | DSAR type VO |
| 4 | `value_object/dsar_status.go` | DSAR status VO |
| 5 | `value_object/account_status.go` | Account status VO |
| 6 | `errors/errors.go` | Sentinel errors |
| 7 | `event/event_metadata.go` | Metadata + Kafka topics |
| 8 | `event/consent_events.go` | Consent domain events |
| 9 | `event/dsar_events.go` | DSAR domain events |
| 10 | `event/account_events.go` | Account domain events |
| 11 | `event/data_events.go` | Data/policy/email events |
| 12 | `entity/consent_log.go` | ConsentLog aggregate |
| 13 | `entity/dsar_request.go` | DSARRequest aggregate (with OTP) |
| 14 | `entity/user_account_status.go` | Account status aggregate |
| 15 | `entity/audit_trail.go` | AuditTrail entity |
| 16 | `entity/privacy_policy.go` | PrivacyPolicy entity |

**หมายเหตุสำคัญ:**
- Domain Layer ไม่ import gorm/gin/chi/sarama/redis เลย ✅
- ทุก Entity มี constructor `New{Entity}` + validation ✅
- เปลี่ยน state ผ่าน behavior methods เท่านั้น (ไม่มี setter) ✅
- Domain Events ฝังใน entity ผ่าน `AddEvent()` / `Events()` ✅
- OTP ถูก hash ด้วย SHA-256 (ไม่เก็บ plain text) ✅

---

**พร้อมดำเนินการส่วนถัดไปหรือไม่?**

ส่วนถัดไปคือ:
- **3.2 Repository Interfaces** (6 ไฟล์) — interfaces เท่านั้น
- **3.3 Domain Services** (3 ไฟล์)
- **3.4 Application Layer — Commands** (10 ไฟล์)
- **3.5 Application Layer — Queries** (4 ไฟล์)
- **3.6 Application Layer — Event Handlers** (8 ไฟล์)
- **3.7 Infrastructure Layer** (Postgres, Redis, Kafka, Outbox, External clients)
- **3.8 Interface Layer** (HTTP handlers, routes, DTOs)
- **3.9 `module.go`** (DI wiring)
- **4. Migrations**
- **5. Docker Compose**
- **6. Testing Strategy**
- **7. Deployment Considerations**

พิมพ์ **"ต่อ"** เพื่อไปส่วนที่ 2 (Repository Interfaces + Domain Services) ครับ