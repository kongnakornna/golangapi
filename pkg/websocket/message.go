package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// ============================================================================
// Message Type Constants
// ============================================================================

// Client-to-Server message types
const (
	MsgTypeSubscribe   = "subscribe"
	MsgTypeUnsubscribe = "unsubscribe"
	MsgTypeJoinRoom    = "join_room"
	MsgTypeLeaveRoom   = "leave_room"
	MsgTypeMessage     = "message"
	MsgTypePing        = "ping"
	MsgTypePong        = "pong"
	MsgTypeAuth        = "auth"
	MsgTypePresence    = "presence"
	MsgTypeTyping      = "typing"
	MsgTypeReadReceipt = "read_receipt"
	MsgTypeAck         = "ack"
)

// Server-to-Client event types
const (
	EventTypeMessage      = "message"
	EventTypeJoined       = "joined"
	EventTypeLeft         = "left"
	EventTypeSubscribed   = "subscribed"
	EventTypeUnsubscribed = "unsubscribed"
	EventTypeError        = "error"
	EventTypeAck          = "ack"
	EventTypePresence     = "presence"
	EventTypeTyping       = "typing"
	EventTypeReadReceipt  = "read_receipt"
	EventTypeSystem       = "system"
	EventTypeAuth         = "auth"
	EventTypePong         = "pong"
)

// ============================================================================
// Base Envelope
// ============================================================================

// Envelope is the wire format for all WebSocket communication.
type Envelope struct {
	ID        string          `json:"id,omitempty"`        // Client-generated correlation ID
	Type      string          `json:"type"`                // Message/event type
	Topic     string          `json:"topic,omitempty"`     // Topic for pub/sub
	Room      string          `json:"room,omitempty"`      // Room for room-based messaging
	Payload   json.RawMessage `json:"payload,omitempty"`   // Typed payload
	Error     *ErrorPayload   `json:"error,omitempty"`     // Error details (server->client)
	Timestamp int64           `json:"timestamp,omitempty"` // Unix milliseconds
	Seq       uint64          `json:"seq,omitempty"`       // Server-assigned sequence number for ordering
	SeqAck    uint64          `json:"seq_ack,omitempty"`   // Acknowledged sequence number (client->server)
}

// ErrorPayload provides structured error information.
type ErrorPayload struct {
	Code    string `json:"code"`              // Machine-readable error code
	Message string `json:"message"`           // Human-readable message
	Details any    `json:"details,omitempty"` // Additional context
}

// NewEnvelope creates a new envelope with timestamp.
func NewEnvelope(msgType string) *Envelope {
	return &Envelope{
		Type:      msgType,
		Timestamp: time.Now().UnixMilli(),
	}
}

// WithID sets the correlation ID.
func (e *Envelope) WithID(id string) *Envelope {
	e.ID = id
	return e
}

// WithTopic sets the topic.
func (e *Envelope) WithTopic(topic string) *Envelope {
	e.Topic = topic
	return e
}

// WithRoom sets the room.
func (e *Envelope) WithRoom(room string) *Envelope {
	e.Room = room
	return e
}

// WithPayload sets the payload (marshalled to JSON).
func (e *Envelope) WithPayload(v any) *Envelope {
	data, _ := json.Marshal(v)
	e.Payload = data
	return e
}

// WithError sets the error payload.
func (e *Envelope) WithError(code, message string, details any) *Envelope {
	e.Error = &ErrorPayload{Code: code, Message: message, Details: details}
	return e
}

// WithSeq sets the sequence number.
func (e *Envelope) WithSeq(seq uint64) *Envelope {
	e.Seq = seq
	return e
}

// WithSeqAck sets the acknowledged sequence number.
func (e *Envelope) WithSeqAck(seqAck uint64) *Envelope {
	e.SeqAck = seqAck
	return e
}

// ============================================================================
// Client-to-Server Payload Types
// ============================================================================

// SubscribePayload for subscribing to a topic.
type SubscribePayload struct {
	Topic string `json:"topic" validate:"required"`
}

// UnsubscribePayload for unsubscribing from a topic.
type UnsubscribePayload struct {
	Topic string `json:"topic" validate:"required"`
}

// JoinRoomPayload for joining a room.
type JoinRoomPayload struct {
	Room string `json:"room" validate:"required"`
}

// LeaveRoomPayload for leaving a room.
type LeaveRoomPayload struct {
	Room string `json:"room" validate:"required"`
}

// MessagePayload for sending a message to topic or room.
type MessagePayload struct {
	Topic   string          `json:"topic,omitempty"`
	Room    string          `json:"room,omitempty"`
	Content json.RawMessage `json:"content" validate:"required"`
	Meta    map[string]any  `json:"meta,omitempty"`
}

// AuthPayload for authentication.
type AuthPayload struct {
	Token string `json:"token" validate:"required"`
}

// PresencePayload for presence updates.
type PresencePayload struct {
	Status string         `json:"status"` // "online", "away", "busy", "offline"
	Meta   map[string]any `json:"meta,omitempty"`
}

// PresenceStatus represents the possible presence states.
type PresenceStatus string

const (
	PresenceOnline  PresenceStatus = "online"
	PresenceAway    PresenceStatus = "away"
	PresenceBusy    PresenceStatus = "busy"
	PresenceOffline PresenceStatus = "offline"
)

// PresenceInfo holds detailed presence information for a user.
type PresenceInfo struct {
	UserID    string         `json:"user_id"`
	Status    PresenceStatus `json:"status"`
	Meta      map[string]any `json:"meta,omitempty"`
	LastSeen  time.Time      `json:"last_seen"`
	Instances []string       `json:"instances,omitempty"` // Instance IDs where user is connected
}

// PresenceQueryParams for querying presence.
type PresenceQueryParams struct {
	UserIDs []string       `json:"user_ids,omitempty"`
	Status  PresenceStatus `json:"status,omitempty"`
	Limit   int            `json:"limit,omitempty"`
	Offset  int            `json:"offset,omitempty"`
}

// PresencePage represents a paginated presence result.
type PresencePage struct {
	Presences []*PresenceInfo `json:"presences"`
	Total     int64           `json:"total"`
	HasMore   bool            `json:"has_more"`
}

// TypingPayload for typing indicators.
type TypingPayload struct {
	Room     string `json:"room,omitempty"`
	Topic    string `json:"topic,omitempty"`
	IsTyping bool   `json:"is_typing"`
}

// ReadReceiptPayload for read receipts.
type ReadReceiptPayload struct {
	MessageIDs []string `json:"message_ids" validate:"required"`
	Room       string   `json:"room,omitempty"`
	Topic      string   `json:"topic,omitempty"`
}

// AckPayload for acknowledgments.
type AckPayload struct {
	OriginalID string `json:"original_id" validate:"required"`
	Status     string `json:"status"` // "ok", "failed"
	Error      string `json:"error,omitempty"`
}

// AckConfig holds configuration for acknowledgment behavior.
type AckConfig struct {
	// Timeout is how long to wait for an ACK before retrying.
	Timeout time.Duration
	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int
	// RetryInterval is the interval between retries.
	RetryInterval time.Duration
}

// DefaultAckConfig returns sensible defaults.
func DefaultAckConfig() *AckConfig {
	return &AckConfig{
		Timeout:       5 * time.Second,
		MaxRetries:    3,
		RetryInterval: 1 * time.Second,
	}
}

// PendingAck represents a message waiting for acknowledgment.
type PendingAck struct {
	ID           string
	Envelope     *Envelope
	Retries      int
	LastSent     time.Time
	ResponseChan chan *Envelope // receives the ACK response
	CancelFunc   context.CancelFunc
}

// ============================================================================
// Server-to-Client Payload Types
// ============================================================================

// MessageEventPayload for incoming messages.
type MessageEventPayload struct {
	ID         string          `json:"id"`
	Topic      string          `json:"topic,omitempty"`
	Room       string          `json:"room,omitempty"`
	SenderID   string          `json:"sender_id"`
	SenderRole string          `json:"sender_role,omitempty"`
	Content    json.RawMessage `json:"content"`
	Meta       map[string]any  `json:"meta,omitempty"`
	Timestamp  int64           `json:"timestamp"`
}

// JoinedEventPayload for successful room/topic join.
type JoinedEventPayload struct {
	Room  string `json:"room,omitempty"`
	Topic string `json:"topic,omitempty"`
	Count int    `json:"count"` // Current member count
}

// LeftEventPayload for room/topic leave.
type LeftEventPayload struct {
	Room  string `json:"room,omitempty"`
	Topic string `json:"topic,omitempty"`
	Count int    `json:"count"` // Remaining member count
}

// SubscribedEventPayload for successful subscription.
type SubscribedEventPayload struct {
	Topic string `json:"topic"`
}

// UnsubscribedEventPayload for successful unsubscription.
type UnsubscribedEventPayload struct {
	Topic string `json:"topic"`
}

// PresenceEventPayload for presence updates.
type PresenceEventPayload struct {
	UserID string         `json:"user_id"`
	Status string         `json:"status"` // "online", "away", "busy", "offline"
	Meta   map[string]any `json:"meta,omitempty"`
}

// TypingEventPayload for typing indicators.
type TypingEventPayload struct {
	UserID   string `json:"user_id"`
	Room     string `json:"room,omitempty"`
	Topic    string `json:"topic,omitempty"`
	IsTyping bool   `json:"is_typing"`
}

// ReadReceiptEventPayload for read receipts.
type ReadReceiptEventPayload struct {
	UserID     string   `json:"user_id"`
	MessageIDs []string `json:"message_ids"`
	Room       string   `json:"room,omitempty"`
	Topic      string   `json:"topic,omitempty"`
}

// AckEventPayload for acknowledgments.
type AckEventPayload struct {
	OriginalID string `json:"original_id"`
	Status     string `json:"status"` // "ok", "failed"
	Error      string `json:"error,omitempty"`
}

// SystemEventPayload for system notifications.
type SystemEventPayload struct {
	Level   string         `json:"level"` // "info", "warning", "critical"
	Message string         `json:"message"`
	Meta    map[string]any `json:"meta,omitempty"`
}

// AuthEventPayload for auth responses.
type AuthEventPayload struct {
	Success bool   `json:"success"`
	UserID  string `json:"user_id,omitempty"`
	Role    string `json:"role,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ============================================================================
// Helper Functions
// ============================================================================

// MarshalPayload marshals a payload to json.RawMessage.
func MarshalPayload(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

// UnmarshalPayload unmarshals payload into target.
func UnmarshalPayload(data json.RawMessage, v any) error {
	if len(data) == 0 {
		return errors.New("empty payload")
	}
	return json.Unmarshal(data, v)
}

// GenerateID generates a unique message ID.
func GenerateID() string {
	return uuid.New().String()
}

// ============================================================================
// Envelope Builders (Fluent API)
// ============================================================================

// BuildSubscribe creates a subscribe envelope.
func BuildSubscribe(topic string) *Envelope {
	return NewEnvelope(MsgTypeSubscribe).
		WithID(GenerateID()).
		WithTopic(topic).
		WithPayload(SubscribePayload{Topic: topic})
}

// BuildUnsubscribe creates an unsubscribe envelope.
func BuildUnsubscribe(topic string) *Envelope {
	return NewEnvelope(MsgTypeUnsubscribe).
		WithID(GenerateID()).
		WithTopic(topic).
		WithPayload(UnsubscribePayload{Topic: topic})
}

// BuildJoinRoom creates a join_room envelope.
func BuildJoinRoom(room string) *Envelope {
	return NewEnvelope(MsgTypeJoinRoom).
		WithID(GenerateID()).
		WithRoom(room).
		WithPayload(JoinRoomPayload{Room: room})
}

// BuildLeaveRoom creates a leave_room envelope.
func BuildLeaveRoom(room string) *Envelope {
	return NewEnvelope(MsgTypeLeaveRoom).
		WithID(GenerateID()).
		WithRoom(room).
		WithPayload(LeaveRoomPayload{Room: room})
}

// BuildMessage creates a message envelope.
func BuildMessage(topic, room string, content any, meta map[string]any) *Envelope {
	return NewEnvelope(MsgTypeMessage).
		WithID(GenerateID()).
		WithTopic(topic).
		WithRoom(room).
		WithPayload(MessagePayload{
			Topic:   topic,
			Room:    room,
			Content: MarshalPayload(content),
			Meta:    meta,
		})
}

// BuildAuth creates an auth envelope.
func BuildAuth(token string) *Envelope {
	return NewEnvelope(MsgTypeAuth).
		WithID(GenerateID()).
		WithPayload(AuthPayload{Token: token})
}

// BuildPing creates a ping envelope.
func BuildPing() *Envelope {
	return NewEnvelope(MsgTypePing).
		WithID(GenerateID())
}

// BuildPong creates a pong envelope.
func BuildPong() *Envelope {
	return NewEnvelope(MsgTypePong).
		WithID(GenerateID())
}

// BuildPresence creates a presence envelope.
func BuildPresence(status string, meta map[string]any) *Envelope {
	return NewEnvelope(MsgTypePresence).
		WithID(GenerateID()).
		WithPayload(PresencePayload{Status: status, Meta: meta})
}

// BuildTyping creates a typing envelope.
func BuildTyping(room, topic string, isTyping bool) *Envelope {
	return NewEnvelope(MsgTypeTyping).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(TypingPayload{Room: room, Topic: topic, IsTyping: isTyping})
}

// BuildReadReceipt creates a read_receipt envelope.
func BuildReadReceipt(messageIDs []string, room, topic string) *Envelope {
	return NewEnvelope(MsgTypeReadReceipt).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(ReadReceiptPayload{MessageIDs: messageIDs, Room: room, Topic: topic})
}

// BuildAck creates an ack envelope.
func BuildAck(originalID, status, errMsg string) *Envelope {
	return NewEnvelope(MsgTypeAck).
		WithID(GenerateID()).
		WithPayload(AckPayload{OriginalID: originalID, Status: status, Error: errMsg})
}

// ============================================================================
// Server Event Builders
// ============================================================================

// BuildMessageEvent creates a message event envelope.
func BuildMessageEvent(id, topic, room, senderID, senderRole string, content any, meta map[string]any) *Envelope {
	return NewEnvelope(EventTypeMessage).
		WithID(GenerateID()).
		WithTopic(topic).
		WithRoom(room).
		WithPayload(MessageEventPayload{
			ID:         id,
			Topic:      topic,
			Room:       room,
			SenderID:   senderID,
			SenderRole: senderRole,
			Content:    MarshalPayload(content),
			Meta:       meta,
			Timestamp:  time.Now().UnixMilli(),
		})
}

// BuildJoinedEvent creates a joined event envelope.
func BuildJoinedEvent(room, topic string, count int) *Envelope {
	return NewEnvelope(EventTypeJoined).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(JoinedEventPayload{Room: room, Topic: topic, Count: count})
}

// BuildLeftEvent creates a left event envelope.
func BuildLeftEvent(room, topic string, count int) *Envelope {
	return NewEnvelope(EventTypeLeft).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(LeftEventPayload{Room: room, Topic: topic, Count: count})
}

// BuildSubscribedEvent creates a subscribed event envelope.
func BuildSubscribedEvent(topic string) *Envelope {
	return NewEnvelope(EventTypeSubscribed).
		WithID(GenerateID()).
		WithTopic(topic).
		WithPayload(SubscribedEventPayload{Topic: topic})
}

// BuildUnsubscribedEvent creates an unsubscribed event envelope.
func BuildUnsubscribedEvent(topic string) *Envelope {
	return NewEnvelope(EventTypeUnsubscribed).
		WithID(GenerateID()).
		WithTopic(topic).
		WithPayload(UnsubscribedEventPayload{Topic: topic})
}

// BuildErrorEvent creates an error event envelope.
func BuildErrorEvent(code, message string, details any) *Envelope {
	return NewEnvelope(EventTypeError).
		WithID(GenerateID()).
		WithError(code, message, details)
}

// BuildAckEvent creates an ack event envelope.
func BuildAckEvent(originalID, status, errMsg string) *Envelope {
	return NewEnvelope(EventTypeAck).
		WithID(GenerateID()).
		WithPayload(AckEventPayload{OriginalID: originalID, Status: status, Error: errMsg})
}

// BuildPresenceEvent creates a presence event envelope.
func BuildPresenceEvent(userID, status string, meta map[string]any) *Envelope {
	return NewEnvelope(EventTypePresence).
		WithID(GenerateID()).
		WithPayload(PresenceEventPayload{UserID: userID, Status: status, Meta: meta})
}

// BuildTypingEvent creates a typing event envelope.
func BuildTypingEvent(userID, room, topic string, isTyping bool) *Envelope {
	return NewEnvelope(EventTypeTyping).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(TypingEventPayload{UserID: userID, Room: room, Topic: topic, IsTyping: isTyping})
}

// BuildReadReceiptEvent creates a read_receipt event envelope.
func BuildReadReceiptEvent(userID string, messageIDs []string, room, topic string) *Envelope {
	return NewEnvelope(EventTypeReadReceipt).
		WithID(GenerateID()).
		WithRoom(room).
		WithTopic(topic).
		WithPayload(ReadReceiptEventPayload{UserID: userID, MessageIDs: messageIDs, Room: room, Topic: topic})
}

// BuildSystemEvent creates a system event envelope.
func BuildSystemEvent(level, message string, meta map[string]any) *Envelope {
	return NewEnvelope(EventTypeSystem).
		WithID(GenerateID()).
		WithPayload(SystemEventPayload{Level: level, Message: message, Meta: meta})
}

// BuildAuthEvent creates an auth event envelope.
func BuildAuthEvent(success bool, userID, role, errMsg string) *Envelope {
	return NewEnvelope(EventTypeAuth).
		WithID(GenerateID()).
		WithPayload(AuthEventPayload{Success: success, UserID: userID, Role: role, Error: errMsg})
}

// BuildPongEvent creates a pong event envelope.
func BuildPongEvent() *Envelope {
	return NewEnvelope(EventTypePong).
		WithID(GenerateID())
}

// ============================================================================
// Acknowledgment Manager
// ============================================================================

// AckManager manages pending acknowledgments with retry logic.
type AckManager struct {
	mu      sync.RWMutex
	pending map[string]*PendingAck
	config  *AckConfig
	logger  logger.Logger
}

// NewAckManager creates a new acknowledgment manager.
func NewAckManager(config *AckConfig, logger logger.Logger) *AckManager {
	if config == nil {
		config = DefaultAckConfig()
	}
	return &AckManager{
		pending: make(map[string]*PendingAck),
		config:  config,
		logger:  logger,
	}
}

// SendWithAck sends an envelope and waits for acknowledgment with retry logic.
// Returns the ACK envelope or an error if all retries exhausted.
func (m *AckManager) SendWithAck(ctx context.Context, client *Client, env *Envelope) (*Envelope, error) {
	if env.ID == "" {
		env.ID = GenerateID()
	}

	// Create pending ack entry
	responseChan := make(chan *Envelope, 1)
	ctx, cancel := context.WithCancel(ctx)
	pending := &PendingAck{
		ID:           env.ID,
		Envelope:     env,
		Retries:      0,
		LastSent:     time.Now(),
		ResponseChan: responseChan,
		CancelFunc:   cancel,
	}

	m.mu.Lock()
	m.pending[env.ID] = pending
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.pending, env.ID)
		m.mu.Unlock()
	}()

	// Send initial message
	client.SendMessage(MarshalPayload(env))

	// Wait for ACK or retry
	for {
		select {
		case ack := <-responseChan:
			if ack.Error != nil {
				return nil, fmt.Errorf("ack error: %s", ack.Error.Message)
			}
			return ack, nil

		case <-time.After(m.config.Timeout):
			m.mu.Lock()
			pending, ok := m.pending[env.ID]
			m.mu.Unlock()

			if !ok {
				return nil, errors.New("ack canceled")
			}

			pending.Retries++
			if pending.Retries >= m.config.MaxRetries {
				return nil, fmt.Errorf("max retries (%d) exceeded", m.config.MaxRetries)
			}

			// Retry
			pending.LastSent = time.Now()
			client.SendMessage(MarshalPayload(env))

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// HandleAck processes an incoming acknowledgment.
func (m *AckManager) HandleAck(ack *Envelope) {
	if ack.Payload == nil {
		return
	}

	var payload AckEventPayload
	if err := UnmarshalPayload(ack.Payload, &payload); err != nil {
		return
	}

	m.mu.Lock()
	pending, ok := m.pending[payload.OriginalID]
	if ok {
		delete(m.pending, payload.OriginalID)
	}
	m.mu.Unlock()

	if ok && pending.ResponseChan != nil {
		select {
		case pending.ResponseChan <- ack:
		default:
		}
		if pending.CancelFunc != nil {
			pending.CancelFunc()
		}
	}
}

// CancelPending cancels a pending acknowledgment.
func (m *AckManager) CancelPending(id string) {
	m.mu.Lock()
	if pending, ok := m.pending[id]; ok {
		if pending.CancelFunc != nil {
			pending.CancelFunc()
		}
		delete(m.pending, id)
	}
	m.mu.Unlock()
}

// PendingCount returns the number of pending acknowledgments.
func (m *AckManager) PendingCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.pending)
}
