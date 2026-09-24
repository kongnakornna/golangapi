package models

import "time"

// PdpaUserRequest maps migrations/20260905_pdpa_user_requests.sql
type PdpaUserRequest struct {
	ID                 string     `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserID             string     `gorm:"column:user_id;type:uuid;not null;index:idx_pdpa_user_requests_user_id"`
	RequestType        int64      `gorm:"column:request_type;type:bigint;not null;default:1"`
	Status             int64      `gorm:"type:bigint;not null;default:1;index:idx_pdpa_user_requests_status"`
	Priority           int64      `gorm:"type:bigint;not null;default:1"`
	RequestChannel     *string    `gorm:"column:request_channel;type:varchar(100)"`
	RequestDetail      *string    `gorm:"column:request_detail;type:text"`
	DeadlineAt         time.Time  `gorm:"column:deadline_at;type:timestamptz;not null;index:idx_pdpa_user_requests_deadline_at"`
	CompletedAt        *time.Time `gorm:"column:completed_at;type:timestamptz"`
	RejectedReason     *string    `gorm:"column:rejected_reason;type:text"`
	AssigneeID         *string    `gorm:"column:assignee_id;type:uuid"`
	RequiresVerification bool     `gorm:"column:requires_verification;type:boolean;not null;default:true"`
	VerificationStatus int64      `gorm:"column:verification_status;type:bigint;not null;default:0"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;not null;default:now();index:idx_pdpa_user_requests_created_at,priority:1,sort:desc"`
	UpdatedAt          time.Time  `gorm:"type:timestamptz;not null;default:now();autoUpdateTime"`
	DeletedAt          *time.Time `gorm:"type:timestamptz;index"`
}

func (PdpaUserRequest) TableName() string {
	return "pdpa_user_requests"
}