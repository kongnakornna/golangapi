package middleware

import (
	"net/http"

	"icmongolang/config"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/logger"
)

type MiddlewareManager struct {
	cfg       *config.Config
	logger    logger.Logger
	usersUC   users.UserUseCaseI
	rateLimit *RateLimitMiddleware
}

func CreateMiddlewareManager(cfg *config.Config, logger logger.Logger, usersUC users.UserUseCaseI) *MiddlewareManager {
	rateLimitMiddleware := NewRateLimitMiddleware(DefaultRateLimitConfig)

	return &MiddlewareManager{
		cfg:       cfg,
		logger:    logger,
		usersUC:   usersUC,
		rateLimit: rateLimitMiddleware,
	}
}

// RateLimit returns the rate limiting middleware handler
func (m *MiddlewareManager) RateLimit() func(http.Handler) http.Handler {
	return m.rateLimit.Handler
}
