// RevokeConsentCommand revokes a user's latest active consent for a given purpose.
// คำสั่งเพิกถอนความยินยอมล่าสุดของผู้ใช้สำหรับวัตถุประสงค์ที่กำหนด
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

// RevokeConsentCommand handles consent revocation.
// จัดการการเพิกถอนความยินยอม
type RevokeConsentCommand struct {
	consentRepo repository.ConsentRepository
	cache       repository.Cache
	producer    *kafka.Producer
	cfg         *config.Config
	logger      logger.Logger
}

// NewRevokeConsentCommand creates a new RevokeConsentCommand.
// สร้างคำสั่งเพิกถอนความยินยอมใหม่
func NewRevokeConsentCommand(consentRepo repository.ConsentRepository, cache repository.Cache, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *RevokeConsentCommand {
	return &RevokeConsentCommand{
		consentRepo: consentRepo,
		cache:       cache,
		producer:    producer,
		cfg:         cfg,
		logger:      logger,
	}
}

// Execute revokes the latest consent for the user and purpose, publishes the
// consent.revoked event, persists the change and refreshes the cache.
// เพิกถอนความยินยอมล่าสุด เผยแพร่เหตุการณ์ consent.revoked บันทึกข้อมูลและอัปเดตแคช
func (c *RevokeConsentCommand) Execute(ctx context.Context, req *dto.RevokeConsentRequest) error {
	consent, err := c.consentRepo.FindLatestByUserAndPurpose(ctx, req.UserID, req.Purpose)
	if err != nil {
		if errors.Is(err, domainerrors.ErrConsentNotFound) {
			return domainerrors.ErrConsentNotFound
		}
		return err
	}
	if consent == nil {
		return domainerrors.ErrConsentNotFound
	}

	if err := consent.Revoke(); err != nil {
		return err
	}

	revokedAt := time.Now()
	if consent.RevokedAt != nil {
		revokedAt = *consent.RevokedAt
	}

	payload := map[string]interface{}{
		"consent_id": consent.ID.String(),
		"user_id":    consent.UserID.String(),
		"purpose":    consent.Purpose,
		"status":     consent.Status.String(),
		"revoked_at": revokedAt.Format(time.RFC3339),
	}
	if err := c.producer.PublishMessage(consent.UserID.String(), payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.consentRepo.Save(ctx, consent); err != nil {
		return err
	}

	if err := c.cache.Set(ctx, consent.UserID, req.Purpose, consent.Status); err != nil {
		return err
	}

	c.logger.Infof("consent revoked: user_id=%s purpose=%s", consent.UserID.String(), req.Purpose.String())
	return nil
}