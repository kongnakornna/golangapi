// Package handlers จัดเตรียม event handlers ที่บันทึก audit trail
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

type HandlerFunc func(ctx context.Context, topic string, payload map[string]interface{}) error

type baseHandler struct {
	auditRepo repository.AuditRepository
	logger    logger.Logger
}

func (h *baseHandler) saveAudit(ctx context.Context, action string, payload map[string]interface{}) error {
	userID := userIDFromPayload(payload)
	audit := entity.NewAuditTrail(userID, action, payload)
	if err := h.auditRepo.Save(ctx, audit); err != nil {
		return fmt.Errorf("save audit trail for action %q: %w", action, err)
	}
	return nil
}

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

func payloadToJSON(payload map[string]interface{}) string {
	b, err := json.Marshal(payload)
	if err != nil {
		return "<unmarshalable payload>"
	}
	return string(b)
}
