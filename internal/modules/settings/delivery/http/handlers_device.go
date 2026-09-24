package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"
)

var deviceCols = []string{
	"setting_id", "type_id", "location_id", "device_name", "sn", "hardware_id",
	"status_warning", "recovery_warning", "status_alert", "recovery_alert",
	"time_life", "period", "work_status", "model", "vendor", "comparevalue",
	"unit", "mqtt_id", "oid", "action_id", "status_alert_id", "measurement",
	"mqtt_data_value", "mqtt_data_control", "mqtt_control_on", "mqtt_control_off",
	"org", "bucket", "status", "mqtt_device_name", "mqtt_status_over_name",
	"mqtt_status_data_name", "mqtt_act_relay_name", "mqtt_control_relay_name",
	"mqtt_config", "max", "min", "layout", "alert_set", "menu",
	"calibration_add", "calibration_subtract", "calibration_type",
}

// ListDeviceAll - GET /api/settings/deviceall
// @Summary List Device All
// @Description Returns every sd_iot_device row unpaginated (ports NestJS device_all: SELECT d.*).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "All device rows"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/deviceall [get]
func (h *settingsHandler) ListDeviceAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := h.uc.RowsAll(r.Context(), settings.ListSpec{
			Table:       "sd_iot_device d",
			Selects:     []string{"d.*"},
			DefaultSort: "d.device_id ASC",
		}, settings.NewListQuery(0, 0, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, items)
	}
}

// ListDevicePage - GET /api/settings/listdevicepage
// @Summary List Device Page
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
// @Router /settings/listdevicepage [get]
func (h *settingsHandler) ListDevicePage() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ListDevicePagess - GET /api/settings/listdevicepagess
// @Summary List Device Pagess
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
// @Router /settings/listdevicepagess [get]
func (h *settingsHandler) ListDevicePagess() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ListDevicePageActive1 - GET /api/settings/listdevicepageactive1
// @Summary List Device Page Active1
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
// @Router /settings/listdevicepageactive1 [get]
func (h *settingsHandler) ListDevicePageActive1() http.HandlerFunc {
	return h.listHandler(settings.DeviceActive1Spec)
}

// ListDevicePageActive - GET /api/settings/listdevicepageactive
// @Summary List Device Page Active
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
// @Router /settings/listdevicepageactive [get]
func (h *settingsHandler) ListDevicePageActive() http.HandlerFunc {
	return h.listHandler(settings.DeviceActiveSpec)
}

// ListDevicePageAll - GET /api/settings/listdevicepageall
// @Summary List Device Page All
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
// @Router /settings/listdevicepageall [get]
func (h *settingsHandler) ListDevicePageAll() http.HandlerFunc {
	return h.listHandler(settings.DeviceAllSpec)
}

// ListDevicePageSensor - GET /api/settings/listdevicepagesensor
// @Summary List Device Page Sensor
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
// @Router /settings/listdevicepagesensor [get]
func (h *settingsHandler) ListDevicePageSensor() http.HandlerFunc {
	return h.listHandler(settings.DeviceAllSpec)
}

// ListDevicePageAllActive - GET /api/settings/listdevicepageallactive
// @Summary List Device Page All Active
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
// @Router /settings/listdevicepageallactive [get]
func (h *settingsHandler) ListDevicePageAllActive() http.HandlerFunc {
	return h.listHandler(settings.DeviceAllActiveSpec)
}

// ListDevicePageAllActiveSchedule - GET /api/settings/listdevicepageallactiveschedule
// @Summary List Device Page All Active Schedule
// @Description Active devices annotated with their schedule membership. Requires schedule_id;
// hardware_id defaults to 3; each row gains schedule_status / count_schedule_device /
// schedule_name / schedule_start / schedule_title.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param schedule_id query string true "schedule_id"
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Items per page (default 1000, max 5000)"
// @Param sort query string false "Sort: field-ASC or field-DESC"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult] "Paginated result envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/listdevicepageallactiveschedule [get]
func (h *settingsHandler) ListDevicePageAllActiveSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lq := listQueryFrom(r)
		scheduleID := lq.Value("schedule_id")
		if scheduleID == "" {
			badRequest(w, r, "schedule_id is null.")
			return
		}
		schedRows, err := h.uc.RowsAll(r.Context(),
			settings.ScheduleSpec.With(settings.FixedFilter{Column: "sc.schedule_id", Value: scheduleID}),
			settings.NewListQuery(1, 1, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		if len(schedRows) == 0 {
			notFound(w, r, "Data not found.")
			return
		}
		sched := schedRows[0]
		// NestJS defaults this view to hardware_id=3 (IO Control devices).
		spec := settings.DeviceAllActiveSchedSpec
		if lq.Value("hardware_id") == "" {
			spec = spec.With(settings.FixedFilter{Column: "d.hardware_id", Value: "3"})
		}
		res, err := h.uc.ListPaginate(r.Context(), spec, lq)
		if err != nil {
			fail(w, r, err)
			return
		}
		for _, row := range res.Data {
			cnt, err := h.uc.CountWhere(r.Context(), "sd_iot_schedule_device",
				"schedule_id = ? AND device_id = ?", coerceID(scheduleID), coerceID(toStr(row["device_id"])))
			if err != nil {
				fail(w, r, err)
				return
			}
			scheduleStatus := int64(0)
			if cnt >= 1 {
				scheduleStatus = 1
			}
			name := toStr(sched["schedule_name"])
			start := toStr(sched["start"])
			row["schedule_id"] = scheduleID
			row["schedule_status"] = scheduleStatus
			row["count_schedule_device"] = cnt
			row["schedule_name"] = name
			row["schedule_start"] = start
			row["schedule_title"] = name + " " + start
		}
		ok(w, r, res)
	}
}

// DeviceEditGet - GET /api/settings/deviceeditget
// @Summary Device Edit Get
// @Description Fetches one record by device_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/deviceeditget [get]
func (h *settingsHandler) DeviceEditGet() http.HandlerFunc {
	return h.getOneByParam(settings.DeviceSpec, "device_id", "d.device_id", true, "device not found")
}

// DeviceDetail - GET /api/settings/devicedetail
// @Summary Device Detail
// @Description Fetches one record by device_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/devicedetail [get]
func (h *settingsHandler) DeviceDetail() http.HandlerFunc {
	return h.getOneByParam(settings.DeviceSpec, "device_id", "d.device_id", true, "device not found")
}

// DeviceDeleteCheck ports /devicedelete: reports whether the device exists.
// DeviceDeleteCheck - GET /api/settings/devicedelete
// @Summary Device Delete Check
// @Description Checks whether a device exists before deletion.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "Device ID"
// @Success 200 {object} responses.SwaggerSuccessResponse "{exists: bool}"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/devicedelete [get]
func (h *settingsHandler) DeviceDeleteCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("device_id")
		if id == "" {
			badRequest(w, r, "device_id is required")
			return
		}
		exists, err := h.uc.ExistsWhere(r.Context(), "sd_iot_device", "device_id = ?", coerceID(id))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, map[string]interface{}{"exists": exists})
	}
}

// DeleteDevice - GET /api/settings/deletedevice
// @Summary Delete Device
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletedevice [get]
func (h *settingsHandler) DeleteDevice() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_device", "device_id", "device_id")
}

// CreateDevice - POST /api/settings/createdevice
// @Summary Create Device
// @Description Creates one row in sd_iot_device (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createdevice [post]
func (h *settingsHandler) CreateDevice() http.HandlerFunc {
	return h.createHandler("sd_iot_device", deviceCols, map[string]interface{}{
		"status": 1, "work_status": 1,
	}, dupCheck{BodyKey: "sn", Column: "sn", Field: "SN"})
}

// UpdateDevice - POST /api/settings/updatedevice
// @Summary Update Device
// @Description Partial update of sd_iot_device by device_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include device_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatedevice [post]
func (h *settingsHandler) UpdateDevice() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_device", "device_id", deviceCols)
}

// UpdateStatusDeviceId ports POST /updatestatusdeviceid.
// UpdateStatusDeviceId - POST /api/settings/updatestatusdeviceid
// @Summary Update Status Device Id
// @Description Updates status and/or work_status of one device. Body requires device_id.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatestatusdeviceid [post]
func (h *settingsHandler) UpdateStatusDeviceId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		id := toStr(body["device_id"])
		if id == "" {
			badRequest(w, r, "device_id is required")
			return
		}
		fields := map[string]interface{}{"updateddate": nowPtr()}
		if v, exists := body["status"]; exists && v != nil {
			fields["status"] = v
		}
		if v, exists := body["work_status"]; exists && v != nil {
			fields["work_status"] = v
		}
		n, err := h.uc.UpdateFields(r.Context(), "sd_iot_device", "device_id", coerceID(id), fields)
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// DeviceActionUser ports POST /deviceactionuser -> create_deviceactionuser.
// DeviceActionUser - POST /api/settings/deviceactionuser
// @Summary Device Action User
// @Description Creates a device-action-user link (fields: alarm_action_id, uid).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deviceactionuser [post]
func (h *settingsHandler) DeviceActionUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		row := filterKeys(body, []string{"alarm_action_id", "uid"})
		if len(row) == 0 {
			badRequest(w, r, "request body has no valid fields")
			return
		}
		if err = h.uc.CreateRow(r.Context(), "sd_iot_device_action_user", row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}
