package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/quotation"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type quotationUseCase struct {
	usecase.UseCase[models.Quotation]
	pgRepo quotation.QuotationPgRepository
}

func CreateQuotationUseCaseI(
	pgRepo quotation.QuotationPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) quotation.QuotationUseCaseI {
	return &quotationUseCase{
		UseCase: usecase.CreateUseCase[models.Quotation](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *quotationUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}

func (u *quotationUseCase) GetQuotationParts(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationPart, error) {
	return u.pgRepo.GetQuotationParts(ctx, quotationID)
}

func (u *quotationUseCase) GetQuotationServices(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationService, error) {
	return u.pgRepo.GetQuotationServices(ctx, quotationID)
}
