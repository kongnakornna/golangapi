package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const idemKeyPrefix = "pdpa:idem"

type IdempotencyStore struct {
	rdb *redis.Client
}

func NewIdempotencyStore(client *redis.Client) *IdempotencyStore {
	return &IdempotencyStore{rdb: client}
}

func (s *IdempotencyStore) IsProcessed(ctx context.Context, eventID, handler string) (bool, error) {
	exists, err := s.rdb.Exists(ctx, idemKey(eventID, handler)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID, handler string, ttl time.Duration) (bool, error) {
	return s.rdb.SetNX(ctx, idemKey(eventID, handler), "1", ttl).Result()
}

func idemKey(eventID, handler string) string {
	return idemKeyPrefix + ":" + handler + ":" + eventID
}
