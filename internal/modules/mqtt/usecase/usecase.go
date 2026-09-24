package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/alarm/models"
	"icmongolang/internal/modules/mqtt/presenter"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/logger"
	mqttPkg "icmongolang/pkg/mqtt"
	"icmongolang/pkg/websocket"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

type MQTTUseCase interface {
	Publish(ctx context.Context, req *presenter.PublishRequest) error
	Subscribe(ctx context.Context, req *presenter.SubscribeRequest) error
	Unsubscribe(ctx context.Context, req *presenter.UnsubscribeRequest) error
	GetSubscriptions() []string
	IsConnected() bool
	SaveAlarmLog(alarmStatus int, subject, content, deviceID string, valueData interface{}) error
	GetTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
	DeviceControl(ctx context.Context, req *presenter.DeviceControlRequest) (*presenter.DeviceControlResponse, error)
}

type AlarmLogRepository interface {
	Create(log *models.AlarmLog) error
}

type subscriptionInfo struct {
	QoS byte
}

type mqttUseCase struct {
	client        mqttPkg.Client
	cfg           *config.Config
	logger        logger.Logger
	subsMu        sync.RWMutex
	subscriptions map[string]subscriptionInfo
	wsBroadcaster websocket.Broadcaster
	db            *gorm.DB
	alarmLogRepo  AlarmLogRepository
	influxClient  influxdb.Client
	cache         Cache
}

type Cache interface {
	Get(ctx context.Context, key string, dst interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

func NewMQTTUseCaseWithClientAndWS(
	client mqttPkg.Client,
	cfg *config.Config,
	log logger.Logger,
	wsBroadcaster websocket.Broadcaster,
	db *gorm.DB,
	alarmLogRepo AlarmLogRepository,
	influxClient influxdb.Client,
	cache Cache,
) (MQTTUseCase, error) {
	if !client.IsConnected() {
		return nil, fmt.Errorf("MQTT client not connected")
	}
	log.Info("MQTT usecase with WebSocket, DB, InfluxDB & Cache")

	if influxClient == nil {
		influxClient = newInfluxClientFromConfig(cfg, log)
	}

	return &mqttUseCase{
		client:        client,
		cfg:           cfg,
		logger:        log,
		subscriptions: make(map[string]subscriptionInfo),
		wsBroadcaster: wsBroadcaster,
		db:            db,
		alarmLogRepo:  alarmLogRepo,
		influxClient:  influxClient,
		cache:         cache,
	}, nil
}

func NewMQTTUseCaseWithClient(client mqttPkg.Client, cfg *config.Config, log logger.Logger) (MQTTUseCase, error) {
	if !client.IsConnected() {
		return nil, fmt.Errorf("MQTT client not connected")
	}
	log.Info("MQTT usecase using existing client")

	influxClient := newInfluxClientFromConfig(cfg, log)

	return &mqttUseCase{
		client:        client,
		cfg:           cfg,
		logger:        log,
		subscriptions: make(map[string]subscriptionInfo),
		influxClient:  influxClient,
	}, nil
}

// newInfluxClientFromConfig สร้าง InfluxDB client จาก config (INFLUXDB_URL / INFLUXDB_TOKEN / INFLUXDB_ORG / INFLUXDB_BUCKET)
// เพื่อให้ MQTT เขียนข้อมูลไปยัง InfluxDB ได้แม้ไม่มี client ส่งมาจากภายนอก
func newInfluxClientFromConfig(cfg *config.Config, log logger.Logger) influxdb.Client {
	if cfg == nil || cfg.InfluxDB.URL == "" || cfg.InfluxDB.Token == "" {
		log.Warn("MQTT: InfluxDB config missing (INFLUXDB_URL / INFLUXDB_TOKEN) – Influx writes skipped")
		return nil
	}
	ic, err := influxdb.NewInfluxClient(cfg, log)
	if err != nil {
		log.Warnf("MQTT: failed to create InfluxDB client from config: %v", err)
		return nil
	}
	log.Info("MQTT: created InfluxDB client from config (INFLUXDB_TOKEN)")
	return ic
}

func (u *mqttUseCase) waitForConnection(timeout time.Duration) error {
	if u.client.IsConnected() {
		return nil
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if u.client.IsConnected() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("MQTT client not connected after %v", timeout)
}

func (u *mqttUseCase) Publish(ctx context.Context, req *presenter.PublishRequest) error {
	if err := u.waitForConnection(5 * time.Second); err != nil {
		return err
	}
	u.logger.Infof("Publishing to topic %s", req.Topic)
	return u.client.Publish(req.Topic, req.QoS, req.Retained, req.Payload)
}

func extractRoomFromTopic(topic string) string {
	parts := strings.Split(topic, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "default"
}

func (u *mqttUseCase) Subscribe(ctx context.Context, req *presenter.SubscribeRequest) error {
	if err := u.waitForConnection(5 * time.Second); err != nil {
		return err
	}
	u.logger.Infof("Subscribing to topic %s", req.Topic)

	u.subsMu.RLock()
	_, exists := u.subscriptions[req.Topic]
	u.subsMu.RUnlock()
	if exists {
		u.logger.Warnf("Already subscribed to %s", req.Topic)
		return nil
	}

	handler := func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		u.logger.Infof("Received on %s: %s", msg.Topic(), payload)

		if u.wsBroadcaster != nil {
			room := extractRoomFromTopic(msg.Topic())
			data := map[string]interface{}{
				"topic":   msg.Topic(),
				"payload": payload,
			}
			u.wsBroadcaster.BroadcastToRoom(room, "mqtt", data)
		}
		go u.saveToInflux(msg.Topic(), msg.Payload())
		if containsAlarmKeyword(msg.Topic()) {
			// TODO: parse payload and call helpers.AlarmDetailValidate
		}
	}

	if err := u.client.Subscribe(req.Topic, req.QoS, handler); err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	u.subsMu.Lock()
	u.subscriptions[req.Topic] = subscriptionInfo{QoS: req.QoS}
	u.subsMu.Unlock()
	return nil
}

func (u *mqttUseCase) Unsubscribe(ctx context.Context, req *presenter.UnsubscribeRequest) error {
	if err := u.waitForConnection(5 * time.Second); err != nil {
		return err
	}
	u.logger.Infof("Unsubscribing from %s", req.Topic)

	u.subsMu.RLock()
	_, exists := u.subscriptions[req.Topic]
	u.subsMu.RUnlock()
	if !exists {
		u.logger.Warnf("Topic %s not subscribed", req.Topic)
		return nil
	}

	if err := u.client.Unsubscribe(req.Topic); err != nil {
		return fmt.Errorf("unsubscribe failed: %w", err)
	}

	u.subsMu.Lock()
	delete(u.subscriptions, req.Topic)
	u.subsMu.Unlock()
	return nil
}

func (u *mqttUseCase) GetSubscriptions() []string {
	u.subsMu.RLock()
	defer u.subsMu.RUnlock()
	topics := make([]string, 0, len(u.subscriptions))
	for topic := range u.subscriptions {
		topics = append(topics, topic)
	}
	return topics
}

func (u *mqttUseCase) SaveAlarmLog(alarmStatus int, subject, content, deviceID string, valueData interface{}) error {
	if u.alarmLogRepo == nil {
		return nil
	}
	now := time.Now()
	log := &models.AlarmLog{
		AlarmStatus: alarmStatus,
		Subject:     subject,
		Content:     content,
		DeviceID:    deviceID,
		DataAlarm:   fmt.Sprintf("%v", valueData),
		Date:        now.Format("2006-01-02"),
		Time:        now.Format("15:04:05"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return u.alarmLogRepo.Create(log)
}

func containsAlarmKeyword(topic string) bool {
	return strings.Contains(topic, "alarm") || strings.HasSuffix(topic, "/alarm")
}

func (u *mqttUseCase) IsConnected() bool {
	return u.client.IsConnected()
}

func (u *mqttUseCase) GetTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
	u.logger.Infof("GetTopic: waiting for connection...")
	if err := u.waitForConnection(timeout); err != nil {
		u.logger.Errorf("GetTopic: connection failed: %v", err)
		return nil, err
	}
	u.logger.Infof("GetTopic: calling client.GetDataFromTopic for %s", topic)
	data, err := u.client.GetDataFromTopic(ctx, topic, timeout)
	if err != nil {
		u.logger.Errorf("GetTopic: client.GetDataFromTopic error: %v", err)
	} else {
		u.logger.Infof("GetTopic: received %d bytes", len(data))
	}
	return data, err
}

func (u *mqttUseCase) DeviceControl(ctx context.Context, req *presenter.DeviceControlRequest) (*presenter.DeviceControlResponse, error) {
	startTime := time.Now()
	if req.Topic == "" || req.Message == "" {
		return nil, fmt.Errorf("topic and message are required")
	}

	cacheKey := fmt.Sprintf("device_ctrl:%s:%s", req.Topic, req.Message)
	if u.cache != nil {
		var cachedResp presenter.DeviceControlResponse
		if err := u.cache.Get(ctx, cacheKey, &cachedResp); err == nil {
			u.logger.Infof("DeviceControl: cache HIT for %s", cacheKey)
			cachedResp.FetchDurationMs = time.Since(startTime).Milliseconds()
			cachedResp.From = "cache"
			return &cachedResp, nil
		}
		u.logger.Debugf("DeviceControl: cache MISS for %s", cacheKey)
	}

	if err := u.waitForConnection(2 * time.Second); err != nil {
		return u.fallbackDeviceControlResponse(req, startTime, err.Error()), nil
	}

	topicData := strings.Replace(req.Topic, "CONTROL", "DATA", 1)
	if topicData == req.Topic {
		topicData = req.Topic + "_DATA"
	}
	u.logger.Infof("DeviceControl: publishing to %s with message %s", req.Topic, req.Message)

	if err := u.client.Publish(req.Topic, 1, false, req.Message); err != nil {
		u.logger.Errorf("DeviceControl: publish failed: %v", err)
		return u.fallbackDeviceControlResponse(req, startTime, err.Error()), nil
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	payload, err := u.client.GetDataFromTopic(ctxTimeout, topicData, 3*time.Second)
	if err != nil {
		u.logger.Errorf("DeviceControl: failed to get response: %v", err)
		return u.fallbackDeviceControlResponse(req, startTime, err.Error()), nil
	}

	var parts []string
	if len(payload) > 0 {
		parts = strings.Split(string(payload), ",")
	}

	msgLower := strings.ToLower(req.Message)
	var status int
	var statusMsg string
	switch msgLower {
	case "1", "on", "a1", "b1", "c1", "d1", "e1", "f1", "g1":
		status = 1
		statusMsg = "ON"
	default:
		status = 0
		statusMsg = "OFF"
	}
	TimeLoc := helpers.GetTimeLocation()
	resp := &presenter.DeviceControlResponse{
		StatusCode:      200,
		Code:            200,
		TopicControl:    req.Topic,
		TopicData:       topicData,
		MessageSent:     req.Message,
		Payload:         string(payload),
		Data:            parts,
		Status:          status,
		StatusMsg:       statusMsg,
		Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
		Message:         fmt.Sprintf("Control sent to %s, response received", req.Topic),
		MessageTh:       fmt.Sprintf("ส่งคำสั่งไปยัง %s และได้รับข้อมูลตอบกลับ", req.Topic),
		From:            "mqtt",
		FetchDurationMs: time.Since(startTime).Milliseconds(),
	}

	if u.cache != nil {
		_ = u.cache.Set(ctx, cacheKey, resp, 5*time.Second)
	}

	return resp, nil
}

func (u *mqttUseCase) fallbackDeviceControlResponse(req *presenter.DeviceControlRequest, startTime time.Time, errorDetail string) *presenter.DeviceControlResponse {
	msgLower := strings.ToLower(req.Message)
	var status int
	var statusMsg string
	switch msgLower {
	case "1", "on", "a1", "b1", "c1", "d1", "e1", "f1", "g1":
		status = 1
		statusMsg = "ON"
	default:
		status = 0
		statusMsg = "OFF"
	}

	TimeLoc := helpers.GetTimeLocation()
	fallbackPayload := "0,0,0,0,0,0,0,0,0"
	fallbackData := strings.Split(fallbackPayload, ",")
	return &presenter.DeviceControlResponse{
		StatusCode:      200,
		Code:            200,
		TopicControl:    req.Topic,
		TopicData:       strings.Replace(req.Topic, "CONTROL", "DATA", 1),
		MessageSent:     req.Message,
		Payload:         fallbackPayload,
		Data:            fallbackData,
		Status:          status,
		StatusMsg:       statusMsg,
		Timestamp:       helpers.TimeConvertermas(time.Now().In(TimeLoc)),
		Message:         "fallback (MQTT timeout/error)",
		MessageTh:       "fallback (MQTT หมดเวลาหรือผิดพลาด)",
		From:            "fallback",
		FetchDurationMs: time.Since(startTime).Milliseconds(),
		ErrorDetail:     errorDetail,
	}
}

func (u *mqttUseCase) saveToInflux(topic string, payload []byte) {
	if u.influxClient == nil {
		return
	}
	if v, ok := u.influxClient.(*influxdb.InfluxClient); ok && v == nil {
		return
	}
	tags := map[string]string{"topic": topic}
	fields := map[string]interface{}{
		"payload": string(payload),
		"length":  len(payload),
	}
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err == nil {
		for k, v := range data {
			if _, ok := v.(float64); ok {
				fields[k] = v
			}
		}
	}
	_ = u.influxClient.WritePoint("mqtt_messages", tags, fields, time.Now())
}
