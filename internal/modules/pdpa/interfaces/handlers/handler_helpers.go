// Package handlers จัดเตรียมตัวจัดการ (handlers) ที่ทำหน้าที่บันทึก audit trail
// ของเหตุการณ์ (events) ที่เผยแพร่จากโมดูล PDPA
//
// Package handlers provides the event handlers that record audit trails for
// events published by the PDPA module.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// HandlerFunc คือสัญญาณ (signature) ของตัวจัดการเหตุการณ์แต่ละประเภท
// ซึ่งจะรับ topic และ payload ของเหตุการณ์จาก consuming layer ในชั้นถัดไป
// (Todo ⑦ / ⑨) เพื่อบันทึก audit trail เพียงอย่างเดียว
// โดยไม่มีการ re-invoke การ publish ย้อนกลับ (ป้องกัน infinite loop)
type HandlerFunc func(ctx context.Context, topic string, payload map[string]interface{}) error

// baseHandler เป็นโครงสร้างร่วมของตัวจัดการทุกตัว
// เพื่อลดโค้ดซ้ำซ้อนในการบันทึก audit trail
type baseHandler struct {
	auditRepo repository.AuditRepository
	logger    logger.Logger
}

// saveAudit สร้างและบันทึก audit trail หนึ่งรายการจาก action (topic)
// และ payload ของเหตุการณ์ จากนั้น wrap error ให้มีบริบทที่ชัดเจน
func (h *baseHandler) saveAudit(ctx context.Context, action string, payload map[string]interface{}) error {
	userID := userIDFromPayload(payload)

	audit := entity.NewAuditTrail(userID, action, payload)
	if err := h.auditRepo.Save(ctx, audit); err != nil {
		return fmt.Errorf("save audit trail for action %q: %w", action, err)
	}

	return nil
}

// userIDFromPayload คืนค่า *uuid.UUID ของผู้ใช้จากฟิลด์ user_id ใน payload
// คืนค่า nil เมื่อฟิลด์ไม่พบ หรือไม่ใช่ UUID ที่ถูกต้อง
func userIDFromPayload(payload map[string]interface{}) *uuid.UUID {
	val, ok := payload["user_id"]
	if !ok {
		return nil
	}

	str, ok := val.(string)
	if !ok {
		return nil
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return nil
	}

	return &id
}

// payloadToJSON แปลง payload เป็น JSON string สำหรับการ log อย่างปลอดภัย
// เมื่อแปลงไม่สำเร็จจะคืนค่าข้อความสั้น ๆ แทนข้อมูลดิบที่อาจอ่านยากหรือเป็นความลับ
func payloadToJSON(payload map[string]interface{}) string {
	b, err := json.Marshal(payload)
	if err != nil {
		return "<unmarshalable payload>"
	}

	return string(b)
}