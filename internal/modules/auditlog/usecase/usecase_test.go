package usecase_test

import (
	"context"
	"errors"
	"testing"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/auditlog/presenter"
	"icmongolang/internal/modules/auditlog/provider"
	"icmongolang/internal/modules/auditlog/repository"
	"icmongolang/internal/modules/auditlog/usecase"
	"icmongolang/pkg/logger"
)

type memRepo struct {
	entries []*models.AuditLogEntry
	anchors []*models.BlockchainAnchor
}

func (m *memRepo) CreateEntry(ctx context.Context, e *models.AuditLogEntry) error {
	e.ID = "id-" + itoa(len(m.entries)+1)
	m.entries = append(m.entries, e)
	return nil
}

func (m *memRepo) GetEntryByID(ctx context.Context, id string) (*models.AuditLogEntry, error) {
	for _, e := range m.entries {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *memRepo) GetLastEntry(ctx context.Context) (*models.AuditLogEntry, error) {
	if len(m.entries) == 0 {
		return nil, errors.New("empty")
	}
	return m.entries[len(m.entries)-1], nil
}

func (m *memRepo) CreateAnchor(ctx context.Context, a *models.BlockchainAnchor) error {
	m.anchors = append(m.anchors, a)
	return nil
}

func (m *memRepo) GetAnchorByChainID(ctx context.Context, chainID string) (*models.BlockchainAnchor, error) {
	return nil, errors.New("none")
}

func (m *memRepo) CountEntries(ctx context.Context) (int64, error) {
	return int64(len(m.entries)), nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := []byte{}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func newUC(r repository.Repository, a provider.AnchorProvider) usecase.AuditLogUseCase {
	return usecase.NewAuditLogUseCase(r, a, logger.NewLogger("audit-test"))
}

func replayAnchor() provider.AnchorProvider {
	return &replayProvider{}
}

type replayProvider struct{}

func (p *replayProvider) Anchor(ctx context.Context, chainID, rootHash string) (provider.AnchorResult, error) {
	return provider.AnchorResult{Anchored: true}, nil
}

func TestCreate_InvalidRequest(t *testing.T) {
	uc := newUC(&memRepo{}, replayAnchor())
	if _, err := uc.Create(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, err := uc.Create(context.Background(), &presenter.CreateRequest{}); err == nil {
		t.Fatal("expected error for empty action")
	}
}

func TestCreate_NilRepo_Graceful(t *testing.T) {
	uc := newUC(nil, provider.NewNoopProvider())
	entry, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "settings.update"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Hash == "" {
		t.Fatal("expected a hash even without repo")
	}
	if entry.Anchored {
		t.Fatalf("expected not anchored with no-op default")
	}
}

func TestCreate_ChainContinuity(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, replayAnchor())
	e1, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "device.control", Resource: "7", After: "ON"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e2, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "settings.update", Resource: "email", After: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e1.Hash == "" || e1.Hash == e2.Hash {
		t.Fatalf("expected distinct non-empty hashes: %q %q", e1.Hash, e2.Hash)
	}
	if e1.PrevHash != "" {
		t.Fatalf("expected first entry prev_hash empty, got %q", e1.PrevHash)
	}
	if e2.PrevHash != e1.Hash {
		t.Fatalf("expected e2.prev_hash == e1.hash, got %q vs %q", e2.PrevHash, e1.Hash)
	}
}

func TestCreate_NoopAnchor_TxHashNil(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, provider.NewNoopProvider())
	entry, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Anchored {
		t.Fatal("expected not anchored with noop provider")
	}
	if entry.TxHash != "" {
		t.Fatalf("expected empty tx_hash, got %q", entry.TxHash)
	}
	if len(repo.anchors) != 0 {
		t.Fatalf("no-op provider must not write anchors, got %d", len(repo.anchors))
	}
}

func TestCreate_ReplayAnchor_WritesAnchor(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, replayAnchor())
	_, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.anchors) != 1 {
		t.Fatalf("expected 1 anchor written, got %d", len(repo.anchors))
	}
	if repo.anchors[0].TxHash != nil {
		t.Fatalf("replay provider has no tx hash, got %v", *repo.anchors[0].TxHash)
	}
}

func TestVerify_Valid_And_Tampered(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, provider.NewNoopProvider())
	entry, err := uc.Create(context.Background(), &presenter.CreateRequest{Action: "device.control", Resource: "7", After: "ON"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := uc.Verify(context.Background(), entry.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.Valid {
		t.Fatalf("expected valid entry, got %+v", v)
	}
	// Tamper with the stored record.
	last := repo.entries[len(repo.entries)-1]
	last.After = "OFF"
	v2, _ := uc.Verify(context.Background(), entry.ID)
	if v2.Valid {
		t.Fatalf("expected tampered entry to fail verification")
	}
}

func TestVerify_NotFound(t *testing.T) {
	uc := newUC(&memRepo{}, provider.NewNoopProvider())
	if _, err := uc.Verify(context.Background(), "missing"); err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestChain_LengthAndHead(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, provider.NewNoopProvider())
	_, _ = uc.Create(context.Background(), &presenter.CreateRequest{Action: "a"})
	_, _ = uc.Create(context.Background(), &presenter.CreateRequest{Action: "b"})
	c, err := uc.Chain(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Length != 2 {
		t.Fatalf("expected length 2, got %d", c.Length)
	}
	if c.HeadID == "" || c.HeadHash == "" {
		t.Fatal("expected head id/hash populated")
	}
}

func TestChain_NilRepo(t *testing.T) {
	uc := newUC(nil, provider.NewNoopProvider())
	c, err := uc.Chain(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Length != 0 {
		t.Fatalf("expected length 0, got %d", c.Length)
	}
}
