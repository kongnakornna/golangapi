```go
package models

import (
	"time"

	"gorm.io/datatypes"
)

// ActivityLog represents the activity_log table
// ActivityLog แทนตาราง activity_log สำหรับบันทึกกิจกรรมต่างๆ
type ActivityLog struct {
	ID            int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Type          string         `gorm:"type:varchar(50);not null;column:type" json:"type"`                                           // ประเภทของกิจกรรม
	DeviceID      *string        `gorm:"type:varchar(50);column:deviceId" json:"deviceId"`                                             // รหัสอุปกรณ์ (ถ้ามี)
	UserID        *string        `gorm:"type:varchar(100);column:userId" json:"userId"`                                                // รหัสผู้ใช้ (ถ้ามี)
	Details       string         `gorm:"type:varchar(500);not null;column:details" json:"details"`                                    // รายละเอียดกิจกรรม
	Data          datatypes.JSON `gorm:"type:jsonb;column:data" json:"data"`                                                           // ข้อมูลเพิ่มเติมในรูปแบบ JSON
	Severity      string         `gorm:"type:varchar(20);not null;default:info;column:severity" json:"severity"`                      // ระดับความรุนแรง (info, warning, error)
	IPAddress     *string        `gorm:"type:varchar(45);column:ipAddress" json:"ipAddress"`                                          // ที่อยู่ IP
	UserAgent     *string        `gorm:"type:varchar(500);column:userAgent" json:"userAgent"`                                         // User agent ของเบราว์เซอร์
	SessionID     *string        `gorm:"type:varchar(100);column:sessionId" json:"sessionId"`                                         // รหัสเซสชัน
	CorrelationID *string        `gorm:"type:varchar(100);column:correlationId" json:"correlationId"`                                 // รหัสสำหรับติดตามธุรกรรมข้ามระบบ
	Timestamp     time.Time      `gorm:"type:timestamptz;not null;column:timestamp" json:"timestamp"`                                 // เวลาที่เกิดเหตุการณ์
	CreatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`                   // เวลาที่สร้างระเบียน
	StackTrace    *string        `gorm:"type:text;column:stackTrace" json:"stackTrace"`                                               // Stack trace (ถ้ามี error)
}

func (ActivityLog) TableName() string {
	return "activity_log"
}

// CommandLog represents the command_log table
// CommandLog แทนตาราง command_log สำหรับบันทึกคำสั่งที่ส่งไปยังอุปกรณ์
type CommandLog struct {
	ID         int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DeviceID   string         `gorm:"type:varchar(50);not null;column:deviceId" json:"deviceId"`                 // รหัสอุปกรณ์เป้าหมาย
	Action     string         `gorm:"type:varchar(100);not null;column:action" json:"action"`                    // ชื่อคำสั่ง
	Parameters datatypes.JSON `gorm:"type:jsonb;column:parameters" json:"parameters"`                            // พารามิเตอร์ของคำสั่ง
	Metadata   datatypes.JSON `gorm:"type:jsonb;column:metadata" json:"metadata"`                                // ข้อมูลเมตาดาต้าเพิ่มเติม
	Status     string         `gorm:"type:varchar(50);not null;default:pending;column:status" json:"status"`    // สถานะคำสั่ง (pending, sent, executed, failed)
	IssuedBy   *string        `gorm:"type:varchar(100);column:issuedBy" json:"issuedBy"`                        // ผู้ส่งคำสั่ง
	ClientIP   *string        `gorm:"type:varchar(45);column:clientIp" json:"clientIp"`                         // IP ของผู้ส่ง
	Response   datatypes.JSON `gorm:"type:jsonb;column:response" json:"response"`                               // ข้อมูลตอบกลับ
	Error      *string        `gorm:"type:varchar(500);column:error" json:"error"`                              // ข้อผิดพลาด (ถ้ามี)
	IssuedAt   time.Time      `gorm:"type:timestamptz;not null;column:issuedAt" json:"issuedAt"`                // เวลาที่ออกคำสั่ง
	SentAt     *time.Time     `gorm:"type:timestamptz;column:sentAt" json:"sentAt"`                             // เวลาที่ส่งคำสั่ง
	ExecutedAt *time.Time     `gorm:"type:timestamptz;column:executedAt" json:"executedAt"`                     // เวลาที่ดำเนินการสำเร็จ
	FailedAt   *time.Time     `gorm:"type:timestamptz;column:failedAt" json:"failedAt"`                         // เวลาที่เกิดความล้มเหลว
	CreatedAt  time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updatedAt" json:"updatedAt"`
}

func (CommandLog) TableName() string {
	return "command_log"
}

// DeviceAlert represents the device_alert table
// DeviceAlert แทนตาราง device_alert สำหรับบันทึกการแจ้งเตือนจากอุปกรณ์
type DeviceAlert struct {
	ID                int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DeviceID          string         `gorm:"type:varchar(50);not null;column:deviceId" json:"deviceId"`                       // รหัสอุปกรณ์
	Type              string         `gorm:"type:varchar(50);not null;column:type" json:"type"`                               // ประเภทการแจ้งเตือน
	Metric            *string        `gorm:"type:varchar(100);column:metric" json:"metric"`                                   // ชื่อเมตริกที่เกี่ยวข้อง
	Value             *float64       `gorm:"column:value" json:"value"`                                                       // ค่าที่เกิดขึ้น
	Threshold         datatypes.JSON `gorm:"type:jsonb;column:threshold" json:"threshold"`                                   // ค่าเกณฑ์ที่กำหนด
	Severity          string         `gorm:"type:varchar(20);not null;default:low;column:severity" json:"severity"`          // ระดับความรุนแรง (low, medium, high)
	Message           string         `gorm:"type:varchar(500);not null;column:message" json:"message"`                       // ข้อความแจ้งเตือน
	Details           datatypes.JSON `gorm:"type:jsonb;column:details" json:"details"`                                       // รายละเอียดเพิ่มเติม
	Resolved          bool           `gorm:"not null;default:false;column:resolved" json:"resolved"`                         // สถานะแก้ไขแล้วหรือไม่
	ResolutionNotes   *string        `gorm:"type:text;column:resolutionNotes" json:"resolutionNotes"`                        // หมายเหตุการแก้ไข
	ResolvedBy        *string        `gorm:"type:varchar(100);column:resolvedBy" json:"resolvedBy"`                          // ผู้แก้ไข
	ResolvedAt        *time.Time     `gorm:"type:timestamptz;column:resolvedAt" json:"resolvedAt"`                           // เวลาที่แก้ไข
	Acknowledged      bool           `gorm:"not null;default:false;column:acknowledged" json:"acknowledged"`                 // สถานะรับทราบแล้ว
	AcknowledgedBy    *string        `gorm:"type:varchar(100);column:acknowledgedBy" json:"acknowledgedBy"`                  // ผู้รับทราบ
	AcknowledgedAt    *time.Time     `gorm:"type:timestamptz;column:acknowledgedAt" json:"acknowledgedAt"`                   // เวลาที่รับทราบ
	Escalation        datatypes.JSON `gorm:"type:jsonb;column:escalation" json:"escalation"`                                 // ข้อมูลการเลื่อนระดับ
	DataID            *int           `gorm:"column:dataId" json:"dataId"`                                                    // อ้างอิง iot_data.id (ถ้ามี)
	CreatedAt         time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updatedAt" json:"updatedAt"`
	ExpiresAt         *time.Time     `gorm:"type:timestamptz;column:expiresAt" json:"expiresAt"`                             // เวลาหมดอายุ
	NotificationCount int            `gorm:"not null;default:0;column:notificationCount" json:"notificationCount"`          // จำนวนครั้งที่แจ้งเตือน
}

func (DeviceAlert) TableName() string {
	return "device_alert"
}

// DeviceConfig represents the device_config table
// DeviceConfig แทนตาราง device_config สำหรับเก็บการตั้งค่าคอนฟิกของอุปกรณ์
type DeviceConfig struct {
	ID            int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DeviceID      string         `gorm:"type:varchar(50);not null;uniqueIndex;column:deviceId" json:"deviceId"` // รหัสอุปกรณ์ (unique)
	Config        datatypes.JSON `gorm:"type:jsonb;column:config" json:"config"`                                // ค่าการตั้งค่าในรูปแบบ JSON
	Status        string         `gorm:"type:varchar(20);not null;default:active;column:status" json:"status"` // สถานะคอนฟิก (active, inactive)
	Notes         *string        `gorm:"type:text;column:notes" json:"notes"`                                   // หมายเหตุ
	UpdatedBy     *string        `gorm:"type:varchar(100);column:updatedBy" json:"updatedBy"`                  // ผู้ปรับปรุงล่าสุด
	CreatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updatedAt" json:"updatedAt"`
	LastAppliedAt *time.Time     `gorm:"type:timestamptz;column:lastAppliedAt" json:"lastAppliedAt"` // เวลาที่นำคอนฟิกไปใช้ล่าสุด
}

func (DeviceConfig) TableName() string {
	return "device_config"
}

// DeviceStatus represents the device_status table
// DeviceStatus แทนตาราง device_status สำหรับเก็บสถานะล่าสุดของอุปกรณ์
type DeviceStatus struct {
	ID               int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DeviceID         string         `gorm:"type:varchar(50);not null;uniqueIndex;column:deviceId" json:"deviceId"` // รหัสอุปกรณ์ (unique)
	IsOnline         bool           `gorm:"not null;default:true;column:isOnline" json:"isOnline"`                 // สถานะออนไลน์
	IsActive         bool           `gorm:"not null;default:true;column:isActive" json:"isActive"`                 // สถานะการทำงาน
	LastSeen         time.Time      `gorm:"type:timestamptz;not null;column:lastSeen" json:"lastSeen"`             // เวลาที่เห็นครั้งล่าสุด
	LastData         datatypes.JSON `gorm:"type:jsonb;column:lastData" json:"lastData"`                            // ข้อมูลล่าสุด
	BatteryLevel     *int           `gorm:"column:batteryLevel" json:"batteryLevel"`                               // ระดับแบตเตอรี่ (%)
	SignalStrength   *int           `gorm:"column:signalStrength" json:"signalStrength"`                           // ความแรงสัญญาณ
	Temperature      *float64       `gorm:"column:temperature" json:"temperature"`                                 // อุณหภูมิ
	Humidity         *float64       `gorm:"column:humidity" json:"humidity"`                                       // ความชื้น
	FirmwareVersion  *string        `gorm:"type:varchar(20);column:firmwareVersion" json:"firmwareVersion"`        // เวอร์ชันเฟิร์มแวร์
	Uptime           *int           `gorm:"column:uptime" json:"uptime"`                                           // เวลาทำงาน (วินาที)
	Location         datatypes.JSON `gorm:"type:jsonb;column:location" json:"location"`                            // ตำแหน่งที่ตั้ง
	NetworkInfo      datatypes.JSON `gorm:"type:jsonb;column:networkInfo" json:"networkInfo"`                      // ข้อมูลเครือข่าย
	HardwareInfo     datatypes.JSON `gorm:"type:jsonb;column:hardwareInfo" json:"hardwareInfo"`                    // ข้อมูลฮาร์ดแวร์
	Metrics          datatypes.JSON `gorm:"type:jsonb;column:metrics" json:"metrics"`                              // เมตริกอื่นๆ
	StatusMessage    *string        `gorm:"type:text;column:statusMessage" json:"statusMessage"`                   // ข้อความสถานะ
	CustomFields     datatypes.JSON `gorm:"type:jsonb;column:customFields" json:"customFields"`                    // ฟิลด์กำหนดเอง
	CreatedAt        time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updatedAt" json:"updatedAt"`
	FirstSeen        *time.Time     `gorm:"type:timestamptz;column:firstSeen" json:"firstSeen"`        // เวลาเห็นครั้งแรก
	LastMaintenance  *time.Time     `gorm:"type:timestamptz;column:lastMaintenance" json:"lastMaintenance"` // เวลาบำรุงรักษาล่าสุด
	ConnectionCount  int            `gorm:"not null;default:0;column:connectionCount" json:"connectionCount"` // จำนวนการเชื่อมต่อ
}

func (DeviceStatus) TableName() string {
	return "device_status"
}

// IoTData represents the iot_data table
// IoTData แทนตาราง iot_data สำหรับเก็บข้อมูลดิบจากอุปกรณ์ IoT
type IoTData struct {
	ID          int            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Data        datatypes.JSON `gorm:"type:jsonb;not null;column:data" json:"data"`               // ข้อมูลหลักในรูปแบบ JSON
	CreatedAt   time.Time      `gorm:"type:timestamp(6);not null;default:now();column:createdAt" json:"createdAt"`
	Location    datatypes.JSON `gorm:"type:jsonb;column:location" json:"location"`                // ตำแหน่งที่ตั้ง
	Metadata    datatypes.JSON `gorm:"type:jsonb;column:metadata" json:"metadata"`                // เมตาดาต้า
	DataType    *string        `gorm:"type:varchar(20);column:dataType" json:"dataType"`          // ประเภทข้อมูล (เช่น sensor, status)
	DataQuality *float64       `gorm:"column:dataQuality" json:"dataQuality"`                     // คุณภาพข้อมูล (0-1)
	DeviceID    string         `gorm:"type:varchar(50);not null;column:deviceId" json:"deviceId"` // รหัสอุปกรณ์
	Timestamp   time.Time      `gorm:"type:timestamptz;not null;default:now();column:timestamp" json:"timestamp"` // เวลาของข้อมูล
}

func (IoTData) TableName() string {
	return "iot_data"
}

// Migration represents the migrations table (for internal use)
// Migration แทนตาราง migrations ใช้สำหรับการจัดการ schema migration ของระบบ
type Migration struct {
	ID        int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Timestamp int64  `gorm:"not null;column:timestamp" json:"timestamp"` // Unix timestamp ของ migration
	Name      string `gorm:"type:varchar;not null;column:name" json:"name"` // ชื่อไฟล์ migration
}

func (Migration) TableName() string {
	return "migrations"
}

// NotiNotificationLog represents the noti_notification_logs table
// NotiNotificationLog แทนตาราง noti_notification_logs สำหรับบันทึกประวัติการส่งการแจ้งเตือน
type NotiNotificationLog struct {
	LogID          int            `gorm:"primaryKey;autoIncrement;column:log_id" json:"log_id"`
	NotificationID string         `gorm:"type:uuid;not null;column:notification_id;index" json:"notification_id"` // รหัสการแจ้งเตือน (อ้างอิง noti_notifications.id)
	Channel        string         `gorm:"type:varchar(50);not null;column:channel" json:"channel"`                // ช่องทางที่ใช้ส่ง (email, line, sms, etc.)
	Payload        datatypes.JSON `gorm:"type:jsonb;not null;column:payload" json:"payload"`                      // ข้อมูล payload ที่ส่ง
	Response       datatypes.JSON `gorm:"type:jsonb;column:response" json:"response"`                             // ข้อมูลตอบกลับจากระบบรับส่ง
	Status         string         `gorm:"type:varchar(50);not null;default:pending;column:status" json:"status"` // สถานะ (pending, sent, failed)
	RetryCount     *int           `gorm:"default:0;column:retry_count" json:"retry_count"`                        // จำนวนครั้งที่ลองส่งซ้ำ
	ErrorMessage   *string        `gorm:"type:text;column:error_message" json:"error_message"`                    // ข้อความผิดพลาด
	CreatedAt      *time.Time     `gorm:"type:timestamp(6);default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	SentAt         *time.Time     `gorm:"type:timestamp(6);column:sent_at" json:"sent_at"`       // เวลาที่ส่ง
	DeliveredAt    *time.Time     `gorm:"type:timestamp(6);column:delivered_at" json:"delivered_at"` // เวลาที่ผู้รับได้รับ (ถ้ามี)
}

func (NotiNotificationLog) TableName() string {
	return "noti_notification_logs"
}

// NotiNotificationRule represents the noti_notification_rules table
// NotiNotificationRule แทนตาราง noti_notification_rules สำหรับเก็บกฎการแจ้งเตือน
type NotiNotificationRule struct {
	RuleID        int            `gorm:"primaryKey;autoIncrement;column:rule_id" json:"rule_id"`
	Name          string         `gorm:"type:varchar(100);not null;column:name" json:"name"`                   // ชื่อกฎ
	Description   string         `gorm:"type:varchar(255);not null;column:description" json:"description"`     // คำอธิบายกฎ
	EventTrigger  string         `gorm:"type:varchar(100);not null;column:event_trigger" json:"event_trigger"` // เหตุการณ์ที่触发 (เช่น device.data.received)
	Conditions    datatypes.JSON `gorm:"type:jsonb;not null;column:conditions" json:"conditions"`             // เงื่อนไขในรูปแบบ JSON
	Actions       datatypes.JSON `gorm:"type:jsonb;not null;column:actions" json:"actions"`                   // การดำเนินการเมื่อเงื่อนไขเป็นจริง
	IsActive      bool           `gorm:"not null;default:true;column:is_active" json:"is_active"`             // สถานะการใช้งาน
	Priority      int            `gorm:"not null;default:1;column:priority" json:"priority"`                  // ค่าความสำคัญ (ยิ่งน้อยยิ่งสำคัญ)
	CreatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updated_at" json:"updated_at"`
}

func (NotiNotificationRule) TableName() string {
	return "noti_notification_rules"
}

// NotiNotificationType represents the noti_notification_types table
// NotiNotificationType แทนตาราง noti_notification_types สำหรับเก็บประเภทของการแจ้งเตือน
type NotiNotificationType struct {
	TypeID          int            `gorm:"primaryKey;autoIncrement;column:type_id" json:"type_id"`
	Name            string         `gorm:"type:varchar(100);not null;column:name" json:"name"`                   // ชื่อประเภท
	Description     string         `gorm:"type:varchar(255);not null;column:description" json:"description"`     // คำอธิบาย
	DefaultTemplate datatypes.JSON `gorm:"type:jsonb;column:default_template" json:"default_template"`           // เทมเพลตเริ่มต้น
	AllowedChannels datatypes.JSON `gorm:"type:jsonb;column:allowed_channels" json:"allowed_channels"`           // ช่องทางที่อนุญาต (array)
	Status          int            `gorm:"not null;default:1;column:status" json:"status"`                      // สถานะ (1=active, 0=inactive)
	CreatedAt       time.Time      `gorm:"type:timestamp(6);not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamp(6);not null;default:now();column:updated_at" json:"updated_at"`
}

func (NotiNotificationType) TableName() string {
	return "noti_notification_types"
}
```