package query

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type GetDSARStatusQuery struct {
	repo domainrepo.DSARRepository
	log  logger.Logger
}

func NewGetDSARStatusQuery(r domainrepo.DSARRepository, l logger.Logger) *GetDSARStatusQuery {
	return &GetDSARStatusQuery{repo: r, log: l}
}

func (q *GetDSARStatusQuery) Execute(ctx context.Context, id uuid.UUID) (*entity.DSARRequest, error) {
	return q.repo.FindByID(ctx, id)
}

type ListDSARByUserQuery struct {
	repo domainrepo.DSARRepository
	log  logger.Logger
}

func NewListDSARByUserQuery(r domainrepo.DSARRepository, l logger.Logger) *ListDSARByUserQuery {
	return &ListDSARByUserQuery{repo: r, log: l}
}

func (q *ListDSARByUserQuery) Execute(ctx context.Context, userID uuid.UUID) ([]entity.DSARRequest, error) {
	return q.repo.FindByUserID(ctx, userID)
}
