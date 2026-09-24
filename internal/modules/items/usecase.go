package items

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type ItemUseCaseI interface {
	internal.UseCaseI[models.Item]
	GetMultiByOwnerId(ctx context.Context, ownerId uuid.UUID, limit, offset int) ([]*models.Item, error)
	CreateWithOwner(ctx context.Context, ownerId uuid.UUID, exp *models.Item) (*models.Item, error)
	DeleteWithoutGet(ctx context.Context, id uuid.UUID) error
	// Count methods for pagination
	Count(ctx context.Context) (int64, error)
	CountByOwnerId(ctx context.Context, ownerId uuid.UUID) (int64, error)
}
