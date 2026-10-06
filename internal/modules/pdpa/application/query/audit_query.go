package query

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type ListAuditTrailsQuery struct {
	repo domainrepo.AuditRepository
	log  logger.Logger
}

func NewListAuditTrailsQuery(r domainrepo.AuditRepository, l logger.Logger) *ListAuditTrailsQuery {
	return &ListAuditTrailsQuery{repo: r, log: l}
}

func (q *ListAuditTrailsQuery) Execute(ctx context.Context, filter domainrepo.AuditFilter) ([]*entity.AuditTrail, error) {
	return q.repo.FindByFilter(ctx, filter)
}
