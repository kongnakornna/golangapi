package repository

import (
	"context"

	"icmongolang/internal/models"
	"gorm.io/gorm"
)

// Repository provides the API-key store (reuses the existing sd_api_key table).
type Repository interface {
	Create(ctx context.Context, key *models.SdApiKey) error
	GetByID(ctx context.Context, id int) (*models.SdApiKey, error)
	GetByKey(ctx context.Context, apiKey string) (*models.SdApiKey, error)
	SetActive(ctx context.Context, id int, active bool) error
	IncrementUsage(ctx context.Context, id int) error
}

type pgRepository struct {
	db *gorm.DB
}

// NewPgRepository creates a DB-backed API manager repository.
func NewPgRepository(db *gorm.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Create(ctx context.Context, key *models.SdApiKey) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *pgRepository) GetByID(ctx context.Context, id int) (*models.SdApiKey, error) {
	var k models.SdApiKey
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&k).Error; err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *pgRepository) GetByKey(ctx context.Context, apiKey string) (*models.SdApiKey, error) {
	var k models.SdApiKey
	if err := r.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&k).Error; err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *pgRepository) SetActive(ctx context.Context, id int, active bool) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Model(&models.SdApiKey{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *pgRepository) IncrementUsage(ctx context.Context, id int) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Model(&models.SdApiKey{}).
		Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).
		UpdateColumn("last_used_at", gorm.Expr("now()")).Error
}
