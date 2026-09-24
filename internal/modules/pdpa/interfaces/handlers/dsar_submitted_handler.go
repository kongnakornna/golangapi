package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DSARSubmittedHandler ตัวจัดการเหตุการณ์ dsar.submitted
// (ผู้ใช้ยื่นคำขอใช้สิทธิ Data Subject Access Request)
// ทำหน้าที่บันทึก audit trail เท่านั้น
type DSARSubmittedHandler struct {
	baseHandler
}

// NewDSARSubmittedHandler สร้าง DSARSubmittedHandler
func NewDSARSubmittedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *DSARSubmittedHandler {
	return &DSARSubmittedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ dsar.submitted
func (h *DSARSubmittedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}