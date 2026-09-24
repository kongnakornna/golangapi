package repository

import (
	"context"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/job"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobPgRepo implements job.JobPgRepository.
// รีโพสิทอรีสำหรับใบรับงานซ่อม
type JobPgRepo struct {
	repository.PgRepo[models.Job]
	DB *gorm.DB
}

// CreateJobPgRepository creates a new Job repository.
// สร้างรีโพสิทอรีสำหรับใบรับงานซ่อม
func CreateJobPgRepository(db *gorm.DB) job.JobPgRepository {
	return &JobPgRepo{
		PgRepo: repository.CreatePgRepo[models.Job](db),
		DB:     db,
	}
}

// ChangeStatus updates the job status and records a history entry.
// อัปเดตสถานะใบรับงานซ่อมและบันทึกประวัติ
func (r *JobPgRepo) ChangeStatus(ctx context.Context, id uuid.UUID, status string, reason *string) (*models.Job, error) {
	jobObj, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	fromStatus := jobObj.Status
	now := time.Now()

	if err := r.DB.WithContext(ctx).Model(&models.Job{}).Where("id = ?", id.String()).Update("status", status).Error; err != nil {
		return nil, err
	}

	history := &models.JobStatusHistory{
		JobID:       id,
		FromStatus:  &fromStatus,
		ToStatus:    status,
		ChangedBy:   jobObj.UserID,
		ChangedAt:   now,
		Reason:      reason,
		WhitelabelID: jobObj.WhitelabelID,
	}

	if err := r.DB.WithContext(ctx).Create(history).Error; err != nil {
		return nil, err
	}

	jobObj.Status = status
	return jobObj, nil
}

// AddService creates a new service record for a job.
// เพิ่มบริการในใบรับงานซ่อม
func (r *JobPgRepo) AddService(ctx context.Context, svc *models.JobService) (*models.JobService, error) {
	if result := r.DB.WithContext(ctx).Create(svc); result.Error != nil {
		return nil, result.Error
	}
	return svc, nil
}

// AddPart creates a new part sales record for a job.
// เพิ่มอะไหล่ในใบรับงานซ่อม
func (r *JobPgRepo) AddPart(ctx context.Context, part *models.JobPartSales) (*models.JobPartSales, error) {
	if result := r.DB.WithContext(ctx).Create(part); result.Error != nil {
		return nil, result.Error
	}
	return part, nil
}

// GetReport retrieves a job with related services and parts.
// ดึงข้อมูลใบรับงานซ่อมพร้อมบริการและอะไหล่
func (r *JobPgRepo) GetReport(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	return r.Get(ctx, id)
}

// GetStatusHistory returns the status change history for a job.
// ดึงประวัติการเปลี่ยนสถานะของใบรับงานซ่อม
func (r *JobPgRepo) GetStatusHistory(ctx context.Context, jobID uuid.UUID) ([]*models.JobStatusHistory, error) {
	var history []*models.JobStatusHistory
	if result := r.DB.WithContext(ctx).Where("job_id = ?", jobID.String()).Order("changed_at DESC").Find(&history); result.Error != nil {
		return nil, result.Error
	}
	return history, nil
}

// GetServices returns all service items for a job.
// ดึงรายการบริการของใบรับงานซ่อม
func (r *JobPgRepo) GetServices(ctx context.Context, jobID uuid.UUID) ([]*models.JobService, error) {
	var services []*models.JobService
	if result := r.DB.WithContext(ctx).Where("job_id = ?", jobID.String()).Find(&services); result.Error != nil {
		return nil, result.Error
	}
	return services, nil
}

// GetParts returns all part sales for a job.
// ดึงรายการอะไหล่ของใบรับงานซ่อม
func (r *JobPgRepo) GetParts(ctx context.Context, jobID uuid.UUID) ([]*models.JobPartSales, error) {
	var parts []*models.JobPartSales
	if result := r.DB.WithContext(ctx).Where("job_id = ?", jobID.String()).Find(&parts); result.Error != nil {
		return nil, result.Error
	}
	return parts, nil
}

func (r *JobPgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&models.Job{}).Where("deleted = ?", false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
