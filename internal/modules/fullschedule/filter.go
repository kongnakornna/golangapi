package fullschedule

import (
	"strings"

	"github.com/google/uuid"
)

// ScheduleFilter carries the list filters shared by the display endpoints
// (listschedulepage / scheduleall). Fields mirror the settings module query
// params (keyword, start, mode, status, ...) so both systems can be driven by
// the same request style: ?page=2&pageSize=10&start=.
// ตัวกรองรายการ schedule (เหมือนรูปแบบของโมดูล settings)
type ScheduleFilter struct {
	Keyword     string
	Mode        *ScheduleMode
	Status      *ScheduleStatusValue
	EventType   *EventType
	EventAction *EventAction
	Start       string
	Scope       ScopeFilter
	Sort        string
}

// ScheduleDeviceFilter filters the schedule-device mapping list.
// ตัวกรองรายการแมป schedule ↔ device
type ScheduleDeviceFilter struct {
	Keyword    string
	ScheduleID *uuid.UUID
	DeviceID   *int
	Status     *ScheduleStatusValue
	Sort       string
}

// ScheduleDeviceRow is one schedule-device mapping used by
// /fullschedule/scheduledevicepage (settings-style display).
// รายการแมป schedule ↔ device 1 แถว
type ScheduleDeviceRow struct {
	ScheduleID   uuid.UUID `json:"schedule_id"`
	ScheduleName string    `json:"schedule_name"`
	Mode         string    `json:"mode"`
	Status       int8    `json:"status"`
	Start        string    `json:"start"`
	EventAction  string    `json:"event_action"`
	DeviceID     int       `json:"device_id"`
	DeviceName   *string   `json:"device_name,omitempty"`
	SN           *string   `json:"sn,omitempty"`
}

// FilterMap echoes the applied filters back to the client.
func (f ScheduleFilter) FilterMap() map[string]interface{} {
	out := map[string]interface{}{}
	if f.Keyword != "" {
		out["keyword"] = f.Keyword
	}
	if f.Mode != nil {
		out["mode"] = string(*f.Mode)
	}
	if f.Status != nil {
		out["status"] = int(*f.Status)
	}
	if f.EventType != nil {
		out["event_type"] = string(*f.EventType)
	}
	if f.EventAction != nil {
		out["event_action"] = string(*f.EventAction)
	}
	if f.Start != "" {
		out["start"] = f.Start
	}
	if f.Scope.GroupID != nil {
		out["group_id"] = f.Scope.GroupID.String()
	}
	if f.Scope.ZoneID != nil {
		out["zone_id"] = f.Scope.ZoneID.String()
	}
	if f.Scope.AreaID != nil {
		out["area_id"] = f.Scope.AreaID.String()
	}
	return out
}

// FilterMap echoes the applied filters back to the client.
func (f ScheduleDeviceFilter) FilterMap() map[string]interface{} {
	out := map[string]interface{}{}
	if f.Keyword != "" {
		out["keyword"] = f.Keyword
	}
	if f.ScheduleID != nil {
		out["schedule_id"] = f.ScheduleID.String()
	}
	if f.DeviceID != nil {
		out["device_id"] = *f.DeviceID
	}
	if f.Status != nil {
		out["status"] = int(*f.Status)
	}
	return out
}

// scheduleSortCols whitelists ORDER BY sources for schedule listing.
var scheduleSortCols = map[string]string{
	"schedule_id":   "id",
	"schedule_name": "name",
	"name":          "name",
	"mode":          "mode",
	"status":        "status",
	"time_start":    "time_start",
	"start":         "time_start",
	"event_type":    "event_type",
	"event_action":  "event_action",
	"created_at":    "created_at",
	"createddate":   "created_at",
	"updateddate":   "updated_at",
	"updated_at":    "updated_at",
	"next_run_at":   "next_run_at",
	"last_run_at":   "last_run_at",
}

// ParseScheduleSort converts "name-DESC" -> "name DESC" using the whitelist.
func ParseScheduleSort(sort string) (string, bool) {
	if strings.TrimSpace(sort) == "" {
		return "created_at DESC", true
	}
	parts := strings.SplitN(sort, "-", 2)
	if len(parts) != 2 {
		return "", false
	}
	col, ok := scheduleSortCols[parts[0]]
	if !ok {
		return "", false
	}
	dir := strings.ToUpper(parts[1])
	if dir != "ASC" && dir != "DESC" {
		return "", false
	}
	return col + " " + dir, true
}

// scheduleDeviceSortCols whitelists ORDER BY sources for the device mapping list.
var scheduleDeviceSortCols = map[string]string{
	"schedule_id":   "s.id",
	"schedule_name": "s.name",
	"start":         "s.time_start",
	"device_id":     "sd.device_id",
	"device_name":   "d.device_name",
	"status":        "s.status",
	"created_at":    "s.created_at",
}

// ParseScheduleDeviceSort converts "start-DESC" -> "s.time_start DESC".
func ParseScheduleDeviceSort(sort string) (string, bool) {
	if strings.TrimSpace(sort) == "" {
		return "", false
	}
	parts := strings.SplitN(sort, "-", 2)
	if len(parts) != 2 {
		return "", false
	}
	col, ok := scheduleDeviceSortCols[parts[0]]
	if !ok {
		return "", false
	}
	dir := strings.ToUpper(parts[1])
	if dir != "ASC" && dir != "DESC" {
		return "", false
	}
	return col + " " + dir, true
}
