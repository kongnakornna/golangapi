package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// AccountSuspendedHandler — เหตุการณ์ account.suspended
type AccountSuspendedHandler struct {
	baseHandler
}

func NewAccountSuspendedHandler(auditRepo repository.AuditRepository, log logger.Logger) *AccountSuspendedHandler {
	return &AccountSuspendedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *AccountSuspendedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
