// ConfirmDeletionCommand confirms permanent deletion of user data, either
// immediately or once the retention period has elapsed.
// คำสั่งยืนยันการลบข้อมูลของผู้ใช้แบบถาวร ทั้งการลบทันทีหรือเมื่อครบระยะเวลาเก็บรักษา
package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// ConfirmDeletionCommand handles deletion confirmation.
// จัดการการยืนยันการลบข้อมูล
type ConfirmDeletionCommand struct {
	accountRepo    repository.UserAccountStatusRepository
	deletionPolicy *service.DeletionPolicyService
	producer       *kafka.Producer
	cfg            *config.Config
	logger         logger.Logger
}

// NewConfirmDeletionCommand creates a new ConfirmDeletionCommand.
// สร้างคำสั่งยืนยันการลบข้อมูลใหม่
func NewConfirmDeletionCommand(accountRepo repository.UserAccountStatusRepository, deletionPolicy *service.DeletionPolicyService, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *ConfirmDeletionCommand {
	return &ConfirmDeletionCommand{
		accountRepo:    accountRepo,
		deletionPolicy: deletionPolicy,
		producer:       producer,
		cfg:            cfg,
		logger:         logger,
	}
}

// Execute confirms deletion, evaluates whether immediate deletion is permitted,
// publishes the data.deletion_requested event and persists the status.
// ยืนยันการลบ ประเมินว่าสามารถลบทันทีได้หรือไม่ เผยแพร่เหตุการณ์ data.deletion_requested และบันทึกสถานะ
func (c *ConfirmDeletionCommand) Execute(ctx context.Context, req *dto.ConfirmDeletionRequest) error {
	account, err := c.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrAccountNotFound) {
			return domainerrors.ErrAccountNotFound
		}
		return err
	}
	if account == nil {
		return domainerrors.ErrAccountNotFound
	}

	account.ConfirmDeletion()

	deletionConfirmedAt := time.Now()
	if account.DeletionConfirmedAt != nil {
		deletionConfirmedAt = *account.DeletionConfirmedAt
	}

	immediateDeletion := c.deletionPolicy.CanImmediateDeletion(account)
	c.logger.Infof("deletion decision: user_id=%s immediate_deletion=%t", req.UserID.String(), immediateDeletion)

	payload := map[string]interface{}{
		"user_id":               req.UserID.String(),
		"deletion_confirmed_at": deletionConfirmedAt.Format(time.RFC3339),
		"immediate_deletion":    immediateDeletion,
	}
	if err := c.producer.PublishMessage(req.UserID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.accountRepo.Save(ctx, account); err != nil {
		return err
	}

	c.logger.Infof("deletion confirmed: user_id=%s", req.UserID.String())
	return nil
}