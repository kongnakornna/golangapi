package purchaseorder

import (
	"context"

	"icmongolang/internal"
	pomodels "icmongolang/internal/modules/purchaseorder/models"

	"github.com/google/uuid"
)

type PurchaseOrderUseCaseI interface {
	internal.UseCaseI[pomodels.PurchaseOrderHeader]
	CreateWithDetails(ctx context.Context, header *pomodels.PurchaseOrderHeader, details []*pomodels.PurchaseOrderDetail) (*pomodels.PurchaseOrderHeader, error)
	CreateFromQuotation(ctx context.Context, quotationID uuid.UUID, userID uuid.UUID, whitelabelID uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	Send(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	Confirm(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	Receive(ctx context.Context, id uuid.UUID, request *ReceiveRequest) (*pomodels.PurchaseOrderHeader, error)
	Cancel(ctx context.Context, id uuid.UUID, reason string, userID uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error)
	List(ctx context.Context, limit, offset int, supplierID *uuid.UUID, status string, startDate, endDate string) ([]*pomodels.PurchaseOrderHeader, int64, error)
	GetSuggestions(ctx context.Context, jobID uuid.UUID) ([]*SuggestionItem, error)
	GetStatusHistory(ctx context.Context, poID uuid.UUID) ([]*pomodels.PurchaseOrderStatusHistory, error)
	GetPDF(ctx context.Context, id uuid.UUID) ([]byte, error)
	CountByFilter(ctx context.Context, supplierID *uuid.UUID, status string, startDate, endDate string) (int64, error)
}

type ReceiveRequest struct {
	Items []ReceiveItem
}

type ReceiveItem struct {
	DetailID        uuid.UUID
	ReceivedQuantity int
}

type SuggestionItem struct {
	PartID         uuid.UUID
	PartName       string
	PartCode       string
	SuggestedQty   int
	CurrentStock   int
	UnitPrice      float64
	FromQuotation  bool
}
