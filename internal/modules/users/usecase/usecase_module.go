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
	cacheKeyModuleList   = "sd_user:module:list:%s"
	cacheKeyModuleDetail = "sd_user:module:%d"
	cacheTTLModule       = 10 * time.Minute
)

type moduleUseCase struct {
	repo   users.ModuleRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewModuleUseCase(repo users.ModuleRepository, cache users.CacheRepository, l logger.Logger) users.ModuleUseCaseI {
	return &moduleUseCase{repo: repo, cache: cache, logger: l}
}

func (u *moduleUseCase) CreateModule(ctx context.Context, m *usermodels.SdModules) (*usermodels.SdModules, error) {
	out, err := u.repo.Create(ctx, m)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:module:list:*", 200)
	return out, nil
}

func (u *moduleUseCase) GetModule(ctx context.Context, id int) (*usermodels.SdModules, error) {
	key := fmt.Sprintf(cacheKeyModuleDetail, id)
	var c usermodels.SdModules
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return &c, nil }
	m, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	_ = u.cache.SetJSON(ctx, key, m, cacheTTLModule)
	return m, nil
}

func (u *moduleUseCase) UpdateModule(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdModules, error) {
	m, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	out, err := u.repo.Update(ctx, m, values)
	if err != nil { return nil, err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyModuleDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:module:list:*", 200)
	return out, nil
}

func (u *moduleUseCase) DeleteModule(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil { return err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyModuleDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:module:list:*", 200)
	return nil
}

func (u *moduleUseCase) ListModules(ctx context.Context, f users.ModuleFilter) ([]*usermodels.SdModules, int64, error) {
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Limit > 100 { f.Limit = 100 }
	if f.Offset < 0 { f.Offset = 0 }
	key := fmt.Sprintf(cacheKeyModuleList, hashFilter(f))
	type cachedList struct {
		Items []*usermodels.SdModules `json:"items"`
		Total int64                   `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, f)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLModule)
	return items, total, nil
}

func (u *moduleUseCase) BulkUpdateModuleSort(ctx context.Context, items []users.SortItem) error {
	if err := u.repo.BulkUpdateSort(ctx, items); err != nil { return err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:module:list:*", 200)
	return nil
}

var _ users.ModuleUseCaseI = (*moduleUseCase)(nil)