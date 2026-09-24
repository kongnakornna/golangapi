package entity

import (
	"time"

	"github.com/google/uuid"

	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type ConsentLog struct {
	ID            uuid.UUID                 `json:"id"`
	UserID        uuid.UUID                 `json:"user_id"`
	Purpose       string                    `json:"purpose"`
	Status        valueobject.ConsentStatus `json:"status"`
	ConsentedAt   time.Time                 `json:"consented_at"`
	ExpiresAt     time.Time                 `json:"expires_at"`
	RevokedAt     *time.Time                `json:"revoked_at,omitempty"`
	AutoDeletedAt *time.Time                `json:"auto_deleted_at,omitempty"`
}

func NewConsentLog(userID uuid.UUID, purpose string, grantedAt time.Time) *ConsentLog {
	return &ConsentLog{
		ID:          uuid.New(),
		UserID:      userID,
		Purpose:     purpose,
		Status:      valueobject.ConsentGranted,
		ConsentedAt: grantedAt,
		ExpiresAt:   grantedAt.AddDate(1, 0, 0),
	}
}

func (c *ConsentLog) Revoke() error {
	if c.Status != valueobject.ConsentGranted {
		return domainerrors.ErrRevokeNotAllowed
	}
	now := time.Now()
	c.Status = valueobject.ConsentRevoked
	c.RevokedAt = &now
	return nil
}

func (c *ConsentLog) MarkDeleted() {
	now := time.Now()
	c.Status = valueobject.ConsentDeleted
	c.AutoDeletedAt = &now
}

func (c *ConsentLog) IsActive() bool {
	return c.Status == valueobject.ConsentGranted && c.ExpiresAt.After(time.Now())
}