package middleware

import (
	"context"
	"errors"
	"net/http"

	"icmongolang/config"
	"icmongolang/internal"
	"icmongolang/internal/modules/auth/domain"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/apikey"
	"icmongolang/pkg/jwe"
	"icmongolang/pkg/logger"
)

type MiddlewareManager struct {
	cfg          *config.Config
	logger       logger.Logger
	usersUC      users.UserUseCaseI
	rateLimit    *RateLimitMiddleware
	jweKP        *jwe.KeyPair
	apiKeyPepper []byte
	redisRepo    internal.RedisRepository[any]
}

func CreateMiddlewareManager(
	cfg *config.Config,
	logger logger.Logger,
	usersUC users.UserUseCaseI,
	redisRepo internal.RedisRepository[any],
) *MiddlewareManager {
	rateLimitMiddleware := NewRateLimitMiddleware(DefaultRateLimitConfig)

	// ── JWE key pair ──
	var jweKP *jwe.KeyPair
	if cfg.Jwt.JweSigningPrivateKey != "" && cfg.Jwt.JweKexPrivateKey != "" {
		kp, err := jwe.LoadKeyPair(cfg.Jwt.JweSigningPrivateKey, cfg.Jwt.JweKexPrivateKey)
		if err != nil {
			logger.Warnf("Failed to load JWE keys for middleware: %v", err)
		} else {
			jweKP = kp
			logger.Info("JWE key pair loaded for middleware")
		}
	} else {
		logger.Warn("JWE keys not configured — JWE auth disabled")
	}

	// ── API key pepper ──
	apiKeyPepper := []byte(cfg.Jwt.ApiKeyPepper)
	if len(apiKeyPepper) == 0 {
		logger.Warn("API key pepper not configured — token revocation checking disabled")
	}

	return &MiddlewareManager{
		cfg:          cfg,
		logger:       logger,
		usersUC:      usersUC,
		rateLimit:    rateLimitMiddleware,
		jweKP:        jweKP,
		apiKeyPepper: apiKeyPepper,
		redisRepo:    redisRepo,
	}
}

// RateLimit handler
func (m *MiddlewareManager) RateLimit() func(http.Handler) http.Handler {
	return m.rateLimit.Handler
}

// JWEKeyPair สำหรับ WebSocket
func (m *MiddlewareManager) JWEKeyPair() *jwe.KeyPair {
	return m.jweKP
}

// AuthenticateWS — ใช้ตอน WS handshake
func (m *MiddlewareManager) AuthenticateWS(ctx context.Context, token string) (*Claims, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}
	if m.jweKP == nil {
		return nil, errors.New("JWE keys not configured")
	}

	raw, err := jwe.OpenFromCookie(m.jweKP, token)
	if err != nil {
		return nil, err
	}

	// user-level revoke
	if m.redisRepo != nil {
		key := "revoke:user:" + raw.UserID
		if v, err := m.redisRepo.GetString(ctx, key); err == nil && v != "" {
			return nil, errors.New("user sessions revoked")
		}
	}

	// token-level revoke
	if len(m.apiKeyPepper) > 0 {
		fp := apikey.Hash(m.apiKeyPepper, token)
		key := "revoke:refresh:" + fp
		if v, err := m.redisRepo.GetString(ctx, key); err == nil && v != "" {
			return nil, errors.New("token revoked")
		}
	}

	role := domain.Role(raw.Role)
	if !role.IsValid() {
		role = domain.RoleUser
	}

	return &Claims{
		ID:        raw.UserID,
		UserID:    raw.UserID,
		Email:     raw.Email,
		Role:      role,
		SessionID: raw.SessionID,
	}, nil
}