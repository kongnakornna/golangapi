package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// GetDSARStatusQuery คิวรีสถานะคำขอใช้สิทธิ์ DSAR ตาม ID
// GetDSARStatusQuery queries the status of a DSAR request by its ID.
type GetDSARStatusQuery struct {
	dsarRepo repository.DSARRepository
	logger   logger.Logger
}

// NewGetDSARStatusQuery สร้าง GetDSARStatusQuery
// NewGetDSARStatusQuery creates a GetDSARStatusQuery.
func NewGetDSARStatusQuery(dsarRepo repository.DSARRepository, logger logger.Logger) *GetDSARStatusQuery {
	return &GetDSARStatusQuery{
		dsarRepo: dsarRepo,
		logger:   logger,
	}
}

// Execute ดึงสถานะคำขอใช้สิทธิ์ DSAR ตาม id
// Execute returns the DSAR request status for the given id.
func (q *GetDSARStatusQuery) Execute(ctx context.Context, id uuid.UUID) (*dto.DSARStatusResponse, error) {
	q.logger.Debugf("GetDSARStatusQuery: querying dsar request id=%s", id)

	req, err := q.dsarRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get dsar status: find by id: %w", err)
	}
	if req == nil {
		q.logger.Warnf("GetDSARStatusQuery: dsar request not found id=%s", id)
		return nil, domainerrors.ErrDSARNotFound
	}

	return &dto.DSARStatusResponse{
		ID:              req.ID,
		UserID:          req.UserID,
		RequestType:     req.RequestType,
		Status:          req.Status,
		RequestedAt:     req.RequestedAt,
		CompletedAt:     req.CompletedAt,
		RejectionReason: req.RejectionReason,
	}, nil
}