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
	cacheKeyMenuList   = "sd_user:menu:list:%s"
	cacheKeyMenuDetail = "sd_user:menu:%d"
	cacheTTLMenu       = 10 * time.Minute
)

type menuUseCase struct {
	repo   users.MenuRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewMenuUseCase(repo users.MenuRepository, cache users.CacheRepository, l logger.Logger) users.MenuUseCaseI {
	return &menuUseCase{repo: repo, cache: cache, logger: l}
}

func (u *menuUseCase) CreateMenu(ctx context.Context, m *usermodels.SdMenu) (*usermodels.SdMenu, error) {
	out, err := u.repo.Create(ctx, m)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:menu:list:*", 200)
	return out, nil
}

func (u *menuUseCase) GetMenu(ctx context.Context, id int) (*usermodels.SdMenu, error) {
	key := fmt.Sprintf(cacheKeyMenuDetail, id)
	var c usermodels.SdMenu
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return &c, nil }
	m, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	_ = u.cache.SetJSON(ctx, key, m, cacheTTLMenu)
	return m, nil
}

func (u *menuUseCase) UpdateMenu(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdMenu, error) {
	m, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	out, err := u.repo.Update(ctx, m, values)
	if err != nil { return nil, err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyMenuDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:menu:list:*", 200)
	return out, nil
}

func (u *menuUseCase) DeleteMenu(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil { return err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyMenuDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:menu:list:*", 200)
	return nil
}

func (u *menuUseCase) ListMenus(ctx context.Context, f users.MenuFilter) ([]*usermodels.SdMenu, int64, error) {
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Limit > 100 { f.Limit = 100 }
	if f.Offset < 0 { f.Offset = 0 }
	key := fmt.Sprintf(cacheKeyMenuList, hashFilter(f))
	type cachedList struct {
		Items []*usermodels.SdMenu `json:"items"`
		Total int64                `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, f)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLMenu)
	return items, total, nil
}

func (u *menuUseCase) BulkUpdateMenuSort(ctx context.Context, items []users.SortItem) error {
	if err := u.repo.BulkUpdateSort(ctx, items); err != nil { return err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:menu:list:*", 200)
	return nil
}

var _ users.MenuUseCaseI = (*menuUseCase)(nil)