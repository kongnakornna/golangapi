package handlers

import (
	"context"
	"fmt"

	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ชื่อ topic ของเหตุการณ์ต่าง ๆ ในกระบวนการ DSAR
const (
	topicDSARProcessing  = "dsar.processing"
	topicDSAROTPVerified = "dsar.otp_verified"
	topicDSARRejected    = "dsar.rejected"
)

// AuditTrailHandler ตัวจัดการแบบกลุ่ม (aggregate) สำหรับเหตุการณ์ต่าง ๆ
// ในกระบวนการ DSAR ได้แก่ dsar.processing, dsar.otp_verified และ dsar.rejected
// โดย mapping ระหว่าง topic และ action ที่ใช้บันทึก audit trail
// ถูกกำหนดไว้ในการสร้างตัวจัดการ (constructor)
type AuditTrailHandler struct {
	baseHandler
	actions map[string]string
}

// NewAuditTrailHandler สร้าง AuditTrailHandler และเตรียมตาราง
// mapping ระหว่าง topic และ action สำหรับบันทึก audit trail
func NewAuditTrailHandler(auditRepo repository.AuditRepository, logger logger.Logger) *AuditTrailHandler {
	return &AuditTrailHandler{
		baseHandler: baseHandler{auditRepo: auditRepo, logger: logger},
		actions: map[string]string{
			topicDSARProcessing:  topicDSARProcessing,
			topicDSAROTPVerified: topicDSAROTPVerified,
			topicDSARRejected:    topicDSARRejected,
		},
	}
}

// Handle บันทึก audit trail สำหรับเหตุการณ์ DSAR ที่รองรับ
// คืนค่า error เมื่อ topic ไม่เป็นที่รู้จักในชุด mapping
func (h *AuditTrailHandler) Handle(ctx context.Context, topic string, payload map[string]interface{}) error {
	action, ok := h.actions[topic]
	if !ok {
		return fmt.Errorf("unsupported dsar topic %q", topic)
	}

	return h.saveAudit(ctx, action, payload)
}