package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PdpaRequestResponse struct {
	ID               uuid.UUID       `json:"id" gorm:"column:id;primaryKey"`
	RequestID        uuid.UUID       `json:"request_id" gorm:"column:request_id"`
	UserID           uuid.UUID       `json:"user_id" gorm:"column:user_id"`
	ResponseType     int64           `json:"response_type" gorm:"column:response_type"`
	FilePath         *string         `json:"file_path,omitempty" gorm:"column:file_path"`
	FileHash         *string         `json:"file_hash,omitempty" gorm:"column:file_hash"`
	Payload          json.RawMessage `json:"payload,omitempty" gorm:"column:payload"`
	Status           int64           `json:"status" gorm:"column:status"`
	DeliveredAt      *time.Time      `json:"delivered_at,omitempty" gorm:"column:delivered_at"`
	DeliveredChannel *string         `json:"delivered_channel,omitempty" gorm:"column:delivered_channel"`
	CreatedAt        time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (PdpaRequestResponse) TableName() string { return "pdpa_request_responses" }
