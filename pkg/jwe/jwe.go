package jwe

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var (
	ErrInvalidKey         = errors.New("invalid key")
	ErrEncryptionFailed   = errors.New("encryption failed")
	ErrDecryptionFailed   = errors.New("decryption failed")
	ErrSigningFailed      = errors.New("signing failed")
	ErrVerificationFailed = errors.New("verification failed")
	ErrNilKeyPair         = errors.New("jwe: nil key pair (JWE keys are not configured)")
)

// KeyPair holds Ed25519 signing keys and X25519 key agreement keys.
type KeyPair struct {
	SigningPrivate ed25519.PrivateKey // 64 bytes (seed + public)
	SigningPublic  ed25519.PublicKey  // 32 bytes
	KEXPrivate     *ecdh.PrivateKey   // X25519 for ECDH-ES
	KEXPublic      *ecdh.PublicKey    // X25519 public
}

// IsReady returns true if the key pair is non-nil and all sub-keys are set.
func (kp *KeyPair) IsReady() bool {
	return kp != nil &&
		len(kp.SigningPrivate) == ed25519.PrivateKeySize &&
		len(kp.SigningPublic) == ed25519.PublicKeySize &&
		kp.KEXPrivate != nil &&
		kp.KEXPublic != nil
}

// GenerateKeyPair creates a new Ed25519/X25519 key pair.
func GenerateKeyPair() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ed25519 generate: %w", err)
	}
	kexPriv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("x25519 generate: %w", err)
	}
	return &KeyPair{
		SigningPrivate: priv,
		SigningPublic:  pub,
		KEXPrivate:     kexPriv,
		KEXPublic:      kexPriv.PublicKey(),
	}, nil
}

// LoadKeyPair loads keys from base64 encoded strings.
func LoadKeyPair(signingPrivB64, kexPrivB64 string) (*KeyPair, error) {
	if signingPrivB64 == "" || kexPrivB64 == "" {
		return nil, ErrInvalidKey
	}

	signingPriv, err := base64.StdEncoding.DecodeString(signingPrivB64)
	if err != nil {
		return nil, fmt.Errorf("decode signing priv: %w", err)
	}
	if len(signingPriv) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("%w: signing priv must be %d bytes, got %d",
			ErrInvalidKey, ed25519.PrivateKeySize, len(signingPriv))
	}

	kexPrivBytes, err := base64.StdEncoding.DecodeString(kexPrivB64)
	if err != nil {
		return nil, fmt.Errorf("decode kex priv: %w", err)
	}
	kexPriv, err := ecdh.X25519().NewPrivateKey(kexPrivBytes)
	if err != nil {
		return nil, fmt.Errorf("parse x25519 priv: %w", err)
	}

	return &KeyPair{
		SigningPrivate: ed25519.PrivateKey(signingPriv),
		SigningPublic:  signingPriv[32:],
		KEXPrivate:     kexPriv,
		KEXPublic:      kexPriv.PublicKey(),
	}, nil
}

// NestedClaims represents the inner JWT payload.
type NestedClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	SessionID string    `json:"session_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	JWTID     string    `json:"jti"`
}

// Seal creates a nested JWT: Ed25519 JWS (inner) + ECDH-ES/A256GCM JWE (outer).
func Seal(kp *KeyPair, recipientPub *ecdh.PublicKey, claims *NestedClaims, ttl time.Duration) (string, error) {
	if kp == nil {
		return "", ErrNilKeyPair
	}
	if !kp.IsReady() {
		return "", fmt.Errorf("%w: key pair is incomplete", ErrInvalidKey)
	}
	if claims == nil {
		return "", errors.New("nil claims")
	}
	if recipientPub == nil {
		recipientPub = kp.KEXPublic
	}

	now := time.Now()
	claims.IssuedAt = now
	claims.ExpiresAt = now.Add(ttl)
	if claims.JWTID == "" {
		claims.JWTID = uuid.New().String()
	}

	innerToken, err := jwt.NewBuilder().
		Claim("user_id", claims.UserID).
		Claim("email", claims.Email).
		Claim("role", claims.Role).
		Claim("session_id", claims.SessionID).
		Claim("jti", claims.JWTID).
		Expiration(claims.ExpiresAt).
		IssuedAt(claims.IssuedAt).
		Build()
	if err != nil {
		return "", fmt.Errorf("build inner jwt: %w", err)
	}

	innerBytes, err := json.Marshal(innerToken)
	if err != nil {
		return "", fmt.Errorf("marshal inner: %w", err)
	}

	signed, err := jws.Sign(innerBytes, jws.WithKey(jwa.EdDSAEd25519(), kp.SigningPrivate))
	if err != nil {
		return "", fmt.Errorf("jws sign: %w: %v", ErrSigningFailed, err)
	}

	encrypted, err := jwe.Encrypt(
		signed,
		jwe.WithKey(jwa.ECDH_ES_A256KW(), recipientPub),
		jwe.WithContentEncryption(jwa.A256GCM()),
	)
	if err != nil {
		return "", fmt.Errorf("jwe encrypt: %w: %v", ErrEncryptionFailed, err)
	}

	return string(encrypted), nil
}

// Open decrypts and verifies a nested JWT.
func Open(kp *KeyPair, senderPub ed25519.PublicKey, compact string) (*NestedClaims, error) {
	if kp == nil {
		return nil, ErrNilKeyPair
	}
	if !kp.IsReady() {
		return nil, fmt.Errorf("%w: key pair is incomplete", ErrInvalidKey)
	}
	if senderPub == nil {
		senderPub = kp.SigningPublic
	}

	decrypted, err := jwe.Decrypt(
		[]byte(compact),
		jwe.WithKey(jwa.ECDH_ES_A256KW(), kp.KEXPrivate),
	)
	if err != nil {
		return nil, fmt.Errorf("jwe decrypt: %w: %v", ErrDecryptionFailed, err)
	}

	verified, err := jws.Verify(decrypted, jws.WithKey(jwa.EdDSAEd25519(), senderPub))
	if err != nil {
		return nil, fmt.Errorf("jws verify: %w: %v", ErrVerificationFailed, err)
	}

	tok, err := jwt.Parse(verified, jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		return nil, fmt.Errorf("parse inner jwt: %w", err)
	}

	claims := &NestedClaims{}
	var s string
	if err := tok.Get("user_id", &s); err == nil {
		claims.UserID = s
	}
	if err := tok.Get("email", &s); err == nil {
		claims.Email = s
	}
	if err := tok.Get("role", &s); err == nil {
		claims.Role = s
	}
	if err := tok.Get("session_id", &s); err == nil {
		claims.SessionID = s
	}
	if exp, ok := tok.Expiration(); ok {
		claims.ExpiresAt = exp
	}
	if iat, ok := tok.IssuedAt(); ok {
		claims.IssuedAt = iat
	}
	if jti, ok := tok.JwtID(); ok {
		claims.JWTID = jti
	}
	return claims, nil
}

// SealForCookie is a convenience wrapper for HTTP-only cookie storage.
func SealForCookie(kp *KeyPair, claims *NestedClaims, ttl time.Duration) (string, error) {
	if kp == nil {
		return "", ErrNilKeyPair
	}
	return Seal(kp, kp.KEXPublic, claims, ttl)
}

// OpenFromCookie decrypts a cookie-stored nested JWT using server's key pair.
func OpenFromCookie(kp *KeyPair, cookieValue string) (*NestedClaims, error) {
	if kp == nil {
		return nil, ErrNilKeyPair
	}
	return Open(kp, kp.SigningPublic, cookieValue)
}

// ---------- Cookie helpers (unchanged) ----------

type CookieConfig struct {
	Name     string
	Domain   string
	Path     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	SameSite string
}

func DefaultCookieConfig() *CookieConfig {
	return &CookieConfig{
		Name: "__session", Path: "/", MaxAge: 900,
		Secure: true, HTTPOnly: true, SameSite: "Strict",
	}
}

func SetCookie(cfg *CookieConfig, token string) string {
	cookie := fmt.Sprintf("%s=%s; Path=%s; MaxAge=%d", cfg.Name, token, cfg.Path, cfg.MaxAge)
	if cfg.Domain != "" {
		cookie += "; Domain=" + cfg.Domain
	}
	if cfg.Secure {
		cookie += "; Secure"
	}
	if cfg.HTTPOnly {
		cookie += "; HttpOnly"
	}
	if cfg.SameSite != "" {
		cookie += "; SameSite=" + cfg.SameSite
	}
	return cookie
}

func ClearCookie(cfg *CookieConfig) string {
	return fmt.Sprintf("%s=; Path=%s; MaxAge=0; Expires=Thu, 01 Jan 1970 00:00:00 GMT",
		cfg.Name, cfg.Path)
}
