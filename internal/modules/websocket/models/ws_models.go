package models

import (
	"encoding/json"
	"time"
)

// WSMessage represents a message exchanged over WebSocket.
type WSMessage struct {
	ID       string          `json:"id,omitempty"`
	Topic    string          `json:"topic"`
	Payload  json.RawMessage `json:"payload"`
	SenderID string          `json:"sender_id,omitempty"`
	SentAt   time.Time       `json:"sent_at,omitempty"`
}

// Session represents a user session.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
