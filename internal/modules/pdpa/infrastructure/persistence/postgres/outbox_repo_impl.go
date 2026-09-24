// Package postgres implements the pdpa module persistence layer on PostgreSQL.
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"icmongolang/internal/modules/pdpa/domain/event"
	repository "icmongolang/internal/modules/pdpa/domain/repository"
	base "icmongolang/internal/repository"
)

const defaultHandler = "pdpa-outbox"

type outboxEventRow struct {
	ID           uuid.UUID              `gorm:"column:id;primaryKey"`
	EventID      uuid.UUID              `gorm:"column:event_id"`
	Topic        string                 `gorm:"column:topic"`
	Key          *string                `gorm:"column:key"`
	Payload      json.RawMessage        `gorm:"column:payload;type:jsonb"`
	Handler      string                 `gorm:"column:handler"`
	Status       event.OutboxEventStatus `gorm:"column:status"`
	Attempts     int                    `gorm:"column:attempts"`
	CreatedAt    time.Time              `gorm:"column:created_at"`
	UpdatedAt    time.Time              `gorm:"column:updated_at"`
	DispatchedAt *time.Time             `gorm:"column:dispatched_at"`
	DeletedAt    gorm.DeletedAt         `gorm:"column:deleted_at"`
}

func (outboxEventRow) TableName() string {
	return "pdpa_outbox"
}

func (r outboxEventRow) toEvent() event.OutboxEvent {
	return event.OutboxEvent{
		ID:           r.EventID,
		Topic:        r.Topic,
		Key:          stringVal(r.Key),
		Payload:      []byte(r.Payload),
		Status:       r.Status,
		Attempts:     r.Attempts,
		CreatedAt:    r.CreatedAt,
		DispatchedAt: r.DispatchedAt,
	}
}

type outboxPgRepository struct {
	base.PgRepo[event.OutboxEvent]
	db *gorm.DB
}

func CreateOutboxPgRepository(db *gorm.DB) repository.OutboxRepository {
	return &outboxPgRepository{
		PgRepo: base.CreatePgRepo[event.OutboxEvent](db),
		db:     db,
	}
}

func (r *outboxPgRepository) Publish(ctx context.Context, evt event.OutboxEvent) error {
	row := outboxEventRow{
		EventID: evt.ID,
		Topic:   evt.Topic,
		Key:     stringPtr(evt.Key),
		Payload: evt.Payload,
		Handler: defaultHandler,
		Status:  event.OutboxEventStatusPending,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&row).Error
}

func (r *outboxPgRepository) ClaimPending(ctx context.Context, limit int, handler string) ([]event.OutboxEvent, error) {
	var rows []outboxEventRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("handler = ? AND status = ?", handler, event.OutboxEventStatusPending).
			Order("created_at ASC").
			Limit(limit).
			Find(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		eventIDs := make([]uuid.UUID, 0, len(rows))
		for _, row := range rows {
			eventIDs = append(eventIDs, row.EventID)
		}
		return tx.Model(&outboxEventRow{}).
			Where("event_id IN ? AND handler = ?", eventIDs, handler).
			Update("status", event.OutboxEventStatusProcessing).Error
	})
	if err != nil {
		return nil, err
	}
	events := make([]event.OutboxEvent, 0, len(rows))
	for i := range rows {
		events = append(events, rows[i].toEvent())
	}
	return events, nil
}

func (r *outboxPgRepository) MarkSuccess(ctx context.Context, eventID uuid.UUID, handler string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&outboxEventRow{}).
		Where("event_id = ? AND handler = ?", eventID, handler).
		Updates(map[string]interface{}{
			"status":        event.OutboxEventStatusSuccess,
			"dispatched_at": now,
			"updated_at":    now,
		}).Error
}

func (r *outboxPgRepository) MarkFailed(ctx context.Context, eventID uuid.UUID, handler string, attempts int) error {
	return r.db.WithContext(ctx).
		Model(&outboxEventRow{}).
		Where("event_id = ? AND handler = ?", eventID, handler).
		Updates(map[string]interface{}{
			"status":     event.OutboxEventStatusFailed,
			"attempts":   attempts,
			"updated_at": time.Now(),
		}).Error
}

func (r *outboxPgRepository) MarkPending(ctx context.Context, eventID uuid.UUID, handler string) error {
	return r.db.WithContext(ctx).
		Model(&outboxEventRow{}).
		Where("event_id = ? AND handler = ?", eventID, handler).
		Updates(map[string]interface{}{
			"status":     event.OutboxEventStatusPending,
			"updated_at": time.Now(),
		}).Error
}

func (r *outboxPgRepository) CleanupPublished(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("status = ? AND dispatched_at < ?", event.OutboxEventStatusSuccess, cutoff).
		Delete(&outboxEventRow{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}