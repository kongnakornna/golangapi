package models

import (
	"time"

	"github.com/google/uuid"
)

// PdpaUserAccountStatus แทนสถานะวงจรชีวิตของบัญชีผู้ใช้
// PdpaUserAccountStatus represents the lifecycle status of a user account.
type PdpaUserAccountStatus struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_pdpa_account_user_id" json:"user_id"`
	Status              int8       `gorm:"not null;default:1;index:idx_pdpa_account_status" json:"status"` // 1=ACTIVE, 2=SUSPENDED, 4=TERMINATED, 8=DELETED
	SuspendedAt         *time.Time `json:"suspended_at,omitempty"`
	TerminatedAt        *time.Time `json:"terminated_at,omitempty"`
	DeletionConfirmedAt *time.Time `json:"deletion_confirmed_at,omitempty"`
	RetentionDeadline   *time.Time `json:"retention_deadline,omitempty"`
	AutoDeletedAt       *time.Time `json:"auto_deleted_at,omitempty"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PdpaUserAccountStatus) TableName() string {
	return "pdpa_user_account_statuses"
}
