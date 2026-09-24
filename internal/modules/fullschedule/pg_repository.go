package fullschedule

import (
	"context"
	"time"

	iotmodels "icmongolang/internal/modules/iot/models"

	"github.com/google/uuid"
)

// FullSchedulePgRepository defines data access for the Full Schedule module.
// ดึงข้อมูล Full Schedule จากฐานข้อมูล (ตาราง fs_*)
type FullSchedulePgRepository interface {
	Create(ctx context.Context, exp *Schedule) (*Schedule, error)
	Get(ctx context.Context, id uuid.UUID) (*Schedule, error)
	GetMulti(ctx context.Context, limit, offset int) ([]*Schedule, error)
	Delete(ctx context.Context, id uuid.UUID) (*Schedule, error)
	Update(ctx context.Context, exp *Schedule, values map[string]interface{}) (*Schedule, error)

	// ListPage lists non-deleted schedules with filters + pagination
	// (settings-style listschedulepage display).
	ListPage(ctx context.Context, f ScheduleFilter, page, pageSize int) ([]*Schedule, int64, error)
	// ListAll lists non-deleted schedules with filters (no pagination).
	ListAll(ctx context.Context, f ScheduleFilter) ([]*Schedule, error)
	// ListScheduleDevicePage lists schedule-device mappings with pagination.
	ListScheduleDevicePage(ctx context.Context, f ScheduleDeviceFilter, page, pageSize int) ([]*ScheduleDeviceRow, int64, error)

	// GetDueSchedules returns active schedules that are due (next_run_at <= now)
	// or flagged with manual_trigger = true.
	GetDueSchedules(ctx context.Context, now time.Time) ([]*Schedule, error)

	// ReplaceDevices replaces the device mapping of a schedule.
	ReplaceDevices(ctx context.Context, scheduleID uuid.UUID, deviceIDs []int) error
	// GetDevices returns the mapped devices of a schedule.
	GetDevices(ctx context.Context, scheduleID uuid.UUID) ([]*ScheduleDevice, error)
	// ListDeviceIDs returns only device IDs mapped to a schedule.
	ListDeviceIDs(ctx context.Context, scheduleID uuid.UUID) ([]int, error)

	// GetDeviceByID resolves an IoT device by its internal ID (read-only).
	GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error)
	// ListDevices lists IoT devices for the device picker (read-only).
	ListDevices(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*iotmodels.Device, int64, error)

	// CreateHistory inserts an execution history row.
	CreateHistory(ctx context.Context, h *ScheduleHistory) (*ScheduleHistory, error)
	// UpdateHistory updates an existing history row (final status etc.).
	UpdateHistory(ctx context.Context, h *ScheduleHistory, values map[string]interface{}) error
	// GetHistoryBySchedule lists history of one schedule.
	GetHistoryBySchedule(ctx context.Context, scheduleID uuid.UUID, limit, offset int) ([]*ScheduleHistory, error)
	// GetHistory lists history across all schedules.
	GetHistory(ctx context.Context, limit, offset int, status *HistoryStatus) ([]*ScheduleHistory, error)

	// UpsertSetting inserts or updates a per-schedule setting.
	UpsertSetting(ctx context.Context, s *ScheduleSetting) (*ScheduleSetting, error)
	// GetSettings returns settings of a schedule.
	GetSettings(ctx context.Context, scheduleID uuid.UUID) ([]*ScheduleSetting, error)

	// GetReport aggregates history rows for reporting.
	GetReport(ctx context.Context, from, to *time.Time, scheduleID *uuid.UUID, scope ScopeFilter) ([]*ScheduleReport, error)

	// Count returns the number of non-deleted schedules.
	Count(ctx context.Context) (int64, error)
}
