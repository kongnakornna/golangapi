package template

import (
	"icmongolang/internal"
	"icmongolang/internal/template/models"
)

type UseCaseI interface {
	internal.UseCaseI[models.Template]
}
