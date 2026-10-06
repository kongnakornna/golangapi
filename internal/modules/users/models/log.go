package models

import (
	"time"

	"github.com/google/uuid"
)

type SdUserLog struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	LogTypeID int       `gorm:"column:log_type_id;primaryKey;index" json:"log_type_id"`
	UID       uuid.UUID `gorm:"column:uid;type:uuid;not null;index" json:"uid"`
	Name      string    `gorm:"column:name;type:text;not null" json:"name"`
	Detail    string    `gorm:"column:detail;type:text;not null" json:"detail"`

	IPAddress    *string `gorm:"column:ipaddress;type:text;index" json:"ipaddress,omitempty"`
	DeviceAccess *string `gorm:"column:device_access;type:text" json:"device_access,omitempty"`
	BrowserAgent *string `gorm:"column:browser_agent;type:text" json:"browser_agent,omitempty"`

	SelectStatus *int      `gorm:"column:select_status" json:"select_status,omitempty"`
	InsertStatus *int      `gorm:"column:insert_status" json:"insert_status,omitempty"`
	UpdateStatus *int      `gorm:"column:update_status" json:"update_status,omitempty"`
	DeleteStatus *int      `gorm:"column:delete_status" json:"delete_status,omitempty"`
	Status       *int      `gorm:"column:status" json:"status,omitempty"`
	Create       time.Time `gorm:"column:create;not null;default:now()" json:"create"`
	Update       time.Time `gorm:"column:update;not null;default:now()" json:"update"`
	Lang         *string   `gorm:"column:lang;type:text" json:"lang,omitempty"`
}

func (SdUserLog) TableName() string { return "sd_user_log" }