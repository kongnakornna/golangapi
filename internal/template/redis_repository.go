package template

import (
	"icmongolang/internal"
	"icmongolang/internal/template/models"
)

type RedisRepository interface {
	internal.RedisRepository[models.Template]
}
