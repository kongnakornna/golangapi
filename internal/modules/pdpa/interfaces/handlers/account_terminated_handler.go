package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// AccountTerminatedHandler ตัวจัดการเหตุการณ์ account.terminated
// (บัญชีผู้ใช้ถูกยุติ / ข้อมูลเข้าสู่กระบวนการลบแบบรอเวลายึดตาม retention)
// ทำหน้าที่บันทึก audit trail เท่านั้น
type AccountTerminatedHandler struct {
	baseHandler
}

// NewAccountTerminatedHandler สร้าง AccountTerminatedHandler
func NewAccountTerminatedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *AccountTerminatedHandler {
	return &AccountTerminatedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ account.terminated
func (h *AccountTerminatedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}