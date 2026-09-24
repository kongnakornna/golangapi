package models

import "time"

type SdUserRolesAccess struct {
	RoleID     int        `gorm:"column:role_id;primaryKey"`
	RoleTypeID int        `gorm:"column:role_type_id;primaryKey"`
	CreateDate *time.Time `gorm:"column:create"`
	UpdateDate *time.Time `gorm:"column:update"`
}

func (SdUserRolesAccess) TableName() string { return "sd_user_roles_access" }
