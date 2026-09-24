package models

import (
	"time"
)

type SdUserRole struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement"`
	RoleID       int        `gorm:"column:role_id;not null"`
	Title        *string    `gorm:"column:title;type:varchar(50)"`
	CreatedDate  *time.Time `gorm:"column:createddate"`
	UpdatedDate  *time.Time `gorm:"column:updateddate"`
	CreateBy     int        `gorm:"column:create_by;not null"`
	LastupdateBy int        `gorm:"column:lastupdate_by;not null"`
	Status       int16      `gorm:"column:status;not null"`
	TypeID       int        `gorm:"column:type_id;not null"`
	Lang         string     `gorm:"column:lang;type:varchar(255);not null"`
}

func (SdUserRole) TableName() string {
	return "sd_user_role"
}
