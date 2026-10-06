package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type DLQHandler struct {
	baseHandler
}

func NewDLQHandler(auditRepo repository.AuditRepository, log logger.Logger) *DLQHandler {
	return &DLQHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *DLQHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	h.logger.Warnf("dlq handler: processing dead-letter event topic=%q payload=%s",
		topic, payloadToJSON(payload))
	if err := h.saveAudit(ctx, topic, payload); err != nil {
		h.logger.Errorf("dlq handler: failed to record audit trail for topic %q: %s",
			topic, err.Error())
	}
	return nil
}
