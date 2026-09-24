package usecase

import (
	"context"
	"errors"

	"icmongolang/internal/modules/realtime/presenter"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"
)

// ErrInvalidRequest is returned when a publish request is malformed.
var ErrInvalidRequest = errors.New("realtime: invalid publish request")

// RealtimeUseCase defines the contract for pushing real-time events to the
// dashboard websocket hub.
type RealtimeUseCase interface {
	Publish(ctx context.Context, req *presenter.PublishRequest) (*presenter.PublishResponse, error)
	Health(ctx context.Context) *presenter.HealthResponse
}

type realtimeUseCase struct {
	broadcaster websocket.Broadcaster
	logger      logger.Logger
}

// NewRealtimeUseCase builds the realtime use case. broadcaster may be nil; in
// that case Publish returns a disabled response without panicking.
func NewRealtimeUseCase(broadcaster websocket.Broadcaster, log logger.Logger) RealtimeUseCase {
	return &realtimeUseCase{broadcaster: broadcaster, logger: log}
}

func (u *realtimeUseCase) Publish(ctx context.Context, req *presenter.PublishRequest) (*presenter.PublishResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	if req.Event == "" {
		return nil, ErrInvalidRequest
	}
	if u.broadcaster == nil {
		// Websocket hub not wired: report disabled, never block existing flows.
		u.logger.Warn("realtime: broadcaster is nil – publish skipped")
		return &presenter.PublishResponse{OK: false, Event: req.Event, Clients: 0}, nil
	}

	var clients int
	if req.Room != "" {
		u.broadcaster.BroadcastToRoom(req.Room, req.Event, req.Data)
		clients = u.broadcaster.GetClientsInRoom(req.Room)
	} else {
		u.broadcaster.BroadcastMessage(req.Event, req.Data)
		clients = len(u.broadcaster.GetRooms())
	}
	u.logger.Infof("realtime: broadcast event=%s room=%q clients=%d", req.Event, req.Room, clients)
	return &presenter.PublishResponse{OK: true, Room: req.Room, Event: req.Event, Clients: clients}, nil
}

func (u *realtimeUseCase) Health(ctx context.Context) *presenter.HealthResponse {
	return &presenter.HealthResponse{Enabled: u.broadcaster != nil}
}
