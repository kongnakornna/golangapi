// Package scheduler implements background jobs ของ PDPA module
package scheduler

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"time"

	"icmongolang/internal/modules/pdpa/application/port"
	"icmongolang/internal/modules/pdpa/domain/event"
)

const outboxHandler = "pdpa-outbox"

const (
	defaultPollInterval  = 5 * time.Second
	defaultBatchLimit    = 100
	defaultMaxDeliveries = 3
	defaultBackoffBase   = 500 * time.Millisecond
	defaultBackoffMax    = 30 * time.Second
)

type EventPublisher interface {
	PublishMessage(key string, msg interface{}) error
}

type Dispatcher struct {
	outboxRepo    port.OutboxRepository
	producer      EventPublisher
	clock         port.Clock
	pollInterval  time.Duration
	batchLimit    int
	maxDeliveries int
	backoffBase   time.Duration
	backoffMax    time.Duration
	logger        *slog.Logger
}

func NewDispatcher(
	outboxRepo port.OutboxRepository,
	producer EventPublisher,
	clock port.Clock,
	pollInterval time.Duration,
	opts ...Option,
) *Dispatcher {
	d := &Dispatcher{
		outboxRepo:    outboxRepo,
		producer:      producer,
		clock:         clock,
		pollInterval:  pollInterval,
		batchLimit:    defaultBatchLimit,
		maxDeliveries: defaultMaxDeliveries,
		backoffBase:   defaultBackoffBase,
		backoffMax:    defaultBackoffMax,
		logger:        slog.Default(),
	}
	for _, opt := range opts {
		opt(d)
	}
	if d.pollInterval <= 0 {
		d.pollInterval = defaultPollInterval
	}
	return d
}

type Option func(*Dispatcher)

func WithBatchLimit(limit int) Option {
	return func(d *Dispatcher) { d.batchLimit = limit }
}

func WithMaxDeliveries(n int) Option {
	return func(d *Dispatcher) { d.maxDeliveries = n }
}

func WithBackoff(base, max time.Duration) Option {
	return func(d *Dispatcher) { d.backoffBase, d.backoffMax = base, max }
}

func WithLogger(l *slog.Logger) Option {
	return func(d *Dispatcher) { d.logger = l }
}

func (d *Dispatcher) Start(ctx context.Context) {
	ticker := time.NewTicker(d.pollInterval)
	defer ticker.Stop()
	d.RunOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.RunOnce(ctx)
		}
	}
}

func (d *Dispatcher) RunOnce(ctx context.Context) {
	events, err := d.outboxRepo.ClaimPending(ctx, d.batchLimit, outboxHandler)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		d.logger.Error("claim pending outbox events", "error", err)
		return
	}
	for _, evt := range events {
		d.dispatch(ctx, evt)
	}
}

func (d *Dispatcher) dispatch(ctx context.Context, evt event.OutboxEvent) {
	if err := d.producer.PublishMessage(evt.Key, json.RawMessage(evt.Payload)); err != nil {
		d.onPublishFailed(ctx, evt, err)
		return
	}
	if err := d.outboxRepo.MarkSuccess(ctx, evt.ID, outboxHandler); err != nil {
		if ctx.Err() != nil {
			return
		}
		d.logger.Error("mark outbox event success", "event_id", evt.ID, "error", err)
	}
}

func (d *Dispatcher) onPublishFailed(ctx context.Context, evt event.OutboxEvent, pubErr error) {
	d.logger.Warn("publish outbox event failed",
		"event_id", evt.ID, "topic", evt.Topic,
		"attempts", evt.Attempts, "error", pubErr)
	if evt.Attempts >= d.maxDeliveries {
		if err := d.outboxRepo.MarkFailed(ctx, evt.ID, outboxHandler, evt.Attempts+1); err != nil {
			if ctx.Err() != nil {
				return
			}
			d.logger.Error("mark outbox event failed", "event_id", evt.ID, "error", err)
		}
		return
	}
	delay := d.backoff(evt.Attempts)
	if delay <= 0 {
		return
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func (d *Dispatcher) backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	exp := 1 << min(attempt, 30)
	base := d.backoffBase * time.Duration(exp)
	if base > d.backoffMax {
		base = d.backoffMax
	}
	return time.Duration(rand.Int64N(int64(base)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
