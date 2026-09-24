package dto

import (
	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// SubmitDSARRequest represents the request payload for submitting a data subject access request.
// SubmitDSARRequest แสดงข้อมูลคำขอสำหรับการยื่นคำร้องขอใช้สิทธิของเจ้าของข้อมูลส่วนบุคคล
type SubmitDSARRequest struct {
	UserID      uuid.UUID            `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	RequestType valueobject.DSARType `json:"request_type" validate:"required,gte=1" example:"1"`
}