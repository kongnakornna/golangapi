package jwe

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSealOpen_RoundTrip(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	claims := &NestedClaims{
		UserID:    uuid.New().String(),
		Email:     "test@example.com",
		Role:      "admin",
		SessionID: uuid.New().String(),
	}

	compact, err := SealForCookie(kp, claims, 15*time.Minute)
	if err != nil {
		t.Fatalf("SealForCookie: %v", err)
	}

	opened, err := OpenFromCookie(kp, compact)
	if err != nil {
		t.Fatalf("OpenFromCookie: %v", err)
	}

	if opened.UserID != claims.UserID {
		t.Fatalf("user_id mismatch: got %s, want %s", opened.UserID, claims.UserID)
	}
	if opened.Email != claims.Email {
		t.Fatalf("email mismatch: got %s, want %s", opened.Email, claims.Email)
	}
	if opened.Role != claims.Role {
		t.Fatalf("role mismatch: got %s, want %s", opened.Role, claims.Role)
	}
	if opened.SessionID != claims.SessionID {
		t.Fatalf("session_id mismatch: got %s, want %s", opened.SessionID, claims.SessionID)
	}
	if opened.JWTID == "" {
		t.Fatal("jti should be set")
	}
	if opened.ExpiresAt.IsZero() {
		t.Fatal("exp should be set")
	}
	if opened.IssuedAt.IsZero() {
		t.Fatal("iat should be set")
	}
}

func TestSealOpen_CrossKeyPair(t *testing.T) {
	// Sender and recipient have different key pairs
	senderKP, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("sender GenerateKeyPair: %v", err)
	}
	recipientKP, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("recipient GenerateKeyPair: %v", err)
	}

	claims := &NestedClaims{
		UserID:    uuid.New().String(),
		Email:     "cross@example.com",
		Role:      "user",
		SessionID: uuid.New().String(),
	}

	// Sender seals for recipient
	compact, err := Seal(senderKP, recipientKP.KEXPublic, claims, 15*time.Minute)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}

	// Recipient opens using sender's public key for JWS verification
	opened, err := Open(recipientKP, senderKP.SigningPublic, compact)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if opened.UserID != claims.UserID {
		t.Fatalf("user_id mismatch")
	}
}

func TestSealOpen_Expired(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	claims := &NestedClaims{
		UserID:    uuid.New().String(),
		Email:     "expired@example.com",
		Role:      "user",
		SessionID: uuid.New().String(),
	}

	// Seal with negative TTL (already expired)
	compact, err := SealForCookie(kp, claims, -time.Hour)
	if err != nil {
		t.Fatalf("SealForCookie: %v", err)
	}

	// Should still decrypt but claims will be expired
	opened, err := OpenFromCookie(kp, compact)
	if err != nil {
		t.Fatalf("OpenFromCookie: %v", err)
	}

	if time.Now().After(opened.ExpiresAt) {
		// This is expected - token is expired
		t.Log("token correctly expired")
	}
}

func TestLoadKeyPair_RoundTrip(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	// Encode to base64
	signingPrivB64 := base64.StdEncoding.EncodeToString(kp.SigningPrivate)
	kexPrivBytes := kp.KEXPrivate.Bytes()
	kexPrivB64 := base64.StdEncoding.EncodeToString(kexPrivBytes)

	// Load back
	loaded, err := LoadKeyPair(signingPrivB64, kexPrivB64)
	if err != nil {
		t.Fatalf("LoadKeyPair: %v", err)
	}

	// Verify keys match
	if !ed25519PrivEqual(loaded.SigningPrivate, kp.SigningPrivate) {
		t.Fatal("signing private key mismatch")
	}
	if !ed25519PubEqual(loaded.SigningPublic, kp.SigningPublic) {
		t.Fatal("signing public key mismatch")
	}
	if !loaded.KEXPublic.Equal(kp.KEXPublic) {
		t.Fatal("kex public key mismatch")
	}
}

func ed25519PrivEqual(a, b ed25519.PrivateKey) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ed25519PubEqual(a, b ed25519.PublicKey) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCookieConfig(t *testing.T) {
	cfg := DefaultCookieConfig()
	token := "test-token"
	cookie := SetCookie(cfg, token)

	if cookie == "" {
		t.Fatal("cookie should not be empty")
	}

	clearCookie := ClearCookie(cfg)
	if clearCookie == "" {
		t.Fatal("clear cookie should not be empty")
	}
}
