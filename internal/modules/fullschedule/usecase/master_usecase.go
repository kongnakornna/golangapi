package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// Errors for master data operations.
var (
	ErrNameRequired    = errors.New("fullschedule: name is required")
	ErrScopeIncomplete = errors.New("fullschedule: zone_id requires the zone to exist; area_id requires the area to exist")
)

// masterUseCase implements fullschedule.MasterUseCaseI.
type masterUseCase struct {
	repo   fullschedule.MasterPgRepository
	cfg    *config.Config
	logger logger.Logger
}

// CreateMasterUseCaseI creates the master data usecase.
// สร้างยูสเคสสำหรับ Master Data (Group/Zone/Area)
func CreateMasterUseCaseI(repo fullschedule.MasterPgRepository, cfg *config.Config, logger logger.Logger) fullschedule.MasterUseCaseI {
	return &masterUseCase{repo: repo, cfg: cfg, logger: logger}
}

// ============================================================================
// Group
// ============================================================================

func (u *masterUseCase) CreateGroup(ctx context.Context, g *fullschedule.Group) (*fullschedule.Group, error) {
	if err := validateMasterName(g.Name); err != nil {
		return nil, err
	}
	if g.Status == "" {
		g.Status = fullschedule.StatusActive
	}
	if g.Status != fullschedule.StatusActive && g.Status != fullschedule.StatusInactive {
		return nil, ErrStatusInvalid
	}
	return u.repo.CreateGroup(ctx, g)
}

func (u *masterUseCase) GetGroup(ctx context.Context, id uuid.UUID) (*fullschedule.Group, error) {
	g, err := u.repo.GetGroup(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrGroupNotFound
	}
	return g, nil
}

func (u *masterUseCase) GetGroups(ctx context.Context, limit, offset int) ([]*fullschedule.Group, error) {
	return u.repo.GetGroups(ctx, limit, offset)
}

func (u *masterUseCase) UpdateGroup(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*fullschedule.Group, error) {
	if _, err := u.repo.GetGroup(ctx, id); err != nil {
		return nil, fullschedule.ErrGroupNotFound
	}
	clean := make(map[string]interface{})
	for k, v := range values {
		switch k {
		case "name":
			name, _ := v.(string)
			if err := validateMasterName(name); err != nil {
				return nil, err
			}
			clean[k] = name
		case "description":
			if s, ok := v.(string); ok {
				clean[k] = s
			}
		case "sort_id":
			if n, ok := extractSortID(v); ok {
				clean[k] = n
			}
		case "status":
			s, _ := v.(string)
			if s != string(fullschedule.StatusActive) && s != string(fullschedule.StatusInactive) {
				return nil, ErrStatusInvalid
			}
			clean[k] = fullschedule.ScheduleStatus(s)
		}
	}
	if len(clean) == 0 {
		return nil, errors.New("fullschedule: nothing to update")
	}
	return u.repo.UpdateGroup(ctx, &fullschedule.Group{ID: id}, clean)
}

func (u *masterUseCase) DeleteGroup(ctx context.Context, id uuid.UUID) (*fullschedule.Group, error) {
	g, err := u.repo.GetGroup(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrGroupNotFound
	}
	count, err := u.repo.CountZonesByGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fullschedule.ErrGroupHasZones
	}
	_, err = u.repo.UpdateGroup(ctx, g, map[string]interface{}{"deleted_at": time.Now()})
	if err != nil {
		return nil, err
	}
	return u.repo.GetGroup(ctx, id)
}

func (u *masterUseCase) CountGroups(ctx context.Context) (int64, error) {
	return u.repo.CountGroups(ctx)
}

// ============================================================================
// Zone
// ============================================================================

func (u *masterUseCase) CreateZone(ctx context.Context, z *fullschedule.Zone) (*fullschedule.Zone, error) {
	if err := validateMasterName(z.Name); err != nil {
		return nil, err
	}
	if z.GroupID == uuid.Nil {
		return nil, ErrScopeIncomplete
	}
	if _, err := u.repo.GetGroup(ctx, z.GroupID); err != nil {
		return nil, fullschedule.ErrGroupNotFound
	}
	return u.repo.CreateZone(ctx, z)
}

func (u *masterUseCase) GetZone(ctx context.Context, id uuid.UUID) (*fullschedule.Zone, error) {
	z, err := u.repo.GetZone(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrZoneNotFound
	}
	return z, nil
}

func (u *masterUseCase) GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*fullschedule.Zone, error) {
	return u.repo.GetZones(ctx, groupID, limit, offset)
}

func (u *masterUseCase) CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error) {
	return u.repo.CountZones(ctx, groupID)
}

func (u *masterUseCase) UpdateZone(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*fullschedule.Zone, error) {
	z, err := u.repo.GetZone(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrZoneNotFound
	}
	clean := make(map[string]interface{})
	for k, v := range values {
		switch k {
		case "name":
			name, _ := v.(string)
			if err := validateMasterName(name); err != nil {
				return nil, err
			}
			clean[k] = name
		case "description":
			if s, ok := v.(string); ok {
				clean[k] = s
			}
		case "sort_id":
			if n, ok := extractSortID(v); ok {
				clean[k] = n
			}
		case "group_id":
			gid, err := uuidParse(v)
			if err != nil {
				return nil, err
			}
			if _, err := u.repo.GetGroup(ctx, gid); err != nil {
				return nil, fullschedule.ErrGroupNotFound
			}
			clean[k] = gid
			z.GroupID = gid
		}
	}
	if len(clean) == 0 {
		return nil, errors.New("fullschedule: nothing to update")
	}
	return u.repo.UpdateZone(ctx, z, clean)
}

func (u *masterUseCase) DeleteZone(ctx context.Context, id uuid.UUID) (*fullschedule.Zone, error) {
	z, err := u.repo.GetZone(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrZoneNotFound
	}
	count, err := u.repo.CountAreasByZone(ctx, id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fullschedule.ErrZoneHasAreas
	}
	_, err = u.repo.UpdateZone(ctx, z, map[string]interface{}{"deleted_at": time.Now()})
	if err != nil {
		return nil, err
	}
	return u.repo.GetZone(ctx, id)
}

// ============================================================================
// Area
// ============================================================================

func (u *masterUseCase) CreateArea(ctx context.Context, a *fullschedule.Area) (*fullschedule.Area, error) {
	if err := validateMasterName(a.Name); err != nil {
		return nil, err
	}
	if a.ZoneID == uuid.Nil {
		return nil, ErrScopeIncomplete
	}
	if _, err := u.repo.GetZone(ctx, a.ZoneID); err != nil {
		return nil, fullschedule.ErrZoneNotFound
	}
	return u.repo.CreateArea(ctx, a)
}

func (u *masterUseCase) GetArea(ctx context.Context, id uuid.UUID) (*fullschedule.Area, error) {
	a, err := u.repo.GetArea(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrAreaNotFound
	}
	return a, nil
}

func (u *masterUseCase) GetAreas(ctx context.Context, zoneID *uuid.UUID, sort fullschedule.SortOrder, limit, offset int) ([]*fullschedule.Area, error) {
	if sort == "" {
		sort = fullschedule.SortOrderASC
	}
	return u.repo.GetAreas(ctx, zoneID, sort, limit, offset)
}

func (u *masterUseCase) CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error) {
	return u.repo.CountAreas(ctx, zoneID)
}

func (u *masterUseCase) UpdateArea(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*fullschedule.Area, error) {
	a, err := u.repo.GetArea(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrAreaNotFound
	}
	clean := make(map[string]interface{})
	for k, v := range values {
		switch k {
		case "name":
			name, _ := v.(string)
			if err := validateMasterName(name); err != nil {
				return nil, err
			}
			clean[k] = name
		case "description":
			if s, ok := v.(string); ok {
				clean[k] = s
			}
		case "sort_id":
			if n, ok := extractSortID(v); ok {
				clean[k] = n
			}
		case "zone_id":
			zid, err := uuidParse(v)
			if err != nil {
				return nil, err
			}
			if _, err := u.repo.GetZone(ctx, zid); err != nil {
				return nil, fullschedule.ErrZoneNotFound
			}
			clean[k] = zid
			a.ZoneID = zid
		}
	}
	if len(clean) == 0 {
		return nil, errors.New("fullschedule: nothing to update")
	}
	return u.repo.UpdateArea(ctx, a, clean)
}

func (u *masterUseCase) DeleteArea(ctx context.Context, id uuid.UUID) (*fullschedule.Area, error) {
	a, err := u.repo.GetArea(ctx, id)
	if err != nil {
		return nil, fullschedule.ErrAreaNotFound
	}
	count, err := u.repo.CountMappingsByArea(ctx, id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fullschedule.ErrAreaHasDevices
	}
	_, err = u.repo.UpdateArea(ctx, a, map[string]interface{}{"deleted_at": time.Now()})
	if err != nil {
		return nil, err
	}
	return u.repo.GetArea(ctx, id)
}

// ============================================================================
// Device <-> Area mapping
// ============================================================================

func (u *masterUseCase) MapDevicesToArea(ctx context.Context, areaID uuid.UUID, deviceIDs []int) error {
	if _, err := u.repo.GetArea(ctx, areaID); err != nil {
		return fullschedule.ErrAreaNotFound
	}
	return u.repo.MapDevicesToArea(ctx, areaID, deviceIDs)
}

func (u *masterUseCase) GetAreaDevices(ctx context.Context, areaID uuid.UUID) ([]*fullschedule.AreaDevice, error) {
	if _, err := u.repo.GetArea(ctx, areaID); err != nil {
		return nil, fullschedule.ErrAreaNotFound
	}
	return u.repo.GetAreaDevices(ctx, areaID)
}

// ============================================================================
// Scope
// ============================================================================

// ValidateScope verifies the declared chain zone⊆group and area⊆zone.
// A fully empty scope is valid (device_ids fallback).
func (u *masterUseCase) ValidateScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) error {
	switch {
	case areaID != nil:
		area, err := u.repo.GetArea(ctx, *areaID)
		if err != nil {
			return fullschedule.ErrAreaNotFound
		}
		if zoneID != nil && area.ZoneID != *zoneID {
			return fullschedule.ErrAreaNotInZone
		}
		zone, err := u.repo.GetZone(ctx, area.ZoneID)
		if err != nil {
			return fullschedule.ErrZoneNotFound
		}
		if groupID != nil && zone.GroupID != *groupID {
			return fullschedule.ErrZoneNotInGroup
		}
	case zoneID != nil:
		zone, err := u.repo.GetZone(ctx, *zoneID)
		if err != nil {
			return fullschedule.ErrZoneNotFound
		}
		if groupID != nil && zone.GroupID != *groupID {
			return fullschedule.ErrZoneNotInGroup
		}
	case groupID != nil:
		if _, err := u.repo.GetGroup(ctx, *groupID); err != nil {
			return fullschedule.ErrGroupNotFound
		}
	}
	return nil
}

// GetDeviceIDsByScope resolves target devices. Returns ErrNoScope when no
// scope ids are provided (callers should fall back to manual device_ids).
func (u *masterUseCase) GetDeviceIDsByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) ([]int, error) {
	if groupID == nil && zoneID == nil && areaID == nil {
		return nil, fullschedule.ErrNoScope
	}
	if err := u.ValidateScope(ctx, groupID, zoneID, areaID); err != nil {
		return nil, err
	}
	return u.repo.GetDeviceIDsByScope(ctx, groupID, zoneID, areaID)
}

// CountDevicesByScope counts devices for a scope preview.
func (u *masterUseCase) CountDevicesByScope(ctx context.Context, groupID, zoneID, areaID *uuid.UUID) (int64, error) {
	if groupID == nil && zoneID == nil && areaID == nil {
		return 0, nil
	}
	if err := u.ValidateScope(ctx, groupID, zoneID, areaID); err != nil {
		return 0, err
	}
	return u.repo.CountDevicesByScope(ctx, groupID, zoneID, areaID)
}

// ============================================================================
// helpers
// ============================================================================

func validateMasterName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrNameRequired
	}
	return nil
}

func uuidParse(v interface{}) (uuid.UUID, error) {
	s, ok := v.(string)
	if !ok {
		return uuid.UUID{}, ErrScopeIncomplete
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, ErrScopeIncomplete
	}
	return parsed, nil
}

func extractSortID(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
}
