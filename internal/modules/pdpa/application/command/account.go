package command

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	valueobject "icmongolang/internal/modules/pdpa/domain/valueobject"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type SuspendAccountCommand struct {
	repo     domainrepo.UserAccountStatusRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewSuspendAccountCommand(
	r domainrepo.UserAccountStatusRepository,
	p *kafka.Producer, cfg *config.Config, l logger.Logger,
) *SuspendAccountCommand {
	return &SuspendAccountCommand{repo: r, producer: p, topic: cfg.Kafka.Topic, log: l}
}

type SuspendAccountInput struct {
	UserID         uuid.UUID
	RetentionYears int
	Reason         string
}

func (c *SuspendAccountCommand) Execute(ctx context.Context, in SuspendAccountInput) error {
	s, err := c.repo.FindByUserID(ctx, in.UserID)
	if err != nil || s == nil {
		s = entity.NewUserAccountStatus(in.UserID, valueobject.AccountActive)
	}
	s.SetSuspended(time.Now(), in.RetentionYears)
	if err := c.repo.Save(ctx, s); err != nil {
		return err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(in.UserID.String(), mustJSON(map[string]interface{}{
			"event":   "account.suspended",
			"user_id": in.UserID.String(),
			"reason":  in.Reason,
		}))
	}
	return nil
}

type TerminateAccountCommand struct {
	repo     domainrepo.UserAccountStatusRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewTerminateAccountCommand(
	r domainrepo.UserAccountStatusRepository,
	p *kafka.Producer, cfg *config.Config, l logger.Logger,
) *TerminateAccountCommand {
	return &TerminateAccountCommand{repo: r, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *TerminateAccountCommand) Execute(ctx context.Context, userID uuid.UUID) error {
	s, err := c.repo.FindByUserID(ctx, userID)
	if err != nil || s == nil {
		s = entity.NewUserAccountStatus(userID, valueobject.AccountActive)
	}
	s.SetTerminated(time.Now())
	if err := c.repo.Save(ctx, s); err != nil {
		return err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(userID.String(), mustJSON(map[string]interface{}{
			"event":   "account.terminated",
			"user_id": userID.String(),
		}))
	}
	return nil
}

type ConfirmDeletionCommand struct {
	repo     domainrepo.UserAccountStatusRepository
	policy   *service.DeletionPolicyService
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewConfirmDeletionCommand(
	r domainrepo.UserAccountStatusRepository,
	p *service.DeletionPolicyService,
	pr *kafka.Producer, cfg *config.Config, l logger.Logger,
) *ConfirmDeletionCommand {
	return &ConfirmDeletionCommand{repo: r, policy: p, producer: pr, topic: cfg.Kafka.Topic, log: l}
}

func (c *ConfirmDeletionCommand) Execute(ctx context.Context, userID uuid.UUID) error {
	s, err := c.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	s.ConfirmDeletion()
	if err := c.repo.Save(ctx, s); err != nil {
		return err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(userID.String(), mustJSON(map[string]interface{}{
			"event":   "data.deletion_requested",
			"user_id": userID.String(),
		}))
	}
	return nil
}

type ImmediateDeletionCommand struct {
	accountRepo domainrepo.UserAccountStatusRepository
	consentRepo domainrepo.ConsentRepository
	dsarRepo    domainrepo.DSARRepository
	policy      *service.DeletionPolicyService
	producer    *kafka.Producer
	topic       string
	log         logger.Logger
}

func NewImmediateDeletionCommand(
	a domainrepo.UserAccountStatusRepository,
	c domainrepo.ConsentRepository,
	d domainrepo.DSARRepository,
	p *service.DeletionPolicyService,
	pr *kafka.Producer, cfg *config.Config, l logger.Logger,
) *ImmediateDeletionCommand {
	return &ImmediateDeletionCommand{
		accountRepo: a, consentRepo: c, dsarRepo: d,
		policy: p, producer: pr, topic: cfg.Kafka.Topic, log: l,
	}
}

func (c *ImmediateDeletionCommand) Execute(ctx context.Context, userID uuid.UUID) error {
	s, err := c.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if !c.policy.CanImmediateDeletion(s) {
		return errors.New("immediate deletion not allowed")
	}
	_ = c.consentRepo.DeleteByUserID(ctx, userID)
	_ = c.dsarRepo.DeleteByUserID(ctx, userID)
	s.MarkDeleted()
	return c.accountRepo.Save(ctx, s)
}

type AutoDeleteExpiredCommand struct {
	accountRepo domainrepo.UserAccountStatusRepository
	consentRepo domainrepo.ConsentRepository
	producer    *kafka.Producer
	topic       string
	log         logger.Logger
}

func NewAutoDeleteExpiredCommand(
	a domainrepo.UserAccountStatusRepository,
	c domainrepo.ConsentRepository,
	p *kafka.Producer, cfg *config.Config, l logger.Logger,
) *AutoDeleteExpiredCommand {
	return &AutoDeleteExpiredCommand{accountRepo: a, consentRepo: c, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *AutoDeleteExpiredCommand) Execute(ctx context.Context) (int, error) {
	rows, err := c.accountRepo.FindSuspendedWithoutDeletion(ctx, time.Now())
	if err != nil {
		return 0, err
	}
	n := 0
	for _, s := range rows {
		if err := c.accountRepo.UpdateAutoDeleted(ctx, s.UserID, time.Now()); err != nil {
			c.log.Warnf("auto-delete failed for user %s: %v", s.UserID, err)
			continue
		}
		_ = c.consentRepo.DeleteByUserID(ctx, s.UserID)
		n++
	}
	return n, nil
}
