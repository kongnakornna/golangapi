package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"
)

// ---- device alarm views / monitors ----
// The NestJS versions join live MQTT payloads; the Go ports paginate over the
// same underlying tables with the documented filters (see plan deviations).

// ListDeviceAlarm - GET /api/settings/listdevicealarm
// @Summary List Device Alarm
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
// @Router /settings/listdevicealarm [get]
func (h *settingsHandler) ListDeviceAlarm() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ListDeviceAlarmAirV1 - GET /api/settings/listdevicealarmairV1
// @Summary List Device Alarm Air V1
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
// @Router /settings/listdevicealarmairV1 [get]
func (h *settingsHandler) ListDeviceAlarmAirV1() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ListDeviceAlarmAir - GET /api/settings/listdevicealarmair
// @Summary List Device Alarm Air
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
// @Router /settings/listdevicealarmair [get]
func (h *settingsHandler) ListDeviceAlarmAir() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ListDeviceAlarmAll - GET /api/settings/listdevicealarmall
// @Summary List Device Alarm All
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
// @Router /settings/listdevicealarmall [get]
func (h *settingsHandler) ListDeviceAlarmAll() http.HandlerFunc {
	return h.listHandler(settings.DeviceActiveSpec)
}

var fanFixed = settings.FixedFilter{Column: "d.hardware_id", Value: "3"}

// UnderscoreListDeviceAlarmFan - GET /api/settings/_listdevicealarmfan
// @Summary Underscore List Device Alarm Fan
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
// @Router /settings/_listdevicealarmfan [get]
func (h *settingsHandler) UnderscoreListDeviceAlarmFan() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec.With(fanFixed))
}

// ListDeviceAlarmFan - GET /api/settings/listdevicealarmfan
// @Summary List Device Alarm Fan
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
// @Router /settings/listdevicealarmfan [get]
func (h *settingsHandler) ListDeviceAlarmFan() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec.With(fanFixed))
}

// ListDeviceAlarmLimit - GET /api/settings/listdevicealarmlimit
// @Summary List Device Alarm Limit
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
// @Router /settings/listdevicealarmlimit [get]
func (h *settingsHandler) ListDeviceAlarmLimit() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec.With(
		settings.FixedFilter{Column: "d.alert_set", Value: "1"}))
}

// UnderscoreDeviceMonitor - GET /api/settings/_devicemonitor
// @Summary Underscore Device Monitor
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
// @Router /settings/_devicemonitor [get]
func (h *settingsHandler) UnderscoreDeviceMonitor() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// DeviceMonitor - GET /api/settings/devicemonitor
// @Summary Device Monitor
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
// @Router /settings/devicemonitor [get]
func (h *settingsHandler) DeviceMonitor() http.HandlerFunc { return h.listHandler(settings.DeviceSpec) }

// DeviceMonitors - GET /api/settings/devicemonitors
// @Summary Device Monitors
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
// @Router /settings/devicemonitors [get]
func (h *settingsHandler) DeviceMonitors() http.HandlerFunc {
	return h.listHandler(settings.DeviceSpec)
}

// ---- process / alarm logs ----

// ScheduleProces returns recent schedule process log rows.
// ScheduleProces - GET /api/settings/scheduleproces
// @Summary Schedule Proces
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
// @Router /settings/scheduleproces [get]
func (h *settingsHandler) ScheduleProces() http.HandlerFunc {
	return h.listHandler(settings.ScheduleProcessLogSpec)
}

// ScheduleProcessLog - GET /api/settings/scheduleprocesslog
// @Summary Schedule Process Log
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
// @Router /settings/scheduleprocesslog [get]
func (h *settingsHandler) ScheduleProcessLog() http.HandlerFunc {
	return h.listHandler(settings.ScheduleProcessLogSpec)
}

// ScheduleProcessLogPaginate - GET /api/settings/scheduleprocesslogpaginate
// @Summary Schedule Process Log Paginate
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
// @Router /settings/scheduleprocesslogpaginate [get]
func (h *settingsHandler) ScheduleProcessLogPaginate() http.HandlerFunc {
	return h.listHandler(settings.ScheduleProcessLogSpec)
}

// MqttErrorLogPaginate - GET /api/settings/mqtterrorlogpaginate
// @Summary Mqtt Error Log Paginate
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
// @Router /settings/mqtterrorlogpaginate [get]
func (h *settingsHandler) MqttErrorLogPaginate() http.HandlerFunc {
	return h.listHandler(settings.MqttErrorLogSpec)
}

// AlarmLogPaginate - GET /api/settings/alarmlogpaginate
// @Summary Alarm Log Paginate
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
// @Router /settings/alarmlogpaginate [get]
func (h *settingsHandler) AlarmLogPaginate() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogSpec)
}

// AlarmLogPaginateEmail - GET /api/settings/alarmlogpaginateemail
// @Summary Alarm Log Paginate Email
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
// @Router /settings/alarmlogpaginateemail [get]
func (h *settingsHandler) AlarmLogPaginateEmail() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogEmailSpec)
}

// AlarmLogPaginateLine - GET /api/settings/alarmlogpaginateline
// @Summary Alarm Log Paginate Line
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
// @Router /settings/alarmlogpaginateline [get]
func (h *settingsHandler) AlarmLogPaginateLine() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogLineSpec)
}

// AlarmLogPaginateSms - GET /api/settings/alarmlogpaginatesms
// @Summary Alarm Log Paginate Sms
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
// @Router /settings/alarmlogpaginatesms [get]
func (h *settingsHandler) AlarmLogPaginateSms() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogSmsSpec)
}

// AlarmLogPaginateTelegram - GET /api/settings/alarmlogpaginatetelegram
// @Summary Alarm Log Paginate Telegram
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
// @Router /settings/alarmlogpaginatetelegram [get]
func (h *settingsHandler) AlarmLogPaginateTelegram() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogTelegramSpec)
}

// AlarmLogPaginateControls - GET /api/settings/alarmlogpaginatecontrols
// @Summary Alarm Log Paginate Controls
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
// @Router /settings/alarmlogpaginatecontrols [get]
func (h *settingsHandler) AlarmLogPaginateControls() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogControlSpec)
}

// AlarmLogPaginateControl - GET /api/settings/alarmlogpaginatecontrol
// @Summary Alarm Log Paginate Control
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
// @Router /settings/alarmlogpaginatecontrol [get]
func (h *settingsHandler) AlarmLogPaginateControl() http.HandlerFunc {
	return h.listHandler(settings.AlarmProcessLogControlSpec)
}
