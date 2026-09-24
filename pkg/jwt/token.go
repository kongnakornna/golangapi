package jwt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"icmongolang/pkg/httpErrors"

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
