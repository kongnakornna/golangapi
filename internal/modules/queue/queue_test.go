package queue

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisQueue_PublishAndConsume_SingleInstance(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	payload := map[string]string{"foo": "bar"}
	err = q.Publish(ctx, "test", payload)
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	assert.NotNil(t, receivedMsg)
	assert.Equal(t, "test", receivedMsg.Topic)
	var result map[string]string
	err = json.Unmarshal(receivedMsg.Payload, &result)
	assert.NoError(t, err)
	assert.Equal(t, "bar", result["foo"])

	q.Close()
}

func TestNoopQueue(t *testing.T) {
	q := NewNoop()
	ctx := context.Background()
	err := q.Publish(ctx, "test", []byte("hello"))
	assert.NoError(t, err)
	err = q.Subscribe(ctx, "test", func(ctx context.Context, msg *Message) error { return nil })
	assert.NoError(t, err)
	err = q.PublishDelayed(ctx, "test", []byte("hello"), time.Second)
	assert.NoError(t, err)
	err = q.Close()
	assert.NoError(t, err)
}
