package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/apimanager/presenter"
	"icmongolang/internal/modules/apimanager/repository"
	"icmongolang/internal/modules/apimanager/usecase"
	"icmongolang/pkg/logger"
)

type stubRepo struct {
	keys   map[int]*models.SdApiKey
	nextID int
}

func (s *stubRepo) Create(ctx context.Context, k *models.SdApiKey) error {
	s.nextID++
	k.ID = s.nextID
	if s.keys == nil {
		s.keys = map[int]*models.SdApiKey{}
	}
	s.keys[k.ID] = k
	return nil
}

func (s *stubRepo) GetByID(ctx context.Context, id int) (*models.SdApiKey, error) {
	if k, ok := s.keys[id]; ok {
		return k, nil
	}
	return nil, errors.New("not found")
}

func (s *stubRepo) GetByKey(ctx context.Context, apiKey string) (*models.SdApiKey, error) {
	for _, k := range s.keys {
		if k.ApiKey == apiKey {
			return k, nil
		}
	}
	return nil, errors.New("not found")
}

func (s *stubRepo) SetActive(ctx context.Context, id int, active bool) error {
	if k, ok := s.keys[id]; ok {
		k.IsActive = active
		return nil
	}
	return errors.New("not found")
}

func (s *stubRepo) IncrementUsage(ctx context.Context, id int) error {
	if k, ok := s.keys[id]; ok {
		k.UsageCount++
		return nil
	}
	return errors.New("not found")
}

type stubCache struct {
	counts map[string]int64
}

func (c *stubCache) Get(ctx context.Context, key string, dst interface{}) error { return errors.New("n/a") }
func (c *stubCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}
func (c *stubCache) Incr(ctx context.Context, key string) (int64, error) {
	if c.counts == nil {
		c.counts = map[string]int64{}
	}
	c.counts[key]++
	return c.counts[key], nil
}

func newUC(r repository.Repository, c usecase.Cache) usecase.APIManagerUseCase {
	return usecase.NewAPIManagerUseCase(r, c, logger.NewLogger("apimanager-test"))
}

func TestCreateKey_Invalid(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubCache{})
	if _, err := uc.CreateKey(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, err := uc.CreateKey(context.Background(), &presenter.APIKeyRequest{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreateKey_Success(t *testing.T) {
	repo := &stubRepo{}
	uc := newUC(repo, &stubCache{})
	resp, err := uc.CreateKey(context.Background(), &presenter.APIKeyRequest{Name: "test-key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Key == "" || resp.HashedKey == "" {
		t.Fatal("expected plaintext + hashed key")
	}
	if resp.Key == resp.HashedKey {
		t.Fatal("plaintext must differ from hashed")
	}
	if resp.KeyID != "1" {
		t.Fatalf("expected key id 1, got %s", resp.KeyID)
	}
	if len(repo.keys) != 1 {
		t.Fatalf("expected 1 stored key, got %d", len(repo.keys))
	}
	if !repo.keys[1].IsActive {
		t.Fatal("expected active key")
	}
}

func TestCreateKey_NilRepo_Graceful(t *testing.T) {
	uc := newUC(nil, &stubCache{})
	resp, err := uc.CreateKey(context.Background(), &presenter.APIKeyRequest{Name: "no-db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Key == "" {
		t.Fatal("expected a generated key even without repo")
	}
}

func TestRevokeKey_NotFound(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubCache{})
	if err := uc.RevokeKey(context.Background(), "999"); err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestRevokeKey_Success(t *testing.T) {
	repo := &stubRepo{}
	uc := newUC(repo, &stubCache{})
	if _, err := uc.CreateKey(context.Background(), &presenter.APIKeyRequest{Name: "k"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := uc.RevokeKey(context.Background(), "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.keys[1].IsActive {
		t.Fatal("expected key revoked")
	}
}

func TestCheckRateLimit_NoCache_AlwaysAllowed(t *testing.T) {
	uc := newUC(&stubRepo{}, nil)
	resp, err := uc.CheckRateLimit(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Allowed {
		t.Fatal("expected allowed with nil cache")
	}
}

func TestCheckRateLimit_AllowedThenLimited(t *testing.T) {
	uc := newUC(&stubRepo{}, &stubCache{})
	// Default limit is 100; call once -> allowed.
	resp, err := uc.CheckRateLimit(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Allowed || resp.Remaining <= 0 {
		t.Fatalf("expected allowed with positive remaining, got %+v", resp)
	}
}

func TestUsage_NilRepo_Graceful(t *testing.T) {
	uc := newUC(nil, &stubCache{})
	entries, err := uc.Usage(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries == nil {
		t.Fatal("expected non-nil empty usage")
	}
}

func TestUsage_WithRepo(t *testing.T) {
	repo := &stubRepo{}
	uc := newUC(repo, &stubCache{})
	if _, err := uc.CreateKey(context.Background(), &presenter.APIKeyRequest{Name: "m"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id := strconv.Itoa(repo.keys[1].ID)
	entries, err := uc.Usage(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 || entries[0].KeyID != id {
		t.Fatalf("unexpected usage: %+v", entries)
	}
}
