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
	cacheKeyFileList   = "sd_user:file:list:%s"
	cacheKeyFileDetail = "sd_user:file:%d"
	cacheTTLFile       = 3 * time.Minute
)

type fileUseCase struct {
	repo   users.FileRepository
	cache  users.CacheRepository
	logger logger.Logger
}

func NewFileUseCase(repo users.FileRepository, cache users.CacheRepository, l logger.Logger) users.FileUseCaseI {
	return &fileUseCase{repo: repo, cache: cache, logger: l}
}

func (u *fileUseCase) CreateFile(ctx context.Context, f *usermodels.SdUserFile) (*usermodels.SdUserFile, error) {
	out, err := u.repo.Create(ctx, f)
	if err != nil { return nil, err }
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:file:list:*", 200)
	return out, nil
}

func (u *fileUseCase) GetFile(ctx context.Context, id int) (*usermodels.SdUserFile, error) {
	key := fmt.Sprintf(cacheKeyFileDetail, id)
	var c usermodels.SdUserFile
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return &c, nil }
	f, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	_ = u.cache.SetJSON(ctx, key, f, cacheTTLFile)
	return f, nil
}

func (u *fileUseCase) UpdateFile(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserFile, error) {
	f, err := u.repo.Get(ctx, id)
	if err != nil { return nil, err }
	out, err := u.repo.Update(ctx, f, values)
	if err != nil { return nil, err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyFileDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:file:list:*", 200)
	return out, nil
}

func (u *fileUseCase) DeleteFile(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil { return err }
	_ = u.cache.Delete(ctx, fmt.Sprintf(cacheKeyFileDetail, id))
	_, _ = u.cache.ScanAndDelete(ctx, "sd_user:file:list:*", 200)
	return nil
}

func (u *fileUseCase) ListFiles(ctx context.Context, f users.FileFilter) ([]*usermodels.SdUserFile, int64, error) {
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Limit > 100 { f.Limit = 100 }
	if f.Offset < 0 { f.Offset = 0 }
	key := fmt.Sprintf(cacheKeyFileList, hashFilter(f))
	type cachedList struct {
		Items []*usermodels.SdUserFile `json:"items"`
		Total int64                    `json:"total"`
	}
	var c cachedList
	if hit, _ := u.cache.GetJSON(ctx, key, &c); hit { return c.Items, c.Total, nil }
	items, total, err := u.repo.List(ctx, f)
	if err != nil { return nil, 0, err }
	_ = u.cache.SetJSON(ctx, key, cachedList{Items: items, Total: total}, cacheTTLFile)
	return items, total, nil
}

var _ users.FileUseCaseI = (*fileUseCase)(nil)