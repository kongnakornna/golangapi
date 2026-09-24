package dto

import (
	"time"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// DSARStatusResponse represents the status of a data subject access request.
// DSARStatusResponse แสดงสถานะของคำขอสิทธิตามพรบ.คุ้มครองข้อมูลส่วนบุคคล
type DSARStatusResponse struct {
	ID              uuid.UUID             `json:"id"`
	UserID          uuid.UUID             `json:"user_id"`
	RequestType     valueobject.DSARType  `json:"request_type"`
	Status          valueobject.DSARStatus `json:"status"`
	RequestedAt     time.Time             `json:"requested_at"`
	CompletedAt     *time.Time            `json:"completed_at,omitempty"`
	RejectionReason string                `json:"rejection_reason,omitempty"`
}

// DSARStatusListResponse represents a list of DSAR statuses.
// DSARStatusListResponse แสดงรายการสถานะของคำขอสิทธิตามพรบ.คุ้มครองข้อมูลส่วนบุคคล
type DSARStatusListResponse struct {
	Requests []DSARStatusResponse `json:"requests"`
	Total    int64                `json:"total"`
}