package models

import "time"

type SdUserLogType struct {
	LogTypeID  int       `gorm:"column:log_type_id;primaryKey;autoIncrement" json:"log_type_id"`
	TypeName   string    `gorm:"column:type_name;type:text;not null" json:"type_name"`
	TypeDetail string    `gorm:"column:type_detail;type:text;not null" json:"type_detail"`
	Status     *int      `gorm:"column:status" json:"status,omitempty"`
	Create     time.Time `gorm:"column:create;not null;default:now()" json:"create"`
	Update     time.Time `gorm:"column:update;not null;default:now()" json:"update"`
}

func (SdUserLogType) TableName() string { return "sd_user_log_type" }