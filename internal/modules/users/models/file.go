package models

import "time"

type SdUserFile struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FileName   string    `gorm:"column:file_name;type:text;not null" json:"file_name"`
	FileType   string    `gorm:"column:file_type;type:text;not null" json:"file_type"`
	FilePath   string    `gorm:"column:file_path;type:text;not null" json:"file_path"`
	FileTypeID int       `gorm:"column:file_type_id;not null;index" json:"file_type_id"`
	UID        *string   `gorm:"column:uid;type:text;index" json:"uid,omitempty"`
	FileDate   time.Time `gorm:"column:file_date;not null;default:now()" json:"file_date"`
	Status     int16     `gorm:"column:status;not null;index" json:"status"`
}

func (SdUserFile) TableName() string { return "sd_user_file" }