package models

// SdUserAccessMenu maps to "sd_user_access_menu".
type SdUserAccessMenu struct {
	UserAccessID int  `gorm:"column:user_access_id;primaryKey;autoIncrement" json:"user_access_id"`
	UserTypeID   *int `gorm:"column:user_type_id;index" json:"user_type_id,omitempty"`
	RoleID       *int `gorm:"column:role_id;index" json:"role_id,omitempty"`
	MenuID       *int `gorm:"column:menu_id;index" json:"menu_id,omitempty"`
	ParentID     *int `gorm:"column:parent_id;index" json:"parent_id,omitempty"`
}

func (SdUserAccessMenu) TableName() string { return "sd_user_access_menu" }