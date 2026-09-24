// SuspendAccountCommand suspends a user account and starts the retention period.
// คำสั่งระงับบัญชีผู้ใช้และเริ่มนับระยะเวลาการเก็บรักษาข้อมูล
package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// SuspendAccountCommand handles account suspension.
// จัดการการระงับบัญชีผู้ใช้
type SuspendAccountCommand struct {
	accountRepo repository.UserAccountStatusRepository
	producer    *kafka.Producer
	cfg         *config.Config
	logger      logger.Logger
}

// NewSuspendAccountCommand creates a new SuspendAccountCommand.
// สร้างคำสั่งระงับบัญชีผู้ใช้ใหม่
func NewSuspendAccountCommand(accountRepo repository.UserAccountStatusRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *SuspendAccountCommand {
	return &SuspendAccountCommand{
		accountRepo: accountRepo,
		producer:    producer,
		cfg:         cfg,
		logger:      logger,
	}
}

// Execute suspends the user account for the configured retention period,
// publishes the account.suspended event and persists the status.
// ระงับบัญชีผู้ใช้ตามระยะเวลาการเก็บรักษาที่กำหนด เผยแพร่เหตุการณ์ account.suspended และบันทึกสถานะ
func (c *SuspendAccountCommand) Execute(ctx context.Context, req *dto.SuspendAccountRequest) error {
	account, err := c.accountRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		if !errors.Is(err, domainerrors.ErrAccountNotFound) {
			return err
		}
		account = entity.NewUserAccountStatus(req.UserID, valueobject.AccountActive)
	} else if account == nil {
		account = entity.NewUserAccountStatus(req.UserID, valueobject.AccountActive)
	}

	now := time.Now()
	account.SetSuspended(now, req.RetentionYears)

	retentionDeadline := now
	if account.RetentionDeadline != nil {
		retentionDeadline = *account.RetentionDeadline
	}

	payload := map[string]interface{}{
		"user_id":            req.UserID.String(),
		"status":             account.Status.String(),
		"suspended_at":       now.Format(time.RFC3339),
		"retention_deadline": retentionDeadline.Format(time.RFC3339),
		"retention_years":    req.RetentionYears,
	}
	if err := c.producer.PublishMessage(req.UserID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.accountRepo.Save(ctx, account); err != nil {
		return err
	}

	c.logger.Infof("account suspended: user_id=%s retention_years=%d", req.UserID.String(), req.RetentionYears)
	return nil
}