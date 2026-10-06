// Package models เก็บ GORM models ของระบบ
package models

import (
	"time"

	"github.com/google/uuid"
)

// PdpaPolicy แทนเวอร์ชันนโยบายความเป็นส่วนตัว
// PdpaPolicy represents a privacy policy version.
type PdpaPolicy struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PolicyKey   string     `gorm:"type:varchar(64);not null;index:idx_pdpa_policy_key" json:"policy_key"`
	Version     string     `gorm:"type:varchar(32);not null" json:"version"`
	TitleEn     string     `gorm:"type:varchar(255);not null" json:"title_en"`
	TitleTh     string     `gorm:"type:varchar(255);not null" json:"title_th"`
	ContentEn   *string    `gorm:"type:text" json:"content_en,omitempty"`
	ContentTh   *string    `gorm:"type:text" json:"content_th,omitempty"`
	Status      int64      `gorm:"not null;default:0;index:idx_pdpa_policy_status" json:"status"` // 0=draft, 1=active, 2=archived
	EffectiveAt *time.Time `json:"effective_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName ระบุชื่อตารางในฐานข้อมูล
func (PdpaPolicy) TableName() string {
	return "pdpa_privacy_policies"
}
