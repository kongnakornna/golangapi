package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// BatchJobRequest represents a batch job creation payload.
// คำขอสร้างงานแบตช์
type BatchJobRequest struct {
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Config   string `json:"config"`
	Schedule string `json:"schedule"`
}

// BatchJobResponse represents a batch job record returned to the client.
// ข้อมูลงานแบตช์
type BatchJobResponse struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Status       string            `json:"status"`
	Config       string            `json:"config"`
	Schedule     *string           `json:"schedule,omitempty"`
	TotalCount   int               `json:"totalCount"`
	SuccessCount int               `json:"successCount"`
	FailCount    int               `json:"failCount"`
	StartedAt    *time.Time        `json:"startedAt,omitempty"`
	FinishedAt   *time.Time        `json:"finishedAt,omitempty"`
	CreatedAt    helpers.LocalTime `json:"createdAt"`
}

// BatchJobLogResponse represents a batch job log entry.
// บันทึกการทำงานของแบตช์
type BatchJobLogResponse struct {
	ID        int               `json:"id"`
	JobID     uuid.UUID         `json:"jobId"`
	Message   string            `json:"message"`
	Level     string            `json:"level"`
	CreatedAt helpers.LocalTime `json:"createdAt"`
}
