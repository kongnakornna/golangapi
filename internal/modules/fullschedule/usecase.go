package fullschedule

import (
	"context"
	"errors"
	"time"

	iotmodels "icmongolang/internal/modules/iot/models"

	"github.com/google/uuid"
)

// ErrDeviceStoreNotWired is returned when the device store is unavailable.

var ErrDeviceStoreNotWired = errors.New("fullschedule: device store not wired")

// ErrInvalidSort is returned when the sort query param is not in the whitelist.
// เกิดเมื่อ sort ไม่อยู่ในรายการที่อนุญาต (กัน SQL injection)
var ErrInvalidSort = errors.New("fullschedule: invalid sort option")

// DeviceListProvider exposes read-only IoT device lookup for the mapping picker.
type DeviceListProvider interface {
	ListDevices(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*iotmodels.Device, int64, error)
}

// FullScheduleUseCaseI defines business logic for the Full Schedule module.
// อินเทอร์เฟซธุรกิจสำหรับโมดูล Full Schedule
type FullScheduleUseCaseI interface {
	Create(ctx context.Context, s *Schedule) (*Schedule, error)
	Get(ctx context.Context, id uuid.UUID) (*Schedule, error)
	GetMulti(ctx context.Context, limit, offset int) ([]*Schedule, error)
	// ListPage lists schedules with filters + pagination (settings-style display).
	ListPage(ctx context.Context, f ScheduleFilter, page, pageSize int) ([]*Schedule, int64, error)
	// ListAll lists schedules with filters (no pagination).
	ListAll(ctx context.Context, f ScheduleFilter) ([]*Schedule, error)
	// ListScheduleDevicePage lists schedule-device mappings with pagination.
	ListScheduleDevicePage(ctx context.Context, f ScheduleDeviceFilter, page, pageSize int) ([]*ScheduleDeviceRow, int64, error)
	Delete(ctx context.Context, id uuid.UUID) (*Schedule, error)
	Update(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*Schedule, error)
	SetStatus(ctx context.Context, id uuid.UUID, status ScheduleStatusValue) (*Schedule, error)
	// Trigger forces immediate execution ignoring day/time/month conditions.
	Trigger(ctx context.Context, id uuid.UUID, by TriggeredBy) (*ScheduleHistory, error)
	// Run executes a schedule right now (used by scheduler + trigger).
	Run(ctx context.Context, s *Schedule, source TriggerSource, by TriggeredBy) error
	// SetSetting upserts a per-schedule setting.
	SetSetting(ctx context.Context, scheduleID uuid.UUID, key string, value []byte) (*ScheduleSetting, error)
	// GetSettings returns all settings for a schedule.
	GetSettings(ctx context.Context, scheduleID uuid.UUID) ([]*ScheduleSetting, error)
	// GetHistoryBySchedule returns history of one schedule.
	GetHistoryBySchedule(ctx context.Context, scheduleID uuid.UUID, limit, offset int) ([]*ScheduleHistory, error)
	// GetHistory returns history across all schedules.
	GetHistory(ctx context.Context, limit, offset int, status *HistoryStatus) ([]*ScheduleHistory, error)
	// GetReport aggregates success/failure for reporting.
	GetReport(ctx context.Context, from, to *time.Time, scheduleID *uuid.UUID, scope ScopeFilter) ([]*ScheduleReport, error)
	// GetDueSchedules returns schedules whose next_run_at <= now (or manual).
	GetDueSchedules(ctx context.Context, now time.Time) ([]*Schedule, error)
	// CalculateNextRun computes the next occurrence from now.
	CalculateNextRun(s *Schedule, now time.Time) (*time.Time, error)
	// Count returns total non-deleted schedules.
	Count(ctx context.Context) (int64, error)
}
