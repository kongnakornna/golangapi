// ProcessDSARCommand transitions a DSAR request into the processing state.
// คำสั่งเปลี่ยนสถานะคำขอใช้สิทธิตาม พ.ร.บ. คุ้มครองข้อมูลส่วนบุคคลเป็นกำลังดำเนินการ
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
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// ProcessDSARCommand handles DSAR processing.
// จัดการการเริ่มดำเนินการตามคำขอใช้สิทธิ
type ProcessDSARCommand struct {
	dsarRepo repository.DSARRepository
	producer *kafka.Producer
	cfg      *config.Config
	logger   logger.Logger
}

// NewProcessDSARCommand creates a new ProcessDSARCommand.
// สร้างคำสั่งเริ่มดำเนินการตามคำขอใช้สิทธิใหม่
func NewProcessDSARCommand(dsarRepo repository.DSARRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *ProcessDSARCommand {
	return &ProcessDSARCommand{
		dsarRepo: dsarRepo,
		producer: producer,
		cfg:      cfg,
		logger:   logger,
	}
}

// Execute marks the DSAR request as processing and publishes the
// dsar.processing audit trail event.
// เปลี่ยนสถานะคำขอเป็นกำลังดำเนินการและเผยแพร่เหตุการณ์ dsar.processing
func (c *ProcessDSARCommand) Execute(ctx context.Context, req *dto.ProcessDSARRequest) error {
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

	if err := dsar.MarkProcessing(); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"action":          "dsar.processing",
		"dsar_request_id": dsar.ID.String(),
		"user_id":         dsar.UserID.String(),
		"status":          dsar.Status.String(),
		"processed_at":    time.Now().Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(dsar.ID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	c.logger.Infof("dsar processing: user_id=%s dsar_request_id=%s", req.UserID.String(), req.DSARRequestID.String())
	return nil
}