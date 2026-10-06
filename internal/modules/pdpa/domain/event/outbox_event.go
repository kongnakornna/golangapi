package event

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEventStatus int8

const (
	OutboxEventStatusPending    OutboxEventStatus = 1
	OutboxEventStatusProcessing OutboxEventStatus = 2
	OutboxEventStatusSuccess    OutboxEventStatus = 3
	OutboxEventStatusFailed     OutboxEventStatus = 4
)

func (s OutboxEventStatus) IsValid() bool {
	switch s {
	case OutboxEventStatusPending, OutboxEventStatusProcessing, OutboxEventStatusSuccess, OutboxEventStatusFailed:
		return true
	default:
		return false
	}
}

func (s OutboxEventStatus) String() string {
	switch s {
	case OutboxEventStatusPending:
		return "PENDING"
	case OutboxEventStatusProcessing:
		return "PROCESSING"
	case OutboxEventStatusSuccess:
		return "SUCCESS"
	case OutboxEventStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type OutboxEvent struct {
	ID           uuid.UUID
	Topic        string
	Key          string
	Payload      []byte
	Status       OutboxEventStatus
	Attempts     int
	CreatedAt    time.Time
	DispatchedAt *time.Time
}
