package entity

import (
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DSARRequest struct {
	ID              uuid.UUID              `json:"id" gorm:"column:id;primaryKey"`
	UserID          uuid.UUID              `json:"user_id" gorm:"column:user_id;index"`
	RequestType     valueobject.DSARType   `json:"request_type" gorm:"column:request_type"`
	Status          valueobject.DSARStatus `json:"status" gorm:"column:status"`
	RequestedAt     time.Time              `json:"requested_at" gorm:"column:requested_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty" gorm:"column:completed_at"`
	DataPayload     []byte                 `json:"data_payload,omitempty" gorm:"column:data_payload"`
	RejectionReason string                 `json:"rejection_reason,omitempty" gorm:"column:rejection_reason"`
	OTPCode         string                 `json:"-" gorm:"column:otp_code"`
	OTPExpiredAt    *time.Time             `json:"-" gorm:"column:otp_expired_at"`
	IPAddress       string                 `json:"ip_address,omitempty" gorm:"column:ip_address"`
	UserAgent       string                 `json:"user_agent,omitempty" gorm:"column:user_agent"`
}

func (DSARRequest) TableName() string { return "pdpa_user_requests" }

func NewDSARRequest(userID uuid.UUID, requestType valueobject.DSARType, otpCode string, otpExpiry time.Time, ip, ua string) *DSARRequest {
	return &DSARRequest{
		ID:           uuid.New(),
		UserID:       userID,
		RequestType:  requestType,
		Status:       valueobject.DSARStatusPending,
		RequestedAt:  time.Now(),
		OTPCode:      otpCode,
		OTPExpiredAt: &otpExpiry,
		IPAddress:    ip,
		UserAgent:    ua,
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
