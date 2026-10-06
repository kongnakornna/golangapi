package repository

import (
	"context"

	"github.com/google/uuid"

	pkgapikey "icmongolang/pkg/apikey"
)

// APIKeyRepository is the persistence contract for API keys.
//
// Implementations MUST map gorm.ErrRecordNotFound (or equivalent) to
// pkgapikey.ErrNotFound so the use-case layer never leaks storage errors.
type APIKeyRepository interface {
	// Create persists a new key.
	Create(ctx context.Context, k *pkgapikey.Key) (*pkgapikey.Key, error)

	// FindByID returns a key by primary key.
	// Returns pkgapikey.ErrNotFound when missing.
	FindByID(ctx context.Context, id uuid.UUID) (*pkgapikey.Key, error)

	// FindByHashedKey is the hot path used by the auth middleware.
	// Returns pkgapikey.ErrNotFound when missing.
	FindByHashedKey(ctx context.Context, hash string) (*pkgapikey.Key, error)

	// List returns every key (admin view), paginated.
	List(ctx context.Context, limit, offset int) ([]*pkgapikey.Key, int64, error)

	// ListByOwner returns keys whose CreatedBy == owner, paginated (manager view).
	ListByOwner(ctx context.Context, owner uuid.UUID, limit, offset int) ([]*pkgapikey.Key, int64, error)

	// Rotate overwrites the fingerprint / prefix / last_four of an existing key.
	Rotate(ctx context.Context, k *pkgapikey.Key) (*pkgapikey.Key, error)

	// Deactivate revokes a key by ID (admin path).
	Deactivate(ctx context.Context, id, updatedBy uuid.UUID) error

	// DeactivateOwned revokes a key only if CreatedBy == owner (manager path).
	DeactivateOwned(ctx context.Context, id, owner, updatedBy uuid.UUID) error
}