package repository

import (
	"context"

	"icmongolang/internal/models"
	"gorm.io/gorm"
)

// Repository provides flow definition persistence.
type Repository interface {
	Create(ctx context.Context, flow *models.FlowDefinition) error
	GetByID(ctx context.Context, id string) (*models.FlowDefinition, error)
	Update(ctx context.Context, flow *models.FlowDefinition) error
}

type pgRepository struct {
	db *gorm.DB
}

// NewPgRepository creates a DB-backed flow engine repository.
func NewPgRepository(db *gorm.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, flow *models.FlowDefinition) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(flow).Error
}

func (r *pgRepository) GetByID(ctx context.Context, id string) (*models.FlowDefinition, error) {
	var f models.FlowDefinition
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *pgRepository) Update(ctx context.Context, flow *models.FlowDefinition) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(flow).Error
}
