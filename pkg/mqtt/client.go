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
	getTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
	GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error)
}

type pendingRequest struct {
	resolve chan<- interface{}
	reject  chan<- error
	timeout *time.Timer
	topic   string
	raw     bool
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
	rm.logger.Info("requestManager initialized with global message handler")
	return rm
}

func (rm *requestManager) handleIncomingMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	rm.mu.Lock()
	defer rm.mu.Unlock()

	pendings, ok := rm.pending[topic]
	if !ok || len(pendings) == 0 {
		rm.logger.Debugf("handleIncomingMessage: no pending request for topic=%s, ignoring", topic)
		return
	}
	req := pendings[0]
	rm.pending[topic] = pendings[1:]
	if len(rm.pending[topic]) == 0 {
		delete(rm.pending, topic)
		rm.decrementSubscription(topic)
	}
	req.timeout.Stop()

	if req.raw {
		rm.logger.Infof("handleIncomingMessage: received message on topic=%s, resolving pending raw request", topic)
		req.resolve <- msg.Payload()
		return
	}

	var result interface{}
	payload := msg.Payload()
	if err := json.Unmarshal(payload, &result); err != nil {
		result = string(payload)
	}
	rm.logger.Infof("handleIncomingMessage: received message on topic=%s, resolving pending request", topic)
	req.resolve <- result
}

func (rm *requestManager) getTopic(ctx context.Context, topic string, timeoutMs int) (interface{}, error) {
	return rm.getTopicInternal(ctx, topic, timeoutMs, false)
}

func (rm *requestManager) getTopicRaw(ctx context.Context, topic string, timeoutMs int) ([]byte, error) {
	result, err := rm.getTopicInternal(ctx, topic, timeoutMs, true)
	if err != nil {
		return nil, err
	}
	switch v := result.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		marshaled, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}
		return marshaled, nil
	}
}

func (rm *requestManager) getTopicInternal(ctx context.Context, topic string, timeoutMs int, raw bool) (interface{}, error) {
	rm.logger.Debugf("getTopicInternal: start for topic=%s timeout=%dms raw=%v", topic, timeoutMs, raw)
	rm.mu.Lock()
	rm.logger.Debugf("getTopicInternal: lock acquired for topic=%s", topic)

	if _, ok := rm.pending[topic]; !ok {
		rm.logger.Debugf("getTopicInternal: no pending request for topic=%s, calling incrementSubscription", topic)
		if err := rm.incrementSubscription(topic); err != nil {
			rm.mu.Unlock()
			return nil, err
		}
	} else {
		rm.logger.Debugf("getTopicInternal: existing pending queue for topic=%s, current length=%d", topic, len(rm.pending[topic]))
	}

	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	timer := time.NewTimer(time.Duration(timeoutMs) * time.Millisecond)

	req := &pendingRequest{
		resolve: resultCh,
		reject:  errCh,
		timeout: timer,
		topic:   topic,
		raw:     raw,
	}
	rm.pending[topic] = append(rm.pending[topic], req)
	rm.logger.Debugf("getTopicInternal: pending request added for topic=%s, queue length=%d", topic, len(rm.pending[topic]))
	rm.mu.Unlock()

	go func() {
		var err error
		select {
		case <-timer.C:
			rm.logger.Warnf("getTopicInternal: timeout occurred for topic=%s after %d ms", topic, timeoutMs)
			err = fmt.Errorf("timeout waiting for message on topic %s after %d ms", topic, timeoutMs)
		case <-ctx.Done():
			rm.logger.Warnf("getTopicInternal: context cancelled while waiting for topic=%s: %v", topic, ctx.Err())
			err = ctx.Err()
		}

		rm.mu.Lock()
		removed := false
		list := rm.pending[topic]
		for i, r := range list {
			if r == req {
				rm.pending[topic] = append(list[:i], list[i+1:]...)
				removed = true
				rm.logger.Debugf("getTopicInternal: removed request from pending for topic=%s, new length=%d", topic, len(rm.pending[topic]))
				break
			}
		}
		if removed && len(rm.pending[topic]) == 0 {
			delete(rm.pending, topic)
			rm.logger.Debugf("getTopicInternal: no more pending requests for topic=%s, deleting from map", topic)
			rm.decrementSubscription(topic)
		}
		rm.mu.Unlock()

		select {
		case errCh <- err:
		default:
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			rm.logger.Errorf("getTopicInternal: returning error for topic=%s: %v", topic, err)
			return nil, err
		}
		res := <-resultCh
		rm.logger.Infof("getTopicInternal: success (via errCh) for topic=%s, result=%v", topic, res)
		return res, nil
	case res := <-resultCh:
		timer.Stop()
		rm.logger.Infof("getTopicInternal: success (via resultCh) for topic=%s, result=%v", topic, res)
		return res, nil
	}
}

func (rm *requestManager) incrementSubscription(topic string) error {
	rm.logger.Debugf("incrementSubscription: start for topic=%s", topic)
	count := rm.subscriptionCount[topic]
	if count == 0 {
		rm.logger.Debugf("incrementSubscription: first subscriber for topic=%s, calling client.Subscribe", topic)
		token := rm.client.Subscribe(topic, 0, rm.handleIncomingMessage)
		if !token.WaitTimeout(5 * time.Second) {
			err := fmt.Errorf("subscribe timeout for topic %s", topic)
			rm.logger.Errorf("incrementSubscription: %v", err)
			rm.rejectPending(topic, err)
			return err
		}
		if token.Error() != nil {
			rm.logger.Errorf("incrementSubscription: subscribe failed for topic=%s: %v", topic, token.Error())
			rm.rejectPending(topic, token.Error())
			return token.Error()
		}
		rm.subscriptionCount[topic] = 1
		rm.logger.Infof("incrementSubscription: subscribed to topic=%s, count=1", topic)
	} else {
		rm.subscriptionCount[topic] = count + 1
		rm.logger.Debugf("incrementSubscription: topic=%s already subscribed, count increased to %d", topic, count+1)
	}
	return nil
}

func (rm *requestManager) rejectPending(topic string, err error) {
	if pendings, ok := rm.pending[topic]; ok {
		for _, req := range pendings {
			req.reject <- err
			req.timeout.Stop()
		}
		delete(rm.pending, topic)
	}
}

func (rm *requestManager) decrementSubscription(topic string) {
	rm.logger.Debugf("decrementSubscription: start for topic=%s", topic)
	count := rm.subscriptionCount[topic]
	if count <= 1 {
		rm.logger.Debugf("decrementSubscription: last subscriber for topic=%s, calling client.Unsubscribe", topic)
		rm.client.Unsubscribe(topic)
		delete(rm.subscriptionCount, topic)
		rm.logger.Infof("decrementSubscription: unsubscribed from topic=%s", topic)
	} else {
		rm.subscriptionCount[topic] = count - 1
		rm.logger.Debugf("decrementSubscription: topic=%s count decreased to %d", topic, count-1)
	}
}

// mqttClient implements Client interface
type mqttClient struct {
	client         mqtt.Client
	cfg            *config.MQTTConfig
	logger         logger.Logger
	opts           *mqtt.ClientOptions
	requestManager *requestManager
}

// New creates a new MQTT client with corrected configuration
func New(cfg *config.MQTTConfig, log logger.Logger) Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)

	// Client ID – generate unique if not provided
	clientID := cfg.ClientID
	if clientID == "" {
		clientID = fmt.Sprintf("go-client-%d", time.Now().UnixNano())
	}

	log.Infof("clientID =%s,", clientID)
	log.Infof("Username =%s,", cfg.Username)
	log.Infof("Password =%s,", cfg.Password)
	log.Infof("Broker =%s,", cfg.Broker)
	log.Infof("Debug:cfg =%s,", cfg)
	log.Infof("Debug: MQTT Config: fg.Broker=%s, cfg.ClientID=%s", cfg)

	opts.SetClientID(clientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)

	// ✅ Critical settings for reliable subscriptions
	opts.SetCleanSession(true)  // Start fresh every connection
	opts.SetAutoReconnect(true) // Auto-reconnect on disconnect
	opts.SetConnectRetry(true)  // Retry initial connection
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetResumeSubs(false) // true,flase Resume subscriptions after reconnect  (FIXED)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(15 * time.Second)
	opts.SetWriteTimeout(15 * time.Second)
	opts.SetMaxReconnectInterval(10 * time.Second)

	// For self-signed certificates (internal testing)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	log.Infof("MQTT Config: broker=%s, clientID=%s, cleanSession=true, autoReconnect=true, resumeSubs=true", cfg.Broker, clientID)

	// Connection event callbacks
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		log.Errorf("MQTT connection lost: %v", err)
	}
	opts.OnReconnecting = func(cl mqtt.Client, opts *mqtt.ClientOptions) {
		log.Info("MQTT reconnecting...")
	}
	opts.OnConnect = func(cl mqtt.Client) {
		log.Info("MQTT connected")
	}
	log.Infof("Debug:opts =%s,", opts)
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
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
	}
	token := c.client.Publish(topic, qos, retained, msgPayload)
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

// RequestData implements request-response using the internal requestManager.
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
	token := c.client.Publish(requestTopic, 1, false, msgPayload)
	token.Wait()
	if token.Error() != nil {
		return nil, fmt.Errorf("publish failed: %w", token.Error())
	}
	c.logger.Debugf("RequestData: published to %s, waiting for response on %s", requestTopic, waitTopic)
	return c.requestManager.getTopic(context.Background(), waitTopic, int(timeout.Milliseconds()))
}

func (c *mqttClient) getTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
	if c.requestManager == nil {
		return nil, fmt.Errorf("request manager not initialized (client not connected)")
	}
	return c.requestManager.getTopicRaw(ctx, topic, int(timeout.Milliseconds()))
}

func (c *mqttClient) GetDataFromTopic(ctx context.Context, topic string, timeout time.Duration) ([]byte, error) {
	if c.client == nil || !c.client.IsConnected() {
		c.logger.Errorf("GetDataFromTopic: MQTT client not connected")
		return nil, fmt.Errorf("MQTT client not connected")
	}
	if c.requestManager == nil {
		return nil, fmt.Errorf("request manager not initialized (client not connected)")
	}
	c.logger.Infof("GetDataFromTopic: waiting for message on topic=%s, timeout=%v", topic, timeout)

	data, err := c.requestManager.getTopicRaw(ctx, topic, int(timeout.Milliseconds()))
	if err != nil {
		c.logger.Warnf("GetDataFromTopic: failed for topic=%s: %v", topic, err)
		return nil, err
	}
	c.logger.Infof("GetDataFromTopic: got %d bytes from %s", len(data), topic)
	return data, nil
}
