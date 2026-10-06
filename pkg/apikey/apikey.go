package apikey

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalid     = errors.New("api key is invalid")
	ErrRevoked     = errors.New("api key has been revoked")
	ErrExpired     = errors.New("api key has expired")
	ErrEmptyPepper = errors.New("pepper must not be empty")
	ErrNotFound   = errors.New("api key not found")
)

// Config holds generation parameters.
type Config struct {
	PrefixLen int           // visible prefix length (default 8)
	SecretLen int           // raw secret entropy bytes (default 32)
	Pepper    []byte        // HMAC key, loaded from env/secret manager
	TTL       time.Duration // optional expiry (0 = no expiry)
	Issuer    string        // optional issuer claim for audit
}

// DefaultConfig returns sane defaults.
func DefaultConfig(pepper []byte) *Config {
	return &Config{
		PrefixLen: 8,
		SecretLen: 32,
		Pepper:    pepper,
		TTL:       0,
		Issuer:    "icmongolang",
	}
}

// GeneratedKey is returned exactly once on Create/Rotate.
// The raw secret is NEVER persisted — only the HMAC fingerprint (Hash) is stored.
type GeneratedKey struct {
	ID        uuid.UUID
	Prefix    string
	LastFour  string
	Hash      string // base64(HMAC-SHA256(pepper, raw))
	Raw       string // prefix + "." + base64(secret) — returned once
	ExpiresAt *time.Time
}

// APIKey is the stored entity (matches domain/apikey/entity.go).
type APIKey struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description *string    `gorm:"column:description" json:"description,omitempty"`
	Prefix      string     `gorm:"column:prefix;not null;index" json:"prefix"`
	LastFour    string     `gorm:"column:last_four;not null" json:"last_four"`
	HashedKey   string     `gorm:"column:hashed_key;not null;uniqueIndex" json:"-"`
	ExpiresAt   *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at" json:"last_used_at,omitempty"`
	CreatedBy   uuid.UUID  `gorm:"column:created_by;not null;index" json:"created_by"`
	UpdatedBy   uuid.UUID  `gorm:"column:updated_by;not null;index" json:"updated_by"`
	IsActive    bool       `gorm:"column:is_active;not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Transient — populated only on Create/Rotate, never persisted.
	PlainKey string `gorm:"-" json:"api_key,omitempty"`
}

// TableName pins the entity.
func (APIKey) TableName() string { return "api_keys" }

// Generate mints a new API key. The raw secret is returned exactly once.
func Generate(cfg *Config) (*GeneratedKey, error) {
	if len(cfg.Pepper) == 0 {
		return nil, ErrEmptyPepper
	}

	// 1. Generate cryptographically random secret
	secret := make([]byte, cfg.SecretLen)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("rand: %w", err)
	}

	// 2. Compute HMAC-SHA256 fingerprint (what we store)
	mac := hmac.New(sha256.New, cfg.Pepper)
	mac.Write(secret)
	hash := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 3. Build raw key: prefix + "." + base64(secret)
	//    Prefix is first N chars of base64(secret) for human readability.
	b64Secret := base64.RawURLEncoding.EncodeToString(secret)
	prefix := b64Secret[:cfg.PrefixLen]
	lastFour := b64Secret[len(b64Secret)-4:]
	raw := prefix + "." + b64Secret

	var exp *time.Time
	if cfg.TTL != 0 {
		t := time.Now().Add(cfg.TTL)
		exp = &t
	}

	return &GeneratedKey{
		ID:        uuid.New(),
		Prefix:    prefix,
		LastFour:  lastFour,
		Hash:      hash,
		Raw:       raw,
		ExpiresAt: exp,
	}, nil
}

// Hash computes the stored fingerprint for a raw key (used on Resolve).
func Hash(pepper []byte, raw string) string {
	// raw format: prefix.base64secret
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	secret, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	mac := hmac.New(sha256.New, pepper)
	mac.Write(secret)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Validate checks expiry and active flag (called after Resolve).
func (k *Key) Validate() error {
	if !k.IsActive {
		return ErrRevoked
	}
	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return ErrExpired
	}
	return nil
}

// TouchLastUsed updates LastUsedAt (call on successful Resolve).
func (k *Key) TouchLastUsed() {
	now := time.Now()
	k.LastUsedAt = &now
	k.UpdatedAt = now
}
