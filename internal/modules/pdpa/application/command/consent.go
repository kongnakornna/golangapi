package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/domain/entity"
	eventd "icmongolang/internal/modules/pdpa/domain/event"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

// ---------- GrantConsentCommand ----------

type GrantConsentCommand struct {
	consentRepo domainrepo.ConsentRepository
	cache       domainrepo.Cache
	producer    *kafka.Producer
	log         logger.Logger
}

func NewGrantConsentCommand(
	consentRepo domainrepo.ConsentRepository,
	cache domainrepo.Cache,
	producer *kafka.Producer,
	cfg *config.Config,
	log logger.Logger,
) *GrantConsentCommand {
	return &GrantConsentCommand{
		consentRepo: consentRepo, cache: cache, producer: producer,
		log: log,
	}
}

type GrantConsentInput struct {
	UserID  uuid.UUID
	Purpose valueobject.ConsentPurpose
	Source  string
	IP      string
	UA      string
}

func (c *GrantConsentCommand) Execute(ctx context.Context, in GrantConsentInput) (*entity.ConsentLog, error) {
	now := time.Now()
	cl := entity.NewConsentLog(in.UserID, in.Purpose.String(), now)
	if err := c.consentRepo.Save(ctx, cl); err != nil {
		return nil, err
	}
	if c.cache != nil {
		_ = c.cache.Set(ctx, in.UserID, in.Purpose, valueobject.ConsentGranted)
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(in.UserID.String(), mustJSON(map[string]interface{}{
			"event":      "consent.granted",
			"user_id":    in.UserID.String(),
			"purpose":    in.Purpose.String(),
			"granted_at": now.Format(time.RFC3339),
			"source":     in.Source,
		}))
	}
	_ = eventd.OutboxEventStatusPending
	return cl, nil
}

// ---------- RevokeConsentCommand ----------

type RevokeConsentCommand struct {
	consentRepo domainrepo.ConsentRepository
	cache       domainrepo.Cache
	producer    *kafka.Producer
	log         logger.Logger
}

func NewRevokeConsentCommand(
	consentRepo domainrepo.ConsentRepository,
	cache domainrepo.Cache,
	producer *kafka.Producer,
	cfg *config.Config,
	log logger.Logger,
) *RevokeConsentCommand {
	return &RevokeConsentCommand{
		consentRepo: consentRepo, cache: cache, producer: producer,
		log: log,
	}
}

type RevokeConsentInput struct {
	UserID  uuid.UUID
	Purpose valueobject.ConsentPurpose
}

func (c *RevokeConsentCommand) Execute(ctx context.Context, in RevokeConsentInput) error {
	if err := c.consentRepo.Revoke(ctx, in.UserID, in.Purpose); err != nil {
		return err
	}
	if c.cache != nil {
		_ = c.cache.DeleteAllByUser(ctx, in.UserID)
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(in.UserID.String(), mustJSON(map[string]interface{}{
			"event":   "consent.revoked",
			"user_id": in.UserID.String(),
			"purpose": in.Purpose.String(),
		}))
	}
	return nil
}
