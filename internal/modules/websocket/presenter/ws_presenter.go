package presenter

// SendMessageRequest represents a request to broadcast a message to a room.
type SendMessageRequest struct {
	Room  string      `json:"room" example:"BAACTW01"`
	Event string      `json:"event" example:"message"`
	Data  interface{} `json:"data"`
}

// SendMessageResponse is the response after broadcasting to a room.
type SendMessageResponse struct {
	Status string `json:"status" example:"sent"`
}

// RoomsResponse lists active WS rooms.
type RoomsResponse struct {
	Rooms []string `json:"rooms" example:"BAACTW01"`
	Count int      `json:"count" example:"3"`
}

// RoomStatsResponse returns the number of clients in a room.
type RoomStatsResponse struct {
	Room    string `json:"room" example:"BAACTW01"`
	Clients int    `json:"clients" example:"2"`
}

// WSMessageItem is a single persisted message history item.
type WSMessageItem struct {
	ID         string      `json:"id,omitempty"`
	Topic      string      `json:"topic" example:"BAACTW01/DATA"`
	Room       string      `json:"room,omitempty" example:"BAACTW01"`
	Payload    interface{} `json:"payload"`
	SenderID   string      `json:"sender_id,omitempty"`
	SenderRole string      `json:"sender_role,omitempty"`
	SentAt     string      `json:"sent_at,omitempty" example:"2024-01-15T10:30:00+07:00"`
}

// MessageQueryRequest represents query parameters for message history pagination.
type MessageQueryRequest struct {
	Topic    string `query:"topic" example:"BAACTW01/DATA"`
	Room     string `query:"room" example:"BAACTW01"`
	SenderID string `query:"sender_id" example:"user123"`
	Before   string `query:"before" example:"2024-01-15T10:30:00Z"` // RFC3339 timestamp
	After    string `query:"after" example:"2024-01-15T10:00:00Z"`  // RFC3339 timestamp
	Limit    int    `query:"limit" example:"50"`
	Offset   int    `query:"offset" example:"0"`
	Order    string `query:"order" example:"desc" enums:"asc,desc"`
}

// MessagePageResponse represents a paginated response of messages.
type MessagePageResponse struct {
	Messages   []WSMessageItem `json:"messages"`
	TotalCount int64           `json:"total_count"`
	HasMore    bool            `json:"has_more"`
	NextCursor string          `json:"next_cursor,omitempty" example:"2024-01-15T10:29:00Z"`
	PrevCursor string          `json:"prev_cursor,omitempty" example:"2024-01-15T10:31:00Z"`
}

// PresenceItem represents a user's presence status.
type PresenceItem struct {
	UserID    string         `json:"user_id" example:"user123"`
	Status    string         `json:"status" example:"online" enums:"online,away,busy,offline"`
	Meta      map[string]any `json:"meta,omitempty"`
	LastSeen  string         `json:"last_seen" example:"2024-01-15T10:30:00Z"`
	Instances []string       `json:"instances,omitempty" example:"[\"abc123\",\"def456\"]"`
}

// PresenceQueryRequest represents query parameters for presence queries.
type PresenceQueryRequest struct {
	UserIDs []string `query:"user_ids" example:"[\"user1\",\"user2\"]"`
	Status  string   `query:"status" example:"online" enums:"online,away,busy,offline"`
	Limit   int      `query:"limit" example:"50"`
	Offset  int      `query:"offset" example:"0"`
}

// PresencePageResponse represents a paginated response of presences.
type PresencePageResponse struct {
	Presences []PresenceItem `json:"presences"`
	Total     int64          `json:"total"`
	HasMore   bool           `json:"has_more"`
}

// SearchQueryRequest represents query parameters for message search.
type SearchQueryRequest struct {
	Query     string `query:"query" example:"error timeout"`
	Topic     string `query:"topic" example:"BAACTW01/DATA"`
	Room      string `query:"room" example:"BAACTW01"`
	SenderID  string `query:"sender_id" example:"user123"`
	Before    string `query:"before" example:"2024-01-15T10:30:00Z"` // RFC3339 timestamp
	After     string `query:"after" example:"2024-01-15T10:00:00Z"`  // RFC3339 timestamp
	Limit     int    `query:"limit" example:"50"`
	Offset    int    `query:"offset" example:"0"`
	Order     string `query:"order" example:"desc" enums:"asc,desc"`
	Highlight bool   `query:"highlight" example:"true"`
}

// SearchResultItem represents a single search result with optional highlights.
type SearchResultItem struct {
	Message    WSMessageItem       `json:"message"`
	Highlights map[string][]string `json:"highlights,omitempty"`
	Rank       float64             `json:"rank,omitempty"`
}

// SearchPageResponse represents a paginated search result.
type SearchPageResponse struct {
	Results    []SearchResultItem `json:"results"`
	TotalCount int64              `json:"total_count"`
	HasMore    bool               `json:"has_more"`
	NextCursor string             `json:"next_cursor,omitempty" example:"2024-01-15T10:29:00Z"`
	PrevCursor string             `json:"prev_cursor,omitempty" example:"2024-01-15T10:31:00Z"`
}

// RetentionPolicyRequest represents a request to create/update a retention policy.
type RetentionPolicyRequest struct {
	Topic    string `json:"topic,omitempty" example:"BAACTW01/*" validate:"required"`
	Room     string `json:"room,omitempty" example:"BAACTW01"`
	MaxAge   string `json:"max_age" example:"168h" validate:"required"` // Duration string (e.g., "168h", "720h")
	MaxCount int    `json:"max_count,omitempty" example:"10000"`
	Enabled  bool   `json:"enabled" example:"true"`
}

// RetentionPolicyResponse represents a retention policy response.
type RetentionPolicyResponse struct {
	ID        string `json:"id" example:"pol-abc123"`
	Topic     string `json:"topic,omitempty" example:"BAACTW01/*"`
	Room      string `json:"room,omitempty" example:"BAACTW01"`
	MaxAge    string `json:"max_age" example:"168h"`
	MaxCount  int    `json:"max_count,omitempty" example:"10000"`
	Enabled   bool   `json:"enabled" example:"true"`
	CreatedAt string `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// RetentionPolicyListResponse represents a list of retention policies.
type RetentionPolicyListResponse struct {
	Policies []RetentionPolicyResponse `json:"policies"`
	Total    int64                     `json:"total"`
}

// RetentionStatsResponse represents statistics about message retention.
type RetentionStatsResponse struct {
	TotalMessages   int64  `json:"total_messages" example:"1000000"`
	MessagesDeleted int64  `json:"messages_deleted" example:"5000"`
	BytesFreed      int64  `json:"bytes_freed" example:"104857600"`
	LastCleanupAt   string `json:"last_cleanup_at,omitempty" example:"2024-01-15T10:30:00Z"`
	NextCleanupAt   string `json:"next_cleanup_at,omitempty" example:"2024-01-16T10:30:00Z"`
	PoliciesCount   int    `json:"policies_count" example:"5"`
	ActivePolicies  int    `json:"active_policies" example:"3"`
}

// CleanupRequest represents a request to run retention cleanup.
type CleanupRequest struct {
	PolicyID string `json:"policy_id,omitempty" example:"pol-abc123"` // If empty, runs all active policies
}

// CleanupResponse represents the response after running retention cleanup.
type CleanupResponse struct {
	Status          string `json:"status" example:"completed"`
	MessagesDeleted int64  `json:"messages_deleted" example:"5000"`
	BytesFreed      int64  `json:"bytes_freed" example:"104857600"`
	Duration        string `json:"duration" example:"1.5s"`
	PoliciesRun     int    `json:"policies_run" example:"3"`
}
