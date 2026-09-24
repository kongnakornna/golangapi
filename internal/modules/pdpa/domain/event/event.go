package event

type WebSocketEvent struct {
	Type    string
	Payload interface{}
}