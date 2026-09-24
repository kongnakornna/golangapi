package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/quotation"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuotationPgRepo struct {
	repository.PgRepo[models.Quotation]
}

func CreateQuotationPgRepository(db *gorm.DB) quotation.QuotationPgRepository {
	return &QuotationPgRepo{
		PgRepo: repository.CreatePgRepo[models.Quotation](db),
	}
}

func (r *QuotationPgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Quotation{}).Count(&count).Error
	return count, err
}

func (r *QuotationPgRepo) GetQuotationParts(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationPart, error) {
	var parts []models.QuotationPart
	err := r.DB.WithContext(ctx).Where("quotation_id = ?", quotationID).Order("created_at ASC").Find(&parts).Error
	return parts, err
}

func (r *QuotationPgRepo) GetQuotationServices(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationService, error) {
	var services []models.QuotationService
	err := r.DB.WithContext(ctx).Where("quotation_id = ?", quotationID).Order("created_at ASC").Find(&services).Error
	return services, err
}
