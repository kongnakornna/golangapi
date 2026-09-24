package models

import "time"

// PdpaConsent maps migrations/20260904_pdpa_consents.sql
type PdpaConsent struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID         string     `gorm:"column:user_id;type:uuid;not null;index:idx_pdpa_consents_user_id"`
	PurposeID      *string    `gorm:"column:purpose_id;type:uuid;index"`
	PolicyID       *string    `gorm:"column:policy_id;type:uuid;index"`
	ConsentType    int64      `gorm:"column:consent_type;type:bigint;not null;default:1"`
	Status         int64      `gorm:"type:bigint;not null;default:1;index:idx_pdpa_consents_status"`
	Source         *string    `gorm:"type:varchar(100)"`
	Channel        *string    `gorm:"type:varchar(100)"`
	ConsentVersion *string    `gorm:"column:consent_version;type:varchar(50)"`
	IPAddress      *string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent      *string    `gorm:"column:user_agent;type:text"`
	ConsentedAt    time.Time  `gorm:"column:consented_at;type:timestamptz;not null;default:now();autoCreateTime"`
	RevokedAt      *time.Time `gorm:"column:revoked_at;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now();autoCreateTime;index:idx_pdpa_consents_created_at,priority:1,sort:desc"`
	DeletedAt      *time.Time `gorm:"type:timestamptz;index"`

	Purpose *PdpaPurpose `gorm:"foreignKey:PurposeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Policy  *PdpaPolicy  `gorm:"foreignKey:PolicyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (PdpaConsent) TableName() string {
	return "pdpa_consents"
}