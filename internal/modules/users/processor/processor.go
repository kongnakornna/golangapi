package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"icmongolang/config"
	"icmongolang/internal/processor"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/sendEmail"

	"github.com/hibiken/asynq"
)

type userRedisTaskProcessor struct {
	processor.RedisTaskProcessor
	emailSender sendEmail.EmailSender
}

func NewUserRedisTaskProcessor(server *asynq.Server, cfg *config.Config, logger logger.Logger, emailSender sendEmail.EmailSender) users.UserRedisTaskProcessor {
	return &userRedisTaskProcessor{
		RedisTaskProcessor: processor.NewRedisTaskProcessor(server, cfg, logger),
		emailSender:        emailSender,
	}
}

func (p *userRedisTaskProcessor) ProcessTaskSendEmail(ctx context.Context, task *asynq.Task) error {
	var payload users.PayloadSendEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	if err := p.emailSender.SendEmail(ctx, payload.From, payload.To, payload.Subject, payload.BodyHtml, payload.BodyPlain); err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	p.Logger.Infof("Type: %v, Msg: email sended", task.Type())
	return nil
}
