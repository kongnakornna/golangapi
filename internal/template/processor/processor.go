package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"icmongolang/config"
	"icmongolang/internal/processor"
	"icmongolang/internal/template"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/sendEmail"

	"github.com/hibiken/asynq"
)

type templateRedisTaskProcessor struct {
	processor.RedisTaskProcessor
	emailSender sendEmail.EmailSender
}

func NewTemplateRedisTaskProcessor(server *asynq.Server, cfg *config.Config, log logger.Logger, emailSender sendEmail.EmailSender) template.TaskProcessor {
	return &templateRedisTaskProcessor{
		RedisTaskProcessor: processor.NewRedisTaskProcessor(server, cfg, log),
		emailSender:        emailSender,
	}
}

func (p *templateRedisTaskProcessor) ProcessTaskSendNotification(ctx context.Context, task *asynq.Task) error {
	var payload template.PayloadSendNotification
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	p.Logger.Infof("Processing notification: %s", payload.Subject)

	if err := p.emailSender.SendEmail(ctx, p.Cfg.Email.From, payload.To, payload.Subject, payload.Body, ""); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	p.Logger.Infof("Notification sent: %s", task.Type())
	return nil
}
