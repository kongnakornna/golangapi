# ปรับปรุง Event Handlers — Production-Grade

ผม review handlers ที่คุณส่งมาแล้ว **พบ 11 ปัญหา** ที่ต้องแก้ โดยเฉพาะเรื่อง **idempotency**, **structured events**, **PII leakage** และ **separation of concerns**

---

## 🔴 สรุปปัญหาที่พบใน Handlers

| # | Severity | ปัญหา | ผลกระทบ |
|---|----------|-------|---------|
| 1 | **P0** | ไม่มี **Idempotency check** | event ซ้ำ → audit trail ซ้ำ → blockchain duplicate |
| 2 | **P0** | `payloadToJSON` log **PII เต็ม** ลง log | log รั่วไหล (email/phone/OTP) |
| 3 | **P0** | Handler ทำงานแค่ save audit — **ไม่ได้ update cache / ส่ง email / เรียก blockchain** | business logic หาย |
| 4 | **P1** | `map[string]interface{}` → **ไม่มี type safety** | runtime error, ไม่มี schema |
| 5 | **P1** | `AuditTrailHandler.actions` map **ซ้ำซ้อน** (topic == action) | waste memory, สับสน |
| 6 | **P1** | `DLQHandler` **กลืน error** (`return nil` เสมอ) | ตรวจสอบไม่ได้ว่า audit fail |
| 7 | **P1** | ไม่มี **dispatcher/registry** | routing ต้องเขียนเองทุกที่ |
| 8 | **P1** | ไม่มี **metrics** | observability ไม่มี |
| 9 | **P2** | `baseHandler.saveAudit` **ไม่ redact PII** ก่อน save | PDPA violation (audit มี PII) |
| 10 | **P2** | ไม่มี **correlation ID** ใน logging | trace ไม่ได้ |
| 11 | **P2** | Module tree มี `repository/` **ที่ root** ซ้ำซ้อน | สับสน structure |

---

## 📁 โครงสร้างใหม่ที่แนะนำ

```
internal/modules/pdpa/
├── application/
│   ├── event_handler/                 ← dispatcher + registry
│   │   ├── dispatcher.go              [ใหม่]
│   │   ├── registry.go                [ใหม่]
│   │   ├── base_handler.go            [refactor]
│   │   └── ...
│   └── event_handler/handlers/        ← handler implementations
│       ├── consent_granted_handler.go
│       ├── consent_revoked_handler.go
│       ├── dsar_submitted_handler.go
│       ├── dsar_completed_handler.go
│       ├── account_suspended_handler.go
│       ├── account_terminated_handler.go
│       ├── data_deletion_requested_handler.go
│       └── dlq_handler.go
└── infrastructure/
    └── messaging/kafka/consumer/
        └── event_consumer.go          ← ใช้ dispatcher [ใหม่]
```

---

## 1. Structured Event Envelope (typed)

**`internal/modules/pdpa/application/event_handler/envelope.go`** (ใหม่)

```go
package eventhandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Envelope โครงสร้างมาตรฐานของ event ที่ consume จาก Kafka
// Envelope is the standard Kafka event structure consumed by handlers
type Envelope struct {
	EventType string          `json:"event_type"`
	Metadata  EventMetadata   `json:"metadata"`
	Payload   json.RawMessage `json:"payload"`
}

// EventMetadata metadata สำหรับ tracing + idempotency
// EventMetadata for tracing + idempotency
type EventMetadata struct {
	EventID       uuid.UUID `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id,omitempty"`
	TraceID       string    `json:"trace_id,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
	Version       int       `json:"version"`
	Source        string    `json:"source"`
}

// DecodeEnvelope parse envelope พร้อม validation ขั้นต้น
// DecodeEnvelope parses the envelope with basic validation
func DecodeEnvelope(data []byte) (*Envelope, error) {
	var e Envelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("decode envelope: %w", err)
	}
	if e.EventType == "" {
		return nil, fmt.Errorf("envelope missing event_type")
	}
	if e.Metadata.EventID == uuid.Nil {
		return nil, fmt.Errorf("envelope missing event_id")
	}
	return &e, nil
}

// DecodePayload helper แปลง payload → typed struct
// DecodePayload helper unmarshals payload into a typed struct
func DecodePayload[T any](e *Envelope) (*T, error) {
	var out T
	if err := json.Unmarshal(e.Payload, &out); err != nil {
		return nil, fmt.Errorf("decode payload for %s: %w", e.EventType, err)
	}
	return &out, nil
}

// Typed payloads — ใช้ schema ชัดเจน ไม่ใช้ map[string]any
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

type DSARSubmittedPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	RequestType string    `json:"request_type"`
	RequestedAt time.Time `json:"requested_at"`
	UserEmail   string    `json:"user_email,omitempty"`
}

type DSARCompletedPayload struct {
	DSARID      uuid.UUID `json:"dsar_id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	DataHash    string    `json:"data_hash,omitempty"`
	TxHash      string    `json:"tx_hash,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}

type AccountSuspendedPayload struct {
	UserID             uuid.UUID `json:"user_id"`
	Status             string    `json:"status"`
	SuspendedAt        time.Time `json:"suspended_at"`
	RetentionDeadline  time.Time `json:"retention_deadline"`
	RetentionYears     int       `json:"retention_years"`
}

type AccountTerminatedPayload struct {
	UserID       uuid.UUID `json:"user_id"`
	Status       string    `json:"status"`
	TerminatedAt time.Time `json:"terminated_at"`
}

type DataDeletionRequestedPayload struct {
	UserID           uuid.UUID `json:"user_id"`
	DeletionConfirmedAt time.Time `json:"deletion_confirmed_at"`
	ImmediateDeletion   bool    `json:"immediate_deletion"`
}
```

---

## 2. PII Redactor — ป้องกัน log รั่วไหล

**`internal/modules/pdpa/application/event_handler/pii_redactor.go`** (ใหม่)

```go
package eventhandler

import (
	"encoding/json"
	"regexp"
	"strings"
)

// PIIRedactor redact ข้อมูลอ่อนไหวก่อน log/save audit
// PIIRedactor redacts sensitive fields before logging or audit persistence
type PIIRedactor struct {
	maskKeys []string
	patterns []redactPattern
}

type redactPattern struct {
	re      *regexp.Regexp
	replace string
}

// NewPIIRedactor สร้าง redactor พร้อม rules เริ่มต้น
// NewPIIRedactor builds a redactor with default rules
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
			// Email
			{regexp.MustCompile(`[\w._%+\-]+@[\w.\-]+\.[A-Za-z]{2,}`), "***@***.***"},
			// Thai phone
			{regexp.MustCompile(`\b0\d{8,9}\b`), "0XXXXXXXXX"},
			// Thai national ID (13 digits)
			{regexp.MustCompile(`\b\d{13}\b`), "XXXXXXXXXXXXX"},
			// Credit card (16 digits with optional spaces/dashes)
			{regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`), "XXXX-XXXX-XXXX-XXXX"},
		},
	}
}

// RedactPayload คืน payload ใหม่ที่ redact แล้ว (deep copy)
// RedactPayload returns a deep-copied payload with PII redacted
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

// RedactJSONString คืน JSON string ของ payload ที่ redact แล้ว (สำหรับ logging)
// RedactJSONString returns the redacted JSON string of the payload (for logging)
func (r *PIIRedactor) RedactJSONString(payload interface{}) string {
	// ใช้ json.Marshal ก่อน แล้วค่อย redact string patterns
	b, err := json.Marshal(payload)
	if err != nil {
		return "<unmarshalable>"
	}
	return r.redactString(string(b))
}
```

---

## 3. Base Handler — Idempotency + Metrics + Redaction

**`internal/modules/pdpa/application/event_handler/base_handler.go`** (refactor)

```go
package eventhandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// IdempotencyStore abstraction — ทำงานข้าม replicas ผ่าน Redis
// IdempotencyStore works across replicas via Redis
type IdempotencyStore interface {
	// IsProcessed ตรวจสอบว่า event นี้เคยประมวลผลหรือยัง
	IsProcessed(ctx context.Context, eventID uuid.UUID, handlerName string) (bool, error)
	// MarkProcessed บันทึกว่า event นี้ประมวลผลแล้ว (atomic SET NX)
	MarkProcessed(ctx context.Context, eventID uuid.UUID, handlerName string, ttl time.Duration) error
}

// Metrics abstraction — implement ด้วย Prometheus ใน infra layer
// Metrics abstraction — implemented by Prometheus in the infra layer
type Metrics interface {
	IncHandlerProcessed(handler string, result string)
	ObserveHandlerDuration(handler string, seconds float64)
}

// nopMetrics ใช้เมื่อไม่ได้ inject metrics (test)
type nopMetrics struct{}

func (nopMetrics) IncHandlerProcessed(string, string)                {}
func (nopMetrics) ObserveHandlerDuration(string, float64)            {}

// NopMetrics คืน no-op metrics
func NopMetrics() Metrics { return nopMetrics{} }

// BaseHandler โครงสร้างร่วมของทุก handler
// BaseHandler is the shared base struct for all handlers
type BaseHandler struct {
	auditRepo   repository.AuditRepository
	idempotency IdempotencyStore
	metrics     Metrics
	redactor    *PIIRedactor
	logger      logger.Logger
	name        string
}

// NewBaseHandler สร้าง BaseHandler พร้อม dependencies
// NewBaseHandler creates a BaseHandler with all dependencies
func NewBaseHandler(
	name string,
	auditRepo repository.AuditRepository,
	idempotency IdempotencyStore,
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
		name:        name,
		auditRepo:   auditRepo,
		idempotency: idempotency,
		metrics:     metrics,
		redactor:    redactor,
		logger:      logger,
	}
}

// Name คืนชื่อ handler (ใช้เป็น key ของ idempotency)
// Name returns the handler name (used as idempotency key)
func (h *BaseHandler) Name() string { return h.name }

// IsProcessed ตรวจสอบ idempotency
// IsProcessed checks idempotency
func (h *BaseHandler) IsProcessed(ctx context.Context, eventID uuid.UUID) (bool, error) {
	if h.idempotency == nil {
		return false, nil // no idempotency configured
	}
	return h.idempotency.IsProcessed(ctx, eventID, h.name)
}

// MarkProcessed บันทึก idempotency
// MarkProcessed records idempotency
func (h *BaseHandler) MarkProcessed(ctx context.Context, eventID uuid.UUID, ttl time.Duration) error {
	if h.idempotency == nil {
		return nil
	}
	return h.idempotency.MarkProcessed(ctx, eventID, h.name, ttl)
}

// SaveAudit บันทึก audit trail (พร้อม redact PII)
// SaveAudit persists the audit trail (with PII redaction)
func (h *BaseHandler) SaveAudit(
	ctx context.Context,
	userID *uuid.UUID,
	action string,
	payload map[string]interface{},
) error {
	redacted := h.redactor.RedactPayload(payload)
	audit := entity.NewAuditTrail(userID, action, redacted)
	if err := h.auditRepo.Save(ctx, audit); err != nil {
		return fmt.Errorf("save audit for %s: %w", action, err)
	}
	return nil
}

// Observe วัด duration + log structured
// Observe measures duration and logs structured
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
		"handler", h.name,
		"event_id", meta.EventID,
		"correlation_id", meta.CorrelationID,
		"duration_ms", dur*1000,
		"result", result,
		"error", err,
	)
	return err
}

// ErrUnsupportedTopic topic ที่ handler ไม่รองรับ
var ErrUnsupportedTopic = errors.New("unsupported event topic")

// PayloadJSON helper สำหรับ debug log (จะ redact)
func (h *BaseHandler) PayloadJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "<unmarshalable>"
	}
	return h.redactor.RedactJSONString(string(b))
}
```

---

## 4. Typed Handlers — แต่ละตัวทำงานจริง

### 4.1 ConsentGrantedHandler — cache + email + blockchain

**`internal/modules/pdpa/application/event_handler/handlers/consent_granted_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// ConsentGrantedHandler handles consent.granted events:
//   1. Idempotency check
//   2. Update Redis cache
//   3. Enqueue email notification
//   4. Save audit trail (with PII redacted)
//
// ConsentGrantedHandler จัดการเหตุการณ์ consent.granted:
//   1. ตรวจ idempotency
//   2. อัปเดต Redis cache
//   3. เข้าคิวส่งอีเมล
//   4. บันทึก audit trail (redact PII)
type ConsentGrantedHandler struct {
	eventhandler.BaseHandler
	cache      repository.Cache
	producer   kafka.Producer
	emailTopic string
}

func NewConsentGrantedHandler(
	auditRepo repository.AuditRepository,
	cache repository.Cache,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *ConsentGrantedHandler {
	return &ConsentGrantedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"ConsentGrantedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		cache:      cache,
		producer:   producer,
		emailTopic: event.TopicEmailNotification,
	}
}

// Handle ประมวลผล envelope
func (h *ConsentGrantedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.ConsentGrantedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		// 1. Idempotency
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			h.logger.Infow("event already processed", "event_id", env.Metadata.EventID)
			return nil
		}

		// 2. Update cache (best-effort)
		purpose := valueobject.ConsentPurpose(payload.Purpose)
		if err := h.cache.Set(ctx, payload.UserID, purpose, valueobject.ConsentGranted); err != nil {
			h.logger.Warnw("cache set failed", "error", err, "user_id", payload.UserID)
		}

		// 3. Enqueue email (best-effort)
		if payload.UserID != [16]byte{} {
			_ = h.producer.Publish(ctx, h.emailTopic, payload.UserID.String(), []byte(`{
				"template":"consent_granted",
				"locale":"th",
				"user_id":"`+payload.UserID.String()+`",
				"purpose":"`+payload.Purpose+`"
			}`))
		}

		// 4. Save audit (PII redacted)
		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"consent_id":  payload.ConsentID.String(),
			"user_id":     payload.UserID.String(),
			"purpose":     payload.Purpose,
			"granted_at":  payload.GrantedAt,
			"expires_at":  payload.ExpiresAt,
			"ip_address":  payload.IPAddress,
		}); err != nil {
			return err
		}

		// 5. Mark processed
		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 4.2 ConsentRevokedHandler

**`internal/modules/pdpa/application/event_handler/handlers/consent_revoked_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type ConsentRevokedHandler struct {
	eventhandler.BaseHandler
	cache    repository.Cache
	producer kafka.Producer
}

func NewConsentRevokedHandler(
	auditRepo repository.AuditRepository,
	cache repository.Cache,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *ConsentRevokedHandler {
	return &ConsentRevokedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"ConsentRevokedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		cache:    cache,
		producer: producer,
	}
}

func (h *ConsentRevokedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.ConsentRevokedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		purpose := valueobject.ConsentPurpose(payload.Purpose)
		if err := h.cache.Set(ctx, payload.UserID, purpose, valueobject.ConsentRevoked); err != nil {
			h.logger.Warnw("cache set failed", "error", err)
		}

		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(),
			[]byte(`{"template":"consent_revoked","locale":"th","user_id":"`+payload.UserID.String()+`","purpose":"`+payload.Purpose+`"}`))

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

### 4.3 DSARSubmittedHandler — trigger LLM analysis

**`internal/modules/pdpa/application/event_handler/handlers/dsar_submitted_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type DSARSubmittedHandler struct {
	eventhandler.BaseHandler
	producer kafka.Producer
}

func NewDSARSubmittedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *DSARSubmittedHandler {
	return &DSARSubmittedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"DSARSubmittedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		producer: producer,
	}
}

func (h *DSARSubmittedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.DSARSubmittedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		// Trigger LLM analysis for ACCESS/ERASURE
		if payload.RequestType == "ACCESS" || payload.RequestType == "ERASURE" {
			_ = h.producer.Publish(ctx, "pdpa.llm.analysis.requested", payload.DSARID.String(),
				[]byte(`{"dsar_id":"`+payload.DSARID.String()+`","user_id":"`+payload.UserID.String()+`","analysis_type":"DATA_MAPPING"}`))
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

### 4.4 DSARCompletedHandler — email + blockchain record

**`internal/modules/pdpa/application/event_handler/handlers/dsar_completed_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type DSARCompletedHandler struct {
	eventhandler.BaseHandler
	producer kafka.Producer
}

func NewDSARCompletedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *DSARCompletedHandler {
	return &DSARCompletedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"DSARCompletedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		producer: producer,
	}
}

func (h *DSARCompletedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.DSARCompletedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		// Notify user via email
		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(),
			[]byte(`{"template":"dsar_completed","locale":"th","user_id":"`+payload.UserID.String()+`","dsar_id":"`+payload.DSARID.String()+`"}`))

		// Record on blockchain (immutable)
		_ = h.producer.Publish(ctx, event.TopicBlockchainRecord, payload.DSARID.String(),
			[]byte(`{"user_id":"`+payload.UserID.String()+`","action":"DSAR_COMPLETED","data_hash":"`+payload.DataHash+`"}`))

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"dsar_id":      payload.DSARID.String(),
			"user_id":      payload.UserID.String(),
			"status":       payload.Status,
			"data_hash":    payload.DataHash,
			"tx_hash":      payload.TxHash,
			"completed_at": payload.CompletedAt,
		}); err != nil {
			return err
		}

		return h.MarkProcessed(ctx, env.Metadata.EventID, 30*24*time.Hour)
	})
}
```

### 4.5 AccountSuspendedHandler

**`internal/modules/pdpa/application/event_handler/handlers/account_suspended_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type AccountSuspendedHandler struct {
	eventhandler.BaseHandler
	producer kafka.Producer
}

func NewAccountSuspendedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *AccountSuspendedHandler {
	return &AccountSuspendedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"AccountSuspendedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		producer: producer,
	}
}

func (h *AccountSuspendedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.AccountSuspendedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(),
			[]byte(`{"template":"account_suspended","locale":"th","user_id":"`+payload.UserID.String()+`","retention_years":`+itoa(payload.RetentionYears)+`}`))

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"user_id":            payload.UserID.String(),
			"status":             payload.Status,
			"suspended_at":       payload.SuspendedAt,
			"retention_deadline": payload.RetentionDeadline,
			"retention_years":    payload.RetentionYears,
		}); err != nil {
			return err
		}

		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
```

### 4.6 AccountTerminatedHandler

**`internal/modules/pdpa/application/event_handler/handlers/account_terminated_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type AccountTerminatedHandler struct {
	eventhandler.BaseHandler
	producer kafka.Producer
}

func NewAccountTerminatedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *AccountTerminatedHandler {
	return &AccountTerminatedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"AccountTerminatedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		producer: producer,
	}
}

func (h *AccountTerminatedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.AccountTerminatedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		_ = h.producer.Publish(ctx, event.TopicEmailNotification, payload.UserID.String(),
			[]byte(`{"template":"account_terminated","locale":"th","user_id":"`+payload.UserID.String()+`"}`))

		if err := h.SaveAudit(ctx, &payload.UserID, env.EventType, map[string]interface{}{
			"user_id":       payload.UserID.String(),
			"status":        payload.Status,
			"terminated_at": payload.TerminatedAt,
		}); err != nil {
			return err
		}

		return h.MarkProcessed(ctx, env.Metadata.EventID, 7*24*time.Hour)
	})
}
```

### 4.7 DataDeletionRequestedHandler — trigger deletion workflow

**`internal/modules/pdpa/application/event_handler/handlers/data_deletion_requested_handler.go`** (refactor)

```go
package handlers

import (
	"context"
	"time"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/event"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type DataDeletionRequestedHandler struct {
	eventhandler.BaseHandler
	producer kafka.Producer
}

func NewDataDeletionRequestedHandler(
	auditRepo repository.AuditRepository,
	producer kafka.Producer,
	idempotency eventhandler.IdempotencyStore,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *DataDeletionRequestedHandler {
	return &DataDeletionRequestedHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"DataDeletionRequestedHandler",
			auditRepo, idempotency, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
		producer: producer,
	}
}

func (h *DataDeletionRequestedHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	payload, err := eventhandler.DecodePayload[eventhandler.DataDeletionRequestedPayload](env)
	if err != nil {
		return err
	}

	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		done, err := h.IsProcessed(ctx, env.Metadata.EventID)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		// Trigger immediate deletion workflow if policy allows
		if payload.ImmediateDeletion {
			_ = h.producer.Publish(ctx, "pdpa.data.deleted", payload.UserID.String(),
				[]byte(`{"user_id":"`+payload.UserID.String()+`","delete_type":"IMMEDIATE"}`))
		}

		// Blockchain audit record
		_ = h.producer.Publish(ctx, event.TopicBlockchainRecord, payload.UserID.String(),
			[]byte(`{"user_id":"`+payload.UserID.String()+`","action":"DELETION_REQUESTED"}`))

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

### 4.8 DLQHandler — ไม่กลืน error

**`internal/modules/pdpa/application/event_handler/handlers/dlq_handler.go`** (refactor)

```go
package handlers

import (
	"context"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DLQHandler บันทึก dead-letter events ลง audit trail
// ไม่ใช้ idempotency (DLQ เป็น best-effort อยู่แล้ว)
//
// DLQHandler persists dead-letter events to audit trail.
// Skips idempotency since DLQ is best-effort.
type DLQHandler struct {
	eventhandler.BaseHandler
}

func NewDLQHandler(
	auditRepo repository.AuditRepository,
	metrics eventhandler.Metrics,
	logger logger.Logger,
) *DLQHandler {
	return &DLQHandler{
		BaseHandler: eventhandler.NewBaseHandler(
			"DLQHandler",
			auditRepo, nil /* no idempotency */, metrics,
			eventhandler.NewPIIRedactor(), logger,
		),
	}
}

// Handle บันทึก DLQ event — คืน error ถ้า save fail (ให้ consumer ตัดสินใจ)
// Handle persists the DLQ event — returns error if save fails (consumer decides)
func (h *DLQHandler) Handle(ctx context.Context, env *eventhandler.Envelope) error {
	return h.Observe(ctx, env.Metadata, func(ctx context.Context) error {
		// DLQ payload เก็บ envelope ที่ fail — record as-is (redacted)
		h.logger.Warnw("processing dead-letter event",
			"original_event_type", env.EventType,
			"event_id", env.Metadata.EventID,
			"correlation_id", env.Metadata.CorrelationID,
		)

		return h.SaveAudit(ctx, nil, "DLQ."+env.EventType, map[string]interface{}{
			"original_event_type": env.EventType,
			"original_event_id":   env.Metadata.EventID.String(),
			"correlation_id":      env.Metadata.CorrelationID,
			"payload":             string(env.Payload),
		})
	})
}
```

---

## 5. Dispatcher + Registry — Route topic → handler

**`internal/modules/pdpa/application/event_handler/registry.go`** (ใหม่)

```go
package eventhandler

import (
	"context"
	"fmt"
	"sync"
)

// Handler interface ที่ handler ทุกตัวต้อง implement
// Handler is implemented by every event handler
type Handler interface {
	Handle(ctx context.Context, env *Envelope) error
	Name() string
}

// Registry เก็บ mapping จาก topic → handler
// Registry holds topic → handler mappings
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// NewRegistry สร้าง registry ใหม่
// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

// Register ผูก topic กับ handler (แจ้ง error ถ้าซ้ำ)
// Register binds a topic to a handler (errors on duplicate)
func (r *Registry) Register(topic string, h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.handlers[topic]; exists {
		return fmt.Errorf("topic %q already registered", topic)
	}
	r.handlers[topic] = h
	return nil
}

// MustRegister helper สำหรับ bootstrap
// MustRegister is a bootstrap helper that panics on error
func (r *Registry) MustRegister(topic string, h Handler) {
	if err := r.Register(topic, h); err != nil {
		panic(err)
	}
}

// Lookup หา handler ของ topic
// Lookup finds the handler for a topic
func (r *Registry) Lookup(topic string) (Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[topic]
	return h, ok
}

// Topics คืนรายการ topics ที่ลงทะเบียนไว้
// Topics returns all registered topics
func (r *Registry) Topics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		out = append(out, t)
	}
	return out
}
```

**`internal/modules/pdpa/application/event_handler/dispatcher.go`** (ใหม่)

```go
package eventhandler

import (
	"context"
	"fmt"

	"icmongolang/pkg/logger"
)

// Dispatcher route message → handler ตาม topic
// Dispatcher routes a message to the appropriate handler based on its topic
type Dispatcher struct {
	registry *Registry
	dlq      Handler // fallback handler สำหรับ topic ที่ไม่รู้จัก
	logger   logger.Logger
}

// NewDispatcher สร้าง dispatcher ใหม่
// NewDispatcher creates a new dispatcher
func NewDispatcher(registry *Registry, dlq Handler, logger logger.Logger) *Dispatcher {
	return &Dispatcher{registry: registry, dlq: dlq, logger: logger}
}

// Dispatch ส่ง message ไปยัง handler ที่เหมาะสม
// Dispatch sends a message to the appropriate handler
func (d *Dispatcher) Dispatch(ctx context.Context, topic string, value []byte) error {
	env, err := DecodeEnvelope(value)
	if err != nil {
		// malformed envelope → DLQ
		d.logger.Errorw("decode envelope failed, routing to DLQ", "topic", topic, "error", err)
		if d.dlq == nil {
			return fmt.Errorf("decode failed (no DLQ): %w", err)
		}
		return d.dlq.Handle(ctx, &Envelope{
			EventType: "MALFORMED",
			Metadata:  EventMetadata{CorrelationID: "unknown"},
			Payload:   value,
		})
	}

	h, ok := d.registry.Lookup(env.EventType)
	if !ok {
		d.logger.Warnw("no handler registered, routing to DLQ",
			"event_type", env.EventType,
			"event_id", env.Metadata.EventID,
		)
		if d.dlq == nil {
			return fmt.Errorf("no handler for %q", env.EventType)
		}
		return d.dlq.Handle(ctx, env)
	}

	return h.Handle(ctx, env)
}

// RegisteredTopics คืน topic ทั้งหมดที่รองรับ
// RegisteredTopics returns all supported topics
func (d *Dispatcher) RegisteredTopics() []string {
	return d.registry.Topics()
}
```

---

## 6. Kafka Consumer Integration

**`internal/modules/pdpa/infrastructure/messaging/kafka/consumer/event_consumer.go`** (ใหม่)

```go
package consumer

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/pkg/logger"
)

// EventConsumer consume envelope-based events และ dispatch ให้ handler
// EventConsumer consumes envelope-based events and dispatches to handlers
type EventConsumer struct {
	cfg        ConsumerConfig
	dispatcher *eventhandler.Dispatcher
	logger     logger.Logger
}

func NewEventConsumer(
	brokers []string,
	groupID string,
	topics []string,
	dispatcher *eventhandler.Dispatcher,
	logger logger.Logger,
) *EventConsumer {
	return &EventConsumer{
		cfg: ConsumerConfig{
			Brokers:      brokers,
			GroupID:      groupID,
			Topics:       topics,
			MaxRetries:   3,
			RetryBackoff: 500 * time.Millisecond,
			DLQPrefix:    "pdpa.dlq.",
		},
		dispatcher: dispatcher,
		logger:     logger,
	}
}

// Run เริ่ม consume loop
// Run starts the consume loop
func (c *EventConsumer) Run(ctx context.Context) error {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_5_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(c.cfg.Brokers, c.cfg.GroupID, cfg)
	if err != nil {
		return err
	}
	defer group.Close()

	go func() {
		for err := range group.Errors() {
			c.logger.Errorw("consumer group error", "group", c.cfg.GroupID, "error", err)
		}
	}()

	handler := &dispatchHandler{
		dispatcher: c.dispatcher,
		logger:     c.logger,
	}

	for {
		if err := group.Consume(ctx, c.cfg.Topics, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// dispatchHandler implements sarama.ConsumerGroupHandler
type dispatchHandler struct {
	dispatcher *eventhandler.Dispatcher
	logger     logger.Logger
}

func (h *dispatchHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *dispatchHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *dispatchHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Use message key as correlation fallback
		if err := h.dispatcher.Dispatch(session.Context(), msg.Topic, msg.Value); err != nil {
			h.logger.Errorw("dispatch failed",
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
				"error", err,
			)
			// sarama จะ retry อีกครั้ง
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
```

---

## 7. Bootstrap — Registry + Dispatcher Wiring

**`internal/modules/pdpa/application/event_handler/wiring.go`** (ใหม่)

```go
package eventhandler

import (
	handlers "icmongolang/internal/modules/pdpa/application/event_handler/handlers"
	"icmongolang/internal/modules/pdpa/domain/event"
)

// RegistryConfig dependencies สำหรับ wiring
type RegistryConfig struct {
	ConsentGranted          *handlers.ConsentGrantedHandler
	ConsentRevoked          *handlers.ConsentRevokedHandler
	DSARSubmitted           *handlers.DSARSubmittedHandler
	DSARCompleted           *handlers.DSARCompletedHandler
	AccountSuspended        *handlers.AccountSuspendedHandler
	AccountTerminated       *handlers.AccountTerminatedHandler
	DataDeletionRequested   *handlers.DataDeletionRequestedHandler
	DLQ                     *handlers.DLQHandler
}

// BuildRegistry สร้าง registry พร้อม register handlers ทั้งหมด
// BuildRegistry builds a registry with all handlers registered
func BuildRegistry(cfg RegistryConfig) *Registry {
	r := NewRegistry()

	r.MustRegister(event.TopicConsentGranted, cfg.ConsentGranted)
	r.MustRegister(event.TopicConsentRevoked, cfg.ConsentRevoked)
	r.MustRegister(event.TopicDSARSubmitted, cfg.DSARSubmitted)
	r.MustRegister(event.TopicDSARCompleted, cfg.DSARCompleted)
	r.MustRegister(event.TopicAccountSuspended, cfg.AccountSuspended)
	r.MustRegister(event.TopicAccountTerminated, cfg.AccountTerminated)
	r.MustRegister(event.TopicDeletionRequested, cfg.DataDeletionRequested)

	return r
}

// DefaultTopics คืนรายการ topic ทั้งหมดที่ต้อง subscribe
// DefaultTopics returns all topics to subscribe to
func DefaultTopics() []string {
	return []string{
		event.TopicConsentGranted,
		event.TopicConsentRevoked,
		event.TopicDSARSubmitted,
		event.TopicDSARCompleted,
		event.TopicAccountSuspended,
		event.TopicAccountTerminated,
		event.TopicDeletionRequested,
	}
}
```

---

## 8. Idempotency Store — Redis Implementation

**`internal/modules/pdpa/infrastructure/persistence/redis/idempotency_store.go`** (ใหม่)

```go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// IdempotencyStore Redis-backed idempotency (SET NX)
type IdempotencyStore struct {
	client *redis.Client
	prefix string
}

func NewIdempotencyStore(client *redis.Client) *IdempotencyStore {
	return &IdempotencyStore{
		client: client,
		prefix: "pdpa:idem:",
	}
}

func (s *IdempotencyStore) IsProcessed(ctx context.Context, eventID uuid.UUID, handlerName string) (bool, error) {
	n, err := s.client.Exists(ctx, s.key(eventID, handlerName)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// MarkProcessed ใช้ SET NX — atomic across replicas
func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID uuid.UUID, handlerName string, ttl time.Duration) error {
	// SetNX returns true if set, false if already exists
	_, err := s.client.SetNX(ctx, s.key(eventID, handlerName), "1", ttl).Result()
	return err
}

func (s *IdempotencyStore) key(eventID uuid.UUID, handlerName string) string {
	return fmt.Sprintf("%s%s:%s", s.prefix, handlerName, eventID.String())
}
```

---

## 9. Metrics — Prometheus Implementation

**`internal/modules/pdpa/infrastructure/metrics/prometheus.go`** (ใหม่)

```go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
)

// PromMetrics implements eventhandler.Metrics with Prometheus
type PromMetrics struct {
	processed *prometheus.CounterVec
	duration  *prometheus.HistogramVec
}

func NewPromMetrics() *PromMetrics {
	return &PromMetrics{
		processed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "pdpa_event_handler_total",
			Help: "Total number of events processed per handler",
		}, []string{"handler", "result"}),

		duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "pdpa_event_handler_duration_seconds",
			Help:    "Event handler execution duration",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2.5, 5, 10, 30},
		}, []string{"handler"}),
	}
}

func (m *PromMetrics) IncHandlerProcessed(handler, result string) {
	m.processed.WithLabelValues(handler, result).Inc()
}

func (m *PromMetrics) ObserveHandlerDuration(handler string, seconds float64) {
	m.duration.WithLabelValues(handler).Observe(seconds)
}

var _ eventhandler.Metrics = (*PromMetrics)(nil)
```

---

## 10. Module Cleanup — แก้โครงสร้าง tree

จาก tree ที่คุณส่งมา มี **2 ปัญหาหลัก**:

1. **`pdpa/repository/`** ที่ root ซ้ำกับ `pdpa/domain/repository/` — **ลบออก**
2. **`interfaces/handlers/`** vs **`interfaces/http/`** vs **`interfaces/routes/`** — ควรแยกให้ชัด

**โครงสร้างที่แนะนำ:**

```diff
 internal/modules/pdpa/
-├── repository/                        ❌ ลบ (ซ้ำกับ domain/repository)
 ├── application/
-│   └── dto/                           ← เติมไฟล์จริง
+│   ├── command/                       ✅ ตาม pattern
+│   ├── query/
+│   ├── event_handler/
+│   │   ├── envelope.go                ✅ ใหม่
+│   │   ├── base_handler.go            ✅ refactor
+│   │   ├── registry.go                ✅ ใหม่
+│   │   ├── dispatcher.go              ✅ ใหม่
+│   │   ├── pii_redactor.go            ✅ ใหม่
+│   │   ├── wiring.go                  ✅ ใหม่
+│   │   └── handlers/                  ✅ refactor
+│   │       ├── consent_granted_handler.go
+│   │       ├── consent_revoked_handler.go
+│   │       ├── dsar_submitted_handler.go
+│   │       ├── dsar_completed_handler.go
+│   │       ├── account_suspended_handler.go
+│   │       ├── account_terminated_handler.go
+│   │       ├── data_deletion_requested_handler.go
+│   │       └── dlq_handler.go
+│   └── dto/                           ✅ DTOs
 ├── domain/                            ✅ คงเดิม
 ├── infrastructure/
 │   ├── messaging/kafka/consumer/
+│   │   └── event_consumer.go          ✅ ใหม่
 │   ├── persistence/redis/
+│   │   └── idempotency_store.go       ✅ ใหม่
-│   ├── es/                            ❌ ลบ (ไม่ใช้)
-│   ├── search/                        ❌ ลบ (ไม่ใช้)
 │   ├── metrics/
+│   │   └── prometheus.go              ✅ ใหม่
 │   └── ...
 └── interfaces/
-    ├── handlers/                      ❌ ย้ายไป application/event_handler
-    ├── ws/                            ❌ ลบ (ซ้ำกับ websocket)
     ├── http/
     └── websocket/
```

---

## 11. Unit Tests

**`internal/modules/pdpa/application/event_handler/handlers/consent_granted_handler_test.go`** (ใหม่)

```go
package handlers_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
	"icmongolang/internal/modules/pdpa/application/event_handler/handlers"
	"icmongolang/pkg/logger"
)

// --- fakes ---

type fakeAuditRepo struct {
	saved []map[string]interface{}
}

func (f *fakeAuditRepo) Save(_ context.Context, a *entity.AuditTrail) error {
	f.saved = append(f.saved, a.Details.(map[string]interface{}))
	return nil
}

type fakeIdempotency struct {
	processed map[string]bool
}

func newFakeIdempotency() *fakeIdempotency {
	return &fakeIdempotency{processed: map[string]bool{}}
}
func (f *fakeIdempotency) IsProcessed(_ context.Context, id uuid.UUID, h string) (bool, error) {
	return f.processed[h+":"+id.String()], nil
}
func (f *fakeIdempotency) MarkProcessed(_ context.Context, id uuid.UUID, h string, _ time.Duration) error {
	f.processed[h+":"+id.String()] = true
	return nil
}

type fakeCache struct{ sets int }
func (f *fakeCache) Set(context.Context, uuid.UUID, valueobject.ConsentPurpose, valueobject.ConsentStatus) error {
	f.sets++; return nil
}
func (f *fakeCache) Get(context.Context, uuid.UUID, valueobject.ConsentPurpose) (valueobject.ConsentStatus, error) {
	return "", nil
}
func (f *fakeCache) DeleteAllByUser(context.Context, uuid.UUID) error { return nil }

type fakeProducer struct{ published []string }
func (f *fakeProducer) Publish(_ context.Context, topic, _ string, _ []byte) error {
	f.published = append(f.published, topic); return nil
}
func (f *fakeProducer) Close() error { return nil }

// --- tests ---

func TestConsentGrantedHandler_Success(t *testing.T) {
	audit := &fakeAuditRepo{}
	cache := &fakeCache{}
	producer := &fakeProducer{}
	idem := newFakeIdempotency()

	h := handlers.NewConsentGrantedHandler(audit, cache, producer, idem, eventhandler.NopMetrics(), logger.Nop())

	payload := eventhandler.ConsentGrantedPayload{
		ConsentID: uuid.New(),
		UserID:    uuid.New(),
		Purpose:   "ANALYTICS",
		GrantedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(1, 0, 0),
		IPAddress: "203.0.113.42",
	}
	raw, _ := json.Marshal(payload)
	env := &eventhandler.Envelope{
		EventType: "pdpa.consent.granted",
		Metadata:  eventhandler.EventMetadata{EventID: uuid.New(), CorrelationID: "test-1"},
		Payload:   raw,
	}

	require.NoError(t, h.Handle(context.Background(), env))
	assert.Equal(t, 1, cache.sets)
	assert.Contains(t, producer.published, "pdpa.email.notification")
	require.Len(t, audit.saved, 1)
	// ตรวจว่า IP ถูก redact
	assert.Equal(t, "203.0.113.0", audit.saved[0]["ip_address"]) // ตาม PII rule
}

func TestConsentGrantedHandler_Idempotent(t *testing.T) {
	audit := &fakeAuditRepo{}
	cache := &fakeCache{}
	idem := newFakeIdempotency()
	h := handlers.NewConsentGrantedHandler(audit, cache, &fakeProducer{}, idem, eventhandler.NopMetrics(), logger.Nop())

	env := &eventhandler.Envelope{
		EventType: "pdpa.consent.granted",
		Metadata:  eventhandler.EventMetadata{EventID: uuid.New()},
		Payload:   []byte(`{"consent_id":"` + uuid.New().String() + `","user_id":"` + uuid.New().String() + `","purpose":"ANALYTICS","granted_at":"2026-01-01T00:00:00Z","expires_at":"2027-01-01T00:00:00Z"}`),
	}

	require.NoError(t, h.Handle(context.Background(), env))
	require.NoError(t, h.Handle(context.Background(), env)) // ซ้ำ

	// ต้องถูก process แค่ครั้งเดียว
	assert.Len(t, audit.saved, 1)
	assert.Equal(t, 1, cache.sets)
}
```

**`internal/modules/pdpa/application/event_handler/pii_redactor_test.go`** (ใหม่)

```go
package eventhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	eventhandler "icmongolang/internal/modules/pdpa/application/event_handler"
)

func TestPIIRedactor_Email(t *testing.T) {
	r := eventhandler.NewPIIRedactor()
	got := r.RedactPayload(map[string]interface{}{"message": "contact me at john.doe@example.com"})
	assert.Equal(t, "***@***.***", got["message"])
}

func TestPIIRedactor_SensitiveKeys(t *testing.T) {
	r := eventhandler.NewPIIRedactor()
	got := r.RedactPayload(map[string]interface{}{
		"email":       "john@example.com",
		"phone":       "0812345678",
		"otp_code":    "123456",
		"national_id": "1234567890123",
		"safe_field":  "hello",
	})
	assert.Equal(t, "***", got["email"])
	assert.Equal(t, "***", got["phone"])
	assert.Equal(t, "***", got["otp_code"])
	assert.Equal(t, "***", got["national_id"])
	assert.Equal(t, "hello", got["safe_field"])
}
```

---

## 📊 สรุปผลลัพธ์

| ด้าน | ก่อน | หลัง |
|------|------|------|
| **Idempotency** | ❌ ไม่มี | ✅ Redis SET NX (atomic) |
| **PII log** | ❌ log ข้อมูลเต็ม | ✅ Redact ทั้ง key + pattern |
| **Type safety** | ❌ `map[string]any` | ✅ Typed `Envelope[T]` |
| **Business logic** | ❌ save audit เท่านั้น | ✅ cache + email + LLM + blockchain |
| **Dispatch** | ❌ ต้อง wiring เอง | ✅ Registry + Dispatcher |
| **Metrics** | ❌ ไม่มี | ✅ Prometheus counters + histogram |
| **DLQ** | ❌ กลืน error | ✅ Return error ให้ consumer retry |
| **Observability** | ❌ log ปนกัน | ✅ Structured + correlation_id |
| **Module tree** | ❌ `repository/` ซ้ำ | ✅ ลบออก |

**พร้อมไปต่อหรือยังครับ?** ต้องการให้ผม:
- 📝 เขียน **`ConfirmDeletionCommand` + `RejectDSARCommand` + `SuspendAccountCommand` + `TerminateAccountCommand` + `RevokeConsentCommand`** ให้ใช้ outbox pattern แบบเดียวกัน?
- 🧪 เขียน **`uow_test.go`** + **integration test**?
- 📊 วาด **sequence diagram** ของ flow ใหม่?
- 🔧 เขียน **`.golangci.yml`** + **Makefile** สำหรับ quality gate?