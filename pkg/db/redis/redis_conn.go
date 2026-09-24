package redis

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.Db,
	})
	return client
}

const redisOpTimeout = 5 * time.Second

func redisCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), redisOpTimeout)
}

type RedisHelper struct {
	client *redis.Client
}

var instance *RedisHelper

func GetRedisHelper() *RedisHelper {
	if instance == nil {
		instance = &RedisHelper{
			client: redis.NewClient(&redis.Options{
				Addr:     "localhost:6380",
				Password: "",
				DB:       0,
			}),
		}
	}
	return instance
}

func (r *RedisHelper) Set(key string, value interface{}, ttlSeconds int) error {
	ctx, cancel := redisCtx()
	defer cancel()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisHelper) Get(key string, dest interface{}) error {
	ctx, cancel := redisCtx()
	defer cancel()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *RedisHelper) Delete(key string) error {
	ctx, cancel := redisCtx()
	defer cancel()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisHelper) Exists(key string) (bool, error) {
	ctx, cancel := redisCtx()
	defer cancel()
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *RedisHelper) TTL(key string) (time.Duration, error) {
	ctx, cancel := redisCtx()
	defer cancel()
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisHelper) Keys(pattern string) ([]string, error) {
	ctx, cancel := redisCtx()
	defer cancel()
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisHelper) FlushAll() error {
	ctx, cancel := redisCtx()
	defer cancel()
	return r.client.FlushAll(ctx).Err()
}

// ---------- Cache interface สำหรับ MQTT handler ----------
type Cache interface {
	Get(ctx context.Context, key string, dst interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type RedisCache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, dst interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}
