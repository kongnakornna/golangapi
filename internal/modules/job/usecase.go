package job

import (
	"context"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

// JobUseCaseI defines business logic methods for the Job module.
// อินเทอร์เฟซธุรกิจสำหรับโมดูลใบรับงานซ่อม
type JobUseCaseI interface {
	internal.UseCaseI[models.Job]
	ChangeStatus(ctx context.Context, id uuid.UUID, status string, reason *string) (*models.Job, error)
	AddService(ctx context.Context, svc *models.JobService) (*models.JobService, error)
	AddPart(ctx context.Context, part *models.JobPartSales) (*models.JobPartSales, error)
	GetReport(ctx context.Context, id uuid.UUID) (*models.Job, error)
	GetStatusHistory(ctx context.Context, jobID uuid.UUID) ([]*models.JobStatusHistory, error)
	GetServices(ctx context.Context, jobID uuid.UUID) ([]*models.JobService, error)
	GetParts(ctx context.Context, jobID uuid.UUID) ([]*models.JobPartSales, error)
	Count(ctx context.Context) (int64, error)
}
