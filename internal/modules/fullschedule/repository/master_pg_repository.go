package repository

import (
	"context"
	"time"

	"icmongolang/internal/modules/fullschedule"
	iotmodels "icmongolang/internal/modules/iot/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MasterPgRepo implements fullschedule.MasterPgRepository (Group / Zone / Area).
// รีโพสิทอรีสำหรับ Master Data ของโมดูล Full Schedule (fs_groups / fs_zones / fs_areas / fs_device_area)
type MasterPgRepo struct {
	DB *gorm.DB
}

// CreateMasterPgRepository creates a new master data repository.
// สร้างรีโพสิทอรีสำหรับ Master Data
func CreateMasterPgRepository(db *gorm.DB) fullschedule.MasterPgRepository {
	return &MasterPgRepo{DB: db}
}

// maxSortID returns the current maximum sort_id (0 when empty) for rows matching where.
// ใช้หาค่า sort_id ลำดับถัดไปสำหรับ auto-assign ตอนสร้างรายการใหม่
func (r *MasterPgRepo) maxSortID(ctx context.Context, m interface{}, where string, args ...interface{}) (int, error) {
	var max int
	q := r.DB.WithContext(ctx).Model(m)
	if where != "" {
		q = q.Where(where, args...)
	}
	if err := q.Select("COALESCE(MAX(sort_id), 0)").Scan(&max).Error; err != nil {
		return 0, err
	}
	return max + 1, nil
}

// ============================================================================
// Group
// ============================================================================

func (r *MasterPgRepo) CreateGroup(ctx context.Context, g *fullschedule.Group) (*fullschedule.Group, error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	if g.Status == "" {
		g.Status = fullschedule.StatusActive
	}
	if g.SortID == 0 {
		next, err := r.maxSortID(ctx, &fullschedule.Group{}, "deleted_at IS NULL")
		if err != nil {
			return nil, err
		}
		g.SortID = next
	}
	if err := r.DB.WithContext(ctx).Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func (r *MasterPgRepo) GetGroup(ctx context.Context, id uuid.UUID) (*fullschedule.Group, error) {
	var g fullschedule.Group
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id.String()).First(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *MasterPgRepo) GetGroups(ctx context.Context, limit, offset int) ([]*fullschedule.Group, error) {
	if limit <= 0 {
		limit = 50
	}
	var rows []*fullschedule.Group
	if err := r.DB.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("sort_id ASC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *MasterPgRepo) UpdateGroup(ctx context.Context, g *fullschedule.Group, values map[string]interface{}) (*fullschedule.Group, error) {
	values["updated_at"] = time.Now()
	if err := r.DB.WithContext(ctx).Model(g).
		Where("id = ? AND deleted_at IS NULL", g.ID.String()).
		Updates(values).Error; err != nil {
		return nil, err
	}
	return r.GetGroup(ctx, g.ID)
}

func (r *MasterPgRepo) CountGroups(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&fullschedule.Group{}).
		Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountZonesByGroup counts zones under a group (delete guard).
func (r *MasterPgRepo) CountZonesByGroup(ctx context.Context, groupID uuid.UUID) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&fullschedule.Zone{}).
		Where("group_id = ? AND deleted_at IS NULL", groupID.String()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================================
// Zone
// ============================================================================

func (r *MasterPgRepo) CreateZone(ctx context.Context, z *fullschedule.Zone) (*fullschedule.Zone, error) {
	if z.ID == uuid.Nil {
		z.ID = uuid.New()
	}
	if z.SortID == 0 {
		next, err := r.maxSortID(ctx, &fullschedule.Zone{}, "group_id = ? AND deleted_at IS NULL", z.GroupID.String())
		if err != nil {
			return nil, err
		}
		z.SortID = next
	}
	if err := r.DB.WithContext(ctx).Create(z).Error; err != nil {
		return nil, err
	}
	return z, nil
}

func (r *MasterPgRepo) GetZone(ctx context.Context, id uuid.UUID) (*fullschedule.Zone, error) {
	var z fullschedule.Zone
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id.String()).First(&z).Error; err != nil {
		return nil, err
	}
	return &z, nil
}

func (r *MasterPgRepo) GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*fullschedule.Zone, error) {
	if limit <= 0 {
		limit = 50
	}
	q := r.DB.WithContext(ctx).
		Model(&fullschedule.Zone{}).
		Select("fs_zones.*, g.name AS group_name").
		Joins("JOIN fs_groups g ON g.id = fs_zones.group_id AND g.deleted_at IS NULL").
		Where("fs_zones.deleted_at IS NULL").
		Order("fs_zones.sort_id ASC, fs_zones.created_at DESC")
	if groupID != nil {
		q = q.Where("fs_zones.group_id = ?", groupID.String())
	}
	var rows []*fullschedule.Zone
	if err := q.Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *MasterPgRepo) UpdateZone(ctx context.Context, z *fullschedule.Zone, values map[string]interface{}) (*fullschedule.Zone, error) {
	values["updated_at"] = time.Now()
	if err := r.DB.WithContext(ctx).Model(z).
		Where("id = ? AND deleted_at IS NULL", z.ID.String()).
		Updates(values).Error; err != nil {
		return nil, err
	}
	return r.GetZone(ctx, z.ID)
}

func (r *MasterPgRepo) CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error) {
	var count int64
	q := r.DB.WithContext(ctx).
		Model(&fullschedule.Zone{}).
		Joins("JOIN fs_groups g ON g.id = fs_zones.group_id AND g.deleted_at IS NULL").
		Where("fs_zones.deleted_at IS NULL")
	if groupID != nil {
		q = q.Where("fs_zones.group_id = ?", groupID.String())
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountAreasByZone counts areas under a zone (delete guard).
func (r *MasterPgRepo) CountAreasByZone(ctx context.Context, zoneID uuid.UUID) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&fullschedule.Area{}).
		Where("zone_id = ? AND deleted_at IS NULL", zoneID.String()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================================
// Area
// ============================================================================

func (r *MasterPgRepo) CreateArea(ctx context.Context, a *fullschedule.Area) (*fullschedule.Area, error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.SortID == 0 {
		next, err := r.maxSortID(ctx, &fullschedule.Area{}, "zone_id = ? AND deleted_at IS NULL", a.ZoneID.String())
		if err != nil {
			return nil, err
		}
		a.SortID = next
	}
	if err := r.DB.WithContext(ctx).Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *MasterPgRepo) GetArea(ctx context.Context, id uuid.UUID) (*fullschedule.Area, error) {
	var a fullschedule.Area
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id.String()).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *MasterPgRepo) GetAreas(ctx context.Context, zoneID *uuid.UUID, sort fullschedule.SortOrder, limit, offset int) ([]*fullschedule.Area, error) {
	if limit <= 0 {
		limit = 50
	}
	order := "fs_areas.sort_id ASC, fs_areas.created_at DESC"
	if sort == fullschedule.SortOrderDESC {
		order = "fs_areas.sort_id DESC, fs_areas.created_at DESC"
	}
	q := r.DB.WithContext(ctx).
		Model(&fullschedule.Area{}).
		Select("fs_areas.*, z.name AS zone_name, g.name AS group_name").
		Joins("JOIN fs_zones z ON z.id = fs_areas.zone_id AND z.deleted_at IS NULL").
		Joins("JOIN fs_groups g ON g.id = z.group_id AND g.deleted_at IS NULL").
		Where("fs_areas.deleted_at IS NULL").
		Order(order)
	if zoneID != nil {
		q = q.Where("fs_areas.zone_id = ?", zoneID.String())
	}
	var rows []*fullschedule.Area
	if err := q.Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *MasterPgRepo) UpdateArea(ctx context.Context, a *fullschedule.Area, values map[string]interface{}) (*fullschedule.Area, error) {
	values["updated_at"] = time.Now()
	if err := r.DB.WithContext(ctx).Model(a).
		Where("id = ? AND deleted_at IS NULL", a.ID.String()).
		Updates(values).Error; err != nil {
		return nil, err
	}
	return r.GetArea(ctx, a.ID)
}

func (r *MasterPgRepo) CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error) {
	var count int64
	q := r.DB.WithContext(ctx).
		Model(&fullschedule.Area{}).
		Joins("JOIN fs_zones z ON z.id = fs_areas.zone_id AND z.deleted_at IS NULL").
		Joins("JOIN fs_groups g ON g.id = z.group_id AND g.deleted_at IS NULL").
		Where("fs_areas.deleted_at IS NULL")
	if zoneID != nil {
		q = q.Where("fs_areas.zone_id = ?", zoneID.String())
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountMappingsByArea counts mapped devices in an area (delete guard).
func (r *MasterPgRepo) CountMappingsByArea(ctx context.Context, areaID uuid.UUID) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&fullschedule.AreaDevice{}).
		Where("area_id = ?", areaID.String()).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================================
// Device <-> Area mapping
// ============================================================================

// MapDevicesToArea replaces the device mapping of an area.
func (r *MasterPgRepo) MapDevicesToArea(ctx context.Context, areaID uuid.UUID, deviceIDs []int) error {
	if err := r.DB.WithContext(ctx).
		Where("area_id = ?", areaID.String()).
		Delete(&fullschedule.AreaDevice{}).Error; err != nil {
		return err
	}
	for _, id := range deviceIDs {
		if id <= 0 {
			continue
		}
		row := &fullschedule.AreaDevice{AreaID: areaID, DeviceID: id, DeviceSN: r.resolveSN(ctx, id)}
		if err := r.DB.WithContext(ctx).Create(row).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAreaDevices returns the device mapping rows of an area.
func (r *MasterPgRepo) GetAreaDevices(ctx context.Context, areaID uuid.UUID) ([]*fullschedule.AreaDevice, error) {
	var rows []*fullschedule.AreaDevice
	if err := r.DB.WithContext(ctx).
		Where("area_id = ?", areaID.String()).
		Order("device_id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetDeviceAreas returns the areas a device belongs to.
func (r *MasterPgRepo) GetDeviceAreas(ctx context.Context, deviceID int) ([]*fullschedule.AreaDevice, error) {
	var rows []*fullschedule.AreaDevice
	if err := r.DB.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ============================================================================
// Scope resolution (target devices)
// ============================================================================

// GetDeviceIDsByScope resolves device IDs from a group/zone/area scope.
// Priority: area -> zone -> group (area covers the narrowest set).
func (r *MasterPgRepo) GetDeviceIDsByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) ([]int, error) {
	var ids []int
	if err := r.scopeQuery(ctx, groupID, zoneID, areaID).
		Distinct("fs_device_area.device_id").
		Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// CountDevicesByScope counts devices covered by the scope (preview).
func (r *MasterPgRepo) CountDevicesByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) (int64, error) {
	var count int64
	if err := r.scopeQuery(ctx, groupID, zoneID, areaID).
		Distinct("fs_device_area.device_id").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// scopeQuery builds the target-device query. Precondition: at least one scope id.
func (r *MasterPgRepo) scopeQuery(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) *gorm.DB {
	q := r.DB.WithContext(ctx).Model(&fullschedule.AreaDevice{}).Select("fs_device_area.device_id")
	switch {
	case areaID != nil:
		q = q.Where("fs_device_area.area_id = ?", areaID.String())
	case zoneID != nil:
		q = q.Joins("JOIN fs_areas a ON a.id = fs_device_area.area_id").
			Where("a.zone_id = ?", zoneID.String())
	case groupID != nil:
		q = q.Joins("JOIN fs_areas a ON a.id = fs_device_area.area_id").
			Joins("JOIN fs_zones z ON z.id = a.zone_id").
			Where("z.group_id = ?", groupID.String())
	default:
		// No scope: match nothing.
		q = q.Where("fs_device_area.area_id = NULL")
	}
	return q
}

func (r *MasterPgRepo) resolveSN(ctx context.Context, deviceID int) *string {
	var d iotmodels.Device
	if err := r.DB.WithContext(ctx).Select("sn").
		Where("device_id = ?", deviceID).First(&d).Error; err != nil || d.SN == "" {
		return nil
	}
	return &d.SN
}

// compile-time assertion
var _ fullschedule.MasterPgRepository = (*MasterPgRepo)(nil)


