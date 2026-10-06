package entity

import (
	"time"

	"github.com/google/uuid"
)

type PdpaPolicy struct {
	ID          uuid.UUID  `json:"id" gorm:"column:id;primaryKey"`
	PolicyKey   string     `json:"policy_key" gorm:"column:policy_key"`
	Version     string     `json:"version" gorm:"column:version"`
	TitleEn     string     `json:"title_en" gorm:"column:title_en"`
	TitleTh     string     `json:"title_th" gorm:"column:title_th"`
	ContentEn   *string    `json:"content_en,omitempty" gorm:"column:content_en"`
	ContentTh   *string    `json:"content_th,omitempty" gorm:"column:content_th"`
	Status      int64      `json:"status" gorm:"column:status"`
	EffectiveAt *time.Time `json:"effective_at,omitempty" gorm:"column:effective_at"`
	PublishedAt *time.Time `json:"published_at,omitempty" gorm:"column:published_at"`
	CreatedBy   *string    `json:"created_by,omitempty" gorm:"column:created_by"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (PdpaPolicy) TableName() string { return "pdpa_policies" }

func NewPdpaPolicy(
	policyKey, version, titleEn, titleTh string,
	contentEn, contentTh *string,
	status int64,
	effectiveAt, publishedAt *time.Time,
	createdBy *string,
) *PdpaPolicy {
	now := time.Now()
	return &PdpaPolicy{
		ID:          uuid.New(),
		PolicyKey:   policyKey,
		Version:     version,
		TitleEn:     titleEn,
		TitleTh:     titleTh,
		ContentEn:   contentEn,
		ContentTh:   contentTh,
		Status:      status,
		EffectiveAt: effectiveAt,
		PublishedAt: publishedAt,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
