package i18n

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// I18nUseCaseI defines business logic methods for translations.
// อินเทอร์เฟซธุรกิจสำหรับการแปลภาษา
type I18nUseCaseI interface {
	internal.UseCaseI[models.Translation]
	GetByLocaleAndKey(ctx context.Context, locale, key string) (*models.Translation, error)
	GetByLocale(ctx context.Context, locale string) ([]*models.Translation, error)
}
