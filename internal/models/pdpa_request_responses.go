package models

import (
	"encoding/json"
	"time"
)

// PdpaRequestResponse maps migrations/20260908_pdpa_request_responses.sql
type PdpaRequestResponse struct {
	ID               string          `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	RequestID        string          `gorm:"column:request_id;type:uuid;not null;uniqueIndex:uq_pdpa_request_responses_request_type,priority:1"`
	UserID           string          `gorm:"column:user_id;type:uuid;not null;index:idx_pdpa_request_responses_user_id"`
	ResponseType     int64           `gorm:"column:response_type;type:bigint;not null;default:1;uniqueIndex:uq_pdpa_request_responses_request_type,priority:2"`
	FilePath         *string         `gorm:"column:file_path;type:varchar(500)"`
	FileHash         *string         `gorm:"column:file_hash;type:varchar(255)"`
	Payload          json.RawMessage `gorm:"type:jsonb"`
	Status           int64           `gorm:"type:bigint;not null;default:1"`
	DeliveredAt      *time.Time      `gorm:"column:delivered_at;type:timestamptz"`
	DeliveredChannel *string         `gorm:"column:delivered_channel;type:varchar(100)"`
	CreatedAt        time.Time       `gorm:"type:timestamptz;not null;default:now();index:idx_pdpa_request_responses_created_at,priority:1,sort:desc"`
	UpdatedAt        time.Time       `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
	DeletedAt        *time.Time      `gorm:"type:timestamptz;index"`

	Request *PdpaUserRequest `gorm:"foreignKey:RequestID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PdpaRequestResponse) TableName() string {
	return "pdpa_request_responses"
}