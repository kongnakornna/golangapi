package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

// ListConsentLogByUserQuery คิวรีประวัติรายการยินยอมทั้งหมดของ User
// ListConsentLogByUserQuery lists all consent logs of a user.
type ListConsentLogByUserQuery struct {
	consentRepo repository.ConsentRepository
	logger      logger.Logger
}

// NewListConsentLogByUserQuery สร้าง ListConsentLogByUserQuery
// NewListConsentLogByUserQuery creates a ListConsentLogByUserQuery.
func NewListConsentLogByUserQuery(consentRepo repository.ConsentRepository, logger logger.Logger) *ListConsentLogByUserQuery {
	return &ListConsentLogByUserQuery{
		consentRepo: consentRepo,
		logger:      logger,
	}
}

// Execute ส่งคืนประวัติสิทธิ์ยินยอมของ userID พร้อมจำนวนรวม
// Execute returns every consent log of the given user along with the total count.
func (q *ListConsentLogByUserQuery) Execute(ctx context.Context, userID uuid.UUID) (*dto.ConsentStatusListResponse, error) {
	q.logger.Debugf("ListConsentLogByUserQuery: listing consent logs user=%s", userID)

	consents, err := q.consentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list consent logs by user: %w", err)
	}

	items := make([]dto.ConsentStatusResponse, 0, len(consents))
	for _, consent := range consents {
		items = append(items, dto.ConsentStatusResponse{
			UserID:    consent.UserID,
			Purpose:   consent.Purpose,
			Status:    consent.Status,
			IsActive:  consent.IsActive(),
			GrantedAt: consent.ConsentedAt,
			ExpiresAt: consent.ExpiresAt,
		})
	}

	return &dto.ConsentStatusListResponse{
		Consents: items,
		Total:    int64(len(consents)),
	}, nil
}