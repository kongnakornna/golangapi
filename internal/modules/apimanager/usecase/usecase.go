package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/apimanager/presenter"
	"icmongolang/internal/modules/apimanager/repository"
	"icmongolang/pkg/logger"
)

// ErrInvalidRequest is returned for a malformed key request.
var ErrInvalidRequest = errors.New("apimanager: invalid request")

// ErrKeyNotFound is returned when a key id does not exist.
var ErrKeyNotFound = errors.New("apimanager: key not found")

// rateLimitWindow is the fixed-window rate-limit window length.
const rateLimitWindow = 60 * time.Second

// defaultRateLimit is the max requests per window per key.
const defaultRateLimit = 100

// Cache is the narrow counter backend the rate limiter needs (Redis-backed
// in production via redisDb.Cache; stubbed in tests).
type Cache interface {
	Get(ctx context.Context, key string, dst interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
}

// APIManagerUseCase defines the API-key management + rate-limit contract.
type APIManagerUseCase interface {
	CreateKey(ctx context.Context, req *presenter.APIKeyRequest) (*presenter.APIKeyResponse, error)
	RevokeKey(ctx context.Context, keyID string) error
	CheckRateLimit(ctx context.Context, keyID string) (*presenter.RateLimitResponse, error)
	Usage(ctx context.Context, keyID string) ([]presenter.UsageEntry, error)
}

type apiManagerUseCase struct {
	repo   repository.Repository
	cache  Cache
	logger logger.Logger
}

// NewAPIManagerUseCase builds the API manager use case. repo/cache may be nil
// (keys become no-op, rate-limit defaults to allowed) — existing flows never break.
func NewAPIManagerUseCase(repo repository.Repository, cache Cache, log logger.Logger) APIManagerUseCase {
	return &apiManagerUseCase{repo: repo, cache: cache, logger: log}
}

func (u *apiManagerUseCase) CreateKey(ctx context.Context, req *presenter.APIKeyRequest) (*presenter.APIKeyResponse, error) {
	if req == nil || req.Name == "" {
		return nil, ErrInvalidRequest
	}
	plainKey := randomKey()
	hashed := hashKey(plainKey)
	keyID := int64(0)
	if u.repo != nil {
		rec := &models.SdApiKey{
			Name:      req.Name,
			ApiKey:    hashed,
			ApiSecret: hashKey(plainKey + ":" + req.Name),
			IsActive:  true,
		}
		if err := u.repo.Create(ctx, rec); err != nil {
			u.logger.Errorf("apimanager: create key error: %v", err)
			return nil, err
		}
		keyID = int64(rec.ID)
	}
	u.logger.Infof("apimanager: created key id=%d name=%s", keyID, req.Name)
	return &presenter.APIKeyResponse{
		Key:       plainKey,
		KeyID:     strconv.FormatInt(keyID, 10),
		HashedKey: hashed,
		Name:      req.Name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (u *apiManagerUseCase) RevokeKey(ctx context.Context, keyID string) error {
	if u.repo == nil {
		return nil
	}
	id, err := parseIntID(keyID)
	if err != nil {
		return ErrInvalidRequest
	}
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return ErrKeyNotFound
	}
	u.logger.Infof("apimanager: revoke key id=%d", id)
	return u.repo.SetActive(ctx, id, false)
}

func (u *apiManagerUseCase) CheckRateLimit(ctx context.Context, keyID string) (*presenter.RateLimitResponse, error) {
	if keyID == "" {
		return nil, ErrInvalidRequest
	}
	if u.cache == nil {
		// No rate-limit backend wired: always allow.
		return &presenter.RateLimitResponse{Allowed: true, Limit: defaultRateLimit, Remaining: defaultRateLimit, ResetIn: int64(rateLimitWindow.Seconds())}, nil
	}
	key := "apikey:ratelimit:" + keyID
	now := time.Now().Unix()
	window := now / int64(rateLimitWindow.Seconds())
	windowKey := key + ":" + strconv.FormatInt(window, 10)
	count, err := u.cache.Incr(ctx, windowKey)
	if err != nil {
		// Fresh window or cache transient error: allow conservatively.
		return &presenter.RateLimitResponse{Allowed: true, Limit: defaultRateLimit, Remaining: defaultRateLimit, ResetIn: int64(rateLimitWindow.Seconds())}, nil
	}
	// Best-effort TTL so old windows expire.
	_ = u.cache.Set(ctx, windowKey, count, rateLimitWindow)
	remaining := defaultRateLimit - count
	allowed := remaining >= 0
	if remaining < 0 {
		remaining = 0
	}
	resetIn := (window + 1) * int64(rateLimitWindow.Seconds())
	if resetTo := resetIn - now; resetTo > 0 {
		resetIn = resetTo
	}
	return &presenter.RateLimitResponse{Allowed: allowed, Limit: defaultRateLimit, Remaining: remaining, ResetIn: resetIn}, nil
}

func (u *apiManagerUseCase) Usage(ctx context.Context, keyID string) ([]presenter.UsageEntry, error) {
	if u.repo == nil {
		return []presenter.UsageEntry{}, nil
	}
	id, err := parseIntID(keyID)
	if err != nil {
		return nil, ErrInvalidRequest
	}
	key, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrKeyNotFound
	}
	out := []presenter.UsageEntry{
		{KeyID: strconv.Itoa(key.ID), Route: "total", Status: 200, At: time.Now().UTC().Format(time.RFC3339), LatencyMs: int64(key.UsageCount)},
	}
	if key.LastUsedAt != nil {
		out[0].At = key.LastUsedAt.Format(time.RFC3339)
	}
	return out, nil
}

func parseIntID(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, ErrInvalidRequest
	}
	return n, nil
}

func randomKey() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "apk_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return "apk_" + hex.EncodeToString(b)
}

func hashKey(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
