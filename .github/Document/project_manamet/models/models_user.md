เราจะสร้าง GORM models สำหรับตารางหลักจาก `icmon.sql` ที่เกี่ยวข้องกับระบบ IoT Gateway, การแจ้งเตือน, ผู้ใช้ และอุปกรณ์ โดยเลือกเฉพาะ core tables ที่จำเป็นสำหรับการทำงานร่วมกับ `pkg/iot_gateway` และ `pkg/email_sender` รวมถึงระบบแจ้งเตือน

โครงสร้างไฟล์ที่จะสร้างใน `internal/models/`:

- `user.go` – ผู้ใช้ (`sd_user`) และบทบาท (`sd_user_role`)
- `device.go` – อุปกรณ์ IoT (`sd_iot_device`), ตำแหน่ง (`sd_iot_location`), กลุ่มอุปกรณ์ (`sd_device_group`, `sd_device_member`)
- `notification.go` – การแจ้งเตือน (`sd_notification_type`, `sd_notification_channel`, `sd_notification_log`, `noti_notifications`, `noti_notification_logs`)
- `alarm.go` – กฎการแจ้งเตือน (`sd_iot_device_alarm_action`)
- `sensor.go` – ข้อมูลเซ็นเซอร์ (`sd_sensor_data`) สำหรับ fallback (แม้หลักใช้ InfluxDB)
- `config.go` – การตั้งค่า MQTT, InfluxDB (`sd_iot_mqtt`, `sd_iot_influxdb`)

---

## 1. Model `User` และ `UserRole` (`internal/models/user.go`)

```go
package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// SdUser corresponds to table "sd_user"
type SdUser struct {
    ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    CreatedDate         time.Time      `gorm:"not null;default:now()" json:"created_date"`
    UpdatedDate         time.Time      `gorm:"not null;default:now()" json:"updated_date"`
    DeletedDate         *time.Time     `json:"deleted_date,omitempty"`
    RoleID              int            `gorm:"not null" json:"role_id"`
    Email               string         `gorm:"size:255;not null" json:"email"`
    Username            string         `gorm:"size:255;not null" json:"username"`
    Password            string         `gorm:"size:255;not null" json:"-"`
    PasswordTemp        *string        `gorm:"size:255" json:"password_temp,omitempty"`
    FirstName           *string        `gorm:"size:255" json:"first_name"`
    LastName            *string        `gorm:"size:255" json:"last_name"`
    FullName            *string        `gorm:"size:255" json:"full_name"`
    Nickname            *string        `gorm:"size:255" json:"nickname"`
    IDCard              *string        `gorm:"size:255" json:"id_card"`
    LastSignInDate      time.Time      `gorm:"not null;default:now()" json:"last_sign_in_date"`
    Status              int16          `gorm:"not null" json:"status"`
    ActiveStatus        *int16         `json:"active_status"`
    NetworkID           *int           `gorm:"default:1" json:"network_id"`
    Remark              *string        `gorm:"size:255" json:"remark"`
    InfoAgreeStatus     *int16         `gorm:"default:0" json:"info_agree_status"`
    Gender              *string        `gorm:"size:255" json:"gender"`
    Birthday            *time.Time     `json:"birthday"`
    OnlineStatus        *string        `gorm:"size:255;default:'0'" json:"online_status"`
    Message             *string        `gorm:"size:255" json:"message"`
    NetworkTypeID       *int           `gorm:"default:0" json:"network_type_id"`
    PublicStatus        *int16         `gorm:"default:0" json:"public_status"`
    TypeID              *int           `gorm:"default:0" json:"type_id"`
    AvatarPath          *string        `gorm:"size:255" json:"avatar_path"`
    Avatar              *string        `gorm:"size:255" json:"avatar"`
    RefreshToken        *string        `gorm:"type:text" json:"refresh_token"`
    LoginFailed         *int16         `json:"login_failed"`
    PublicNotification  *int16         `gorm:"default:0" json:"public_notification"`
    SMSNotification     *int16         `gorm:"default:0" json:"sms_notification"`
    EmailNotification   *int16         `gorm:"default:0" json:"email_notification"`
    LineNotification    *int16         `gorm:"default:0" json:"line_notification"`
    MobileNumber        *string        `gorm:"size:255;default:'0'" json:"mobile_number"`
    PhoneNumber         *string        `gorm:"size:255;default:'0'" json:"phone_number"`
    LineID              *string        `gorm:"size:255;default:'0'" json:"line_id"`
    SystemID            *string        `gorm:"size:255;default:'1'" json:"system_id"`
    LocationID          *string        `gorm:"size:255;default:'1'" json:"location_id"`

    // Relationships
    Role *SdUserRole `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

func (SdUser) TableName() string {
    return "sd_user"
}

// SdUserRole corresponds to table "sd_user_role"
type SdUserRole struct {
    ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
    RoleID       int       `gorm:"not null" json:"role_id"`
    Title        *string   `gorm:"size:50" json:"title"`
    CreatedDate  *time.Time `gorm:"default:now()" json:"created_date"`
    UpdatedDate  *time.Time `gorm:"default:now()" json:"updated_date"`
    CreateBy     int       `gorm:"not null" json:"create_by"`
    LastUpdateBy int       `gorm:"not null" json:"last_update_by"`
    Status       int16     `gorm:"not null" json:"status"`
    TypeID       int       `gorm:"not null" json:"type_id"`
    Lang         string    `gorm:"size:255;not null" json:"lang"`
}

func (SdUserRole) TableName() string {
    return "sd_user_role"
}
```

---

## 2. Model อุปกรณ์และตำแหน่ง (`internal/models/device.go`)

```go
package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
    "gorm.io/gorm"
)

// SdIoTDevice corresponds to table "sd_iot_device"
type SdIoTDevice struct {
    DeviceID            int            `gorm:"primaryKey;autoIncrement" json:"device_id"`
    SettingID           *int           `json:"setting_id"`
    TypeID              *int           `json:"type_id"`
    LocationID          *int           `json:"location_id"`
    DeviceName          *string        `gorm:"size:255" json:"device_name"`
    SN                  *string        `gorm:"size:255;uniqueIndex" json:"sn"`
    HardwareID          *int           `json:"hardware_id"`
    StatusWarning       *string        `gorm:"size:150" json:"status_warning"`
    RecoveryWarning     *string        `gorm:"size:150" json:"recovery_warning"`
    StatusAlert         *string        `gorm:"size:150" json:"status_alert"`
    RecoveryAlert       *string        `gorm:"size:150" json:"recovery_alert"`
    TimeLife            *int           `gorm:"default:1" json:"time_life"`
    Period              *string        `gorm:"size:150" json:"period"`
    WorkStatus          *int           `gorm:"default:1" json:"work_status"`
    Model               *string        `gorm:"size:255" json:"model"`
    Vendor              *string        `gorm:"size:255" json:"vendor"`
    CompareValue        *string        `gorm:"size:255" json:"compare_value"`
    Unit                *string        `gorm:"size:255" json:"unit"`
    MqttID              *int           `json:"mqtt_id"`
    OID                 *string        `gorm:"size:255" json:"oid"`
    ActionID            *int           `json:"action_id"`
    StatusAlertID       *int           `json:"status_alert_id"`
    MqttDataValue       *string        `gorm:"size:255" json:"mqtt_data_value"`
    MqttDataControl     *string        `gorm:"size:255" json:"mqtt_data_control"`
    Measurement         *string        `gorm:"size:255" json:"measurement"`
    MqttControlOn       *string        `gorm:"size:255;default:'1'" json:"mqtt_control_on"`
    MqttControlOff      *string        `gorm:"size:255;default:'0'" json:"mqtt_control_off"`
    Org                 string         `gorm:"size:255;not null" json:"org"`
    Bucket              string         `gorm:"size:255;not null" json:"bucket"`
    Status              *int           `json:"status"`
    MqttDeviceName      string         `gorm:"size:255;not null" json:"mqtt_device_name"`
    MqttStatusOverName  *string        `gorm:"type:text" json:"mqtt_status_over_name"`
    MqttStatusDataName  *string        `gorm:"type:text" json:"mqtt_status_data_name"`
    MqttActRelayName    *string        `gorm:"type:text" json:"mqtt_act_relay_name"`
    MqttControlRelayName *string       `gorm:"type:text" json:"mqtt_control_relay_name"`
    MqttConfig          *string        `gorm:"type:text" json:"mqtt_config"`
    CreatedDate         time.Time      `gorm:"not null;default:now()" json:"created_date"`
    UpdatedDate         time.Time      `gorm:"not null;default:now()" json:"updated_date"`
    Max                 *string        `gorm:"size:255" json:"max"`
    Min                 *string        `gorm:"size:255" json:"min"`
    Layout              *int           `gorm:"default:1" json:"layout"`
    AlertSet            *int           `gorm:"default:1" json:"alert_set"`
    IconNormal          string         `gorm:"type:text;not null" json:"icon_normal"`
    IconWarning         string         `gorm:"type:text;not null" json:"icon_warning"`
    IconAlert           string         `gorm:"type:text;not null" json:"icon_alert"`
    Icon                string         `gorm:"type:text;not null" json:"icon"`
    ColorNormal         string         `gorm:"not null;default:'#22C55E'" json:"color_normal"`
    ColorWarning        string         `gorm:"not null;default:'#F59E0B'" json:"color_warning"`
    ColorAlarm          string         `gorm:"not null;default:'#EF4444'" json:"color_alarm"`
    Code                string         `gorm:"not null;default:'normal'" json:"code"`
    Menu                *int           `gorm:"default:1" json:"menu"`
    IconOn              string         `gorm:"type:text;not null" json:"icon_on"`
    IconOff             string         `gorm:"type:text;not null" json:"icon_off"`
    CalibrationAdd      *string        `gorm:"size:250;default:'0'" json:"calibration_add"`
    CalibrationSubtract *string        `gorm:"size:250;default:'0'" json:"calibration_subtract"`
    CalibrationType     *int           `gorm:"default:3" json:"calibration_type"`

    // Relationships
    Location   *SdIoTLocation            `gorm:"foreignKey:LocationID" json:"location,omitempty"`
    AlarmAction *SdIoTDeviceAlarmAction  `gorm:"foreignKey:ActionID" json:"alarm_action,omitempty"`
    MqttConfig  *SdIoTMqtt               `gorm:"foreignKey:MqttID" json:"mqtt_config,omitempty"`
}

func (SdIoTDevice) TableName() string {
    return "sd_iot_device"
}

// SdIoTLocation corresponds to table "sd_iot_location"
type SdIoTLocation struct {
    LocationID     int        `gorm:"primaryKey;autoIncrement" json:"location_id"`
    LocationName   string     `gorm:"size:255;not null" json:"location_name"`
    IPAddress      string     `gorm:"size:255;not null" json:"ip_address"`
    LocationDetail string     `gorm:"not null" json:"location_detail"`
    CreatedDate    time.Time  `gorm:"not null;default:now()" json:"created_date"`
    UpdatedDate    time.Time  `gorm:"not null;default:now()" json:"updated_date"`
    Status         *int       `json:"status"`
    ConfigData     *string    `gorm:"type:text" json:"config_data"`
}

func (SdIoTLocation) TableName() string {
    return "sd_iot_location"
}

// SdDeviceGroup corresponds to table "sd_device_group"
type SdDeviceGroup struct {
    ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string         `gorm:"size:200;not null" json:"name"`
    Description *string        `gorm:"type:text" json:"description"`
    GroupType   string         `gorm:"size:50;not null;default:'custom'" json:"group_type"`
    IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
    Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
    CreatedAt   time.Time      `gorm:"not null;default:now()" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"not null;default:now()" json:"updated_at"`

    Members []SdDeviceMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

func (SdDeviceGroup) TableName() string {
    return "sd_device_group"
}

// SdDeviceMember corresponds to table "sd_device_member"
type SdDeviceMember struct {
    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
    DeviceID  int       `gorm:"column:Device_id;not null" json:"device_id"`
    GroupID   int       `gorm:"not null" json:"group_id"`
    Role      string    `gorm:"size:50;not null;default:'member'" json:"role"`
    Priority  int       `gorm:"not null;default:1" json:"priority"`
    IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
    CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`

    Device SdIoTDevice  `gorm:"foreignKey:DeviceID;references:DeviceID" json:"device,omitempty"`
    Group  SdDeviceGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

func (SdDeviceMember) TableName() string {
    return "sd_device_member"
}
```

---

## 3. Model การแจ้งเตือน (`internal/models/notification.go`)

```go
package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// SdNotificationType corresponds to table "sd_notification_type"
type SdNotificationType struct {
    ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Name            string    `gorm:"size:50;not null" json:"name"`
    Description     *string   `gorm:"type:text" json:"description"`
    CooldownMinutes int       `gorm:"not null;default:10" json:"cooldown_minutes"`
    IsActive        bool      `gorm:"not null;default:true" json:"is_active"`
    Icon            *string   `gorm:"size:100" json:"icon"`
    Color           *string   `gorm:"size:20" json:"color"`
    CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
    UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (SdNotificationType) TableName() string {
    return "sd_notification_type"
}

// SdNotificationChannel corresponds to table "sd_notification_channel"
type SdNotificationChannel struct {
    ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Description *string   `gorm:"type:text" json:"description"`
    Icon        *string   `gorm:"size:100" json:"icon"`
    HandlerClass *string  `gorm:"size:200" json:"handler_class"`
    IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
    CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

func (SdNotificationChannel) TableName() string {
    return "sd_notification_channel"
}

// SdNotificationLog corresponds to table "sd_notification_log"
type SdNotificationLog struct {
    ID                   int            `gorm:"primaryKey;autoIncrement" json:"id"`
    DeviceID             *int           `json:"device_id"`
    NotificationTypeID   *int           `json:"notification_type_id"`
    NotificationChannelID *int          `json:"notification_channel_id"`
    Message              string         `gorm:"type:text;not null" json:"message"`
    ResponseData         datatypes.JSON `gorm:"type:jsonb" json:"response_data"`
    SentAt               *time.Time     `json:"sent_at"`
    CreatedAt            time.Time      `gorm:"not null;default:now()" json:"created_at"`
    TemplateID           *int           `json:"template_id"`
    DeliveredAt          *time.Time     `json:"delivered_at"`
    ReadAt               *time.Time     `json:"read_at"`
    RetryCount           int            `gorm:"not null;default:0" json:"retry_count"`
    ErrorMessage         *string        `gorm:"type:text" json:"error_message"`
    MessageID            *string        `gorm:"size:100" json:"message_id"`
    Recipient            *string        `gorm:"size:255" json:"recipient"`
    Status               string         `gorm:"size:20;not null;default:'pending'" json:"status"`

    Device   *SdIoTDevice           `gorm:"foreignKey:DeviceID;references:DeviceID" json:"device,omitempty"`
    Type     *SdNotificationType    `gorm:"foreignKey:NotificationTypeID" json:"type,omitempty"`
    Channel  *SdNotificationChannel `gorm:"foreignKey:NotificationChannelID" json:"channel,omitempty"`
}

func (SdNotificationLog) TableName() string {
    return "sd_notification_log"
}

// NotiNotification corresponds to table "noti_notifications" (newer style)
type NotiNotification struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    Title        string         `gorm:"size:255;not null" json:"title"`
    Message      string         `gorm:"type:text;not null" json:"message"`
    Type         string         `gorm:"size:50;not null" json:"type"`
    Priority     string         `gorm:"size:50;not null" json:"priority"`
    Category     *string        `gorm:"size:255" json:"category"`
    UserID       *int           `json:"user_id"`
    UserUUID     *uuid.UUID     `gorm:"type:uuid" json:"user_uuid"`
    Metadata     datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
    IsRead       bool           `gorm:"not null;default:false" json:"is_read"`
    ReadAt       *time.Time     `json:"read_at"`
    IsSent       bool           `gorm:"not null;default:false" json:"is_sent"`
    ChannelsSent datatypes.JSON `gorm:"type:jsonb" json:"channels_sent"`
    ScheduledAt  *time.Time     `json:"scheduled_at"`
    ExpiresAt    *time.Time     `json:"expires_at"`
    Status       int            `gorm:"not null;default:1" json:"status"`
    CreatedAt    time.Time      `gorm:"not null;default:now()" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"not null;default:now()" json:"updated_at"`
    DeletedAt    *time.Time     `json:"deleted_at"`
}

func (NotiNotification) TableName() string {
    return "noti_notifications"
}

// NotiNotificationLog corresponds to table "noti_notification_logs"
type NotiNotificationLog struct {
    LogID          int            `gorm:"primaryKey;autoIncrement" json:"log_id"`
    NotificationID uuid.UUID      `gorm:"type:uuid;not null" json:"notification_id"`
    Channel        string         `gorm:"size:50;not null" json:"channel"`
    Payload        datatypes.JSON `gorm:"type:jsonb;not null" json:"payload"`
    Response       datatypes.JSON `gorm:"type:jsonb" json:"response"`
    Status         string         `gorm:"size:50;not null;default:'pending'" json:"status"`
    RetryCount     int            `gorm:"default:0" json:"retry_count"`
    ErrorMessage   *string        `gorm:"type:text" json:"error_message"`
    CreatedAt      *time.Time     `gorm:"default:now()" json:"created_at"`
    SentAt         *time.Time     `json:"sent_at"`
    DeliveredAt    *time.Time     `json:"delivered_at"`

    Notification NotiNotification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
}

func (NotiNotificationLog) TableName() string {
    return "noti_notification_logs"
}
```

---

## 4. Model กฎการแจ้งเตือน (`internal/models/alarm.go`)

```go
package models

// SdIoTDeviceAlarmAction corresponds to table "sd_iot_device_alarm_action"
type SdIoTDeviceAlarmAction struct {
    AlarmActionID    int     `gorm:"primaryKey;autoIncrement" json:"alarm_action_id"`
    ActionName       *string `gorm:"size:255" json:"action_name"`
    StatusWarning    *string `gorm:"size:150" json:"status_warning"`
    RecoveryWarning  *string `gorm:"size:150" json:"recovery_warning"`
    StatusAlert      *string `gorm:"size:150" json:"status_alert"`
    RecoveryAlert    *string `gorm:"size:150" json:"recovery_alert"`
    EmailAlarm       *int    `json:"email_alarm"`
    LineAlarm        *int    `json:"line_alarm"`
    TelegramAlarm    *int    `json:"telegram_alarm"`
    SMSAlarm         *int    `json:"sms_alarm"`
    NoncAlarm        *int    `json:"nonc_alarm"`
    TimeLife         *int    `json:"time_life"`
    Event            *int    `json:"event"`
    Status           *int    `json:"status"`
}

func (SdIoTDeviceAlarmAction) TableName() string {
    return "sd_iot_device_alarm_action"
}
```

---

## 5. Model ข้อมูลเซ็นเซอร์ (fallback) (`internal/models/sensor.go`)

```go
package models

import (
    "time"

    "gorm.io/datatypes"
)

// SdSensorData corresponds to table "sd_sensor_data" (optional, for PostgreSQL fallback)
type SdSensorData struct {
    ID                 int            `gorm:"primaryKey;autoIncrement" json:"id"`
    DeviceID           int            `gorm:"not null" json:"device_id"`
    Value              float64        `gorm:"type:numeric(10,2);not null" json:"value"`
    RawData            datatypes.JSON `gorm:"type:jsonb" json:"raw_data"`
    NotificationTypeID *int           `json:"notification_type_id"`
    Timestamp          time.Time      `gorm:"not null;default:now()" json:"timestamp"`
    CreatedAt          time.Time      `gorm:"not null;default:now()" json:"created_at"`
    BatteryLevel       *float64       `gorm:"type:numeric(5,2)" json:"battery_level"`
    SignalStrength     *int           `json:"signal_strength"`

    Device SdIoTDevice `gorm:"foreignKey:DeviceID;references:DeviceID" json:"device,omitempty"`
}

func (SdSensorData) TableName() string {
    return "sd_sensor_data"
}
```

---

## 6. Model การตั้งค่า MQTT และ InfluxDB (`internal/models/config.go`)

```go
package models

import (
    "time"

    "gorm.io/datatypes"
)

// SdIoTMqtt corresponds to table "sd_iot_mqtt"
type SdIoTMqtt struct {
    MqttID        int            `gorm:"primaryKey;autoIncrement" json:"mqtt_id"`
    MqttTypeID    *int           `json:"mqtt_type_id"`
    Sort          int            `gorm:"not null;default:1" json:"sort"`
    MqttName      *string        `gorm:"type:varchar" json:"mqtt_name"`
    Host          *string        `gorm:"type:varchar" json:"host"`
    Port          *int           `json:"port"`
    Username      *string        `gorm:"type:varchar" json:"username"`
    Password      *string        `gorm:"type:varchar" json:"password"`
    Secret        *string        `gorm:"type:varchar" json:"secret"`
    ExpireIn      *string        `gorm:"type:varchar" json:"expire_in"`
    TokenValue    *string        `gorm:"type:varchar" json:"token_value"`
    Org           *string        `gorm:"type:varchar" json:"org"`
    Bucket        *string        `gorm:"type:varchar" json:"bucket"`
    Environment   *string        `gorm:"type:varchar" json:"environment"`
    CreatedDate   time.Time      `gorm:"not null;default:now()" json:"created_date"`
    UpdatedDate   time.Time      `gorm:"not null;default:now()" json:"updated_date"`
    Status        int            `gorm:"not null;default:1" json:"status"`
    LocationID    *int           `gorm:"default:1" json:"location_id"`
    Latitude      *string        `gorm:"size:255" json:"latitude"`
    Longitude     *string        `gorm:"size:255" json:"longitude"`
    MqttMainID    int            `gorm:"not null;default:1" json:"mqtt_main_id"`
    Configuration datatypes.JSON `gorm:"type:jsonb;default:'{\"0\":\"temperature1\",\"1\":\"humidity1\"}'" json:"configuration"`
    Zoom          int            `gorm:"not null;default:6" json:"zoom"`
}

func (SdIoTMqtt) TableName() string {
    return "sd_iot_mqtt"
}

// SdIoTInfluxdb corresponds to table "sd_iot_influxdb"
type SdIoTInfluxdb struct {
    InfluxdbID   uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"influxdb_id"`
    InfluxdbName *string    `gorm:"type:text" json:"influxdb_name"`
    Host         *string    `gorm:"type:text" json:"host"`
    Port         string     `gorm:"size:255;not null" json:"port"`
    Username     string     `gorm:"size:255;not null" json:"username"`
    Password     string     `gorm:"size:255;not null" json:"password"`
    TokenValue   *string    `gorm:"type:text" json:"token_value"`
    Buckets      *string    `gorm:"type:text" json:"buckets"`
    CreatedDate  time.Time  `gorm:"not null;default:now()" json:"created_date"`
    UpdatedDate  time.Time  `gorm:"not null;default:now()" json:"updated_date"`
    Status       *int       `json:"status"`
}

func (SdIoTInfluxdb) TableName() string {
    return "sd_iot_influxdb"
}
```

---

## 7. การติดตั้ง Dependencies สำหรับ datatypes.JSON

ใน `go.mod` ต้องมี:

```bash
go get gorm.io/datatypes
go get github.com/google/uuid
```

---

## สรุป

เราได้สร้าง GORM models สำหรับตารางหลักในระบบ `icmon` ที่จำเป็นสำหรับ:
- การจัดการผู้ใช้และบทบาท
- การจัดการอุปกรณ์ IoT, ตำแหน่ง, กลุ่มอุปกรณ์
- การตั้งค่า MQTT และ InfluxDB
- กฎการแจ้งเตือน (alarm action)
- การแจ้งเตือนผ่านช่องทางต่าง ๆ และ logs
- ข้อมูลเซ็นเซอร์สำรอง (fallback)

models เหล่านี้สามารถนำไปใช้ใน `internal/repository`, `internal/usecase`, และ `pkg/iot_gateway` ได้ทันที โดยเฉพาะการดึง metadata ของอุปกรณ์จาก PostgreSQL มาใช้ใน `ProcessAlarm` หรือการส่งการแจ้งเตือนผ่าน `email_sender`