package entity

import (
	"time"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type UserAccountStatus struct {
	ID                  uuid.UUID                  `json:"id"`
	UserID              uuid.UUID                  `json:"user_id"`
	Status              valueobject.AccountStatus  `json:"status"`
	SuspendedAt         *time.Time                 `json:"suspended_at,omitempty"`
	TerminatedAt        *time.Time                 `json:"terminated_at,omitempty"`
	DeletionConfirmedAt *time.Time                 `json:"deletion_confirmed_at,omitempty"`
	RetentionDeadline   *time.Time                 `json:"retention_deadline,omitempty"`
	AutoDeletedAt       *time.Time                 `json:"auto_deleted_at,omitempty"`
	UpdatedAt           time.Time                  `json:"updated_at"`
}

func NewUserAccountStatus(userID uuid.UUID, status valueobject.AccountStatus) *UserAccountStatus {
	return &UserAccountStatus{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    status,
		UpdatedAt: time.Now(),
	}
}

func (u *UserAccountStatus) SetSuspended(suspendedAt time.Time, retentionYears int) {
	u.Status = valueobject.AccountSuspended
	u.SuspendedAt = &suspendedAt
	deadline := suspendedAt.AddDate(retentionYears, 0, 0)
	u.RetentionDeadline = &deadline
}

func (u *UserAccountStatus) SetTerminated(terminatedAt time.Time) {
	u.Status = valueobject.AccountTerminated
	u.TerminatedAt = &terminatedAt
}

func (u *UserAccountStatus) ConfirmDeletion() {
	now := time.Now()
	u.DeletionConfirmedAt = &now
}

func (u *UserAccountStatus) MarkDeleted() {
	now := time.Now()
	u.Status = valueobject.AccountTerminated
	u.AutoDeletedAt = &now
}

func (u *UserAccountStatus) IsActive() bool {
	return u.Status == valueobject.AccountActive
}

func (u *UserAccountStatus) IsSuspended() bool {
	return u.Status == valueobject.AccountSuspended
}

func (u *UserAccountStatus) IsTerminated() bool {
	return u.Status == valueobject.AccountTerminated
}

func (u *UserAccountStatus) IsDeleted() bool {
	return u.Status == valueobject.AccountTerminated
}