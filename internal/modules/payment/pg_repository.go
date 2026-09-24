package payment

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type PaymentPgRepository interface {
	internal.PgRepository[models.Payment]
	GetByInvoiceId(ctx context.Context, invoiceId uuid.UUID) (*models.Payment, error)
	Search(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*models.Payment, error)
	GetByCustomerId(ctx context.Context, customerId uuid.UUID, limit, offset int) ([]*models.Payment, error)
	GetOutstandingByCustomerId(ctx context.Context, customerId uuid.UUID) ([]*models.OutstandingBalance, error)
	GetPaymentHistory(ctx context.Context, customerId uuid.UUID) ([]*models.PaymentHistory, error)
	Count(ctx context.Context) (int64, error)
	CountByFilter(ctx context.Context, filter map[string]interface{}) (int64, error)
	// Receipt
	GetReceipt(ctx context.Context, id uuid.UUID) (*models.Receipt, error)
	GetReceiptByPaymentId(ctx context.Context, paymentId uuid.UUID) (*models.Receipt, error)
	CancelReceipt(ctx context.Context, id uuid.UUID, reason string) error
}
