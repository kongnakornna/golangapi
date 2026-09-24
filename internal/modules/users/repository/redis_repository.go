package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/repository"
	"icmongolang/internal/modules/users"

	"github.com/redis/go-redis/v9"
)

type UserRedisRepo struct {
	repository.RedisRepo[models.SdUser]
}

func CreateUserRedisRepository(redisClient *redis.Client) users.UserRedisRepository {
	return &UserRedisRepo{
		RedisRepo: repository.CreateRedisRepo[models.SdUser](redisClient),
	}
}

func (r *UserRedisRepo) SetDetail(ctx context.Context, key string, d *users.UserProfileDetail, seconds int) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return r.RedisClient.Set(ctx, key, b, time.Second*time.Duration(seconds)).Err()
}

func (r *UserRedisRepo) GetDetail(ctx context.Context, key string) (*users.UserProfileDetail, error) {
	b, err := r.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var d users.UserProfileDetail
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return &d, nil
}
