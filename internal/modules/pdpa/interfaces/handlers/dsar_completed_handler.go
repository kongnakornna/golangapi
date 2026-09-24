package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DSARCompletedHandler ตัวจัดการเหตุการณ์ dsar.completed
// (ผู้ให้บริการดำเนินการตามคำขอใช้สิทธิเสร็จสิ้น)
// ทำหน้าที่บันทึก audit trail เท่านั้น
//
// หมายเหตุ: ฟิลด์ tx_hash ใน payload อาจเป็นค่าว่าง (ไม่ใช่ข้อมูลบังคับ)
// ซึ่งไม่ส่งผลต่อการบันทึก audit trail เนื่องจากตัวจัดการนี้
// ไม่ได้อ่านค่าในฟิลด์ดังกล่าว
type DSARCompletedHandler struct {
	baseHandler
}

// NewDSARCompletedHandler สร้าง DSARCompletedHandler
func NewDSARCompletedHandler(auditRepo repository.AuditRepository, logger logger.Logger) *DSARCompletedHandler {
	return &DSARCompletedHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึก audit trail ของเหตุการณ์ dsar.completed
func (h *DSARCompletedHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	return h.saveAudit(ctx, topic, payload)
}