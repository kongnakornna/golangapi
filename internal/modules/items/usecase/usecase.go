package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/modules/items"
	"icmongolang/internal/models"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type itemUseCase struct {
	usecase.UseCase[models.Item]
	pgRepo items.ItemPgRepository
}

func CreateItemUseCaseI(
	pgRepo items.ItemPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) items.ItemUseCaseI {
	return &itemUseCase{
		UseCase: usecase.CreateUseCase[models.Item](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *itemUseCase) GetMultiByOwnerId(ctx context.Context, ownerId uuid.UUID, limit, offset int) ([]*models.Item, error) {
	return u.pgRepo.GetMultiByOwnerId(ctx, ownerId, limit, offset)
}

func (u *itemUseCase) CreateWithOwner(ctx context.Context, ownerId uuid.UUID, exp *models.Item) (*models.Item, error) {
	return u.pgRepo.CreateWithOwner(ctx, ownerId, exp)
}

func (u *itemUseCase) DeleteWithoutGet(ctx context.Context, id uuid.UUID) error {
	return u.pgRepo.DeleteWithoutGet(ctx, id)
}

func (u *itemUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}

func (u *itemUseCase) CountByOwnerId(ctx context.Context, ownerId uuid.UUID) (int64, error) {
	return u.pgRepo.CountByOwnerId(ctx, ownerId)
}
