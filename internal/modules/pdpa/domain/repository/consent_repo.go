package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentRepository interface {
	Save(ctx context.Context, c *entity.ConsentLog) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error)
	FindLatestByUserAndPurpose(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (*entity.ConsentLog, error)
	Revoke(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) error
	FindExpiredRevoked(ctx context.Context, cutoff time.Time) ([]entity.ConsentLog, error)
	DeletePermanently(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	IsConsentActive(ctx context.Context, userID uuid.UUID, purpose valueobject.ConsentPurpose) (bool, error)
}
