// Package ws จัดเตรียม WebSocket handler สำหรับโมดูล kafka
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

// ServeWS godoc
//
// @Summary      WebSocket endpoint สำหรับ real-time order updates
// @Description  อัปเกรด HTTP เป็น WebSocket เพื่อรับการแจ้งเตือนสถานะ order แบบ real-time
// @Description
// @Description  **Flow:**
// @Description  1. Connect: `ws://<host>/ws?user_id=<user>&room=<room>`
// @Description  2. Server จะส่งข้อความ JSON envelope: `{"event":"...","data":...}`
// @Description  3. Client ส่งข้อความ JSON กลับได้ — server จะ route ผ่าน usecase (เช่น publish ไป Kafka)
// @Description
// @Description  **Query Parameters:**
// @Description  - `user_id` (optional) — ระบุตัวตนผู้ใช้ (default: `anonymous`)
// @Description  - `room`    (optional) — เข้าห้องที่ระบุ (เช่น `order_id`) เพื่อรับ event เฉพาะห้อง
// @Description
// @Description  **Close Codes:**
// @Description  - `1000` — ปิดปกติ
// @Description  - `1001` — server กำลัง shutdown
// @Description  - `1011` — server error
// @Tags         Websocket
// @Accept       json
// @Produce      json
// @Param        user_id query string false "User ID (default anonymous)" example(user123)
// @Param        room    query string false "Room to join (e.g. order_id)" example(BAACTW01)
// @Success      101 {string} string "Switching Protocols — WebSocket handshake succeeded"
// @Failure      400 {string} string "Bad Request — invalid WebSocket upgrade headers"
// @Failure      500 {string} string "Internal Server Error — upgrade failed"
// @Security     BearerAuth
// @Router       /ws [get]
func (h *WsHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Errorf("WebSocket upgrade error: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	room := r.URL.Query().Get("room")

	// สร้างและลงทะเบียน client ผ่าน Hub
	client := h.hub.AddClient(conn, userID, room)
	h.log.Infof("websocket client connected: user=%s room=%s", userID, room)

	// เริ่ม write pump เพื่อส่งข้อความไปยัง client
	go client.WritePump()

	// อ่านข้อความจาก client และส่งไปยัง usecase (ซึ่งจะ publish ไปยัง Kafka)
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			h.log.Infof("websocket client disconnected: user=%s room=%s err=%v", userID, room, err)
			break
		}
		if err := h.usecase.HandleWebSocketMessage(client, msgBytes); err != nil {
			h.log.Errorf("HandleWebSocketMessage error: %v", err)
		}
	}
}
