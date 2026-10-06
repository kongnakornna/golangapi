// Package redis provides Redis-backed caches for the PDPA module.
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const consentKeyPrefix = "pdpa:consent"

type ConsentCacheImpl struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewConsentCache(client *redis.Client, ttl time.Duration) *ConsentCacheImpl {
	return &ConsentCacheImpl{rdb: client, ttl: ttl}
}

func (c *ConsentCacheImpl) Set(ctx context.Context, userID, purpose, status string) error {
	return c.rdb.SetEx(ctx, consentKey(userID, purpose), status, c.ttl).Err()
}

func (c *ConsentCacheImpl) Get(ctx context.Context, userID, purpose string) (string, error) {
	status, err := c.rdb.Get(ctx, consentKey(userID, purpose)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return status, err
}

func (c *ConsentCacheImpl) DeleteAllByUser(ctx context.Context, userID string) error {
	pattern := consentKeyPrefix + ":" + userID + ":*"
	var cursor uint64
	for {
		keys, next, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func consentKey(userID, purpose string) string {
	return consentKeyPrefix + ":" + userID + ":" + purpose
}
