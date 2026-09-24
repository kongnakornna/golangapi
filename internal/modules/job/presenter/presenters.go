package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// JobCreate represents the payload for creating a new job.
// คำขอสร้างใบรับงานซ่อม
type JobCreate struct {
	JobNo         string    `json:"jobNo" validate:"required"`
	CustomerID    uuid.UUID `json:"customerId" validate:"required"`
	CarID         uuid.UUID `json:"carId" validate:"required"`
	MechanicID    uuid.UUID `json:"mechanicId" validate:"required"`
	Symptom       *string   `json:"symptom,omitempty"`
	DiagnosisNote *string   `json:"diagnosisNote,omitempty"`
	Mileage       *int      `json:"mileage,omitempty"`
	EstimatedCost *float64  `json:"estimatedCost,omitempty"`
	Priority      string    `json:"priority"`
}

// JobUpdate represents the payload for updating a job.
// คำขอแก้ไขใบรับงานซ่อม
type JobUpdate struct {
	MechanicID    *uuid.UUID `json:"mechanicId,omitempty"`
	Symptom       *string    `json:"symptom,omitempty"`
	DiagnosisNote *string    `json:"diagnosisNote,omitempty"`
	Mileage       *int       `json:"mileage,omitempty"`
	EstimatedCost *float64   `json:"estimatedCost,omitempty"`
	ActualCost    *float64   `json:"actualCost,omitempty"`
	Priority      *string    `json:"priority,omitempty"`
}

// JobStatusChange represents a status change request.
// คำขอเปลี่ยนสถานะใบรับงานซ่อม
type JobStatusChange struct {
	Status string  `json:"status" validate:"required"`
	Reason *string `json:"reason,omitempty"`
}

// JobResponse represents a job record returned to the client.
// ข้อมูลใบรับงานซ่อมที่ส่งกลับไปยังผู้ใช้
type JobResponse struct {
	ID            uuid.UUID          `json:"id"`
	JobNo         string             `json:"jobNo"`
	CustomerID    uuid.UUID          `json:"customerId"`
	CarID         uuid.UUID          `json:"carId"`
	MechanicID    uuid.UUID          `json:"mechanicId"`
	Status        string             `json:"status"`
	StartDate     time.Time          `json:"startDate"`
	EndDate       *time.Time         `json:"endDate,omitempty"`
	Symptom       *string            `json:"symptom,omitempty"`
	DiagnosisNote *string            `json:"diagnosisNote,omitempty"`
	Mileage       *int               `json:"mileage,omitempty"`
	EstimatedCost *float64           `json:"estimatedCost,omitempty"`
	ActualCost    *float64           `json:"actualCost,omitempty"`
	Priority      string             `json:"priority"`
	UserID        uuid.UUID          `json:"userId"`
	WhitelabelID  uuid.UUID          `json:"whitelabelId"`
	CreatedAt     helpers.LocalTime  `json:"createdAt"`
	UpdatedAt     *helpers.LocalTime `json:"updatedAt,omitempty"`
}

// PaginatedJobResponse holds paginated job results.
// ผลลัพธ์ใบรับงานซ่อมแบบแบ่งหน้า
type PaginatedJobResponse struct {
	Jobs       []*JobResponse `json:"jobs"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// JobServiceRequest represents a service item to add to a job.
// คำขอเพิ่มบริการในใบรับงานซ่อม
type JobServiceRequest struct {
	ServiceID uuid.UUID `json:"serviceId" validate:"required"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unitPrice" validate:"required"`
	Discount  float64   `json:"discount"`
	Note      *string   `json:"note,omitempty"`
}

// JobServiceResponse represents a service item on a job.
// ข้อมูลบริการในใบรับงานซ่อม
type JobServiceResponse struct {
	ID        uuid.UUID         `json:"id"`
	JobID     uuid.UUID         `json:"jobId"`
	ServiceID uuid.UUID         `json:"serviceId"`
	Quantity  int               `json:"quantity"`
	UnitPrice float64           `json:"unitPrice"`
	Discount  float64           `json:"discount"`
	Note      *string           `json:"note,omitempty"`
	CreatedAt helpers.LocalTime `json:"createdAt"`
}

// JobPartRequest represents a part to add to a job.
// คำขอเพิ่มอะไหล่ในใบรับงานซ่อม
type JobPartRequest struct {
	PartID    uuid.UUID `json:"partId" validate:"required"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unitPrice" validate:"required"`
	Discount  float64   `json:"discount"`
	Note      *string   `json:"note,omitempty"`
}

// JobPartResponse represents a part used in a job.
// ข้อมูลอะไหล่ในใบรับงานซ่อม
type JobPartResponse struct {
	ID        uuid.UUID         `json:"id"`
	JobID     uuid.UUID         `json:"jobId"`
	PartID    uuid.UUID         `json:"partId"`
	Quantity  int               `json:"quantity"`
	UnitPrice float64           `json:"unitPrice"`
	Discount  float64           `json:"discount"`
	Note      *string           `json:"note,omitempty"`
	CreatedAt helpers.LocalTime `json:"createdAt"`
}

// JobStatusHistoryResponse represents a status change record.
// ประวัติการเปลี่ยนสถานะ
type JobStatusHistoryResponse struct {
	ID         uuid.UUID `json:"id"`
	JobID      uuid.UUID `json:"jobId"`
	FromStatus *string   `json:"fromStatus,omitempty"`
	ToStatus   string    `json:"toStatus"`
	ChangedBy  uuid.UUID `json:"changedBy"`
	ChangedAt  time.Time `json:"changedAt"`
	Reason     *string   `json:"reason,omitempty"`
}
