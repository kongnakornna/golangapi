package websocket

// Broadcaster defines the interface for sending messages to WebSocket clients
type Broadcaster interface {
	BroadcastToRoom(room, event string, data interface{})
	BroadcastToTopic(topic, event string, data interface{})
	BroadcastMessage(event string, data interface{})
	GetRooms() []string
	GetClientsInRoom(room string) int
}
