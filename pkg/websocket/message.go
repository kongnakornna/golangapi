package websocket

import "encoding/json"

// Envelope defines the wire format for WebSocket communication.
type Envelope struct {
	Type    string          `json:"type"` // "subscribe", "unsubscribe", "message", "error"
	Topic   string          `json:"topic,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Error   string          `json:"error,omitempty"`
}
