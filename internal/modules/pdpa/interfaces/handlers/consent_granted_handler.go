package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ConsentGrantedHandler — เหตุการณ์ consent.granted
type ConsentGrantedHandler struct {
	baseHandler
}

func NewConsentGrantedHandler(auditRepo repository.AuditRepository, log logger.Logger) *ConsentGrantedHandler {
	return &ConsentGrantedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
	}
}

func (h *ConsentGrantedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}
