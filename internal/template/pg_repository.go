package template

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/template/models"
)

type PgRepository interface {
	internal.PgRepository[models.Template]
	GetByName(ctx context.Context, name string) (*models.Template, error)
}
