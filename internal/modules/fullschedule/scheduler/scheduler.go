package scheduler

import (
	"context"
	"fmt"
	"time"

	"icmongolang/internal/modules/fullschedule"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	// defaultInterval is how often due schedules are swept.
	defaultInterval = 30 * time.Second
	// lockTTL protects against a worker dying mid-run (stale lock auto-expires).
	lockTTL = 5 * time.Minute
)

// Scheduler sweeps due schedules and runs them, guarded by Redis locks so the
// same schedule can never run twice (idempotency).
// ตัวประมวลผล: ตรวจสอบ schedule ที่ถึงกำหนดและสั่งทำงาน กันการสั่งซ้ำด้วยล็อก Redis
type Scheduler struct {
	uc       fullschedule.FullScheduleUseCaseI
	rdb      *redis.Client
	enabled  bool
	interval time.Duration
	logger   logger.Logger
	cancel   context.CancelFunc
}

// New creates a scheduler. A nil redis client disables execution.
func New(uc fullschedule.FullScheduleUseCaseI, rdb *redis.Client, log logger.Logger, interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = defaultInterval
	}
	return &Scheduler{
		uc:       uc,
		rdb:      rdb,
		enabled:  rdb != nil,
		interval: interval,
		logger:   log,
	}
}

// Enabled reports whether scheduling is active.
func (s *Scheduler) Enabled() bool { return s.enabled }

// Start launches the scheduler loop in the background.
func (s *Scheduler) Start(ctx context.Context) {
	if !s.enabled {
		s.logger.Warn("⚠️ fullschedule: scheduler disabled (redis client is nil)")
		return
	}
	loopCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go s.loop(loopCtx)
	s.logger.Info("✅ fullschedule: scheduler started")
}

// Stop shuts the scheduler loop down.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Scheduler) loop(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	due, err := s.uc.GetDueSchedules(ctx, time.Now())
	if err != nil {
		s.logger.Errorf("fullschedule: get due schedules failed: %v", err)
		return
	}
	for _, sch := range due {
		if sch == nil {
			continue
		}
		if !s.tryLock(ctx, sch.ID) {
			continue // อีกเครื่อง/รอบก่อนเคยสั่งไปแล้ว → ข้าม (idempotency)
		}
		source := fullschedule.TriggerAutomatic
		if sch.Mode == fullschedule.ModeBatch && !sch.ManualTrigger {
			source = fullschedule.TriggerCron
		}
		if sch.ManualTrigger {
			source = fullschedule.TriggerManual
		}
		if err := s.uc.Run(ctx, sch, source, fullschedule.TriggeredSystem); err != nil {
			s.logger.Errorf("fullschedule: run %s failed: %v", sch.ID, err)
		}
		s.unlock(ctx, sch.ID)
	}
}

// tryLock atomically acquires the per-schedule Redis lock (NX + EX).
func (s *Scheduler) tryLock(ctx context.Context, id uuid.UUID) bool {
	ok, err := s.rdb.SetNX(ctx, lockKey(id), time.Now().UnixMilli(), lockTTL).Result()
	if err != nil {
		s.logger.Errorf("fullschedule: lock failed for %s: %v", id, err)
		return false
	}
	return ok
}

func (s *Scheduler) unlock(ctx context.Context, id uuid.UUID) {
	if err := s.rdb.Del(ctx, lockKey(id)).Err(); err != nil {
		s.logger.Errorf("fullschedule: unlock failed for %s: %v", id, err)
	}
}

func lockKey(id uuid.UUID) string {
	return fmt.Sprintf("fs:schedule:lock:%s", id.String())
}
