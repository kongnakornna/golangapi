package server

import (
	"context"
	"time"

	apimanagerUsecase "icmongolang/internal/modules/apimanager/usecase"
	redisDb "icmongolang/pkg/db/redis"

	"github.com/redis/go-redis/v9"
)

// apiCacheAdapter adapts the existing redisDb.Cache (Get/Set) plus a raw Redis
// client (INCR) to the narrow Cache interface the API manager needs for
// atomic rate-limit counting. Additive and nil-safe.
type apiCacheAdapter struct {
	cache redisDb.Cache
	redis *redis.Client
}

func (a *apiCacheAdapter) Get(ctx context.Context, key string, dst interface{}) error {
	return a.cache.Get(ctx, key, dst)
}

func (a *apiCacheAdapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return a.cache.Set(ctx, key, value, ttl)
}

func (a *apiCacheAdapter) Incr(ctx context.Context, key string) (int64, error) {
	if a.redis == nil {
		return 0, redis.Nil
	}
	return a.redis.Incr(ctx, key).Result()
}

var _ apimanagerUsecase.Cache = (*apiCacheAdapter)(nil)
