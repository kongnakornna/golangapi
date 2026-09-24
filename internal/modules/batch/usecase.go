package batch

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

// BatchUseCaseI defines business logic methods for batch jobs.
// อินเทอร์เฟซธุรกิจสำหรับงานแบตช์
type BatchUseCaseI interface {
	internal.UseCaseI[models.BatchJob]
	RunJob(ctx context.Context, id uuid.UUID) error
	GetJobLogs(ctx context.Context, jobID uuid.UUID) ([]*models.BatchJobLog, error)
}
