package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// AccountTerminatedHandler — เหตุการณ์ account.terminated
type AccountTerminatedHandler struct {
	baseHandler
}

func NewAccountTerminatedHandler(auditRepo repository.AuditRepository, log logger.Logger) *AccountTerminatedHandler {
	return &AccountTerminatedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *AccountTerminatedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
