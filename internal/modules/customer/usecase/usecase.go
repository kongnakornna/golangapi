package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/customer"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// ─── Customer Use Case ────────────────────────────────────────────────────

type customerUseCase struct {
	usecase.UseCase[models.Customer]
	pgRepo customer.CustomerPgRepository
}

func CreateCustomerUseCaseI(
	pgRepo customer.CustomerPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) customer.CustomerUseCaseI {
	return &customerUseCase{
		UseCase: usecase.CreateUseCase[models.Customer](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *customerUseCase) GetMultiByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Customer, error) {
	return u.pgRepo.GetMultiByUserID(ctx, userID, limit, offset)
}

func (u *customerUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}

func (u *customerUseCase) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return u.pgRepo.CountByUserID(ctx, userID)
}

// ─── Car Use Case ─────────────────────────────────────────────────────────

type carUseCase struct {
	usecase.UseCase[models.Car]
	pgRepo customer.CarPgRepository
}

func CreateCarUseCaseI(
	pgRepo customer.CarPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) customer.CarUseCaseI {
	return &carUseCase{
		UseCase: usecase.CreateUseCase[models.Car](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *carUseCase) GetMultiByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Car, error) {
	return u.pgRepo.GetMultiByCustomerID(ctx, customerID, limit, offset)
}

func (u *carUseCase) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error) {
	return u.pgRepo.CountByCustomerID(ctx, customerID)
}
