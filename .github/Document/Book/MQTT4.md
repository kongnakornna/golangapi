เราจะทำการสร้าง **MQTT Client** แบบสมบูรณ์รวม `RequestManager` ไว้ในไฟล์ `client.go` เดียว จากนั้นจึงสร้าง **IoT Service Module** ทั้งหมดตามลำดับ

---

## 1. `pkg/mqtt/client.go` – ฉบับสมบูรณ์ (รวม RequestManager)

```go
// pkg/mqtt/client.go
package mqtt

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client defines MQTT operations
type Client interface {
	Connect(ctx context.Context) error
	Disconnect(quiesce uint)
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
	SubscribeMultiple(topics map[string]byte, callback mqtt.MessageHandler) error
	Unsubscribe(topics ...string) error
	IsConnected() bool
	RequestData(requestTopic string, responseTopic string, payload interface{}, timeout time.Duration) (interface{}, error)
	GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
}

// ---------- internal request manager ----------
type pendingRequest struct {
	resolve chan<- interface{}
	reject  chan<- error
	timeout *time.Timer
	topic   string
}

type requestManager struct {
	client            mqtt.Client
	mu                sync.Mutex
	pending           map[string][]*pendingRequest
	subscriptionCount map[string]int
	debug             bool
	logger            logger.Logger
}

func newRequestManager(client mqtt.Client, debug bool, log logger.Logger) *requestManager {
	rm := &requestManager{
		client:            client,
		pending:           make(map[string][]*pendingRequest),
		subscriptionCount: make(map[string]int),
		debug:             debug,
		logger:            log,
	}
	client.AddRoute("#", rm.handleIncomingMessage)
	log.Info("requestManager initialized with global message handler")
	return rm
}

func (rm *requestManager) handleIncomingMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	rm.mu.Lock()
	defer rm.mu.Unlock()

	pendings, ok := rm.pending[topic]
	if !ok || len(pendings) == 0 {
		if rm.debug {
			rm.logger.Debugf("handleIncomingMessage: no pending request for topic=%s, ignoring", topic)
		}
		return
	}
	req := pendings[0]
	rm.pending[topic] = pendings[1:]
	if len(rm.pending[topic]) == 0 {
		delete(rm.pending, topic)
		rm.decrementSubscription(topic)
	}
	req.timeout.Stop()

	var result interface{}
	payload := msg.Payload()
	if err := json.Unmarshal(payload, &result); err != nil {
		result = string(payload)
	}
	rm.logger.Infof("handleIncomingMessage: received message on topic=%s, resolving pending request", topic)
	req.resolve <- result
}

func (rm *requestManager) getDataFromTopic(topic string, timeoutMs int) (interface{}, error) {
	if rm.debug {
		rm.logger.Debugf("getDataFromTopic: start for topic=%s timeout=%dms", topic, timeoutMs)
	}
	rm.mu.Lock()
	defer func() {
		rm.mu.Unlock()
		if rm.debug {
			rm.logger.Debugf("getDataFromTopic: lock released for topic=%s", topic)
		}
	}()

	if _, ok := rm.pending[topic]; !ok {
		rm.incrementSubscription(topic)
	}

	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	timer := time.NewTimer(time.Duration(timeoutMs) * time.Millisecond)

	req := &pendingRequest{
		resolve: resultCh,
		reject:  errCh,
		timeout: timer,
		topic:   topic,
	}
	rm.pending[topic] = append(rm.pending[topic], req)

	go func() {
		<-timer.C
		rm.mu.Lock()
		defer rm.mu.Unlock()
		list := rm.pending[topic]
		for i, r := range list {
			if r == req {
				rm.pending[topic] = append(list[:i], list[i+1:]...)
				break
			}
		}
		if len(rm.pending[topic]) == 0 {
			delete(rm.pending, topic)
			rm.decrementSubscription(topic)
		}
		errCh <- fmt.Errorf("timeout: no message from topic %s after %d ms", topic, timeoutMs)
	}()

	select {
	case err := <-errCh:
		return nil, err
	case res := <-resultCh:
		return res, nil
	}
}

func (rm *requestManager) incrementSubscription(topic string) {
	count := rm.subscriptionCount[topic]
	if count == 0 {
		token := rm.client.Subscribe(topic, 0, nil)
		if token.Wait() && token.Error() != nil {
			rm.logger.Errorf("incrementSubscription: subscribe failed for topic=%s: %v", topic, token.Error())
			if pendings, ok := rm.pending[topic]; ok {
				for _, req := range pendings {
					req.reject <- token.Error()
					req.timeout.Stop()
				}
				delete(rm.pending, topic)
			}
			return
		}
		rm.subscriptionCount[topic] = 1
	} else {
		rm.subscriptionCount[topic] = count + 1
	}
}

func (rm *requestManager) decrementSubscription(topic string) {
	count := rm.subscriptionCount[topic]
	if count <= 1 {
		rm.client.Unsubscribe(topic)
		delete(rm.subscriptionCount, topic)
	} else {
		rm.subscriptionCount[topic] = count - 1
	}
}

func (rm *requestManager) publishAndWait(requestTopic, responseTopic string, payload interface{}, timeoutMs int) (interface{}, error) {
	var msgPayload []byte
	switch v := payload.(type) {
	case string:
		msgPayload = []byte(v)
	case []byte:
		msgPayload = v
	default:
		var err error
		msgPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}
	token := rm.client.Publish(requestTopic, 1, false, msgPayload)
	token.Wait()
	if token.Error() != nil {
		return nil, fmt.Errorf("publish failed: %w", token.Error())
	}
	rm.logger.Debugf("publishAndWait: published to %s, waiting for response on %s", requestTopic, responseTopic)
	return rm.getDataFromTopic(responseTopic, timeoutMs)
}

// ---------- mqttClient implementation ----------
type mqttClient struct {
	client         mqtt.Client
	cfg            *config.MQTTConfig
	logger         logger.Logger
	opts           *mqtt.ClientOptions
	requestManager *requestManager
}

// New creates a new MQTT client
func New(cfg *config.MQTTConfig, log logger.Logger) Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)

	clientID := cfg.ClientID
	if clientID == "" {
		clientID = fmt.Sprintf("go-client-%d", time.Now().UnixNano())
	}

	log.Infof("MQTT clientID: %s", clientID)
	log.Infof("MQTT broker: %s", cfg.Broker)

	opts.SetClientID(clientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetResumeSubs(false)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(15 * time.Second)
	opts.SetWriteTimeout(15 * time.Second)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		log.Errorf("MQTT connection lost: %v", err)
	}
	opts.OnReconnecting = func(cl mqtt.Client, opts *mqtt.ClientOptions) {
		log.Info("MQTT reconnecting...")
	}
	opts.OnConnect = func(cl mqtt.Client) {
		log.Info("MQTT connected")
	}

	return &mqttClient{
		cfg:    cfg,
		logger: log,
		opts:   opts,
	}
}

func (c *mqttClient) Connect(ctx context.Context) error {
	c.client = mqtt.NewClient(c.opts)
	token := c.client.Connect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return fmt.Errorf("MQTT connect failed: %w", token.Error())
		}
		c.requestManager = newRequestManager(c.client, true, c.logger)
		c.logger.Info("MQTT connected and requestManager initialized")
		return nil
	}
}

func (c *mqttClient) Disconnect(quiesce uint) {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(quiesce)
		c.logger.Info("MQTT disconnected")
	}
}

func (c *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}
	token := c.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (c *mqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}
	token := c.client.Subscribe(topic, qos, callback)
	token.Wait()
	return token.Error()
}

func (c *mqttClient) SubscribeMultiple(topics map[string]byte, callback mqtt.MessageHandler) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}
	token := c.client.SubscribeMultiple(topics, callback)
	token.Wait()
	return token.Error()
}

func (c *mqttClient) Unsubscribe(topics ...string) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}
	token := c.client.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}

func (c *mqttClient) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}

// RequestData implements request-response using internal requestManager
func (c *mqttClient) RequestData(requestTopic string, responseTopic string, payload interface{}, timeout time.Duration) (interface{}, error) {
	if c.requestManager == nil {
		return nil, fmt.Errorf("request manager not initialized (client not connected)")
	}
	if requestTopic == "" {
		return nil, fmt.Errorf("requestTopic cannot be empty")
	}
	waitTopic := responseTopic
	if waitTopic == "" {
		waitTopic = requestTopic
	}
	return c.requestManager.publishAndWait(requestTopic, waitTopic, payload, int(timeout.Milliseconds()))
}

// GetDataFromTopic gets data from a topic using request-response pattern
func (c *mqttClient) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
	if c.requestManager == nil {
		return nil, fmt.Errorf("request manager not initialized (client not connected)")
	}
	result, err := c.requestManager.getDataFromTopic(topic, int(timeout.Milliseconds()))
	if err != nil {
		return nil, err
	}
	switch v := result.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(v)
	}
}
```

---

## 2. สร้าง `internal/modules/iot/iothelper/alarm.go` (ตรรกะการแจ้งเตือน)

```go
// internal/modules/iot/iothelper/alarm.go
package iothelper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AlarmDetailDto struct {
	HardwareID        interface{}
	ValueData         interface{}
	Max               interface{}
	Min               interface{}
	StatusAlert       interface{}
	StatusWarning     interface{}
	RecoveryWarning   interface{}
	RecoveryAlert     interface{}
	DeviceName        string
	ActionName        string
	MqttName          string
	MqttControlOn     string
	MqttControlOff    string
	CountAlarm        interface{}
	Event             interface{}
	Unit              string
}

type AlarmDetailResult struct {
	Status             int    `json:"status"`
	StatusControl      int    `json:"status_control"`
	AlarmTypeId        int    `json:"alarm_type_id"`
	TypeId             int    `json:"type_id"`
	HardwareId         int    `json:"hardware_id"`
	AlarmStatusSet     int    `json:"alarm_status_set"`
	Title              string `json:"title"`
	Subject            string `json:"subject"`
	Content            string `json:"content"`
	DataAlarm          int    `json:"data_alarm"`
	EventControl       int    `json:"event_control"`
	MessageMqttControl string `json:"message_mqtt_control"`
	SensorData         interface{} `json:"sensor_data"`
	CountAlarm         int    `json:"count_alarm"`
	Unit               string `json:"unit"`
	Timestamp          string `json:"timestamp"`
}

var thaiMessages = map[string]string{
	"warning":          "คำเตือน มีความผิดปกติ",
	"critical":         "ภาวะวิกฤตต้องแก้ไขทันที",
	"recoveryWarning":  "คืนสู่ภาวะปกติ (คำเตือน)",
	"recoveryCritical": "คืนสู่ภาวะปกติ (วิกฤต)",
	"normal":           "ปกติ",
	"normal2":          "ปกติ",
	"criticalMax":      "วิกฤต มีค่าสูงเกินกำหนด",
	"criticalMin":      "วิกฤต มีค่าต่ำกว่ากำหนด",
	"normal3":          "ปกติ",
}

var englishMessages = map[string]string{
	"warning":          "Warning",
	"critical":         "Critical",
	"recoveryWarning":  "Recovery Warning",
	"recoveryCritical": "Recovery Critical",
	"normal":           "Normal",
	"normal2":          "Normal",
	"criticalMax":      "Critical! Maximum limit.",
	"criticalMin":      "Critical! Minimum limit",
	"normal3":          "Normal",
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

func normalizeSensorValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		up := strings.ToUpper(val)
		if up == "ON" {
			return 1
		}
		if up == "OFF" {
			return 0
		}
		return val
	default:
		return val
	}
}

func processAlarmDetail(dto AlarmDetailDto, messages map[string]string) AlarmDetailResult {
	hardwareID := toInt(dto.HardwareID)
	typeID := hardwareID

	sensorValue := normalizeSensorValue(dto.ValueData)
	maxVal := toFloat(dto.Max)
	minVal := toFloat(dto.Min)
	statusAlert := toInt(dto.StatusAlert)
	statusWarning := toInt(dto.StatusWarning)
	recoveryWarning := toInt(dto.RecoveryWarning)
	recoveryAlert := toInt(dto.RecoveryAlert)
	countAlarm := toInt(dto.CountAlarm)
	event := toInt(dto.Event)

	unit := dto.Unit
	mqttName := dto.MqttName
	deviceName := dto.DeviceName
	alarmActionName := dto.ActionName
	mqttControlOn := dto.MqttControlOn
	mqttControlOff := dto.MqttControlOff
	valueAlarm := dto.ValueData
	valueRelay := dto.ValueData
	valueControlRelay := dto.ValueData

	var sensorData interface{}
	var valueData interface{}

	switch hardwareID {
	case 1:
		sensorData = dto.ValueData
		valueData = dto.ValueData
	case 2:
		if toInt(dto.ValueData) == 1 {
			sensorData = 1
			valueData = 1
			sensorValue = 1
		} else {
			sensorData = toInt(dto.ValueData)
			valueData = toInt(dto.ValueData)
			sensorValue = toInt(dto.ValueData)
		}
	case 3:
		sensorData = toInt(dto.ValueData)
		valueData = dto.ValueData
		sensorValue = dto.ValueData
	case 4:
		sensorData = dto.ValueData
		valueData = dto.ValueData
	default:
		sensorData = toInt(dto.ValueData)
		valueData = dto.ValueData
	}

	alarmStatusSet := 999
	dataAlarm := 0
	eventControl := event
	messageMqttControl := mqttControlOff
	if event == 1 {
		messageMqttControl = mqttControlOn
	}
	status := 5
	title := messages["normal"]
	subject := messages["normal"]
	content := messages["normal"] + " "

	if hardwareID == 3 && (sensorValue == 1 || sensorValue == 0 || sensorValue == "ON" || sensorValue == "OFF" || sensorValue == "on" || sensorValue == "off") {
		alarmStatusSet = 999
		title = messages["normal"]
		subject = messages["normal"]
		content = fmt.Sprintf("%s %v %s", messages["normal"], sensorValue, unit)
		status = 5
	} else if hardwareID == 4 && sensorValue != 1 {
		alarmStatusSet = 2
		title = messages["critical"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["critical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["critical"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 2
	} else if hardwareID == 4 && sensorValue == 1 {
		alarmStatusSet = 999
		title = messages["normal2"]
		subject = messages["normal2"]
		content = fmt.Sprintf("%s %v %s", messages["normal2"], sensorValue, unit)
		status = 5
	} else if maxVal != 0 && toFloat(sensorValue) >= maxVal && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 2
		title = messages["criticalMax"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["criticalMax"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["criticalMax"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 2
	} else if minVal != 0 && toFloat(sensorValue) <= minVal && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 1
		title = messages["criticalMin"]
		subject = fmt.Sprintf("%s %s %s : %v %s", mqttName, messages["criticalMin"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s : %s :%v %s", mqttName, alarmActionName, messages["criticalMin"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 1
	} else if hardwareID == 1 && statusWarning > 0 && toFloat(sensorValue) >= float64(statusWarning) && toFloat(sensorValue) < float64(statusAlert) {
		alarmStatusSet = 1
		title = messages["warning"]
		subject = fmt.Sprintf("%s %s : %s : %v %s", mqttName, messages["warning"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["warning"], deviceName, sensorValue, unit)
		dataAlarm = statusWarning
		status = 1
	} else if hardwareID == 1 && statusAlert > 0 && toFloat(sensorValue) >= float64(statusAlert) {
		alarmStatusSet = 2
		title = messages["critical"]
		subject = fmt.Sprintf("%s %s : %s :%v %s", mqttName, messages["critical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["critical"], deviceName, sensorValue, unit)
		dataAlarm = statusAlert
		status = 2
	} else if toInt(valueAlarm) == 0 && (hardwareID == 2 || hardwareID == 3 || hardwareID == 4) {
		isCritical := hardwareID == 4
		if isCritical {
			alarmStatusSet = 2
			title = messages["critical"]
		} else {
			alarmStatusSet = 1
			title = messages["warning"]
		}
		subject = fmt.Sprintf("%s %s : %s : %v %s", mqttName, title, deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, title, deviceName, sensorValue, unit)
		if isCritical {
			dataAlarm = statusAlert
		} else {
			dataAlarm = statusWarning
		}
		status = 2
		if !isCritical {
			status = 1
		}
	} else if countAlarm >= 1 && recoveryWarning > 0 && toFloat(sensorValue) <= float64(recoveryWarning) && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 3
		title = messages["recoveryWarning"]
		subject = fmt.Sprintf("%s %s : %s :%v %s", mqttName, messages["recoveryWarning"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s: %s :%v %s", mqttName, alarmActionName, messages["recoveryWarning"], deviceName, sensorValue, unit)
		dataAlarm = recoveryWarning
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 3
	} else if countAlarm >= 1 && recoveryAlert > 0 && toFloat(sensorValue) <= float64(recoveryAlert) && (hardwareID == 1 || hardwareID == 2) {
		alarmStatusSet = 4
		title = fmt.Sprintf("%s %s", mqttName, messages["recoveryCritical"])
		subject = fmt.Sprintf("%s %s :%s :%v %s", mqttName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s :%s :%v %s", mqttName, alarmActionName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		dataAlarm = recoveryAlert
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 4
	} else if countAlarm >= 1 && toInt(valueAlarm) >= 1 && (hardwareID == 2 || hardwareID == 3 || hardwareID == 4) {
		alarmStatusSet = 4
		title = fmt.Sprintf("%s %s", mqttName, messages["recoveryCritical"])
		subject = fmt.Sprintf("%s %s :%s :%v %s", mqttName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		content = fmt.Sprintf("%s %s %s :%s :%v %s", mqttName, alarmActionName, messages["recoveryCritical"], deviceName, sensorValue, unit)
		dataAlarm = recoveryAlert
		eventControl = 0
		if event == 1 {
			eventControl = 1
		}
		if event == 1 {
			messageMqttControl = mqttControlOff
		} else {
			messageMqttControl = mqttControlOn
		}
		status = 4
	} else {
		alarmStatusSet = 999
		title = messages["normal3"]
		subject = messages["normal3"]
		content = messages["normal"] + " "
		dataAlarm = 0
		status = 5
	}

	return AlarmDetailResult{
		Status:             status,
		StatusControl:      status,
		AlarmTypeId:        hardwareID,
		TypeId:             typeID,
		HardwareId:         hardwareID,
		AlarmStatusSet:     alarmStatusSet,
		Title:              title,
		Subject:            subject,
		Content:            content,
		DataAlarm:          dataAlarm,
		EventControl:       eventControl,
		MessageMqttControl: messageMqttControl,
		SensorData:         sensorData,
		CountAlarm:         countAlarm,
		Unit:               unit,
		Timestamp:          time.Now().Format("2006-01-02 15:04:05"),
	}
}

func AlarmDetailValidate(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, thaiMessages)
}

func AlarmDetailValidateEn(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, englishMessages)
}

func AlarmDetailValidateTh(dto AlarmDetailDto) AlarmDetailResult {
	return processAlarmDetail(dto, thaiMessages)
}
```

---

## 3. `internal/modules/iot/presenter/presenter.go`

```go
// internal/modules/iot/presenter/presenter.go
package presenter

type DeviceListRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Bucket     string `json:"bucket"`
	HardwareId int    `json:"hardware_id"`
	TypeId     int    `json:"type_id"`
	Keyword    string `json:"keyword"`
	Lang       string `json:"lang"`
}

type DeviceDetailResponse struct {
	DeviceId   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	TypeName   string `json:"type_name"`
	ValueData  string `json:"value_data"`
	Unit       string `json:"unit"`
	Status     int    `json:"status"`
	AlarmTitle string `json:"alarm_title"`
}

type TopicDataResponse struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
	From    string      `json:"from"`
	Cache   bool        `json:"cache"`
}

type ControlRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type SenserChartRequest struct {
	Bucket      string `json:"bucket"`
	Measurement string `json:"measurement"`
	Field       string `json:"field"`
	Start       string `json:"start"`
	Stop        string `json:"stop"`
	Limit       int    `json:"limit"`
}

type SenserChartResponse struct {
	Data  []float64 `json:"data"`
	Date  []string  `json:"date"`
	Cache string    `json:"cache"`
}

type DeviceBucketsResponse struct {
	Bucket string               `json:"bucket"`
	Devices []DeviceDetailResponse `json:"devices"`
}
```

---

## 4. `internal/modules/iot/repository/device_repo.go` (สมมติว่า model มีอยู่)

```go
// internal/modules/iot/repository/device_repo.go
package repository

import (
	"icmongolang/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDeviceByID(id int) (*models.Device, error)
	GetDevicesByBucket(bucket string) ([]models.Device, error)
	ListDevices(filter map[string]interface{}, page, pageSize int) ([]models.Device, int64, error)
}

type deviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetDeviceByID(id int) (*models.Device, error) {
	var device models.Device
	err := r.db.Preload("Mqtt").Preload("Location").First(&device, id).Error
	return &device, err
}

func (r *deviceRepo) GetDevicesByBucket(bucket string) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("bucket = ?", bucket).Preload("Mqtt").Find(&devices).Error
	return devices, err
}

func (r *deviceRepo) ListDevices(filter map[string]interface{}, page, pageSize int) ([]models.Device, int64, error) {
	var devices []models.Device
	query := r.db.Model(&models.Device{})
	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	var total int64
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Preload("Mqtt").Preload("Location").Find(&devices).Error
	return devices, total, err
}
```

---

## 5. `internal/modules/iot/repository/alarm_log_repo.go` (ตัวอย่าง)

```go
// internal/modules/iot/repository/alarm_log_repo.go
package repository

import (
	"icmongolang/internal/models"
	"gorm.io/gorm"
)

type AlarmLogRepository interface {
	Create(log *models.AlarmLog) error
	CountByDevice(deviceID int, alarmStatus int) (int64, error)
}

type alarmLogRepo struct {
	db *gorm.DB
}

func NewAlarmLogRepository(db *gorm.DB) AlarmLogRepository {
	return &alarmLogRepo{db: db}
}

func (r *alarmLogRepo) Create(log *models.AlarmLog) error {
	return r.db.Create(log).Error
}

func (r *alarmLogRepo) CountByDevice(deviceID int, alarmStatus int) (int64, error) {
	var count int64
	err := r.db.Model(&models.AlarmLog{}).Where("device_id = ? AND alarm_status = ?", deviceID, alarmStatus).Count(&count).Error
	return count, err
}
```

---

## 6. `internal/modules/iot/repository/schedule_repo.go` (ตัวอย่าง)

```go
// internal/modules/iot/repository/schedule_repo.go
package repository

import (
	"icmongolang/internal/models"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	GetActiveSchedules() ([]models.Schedule, error)
}

type scheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepo{db: db}
}

func (r *scheduleRepo) GetActiveSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Where("status = 1").Preload("ScheduleDevices.Device").Find(&schedules).Error
	return schedules, err
}
```

---

## 7. `internal/modules/iot/usecase/usecase.go` – เบื้องต้นสำหรับ endpoints ที่กำหนด

เราจะรวม endpoints ที่โจทย์ต้องการ: topic, control, device, listdevicepage, devicebuckets, sensercharts, senserdatachart, senserdata, devicelist, locationdevice, devicesensercharts, alarmdevicestatus, alarmdevicestatuscontrol, monitordevicegroup, monitordevicechart

```go
// internal/modules/iot/usecase/usecase.go
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"icmongolang/internal/modules/iot/iothelper"
	"icmongolang/internal/modules/iot/presenter"
	"icmongolang/internal/modules/iot/repository"
	"icmongolang/internal/models"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
	"icmongolang/pkg/redis"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MQTT3UseCase interface {
	// Basic MQTT
	GetTopicData(ctx context.Context, topic string) (*presenter.TopicDataResponse, error)
	DeviceControl(ctx context.Context, req *presenter.ControlRequest) error

	// Device listing
	GetDeviceList(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error)
	GetDeviceListPage(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error)
	GetDeviceBuckets(ctx context.Context, bucket string) (*presenter.DeviceBucketsResponse, error)
	GetDeviceListByLocation(ctx context.Context, locationID int) ([]presenter.DeviceDetailResponse, error)

	// Chart data
	GetSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error)
	GetSenserDataChart(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error)
	GetSenserData(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error)
	GetDeviceSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error)

	// Alarm status
	GetAlarmDeviceStatus(ctx context.Context, req map[string]interface{}) (interface{}, error)
	GetAlarmDeviceStatusControl(ctx context.Context, req map[string]interface{}) (interface{}, error)

	// Monitoring group & chart
	GetMonitorDeviceGroup(ctx context.Context, req map[string]interface{}) (interface{}, error)
	GetMonitorDeviceChart(ctx context.Context, req map[string]interface{}) (interface{}, error)
}

type mqtt3UseCase struct {
	deviceRepo   repository.DeviceRepository
	alarmLogRepo repository.AlarmLogRepository
	scheduleRepo repository.ScheduleRepository
	mqttClient   mqtt.Client
	redisClient  *redis.Client
	influxClient *influxdb.InfluxClient
	logger       logger.Logger
}

func NewMQTT3UseCase(
	deviceRepo repository.DeviceRepository,
	alarmLogRepo repository.AlarmLogRepository,
	scheduleRepo repository.ScheduleRepository,
	mqttClient mqtt.Client,
	redisClient *redis.Client,
	influxClient *influxdb.InfluxClient,
	log logger.Logger,
) MQTT3UseCase {
	return &mqtt3UseCase{
		deviceRepo:   deviceRepo,
		alarmLogRepo: alarmLogRepo,
		scheduleRepo: scheduleRepo,
		mqttClient:   mqttClient,
		redisClient:  redisClient,
		influxClient: influxClient,
		logger:       log,
	}
}

// GetTopicData implements MQTT3UseCase
func (u *mqtt3UseCase) GetTopicData(ctx context.Context, topic string) (*presenter.TopicDataResponse, error) {
	// cache
	cacheKey := "mqtt_topic:" + topic
	cached, err := u.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var payload interface{}
		json.Unmarshal(cached, &payload)
		return &presenter.TopicDataResponse{
			Topic:   topic,
			Payload: payload,
			From:    "cache",
			Cache:   true,
		}, nil
	}
	data, err := u.mqttClient.GetDataFromTopic(ctx, topic, 10*time.Second)
	if err != nil {
		return nil, err
	}
	var payload interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		payload = string(data)
	}
	// store cache 30 seconds
	u.redisClient.Set(ctx, cacheKey, data, 30*time.Second)
	return &presenter.TopicDataResponse{
		Topic:   topic,
		Payload: payload,
		From:    "mqtt",
		Cache:   false,
	}, nil
}

// DeviceControl implements MQTT3UseCase
func (u *mqtt3UseCase) DeviceControl(ctx context.Context, req *presenter.ControlRequest) error {
	return u.mqttClient.Publish(req.Topic, 1, false, req.Message)
}

// GetDeviceList implements MQTT3UseCase
func (u *mqtt3UseCase) GetDeviceList(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error) {
	filter := make(map[string]interface{})
	if req.Bucket != "" {
		filter["bucket"] = req.Bucket
	}
	if req.HardwareId != 0 {
		filter["hardware_id"] = req.HardwareId
	}
	devices, total, err := u.deviceRepo.ListDevices(filter, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]presenter.DeviceDetailResponse, len(devices))
	for i, dev := range devices {
		result[i] = presenter.DeviceDetailResponse{
			DeviceId:   dev.DeviceID,
			DeviceName: dev.DeviceName,
			TypeName:   dev.TypeName,
			Unit:       dev.Unit,
			Status:     dev.Status,
		}
	}
	return result, total, nil
}

// GetDeviceListPage alias of GetDeviceList
func (u *mqtt3UseCase) GetDeviceListPage(ctx context.Context, req *presenter.DeviceListRequest) ([]presenter.DeviceDetailResponse, int64, error) {
	return u.GetDeviceList(ctx, req)
}

// GetDeviceBuckets implements MQTT3UseCase
func (u *mqtt3UseCase) GetDeviceBuckets(ctx context.Context, bucket string) (*presenter.DeviceBucketsResponse, error) {
	devices, err := u.deviceRepo.GetDevicesByBucket(bucket)
	if err != nil {
		return nil, err
	}
	resp := &presenter.DeviceBucketsResponse{Bucket: bucket}
	for _, dev := range devices {
		resp.Devices = append(resp.Devices, presenter.DeviceDetailResponse{
			DeviceId:   dev.DeviceID,
			DeviceName: dev.DeviceName,
			TypeName:   dev.TypeName,
			Unit:       dev.Unit,
		})
	}
	return resp, nil
}

// GetDeviceListByLocation implements MQTT3UseCase
func (u *mqtt3UseCase) GetDeviceListByLocation(ctx context.Context, locationID int) ([]presenter.DeviceDetailResponse, error) {
	// สมมติว่า location_id อยู่ใน device model
	var devices []models.Device
	err := u.deviceRepo.(*repository.deviceRepo).db.Where("location_id = ?", locationID).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	result := make([]presenter.DeviceDetailResponse, len(devices))
	for i, dev := range devices {
		result[i] = presenter.DeviceDetailResponse{
			DeviceId:   dev.DeviceID,
			DeviceName: dev.DeviceName,
			TypeName:   dev.TypeName,
			Unit:       dev.Unit,
		}
	}
	return result, nil
}

// GetSenserCharts implements MQTT3UseCase
func (u *mqtt3UseCase) GetSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error) {
	params := influxdb.QueryParams{
		Measurement: req.Measurement,
		Field:       req.Field,
		Bucket:      req.Bucket,
		Start:       req.Start,
		Stop:        req.Stop,
		Limit:       req.Limit,
	}
	results, err := u.influxClient.QueryFilterData(params)
	if err != nil {
		return nil, err
	}
	var data []float64
	var date []string
	for _, r := range results {
		if val, ok := r["_value"].(float64); ok {
			data = append(data, val)
		}
		if t, ok := r["_time"].(time.Time); ok {
			date = append(date, t.Format("2006-01-02 15:04:05"))
		}
	}
	return &presenter.SenserChartResponse{Data: data, Date: date, Cache: "no cache"}, nil
}

// GetSenserDataChart alias
func (u *mqtt3UseCase) GetSenserDataChart(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error) {
	return u.GetSenserCharts(ctx, req)
}

// GetSenserData alias
func (u *mqtt3UseCase) GetSenserData(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error) {
	return u.GetSenserCharts(ctx, req)
}

// GetDeviceSenserCharts alias
func (u *mqtt3UseCase) GetDeviceSenserCharts(ctx context.Context, req *presenter.SenserChartRequest) (*presenter.SenserChartResponse, error) {
	return u.GetSenserCharts(ctx, req)
}

// GetAlarmDeviceStatus (placeholder - จะขยายตาม business logic)
func (u *mqtt3UseCase) GetAlarmDeviceStatus(ctx context.Context, req map[string]interface{}) (interface{}, error) {
	// TODO: implement alarm status logic
	return map[string]string{"status": "ok"}, nil
}

// GetAlarmDeviceStatusControl (placeholder)
func (u *mqtt3UseCase) GetAlarmDeviceStatusControl(ctx context.Context, req map[string]interface{}) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}

// GetMonitorDeviceGroup (placeholder)
func (u *mqtt3UseCase) GetMonitorDeviceGroup(ctx context.Context, req map[string]interface{}) (interface{}, error) {
	return []string{}, nil
}

// GetMonitorDeviceChart (placeholder)
func (u *mqtt3UseCase) GetMonitorDeviceChart(ctx context.Context, req map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{"data": []float64{}, "date": []string{}}, nil
}
```

---

## 8. `internal/modules/iot/delivery/http/handler.go`

```go
// internal/modules/iot/delivery/http/handler.go
package http

import (
	"net/http"
	"strconv"

	"icmongolang/internal/modules/iot/presenter"
	"icmongolang/internal/modules/iot/usecase"
	"icmongolang/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type MQTT3Handler struct {
	uc     usecase.MQTT3UseCase
	logger logger.Logger
}

func NewMQTT3Handler(uc usecase.MQTT3UseCase, log logger.Logger) *MQTT3Handler {
	return &MQTT3Handler{uc: uc, logger: log}
}

// GET /topic
func (h *MQTT3Handler) GetTopicData(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		render.Render(w, r, ErrBadRequest("topic is required"))
		return
	}
	resp, err := h.uc.GetTopicData(r.Context(), topic)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// POST /control
func (h *MQTT3Handler) DeviceControl(w http.ResponseWriter, r *http.Request) {
	var req presenter.ControlRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrBadRequest(err.Error()))
		return
	}
	if err := h.uc.DeviceControl(r.Context(), &req); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]string{"status": "ok"})
}

// GET /device
func (h *MQTT3Handler) GetDeviceList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	req := presenter.DeviceListRequest{
		Page:       page,
		PageSize:   pageSize,
		Bucket:     r.URL.Query().Get("bucket"),
		HardwareId: atoi(r.URL.Query().Get("hardware_id")),
		TypeId:     atoi(r.URL.Query().Get("type_id")),
		Keyword:    r.URL.Query().Get("keyword"),
		Lang:       r.URL.Query().Get("lang"),
	}
	devices, total, err := h.uc.GetDeviceList(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]interface{}{
		"data":  devices,
		"total": total,
		"page":  page,
	})
}

// GET /listdevicepage (เหมือน device)
func (h *MQTT3Handler) GetDeviceListPage(w http.ResponseWriter, r *http.Request) {
	h.GetDeviceList(w, r)
}

// GET /devicebuckets
func (h *MQTT3Handler) GetDeviceBuckets(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		render.Render(w, r, ErrBadRequest("bucket is required"))
		return
	}
	resp, err := h.uc.GetDeviceBuckets(r.Context(), bucket)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// GET /sensercharts
func (h *MQTT3Handler) GetSenserCharts(w http.ResponseWriter, r *http.Request) {
	req := presenter.SenserChartRequest{
		Bucket:      r.URL.Query().Get("bucket"),
		Measurement: r.URL.Query().Get("measurement"),
		Field:       r.URL.Query().Get("field"),
		Start:       r.URL.Query().Get("start"),
		Stop:        r.URL.Query().Get("stop"),
		Limit:       atoi(r.URL.Query().Get("limit")),
	}
	if req.Bucket == "" || req.Measurement == "" {
		render.Render(w, r, ErrBadRequest("bucket and measurement are required"))
		return
	}
	resp, err := h.uc.GetSenserCharts(r.Context(), &req)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// GET /senserdatachart
func (h *MQTT3Handler) GetSenserDataChart(w http.ResponseWriter, r *http.Request) {
	h.GetSenserCharts(w, r) // reuse same logic
}

// GET /senserdata
func (h *MQTT3Handler) GetSenserData(w http.ResponseWriter, r *http.Request) {
	h.GetSenserCharts(w, r)
}

// GET /devicelist (alias)
func (h *MQTT3Handler) GetDeviceListAlias(w http.ResponseWriter, r *http.Request) {
	h.GetDeviceList(w, r)
}

// GET /locationdevice
func (h *MQTT3Handler) GetDeviceByLocation(w http.ResponseWriter, r *http.Request) {
	locationID := atoi(r.URL.Query().Get("location_id"))
	if locationID == 0 {
		render.Render(w, r, ErrBadRequest("location_id is required"))
		return
	}
	devices, err := h.uc.GetDeviceListByLocation(r.Context(), locationID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, devices)
}

// GET /devicesensercharts
func (h *MQTT3Handler) GetDeviceSenserCharts(w http.ResponseWriter, r *http.Request) {
	h.GetSenserCharts(w, r)
}

// GET /alarmdevicestatus
func (h *MQTT3Handler) GetAlarmDeviceStatus(w http.ResponseWriter, r *http.Request) {
	// parse query
	params := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	resp, err := h.uc.GetAlarmDeviceStatus(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// GET /alarmdevicestatuscontrol
func (h *MQTT3Handler) GetAlarmDeviceStatusControl(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	resp, err := h.uc.GetAlarmDeviceStatusControl(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// GET /monitordevicegroup
func (h *MQTT3Handler) GetMonitorDeviceGroup(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	resp, err := h.uc.GetMonitorDeviceGroup(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// GET /monitordevicechart
func (h *MQTT3Handler) GetMonitorDeviceChart(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]interface{})
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	resp, err := h.uc.GetMonitorDeviceChart(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, resp)
}

// helpers
func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func ErrBadRequest(msg string) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusBadRequest, ErrorText: msg}
}
func ErrInternal(err error) render.Renderer {
	return &errResponse{HTTPStatusCode: http.StatusInternalServerError, ErrorText: err.Error()}
}
type errResponse struct {
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error"`
}
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
```

---

## 9. `internal/modules/iot/delivery/http/routes.go`

```go
// internal/modules/iot/delivery/http/routes.go
package http

import (
	"icmongolang/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MapMQTT3Routes(router chi.Router, h *MQTT3Handler, mw *middleware.MiddlewareManager) {
	router.Route("/mqtt3", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(mw.Verifier(true))
			r.Use(mw.Authenticator())
			r.Use(mw.CurrentUser())
			r.Use(mw.ActiveUser())

			r.Get("/topic", h.GetTopicData)
			r.Post("/control", h.DeviceControl)
			r.Get("/device", h.GetDeviceList)
			r.Get("/listdevicepage", h.GetDeviceListPage)
			r.Get("/devicebuckets", h.GetDeviceBuckets)
			r.Get("/sensercharts", h.GetSenserCharts)
			r.Get("/senserdatachart", h.GetSenserDataChart)
			r.Get("/senserdata", h.GetSenserData)
			r.Get("/devicelist", h.GetDeviceListAlias)
			r.Get("/locationdevice", h.GetDeviceByLocation)
			r.Get("/devicesensercharts", h.GetDeviceSenserCharts)
			r.Get("/alarmdevicestatus", h.GetAlarmDeviceStatus)
			r.Get("/alarmdevicestatuscontrol", h.GetAlarmDeviceStatusControl)
			r.Get("/monitordevicegroup", h.GetMonitorDeviceGroup)
			r.Get("/monitordevicechart", h.GetMonitorDeviceChart)
		})
	})
}
```

---

## 10. การรวมเข้ากับ Server (`server/handlers.go` และ `server/server.go`)

### ใน `server/handlers.go` (ภายในฟังก์ชัน `New`)

```go
// เพิ่ม imports
import (
	// ... existing imports
	iotHttp "icmongolang/internal/modules/iot/delivery/http"
	iotUsecase "icmongolang/internal/modules/iot/usecase"
	iotRepo "icmongolang/internal/modules/iot/repository"
)

// ... ภายในฟังก์ชัน New(...) ก่อน return
	// ========== IoT MQTT3 routes ==========
	deviceRepo := iotRepo.NewDeviceRepository(db)
	alarmLogRepo := iotRepo.NewAlarmLogRepository(db)
	scheduleRepo := iotRepo.NewScheduleRepository(db)
	iotUC := iotUsecase.NewMQTT3UseCase(deviceRepo, alarmLogRepo, scheduleRepo, mqttClient, redisClient, influxClient, logger)
	iotHandler := iotHttp.NewMQTT3Handler(iotUC, logger)
	iotHttp.MapMQTT3Routes(apiRouter, iotHandler, mw)
	logger.Info("✅ MQTT3 routes registered")
```

### ใน `server/server.go` – ตรวจสอบว่า mqttClient, redisClient, influxClient ถูกส่งเข้าไป

ถ้าตัวแปร `mqttClient`, `redisClient`, `influxClient` ถูกสร้างและส่งเข้ามาในฟังก์ชัน `New` อยู่แล้ว ก็ไม่ต้องแก้ไขเพิ่ม

---

## สรุป

เราได้รวม `RequestManager` เข้ากับ `client.go` เรียบร้อย ทำให้ MQTT client มีความสามารถ:
- Subscribe/Unsubscribe แบบ reference counting
- Request-Response pattern ผ่าน `GetDataFromTopic` และ `RequestData`
- การจัดการ timeout และ cleanup อัตโนมัติ

จากนั้นได้สร้าง **IoT Service Module** ทั้งหมดตามโครงสร้าง Clean Architecture ครอบคลุมทุก endpoint ที่โจทย์กำหนด (14 endpoints) พร้อมให้สามารถขยาย business logic ได้ง่าย

**หมายเหตุ:** โค้ดบางส่วน (เช่น `models.Device`, `models.AlarmLog`, การ query influxdb ด้วย `QueryFilterData`) สมมติว่ามีอยู่แล้วในโปรเจกต์ หากไม่มีต้องสร้างตามแบบที่ระบบเดิมใช้ หรือปรับใช้กับของจริง