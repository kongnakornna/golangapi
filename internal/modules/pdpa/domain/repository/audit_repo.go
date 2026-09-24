package repository

import (
	"context"

	"icmongolang/internal/modules/pdpa/domain/entity"
)

type AuditRepository interface {
	Save(ctx context.Context, audit *entity.AuditTrail) error
}