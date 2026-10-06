package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/websocket/usecase"
	"icmongolang/pkg/logger"
	pkgws "icmongolang/pkg/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Origin validation is done by WSOriginGuard middleware
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WsHandler struct {
	hub     *pkgws.Hub
	usecase usecase.WSUsecase
	log     logger.Logger
	mw      *middleware.MiddlewareManager
}

func NewWsHandler(hub *pkgws.Hub, uc usecase.WSUsecase, log logger.Logger, mw *middleware.MiddlewareManager) *WsHandler {
	return &WsHandler{hub: hub, usecase: uc, log: log, mw: mw}
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
	// Extract token from query param or Authorization header
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
	}

	// Use JWE middleware for authentication
	if h.mw != nil && h.mw.JWEKeyPair() != nil {
		claims, err := h.mw.AuthenticateWS(r.Context(), token)
		if err != nil {
			h.log.Errorf("WS auth failed: %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			h.log.Errorf("Upgrade error: %v", err)
			return
		}

		client := h.hub.AddClient(conn, claims.UserID, "")
		// Store role for broadcast targeting
		client.Role = string(claims.Role)
		go client.WritePump()
		h.readPump(client, r.Context())
		return
	}

	// Fallback to legacy authentication
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

// HandlerFunc returns an http.HandlerFunc for use with chi router
func (h *WsHandler) HandlerFunc() http.HandlerFunc {
	return h.ServeHTTP
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
		var env pkgws.Envelope
		err := client.Conn.ReadJSON(&env)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Errorf("WS read error: %v", err)
			}
			break
		}

		pkgws.RecordReceivedMessage(env.Type)

		// Handle client-to-server messages
		switch env.Type {
		case pkgws.MsgTypeSubscribe:
			var payload pkgws.SubscribePayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil && payload.Topic != "" {
				h.hub.Subscribe(client, payload.Topic)
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeUnsubscribe:
			var payload pkgws.UnsubscribePayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil && payload.Topic != "" {
				h.hub.Unsubscribe(client, payload.Topic)
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeJoinRoom:
			var payload pkgws.JoinRoomPayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil && payload.Room != "" {
				h.hub.JoinRoom(client, payload.Room)
				h.sendEvent(client, pkgws.BuildJoinedEvent(payload.Room, "", 0))
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeLeaveRoom:
			var payload pkgws.LeaveRoomPayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil && payload.Room != "" {
				h.hub.LeaveRoom(client, payload.Room)
				h.sendEvent(client, pkgws.BuildLeftEvent(payload.Room, "", 0))
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeMessage:
			var payload pkgws.MessagePayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil {
				// Call OnMessage callback - return false to prevent broadcast
				allow := true
				if h.hub.Callbacks.OnMessage != nil {
					allow = h.hub.Callbacks.OnMessage(client, payload.Topic, payload.Room, payload.Content)
				}
				if allow {
					if payload.Topic != "" {
						h.hub.BroadcastToTopic(payload.Topic, pkgws.EventTypeMessage, pkgws.MessageEventPayload{
							ID:         pkgws.GenerateID(),
							Topic:      payload.Topic,
							Room:       payload.Room,
							SenderID:   client.UserID,
							SenderRole: client.Role,
							Content:    payload.Content,
							Meta:       payload.Meta,
							Timestamp:  time.Now().UnixMilli(),
						})
					} else if payload.Room != "" {
						h.hub.BroadcastToRoom(payload.Room, pkgws.EventTypeMessage, pkgws.MessageEventPayload{
							ID:         pkgws.GenerateID(),
							Topic:      payload.Topic,
							Room:       payload.Room,
							SenderID:   client.UserID,
							SenderRole: client.Role,
							Content:    payload.Content,
							Meta:       payload.Meta,
							Timestamp:  time.Now().UnixMilli(),
						})
					}
				}
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypePing:
			h.sendEvent(client, pkgws.BuildPongEvent())

		case pkgws.MsgTypeAuth:
			var payload pkgws.AuthPayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil {
				// Re-authenticate with new token
				if h.mw != nil && h.mw.JWEKeyPair() != nil {
					claims, err := h.mw.AuthenticateWS(ctx, payload.Token)
					if err != nil {
						h.sendEvent(client, pkgws.BuildAuthEvent(false, "", "", err.Error()))
					} else {
						client.UserID = claims.UserID
						client.Role = string(claims.Role)
						h.sendEvent(client, pkgws.BuildAuthEvent(true, claims.UserID, string(claims.Role), ""))
					}
				}
			}

		case pkgws.MsgTypePresence:
			var payload pkgws.PresencePayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil {
				h.hub.BroadcastToRoom("", pkgws.EventTypePresence, pkgws.PresenceEventPayload{
					UserID: client.UserID,
					Status: payload.Status,
					Meta:   payload.Meta,
				})
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeTyping:
			var payload pkgws.TypingPayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil {
				if payload.Room != "" {
					h.hub.BroadcastToRoom(payload.Room, pkgws.EventTypeTyping, pkgws.TypingEventPayload{
						UserID:   client.UserID,
						Room:     payload.Room,
						Topic:    payload.Topic,
						IsTyping: payload.IsTyping,
					})
				} else if payload.Topic != "" {
					h.hub.BroadcastToTopic(payload.Topic, pkgws.EventTypeTyping, pkgws.TypingEventPayload{
						UserID:   client.UserID,
						Room:     payload.Room,
						Topic:    payload.Topic,
						IsTyping: payload.IsTyping,
					})
				}
			}

		case pkgws.MsgTypeReadReceipt:
			var payload pkgws.ReadReceiptPayload
			if err := pkgws.UnmarshalPayload(env.Payload, &payload); err == nil {
				if payload.Room != "" {
					h.hub.BroadcastToRoom(payload.Room, pkgws.EventTypeReadReceipt, pkgws.ReadReceiptEventPayload{
						UserID:     client.UserID,
						MessageIDs: payload.MessageIDs,
						Room:       payload.Room,
						Topic:      payload.Topic,
					})
				} else if payload.Topic != "" {
					h.hub.BroadcastToTopic(payload.Topic, pkgws.EventTypeReadReceipt, pkgws.ReadReceiptEventPayload{
						UserID:     client.UserID,
						MessageIDs: payload.MessageIDs,
						Room:       payload.Room,
						Topic:      payload.Topic,
					})
				}
				h.sendAck(client, env.ID, "ok", "")
			}

		case pkgws.MsgTypeAck:
			// Handle incoming acknowledgment from client
			if client.AckManager != nil {
				client.AckManager.HandleAck(&env)
			}
			// Also handle sequence acknowledgment
			if env.SeqAck > 0 {
				client.HandleSeqAck(env.SeqAck)
			}

		default:
			h.log.Warnf("Unknown message type: %s", env.Type)
			h.sendEvent(client, pkgws.BuildErrorEvent("unknown_type", "unknown message type: "+env.Type, nil))
		}
	}
}

// sendEvent sends an event envelope to a specific client.
func (h *WsHandler) sendEvent(client *pkgws.Client, env *pkgws.Envelope) {
	// Include the client's highest acknowledged sequence number
	env.SeqAck = client.GetAckedSeq()
	data, err := json.Marshal(env)
	if err != nil {
		h.log.Errorf("Failed to marshal event: %v", err)
		return
	}
	client.SendMessage(data)
}

// sendAck sends an acknowledgment to the client.
func (h *WsHandler) sendAck(client *pkgws.Client, originalID, status, errMsg string) {
	if originalID == "" {
		return
	}
	env := pkgws.BuildAckEvent(originalID, status, errMsg)
	// Include the client's highest acknowledged sequence number
	env.SeqAck = client.GetAckedSeq()
	data, err := json.Marshal(env)
	if err != nil {
		h.log.Errorf("Failed to marshal ack: %v", err)
		return
	}
	client.SendMessage(data)
}
