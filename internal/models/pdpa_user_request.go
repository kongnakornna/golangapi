package models

import (
	"time"

	"github.com/google/uuid"
)

// PdpaUserRequest แทนคำร้องขอใช้สิทธิ์ของเจ้าของข้อมูล (DSAR)
// PdpaUserRequest represents a Data Subject Access Request (DSAR).
type PdpaUserRequest struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index:idx_pdpa_user_request_user_id" json:"user_id"`
	RequestType     int8       `gorm:"not null" json:"request_type"`                                        // 1=ACCESS, 3=ERASURE, 7=WITHDRAW_CONSENT
	Status          int8       `gorm:"not null;default:1;index:idx_pdpa_user_request_status" json:"status"` // 1=PENDING, 3=PROCESSING, 4=COMPLETED, 5=REJECTED
	RequestedAt     time.Time  `gorm:"not null;default:NOW()" json:"requested_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	DataPayload     []byte     `gorm:"type:bytea" json:"data_payload,omitempty"`
	RejectionReason string     `gorm:"type:text" json:"rejection_reason,omitempty"`
	OTPCode         string     `gorm:"type:varchar(16)" json:"-"`
	OTPExpiredAt    *time.Time `json:"-"`
	IPAddress       string     `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent       string     `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PdpaUserRequest) TableName() string {
	return "pdpa_dsar_requests"
}
