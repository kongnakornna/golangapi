package entity

import (
	"time"

	"github.com/google/uuid"
)

type PdpaPolicy struct {
	ID          uuid.UUID  `json:"id"`
	PolicyKey   string     `json:"policy_key"`
	Version     string     `json:"version"`
	TitleEn     string     `json:"title_en"`
	TitleTh     string     `json:"title_th"`
	ContentEn   *string    `json:"content_en,omitempty"`
	ContentTh   *string    `json:"content_th,omitempty"`
	Status      int64      `json:"status"`
	EffectiveAt *time.Time `json:"effective_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedBy   *string    `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewPdpaPolicy(
	policyKey string,
	version string,
	titleEn string,
	titleTh string,
	contentEn *string,
	contentTh *string,
	status int64,
	effectiveAt *time.Time,
	publishedAt *time.Time,
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