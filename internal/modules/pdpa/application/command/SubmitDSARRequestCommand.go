// SubmitDSARRequestCommand creates a new DSAR request secured by an OTP code.
// คำสั่งสร้างคำขอใช้สิทธิตาม พ.ร.บ. คุ้มครองข้อมูลส่วนบุคคลพร้อมรหัส OTP สำหรับยืนยันตัวตน
package command

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// SubmitDSARRequestCommand handles DSAR submission.
// จัดการการยื่นคำขอใช้สิทธิตามกฎหมายคุ้มครองข้อมูลส่วนบุคคล
type SubmitDSARRequestCommand struct {
	dsarRepo repository.DSARRepository
	producer *kafka.Producer
	cfg      *config.Config
	logger   logger.Logger
}

// NewSubmitDSARRequestCommand creates a new SubmitDSARRequestCommand.
// สร้างคำสั่งยื่นคำขอใช้สิทธิใหม่
func NewSubmitDSARRequestCommand(dsarRepo repository.DSARRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *SubmitDSARRequestCommand {
	return &SubmitDSARRequestCommand{
		dsarRepo: dsarRepo,
		producer: producer,
		cfg:      cfg,
		logger:   logger,
	}
}

// Execute persists a new DSAR request and publishes the dsar.submitted event.
// บันทึกคำขอใหม่และเผยแพร่เหตุการณ์ dsar.submitted
func (c *SubmitDSARRequestCommand) Execute(ctx context.Context, req *dto.SubmitDSARRequest) error {
	now := time.Now()
	otpCode := generateOTPCode()
	otpExpiry := now.Add(5 * time.Minute)

	dsar := entity.NewDSARRequest(req.UserID, req.RequestType, otpCode, otpExpiry, "", "")

	payload := map[string]interface{}{
		"dsar_request_id": dsar.ID.String(),
		"user_id":         dsar.UserID.String(),
		"request_type":    dsar.RequestType.String(),
		"status":          dsar.Status.String(),
		"requested_at":    dsar.RequestedAt.Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(dsar.ID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.dsarRepo.Save(ctx, dsar); err != nil {
		return err
	}

	c.logger.Infof("dsar submitted: user_id=%s request_type=%s", req.UserID.String(), req.RequestType.String())
	return nil
}

// generateOTPCode generates a 6-digit OTP code using a cryptographic random source.
// สร้างรหัส OTP 6 หลักจากแหล่งสุ่มเชิงเข้ารหัส
func generateOTPCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}