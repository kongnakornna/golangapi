package models

import "time"

// PdpaPurpose maps migrations/20260903_pdpa_purposes.sql
type PdpaPurpose struct {
	ID              string     `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	PurposeCode     string     `gorm:"column:purpose_code;type:varchar(100);not null;uniqueIndex:uq_pdpa_purposes_purpose_code"`
	PurposeType     int64      `gorm:"column:purpose_type;type:bigint;not null;default:1"`
	NameEn          string     `gorm:"column:name_en;type:varchar(255);not null"`
	NameTh          string     `gorm:"column:name_th;type:varchar(255);not null"`
	DescriptionEn   *string    `gorm:"column:description_en;type:text"`
	DescriptionTh   *string    `gorm:"column:description_th;type:text"`
	ConsentRequired bool       `gorm:"column:consent_required;type:boolean;not null;default:true"`
	RetentionDays   *int32     `gorm:"column:retention_days;type:int"`
	Status          int64      `gorm:"type:bigint;not null;default:1;index:idx_pdpa_purposes_status"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now();index:idx_pdpa_purposes_created_at,priority:1,sort:desc"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"type:timestamptz;index"`
}

func (PdpaPurpose) TableName() string {
	return "pdpa_purposes"
}