package repository

import (
	"context"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type Cache interface {
	Set(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose, status valueobject.ConsentStatus) error
	Get(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (valueobject.ConsentStatus, error)
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error
}