package models

import (
	"time"
)

type Device struct {
	DeviceID             int `gorm:"primaryKey;autoIncrement"`
	TypeID               int `gorm:"index"`
	HardwareID           int
	DeviceName           string
	MQTTDataValue        string
	MQTTDataControl      string
	MQTTControlOn        string
	MQTTControlOff       string
	Org                  string
	Bucket               string
	DataStatus           int
	MQTTID               int
	SettingID            int
	SN                   string
	StatusWarning        string
	RecoveryWarning      string
	StatusAlert          string
	RecoveryAlert        string
	TimeLife             int
	Period               string
	WorkStatus           int
	Max                  string
	Min                  string
	OID                  string
	CompareValue         string
	Unit                 string
	ActionID             int
	StatusAlertID        int
	Measurement          string
	LocationID           int
	ConfigData           string
	MQTTDeviceName       string
	MQTTStatusOverName   string
	MQTTStatusDataName   string
	MQTTActRelayName     string
	MQTTControlRelayName string
	Layout               int
	AlertSet             int
	IconNormal           string
	IconWarning          string
	IconAlert            string
	Icon                 string
	IconOn               string
	IconOff              string
	ColorNormal          string
	ColorWarning         string
	ColorAlert           string
	Code                 string
	CalibrationAdd       float64
	CalibrationSubtract  float64
	CalibrationType      int
	CreatedDate          time.Time `gorm:"default:now()"`
	UpdatedDate          time.Time `gorm:"default:now()"`
	Status               int       `gorm:"default:1"`
}

func (Device) TableName() string {
	return "sd_iot_device"
}
