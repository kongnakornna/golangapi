package customer

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type CustomerPgRepository interface {
	internal.PgRepository[models.Customer]
	GetMultiByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Customer, error)
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type CarPgRepository interface {
	internal.PgRepository[models.Car]
	GetMultiByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Car, error)
	CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error)
}
