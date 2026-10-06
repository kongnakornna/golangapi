package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PdpaAuditTrail แทนบันทึกการกระทำที่ตรวจสอบได้
// PdpaAuditTrail represents an auditable action record.
type PdpaAuditTrail struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index:idx_pdpa_audit_user_id" json:"user_id,omitempty"`
	Action    string         `gorm:"type:varchar(64);not null;index:idx_pdpa_audit_action" json:"action"`
	Details   datatypes.JSON `gorm:"type:jsonb;index:idx_pdpa_audit_details,type:gin" json:"details,omitempty"`
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent string         `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index:idx_pdpa_audit_created_at" json:"created_at"`
}

func (PdpaAuditTrail) TableName() string {
	return "pdpa_audit_trails"
}
