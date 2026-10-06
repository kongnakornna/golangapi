package models

import (
	"time"

	"github.com/google/uuid"
)

// PdpaConsent แทนบันทึกความยินยอมของผู้ใช้
// PdpaConsent represents a user's consent record.
type PdpaConsent struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index:idx_pdpa_consent_user_id" json:"user_id"`
	PurposeCode   string     `gorm:"type:varchar(64);not null;index:idx_pdpa_consent_purpose" json:"purpose_code"`
	Status        int8       `gorm:"not null;default:1;index:idx_pdpa_consent_status" json:"status"` // 0=revoked, 1=granted, 2=expired, 3=deleted
	ConsentedAt   time.Time  `gorm:"not null;default:NOW()" json:"consented_at"`
	ExpiresAt     time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time `json:"auto_deleted_at,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PdpaConsent) TableName() string {
	return "pdpa_consents"
}
