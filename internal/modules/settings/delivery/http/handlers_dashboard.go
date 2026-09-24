package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"

	"github.com/go-chi/render"
)

var dashboardCols = []string{"location_id", "name", "config_data", "status"}

// CreateDashboardConfig ports POST /dashboardconfig -> createDashboardConfig:
// an upsert keyed on (name, location_id). When a matching row exists its
// config_data is replaced; otherwise a new row is inserted.
// CreateDashboardConfig - POST /api/settings/dashboardconfig
// @Summary Create Dashboard Config
// @Description Upserts one row in sd_dashboard_config by name + location_id (body: name, location_id, config).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "name, location_id, config"
// @Success 200 {object} responses.SwaggerSuccessResponse "Upserted row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/dashboardconfig [post]
func (h *settingsHandler) CreateDashboardConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		name := toStr(body["name"])
		locationID := toStr(body["location_id"])
		config, hasConfig := body["config"]
		if name == "" || locationID == "" || !hasConfig {
			badRequest(w, r, "name, location_id and config are required")
			return
		}
		spec := settings.DashboardConfigSpec.With(
			settings.FixedFilter{Column: "dc.name", Value: name},
			settings.FixedFilter{Column: "dc.location_id", Value: locationID},
		)
		items, err := h.uc.RowsAll(r.Context(), spec, settings.NewListQuery(1, 1, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		if len(items) > 0 {
			fields := map[string]interface{}{"config_data": config, "updated_date": nowPtr()}
			if _, err = h.uc.UpdateFields(r.Context(), "sd_dashboard_config", "id", items[0]["id"], fields); err != nil {
				fail(w, r, err)
				return
			}
			items, err = h.uc.RowsAll(r.Context(), spec, settings.NewListQuery(1, 1, "", nil))
			if err != nil {
				fail(w, r, err)
				return
			}
			ok(w, r, items[0])
			return
		}
		row := map[string]interface{}{
			"location_id":  coerceID(locationID),
			"name":         name,
			"config_data":  config,
			"created_date": nowPtr(),
			"updated_date": nowPtr(),
		}
		if err = h.uc.CreateRow(r.Context(), "sd_dashboard_config", row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}

// DashboardConfigByLocation ports GET /dashboardconfig_1 (location_id required).
// DashboardConfigByLocation - GET /api/settings/dashboardconfig_1
// @Summary Dashboard Config By Location
// @Description Fetches one record by location_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param location_id query string true "location_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/dashboardconfig_1 [get]
func (h *settingsHandler) DashboardConfigByLocation() http.HandlerFunc {
	return h.getOneByParam(settings.DashboardConfigSpec, "location_id", "dc.location_id", true, "location_id is required")
}

// ListDashboardConfig - GET /api/settings/dashboardconfig
// @Summary List Dashboard Config
// @Description Paginated listing (filters via query params, sort=field-ASC|DESC).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Items per page (default 1000, max 5000)"
// @Param sort query string false "Sort: field-ASC or field-DESC"
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult] "Paginated result envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/dashboardconfig [get]
func (h *settingsHandler) ListDashboardConfig() http.HandlerFunc {
	return h.listHandler(settings.DashboardConfigSpec)
}

// FindOrCreateDashboardConfig ports GET /dashboardconfig/search -> findByCriteria:
// returns the first match for name+location_id, creating it when absent.
// FindOrCreateDashboardConfig - GET /api/settings/dashboardconfig/search
// @Summary Find Or Create Dashboard Config
// @Description Finds a dashboard config by name + location_id, creating a default one when absent.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param name query string true "Config name"
// @Param location_id query int true "Location ID"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/dashboardconfig/search [get]
func (h *settingsHandler) FindOrCreateDashboardConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		name := q.Get("name")
		locationID := q.Get("location_id")
		if name == "" || locationID == "" {
			render422(w, r, "name and location_id are required")
			return
		}
		spec := settings.DashboardConfigSpec.With(
			settings.FixedFilter{Column: "dc.name", Value: name},
			settings.FixedFilter{Column: "dc.location_id", Value: locationID},
		)
		items, err := h.uc.RowsAll(r.Context(), spec, settings.NewListQuery(1, 1, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		if len(items) > 0 {
			ok(w, r, items[0])
			return
		}
		row := map[string]interface{}{
			"location_id": coerceID(locationID),
			"name":        name,
			"config_data": map[string]interface{}{},
			"status":      1,
		}
		if err = h.uc.CreateRow(r.Context(), "sd_dashboard_config", row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}

// GetDashboardConfig - GET /api/settings/dashboardconfig/{id}
// @Summary Get Dashboard Config
// @Description Fetches one record by id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/dashboardconfig/{id} [get]
func (h *settingsHandler) GetDashboardConfig() http.HandlerFunc {
	return h.getOneByParam(settings.DashboardConfigSpec, "id", "dc.id", true, "dashboard config not found")
}

// UpdateDashboardConfig ports PATCH /dashboardconfig/:id.
// UpdateDashboardConfig - PATCH /api/settings/dashboardconfig/{id}
// @Summary Update Dashboard Config
// @Description Partial update of one dashboard config by path id.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/dashboardconfig/{id} [patch]
func (h *settingsHandler) UpdateDashboardConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := urlID(r)
		if id == "" {
			badRequest(w, r, "id is required")
			return
		}
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		fields := filterKeys(body, dashboardCols)
		fields["updated_date"] = nowPtr()
		n, err := h.uc.UpdateFields(r.Context(), "sd_dashboard_config", "id", coerceID(id), fields)
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// RemoveDashboardConfig - DELETE /api/settings/dashboardconfig/{id}
// @Summary Remove Dashboard Config
// @Description Deletes one row by path id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param id path string true "Row ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/dashboardconfig/{id} [delete]
func (h *settingsHandler) RemoveDashboardConfig() http.HandlerFunc {
	return h.deleteParamHandler("sd_dashboard_config", "id")
}

func render422(w http.ResponseWriter, r *http.Request, msg string) {
	render.Render(w, r, response422(msg))
}
