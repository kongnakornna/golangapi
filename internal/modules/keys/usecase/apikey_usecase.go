package usecase

import (
	"context"

	"github.com/google/uuid"

	pkgapikey "icmongolang/pkg/apikey"
)

// Actor identifies the caller driving a use-case invocation.
//
// IsAdmin is computed by the HTTP layer from JWT claims; keeping it as a
// plain bool here avoids dragging auth-domain types into the use case.
type Actor struct {
	ID      uuid.UUID
	IsAdmin bool
}

// APIKeyUseCase drives the API-key lifecycle.
type APIKeyUseCase interface {
	// Create mints a new key. The raw secret is returned exactly once.
	Create(ctx context.Context, name string, description *string, expiresIn *int, actor Actor) (*pkgapikey.Key, error)

	// List returns paginated API keys visible to actor.
	// Admins see every key; managers see only their own.
	List(ctx context.Context, limit, offset int, actor Actor) ([]*pkgapikey.Key, int64, error)

	// GetByID returns a key by ID. Non-admins must own it.
	GetByID(ctx context.Context, id uuid.UUID, actor Actor) (*pkgapikey.Key, error)

	// Rotate revokes the outgoing secret and issues a fresh one.
	// Non-admins must own the key.
	Rotate(ctx context.Context, id uuid.UUID, actor Actor) (*pkgapikey.Key, error)

	// Revoke soft-deletes a key. Non-admins must own it.
	Revoke(ctx context.Context, id uuid.UUID, actor Actor) error
}