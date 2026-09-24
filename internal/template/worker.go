package template

import (
	"context"

	"github.com/hibiken/asynq"
)

const TaskSendNotification = "template:send_notification"

type PayloadSendNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type TaskDistributor interface {
	DistributeTaskSendNotification(ctx context.Context, payload *PayloadSendNotification, opts ...asynq.Option) error
}

type TaskProcessor interface {
	ProcessTaskSendNotification(ctx context.Context, task *asynq.Task) error
}
