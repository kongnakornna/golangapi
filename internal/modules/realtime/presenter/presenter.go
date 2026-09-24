package presenter

// PublishRequest is an inbound real-time event to broadcast to dashboard rooms.
type PublishRequest struct {
	Room  string         `json:"room,omitempty"` // target room; empty = global broadcast
	Event string         `json:"event"`          // e.g. "monitor", "alarm", "device"
	Data  map[string]any `json:"data"`
}

// PublishResponse reports the broadcast result.
type PublishResponse struct {
	OK      bool   `json:"ok"`
	Room    string `json:"room,omitempty"`
	Event   string `json:"event"`
	Clients int    `json:"clients"`
}

// HealthResponse reports whether the realtime service is wired.
type HealthResponse struct {
	Enabled bool `json:"enabled"`
}
