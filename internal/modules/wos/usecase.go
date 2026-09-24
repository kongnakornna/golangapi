package wos

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

// WosUseCaseI defines business logic methods for the Web Order System.
// อินเทอร์เฟซธุรกิจสำหรับระบบสั่งซื้อออนไลน์
type WosUseCaseI interface {
	internal.UseCaseI[models.WosOrder]
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.WosOrder, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
