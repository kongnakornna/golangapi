package models

import (
	"time"

	"github.com/google/uuid"
)

// PdpaPurpose แทนวัตถุประสงค์การประมวลผลข้อมูลส่วนบุคคล
// PdpaPurpose represents a lawful purpose for processing personal data.
type PdpaPurpose struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PurposeCode     string     `gorm:"type:varchar(64);not null;uniqueIndex:uq_pdpa_purpose_code" json:"purpose_code"`
	PurposeType     int64      `gorm:"not null;default:1" json:"purpose_type"`
	NameEn          string     `gorm:"type:varchar(255);not null" json:"name_en"`
	NameTh          string     `gorm:"type:varchar(255);not null" json:"name_th"`
	DescriptionEn   *string    `gorm:"type:text" json:"description_en,omitempty"`
	DescriptionTh   *string    `gorm:"type:text" json:"description_th,omitempty"`
	ConsentRequired bool       `gorm:"not null;default:true" json:"consent_required"`
	RetentionDays   *int32     `json:"retention_days,omitempty"`
	Status          int64      `gorm:"not null;default:1;index:idx_pdpa_purpose_status" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (PdpaPurpose) TableName() string {
	return "pdpa_purposes"
}
