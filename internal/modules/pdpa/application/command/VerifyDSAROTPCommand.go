// VerifyDSAROTPCommand verifies the OTP code of a submitted DSAR request with rate limiting.
// คำสั่งตรวจสอบรหัส OTP ของคำขอใช้สิทธิตาม พ.ร.บ. คุ้มครองข้อมูลส่วนบุคคลพร้อมการจำกัดความถี่
package command

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	domainerrors "icmongolang/internal/modules/pdpa/domain/errors"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

const maxOTPAttempts = 5

// VerifyDSAROTPCommand handles OTP verification for DSAR requests.
// จัดการการตรวจสอบรหัส OTP สำหรับคำขอใช้สิทธิ
type VerifyDSAROTPCommand struct {
	dsarRepo repository.DSARRepository
	producer *kafka.Producer
	cfg      *config.Config
	logger   logger.Logger
	attempts sync.Map
}

// NewVerifyDSAROTPCommand creates a new VerifyDSAROTPCommand.
// สร้างคำสั่งตรวจสอบรหัส OTP ใหม่
func NewVerifyDSAROTPCommand(dsarRepo repository.DSARRepository, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *VerifyDSAROTPCommand {
	return &VerifyDSAROTPCommand{
		dsarRepo: dsarRepo,
		producer: producer,
		cfg:      cfg,
		logger:   logger,
		attempts: sync.Map{},
	}
}

// Execute verifies the OTP, enforcing a five-attempt limit per DSAR request and
// publishing the dsar.otp_verified audit trail event on success.
// ตรวจสอบรหัส OTP จำกัดจำนวนครั้งไม่เกินห้าครั้งต่อคำขอ และเผยแพร่เหตุการณ์ dsar.otp_verified เมื่อสำเร็จ
func (c *VerifyDSAROTPCommand) Execute(ctx context.Context, req *dto.VerifyOTPRequest) error {
	key := req.UserID.String() + ":" + req.DSARRequestID.String()

	count, loaded := c.attempts.LoadOrStore(key, 1)
	if loaded {
		attempts, ok := count.(int)
		if !ok {
			attempts = 0
		}
		attempts++
		if attempts > maxOTPAttempts {
			return domainerrors.ErrRateLimitExceeded
		}
		c.attempts.Store(key, attempts)
	}

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

	if !dsar.VerifyOTP(req.OTPCode) {
		return domainerrors.ErrOTPExpired
	}

	payload := map[string]interface{}{
		"action":          "dsar.otp_verified",
		"dsar_request_id": dsar.ID.String(),
		"user_id":         dsar.UserID.String(),
		"verified_at":     time.Now().Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(dsar.ID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	c.attempts.Delete(key)
	c.logger.Infof("dsar otp verified: user_id=%s dsar_request_id=%s", req.UserID.String(), req.DSARRequestID.String())
	return nil
}