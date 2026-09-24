package dto

import "github.com/google/uuid"

// CompleteDSARRequest represents the request payload for completing a data subject access request.
// CompleteDSARRequest แสดงข้อมูลคำขอสำหรับการดำเนินการคำร้องของเจ้าของข้อมูลส่วนบุคคลให้เสร็จสิ้น
type CompleteDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ResponseData  []byte    `json:"response_data" validate:"required" example:"eyJkYXRhIjoiZXhhbXBsZSJ9"`
}