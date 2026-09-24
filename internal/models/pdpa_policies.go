package models

import "time"

// PdpaPolicy maps migrations/20260902_pdpa_policies.sql
type PdpaPolicy struct {
	ID          string     `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	PolicyKey   string     `gorm:"column:policy_key;type:varchar(100);not null;uniqueIndex:uq_pdpa_policies_policy_key_version"`
	Version     string     `gorm:"type:varchar(50);not null;uniqueIndex:uq_pdpa_policies_policy_key_version"`
	TitleEn     string     `gorm:"column:title_en;type:varchar(255);not null"`
	TitleTh     string     `gorm:"column:title_th;type:varchar(255);not null"`
	ContentEn   *string    `gorm:"column:content_en;type:text"`
	ContentTh   *string    `gorm:"column:content_th;type:text"`
	Status      int64      `gorm:"type:bigint;not null;default:1;index:idx_pdpa_policies_status"`
	EffectiveAt *time.Time `gorm:"column:effective_at;type:timestamptz;index"`
	PublishedAt *time.Time `gorm:"column:published_at;type:timestamptz"`
	CreatedBy   *string    `gorm:"column:created_by;type:uuid"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null;default:now();index:idx_pdpa_policies_created_at,priority:1,sort:desc"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"type:timestamptz;index"`
}

func (PdpaPolicy) TableName() string {
	return "pdpa_policies"
}