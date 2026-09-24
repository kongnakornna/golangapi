package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/wos"
	"icmongolang/internal/repository"

	"gorm.io/gorm"
)

// WosPgRepo implements wos.WosPgRepository.
// รีโพสิทอรีสำหรับระบบสั่งซื้อออนไลน์
type WosPgRepo struct {
	repository.PgRepo[models.WosOrder]
	DB *gorm.DB
}

// CreateWosPgRepository creates a new WOS repository.
// สร้างรีโพสิทอรีสำหรับระบบสั่งซื้อออนไลน์
func CreateWosPgRepository(db *gorm.DB) wos.WosPgRepository {
	return &WosPgRepo{
		PgRepo: repository.CreatePgRepo[models.WosOrder](db),
		DB:     db,
	}
}

// GetByOrderNumber finds an order by its order number.
// ค้นหาคำสั่งซื้อตามหมายเลขคำสั่งซื้อ
func (r *WosPgRepo) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.WosOrder, error) {
	var order models.WosOrder
	if result := r.DB.WithContext(ctx).Where("order_number = ?", orderNumber).First(&order); result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

// UpdateStatus updates the status of an order by ID.
// อัปเดตสถานะคำสั่งซื้อ
func (r *WosPgRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.DB.WithContext(ctx).Model(&models.WosOrder{}).Where("id = ?", id).Update("status", status).Error
}
