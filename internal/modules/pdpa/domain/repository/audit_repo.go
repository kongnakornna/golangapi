package repository

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type AuditFilter struct {
	UserID *string
	Action string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

type AuditRepository interface {
	Save(ctx context.Context, a *entity.AuditTrail) error
	FindByUser(ctx context.Context, userID string, limit int) ([]*entity.AuditTrail, error)
	FindByAction(ctx context.Context, action string, limit int) ([]*entity.AuditTrail, error)
	FindByFilter(ctx context.Context, filter AuditFilter) ([]*entity.AuditTrail, error)
	CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error)
}
