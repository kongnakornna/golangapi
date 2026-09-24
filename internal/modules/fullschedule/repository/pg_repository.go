package repository

import (
	"context"
	"time"

	"icmongolang/internal/modules/fullschedule"
	iotmodels "icmongolang/internal/modules/iot/models"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FullSchedulePgRepo implements fullschedule.FullSchedulePgRepository.
// รีโพสิทอรีสำหรับโมดูล Full Schedule (ตาราง fs_*)
type FullSchedulePgRepo struct {
	repository.PgRepo[fullschedule.Schedule]
	DB *gorm.DB
}

// CreateFullSchedulePgRepository creates a new repository instance.
// สร้างรีโพสิทอรีสำหรับ Full Schedule
func CreateFullSchedulePgRepository(db *gorm.DB) fullschedule.FullSchedulePgRepository {
	return &FullSchedulePgRepo{
		PgRepo: repository.CreatePgRepo[fullschedule.Schedule](db),
		DB:     db,
	}
}

// Get returns a non-deleted schedule by ID.
// ดึง schedule (ไม่รวมที่ถูกลบแบบ soft delete)
func (r *FullSchedulePgRepo) Get(ctx context.Context, id uuid.UUID) (*fullschedule.Schedule, error) {
	var s fullschedule.Schedule
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id.String()).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// GetMulti lists non-deleted schedules (newest first).
// รายการ schedule ที่ไม่ถูกลบ เรียงตามสร้างล่าสุด
func (r *FullSchedulePgRepo) GetMulti(ctx context.Context, limit, offset int) ([]*fullschedule.Schedule, error) {
	if limit <= 0 {
		limit = 50
	}
	var list []*fullschedule.Schedule
	if err := r.DB.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetDevices returns the mapped IoT devices of a schedule.
// ดึงอุปกรณ์ IoT ที่แมปกับ schedule
func (r *FullSchedulePgRepo) GetDevices(ctx context.Context, scheduleID uuid.UUID) ([]*fullschedule.ScheduleDevice, error) {
	var rows []*fullschedule.ScheduleDevice
	if err := r.DB.WithContext(ctx).Where("schedule_id = ?", scheduleID.String()).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ListDeviceIDs returns the raw device IDs mapped to a schedule.
// ดึงรหัสอุปกรณ์ที่แมปกับ schedule
func (r *FullSchedulePgRepo) ListDeviceIDs(ctx context.Context, scheduleID uuid.UUID) ([]int, error) {
	var rows []*fullschedule.ScheduleDevice
	rows, err := r.GetDevices(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	out := make([]int, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.DeviceID)
	}
	return out, nil
}

// ReplaceDevices wipes and re-inserts the device mapping.
// ลบและแมปอุปกรณ์ใหม่ทั้งหมด
func (r *FullSchedulePgRepo) ReplaceDevices(ctx context.Context, scheduleID uuid.UUID, deviceIDs []int) error {
	if err := r.DB.WithContext(ctx).
		Where("schedule_id = ?", scheduleID.String()).
		Delete(&fullschedule.ScheduleDevice{}).Error; err != nil {
		return err
	}
	for _, id := range deviceIDs {
		row := &fullschedule.ScheduleDevice{ScheduleID: scheduleID, DeviceID: id}
		if err := r.DB.WithContext(ctx).Create(row).Error; err != nil {
			return err
		}
	}
	return nil
}

// ListPage lists non-deleted schedules with filters and pagination.
// รายการ schedule พร้อมตัวกรองและแบ่งหน้า (แบบ listschedulepage)
func (r *FullSchedulePgRepo) ListPage(ctx context.Context, f fullschedule.ScheduleFilter, page, pageSize int) ([]*fullschedule.Schedule, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	const maxPageSize = 100
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	q := r.listFilter(ctx, f)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if o, ok := fullschedule.ParseScheduleSort(f.Sort); ok {
		order = o
	}
	var list []*fullschedule.Schedule
	if err := q.Order(order).Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	if err := r.attachDeviceCounts(ctx, list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListAll lists non-deleted schedules with filters (no pagination).
// รายการ schedule ทั้งหมด (ไม่มีแบ่งหน้า)
func (r *FullSchedulePgRepo) ListAll(ctx context.Context, f fullschedule.ScheduleFilter) ([]*fullschedule.Schedule, error) {
	order := "created_at DESC"
	if o, ok := fullschedule.ParseScheduleSort(f.Sort); ok {
		order = o
	}
	var list []*fullschedule.Schedule
	if err := r.listFilter(ctx, f).Order(order).Find(&list).Error; err != nil {
		return nil, err
	}
	if err := r.attachDeviceCounts(ctx, list); err != nil {
		return nil, err
	}
	return list, nil
}

// listFilter builds the WHERE clause for schedule listing.
func (r *FullSchedulePgRepo) listFilter(ctx context.Context, f fullschedule.ScheduleFilter) *gorm.DB {
	q := r.DB.WithContext(ctx).Model(&fullschedule.Schedule{}).Where("deleted_at IS NULL")
	if f.Keyword != "" {
		q = q.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	if f.Mode != nil {
		q = q.Where("mode = ?", *f.Mode)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.EventType != nil {
		q = q.Where("event_type = ?", *f.EventType)
	}
	if f.EventAction != nil {
		q = q.Where("event_action = ?", *f.EventAction)
	}
	if f.Start != "" {
		q = q.Where("time_start = ?", f.Start)
	}
	if f.Scope.GroupID != nil {
		q = q.Where("group_id = ?", f.Scope.GroupID.String())
	}
	if f.Scope.ZoneID != nil {
		q = q.Where("zone_id = ?", f.Scope.ZoneID.String())
	}
	if f.Scope.AreaID != nil {
		q = q.Where("area_id = ?", f.Scope.AreaID.String())
	}
	return q
}

// attachDeviceCounts fills DeviceCount for each schedule in one grouped query.
// ใส่จำนวน device ที่แมปให้แต่ละ schedule (query เดียว)
func (r *FullSchedulePgRepo) attachDeviceCounts(ctx context.Context, list []*fullschedule.Schedule) error {
	if len(list) == 0 {
		return nil
	}
	ids := make([]string, 0, len(list))
	for _, s := range list {
		ids = append(ids, s.ID.String())
	}
	var counts []struct {
		ScheduleID string
		Cnt        int64
	}
	if err := r.DB.WithContext(ctx).
		Table("fs_schedule_device").
		Select("schedule_id, COUNT(*) AS cnt").
		Where("schedule_id IN ?", ids).
		Group("schedule_id").
		Scan(&counts).Error; err != nil {
		return err
	}
	m := make(map[string]int64, len(counts))
	for _, c := range counts {
		m[c.ScheduleID] = c.Cnt
	}
	for _, s := range list {
		s.DeviceCount = m[s.ID.String()]
	}
	return nil
}

// ListScheduleDevicePage lists schedule-device mappings with pagination.
// รายการแมป schedule ↔ device แบบแบ่งหน้า (แบบ scheduledevicepage)
func (r *FullSchedulePgRepo) ListScheduleDevicePage(ctx context.Context, f fullschedule.ScheduleDeviceFilter, page, pageSize int) ([]*fullschedule.ScheduleDeviceRow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	const maxPageSize = 100
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	base := func(db *gorm.DB) *gorm.DB {
		db = db.Table("fs_schedule_device sd").
			Joins("JOIN fs_schedule s ON s.id = sd.schedule_id AND s.deleted_at IS NULL").
			Joins("LEFT JOIN sd_iot_device d ON d.device_id = sd.device_id")
		if f.Keyword != "" {
			db = db.Where("(s.name ILIKE ? OR d.device_name ILIKE ?)", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
		}
		if f.ScheduleID != nil {
			db = db.Where("sd.schedule_id = ?", f.ScheduleID.String())
		}
		if f.DeviceID != nil {
			db = db.Where("sd.device_id = ?", *f.DeviceID)
		}
		if f.Status != nil {
			db = db.Where("s.status = ?", *f.Status)
		}
		return db
	}

	var total int64
	if err := base(r.DB).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "s.time_start ASC, sd.device_id ASC"
	if o, ok := fullschedule.ParseScheduleDeviceSort(f.Sort); ok {
		order = o
	}
	var rows []*fullschedule.ScheduleDeviceRow
	if err := base(r.DB).
		Select(`sd.schedule_id AS schedule_id, s.name AS schedule_name, s.mode AS mode,
		        s.status AS status,
		        s.time_start AS start, s.event_action AS event_action,
		        sd.device_id AS device_id, d.device_name AS device_name, d.sn AS sn`).
		Order(order).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GetDueSchedules returns active schedules that are due or manually flagged.
// ดึง schedule ที่ถึงกำหนดทำงาน (อัตโนมัติ/สั่งด้วยมือ)
func (r *FullSchedulePgRepo) GetDueSchedules(ctx context.Context, now time.Time) ([]*fullschedule.Schedule, error) {
	var list []*fullschedule.Schedule
	err := r.DB.WithContext(ctx).
		Preload("Devices").
		Preload("Settings").
		Where("deleted_at IS NULL").
		Where("status = ?", fullschedule.StatusValueActive).
		Where("(manual_trigger = true) OR (next_run_at IS NOT NULL AND next_run_at <= ?)", now).
		Order("next_run_at ASC NULLS FIRST").
		Limit(500).
		Find(&list).Error
	return list, err
}

// GetDeviceByID resolves an IoT device by internal ID (read-only).
// ดึงข้อมูลอุปกรณ์ IoT (อ่านอย่างเดียว)
func (r *FullSchedulePgRepo) GetDeviceByID(ctx context.Context, deviceID int) (*iotmodels.Device, error) {
	var d iotmodels.Device
	if err := r.DB.WithContext(ctx).Where("device_id = ?", deviceID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

// ListDevices finds IoT devices for the mapping picker (read-only).
// รายการอุปกรณ์ IoT สำหรับเลือกแมป
func (r *FullSchedulePgRepo) ListDevices(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*iotmodels.Device, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	const maxPageSize = 100
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	q := r.DB.WithContext(ctx).Model(&iotmodels.Device{})
	if keyword, ok := filter["keyword"].(string); ok && keyword != "" {
		q = q.Where("device_name ILIKE ? OR sn ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status, ok := filter["status"].(int); ok && status != 0 {
		q = q.Where("status = ?", status)
	}
	if locationID, ok := filter["location_id"].(int); ok && locationID != 0 {
		q = q.Where("location_id = ?", locationID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []*iotmodels.Device
	if err := q.Order("device_id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// CreateHistory inserts a history row.
// สร้างประวัติการทำงาน
func (r *FullSchedulePgRepo) CreateHistory(ctx context.Context, h *fullschedule.ScheduleHistory) (*fullschedule.ScheduleHistory, error) {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	if err := r.DB.WithContext(ctx).Create(h).Error; err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHistory patches a history row.
// อัปเดตประวัติการทำงาน
func (r *FullSchedulePgRepo) UpdateHistory(ctx context.Context, h *fullschedule.ScheduleHistory, values map[string]interface{}) error {
	if values == nil {
		values = map[string]interface{}{}
	}
	values["updated_at"] = time.Now()
	return r.DB.WithContext(ctx).Model(h).Where("id = ?", h.ID.String()).Updates(values).Error
}

// GetHistoryBySchedule lists history of a single schedule (newest first).
// ดูประวัติของ schedule หนึ่งรายการ
func (r *FullSchedulePgRepo) GetHistoryBySchedule(ctx context.Context, scheduleID uuid.UUID, limit, offset int) ([]*fullschedule.ScheduleHistory, error) {
	var rows []*fullschedule.ScheduleHistory
	q := r.DB.WithContext(ctx).Where("schedule_id = ?", scheduleID.String()).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetHistory lists history across all schedules.
// ดูประวัติทั้งหมด
func (r *FullSchedulePgRepo) GetHistory(ctx context.Context, limit, offset int, status *fullschedule.HistoryStatus) ([]*fullschedule.ScheduleHistory, error) {
	var rows []*fullschedule.ScheduleHistory
	q := r.DB.WithContext(ctx).Order("created_at DESC")
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// UpsertSetting inserts or updates a per-schedule setting.
// เพิ่ม/อัปเดตการตั้งค่าของ schedule
func (r *FullSchedulePgRepo) UpsertSetting(ctx context.Context, s *fullschedule.ScheduleSetting) (*fullschedule.ScheduleSetting, error) {
	var existing fullschedule.ScheduleSetting
	err := r.DB.WithContext(ctx).
		Where("schedule_id = ? AND key = ?", s.ScheduleID.String(), s.Key).
		First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}
		if err := r.DB.WithContext(ctx).Create(s).Error; err != nil {
			return nil, err
		}
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	existing.Value = s.Value
	if err := r.DB.WithContext(ctx).Model(&existing).
		Update("value", s.Value).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// GetSettings returns all settings of a schedule.
// ดึงการตั้งค่าทั้งหมดของ schedule
func (r *FullSchedulePgRepo) GetSettings(ctx context.Context, scheduleID uuid.UUID) ([]*fullschedule.ScheduleSetting, error) {
	var rows []*fullschedule.ScheduleSetting
	if err := r.DB.WithContext(ctx).
		Where("schedule_id = ?", scheduleID.String()).
		Order("key ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetReport aggregates history rows for the report endpoint.
// สรุปรายงาน success/failed จากประวัติ
func (r *FullSchedulePgRepo) GetReport(ctx context.Context, from, to *time.Time, scheduleID *uuid.UUID, scope fullschedule.ScopeFilter) ([]*fullschedule.ScheduleReport, error) {
	q := r.DB.WithContext(ctx).Table("fs_schedule_history h").
		Select(`h.schedule_id AS schedule_id,
		        COALESCE(s.name, '') AS name,
		        COALESCE(s.mode, '') AS mode,
		        COUNT(h.id) AS total_runs,
		        COUNT(*) FILTER (WHERE h.status = 'success') AS success_count,
		        COUNT(*) FILTER (WHERE h.status = 'failed')  AS failed_count,
		        COUNT(*) FILTER (WHERE h.status = 'skipped') AS skipped_count,
		        COUNT(*) FILTER (WHERE h.status = 'processing') AS processing,
		        ROUND(100.0 * COUNT(*) FILTER (WHERE h.status = 'success') / NULLIF(COUNT(h.id) FILTER (WHERE h.status IN ('success','failed')), 0), 2) AS success_rate,
		        COALESCE(AVG(h.duration_ms), 0) AS avg_duration,
		        MAX(h.executed_at) AS last_run_at`).
		Joins("LEFT JOIN fs_schedule s ON s.id = h.schedule_id").
		Where("h.created_at >= ?", from).
		Where("h.created_at <= ?", to).
		Group("h.schedule_id, s.name, s.mode").
		Order("h.schedule_id ASC")
	if scheduleID != nil {
		q = q.Where("h.schedule_id = ?", scheduleID.String())
	}
	if scope.GroupID != nil {
		q = q.Where("s.group_id = ?", scope.GroupID.String())
	}
	if scope.ZoneID != nil {
		q = q.Where("s.zone_id = ?", scope.ZoneID.String())
	}
	if scope.AreaID != nil {
		q = q.Where("s.area_id = ?", scope.AreaID.String())
	}
	var rows []*fullschedule.ScheduleReport
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// Count returns the number of non-deleted schedules.
// นับจำนวน schedule ทั้งหมด
func (r *FullSchedulePgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&fullschedule.Schedule{}).
		Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
