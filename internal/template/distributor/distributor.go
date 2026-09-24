package distributor

import (
	"context"
	"encoding/json"
	"fmt"

	"icmongolang/config"
	"icmongolang/internal/distributor"
	"icmongolang/internal/template"
	"icmongolang/pkg/logger"

	"github.com/hibiken/asynq"
)

type templateRedisTaskDistributor struct {
	distributor.RedisTaskDistributor
}

func NewTemplateRedisTaskDistributor(redisClient *asynq.Client, cfg *config.Config, log logger.Logger) template.TaskDistributor {
	return &templateRedisTaskDistributor{
		RedisTaskDistributor: distributor.NewRedisTaskDistributor(redisClient, cfg, log),
	}
}

func (d *templateRedisTaskDistributor) DistributeTaskSendNotification(ctx context.Context, payload *template.PayloadSendNotification, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(template.TaskSendNotification, jsonPayload, opts...)
	info, err := d.RedisClient.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	d.Logger.Infof("Type: %v, Queue: %v, MaxRetry: %v, Msg: queued", task.Type(), info.Queue, info.MaxRetry)
	return nil
}
