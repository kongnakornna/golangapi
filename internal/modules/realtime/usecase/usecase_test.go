package usecase_test

import (
	"context"
	"testing"

	"icmongolang/internal/modules/realtime/presenter"
	"icmongolang/internal/modules/realtime/usecase"
	"icmongolang/pkg/logger"
)

type stubBroadcaster struct {
	rooms          map[string]bool
	broadcastedRoom string
	broadcastedGlobal bool
	lastEvent      string
	lastData       map[string]any
}

func newStub() *stubBroadcaster {
	return &stubBroadcaster{rooms: map[string]bool{"roomA": true}}
}

func (s *stubBroadcaster) BroadcastToRoom(room, event string, data interface{}) {
	s.broadcastedRoom = room
	s.lastEvent = event
	s.lastData, _ = data.(map[string]any)
}
func (s *stubBroadcaster) BroadcastToTopic(topic, event string, data interface{}) {
	s.broadcastedRoom = topic
	s.lastEvent = event
	s.lastData, _ = data.(map[string]any)
}
func (s *stubBroadcaster) BroadcastMessage(event string, data interface{}) {
	s.broadcastedGlobal = true
	s.lastEvent = event
	s.lastData, _ = data.(map[string]any)
}
func (s *stubBroadcaster) GetRooms() []string { return []string{"roomA", "roomB"} }
func (s *stubBroadcaster) GetClientsInRoom(room string) int {
	if _, ok := s.rooms[room]; ok {
		return 2
	}
	return 0
}

func TestPublish_ToRoom(t *testing.T) {
	hub := newStub()
	uc := usecase.NewRealtimeUseCase(hub, logger.NewLogger("realtime-test"))
	req := &presenter.PublishRequest{Room: "roomA", Event: "monitor", Data: map[string]any{"value": 25.5}}

	resp, err := uc.Publish(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected OK=true")
	}
	if hub.broadcastedRoom != "roomA" {
		t.Fatalf("expected roomA, got %q", hub.broadcastedRoom)
	}
	if hub.lastEvent != "monitor" {
		t.Fatalf("expected event monitor, got %q", hub.lastEvent)
	}
	if resp.Clients != 2 {
		t.Fatalf("expected 2 clients, got %d", resp.Clients)
	}
}

func TestPublish_Global(t *testing.T) {
	hub := newStub()
	uc := usecase.NewRealtimeUseCase(hub, logger.NewLogger("realtime-test"))
	req := &presenter.PublishRequest{Event: "alarm", Data: map[string]any{"id": "x"}}

	resp, err := uc.Publish(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hub.broadcastedGlobal {
		t.Fatalf("expected global broadcast")
	}
	if resp.Clients != 2 {
		t.Fatalf("expected rooms=2 clients, got %d", resp.Clients)
	}
}

func TestPublish_NilBroadcaster_DoesNotPanic(t *testing.T) {
	uc := usecase.NewRealtimeUseCase(nil, logger.NewLogger("realtime-test"))
	req := &presenter.PublishRequest{Event: "monitor", Data: map[string]any{}}

	resp, err := uc.Publish(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OK {
		t.Fatalf("expected OK=false when broadcaster nil")
	}
}

func TestPublish_InvalidRequest(t *testing.T) {
	uc := usecase.NewRealtimeUseCase(nil, logger.NewLogger("realtime-test"))
	if _, err := uc.Publish(context.Background(), nil); err == nil {
		t.Fatalf("expected error for nil request")
	}
	empty := &presenter.PublishRequest{Data: map[string]any{}}
	if _, err := uc.Publish(context.Background(), empty); err == nil {
		t.Fatalf("expected error for empty event")
	}
}
