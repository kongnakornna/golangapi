package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/wos"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"
)

type wosUseCase struct {
	usecase.UseCase[models.WosOrder]
	pgRepo wos.WosPgRepository
}

// CreateWosUseCaseI creates a new WOS use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจระบบสั่งซื้อออนไลน์
func CreateWosUseCaseI(
	pgRepo wos.WosPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) wos.WosUseCaseI {
	return &wosUseCase{
		UseCase: usecase.CreateUseCase[models.WosOrder](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *wosUseCase) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.WosOrder, error) {
	return u.pgRepo.GetByOrderNumber(ctx, orderNumber)
}

func (u *wosUseCase) UpdateStatus(ctx context.Context, id string, status string) error {
	return u.pgRepo.UpdateStatus(ctx, id, status)
}
