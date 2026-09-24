package dto

import (
	"time"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// ConsentStatusResponse represents the current consent status for a user and purpose.
// ConsentStatusResponse แสดงสถานะความยินยอมปัจจุบันของผู้ใช้สำหรับวัตถุประสงค์หนึ่ง
type ConsentStatusResponse struct {
	UserID    uuid.UUID                    `json:"user_id"`
	Purpose   string                       `json:"purpose"`
	Status    valueobject.ConsentStatus    `json:"status"`
	IsActive  bool                         `json:"is_active"`
	GrantedAt time.Time                    `json:"granted_at"`
	ExpiresAt time.Time                    `json:"expires_at"`
}

// ConsentStatusListResponse represents a list of consent statuses.
// ConsentStatusListResponse แสดงรายการสถานะความยินยอม
type ConsentStatusListResponse struct {
	Consents []ConsentStatusResponse `json:"consents"`
	Total    int64                   `json:"total"`
}