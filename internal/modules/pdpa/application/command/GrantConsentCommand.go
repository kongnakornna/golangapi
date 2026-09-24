package command

import (
	"context"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/dto"
	"icmongolang/internal/modules/pdpa/domain/entity"
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// GrantConsentCommand ใช้สำหรับให้สิทธิ์การยินยอม (consent) แก่ผู้ใช้งานตามวัตถุประสงค์ที่กำหนด.
// GrantConsentCommand grants consent to a user for the specified purpose.
type GrantConsentCommand struct {
	consentRepo repository.ConsentRepository
	cache       repository.Cache
	producer    *kafka.Producer
	cfg         *config.Config
	logger      logger.Logger
}

// NewGrantConsentCommand สร้าง GrantConsentCommand พร้อม dependency ที่จำเป็นทั้งหมด.
// NewGrantConsentCommand creates a GrantConsentCommand with all required dependencies.
func NewGrantConsentCommand(consentRepo repository.ConsentRepository, cache repository.Cache, producer *kafka.Producer, cfg *config.Config, logger logger.Logger) *GrantConsentCommand {
	return &GrantConsentCommand{
		consentRepo: consentRepo,
		cache:       cache,
		producer:    producer,
		cfg:         cfg,
		logger:      logger,
	}
}

// Execute บันทึก consent log, ประกาศเหตุการณ์ consent.granted แล้วอัปเดตสถานะในแคช.
// Execute stores the consent log, publishes the consent.granted event, and updates the cache.
func (c *GrantConsentCommand) Execute(ctx context.Context, req dto.GrantConsentRequest) error {
	now := time.Now()
	consent := entity.NewConsentLog(req.UserID, req.Purpose.String(), now)

	if err := c.producer.PublishMessage(consent.UserID.String(), map[string]interface{}{
		"consent_id":   consent.ID.String(),
		"user_id":      consent.UserID.String(),
		"purpose":      consent.Purpose,
		"status":       consent.Status.String(),
		"consented_at": consent.ConsentedAt.Format(time.RFC3339),
		"expires_at":   consent.ExpiresAt.Format(time.RFC3339),
	}); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	if err := c.consentRepo.Save(ctx, consent); err != nil {
		c.logger.Errorf("Failed to save consent log: %v", err)
		return err
	}

	if err := c.cache.Set(ctx, consent.UserID, req.Purpose, consent.Status); err != nil {
		c.logger.Errorf("Failed to cache consent status: %v", err)
		return err
	}

	c.logger.Infof("Consent granted for user %s purpose %s", consent.UserID.String(), consent.Purpose)
	return nil
}