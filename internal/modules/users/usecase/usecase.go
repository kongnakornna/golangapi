package usecase

import (
	"context"
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
	"icmongolang/internal/usecase"
	"icmongolang/internal/modules/users"
	"icmongolang/internal/worker"
	"icmongolang/pkg/cryptpass"
	"icmongolang/pkg/emailTemplates"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/jwt"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/secureRandom"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// userUseCase implements UserUseCaseI interface
// userUseCase implements อินเทอร์เฟซ UserUseCaseI
type userUseCase struct {
	usecase.UseCase[models.SdUser]
	pgRepo                 users.UserPgRepository
	redisRepo              users.UserRedisRepository
	emailTemplateGenerator emailTemplates.EmailTemplatesGenerator
	redisTaskDistributor   users.UserRedisTaskDistributor
	logger                 logger.Logger // Logger instance for structured logging
}

// CreateUserUseCaseI creates new user use case instance
// CreateUserUseCaseI สร้างอินสแตนซ์ use case สำหรับผู้ใช้
func CreateUserUseCaseI(
	pgRepo users.UserPgRepository,
	redisRepo users.UserRedisRepository,
	redisTaskDistributor users.UserRedisTaskDistributor,
	cfg *config.Config,
	logger logger.Logger,
) users.UserUseCaseI {
	logger.Info("Initializing user use case")
	return &userUseCase{
		UseCase:                usecase.CreateUseCase[models.SdUser](pgRepo, cfg, logger),
		pgRepo:                 pgRepo,
		redisRepo:              redisRepo,
		emailTemplateGenerator: emailTemplates.NewEmailTemplatesGenerator(cfg),
		redisTaskDistributor:   redisTaskDistributor,
		logger:                 logger,
	}
}

func (u *userUseCase) SignInByUsername(ctx context.Context, username, password string) (string, string, error) {
	u.logger.Info("Sign in attempt")
	user, err := u.pgRepo.GetByEmail(ctx, username)
	if err != nil {
		u.logger.Warn("User not found")
		return "", "", httpErrors.ErrWrongPassword(errors.New("invalid credentials"))
	}

	if !cryptpass.ComparePassword(password, user.Password) {
		u.logger.Warn("Invalid password")
		return "", "", httpErrors.ErrWrongPassword(errors.New("wrong password"))
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

// storeRefreshToken stores the refresh token in the Redis set and applies TTL
// storeRefreshToken เก็บ refresh token ลง Redis set พร้อมกำหนด TTL
func (u *userUseCase) storeRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	key := u.GenerateRedisRefreshTokenKey(id)
	if err := u.redisRepo.Sadd(ctx, key, refreshToken); err != nil {
		return err
	}
	ttlSeconds := int(u.Cfg.Jwt.RefreshTokenExpireDuration) * 60 // config หน่วยนาที
	if err := u.redisRepo.Expire(ctx, key, ttlSeconds); err != nil {
		u.logger.Warnf("Failed to set TTL on refresh token set %s: %v", key, err)
	}
	return nil
}

// Get retrieves user by ID (with Redis cache)
// Get ดึงข้อมูลผู้ใช้ตาม ID (พร้อม Redis cache)
func (u *userUseCase) Get(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	u.logger.Infof("Getting user by ID: %s", id)
	cacheKey := u.GenerateRedisUserKey(id)

	// Try to get from cache
	cachedUser, err := u.redisRepo.Get(ctx, cacheKey)
	if err != nil {
		u.logger.Errorf("Redis get error for key %s: %v", cacheKey, err)
		return nil, err
	}
	if cachedUser != nil {
		u.logger.Debugf("Cache hit for user %s", id)
		return cachedUser, nil
	}
	u.logger.Debugf("Cache miss for user %s, fetching from DB", id)

	// Get from database
	user, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		u.logger.Errorf("Database get error for user %s: %v", id, err)
		return nil, err
	}

	// Store in cache for 1 hour
	if err = u.redisRepo.Create(ctx, cacheKey, user, 3600); err != nil {
		u.logger.Warnf("Failed to cache user %s: %v", id, err)
		// Non-fatal, continue
	}
	u.logger.Infof("User %s retrieved successfully", id)
	return user, nil
}

// Delete removes user from DB and clears Redis cache
// Delete ลบผู้ใช้ออกจากฐานข้อมูลและล้าง Redis cache
func (u *userUseCase) Delete(ctx context.Context, id uuid.UUID) (*models.SdUser, error) {
	u.logger.Infof("Deleting user: %s", id)
	user, err := u.pgRepo.Delete(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to delete user %s from DB: %v", id, err)
		return nil, err
	}

	// Clear cache
	u.invalidateUserCaches(ctx, id)
	if err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		u.logger.Warnf("Failed to delete refresh token cache for %s: %v", id, err)
	}
	u.logger.Infof("User %s deleted successfully", id)
	return user, nil
}

// Update updates user fields (values map) and invalidates cache
// Update อัปเดตฟิลด์ผู้ใช้ (ผ่าน map) และล้าง cache
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

	// Invalidate cache
	u.invalidateUserCaches(ctx, id)
	u.logger.Infof("User %s updated successfully", id)
	return user, nil
}

// Create – inserts a new user and sends verification email
// Create – สร้างผู้ใช้ใหม่และส่งอีเมลยืนยันตัวตน
func (u *userUseCase) Create(ctx context.Context, exp *models.SdUser) (*models.SdUser, error) {
	u.logger.Info("Creating new user")

	exp.Email = strings.ToLower(strings.TrimSpace(exp.Email))
	exp.Password = strings.TrimSpace(exp.Password)

	// Pre-check for duplicate email / username before hitting the DB unique constraints
	// ตรวจสอบ email / username ซ้ำก่อน insert เพื่อให้ error ชัดเจน
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

	// Hash password
	hashedPassword, err := cryptpass.HashPassword(exp.Password)
	if err != nil {
		u.logger.Errorf("Password hashing failed for %s: %v", exp.Email, err)
		return nil, err
	}
	exp.Password = hashedPassword
	exp.PasswordTemp = nil // ห้ามเก็บ plaintext ค้างใน model
	u.logger.Debug("Password hashed successfully")

	// Set username default to email if empty
	if exp.Username == "" {
		exp.Username = exp.Email
		u.logger.Debugf("Username set to email: %s", exp.Username)
	}

	// Default role if not provided
	if exp.RoleID == 0 {
		exp.RoleID = 2 // adjust default role ID as needed
		u.logger.Debugf("Default role_id=2 assigned to %s", exp.Email)
	}

	// Save to database
	user, err := u.pgRepo.Create(ctx, exp)
	if err != nil {
		u.logger.Errorf("Database create failed for %s: %v", exp.Email, err)
		return nil, err
	}
	u.logger.Infof("User created with ID: %s, verified=%v", user.ID, user.Verified)

	// If already verified, skip email sending
	if user.Verified {
		u.logger.Infof("User %s already verified, skipping verification email", user.Email)
		return user, nil
	}

	// Generate verification code
	verificationCode, err := secureRandom.RandomHex(16)
	if err != nil {
		u.logger.Errorf("Failed to generate verification code for %s: %v", user.Email, err)
		return nil, err
	}
	u.logger.Debugf("Verification code generated for %s", user.Email)

	// Update user with verification code
	updatedUser, err := u.pgRepo.UpdateVerificationCode(ctx, user, verificationCode)
	if err != nil {
		u.logger.Errorf("Failed to save verification code for %s: %v", user.Email, err)
		return nil, err
	}
	u.logger.Debugf("Verification code saved for user %s", updatedUser.ID)

	// Prepare name for email template
	name := ""
	if updatedUser.Fullname != nil {
		name = *updatedUser.Fullname
	} else {
		name = updatedUser.Email
	}

	// Generate email content
	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GenerateVerificationCodeTemplate(
		ctx,
		name,
		fmt.Sprintf("%s/auth/verifyemail?code=%s", strings.TrimRight(u.Cfg.Server.BaseUrl, "/"), verificationCode),
	)
	if err != nil {
		u.logger.Errorf("Failed to generate email template for %s: %v", updatedUser.Email, err)
		return nil, err
	}
	u.logger.Debug("Email template generated")

	// Send email via async task
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

// CreateUser – wrapper with password confirmation
// CreateUser – ฟังก์ชันครอบที่มีการยืนยันรหัสผ่าน
func (u *userUseCase) CreateUser(ctx context.Context, exp *models.SdUser, confirmPassword string) (*models.SdUser, error) {
	u.logger.Debugf("CreateUser called for email: %s", exp.Email)
	if exp.Password != confirmPassword {
		u.logger.Warnf("Password confirmation mismatch for %s", exp.Email)
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}
	return u.Create(ctx, exp)
}

// createToken generates access and refresh tokens
// createToken สร้าง access token และ refresh token
func (u *userUseCase) createToken(ctx context.Context, exp models.SdUser) (string, string, error) {
	u.logger.Debugf("Generating tokens for user: %s", exp.ID)
	accessToken, err := jwt.CreateAccessTokenRS256(
		exp.ID.String(),
		exp.Email,
		u.Cfg.Jwt.AccessTokenPrivateKey,
		u.Cfg.Jwt.AccessTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.Issuer,
	)
	if err != nil {
		u.logger.Errorf("Access token generation failed for %s: %v", exp.ID, err)
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(
		exp.ID.String(),
		exp.Email,
		u.Cfg.Jwt.RefreshTokenPrivateKey,
		u.Cfg.Jwt.RefreshTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.Issuer,
	)
	if err != nil {
		u.logger.Errorf("Refresh token generation failed for %s: %v", exp.ID, err)
		return "", "", err
	}
	u.logger.Debugf("Tokens generated for user %s", exp.ID)
	return accessToken, refreshToken, nil
}

// SignIn authenticates user and returns tokens
// SignIn ตรวจสอบสิทธิ์ผู้ใช้และคืนค่า token
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

// IsActive checks if user account is active
// IsActive ตรวจสอบว่าบัญชีผู้ใช้ active หรือไม่
func (u *userUseCase) IsActive(ctx context.Context, exp models.SdUser) bool {
	active := exp.Status == 1
	u.logger.Debugf("User %s active status: %v", exp.ID, active)
	return active
}

// IsSuper checks if user has superuser privileges
// IsSuper ตรวจสอบว่าผู้ใช้มีสิทธิ์ superuser หรือไม่
func (u *userUseCase) IsSuper(ctx context.Context, exp models.SdUser) bool {
	super := exp.IsSuperUser
	u.logger.Debugf("User %s superuser status: %v", exp.ID, super)
	return super
}

// CreateSuperUserIfNotExist creates the first superuser from config
// CreateSuperUserIfNotExist สร้าง superuser รายแรกจาก config ถ้ายังไม่มี
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

// UpdatePassword changes user password and invalidates all sessions
// UpdatePassword เปลี่ยนรหัสผ่านผู้ใช้และยกเลิกทุก session
func (u *userUseCase) UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword, confirmPassword string) (*models.SdUser, error) {
	u.logger.Infof("Password change request for user: %s", id)
	if newPassword != confirmPassword {
		u.logger.Warnf("Password confirmation mismatch for user %s", id)
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.Get(ctx, id)
	if err != nil {
		u.logger.Errorf("User %s not found for password update: %v", id, err)
		return nil, err
	}

	if !cryptpass.ComparePassword(oldPassword, user.Password) {
		u.logger.Warnf("Old password mismatch for user %s", id)
		return nil, httpErrors.ErrWrongPassword(errors.New("old password and new password not same"))
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		u.logger.Errorf("Failed to hash new password for %s: %v", id, err)
		return nil, err
	}

	updatedUser, err := u.pgRepo.UpdatePassword(ctx, user, hashedPassword)
	if err != nil {
		u.logger.Errorf("Database password update failed for %s: %v", id, err)
		return nil, err
	}

	// Invalidate all caches and sessions
	u.invalidateUserCaches(ctx, id)
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id))
	u.logger.Infof("Password changed successfully for user %s, all sessions invalidated", id)
	return updatedUser, nil
}

// ParseIdFromRefreshToken extracts user ID from refresh token
// ParseIdFromRefreshToken แยก user ID จาก refresh token
func (u *userUseCase) ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	u.logger.Debug("Parsing refresh token")
	id, _, err := jwt.ParseTokenRS256(refreshToken, u.Cfg.Jwt.RefreshTokenPublicKey)
	if err != nil {
		u.logger.Errorf("Failed to parse refresh token: %v", err)
		return uuid.UUID{}, err
	}
	idParsed, err := uuid.Parse(id)
	if err != nil {
		u.logger.Errorf("Invalid user ID in token: %s, error: %v", id, err)
		return uuid.UUID{}, httpErrors.ErrInvalidJWTClaims(errors.New("can not convert id to uuid from id in token"))
	}
	u.logger.Debugf("Refresh token belongs to user: %s", idParsed)
	return idParsed, nil
}

// Refresh issues new tokens using a valid refresh token
// Refresh ออก token ใหม่โดยใช้ refresh token ที่ยังไม่หมดอายุ
func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	u.logger.Info("Token refresh request")
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	isMember, err := u.redisRepo.SIsMember(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken)
	if err != nil {
		u.logger.Errorf("Redis SIsMember error for user %s: %v", idParsed, err)
		return "", "", err
	}
	if !isMember {
		u.logger.Warnf("Refresh token not found in Redis for user %s", idParsed)
		return "", "", httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("not found refresh token in redis"))
	}

	// Remove old refresh token (rotation)
	if err = u.redisRepo.Srem(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken); err != nil {
		u.logger.Errorf("Failed to remove old refresh token for %s: %v", idParsed, err)
		return "", "", err
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

	if err = u.storeRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		u.logger.Errorf("Failed to store new refresh token for %s: %v", user.ID, err)
		return "", "", err
	}
	u.logger.Infof("Tokens refreshed successfully for user %s", user.ID)
	return accessToken, newRefreshToken, nil
}

// Logout removes the specific refresh token from Redis
// Logout ลบ refresh token ที่ระบุออกจาก Redis
func (u *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	u.logger.Info("Logout request")
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	err = u.redisRepo.Srem(ctx, u.GenerateRedisRefreshTokenKey(idParsed), refreshToken)
	if err != nil {
		u.logger.Errorf("Failed to logout user %s: %v", idParsed, err)
	} else {
		u.logger.Infof("User %s logged out successfully", idParsed)
	}
	return err
}

// LogoutAll removes all refresh tokens of a user
// LogoutAll ลบ refresh token ทั้งหมดของผู้ใช้
func (u *userUseCase) LogoutAll(ctx context.Context, id uuid.UUID) error {
	u.logger.Infof("Logout all sessions for user: %s", id)
	err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id))
	if err != nil {
		u.logger.Errorf("Failed to logout all for user %s: %v", id, err)
	} else {
		u.logger.Infof("All sessions terminated for user %s", id)
	}
	return err
}

// Verify confirms user email using verification code
// Verify ยืนยันอีเมลผู้ใช้ด้วย verification code
func (u *userUseCase) Verify(ctx context.Context, verificationCode string) error {
	u.logger.Infof("Email verification attempt with code: %s", verificationCode)
	user, err := u.pgRepo.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		u.logger.Warnf("Invalid verification code: %s, error: %v", verificationCode, err)
		return err
	}
	if user.Verified {
		u.logger.Warnf("User %s already verified", user.Email)
		return httpErrors.ErrUserAlreadyVerified(errors.New("user already verified"))
	}

	updatedUser, err := u.pgRepo.UpdateVerification(ctx, user, "", true)
	if err != nil {
		u.logger.Errorf("Failed to update verification status for %s: %v", user.Email, err)
		return err
	}
	u.invalidateUserCaches(ctx, updatedUser.ID)
	u.logger.Infof("Email verified successfully for user: %s", updatedUser.Email)
	return nil
}

// ForgotPassword sends a password reset email
// ForgotPassword ส่งอีเมลรีเซ็ตรหัสผ่าน
func (u *userUseCase) ForgotPassword(ctx context.Context, email string) error {
	u.logger.Infof("Forgot password request for email: %s", email)
	user, err := u.pgRepo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.Warnf("Email not found: %s", email)
		return httpErrors.ErrNotFound(err)
	}
	if !user.Verified {
		u.logger.Warnf("Unverified user attempted password reset, continuing: %s", email)
	}

	resetToken, err := secureRandom.RandomHex(16)
	if err != nil {
		u.logger.Errorf("Failed to generate reset token for %s: %v", email, err)
		return err
	}
	u.logger.Debugf("Reset token generated for %s", email)

	updatedUser, err := u.pgRepo.UpdatePasswordReset(ctx, user, resetToken, time.Now().In(helpers.GetTimeLocation()).Add(15*time.Minute))
	if err != nil {
		u.logger.Errorf("Failed to save reset token for %s: %v", email, err)
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
		u.logger.Errorf("Failed to generate reset email template for %s: %v", email, err)
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
		u.logger.Errorf("Failed to enqueue reset email for %s: %v", email, err)
		return err
	}
	u.logger.Infof("Password reset email sent to %s", email)
	return nil
}

// ResetPassword performs password reset using token
// ResetPassword ดำเนินการรีเซ็ตรหัสผ่านด้วย token
func (u *userUseCase) ResetPassword(ctx context.Context, resetToken, newPassword, confirmPassword string) error {
	u.logger.Info("Password reset attempt")
	if newPassword != confirmPassword {
		u.logger.Warn("Password confirmation mismatch during reset")
		return httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.pgRepo.GetByResetTokenResetAt(ctx, resetToken, time.Now().In(helpers.GetTimeLocation()))
	if err != nil {
		u.logger.Warnf("Invalid or expired reset token: %s", resetToken)
		return err
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		u.logger.Errorf("Failed to hash new password for %s: %v", user.Email, err)
		return err
	}

	updatedUser, err := u.pgRepo.UpdatePasswordResetToken(ctx, user, hashedPassword, "")
	if err != nil {
		u.logger.Errorf("Failed to update password for %s: %v", user.Email, err)
		return err
	}
	u.invalidateUserCaches(ctx, updatedUser.ID)
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(updatedUser.ID))
	u.logger.Infof("Password reset successfully for user: %s", updatedUser.Email)
	return nil
}

// GenerateRedisUserKey returns Redis key for user cache
// GenerateRedisUserKey คืนค่า Redis key สำหรับ cache ผู้ใช้
func (u *userUseCase) GenerateRedisUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", models.SdUser{}.TableName(), id.String())
}

// GenerateRedisRefreshTokenKey returns Redis key for user's refresh token set
// GenerateRedisRefreshTokenKey คืนค่า Redis key สำหรับชุด refresh token ของผู้ใช้
func (u *userUseCase) GenerateRedisRefreshTokenKey(id uuid.UUID) string {
	return fmt.Sprintf("RefreshToken:%s", id.String())
}

// GetMultiWithTotal คืนค่ารายชื่อผู้ใช้และจำนวนทั้งหมด (รองรับ filter และ sort)
func (u *userUseCase) GetMultiWithTotal(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.SdUser, int64, error) {
	u.logger.Infof("GetMultiWithTotal called: limit=%d, offset=%d, filters=%v", limit, offset, filters)
	return u.pgRepo.GetMultiWithTotal(ctx, limit, offset, filters)
}

// ถ้าต้องการคง method เดิม GetMulti ไว้ก็เรียกใช้ตัวใหม่
func (u *userUseCase) GetMulti(ctx context.Context, limit, offset int) ([]*models.SdUser, error) {
	u.logger.Infof("GetMulti called with limit=%d, offset=%d", limit, offset)
	users, _, err := u.GetMultiWithTotal(ctx, limit, offset, nil)
	return users, err
}

// Helper function to get keys from map for logging
func keys(m map[string]interface{}) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func (u *userUseCase) invalidateUserCaches(ctx context.Context, id uuid.UUID) {
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id))
	_ = u.redisRepo.Delete(ctx, u.GenerateRedisProfileKey(id))
}

func (u *userUseCase) GenerateRedisProfileKey(id uuid.UUID) string {
	return fmt.Sprintf("sd_user_profile:%s", id.String())
}

// Profile – joined profile + role title + permission flags (cache-through TTL 1h)
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

func (u *userUseCase) ListUsers(ctx context.Context, f users.UserFilter) ([]*models.SdUser, int64, error) {
	u.logger.Infof("ListUsers keyword=%q sort=%s limit=%d offset=%d", f.Keyword, f.SortField, f.Limit, f.Offset)
	return u.pgRepo.ListUsers(ctx, f)
}

func (u *userUseCase) Statistics(ctx context.Context) ([]users.UserStatusCount, error) {
	return u.pgRepo.GetStatistics(ctx)
}

func (u *userUseCase) ActiveByNotification(ctx context.Context, ch users.NotificationChannel) ([]*models.SdUser, error) {
	switch ch {
	case users.NotifyEmail, users.NotifySms, users.NotifyLine:
	default:
		return nil, httpErrors.ErrValidation(fmt.Errorf("invalid notification channel: %s", ch))
	}
	return u.pgRepo.GetActiveByNotification(ctx, ch)
}

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

// UpdateAvatar – validate size/type → save disk → update DB → invalidate cache
func (u *userUseCase) UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error {
	maxBytes := u.Cfg.Upload.MaxAvatarBytes
	if maxBytes <= 0 {
		maxBytes = 2 << 20 // 2MB default
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
		_ = os.Remove(dstPath) // rollback file ถ้า DB fail
		return err
	}
	u.invalidateUserCaches(ctx, id)
	return nil
}
