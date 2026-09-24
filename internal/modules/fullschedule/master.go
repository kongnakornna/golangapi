package fullschedule

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// Errors for the master data (Group / Zone / Area) sub-module.
var (
	// ErrNoTargetDevices is returned when a schedule has neither a device
	// mapping nor a group/zone/area scope covering any device.
	ErrNoTargetDevices = errors.New("fullschedule: no target devices specified (set device_ids or a group/zone/area scope)")
	// ErrNoScope is returned when device resolution is asked without any scope.
	ErrNoScope = errors.New("fullschedule: no group/zone/area scope specified")
	// ErrZoneNotInGroup is returned when zone_id does not belong to group_id.
	ErrZoneNotInGroup = errors.New("fullschedule: zone_id does not belong to the given group_id")
	// ErrAreaNotInZone is returned when area_id does not belong to zone_id.
	ErrAreaNotInZone = errors.New("fullschedule: area_id does not belong to the given zone_id")
	// ErrGroupNotFound is returned when a group does not exist.
	ErrGroupNotFound = errors.New("fullschedule: group not found")
	// ErrZoneNotFound is returned when a zone does not exist.
	ErrZoneNotFound = errors.New("fullschedule: zone not found")
	// ErrAreaNotFound is returned when an area does not exist.
	ErrAreaNotFound = errors.New("fullschedule: area not found")
	// ErrGroupHasZones blocks deleting a group that still has child zones.
	ErrGroupHasZones = errors.New("fullschedule: group still has zones, cannot delete")
	// ErrZoneHasAreas blocks deleting a zone that still has child areas.
	ErrZoneHasAreas = errors.New("fullschedule: zone still has areas, cannot delete")
	// ErrAreaHasDevices blocks deleting an area that still has mapped devices.
	ErrAreaHasDevices = errors.New("fullschedule: area still has mapped devices, cannot delete")
	// ErrMasterStoreNotWired is returned when the master store is unavailable.
	ErrMasterStoreNotWired = errors.New("fullschedule: master data store not wired")
)

// SortOrder is an ordering direction for list endpoints.
type SortOrder string

const (
	// SortOrderASC sorts ascending (default).
	SortOrderASC SortOrder = "asc"
	// SortOrderDESC sorts descending.
	SortOrderDESC SortOrder = "desc"
)

// MasterPgRepository defines data access for Group / Zone / Area + device mapping.
// เก็บข้อมูล Master Data (fs_groups / fs_zones / fs_areas / fs_device_area)
type MasterPgRepository interface {
	// --- Group ---
	CreateGroup(ctx context.Context, g *Group) (*Group, error)
	GetGroup(ctx context.Context, id uuid.UUID) (*Group, error)
	GetGroups(ctx context.Context, limit, offset int) ([]*Group, error)
	UpdateGroup(ctx context.Context, g *Group, values map[string]interface{}) (*Group, error)
	CountGroups(ctx context.Context) (int64, error)
	// CountZonesByGroup counts zones under a group (delete guard).
	CountZonesByGroup(ctx context.Context, groupID uuid.UUID) (int64, error)

	// --- Zone ---
	CreateZone(ctx context.Context, z *Zone) (*Zone, error)
	GetZone(ctx context.Context, id uuid.UUID) (*Zone, error)
	GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*Zone, error)
	UpdateZone(ctx context.Context, z *Zone, values map[string]interface{}) (*Zone, error)
	// CountZones counts zones (optionally filtered by group).
	CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error)
	// CountAreasByZone counts areas under a zone (delete guard).
	CountAreasByZone(ctx context.Context, zoneID uuid.UUID) (int64, error)

	// --- Area ---
	CreateArea(ctx context.Context, a *Area) (*Area, error)
	GetArea(ctx context.Context, id uuid.UUID) (*Area, error)
	GetAreas(ctx context.Context, zoneID *uuid.UUID, sort SortOrder, limit, offset int) ([]*Area, error)
	UpdateArea(ctx context.Context, a *Area, values map[string]interface{}) (*Area, error)
	// CountAreas counts areas (optionally filtered by zone).
	CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error)
	// CountMappingsByArea counts mapped devices in an area (delete guard).
	CountMappingsByArea(ctx context.Context, areaID uuid.UUID) (int64, error)

	// --- Device <-> Area mapping ---
	// MapDevicesToArea replaces the device mapping of an area.
	MapDevicesToArea(ctx context.Context, areaID uuid.UUID, deviceIDs []int) error
	// GetAreaDevices returns the device mapping rows of an area.
	GetAreaDevices(ctx context.Context, areaID uuid.UUID) ([]*AreaDevice, error)
	// GetDeviceAreas returns the areas a device belongs to.
	GetDeviceAreas(ctx context.Context, deviceID int) ([]*AreaDevice, error)

	// --- Scope resolution (target devices) ---
	// GetDeviceIDsByScope resolves IoT device IDs from the group/zone/area scope.
	// Returns an empty slice when the scope covers no devices.
	GetDeviceIDsByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) ([]int, error)
	// CountDevicesByScope counts devices covered by the scope.
	CountDevicesByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) (int64, error)
}

// MasterUseCaseI defines business logic for Group / Zone / Area master data.
// ธุรกิจสำหรับ Master Data ของโมดูล Full Schedule
type MasterUseCaseI interface {
	// --- Group ---
	CreateGroup(ctx context.Context, g *Group) (*Group, error)
	GetGroup(ctx context.Context, id uuid.UUID) (*Group, error)
	GetGroups(ctx context.Context, limit, offset int) ([]*Group, error)
	UpdateGroup(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*Group, error)
	DeleteGroup(ctx context.Context, id uuid.UUID) (*Group, error)
	CountGroups(ctx context.Context) (int64, error)

	// --- Zone ---
	CreateZone(ctx context.Context, z *Zone) (*Zone, error)
	GetZone(ctx context.Context, id uuid.UUID) (*Zone, error)
	GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*Zone, error)
	UpdateZone(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*Zone, error)
	DeleteZone(ctx context.Context, id uuid.UUID) (*Zone, error)
	// CountZones counts zones (optionally filtered by group).
	CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error)

	// --- Area ---
	CreateArea(ctx context.Context, a *Area) (*Area, error)
	GetArea(ctx context.Context, id uuid.UUID) (*Area, error)
	GetAreas(ctx context.Context, zoneID *uuid.UUID, sort SortOrder, limit, offset int) ([]*Area, error)
	UpdateArea(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*Area, error)
	DeleteArea(ctx context.Context, id uuid.UUID) (*Area, error)
	// CountAreas counts areas (optionally filtered by zone).
	CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error)

	// --- Device <-> Area mapping ---
	MapDevicesToArea(ctx context.Context, areaID uuid.UUID, deviceIDs []int) error
	GetAreaDevices(ctx context.Context, areaID uuid.UUID) ([]*AreaDevice, error)

	// ValidateScope verifies the zone/area belong to their parents. Returns
	// ErrNoScope when no scope is set, which is legitimate.
	ValidateScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) error
	// GetDeviceIDsByScope resolves target devices from the scope.
	GetDeviceIDsByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) ([]int, error)
	// CountDevicesByScope counts target devices for a scope preview.
	CountDevicesByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) (int64, error)
}

// MasterHandlers defines HTTP handlers for Group / Zone / Area master data.
// ตัวจัดการ HTTP สำหรับ Master Data
type MasterHandlers interface {
	// --- Group ---
	CreateGroup() func(w http.ResponseWriter, r *http.Request)
	ListGroups() func(w http.ResponseWriter, r *http.Request)
	GetGroup() func(w http.ResponseWriter, r *http.Request)
	UpdateGroup() func(w http.ResponseWriter, r *http.Request)
	DeleteGroup() func(w http.ResponseWriter, r *http.Request)

	// --- Zone ---
	CreateZone() func(w http.ResponseWriter, r *http.Request)
	ListZones() func(w http.ResponseWriter, r *http.Request)
	GetZone() func(w http.ResponseWriter, r *http.Request)
	UpdateZone() func(w http.ResponseWriter, r *http.Request)
	DeleteZone() func(w http.ResponseWriter, r *http.Request)

	// --- Area ---
	CreateArea() func(w http.ResponseWriter, r *http.Request)
	ListAreas() func(w http.ResponseWriter, r *http.Request)
	GetArea() func(w http.ResponseWriter, r *http.Request)
	UpdateArea() func(w http.ResponseWriter, r *http.Request)
	DeleteArea() func(w http.ResponseWriter, r *http.Request)

	// MapAreaDevices binds IoT devices to an area.
	MapAreaDevices() func(w http.ResponseWriter, r *http.Request)
	// AreaDevices lists devices mapped to an area.
	AreaDevices() func(w http.ResponseWriter, r *http.Request)
	// ScopePreview counts devices covered by a group/zone/area selection.
	ScopePreview() func(w http.ResponseWriter, r *http.Request)
}
