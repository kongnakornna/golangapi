package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"icmongolang/internal/modules/queue"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// HubCallbacks defines lifecycle hooks for client connections.
type HubCallbacks struct {
	// OnConnect is called when a client successfully connects and registers.
	// Return an error to reject the connection (client will be disconnected).
	OnConnect func(client *Client) error

	// OnDisconnect is called when a client disconnects (cleanly or due to error).
	// The client has already been unregistered from topics/rooms at this point.
	OnDisconnect func(client *Client)

	// OnSubscribe is called when a client subscribes to a topic.
	OnSubscribe func(client *Client, topic string)

	// OnUnsubscribe is called when a client unsubscribes from a topic.
	OnUnsubscribe func(client *Client, topic string)

	// OnJoinRoom is called when a client joins a room.
	OnJoinRoom func(client *Client, room string)

	// OnLeaveRoom is called when a client leaves a room.
	OnLeaveRoom func(client *Client, room string)

	// OnMessage is called when a client sends a message (before broadcasting).
	// Return false to prevent the message from being broadcast.
	OnMessage func(client *Client, topic, room string, payload []byte) bool
}

// HubConfig holds configuration for the Hub.
type HubConfig struct {
	Callbacks      HubCallbacks
	Queue          queue.Queue
	Log            logger.Logger
	RedisClient    *redis.Client // Redis client for horizontal scaling
	InstanceID     string        // Unique identifier for this hub instance
	EnablePresence bool          // Enable presence tracking and broadcasting
}

// Hub ดูแล client และ topic subscriptions (ใช้ Redis queue)
type Hub struct {
	clients       map[string]*Client
	topics        map[string]map[string]bool // topic -> clientID -> bool
	rooms         map[string]map[string]bool // room -> clientID -> bool
	topicCancel   map[string]context.CancelFunc
	mu            sync.RWMutex
	queue         queue.Queue
	log           logger.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	Callbacks     HubCallbacks
	redisClient   *redis.Client
	instanceID    string
	pubsub        *redis.PubSub
	redisChannels map[string]*redis.PubSub // topic/room -> pubsub

	// Presence tracking
	userPresence   map[string]*PresenceInfo // userID -> PresenceInfo
	presenceMu     sync.RWMutex
	enablePresence bool
}

func NewHub(cfg HubConfig) *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	instanceID := uuid.New().String()[:8]
	h := &Hub{
		clients:        make(map[string]*Client),
		topics:         make(map[string]map[string]bool),
		rooms:          make(map[string]map[string]bool),
		topicCancel:    make(map[string]context.CancelFunc),
		queue:          cfg.Queue,
		log:            cfg.Log,
		ctx:            ctx,
		cancel:         cancel,
		Callbacks:      cfg.Callbacks,
		redisClient:    cfg.RedisClient,
		instanceID:     instanceID,
		redisChannels:  make(map[string]*redis.PubSub),
		userPresence:   make(map[string]*PresenceInfo),
		enablePresence: cfg.EnablePresence,
	}

	// Start Redis pub/sub listener if Redis client is provided
	if cfg.RedisClient != nil {
		go h.listenRedisPubSub()
	}

	return h
}

func (h *Hub) Run() {
	// ไม่มี loop หลัก, การ subscribe จะเกิดขึ้นเมื่อมี client subscribe topic แรก
}

// AddClient creates a new client, registers it, and optionally joins a room.
// This is the preferred way to add a client from the WebSocket handler.
func (h *Hub) AddClient(conn *websocket.Conn, userID, room string) *Client {
	client := NewClient(generateClientID(), userID, conn, h)
	// Initialize AckManager for this client
	client.AckManager = NewAckManager(DefaultAckConfig(), h.log)
	h.Register(client)
	if room != "" {
		h.JoinRoom(client, room)
	}
	return client
}

// generateClientID creates a unique client identifier.
func generateClientID() string {
	return uuid.New().String()
}

// Register adds a client and calls OnConnect callback.
func (h *Hub) Register(client *Client) error {
	h.mu.Lock()
	h.clients[client.ID] = client
	wsConnectedClients.Inc()
	wsConnectionsTotal.Inc()
	h.mu.Unlock()

	// Update presence - set user as online
	h.updateUserPresence(client.UserID, PresenceOnline, map[string]any{
		"connected_at": time.Now().Format(time.RFC3339),
		"instance_id":  h.instanceID,
	})

	// Call OnConnect callback (outside lock to avoid deadlock)
	if h.Callbacks.OnConnect != nil {
		if err := h.Callbacks.OnConnect(client); err != nil {
			// Reject connection: remove client and return error
			h.mu.Lock()
			delete(h.clients, client.ID)
			wsConnectedClients.Dec()
			h.mu.Unlock()
			return err
		}
	}
	return nil
}

// Unregister removes client and cleans subscriptions
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	// remove from topics
	for topic := range client.Topics {
		if clients, ok := h.topics[topic]; ok {
			delete(clients, client.ID)
			if len(clients) == 0 {
				delete(h.topics, topic)
				wsActiveTopics.Dec()
			}
		}
	}
	// remove from rooms
	for room := range client.Rooms {
		if roomClients, ok := h.rooms[room]; ok {
			delete(roomClients, client.ID)
			if len(roomClients) == 0 {
				delete(h.rooms, room)
				wsActiveRooms.Dec()
			}
		}
	}
	delete(h.clients, client.ID)
	wsConnectedClients.Dec()

	// Check if user has any other connections
	hasOtherConnections := false
	for _, c := range h.clients {
		if c.UserID == client.UserID {
			hasOtherConnections = true
			break
		}
	}
	h.mu.Unlock()

	client.close()

	// Update presence - set user as offline if no other connections
	if !hasOtherConnections {
		h.updateUserPresence(client.UserID, PresenceOffline, map[string]any{
			"disconnected_at": time.Now().Format(time.RFC3339),
		})
	}

	// Call OnDisconnect callback (outside lock)
	if h.Callbacks.OnDisconnect != nil {
		h.Callbacks.OnDisconnect(client)
	}
}

// Subscribe client to a topic (และเริ่ม queue subscriber ถ้ายังไม่มี)
func (h *Hub) Subscribe(client *Client, topic string) error {
	h.mu.Lock()
	if h.topics[topic] == nil {
		h.topics[topic] = make(map[string]bool)
		wsActiveTopics.Inc()
		topicCtx, topicCancel := context.WithCancel(h.ctx)
		h.topicCancel[topic] = topicCancel
		go h.subscribeQueueTopic(topicCtx, topic)
	}
	h.topics[topic][client.ID] = true
	client.Topics[topic] = true
	h.mu.Unlock()

	// Call OnSubscribe callback
	if h.Callbacks.OnSubscribe != nil {
		h.Callbacks.OnSubscribe(client, topic)
	}
	return nil
}

// Unsubscribe client from topic
func (h *Hub) Unsubscribe(client *Client, topic string) error {
	h.mu.Lock()
	if clients, ok := h.topics[topic]; ok {
		delete(clients, client.ID)
		if len(clients) == 0 {
			delete(h.topics, topic)
			wsActiveTopics.Dec()
			if cancel, ok := h.topicCancel[topic]; ok {
				cancel()
				delete(h.topicCancel, topic)
			}
		}
	}
	delete(client.Topics, topic)
	h.mu.Unlock()

	// Call OnUnsubscribe callback
	if h.Callbacks.OnUnsubscribe != nil {
		h.Callbacks.OnUnsubscribe(client, topic)
	}
	return nil
}

// ============================================================================
// Presence Tracking
// ============================================================================

// updateUserPresence updates a user's presence status and broadcasts the change.
func (h *Hub) updateUserPresence(userID string, status PresenceStatus, meta map[string]any) {
	if !h.enablePresence {
		return
	}
	h.presenceMu.Lock()
	defer h.presenceMu.Unlock()

	_ = time.Now()
	existing, exists := h.userPresence[userID]

	var instances []string
	if exists {
		instances = existing.Instances
		// Update instance list
		found := false
		for _, inst := range instances {
			if inst == h.instanceID {
				found = true
				break
			}
		}
		if !found {
			instances = append(instances, h.instanceID)
		}
	} else {
		instances = []string{h.instanceID}
	}

	// Merge meta
	mergedMeta := make(map[string]any)
	if exists {
		for k, v := range existing.Meta {
			mergedMeta[k] = v
		}
	}
	for k, v := range meta {
		mergedMeta[k] = v
	}

	h.userPresence[userID] = &PresenceInfo{
		UserID:    userID,
		Status:    status,
		Meta:      mergedMeta,
		LastSeen:  time.Now(),
		Instances: instances,
	}

	// Broadcast presence change
	h.broadcastPresenceChange(userID, status, mergedMeta)

	// Publish to Redis for other instances
	if h.redisClient != nil {
		ctx := context.Background()
		h.publishPresenceChange(ctx, userID, status, mergedMeta)
	}
}

// GetUserPresence returns the presence info for a user.
func (h *Hub) GetUserPresence(userID string) (*PresenceInfo, bool) {
	h.presenceMu.RLock()
	defer h.presenceMu.RUnlock()
	info, exists := h.userPresence[userID]
	return info, exists
}

// GetAllPresence returns all user presence info.
func (h *Hub) GetAllPresence() map[string]*PresenceInfo {
	h.presenceMu.RLock()
	defer h.presenceMu.RUnlock()
	result := make(map[string]*PresenceInfo, len(h.userPresence))
	for k, v := range h.userPresence {
		result[k] = v
	}
	return result
}

// GetPresenceByStatus returns all users with a specific status.
func (h *Hub) GetPresenceByStatus(status PresenceStatus) []*PresenceInfo {
	h.presenceMu.RLock()
	defer h.presenceMu.RUnlock()
	var result []*PresenceInfo
	for _, info := range h.userPresence {
		if info.Status == status {
			result = append(result, info)
		}
	}
	return result
}

// broadcastPresenceChange broadcasts a presence change to all connected clients.
func (h *Hub) broadcastPresenceChange(userID string, status PresenceStatus, meta map[string]any) {
	if !h.enablePresence {
		return
	}
	env := BuildPresenceEvent(userID, string(status), meta)
	data, _ := json.Marshal(env)
	h.mu.RLock()
	for _, client := range h.clients {
		client.safeSend(data)
	}
	h.mu.RUnlock()
}

// publishPresenceChange publishes a presence change to Redis for other instances.
func (h *Hub) publishPresenceChange(ctx context.Context, userID string, status PresenceStatus, meta map[string]any) error {
	if !h.enablePresence || h.redisClient == nil {
		return nil
	}

	data, _ := json.Marshal(map[string]any{
		"user_id": userID,
		"status":  status,
		"meta":    meta,
	})

	msg := RedisMessage{
		InstanceID: h.instanceID,
		Type:       "presence_change",
		Event:      "presence",
		Data:       data,
	}
	msgBytes, _ := json.Marshal(msg)
	return h.redisClient.Publish(ctx, "websocket:hub:global", msgBytes).Err()
}

// handlePresenceChange handles a presence change message from another instance.
func (h *Hub) handlePresenceChange(payload string) {
	if !h.enablePresence {
		return
	}
	var msg map[string]any
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		h.log.Errorf("Failed to unmarshal presence change: %v", err)
		return
	}

	userID, _ := msg["user_id"].(string)
	statusStr, _ := msg["status"].(string)
	meta, _ := msg["meta"].(map[string]any)

	status := PresenceStatus(statusStr)

	h.presenceMu.Lock()
	defer h.presenceMu.Unlock()

	existing, exists := h.userPresence[userID]
	instances := []string{}
	if exists {
		instances = existing.Instances
	}
	// Add instance if not present
	found := false
	for _, inst := range instances {
		if inst == msg["instance_id"] {
			found = true
			break
		}
	}
	if !found {
		instances = append(instances, msg["instance_id"].(string))
	}

	mergedMeta := make(map[string]any)
	if exists {
		for k, v := range existing.Meta {
			mergedMeta[k] = v
		}
	}
	if meta != nil {
		for k, v := range meta {
			mergedMeta[k] = v
		}
	}

	h.userPresence[userID] = &PresenceInfo{
		UserID:    userID,
		Status:    status,
		Meta:      mergedMeta,
		LastSeen:  time.Now(),
		Instances: instances,
	}

	// Broadcast locally
	h.broadcastPresenceChange(userID, status, mergedMeta)
}

// JoinRoom adds client to a room
func (h *Hub) JoinRoom(client *Client, room string) {
	h.mu.Lock()
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[string]bool)
		wsActiveRooms.Inc()
	}
	h.rooms[room][client.ID] = true
	client.Rooms[room] = true
	h.mu.Unlock()

	// Call OnJoinRoom callback
	if h.Callbacks.OnJoinRoom != nil {
		h.Callbacks.OnJoinRoom(client, room)
	}
}

// LeaveRoom removes client from room
func (h *Hub) LeaveRoom(client *Client, room string) {
	h.mu.Lock()
	if clients, ok := h.rooms[room]; ok {
		delete(clients, client.ID)
		if len(clients) == 0 {
			delete(h.rooms, room)
			wsActiveRooms.Dec()
		}
	}
	delete(client.Rooms, room)
	h.mu.Unlock()

	// Call OnLeaveRoom callback
	if h.Callbacks.OnLeaveRoom != nil {
		h.Callbacks.OnLeaveRoom(client, room)
	}
}

// UpdateClientRoom updates a client's current room for presence tracking.
func (h *Hub) UpdateClientRoom(client *Client, oldRoom, newRoom string) {
	if oldRoom == newRoom {
		return
	}

	h.mu.Lock()
	// Remove from old room's presence tracking
	if oldRoom != "" {
		if clients, ok := h.rooms[oldRoom]; ok {
			delete(clients, client.ID)
			if len(clients) == 0 {
				delete(h.rooms, oldRoom)
				wsActiveRooms.Dec()
			}
		}
	}
	// Add to new room's presence tracking
	if newRoom != "" {
		if h.rooms[newRoom] == nil {
			h.rooms[newRoom] = make(map[string]bool)
			wsActiveRooms.Inc()
		}
		h.rooms[newRoom][client.ID] = true
		client.Rooms[newRoom] = true
	}
	h.mu.Unlock()

	// Update presence meta with current room
	h.updateUserPresence(client.UserID, PresenceOnline, map[string]any{
		"current_room": newRoom,
		"updated_at":   time.Now().Format(time.RFC3339),
	})

	// Call OnLeaveRoom for old room
	if oldRoom != "" && h.Callbacks.OnLeaveRoom != nil {
		h.Callbacks.OnLeaveRoom(client, oldRoom)
	}
	// Call OnJoinRoom for new room
	if newRoom != "" && h.Callbacks.OnJoinRoom != nil {
		h.Callbacks.OnJoinRoom(client, newRoom)
	}
}

// GetUserCurrentRoom returns the current room for a user (if they have an active connection).
func (h *Hub) GetUserCurrentRoom(userID string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients {
		if client.UserID == userID {
			return client.CurrentRoom
		}
	}
	return ""
}

// GetRoomPresence returns all users currently in a specific room.
func (h *Hub) GetRoomPresence(room string) []*PresenceInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []*PresenceInfo
	if clients, ok := h.rooms[room]; ok {
		for clientID := range clients {
			if client, ok := h.clients[clientID]; ok {
				h.presenceMu.RLock()
				if info, exists := h.userPresence[client.UserID]; exists {
					result = append(result, info)
				}
				h.presenceMu.RUnlock()
			}
		}
	}
	return result
}

// BroadcastToRoom sends event/data to all clients in room with sequence numbers
func (h *Hub) BroadcastToRoom(room, event string, data interface{}) {
	h.mu.RLock()
	clientsMap := h.rooms[room]
	var clientIDs []string
	for cid := range clientsMap {
		clientIDs = append(clientIDs, cid)
	}
	h.mu.RUnlock()

	wsMessagesBroadcastTotal.WithLabelValues("room").Add(float64(len(clientIDs)))

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			// Assign sequence number for this client
			seq := client.NextSeq()
			env := NewEnvelope("message").
				WithID(GenerateID()).
				WithSeq(seq).
				WithTopic("").
				WithRoom(room).
				WithPayload(map[string]any{
					"event": event,
					"data":  data,
				})
			// Track for acknowledgment
			client.TrackPending(env)
			envBytes, _ := json.Marshal(env)
			client.safeSend(envBytes)
		}
	}

	// Also publish to Redis for other instances
	if h.redisClient != nil {
		ctx := context.Background()
		h.PublishToRoom(ctx, room, event, data)
	}
}

// BroadcastToTopic sends event/data to all clients subscribed to a topic
func (h *Hub) BroadcastToTopic(topic, event string, data interface{}) {
	h.mu.RLock()
	clientIDs := make([]string, 0, len(h.topics[topic]))
	for cid := range h.topics[topic] {
		clientIDs = append(clientIDs, cid)
	}
	h.mu.RUnlock()

	wsMessagesBroadcastTotal.WithLabelValues("topic").Add(float64(len(clientIDs)))

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			// Assign sequence number for this client
			seq := client.NextSeq()
			env := NewEnvelope("message").
				WithID(GenerateID()).
				WithSeq(seq).
				WithTopic(topic).
				WithRoom("").
				WithPayload(map[string]any{
					"event": event,
					"data":  data,
				})
			// Track for acknowledgment
			client.TrackPending(env)
			envBytes, _ := json.Marshal(env)
			client.safeSend(envBytes)
		}
	}

	// Also publish to Redis for other instances
	if h.redisClient != nil {
		ctx := context.Background()
		h.PublishToTopic(ctx, topic, event, data)
	}
}

// GetRooms returns list of active room names
func (h *Hub) GetRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	rooms := make([]string, 0, len(h.rooms))
	for room := range h.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// GetClientsInRoom returns count
func (h *Hub) GetClientsInRoom(room string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[room])
}

// subscribeQueueTopic เริ่ม subscriber ของ queue
func (h *Hub) subscribeQueueTopic(ctx context.Context, topic string) {
	handler := func(ctx context.Context, msg *queue.Message) error {
		// ส่ง payload ไปยังทุก client ที่ subscribe topic นี้
		h.BroadcastToTopic(topic, "message", json.RawMessage(msg.Payload))
		return nil
	}
	if err := h.queue.Subscribe(ctx, topic, handler); err != nil {
		h.log.Errorf("Failed to subscribe to queue topic %s: %v", topic, err)
	}
}

// PublishMessage publishes a message to queue and optionally saves to DB
func (h *Hub) PublishMessage(ctx context.Context, topic string, payload []byte, senderID string) error {
	return h.queue.Publish(ctx, topic, json.RawMessage(payload))
}

// BroadcastToRole sends event/data to all clients with a specific role
func (h *Hub) BroadcastToRole(role, event string, data interface{}) {
	h.mu.RLock()
	var clientIDs []string
	for _, client := range h.clients {
		if client.Role == role {
			clientIDs = append(clientIDs, client.ID)
		}
	}
	h.mu.RUnlock()

	payload := map[string]interface{}{
		"event": event,
		"data":  data,
	}
	msgBytes, _ := json.Marshal(payload)

	wsMessagesBroadcastTotal.WithLabelValues("role").Add(float64(len(clientIDs)))

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			client.safeSend(msgBytes)
		}
	}
}

// BroadcastToRoles sends event/data to all clients with any of the specified roles
func (h *Hub) BroadcastToRoles(roles []string, event string, data interface{}) {
	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r] = true
	}

	h.mu.RLock()
	var clientIDs []string
	for _, client := range h.clients {
		if roleSet[client.Role] {
			clientIDs = append(clientIDs, client.ID)
		}
	}
	h.mu.RUnlock()

	payload := map[string]interface{}{
		"event": event,
		"data":  data,
	}
	msgBytes, _ := json.Marshal(payload)

	wsMessagesBroadcastTotal.WithLabelValues("roles").Add(float64(len(clientIDs)))

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			client.safeSend(msgBytes)
		}
	}

	// Also publish to Redis for other instances
	if h.redisClient != nil {
		ctx := context.Background()
		h.PublishToRoles(ctx, roles, event, data)
	}
}

// BroadcastMessage ส่ง event/data ไปยังทุก client (ทั่วโลก)
func (h *Hub) BroadcastMessage(event string, data interface{}) {
	dataBytes, _ := json.Marshal(data)
	payload := map[string]interface{}{
		"event": event,
		"data":  dataBytes,
	}
	msgBytes, _ := json.Marshal(payload)

	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	wsMessagesBroadcastTotal.WithLabelValues("global").Add(float64(len(clients)))

	for _, client := range clients {
		client.safeSend(msgBytes)
	}

	// Also publish to Redis for other instances
	if h.redisClient != nil {
		ctx := context.Background()
		h.PublishGlobal(ctx, event, data)
	}
}

// ============================================================================
// Redis Pub/Sub for Horizontal Scaling
// ============================================================================

// RedisMessage represents a message sent via Redis pub/sub for cross-instance communication.
type RedisMessage struct {
	InstanceID string          `json:"instance_id"` // Originating instance ID
	Type       string          `json:"type"`        // "broadcast_room", "broadcast_topic", "broadcast_role", "broadcast_global", "presence_change"
	Event      string          `json:"event"`
	Data       json.RawMessage `json:"data"`
	Room       string          `json:"room,omitempty"`
	Topic      string          `json:"topic,omitempty"`
	Roles      []string        `json:"roles,omitempty"`
}

// listenRedisPubSub listens for messages from other hub instances via Redis pub/sub.
func (h *Hub) listenRedisPubSub() {
	if h.redisClient == nil {
		return
	}

	// Subscribe to the global hub channel
	pubsub := h.redisClient.Subscribe(h.ctx, "websocket:hub:global")
	h.pubsub = pubsub

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			h.handleRedisMessage(msg.Payload)
		case <-h.ctx.Done():
			pubsub.Close()
			return
		}
	}
}

// handleRedisMessage processes a message received from another hub instance.
func (h *Hub) handleRedisMessage(payload string) {
	var msg RedisMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		h.log.Errorf("Failed to unmarshal Redis message: %v", err)
		return
	}

	// Ignore messages from ourselves
	if msg.InstanceID == h.instanceID {
		return
	}

	h.log.Debugf("Received Redis message from instance %s: type=%s event=%s", msg.InstanceID, msg.Type, msg.Event)

	switch msg.Type {
	case "broadcast_room":
		if msg.Room != "" {
			h.broadcastToRoomLocal(msg.Room, msg.Event, msg.Data)
		}
	case "broadcast_topic":
		if msg.Topic != "" {
			h.broadcastToTopicLocal(msg.Topic, msg.Event, msg.Data)
		}
	case "broadcast_role":
		if len(msg.Roles) > 0 {
			h.broadcastToRolesLocal(msg.Roles, msg.Event, msg.Data)
		}
	case "broadcast_global":
		h.broadcastGlobalLocal(msg.Event, msg.Data)
	case "presence_change":
		h.handlePresenceChange(string(msg.Data))
	}
}

// broadcastToRoomLocal broadcasts to a room on this instance only (called from Redis message handler).
func (h *Hub) broadcastToRoomLocal(room, event string, data json.RawMessage) {
	h.mu.RLock()
	clientIDs := make([]string, 0, len(h.rooms[room]))
	for cid := range h.rooms[room] {
		clientIDs = append(clientIDs, cid)
	}
	h.mu.RUnlock()

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			client.safeSend(data)
		}
	}
}

// broadcastToTopicLocal broadcasts to a topic on this instance only.
func (h *Hub) broadcastToTopicLocal(topic, event string, data json.RawMessage) {
	h.mu.RLock()
	clientIDs := make([]string, 0, len(h.topics[topic]))
	for cid := range h.topics[topic] {
		clientIDs = append(clientIDs, cid)
	}
	h.mu.RUnlock()

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			client.safeSend(data)
		}
	}
}

// broadcastToRolesLocal broadcasts to roles on this instance only.
func (h *Hub) broadcastToRolesLocal(roles []string, event string, data json.RawMessage) {
	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r] = true
	}

	h.mu.RLock()
	var clientIDs []string
	for _, client := range h.clients {
		if roleSet[client.Role] {
			clientIDs = append(clientIDs, client.ID)
		}
	}
	h.mu.RUnlock()

	for _, cid := range clientIDs {
		h.mu.RLock()
		client := h.clients[cid]
		h.mu.RUnlock()
		if client != nil {
			client.safeSend(data)
		}
	}
}

// broadcastGlobalLocal broadcasts to all clients on this instance only.
func (h *Hub) broadcastGlobalLocal(event string, data json.RawMessage) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		client.safeSend(data)
	}
}

// PublishToRoom publishes a message to a room across all hub instances via Redis.
func (h *Hub) PublishToRoom(ctx context.Context, room, event string, data interface{}) error {
	if h.redisClient == nil {
		// Fallback to local broadcast only
		h.BroadcastToRoom(room, event, data)
		return nil
	}

	dataBytes, _ := json.Marshal(data)
	msg := RedisMessage{
		InstanceID: h.instanceID,
		Type:       "broadcast_room",
		Event:      event,
		Data:       dataBytes,
		Room:       room,
	}
	msgBytes, _ := json.Marshal(msg)
	return h.redisClient.Publish(ctx, "websocket:hub:global", msgBytes).Err()
}

// PublishToTopic publishes a message to a topic across all hub instances via Redis.
func (h *Hub) PublishToTopic(ctx context.Context, topic, event string, data interface{}) error {
	if h.redisClient == nil {
		h.BroadcastToTopic(topic, event, data)
		return nil
	}

	dataBytes, _ := json.Marshal(data)
	msg := RedisMessage{
		InstanceID: h.instanceID,
		Type:       "broadcast_topic",
		Event:      event,
		Data:       dataBytes,
		Topic:      topic,
	}
	msgBytes, _ := json.Marshal(msg)
	return h.redisClient.Publish(ctx, "websocket:hub:global", msgBytes).Err()
}

// PublishToRoles publishes a message to roles across all hub instances via Redis.
func (h *Hub) PublishToRoles(ctx context.Context, roles []string, event string, data interface{}) error {
	if h.redisClient == nil {
		h.BroadcastToRoles(roles, event, data)
		return nil
	}

	dataBytes, _ := json.Marshal(data)
	msg := RedisMessage{
		InstanceID: h.instanceID,
		Type:       "broadcast_role",
		Event:      event,
		Data:       dataBytes,
		Roles:      roles,
	}
	msgBytes, _ := json.Marshal(msg)
	return h.redisClient.Publish(ctx, "websocket:hub:global", msgBytes).Err()
}

// PublishGlobal publishes a message to all clients across all hub instances via Redis.
func (h *Hub) PublishGlobal(ctx context.Context, event string, data interface{}) error {
	if h.redisClient == nil {
		h.BroadcastMessage(event, data)
		return nil
	}

	dataBytes, _ := json.Marshal(data)
	msg := RedisMessage{
		InstanceID: h.instanceID,
		Type:       "broadcast_global",
		Event:      event,
		Data:       dataBytes,
	}
	msgBytes, _ := json.Marshal(msg)
	return h.redisClient.Publish(ctx, "websocket:hub:global", msgBytes).Err()
}

// SubscribeToRoomChannel subscribes to a room-specific Redis channel for targeted delivery.
func (h *Hub) SubscribeToRoomChannel(ctx context.Context, room string) error {
	if h.redisClient == nil {
		return nil
	}

	channel := "websocket:room:" + room
	pubsub := h.redisClient.Subscribe(ctx, channel)
	h.redisChannels[room] = pubsub

	go func() {
		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				h.handleRoomChannelMessage(room, msg.Payload)
			case <-h.ctx.Done():
				pubsub.Close()
				return
			}
		}
	}()

	return nil
}

// handleRoomChannelMessage handles messages from room-specific Redis channels.
func (h *Hub) handleRoomChannelMessage(room, payload string) {
	var msg RedisMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		h.log.Errorf("Failed to unmarshal room channel message: %v", err)
		return
	}

	if msg.InstanceID == h.instanceID {
		return
	}

	h.broadcastToRoomLocal(msg.Room, msg.Event, msg.Data)
}

// UnsubscribeFromRoomChannel unsubscribes from a room-specific Redis channel.
func (h *Hub) UnsubscribeFromRoomChannel(room string) {
	if pubsub, ok := h.redisChannels[room]; ok {
		pubsub.Close()
		delete(h.redisChannels, room)
	}
}

// Shutdown gracefully shuts down the hub and closes Redis connections.
func (h *Hub) Shutdown() {
	h.cancel()

	if h.pubsub != nil {
		h.pubsub.Close()
	}

	for _, pubsub := range h.redisChannels {
		pubsub.Close()
	}

	if h.queue != nil {
		h.queue.Close()
	}
}
