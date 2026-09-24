package users

import (
	"context"
	"mime/multipart"

	"icmongolang/internal"
	"icmongolang/internal/models"

	"github.com/google/uuid"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.SdUser]
	CreateUser(ctx context.Context, exp *models.SdUser, confirmPassword string) (*models.SdUser, error)
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	SignInByUsername(ctx context.Context, username, password string) (string, string, error)
	IsActive(ctx context.Context, exp models.SdUser) bool
	IsSuper(ctx context.Context, exp models.SdUser) bool
	CreateSuperUserIfNotExist(context.Context) (bool, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string, confirmPassword string) (*models.SdUser, error)
	ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	GenerateRedisUserKey(id uuid.UUID) string
	GenerateRedisRefreshTokenKey(id uuid.UUID) string
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, id uuid.UUID) error
	Verify(ctx context.Context, verificationCode string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken string, newPassword string, confirmPassword string) error
	GetMultiWithTotal(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.SdUser, int64, error)
	GenerateRedisProfileKey(id uuid.UUID) string
	Profile(ctx context.Context, id uuid.UUID) (*UserProfileDetail, error)
	ListUsers(ctx context.Context, f UserFilter) ([]*models.SdUser, int64, error)
	Statistics(ctx context.Context) ([]UserStatusCount, error)
	ActiveByNotification(ctx context.Context, ch NotificationChannel) ([]*models.SdUser, error)
	UpdateActiveStatus(ctx context.Context, id uuid.UUID, status int16) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, fh *multipart.FileHeader) error
}
