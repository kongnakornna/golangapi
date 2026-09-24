package i18n

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// I18nPgRepository defines data access methods for translations.
// ดึงข้อมูลการแปลภาษาจากฐานข้อมูล
type I18nPgRepository interface {
	internal.PgRepository[models.Translation]
	GetByLocaleAndKey(ctx context.Context, locale, key string) (*models.Translation, error)
	GetByLocale(ctx context.Context, locale string) ([]*models.Translation, error)
}
