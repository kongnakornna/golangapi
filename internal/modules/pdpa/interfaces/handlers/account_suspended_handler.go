package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// AccountSuspendedHandler ตัวจัดการเหตุการณ์ account.suspended
// (บัญชีผู้ใช้ถูกระงับ) ทำหน้าที่บันทึก audit trail เท่านั้น
type AccountSuspendedHandler struct {
	baseHandler
}

// NewAccountSuspendedHandler สร้าง AccountSuspendedHandler
func NewAccountSuspendedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *AccountSuspendedHandler {
	return &AccountSuspendedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ account.suspended
func (h *AccountSuspendedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}