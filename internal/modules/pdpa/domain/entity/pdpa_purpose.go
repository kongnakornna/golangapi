package entity

import (
	"time"

	"github.com/google/uuid"
)

type PdpaPurpose struct {
	ID              uuid.UUID  `json:"id"`
	PurposeCode     string     `json:"purpose_code"`
	PurposeType     int64      `json:"purpose_type"`
	NameEn          string     `json:"name_en"`
	NameTh          string     `json:"name_th"`
	DescriptionEn   *string    `json:"description_en,omitempty"`
	DescriptionTh   *string    `json:"description_th,omitempty"`
	ConsentRequired bool       `json:"consent_required"`
	RetentionDays   *int32     `json:"retention_days,omitempty"`
	Status          int64      `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewPdpaPurpose(
	purposeCode string,
	purposeType int64,
	nameEn string,
	nameTh string,
	descriptionEn *string,
	descriptionTh *string,
	consentRequired bool,
	retentionDays *int32,
	status int64,
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