package http

import (
	"encoding/json"
	"net/http"

	"icmongolang/config"
	"icmongolang/internal/modules/fullschedule"
	"icmongolang/internal/modules/fullschedule/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type masterHandler struct {
	cfg    *config.Config
	uc     fullschedule.MasterUseCaseI
	logger logger.Logger
}

// CreateMasterHandler creates the master data HTTP handler.
// สร้างฮาร์ดเลอร์สำหรับ Master Data (Group/Zone/Area)
func CreateMasterHandler(uc fullschedule.MasterUseCaseI, cfg *config.Config, logger logger.Logger) fullschedule.MasterHandlers {
	return &masterHandler{cfg: cfg, uc: uc, logger: logger}
}

// ============================================================================
// Group
// ============================================================================

// CreateGroup godoc
// @Summary Create a group
// @Description Create a top-level organization group.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param group body presenter.GroupCreate true "Group details"
// @Success 200 {object} responses.SuccessResponse[presenter.GroupResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/groups [post]
func (h *masterHandler) CreateGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.GroupCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		status := fullschedule.StatusActive
		if req.Status != "" {
			status = fullschedule.ScheduleStatus(req.Status)
		}
		g, err := h.uc.CreateGroup(r.Context(), &fullschedule.Group{
			Name:        req.Name,
			Description: req.Description,
			SortID:      req.SortID,
			Status:      status,
		})
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapGroupResponse(g)))
	}
}

// ListGroups godoc
// @Summary List groups
// @Description List groups with pagination.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedGroupResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/groups [get]
func (h *masterHandler) ListGroups() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, page, perPage := parsePagination(r)
		groups, err := h.uc.GetGroups(r.Context(), limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		total, err := h.uc.CountGroups(r.Context())
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := &presenter.PaginatedGroupResponse{
			Groups:     mapGroupsResponse(groups),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages(total, perPage),
		}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// GetGroup godoc
// @Summary Get a group
// @Description Get a group by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} responses.SuccessResponse[presenter.GroupResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/groups/{id} [get]
func (h *masterHandler) GetGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		g, err := h.uc.GetGroup(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapGroupResponse(g)))
	}
}

// UpdateGroup godoc
// @Summary Update a group
// @Description Update group fields by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param group body presenter.GroupUpdate true "Group update fields"
// @Success 200 {object} responses.SuccessResponse[presenter.GroupResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/groups/{id} [put]
func (h *masterHandler) UpdateGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.GroupUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		values := map[string]interface{}{}
		if req.Name != nil {
			values["name"] = *req.Name
		}
		if req.Description != nil {
			values["description"] = *req.Description
		}
		if req.SortID != nil {
			values["sort_id"] = *req.SortID
		}
		if req.Status != nil {
			values["status"] = *req.Status
		}
		g, err := h.uc.UpdateGroup(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapGroupResponse(g)))
	}
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Soft-delete a group (only if it has no child zones).
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} responses.SuccessResponse[presenter.GroupResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/groups/{id} [delete]
func (h *masterHandler) DeleteGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		g, err := h.uc.DeleteGroup(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapGroupResponse(g)))
	}
}

// ============================================================================
// Zone
// ============================================================================

// CreateZone godoc
// @Summary Create a zone
// @Description Create a zone under a group.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param zone body presenter.ZoneCreate true "Zone details"
// @Success 200 {object} responses.SuccessResponse[presenter.ZoneResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/zones [post]
func (h *masterHandler) CreateZone() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.ZoneCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		z, err := h.uc.CreateZone(r.Context(), &fullschedule.Zone{
			GroupID:     req.GroupID,
			Name:        req.Name,
			Description: req.Description,
			SortID:      req.SortID,
		})
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapZoneResponse(z)))
	}
}

// ListZones godoc
// @Summary List zones
// @Description List zones, optionally filtered by group_id.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param group_id query string false "Filter by group"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedZoneResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/zones [get]
func (h *masterHandler) ListZones() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, page, perPage := parsePagination(r)
		q := r.URL.Query()
		var groupID *uuid.UUID
		if raw := q.Get("group_id"); raw != "" {
			parsed, err := uuid.Parse(raw)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}
			groupID = &parsed
		}
		zones, err := h.uc.GetZones(r.Context(), groupID, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		total, err := h.uc.CountZones(r.Context(), groupID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(&presenter.PaginatedZoneResponse{
			Zones:      mapZonesResponse(zones),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages(total, perPage),
		}))
	}
}

// GetZone godoc
// @Summary Get a zone
// @Description Get a zone by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Zone ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ZoneResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/zones/{id} [get]
func (h *masterHandler) GetZone() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		z, err := h.uc.GetZone(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapZoneResponse(z)))
	}
}

// UpdateZone godoc
// @Summary Update a zone
// @Description Update zone fields by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Zone ID"
// @Param zone body presenter.ZoneUpdate true "Zone update fields"
// @Success 200 {object} responses.SuccessResponse[presenter.ZoneResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/zones/{id} [put]
func (h *masterHandler) UpdateZone() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.ZoneUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		values := map[string]interface{}{}
		if req.Name != nil {
			values["name"] = *req.Name
		}
		if req.Description != nil {
			values["description"] = *req.Description
		}
		if req.SortID != nil {
			values["sort_id"] = *req.SortID
		}
		if req.GroupID != nil {
			values["group_id"] = *req.GroupID
		}
		z, err := h.uc.UpdateZone(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapZoneResponse(z)))
	}
}

// DeleteZone godoc
// @Summary Delete a zone
// @Description Soft-delete a zone (only if it has no child areas).
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Zone ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ZoneResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/zones/{id} [delete]
func (h *masterHandler) DeleteZone() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		z, err := h.uc.DeleteZone(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapZoneResponse(z)))
	}
}

// ============================================================================
// Area
// ============================================================================

// CreateArea godoc
// @Summary Create an area
// @Description Create an area under a zone.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param area body presenter.AreaCreate true "Area details"
// @Success 200 {object} responses.SuccessResponse[presenter.AreaResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas [post]
func (h *masterHandler) CreateArea() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.AreaCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		a, err := h.uc.CreateArea(r.Context(), &fullschedule.Area{
			ZoneID:      req.ZoneID,
			Name:        req.Name,
			Description: req.Description,
			SortID:      req.SortID,
		})
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapAreaResponse(a)))
	}
}

// ListAreas godoc
// @Summary List areas
// @Description List areas, optionally filtered by zone_id. Order by sort_id (sort=asc/desc).
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param zone_id query string false "Filter by zone"
// @Param sort query string false "Sort order: asc (default) or desc"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} responses.SuccessResponse[presenter.PaginatedAreaResponse]
// @Failure 401 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas [get]
func (h *masterHandler) ListAreas() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, page, perPage := parsePagination(r)
		q := r.URL.Query()
		var zoneID *uuid.UUID
		if raw := q.Get("zone_id"); raw != "" {
			parsed, err := uuid.Parse(raw)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}
			zoneID = &parsed
		}
		sort := parseSortOrder(q.Get("sort"))
		areas, err := h.uc.GetAreas(r.Context(), zoneID, sort, limit, offset)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		total, err := h.uc.CountAreas(r.Context(), zoneID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(&presenter.PaginatedAreaResponse{
			Areas:      mapAreasResponse(areas),
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages(total, perPage),
		}))
	}
}

// GetArea godoc
// @Summary Get an area
// @Description Get an area by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Area ID"
// @Success 200 {object} responses.SuccessResponse[presenter.AreaResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas/{id} [get]
func (h *masterHandler) GetArea() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		a, err := h.uc.GetArea(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapAreaResponse(a)))
	}
}

// UpdateArea godoc
// @Summary Update an area
// @Description Update area fields by ID.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Area ID"
// @Param area body presenter.AreaUpdate true "Area update fields"
// @Success 200 {object} responses.SuccessResponse[presenter.AreaResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas/{id} [put]
func (h *masterHandler) UpdateArea() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.AreaUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		values := map[string]interface{}{}
		if req.Name != nil {
			values["name"] = *req.Name
		}
		if req.Description != nil {
			values["description"] = *req.Description
		}
		if req.SortID != nil {
			values["sort_id"] = *req.SortID
		}
		if req.ZoneID != nil {
			values["zone_id"] = *req.ZoneID
		}
		a, err := h.uc.UpdateArea(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapAreaResponse(a)))
	}
}

// DeleteArea godoc
// @Summary Delete an area
// @Description Soft-delete an area (only if it has no mapped devices).
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Area ID"
// @Success 200 {object} responses.SuccessResponse[presenter.AreaResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas/{id} [delete]
func (h *masterHandler) DeleteArea() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		a, err := h.uc.DeleteArea(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapAreaResponse(a)))
	}
}

// ============================================================================
// Device mapping + scope preview
// ============================================================================

// MapAreaDevices godoc
// @Summary Map devices to an area
// @Description Replace the IoT device mapping of an area (device_ids auto-resolves the schedule scope).
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Area ID"
// @Param body body presenter.DeviceMapRequest true "Device IDs"
// @Success 200 {object} responses.SuccessResponse[presenter.MessageResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas/{id}/devices [post]
func (h *masterHandler) MapAreaDevices() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		req := new(presenter.DeviceMapRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		if err := h.uc.MapDevicesToArea(r.Context(), id, req.DeviceIDs); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(map[string]interface{}{"message": "devices mapped"}))
	}
}

// AreaDevices godoc
// @Summary List devices of an area
// @Description List IoT devices mapped to an area.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param id path string true "Area ID"
// @Success 200 {object} responses.SuccessResponse[[]presenter.AreaDeviceResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/areas/{id}/devices [get]
func (h *masterHandler) AreaDevices() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}
		rows, err := h.uc.GetAreaDevices(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		out := make([]*presenter.AreaDeviceResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, &presenter.AreaDeviceResponse{DeviceID: row.DeviceID, DeviceSN: row.DeviceSN})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// ScopePreview godoc
// @Summary Preview device count for a scope
// @Description Count devices covered by a group/zone/area selection before saving a schedule.
// @Tags fullschedule-master
// @Accept json
// @Produce json
// @Param group_id query string false "Group ID"
// @Param zone_id query string false "Zone ID"
// @Param area_id query string false "Area ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ScopePreviewResponse]
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Security OAuth2Password
// @Security BearerAuth
// @Router /fullschedule/scope/preview [get]
func (h *masterHandler) ScopePreview() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var groupID, zoneID, areaID *uuid.UUID
		for key, dst := range map[string]**uuid.UUID{
			"group_id": &groupID,
			"zone_id":  &zoneID,
			"area_id":  &areaID,
		} {
			raw := firstParam(q, key)
			if raw == "" {
				continue
			}
			parsed, err := uuid.Parse(raw)
			if err != nil {
				render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
				return
			}
			*dst = &parsed
		}
		count, err := h.uc.CountDevicesByScope(r.Context(), groupID, zoneID, areaID)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		res := &presenter.ScopePreviewResponse{GroupID: groupID, ZoneID: zoneID, AreaID: areaID, DeviceCount: count}
		render.Respond(w, r, responses.CreateSuccessResponse(res))
	}
}

// ============================================================================
// helpers
// ============================================================================

func totalPages(total int64, perPage int) int {
	if perPage <= 0 {
		perPage = 1
	}
	pages := int(total) / perPage
	if int(total)%perPage != 0 {
		pages++
	}
	if pages < 1 {
		pages = 1
	}
	return pages
}

func mapGroupResponse(g *fullschedule.Group) *presenter.GroupResponse {
	return &presenter.GroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		SortID:      g.SortID,
		Status:      string(g.Status),
		CreatedBy:   g.CreatedBy,
		CreatedAt:   helpers.NewLocalTime(g.CreatedAt),
		UpdatedBy:   g.UpdatedBy,
		UpdatedAt:   helpers.NewLocalTime(g.UpdatedAt),
		Version:     g.Version,
	}
}

func mapGroupsResponse(list []*fullschedule.Group) []*presenter.GroupResponse {
	out := make([]*presenter.GroupResponse, 0, len(list))
	for _, g := range list {
		out = append(out, mapGroupResponse(g))
	}
	return out
}

func mapZoneResponse(z *fullschedule.Zone) *presenter.ZoneResponse {
	return &presenter.ZoneResponse{
		ID:          z.ID,
		GroupID:     z.GroupID,
		GroupName:   z.GroupName,
		Name:        z.Name,
		Description: z.Description,
		SortID:      z.SortID,
		CreatedBy:   z.CreatedBy,
		CreatedAt:   helpers.NewLocalTime(z.CreatedAt),
		UpdatedBy:   z.UpdatedBy,
		UpdatedAt:   helpers.NewLocalTime(z.UpdatedAt),
		Version:     z.Version,
	}
}

func mapZonesResponse(list []*fullschedule.Zone) []*presenter.ZoneResponse {
	out := make([]*presenter.ZoneResponse, 0, len(list))
	for _, z := range list {
		out = append(out, mapZoneResponse(z))
	}
	return out
}

func mapAreaResponse(a *fullschedule.Area) *presenter.AreaResponse {
	return &presenter.AreaResponse{
		ID:          a.ID,
		ZoneID:      a.ZoneID,
		ZoneName:    a.ZoneName,
		GroupName:   a.GroupName,
		Name:        a.Name,
		Description: a.Description,
		SortID:      a.SortID,
		CreatedBy:   a.CreatedBy,
		CreatedAt:   helpers.NewLocalTime(a.CreatedAt),
		UpdatedBy:   a.UpdatedBy,
		UpdatedAt:   helpers.NewLocalTime(a.UpdatedAt),
		Version:     a.Version,
	}
}

func mapAreasResponse(list []*fullschedule.Area) []*presenter.AreaResponse {
	out := make([]*presenter.AreaResponse, 0, len(list))
	for _, a := range list {
		out = append(out, mapAreaResponse(a))
	}
	return out
}
