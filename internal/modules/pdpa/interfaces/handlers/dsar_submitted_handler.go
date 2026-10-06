package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DSARSubmittedHandler — เหตุการณ์ dsar.submitted
type DSARSubmittedHandler struct {
	baseHandler
}

func NewDSARSubmittedHandler(auditRepo repository.AuditRepository, log logger.Logger) *DSARSubmittedHandler {
	return &DSARSubmittedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *DSARSubmittedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
