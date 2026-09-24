package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testDownloadSecret = "test-download-secret-at-least-32-bytes-long!!"

// rsaKeysB64 generates an RSA-2048 key pair and returns base64-encoded PEM strings,
// matching the format used in config files.
func rsaKeysB64(t *testing.T) (privB64, pubB64 string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("MarshalPKIXPublicKey: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	return base64.StdEncoding.EncodeToString(privPEM), base64.StdEncoding.EncodeToString(pubPEM)
}

func TestCreateDownloadTokenHS256_ShortSecret_Errors(t *testing.T) {
	_, err := CreateDownloadTokenHS256("user-1", "invoice", "src-1", "too-short", time.Minute, "issuer")
	assert.Error(t, err)
}

func TestCreateAndParseDownloadTokenHS256_RoundTrip(t *testing.T) {
	token, err := CreateDownloadTokenHS256("user-1", "invoice", "src-1", testDownloadSecret, 5*time.Minute, "test-issuer")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseDownloadTokenHS256(token, testDownloadSecret)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "user-1", claims.Id)
	assert.Equal(t, "invoice", claims.ReportType)
	assert.Equal(t, "src-1", claims.Source)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.Equal(t, "report", string(claims.Audience[0]))
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestParseDownloadTokenHS256_WrongSecret_Fails(t *testing.T) {
	token, err := CreateDownloadTokenHS256("user-1", "daily_sales", "", testDownloadSecret, time.Minute, "issuer")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ParseDownloadTokenHS256(token, "a-different-secret-that-is-long-enough-ok")
	assert.Error(t, err)
}

func TestParseDownloadTokenHS256_Expired_Fails(t *testing.T) {
	token, err := CreateDownloadTokenHS256("user-1", "invoice", "src-1", testDownloadSecret, -time.Minute, "issuer")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ParseDownloadTokenHS256(token, testDownloadSecret)
	assert.Error(t, err)
}

// RS256 tests — regression: empty PEM key previously caused 400 on login.

func TestCreateAndParseAccessTokenRS256_RoundTrip(t *testing.T) {
	privB64, pubB64 := rsaKeysB64(t)

	token, err := CreateAccessTokenRS256("user-42", "test@example.com", privB64, int64(time.Hour), "test-issuer")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, email, err := ParseTokenRS256(token, pubB64)
	assert.NoError(t, err)
	assert.Equal(t, "user-42", id)
	assert.Equal(t, "test@example.com", email)
}

func TestCreateAccessTokenRS256_EmptyPrivateKey_Fails(t *testing.T) {
	_, err := CreateAccessTokenRS256("user-1", "a@b.com", "", 3600, "issuer")
	assert.Error(t, err)
}

func TestCreateAccessTokenRS256_InvalidBase64_Fails(t *testing.T) {
	_, err := CreateAccessTokenRS256("user-1", "a@b.com", "!!!not-base64!!!", 3600, "issuer")
	assert.Error(t, err)
}

func TestParseTokenRS256_EmptyPublicKey_Fails(t *testing.T) {
	privB64, _ := rsaKeysB64(t)
	token, err := CreateAccessTokenRS256("user-1", "a@b.com", privB64, 3600, "issuer")
	assert.NoError(t, err)

	_, _, err = ParseTokenRS256(token, "")
	assert.Error(t, err)
}

func TestParseTokenRS256_WrongPublicKey_Fails(t *testing.T) {
	privB64, _ := rsaKeysB64(t)
	_, wrongPubB64 := rsaKeysB64(t)

	token, err := CreateAccessTokenRS256("user-1", "a@b.com", privB64, 3600, "issuer")
	assert.NoError(t, err)

	_, _, err = ParseTokenRS256(token, wrongPubB64)
	assert.Error(t, err)
}

func TestParseTokenRS256_ExpiredToken_Fails(t *testing.T) {
	privB64, pubB64 := rsaKeysB64(t)

	token, err := CreateAccessTokenRS256("user-1", "a@b.com", privB64, -int64(time.Hour), "issuer")
	assert.NoError(t, err)

	_, _, err = ParseTokenRS256(token, pubB64)
	assert.Error(t, err)
}
