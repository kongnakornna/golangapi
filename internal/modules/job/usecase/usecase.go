package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/job"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type jobUseCase struct {
	usecase.UseCase[models.Job]
	pgRepo job.JobPgRepository
}

func CreateJobUseCaseI(
	pgRepo job.JobPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) job.JobUseCaseI {
	return &jobUseCase{
		UseCase: usecase.CreateUseCase[models.Job](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *jobUseCase) ChangeStatus(ctx context.Context, id uuid.UUID, status string, reason *string) (*models.Job, error) {
	return u.pgRepo.ChangeStatus(ctx, id, status, reason)
}

func (u *jobUseCase) AddService(ctx context.Context, svc *models.JobService) (*models.JobService, error) {
	return u.pgRepo.AddService(ctx, svc)
}

func (u *jobUseCase) AddPart(ctx context.Context, part *models.JobPartSales) (*models.JobPartSales, error) {
	return u.pgRepo.AddPart(ctx, part)
}

func (u *jobUseCase) GetReport(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	return u.pgRepo.GetReport(ctx, id)
}

func (u *jobUseCase) GetStatusHistory(ctx context.Context, jobID uuid.UUID) ([]*models.JobStatusHistory, error) {
	return u.pgRepo.GetStatusHistory(ctx, jobID)
}

func (u *jobUseCase) GetServices(ctx context.Context, jobID uuid.UUID) ([]*models.JobService, error) {
	return u.pgRepo.GetServices(ctx, jobID)
}

func (u *jobUseCase) GetParts(ctx context.Context, jobID uuid.UUID) ([]*models.JobPartSales, error) {
	return u.pgRepo.GetParts(ctx, jobID)
}

func (u *jobUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}
