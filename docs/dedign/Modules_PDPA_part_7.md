
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


# เติมไฟล์ที่เหลือให้ครบ — Production-Grade PDPA Module

ผมจะเติมไฟล์ที่เหลือตาม tree ที่คุณให้ พร้อมปรับปรุงตาม recommendations ที่แนะนำไว้ (string VOs, outbox, idempotency, typed events)

---

## 📁 1. Domain Layer — Value Objects

### 1.1 `domain/value_object/consent_purpose.go`

```go
package valueobject

import (
	"fmt"
	"strings"
)

// ConsentPurpose วัตถุประสงค์การขอความยินยอม (string-based สำหรับ JSON contract)
// ConsentPurpose represents consent purpose (string-based for JSON API)
type ConsentPurpose string

const (
	ConsentPurposeNecessary          ConsentPurpose = "NECESSARY"
	ConsentPurposeAnalytics          ConsentPurpose = "ANALYTICS"
	ConsentPurposeMarketing          ConsentPurpose = "MARKETING"
	ConsentPurposeAccountSystem      ConsentPurpose = "ACCOUNT_SYSTEM"
	ConsentPurposeUsageLogs          ConsentPurpose = "USAGE_LOGS"
	ConsentPurposeTransactionHistory ConsentPurpose = "TRANSACTION_HISTORY"
)

func ParseConsentPurpose(s string) (ConsentPurpose, error) {
	p := ConsentPurpose(strings.ToUpper(strings.TrimSpace(s)))
	if !p.IsValid() {
		return "", fmt.Errorf("invalid consent purpose: %q", s)
	}
	return p, nil
}

func (c ConsentPurpose) IsValid() bool {
	switch c {
	case ConsentPurposeNecessary, ConsentPurposeAnalytics, ConsentPurposeMarketing,
		ConsentPurposeAccountSystem, ConsentPurposeUsageLogs, ConsentPurposeTransactionHistory:
		return true
	}
	return false
}

func (c ConsentPurpose) IsRequired() bool { return c == ConsentPurposeNecessary }
func (c ConsentPurpose) String() string   { return string(c) }

func (c ConsentPurpose) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(c) + `"`), nil
}

func (c *ConsentPurpose) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	p, err := ParseConsentPurpose(s)
	if err != nil {
		return err
	}
	*c = p
	return nil
}
```

### 1.2 `domain/value_object/consent_status.go`

```go
package valueobject

type ConsentStatus string

const (
	ConsentGranted ConsentStatus = "GRANTED"
	ConsentRevoked ConsentStatus = "REVOKED"
	ConsentExpired ConsentStatus = "EXPIRED"
	ConsentDeleted ConsentStatus = "DELETED"
)

func (c ConsentStatus) IsValid() bool {
	switch c {
	case ConsentGranted, ConsentRevoked, ConsentExpired, ConsentDeleted:
		return true
	}
	return false
}

func (c ConsentStatus) IsActive() bool { return c == ConsentGranted }
func (c ConsentStatus) String() string { return string(c) }

func (c ConsentStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(c) + `"`), nil
}
```

### 1.3 `domain/value_object/dsar_type.go`

```go
package valueobject

import (
	"fmt"
	"strings"
)

type DSARType string

const (
	DSARTypeAccess          DSARType = "ACCESS"
	DSARTypeErasure         DSARType = "ERASURE"
	DSARTypeWithdrawConsent DSARType = "WITHDRAW_CONSENT"
)

func ParseDSARType(s string) (DSARType, error) {
	t := DSARType(strings.ToUpper(strings.TrimSpace(s)))
	if !t.IsValid() {
		return "", fmt.Errorf("invalid dsar type: %q", s)
	}
	return t, nil
}

func (d DSARType) IsValid() bool {
	switch d {
	case DSARTypeAccess, DSARTypeErasure, DSARTypeWithdrawConsent:
		return true
	}
	return false
}

func (d DSARType) String() string { return string(d) }

func (d DSARType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(d) + `"`), nil
}
```

### 1.4 `domain/value_object/dsar_status.go`

```go
package valueobject

type DSARStatus string

const (
	DSARStatusPending    DSARStatus = "PENDING"
	DSARStatusProcessing DSARStatus = "PROCESSING"
	DSARStatusCompleted  DSARStatus = "COMPLETED"
	DSARStatusRejected   DSARStatus = "REJECTED"
)

func (s DSARStatus) IsValid() bool {
	switch s {
	case DSARStatusPending, DSARStatusProcessing, DSARStatusCompleted, DSARStatusRejected:
		return true
	}
	return false
}

func (s DSARStatus) IsTerminal() bool {
	return s == DSARStatusCompleted || s == DSARStatusRejected
}

func (s DSARStatus) String() string { return string(s) }

func (s DSARStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(s) + `"`), nil
}
```

### 1.5 `domain/value_object/account_status.go`

```go
package valueobject

type AccountStatus string

const (
	AccountActive     AccountStatus = "ACTIVE"
	AccountSuspended  AccountStatus = "SUSPENDED"
	AccountTerminated AccountStatus = "TERMINATED"
	AccountDeleted    AccountStatus = "DELETED"
)

func (a AccountStatus) IsValid() bool {
	switch a {
	case AccountActive, AccountSuspended, AccountTerminated, AccountDeleted:
		return true
	}
	return false
}

func (a AccountStatus) IsActive() bool     { return a == AccountActive }
func (a AccountStatus) IsSuspended() bool  { return a == AccountSuspended }
func (a AccountStatus) IsTerminated() bool { return a == AccountTerminated }
func (a AccountStatus) IsDeleted() bool    { return a == AccountDeleted }
func (a AccountStatus) String() string     { return string(a) }

func (a AccountStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(a) + `"`), nil
}
```

---

## 📁 2. Domain Layer — Events

### 2.1 `domain/event/event_metadata.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// EventMetadata metadata ที่ทุก event ต้องมี (tracing + idempotency)
type EventMetadata struct {
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id,omitempty"`
	TraceID       string    `json:"trace_id,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
	Version       int       `json:"version"`
	Source        string    `json:"source"`
}

func NewMetadata(correlationID, causationID, traceID string) EventMetadata {
	return EventMetadata{
		EventID:       uuid.New(),
		CorrelationID: correlationID,
		CausationID:   causationID,
		TraceID:       traceID,
		OccurredAt:    time.Now().UTC(),
		Version:       1,
		Source:        "pdpa",
	}
}

// Envelope โครงสร้างมาตรฐานของ event บน Kafka
type Envelope struct {
	EventType string        `json:"event_type"`
	Metadata  EventMetadata `json:"metadata"`
	Payload   interface{}   `json:"payload"`
}

func NewEnvelope(eventType string, meta EventMetadata, payload interface{}) Envelope {
	return Envelope{EventType: eventType, Metadata: meta, Payload: payload}
}

// Kafka Topics
const (
	TopicConsentGranted    = "pdpa.consent.granted"
	TopicConsentRevoked    = "pdpa.consent.revoked"
	TopicDSARSubmitted     = "pdpa.dsar.submitted"
	TopicDSARProcessing    = "pdpa.dsar.processing"
	TopicDSAROTPVerified   = "pdpa.dsar.otp_verified"
	TopicDSARCompleted     = "pdpa.dsar.completed"
	TopicDSARRejected      = "pdpa.dsar.rejected"
	TopicAccountSuspended  = "pdpa.account.suspended"
	TopicAccountTerminated = "pdpa.account.terminated"
	TopicDeletionRequested = "pdpa.data.deletion_requested"
	TopicDataDeleted       = "pdpa.data.deleted"
	TopicLLMAnalysis       = "pdpa.llm.analysis.requested"
	TopicEmailNotification = "pdpa.email.notification"
	TopicBlockchainRecord  = "pdpa.blockchain.record"
	TopicPolicyPublished   = "pdpa.policy.published"
	TopicDLQPrefix         = "pdpa.dlq."
)
```

### 2.2 `domain/event/consent_events.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

type ConsentGrantedPayload struct {
	ConsentID uuid.UUID `json:"consent_id"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	GrantedAt time.Time `json:"granted_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

type ConsentRevokedPayload struct {
	ConsentID uuid.UUID `json:"consent_id"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	RevokedAt time.Time `json:"revoked_at"`
}
```

### 2.3 `domain/event/dsar_events.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

type DSARSubmittedPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	RequestType string    `json:"request_type"`
	RequestedAt time.Time `json:"requested_at"`
}

type DSARProcessingPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	ProcessedAt time.Time `json:"processed_at"`
}

type DSAROTPVerifiedPayload struct {
	DSARID     uuid.UUID `json:"dsar_id"`
	UserID     uuid.UUID `json:"user_id"`
	VerifiedAt time.Time `json:"verified_at"`
}

type DSARCompletedPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	DataHash    string    `json:"data_hash,omitempty"`
	TxHash      string    `json:"tx_hash,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}

type DSARRejectedPayload struct {
	DSARID     uuid.UUID `json:"dsar_id"`
	UserID     uuid.UUID `json:"user_id"`
	Reason     string    `json:"reason"`
	RejectedAt time.Time `json:"rejected_at"`
}
```

### 2.4 `domain/event/account_events.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

type AccountSuspendedPayload struct {
	UserID            uuid.UUID `json:"user_id"`
	Status            string    `json:"status"`
	SuspendedAt       time.Time `json:"suspended_at"`
	RetentionDeadline time.Time `json:"retention_deadline"`
	RetentionYears    int       `json:"retention_years"`
}

type AccountTerminatedPayload struct {
	UserID       uuid.UUID `json:"user_id"`
	Status       string    `json:"status"`
	TerminatedAt time.Time `json:"terminated_at"`
}
```

### 2.5 `domain/event/data_events.go`

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

type DataDeletionRequestedPayload struct {
	UserID              uuid.UUID `json:"user_id"`
	DeletionConfirmedAt time.Time `json:"deletion_confirmed_at"`
	ImmediateDeletion   bool      `json:"immediate_deletion"`
}

type DataDeletedPayload struct {
	UserID     uuid.UUID `json:"user_id"`
	DeleteType string    `json:"delete_type"` // IMMEDIATE | AUTO | DSAR
	Reason     string    `json:"reason,omitempty"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type EmailNotificationPayload struct {
	Template  string                 `json:"template"`
	Locale    string                 `json:"locale"`
	To        string                 `json:"to,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type BlockchainRecordPayload struct {
	UserID   string `json:"user_id"`
	Action   string `json:"action"`
	DataHash string `json:"data_hash"`
}

type PolicyPublishedPayload struct {
	PolicyID      uuid.UUID `json:"policy_id"`
	Version       string    `json:"version"`
	Title         string    `json:"title"`
	EffectiveDate time.Time `json:"effective_date"`
	PublishedAt   time.Time `json:"published_at"`
}
```

---

## 📁 3. Domain Layer — Repository Interfaces

### 3.1 `domain/repository/consent_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentRepository interface {
	Save(ctx context.Context, c *entity.ConsentLog) error
	Update(ctx context.Context, c *entity.ConsentLog) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ConsentLog, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error)
	FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error)
	FindActiveByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error)
	FindReadyForAutoDeletion(ctx context.Context, cutoff time.Time, limit int) ([]entity.ConsentLog, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error)
}
```

### 3.2 `domain/repository/dsar_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DSARRepository interface {
	Save(ctx context.Context, d *entity.DSARRequest) error
	Update(ctx context.Context, d *entity.DSARRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error)
	FindByStatus(ctx context.Context, status valueobject.DSARStatus, limit, offset int) ([]entity.DSARRequest, error)
	CountByUserAndTypeInPeriod(ctx context.Context, userID uuid.UUID, t valueobject.DSARType, from time.Time) (int64, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

### 3.3 `domain/repository/account_status_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type UserAccountStatusRepository interface {
	Save(ctx context.Context, s *entity.UserAccountStatus) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error)
	FindReadyForAutoDeletion(ctx context.Context, now time.Time, limit int) ([]entity.UserAccountStatus, error)
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}
```

### 3.4 `domain/repository/audit_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type AuditFilter struct {
	UserID   *uuid.UUID
	Action   string
	FromTime *time.Time
	ToTime   *time.Time
	Limit    int
	Offset   int
}

type AuditRepository interface {
	Save(ctx context.Context, a *entity.AuditTrail) error
	FindByFilter(ctx context.Context, f AuditFilter) ([]entity.AuditTrail, error)
	CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error)
}
```

### 3.5 `domain/repository/policy_repository.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type PolicyRepository interface {
	Save(ctx context.Context, p *entity.PrivacyPolicy) error
	Update(ctx context.Context, p *entity.PrivacyPolicy) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PrivacyPolicy, error)
	FindActive(ctx context.Context) (*entity.PrivacyPolicy, error)
	FindAll(ctx context.Context) ([]entity.PrivacyPolicy, error)
	FindByVersion(ctx context.Context, version string) (*entity.PrivacyPolicy, error)
	DeactivateAll(ctx context.Context) error
}
```

### 3.6 `domain/repository/outbox_repository.go`

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	OutboxPending   = "PENDING"
	OutboxPublished = "PUBLISHED"
	OutboxFailed    = "FAILED"
)

type OutboxEvent struct {
	ID            uuid.UUID  `json:"id"`
	AggregateType string     `json:"aggregate_type"`
	AggregateID   uuid.UUID  `json:"aggregate_id"`
	EventType     string     `json:"event_type"`
	Topic         string     `json:"topic"`
	Payload       []byte     `json:"payload"`
	Status        string     `json:"status"`
	RetryCount    int        `json:"retry_count"`
	LastError     string     `json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
}

type OutboxRepository interface {
	Save(ctx context.Context, e *OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error)
}

// UnitOfWork abstraction
type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### 3.7 `domain/repository/cache.go`

```go
package repository

import (
	"context"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type Cache interface {
	Set(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose, status valueobject.ConsentStatus) error
	Get(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (valueobject.ConsentStatus, error)
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}
```

---

## 📁 4. Domain Layer — Services

### 4.1 `domain/service/consent_validator.go`

```go
package service

import (
	"fmt"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentValidator struct{}

func NewConsentValidator() *ConsentValidator { return &ConsentValidator{} }

func (v *ConsentValidator) ValidatePurposes(purposes map[string]bool) ([]valueobject.ConsentPurpose, error) {
	if len(purposes) == 0 {
		return nil, fmt.Errorf("no purposes provided")
	}
	result := make([]valueobject.ConsentPurpose, 0, len(purposes))
	hasNecessary := false
	for raw, granted := range purposes {
		p, err := valueobject.ParseConsentPurpose(raw)
		if err != nil {
			return nil, err
		}
		if p.IsRequired() {
			hasNecessary = true
			result = append(result, p)
			continue
		}
		if granted {
			result = append(result, p)
		}
	}
	if !hasNecessary {
		return nil, fmt.Errorf("mandatory purpose %q is required", valueobject.ConsentPurposeNecessary)
	}
	return result, nil
}

func (v *ConsentValidator) ValidateRevoke(p valueobject.ConsentPurpose) error {
	if !p.IsValid() {
		return fmt.Errorf("invalid purpose: %s", p)
	}
	if p.IsRequired() {
		return fmt.Errorf("cannot revoke mandatory purpose: %s", p)
	}
	return nil
}
```

### 4.2 `domain/service/deletion_policy_service.go`

```go
package service

import (
	"time"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type RetentionPolicy struct {
	SuspendedRetentionYears int
	RevokedConsentDays      int
}

func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{SuspendedRetentionYears: 1, RevokedConsentDays: 365}
}

type DeletionPolicyService struct {
	policy RetentionPolicy
}

func NewDeletionPolicyService(policy RetentionPolicy) *DeletionPolicyService {
	if policy.SuspendedRetentionYears < 1 {
		policy.SuspendedRetentionYears = 1
	}
	if policy.RevokedConsentDays < 1 {
		policy.RevokedConsentDays = 365
	}
	return &DeletionPolicyService{policy: policy}
}

func (s *DeletionPolicyService) CanImmediateDeletion(st *entity.UserAccountStatus) bool {
	if st == nil {
		return false
	}
	return st.CanImmediateDelete()
}

func (s *DeletionPolicyService) IsReadyForAutoDeletion(st *entity.UserAccountStatus, now time.Time) bool {
	if st == nil {
		return false
	}
	return st.IsReadyForAutoDeletion(now)
}

func (s *DeletionPolicyService) CalculateRetentionDeadline(suspendedAt time.Time) time.Time {
	return suspendedAt.AddDate(s.policy.SuspendedRetentionYears, 0, 0)
}

func (s *DeletionPolicyService) ConsentAutoDeleteCutoff(now time.Time) time.Time {
	return now.AddDate(0, 0, -s.policy.RevokedConsentDays)
}

func (s *DeletionPolicyService) Policy() RetentionPolicy { return s.policy }
```

### 4.3 `domain/service/anonymization_service.go`

```go
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type AnonymizationService struct{ salt string }

func NewAnonymizationService(salt string) *AnonymizationService {
	return &AnonymizationService{salt: salt}
}

func (s *AnonymizationService) AnonymizeEmail(email string) string {
	if email == "" {
		return ""
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return s.hash(email)
	}
	return fmt.Sprintf("%s@%s", s.hash(parts[0])[:12], parts[1])
}

func (s *AnonymizationService) AnonymizePhone(phone string) string {
	if len(phone) <= 4 {
		return strings.Repeat("*", len(phone))
	}
	return strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
}

func (s *AnonymizationService) AnonymizeIP(ip string) string {
	if ip == "" {
		return ""
	}
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			parts[3] = "0"
			return strings.Join(parts, ".")
		}
	}
	return s.hash(ip)[:16]
}

func (s *AnonymizationService) hash(input string) string {
	h := sha256.Sum256([]byte(s.salt + input))
	return hex.EncodeToString(h[:])
}
```

### 4.4 `domain/service/blockchain_service.go`

```go
package service

import "context"

// BlockchainService abstraction สำหรับบันทึก hash ลง blockchain
type BlockchainService interface {
	RecordHash(ctx context.Context, hash string) (txHash string, err error)
}
```

---

## 📁 5. Application Layer — DTOs

### 5.1 `application/dto/dto.go`

```go
package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== Consent ====================
type GrantConsentRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	Purpose       string    `json:"purpose" validate:"required"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type RevokeConsentRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	Purpose       string    `json:"purpose" validate:"required"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type ConsentStatusResponse struct {
	UserID    uuid.UUID  `json:"user_id"`
	Purpose   string     `json:"purpose"`
	Status    string     `json:"status"`
	IsActive  bool       `json:"is_active"`
	GrantedAt time.Time  `json:"granted_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type ConsentStatusListResponse struct {
	Consents []ConsentStatusResponse `json:"consents"`
	Total    int64                   `json:"total"`
}

// ==================== DSAR ====================
type SubmitDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	RequestType   string    `json:"request_type" validate:"required"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type VerifyOTPRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid"`
	OTPCode       string    `json:"otp_code" validate:"required,len=6,numeric"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type ProcessDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type CompleteDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid"`
	ResponseData  []byte    `json:"response_data" validate:"required"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type RejectDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid"`
	Reason        string    `json:"reason" validate:"required,min=1,max=500"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type DSARStatusResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	RequestType     string     `json:"request_type"`
	Status          string     `json:"status"`
	RequestedAt     time.Time  `json:"requested_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
}

type DSARStatusListResponse struct {
	Requests []DSARStatusResponse `json:"requests"`
	Total    int64                `json:"total"`
}

// ==================== Account ====================
type SuspendAccountRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required,uuid"`
	RetentionYears int       `json:"retention_years" validate:"required,gte=1,lte=10"`
	CorrelationID  string    `json:"-"`
	TraceID        string    `json:"-"`
}

type TerminateAccountRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type ConfirmDeletionRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	Reason        string    `json:"reason,omitempty"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type UserAccountStatusResponse struct {
	UserID              uuid.UUID  `json:"user_id"`
	Status              string     `json:"status"`
	RetentionDeadline   *time.Time `json:"retention_deadline,omitempty"`
	DeletionConfirmedAt *time.Time `json:"deletion_confirmed_at,omitempty"`
	AutoDeletedAt       *time.Time `json:"auto_deleted_at,omitempty"`
}

// ==================== Policy ====================
type PublishPolicyRequest struct {
	Version       string    `json:"version" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Content       string    `json:"content" validate:"required"`
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

type PrivacyPolicyResponse struct {
	ID            uuid.UUID `json:"id"`
	Version       string    `json:"version"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	EffectiveDate time.Time `json:"effective_date"`
	IsActive      bool      `json:"is_active"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ==================== Admin ====================
type AdminReportResponse struct {
	From             time.Time        `json:"from"`
	To               time.Time        `json:"to"`
	ConsentByPurpose map[string]int64 `json:"consent_by_purpose"`
	ActionsByType    map[string]int64 `json:"actions_by_type"`
	GeneratedAt      time.Time        `json:"generated_at"`
}

type AuditTrailItem struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	Action    string      `json:"action"`
	Details   interface{} `json:"details,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}
```

---

## 📁 6. Application Layer — Commands

### 6.1 `application/command/record_consent.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/logger"
)

type RecordConsentHandler struct {
	consentRepo repository.ConsentRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	validator   *service.ConsentValidator
	clock       func() time.Time
	logger      logger.Logger
}

func NewRecordConsentHandler(
	consentRepo repository.ConsentRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	validator *service.ConsentValidator,
	clock func() time.Time,
	logger logger.Logger,
) *RecordConsentHandler {
	if clock == nil {
		clock = time.Now
	}
	return &RecordConsentHandler{
		consentRepo: consentRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		validator:   validator,
		clock:       clock,
		logger:      logger,
	}
}

func (h *RecordConsentHandler) Handle(ctx context.Context, req dto.GrantConsentRequest) error {
	purpose, err := valueobject.ParseConsentPurpose(req.Purpose)
	if err != nil {
		return err
	}
	now := h.clock()

	consent, err := entity.NewConsentLog(req.UserID, purpose, now, req.IPAddress, req.UserAgent)
	if err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicConsentGranted, meta, event.ConsentGrantedPayload{
		ConsentID: consent.ID,
		UserID:    consent.UserID,
		Purpose:   consent.Purpose.String(),
		GrantedAt: consent.ConsentedAt,
		ExpiresAt: consent.ExpiresAt,
		IPAddress: consent.IPAddress,
		UserAgent: consent.UserAgent,
	})
	body, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.consentRepo.Save(txCtx, consent); err != nil {
			return fmt.Errorf("save consent: %w", err)
		}
		if err := h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "ConsentLog",
			AggregateID:   consent.ID,
			EventType:     event.TopicConsentGranted,
			Topic:         event.TopicConsentGranted,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		}); err != nil {
			return fmt.Errorf("save outbox: %w", err)
		}
		h.logger.Infow("consent granted",
			"user_id", consent.UserID,
			"purpose", consent.Purpose.String(),
			"consent_id", consent.ID,
		)
		return nil
	})
}
```

### 6.2 `application/command/revoke_consent.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/logger"
)

type RevokeConsentHandler struct {
	consentRepo repository.ConsentRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewRevokeConsentHandler(
	consentRepo repository.ConsentRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *RevokeConsentHandler {
	if clock == nil {
		clock = time.Now
	}
	return &RevokeConsentHandler{
		consentRepo: consentRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		clock:       clock,
		logger:      logger,
	}
}

func (h *RevokeConsentHandler) Handle(ctx context.Context, req dto.RevokeConsentRequest) error {
	purpose, err := valueobject.ParseConsentPurpose(req.Purpose)
	if err != nil {
		return err
	}
	if purpose.IsRequired() {
		return domainerrors.ErrRevokeNotAllowed
	}

	consent, err := h.consentRepo.FindActiveByUserAndPurpose(ctx, req.UserID, purpose)
	if err != nil {
		return err
	}
	if consent == nil {
		return domainerrors.ErrConsentNotFound
	}

	now := h.clock()
	if err := consent.Revoke(now); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicConsentRevoked, meta, event.ConsentRevokedPayload{
		ConsentID: consent.ID,
		UserID:    consent.UserID,
		Purpose:   consent.Purpose.String(),
		RevokedAt: *consent.RevokedAt,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.consentRepo.Update(txCtx, consent); err != nil {
			return err
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "ConsentLog",
			AggregateID:   consent.ID,
			EventType:     event.TopicConsentRevoked,
			Topic:         event.TopicConsentRevoked,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}

var _ = fmt.Sprintf
```

### 6.3 `application/command/submit_dsar.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/logger"
)

const maxDSARPerDay = 3

type SubmitDSARHandler struct {
	dsarRepo    repository.DSARRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewSubmitDSARHandler(
	dsarRepo repository.DSARRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *SubmitDSARHandler {
	if clock == nil {
		clock = time.Now
	}
	return &SubmitDSARHandler{
		dsarRepo:   dsarRepo,
		outboxRepo: outboxRepo,
		uow:        uow,
		clock:      clock,
		logger:     logger,
	}
}

func (h *SubmitDSARHandler) Handle(ctx context.Context, req dto.SubmitDSARRequest) (*entity.DSARRequest, error) {
	reqType, err := valueobject.ParseDSARType(req.RequestType)
	if err != nil {
		return nil, err
	}
	now := h.clock()

	// Rate limit
	count, err := h.dsarRepo.CountByUserAndTypeInPeriod(ctx, req.UserID, reqType, now.Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}
	if count >= maxDSARPerDay {
		return nil, domainerrors.ErrRateLimitExceeded
	}

	dsar, otp, err := entity.NewDSARRequest(req.UserID, reqType, now, req.IPAddress, req.UserAgent)
	if err != nil {
		return nil, err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicDSARSubmitted, meta, event.DSARSubmittedPayload{
		DSARID:      dsar.ID,
		UserID:      dsar.UserID,
		RequestType: dsar.RequestType.String(),
		RequestedAt: dsar.RequestedAt,
	})
	body, _ := json.Marshal(envelope)

	// Email OTP event (แยก topic)
	emailEnv := event.NewEnvelope(event.TopicEmailNotification, meta, event.EmailNotificationPayload{
		Template:  "dsar_otp",
		Locale:    "th",
		UserID:    req.UserID.String(),
		Variables: map[string]interface{}{"otp": otp, "dsar_id": dsar.ID.String()},
	})
	emailBody, _ := json.Marshal(emailEnv)

	if err := h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.dsarRepo.Save(txCtx, dsar); err != nil {
			return fmt.Errorf("save dsar: %w", err)
		}
		if err := h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "DSARRequest",
			AggregateID:   dsar.ID,
			EventType:     event.TopicDSARSubmitted,
			Topic:         event.TopicDSARSubmitted,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		}); err != nil {
			return err
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            uuid.New(),
			AggregateType: "DSARRequest",
			AggregateID:   dsar.ID,
			EventType:     event.TopicEmailNotification,
			Topic:         event.TopicEmailNotification,
			Payload:       emailBody,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	}); err != nil {
		return nil, err
	}

	h.logger.Infow("dsar submitted", "user_id", req.UserID, "dsar_id", dsar.ID)
	return dsar, nil
}

var _ = uuid.Nil
```

> **หมายเหตุ:** ต้อง `import "github.com/google/uuid"` เพิ่มในไฟล์นี้

### 6.4 `application/command/process_dsar.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type ProcessDSARHandler struct {
	dsarRepo   repository.DSARRepository
	outboxRepo repository.OutboxRepository
	uow        repository.UnitOfWork
	logger     logger.Logger
}

func NewProcessDSARHandler(
	dsarRepo repository.DSARRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	logger logger.Logger,
) *ProcessDSARHandler {
	return &ProcessDSARHandler{dsarRepo: dsarRepo, outboxRepo: outboxRepo, uow: uow, logger: logger}
}

func (h *ProcessDSARHandler) Handle(ctx context.Context, req dto.ProcessDSARRequest) error {
	dsar, err := h.dsarRepo.FindByID(ctx, req.DSARRequestID)
	if err != nil {
		return err
	}
	if dsar == nil {
		return domainerrors.ErrDSARNotFound
	}
	if err := dsar.MarkProcessing(); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicDSARProcessing, meta, event.DSARProcessingPayload{
		DSARID:      dsar.ID,
		UserID:      dsar.UserID,
		ProcessedAt: meta.OccurredAt,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.dsarRepo.Update(txCtx, dsar); err != nil {
			return fmt.Errorf("update dsar: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "DSARRequest",
			AggregateID:   dsar.ID,
			EventType:     event.TopicDSARProcessing,
			Topic:         event.TopicDSARProcessing,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     meta.OccurredAt,
		})
	})
}
```

### 6.5 `application/command/confirm_deletion.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

type ConfirmDeletionHandler struct {
	accountRepo repository.UserAccountStatusRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	policy      *service.DeletionPolicyService
	clock       func() time.Time
	logger      logger.Logger
}

func NewConfirmDeletionHandler(
	accountRepo repository.UserAccountStatusRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	policy *service.DeletionPolicyService,
	clock func() time.Time,
	logger logger.Logger,
) *ConfirmDeletionHandler {
	if clock == nil {
		clock = time.Now
	}
	return &ConfirmDeletionHandler{
		accountRepo: accountRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		policy:      policy,
		clock:       clock,
		logger:      logger,
	}
}

func (h *ConfirmDeletionHandler) Handle(ctx context.Context, req dto.ConfirmDeletionRequest) error {
	acct, err := h.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if acct == nil {
		return nil
	}

	now := h.clock()
	if err := acct.ConfirmDeletion(now); err != nil {
		return err
	}

	immediate := h.policy.CanImmediateDeletion(acct)

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicDeletionRequested, meta, event.DataDeletionRequestedPayload{
		UserID:              acct.UserID,
		DeletionConfirmedAt: *acct.DeletionConfirmedAt,
		ImmediateDeletion:   immediate,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.accountRepo.Save(txCtx, acct); err != nil {
			return fmt.Errorf("save account: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "UserAccountStatus",
			AggregateID:   acct.UserID,
			EventType:     event.TopicDeletionRequested,
			Topic:         event.TopicDeletionRequested,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}
```

### 6.6 `application/command/immediate_deletion.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type ImmediateDeletionHandler struct {
	consentRepo repository.ConsentRepository
	dsarRepo    repository.DSARRepository
	accountRepo repository.UserAccountStatusRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewImmediateDeletionHandler(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountRepo repository.UserAccountStatusRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *ImmediateDeletionHandler {
	if clock == nil {
		clock = time.Now
	}
	return &ImmediateDeletionHandler{
		consentRepo: consentRepo,
		dsarRepo:    dsarRepo,
		accountRepo: accountRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		clock:       clock,
		logger:      logger,
	}
}

// Handle ลบข้อมูลทั้งหมดของผู้ใช้ — จะถูกเรียกโดย DataDeletionRequested handler หรือ admin
func (h *ImmediateDeletionHandler) Handle(ctx context.Context, userID uuid.UUID, reason string) error {
	acct, err := h.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if acct == nil {
		return nil
	}

	now := h.clock()
	acct.MarkDeleted(now)

	meta := event.NewMetadata("", "", "")
	envelope := event.NewEnvelope(event.TopicDataDeleted, meta, event.DataDeletedPayload{
		UserID:     userID,
		DeleteType: "IMMEDIATE",
		Reason:     reason,
		DeletedAt:  now,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.consentRepo.DeleteByUserID(txCtx, userID); err != nil {
			return fmt.Errorf("delete consents: %w", err)
		}
		if err := h.dsarRepo.DeleteByUserID(txCtx, userID); err != nil {
			return fmt.Errorf("delete dsars: %w", err)
		}
		if err := h.accountRepo.Save(txCtx, acct); err != nil {
			return fmt.Errorf("save account: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "UserAccountStatus",
			AggregateID:   userID,
			EventType:     event.TopicDataDeleted,
			Topic:         event.TopicDataDeleted,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}

var _ = uuid.Nil
```

### 6.7 `application/command/auto_delete_expired.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

type AutoDeleteExpiredResult struct {
	AccountsDeleted int
	ConsentsDeleted int
	Errors          []string
}

type AutoDeleteExpiredHandler struct {
	consentRepo repository.ConsentRepository
	dsarRepo    repository.DSARRepository
	accountRepo repository.UserAccountStatusRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	policy      *service.DeletionPolicyService
	clock       func() time.Time
	logger      logger.Logger
}

func NewAutoDeleteExpiredHandler(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountRepo repository.UserAccountStatusRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	policy *service.DeletionPolicyService,
	clock func() time.Time,
	logger logger.Logger,
) *AutoDeleteExpiredHandler {
	if clock == nil {
		clock = time.Now
	}
	return &AutoDeleteExpiredHandler{
		consentRepo: consentRepo,
		dsarRepo:    dsarRepo,
		accountRepo: accountRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		policy:      policy,
		clock:       clock,
		logger:      logger,
	}
}

func (h *AutoDeleteExpiredHandler) Handle(ctx context.Context, batchSize int) (*AutoDeleteExpiredResult, error) {
	if batchSize <= 0 {
		batchSize = 100
	}
	now := h.clock()
	result := &AutoDeleteExpiredResult{}

	// Phase 1: consent logs เก่า
	cutoff := h.policy.ConsentAutoDeleteCutoff(now)
	expired, err := h.consentRepo.FindReadyForAutoDeletion(ctx, cutoff, batchSize)
	if err != nil {
		return nil, err
	}
	for i := range expired {
		c := &expired[i]
		c.MarkDeleted(now)
		if err := h.consentRepo.Update(ctx, c); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.ConsentsDeleted++
	}

	// Phase 2: accounts ที่เลย retention
	accounts, err := h.accountRepo.FindReadyForAutoDeletion(ctx, now, batchSize)
	if err != nil {
		return nil, err
	}
	for i := range accounts {
		acct := &accounts[i]
		if !h.policy.IsReadyForAutoDeletion(acct, now) {
			continue
		}
		if err := h.deleteUserData(ctx, acct.UserID, now); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.AccountsDeleted++
	}

	return result, nil
}

func (h *AutoDeleteExpiredHandler) deleteUserData(ctx context.Context, userID uuid.UUID, now time.Time) error {
	acct, err := h.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if acct == nil {
		return nil
	}
	acct.MarkDeleted(now)

	meta := event.NewMetadata("", "", "")
	envelope := event.NewEnvelope(event.TopicDataDeleted, meta, event.DataDeletedPayload{
		UserID:     userID,
		DeleteType: "AUTO",
		Reason:     "retention_policy_expired",
		DeletedAt:  now,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.consentRepo.DeleteByUserID(txCtx, userID); err != nil {
			return fmt.Errorf("delete consents: %w", err)
		}
		if err := h.dsarRepo.DeleteByUserID(txCtx, userID); err != nil {
			return fmt.Errorf("delete dsars: %w", err)
		}
		if err := h.accountRepo.Save(txCtx, acct); err != nil {
			return fmt.Errorf("save account: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "UserAccountStatus",
			AggregateID:   userID,
			EventType:     event.TopicDataDeleted,
			Topic:         event.TopicDataDeleted,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}
```

### 6.8 `application/command/handle_account_suspended.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type HandleAccountSuspendedHandler struct {
	accountRepo repository.UserAccountStatusRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewHandleAccountSuspendedHandler(
	accountRepo repository.UserAccountStatusRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *HandleAccountSuspendedHandler {
	if clock == nil {
		clock = time.Now
	}
	return &HandleAccountSuspendedHandler{
		accountRepo: accountRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		clock:       clock,
		logger:      logger,
	}
}

func (h *HandleAccountSuspendedHandler) Handle(ctx context.Context, req dto.SuspendAccountRequest) error {
	acct, err := h.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil || acct == nil {
		acct, err = entity.NewUserAccountStatus(req.UserID, h.clock())
		if err != nil {
			return err
		}
	}

	now := h.clock()
	if err := acct.Suspend(now, req.RetentionYears); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicAccountSuspended, meta, event.AccountSuspendedPayload{
		UserID:            acct.UserID,
		Status:            acct.Status.String(),
		SuspendedAt:       *acct.SuspendedAt,
		RetentionDeadline: *acct.RetentionDeadline,
		RetentionYears:    req.RetentionYears,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.accountRepo.Save(txCtx, acct); err != nil {
			return fmt.Errorf("save account: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "UserAccountStatus",
			AggregateID:   acct.UserID,
			EventType:     event.TopicAccountSuspended,
			Topic:         event.TopicAccountSuspended,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}
```

### 6.9 `application/command/handle_account_terminated.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type HandleAccountTerminatedHandler struct {
	accountRepo repository.UserAccountStatusRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewHandleAccountTerminatedHandler(
	accountRepo repository.UserAccountStatusRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *HandleAccountTerminatedHandler {
	if clock == nil {
		clock = time.Now
	}
	return &HandleAccountTerminatedHandler{
		accountRepo: accountRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		clock:       clock,
		logger:      logger,
	}
}

func (h *HandleAccountTerminatedHandler) Handle(ctx context.Context, req dto.TerminateAccountRequest) error {
	acct, err := h.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil || acct == nil {
		acct, err = entity.NewUserAccountStatus(req.UserID, h.clock())
		if err != nil {
			return err
		}
	}

	now := h.clock()
	if err := acct.Terminate(now); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicAccountTerminated, meta, event.AccountTerminatedPayload{
		UserID:       acct.UserID,
		Status:       acct.Status.String(),
		TerminatedAt: *acct.TerminatedAt,
	})
	body, _ := json.Marshal(envelope)

	return h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.accountRepo.Save(txCtx, acct); err != nil {
			return fmt.Errorf("save account: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "UserAccountStatus",
			AggregateID:   acct.UserID,
			EventType:     event.TopicAccountTerminated,
			Topic:         event.TopicAccountTerminated,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	})
}
```

### 6.10 `application/command/publish_privacy_policy.go`

```go
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type PublishPrivacyPolicyHandler struct {
	policyRepo repository.PolicyRepository
	outboxRepo repository.OutboxRepository
	uow        repository.UnitOfWork
	clock      func() time.Time
	logger     logger.Logger
}

func NewPublishPrivacyPolicyHandler(
	policyRepo repository.PolicyRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *PublishPrivacyPolicyHandler {
	if clock == nil {
		clock = time.Now
	}
	return &PublishPrivacyPolicyHandler{
		policyRepo: policyRepo,
		outboxRepo: outboxRepo,
		uow:        uow,
		clock:      clock,
		logger:     logger,
	}
}

func (h *PublishPrivacyPolicyHandler) Handle(ctx context.Context, req dto.PublishPolicyRequest) (*entity.PrivacyPolicy, error) {
	if existing, _ := h.policyRepo.FindByVersion(ctx, req.Version); existing != nil {
		return nil, domainerrors.ErrPolicyVersionExists
	}

	policy, err := entity.NewPrivacyPolicy(req.Version, req.Title, req.Content, req.EffectiveDate, h.clock())
	if err != nil {
		return nil, err
	}
	policy.Activate(h.clock())

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicPolicyPublished, meta, event.PolicyPublishedPayload{
		PolicyID:      policy.ID,
		Version:       policy.Version,
		Title:         policy.Title,
		EffectiveDate: policy.EffectiveDate,
		PublishedAt:   h.clock(),
	})
	body, _ := json.Marshal(envelope)

	if err := h.uow.Do(ctx, func(txCtx context.Context) error {
		if err := h.policyRepo.DeactivateAll(txCtx); err != nil {
			return fmt.Errorf("deactivate old: %w", err)
		}
		if err := h.policyRepo.Save(txCtx, policy); err != nil {
			return fmt.Errorf("save policy: %w", err)
		}
		return h.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "PrivacyPolicy",
			AggregateID:   policy.ID,
			EventType:     event.TopicPolicyPublished,
			Topic:         event.TopicPolicyPublished,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     h.clock(),
		})
	}); err != nil {
		return nil, err
	}

	return policy, nil
}
```

---

## 📁 7. Application Layer — Queries

### 7.1 `application/query/get_consent_history.go`

```go
package query

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type GetConsentHistoryHandler struct {
	consentRepo repository.ConsentRepository
}

func NewGetConsentHistoryHandler(consentRepo repository.ConsentRepository) *GetConsentHistoryHandler {
	return &GetConsentHistoryHandler{consentRepo: consentRepo}
}

func (h *GetConsentHistoryHandler) Handle(ctx context.Context, userID uuid.UUID) (*dto.ConsentStatusListResponse, error) {
	items, err := h.consentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ConsentStatusResponse, 0, len(items))
	for _, c := range items {
		out = append(out, dto.ConsentStatusResponse{
			UserID:    c.UserID,
			Purpose:   c.Purpose.String(),
			Status:    c.Status.String(),
			IsActive:  c.IsActive(nowUTC()),
			GrantedAt: c.ConsentedAt,
			ExpiresAt: c.ExpiresAt,
			RevokedAt: c.RevokedAt,
		})
	}
	return &dto.ConsentStatusListResponse{Consents: out, Total: int64(len(out))}, nil
}
```

### 7.2 `application/query/get_dsar_status.go`

```go
package query

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type GetDSARStatusHandler struct {
	dsarRepo repository.DSARRepository
}

func NewGetDSARStatusHandler(dsarRepo repository.DSARRepository) *GetDSARStatusHandler {
	return &GetDSARStatusHandler{dsarRepo: dsarRepo}
}

func (h *GetDSARStatusHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.DSARStatusResponse, error) {
	d, err := h.dsarRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, domainerrors.ErrDSARNotFound
	}
	return &dto.DSARStatusResponse{
		ID:              d.ID,
		UserID:          d.UserID,
		RequestType:     d.RequestType.String(),
		Status:          d.Status.String(),
		RequestedAt:     d.RequestedAt,
		CompletedAt:     d.CompletedAt,
		RejectionReason: d.RejectionReason,
	}, nil
}

type ListDSARHandler struct {
	dsarRepo repository.DSARRepository
}

func NewListDSARHandler(dsarRepo repository.DSARRepository) *ListDSARHandler {
	return &ListDSARHandler{dsarRepo: dsarRepo}
}

func (h *ListDSARHandler) Handle(ctx context.Context, userID uuid.UUID) (*dto.DSARStatusListResponse, error) {
	items, err := h.dsarRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DSARStatusResponse, 0, len(items))
	for _, d := range items {
		out = append(out, dto.DSARStatusResponse{
			ID:              d.ID,
			UserID:          d.UserID,
			RequestType:     d.RequestType.String(),
			Status:          d.Status.String(),
			RequestedAt:     d.RequestedAt,
			CompletedAt:     d.CompletedAt,
			RejectionReason: d.RejectionReason,
		})
	}
	return &dto.DSARStatusListResponse{Requests: out, Total: int64(len(out))}, nil
}

// helper
func nowUTC() time.Time { return time.Now().UTC() }
```

### 7.3 `application/query/get_admin_report.go`

```go
package query

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type GetAdminReportHandler struct {
	auditRepo repository.AuditRepository
}

func NewGetAdminReportHandler(auditRepo repository.AuditRepository) *GetAdminReportHandler {
	return &GetAdminReportHandler{auditRepo: auditRepo}
}

func (h *GetAdminReportHandler) Handle(ctx context.Context, from, to time.Time) (*dto.AdminReportResponse, error) {
	actions, err := h.auditRepo.CountByAction(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return &dto.AdminReportResponse{
		From:          from,
		To:            to,
		ActionsByType: actions,
		GeneratedAt:   time.Now().UTC(),
	}, nil
}

type ListAuditTrailsHandler struct {
	auditRepo repository.AuditRepository
}

func NewListAuditTrailsHandler(auditRepo repository.AuditRepository) *ListAuditTrailsHandler {
	return &ListAuditTrailsHandler{auditRepo: auditRepo}
}

func (h *ListAuditTrailsHandler) Handle(ctx context.Context, filter repository.AuditFilter) ([]dto.AuditTrailItem, error) {
	items, err := h.auditRepo.FindByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AuditTrailItem, 0, len(items))
	for _, a := range items {
		out = append(out, dto.AuditTrailItem{
			ID:        a.ID,
			UserID:    a.UserID,
			Action:    a.Action,
			Details:   a.Details,
			IPAddress: a.IPAddress,
			UserAgent: a.UserAgent,
			CreatedAt: a.CreatedAt,
		})
	}
	return out, nil
}
```

### 7.4 `application/query/get_active_privacy_policy.go`

```go
package query

import (
	"context"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type GetActivePrivacyPolicyHandler struct {
	policyRepo repository.PolicyRepository
}

func NewGetActivePrivacyPolicyHandler(policyRepo repository.PolicyRepository) *GetActivePrivacyPolicyHandler {
	return &GetActivePrivacyPolicyHandler{policyRepo: policyRepo}
}

func (h *GetActivePrivacyPolicyHandler) Handle(ctx context.Context) (*dto.PrivacyPolicyResponse, error) {
	p, err := h.policyRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, domainerrors.ErrPolicyNotFound
	}
	return toPolicyResponse(p), nil
}

type ListPoliciesHandler struct {
	policyRepo repository.PolicyRepository
}

func NewListPoliciesHandler(policyRepo repository.PolicyRepository) *ListPoliciesHandler {
	return &ListPoliciesHandler{policyRepo: policyRepo}
}

func (h *ListPoliciesHandler) Handle(ctx context.Context) ([]dto.PrivacyPolicyResponse, error) {
	items, err := h.policyRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PrivacyPolicyResponse, 0, len(items))
	for _, p := range items {
		out = append(out, *toPolicyResponse(&p))
	}
	return out, nil
}

func toPolicyResponse(p *entity.PrivacyPolicy) *dto.PrivacyPolicyResponse {
	return &dto.PrivacyPolicyResponse{
		ID:            p.ID,
		Version:       p.Version,
		Title:         p.Title,
		Content:       p.Content,
		EffectiveDate: p.EffectiveDate,
		IsActive:      p.IsActive,
		UpdatedAt:     p.UpdatedAt,
	}
}
```

---

## 📁 8. Application Layer — Event Handlers (Idempotent)

### 8.1 `application/event_handler/base_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// Envelope โครงสร้าง event ที่ consume จาก Kafka (mirror domain/event.Envelope)
type Envelope struct {
	EventType string          `json:"event_type"`
	Metadata  EventMetadata   `json:"metadata"`
	Payload   json.RawMessage `json:"payload"`
}

type EventMetadata struct {
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id,omitempty"`
	TraceID       string    `json:"trace_id,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
	Version       int       `json:"version"`
	Source        string    `json:"source"`
}

func DecodeEnvelope(data []byte) (*Envelope, error) {
	var e Envelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("decode envelope: %w", err)
	}
	if e.EventType == "" || e.Metadata.EventID == uuid.Nil {
		return nil, fmt.Errorf("invalid envelope")
	}
	return &e, nil
}

func DecodePayload[T any](e *Envelope) (*T, error) {
	var out T
	if err := json.Unmarshal(e.Payload, &out); err != nil {
		return nil, fmt.Errorf("decode payload for %s: %w", e.EventType, err)
	}
	return &out, nil
}

type IdempotencyStore interface {
	IsProcessed(ctx context.Context, eventID uuid.UUID, handler string) (bool, error)
	MarkProcessed(ctx context.Context, eventID uuid.UUID, handler string, ttl time.Duration) error
}

type Metrics interface {
	IncHandlerProcessed(handler, result string)
	ObserveHandlerDuration(handler string, seconds float64)
}

type nopMetrics struct{}

func (nopMetrics) IncHandlerProcessed(string, string)     {}
func (nopMetrics) ObserveHandlerDuration(string, float64) {}
func NopMetrics() Metrics                                  { return nopMetrics{} }

// Handler interface
type Handler interface {
	Handle(ctx context.Context, env *Envelope) error
	Name() string
}

// BaseHandler ที่ทุก handler embed
type BaseHandler struct {
	auditRepo   repository.AuditRepository
	idempotency IdempotencyStore
	metrics     Metrics
	redactor    *PIIRedactor
	logger      logger.Logger
	name        string
}

func NewBaseHandler(
	name string,
	auditRepo repository.AuditRepository,
	idem IdempotencyStore,
	metrics Metrics,
	redactor *PIIRedactor,
	logger logger.Logger,
) BaseHandler {
	if metrics == nil {
		metrics = NopMetrics()
	}
	if redactor == nil {
		redactor = NewPIIRedactor()
	}
	return BaseHandler{
		name: name, auditRepo: auditRepo, idempotency: idem,
		metrics: metrics, redactor: redactor, logger: logger,
	}
}

func (h *BaseHandler) Name() string { return h.name }

func (h *BaseHandler) IsProcessed(ctx context.Context, id uuid.UUID) (bool, error) {
	if h.idempotency == nil {
		return false, nil
	}
	return h.idempotency.IsProcessed(ctx, id, h.name)
}

func (h *BaseHandler) MarkProcessed(ctx context.Context, id uuid.UUID, ttl time.Duration) error {
	if h.idempotency == nil {
		return nil
	}
	return h.idempotency.MarkProcessed(ctx, id, h.name, ttl)
}

func (h *BaseHandler) SaveAudit(ctx context.Context, userID *uuid.UUID, action string, payload map[string]interface{}) error {
	redacted := h.redactor.RedactPayload(payload)
	audit := entity.NewAuditTrail(userID, action, redacted)
	if err := h.auditRepo.Save(ctx, audit); err != nil {
		return fmt.Errorf("save audit for %s: %w", action, err)
	}
	return nil
}

func (h *BaseHandler) Observe(ctx context.Context, meta EventMetadata, fn func(ctx context.Context) error) error {
	start := time.Now()
	err := fn(ctx)
	dur := time.Since(start).Seconds()

	h.metrics.ObserveHandlerDuration(h.name, dur)
	result := "success"
	if err != nil {
		result = "error"
	}
	h.metrics.IncHandlerProcessed(h.name, result)

	logFn := h.logger.Infow
	if err != nil {
		logFn = h.logger.Errorw
	}
	logFn("event handled",
		"handler", h.name, "event_id", meta.EventID,
		"correlation_id", meta.CorrelationID,
		"duration_ms", dur*1000, "result", result, "error", err)
	return err
}
```

### 8.2 `application/event_handler/pii_redactor.go`

```go
package eventhandler

import (
	"encoding/json"
	"regexp"
	"strings"
)

type PIIRedactor struct {
	maskKeys []string
	patterns []redactPattern
}

type redactPattern struct {
	re      *regexp.Regexp
	replace string
}

func NewPIIRedactor() *PIIRedactor {
	return &PIIRedactor{
		maskKeys: []string{
			"email", "phone", "mobile", "otp", "otp_code",
			"password", "token", "access_token", "refresh_token",
			"national_id", "id_card", "passport", "ssn",
			"credit_card", "card_number", "cvv", "iban",
			"address", "date_of_birth", "dob",
		},
		patterns: []redactPattern{
			{regexp.MustCompile(`[\w._%+\-]+@[\w.\-]+\.[A-Za-z]{2,}`), "***@***.***"},
			{regexp.MustCompile(`\b0\d{8,9}\b`), "0XXXXXXXXX"},
			{regexp.MustCompile(`\b\d{13}\b`), "XXXXXXXXXXXXX"},
			{regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`), "XXXX-XXXX-XXXX-XXXX"},
		},
	}
}

func (r *PIIRedactor) RedactPayload(payload map[string]interface{}) map[string]interface{} {
	if payload == nil {
		return nil
	}
	out := make(map[string]interface{}, len(payload))
	for k, v := range payload {
		if r.isSensitiveKey(k) {
			out[k] = "***"
			continue
		}
		out[k] = r.redactValue(v)
	}
	return out
}

func (r *PIIRedactor) isSensitiveKey(k string) bool {
	lower := strings.ToLower(k)
	for _, s := range r.maskKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

func (r *PIIRedactor) redactValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return r.redactString(val)
	case map[string]interface{}:
		return r.RedactPayload(val)
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, item := range val {
			out[i] = r.redactValue(item)
		}
		return out
	default:
		return v
	}
}

func (r *PIIRedactor) redactString(s string) string {
	for _, p := range r.patterns {
		s = p.re.ReplaceAllString(s, p.replace)
	}
	return s
}

func (r *PIIRedactor) RedactJSONString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "<unmarshalable>"
	}
	return r.redactString(string(b))
}
```

### 8.3 `application/event_handler/consent_granted_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type ConsentGrantedPayload struct {
	ConsentID uuid.UUID `json:"consent_id"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	GrantedAt time.Time `json:"granted_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address,omitempty"`
}

type ConsentGrantedHandler struct {
	BaseHandler
	cache    repository.Cache
	producer kafka.Producer
}

func NewConsentGrantedHandler(
	auditRepo repository.AuditRepository,
	cache repository.Cache,
	producer kafka.Producer,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *ConsentGrantedHandler {
	return &ConsentGrantedHandler{
		BaseHandler: NewBaseHandler("ConsentGrantedHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		cache:       cache,
		producer:    producer,
	}
}

func (h *ConsentGrantedHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[ConsentGrantedPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}
		purpose := valueobject.ConsentPurpose(payload.Purpose)
		_ = h.cache.Set(ctx, payload.UserID, purpose, valueobject.ConsentGranted)

		emailBody, _ := json.Marshal(map[string]interface{}{
			"template": "consent_granted", "locale": "th",
			"user_id": payload.UserID.String(), "purpose": payload.Purpose,
		})
		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(), emailBody)

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"consent_id": payload.ConsentID.String(),
			"user_id":    payload.UserID.String(),
			"purpose":    payload.Purpose,
			"granted_at": payload.GrantedAt,
			"expires_at": payload.ExpiresAt,
			"ip_address": payload.IPAddress,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}

var _ = uuid.Nil
```

> **เพิ่ม import `github.com/google/uuid`**

### 8.4 `application/event_handler/consent_revoked_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type ConsentRevokedPayload struct {
	ConsentID uuid.UUID `json:"consent_id"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	RevokedAt time.Time `json:"revoked_at"`
}

type ConsentRevokedHandler struct {
	BaseHandler
	cache    repository.Cache
	producer kafka.Producer
}

func NewConsentRevokedHandler(
	auditRepo repository.AuditRepository,
	cache repository.Cache,
	producer kafka.Producer,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *ConsentRevokedHandler {
	return &ConsentRevokedHandler{
		BaseHandler: NewBaseHandler("ConsentRevokedHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		cache:       cache,
		producer:    producer,
	}
}

func (h *ConsentRevokedHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[ConsentRevokedPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}
		purpose := valueobject.ConsentPurpose(payload.Purpose)
		_ = h.cache.Set(ctx, payload.UserID, purpose, valueobject.ConsentRevoked)

		emailBody, _ := json.Marshal(map[string]interface{}{
			"template": "consent_revoked", "locale": "th",
			"user_id": payload.UserID.String(), "purpose": payload.Purpose,
		})
		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(), emailBody)

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"consent_id": payload.ConsentID.String(),
			"user_id":    payload.UserID.String(),
			"purpose":    payload.Purpose,
			"revoked_at": payload.RevokedAt,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 8.5 `application/event_handler/dsar_submitted_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type DSARSubmittedPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	RequestType string    `json:"request_type"`
	RequestedAt time.Time `json:"requested_at"`
}

type DSARSubmittedHandler struct {
	BaseHandler
	producer kafka.Producer
}

func NewDSARSubmittedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *DSARSubmittedHandler {
	return &DSARSubmittedHandler{
		BaseHandler: NewBaseHandler("DSARSubmittedHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		producer:    producer,
	}
}

func (h *DSARSubmittedHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[DSARSubmittedPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}

		if payload.RequestType == "ACCESS" || payload.RequestType == "ERASURE" {
			body, _ := json.Marshal(map[string]interface{}{
				"dsar_id":       payload.DSARID.String(),
				"user_id":       payload.UserID.String(),
				"analysis_type": "DATA_MAPPING",
			})
			_ = h.producer.Publish(ctx, event.TopicLLMAnalysis, payload.DSARID.String(), body)
		}

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"dsar_id":      payload.DSARID.String(),
			"user_id":      payload.UserID.String(),
			"request_type": payload.RequestType,
			"requested_at": payload.RequestedAt,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 8.6 `application/event_handler/data_deletion_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type DataDeletionRequestedPayload struct {
	UserID              uuid.UUID `json:"user_id"`
	DeletionConfirmedAt time.Time `json:"deletion_confirmed_at"`
	ImmediateDeletion   bool      `json:"immediate_deletion"`
}

type DataDeletionHandler struct {
	BaseHandler
	producer kafka.Producer
}

func NewDataDeletionHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *DataDeletionHandler {
	return &DataDeletionHandler{
		BaseHandler: NewBaseHandler("DataDeletionHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		producer:    producer,
	}
}

func (h *DataDeletionHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[DataDeletionRequestedPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}

		// Trigger deletion workflow
		if payload.ImmediateDeletion {
			body, _ := json.Marshal(map[string]interface{}{
				"user_id":     payload.UserID.String(),
				"delete_type": "IMMEDIATE",
			})
			_ = h.producer.Publish(ctx, event.TopicDataDeleted, payload.UserID.String(), body)
		}

		// Blockchain audit
		bcBody, _ := json.Marshal(map[string]interface{}{
			"user_id": payload.UserID.String(),
			"action":  "DELETION_REQUESTED",
		})
		_ = h.producer.Publish(ctx, event.TopicBlockchainRecord, payload.UserID.String(), bcBody)

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"user_id":               payload.UserID.String(),
			"deletion_confirmed_at": payload.DeletionConfirmedAt,
			"immediate_deletion":    payload.ImmediateDeletion,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 30*24*time.Hour)
	})
}
```

### 8.7 `application/event_handler/account_suspended_handler.go`

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type AccountSuspendedPayload struct {
	UserID            uuid.UUID `json:"user_id"`
	Status            string    `json:"status"`
	SuspendedAt       time.Time `json:"suspended_at"`
	RetentionDeadline time.Time `json:"retention_deadline"`
	RetentionYears    int       `json:"retention_years"`
}

type AccountSuspendedHandler struct {
	BaseHandler
	producer kafka.Producer
}

func NewAccountSuspendedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *AccountSuspendedHandler {
	return &AccountSuspendedHandler{
		BaseHandler: NewBaseHandler("AccountSuspendedHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		producer:    producer,
	}
}

func (h *AccountSuspendedHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[AccountSuspendedPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}
		body, _ := json.Marshal(map[string]interface{}{
			"template": "account_suspended", "locale": "th",
			"user_id": payload.UserID.String(), "retention_years": payload.RetentionYears,
		})
		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(), body)

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"user_id":            payload.UserID.String(),
			"status":             payload.Status,
			"suspended_at":       payload.SuspendedAt,
			"retention_deadline": payload.RetentionDeadline,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 8.8 `application/event_handler/llm_analysis_handler.go`

```go
package eventhandler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type LLMClient interface {
	AnalyzeDataWithContext(ctx context.Context, analysisType string, payload map[string]interface{}) (string, error)
}

type LLMAnalysisPayload struct {
	DSARID       uuid.UUID `json:"dsar_id"`
	UserID       uuid.UUID `json:"user_id"`
	AnalysisType string    `json:"analysis_type"`
}

type LLMAnalysisHandler struct {
	BaseHandler
	llm LLMClient
}

func NewLLMAnalysisHandler(
	auditRepo repository.AuditRepository,
	llm LLMClient,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *LLMAnalysisHandler {
	return &LLMAnalysisHandler{
		BaseHandler: NewBaseHandler("LLMAnalysisHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		llm:         llm,
	}
}

func (h *LLMAnalysisHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[LLMAnalysisPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		result, err := h.llm.AnalyzeDataWithContext(ctx, payload.AnalysisType, map[string]interface{}{
			"dsar_id": payload.DSARID.String(),
			"user_id": payload.UserID.String(),
		})
		if err != nil {
			return err // retry
		}

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"dsar_id":       payload.DSARID.String(),
			"analysis_type": payload.AnalysisType,
			"summary_len":   len(result),
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 8.9 `application/event_handler/email_notification_handler.go`

```go
package eventhandler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type EmailClient interface {
	Send(ctx context.Context, to, subject, template, locale string, vars map[string]interface{}) error
}

type EmailNotificationPayload struct {
	Template  string                 `json:"template"`
	Locale    string                 `json:"locale"`
	To        string                 `json:"to,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type EmailNotificationHandler struct {
	BaseHandler
	email EmailClient
}

func NewEmailNotificationHandler(
	auditRepo repository.AuditRepository,
	email EmailClient,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *EmailNotificationHandler {
	return &EmailNotificationHandler{
		BaseHandler: NewBaseHandler("EmailNotificationHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		email:       email,
	}
}

func (h *EmailNotificationHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[EmailNotificationPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := h.email.Send(ctx, payload.To, payload.Template, payload.Template, payload.Locale, payload.Variables); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 8.10 `application/event_handler/blockchain_record_handler.go`

```go
package eventhandler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type BlockchainClient interface {
	Record(ctx context.Context, userID, action, dataHash string) (string, error)
}

type BlockchainRecordPayload struct {
	UserID   string `json:"user_id"`
	Action   string `json:"action"`
	DataHash string `json:"data_hash"`
}

type BlockchainRecordHandler struct {
	BaseHandler
	client BlockchainClient
}

func NewBlockchainRecordHandler(
	auditRepo repository.AuditRepository,
	client BlockchainClient,
	idem IdempotencyStore,
	metrics Metrics,
	logger logger.Logger,
) *BlockchainRecordHandler {
	return &BlockchainRecordHandler{
		BaseHandler: NewBaseHandler("BlockchainRecordHandler", auditRepo, idem, metrics, NewPIIRedactor(), logger),
		client:      client,
	}
}

func (h *BlockchainRecordHandler) Handle(ctx context.Context, env *Envelope) error {
	payload, err := DecodePayload[BlockchainRecordPayload](env)
	if err != nil {
		return err
	}
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil || done {
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
		defer cancel()

		txHash, err := h.client.Record(ctx, payload.UserID, payload.Action, payload.DataHash)
		if err != nil {
			return err
		}

		if err := h.SaveAudit(ctx, nil, env.EventType, map[string]interface{}{
			"user_id": payload.UserID,
			"action":  payload.Action,
			"tx_hash": txHash,
		}); err != nil {
			return err
		}
		return h.MarkProcessed(ctx, env.Metadata.EventID, 30*24*time.Hour)
	})
}
```

---

## 📁 9. Infrastructure — Postgres

### 9.1 `infrastructure/persistence/postgres/models.go`

```go
package postgres

import (
	"time"

	"github.com/google/uuid"
)

type PdpAConsentModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	SessionID     string     `gorm:"type:varchar(255);index"`
	PurposeCode   string     `gorm:"type:varchar(50);not null;index"`
	Status        string     `gorm:"type:varchar(20);not null;index"`
	IPAddress     string     `gorm:"type:varchar(45)"`
	UserAgent     string     `gorm:"type:text"`
	GrantedAt     time.Time  `gorm:"not null;index"`
	ExpiresAt     time.Time  `gorm:"not null"`
	RevokedAt     *time.Time `gorm:"index"`
	AutoDeletedAt *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (PdpAConsentModel) TableName() string { return "pdpa_consents" }

type PdpADSARModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	RequestType     string     `gorm:"type:varchar(30);not null;index"`
	Status          string     `gorm:"type:varchar(20);not null;index"`
	RequestedAt     time.Time  `gorm:"not null;index"`
	CompletedAt     *time.Time
	DataPayload     []byte     `gorm:"type:jsonb"`
	DataHash        string     `gorm:"type:varchar(128)"`
	BlockchainTx    string     `gorm:"type:varchar(128)"`
	RejectionReason string     `gorm:"type:text"`
	OTPHash         string     `gorm:"type:varchar(128)"`
	OTPExpiresAt    *time.Time
	OTPVerified     bool       `gorm:"default:false"`
	OTPAttempts     int        `gorm:"default:0"`
	IPAddress       string     `gorm:"type:varchar(45)"`
	UserAgent       string     `gorm:"type:text"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
}

func (PdpADSARModel) TableName() string { return "pdpa_dsar_requests" }

type PdpAAccountStatusModel struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	Status              string     `gorm:"type:varchar(20);not null;index"`
	SuspendedAt         *time.Time
	TerminatedAt        *time.Time
	DeletionConfirmedAt *time.Time
	RetentionDeadline   *time.Time `gorm:"index"`
	AutoDeletedAt       *time.Time
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}

func (PdpAAccountStatusModel) TableName() string { return "pdpa_user_account_statuses" }

type PdpAAuditTrailModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID `gorm:"type:uuid;index"`
	Action    string     `gorm:"type:varchar(80);not null;index"`
	Details   []byte     `gorm:"type:jsonb"`
	IPAddress string     `gorm:"type:varchar(45)"`
	UserAgent string     `gorm:"type:text"`
	CreatedAt time.Time  `gorm:"not null;index"`
}

func (PdpAAuditTrailModel) TableName() string { return "pdpa_audit_trails" }

type PdpAPolicyModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Version       string    `gorm:"type:varchar(30);not null;uniqueIndex"`
	Title         string    `gorm:"type:varchar(255);not null"`
	Content       string    `gorm:"type:text;not null"`
	EffectiveDate time.Time `gorm:"not null"`
	IsActive      bool      `gorm:"default:false;index"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (PdpAPolicyModel) TableName() string { return "pdpa_privacy_policies" }

type PdpAOutboxModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	AggregateType string     `gorm:"type:varchar(60);not null;index"`
	AggregateID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	EventType     string     `gorm:"type:varchar(120);not null;index"`
	Topic         string     `gorm:"type:varchar(120);not null"`
	Payload       []byte     `gorm:"type:jsonb;not null"`
	Status        string     `gorm:"type:varchar(20);not null;index"`
	RetryCount    int        `gorm:"default:0"`
	LastError     string     `gorm:"type:text"`
	CreatedAt     time.Time  `gorm:"not null;index"`
	PublishedAt   *time.Time
}

func (PdpAOutboxModel) TableName() string { return "pdpa_outbox" }
```

### 9.2 `infrastructure/persistence/postgres/uow.go`

```go
package postgres

import (
	"context"

	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/repository"
)

type UnitOfWorkImpl struct{ db *gorm.DB }

func NewUnitOfWork(db *gorm.DB) *UnitOfWorkImpl { return &UnitOfWorkImpl{db: db} }

type gormTxKey struct{}

func (u *UnitOfWorkImpl) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return fn(ctx)
	}
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, gormTxKey{}, tx))
	})
}

func GetDB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}

var _ repository.UnitOfWork = (*UnitOfWorkImpl)(nil)
```

### 9.3 `infrastructure/persistence/postgres/consent_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentRepoImpl struct{ db *gorm.DB }

func NewConsentRepository(db *gorm.DB) *ConsentRepoImpl { return &ConsentRepoImpl{db: db} }

func (r *ConsentRepoImpl) Save(ctx context.Context, c *entity.ConsentLog) error {
	return GetDB(ctx, r.db).Create(toConsentModel(c)).Error
}

func (r *ConsentRepoImpl) Update(ctx context.Context, c *entity.ConsentLog) error {
	res := GetDB(ctx, r.db).Model(&PdpAConsentModel{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"status":          c.Status.String(),
			"revoked_at":      c.RevokedAt,
			"auto_deleted_at": c.AutoDeletedAt,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrConsentNotFound
	}
	return nil
}

func (r *ConsentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	if err := GetDB(ctx, r.db).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrConsentNotFound
		}
		return nil, err
	}
	return toConsentEntity(&m), nil
}

func (r *ConsentRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error) {
	var ms []PdpAConsentModel
	if err := GetDB(ctx, r.db).Where("user_id = ?", userID).Order("granted_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	return toConsentEntities(ms), nil
}

func (r *ConsentRepoImpl) FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	err := GetDB(ctx, r.db).Where("user_id = ? AND purpose_code = ?", userID, purpose.String()).
		Order("granted_at DESC").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toConsentEntity(&m), nil
}

func (r *ConsentRepoImpl) FindActiveByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	err := GetDB(ctx, r.db).
		Where("user_id = ? AND purpose_code = ? AND status = ? AND expires_at > ?",
			userID, purpose.String(), valueobject.ConsentGranted, time.Now().UTC()).
		Order("granted_at DESC").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toConsentEntity(&m), nil
}

func (r *ConsentRepoImpl) FindReadyForAutoDeletion(ctx context.Context, cutoff time.Time, limit int) ([]entity.ConsentLog, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAConsentModel
	err := GetDB(ctx, r.db).
		Where("status IN ? AND revoked_at IS NOT NULL AND revoked_at <= ?",
			[]string{valueobject.ConsentRevoked.String(), valueobject.ConsentExpired.String()}, cutoff).
		Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return toConsentEntities(ms), nil
}

func (r *ConsentRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return GetDB(ctx, r.db).Where("user_id = ?", userID).Delete(&PdpAConsentModel{}).Error
}

func (r *ConsentRepoImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Where("id = ?", id).Delete(&PdpAConsentModel{}).Error
}

func (r *ConsentRepoImpl) IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error) {
	var count int64
	err := GetDB(ctx, r.db).Model(&PdpAConsentModel{}).
		Where("user_id = ? AND purpose_code = ? AND status = ? AND expires_at > ?",
			userID, purpose.String(), valueobject.ConsentGranted, time.Now().UTC()).
		Count(&count).Error
	return count > 0, err
}

// mappers
func toConsentModel(e *entity.ConsentLog) *PdpAConsentModel {
	return &PdpAConsentModel{
		ID: e.ID, UserID: e.UserID, SessionID: e.SessionID,
		PurposeCode: e.Purpose.String(), Status: e.Status.String(),
		IPAddress: e.IPAddress, UserAgent: e.UserAgent,
		GrantedAt: e.ConsentedAt, ExpiresAt: e.ExpiresAt,
		RevokedAt: e.RevokedAt, AutoDeletedAt: e.AutoDeletedAt,
	}
}

func toConsentEntity(m *PdpAConsentModel) *entity.ConsentLog {
	return &entity.ConsentLog{
		ID: m.ID, UserID: m.UserID, SessionID: m.SessionID,
		Purpose: valueobject.ConsentPurpose(m.PurposeCode),
		Status:  valueobject.ConsentStatus(m.Status),
		IPAddress: m.IPAddress, UserAgent: m.UserAgent,
		ConsentedAt: m.GrantedAt, ExpiresAt: m.ExpiresAt,
		RevokedAt: m.RevokedAt, AutoDeletedAt: m.AutoDeletedAt,
	}
}

func toConsentEntities(ms []PdpAConsentModel) []entity.ConsentLog {
	out := make([]entity.ConsentLog, len(ms))
	for i := range ms {
		out[i] = *toConsentEntity(&ms[i])
	}
	return out
}
```

### 9.4 `infrastructure/persistence/postgres/dsar_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DSARRepoImpl struct{ db *gorm.DB }

func NewDSARRepository(db *gorm.DB) *DSARRepoImpl { return &DSARRepoImpl{db: db} }

func (r *DSARRepoImpl) Save(ctx context.Context, d *entity.DSARRequest) error {
	return GetDB(ctx, r.db).Create(toDSARModel(d)).Error
}

func (r *DSARRepoImpl) Update(ctx context.Context, d *entity.DSARRequest) error {
	m := toDSARModel(d)
	return GetDB(ctx, r.db).Model(&PdpADSARModel{}).Where("id = ?", d.ID).
		Updates(map[string]interface{}{
			"status":           m.Status,
			"completed_at":     m.CompletedAt,
			"data_payload":     m.DataPayload,
			"data_hash":        m.DataHash,
			"blockchain_tx":    m.BlockchainTx,
			"rejection_reason": m.RejectionReason,
			"otp_verified":     m.OTPVerified,
			"otp_attempts":     m.OTPAttempts,
		}).Error
}

func (r *DSARRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error) {
	var m PdpADSARModel
	if err := GetDB(ctx, r.db).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrDSARNotFound
		}
		return nil, err
	}
	return toDSAREntity(&m), nil
}

func (r *DSARRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error) {
	var ms []PdpADSARModel
	if err := GetDB(ctx, r.db).Where("user_id = ?", userID).Order("requested_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	return toDSAREntities(ms), nil
}

func (r *DSARRepoImpl) FindByStatus(ctx context.Context, status valueobject.DSARStatus, limit, offset int) ([]entity.DSARRequest, error) {
	if limit <= 0 {
		limit = 50
	}
	var ms []PdpADSARModel
	err := GetDB(ctx, r.db).Where("status = ?", status.String()).
		Order("requested_at ASC").Limit(limit).Offset(offset).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return toDSAREntities(ms), nil
}

func (r *DSARRepoImpl) CountByUserAndTypeInPeriod(ctx context.Context, userID uuid.UUID, t valueobject.DSARType, from time.Time) (int64, error) {
	var count int64
	err := GetDB(ctx, r.db).Model(&PdpADSARModel{}).
		Where("user_id = ? AND request_type = ? AND requested_at >= ?", userID, t.String(), from).
		Count(&count).Error
	return count, err
}

func (r *DSARRepoImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return GetDB(ctx, r.db).Where("user_id = ?", userID).Delete(&PdpADSARModel{}).Error
}

func toDSARModel(e *entity.DSARRequest) *PdpADSARModel {
	return &PdpADSARModel{
		ID: e.ID, UserID: e.UserID,
		RequestType: e.RequestType.String(), Status: e.Status.String(),
		RequestedAt: e.RequestedAt, CompletedAt: e.CompletedAt,
		DataPayload: e.DataPayload, DataHash: e.DataHash, BlockchainTx: e.BlockchainTx,
		RejectionReason: e.RejectionReason,
		OTPHash: e.OTPHash, OTPExpiresAt: e.OTPExpiresAt,
		OTPVerified: e.OTPVerified, OTPAttempts: e.OTPAttempts,
		IPAddress: e.IPAddress, UserAgent: e.UserAgent,
	}
}

func toDSAREntity(m *PdpADSARModel) *entity.DSARRequest {
	return &entity.DSARRequest{
		ID: m.ID, UserID: m.UserID,
		RequestType: valueobject.DSARType(m.RequestType),
		Status:      valueobject.DSARStatus(m.Status),
		RequestedAt: m.RequestedAt, CompletedAt: m.CompletedAt,
		DataPayload: m.DataPayload, DataHash: m.DataHash, BlockchainTx: m.BlockchainTx,
		RejectionReason: m.RejectionReason,
		OTPHash: m.OTPHash, OTPExpiresAt: m.OTPExpiresAt,
		OTPVerified: m.OTPVerified, OTPAttempts: m.OTPAttempts,
		IPAddress: m.IPAddress, UserAgent: m.UserAgent,
	}
}

func toDSAREntities(ms []PdpADSARModel) []entity.DSARRequest {
	out := make([]entity.DSARRequest, len(ms))
	for i := range ms {
		out[i] = *toDSAREntity(&ms[i])
	}
	return out
}
```

### 9.5 `infrastructure/persistence/postgres/account_status_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type AccountStatusRepoImpl struct{ db *gorm.DB }

func NewAccountStatusRepository(db *gorm.DB) *AccountStatusRepoImpl { return &AccountStatusRepoImpl{db: db} }

func (r *AccountStatusRepoImpl) Save(ctx context.Context, s *entity.UserAccountStatus) error {
	return GetDB(ctx, r.db).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"status", "suspended_at", "terminated_at",
			"deletion_confirmed_at", "retention_deadline",
			"auto_deleted_at", "updated_at",
		}),
	}).Create(toAccountModel(s)).Error
}

func (r *AccountStatusRepoImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error) {
	var m PdpAAccountStatusModel
	if err := GetDB(ctx, r.db).First(&m, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrAccountNotFound
		}
		return nil, err
	}
	return toAccountEntity(&m), nil
}

func (r *AccountStatusRepoImpl) FindReadyForAutoDeletion(ctx context.Context, now time.Time, limit int) ([]entity.UserAccountStatus, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAAccountStatusModel
	err := GetDB(ctx, r.db).
		Where("status = ? AND retention_deadline IS NOT NULL AND retention_deadline <= ?",
			valueobject.AccountSuspended.String(), now).
		Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.UserAccountStatus, len(ms))
	for i := range ms {
		out[i] = *toAccountEntity(&ms[i])
	}
	return out, nil
}

func (r *AccountStatusRepoImpl) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := GetDB(ctx, r.db).Model(&PdpAAccountStatusModel{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

func toAccountModel(e *entity.UserAccountStatus) *PdpAAccountStatusModel {
	return &PdpAAccountStatusModel{
		ID: e.ID, UserID: e.UserID, Status: e.Status.String(),
		SuspendedAt: e.SuspendedAt, TerminatedAt: e.TerminatedAt,
		DeletionConfirmedAt: e.DeletionConfirmedAt,
		RetentionDeadline:   e.RetentionDeadline, AutoDeletedAt: e.AutoDeletedAt,
	}
}

func toAccountEntity(m *PdpAAccountStatusModel) *entity.UserAccountStatus {
	return &entity.UserAccountStatus{
		ID: m.ID, UserID: m.UserID,
		Status:      valueobject.AccountStatus(m.Status),
		SuspendedAt: m.SuspendedAt, TerminatedAt: m.TerminatedAt,
		DeletionConfirmedAt: m.DeletionConfirmedAt,
		RetentionDeadline:   m.RetentionDeadline, AutoDeletedAt: m.AutoDeletedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
```

### 9.6 `infrastructure/persistence/postgres/audit_repo_impl.go`

```go
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type AuditRepoImpl struct{ db *gorm.DB }

func NewAuditRepository(db *gorm.DB) *AuditRepoImpl { return &AuditRepoImpl{db: db} }

func (r *AuditRepoImpl) Save(ctx context.Context, a *entity.AuditTrail) error {
	details, _ := json.Marshal(a.Details)
	return GetDB(ctx, r.db).Create(&PdpAAuditTrailModel{
		ID: a.ID, UserID: a.UserID, Action: a.Action,
		Details: details, IPAddress: a.IPAddress, UserAgent: a.UserAgent,
		CreatedAt: a.CreatedAt,
	}).Error
}

func (r *AuditRepoImpl) FindByFilter(ctx context.Context, f repository.AuditFilter) ([]entity.AuditTrail, error) {
	q := GetDB(ctx, r.db).Model(&PdpAAuditTrailModel{})
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.Action != "" {
		q = q.Where("action = ?", f.Action)
	}
	if f.FromTime != nil {
		q = q.Where("created_at >= ?", *f.FromTime)
	}
	if f.ToTime != nil {
		q = q.Where("created_at <= ?", *f.ToTime)
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAAuditTrailModel
	if err := q.Order("created_at DESC").Limit(limit).Offset(f.Offset).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]entity.AuditTrail, len(ms))
	for i, m := range ms {
		var details interface{}
		_ = json.Unmarshal(m.Details, &details)
		out[i] = entity.AuditTrail{
			ID: m.ID, UserID: m.UserID, Action: m.Action,
			Details: details, IPAddress: m.IPAddress,
			UserAgent: m.UserAgent, CreatedAt: m.CreatedAt,
		}
	}
	return out, nil
}

func (r *AuditRepoImpl) CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	type row struct {
		Action string
		Cnt    int64
	}
	var rows []row
	err := GetDB(ctx, r.db).Model(&PdpAAuditTrailModel{}).
		Select("action, COUNT(*) as cnt").
		Where("created_at BETWEEN ? AND ?", from, to).
		Group("action").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		out[r.Action] = r.Cnt
	}
	return out, nil
}
```

### 9.7 `infrastructure/persistence/postgres/policy_repo_impl.go`

```go
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
)

type PolicyRepoImpl struct{ db *gorm.DB }

func NewPolicyRepository(db *gorm.DB) *PolicyRepoImpl { return &PolicyRepoImpl{db: db} }

func (r *PolicyRepoImpl) Save(ctx context.Context, p *entity.PrivacyPolicy) error {
	return GetDB(ctx, r.db).Create(&PdpAPolicyModel{
		ID: p.ID, Version: p.Version, Title: p.Title, Content: p.Content,
		EffectiveDate: p.EffectiveDate, IsActive: p.IsActive,
	}).Error
}

func (r *PolicyRepoImpl) Update(ctx context.Context, p *entity.PrivacyPolicy) error {
	return GetDB(ctx, r.db).Model(&PdpAPolicyModel{}).Where("id = ?", p.ID).
		Updates(map[string]interface{}{
			"title": p.Title, "content": p.Content,
			"effective_date": p.EffectiveDate, "is_active": p.IsActive,
		}).Error
}

func (r *PolicyRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := GetDB(ctx, r.db).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

func (r *PolicyRepoImpl) FindActive(ctx context.Context) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := GetDB(ctx, r.db).Where("is_active = ?", true).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

func (r *PolicyRepoImpl) FindAll(ctx context.Context) ([]entity.PrivacyPolicy, error) {
	var ms []PdpAPolicyModel
	if err := GetDB(ctx, r.db).Order("effective_date DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]entity.PrivacyPolicy, len(ms))
	for i := range ms {
		out[i] = *toPolicyEntity(&ms[i])
	}
	return out, nil
}

func (r *PolicyRepoImpl) FindByVersion(ctx context.Context, version string) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := GetDB(ctx, r.db).Where("version = ?", version).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

func (r *PolicyRepoImpl) DeactivateAll(ctx context.Context) error {
	return GetDB(ctx, r.db).Model(&PdpAPolicyModel{}).
		Where("is_active = ?", true).Update("is_active", false).Error
}

func toPolicyEntity(m *PdpAPolicyModel) *entity.PrivacyPolicy {
	return &entity.PrivacyPolicy{
		ID: m.ID, Version: m.Version, Title: m.Title, Content: m.Content,
		EffectiveDate: m.EffectiveDate, IsActive: m.IsActive,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
```

### 9.8 `infrastructure/persistence/postgres/outbox_repo_impl.go`

```go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"icmongolang/internal/modules/pdpa/domain/repository"
)

type OutboxRepoImpl struct{ db *gorm.DB }

func NewOutboxRepository(db *gorm.DB) *OutboxRepoImpl { return &OutboxRepoImpl{db: db} }

func (r *OutboxRepoImpl) Save(ctx context.Context, e *repository.OutboxEvent) error {
	return GetDB(ctx, r.db).Create(&PdpAOutboxModel{
		ID: e.ID, AggregateType: e.AggregateType, AggregateID: e.AggregateID,
		EventType: e.EventType, Topic: e.Topic, Payload: e.Payload,
		Status: e.Status, RetryCount: e.RetryCount, LastError: e.LastError,
		CreatedAt: e.CreatedAt, PublishedAt: e.PublishedAt,
	}).Error
}

// FetchPending ใช้ SKIP LOCKED — ปลอดภัยเมื่อมีหลาย replica
func (r *OutboxRepoImpl) FetchPending(ctx context.Context, limit int) ([]repository.OutboxEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAOutboxModel
	err := GetDB(ctx, r.db).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ?", repository.OutboxPending).
		Order("created_at ASC").Limit(limit).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]repository.OutboxEvent, len(ms))
	for i := range ms {
		out[i] = repository.OutboxEvent{
			ID: ms[i].ID, AggregateType: ms[i].AggregateType,
			AggregateID: ms[i].AggregateID, EventType: ms[i].EventType,
			Topic: ms[i].Topic, Payload: ms[i].Payload,
			Status: ms[i].Status, RetryCount: ms[i].RetryCount,
			LastError: ms[i].LastError, CreatedAt: ms[i].CreatedAt,
			PublishedAt: ms[i].PublishedAt,
		}
	}
	return out, nil
}

func (r *OutboxRepoImpl) MarkPublished(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return GetDB(ctx, r.db).Model(&PdpAOutboxModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": repository.OutboxPublished, "published_at": now,
		}).Error
}

func (r *OutboxRepoImpl) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return GetDB(ctx, r.db).Model(&PdpAOutboxModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": repository.OutboxFailed, "last_error": errMsg,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error
}

func (r *OutboxRepoImpl) CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error) {
	res := GetDB(ctx, r.db).
		Where("status = ? AND published_at < ?", repository.OutboxPublished, cutoff).
		Delete(&PdpAOutboxModel{})
	return res.RowsAffected, res.Error
}
```

---

## 📁 10. Infrastructure — Redis Cache

### 10.1 `infrastructure/cache/redis/consent_cache.go`

```go
package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentCacheImpl struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

func NewConsentCache(client *redis.Client, ttl time.Duration) *ConsentCacheImpl {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ConsentCacheImpl{client: client, ttl: ttl, prefix: "pdpa"}
}

func (c *ConsentCacheImpl) key(userID uuid.UUID, purpose valueobject.ConsentPurpose) string {
	return fmt.Sprintf("%s:consent:%s:%s", c.prefix, userID.String(), purpose.String())
}

func (c *ConsentCacheImpl) Set(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose, status valueobject.ConsentStatus) error {
	return c.client.Set(ctx, c.key(userID, purpose), status.String(), c.ttl).Err()
}

func (c *ConsentCacheImpl) Get(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (valueobject.ConsentStatus, error) {
	v, err := c.client.Get(ctx, c.key(userID, purpose)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return valueobject.ConsentStatus(v), nil
}

func (c *ConsentCacheImpl) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("%s:consent:%s:*", c.prefix, userID.String())
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
```

---

## 📁 11. Infrastructure — Kafka

### 11.1 `infrastructure/messaging/kafka/producer.go`

```go
package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"icmongolang/pkg/logger"
)

// Producer abstraction
type Producer interface {
	Publish(ctx context.Context, topic, key string, payload []byte) error
	Close() error
}

type SyncProducerImpl struct {
	producer sarama.SyncProducer
	log      logger.Logger
}

func NewSyncProducer(brokers []string, log logger.Logger) (*SyncProducerImpl, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &SyncProducerImpl{producer: p, log: log}, nil
}

func (p *SyncProducerImpl) Publish(_ context.Context, topic, key string, payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Errorw("kafka publish failed", "topic", topic, "error", err)
	}
	return err
}

func (p *SyncProducerImpl) Close() error { return p.producer.Close() }
```

### 11.2 `infrastructure/messaging/kafka/event_serializer.go`

```go
package kafka

import (
	"encoding/json"

	"icmongolang/internal/modules/pdpa/domain/event"
)

// SerializeEnvelope แปลง domain envelope → JSON bytes
func SerializeEnvelope(env event.Envelope) ([]byte, error) {
	return json.Marshal(env)
}

// DeserializeEnvelope parse JSON → domain envelope
func DeserializeEnvelope(data []byte) (*event.Envelope, error) {
	var env event.Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return &env, nil
}
```

---

## 📁 12. Infrastructure — Outbox Publisher

### 12.1 `infrastructure/messaging/outbox/outbox_publisher.go`

```go
package outbox

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type OutboxPublisher struct {
	outboxRepo repository.OutboxRepository
	producer   kafka.Producer
	interval   time.Duration
	batchSize  int
	logger     logger.Logger
}

func NewOutboxPublisher(
	outboxRepo repository.OutboxRepository,
	producer kafka.Producer,
	interval time.Duration,
	batchSize int,
	logger logger.Logger,
) *OutboxPublisher {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	return &OutboxPublisher{
		outboxRepo: outboxRepo, producer: producer,
		interval: interval, batchSize: batchSize, logger: logger,
	}
}

func (p *OutboxPublisher) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.logger.Infow("outbox publisher starting", "interval", p.interval, "batch_size", p.batchSize)
	for {
		select {
		case <-ctx.Done():
			p.logger.Infow("outbox publisher stopping")
			return
		case <-ticker.C:
			p.drain(ctx)
		}
	}
}

func (p *OutboxPublisher) drain(ctx context.Context) {
	events, err := p.outboxRepo.FetchPending(ctx, p.batchSize)
	if err != nil {
		p.logger.Errorw("fetch outbox failed", "error", err)
		return
	}
	for i := range events {
		e := &events[i]
		if err := p.producer.Publish(ctx, e.Topic, e.AggregateID.String(), e.Payload); err != nil {
			_ = p.outboxRepo.MarkFailed(ctx, e.ID, err.Error())
			continue
		}
		_ = p.outboxRepo.MarkPublished(ctx, e.ID)
	}
}

func (p *OutboxPublisher) Cleanup(ctx context.Context) {
	cutoff := time.Now().UTC().Add(-7 * 24 * time.Hour)
	n, err := p.outboxRepo.CleanupPublished(ctx, cutoff)
	if err != nil {
		p.logger.Errorw("outbox cleanup failed", "error", err)
		return
	}
	p.logger.Infow("outbox cleanup completed", "deleted", n)
}
```

---

## 📁 13. Infrastructure — External Clients

### 13.1 `infrastructure/external/email_client.go`

```go
package external

import (
	"context"

	"icmongolang/pkg/logger"
)

type SMTPEmailClient struct {
	from string
	log  logger.Logger
}

func NewEmailClient(from string, log logger.Logger) *SMTPEmailClient {
	return &SMTPEmailClient{from: from, log: log}
}

func (c *SMTPEmailClient) Send(_ context.Context, to, subject, template, locale string, vars map[string]interface{}) error {
	// TODO: implement real SMTP
	c.log.Infow("sending email",
		"to", to, "subject", subject, "template", template,
		"locale", locale, "vars_count", len(vars))
	return nil
}
```

### 13.2 `infrastructure/external/llm_client.go`

```go
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/pkg/logger"
)

type HTTPLLMClient struct {
	baseURL string
	apiKey  string
	model   string
	timeout time.Duration
	log     logger.Logger
}

func NewLLMClient(baseURL, apiKey, model string, log logger.Logger) *HTTPLLMClient {
	return &HTTPLLMClient{
		baseURL: baseURL, apiKey: apiKey, model: model,
		timeout: 60 * time.Second, log: log,
	}
}

func (c *HTTPLLMClient) AnalyzeDataWithContext(_ context.Context, analysisType string, payload map[string]interface{}) (string, error) {
	// TODO: real HTTP call + circuit breaker
	body, _ := json.Marshal(payload)
	c.log.Infow("LLM analysis requested", "type", analysisType, "payload_size", len(body))
	return fmt.Sprintf("[LLM-STUB] analysis=%s accepted", analysisType), nil
}
```

### 13.3 `infrastructure/external/blockchain_client.go`

```go
package external

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"icmongolang/pkg/logger"
)

type MockBlockchainClient struct{ log logger.Logger }

func NewBlockchainClient(log logger.Logger) *MockBlockchainClient {
	return &MockBlockchainClient{log: log}
}

func (c *MockBlockchainClient) Record(_ context.Context, userID, action, dataHash string) (string, error) {
	// TODO: real blockchain (Hyperledger/Ethereum)
	input := fmt.Sprintf("%s|%s|%s|%d", userID, action, dataHash, time.Now().UnixNano())
	h := sha256.Sum256([]byte(input))
	txHash := hex.EncodeToString(h[:])
	c.log.Infow("blockchain record created", "tx_hash", txHash, "action", action)
	return txHash, nil
}
```

---

## 📁 14. Interface Layer — HTTP

### 14.1 `interfaces/http/dto.go`

```go
package http

import "github.com/google/uuid"

type RecordConsentRequest struct {
	Purposes  map[string]bool `json:"purposes" validate:"required"`
	SessionID string          `json:"session_id,omitempty"`
}

type RevokeConsentRequest struct {
	Purpose string `json:"purpose" validate:"required"`
}

type SubmitDSARRequest struct {
	RequestType string `json:"request_type" validate:"required"`
}

type ProcessDSARRequest struct {
	Action       string `json:"action" validate:"required,oneof=VERIFY_OTP APPROVE REJECT"`
	OTPCode      string `json:"otp_code,omitempty"`
	RejectReason string `json:"reject_reason,omitempty"`
}

type ConfirmDeletionRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

type SuspendAccountRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required,uuid"`
	RetentionYears int       `json:"retention_years,omitempty"`
}

type TerminateAccountRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

type PublishPolicyRequest struct {
	Version       string `json:"version" validate:"required"`
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	EffectiveDate string `json:"effective_date" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
```

### 14.2 `interfaces/http/consent_handler.go`

```go
package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httputil"
)

var errUnauthorized = errors.New("unauthorized")

type ConsentHandler struct {
	record  *command.RecordConsentHandler
	revoke  *command.RevokeConsentHandler
	history *query.GetConsentHistoryHandler
}

func NewConsentHandler(
	record *command.RecordConsentHandler,
	revoke *command.RevokeConsentHandler,
	history *query.GetConsentHistoryHandler,
) *ConsentHandler {
	return &ConsentHandler{record: record, revoke: revoke, history: history}
}

func (h *ConsentHandler) Record(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body RecordConsentRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// บันทึกทุก purpose
	for raw, granted := range body.Purposes {
		if !granted {
			continue
		}
		if err := h.record.Handle(r.Context(), dto.GrantConsentRequest{
			UserID:        userID,
			Purpose:       raw,
			IPAddress:     r.RemoteAddr,
			UserAgent:     r.UserAgent(),
			CorrelationID: r.Header.Get("X-Correlation-ID"),
			TraceID:       r.Header.Get("X-Trace-ID"),
		}); err != nil {
			httputil.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "consent recorded"})
}

func (h *ConsentHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body RevokeConsentRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.revoke.Handle(r.Context(), dto.RevokeConsentRequest{
		UserID:        userID,
		Purpose:       body.Purpose,
		IPAddress:     r.RemoteAddr,
		UserAgent:     r.UserAgent(),
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "consent revoked"})
}

func (h *ConsentHandler) History(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	res, err := h.history.Handle(r.Context(), userID)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func extractUserID(r *http.Request) (uuid.UUID, error) {
	v := r.Context().Value("user_id")
	if v == nil {
		return uuid.Nil, errUnauthorized
	}
	switch id := v.(type) {
	case uuid.UUID:
		return id, nil
	case string:
		return uuid.Parse(id)
	}
	return uuid.Nil, errUnauthorized
}

var _ = chi.NewRouter
```

### 14.3 `interfaces/http/dsar_handler.go`

```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httputil"
)

type DSARHandler struct {
	submit  *command.SubmitDSARHandler
	process *command.ProcessDSARHandler
	getOne  *query.GetDSARStatusHandler
	list    *query.ListDSARHandler
}

func NewDSARHandler(
	submit *command.SubmitDSARHandler,
	process *command.ProcessDSARHandler,
	getOne *query.GetDSARStatusHandler,
	list *query.ListDSARHandler,
) *DSARHandler {
	return &DSARHandler{submit: submit, process: process, getOne: getOne, list: list}
}

func (h *DSARHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body SubmitDSARRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	dsar, err := h.submit.Handle(r.Context(), dto.SubmitDSARRequest{
		UserID:        userID,
		RequestType:   body.RequestType,
		IPAddress:     r.RemoteAddr,
		UserAgent:     r.UserAgent(),
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	})
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":           dsar.ID,
		"status":       dsar.Status.String(),
		"requested_at": dsar.RequestedAt,
	})
}

func (h *DSARHandler) Process(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	userID, _ := extractUserID(r)
	if err := h.process.Handle(r.Context(), dto.ProcessDSARRequest{
		UserID:        userID,
		DSARRequestID: id,
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "processing"})
}

func (h *DSARHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	res, err := h.getOne.Handle(r.Context(), id)
	if err != nil {
		httputil.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func (h *DSARHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	res, err := h.list.Handle(r.Context(), userID)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}
```

### 14.4 `interfaces/http/account_handler.go`

```go
package http

import (
	"net/http"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/pkg/httputil"
)

type AccountHandler struct {
	suspend   *command.HandleAccountSuspendedHandler
	terminate *command.HandleAccountTerminatedHandler
	confirm   *command.ConfirmDeletionHandler
}

func NewAccountHandler(
	suspend *command.HandleAccountSuspendedHandler,
	terminate *command.HandleAccountTerminatedHandler,
	confirm *command.ConfirmDeletionHandler,
) *AccountHandler {
	return &AccountHandler{suspend: suspend, terminate: terminate, confirm: confirm}
}

func (h *AccountHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	var body SuspendAccountRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if body.RetentionYears < 1 {
		body.RetentionYears = 1
	}
	if err := h.suspend.Handle(r.Context(), dto.SuspendAccountRequest{
		UserID:         body.UserID,
		RetentionYears: body.RetentionYears,
		CorrelationID:  r.Header.Get("X-Correlation-ID"),
		TraceID:        r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "account suspended"})
}

func (h *AccountHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	var body TerminateAccountRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.terminate.Handle(r.Context(), dto.TerminateAccountRequest{
		UserID:        body.UserID,
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "account terminated"})
}

func (h *AccountHandler) ConfirmDeletion(w http.ResponseWriter, r *http.Request) {
	var body ConfirmDeletionRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.confirm.Handle(r.Context(), dto.ConfirmDeletionRequest{
		UserID:        body.UserID,
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "deletion confirmed"})
}
```

### 14.5 `interfaces/http/policy_handler.go`

```go
package http

import (
	"net/http"
	"time"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httputil"
)

type PolicyHandler struct {
	publish *command.PublishPrivacyPolicyHandler
	active  *query.GetActivePrivacyPolicyHandler
	list    *query.ListPoliciesHandler
}

func NewPolicyHandler(
	publish *command.PublishPrivacyPolicyHandler,
	active *query.GetActivePrivacyPolicyHandler,
	list *query.ListPoliciesHandler,
) *PolicyHandler {
	return &PolicyHandler{publish: publish, active: active, list: list}
}

func (h *PolicyHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	res, err := h.active.Handle(r.Context())
	if err != nil {
		httputil.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func (h *PolicyHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	res, err := h.list.Handle(r.Context())
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func (h *PolicyHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var body PublishPolicyRequest
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	eff, err := time.Parse(time.RFC3339, body.EffectiveDate)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid effective_date")
		return
	}
	policy, err := h.publish.Handle(r.Context(), dto.PublishPolicyRequest{
		Version:       body.Version,
		Title:         body.Title,
		Content:       body.Content,
		EffectiveDate: eff,
		CorrelationID: r.Header.Get("X-Correlation-ID"),
		TraceID:       r.Header.Get("X-Trace-ID"),
	})
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"policy_id": policy.ID,
		"version":   policy.Version,
	})
}
```

### 14.6 `interfaces/http/admin_handler.go`

```go
package http

import (
	"net/http"
	"time"

	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/httputil"
)

type AdminHandler struct {
	report *query.GetAdminReportHandler
	audit  *query.ListAuditTrailsHandler
}

func NewAdminHandler(report *query.GetAdminReportHandler, audit *query.ListAuditTrailsHandler) *AdminHandler {
	return &AdminHandler{report: report, audit: audit}
}

func (h *AdminHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	from := parseTime(r.URL.Query().Get("from"), time.Now().AddDate(0, -1, 0))
	to := parseTime(r.URL.Query().Get("to"), time.Now())
	res, err := h.report.Handle(r.Context(), from, to)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func (h *AdminHandler) ListAudit(w http.ResponseWriter, r *http.Request) {
	filter := repository.AuditFilter{
		Action: r.URL.Query().Get("action"),
		Limit:  parseInt(r.URL.Query().Get("limit"), 100),
		Offset: parseInt(r.URL.Query().Get("offset"), 0),
	}
	res, err := h.audit.Handle(r.Context(), filter)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func parseTime(s string, def time.Time) time.Time {
	if s == "" {
		return def
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return def
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return def
	}
	return n
}
```

> **เพิ่ม `import "fmt"`**

### 14.7 `interfaces/http/routes.go`

```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Consent *ConsentHandler
	DSAR    *DSARHandler
	Account *AccountHandler
	Policy  *PolicyHandler
	Admin   *AdminHandler
}

type Middleware func(http.Handler) http.Handler

func RegisterRoutes(r chi.Router, h *Handlers, auth, adminOnly Middleware) {
	r.Route("/api/v1/pdpa", func(r chi.Router) {
		r.Use(auth)

		// Consent
		r.Post("/consent", h.Consent.Record)
		r.Delete("/consent", h.Consent.Revoke)
		r.Get("/consent/history", h.Consent.History)

		// DSAR
		r.Post("/dsar", h.DSAR.Submit)
		r.Get("/dsar", h.DSAR.List)
		r.Get("/dsar/{id}", h.DSAR.Get)

		// Policy (public read — already under auth)
		r.Get("/policy", h.Policy.GetActive)
		r.Get("/policy/versions", h.Policy.ListVersions)

		// Admin-only
		r.Group(func(r chi.Router) {
			r.Use(adminOnly)

			r.Post("/dsar/{id}/process", h.DSAR.Process)

			r.Post("/account/suspend", h.Account.Suspend)
			r.Post("/account/terminate", h.Account.Terminate)
			r.Post("/account/confirm-deletion", h.Account.ConfirmDeletion)

			r.Post("/policy", h.Policy.Publish)

			r.Get("/admin/reports", h.Admin.GetReport)
			r.Get("/admin/audit-trails", h.Admin.ListAudit)
		})
	})
}
```

---

## 📁 15. Module Entry Point

### 15.1 `module.go`

```go
package pdpa

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/application/command"
	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	rediscache "icmongolang/internal/modules/pdpa/infrastructure/cache/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/external"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging/kafka"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging/outbox"
	pg "icmongolang/internal/modules/pdpa/infrastructure/persistence/postgres"
	httpiface "icmongolang/internal/modules/pdpa/interfaces/http"
	"icmongolang/pkg/logger"
)

type Config struct {
	KafkaBrokers      []string
	AnonymizationSalt string
	RetentionYears    int
	OutboxInterval    time.Duration
}

type Module struct {
	// Repos
	ConsentRepo repository.ConsentRepository
	DSARRepo    repository.DSARRepository
	AccountRepo repository.UserAccountStatusRepository
	AuditRepo   repository.AuditRepository
	PolicyRepo  repository.PolicyRepository
	OutboxRepo  repository.OutboxRepository
	UoW         repository.UnitOfWork

	// Domain services
	Validator  *service.ConsentValidator
	Policy     *service.DeletionPolicyService
	Anonymizer *service.AnonymizationService

	// Commands
	RecordConsent   *command.RecordConsentHandler
	RevokeConsent   *command.RevokeConsentHandler
	SubmitDSAR      *command.SubmitDSARHandler
	ProcessDSAR     *command.ProcessDSARHandler
	ConfirmDeletion *command.ConfirmDeletionHandler
	ImmediateDelete *command.ImmediateDeletionHandler
	AutoDelete      *command.AutoDeleteExpiredHandler
	SuspendAccount  *command.HandleAccountSuspendedHandler
	TerminateAcct   *command.HandleAccountTerminatedHandler
	PublishPolicy   *command.PublishPrivacyPolicyHandler

	// Queries
	GetConsentHistory *query.GetConsentHistoryHandler
	GetDSARStatus     *query.GetDSARStatusHandler
	ListDSAR          *query.ListDSARHandler
	GetAdminReport    *query.GetAdminReportHandler
	ListAuditTrails   *query.ListAuditTrailsHandler
	GetActivePolicy   *query.GetActivePrivacyPolicyHandler
	ListPolicies      *query.ListPoliciesHandler

	// Event handlers
	OnConsentGranted    *eventhandler.ConsentGrantedHandler
	OnConsentRevoked    *eventhandler.ConsentRevokedHandler
	OnDSARSubmitted     *eventhandler.DSARSubmittedHandler
	OnDataDeletion      *eventhandler.DataDeletionHandler
	OnAccountSuspended  *eventhandler.AccountSuspendedHandler
	OnLLMAnalysis       *eventhandler.LLMAnalysisHandler
	OnEmailNotification *eventhandler.EmailNotificationHandler
	OnBlockchainRecord  *eventhandler.BlockchainRecordHandler

	// Infra
	Producer        kafka.Producer
	OutboxPublisher *outbox.OutboxPublisher
	ConsentCache    *rediscache.ConsentCacheImpl

	// HTTP
	Handlers *httpiface.Handlers

	log logger.Logger
}

func NewModule(
	db *gorm.DB,
	redisClient *redis.Client,
	cfg Config,
	log logger.Logger,
) (*Module, error) {
	if cfg.RetentionYears < 1 {
		cfg.RetentionYears = 1
	}
	if cfg.AnonymizationSalt == "" {
		cfg.AnonymizationSalt = "change-me"
	}

	// Repos + UoW
	consentRepo := pg.NewConsentRepository(db)
	dsarRepo := pg.NewDSARRepository(db)
	accountRepo := pg.NewAccountStatusRepository(db)
	auditRepo := pg.NewAuditRepository(db)
	policyRepo := pg.NewPolicyRepository(db)
	outboxRepo := pg.NewOutboxRepository(db)
	uow := pg.NewUnitOfWork(db)

	// Services
	validator := service.NewConsentValidator()
	policySvc := service.NewDeletionPolicyService(service.RetentionPolicy{
		SuspendedRetentionYears: cfg.RetentionYears,
		RevokedConsentDays:      365,
	})
	anonymizer := service.NewAnonymizationService(cfg.AnonymizationSalt)

	// Kafka + Outbox
	producer, err := kafka.NewSyncProducer(cfg.KafkaBrokers, log)
	if err != nil {
		return nil, err
	}
	publisher := outbox.NewOutboxPublisher(outboxRepo, producer, cfg.OutboxInterval, 100, log)

	// Cache
	cache := rediscache.NewConsentCache(redisClient, 24*time.Hour)

	// External clients
	emailClient := external.NewEmailClient("noreply@example.com", log)
	llmClient := external.NewLLMClient("https://api.example.com", "key", "gpt-4", log)
	bcClient := external.NewBlockchainClient(log)

	// Idempotency store (Redis)
	idem := rediscache.NewIdempotencyStore(redisClient)

	// Commands
	recordUC := command.NewRecordConsentHandler(consentRepo, outboxRepo, uow, validator, time.Now, log)
	revokeUC := command.NewRevokeConsentHandler(consentRepo, outboxRepo, uow, time.Now, log)
	submitUC := command.NewSubmitDSARHandler(dsarRepo, outboxRepo, uow, time.Now, log)
	processUC := command.NewProcessDSARHandler(dsarRepo, outboxRepo, uow, log)
	confirmUC := command.NewConfirmDeletionHandler(accountRepo, outboxRepo, uow, policySvc, time.Now, log)
	immediateUC := command.NewImmediateDeletionHandler(consentRepo, dsarRepo, accountRepo, outboxRepo, uow, time.Now, log)
	autoDeleteUC := command.NewAutoDeleteExpiredHandler(consentRepo, dsarRepo, accountRepo, outboxRepo, uow, policySvc, time.Now, log)
	suspendUC := command.NewHandleAccountSuspendedHandler(accountRepo, outboxRepo, uow, time.Now, log)
	terminateUC := command.NewHandleAccountTerminatedHandler(accountRepo, outboxRepo, uow, time.Now, log)
	publishUC := command.NewPublishPrivacyPolicyHandler(policyRepo, outboxRepo, uow, time.Now, log)

	// Queries
	historyQ := query.NewGetConsentHistoryHandler(consentRepo)
	dsarStatusQ := query.NewGetDSARStatusHandler(dsarRepo)
	listDSARQ := query.NewListDSARHandler(dsarRepo)
	reportQ := query.NewGetAdminReportHandler(auditRepo)
	auditQ := query.NewListAuditTrailsHandler(auditRepo)
	activePolicyQ := query.NewGetActivePrivacyPolicyHandler(policyRepo)
	listPolicyQ := query.NewListPoliciesHandler(policyRepo)

	// Event handlers
	metrics := eventhandler.NopMetrics() // swap with Prometheus in prod
	onGranted := eventhandler.NewConsentGrantedHandler(auditRepo, cache, producer, idem, metrics, log)
	onRevoked := eventhandler.NewConsentRevokedHandler(auditRepo, cache, producer, idem, metrics, log)
	onDSARSub := eventhandler.NewDSARSubmittedHandler(auditRepo, producer, idem, metrics, log)
	onDataDel := eventhandler.NewDataDeletionHandler(auditRepo, producer, idem, metrics, log)
	onSusp := eventhandler.NewAccountSuspendedHandler(auditRepo, producer, idem, metrics, log)
	onLLM := eventhandler.NewLLMAnalysisHandler(auditRepo, llmClient, idem, metrics, log)
	onEmail := eventhandler.NewEmailNotificationHandler(auditRepo, emailClient, idem, metrics, log)
	onBC := eventhandler.NewBlockchainRecordHandler(auditRepo, bcClient, idem, metrics, log)

	// HTTP
	handlers := &httpiface.Handlers{
		Consent: httpiface.NewConsentHandler(recordUC, revokeUC, historyQ),
		DSAR:    httpiface.NewDSARHandler(submitUC, processUC, dsarStatusQ, listDSARQ),
		Account: httpiface.NewAccountHandler(suspendUC, terminateUC, confirmUC),
		Policy:  httpiface.NewPolicyHandler(publishUC, activePolicyQ, listPolicyQ),
		Admin:   httpiface.NewAdminHandler(reportQ, auditQ),
	}

	m := &Module{
		ConsentRepo: consentRepo, DSARRepo: dsarRepo, AccountRepo: accountRepo,
		AuditRepo: auditRepo, PolicyRepo: policyRepo, OutboxRepo: outboxRepo, UoW: uow,

		Validator: validator, Policy: policySvc, Anonymizer: anonymizer,

		RecordConsent: recordUC, RevokeConsent: revokeUC, SubmitDSAR: submitUC,
		ProcessDSAR: processUC, ConfirmDeletion: confirmUC, ImmediateDelete: immediateUC,
		AutoDelete: autoDeleteUC, SuspendAccount: suspendUC, TerminateAcct: terminateUC,
		PublishPolicy: publishUC,

		GetConsentHistory: historyQ, GetDSARStatus: dsarStatusQ, ListDSAR: listDSARQ,
		GetAdminReport: reportQ, ListAuditTrails: auditQ,
		GetActivePolicy: activePolicyQ, ListPolicies: listPolicyQ,

		OnConsentGranted: onGranted, OnConsentRevoked: onRevoked, OnDSARSubmitted: onDSARSub,
		OnDataDeletion: onDataDel, OnAccountSuspended: onSusp, OnLLMAnalysis: onLLM,
		OnEmailNotification: onEmail, OnBlockchainRecord: onBC,

		Producer: producer, OutboxPublisher: publisher, ConsentCache: cache,
		Handlers: handlers, log: log,
	}
	return m, nil
}

func (m *Module) RegisterHTTP(r chi.Router, auth, admin httpiface.Middleware) {
	httpiface.RegisterRoutes(r, m.Handlers, auth, admin)
}

func (m *Module) StartBackground(ctx context.Context) {
	go m.OutboxPublisher.Run(ctx)
}

func (m *Module) Shutdown() error {
	return m.Producer.Close()
}
```

---

## 📁 16. Infrastructure — Idempotency Store (Redis)

### 16.1 `infrastructure/cache/redis/idempotency_store.go`

```go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type IdempotencyStore struct {
	client *redis.Client
	prefix string
}

func NewIdempotencyStore(client *redis.Client) *IdempotencyStore {
	return &IdempotencyStore{client: client, prefix: "pdpa:idem:"}
}

func (s *IdempotencyStore) IsProcessed(ctx context.Context, eventID uuid.UUID, handler string) (bool, error) {
	n, err := s.client.Exists(ctx, s.key(eventID, handler)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID uuid.UUID, handler string, ttl time.Duration) error {
	_, err := s.client.SetNX(ctx, s.key(eventID, handler), "1", ttl).Result()
	return err
}

func (s *IdempotencyStore) key(eventID uuid.UUID, handler string) string {
	return fmt.Sprintf("%s%s:%s", s.prefix, handler, eventID.String())
}
```

---

## ✅ สรุปความครบถ้วน

| Layer | ไฟล์ | สถานะ |
|-------|------|-------|
| **domain/value_object** | 5 | ✅ |
| **domain/entity** | 5 | ✅ (มีในระบบเดิม) |
| **domain/event** | 5 | ✅ |
| **domain/repository** | 7 | ✅ |
| **domain/service** | 4 | ✅ |
| **domain/errors** | 1 | ✅ (มีแล้ว) |
| **application/dto** | 1 | ✅ |
| **application/command** | 10 | ✅ |
| **application/query** | 4 | ✅ |
| **application/event_handler** | 10+ | ✅ |
| **infrastructure/postgres** | 8 | ✅ |
| **infrastructure/cache/redis** | 2 | ✅ |
| **infrastructure/messaging/kafka** | 2 | ✅ |
| **infrastructure/messaging/outbox** | 1 | ✅ |
| **infrastructure/external** | 3 | ✅ |
| **interfaces/http** | 7 | ✅ |
| **module.go** | 1 | ✅ |
| **รวม** | **~75 ไฟล์** | **ครบ** |

### สิ่งที่ต้องเพิ่มเองในโปรเจกต์ของคุณ:

1. **`pkg/httputil`** — `RespondJSON`, `RespondError`, `DecodeJSON`
2. **`pkg/logger`** — `Infow`, `Errorw`, `Warnw`, `Debugf`, `Infof`, `Nop()`
3. **`pkg/kafka`** — แค่ re-export จาก `infrastructure/messaging/kafka`
4. **`config.Config`** — ปรับตาม module
5. **migration SQL** — `pdpa_consents`, `pdpa_dsar_requests`, `pdpa_user_account_statuses`, `pdpa_audit_trails`, `pdpa_privacy_policies`, `pdpa_outbox` (ผมเขียนไว้ในส่วนที่ 3 แล้ว)

 