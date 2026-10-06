package entity

import (
	"time"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type UserAccountStatus struct {
	ID                  uuid.UUID                 `json:"id" gorm:"column:id;primaryKey"`
	UserID              uuid.UUID                 `json:"user_id" gorm:"column:user_id;uniqueIndex"`
	Status              valueobject.AccountStatus `json:"status" gorm:"column:status"`
	SuspendedAt         *time.Time                `json:"suspended_at,omitempty" gorm:"column:suspended_at"`
	TerminatedAt        *time.Time                `json:"terminated_at,omitempty" gorm:"column:terminated_at"`
	DeletionConfirmedAt *time.Time                `json:"deletion_confirmed_at,omitempty" gorm:"column:deletion_confirmed_at"`
	RetentionDeadline   *time.Time                `json:"retention_deadline,omitempty" gorm:"column:retention_deadline"`
	AutoDeletedAt       *time.Time                `json:"auto_deleted_at,omitempty" gorm:"column:auto_deleted_at"`
	UpdatedAt           time.Time                 `json:"updated_at" gorm:"column:updated_at"`
}

func (UserAccountStatus) TableName() string { return "pdpa_user_account_statuses" }

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

func (u *UserAccountStatus) IsActive() bool     { return u.Status == valueobject.AccountActive }
func (u *UserAccountStatus) IsSuspended() bool  { return u.Status == valueobject.AccountSuspended }
func (u *UserAccountStatus) IsTerminated() bool { return u.Status == valueobject.AccountTerminated }
func (u *UserAccountStatus) IsDeleted() bool    { return u.Status == valueobject.AccountTerminated }
