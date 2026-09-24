package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	fshttp "icmongolang/internal/modules/fullschedule/delivery/http"
	"icmongolang/internal/modules/fullschedule/presenter"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// listZonesStub stubs only the usecase methods invoked by ListZones.
type listZonesStub struct {
	fullschedule.MasterUseCaseI
	zones   []*fullschedule.Zone
	total   int64
	groupID *uuid.UUID
}

func (s *listZonesStub) GetZones(ctx context.Context, groupID *uuid.UUID, limit, offset int) ([]*fullschedule.Zone, error) {
	if s.groupID != nil && (groupID == nil || *groupID != *s.groupID) {
		return nil, nil
	}
	return s.zones, nil
}

func (s *listZonesStub) CountZones(ctx context.Context, groupID *uuid.UUID) (int64, error) {
	return s.total, nil
}

func TestListZones_CorrectPagination(t *testing.T) {
	groupID := uuid.New()

	zones := make([]*fullschedule.Zone, 5)
	for i := range zones {
		zones[i] = &fullschedule.Zone{
			ID:      uuid.New(),
			GroupID: groupID,
			Name:    "Zone",
			SortID:  i,
		}
	}

	uc := &listZonesStub{zones: zones, total: 23, groupID: &groupID}
	h := fshttp.CreateMasterHandler(uc, &config.Config{}, logger.NewApiLogger(&config.Config{}))

	req := httptest.NewRequest(http.MethodGet, "/api/zones?group_id="+groupID.String()+"&page=2&per_page=5", nil)
	rec := httptest.NewRecorder()

	h.ListZones()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ListZones status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	var body struct {
		Data presenter.PaginatedZoneResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error = %v; body=%s", err, rec.Body.String())
	}

	if body.Data.Total != 23 {
		t.Fatalf("Total = %d, want 23 (real count, not len(zones)=5)", body.Data.Total)
	}
	if body.Data.Page != 2 {
		t.Fatalf("Page = %d, want 2", body.Data.Page)
	}
	if body.Data.PerPage != 5 {
		t.Fatalf("PerPage = %d, want 5", body.Data.PerPage)
	}
	if body.Data.TotalPages != 5 {
		t.Fatalf("TotalPages = %d, want 5", body.Data.TotalPages)
	}
	if len(body.Data.Zones) != 5 {
		t.Fatalf("len(Zones) = %d, want 5", len(body.Data.Zones))
	}
	for i, z := range body.Data.Zones {
		if z.SortID != i {
			t.Fatalf("Zones[%d].SortID = %d, want %d", i, z.SortID, i)
		}
	}
}

func TestListZones_InvalidGroupID(t *testing.T) {
	h := fshttp.CreateMasterHandler(&listZonesStub{}, &config.Config{}, logger.NewApiLogger(&config.Config{}))

	req := httptest.NewRequest(http.MethodGet, "/api/zones?group_id=not-a-uuid", nil)
	rec := httptest.NewRecorder()

	h.ListZones()(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ListZones(bad group_id) status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
}

// listAreasStub stubs GetAreas/CountAreas for ListAreas.
type listAreasStub struct {
	fullschedule.MasterUseCaseI
	areas   []*fullschedule.Area
	total   int64
	gotSort fullschedule.SortOrder
}

func (s *listAreasStub) GetAreas(ctx context.Context, zoneID *uuid.UUID, sort fullschedule.SortOrder, limit, offset int) ([]*fullschedule.Area, error) {
	s.gotSort = sort
	return s.areas, nil
}

func (s *listAreasStub) CountAreas(ctx context.Context, zoneID *uuid.UUID) (int64, error) {
	return s.total, nil
}

func TestListAreas_GroupNameAndSort(t *testing.T) {
	zoneID := uuid.New()

	uc := &listAreasStub{
		areas: []*fullschedule.Area{{
			ID:        uuid.New(),
			ZoneID:    zoneID,
			ZoneName:  "Zone A",
			GroupName: "Group 1",
			Name:      "Area 1",
			SortID:    3,
		}},
		total: 1,
	}
	h := fshttp.CreateMasterHandler(uc, &config.Config{}, logger.NewApiLogger(&config.Config{}))

	req := httptest.NewRequest(http.MethodGet, "/api/areas?zone_id="+zoneID.String()+"&sort=desc", nil)
	rec := httptest.NewRecorder()

	h.ListAreas()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ListAreas status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if uc.gotSort != fullschedule.SortOrderDESC {
		t.Fatalf("uc sort = %q, want desc", uc.gotSort)
	}

	var body struct {
		Data presenter.PaginatedAreaResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error = %v; body=%s", err, rec.Body.String())
	}
	if len(body.Data.Areas) != 1 {
		t.Fatalf("len(Areas) = %d, want 1", len(body.Data.Areas))
	}
	a := body.Data.Areas[0]
	if a.GroupName != "Group 1" {
		t.Fatalf("group_name = %q, want Group 1", a.GroupName)
	}
	if a.ZoneName != "Zone A" {
		t.Fatalf("zone_name = %q, want Zone A", a.ZoneName)
	}
	if a.SortID != 3 {
		t.Fatalf("sort_id = %d, want 3", a.SortID)
	}
}

func TestListAreas_DefaultSortAsc(t *testing.T) {
	uc := &listAreasStub{}
	h := fshttp.CreateMasterHandler(uc, &config.Config{}, logger.NewApiLogger(&config.Config{}))

	req := httptest.NewRequest(http.MethodGet, "/api/areas", nil)
	rec := httptest.NewRecorder()

	h.ListAreas()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ListAreas status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if uc.gotSort != fullschedule.SortOrderASC {
		t.Fatalf("uc sort = %q, want default asc", uc.gotSort)
	}
}
