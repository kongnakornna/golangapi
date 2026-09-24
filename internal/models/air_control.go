package models

import "time"

type AirControl struct {
	AirControlID int       `gorm:"primaryKey;autoIncrement" json:"air_control_id"`
	Name         string    `gorm:"size:255" json:"name"`
	Data         string    `gorm:"type:text" json:"data"`
	Active       int       `gorm:"default:0" json:"active"`
	Status       int       `gorm:"default:1" json:"status"`
	CreatedDate  time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate  time.Time `gorm:"default:now()" json:"updated_date"`
}

func (AirControl) TableName() string { return "sd_air_control" }

type AirMod struct {
	AirModID    int       `gorm:"primaryKey;autoIncrement" json:"air_mod_id"`
	Name        string    `gorm:"size:255" json:"name"`
	Data        string    `gorm:"type:text" json:"data"`
	Active      int       `gorm:"default:0" json:"active"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedDate time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate time.Time `gorm:"default:now()" json:"updated_date"`
}

func (AirMod) TableName() string { return "sd_air_mod" }

type AirPeriod struct {
	AirPeriodID int       `gorm:"primaryKey;autoIncrement" json:"air_period_id"`
	Name        string    `gorm:"size:255" json:"name"`
	Data        string    `gorm:"type:text" json:"data"`
	Active      int       `gorm:"default:0" json:"active"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedDate time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate time.Time `gorm:"default:now()" json:"updated_date"`
}

func (AirPeriod) TableName() string { return "sd_air_period" }

type AirSettingWarning struct {
	AirSettingWarningID int       `gorm:"primaryKey;autoIncrement" json:"air_setting_warning_id"`
	TypeID              int       `json:"type_id"`
	DeviceID            int       `json:"device_id"`
	PeriodID            int       `json:"period_id"`
	EventName           string    `json:"event_name"`
	Name                string    `json:"name"`
	Data                string    `gorm:"type:text" json:"data"`
	Active              int       `gorm:"default:0" json:"active"`
	Status              int       `gorm:"default:1" json:"status"`
	CreatedDate         time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate         time.Time `gorm:"default:now()" json:"updated_date"`
}

func (AirSettingWarning) TableName() string { return "sd_air_setting_warning" }

type AirWarning struct {
	AirWarningID int       `gorm:"primaryKey;autoIncrement" json:"air_warning_id"`
	Name         string    `json:"name"`
	Data         string    `gorm:"type:text" json:"data"`
	Active       int       `gorm:"default:0" json:"active"`
	Status       int       `gorm:"default:1" json:"status"`
	CreatedDate  time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate  time.Time `gorm:"default:now()" json:"updated_date"`
}

func (AirWarning) TableName() string { return "sd_air_warning" }
