package repository

import (
	"context"
	"fmt"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/payment"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentPgRepo struct {
	repository.PgRepo[models.Payment]
}

func CreatePaymentPgRepository(db *gorm.DB) payment.PaymentPgRepository {
	return &PaymentPgRepo{
		PgRepo: repository.CreatePgRepo[models.Payment](db),
	}
}

func (r *PaymentPgRepo) GetByInvoiceId(ctx context.Context, invoiceId uuid.UUID) (*models.Payment, error) {
	var obj *models.Payment
	result := r.DB.WithContext(ctx).Where("invoice_id = ?", invoiceId.String()).First(&obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PaymentPgRepo) Search(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*models.Payment, error) {
	var objs []*models.Payment
	query := r.DB.WithContext(ctx).Model(&models.Payment{})

	if customerId, ok := filter["customer_id"]; ok {
		query = query.Where("customer_id = ?", customerId)
	}
	if invoiceId, ok := filter["invoice_id"]; ok {
		query = query.Where("invoice_id = ?", invoiceId)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if paymentMethodId, ok := filter["payment_method_id"]; ok {
		query = query.Where("payment_method_id = ?", paymentMethodId)
	}
	if dateFrom, ok := filter["date_from"]; ok {
		query = query.Where("payment_date >= ?", dateFrom)
	}
	if dateTo, ok := filter["date_to"]; ok {
		query = query.Where("payment_date <= ?", dateTo)
	}

	if err := query.Where("deleted = ?", false).Limit(limit).Offset(offset).Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *PaymentPgRepo) GetByCustomerId(ctx context.Context, customerId uuid.UUID, limit, offset int) ([]*models.Payment, error) {
	var objs []*models.Payment
	if err := r.DB.WithContext(ctx).
		Where("customer_id = ?", customerId.String()).
		Where("deleted = ?", false).
		Limit(limit).
		Offset(offset).
		Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *PaymentPgRepo) GetOutstandingByCustomerId(ctx context.Context, customerId uuid.UUID) ([]*models.OutstandingBalance, error) {
	var objs []*models.OutstandingBalance
	if err := r.DB.WithContext(ctx).
		Where("customer_id = ?", customerId.String()).
		Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *PaymentPgRepo) GetPaymentHistory(ctx context.Context, customerId uuid.UUID) ([]*models.PaymentHistory, error) {
	var objs []*models.PaymentHistory
	if err := r.DB.WithContext(ctx).
		Joins("JOIN t_payment ON t_payment.id = t_payment_history.payment_id").
		Where("t_payment.customer_id = ?", customerId.String()).
		Order("t_payment_history.changed_at DESC").
		Find(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *PaymentPgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Payment{}).Where("deleted = ?", false).Count(&count).Error
	return count, err
}

func (r *PaymentPgRepo) CountByFilter(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.DB.WithContext(ctx).Model(&models.Payment{})

	if customerId, ok := filter["customer_id"]; ok {
		query = query.Where("customer_id = ?", customerId)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("status = ?", status)
	}

	err := query.Where("deleted = ?", false).Count(&count).Error
	return count, err
}

// Receipt methods
func (r *PaymentPgRepo) GetReceipt(ctx context.Context, id uuid.UUID) (*models.Receipt, error) {
	var obj *models.Receipt
	result := r.DB.WithContext(ctx).First(&obj, "id = ?", id.String())
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PaymentPgRepo) GetReceiptByPaymentId(ctx context.Context, paymentId uuid.UUID) (*models.Receipt, error) {
	var obj *models.Receipt
	result := r.DB.WithContext(ctx).Where("payment_id = ?", paymentId.String()).First(&obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PaymentPgRepo) CancelReceipt(ctx context.Context, id uuid.UUID, reason string) error {
	result := r.DB.WithContext(ctx).Model(&models.Receipt{}).
		Where("id = ?", id.String()).
		Updates(map[string]interface{}{
			"status": "CANCELLED",
			"notes":  gorm.Expr("COALESCE(notes, '') || ?", fmt.Sprintf("\nCancelled: %s", reason)),
		})
	return result.Error
}
