package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditTrail struct {
	ID        uuid.UUID      `json:"id"`
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	Action    string         `json:"action"`
	Details   interface{}    `json:"details,omitempty"`
	IPAddress string         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

func (AuditTrail) TableName() string {
	return "pdpa_audit_trails"
}

func NewAuditTrail(userID *uuid.UUID, action string, details interface{}) *AuditTrail {
	return &AuditTrail{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}
}

func NewAuditTrailWithMeta(userID *uuid.UUID, action string, details interface{}, ipAddress string, userAgent string) *AuditTrail {
	at := NewAuditTrail(userID, action, details)
	at.IPAddress = ipAddress
	at.UserAgent = userAgent
	return at
}