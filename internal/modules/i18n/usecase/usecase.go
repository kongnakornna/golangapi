package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/i18n"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"
)

type i18nUseCase struct {
	usecase.UseCase[models.Translation]
	pgRepo i18n.I18nPgRepository
}

// CreateI18nUseCaseI creates a new i18n use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจการแปลภาษา
func CreateI18nUseCaseI(
	pgRepo i18n.I18nPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) i18n.I18nUseCaseI {
	return &i18nUseCase{
		UseCase: usecase.CreateUseCase[models.Translation](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *i18nUseCase) GetByLocaleAndKey(ctx context.Context, locale, key string) (*models.Translation, error) {
	return u.pgRepo.GetByLocaleAndKey(ctx, locale, key)
}

func (u *i18nUseCase) GetByLocale(ctx context.Context, locale string) ([]*models.Translation, error) {
	return u.pgRepo.GetByLocale(ctx, locale)
}
