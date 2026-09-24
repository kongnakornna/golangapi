package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestIdempotencyStore(t *testing.T) (*redis.Client, *IdempotencyStore, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	store := NewIdempotencyStore(client)

	return client, store, mr
}

func TestIdempotencyStore_MarkAndCheckProcessed(t *testing.T) {
	_, store, _ := newTestIdempotencyStore(t)

	ctx := context.Background()
	const (
		eventID = "evt-001"
		handler = "consent-event"
	)

	processed, err := store.IsProcessed(ctx, eventID, handler)
	require.NoError(t, err)
	assert.False(t, processed)

	marked, err := store.MarkProcessed(ctx, eventID, handler, time.Minute)
	require.NoError(t, err)
	assert.True(t, marked)

	markedAgain, err := store.MarkProcessed(ctx, eventID, handler, time.Minute)
	require.NoError(t, err)
	assert.False(t, markedAgain)

	processed, err = store.IsProcessed(ctx, eventID, handler)
	require.NoError(t, err)
	assert.True(t, processed)
}

func TestIdempotencyStore_NotProcessedBeforeMark(t *testing.T) {
	_, store, _ := newTestIdempotencyStore(t)

	ctx := context.Background()

	processed, err := store.IsProcessed(ctx, "evt-ghost", "consent-event")
	require.NoError(t, err)
	assert.False(t, processed)
}

func TestIdempotencyStore_TTLExpiry(t *testing.T) {
	_, store, mr := newTestIdempotencyStore(t)

	ctx := context.Background()
	const (
		eventID = "evt-expire"
		handler = "consent-event"
	)

	marked, err := store.MarkProcessed(ctx, eventID, handler, time.Second)
	require.NoError(t, err)
	assert.True(t, marked)

	mr.FastForward(2 * time.Second)

	processed, err := store.IsProcessed(ctx, eventID, handler)
	require.NoError(t, err)
	assert.False(t, processed)

	marked, err = store.MarkProcessed(ctx, eventID, handler, time.Minute)
	require.NoError(t, err)
	assert.True(t, marked)
}

func TestIdempotencyStore_DifferentEventsScopedSeparately(t *testing.T) {
	_, store, _ := newTestIdempotencyStore(t)

	ctx := context.Background()

	marked, err := store.MarkProcessed(ctx, "evt-001", "consent-event", time.Minute)
	require.NoError(t, err)
	assert.True(t, marked)

	others := []struct {
		eventID string
		handler string
	}{
		{eventID: "evt-002", handler: "consent-event"},
		{eventID: "evt-001", handler: "consent-event-other"},
	}

	for _, o := range others {
		processed, err := store.IsProcessed(ctx, o.eventID, o.handler)
		require.NoError(t, err)
		assert.False(t, processed)

		marked, err := store.MarkProcessed(ctx, o.eventID, o.handler, time.Minute)
		require.NoError(t, err)
		assert.True(t, marked)
	}
}