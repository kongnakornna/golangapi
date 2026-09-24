package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ConsentGrantedHandler ตัวจัดการเหตุการณ์ consent.granted
// (ผู้ใช้ให้ความยินยอม) ทำหน้าที่บันทึก audit trail เท่านั้น
type ConsentGrantedHandler struct {
	baseHandler
}

// NewConsentGrantedHandler สร้าง ConsentGrantedHandler
func NewConsentGrantedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *ConsentGrantedHandler {
	return &ConsentGrantedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ consent.granted
func (h *ConsentGrantedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}