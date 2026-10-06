package apikey

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerate_RoundTrip(t *testing.T) {
	pepper := []byte("test-pepper-32-bytes-minimum!!")
	cfg := DefaultConfig(pepper)

	gk, err := Generate(cfg)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Hash(raw) must equal stored Hash
	if Hash(pepper, gk.Raw) != gk.Hash {
		t.Fatalf("hash mismatch: got %s, want %s", Hash(pepper, gk.Raw), gk.Hash)
	}

	// Raw format: prefix.base64secret
	parts := splitRaw(gk.Raw)
	if parts[0] != gk.Prefix {
		t.Fatalf("prefix mismatch: got %s, want %s", parts[0], gk.Prefix)
	}
	if parts[1][len(parts[1])-4:] != gk.LastFour {
		t.Fatalf("lastFour mismatch")
	}
}

func TestGenerate_DeterministicHash(t *testing.T) {
	pepper := []byte("another-pepper-32-bytes-min!!")
	cfg := DefaultConfig(pepper)

	gk1, _ := Generate(cfg)
	gk2, _ := Generate(cfg)

	// Different secrets → different hashes
	if gk1.Hash == gk2.Hash {
		t.Fatal("hashes should differ")
	}
	// But same raw → same hash
	if Hash(pepper, gk1.Raw) != gk1.Hash {
		t.Fatal("hash not deterministic")
	}
}

func TestHash_InvalidRaw(t *testing.T) {
	pepper := []byte("pepper-pepper-pepper-pepper-pepp")
	if h := Hash(pepper, "no-dot"); h != "" {
		t.Fatalf("expected empty hash for malformed raw, got %s", h)
	}
	if h := Hash(pepper, "prefix.notbase64!"); h != "" {
		t.Fatalf("expected empty hash for invalid base64, got %s", h)
	}
}

func TestValidate_Expiry(t *testing.T) {
	pepper := []byte("pepper-pepper-pepper-pepper-pepp")
	cfg := DefaultConfig(pepper)
	cfg.TTL = -time.Hour // already expired

	gk, _ := Generate(cfg)
	k := &Key{
		ID:        gk.ID,
		Prefix:    gk.Prefix,
		LastFour:  gk.LastFour,
		HashedKey: gk.Hash,
		ExpiresAt: gk.ExpiresAt,
		IsActive:  true,
	}
	if err := k.Validate(); err != ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestValidate_Revoked(t *testing.T) {
	k := &Key{IsActive: false}
	if err := k.Validate(); err != ErrRevoked {
		t.Fatalf("expected ErrRevoked, got %v", err)
	}
}

func TestValidate_OK(t *testing.T) {
	pepper := []byte("pepper-pepper-pepper-pepper-pepp")
	cfg := DefaultConfig(pepper)
	gk, _ := Generate(cfg)
	k := &Key{
		ID:        gk.ID,
		Prefix:    gk.Prefix,
		LastFour:  gk.LastFour,
		HashedKey: gk.Hash,
		ExpiresAt: gk.ExpiresAt,
		IsActive:  true,
	}
	if err := k.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTouchLastUsed(t *testing.T) {
	k := &Key{ID: uuid.New()}
	before := k.UpdatedAt
	k.TouchLastUsed()
	if k.LastUsedAt == nil || k.UpdatedAt.Equal(before) {
		t.Fatal("LastUsedAt/UpdatedAt not updated")
	}
}

func splitRaw(raw string) []string {
	// helper for test only
	for i := 0; i < len(raw); i++ {
		if raw[i] == '.' {
			return []string{raw[:i], raw[i+1:]}
		}
	}
	return []string{raw, ""}
}
