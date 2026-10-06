package query

import (
	"context"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type ListConsentLogByUserQuery struct {
	repo domainrepo.ConsentRepository
	log  logger.Logger
}

func NewListConsentLogByUserQuery(r domainrepo.ConsentRepository, l logger.Logger) *ListConsentLogByUserQuery {
	return &ListConsentLogByUserQuery{repo: r, log: l}
}

func (q *ListConsentLogByUserQuery) Execute(ctx context.Context, userID uuid.UUID) ([]entity.ConsentLog, error) {
	return q.repo.FindByUserID(ctx, userID)
}
