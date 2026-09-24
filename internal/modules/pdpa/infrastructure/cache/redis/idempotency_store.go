package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const idemKeyPrefix = "pdpa:idem"

// IdempotencyStore deduplicates event processing by handler and event ID.
type IdempotencyStore struct {
	rdb *redis.Client
}

// NewIdempotencyStore creates an IdempotencyStore backed by the given Redis client.
func NewIdempotencyStore(client *redis.Client) *IdempotencyStore {
	return &IdempotencyStore{rdb: client}
}

// IsProcessed reports whether the event ID was already processed by the handler.
func (s *IdempotencyStore) IsProcessed(ctx context.Context, eventID, handler string) (bool, error) {
	exists, err := s.rdb.Exists(ctx, idemKey(eventID, handler)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// MarkProcessed atomically records the event ID as processed by the handler.
// It returns false when the event was already processed.
func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID, handler string, ttl time.Duration) (bool, error) {
	ok, err := s.rdb.SetNX(ctx, idemKey(eventID, handler), "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func idemKey(eventID, handler string) string {
	return idemKeyPrefix + ":" + handler + ":" + eventID
}