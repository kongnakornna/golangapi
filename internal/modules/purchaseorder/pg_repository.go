package purchaseorder

import (
	"context"

	"icmongolang/internal"
	pomodels "icmongolang/internal/modules/purchaseorder/models"

	"github.com/google/uuid"
)

type PurchaseOrderPgRepository interface {
	internal.PgRepository[pomodels.PurchaseOrderHeader]
	CreateWithDetails(ctx context.Context, header *pomodels.PurchaseOrderHeader, details []*pomodels.PurchaseOrderDetail) (*pomodels.PurchaseOrderHeader, error)
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	List(ctx context.Context, limit, offset int, supplierID *uuid.UUID, status string, startDate, endDate string) ([]*pomodels.PurchaseOrderHeader, error)
	CountByFilter(ctx context.Context, supplierID *uuid.UUID, status string, startDate, endDate string) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, fields map[string]interface{}) error
	SaveStatusHistory(ctx context.Context, history *pomodels.PurchaseOrderStatusHistory) error
	GetStatusHistory(ctx context.Context, poID uuid.UUID) ([]*pomodels.PurchaseOrderStatusHistory, error)
	GetDetailsByHeaderID(ctx context.Context, headerID uuid.UUID) ([]*pomodels.PurchaseOrderDetail, error)
	UpdateDetail(ctx context.Context, detail *pomodels.PurchaseOrderDetail) error
	UpdateDetails(ctx context.Context, details []*pomodels.PurchaseOrderDetail) error
	IncrementReceivedQuantity(ctx context.Context, detailID uuid.UUID, quantity int) error
}
