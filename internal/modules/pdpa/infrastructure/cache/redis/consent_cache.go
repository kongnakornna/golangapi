// Package redis provides Redis-backed caches and stores for the PDPA module.
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const consentKeyPrefix = "pdpa:consent"

// ConsentCacheImpl caches a user's consent status by purpose.
type ConsentCacheImpl struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewConsentCache creates a ConsentCacheImpl with the given client and ttl.
func NewConsentCache(client *redis.Client, ttl time.Duration) *ConsentCacheImpl {
	return &ConsentCacheImpl{
		rdb: client,
		ttl: ttl,
	}
}

// Set stores the consent status for the user and purpose with a ttl.
func (c *ConsentCacheImpl) Set(ctx context.Context, userID string, purpose string, status string) error {
	key := consentKey(userID, purpose)
	return c.rdb.SetEx(ctx, key, status, c.ttl).Err()
}

// Get returns the consent status for the user and purpose, or "" when absent.
func (c *ConsentCacheImpl) Get(ctx context.Context, userID string, purpose string) (string, error) {
	key := consentKey(userID, purpose)
	status, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return status, nil
}

// DeleteAllByUser removes all cached consent entries for the user.
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

func consentKey(userID string, purpose string) string {
	return consentKeyPrefix + ":" + userID + ":" + purpose
}