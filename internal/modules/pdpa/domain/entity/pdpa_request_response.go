package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PdpaRequestResponse struct {
	ID               uuid.UUID       `json:"id"`
	RequestID        uuid.UUID       `json:"request_id"`
	UserID           uuid.UUID       `json:"user_id"`
	ResponseType     int64           `json:"response_type"`
	FilePath         *string         `json:"file_path,omitempty"`
	FileHash         *string         `json:"file_hash,omitempty"`
	Payload          json.RawMessage `json:"payload,omitempty"`
	Status           int64           `json:"status"`
	DeliveredAt      *time.Time      `json:"delivered_at,omitempty"`
	DeliveredChannel *string         `json:"delivered_channel,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at,omitempty"`
}

func NewPdpaRequestResponse(
	requestID uuid.UUID,
	userID uuid.UUID,
	responseType int64,
	filePath *string,
	fileHash *string,
	payload json.RawMessage,
	status int64,
	deliveredAt *time.Time,
	deliveredChannel *string,
) *PdpaRequestResponse {
	now := time.Now()
	return &PdpaRequestResponse{
		ID:               uuid.New(),
		RequestID:        requestID,
		UserID:           userID,
		ResponseType:     responseType,
		FilePath:         filePath,
		FileHash:         fileHash,
		Payload:          payload,
		Status:           status,
		DeliveredAt:      deliveredAt,
		DeliveredChannel: deliveredChannel,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}