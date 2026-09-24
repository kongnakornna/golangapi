package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"icmongolang/internal/modules/websocket/usecase"
	"icmongolang/pkg/logger"
	pkgws "icmongolang/pkg/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsHandler struct {
	hub     *pkgws.Hub
	usecase usecase.WSUsecase
	log     logger.Logger
}

func NewWsHandler(hub *pkgws.Hub, uc usecase.WSUsecase, log logger.Logger) *WsHandler {
	return &WsHandler{hub: hub, usecase: uc, log: log}
}

// ServeHTTP godoc
// @Summary      WebSocket endpoint (real-time MQTT/IoT stream)
// @Description  Upgrades an HTTP connection to WebSocket. Clients send JSON envelopes with "type": subscribe | unsubscribe | join_room | leave_room | message and a "topic" or "room". MQTT messages are broadcast to clients joined to the room named by the first segment of the MQTT topic (e.g. topic BAACTW05/DATA -> room BAACTW05); queue messages are delivered to topic subscribers. Server frames look like {"event":"...","data":...}.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        token query string false "Bearer token for authentication (or Authorization header)"
// @Success      101 {string} string "Switching Protocols"
// @Failure      401 {string} string "unauthorized"
// @Router       /ws [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
	}

	// usecase เป็น nil ได้ (เช่นใน main API server ที่ไม่ต้องบันทึก history)
	userID := "anonymous"
	if h.usecase != nil {
		uid, err := h.usecase.Authenticate(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID = uid
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Errorf("Upgrade error: %v", err)
		return
	}

	client := h.hub.AddClient(conn, userID, "")
	go client.WritePump()
	h.readPump(client, r.Context())
}

func (h *WsHandler) readPump(client *pkgws.Client, ctx context.Context) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512 * 1024)
	client.Conn.SetReadDeadline(time.Now().Add(pkgws.PongWait()))
	client.Conn.SetPongHandler(func(string) error {
		return client.Conn.SetReadDeadline(time.Now().Add(pkgws.PongWait()))
	})

	for {
		var envelope struct {
			Type    string          `json:"type"`
			Topic   string          `json:"topic,omitempty"`
			Room    string          `json:"room,omitempty"`
			Payload json.RawMessage `json:"payload,omitempty"`
		}
		err := client.Conn.ReadJSON(&envelope)
		if err != nil {
			break
		}
		pkgws.RecordReceivedMessage(envelope.Type)
		switch envelope.Type {
		case "subscribe":
			if envelope.Topic != "" {
				h.hub.Subscribe(client, envelope.Topic)
			}
		case "unsubscribe":
			if envelope.Topic != "" {
				h.hub.Unsubscribe(client, envelope.Topic)
			}
		case "join_room":
			if envelope.Room != "" {
				h.hub.JoinRoom(client, envelope.Room)
			}
		case "leave_room":
			if envelope.Room != "" {
				h.hub.LeaveRoom(client, envelope.Room)
			}
		case "message":
			if envelope.Topic == "" && envelope.Room == "" {
				// invalid
				continue
			}
			// ส่งไปยัง queue พร้อมบันทึก history
			if h.usecase != nil {
				if err := h.usecase.HandleIncomingMessage(ctx, envelope.Topic, envelope.Room, envelope.Payload, client.UserID); err != nil {
					h.log.Errorf("HandleIncomingMessage error: %v", err)
				}
			}
		default:
			// unknown
		}
	}
}
