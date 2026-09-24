package users

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"
)

type UserRedisRepository interface {
	internal.RedisRepository[models.SdUser]
	SetDetail(ctx context.Context, key string, exp *UserProfileDetail, seconds int) error
	GetDetail(ctx context.Context, key string) (*UserProfileDetail, error)
}
