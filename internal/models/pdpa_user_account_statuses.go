package models

import "time"

// PdpaUserAccountStatus maps migrations/20260906_pdpa_user_account_statuses.sql
type PdpaUserAccountStatus struct {
	ID                 string     `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID             string     `gorm:"column:user_id;type:uuid;not null;uniqueIndex:uq_pdpa_user_account_statuses_user_id"`
	AccountStatus      int64      `gorm:"column:account_status;type:bigint;not null;default:1;index:idx_pdpa_user_account_statuses_status"`
	DeletionRequestedAt *time.Time `gorm:"column:deletion_requested_at;type:timestamptz"`
	DeletionDueAt      *time.Time `gorm:"column:deletion_due_at;type:timestamptz;index:idx_pdpa_user_account_statuses_deletion_due_at"`
	GracePeriodDays    *int32     `gorm:"column:grace_period_days;type:int"`
	DeletionReason     *string    `gorm:"column:deletion_reason;type:varchar(255)"`
	ErasedAt           *time.Time `gorm:"column:erased_at;type:timestamptz"`
	ManuallyDeleted    bool       `gorm:"column:manually_deleted;type:boolean;not null;default:false"`
	DeletedBy          *string    `gorm:"column:deleted_by;type:uuid"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt          time.Time  `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
}

func (PdpaUserAccountStatus) TableName() string {
	return "pdpa_user_account_statuses"
}