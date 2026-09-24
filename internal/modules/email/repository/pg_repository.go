package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/email"
	"icmongolang/internal/repository"

	"gorm.io/gorm"
)

// EmailPgRepo implements email.EmailPgRepository.
// รีโพสิทอรีสำหรับอีเมล
type EmailPgRepo struct {
	repository.PgRepo[models.EmailLog]
	DB *gorm.DB
}

// CreateEmailPgRepository creates a new email repository.
// สร้างรีโพสิทอรีสำหรับอีเมล
func CreateEmailPgRepository(db *gorm.DB) email.EmailPgRepository {
	return &EmailPgRepo{
		PgRepo: repository.CreatePgRepo[models.EmailLog](db),
		DB:     db,
	}
}

// GetConfig returns the active email configuration.
// ดึงการตั้งค่าอีเมลล่าสุด
func (r *EmailPgRepo) GetConfig(ctx context.Context) (*models.EmailConfig, error) {
	var cfg models.EmailConfig
	if result := r.DB.WithContext(ctx).Where("is_active = ?", true).First(&cfg); result.Error != nil {
		return nil, result.Error
	}
	return &cfg, nil
}

// UpdateConfig updates the email configuration.
// อัปเดตการตั้งค่าอีเมล
func (r *EmailPgRepo) UpdateConfig(ctx context.Context, cfg *models.EmailConfig) (*models.EmailConfig, error) {
	if result := r.DB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(cfg); result.Error != nil {
		return nil, result.Error
	}
	return cfg, nil
}
