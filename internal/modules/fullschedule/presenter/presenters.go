package presenter

import (
	"time"

	"icmongolang/internal/modules/fullschedule"
	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// ScheduleCreate is the payload for creating a schedule.
// คำขอสร้าง schedule
type ScheduleCreate struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Mode        string  `json:"mode" validate:"required,oneof=normal full batch"`
	Status      *int8   `json:"status,omitempty" validate:"omitempty,oneof=0 1 3"`
	TimeStart   string  `json:"time_start" validate:"omitempty,datetime=15:04"`
	EventType   string  `json:"event_type" validate:"required,oneof=device email"`
	EventAction string  `json:"event_action" validate:"required,oneof=ON OFF start stop"`
	Event       *int8   `json:"event,omitempty"` // 1=ON,0=OFF,2=start,3=stop (ทางเลือก)
	Sunday      *int8   `json:"sunday,omitempty"`
	Monday      *int8   `json:"monday,omitempty"`
	Tuesday     *int8   `json:"tuesday,omitempty"`
	Wednesday   *int8   `json:"wednesday,omitempty"`
	Thursday    *int8   `json:"thursday,omitempty"`
	Friday      *int8   `json:"friday,omitempty"`
	Saturday    *int8   `json:"saturday,omitempty"`
	Months      fullschedule.IntArray `json:"months"`
	Dates       fullschedule.IntArray `json:"dates"`
	CronExpr    *string `json:"cron_expr,omitempty"`
	DeviceIDs   []int   `json:"device_ids"`
	// Scope — ขอบเขตการทำงาน (ทั้งหมด optional)
	GroupID   *uuid.UUID        `json:"group_id,omitempty"`
	ZoneID    *uuid.UUID        `json:"zone_id,omitempty"`
	AreaID    *uuid.UUID        `json:"area_id,omitempty"`
	CreatedBy *uuid.UUID        `json:"created_by,omitempty"`
	Settings  map[string]string `json:"settings,omitempty"`
}

// ScheduleUpdate is the payload for updating a schedule.
// คำขอแก้ไข schedule
type ScheduleUpdate struct {
	Name        *string `json:"name,omitempty"`
	Mode        *string `json:"mode,omitempty"`
	Status      *int8   `json:"status,omitempty"`
	TimeStart   *string `json:"time_start,omitempty"`
	EventType   *string `json:"event_type,omitempty"`
	EventAction *string `json:"event_action,omitempty"`
	Event       *int8   `json:"event,omitempty"` // 1=ON,0=OFF,2=start,3=stop
	Sunday      *int8   `json:"sunday,omitempty"`
	Monday      *int8   `json:"monday,omitempty"`
	Tuesday     *int8   `json:"tuesday,omitempty"`
	Wednesday   *int8   `json:"wednesday,omitempty"`
	Thursday    *int8   `json:"thursday,omitempty"`
	Friday      *int8   `json:"friday,omitempty"`
	Saturday    *int8   `json:"saturday,omitempty"`
	Months      fullschedule.IntArray `json:"months,omitempty"`
	Dates       fullschedule.IntArray `json:"dates,omitempty"`
	CronExpr    *string `json:"cron_expr,omitempty"`
	DeviceIDs   []int   `json:"device_ids,omitempty"`
	// Scope — ขอบเขตการทำงาน (ทั้งหมด optional)
	GroupID *uuid.UUID `json:"group_id,omitempty"`
	ZoneID  *uuid.UUID `json:"zone_id,omitempty"`
	AreaID  *uuid.UUID `json:"area_id,omitempty"`
}

// ScheduleStatusChange changes the active/inactive/draft state.
// คำขอเปลี่ยนสถานะ (1=active, 0=inactive, 3=draft)
type ScheduleStatusChange struct {
	Status int8 `json:"status" validate:"required,oneof=0 1 3"`
}

// ScheduleDayStatusChange toggles a single weekday flag (sunday..saturday).
// คำขอเปลี่ยนวันทำงานวันเดียว (sunday..saturday)
type ScheduleDayStatusChange struct {
	Day   string `json:"day" validate:"required,oneof=sunday monday tuesday wednesday thursday friday saturday"`
	Value int8   `json:"value" validate:"oneof=0 1"`
}

// ScheduleSettingRequest is a single setting upsert.
// คำขอตั้งค่า schedule
type ScheduleSettingRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// ScheduleResponse is the response shape of a schedule.
// ข้อมูล schedule ที่ส่งกลับ
type ScheduleResponse struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	Mode          string            `json:"mode"`
	Status        int8              `json:"status"`
	TimeStart     string            `json:"time_start"`
	EventType     string            `json:"event_type"`
	EventAction   string            `json:"event_action"`
	Event         int8              `json:"event"`
	Sunday        int8              `json:"sunday"`
	Monday        int8              `json:"monday"`
	Tuesday       int8              `json:"tuesday"`
	Wednesday     int8              `json:"wednesday"`
	Thursday      int8              `json:"thursday"`
	Friday        int8              `json:"friday"`
	Saturday      int8              `json:"saturday"`
	Months        fullschedule.IntArray `json:"months"`
	Dates         fullschedule.IntArray `json:"dates"`
	CronExpr      *string               `json:"cron_expr,omitempty"`
	ManualTrigger bool              `json:"manual_trigger"`
	LastRunAt     *time.Time        `json:"last_run_at,omitempty"`
	NextRunAt     *time.Time        `json:"next_run_at,omitempty"`
	RunCount      int64             `json:"run_count"`
	SuccessCount  int64             `json:"success_count"`
	FailedCount   int64             `json:"failed_count"`
	DeviceCount   int64             `json:"device_count"`
	GroupID       *uuid.UUID        `json:"group_id,omitempty"`
	ZoneID        *uuid.UUID        `json:"zone_id,omitempty"`
	AreaID        *uuid.UUID        `json:"area_id,omitempty"`
	CreatedBy     *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt     helpers.LocalTime `json:"created_at"`
	UpdatedBy     *uuid.UUID        `json:"updated_by,omitempty"`
	UpdatedAt     helpers.LocalTime `json:"updated_at"`
	Version       int64             `json:"version"`
	Devices       []DeviceResponse  `json:"devices,omitempty"`
}

// DeviceResponse describes a mapped IoT device.
// ข้อมูลอุปกรณ์ที่แมป
type DeviceResponse struct {
	DeviceID int    `json:"device_id"`
	DeviceSN string `json:"device_sn,omitempty"`
}

// PaginatedScheduleResponse holds paginated results.
// ผลลัพธ์แบบแบ่งหน้า
type PaginatedScheduleResponse struct {
	Schedules  []*ScheduleResponse `json:"schedules"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}

// ListResult is the settings-style paginated envelope used by the
// listschedulepage / scheduledevicepage endpoints. The shape matches the
// /api/settings/listschedulepage response so the two modules can be shown
// together on the same frontend page.
type ListResult struct {
	Page        int                    `json:"page"`
	CurrentPage int                    `json:"currentPage"`
	PageSize    int                    `json:"pageSize"`
	TotalPages  int64                  `json:"totalPages"`
	Total       int64                  `json:"total"`
	Filter      map[string]interface{} `json:"filter"`
	Data        interface{}            `json:"data"`
}

// NewScheduleListResult builds a settings-style paginated envelope.
// สร้างผลลัพธ์แบบแบ่งหน้า (รูปแบบเดียวกับ /settings/listschedulepage)
func NewScheduleListResult(data interface{}, total int64, page, pageSize int, filter map[string]interface{}) *ListResult {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	tp := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		tp++
	}
	if tp < 1 {
		tp = 1
	}
	if filter == nil {
		filter = map[string]interface{}{}
	}
	return &ListResult{
		Page:        page,
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  tp,
		Total:       total,
		Filter:      filter,
		Data:        data,
	}
}

// HistoryResponse is one execution history row.
// ประวัติการทำงาน 1 รายการ
type HistoryResponse struct {
	ID            uuid.UUID         `json:"id"`
	ScheduleID    uuid.UUID         `json:"schedule_id"`
	DeviceID      *int              `json:"device_id,omitempty"`
	TriggeredBy   string            `json:"triggered_by"`
	TriggerSource string            `json:"trigger_source"`
	Status        string            `json:"status"`
	EventAction   string            `json:"event_action"`
	Payload       *string           `json:"payload,omitempty"`
	Message       *string           `json:"message,omitempty"`
	DurationMs    int64             `json:"duration_ms"`
	RetryCount    int               `json:"retry_count"`
	ExecutedAt    *time.Time        `json:"executed_at,omitempty"`
	Timezone      *string           `json:"timezone,omitempty"`
	Date          *string           `json:"date,omitempty"`
	Time          *string           `json:"time,omitempty"`
	CreatedAt     helpers.LocalTime `json:"created_at"`
}

// PaginatedHistoryResponse holds paginated history.
// ประวัติแบบแบ่งหน้า
type PaginatedHistoryResponse struct {
	History []*HistoryResponse `json:"history"`
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	PerPage int                `json:"per_page"`
}

// ReportResponse is one aggregated row.
// รายงานสรุป 1 รายการ
type ReportResponse struct {
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

// TriggerResponse is the manual trigger result.
// ผลการสั่งงานด้วยตนเอง
type TriggerResponse struct {
	Message   string    `json:"message"`
	HistoryID uuid.UUID `json:"history_id"`
}

// SettingResponse is a schedule setting.
// ค่าตั้ง schedule
type SettingResponse struct {
	ID         uuid.UUID         `json:"id"`
	ScheduleID uuid.UUID         `json:"schedule_id"`
	Key        string            `json:"key"`
	Value      string            `json:"value"`
	UpdatedAt  helpers.LocalTime `json:"updated_at"`
}

// IoTDeviceBrief is a lightweight device for the picker endpoint.
// ข้อมูลอุปกรณ์ย่อสำหรับเลือกแมป
type IoTDeviceBrief struct {
	DeviceID   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	SN         string `json:"sn"`
	Status     int    `json:"status"`
}
