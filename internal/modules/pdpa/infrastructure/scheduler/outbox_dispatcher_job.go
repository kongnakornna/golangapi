// Package scheduler implements the background jobs of the pdpa module.
//
// The outbox dispatcher drains the pdpa outbox table: it polls for pending
// events, publishes each one to the configured Kafka topic, and records
// success or failure in the outbox store so that delivery stays at-least-once
// and recoverable across restarts.
package scheduler

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"icmongolang/internal/modules/pdpa/application/port"
	"icmongolang/internal/modules/pdpa/domain/event"
)

// outboxHandler is the handler identity this dispatcher records in the outbox
// store. It scopes ClaimPending and MarkSuccess/MarkFailed calls so that a
// batch claimed by another consumer is never finalized twice.
const outboxHandler = "pdpa-outbox"

// Default dispatcher tuning applied unless overridden by an option.
const (
	defaultPollInterval  = 5 * time.Second
	defaultBatchLimit    = 100
	defaultMaxDeliveries = 3
	defaultBackoffBase   = 500 * time.Millisecond
	defaultBackoffMax    = 30 * time.Second
)

// EventPublisher publishes a single outbox event payload. It is satisfied by
// *kafka.Producer, which is configured with the destination topic at
// construction.
type EventPublisher interface {
	PublishMessage(key string, msg interface{}) error
}

// Dispatcher polls the pdpa outbox store and publishes pending events to the
// configured Kafka topic, completing or failing each event as the publish
// result requires.
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

// NewDispatcher builds an outbox dispatcher that drains outboxRepo and
// publishes claimed events through producer. pollInterval sets how often the
// store is polled; values <= 0 fall back to the default. Clock backs
// time-dependent decisions (currently unused by core dispatch, but kept so the
// dispatcher can grow TTL-aware reclaim without changing its signature).
func NewDispatcher(outboxRepo port.OutboxRepository, producer EventPublisher, clock port.Clock, pollInterval time.Duration, opts ...Option) *Dispatcher {
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

// Option customizes a Dispatcher.
type Option func(*Dispatcher)

// WithBatchLimit caps how many events a single poll cycle claims.
func WithBatchLimit(limit int) Option {
	return func(d *Dispatcher) { d.batchLimit = limit }
}

// WithMaxDeliveries sets how many delivery attempts an event gets before it is
// marked failed and left for dead-letter handling.
func WithMaxDeliveries(n int) Option {
	return func(d *Dispatcher) { d.maxDeliveries = n }
}

// WithBackoff sets the base and maximum retry backoff durations applied after
// a failed publish.
func WithBackoff(base, max time.Duration) Option {
	return func(d *Dispatcher) {
		d.backoffBase = base
		d.backoffMax = max
	}
}

// WithLogger sets the logger used for dispatch activity.
func WithLogger(logger *slog.Logger) Option {
	return func(d *Dispatcher) { d.logger = logger }
}

// Start runs the poll loop until ctx is cancelled. It returns after the
// dispatcher has stopped. RunOnce is available for a single drain in tests or
// on-demand operation.
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

// RunOnce claims and dispatches a single batch of pending events. It is safe
// to call concurrently with itself: the store's atomic claim prevents two
// workers from claiming the same event concurrently.
func (d *Dispatcher) RunOnce(ctx context.Context) {
	evts, err := d.outboxRepo.ClaimPending(ctx, d.batchLimit, outboxHandler)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		d.logger.Error("claim pending outbox events", "error", err)
		return
	}
	if len(evts) > 0 {
		d.logger.Info("dispatch outbox events", "count", len(evts))
	}
	for _, evt := range evts {
		d.dispatch(ctx, evt)
	}
}

// dispatch publishes one event and finalizes its outbox record.
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

// onPublishFailed applies the retry policy. A failed publish is backed off so
// the store can re-claim the event on a later cycle; once maxDeliveries is
// reached the event is marked failed and left for a dead-letter consumer (see
// the module plan, "Automated remediation and dead letters (P2)").
func (d *Dispatcher) onPublishFailed(ctx context.Context, evt event.OutboxEvent, pubErr error) {
	d.logger.Warn("publish outbox event failed",
		"event_id", evt.ID,
		"topic", evt.Topic,
		"attempts", evt.Attempts,
		"error", pubErr,
	)

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

// backoff returns a jittered, exponentially growing delay for the given
// zero-based attempt, capped at backoffMax. Jitter stops concurrent retriers
// from resurfacing in lockstep.
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

var _ = uuid.Nil