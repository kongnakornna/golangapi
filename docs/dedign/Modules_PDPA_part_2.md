# ระบบ PDPA Module - ส่วนที่ 2

## 3.2 Repository Interfaces

### 3.2.1 ConsentRepository

**`internal/modules/pdpa/domain/repository/consent_repository.go`**

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentRepository interface สำหรับจัดการ ConsentLog
// ConsentRepository is the persistence interface for ConsentLog
type ConsentRepository interface {
	// Save บันทึก consent log ใหม่ / Save persists a new consent log
	Save(ctx context.Context, log *entity.ConsentLog) error

	// Update อัปเดต consent log ที่มีอยู่ / Update updates an existing consent log
	Update(ctx context.Context, log *entity.ConsentLog) error

	// FindByID ค้นหา consent ตาม ID / FindByID finds consent by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ConsentLog, error)

	// FindByUserID ค้นหา consent ทั้งหมดของ user / FindByUserID finds all consents for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error)

	// FindActiveByUserAndPurpose ค้นหา consent ที่ยังมีผล / Find active consent by user & purpose
	FindActiveByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (*entity.ConsentLog, error)

	// FindLatestByUserAndPurpose ค้นหา consent ล่าสุด (ไม่จำกัด status) / Find latest regardless of status
	FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (*entity.ConsentLog, error)

	// FindReadyForAutoDeletion ค้นหา consent ที่พร้อมลบอัตโนมัติ / Find consents ready for auto-deletion
	FindReadyForAutoDeletion(ctx context.Context, cutoff time.Time, limit int) ([]entity.ConsentLog, error)

	// DeleteByUserID ลบ consent ทั้งหมดของ user / Delete all consents for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteByID ลบ consent ตาม ID / Delete consent by ID
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// IsConsentActive ตรวจสอบว่า consent ยังมีผล / Check if consent is currently active
	IsConsentActive(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (bool, error)

	// CountActiveByPurpose นับจำนวน consent ที่ active แยกตาม purpose / Count active consents grouped by purpose
	CountActiveByPurpose(ctx context.Context, from, to time.Time) (map[string]int64, error)
}
```

### 3.2.2 DSARRepository

**`internal/modules/pdpa/domain/repository/dsar_repository.go`**

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// DSARRepository interface สำหรับจัดการ DSARRequest
// DSARRepository is the persistence interface for DSARRequest
type DSARRepository interface {
	// Save บันทึก DSAR ใหม่ / Save persists a new DSAR
	Save(ctx context.Context, req *entity.DSARRequest) error

	// Update อัปเดต DSAR ที่มีอยู่ / Update updates an existing DSAR
	Update(ctx context.Context, req *entity.DSARRequest) error

	// FindByID ค้นหา DSAR ตาม ID / FindByID finds DSAR by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error)

	// FindByUserID ค้นหา DSAR ทั้งหมดของ user / Find all DSARs for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error)

	// FindByStatus ค้นหา DSAR ตาม status / Find DSARs by status
	FindByStatus(ctx context.Context, status vo.DSARStatus, limit, offset int) ([]entity.DSARRequest, error)

	// CountByUserAndTypeInPeriod นับจำนวนคำร้องของ user ในช่วงเวลา (rate limit)
	// Count user requests in a period (for rate limiting)
	CountByUserAndTypeInPeriod(ctx context.Context, userID uuid.UUID, requestType vo.DSARType, from time.Time) (int64, error)

	// DeleteByUserID ลบ DSAR ทั้งหมดของ user / Delete all DSARs for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

### 3.2.3 UserAccountStatusRepository

**`internal/modules/pdpa/domain/repository/account_status_repository.go`**

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
)

// UserAccountStatusRepository interface สำหรับจัดการสถานะบัญชี
// UserAccountStatusRepository is the persistence interface for account status
type UserAccountStatusRepository interface {
	// Save บันทึก (upsert) สถานะบัญชี / Save (upsert) account status
	Save(ctx context.Context, status *entity.UserAccountStatus) error

	// FindByUserID ค้นหาสถานะบัญชีของ user / Find account status by user ID
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error)

	// FindReadyForAutoDeletion ค้นหาบัญชีที่พร้อมลบอัตโนมัติ (SUSPENDED + past deadline)
	// Find accounts ready for auto-deletion
	FindReadyForAutoDeletion(ctx context.Context, now time.Time, limit int) ([]entity.UserAccountStatus, error)

	// ExistsByUserID ตรวจสอบว่ามีสถานะบัญชีหรือยัง / Check if account status exists
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}
```

### 3.2.4 AuditRepository

**`internal/modules/pdpa/domain/repository/audit_repository.go`**

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
)

// AuditFilter ตัวกรองสำหรับค้นหา audit logs
// AuditFilter holds filter criteria for audit queries
type AuditFilter struct {
	UserID    *uuid.UUID
	Action    string
	FromTime  *time.Time
	ToTime    *time.Time
	Limit     int
	Offset    int
}

// AuditRepository interface สำหรับจัดการ audit trail
// AuditRepository is the persistence interface for audit trails
type AuditRepository interface {
	// Save บันทึก audit log / Save an audit log
	Save(ctx context.Context, audit *entity.AuditTrail) error

	// FindByFilter ค้นหา audit logs ตามเงื่อนไข / Find audit logs with filter
	FindByFilter(ctx context.Context, filter AuditFilter) ([]entity.AuditTrail, error)

	// CountByAction นับจำนวน audit logs แยกตาม action ในช่วงเวลา
	// Count audit logs by action in a time range
	CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error)
}
```

### 3.2.5 PolicyRepository

**`internal/modules/pdpa/domain/repository/policy_repository.go`**

```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
)

// PolicyRepository interface สำหรับจัดการ privacy policy
// PolicyRepository is the persistence interface for privacy policies
type PolicyRepository interface {
	// Save บันทึกนโยบายใหม่ / Save a new policy
	Save(ctx context.Context, policy *entity.PrivacyPolicy) error

	// Update อัปเดตนโยบาย / Update an existing policy
	Update(ctx context.Context, policy *entity.PrivacyPolicy) error

	// FindByID ค้นหานโยบายตาม ID / Find policy by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PrivacyPolicy, error)

	// FindActive ค้นหานโยบายที่เปิดใช้งาน / Find the currently active policy
	FindActive(ctx context.Context) (*entity.PrivacyPolicy, error)

	// FindAll ค้นหานโยบายทั้งหมด เรียงจากใหม่ไปเก่า / Find all policies, newest first
	FindAll(ctx context.Context) ([]entity.PrivacyPolicy, error)

	// FindByVersion ค้นหานโยบายตามเวอร์ชัน / Find policy by version
	FindByVersion(ctx context.Context, version string) (*entity.PrivacyPolicy, error)

	// DeactivateAll ปิดการใช้งานนโยบายทั้งหมด / Deactivate all policies
	DeactivateAll(ctx context.Context) error
}
```

### 3.2.6 OutboxRepository (Transactional Outbox Pattern)

**`internal/modules/pdpa/domain/repository/outbox_repository.go`**

```go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OutboxStatus สถานะของ outbox event
// OutboxStatus represents the status of an outbox event
type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "PENDING"
	OutboxStatusPublished OutboxStatus = "PUBLISHED"
	OutboxStatusFailed    OutboxStatus = "FAILED"
)

// OutboxEvent เหตุการณ์ที่รอ publish ไปยัง Kafka
// OutboxEvent is a domain event waiting to be published to Kafka
type OutboxEvent struct {
	ID            uuid.UUID    `json:"id"`
	AggregateType string       `json:"aggregate_type"`
	AggregateID   uuid.UUID    `json:"aggregate_id"`
	EventType     string       `json:"event_type"`
	Topic         string       `json:"topic"`
	Payload       []byte       `json:"payload"`
	Metadata      []byte       `json:"metadata"`
	Status        OutboxStatus `json:"status"`
	RetryCount    int          `json:"retry_count"`
	LastError     string       `json:"last_error,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	PublishedAt   *time.Time   `json:"published_at,omitempty"`
}

// OutboxRepository interface สำหรับ transactional outbox pattern
// OutboxRepository is the persistence interface for the transactional outbox
type OutboxRepository interface {
	// Save บันทึก outbox event (ต้องอยู่ใน transaction เดียวกับ aggregate)
	// Save persists an outbox event (must be within the same transaction as the aggregate)
	Save(ctx context.Context, event *OutboxEvent) error

	// FetchPending ดึง events ที่ยังไม่ publish / Fetch events that haven't been published
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)

	// MarkPublished ทำเครื่องหมายว่า publish สำเร็จ / Mark as published
	MarkPublished(ctx context.Context, id uuid.UUID) error

	// MarkFailed ทำเครื่องหมายว่า publish ล้มเหลว / Mark as failed
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error

	// CleanupPublished ลบ events ที่ publish แล้วและเก่ากว่า cutoff / Cleanup old published events
	CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error)
}
```

---

## 3.3 Domain Services

### 3.3.1 ConsentValidator

**`internal/modules/pdpa/domain/service/consent_validator.go`**

```go
package service

import (
	"fmt"

	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentValidator Domain Service สำหรับ validate การขอความยินยอม
// ConsentValidator is a Domain Service that validates consent requests
type ConsentValidator struct{}

// NewConsentValidator สร้าง ConsentValidator ใหม่
// NewConsentValidator creates a new ConsentValidator
func NewConsentValidator() *ConsentValidator {
	return &ConsentValidator{}
}

// ValidatePurposes ตรวจสอบรายการ purposes ทั้งหมด
// ValidatePurposes validates a collection of purposes
// - ต้องมี purpose ที่บังคับ (NECESSARY) เสมอ
// - Purpose ที่ไม่รู้จักจะถูก reject
// Must include mandatory purposes (NECESSARY); unknown purposes are rejected
func (v *ConsentValidator) ValidatePurposes(purposes map[string]bool) ([]vo.ConsentPurpose, error) {
	if len(purposes) == 0 {
		return nil, fmt.Errorf("no purposes provided")
	}

	result := make([]vo.ConsentPurpose, 0, len(purposes))
	hasNecessary := false

	for raw, granted := range purposes {
		p, err := vo.NewConsentPurpose(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid purpose %q: %w", raw, err)
		}
		if p.IsRequired() {
			hasNecessary = true
			// Purpose ที่บังคับต้อง granted เสมอ
			// Mandatory purpose must always be granted
			result = append(result, p)
			continue
		}
		if granted {
			result = append(result, p)
		}
	}

	if !hasNecessary {
		return nil, fmt.Errorf("mandatory purpose %q is required", vo.PurposeNecessary)
	}

	return result, nil
}

// ValidateRevoke ตรวจสอบว่าสามารถถอนความยินยอมได้หรือไม่
// ValidateRevoke validates whether a consent can be revoked
func (v *ConsentValidator) ValidateRevoke(purpose vo.ConsentPurpose) error {
	if !purpose.IsValid() {
		return fmt.Errorf("invalid purpose: %s", purpose)
	}
	if purpose.IsRequired() {
		return fmt.Errorf("cannot revoke mandatory purpose: %s", purpose)
	}
	return nil
}
```

### 3.3.2 DeletionPolicyService

**`internal/modules/pdpa/domain/service/deletion_policy_service.go`**

```go
package service

import (
	"time"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

// RetentionPolicy นโยบายการเก็บรักษาข้อมูล
// RetentionPolicy defines data retention policy
type RetentionPolicy struct {
	SuspendedRetentionYears int // จำนวนปีที่เก็บข้อมูลเมื่อถูกระงับ (default 1)
	RevokedConsentDays      int // จำนวนวันหลัง revoke ก่อน auto-delete consent log
}

// DefaultRetentionPolicy นโยบายค่าเริ่มต้น
// DefaultRetentionPolicy is the default retention policy
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		SuspendedRetentionYears: 1,
		RevokedConsentDays:      365,
	}
}

// DeletionPolicyService Domain Service สำหรับตัดสินใจการลบข้อมูล
// DeletionPolicyService is a Domain Service that decides deletion eligibility
type DeletionPolicyService struct {
	policy RetentionPolicy
}

// NewDeletionPolicyService สร้าง service ใหม่
// NewDeletionPolicyService creates a new service
func NewDeletionPolicyService(policy RetentionPolicy) *DeletionPolicyService {
	if policy.SuspendedRetentionYears < 1 {
		policy.SuspendedRetentionYears = 1
	}
	if policy.RevokedConsentDays < 1 {
		policy.RevokedConsentDays = 365
	}
	return &DeletionPolicyService{policy: policy}
}

// CanImmediateDelete ตรวจสอบว่าสามารถลบทันทีได้หรือไม่
// CanImmediateDelete checks whether immediate deletion is allowed
// (เฉพาะ account ที่ TERMINATED และยืนยันการลบแล้ว)
func (s *DeletionPolicyService) CanImmediateDelete(status *entity.UserAccountStatus) bool {
	return status.CanImmediateDelete()
}

// IsReadyForAutoDeletion ตรวจสอบว่าถึงเวลาลบอัตโนมัติหรือยัง
// IsReadyForAutoDeletion checks whether auto-deletion is due
func (s *DeletionPolicyService) IsReadyForAutoDeletion(status *entity.UserAccountStatus, now time.Time) bool {
	return status.IsReadyForAutoDeletion(now)
}

// CalculateRetentionDeadline คำนวณวันหมดอายุการเก็บข้อมูล
// CalculateRetentionDeadline computes the retention deadline
func (s *DeletionPolicyService) CalculateRetentionDeadline(suspendedAt time.Time) time.Time {
	return suspendedAt.AddDate(s.policy.SuspendedRetentionYears, 0, 0)
}

// ConsentAutoDeleteCutoff คำนวณ cutoff สำหรับ consent ที่ถูก revoke
// ConsentAutoDeleteCutoff computes the cutoff for revoked consents
func (s *DeletionPolicyService) ConsentAutoDeleteCutoff(now time.Time) time.Time {
	return now.AddDate(0, 0, -s.policy.RevokedConsentDays)
}

// Policy คืนค่านโยบายปัจจุบัน
// Policy returns the current retention policy
func (s *DeletionPolicyService) Policy() RetentionPolicy {
	return s.policy
}
```

### 3.3.3 AnonymizationService

**`internal/modules/pdpa/domain/service/anonymization_service.go`**

```go
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// AnonymizationService Domain Service สำหรับ anonymize ข้อมูลส่วนบุคคล
// AnonymizationService is a Domain Service for PII anonymization
type AnonymizationService struct {
	salt string
}

// NewAnonymizationService สร้าง service ใหม่พร้อม salt
// NewAnonymizationService creates a new service with a salt
func NewAnonymizationService(salt string) *AnonymizationService {
	return &AnonymizationService{salt: salt}
}

// AnonymizeEmail anonymize email address (keep domain, hash local part)
// AnonymizeEmail anonymizes an email address (keeps domain, hashes local part)
// Example: "john.doe@example.com" -> "a3f8...@example.com"
func (s *AnonymizationService) AnonymizeEmail(email string) string {
	if email == "" {
		return ""
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return s.hashString(email)
	}
	return fmt.Sprintf("%s@%s", s.hashString(parts[0])[:12], parts[1])
}

// AnonymizePhone anonymize เบอร์โทรศัพท์ (keep last 4 digits)
// AnonymizePhone anonymizes a phone number (keeps last 4 digits)
// Example: "0812345678" -> "******5678"
func (s *AnonymizationService) AnonymizePhone(phone string) string {
	if len(phone) <= 4 {
		return strings.Repeat("*", len(phone))
	}
	return strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
}

// AnonymizeIP anonymize IP address (mask last octet for IPv4)
// AnonymizeIP anonymizes an IP address (masks last octet for IPv4)
// Example: "192.168.1.100" -> "192.168.1.0"
func (s *AnonymizationService) AnonymizeIP(ip string) string {
	if ip == "" {
		return ""
	}
	// IPv4
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			parts[3] = "0"
			return strings.Join(parts, ".")
		}
	}
	// IPv6 or other - hash
	return s.hashString(ip)[:16]
}

// AnonymizeUserAgent anonymize user agent (keep only browser + OS family)
// AnonymizeUserAgent anonymizes user agent (keeps only browser + OS family)
func (s *AnonymizationService) AnonymizeUserAgent(ua string) string {
	if ua == "" {
		return ""
	}
	// Keep only first 40 chars, hash rest
	if len(ua) > 40 {
		return ua[:40] + "..."
	}
	return ua
}

// HashUserID สร้าง deterministic hash ของ user ID (ใช้เป็น pseudonym)
// HashUserID generates a deterministic hash of user ID (pseudonym)
func (s *AnonymizationService) HashUserID(userID string) string {
	return s.hashString(userID)[:16]
}

// hashString hash string ด้วย SHA-256 + salt
// hashString hashes a string with SHA-256 + salt
func (s *AnonymizationService) hashString(input string) string {
	h := sha256.Sum256([]byte(s.salt + input))
	return hex.EncodeToString(h[:])
}
```

---

## 3.4 Application Layer — Commands

### 3.4.1 RecordConsentCommand

**`internal/modules/pdpa/application/command/record_consent.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// RecordConsentCommand คำสั่งสำหรับบันทึกความยินยอม
// RecordConsentCommand is the command to record consent
type RecordConsentCommand struct {
	UserID    uuid.UUID
	SessionID string
	Purposes  map[string]bool
	IPAddress string
	UserAgent string
	TraceID   string
}

// RecordConsentResult ผลลัพธ์จาก command
// RecordConsentResult is the command result
type RecordConsentResult struct {
	RecordedPurposes []string
	ConsentIDs       []uuid.UUID
}

// RecordConsentHandler handles the RecordConsentCommand
type RecordConsentHandler struct {
	consentRepo repository.ConsentRepository
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	validator   *service.ConsentValidator
	log         logger.Logger
}

// NewRecordConsentHandler สร้าง handler ใหม่
// NewRecordConsentHandler creates a new handler
func NewRecordConsentHandler(
	consentRepo repository.ConsentRepository,
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	validator *service.ConsentValidator,
	log logger.Logger,
) *RecordConsentHandler {
	return &RecordConsentHandler{
		consentRepo: consentRepo,
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		validator:   validator,
		log:         log,
	}
}

// Handle ประมวลผล command
// Handle executes the command
func (h *RecordConsentHandler) Handle(ctx context.Context, cmd RecordConsentCommand) (*RecordConsentResult, error) {
	// 1. ตรวจสอบว่าบัญชียัง active / Verify account is active
	if status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID); err == nil {
		if !status.Status.IsActive() {
			return nil, domainerrors.ErrAccountNotActive
		}
	}

	// 2. Validate purposes / Validate the purposes
	purposes, err := h.validator.ValidatePurposes(cmd.Purposes)
	if err != nil {
		return nil, err
	}

	result := &RecordConsentResult{
		RecordedPurposes: make([]string, 0, len(purposes)),
		ConsentIDs:       make([]uuid.UUID, 0, len(purposes)),
	}

	// 3. Loop purposes and save
	for _, purpose := range purposes {
		log, err := entity.NewConsentLog(cmd.UserID, cmd.SessionID, purpose, cmd.IPAddress, cmd.UserAgent)
		if err != nil {
			return nil, err
		}

		if err := h.consentRepo.Save(ctx, log); err != nil {
			h.log.Error("failed to save consent", "error", err, "user_id", cmd.UserID)
			return nil, err
		}

		// 4. Build domain event / สร้าง domain event
		meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
		evt := event.ConsentGrantedEvent{
			Metadata:  meta,
			ConsentID: log.ID,
			UserID:    log.UserID,
			SessionID: log.SessionID,
			Purpose:   log.Purpose.String(),
			GrantedAt: log.GrantedAt,
			ExpiresAt: log.ExpiresAt,
			IPAddress: log.IPAddress,
			UserAgent: log.UserAgent,
		}

		// 5. Save to outbox (transactional outbox pattern)
		if err := h.saveOutbox(ctx, "ConsentLog", log.ID, evt, event.TopicConsentGranted); err != nil {
			return nil, err
		}

		result.RecordedPurposes = append(result.RecordedPurposes, purpose.String())
		result.ConsentIDs = append(result.ConsentIDs, log.ID)
	}

	// 6. Audit trail
	audit := entity.NewAuditTrailWithMeta(
		&cmd.UserID,
		entity.ActionConsentGranted,
		map[string]interface{}{"purposes": result.RecordedPurposes},
		cmd.IPAddress,
		cmd.UserAgent,
	)
	if err := h.auditRepo.Save(ctx, audit); err != nil {
		h.log.Warn("failed to save audit trail", "error", err)
	}

	return result, nil
}

// saveOutbox บันทึก event ลง outbox
// saveOutbox saves the event to the outbox table
func (h *RecordConsentHandler) saveOutbox(
	ctx context.Context,
	aggregateType string,
	aggregateID uuid.UUID,
	evt interface{},
	topic string,
) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     topic,
		Topic:         topic,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})
}
```

### 3.4.2 RevokeConsentCommand

**`internal/modules/pdpa/application/command/revoke_consent.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// RevokeConsentCommand คำสั่งสำหรับถอนความยินยอม
// RevokeConsentCommand is the command to revoke consent
type RevokeConsentCommand struct {
	UserID    uuid.UUID
	Purpose   vo.ConsentPurpose
	IPAddress string
	UserAgent string
	TraceID   string
}

// RevokeConsentHandler handles RevokeConsentCommand
type RevokeConsentHandler struct {
	consentRepo repository.ConsentRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	validator   *service.ConsentValidator
	log         logger.Logger
}

// NewRevokeConsentHandler สร้าง handler ใหม่
// NewRevokeConsentHandler creates a new handler
func NewRevokeConsentHandler(
	consentRepo repository.ConsentRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	validator *service.ConsentValidator,
	log logger.Logger,
) *RevokeConsentHandler {
	return &RevokeConsentHandler{
		consentRepo: consentRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		validator:   validator,
		log:         log,
	}
}

// Handle ประมวลผล command
// Handle executes the command
func (h *RevokeConsentHandler) Handle(ctx context.Context, cmd RevokeConsentCommand) error {
	// 1. Validate / ตรวจสอบ
	if err := h.validator.ValidateRevoke(cmd.Purpose); err != nil {
		return err
	}

	// 2. Find active consent / ค้นหา consent ที่ active
	consent, err := h.consentRepo.FindActiveByUserAndPurpose(ctx, cmd.UserID, cmd.Purpose)
	if err != nil {
		return err
	}
	if consent == nil {
		return domainerrors.ErrConsentNotFound
	}

	// 3. Revoke via behavior method / ถอนผ่าน behavior method
	if err := consent.Revoke(cmd.IPAddress, cmd.UserAgent); err != nil {
		return err
	}

	// 4. Persist / บันทึก
	if err := h.consentRepo.Update(ctx, consent); err != nil {
		return err
	}

	// 5. Publish domain event via outbox
	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.ConsentRevokedEvent{
		Metadata:  meta,
		ConsentID: consent.ID,
		UserID:    consent.UserID,
		Purpose:   consent.Purpose.String(),
		RevokedAt: *consent.RevokedAt,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
	}
	payload, _ := json.Marshal(evt)
	if err := h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "ConsentLog",
		AggregateID:   consent.ID,
		EventType:     event.TopicConsentRevoked,
		Topic:         event.TopicConsentRevoked,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	}); err != nil {
		return err
	}

	// 6. Audit
	audit := entity.NewAuditTrailWithMeta(
		&cmd.UserID,
		entity.ActionConsentRevoked,
		map[string]interface{}{"purpose": cmd.Purpose.String()},
		cmd.IPAddress,
		cmd.UserAgent,
	)
	_ = h.auditRepo.Save(ctx, audit)

	return nil
}
```

### 3.4.3 SubmitDSARCommand

**`internal/modules/pdpa/application/command/submit_dsar.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// MaxDSARPerDay จำนวนคำร้องสูงสุดต่อ user ต่อวัน
// MaxDSARPerDay is the max number of DSARs per user per day
const MaxDSARPerDay = 3

// SubmitDSARCommand คำสั่งสำหรับส่งคำร้อง DSAR
// SubmitDSARCommand submits a new DSAR
type SubmitDSARCommand struct {
	UserID      uuid.UUID
	RequestType vo.DSARType
	IPAddress   string
	UserAgent   string
	TraceID     string
}

// SubmitDSARResult ผลลัพธ์
// SubmitDSARResult is the result
type SubmitDSARResult struct {
	RequestID   uuid.UUID
	Status      string
	RequestedAt time.Time
}

// SubmitDSARHandler handles SubmitDSARCommand
type SubmitDSARHandler struct {
	dsarRepo    repository.DSARRepository
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	log         logger.Logger
}

// NewSubmitDSARHandler สร้าง handler ใหม่
// NewSubmitDSARHandler creates a new handler
func NewSubmitDSARHandler(
	dsarRepo repository.DSARRepository,
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *SubmitDSARHandler {
	return &SubmitDSARHandler{
		dsarRepo:    dsarRepo,
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		log:         log,
	}
}

// Handle ประมวลผล command
// Handle executes the command
func (h *SubmitDSARHandler) Handle(ctx context.Context, cmd SubmitDSARCommand) (*SubmitDSARResult, error) {
	// 1. Rate limit check / ตรวจสอบ rate limit
	since := time.Now().UTC().Add(-24 * time.Hour)
	count, err := h.dsarRepo.CountByUserAndTypeInPeriod(ctx, cmd.UserID, cmd.RequestType, since)
	if err != nil {
		return nil, err
	}
	if count >= MaxDSARPerDay {
		return nil, domainerrors.ErrRateLimitExceeded
	}

	// 2. Check account is active / ตรวจสอบบัญชี
	if status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID); err == nil {
		if status.Status.IsDeleted() {
			return nil, domainerrors.ErrAccountNotActive
		}
	}

	// 3. Create aggregate / สร้าง aggregate
	req, otp, err := entity.NewDSARRequest(cmd.UserID, cmd.RequestType, cmd.IPAddress, cmd.UserAgent)
	if err != nil {
		return nil, err
	}

	// 4. Persist / บันทึก
	if err := h.dsarRepo.Save(ctx, req); err != nil {
		return nil, err
	}

	// 5. Publish event with OTP (in plaintext inside event only for email worker)
	// OTP ถูกส่งผ่าน event เพื่อให้ email worker ส่งให้ user
	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.DSARSubmittedEvent{
		Metadata:    meta,
		DSARID:      req.ID,
		UserID:      req.UserID,
		RequestType: req.RequestType.String(),
		RequestedAt: req.RequestedAt,
		IPAddress:   cmd.IPAddress,
		UserAgent:   cmd.UserAgent,
	}
	payload, _ := json.Marshal(evt)
	if err := h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "DSARRequest",
		AggregateID:   req.ID,
		EventType:     event.TopicDSARSubmitted,
		Topic:         event.TopicDSARSubmitted,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

	// 6. Send OTP via email event (different topic)
	emailEvt := event.EmailNotificationEvent{
		Metadata:  meta,
		Template:  "dsar_otp",
		Locale:    "th",
		Variables: map[string]interface{}{"otp": otp, "request_type": cmd.RequestType.String()},
		RefType:   "DSARRequest",
		RefID:     req.ID.String(),
		CreatedAt: time.Now().UTC(),
	}
	// Note: recipient email จะถูกเติมโดย email worker จาก user service
	emailPayload, _ := json.Marshal(emailEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "DSARRequest",
		AggregateID:   req.ID,
		EventType:     event.TopicEmailNotification,
		Topic:         event.TopicEmailNotification,
		Payload:       emailPayload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	// 7. Audit
	audit := entity.NewAuditTrailWithMeta(
		&cmd.UserID,
		entity.ActionDSARSubmitted,
		map[string]interface{}{
			"request_id":   req.ID.String(),
			"request_type": cmd.RequestType.String(),
		},
		cmd.IPAddress,
		cmd.UserAgent,
	)
	_ = h.auditRepo.Save(ctx, audit)

	return &SubmitDSARResult{
		RequestID:   req.ID,
		Status:      string(req.Status),
		RequestedAt: req.RequestedAt,
	}, nil
}
```

### 3.4.4 ProcessDSARCommand

**`internal/modules/pdpa/application/command/process_dsar.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// DSARAction การกระทำที่ admin เลือก
// DSARAction is the action selected by admin
type DSARAction string

const (
	DSARActionVerifyOTP DSARAction = "VERIFY_OTP"
	DSARActionApprove   DSARAction = "APPROVE"
	DSARActionReject    DSARAction = "REJECT"
)

// ProcessDSARCommand คำสั่งประมวลผล DSAR (admin)
// ProcessDSARCommand processes a DSAR (admin only)
type ProcessDSARCommand struct {
	DSARID        uuid.UUID
	Action        DSARAction
	OTPCode       string // required when Action = VERIFY_OTP
	Payload       []byte // required when Action = APPROVE
	RejectReason  string // required when Action = REJECT
	OperatorID    uuid.UUID
	TraceID       string
}

// ProcessDSARHandler handles ProcessDSARCommand
type ProcessDSARHandler struct {
	dsarRepo   repository.DSARRepository
	auditRepo  repository.AuditRepository
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewProcessDSARHandler สร้าง handler ใหม่
// NewProcessDSARHandler creates a new handler
func NewProcessDSARHandler(
	dsarRepo repository.DSARRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *ProcessDSARHandler {
	return &ProcessDSARHandler{
		dsarRepo:   dsarRepo,
		auditRepo:  auditRepo,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล command
// Handle executes the command
func (h *ProcessDSARHandler) Handle(ctx context.Context, cmd ProcessDSARCommand) error {
	req, err := h.dsarRepo.FindByID(ctx, cmd.DSARID)
	if err != nil {
		return err
	}
	if req.Status.IsTerminal() {
		return domainerrors.ErrDSARAlreadyProcessed
	}

	switch cmd.Action {
	case DSARActionVerifyOTP:
		if err := req.VerifyOTP(cmd.OTPCode); err != nil {
			return err
		}
		if err := req.MarkProcessing(); err != nil {
			return err
		}
		// mark OTP verified separately could be handled by Update
	case DSARActionApprove:
		if err := req.MarkCompleted(cmd.Payload); err != nil {
			return err
		}
	case DSARActionReject:
		if err := req.MarkRejected(cmd.RejectReason); err != nil {
			return err
		}
	default:
		return domainerrors.ErrInvalidInput
	}

	if err := h.dsarRepo.Update(ctx, req); err != nil {
		return err
	}

	// Publish completion event เมื่อ approve หรือ reject
	if req.Status == vo.DSARStatusCompleted || req.Status == vo.DSARStatusRejected {
		meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
		evt := event.DSARCompletedEvent{
			Metadata:    meta,
			DSARID:      req.ID,
			UserID:      req.UserID,
			RequestType: req.RequestType.String(),
			Status:      req.Status.String(),
			CompletedAt: time.Now().UTC(),
			DataPayload: req.DataPayload,
			Reason:      req.RejectionReason,
		}
		payload, _ := json.Marshal(evt)
		_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
			ID:            uuid.New(),
			AggregateType: "DSARRequest",
			AggregateID:   req.ID,
			EventType:     event.TopicDSARCompleted,
			Topic:         event.TopicDSARCompleted,
			Payload:       payload,
			Status:        repository.OutboxStatusPending,
			CreatedAt:     time.Now().UTC(),
		})
	}

	return nil
}
```

### 3.4.5 ConfirmDeletionCommand

**`internal/modules/pdpa/application/command/confirm_deletion.go`**

```go
package command

import (
	"context"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ConfirmDeletionCommand ยืนยันการลบข้อมูล (admin/DPO)
// ConfirmDeletionCommand confirms data deletion (admin/DPO)
type ConfirmDeletionCommand struct {
	UserID     uuid.UUID
	OperatorID uuid.UUID
	TraceID    string
}

// ConfirmDeletionHandler handles ConfirmDeletionCommand
type ConfirmDeletionHandler struct {
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	log         logger.Logger
}

// NewConfirmDeletionHandler สร้าง handler ใหม่
// NewConfirmDeletionHandler creates a new handler
func NewConfirmDeletionHandler(
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	log logger.Logger,
) *ConfirmDeletionHandler {
	return &ConfirmDeletionHandler{
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		log:         log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *ConfirmDeletionHandler) Handle(ctx context.Context, cmd ConfirmDeletionCommand) error {
	status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if err := status.ConfirmDeletion(); err != nil {
		if err == domainerrors.ErrAccountNotTerminated {
			return err
		}
		return err
	}
	if err := h.accountRepo.Save(ctx, status); err != nil {
		return err
	}
	return nil
}
```

### 3.4.6 ImmediateDeletionCommand

**`internal/modules/pdpa/application/command/immediate_deletion.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

// ImmediateDeletionCommand ลบข้อมูลทันที (หลัง confirm)
// ImmediateDeletionCommand deletes data immediately (after confirmation)
type ImmediateDeletionCommand struct {
	UserID     uuid.UUID
	Reason     string
	OperatorID uuid.UUID
	TraceID    string
}

// ImmediateDeletionHandler handles ImmediateDeletionCommand
type ImmediateDeletionHandler struct {
	consentRepo repository.ConsentRepository
	dsarRepo    repository.DSARRepository
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	policy      *service.DeletionPolicyService
	log         logger.Logger
}

// NewImmediateDeletionHandler สร้าง handler ใหม่
// NewImmediateDeletionHandler creates a new handler
func NewImmediateDeletionHandler(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	policy *service.DeletionPolicyService,
	log logger.Logger,
) *ImmediateDeletionHandler {
	return &ImmediateDeletionHandler{
		consentRepo: consentRepo,
		dsarRepo:    dsarRepo,
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		policy:      policy,
		log:         log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *ImmediateDeletionHandler) Handle(ctx context.Context, cmd ImmediateDeletionCommand) error {
	// 1. Verify policy allows immediate deletion
	status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if !h.policy.CanImmediateDelete(status) {
		return errImmediateDeletionNotAllowed()
	}

	// 2. Delete consent logs
	if err := h.consentRepo.DeleteByUserID(ctx, cmd.UserID); err != nil {
		return err
	}
	// 3. Delete DSARs
	if err := h.dsarRepo.DeleteByUserID(ctx, cmd.UserID); err != nil {
		return err
	}

	// 4. Mark account as deleted
	status.MarkDeleted()
	if err := h.accountRepo.Save(ctx, status); err != nil {
		return err
	}

	// 5. Publish DataDeleted event
	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.DataDeletedEvent{
		Metadata:   meta,
		UserID:     cmd.UserID,
		DeletedAt:  time.Now().UTC(),
		DeleteType: "IMMEDIATE",
		Reason:     cmd.Reason,
	}
	payload, _ := json.Marshal(evt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "UserAccountStatus",
		AggregateID:   cmd.UserID,
		EventType:     event.TopicDataDeleted,
		Topic:         event.TopicDataDeleted,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	// 6. Audit
	audit := entity.NewAuditTrail(&cmd.UserID, entity.ActionImmediateDeletion, map[string]interface{}{
		"reason":      cmd.Reason,
		"operator_id": cmd.OperatorID.String(),
	})
	_ = h.auditRepo.Save(ctx, audit)

	return nil
}

// error helper เลี่ยง cyclic import
// error helper to avoid cyclic import
func errImmediateDeletionNotAllowed() error {
	return errImmediateDeletion
}
```

> **หมายเหตุ:** ต้องเพิ่ม sentinel error ที่ `application/command/errors.go` เพื่อหลีกเลี่ยง circular import

**`internal/modules/pdpa/application/command/errors.go`**

```go
package command

import "errors"

// errImmediateDeletion ใช้ภายใน handler เพื่อไม่ให้ import domainerrors ซ้ำ
// internal helper error
var errImmediateDeletion = errors.New("immediate deletion is not allowed for this account")
```

### 3.4.7 AutoDeleteExpiredCommand

**`internal/modules/pdpa/application/command/auto_delete_expired.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

// AutoDeleteExpiredCommand ลบข้อมูลที่หมดอายุอัตโนมัติ (cron job)
// AutoDeleteExpiredCommand auto-deletes expired data (cron job)
type AutoDeleteExpiredCommand struct {
	BatchSize int
	TraceID   string
}

// AutoDeleteExpiredResult ผลลัพธ์
// AutoDeleteExpiredResult is the result
type AutoDeleteExpiredResult struct {
	AccountsDeleted int
	ConsentsDeleted int
	Errors          []string
}

// AutoDeleteExpiredHandler handles AutoDeleteExpiredCommand
type AutoDeleteExpiredHandler struct {
	consentRepo repository.ConsentRepository
	dsarRepo    repository.DSARRepository
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	policy      *service.DeletionPolicyService
	log         logger.Logger
}

// NewAutoDeleteExpiredHandler สร้าง handler ใหม่
// NewAutoDeleteExpiredHandler creates a new handler
func NewAutoDeleteExpiredHandler(
	consentRepo repository.ConsentRepository,
	dsarRepo repository.DSARRepository,
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	policy *service.DeletionPolicyService,
	log logger.Logger,
) *AutoDeleteExpiredHandler {
	return &AutoDeleteExpiredHandler{
		consentRepo: consentRepo,
		dsarRepo:    dsarRepo,
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		policy:      policy,
		log:         log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *AutoDeleteExpiredHandler) Handle(ctx context.Context, cmd AutoDeleteExpiredCommand) (*AutoDeleteExpiredResult, error) {
	if cmd.BatchSize <= 0 {
		cmd.BatchSize = 100
	}
	now := time.Now().UTC()
	result := &AutoDeleteExpiredResult{}

	// Phase 1: Auto-delete consent logs ที่ revoke เกิน retention
	// Phase 1: Auto-delete revoked consent logs past retention
	cutoff := h.policy.ConsentAutoDeleteCutoff(now)
	expiredConsents, err := h.consentRepo.FindReadyForAutoDeletion(ctx, cutoff, cmd.BatchSize)
	if err != nil {
		return nil, err
	}
	for i := range expiredConsents {
		c := &expiredConsents[i]
		c.MarkDeleted()
		if err := h.consentRepo.Update(ctx, c); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		result.ConsentsDeleted++
	}

	// Phase 2: Auto-delete suspended accounts past retention deadline
	// Phase 2: Auto-delete suspended accounts
	expiredAccounts, err := h.accountRepo.FindReadyForAutoDeletion(ctx, now, cmd.BatchSize)
	if err != nil {
		return nil, err
	}
	for i := range expiredAccounts {
		st := &expiredAccounts[i]
		if !h.policy.IsReadyForAutoDeletion(st, now) {
			continue
		}
		if err := h.consentRepo.DeleteByUserID(ctx, st.UserID); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}
		if err := h.dsarRepo.DeleteByUserID(ctx, st.UserID); err != nil {
			result.Errors = append(result.Errors, err.Error())
		}
		st.MarkDeleted()
		if err := h.accountRepo.Save(ctx, st); err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}

		// Publish event
		meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
		evt := event.DataDeletedEvent{
			Metadata:   meta,
			UserID:     st.UserID,
			DeletedAt:  now,
			DeleteType: "AUTO",
			Reason:     "retention_policy_expired",
		}
		payload, _ := json.Marshal(evt)
		_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
			ID:            uuid.New(),
			AggregateType: "UserAccountStatus",
			AggregateID:   st.UserID,
			EventType:     event.TopicDataDeleted,
			Topic:         event.TopicDataDeleted,
			Payload:       payload,
			Status:        repository.OutboxStatusPending,
			CreatedAt:     now,
		})

		// Audit
		audit := entity.NewAuditTrail(&st.UserID, entity.ActionAutoDeletion, map[string]interface{}{
			"retention_deadline": st.RetentionDeadline,
		})
		_ = h.auditRepo.Save(ctx, audit)

		result.AccountsDeleted++
	}

	return result, nil
}
```

### 3.4.8 HandleAccountSuspendedCommand

**`internal/modules/pdpa/application/command/handle_account_suspended.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/logger"
)

// HandleAccountSuspendedCommand จัดการเมื่อบัญชีถูกระงับ
// HandleAccountSuspendedCommand handles account suspension
type HandleAccountSuspendedCommand struct {
	UserID         uuid.UUID
	SuspendedAt    time.Time
	RetentionYears int
	TraceID        string
}

// HandleAccountSuspendedHandler handles HandleAccountSuspendedCommand
type HandleAccountSuspendedHandler struct {
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	policy      *service.DeletionPolicyService
	log         logger.Logger
}

// NewHandleAccountSuspendedHandler สร้าง handler ใหม่
// NewHandleAccountSuspendedHandler creates a new handler
func NewHandleAccountSuspendedHandler(
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	policy *service.DeletionPolicyService,
	log logger.Logger,
) *HandleAccountSuspendedHandler {
	return &HandleAccountSuspendedHandler{
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		policy:      policy,
		log:         log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *HandleAccountSuspendedHandler) Handle(ctx context.Context, cmd HandleAccountSuspendedCommand) error {
	years := cmd.RetentionYears
	if years < 1 {
		years = h.policy.Policy().SuspendedRetentionYears
	}

	status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil && err != domainerrors.ErrAccountNotFound {
		return err
	}
	if status == nil {
		status, err = entity.NewUserAccountStatus(cmd.UserID)
		if err != nil {
			return err
		}
	}
	if err := status.Suspend(cmd.SuspendedAt, years); err != nil {
		return err
	}
	if err := h.accountRepo.Save(ctx, status); err != nil {
		return err
	}

	// Publish event
	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.AccountSuspendedEvent{
		Metadata:       meta,
		UserID:         cmd.UserID,
		SuspendedAt:    cmd.SuspendedAt,
		RetentionYears: years,
		Deadline:       *status.RetentionDeadline,
	}
	payload, _ := json.Marshal(evt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "UserAccountStatus",
		AggregateID:   cmd.UserID,
		EventType:     event.TopicAccountSuspended,
		Topic:         event.TopicAccountSuspended,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	audit := entity.NewAuditTrail(&cmd.UserID, entity.ActionAccountSuspended, map[string]interface{}{
		"suspended_at":       cmd.SuspendedAt,
		"retention_deadline": status.RetentionDeadline,
	})
	_ = h.auditRepo.Save(ctx, audit)

	return nil
}
```

### 3.4.9 HandleAccountTerminatedCommand

**`internal/modules/pdpa/application/command/handle_account_terminated.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// HandleAccountTerminatedCommand จัดการเมื่อบัญชีถูกยกเลิก
// HandleAccountTerminatedCommand handles account termination
type HandleAccountTerminatedCommand struct {
	UserID       uuid.UUID
	TerminatedAt time.Time
	TraceID      string
}

// HandleAccountTerminatedHandler handles HandleAccountTerminatedCommand
type HandleAccountTerminatedHandler struct {
	accountRepo repository.UserAccountStatusRepository
	auditRepo   repository.AuditRepository
	outboxRepo  repository.OutboxRepository
	log         logger.Logger
}

// NewHandleAccountTerminatedHandler สร้าง handler ใหม่
// NewHandleAccountTerminatedHandler creates a new handler
func NewHandleAccountTerminatedHandler(
	accountRepo repository.UserAccountStatusRepository,
	auditRepo repository.AuditRepository,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *HandleAccountTerminatedHandler {
	return &HandleAccountTerminatedHandler{
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
		outboxRepo:  outboxRepo,
		log:         log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *HandleAccountTerminatedHandler) Handle(ctx context.Context, cmd HandleAccountTerminatedCommand) error {
	status, err := h.accountRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil && err != domainerrors.ErrAccountNotFound {
		return err
	}
	if status == nil {
		status, err = entity.NewUserAccountStatus(cmd.UserID)
		if err != nil {
			return err
		}
	}
	if err := status.Terminate(cmd.TerminatedAt); err != nil {
		return err
	}
	if err := h.accountRepo.Save(ctx, status); err != nil {
		return err
	}

	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.AccountTerminatedEvent{
		Metadata:     meta,
		UserID:       cmd.UserID,
		TerminatedAt: cmd.TerminatedAt,
	}
	payload, _ := json.Marshal(evt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "UserAccountStatus",
		AggregateID:   cmd.UserID,
		EventType:     event.TopicAccountTerminated,
		Topic:         event.TopicAccountTerminated,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	audit := entity.NewAuditTrail(&cmd.UserID, entity.ActionAccountTerminated, map[string]interface{}{
		"terminated_at": cmd.TerminatedAt,
	})
	_ = h.auditRepo.Save(ctx, audit)

	return nil
}
```

### 3.4.10 PublishPrivacyPolicyCommand

**`internal/modules/pdpa/application/command/publish_privacy_policy.go`**

```go
package command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// PublishPrivacyPolicyCommand เผยแพร่นโยบายใหม่
// PublishPrivacyPolicyCommand publishes a new privacy policy
type PublishPrivacyPolicyCommand struct {
	Version       string
	Title         string
	Content       string
	EffectiveDate time.Time
	OperatorID    uuid.UUID
	TraceID       string
}

// PublishPrivacyPolicyResult ผลลัพธ์
// PublishPrivacyPolicyResult is the result
type PublishPrivacyPolicyResult struct {
	PolicyID uuid.UUID
	Version  string
}

// PublishPrivacyPolicyHandler handles PublishPrivacyPolicyCommand
type PublishPrivacyPolicyHandler struct {
	policyRepo repository.PolicyRepository
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewPublishPrivacyPolicyHandler สร้าง handler ใหม่
// NewPublishPrivacyPolicyHandler creates a new handler
func NewPublishPrivacyPolicyHandler(
	policyRepo repository.PolicyRepository,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *PublishPrivacyPolicyHandler {
	return &PublishPrivacyPolicyHandler{
		policyRepo: policyRepo,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle executes
func (h *PublishPrivacyPolicyHandler) Handle(ctx context.Context, cmd PublishPrivacyPolicyCommand) (*PublishPrivacyPolicyResult, error) {
	// ตรวจสอบว่า version ซ้ำหรือไม่
	// Check if version exists
	if existing, _ := h.policyRepo.FindByVersion(ctx, cmd.Version); existing != nil {
		return nil, domainerrors.ErrPolicyVersionExists
	}

	policy, err := entity.NewPrivacyPolicy(cmd.Version, cmd.Title, cmd.Content, cmd.EffectiveDate)
	if err != nil {
		return nil, err
	}
	policy.Activate()

	// ปิดการใช้งานนโยบายเก่า
	// Deactivate old policies
	if err := h.policyRepo.DeactivateAll(ctx); err != nil {
		return nil, err
	}
	if err := h.policyRepo.Save(ctx, policy); err != nil {
		return nil, err
	}

	// Publish event
	meta := event.NewEventMetadata(cmd.TraceID, "", cmd.TraceID)
	evt := event.PrivacyPolicyPublishedEvent{
		Metadata:      meta,
		PolicyID:      policy.ID,
		Version:       policy.Version,
		Title:         policy.Title,
		EffectiveDate: policy.EffectiveDate,
		PublishedAt:   time.Now().UTC(),
	}
	payload, _ := json.Marshal(evt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "PrivacyPolicy",
		AggregateID:   policy.ID,
		EventType:     evt.EventName(),
		Topic:         evt.EventName(),
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	return &PublishPrivacyPolicyResult{
		PolicyID: policy.ID,
		Version:  policy.Version,
	}, nil
}
```

---

## 3.5 Application Layer — Queries

### 3.5.1 GetConsentHistoryQuery

**`internal/modules/pdpa/application/query/get_consent_history.go`**

```go
package query

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// GetConsentHistoryQuery ดึงประวัติความยินยอม
// GetConsentHistoryQuery retrieves consent history
type GetConsentHistoryQuery struct {
	UserID uuid.UUID
}

// GetConsentHistoryHandler handles GetConsentHistoryQuery
type GetConsentHistoryHandler struct {
	consentRepo repository.ConsentRepository
}

// NewGetConsentHistoryHandler สร้าง handler ใหม่
// NewGetConsentHistoryHandler creates a new handler
func NewGetConsentHistoryHandler(consentRepo repository.ConsentRepository) *GetConsentHistoryHandler {
	return &GetConsentHistoryHandler{consentRepo: consentRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *GetConsentHistoryHandler) Handle(ctx context.Context, q GetConsentHistoryQuery) ([]dto.ConsentHistoryItem, error) {
	items, err := h.consentRepo.FindByUserID(ctx, q.UserID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ConsentHistoryItem, 0, len(items))
	for _, c := range items {
		result = append(result, dto.ConsentHistoryItem{
			ID:            c.ID,
			Purpose:       c.Purpose.String(),
			Status:        c.Status.String(),
			GrantedAt:     c.GrantedAt,
			ExpiresAt:     c.ExpiresAt,
			RevokedAt:     c.RevokedAt,
			AutoDeletedAt: c.AutoDeletedAt,
		})
	}
	return result, nil
}
```

### 3.5.2 GetDSARStatusQuery

**`internal/modules/pdpa/application/query/get_dsar_status.go`**

```go
package query

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// GetDSARStatusQuery ดึงสถานะคำร้อง
// GetDSARStatusQuery retrieves DSAR status
type GetDSARStatusQuery struct {
	DSARID uuid.UUID
	UserID uuid.UUID // ถ้าไม่ใช่ admin ต้องตรงกับ owner
}

// GetDSARStatusHandler handles GetDSARStatusQuery
type GetDSARStatusHandler struct {
	dsarRepo repository.DSARRepository
}

// NewGetDSARStatusHandler สร้าง handler ใหม่
// NewGetDSARStatusHandler creates a new handler
func NewGetDSARStatusHandler(dsarRepo repository.DSARRepository) *GetDSARStatusHandler {
	return &GetDSARStatusHandler{dsarRepo: dsarRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *GetDSARStatusHandler) Handle(ctx context.Context, q GetDSARStatusQuery) (*dto.DSARStatusResponse, error) {
	req, err := h.dsarRepo.FindByID(ctx, q.DSARID)
	if err != nil {
		return nil, err
	}
	if q.UserID != uuid.Nil && req.UserID != q.UserID {
		return nil, nil // caller must check and return 403
	}
	return &dto.DSARStatusResponse{
		ID:              req.ID,
		UserID:          req.UserID,
		RequestType:     req.RequestType.String(),
		Status:          req.Status.String(),
		RequestedAt:     req.RequestedAt,
		CompletedAt:     req.CompletedAt,
		RejectionReason: req.RejectionReason,
	}, nil
}

// ListDSARQuery ดึงรายการ DSAR
// ListDSARQuery lists DSARs for a user or admin
type ListDSARQuery struct {
	UserID *uuid.UUID
	Status string
	Limit  int
	Offset int
}

// ListDSARHandler handles ListDSARQuery
type ListDSARHandler struct {
	dsarRepo repository.DSARRepository
}

// NewListDSARHandler สร้าง handler ใหม่
// NewListDSARHandler creates a new handler
func NewListDSARHandler(dsarRepo repository.DSARRepository) *ListDSARHandler {
	return &ListDSARHandler{dsarRepo: dsarRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *ListDSARHandler) Handle(ctx context.Context, q ListDSARQuery) ([]dto.DSARStatusResponse, error) {
	var reqs []dto.DSARStatusResponse
	if q.UserID != nil {
		items, err := h.dsarRepo.FindByUserID(ctx, *q.UserID)
		if err != nil {
			return nil, err
		}
		for _, r := range items {
			reqs = append(reqs, dto.DSARStatusResponse{
				ID:              r.ID,
				UserID:          r.UserID,
				RequestType:     r.RequestType.String(),
				Status:          r.Status.String(),
				RequestedAt:     r.RequestedAt,
				CompletedAt:     r.CompletedAt,
				RejectionReason: r.RejectionReason,
			})
		}
	}
	return reqs, nil
}
```

### 3.5.3 GetAdminReportQuery

**`internal/modules/pdpa/application/query/get_admin_report.go`**

```go
package query

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// GetAdminReportQuery ดึงรายงานสำหรับ admin
// GetAdminReportQuery generates admin reports
type GetAdminReportQuery struct {
	From time.Time
	To   time.Time
}

// GetAdminReportHandler handles GetAdminReportQuery
type GetAdminReportHandler struct {
	consentRepo repository.ConsentRepository
	auditRepo   repository.AuditRepository
}

// NewGetAdminReportHandler สร้าง handler ใหม่
// NewGetAdminReportHandler creates a new handler
func NewGetAdminReportHandler(
	consentRepo repository.ConsentRepository,
	auditRepo repository.AuditRepository,
) *GetAdminReportHandler {
	return &GetAdminReportHandler{consentRepo: consentRepo, auditRepo: auditRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *GetAdminReportHandler) Handle(ctx context.Context, q GetAdminReportQuery) (*dto.AdminReportResponse, error) {
	consentCounts, err := h.consentRepo.CountActiveByPurpose(ctx, q.From, q.To)
	if err != nil {
		return nil, err
	}
	actionCounts, err := h.auditRepo.CountByAction(ctx, q.From, q.To)
	if err != nil {
		return nil, err
	}
	return &dto.AdminReportResponse{
		From:              q.From,
		To:                q.To,
		ConsentByPurpose:  consentCounts,
		ActionsByType:     actionCounts,
		GeneratedAt:       time.Now().UTC(),
	}, nil
}
```

### 3.5.4 GetActivePrivacyPolicyQuery

**`internal/modules/pdpa/application/query/get_active_privacy_policy.go`**

```go
package query

import (
	"context"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// GetActivePrivacyPolicyQuery ดึงนโยบายที่ใช้งานอยู่
// GetActivePrivacyPolicyQuery retrieves the active privacy policy
type GetActivePrivacyPolicyQuery struct{}

// GetActivePrivacyPolicyHandler handles GetActivePrivacyPolicyQuery
type GetActivePrivacyPolicyHandler struct {
	policyRepo repository.PolicyRepository
}

// NewGetActivePrivacyPolicyHandler สร้าง handler ใหม่
// NewGetActivePrivacyPolicyHandler creates a new handler
func NewGetActivePrivacyPolicyHandler(policyRepo repository.PolicyRepository) *GetActivePrivacyPolicyHandler {
	return &GetActivePrivacyPolicyHandler{policyRepo: policyRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *GetActivePrivacyPolicyHandler) Handle(ctx context.Context, _ GetActivePrivacyPolicyQuery) (*dto.PrivacyPolicyResponse, error) {
	p, err := h.policyRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.PrivacyPolicyResponse{
		ID:            p.ID,
		Version:       p.Version,
		Title:         p.Title,
		Content:       p.Content,
		EffectiveDate: p.EffectiveDate,
		IsActive:      p.IsActive,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

// ListPrivacyPoliciesQuery ดึงรายการนโยบายทั้งหมด
// ListPrivacyPoliciesQuery lists all privacy policies
type ListPrivacyPoliciesQuery struct{}

// ListPrivacyPoliciesHandler handles ListPrivacyPoliciesQuery
type ListPrivacyPoliciesHandler struct {
	policyRepo repository.PolicyRepository
}

// NewListPrivacyPoliciesHandler สร้าง handler ใหม่
// NewListPrivacyPoliciesHandler creates a new handler
func NewListPrivacyPoliciesHandler(policyRepo repository.PolicyRepository) *ListPrivacyPoliciesHandler {
	return &ListPrivacyPoliciesHandler{policyRepo: policyRepo}
}

// Handle ประมวลผล
// Handle executes
func (h *ListPrivacyPoliciesHandler) Handle(ctx context.Context, _ ListPrivacyPoliciesQuery) ([]dto.PrivacyPolicyResponse, error) {
	items, err := h.policyRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PrivacyPolicyResponse, 0, len(items))
	for _, p := range items {
		out = append(out, dto.PrivacyPolicyResponse{
			ID:            p.ID,
			Version:       p.Version,
			Title:         p.Title,
			Content:       p.Content,
			EffectiveDate: p.EffectiveDate,
			IsActive:      p.IsActive,
			UpdatedAt:     p.UpdatedAt,
		})
	}
	return out, nil
}
```

---

## 3.6 Application Layer — DTOs

**`internal/modules/pdpa/application/dto/dto.go`**

```go
package dto

import (
	"time"

	"github.com/google/uuid"
)

// ConsentHistoryItem รายการประวัติความยินยอม
// ConsentHistoryItem is a single consent history record
type ConsentHistoryItem struct {
	ID            uuid.UUID  `json:"id"`
	Purpose       string     `json:"purpose"`
	Status        string     `json:"status"`
	GrantedAt     time.Time  `json:"granted_at"`
	ExpiresAt     time.Time  `json:"expires_at"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time `json:"auto_deleted_at,omitempty"`
}

// DSARStatusResponse สถานะคำร้อง DSAR
// DSARStatusResponse represents a DSAR status
type DSARStatusResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	RequestType     string     `json:"request_type"`
	Status          string     `json:"status"`
	RequestedAt     time.Time  `json:"requested_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
}

// AdminReportResponse รายงานสำหรับ admin
// AdminReportResponse is the admin report
type AdminReportResponse struct {
	From             time.Time        `json:"from"`
	To               time.Time        `json:"to"`
	ConsentByPurpose map[string]int64 `json:"consent_by_purpose"`
	ActionsByType    map[string]int64 `json:"actions_by_type"`
	GeneratedAt      time.Time        `json:"generated_at"`
}

// PrivacyPolicyResponse นโยบายความเป็นส่วนตัว
// PrivacyPolicyResponse is a privacy policy
type PrivacyPolicyResponse struct {
	ID            uuid.UUID `json:"id"`
	Version       string    `json:"version"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	EffectiveDate time.Time `json:"effective_date"`
	IsActive      bool      `json:"is_active"`
	UpdatedAt     time.Time `json:"updated_at"`
}
```

---

## 3.7 Application Layer — Event Handlers

### 3.7.1 ConsentGrantedHandler

**`internal/modules/pdpa/application/event_handler/consent_granted_handler.go`**

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// IdempotencyStore interface สำหรับเช็ค idempotency
// IdempotencyStore interface for idempotency checking
type IdempotencyStore interface {
	// IsProcessed ตรวจสอบว่า event นี้ถูกประมวลผลแล้วหรือยัง
	// IsProcessed checks if the event has been processed
	IsProcessed(ctx context.Context, eventID uuid.UUID, handlerName string) (bool, error)
	// MarkProcessed บันทึกว่า event นี้ประมวลผลแล้ว
	// MarkProcessed records that the event has been processed
	MarkProcessed(ctx context.Context, eventID uuid.UUID, handlerName string, ttl time.Duration) error
}

// ConsentCache interface สำหรับ cache
// ConsentCache interface for caching consent status
type ConsentCache interface {
	Set(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose, status vo.ConsentStatus) error
	Get(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (vo.ConsentStatus, error)
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}

// ConsentGrantedHandler จัดการ event ConsentGranted
// ConsentGrantedHandler handles ConsentGranted events
type ConsentGrantedHandler struct {
	cache      ConsentCache
	idempotent IdempotencyStore
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewConsentGrantedHandler สร้าง handler ใหม่
// NewConsentGrantedHandler creates a new handler
func NewConsentGrantedHandler(
	cache ConsentCache,
	idempotent IdempotencyStore,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *ConsentGrantedHandler {
	return &ConsentGrantedHandler{
		cache:      cache,
		idempotent: idempotent,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล event
// Handle processes the event
func (h *ConsentGrantedHandler) Handle(ctx context.Context, evt event.ConsentGrantedEvent) error {
	const handler = "ConsentGrantedHandler"

	// 1. Idempotency check
	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		h.log.Info("event already processed", "event_id", evt.Metadata.EventID)
		return nil
	}

	// 2. Update cache
	purpose := vo.ConsentPurpose(evt.Purpose)
	if err := h.cache.Set(ctx, evt.UserID, purpose, vo.ConsentGranted); err != nil {
		h.log.Error("cache set failed", "error", err)
	}

	// 3. Emit email notification event (via outbox)
	emailEvt := event.EmailNotificationEvent{
		Metadata:  evt.Metadata,
		Template:  "consent_granted",
		Locale:    "th",
		Variables: map[string]interface{}{"purpose": evt.Purpose},
		RefType:   "ConsentLog",
		RefID:     evt.ConsentID.String(),
		CreatedAt: time.Now().UTC(),
	}
	payload, _ := json.Marshal(emailEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "ConsentLog",
		AggregateID:   evt.ConsentID,
		EventType:     event.TopicEmailNotification,
		Topic:         event.TopicEmailNotification,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	// 4. Emit blockchain record event (audit log immutable)
	blockchainEvt := event.BlockchainRecordEvent{
		Metadata:  evt.Metadata,
		UserID:    evt.UserID,
		Action:    "CONSENT_GRANTED",
		DataHash:  hashEvent(evt),
		CreatedAt: time.Now().UTC(),
	}
	bcPayload, _ := json.Marshal(blockchainEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "ConsentLog",
		AggregateID:   evt.ConsentID,
		EventType:     event.TopicBlockchainRecord,
		Topic:         event.TopicBlockchainRecord,
		Payload:       bcPayload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	// 5. Mark processed
	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}

func hashEvent(evt interface{}) string {
	b, _ := json.Marshal(evt)
	return hashSHA256(string(b))
}
```

### 3.7.2 ConsentRevokedHandler

**`internal/modules/pdpa/application/event_handler/consent_revoked_handler.go`**

```go
package eventhandler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/logger"
)

// ConsentRevokedHandler handles ConsentRevokedEvent
type ConsentRevokedHandler struct {
	cache      ConsentCache
	idempotent IdempotencyStore
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewConsentRevokedHandler สร้าง handler ใหม่
// NewConsentRevokedHandler creates a new handler
func NewConsentRevokedHandler(
	cache ConsentCache,
	idempotent IdempotencyStore,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *ConsentRevokedHandler {
	return &ConsentRevokedHandler{
		cache:      cache,
		idempotent: idempotent,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *ConsentRevokedHandler) Handle(ctx context.Context, evt event.ConsentRevokedEvent) error {
	const handler = "ConsentRevokedHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	// Update cache
	purpose := vo.ConsentPurpose(evt.Purpose)
	_ = h.cache.Set(ctx, evt.UserID, purpose, vo.ConsentRevoked)

	// Emit email
	emailEvt := event.EmailNotificationEvent{
		Metadata:  evt.Metadata,
		Template:  "consent_revoked",
		Locale:    "th",
		Variables: map[string]interface{}{"purpose": evt.Purpose},
		RefType:   "ConsentLog",
		RefID:     evt.ConsentID.String(),
		CreatedAt: time.Now().UTC(),
	}
	payload, _ := json.Marshal(emailEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "ConsentLog",
		AggregateID:   evt.ConsentID,
		EventType:     event.TopicEmailNotification,
		Topic:         event.TopicEmailNotification,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}

// hashSHA256 helper
func hashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
```

### 3.7.3 DSARSubmittedHandler

**`internal/modules/pdpa/application/event_handler/dsar_submitted_handler.go`**

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DSARSubmittedHandler handles DSARSubmittedEvent
type DSARSubmittedHandler struct {
	idempotent IdempotencyStore
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewDSARSubmittedHandler สร้าง handler ใหม่
// NewDSARSubmittedHandler creates a new handler
func NewDSARSubmittedHandler(
	idempotent IdempotencyStore,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *DSARSubmittedHandler {
	return &DSARSubmittedHandler{
		idempotent: idempotent,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *DSARSubmittedHandler) Handle(ctx context.Context, evt event.DSARSubmittedEvent) error {
	const handler = "DSARSubmittedHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	// Trigger LLM analysis for ACCESS/ERASURE requests
	// (LLM จะช่วยสร้างรายงานข้อมูลที่มี)
	if evt.RequestType == "ACCESS" || evt.RequestType == "ERASURE" {
		llmEvt := event.LLMAnalysisRequestedEvent{
			Metadata:     evt.Metadata,
			DSARID:       evt.DSARID,
			UserID:       evt.UserID,
			AnalysisType: "DATA_MAPPING",
			RequestedAt:  time.Now().UTC(),
		}
		payload, _ := json.Marshal(llmEvt)
		_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
			ID:            uuid.New(),
			AggregateType: "DSARRequest",
			AggregateID:   evt.DSARID,
			EventType:     event.TopicLLMAnalysisRequest,
			Topic:         event.TopicLLMAnalysisRequest,
			Payload:       payload,
			Status:        repository.OutboxStatusPending,
			CreatedAt:     time.Now().UTC(),
		})
	}

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}
```

### 3.7.4 DataDeletionHandler

**`internal/modules/pdpa/application/event_handler/data_deletion_handler.go`**

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DataDeletionHandler handles DataDeletedEvent (post-deletion audit)
type DataDeletionHandler struct {
	idempotent IdempotencyStore
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewDataDeletionHandler สร้าง handler ใหม่
// NewDataDeletionHandler creates a new handler
func NewDataDeletionHandler(
	idempotent IdempotencyStore,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *DataDeletionHandler {
	return &DataDeletionHandler{
		idempotent: idempotent,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *DataDeletionHandler) Handle(ctx context.Context, evt event.DataDeletedEvent) error {
	const handler = "DataDeletionHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	// Record on blockchain
	blockchainEvt := event.BlockchainRecordEvent{
		Metadata:  evt.Metadata,
		UserID:    evt.UserID,
		Action:    "DATA_DELETED_" + evt.DeleteType,
		DataHash:  evt.Metadata.EventID.String(),
		CreatedAt: time.Now().UTC(),
	}
	payload, _ := json.Marshal(blockchainEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "UserAccountStatus",
		AggregateID:   evt.UserID,
		EventType:     event.TopicBlockchainRecord,
		Topic:         event.TopicBlockchainRecord,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 30*24*time.Hour)
}
```

### 3.7.5 AccountSuspendedHandler

**`internal/modules/pdpa/application/event_handler/account_suspended_handler.go`**

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// AccountSuspendedHandler handles AccountSuspendedEvent
type AccountSuspendedHandler struct {
	idempotent IdempotencyStore
	outboxRepo repository.OutboxRepository
	log        logger.Logger
}

// NewAccountSuspendedHandler สร้าง handler ใหม่
// NewAccountSuspendedHandler creates a new handler
func NewAccountSuspendedHandler(
	idempotent IdempotencyStore,
	outboxRepo repository.OutboxRepository,
	log logger.Logger,
) *AccountSuspendedHandler {
	return &AccountSuspendedHandler{
		idempotent: idempotent,
		outboxRepo: outboxRepo,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *AccountSuspendedHandler) Handle(ctx context.Context, evt event.AccountSuspendedEvent) error {
	const handler = "AccountSuspendedHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	// Send email informing user of retention policy
	emailEvt := event.EmailNotificationEvent{
		Metadata: evt.Metadata,
		Template: "account_suspended",
		Locale:   "th",
		Variables: map[string]interface{}{
			"retention_years":    evt.RetentionYears,
			"retention_deadline": evt.Deadline.Format(time.RFC3339),
		},
		RefType:   "UserAccountStatus",
		RefID:     evt.UserID.String(),
		CreatedAt: time.Now().UTC(),
	}
	payload, _ := json.Marshal(emailEvt)
	_ = h.outboxRepo.Save(ctx, &repository.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "UserAccountStatus",
		AggregateID:   evt.UserID,
		EventType:     event.TopicEmailNotification,
		Topic:         event.TopicEmailNotification,
		Payload:       payload,
		Status:        repository.OutboxStatusPending,
		CreatedAt:     time.Now().UTC(),
	})

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}
```

### 3.7.6 LLMAnalysisHandler

**`internal/modules/pdpa/application/event_handler/llm_analysis_handler.go`**

```go
package eventhandler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/pkg/logger"
)

// LLMClient interface สำหรับเรียก LLM
// LLMClient interface for calling the LLM service
type LLMClient interface {
	// AnalyzeDataWithContext ส่งข้อมูลให้ LLM วิเคราะห์
	// AnalyzeDataWithContext sends data to LLM for analysis
	AnalyzeDataWithContext(ctx context.Context, analysisType string, payload map[string]interface{}) (string, error)
}

// LLMAnalysisHandler handles LLMAnalysisRequestedEvent (async worker)
type LLMAnalysisHandler struct {
	llm        LLMClient
	idempotent IdempotencyStore
	log        logger.Logger
}

// NewLLMAnalysisHandler สร้าง handler ใหม่
// NewLLMAnalysisHandler creates a new handler
func NewLLMAnalysisHandler(
	llm LLMClient,
	idempotent IdempotencyStore,
	log logger.Logger,
) *LLMAnalysisHandler {
	return &LLMAnalysisHandler{
		llm:        llm,
		idempotent: idempotent,
		log:        log,
	}
}

// Handle ประมวลผล (เรียก LLM แบบ async พร้อม circuit breaker)
// Handle processes (async LLM call with circuit breaker in client)
func (h *LLMAnalysisHandler) Handle(ctx context.Context, evt event.LLMAnalysisRequestedEvent) error {
	const handler = "LLMAnalysisHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	// Call LLM with timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	result, err := h.llm.AnalyzeDataWithContext(ctx, evt.AnalysisType, map[string]interface{}{
		"dsar_id": evt.DSARID.String(),
		"user_id": evt.UserID.String(),
	})
	if err != nil {
		h.log.Error("LLM analysis failed", "error", err, "dsar_id", evt.DSARID)
		// Do NOT mark as processed so it can be retried
		return err
	}
	h.log.Info("LLM analysis completed", "dsar_id", evt.DSARID, "summary_len", len(result))

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}

// Import guard (compile check)
var _ = uuid.Nil
```

### 3.7.7 EmailNotificationHandler

**`internal/modules/pdpa/application/event_handler/email_notification_handler.go`**

```go
package eventhandler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/pkg/logger"
)

// EmailClient interface สำหรับส่ง email
// EmailClient interface for sending email
type EmailClient interface {
	// Send ส่ง email ด้วย template
	// Send sends an email using a template
	Send(ctx context.Context, to, subject, template string, locale string, vars map[string]interface{}) error
}

// EmailNotificationHandler handles EmailNotificationEvent (async worker)
type EmailNotificationHandler struct {
	email      EmailClient
	idempotent IdempotencyStore
	log        logger.Logger
}

// NewEmailNotificationHandler สร้าง handler ใหม่
// NewEmailNotificationHandler creates a new handler
func NewEmailNotificationHandler(
	email EmailClient,
	idempotent IdempotencyStore,
	log logger.Logger,
) *EmailNotificationHandler {
	return &EmailNotificationHandler{
		email:      email,
		idempotent: idempotent,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *EmailNotificationHandler) Handle(ctx context.Context, evt event.EmailNotificationEvent) error {
	const handler = "EmailNotificationHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	subject := evt.Subject
	if subject == "" {
		subject = evt.Template
	}
	if err := h.email.Send(ctx, evt.To, subject, evt.Template, evt.Locale, evt.Variables); err != nil {
		h.log.Error("email send failed", "error", err, "template", evt.Template)
		return err // retry
	}

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 7*24*time.Hour)
}
```

### 3.7.8 BlockchainRecordHandler

**`internal/modules/pdpa/application/event_handler/blockchain_record_handler.go`**

```go
package eventhandler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/pkg/logger"
)

// BlockchainClient interface สำหรับบันทึกบน blockchain
// BlockchainClient interface for recording on blockchain
type BlockchainClient interface {
	// Record บันทึกข้อมูลลง blockchain
	// Record writes data to blockchain
	Record(ctx context.Context, userID, action, dataHash string) (string, error)
}

// BlockchainRecordHandler handles BlockchainRecordEvent (async worker)
type BlockchainRecordHandler struct {
	client     BlockchainClient
	idempotent IdempotencyStore
	log        logger.Logger
}

// NewBlockchainRecordHandler สร้าง handler ใหม่
// NewBlockchainRecordHandler creates a new handler
func NewBlockchainRecordHandler(
	client BlockchainClient,
	idempotent IdempotencyStore,
	log logger.Logger,
) *BlockchainRecordHandler {
	return &BlockchainRecordHandler{
		client:     client,
		idempotent: idempotent,
		log:        log,
	}
}

// Handle ประมวลผล
// Handle processes
func (h *BlockchainRecordHandler) Handle(ctx context.Context, evt event.BlockchainRecordEvent) error {
	const handler = "BlockchainRecordHandler"

	processed, err := h.idempotent.IsProcessed(ctx, evt.Metadata.EventID, handler)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	txHash, err := h.client.Record(ctx, evt.UserID.String(), evt.Action, evt.DataHash)
	if err != nil {
		h.log.Error("blockchain record failed", "error", err, "action", evt.Action)
		return err // retry
	}
	h.log.Info("blockchain recorded", "tx_hash", txHash, "action", evt.Action)

	return h.idempotent.MarkProcessed(ctx, evt.Metadata.EventID, handler, 30*24*time.Hour)
}
```

---

## 3.8 Infrastructure Layer — Postgres

### 3.8.1 GORM Models

**`internal/modules/pdpa/infrastructure/persistence/postgres/models.go`**

```go
package postgres

import (
	"time"

	"github.com/google/uuid"
)

// PdpAConsentModel GORM model สำหรับ pdpa_consents
// PdpAConsentModel is the GORM model for pdpa_consents
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
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

// TableName ตั้งชื่อตาราง
// TableName overrides the table name
func (PdpAConsentModel) TableName() string { return "pdpa_consents" }

// PdpADSARModel GORM model สำหรับ pdpa_dsar_requests
// PdpADSARModel is the GORM model for pdpa_dsar_requests
type PdpADSARModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	RequestType     string     `gorm:"type:varchar(30);not null;index"`
	Status          string     `gorm:"type:varchar(20);not null;index"`
	RequestedAt     time.Time  `gorm:"not null;index"`
	CompletedAt     *time.Time
	DataPayload     []byte     `gorm:"type:jsonb"`
	RejectionReason string     `gorm:"type:text"`
	OTPHash         string     `gorm:"type:varchar(128)"`
	OTPExpiredAt    time.Time
	OTPVerified     bool       `gorm:"default:false"`
	IPAddress       string     `gorm:"type:varchar(45)"`
	UserAgent       string     `gorm:"type:text"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
}

func (PdpADSARModel) TableName() string { return "pdpa_dsar_requests" }

// PdpAAccountStatusModel GORM model สำหรับ pdpa_user_account_statuses
// PdpAAccountStatusModel is the GORM model for pdpa_user_account_statuses
type PdpAAccountStatusModel struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	Status              string     `gorm:"type:varchar(20);not null;index"`
	SuspendedAt         *time.Time
	TerminatedAt        *time.Time
	DeletionConfirmedAt *time.Time
	RetentionDeadline   *time.Time `gorm:"index"`
	AutoDeletedAt       *time.Time
	CreatedAt           time.Time  `gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime"`
}

func (PdpAAccountStatusModel) TableName() string { return "pdpa_user_account_statuses" }

// PdpAAuditTrailModel GORM model สำหรับ pdpa_audit_trails
// PdpAAuditTrailModel is the GORM model for pdpa_audit_trails
type PdpAAuditTrailModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID `gorm:"type:uuid;index"`
	Action    string     `gorm:"type:varchar(60);not null;index"`
	Details   []byte     `gorm:"type:jsonb"`
	IPAddress string     `gorm:"type:varchar(45)"`
	UserAgent string     `gorm:"type:text"`
	CreatedAt time.Time  `gorm:"not null;index"`
}

func (PdpAAuditTrailModel) TableName() string { return "pdpa_audit_trails" }

// PdpAPolicyModel GORM model สำหรับ pdpa_privacy_policies
// PdpAPolicyModel is the GORM model for pdpa_privacy_policies
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

// PdpAOutboxModel GORM model สำหรับ pdpa_outbox
// PdpAOutboxModel is the GORM model for pdpa_outbox
type PdpAOutboxModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	AggregateType string     `gorm:"type:varchar(60);not null;index"`
	AggregateID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	EventType     string     `gorm:"type:varchar(120);not null;index"`
	Topic         string     `gorm:"type:varchar(120);not null"`
	Payload       []byte     `gorm:"type:jsonb;not null"`
	Metadata      []byte     `gorm:"type:jsonb"`
	Status        string     `gorm:"type:varchar(20);not null;index"`
	RetryCount    int        `gorm:"default:0"`
	LastError     string     `gorm:"type:text"`
	CreatedAt     time.Time  `gorm:"not null;index"`
	PublishedAt   *time.Time
}

func (PdpAOutboxModel) TableName() string { return "pdpa_outbox" }

// PdpAProcessedEventModel GORM model สำหรับ idempotency
// PdpAProcessedEventModel tracks processed events for idempotency
type PdpAProcessedEventModel struct {
	EventID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	HandlerName string    `gorm:"type:varchar(80);primaryKey"`
	ProcessedAt time.Time `gorm:"not null;index"`
	ExpiresAt   time.Time `gorm:"not null;index"`
}

func (PdpAProcessedEventModel) TableName() string { return "pdpa_processed_events" }
```

### 3.8.2 Consent Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/consent_repo_impl.go`**

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
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentRepositoryImpl การ implement ConsentRepository บน PostgreSQL
// ConsentRepositoryImpl implements ConsentRepository on PostgreSQL
type ConsentRepositoryImpl struct {
	db *gorm.DB
}

// NewConsentRepository สร้าง repository ใหม่
// NewConsentRepository creates a new repository
func NewConsentRepository(db *gorm.DB) *ConsentRepositoryImpl {
	return &ConsentRepositoryImpl{db: db}
}

// Save บันทึก consent ใหม่
// Save persists a new consent
func (r *ConsentRepositoryImpl) Save(ctx context.Context, log *entity.ConsentLog) error {
	m := toConsentModel(log)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update อัปเดต consent ที่มีอยู่
// Update updates an existing consent
func (r *ConsentRepositoryImpl) Update(ctx context.Context, log *entity.ConsentLog) error {
	m := toConsentModel(log)
	return r.db.WithContext(ctx).
		Model(&PdpAConsentModel{}).
		Where("id = ?", log.ID).
		Updates(map[string]interface{}{
			"status":          m.Status,
			"revoked_at":      m.RevokedAt,
			"auto_deleted_at": m.AutoDeletedAt,
			"ip_address":      m.IPAddress,
			"user_agent":      m.UserAgent,
		}).Error
}

// FindByID ค้นหาตาม ID
// FindByID finds by ID
func (r *ConsentRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrConsentNotFound
		}
		return nil, err
	}
	return toConsentEntity(&m), nil
}

// FindByUserID ค้นหาทั้งหมดของ user
// FindByUserID finds all consents for a user
func (r *ConsentRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error) {
	var ms []PdpAConsentModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("granted_at DESC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	return toConsentEntities(ms), nil
}

// FindActiveByUserAndPurpose ค้นหา active consent
// FindActiveByUserAndPurpose finds the active consent
func (r *ConsentRepositoryImpl) FindActiveByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose_code = ? AND status = ? AND expires_at > ?",
			userID, purpose.String(), vo.ConsentGranted, time.Now().UTC()).
		Order("granted_at DESC").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toConsentEntity(&m), nil
}

// FindLatestByUserAndPurpose ค้นหา consent ล่าสุด
// FindLatestByUserAndPurpose finds the latest consent
func (r *ConsentRepositoryImpl) FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (*entity.ConsentLog, error) {
	var m PdpAConsentModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose_code = ?", userID, purpose.String()).
		Order("granted_at DESC").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrConsentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toConsentEntity(&m), nil
}

// FindReadyForAutoDeletion ค้นหา consent ที่พร้อมลบ
// FindReadyForAutoDeletion finds consents ready for auto-deletion
func (r *ConsentRepositoryImpl) FindReadyForAutoDeletion(ctx context.Context, cutoff time.Time, limit int) ([]entity.ConsentLog, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAConsentModel
	err := r.db.WithContext(ctx).
		Where("status IN ? AND revoked_at IS NOT NULL AND revoked_at <= ?",
			[]string{vo.ConsentRevoked.String(), vo.ConsentExpired.String()}, cutoff).
		Limit(limit).
		Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return toConsentEntities(ms), nil
}

// DeleteByUserID ลบทั้งหมดของ user
// DeleteByUserID deletes all consents for a user
func (r *ConsentRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&PdpAConsentModel{}).Error
}

// DeleteByID ลบตาม ID
// DeleteByID deletes by ID
func (r *ConsentRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&PdpAConsentModel{}).Error
}

// IsConsentActive ตรวจสอบว่า consent active
// IsConsentActive checks if consent is active
func (r *ConsentRepositoryImpl) IsConsentActive(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&PdpAConsentModel{}).
		Where("user_id = ? AND purpose_code = ? AND status = ? AND expires_at > ?",
			userID, purpose.String(), vo.ConsentGranted, time.Now().UTC()).
		Count(&count).Error
	return count > 0, err
}

// CountActiveByPurpose นับจำนวน active consent แยกตาม purpose
// CountActiveByPurpose counts active consents grouped by purpose
func (r *ConsentRepositoryImpl) CountActiveByPurpose(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	type row struct {
		PurposeCode string
		Cnt         int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Model(&PdpAConsentModel{}).
		Select("purpose_code, COUNT(*) as cnt").
		Where("granted_at BETWEEN ? AND ? AND status = ?", from, to, vo.ConsentGranted).
		Group("purpose_code").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]int64, len(rows))
	for _, r := range rows {
		result[r.PurposeCode] = r.Cnt
	}
	return result, nil
}

// toConsentModel แปลง entity -> model
func toConsentModel(e *entity.ConsentLog) *PdpAConsentModel {
	return &PdpAConsentModel{
		ID:            e.ID,
		UserID:        e.UserID,
		SessionID:     e.SessionID,
		PurposeCode:   e.Purpose.String(),
		Status:        e.Status.String(),
		IPAddress:     e.IPAddress,
		UserAgent:     e.UserAgent,
		GrantedAt:     e.GrantedAt,
		ExpiresAt:     e.ExpiresAt,
		RevokedAt:     e.RevokedAt,
		AutoDeletedAt: e.AutoDeletedAt,
	}
}

// toConsentEntity แปลง model -> entity
func toConsentEntity(m *PdpAConsentModel) *entity.ConsentLog {
	return &entity.ConsentLog{
		ID:            m.ID,
		UserID:        m.UserID,
		SessionID:     m.SessionID,
		Purpose:       vo.ConsentPurpose(m.PurposeCode),
		Status:        vo.ConsentStatus(m.Status),
		IPAddress:     m.IPAddress,
		UserAgent:     m.UserAgent,
		GrantedAt:     m.GrantedAt,
		ExpiresAt:     m.ExpiresAt,
		RevokedAt:     m.RevokedAt,
		AutoDeletedAt: m.AutoDeletedAt,
	}
}

// toConsentEntities แปลง slice
func toConsentEntities(ms []PdpAConsentModel) []entity.ConsentLog {
	out := make([]entity.ConsentLog, len(ms))
	for i := range ms {
		out[i] = *toConsentEntity(&ms[i])
	}
	return out
}
```

### 3.8.3 DSAR Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/dsar_repo_impl.go`**

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
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// DSARRepositoryImpl implements DSARRepository on PostgreSQL
type DSARRepositoryImpl struct {
	db *gorm.DB
}

// NewDSARRepository สร้าง repository ใหม่
// NewDSARRepository creates a new repository
func NewDSARRepository(db *gorm.DB) *DSARRepositoryImpl {
	return &DSARRepositoryImpl{db: db}
}

// Save บันทึก DSAR ใหม่
func (r *DSARRepositoryImpl) Save(ctx context.Context, req *entity.DSARRequest) error {
	return r.db.WithContext(ctx).Create(toDSARModel(req)).Error
}

// Update อัปเดต DSAR
func (r *DSARRepositoryImpl) Update(ctx context.Context, req *entity.DSARRequest) error {
	m := toDSARModel(req)
	return r.db.WithContext(ctx).
		Model(&PdpADSARModel{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status":           m.Status,
			"completed_at":     m.CompletedAt,
			"data_payload":     m.DataPayload,
			"rejection_reason": m.RejectionReason,
			"otp_verified":     m.OTPVerified,
		}).Error
}

// FindByID ค้นหาตาม ID
func (r *DSARRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error) {
	var m PdpADSARModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrDSARNotFound
		}
		return nil, err
	}
	return toDSAREntity(&m), nil
}

// FindByUserID ค้นหาทั้งหมดของ user
func (r *DSARRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error) {
	var ms []PdpADSARModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("requested_at DESC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	return toDSAREntities(ms), nil
}

// FindByStatus ค้นหาตาม status
func (r *DSARRepositoryImpl) FindByStatus(ctx context.Context, status vo.DSARStatus, limit, offset int) ([]entity.DSARRequest, error) {
	if limit <= 0 {
		limit = 50
	}
	var ms []PdpADSARModel
	err := r.db.WithContext(ctx).
		Where("status = ?", status.String()).
		Order("requested_at ASC").
		Limit(limit).Offset(offset).
		Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return toDSAREntities(ms), nil
}

// CountByUserAndTypeInPeriod นับจำนวน
func (r *DSARRepositoryImpl) CountByUserAndTypeInPeriod(ctx context.Context, userID uuid.UUID, requestType vo.DSARType, from time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&PdpADSARModel{}).
		Where("user_id = ? AND request_type = ? AND requested_at >= ?", userID, requestType.String(), from).
		Count(&count).Error
	return count, err
}

// DeleteByUserID ลบทั้งหมดของ user
func (r *DSARRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PdpADSARModel{}).Error
}

func toDSARModel(e *entity.DSARRequest) *PdpADSARModel {
	return &PdpADSARModel{
		ID:              e.ID,
		UserID:          e.UserID,
		RequestType:     e.RequestType.String(),
		Status:          e.Status.String(),
		RequestedAt:     e.RequestedAt,
		CompletedAt:     e.CompletedAt,
		DataPayload:     e.DataPayload,
		RejectionReason: e.RejectionReason,
		OTPHash:         e.OTPHash,
		OTPExpiredAt:    e.OTPExpiredAt,
		OTPVerified:     e.OTPVerified,
		IPAddress:       e.IPAddress,
		UserAgent:       e.UserAgent,
	}
}

func toDSAREntity(m *PdpADSARModel) *entity.DSARRequest {
	return &entity.DSARRequest{
		ID:              m.ID,
		UserID:          m.UserID,
		RequestType:     vo.DSARType(m.RequestType),
		Status:          vo.DSARStatus(m.Status),
		RequestedAt:     m.RequestedAt,
		CompletedAt:     m.CompletedAt,
		DataPayload:     m.DataPayload,
		RejectionReason: m.RejectionReason,
		OTPHash:         m.OTPHash,
		OTPExpiredAt:    m.OTPExpiredAt,
		OTPVerified:     m.OTPVerified,
		IPAddress:       m.IPAddress,
		UserAgent:       m.UserAgent,
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

### 3.8.4 Account Status Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/account_status_repo_impl.go`**

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
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// AccountStatusRepositoryImpl implements UserAccountStatusRepository
type AccountStatusRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountStatusRepository สร้าง repository ใหม่
func NewAccountStatusRepository(db *gorm.DB) *AccountStatusRepositoryImpl {
	return &AccountStatusRepositoryImpl{db: db}
}

// Save upsert สถานะบัญชี
func (r *AccountStatusRepositoryImpl) Save(ctx context.Context, st *entity.UserAccountStatus) error {
	m := toAccountModel(st)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"status", "suspended_at", "terminated_at",
				"deletion_confirmed_at", "retention_deadline",
				"auto_deleted_at", "updated_at",
			}),
		}).
		Create(m).Error
}

// FindByUserID ค้นหาตาม user
func (r *AccountStatusRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error) {
	var m PdpAAccountStatusModel
	if err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrAccountNotFound
		}
		return nil, err
	}
	return toAccountEntity(&m), nil
}

// FindReadyForAutoDeletion ค้นหาบัญชีที่พร้อมลบ
func (r *AccountStatusRepositoryImpl) FindReadyForAutoDeletion(ctx context.Context, now time.Time, limit int) ([]entity.UserAccountStatus, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAAccountStatusModel
	err := r.db.WithContext(ctx).
		Where("status = ? AND retention_deadline IS NOT NULL AND retention_deadline <= ?",
			vo.AccountSuspended.String(), now).
		Limit(limit).
		Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.UserAccountStatus, len(ms))
	for i := range ms {
		out[i] = *toAccountEntity(&ms[i])
	}
	return out, nil
}

// ExistsByUserID ตรวจสอบว่ามีหรือยัง
func (r *AccountStatusRepositoryImpl) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&PdpAAccountStatusModel{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count > 0, err
}

func toAccountModel(e *entity.UserAccountStatus) *PdpAAccountStatusModel {
	return &PdpAAccountStatusModel{
		ID:                  e.ID,
		UserID:              e.UserID,
		Status:              e.Status.String(),
		SuspendedAt:         e.SuspendedAt,
		TerminatedAt:        e.TerminatedAt,
		DeletionConfirmedAt: e.DeletionConfirmedAt,
		RetentionDeadline:   e.RetentionDeadline,
		AutoDeletedAt:       e.AutoDeletedAt,
	}
}

func toAccountEntity(m *PdpAAccountStatusModel) *entity.UserAccountStatus {
	return &entity.UserAccountStatus{
		ID:                  m.ID,
		UserID:              m.UserID,
		Status:              vo.AccountStatus(m.Status),
		SuspendedAt:         m.SuspendedAt,
		TerminatedAt:        m.TerminatedAt,
		DeletionConfirmedAt: m.DeletionConfirmedAt,
		RetentionDeadline:   m.RetentionDeadline,
		AutoDeletedAt:       m.AutoDeletedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}
```

### 3.8.5 Audit Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/audit_repo_impl.go`**

```go
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// AuditRepositoryImpl implements AuditRepository
type AuditRepositoryImpl struct {
	db *gorm.DB
}

// NewAuditRepository สร้าง repository ใหม่
func NewAuditRepository(db *gorm.DB) *AuditRepositoryImpl {
	return &AuditRepositoryImpl{db: db}
}

// Save บันทึก audit log
func (r *AuditRepositoryImpl) Save(ctx context.Context, a *entity.AuditTrail) error {
	details, _ := json.Marshal(a.Details)
	m := &PdpAAuditTrailModel{
		ID:        a.ID,
		UserID:    a.UserID,
		Action:    a.Action,
		Details:   details,
		IPAddress: a.IPAddress,
		UserAgent: a.UserAgent,
		CreatedAt: a.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// FindByFilter ค้นหาตามเงื่อนไข
func (r *AuditRepositoryImpl) FindByFilter(ctx context.Context, f repository.AuditFilter) ([]entity.AuditTrail, error) {
	q := r.db.WithContext(ctx).Model(&PdpAAuditTrailModel{})
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
			ID:        m.ID,
			UserID:    m.UserID,
			Action:    m.Action,
			Details:   details,
			IPAddress: m.IPAddress,
			UserAgent: m.UserAgent,
			CreatedAt: m.CreatedAt,
		}
	}
	return out, nil
}

// CountByAction นับตาม action
func (r *AuditRepositoryImpl) CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	type row struct {
		Action string
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Model(&PdpAAuditTrailModel{}).
		Select("action, COUNT(*) as cnt").
		Where("created_at BETWEEN ? AND ?", from, to).
		Group("action").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		out[r.Action] = r.Cnt
	}
	return out, nil
}

// unused import guard
var _ = uuid.Nil
```

### 3.8.6 Policy Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/policy_repo_impl.go`**

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

// PolicyRepositoryImpl implements PolicyRepository
type PolicyRepositoryImpl struct {
	db *gorm.DB
}

// NewPolicyRepository สร้าง repository ใหม่
func NewPolicyRepository(db *gorm.DB) *PolicyRepositoryImpl {
	return &PolicyRepositoryImpl{db: db}
}

// Save บันทึก
func (r *PolicyRepositoryImpl) Save(ctx context.Context, p *entity.PrivacyPolicy) error {
	m := &PdpAPolicyModel{
		ID:            p.ID,
		Version:       p.Version,
		Title:         p.Title,
		Content:       p.Content,
		EffectiveDate: p.EffectiveDate,
		IsActive:      p.IsActive,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// Update อัปเดต
func (r *PolicyRepositoryImpl) Update(ctx context.Context, p *entity.PrivacyPolicy) error {
	return r.db.WithContext(ctx).
		Model(&PdpAPolicyModel{}).
		Where("id = ?", p.ID).
		Updates(map[string]interface{}{
			"title":          p.Title,
			"content":        p.Content,
			"effective_date": p.EffectiveDate,
			"is_active":      p.IsActive,
		}).Error
}

// FindByID ค้นหาตาม ID
func (r *PolicyRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

// FindActive ค้นหาที่ active
func (r *PolicyRepositoryImpl) FindActive(ctx context.Context) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

// FindAll ค้นหาทั้งหมด
func (r *PolicyRepositoryImpl) FindAll(ctx context.Context) ([]entity.PrivacyPolicy, error) {
	var ms []PdpAPolicyModel
	if err := r.db.WithContext(ctx).Order("effective_date DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]entity.PrivacyPolicy, len(ms))
	for i := range ms {
		out[i] = *toPolicyEntity(&ms[i])
	}
	return out, nil
}

// FindByVersion ค้นหาตาม version
func (r *PolicyRepositoryImpl) FindByVersion(ctx context.Context, version string) (*entity.PrivacyPolicy, error) {
	var m PdpAPolicyModel
	if err := r.db.WithContext(ctx).Where("version = ?", version).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrPolicyNotFound
		}
		return nil, err
	}
	return toPolicyEntity(&m), nil
}

// DeactivateAll ปิดทั้งหมด
func (r *PolicyRepositoryImpl) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&PdpAPolicyModel{}).
		Where("is_active = ?", true).
		Update("is_active", false).Error
}

func toPolicyEntity(m *PdpAPolicyModel) *entity.PrivacyPolicy {
	return &entity.PrivacyPolicy{
		ID:            m.ID,
		Version:       m.Version,
		Title:         m.Title,
		Content:       m.Content,
		EffectiveDate: m.EffectiveDate,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
```

### 3.8.7 Outbox Repository Implementation

**`internal/modules/pdpa/infrastructure/persistence/postgres/outbox_repo_impl.go`**

```go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

// OutboxRepositoryImpl implements OutboxRepository
type OutboxRepositoryImpl struct {
	db *gorm.DB
}

// NewOutboxRepository สร้าง repository ใหม่
func NewOutboxRepository(db *gorm.DB) *OutboxRepositoryImpl {
	return &OutboxRepositoryImpl{db: db}
}

// Save บันทึก outbox event
func (r *OutboxRepositoryImpl) Save(ctx context.Context, e *repository.OutboxEvent) error {
	m := &PdpAOutboxModel{
		ID:            e.ID,
		AggregateType: e.AggregateType,
		AggregateID:   e.AggregateID,
		EventType:     e.EventType,
		Topic:         e.Topic,
		Payload:       e.Payload,
		Metadata:      e.Metadata,
		Status:        string(e.Status),
		RetryCount:    e.RetryCount,
		LastError:     e.LastError,
		CreatedAt:     e.CreatedAt,
		PublishedAt:   e.PublishedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// FetchPending ดึง events ที่รอ publish
// FetchPending fetches pending events (with row lock)
func (r *OutboxRepositoryImpl) FetchPending(ctx context.Context, limit int) ([]repository.OutboxEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	var ms []PdpAOutboxModel
	err := r.db.WithContext(ctx).
		Where("status = ?", string(repository.OutboxStatusPending)).
		Order("created_at ASC").
		Limit(limit).
		Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]repository.OutboxEvent, len(ms))
	for i := range ms {
		out[i] = repository.OutboxEvent{
			ID:            ms[i].ID,
			AggregateType: ms[i].AggregateType,
			AggregateID:   ms[i].AggregateID,
			EventType:     ms[i].EventType,
			Topic:         ms[i].Topic,
			Payload:       ms[i].Payload,
			Metadata:      ms[i].Metadata,
			Status:        repository.OutboxStatus(ms[i].Status),
			RetryCount:    ms[i].RetryCount,
			LastError:     ms[i].LastError,
			CreatedAt:     ms[i].CreatedAt,
			PublishedAt:   ms[i].PublishedAt,
		}
	}
	return out, nil
}

// MarkPublished ทำเครื่องหมายว่า publish สำเร็จ
func (r *OutboxRepositoryImpl) MarkPublished(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&PdpAOutboxModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       string(repository.OutboxStatusPublished),
			"published_at": now,
		}).Error
}

// MarkFailed ทำเครื่องหมายว่าล้มเหลว
func (r *OutboxRepositoryImpl) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return r.db.WithContext(ctx).
		Model(&PdpAOutboxModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      string(repository.OutboxStatusFailed),
			"last_error":  errMsg,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error
}

// CleanupPublished ลบ events เก่าที่ publish แล้ว
func (r *OutboxRepositoryImpl) CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("status = ? AND published_at < ?",
			string(repository.OutboxStatusPublished), cutoff).
		Delete(&PdpAOutboxModel{})
	return res.RowsAffected, res.Error
}
```

### 3.8.8 Processed Event Repository (Idempotency)

**`internal/modules/pdpa/infrastructure/persistence/postgres/processed_event_repo_impl.go`**

```go
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProcessedEventRepositoryImpl manages idempotency tracking
type ProcessedEventRepositoryImpl struct {
	db *gorm.DB
}

// NewProcessedEventRepository สร้าง repository ใหม่
func NewProcessedEventRepository(db *gorm.DB) *ProcessedEventRepositoryImpl {
	return &ProcessedEventRepositoryImpl{db: db}
}

// IsProcessed ตรวจสอบว่า event ถูกประมวลผลแล้วหรือยัง
// IsProcessed checks whether the event has been processed
func (r *ProcessedEventRepositoryImpl) IsProcessed(ctx context.Context, eventID uuid.UUID, handlerName string) (bool, error) {
	var m PdpAProcessedEventModel
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND handler_name = ?", eventID, handlerName).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// MarkProcessed บันทึกว่า event ถูกประมวลผลแล้ว
// MarkProcessed records the event as processed
func (r *ProcessedEventRepositoryImpl) MarkProcessed(ctx context.Context, eventID uuid.UUID, handlerName string, ttl time.Duration) error {
	m := &PdpAProcessedEventModel{
		EventID:     eventID,
		HandlerName: handlerName,
		ProcessedAt: time.Now().UTC(),
		ExpiresAt:   time.Now().UTC().Add(ttl),
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(m).Error
}

// CleanupExpired ลบ records ที่หมดอายุ
// CleanupExpired removes expired records
func (r *ProcessedEventRepositoryImpl) CleanupExpired(ctx context.Context) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&PdpAProcessedEventModel{})
	return res.RowsAffected, res.Error
}
```

---

## 3.9 Infrastructure — Cache (Redis)

**`internal/modules/pdpa/infrastructure/cache/redis/consent_cache.go`**

```go
package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
)

// ConsentCacheImpl implements ConsentCache with Redis
type ConsentCacheImpl struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

// NewConsentCache สร้าง cache ใหม่
// NewConsentCache creates a new consent cache
func NewConsentCache(client *redis.Client, ttl time.Duration) *ConsentCacheImpl {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ConsentCacheImpl{
		client: client,
		ttl:    ttl,
		prefix: "pdpa",
	}
}

// key pattern: pdpa:consent:{userID}:{purpose}
func (c *ConsentCacheImpl) key(userID uuid.UUID, purpose vo.ConsentPurpose) string {
	return fmt.Sprintf("%s:consent:%s:%s", c.prefix, userID.String(), purpose.String())
}

// Set บันทึก status
func (c *ConsentCacheImpl) Set(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose, status vo.ConsentStatus) error {
	return c.client.Set(ctx, c.key(userID, purpose), status.String(), c.ttl).Err()
}

// Get อ่าน status
func (c *ConsentCacheImpl) Get(ctx context.Context, userID uuid.UUID, purpose vo.ConsentPurpose) (vo.ConsentStatus, error) {
	v, err := c.client.Get(ctx, c.key(userID, purpose)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return vo.ConsentStatus(v), nil
}

// DeleteAllByUser ลบ key ทั้งหมดของ user
// DeleteAllByUser deletes all keys for a user
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

## 3.10 Infrastructure — Kafka Producer

**`internal/modules/pdpa/infrastructure/messaging/kafka/producer.go`**

```go
package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"icmongolang/pkg/logger"
)

// Producer abstraction
// Producer is a Kafka producer abstraction
type Producer interface {
	Publish(ctx context.Context, topic string, key string, payload []byte) error
	Close() error
}

// SyncProducer wraps sarama.SyncProducer
type SyncProducer struct {
	producer sarama.SyncProducer
	log      logger.Logger
}

// NewSyncProducer สร้าง producer ใหม่
// NewSyncProducer creates a new producer
func NewSyncProducer(brokers []string, log logger.Logger) (*SyncProducer, error) {
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
	return &SyncProducer{producer: p, log: log}, nil
}

// Publish ส่งข้อความ
// Publish sends a message
func (p *SyncProducer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Error("kafka publish failed", "topic", topic, "error", err)
	}
	return err
}

// Close ปิด
func (p *SyncProducer) Close() error {
	return p.producer.Close()
}

// unused import guard
var _ = time.Now
```

---

## 3.11 Infrastructure — Outbox Publisher

**`internal/modules/pdpa/infrastructure/messaging/outbox/outbox_publisher.go`**

```go
package outbox

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging/kafka"
	"icmongolang/pkg/logger"
)

// OutboxPublisher publishes outbox events to Kafka in background
// OutboxPublisher เป็น background worker ที่คอย publish outbox ไป Kafka
type OutboxPublisher struct {
	outboxRepo repository.OutboxRepository
	producer   kafka.Producer
	log        logger.Logger
	interval   time.Duration
	batchSize  int
}

// NewOutboxPublisher สร้าง publisher ใหม่
// NewOutboxPublisher creates a new publisher
func NewOutboxPublisher(
	outboxRepo repository.OutboxRepository,
	producer kafka.Producer,
	log logger.Logger,
) *OutboxPublisher {
	return &OutboxPublisher{
		outboxRepo: outboxRepo,
		producer:   producer,
		log:        log,
		interval:   500 * time.Millisecond,
		batchSize:  100,
	}
}

// Run เริ่ม loop
// Run starts the publishing loop
func (p *OutboxPublisher) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			p.log.Info("outbox publisher stopping")
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxPublisher) processBatch(ctx context.Context) {
	events, err := p.outboxRepo.FetchPending(ctx, p.batchSize)
	if err != nil {
		p.log.Error("fetch outbox failed", "error", err)
		return
	}
	for i := range events {
		e := &events[i]
		if err := p.producer.Publish(ctx, e.Topic, e.AggregateID.String(), e.Payload); err != nil {
			_ = p.outboxRepo.MarkFailed(ctx, e.ID, err.Error())
			continue
		}
		if err := p.outboxRepo.MarkPublished(ctx, e.ID); err != nil {
			p.log.Warn("mark published failed", "id", e.ID, "error", err)
		}
	}
}

// Cleanup ลบ events เก่า
// Cleanup removes old events
func (p *OutboxPublisher) Cleanup(ctx context.Context) {
	cutoff := time.Now().UTC().Add(-7 * 24 * time.Hour)
	n, err := p.outboxRepo.CleanupPublished(ctx, cutoff)
	if err != nil {
		p.log.Error("outbox cleanup failed", "error", err)
		return
	}
	p.log.Info("outbox cleanup completed", "deleted", n)
}
```

---

## 3.12 Infrastructure — External Clients

**`internal/modules/pdpa/infrastructure/external/email_client.go`**

```go
package external

import (
	"context"
	"fmt"

	"icmongolang/pkg/logger"
)

// SMTPEmailClient implements EmailClient
// (ในระบบจริงสามารถใช้ SMTP, SendGrid, Amazon SES ฯลฯ)
// SMTPEmailClient is a simple email client placeholder
type SMTPEmailClient struct {
	from string
	log  logger.Logger
}

// NewEmailClient สร้าง email client ใหม่
// NewEmailClient creates a new email client
func NewEmailClient(from string, log logger.Logger) *SMTPEmailClient {
	return &SMTPEmailClient{from: from, log: log}
}

// Send ส่ง email
// Send sends an email
func (c *SMTPEmailClient) Send(ctx context.Context, to, subject, template, locale string, vars map[string]interface{}) error {
	// TODO: render template ตาม locale
	// TODO: render template based on locale
	// TODO: send via SMTP provider
	c.log.Info("sending email",
		"to", to,
		"subject", subject,
		"template", template,
		"locale", locale,
	)
	return nil
}

// CircuitBreakerConfig ตัวอย่าง config
// CircuitBreakerConfig example
type CircuitBreakerConfig struct {
	MaxRequests   uint32
	Interval      int // seconds
	Timeout       int // seconds
	ReadyToTrip   func(total, failed uint32) bool
}

// DefaultCircuitBreakerConfig ค่าเริ่มต้น
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxRequests: 5,
		Interval:    60,
		Timeout:     30,
		ReadyToTrip: func(total, failed uint32) bool {
			return total >= 5 && float64(failed)/float64(total) >= 0.5
		},
	}
}

// unused import guard
var _ = fmt.Sprintf
```

**`internal/modules/pdpa/infrastructure/external/llm_client.go`**

```go
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"icmongolang/pkg/logger"
)

// HTTPLLMClient ตัวอย่าง LLM client (OpenAI-compatible)
// HTTPLLMClient is a sample LLM client (OpenAI-compatible)
type HTTPLLMClient struct {
	baseURL string
	apiKey  string
	model   string
	timeout time.Duration
	log     logger.Logger
}

// NewLLMClient สร้าง LLM client ใหม่
func NewLLMClient(baseURL, apiKey, model string, log logger.Logger) *HTTPLLMClient {
	return &HTTPLLMClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		timeout: 60 * time.Second,
		log:     log,
	}
}

// AnalyzeDataWithContext ส่งข้อมูลให้ LLM วิเคราะห์
// AnalyzeDataWithContext sends data to LLM for analysis
func (c *HTTPLLMClient) AnalyzeDataWithContext(ctx context.Context, analysisType string, payload map[string]interface{}) (string, error) {
	// TODO: replace with real HTTP call + circuit breaker (e.g., gobreaker)
	// TODO: แทนด้วย HTTP call จริง พร้อม circuit breaker
	body, _ := json.Marshal(payload)
	c.log.Info("LLM analysis requested",
		"type", analysisType,
		"payload_size", len(body),
		"model", c.model,
	)
	return fmt.Sprintf("[LLM-STUB] analysis=%s accepted", analysisType), nil
}
```

**`internal/modules/pdpa/infrastructure/external/blockchain_client.go`**

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

// MockBlockchainClient ตัวอย่าง blockchain client
// MockBlockchainClient is a sample blockchain client
type MockBlockchainClient struct {
	log logger.Logger
}

// NewBlockchainClient สร้าง client ใหม่
func NewBlockchainClient(log logger.Logger) *MockBlockchainClient {
	return &MockBlockchainClient{log: log}
}

// Record บันทึกข้อมูลลง blockchain
// Record writes data to blockchain
func (c *MockBlockchainClient) Record(ctx context.Context, userID, action, dataHash string) (string, error) {
	// TODO: replace with real blockchain (Hyperledger Fabric, Ethereum, etc.)
	// TODO: แทนด้วย blockchain จริง
	txInput := fmt.Sprintf("%s|%s|%s|%d", userID, action, dataHash, time.Now().UnixNano())
	h := sha256.Sum256([]byte(txInput))
	txHash := hex.EncodeToString(h[:])
	c.log.Info("blockchain record created",
		"tx_hash", txHash,
		"action", action,
		"user_id", userID,
	)
	return txHash, nil
}
```

---

## 3.13 Infrastructure — Scheduler

**`internal/modules/pdpa/infrastructure/scheduler/consent_cleanup_job.go`**

```go
package scheduler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/pkg/logger"
)

// ConsentCleanupJob รัน AutoDeleteExpired ทุกคืน
// ConsentCleanupJob runs AutoDeleteExpired nightly
type ConsentCleanupJob struct {
	handler *command.AutoDeleteExpiredHandler
	log     logger.Logger
}

// NewConsentCleanupJob สร้าง job ใหม่
func NewConsentCleanupJob(h *command.AutoDeleteExpiredHandler, log logger.Logger) *ConsentCleanupJob {
	return &ConsentCleanupJob{handler: h, log: log}
}

// Run เรียกใช้ handler
func (j *ConsentCleanupJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	res, err := j.handler.Handle(ctx, command.AutoDeleteExpiredCommand{
		BatchSize: 500,
		TraceID:   "cron-consent-cleanup",
	})
	if err != nil {
		j.log.Error("consent cleanup failed", "error", err)
		return
	}
	j.log.Info("consent cleanup completed",
		"accounts_deleted", res.AccountsDeleted,
		"consents_deleted", res.ConsentsDeleted,
		"errors", len(res.Errors),
	)
}
```

---

## 3.14 Interface Layer — DTOs

**`internal/modules/pdpa/interfaces/http/dto.go`**

```go
package http

import (
	"time"

	"github.com/google/uuid"
)

// RecordConsentRequest request body
type RecordConsentRequest struct {
	Purposes  map[string]bool `json:"purposes" validate:"required"`
	SessionID string          `json:"session_id,omitempty"`
}

// RecordConsentResponse response body
type RecordConsentResponse struct {
	Purposes []string    `json:"purposes"`
	IDs      []uuid.UUID `json:"ids"`
}

// RevokeConsentRequest request body
type RevokeConsentRequest struct {
	Purpose string `json:"purpose" validate:"required"`
}

// SubmitDSARRequest request body
type SubmitDSARRequest struct {
	RequestType string `json:"request_type" validate:"required"`
}

// SubmitDSARResponse response
type SubmitDSARResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
	Message     string    `json:"message"`
}

// ProcessDSARRequest request body
type ProcessDSARRequest struct {
	Action       string `json:"action" validate:"required"`
	OTPCode      string `json:"otp_code,omitempty"`
	RejectReason string `json:"reject_reason,omitempty"`
}

// ConfirmDeletionRequest request body
type ConfirmDeletionRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// AccountSuspendedRequest request body
type AccountSuspendedRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	SuspendedAt    time.Time `json:"suspended_at"`
	RetentionYears int       `json:"retention_years,omitempty"`
}

// AccountTerminatedRequest request body
type AccountTerminatedRequest struct {
	UserID       uuid.UUID `json:"user_id" validate:"required"`
	TerminatedAt time.Time `json:"terminated_at"`
}

// PublishPolicyRequest request body
type PublishPolicyRequest struct {
	Version       string    `json:"version" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Content       string    `json:"content" validate:"required"`
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
}

// ErrorResponse รูปแบบ error
type ErrorResponse struct {
	Error string `json:"error"`
}
```

---

## 3.15 Interface Layer — Handlers

**`internal/modules/pdpa/interfaces/http/consent_handler.go`**

```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/query"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/httputil"
)

// ConsentHandler HTTP handler สำหรับ consent
type ConsentHandler struct {
	record     *command.RecordConsentHandler
	revoke     *command.RevokeConsentHandler
	history    *query.GetConsentHistoryHandler
}

// NewConsentHandler สร้าง handler ใหม่
func NewConsentHandler(
	record *command.RecordConsentHandler,
	revoke *command.RevokeConsentHandler,
	history *query.GetConsentHistoryHandler,
) *ConsentHandler {
	return &ConsentHandler{record: record, revoke: revoke, history: history}
}

// RecordConsent POST /consent
func (h *ConsentHandler) RecordConsent(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req RecordConsentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	traceID := r.Header.Get("X-Trace-ID")
	result, err := h.record.Handle(r.Context(), command.RecordConsentCommand{
		UserID:    userID,
		SessionID: req.SessionID,
		Purposes:  req.Purposes,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		TraceID:   traceID,
	})
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, RecordConsentResponse{
		Purposes: result.RecordedPurposes,
		IDs:      result.ConsentIDs,
	})
}

// RevokeConsent DELETE /consent
func (h *ConsentHandler) RevokeConsent(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req RevokeConsentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	purpose, err := vo.NewConsentPurpose(req.Purpose)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	traceID := r.Header.Get("X-Trace-ID")
	if err := h.revoke.Handle(r.Context(), command.RevokeConsentCommand{
		UserID:    userID,
		Purpose:   purpose,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		TraceID:   traceID,
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "consent revoked"})
}

// GetHistory GET /consent/history
func (h *ConsentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.history.Handle(r.Context(), query.GetConsentHistoryQuery{UserID: userID})
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, items)
}

// extractUserID ดึง user ID จาก context ที่ middleware ใส่ไว้
// extractUserID retrieves user ID from context set by middleware
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

// unused import guard
var _ = chi.NewRouter
```

**`internal/modules/pdpa/interfaces/http/dsar_handler.go`**

```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/query"
	vo "icmongolang/internal/modules/pdpa/domain/value_object"
	"icmongolang/pkg/httputil"
)

// DSARHandler HTTP handler สำหรับ DSAR
type DSARHandler struct {
	submit  *command.SubmitDSARHandler
	process *command.ProcessDSARHandler
	getOne  *query.GetDSARStatusHandler
	list    *query.ListDSARHandler
}

// NewDSARHandler สร้าง handler ใหม่
func NewDSARHandler(
	submit *command.SubmitDSARHandler,
	process *command.ProcessDSARHandler,
	getOne *query.GetDSARStatusHandler,
	list *query.ListDSARHandler,
) *DSARHandler {
	return &DSARHandler{submit: submit, process: process, getOne: getOne, list: list}
}

// Submit POST /dsar
func (h *DSARHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req SubmitDSARRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	rt, err := vo.NewDSARType(req.RequestType)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.submit.Handle(r.Context(), command.SubmitDSARCommand{
		UserID:      userID,
		RequestType: rt,
		IPAddress:   r.RemoteAddr,
		UserAgent:   r.UserAgent(),
		TraceID:     r.Header.Get("X-Trace-ID"),
	})
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusCreated, SubmitDSARResponse{
		ID:          res.RequestID,
		Status:      res.Status,
		RequestedAt: res.RequestedAt,
		Message:     "DSAR submitted, please verify OTP sent to your email",
	})
}

// Process POST /dsar/{id}/process (admin)
func (h *DSARHandler) Process(w http.ResponseWriter, r *http.Request) {
	dsarID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid dsar id")
		return
	}
	var req ProcessDSARRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	operatorID, _ := extractUserID(r)
	if err := h.process.Handle(r.Context(), command.ProcessDSARCommand{
		DSARID:       dsarID,
		Action:       command.DSARAction(req.Action),
		OTPCode:      req.OTPCode,
		RejectReason: req.RejectReason,
		OperatorID:   operatorID,
		TraceID:      r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "processed"})
}

// Get GET /dsar/{id}
func (h *DSARHandler) Get(w http.ResponseWriter, r *http.Request) {
	dsarID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid dsar id")
		return
	}
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	res, err := h.getOne.Handle(r.Context(), query.GetDSARStatusQuery{DSARID: dsarID, UserID: userID})
	if err != nil {
		httputil.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	if res == nil {
		httputil.RespondError(w, http.StatusForbidden, "forbidden")
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

// List GET /dsar
func (h *DSARHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	if err != nil {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	res, err := h.list.Handle(r.Context(), query.ListDSARQuery{UserID: &userID})
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}
```

**`internal/modules/pdpa/interfaces/http/account_handler.go`**

```go
package http

import (
	"net/http"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/pkg/httputil"
)

// AccountHandler HTTP handler สำหรับ account
type AccountHandler struct {
	suspend  *command.HandleAccountSuspendedHandler
	terminate *command.HandleAccountTerminatedHandler
	confirm  *command.ConfirmDeletionHandler
	immediate *command.ImmediateDeletionHandler
}

// NewAccountHandler สร้าง handler ใหม่
func NewAccountHandler(
	suspend *command.HandleAccountSuspendedHandler,
	terminate *command.HandleAccountTerminatedHandler,
	confirm *command.ConfirmDeletionHandler,
	immediate *command.ImmediateDeletionHandler,
) *AccountHandler {
	return &AccountHandler{suspend: suspend, terminate: terminate, confirm: confirm, immediate: immediate}
}

// Suspend POST /account/suspend
func (h *AccountHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	var req AccountSuspendedRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.suspend.Handle(r.Context(), command.HandleAccountSuspendedCommand{
		UserID:         req.UserID,
		SuspendedAt:    req.SuspendedAt,
		RetentionYears: req.RetentionYears,
		TraceID:        r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "account suspended"})
}

// Terminate POST /account/terminate
func (h *AccountHandler) Terminate(w http.ResponseWriter, r *http.Request) {
	var req AccountTerminatedRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.terminate.Handle(r.Context(), command.HandleAccountTerminatedCommand{
		UserID:       req.UserID,
		TerminatedAt: req.TerminatedAt,
		TraceID:      r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "account terminated"})
}

// ConfirmDeletion POST /account/confirm-deletion
func (h *AccountHandler) ConfirmDeletion(w http.ResponseWriter, r *http.Request) {
	var req ConfirmDeletionRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	operatorID, _ := extractUserID(r)
	if err := h.confirm.Handle(r.Context(), command.ConfirmDeletionCommand{
		UserID:     req.UserID,
		OperatorID: operatorID,
		TraceID:    r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.immediate.Handle(r.Context(), command.ImmediateDeletionCommand{
		UserID:     req.UserID,
		Reason:     "termination_confirmed",
		OperatorID: operatorID,
		TraceID:    r.Header.Get("X-Trace-ID"),
	}); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]string{"message": "data deleted"})
}
```

**`internal/modules/pdpa/interfaces/http/policy_handler.go`**

```go
package http

import (
	"net/http"

	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httputil"
)

// PolicyHandler HTTP handler สำหรับ privacy policy
type PolicyHandler struct {
	publish *command.PublishPrivacyPolicyHandler
	getOne  *query.GetActivePrivacyPolicyHandler
	list    *query.ListPrivacyPoliciesHandler
}

// NewPolicyHandler สร้าง handler ใหม่
func NewPolicyHandler(
	publish *command.PublishPrivacyPolicyHandler,
	getOne *query.GetActivePrivacyPolicyHandler,
	list *query.ListPrivacyPoliciesHandler,
) *PolicyHandler {
	return &PolicyHandler{publish: publish, getOne: getOne, list: list}
}

// GetActive GET /policy
func (h *PolicyHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	res, err := h.getOne.Handle(r.Context(), query.GetActivePrivacyPolicyQuery{})
	if err != nil {
		httputil.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

// ListVersions GET /policy/versions
func (h *PolicyHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	res, err := h.list.Handle(r.Context(), query.ListPrivacyPoliciesQuery{})
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

// Publish POST /policy (admin)
func (h *PolicyHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req PublishPolicyRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	operatorID, _ := extractUserID(r)
	res, err := h.publish.Handle(r.Context(), command.PublishPrivacyPolicyCommand{
		Version:       req.Version,
		Title:         req.Title,
		Content:       req.Content,
		EffectiveDate: req.EffectiveDate,
		OperatorID:    operatorID,
		TraceID:       r.Header.Get("X-Trace-ID"),
	})
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusCreated, res)
}
```

**`internal/modules/pdpa/interfaces/http/admin_handler.go`**

```go
package http

import (
	"net/http"
	"time"

	"icmongolang/internal/modules/pdpa/application/query"
	"icmongolang/pkg/httputil"
)

// AdminHandler HTTP handler สำหรับ admin reports
type AdminHandler struct {
	report *query.GetAdminReportHandler
}

// NewAdminHandler สร้าง handler ใหม่
func NewAdminHandler(report *query.GetAdminReportHandler) *AdminHandler {
	return &AdminHandler{report: report}
}

// GetReport GET /admin/reports
func (h *AdminHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	from := parseTimeQuery(r.URL.Query().Get("from"), time.Now().AddDate(0, -1, 0))
	to := parseTimeQuery(r.URL.Query().Get("to"), time.Now())
	res, err := h.report.Handle(r.Context(), query.GetAdminReportQuery{From: from, To: to})
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.RespondJSON(w, http.StatusOK, res)
}

func parseTimeQuery(s string, def time.Time) time.Time {
	if s == "" {
		return def
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return def
}
```

---

## 3.16 Interface Layer — Routes

**`internal/modules/pdpa/interfaces/http/routes.go`**

```go
package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// errUnauthorized ใช้ภายใน handler helper
var errUnauthorized = errors.New("unauthorized")

// Handlers รวม handler ทั้งหมด
type Handlers struct {
	Consent *ConsentHandler
	DSAR    *DSARHandler
	Account *AccountHandler
	Policy  *PolicyHandler
	Admin   *AdminHandler
}

// AuthMiddleware interface เพื่อหลีกเลี่ยง cyclic import กับ auth module
type AuthMiddleware func(http.Handler) http.Handler

// RegisterRoutes ลงทะเบียน routes ทั้งหมด
// RegisterRoutes registers all PDPA routes
func RegisterRoutes(r chi.Router, h *Handlers, auth AuthMiddleware, adminOnly AuthMiddleware) {
	r.Route("/api/v1/pdpa", func(r chi.Router) {
		r.Use(auth)

		// Consent
		r.Post("/consent", h.Consent.RecordConsent)
		r.Delete("/consent", h.Consent.RevokeConsent)
		r.Get("/consent/history", h.Consent.GetHistory)

		// DSAR
		r.Post("/dsar", h.DSAR.Submit)
		r.Get("/dsar", h.DSAR.List)
		r.Get("/dsar/{id}", h.DSAR.Get)

		// Policy (public read)
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
		})
	})
}
```

---

## 3.17 Module Entry Point

**`internal/modules/pdpa/module.go`**

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
	redisinfra "icmongolang/internal/modules/pdpa/infrastructure/cache/redis"
	"icmongolang/internal/modules/pdpa/infrastructure/external"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging/kafka"
	"icmongolang/internal/modules/pdpa/infrastructure/messaging/outbox"
	pg "icmongolang/internal/modules/pdpa/infrastructure/persistence/postgres"
	"icmongolang/internal/modules/pdpa/infrastructure/scheduler"
	httpiface "icmongolang/internal/modules/pdpa/interfaces/http"
	"icmongolang/pkg/logger"
)

// Config ค่า config ของ module
// Config holds module configuration
type Config struct {
	KafkaBrokers    []string
	AnonymizationSalt string
	RetentionYears  int
	OutboxInterval  time.Duration
}

// Module รวม dependencies ทั้งหมด
// Module aggregates all module dependencies
type Module struct {
	// Repositories
	ConsentRepo    repository.ConsentRepository
	DSARRepo       repository.DSARRepository
	AccountRepo    repository.UserAccountStatusRepository
	AuditRepo      repository.AuditRepository
	PolicyRepo     repository.PolicyRepository
	OutboxRepo     repository.OutboxRepository
	ProcessedRepo  *pg.ProcessedEventRepositoryImpl

	// Services
	Validator    *service.ConsentValidator
	PolicySvc    *service.DeletionPolicyService
	Anonymizer   *service.AnonymizationService

	// Command handlers
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

	// Query handlers
	GetConsentHistory *query.GetConsentHistoryHandler
	GetDSARStatus     *query.GetDSARStatusHandler
	ListDSAR          *query.ListDSARHandler
	GetAdminReport    *query.GetAdminReportHandler
	GetActivePolicy   *query.GetActivePrivacyPolicyHandler
	ListPolicies      *query.ListPrivacyPoliciesHandler

	// Event handlers
	OnConsentGranted  *eventhandler.ConsentGrantedHandler
	OnConsentRevoked  *eventhandler.ConsentRevokedHandler
	OnDSARSubmitted   *eventhandler.DSARSubmittedHandler
	OnDataDeleted     *eventhandler.DataDeletionHandler
	OnAccountSuspended *eventhandler.AccountSuspendedHandler
	OnLLMAnalysis     *eventhandler.LLMAnalysisHandler
	OnEmailNotify     *eventhandler.EmailNotificationHandler
	OnBlockchainRec   *eventhandler.BlockchainRecordHandler

	// Infra
	Producer        kafka.Producer
	OutboxPublisher *outbox.OutboxPublisher
	ConsentCache    *redisinfra.ConsentCacheImpl
	CleanupJob      *scheduler.ConsentCleanupJob

	// HTTP
	Handlers *httpiface.Handlers

	log logger.Logger
}

// NewModule สร้าง module และ wire dependencies ทั้งหมด
// NewModule constructs the module and wires all dependencies
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
		cfg.AnonymizationSalt = "change-me-in-production"
	}

	// Repositories
	consentRepo := pg.NewConsentRepository(db)
	dsarRepo := pg.NewDSARRepository(db)
	accountRepo := pg.NewAccountStatusRepository(db)
	auditRepo := pg.NewAuditRepository(db)
	policyRepo := pg.NewPolicyRepository(db)
	outboxRepo := pg.NewOutboxRepository(db)
	processedRepo := pg.NewProcessedEventRepository(db)

	// Domain services
	validator := service.NewConsentValidator()
	policySvc := service.NewDeletionPolicyService(service.RetentionPolicy{
		SuspendedRetentionYears: cfg.RetentionYears,
		RevokedConsentDays:      365,
	})
	anonymizer := service.NewAnonymizationService(cfg.AnonymizationSalt)

	// Kafka producer
	producer, err := kafka.NewSyncProducer(cfg.KafkaBrokers, log)
	if err != nil {
		return nil, err
	}

	// Outbox publisher
	pub := outbox.NewOutboxPublisher(outboxRepo, producer, log)

	// Cache
	cache := redisinfra.NewConsentCache(redisClient, 24*time.Hour)

	// External clients
	emailClient := external.NewEmailClient("noreply@example.com", log)
	llmClient := external.NewLLMClient("https://api.example.com", "key", "gpt-4", log)
	blockchainClient := external.NewBlockchainClient(log)

	// Command handlers
	recordUC := command.NewRecordConsentHandler(consentRepo, accountRepo, auditRepo, outboxRepo, validator, log)
	revokeUC := command.NewRevokeConsentHandler(consentRepo, auditRepo, outboxRepo, validator, log)
	submitUC := command.NewSubmitDSARHandler(dsarRepo, accountRepo, auditRepo, outboxRepo, log)
	processUC := command.NewProcessDSARHandler(dsarRepo, auditRepo, outboxRepo, log)
	confirmUC := command.NewConfirmDeletionHandler(accountRepo, auditRepo, log)
	immediateUC := command.NewImmediateDeletionHandler(consentRepo, dsarRepo, accountRepo, auditRepo, outboxRepo, policySvc, log)
	autoDeleteUC := command.NewAutoDeleteExpiredHandler(consentRepo, dsarRepo, accountRepo, auditRepo, outboxRepo, policySvc, log)
	suspendUC := command.NewHandleAccountSuspendedHandler(accountRepo, auditRepo, outboxRepo, policySvc, log)
	terminateUC := command.NewHandleAccountTerminatedHandler(accountRepo, auditRepo, outboxRepo, log)
	publishUC := command.NewPublishPrivacyPolicyHandler(policyRepo, outboxRepo, log)

	// Query handlers
	historyQ := query.NewGetConsentHistoryHandler(consentRepo)
	dsarStatusQ := query.NewGetDSARStatusHandler(dsarRepo)
	listDSARQ := query.NewListDSARHandler(dsarRepo)
	reportQ := query.NewGetAdminReportHandler(consentRepo, auditRepo)
	activePolicyQ := query.NewGetActivePrivacyPolicyHandler(policyRepo)
	listPolicyQ := query.NewListPrivacyPoliciesHandler(policyRepo)

	// Event handlers
	onGranted := eventhandler.NewConsentGrantedHandler(cache, processedRepo, outboxRepo, log)
	onRevoked := eventhandler.NewConsentRevokedHandler(cache, processedRepo, outboxRepo, log)
	onDSARSub := eventhandler.NewDSARSubmittedHandler(processedRepo, outboxRepo, log)
	onDataDel := eventhandler.NewDataDeletionHandler(processedRepo, outboxRepo, log)
	onSusp := eventhandler.NewAccountSuspendedHandler(processedRepo, outboxRepo, log)
	onLLM := eventhandler.NewLLMAnalysisHandler(llmClient, processedRepo, log)
	onEmail := eventhandler.NewEmailNotificationHandler(emailClient, processedRepo, log)
	onBC := eventhandler.NewBlockchainRecordHandler(blockchainClient, processedRepo, log)

	// HTTP
	handlers := &httpiface.Handlers{
		Consent: httpiface.NewConsentHandler(recordUC, revokeUC, historyQ),
		DSAR:    httpiface.NewDSARHandler(submitUC, processUC, dsarStatusQ, listDSARQ),
		Account: httpiface.NewAccountHandler(suspendUC, terminateUC, confirmUC, immediateUC),
		Policy:  httpiface.NewPolicyHandler(publishUC, activePolicyQ, listPolicyQ),
		Admin:   httpiface.NewAdminHandler(reportQ),
	}

	m := &Module{
		ConsentRepo: consentRepo, DSARRepo: dsarRepo, AccountRepo: accountRepo,
		AuditRepo: auditRepo, PolicyRepo: policyRepo, OutboxRepo: outboxRepo,
		ProcessedRepo: processedRepo,

		Validator: validator, PolicySvc: policySvc, Anonymizer: anonymizer,

		RecordConsent: recordUC, RevokeConsent: revokeUC, SubmitDSAR: submitUC,
		ProcessDSAR: processUC, ConfirmDeletion: confirmUC, ImmediateDelete: immediateUC,
		AutoDelete: autoDeleteUC, SuspendAccount: suspendUC, TerminateAcct: terminateUC,
		PublishPolicy: publishUC,

		GetConsentHistory: historyQ, GetDSARStatus: dsarStatusQ, ListDSAR: listDSARQ,
		GetAdminReport: reportQ, GetActivePolicy: activePolicyQ, ListPolicies: listPolicyQ,

		OnConsentGranted: onGranted, OnConsentRevoked: onRevoked, OnDSARSubmitted: onDSARSub,
		OnDataDeleted: onDataDel, OnAccountSuspended: onSusp, OnLLMAnalysis: onLLM,
		OnEmailNotify: onEmail, OnBlockchainRec: onBC,

		Producer: producer, OutboxPublisher: pub, ConsentCache: cache,
		CleanupJob: scheduler.NewConsentCleanupJob(autoDeleteUC, log),

		Handlers: handlers,
		log:      log,
	}
	return m, nil
}

// RegisterHTTP ลงทะเบียน routes เข้า chi router
// RegisterHTTP registers routes onto a chi router
func (m *Module) RegisterHTTP(r chi.Router, auth httpiface.AuthMiddleware, adminOnly httpiface.AuthMiddleware) {
	httpiface.RegisterRoutes(r, m.Handlers, auth, adminOnly)
}

// StartBackground เริ่ม background workers (outbox publisher)
// StartBackground starts background workers
func (m *Module) StartBackground(ctx context.Context) {
	go m.OutboxPublisher.Run(ctx)
}

// Shutdown ปิด connections
// Shutdown closes connections
func (m *Module) Shutdown() error {
	return m.Producer.Close()
}
```

---

## สรุปส่วนที่ 2

ในส่วนนี้ได้สร้าง:

| Layer | จำนวนไฟล์ | รายละเอียด |
|-------|-----------|-----------|
| Repository Interfaces | 6 | consent, dsar, account, audit, policy, outbox |
| Domain Services | 3 | validator, deletion policy, anonymization |
| Commands | 10 (+ errors.go) | ทุก write use case |
| Queries | 4 (+ list DSAR) | read use cases |
| DTOs | 1 | application DTOs |
| Event Handlers | 8 | idempotent + DLQ ready |
| Postgres Repos | 8 | models + 7 repo impls |
| Redis Cache | 1 | consent cache |
| Kafka Producer | 1 | sync producer |
| Outbox Publisher | 1 | background publisher |
| External Clients | 3 | email, LLM, blockchain |
| Scheduler | 1 | cleanup cron |
| HTTP Layer | 5 handlers + routes + DTOs | consent, dsar, account, policy, admin |
| Module Entry | 1 | DI wiring |

**จุดเด่น:**
- **Transactional Outbox Pattern** — ทุก command เขียน event ลง outbox ใน transaction เดียวกัน, background publisher คอยส่ง Kafka
- **Idempotency** — event handlers เช็ค `processed_events` table ก่อนประมวลผล
- **Correlation/TraceID** — ฝังใน `EventMetadata` ทุก event
- **Rate limiting** — DSAR จำกัด 3 ครั้ง/วัน/user
- **OTP hashing** — เก็บเป็น SHA-256 ไม่ใช่ plaintext
- **Retention policy** — configurable ผ่าน `DeletionPolicyService`

---

**ส่วนถัดไป (ส่วนที่ 3):**
- **4. Database Migrations** (schema SQL ครบ + outbox + processed_events)
- **5. Docker Compose** (postgres, redis, kafka, zookeeper, elasticsearch, kafka-ui)
- **6. Testing Strategy** (unit + integration + e2e)
- **7. Deployment Considerations** (horizontal scaling, consumer groups, graceful shutdown, health checks)
- **8. Kafka Consumers** (worker สำหรับ event handlers)
- **9. `.env.example` + Makefile**

พิมพ์ **"ต่อ"** เพื่อไปส่วนที่ 3 ครับ