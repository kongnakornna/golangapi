package handlers

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// DLQHandler ใช้จัดการเหตุการณ์ที่ถูกส่งเข้าสู่ dead-letter queue (DLQ)
// เช่น เหตุการณ์ที่ decode ไม่สำเร็จ หรือประมวลผลไม่สำเร็จซ้ำหลายครั้ง
//
// DLQHandler จะพยายามบันทึกเหตุการณ์ลง audit trail เพื่อให้ตรวจสอบย้อนหลังได้
// แต่จะไม่คืน error หรือ panic ออกจากตัวจัดการ
// เพื่อไม่ให้การจัดคิว (queue consuming) ติดขัดแม้การบันทึกจะล้มเหลว
type DLQHandler struct {
	baseHandler
}

// NewDLQHandler สร้าง DLQHandler
func NewDLQHandler(auditRepo repository.AuditRepository, logger logger.Logger) *DLQHandler {
	return &DLQHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
	}
}

// Handle บันทึกเหตุการณ์จาก DLQ ลง audit trail
// หากบันทึกไม่สำเร็จจะ log เฉพาะ และคืนค่า nil เสมอ
func (h *DLQHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	h.logger.Warnf("dlq handler: processing dead-letter event topic=%q payload=%s", topic, payloadToJSON(payload))

	if err := h.saveAudit(ctx, topic, payload); err != nil {
		h.logger.Errorf("dlq handler: failed to record audit trail for topic %q: %s", topic, err.Error())
	}

	return nil
}