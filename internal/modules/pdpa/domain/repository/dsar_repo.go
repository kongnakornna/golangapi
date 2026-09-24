package repository

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DSARRepository interface {
	Save(ctx context.Context, dsar *entity.DSARRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status valueobject.DSARStatus) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}