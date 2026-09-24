// internal/modules/auth/domain/errors/errors.go

package errors

import "errors"

var (
	// User errors.
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyActive    = errors.New("user already active")
	ErrUserAlreadyInactive  = errors.New("user already inactive")
	ErrUserAlreadySuspended = errors.New("user already suspended")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidEmail         = errors.New("invalid email format")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrAccountInactive      = errors.New("account is inactive")
	ErrAccountSuspended     = errors.New("account is suspended")

	// Auth errors.
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenRevoked = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenAlreadyUsed    = errors.New("token already used")
	ErrTokenExpired        = errors.New("token expired")

	// Permission errors.
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidRole      = errors.New("invalid role")

	// Session errors.
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionRevoked  = errors.New("session revoked")

	// Validation errors.
	ErrValidationFailed = errors.New("validation failed")
	ErrRequiredField    = errors.New("required field missing")
)