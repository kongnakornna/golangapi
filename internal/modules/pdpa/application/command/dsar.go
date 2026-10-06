package command

import (
	"context"
	"crypto/rand"
	"fmt"
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

type SubmitDSARRequestCommand struct {
	dsarRepo domainrepo.DSARRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewSubmitDSARRequestCommand(
	d domainrepo.DSARRepository, p *kafka.Producer,
	cfg *config.Config, l logger.Logger,
) *SubmitDSARRequestCommand {
	return &SubmitDSARRequestCommand{dsarRepo: d, producer: p, topic: cfg.Kafka.Topic, log: l}
}

type SubmitDSARInput struct {
	UserID      uuid.UUID
	RequestType valueobject.DSARType
	IP          string
	UA          string
}

func (c *SubmitDSARRequestCommand) Execute(ctx context.Context, in SubmitDSARInput) (*entity.DSARRequest, error) {
	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}
	expiry := time.Now().Add(15 * time.Minute)
	d := entity.NewDSARRequest(in.UserID, in.RequestType, otp, expiry, in.IP, in.UA)
	if err := c.dsarRepo.Save(ctx, d); err != nil {
		return nil, err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(in.UserID.String(), mustJSON(map[string]interface{}{
			"event":      "dsar.submitted",
			"user_id":    in.UserID.String(),
			"request_id": d.ID.String(),
			"type":       in.RequestType.String(),
		}))
	}
	return d, nil
}

type VerifyDSAROTPCommand struct {
	dsarRepo domainrepo.DSARRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewVerifyDSAROTPCommand(d domainrepo.DSARRepository, p *kafka.Producer, cfg *config.Config, l logger.Logger) *VerifyDSAROTPCommand {
	return &VerifyDSAROTPCommand{dsarRepo: d, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *VerifyDSAROTPCommand) Execute(ctx context.Context, requestID uuid.UUID, otp string) error {
	d, err := c.dsarRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if !d.VerifyOTP(otp) {
		return fmt.Errorf("OTP expired or invalid")
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(d.UserID.String(), mustJSON(map[string]interface{}{
			"event":      "dsar.otp_verified",
			"user_id":    d.UserID.String(),
			"request_id": d.ID.String(),
		}))
	}
	return nil
}

type ProcessDSARCommand struct {
	dsarRepo domainrepo.DSARRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewProcessDSARCommand(d domainrepo.DSARRepository, p *kafka.Producer, cfg *config.Config, l logger.Logger) *ProcessDSARCommand {
	return &ProcessDSARCommand{dsarRepo: d, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *ProcessDSARCommand) Execute(ctx context.Context, requestID uuid.UUID) error {
	d, err := c.dsarRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if err := d.MarkProcessing(); err != nil {
		return err
	}
	if err := c.dsarRepo.Update(ctx, d); err != nil {
		return err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(d.UserID.String(), mustJSON(map[string]interface{}{
			"event":      "dsar.processing",
			"user_id":    d.UserID.String(),
			"request_id": d.ID.String(),
		}))
	}
	return nil
}

type CompleteDSARCommand struct {
	dsarRepo domainrepo.DSARRepository
	outbox   domainrepo.OutboxRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewCompleteDSARCommand(
	d domainrepo.DSARRepository,
	o domainrepo.OutboxRepository,
	p *kafka.Producer, cfg *config.Config, l logger.Logger,
) *CompleteDSARCommand {
	return &CompleteDSARCommand{dsarRepo: d, outbox: o, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *CompleteDSARCommand) Execute(ctx context.Context, requestID uuid.UUID, payload map[string]interface{}) error {
	d, err := c.dsarRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	body := mustJSON(payload)
	if err := d.MarkCompleted(body); err != nil {
		return err
	}
	if err := c.dsarRepo.Update(ctx, d); err != nil {
		return err
	}
	evt := eventd.OutboxEvent{
		ID:    uuid.New(),
		Topic: c.topic,
		Key:   d.UserID.String(),
		Payload: mustJSON(map[string]interface{}{
			"event":      "dsar.completed",
			"user_id":    d.UserID.String(),
			"request_id": d.ID.String(),
		}),
		Status:    eventd.OutboxEventStatusPending,
		CreatedAt: time.Now(),
	}
	if c.outbox != nil {
		_ = c.outbox.Publish(ctx, evt)
	} else if c.producer != nil {
		_ = c.producer.PublishMessage(evt.Key, evt.Payload)
	}
	return nil
}

type RejectDSARCommand struct {
	dsarRepo domainrepo.DSARRepository
	producer *kafka.Producer
	topic    string
	log      logger.Logger
}

func NewRejectDSARCommand(d domainrepo.DSARRepository, p *kafka.Producer, cfg *config.Config, l logger.Logger) *RejectDSARCommand {
	return &RejectDSARCommand{dsarRepo: d, producer: p, topic: cfg.Kafka.Topic, log: l}
}

func (c *RejectDSARCommand) Execute(ctx context.Context, requestID uuid.UUID, reason string) error {
	d, err := c.dsarRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if err := d.MarkRejected(reason); err != nil {
		return err
	}
	if err := c.dsarRepo.Update(ctx, d); err != nil {
		return err
	}
	if c.producer != nil {
		_ = c.producer.PublishMessage(d.UserID.String(), mustJSON(map[string]interface{}{
			"event":      "dsar.rejected",
			"user_id":    d.UserID.String(),
			"request_id": d.ID.String(),
			"reason":     reason,
		}))
	}
	return nil
}

func generateOTP() (string, error) {
	const digits = "0123456789"
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, 6)
	for i, v := range b {
		out[i] = digits[int(v)%len(digits)]
	}
	return fmt.Sprintf("%s", out), nil
}
