package command

import (
	"context"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/domain/entity"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type PublishPrivacyPolicyCommand struct {
	repo     domainrepo.PolicyRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewPublishPrivacyPolicyCommand(
	r domainrepo.PolicyRepository,
	p *kafka.Producer, cfg *config.Config, l logger.Logger,
) *PublishPrivacyPolicyCommand {
	return &PublishPrivacyPolicyCommand{repo: r, producer: p, topic: cfg.Kafka.Topic, log: l}
}

type PublishPolicyInput struct {
	PolicyKey   string
	Version     string
	TitleEn     string
	TitleTh     string
	ContentEn   *string
	ContentTh   *string
	EffectiveAt *time.Time
}

func (c *PublishPrivacyPolicyCommand) Execute(ctx context.Context, in PublishPolicyInput) (*entity.PdpaPolicy, error) {
	now := time.Now()
	p := entity.NewPdpaPolicy(
		in.PolicyKey, in.Version, in.TitleEn, in.TitleTh,
		in.ContentEn, in.ContentTh, 1, in.EffectiveAt, &now, nil,
	)
	if err := c.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
