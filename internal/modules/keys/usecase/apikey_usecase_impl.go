package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	keyrepo "icmongolang/internal/modules/keys/repository"
	pkgapikey "icmongolang/pkg/apikey"
)

// APIKeyUseCaseImpl drives the API-key lifecycle.
type APIKeyUseCaseImpl struct {
	repo   keyrepo.APIKeyRepository
	pepper []byte
}

// NewAPIKeyUseCase wires dependencies.
func NewAPIKeyUseCase(repo keyrepo.APIKeyRepository, pepper []byte) *APIKeyUseCaseImpl {
	return &APIKeyUseCaseImpl{repo: repo, pepper: pepper}
}

// ---------------------------------------------------------------- create

func (u *APIKeyUseCaseImpl) Create(
	ctx context.Context,
	name string,
	description *string,
	expiresIn *int,
	actor Actor,
) (*pkgapikey.Key, error) {
	cfg := &pkgapikey.Config{Pepper: u.pepper}
	if expiresIn != nil && *expiresIn > 0 {
		cfg.TTL = time.Duration(*expiresIn) * time.Minute
	}

	gen, err := pkgapikey.Generate(cfg)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	k := &pkgapikey.Key{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Prefix:      gen.Prefix,
		LastFour:    gen.LastFour,
		HashedKey:   gen.Hash,
		CreatedBy:   actor.ID,
		UpdatedBy:   actor.ID,
		IsActive:    true,
		ExpiresAt:   gen.ExpiresAt,
	}

	created, err := u.repo.Create(ctx, k)
	if err != nil {
		return nil, err
	}
	created.PlainKey = gen.Raw
	return created, nil
}

// ---------------------------------------------------------------- list

func (u *APIKeyUseCaseImpl) List(
	ctx context.Context,
	limit, offset int,
	actor Actor,
) ([]*pkgapikey.Key, int64, error) {
	if actor.IsAdmin {
		return u.repo.List(ctx, limit, offset)
	}
	return u.repo.ListByOwner(ctx, actor.ID, limit, offset)
}

// ---------------------------------------------------------------- get

func (u *APIKeyUseCaseImpl) GetByID(ctx context.Context, id uuid.UUID, actor Actor) (*pkgapikey.Key, error) {
	k, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.IsAdmin && k.CreatedBy != actor.ID {
		// Return NotFound rather than Forbidden so non-admins cannot
		// probe whether a given key ID exists.
		return nil, pkgapikey.ErrNotFound
	}
	return k, nil
}

// ---------------------------------------------------------------- rotate

func (u *APIKeyUseCaseImpl) Rotate(ctx context.Context, id uuid.UUID, actor Actor) (*pkgapikey.Key, error) {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.IsAdmin && existing.CreatedBy != actor.ID {
		return nil, pkgapikey.ErrNotFound
	}

	gen, err := pkgapikey.Generate(&pkgapikey.Config{Pepper: u.pepper})
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	existing.Prefix = gen.Prefix
	existing.LastFour = gen.LastFour
	existing.HashedKey = gen.Hash
	existing.UpdatedBy = actor.ID

	rotated, err := u.repo.Rotate(ctx, existing)
	if err != nil {
		return nil, err
	}
	rotated.PlainKey = gen.Raw
	return rotated, nil
}

// ---------------------------------------------------------------- revoke

func (u *APIKeyUseCaseImpl) Revoke(ctx context.Context, id uuid.UUID, actor Actor) error {
	if actor.IsAdmin {
		return u.repo.Deactivate(ctx, id, actor.ID)
	}

	// Ownership check up-front so we can return a clean domain error
	// instead of relying on the UPDATE's WHERE clause to filter silently.
	k, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if k.CreatedBy != actor.ID {
		return pkgapikey.ErrNotFound
	}
	return u.repo.DeactivateOwned(ctx, id, actor.ID, actor.ID)
}