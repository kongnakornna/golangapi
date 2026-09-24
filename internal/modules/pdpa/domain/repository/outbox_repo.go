package repository

import (
	"context"
	"icmongolang/internal/modules/pdpa/domain/event"
	"time"

	"github.com/google/uuid"
)

type OutboxRepository interface {
	Publish(ctx context.Context, evt event.OutboxEvent) error
	ClaimPending(ctx context.Context, limit int, handler string) ([]event.OutboxEvent, error)
	MarkSuccess(ctx context.Context, eventID uuid.UUID, handler string) error
	MarkFailed(ctx context.Context, eventID uuid.UUID, handler string, attempts int) error
	MarkPending(ctx context.Context, eventID uuid.UUID, handler string) error
	CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error)
}