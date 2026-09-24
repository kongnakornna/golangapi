package persistence

import (
	"context"

	"gorm.io/gorm"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
)

type pgRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) repository.AuditRepository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Save(ctx context.Context, audit *entity.AuditTrail) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(audit).Error
}