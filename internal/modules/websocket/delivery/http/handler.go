package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"icmongolang/internal/modules/websocket/presenter"
	"icmongolang/internal/modules/websocket/usecase"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"

	"github.com/gorilla/mux"
)

type WsRestHandler struct {
	hub     *websocket.Hub
	usecase usecase.WSUsecase
	logger  logger.Logger
}

func NewWsRestHandler(hub *websocket.Hub, uc usecase.WSUsecase, log logger.Logger) *WsRestHandler {
	return &WsRestHandler{hub: hub, usecase: uc, logger: log}
}

// SendMessage godoc
// @Summary      Send a message to a WebSocket room
// @Description  Broadcasts an event/data payload to all WebSocket clients currently joined to the room.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        request body presenter.SendMessageRequest true "Message to broadcast"
// @Success      200 {object} presenter.SendMessageResponse
// @Failure      400 {string} string "room and event are required"
// @Router       /api/ws/messages [post]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req presenter.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Room == "" || req.Event == "" {
		http.Error(w, "room and event are required", http.StatusBadRequest)
		return
	}
	h.hub.BroadcastToRoom(req.Room, req.Event, req.Data)
	go h.usecase.SaveMessage(r.Context(), req.Room, []byte{}, "api") // หรือจะบันทึกเพิ่มเติม
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presenter.SendMessageResponse{Status: "sent"})
}

// GetMessages godoc
// @Summary      Get WebSocket message history
// @Description  Returns the last N messages sent on a topic/room, persisted in PostgreSQL.
// @Tags         websocket
// @Produce      json
// @Param        room   query string false "Room/topic to fetch history for"
// @Param        limit  query int    false "Max messages (default 20)"
// @Success      200 {array} presenter.WSMessageItem
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/messages [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	messages, err := h.usecase.GetTopicHistory(r.Context(), room, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetRooms godoc
// @Summary      List active WebSocket rooms
// @Description  Returns all rooms that currently have at least one connected WebSocket client.
// @Tags         websocket
// @Produce      json
// @Success      200 {object} presenter.RoomsResponse
// @Router       /api/ws/rooms [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := h.hub.GetRooms()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presenter.RoomsResponse{Rooms: rooms, Count: len(rooms)})
}

// GetRoomStats godoc
// @Summary      Get WebSocket room stats
// @Description  Returns the number of connected WebSocket clients in a room.
// @Tags         websocket
// @Produce      json
// @Param        room path string true "Room name"
// @Success      200 {object} presenter.RoomStatsResponse
// @Failure      400 {string} string "room required"
// @Router       /api/ws/rooms/{room}/stats [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetRoomStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room := vars["room"]
	if room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}
	count := h.hub.GetClientsInRoom(room)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presenter.RoomStatsResponse{Room: room, Clients: count})
}
