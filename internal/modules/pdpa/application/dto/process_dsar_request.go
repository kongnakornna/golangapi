package dto

import "github.com/google/uuid"

// ProcessDSARRequest represents the request payload for marking a data subject access request as processing.
// ProcessDSARRequest แสดงข้อมูลคำขอสำหรับการทำเครื่องหมายคำร้องของเจ้าของข้อมูลส่วนบุคคลว่าอยู่ระหว่างดำเนินการ
type ProcessDSARRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}