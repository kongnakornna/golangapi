package dto

import (
	"time"

	"github.com/google/uuid"

	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

// UserAccountStatusResponse represents the retention and lifecycle status of a user account.
// UserAccountStatusResponse แสดงสถานะวงจรชีวิตและระยะเวลาการเก็บรักษาบัญชีผู้ใช้
type UserAccountStatusResponse struct {
	UserID               uuid.UUID                 `json:"user_id"`
	Status               valueobject.AccountStatus `json:"status"`
	RetentionDeadline    *time.Time                `json:"retention_deadline,omitempty"`
	DeletionConfirmedAt  *time.Time                `json:"deletion_confirmed_at,omitempty"`
	AutoDeletedAt        *time.Time                `json:"auto_deleted_at,omitempty"`
}

// UserAccountStatusListResponse represents a list of user account statuses.
// UserAccountStatusListResponse แสดงรายการสถานะวงจรชีวิตบัญชีผู้ใช้
type UserAccountStatusListResponse struct {
	Accounts []UserAccountStatusResponse `json:"accounts"`
	Total    int64                       `json:"total"`
}