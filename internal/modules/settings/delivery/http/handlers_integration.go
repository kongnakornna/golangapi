package http

import (
	"net/http"
	"strings"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"

	"icmongolang/internal/modules/settings"

	"github.com/go-chi/render"
)

var mqttCols = []string{
	"mqtt_type_id", "sort", "mqtt_name", "host", "port", "username", "password",
	"secret", "expire_in", "token_value", "org", "bucket", "envavorment",
	"location_id", "latitude", "longitude", "mqtt_main_id", "configuration",
	"zoom", "status",
}
var mqttHostCols = []string{"hostname", "host", "port", "username", "password", "idhost", "status"}
var apiCols = []string{"api_name", "host", "port", "token_value", "status"}
var emailCols = []string{"email_name", "host", "port", "username", "password", "status"}
var hostCols = []string{"host_name", "port", "username", "password", "idhost", "status"}
var influxdbCols = []string{"influxdb_name", "host", "port", "username", "password", "token_value", "buckets", "status"}
var lineCols = []string{"line_name", "client_id", "client_secret", "secret_key", "redirect_uri", "grant_type", "code", "accesstoken", "status"}
var noderedCols = []string{"nodered_name", "host", "port", "routing", "client_id", "grant_type", "scope", "username", "password", "status"}
var smsCols = []string{"sms_name", "host", "port", "username", "password", "apikey", "originator", "status"}
var tokenCols = []string{"token_name", "host", "port", "token_value", "status"}
var telegramCols = []string{"telegram_name", "port", "username", "password", "status"}

// ---- mqtt ----

// ListMqtt - GET /api/settings/lismqtt
// @Summary List Mqtt
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
// @Router /settings/lismqtt [get]
func (h *settingsHandler) ListMqtt() http.HandlerFunc { return h.listHandler(settings.MqttSpec) }

// ListMqttAlt - GET /api/settings/listmqtt
// @Summary List Mqtt Alt
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
// @Router /settings/listmqtt [get]
func (h *settingsHandler) ListMqttAlt() http.HandlerFunc { return h.listHandler(settings.MqttSpec) }

// MqttAll - GET /api/settings/lismqttall
// @Summary Mqtt All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/lismqttall [get]
func (h *settingsHandler) MqttAll() http.HandlerFunc { return h.allHandler(settings.MqttSpec) }

// ListMqttPaginate - GET /api/settings/listmqttpaginate
// @Summary List Mqtt Paginate
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
// @Router /settings/listmqttpaginate [get]
func (h *settingsHandler) ListMqttPaginate() http.HandlerFunc {
	return h.listHandler(settings.MqttSpec)
}

// ListMqttPaginateActive - GET /api/settings/listmqttpaginateactive
// @Summary List Mqtt Paginate Active
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
// @Router /settings/listmqttpaginateactive [get]
func (h *settingsHandler) ListMqttPaginateActive() http.HandlerFunc {
	return h.listHandler(settings.MqttSpec.With(settings.FixedFilter{Column: "m.status", Value: "1"}))
}

// ListMqttDevicePaginate ports /listmqttdevicepaginate: devices of one broker.
// ListMqttDevicePaginate - GET /api/settings/listmqttdevicepaginate
// @Summary List Mqtt Device Paginate
// @Description Paginated devices of one MQTT broker.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param mqtt_id query int false "MQTT broker ID"
// @Success 200 {object} responses.SuccessResponse[presenter.ListResult] "Paginated result envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/listmqttdevicepaginate [get]
func (h *settingsHandler) ListMqttDevicePaginate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mqttID := r.URL.Query().Get("mqtt_id")
		spec := settings.DeviceSpec
		if mqttID != "" {
			spec = spec.With(settings.FixedFilter{Column: "d.mqtt_id", Value: mqttID})
		}
		res, err := h.uc.ListPaginate(r.Context(), spec, listQueryFrom(r))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, res)
	}
}

// GetMqttDetail ports /getmqttdetail: one broker row by mqtt_id.
// GetMqttDetail - GET /api/settings/getmqttdetail
// @Summary Get Mqtt Detail
// @Description Fetches one record by mqtt_id.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param mqtt_id query string true "mqtt_id"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/getmqttdetail [get]
func (h *settingsHandler) GetMqttDetail() http.HandlerFunc {
	return h.getOneByParam(settings.MqttSpec, "mqtt_id", "m.mqtt_id", true, "mqtt not found")
}

// DeleteMqtt - GET /api/settings/mqttdelete
// @Summary Delete Mqtt
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param mqtt_id query string true "mqtt_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/mqttdelete [get]
func (h *settingsHandler) DeleteMqtt() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_mqtt", "mqtt_id", "mqtt_id")
}

// CreateMqtt - POST /api/settings/createmqtt
// @Summary Create Mqtt
// @Description Creates one row in sd_iot_mqtt (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createmqtt [post]
func (h *settingsHandler) CreateMqtt() http.HandlerFunc {
	return h.createHandler("sd_iot_mqtt", mqttCols, map[string]interface{}{"status": 1, "sort": 1},
		dupCheck{BodyKey: "mqtt_name", Column: "mqtt_name", Field: "mqtt_name"},
		dupCheck{BodyKey: "bucket", Column: "bucket", Field: "bucket"})
}

// UpdateMqtt - POST /api/settings/updatemqtt
// @Summary Update Mqtt
// @Description Partial update of sd_iot_mqtt by mqtt_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include mqtt_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatemqtt [post]
func (h *settingsHandler) UpdateMqtt() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_mqtt", "mqtt_id", mqttCols)
}

// UpdateMqttStatus - POST /api/settings/updatemqttstatus
// @Summary Update Mqtt Status
// @Description Sets one row's status by mqtt_id (target-only update, 404 when missing).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "{mqtt_id, status}"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatemqttstatus [post]
func (h *settingsHandler) UpdateMqttStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_mqtt", "mqtt_id", false)
}

// UpdateMqtttSort ports POST /mqtttsort -> update_mqttt_sort.
// UpdateMqtttSort - POST /api/settings/mqtttsort
// @Summary Update Mqttt Sort
// @Description Updates sort order of one MQTT broker. Body requires mqtt_id plus optional sort.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/mqtttsort [post]
func (h *settingsHandler) UpdateMqtttSort() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		id := toStr(body["mqtt_id"])
		if id == "" {
			badRequest(w, r, "mqtt_id is required")
			return
		}
		fields := filterKeys(body, []string{"sort"})
		fields["updateddate"] = nowPtr()
		n, err := h.uc.UpdateFields(r.Context(), "sd_iot_mqtt", "mqtt_id", coerceID(id), fields)
		if err != nil {
			fail(w, r, err)
			return
		}
		affected(w, r, n)
	}
}

// ---- mqtthost / api / email / host / influxdb / line / nodered / sms / token / telegram ----

// ListMqttHost - GET /api/settings/listmqtthost
// @Summary List Mqtt Host
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
// @Router /settings/listmqtthost [get]
func (h *settingsHandler) ListMqttHost() http.HandlerFunc {
	return h.listHandler(settings.MqttHostSpec)
}

// MqttHostAll - GET /api/settings/mqtthostall
// @Summary Mqtt Host All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/mqtthostall [get]
func (h *settingsHandler) MqttHostAll() http.HandlerFunc { return h.allHandler(settings.MqttHostSpec) }

// CreateMqttHost - POST /api/settings/createmqtthost
// @Summary Create Mqtt Host
// @Description Creates one row in sd_mqtt_host (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createmqtthost [post]
func (h *settingsHandler) CreateMqttHost() http.HandlerFunc {
	return h.createHandler("sd_mqtt_host", mqttHostCols, statusDefault,
		dupCheck{BodyKey: "hostname", Column: "hostname", Field: "hostname"})
}

// UpdateMqttHost - POST /api/settings/updatemqtthost
// @Summary Update Mqtt Host
// @Description Partial update of sd_mqtt_host by id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatemqtthost [post]
func (h *settingsHandler) UpdateMqttHost() http.HandlerFunc {
	return h.updateBodyHandler("sd_mqtt_host", "id", mqttHostCols)
}

// UpdateMqttHostStatus - POST /api/settings/updatemqtthoststatus
// @Summary Update Mqtt Host Status
// @Description Partial update of sd_mqtt_host by id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatemqtthoststatus [post]
func (h *settingsHandler) UpdateMqttHostStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_mqtt_host", "id", true)
}

// DeleteMqttHost - GET /api/settings/deletemqtthost
// @Summary Delete Mqtt Host
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletemqtthost [get]
func (h *settingsHandler) DeleteMqttHost() http.HandlerFunc {
	return h.deleteGetHandler("sd_mqtt_host", "id::text", "id")
}

// ApiAll - GET /api/settings/apiall
// @Summary Api All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/apiall [get]
func (h *settingsHandler) ApiAll() http.HandlerFunc { return h.allHandler(settings.ApiSpec) }

// ListApiPage - GET /api/settings/listapipage
// @Summary List Api Page
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
// @Router /settings/listapipage [get]
func (h *settingsHandler) ListApiPage() http.HandlerFunc { return h.listHandler(settings.ApiSpec) }

// CreateApi - POST /api/settings/createapi
// @Summary Create Api
// @Description Creates one row in sd_iot_api (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createapi [post]
func (h *settingsHandler) CreateApi() http.HandlerFunc {
	return h.createHandler("sd_iot_api", apiCols, statusDefault,
		dupCheck{BodyKey: "api_name", Column: "api_name", Field: "api_name"})
}

// UpdateApi - POST /api/settings/updateapi
// @Summary Update Api
// @Description Partial update of sd_iot_api by api_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include api_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateapi [post]
func (h *settingsHandler) UpdateApi() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_api", "api_id", apiCols)
}

// DeleteApi - GET /api/settings/deleteapi
// @Summary Delete Api
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param api_id query string true "api_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deleteapi [get]
func (h *settingsHandler) DeleteApi() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_api", "api_id", "api_id")
}

// ListEmail - GET /api/settings/listemail
// @Summary List Email
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
// @Router /settings/listemail [get]
func (h *settingsHandler) ListEmail() http.HandlerFunc { return h.listHandler(settings.EmailSpec) }

// EmailAll - GET /api/settings/emailall
// @Summary Email All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/emailall [get]
func (h *settingsHandler) EmailAll() http.HandlerFunc { return h.allHandler(settings.EmailSpec) }

// CreateEmail - POST /api/settings/createemail
// @Summary Create Email
// @Description Creates one row in sd_iot_email (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createemail [post]
func (h *settingsHandler) CreateEmail() http.HandlerFunc {
	return h.createHandler("sd_iot_email", emailCols, statusDefault,
		dupCheck{BodyKey: "email_name", Column: "email_name", Field: "email_name"})
}

// UpdateEmail - POST /api/settings/updateemail
// @Summary Update Email
// @Description Partial update of sd_iot_email by email_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include email_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateemail [post]
func (h *settingsHandler) UpdateEmail() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_email", "email_id", emailCols)
}

// UpdateEmailStatus - POST /api/settings/updateemailstatus
// @Summary Update Email Status
// @Description Partial update of sd_iot_email by email_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include email_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateemailstatus [post]
func (h *settingsHandler) UpdateEmailStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_email", "email_id", true)
}

// DeleteEmail - GET /api/settings/deleteemail
// @Summary Delete Email
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param email_id query string true "email_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deleteemail [get]
func (h *settingsHandler) DeleteEmail() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_email", "email_id::text", "email_id")
}

// HostAll - GET /api/settings/hostall
// @Summary Host All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/hostall [get]
func (h *settingsHandler) HostAll() http.HandlerFunc { return h.allHandler(settings.HostSpec) }

// ListHostPage - GET /api/settings/listhostpage
// @Summary List Host Page
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
// @Router /settings/listhostpage [get]
func (h *settingsHandler) ListHostPage() http.HandlerFunc { return h.listHandler(settings.HostSpec) }

// CreateHost - POST /api/settings/createhost
// @Summary Create Host
// @Description Creates one row in sd_iot_host (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createhost [post]
func (h *settingsHandler) CreateHost() http.HandlerFunc {
	return h.createHandler("sd_iot_host", hostCols, statusDefault,
		dupCheck{BodyKey: "host_name", Column: "host_name", Field: "host_name"})
}

// UpdateHost - POST /api/settings/updatehost
// @Summary Update Host
// @Description Partial update of sd_iot_host by host_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include host_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatehost [post]
func (h *settingsHandler) UpdateHost() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_host", "host_id", hostCols)
}

// DeleteHost - GET /api/settings/deletehost
// @Summary Delete Host
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param host_id query string true "host_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletehost [get]
func (h *settingsHandler) DeleteHost() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_host", "host_id::text", "host_id")
}

// InfluxdbAll - GET /api/settings/influxdball
// @Summary Influxdb All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/influxdball [get]
func (h *settingsHandler) InfluxdbAll() http.HandlerFunc { return h.allHandler(settings.InfluxdbSpec) }

// ListInfluxdbPage - GET /api/settings/listinfluxdbpage
// @Summary List Influxdb Page
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
// @Router /settings/listinfluxdbpage [get]
func (h *settingsHandler) ListInfluxdbPage() http.HandlerFunc {
	return h.listHandler(settings.InfluxdbSpec)
}

// CreateInfluxdb - POST /api/settings/createinfluxdb
// @Summary Create Influxdb
// @Description Creates one row in sd_iot_influxdb (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createinfluxdb [post]
func (h *settingsHandler) CreateInfluxdb() http.HandlerFunc {
	return h.createHandler("sd_iot_influxdb", influxdbCols, statusDefault)
}

// UpdateInfluxdb - POST /api/settings/updateinfluxdb
// @Summary Update Influxdb
// @Description Partial update of sd_iot_influxdb by influxdb_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include influxdb_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateinfluxdb [post]
func (h *settingsHandler) UpdateInfluxdb() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_influxdb", "influxdb_id", influxdbCols)
}

// UpdateInfluxdbStatus - POST /api/settings/updateinfluxdbstatus
// @Summary Update Influxdb Status
// @Description Activates/deactivates one row by influxdb_id; activating resets all other rows to status 0 first.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "{influxdb_id, status}"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateinfluxdbstatus [post]
func (h *settingsHandler) UpdateInfluxdbStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_influxdb", "influxdb_id", true)
}

// DeleteInfluxdb - GET /api/settings/deleteinfluxdb
// @Summary Delete Influxdb
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param influxdb_id query string true "influxdb_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deleteinfluxdb [get]
func (h *settingsHandler) DeleteInfluxdb() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_influxdb", "influxdb_id::text", "influxdb_id")
}

// LineAll - GET /api/settings/lineall
// @Summary Line All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/lineall [get]
func (h *settingsHandler) LineAll() http.HandlerFunc { return h.allHandler(settings.LineSpec) }

// ListLinePage - GET /api/settings/listlinepage
// @Summary List Line Page
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
// @Router /settings/listlinepage [get]
func (h *settingsHandler) ListLinePage() http.HandlerFunc { return h.listHandler(settings.LineSpec) }

// CreateLine - POST /api/settings/createline
// @Summary Create Line
// @Description Creates one row in sd_iot_line (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createline [post]
func (h *settingsHandler) CreateLine() http.HandlerFunc {
	return h.createHandler("sd_iot_line", lineCols, statusDefault,
		dupCheck{BodyKey: "line_name", Column: "line_name", Field: "line_name"})
}

// UpdateLine - POST /api/settings/updateline
// @Summary Update Line
// @Description Partial update of sd_iot_line by line_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include line_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updateline [post]
func (h *settingsHandler) UpdateLine() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_line", "line_id", lineCols)
}

// UpdateLineStatus - POST /api/settings/updatelinestatus
// @Summary Update Line Status
// @Description Activates/deactivates one row by line_id; activating resets all other rows to status 0 first.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "{line_id, status}"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatelinestatus [post]
func (h *settingsHandler) UpdateLineStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_line", "line_id", true)
}

// DeleteLine - GET /api/settings/deleteline
// @Summary Delete Line
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param line_id query string true "line_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deleteline [get]
func (h *settingsHandler) DeleteLine() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_line", "line_id::text", "line_id")
}

// NoderedAll - GET /api/settings/noderedall
// @Summary Nodered All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/noderedall [get]
func (h *settingsHandler) NoderedAll() http.HandlerFunc { return h.allHandler(settings.NoderedSpec) }

// ListNoderedPaginate - GET /api/settings/listnoderedpaginate
// @Summary List Nodered Paginate
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
// @Router /settings/listnoderedpaginate [get]
func (h *settingsHandler) ListNoderedPaginate() http.HandlerFunc {
	return h.listHandler(settings.NoderedSpec)
}

// CreateNodered - POST /api/settings/createnodered
// @Summary Create Nodered
// @Description Creates one row in sd_iot_nodered (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createnodered [post]
func (h *settingsHandler) CreateNodered() http.HandlerFunc {
	return h.createHandler("sd_iot_nodered", noderedCols, statusDefault,
		dupCheck{BodyKey: "nodered_name", Column: "nodered_name", Field: "nodered_name"})
}

// UpdateNodered - POST /api/settings/updatenodered
// @Summary Update Nodered
// @Description Partial update of sd_iot_nodered by nodered_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include nodered_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatenodered [post]
func (h *settingsHandler) UpdateNodered() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_nodered", "nodered_id", noderedCols)
}

// UpdateNoderedStatus - POST /api/settings/updatenoderedstatus
// @Summary Update Nodered Status
// @Description Activates/deactivates one row by nodered_id; activating resets all other rows to status 0 first.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "{nodered_id, status}"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatenoderedstatus [post]
func (h *settingsHandler) UpdateNoderedStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_nodered", "nodered_id", true)
}

// DeleteNodered - GET /api/settings/deletenodered
// @Summary Delete Nodered
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param nodered_id query string true "nodered_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletenodered [get]
func (h *settingsHandler) DeleteNodered() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_nodered", "nodered_id::text", "nodered_id")
}

// SmsAll - GET /api/settings/smsall
// @Summary Sms All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/smsall [get]
func (h *settingsHandler) SmsAll() http.HandlerFunc { return h.allHandler(settings.SmsSpec) }

// ListSmsPage - GET /api/settings/listsmspage
// @Summary List Sms Page
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
// @Router /settings/listsmspage [get]
func (h *settingsHandler) ListSmsPage() http.HandlerFunc { return h.listHandler(settings.SmsSpec) }

// CreateSms - POST /api/settings/createsms
// @Summary Create Sms
// @Description Creates one row in sd_iot_sms (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createsms [post]
func (h *settingsHandler) CreateSms() http.HandlerFunc {
	return h.createHandler("sd_iot_sms", smsCols, statusDefault,
		dupCheck{BodyKey: "sms_name", Column: "sms_name", Field: "sms_name"})
}

// UpdateSms - POST /api/settings/updatesms
// @Summary Update Sms
// @Description Partial update of sd_iot_sms by sms_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include sms_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatesms [post]
func (h *settingsHandler) UpdateSms() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_sms", "sms_id", smsCols)
}

// UpdateSmsStatus - POST /api/settings/updatesmsstatus
// @Summary Update Sms Status
// @Description Activates/deactivates one row by sms_id; activating resets all other rows to status 0 first.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "{sms_id, status}"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatesmsstatus [post]
func (h *settingsHandler) UpdateSmsStatus() http.HandlerFunc {
	return h.statusBodyHandler("sd_iot_sms", "sms_id", true)
}

// DeleteSms - GET /api/settings/deletesms
// @Summary Delete Sms
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param sms_id query string true "sms_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletesms [get]
func (h *settingsHandler) DeleteSms() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_sms", "sms_id::text", "sms_id")
}

// TokenAll - GET /api/settings/tokenall
// @Summary Token All
// @Description Full listing without pagination.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Keyword filter (LIKE)"
// @Success 200 {object} responses.SuccessResponse[[]presenter.Row] "Full list envelope"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/tokenall [get]
func (h *settingsHandler) TokenAll() http.HandlerFunc { return h.allHandler(settings.TokenSpec) }

// ListTokenPage - GET /api/settings/tokensmspage
// @Summary List Token Page
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
// @Router /settings/tokensmspage [get]
func (h *settingsHandler) ListTokenPage() http.HandlerFunc { return h.listHandler(settings.TokenSpec) }

// CreateToken - POST /api/settings/createtoken
// @Summary Create Token
// @Description Creates one row in sd_iot_token (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createtoken [post]
func (h *settingsHandler) CreateToken() http.HandlerFunc {
	return h.createHandler("sd_iot_token", tokenCols, statusDefault,
		dupCheck{BodyKey: "token_name", Column: "token_name", Field: "token_name"})
}

// UpdateToken - POST /api/settings/updatetoken
// @Summary Update Token
// @Description Partial update of sd_iot_token by token_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include token_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatetoken [post]
func (h *settingsHandler) UpdateToken() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_token", "token_id", tokenCols)
}

// DeleteToken - GET /api/settings/deletetoken
// @Summary Delete Token
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param token_id query string true "token_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletetoken [get]
func (h *settingsHandler) DeleteToken() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_token", "token_id", "token_id")
}

// CreateTelegram - POST /api/settings/createtelegram
// @Summary Create Telegram
// @Description Creates one row in sd_iot_telegram (whitelisted fields; defaults applied).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Resource fields"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created / fetched row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createtelegram [post]
func (h *settingsHandler) CreateTelegram() http.HandlerFunc {
	return h.createHandler("sd_iot_telegram", telegramCols, statusDefault,
		dupCheck{BodyKey: "telegram_name", Column: "telegram_name", Field: "telegram_name"})
}

// UpdateTelegram - POST /api/settings/updatetelegram
// @Summary Update Telegram
// @Description Partial update of sd_iot_telegram by telegram_id in body (only provided fields + updateddate).
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Fields to update (must include telegram_id)"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/updatetelegram [post]
func (h *settingsHandler) UpdateTelegram() http.HandlerFunc {
	return h.updateBodyHandler("sd_iot_telegram", "telegram_id", telegramCols)
}

// DeleteTelegram - GET /api/settings/deletetelegram
// @Summary Delete Telegram
// @Description Deletes matching rows; returns affected count.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param telegram_id query string true "telegram_id"
// @Success 200 {object} responses.SuccessResponse[presenter.Affected] "Affected rows count"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Failure 404 {object} responses.ErrorResponse "Not found"
// @Failure 500 {object} responses.ErrorResponse "Internal error"
// @Router /settings/deletetelegram [get]
func (h *settingsHandler) DeleteTelegram() http.HandlerFunc {
	return h.deleteGetHandler("sd_iot_telegram", "telegram_id::text", "telegram_id")
}

// DeleteMqttAlt ports GET /deletemqtt: existence-checked delete by mqtt_id.
// DeleteMqttAlt - GET /api/settings/deletemqtt
// @Summary Delete Mqtt (checked)
// @Description Deletes one broker after verifying it exists (NestJS /deletemqtt).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param mqtt_id query string true "mqtt_id"
// @Success 200 {object} responses.SuccessResponse[string] "Deleted complete"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/deletemqtt [get]
func (h *settingsHandler) DeleteMqttAlt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.URL.Query().Get("mqtt_id"))
		if id == "" {
			render.Render(w, r, response422("mqtt_id is null."))
			return
		}
		exists, err := h.uc.ExistsWhere(r.Context(), "sd_iot_mqtt", "mqtt_id = ?", coerceID(id))
		if err != nil {
			fail(w, r, err)
			return
		}
		if !exists {
			render.Render(w, r, response422("uid is null."))
			return
		}
		n, err := h.uc.DeleteRow(r.Context(), "sd_iot_mqtt", "mqtt_id = ?", coerceID(id))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, map[string]interface{}{"mqtt_id": id, "affected": n})
	}
}

// CreateAlarmDevicePaginate ports POST /createalarmdevicepaginate:
// creates a setting row after rejecting duplicate SN values.
// CreateAlarmDevicePaginate - POST /api/settings/createalarmdevicepaginate
// @Summary Create Alarm Device Paginate
// @Description Creates a setting (sd_iot_setting); rejects duplicate sn with 422.
// @Tags settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Setting fields (sn required)"
// @Success 200 {object} responses.SwaggerSuccessResponse "Created row"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/createalarmdevicepaginate [post]
func (h *settingsHandler) CreateAlarmDevicePaginate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := decodeBody(r)
		if err != nil {
			badRequest(w, r, err.Error())
			return
		}
		sn := strings.TrimSpace(toStr(body["sn"]))
		if sn == "" {
			badRequest(w, r, "sn is required")
			return
		}
		dup, err := h.uc.ExistsWhere(r.Context(), "sd_iot_setting", "sn = ?", sn)
		if err != nil {
			fail(w, r, err)
			return
		}
		if dup {
			render.Render(w, r, response422("The SN duplicate this data cannot createddate."))
			return
		}
		row := filterKeys(body, settingCols)
		if len(row) == 0 {
			badRequest(w, r, "request body has no valid fields")
			return
		}
		for k, v := range statusDefault {
			if _, exists := row[k]; !exists {
				row[k] = v
			}
		}
		if err := h.uc.CreateRow(r.Context(), "sd_iot_setting", row); err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, row)
	}
}
