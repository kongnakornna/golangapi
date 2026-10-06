package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"icmongolang/internal/modules/websocket/models"
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
// @Summary      Get WebSocket message history (legacy)
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

// GetMessagesPaginated godoc
// @Summary      Get WebSocket message history with pagination
// @Description  Returns paginated message history with cursor-based and offset-based pagination support.
// @Tags         websocket
// @Produce      json
// @Param        topic     query string false "Topic to filter messages"
// @Param        room      query string false "Room to filter messages"
// @Param        sender_id query string false "Sender ID to filter messages"
// @Param        before    query string false "RFC3339 timestamp - messages before this time"
// @Param        after     query string false "RFC3339 timestamp - messages after this time"
// @Param        limit     query int    false "Max messages per page (default 50, max 1000)"
// @Param        offset    query int    false "Offset for offset-based pagination"
// @Param        order     query string false "Sort order: asc or desc (default desc)"
// @Success      200 {object} presenter.MessagePageResponse
// @Failure      400 {string} string "invalid query parameters"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/messages/paginated [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetMessagesPaginated(w http.ResponseWriter, r *http.Request) {
	params, err := h.parseMessageQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.usecase.GetMessages(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := presenter.MessagePageResponse{
		Messages:   make([]presenter.WSMessageItem, len(page.Messages)),
		TotalCount: page.TotalCount,
		HasMore:    page.HasMore,
	}
	for i, msg := range page.Messages {
		response.Messages[i] = presenter.WSMessageItem{
			ID:         msg.ID,
			Topic:      msg.Topic,
			Room:       msg.Room,
			Payload:    msg.Payload,
			SenderID:   msg.SenderID,
			SenderRole: msg.SenderRole,
			SentAt:     msg.SentAt.Format(time.RFC3339),
		}
	}
	if page.NextCursor != nil {
		response.NextCursor = page.NextCursor.Format(time.RFC3339)
	}
	if page.PrevCursor != nil {
		response.PrevCursor = page.PrevCursor.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchMessages godoc
// @Summary      Search messages with full-text search
// @Description  Performs full-text search on persisted WebSocket messages with optional filters and highlighting.
// @Tags         websocket
// @Produce      json
// @Param        query     query string true "Search query string"
// @Param        topic     query string false "Filter by topic"
// @Param        room      query string false "Filter by room"
// @Param        sender_id query string false "Filter by sender ID"
// @Param        before    query string false "RFC3339 timestamp - messages before this time"
// @Param        after     query string false "RFC3339 timestamp - messages after this time"
// @Param        limit     query int    false "Max results per page (default 50, max 1000)"
// @Param        offset    query int    false "Offset for pagination"
// @Param        order     query string false "Sort order: asc or desc (default desc by rank)"
// @Param        highlight query bool   false "Include highlighted snippets"
// @Success      200 {object} presenter.SearchPageResponse
// @Failure      400 {string} string "invalid query parameters"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/messages/search [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	params, err := h.parseSearchQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.usecase.SearchMessages(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := presenter.SearchPageResponse{
		Results:    make([]presenter.SearchResultItem, len(page.Results)),
		TotalCount: page.TotalCount,
		HasMore:    page.HasMore,
	}
	for i, result := range page.Results {
		item := presenter.SearchResultItem{
			Message: presenter.WSMessageItem{
				ID:         result.Message.ID,
				Topic:      result.Message.Topic,
				Room:       result.Message.Room,
				Payload:    result.Message.Payload,
				SenderID:   result.Message.SenderID,
				SenderRole: result.Message.SenderRole,
				SentAt:     result.Message.SentAt.Format(time.RFC3339),
			},
			Rank: result.Rank,
		}
		if result.Highlights != nil {
			item.Highlights = result.Highlights
		}
		response.Results[i] = item
	}
	if page.NextCursor != nil {
		response.NextCursor = page.NextCursor.Format(time.RFC3339)
	}
	if page.PrevCursor != nil {
		response.PrevCursor = page.PrevCursor.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseSearchQueryParams parses query parameters for message search.
func (h *WsRestHandler) parseSearchQueryParams(r *http.Request) (models.SearchQueryParams, error) {
	params := models.SearchQueryParams{
		Query:     r.URL.Query().Get("query"),
		Topic:     r.URL.Query().Get("topic"),
		Room:      r.URL.Query().Get("room"),
		SenderID:  r.URL.Query().Get("sender_id"),
		Order:     r.URL.Query().Get("order"),
		Highlight: r.URL.Query().Get("highlight") == "true",
	}

	if params.Query == "" {
		return params, fmt.Errorf("query parameter is required")
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return params, fmt.Errorf("invalid limit: %w", err)
		}
		params.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return params, fmt.Errorf("invalid offset: %w", err)
		}
		params.Offset = offset
	}

	if beforeStr := r.URL.Query().Get("before"); beforeStr != "" {
		t, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return params, fmt.Errorf("invalid before timestamp: %w", err)
		}
		params.Before = &t
	}

	if afterStr := r.URL.Query().Get("after"); afterStr != "" {
		t, err := time.Parse(time.RFC3339, afterStr)
		if err != nil {
			return params, fmt.Errorf("invalid after timestamp: %w", err)
		}
		params.After = &t
	}

	if orderStr := r.URL.Query().Get("order"); orderStr != "" {
		if orderStr != "asc" && orderStr != "desc" {
			return params, fmt.Errorf("invalid order: must be asc or desc")
		}
		params.Order = orderStr
	}

	if highlightStr := r.URL.Query().Get("highlight"); highlightStr != "" {
		params.Highlight = highlightStr == "true"
	}

	return params, nil
}

// GetMessagesByRoom godoc
// @Summary      Get WebSocket message history for a room with pagination
// @Description  Returns paginated message history for a specific room.
// @Tags         websocket
// @Produce      json
// @Param        room      path string true "Room name"
// @Param        before    query string false "RFC3339 timestamp - messages before this time"
// @Param        after     query string false "RFC3339 timestamp - messages after this time"
// @Param        limit     query int    false "Max messages per page (default 50, max 1000)"
// @Param        offset    query int    false "Offset for offset-based pagination"
// @Param        order     query string false "Sort order: asc or desc (default desc)"
// @Success      200 {object} presenter.MessagePageResponse
// @Failure      400 {string} string "invalid query parameters"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/rooms/{room}/messages [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetMessagesByRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room := vars["room"]
	if room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}

	params, err := h.parseMessageQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.usecase.GetMessagesByRoom(r.Context(), room, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := presenter.MessagePageResponse{
		Messages:   make([]presenter.WSMessageItem, len(page.Messages)),
		TotalCount: page.TotalCount,
		HasMore:    page.HasMore,
	}
	for i, msg := range page.Messages {
		response.Messages[i] = presenter.WSMessageItem{
			ID:         msg.ID,
			Topic:      msg.Topic,
			Room:       msg.Room,
			Payload:    msg.Payload,
			SenderID:   msg.SenderID,
			SenderRole: msg.SenderRole,
			SentAt:     msg.SentAt.Format(time.RFC3339),
		}
	}
	if page.NextCursor != nil {
		response.NextCursor = page.NextCursor.Format(time.RFC3339)
	}
	if page.PrevCursor != nil {
		response.PrevCursor = page.PrevCursor.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMessagesBySender godoc
// @Summary      Get WebSocket message history for a sender with pagination
// @Description  Returns paginated message history for a specific sender.
// @Tags         websocket
// @Produce      json
// @Param        sender_id path string true "Sender ID"
// @Param        before    query string false "RFC3339 timestamp - messages before this time"
// @Param        after     query string false "RFC3339 timestamp - messages after this time"
// @Param        limit     query int    false "Max messages per page (default 50, max 1000)"
// @Param        offset    query int    false "Offset for offset-based pagination"
// @Param        order     query string false "Sort order: asc or desc (default desc)"
// @Success      200 {object} presenter.MessagePageResponse
// @Failure      400 {string} string "invalid query parameters"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/senders/{sender_id}/messages [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetMessagesBySender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderID := vars["sender_id"]
	if senderID == "" {
		http.Error(w, "sender_id required", http.StatusBadRequest)
		return
	}

	params, err := h.parseMessageQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := h.usecase.GetMessagesBySender(r.Context(), senderID, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := presenter.MessagePageResponse{
		Messages:   make([]presenter.WSMessageItem, len(page.Messages)),
		TotalCount: page.TotalCount,
		HasMore:    page.HasMore,
	}
	for i, msg := range page.Messages {
		response.Messages[i] = presenter.WSMessageItem{
			ID:         msg.ID,
			Topic:      msg.Topic,
			Room:       msg.Room,
			Payload:    msg.Payload,
			SenderID:   msg.SenderID,
			SenderRole: msg.SenderRole,
			SentAt:     msg.SentAt.Format(time.RFC3339),
		}
	}
	if page.NextCursor != nil {
		response.NextCursor = page.NextCursor.Format(time.RFC3339)
	}
	if page.PrevCursor != nil {
		response.PrevCursor = page.PrevCursor.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

// parseMessageQueryParams parses query parameters for message pagination.
func (h *WsRestHandler) parseMessageQueryParams(r *http.Request) (models.MessageQueryParams, error) {
	params := models.MessageQueryParams{
		Topic:    r.URL.Query().Get("topic"),
		Room:     r.URL.Query().Get("room"),
		SenderID: r.URL.Query().Get("sender_id"),
		Order:    r.URL.Query().Get("order"),
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return params, fmt.Errorf("invalid limit: %w", err)
		}
		params.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return params, fmt.Errorf("invalid offset: %w", err)
		}
		params.Offset = offset
	}

	if beforeStr := r.URL.Query().Get("before"); beforeStr != "" {
		t, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return params, fmt.Errorf("invalid before timestamp: %w", err)
		}
		params.Before = &t
	}

	if afterStr := r.URL.Query().Get("after"); afterStr != "" {
		t, err := time.Parse(time.RFC3339, afterStr)
		if err != nil {
			return params, fmt.Errorf("invalid after timestamp: %w", err)
		}
		params.After = &t
	}

	return params, nil
}

// GetPresence godoc
// @Summary      Get user presence status
// @Description  Returns the presence status of a specific user.
// @Tags         websocket
// @Produce      json
// @Param        user_id path string true "User ID"
// @Success      200 {object} presenter.PresenceItem
// @Failure      404 {string} string "user not found"
// @Router       /api/ws/presence/{user_id} [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetPresence(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	info, exists := h.hub.GetUserPresence(userID)
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := presenter.PresenceItem{
		UserID:    info.UserID,
		Status:    string(info.Status),
		Meta:      info.Meta,
		LastSeen:  info.LastSeen.Format(time.RFC3339),
		Instances: info.Instances,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllPresence godoc
// @Summary      Get all user presence statuses
// @Description  Returns the presence status of all users with pagination.
// @Tags         websocket
// @Produce      json
// @Param        user_ids query string false "Comma-separated user IDs to filter"
// @Param        status    query string false "Filter by status (online, away, busy, offline)"
// @Param        limit     query int    false "Max results per page (default 50)"
// @Param        offset    query int    false "Offset for pagination"
// @Success      200 {object} presenter.PresencePageResponse
// @Failure      400 {string} string "invalid query parameters"
// @Router       /api/ws/presence [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetAllPresence(w http.ResponseWriter, r *http.Request) {
	params := websocket.PresenceQueryParams{
		Status: websocket.PresenceStatus(r.URL.Query().Get("status")),
	}

	if userIDsStr := r.URL.Query().Get("user_ids"); userIDsStr != "" {
		params.UserIDs = strings.Split(userIDsStr, ",")
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		params.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		params.Offset = offset
	}

	var presences []*websocket.PresenceInfo
	if len(params.UserIDs) > 0 {
		for _, userID := range params.UserIDs {
			if info, exists := h.hub.GetUserPresence(userID); exists {
				presences = append(presences, info)
			}
		}
	} else if params.Status != "" {
		presences = h.hub.GetPresenceByStatus(websocket.PresenceStatus(params.Status))
	} else {
		allPresence := h.hub.GetAllPresence()
		for _, info := range allPresence {
			presences = append(presences, info)
		}
	}

	// Apply pagination
	start := params.Offset
	end := start + params.Limit
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if start > len(presences) {
		presences = []*websocket.PresenceInfo{}
	} else {
		if end > len(presences) {
			end = len(presences)
		}
		presences = presences[start:end]
	}

	response := presenter.PresencePageResponse{
		Presences: make([]presenter.PresenceItem, len(presences)),
		Total:     int64(len(presences)),
		HasMore:   false,
	}
	for i, info := range presences {
		response.Presences[i] = presenter.PresenceItem{
			UserID:    info.UserID,
			Status:    string(info.Status),
			Meta:      info.Meta,
			LastSeen:  info.LastSeen.Format(time.RFC3339),
			Instances: info.Instances,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateRetentionPolicy godoc
// @Summary      Create a new retention policy
// @Description  Creates a new message retention policy for topics/rooms.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        request body presenter.RetentionPolicyRequest true "Retention policy to create"
// @Success      201 {object} presenter.RetentionPolicyResponse
// @Failure      400 {string} string "invalid request"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/policies [post]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) CreateRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	var req presenter.RetentionPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Parse duration
	maxAge, err := time.ParseDuration(req.MaxAge)
	if err != nil {
		http.Error(w, "invalid max_age format: "+err.Error(), http.StatusBadRequest)
		return
	}

	policyReq := &models.RetentionPolicyRequest{
		Topic:    req.Topic,
		Room:     req.Room,
		MaxAge:   maxAge,
		MaxCount: req.MaxCount,
		Enabled:  req.Enabled,
	}

	response, err := h.usecase.CreateRetentionPolicy(r.Context(), policyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetRetentionPolicy godoc
// @Summary      Get a retention policy by ID
// @Description  Retrieves a specific retention policy.
// @Tags         websocket
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} presenter.RetentionPolicyResponse
// @Failure      404 {string} string "policy not found"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/policies/{id} [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	response, err := h.usecase.GetRetentionPolicy(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListRetentionPolicies godoc
// @Summary      List all retention policies
// @Description  Returns all retention policies with pagination.
// @Tags         websocket
// @Produce      json
// @Param        limit  query int false "Max results per page (default 50, max 1000)"
// @Param        offset query int false "Offset for pagination"
// @Success      200 {object} presenter.RetentionPolicyListResponse
// @Failure      400 {string} string "invalid query parameters"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/policies [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) ListRetentionPolicies(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	response, err := h.usecase.ListRetentionPolicies(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateRetentionPolicy godoc
// @Summary      Update a retention policy
// @Description  Updates an existing retention policy.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        id      path string true "Policy ID"
// @Param        request body presenter.RetentionPolicyRequest true "Retention policy to update"
// @Success      200 {object} presenter.RetentionPolicyResponse
// @Failure      400 {string} string "invalid request"
// @Failure      404 {string} string "policy not found"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/policies/{id} [put]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) UpdateRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	var req presenter.RetentionPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	maxAge, err := time.ParseDuration(req.MaxAge)
	if err != nil {
		http.Error(w, "invalid max_age format: "+err.Error(), http.StatusBadRequest)
		return
	}

	policyReq := &models.RetentionPolicyRequest{
		Topic:    req.Topic,
		Room:     req.Room,
		MaxAge:   maxAge,
		MaxCount: req.MaxCount,
		Enabled:  req.Enabled,
	}

	response, err := h.usecase.UpdateRetentionPolicy(r.Context(), id, policyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteRetentionPolicy godoc
// @Summary      Delete a retention policy
// @Description  Deletes a retention policy by ID.
// @Tags         websocket
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      204 {string} string "deleted"
// @Failure      404 {string} string "policy not found"
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/policies/{id} [delete]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) DeleteRetentionPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	err := h.usecase.DeleteRetentionPolicy(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRetentionStats godoc
// @Summary      Get retention statistics
// @Description  Returns statistics about message retention.
// @Tags         websocket
// @Produce      json
// @Success      200 {object} presenter.RetentionStatsResponse
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/stats [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetRetentionStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.usecase.GetRetentionStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := presenter.RetentionStatsResponse{
		TotalMessages:   stats.TotalMessages,
		MessagesDeleted: stats.MessagesDeleted,
		BytesFreed:      stats.BytesFreed,
		PoliciesCount:   stats.PoliciesCount,
		ActivePolicies:  stats.ActivePolicies,
	}
	if !stats.LastCleanupAt.IsZero() {
		response.LastCleanupAt = stats.LastCleanupAt.Format(time.RFC3339)
	}
	if !stats.NextCleanupAt.IsZero() {
		response.NextCleanupAt = stats.NextCleanupAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RunRetentionCleanup godoc
// @Summary      Run retention cleanup
// @Description  Manually triggers retention cleanup for all active policies.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Param        request body presenter.CleanupRequest true "Cleanup request"
// @Success      200 {object} presenter.CleanupResponse
// @Failure      500 {string} string "internal server error"
// @Router       /api/ws/retention/cleanup [post]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) RunRetentionCleanup(w http.ResponseWriter, r *http.Request) {
	var req presenter.CleanupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body for running all policies
		req = presenter.CleanupRequest{}
	}

	start := time.Now()
	stats, err := h.usecase.RunRetentionCleanup(r.Context())
	duration := time.Since(start)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := presenter.CleanupResponse{
		Status:          "completed",
		MessagesDeleted: stats.MessagesDeleted,
		BytesFreed:      stats.BytesFreed,
		Duration:        duration.String(),
		PoliciesRun:     stats.ActivePolicies,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseRetentionPolicyRequest parses a retention policy request from JSON.
func (h *WsRestHandler) parseRetentionPolicyRequest(r *http.Request) (*models.RetentionPolicyRequest, error) {
	var req presenter.RetentionPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	maxAge, err := time.ParseDuration(req.MaxAge)
	if err != nil {
		return nil, fmt.Errorf("invalid max_age format: %w", err)
	}

	return &models.RetentionPolicyRequest{
		Topic:    req.Topic,
		Room:     req.Room,
		MaxAge:   maxAge,
		MaxCount: req.MaxCount,
		Enabled:  req.Enabled,
	}, nil
}

// GetRoomPresence godoc
// @Summary      Get users currently in a room
// @Description  Returns the presence status of all users currently in a specific room.
// @Tags         websocket
// @Produce      json
// @Param        room path string true "Room name"
// @Success      200 {array} presenter.PresenceItem
// @Failure      400 {string} string "room required"
// @Router       /api/ws/rooms/{room}/presence [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetRoomPresence(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room := vars["room"]
	if room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}

	presences := h.hub.GetRoomPresence(room)

	response := make([]presenter.PresenceItem, len(presences))
	for i, info := range presences {
		response[i] = presenter.PresenceItem{
			UserID:    info.UserID,
			Status:    string(info.Status),
			Meta:      info.Meta,
			LastSeen:  info.LastSeen.Format(time.RFC3339),
			Instances: info.Instances,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserCurrentRoom godoc
// @Summary      Get user's current room
// @Description  Returns the current room for a specific user (if they have an active connection).
// @Tags         websocket
// @Produce      json
// @Param        user_id path string true "User ID"
// @Success      200 {object} map[string]string
// @Failure      404 {string} string "user not found or not in a room"
// @Router       /api/ws/users/{user_id}/room [get]
// @Security     OAuth2Password
// @Security     BearerAuth
func (h *WsRestHandler) GetUserCurrentRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	room := h.hub.GetUserCurrentRoom(userID)
	if room == "" {
		http.Error(w, "user not found or not in a room", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID, "room": room})
}
