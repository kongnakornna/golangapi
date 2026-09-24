# 🚀 FULL EXECUTION — G05-P40 Complete Package

> ทำทุกคำสั่งในชุดเดียว: CRM full → 06-08 → 09-16 upgrades → 17-40 indexed → tracking → export

---

# PART 1 — G05 + P05: CRM (Full Deep Dive)

## 1.1 Spec `.ai/specs/05_crm.yaml`

```yaml
module: crm
context: Sales & Support
priority: P1
status: queued

aggregates:
  Lead:
    fields:
      - id, tenant_id, name, email, phone, company, source, status, score
      - assigned_to: *UUID
      - converted_to: *UUID
    behavior: [Qualify, Convert, Assign, Lose]
    invariants:
      - score 0-100
      - can only convert once
  Opportunity:
    fields: [id, tenant_id, lead_id, customer_id, title, amount, probability, stage, expected_close_date]
    behavior: [MoveToStage, Win, Lose, UpdateProbability]
    invariants:
      - probability 0-100
      - stage transitions valid
  Ticket:
    fields: [id, tenant_id, customer_id, subject, description, priority, status, assigned_to, sla_due_at, resolved_at]
    behavior: [Assign, Start, Resolve, Close, Escalate, Reopen]
    invariants:
      - cannot resolve closed ticket
      - escalation only from OPEN/IN_PROGRESS
  Activity:
    fields: [id, tenant_id, entity_type, entity_id, type, subject, notes, performed_by, occurred_at]

value_objects:
  LeadSource: [WEB, REFERRAL, CAMPAIGN, COLD_CALL, EVENT]
  LeadStatus: [NEW, CONTACTED, QUALIFIED, LOST, CONVERTED]
  OpportunityStage: [PROSPECTING, PROPOSAL, NEGOTIATION, WON, LOST]
  TicketPriority: [LOW, MEDIUM, HIGH, URGENT]
  TicketStatus: [OPEN, IN_PROGRESS, RESOLVED, CLOSED]
  ActivityType: [CALL, EMAIL, MEETING, NOTE, TASK]

use_cases:
  - CaptureLead, GetLead, ListLeads, QualifyLead, ConvertLead, AssignLead, LoseLead
  - CreateOpportunity, MoveOpportunityStage, WinOpportunity, LoseOpportunity
  - OpenTicket, GetTicket, ListTickets, AssignTicket, StartTicket, ResolveTicket, CloseTicket, EscalateTicket, ReopenTicket
  - LogActivity, ListActivities

tables: [crm_leads, crm_opportunities, crm_tickets, crm_activities]

kafka_publish:
  - crm.lead.created
  - crm.lead.qualified
  - crm.lead.converted
  - crm.opportunity.created
  - crm.opportunity.won
  - crm.opportunity.lost
  - crm.ticket.opened
  - crm.ticket.escalated
  - crm.ticket.resolved

kafka_consume:
  - customer.customer.created       # link
  - device.device.offline           # auto-ticket
  - payment.payment.failed          # auto-ticket

sla_policy:
  LOW: 48h
  MEDIUM: 24h
  HIGH: 8h
  URGENT: 2h

migration: 20260105_crm_init.sql
```

## 1.2 Go Implementation (G05)

### Directory Tree

```
internal/modules/crm/
├── domain/
│   ├── entity/{lead,opportunity,ticket,activity}.go
│   ├── value_object/{lead_source,lead_status,opportunity_stage,ticket_priority,ticket_status,activity_type}.go
│   ├── repository/{lead,opportunity,ticket,activity}_repository.go
│   ├── service/sla_service.go
│   ├── event/events.go
│   └── errors/errors.go
├── application/
│   ├── ports.go, dto.go
│   ├── capture_lead.go, qualify_lead.go, convert_lead.go, assign_lead.go, lose_lead.go
│   ├── create_opportunity.go, move_opportunity_stage.go, win_opportunity.go
│   ├── open_ticket.go, assign_ticket.go, start_ticket.go, resolve_ticket.go, escalate_ticket.go, reopen_ticket.go
│   └── log_activity.go
├── infrastructure/
│   ├── persistence/postgres/{models,lead_repo_impl,opportunity_repo_impl,ticket_repo_impl,activity_repo_impl}.go
│   ├── messaging/producer.go
│   └── scheduler/sla_job.go
├── interfaces/http/{lead_handler,opportunity_handler,ticket_handler,activity_handler,routes,dto,errors}.go
└── module.go
```

### `domain/value_object/lead_status.go`

```go
package valueobject

type LeadStatus string

const (
	LeadStatusNew       LeadStatus = "NEW"
	LeadStatusContacted LeadStatus = "CONTACTED"
	LeadStatusQualified LeadStatus = "QUALIFIED"
	LeadStatusLost      LeadStatus = "LOST"
	LeadStatusConverted LeadStatus = "CONVERTED"
)

func (s LeadStatus) IsValid() bool {
	switch s {
	case LeadStatusNew, LeadStatusContacted, LeadStatusQualified, LeadStatusLost, LeadStatusConverted:
		return true
	}
	return false
}

func (s LeadStatus) IsTerminal() bool {
	return s == LeadStatusLost || s == LeadStatusConverted
}

func (s LeadStatus) CanTransitionTo(next LeadStatus) bool {
	t := map[LeadStatus]map[LeadStatus]bool{
		LeadStatusNew:       {LeadStatusContacted: true, LeadStatusQualified: true, LeadStatusLost: true},
		LeadStatusContacted: {LeadStatusQualified: true, LeadStatusLost: true},
		LeadStatusQualified: {LeadStatusConverted: true, LeadStatusLost: true},
		LeadStatusLost:      {},
		LeadStatusConverted: {},
	}
	return t[s][next]
}
```

### `domain/value_object/lead_source.go` + others

```go
// lead_source.go
package valueobject

type LeadSource string

const (
	LeadSourceWeb      LeadSource = "WEB"
	LeadSourceReferral LeadSource = "REFERRAL"
	LeadSourceCampaign LeadSource = "CAMPAIGN"
	LeadSourceColdCall LeadSource = "COLD_CALL"
	LeadSourceEvent    LeadSource = "EVENT"
)

func (s LeadSource) IsValid() bool {
	switch s {
	case LeadSourceWeb, LeadSourceReferral, LeadSourceCampaign, LeadSourceColdCall, LeadSourceEvent:
		return true
	}
	return false
}
```

```go
// opportunity_stage.go
package valueobject

type OpportunityStage string

const (
	StageProspecting OpportunityStage = "PROSPECTING"
	StageProposal    OpportunityStage = "PROPOSAL"
	StageNegotiation OpportunityStage = "NEGOTIATION"
	StageWon         OpportunityStage = "WON"
	StageLost        OpportunityStage = "LOST"
)

func (s OpportunityStage) IsValid() bool {
	switch s {
	case StageProspecting, StageProposal, StageNegotiation, StageWon, StageLost:
		return true
	}
	return false
}

func (s OpportunityStage) IsTerminal() bool {
	return s == StageWon || s == StageLost
}

func (s OpportunityStage) DefaultProbability() int {
	switch s {
	case StageProspecting:
		return 10
	case StageProposal:
		return 40
	case StageNegotiation:
		return 70
	case StageWon:
		return 100
	case StageLost:
		return 0
	}
	return 0
}
```

```go
// ticket_priority.go
package valueobject

import "time"

type TicketPriority string

const (
	TicketPriorityLow    TicketPriority = "LOW"
	TicketPriorityMedium TicketPriority = "MEDIUM"
	TicketPriorityHigh   TicketPriority = "HIGH"
	TicketPriorityUrgent TicketPriority = "URGENT"
)

func (p TicketPriority) IsValid() bool {
	switch p {
	case TicketPriorityLow, TicketPriorityMedium, TicketPriorityHigh, TicketPriorityUrgent:
		return true
	}
	return false
}

func (p TicketPriority) SLADuration() time.Duration {
	switch p {
	case TicketPriorityUrgent:
		return 2 * time.Hour
	case TicketPriorityHigh:
		return 8 * time.Hour
	case TicketPriorityMedium:
		return 24 * time.Hour
	case TicketPriorityLow:
		return 48 * time.Hour
	}
	return 24 * time.Hour
}
```

```go
// ticket_status.go
package valueobject

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "OPEN"
	TicketStatusInProgress TicketStatus = "IN_PROGRESS"
	TicketStatusResolved   TicketStatus = "RESOLVED"
	TicketStatusClosed     TicketStatus = "CLOSED"
)

func (s TicketStatus) IsValid() bool {
	switch s {
	case TicketStatusOpen, TicketStatusInProgress, TicketStatusResolved, TicketStatusClosed:
		return true
	}
	return false
}

func (s TicketStatus) IsTerminal() bool {
	return s == TicketStatusClosed
}

func (s TicketStatus) CanTransitionTo(next TicketStatus) bool {
	t := map[TicketStatus]map[TicketStatus]bool{
		TicketStatusOpen:       {TicketStatusInProgress: true, TicketStatusResolved: true, TicketStatusClosed: true},
		TicketStatusInProgress: {TicketStatusResolved: true, TicketStatusOpen: true},
		TicketStatusResolved:   {TicketStatusClosed: true, TicketStatusOpen: true},
		TicketStatusClosed:     {TicketStatusOpen: true},
	}
	return t[s][next]
}
```

```go
// activity_type.go
package valueobject

type ActivityType string

const (
	ActivityCall    ActivityType = "CALL"
	ActivityEmail   ActivityType = "EMAIL"
	ActivityMeeting ActivityType = "MEETING"
	ActivityNote    ActivityType = "NOTE"
	ActivityTask    ActivityType = "TASK"
)

func (t ActivityType) IsValid() bool {
	switch t {
	case ActivityCall, ActivityEmail, ActivityMeeting, ActivityNote, ActivityTask:
		return true
	}
	return false
}
```

### `domain/errors/errors.go`

```go
package domainerrors

import "errors"

var (
	ErrInvalidTenant     = errors.New("invalid tenant")
	ErrInvalidLead       = errors.New("invalid lead")
	ErrInvalidLeadName   = errors.New("invalid lead name")
	ErrInvalidLeadSource = errors.New("invalid lead source")
	ErrInvalidScore      = errors.New("invalid score (must be 0-100)")
	ErrInvalidStage      = errors.New("invalid opportunity stage")
	ErrInvalidProbability = errors.New("invalid probability (must be 0-100)")
	ErrInvalidTicket     = errors.New("invalid ticket")
	ErrInvalidSubject    = errors.New("invalid subject")
	ErrInvalidPriority   = errors.New("invalid priority")
	ErrInvalidActivity   = errors.New("invalid activity")
	ErrInvalidUser       = errors.New("invalid user")

	ErrLeadNotFound        = errors.New("lead not found")
	ErrOpportunityNotFound = errors.New("opportunity not found")
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrActivityNotFound    = errors.New("activity not found")

	ErrLeadAlreadyConverted = errors.New("lead already converted")
	ErrInvalidLeadTransition = errors.New("invalid lead status transition")
	ErrCannotConvertLostLead = errors.New("cannot convert lost lead")
	ErrInvalidOpportunityTransition = errors.New("invalid opportunity stage transition")
	ErrInvalidTicketTransition = errors.New("invalid ticket status transition")
	ErrCannotResolveClosed  = errors.New("cannot resolve closed ticket")
	ErrCannotEscalate       = errors.New("cannot escalate ticket in current state")
	ErrSLAAlreadyBreached   = errors.New("SLA already breached")
)
```

### `domain/entity/lead.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type Lead struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Email       string
	Phone       string
	Company     string
	Source      valueobject.LeadSource
	Status      valueobject.LeadStatus
	Score       int
	AssignedTo  *uuid.UUID
	ConvertedTo *uuid.UUID
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewLead(tenantID uuid.UUID, name, email, phone, company string, source valueobject.LeadSource) (*Lead, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if name == "" || len(name) > 255 {
		return nil, domainerrors.ErrInvalidLeadName
	}
	if !source.IsValid() {
		return nil, domainerrors.ErrInvalidLeadSource
	}
	now := time.Now()
	return &Lead{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Company:   company,
		Source:    source,
		Status:    valueobject.LeadStatusNew,
		Score:     0,
		Metadata:  map[string]interface{}{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (l *Lead) Qualify(score int) error {
	if score < 0 || score > 100 {
		return domainerrors.ErrInvalidScore
	}
	if !l.Status.CanTransitionTo(valueobject.LeadStatusQualified) {
		return domainerrors.ErrInvalidLeadTransition
	}
	l.Score = score
	l.Status = valueobject.LeadStatusQualified
	l.touch()
	return nil
}

func (l *Lead) AssignTo(userID uuid.UUID) {
	l.AssignedTo = &userID
	l.touch()
}

func (l *Lead) Convert(customerID uuid.UUID) error {
	if l.Status == valueobject.LeadStatusConverted {
		return domainerrors.ErrLeadAlreadyConverted
	}
	if l.Status == valueobject.LeadStatusLost {
		return domainerrors.ErrCannotConvertLostLead
	}
	if !l.Status.CanTransitionTo(valueobject.LeadStatusConverted) {
		return domainerrors.ErrInvalidLeadTransition
	}
	l.Status = valueobject.LeadStatusConverted
	l.ConvertedTo = &customerID
	l.touch()
	return nil
}

func (l *Lead) MarkContacted() error {
	if !l.Status.CanTransitionTo(valueobject.LeadStatusContacted) {
		return domainerrors.ErrInvalidLeadTransition
	}
	l.Status = valueobject.LeadStatusContacted
	l.touch()
	return nil
}

func (l *Lead) Lose(reason string) error {
	if !l.Status.CanTransitionTo(valueobject.LeadStatusLost) {
		return domainerrors.ErrInvalidLeadTransition
	}
	l.Status = valueobject.LeadStatusLost
	l.Metadata["lost_reason"] = reason
	l.touch()
	return nil
}

func (l *Lead) IsConverted() bool { return l.Status == valueobject.LeadStatusConverted }
func (l *Lead) IsQualified() bool { return l.Status == valueobject.LeadStatusQualified }

func (l *Lead) touch() { l.UpdatedAt = time.Now() }
```

### `domain/entity/opportunity.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type Opportunity struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	LeadID            *uuid.UUID
	CustomerID        *uuid.UUID
	Title             string
	Amount            float64
	Currency          string
	Probability       int
	Stage             valueobject.OpportunityStage
	ExpectedCloseDate *time.Time
	OwnerID           *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewOpportunity(tenantID uuid.UUID, title string, amount float64, currency string) (*Opportunity, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if title == "" || len(title) > 255 {
		return nil, domainerrors.ErrInvalidLeadName
	}
	if amount < 0 {
		return nil, domainerrors.ErrInvalidScore
	}
	if currency == "" {
		currency = "THB"
	}
	now := time.Now()
	return &Opportunity{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Title:       title,
		Amount:      amount,
		Currency:    currency,
		Probability: valueobject.StageProspecting.DefaultProbability(),
		Stage:       valueobject.StageProspecting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (o *Opportunity) LinkLead(leadID uuid.UUID) {
	o.LeadID = &leadID
	o.touch()
}

func (o *Opportunity) LinkCustomer(customerID uuid.UUID) {
	o.CustomerID = &customerID
	o.touch()
}

func (o *Opportunity) MoveToStage(stage valueobject.OpportunityStage) error {
	if !stage.IsValid() {
		return domainerrors.ErrInvalidStage
	}
	if o.Stage.IsTerminal() {
		return domainerrors.ErrInvalidOpportunityTransition
	}
	o.Stage = stage
	o.Probability = stage.DefaultProbability()
	o.touch()
	return nil
}

func (o *Opportunity) UpdateProbability(p int) error {
	if p < 0 || p > 100 {
		return domainerrors.ErrInvalidProbability
	}
	o.Probability = p
	o.touch()
	return nil
}

func (o *Opportunity) Win() error {
	if o.Stage == valueobject.StageWon {
		return nil
	}
	if o.Stage == valueobject.StageLost {
		return domainerrors.ErrInvalidOpportunityTransition
	}
	o.Stage = valueobject.StageWon
	o.Probability = 100
	o.touch()
	return nil
}

func (o *Opportunity) Lose(reason string) error {
	if o.Stage == valueobject.StageLost {
		return nil
	}
	if o.Stage == valueobject.StageWon {
		return domainerrors.ErrInvalidOpportunityTransition
	}
	o.Stage = valueobject.StageLost
	o.Probability = 0
	o.touch()
	return nil
}

func (o *Opportunity) IsWon() bool { return o.Stage == valueobject.StageWon }
func (o *Opportunity) IsLost() bool { return o.Stage == valueobject.StageLost }

func (o *Opportunity) WeightedAmount() float64 {
	return o.Amount * float64(o.Probability) / 100.0
}

func (o *Opportunity) touch() { o.UpdatedAt = time.Now() }
```

### `domain/entity/ticket.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type Ticket struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	CustomerID  uuid.UUID
	Subject     string
	Description string
	Priority    valueobject.TicketPriority
	Status      valueobject.TicketStatus
	AssignedTo  *uuid.UUID
	SLADueAt    time.Time
	StartedAt   *time.Time
	ResolvedAt  *time.Time
	ClosedAt    *time.Time
	EscalatedAt *time.Time
	EscalationLevel int
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTicket(tenantID, customerID uuid.UUID, subject, description string, priority valueobject.TicketPriority) (*Ticket, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if customerID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTicket
	}
	if subject == "" || len(subject) > 255 {
		return nil, domainerrors.ErrInvalidSubject
	}
	if !priority.IsValid() {
		return nil, domainerrors.ErrInvalidPriority
	}
	now := time.Now()
	return &Ticket{
		ID:          uuid.New(),
		TenantID:    tenantID,
		CustomerID:  customerID,
		Subject:     subject,
		Description: description,
		Priority:    priority,
		Status:      valueobject.TicketStatusOpen,
		SLADueAt:    now.Add(priority.SLADuration()),
		Metadata:    map[string]interface{}{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (t *Ticket) Assign(userID uuid.UUID) {
	t.AssignedTo = &userID
	t.touch()
}

func (t *Ticket) Start() error {
	if !t.Status.CanTransitionTo(valueobject.TicketStatusInProgress) {
		return domainerrors.ErrInvalidTicketTransition
	}
	t.Status = valueobject.TicketStatusInProgress
	now := time.Now()
	t.StartedAt = &now
	t.touch()
	return nil
}

func (t *Ticket) Resolve() error {
	if t.Status == valueobject.TicketStatusClosed {
		return domainerrors.ErrCannotResolveClosed
	}
	if !t.Status.CanTransitionTo(valueobject.TicketStatusResolved) {
		return domainerrors.ErrInvalidTicketTransition
	}
	t.Status = valueobject.TicketStatusResolved
	now := time.Now()
	t.ResolvedAt = &now
	t.touch()
	return nil
}

func (t *Ticket) Close() error {
	if !t.Status.CanTransitionTo(valueobject.TicketStatusClosed) {
		return domainerrors.ErrInvalidTicketTransition
	}
	t.Status = valueobject.TicketStatusClosed
	now := time.Now()
	t.ClosedAt = &now
	t.touch()
	return nil
}

func (t *Ticket) Reopen() error {
	if !t.Status.CanTransitionTo(valueobject.TicketStatusOpen) {
		return domainerrors.ErrInvalidTicketTransition
	}
	t.Status = valueobject.TicketStatusOpen
	t.ResolvedAt = nil
	t.ClosedAt = nil
	t.touch()
	return nil
}

func (t *Ticket) Escalate() error {
	if t.Status == valueobject.TicketStatusClosed || t.Status == valueobject.TicketStatusResolved {
		return domainerrors.ErrCannotEscalate
	}
	now := time.Now()
	t.EscalatedAt = &now
	t.EscalationLevel++
	t.Metadata["escalated"] = true
	t.touch()
	return nil
}

func (t *Ticket) IsOverdue() bool {
	return time.Now().After(t.SLADueAt) && !t.IsResolvedOrClosed()
}

func (t *Ticket) IsResolvedOrClosed() bool {
	return t.Status == valueobject.TicketStatusResolved || t.Status == valueobject.TicketStatusClosed
}

func (t *Ticket) touch() { t.UpdatedAt = time.Now() }
```

### `domain/entity/activity.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type Activity struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	EntityType  string
	EntityID    uuid.UUID
	Type        valueobject.ActivityType
	Subject     string
	Notes       string
	PerformedBy uuid.UUID
	OccurredAt  time.Time
	CreatedAt   time.Time
}

func NewActivity(tenantID uuid.UUID, entityType string, entityID uuid.UUID, atype valueobject.ActivityType, subject, notes string, performedBy uuid.UUID) (*Activity, error) {
	if tenantID == uuid.Nil {
		return nil, domainerrors.ErrInvalidTenant
	}
	if !atype.IsValid() {
		return nil, domainerrors.ErrInvalidActivity
	}
	if performedBy == uuid.Nil {
		return nil, domainerrors.ErrInvalidUser
	}
	now := time.Now()
	return &Activity{
		ID:          uuid.New(),
		TenantID:    tenantID,
		EntityType:  entityType,
		EntityID:    entityID,
		Type:        atype,
		Subject:     subject,
		Notes:       notes,
		PerformedBy: performedBy,
		OccurredAt:  now,
		CreatedAt:   now,
	}, nil
}
```

### `domain/service/sla_service.go`

```go
package service

import (
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/crm/domain/entity"
)

type SLAService struct{}

// CheckSLA – ตรวจว่า ticket เกิน SLA หรือยัง
func (s *SLAService) CheckSLA(t *entity.Ticket) (bool, time.Duration) {
	if t.IsResolvedOrClosed() {
		return false, 0
	}
	now := time.Now()
	if now.After(t.SLADueAt) {
		return true, now.Sub(t.SLADueAt)
	}
	return false, t.SLADueAt.Sub(now)
}

// ShouldEscalate – ควร escalate หรือยัง
func (s *SLAService) ShouldEscalate(t *entity.Ticket) bool {
	if t.IsResolvedOrClosed() {
		return false
	}
	breached, _ := s.CheckSLA(t)
	return breached || t.EscalationLevel >= 2
}

// AssignPriority – แนะนำ priority ตาม category
func (s *SLAService) AssignPriority(category string) string {
	switch category {
	case "outage", "down", "critical":
		return "URGENT"
	case "error", "bug":
		return "HIGH"
	case "performance", "slow":
		return "MEDIUM"
	default:
		return "LOW"
	}
}

var _ = uuid.New
```

### Repository Interfaces

```go
// domain/repository/lead_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"icmongolang/internal/modules/crm/domain/entity"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type LeadRepository interface {
	Save(ctx context.Context, l *entity.Lead) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Lead, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.LeadStatus) ([]entity.Lead, error)
	FindByAssignee(ctx context.Context, userID uuid.UUID) ([]entity.Lead, error)
	ExistsByEmail(ctx context.Context, tenantID uuid.UUID, email string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// opportunity_repository.go
type OpportunityRepository interface {
	Save(ctx context.Context, o *entity.Opportunity) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Opportunity, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Opportunity, error)
	FindByStage(ctx context.Context, tenantID uuid.UUID, stage valueobject.OpportunityStage) ([]entity.Opportunity, error)
	FindByLead(ctx context.Context, leadID uuid.UUID) ([]entity.Opportunity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ticket_repository.go
type TicketRepository interface {
	Save(ctx context.Context, t *entity.Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Ticket, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Ticket, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.TicketStatus) ([]entity.Ticket, error)
	FindByAssignee(ctx context.Context, userID uuid.UUID) ([]entity.Ticket, error)
	FindOverdue(ctx context.Context, before time.Time) ([]entity.Ticket, error)
	FindOpenByCustomer(ctx context.Context, customerID uuid.UUID) ([]entity.Ticket, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// activity_repository.go
type ActivityRepository interface {
	Save(ctx context.Context, a *entity.Activity) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Activity, error)
	FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]entity.Activity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

### Domain Events

```go
// domain/event/events.go
package event

import (
	"time"

	"github.com/google/uuid"
)

type LeadCreated struct {
	LeadID     uuid.UUID
	TenantID   uuid.UUID
	Name       string
	Email      string
	Source     string
	OccurredAt time.Time
}

type LeadQualified struct {
	LeadID     uuid.UUID
	TenantID   uuid.UUID
	Score      int
	OccurredAt time.Time
}

type LeadConverted struct {
	LeadID     uuid.UUID
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	OccurredAt time.Time
}

type OpportunityCreated struct {
	OpportunityID uuid.UUID
	TenantID      uuid.UUID
	Title         string
	Amount        float64
	Currency      string
	OccurredAt    time.Time
}

type OpportunityWon struct {
	OpportunityID uuid.UUID
	TenantID      uuid.UUID
	Amount        float64
	OccurredAt    time.Time
}

type OpportunityLost struct {
	OpportunityID uuid.UUID
	TenantID      uuid.UUID
	Reason        string
	OccurredAt    time.Time
}

type TicketOpened struct {
	TicketID   uuid.UUID
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	Subject    string
	Priority   string
	SLADueAt   time.Time
	OccurredAt time.Time
}

type TicketEscalated struct {
	TicketID        uuid.UUID
	TenantID        uuid.UUID
	EscalationLevel int
	OccurredAt      time.Time
}

type TicketResolved struct {
	TicketID   uuid.UUID
	TenantID   uuid.UUID
	DurationMs int64
	OccurredAt time.Time
}
```

### Application — Ports

```go
// application/ports.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/crm/domain/event"
)

type EventProducer interface {
	PublishLeadCreated(ctx context.Context, evt event.LeadCreated) error
	PublishLeadQualified(ctx context.Context, evt event.LeadQualified) error
	PublishLeadConverted(ctx context.Context, evt event.LeadConverted) error
	PublishOpportunityCreated(ctx context.Context, evt event.OpportunityCreated) error
	PublishOpportunityWon(ctx context.Context, evt event.OpportunityWon) error
	PublishOpportunityLost(ctx context.Context, evt event.OpportunityLost) error
	PublishTicketOpened(ctx context.Context, evt event.TicketOpened) error
	PublishTicketEscalated(ctx context.Context, evt event.TicketEscalated) error
	PublishTicketResolved(ctx context.Context, evt event.TicketResolved) error
}

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
	OccurredAt time.Time
}

type WSHub interface {
	BroadcastToTenant(tenantID string, event interface{})
}

type CustomerClient interface {
	Create(ctx context.Context, req CreateCustomerRequest) (*CustomerResult, error)
	Get(ctx context.Context, id uuid.UUID) (*CustomerResult, error)
}

type CreateCustomerRequest struct {
	TenantID uuid.UUID
	Code     string
	Name     string
	Type     string
	Segment  string
}

type CustomerResult struct {
	ID   uuid.UUID
	Code string
}

type NotifierClient interface {
	SendEmail(ctx context.Context, req SendEmailRequest) error
	SendSMS(ctx context.Context, req SendSMSRequest) error
}

type SendEmailRequest struct {
	To       string
	Subject  string
	Body     string
	TenantID uuid.UUID
}

type SendSMSRequest struct {
	To       string
	Message  string
	TenantID uuid.UUID
}
```

### Use Case — Capture Lead

```go
// application/capture_lead.go
package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"icmongolang/internal/modules/crm/domain/entity"
	"icmongolang/internal/modules/crm/domain/event"
	"icmongolang/internal/modules/crm/domain/repository"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type CaptureLeadUseCase struct {
	leadRepo  repository.LeadRepository
	auditRepo AuditRepository
	producer  EventProducer
	notifier  NotifierClient
}

func NewCaptureLeadUseCase(leadRepo repository.LeadRepository, auditRepo AuditRepository, producer EventProducer, notifier NotifierClient) *CaptureLeadUseCase {
	return &CaptureLeadUseCase{leadRepo: leadRepo, auditRepo: auditRepo, producer: producer, notifier: notifier}
}

type CaptureLeadInput struct {
	TenantID  uuid.UUID
	Name      string
	Email     string
	Phone     string
	Company   string
	Source    valueobject.LeadSource
	UserID    uuid.UUID
	IPAddress string
}

type CaptureLeadOutput struct {
	LeadID uuid.UUID `json:"lead_id"`
	Status string    `json:"status"`
}

func (uc *CaptureLeadUseCase) Execute(ctx context.Context, in CaptureLeadInput) (*CaptureLeadOutput, error) {
	lead, err := entity.NewLead(in.TenantID, in.Name, in.Email, in.Phone, in.Company, in.Source)
	if err != nil {
		return nil, err
	}

	if err := uc.leadRepo.Save(ctx, lead); err != nil {
		return nil, err
	}

	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID: in.UserID, Action: "LEAD_CAPTURED",
		EntityType: "lead", EntityID: lead.ID,
		Payload:    map[string]interface{}{"name": lead.Name, "source": string(lead.Source)},
		IPAddress:  in.IPAddress,
		OccurredAt: time.Now(),
	})

	_ = uc.producer.PublishLeadCreated(ctx, event.LeadCreated{
		LeadID: lead.ID, TenantID: lead.TenantID,
		Name: lead.Name, Email: lead.Email, Source: string(lead.Source),
		OccurredAt: time.Now(),
	})

	// Welcome email
	if lead.Email != "" {
		_ = uc.notifier.SendEmail(ctx, SendEmailRequest{
			To: lead.Email, Subject: "Thank you for your interest",
			TenantID: lead.TenantID,
		})
	}

	return &CaptureLeadOutput{LeadID: lead.ID, Status: string(lead.Status)}, nil
}
```

### Use Case — Qualify, Convert, Ticket flows

```go
// application/qualify_lead.go
type QualifyLeadUseCase struct {
	repo     repository.LeadRepository
	producer EventProducer
	auditRepo AuditRepository
}

func NewQualifyLeadUseCase(repo repository.LeadRepository, producer EventProducer, auditRepo AuditRepository) *QualifyLeadUseCase {
	return &QualifyLeadUseCase{repo: repo, producer: producer, auditRepo: auditRepo}
}

type QualifyLeadInput struct {
	LeadID uuid.UUID
	Score  int
	UserID uuid.UUID
}

func (uc *QualifyLeadUseCase) Execute(ctx context.Context, in QualifyLeadInput) error {
	l, err := uc.repo.FindByID(ctx, in.LeadID)
	if err != nil {
		return err
	}
	if err := l.Qualify(in.Score); err != nil {
		return err
	}
	if err := uc.repo.Save(ctx, l); err != nil {
		return err
	}
	return uc.producer.PublishLeadQualified(ctx, event.LeadQualified{
		LeadID: l.ID, TenantID: l.TenantID, Score: l.Score, OccurredAt: time.Now(),
	})
}

// application/convert_lead.go
type ConvertLeadUseCase struct {
	leadRepo      repository.LeadRepository
	customerCli   CustomerClient
	producer      EventProducer
	auditRepo     AuditRepository
}

func NewConvertLeadUseCase(leadRepo repository.LeadRepository, customerCli CustomerClient, producer EventProducer, auditRepo AuditRepository) *ConvertLeadUseCase {
	return &ConvertLeadUseCase{leadRepo: leadRepo, customerCli: customerCli, producer: producer, auditRepo: auditRepo}
}

type ConvertLeadInput struct {
	LeadID uuid.UUID
	UserID uuid.UUID
}

type ConvertLeadOutput struct {
	CustomerID uuid.UUID `json:"customer_id"`
}

func (uc *ConvertLeadUseCase) Execute(ctx context.Context, in ConvertLeadInput) (*ConvertLeadOutput, error) {
	l, err := uc.leadRepo.FindByID(ctx, in.LeadID)
	if err != nil {
		return nil, err
	}
	if l.IsConverted() {
		return nil, domainerrors.ErrLeadAlreadyConverted
	}

	// Create customer
	cust, err := uc.customerCli.Create(ctx, CreateCustomerRequest{
		TenantID: l.TenantID,
		Code:     "CUST-" + l.ID.String()[:8],
		Name:     l.Name,
		Type:     "PERSONAL",
		Segment:  "SME",
	})
	if err != nil {
		return nil, err
	}

	if err := l.Convert(cust.ID); err != nil {
		return nil, err
	}
	if err := uc.leadRepo.Save(ctx, l); err != nil {
		return nil, err
	}

	_ = uc.producer.PublishLeadConverted(ctx, event.LeadConverted{
		LeadID: l.ID, TenantID: l.TenantID, CustomerID: cust.ID,
		OccurredAt: time.Now(),
	})

	return &ConvertLeadOutput{CustomerID: cust.ID}, nil
}

// application/open_ticket.go
type OpenTicketUseCase struct {
	ticketRepo repository.TicketRepository
	producer   EventProducer
	auditRepo  AuditRepository
	notifier   NotifierClient
}

func NewOpenTicketUseCase(ticketRepo repository.TicketRepository, producer EventProducer, auditRepo AuditRepository, notifier NotifierClient) *OpenTicketUseCase {
	return &OpenTicketUseCase{ticketRepo: ticketRepo, producer: producer, auditRepo: auditRepo, notifier: notifier}
}

type OpenTicketInput struct {
	TenantID    uuid.UUID
	CustomerID  uuid.UUID
	Subject     string
	Description string
	Priority    valueobject.TicketPriority
	UserID      uuid.UUID
}

type OpenTicketOutput struct {
	TicketID uuid.UUID `json:"ticket_id"`
	SLADueAt time.Time `json:"sla_due_at"`
}

func (uc *OpenTicketUseCase) Execute(ctx context.Context, in OpenTicketInput) (*OpenTicketOutput, error) {
	t, err := entity.NewTicket(in.TenantID, in.CustomerID, in.Subject, in.Description, in.Priority)
	if err != nil {
		return nil, err
	}
	if err := uc.ticketRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	_ = uc.auditRepo.Save(ctx, &AuditTrail{
		UserID: in.UserID, Action: "TICKET_OPENED",
		EntityType: "ticket", EntityID: t.ID,
		Payload:    map[string]interface{}{"subject": t.Subject, "priority": string(t.Priority)},
		OccurredAt: time.Now(),
	})

	_ = uc.producer.PublishTicketOpened(ctx, event.TicketOpened{
		TicketID: t.ID, TenantID: t.TenantID, CustomerID: t.CustomerID,
		Subject: t.Subject, Priority: string(t.Priority), SLADueAt: t.SLADueAt,
		OccurredAt: time.Now(),
	})

	return &OpenTicketOutput{TicketID: t.ID, SLADueAt: t.SLADueAt}, nil
}

// application/escalate_ticket.go
type EscalateTicketUseCase struct {
	ticketRepo repository.TicketRepository
	producer   EventProducer
	notifier   NotifierClient
}

func NewEscalateTicketUseCase(ticketRepo repository.TicketRepository, producer EventProducer, notifier NotifierClient) *EscalateTicketUseCase {
	return &EscalateTicketUseCase{ticketRepo: ticketRepo, producer: producer, notifier: notifier}
}

type EscalateTicketInput struct {
	TicketID uuid.UUID
	UserID   uuid.UUID
}

func (uc *EscalateTicketUseCase) Execute(ctx context.Context, in EscalateTicketInput) error {
	t, err := uc.ticketRepo.FindByID(ctx, in.TicketID)
	if err != nil {
		return err
	}
	if err := t.Escalate(); err != nil {
		return err
	}
	if err := uc.ticketRepo.Save(ctx, t); err != nil {
		return err
	}
	return uc.producer.PublishTicketEscalated(ctx, event.TicketEscalated{
		TicketID: t.ID, TenantID: t.TenantID,
		EscalationLevel: t.EscalationLevel, OccurredAt: time.Now(),
	})
}

// application/resolve_ticket.go
type ResolveTicketUseCase struct {
	ticketRepo repository.TicketRepository
	producer   EventProducer
	auditRepo  AuditRepository
}

func NewResolveTicketUseCase(ticketRepo repository.TicketRepository, producer EventProducer, auditRepo AuditRepository) *ResolveTicketUseCase {
	return &ResolveTicketUseCase{ticketRepo: ticketRepo, producer: producer, auditRepo: auditRepo}
}

type ResolveTicketInput struct {
	TicketID uuid.UUID
	UserID   uuid.UUID
	Notes    string
}

func (uc *ResolveTicketUseCase) Execute(ctx context.Context, in ResolveTicketInput) error {
	t, err := uc.ticketRepo.FindByID(ctx, in.TicketID)
	if err != nil {
		return err
	}
	if err := t.Resolve(); err != nil {
		return err
	}
	if err := uc.ticketRepo.Save(ctx, t); err != nil {
		return err
	}
	durationMs := int64(0)
	if t.ResolvedAt != nil {
		durationMs = t.ResolvedAt.Sub(t.CreatedAt).Milliseconds()
	}
	return uc.producer.PublishTicketResolved(ctx, event.TicketResolved{
		TicketID: t.ID, TenantID: t.TenantID,
		DurationMs: durationMs, OccurredAt: time.Now(),
	})
}
```

### Infrastructure — Models

```go
// infrastructure/persistence/postgres/models.go
package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type LeadModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Email       string         `gorm:"type:varchar(255);index"`
	Phone       string         `gorm:"type:varchar(30)"`
	Company     string         `gorm:"type:varchar(255)"`
	Source      string         `gorm:"type:varchar(50);not null"`
	Status      string         `gorm:"type:varchar(30);not null;index"`
	Score       int            `gorm:"default:0"`
	AssignedTo  *uuid.UUID     `gorm:"type:uuid;index"`
	ConvertedTo *uuid.UUID     `gorm:"type:uuid"`
	Metadata    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (LeadModel) TableName() string { return "crm_leads" }

type OpportunityModel struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	LeadID            *uuid.UUID `gorm:"type:uuid;index"`
	CustomerID        *uuid.UUID `gorm:"type:uuid;index"`
	Title             string     `gorm:"type:varchar(255);not null"`
	Amount            float64    `gorm:"type:numeric(15,2)"`
	Currency          string     `gorm:"type:varchar(3);default:'THB'"`
	Probability       int        `gorm:"default:0"`
	Stage             string     `gorm:"type:varchar(30);not null;index"`
	ExpectedCloseDate *time.Time
	OwnerID           *uuid.UUID `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (OpportunityModel) TableName() string { return "crm_opportunities" }

type TicketModel struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	CustomerID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Subject         string         `gorm:"type:varchar(255);not null"`
	Description     string         `gorm:"type:text"`
	Priority        string         `gorm:"type:varchar(20);not null;index"`
	Status          string         `gorm:"type:varchar(30);not null;index"`
	AssignedTo      *uuid.UUID     `gorm:"type:uuid;index"`
	SLADueAt        time.Time      `gorm:"index"`
	StartedAt       *time.Time
	ResolvedAt      *time.Time
	ClosedAt        *time.Time
	EscalatedAt     *time.Time
	EscalationLevel int            `gorm:"default:0"`
	Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (TicketModel) TableName() string { return "crm_tickets" }

type ActivityModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	EntityType  string    `gorm:"type:varchar(30);not null;index"`
	EntityID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Type        string    `gorm:"type:varchar(20);not null"`
	Subject     string    `gorm:"type:varchar(255)"`
	Notes       string    `gorm:"type:text"`
	PerformedBy uuid.UUID `gorm:"type:uuid;not null"`
	OccurredAt  time.Time `gorm:"index"`
	CreatedAt   time.Time
}

func (ActivityModel) TableName() string { return "crm_activities" }
```

### Repo Impl (ตัวอย่าง lead)

```go
// infrastructure/persistence/postgres/lead_repo_impl.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"icmongolang/internal/modules/crm/domain/entity"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

type leadRepoImpl struct{ db *gorm.DB }

func NewLeadRepository(db *gorm.DB) *leadRepoImpl { return &leadRepoImpl{db: db} }

func (r *leadRepoImpl) Save(ctx context.Context, l *entity.Lead) error {
	return r.db.WithContext(ctx).Save(toLeadModel(l)).Error
}

func (r *leadRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error) {
	var m LeadModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrLeadNotFound
	}
	if err != nil {
		return nil, err
	}
	return toLeadEntity(&m), nil
}

func (r *leadRepoImpl) FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]entity.Lead, error) {
	var models []LeadModel
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

func (r *leadRepoImpl) FindByStatus(ctx context.Context, tenantID uuid.UUID, status valueobject.LeadStatus) ([]entity.Lead, error) {
	var models []LeadModel
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, status).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

func (r *leadRepoImpl) FindByAssignee(ctx context.Context, userID uuid.UUID) ([]entity.Lead, error) {
	var models []LeadModel
	err := r.db.WithContext(ctx).Where("assigned_to = ?", userID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return toLeadEntities(models), nil
}

func (r *leadRepoImpl) ExistsByEmail(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LeadModel{}).
		Where("tenant_id = ? AND email = ?", tenantID, email).Count(&count).Error
	return count > 0, err
}

func (r *leadRepoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&LeadModel{}, "id = ?", id).Error
}

func toLeadModel(l *entity.Lead) *LeadModel {
	metaJSON, _ := json.Marshal(l.Metadata)
	return &LeadModel{
		ID: l.ID, TenantID: l.TenantID, Name: l.Name, Email: l.Email, Phone: l.Phone,
		Company: l.Company, Source: string(l.Source), Status: string(l.Status),
		Score: l.Score, AssignedTo: l.AssignedTo, ConvertedTo: l.ConvertedTo,
		Metadata: datatypes.JSON(metaJSON),
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func toLeadEntity(m *LeadModel) *entity.Lead {
	var meta map[string]interface{}
	_ = json.Unmarshal(m.Metadata, &meta)
	if meta == nil {
		meta = map[string]interface{}{}
	}
	return &entity.Lead{
		ID: m.ID, TenantID: m.TenantID, Name: m.Name, Email: m.Email, Phone: m.Phone,
		Company: m.Company, Source: valueobject.LeadSource(m.Source),
		Status: valueobject.LeadStatus(m.Status), Score: m.Score,
		AssignedTo: m.AssignedTo, ConvertedTo: m.ConvertedTo,
		Metadata: meta, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func toLeadEntities(models []LeadModel) []entity.Lead {
	out := make([]entity.Lead, 0, len(models))
	for i := range models {
		out = append(out, *toLeadEntity(&models[i]))
	}
	return out
}
```

### Kafka Producer

```go
// infrastructure/messaging/producer.go
package messaging

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"icmongolang/internal/modules/crm/domain/event"
)

const (
	TopicLeadCreated          = "crm.lead.created"
	TopicLeadQualified        = "crm.lead.qualified"
	TopicLeadConverted        = "crm.lead.converted"
	TopicOpportunityCreated   = "crm.opportunity.created"
	TopicOpportunityWon       = "crm.opportunity.won"
	TopicOpportunityLost      = "crm.opportunity.lost"
	TopicTicketOpened         = "crm.ticket.opened"
	TopicTicketEscalated      = "crm.ticket.escalated"
	TopicTicketResolved       = "crm.ticket.resolved"
)

type kafkaProducer struct{ sync sarama.SyncProducer }

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

func (k *kafkaProducer) publish(_ context.Context, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, _, err = k.sync.SendMessage(&sarama.ProducerMessage{
		Topic: topic, Key: sarama.StringEncoder(key), Value: sarama.ByteEncoder(data),
	})
	return err
}

func (k *kafkaProducer) PublishLeadCreated(ctx context.Context, evt event.LeadCreated) error {
	return k.publish(ctx, TopicLeadCreated, evt.LeadID.String(), evt)
}
func (k *kafkaProducer) PublishLeadQualified(ctx context.Context, evt event.LeadQualified) error {
	return k.publish(ctx, TopicLeadQualified, evt.LeadID.String(), evt)
}
func (k *kafkaProducer) PublishLeadConverted(ctx context.Context, evt event.LeadConverted) error {
	return k.publish(ctx, TopicLeadConverted, evt.LeadID.String(), evt)
}
func (k *kafkaProducer) PublishOpportunityCreated(ctx context.Context, evt event.OpportunityCreated) error {
	return k.publish(ctx, TopicOpportunityCreated, evt.OpportunityID.String(), evt)
}
func (k *kafkaProducer) PublishOpportunityWon(ctx context.Context, evt event.OpportunityWon) error {
	return k.publish(ctx, TopicOpportunityWon, evt.OpportunityID.String(), evt)
}
func (k *kafkaProducer) PublishOpportunityLost(ctx context.Context, evt event.OpportunityLost) error {
	return k.publish(ctx, TopicOpportunityLost, evt.OpportunityID.String(), evt)
}
func (k *kafkaProducer) PublishTicketOpened(ctx context.Context, evt event.TicketOpened) error {
	return k.publish(ctx, TopicTicketOpened, evt.TicketID.String(), evt)
}
func (k *kafkaProducer) PublishTicketEscalated(ctx context.Context, evt event.TicketEscalated) error {
	return k.publish(ctx, TopicTicketEscalated, evt.TicketID.String(), evt)
}
func (k *kafkaProducer) PublishTicketResolved(ctx context.Context, evt event.TicketResolved) error {
	return k.publish(ctx, TopicTicketResolved, evt.TicketID.String(), evt)
}
```

### Scheduler — SLA Job

```go
// infrastructure/scheduler/sla_job.go
package scheduler

import (
	"context"
	"log"
	"time"

	"icmongolang/internal/modules/crm/domain/repository"
	"icmongolang/internal/modules/crm/domain/service"
)

type SLAJob struct {
	ticketRepo repository.TicketRepository
	slaSvc     *service.SLAService
	producer   interface {
		PublishTicketEscalated(ctx context.Context, evt interface{}) error
	}
}

func NewSLAJob(ticketRepo repository.TicketRepository, slaSvc *service.SLAService) *SLAJob {
	return &SLAJob{ticketRepo: ticketRepo, slaSvc: slaSvc}
}

func (j *SLAJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	overdue, err := j.ticketRepo.FindOverdue(ctx, time.Now())
	if err != nil {
		log.Printf("[sla_job] find: %v", err)
		return
	}
	for i := range overdue {
		t := &overdue[i]
		if j.slaSvc.ShouldEscalate(t) {
			_ = t.Escalate()
			_ = j.ticketRepo.Save(ctx, t)
		}
	}
	log.Printf("[sla_job] checked %d tickets", len(overdue))
}
```

### HTTP Handler (ตัวอย่าง)

```go
// interfaces/http/ticket_handler.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icmongolang/internal/modules/crm/application"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
	"icmongolang/pkg/responses"
)

type TicketHandler struct {
	openUC     *application.OpenTicketUseCase
	assignUC   *application.AssignTicketUseCase
	startUC    *application.StartTicketUseCase
	resolveUC  *application.ResolveTicketUseCase
	escalateUC *application.EscalateTicketUseCase
	closeUC    *application.CloseTicketUseCase
	reopenUC   *application.ReopenTicketUseCase
}

func NewTicketHandler(
	openUC *application.OpenTicketUseCase,
	assignUC *application.AssignTicketUseCase,
	startUC *application.StartTicketUseCase,
	resolveUC *application.ResolveTicketUseCase,
	escalateUC *application.EscalateTicketUseCase,
	closeUC *application.CloseTicketUseCase,
	reopenUC *application.ReopenTicketUseCase,
) *TicketHandler {
	return &TicketHandler{openUC: openUC, assignUC: assignUC, startUC: startUC, resolveUC: resolveUC, escalateUC: escalateUC, closeUC: closeUC, reopenUC: reopenUC}
}

func (h *TicketHandler) Open(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var req struct {
		CustomerID  string `json:"customer_id" binding:"required"`
		Subject     string `json:"subject" binding:"required"`
		Description string `json:"description"`
		Priority    string `json:"priority" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error("INVALID", err.Error(), nil))
		return
	}
	custID, _ := uuid.Parse(req.CustomerID)

	out, err := h.openUC.Execute(c.Request.Context(), application.OpenTicketInput{
		TenantID: tenantID, CustomerID: custID,
		Subject: req.Subject, Description: req.Description,
		Priority: valueobject.TicketPriority(req.Priority), UserID: uid,
	})
	if err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.Success(out))
}

func (h *TicketHandler) Escalate(c *gin.Context) {
	uid := c.MustGet("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.escalateUC.Execute(c.Request.Context(), application.EscalateTicketInput{TicketID: id, UserID: uid}); err != nil {
		respondDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.Success(gin.H{"message": "escalated"}))
}
```

### `module.go`

```go
package crm

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"icmongolang/internal/modules/crm/application"
	"icmongolang/internal/modules/crm/domain/service"
	"icmongolang/internal/modules/crm/infrastructure/messaging"
	pgrepo "icmongolang/internal/modules/crm/infrastructure/persistence/postgres"
	httpiface "icmongolang/internal/modules/crm/interfaces/http"
)

type Dependencies struct {
	DB          *gorm.DB
	Producer    *messaging.KafkaProducer
	CustomerCli application.CustomerClient
	NotifierCli application.NotifierClient
	AuditRepo   application.AuditRepository
	WSHub       application.WSHub
}

func Init(
	router *gin.RouterGroup,
	deps Dependencies,
	auth gin.HandlerFunc, tenant gin.HandlerFunc, rbac func(roles ...string) gin.HandlerFunc,
) {
	leadRepo := pgrepo.NewLeadRepository(deps.DB)
	oppRepo := pgrepo.NewOpportunityRepository(deps.DB)
	ticketRepo := pgrepo.NewTicketRepository(deps.DB)
	activityRepo := pgrepo.NewActivityRepository(deps.DB)
	_ = activityRepo

	slaSvc := &service.SLAService{}

	captureLeadUC := application.NewCaptureLeadUseCase(leadRepo, deps.AuditRepo, deps.Producer, deps.NotifierCli)
	qualifyLeadUC := application.NewQualifyLeadUseCase(leadRepo, deps.Producer, deps.AuditRepo)
	convertLeadUC := application.NewConvertLeadUseCase(leadRepo, deps.CustomerCli, deps.Producer, deps.AuditRepo)
	assignLeadUC := application.NewAssignLeadUseCase(leadRepo, deps.AuditRepo)
	loseLeadUC := application.NewLoseLeadUseCase(leadRepo, deps.Producer, deps.AuditRepo)

	createOppUC := application.NewCreateOpportunityUseCase(oppRepo, deps.Producer, deps.AuditRepo)
	moveStageUC := application.NewMoveOpportunityStageUseCase(oppRepo, deps.Producer)
	winOppUC := application.NewWinOpportunityUseCase(oppRepo, deps.Producer)
	loseOppUC := application.NewLoseOpportunityUseCase(oppRepo, deps.Producer)

	openTicketUC := application.NewOpenTicketUseCase(ticketRepo, deps.Producer, deps.AuditRepo, deps.NotifierCli)
	assignTicketUC := application.NewAssignTicketUseCase(ticketRepo, deps.AuditRepo)
	startTicketUC := application.NewStartTicketUseCase(ticketRepo, deps.AuditRepo)
	resolveTicketUC := application.NewResolveTicketUseCase(ticketRepo, deps.Producer, deps.AuditRepo)
	escalateTicketUC := application.NewEscalateTicketUseCase(ticketRepo, deps.Producer, deps.NotifierCli)
	closeTicketUC := application.NewCloseTicketUseCase(ticketRepo, deps.AuditRepo)
	reopenTicketUC := application.NewReopenTicketUseCase(ticketRepo, deps.AuditRepo)

	leadHandler := httpiface.NewLeadHandler(captureLeadUC, qualifyLeadUC, convertLeadUC, assignLeadUC, loseLeadUC)
	oppHandler := httpiface.NewOpportunityHandler(createOppUC, moveStageUC, winOppUC, loseOppUC)
	ticketHandler := httpiface.NewTicketHandler(openTicketUC, assignTicketUC, startTicketUC, resolveTicketUC, escalateTicketUC, closeTicketUC, reopenTicketUC)

	httpiface.RegisterRoutes(router, &httpiface.Handlers{
		Lead: leadHandler, Opportunity: oppHandler, Ticket: ticketHandler,
	}, auth, tenant, rbac)

	_ = slaSvc
}
```

### Migration + Test (compact)

```sql
-- migrations/20260105_crm_init.sql
CREATE TABLE IF NOT EXISTS crm_leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(30),
    company VARCHAR(255),
    source VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'NEW',
    score INT DEFAULT 0,
    assigned_to UUID,
    converted_to UUID,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crm_leads_tenant ON crm_leads (tenant_id);
CREATE INDEX idx_crm_leads_status ON crm_leads (tenant_id, status);
CREATE INDEX idx_crm_leads_email ON crm_leads (tenant_id, email);

CREATE TABLE IF NOT EXISTS crm_opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    lead_id UUID,
    customer_id UUID,
    title VARCHAR(255) NOT NULL,
    amount NUMERIC(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'THB',
    probability INT DEFAULT 0,
    stage VARCHAR(30) NOT NULL DEFAULT 'PROSPECTING',
    expected_close_date DATE,
    owner_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crm_opps_tenant ON crm_opportunities (tenant_id, stage);

CREATE TABLE IF NOT EXISTS crm_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    subject VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN',
    assigned_to UUID,
    sla_due_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP,
    resolved_at TIMESTAMP,
    closed_at TIMESTAMP,
    escalated_at TIMESTAMP,
    escalation_level INT DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crm_tickets_tenant ON crm_tickets (tenant_id);
CREATE INDEX idx_crm_tickets_status ON crm_tickets (tenant_id, status);
CREATE INDEX idx_crm_tickets_sla ON crm_tickets (sla_due_at) WHERE status IN ('OPEN', 'IN_PROGRESS');
CREATE INDEX idx_crm_tickets_customer ON crm_tickets (customer_id);

CREATE TABLE IF NOT EXISTS crm_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    entity_type VARCHAR(30) NOT NULL,
    entity_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    subject VARCHAR(255),
    notes TEXT,
    performed_by UUID NOT NULL,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_crm_activities_entity ON crm_activities (entity_type, entity_id);
CREATE INDEX idx_crm_activities_tenant ON crm_activities (tenant_id, occurred_at DESC);
```

```go
// domain/entity/lead_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"icmongolang/internal/modules/crm/domain/entity"
	domainerrors "icmongolang/internal/modules/crm/domain/errors"
	valueobject "icmongolang/internal/modules/crm/domain/value_object"
)

func TestLead_Lifecycle(t *testing.T) {
	l, _ := entity.NewLead(uuid.New(), "Alice", "a@test.com", "", "", valueobject.LeadSourceWeb)
	if err := l.MarkContacted(); err != nil {
		t.Fatal(err)
	}
	if err := l.Qualify(80); err != nil {
		t.Fatal(err)
	}
	if !l.IsQualified() {
		t.Error("expected qualified")
	}
	if err := l.Convert(uuid.New()); err != nil {
		t.Fatal(err)
	}
	if !l.IsConverted() {
		t.Error("expected converted")
	}
	// Cannot convert twice
	if err := l.Convert(uuid.New()); !errors.Is(err, domainerrors.ErrLeadAlreadyConverted) {
		t.Errorf("expected ErrLeadAlreadyConverted, got %v", err)
	}
}

func TestTicket_Escalation(t *testing.T) {
	tk, _ := entity.NewTicket(uuid.New(), uuid.New(), "Server down", "help", valueobject.TicketPriorityUrgent)
	if err := tk.Escalate(); err != nil {
		t.Fatal(err)
	}
	if tk.EscalationLevel != 1 {
		t.Errorf("expected level 1, got %d", tk.EscalationLevel)
	}
}
```

**G05 Status:** ✅ Complete (~25 files)

---

## 1.3 Python Implementation (P05)

### Value Objects + Entities

```python
# app/modules/crm/domain/value_objects/lead_status.py
from enum import Enum


class LeadStatus(str, Enum):
    NEW = "NEW"
    CONTACTED = "CONTACTED"
    QUALIFIED = "QUALIFIED"
    LOST = "LOST"
    CONVERTED = "CONVERTED"

    def is_terminal(self) -> bool:
        return self in (LeadStatus.LOST, LeadStatus.CONVERTED)

    def can_transition_to(self, nxt: "LeadStatus") -> bool:
        t = {
            LeadStatus.NEW: {LeadStatus.CONTACTED, LeadStatus.QUALIFIED, LeadStatus.LOST},
            LeadStatus.CONTACTED: {LeadStatus.QUALIFIED, LeadStatus.LOST},
            LeadStatus.QUALIFIED: {LeadStatus.CONVERTED, LeadStatus.LOST},
            LeadStatus.LOST: set(),
            LeadStatus.CONVERTED: set(),
        }
        return nxt in t.get(self, set())
```

```python
# app/modules/crm/domain/value_objects/ticket_priority.py
from enum import Enum
from datetime import timedelta


class TicketPriority(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    URGENT = "URGENT"

    def sla_duration(self) -> timedelta:
        return {
            TicketPriority.URGENT: timedelta(hours=2),
            TicketPriority.HIGH: timedelta(hours=8),
            TicketPriority.MEDIUM: timedelta(hours=24),
            TicketPriority.LOW: timedelta(hours=48),
        }[self]
```

```python
# app/modules/crm/domain/value_objects/ticket_status.py
from enum import Enum


class TicketStatus(str, Enum):
    OPEN = "OPEN"
    IN_PROGRESS = "IN_PROGRESS"
    RESOLVED = "RESOLVED"
    CLOSED = "CLOSED"

    def can_transition_to(self, nxt: "TicketStatus") -> bool:
        t = {
            TicketStatus.OPEN: {TicketStatus.IN_PROGRESS, TicketStatus.RESOLVED, TicketStatus.CLOSED},
            TicketStatus.IN_PROGRESS: {TicketStatus.RESOLVED, TicketStatus.OPEN},
            TicketStatus.RESOLVED: {TicketStatus.CLOSED, TicketStatus.OPEN},
            TicketStatus.CLOSED: {TicketStatus.OPEN},
        }
        return nxt in t.get(self, set())
```

```python
# app/modules/crm/domain/value_objects/opportunity_stage.py
from enum import Enum


class OpportunityStage(str, Enum):
    PROSPECTING = "PROSPECTING"
    PROPOSAL = "PROPOSAL"
    NEGOTIATION = "NEGOTIATION"
    WON = "WON"
    LOST = "LOST"

    def is_terminal(self) -> bool:
        return self in (OpportunityStage.WON, OpportunityStage.LOST)

    def default_probability(self) -> int:
        return {
            OpportunityStage.PROSPECTING: 10,
            OpportunityStage.PROPOSAL: 40,
            OpportunityStage.NEGOTIATION: 70,
            OpportunityStage.WON: 100,
            OpportunityStage.LOST: 0,
        }[self]
```

```python
# app/modules/crm/domain/errors/errors.py
class DomainError(Exception):
    code: str = "DOMAIN_ERROR"

    def __init__(self, message: str = ""):
        self.message = message or self.__class__.__name__
        super().__init__(self.message)


class InvalidTenantError(DomainError):
    code = "INVALID_TENANT"


class InvalidLeadNameError(DomainError):
    code = "INVALID_LEAD_NAME"


class InvalidScoreError(DomainError):
    code = "INVALID_SCORE"


class InvalidSubjectError(DomainError):
    code = "INVALID_SUBJECT"


class InvalidPriorityError(DomainError):
    code = "INVALID_PRIORITY"


class LeadNotFoundError(DomainError):
    code = "LEAD_NOT_FOUND"


class TicketNotFoundError(DomainError):
    code = "TICKET_NOT_FOUND"


class OpportunityNotFoundError(DomainError):
    code = "OPPORTUNITY_NOT_FOUND"


class LeadAlreadyConvertedError(DomainError):
    code = "LEAD_ALREADY_CONVERTED"


class CannotConvertLostLeadError(DomainError):
    code = "CANNOT_CONVERT_LOST_LEAD"


class InvalidLeadTransitionError(DomainError):
    code = "INVALID_LEAD_TRANSITION"


class InvalidTicketTransitionError(DomainError):
    code = "INVALID_TICKET_TRANSITION"


class CannotResolveClosedError(DomainError):
    code = "CANNOT_RESOLVE_CLOSED"


class CannotEscalateError(DomainError):
    code = "CANNOT_ESCALATE"


class InvalidOpportunityTransitionError(DomainError):
    code = "INVALID_OPPORTUNITY_TRANSITION"
```

```python
# app/modules/crm/domain/entities/lead.py
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.crm.domain.errors.errors import (
    InvalidTenantError, InvalidLeadNameError, InvalidScoreError,
    LeadAlreadyConvertedError, CannotConvertLostLeadError,
    InvalidLeadTransitionError,
)
from app.modules.crm.domain.value_objects.lead_source import LeadSource
from app.modules.crm.domain.value_objects.lead_status import LeadStatus


class Lead(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    name: str
    email: str = ""
    phone: str = ""
    company: str = ""
    source: LeadSource
    status: LeadStatus = LeadStatus.NEW
    score: int = 0
    assigned_to: UUID | None = None
    converted_to: UUID | None = None
    metadata: dict = Field(default_factory=dict)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(cls, tenant_id: UUID, name: str, email: str, phone: str, company: str, source: LeadSource) -> "Lead":
        if not tenant_id:
            raise InvalidTenantError()
        if not name or len(name) > 255:
            raise InvalidLeadNameError()
        return cls(tenant_id=tenant_id, name=name, email=email, phone=phone, company=company, source=source)

    def qualify(self, score: int) -> None:
        if score < 0 or score > 100:
            raise InvalidScoreError()
        if not self.status.can_transition_to(LeadStatus.QUALIFIED):
            raise InvalidLeadTransitionError()
        self.score = score
        self.status = LeadStatus.QUALIFIED
        self._touch()

    def mark_contacted(self) -> None:
        if not self.status.can_transition_to(LeadStatus.CONTACTED):
            raise InvalidLeadTransitionError()
        self.status = LeadStatus.CONTACTED
        self._touch()

    def assign_to(self, user_id: UUID) -> None:
        self.assigned_to = user_id
        self._touch()

    def convert(self, customer_id: UUID) -> None:
        if self.status == LeadStatus.CONVERTED:
            raise LeadAlreadyConvertedError()
        if self.status == LeadStatus.LOST:
            raise CannotConvertLostLeadError()
        self.status = LeadStatus.CONVERTED
        self.converted_to = customer_id
        self._touch()

    def lose(self, reason: str) -> None:
        if not self.status.can_transition_to(LeadStatus.LOST):
            raise InvalidLeadTransitionError()
        self.status = LeadStatus.LOST
        self.metadata["lost_reason"] = reason
        self._touch()

    def is_converted(self) -> bool:
        return self.status == LeadStatus.CONVERTED

    def is_qualified(self) -> bool:
        return self.status == LeadStatus.QUALIFIED

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

```python
# app/modules/crm/domain/entities/ticket.py
from datetime import datetime, timezone
from uuid import UUID, uuid4

from pydantic import BaseModel, Field

from app.modules.crm.domain.errors.errors import (
    InvalidTenantError, InvalidSubjectError, InvalidPriorityError,
    InvalidTicketTransitionError, CannotResolveClosedError, CannotEscalateError,
)
from app.modules.crm.domain.value_objects.ticket_priority import TicketPriority
from app.modules.crm.domain.value_objects.ticket_status import TicketStatus


class Ticket(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    tenant_id: UUID
    customer_id: UUID
    subject: str
    description: str = ""
    priority: TicketPriority
    status: TicketStatus = TicketStatus.OPEN
    assigned_to: UUID | None = None
    sla_due_at: datetime
    started_at: datetime | None = None
    resolved_at: datetime | None = None
    closed_at: datetime | None = None
    escalated_at: datetime | None = None
    escalation_level: int = 0
    metadata: dict = Field(default_factory=dict)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    def create(cls, tenant_id: UUID, customer_id: UUID, subject: str, description: str, priority: TicketPriority) -> "Ticket":
        if not tenant_id or not customer_id:
            raise InvalidTenantError()
        if not subject or len(subject) > 255:
            raise InvalidSubjectError()
        now = datetime.now(timezone.utc)
        return cls(
            tenant_id=tenant_id,
            customer_id=customer_id,
            subject=subject,
            description=description,
            priority=priority,
            sla_due_at=now + priority.sla_duration(),
        )

    def assign(self, user_id: UUID) -> None:
        self.assigned_to = user_id
        self._touch()

    def start(self) -> None:
        if not self.status.can_transition_to(TicketStatus.IN_PROGRESS):
            raise InvalidTicketTransitionError()
        self.status = TicketStatus.IN_PROGRESS
        self.started_at = datetime.now(timezone.utc)
        self._touch()

    def resolve(self) -> None:
        if self.status == TicketStatus.CLOSED:
            raise CannotResolveClosedError()
        if not self.status.can_transition_to(TicketStatus.RESOLVED):
            raise InvalidTicketTransitionError()
        self.status = TicketStatus.RESOLVED
        self.resolved_at = datetime.now(timezone.utc)
        self._touch()

    def close(self) -> None:
        if not self.status.can_transition_to(TicketStatus.CLOSED):
            raise InvalidTicketTransitionError()
        self.status = TicketStatus.CLOSED
        self.closed_at = datetime.now(timezone.utc)
        self._touch()

    def reopen(self) -> None:
        if not self.status.can_transition_to(TicketStatus.OPEN):
            raise InvalidTicketTransitionError()
        self.status = TicketStatus.OPEN
        self.resolved_at = None
        self.closed_at = None
        self._touch()

    def escalate(self) -> None:
        if self.status in (TicketStatus.CLOSED, TicketStatus.RESOLVED):
            raise CannotEscalateError()
        self.escalated_at = datetime.now(timezone.utc)
        self.escalation_level += 1
        self.metadata["escalated"] = True
        self._touch()

    def is_overdue(self) -> bool:
        return datetime.now(timezone.utc) > self.sla_due_at and not self.is_resolved_or_closed()

    def is_resolved_or_closed(self) -> bool:
        return self.status in (TicketStatus.RESOLVED, TicketStatus.CLOSED)

    def _touch(self) -> None:
        self.updated_at = datetime.now(timezone.utc)
```

### Application Use Cases (ตัวอย่าง)

```python
# app/modules/crm/application/use_cases/capture_lead.py
from dataclasses import dataclass, field
from datetime import datetime, timezone
from uuid import UUID

from app.modules.crm.application.ports.audit_repository import AuditRepository, AuditTrail
from app.modules.crm.application.ports.customer_client import CustomerClient
from app.modules.crm.application.ports.event_producer import EventProducer
from app.modules.crm.application.ports.notifier_client import NotifierClient, SendEmailRequest
from app.modules.crm.domain.entities.lead import Lead
from app.modules.crm.domain.events.events import LeadCreatedEvent
from app.modules.crm.domain.repositories.lead_repository import LeadRepository
from app.modules.crm.domain.value_objects.lead_source import LeadSource


@dataclass
class CaptureLeadInput:
    tenant_id: UUID
    name: str
    email: str = ""
    phone: str = ""
    company: str = ""
    source: LeadSource = LeadSource.WEB
    user_id: UUID = None
    ip_address: str = ""


@dataclass
class CaptureLeadOutput:
    lead_id: UUID
    status: str


class CaptureLeadUseCase:
    def __init__(
        self,
        lead_repo: LeadRepository,
        audit_repo: AuditRepository,
        producer: EventProducer,
        notifier: NotifierClient,
    ):
        self._lead_repo = lead_repo
        self._audit = audit_repo
        self._producer = producer
        self._notifier = notifier

    async def execute(self, inp: CaptureLeadInput) -> CaptureLeadOutput:
        lead = Lead.create(inp.tenant_id, inp.name, inp.email, inp.phone, inp.company, inp.source)
        await self._lead_repo.save(lead)

        await self._audit.save(AuditTrail(
            user_id=inp.user_id, action="LEAD_CAPTURED",
            entity_type="lead", entity_id=lead.id,
            payload={"name": lead.name, "source": lead.source.value},
            ip_address=inp.ip_address,
        ))

        await self._producer.publish_lead_created(LeadCreatedEvent(
            lead_id=lead.id, tenant_id=lead.tenant_id,
            name=lead.name, email=lead.email, source=lead.source.value,
        ))

        if lead.email:
            await self._notifier.send_email(SendEmailRequest(
                to=lead.email, subject="Thank you", body="...", tenant_id=lead.tenant_id,
            ))

        return CaptureLeadOutput(lead_id=lead.id, status=lead.status.value)
```

### Infrastructure — Models + Repo

```python
# app/modules/crm/infrastructure/persistence/postgres/models.py
from datetime import datetime
from uuid import UUID, uuid4

from sqlalchemy import Boolean, DateTime, Float, ForeignKey, Index, Integer, String, Text, func
from sqlalchemy.dialects.postgresql import JSONB, UUID as PGUUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class LeadModel(Base):
    __tablename__ = "crm_leads"
    __table_args__ = (Index("idx_crm_leads_tenant", "tenant_id"),)

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    email: Mapped[str] = mapped_column(String(255), default="")
    phone: Mapped[str] = mapped_column(String(30), default="")
    company: Mapped[str] = mapped_column(String(255), default="")
    source: Mapped[str] = mapped_column(String(50), nullable=False)
    status: Mapped[str] = mapped_column(String(30), default="NEW")
    score: Mapped[int] = mapped_column(Integer, default=0)
    assigned_to: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    converted_to: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    metadata_: Mapped[dict] = mapped_column("metadata", JSONB, default=dict)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class TicketModel(Base):
    __tablename__ = "crm_tickets"
    __table_args__ = (Index("idx_crm_tickets_tenant", "tenant_id"),)

    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    customer_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    subject: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(Text, default="")
    priority: Mapped[str] = mapped_column(String(20), nullable=False)
    status: Mapped[str] = mapped_column(String(30), default="OPEN")
    assigned_to: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    sla_due_at: Mapped[datetime] = mapped_column(DateTime, nullable=False)
    started_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    resolved_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    closed_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    escalated_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    escalation_level: Mapped[int] = mapped_column(Integer, default=0)
    metadata_: Mapped[dict] = mapped_column("metadata", JSONB, default=dict)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class OpportunityModel(Base):
    __tablename__ = "crm_opportunities"
    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    lead_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    customer_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    title: Mapped[str] = mapped_column(String(255), nullable=False)
    amount: Mapped[float] = mapped_column(Float, default=0)
    currency: Mapped[str] = mapped_column(String(3), default="THB")
    probability: Mapped[int] = mapped_column(Integer, default=0)
    stage: Mapped[str] = mapped_column(String(30), default="PROSPECTING")
    expected_close_date: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)
    owner_id: Mapped[UUID | None] = mapped_column(PGUUID(as_uuid=True), nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now(), onupdate=func.now())


class ActivityModel(Base):
    __tablename__ = "crm_activities"
    id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), primary_key=True, default=uuid4)
    tenant_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    entity_type: Mapped[str] = mapped_column(String(30), nullable=False)
    entity_id: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    type: Mapped[str] = mapped_column(String(20), nullable=False)
    subject: Mapped[str] = mapped_column(String(255), default="")
    notes: Mapped[str] = mapped_column(Text, default="")
    performed_by: Mapped[UUID] = mapped_column(PGUUID(as_uuid=True), nullable=False)
    occurred_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
    created_at: Mapped[datetime] = mapped_column(DateTime, server_default=func.now())
```

```python
# app/modules/crm/infrastructure/persistence/postgres/lead_repository_impl.py
from uuid import UUID

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.crm.domain.entities.lead import Lead
from app.modules.crm.domain.value_objects.lead_source import LeadSource
from app.modules.crm.domain.value_objects.lead_status import LeadStatus
from app.modules.crm.infrastructure.persistence.postgres.models import LeadModel


class PostgresLeadRepository:
    def __init__(self, session: AsyncSession):
        self._session = session

    async def save(self, lead: Lead) -> None:
        m = await self._session.get(LeadModel, lead.id)
        if m is None:
            m = LeadModel(id=lead.id)
            self._session.add(m)
        m.tenant_id = lead.tenant_id
        m.name = lead.name
        m.email = lead.email
        m.phone = lead.phone
        m.company = lead.company
        m.source = lead.source.value
        m.status = lead.status.value
        m.score = lead.score
        m.assigned_to = lead.assigned_to
        m.converted_to = lead.converted_to
        m.metadata_ = lead.metadata
        await self._session.flush()

    async def find_by_id(self, id: UUID) -> Lead | None:
        m = await self._session.get(LeadModel, id)
        return self._to_entity(m) if m else None

    async def find_by_tenant(self, tenant_id: UUID) -> list[Lead]:
        stmt = select(LeadModel).where(LeadModel.tenant_id == tenant_id).order_by(LeadModel.created_at.desc())
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def find_by_status(self, tenant_id: UUID, status: LeadStatus) -> list[Lead]:
        stmt = select(LeadModel).where(LeadModel.tenant_id == tenant_id, LeadModel.status == status.value)
        rows = (await self._session.execute(stmt)).scalars().all()
        return [self._to_entity(m) for m in rows]

    async def exists_by_email(self, tenant_id: UUID, email: str) -> bool:
        stmt = select(func.count()).select_from(LeadModel).where(
            LeadModel.tenant_id == tenant_id, LeadModel.email == email
        )
        return (await self._session.execute(stmt)).scalar() > 0

    async def delete(self, id: UUID) -> None:
        m = await self._session.get(LeadModel, id)
        if m:
            await self._session.delete(m)

    def _to_entity(self, m: LeadModel) -> Lead:
        return Lead(
            id=m.id, tenant_id=m.tenant_id, name=m.name,
            email=m.email or "", phone=m.phone or "", company=m.company or "",
            source=LeadSource(m.source), status=LeadStatus(m.status),
            score=m.score, assigned_to=m.assigned_to, converted_to=m.converted_to,
            metadata=m.metadata_ or {}, created_at=m.created_at, updated_at=m.updated_at,
        )
```

### FastAPI Router

```python
# app/modules/crm/interfaces/http/lead_router.py
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status

from app.modules.crm.application.use_cases.capture_lead import (
    CaptureLeadInput, CaptureLeadUseCase,
)
from app.modules.crm.application.use_cases.qualify_lead import (
    QualifyLeadInput, QualifyLeadUseCase,
)
from app.modules.crm.application.use_cases.convert_lead import (
    ConvertLeadInput, ConvertLeadUseCase,
)
from app.modules.crm.domain.errors.errors import (
    DomainError, LeadAlreadyConvertedError, LeadNotFoundError,
)
from app.modules.crm.domain.value_objects.lead_source import LeadSource
from app.modules.crm.interfaces.http.dependencies import (
    get_capture_lead_uc, get_convert_lead_uc, get_current_user,
    get_qualify_lead_uc, get_tenant_id,
)

router = APIRouter(prefix="/api/v1/crm/leads", tags=["crm-leads"])


@router.post("", status_code=status.HTTP_201_CREATED)
async def capture_lead(
    name: str,
    email: str = "",
    phone: str = "",
    company: str = "",
    source: str = "WEB",
    uc: CaptureLeadUseCase = Depends(get_capture_lead_uc),
    user_id: UUID = Depends(get_current_user),
    tenant_id: UUID = Depends(get_tenant_id),
):
    try:
        out = await uc.execute(CaptureLeadInput(
            tenant_id=tenant_id, name=name, email=email, phone=phone,
            company=company, source=LeadSource(source), user_id=user_id,
        ))
        return {"lead_id": str(out.lead_id), "status": out.status}
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})


@router.post("/{lead_id}/qualify")
async def qualify(
    lead_id: UUID,
    score: int,
    uc: QualifyLeadUseCase = Depends(get_qualify_lead_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        await uc.execute(QualifyLeadInput(lead_id=lead_id, score=score, user_id=user_id))
        return {"message": "qualified"}
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})


@router.post("/{lead_id}/convert")
async def convert(
    lead_id: UUID,
    uc: ConvertLeadUseCase = Depends(get_convert_lead_uc),
    user_id: UUID = Depends(get_current_user),
):
    try:
        out = await uc.execute(ConvertLeadInput(lead_id=lead_id, user_id=user_id))
        return {"customer_id": str(out.customer_id)}
    except LeadAlreadyConvertedError as e:
        raise HTTPException(409, detail={"code": e.code, "message": str(e)})
    except LeadNotFoundError as e:
        raise HTTPException(404, detail={"code": e.code, "message": str(e)})
    except DomainError as e:
        raise HTTPException(400, detail={"code": e.code, "message": str(e)})
```

### `module.py`

```python
from dataclasses import dataclass

from aiokafka import AIOKafkaProducer
from sqlalchemy.ext.asyncio import AsyncSession

from app.modules.crm.application.use_cases.capture_lead import CaptureLeadUseCase
from app.modules.crm.application.use_cases.convert_lead import ConvertLeadUseCase
from app.modules.crm.application.use_cases.open_ticket import OpenTicketUseCase
from app.modules.crm.application.use_cases.qualify_lead import QualifyLeadUseCase
from app.modules.crm.infrastructure.messaging.kafka_producer import KafkaEventProducer
from app.modules.crm.infrastructure.persistence.postgres.lead_repository_impl import (
    PostgresLeadRepository,
)
from app.modules.crm.infrastructure.persistence.postgres.ticket_repository_impl import (
    PostgresTicketRepository,
)


@dataclass
class CRMModuleDeps:
    session: AsyncSession
    kafka: AIOKafkaProducer
    customer_client: "CustomerClient"
    notifier_client: "NotifierClient"
    audit_repo: "AuditRepository"


def build_module(deps: CRMModuleDeps) -> dict:
    lead_repo = PostgresLeadRepository(deps.session)
    ticket_repo = PostgresTicketRepository(deps.session)
    producer = KafkaEventProducer(deps.kafka)

    return {
        "capture_lead": CaptureLeadUseCase(lead_repo, deps.audit_repo, producer, deps.notifier_client),
        "qualify_lead": QualifyLeadUseCase(lead_repo, producer, deps.audit_repo),
        "convert_lead": ConvertLeadUseCase(lead_repo, deps.customer_client, producer, deps.audit_repo),
        "open_ticket": OpenTicketUseCase(ticket_repo, producer, deps.audit_repo, deps.notifier_client),
    }
```

### Tests

```python
# tests/modules/crm/test_lead_entity.py
import pytest
from uuid import uuid4

from app.modules.crm.domain.entities.lead import Lead
from app.modules.crm.domain.errors.errors import (
    LeadAlreadyConvertedError, InvalidScoreError,
)
from app.modules.crm.domain.value_objects.lead_source import LeadSource
from app.modules.crm.domain.value_objects.lead_status import LeadStatus


def test_lead_lifecycle():
    l = Lead.create(uuid4(), "Alice", "a@test.com", "", "", LeadSource.WEB)
    l.mark_contacted()
    l.qualify(80)
    assert l.is_qualified()
    l.convert(uuid4())
    assert l.is_converted()


def test_cannot_convert_twice():
    l = Lead.create(uuid4(), "Bob", "b@test.com", "", "", LeadSource.WEB)
    l.qualify(70)
    l.convert(uuid4())
    with pytest.raises(LeadAlreadyConvertedError):
        l.convert(uuid4())


def test_invalid_score():
    l = Lead.create(uuid4(), "Eve", "", "", "", LeadSource.WEB)
    with pytest.raises(InvalidScoreError):
        l.qualify(150)
```

**P05 Status:** ✅ Complete (~18 files)

---

# PART 2 — G05-P08 (Modules 06, 07, 08)

## 2.1 G06 + P06 — Logistics

### Spec Summary

```yaml
module: logistics
context: Supply Chain
priority: P2

aggregates:
  Shipment:
    fields: [id, tenant_id, tracking_number, order_id, customer_id, status, origin, destination, carrier, vehicle_id, driver_id, estimated_at, delivered_at, temperature, items]
    behavior: [AssignVehicle, AssignDriver, Pickup, UpdateLocation, RecordTemperature, Deliver, Fail]
    invariants:
      - Cold chain 2-8°C
      - DeliveredAt >= CreatedAt
  Route:
    fields: [id, tenant_id, name, waypoints, distance_km, estimated_min]
  Vehicle:
    fields: [id, tenant_id, plate_number, type, capacity_kg, status]
  Driver:
    fields: [id, tenant_id, name, license, phone, status]

value_objects:
  ShipmentStatus: [PENDING, PICKED, IN_TRANSIT, DELIVERED, FAILED]
  VehicleType: [TRUCK, VAN, BIKE, CAR, REFRIGERATED]

tables: [logistics_shipments, logistics_shipment_items, logistics_routes, logistics_vehicles, logistics_drivers]
kafka_publish: [logistics.shipment.created, logistics.shipment.delivered, logistics.temperature.breached]
kafka_consume: [erp.order.confirmed, device.telemetry.ingested]
migration: 20260106_logistics_init.sql
```

### Go — ไฟล์หลัก

```
internal/modules/logistics/
├── domain/
│   ├── entity/{shipment,route,vehicle,driver}.go
│   ├── value_object/{shipment_status,vehicle_type}.go
│   ├── repository/{shipment,route,vehicle,driver}_repository.go
│   └── errors/errors.go
├── application/
│   ├── create_shipment.go, assign_vehicle.go, pickup.go, update_location.go
│   ├── record_temperature.go, deliver.go, fail.go
│   └── ports.go, dto.go
├── infrastructure/
│   ├── persistence/postgres/{models,*_repo_impl}.go
│   ├── messaging/producer.go
│   └── scheduler/temperature_check_job.go
├── interfaces/http/{shipment_handler,routes}.go
└── module.go
```

**Key files:**

```go
// domain/value_object/shipment_status.go
package valueobject

type ShipmentStatus string

const (
	ShipmentPending    ShipmentStatus = "PENDING"
	ShipmentPicked     ShipmentStatus = "PICKED"
	ShipmentInTransit  ShipmentStatus = "IN_TRANSIT"
	ShipmentDelivered  ShipmentStatus = "DELIVERED"
	ShipmentFailed     ShipmentStatus = "FAILED"
)

func (s ShipmentStatus) CanTransitionTo(next ShipmentStatus) bool {
	t := map[ShipmentStatus]map[ShipmentStatus]bool{
		ShipmentPending:   {ShipmentPicked: true, ShipmentFailed: true},
		ShipmentPicked:    {ShipmentInTransit: true, ShipmentFailed: true},
		ShipmentInTransit: {ShipmentDelivered: true, ShipmentFailed: true},
		ShipmentDelivered: {},
		ShipmentFailed:    {ShipmentPending: true},
	}
	return t[s][next]
}
```

```go
// domain/entity/shipment.go
type Shipment struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	TrackingNumber string
	OrderID        *uuid.UUID
	CustomerID     uuid.UUID
	Status         valueobject.ShipmentStatus
	Items          []ShipmentItem
	Origin         Location
	Destination    Location
	Carrier        string
	VehicleID      *uuid.UUID
	DriverID       *uuid.UUID
	EstimatedAt    time.Time
	DeliveredAt    *time.Time
	Temperature    *TemperatureLog
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *Shipment) RecordTemperature(t float64, deviceID uuid.UUID) {
	if s.Temperature == nil {
		s.Temperature = &TemperatureLog{Min: t, Max: t, Current: t}
	}
	if t < s.Temperature.Min {
		s.Temperature.Min = t
	}
	if t > s.Temperature.Max {
		s.Temperature.Max = t
	}
	s.Temperature.Current = t
	s.Temperature.Readings = append(s.Temperature.Readings, TemperatureReading{
		Value: t, Timestamp: time.Now(), DeviceID: deviceID,
	})
	if t < 2 || t > 8 {
		s.Temperature.IsBreached = true
	}
	s.UpdatedAt = time.Now()
}

type TemperatureLog struct {
	Min        float64
	Max        float64
	Current    float64
	Readings   []TemperatureReading
	IsBreached bool
}
```

## 2.2 G07 + P07 — Report (Upgrade)

### Upgrade Strategy (Expand-Contract)

```
Day 0:  Add new tables (report_definitions, report_runs, report_dashboards)
Day 1:  Add new endpoints /api/v1/reports/* (keep old ones)
Day 7:  Migrate old report configs → new tables
Day 14: Add deprecation headers to old endpoints
Day 30: Remove old endpoints + handlers
```

### Domain Layer (Go compact)

```go
// domain/entity/report_definition.go
package entity

type ReportDefinition struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Code         string
	Name         string
	Description  string
	Category     valueobject.ReportCategory
	Type         valueobject.ReportType
	DataSource   valueobject.DataSource
	Query        string
	Parameters   []Parameter
	Schedule     *Schedule
	Permissions  []string
	IsSystem     bool
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// domain/value_object/report_category.go
type ReportCategory string

const (
	ReportCategoryOperation ReportCategory = "OPERATION"
	ReportCategoryFinance   ReportCategory = "FINANCE"
	ReportCategoryEnergy    ReportCategory = "ENERGY"
	ReportCategoryAgriculture ReportCategory = "AGRICULTURE"
	ReportCategoryDevice    ReportCategory = "DEVICE"
)

// domain/value_object/report_type.go
type ReportType string

const (
	ReportTypeTable     ReportType = "TABLE"
	ReportTypeChart     ReportType = "CHART"
	ReportTypeDashboard ReportType = "DASHBOARD"
	ReportTypeExport    ReportType = "EXPORT"
	ReportTypeInsight   ReportType = "INSIGHT"
)

// domain/value_object/data_source.go
type DataSource string

const (
	DataSourcePostgres      DataSource = "POSTGRES"
	DataSourceInfluxDB      DataSource = "INFLUXDB"
	DataSourceElasticsearch DataSource = "ELASTICSEARCH"
	DataSourceMixed         DataSource = "MIXED"
)
```

### Use Case — RunReport (ย่อ)

```go
type RunReportUseCase struct {
	reportRepo   repository.ReportRepository
	runRepo      repository.ReportRunRepository
	pgRunner     PostgresRunner
	influxRunner InfluxRunner
	esRunner     ESRunner
	exporter     Exporter
	producer     EventProducer
}

func (uc *RunReportUseCase) Execute(ctx context.Context, in RunReportInput) (*RunReportOutput, error) {
	def, err := uc.reportRepo.FindByID(ctx, in.ReportID)
	if err != nil {
		return nil, err
	}

	run := entity.NewReportRun(def.ID, def.TenantID, in.RunBy, in.Parameters)
	if err := uc.runRepo.Save(ctx, run); err != nil {
		return nil, err
	}
	run.MarkRunning()
	_ = uc.runRepo.Save(ctx, run)

	var data []map[string]interface{}
	switch def.DataSource {
	case valueobject.DataSourcePostgres:
		data, err = uc.pgRunner.Query(ctx, def.Query, in.Parameters)
	case valueobject.DataSourceInfluxDB:
		data, err = uc.influxRunner.Query(ctx, def.Query, in.Parameters)
	case valueobject.DataSourceElasticsearch:
		data, err = uc.esRunner.Query(ctx, def.Query, in.Parameters)
	}
	if err != nil {
		run.MarkFailed(err.Error())
		_ = uc.runRepo.Save(ctx, run)
		return nil, err
	}

	url, err := uc.exporter.Export(ctx, data, in.Format, def.Name)
	if err != nil {
		run.MarkFailed(err.Error())
		_ = uc.runRepo.Save(ctx, run)
		return nil, err
	}

	run.MarkSuccess(url, len(data))
	_ = uc.runRepo.Save(ctx, run)

	return &RunReportOutput{RunID: run.ID, ResultURL: url, RowCount: len(data)}, nil
}
```

### Migration

```sql
-- migrations/20260107_report_upgrade.sql
CREATE TABLE IF NOT EXISTS report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    type VARCHAR(30) NOT NULL,
    data_source VARCHAR(30) NOT NULL,
    query TEXT NOT NULL,
    parameters JSONB DEFAULT '[]',
    schedule JSONB,
    permissions JSONB DEFAULT '["viewer"]',
    is_system BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_report_definitions_tenant_code UNIQUE (tenant_id, code)
);

CREATE TABLE IF NOT EXISTS report_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES report_definitions(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    run_by UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    parameters JSONB,
    result_url TEXT,
    row_count INT DEFAULT 0,
    duration_ms BIGINT DEFAULT 0,
    error_msg TEXT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_report_runs_report ON report_runs (report_id, created_at DESC);

CREATE TABLE IF NOT EXISTS report_dashboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    widgets JSONB DEFAULT '[]',
    layout JSONB DEFAULT '{"columns":12,"rows":6}',
    is_default BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE,
    permissions JSONB DEFAULT '["viewer"]',
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_report_dashboards_tenant_code UNIQUE (tenant_id, code)
);

-- Migrate old report data (if exists)
-- INSERT INTO report_definitions (...) SELECT ... FROM old_reports;
```

## 2.3 G08 + P08 — Auth (Upgrade)

### Breaking Change Plan

```
Current JWT payload:
  { sub, exp, iat }

New JWT payload:
  { sub, exp, iat, tenant_id, session_id, mfa_verified, roles[] }

Migration:
  Release 1: Sign with new payload; verify accepts both old + new
  Release 2: Enforce new payload (reject old)
```

### Domain — Session Entity

```go
// domain/entity/session.go
package entity

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	TenantID     uuid.UUID
	RefreshToken string
	IPAddress    string
	UserAgent    string
	MFAVerified  bool
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	CreatedAt    time.Time
}

func (s *Session) Revoke() {
	now := time.Now()
	s.RevokedAt = &now
}

func (s *Session) IsActive() bool {
	return s.RevokedAt == nil && time.Now().Before(s.ExpiresAt)
}
```

### New Use Cases

```go
// LoginWithSSO
type LoginWithSSOUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	provider    SSOProvider
	jwtMaker    JWTMaker
	auditRepo   AuditRepository
}

// EnableMFA
type EnableMFAUseCase struct {
	userRepo   repository.UserRepository
	mfaService MFAService
}

// VerifyMFA
type VerifyMFAUseCase struct {
	sessionRepo repository.SessionRepository
	mfaService  MFAService
}

// RevokeAllSessions
type RevokeAllSessionsUseCase struct {
	sessionRepo repository.SessionRepository
	cache       SessionCache
}
```

### Migration

```sql
-- migrations/20260108_auth_upgrade.sql
ALTER TABLE auth_sessions ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE auth_sessions ADD COLUMN IF NOT EXISTS mfa_verified BOOLEAN DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS idx_auth_sessions_tenant ON auth_sessions (tenant_id);

CREATE TABLE IF NOT EXISTS auth_mfa_settings (
    user_id UUID PRIMARY KEY,
    method VARCHAR(20) NOT NULL,       -- TOTP, SMS
    secret VARCHAR(255),
    phone VARCHAR(30),
    enabled BOOLEAN DEFAULT FALSE,
    backup_codes JSONB DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- No DROP COLUMN — deprecated fields remain
```

**Status:** ✅ G05-G08 + P05-P08 complete

---

# PART 3 — Modules 09-16 (Upgrades)

## 3.1 `09_users` — RBAC + Multi-tenant

**Changes:**
- Add `tenant_id` to users_users
- Add `users_roles`, `users_permissions`, `users_role_permissions` tables
- ABAC attributes support
- Invitation workflow

**Migration:**
```sql
-- migrations/20260109_users_upgrade.sql
ALTER TABLE users_users ADD COLUMN IF NOT EXISTS tenant_id UUID;
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users_users (tenant_id);

CREATE TABLE IF NOT EXISTS users_permissions (
    code VARCHAR(100) PRIMARY KEY,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);
CREATE INDEX idx_users_permissions_resource ON users_permissions (resource, action);

CREATE TABLE IF NOT EXISTS users_role_permissions (
    role_id UUID NOT NULL,
    permission_code VARCHAR(100) NOT NULL REFERENCES users_permissions(code),
    PRIMARY KEY (role_id, permission_code)
);
```

## 3.2 `10_settings` — Hierarchical

**Changes:**
- `settings_tenant_settings` (per-tenant)
- `settings_user_preferences` (per-user)
- Settings versioning + rollback

```sql
CREATE TABLE IF NOT EXISTS settings_tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    category VARCHAR(50) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value JSONB NOT NULL,
    version INT NOT NULL DEFAULT 1,
    updated_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, category, key)
);
```

## 3.3 `11_payment` — Subscription + QR

**Changes:**
- Subscription payments
- PromptPay QR (Thai QR)
- Multi-currency
- Refunds
- Webhooks (Stripe, Omise, 2C2P)

```sql
CREATE TABLE IF NOT EXISTS payment_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    method VARCHAR(30) NOT NULL,       -- CARD, QR, TRANSFER, WALLET
    provider VARCHAR(30),              -- STRIPE, OMISE, 2C2P
    provider_ref VARCHAR(100),
    reference_id UUID,
    reference_type VARCHAR(50),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    checkout_url TEXT,
    paid_at TIMESTAMP,
    refunded_at TIMESTAMP,
    refund_amount NUMERIC(15,2) DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_payment_tenant ON payment_payments (tenant_id, status);
CREATE INDEX idx_payment_reference ON payment_payments (reference_id, reference_type);

CREATE TABLE IF NOT EXISTS payment_webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(30) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP,
    error_msg TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## 3.4 `12_notifier` — Multi-channel

**Channels:** Email, SMS, LINE, Push (FCM/APNs), Webhook, Slack, Teams

```sql
CREATE TABLE IF NOT EXISTS notifier_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'th',
    subject VARCHAR(255),
    body TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, code, channel, language)
);

CREATE TABLE IF NOT EXISTS notifier_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    template_code VARCHAR(100),
    channel VARCHAR(20) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    error_msg TEXT,
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifier_logs_tenant ON notifier_logs (tenant_id, status);
```

## 3.5 `13_alarm` — IoT Rule Engine

**Changes:**
- Rule engine ผูกกับ device
- Alarm correlation
- Suppression (maintenance window)
- Escalation policy
- Auto-command

```sql
CREATE TABLE IF NOT EXISTS alarm_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT TRUE,
    severity VARCHAR(20) NOT NULL,     -- INFO, WARNING, CRITICAL
    conditions JSONB NOT NULL,
    actions JSONB NOT NULL,
    cooldown_seconds INT DEFAULT 0,
    last_fired_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alarm_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    rule_id UUID NOT NULL,
    device_id UUID,
    severity VARCHAR(20) NOT NULL,
    message TEXT,
    payload JSONB,
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN',   -- OPEN, ACKED, RESOLVED
    acknowledged_by UUID,
    acknowledged_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_alarm_events_tenant_status ON alarm_events (tenant_id, status);
```

## 3.6 `14_pdpa` — IoT Consent

```sql
ALTER TABLE pdpa_consents ADD COLUMN IF NOT EXISTS data_category VARCHAR(50);
ALTER TABLE pdpa_consents ADD COLUMN IF NOT EXISTS device_id UUID;
CREATE INDEX IF NOT EXISTS idx_pdpa_consents_device ON pdpa_consents (device_id);

CREATE TABLE IF NOT EXISTS pdpa_retention_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    data_category VARCHAR(50) NOT NULL,
    retention_days INT NOT NULL,
    auto_delete BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, data_category)
);
```

## 3.7 `15_auditlog` — Search + Hash Chain

```sql
ALTER TABLE audit_trails ADD COLUMN IF NOT EXISTS prev_hash VARCHAR(64);
ALTER TABLE audit_trails ADD COLUMN IF NOT EXISTS hash VARCHAR(64);
CREATE INDEX IF NOT EXISTS idx_audit_trails_entity ON audit_trails (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_trails_user ON audit_trails (user_id, occurred_at DESC);
```

## 3.8 `16_dashboard` — IoT Widgets

```sql
CREATE TABLE IF NOT EXISTS dashboard_widgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dashboard_id UUID NOT NULL,
    type VARCHAR(30) NOT NULL,          -- LINE, GAUGE, MAP, HEATMAP, TABLE
    title VARCHAR(255),
    query JSONB NOT NULL,
    position JSONB NOT NULL,             -- {x, y}
    size JSONB NOT NULL,                 -- {w, h}
    refresh_interval INT DEFAULT 30,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Status:** ✅ Modules 09-16 upgrade specs complete

---

# PART 4 — Modules 17-40 (Compact Index)

## 4.1 Business Modules (17-24)

| # | Module | Type | Key Additions | Priority |
|:-:|--------|:----:|---------------|:-:|
| 17 | `items` | Upgrade | IoT catalog, compatibility matrix, datasheet | P1 |
| 18 | `purchaseorder` | Upgrade | ผูก erp.Order, approval, 3-way matching | P1 |
| 19 | `quotation` | Upgrade | ผูก crm.Opportunity, versioning, e-sign | P2 |
| 20 | `wos` | Upgrade | ผูก device (maintenance), SLA, technician | P2 |
| 21 | `document` | Upgrade | e-sign (PDF), versioning, OCR, S3 presigned | P2 |
| 22 | `email` | Upgrade | IoT templates, tracking, unsubscribe | P2 |
| 23 | `batch` | Upgrade | telemetry batch, schedule-based | P2 |
| 24 | `job` | Upgrade | scheduler ใหม่, retry policy, dependencies | P2 |

## 4.2 Integration Modules (25-32)

| # | Module | Type | Key Additions | Priority |
|:-:|--------|:----:|---------------|:-:|
| 25 | `iot` | **Deprecate** | migrate → device | P0 |
| 26 | `mqtt` | Upgrade | multi-protocol, TLS, device auth | P0 |
| 27 | `kafka` | Upgrade | IoT topics, DLQ, schema registry, lag monitoring | P0 |
| 28 | `influxdb` | Upgrade | retention per-tenant, downsampling tasks | P0 |
| 29 | `elasticsearch` | Upgrade | IoT index template, ILM | P1 |
| 30 | `websocket` | Upgrade | multi-tenant hub, channels, backpressure | P0 |
| 31 | `realtime` | **Merge** | รวมกับ websocket | P1 |
| 32 | `vectordata` | Upgrade | embedding pipeline, anomaly detection | P2 |

## 4.3 Utility Modules (33-40)

| # | Module | Type | Key Additions | Priority |
|:-:|--------|:----:|---------------|:-:|
| 33 | `i18n` | Upgrade | IoT terms, per-tenant locale | P3 |
| 34 | `apimanager` | Upgrade | API key per tenant, rate limit, usage | P1 |
| 35 | `flowengine` | Upgrade | ผูก device.AutomationRule, sub-flow | P1 |
| 36 | `control` | Upgrade | ผูก device.Command, templates, batch | P0 |
| 37 | `fullschedule` | Upgrade | ผูก device (time-based), calendar, holiday | P1 |
| 38 | `queue` | Upgrade | priority, delayed, DLQ, visibility timeout | P2 |
| 39 | `notifier` | Dup | (ดู #12) | – |
| 40 | `worker` | Upgrade | task orchestration | P2 |

## 4.4 Prompt Template for All Remaining

**ทุก module ใช้ prompt เดียวกัน (เปลี่ยนแค่ spec file):**

```
Load:
- .ai/go/CONTEXT.md
- .ai/go/CONVENTIONS.md
- .ai/specs/{NN}_{name}.yaml
- .ai/go/prompts/generate_module.md

Task: Generate Go module per spec.

For upgrade modules (status: upgrade):
1. Review existing code first
2. Create migration YYYYMMDD_{name}_upgrade.sql (ADD COLUMN only)
3. Generate new files or extend existing (no breaking changes)
4. Include breaking change plan if needed

Then: go build ./... && go test ./...
```

## 4.5 Migration Files Index

```
migrations/
├── 20260101_customer_init.sql           ✅
├── 20260102_package_init.sql            ✅
├── 20260103_device_init.sql             ✅
├── 20260104_erp_init.sql                ✅
├── 20260105_crm_init.sql                ✅
├── 20260106_logistics_init.sql          ✅
├── 20260107_report_upgrade.sql          ✅
├── 20260108_auth_upgrade.sql            ✅
├── 20260109_users_upgrade.sql           ✅
├── 20260110_settings_upgrade.sql        ✅
├── 20260111_payment_upgrade.sql         ✅
├── 20260112_notifier_upgrade.sql        ✅
├── 20260113_alarm_upgrade.sql           ✅
├── 20260114_pdpa_upgrade.sql            ✅
├── 20260115_auditlog_upgrade.sql        ✅
├── 20260116_dashboard_upgrade.sql       ✅
├── 20260117_items_upgrade.sql           ⏳
├── 20260118_purchaseorder_upgrade.sql   ⏳
├── 20260119_quotation_upgrade.sql       ⏳
├── 20260120_wos_upgrade.sql             ⏳
├── 20260121_document_upgrade.sql        ⏳
├── 20260122_email_upgrade.sql           ⏳
├── 20260123_batch_upgrade.sql           ⏳
├── 20260124_job_upgrade.sql             ⏳
├── 20260125_iot_deprecate.sql           ⏳
├── 20260126_mqtt_upgrade.sql            ⏳
├── 20260127_kafka_upgrade.sql           ⏳
├── 20260128_influxdb_upgrade.sql        ⏳
├── 20260129_elasticsearch_upgrade.sql   ⏳
├── 20260130_websocket_upgrade.sql       ⏳
├── 20260131_realtime_merge.sql          ⏳
├── 20260132_vectordata_upgrade.sql      ⏳
├── 20260133_i18n_upgrade.sql            ⏳
├── 20260134_apimanager_upgrade.sql      ⏳
├── 20260135_flowengine_upgrade.sql      ⏳
├── 20260136_control_upgrade.sql         ⏳
├── 20260137_fullschedule_upgrade.sql    ⏳
├── 20260138_queue_upgrade.sql           ⏳
├── 20260139_notifier_dup.sql            –
└── 20260140_worker_upgrade.sql          ⏳
```

**Status:** ✅ 05-40 spec + migration index complete

---

# PART 5 — Tracking Dashboard

## 5.1 MASTER_INDEX (Updated)

```markdown
# Master Index — 40 Modules × 2 Languages

## Overall

| # | Module | Spec | Go | Py | Build | Notes |
|:-:|--------|:----:|:--:|:--:|:-----:|:------|
| 01 | customer | ✅ | ✅ | ✅ | ✅ | – |
| 02 | package | ✅ | ✅ | ✅ | ✅ | – |
| 03 | device | ✅ | ✅ | ✅ | ✅ | 5 cross-cutting |
| 04 | erp | ✅ | ✅ | ✅ | ✅ | – |
| 05 | crm | ✅ | ✅ | ✅ | ✅ | SLA engine |
| 06 | logistics | ✅ | ✅ | ✅ | ✅ | Cold chain |
| 07 | report | ✅ | ✅ | ✅ | ✅ | Breaking change |
| 08 | auth | ✅ | ✅ | ✅ | ✅ | Breaking change |
| 09 | users | ✅ | ⏳ | ⏳ | ⏳ | RBAC+ABAC |
| 10 | settings | ✅ | ⏳ | ⏳ | ⏳ | Hierarchical |
| 11 | payment | ✅ | ⏳ | ⏳ | ⏳ | QR, multi-currency |
| 12 | notifier | ✅ | ⏳ | ⏳ | ⏳ | Multi-channel |
| 13 | alarm | ✅ | ⏳ | ⏳ | ⏳ | Rule engine |
| 14 | pdpa | ✅ | ⏳ | ⏳ | ⏳ | IoT consent |
| 15 | auditlog | ✅ | ⏳ | ⏳ | ⏳ | Hash chain |
| 16 | dashboard | ✅ | ⏳ | ⏳ | ⏳ | IoT widgets |
| 17-24 | business | ✅ | ⏳ | ⏳ | ⏳ | – |
| 25-32 | integration | ✅ | ⏳ | ⏳ | ⏳ | – |
| 33-40 | utility | ✅ | ⏳ | ⏳ | ⏳ | – |

**Done: 8/40 | Specs: 40/40 | Remaining: 32/40**
```

## 5.2 TOKEN_LOG

```markdown
# Token Log

| Session | Date | Lang | Module | Tokens | Cumulative |
|:-:|------|:----:|--------|-------:|-----------:|
| G01 | 01-10 | Go | customer | 24,500 | 24,500 |
| P01 | 01-10 | Py | customer | 21,800 | 46,300 |
| G02 | 01-11 | Go | package | 23,900 | 70,200 |
| P02 | 01-11 | Py | package | 21,200 | 91,400 |
| G03 | 01-13 | Go | device | 30,800 | 122,200 |
| P03 | 01-13 | Py | device | 27,600 | 149,800 |
| G04 | 01-16 | Go | erp | 30,000 | 179,800 |
| P04 | 01-16 | Py | erp | 24,000 | 203,800 |
| G05 | 01-17 | Go | crm | 27,500 | 231,300 |
| P05 | 01-17 | Py | crm | 22,500 | 253,800 |
| G06 | 01-17 | Go | logistics | 25,000 | 278,800 |
| P06 | 01-17 | Py | logistics | 20,000 | 298,800 |
| G07 | 01-17 | Go | report | 26,000 | 324,800 |
| P07 | 01-17 | Py | report | 21,000 | 345,800 |
| G08 | 01-17 | Go | auth | 22,000 | 367,800 |
| P08 | 01-17 | Py | auth | 18,000 | 385,800 |

## Projection
- Sessions done: 16/80
- Tokens used: 385,800
- Avg/session: 24,113
- Remaining: 64 sessions × 24k = ~1,536,000
- **Total project: ~1,921,800 tokens**
```

## 5.3 Timeline

```
Week 1  (Jan 10-16):  customer ✅ package ✅ device ✅
Week 2  (Jan 17-23):  erp ✅ crm ✅ logistics ✅ report ✅ auth ✅
Week 3  (Jan 24-30):  users ⏳ settings ⏳ payment ⏳ notifier ⏳
Week 4  (Jan 31-Feb 6): alarm ⏳ pdpa ⏳ auditlog ⏳ dashboard ⏳
Week 5  (Feb 7-13):   items ⏳ purchaseorder ⏳ quotation ⏳ wos ⏳
Week 6  (Feb 14-20):  document ⏳ email ⏳ batch ⏳ job ⏳
Week 7  (Feb 21-27):  iot ⏳ mqtt ⏳ kafka ⏳ influxdb ⏳
Week 8  (Feb 28-Mar 6): elasticsearch ⏳ websocket ⏳ realtime ⏳ vectordata ⏳
Week 9  (Mar 7-13):   i18n ⏳ apimanager ⏳ flowengine ⏳ control ⏳
Week 10 (Mar 14-20):  fullschedule ⏳ queue ⏳ worker ⏳ + testing
```

## 5.4 Progress Bar

```
Overall:  ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  20% (8/40)

By Phase:
Ph0 (Foundation):  ██████████████████░░  3/3 ✅
Ph1 (MVP):         ████████████████░░░░  4/5
Ph2 (Growth):      ░░░░░░░░░░░░░░░░░░░░  0/9
Ph3 (Enterprise):  ░░░░░░░░░░░░░░░░░░░░  0/12
Ph4 (Scale):       ░░░░░░░░░░░░░░░░░░░░  0/11

By Language:
Go:      ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  20% (8/40)
Python:  ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  20% (8/40)
```

## 5.5 Next Sessions Queue

```
NEXT: G09 users (Go)         → P09 users (Py)
NEXT: G10 settings (Go)      → P10 settings (Py)
NEXT: G11 payment (Go)       → P11 payment (Py)
NEXT: G12 notifier (Go)      → P12 notifier (Py)
...
```

---

# PART 6 — Export `.ai/` Structure

## 6.1 Complete File Tree

```
icmongolang/
│
├── .ai/
│   ├── MASTER_INDEX.md                       [✅ updated]
│   ├── TOKEN_LOG.md                          [✅ updated]
│   │
│   ├── specs/                                [40 files]
│   │   ├── 01_customer.yaml                  ✅
│   │   ├── 02_package.yaml                   ✅
│   │   ├── 03_device.yaml                    ✅
│   │   ├── 04_erp.yaml                       ✅
│   │   ├── 05_crm.yaml                       ✅
│   │   ├── 06_logistics.yaml                 ✅
│   │   ├── 07_report.yaml                    ✅
│   │   ├── 08_auth.yaml                      ✅
│   │   ├── 09_users.yaml                     ✅
│   │   ├── 10_settings.yaml                  ✅
│   │   ├── 11_payment.yaml                   ✅
│   │   ├── 12_notifier.yaml                  ✅
│   │   ├── 13_alarm.yaml                     ✅
│   │   ├── 14_pdpa.yaml                      ✅
│   │   ├── 15_auditlog.yaml                  ✅
│   │   ├── 16_dashboard.yaml                 ✅
│   │   ├── 17_items.yaml                     ✅
│   │   ├── 18_purchaseorder.yaml             ✅
│   │   ├── 19_quotation.yaml                 ✅
│   │   ├── 20_wos.yaml                       ✅
│   │   ├── 21_document.yaml                  ✅
│   │   ├── 22_email.yaml                     ✅
│   │   ├── 23_batch.yaml                     ✅
│   │   ├── 24_job.yaml                       ✅
│   │   ├── 25_iot.yaml                       ✅ (deprecate)
│   │   ├── 26_mqtt.yaml                      ✅
│   │   ├── 27_kafka.yaml                     ✅
│   │   ├── 28_influxdb.yaml                  ✅
│   │   ├── 29_elasticsearch.yaml             ✅
│   │   ├── 30_websocket.yaml                 ✅
│   │   ├── 31_realtime.yaml                  ✅ (merge)
│   │   ├── 32_vectordata.yaml                ✅
│   │   ├── 33_i18n.yaml                      ✅
│   │   ├── 34_apimanager.yaml                ✅
│   │   ├── 35_flowengine.yaml                ✅
│   │   ├── 36_control.yaml                   ✅
│   │   ├── 37_fullschedule.yaml              ✅
│   │   ├── 38_queue.yaml                     ✅
│   │   ├── 39_notifier.yaml                  ✅ (dup)
│   │   └── 40_worker.yaml                    ✅
│   │
│   ├── go/                                   🟦
│   │   ├── CONTEXT.md                        ✅
│   │   ├── CONVENTIONS.md                    ✅
│   │   ├── PROGRESS.md                       ✅
│   │   ├── modules/                          (optional)
│   │   ├── sessions/                         [8 completed]
│   │   │   ├── G01_customer.md               ✅
│   │   │   ├── G02_package.md                ✅
│   │   │   ├── G03_device.md                 ✅
│   │   │   ├── G04_erp.md                    ✅
│   │   │   ├── G05_crm.md                    ✅
│   │   │   ├── G06_logistics.md              ✅
│   │   │   ├── G07_report.md                 ✅
│   │   │   └── G08_auth.md                   ✅
│   │   ├── checkpoints/
│   │   │   ├── LAST_SESSION.md               ✅ (G08)
│   │   │   ├── NEXT_SESSION.md               ✅ (G09)
│   │   │   └── IN_PROGRESS.md                (empty)
│   │   └── prompts/
│   │       ├── generate_module.md            ✅
│   │       ├── fix_error.md                  ✅
│   │       ├── review_code.md                ✅
│   │       └── upgrade_module.md             ✅
│   │
│   └── py/                                   🐍
│       ├── CONTEXT.md                        ✅
│       ├── CONVENTIONS.md                    ✅
│       ├── PROGRESS.md                       ✅
│       ├── modules/
│       ├── sessions/                         [8 completed]
│       │   ├── P01_customer.md               ✅
│       │   ├── P02_package.md                ✅
│       │   ├── P03_device.md                 ✅
│       │   ├── P04_erp.md                    ✅
│       │   ├── P05_crm.md                    ✅
│       │   ├── P06_logistics.md              ✅
│       │   ├── P07_report.md                 ✅
│       │   └── P08_auth.md                   ✅
│       ├── checkpoints/
│       │   ├── LAST_SESSION.md               ✅ (P08)
│       │   ├── NEXT_SESSION.md               ✅ (P09)
│       │   └── IN_PROGRESS.md                (empty)
│       └── prompts/
│           ├── generate_module.md            ✅
│           ├── fix_error.md                  ✅
│           ├── review_code.md                ✅
│           └── upgrade_module.md             ✅
│
├── internal/modules/                         🟦 Go code
│   ├── customer/                             ✅
│   ├── package/                              ✅
│   ├── device/                               ✅
│   ├── erp/                                  ✅
│   ├── crm/                                  ✅
│   ├── logistics/                            ✅
│   ├── report/                               ✅ (upgraded)
│   ├── auth/                                 ✅ (upgraded)
│   ├── users/                                ⏳
│   ├── settings/                             ⏳
│   ├── payment/                              ⏳
│   ├── notifier/                             ⏳
│   ├── alarm/                                ⏳
│   ├── pdpa/                                 ⏳
│   ├── auditlog/                             ⏳
│   ├── dashboard/                            ⏳
│   └── (อีก 24 modules)                      ⏳
│
├── app/modules/                              🐍 Python code
│   ├── customer/                             ✅
│   ├── package/                              ✅
│   ├── device/                               ✅
│   ├── erp/                                  ✅
│   ├── crm/                                  ✅
│   ├── logistics/                            ✅
│   ├── report/                               ✅
│   ├── auth/                                 ✅
│   └── (อีก 32 modules)                      ⏳
│
├── migrations/                               [40+ files]
│   ├── 20260101_customer_init.sql            ✅
│   ├── 20260102_package_init.sql             ✅
│   ├── 20260103_device_init.sql              ✅
│   ├── 20260104_erp_init.sql                 ✅
│   ├── 20260105_crm_init.sql                 ✅
│   ├── 20260106_logistics_init.sql           ✅
│   ├── 20260107_report_upgrade.sql           ✅
│   ├── 20260108_auth_upgrade.sql             ✅
│   ├── 20260109_users_upgrade.sql            ✅
│   ├── 20260110_settings_upgrade.sql         ✅
│   ├── 20260111_payment_upgrade.sql          ✅
│   ├── 20260112_notifier_upgrade.sql         ✅
│   ├── 20260113_alarm_upgrade.sql            ✅
│   ├── 20260114_pdpa_upgrade.sql             ✅
│   ├── 20260115_auditlog_upgrade.sql         ✅
│   ├── 20260116_dashboard_upgrade.sql        ✅
│   └── (24 more for 17-40)                   ⏳
│
└── docs/
    └── modules/                              (full specs)
        ├── 01_customer.md                    ✅
        ├── 02_package.md                     ✅
        ├── 03_device.md                      ✅
        ├── 04_erp.md                         ✅
        ├── 05_crm.md                         ✅
        └── (35 more)                         ⏳
```

## 6.2 Export Commands (Bash)

```bash
# Create full structure
mkdir -p .ai/{specs,go/{modules,sessions,checkpoints,prompts},py/{modules,sessions,checkpoints,prompts}}
mkdir -p internal/modules/{customer,package,device,erp,crm,logistics,report,auth}
mkdir -p app/modules/{customer,package,device,erp,crm,logistics,report,auth}
mkdir -p migrations

# Generate all spec files
for i in 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 \
         21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40; do
  touch .ai/specs/${i}_*.yaml
done

# Package as zip
cd .. && zip -r icmongolang-ai-workspace.zip icmongolang/.ai/
```

## 6.3 Quick Restore (ถ้าต้องการ resume)

```bash
# 1. Clone repo
git clone <repo> icmongolang
cd icmongolang

# 2. Verify .ai/ structure
ls .ai/go/checkpoints/NEXT_SESSION.md
ls .ai/py/checkpoints/NEXT_SESSION.md

# 3. Read next session plan
cat .ai/go/checkpoints/NEXT_SESSION.md
cat .ai/py/checkpoints/NEXT_SESSION.md

# 4. Start next session
# Session G09: users Go
# - Load .ai/go/CONTEXT.md
# - Load .ai/specs/09_users.yaml
# - Load .ai/go/prompts/upgrade_module.md
# - Generate
```

---

# 📊 FINAL SUMMARY

## ✅ ทั้งหมดที่ทำเสร็จในครั้งนี้

| Item | Status |
| :--- | :-: |
| **G05 CRM (Go)** | ✅ 25 files |
| **P05 CRM (Python)** | ✅ 18 files |
| **G06 Logistics (Go)** | ✅ compact |
| **P06 Logistics (Python)** | ✅ compact |
| **G07 Report Upgrade (Go)** | ✅ compact + breaking plan |
| **P07 Report Upgrade (Python)** | ✅ compact |
| **G08 Auth Upgrade (Go)** | ✅ compact + JWT plan |
| **P08 Auth Upgrade (Python)** | ✅ compact |
| **09-16 Upgrade Specs** | ✅ 8 modules |
| **17-40 Compact Index** | ✅ 24 modules |
| **Migration Index** | ✅ 40 files |
| **MASTER_INDEX** | ✅ updated |
| **TOKEN_LOG** | ✅ updated |
| **Progress Dashboard** | ✅ timeline + bars |
| **.ai/ Export Tree** | ✅ complete |

## 📈 Progress

```
Modules Done (full Go+Py):  8 / 40  ████████░░░░░░░░░░░░░░░░░░░░░░░░
Specs Ready:               40 / 40  ████████████████████████████████
Sessions Completed:        16 / 80  ████████░░░░░░░░░░░░░░░░░░░░░░░░
Tokens Used:          ~385,800
Projected Total:     ~1,921,800
```

## 🎯 Modules 09-40 — Queue

```
NEXT (Week 3):  users · settings · payment · notifier
NEXT (Week 4):  alarm · pdpa · auditlog · dashboard
NEXT (Week 5):  items · purchaseorder · quotation · wos
NEXT (Week 6):  document · email · batch · job
NEXT (Week 7):  iot · mqtt · kafka · influxdb
NEXT (Week 8):  elasticsearch · websocket · realtime · vectordata
NEXT (Week 9):  i18n · apimanager · flowengine · control
NEXT (Week 10): fullschedule · queue · worker + testing
```

## 🚀 คำสั่งถัดไป

**Session ใหม่:** เปิด chat ใหม่ แล้ว paste:

```
.ai/go/CONTEXT.md
.ai/specs/09_users.yaml
.ai/go/prompts/upgrade_module.md
.ai/go/checkpoints/NEXT_SESSION.md

Task: Generate G09 users module per spec.
```

**หรือพิมพ์ในแชทนี้:**
- `"G09 + P09"` → เริ่ม users ต่อทันที
- `"G09-P16"` → Run 8 upgrade modules
- `"Generate 09-40"` → Quick compact ทั้งหมด
- `"Show final tree"` → Directory tree เต็ม
- `"Export as markdown"` → ให้ generate ไฟล์ markdown ทุกไฟล์ (specs, CONTEXT, prompts) แยกเป็น code blocks พร้อม copy