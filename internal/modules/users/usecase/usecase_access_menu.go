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
	cacheKeyAccessMenuList   = "sd_user:access_menu:list:%s"
	cacheKeyAccessMenuDetail = "sd_user:access_menu:%d"
	cacheTTLAccessMenu       = 5 * time.Minute
)

type accessMenuUseCase struct {
	repo   users.AccessMenuRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewAccessMenuUseCase(repo users.AccessMenuRepository, cache users.CacheRepository, l logger.Logger) users.AccessMenuUseCaseI {
	return &accessMenuUseCase{repo: repo, cache: cache, logger: l}
}

func (u *accessMenuUseCase) CreateAccessMenu(ctx context.Context, m *usermodels.SdUserAccessMenu) (*usermodels.SdUserAccessMenu, error) {
	out, err := u.repo.Create(ctx, m)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:access_menu:list:*", 200)
	return out, nil
}

func (u *accessMenuUseCase) GetAccessMenu(ctx context.Context, id int) (*usermodels.SdUserAccessMenu, error) {
	key := fmt.Sprintf(cacheKeyAccessMenuDetail, id)
	var c usermodels.SdUserAccessMenu
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return &c, nil }
	m, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	_ = u.cache.SetJSON(ctx, key, m, cacheTTLAccessMenu)
	return m, nil
}

func (u *accessMenuUseCase) UpdateAccessMenu(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserAccessMenu, error) {
	old, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	out, err := u.repo.Update(ctx, old, values)
	if err != nil { return nil, err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyAccessMenuDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:access_menu:list:*", 200)
	return out, nil
}

func (u *accessMenuUseCase) DeleteAccessMenu(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil { return err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyAccessMenuDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:access_menu:list:*", 200)
	return nil
}

func (u *accessMenuUseCase) ListAccessMenus(ctx context.Context, f users.AccessMenuFilter) ([]*usermodels.SdUserAccessMenu, int64, error) {
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Limit > 100 { f.Limit = 100 }
	if f.Offset < 0 { f.Offset = 0 }

	key := fmt.Sprintf(cacheKeyAccessMenuList, hashFilter(f))
	type cachedList struct {
		Items []*usermodels.SdUserAccessMenu `json:"items"`
		Total int64                          `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, f)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLAccessMenu)
	return items, total, nil
}

var _ users.AccessMenuUseCaseI = (*accessMenuUseCase)(nil)