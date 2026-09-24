package http

import (
	"net/http"

	_ "icmongolang/internal/modules/settings/presenter"
	_ "icmongolang/pkg/responses"
)

// SensorType ports GET /sensortype (query lang=th selects the Thai list).
// SensorType - GET /api/settings/sensortype
// @Summary Sensor Type
// @Description Static sensor-type catalog; lang=th returns Thai names.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param lang query string false "Language code: en (default) or th"
// @Success 200 {object} responses.SwaggerSuccessResponse "Sensor type rows"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/sensortype [get]
func (h *settingsHandler) SensorType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok(w, r, h.uc.SensorTypes(r.URL.Query().Get("lang")))
	}
}

// TestGmailConnection ports GET /testgemail.
// TestGmailConnection - GET /api/settings/testgemail
// @Summary Test Gmail Connection
// @Description Verifies smtp.gmail.com on ports 465/587 using env credentials and sends a test mail.
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} responses.SuccessResponse[presenter.SendEmailResult] "Email test/send result"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/testgemail [get]
func (h *settingsHandler) TestGmailConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := h.uc.TestGmailConnection(r.Context())
		ok(w, r, res)
	}
}

// SendEmail ports GET /sendemail (stub success — original never sends).
// SendEmail - GET /api/settings/sendemail
// @Summary Send Email
// @Description Stub send-email endpoint mirroring NestJS behavior (returns success payload without sending).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param to query string false "Recipient email"
// @Param subject query string false "Subject"
// @Param content query string false "Content"
// @Success 200 {object} responses.SuccessResponse[presenter.SendEmailResult] "Email test/send result"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/sendemail [get]
func (h *settingsHandler) SendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		res := h.uc.SendEmailStub(q.Get("to"), q.Get("subject"), q.Get("content"))
		ok(w, r, res)
	}
}

// MqttData ports GET /mqttdata: Redis cache -> MQTT request-and-wait.
// MqttData - GET /api/settings/mqttdata
// @Summary Mqtt Data
// @Description Reads live MQTT payload for a topic: Redis cache first, falls back to request-and-wait (result cached 30s).
// @Tags settings
// @Security BearerAuth
// @Produce json
// @Param mqttdata query string true "MQTT topic"
// @Success 200 {object} responses.SuccessResponse[presenter.MqttDataResult] "Cached or live MQTT payload"
// @Failure 400 {object} responses.ErrorResponse "Bad request / validation error"
// @Router /settings/mqttdata [get]
func (h *settingsHandler) MqttData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.MqttData(r.Context(), r.URL.Query().Get("mqttdata"))
		if err != nil {
			fail(w, r, err)
			return
		}
		ok(w, r, res)
	}
}
