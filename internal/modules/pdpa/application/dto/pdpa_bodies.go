// Package dto — request body types สำหรับ PDPA HTTP endpoints
package dto

import "time"

// ============================================================
// Consent
// ============================================================

// RecordConsentBody — body ของ POST /v1/pdpa/consent
type RecordConsentBody struct {
	Purpose int8   `json:"purpose" example:"1"`  // valueobject.ConsentPurpose: 1=NECESSARY, 5=MARKETING, 6=ANALYTICS
	Source  string `json:"source" example:"web"` // web | mobile | api | offline
}

// RevokeConsentBody — body ของ DELETE /v1/pdpa/consent
type RevokeConsentBody struct {
	Purpose int8 `json:"purpose" example:"5"`
}

// ============================================================
// DSAR
// ============================================================

// SubmitDSARBody — body ของ POST /v1/pdpa/dsar
type SubmitDSARBody struct {
	RequestType int8 `json:"request_type" example:"1"` // valueobject.DSARType: 1=ACCESS, 3=ERASURE, 7=WITHDRAW_CONSENT
}

// ProcessDSARBody — body ของ POST /v1/pdpa/dsar/{id}/process (admin)
type ProcessDSARBody struct {
	Status          int8   `json:"status" example:"3"` // valueobject.DSARStatus: 3=PROCESSING, 4=COMPLETED, 5=REJECTED
	RejectionReason string `json:"rejection_reason,omitempty"`
	Payload         string `json:"payload,omitempty"`
}

// ============================================================
// Account
// ============================================================

// SuspendAccountBody — body ของ POST /v1/pdpa/account/suspend (admin)
type SuspendAccountBody struct {
	UserID         string `json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	RetentionYears int    `json:"retention_years" example:"5"`
	Reason         string `json:"reason,omitempty" example:"PDPA request"`
}

// TerminateAccountBody — body ของ POST /v1/pdpa/account/terminate (admin)
type TerminateAccountBody struct {
	UserID string `json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

// ConfirmDeletionBody — body ของ POST /v1/pdpa/account/confirm-deletion (admin)
type ConfirmDeletionBody struct {
	UserID string `json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

// ============================================================
// Policy
// ============================================================

// PublishPolicyBody — body ของ POST /v1/pdpa/policy (admin)
type PublishPolicyBody struct {
	PolicyKey   string     `json:"policy_key" example:"privacy" validate:"required"`
	Version     string     `json:"version" example:"1.0" validate:"required"`
	TitleEn     string     `json:"title_en" example:"Privacy Policy" validate:"required"`
	TitleTh     string     `json:"title_th" example:"นโยบายความเป็นส่วนตัว" validate:"required"`
	ContentEn   *string    `json:"content_en,omitempty"`
	ContentTh   *string    `json:"content_th,omitempty"`
	EffectiveAt *time.Time `json:"effective_at,omitempty" example:"2026-01-01T00:00:00Z"`
}
