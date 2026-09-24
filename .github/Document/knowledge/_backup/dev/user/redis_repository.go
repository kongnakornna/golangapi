package users

import (
	"icmongolang/internal"
	"icmongolang/internal/models"
)

type UserRedisRepository interface {
	internal.RedisRepository[models.User]
}
