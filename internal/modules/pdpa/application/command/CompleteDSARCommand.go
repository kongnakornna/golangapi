// CompleteDSARCommand completes a DSAR request and records the response data hash
// on the blockchain for non-repudiation.
// คำสั่งเสร็จสิ้นการดำเนินการตามคำขอใช้สิทธิ พร้อมบันทึกแฮชข้อมูลลงบล็อกเชน
package command

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// CompleteDSARCommand handles DSAR completion with a blockchain audit record.
// จัดการการเสร็จสิ้นคำขอใช้สิทธิพร้อมบันทึกการตรวจสอบลงบล็อกเชน
type CompleteDSARCommand struct {
	dsarRepo      repository.DSARRepository
	blockchainSvc service.BlockchainService
	producer      *kafka.Producer
	cfg           *config.Config
	logger        logger.Logger
}

// NewCompleteDSARCommand creates a new CompleteDSARCommand.
// สร้างคำสั่งเสร็จสิ้นคำขอใช้สิทธิใหม่
func NewCompleteDSARCommand(dsarRepo repository.DSARRepository, blockchainSvc service.BlockchainService, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *CompleteDSARCommand {
	return &CompleteDSARCommand{
		dsarRepo:      dsarRepo,
		blockchainSvc: blockchainSvc,
		producer:      producer,
		cfg:           cfg,
		logger:        logger,
	}
}

// Execute finalizes the DSAR request, publishes the dsar.completed event and
// persists the completed status.
// ทำให้คำขอใช้สิทธิเสร็จสิ้น เผยแพร่เหตุการณ์ dsar.completed และบันทึกสถานะเสร็จสิ้น
func (c *CompleteDSARCommand) Execute(ctx context.Context, req *dto.CompleteDSARRequest) error {
	dsar, err := c.dsarRepo.FindByID(ctx, req.DSARRequestID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrDSARNotFound) {
			return domainerrors.ErrDSARNotFound
		}
		return err
	}
	if dsar == nil {
		return domainerrors.ErrDSARNotFound
	}

	dataHash := fmt.Sprintf("%x", sha256.Sum256(req.ResponseData))

	txHash := ""
	if hash, recErr := c.blockchainSvc.RecordHash(ctx, dataHash); recErr != nil {
		c.logger.Warnf("blockchain record failed: %v", recErr)
	} else {
		txHash = hash
	}

	payload := map[string]interface{}{
		"status":          valueobject.DSARStatusCompleted.String(),
		"completed_at":    time.Now().Format(time.RFC3339),
		"data_hash":       dataHash,
		"tx_hash":         txHash,
		"dsar_request_id": dsar.ID.String(),
		"user_id":         dsar.UserID.String(),
	}
	if err := c.producer.PublishMessage(dsar.ID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.dsarRepo.UpdateStatus(ctx, dsar.ID, valueobject.DSARStatusCompleted); err != nil {
		return err
	}

	c.logger.Infof("dsar completed: user_id=%s dsar_request_id=%s", req.UserID.String(), req.DSARRequestID.String())
	return nil
}