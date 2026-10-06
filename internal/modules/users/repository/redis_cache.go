package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"icmongolang/internal/modules/users"

	"github.com/redis/go-redis/v9"
)

type cacheRepo struct{ client *redis.Client }

func CreateCacheRepository(client *redis.Client) users.CacheRepository {
	return &cacheRepo{client: client}
}

func (r *cacheRepo) GetJSON(ctx context.Context, key string, dst interface{}) (bool, error) {
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) { return false, nil }
		return false, err
	}
	if err := json.Unmarshal(b, dst); err != nil { return false, err }
	return true, nil
}

func (r *cacheRepo) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil { return err }
	return r.client.Set(ctx, key, b, ttl).Err()
}

func (r *cacheRepo) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 { return nil }
	return r.client.Del(ctx, keys...).Err()
}

func (r *cacheRepo) DeleteByPattern(ctx context.Context, pattern string) error {
	_, err := r.ScanAndDelete(ctx, pattern, 200)
	return err
}

func (r *cacheRepo) ScanAndDelete(ctx context.Context, pattern string, batchSize int) (int, error) {
	if batchSize <= 0 { batchSize = 200 }
	var cursor uint64
	deleted := 0
	for {
		keys, next, err := r.client.Scan(ctx, cursor, pattern, int64(batchSize)).Result()
		if err != nil { return deleted, err }
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil { return deleted, err }
			deleted += len(keys)
		}
		cursor = next
		if cursor == 0 { break }
	}
	return deleted, nil
}

var _ users.CacheRepository = (*cacheRepo)(nil)