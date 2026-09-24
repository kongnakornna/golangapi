package entity

import (
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DSARRequest struct {
	ID              uuid.UUID               `json:"id"`
	UserID          uuid.UUID               `json:"user_id"`
	RequestType     valueobject.DSARType    `json:"request_type"`
	Status          valueobject.DSARStatus  `json:"status"`
	RequestedAt     time.Time               `json:"requested_at"`
	CompletedAt     *time.Time              `json:"completed_at,omitempty"`
	DataPayload     []byte                  `json:"data_payload,omitempty"`
	RejectionReason string                  `json:"rejection_reason,omitempty"`
	OTPCode         string                  `json:"-"`
	OTPExpiredAt    *time.Time              `json:"-"`
	IPAddress       string                  `json:"ip_address,omitempty"`
	UserAgent       string                  `json:"user_agent,omitempty"`
}

func NewDSARRequest(userID uuid.UUID, requestType valueobject.DSARType, otpCode string, otpExpiry time.Time, ip string, userAgent string) *DSARRequest {
	return &DSARRequest{
		ID:           uuid.New(),
		UserID:       userID,
		RequestType:  requestType,
		Status:       valueobject.DSARStatusPending,
		RequestedAt:  time.Now(),
		OTPCode:      otpCode,
		OTPExpiredAt: &otpExpiry,
		IPAddress:    ip,
		UserAgent:    userAgent,
	}
}

func (d *DSARRequest) VerifyOTP(otpCode string) bool {
	return d.OTPCode == otpCode && d.OTPExpiredAt != nil && time.Now().Before(*d.OTPExpiredAt)
}

func (d *DSARRequest) MarkProcessing() error {
	if d.Status != valueobject.DSARStatusPending {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusProcessing
	return nil
}

func (d *DSARRequest) MarkCompleted(payload []byte) error {
	if d.Status == valueobject.DSARStatusCompleted {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	now := time.Now()
	d.Status = valueobject.DSARStatusCompleted
	d.CompletedAt = &now
	d.DataPayload = payload
	return nil
}

func (d *DSARRequest) MarkRejected(reason string) error {
	if d.Status == valueobject.DSARStatusCompleted {
		return domainerrors.ErrDSARAlreadyProcessed
	}
	d.Status = valueobject.DSARStatusRejected
	d.RejectionReason = reason
	return nil
}