package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"
)

// ---- column whitelists (writable columns per table, from internal/models) ----

var settingCols = []string{"location_id", "setting_type_id", "setting_name", "sn", "status"}
var locationCols = []string{"location_name", "ipaddress", "location_detail", "configdata", "status"}
var typeCols = []string{"type_name", "group_id", "status"}
var deviceTypeCols = []string{"type_name", "status"}
var groupCols = []string{"group_name", "status"}
var sensorCols = []string{
	"setting_id", "setting_type_id", "sensor_name", "sn", "max", "min",
	"hardware_id", "status_high", "status_warning", "status_alert",
	"model", "vendor", "comparevalue", "unit", "mqtt_id", "oid",
	"action_id", "status_alert_id", "mqtt_data_value", "mqtt_data_control", "status",
}

var statusDefault = map[string]interface{}{"status": 1}

// ---- setting ----

// ListSetting - GET /api/settings/listsetting
// @Summary List Setting
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
// @Router /settings/listsetting [get]
func (h *settingsHandler) ListSetting() http.HandlerFunc { return h.listHandler(settings.SettingSpec) }

// SettingAll - GET /api/settings/settingall
// @Summary Setting All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/settingall [get]
func (h *settingsHandler) SettingAll() http.HandlerFunc { return h.allHandler(settings.SettingSpec) }

// CreateSetting - POST /api/settings/createsetting
// @Summary Create Setting
// @Description Creates one row in sd_iot_setting (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createsetting [post]
func (h *settingsHandler) CreateSetting() http.HandlerFunc {
	return h.createHandler("sd_iot_setting", settingCols, statusDefault,
		dupCheck{BodyKey: "sn", Column: "sn", Field: "SN"})
}

// UpdateSetting - POST /api/settings/updatesetting
// @Summary Update Setting
// @Description Partial update of sd_iot_setting by setting_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include setting_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatesetting [post]
func (h *settingsHandler) UpdateSetting() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_setting", "setting_id", settingCols)
}

// DeleteSettingViaGet - GET /api/settings/deletesetting
// @Summary Delete Setting Via Get
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param setting_id query string true "setting_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletesetting [get]
func (h *settingsHandler) DeleteSettingViaGet() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_setting", "setting_id", "setting_id")
}

// DeleteSetting - DELETE /api/settings/deletesetting
// @Summary Delete Setting
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param setting_id query string true "setting_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletesetting [delete]
func (h *settingsHandler) DeleteSetting() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_setting", "setting_id", "setting_id")
}

// ---- location ----

// ListLocation - GET /api/settings/listlocation
// @Summary List Location
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
// @Router /settings/listlocation [get]
func (h *settingsHandler) ListLocation() http.HandlerFunc {
	return h.listHandler(settings.LocationSpec)
}

// LocationAll - GET /api/settings/locationall
// @Summary Location All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/locationall [get]
func (h *settingsHandler) LocationAll() http.HandlerFunc { return h.allHandler(settings.LocationSpec) }

// CreateLocation - POST /api/settings/createlocation
// @Summary Create Location
// @Description Creates one row in sd_iot_location (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createlocation [post]
func (h *settingsHandler) CreateLocation() http.HandlerFunc {
	return h.createHandler("sd_iot_location", locationCols, statusDefault,
		dupCheck{BodyKey: "ipaddress", Column: "ipaddress", Field: "ipaddress"})
}

// UpdateLocation - POST /api/settings/updatelocation
// @Summary Update Location
// @Description Partial update of sd_iot_location by location_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include location_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatelocation [post]
func (h *settingsHandler) UpdateLocation() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_location", "location_id", locationCols)
}

// DeleteLocation - GET /api/settings/deletelocation
// @Summary Delete Location
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param location_id query string true "location_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletelocation [get]
func (h *settingsHandler) DeleteLocation() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_location", "location_id", "location_id")
}

// ---- type ----

// ListType - GET /api/settings/listtype
// @Summary List Type
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
// @Router /settings/listtype [get]
func (h *settingsHandler) ListType() http.HandlerFunc { return h.listHandler(settings.TypeSpec) }

// TypeAll - GET /api/settings/typeall
// @Summary Type All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/typeall [get]
func (h *settingsHandler) TypeAll() http.HandlerFunc { return h.allHandler(settings.TypeSpec) }

// CreateType - POST /api/settings/createtype
// @Summary Create Type
// @Description Creates one row in sd_iot_type (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createtype [post]
func (h *settingsHandler) CreateType() http.HandlerFunc {
	return h.createHandler("sd_iot_type", typeCols, statusDefault,
		dupCheck{BodyKey: "type_name", Column: "type_name", Field: "type_name"})
}

// UpdateType - POST /api/settings/updatetype
// @Summary Update Type
// @Description Partial update of sd_iot_type by type_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include type_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatetype [post]
func (h *settingsHandler) UpdateType() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_type", "type_id", typeCols)
}

// DeleteType - GET /api/settings/deletetype
// @Summary Delete Type
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param type_id query string true "type_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletetype [get]
func (h *settingsHandler) DeleteType() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_type", "type_id", "type_id")
}

// ---- devicetype ----

// ListDeviceTypePage - GET /api/settings/listdevicetype
// @Summary List Device Type Page
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
// @Router /settings/listdevicetype [get]
func (h *settingsHandler) ListDeviceTypePage() http.HandlerFunc {
	return h.listHandler(settings.DeviceTypeSpec)
}

// DeviceTypeAll - GET /api/settings/devicetypeall
// @Summary Device Type All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/devicetypeall [get]
func (h *settingsHandler) DeviceTypeAll() http.HandlerFunc {
	return h.allHandler(settings.DeviceTypeSpec)
}

// DeviceTypeAllControl - GET /api/settings/devicetypeallcontrol
// @Summary Device Type All Control
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/devicetypeallcontrol [get]
func (h *settingsHandler) DeviceTypeAllControl() http.HandlerFunc {
	return h.allHandler(settings.DeviceTypeSpec.With())
}

// CreateDeviceType - POST /api/settings/createdevicetype
// @Summary Create Device Type
// @Description Creates one row in sd_iot_device_type (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createdevicetype [post]
func (h *settingsHandler) CreateDeviceType() http.HandlerFunc {
	return h.createHandler("sd_iot_device_type", deviceTypeCols, statusDefault,
		dupCheck{BodyKey: "type_name", Column: "type_name", Field: "type_name"})
}

// UpdateDeviceType - POST /api/settings/updatedevicetype
// @Summary Update Device Type
// @Description Partial update of sd_iot_device_type by type_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include type_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatedevicetype [post]
func (h *settingsHandler) UpdateDeviceType() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_device_type", "type_id", deviceTypeCols)
}

// DeleteDeviceType - GET /api/settings/deletedevicetype
// @Summary Delete Device Type
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param type_id query string true "type_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletedevicetype [get]
func (h *settingsHandler) DeleteDeviceType() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_device_type", "type_id", "type_id")
}

// ---- group ----

// ListGroup - GET /api/settings/lisgroup
// @Summary List Group
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
// @Router /settings/lisgroup [get]
func (h *settingsHandler) ListGroup() http.HandlerFunc { return h.listHandler(settings.GroupSpec) }

// GroupAll - GET /api/settings/lisgroupall
// @Summary Group All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/lisgroupall [get]
func (h *settingsHandler) GroupAll() http.HandlerFunc { return h.allHandler(settings.GroupSpec) }

// ListGroupPage - GET /api/settings/listgrouppage
// @Summary List Group Page
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
// @Router /settings/listgrouppage [get]
func (h *settingsHandler) ListGroupPage() http.HandlerFunc { return h.listHandler(settings.GroupSpec) }

// CreateGroup - POST /api/settings/creategroup
// @Summary Create Group
// @Description Creates one row in sd_iot_group (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/creategroup [post]
func (h *settingsHandler) CreateGroup() http.HandlerFunc {
	return h.createHandler("sd_iot_group", groupCols, statusDefault,
		dupCheck{BodyKey: "group_name", Column: "group_name", Field: "group_name"})
}

// UpdateGroup - POST /api/settings/updategroup
// @Summary Update Group
// @Description Partial update of sd_iot_group by group_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include group_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updategroup [post]
func (h *settingsHandler) UpdateGroup() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_group", "group_id", groupCols)
}

// DeleteGroup - GET /api/settings/deletegroup
// @Summary Delete Group
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param group_id query string true "group_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletegroup [get]
func (h *settingsHandler) DeleteGroup() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_group", "group_id", "group_id")
}

// ---- sensor ----

// ListSensor - GET /api/settings/listsensor
// @Summary List Sensor
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
// @Router /settings/listsensor [get]
func (h *settingsHandler) ListSensor() http.HandlerFunc { return h.listHandler(settings.SensorSpec) }

// SensorAll - GET /api/settings/sensorall
// @Summary Sensor All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/sensorall [get]
func (h *settingsHandler) SensorAll() http.HandlerFunc { return h.allHandler(settings.SensorSpec) }

// CreateSensor - POST /api/settings/createsensor
// @Summary Create Sensor
// @Description Creates one row in sd_iot_sensor (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createsensor [post]
func (h *settingsHandler) CreateSensor() http.HandlerFunc {
	return h.createHandler("sd_iot_sensor", sensorCols, statusDefault,
		dupCheck{BodyKey: "sensor_name", Column: "sensor_name", Field: "sensor_name"})
}

// UpdateSensor - POST /api/settings/updatesensor
// @Summary Update Sensor
// @Description Partial update of sd_iot_sensor by sensor_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include sensor_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatesensor [post]
func (h *settingsHandler) UpdateSensor() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_sensor", "sensor_id", sensorCols)
}

// DeleteSensor - GET /api/settings/deletesensor
// @Summary Delete Sensor
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param sensor_id query string true "sensor_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletesensor [get]
func (h *settingsHandler) DeleteSensor() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_sensor", "sensor_id", "sensor_id")
}
