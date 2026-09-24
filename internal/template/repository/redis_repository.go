package repository

import (
	"context"
	"encoding/json"
	"time"

	"icmongolang/internal/template"
	"icmongolang/internal/template/models"

	"github.com/redis/go-redis/v9"
)

type templateRedisRepo struct {
	rdb *redis.Client
}

func NewTemplateRedisRepository(rdb *redis.Client) template.RedisRepository {
	return &templateRedisRepo{rdb: rdb}
}

func (r *templateRedisRepo) Create(ctx context.Context, key string, exp *models.Template, seconds int) error {
	data, err := json.Marshal(exp)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, data, time.Duration(seconds)*time.Second).Err()
}

func (r *templateRedisRepo) Get(ctx context.Context, key string) (*models.Template, error) {
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var obj models.Template
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (r *templateRedisRepo) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *templateRedisRepo) Sadd(ctx context.Context, key string, value string) error {
	return r.rdb.SAdd(ctx, key, value).Err()
}

func (r *templateRedisRepo) Sadds(ctx context.Context, key string, values []string) error {
	items := make([]interface{}, len(values))
	for i, v := range values {
		items[i] = v
	}
	return r.rdb.SAdd(ctx, key, items...).Err()
}

func (r *templateRedisRepo) Srem(ctx context.Context, key string, value string) error {
	return r.rdb.SRem(ctx, key, value).Err()
}

func (r *templateRedisRepo) SIsMember(ctx context.Context, key string, value string) (bool, error) {
	return r.rdb.SIsMember(ctx, key, value).Result()
}

func (r *templateRedisRepo) Expire(ctx context.Context, key string, seconds int) error {
	return r.rdb.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}
