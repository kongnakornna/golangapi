package models

import (
	"time"
)

type AlarmLog struct {
	ID              uint   `gorm:"primaryKey"`
	AlarmActionID   int    `gorm:"index"`
	DeviceID        string `gorm:"size:36;index"`
	TypeID          int
	Event           int
	Status          int
	AlarmType       int
	StatusWarning   int
	RecoveryWarning int
	StatusAlert     int
	RecoveryAlert   int
	EmailAlarm      int
	LineAlarm       int
	TelegramAlarm   int
	SmsAlarm        int
	NoncAlarm       int
	Date            string `gorm:"size:10"`
	Time            string `gorm:"size:8"`
	Data            string
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	DataAlarm       string
	AlarmStatus     int
	Subject         string
	Content         string
}
