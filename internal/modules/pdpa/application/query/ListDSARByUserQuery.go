package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ListDSARByUserQuery คิวรีรายการคำขอใช้สิทธิ์ DSAR ทั้งหมดของ User
// ListDSARByUserQuery lists all DSAR requests of a user.
type ListDSARByUserQuery struct {
	dsarRepo repository.DSARRepository
	logger   logger.Logger
}

// NewListDSARByUserQuery สร้าง ListDSARByUserQuery
// NewListDSARByUserQuery creates a ListDSARByUserQuery.
func NewListDSARByUserQuery(dsarRepo repository.DSARRepository, logger logger.Logger) *ListDSARByUserQuery {
	return &ListDSARByUserQuery{
		dsarRepo: dsarRepo,
		logger:   logger,
	}
}

// Execute ส่งคืนรายการคำขอ DSAR ของ userID พร้อมจำนวนรวม
// Execute returns every DSAR request of the given user along with the total count.
func (q *ListDSARByUserQuery) Execute(ctx context.Context, userID uuid.UUID) (*dto.DSARStatusListResponse, error) {
	q.logger.Debugf("ListDSARByUserQuery: listing dsar requests user=%s", userID)

	requests, err := q.dsarRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list dsar by user: %w", err)
	}

	items := make([]dto.DSARStatusResponse, 0, len(requests))
	for _, req := range requests {
		items = append(items, dto.DSARStatusResponse{
			ID:              req.ID,
			UserID:          req.UserID,
			RequestType:     req.RequestType,
			Status:          req.Status,
			RequestedAt:     req.RequestedAt,
			CompletedAt:     req.CompletedAt,
			RejectionReason: req.RejectionReason,
		})
	}

	return &dto.DSARStatusListResponse{
		Requests: items,
		Total:    int64(len(requests)),
	}, nil
}