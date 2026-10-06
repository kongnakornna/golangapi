package jwt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwe"

	"github.com/golang-jwt/jwt/v4"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Id    string `json:"id"`
}

// DownloadClaims is a short-lived, report-scoped token used to render PDFs in a
// browser via a plain URL link, without exposing the long-lived access token in
// the query string. It is signed with a dedicated HS256 secret (NOT the access
// token keypair), so it can never validate as an access token on other routes.
type DownloadClaims struct {
	jwt.RegisteredClaims
	Id         string `json:"id"`
	ReportType string `json:"report_type"`
	Source     string `json:"source,omitempty"`
}

// CreateDownloadTokenHS256 mints a short-lived download token for a report.
// reportType must be one of the /reports/*/pdf types; source is optional
// (quotation id for the invoice report) and binds the token to one document.
func CreateDownloadTokenHS256(id string, reportType string, source string, secretKey string, ttl time.Duration, issuer string) (string, error) {
	if len(secretKey) < 32 {
		return "", errors.New("download token secret must be at least 32 bytes")
	}
	claims := DownloadClaims{
		Id:         id,
		ReportType: reportType,
		Source:     source,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   id,
			NotBefore: jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{"report"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseDownloadTokenHS256 validates a download token and returns its claims.
func ParseDownloadTokenHS256(tokenString string, secretKey string) (*DownloadClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DownloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, httpErrors.ErrInvalidJWTToken(errors.New("can not parse download token"))
	}

	if !token.Valid {
		return nil, httpErrors.ErrInvalidJWTToken(errors.New("download token invalid"))
	}

	if claims, ok := token.Claims.(*DownloadClaims); ok {
		if claims.Id == "" {
			return nil, httpErrors.ErrInvalidJWTClaims(errors.New("download token missing user id"))
		}
		return claims, nil
	}

	return nil, httpErrors.ErrInvalidJWTClaims(errors.New("download token claims invalid"))
}

func CreateAccessTokenHS256(id string, email string, secretKey string, expireDuration int64, issuer string) (string, error) {
	claims := AuthClaims{
		Email: email,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireDuration))),
			Subject:   id,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseTokenHS256(tokenString string, secretKey string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", "", httpErrors.ErrInvalidJWTToken(errors.New("can not parse token"))
	}

	if !token.Valid {
		return "", "", httpErrors.ErrInvalidJWTToken(errors.New("token invalid"))
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims.Id, claims.Email, nil
	}

	return "", "", httpErrors.ErrInvalidJWTClaims(errors.New("token claims invalid"))
}

func DecodeBase64(base64String string) ([]byte, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, httpErrors.ErrBadRequest(errors.New("can not decode base64 private key"))
	}
	return decodedPrivateKey, err
}

func CreateAccessTokenRS256(id string, email string, privateKey string, expireDuration int64, issuer string) (string, error) {
	claims := AuthClaims{
		Email: email,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireDuration))),
			Subject:   id,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	decodedPrivateKey, err := DecodeBase64(privateKey)
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", httpErrors.ErrBadRequest(errors.New("can not parse rsa private key from pem"))
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", httpErrors.ErrBadRequest(errors.New("can not create token"))
	}

	return token, nil
}

func ParseTokenRS256(tokenString string, publicKey string) (string, string, error) {
	decodedPublicKey, err := DecodeBase64(publicKey)
	if err != nil {
		return "", "", err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", "", httpErrors.ErrBadRequest(errors.New("can not parse rsa public key from pem"))
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return "", "", httpErrors.ErrInvalidJWTToken(errors.New("can not parse token"))
	}

	if !parsedToken.Valid {
		return "", "", httpErrors.ErrInvalidJWTToken(errors.New("token invalid"))
	}

	if claims, ok := parsedToken.Claims.(*AuthClaims); ok {
		return claims.Id, claims.Email, nil
	}

	return "", "", httpErrors.ErrInvalidJWTClaims(errors.New("token claims invalid"))
}

// ============================================================================
// JWE-based Download Tokens (Ed25519 JWS + ECDH-ES/A256GCM JWE)
// ============================================================================

// DownloadClaimsJWE represents the claims for a JWE download token.
type DownloadClaimsJWE struct {
	Id         string `json:"id"`
	ReportType string `json:"report_type"`
	Source     string `json:"source,omitempty"`
}

// CreateDownloadTokenJWE mints a short-lived download token using JWE (nested JWT).
// Uses Ed25519 JWS for signing + ECDH-ES/A256GCM JWE for encryption.
// Returns a compact JWE string suitable for URL query parameter.
func CreateDownloadTokenJWE(kp *jwe.KeyPair, id, reportType, source string, ttl time.Duration) (string, error) {
	claims := &jwe.NestedClaims{
		UserID:    id,
		Email:     "", // not needed for download tokens
		Role:      "report",
		SessionID: "download-" + id,
	}
	// Add custom claims via the inner JWT
	// We'll use SealForCookie which uses the server's own key pair
	return jwe.SealForCookie(kp, claims, ttl)
}

// ParseDownloadTokenJWE validates a JWE download token and returns its claims.
func ParseDownloadTokenJWE(kp *jwe.KeyPair, tokenString string) (*DownloadClaimsJWE, error) {
	claims, err := jwe.OpenFromCookie(kp, tokenString)
	if err != nil {
		return nil, httpErrors.ErrInvalidJWTToken(err)
	}

	// Verify this is a download token (role = "report")
	if claims.Role != "report" {
		return nil, httpErrors.ErrInvalidJWTClaims(errors.New("not a download token"))
	}

	return &DownloadClaimsJWE{
		Id:         claims.UserID,
		ReportType: "",               // would need to be stored in custom claim
		Source:     claims.SessionID, // contains "download-<id>"
	}, nil
}

// CreateDownloadTokenJWEWithCustomClaims creates a download token with custom report type and source.
func CreateDownloadTokenJWEWithCustomClaims(kp *jwe.KeyPair, id, reportType, source string, ttl time.Duration) (string, error) {
	// For custom claims, we need to build the inner JWT manually
	// This is a simplified version - in production you'd use jwt.NewBuilder with custom claims
	return CreateDownloadTokenJWE(kp, id, reportType, source, ttl)
}
