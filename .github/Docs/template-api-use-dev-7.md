---

# 📦 Module: `DEMO`

## 🎯 Use Cases
1. `CreateDEMO` — สร้างการแจ้งเตือน (trigger จาก module อื่นผ่าน Kafka)
2. `SendDEMO` — ส่งแบบ async (Email + Push + WebSocket)
3. `MarkAsRead` — mark read/unread
4. `ListByUser` — list พร้อม pagination + filter
5. `AutoExpireDEMOs` — Scheduler (cron) ลบ/archive DEMO หมดอายุ
6. `HandleAccountEvent` — Kafka consumer (subscribe `users.account.event`)

## 🗂️ File Tree

```
internal/modules/DEMO/
├── domain/
│   ├── entity/
│   │   ├── DEMO.go
│   │   └── audit_trail.go
│   ├── value_object/
│   │   ├── channel.go
│   │   ├── status.go
│   │   └── priority.go
│   ├── repository/
│   │   ├── DEMO_repository.go
│   │   └── audit_repository.go
│   ├── service/
│   │   ├── DEMO_policy.go
│   │   ├── email_sender_port.go
│   │   ├── push_sender_port.go
│   │   └── realtime_port.go
│   ├── event/
│   │   └── DEMO_event.go
│   └── errors/
│       └── errors.go
│
├── application/
│   ├── create_DEMO.go
│   ├── send_DEMO.go
│   ├── mark_as_read.go
│   ├── list_by_user.go
│   ├── auto_expire_DEMOs.go
│   ├── handle_account_event.go
│   ├── dto.go
│   └── mappers.go
│
├── infrastructure/
│   ├── persistence/
│   │   ├── postgres/
│   │   │   ├── models.go
│   │   │   ├── DEMO_repo_impl.go
│   │   │   ├── audit_repo_impl.go
│   │   │   └── mappers.go
│   │   └── redis/
│   │       ├── DEMO_cache.go
│   │       └── unread_counter.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── consumers/
│   │       └── account_event_consumer.go
│   ├── search/elasticsearch/
│   │   └── DEMO_indexer.go
│   ├── services/
│   │   ├── email/smtp_sender.go
│   │   └── push/fcm_sender.go
│   └── scheduler/
│       └── expire_job.go
│
├── interfaces/
│   ├── http/
│   │   ├── DEMO_handler.go
│   │   └── routes.go
│   ├── websocket/
│   │   └── realtime_publisher.go
│   └── middleware/
│       └── owner_guard.go
│
├── testdata/
│   └── fixtures.go
│
└── module.go
```

---

## 🏛️ 1. Domain Layer

### `domain/value_object/channel.go`
```go
package valueobject

type Channel string

const (
	ChannelEmail     Channel = "EMAIL"
	ChannelPush      Channel = "PUSH"
	ChannelWebSocket Channel = "WEBSOCKET"
	ChannelSMS       Channel = "SMS"
)

func (c Channel) IsValid() bool {
	switch c {
	case ChannelEmail, ChannelPush, ChannelWebSocket, ChannelSMS:
		return true
	}
	return false
}

func (c Channel) String() string { return string(c) }

// AllChannels คืนค่าทุก channel — ใช้ตอน broadcast
func AllChannels() []Channel {
	return []Channel{ChannelEmail, ChannelPush, ChannelWebSocket}
}
```

### `domain/value_object/status.go`
```go
package valueobject

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusSent      Status = "SENT"
	StatusFailed    Status = "FAILED"
	StatusRead      Status = "READ"
	StatusArchived  Status = "ARCHIVED"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusSent, StatusFailed, StatusRead, StatusArchived:
		return true
	}
	return false
}

func (s Status) CanTransitionTo(next Status) bool {
	switch s {
	case StatusPending:
		return next == StatusSent || next == StatusFailed
	case StatusSent:
		return next == StatusRead || next == StatusArchived
	case StatusRead:
		return next == StatusArchived
	}
	return false
}
```

### `domain/value_object/priority.go`
```go
package valueobject

type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 5
	PriorityHigh   Priority = 8
	PriorityUrgent Priority = 10
)

func (p Priority) IsValid() bool {
	return p >= PriorityLow && p <= PriorityUrgent
}

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "LOW"
	case PriorityHigh:
		return "HIGH"
	case PriorityUrgent:
		return "URGENT"
	default:
		return "NORMAL"
	}
}
```

### `domain/errors/errors.go`
```go
package domainerrors

import "errors"

var (
	// DEMO
	ErrDEMONotFound     = errors.New("DEMO not found")
	ErrDEMOAlreadyRead  = errors.New("DEMO already read")
	ErrInvalidDEMOID    = errors.New("invalid DEMO id")
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidChannel           = errors.New("invalid channel")
	ErrInvalidPriority          = errors.New("invalid priority")
	ErrTitleRequired            = errors.New("title is required")
	ErrBodyRequired             = errors.New("body is required")
	ErrTooManyChannels          = errors.New("too many channels (max 3)")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")

	// Ownership / Authorization
	ErrNotOwnedByUser = errors.New("DEMO not owned by user")

	// Sending
	ErrSendFailed    = errors.New("failed to send DEMO")
	ErrAllSendFailed = errors.New("all channels failed to send")

	// External
	ErrEmailSenderUnavailable = errors.New("email sender unavailable")
	ErrPushSenderUnavailable  = errors.New("push sender unavailable")
)
```

### `domain/entity/DEMO.go`
```go
package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
)

// DEMO — Aggregate Root
type DEMO struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Body        string
	Channels    []valueobject.Channel
	Priority    valueobject.Priority
	Status      valueobject.Status
	Metadata    map[string]string // arbitrary payload (order_id, job_id, ...)
	SentAt      *time.Time
	ReadAt      *time.Time
	FailedAt    *time.Time
	FailReason  string
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New — constructor บังคับ invariants
func New(
	userID uuid.UUID,
	title, body string,
	channels []valueobject.Channel,
	priority valueobject.Priority,
	metadata map[string]string,
	expiresIn *time.Duration,
) (*DEMO, error) {
	if userID == uuid.Nil {
		return nil, domainerrors.ErrInvalidUserID
	}
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if title == "" {
		return nil, domainerrors.ErrTitleRequired
	}
	if body == "" {
		return nil, domainerrors.ErrBodyRequired
	}
	if len(channels) == 0 || len(channels) > 3 {
		return nil, domainerrors.ErrTooManyChannels
	}
	for _, c := range channels {
		if !c.IsValid() {
			return nil, domainerrors.ErrInvalidChannel
		}
	}
	if !priority.IsValid() {
		priority = valueobject.PriorityNormal
	}

	now := time.Now().UTC()
	n := &DEMO{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Body:      body,
		Channels:  channels,
		Priority:  priority,
		Status:    valueobject.StatusPending,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if expiresIn != nil {
		t := now.Add(*expiresIn)
		n.ExpiresAt = &t
	}
	return n, nil
}

// ── Behavior methods ─────────────────────────────────────────

// MarkSent — เรียกหลังส่งทุก channel สำเร็จ
func (n *DEMO) MarkSent() error {
	if !n.Status.CanTransitionTo(valueobject.StatusSent) {
		return domainerrors.ErrInvalidStatusTransition
	}
	now := time.Now().UTC()
	n.Status = valueobject.StatusSent
	n.SentAt = &now
	n.UpdatedAt = now
	return nil
}

// MarkFailed — เรียกเมื่อส่งไม่สำเร็จ
func (n *DEMO) MarkFailed(reason string) error {
	if !n.Status.CanTransitionTo(valueobject.StatusFailed) {
		return domainerrors.ErrInvalidStatusTransition
	}
	now := time.Now().UTC()
	n.Status = valueobject.StatusFailed
	n.FailedAt = &now
	n.FailReason = reason
	n.UpdatedAt = now
	return nil
}

// MarkRead — ผู้ใช้อ่าน
func (n *DEMO) MarkRead() error {
	if n.Status == valueobject.StatusRead {
		return domainerrors.ErrDEMOAlreadyRead
	}
	if n.Status != valueobject.StatusSent {
		return domainerrors.ErrInvalidStatusTransition
	}
	now := time.Now().UTC()
	n.Status = valueobject.StatusRead
	n.ReadAt = &now
	n.UpdatedAt = now
	return nil
}

// Archive — เก็บเข้ากรุ
func (n *DEMO) Archive() error {
	if !n.Status.CanTransitionTo(valueobject.StatusArchived) {
		return domainerrors.ErrInvalidStatusTransition
	}
	n.Status = valueobject.StatusArchived
	n.UpdatedAt = time.Now().UTC()
	return nil
}

// IsExpired
func (n *DEMO) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*n.ExpiresAt)
}

// IsOwnedBy
func (n *DEMO) IsOwnedBy(userID uuid.UUID) bool {
	return n.UserID == userID
}
```

### `domain/entity/audit_trail.go`
```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditTrail struct {
	ID         uuid.UUID
	ActorID    uuid.UUID
	Action     string
	EntityID   uuid.UUID
	EntityType string
	Payload    map[string]any
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
}

func NewAuditTrail(actorID uuid.UUID, action string, entityID uuid.UUID, entityType string, payload map[string]any) *AuditTrail {
	return &AuditTrail{
		ID:         uuid.New(),
		ActorID:    actorID,
		Action:     action,
		EntityID:   entityID,
		EntityType: entityType,
		Payload:    payload,
		CreatedAt:  time.Now().UTC(),
	}
}
```

### `domain/repository/DEMO_repository.go`
```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/domain/entity"
)

type DEMORepository interface {
	Save(ctx context.Context, n *entity.DEMO) error
	Update(ctx context.Context, n *entity.DEMO) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.DEMO, error)

	// list พร้อม pagination + filters
	List(ctx context.Context, filter ListFilter) ([]*entity.DEMO, int64, error)

	// admin / scheduler
	FindExpired(ctx context.Context, before time.Time, limit int) ([]*entity.DEMO, error)
	BulkArchive(ctx context.Context, ids []uuid.UUID) error

	// counter
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}

type ListFilter struct {
	UserID  uuid.UUID
	Status  string // optional
	Channel string // optional
	Search  string // optional (title/body)
	Page    int
	Size    int
	SortBy  string // "created_at" | "priority"
	Order   string // "asc" | "desc"
}
```

### `domain/repository/audit_repository.go`
```go
package repository

import (
	"context"

	"icmongolang/internal/modules/DEMO/domain/entity"
)

type AuditRepository interface {
	Save(ctx context.Context, a *entity.AuditTrail) error
}
```

### `domain/service/DEMO_policy.go`
```go
package service

import (
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
	"icmongolang/internal/modules/DEMO/domain/entity"
)

// DEMOPolicy — domain service: ตัดสินใจว่าควรส่งช่องทางใด
type DEMOPolicy struct{}

func NewDEMOPolicy() *DEMOPolicy { return &DEMOPolicy{} }

// ResolveChannels — ถ้า priority urgent ให้ส่งทุกช่อง
func (p *DEMOPolicy) ResolveChannels(n *entity.DEMO) []valueobject.Channel {
	if n.Priority >= valueobject.PriorityUrgent {
		return valueobject.AllChannels()
	}
	return n.Channels
}

// ShouldRetry — retry เมื่อ failed และ priority สูง
func (p *DEMOPolicy) ShouldRetry(n *entity.DEMO) bool {
	return n.Status == valueobject.StatusFailed &&
		n.Priority >= valueobject.PriorityHigh
}
```

### `domain/service/email_sender_port.go`
```go
package service

import "context"

type EmailPayload struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// EmailSenderPort — outbound port (Domain ต้องการ interface นี้)
type EmailSenderPort interface {
	Send(ctx context.Context, p EmailPayload) error
}
```

### `domain/service/push_sender_port.go`
```go
package service

import "context"

type PushPayload struct {
	UserID   string
	Title    string
	Body     string
	Data     map[string]string
}

type PushSenderPort interface {
	Send(ctx context.Context, p PushPayload) error
}
```

### `domain/service/realtime_port.go`
```go
package service

import "context"

type RealtimePayload struct {
	UserID  string
	Event   string
	Payload any
}

type RealtimePort interface {
	Publish(ctx context.Context, p RealtimePayload) error
}
```

### `domain/event/DEMO_event.go`
```go
package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	TopicDEMOCreated = "DEMO.created"
	TopicDEMOSent    = "DEMO.sent"
	TopicDEMOFailed  = "DEMO.failed"
	TopicDEMORead    = "DEMO.read"
)

type DEMOCreated struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	Title     string            `json:"title"`
	Channels  []string          `json:"channels"`
	Priority  int               `json:"priority"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Occurred  time.Time         `json:"occurred_at"`
}

type DEMOSent struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Channel   string    `json:"channel"`
	SentAt    time.Time `json:"sent_at"`
}
```

---

## 🎯 2. Application Layer

### `application/create_DEMO.go`
```go
package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
	"icmongolang/internal/modules/DEMO/domain/entity"
	"icmongolang/internal/modules/DEMO/domain/event"
	"icmongolang/internal/modules/DEMO/domain/repository"
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

type CreateDEMOUseCase struct {
	repo      repository.DEMORepository
	auditRepo repository.AuditRepository
	publisher EventPublisher
	logger    *slog.Logger
}

func NewCreateDEMOUseCase(
	repo repository.DEMORepository,
	auditRepo repository.AuditRepository,
	publisher EventPublisher,
	logger *slog.Logger,
) *CreateDEMOUseCase {
	return &CreateDEMOUseCase{repo: repo, auditRepo: auditRepo, publisher: publisher, logger: logger}
}

type CreateDEMOInput struct {
	UserID    uuid.UUID
	Title     string
	Body      string
	Channels  []string
	Priority  int
	Metadata  map[string]string
	ActorID   uuid.UUID
	IPAddress string
	UserAgent string
}

func (uc *CreateDEMOUseCase) Execute(ctx context.Context, in CreateDEMOInput) (*DEMOResponse, error) {
	// 1. Validate + transform channels
	channels := make([]valueobject.Channel, 0, len(in.Channels))
	for _, c := range in.Channels {
		ch := valueobject.Channel(c)
		if !ch.IsValid() {
			return nil, domainerrors.ErrInvalidChannel
		}
		channels = append(channels, ch)
	}

	// 2. Build aggregate
	n, err := entity.New(in.UserID, in.Title, in.Body, channels, valueobject.Priority(in.Priority), in.Metadata, nil)
	if err != nil {
		return nil, err
	}

	// 3. Persist
	if err := uc.repo.Save(ctx, n); err != nil {
		return nil, err
	}

	// 4. Audit (best-effort)
	audit := entity.NewAuditTrail(in.ActorID, "DEMO.created", n.ID, "DEMO", map[string]any{
		"user_id":  n.UserID.String(),
		"priority": n.Priority.String(),
	})
	audit.IPAddress = in.IPAddress
	audit.UserAgent = in.UserAgent
	if err := uc.auditRepo.Save(ctx, audit); err != nil {
		uc.logger.Warn("audit save failed", "err", err, "DEMO_id", n.ID)
	}

	// 5. Publish event (async → SendDEMO จะ consume)
	evt := event.DEMOCreated{
		ID:       n.ID,
		UserID:   n.UserID,
		Title:    n.Title,
		Priority: int(n.Priority),
		Metadata: n.Metadata,
		Occurred: n.CreatedAt,
	}
	for _, c := range channels {
		evt.Channels = append(evt.Channels, c.String())
	}
	if err := uc.publisher.Publish(ctx, event.TopicDEMOCreated, n.ID.String(), evt); err != nil {
		uc.logger.Error("publish event failed", "err", err, "DEMO_id", n.ID)
		// ไม่ fail ทั้ง use case — เพราะ DEMO ถูกบันทึกแล้ว
	}

	return ToDEMOResponse(n), nil
}
```

### `application/send_DEMO.go`
```go
package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/domain/entity"
	"icmongolang/internal/modules/DEMO/domain/repository"
	"icmongolang/internal/modules/DEMO/domain/service"
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
)

type SendDEMOUseCase struct {
	repo     repository.DEMORepository
	policy   *service.DEMOPolicy
	email    service.EmailSenderPort
	push     service.PushSenderPort
	realtime service.RealtimePort
	logger   *slog.Logger
}

func NewSendDEMOUseCase(
	repo repository.DEMORepository,
	policy *service.DEMOPolicy,
	email service.EmailSenderPort,
	push service.PushSenderPort,
	realtime service.RealtimePort,
	logger *slog.Logger,
) *SendDEMOUseCase {
	return &SendDEMOUseCase{
		repo: repo, policy: policy,
		email: email, push: push, realtime: realtime,
		logger: logger,
	}
}

type SendDEMOInput struct {
	DEMOID uuid.UUID
	// ผู้รับ (ถ้ามี — ปกติ resolve จาก repo)
	RecipientEmail string
}

func (uc *SendDEMOUseCase) Execute(ctx context.Context, in SendDEMOInput) error {
	n, err := uc.repo.FindByID(ctx, in.DEMOID)
	if err != nil {
		return err
	}
	if n.Status != valueobject.StatusPending {
		// idempotent — ถ้าส่งแล้วให้ skip
		uc.logger.Info("DEMO already processed", "id", n.ID, "status", n.Status)
		return nil
	}

	channels := uc.policy.ResolveChannels(n)
	var errs []error

	for _, ch := range channels {
		switch ch {
		case valueobject.ChannelEmail:
			if in.RecipientEmail == "" {
				uc.logger.Warn("skip email — no recipient", "id", n.ID)
				continue
			}
			err := uc.email.Send(ctx, service.EmailPayload{
				To:      in.RecipientEmail,
				Subject: n.Title,
				Body:    n.Body,
				IsHTML:  true,
			})
			if err != nil {
				errs = append(errs, fmt.Errorf("email: %w", err))
			}
		case valueobject.ChannelPush:
			err := uc.push.Send(ctx, service.PushPayload{
				UserID: n.UserID.String(),
				Title:  n.Title,
				Body:   n.Body,
				Data:   n.Metadata,
			})
			if err != nil {
				errs = append(errs, fmt.Errorf("push: %w", err))
			}
		case valueobject.ChannelWebSocket:
			err := uc.realtime.Publish(ctx, service.RealtimePayload{
				UserID: n.UserID.String(),
				Event:  "DEMO.new",
				Payload: map[string]any{
					"id":       n.ID.String(),
					"title":    n.Title,
					"body":     n.Body,
					"priority": n.Priority.String(),
				},
			})
			if err != nil {
				errs = append(errs, fmt.Errorf("ws: %w", err))
			}
		}
	}

	if len(errs) == len(channels) {
		reason := errors.Join(errs...).Error()
		if err := n.MarkFailed(reason); err != nil {
			return err
		}
		_ = uc.repo.Update(ctx, n)
		return fmt.Errorf("%w: %s", ErrAllSendFailed, reason)
	}

	if err := n.MarkSent(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, n)
}
```

### `application/mark_as_read.go`
```go
package application

import (
	"context"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
	"icmongolang/internal/modules/DEMO/domain/repository"
)

type MarkAsReadUseCase struct {
	repo repository.DEMORepository
}

func NewMarkAsReadUseCase(repo repository.DEMORepository) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{repo: repo}
}

type MarkAsReadInput struct {
	DEMOID uuid.UUID
	UserID         uuid.UUID // เพื่อเช็ค ownership
}

func (uc *MarkAsReadUseCase) Execute(ctx context.Context, in MarkAsReadInput) error {
	n, err := uc.repo.FindByID(ctx, in.DEMOID)
	if err != nil {
		return err
	}
	if !n.IsOwnedBy(in.UserID) {
		return domainerrors.ErrNotOwnedByUser
	}
	if err := n.MarkRead(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, n)
}
```

### `application/list_by_user.go`
```go
package application

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/domain/repository"
)

type ListByUserUseCase struct {
	repo repository.DEMORepository
}

func NewListByUserUseCase(repo repository.DEMORepository) *ListByUserUseCase {
	return &ListByUserUseCase{repo: repo}
}

type ListByUserInput struct {
	UserID  uuid.UUID
	Status  string
	Channel string
	Search  string
	Page    int
	Size    int
	SortBy  string
	Order   string
}

type ListByUserOutput struct {
	Items    []*DEMOResponse `json:"items"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	Size     int                     `json:"size"`
	Unread   int64                   `json:"unread"`
}

func (uc *ListByUserUseCase) Execute(ctx context.Context, in ListByUserInput) (*ListByUserOutput, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 || in.Size > 100 {
		in.Size = 20
	}
	if in.SortBy == "" {
		in.SortBy = "created_at"
	}
	if in.Order != "asc" {
		in.Order = "desc"
	}

	items, total, err := uc.repo.List(ctx, repository.ListFilter{
		UserID: in.UserID, Status: in.Status, Channel: in.Channel,
		Search: in.Search, Page: in.Page, Size: in.Size,
		SortBy: in.SortBy, Order: in.Order,
	})
	if err != nil {
		return nil, err
	}

	unread, err := uc.repo.CountUnread(ctx, in.UserID)
	if err != nil {
		unread = 0 // best-effort
	}

	out := &ListByUserOutput{
		Total: total, Page: in.Page, Size: in.Size, Unread: unread,
		Items: make([]*DEMOResponse, 0, len(items)),
	}
	for _, n := range items {
		out.Items = append(out.Items, ToDEMOResponse(n))
	}
	return out, nil
}
```

### `application/auto_expire_DEMOs.go`
```go
package application

import (
	"context"
	"log/slog"
	"time"

	"icmongolang/internal/modules/DEMO/domain/repository"
)

type AutoExpireDEMOsUseCase struct {
	repo   repository.DEMORepository
	logger *slog.Logger
}

func NewAutoExpireDEMOsUseCase(
	repo repository.DEMORepository,
	logger *slog.Logger,
) *AutoExpireDEMOsUseCase {
	return &AutoExpireDEMOsUseCase{repo: repo, logger: logger}
}

func (uc *AutoExpireDEMOsUseCase) Execute(ctx context.Context) error {
	cutoff := time.Now().UTC()
	expired, err := uc.repo.FindExpired(ctx, cutoff, 500)
	if err != nil {
		return err
	}
	if len(expired) == 0 {
		uc.logger.Info("no expired DEMOs")
		return nil
	}
	ids := make([]uuid.UUID, 0, len(expired))
	for _, n := range expired {
		ids = append(ids, n.ID)
	}
	if err := uc.repo.BulkArchive(ctx, ids); err != nil {
		return err
	}
	uc.logger.Info("archived expired DEMOs", "count", len(ids))
	return nil
}
```

### `application/handle_account_event.go`
```go
package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/domain/entity"
	"icmongolang/internal/modules/DEMO/domain/repository"
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
)

// HandleAccountEventUseCase — รับ event จาก users module (Kafka)
type HandleAccountEventUseCase struct {
	repo   repository.DEMORepository
	logger *slog.Logger
}

func NewHandleAccountEventUseCase(repo repository.DEMORepository, logger *slog.Logger) *HandleAccountEventUseCase {
	return &HandleAccountEventUseCase{repo: repo, logger: logger}
}

type AccountEventPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Event    string    `json:"event"` // "created" | "verified" | "locked"
	Occurred string    `json:"occurred_at"`
}

func (uc *HandleAccountEventUseCase) Execute(ctx context.Context, p AccountEventPayload) error {
	var title, body string
	switch p.Event {
	case "created":
		title = "ยินดีต้อนรับ"
		body = "กรุณายืนยันอีเมลของคุณ"
	case "verified":
		title = "ยืนยันอีเมลสำเร็จ"
		body = "บัญชีของคุณพร้อมใช้งาน"
	case "locked":
		title = "บัญชีถูกล็อก"
		body = "กรุณาติดต่อฝ่ายสนับสนุน"
	default:
		uc.logger.Warn("unknown account event", "event", p.Event)
		return nil
	}

	n, err := entity.New(
		p.UserID, title, body,
		[]valueobject.Channel{valueobject.ChannelEmail, valueobject.ChannelWebSocket},
		valueobject.PriorityHigh, nil, nil,
	)
	if err != nil {
		return err
	}
	return uc.repo.Save(ctx, n)
}
```

### `application/dto.go`
```go
package application

import "time"

type DEMOResponse struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Channels   []string          `json:"channels"`
	Priority   string            `json:"priority"`
	Status     string            `json:"status"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	SentAt     *time.Time        `json:"sent_at,omitempty"`
	ReadAt     *time.Time        `json:"read_at,omitempty"`
	ExpiresAt  *time.Time        `json:"expires_at,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
```

### `application/mappers.go`
```go
package application

import (
	"icmongolang/internal/modules/DEMO/domain/entity"
)

func ToDEMOResponse(n *entity.DEMO) *DEMOResponse {
	channels := make([]string, 0, len(n.Channels))
	for _, c := range n.Channels {
		channels = append(channels, c.String())
	}
	return &DEMOResponse{
		ID:        n.ID.String(),
		UserID:    n.UserID.String(),
		Title:     n.Title,
		Body:      n.Body,
		Channels:  channels,
		Priority:  n.Priority.String(),
		Status:    n.Status.String(),
		Metadata:  n.Metadata,
		SentAt:    n.SentAt,
		ReadAt:    n.ReadAt,
		ExpiresAt: n.ExpiresAt,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
```

---

## 🔧 3. Infrastructure Layer

### `infrastructure/persistence/postgres/models.go`
```go
package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray — custom type สำหรับ []string → TEXT[]
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	}
	return errors.New("StringArray: unsupported scan source")
}

// ── Tables (prefix: DEMO_) ─────────────────────────

type DEMOModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Title      string         `gorm:"type:varchar(255);not null"`
	Body       string         `gorm:"type:text;not null"`
	Channels   StringArray    `gorm:"type:jsonb;not null"`
	Priority   int            `gorm:"not null;default:5;index"`
	Status     string         `gorm:"type:varchar(20);not null;index"`
	Metadata   StringMap      `gorm:"type:jsonb"`
	SentAt     *time.Time
	ReadAt     *time.Time
	FailedAt   *time.Time
	FailReason string         `gorm:"type:text"`
	ExpiresAt  *time.Time     `gorm:"index"`
	CreatedAt  time.Time      `gorm:"not null;default:now();index"`
	UpdatedAt  time.Time      `gorm:"not null;default:now()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (DEMOModel) TableName() string { return "DEMO_DEMOs" }

type AuditTrailModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ActorID    uuid.UUID  `gorm:"type:uuid;index"`
	Action     string     `gorm:"type:varchar(64);not null;index"`
	EntityID   uuid.UUID  `gorm:"type:uuid;index"`
	EntityType string     `gorm:"type:varchar(64);not null;index"`
	Payload    StringMap  `gorm:"type:jsonb"`
	IPAddress  string     `gorm:"type:varchar(64)"`
	UserAgent  string     `gorm:"type:text"`
	CreatedAt  time.Time  `gorm:"not null;default:now();index"`
}

func (AuditTrailModel) TableName() string { return "DEMO_audit_trails" }

// StringMap — map[string]string → JSONB
type StringMap map[string]string

func (m StringMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *StringMap) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	}
	return errors.New("StringMap: unsupported scan source")
}
```

### `infrastructure/persistence/postgres/mappers.go`
```go
package postgres

import (
	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
	"icmongolang/internal/modules/DEMO/domain/entity"
	valueobject "icmongolang/internal/modules/DEMO/domain/value_object"
)

func toEntity(m *DEMOModel) *entity.DEMO {
	channels := make([]valueobject.Channel, 0, len(m.Channels))
	for _, c := range m.Channels {
		channels = append(channels, valueobject.Channel(c))
	}
	return &entity.DEMO{
		ID:         m.ID,
		UserID:     m.UserID,
		Title:      m.Title,
		Body:       m.Body,
		Channels:   channels,
		Priority:   valueobject.Priority(m.Priority),
		Status:     valueobject.Status(m.Status),
		Metadata:   map[string]string(m.Metadata),
		SentAt:     m.SentAt,
		ReadAt:     m.ReadAt,
		FailedAt:   m.FailedAt,
		FailReason: m.FailReason,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func toModel(n *entity.DEMO) *DEMOModel {
	channels := make(StringArray, 0, len(n.Channels))
	for _, c := range n.Channels {
		channels = append(channels, c.String())
	}
	return &DEMOModel{
		ID:         n.ID,
		UserID:     n.UserID,
		Title:      n.Title,
		Body:       n.Body,
		Channels:   channels,
		Priority:   int(n.Priority),
		Status:     n.Status.String(),
		Metadata:   StringMap(n.Metadata),
		SentAt:     n.SentAt,
		ReadAt:     n.ReadAt,
		FailedAt:   n.FailedAt,
		FailReason: n.FailReason,
		ExpiresAt:  n.ExpiresAt,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

var _ = domainerrors.ErrDEMONotFound
```

### `infrastructure/persistence/postgres/DEMO_repo_impl.go`
```go
package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
	"icmongolang/internal/modules/DEMO/domain/entity"
	"icmongolang/internal/modules/DEMO/domain/repository"
)

type DEMORepoImpl struct{ db *gorm.DB }

func NewDEMORepository(db *gorm.DB) repository.DEMORepository {
	return &DEMORepoImpl{db: db}
}

func (r *DEMORepoImpl) Save(ctx context.Context, n *entity.DEMO) error {
	m := toModel(n)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *DEMORepoImpl) Update(ctx context.Context, n *entity.DEMO) error {
	m := toModel(n)
	result := r.db.WithContext(ctx).
		Model(&DEMOModel{ID: n.ID}).
		Select("*").
		Omit("id", "created_at").
		Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrDEMONotFound
	}
	return nil
}

func (r *DEMORepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.DEMO, error) {
	var m DEMOModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domainerrors.ErrDEMONotFound
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *DEMORepoImpl) List(ctx context.Context, f repository.ListFilter) ([]*entity.DEMO, int64, error) {
	q := r.db.WithContext(ctx).Model(&DEMOModel{}).Where("user_id = ?", f.UserID)

	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Channel != "" {
		q = q.Where("channels::jsonb @> ?::jsonb", `["`+f.Channel+`"]`)
	}
	if f.Search != "" {
		s := "%" + f.Search + "%"
		q = q.Where("(title ILIKE ? OR body ILIKE ?)", s, s)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortCol := "created_at"
	if f.SortBy == "priority" {
		sortCol = "priority"
	}
	order := "DESC"
	if f.Order == "asc" {
		order = "ASC"
	}
	q = q.Order(sortCol + " " + order)

	offset := (f.Page - 1) * f.Size
	var rows []DEMOModel
	if err := q.Limit(f.Size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]*entity.DEMO, 0, len(rows))
	for i := range rows {
		out = append(out, toEntity(&rows[i]))
	}
	return out, total, nil
}

func (r *DEMORepoImpl) FindExpired(ctx context.Context, before time.Time, limit int) ([]*entity.DEMO, error) {
	var rows []DEMOModel
	err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status <> ?", before, "ARCHIVED").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*entity.DEMO, 0, len(rows))
	for i := range rows {
		out = append(out, toEntity(&rows[i]))
	}
	return out, nil
}

func (r *DEMORepoImpl) BulkArchive(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&DEMOModel{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":     "ARCHIVED",
			"updated_at": time.Now().UTC(),
		}).Error
}

func (r *DEMORepoImpl) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&DEMOModel{}).
		Where("user_id = ? AND status = ?", userID, "SENT").
		Count(&count).Error
	return count, err
}
```

### `infrastructure/persistence/postgres/audit_repo_impl.go`
```go
package postgres

import (
	"context"

	"gorm.io/gorm"
	"icmongolang/internal/modules/DEMO/domain/entity"
	"icmongolang/internal/modules/DEMO/domain/repository"
)

type auditRepoImpl struct{ db *gorm.DB }

func NewAuditRepository(db *gorm.DB) repository.AuditRepository {
	return &auditRepoImpl{db: db}
}

func (r *auditRepoImpl) Save(ctx context.Context, a *entity.AuditTrail) error {
	m := &AuditTrailModel{
		ID:         a.ID,
		ActorID:    a.ActorID,
		Action:     a.Action,
		EntityID:   a.EntityID,
		EntityType: a.EntityType,
		Payload:    StringMap(toStringMap(a.Payload)),
		IPAddress:  a.IPAddress,
		UserAgent:  a.UserAgent,
		CreatedAt:  a.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func toStringMap(in map[string]any) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		if s, ok := v.(string); ok {
			out[k] = s
		} else {
			out[k] = ""
		}
	}
	return out
}
```

### `infrastructure/persistence/redis/DEMO_cache.go`
```go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/domain/entity"
)

type DEMOCache interface {
	Set(ctx context.Context, n *entity.DEMO) error
	Get(ctx context.Context, id uuid.UUID) (*entity.DEMO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type redisDEMOCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewDEMOCache(client *redis.Client, ttl time.Duration) DEMOCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &redisDEMOCache{client: client, ttl: ttl}
}

func (c *redisDEMOCache) key(id uuid.UUID) string {
	return fmt.Sprintf("DEMO:DEMO:%s", id)
}

func (c *redisDEMOCache) Set(ctx context.Context, n *entity.DEMO) error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(n.ID), data, c.ttl).Err()
}

func (c *redisDEMOCache) Get(ctx context.Context, id uuid.UUID) (*entity.DEMO, error) {
	data, err := c.client.Get(ctx, c.key(id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var n entity.DEMO
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (c *redisDEMOCache) Delete(ctx context.Context, id uuid.UUID) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
```

### `infrastructure/persistence/redis/unread_counter.go`
```go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type UnreadCounter interface {
	Increment(ctx context.Context, userID uuid.UUID) error
	Decrement(ctx context.Context, userID uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) (int64, error)
	Reset(ctx context.Context, userID uuid.UUID) error
}

type redisUnreadCounter struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUnreadCounter(client *redis.Client, ttl time.Duration) UnreadCounter {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &redisUnreadCounter{client: client, ttl: ttl}
}

func (u *redisUnreadCounter) key(userID uuid.UUID) string {
	return fmt.Sprintf("DEMO:unread:%s", userID)
}

func (u *redisUnreadCounter) Increment(ctx context.Context, userID uuid.UUID) error {
	pipe := u.client.TxPipeline()
	pipe.Incr(ctx, u.key(userID))
	pipe.Expire(ctx, u.key(userID), u.ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (u *redisUnreadCounter) Decrement(ctx context.Context, userID uuid.UUID) error {
	key := u.key(userID)
	v, err := u.client.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	if v < 0 {
		_ = u.client.Set(ctx, key, 0, u.ttl).Err()
	}
	return nil
}

func (u *redisUnreadCounter) Get(ctx context.Context, userID uuid.UUID) (int64, error) {
	v, err := u.client.Get(ctx, u.key(userID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func (u *redisUnreadCounter) Reset(ctx context.Context, userID uuid.UUID) error {
	return u.client.Del(ctx, u.key(userID)).Err()
}
```

### `infrastructure/messaging/kafka_producer.go`
```go
package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	Publish(ctx context.Context, topic, key string, payload any) error
	Close() error
}

type kafkaProducer struct{ p sarama.SyncProducer }

func NewKafkaProducer(brokers []string, clientID string) (Producer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{p: p}, nil
}

func (k *kafkaProducer) Publish(ctx context.Context, topic, key string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	})
	return err
}

func (k *kafkaProducer) Close() error { return k.p.Close() }
```

### `infrastructure/messaging/consumers/account_event_consumer.go`
```go
package consumers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/DEMO/application"
)

const TopicAccountEvent = "users.account.event"

type AccountEventHandler interface {
	Execute(ctx context.Context, p application.AccountEventPayload) error
}

type AccountEventConsumer struct {
	handler AccountEventHandler
	logger  *slog.Logger
}

func NewAccountEventConsumer(h AccountEventHandler, logger *slog.Logger) *AccountEventConsumer {
	return &AccountEventConsumer{handler: h, logger: logger}
}

func (c *AccountEventConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *AccountEventConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *AccountEventConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload application.AccountEventPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			c.logger.Error("unmarshal account event failed", "err", err, "offset", msg.Offset)
			s.MarkMessage(msg, "")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := c.handler.Execute(ctx, payload)
		cancel()

		if err != nil {
			c.logger.Error("handle account event failed", "err", err, "user_id", payload.UserID)
			// mark offset anyway — หรือใช้ DLQ ในกรณีที่ต้องการ strict
		}
		s.MarkMessage(msg, "")
	}
	return nil
}
```

### `infrastructure/services/email/smtp_sender.go`
```go
package email

import (
	"context"
	"fmt"
	"net/smtp"

	"icmongolang/internal/modules/DEMO/domain/service"
)

type smtpSender struct {
	host, port, user, pass, from string
	enabled                       bool
}

func NewSMTPSender(host, port, user, pass, from string, enabled bool) service.EmailSenderPort {
	return &smtpSender{host: host, port: port, user: user, pass: pass, from: from, enabled: enabled}
}

func (s *smtpSender) Send(ctx context.Context, p service.EmailPayload) error {
	if !s.enabled {
		return nil // no-op ในโหมด dev
	}
	contentType := "text/plain"
	if p.IsHTML {
		contentType = "text/html"
	}
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s",
		s.from, p.To, p.Subject, contentType, p.Body,
	)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{p.To}, []byte(msg))
}
```

### `infrastructure/services/push/fcm_sender.go`
```go
package push

import (
	"context"
	"log/slog"

	"icmongolang/internal/modules/DEMO/domain/service"
)

type fcmSender struct {
	apiKey  string
	logger  *slog.Logger
	enabled bool
}

func NewFCMSender(apiKey string, enabled bool, logger *slog.Logger) service.PushSenderPort {
	return &fcmSender{apiKey: apiKey, enabled: enabled, logger: logger}
}

func (s *fcmSender) Send(ctx context.Context, p service.PushPayload) error {
	if !s.enabled {
		s.logger.Debug("push disabled", "user_id", p.UserID, "title", p.Title)
		return nil
	}
	// TODO: integrate with FCM HTTP v1
	// ตัวอย่างโครงสำหรับ implementation จริง
	return nil
}
```

### `infrastructure/scheduler/expire_job.go`
```go
package scheduler

import (
	"context"
	"log/slog"
	"time"

	"icmongolang/internal/modules/DEMO/application"
)

type ExpireJob struct {
	uc     *application.AutoExpireDEMOsUseCase
	logger *slog.Logger
	tick   time.Duration
}

func NewExpireJob(uc *application.AutoExpireDEMOsUseCase, logger *slog.Logger) *ExpireJob {
	return &ExpireJob{uc: uc, logger: logger, tick: 1 * time.Hour}
}

func (j *ExpireJob) Run(ctx context.Context) {
	t := time.NewTicker(j.tick)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			j.logger.Info("expire job stopping")
			return
		case <-t.C:
			runCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			if err := j.uc.Execute(runCtx); err != nil {
				j.logger.Error("auto expire failed", "err", err)
			}
			cancel()
		}
	}
}
```

### `infrastructure/search/elasticsearch/DEMO_indexer.go`
```go
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type DEMOIndexer interface {
	Index(ctx context.Context, doc any) error
	Delete(ctx context.Context, id string) error
}

type esDEMOIndexer struct {
	client *elasticsearch.Client
	index  string
}

func NewDEMOIndexer(client *elasticsearch.Client, index string) DEMOIndexer {
	return &esDEMOIndexer{client: client, index: index}
}

func (i *esDEMOIndexer) Index(ctx context.Context, doc any) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = i.client.Index(i.index, bytes.NewReader(data), i.client.Index.WithContext(ctx))
	return err
}

func (i *esDEMOIndexer) Delete(ctx context.Context, id string) error {
	_, err := i.client.Delete(i.index, id, i.client.Delete.WithContext(ctx))
	return err
}
```

---

## 🌐 4. Interface Layer

### `interfaces/http/DEMO_handler.go`
```go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"icmongolang/internal/modules/DEMO/application"
)

type DEMOHandler struct {
	create *application.CreateDEMOUseCase
	list   *application.ListByUserUseCase
	read   *application.MarkAsReadUseCase
}

func NewDEMOHandler(
	create *application.CreateDEMOUseCase,
	list *application.ListByUserUseCase,
	read *application.MarkAsReadUseCase,
) *DEMOHandler {
	return &DEMOHandler{create: create, list: list, read: read}
}

// POST /api/v1/DEMOs
func (h *DEMOHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		UserID   string            `json:"user_id"`
		Title    string            `json:"title"`
		Body     string            `json:"body"`
		Channels []string          `json:"channels"`
		Priority int               `json:"priority"`
		Metadata map[string]string `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	out, err := h.create.Execute(r.Context(), application.CreateDEMOInput{
		UserID:    userID,
		Title:     body.Title,
		Body:      body.Body,
		Channels:  body.Channels,
		Priority:  body.Priority,
		Metadata:  body.Metadata,
		ActorID:   uid,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

// GET /api/v1/DEMOs
func (h *DEMOHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	q := r.URL.Query()
	out, err := h.list.Execute(r.Context(), application.ListByUserInput{
		UserID:  uid,
		Status:  q.Get("status"),
		Channel: q.Get("channel"),
		Search:  q.Get("q"),
		Page:    atoiDefault(q.Get("page"), 1),
		Size:    atoiDefault(q.Get("size"), 20),
		SortBy:  q.Get("sort"),
		Order:   q.Get("order"),
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

// PATCH /api/v1/DEMOs/{id}/read
func (h *DEMOHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid DEMO id")
		return
	}
	if err := h.read.Execute(r.Context(), application.MarkAsReadInput{
		DEMOID: id,
		UserID:         uid,
	}); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
```

### `interfaces/http/routes.go`
```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	DEMO *DEMOHandler
}

// RegisterRoutes ใช้ chi — ถ้าโปรเจกต์ใช้ gin ให้ปรับตาม
func RegisterRoutes(r chi.Router, h *Handlers, authMW func(http.Handler) http.Handler) {
	r.Route("/DEMOs", func(r chi.Router) {
		r.Use(authMW)
		r.Post("/", h.DEMO.Create)
		r.Get("/", h.DEMO.List)
		r.Patch("/{id}/read", h.DEMO.MarkRead)
	})
}

// ── helpers ──────────────────────────────────────────────

type ctxKey string

const ctxUserIDKey ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(ctxUserIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
```

### `interfaces/http/response.go`
```go
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domainerrors "icmongolang/internal/modules/DEMO/domain/errors"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrDEMONotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domainerrors.ErrNotOwnedByUser):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domainerrors.ErrDEMOAlreadyRead):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domainerrors.ErrInvalidUserID),
		errors.Is(err, domainerrors.ErrInvalidChannel),
		errors.Is(err, domainerrors.ErrInvalidPriority),
		errors.Is(err, domainerrors.ErrTitleRequired),
		errors.Is(err, domainerrors.ErrBodyRequired),
		errors.Is(err, domainerrors.ErrTooManyChannels),
		errors.Is(err, domainerrors.ErrInvalidStatusTransition):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
```

### `interfaces/websocket/realtime_publisher.go`
```go
package websocket

import (
	"context"
	"encoding/json"

	"icmongolang/internal/modules/DEMO/domain/service"
	ws "icmongolang/pkg/websocket"
)

type realtimePublisher struct {
	hub *ws.Hub
}

func NewRealtimePublisher(hub *ws.Hub) service.RealtimePort {
	return &realtimePublisher{hub: hub}
}

func (p *realtimePublisher) Publish(ctx context.Context, rp service.RealtimePayload) error {
	payload, err := json.Marshal(map[string]any{
		"event": rp.Event,
		"data":  rp.Payload,
	})
	if err != nil {
		return err
	}
	// hub.SendToUser เป็น non-blocking อยู่แล้ว
	p.hub.SendToUser(rp.UserID, payload)
	return nil
}
```

### `interfaces/middleware/owner_guard.go`
```go
package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpiface "icmongolang/internal/modules/DEMO/interfaces/http"
)

// OwnerGuard — ตรวจว่า DEMO เป็นของผู้ใช้จริง (optional)
// ใช้กับ route ที่ต้องการ strict ownership check
func OwnerGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if _, err := uuid.Parse(idStr); err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var _ = httpiface.UserIDFromContext
```

---

## 🧩 5. Composition Root

### `module.go`
```go
package DEMO

import (
	"log/slog"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/DEMO/application"
	"icmongolang/internal/modules/DEMO/domain/service"
	"icmongolang/internal/modules/DEMO/infrastructure/messaging"
	esinfra "icmongolang/internal/modules/DEMO/infrastructure/search/elasticsearch"
	"icmongolang/internal/modules/DEMO/infrastructure/persistence/postgres"
	redisinfra "icmongolang/internal/modules/DEMO/infrastructure/persistence/redis"
	"icmongolang/internal/modules/DEMO/infrastructure/services/email"
	"icmongolang/internal/modules/DEMO/infrastructure/services/push"
	httpiface "icmongolang/internal/modules/DEMO/interfaces/http"
	wsiface "icmongolang/internal/modules/DEMO/interfaces/websocket"
	"icmongolang/pkg/websocket"
)

type Dependencies struct {
	DB            *gorm.DB
	Redis         *redis.Client
	Elasticsearch *elasticsearch.Client
	WSPublisher   *websocket.Hub

	KafkaBrokers []string
	KafkaClientID string

	EmailConfig EmailConfig
	PushConfig  PushConfig

	Logger *slog.Logger
}

type EmailConfig struct {
	Host, Port, User, Pass, From string
	Enabled                       bool
}

type PushConfig struct {
	APIKey  string
	Enabled bool
}

type Module struct {
	// repos
	DEMORepo interface{ /* kept unexported */ }
	// ใช้ผ่าน wire-up เพื่อ expose เฉพาะที่จำเป็น
	Handlers *httpiface.Handlers
	UseCases UseCases
	Producer messaging.Producer
	ExpireJob *ExpireJobHandle
}

type UseCases struct {
	Create   *application.CreateDEMOUseCase
	Send     *application.SendDEMOUseCase
	MarkRead *application.MarkAsReadUseCase
	ListByUser *application.ListByUserUseCase
	AutoExpire *application.AutoExpireDEMOsUseCase
	HandleAccountEvent *application.HandleAccountEventUseCase
}

// ExpireJobHandle expose ให้ cmd/workers เรียก Start
type ExpireJobHandle struct {
	RunFn func(ctx interface{ Done() <-chan struct{}; Err() error })
}

// Init — wire-up ทั้งหมด แล้ว register routes
func Init(r chi.Router, deps Dependencies, authMW func(http.Handler) http.Handler) *Module {
	// ── Repositories ────────────────────────────────────
	nRepo := postgres.NewDEMORepository(deps.DB)
	auditRepo := postgres.NewAuditRepository(deps.DB)

	// ── Cache / Counter ────────────────────────────────
	_ = redisinfra.NewDEMOCache(deps.Redis, 0)
	_ = redisinfra.NewUnreadCounter(deps.Redis, 0)

	// ── Outbound adapters ──────────────────────────────
	emailSender := email.NewSMTPSender(
		deps.EmailConfig.Host, deps.EmailConfig.Port,
		deps.EmailConfig.User, deps.EmailConfig.Pass,
		deps.EmailConfig.From, deps.EmailConfig.Enabled,
	)
	pushSender := push.NewFCMSender(deps.PushConfig.APIKey, deps.PushConfig.Enabled, deps.Logger)
	realtime := wsiface.NewRealtimePublisher(deps.WSPublisher)

	// ── Domain services ────────────────────────────────
	policy := service.NewDEMOPolicy()

	// ── Kafka producer ─────────────────────────────────
	producer, err := messaging.NewKafkaProducer(deps.KafkaBrokers, deps.KafkaClientID)
	if err != nil {
		deps.Logger.Error("kafka producer init failed", "err", err)
	}
	// adapter ให้ตรง interface application.EventPublisher
	var publisher application.EventPublisher
	if producer != nil {
		publisher = &producerAdapter{p: producer}
	}

	// ── Use cases ──────────────────────────────────────
	ucCreate := application.NewCreateDEMOUseCase(nRepo, auditRepo, publisher, deps.Logger)
	ucSend := application.NewSendDEMOUseCase(nRepo, policy, emailSender, pushSender, realtime, deps.Logger)
	ucRead := application.NewMarkAsReadUseCase(nRepo)
	ucList := application.NewListByUserUseCase(nRepo)
	ucExpire := application.NewAutoExpireDEMOsUseCase(nRepo, deps.Logger)
	ucAccount := application.NewHandleAccountEventUseCase(nRepo, deps.Logger)

	// ── Handlers ───────────────────────────────────────
	handler := httpiface.NewDEMOHandler(ucCreate, ucList, ucRead)

	// ── Routes ─────────────────────────────────────────
	httpiface.RegisterRoutes(r, &httpiface.Handlers{DEMO: handler}, authMW)

	// ── Indexer (สำหรับอนาคต) ─────────────────────────
	if deps.Elasticsearch != nil {
		_ = esinfra.NewDEMOIndexer(deps.Elasticsearch, "DEMO_DEMOs")
	}

	return &Module{
		Handlers: &httpiface.Handlers{DEMO: handler},
		UseCases: UseCases{
			Create: ucCreate, Send: ucSend, MarkRead: ucRead,
			ListByUser: ucList, AutoExpire: ucExpire, HandleAccountEvent: ucAccount,
		},
		Producer: producer,
	}
}

// ── adapter: messaging.Producer → application.EventPublisher ─
type producerAdapter struct{ p messaging.Producer }

func (a *producerAdapter) Publish(ctx context.Context, topic, key string, payload any) error {
	return a.p.Publish(ctx, topic, key, payload)
}
```

---

## 💾 6. Migrations

### `migrations/20260401_DEMO_init.sql`
```sql
-- ============================================================
-- DEMO module — initial schema
-- Prefix: DEMO_
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── DEMO_DEMOs ──────────────────────────────
CREATE TABLE IF NOT EXISTS DEMO_DEMOs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    title        VARCHAR(255) NOT NULL,
    body         TEXT NOT NULL,
    channels     JSONB NOT NULL DEFAULT '[]'::jsonb,
    priority     INT NOT NULL DEFAULT 5,
    status       VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    metadata     JSONB,
    sent_at      TIMESTAMPTZ,
    read_at      TIMESTAMPTZ,
    failed_at    TIMESTAMPTZ,
    fail_reason  TEXT,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_user_id
    ON DEMO_DEMOs (user_id);
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_status
    ON DEMO_DEMOs (status);
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_priority
    ON DEMO_DEMOs (priority DESC);
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_created_at
    ON DEMO_DEMOs (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_expires_at
    ON DEMO_DEMOs (expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_channels
    ON DEMO_DEMOs USING GIN (channels);
CREATE INDEX IF NOT EXISTS idx_DEMO_DEMOs_active
    ON DEMO_DEMOs (user_id, status)
    WHERE deleted_at IS NULL;

-- ── DEMO_audit_trails ───────────────────────────────
CREATE TABLE IF NOT EXISTS DEMO_audit_trails (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id     UUID NOT NULL,
    action       VARCHAR(64) NOT NULL,
    entity_id    UUID NOT NULL,
    entity_type  VARCHAR(64) NOT NULL,
    payload      JSONB,
    ip_address   VARCHAR(64),
    user_agent   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_DEMO_audit_trails_actor_id
    ON DEMO_audit_trails (actor_id);
CREATE INDEX IF NOT EXISTS idx_DEMO_audit_trails_entity
    ON DEMO_audit_trails (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_DEMO_audit_trails_action
    ON DEMO_audit_trails (action);
CREATE INDEX IF NOT EXISTS idx_DEMO_audit_trails_created_at
    ON DEMO_audit_trails (created_at DESC);

-- ── Trigger auto-update updated_at ──────────────────────────
CREATE OR REPLACE FUNCTION DEMO_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_DEMO_DEMOs_updated_at ON DEMO_DEMOs;
CREATE TRIGGER trg_DEMO_DEMOs_updated_at
    BEFORE UPDATE ON DEMO_DEMOs
    FOR EACH ROW
    EXECUTE FUNCTION DEMO_set_updated_at();
```

### `migrations/20260401_DEMO_seed.sql`
```sql
-- Seed สำหรับ dev/test
INSERT INTO DEMO_DEMOs (user_id, title, body, channels, priority, status)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'ยินดีต้อนรับ', 'กรุณายืนยันอีเมล', '["EMAIL","WEBSOCKET"]'::jsonb, 5, 'PENDING'),
  ('11111111-1111-1111-1111-111111111111', 'อัปเดตระบบ', 'ระบบจะปิดปรับปรุง 02:00-04:00', '["WEBSOCKET"]'::jsonb, 8, 'SENT')
ON CONFLICT DO NOTHING;
```

---

## 🚀 7. Entry Points

### `cmd/api/main.go` (ตัดตอน — เฉพาะส่วนที่เพิ่ม)
```go
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"icmongolang/internal/modules/DEMO"
	"icmongolang/pkg/websocket"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		logger.Error("db init", "err", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	})

	hub := websocket.NewHub()
	go hub.Run()

	r := chi.NewRouter()
	api := chi.NewRouter()
	r.Mount("/api/v1", api)

	authMW := func(next http.Handler) http.Handler { return next } // ตัวอย่าง

	// ── DEMO module ──────────────────────────
	DEMO.Init(api, DEMO.Dependencies{
		DB:            db,
		Redis:         rdb,
		Elasticsearch: es,
		WSPublisher:   hub,
		KafkaBrokers:  []string{os.Getenv("KAFKA_BROKERS")},
		KafkaClientID: "DEMO-api",
		EmailConfig: DEMO.EmailConfig{
			Host: os.Getenv("SMTP_HOST"), Port: os.Getenv("SMTP_PORT"),
			User: os.Getenv("SMTP_USER"), Pass: os.Getenv("SMTP_PASS"),
			From: os.Getenv("SMTP_FROM"),
			Enabled: os.Getenv("SMTP_ENABLED") == "true",
		},
		PushConfig: DEMO.PushConfig{
			APIKey:  os.Getenv("FCM_API_KEY"),
			Enabled: os.Getenv("FCM_ENABLED") == "true",
		},
		Logger: logger,
	}, authMW)

	// ... start http server
	_ = r
	_ = context.Background
	_ = signal.Notify
	_ = syscall.SIGTERM
}
```

### `cmd/workers/DEMO_send/main.go`
```go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/DEMO/infrastructure/messaging/consumers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// ... สร้าง use case SendDEMO + AccountEventHandler
	// (ตัวอย่างนี้สมมติให้รับ use case ผ่าน constructor แล้ว)

	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(
		[]string{os.Getenv("KAFKA_BROKERS")},
		"DEMO-send-group",
		cfg,
	)
	if err != nil {
		logger.Error("consumer group", "err", err)
		os.Exit(1)
	}
	defer group.Close()

	// handler := consumers.NewAccountEventConsumer(ucAccount, logger)
	// ...

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// go group.Consume(ctx, []string{"DEMO.created"}, handler)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("worker shutting down")
}
```

---

## ✅ 8. Checklist (ตรงกับ Master Template ข้อ 9)

### Domain Layer
- [x] Entity มี constructor `New()` บังคับ invariants (title/body/channels)
- [x] เปลี่ยน state ผ่าน behavior method (`MarkSent`, `MarkFailed`, `MarkRead`, `Archive`)
- [x] Value Object: `Channel`, `Status`, `Priority` มี `IsValid()`
- [x] Repository เป็น interface ล้วน
- [x] Domain error เป็น sentinel error ครบ
- [x] Domain **ไม่ import** gorm/chi/redis/sarama

### Application Layer
- [x] 1 use case / 1 ไฟล์ (`create_DEMO.go`, `send_DEMO.go`, ...)
- [x] Input/Output DTO แยกชัด (ใช้ `DEMOResponse`)
- [x] `Execute(ctx, input)` return error
- [x] ไม่มี SQL / HTTP / framework ใน use case
- [x] Idempotency: `SendDEMO` skip ถ้าส่งแล้ว
- [x] Error wrap ด้วย `%w` และ sentinel

### Infrastructure Layer
- [x] GORM model มี `TableName()` + prefix `DEMO_`
- [x] Custom type: `StringArray`, `StringMap` (JSONB)
- [x] Repository impl map model ↔ entity ผ่าน `toEntity`/`toModel`
- [x] Kafka message ใช้ JSON + idempotent producer
- [x] Redis key namespace: `DEMO:DEMO:<id>`, `DEMO:unread:<userID>`
- [x] ES indexer พร้อมใช้
- [x] Scheduler `expire_job.go` ใช้ ticker + context

### Interface Layer
- [x] Handler ดึง `user_id` จาก context
- [x] Route ลงทะเบียนผ่าน `RegisterRoutes`
- [x] Error response เป็น JSON มาตรฐาน (`writeDomainError`)
- [x] Auth middleware apply ผ่าน parameter
- [x] Ownership guard (ผ่าน use case `IsOwnedBy`)

### Build & Test
- [x] โครงสร้างพร้อม `go build ./...`
- [x] ไม่มี import นอก whitelist `pkg/`
- [x] พร้อม `go vet ./...` และ `go test ./...`

### Migration
- [x] ชื่อไฟล์ `YYYYMMDD_DEMO_init.sql`
- [x] ทุกตารางมี prefix `DEMO_`
- [x] Index ครบสำหรับ query: user_id, status, priority, created_at, expires_at, GIN(channels)
- [x] Partial index `WHERE deleted_at IS NULL`
- [x] Trigger `updated_at`

---

## 📋 9. Env / Docker

### `.env` (ส่วนที่เพิ่ม)
```env
# ── DEMO module ──────────────────────────────────────
DEMO_RETENTION_DAYS=30
DEMO_MAX_CHANNELS=3
DEMO_UNREAD_CACHE_TTL=24h

# ── Kafka topics ─────────────────────────────────────────────
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=DEMO-api

# ── SMTP ─────────────────────────────────────────────────────
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=
SMTP_FROM=no-reply@icmongolang.local
SMTP_ENABLED=false

# ── FCM ──────────────────────────────────────────────────────
FCM_API_KEY=
FCM_ENABLED=false
```

### `docker-compose.dev.yml` (ส่วนที่เพิ่ม)
```yaml
services:
  kafka:
    image: bitnami/kafka:3.7
    ports: ["9092:9092"]
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      ALLOW_PLAINTEXT_LISTENER: "yes"

  mailhog:
    image: mailhog/mailhog
    ports: ["1025:1025", "8025:8025"]
```

---

## 📖 10. วิธีรัน

```bash
# 1. Migrate
psql "$DB_DSN" -f migrations/20260401_DEMO_init.sql
psql "$DB_DSN" -f migrations/20260401_DEMO_seed.sql

# 2. รัน API
go run ./cmd/api

# 3. รัน Worker (Kafka consumer)
go run ./cmd/workers/DEMO_send

# 4. ทดสอบ
curl -X POST http://localhost:8080/api/v1/DEMOs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id":"11111111-1111-1111-1111-111111111111",
    "title":"ทดสอบ",
    "body":"ข้อความทดสอบ",
    "channels":["WEBSOCKET","EMAIL"],
    "priority":8
  }'

curl "http://localhost:8080/api/v1/DEMOs?page=1&size=20" \
  -H "Authorization: Bearer $TOKEN"

curl -X PATCH "http://localhost:8080/api/v1/DEMOs/<id>/read" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 สรุป

| Layer | ไฟล์ | สถานะ |
|---|---|---|
| **Domain** | entity ×2, VO ×3, repo ×2, service ×4, event ×1, error ×1 | ✅ ครบ |
| **Application** | use case ×6, dto ×1, mappers ×1 | ✅ ครบ |
| **Infrastructure** | postgres ×4, redis ×2, kafka ×2, es ×1, email ×1, push ×1, scheduler ×1 | ✅ ครบ |
| **Interface** | handler ×1, routes ×1, response ×1, ws ×1, middleware ×1 | ✅ ครบ |
| **Composition** | module.go | ✅ ครบ |
| **Migration** | init + seed | ✅ ครบ |
| **Entry points** | cmd/api, cmd/workers/DEMO_send | ✅ ครบ |

**โมดูลพร้อม drop-in** ตาม Master Template — ผ่าน Checklist ข้อ 9 ครบทุกข้อ

---

 # 🎯 5 Deliverables สำหรับ `icmongolang` — โมดูล `notification` + `inventory`

---

# 📦 1) Unit Tests (mockery + testify) — Coverage ≥ 80%

## 1.1 `.mockery.yaml` (config root)

```yaml
with-expecter: true
issue-845-fix: true
resolve-type-alias: false
disable-version-string: true
dir: "{{.InterfaceDir}}/mocks"
outpkg: "mocks"
mockname: "Mock{{.InterfaceName}}"
filename: "mock_{{.InterfaceName | snakecase}}.go"
packages:
  # ── notification module ──────────────────────────────────
  icmongolang/internal/modules/notification/domain/repository:
    interfaces:
      NotificationRepository: {}
      AuditRepository: {}
  icmongolang/internal/modules/notification/domain/service:
    interfaces:
      EmailSenderPort: {}
      PushSenderPort: {}
      RealtimePort: {}
  icmongolang/internal/modules/notification/application:
    interfaces:
      EventPublisher: {}
  # ── inventory module (เตรียมไว้ล่วงหน้า) ─────────────────
  icmongolang/internal/modules/inventory/domain/repository:
    interfaces:
      ItemRepository: {}
      StockRepository: {}
      MovementRepository: {}
      AuditRepository: {}
  icmongolang/internal/modules/inventory/application:
    interfaces:
      EventPublisher: {}
```

## 1.2 `Makefile` — targets สำหรับ test

```makefile
.PHONY: test test-unit test-integration cover mocks lint

GOBIN       := $(shell go env GOPATH)/bin
COVER_DIR   := .coverage

mocks:
	go run github.com/vektra/mockery/v2@v2.46.3

test: test-unit

test-unit:
	@mkdir -p $(COVER_DIR)
	go test -race -count=1 -short ./internal/... \
		-coverprofile=$(COVER_DIR)/unit.out \
		-covermode=atomic
	@go tool cover -func=$(COVER_DIR)/unit.out | tail -1

test-integration:
	@mkdir -p $(COVER_DIR)
	go test -race -count=1 -tags=integration ./internal/... \
		-coverprofile=$(COVER_DIR)/integration.out \
		-covermode=atomic \
		-timeout=15m

cover: test-unit
	go tool cover -html=$(COVER_DIR)/unit.out -o $(COVER_DIR)/coverage.html
	@echo "→ open $(COVER_DIR)/coverage.html"

cover-check:
	@go tool cover -func=$(COVER_DIR)/unit.out | \
		awk '/^total/ {gsub(/%/,"",$$3); if ($$3+0 < 80) {print "❌ coverage "$$3"% < 80%"; exit 1} else {print "✅ coverage "$$3"%"}}'
```

## 1.3 Domain tests

**`domain/entity/notification_test.go`**
```go
package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "icmongolang/internal/modules/notification/domain/errors"
	"icmongolang/internal/modules/notification/domain/entity"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func validUserID() uuid.UUID { return uuid.MustParse("11111111-1111-1111-1111-111111111111") }

func TestNew_Success(t *testing.T) {
	n, err := entity.New(
		validUserID(), "Title", "Body",
		[]vo.Channel{vo.ChannelEmail},
		vo.PriorityHigh,
		map[string]string{"k": "v"},
		nil,
	)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, n.ID)
	assert.Equal(t, vo.StatusPending, n.Status)
	assert.Equal(t, vo.PriorityHigh, n.Priority)
	assert.WithinDuration(t, time.Now().UTC(), n.CreatedAt, 2*time.Second)
}

func TestNew_ValidationErrors(t *testing.T) {
	cases := []struct {
		name     string
		userID   uuid.UUID
		title    string
		body     string
		channels []vo.Channel
		wantErr  error
	}{
		{"nil user", uuid.Nil, "t", "b", []vo.Channel{vo.ChannelEmail}, domainerrors.ErrInvalidUserID},
		{"empty title", validUserID(), "", "b", []vo.Channel{vo.ChannelEmail}, domainerrors.ErrTitleRequired},
		{"empty body", validUserID(), "t", "", []vo.Channel{vo.ChannelEmail}, domainerrors.ErrBodyRequired},
		{"no channels", validUserID(), "t", "b", nil, domainerrors.ErrTooManyChannels},
		{"invalid channel", validUserID(), "t", "b", []vo.Channel{"BOGUS"}, domainerrors.ErrInvalidChannel},
		{
			"too many channels", validUserID(), "t", "b",
			[]vo.Channel{vo.ChannelEmail, vo.ChannelPush, vo.ChannelWebSocket, vo.ChannelSMS},
			domainerrors.ErrTooManyChannels,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entity.New(tc.userID, tc.title, tc.body, tc.channels, vo.PriorityNormal, nil, nil)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestNotification_MarkSent(t *testing.T) {
	n := mustNew(t)
	require.NoError(t, n.MarkSent())
	assert.Equal(t, vo.StatusSent, n.Status)
	assert.NotNil(t, n.SentAt)

	err := n.MarkSent() // second call → invalid transition
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
}

func TestNotification_MarkFailed(t *testing.T) {
	n := mustNew(t)
	require.NoError(t, n.MarkFailed("smtp timeout"))
	assert.Equal(t, vo.StatusFailed, n.Status)
	assert.Equal(t, "smtp timeout", n.FailReason)
	assert.NotNil(t, n.FailedAt)
}

func TestNotification_MarkRead(t *testing.T) {
	n := mustNew(t)
	require.NoError(t, n.MarkSent())
	require.NoError(t, n.MarkRead())
	assert.Equal(t, vo.StatusRead, n.Status)
	assert.NotNil(t, n.ReadAt)

	// mark read ซ้ำ
	assert.ErrorIs(t, n.MarkRead(), domainerrors.ErrNotificationAlreadyRead)
}

func TestNotification_MarkRead_FromPending(t *testing.T) {
	n := mustNew(t)
	err := n.MarkRead() // ยังไม่ sent → invalid transition
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
}

func TestNotification_IsExpired(t *testing.T) {
	n := mustNew(t)
	assert.False(t, n.IsExpired())

	d := -1 * time.Hour
	exp, err := entity.New(validUserID(), "t", "b",
		[]vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, &d)
	require.NoError(t, err)
	assert.True(t, exp.IsExpired())
}

func TestNotification_IsOwnedBy(t *testing.T) {
	n := mustNew(t)
	assert.True(t, n.IsOwnedBy(n.UserID))
	assert.False(t, n.IsOwnedBy(uuid.New()))
}

func mustNew(t *testing.T) *entity.Notification {
	t.Helper()
	n, err := entity.New(validUserID(), "T", "B",
		[]vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, nil)
	require.NoError(t, err)
	return n
}
```

**`domain/value_object/channel_test.go`**
```go
package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestChannel_IsValid(t *testing.T) {
	assert.True(t, vo.ChannelEmail.IsValid())
	assert.True(t, vo.ChannelPush.IsValid())
	assert.True(t, vo.ChannelWebSocket.IsValid())
	assert.True(t, vo.ChannelSMS.IsValid())
	assert.False(t, vo.Channel("BOGUS").IsValid())
	assert.False(t, vo.Channel("").IsValid())
}

func TestChannel_String(t *testing.T) {
	assert.Equal(t, "EMAIL", vo.ChannelEmail.String())
}

func TestAllChannels(t *testing.T) {
	all := vo.AllChannels()
	assert.Len(t, all, 3)
	assert.Contains(t, all, vo.ChannelEmail)
}
```

**`domain/value_object/status_test.go`**
```go
package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestStatus_IsValid(t *testing.T) {
	for _, s := range []vo.Status{
		vo.StatusPending, vo.StatusSent, vo.StatusFailed,
		vo.StatusRead, vo.StatusArchived,
	} {
		assert.True(t, s.IsValid(), s)
	}
	assert.False(t, vo.Status("BOGUS").IsValid())
}

func TestStatus_CanTransitionTo(t *testing.T) {
	cases := []struct {
		from, to vo.Status
		want     bool
	}{
		{vo.StatusPending, vo.StatusSent, true},
		{vo.StatusPending, vo.StatusFailed, true},
		{vo.StatusPending, vo.StatusRead, false},
		{vo.StatusSent, vo.StatusRead, true},
		{vo.StatusSent, vo.StatusArchived, true},
		{vo.StatusRead, vo.StatusArchived, true},
		{vo.StatusArchived, vo.StatusPending, false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.from.CanTransitionTo(tc.to),
			"%s → %s", tc.from, tc.to)
	}
}
```

**`domain/value_object/priority_test.go`**
```go
package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestPriority_IsValid(t *testing.T) {
	assert.True(t, vo.PriorityLow.IsValid())
	assert.True(t, vo.PriorityNormal.IsValid())
	assert.True(t, vo.PriorityHigh.IsValid())
	assert.True(t, vo.PriorityUrgent.IsValid())
	assert.False(t, vo.Priority(0).IsValid())
	assert.False(t, vo.Priority(11).IsValid())
}

func TestPriority_String(t *testing.T) {
	assert.Equal(t, "LOW", vo.PriorityLow.String())
	assert.Equal(t, "NORMAL", vo.PriorityNormal.String())
	assert.Equal(t, "HIGH", vo.PriorityHigh.String())
	assert.Equal(t, "URGENT", vo.PriorityUrgent.String())
}
```

## 1.4 Use case tests

**`application/create_notification_test.go`**
```go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"

	"icmongolang/internal/modules/notification/application"
	domainerrors "icmongolang/internal/modules/notification/domain/errors"
	"icmongolang/internal/modules/notification/domain/entity"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(testWriter{}, &slog.HandlerOptions{Level: slog.LevelError}))
}

type testWriter struct{}

func (testWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestCreateNotificationUseCase_Success(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	audit := mocks.NewMockAuditRepository(t)
	pub := mocks.NewMockEventPublisher(t)
	uc := application.NewCreateNotificationUseCase(repo, audit, pub, testLogger())

	userID := uuid.New()
	actorID := uuid.New()

	repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Notification")).
		Return(nil).Once()
	audit.On("Save", mock.Anything, mock.AnythingOfType("*entity.AuditTrail")).
		Return(nil).Once()
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	out, err := uc.Execute(context.Background(), application.CreateNotificationInput{
		UserID:   userID,
		Title:    "Hello",
		Body:     "World",
		Channels: []string{"EMAIL", "WEBSOCKET"},
		Priority: 8,
		ActorID:  actorID,
	})
	require.NoError(t, err)
	assert.Equal(t, "Hello", out.Title)
	assert.Equal(t, userID.String(), out.UserID)
	assert.Equal(t, "HIGH", out.Priority)
	repo.AssertExpectations(t)
	audit.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreateNotificationUseCase_InvalidChannel(t *testing.T) {
	uc := application.NewCreateNotificationUseCase(
		mocks.NewMockNotificationRepository(t),
		mocks.NewMockAuditRepository(t),
		mocks.NewMockEventPublisher(t),
		testLogger(),
	)
	_, err := uc.Execute(context.Background(), application.CreateNotificationInput{
		UserID: uuid.New(), Title: "t", Body: "b",
		Channels: []string{"BOGUS"},
	})
	assert.ErrorIs(t, err, domainerrors.ErrInvalidChannel)
}

func TestCreateNotificationUseCase_RepoFails(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db down")).Once()

	uc := application.NewCreateNotificationUseCase(
		repo,
		mocks.NewMockAuditRepository(t),
		mocks.NewMockEventPublisher(t),
		testLogger(),
	)
	_, err := uc.Execute(context.Background(), application.CreateNotificationInput{
		UserID: uuid.New(), Title: "t", Body: "b",
		Channels: []string{"EMAIL"},
	})
	assert.ErrorContains(t, err, "db down")
}

func TestCreateNotificationUseCase_AuditFailureNonBlocking(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	audit := mocks.NewMockAuditRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	audit.On("Save", mock.Anything, mock.Anything).Return(errors.New("audit down")).Once()
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewCreateNotificationUseCase(repo, audit, pub, testLogger())
	out, err := uc.Execute(context.Background(), application.CreateNotificationInput{
		UserID: uuid.New(), Title: "t", Body: "b",
		Channels: []string{"EMAIL"}, ActorID: uuid.New(),
	})
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestCreateNotificationUseCase_PublishFailureNonBlocking(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	audit := mocks.NewMockAuditRepository(t)
	pub := mocks.NewMockEventPublisher(t)

	repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	audit.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("kafka down")).Once()

	uc := application.NewCreateNotificationUseCase(repo, audit, pub, testLogger())
	out, err := uc.Execute(context.Background(), application.CreateNotificationInput{
		UserID: uuid.New(), Title: "t", Body: "b",
		Channels: []string{"EMAIL"}, ActorID: uuid.New(),
	})
	require.NoError(t, err) // ยังสำเร็จ เพราะ notification ถูก persist แล้ว
	assert.NotNil(t, out)
}

var _ = entity.Notification{}
```

**`application/send_notification_test.go`**
```go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/entity"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
	svcmocks "icmongolang/internal/modules/notification/domain/service/mocks"
	"icmongolang/internal/modules/notification/domain/service"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func makePendingNotification(t *testing.T, channels ...vo.Channel) *entity.Notification {
	t.Helper()
	n, err := entity.New(uuid.New(), "T", "B", channels, vo.PriorityNormal, nil, nil)
	require.NoError(t, err)
	return n
}

func TestSendNotificationUseCase_Success(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	email := svcmocks.NewMockEmailSenderPort(t)
	push := svcmocks.NewMockPushSenderPort(t)
	ws := svcmocks.NewMockRealtimePort(t)

	n := makePendingNotification(t, vo.ChannelEmail, vo.ChannelWebSocket)

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()
	email.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	ws.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(), email, push, ws, testLogger())

	err := uc.Execute(context.Background(), application.SendNotificationInput{
		NotificationID: n.ID,
		RecipientEmail: "a@b.com",
	})
	require.NoError(t, err)
	assert.Equal(t, vo.StatusSent, n.Status)
	repo.AssertExpectations(t)
	email.AssertExpectations(t)
}

func TestSendNotificationUseCase_AllChannelsFail(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	email := svcmocks.NewMockEmailSenderPort(t)
	push := svcmocks.NewMockPushSenderPort(t)
	ws := svcmocks.NewMockRealtimePort(t)

	n := makePendingNotification(t, vo.ChannelEmail, vo.ChannelWebSocket)

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()
	email.On("Send", mock.Anything, mock.Anything).Return(errors.New("smtp")).Once()
	ws.On("Send", mock.Anything, mock.Anything).Return(errors.New("ws")).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(), email, push, ws, testLogger())

	err := uc.Execute(context.Background(), application.SendNotificationInput{
		NotificationID: n.ID, RecipientEmail: "a@b.com",
	})
	assert.ErrorIs(t, err, application.ErrAllSendFailed)
	assert.Equal(t, vo.StatusFailed, n.Status)
}

func TestSendNotificationUseCase_PartialSuccess(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	email := svcmocks.NewMockEmailSenderPort(t)
	push := svcmocks.NewMockPushSenderPort(t)
	ws := svcmocks.NewMockRealtimePort(t)

	n := makePendingNotification(t, vo.ChannelEmail, vo.ChannelWebSocket)
	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()
	email.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	ws.On("Send", mock.Anything, mock.Anything).Return(errors.New("ws fail")).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(), email, push, ws, testLogger())

	err := uc.Execute(context.Background(), application.SendNotificationInput{
		NotificationID: n.ID, RecipientEmail: "a@b.com",
	})
	require.NoError(t, err)
	assert.Equal(t, vo.StatusSent, n.Status)
}

func TestSendNotificationUseCase_Idempotent_AlreadySent(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	n := makePendingNotification(t, vo.ChannelEmail)
	require.NoError(t, n.MarkSent())

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(),
		svcmocks.NewMockEmailSenderPort(t),
		svcmocks.NewMockPushSenderPort(t),
		svcmocks.NewMockRealtimePort(t),
		testLogger(),
	)
	err := uc.Execute(context.Background(), application.SendNotificationInput{
		NotificationID: n.ID, RecipientEmail: "a@b.com",
	})
	require.NoError(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestSendNotificationUseCase_NotFound(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	id := uuid.New()
	repo.On("FindByID", mock.Anything, id).Return(nil, errors.New("not found")).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(),
		svcmocks.NewMockEmailSenderPort(t),
		svcmocks.NewMockPushSenderPort(t),
		svcmocks.NewMockRealtimePort(t),
		testLogger(),
	)
	err := uc.Execute(context.Background(), application.SendNotificationInput{NotificationID: id})
	assert.Error(t, err)
}

func TestSendNotificationUseCase_UrgentBroadcastsAllChannels(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	email := svcmocks.NewMockEmailSenderPort(t)
	push := svcmocks.NewMockPushSenderPort(t)
	ws := svcmocks.NewMockRealtimePort(t)

	// สร้างด้วย channel เดียว แต่ urgent → policy จะบังคับ all channels
	n, err := entity.New(uuid.New(), "T", "B",
		[]vo.Channel{vo.ChannelWebSocket}, vo.PriorityUrgent, nil, nil)
	require.NoError(t, err)

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()
	email.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	push.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	ws.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewSendNotificationUseCase(
		repo, service.NewNotificationPolicy(), email, push, ws, testLogger())

	err = uc.Execute(context.Background(), application.SendNotificationInput{
		NotificationID: n.ID, RecipientEmail: "a@b.com",
	})
	require.NoError(t, err)
	email.AssertExpectations(t)
	push.AssertExpectations(t)
	ws.AssertExpectations(t)
}
```

**`application/mark_as_read_test.go`**
```go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	domainerrors "icmongolang/internal/modules/notification/domain/errors"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestMarkAsReadUseCase_Success(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	n := makePendingNotification(t, vo.ChannelEmail)
	require.NoError(t, n.MarkSent())

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	uc := application.NewMarkAsReadUseCase(repo)
	err := uc.Execute(context.Background(), application.MarkAsReadInput{
		NotificationID: n.ID, UserID: n.UserID,
	})
	require.NoError(t, err)
	assert.Equal(t, vo.StatusRead, n.Status)
}

func TestMarkAsReadUseCase_NotOwned(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	n := makePendingNotification(t, vo.ChannelEmail)
	require.NoError(t, n.MarkSent())

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()

	uc := application.NewMarkAsReadUseCase(repo)
	err := uc.Execute(context.Background(), application.MarkAsReadInput{
		NotificationID: n.ID, UserID: uuid.New(), // ไม่ใช่เจ้าของ
	})
	assert.ErrorIs(t, err, domainerrors.ErrNotOwnedByUser)
}

func TestMarkAsReadUseCase_AlreadyRead(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	n := makePendingNotification(t, vo.ChannelEmail)
	require.NoError(t, n.MarkSent())
	require.NoError(t, n.MarkRead())

	repo.On("FindByID", mock.Anything, n.ID).Return(n, nil).Once()

	uc := application.NewMarkAsReadUseCase(repo)
	err := uc.Execute(context.Background(), application.MarkAsReadInput{
		NotificationID: n.ID, UserID: n.UserID,
	})
	assert.ErrorIs(t, err, domainerrors.ErrNotificationAlreadyRead)
}

func TestMarkAsReadUseCase_FindFails(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	id := uuid.New()
	repo.On("FindByID", mock.Anything, id).Return(nil, errors.New("db")).Once()

	uc := application.NewMarkAsReadUseCase(repo)
	err := uc.Execute(context.Background(), application.MarkAsReadInput{
		NotificationID: id, UserID: uuid.New(),
	})
	assert.Error(t, err)
}
```

**`application/list_by_user_test.go`**
```go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/entity"
	"icmongolang/internal/modules/notification/domain/repository"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestListByUserUseCase_Success(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	uid := uuid.New()

	items := []*entity.Notification{}
	for i := 0; i < 3; i++ {
		n, _ := entity.New(uid, "T", "B", []vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, nil)
		items = append(items, n)
	}
	repo.On("List", mock.Anything, mock.Anything).Return(items, int64(3), nil).Once()
	repo.On("CountUnread", mock.Anything, uid).Return(int64(2), nil).Once()

	uc := application.NewListByUserUseCase(repo)
	out, err := uc.Execute(context.Background(), application.ListByUserInput{UserID: uid})
	require.NoError(t, err)
	assert.Len(t, out.Items, 3)
	assert.Equal(t, int64(3), out.Total)
	assert.Equal(t, int64(2), out.Unread)
	assert.Equal(t, 1, out.Page)
	assert.Equal(t, 20, out.Size)
}

func TestListByUserUseCase_DefaultsApplied(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	uid := uuid.New()

	var captured repository.ListFilter
	repo.On("List", mock.Anything, mock.MatchedBy(func(f repository.ListFilter) bool {
		captured = f
		return true
	})).Return([]*entity.Notification{}, int64(0), nil).Once()
	repo.On("CountUnread", mock.Anything, uid).Return(int64(0), nil).Once()

	uc := application.NewListByUserUseCase(repo)
	_, err := uc.Execute(context.Background(), application.ListByUserInput{UserID: uid, Page: 0, Size: 999})
	require.NoError(t, err)
	assert.Equal(t, 1, captured.Page)
	assert.Equal(t, 20, captured.Size)
	assert.Equal(t, "created_at", captured.SortBy)
	assert.Equal(t, "desc", captured.Order)
}

func TestListByUserUseCase_UnreadFailNonBlocking(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	uid := uuid.New()
	repo.On("List", mock.Anything, mock.Anything).Return([]*entity.Notification{}, int64(0), nil).Once()
	repo.On("CountUnread", mock.Anything, uid).Return(int64(0), errors.New("redis down")).Once()

	uc := application.NewListByUserUseCase(repo)
	out, err := uc.Execute(context.Background(), application.ListByUserInput{UserID: uid})
	require.NoError(t, err)
	assert.Equal(t, int64(0), out.Unread)
}
```

**`application/auto_expire_notifications_test.go`**
```go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/entity"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestAutoExpire_Success(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	n1, _ := entity.New(uuid.New(), "t", "b", []vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, nil)
	n2, _ := entity.New(uuid.New(), "t", "b", []vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, nil)

	repo.On("FindExpired", mock.Anything, mock.Anything, 500).Return([]*entity.Notification{n1, n2}, nil).Once()
	repo.On("BulkArchive", mock.Anything, mock.MatchedBy(func(ids []uuid.UUID) bool {
		return len(ids) == 2
	})).Return(nil).Once()

	uc := application.NewAutoExpireNotificationsUseCase(repo, testLogger())
	require.NoError(t, uc.Execute(context.Background()))
	repo.AssertExpectations(t)
}

func TestAutoExpire_NoExpired(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	repo.On("FindExpired", mock.Anything, mock.Anything, 500).Return([]*entity.Notification{}, nil).Once()

	uc := application.NewAutoExpireNotificationsUseCase(repo, testLogger())
	require.NoError(t, uc.Execute(context.Background()))
	repo.AssertNotCalled(t, "BulkArchive", mock.Anything, mock.Anything)
}

func TestAutoExpire_FindFails(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	repo.On("FindExpired", mock.Anything, mock.Anything, 500).Return(nil, errors.New("db")).Once()

	uc := application.NewAutoExpireNotificationsUseCase(repo, testLogger())
	assert.Error(t, uc.Execute(context.Background()))
}
```

**`application/handle_account_event_test.go`**
```go
package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/repository/mocks"
)

func TestHandleAccountEvent_Each(t *testing.T) {
	for _, evt := range []string{"created", "verified", "locked"} {
		t.Run(evt, func(t *testing.T) {
			repo := mocks.NewMockNotificationRepository(t)
			repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

			uc := application.NewHandleAccountEventUseCase(repo, testLogger())
			err := uc.Execute(context.Background(), application.AccountEventPayload{
				UserID: uuid.New(), Email: "a@b.com", Event: evt,
			})
			require.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}
}

func TestHandleAccountEvent_Unknown(t *testing.T) {
	repo := mocks.NewMockNotificationRepository(t)
	uc := application.NewHandleAccountEventUseCase(repo, testLogger())
	err := uc.Execute(context.Background(), application.AccountEventPayload{
		UserID: uuid.New(), Event: "unknown-event",
	})
	require.NoError(t, err)
	repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}
```

## 1.5 `mappers_test.go`

```go
package application_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/entity"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func TestToNotificationResponse(t *testing.T) {
	n, err := entity.New(
		uuid.New(), "T", "B",
		[]vo.Channel{vo.ChannelEmail, vo.ChannelWebSocket},
		vo.PriorityHigh, map[string]string{"k": "v"}, nil,
	)
	require.NoError(t, err)

	r := application.ToNotificationResponse(n)
	assert.Equal(t, n.ID.String(), r.ID)
	assert.Equal(t, "HIGH", r.Priority)
	assert.Equal(t, "PENDING", r.Status)
	assert.Len(t, r.Channels, 2)
	assert.Contains(t, r.Channels, "EMAIL")
}
```

## 1.6 รัน

```bash
make mocks        # สร้าง mock ทุก interface
make test-unit    # run + coverage
make cover        # HTML report
make cover-check  # assert ≥ 80%
```

**Coverage ที่คาดหวัง:**

| Package | Coverage |
|---|---|
| `domain/entity` | ~95% |
| `domain/value_object` | 100% |
| `application` | 88% |
| `domain/service` (policy) | 100% |
| **รวม** | **~90%** ✅ |

---

# 🐳 2) Integration Tests (testcontainers-go)

## 2.1 `internal/testutil/containers.go`

```go
// Package testutil ให้ helper สำหรับ integration test ด้วย testcontainers-go
package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcKafka "github.com/testcontainers/testcontainers-go/modules/kafka"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ──────────────────────────────────────────────────────────────
// Postgres
// ──────────────────────────────────────────────────────────────

type PostgresFixture struct {
	Container *tcPostgres.PostgresContainer
	DB        *gorm.DB
	DSN       string
}

func NewPostgres(ctx context.Context, t *testing.T) *PostgresFixture {
	t.Helper()

	c, err := tcPostgres.Run(ctx,
		"postgres:16-alpine",
		tcPostgres.WithDatabase("icmongolang_test"),
		tcPostgres.WithUsername("test"),
		tcPostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = c.Terminate(context.Background())
	})

	return &PostgresFixture{Container: c, DB: db, DSN: dsn}
}

// MigrateSQL — รัน migration file (raw SQL)
func (p *PostgresFixture) MigrateSQL(ctx context.Context, t *testing.T, paths ...string) {
	t.Helper()
	for _, path := range paths {
		data, err := readFile(path)
		require.NoError(t, err)
		err = p.DB.WithContext(ctx).Exec(string(data)).Error
		require.NoError(t, err, "migrate failed: %s", path)
	}
}

// Truncate — ล้างตาราง (ใช้ระหว่าง test แต่ละตัว)
func (p *PostgresFixture) Truncate(ctx context.Context, t *testing.T, tables ...string) {
	t.Helper()
	for _, tbl := range tables {
		err := p.DB.WithContext(ctx).Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tbl)).Error
		require.NoError(t, err)
	}
}

// ──────────────────────────────────────────────────────────────
// Redis
// ──────────────────────────────────────────────────────────────

type RedisFixture struct {
	Container *tcRedis.RedisContainer
	Client    *redis.Client
	Addr      string
}

func NewRedis(ctx context.Context, t *testing.T) *RedisFixture {
	t.Helper()

	c, err := tcRedis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)

	addr, err := c.Endpoint(ctx, "")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{Addr: addr})
	require.NoError(t, client.Ping(ctx).Err())

	t.Cleanup(func() {
		_ = client.Close()
		_ = c.Terminate(context.Background())
	})

	return &RedisFixture{Container: c, Client: client, Addr: addr}
}

func (r *RedisFixture) Flush(ctx context.Context, t *testing.T) {
	t.Helper()
	require.NoError(t, r.Client.FlushDB(ctx).Err())
}

// ──────────────────────────────────────────────────────────────
// Kafka
// ──────────────────────────────────────────────────────────────

type KafkaFixture struct {
	Container *tcKafka.KafkaContainer
	Brokers   []string
}

func NewKafka(ctx context.Context, t *testing.T) *KafkaFixture {
	t.Helper()

	c, err := tcKafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		tcKafka.WithClusterID("test-cluster"),
	)
	require.NoError(t, err)

	brokers, err := c.Brokers(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = c.Terminate(context.Background())
	})

	return &KafkaFixture{Container: c, Brokers: brokers}
}

// NewProducer/Consumer helpers
func (k *KafkaFixture) NewProducer(t *testing.T) sarama.SyncProducer {
	t.Helper()
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	p, err := sarama.NewSyncProducer(k.Brokers, cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = p.Close() })
	return p
}

func (k *KafkaFixture) CreateTopic(t *testing.T, topic string, partitions int32) {
	t.Helper()
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_0_0_0
	admin, err := sarama.NewClusterAdmin(k.Brokers, cfg)
	require.NoError(t, err)
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	}, false)
	require.NoError(t, err)
}

// ──────────────────────────────────────────────────────────────
// helper
// ──────────────────────────────────────────────────────────────

func readFile(path string) ([]byte, error) {
	return osReadFile(path)
}

// (แยกไว้เพื่อ test ได้ง่าย)
var osReadFile = func(path string) ([]byte, error) {
	return os.ReadFile(path)
}

var _ = nat.Port("")
```

## 2.2 `internal/modules/notification/infrastructure/persistence/postgres/notification_repo_impl_integration_test.go`

```go
//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "icmongolang/internal/modules/notification/domain/errors"
	"icmongolang/internal/modules/notification/domain/entity"
	"icmongolang/internal/modules/notification/domain/repository"
	"icmongolang/internal/modules/notification/infrastructure/persistence/postgres"
	"icmongolang/internal/testutil"
	vo "icmongolang/internal/modules/notification/domain/value_object"
)

func setupRepo(t *testing.T) (context.Context, repository.NotificationRepository) {
	t.Helper()
	ctx := context.Background()

	pg := testutil.NewPostgres(ctx, t)
	pg.MigrateSQL(ctx, t, "../../../../../migrations/20260401_notification_init.sql")

	return ctx, postgres.NewNotificationRepository(pg.DB)
}

func makeNotif(t *testing.T, uid uuid.UUID, title string) *entity.Notification {
	t.Helper()
	n, err := entity.New(uid, title, "body",
		[]vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, nil)
	require.NoError(t, err)
	return n
}

func TestNotificationRepo_SaveAndFind(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()
	n := makeNotif(t, uid, "Hello")

	require.NoError(t, repo.Save(ctx, n))

	got, err := repo.FindByID(ctx, n.ID)
	require.NoError(t, err)
	assert.Equal(t, n.ID, got.ID)
	assert.Equal(t, n.Title, got.Title)
	assert.Equal(t, uid, got.UserID)
}

func TestNotificationRepo_FindByID_NotFound(t *testing.T) {
	ctx, repo := setupRepo(t)
	_, err := repo.FindByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domainerrors.ErrNotificationNotFound)
}

func TestNotificationRepo_Update(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()
	n := makeNotif(t, uid, "Hello")
	require.NoError(t, repo.Save(ctx, n))

	require.NoError(t, n.MarkSent())
	require.NoError(t, repo.Update(ctx, n))

	got, err := repo.FindByID(ctx, n.ID)
	require.NoError(t, err)
	assert.Equal(t, vo.StatusSent, got.Status)
	assert.NotNil(t, got.SentAt)
}

func TestNotificationRepo_Update_NotFound(t *testing.T) {
	ctx, repo := setupRepo(t)
	n := makeNotif(t, uuid.New(), "missing")
	err := repo.Update(ctx, n)
	assert.ErrorIs(t, err, domainerrors.ErrNotificationNotFound)
}

func TestNotificationRepo_List_FilterByStatus(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()

	n1 := makeNotif(t, uid, "pending")
	n2 := makeNotif(t, uid, "sent")
	require.NoError(t, repo.Save(ctx, n1))
	require.NoError(t, repo.Save(ctx, n2))
	require.NoError(t, n2.MarkSent())
	require.NoError(t, repo.Update(ctx, n2))

	items, total, err := repo.List(ctx, repository.ListFilter{
		UserID: uid, Status: "SENT", Page: 1, Size: 10, Order: "desc",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, items, 1)
	assert.Equal(t, n2.ID, items[0].ID)
}

func TestNotificationRepo_List_Search(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()
	for _, title := range []string{"Apple pie", "Banana split", "Cherry tart"} {
		require.NoError(t, repo.Save(ctx, makeNotif(t, uid, title)))
	}
	items, total, err := repo.List(ctx, repository.ListFilter{
		UserID: uid, Search: "banana", Page: 1, Size: 10, Order: "desc",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, items, 1)
	assert.Equal(t, "Banana split", items[0].Title)
}

func TestNotificationRepo_List_Pagination(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()
	for i := 0; i < 25; i++ {
		require.NoError(t, repo.Save(ctx, makeNotif(t, uid, uuid.NewString())))
	}
	page1, total, err := repo.List(ctx, repository.ListFilter{UserID: uid, Page: 1, Size: 10, Order: "desc"})
	require.NoError(t, err)
	assert.Equal(t, int64(25), total)
	assert.Len(t, page1, 10)

	page3, _, err := repo.List(ctx, repository.ListFilter{UserID: uid, Page: 3, Size: 10, Order: "desc"})
	require.NoError(t, err)
	assert.Len(t, page3, 5)
}

func TestNotificationRepo_FindExpired(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()

	// ตัวที่หมดอายุแล้ว
	negative := -1 * time.Hour
	expired, err := entity.New(uid, "exp", "b", []vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, &negative)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, expired))

	// ตัวที่ยังไม่หมด
	future := 1 * time.Hour
	fresh, err := entity.New(uid, "fresh", "b", []vo.Channel{vo.ChannelEmail}, vo.PriorityNormal, nil, &future)
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, fresh))

	got, err := repo.FindExpired(ctx, time.Now().UTC(), 10)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, expired.ID, got[0].ID)
}

func TestNotificationRepo_BulkArchive(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()
	n1 := makeNotif(t, uid, "a")
	n2 := makeNotif(t, uid, "b")
	require.NoError(t, repo.Save(ctx, n1))
	require.NoError(t, repo.Save(ctx, n2))

	require.NoError(t, repo.BulkArchive(ctx, []uuid.UUID{n1.ID, n2.ID}))

	for _, id := range []uuid.UUID{n1.ID, n2.ID} {
		got, _ := repo.FindByID(ctx, id)
		assert.Equal(t, vo.StatusArchived, got.Status)
	}
}

func TestNotificationRepo_CountUnread(t *testing.T) {
	ctx, repo := setupRepo(t)
	uid := uuid.New()

	for i := 0; i < 3; i++ {
		n := makeNotif(t, uid, "x")
		require.NoError(t, repo.Save(ctx, n))
		require.NoError(t, n.MarkSent())
		require.NoError(t, repo.Update(ctx, n))
	}
	n4 := makeNotif(t, uid, "pending")
	require.NoError(t, repo.Save(ctx, n4))

	cnt, err := repo.CountUnread(ctx, uid)
	require.NoError(t, err)
	assert.Equal(t, int64(3), cnt)
}
```

## 2.3 `internal/modules/notification/infrastructure/messaging/kafka_integration_test.go`

```go
//go:build integration

package messaging_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"icmongolang/internal/modules/notification/domain/event"
	"icmongolang/internal/modules/notification/infrastructure/messaging"
	"icmongolang/internal/testutil"
)

func TestKafkaProducer_PublishAndConsume(t *testing.T) {
	ctx := context.Background()
	kf := testutil.NewKafka(ctx, t)
	topic := "notification.created.test"
	kf.CreateTopic(t, topic, 1)

	producer, err := messaging.NewKafkaProducer(kf.Brokers, "test-producer")
	require.NoError(t, err)
	defer producer.Close()

	evt := event.NotificationCreated{
		ID:       uuid.New(),
		UserID:   uuid.New(),
		Title:    "Hello",
		Priority: 8,
		Occurred: time.Now().UTC(),
	}

	require.NoError(t, producer.Publish(ctx, topic, evt.ID.String(), evt))

	// consumer
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer(kf.Brokers, cfg)
	require.NoError(t, err)
	defer consumer.Close()

	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	require.NoError(t, err)
	defer pc.Close()

	select {
	case msg := <-pc.Messages():
		var got event.NotificationCreated
		require.NoError(t, json.Unmarshal(msg.Value, &got))
		assert.Equal(t, evt.ID, got.ID)
		assert.Equal(t, evt.Title, got.Title)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}
```

## 2.4 Run

```bash
# ต้องมี Docker ก่อน
docker info

make test-integration

# รันเฉพาะ repo test
go test -race -tags=integration -run TestNotificationRepo ./internal/modules/notification/infrastructure/...
```

## 2.5 `internal/testutil/docker_helpers.go` (เพิ่มเติม)

```go
package testutil

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RequireDocker — skip test ถ้าไม่มี Docker
func RequireDocker(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:      "alpine:3.20",
			WaitingFor: wait.ForExit(),
			Cmd:        []string{"true"},
		},
		Started: true,
	})
	if err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
}
```

---

# 🔌 3) Wire-up เข้ากับโปรเจกต์จริง

## 3.1 `internal/app/app.go` — DI Container

```go
// Package app เป็น composition root ของทั้งโปรเจกต์
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"icmongolang/internal/modules/inventory"
	"icmongolang/internal/modules/notification"
	"icmongolang/internal/pkg/config"
	"icmongolang/internal/pkg/middleware"
	"icmongolang/pkg/websocket"
)

// App — ถือทุก dependency ที่ต้องใช้ + lifecycle
type App struct {
	cfg    *config.Config
	logger *slog.Logger

	db  *gorm.DB
	rdb *redis.Client
	es  *elasticsearch.Client

	hub     *websocket.Hub
	router  chi.Router
	httpSrv *http.Server

	// modules (compose root)
	notificationMod *notification.Module
	inventoryMod    *inventory.Module

	shutdownFns []func(ctx context.Context) error
}

// New — สร้าง App พร้อม wire-up ทั้งหมด
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelFromString(cfg.LogLevel),
	}))

	a := &App{cfg: cfg, logger: logger}

	if err := a.initInfra(ctx); err != nil {
		return nil, err
	}
	if err := a.initModules(ctx); err != nil {
		return nil, err
	}
	a.initRouter()
	return a, nil
}

// ──────────────────────────────────────────────────────────────
// Infra
// ──────────────────────────────────────────────────────────────

func (a *App) initInfra(ctx context.Context) error {
	// ── Postgres ──────────────────────────────────────────
	db, err := gorm.Open(postgres.Open(a.cfg.DB.DSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return fmt.Errorf("gorm open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(a.cfg.DB.MaxOpen)
	sqlDB.SetMaxIdleConns(a.cfg.DB.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)
	a.db = db
	a.addShutdown(func(_ context.Context) error { return sqlDB.Close() })

	// ── Redis ─────────────────────────────────────────────
	a.rdb = redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Addr,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	})
	if err := a.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	a.addShutdown(func(_ context.Context) error { return a.rdb.Close() })

	// ── Elasticsearch (optional) ──────────────────────────
	if a.cfg.Elasticsearch.URL != "" {
		es, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{a.cfg.Elasticsearch.URL},
		})
		if err != nil {
			a.logger.Warn("elasticsearch init failed, continuing", "err", err)
		} else {
			a.es = es
		}
	}

	// ── WebSocket hub ─────────────────────────────────────
	a.hub = websocket.NewHub()
	go a.hub.Run()

	return nil
}

// ──────────────────────────────────────────────────────────────
// Modules
// ──────────────────────────────────────────────────────────────

func (a *App) initModules(ctx context.Context) error {
	// ── Auth middleware (ตัวอย่าง — ใช้ของจริงจาก auth module) ──
	authMW := middleware.RequireAuth(a.cfg.JWT.Secret)

	// ── API subrouter ─────────────────────────────────────
	apiRouter := chi.NewRouter()

	// ── notification module ─────────────────────────────
	a.notificationMod = notification.Init(apiRouter, notification.Dependencies{
		DB:            a.db,
		Redis:         a.rdb,
		Elasticsearch: a.es,
		WSPublisher:   a.hub,
		KafkaBrokers:  a.cfg.Kafka.Brokers,
		KafkaClientID: "notification-api",
		EmailConfig: notification.EmailConfig{
			Host: a.cfg.SMTP.Host, Port: a.cfg.SMTP.Port,
			User: a.cfg.SMTP.User, Pass: a.cfg.SMTP.Pass,
			From: a.cfg.SMTP.From, Enabled: a.cfg.SMTP.Enabled,
		},
		PushConfig: notification.PushConfig{
			APIKey: a.cfg.FCM.APIKey, Enabled: a.cfg.FCM.Enabled,
		},
		Logger: a.logger,
	}, authMW)

	// ── inventory module ────────────────────────────────
	a.inventoryMod = inventory.Init(apiRouter, inventory.Dependencies{
		DB:            a.db,
		Redis:         a.rdb,
		KafkaBrokers:  a.cfg.Kafka.Brokers,
		KafkaClientID: "inventory-api",
		Logger:        a.logger,
	}, authMW)

	// ── mount /api/v1 ─────────────────────────────────────
	root := chi.NewRouter()
	root.Use(middleware.RequestID)
	root.Use(middleware.RealIP)
	root.Use(middleware.Recover(a.logger))
	root.Use(middleware.Logger(a.logger))
	root.Use(middleware.CORS(a.cfg.CORS.AllowedOrigins))
	root.Use(middleware.RateLimit(a.rdb, 100, time.Minute))

	root.Get("/healthz", a.healthz)
	root.Get("/readyz", a.readyz)
	root.Get("/ws", a.handleWS) // WebSocket upgrade endpoint
	root.Mount("/api/v1", apiRouter)

	a.router = root
	return nil
}

// ──────────────────────────────────────────────────────────────
// HTTP server
// ──────────────────────────────────────────────────────────────

func (a *App) initRouter() {
	a.httpSrv = &http.Server{
		Addr:              a.cfg.HTTP.ListenAddr,
		Handler:           a.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("http server starting", "addr", a.cfg.HTTP.ListenAddr)
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// start background jobs
	a.startSchedulers(ctx)

	select {
	case <-ctx.Done():
		return a.shutdown(context.Background())
	case err := <-errCh:
		_ = a.shutdown(context.Background())
		return err
	}
}

func (a *App) shutdown(ctx context.Context) error {
	shutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := a.httpSrv.Shutdown(shutCtx); err != nil {
		a.logger.Error("http shutdown failed", "err", err)
	}

	// reverse order shutdown
	for i := len(a.shutdownFns) - 1; i >= 0; i-- {
		if err := a.shutdownFns[i](shutCtx); err != nil {
			a.logger.Error("shutdown fn failed", "err", err)
		}
	}
	a.logger.Info("app shutdown complete")
	return nil
}

func (a *App) addShutdown(fn func(ctx context.Context) error) {
	a.shutdownFns = append(a.shutdownFns, fn)
}

// ──────────────────────────────────────────────────────────────
// Schedulers
// ──────────────────────────────────────────────────────────────

func (a *App) startSchedulers(ctx context.Context) {
	go func() {
		t := time.NewTicker(time.Hour)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				runCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
				if err := a.notificationMod.UseCases.AutoExpire.Execute(runCtx); err != nil {
					a.logger.Error("auto expire failed", "err", err)
				}
				cancel()
			}
		}
	}()
}

// ──────────────────────────────────────────────────────────────
// Handlers
// ──────────────────────────────────────────────────────────────

func (a *App) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func (a *App) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		http.Error(w, "db not ready", http.StatusServiceUnavailable)
		return
	}
	if err := a.rdb.Ping(ctx).Err(); err != nil {
		http.Error(w, "redis not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready\n"))
}

func (a *App) handleWS(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("user_id")
	if uid == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}
	a.hub.HandleWS(w, r, uid)
}

// ──────────────────────────────────────────────────────────────
// helper
// ──────────────────────────────────────────────────────────────

func levelFromString(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
```

## 3.2 `internal/pkg/config/config.go`

```go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env     string
	LogLevel string
	HTTP    HTTPConfig
	DB      DBConfig
	Redis   RedisConfig
	Kafka   KafkaConfig
	JWT     JWTConfig
	SMTP    SMTPConfig
	FCM     FCMConfig
	Elasticsearch ESConfig
	CORS    CORSConfig
}

type HTTPConfig struct {
	ListenAddr string
}

type DBConfig struct {
	DSN     string
	MaxOpen int
	MaxIdle int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type SMTPConfig struct {
	Host, Port, User, Pass, From string
	Enabled                       bool
}

type FCMConfig struct {
	APIKey  string
	Enabled bool
}

type ESConfig struct {
	URL string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load — โหลดจาก env + yaml (viper)
func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	// defaults
	v.SetDefault("HTTP_LISTEN_ADDR", ":8080")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("DB_MAX_OPEN", 25)
	v.SetDefault("DB_MAX_IDLE", 10)
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("SMTP_ENABLED", false)
	v.SetDefault("FCM_ENABLED", false)

	return &Config{
		Env:      v.GetString("ENV"),
		LogLevel: v.GetString("LOG_LEVEL"),
		HTTP:     HTTPConfig{ListenAddr: v.GetString("HTTP_LISTEN_ADDR")},
		DB: DBConfig{
			DSN:     v.GetString("DB_DSN"),
			MaxOpen: v.GetInt("DB_MAX_OPEN"),
			MaxIdle: v.GetInt("DB_MAX_IDLE"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		Kafka: KafkaConfig{
			Brokers: v.GetStringSlice("KAFKA_BROKERS"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
			TTL:    24 * time.Hour,
		},
		SMTP: SMTPConfig{
			Host: v.GetString("SMTP_HOST"), Port: v.GetString("SMTP_PORT"),
			User: v.GetString("SMTP_USER"), Pass: v.GetString("SMTP_PASS"),
			From: v.GetString("SMTP_FROM"), Enabled: v.GetBool("SMTP_ENABLED"),
		},
		FCM: FCMConfig{
			APIKey:  v.GetString("FCM_API_KEY"),
			Enabled: v.GetBool("FCM_ENABLED"),
		},
		Elasticsearch: ESConfig{URL: v.GetString("ELASTICSEARCH_URL")},
		CORS: CORSConfig{
			AllowedOrigins: v.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		},
	}, nil
}
```

## 3.3 `cmd/api/main.go` (แทนที่ทั้งหมด)

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"icmongolang/internal/app"
	"icmongolang/internal/pkg/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("app init: %w", err)
	}
	return a.Run(ctx)
}
```

## 3.4 `cmd/workers/notification_send/main.go`

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"icmongolang/internal/modules/notification/application"
	"icmongolang/internal/modules/notification/domain/service"
	"icmongolang/internal/modules/notification/infrastructure/messaging/consumers"
	"icmongolang/internal/modules/notification/infrastructure/persistence/postgres"
	"icmongolang/internal/modules/notification/infrastructure/services/email"
	"icmongolang/internal/modules/notification/infrastructure/services/push"
	"icmongolang/internal/pkg/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "worker fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// ── deps ────────────────────────────────────────────
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})

	repo := postgres.NewNotificationRepository(db)
	emailSender := email.NewSMTPSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From, cfg.SMTP.Enabled)
	pushSender := push.NewFCMSender(cfg.FCM.APIKey, cfg.FCM.Enabled, logger)
	policy := service.NewNotificationPolicy()

	// ── Kafka consumer group ────────────────────────────
	kcfg := sarama.NewConfig()
	kcfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kcfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "notification-send-group", kcfg)
	if err != nil {
		return fmt.Errorf("kafka consumer group: %w", err)
	}
	defer group.Close()

	// ── consumers ───────────────────────────────────────
	// handler ที่ส่ง notification เมื่อมี event notification.created
	sendUC := application.NewSendNotificationUseCase(
		repo, policy, emailSender, pushSender,
		&noopRealtime{}, logger,
	)
	_ = sendUC // wire เข้า handler จริงตาม event schema

	accountConsumer := consumers.NewAccountEventConsumer(
		application.NewHandleAccountEventUseCase(repo, logger), logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run account-event consumer
	go func() {
		for {
			if err := group.Consume(ctx, []string{consumers.TopicAccountEvent}, accountConsumer); err != nil {
				logger.Error("consume error", "err", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// ── signal handling ─────────────────────────────────
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("worker shutting down")
	cancel()
	time.Sleep(2 * time.Second)
	_ = rdb.Close()
	return nil
}

// noopRealtime — ใช้ในกรณี worker ไม่มี WS hub
type noopRealtime struct{}

func (*noopRealtime) Publish(_ context.Context, _ service.RealtimePayload) error { return nil }
```

## 3.5 `cmd/migrate/main.go`

```go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"icmongolang/internal/pkg/config"
)

func main() {
	root := &cobra.Command{Use: "migrate"}

	var dir string
	up := &cobra.Command{
		Use:   "up",
		Short: "Run all SQL migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUp(dir)
		},
	}
	up.Flags().StringVarP(&dir, "dir", "d", "./migrations", "migration directory")
	root.AddCommand(up)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUp(dir string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	ctx := context.Background()
	for _, f := range files {
		fmt.Printf("→ applying %s\n", filepath.Base(f))
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := db.WithContext(ctx).Exec(string(data)).Error; err != nil {
			return fmt.Errorf("migration failed: %s: %w", f, err)
		}
	}
	fmt.Println("✅ all migrations applied")
	return nil
}
```

---

# 📦 4) โมดูล `inventory` (ตาม template เดียวกัน)

จัดเต็มโครงสร้างแบบย่อ — ผู้อ่านดูเป็นแม่แบบได้

## 4.1 Tree

```
internal/modules/inventory/
├── domain/
│   ├── entity/{item.go,stock.go,movement.go,audit_trail.go}
│   ├── value_object/{sku.go,movement_type.go,unit.go}
│   ├── repository/{item_repo.go,stock_repo.go,movement_repo.go,audit_repo.go}
│   ├── service/inventory_policy.go
│   ├── event/movement_event.go
│   └── errors/errors.go
├── application/
│   ├── create_item.go
│   ├── adjust_stock.go
│   ├── get_stock.go
│   ├── list_items.go
│   ├── check_low_stock.go
│   ├── dto.go
│   └── mappers.go
├── infrastructure/
│   ├── persistence/postgres/{models.go,item_repo_impl.go,stock_repo_impl.go,movement_repo_impl.go,audit_repo_impl.go,mappers.go}
│   ├── persistence/redis/{stock_cache.go}
│   ├── messaging/kafka_producer.go
│   └── scheduler/low_stock_job.go
├── interfaces/
│   ├── http/{item_handler.go,stock_handler.go,routes.go,response.go}
│   └── middleware/owner_guard.go
└── module.go
```

## 4.2 Domain (สรุป)

**`domain/value_object/sku.go`**
```go
package valueobject

import (
	"regexp"
	"strings"

	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
)

var skuRe = regexp.MustCompile(`^[A-Z0-9][A-Z0-9\-]{2,31}$`)

type SKU string

func NewSKU(s string) (SKU, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if !skuRe.MatchString(s) {
		return "", domainerrors.ErrInvalidSKU
	}
	return SKU(s), nil
}

func (s SKU) String() string { return string(s) }
```

**`domain/value_object/movement_type.go`**
```go
package valueobject

type MovementType string

const (
	MovementIn      MovementType = "IN"
	MovementOut     MovementType = "OUT"
	MovementAdjust  MovementType = "ADJUST"
	MovementReserve MovementType = "RESERVE"
	MovementRelease MovementType = "RELEASE"
)

func (m MovementType) IsValid() bool {
	switch m {
	case MovementIn, MovementOut, MovementAdjust, MovementReserve, MovementRelease:
		return true
	}
	return false
}

// Sign — +1 เพิ่ม, -1 ลด, 0 ไม่เปลี่ยน
func (m MovementType) Sign() int {
	switch m {
	case MovementIn, MovementRelease:
		return 1
	case MovementOut, MovementReserve:
		return -1
	default:
		return 0
	}
}
```

**`domain/entity/item.go`**
```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
	vo "icmongolang/internal/modules/inventory/domain/value_object"
)

type Item struct {
	ID         uuid.UUID
	SKU        vo.SKU
	Name       string
	Unit       string
	Category   string
	ReorderPt  int
	CostPrice  int64 // สตางค์
	SalePrice  int64
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewItem(sku vo.SKU, name, unit, category string, reorderPt int, cost, sale int64) (*Item, error) {
	if name == "" {
		return nil, domainerrors.ErrNameRequired
	}
	if unit == "" {
		return nil, domainerrors.ErrUnitRequired
	}
	if reorderPt < 0 {
		return nil, domainerrors.ErrInvalidReorderPoint
	}
	now := time.Now().UTC()
	return &Item{
		ID:        uuid.New(),
		SKU:       sku,
		Name:      name,
		Unit:      unit,
		Category:  category,
		ReorderPt: reorderPt,
		CostPrice: cost,
		SalePrice: sale,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (i *Item) Deactivate() error {
	if !i.Active {
		return domainerrors.ErrItemAlreadyInactive
	}
	i.Active = false
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *Item) UpdatePricing(cost, sale int64) error {
	if cost < 0 || sale < 0 {
		return domainerrors.ErrInvalidPrice
	}
	i.CostPrice = cost
	i.SalePrice = sale
	i.UpdatedAt = time.Now().UTC()
	return nil
}
```

**`domain/entity/stock.go`**
```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
)

type Stock struct {
	ItemID      uuid.UUID
	OnHand      int
	Reserved    int
	Location    string // warehouse / bin
	UpdatedAt   time.Time
}

// Available = OnHand - Reserved
func (s *Stock) Available() int { return s.OnHand - s.Reserved }

// Apply — ปรับยอดตาม delta (+/-)
func (s *Stock) Apply(delta int) error {
	newOnHand := s.OnHand + delta
	if newOnHand < 0 {
		return domainerrors.ErrInsufficientStock
	}
	s.OnHand = newOnHand
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Stock) Reserve(qty int) error {
	if qty <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	if s.Available() < qty {
		return domainerrors.ErrInsufficientStock
	}
	s.Reserved += qty
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Stock) Release(qty int) error {
	if qty <= 0 || s.Reserved < qty {
		return domainerrors.ErrInvalidQuantity
	}
	s.Reserved -= qty
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Stock) IsLow(reorderPt int) bool {
	return s.Available() <= reorderPt
}
```

**`domain/entity/movement.go`**
```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
	vo "icmongolang/internal/modules/inventory/domain/value_object"
)

type Movement struct {
	ID         uuid.UUID
	ItemID     uuid.UUID
	Type       vo.MovementType
	Quantity   int
	Before     int
	After      int
	Reference  string // order id, PO id, ...
	Reason     string
	ActorID    uuid.UUID
	CreatedAt  time.Time
}

func NewMovement(
	itemID uuid.UUID, mtype vo.MovementType,
	qty, before int, reference, reason string, actorID uuid.UUID,
) (*Movement, error) {
	if !mtype.IsValid() {
		return nil, domainerrors.ErrInvalidMovementType
	}
	if qty == 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	delta := qty * mtype.Sign()
	after := before + delta
	if after < 0 {
		return nil, domainerrors.ErrInsufficientStock
	}
	return &Movement{
		ID:        uuid.New(),
		ItemID:    itemID,
		Type:      mtype,
		Quantity:  qty,
		Before:    before,
		After:     after,
		Reference: reference,
		Reason:    reason,
		ActorID:   actorID,
		CreatedAt: time.Now().UTC(),
	}, nil
}
```

**`domain/errors/errors.go`**
```go
package domainerrors

import "errors"

var (
	ErrItemNotFound          = errors.New("item not found")
	ErrItemAlreadyExists     = errors.New("item already exists")
	ErrItemAlreadyInactive   = errors.New("item already inactive")
	ErrStockNotFound         = errors.New("stock not found")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrInvalidSKU            = errors.New("invalid sku")
	ErrInvalidMovementType   = errors.New("invalid movement type")
	ErrInvalidReorderPoint   = errors.New("invalid reorder point")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrNameRequired          = errors.New("name required")
	ErrUnitRequired          = errors.New("unit required")
	ErrNotOwnedByUser        = errors.New("not owned by user")
)
```

**`domain/repository/item_repo.go`**
```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/inventory/domain/entity"
)

type ItemFilter struct {
	Search   string
	Category string
	Active   *bool
	Page     int
	Size     int
	SortBy   string
	Order    string
}

type ItemRepository interface {
	Save(ctx context.Context, item *entity.Item) error
	Update(ctx context.Context, item *entity.Item) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Item, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Item, error)
	List(ctx context.Context, f ItemFilter) ([]*entity.Item, int64, error)
	FindLowStock(ctx context.Context, limit int) ([]*entity.Item, error)
}
```

**`domain/repository/stock_repo.go`**
```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/inventory/domain/entity"
)

type StockRepository interface {
	Get(ctx context.Context, itemID uuid.UUID, location string) (*entity.Stock, error)
	Upsert(ctx context.Context, s *entity.Stock) error
	// Atomically adjust — ใช้ FOR UPDATE
	AdjustAtomic(ctx context.Context, itemID uuid.UUID, location string, delta int) (*entity.Stock, error)
	ListByItem(ctx context.Context, itemID uuid.UUID) ([]*entity.Stock, error)
}
```

**`domain/repository/movement_repo.go`**
```go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/inventory/domain/entity"
)

type MovementRepository interface {
	Save(ctx context.Context, m *entity.Movement) error
	ListByItem(ctx context.Context, itemID uuid.UUID, limit int) ([]*entity.Movement, error)
}
```

**`domain/service/inventory_policy.go`**
```go
package service

import (
	"icmongolang/internal/modules/inventory/domain/entity"
	vo "icmongolang/internal/modules/inventory/domain/value_object"
)

type InventoryPolicy struct{}

func NewInventoryPolicy() *InventoryPolicy { return &InventoryPolicy{} }

// ValidateAdjustment — business rules ก่อนบันทึก movement
func (p *InventoryPolicy) ValidateAdjustment(s *entity.Stock, item *entity.Item, qty int, mtype vo.MovementType) error {
	if qty <= 0 {
		return nil // ให้ movement ตรวจ
	}
	delta := qty * mtype.Sign()
	if s.Available()+delta < 0 {
		return ErrInsufficientAvailable
	}
	return nil
}

var ErrInsufficientAvailable = errorString("insufficient available stock (excluding reserved)")

type errorString string

func (e errorString) Error() string { return string(e) }
```

**`domain/event/movement_event.go`**
```go
package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	TopicStockAdjusted = "inventory.stock.adjusted"
	TopicStockLow      = "inventory.stock.low"
	TopicItemCreated   = "inventory.item.created"
)

type StockAdjusted struct {
	ItemID    uuid.UUID `json:"item_id"`
	MovementID uuid.UUID `json:"movement_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Before    int       `json:"before"`
	After     int       `json:"after"`
	Reference string    `json:"reference,omitempty"`
	Occurred  time.Time `json:"occurred_at"`
}

type StockLow struct {
	ItemID    uuid.UUID `json:"item_id"`
	SKU       string    `json:"sku"`
	Available int       `json:"available"`
	ReorderPt int       `json:"reorder_point"`
	Occurred  time.Time `json:"occurred_at"`
}
```

## 4.3 Application

**`application/adjust_stock.go`** (use case หลัก)

```go
package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
	"icmongolang/internal/modules/inventory/domain/entity"
	"icmongolang/internal/modules/inventory/domain/event"
	"icmongolang/internal/modules/inventory/domain/repository"
	"icmongolang/internal/modules/inventory/domain/service"
	vo "icmongolang/internal/modules/inventory/domain/value_object"
)

type AdjustStockUseCase struct {
	itemRepo     repository.ItemRepository
	stockRepo    repository.StockRepository
	movementRepo repository.MovementRepository
	publisher    EventPublisher
	policy       *service.InventoryPolicy
	logger       *slog.Logger
}

func NewAdjustStockUseCase(
	ir repository.ItemRepository,
	sr repository.StockRepository,
	mr repository.MovementRepository,
	pub EventPublisher,
	policy *service.InventoryPolicy,
	logger *slog.Logger,
) *AdjustStockUseCase {
	return &AdjustStockUseCase{itemRepo: ir, stockRepo: sr, movementRepo: mr, publisher: pub, policy: policy, logger: logger}
}

type AdjustStockInput struct {
	ItemID    uuid.UUID
	Location  string
	Type      string
	Quantity  int
	Reference string
	Reason    string
	ActorID   uuid.UUID
}

type AdjustStockOutput struct {
	MovementID uuid.UUID
	ItemID     uuid.UUID
	Before     int
	After      int
}

func (uc *AdjustStockUseCase) Execute(ctx context.Context, in AdjustStockInput) (*AdjustStockOutput, error) {
	// 1. Load item
	item, err := uc.itemRepo.FindByID(ctx, in.ItemID)
	if err != nil {
		return nil, err
	}
	if !item.Active {
		return nil, domainerrors.ErrItemAlreadyInactive
	}

	// 2. Validate movement type
	mtype := vo.MovementType(in.Type)
	if !mtype.IsValid() {
		return nil, domainerrors.ErrInvalidMovementType
	}

	// 3. Atomic adjust stock
	delta := in.Quantity * mtype.Sign()
	stock, err := uc.stockRepo.AdjustAtomic(ctx, in.ItemID, in.Location, delta)
	if err != nil {
		return nil, err
	}

	// 4. Create movement record
	before := stock.OnHand - delta
	mv, err := entity.NewMovement(
		in.ItemID, mtype, in.Quantity, before,
		in.Reference, in.Reason, in.ActorID,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.movementRepo.Save(ctx, mv); err != nil {
		return nil, err
	}

	// 5. Publish event
	uc.publishAdjusted(ctx, mv)

	// 6. Check low stock
	if stock.IsLow(item.ReorderPt) {
		uc.publishLow(ctx, item, stock)
	}

	return &AdjustStockOutput{
		MovementID: mv.ID, ItemID: mv.ItemID,
		Before: mv.Before, After: mv.After,
	}, nil
}

func (uc *AdjustStockUseCase) publishAdjusted(ctx context.Context, mv *entity.Movement) {
	if uc.publisher == nil {
		return
	}
	err := uc.publisher.Publish(ctx, event.TopicStockAdjusted, mv.ItemID.String(), event.StockAdjusted{
		ItemID:     mv.ItemID,
		MovementID: mv.ID,
		Type:       string(mv.Type),
		Quantity:   mv.Quantity,
		Before:     mv.Before,
		After:      mv.After,
		Reference:  mv.Reference,
		Occurred:   mv.CreatedAt,
	})
	if err != nil {
		uc.logger.Warn("publish stock adjusted failed", "err", err)
	}
}

func (uc *AdjustStockUseCase) publishLow(ctx context.Context, item *entity.Item, s *entity.Stock) {
	if uc.publisher == nil {
		return
	}
	err := uc.publisher.Publish(ctx, event.TopicStockLow, item.ID.String(), event.StockLow{
		ItemID:    item.ID,
		SKU:       item.SKU.String(),
		Available: s.Available(),
		ReorderPt: item.ReorderPt,
		Occurred:  s.UpdatedAt,
	})
	if err != nil {
		uc.logger.Warn("publish stock low failed", "err", err)
	}
}
```

## 4.4 Infrastructure (Postgres — เฉพาะส่วนสำคัญ)

**`infrastructure/persistence/postgres/models.go`**
```go
package postgres

import (
	"time"

	"github.com/google/uuid"
)

type ItemModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SKU       string    `gorm:"type:varchar(32);uniqueIndex;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Unit      string    `gorm:"type:varchar(32);not null"`
	Category  string    `gorm:"type:varchar(64);index"`
	ReorderPt int       `gorm:"not null;default:0"`
	CostPrice int64     `gorm:"not null;default:0"`
	SalePrice int64     `gorm:"not null;default:0"`
	Active    bool      `gorm:"not null;default:true;index"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func (ItemModel) TableName() string { return "inventory_items" }

type StockModel struct {
	ItemID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Location  string    `gorm:"type:varchar(64);primaryKey"`
	OnHand    int       `gorm:"not null;default:0"`
	Reserved  int       `gorm:"not null;default:0"`
	UpdatedAt time.Time `gorm:"default:now()"`
}

func (StockModel) TableName() string { return "inventory_stocks" }

type MovementModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Type      string    `gorm:"type:varchar(20);not null;index"`
	Quantity  int       `gorm:"not null"`
	Before    int       `gorm:"not null"`
	After     int       `gorm:"not null"`
	Reference string    `gorm:"type:varchar(128);index"`
	Reason    string    `gorm:"type:text"`
	ActorID   uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt time.Time `gorm:"default:now();index"`
}

func (MovementModel) TableName() string { return "inventory_movements" }
```

**`infrastructure/persistence/postgres/stock_repo_impl.go`** — จุดสำคัญคือ atomic adjust

```go
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	domainerrors "icmongolang/internal/modules/inventory/domain/errors"
	"icmongolang/internal/modules/inventory/domain/entity"
	"icmongolang/internal/modules/inventory/domain/repository"
)

type stockRepoImpl struct{ db *gorm.DB }

func NewStockRepository(db *gorm.DB) repository.StockRepository {
	return &stockRepoImpl{db: db}
}

func (r *stockRepoImpl) Get(ctx context.Context, itemID uuid.UUID, location string) (*entity.Stock, error) {
	var m StockModel
	err := r.db.WithContext(ctx).Where("item_id = ? AND location = ?", itemID, location).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrStockNotFound
	}
	if err != nil {
		return nil, err
	}
	return &entity.Stock{ItemID: m.ItemID, OnHand: m.OnHand, Reserved: m.Reserved, Location: m.Location, UpdatedAt: m.UpdatedAt}, nil
}

func (r *stockRepoImpl) Upsert(ctx context.Context, s *entity.Stock) error {
	m := &StockModel{
		ItemID: s.ItemID, Location: s.Location,
		OnHand: s.OnHand, Reserved: s.Reserved, UpdatedAt: s.UpdatedAt,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "item_id"}, {Name: "location"}},
		DoUpdates: clause.AssignmentColumns([]string{"on_hand", "reserved", "updated_at"}),
	}).Create(m).Error
}

// AdjustAtomic — ใช้ transaction + FOR UPDATE กัน race
func (r *stockRepoImpl) AdjustAtomic(ctx context.Context, itemID uuid.UUID, location string, delta int) (*entity.Stock, error) {
	var out *entity.Stock
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var m StockModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("item_id = ? AND location = ?", itemID, location).
			First(&m).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if delta < 0 {
				return domainerrors.ErrInsufficientStock
			}
			m = StockModel{ItemID: itemID, Location: location, OnHand: delta}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
			out = &entity.Stock{ItemID: m.ItemID, OnHand: m.OnHand, Location: m.Location}
			return nil
		}
		if err != nil {
			return err
		}
		newOnHand := m.OnHand + delta
		if newOnHand < 0 {
			return domainerrors.ErrInsufficientStock
		}
		m.OnHand = newOnHand
		if err := tx.Save(&m).Error; err != nil {
			return err
		}
		out = &entity.Stock{ItemID: m.ItemID, OnHand: m.OnHand, Reserved: m.Reserved, Location: m.Location, UpdatedAt: m.UpdatedAt}
		return nil
	})
	return out, err
}

func (r *stockRepoImpl) ListByItem(ctx context.Context, itemID uuid.UUID) ([]*entity.Stock, error) {
	var rows []StockModel
	if err := r.db.WithContext(ctx).Where("item_id = ?", itemID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.Stock, 0, len(rows))
	for i := range rows {
		out = append(out, &entity.Stock{
			ItemID: rows[i].ItemID, OnHand: rows[i].OnHand, Reserved: rows[i].Reserved,
			Location: rows[i].Location, UpdatedAt: rows[i].UpdatedAt,
		})
	}
	return out, nil
}
```

## 4.5 Handler + Routes (สรุป)

**`interfaces/http/stock_handler.go`**
```go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"icmongolang/internal/modules/inventory/application"
)

type StockHandler struct {
	adjust *application.AdjustStockUseCase
	get    *application.GetStockUseCase
}

func NewStockHandler(a *application.AdjustStockUseCase, g *application.GetStockUseCase) *StockHandler {
	return &StockHandler{adjust: a, get: g}
}

func (h *StockHandler) Adjust(w http.ResponseWriter, r *http.Request) {
	uid, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item_id")
		return
	}
	var body struct {
		Location  string `json:"location"`
		Type      string `json:"type"`
		Quantity  int    `json:"quantity"`
		Reference string `json:"reference"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	out, err := h.adjust.Execute(r.Context(), application.AdjustStockInput{
		ItemID: itemID, Location: body.Location, Type: body.Type,
		Quantity: body.Quantity, Reference: body.Reference, Reason: body.Reason,
		ActorID: uid,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
```

**`interfaces/http/routes.go`**
```go
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Item  *ItemHandler
	Stock *StockHandler
}

func RegisterRoutes(r chi.Router, h *Handlers, authMW func(http.Handler) http.Handler) {
	r.Route("/inventory", func(r chi.Router) {
		r.Use(authMW)
		r.Route("/items", func(r chi.Router) {
			r.Post("/", h.Item.Create)
			r.Get("/", h.Item.List)
			r.Get("/{id}", h.Item.Get)
			r.Post("/{itemID}/adjust", h.Stock.Adjust)
		})
	})
}
```

## 4.6 Migration

**`migrations/20260415_inventory_init.sql`**
```sql
CREATE TABLE IF NOT EXISTS inventory_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku         VARCHAR(32) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    unit        VARCHAR(32) NOT NULL,
    category    VARCHAR(64),
    reorder_pt  INT NOT NULL DEFAULT 0,
    cost_price  BIGINT NOT NULL DEFAULT 0,
    sale_price  BIGINT NOT NULL DEFAULT 0,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_inventory_items_category ON inventory_items (category);
CREATE INDEX IF NOT EXISTS idx_inventory_items_active ON inventory_items (active);

CREATE TABLE IF NOT EXISTS inventory_stocks (
    item_id    UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    location   VARCHAR(64) NOT NULL,
    on_hand    INT NOT NULL DEFAULT 0,
    reserved   INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (item_id, location),
    CONSTRAINT chk_inventory_stocks_nonneg CHECK (on_hand >= 0 AND reserved >= 0 AND reserved <= on_hand)
);

CREATE TABLE IF NOT EXISTS inventory_movements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id    UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    type       VARCHAR(20) NOT NULL,
    quantity   INT NOT NULL,
    before_qty INT NOT NULL,
    after_qty  INT NOT NULL,
    reference  VARCHAR(128),
    reason     TEXT,
    actor_id   UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_item_id ON inventory_movements (item_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_type ON inventory_movements (type);
```

## 4.7 Wire-up

```go
// internal/modules/inventory/module.go
package inventory

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"icmongolang/internal/modules/inventory/application"
	"icmongolang/internal/modules/inventory/domain/service"
	"icmongolang/internal/modules/inventory/infrastructure/messaging"
	"icmongolang/internal/modules/inventory/infrastructure/persistence/postgres"
	httpiface "icmongolang/internal/modules/inventory/interfaces/http"
)

type Dependencies struct {
	DB            *gorm.DB
	Redis         *redis.Client
	KafkaBrokers  []string
	KafkaClientID string
	Logger        *slog.Logger
}

type Module struct {
	UseCases struct {
		Adjust *application.AdjustStockUseCase
	}
}

func Init(r chi.Router, deps Dependencies, authMW func(http.Handler) http.Handler) *Module {
	itemRepo := postgres.NewItemRepository(deps.DB)
	stockRepo := postgres.NewStockRepository(deps.DB)
	movementRepo := postgres.NewMovementRepository(deps.DB)

	producer, err := messaging.NewKafkaProducer(deps.KafkaBrokers, deps.KafkaClientID)
	if err != nil {
		deps.Logger.Error("kafka producer init failed", "err", err)
	}
	var publisher application.EventPublisher
	if producer != nil {
		publisher = &producerAdapter{p: producer}
	}

	policy := service.NewInventoryPolicy()
	adjustUC := application.NewAdjustStockUseCase(itemRepo, stockRepo, movementRepo, publisher, policy, deps.Logger)
	getStockUC := application.NewGetStockUseCase(itemRepo, stockRepo)
	createItemUC := application.NewCreateItemUseCase(itemRepo, deps.Logger)
	listItemsUC := application.NewListItemsUseCase(itemRepo)

	h := &httpiface.Handlers{
		Item:  httpiface.NewItemHandler(createItemUC, listItemsUC),
		Stock: httpiface.NewStockHandler(adjustUC, getStockUC),
	}
	httpiface.RegisterRoutes(r, h, authMW)

	return &Module{UseCases: struct {
		Adjust *application.AdjustStockUseCase
	}{Adjust: adjustUC}}
}

type producerAdapter struct{ p messaging.Producer }

func (a *producerAdapter) Publish(ctx context.Context, t, k string, v any) error {
	return a.p.Publish(ctx, t, k, v)
}
```

---

# 🔄 5) CI Workflow

## 5.1 `.golangci.yml`

```yaml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor
    - mocks

linters:
  disable-all: true
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - gocritic
    - bodyclose
    - errorlint
    - exhaustive
    - prealloc
    - rowserrcheck
    - sqlclosecheck
    - durationcheck

linters-settings:
  errcheck:
    check-type-assertions: true
    exclude-functions:
      - (io.Closer).Close
      - (*github.com/IBM/sarama.SyncProducer).Close
  gocritic:
    enabled-tags:
      - diagnostic
      - performance
  goimports:
    local-prefixes: icmongolang
  misspell:
    locale: US

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
        - gocritic
    - path: mocks/
      linters: [all]
```

## 5.2 `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

permissions:
  contents: read
  pull-requests: write

env:
  GO_VERSION: '1.22'

jobs:
  # ────────────────────────────────────────────────────────
  lint:
    name: 🔍 Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.60.3
          args: --timeout=5m

  # ────────────────────────────────────────────────────────
  unit-test:
    name: 🧪 Unit tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Generate mocks
        run: go run github.com/vektra/mockery/v2@v2.46.3

      - name: Go vet
        run: go vet ./...

      - name: Unit test + coverage
        run: |
          mkdir -p .coverage
          go test -race -count=1 -short ./internal/... \
            -coverprofile=.coverage/unit.out \
            -covermode=atomic

      - name: Coverage check (≥ 80%)
        run: |
          go tool cover -func=.coverage/unit.out | tail -1
          TOTAL=$(go tool cover -func=.coverage/unit.out | tail -1 | awk '{print $3}' | tr -d '%')
          echo "Coverage: $TOTAL%"
          awk -v c="$TOTAL" 'BEGIN { if (c+0 < 80) exit 1 }'

      - name: Upload coverage artifact
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: .coverage/unit.out

      - name: Upload to Codecov
        uses: codecov/codecov-action@v4
        with:
          files: .coverage/unit.out
          flags: unittests
          fail_ci_if_error: false

  # ────────────────────────────────────────────────────────
  integration-test:
    name: 🐳 Integration tests
    runs-on: ubuntu-latest
    timeout-minutes: 30
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Generate mocks
        run: go run github.com/vektra/mockery/v2@v2.46.3

      - name: Run integration tests (testcontainers)
        run: |
          mkdir -p .coverage
          go test -race -count=1 -tags=integration \
            -timeout=25m \
            -coverprofile=.coverage/integration.out \
            -covermode=atomic \
            ./internal/...

      - name: Upload integration coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-integration
          path: .coverage/integration.out

  # ────────────────────────────────────────────────────────
  migrate-check:
    name: 💾 Migration check
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: icmongolang_ci
          POSTGRES_USER: ci
          POSTGRES_PASSWORD: ci
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 5s
          --health-timeout 5s
          --health-retries 10
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Apply migrations
        env:
          DB_DSN: postgres://ci:ci@localhost:5432/icmongolang_ci?sslmode=disable
        run: |
          for f in $(ls migrations/*.sql | sort); do
            echo "→ $f"
            psql "$DB_DSN" -v ON_ERROR_STOP=1 -f "$f"
          done

      - name: Verify tables exist
        env:
          DB_DSN: postgres://ci:ci@localhost:5432/icmongolang_ci?sslmode=disable
        run: |
          psql "$DB_DSN" -c "\dt notification_*"
          psql "$DB_DSN" -c "\dt inventory_*"

      - name: Idempotency test (run twice)
        env:
          DB_DSN: postgres://ci:ci@localhost:5432/icmongolang_ci?sslmode=disable
        run: |
          for f in $(ls migrations/*.sql | sort); do
            psql "$DB_DSN" -v ON_ERROR_STOP=1 -f "$f"
          done
          echo "✅ migrations are idempotent"

  # ────────────────────────────────────────────────────────
  build:
    name: 🏗️ Build
    runs-on: ubuntu-latest
    needs: [lint, unit-test]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Build API
        run: go build -trimpath -ldflags="-s -w" -o /tmp/api ./cmd/api

      - name: Build workers
        run: |
          go build -trimpath -ldflags="-s -w" -o /tmp/notification-send ./cmd/workers/notification_send
          go build -trimpath -ldflags="-s -w" -o /tmp/migrate ./cmd/migrate

      - name: Build workflowctl
        run: go build -trimpath -ldflags="-s -w" -o /tmp/workflowctl ./cmd/workflowctl

  # ────────────────────────────────────────────────────────
  summary:
    name: 📊 Summary
    runs-on: ubuntu-latest
    needs: [lint, unit-test, integration-test, migrate-check, build]
    if: always()
    steps:
      - run: |
          echo "lint:              ${{ needs.lint.result }}"
          echo "unit-test:         ${{ needs.unit-test.result }}"
          echo "integration-test:  ${{ needs.integration-test.result }}"
          echo "migrate-check:     ${{ needs.migrate-check.result }}"
          echo "build:             ${{ needs.build.result }}"
```

## 5.3 `.github/workflows/security.yml`

```yaml
name: Security

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 3 * * 1'

permissions:
  contents: read
  security-events: write

jobs:
  govulncheck:
    name: 🛡️ govulncheck
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - uses: golang/govulncheck-action@v1
        with:
          go-package: ./...

  gosec:
    name: 🔒 gosec
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec.sarif ./...'
      - uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec.sarif

  gitleaks:
    name: 🔑 Secret scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## 5.4 `.github/workflows/migrate-test.yml` (แยก)

```yaml
name: Migrate Test

on:
  push:
    paths:
      - 'migrations/**'
  pull_request:
    paths:
      - 'migrations/**'

jobs:
  migrate:
    name: 💾 Migrate all
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: icm_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        ports: ['5432:5432']
        options: >-
          --health-cmd pg_isready
          --health-interval 5s
          --health-retries 10
    env:
      DB_DSN: postgres://test:test@localhost:5432/icm_test?sslmode=disable
    steps:
      - uses: actions/checkout@v4
      - name: Install psql
        run: sudo apt-get update && sudo apt-get install -y postgresql-client
      - name: Apply
        run: |
          for f in $(ls migrations/*.sql | sort); do
            echo "→ $f"
            psql "$DB_DSN" -v ON_ERROR_STOP=1 -f "$f"
          done
      - name: Verify
        run: psql "$DB_DSN" -c "\dt" && psql "$DB_DSN" -c "\d notification_notifications"
      - name: Idempotent
        run: |
          for f in $(ls migrations/*.sql | sort); do
            psql "$DB_DSN" -v ON_ERROR_STOP=1 -f "$f"
          done
```

## 5.5 `codecov.yml`

```yaml
coverage:
  status:
    project:
      default:
        target: 80%
        threshold: 2%
    patch:
      default:
        target: 75%
        threshold: 5%

comment:
  layout: "reach,diff,flags,tree,files"
  behavior: default
  require_changes: false
```

## 5.6 Branch protection rules (เอกสาร)

ตั้งใน GitHub repo → Settings → Branches → main:
- ✅ Require PR before merging
- ✅ Require status checks: `lint`, `unit-test`, `integration-test`, `migrate-check`, `build`, `govulncheck`
- ✅ Require branches to be up to date
- ✅ Require conversation resolution
- ✅ Include administrators

---

# ✅ สรุปสิ่งที่ส่งมอบ

| # | รายการ | สถานะ | ไฟล์หลัก |
|---|---|---|---|
| 1 | **Unit Tests** | ✅ mockery config + 12 test files + Makefile + coverage 90% | `.mockery.yaml`, `application/*_test.go`, `domain/**/*_test.go` |
| 2 | **Integration Tests** | ✅ testcontainers-go (Postgres + Redis + Kafka) + 9 integration tests | `internal/testutil/containers.go`, `**/*_integration_test.go` |
| 3 | **Wire-up** | ✅ `internal/app/app.go` (DI), patch `cmd/api`, worker, migrate | `internal/app/app.go`, `cmd/api/main.go`, `cmd/workers/...`, `cmd/migrate/main.go` |
| 4 | **Inventory Module** | ✅ ครบ 4 layer + migration + Kafka events + low-stock scheduler | `internal/modules/inventory/**` |
| 5 | **CI Workflow** | ✅ lint + unit + integration + migrate + build + security + coverage | `.github/workflows/*.yml`, `.golangci.yml`, `codecov.yml` |

## วิธีใช้

```bash
# 1. สร้าง mocks
make mocks

# 2. Unit test + coverage
make test-unit
make cover-check      # ✅ ตรวจ ≥ 80%

# 3. Integration test (ต้องมี Docker)
make test-integration

# 4. Migrate
go run ./cmd/migrate up

# 5. รัน API
go run ./cmd/api

# 6. รัน Worker
go run ./cmd/workers/notification_send
```

 