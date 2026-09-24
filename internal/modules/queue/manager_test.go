package queue

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisQueue_PublishAndConsume(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	q := NewRedisQueue(client, 2)
	ctx := context.Background()

	var receivedMsg *Message
	handler := func(ctx context.Context, msg *Message) error {
		receivedMsg = msg
		return nil
	}

	err = q.Subscribe(ctx, "test", handler)
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	payload := map[string]string{"foo": "bar"}
	err = q.Publish(ctx, "test", payload)
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	// รอให้ consumer ประมวลผล
	time.Sleep(200 * time.Millisecond)

	if receivedMsg == nil {
		t.Fatal("expected received message")
	}
	if receivedMsg.Topic != "test" {
		t.Fatalf("expected topic %q, got %q", "test", receivedMsg.Topic)
	}

	var result map[string]string
	err = json.Unmarshal(receivedMsg.Payload, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if result["foo"] != "bar" {
		t.Fatalf("expected payload foo=bar, got %q", result["foo"])
	}
}
