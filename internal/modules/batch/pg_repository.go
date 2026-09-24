package batch

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

// BatchPgRepository defines data access methods for batch jobs.
// ดึงข้อมูลงานแบตช์จากฐานข้อมูล
type BatchPgRepository interface {
	internal.PgRepository[models.BatchJob]
	GetLogsByJobId(ctx context.Context, jobID uuid.UUID) ([]*models.BatchJobLog, error)
	CreateLog(ctx context.Context, log *models.BatchJobLog) (*models.BatchJobLog, error)
}
