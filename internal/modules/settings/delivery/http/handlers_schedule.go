package http

import (
	"net/http"
	"sort"
	"strings"

	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"
	"icmongolang/internal/modules/settings/presenter"
)

var scheduleCols = []string{
	"schedule_name", "device_id", "start", "event",
	"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
	"status",
}

var scheduleDayCols = []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

// ListScheduleDevice - GET /api/settings/listscheduledevice
// @Summary List Schedule Device
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
// @Router /settings/listscheduledevice [get]
func (h *settingsHandler) ListScheduleDevice() http.HandlerFunc {
	return h.listHandler(settings.ScheduleDeviceJoinSpec)
}

// FindScheduleDeviceChk - GET /api/settings/findscheduledevicechk
// @Summary Find Schedule Device Chk
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
// @Router /settings/findscheduledevicechk [get]
func (h *settingsHandler) FindScheduleDeviceChk() http.HandlerFunc {
	return h.listHandler(settings.ScheduleDeviceJoinSpec)
}

// ScheduleList - GET /api/settings/schedulelist
// @Summary Schedule List
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
// @Router /settings/schedulelist [get]
func (h *settingsHandler) ScheduleList() http.HandlerFunc {
	return h.listHandler(settings.ScheduleSpec)
}

// ScheduleAll - GET /api/settings/scheduleall
// @Summary Schedule All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/scheduleall [get]
func (h *settingsHandler) ScheduleAll() http.HandlerFunc { return h.allHandler(settings.ScheduleSpec) }

// ListSchedulePage - GET /api/settings/listschedulepage
// @Summary List Schedule Page
// @Description Paginated listing (filters via query params, sort=field-ASC|DESC).
// Merges sd_iot_schedule (ScheduleSpec) and fs_schedule (FullScheduleSpec) so
// the new Full Schedule module shows up in the same page.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Items per page (default 1000, max 5000)"
// @Param sort query string false "Sort: field-ASC or field-DESC"
// @Param keyword query string false "Keyword filter (LIKE)"
// @Param event query string false "Event filter (exact match on sc.event)"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult] "Paginated result envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/listschedulepage [get]
func (h *settingsHandler) ListSchedulePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := listQueryFrom(r)
		ctx := r.Context()

		sdRows, err := h.uc.RowsAll(ctx, settings.ScheduleSpec, noPageQuery(q))
		if err != nil {
			fail(w, r, err)
			return
		}
		fsRows, err := h.uc.RowsAll(ctx, settings.FullScheduleSpec, noPageQuery(q))
		if err != nil {
			fail(w, r, err)
			return
		}

		rows := make([]map[string]interface{}, 0, len(sdRows)+len(fsRows))
		rows = append(rows, sdRows...)
		for _, row := range fsRows {
			rows = append(rows, normalizeFullScheduleRow(row))
		}

		convertDatesToBangkok(rows)
		rows = sortScheduleRows(rows, q.Sort)
		total := int64(len(rows))

		page := q.Page
		if page < 1 {
			page = 1
		}
		pageSize := q.PageSize
		if pageSize < 1 || pageSize > 5000 {
			pageSize = 1000
		}
		start := (page - 1) * pageSize
		switch {
		case start >= len(rows):
			rows = []map[string]interface{}{}
		case start+pageSize < len(rows):
			rows = rows[start : start+pageSize]
		default:
			rows = rows[start:]
		}

		ok(w, r, presenter.NewListResult(rows, total, page, pageSize, q.FilterMap()))
	}
}

// noPageQuery returns a copy of q with paging disabled (RowsAll then runs the
// full unfiltered query which the handler re-slices in memory).
// ส่งคืน ListQuery ที่ไม่แบ่งหน้า (ให้ RowsAll ดึงทั้งหมดแล้วแบ่งหน้าใน Go)
func noPageQuery(q settings.ListQuery) settings.ListQuery {
	q.PageSize = 0
	return q
}

// normalizeFullScheduleRow maps fs_schedule rows into the sd_iot_schedule
// display shape — status int8 (1=active, 0=inactive, 3=draft) emitted as a
// number, consistent with the sd_iot_schedule side.
// ปรับแถว fs_schedule ให้ตรงกับรูปแบบของ sd_iot_schedule (status ตัวเลข 1/0/3)
func normalizeFullScheduleRow(row map[string]interface{}) map[string]interface{} {
	switch v := row["status"].(type) {
	case string:
		switch v {
		case "active":
			row["status"] = int64(1)
		case "inactive":
			row["status"] = int64(0)
		case "draft":
			row["status"] = int64(3)
		}
	case int64:
		// already numeric (1/0/3) — keep as-is
	}
	return row
}

// sortScheduleRows sorts the merged rows by q.Sort (default "start ASC") using
// the settings-visible column names. Dates already carry "2006-01-02 15:04:05"
// so plain string comparison is stable across both sources.
// เรียงแถวที่รวมแล้วตาม sort (ค่าเริ่มต้น start ASC)
func sortScheduleRows(rows []map[string]interface{}, sortParam string) []map[string]interface{} {
	col, ok, desc := splitScheduleSort(sortParam)
	if !ok {
		col, desc = "start", false
	}
	sort.SliceStable(rows, func(i, j int) bool {
		a := toStr(rows[i][col])
		b := toStr(rows[j][col])
		if desc {
			return a > b
		}
		return a < b
	})
	return rows
}

// splitScheduleSort parses "field-ASC|DESC" for the merged schedule page.
func splitScheduleSort(sortParam string) (col string, ok bool, desc bool) {
	parts := strings.SplitN(sortParam, "-", 2)
	if len(parts) != 2 {
		return "", false, false
	}
	switch parts[0] {
	case "schedule_id", "schedule_name", "start", "createddate", "updateddate":
	default:
		return "", false, false
	}
	switch strings.ToUpper(parts[1]) {
	case "ASC":
		return parts[0], true, false
	case "DESC":
		return parts[0], true, true
	}
	return "", false, false
}

// ScheduleDevicePage - GET /api/settings/scheduledevicepage
// @Summary Schedule Device Page
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
// @Router /settings/scheduledevicepage [get]
func (h *settingsHandler) ScheduleDevicePage() http.HandlerFunc {
	return h.listHandler(settings.ScheduleDeviceJoinSpec)
}

// ListDeviceScheduleData ports /listdevicescheduledata -> get_data_schedule_device.
// ListDeviceScheduleData - GET /api/settings/listdevicescheduledata
// @Summary List Device Schedule Data
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
// @Router /settings/listdevicescheduledata [get]
func (h *settingsHandler) ListDeviceScheduleData() http.HandlerFunc {
	return h.listHandler(settings.ScheduleDeviceJoinSpec)
}

// CreateScheduleDeviceViaGet ports GET /createscheduledevice.
// CreateScheduleDeviceViaGet - GET /api/settings/createscheduledevice
// @Summary Create Schedule Device Via Get
// @Description Creates a link row from required query params.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param schedule_id query string true "schedule_id"
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createscheduledevice [get]
func (h *settingsHandler) CreateScheduleDeviceViaGet() http.HandlerFunc {
	return h.insertLinkViaGet("sd_iot_schedule_device", "schedule_id", "device_id")
}

// CreateScheduleDevice ports POST /createscheduledevice.
// CreateScheduleDevice - POST /api/settings/createscheduledevice
// @Summary Create Schedule Device
// @Description Links a schedule to a device (schedule_id + device_id required).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/createscheduledevice [post]
func (h *settingsHandler) CreateScheduleDevice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		row := filterKeys(body, []string{"schedule_id", "device_id"})
		if row["schedule_id"] == nil || row["device_id"] == nil {
			badRequest(w, r, "schedule_id and device_id are required")
			return
		}
		if err = h.uc.CreateRow(r.Context(), "sd_iot_schedule_device", row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}

// DeleteScheduleDevices ports GET /deletescheduledevice (by schedule_id + device_id).
// DeleteScheduleDevices - GET /api/settings/deletescheduledevice
// @Summary Delete Schedule Devices
// @Description Deletes one schedule-device link by schedule_id + device_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param schedule_id query string true "Schedule ID"
// @Param device_id query string true "Device ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletescheduledevice [get]
func (h *settingsHandler) DeleteScheduleDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		scheduleID := q.Get("schedule_id")
		deviceID := q.Get("device_id")
		if scheduleID == "" || deviceID == "" {
			badRequest(w, r, "schedule_id and device_id are required")
			return
		}
		n, err := h.uc.DeleteRow(r.Context(), "sd_iot_schedule_device",
			"schedule_id = ? AND device_id = ?", coerceID(scheduleID), coerceID(deviceID))
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// DeleteDeviceSchedule ports GET /deletedeviceschedule (by device_id [+schedule_id]).
// DeleteDeviceSchedule - GET /api/settings/deletedeviceschedule
// @Summary Delete Device Schedule
// @Description Deletes schedule-device links by device_id (optionally filtered by schedule_id).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "Device ID"
// @Param schedule_id query string false "Schedule ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletedeviceschedule [get]
func (h *settingsHandler) DeleteDeviceSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("device_id")
		if deviceID == "" {
			badRequest(w, r, "device_id is required")
			return
		}
		query := "device_id = ?"
		args := []interface{}{coerceID(deviceID)}
		if sid := r.URL.Query().Get("schedule_id"); sid != "" {
			query += " AND schedule_id = ?"
			args = append(args, coerceID(sid))
		}
		n, err := h.uc.DeleteRow(r.Context(), "sd_iot_schedule_device", query, args...)
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// DeleteDeviceAndSchedule ports GET /deletedeviceandschedule: removes the
// schedule and its device links in one call.
// DeleteDeviceAndSchedule - GET /api/settings/deletedeviceandschedule
// @Summary Delete Device And Schedule
// @Description Deletes a schedule and/or its device links (schedule_id and/or device_id).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param schedule_id query string false "Schedule ID"
// @Param device_id query string false "Device ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletedeviceandschedule [get]
func (h *settingsHandler) DeleteDeviceAndSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid := r.URL.Query().Get("schedule_id")
		did := r.URL.Query().Get("device_id")
		if sid == "" && did == "" {
			badRequest(w, r, "schedule_id or device_id is required")
			return
		}
		total := int64(0)
		if sid != "" {
			n, err := h.uc.DeleteRow(r.Context(), "sd_iot_schedule", "schedule_id = ?", coerceID(sid))
			if err != nil {
				fail(w, r, err)
				return
			}
			total += n
		}
		if did != "" {
			n, err := h.uc.DeleteRow(r.Context(), "sd_iot_schedule_device", "device_id = ?", coerceID(did))
			if err != nil {
				fail(w, r, err)
				return
			}
			total += n
		}
		affected(w, r, total)
	}
}

// CreateSchedule - POST /api/settings/createschedule
// @Summary Create Schedule
// @Description Creates one row in sd_iot_schedule (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createschedule [post]
func (h *settingsHandler) CreateSchedule() http.HandlerFunc {
	return h.createHandler("sd_iot_schedule", scheduleCols, statusDefault,
		dupCheck{BodyKey: "schedule_name", Column: "schedule_name", Field: "schedule_name"})
}

// UpdateSchedule - POST /api/settings/updateschedule
// @Summary Update Schedule
// @Description Partial update of sd_iot_schedule by schedule_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include schedule_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateschedule [post]
func (h *settingsHandler) UpdateSchedule() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_schedule", "schedule_id", scheduleCols)
}

// UpdateScheduleStatus - POST /api/settings/updateschedulestatus
// @Summary Update Schedule Status
// @Description Partial update of sd_iot_schedule by schedule_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include schedule_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateschedulestatus [post]
func (h *settingsHandler) UpdateScheduleStatus() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_schedule", "schedule_id", []string{"status"})
}

// UpdateScheduleDayStatus ports POST /updatescheduledaystatus: toggles one day column.
// UpdateScheduleDayStatus - POST /api/settings/updatescheduledaystatus
// @Summary Update Schedule Day Status
// @Description Toggles day columns (sunday..saturday) of one schedule. Body requires schedule_id plus at least one day field.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatescheduledaystatus [post]
func (h *settingsHandler) UpdateScheduleDayStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		id := toStr(body["schedule_id"])
		if id == "" {
			badRequest(w, r, "schedule_id is required")
			return
		}
		fields := map[string]interface{}{"updateddate": nowPtr()}
		for _, day := range scheduleDayCols {
			if v, exists := body[day]; exists && v != nil {
				fields[day] = v
			}
		}
		if len(fields) == 1 {
			badRequest(w, r, "no day field provided")
			return
		}
		n, err := h.uc.UpdateFields(r.Context(), "sd_iot_schedule", "schedule_id", coerceID(id), fields)
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// DeleteSchedule - GET /api/settings/deleteschedule
// @Summary Delete Schedule
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param schedule_id query string true "schedule_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deleteschedule [get]
func (h *settingsHandler) DeleteSchedule() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_schedule", "schedule_id", "schedule_id")
}
