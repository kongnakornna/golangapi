package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PdpaRequestResponse แทนการตอบกลับคำร้อง DSAR
// PdpaRequestResponse represents the response to a DSAR request.
type PdpaRequestResponse struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RequestID        uuid.UUID       `gorm:"type:uuid;not null;index:idx_pdpa_response_request_id" json:"request_id"`
	UserID           uuid.UUID       `gorm:"type:uuid;not null;index:idx_pdpa_response_user_id" json:"user_id"`
	ResponseType     int64           `gorm:"not null" json:"response_type"`
	FilePath         *string         `gorm:"type:text" json:"file_path,omitempty"`
	FileHash         *string         `gorm:"type:varchar(128)" json:"file_hash,omitempty"`
	Payload          json.RawMessage `gorm:"type:jsonb" json:"payload,omitempty"`
	Status           int64           `gorm:"not null;default:1" json:"status"`
	DeliveredAt      *time.Time      `json:"delivered_at,omitempty"`
	DeliveredChannel *string         `gorm:"type:varchar(32)" json:"delivered_channel,omitempty"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        *time.Time      `gorm:"index" json:"deleted_at,omitempty"`
}

func (PdpaRequestResponse) TableName() string {
	return "pdpa_request_responses"
}
