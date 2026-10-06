package usecase

import (
	"context"
	"time"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

const cacheTTLLog = 30 * time.Second

type logUseCase struct {
	repo   users.LogRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewLogUseCase(repo users.LogRepository, cache users.CacheRepository, l logger.Logger) users.LogUseCaseI {
	return &logUseCase{repo: repo, cache: cache, logger: l}
}

func (u *logUseCase) CreateLog(ctx context.Context, l *usermodels.SdUserLog) (*usermodels.SdUserLog, error) {
	if l.Create.IsZero() { l.Create = time.Now() }
	if l.Update.IsZero() { l.Update = time.Now() }
	out, err := u.repo.Create(ctx, l)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:log:list:*", 200)
	return out, nil
}

func (u *logUseCase) GetLog(ctx context.Context, id int, logTypeID int) (*usermodels.SdUserLog, error) {
	return u.repo.Get(ctx, id, logTypeID)
}

func (u *logUseCase) ListLogs(ctx context.Context, f users.LogFilter) ([]*usermodels.SdUserLog, int64, error) {
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Limit > 200 { f.Limit = 200 }
	if f.Offset < 0 { f.Offset = 0 }
	key := "sd_user:log:list:" + hashFilter(f)
	type cachedList struct {
		Items []*usermodels.SdUserLog `json:"items"`
		Total int64                   `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, f)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLLog)
	return items, total, nil
}

func (u *logUseCase) DeleteLog(ctx context.Context, id int, logTypeID int) error {
	if err := u.repo.Delete(ctx, id, logTypeID); err != nil { return err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:log:list:*", 200)
	return nil
}

func (u *logUseCase) WriteAuditLogWithMeta(
	ctx context.Context,
	uid uuid.UUID,
	logTypeID int,
	name, detail, ip, device, browser string,
) {
	if u == nil || u.repo == nil { return }
	now := time.Now()
	log := &usermodels.SdUserLog{
		LogTypeID:    logTypeID,
		UID:          uid,
		Name:         name,
		Detail:       detail,
		IPAddress:    strPtr(ip),
		DeviceAccess: strPtr(device),
		BrowserAgent: strPtr(browser),
		Create:       now,
		Update:       now,
	}
	if _, err := u.repo.Create(ctx, log); err != nil {
		u.logger.Warnf("WriteAuditLog failed (uid=%s, type=%d): %v", uid, logTypeID, err)
	}
}

var _ users.LogUseCaseI = (*logUseCase)(nil)