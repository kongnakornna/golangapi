package dto

import "github.com/google/uuid"

// RejectDSARRequest represents the request payload for rejecting a data subject access request.
// RejectDSARRequest แสดงข้อมูลคำขอสำหรับการปฏิเสธคำร้องของเจ้าของข้อมูลส่วนบุคคล
type RejectDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Reason        string    `json:"reason" validate:"required,min=1,max=500" example:"ข้อมูลไม่เพียงพอ"`
}