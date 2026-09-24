package models

import (
	"time"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

// ActivityLog represents the activity_log table
// ActivityLog แทนตาราง activity_log สำหรับบันทึกกิจกรรมต่างๆ
type ActivityLog struct {
	ID            int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Type          string         `gorm:"type:varchar(50);not null;column:type" json:"type"`                         // ประเภทของกิจกรรม
	DeviceID      *string        `gorm:"type:varchar(50);column:deviceId" json:"deviceId"`                           // รหัสอุปกรณ์ (ถ้ามี)
	UserID        *string        `gorm:"type:varchar(100);column:userId" json:"userId"`                              // รหัสผู้ใช้ (ถ้ามี)
	Details       string         `gorm:"type:varchar(500);not null;column:details" json:"details"`                  // รายละเอียดกิจกรรม
	Data          datatypes.JSON `gorm:"type:jsonb;column:data" json:"data"`                                        // ข้อมูลเพิ่มเติมในรูปแบบ JSON
	Severity      string         `gorm:"type:varchar(20);not null;default:info;column:severity" json:"severity"`   // ระดับความรุนแรง (info, warning, error)
	IPAddress     *string        `gorm:"type:varchar(45);column:ipAddress" json:"ipAddress"`                        // ที่อยู่ IP
	UserAgent     *string        `gorm:"type:varchar(500);column:userAgent" json:"userAgent"`                       // User agent ของเบราว์เซอร์
	SessionID     *string        `gorm:"type:varchar(100);column:sessionId" json:"sessionId"`                       // รหัสเซสชัน
	CorrelationID *string        `gorm:"type:varchar(100);column:correlationId" json:"correlationId"`               // รหัสสำหรับติดตามธุรกรรมข้ามระบบ
	Timestamp     time.Time      `gorm:"type:timestamptz;not null;column:timestamp" json:"timestamp"`               // เวลาที่เกิดเหตุการณ์
	CreatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"` // เวลาที่สร้างระเบียน
	StackTrace    *string        `gorm:"type:text;column:stackTrace" json:"stackTrace"`                             // Stack trace (ถ้ามี error)
}

func (ActivityLog) TableName() string {
	return "activity_log"
}