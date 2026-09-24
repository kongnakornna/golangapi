package http

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"icmongolang/internal/modules/mqtt/presenter"
	"icmongolang/internal/modules/mqtt/usecase"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/logger"

	"github.com/go-chi/render"
)

// Cache interface for caching
type Cache interface {
	Get(ctx context.Context, key string, dst interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// MQTTHandler handles MQTT HTTP endpoints
type MQTTHandler struct {
	uc     usecase.MQTTUseCase
	logger logger.Logger
	cache  Cache
}

type cachedTopicData struct {
	Data     []byte    `json:"data"`
	CachedAt time.Time `json:"cached_at"`
}

type cachedDeviceControl struct {
	Response presenter.DeviceControlResponse `json:"response"`
	CachedAt time.Time                       `json:"cached_at"`
}

// NewMQTTHandler creates a new MQTTHandler
func NewMQTTHandler(uc usecase.MQTTUseCase, log logger.Logger, cache Cache) *MQTTHandler {
	return &MQTTHandler{uc: uc, logger: log, cache: cache}
}

// CreateMQTTHandler is an alias for NewMQTTHandler
func CreateMQTTHandler(uc usecase.MQTTUseCase, log logger.Logger, cache Cache) *MQTTHandler {
	return &MQTTHandler{uc: uc, logger: log, cache: cache}
}

// Publish godoc
// @Summary      Publish an MQTT message
// @Description  Publishes a payload to a specified MQTT topic. Supports QoS 0,1,2 and retained flag.
// @Tags         mqtt
// @Accept       json
// @Produce      json
// @Param        request body presenter.PublishRequest true "Publish request"
// @Success      200 {object} presenter.PublishResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /mqtt/publish [post]
// @Security     BearerAuth
func (h *MQTTHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req presenter.PublishRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Errorf("decode publish request error: %v", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if req.Topic == "" {
		render.Render(w, r, ErrInvalidRequestString("topic is required"))
		return
	}
	if req.Payload == nil {
		render.Render(w, r, ErrInvalidRequestString("payload is required"))
		return
	}
	if err := h.uc.Publish(r.Context(), &req); err != nil {
		h.logger.Errorf("publish error: %v", err)
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, presenter.PublishResponse{Success: true, Message: "published successfully"})
}

// Subscribe godoc
// @Summary      Subscribe to an MQTT topic
// @Description  Subscribes to a topic and starts receiving messages (persistent subscription). The connection will stay open and messages will be broadcast via WebSocket.
// @Tags         mqtt
// @Accept       json
// @Produce      json
// @Param        request body presenter.SubscribeRequest true "Subscribe request"
// @Success      200 {object} presenter.SubscribeResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /mqtt/subscribe [post]
// @Security     BearerAuth
func (h *MQTTHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req presenter.SubscribeRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Errorf("decode subscribe request error: %v", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if req.Topic == "" {
		render.Render(w, r, ErrInvalidRequestString("topic is required"))
		return
	}
	if req.QoS > 2 {
		render.Render(w, r, ErrInvalidRequestString("QoS must be 0, 1, or 2"))
		return
	}
	if err := h.uc.Subscribe(r.Context(), &req); err != nil {
		h.logger.Errorf("subscribe error: %v", err)
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, presenter.SubscribeResponse{
		Topic:        req.Topic,
		QoS:          req.QoS,
		SubscribedAt: time.Now(),
		Message:      "subscribed successfully",
	})
}

// Unsubscribe godoc
// @Summary      Unsubscribe from an MQTT topic
// @Description  Removes a previous subscription from a topic. Messages will no longer be received.
// @Tags         mqtt
// @Accept       json
// @Produce      json
// @Param        request body presenter.UnsubscribeRequest true "Unsubscribe request"
// @Success      200 {object} presenter.UnsubscribeResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /mqtt/unsubscribe [post]
// @Security     BearerAuth
func (h *MQTTHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	var req presenter.UnsubscribeRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Errorf("decode unsubscribe request error: %v", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if req.Topic == "" {
		render.Render(w, r, ErrInvalidRequestString("topic is required"))
		return
	}
	if err := h.uc.Unsubscribe(r.Context(), &req); err != nil {
		h.logger.Errorf("unsubscribe error: %v", err)
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, presenter.UnsubscribeResponse{Success: true, Topic: req.Topic, Message: "unsubscribed successfully"})
}

// Subscriptions godoc
// @Summary      List all subscribed topics
// @Description  Returns an array of topics that the MQTT client is currently subscribed to (persistent subscriptions).
// @Tags         mqtt
// @Produce      json
// @Success      200 {object} presenter.SubscriptionsResponse
// @Failure      500 {object} errResponse
// @Router       /mqtt/subscriptions [get]
// @Security     BearerAuth
func (h *MQTTHandler) Subscriptions(w http.ResponseWriter, r *http.Request) {
	topics := h.uc.GetSubscriptions()
	render.JSON(w, r, presenter.SubscriptionsResponse{Topics: topics, Count: len(topics)})
}

// Status godoc
// @Summary      Get MQTT connection status
// @Description  Returns the current connection state of the MQTT client to the broker, along with a timestamp (Asia/Bangkok).
// @Tags         mqtt
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} errResponse
// @Router       /mqtt/status [get]
// @Security     BearerAuth
func (h *MQTTHandler) Status(w http.ResponseWriter, r *http.Request) {
	connected := h.uc.IsConnected()
	status := "connected"
	if !connected {
		status = "disconnected"
	}
	TimeLoc := helpers.GetTimeLocation()
	render.JSON(w, r, map[string]interface{}{
		"connected": connected,
		"status":    status,
		"timestamp": helpers.TimeConvertermas(time.Now().In(TimeLoc)),
	})
}

// GetTopicData godoc
// @Summary      Get live or cached MQTT topic data (request‑response)
// @Description  Subscribes to the given topic, waits for a single message, then unsubscribes. Caches the result for 60 seconds.
// @Description  On MQTT timeout the endpoint returns fallback data with HTTP 504 (see `error_detail` / `from:fallback`).
// @Tags         mqtt
// @Accept       json
// @Produce      json
// @Param        topic query string true "MQTT topic name (e.g. BAACTW05/DATA)"
// @Success      200 {object} presenter.MQTTTopicDataResponse
// @Failure      400 {object} errResponse "missing topic parameter"
// @Failure      503 {object} presenter.MQTTTopicDataResponse "MQTT broker not connected"
// @Failure      504 {object} presenter.MQTTTopicDataResponse "MQTT timeout (fallback data returned)"
// @Failure      500 {object} presenter.MQTTTopicDataResponse "internal error"
// @Router       /mqtt/gettopicdata [get]
// @Security     OAuth2Password
// @Security BearerAuth
func (h *MQTTHandler) GetTopicData(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	topic := strings.TrimSpace(r.URL.Query().Get("topic"))
	if topic == "" {
		h.logger.Warn("GetTopicData: missing topic parameter")
		render.Render(w, r, ErrInvalidRequestString("topic is required"))
		return
	}
	h.logger.Infof("GetTopicData: start processing request for topic=%s", topic)

	TimeLoc := helpers.GetTimeLocation()

	mqttConnected := h.uc.IsConnected()
	if !mqttConnected {
		h.logger.Errorf("GetTopicData: MQTT client NOT connected for topic %s", topic)
		render.Status(r, http.StatusServiceUnavailable)
		render.JSON(w, r, presenter.MQTTTopicDataResponse{
			StatusCode:    http.StatusServiceUnavailable,
			Code:          http.StatusServiceUnavailable,
			Topic:         topic,
			Payload:       nil,
			From:          "error",
			Timestamp:     helpers.TimeConvertermas(time.Now().In(TimeLoc)),
			Message:       "MQTT client not connected to broker",
			MqttConnected: false,
			CacheEnabled:  h.cache != nil,
			CacheHit:      false,
			ErrorDetail:   "MQTT broker connection lost",
		})
		return
	}
	h.logger.Infof("GetTopicData: MQTT connected, starting fetch for topic %s", topic)

	ctx := r.Context()
	cacheKey := "mqtt_topic_data:" + fmt.Sprintf("%x", md5.Sum([]byte(topic)))
	cacheEnabled := h.cache != nil

	var cached cachedTopicData
	cacheHit := false
	var cachedAt time.Time
	if cacheEnabled {
		cacheGetStart := time.Now()
		if err := h.cache.Get(ctx, cacheKey, &cached); err == nil && len(cached.Data) > 0 {
			cacheHit = true
			cachedAt = cached.CachedAt
			h.logger.Infof("GetTopicData: cache HIT for topic=%s, dataLen=%d, cachedAt=%v, cacheGetDuration=%dms",
				topic, len(cached.Data), cachedAt, time.Since(cacheGetStart).Milliseconds())
			render.JSON(w, r, presenter.MQTTTopicDataResponse{
				StatusCode:      http.StatusOK,
				Code:            200,
				Topic:           topic,
				Payload:         processPayload(cached.Data),
				From:            "cache",
				Timestamp:       helpers.TimeConvertermas(cachedAt.In(TimeLoc)),
				Message:         "from cache",
				MqttConnected:   mqttConnected,
				CacheEnabled:    cacheEnabled,
				CacheHit:        cacheHit,
				DataLength:      len(cached.Data),
				FetchDurationMs: time.Since(startTime).Milliseconds(),
			})
			return
		}
		h.logger.Infof("GetTopicData: cache MISS for topic=%s", topic)
	}

	messageTimeout := 10 * time.Second
	h.logger.Infof("GetTopicData: attempting to subscribe and wait for topic %s (timeout=%v)", topic, messageTimeout)
	data, err := h.uc.GetTopic(ctx, topic, messageTimeout)
	fetchDuration := time.Since(startTime).Milliseconds()

	if err != nil {
		h.logger.Errorf("GetTopicData: MQTT get failed: %v", err)
		errorStr := err.Error()
		// Match every error the MQTT client can produce when the device does not
		// publish within the timeout (e.g. "timeout waiting for message on topic ...",
		// "subscribe timeout", "context deadline exceeded", "ResumeSubs", ...).
		isTimeout := strings.Contains(errorStr, "timeout") ||
			strings.Contains(errorStr, "context deadline exceeded") ||
			strings.Contains(errorStr, "ResumeSubs") ||
			strings.Contains(errorStr, "not currently connected")

		if isTimeout {
			h.logger.Warn("GetTopicData: MQTT timeout, returning fallback data")
			fallbackData := []byte("30.50,80.5,25.50,75.5")
			render.Status(r, http.StatusGatewayTimeout)
			render.JSON(w, r, presenter.MQTTTopicDataResponse{
				StatusCode:      http.StatusGatewayTimeout,
				Code:            http.StatusGatewayTimeout,
				Topic:           topic,
				Payload:         processPayload(fallbackData),
				From:            "fallback",
				Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
				Message:         "fallback (MQTT timeout)",
				MqttConnected:   mqttConnected,
				CacheEnabled:    cacheEnabled,
				CacheHit:        false,
				DataLength:      len(fallbackData),
				FetchDurationMs: fetchDuration,
				ErrorDetail:     err.Error(),
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, presenter.MQTTTopicDataResponse{
			StatusCode:      http.StatusInternalServerError,
			Code:            http.StatusInternalServerError,
			Topic:           topic,
			Payload:         nil,
			From:            "error",
			Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
			Message:         "failed to fetch from MQTT",
			MqttConnected:   mqttConnected,
			CacheEnabled:    cacheEnabled,
			CacheHit:        false,
			FetchDurationMs: fetchDuration,
			ErrorDetail:     err.Error(),
		})
		return
	}

	if cacheEnabled {
		cacheSetStart := time.Now()
		cacheValue := cachedTopicData{
			Data:     data,
			CachedAt: time.Now(),
		}
		if err := h.cache.Set(ctx, cacheKey, cacheValue, 60*time.Second); err != nil {
			h.logger.Warnf("GetTopicData: failed to set cache for key=%s: %v (took %dms)", cacheKey, err, time.Since(cacheSetStart).Milliseconds())
		} else {
			h.logger.Infof("GetTopicData: cached data for topic=%s, dataLen=%d, cacheWriteDuration=%dms",
				topic, len(data), time.Since(cacheSetStart).Milliseconds())
		}
	}

	renderStart := time.Now()
	render.JSON(w, r, presenter.MQTTTopicDataResponse{
		StatusCode:      http.StatusOK,
		Code:            200,
		Topic:           topic,
		Payload:         processPayload(data),
		From:            "mqtt",
		Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
		Message:         "live data",
		MqttConnected:   mqttConnected,
		CacheEnabled:    cacheEnabled,
		CacheHit:        false,
		DataLength:      len(data),
		FetchDurationMs: fetchDuration,
	})
	h.logger.Infof("GetTopicData: response rendered for topic=%s - totalDuration=%dms, renderDuration=%dms",
		topic, time.Since(startTime).Milliseconds(), time.Since(renderStart).Milliseconds())
}

// DeviceControl godoc
// @Summary      Send control command to an MQTT device and wait for response
// @Description  Publishes a control message to a topic (usually ending with CONTROL), then listens on the corresponding DATA topic (CONTROL replaced by DATA) for a response.
// @Tags         mqtt
// @Accept       json
// @Produce      json
// @Param        topic query string true "MQTT topic to send control message (e.g., BAACTW02/CONTROL)"
// @Param        message query string true "Control message (e.g., ON, OFF, 1, 0)"
// @Success      200 {object} presenter.DeviceControlResponse
// @Failure      400 {object} errResponse
// @Failure      500 {object} errResponse
// @Router       /mqtt/devicecontrol [get]
// @Security     OAuth2Password
// @Security BearerAuth
func (h *MQTTHandler) DeviceControl(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	topic := strings.TrimSpace(r.URL.Query().Get("topic"))
	message := strings.TrimSpace(r.URL.Query().Get("message"))

	TimeLoc := helpers.GetTimeLocation()

	if topic == "" || message == "" {
		render.Render(w, r, ErrInvalidRequestString("topic and message are required"))
		return
	}

	ctx := r.Context()
	cacheKey := "device_control:" + fmt.Sprintf("%x", md5.Sum([]byte(topic+":"+message)))
	cacheEnabled := h.cache != nil

	if cacheEnabled {
		var cached cachedDeviceControl
		if err := h.cache.Get(ctx, cacheKey, &cached); err == nil {
			h.logger.Infof("DeviceControl: cache HIT for %s (cachedAt=%v)", cacheKey, cached.CachedAt)
			cached.Response.FetchDurationMs = time.Since(startTime).Milliseconds()
			cached.Response.From = "cache"
			cached.Response.Timestamp = helpers.TimeConvertermas(cached.CachedAt.In(TimeLoc))
			render.JSON(w, r, cached.Response)
			return
		}
		h.logger.Debugf("DeviceControl: cache MISS for %s", cacheKey)
	}

	req := &presenter.DeviceControlRequest{Topic: topic, Message: message}
	resp, err := h.uc.DeviceControl(r.Context(), req)

	if err != nil {
		h.logger.Errorf("DeviceControl error: %v", err)
		render.JSON(w, r, presenter.DeviceControlResponse{
			StatusCode:      http.StatusInternalServerError,
			Code:            500,
			TopicControl:    topic,
			MessageSent:     message,
			Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
			Message:         "failed to execute device control",
			MessageTh:       "ไม่สามารถสั่งงานอุปกรณ์ได้",
			FetchDurationMs: time.Since(startTime).Milliseconds(),
			ErrorDetail:     err.Error(),
		})
		return
	}

	if cacheEnabled && resp.StatusCode == 200 {
		cacheValue := cachedDeviceControl{
			Response: *resp,
			CachedAt: time.Now(),
		}
		cacheValue.Response.Timestamp = helpers.TimeConvertermas(cacheValue.CachedAt.In(TimeLoc))
		if err := h.cache.Set(ctx, cacheKey, cacheValue, 5*time.Second); err != nil {
			h.logger.Warnf("DeviceControl: failed to set cache: %v", err)
		} else {
			h.logger.Infof("DeviceControl: cached response for %s (TTL=5s)", cacheKey)
		}
		render.JSON(w, r, resp)
		return
	}

	render.JSON(w, r, resp)
}

// processPayload converts byte data to appropriate type (string or []string)
func processPayload(data []byte) interface{} {
	str := strings.TrimSpace(string(data))
	if strings.Contains(str, ",") {
		return strings.Split(str, ",")
	}
	return str
}

// Error helpers

// ErrInvalidRequest returns a 400 Bad Request error
func ErrInvalidRequest(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: err.Error()}
}

// ErrInternal returns a 500 Internal Server Error
func ErrInternal(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: err.Error()}
}

// ErrInvalidRequestString returns a 400 Bad Request with a custom message
func ErrInvalidRequestString(msg string) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: msg}
}

// errResponse represents an error response
type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}

// Render implements render.Renderer interface
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
