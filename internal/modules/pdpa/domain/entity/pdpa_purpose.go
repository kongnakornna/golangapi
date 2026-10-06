package entity

import (
	"time"

	"github.com/google/uuid"
)

type PdpaPurpose struct {
	ID              uuid.UUID  `json:"id" gorm:"column:id;primaryKey"`
	PurposeCode     string     `json:"purpose_code" gorm:"column:purpose_code"`
	PurposeType     int64      `json:"purpose_type" gorm:"column:purpose_type"`
	NameEn          string     `json:"name_en" gorm:"column:name_en"`
	NameTh          string     `json:"name_th" gorm:"column:name_th"`
	DescriptionEn   *string    `json:"description_en,omitempty" gorm:"column:description_en"`
	DescriptionTh   *string    `json:"description_th,omitempty" gorm:"column:description_th"`
	ConsentRequired bool       `json:"consent_required" gorm:"column:consent_required"`
	RetentionDays   *int32     `json:"retention_days,omitempty" gorm:"column:retention_days"`
	Status          int64      `json:"status" gorm:"column:status"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (PdpaPurpose) TableName() string { return "pdpa_purposes" }

func NewPdpaPurpose(
	purposeCode string, purposeType int64,
	nameEn, nameTh string,
	descriptionEn, descriptionTh *string,
	consentRequired bool, retentionDays *int32, status int64,
) *PdpaPurpose {
	now := time.Now()
	return &PdpaPurpose{
		ID:              uuid.New(),
		PurposeCode:     purposeCode,
		PurposeType:     purposeType,
		NameEn:          nameEn,
		NameTh:          nameTh,
		DescriptionEn:   descriptionEn,
		DescriptionTh:   descriptionTh,
		ConsentRequired: consentRequired,
		RetentionDays:   retentionDays,
		Status:          status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
