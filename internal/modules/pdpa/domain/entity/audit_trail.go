package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditTrail — maps pdpa_audit_trails
type AuditTrail struct {
	ID        uuid.UUID   `json:"id" gorm:"column:id;primaryKey"`
	UserID    *uuid.UUID  `json:"user_id,omitempty" gorm:"column:user_id;index"`
	Action    string      `json:"action" gorm:"column:action"`
	Details   interface{} `json:"details,omitempty" gorm:"-"`
	IPAddress string      `json:"ip_address,omitempty" gorm:"column:ip_address"`
	UserAgent string      `json:"user_agent,omitempty" gorm:"column:user_agent"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at"`
}

func (AuditTrail) TableName() string { return "pdpa_audit_trails" }

func NewAuditTrail(userID *uuid.UUID, action string, details interface{}) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}
}

func NewAuditTrailWithMeta(userID *uuid.UUID, action string, details interface{}, ip, ua string) *AuditTrail {
	at := NewAuditTrail(userID, action, details)
	at.IPAddress = ip
	at.UserAgent = ua
	return at
}
