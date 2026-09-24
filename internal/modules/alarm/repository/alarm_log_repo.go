package repository

import (
	"icmongolang/internal/modules/alarm/models"
	"gorm.io/gorm"
)

var alarmLogAllowedFilters = map[string]bool{
	"alarm_action_id": true,
	"device_id":       true,
	"type_id":         true,
	"event":           true,
	"status":          true,
	"alarm_type":      true,
}

type AlarmLogRepository interface {
	Create(log *models.AlarmLog) error
	GetByID(id uint) (*models.AlarmLog, error)
	List(filter map[string]interface{}, limit, offset int) ([]models.AlarmLog, int64, error)
}

type alarmLogRepo struct {
	db *gorm.DB
}

func NewAlarmLogRepository(db *gorm.DB) AlarmLogRepository {
	return &alarmLogRepo{db: db}
}

func (r *alarmLogRepo) Create(log *models.AlarmLog) error {
	return r.db.Create(log).Error
}

func (r *alarmLogRepo) GetByID(id uint) (*models.AlarmLog, error) {
	var log models.AlarmLog
	err := r.db.First(&log, id).Error
	return &log, err
}

func (r *alarmLogRepo) List(filter map[string]interface{}, limit, offset int) ([]models.AlarmLog, int64, error) {
	var logs []models.AlarmLog
	var total int64
	query := r.db.Model(&models.AlarmLog{})
	for k, v := range filter {
		if !alarmLogAllowedFilters[k] {
			continue
		}
		query = query.Where(k+" = ?", v)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}