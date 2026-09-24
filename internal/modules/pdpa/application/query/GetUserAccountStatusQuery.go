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

// GetUserAccountStatusQuery คิวรีสถานะบัญชีผู้ใช้เพื่อรองรับรอบการลบข้อมูล (retention)
// GetUserAccountStatusQuery queries the account status of a user to support the retention lifecycle.
type GetUserAccountStatusQuery struct {
	accountRepo repository.UserAccountStatusRepository
	logger      logger.Logger
}

// NewGetUserAccountStatusQuery สร้าง GetUserAccountStatusQuery
// NewGetUserAccountStatusQuery creates a GetUserAccountStatusQuery.
func NewGetUserAccountStatusQuery(accountRepo repository.UserAccountStatusRepository, logger logger.Logger) *GetUserAccountStatusQuery {
	return &GetUserAccountStatusQuery{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

// Execute ดึงสถานะบัญชีผู้ใช้ของ userID
// Execute returns the user account status for the given userID.
func (q *GetUserAccountStatusQuery) Execute(ctx context.Context, userID uuid.UUID) (*dto.UserAccountStatusResponse, error) {
	q.logger.Debugf("GetUserAccountStatusQuery: querying account status user=%s", userID)

	acct, err := q.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user account status: find by user id: %w", err)
	}
	if acct == nil {
		q.logger.Warnf("GetUserAccountStatusQuery: account status not found user=%s", userID)
		return nil, domainerrors.ErrAccountNotFound
	}

	return &dto.UserAccountStatusResponse{
		UserID:              acct.UserID,
		Status:              acct.Status,
		RetentionDeadline:   acct.RetentionDeadline,
		DeletionConfirmedAt: acct.DeletionConfirmedAt,
		AutoDeletedAt:       acct.AutoDeletedAt,
	}, nil
}