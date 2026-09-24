package websocket

import (
	"context"
	"testing"
	"time"

	"icmongolang/internal/modules/queue"
	"icmongolang/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func newTestHub() *Hub {
	return NewHub(queue.NewNoop(), logger.NewLogger("websocket-test"))
}

func TestBroadcastAfterUnregister_NoPanic(t *testing.T) {
	hub := newTestHub()

	client := hub.AddClient(&websocket.Conn{}, "user1", "room1")
	hub.BroadcastMessage("x", "before")

	hub.Unregister(client)

	// These used to panic with "send on closed channel".
	hub.BroadcastMessage("after", "y")
	hub.BroadcastToRoom("room1", "after", "y")
	hub.BroadcastToTopic("topic1", "after", "y")
}

func TestBroadcastToTopic_DeliversToSubscriber(t *testing.T) {
	hub := newTestHub()

	client := hub.AddClient(&websocket.Conn{}, "user1", "")
	assert.NoError(t, hub.Subscribe(client, "BAACTW05/DATA"))

	hub.BroadcastToTopic("BAACTW05/DATA", "message", map[string]interface{}{"topic": "BAACTW05/DATA"})

	select {
	case msg := <-client.Send:
		assert.Contains(t, string(msg), `"event":"message"`)
	case <-time.After(time.Second):
		t.Fatal("expected message to be delivered to topic subscriber")
	}
}

func TestBroadcastToTopic_NotDeliveredToRoomOnlyClient(t *testing.T) {
	hub := newTestHub()

	roomClient := hub.AddClient(&websocket.Conn{}, "user1", "BAACTW05")
	hub.JoinRoom(roomClient, "BAACTW05")

	hub.BroadcastToTopic("BAACTW05/DATA", "message", "payload")

	select {
	case msg := <-roomClient.Send:
		t.Fatalf("room-only client should NOT receive topic broadcast, got %s", msg)
	case <-time.After(100 * time.Millisecond):
		// expected: nothing delivered
	}
}

func TestBroadcastToRoom_NotDeliveredToTopicSubscriber(t *testing.T) {
	hub := newTestHub()

	topicClient := hub.AddClient(&websocket.Conn{}, "user1", "")
	assert.NoError(t, hub.Subscribe(topicClient, "BAACTW05/DATA"))

	hub.BroadcastToRoom("BAACTW05", "mqtt", map[string]interface{}{"topic": "BAACTW05/DATA"})

	select {
	case msg := <-topicClient.Send:
		t.Fatalf("topic subscriber should NOT receive room broadcast, got %s", msg)
	case <-time.After(100 * time.Millisecond):
		// expected: nothing delivered
	}
}

func TestBroadcastToRoom_DeliversToRoomMember(t *testing.T) {
	hub := newTestHub()

	client := hub.AddClient(&websocket.Conn{}, "user1", "BAACTW05")
	hub.BroadcastToRoom("BAACTW05", "mqtt", map[string]interface{}{"topic": "BAACTW05/DATA"})

	select {
	case msg := <-client.Send:
		assert.Contains(t, string(msg), `"event":"mqtt"`)
	case <-time.After(time.Second):
		t.Fatal("expected message to be delivered to room member")
	}
}

func TestPublishMessageToQueueAndReceive(t *testing.T) {
	hub := newTestHub()

	client := hub.AddClient(&websocket.Conn{}, "user1", "")
	assert.NoError(t, hub.Subscribe(client, "topic1"))

	// Use a fake conn-less check: PublishMessage goes to the (noop) queue and
	// returns without error; delivery to subscribers is covered by
	// TestBroadcastToTopic_DeliversToSubscriber.
	ctx := context.Background()
	assert.NoError(t, hub.PublishMessage(ctx, "topic1", []byte(`{"a":1}`), "user1"))
}
