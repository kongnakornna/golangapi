package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/batch"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BatchPgRepo implements batch.BatchPgRepository.
// รีโพสิทอรีสำหรับงานแบตช์
type BatchPgRepo struct {
	repository.PgRepo[models.BatchJob]
	DB *gorm.DB
}

// CreateBatchPgRepository creates a new batch repository.
// สร้างรีโพสิทอรีสำหรับงานแบตช์
func CreateBatchPgRepository(db *gorm.DB) batch.BatchPgRepository {
	return &BatchPgRepo{
		PgRepo: repository.CreatePgRepo[models.BatchJob](db),
		DB:     db,
	}
}

// GetLogsByJobId returns all logs for a specific job.
// ดึงบันทึกของงานตาม ID
func (r *BatchPgRepo) GetLogsByJobId(ctx context.Context, jobID uuid.UUID) ([]*models.BatchJobLog, error) {
	var logs []*models.BatchJobLog
	r.DB.WithContext(ctx).Where("job_id = ?", jobID.String()).Order("created_at ASC").Find(&logs)
	return logs, nil
}

// CreateLog creates a new batch job log entry.
// สร้างบันทึกการทำงานของแบตช์
func (r *BatchPgRepo) CreateLog(ctx context.Context, log *models.BatchJobLog) (*models.BatchJobLog, error) {
	if result := r.DB.WithContext(ctx).Create(log); result.Error != nil {
		return nil, result.Error
	}
	return log, nil
}
