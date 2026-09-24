package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/batch"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type batchUseCase struct {
	usecase.UseCase[models.BatchJob]
	pgRepo batch.BatchPgRepository
}

// CreateBatchUseCaseI creates a new batch use case instance.
// สร้างอินสแตนซ์สำหรับธุรกิจแบตช์
func CreateBatchUseCaseI(
	pgRepo batch.BatchPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) batch.BatchUseCaseI {
	return &batchUseCase{
		UseCase: usecase.CreateUseCase[models.BatchJob](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *batchUseCase) RunJob(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement actual job execution logic
	_, err := u.pgRepo.CreateLog(ctx, &models.BatchJobLog{
		JobID:   id,
		Message: "Job triggered manually",
		Level:   "info",
	})
	return err
}

func (u *batchUseCase) GetJobLogs(ctx context.Context, jobID uuid.UUID) ([]*models.BatchJobLog, error) {
	return u.pgRepo.GetLogsByJobId(ctx, jobID)
}
