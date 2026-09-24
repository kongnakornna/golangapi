package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	"icmongolang/internal/modules/fullschedule/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// memMasterRepo implements fullschedule.MasterPgRepository in-memory.
// Menyimpan Group/Zone/Area แบบ simple map ไม่มี DB จริง.
type memMasterRepo struct {
	fullschedule.MasterPgRepository

	groups map[uuid.UUID]*fullschedule.Group
	zones  map[uuid.UUID]*fullschedule.Zone
	areas  map[uuid.UUID]*fullschedule.Area

	lastAreaSort fullschedule.SortOrder
}

func newMemRepo() *memMasterRepo {
	return &memMasterRepo{
		groups: map[uuid.UUID]*fullschedule.Group{},
		zones:  map[uuid.UUID]*fullschedule.Zone{},
		areas:  map[uuid.UUID]*fullschedule.Area{},
	}
}

func (m *memMasterRepo) CreateGroup(ctx context.Context, g *fullschedule.Group) (*fullschedule.Group, error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	cp := *g
	m.groups[cp.ID] = &cp
	return &cp, nil
}

func (m *memMasterRepo) GetGroup(ctx context.Context, id uuid.UUID) (*fullschedule.Group, error) {
	g, ok := m.groups[id]
	if !ok || g.DeletedAt != nil {
		return nil, errors.New("not found")
	}
	cp := *g
	return &cp, nil
}

func (m *memMasterRepo) GetGroups(ctx context.Context, limit, offset int) ([]*fullschedule.Group, error) {
	var out []*fullschedule.Group
	for _, g := range m.groups {
		if g.DeletedAt == nil {
			cp := *g
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (m *memMasterRepo) CountGroups(ctx context.Context) (int64, error) {
	var n int64
	for _, g := range m.groups {
		if g.DeletedAt == nil {
			n++
		}
	}
	return n, nil
}

func (m *memMasterRepo) UpdateGroup(ctx context.Context, g *fullschedule.Group, values map[string]interface{}) (*fullschedule.Group, error) {
	stored, ok := m.groups[g.ID]
	if !ok {
		return nil, errors.New("not found")
	}
	applyGroup(stored, values)
	cp := *stored
	return &cp, nil
}

func (m *memMasterRepo) CountZonesByGroup(ctx context.Context, groupID uuid.UUID) (int64, error) {
	var n int64
	for _, z := range m.zones {
		if z.GroupID == groupID && z.DeletedAt == nil {
			n++
		}
	}
	return n, nil
}

func (m *memMasterRepo) CreateZone(ctx context.Context, z *fullschedule.Zone) (*fullschedule.Zone, error) {
	if z.ID == uuid.Nil {
		z.ID = uuid.New()
	}
	cp := *z
	m.zones[cp.ID] = &cp
	return &cp, nil
}

func (m *memMasterRepo) GetZone(ctx context.Context, id uuid.UUID) (*fullschedule.Zone, error) {
	z, ok := m.zones[id]
	if !ok || z.DeletedAt != nil {
		return nil, errors.New("not found")
	}
	cp := *z
	return &cp, nil
}

func (m *memMasterRepo) GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*fullschedule.Zone, error) {
	var out []*fullschedule.Zone
	for _, z := range m.zones {
		if z.DeletedAt != nil {
			continue
		}
		if groupID != nil && z.GroupID != *groupID {
			continue
		}
		cp := *z
		out = append(out, &cp)
	}
	return out, nil
}

func (m *memMasterRepo) UpdateZone(ctx context.Context, z *fullschedule.Zone, values map[string]interface{}) (*fullschedule.Zone, error) {
	stored, ok := m.zones[z.ID]
	if !ok {
		return nil, errors.New("not found")
	}
	applyZone(stored, values)
	cp := *stored
	return &cp, nil
}

func (m *memMasterRepo) CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error) {
	var n int64
	for _, z := range m.zones {
		if z.DeletedAt != nil {
			continue
		}
		if groupID != nil && z.GroupID != *groupID {
			continue
		}
		n++
	}
	return n, nil
}

func (m *memMasterRepo) CountAreasByZone(ctx context.Context, zoneID uuid.UUID) (int64, error) {
	var n int64
	for _, a := range m.areas {
		if a.ZoneID == zoneID && a.DeletedAt == nil {
			n++
		}
	}
	return n, nil
}

func (m *memMasterRepo) CreateArea(ctx context.Context, a *fullschedule.Area) (*fullschedule.Area, error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	cp := *a
	m.areas[cp.ID] = &cp
	return &cp, nil
}

func (m *memMasterRepo) GetArea(ctx context.Context, id uuid.UUID) (*fullschedule.Area, error) {
	a, ok := m.areas[id]
	if !ok || a.DeletedAt != nil {
		return nil, errors.New("not found")
	}
	cp := *a
	return &cp, nil
}

func (m *memMasterRepo) GetAreas(ctx context.Context, zoneID *uuid.UUID, sort fullschedule.SortOrder, limit, offset int) ([]*fullschedule.Area, error) {
	m.lastAreaSort = sort
	var out []*fullschedule.Area
	for _, a := range m.areas {
		if a.DeletedAt != nil {
			continue
		}
		if zoneID != nil && a.ZoneID != *zoneID {
			continue
		}
		cp := *a
		out = append(out, &cp)
	}
	return out, nil
}

func (m *memMasterRepo) UpdateArea(ctx context.Context, a *fullschedule.Area, values map[string]interface{}) (*fullschedule.Area, error) {
	stored, ok := m.areas[a.ID]
	if !ok {
		return nil, errors.New("not found")
	}
	applyArea(stored, values)
	cp := *stored
	return &cp, nil
}

func (m *memMasterRepo) CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error) {
	var n int64
	for _, a := range m.areas {
		if a.DeletedAt != nil {
			continue
		}
		if zoneID != nil && a.ZoneID != *zoneID {
			continue
		}
		n++
	}
	return n, nil
}

func strPtr(s string) *string { return &s }

func applyGroup(g *fullschedule.Group, values map[string]interface{}) {
	if v, ok := values["name"]; ok {
		g.Name = v.(string)
	}
	if v, ok := values["description"]; ok {
		g.Description = strPtr(v.(string))
	}
	if v, ok := values["sort_id"]; ok {
		g.SortID = v.(int)
	}
	if v, ok := values["status"]; ok {
		g.Status = fullschedule.ScheduleStatus(v.(string))
	}
	if v, ok := values["deleted_at"]; ok {
		switch tv := v.(type) {
		case time.Time:
			g.DeletedAt = &tv
		case *time.Time:
			g.DeletedAt = tv
		}
	}
}

func applyZone(z *fullschedule.Zone, values map[string]interface{}) {
	if v, ok := values["name"]; ok {
		z.Name = v.(string)
	}
	if v, ok := values["description"]; ok {
		z.Description = strPtr(v.(string))
	}
	if v, ok := values["sort_id"]; ok {
		z.SortID = v.(int)
	}
	if v, ok := values["group_id"]; ok {
		z.GroupID = v.(uuid.UUID)
	}
	if v, ok := values["deleted_at"]; ok {
		switch tv := v.(type) {
		case time.Time:
			z.DeletedAt = &tv
		case *time.Time:
			z.DeletedAt = tv
		}
	}
}

func applyArea(a *fullschedule.Area, values map[string]interface{}) {
	if v, ok := values["name"]; ok {
		a.Name = v.(string)
	}
	if v, ok := values["description"]; ok {
		a.Description = strPtr(v.(string))
	}
	if v, ok := values["sort_id"]; ok {
		a.SortID = v.(int)
	}
	if v, ok := values["zone_id"]; ok {
		a.ZoneID = v.(uuid.UUID)
	}
	if v, ok := values["deleted_at"]; ok {
		switch tv := v.(type) {
		case time.Time:
			a.DeletedAt = &tv
		case *time.Time:
			a.DeletedAt = tv
		}
	}
}

func newMasterUC(repo fullschedule.MasterPgRepository) fullschedule.MasterUseCaseI {
	return usecase.CreateMasterUseCaseI(repo, &config.Config{}, logger.NewApiLogger(&config.Config{}))
}

func TestCreateGroup_Validation(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	if _, err := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "  "}); err != usecase.ErrNameRequired {
		t.Fatalf("CreateGroup(empty name) error = %v, want ErrNameRequired", err)
	}

	if _, err := uc.CreateGroup(context.Background(), &fullschedule.Group{
		Name:   "G",
		Status: fullschedule.ScheduleStatus("bogus"),
	}); err != usecase.ErrStatusInvalid {
		t.Fatalf("CreateGroup(bad status) error = %v, want ErrStatusInvalid", err)
	}
}

func TestCreateZone_RequiresExistingGroup(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	_, err := uc.CreateZone(context.Background(), &fullschedule.Zone{
		GroupID: uuid.New(),
		Name:    "Zone A",
	})
	if err != fullschedule.ErrGroupNotFound {
		t.Fatalf("CreateZone(bad group) error = %v, want ErrGroupNotFound", err)
	}
}

func TestCreateZone_SuccessWithSortID(t *testing.T) {
	repo := newMemRepo()
	uc := newMasterUC(repo)

	g, err := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "Group 1"})
	if err != nil {
		t.Fatalf("CreateGroup error = %v", err)
	}

	z, err := uc.CreateZone(context.Background(), &fullschedule.Zone{
		GroupID: g.ID,
		Name:    "Zone A",
		SortID:  5,
	})
	if err != nil {
		t.Fatalf("CreateZone error = %v", err)
	}
	if z.SortID != 5 {
		t.Fatalf("CreateZone SortID = %d, want 5", z.SortID)
	}
	got, err := repo.GetZone(context.Background(), z.ID)
	if err != nil {
		t.Fatalf("stored zone error = %v", err)
	}
	if got.SortID != 5 {
		t.Fatalf("stored zone SortID = %d, want 5", got.SortID)
	}
}

func TestDeleteGroup_BlockedWhenHasZones(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	if _, err := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Z"}); err != nil {
		t.Fatalf("CreateZone error = %v", err)
	}

	if _, err := uc.DeleteGroup(context.Background(), g.ID); err != fullschedule.ErrGroupHasZones {
		t.Fatalf("DeleteGroup error = %v, want ErrGroupHasZones", err)
	}
}

func TestDeleteGroup_SoftDeletes(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	uc.DeleteGroup(context.Background(), g.ID)
	if _, err := uc.GetGroup(context.Background(), g.ID); err != fullschedule.ErrGroupNotFound {
		t.Fatalf("GetGroup after delete error = %v, want ErrGroupNotFound", err)
	}
}

func TestUpdateZone_SortIDAndName(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	z, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Zone A", SortID: 1})

	updated, err := uc.UpdateZone(context.Background(), z.ID, map[string]interface{}{
		"name":    "Zone A+",
		"sort_id": 7,
	})
	if err != nil {
		t.Fatalf("UpdateZone error = %v", err)
	}
	if updated.Name != "Zone A+" {
		t.Fatalf("updated name = %q, want Zone A+", updated.Name)
	}
	if updated.SortID != 7 {
		t.Fatalf("updated SortID = %d, want 7", updated.SortID)
	}
}

func TestUpdateZone_InvalidSortIDTypeIgnored(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	z, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Zone A"})

	updated, err := uc.UpdateZone(context.Background(), z.ID, map[string]interface{}{
		"sort_id": "not-a-number",
	})
	if err == nil {
		t.Fatalf("UpdateZone with only invalid sort_id should error (nothing to update), got nil")
	}
	if updated != nil {
		t.Fatalf("expected nil zone on error, got %+v", updated)
	}
}

func TestCountZones_FilteredByGroup(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g1, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G1"})
	g2, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G2"})
	for i := 0; i < 3; i++ {
		if _, err := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g1.ID, Name: "Z"}); err != nil {
			t.Fatalf("CreateZone error = %v", err)
		}
	}
	if _, err := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g2.ID, Name: "Z"}); err != nil {
		t.Fatalf("CreateZone error = %v", err)
	}

	total, err := uc.CountZones(context.Background(), &g1.ID)
	if err != nil {
		t.Fatalf("CountZones error = %v", err)
	}
	if total != 3 {
		t.Fatalf("CountZones(g1) = %d, want 3", total)
	}

	all, err := uc.CountZones(context.Background(), nil)
	if err != nil {
		t.Fatalf("CountZones(all) error = %v", err)
	}
	if all != 4 {
		t.Fatalf("CountZones(all) = %d, want 4", all)
	}
}

func TestDeleteZone_BlockedWhenHasAreas(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	z, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Z"})
	if _, err := uc.CreateArea(context.Background(), &fullschedule.Area{ZoneID: z.ID, Name: "A"}); err != nil {
		t.Fatalf("CreateArea error = %v", err)
	}

	if _, err := uc.DeleteZone(context.Background(), z.ID); err != fullschedule.ErrZoneHasAreas {
		t.Fatalf("DeleteZone error = %v, want ErrZoneHasAreas", err)
	}
}

func TestCountAreas_FilteredByZone(t *testing.T) {
	uc := newMasterUC(newMemRepo())

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	z1, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Z1"})
	z2, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Z2"})
	for i := 0; i < 2; i++ {
		if _, err := uc.CreateArea(context.Background(), &fullschedule.Area{ZoneID: z1.ID, Name: "A"}); err != nil {
			t.Fatalf("CreateArea error = %v", err)
		}
	}
	if _, err := uc.CreateArea(context.Background(), &fullschedule.Area{ZoneID: z2.ID, Name: "A"}); err != nil {
		t.Fatalf("CreateArea error = %v", err)
	}

	total, err := uc.CountAreas(context.Background(), &z1.ID)
	if err != nil {
		t.Fatalf("CountAreas error = %v", err)
	}
	if total != 2 {
		t.Fatalf("CountAreas(z1) = %d, want 2", total)
	}
}

func TestGetAreas_SortOrderPassThrough(t *testing.T) {
	repo := newMemRepo()
	uc := newMasterUC(repo)

	g, _ := uc.CreateGroup(context.Background(), &fullschedule.Group{Name: "G"})
	z, _ := uc.CreateZone(context.Background(), &fullschedule.Zone{GroupID: g.ID, Name: "Z"})
	if _, err := uc.CreateArea(context.Background(), &fullschedule.Area{ZoneID: z.ID, Name: "A"}); err != nil {
		t.Fatalf("CreateArea error = %v", err)
	}

	if _, err := uc.GetAreas(context.Background(), nil, fullschedule.SortOrderDESC, 10, 0); err != nil {
		t.Fatalf("GetAreas(desc) error = %v", err)
	}
	if repo.lastAreaSort != fullschedule.SortOrderDESC {
		t.Fatalf("repo sort = %q, want %q", repo.lastAreaSort, fullschedule.SortOrderDESC)
	}

	// Empty sort should be normalized to ASC by the usecase.
	if _, err := uc.GetAreas(context.Background(), nil, "", 10, 0); err != nil {
		t.Fatalf("GetAreas(default) error = %v", err)
	}
	if repo.lastAreaSort != fullschedule.SortOrderASC {
		t.Fatalf("repo sort = %q, want default %q", repo.lastAreaSort, fullschedule.SortOrderASC)
	}
}
