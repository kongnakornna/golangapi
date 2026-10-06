package repository

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type PolicyRepository interface {
	Save(ctx context.Context, p *entity.PdpaPolicy) error
	Update(ctx context.Context, p *entity.PdpaPolicy) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PdpaPolicy, error)
	FindByVersion(ctx context.Context, policyKey, version string) (*entity.PdpaPolicy, error)
	FindActive(ctx context.Context) ([]entity.PdpaPolicy, error)
	FindAll(ctx context.Context) ([]entity.PdpaPolicy, error)
	DeactivateAll(ctx context.Context) error
}
