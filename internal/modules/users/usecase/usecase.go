package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/users"
	"icmongolang/internal/usecase"
	"icmongolang/internal/worker"
	"icmongolang/pkg/cryptpass"
	"icmongolang/pkg/emailTemplates"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwe"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/secureRandom"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// userUseCase implements UserUseCaseI interface.
type userUseCase struct {
	usecase.UseCase[models.SdUser]
	pgRepo                 users.UserPgRepository
	redisRepo              users.UserRedisRepository
	emailTemplateGenerator emailTemplates.EmailTemplatesGenerator
	redisTaskDistributor   users.UserRedisTaskDistributor
	logger                 logger.Logger
	jweKP                  *jwe.KeyPair
	apiKeyPepper           []byte
	useJWE                 bool
}
// CreateUserUseCaseI creates new user use case instance.
func CreateUserUseCaseI(
	pgRepo users.UserPgRepository,
	redisRepo users.UserRedisRepository,
	redisTaskDistributor users.UserRedisTaskDistributor,
	cfg *config.Config,
	logger logger.Logger,
) users.UserUseCaseI {
	logger.Info("Initializing user use case")

	// ---- JWE keys (optional) ----
	var jweKP *jwe.KeyPair
	useJWE := false

	if cfg.Jwt.JweSigningPrivateKey != "" && cfg.Jwt.JweKexPrivateKey != "" {
		kp, err := jwe.LoadKeyPair(cfg.Jwt.JweSigningPrivateKey, cfg.Jwt.JweKexPrivateKey)
		if err != nil {
			logger.Warnf("Failed to load JWE keys, falling back to RS256: %v", err)
		} else {
			jweKP = kp
			useJWE = true
			logger.Info("JWE key pair loaded for nested JWT")
		}
	} else {
		logger.Warn("JWE keys not configured — using RS256 for tokens")
	}

	// ---- API key pepper (optional) ----
	apiKeyPepper := []byte(cfg.Jwt.ApiKeyPepper)
	if len(apiKeyPepper) == 0 {
		logger.Warn("API key pepper not configured, token fingerprinting will use plain SHA-256")
	}

	return &userUseCase{
		UseCase:                usecase.CreateUseCase[models.SdUser](pgRepo, cfg, logger),
		pgRepo:                 pgRepo,
		redisRepo:              redisRepo,
		emailTemplateGenerator: emailTemplates.NewEmailTemplatesGenerator(cfg),
		redisTaskDistributor:   redisTaskDistributor,
		logger:                 logger,
		jweKP:                  jweKP,
		apiKeyPepper:           apiKeyPepper,
		useJWE:                 useJWE,
	}
}

// refreshTokenFingerprint returns a short, deterministic fingerprint of the
// refresh token (HMAC-SHA256 with pepper when available, SHA-256 otherwise).
// Used as the Redis revocation key.
func (u *userUseCase) refreshTokenFingerprint(token string) string {
	if len(u.apiKeyPepper) > 0 {
		mac := hmac.New(sha256.New, u.apiKeyPepper)
		mac.Write([]byte(token))
		return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	}
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// storeRefreshToken stores the refresh token in the Redis set and applies TTL.
func (u *userUseCase) storeRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	key := u.GenerateRedisRefreshTokenKey(id)
	if err := u.redisRepo.Sadd(ctx, key, refreshToken); err != nil {
		return err
	}
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60
	if ttlSeconds <= 0 {
		ttlSeconds = 24 * 3600
	}
	if err := u.redisRepo.Expire(ctx, key, ttlSeconds); err != nil {
		u.logger.Warnf("Failed to set TTL on refresh token set %s: %v", key, err)
	}
	return nil
}

// Get retrieves user by ID (with Redis cache).
func (u *userUseCase) Get(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	u.logger.Infof("Getting user by ID: %s", id)
	cacheKey := u.GenerateRedisUserKey(id)

	cachedUser, err := u.redisRepo.Get(ctx, cacheKey)
	if err != nil {
		u.logger.Errorf("Redis get error for key %s: %v", cacheKey, err)
		return nil, err
	}
	if cachedUser != nil {
		u.logger.Debugf("Cache hit for user %s", id)
		return cachedUser, nil
	}

	user, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		u.logger.Errorf("Database get error for user %s: %v", id, err)
		return nil, err
	}

	if err = u.redisRepo.Create(ctx, cacheKey, user, 3600); err != nil {
		u.logger.Warnf("Failed to cache user %s: %v", id, err)
	}
	return user, nil
}

// Delete removes user from DB and clears Redis cache.
func (u *userUseCase) Delete(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	u.logger.Infof("Deleting user: %s", id)
	user, err := u.pgRepo.Delete(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to delete user %s from DB: %v", id, err)
		return nil, err
	}
	u.invalidateUserCaches(ctx, id)
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id))
	u.logger.Infof("User %s deleted successfully", id)
	return user, nil
}

// Update updates user fields and invalidates cache.
func (u *userUseCase) Update(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*models.SdUser, error) {
	u.logger.Infof("Updating user %s with fields: %v", id, keys(values))
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		u.logger.Errorf("User %s not found for update: %v", id, err)
		return nil, err
	}

	user, err := u.pgRepo.Update(ctx, obj, values)
	if err != nil {
		u.logger.Errorf("Database update failed for user %s: %v", id, err)
		return nil, err
	}
	u.invalidateUserCaches(ctx, id)
	u.logger.Infof("User %s updated successfully", id)
	return user, nil
}

// Create – inserts a new user and sends verification email.
func (u *userUseCase) Create(ctx context.Context, exp *models.SdUser) (*models.SdUser, error) {
	u.logger.Info("Creating new user")

	exp.Email = strings.ToLower(strings.TrimSpace(exp.Email))
	exp.Password = strings.TrimSpace(exp.Password)

	if existing, err := u.pgRepo.GetByEmailStrict(ctx, exp.Email); err == nil && existing != nil {
		u.logger.Warnf("Email already registered: %s", exp.Email)
		return nil, httpErrors.ErrExistsEmail(errors.New("email already registered"))
	}
	effectiveUsername := exp.Username
	if effectiveUsername == "" {
		effectiveUsername = exp.Email
	}
	if existing, err := u.pgRepo.GetByUsernameStrict(ctx, effectiveUsername); err == nil && existing != nil {
		u.logger.Warnf("Username already taken: %s", effectiveUsername)
		return nil, httpErrors.ErrExistsUsername(errors.New("username already taken"))
	}

	hashedPassword, err := cryptpass.HashPassword(exp.Password)
	if err != nil {
		u.logger.Errorf("Password hashing failed for %s: %v", exp.Email, err)
		return nil, err
	}
	exp.Password = hashedPassword
	exp.PasswordTemp = nil

	if exp.Username == "" {
		exp.Username = exp.Email
	}
	if exp.RoleID == 0 {
		exp.RoleID = 2
	}

	user, err := u.pgRepo.Create(ctx, exp)
	if err != nil {
		u.logger.Errorf("Database create failed for %s: %v", exp.Email, err)
		return nil, err
	}
	u.logger.Infof("User created with ID: %s, verified=%v", user.ID, user.Verified)

	if user.Verified {
		return user, nil
	}

	verificationCode, err := secureRandom.RandomHex(16)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.pgRepo.UpdateVerificationCode(ctx, user, verificationCode)
	if err != nil {
		return nil, err
	}

	name := ""
	if updatedUser.Fullname != nil {
		name = *updatedUser.Fullname
	} else {
		name = updatedUser.Email
	}

	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GenerateVerificationCodeTemplate(
		ctx,
		name,
		fmt.Sprintf("%s/auth/verifyemail?code=%s", strings.TrimRight(u.Cfg.Server.BaseUrl, "/"), verificationCode),
	)
	if err != nil {
		u.logger.Errorf("Failed to generate email template for %s: %v", updatedUser.Email, err)
		return nil, err
	}

	err = u.redisTaskDistributor.DistributeTaskSendEmail(ctx, &users.PayloadSendEmail{
		From:      u.Cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.Cfg.Email.VerificationSubject,
		BodyHtml:  bodyHtml,
		BodyPlain: bodyPlain,
	}, asynq.MaxRetry(10), asynq.ProcessIn(10*time.Second), asynq.Queue(worker.QueueCritical))
	if err != nil {
		u.logger.Errorf("Failed to enqueue verification email for %s: %v", updatedUser.Email, err)
		return nil, err
	}
	u.logger.Infof("Verification email task queued for %s", updatedUser.Email)
	return updatedUser, nil
}

// CreateUser – wrapper with password confirmation.
func (u *userUseCase) CreateUser(ctx context.Context, exp *models.SdUser, confirmPassword string) (*models.SdUser, error) {
	if exp.Password != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}
	return u.Create(ctx, exp)
}

// createToken generates access and refresh tokens.
// Uses JWE (nested JWT) when configured, otherwise falls back to RS256.
func (u *userUseCase) createToken(ctx context.Context, exp models.SdUser) (string, string, error) {
	u.logger.Debugf("Generating tokens for user: %s", exp.ID)

	role := "user"
	switch exp.RoleID {
	case 1:
		role = "admin"
	case 2:
		role = "manager"
	}

	// ---- Resolve TTLs ----
	accessTTL := time.Duration(u.Cfg.Jwt.AccessTokenExpireDuration) * time.Minute
	if u.Cfg.Jwt.JweAccessTokenTTL > 0 {
		accessTTL = time.Duration(u.Cfg.Jwt.JweAccessTokenTTL) * time.Minute
	}
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}

	refreshTTL := time.Duration(u.Cfg.Jwt.RefreshTokenExpireDuration) * time.Minute
	if u.Cfg.Jwt.JweRefreshTokenTTL > 0 {
		refreshTTL = time.Duration(u.Cfg.Jwt.JweRefreshTokenTTL) * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 24 * time.Hour
	}

	// ---------- JWE path ----------
	if u.useJWE && u.jweKP != nil {
		accessClaims := &jwe.NestedClaims{
			UserID:    exp.ID.String(),
			Email:     exp.Email,
			Role:      role,
			SessionID: uuid.New().String(),
		}
		accessToken, err := jwe.SealForCookie(u.jweKP, accessClaims, accessTTL)
		if err != nil {
			u.logger.Errorf("JWE access token generation failed for %s: %v", exp.ID, err)
			return "", "", err
		}

		refreshClaims := &jwe.NestedClaims{
			UserID:    exp.ID.String(),
			Email:     exp.Email,
			Role:      role,
			SessionID: accessClaims.SessionID,
		}
		refreshToken, err := jwe.SealForCookie(u.jweKP, refreshClaims, refreshTTL)
		if err != nil {
			u.logger.Errorf("JWE refresh token generation failed for %s: %v", exp.ID, err)
			return "", "", err
		}

		u.logger.Debugf("JWE tokens generated for user %s", exp.ID)
		return accessToken, refreshToken, nil
	}

	// ---------- RS256 fallback ----------
	u.logger.Warn("Falling back to RS256 token generation")

	accessToken, err := jwt.CreateAccessTokenRS256(
		exp.ID.String(),
		exp.Email,
		u.Cfg.Jwt.AccessTokenPrivateKey,
		int64(accessTTL),
		u.Cfg.Jwt.Issuer,
	)
	if err != nil {
		u.logger.Errorf("RS256 access token generation failed for %s: %v", exp.ID, err)
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(
		exp.ID.String(),
		exp.Email,
		u.Cfg.Jwt.RefreshTokenPrivateKey,
		int64(refreshTTL),
		u.Cfg.Jwt.Issuer,
	)
	if err != nil {
		u.logger.Errorf("RS256 refresh token generation failed for %s: %v", exp.ID, err)
		return "", "", err
	}

	u.logger.Debugf("RS256 tokens generated for user %s", exp.ID)
	return accessToken, refreshToken, nil
}

// SignIn authenticates user (via email) and returns tokens.
func (u *userUseCase) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	u.logger.Info("Sign in attempt")
	user, err := u.pgRepo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.Warn("Invalid credentials")
		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
	}

	if !cryptpass.ComparePassword(password, user.Password) {
		u.logger.Warn("Invalid credentials")
		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		u.logger.Error("Token creation failed")
		return "", "", err
	}

	if err = u.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		u.logger.Errorf("Failed to store refresh token for user %s: %v", user.ID, err)
		return "", "", err
	}
	u.logger.Info("Sign in successful")
	return accessToken, refreshToken, nil
}

// SignInByUsername authenticates user (via username or email).
func (u *userUseCase) SignInByUsername(ctx context.Context, username, password string) (string, string, error) {
	u.logger.Info("Sign in attempt")
	user, err := u.pgRepo.GetByEmail(ctx, username)
	if err != nil {
		u.logger.Warn("User not found")
		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
	}

	if !cryptpass.ComparePassword(password, user.Password) {
		u.logger.Warn("Invalid password")
		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		u.logger.Error("Token creation failed")
		return "", "", err
	}

	if err = u.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		u.logger.Errorf("Failed to store refresh token for user %s: %v", user.ID, err)
		return "", "", err
	}
	u.logger.Info("Sign in successful")
	return accessToken, refreshToken, nil
}

// IsActive checks if user account is active.
func (u *userUseCase) IsActive(ctx context.Context, exp models.SdUser) bool {
	return exp.Status == 1
}

// IsSuper checks if user has superuser privileges.
func (u *userUseCase) IsSuper(ctx context.Context, exp models.SdUser) bool {
	return exp.IsSuperUser
}

// CreateSuperUserIfNotExist creates the first superuser from config.
func (u *userUseCase) CreateSuperUserIfNotExist(ctx context.Context) (bool, error) {
	u.logger.Info("Checking if superuser exists")
	user, err := u.pgRepo.GetByEmail(ctx, u.Cfg.FirstSuperUser.Email)
	if err != nil || user == nil {
		u.logger.Infof("Superuser not found, creating from config: %s", u.Cfg.FirstSuperUser.Email)
		fullname := u.Cfg.FirstSuperUser.Name
		_, err := u.Create(ctx, &models.SdUser{
			Fullname:    &fullname,
			Email:       u.Cfg.FirstSuperUser.Email,
			Username:    u.Cfg.FirstSuperUser.Email,
			Password:    u.Cfg.FirstSuperUser.Password,
			RoleID:      1,
			Status:      1,
			IsSuperUser: true,
			Verified:    true,
		})
		if err != nil {
			u.logger.Errorf("Failed to create superuser: %v", err)
			return false, err
		}
		u.logger.Info("Superuser created successfully")
		return true, nil
	}
	u.logger.Info("Superuser already exists")
	return false, nil
}

// UpdatePassword changes user password and invalidates all sessions.
func (u *userUseCase) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword, confirmPassword string) (*models.SdUser, error) {
	u.logger.Infof("Password change request for user: %s", id)
	if newPassword != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !cryptpass.ComparePassword(oldPassword, user.Password) {
		return nil, httpErrors.ErrWrongPassword(errors.New("old password and new password not same"))
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.pgRepo.UpdatePassword(ctx, user, hashedPassword)
	if err != nil {
		return nil, err
	}

	u.invalidateUserCaches(ctx, id)
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id))
	u.logger.Infof("Password changed successfully for user %s, all sessions invalidated", id)
	return updatedUser, nil
}

// ParseIdFromRefreshToken extracts user ID from a refresh token.
// Supports both JWE (5 segments) and RS256 (3 segments) token formats.
func (u *userUseCase) ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	u.logger.Debug("Parsing refresh token")
	if refreshToken == "" {
		return uuid.UUID{}, httpErrors.ErrInvalidJWTToken(errors.New("empty token"))
	}

	parts := strings.Split(refreshToken, ".")

	// ---------- JWE path (5 segments: header.key.iv.ciphertext.tag) ----------
	if len(parts) == 5 {
		if u.jweKP == nil {
			return uuid.UUID{}, httpErrors.ErrInvalidJWTToken(errors.New("JWE keys not configured"))
		}

		claims, err := jwe.OpenFromCookie(u.jweKP, refreshToken)
		if err != nil {
			u.logger.Errorf("Failed to parse JWE refresh token: %v", err)
			return uuid.UUID{}, err
		}

		// Revocation check
		fingerprint := u.refreshTokenFingerprint(refreshToken)
		revokeKey := fmt.Sprintf("revoke:refresh:%s", fingerprint)
		if revoked, _ := u.redisRepo.Get(ctx, revokeKey); revoked != nil {
			u.logger.Warnf("Refresh token revoked (fingerprint: %s)", fingerprint)
			return uuid.UUID{}, httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("refresh token revoked"))
		}

		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			return uuid.UUID{}, httpErrors.ErrInvalidJWTClaims(errors.New("invalid uuid in token"))
		}
		return id, nil
	}

	// ---------- RS256 fallback (3 segments) ----------
	if len(parts) == 3 {
		idStr, _, err := jwt.ParseTokenRS256(refreshToken, u.Cfg.Jwt.RefreshTokenPublicKey)
		if err != nil {
			u.logger.Errorf("Failed to parse RS256 refresh token: %v", err)
			return uuid.UUID{}, err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.UUID{}, httpErrors.ErrInvalidJWTClaims(errors.New("invalid uuid in token"))
		}
		return id, nil
	}

	return uuid.UUID{}, httpErrors.ErrInvalidJWTToken(errors.New("unknown token format"))
}

// Refresh issues new tokens using a valid refresh token (rotation).
func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	u.logger.Info("Token refresh request")

	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	// Verify the token is still in the user's active set (rotation check)
	setKey := u.GenerateRedisRefreshTokenKey(idParsed)
	isMember, err := u.redisRepo.SIsMember(ctx, setKey, refreshToken)
	if err != nil {
		u.logger.Errorf("Redis SIsMember error: %v", err)
		return "", "", err
	}
	if !isMember {
		u.logger.Warnf("Refresh token not found in active set for user %s (already used or revoked)", idParsed)
		return "", "", httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("refresh token not found or already used"))
	}

	// Remove old token (rotation)
	if err := u.redisRepo.Srem(ctx, setKey, refreshToken); err != nil {
		u.logger.Warnf("Failed to remove old refresh token: %v", err)
	}

	// Also mark fingerprint as revoked (belt-and-suspenders)
	fingerprint := u.refreshTokenFingerprint(refreshToken)
	revokeKey := fmt.Sprintf("revoke:refresh:%s", fingerprint)
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60
	if ttlSeconds <= 0 {
		ttlSeconds = 24 * 3600
	}
	if err := u.redisRepo.Set(ctx, revokeKey, "1", ttlSeconds); err != nil {
		u.logger.Warnf("Failed to mark old refresh token as revoked: %v", err)
	}

	user, err := u.Get(ctx, idParsed)
	if err != nil {
		u.logger.Errorf("User %s not found after token validation: %v", idParsed, err)
		return "", "", err
	}

	accessToken, newRefreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		return "", "", err
	}

	if err := u.storeRefreshToken(ctx, idParsed, newRefreshToken); err != nil {
		u.logger.Errorf("Failed to store new refresh token: %v", err)
		return "", "", err
	}

	u.logger.Infof("Tokens refreshed successfully for user %s", idParsed)
	return accessToken, newRefreshToken, nil
}

// Logout revokes the specific refresh token.
func (u *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	u.logger.Info("Logout request")
	if refreshToken == "" {
		return httpErrors.ErrInvalidJWTToken(errors.New("empty refresh token"))
	}

	// Try set-based removal first (works for both JWE & RS256)
	if id, err := u.ParseIdFromRefreshToken(ctx, refreshToken); err == nil {
		if err := u.redisRepo.Srem(ctx, u.GenerateRedisRefreshTokenKey(id), refreshToken); err != nil {
			u.logger.Warnf("Failed to remove refresh token from set: %v", err)
		}
	}

	// Belt-and-suspenders: mark fingerprint as revoked
	fingerprint := u.refreshTokenFingerprint(refreshToken)
	revokeKey := fmt.Sprintf("revoke:refresh:%s", fingerprint)
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60
	if ttlSeconds <= 0 {
		ttlSeconds = 24 * 3600
	}
	if err := u.redisRepo.Set(ctx, revokeKey, "1", ttlSeconds); err != nil {
		u.logger.Errorf("Failed to revoke refresh token: %v", err)
		return err
	}

	u.logger.Infof("Refresh token revoked (fingerprint: %s)", fingerprint)
	return nil
}

// LogoutAll revokes all refresh tokens for a user.
func (u *userUseCase) LogoutAll(ctx context.Context, id uuid.UUID) error {
	u.logger.Infof("Logout all sessions for user: %s", id)

	// Clear the entire set of refresh tokens for this user
	if err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		u.logger.Errorf("Failed to delete refresh token set for user %s: %v", id, err)
		return err
	}

	// Also mark user-level revocation flag
	revokeKey := fmt.Sprintf("revoke:user:%s", id.String())
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60
	if ttlSeconds <= 0 {
		ttlSeconds = 24 * 3600
	}
	if err := u.redisRepo.Set(ctx, revokeKey, "1", ttlSeconds); err != nil {
		u.logger.Warnf("Failed to set user-level revocation flag: %v", err)
	}

	u.invalidateUserCaches(ctx, id)
	u.logger.Infof("All sessions revoked for user %s", id)
	return nil
}

// Verify confirms user email using verification code.
func (u *userUseCase) Verify(ctx context.Context, verificationCode string) error {
	u.logger.Infof("Email verification attempt with code: %s", verificationCode)
	user, err := u.pgRepo.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return err
	}
	if user.Verified {
		return httpErrors.ErrUserAlreadyVerified(errors.New("user already verified"))
	}

	updatedUser, err := u.pgRepo.UpdateVerification(ctx, user, "", true)
	if err != nil {
		return err
	}
	u.invalidateUserCaches(ctx, updatedUser.ID)
	u.logger.Infof("Email verified successfully for user: %s", updatedUser.Email)
	return nil
}

// ForgotPassword sends a password reset email.
func (u *userUseCase) ForgotPassword(ctx context.Context, email string) error {
	u.logger.Infof("Forgot password request for email: %s", email)
	user, err := u.pgRepo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.Warnf("Email not found: %s", email)
		return httpErrors.ErrNotFound(err)
	}
	if !user.Verified {
		u.logger.Warnf("Unverified user attempted password reset: %s", email)
	}

	resetToken, err := secureRandom.RandomHex(16)
	if err != nil {
		return err
	}

	updatedUser, err := u.pgRepo.UpdatePasswordReset(ctx, user, resetToken,
		time.Now().In(helpers.GetTimeLocation()).Add(15*time.Minute))
	if err != nil {
		return err
	}
	u.invalidateUserCaches(ctx, updatedUser.ID)

	name := ""
	if updatedUser.Fullname != nil {
		name = *updatedUser.Fullname
	} else {
		name = updatedUser.Email
	}
	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GeneratePasswordResetTemplate(
		ctx,
		name,
		fmt.Sprintf("%s/auth/resetpassword?code=%s", strings.TrimRight(u.Cfg.Server.BaseUrl, "/"), resetToken),
	)
	if err != nil {
		return err
	}

	err = u.redisTaskDistributor.DistributeTaskSendEmail(ctx, &users.PayloadSendEmail{
		From:      u.Cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.Cfg.Email.ResetSubject,
		BodyHtml:  bodyHtml,
		BodyPlain: bodyPlain,
	}, asynq.MaxRetry(10), asynq.ProcessIn(10*time.Second), asynq.Queue(worker.QueueCritical))
	if err != nil {
		return err
	}
	u.logger.Infof("Password reset email sent to %s", email)
	return nil
}

// ResetPassword performs password reset using token.
func (u *userUseCase) ResetPassword(ctx context.Context, resetToken, newPassword, confirmPassword string) error {
	u.logger.Info("Password reset attempt")
	if newPassword != confirmPassword {
		return httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.pgRepo.GetByResetTokenResetAt(ctx, resetToken, time.Now().In(helpers.GetTimeLocation()))
	if err != nil {
		return err
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		return err
	}

	updatedUser, err := u.pgRepo.UpdatePasswordResetToken(ctx, user, hashedPassword, "")
	if err != nil {
		return err
	}
	u.invalidateUserCaches(ctx, updatedUser.ID)
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(updatedUser.ID))
	u.logger.Infof("Password reset successfully for user: %s", updatedUser.Email)
	return nil
}

// GenerateRedisUserKey returns Redis key for user cache.
func (u *userUseCase) GenerateRedisUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", models.SdUser{}.TableName(), id.String())
}

// GenerateRedisRefreshTokenKey returns Redis key for user's refresh token set.
func (u *userUseCase) GenerateRedisRefreshTokenKey(id uuid.UUID) string {
	return fmt.Sprintf("RefreshToken:%s", id.String())
}

// GenerateRedisProfileKey returns Redis key for cached user profile.
func (u *userUseCase) GenerateRedisProfileKey(id uuid.UUID) string {
	return fmt.Sprintf("sd_user_profile:%s", id.String())
}

// GetMultiWithTotal returns users with total count (supports filter and sort).
func (u *userUseCase) GetMultiWithTotal(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.SdUser, int64, error) {
	u.logger.Infof("GetMultiWithTotal called: limit=%d, offset=%d", limit, offset)
	return u.pgRepo.GetMultiWithTotal(ctx, limit, offset, filters)
}

// GetMulti – legacy wrapper.
func (u *userUseCase) GetMulti(ctx context.Context, limit, offset int) ([]*models.SdUser, error) {
	usersList, _, err := u.GetMultiWithTotal(ctx, limit, offset, nil)
	return usersList, err
}

// keys is a helper to log map keys.
func keys(m map[string]interface{}) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

// invalidateUserCaches clears user + profile caches.
func (u *userUseCase) invalidateUserCaches(ctx context.Context, id uuid.UUID) {
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id))
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisProfileKey(id))
}

// Profile – joined profile + role title + permission flags (cache-through TTL 1h).
func (u *userUseCase) Profile(ctx context.Context, id uuid.UUID) (*users.UserProfileDetail, error) {
	cacheKey := u.GenerateRedisProfileKey(id)
	if cached, _ := u.redisRepo.GetDetail(ctx, cacheKey); cached != nil {
		return cached, nil
	}
	detail, err := u.pgRepo.GetProfileWithPermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.redisRepo.SetDetail(ctx, cacheKey, detail, 3600); err != nil {
		u.logger.Warnf("Failed to cache profile %s: %v", id, err)
	}
	return detail, nil
}

// ListUsers – admin listing with filters.
func (u *userUseCase) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	u.logger.Infof("ListUsers keyword=%q sort=%s limit=%d offset=%d", f.Keyword, f.SortField, f.Limit, f.Offset)
	return u.pgRepo.ListUsers(ctx, f)
}

// Statistics – aggregated user stats.
func (u *userUseCase) Statistics(ctx context.Context) ([]users.UserStatusCount, error) {
	return u.pgRepo.GetStatistics(ctx)
}

// ActiveByNotification – target list by channel.
func (u *userUseCase) ActiveByNotification(ctx context.Context, ch users.NotificationChannel) ([]*models.SdUser, error) {
	switch ch {
	case users.NotifyEmail, users.NotifySms, users.NotifyLine:
	default:
		return nil, httpErrors.ErrValidation(fmt.Errorf("invalid notification channel: %s", ch))
	}
	return u.pgRepo.GetActiveByNotification(ctx, ch)
}

// UpdateActiveStatus – toggle active_status (0/1).
func (u *userUseCase) UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error {
	if status != 0 && status != 1 {
		return httpErrors.ErrValidation(errors.New("active_status must be 0 or 1"))
	}
	if _, err := u.Get(ctx, id); err != nil {
		return err
	}
	if err := u.pgRepo.UpdateActiveStatus(ctx, id, status); err != nil {
		return err
	}
	u.invalidateUserCaches(ctx, id)
	return nil
}

// UpdateAvatar – validate size/type → save disk → update DB → invalidate cache.
func (u *userUseCase) UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error {
	maxBytes := u.Cfg.Upload.MaxAvatarBytes
	if maxBytes <= 0 {
		maxBytes = 2 << 20
	}
	if fh.Size > maxBytes {
		return httpErrors.ErrValidation(fmt.Errorf("avatar exceeds max size %d bytes", maxBytes))
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png":
	default:
		return httpErrors.ErrValidation(errors.New("avatar must be jpg/jpeg/png"))
	}

	dir := u.Cfg.Upload.Dir
	if dir == "" {
		dir = "./uploads"
	}
	avatarDir := filepath.Join(dir, "avatar")
	if err := os.MkdirAll(avatarDir, 0o755); err != nil {
		return err
	}
	filename := id.String() + ext
	dstPath := filepath.Join(avatarDir, filename)

	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	publicPath := "/uploads/avatar/" + filename
	if err := u.pgRepo.UpdateAvatar(ctx, id, publicPath, filename); err != nil {
		_ = os.Remove(dstPath)
		return err
	}
	u.invalidateUserCaches(ctx, id)
	return nil
}

func (u *userUseCase) GetByEmail(ctx context.Context, email string) (*models.SdUser, error) {
	return u.pgRepo.GetByEmail(ctx, email)
}