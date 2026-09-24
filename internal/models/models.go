// models.go
// Package models defines all database models for the IoT monitoring system.
// โมเดลทั้งหมดของระบบตรวจสอบ IoT
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// ENUM Types (ประเภท Enum)
// ============================================================================

// UserRoleEnum defines user role levels.
// ระดับบทบาทของผู้ใช้
type UserRoleEnum string

const (
	RoleSuperAdmin UserRoleEnum = "SUPERADMIN"
	RoleAdmin      UserRoleEnum = "ADMIN"
	RoleEditor     UserRoleEnum = "EDITOR"
	RoleMonitor    UserRoleEnum = "MONITOR"
	RoleUser       UserRoleEnum = "USER"
)

// UserUsertypeEnum defines user type categories.
// ประเภทของผู้ใช้
type UserUsertypeEnum string

const (
	UsertypeTherapist  UserUsertypeEnum = "therapist"
	UsertypeSupervisor UserUsertypeEnum = "supervisor"
	UsertypeSuperadmin UserUsertypeEnum = "superadmin"
	UsertypeSystem     UserUsertypeEnum = "system"
	UsertypeAdmin      UserUsertypeEnum = "admin"
	UsertypeSupport    UserUsertypeEnum = "support"
	UsertypeEnduser    UserUsertypeEnum = "enduser"
)

type Item struct {
	Id          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(200);not null"`
	OwnerId     uuid.UUID
}

func (Item) TableName() string { return "item" }

// ============================================================================
// Core Tables (ตารางหลัก)
// ============================================================================

// ActivityLog records system activities.
// บันทึกกิจกรรมของระบบ
type ActivityLog struct {
	ID            int             `gorm:"column:id;primaryKey;autoIncrement"`
	Type          string          `gorm:"column:type;not null"`
	DeviceID      *string         `gorm:"column:deviceId"`
	UserID        *string         `gorm:"column:userId"`
	Details       string          `gorm:"column:details;not null"`
	Data          json.RawMessage `gorm:"column:data;type:jsonb"`
	Severity      string          `gorm:"column:severity;not null;default:info"`
	IPAddress     *string         `gorm:"column:ipAddress"`
	UserAgent     *string         `gorm:"column:userAgent"`
	SessionID     *string         `gorm:"column:sessionId"`
	CorrelationID *string         `gorm:"column:correlationId"`
	Timestamp     time.Time       `gorm:"column:timestamp;not null"`
	CreatedAt     time.Time       `gorm:"column:createdAt;not null;default:now()"`
	StackTrace    *string         `gorm:"column:stackTrace;type:text"`
}

func (ActivityLog) TableName() string { return "activity_log" }

// CommandLog logs commands sent to devices.
// บันทึกคำสั่งที่ส่งไปยังอุปกรณ์
type CommandLog struct {
	ID         int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID   string          `gorm:"column:deviceId;not null"`
	Action     string          `gorm:"column:action;not null"`
	Parameters json.RawMessage `gorm:"column:parameters;type:jsonb"`
	Metadata   json.RawMessage `gorm:"column:metadata;type:jsonb"`
	Status     string          `gorm:"column:status;not null;default:pending"`
	IssuedBy   *string         `gorm:"column:issuedBy"`
	ClientIP   *string         `gorm:"column:clientIp"`
	Response   json.RawMessage `gorm:"column:response;type:jsonb"`
	Error      *string         `gorm:"column:error"`
	IssuedAt   time.Time       `gorm:"column:issuedAt;not null"`
	SentAt     *time.Time      `gorm:"column:sentAt"`
	ExecutedAt *time.Time      `gorm:"column:executedAt"`
	FailedAt   *time.Time      `gorm:"column:failedAt"`
	CreatedAt  time.Time       `gorm:"column:createdAt;not null;default:now()"`
	UpdatedAt  time.Time       `gorm:"column:updatedAt;not null;default:now()"`
}

func (CommandLog) TableName() string { return "command_log" }

// DeviceAlert represents device alerts and notifications.
// การแจ้งเตือนจากอุปกรณ์
type DeviceAlert struct {
	ID                int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID          string          `gorm:"column:deviceId;not null"`
	Type              string          `gorm:"column:type;not null"`
	Metric            *string         `gorm:"column:metric"`
	Value             *float64        `gorm:"column:value;type:float8"`
	Threshold         json.RawMessage `gorm:"column:threshold;type:jsonb"`
	Severity          string          `gorm:"column:severity;not null;default:low"`
	Message           string          `gorm:"column:message;not null"`
	Details           json.RawMessage `gorm:"column:details;type:jsonb"`
	Resolved          bool            `gorm:"column:resolved;not null;default:false"`
	ResolutionNotes   *string         `gorm:"column:resolutionNotes;type:text"`
	ResolvedBy        *string         `gorm:"column:resolvedBy"`
	ResolvedAt        *time.Time      `gorm:"column:resolvedAt"`
	Acknowledged      bool            `gorm:"column:acknowledged;not null;default:false"`
	AcknowledgedBy    *string         `gorm:"column:acknowledgedBy"`
	AcknowledgedAt    *time.Time      `gorm:"column:acknowledgedAt"`
	Escalation        json.RawMessage `gorm:"column:escalation;type:jsonb"`
	DataID            *int            `gorm:"column:dataId"`
	CreatedAt         time.Time       `gorm:"column:createdAt;not null;default:now()"`
	UpdatedAt         time.Time       `gorm:"column:updatedAt;not null;default:now()"`
	ExpiresAt         *time.Time      `gorm:"column:expiresAt"`
	NotificationCount int             `gorm:"column:notificationCount;not null;default:0"`
}

func (DeviceAlert) TableName() string { return "device_alert" }

// DeviceConfig stores device configuration.
// การตั้งค่าคอนฟิกของอุปกรณ์
type DeviceConfig struct {
	ID            int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID      string          `gorm:"column:deviceId;not null;unique"`
	Config        json.RawMessage `gorm:"column:config;type:jsonb"`
	Status        string          `gorm:"column:status;not null;default:active"`
	Notes         *string         `gorm:"column:notes;type:text"`
	UpdatedBy     *string         `gorm:"column:updatedBy"`
	CreatedAt     time.Time       `gorm:"column:createdAt;not null;default:now()"`
	UpdatedAt     time.Time       `gorm:"column:updatedAt;not null;default:now()"`
	LastAppliedAt *time.Time      `gorm:"column:lastAppliedAt"`
}

func (DeviceConfig) TableName() string { return "device_config" }

// DeviceStatus tracks current device status.
// สถานะปัจจุบันของอุปกรณ์
type DeviceStatus struct {
	ID              int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID        string          `gorm:"column:deviceId;not null;unique"`
	IsOnline        bool            `gorm:"column:isOnline;not null;default:true"`
	IsActive        bool            `gorm:"column:isActive;not null;default:true"`
	LastSeen        time.Time       `gorm:"column:lastSeen;not null"`
	LastData        json.RawMessage `gorm:"column:lastData;type:jsonb"`
	BatteryLevel    *int            `gorm:"column:batteryLevel"`
	SignalStrength  *int            `gorm:"column:signalStrength"`
	Temperature     *float64        `gorm:"column:temperature;type:float8"`
	Humidity        *float64        `gorm:"column:humidity;type:float8"`
	FirmwareVersion *string         `gorm:"column:firmwareVersion"`
	Uptime          *int            `gorm:"column:uptime"`
	Location        json.RawMessage `gorm:"column:location;type:jsonb"`
	NetworkInfo     json.RawMessage `gorm:"column:networkInfo;type:jsonb"`
	HardwareInfo    json.RawMessage `gorm:"column:hardwareInfo;type:jsonb"`
	Metrics         json.RawMessage `gorm:"column:metrics;type:jsonb"`
	StatusMessage   *string         `gorm:"column:statusMessage;type:text"`
	CustomFields    json.RawMessage `gorm:"column:customFields;type:jsonb"`
	CreatedAt       time.Time       `gorm:"column:createdAt;not null;default:now()"`
	UpdatedAt       time.Time       `gorm:"column:updatedAt;not null;default:now()"`
	FirstSeen       *time.Time      `gorm:"column:firstSeen"`
	LastMaintenance *time.Time      `gorm:"column:lastMaintenance"`
	ConnectionCount int             `gorm:"column:connectionCount;not null;default:0"`
}

func (DeviceStatus) TableName() string { return "device_status" }

// IotData stores raw IoT device data.
// ข้อมูลดิบจากอุปกรณ์ IoT
type IotData struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	Data        json.RawMessage `gorm:"column:data;type:jsonb;not null"`
	CreatedAt   time.Time       `gorm:"column:createdAt;not null;default:now()"`
	Location    json.RawMessage `gorm:"column:location;type:jsonb"`
	Metadata    json.RawMessage `gorm:"column:metadata;type:jsonb"`
	DataType    *string         `gorm:"column:dataType"`
	DataQuality *float64        `gorm:"column:dataQuality;type:float8"`
	DeviceID    string          `gorm:"column:deviceId;not null"`
	Timestamp   time.Time       `gorm:"column:timestamp;not null;default:now()"`
}

func (IotData) TableName() string { return "iot_data" }

// Migration records database migrations.
// บันทึกการย้ายฐานข้อมูล
type Migration struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Timestamp int64  `gorm:"column:timestamp;not null"`
	Name      string `gorm:"column:name;not null"`
}

func (Migration) TableName() string { return "migrations" }

// ============================================================================
// Notification Tables (ตารางการแจ้งเตือน)
// ============================================================================

// NotiNotificationLog logs notification deliveries.
// บันทึกการส่งการแจ้งเตือน
type NotiNotificationLog struct {
	LogID          int             `gorm:"column:log_id;primaryKey;autoIncrement"`
	NotificationID string          `gorm:"column:notification_id;type:uuid;not null"`
	Channel        string          `gorm:"column:channel;not null"`
	Payload        json.RawMessage `gorm:"column:payload;type:jsonb;not null"`
	Response       json.RawMessage `gorm:"column:response;type:jsonb"`
	Status         string          `gorm:"column:status;not null;default:pending"`
	RetryCount     *int            `gorm:"column:retry_count;default:0"`
	ErrorMessage   *string         `gorm:"column:error_message;type:text"`
	CreatedAt      time.Time       `gorm:"column:created_at;default:now()"`
	SentAt         *time.Time      `gorm:"column:sent_at"`
	DeliveredAt    *time.Time      `gorm:"column:delivered_at"`

	// Relationships
	Notification *NotiNotification `gorm:"foreignKey:NotificationID;references:ID"`
}

func (NotiNotificationLog) TableName() string { return "noti_notification_logs" }

// NotiNotificationRule defines notification rules.
// กฎการแจ้งเตือน
type NotiNotificationRule struct {
	RuleID       int             `gorm:"column:rule_id;primaryKey;autoIncrement"`
	Name         string          `gorm:"column:name;not null"`
	Description  string          `gorm:"column:description;not null"`
	EventTrigger string          `gorm:"column:event_trigger;not null"`
	Conditions   json.RawMessage `gorm:"column:conditions;type:jsonb;not null"`
	Actions      json.RawMessage `gorm:"column:actions;type:jsonb;not null"`
	IsActive     bool            `gorm:"column:is_active;not null;default:true"`
	Priority     int             `gorm:"column:priority;not null;default:1"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;not null;default:now()"`
}

func (NotiNotificationRule) TableName() string { return "noti_notification_rules" }

// NotiNotificationType defines notification types.
// ประเภทการแจ้งเตือน
type NotiNotificationType struct {
	TypeID          int             `gorm:"column:type_id;primaryKey;autoIncrement"`
	Name            string          `gorm:"column:name;not null"`
	Description     string          `gorm:"column:description;not null"`
	DefaultTemplate json.RawMessage `gorm:"column:default_template;type:jsonb"`
	AllowedChannels json.RawMessage `gorm:"column:allowed_channels;type:jsonb"`
	Status          int             `gorm:"column:status;not null;default:1"`
	CreatedAt       time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;not null;default:now()"`
}

func (NotiNotificationType) TableName() string { return "noti_notification_types" }

// NotiNotification represents a notification message.
// ข้อความแจ้งเตือนหลัก
type NotiNotification struct {
	ID           string          `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title        string          `gorm:"column:title;not null"`
	Message      string          `gorm:"column:message;type:text;not null"`
	Type         string          `gorm:"column:type;not null"`
	Priority     string          `gorm:"column:priority;not null"`
	Category     *string         `gorm:"column:category"`
	UserID       *int            `gorm:"column:user_id"`
	UserUUID     *string         `gorm:"column:user_uuid;type:uuid"`
	Metadata     json.RawMessage `gorm:"column:metadata;type:jsonb"`
	IsRead       bool            `gorm:"column:is_read;default:false"`
	ReadAt       *time.Time      `gorm:"column:read_at"`
	IsSent       bool            `gorm:"column:is_sent;default:false"`
	ChannelsSent json.RawMessage `gorm:"column:channels_sent;type:jsonb"`
	ScheduledAt  *time.Time      `gorm:"column:scheduled_at"`
	ExpiresAt    *time.Time      `gorm:"column:expires_at"`
	Status       int             `gorm:"column:status;not null;default:1"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;not null;default:now()"`
	DeletedAt    *time.Time      `gorm:"column:deleted_at"`

	// Relationships
	NotificationLogs []NotiNotificationLog `gorm:"foreignKey:NotificationID"`
}

func (NotiNotification) TableName() string { return "noti_notifications" }

// NotificationDevice represents devices used for notifications.
// อุปกรณ์สำหรับการแจ้งเตือน
type NotificationDevice struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID    int        `gorm:"column:device_id;unique"`
	DeviceName  string     `gorm:"column:device_name;not null"`
	DeviceType  string     `gorm:"column:device_type;not null"`
	MqttTopic   *string    `gorm:"column:mqtt_topic"`
	MqttOn      *string    `gorm:"column:mqtt_on"`
	MqttOff     *string    `gorm:"column:mqtt_off"`
	Location    *string    `gorm:"column:location"`
	Unit        *string    `gorm:"column:unit"`
	LastValue   *string    `gorm:"column:last_value"`
	LastStatus  *int       `gorm:"column:last_status"`
	LastUpdated *time.Time `gorm:"column:last_updated"`
	IsOnline    bool       `gorm:"column:is_online;not null;default:false"`
	IsActive    bool       `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;default:now()"`

	// Relationships many-to-many with NotificationGroup
	NotificationGroups []NotificationGroup `gorm:"many2many:notification_groups_devices_notification_devices;"`
}

func (NotificationDevice) TableName() string { return "notification_devices" }

// NotificationGroup groups notification devices.
// กลุ่มของอุปกรณ์แจ้งเตือน
type NotificationGroup struct {
	ID        int             `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string          `gorm:"column:name;not null"`
	DeviceIDs json.RawMessage `gorm:"column:device_ids;type:json"`
	IsActive  bool            `gorm:"column:isActive;not null;default:true"`
	CreatedAt time.Time       `gorm:"column:createdAt;not null;default:now()"`
	UpdatedAt time.Time       `gorm:"column:updatedAt;not null;default:now()"`

	Devices []NotificationDevice `gorm:"many2many:notification_groups_devices_notification_devices;"`
}

func (NotificationGroup) TableName() string { return "notification_groups" }

// NotificationGroupsDevicesNotificationDevice is the join table for many-to-many.
// ตารางเชื่อม many-to-many ระหว่าง NotificationGroup และ NotificationDevice
type NotificationGroupsDevicesNotificationDevice struct {
	NotificationGroupsID  int `gorm:"column:notificationGroupsId;not null"`
	NotificationDevicesID int `gorm:"column:notificationDevicesId;not null"`
}

func (NotificationGroupsDevicesNotificationDevice) TableName() string {
	return "notification_groups_devices_notification_devices"
}

// NotificationLog logs device notifications.
// บันทึกการแจ้งเตือนของอุปกรณ์
type NotificationLog struct {
	ID                 int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID           int             `gorm:"column:device_id;not null"`
	DeviceName         string          `gorm:"column:device_name;not null"`
	DeviceType         string          `gorm:"column:device_type;not null"`
	ValueData          string          `gorm:"column:value_data;not null"`
	NumericValue       *float64        `gorm:"column:numeric_value;type:float8"`
	NotificationTypeID int             `gorm:"column:notification_type_id;not null"`
	Status             int             `gorm:"column:status;not null"`
	Title              string          `gorm:"column:title;not null"`
	Message            string          `gorm:"column:message;not null"`
	ChannelsSent       json.RawMessage `gorm:"column:channels_sent;type:jsonb;not null"`
	ControlAction      json.RawMessage `gorm:"column:control_action;type:jsonb"`
	RedisKey           *string         `gorm:"column:redis_key"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt          time.Time       `gorm:"column:updated_at;not null;default:now()"`
	ConfigID           *int            `gorm:"column:config_id"`
}

func (NotificationLog) TableName() string { return "notification_logs" }

// NotificationType defines basic notification types.
// ประเภทการแจ้งเตือนพื้นฐาน
type NotificationType struct {
	ID             int     `gorm:"column:id;primaryKey;autoIncrement"`
	Name           string  `gorm:"column:name;not null;unique"`
	Code           string  `gorm:"column:code;not null;unique"`
	Description    string  `gorm:"column:description;not null"`
	Icon           *string `gorm:"column:icon"`
	Color          *string `gorm:"column:color"`
	RepeatCooldown int     `gorm:"column:repeat_cooldown;not null;default:10"`
	ShouldNotify   bool    `gorm:"column:should_notify;not null;default:true"`
}

func (NotificationType) TableName() string { return "notification_types" }

// ============================================================================
// SD Tables (System Definition tables)
// ============================================================================

// SdActivityLog logs system activities with module references.
// บันทึกกิจกรรมของระบบ (อ้างอิงโมดูล)
type SdActivityLog struct {
	ID           int                `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       *string            `gorm:"column:user_id"`
	TypeID       *int               `gorm:"column:type_id"`
	ModulesID    *int               `gorm:"column:modules_id"`
	Name         *string            `gorm:"column:name"`
	Event        *string            `gorm:"column:event;type:text"`
	Detail       *string            `gorm:"column:detail;type:text"`
	Location     *string            `gorm:"column:location;type:text"`
	Date         time.Time          `gorm:"column:date;not null;default:now()"`
	ActivityType *SdActivityTypeLog `gorm:"-"`
	Module       *SdModuleLog       `gorm:"-"`
}

func (SdActivityLog) TableName() string { return "sd_activity_log" }

// SdActivityTypeLog defines activity types.
// ประเภทของกิจกรรม
type SdActivityTypeLog struct {
	TypeID   int    `gorm:"column:typeId;primaryKey;autoIncrement"`
	TypeName string `gorm:"column:type_name;not null"`

	// has many activity logs (แต่ไม่ต้องสร้าง foreign key อัตโนมัติ เพราะ FK อยู่ที่ SdActivityLog แล้ว)
	ActivityLogs []SdActivityLog `gorm:"foreignKey:TypeID;constraint:false"`
}

func (SdActivityTypeLog) TableName() string { return "sd_activity_type_log" }

// SdAdminAccessMenu defines admin menu access rights.
// สิทธิ์การเข้าถึงเมนูสำหรับผู้ดูแล
type SdAdminAccessMenu struct {
	AdminAccessID int  `gorm:"column:admin_access_id;primaryKey;autoIncrement"`
	AdminTypeID   *int `gorm:"column:admin_type_id"`
	AdminMenuID   *int `gorm:"column:admin_menu_id"`
}

func (SdAdminAccessMenu) TableName() string { return "sd_admin_access_menu" }

// SdAirControl controls air conditioning settings.
// การควบคุมระบบปรับอากาศ
type SdAirControl struct {
	AirControlID int       `gorm:"column:air_control_id;primaryKey;autoIncrement"`
	Name         *string   `gorm:"column:name"`
	Data         *string   `gorm:"column:data"`
	Status       *string   `gorm:"column:status"`
	Active       *int      `gorm:"column:active"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`

	// Relationships
	Devices []SdIotDevice `gorm:"many2many:sd_air_control_device_map;"`
}

func (SdAirControl) TableName() string { return "sd_air_control" }

// SdAirControlDeviceMap is the join table for SdAirControl and SdIotDevice.
// ตารางเชื่อมระหว่าง SdAirControl และ SdIotDevice
type SdAirControlDeviceMap struct {
	ID           string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AirControlID *int   `gorm:"column:air_control_id"`
	DeviceID     *int   `gorm:"column:device_id"`
}

func (SdAirControlDeviceMap) TableName() string { return "sd_air_control_device_map" }

// SdAirControlLog logs air control actions.
// บันทึกการทำงานของระบบปรับอากาศ
type SdAirControlLog struct {
	ID                  string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID       *int      `gorm:"column:alarm_action_id"`
	AirControlID        *int      `gorm:"column:air_control_id"`
	DeviceID            *int      `gorm:"column:device_id"`
	TypeID              *int      `gorm:"column:type_id"`
	Temperature         *string   `gorm:"column:temperature"`
	Warning             *string   `gorm:"column:warning"`
	Recovery            *string   `gorm:"column:recovery"`
	Period              *string   `gorm:"column:period"`
	Percent             *string   `gorm:"column:percent"`
	Firealarm           *string   `gorm:"column:firealarm"`
	Humidityalarm       *string   `gorm:"column:humidityalarm"`
	Air2Alarm           *string   `gorm:"column:air2_alarm"`
	Air1Alarm           *string   `gorm:"column:air1_alarm"`
	Temperaturealarm    *string   `gorm:"column:temperaturealarm"`
	Mode                *string   `gorm:"column:mode"`
	StateAir1           *string   `gorm:"column:state_air1"`
	StateAir2           *string   `gorm:"column:state_air2"`
	Temperaturealarmoff *string   `gorm:"column:temperaturealarmoff"`
	UpsAlarm            *string   `gorm:"column:ups_alarm"`
	Ups2Alarm           *string   `gorm:"column:ups2_alarm"`
	Hssdalarm           *string   `gorm:"column:hssdalarm"`
	Waterleakalarm      *string   `gorm:"column:waterleakalarm"`
	Date                string    `gorm:"column:date;not null"`
	Time                string    `gorm:"column:time;not null"`
	Data                *string   `gorm:"column:data"`
	Status              *string   `gorm:"column:status"`
	CreatedDate         time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate         time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAirControlLog) TableName() string { return "sd_air_control_log" }

// SdAirMod defines air conditioner modes.
// โหมดเครื่องปรับอากาศ
type SdAirMod struct {
	AirModID    int       `gorm:"column:air_mod_id;primaryKey;autoIncrement"`
	Name        *string   `gorm:"column:name"`
	Data        *string   `gorm:"column:data"`
	Status      *string   `gorm:"column:status"`
	Active      *int      `gorm:"column:active"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAirMod) TableName() string { return "sd_air_mod" }

// SdAirModDeviceMap join table for SdAirMod and devices.
// ตารางเชื่อม SdAirMod และอุปกรณ์
type SdAirModDeviceMap struct {
	ID           string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AirModID     *int   `gorm:"column:air_mod_id"`
	AirControlID *int   `gorm:"column:air_control_id"`
	DeviceID     *int   `gorm:"column:device_id"`
}

func (SdAirModDeviceMap) TableName() string { return "sd_air_mod_device_map" }

// SdAirPeriod defines time periods for air control.
// ช่วงเวลาสำหรับการควบคุมอากาศ
type SdAirPeriod struct {
	AirPeriodID int       `gorm:"column:air_period_id;primaryKey;autoIncrement"`
	Name        *string   `gorm:"column:name"`
	Data        *string   `gorm:"column:data"`
	Status      *string   `gorm:"column:status"`
	Active      *int      `gorm:"column:active"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAirPeriod) TableName() string { return "sd_air_period" }

// SdAirPeriodDeviceMap join table for SdAirPeriod and devices.
// ตารางเชื่อม SdAirPeriod และอุปกรณ์
type SdAirPeriodDeviceMap struct {
	ID           string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AirPeriodID  *int   `gorm:"column:air_period_id"`
	AirControlID *int   `gorm:"column:air_control_id"`
	DeviceID     *int   `gorm:"column:device_id"`
}

func (SdAirPeriodDeviceMap) TableName() string { return "sd_air_period_device_map" }

// SdAirSettingWarning defines warning settings for air control.
// การตั้งค่าเตือนสำหรับระบบอากาศ
type SdAirSettingWarning struct {
	AirSettingWarningID int       `gorm:"column:air_setting_warning_id;primaryKey;autoIncrement"`
	TypeID              *int      `gorm:"column:type_id"`
	DeviceID            *int      `gorm:"column:device_id"`
	PeriodID            *int      `gorm:"column:period_id"`
	EventName           *string   `gorm:"column:event_name"`
	Date                string    `gorm:"column:date;not null"`
	Time                string    `gorm:"column:time;not null"`
	Data                *string   `gorm:"column:data"`
	Status              *string   `gorm:"column:status"`
	Active              *int      `gorm:"column:active"`
	CreatedDate         time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate         time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAirSettingWarning) TableName() string { return "sd_air_setting_warning" }

// SdAirSettingWarningDeviceMap join table.
// ตารางเชื่อม
type SdAirSettingWarningDeviceMap struct {
	ID                  string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AirSettingWarningID *int   `gorm:"column:air_setting_warning_id"`
	AirControlID        *int   `gorm:"column:air_control_id"`
	DeviceID            *int   `gorm:"column:device_id"`
}

func (SdAirSettingWarningDeviceMap) TableName() string { return "sd_air_setting_warning_device_map" }

// SdAirWarning defines warnings for air systems.
// คำเตือนสำหรับระบบอากาศ
type SdAirWarning struct {
	AirWarningID int       `gorm:"column:air_warning_id;primaryKey;autoIncrement"`
	Name         *string   `gorm:"column:name"`
	Data         *string   `gorm:"column:data"`
	Status       *string   `gorm:"column:status"`
	Active       *int      `gorm:"column:active"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAirWarning) TableName() string { return "sd_air_warning" }

// SdAirWarningDeviceMap join table.
// ตารางเชื่อม
type SdAirWarningDeviceMap struct {
	ID           string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AirWarningID *int   `gorm:"column:air_warning_id"`
	AirControlID *int   `gorm:"column:air_control_id"`
	DeviceID     *int   `gorm:"column:device_id"`
}

func (SdAirWarningDeviceMap) TableName() string { return "sd_air_warning_device_map" }

// SdAlarmProcessLog logs alarm processing.
// บันทึกการประมวลผลแจ้งเตือน
type SdAlarmProcessLog struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLog) TableName() string { return "sd_alarm_process_log" }

// SdAlarmProcessLogEmail logs email alarms.
// บันทึกการแจ้งเตือนทางอีเมล
type SdAlarmProcessLogEmail struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogEmail) TableName() string { return "sd_alarm_process_log_email" }

// SdAlarmProcessLogLine logs LINE notifications.
// บันทึกการแจ้งเตือนทาง LINE
type SdAlarmProcessLogLine struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogLine) TableName() string { return "sd_alarm_process_log_line" }

// SdAlarmProcessLogMqtt logs MQTT alarms.
// บันทึกการแจ้งเตือนทาง MQTT
type SdAlarmProcessLogMqtt struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogMqtt) TableName() string { return "sd_alarm_process_log_mqtt" }

// SdAlarmProcessLogSms logs SMS alarms.
// บันทึกการแจ้งเตือนทาง SMS
type SdAlarmProcessLogSms struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogSms) TableName() string { return "sd_alarm_process_log_sms" }

// SdAlarmProcessLogTelegram logs Telegram alarms.
// บันทึกการแจ้งเตือนทาง Telegram
type SdAlarmProcessLogTelegram struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogTelegram) TableName() string { return "sd_alarm_process_log_telegram" }

// SdAlarmProcessLogTemp temporary alarm log.
// บันทึกการแจ้งเตือนชั่วคราว
type SdAlarmProcessLogTemp struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID   *int      `gorm:"column:alarm_action_id"`
	DeviceID        *int      `gorm:"column:device_id"`
	TypeID          *int      `gorm:"column:type_id"`
	Event           *string   `gorm:"column:event"`
	AlarmType       *string   `gorm:"column:alarm_type"`
	StatusWarning   *string   `gorm:"column:status_warning"`
	RecoveryWarning *string   `gorm:"column:recovery_warning"`
	StatusAlert     *string   `gorm:"column:status_alert"`
	RecoveryAlert   *string   `gorm:"column:recovery_alert"`
	EmailAlarm      *int      `gorm:"column:email_alarm"`
	LineAlarm       *int      `gorm:"column:line_alarm"`
	TelegramAlarm   *int      `gorm:"column:telegram_alarm"`
	SmsAlarm        *int      `gorm:"column:sms_alarm"`
	NoncAlarm       *int      `gorm:"column:nonc_alarm"`
	Status          *string   `gorm:"column:status"`
	Date            string    `gorm:"column:date;not null"`
	Time            string    `gorm:"column:time;not null"`
	Data            *string   `gorm:"column:data"`
	DataAlarm       *string   `gorm:"column:data_alarm"`
	AlarmStatus     *string   `gorm:"column:alarm_status"`
	Subject         *string   `gorm:"column:subject"`
	Content         *string   `gorm:"column:content"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdAlarmProcessLogTemp) TableName() string { return "sd_alarm_process_log_temp" }

// SdApiKey stores API keys for external access.
// คีย์ API สำหรับการเข้าถึงจากภายนอก
type SdApiKey struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string          `gorm:"column:name;not null"`
	Description *string         `gorm:"column:description;type:text"`
	ApiKey      string          `gorm:"column:api_key;not null;unique"`
	ApiSecret   string          `gorm:"column:api_secret;not null"`
	UserID      *string         `gorm:"column:user_id"`
	Permissions json.RawMessage `gorm:"column:permissions;type:jsonb"`
	ExpiresAt   *time.Time      `gorm:"column:expires_at"`
	LastUsedAt  *time.Time      `gorm:"column:last_used_at"`
	UsageCount  int             `gorm:"column:usage_count;not null;default:0"`
	IsActive    bool            `gorm:"column:is_active;not null;default:true"`
	IPWhitelist json.RawMessage `gorm:"column:ip_whitelist;type:jsonb"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;not null;default:now()"`
}

func (SdApiKey) TableName() string { return "sd_api_key" }

// SdAuditLog audits data changes.
// บันทึกการตรวจสอบการเปลี่ยนแปลงข้อมูล
type SdAuditLog struct {
	AuditID     int             `gorm:"column:audit_id;primaryKey;autoIncrement"`
	UserID      *string         `gorm:"column:user_id"`
	UserName    *string         `gorm:"column:user_name"`
	Action      string          `gorm:"column:action;not null"`
	EntityType  string          `gorm:"column:entity_type;not null"`
	EntityID    int             `gorm:"column:entity_id;not null"`
	Before      json.RawMessage `gorm:"column:before;type:jsonb"`
	After       json.RawMessage `gorm:"column:after;type:jsonb"`
	Changes     json.RawMessage `gorm:"column:changes;type:jsonb"`
	IPAddress   *string         `gorm:"column:ip_address"`
	UserAgent   *string         `gorm:"column:user_agent;type:text"`
	ActionTime  time.Time       `gorm:"column:action_time;not null;default:now()"`
	Description *string         `gorm:"column:description;type:text"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:now()"`
}

func (SdAuditLog) TableName() string { return "sd_audit_log" }

// AuditLogEntry is a hash-chained tamper-evident audit record (Phase 5).
// Each entry links to the previous entry via PrevHash forming a hash chain.
type AuditLogEntry struct {
	ID         string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(100)" json:"user_id,omitempty"`
	Action     string    `gorm:"column:action;type:varchar(100);not null" json:"action"`
	Resource   string    `gorm:"column:resource;type:varchar(200)" json:"resource,omitempty"`
	ResourceID string    `gorm:"column:resource_id;type:varchar(200)" json:"resource_id,omitempty"`
	RequestID  string    `gorm:"column:request_id;type:varchar(100)" json:"request_id,omitempty"`
	Before     string    `gorm:"column:before;type:text" json:"before,omitempty"`
	After      string    `gorm:"column:after;type:text" json:"after,omitempty"`
	Hash       string    `gorm:"column:hash;type:varchar(64);not null" json:"hash"`
	PrevHash   string    `gorm:"column:prev_hash;type:varchar(64)" json:"prev_hash,omitempty"`
	Nonce      int64     `gorm:"column:nonce;not null" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
}

func (AuditLogEntry) TableName() string { return "audit_logs" }

// BlockchainAnchor is an off-chain anchoring record for a hash-chain root.
// TxHash stays nil when using the default no-op anchor provider.
type BlockchainAnchor struct {
	ID         string     `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ChainID    string     `gorm:"column:chain_id;type:varchar(100);not null" json:"chain_id"`
	Payload    string     `gorm:"column:payload;type:text" json:"payload,omitempty"`
	Hash       string     `gorm:"column:hash;type:varchar(64);not null" json:"hash"`
	PrevHash   string     `gorm:"column:prev_hash;type:varchar(64)" json:"prev_hash,omitempty"`
	AnchoredAt *time.Time `gorm:"column:anchored_at" json:"anchored_at,omitempty"`
	TxHash     *string    `gorm:"column:tx_hash;type:varchar(200)" json:"tx_hash,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
}

func (BlockchainAnchor) TableName() string { return "blockchain_anchor" }

// FlowDefinition is a persisted Node-RED-like workflow (Phase 7).
// Nodes/Edges are stored as JSONB for graph flexibility.
type FlowDefinition struct {
	ID        string          `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string          `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Nodes     json.RawMessage `gorm:"column:nodes;type:jsonb" json:"nodes"`
	Edges     json.RawMessage `gorm:"column:edges;type:jsonb" json:"edges"`
	Enabled   bool            `gorm:"column:enabled;default:true" json:"enabled"`
	CreatedAt time.Time       `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (FlowDefinition) TableName() string { return "flow_definitions" }



// SdChannelTemplate defines message templates for notification channels.
// เทมเพลตข้อความสำหรับช่องทางการแจ้งเตือน
type SdChannelTemplate struct {
	ID                 int             `gorm:"column:id;primaryKey;autoIncrement"`
	Name               string          `gorm:"column:name;not null"`
	Description        *string         `gorm:"column:description;type:text"`
	ChannelID          int             `gorm:"column:channel_id;not null"`
	NotificationTypeID int             `gorm:"column:notification_type_id;not null"`
	Template           string          `gorm:"column:template;type:text;not null"`
	Variables          json.RawMessage `gorm:"column:variables;type:jsonb"`
	IsActive           bool            `gorm:"column:is_active;not null;default:true"`
	IsDefault          bool            `gorm:"column:is_default;not null;default:false"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt          time.Time       `gorm:"column:updated_at;not null;default:now()"`

	// Relationships
	Channel          *SdNotificationChannel `gorm:"foreignKey:ChannelID"`
	NotificationType *SdNotificationType    `gorm:"foreignKey:NotificationTypeID"`
}

func (SdChannelTemplate) TableName() string { return "sd_channel_template" }

// SdDashboardConfig stores dashboard configurations.
// การตั้งค่าแดชบอร์ด
type SdDashboardConfig struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	LocationID  int             `gorm:"column:location_id;not null"`
	Name        string          `gorm:"column:name;not null"`
	ConfigData  json.RawMessage `gorm:"column:config_data;type:json;not null"`
	Status      int             `gorm:"column:status;not null;default:1"`
	CreatedDate time.Time       `gorm:"column:created_date;not null;default:now()"`
	UpdatedDate time.Time       `gorm:"column:updated_date;not null;default:now()"`
}

func (SdDashboardConfig) TableName() string { return "sd_dashboard_config" }

// SdDeviceCategory categorizes devices.
// หมวดหมู่อุปกรณ์
type SdDeviceCategory struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;not null"`
	Description *string   `gorm:"column:description;type:text"`
	Icon        *string   `gorm:"column:icon"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (SdDeviceCategory) TableName() string { return "sd_device_category" }

// SdDeviceGroup groups devices.
// กลุ่มอุปกรณ์
type SdDeviceGroup struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string          `gorm:"column:name;not null"`
	Description *string         `gorm:"column:description;type:text"`
	GroupType   string          `gorm:"column:group_type;not null;default:custom"`
	IsActive    bool            `gorm:"column:is_active;not null;default:true"`
	Config      json.RawMessage `gorm:"column:config;type:jsonb"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;not null;default:now()"`

	// Relationships
	Members []SdDeviceMember `gorm:"foreignKey:GroupID"`
}

func (SdDeviceGroup) TableName() string { return "sd_device_group" }

// SdDeviceLog logs device activities.
// บันทึกกิจกรรมของอุปกรณ์
type SdDeviceLog struct {
	ID       int       `gorm:"column:id;primaryKey;autoIncrement"`
	TypeID   int       `gorm:"column:type_id;not null"`
	SensorID int       `gorm:"column:sensor_id;not null"`
	Name     string    `gorm:"column:name;not null"`
	Data     string    `gorm:"column:data;not null"`
	Status   *int      `gorm:"column:status"`
	Lang     *string   `gorm:"column:lang"`
	Create   time.Time `gorm:"column:create;not null;default:now()"`
}

func (SdDeviceLog) TableName() string { return "sd_device_log" }

// SdDeviceMember links devices to groups.
// เชื่อมอุปกรณ์กับกลุ่ม
type SdDeviceMember struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID  int       `gorm:"column:Device_id;not null"`
	GroupID   int       `gorm:"column:group_id;not null"`
	Role      string    `gorm:"column:role;not null;default:member"`
	Priority  int       `gorm:"column:priority;not null;default:1"`
	IsActive  bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`

	// Relationships
	Device *SdIotDevice   `gorm:"foreignKey:DeviceID"`
	Group  *SdDeviceGroup `gorm:"foreignKey:GroupID"`
}

func (SdDeviceMember) TableName() string { return "sd_device_member" }

// SdDeviceNotificationConfig configures device notification settings.
// การตั้งค่าการแจ้งเตือนของอุปกรณ์
type SdDeviceNotificationConfig struct {
	ID                    int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID              int             `gorm:"column:device_id;not null"`
	NotificationChannelID int             `gorm:"column:notification_channel_id;not null"`
	NotificationTypeID    int             `gorm:"column:notification_type_id;not null"`
	Config                json.RawMessage `gorm:"column:config;type:jsonb"`
	IsActive              bool            `gorm:"column:is_active;not null;default:true"`
	RetryCount            int             `gorm:"column:retry_count;not null;default:3"`
	RetryDelayMinutes     int             `gorm:"column:retry_delay_minutes;not null;default:5"`
	CreatedAt             time.Time       `gorm:"column:created_at;not null;default:now()"`

	// Relationships
	Device  *SdIotDevice           `gorm:"foreignKey:DeviceID"`
	Channel *SdNotificationChannel `gorm:"foreignKey:NotificationChannelID"`
	Type    *SdNotificationType    `gorm:"foreignKey:NotificationTypeID"`
}

func (SdDeviceNotificationConfig) TableName() string { return "sd_device_notification_config" }

// SdDeviceSchedule schedules device actions.
// ตารางเวลาการทำงานของอุปกรณ์
type SdDeviceSchedule struct {
	ID             int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID       int             `gorm:"column:device_id;not null"`
	Name           string          `gorm:"column:name;not null"`
	Description    *string         `gorm:"column:description;type:text"`
	ScheduleType   string          `gorm:"column:schedule_type;not null"`
	ScheduleConfig json.RawMessage `gorm:"column:schedule_config;type:jsonb;not null"`
	Action         json.RawMessage `gorm:"column:action;type:jsonb;not null"`
	IsActive       bool            `gorm:"column:is_active;not null;default:true"`
	LastRunAt      *time.Time      `gorm:"column:last_run_at"`
	NextRunAt      *time.Time      `gorm:"column:next_run_at"`
	RunCount       int             `gorm:"column:run_count;not null;default:0"`
	CreatedAt      time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;not null;default:now()"`

	Device *SdIotDevice `gorm:"foreignKey:DeviceID"`
}

func (SdDeviceSchedule) TableName() string { return "sd_device_schedule" }

// SdDeviceStatusHistory tracks device status changes over time.
// ประวัติการเปลี่ยนแปลงสถานะอุปกรณ์
type SdDeviceStatusHistory struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID           int       `gorm:"column:device_id;not null"`
	Status             *string   `gorm:"column:status"`
	Value              *float64  `gorm:"column:value;type:numeric(10,2)"`
	NotificationTypeID *int      `gorm:"column:notification_type_id"`
	DurationMinutes    *int      `gorm:"column:duration_minutes"`
	PreviousStatus     *string   `gorm:"column:previous_status"`
	PreviousValue      *float64  `gorm:"column:previous_value;type:numeric(10,2)"`
	ChangeReason       *string   `gorm:"column:change_reason;type:text"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()"`

	Device *SdIotDevice        `gorm:"foreignKey:DeviceID"`
	Type   *SdNotificationType `gorm:"foreignKey:NotificationTypeID"`
}

func (SdDeviceStatusHistory) TableName() string { return "sd_device_status_history" }

// SdGroupNotificationConfig configures notification settings for device groups.
// การตั้งค่าการแจ้งเตือนสำหรับกลุ่มอุปกรณ์
type SdGroupNotificationConfig struct {
	ID                     int             `gorm:"column:id;primaryKey;autoIncrement"`
	GroupID                int             `gorm:"column:group_id;not null"`
	NotificationChannelID  int             `gorm:"column:notification_channel_id;not null"`
	NotificationTypeID     int             `gorm:"column:notification_type_id;not null"`
	Config                 json.RawMessage `gorm:"column:config;type:jsonb"`
	IsActive               bool            `gorm:"column:is_active;not null;default:true"`
	EscalationLevel        int             `gorm:"column:escalation_level;not null;default:1"`
	EscalationDelayMinutes int             `gorm:"column:escalation_delay_minutes;not null;default:30"`
	CreatedAt              time.Time       `gorm:"column:created_at;not null;default:now()"`

	Group   *SdDeviceGroup         `gorm:"foreignKey:GroupID"`
	Channel *SdNotificationChannel `gorm:"foreignKey:NotificationChannelID"`
	Type    *SdNotificationType    `gorm:"foreignKey:NotificationTypeID"`
}

func (SdGroupNotificationConfig) TableName() string { return "sd_group_notification_config" }

// SdIotAlarmDevice links alarm actions to devices.
// เชื่อมการแจ้งเตือนกับอุปกรณ์
type SdIotAlarmDevice struct {
	ID            string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID *int   `gorm:"column:alarm_action_id"`
	DeviceID      *int   `gorm:"column:device_id"`
}

func (SdIotAlarmDevice) TableName() string { return "sd_iot_alarm_device" }

// SdIotAlarmDeviceEvent links alarm events to devices.
// เชื่อมเหตุการณ์แจ้งเตือนกับอุปกรณ์
type SdIotAlarmDeviceEvent struct {
	ID            string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	AlarmActionID *int   `gorm:"column:alarm_action_id"`
	DeviceID      *int   `gorm:"column:device_id"`
}

func (SdIotAlarmDeviceEvent) TableName() string { return "sd_iot_alarm_device_event" }

// SdIotApi stores API configurations.
// การตั้งค่า API
type SdIotApi struct {
	APIID       int       `gorm:"column:api_id;primaryKey;autoIncrement"`
	APIName     string    `gorm:"column:api_name;not null"`
	Host        *int      `gorm:"column:host"`
	Port        string    `gorm:"column:port;not null"`
	TokenValue  *string   `gorm:"column:token_value;type:text"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotApi) TableName() string { return "sd_iot_api" }

// SdIotDevice represents an IoT device.
// อุปกรณ์ IoT
type SdIotDevice struct {
	DeviceID             int       `gorm:"column:device_id;primaryKey;autoIncrement"`
	SettingID            *int      `gorm:"column:setting_id"`
	TypeID               *int      `gorm:"column:type_id"`
	LocationID           *int      `gorm:"column:location_id"`
	DeviceName           *string   `gorm:"column:device_name"`
	Sn                   *string   `gorm:"column:sn;unique"`
	HardwareID           *int      `gorm:"column:hardware_id"`
	StatusWarning        *string   `gorm:"column:status_warning"`
	RecoveryWarning      *string   `gorm:"column:recovery_warning"`
	StatusAlert          *string   `gorm:"column:status_alert"`
	RecoveryAlert        *string   `gorm:"column:recovery_alert"`
	TimeLife             *int      `gorm:"column:time_life;default:1"`
	Period               *string   `gorm:"column:period"`
	WorkStatus           *int      `gorm:"column:work_status;default:1"`
	Model                *string   `gorm:"column:model"`
	Vendor               *string   `gorm:"column:vendor"`
	Comparevalue         *string   `gorm:"column:comparevalue"`
	Unit                 *string   `gorm:"column:unit"`
	MqttID               *int      `gorm:"column:mqtt_id"`
	Oid                  *string   `gorm:"column:oid"`
	ActionID             *int      `gorm:"column:action_id"`
	StatusAlertID        *int      `gorm:"column:status_alert_id"`
	MqttDataValue        *string   `gorm:"column:mqtt_data_value"`
	MqttDataControl      *string   `gorm:"column:mqtt_data_control"`
	Measurement          *string   `gorm:"column:measurement"`
	MqttControlOn        *string   `gorm:"column:mqtt_control_on;default:1"`
	MqttControlOff       *string   `gorm:"column:mqtt_control_off;default:0"`
	Org                  string    `gorm:"column:org;not null"`
	Bucket               string    `gorm:"column:bucket;not null"`
	Status               *int      `gorm:"column:status"`
	MqttDeviceName       string    `gorm:"column:mqtt_device_name;not null"`
	MqttStatusOverName   *string   `gorm:"column:mqtt_status_over_name;type:text"`
	MqttStatusDataName   *string   `gorm:"column:mqtt_status_data_name;type:text"`
	MqttActRelayName     *string   `gorm:"column:mqtt_act_relay_name;type:text"`
	MqttControlRelayName *string   `gorm:"column:mqtt_control_relay_name;type:text"`
	MqttConfig           *string   `gorm:"column:mqtt_config;type:text"`
	CreatedDate          time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate          time.Time `gorm:"column:updateddate;not null;default:now()"`
	Max                  *string   `gorm:"column:max"`
	Min                  *string   `gorm:"column:min"`
	Layout               *int      `gorm:"column:layout;default:1"`
	AlertSet             *int      `gorm:"column:alert_set;default:1"`
	IconNormal           string    `gorm:"column:icon_normal;type:text;not null;default:'<svg...>'"`
	IconWarning          string    `gorm:"column:icon_warning;type:text;not null;default:'<svg...>'"`
	IconAlert            string    `gorm:"column:icon_alert;type:text;not null;default:'<svg...>'"`
	Icon                 string    `gorm:"column:icon;type:text;not null;default:'<svg...>'"`
	ColorNormal          string    `gorm:"column:color_normal;not null;default:'#22C55E'"`
	ColorWarning         string    `gorm:"column:color_warning;not null;default:'#F59E0B'"`
	ColorAlarm           string    `gorm:"column:color_alarm;not null;default:'#EF4444'"`
	Code                 string    `gorm:"column:code;not null;default:'normal'"`
	Menu                 *int      `gorm:"column:menu;default:1"`
	IconOn               string    `gorm:"column:icon_on;type:text;not null;default:'<svg...>'"`
	IconOff              string    `gorm:"column:icon_off;type:text;not null;default:'<svg...>'"`
	CalibrationAdd       *string   `gorm:"column:calibration_add;default:'0'"`
	CalibrationSubtract  *string   `gorm:"column:calibration_subtract;default:'0'"`
	CalibrationType      *int      `gorm:"column:calibration_type;default:3"`

	// Relationships
	Location *SdIotLocation `gorm:"foreignKey:LocationID"`
	Mqtt     *SdIotMqtt     `gorm:"foreignKey:MqttID"`
}

func (SdIotDevice) TableName() string { return "sd_iot_device" }

// SdIotDeviceAction links device actions.
// เชื่อมการทำงานของอุปกรณ์
type SdIotDeviceAction struct {
	DeviceActionUserID int  `gorm:"column:device_action_user_id;primaryKey;autoIncrement"`
	AlarmActionID      *int `gorm:"column:alarm_action_id"`
	DeviceID           *int `gorm:"column:device_id"`
}

func (SdIotDeviceAction) TableName() string { return "sd_iot_device_action" }

// SdIotDeviceActionLog logs device action executions.
// บันทึกการทำงานของอุปกรณ์
type SdIotDeviceActionLog struct {
	LogID         int       `gorm:"column:log_id;primaryKey;autoIncrement"`
	AlarmActionID *int      `gorm:"column:alarm_action_id"`
	DeviceID      *int      `gorm:"column:device_id"`
	UID           *string   `gorm:"column:uid"`
	Status        *int      `gorm:"column:status"`
	CreatedDate   time.Time `gorm:"column:createddate;not null;default:now()"`
}

func (SdIotDeviceActionLog) TableName() string { return "sd_iot_device_action_log" }

// SdIotDeviceActionUser links actions to users.
// เชื่อมการทำงานกับผู้ใช้
type SdIotDeviceActionUser struct {
	DeviceActionUserID int     `gorm:"column:device_action_user_id;primaryKey;autoIncrement"`
	AlarmActionID      *int    `gorm:"column:alarm_action_id"`
	UID                *string `gorm:"column:uid"`
}

func (SdIotDeviceActionUser) TableName() string { return "sd_iot_device_action_user" }

// SdIotDeviceAlarmAction defines alarm actions for devices.
// การดำเนินการเมื่อเกิดการแจ้งเตือนของอุปกรณ์
type SdIotDeviceAlarmAction struct {
	AlarmActionID   int     `gorm:"column:alarm_action_id;primaryKey;autoIncrement"`
	ActionName      *string `gorm:"column:action_name"`
	StatusWarning   *string `gorm:"column:status_warning"`
	RecoveryWarning *string `gorm:"column:recovery_warning"`
	StatusAlert     *string `gorm:"column:status_alert"`
	RecoveryAlert   *string `gorm:"column:recovery_alert"`
	EmailAlarm      *int    `gorm:"column:email_alarm"`
	LineAlarm       *int    `gorm:"column:line_alarm"`
	TelegramAlarm   *int    `gorm:"column:telegram_alarm"`
	SmsAlarm        *int    `gorm:"column:sms_alarm"`
	NoncAlarm       *int    `gorm:"column:nonc_alarm"`
	TimeLife        *int    `gorm:"column:time_life"`
	Event           *int    `gorm:"column:event"`
	Status          *int    `gorm:"column:status"`
}

func (SdIotDeviceAlarmAction) TableName() string { return "sd_iot_device_alarm_action" }

// SdIotDeviceType defines device types.
// ประเภทอุปกรณ์
type SdIotDeviceType struct {
	TypeID      int       `gorm:"column:type_id;primaryKey;autoIncrement"`
	TypeName    string    `gorm:"column:type_name;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotDeviceType) TableName() string { return "sd_iot_device_type" }

// SdIotEmail stores email server configurations.
// การตั้งค่าเซิร์ฟเวอร์อีเมล
type SdIotEmail struct {
	EmailID     string    `gorm:"column:email_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	EmailName   string    `gorm:"column:email_name;not null"`
	Host        string    `gorm:"column:host;not null"`
	Port        *int      `gorm:"column:port"`
	Username    string    `gorm:"column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotEmail) TableName() string { return "sd_iot_email" }

// SdIotGroup defines device groups.
// กลุ่มอุปกรณ์ IoT
type SdIotGroup struct {
	GroupID     int       `gorm:"column:group_id;primaryKey;autoIncrement"`
	GroupName   string    `gorm:"column:group_name;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotGroup) TableName() string { return "sd_iot_group" }

// SdIotHost stores host configurations.
// การตั้งค่าโฮสต์
type SdIotHost struct {
	HostID      string    `gorm:"column:host_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	HostName    string    `gorm:"column:host_name;not null"`
	Port        string    `gorm:"column:port;not null"`
	Username    string    `gorm:"column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
	IDHost      *int      `gorm:"column:idhost"`
}

func (SdIotHost) TableName() string { return "sd_iot_host" }

// SdIotInfluxdb stores InfluxDB configurations.
// การตั้งค่า InfluxDB
type SdIotInfluxdb struct {
	InfluxdbID   string    `gorm:"column:influxdb_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	InfluxdbName *string   `gorm:"column:influxdb_name;type:text"`
	Host         *string   `gorm:"column:host;type:text"`
	Port         string    `gorm:"column:port;not null"`
	Username     string    `gorm:"column:username;not null"`
	Password     string    `gorm:"column:password;not null"`
	TokenValue   *string   `gorm:"column:token_value;type:text"`
	Buckets      *string   `gorm:"column:buckets;type:text"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status       *int      `gorm:"column:status"`
}

func (SdIotInfluxdb) TableName() string { return "sd_iot_influxdb" }

// SdIotLine stores LINE API configurations.
// การตั้งค่า LINE API
type SdIotLine struct {
	LineID       string    `gorm:"column:line_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	LineName     *string   `gorm:"column:line_name;type:text"`
	ClientID     *string   `gorm:"column:client_id;type:text"`
	ClientSecret *string   `gorm:"column:client_secret;type:text"`
	SecretKey    *string   `gorm:"column:secret_key;type:text"`
	RedirectURI  *string   `gorm:"column:redirect_uri;type:text"`
	GrantType    string    `gorm:"column:grant_type;not null"`
	Code         string    `gorm:"column:code;not null"`
	Accesstoken  *string   `gorm:"column:accesstoken;type:text"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status       *int      `gorm:"column:status"`
}

func (SdIotLine) TableName() string { return "sd_iot_line" }

// SdIotLocation defines device locations.
// ตำแหน่งที่ตั้งอุปกรณ์
type SdIotLocation struct {
	LocationID     int       `gorm:"column:location_id;primaryKey;autoIncrement"`
	LocationName   string    `gorm:"column:location_name;not null"`
	IPAddress      string    `gorm:"column:ipaddress;not null"`
	LocationDetail string    `gorm:"column:location_detail;not null"`
	CreatedDate    time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate    time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status         *int      `gorm:"column:status"`
	Configdata     *string   `gorm:"column:configdata;type:text"`
}

func (SdIotLocation) TableName() string { return "sd_iot_location" }

// SdIotMqtt stores MQTT broker configurations.
// การตั้งค่า MQTT broker
type SdIotMqtt struct {
	MqttID        int       `gorm:"column:mqtt_id;primaryKey;autoIncrement"`
	MqttTypeID    *int      `gorm:"column:mqtt_type_id"`
	Sort          int       `gorm:"column:sort;not null;default:1"`
	MqttName      *string   `gorm:"column:mqtt_name"`
	Host          *string   `gorm:"column:host"`
	Port          *int      `gorm:"column:port"`
	Username      *string   `gorm:"column:username"`
	Password      *string   `gorm:"column:password"`
	Secret        *string   `gorm:"column:secret"`
	ExpireIn      *string   `gorm:"column:expire_in"`
	TokenValue    *string   `gorm:"column:token_value"`
	Org           *string   `gorm:"column:org"`
	Bucket        *string   `gorm:"column:bucket"`
	Envavorment   *string   `gorm:"column:envavorment"`
	CreatedDate   time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate   time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status        int       `gorm:"column:status;not null;default:1"`
	LocationID    *int      `gorm:"column:location_id;default:1"`
	Latitude      *string   `gorm:"column:latitude"`
	Longitude     *string   `gorm:"column:longitude"`
	MqttMainID    int       `gorm:"column:mqtt_main_id;not null;default:1"`
	Configuration *string   `gorm:"column:configuration;type:text;default:'{\"0\":\"temperature1\",\"1\":\"humidity1\"}'"`
	Zoom          int       `gorm:"column:zoom;not null;default:6"`
}

func (SdIotMqtt) TableName() string { return "sd_iot_mqtt" }

// SdIotNodered stores Node-RED configurations.
// การตั้งค่า Node-RED
type SdIotNodered struct {
	NoderedID   string    `gorm:"column:nodered_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	NoderedName *string   `gorm:"column:nodered_name;type:text"`
	Host        string    `gorm:"column:host;not null"`
	Port        string    `gorm:"column:port;not null"`
	Routing     *string   `gorm:"column:routing;type:text"`
	ClientID    *string   `gorm:"column:client_id;type:text"`
	GrantType   string    `gorm:"column:grant_type;not null"`
	Scope       string    `gorm:"column:scope;not null"`
	Username    string    `gorm:"column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotNodered) TableName() string { return "sd_iot_nodered" }

// SdIotSchedule defines scheduled tasks.
// งานตามตารางเวลา
type SdIotSchedule struct {
	ScheduleID   int       `gorm:"column:schedule_id;primaryKey;autoIncrement"`
	ScheduleName string    `gorm:"column:schedule_name;not null"`
	DeviceID     *int      `gorm:"column:device_id"`
	Start        string    `gorm:"column:start;not null"`
	Event        *int      `gorm:"column:event"`
	Sunday       *int      `gorm:"column:sunday"`
	Monday       *int      `gorm:"column:monday"`
	Tuesday      *int      `gorm:"column:tuesday"`
	Wednesday    *int      `gorm:"column:wednesday"`
	Thursday     *int      `gorm:"column:thursday"`
	Friday       *int      `gorm:"column:friday"`
	Saturday     *int      `gorm:"column:saturday"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status       *int      `gorm:"column:status"`
}

func (SdIotSchedule) TableName() string { return "sd_iot_schedule" }

// SdIotScheduleDevice links schedules to devices.
// เชื่อมตารางเวลากับอุปกรณ์
type SdIotScheduleDevice struct {
	ID         string `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	ScheduleID *int   `gorm:"column:schedule_id"`
	DeviceID   *int   `gorm:"column:device_id"`
}

func (SdIotScheduleDevice) TableName() string { return "sd_iot_schedule_device" }

// SdIotSensor represents IoT sensors.
// เซ็นเซอร์ IoT
type SdIotSensor struct {
	SensorID        int       `gorm:"column:sensor_id;primaryKey;autoIncrement"`
	SettingID       *int      `gorm:"column:setting_id"`
	SettingTypeID   *int      `gorm:"column:setting_type_id"`
	SensorName      string    `gorm:"column:sensor_name;not null"`
	Sn              string    `gorm:"column:sn;not null"`
	Max             string    `gorm:"column:max;not null"`
	Min             string    `gorm:"column:min;not null"`
	HardwareID      *int      `gorm:"column:hardware_id"`
	StatusHigh      string    `gorm:"column:status_high;not null"`
	StatusWarning   string    `gorm:"column:status_warning;not null"`
	StatusAlert     string    `gorm:"column:status_alert;not null"`
	Model           string    `gorm:"column:model;not null"`
	Vendor          string    `gorm:"column:vendor;not null"`
	Comparevalue    string    `gorm:"column:comparevalue;not null"`
	CreatedDate     time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate     time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status          *int      `gorm:"column:status"`
	Unit            string    `gorm:"column:unit;not null"`
	MqttID          *int      `gorm:"column:mqtt_id"`
	Oid             string    `gorm:"column:oid;not null"`
	ActionID        *int      `gorm:"column:action_id"`
	StatusAlertID   *int      `gorm:"column:status_alert_id"`
	MqttDataValue   string    `gorm:"column:mqtt_data_value;not null"`
	MqttDataControl string    `gorm:"column:mqtt_data_control;not null"`
}

func (SdIotSensor) TableName() string { return "sd_iot_sensor" }

// SdIotSetting stores IoT settings.
// การตั้งค่า IoT
type SdIotSetting struct {
	SettingID     int       `gorm:"column:setting_id;primaryKey;autoIncrement"`
	LocationID    *int      `gorm:"column:location_id"`
	SettingTypeID *int      `gorm:"column:setting_type_id"`
	SettingName   string    `gorm:"column:setting_name;not null"`
	Sn            string    `gorm:"column:sn;not null"`
	CreatedDate   time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate   time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status        *int      `gorm:"column:status"`
}

func (SdIotSetting) TableName() string { return "sd_iot_setting" }

// SdIotSms stores SMS gateway configurations.
// การตั้งค่า SMS gateway
type SdIotSms struct {
	SmsID       string    `gorm:"column:sms_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	SmsName     string    `gorm:"column:sms_name;not null"`
	Host        string    `gorm:"column:host;not null"`
	Port        *int      `gorm:"column:port"`
	Username    string    `gorm:"column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	Apikey      string    `gorm:"column:apikey;not null"`
	Originator  string    `gorm:"column:originator;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotSms) TableName() string { return "sd_iot_sms" }

// SdIotTelegram stores Telegram bot configurations.
// การตั้งค่า Telegram bot
type SdIotTelegram struct {
	TelegramID   string    `gorm:"column:telegram_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TelegramName string    `gorm:"column:telegram_name;not null"`
	Port         string    `gorm:"column:port;not null"`
	Username     string    `gorm:"column:username;not null"`
	Password     string    `gorm:"column:password;not null"`
	CreatedDate  time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate  time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status       *int      `gorm:"column:status"`
}

func (SdIotTelegram) TableName() string { return "sd_iot_telegram" }

// SdIotToken stores token configurations.
// การตั้งค่าโทเค็น
type SdIotToken struct {
	TokenID     int       `gorm:"column:token_id;primaryKey;autoIncrement"`
	TokenName   string    `gorm:"column:token_name;not null"`
	Host        *int      `gorm:"column:host"`
	Port        string    `gorm:"column:port;not null"`
	TokenValue  *string   `gorm:"column:token_value;type:text"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotToken) TableName() string { return "sd_iot_token" }

// SdIotType defines IoT types.
// ประเภท IoT
type SdIotType struct {
	TypeID      int       `gorm:"column:type_id;primaryKey;autoIncrement"`
	TypeName    string    `gorm:"column:type_name;not null"`
	GroupID     *int      `gorm:"column:group_id"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
}

func (SdIotType) TableName() string { return "sd_iot_type" }

// SdModuleLog defines system modules.
// โมดูลของระบบ
type SdModuleLog struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement"`
	ModuleName string `gorm:"column:module_name;not null"`
}

func (SdModuleLog) TableName() string { return "sd_module_log" }

// SdMqttHost stores MQTT host configurations.
// การตั้งค่าโฮสต์ MQTT
type SdMqttHost struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Hostname    string    `gorm:"column:hostname;not null"`
	Host        string    `gorm:"column:host;not null"`
	Port        string    `gorm:"column:port;not null"`
	Username    string    `gorm:"column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	Status      *int      `gorm:"column:status"`
	IDHost      *int      `gorm:"column:idhost"`
}

func (SdMqttHost) TableName() string { return "sd_mqtt_host" }

// SdMqttLog logs MQTT activities.
// บันทึกกิจกรรม MQTT
type SdMqttLog struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        *string   `gorm:"column:name"`
	Statusmqtt  *string   `gorm:"column:statusmqtt"`
	Msg         *string   `gorm:"column:msg"`
	TypeID      *int      `gorm:"column:type_id"`
	Date        string    `gorm:"column:date;not null"`
	Time        string    `gorm:"column:time;not null"`
	Data        *string   `gorm:"column:data"`
	Status      *string   `gorm:"column:status"`
	CreatedDate time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate time.Time `gorm:"column:updateddate;not null;default:now()"`
	DeviceID    *int      `gorm:"column:Device_id"`
	DeviceName  *string   `gorm:"column:Device_name"`
}

func (SdMqttLog) TableName() string { return "sd_mqtt_log" }

// SdNotificationChannel defines notification channels (email, LINE, etc.).
// ช่องทางการแจ้งเตือน
type SdNotificationChannel struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string    `gorm:"column:name;not null"`
	Description  *string   `gorm:"column:description;type:text"`
	Icon         *string   `gorm:"column:icon"`
	HandlerClass *string   `gorm:"column:handler_class"`
	IsActive     bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (SdNotificationChannel) TableName() string { return "sd_notification_channel" }

// SdNotificationCondition defines conditions for notifications.
// เงื่อนไขการแจ้งเตือน
type SdNotificationCondition struct {
	ID                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID           int       `gorm:"column:device_id;not null"`
	NotificationTypeID int       `gorm:"column:notification_type_id;not null"`
	ConditionOperator  string    `gorm:"column:condition_operator;not null;default:between"`
	Priority           int       `gorm:"column:priority;not null;default:1"`
	IsActive           bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()"`
	MinValue           *float64  `gorm:"column:minValue;type:numeric(10,2)"`
	MaxValue           *float64  `gorm:"column:maxValue;type:numeric(10,2)"`

	Device *SdIotDevice        `gorm:"foreignKey:DeviceID"`
	Type   *SdNotificationType `gorm:"foreignKey:NotificationTypeID"`
}

func (SdNotificationCondition) TableName() string { return "sd_notification_condition" }

// SdNotificationLog logs notification deliveries.
// บันทึกการส่งการแจ้งเตือนของระบบ
type SdNotificationLog struct {
	ID                    int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID              *int            `gorm:"column:device_id"`
	NotificationTypeID    *int            `gorm:"column:notification_type_id"`
	NotificationChannelID *int            `gorm:"column:notification_channel_id"`
	Message               string          `gorm:"column:message;type:text;not null"`
	ResponseData          json.RawMessage `gorm:"column:response_data;type:jsonb"`
	SentAt                *time.Time      `gorm:"column:sent_at"`
	CreatedAt             time.Time       `gorm:"column:created_at;not null;default:now()"`
	TemplateID            *int            `gorm:"column:template_id"`
	DeliveredAt           *time.Time      `gorm:"column:delivered_at"`
	ReadAt                *time.Time      `gorm:"column:read_at"`
	RetryCount            int             `gorm:"column:retry_count;not null;default:0"`
	ErrorMessage          *string         `gorm:"column:error_message;type:text"`
	MessageID             *string         `gorm:"column:message_id"`
	Recipient             *string         `gorm:"column:recipient"`
	Status                string          `gorm:"column:status;not null;default:pending"`

	Device  *SdIotDevice           `gorm:"foreignKey:DeviceID"`
	Type    *SdNotificationType    `gorm:"foreignKey:NotificationTypeID"`
	Channel *SdNotificationChannel `gorm:"foreignKey:NotificationChannelID"`
}

func (SdNotificationLog) TableName() string { return "sd_notification_log" }

// SdNotificationType defines notification types for SD system.
// ประเภทการแจ้งเตือนของระบบ SD
type SdNotificationType struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name            string    `gorm:"column:name;not null"`
	Description     *string   `gorm:"column:description;type:text"`
	CooldownMinutes int       `gorm:"column:cooldown_minutes;not null;default:10"`
	IsActive        bool      `gorm:"column:is_active;not null;default:true"`
	Icon            *string   `gorm:"column:icon"`
	Color           *string   `gorm:"column:color"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (SdNotificationType) TableName() string { return "sd_notification_type" }

// SdReportData stores generated reports.
// ข้อมูลรายงานที่สร้างขึ้น
type SdReportData struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID    int             `gorm:"column:device_id;not null"`
	ReportType  string          `gorm:"column:report_type;not null"`
	Data        json.RawMessage `gorm:"column:data;type:jsonb;not null"`
	PeriodStart time.Time       `gorm:"column:period_start;not null"`
	PeriodEnd   time.Time       `gorm:"column:period_end;not null"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:now()"`
	TemplateID  *int            `gorm:"column:template_id"`
	GeneratedAt time.Time       `gorm:"column:generated_at;not null;default:now()"`
	FilePath    *string         `gorm:"column:file_path"`
	FileFormat  *string         `gorm:"column:file_format"`
	IsExported  bool            `gorm:"column:is_exported;not null;default:false"`
	ExportedAt  *time.Time      `gorm:"column:exported_at"`

	Device *SdIotDevice `gorm:"foreignKey:DeviceID"`
}

func (SdReportData) TableName() string { return "sd_report_data" }

// SdScheduleProcessLog logs schedule executions.
// บันทึกการทำงานตามตารางเวลา
type SdScheduleProcessLog struct {
	ID                 string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	ScheduleID         *int      `gorm:"column:schedule_id"`
	DeviceID           *int      `gorm:"column:device_id"`
	ScheduleEventStart string    `gorm:"column:schedule_event_start;not null"`
	Day                string    `gorm:"column:day;not null"`
	Doday              string    `gorm:"column:doday;not null"`
	Dotime             string    `gorm:"column:dotime;not null"`
	ScheduleEvent      string    `gorm:"column:schedule_event;not null"`
	DeviceStatus       string    `gorm:"column:device_status;not null"`
	Status             string    `gorm:"column:status;not null"`
	Date               string    `gorm:"column:date;not null"`
	Time               string    `gorm:"column:time;not null"`
	CreatedDate        time.Time `gorm:"column:createddate;not null;default:now()"`
	UpdatedDate        time.Time `gorm:"column:updateddate;not null;default:now()"`
}

func (SdScheduleProcessLog) TableName() string { return "sd_schedule_process_log" }

// SdSensorData stores sensor readings.
// ค่าที่อ่านได้จากเซ็นเซอร์
type SdSensorData struct {
	ID                 int             `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID           int             `gorm:"column:device_id;not null"`
	Value              float64         `gorm:"column:value;type:numeric(10,2);not null"`
	RawData            json.RawMessage `gorm:"column:raw_data;type:jsonb"`
	NotificationTypeID *int            `gorm:"column:notification_type_id"`
	Timestamp          time.Time       `gorm:"column:timestamp;not null;default:now()"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null;default:now()"`
	BatteryLevel       *float64        `gorm:"column:battery_level;type:numeric(5,2)"`
	SignalStrength     *int            `gorm:"column:signal_strength"`

	Device *SdIotDevice        `gorm:"foreignKey:DeviceID"`
	Type   *SdNotificationType `gorm:"foreignKey:NotificationTypeID"`
}

func (SdSensorData) TableName() string { return "sd_sensor_data" }

// SdSystemSetting stores system-wide settings.
// การตั้งค่าระดับระบบ
type SdSystemSetting struct {
	ID          int             `gorm:"column:id;primaryKey;autoIncrement"`
	Key         string          `gorm:"column:key;not null;unique"`
	Value       json.RawMessage `gorm:"column:value;type:jsonb;not null"`
	Category    *string         `gorm:"column:category"`
	Description *string         `gorm:"column:description;type:text"`
	IsPublic    bool            `gorm:"column:is_public;not null;default:false"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;not null;default:now()"`
}

func (SdSystemSetting) TableName() string { return "sd_system_setting" }

// SdUserFile stores user file metadata.
// ข้อมูลไฟล์ของผู้ใช้
type SdUserFile struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	FileName   string    `gorm:"column:file_name;not null"`
	FileType   string    `gorm:"column:file_type;not null"`
	FilePath   string    `gorm:"column:file_path;not null"`
	FileTypeID int       `gorm:"column:file_type_id;not null"`
	UID        *string   `gorm:"column:uid"`
	FileDate   time.Time `gorm:"column:file_date;not null;default:now()"`
	Status     int16     `gorm:"column:status;not null"`
}

func (SdUserFile) TableName() string { return "sd_user_file" }

// SdUserLog logs user activities.
// บันทึกกิจกรรมของผู้ใช้
type SdUserLog struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	LogTypeID    int       `gorm:"column:log_type_id;not null"`
	UID          uuid.UUID `gorm:"column:uid;type:uuid;not null"`
	Name         string    `gorm:"column:name;not null"`
	Detail       string    `gorm:"column:detail;not null"`
	SelectStatus *int      `gorm:"column:select_status"`
	InsertStatus *int      `gorm:"column:insert_status"`
	UpdateStatus *int      `gorm:"column:update_status"`
	DeleteStatus *int      `gorm:"column:delete_status"`
	Status       *int      `gorm:"column:status"`
	Create       time.Time `gorm:"column:create;not null;default:now()"`
	Update       time.Time `gorm:"column:update;not null;default:now()"`
	Lang         *string   `gorm:"column:lang"`
}

func (SdUserLog) TableName() string { return "sd_user_log" }

// SdUserLogType defines user log types.
// ประเภทบันทึกกิจกรรมผู้ใช้
type SdUserLogType struct {
	LogTypeID  int       `gorm:"column:log_type_id;primaryKey;autoIncrement"`
	TypeName   string    `gorm:"column:type_name;not null"`
	TypeDetail string    `gorm:"column:type_detail;not null"`
	Status     *int      `gorm:"column:status"`
	Create     time.Time `gorm:"column:create;not null;default:now()"`
	Update     time.Time `gorm:"column:update;not null;default:now()"`
}

func (SdUserLogType) TableName() string { return "sd_user_log_type" }

// SdUserAccessMenu defines user menu access rights.
// สิทธิ์การเข้าถึงเมนูของผู้ใช้
type SdUserAccessMenu struct {
	UserAccessID int64  `gorm:"column:user_access_id;primaryKey;autoIncrement"`
	UserTypeID   *int64 `gorm:"column:user_type_id"` // nullable
	MenuID       *int64 `gorm:"column:menu_id"`      // nullable
	ParentID     *int64 `gorm:"column:parent_id"`    // nullable
}

func (SdUserAccessMenu) TableName() string {
	return "sd_user_access_menu"
}

// SdUserRolesPermision defines role permissions.
// สิทธิ์ของบทบาท
type SdUserRolesPermision struct {
	RoleTypeID int        `gorm:"column:role_type_id;primaryKey"`
	Name       string     `gorm:"column:name;not null"`
	Detail     *string    `gorm:"column:detail;type:text"`
	Created    time.Time  `gorm:"column:created;not null"`
	Updated    *time.Time `gorm:"column:updated"`
	Insert     *int       `gorm:"column:insert"`
	Update     *int       `gorm:"column:update"`
	Delete     *int       `gorm:"column:delete"`
	Select     *int       `gorm:"column:select"`
	Log        *int       `gorm:"column:log"`
	Config     *int       `gorm:"column:config"`
	Truncate   *int       `gorm:"column:truncate"`
}

func (SdUserRolesPermision) TableName() string { return "sd_user_roles_permision" }

// Tnb is an alias for dashboard config (same structure).
// อะไหล่ของแดชบอร์ด (โครงสร้างเดียวกับ sd_dashboard_config)
type Tnb SdDashboardConfig

func (Tnb) TableName() string { return "tnb" }

// ============================================================================
// Job Module (โมดูลใบรับงานซ่อม)
// ============================================================================

// Job represents a repair job card.
// ใบรับงานซ่อม
type Job struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobNo         string     `gorm:"column:job_no;type:varchar(20);unique;not null"`
	CustomerID    uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	CarID         uuid.UUID  `gorm:"column:car_id;type:uuid;not null"`
	MechanicID    uuid.UUID  `gorm:"column:mechanic_id;type:uuid;not null"`
	Status        string     `gorm:"column:status;type:varchar(30);not null;default:OPEN"`
	StartDate     time.Time  `gorm:"column:start_date;not null;default:now()"`
	EndDate       *time.Time `gorm:"column:end_date"`
	Symptom       *string    `gorm:"column:symptom;type:text"`
	DiagnosisNote *string    `gorm:"column:diagnosis_note;type:text"`
	Mileage       *int       `gorm:"column:mileage"`
	EstimatedCost *float64   `gorm:"column:estimated_cost;type:decimal(15,2)"`
	ActualCost    *float64   `gorm:"column:actual_cost;type:decimal(15,2)"`
	Priority      string     `gorm:"column:priority;type:varchar(20);default:NORMAL"`
	UserID        uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID  uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
	Deleted       bool       `gorm:"column:deleted;default:false"`
}

func (Job) TableName() string { return "t_job" }

// JobService represents a service item in a job card.
// รายการบริการในใบรับงานซ่อม
type JobService struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobID        uuid.UUID  `gorm:"column:job_id;type:uuid;not null"`
	ServiceID    uuid.UUID  `gorm:"column:service_id;type:uuid;not null"`
	Quantity     int        `gorm:"column:quantity;not null;default:1"`
	UnitPrice    float64    `gorm:"column:unit_price;type:decimal(15,2);not null"`
	Discount     float64    `gorm:"column:discount;type:decimal(15,2);default:0"`
	Note         *string    `gorm:"column:note;type:text"`
	UserID       uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (JobService) TableName() string { return "t_job_service" }

// JobPartSales represents a part used in a job card.
// อะไหล่ที่ใช้ในใบรับงานซ่อม
type JobPartSales struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobID        uuid.UUID  `gorm:"column:job_id;type:uuid;not null"`
	PartID       uuid.UUID  `gorm:"column:part_id;type:uuid;not null"`
	Quantity     int        `gorm:"column:quantity;not null;default:1"`
	UnitPrice    float64    `gorm:"column:unit_price;type:decimal(15,2);not null"`
	Discount     float64    `gorm:"column:discount;type:decimal(15,2);default:0"`
	Note         *string    `gorm:"column:note;type:text"`
	UserID       uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (JobPartSales) TableName() string { return "t_job_part_sales" }

// JobSymptom represents a car symptom reported for a job.
// อาการรถที่แจ้งในใบรับงานซ่อม
type JobSymptom struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobID              uuid.UUID `gorm:"column:job_id;type:uuid;not null"`
	SymptomCode        *string   `gorm:"column:symptom_code;type:varchar(20)"`
	SymptomDescription string    `gorm:"column:symptom_description;type:text;not null"`
	Severity           string    `gorm:"column:severity;type:varchar(20);default:MEDIUM"`
	ReportedBy         *string   `gorm:"column:reported_by;type:varchar(100)"`
	UserID             uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID       uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (JobSymptom) TableName() string { return "t_job_service_car_symptom" }

// JobDiagTroubleCode stores diagnostic trouble codes for a job.
// รหัสข้อผิดพลาด DTC สำหรับงานซ่อม
type JobDiagTroubleCode struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobID       uuid.UUID `gorm:"column:job_id;type:uuid;not null"`
	TroubleCode string    `gorm:"column:trouble_code;type:varchar(20);not null"`
	Description *string   `gorm:"column:description;type:text"`
	System      *string   `gorm:"column:system;type:varchar(50)"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (JobDiagTroubleCode) TableName() string { return "t_job_diag_trouble_code" }

// JobStatusHistory tracks job status changes.
// ประวัติการเปลี่ยนสถานะของใบรับงานซ่อม
type JobStatusHistory struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	JobID        uuid.UUID  `gorm:"column:job_id;type:uuid;not null"`
	FromStatus   *string    `gorm:"column:from_status;type:varchar(30)"`
	ToStatus     string     `gorm:"column:to_status;type:varchar(30);not null"`
	ChangedBy    uuid.UUID  `gorm:"column:changed_by;type:uuid;not null"`
	ChangedAt    time.Time  `gorm:"column:changed_at;not null;default:now()"`
	Reason       *string    `gorm:"column:reason;type:text"`
	WhitelabelID uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (JobStatusHistory) TableName() string { return "t_job_status_history" }

// ============================================================================
// Document Module (โมดูลเอกสาร)
// ============================================================================

// Document stores uploaded file metadata.
// เก็บข้อมูลไฟล์ที่อัปโหลด
type Document struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Filename     string    `gorm:"column:filename;type:varchar(255);not null"`
	OriginalName string    `gorm:"column:original_name;type:varchar(255);not null"`
	MimeType     string    `gorm:"column:mime_type;type:varchar(100);not null"`
	Size         int64     `gorm:"column:size;not null"`
	Path         string    `gorm:"column:path;type:varchar(500);not null"`
	UploadedBy   uuid.UUID `gorm:"column:uploaded_by;type:uuid"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Document) TableName() string { return "document" }

// ============================================================================
// Email Module (โมดูลอีเมล)
// ============================================================================

// EmailLog stores sent email records.
// บันทึกการส่งอีเมล
type EmailLog struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement"`
	To           string     `gorm:"column:to;type:varchar(255);not null"`
	Cc           *string    `gorm:"column:cc;type:varchar(500)"`
	Bcc          *string    `gorm:"column:bcc;type:varchar(500)"`
	Subject      string     `gorm:"column:subject;type:varchar(500);not null"`
	Body         string     `gorm:"column:body;type:text;not null"`
	Status       string     `gorm:"column:status;type:varchar(20);not null;default:pending"`
	ErrorMessage *string    `gorm:"column:error_message;type:text"`
	SentAt       *time.Time `gorm:"column:sent_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
}

func (EmailLog) TableName() string { return "email_log" }

// EmailConfig stores SMTP configuration.
// การตั้งค่า SMTP
type EmailConfig struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	SmtpHost    string    `gorm:"column:smtp_host;type:varchar(255);not null"`
	SmtpPort    int       `gorm:"column:smtp_port;not null"`
	SmtpUser    string    `gorm:"column:smtp_user;type:varchar(255);not null"`
	SmtpPass    string    `gorm:"column:smtp_pass;type:varchar(255);not null"`
	FromEmail   string    `gorm:"column:from_email;type:varchar(255);not null"`
	FromName    string    `gorm:"column:from_name;type:varchar(255);not null"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (EmailConfig) TableName() string { return "email_config" }

// ============================================================================
// Batch Module (โมดูลแบตช์)
// ============================================================================

// BatchJob represents a batch processing job.
// งานประมวลผลแบบแบตช์
type BatchJob struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string     `gorm:"column:name;type:varchar(255);not null"`
	Type         string     `gorm:"column:type;type:varchar(100);not null"`
	Status       string     `gorm:"column:status;type:varchar(20);not null;default:pending"`
	Config       string     `gorm:"column:config;type:text"`
	Schedule     *string    `gorm:"column:schedule;type:varchar(100)"`
	TotalCount   int        `gorm:"column:total_count;default:0"`
	SuccessCount int        `gorm:"column:success_count;default:0"`
	FailCount    int        `gorm:"column:fail_count;default:0"`
	StartedAt    *time.Time `gorm:"column:started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;default:now()"`
}

func (BatchJob) TableName() string { return "t_batch_job" }

// BatchJobLog stores batch job execution logs.
// บันทึกการทำงานของแบตช์
type BatchJobLog struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	JobID     uuid.UUID `gorm:"column:job_id;type:uuid;not null"`
	Message   string    `gorm:"column:message;type:text;not null"`
	Level     string    `gorm:"column:level;type:varchar(20);not null;default:info"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (BatchJobLog) TableName() string { return "t_batch_job_log" }

// ============================================================================
// I18n Module (โมดูลแปลภาษา)
// ============================================================================

// Translation stores multilingual key-value pairs.
// เก็บค่าข้อความหลายภาษา
type Translation struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Locale    string    `gorm:"column:locale;type:varchar(10);not null;index:idx_locale_key,unique"`
	Key       string    `gorm:"column:key;type:varchar(255);not null;index:idx_locale_key,unique"`
	Value     string    `gorm:"column:value;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Translation) TableName() string { return "m_translation" }

// ============================================================================
// WOS Module — Web Order System (โมดูลสั่งซื้อออนไลน์)
// ============================================================================

// WosOrder represents a customer order from the web order system.
// คำสั่งซื้อจากระบบสั่งซื้อออนไลน์
type WosOrder struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderNumber   string          `gorm:"column:order_number;type:varchar(50);not null;unique"`
	CustomerName  string          `gorm:"column:customer_name;type:varchar(255);not null"`
	CustomerEmail string          `gorm:"column:customer_email;type:varchar(255);not null"`
	CustomerPhone *string         `gorm:"column:customer_phone;type:varchar(50)"`
	Items         json.RawMessage `gorm:"column:items;type:jsonb;not null"`
	TotalAmount   float64         `gorm:"column:total_amount;type:decimal(12,2);not null"`
	Status        string          `gorm:"column:status;type:varchar(20);not null;default:pending"`
	Notes         *string         `gorm:"column:notes;type:text"`
	CreatedAt     time.Time       `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;not null;default:now()"`
}

func (WosOrder) TableName() string { return "wos_order" }

// ============================================================================
// Customer Module (โมดูลลูกค้า)
// ============================================================================

// Customer represents a customer (individual or corporate).
type Customer struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerCode    string     `gorm:"column:customer_code;type:varchar(20);unique;not null"`
	FullName        string     `gorm:"column:full_name;type:varchar(200);not null"`
	DisplayName     *string    `gorm:"column:display_name;type:varchar(200)"`
	CustomerType    string     `gorm:"column:customer_type;type:varchar(20);not null;default:INDIVIDUAL"`
	Status          string     `gorm:"column:status;type:varchar(20);not null;default:ACTIVE"`
	TaxID           *string    `gorm:"column:tax_id;type:varchar(20)"`
	Email           *string    `gorm:"column:email;type:varchar(100)"`
	PhoneNumber     string     `gorm:"column:phone_number;type:varchar(20);not null"`
	SecondaryPhone  *string    `gorm:"column:secondary_phone;type:varchar(20)"`
	Address         *string    `gorm:"column:address;type:text"`
	Province        *string    `gorm:"column:province;type:varchar(100)"`
	City            *string    `gorm:"column:city;type:varchar(100)"`
	District        *string    `gorm:"column:district;type:varchar(100)"`
	PostalCode      *string    `gorm:"column:postal_code;type:varchar(10)"`
	Country         *string    `gorm:"column:country;type:varchar(50);default:Thailand"`
	ContactPerson   *string    `gorm:"column:contact_person;type:varchar(100)"`
	ContactPhone    *string    `gorm:"column:contact_phone;type:varchar(20)"`
	Notes           *string    `gorm:"column:notes;type:text"`
	LastVisitDate   *time.Time `gorm:"column:last_visit_date"`
	TotalVisitCount int        `gorm:"column:total_visit_count;default:0"`
	TotalSpent      float64    `gorm:"column:total_spent;type:decimal(15,2);default:0"`
	UserID          uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID    uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
	Deleted         bool       `gorm:"column:deleted;default:false"`
}

func (Customer) TableName() string { return "m_customer" }

// Car represents a customer's vehicle.
type Car struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID        uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	LicensePlate      string     `gorm:"column:license_plate;type:varchar(20);unique;not null"`
	Province          *string    `gorm:"column:province;type:varchar(50)"`
	Brand             string     `gorm:"column:brand;type:varchar(50);not null"`
	Model             string     `gorm:"column:model;type:varchar(100);not null"`
	SubModel          *string    `gorm:"column:sub_model;type:varchar(100)"`
	Year              *int       `gorm:"column:year"`
	Color             *string    `gorm:"column:color;type:varchar(30)"`
	EngineNumber      *string    `gorm:"column:engine_number;type:varchar(50)"`
	ChassisNumber     *string    `gorm:"column:chassis_number;type:varchar(50)"`
	FuelType          *string    `gorm:"column:fuel_type;type:varchar(20)"`
	TransmissionType  *string    `gorm:"column:transmission_type;type:varchar(20)"`
	EngineCC          *int       `gorm:"column:engine_cc"`
	SeatingCapacity   *int       `gorm:"column:seating_capacity"`
	Mileage           int        `gorm:"column:mileage;default:0"`
	LastServiceDate   *time.Time `gorm:"column:last_service_date"`
	NextServiceMileage *int      `gorm:"column:next_service_mileage"`
	Notes             *string    `gorm:"column:notes;type:text"`
	IsActive          bool       `gorm:"column:is_active;default:true"`
	UserID            uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID      uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt         *time.Time `gorm:"column:updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at"`
	Deleted           bool       `gorm:"column:deleted;default:false"`
}

func (Car) TableName() string { return "m_car" }

// CarServiceHistory records vehicle service history.
type CarServiceHistory struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CarID           uuid.UUID `gorm:"column:car_id;type:uuid;not null"`
	JobID           uuid.UUID `gorm:"column:job_id;type:uuid;not null"`
	ServiceDate     time.Time `gorm:"column:service_date;not null"`
	ServiceType     *string   `gorm:"column:service_type;type:varchar(50)"`
	Description     *string   `gorm:"column:description;type:text"`
	TotalCost       *float64  `gorm:"column:total_cost;type:decimal(15,2)"`
	MileageAtService *int     `gorm:"column:mileage_at_service"`
	MechanicName    *string   `gorm:"column:mechanic_name;type:varchar(100)"`
	WhitelabelID    uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (CarServiceHistory) TableName() string { return "m_car_service_history" }

// ============================================================================
// Quotation Module (โมดูลใบเสนอราคา)
// ============================================================================

// Quotation represents a quotation document.
type Quotation struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuotationNo      string     `gorm:"column:quotation_no;type:varchar(20);unique;not null"`
	JobID            uuid.UUID  `gorm:"column:job_id;type:uuid;not null"`
	CustomerID       uuid.UUID  `gorm:"column:customer_id;type:uuid;not null"`
	QuotationDate    time.Time  `gorm:"column:quotation_date;not null;default:now()"`
	ExpiryDate       time.Time  `gorm:"column:expiry_date;not null"`
	Status           string     `gorm:"column:status;type:varchar(20);not null;default:DRAFT"`
	Subtotal         float64    `gorm:"column:subtotal;type:decimal(15,2);not null;default:0"`
	TaxRate          float64    `gorm:"column:tax_rate;type:decimal(5,2);default:7.00"`
	TaxAmount        float64    `gorm:"column:tax_amount;type:decimal(15,2);default:0"`
	DiscountType     *string    `gorm:"column:discount_type;type:varchar(20)"`
	DiscountValue    float64    `gorm:"column:discount_value;type:decimal(15,2);default:0"`
	Total            float64    `gorm:"column:total;type:decimal(15,2);not null;default:0"`
	AmountInWordsTh  *string    `gorm:"column:amount_in_words_th;type:text"`
	AmountInWordsEn  *string    `gorm:"column:amount_in_words_en;type:text"`
	Currency         *string    `gorm:"column:currency;type:varchar(10);default:THB"`
	ExchangeRate     float64    `gorm:"column:exchange_rate;type:decimal(10,4);default:1.0000"`
	Notes            *string    `gorm:"column:notes;type:text"`
	TermsAndConditions *string  `gorm:"column:terms_and_conditions;type:text"`
	ApprovedBy       *uuid.UUID `gorm:"column:approved_by;type:uuid"`
	ApprovedAt       *time.Time `gorm:"column:approved_at"`
	RejectedReason   *string    `gorm:"column:rejected_reason;type:text"`
	ConvertedToPo    bool       `gorm:"column:converted_to_po;default:false"`
	UserID           uuid.UUID  `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID     uuid.UUID  `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt        *time.Time `gorm:"column:updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at"`
	Deleted          bool       `gorm:"column:deleted;default:false"`
}

func (Quotation) TableName() string { return "t_quotation" }

// QuotationPart represents a part item in a quotation.
type QuotationPart struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuotationID uuid.UUID `gorm:"column:quotation_id;type:uuid;not null"`
	PartID      uuid.UUID `gorm:"column:part_id;type:uuid;not null"`
	Quantity    int       `gorm:"column:quantity;not null;default:1"`
	UnitPrice   float64   `gorm:"column:unit_price;type:decimal(15,2);not null"`
	Discount    float64   `gorm:"column:discount;type:decimal(15,2);default:0"`
	Note        *string   `gorm:"column:note;type:text"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (QuotationPart) TableName() string { return "t_quotation_part" }

// QuotationService represents a service item in a quotation.
type QuotationService struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuotationID uuid.UUID `gorm:"column:quotation_id;type:uuid;not null"`
	ServiceID   uuid.UUID `gorm:"column:service_id;type:uuid;not null"`
	Quantity    int       `gorm:"column:quantity;not null;default:1"`
	UnitPrice   float64   `gorm:"column:unit_price;type:decimal(15,2);not null"`
	Discount    float64   `gorm:"column:discount;type:decimal(15,2);default:0"`
	Note        *string   `gorm:"column:note;type:text"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	WhitelabelID uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (QuotationService) TableName() string { return "t_quotation_service" }

// QuotationStatusHistory tracks quotation status changes.
type QuotationStatusHistory struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuotationID  uuid.UUID  `gorm:"column:quotation_id;type:uuid;not null"`
	FromStatus   *string    `gorm:"column:from_status;type:varchar(20)"`
	ToStatus     string     `gorm:"column:to_status;type:varchar(20);not null"`
	ChangedBy    uuid.UUID  `gorm:"column:changed_by;type:uuid;not null"`
	ChangedAt    time.Time  `gorm:"column:changed_at;not null;default:now()"`
	Reason       *string    `gorm:"column:reason;type:text"`
	WhitelabelID uuid.UUID `gorm:"column:whitelabel_id;type:uuid;not null"`
}

func (QuotationStatusHistory) TableName() string { return "t_quotation_status_history" }

/*
	## คำอธิบายเพิ่มเติม
	- **ความสัมพันธ์ (Relationships)**: เราได้เพิ่ม foreign key references ใน struct ต่าง ๆ เช่น `gorm:"foreignKey:..."` เพื่อให้ GORM จัดการความสัมพันธ์อัตโนมัติ
	- **JSONB fields**: ใช้ `json.RawMessage` เพื่อรองรับการ Marshal/Unmarshal แบบยืดหยุ่น
	- **UUID**: สำหรับ primary key ที่เป็น uuid เราใช้ `github.com/google/uuid` และ default `gen_random_uuid()`
	- **Nullable fields**: ใช้ pointer (`*string`, `*int`, `*time.Time`) สำหรับคอลัมน์ที่อนุญาต NULL
	- **TableName**: ทุกโมเดลมีเมธอด `TableName()` ชี้ไปยังชื่อตารางจริงในฐานข้อมูล

	โมเดลทั้งหมดนี้ครอบคลุมทุกตารางในไฟล์ SQL ที่ให้มา พร้อมใช้งานกับ GORM และ PostgreSQL ได้ทันที

	Important: The SQL has many tables with prefix "sd_" and others. We'll create models for:
	Table
		activity_log
		command_log
		device_alert
		device_config
		device_status
		iot_data
		migrations (skip? but include)
		noti_notification_logs
		noti_notification_rules
		noti_notification_types
		noti_notifications
		notification_devices
		notification_groups
		notification_groups_devices_notification_devices (join table, may be omitted if using many2many)
		notification_logs
		notification_types
		sd_activity_log
		sd_activity_type_log
		sd_admin_access_menu
		sd_air_control
		sd_air_control_device_map
		sd_air_control_log
		sd_air_mod
		sd_air_mod_device_map
		sd_air_period
		sd_air_period_device_map
		sd_air_setting_warning
		sd_air_setting_warning_device_map
		sd_air_warning
		sd_air_warning_device_map
		sd_alarm_process_log
		sd_alarm_process_log_email
		sd_alarm_process_log_line
		sd_alarm_process_log_mqtt
		sd_alarm_process_log_sms
		sd_alarm_process_log_telegram
		sd_alarm_process_log_temp
		sd_api_key
		sd_audit_log
		sd_channel_template
		sd_dashboard_config
		sd_device_category
		sd_device_group
		sd_device_log
		sd_device_member
		sd_device_notification_config
		sd_device_schedule
		sd_device_status_history
		sd_group_notification_config
		sd_iot_alarm_device
		sd_iot_alarm_device_event
		sd_iot_api
		sd_iot_device
		sd_iot_device_action
		sd_iot_device_action_log
		sd_iot_device_action_user
		sd_iot_device_alarm_action
		sd_iot_device_type
		sd_iot_email
		sd_iot_group
		sd_iot_host
		sd_iot_influxdb
		sd_iot_line
		sd_iot_location
		sd_iot_mqtt
		sd_iot_nodered
		sd_iot_schedule
		sd_iot_schedule_device
		sd_iot_sensor
		sd_iot_setting
		sd_iot_sms
		sd_iot_telegram
		sd_iot_token
		sd_iot_type
		sd_module_log
		sd_mqtt_host
		sd_mqtt_log
		sd_notification_channel
		sd_notification_condition
		sd_notification_log
		sd_notification_type
		sd_report_data
		sd_schedule_process_log
		sd_sensor_data
		sd_system_setting
		sd_user
		sd_user_access_menu
		sd_user_file
		sd_user_log
		sd_user_log_type
		sd_user_role
		sd_user_roles_access
		sd_user_roles_permision
		tnb

*/
