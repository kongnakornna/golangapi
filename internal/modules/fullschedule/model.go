package fullschedule

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ============================================================================
// Enums
// ============================================================================

// ScheduleMode is the recurrence mode of a schedule.
type ScheduleMode string

const (
	ModeWeekly ScheduleMode = "weekly" // วันในสัปดาห์ (Sunday–Saturday)
	ModeFull   ScheduleMode = "full"   // ปฏิทิน (เดือน + วันที่)
	ModeBatch  ScheduleMode = "batch"  // Cron Expression
)

// ScheduleStatus is the active/inactive state (string) used by master-data tables
// (fs_groups / fs_zones / fs_areas) which still store status as varchar.
type ScheduleStatus string

const (
	StatusActive   ScheduleStatus = "active"
	StatusInactive ScheduleStatus = "inactive"
)

// Scan implements sql.Scanner: รับค่า int8 (หรือ string) จาก DB แล้วแปลงเป็นสถานะ string.
func (s *ScheduleStatus) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*s = StatusInactive
	case int64:
		if v == 1 {
			*s = StatusActive
		} else {
			*s = StatusInactive
		}
	case []byte:
		switch string(v) {
		case "1", "active":
			*s = StatusActive
		default:
			*s = StatusInactive
		}
	case string:
		switch v {
		case "1", "active":
			*s = StatusActive
		default:
			*s = StatusInactive
		}
	default:
		*s = StatusInactive
	}
	return nil
}

// Value implements driver.Valuer: เขียนสถานะ string ลงคอลัมน์ int8 (1=active, 0=inactive).
func (s ScheduleStatus) Value() (interface{}, error) {
	if s == StatusActive {
		return int64(1), nil
	}
	return int64(0), nil
}

// ScheduleStatusValue is the active/inactive/draft state stored as int8
// (1=active, 0=inactive, 3=draft) in fs_schedule, mirroring sd_iot_schedule.
type ScheduleStatusValue int8

const (
	StatusValueInactive ScheduleStatusValue = 0
	StatusValueActive   ScheduleStatusValue = 1
	StatusValueDraft    ScheduleStatusValue = 3
)

// String returns the human-readable form ("active"/"inactive"/"draft").
func (s ScheduleStatusValue) String() string {
	switch s {
	case StatusValueActive:
		return "active"
	case StatusValueDraft:
		return "draft"
	default:
		return "inactive"
	}
}

// Scan implements sql.Scanner: รับค่า int8 (หรือ string เดิม) จาก DB แล้วแปลงเป็น int8.
func (s *ScheduleStatusValue) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*s = StatusValueInactive
	case int64:
		*s = ScheduleStatusValue(v)
	case []byte:
		switch string(v) {
		case "1", "active":
			*s = StatusValueActive
		case "3", "draft":
			*s = StatusValueDraft
		default:
			*s = StatusValueInactive
		}
	case string:
		switch v {
		case "1", "active":
			*s = StatusValueActive
		case "3", "draft":
			*s = StatusValueDraft
		default:
			*s = StatusValueInactive
		}
	default:
		*s = StatusValueInactive
	}
	return nil
}

// Value implements driver.Valuer: เขียนสถานะ int8 ลงคอลัมน์ int8.
func (s ScheduleStatusValue) Value() (driver.Value, error) {
	return int64(s), nil
}

// EventType is the kind of event the schedule fires.
type EventType string

const (
	EventTypeDevice EventType = "device" // เปิด/ปิดอุปกรณ์ IoT
	EventTypeEmail  EventType = "email"  // ส่ง email alert (start/stop)
)

// EventAction is the action for the event.
type EventAction string

const (
	ActionOn    EventAction = "ON"    // device ON
	ActionOff   EventAction = "OFF"   // device OFF
	ActionStart EventAction = "start" // email alert start
	ActionStop  EventAction = "stop"  // email alert stop
)

// HistoryStatus is the execution status of a history row.
type HistoryStatus string

const (
	HistoryProcessing HistoryStatus = "processing"
	HistorySuccess    HistoryStatus = "success"
	HistoryFailed     HistoryStatus = "failed"
	HistorySkipped    HistoryStatus = "skipped"
)

// TriggerSource identifies who fired the schedule.
type TriggerSource string

const (
	TriggerAutomatic TriggerSource = "automatic"
	TriggerManual    TriggerSource = "manual"
	TriggerCron      TriggerSource = "cron"
)

// TriggeredBy identifies who asked for the run.
type TriggeredBy string

const (
	TriggeredSystem TriggeredBy = "system"
	TriggeredManual TriggeredBy = "manual"
)

// ============================================================================
// Models (ตาราง fs_* แยกจากระบบเดิมโดยสิ้นเชิง)
// ============================================================================

// IntArray is a []int that can scan/scan a PostgreSQL integer[] column.
// ช่วยให้ GORM อ่าน/เขียนคอลัมน์ integer[] (months/dates) ได้ โดยไม่ต้องพึ่ง
// driver-specific array type — รองรับทั้งรูปแบบ text "{1,15}" และ []byte/[]int64
// ที่ driver บางตัว (pgx/libpq) ส่งกลับมา
type IntArray []int

// Scan implements sql.Scanner: แปลงค่าจาก DB (string/[]byte/[]int64/[]interface{})
// ที่เป็น PostgreSQL array literal แล้วเก็บเป็น []int
func (a *IntArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*a = nil
		return nil
	case []byte:
		return a.parse(string(v))
	case string:
		return a.parse(v)
	case []int64:
		out := make([]int, 0, len(v))
		for _, n := range v {
			out = append(out, int(n))
		}
		*a = out
		return nil
	case []int:
		*a = append([]int(nil), v...)
		return nil
	case []interface{}:
		out := make([]int, 0, len(v))
		for _, e := range v {
			n, err := toInt(e)
			if err != nil {
				return err
			}
			out = append(out, n)
		}
		*a = out
		return nil
	default:
		return fmt.Errorf("IntArray: unsupported Scan type %T", src)
	}
}

// Value implements driver.Valuer: เขียน []int เป็น PostgreSQL array literal.
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return "{}", nil
	}
	parts := make([]string, 0, len(a))
	for _, n := range a {
		parts = append(parts, strconv.Itoa(n))
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

// parse แปลง PostgreSQL array literal ("{1,15}" / "{}") เป็น []int.
func (a *IntArray) parse(s string) error {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return fmt.Errorf("IntArray: invalid array literal %q", s)
	}
	inner := s[1 : len(s)-1]
	if strings.TrimSpace(inner) == "" {
		*a = []int{}
		return nil
	}
	raw := strings.Split(inner, ",")
	out := make([]int, 0, len(raw))
	for _, tok := range raw {
		tok = strings.TrimSpace(tok)
		if tok == "NULL" {
			return errors.New("IntArray: NULL element not supported")
		}
		n, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("IntArray: invalid element %q: %w", tok, err)
		}
		out = append(out, n)
	}
	*a = out
	return nil
}

// toInt converts a driver scalar to int.
func toInt(v interface{}) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case int32:
		return int(n), nil
	case float64:
		return int(n), nil
	case []byte:
		return strconv.Atoi(string(n))
	case string:
		return strconv.Atoi(n)
	default:
		return 0, fmt.Errorf("IntArray: cannot convert %T", v)
	}
}

// Schedule is the main fullschedule record (fs_schedule).
type Schedule struct {
	ID          uuid.UUID           `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string              `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Mode        ScheduleMode        `gorm:"column:mode;type:varchar(20);not null;default:weekly" json:"mode"`
	Status      ScheduleStatusValue `gorm:"column:status;type:int8;not null;default:1" json:"status"`
	TimeStart   string              `gorm:"column:time_start;type:varchar(5);not null;default:00:00" json:"time_start"`
	EventType   EventType           `gorm:"column:event_type;type:varchar(20);not null;default:device" json:"event_type"`
	EventAction EventAction         `gorm:"column:event_action;type:varchar(20);not null;default:ON" json:"event_action"`
	// sd_iot_schedule_style columns — เก็บแบบเดียวกับตารางเดิม
	Event         int8       `gorm:"column:event;type:int8;default:1" json:"event"` // 1=ON, 0=OFF, 2=start, 3=stop
	Sunday        int8       `gorm:"column:sunday;type:int8;default:1" json:"sunday"`
	Monday        int8       `gorm:"column:monday;type:int8;default:0" json:"monday"`
	Tuesday       int8       `gorm:"column:tuesday;type:int8;default:0" json:"tuesday"`
	Wednesday     int8       `gorm:"column:wednesday;type:int8;default:0" json:"wednesday"`
	Thursday      int8       `gorm:"column:thursday;type:int8;default:0" json:"thursday"`
	Friday        int8       `gorm:"column:friday;type:int8;default:0" json:"friday"`
	Saturday      int8       `gorm:"column:saturday;type:int8;default:0" json:"saturday"`
	Months        IntArray     `gorm:"column:months;type:integer[];not null;default:'{}'" json:"months"`
	Dates         IntArray     `gorm:"column:dates;type:integer[];not null;default:'{}'" json:"dates"`
	CronExpr      *string    `gorm:"column:cron_expr;type:varchar(100)" json:"cron_expr,omitempty"`
	ManualTrigger bool       `gorm:"column:manual_trigger;not null;default:false" json:"manual_trigger"`
	LastRunAt     *time.Time `gorm:"column:last_run_at" json:"last_run_at,omitempty"`
	NextRunAt     *time.Time `gorm:"column:next_run_at" json:"next_run_at,omitempty"`
	RunCount      int64      `gorm:"column:run_count;not null;default:0" json:"run_count"`
	SuccessCount  int64      `gorm:"column:success_count;not null;default:0" json:"success_count"`
	FailedCount   int64      `gorm:"column:failed_count;not null;default:0" json:"failed_count"`
	DeviceCount   int64      `gorm:"-" json:"device_count,omitempty"`
	CreatedBy     *uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedBy     *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	Version       int64      `gorm:"column:version;not null;default:1" json:"version"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	// Scope (ขอบเขตการทำงาน Group → Zone → Area) — ทั้งหมดเป็น optional.
	GroupID *uuid.UUID `gorm:"column:group_id;type:uuid;index" json:"group_id,omitempty"`
	ZoneID  *uuid.UUID `gorm:"column:zone_id;type:uuid;index" json:"zone_id,omitempty"`
	AreaID  *uuid.UUID `gorm:"column:area_id;type:uuid;index" json:"area_id,omitempty"`

	Devices  []ScheduleDevice  `gorm:"foreignKey:ScheduleID;references:ID" json:"devices,omitempty"`
	Settings []ScheduleSetting `gorm:"foreignKey:ScheduleID;references:ID" json:"settings,omitempty"`
}

// TableName overrides the default GORM table name.
func (Schedule) TableName() string { return "fs_schedule" }

// ScheduleDevice maps an IoT device to a schedule (fs_schedule_device).
type ScheduleDevice struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ScheduleID uuid.UUID  `gorm:"column:schedule_id;type:uuid;not null" json:"schedule_id"`
	DeviceID   int        `gorm:"column:device_id;not null" json:"device_id"`
	DeviceSN   *string    `gorm:"column:device_sn;type:varchar(255)" json:"device_sn,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName overrides the default GORM table name.
func (ScheduleDevice) TableName() string { return "fs_schedule_device" }

// ScheduleHistory logs each execution (fs_schedule_history).
type ScheduleHistory struct {
	ID            uuid.UUID     `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ScheduleID    uuid.UUID     `gorm:"column:schedule_id;type:uuid;not null" json:"schedule_id"`
	DeviceID      *int          `gorm:"column:device_id" json:"device_id,omitempty"`
	TriggeredBy   TriggeredBy   `gorm:"column:triggered_by;type:varchar(20);not null;default:system" json:"triggered_by"`
	TriggerSource TriggerSource `gorm:"column:trigger_source;type:varchar(20);not null;default:automatic" json:"trigger_source"`
	Status        HistoryStatus `gorm:"column:status;type:varchar(20);not null;default:processing" json:"status"`
	EventAction   EventAction   `gorm:"column:event_action;type:varchar(20);not null;default:''" json:"event_action"`
	Payload       *string       `gorm:"column:payload;type:text" json:"payload,omitempty"`
	Message       *string       `gorm:"column:message;type:text" json:"message,omitempty"`
	DurationMs    int64         `gorm:"column:duration_ms;not null;default:0" json:"duration_ms"`
	RetryCount    int           `gorm:"column:retry_count;not null;default:0" json:"retry_count"`
	ExecutedAt    *time.Time    `gorm:"column:executed_at" json:"executed_at,omitempty"`
	Timezone      *string       `gorm:"column:timezone;type:varchar(50)" json:"timezone,omitempty"`
	Date          *string       `gorm:"column:date;type:varchar(20)" json:"date,omitempty"`
	Time          *string       `gorm:"column:time;type:varchar(20)" json:"time,omitempty"`
	CreatedAt     time.Time     `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at;default:now()" json:"updated_at"`
}

// TableName overrides the default GORM table name.
func (ScheduleHistory) TableName() string { return "fs_schedule_history" }

// ScheduleSetting is a per-schedule configuration (fs_schedule_settings).
type ScheduleSetting struct {
	ID         uuid.UUID      `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ScheduleID uuid.UUID      `gorm:"column:schedule_id;type:uuid;not null" json:"schedule_id"`
	Key        string         `gorm:"column:key;type:varchar(100);not null" json:"key"`
	Value      datatypes.JSON `gorm:"column:value;type:jsonb;not null;default:'{}'" json:"value"`
	CreatedAt  time.Time      `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;default:now()" json:"updated_at"`
}

// TableName overrides the default GORM table name.
func (ScheduleSetting) TableName() string { return "fs_schedule_settings" }

// Group is the top-level organizational unit (fs_groups).
// กลุ่มหลัก เช่น กลุ่มโรงงานภาคตะวันออก
type Group struct {
	ID          uuid.UUID      `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"column:description;type:text" json:"description,omitempty"`
	SortID      int            `gorm:"column:sort_id;type:int;not null;default:0" json:"sort_id"`
	Status      ScheduleStatus `gorm:"column:status;type:varchar(20);not null;default:active" json:"status"`
	CreatedBy   *uuid.UUID     `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedBy   *uuid.UUID     `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;default:now()" json:"updated_at"`
	Version     int64          `gorm:"column:version;not null;default:1" json:"version"`
	DeletedAt   *time.Time     `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	Zones       []Zone         `gorm:"foreignKey:GroupID;references:ID" json:"zones,omitempty"`
}

// TableName overrides the default GORM table name.
func (Group) TableName() string { return "fs_groups" }

// Zone is a unit under a group (fs_zones).
// โซน เช่น อาคารคลังสินค้า
type Zone struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	GroupID     uuid.UUID  `gorm:"column:group_id;type:uuid;not null;index" json:"group_id"`
	GroupName   string     `gorm:"->;column:group_name" json:"group_name"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"column:description;type:text" json:"description,omitempty"`
	SortID      int        `gorm:"column:sort_id;type:int;not null;default:0" json:"sort_id"`
	CreatedBy   *uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedBy   *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	Version     int64      `gorm:"column:version;not null;default:1" json:"version"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	Group       *Group     `gorm:"foreignKey:GroupID;references:ID" json:"-"`
	Areas       []Area     `gorm:"foreignKey:ZoneID;references:ID" json:"areas,omitempty"`
}

// TableName overrides the default GORM table name.
func (Zone) TableName() string { return "fs_zones" }

// Area is a unit under a zone (fs_areas).
// พื้นที่ย่อย เช่น ห้องเซิร์ฟเวอร์ A1
type Area struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ZoneID      uuid.UUID  `gorm:"column:zone_id;type:uuid;not null;index" json:"zone_id"`
	ZoneName    string     `gorm:"->;column:zone_name" json:"zone_name"`
	GroupName   string     `gorm:"->;column:group_name" json:"group_name"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"column:description;type:text" json:"description,omitempty"`
	SortID      int        `gorm:"column:sort_id;type:int;not null;default:0" json:"sort_id"`
	CreatedBy   *uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedBy   *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	Version     int64      `gorm:"column:version;not null;default:1" json:"version"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	Zone        *Zone      `gorm:"foreignKey:ZoneID;references:ID" json:"-"`
}

// TableName overrides the default GORM table name.
func (Area) TableName() string { return "fs_areas" }

// AreaDevice maps an IoT device (sd_iot_device.device_id) to an area (fs_device_area).
// การแมปอุปกรณ์ IoT เข้ากับพื้นที่
type AreaDevice struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AreaID    uuid.UUID  `gorm:"column:area_id;type:uuid;not null;index" json:"area_id"`
	DeviceID  int        `gorm:"column:device_id;not null;index" json:"device_id"`
	DeviceSN  *string    `gorm:"column:device_sn;type:varchar(255)" json:"device_sn,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName overrides the default GORM table name.
func (AreaDevice) TableName() string { return "fs_device_area" }

// ScopeFilter filters schedules by their group/zone/area scope.
// ตัวกรองรายงาน/รายการตามขอบเขต Group/Zone/Area
type ScopeFilter struct {
	GroupID *uuid.UUID `json:"group_id,omitempty"`
	ZoneID  *uuid.UUID `json:"zone_id,omitempty"`
	AreaID  *uuid.UUID `json:"area_id,omitempty"`
}

// ScheduleReport is a summary row used by the report endpoint.
type ScheduleReport struct {
	ScheduleID   uuid.UUID  `json:"schedule_id"`
	Name         string     `json:"name"`
	Mode         string     `json:"mode"`
	TotalRuns    int64      `json:"total_runs"`
	SuccessCount int64      `json:"success_count"`
	FailedCount  int64      `json:"failed_count"`
	SkippedCount int64      `json:"skipped_count"`
	Processing   int64      `json:"processing"`
	SuccessRate  float64    `json:"success_rate"`
	AvgDuration  float64    `json:"avg_duration_ms"`
	LastRunAt    *time.Time `json:"last_run_at,omitempty"`
}
