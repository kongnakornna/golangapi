package dto

import "github.com/google/uuid"

// ConfirmDeletionRequest represents the request payload for confirming the deletion of a user's personal data.
// ConfirmDeletionRequest แสดงข้อมูลคำขอสำหรับการยืนยันการลบข้อมูลส่วนบุคคลของผู้ใช้
type ConfirmDeletionRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}