package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	"icmongolang/internal/modules/fullschedule/executor"
	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// Errors returned by the usecase.
var (
	ErrModeInvalid      = errors.New("fullschedule: mode must be normal, full or batch")
	ErrStatusInvalid    = errors.New("fullschedule: status must be active or inactive")
	ErrEventInvalid     = errors.New("fullschedule: invalid event_type/event_action combination")
	ErrTimeStartInvalid = errors.New("fullschedule: time_start must be HH:MM between 00:00-23:59")
	ErrDaysRequired     = errors.New("fullschedule: mode=normal requires at least one weekday (sunday..saturday=1)")
	ErrCalendarRequired = errors.New("fullschedule: mode=full requires months (1-12) and dates (1-31)")
	ErrCronRequired     = errors.New("fullschedule: mode=batch requires valid cron_expr")
	ErrNoSchedule       = errors.New("fullschedule: schedule not found")
)

// fullscheduleUseCase implements fullschedule.FullScheduleUseCaseI.
type fullscheduleUseCase struct {
	usecase.UseCase[fullschedule.Schedule]
	pgRepo   fullschedule.FullSchedulePgRepository
	masterI  fullschedule.MasterUseCaseI
	exec     *executor.Executor
	logger   logger.Logger
	timezone *time.Location
}

// CreateFullScheduleUseCaseI creates the usecase instance.
// สร้างยูสเคสสำหรับโมดูล Full Schedule
func CreateFullScheduleUseCaseI(
	pgRepo fullschedule.FullSchedulePgRepository,
	exec *executor.Executor,
	masterI fullschedule.MasterUseCaseI,
	cfg *config.Config,
	logger logger.Logger,
) fullschedule.FullScheduleUseCaseI {
	loc := time.Local
	if cfg != nil && cfg.Server.Timezone != "" {
		if l, err := time.LoadLocation(cfg.Server.Timezone); err == nil {
			loc = l
		}
	}
	return &fullscheduleUseCase{
		UseCase:  usecase.CreateUseCase[fullschedule.Schedule](pgRepo, cfg, logger),
		pgRepo:   pgRepo,
		masterI:  masterI,
		exec:     exec,
		logger:   logger,
		timezone: loc,
	}
}

// validate checks the schedule fields based on its mode.
func validate(s *fullschedule.Schedule) error {
	switch s.Mode {
	case fullschedule.ModeWeekly, fullschedule.ModeFull:
		if _, _, err := parseTimeStart(s.TimeStart); err != nil {
			return ErrTimeStartInvalid
		}
	case fullschedule.ModeBatch:
		// cron_expr validated by calcNextRun below.
	default:
		return ErrModeInvalid
	}

	if s.Status != fullschedule.StatusValueActive && s.Status != fullschedule.StatusValueInactive && s.Status != fullschedule.StatusValueDraft {
		return ErrStatusInvalid
	}

	switch s.EventType {
	case fullschedule.EventTypeDevice:
		if s.EventAction != fullschedule.ActionOn && s.EventAction != fullschedule.ActionOff {
			return ErrEventInvalid
		}
	case fullschedule.EventTypeEmail:
		if s.EventAction != fullschedule.ActionStart && s.EventAction != fullschedule.ActionStop {
			return ErrEventInvalid
		}
	default:
		return ErrEventInvalid
	}

	switch s.Mode {
	case fullschedule.ModeWeekly:
		if !anyWeekday(s) {
			return ErrDaysRequired
		}
	case fullschedule.ModeFull:
		if len(s.Months) == 0 || len(s.Dates) == 0 {
			return ErrCalendarRequired
		}
		for _, m := range s.Months {
			if m < 1 || m > 12 {
				return ErrCalendarRequired
			}
		}
		for _, d := range s.Dates {
			if d < 1 || d > 31 {
				return ErrCalendarRequired
			}
		}
	case fullschedule.ModeBatch:
		if s.CronExpr == nil || strings.TrimSpace(*s.CronExpr) == "" {
			return ErrCronRequired
		}
	}
	return nil
}

func (u *fullscheduleUseCase) applyNextRun(s *fullschedule.Schedule, now time.Time) error {
	if s.Status == fullschedule.StatusValueInactive || s.Status == fullschedule.StatusValueDraft {
		s.NextRunAt = nil
		return nil
	}
	next, err := u.CalculateNextRun(s, now)
	if err != nil {
		return err
	}
	s.NextRunAt = next
	return nil
}

// CalculateNextRun delegates to the recurrence calculator.
// คำนวณ `next_run_at` จากโหมด
func (u *fullscheduleUseCase) CalculateNextRun(s *fullschedule.Schedule, now time.Time) (*time.Time, error) {
	return calcNextRun(s, now.In(u.timezone))
}

// Create validates, computes next_run and persists the schedule with devices.
// สร้าง schedule ใหม่
func (u *fullscheduleUseCase) Create(ctx context.Context, s *fullschedule.Schedule) (*fullschedule.Schedule, error) {
	if err := validate(s); err != nil {
		return nil, err
	}
	if s.GroupID != nil || s.ZoneID != nil || s.AreaID != nil {
		if u.masterI == nil {
			return nil, fullschedule.ErrMasterStoreNotWired
		}
		if err := u.masterI.ValidateScope(ctx, s.GroupID, s.ZoneID, s.AreaID); err != nil {
			return nil, err
		}
	}
	now := time.Now().In(u.timezone)
	if err := u.applyNextRun(s, now); err != nil {
		return nil, err
	}
	syncSDIoTFields(s)
	created, err := u.pgRepo.Create(ctx, s)
	if err != nil {
		return nil, err
	}
	if err := u.replaceDevices(ctx, created.ID, deviceIDsFrom(created.Devices)); err != nil {
		u.logger.Errorf("fullschedule: replace devices failed: %v", err)
	}
	return u.Get(ctx, created.ID)
}

// Get loads a schedule including devices and settings.
// ดึง schedule พร้อมอุปกรณ์และการตั้งค่า
func (u *fullscheduleUseCase) Get(ctx context.Context, id uuid.UUID) (*fullschedule.Schedule, error) {
	s, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	devices, err := u.pgRepo.GetDevices(ctx, id)
	if err != nil {
		return nil, err
	}
	settings, err := u.pgRepo.GetSettings(ctx, id)
	if err != nil {
		return nil, err
	}
	s.Devices = valueDevices(devices)
	s.Settings = valueSettings(settings)
	return s, nil
}

// GetMulti lists schedules (paged).
// รายการ schedule แบบแบ่งหน้า
func (u *fullscheduleUseCase) GetMulti(ctx context.Context, limit, offset int) ([]*fullschedule.Schedule, error) {
	return u.pgRepo.GetMulti(ctx, limit, offset)
}

// ListPage lists schedules with filters + pagination (settings-style display).
// รายการ schedule พร้อมตัวกรองและแบ่งหน้า
func (u *fullscheduleUseCase) ListPage(ctx context.Context, f fullschedule.ScheduleFilter, page, pageSize int) ([]*fullschedule.Schedule, int64, error) {
	return u.pgRepo.ListPage(ctx, f, page, pageSize)
}

// ListAll lists schedules with filters (no pagination).
// รายการ schedule ทั้งหมดพร้อมตัวกรอง
func (u *fullscheduleUseCase) ListAll(ctx context.Context, f fullschedule.ScheduleFilter) ([]*fullschedule.Schedule, error) {
	return u.pgRepo.ListAll(ctx, f)
}

// ListScheduleDevicePage lists schedule-device mappings with pagination.
// รายการแมป schedule ↔ device แบบแบ่งหน้า
func (u *fullscheduleUseCase) ListScheduleDevicePage(ctx context.Context, f fullschedule.ScheduleDeviceFilter, page, pageSize int) ([]*fullschedule.ScheduleDeviceRow, int64, error) {
	return u.pgRepo.ListScheduleDevicePage(ctx, f, page, pageSize)
}

// Update patches a schedule and recomputes next_run.
// แก้ไข schedule
func (u *fullscheduleUseCase) Update(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*fullschedule.Schedule, error) {
	existing, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// Build an updated view of fields for validation.
	mode := existing.Mode
	status := existing.Status
	timeStart := existing.TimeStart
	eventType := existing.EventType
	eventAction := existing.EventAction
	sunday := existing.Sunday
	monday := existing.Monday
	tuesday := existing.Tuesday
	wednesday := existing.Wednesday
	thursday := existing.Thursday
	friday := existing.Friday
	saturday := existing.Saturday
	months := existing.Months
	dates := existing.Dates
	cronExpr := existing.CronExpr

	if v, ok := values["mode"].(string); ok {
		mode = fullschedule.ScheduleMode(v)
	}
	if v, ok := values["status"]; ok {
		switch tv := v.(type) {
		case int8:
			status = fullschedule.ScheduleStatusValue(tv)
		case int:
			status = fullschedule.ScheduleStatusValue(tv)
		case int64:
			status = fullschedule.ScheduleStatusValue(tv)
		case float64:
			status = fullschedule.ScheduleStatusValue(int8(tv))
		case string:
			switch tv {
			case "1", "active":
				status = fullschedule.StatusValueActive
			case "3", "draft":
				status = fullschedule.StatusValueDraft
			case "0", "inactive":
				status = fullschedule.StatusValueInactive
			}
		}
	}
	if v, ok := values["time_start"].(string); ok {
		timeStart = v
	}
	if v, ok := values["event_type"].(string); ok {
		eventType = fullschedule.EventType(v)
	}
	if v, ok := values["event_action"].(string); ok {
		eventAction = fullschedule.EventAction(v)
	}
	if v, ok := values["sunday"].(int8); ok {
		sunday = v
	}
	if v, ok := values["monday"].(int8); ok {
		monday = v
	}
	if v, ok := values["tuesday"].(int8); ok {
		tuesday = v
	}
	if v, ok := values["wednesday"].(int8); ok {
		wednesday = v
	}
	if v, ok := values["thursday"].(int8); ok {
		thursday = v
	}
	if v, ok := values["friday"].(int8); ok {
		friday = v
	}
	if v, ok := values["saturday"].(int8); ok {
		saturday = v
	}
	if v, ok := values["months"].(fullschedule.IntArray); ok {
		months = v
	}
	if v, ok := values["dates"].(fullschedule.IntArray); ok {
		dates = v
	}
	if v, ok := values["cron_expr"].(*string); ok {
		cronExpr = v
	}

	probe := &fullschedule.Schedule{
		Mode:        mode,
		Status:      status,
		TimeStart:   timeStart,
		EventType:   eventType,
		EventAction: eventAction,
		Sunday:      sunday,
		Monday:      monday,
		Tuesday:     tuesday,
		Wednesday:   wednesday,
		Thursday:    thursday,
		Friday:      friday,
		Saturday:    saturday,
		Months:      months,
		Dates:       dates,
		CronExpr:    cronExpr,
	}
	if err := validate(probe); err != nil {
		return nil, err
	}

	// Scope validation: effective scope = new values, else existing.
	scopeGroup, scopeZone, scopeArea := existing.GroupID, existing.ZoneID, existing.AreaID
	if v, ok := values["group_id"].(uuid.UUID); ok {
		scopeGroup = &v
	}
	if v, ok := values["zone_id"].(uuid.UUID); ok {
		scopeZone = &v
	}
	if v, ok := values["area_id"].(uuid.UUID); ok {
		scopeArea = &v
	}
	if scopeGroup != nil || scopeZone != nil || scopeArea != nil {
		if u.masterI == nil {
			return nil, fullschedule.ErrMasterStoreNotWired
		}
		if err := u.masterI.ValidateScope(ctx, scopeGroup, scopeZone, scopeArea); err != nil {
			return nil, err
		}
	}

	// Sync the sd_iot_schedule_style int8 columns into the update values so the
	// DB keeps `event` consistent with event_type/event_action; the weekday flags
	// (sunday..saturday) are already the source of truth.
	syncSDIoTFields(probe)
	values["event"] = probe.Event
	values["sunday"] = probe.Sunday
	values["monday"] = probe.Monday
	values["tuesday"] = probe.Tuesday
	values["wednesday"] = probe.Wednesday
	values["thursday"] = probe.Thursday
	values["friday"] = probe.Friday
	values["saturday"] = probe.Saturday
	values["status"] = probe.Status

	values["updated_at"] = time.Now().In(u.timezone)
	updated, err := u.pgRepo.Update(ctx, existing, values)
	if err != nil {
		return nil, err
	}

	// Recompute next_run for the updated schedule.
	if deviceIDs, ok := values["_device_ids"].([]int); ok {
		delete(values, "_device_ids")
		if err := u.replaceDevices(ctx, id, deviceIDs); err != nil {
			u.logger.Errorf("fullschedule: replace devices failed: %v", err)
		}
	}
	if err := u.applyNextRun(updated, time.Now().In(u.timezone)); err != nil {
		u.logger.Errorf("fullschedule: recompute next_run failed: %v", err)
	} else {
		_, _ = u.pgRepo.Update(ctx, updated, map[string]interface{}{"next_run_at": updated.NextRunAt})
	}
	return u.Get(ctx, id)
}

// Delete soft-deletes a schedule.
// ลบ schedule (soft delete)
func (u *fullscheduleUseCase) Delete(ctx context.Context, id uuid.UUID) (*fullschedule.Schedule, error) {
	existing, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now().In(u.timezone)
	_, err = u.pgRepo.Update(ctx, existing, map[string]interface{}{"deleted_at": now})
	if err != nil {
		return nil, err
	}
	return u.Get(ctx, id)
}

// SetStatus toggles active/inactive.
// ตั้งค่าสถานะ schedule
func (u *fullscheduleUseCase) SetStatus(ctx context.Context, id uuid.UUID, status fullschedule.ScheduleStatusValue) (*fullschedule.Schedule, error) {
	if status != fullschedule.StatusValueActive && status != fullschedule.StatusValueInactive && status != fullschedule.StatusValueDraft {
		return nil, ErrStatusInvalid
	}
	existing, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	updated, err := u.pgRepo.Update(ctx, existing, map[string]interface{}{"status": status})
	if err != nil {
		return nil, err
	}
	// Recompute next_run on activation, clear on deactivation.
	if err := u.applyNextRun(updated, time.Now().In(u.timezone)); err != nil {
		u.logger.Errorf("fullschedule: recompute next_run on status change failed: %v", err)
	} else {
		_, _ = u.pgRepo.Update(ctx, updated, map[string]interface{}{"next_run_at": updated.NextRunAt})
	}
	return u.Get(ctx, id)
}

// Trigger forces execution now ignoring day/time conditions.
// สั่งงานด้วยตนเอง: ข้ามเงื่อนไขวัน/เวลา
func (u *fullscheduleUseCase) Trigger(ctx context.Context, id uuid.UUID, by fullschedule.TriggeredBy) (*fullschedule.ScheduleHistory, error) {
	s, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.Mode == fullschedule.ModeBatch {
		return nil, errors.New("fullschedule: batch/cron schedules cannot be manually triggered; use POST /run instead")
	}
	// Mark manual flag so the execution is recorded and never rescheduled twice.
	now := time.Now().In(u.timezone)
	past := now.Add(-1 * time.Second)
	_, err = u.pgRepo.Update(ctx, s, map[string]interface{}{
		"manual_trigger": true,
		"next_run_at":    past,
	})
	if err != nil {
		return nil, err
	}
	s.ManualTrigger = true
	s.NextRunAt = &past

	var last *fullschedule.ScheduleHistory
	u.runLocked(ctx, s, fullschedule.TriggerManual, fullschedule.TriggeredManual, func(h *fullschedule.ScheduleHistory) {
		last = h
	})
	return last, nil
}

// GetDueSchedules returns schedules due now.
// ดึง schedule ที่ถึงกำหนด
func (u *fullscheduleUseCase) GetDueSchedules(ctx context.Context, now time.Time) ([]*fullschedule.Schedule, error) {
	return u.pgRepo.GetDueSchedules(ctx, now)
}

// Run executes a schedule and records history. It is idempotent via lock.
// ทำงาน schedule พร้อมบันทึกประวัติและกันการสั่งซ้ำ
func (u *fullscheduleUseCase) Run(ctx context.Context, s *fullschedule.Schedule, source fullschedule.TriggerSource, by fullschedule.TriggeredBy) error {
	return u.runLocked(ctx, s, source, by, nil)
}

// runLocked is the shared execution path (scheduler + manual trigger).
func (u *fullscheduleUseCase) runLocked(
	ctx context.Context,
	s *fullschedule.Schedule,
	source fullschedule.TriggerSource,
	by fullschedule.TriggeredBy,
	onResult func(*fullschedule.ScheduleHistory),
) error {
	now := time.Now().In(u.timezone)

	deviceIDs, err := u.deviceIDs(ctx, s)
	if err != nil {
		return err
	}

	var results []executor.Result
	if u.exec != nil {
		results = u.exec.Execute(ctx, s, deviceIDs)
	} else {
		results = []executor.Result{{OK: false, Message: "executor not wired"}}
	}

	// Record one history row per result (email produces a single row).
	var firstHist *fullschedule.ScheduleHistory
	allOK := true
	for _, r := range results {
		status := fullschedule.HistorySuccess
		if !r.OK {
			status = fullschedule.HistoryFailed
			allOK = false
		}
		msg := r.Message
		payload := r.Payload
		date := now.Format("2006-01-02")
		tim := now.Format("15:04:05")
		locName := u.timezone.String()
		h := &fullschedule.ScheduleHistory{
			ScheduleID:    s.ID,
			DeviceID:      r.DeviceID,
			TriggeredBy:   by,
			TriggerSource: source,
			Status:        status,
			EventAction:   s.EventAction,
			Payload:       &payload,
			Message:       &msg,
			DurationMs:    r.Duration.Milliseconds(),
			ExecutedAt:    &now,
			Timezone:      &locName,
			Date:          &date,
			Time:          &tim,
		}
		created, cerr := u.pgRepo.CreateHistory(ctx, h)
		if cerr != nil {
			u.logger.Errorf("fullschedule: create history failed: %v", cerr)
			continue
		}
		if onResult != nil {
			onResult(created)
		}
		if firstHist == nil {
			firstHist = created
		}
	}
	if firstHist == nil {
		return errors.New("fullschedule: failed to record history")
	}

	// Update schedule counters + next_run.
	next, nextErr := u.CalculateNextRun(s, now)
	s.RunCount++
	if allOK {
		s.SuccessCount++
	} else {
		s.FailedCount++
	}
	s.LastRunAt = &now
	s.NextRunAt = next
	if nextErr != nil {
		s.NextRunAt = nil
		u.logger.Errorf("fullschedule: next_run computation failed for %s: %v", s.ID, nextErr)
	}
	values := map[string]interface{}{
		"manual_trigger": false,
		"last_run_at":    s.LastRunAt,
		"next_run_at":    s.NextRunAt,
		"run_count":      s.RunCount,
		"success_count":  s.SuccessCount,
		"failed_count":   s.FailedCount,
	}
	if _, err := u.pgRepo.Update(ctx, s, values); err != nil {
		u.logger.Errorf("fullschedule: update counters failed: %v", err)
	}
	return nil
}

// SetSetting upserts a per-schedule setting.
// ตั้งค่าของ schedule
func (u *fullscheduleUseCase) SetSetting(ctx context.Context, scheduleID uuid.UUID, key string, value []byte) (*fullschedule.ScheduleSetting, error) {
	if _, err := u.pgRepo.Get(ctx, scheduleID); err != nil {
		return nil, ErrNoSchedule
	}
	if strings.TrimSpace(key) == "" {
		return nil, errors.New("fullschedule: setting key is required")
	}
	check := &fullschedule.ScheduleSetting{
		ScheduleID: scheduleID,
		Key:        key,
		Value:      value,
	}
	return u.pgRepo.UpsertSetting(ctx, check)
}

// GetSettings returns all settings of a schedule.
// ดูค่าตั้งของ schedule
func (u *fullscheduleUseCase) GetSettings(ctx context.Context, scheduleID uuid.UUID) ([]*fullschedule.ScheduleSetting, error) {
	return u.pgRepo.GetSettings(ctx, scheduleID)
}

// GetHistoryBySchedule lists history of one schedule.
// ดูประวัติของ schedule
func (u *fullscheduleUseCase) GetHistoryBySchedule(ctx context.Context, scheduleID uuid.UUID, limit, offset int) ([]*fullschedule.ScheduleHistory, error) {
	return u.pgRepo.GetHistoryBySchedule(ctx, scheduleID, limit, offset)
}

// GetHistory lists history across all schedules.
// ดูประวัติทั้งหมด
func (u *fullscheduleUseCase) GetHistory(ctx context.Context, limit, offset int, status *fullschedule.HistoryStatus) ([]*fullschedule.ScheduleHistory, error) {
	return u.pgRepo.GetHistory(ctx, limit, offset, status)
}

// GetReport aggregates history for reporting.
// รายงานสรุปความสำเร็จ
func (u *fullscheduleUseCase) GetReport(ctx context.Context, from, to *time.Time, scheduleID *uuid.UUID, scope fullschedule.ScopeFilter) ([]*fullschedule.ScheduleReport, error) {
	return u.pgRepo.GetReport(ctx, from, to, scheduleID, scope)
}

// Count returns total schedules.
// นับจำนวน schedule
func (u *fullscheduleUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}

// ListDevices exposes read-only IoT device lookup for the mapping picker.
// รายการอุปกรณ์ IoT สำหรับเลือกแมป
func (u *fullscheduleUseCase) ListDevices(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*iotmodels.Device, int64, error) {
	return u.pgRepo.ListDevices(ctx, filter, page, pageSize)
}

// replaceDevices validates and replaces the device mapping.
func (u *fullscheduleUseCase) replaceDevices(ctx context.Context, scheduleID uuid.UUID, deviceIDs []int) error {
	if len(deviceIDs) == 0 {
		return u.pgRepo.ReplaceDevices(ctx, scheduleID, nil)
	}
	valid := make([]int, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		if id <= 0 {
			continue
		}
		valid = append(valid, id)
	}
	return u.pgRepo.ReplaceDevices(ctx, scheduleID, valid)
}

// deviceIDs returns the device ids for a schedule. Manual device_ids override
// the scope; otherwise devices are resolved from the group/zone/area scope.
func (u *fullscheduleUseCase) deviceIDs(ctx context.Context, s *fullschedule.Schedule) ([]int, error) {
	if len(s.Devices) > 0 {
		out := make([]int, 0, len(s.Devices))
		for _, d := range s.Devices {
			out = append(out, d.DeviceID)
		}
		return out, nil
	}
	if u.masterI != nil {
		if ids, err := u.masterI.GetDeviceIDsByScope(ctx, s.GroupID, s.ZoneID, s.AreaID); err == nil {
			if len(ids) > 0 {
				return ids, nil
			}
		} else if !errors.Is(err, fullschedule.ErrNoScope) {
			return nil, err
		}
	}
	ids, err := u.pgRepo.ListDeviceIDs(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fullschedule.ErrNoTargetDevices
	}
	return ids, nil
}

func deviceIDsFrom(devices []fullschedule.ScheduleDevice) []int {
	out := make([]int, 0, len(devices))
	for _, d := range devices {
		out = append(out, d.DeviceID)
	}
	return out
}

// syncSDIoTFields derives the sd_iot_schedule_style `event` int8 column from the
// richer event_type/event_action representation before persisting. The weekday
// columns (sunday..saturday int8) are now the source of truth and are left as-is.
func syncSDIoTFields(s *fullschedule.Schedule) {
	// event int8: device ON=1/OFF=0, email start=2/stop=3
	switch s.EventType {
	case fullschedule.EventTypeDevice:
		if s.EventAction == fullschedule.ActionOff {
			s.Event = 0
		} else {
			s.Event = 1
		}
	case fullschedule.EventTypeEmail:
		if s.EventAction == fullschedule.ActionStop {
			s.Event = 3
		} else if s.EventAction == fullschedule.ActionStart {
			s.Event = 2
		} else {
			s.Event = 1
		}
	default:
		s.Event = 1
	}
}

// anyWeekday reports whether at least one weekday flag (sunday..saturday) is on.
func anyWeekday(s *fullschedule.Schedule) bool {
	return s.Sunday == 1 || s.Monday == 1 || s.Tuesday == 1 ||
		s.Wednesday == 1 || s.Thursday == 1 || s.Friday == 1 || s.Saturday == 1
}

func valueDevices(rows []*fullschedule.ScheduleDevice) []fullschedule.ScheduleDevice {
	out := make([]fullschedule.ScheduleDevice, 0, len(rows))
	for _, r := range rows {
		if r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func valueSettings(rows []*fullschedule.ScheduleSetting) []fullschedule.ScheduleSetting {
	out := make([]fullschedule.ScheduleSetting, 0, len(rows))
	for _, r := range rows {
		if r != nil {
			out = append(out, *r)
		}
	}
	return out
}

// String gives a compact debug representation.
func (u *fullscheduleUseCase) String() string {
	return fmt.Sprintf("fullschedule-use-case(tz=%s)", u.timezone)
}
