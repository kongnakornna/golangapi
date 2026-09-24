// TerminateAccountCommand terminates a user account immediately.
// คำสั่งยุติบัญชีผู้ใช้ทันที
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

// TerminateAccountCommand handles account termination.
// จัดการการยุติบัญชีผู้ใช้
type TerminateAccountCommand struct {
	accountRepo repository.UserAccountStatusRepository
	producer    *kafka.Producer
	cfg         *config.Config
	logger      logger.Logger
}

// NewTerminateAccountCommand creates a new TerminateAccountCommand.
// สร้างคำสั่งยุติบัญชีผู้ใช้ใหม่
func NewTerminateAccountCommand(accountRepo repository.UserAccountStatusRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *TerminateAccountCommand {
	return &TerminateAccountCommand{
		accountRepo: accountRepo,
		producer:    producer,
		cfg:         cfg,
		logger:      logger,
	}
}

// Execute terminates the user account, publishes the account.terminated event
// and persists the status.
// ยุติบัญชีผู้ใช้ เผยแพร่เหตุการณ์ account.terminated และบันทึกสถานะ
func (c *TerminateAccountCommand) Execute(ctx context.Context, req *dto.TerminateAccountRequest) error {
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
	account.SetTerminated(now)

	payload := map[string]interface{}{
		"user_id":       req.UserID.String(),
		"status":        account.Status.String(),
		"terminated_at": now.Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(req.UserID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.accountRepo.Save(ctx, account); err != nil {
		return err
	}

	c.logger.Infof("account terminated: user_id=%s", req.UserID.String())
	return nil
}