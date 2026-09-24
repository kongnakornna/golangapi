package models

import "time"

type SdUserRolePermission struct {
	RoleTypeID int        `gorm:"column:role_type_id;primaryKey"`
	Name       string     `gorm:"column:name"`
	Detail     *string    `gorm:"column:detail"`
	CreatedAt  *time.Time `gorm:"column:created"`
	UpdatedAt  *time.Time `gorm:"column:updated"`
	Insert     int        `gorm:"column:insert"`
	Update     int        `gorm:"column:update"`
	Delete     int        `gorm:"column:delete"`
	Select     int        `gorm:"column:select"`
	Log        int        `gorm:"column:log"`
	Config     int        `gorm:"column:config"`
	Truncate   int        `gorm:"column:truncate"`
}

func (SdUserRolePermission) TableName() string { return "sd_user_roles_permision" }
