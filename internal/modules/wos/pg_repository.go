package wos

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// WosPgRepository defines data access methods for web orders.
// ดึงข้อมูลคำสั่งซื้อจากฐานข้อมูล
type WosPgRepository interface {
	internal.PgRepository[models.WosOrder]
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.WosOrder, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
