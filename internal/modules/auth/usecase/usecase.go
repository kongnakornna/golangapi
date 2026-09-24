package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/modules/auth"
	"icmongolang/internal/modules/users"

	"github.com/google/uuid"
	"icmongolang/pkg/logger"
)

// authUseCase implements AuthUseCaseI by delegating to users.UserUseCaseI.
// authUseCase implement อินเทอร์เฟซ AuthUseCaseI โดยมอบหมายงานให้ users.UserUseCaseI
type authUseCase struct {
	usersUC users.UserUseCaseI
	cfg     *config.Config
	logger  logger.Logger
}

// CreateAuthUseCaseI creates a new AuthUseCaseI instance.
// CreateAuthUseCaseI สร้างอินสแตนซ์ AuthUseCaseI ใหม่
func CreateAuthUseCaseI(usersUC users.UserUseCaseI, cfg *config.Config, logger logger.Logger) auth.AuthUseCaseI {
	logger.Info("Initializing auth use case")
	return &authUseCase{
		usersUC: usersUC,
		cfg:     cfg,
		logger:  logger,
	}
}

func (u *authUseCase) SignIn(ctx context.Context, email, password string) (string, string, error) {
	u.logger.Info("Auth use case: SignIn")
	return u.usersUC.SignIn(ctx, email, password)
}

func (u *authUseCase) Login(ctx context.Context, username, password string) (string, string, error) {
	u.logger.Info("Auth use case: Login")
	return u.usersUC.SignInByUsername(ctx, username, password)
}

func (u *authUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	u.logger.Info("Auth use case: Refresh")
	return u.usersUC.Refresh(ctx, refreshToken)
}

func (u *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	u.logger.Info("Auth use case: Logout")
	return u.usersUC.Logout(ctx, refreshToken)
}

func (u *authUseCase) LogoutAll(ctx context.Context, userID string) error {
	u.logger.Info("Auth use case: LogoutAll")
	id, err := uuid.Parse(userID)
	if err != nil {
		u.logger.Errorf("Invalid user ID for LogoutAll: %s, error: %v", userID, err)
		return err
	}
	return u.usersUC.LogoutAll(ctx, id)
}

func (u *authUseCase) Verify(ctx context.Context, verificationCode string) error {
	u.logger.Info("Auth use case: Verify")
	return u.usersUC.Verify(ctx, verificationCode)
}

func (u *authUseCase) ForgotPassword(ctx context.Context, email string) error {
	u.logger.Info("Auth use case: ForgotPassword")
	return u.usersUC.ForgotPassword(ctx, email)
}

func (u *authUseCase) ResetPassword(ctx context.Context, resetToken, newPassword, confirmPassword string) error {
	u.logger.Info("Auth use case: ResetPassword")
	return u.usersUC.ResetPassword(ctx, resetToken, newPassword, confirmPassword)
}
