package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenFromHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo", nil)
	r.Header.Set("Authorization", "Bearer abc123")

	assert.Equal(t, "abc123", TokenFromHeader(r))
}

func TestTokenFromHeader_NoToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo", nil)
	assert.Equal(t, "", TokenFromHeader(r))
}

func TestTokenFromRequest_Header(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo?access_token=query-token&token=other-token", nil)
	r.Header.Set("Authorization", "Bearer header-token")
	r.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: "cookie-token"})

	assert.Equal(t, "header-token", TokenFromRequest(r))
}

func TestTokenFromRequest_Cookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo?access_token=query-token", nil)
	r.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: "cookie-token"})

	assert.Equal(t, "cookie-token", TokenFromRequest(r))
}

func TestTokenFromRequest_QueryAccessToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo?access_token=query-token", nil)

	assert.Equal(t, "query-token", TokenFromRequest(r))
}

func TestTokenFromRequest_QueryToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo?token=query-token", nil)

	assert.Equal(t, "query-token", TokenFromRequest(r))
}

func TestTokenFromRequest_None(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/foo", nil)

	assert.Equal(t, "", TokenFromRequest(r))
}

func TestTokenFromHeaderOrCookie(t *testing.T) {
	t.Run("header wins over cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/foo", nil)
		r.Header.Set("Authorization", "Bearer header-token")
		r.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: "cookie-token"})

		assert.Equal(t, "header-token", tokenFromHeaderOrCookie(r))
	})

	t.Run("cookie used when no header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/foo", nil)
		r.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: "cookie-token"})

		assert.Equal(t, "cookie-token", tokenFromHeaderOrCookie(r))
	})

	t.Run("query param is ignored", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/foo?token=query-token", nil)

		assert.Equal(t, "", tokenFromHeaderOrCookie(r))
	})
}
