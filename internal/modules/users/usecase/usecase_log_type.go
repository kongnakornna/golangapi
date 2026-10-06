package usecase

import (
	"context"
	"fmt"
	"time"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/pkg/logger"
)

const (
	cacheKeyLogTypeList   = "sd_user:log_type:list:%s"
	cacheKeyLogTypeDetail = "sd_user:log_type:%d"
	cacheTTLLogType       = 10 * time.Minute
)

type logTypeUseCase struct {
	repo   users.LogTypeRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewLogTypeUseCase(repo users.LogTypeRepository, cache users.CacheRepository, l logger.Logger) users.LogTypeUseCaseI {
	return &logTypeUseCase{repo: repo, cache: cache, logger: l}
}

func (u *logTypeUseCase) CreateLogType(ctx context.Context, lt *usermodels.SdUserLogType) (*usermodels.SdUserLogType, error) {
	out, err := u.repo.Create(ctx, lt)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:log_type:list:*", 200)
	return out, nil
}

func (u *logTypeUseCase) GetLogType(ctx context.Context, id int) (*usermodels.SdUserLogType, error) {
	key := fmt.Sprintf(cacheKeyLogTypeDetail, id)
	var c usermodels.SdUserLogType
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return &c, nil }
	lt, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	_ = u.cache.SetJSON(ctx, key, lt, cacheTTLLogType)
	return lt, nil
}

func (u *logTypeUseCase) UpdateLogType(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserLogType, error) {
	lt, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	out, err := u.repo.Update(ctx, lt, values)
	if err != nil { return nil, err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyLogTypeDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:log_type:list:*", 200)
	return out, nil
}

func (u *logTypeUseCase) DeleteLogType(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil { return err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyLogTypeDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:log_type:list:*", 200)
	return nil
}

func (u *logTypeUseCase) ListLogTypes(ctx context.Context, limit, offset int) ([]*usermodels.SdUserLogType, int64, error) {
	if limit <= 0 { limit = 20 }
	if limit > 100 { limit = 100 }
	if offset < 0 { offset = 0 }
	key := fmt.Sprintf(cacheKeyLogTypeList, hashFilter(map[string]int{"limit": limit, "offset": offset}))
	type cachedList struct {
		Items []*usermodels.SdUserLogType `json:"items"`
		Total int64                       `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, limit, offset)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLLogType)
	return items, total, nil
}

var _ users.LogTypeUseCaseI = (*logTypeUseCase)(nil)