package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"
)

var alarmActionCols = []string{
	"action_name", "status_warning", "recovery_warning", "status_alert",
	"recovery_alert", "email_alarm", "line_alarm", "telegram_alarm", "sms_alarm",
	"nonc_alarm", "time_life", "event", "status",
}

// ---- alarm device / event links (sd_iot_alarm_device / _event) ----

// ListAlarmDevicePage - GET /api/settings/listalarmdevicepage
// @Summary List Alarm Device Page
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
// @Router /settings/listalarmdevicepage [get]
func (h *settingsHandler) ListAlarmDevicePage() http.HandlerFunc {
	return h.listHandler(settings.AlarmDeviceJoinSpec)
}

// ListAlarmEventDevicePage - GET /api/settings/listalarmeventdevicepage
// @Summary List Alarm Event Device Page
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
// @Router /settings/listalarmeventdevicepage [get]
func (h *settingsHandler) ListAlarmEventDevicePage() http.HandlerFunc {
	return h.listHandler(settings.AlarmEventDeviceJoinSpec)
}

// ListAlarmDeviceActivePage - GET /api/settings/listalarmdeviceactivepage
// @Summary List Alarm Device Active Page
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
// @Router /settings/listalarmdeviceactivepage [get]
func (h *settingsHandler) ListAlarmDeviceActivePage() http.HandlerFunc {
	return h.listHandler(settings.AlarmDeviceJoinSpec.With(settings.FixedFilter{Column: "aa.status", Value: "1"}))
}

// ListAlarmEventDeviceControlPage - GET /api/settings/listalarmeventdevicecontrolpage
// @Summary List Alarm Event Device Control Page
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
// @Router /settings/listalarmeventdevicecontrolpage [get]
func (h *settingsHandler) ListAlarmEventDeviceControlPage() http.HandlerFunc {
	return h.listHandler(settings.AlarmEventDeviceJoinSpec)
}

// ActiveAlarmDevicePage - GET /api/settings/activealarmdevicepage
// @Summary Active Alarm Device Page
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
// @Router /settings/activealarmdevicepage [get]
func (h *settingsHandler) ActiveAlarmDevicePage() http.HandlerFunc {
	return h.listHandler(settings.DeviceActiveSpec)
}

// ActiveAlarmEventDeviceEventPage - GET /api/settings/activealarmeventdeviceeventpage
// @Summary Active Alarm Event Device Event Page
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
// @Router /settings/activealarmeventdeviceeventpage [get]
func (h *settingsHandler) ActiveAlarmEventDeviceEventPage() http.HandlerFunc {
	return h.listHandler(settings.DeviceActiveSpec)
}

// DeviceActiveMqttAlarm - GET /api/settings/deviceactivemqttalarm
// @Summary Device Active Mqtt Alarm
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
// @Router /settings/deviceactivemqttalarm [get]
func (h *settingsHandler) DeviceActiveMqttAlarm() http.HandlerFunc {
	return h.listHandler(settings.DeviceActiveSpec)
}

// AlarmDevice - GET /api/settings/alarmdevice
// @Summary Alarm Device
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
// @Router /settings/alarmdevice [get]
func (h *settingsHandler) AlarmDevice() http.HandlerFunc {
	return h.listHandler(settings.AlarmDeviceJoinSpec)
}

// DeviceAlarm - GET /api/settings/devicealarm
// @Summary Device Alarm
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
// @Router /settings/devicealarm [get]
func (h *settingsHandler) DeviceAlarm() http.HandlerFunc { return h.listHandler(settings.DeviceSpec) }

// AlarmDeviceStatus ports /alarmdevicestatus: work_status summary of one device.
// AlarmDeviceStatus - GET /api/settings/alarmdevicestatus
// @Summary Alarm Device Status
// @Description Returns work/alarm status snapshot of one device.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "Device ID"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/alarmdevicestatus [get]
func (h *settingsHandler) AlarmDeviceStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("device_id")
		if id == "" {
			badRequest(w, r, "device_id is required")
			return
		}
		items, err := h.uc.RowsAll(r.Context(),
			settings.DeviceSpec.With(settings.FixedFilter{Column: "d.device_id", Value: id}),
			settings.NewListQuery(1, 1, "", nil))
		if err != nil {
			fail(w, r, err)
			return
		}
		if len(items) == 0 {
			notFound(w, r, "device not found")
			return
		}
		ok(w, r, items[0])
	}
}

// CreateAlarmDeviceViaGet - GET /api/settings/createalarmdevice
// @Summary Create Alarm Device Via Get
// @Description Creates a link row from required query params.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param alarm_action_id query string true "alarm_action_id"
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createalarmdevice [get]
func (h *settingsHandler) CreateAlarmDeviceViaGet() http.HandlerFunc {
	return h.insertLinkViaGet("sd_iot_alarm_device", "alarm_action_id", "device_id")
}

// DeleteAlarmDevices - GET /api/settings/deletealarmdevice
// @Summary Delete Alarm Devices
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletealarmdevice [get]
func (h *settingsHandler) DeleteAlarmDevices() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_alarm_device", "device_id", "device_id")
}

// CreateAlarmEventDeviceViaGet - GET /api/settings/createalarmeventdevice
// @Summary Create Alarm Event Device Via Get
// @Description Creates a link row from required query params.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param alarm_action_id query string true "alarm_action_id"
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createalarmeventdevice [get]
func (h *settingsHandler) CreateAlarmEventDeviceViaGet() http.HandlerFunc {
	return h.insertLinkViaGet("sd_iot_alarm_device_event", "alarm_action_id", "device_id")
}

// DeleteAlarmEventDevices - GET /api/settings/deletealarmeventdevice
// @Summary Delete Alarm Event Devices
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param device_id query string true "device_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletealarmeventdevice [get]
func (h *settingsHandler) DeleteAlarmEventDevices() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_alarm_device_event", "device_id", "device_id")
}

// DeleteArmDevice ports GET /deletearmdevice (by alarm_action_id + device_id).
// DeleteArmDevice - GET /api/settings/deletearmdevice
// @Summary Delete Arm Device
// @Description Deletes an alarm-device link by alarm_action_id + device_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param alarm_action_id query string true "Alarm action ID"
// @Param device_id query string true "Device ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletearmdevice [get]
func (h *settingsHandler) DeleteArmDevice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		actionID := q.Get("alarm_action_id")
		deviceID := q.Get("device_id")
		if actionID == "" || deviceID == "" {
			badRequest(w, r, "alarm_action_id and device_id are required")
			return
		}
		n, err := h.uc.DeleteRow(r.Context(), "sd_iot_alarm_device",
			"alarm_action_id = ? AND device_id = ?", coerceID(actionID), coerceID(deviceID))
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// DeleteArmDeviceV2 ports GET /deletearmdevicev2: same as v1 on the link table.
// DeleteArmDeviceV2 - GET /api/settings/deletearmdevicev2
// @Summary Delete Arm Device V2
// @Description Deletes an alarm-device link by alarm_action_id + device_id (v2 route alias).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param alarm_action_id query string true "Alarm action ID"
// @Param device_id query string true "Device ID"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletearmdevicev2 [get]
func (h *settingsHandler) DeleteArmDeviceV2() http.HandlerFunc { return h.DeleteArmDevice() }

// CreateAlarmDevice ports POST /createalarmDevice -> create_alarm_device.
// CreateAlarmDevice - POST /api/settings/createalarmDevice
// @Summary Create Alarm Device
// @Description Creates one row in sd_iot_device_alarm_action (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createalarmDevice [post]
func (h *settingsHandler) CreateAlarmDevice() http.HandlerFunc {
	return h.createHandler("sd_iot_device_alarm_action", alarmActionCols, statusDefault)
}

// UpdateAlarmDevice - POST /api/settings/updatealarmdevice
// @Summary Update Alarm Device
// @Description Partial update of sd_iot_device_alarm_action by alarm_action_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include alarm_action_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatealarmdevice [post]
func (h *settingsHandler) UpdateAlarmDevice() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_device_alarm_action", "alarm_action_id", alarmActionCols)
}

// UpdateAlarmStatus - POST /api/settings/updatealarmstatus
// @Summary Update Alarm Status
// @Description Partial update of sd_iot_device_alarm_action by alarm_action_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include alarm_action_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatealarmstatus [post]
func (h *settingsHandler) UpdateAlarmStatus() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_device_alarm_action", "alarm_action_id", []string{"status"})
}

// CreateDeviceAlarmAction ports POST /createdevicealarmaction.
// CreateDeviceAlarmAction - POST /api/settings/createdevicealarmaction
// @Summary Create Device Alarm Action
// @Description Creates a device-action-user link in sd_iot_device_action_user (fields: alarm_action_id, uid).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/createdevicealarmaction [post]
func (h *settingsHandler) CreateDeviceAlarmAction() http.HandlerFunc {
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
