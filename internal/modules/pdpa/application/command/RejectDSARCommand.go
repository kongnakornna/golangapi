// RejectDSARCommand rejects a DSAR request and records the rejection reason.
// คำสั่งปฏิเสธคำขอใช้สิทธิตาม พ.ร.บ. คุ้มครองข้อมูลส่วนบุคคลพร้อมบันทึกเหตุผล
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
	"icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// RejectDSARCommand handles DSAR rejection.
// จัดการการปฏิเสธคำขอใช้สิทธิ
type RejectDSARCommand struct {
	dsarRepo repository.DSARRepository
	producer *kafka.Producer
	cfg      *config.Config
	logger   logger.Logger
}

// NewRejectDSARCommand creates a new RejectDSARCommand.
// สร้างคำสั่งปฏิเสธคำขอใช้สิทธิใหม่
func NewRejectDSARCommand(dsarRepo repository.DSARRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *RejectDSARCommand {
	return &RejectDSARCommand{
		dsarRepo: dsarRepo,
		producer: producer,
		cfg:      cfg,
		logger:   logger,
	}
}

// Execute records the rejection reason, publishes the dsar.rejected audit trail
// event and persists the rejected status.
// บันทึกเหตุผลการปฏิเสธ เผยแพร่เหตุการณ์ dsar.rejected และบันทึกสถานะปฏิเสธ
func (c *RejectDSARCommand) Execute(ctx context.Context, req *dto.RejectDSARRequest) error {
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

	dsar.RejectionReason = req.Reason

	payload := map[string]interface{}{
		"action":          "dsar.rejected",
		"dsar_request_id": dsar.ID.String(),
		"user_id":         dsar.UserID.String(),
		"reason":          req.Reason,
		"status":          valueobject.DSARStatusRejected.String(),
		"rejected_at":     time.Now().Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(dsar.ID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.dsarRepo.UpdateStatus(ctx, dsar.ID, valueobject.DSARStatusRejected); err != nil {
		return err
	}

	c.logger.Infof("dsar rejected: user_id=%s dsar_request_id=%s", req.UserID.String(), req.DSARRequestID.String())
	return nil
}