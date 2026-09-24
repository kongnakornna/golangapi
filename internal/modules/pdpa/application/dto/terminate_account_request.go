package dto

import "github.com/google/uuid"

// TerminateAccountRequest represents the request payload for terminating a user account.
// TerminateAccountRequest แสดงข้อมูลคำขอสำหรับการยุติบัญชีผู้ใช้ของผู้ควบคุมข้อมูลส่วนบุคคล
type TerminateAccountRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}