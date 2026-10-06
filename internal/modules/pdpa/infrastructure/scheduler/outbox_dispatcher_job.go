package scheduler

import (
	"context"
	"time"

	"icmongolang/internal/modules/pdpa/application/port"
)

// OutboxDispatcherJob — wrapper ให้ตรงกับ API ที่ module.go ใช้
type OutboxDispatcherJob struct {
	inner *Dispatcher
}

func NewOutboxDispatcherJob(
	repo port.OutboxRepository,
	publisher EventPublisher,
	clock port.Clock,
	interval time.Duration,
	opts ...Option,
) *OutboxDispatcherJob {
	return &OutboxDispatcherJob{
		inner: NewDispatcher(repo, publisher, clock, interval, opts...),
	}
}

func (j *OutboxDispatcherJob) Start(ctx context.Context) {
	j.inner.Start(ctx)
}

func (j *OutboxDispatcherJob) RunOnce(ctx context.Context) {
	j.inner.RunOnce(ctx)
}
