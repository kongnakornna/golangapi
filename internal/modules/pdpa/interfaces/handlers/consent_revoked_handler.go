package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ConsentRevokedHandler ตัวจัดการเหตุการณ์ consent.revoked
// (ผู้ใช้ถอนความยินยอม) ทำหน้าที่บันทึก audit trail เท่านั้น
type ConsentRevokedHandler struct {
	baseHandler
}

// NewConsentRevokedHandler สร้าง ConsentRevokedHandler
func NewConsentRevokedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *ConsentRevokedHandler {
	return &ConsentRevokedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ consent.revoked
func (h *ConsentRevokedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}