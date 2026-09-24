package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/i18n"
	"icmongolang/internal/repository"

	"gorm.io/gorm"
)

// I18nPgRepo implements i18n.I18nPgRepository.
// รีโพสิทอรีสำหรับการแปลภาษา
type I18nPgRepo struct {
	repository.PgRepo[models.Translation]
	DB *gorm.DB
}

// CreateI18nPgRepository creates a new i18n repository.
// สร้างรีโพสิทอรีสำหรับการแปลภาษา
func CreateI18nPgRepository(db *gorm.DB) i18n.I18nPgRepository {
	return &I18nPgRepo{
		PgRepo: repository.CreatePgRepo[models.Translation](db),
		DB:     db,
	}
}

// GetByLocaleAndKey finds a translation by locale and key.
// ค้นหาคำแปลตามภาษาและคีย์
func (r *I18nPgRepo) GetByLocaleAndKey(ctx context.Context, locale, key string) (*models.Translation, error) {
	var t models.Translation
	if result := r.DB.WithContext(ctx).Where("locale = ? AND key = ?", locale, key).First(&t); result.Error != nil {
		return nil, result.Error
	}
	return &t, nil
}

// GetByLocale returns all translations for a specific locale.
// ดึงคำแปลทั้งหมดตามภาษา
func (r *I18nPgRepo) GetByLocale(ctx context.Context, locale string) ([]*models.Translation, error) {
	var translations []*models.Translation
	r.DB.WithContext(ctx).Where("locale = ?", locale).Find(&translations)
	return translations, nil
}
