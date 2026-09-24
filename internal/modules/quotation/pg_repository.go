package quotation

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type QuotationPgRepository interface {
	internal.PgRepository[models.Quotation]
	Count(ctx context.Context) (int64, error)
	GetQuotationParts(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationPart, error)
	GetQuotationServices(ctx context.Context, quotationID uuid.UUID) ([]models.QuotationService, error)
}
