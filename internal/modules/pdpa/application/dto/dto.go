package dto

import "time"

// ---------- Consent ----------

type RecordConsentRequest struct {
	PurposeID      *string `json:"purpose_id,omitempty"`
	PolicyID       *string `json:"policy_id,omitempty"`
	ConsentType    int64   `json:"consent_type"`
	Source         string  `json:"source,omitempty"`
	Channel        string  `json:"channel,omitempty"`
	ConsentVersion string  `json:"consent_version,omitempty"`
}

type RevokeConsentRequest struct {
	PurposeID *string `json:"purpose_id,omitempty"`
}

// ---------- DSAR ----------

type SubmitDSARRequest struct {
	RequestType    int64  `json:"request_type"`
	RequestDetail  string `json:"request_detail,omitempty"`
	RequestChannel string `json:"request_channel,omitempty"`
}

type ProcessDSARRequest struct {
	Status          int64  `json:"status"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	Payload         string `json:"payload,omitempty"`
}

// ---------- Account ----------

type SuspendAccountRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Reason string `json:"reason,omitempty"`
}

type TerminateAccountRequest struct {
	UserID         string `json:"user_id" validate:"required"`
	RetentionYears int    `json:"retention_years"`
}

type ConfirmDeletionRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// ---------- Policy ----------

type PublishPolicyRequest struct {
	PolicyKey   string     `json:"policy_key" validate:"required"`
	Version     string     `json:"version" validate:"required"`
	TitleEn     string     `json:"title_en" validate:"required"`
	TitleTh     string     `json:"title_th" validate:"required"`
	ContentEn   *string    `json:"content_en,omitempty"`
	ContentTh   *string    `json:"content_th,omitempty"`
	EffectiveAt *time.Time `json:"effective_at,omitempty"`
}

// ---------- Report ----------

type AdminReportResponse struct {
	TotalConsents   int64 `json:"total_consents"`
	TotalDSAR       int64 `json:"total_dsar"`
	PendingDSAR     int64 `json:"pending_dsar"`
	CompletedDSAR   int64 `json:"completed_dsar"`
	SuspendedUsers  int64 `json:"suspended_users"`
	TerminatedUsers int64 `json:"terminated_users"`
	AuditTrailCount int64 `json:"audit_trail_count"`
}
