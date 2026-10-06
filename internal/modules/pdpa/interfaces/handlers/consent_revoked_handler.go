package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ConsentRevokedHandler — เหตุการณ์ consent.revoked
type ConsentRevokedHandler struct {
	baseHandler
}

func NewConsentRevokedHandler(auditRepo repository.AuditRepository, log logger.Logger) *ConsentRevokedHandler {
	return &ConsentRevokedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *ConsentRevokedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
