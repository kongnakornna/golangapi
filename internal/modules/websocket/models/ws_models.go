package models

import (
	"encoding/json"
	"time"
)

// WSMessage represents a message exchanged over WebSocket.
type WSMessage struct {
	ID         string          `json:"id,omitempty"`
	Topic      string          `json:"topic"`
	Room       string          `json:"room,omitempty"`
	Payload    json.RawMessage `json:"payload"`
	SenderID   string          `json:"sender_id,omitempty"`
	SenderRole string          `json:"sender_role,omitempty"`
	SentAt     time.Time       `json:"sent_at,omitempty"`
}

// Session represents a user session.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// MessageQueryParams holds parameters for querying messages with pagination.
type MessageQueryParams struct {
	Topic    string     `json:"topic,omitempty"`
	Room     string     `json:"room,omitempty"`
	SenderID string     `json:"sender_id,omitempty"`
	Before   *time.Time `json:"before,omitempty"` // Cursor-based pagination: messages before this time
	After    *time.Time `json:"after,omitempty"`  // Cursor-based pagination: messages after this time
	Limit    int        `json:"limit,omitempty"`  // Max number of messages to return
	Offset   int        `json:"offset,omitempty"` // Offset-based pagination
	Order    string     `json:"order,omitempty"`  // "asc" or "desc"
}

// MessagePage represents a paginated result of messages.
type MessagePage struct {
	Messages   []*WSMessage `json:"messages"`
	TotalCount int64        `json:"total_count"`
	HasMore    bool         `json:"has_more"`
	NextCursor *time.Time   `json:"next_cursor,omitempty"`
	PrevCursor *time.Time   `json:"prev_cursor,omitempty"`
}

// SearchQueryParams holds parameters for full-text search on messages.
type SearchQueryParams struct {
	Query     string     `json:"query"`               // Search query string
	Topic     string     `json:"topic,omitempty"`     // Filter by topic
	Room      string     `json:"room,omitempty"`      // Filter by room
	SenderID  string     `json:"sender_id,omitempty"` // Filter by sender
	Before    *time.Time `json:"before,omitempty"`    // Messages before this time
	After     *time.Time `json:"after,omitempty"`     // Messages after this time
	Limit     int        `json:"limit,omitempty"`     // Max results (default 50, max 1000)
	Offset    int        `json:"offset,omitempty"`    // Offset for pagination
	Order     string     `json:"order,omitempty"`     // "asc" or "desc" (default desc by rank)
	Highlight bool       `json:"highlight,omitempty"` // Include highlighted snippets
}

// SearchResult represents a single search result with optional highlights.
type SearchResult struct {
	Message    *WSMessage          `json:"message"`
	Highlights map[string][]string `json:"highlights,omitempty"` // Field -> highlighted snippets
	Rank       float64             `json:"rank,omitempty"`       // Search relevance rank
}

// SearchPage represents a paginated search result.
type SearchPage struct {
	Results    []*SearchResult `json:"results"`
	TotalCount int64           `json:"total_count"`
	HasMore    bool            `json:"has_more"`
	NextCursor *time.Time      `json:"next_cursor,omitempty"`
	PrevCursor *time.Time      `json:"prev_cursor,omitempty"`
}

// RetentionPolicy defines the message retention policy for a topic/room.
type RetentionPolicy struct {
	ID        string        `json:"id,omitempty"`
	Topic     string        `json:"topic,omitempty"`     // Topic pattern (supports wildcards)
	Room      string        `json:"room,omitempty"`      // Room pattern (supports wildcards)
	MaxAge    time.Duration `json:"max_age"`             // Maximum age of messages to retain
	MaxCount  int           `json:"max_count,omitempty"` // Maximum number of messages to retain per topic/room
	Enabled   bool          `json:"enabled"`             // Whether the policy is active
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// RetentionPolicyRequest represents the request to create/update a retention policy.
type RetentionPolicyRequest struct {
	Topic    string        `json:"topic,omitempty" validate:"required"` // Topic pattern (supports wildcards)
	Room     string        `json:"room,omitempty"`                      // Room pattern (supports wildcards)
	MaxAge   time.Duration `json:"max_age" validate:"required"`         // Maximum age of messages to retain
	MaxCount int           `json:"max_count,omitempty"`                 // Maximum number of messages to retain per topic/room
	Enabled  bool          `json:"enabled"`                             // Whether the policy is active
}

// RetentionPolicyResponse represents a retention policy response.
type RetentionPolicyResponse struct {
	ID        string        `json:"id"`
	Topic     string        `json:"topic,omitempty"`
	Room      string        `json:"room,omitempty"`
	MaxAge    time.Duration `json:"max_age"`
	MaxCount  int           `json:"max_count,omitempty"`
	Enabled   bool          `json:"enabled"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// RetentionPolicyListResponse represents a list of retention policies.
type RetentionPolicyListResponse struct {
	Policies []RetentionPolicyResponse `json:"policies"`
	Total    int64                     `json:"total"`
}

// RetentionStats represents statistics about message retention.
type RetentionStats struct {
	TotalMessages   int64     `json:"total_messages"`
	MessagesDeleted int64     `json:"messages_deleted"`
	BytesFreed      int64     `json:"bytes_freed"`
	LastCleanupAt   time.Time `json:"last_cleanup_at"`
	NextCleanupAt   time.Time `json:"next_cleanup_at"`
	PoliciesCount   int       `json:"policies_count"`
	ActivePolicies  int       `json:"active_policies"`
}
