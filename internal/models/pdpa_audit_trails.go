package models

import (
	"encoding/json"
	"time"
)

// PdpaAuditTrail maps migrations/20260907_pdpa_audit_trails.sql
type PdpaAuditTrail struct {
	ID         string          `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID     *string         `gorm:"column:user_id;type:uuid;index:idx_pdpa_audit_trails_user_id"`
	ActorID    *string         `gorm:"column:actor_id;type:uuid;index:idx_pdpa_audit_trails_actor_id"`
	Action     string          `gorm:"type:varchar(100);not null"`
	EntityType string          `gorm:"column:entity_type;type:varchar(100);not null;index:idx_pdpa_audit_trails_entity,priority:1"`
	EntityID   *string         `gorm:"column:entity_id;type:uuid;index:idx_pdpa_audit_trails_entity,priority:2"`
	DataBefore json.RawMessage `gorm:"column:data_before;type:jsonb"`
	DataAfter  json.RawMessage `gorm:"column:data_after;type:jsonb"`
	IPAddress  *string         `gorm:"column:ip_address;type:varchar(45)"`
	Source     *string         `gorm:"type:varchar(100)"`
	CreatedAt  time.Time       `gorm:"type:timestamptz;not null;default:now();index:idx_pdpa_audit_trails_created_at,priority:1,sort:desc"`
}

func (PdpaAuditTrail) TableName() string {
	return "pdpa_audit_trails"
}