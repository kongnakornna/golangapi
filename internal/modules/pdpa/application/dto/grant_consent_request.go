package dto

import (
	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// GrantConsentRequest represents the request payload for granting a data subject's consent.
// GrantConsentRequest แสดงข้อมูลคำขอสำหรับการให้ความยินยอมของเจ้าของข้อมูลส่วนบุคคล
type GrantConsentRequest struct {
	UserID  uuid.UUID                  `json:"user_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Purpose valueobject.ConsentPurpose `json:"purpose" validate:"required,gte=1" example:"1"`
}