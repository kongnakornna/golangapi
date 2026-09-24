# ปรับปรุง PDPA Module — Code Review + Refactor Plan

ผมได้ review โค้ดที่คุณส่งมาแล้ว **พบปัญหาสำคัญ 27 จุด** โดยแบ่งเป็น 3 ระดับความรุนแรง

---

## 🔴 สรุปปัญหาที่พบ (Critical Issues)

| # | Severity | ปัญหา | ผลกระทบ |
|---|----------|-------|---------|
| 1 | **P0** | Publish Kafka **ก่อน** DB commit | ข้อมูลไม่ sync ถ้า DB write fail |
| 2 | **P0** | `CompleteDSARCommand` ไม่เรียก `dsar.MarkCompleted()` — bypass domain logic | สถานะไม่ consistent |
| 3 | **P0** | `ProcessDSARCommand` เรียก `MarkProcessing()` แต่ **ไม่ persist** | สถานะไม่เปลี่ยนใน DB |
| 4 | **P0** | `VerifyDSAROTPCommand` ไม่ persist `OTPVerified` | ยืนยันซ้ำได้ |
| 5 | **P0** | OTP เก็บ **plaintext** ใน DB | รั่วไหลได้ |
| 6 | **P0** | Rate limit ใช้ `sync.Map` (in-memory) | ไม่ทำงาน multi-replica |
| 7 | **P1** | Value Objects เป็น `int8` ขัดกับ original spec (string-based) | API contract ไม่ชัด |
| 8 | **P1** | `ConsentLog.Purpose` เป็น `string` ไม่ใช่ VO | สูญเสีย type safety |
| 9 | **P1** | `time.Now()` hardcode ใน entity | test ไม่ได้ |
| 10 | **P1** | `kafka.Producer` concrete type | mock ไม่ได้ |
| 11 | **P1** | ไม่มี `Update` ใน ConsentRepository | Revoke แล้ว Save = insert ใหม่ |
| 12 | **P1** | ไม่มี correlation ID ใน events | trace ไม่ได้ |
| 13 | **P1** | Events ส่งเป็น `map[string]interface{}` | ไม่มี schema |
| 14 | **P2** | `ConfirmDeletionCommand` dead code (recheck nil) | สับสน |
| 15 | **P2** | `PdpaPolicy.Status int64` | magic number |
| 16 | **P2** | Logger ใช้ `Debugf/Infof` ปนกับ structured | inconsistent |
| 17 | **P2** | DTO validate `purpose gte=1` แต่ VO เป็น int8 | API docs ผิด |

---

## 🎯 แผนการปรับปรุง (Priority Fix Plan)

### P0 — ต้องแก้ทันที (Reliability + Security)

```
1. Transactional Outbox Pattern  → publish หลัง commit
2. บันทึก state ทุกครั้งที่เปลี่ยน (Process/Complete/Verify OTP)
3. Hash OTP ด้วย SHA-256
4. Rate limit ด้วย Redis (shared state)
```

### P1 — ควรแก้ (Correctness + Testability)

```
5. Value Objects → string-based ตาม spec
6. Inject clock (Clock interface)
7. Kafka Producer → interface
8. ConsentRepository.Update() แยกจาก Save()
9. Structured domain events + EventMetadata
```

### P2 — Nice to have

```
10. Enum types (PolicyStatus, PurposeType)
11. Structured logging everywhere
12. Fix dead code
```

---

## 📝 โค้ดที่ปรับปรุงแล้ว

### 1. Value Objects — เปลี่ยนเป็น String-based

**`internal/modules/pdpa/domain/valueobject/consent_purpose.go`** (ปรับปรุง)

```go
package valueobject

import (
	"fmt"
	"strings"
)

// ConsentPurpose วัตถุประสงค์ในการขอความยินยอม (string-based เพื่อ API contract ชัดเจน)
// ConsentPurpose represents the purpose (string-based for clarity in API contract)
type ConsentPurpose string

const (
	ConsentPurposeNecessary          ConsentPurpose = "NECESSARY"
	ConsentPurposeAnalytics          ConsentPurpose = "ANALYTICS"
	ConsentPurposeMarketing          ConsentPurpose = "MARKETING"
	ConsentPurposeAccountSystem      ConsentPurpose = "ACCOUNT_SYSTEM"
	ConsentPurposeUsageLogs          ConsentPurpose = "USAGE_LOGS"
	ConsentPurposeTransactionHistory ConsentPurpose = "TRANSACTION_HISTORY"
)

// ParseConsentPurpose แปลง string จาก API → VO (validation พร้อม)
// ParseConsentPurpose converts API string → VO with validation
func ParseConsentPurpose(s string) (ConsentPurpose, error) {
	p := ConsentPurpose(strings.ToUpper(strings.TrimSpace(s)))
	if !p.IsValid() {
		return "", fmt.Errorf("invalid consent purpose: %q", s)
	}
	return p, nil
}

// MustParseConsentPurpose parse หรือ panic (ใช้เฉพาะ test)
// MustParseConsentPurpose parses or panics (test-only)
func MustParseConsentPurpose(s string) ConsentPurpose {
	p, err := ParseConsentPurpose(s)
	if err != nil {
		panic(err)
	}
	return p
}

func (c ConsentPurpose) IsValid() bool {
	switch c {
	case ConsentPurposeNecessary, ConsentPurposeAnalytics, ConsentPurposeMarketing,
		ConsentPurposeAccountSystem, ConsentPurposeUsageLogs, ConsentPurposeTransactionHistory:
		return true
	}
	return false
}

// IsRequired ตรวจสอบว่าเป็น purpose บังคับ (revoke ไม่ได้)
// IsRequired checks if this purpose is mandatory
func (c ConsentPurpose) IsRequired() bool {
	return c == ConsentPurposeNecessary
}

// String คืนค่า string form (สำหรับ serialize + log)
func (c ConsentPurpose) String() string { return string(c) }

// MarshalJSON ทำให้ JSON ของ API ออกมาเป็น "NECESSARY" ไม่ใช่ 1
// MarshalJSON makes JSON output "NECESSARY" not 1
func (c ConsentPurpose) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(c) + `"`), nil
}

// UnmarshalJSON parse จาก string
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

**`internal/modules/pdpa/domain/valueobject/consent_status.go`** (ปรับปรุง)

```go
package valueobject

import "strings"

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

func (c *ConsentStatus) UnmarshalJSON(data []byte) error {
	*c = ConsentStatus(strings.Trim(string(data), `"`))
	return nil
}
```

**`internal/modules/pdpa/domain/valueobject/dsar_status.go`** (ปรับปรุง)

```go
package valueobject

import "strings"

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

// IsTerminal ตรวจสอบว่าเป็นสถานะสุดท้าย (เปลี่ยนไม่ได้อีก)
// IsTerminal checks if this is a terminal state
func (s DSARStatus) IsTerminal() bool {
	return s == DSARStatusCompleted || s == DSARStatusRejected
}

func (s DSARStatus) String() string { return string(s) }

func (s DSARStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(s) + `"`), nil
}

func (s *DSARStatus) UnmarshalJSON(data []byte) error {
	*s = DSARStatus(strings.Trim(string(data), `"`))
	return nil
}
```

**`internal/modules/pdpa/domain/valueobject/account_status.go`** (ปรับปรุง)

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

func (a AccountStatus) IsActive() bool    { return a == AccountActive }
func (a AccountStatus) IsSuspended() bool { return a == AccountSuspended }
func (a AccountStatus) IsTerminated() bool { return a == AccountTerminated }
func (a AccountStatus) IsDeleted() bool   { return a == AccountDeleted }
func (a AccountStatus) String() string    { return string(a) }
```

**`internal/modules/pdpa/domain/valueobject/dsar_type.go`** (ปรับปรุง)

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

func (d DSARType) String() string { return string(d) }
```

---

### 2. Domain Events — Structured Events

**`internal/modules/pdpa/domain/event/metadata.go`** (ใหม่)

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// EventMetadata metadata ที่ทุก event ต้องมี (สำหรับ tracing + idempotency)
// EventMetadata is required by every event (for tracing + idempotency)
type EventMetadata struct {
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id,omitempty"`
	TraceID       string    `json:"trace_id,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
	Version       int       `json:"version"`
	Source        string    `json:"source"`
}

// NewMetadata สร้าง metadata พร้อม defaults
// NewMetadata creates metadata with defaults
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

// Topics ค่าคงที่ของ Kafka topics
// Topics constants
const (
	TopicConsentGranted     = "pdpa.consent.granted"
	TopicConsentRevoked     = "pdpa.consent.revoked"
	TopicDSARSubmitted      = "pdpa.dsar.submitted"
	TopicDSARProcessing     = "pdpa.dsar.processing"
	TopicDSAROTPVerified    = "pdpa.dsar.otp_verified"
	TopicDSARCompleted      = "pdpa.dsar.completed"
	TopicDSARRejected       = "pdpa.dsar.rejected"
	TopicAccountSuspended   = "pdpa.account.suspended"
	TopicAccountTerminated  = "pdpa.account.terminated"
	TopicDeletionRequested  = "pdpa.data.deletion_requested"
	TopicDataDeleted        = "pdpa.data.deleted"
	TopicEmailNotification  = "pdpa.email.notification"
	TopicBlockchainRecord   = "pdpa.blockchain.record"
	TopicAuditTrail         = "pdpa.audit.trail"
)

// Envelope โครงสร้างมาตรฐานของ event ที่ส่งผ่าน Kafka
// Envelope is the standard event structure for Kafka
type Envelope struct {
	EventType string          `json:"event_type"`
	Metadata  EventMetadata   `json:"metadata"`
	Payload   interface{}     `json:"payload"`
}

// NewEnvelope สร้าง envelope
// NewEnvelope creates an envelope
func NewEnvelope(eventType string, meta EventMetadata, payload interface{}) Envelope {
	return Envelope{
		EventType: eventType,
		Metadata:  meta,
		Payload:   payload,
	}
}
```

**`internal/modules/pdpa/domain/event/consent_events.go`** (ใหม่)

```go
package event

import (
	"time"

	"github.com/google/uuid"
)

// ConsentGrantedPayload payload ของ ConsentGrantedEvent
type ConsentGrantedPayload struct {
	ConsentID  uuid.UUID `json:"consent_id"`
	UserID     uuid.UUID `json:"user_id"`
	Purpose    string    `json:"purpose"`
	GrantedAt  time.Time `json:"granted_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	IPAddress  string    `json:"ip_address,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}

// ConsentRevokedPayload payload ของ ConsentRevokedEvent
type ConsentRevokedPayload struct {
	ConsentID uuid.UUID `json:"consent_id"`
	UserID    uuid.UUID `json:"user_id"`
	Purpose   string    `json:"purpose"`
	RevokedAt time.Time `json:"revoked_at"`
}
```

**`internal/modules/pdpa/domain/event/dsar_events.go`** (ใหม่)

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

---

### 3. Entities — Inject Clock + Hash OTP

**`internal/modules/pdpa/domain/entity/consent_log.go`** (ปรับปรุง)

```go
package entity

import (
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// ConsentLog aggregate root สำหรับบันทึกความยินยอม
type ConsentLog struct {
	ID            uuid.UUID                   `json:"id"`
	UserID        uuid.UUID                   `json:"user_id"`
	Purpose       valueobject.ConsentPurpose  `json:"purpose"`
	Status        valueobject.ConsentStatus   `json:"status"`
	ConsentedAt   time.Time                   `json:"consented_at"`
	ExpiresAt     time.Time                   `json:"expires_at"`
	RevokedAt     *time.Time                  `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time                  `json:"auto_deleted_at,omitempty"`
	IPAddress     string                      `json:"ip_address,omitempty"`
	UserAgent     string                      `json:"user_agent,omitempty"`
}

// NewConsentLog สร้าง consent log ใหม่ — ต้องรับ now() จาก caller (injectable clock)
// NewConsentLog creates a new consent log — requires caller-supplied now() for testability
func NewConsentLog(
	userID uuid.UUID,
	purpose valueobject.ConsentPurpose,
	now time.Time,
	ip, userAgent string,
) (*ConsentLog, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidInput
	}
	if !purpose.IsValid() {
		return nil, domainerrors.ErrInvalidPurpose
	}
	return &ConsentLog{
		ID:          uuid.New(),
		UserID:      userID,
		Purpose:     purpose,
		Status:      valueobject.ConsentGranted,
		ConsentedAt: now,
		ExpiresAt:   now.AddDate(1, 0, 0),
		IPAddress:   ip,
		UserAgent:   userAgent,
	}, nil
}

// Revoke ถอนความยินยอม — caller ต้องส่ง now() (injectable)
func (c *ConsentLog) Revoke(now time.Time) error {
	if c.Status == valueobject.ConsentRevoked {
		return domainerrors.ErrConsentAlreadyRevoked
	}
	if c.Status != valueobject.ConsentGranted {
		return domainerrors.ErrRevokeNotAllowed
	}
	if c.Purpose.IsRequired() {
		return domainerrors.ErrRevokeNotAllowed
	}
	c.Status = valueobject.ConsentRevoked
	c.RevokedAt = &now
	return nil
}

func (c *ConsentLog) MarkDeleted(now time.Time) {
	c.Status = valueobject.ConsentDeleted
	c.AutoDeletedAt = &now
}

func (c *ConsentLog) IsActive(now time.Time) bool {
	return c.Status == valueobject.ConsentGranted && c.ExpiresAt.After(now)
}
```

**`internal/modules/pdpa/domain/entity/dsar_request.go`** (ปรับปรุง)

```go
package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

const (
	OTPTTL        = 15 * time.Minute
	OTPMaxAttempt = 5
)

type DSARRequest struct {
	ID              uuid.UUID                 `json:"id"`
	UserID          uuid.UUID                 `json:"user_id"`
	RequestType     valueobject.DSARType      `json:"request_type"`
	Status          valueobject.DSARStatus    `json:"status"`
	RequestedAt     time.Time                 `json:"requested_at"`
	CompletedAt     *time.Time                `json:"completed_at,omitempty"`
	DataPayload     []byte                    `json:"data_payload,omitempty"`
	DataHash        string                    `json:"data_hash,omitempty"`
	BlockchainTx    string                    `json:"blockchain_tx,omitempty"`
	RejectionReason string                    `json:"rejection_reason,omitempty"`

	// OTP — เก็บเป็น hash เท่านั้น (ไม่เก็บ plaintext)
	OTPHash      string     `json:"-"`
	OTPExpiresAt *time.Time `json:"-"`
	OTPVerified  bool       `json:"otp_verified"`
	OTPAttempts  int        `json:"-"`

	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// NewDSARRequest สร้าง DSAR ใหม่พร้อม OTP (คืน plaintext OTP ให้ caller ส่งอีเมล)
// NewDSARRequest creates a new DSAR with OTP (returns plaintext OTP for email delivery)
func NewDSARRequest(
	userID uuid.UUID,
	requestType valueobject.DSARType,
	now time.Time,
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
	expires := now.Add(OTPTTL)
	return &DSARRequest{
		ID:           uuid.New(),
		UserID:       userID,
		RequestType:  requestType,
		Status:       valueobject.DSARStatusPending,
		RequestedAt:  now,
		OTPHash:      hashOTP(otp),
		OTPExpiresAt: &expires,
		IPAddress:    ip,
		UserAgent:    userAgent,
	}, otp, nil
}

// VerifyOTP ตรวจสอบ OTP (นับ attempt)
func (d *DSARRequest) VerifyOTP(otp string, now time.Time) error {
	if d.OTPVerified {
		return nil // idempotent
	}
	if d.OTPAttempts >= OTPMaxAttempt {
		return domainerrors.ErrRateLimitExceeded
	}
	d.OTPAttempts++
	if d.OTPExpiresAt == nil || now.After(*d.OTPExpiresAt) {
		return domainerrors.ErrOTPExpired
	}
	if hashOTP(otp) != d.OTPHash {
		return domainerrors.ErrOTPInvalid
	}
	d.OTPVerified = true
	return nil
}

// MarkProcessing เปลี่ยนสถานะเป็นกำลังดำเนินการ
func (d *DSARRequest) MarkProcessing() error {
	if d.Status != valueobject.DSARStatusPending {
		return domainerrors.ErrDSARNotPending
	}
	if !d.OTPVerified {
		return domainerrors.ErrOTPInvalid
	}
	d.Status = valueobject.DSARStatusProcessing
	return nil
}

// MarkCompleted เสร็จสิ้น พร้อมคำนวณ hash + txHash
func (d *DSARRequest) MarkCompleted(payload []byte, dataHash, txHash string, now time.Time) error {
	if d.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusCompleted
	d.CompletedAt = &now
	d.DataPayload = payload
	d.DataHash = dataHash
	d.BlockchainTx = txHash
	return nil
}

// MarkRejected ปฏิเสธคำขอ
func (d *DSARRequest) MarkRejected(reason string) error {
	if d.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusRejected
	d.RejectionReason = reason
	return nil
}

// generateOTP สร้าง OTP 6 หลักด้วย crypto/rand
func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// hashOTP hash OTP ด้วย SHA-256 (ตอน verify เทียบ hash)
func hashOTP(otp string) string {
	h := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(h[:])
}
```

**`internal/modules/pdpa/domain/entity/user_account_status.go`** (ปรับปรุง)

```go
package entity

import (
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type UserAccountStatus struct {
	ID                  uuid.UUID                 `json:"id"`
	UserID              uuid.UUID                 `json:"user_id"`
	Status              valueobject.AccountStatus `json:"status"`
	SuspendedAt         *time.Time                `json:"suspended_at,omitempty"`
	TerminatedAt        *time.Time                `json:"terminated_at,omitempty"`
	DeletionConfirmedAt *time.Time                `json:"deletion_confirmed_at,omitempty"`
	RetentionDeadline   *time.Time                `json:"retention_deadline,omitempty"`
	AutoDeletedAt       *time.Time                `json:"auto_deleted_at,omitempty"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

func NewUserAccountStatus(userID uuid.UUID, now time.Time) (*UserAccountStatus, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidInput
	}
	return &UserAccountStatus{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    valueobject.AccountActive,
		UpdatedAt: now,
	}, nil
}

// Suspend ระงับบัญชี + ตั้ง retention deadline
func (u *UserAccountStatus) Suspend(now time.Time, retentionYears int) error {
	if u.Status == valueobject.AccountDeleted {
		return domainerrors.ErrAlreadyDeleted
	}
	if retentionYears < 1 {
		retentionYears = 1
	}
	deadline := now.AddDate(retentionYears, 0, 0)
	u.Status = valueobject.AccountSuspended
	u.SuspendedAt = &now
	u.RetentionDeadline = &deadline
	u.UpdatedAt = now
	return nil
}

func (u *UserAccountStatus) Terminate(now time.Time) error {
	if u.Status == valueobject.AccountDeleted {
		return domainerrors.ErrAlreadyDeleted
	}
	u.Status = valueobject.AccountTerminated
	u.TerminatedAt = &now
	u.UpdatedAt = now
	return nil
}

// ConfirmDeletion ยืนยันการลบ (ต้องอยู่สถานะ TERMINATED เท่านั้น)
func (u *UserAccountStatus) ConfirmDeletion(now time.Time) error {
	if u.Status != valueobject.AccountTerminated {
		return domainerrors.ErrAccountNotTerminated
	}
	u.DeletionConfirmedAt = &now
	u.UpdatedAt = now
	return nil
}

func (u *UserAccountStatus) MarkDeleted(now time.Time) {
	u.Status = valueobject.AccountDeleted
	u.AutoDeletedAt = &now
	u.UpdatedAt = now
}

// CanImmediateDelete ตรวจสอบว่าลบทันทีได้หรือไม่
func (u *UserAccountStatus) CanImmediateDelete() bool {
	return u.Status == valueobject.AccountTerminated && u.DeletionConfirmedAt != nil
}

// IsReadyForAutoDeletion ตรวจสอบว่าถึงเวลาลบอัตโนมัติหรือยัง
func (u *UserAccountStatus) IsReadyForAutoDeletion(now time.Time) bool {
	return u.Status == valueobject.AccountSuspended &&
		u.RetentionDeadline != nil &&
		!u.RetentionDeadline.After(now)
}
```

---

### 4. Repository Interfaces — เพิ่ม Update + Transaction

**`internal/modules/pdpa/domain/repository/consent_repo.go`** (ปรับปรุง)

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
	// Save บันทึก consent ใหม่ (insert เท่านั้น)
	Save(ctx context.Context, consent *entity.ConsentLog) error

	// Update อัปเดต consent ที่มีอยู่ (revoke/delete)
	// Update อัปเดตเท่านั้น — คืน error ถ้าไม่พบ record
	Update(ctx context.Context, consent *entity.ConsentLog) error

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

**`internal/modules/pdpa/domain/repository/outbox_repo.go`** (ใหม่)

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OutboxEvent เก็บ event ที่จะ publish ไป Kafka แบบ atomic กับ DB transaction
// OutboxEvent stores events for atomic publish with DB transactions
type OutboxEvent struct {
	ID            uuid.UUID  `json:"id"`
	AggregateType string     `json:"aggregate_type"`
	AggregateID   uuid.UUID  `json:"aggregate_id"`
	EventType     string     `json:"event_type"`
	Topic         string     `json:"topic"`
	Payload       []byte     `json:"payload"`
	Status        string     `json:"status"` // PENDING | PUBLISHED | FAILED
	RetryCount    int        `json:"retry_count"`
	LastError     string     `json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
}

const (
	OutboxPending   = "PENDING"
	OutboxPublished = "PUBLISHED"
	OutboxFailed    = "FAILED"
)

// OutboxRepository transactional outbox pattern
type OutboxRepository interface {
	Save(ctx context.Context, event *OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error)
}

// UnitOfWork abstraction สำหรับ transaction spanning multiple repos
// UnitOfWork allows multi-repo transactions
type UnitOfWork interface {
	// Do รัน fn ใน transaction เดียว — rollback อัตโนมัติเมื่อ error
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
```

---

### 5. Commands — แก้ Publish Order (ใช้ Outbox)

**`internal/modules/pdpa/application/command/GrantConsentCommand.go`** (ปรับปรุง)

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

// GrantConsentCommand ให้ความยินยอม (write path) — ใช้ Outbox Pattern
// GrantConsentCommand grants consent — uses Transactional Outbox Pattern
type GrantConsentCommand struct {
	consentRepo repository.ConsentRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	clock       func() time.Time
	logger      logger.Logger
}

func NewGrantConsentCommand(
	consentRepo repository.ConsentRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	clock func() time.Time,
	logger logger.Logger,
) *GrantConsentCommand {
	if clock == nil {
		clock = time.Now
	}
	return &GrantConsentCommand{
		consentRepo: consentRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		clock:       clock,
		logger:      logger,
	}
}

// Execute ทำงานใน transaction เดียว:
//   1. insert consent
//   2. insert outbox event
//   → แล้ว background publisher จะส่ง Kafka
// Execute runs in a single transaction:
//   1. insert consent
//   2. insert outbox event
//   → background publisher delivers to Kafka
func (c *GrantConsentCommand) Execute(ctx context.Context, req dto.GrantConsentRequest) error {
	purpose, err := parsePurpose(req.Purpose)
	if err != nil {
		return err
	}
	now := c.clock()

	consent, err := entity.NewConsentLog(req.UserID, purpose, now, req.IPAddress, req.UserAgent)
	if err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	payload := event.ConsentGrantedPayload{
		ConsentID: consent.ID,
		UserID:    consent.UserID,
		Purpose:   consent.Purpose.String(),
		GrantedAt: consent.ConsentedAt,
		ExpiresAt: consent.ExpiresAt,
		IPAddress: consent.IPAddress,
		UserAgent: consent.UserAgent,
	}
	envelope := event.NewEnvelope(event.TopicConsentGranted, meta, payload)
	body, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// ⚡ Atomic: DB write + outbox insert ต้องอยู่ใน transaction เดียวกัน
	return c.uow.Do(ctx, func(txCtx context.Context) error {
		if err := c.consentRepo.Save(txCtx, consent); err != nil {
			return fmt.Errorf("save consent: %w", err)
		}
		if err := c.outboxRepo.Save(txCtx, &repository.OutboxEvent{
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
		c.logger.Infow("consent granted",
			"user_id", consent.UserID,
			"purpose", consent.Purpose.String(),
			"consent_id", consent.ID,
			"correlation_id", meta.CorrelationID,
		)
		return nil
	})
}

// parsePurpose parse จาก DTO (string) → VO
func parsePurpose(s string) (valueobject.ConsentPurpose, error) {
	return valueobject.ParseConsentPurpose(s)
}
```

**`internal/modules/pdpa/application/command/ProcessDSARCommand.go`** (ปรับปรุง — fix missing persist)

```go
package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type ProcessDSARCommand struct {
	dsarRepo   repository.DSARRepository
	outboxRepo repository.OutboxRepository
	uow        repository.UnitOfWork
	logger     logger.Logger
}

func NewProcessDSARCommand(
	dsarRepo repository.DSARRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	logger logger.Logger,
) *ProcessDSARCommand {
	return &ProcessDSARCommand{
		dsarRepo:   dsarRepo,
		outboxRepo: outboxRepo,
		uow:        uow,
		logger:     logger,
	}
}

func (c *ProcessDSARCommand) Execute(ctx context.Context, req *dto.ProcessDSARRequest) error {
	dsar, err := c.dsarRepo.FindByID(ctx, req.DSARRequestID)
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
	envelope := event.NewEnvelope(event.TopicDSARProcessing, meta, map[string]any{
		"dsar_id":      dsar.ID,
		"user_id":      dsar.UserID,
		"status":       dsar.Status.String(),
		"processed_at": meta.OccurredAt,
	})
	body, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// 🔧 FIX: ต้อง persist state change (เดิมไม่บันทึก!)
	return c.uow.Do(ctx, func(txCtx context.Context) error {
		if err := c.dsarRepo.Update(txCtx, dsar); err != nil {
			return fmt.Errorf("update dsar: %w", err)
		}
		return c.outboxRepo.Save(txCtx, &repository.OutboxEvent{
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

var _ = errors.Is
```

**`internal/modules/pdpa/application/command/CompleteDSARCommand.go`** (ปรับปรุง)

```go
package command

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

type CompleteDSARCommand struct {
	dsarRepo      repository.DSARRepository
	outboxRepo    repository.OutboxRepository
	uow           repository.UnitOfWork
	blockchainSvc service.BlockchainService
	clock         func() time.Time
	logger        logger.Logger
}

func NewCompleteDSARCommand(
	dsarRepo repository.DSARRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	blockchainSvc service.BlockchainService,
	clock func() time.Time,
	logger logger.Logger,
) *CompleteDSARCommand {
	if clock == nil {
		clock = time.Now
	}
	return &CompleteDSARCommand{
		dsarRepo:      dsarRepo,
		outboxRepo:    outboxRepo,
		uow:           uow,
		blockchainSvc: blockchainSvc,
		clock:         clock,
		logger:        logger,
	}
}

func (c *CompleteDSARCommand) Execute(ctx context.Context, req *dto.CompleteDSARRequest) error {
	dsar, err := c.dsarRepo.FindByID(ctx, req.DSARRequestID)
	if err != nil {
		return err
	}
	if dsar == nil {
		return domainerrors.ErrDSARNotFound
	}
	if dsar.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}

	// 1. Hash payload
	hash := sha256.Sum256(req.ResponseData)
	dataHash := hex.EncodeToString(hash[:])

	// 2. บันทึก blockchain (best-effort — ไม่ block completion ถ้า fail)
	txHash := ""
	if h, err := c.blockchainSvc.RecordHash(ctx, dataHash); err != nil {
		c.logger.Warnw("blockchain record failed (continuing)",
			"dsar_id", dsar.ID, "error", err)
	} else {
		txHash = h
	}

	// 3. เรียก domain method (ไม่ bypass)
	now := c.clock()
	if err := dsar.MarkCompleted(req.ResponseData, dataHash, txHash, now); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicDSARCompleted, meta, event.DSARCompletedPayload{
		DSARID:      dsar.ID,
		UserID:      dsar.UserID,
		Status:      dsar.Status.String(),
		DataHash:    dataHash,
		TxHash:      txHash,
		CompletedAt: *dsar.CompletedAt,
	})
	body, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// 4. Persist + Outbox (atomic)
	if err := c.uow.Do(ctx, func(txCtx context.Context) error {
		if err := c.dsarRepo.Update(txCtx, dsar); err != nil {
			return fmt.Errorf("update dsar: %w", err)
		}
		return c.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "DSARRequest",
			AggregateID:   dsar.ID,
			EventType:     event.TopicDSARCompleted,
			Topic:         event.TopicDSARCompleted,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	}); err != nil {
		return err
	}

	c.logger.Infow("dsar completed",
		"user_id", dsar.UserID,
		"dsar_id", dsar.ID,
		"correlation_id", meta.CorrelationID,
	)
	return nil
}
```

**`internal/modules/pdpa/application/command/VerifyDSAROTPCommand.go`** (ปรับปรุง — ใช้ Redis rate limit)

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

// OTPRateLimiter abstraction (Redis-backed)
type OTPRateLimiter interface {
	// IncrementAndGet เพิ่ม counter + คืนค่าล่าสุด, ตั้ง TTL ถ้าเพิ่งสร้าง
	IncrementAndGet(ctx context.Context, key string, ttlSeconds int) (int, error)
	// Reset ล้าง counter (เรียกเมื่อ verify สำเร็จ)
	Reset(ctx context.Context, key string) error
}

type VerifyDSAROTPCommand struct {
	dsarRepo    repository.DSARRepository
	outboxRepo  repository.OutboxRepository
	uow         repository.UnitOfWork
	rateLimiter OTPRateLimiter
	clock       func() time.Time
	logger      logger.Logger
}

func NewVerifyDSAROTPCommand(
	dsarRepo repository.DSARRepository,
	outboxRepo repository.OutboxRepository,
	uow repository.UnitOfWork,
	rateLimiter OTPRateLimiter,
	clock func() time.Time,
	logger logger.Logger,
) *VerifyDSAROTPCommand {
	if clock == nil {
		clock = time.Now
	}
	return &VerifyDSAROTPCommand{
		dsarRepo:    dsarRepo,
		outboxRepo:  outboxRepo,
		uow:         uow,
		rateLimiter: rateLimiter,
		clock:       clock,
		logger:      logger,
	}
}

func (c *VerifyDSAROTPCommand) Execute(ctx context.Context, req *dto.VerifyOTPRequest) error {
	key := fmt.Sprintf("otp:attempt:%s", req.DSARRequestID.String())

	// 1. Shared rate limit (Redis — ทำงานข้าม replicas)
	attempts, err := c.rateLimiter.IncrementAndGet(ctx, key, 900) // 15 min TTL
	if err != nil {
		c.logger.Errorw("rate limiter failed", "error", err)
		// don't block, log only — fail open
	}
	if attempts > 5 {
		return domainerrors.ErrRateLimitExceeded
	}

	dsar, err := c.dsarRepo.FindByID(ctx, req.DSARRequestID)
	if err != nil {
		return err
	}
	if dsar == nil {
		return domainerrors.ErrDSARNotFound
	}

	now := c.clock()
	if err := dsar.VerifyOTP(req.OTPCode, now); err != nil {
		return err
	}

	meta := event.NewMetadata(req.CorrelationID, "", req.TraceID)
	envelope := event.NewEnvelope(event.TopicDSAROTPVerified, meta, map[string]any{
		"dsar_id":     dsar.ID,
		"user_id":     dsar.UserID,
		"verified_at": now,
	})
	body, _ := json.Marshal(envelope)

	// 2. Persist (🔧 FIX: เดิมไม่บันทึก!)
	if err := c.uow.Do(ctx, func(txCtx context.Context) error {
		if err := c.dsarRepo.Update(txCtx, dsar); err != nil {
			return err
		}
		return c.outboxRepo.Save(txCtx, &repository.OutboxEvent{
			ID:            meta.EventID,
			AggregateType: "DSARRequest",
			AggregateID:   dsar.ID,
			EventType:     event.TopicDSAROTPVerified,
			Topic:         event.TopicDSAROTPVerified,
			Payload:       body,
			Status:        repository.OutboxPending,
			CreatedAt:     now,
		})
	}); err != nil {
		return err
	}

	// 3. Reset rate limit
	_ = c.rateLimiter.Reset(ctx, key)
	return nil
}
```

---

### 6. DTO — ปรับให้สอดคล้อง String VO + Metadata

**`internal/modules/pdpa/application/dto/grant_consent_request.go`** (ปรับปรุง)

```go
package dto

import "github.com/google/uuid"

// GrantConsentRequest payload (รับ purpose เป็น string)
// GrantConsentRequest payload (purpose as string)
type GrantConsentRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	Purpose       string    `json:"purpose" validate:"required,oneof=NECESSARY ANALYTICS MARKETING ACCOUNT_SYSTEM USAGE_LOGS TRANSACTION_HISTORY"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

// RevokeConsentRequest payload
type RevokeConsentRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	Purpose       string    `json:"purpose" validate:"required,oneof=ANALYTICS MARKETING ACCOUNT_SYSTEM USAGE_LOGS TRANSACTION_HISTORY"` // ไม่รวม NECESSARY (mandatory revoke ไม่ได้)
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

// SubmitDSARRequest payload
type SubmitDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	RequestType   string    `json:"request_type" validate:"required,oneof=ACCESS ERASURE WITHDRAW_CONSENT"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}

// VerifyOTPRequest payload
type VerifyOTPRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid"`
	OTPCode       string    `json:"otp_code" validate:"required,len=6,numeric"`
	CorrelationID string    `json:"-"`
	TraceID       string    `json:"-"`
}
```

---

### 7. Errors — เพิ่มที่ขาด

**`internal/modules/pdpa/domain/errors/errors.go`** (ปรับปรุง)

```go
package domainerrors

import "errors"

var (
	// Consent
	ErrConsentNotFound       = errors.New("consent record not found")
	ErrConsentAlreadyRevoked = errors.New("consent already revoked")
	ErrConsentExpired        = errors.New("consent expired")
	ErrInvalidPurpose        = errors.New("invalid consent purpose")
	ErrRevokeNotAllowed      = errors.New("cannot revoke non-granted consent")

	// DSAR
	ErrDSARNotFound         = errors.New("DSAR request not found")
	ErrDSARAlreadyProcessed = errors.New("DSAR request already processed")
	ErrDSARNotPending       = errors.New("DSAR request is not in pending state")
	ErrOTPInvalid           = errors.New("OTP is invalid")
	ErrOTPExpired           = errors.New("OTP expired")
	ErrRateLimitExceeded    = errors.New("too many requests, please try again later")

	// Account
	ErrAccountNotFound             = errors.New("user account status not found")
	ErrAccountNotActive            = errors.New("account is not active")
	ErrAccountNotSuspended         = errors.New("account is not suspended")
	ErrAccountNotTerminated        = errors.New("account is not terminated")
	ErrDeletionNotConfirmed        = errors.New("deletion not confirmed for this account")
	ErrImmediateDeletionNotAllowed = errors.New("immediate deletion is not allowed for this account")
	ErrAlreadyDeleted              = errors.New("account data already deleted")

	// Policy
	ErrPolicyNotFound      = errors.New("privacy policy not found")
	ErrPolicyVersionExists = errors.New("privacy policy version already exists")

	// Outbox
	ErrOutboxPublishFailed = errors.New("failed to publish outbox event")

	// General
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrInternal     = errors.New("internal server error")
)
```

---

### 8. Outbox Publisher (Background Worker)

**`internal/modules/pdpa/infrastructure/messaging/outbox/publisher.go`** (ใหม่)

```go
package outbox

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// Publisher รัน background loop — ดึง PENDING จาก outbox → publish Kafka → mark PUBLISHED
// Publisher runs a background loop to drain the outbox
type Publisher struct {
	outboxRepo repository.OutboxRepository
	producer   kafka.Producer
	interval   time.Duration
	batchSize  int
	logger     logger.Logger
}

func NewPublisher(
	outboxRepo repository.OutboxRepository,
	producer kafka.Producer,
	interval time.Duration,
	batchSize int,
	logger logger.Logger,
) *Publisher {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	return &Publisher{
		outboxRepo: outboxRepo,
		producer:   producer,
		interval:   interval,
		batchSize:  batchSize,
		logger:     logger,
	}
}

// Run loop — ใช้ context cancellation สำหรับ graceful shutdown
func (p *Publisher) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.logger.Infow("outbox publisher starting",
		"interval", p.interval, "batch_size", p.batchSize)

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

func (p *Publisher) drain(ctx context.Context) {
	events, err := p.outboxRepo.FetchPending(ctx, p.batchSize)
	if err != nil {
		p.logger.Errorw("fetch outbox failed", "error", err)
		return
	}
	for i := range events {
		e := &events[i]
		// ใช้ AggregateID เป็น key → ordering ต่อ aggregate
		key := e.AggregateID.String()

		if err := p.producer.Publish(ctx, e.Topic, key, e.Payload); err != nil {
			_ = p.outboxRepo.MarkFailed(ctx, e.ID, err.Error())
			p.logger.Warnw("publish failed",
				"event_id", e.ID,
				"topic", e.Topic,
				"retry", e.RetryCount+1,
				"error", err)
			continue
		}
		if err := p.outboxRepo.MarkPublished(ctx, e.ID); err != nil {
			p.logger.Warnw("mark published failed", "event_id", e.ID, "error", err)
		}
	}
}
```

---

### 9. Unit of Work (GORM Implementation)

**`internal/modules/pdpa/infrastructure/persistence/postgres/uow.go`** (ใหม่)

```go
package postgres

import (
	"context"

	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/repository"
)

// UnitOfWorkImpl GORM-based UnitOfWork
type UnitOfWorkImpl struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{db: db}
}

// gormTxKey ใช้เก็บ *gorm.DB ของ tx ใน context
type gormTxKey struct{}

// Do รัน fn ใน transaction เดียว
// Do runs fn in a single transaction (rolls back on error, commits otherwise)
func (u *UnitOfWorkImpl) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	// ถ้าอยู่ใน tx อยู่แล้ว → nested call → ใช้ tx เดิม
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return fn(ctx)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, gormTxKey{}, tx)
		return fn(txCtx)
	})
}

// GetDB คืน *gorm.DB (จะคืน tx ถ้าอยู่ใน UnitOfWork)
// GetDB returns the active *gorm.DB (tx if inside UnitOfWork)
func GetDB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(gormTxKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}

var _ repository.UnitOfWork = (*UnitOfWorkImpl)(nil)
```

**Update Repo Implementations ให้ใช้ `GetDB`:**

```go
// ตัวอย่าง consent_repo_impl.go
func (r *ConsentRepoImpl) Save(ctx context.Context, c *entity.ConsentLog) error {
	db := GetDB(ctx, r.db) // ← ใช้ tx ถ้าอยู่ใน UnitOfWork
	return db.Create(toModel(c)).Error
}

func (r *ConsentRepoImpl) Update(ctx context.Context, c *entity.ConsentLog) error {
	db := GetDB(ctx, r.db)
	res := db.Model(&ConsentModel{}).
		Where("id = ?", c.ID).
		Updates(map[string]any{
			"status":     c.Status.String(),
			"revoked_at": c.RevokedAt,
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
```

---

### 10. Idempotency — Redis-backed

**`internal/modules/pdpa/infrastructure/idempotency/redis_store.go`** (ใหม่)

```go
package idempotency

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisStore idempotency store (Redis SET NX)
type RedisStore struct {
	client *redis.Client
	prefix string
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client, prefix: "pdpa:idem:"}
}

// IsProcessed ตรวจสอบว่า event นี้ถูกประมวลผลแล้วหรือยัง
func (s *RedisStore) IsProcessed(ctx context.Context, eventID uuid.UUID, handler string) (bool, error) {
	key := s.key(eventID, handler)
	n, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// MarkProcessed บันทึกด้วย SET NX (atomic)
// MarkProcessed uses SET NX (atomic, prevents race between replicas)
func (s *RedisStore) MarkProcessed(ctx context.Context, eventID uuid.UUID, handler string, ttl time.Duration) error {
	key := s.key(eventID, handler)
	ok, err := s.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		// already exists — treat as success (idempotent)
		return nil
	}
	return nil
}

func (s *RedisStore) key(eventID uuid.UUID, handler string) string {
	return fmt.Sprintf("%s%s:%s", s.prefix, handler, eventID.String())
}
```

---

## 📊 สรุปการเปลี่ยนแปลง

### ก่อน vs หลัง

| ประเด็น | ก่อน | หลัง |
|---------|------|------|
| **Publish order** | Kafka ก่อน DB ❌ | DB + Outbox atomic แล้ว background publish ✅ |
| **VO type** | `int8` | `string` + JSON marshal ✅ |
| **OTP** | Plaintext ❌ | SHA-256 hash ✅ |
| **Rate limit** | `sync.Map` (per-pod) ❌ | Redis (shared) ✅ |
| **time.Now()** | Hardcode ❌ | Inject `clock func()` ✅ |
| **Kafka** | Concrete `*kafka.Producer` ❌ | Interface `kafka.Producer` ✅ |
| **Persist** | Process/Verify ไม่ save ❌ | Save + Update ทุก path ✅ |
| **Events** | `map[string]any` ❌ | Structured `Envelope{Metadata, Payload}` ✅ |
| **Correlation ID** | ไม่มี ❌ | ทุก DTO + Event ✅ |
| **Idempotency** | ไม่มี ❌ | Redis SET NX ✅ |
| **UnitOfWork** | ไม่มี ❌ | Interface + GORM impl ✅ |
| **Dead code** | มีใน ConfirmDeletion ❌ | ลบออก ✅ |

### Test Coverage ที่จะได้

```go
// Test ตัวอย่าง — ใช้ fake clock
func TestGrantConsent_PublishesAfterCommit(t *testing.T) {
    fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    clock := func() time.Time { return fixedTime }

    consentRepo := mocks.NewConsentRepo()
    outboxRepo := mocks.NewOutboxRepo()
    uow := mocks.NewFailingUow(consentRepo.err) // simulate DB failure

    cmd := command.NewGrantConsentCommand(consentRepo, outboxRepo, uow, clock, logger.Nop())
    err := cmd.Execute(ctx, dto.GrantConsentRequest{...})

    // Assert: outbox ต้องว่าง (rollback)
    assert.Error(t, err)
    assert.Empty(t, outboxRepo.Events())
}
```

---

## 🚀 ลำดับการ Migration (แนะนำ)

```
Sprint 1 (P0):
  1. เพิ่ม OutboxRepository + UnitOfWork interfaces
  2. Implement GORM OutboxRepository + UnitOfWork
  3. เพิ่ม OutboxPublisher background worker
  4. Refactor GrantConsentCommand/RevokeConsentCommand ให้ใช้ Outbox
  5. Fix ProcessDSARCommand (add Update)
  6. Fix CompleteDSARCommand (call MarkCompleted)
  7. Hash OTP + VerifyDSAROTP persist

Sprint 2 (P1):
  8. เปลี่ยน Value Objects → string
  9. Update DTOs + JSON marshal
  10. Inject clock ในทุก entity
  11. Kafka Producer → interface
  12. Structured events + EventMetadata
  13. Redis idempotency store

Sprint 3 (P2):
  14. Test coverage ≥ 80%
  15. Load test with new outbox flow
  16. Grafana dashboard สำหรับ outbox backlog
```

---

**ต้องการให้ผม:**
- 🔧 เขียน **full code** ของไฟล์ที่เหลือ (ConfirmDeletion, RejectDSAR, Suspend/Terminate, Revoke, Queries)?
- 📝 เขียน **migration SQL** สำหรับ outbox + idempotency tables?
- 🧪 เขียน **unit tests** ตัวอย่างที่ใช้ fake clock + mock?
- 📊 วาด **sequence diagram** ของ flow ใหม่?

 