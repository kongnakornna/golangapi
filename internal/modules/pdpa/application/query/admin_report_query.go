package query

import (
	"context"
	"time"

	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type GetAdminReportQuery struct {
	auditRepo domainrepo.AuditRepository
	log       logger.Logger
}

func NewGetAdminReportQuery(a domainrepo.AuditRepository, l logger.Logger) *GetAdminReportQuery {
	return &GetAdminReportQuery{auditRepo: a, log: l}
}

func (q *GetAdminReportQuery) Execute(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	return q.auditRepo.CountByAction(ctx, from, to)
}
