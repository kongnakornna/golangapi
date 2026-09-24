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
	ID       string      `json:"id,omitempty"`
	Topic    string      `json:"topic" example:"BAACTW01/DATA"`
	Payload  interface{} `json:"payload"`
	SenderID string      `json:"sender_id,omitempty"`
	SentAt   string      `json:"sent_at,omitempty" example:"2024-01-15T10:30:00+07:00"`
}
