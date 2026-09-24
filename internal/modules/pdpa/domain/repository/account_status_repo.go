package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
)

type UserAccountStatusRepository interface {
	Save(ctx context.Context, status *entity.UserAccountStatus) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAccountStatus, error)
	FindSuspendedWithoutDeletion(ctx context.Context, cutoffDate time.Time) ([]entity.UserAccountStatus, error)
	UpdateStatus(ctx context.Context, userID uuid.UUID, status valueobject.AccountStatus) error
	UpdateDeletionConfirmed(ctx context.Context, userID uuid.UUID, confirmedAt time.Time) error
	UpdateAutoDeleted(ctx context.Context, userID uuid.UUID, deletedAt time.Time) error
}