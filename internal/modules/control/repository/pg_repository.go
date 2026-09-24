package repository

import (
	"context"

	iotmodels "icmongolang/internal/modules/iot/models"
	"gorm.io/gorm"
)

// ScheduledDevice joins a schedule with the devices it controls.
type ScheduledDevice struct {
	Device iotmodels.Device
	Event  int
}

// Repository provides device + schedule access for control execution.
type Repository interface {
	GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error)
	GetActiveSchedules(ctx context.Context) ([]iotmodels.Schedule, error)
	GetScheduleByID(ctx context.Context, scheduleID int) (*iotmodels.Schedule, error)
	GetScheduleDevices(ctx context.Context, scheduleID int) ([]int, error)
	CreateScheduleLog(ctx context.Context, log *iotmodels.ScheduleProcessLog) error
}

type pgRepository struct {
	db *gorm.DB
}

// NewPgRepository creates a DB-backed control repository.
func NewPgRepository(db *gorm.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error) {
	var d iotmodels.Device
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *pgRepository) GetActiveSchedules(ctx context.Context) ([]iotmodels.Schedule, error) {
	var list []iotmodels.Schedule
	if err := r.db.WithContext(ctx).Where("status = 1").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgRepository) GetScheduleByID(ctx context.Context, scheduleID int) (*iotmodels.Schedule, error) {
	var s iotmodels.Schedule
	if err := r.db.WithContext(ctx).Where("schedule_id = ?", scheduleID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *pgRepository) GetScheduleDevices(ctx context.Context, scheduleID int) ([]int, error) {
	var rows []iotmodels.ScheduleDevice
	if err := r.db.WithContext(ctx).Where("schedule_id = ?", scheduleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]int, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.DeviceID)
	}
	return out, nil
}

func (r *pgRepository) CreateScheduleLog(ctx context.Context, log *iotmodels.ScheduleProcessLog) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(log).Error
}
