package dto

import "github.com/google/uuid"

// VerifyOTPRequest represents the request payload for verifying an OTP to confirm a data subject request.
// VerifyOTPRequest แสดงข้อมูลคำขอสำหรับการยืนยัน OTP เพื่อยืนยันคำร้องของเจ้าของข้อมูลส่วนบุคคล
type VerifyOTPRequest struct {
	UserID        uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	DSARRequestID uuid.UUID `json:"dsar_request_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	OTPCode       string    `json:"otp_code" validate:"required,min=6,max=8" example:"123456"`
}