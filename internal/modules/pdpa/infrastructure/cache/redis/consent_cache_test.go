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

func newTestConsentCache(t *testing.T) (*redis.Client, *ConsentCacheImpl, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	cache := NewConsentCache(client, 5*time.Minute)

	return client, cache, mr
}

func TestConsentCache_SetAndGet(t *testing.T) {
	_, cache, _ := newTestConsentCache(t)

	err := cache.Set(context.Background(), "user-1", "marketing", "granted")
	require.NoError(t, err)

	status, err := cache.Get(context.Background(), "user-1", "marketing")
	require.NoError(t, err)
	assert.Equal(t, "granted", status)
}

func TestConsentCache_GetMiss(t *testing.T) {
	_, cache, _ := newTestConsentCache(t)

	status, err := cache.Get(context.Background(), "user-1", "marketing")
	require.NoError(t, err)
	assert.Empty(t, status)
}

func TestConsentCache_ScopedPerUserAndPurpose(t *testing.T) {
	_, cache, _ := newTestConsentCache(t)

	require.NoError(t, cache.Set(context.Background(), "user-1", "marketing", "granted"))
	require.NoError(t, cache.Set(context.Background(), "user-1", "analytics", "denied"))
	require.NoError(t, cache.Set(context.Background(), "user-2", "marketing", "granted"))

	status, err := cache.Get(context.Background(), "user-1", "analytics")
	require.NoError(t, err)
	assert.Equal(t, "denied", status)
}

func TestConsentCache_DeleteAllByUser(t *testing.T) {
	_, cache, _ := newTestConsentCache(t)
	ctx := context.Background()

	require.NoError(t, cache.Set(ctx, "user-1", "marketing", "granted"))
	require.NoError(t, cache.Set(ctx, "user-1", "analytics", "denied"))
	require.NoError(t, cache.Set(ctx, "user-2", "marketing", "granted"))

	require.NoError(t, cache.DeleteAllByUser(ctx, "user-1"))

	status, err := cache.Get(ctx, "user-1", "marketing")
	require.NoError(t, err)
	assert.Empty(t, status)

	status, err = cache.Get(ctx, "user-1", "analytics")
	require.NoError(t, err)
	assert.Empty(t, status)

	status, err = cache.Get(ctx, "user-2", "marketing")
	require.NoError(t, err)
	assert.Equal(t, "granted", status)
}

func TestConsentCache_TTLExpiry(t *testing.T) {
	_, cache, mr := newTestConsentCache(t)

	require.NoError(t, cache.Set(context.Background(), "user-1", "marketing", "granted"))

	mr.FastForward(10 * time.Minute)

	status, err := cache.Get(context.Background(), "user-1", "marketing")
	require.NoError(t, err)
	assert.Empty(t, status)
}