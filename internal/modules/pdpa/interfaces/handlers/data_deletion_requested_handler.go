package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DataDeletionRequestedHandler ตัวจัดการเหตุการณ์ data.deletion_requested
// (ผู้ใช้ยืนยันการลบข้อมูลส่วนบุคคล) ทำหน้าที่บันทึก audit trail เท่านั้น
type DataDeletionRequestedHandler struct {
	baseHandler
}

// NewDataDeletionRequestedHandler สร้าง DataDeletionRequestedHandler
func NewDataDeletionRequestedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *DataDeletionRequestedHandler {
	return &DataDeletionRequestedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ data.deletion_requested
func (h *DataDeletionRequestedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}