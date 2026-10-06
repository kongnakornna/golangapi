package models

import "time"

type SdModules struct {
	ModuleID       int       `gorm:"column:module_id;primaryKey;autoIncrement" json:"module_id"`
	ModuleParentID *int      `gorm:"column:module_parent_id;index" json:"module_parent_id,omitempty"`
	Name           string    `gorm:"column:name;type:text;not null" json:"name"`
	Path           *string   `gorm:"column:path;type:text" json:"path,omitempty"`
	Icon           *string   `gorm:"column:icon;type:text" json:"icon,omitempty"`
	Sort           int       `gorm:"column:sort;not null;default:0;index" json:"sort"`
	CreateDate     time.Time `gorm:"column:create_date;not null;default:now()" json:"create_date"`
	UpdateDate     time.Time `gorm:"column:update_date;not null;default:now()" json:"update_date"`
	Status         int16     `gorm:"column:status;not null;default:1;index" json:"status"`
}

func (SdModules) TableName() string { return "sd_modules" }