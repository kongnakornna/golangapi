package dto

import "github.com/google/uuid"

// SuspendAccountRequest represents the request payload for suspending a user account.
// SuspendAccountRequest แสดงข้อมูลคำขอสำหรับการระงับบัญชีผู้ใช้ของผู้ควบคุมข้อมูลส่วนบุคคล
type SuspendAccountRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	RetentionYears int       `json:"retention_years" validate:"required,gte=1" example:"2"`
}