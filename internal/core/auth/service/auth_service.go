package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	authdto "github.com/vadxq/go-rest-starter/internal/core/auth/dto"
	userrepo "github.com/vadxq/go-rest-starter/internal/core/user/repository"
	"github.com/vadxq/go-rest-starter/pkg/cache"
	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
	"github.com/vadxq/go-rest-starter/pkg/jwt"
)

const (
	// Token cache key prefix
	tokenCachePrefix = "token:"

	// Token blacklist cache key prefix
	tokenBlacklistPrefix = "blacklist:"
)

// AuthService interface
type AuthService interface {
	Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*authdto.TokenResponse, error)
	Logout(ctx context.Context, accessToken string) error
}

// authService implements AuthService
type authService struct {
	userRepo  userrepo.UserRepository
	validator *validator.Validate
	db        *gorm.DB
	jwtConfig *jwt.Config
	cache     cache.Cache
}

// NewAuthService creates a new auth service
func NewAuthService(ur userrepo.UserRepository, v *validator.Validate, db *gorm.DB, jwtConfig *jwt.Config, c cache.Cache) AuthService {
	return &authService{
		userRepo:  ur,
		validator: v,
		db:        db,
		jwtConfig: jwtConfig,
		cache:     c,
	}
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, apperrors.ValidationError("การตรวจสอบข้อมูลล้มเหลว", err)
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Return same error for both not found and db errors to prevent enumeration
		return nil, apperrors.UnauthorizedError("อีเมลหรือรหัสผ่านไม่ถูกต้อง", nil)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.UnauthorizedError("อีเมลหรือรหัสผ่านไม่ถูกต้อง", nil)
	}

	// Generate access token
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Role, s.jwtConfig)
	if err != nil {
		return nil, apperrors.InternalError("สร้าง access token ไม่สำเร็จ", err)
	}

	// Generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperrors.InternalError("สร้าง refresh token ไม่สำเร็จ", err)
	}

	// Cache tokens for quick validation/tracking
	tokenKey := fmt.Sprintf("%s%d", tokenCachePrefix, user.ID)
	if s.cache != nil {
		_ = s.cache.SetObject(ctx, tokenKey, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}, s.jwtConfig.AccessTokenExp)
	}

	return &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessTokenExp.Seconds()),
		TokenType:    "Bearer",
		User: authdto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// RefreshToken issues a new access token using a refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*authdto.TokenResponse, error) {
	// Check if token is blacklisted
	blacklistKey := fmt.Sprintf("%s%s", tokenBlacklistPrefix, refreshToken)
	var blacklisted bool
	if s.cache != nil {
		err := s.cache.GetObject(ctx, blacklistKey, &blacklisted)
		if err == nil && blacklisted {
			return nil, apperrors.UnauthorizedError("refresh token ถูกเพิกถอนแล้ว", nil)
		}
	}

	// Parse refresh token
	userId, err := jwt.ParseRefreshToken(refreshToken, s.jwtConfig.Secret)
	if err != nil {
		return nil, apperrors.UnauthorizedError("refresh token ไม่ถูกต้อง", nil)
	}

	// Convert user ID to string
	userIdStr := fmt.Sprintf("%d", userId)

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userIdStr)
	if err != nil {
		return nil, apperrors.UnauthorizedError("ไม่มีผู้ใช้", nil)
	}

	// Generate new access token
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Role, s.jwtConfig)
	if err != nil {
		return nil, apperrors.InternalError("สร้าง access token ไม่สำเร็จ", err)
	}

	// Update cache
	tokenKey := fmt.Sprintf("%s%d", tokenCachePrefix, user.ID)
	if s.cache != nil {
		_ = s.cache.SetObject(ctx, tokenKey, map[string]string{
			"access_token": accessToken,
		}, s.jwtConfig.AccessTokenExp)
	}

	return &authdto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(s.jwtConfig.AccessTokenExp.Seconds()),
		TokenType:   "Bearer",
	}, nil
}

// Logout invalidates the access token
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	// Parse token to get user ID
	claims, err := jwt.ParseToken(accessToken, s.jwtConfig.Secret)
	if err != nil {
		return apperrors.UnauthorizedError("access token ไม่ถูกต้อง", nil)
	}

	// Add token to blacklist
	if s.cache != nil {
		blacklistKey := fmt.Sprintf("%s%s", tokenBlacklistPrefix, accessToken)
		_ = s.cache.SetObject(ctx, blacklistKey, true, s.jwtConfig.AccessTokenExp)

		// Clear user token cache
		tokenKey := fmt.Sprintf("%s%d", tokenCachePrefix, claims.UserID)
		_ = s.cache.Delete(ctx, tokenKey)
	}

	return nil
}