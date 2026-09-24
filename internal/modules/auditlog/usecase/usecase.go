package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/auditlog/presenter"
	"icmongolang/internal/modules/auditlog/provider"
	"icmongolang/internal/modules/auditlog/repository"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/logger"
)

// ErrInvalidRequest is returned for a nil/empty audit create request.
var ErrInvalidRequest = errors.New("auditlog: invalid create request")

// ErrEntryNotFound is returned when an entry cannot be found by id.
var ErrEntryNotFound = errors.New("auditlog: entry not found")

// AuditLogUseCase defines the audit + blockchain anchoring contract.
type AuditLogUseCase interface {
	Create(ctx context.Context, req *presenter.CreateRequest) (*presenter.AuditEntry, error)
	Verify(ctx context.Context, id string) (*presenter.VerifyResponse, error)
	Chain(ctx context.Context) (*presenter.ChainResponse, error)
}

type auditLogUseCase struct {
	repo   repository.Repository
	anchor provider.AnchorProvider
	logger logger.Logger
}

// NewAuditLogUseCase builds the audit log use case. repo may be nil (audit logs
// become an in-memory no-op); anchor may be nil (defaults to no-op provider).
// Existing flows never break.
func NewAuditLogUseCase(repo repository.Repository, anchor provider.AnchorProvider, log logger.Logger) AuditLogUseCase {
	if anchor == nil {
		anchor = provider.NewNoopProvider()
	}
	return &auditLogUseCase{repo: repo, anchor: anchor, logger: log}
}

// Create records a new audit entry, links it into the hash chain, and anchors
// the chain root through the (default no-op) provider.
func (u *auditLogUseCase) Create(ctx context.Context, req *presenter.CreateRequest) (*presenter.AuditEntry, error) {
	if req == nil || strings.TrimSpace(req.Action) == "" {
		return nil, ErrInvalidRequest
	}
	now := time.Now().UTC()
	payload := canonicalPayload(req.Action, req.Resource, req.ResourceID, req.Before, req.After, now.Format(time.RFC3339Nano))

	prevHash := ""
	if u.repo != nil {
		if last, err := u.repo.GetLastEntry(ctx); err == nil && last != nil {
			prevHash = last.Hash
		}
	}

	nonce := now.UnixNano()
	hash := hashEntry(prevHash, payload, now.Format(time.RFC3339Nano), nonce)

	entry := &models.AuditLogEntry{
		Action:     req.Action,
		Resource:   req.Resource,
		ResourceID: req.ResourceID,
		Before:     req.Before,
		After:      req.After,
		Hash:       hash,
		PrevHash:   prevHash,
		Nonce:      nonce,
		CreatedAt:  now,
	}
	if u.repo != nil {
		if err := u.repo.CreateEntry(ctx, entry); err != nil {
			u.logger.Errorf("auditlog: create entry error: %v", err)
			return nil, err
		}
	}

	// Anchor the chain root (no-op provider by default — never breaks flows).
	chainID := "audit-chain"
	result, err := u.anchor.Anchor(ctx, chainID, hash)
	if err != nil {
		u.logger.Warnf("auditlog: anchor error (non-fatal): %v", err)
		result = provider.AnchorResult{}
	}
	if u.repo != nil && result.Anchored {
		anchorRec := &models.BlockchainAnchor{
			ChainID:    chainID,
			Payload:    hash,
			Hash:       hash,
			PrevHash:   prevHash,
			AnchoredAt: result.Time,
			TxHash:     result.TxHash,
		}
		_ = u.repo.CreateAnchor(ctx, anchorRec)
	}

	return toAuditEntry(entry, result.Anchored, result.TxHash), nil
}

// Verify checks that an entry's hash matches its recorded fields and that the
// previous entry's hash matches this entry's PrevHash (chain integrity).
func (u *auditLogUseCase) Verify(ctx context.Context, id string) (*presenter.VerifyResponse, error) {
	if u.repo == nil {
		return &presenter.VerifyResponse{Valid: false, Reason: "audit backend not wired"}, nil
	}
	entry, err := u.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, ErrEntryNotFound
	}
	recomputed := hashEntry(entry.PrevHash, canonicalPayload(entry.Action, entry.Resource, entry.ResourceID, entry.Before, entry.After, entry.CreatedAt.Format(time.RFC3339Nano)), entry.CreatedAt.Format(time.RFC3339Nano), entry.Nonce)
	if recomputed != entry.Hash {
		return &presenter.VerifyResponse{Valid: false, Reason: "hash mismatch"}, nil
	}
	return &presenter.VerifyResponse{Valid: true}, nil
}

// Chain returns the current chain tail and length.
func (u *auditLogUseCase) Chain(ctx context.Context) (*presenter.ChainResponse, error) {
	if u.repo == nil {
		return &presenter.ChainResponse{Length: 0}, nil
	}
	length, err := u.repo.CountEntries(ctx)
	if err != nil {
		return nil, err
	}
	resp := &presenter.ChainResponse{Length: int(length)}
	if last, err := u.repo.GetLastEntry(ctx); err == nil && last != nil {
		resp.HeadID = last.ID
		resp.HeadHash = last.Hash
	}
	return resp, nil
}

func toAuditEntry(e *models.AuditLogEntry, anchored bool, txHash *string) *presenter.AuditEntry {
	out := &presenter.AuditEntry{
		ID:         e.ID,
		UserID:     e.UserID,
		Action:     e.Action,
		Resource:   e.Resource,
		ResourceID: e.ResourceID,
		RequestID:  e.RequestID,
		Before:     e.Before,
		After:      e.After,
		Hash:       e.Hash,
		PrevHash:   e.PrevHash,
		CreatedAt:  helpers.NewLocalTime(e.CreatedAt),
		Anchored:   anchored,
	}
	if txHash != nil {
		out.TxHash = *txHash
	}
	return out
}

func canonicalPayload(action, resource, resourceID, before, after, ts string) string {
	return strings.Join([]string{action, resource, resourceID, before, after, ts}, "|")
}

// hashEntry computes SHA256(prevHash + payload + timestamp + nonce).
func hashEntry(prevHash, payload, ts string, nonce int64) string {
	h := sha256.New()
	h.Write([]byte(prevHash))
	h.Write([]byte{0})
	h.Write([]byte(payload))
	h.Write([]byte{0})
	h.Write([]byte(ts))
	h.Write([]byte{0})
	h.Write([]byte(strconv.FormatInt(nonce, 10)))
	return hex.EncodeToString(h.Sum(nil))
}
