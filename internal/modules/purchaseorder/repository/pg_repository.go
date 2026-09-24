package repository

import (
	"context"
	"fmt"

	"icmongolang/internal/modules/purchaseorder"
	pomodels "icmongolang/internal/modules/purchaseorder/models"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type purchaseOrderPgRepository struct {
	repository.PgRepo[pomodels.PurchaseOrderHeader]
	db *gorm.DB
}

func CreatePurchaseOrderPgRepository(db *gorm.DB) purchaseorder.PurchaseOrderPgRepository {
	return &purchaseOrderPgRepository{
		PgRepo: repository.CreatePgRepo[pomodels.PurchaseOrderHeader](db),
		db:     db,
	}
}

func (r *purchaseOrderPgRepository) CreateWithDetails(ctx context.Context, header *pomodels.PurchaseOrderHeader, details []*pomodels.PurchaseOrderDetail) (*pomodels.PurchaseOrderHeader, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(header).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, detail := range details {
		detail.PoHeaderID = header.ID
		detail.UserID = header.UserID
		detail.WhitelabelID = header.WhitelabelID
		if err := tx.Create(detail).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	subtotal := 0.0
	for _, d := range details {
		subtotal += d.NetPrice
	}
	discount := 0.0
	if header.DiscountType != nil && *header.DiscountType == "PERCENTAGE" {
		discount = subtotal * (header.DiscountValue / 100.0)
	} else {
		discount = header.DiscountValue
	}
	taxAmount := (subtotal - discount) * (header.TaxRate / 100.0)
	total := subtotal - discount + taxAmount + header.ShippingCost

	if err := tx.Model(header).Updates(map[string]interface{}{
		"subtotal":   subtotal,
		"tax_amount": taxAmount,
		"total":      total,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return r.GetByIDWithDetails(ctx, header.ID)
}

func (r *purchaseOrderPgRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	var header pomodels.PurchaseOrderHeader
	if err := r.db.WithContext(ctx).
		Preload("Details").
		Where("id = ? AND deleted = ?", id.String(), false).
		First(&header).Error; err != nil {
		return nil, err
	}
	return &header, nil
}

func (r *purchaseOrderPgRepository) List(ctx context.Context, limit, offset int, supplierID *uuid.UUID, status string, startDate, endDate string) ([]*pomodels.PurchaseOrderHeader, error) {
	var headers []*pomodels.PurchaseOrderHeader
	query := r.db.WithContext(ctx).Where("deleted = ?", false)

	if supplierID != nil {
		query = query.Where("supplier_id = ?", supplierID.String())
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("po_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("po_date <= ?", endDate)
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&headers).Error; err != nil {
		return nil, err
	}
	return headers, nil
}

func (r *purchaseOrderPgRepository) CountByFilter(ctx context.Context, supplierID *uuid.UUID, status string, startDate, endDate string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&pomodels.PurchaseOrderHeader{}).Where("deleted = ?", false)

	if supplierID != nil {
		query = query.Where("supplier_id = ?", supplierID.String())
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("po_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("po_date <= ?", endDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *purchaseOrderPgRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]interface{}) error {
	updates := map[string]interface{}{
		"status": status,
	}
	for k, v := range fields {
		updates[k] = v
	}

	if result := r.db.WithContext(ctx).Model(&pomodels.PurchaseOrderHeader{}).
		Where("id = ?", id.String()).
		Updates(updates); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *purchaseOrderPgRepository) SaveStatusHistory(ctx context.Context, history *pomodels.PurchaseOrderStatusHistory) error {
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		return fmt.Errorf("บันทึกประวัติสถานะไม่สำเร็จ: %w", err)
	}
	return nil
}

func (r *purchaseOrderPgRepository) GetStatusHistory(ctx context.Context, poID uuid.UUID) ([]*pomodels.PurchaseOrderStatusHistory, error) {
	var history []*pomodels.PurchaseOrderStatusHistory
	if err := r.db.WithContext(ctx).
		Where("po_header_id = ?", poID.String()).
		Order("changed_at ASC").
		Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *purchaseOrderPgRepository) GetDetailsByHeaderID(ctx context.Context, headerID uuid.UUID) ([]*pomodels.PurchaseOrderDetail, error) {
	var details []*pomodels.PurchaseOrderDetail
	if err := r.db.WithContext(ctx).
		Where("po_header_id = ?", headerID.String()).
		Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *purchaseOrderPgRepository) UpdateDetail(ctx context.Context, detail *pomodels.PurchaseOrderDetail) error {
	if err := r.db.WithContext(ctx).Save(detail).Error; err != nil {
		return err
	}
	return nil
}

func (r *purchaseOrderPgRepository) UpdateDetails(ctx context.Context, details []*pomodels.PurchaseOrderDetail) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, detail := range details {
		if err := tx.Save(detail).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *purchaseOrderPgRepository) IncrementReceivedQuantity(ctx context.Context, detailID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&pomodels.PurchaseOrderDetail{}).
		Where("id = ?", detailID).
		Update("quantity_received", gorm.Expr("quantity_received + ?", quantity)).Error
}
