package query

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type GetActivePrivacyPolicyQuery struct {
	repo domainrepo.PolicyRepository
	log  logger.Logger
}

func NewGetActivePrivacyPolicyQuery(r domainrepo.PolicyRepository, l logger.Logger) *GetActivePrivacyPolicyQuery {
	return &GetActivePrivacyPolicyQuery{repo: r, log: l}
}

func (q *GetActivePrivacyPolicyQuery) Execute(ctx context.Context) ([]entity.PdpaPolicy, error) {
	return q.repo.FindActive(ctx)
}
