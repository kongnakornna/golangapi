package handlers

import (
	"context"
	"fmt"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

const (
	topicDSARProcessing  = "dsar.processing"
	topicDSAROTPVerified = "dsar.otp_verified"
	topicDSARRejected    = "dsar.rejected"
)

type AuditTrailHandler struct {
	baseHandler
	actions map[string]string
}

func NewAuditTrailHandler(auditRepo repository.AuditRepository, log logger.Logger) *AuditTrailHandler {
	return &AuditTrailHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: log},
		actions: map[string]string{
			topicDSARProcessing:  topicDSARProcessing,
			topicDSAROTPVerified: topicDSAROTPVerified,
			topicDSARRejected:    topicDSARRejected,
		},
	}
}

func (h *AuditTrailHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	action, ok := h.actions[topic]
	if !ok {
		return fmt.Errorf("unsupported dsar topic %q", topic)
	}
	return h.saveAudit(ctx, action, payload)
}
