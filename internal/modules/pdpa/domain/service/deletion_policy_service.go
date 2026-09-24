package service

import (
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/valueobject"
)

type DeletionPolicyService struct{}

func NewDeletionPolicyService() *DeletionPolicyService {
	return &DeletionPolicyService{}
}

func (s *DeletionPolicyService) CanImmediateDeletion(status *entity.UserAccountStatus) bool {
	if status == nil {
		return false
	}
	return status.Status == valueobject.AccountTerminated && status.DeletionConfirmedAt != nil
}

func (s *DeletionPolicyService) IsReadyForAutoDeletion(status *entity.UserAccountStatus, now time.Time) bool {
	if status == nil {
		return false
	}
	if status.Status != valueobject.AccountSuspended {
		return false
	}
	if status.DeletionConfirmedAt == nil || status.RetentionDeadline == nil {
		return false
	}
	return now.After(*status.RetentionDeadline)
}

func (s *DeletionPolicyService) AnonymizeUserData(userID uuid.UUID, data map[string]interface{}) map[string]interface{} {
	anonymized := make(map[string]interface{}, len(data))
	for key, value := range data {
		switch v := value.(type) {
		case string:
			anonymized[key] = "****"
		default:
			anonymized[key] = v
		}
	}
	return anonymized
}