package ws

import (
	"net/http"

	"icmongolang/internal/modules/kafka/usecase"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"

	gorillaWS "github.com/gorilla/websocket"
)

var upgrader = gorillaWS.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WsHandler struct {
	hub     *websocket.Hub
	usecase usecase.OrderUsecase
	log     logger.Logger
}

func NewWsHandler(hub *websocket.Hub, uc usecase.OrderUsecase, log logger.Logger) *WsHandler {
	return &WsHandler{
		hub:     hub,
		usecase: uc,
		log:     log,
	}
}

// ServeWS upgrades HTTP to WebSocket and handles incoming messages
func (h *WsHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Errorf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	room := r.URL.Query().Get("room")

	// สร้างและลงทะเบียน client ผ่าน Hub
	client := h.hub.AddClient(conn, userID, room)

	// เมื่อออกจากฟังก์ชัน ให้ยกเลิกการลงทะเบียน client
	//defer h.hub.RemoveClient(client)

	// เริ่ม write pump เพื่อส่งข้อความไปยัง client
	go client.WritePump()

	// อ่านข้อความจาก client และส่งไปยัง usecase (ซึ่งจะ publish ไปยัง Kafka)
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := h.usecase.HandleWebSocketMessage(client, msgBytes); err != nil {
			h.log.Errorf("HandleWebSocketMessage error: %v", err)
		}
	}
}
