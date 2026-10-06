package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DSARCompletedHandler — เหตุการณ์ dsar.completed
type DSARCompletedHandler struct {
	baseHandler
}

func NewDSARCompletedHandler(auditRepo repository.AuditRepository, log logger.Logger) *DSARCompletedHandler {
	return &DSARCompletedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *DSARCompletedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
