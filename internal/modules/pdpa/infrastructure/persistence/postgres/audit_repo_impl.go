// Package postgres ให้ persistence adapter สำหรับ audit trail
package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	base "icmongolang/internal/repository"
)

type AuditRepoImpl struct {
	base.PgRepo[entity.AuditTrail]
	db *gorm.DB
}

func NewAuditPgRepository(db *gorm.DB) repository.AuditRepository {
	return &AuditRepoImpl{
		PgRepo: base.CreatePgRepo[entity.AuditTrail](db),
		db:     db,
	}
}

func (r *AuditRepoImpl) Save(ctx context.Context, audit *entity.AuditTrail) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(audit).Error
}

func (r *AuditRepoImpl) FindByUser(ctx context.Context, userID string, limit int) ([]*entity.AuditTrail, error) {
	if limit <= 0 {
		limit = 100
	}
	var out []*entity.AuditTrail
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Find(&out).Error
	return out, err
}

func (r *AuditRepoImpl) FindByAction(ctx context.Context, action string, limit int) ([]*entity.AuditTrail, error) {
	if limit <= 0 {
		limit = 100
	}
	var out []*entity.AuditTrail
	err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).Find(&out).Error
	return out, err
}

func (r *AuditRepoImpl) FindByFilter(ctx context.Context, filter repository.AuditFilter) ([]*entity.AuditTrail, error) {
	q := r.db.WithContext(ctx).Model(&entity.AuditTrail{})
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != "" {
		q = q.Where("action = ?", filter.Action)
	}
	if filter.From != nil {
		q = q.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		q = q.Where("created_at <= ?", *filter.To)
	}
	if filter.Limit > 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}
	var out []*entity.AuditTrail
	err := q.Order("created_at DESC").Find(&out).Error
	return out, err
}

func (r *AuditRepoImpl) CountByAction(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	type row struct {
		Action string
		Count  int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Model(&entity.AuditTrail{}).
		Select("action, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", from, to).
		Group("action").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		out[r.Action] = r.Count
	}
	return out, nil
}
