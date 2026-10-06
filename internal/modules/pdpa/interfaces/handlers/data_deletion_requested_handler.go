package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DataDeletionRequestedHandler — เหตุการณ์ data.deletion_requested
type DataDeletionRequestedHandler struct {
	baseHandler
}

func NewDataDeletionRequestedHandler(auditRepo repository.AuditRepository, log logger.Logger) *DataDeletionRequestedHandler {
	return &DataDeletionRequestedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *DataDeletionRequestedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
