package repository

import (
	"context"

	"icmongolang/internal/modules/items"
	"icmongolang/internal/models"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemPgRepo struct {
	repository.PgRepo[models.Item]
}

func CreateItemPgRepository(db *gorm.DB) items.ItemPgRepository {
	return &ItemPgRepo{
		PgRepo: repository.CreatePgRepo[models.Item](db),
	}
}

func (r *ItemPgRepo) GetMultiByOwnerId(ctx context.Context, ownerId uuid.UUID, limit, offset int) ([]*models.Item, error) {
	var objs []*models.Item
	r.DB.WithContext(ctx).Where("owner_id = ?", ownerId.String()).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *ItemPgRepo) CreateWithOwner(ctx context.Context, ownerId uuid.UUID, exp *models.Item) (*models.Item, error) {
	exp.OwnerId = ownerId
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *ItemPgRepo) DeleteWithoutGet(ctx context.Context, id uuid.UUID) error {
	if result := r.DB.WithContext(ctx).Delete(&models.Item{}, "id = ?", id.String()); result.Error != nil {
		return result.Error
	}
	return nil
}

// Count returns total number of items (for superuser)
func (r *ItemPgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Item{}).Count(&count).Error
	return count, err
}

// CountByOwnerId returns total number of items owned by a specific user
func (r *ItemPgRepo) CountByOwnerId(ctx context.Context, ownerId uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Item{}).Where("owner_id = ?", ownerId.String()).Count(&count).Error
	return count, err
}
