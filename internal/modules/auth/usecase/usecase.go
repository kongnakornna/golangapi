package usecase

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/modules/auth"
	"icmongolang/internal/modules/users"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

// authUseCase implements AuthUseCaseI by delegating to users.UserUseCaseI.
type authUseCase struct {
	usersUC users.UserUseCaseI
	cfg     *config.Config
	logger  logger.Logger
}

// CreateAuthUseCaseI creates a new AuthUseCaseI instance.
func CreateAuthUseCaseI(usersUC users.UserUseCaseI, cfg *config.Config, logger logger.Logger) auth.AuthUseCaseI {
	logger.Info("Initializing auth use case")
	return &authUseCase{
		usersUC: usersUC,
		cfg:     cfg,
		logger:  logger,
	}
}

func (u *authUseCase) SignIn(ctx context.Context, email, password string) (string, string, error) {
	return u.usersUC.SignIn(ctx, email, password)
}

func (u *authUseCase) Login(ctx context.Context, username, password string) (string, string, error) {
	return u.usersUC.SignInByUsername(ctx, username, password)
}

func (u *authUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	return u.usersUC.Refresh(ctx, refreshToken)
}

func (u *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	return u.usersUC.Logout(ctx, refreshToken)
}

// LogoutAll accepts a userID string to match AuthUseCaseI interface.
// The current handler calls usersUC.LogoutAll(uuid) directly; this method
// exists to keep the interface satisfied and for future use.
func (u *authUseCase) LogoutAll(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		u.logger.Errorf("Invalid user ID for LogoutAll: %s, error: %v", userID, err)
		return err
	}
	return u.usersUC.LogoutAll(ctx, id)
}

func (u *authUseCase) Verify(ctx context.Context, verificationCode string) error {
	return u.usersUC.Verify(ctx, verificationCode)
}

func (u *authUseCase) ForgotPassword(ctx context.Context, email string) error {
	return u.usersUC.ForgotPassword(ctx, email)
}

func (u *authUseCase) ResetPassword(ctx context.Context, resetToken, newPassword, confirmPassword string) error {
	return u.usersUC.ResetPassword(ctx, resetToken, newPassword, confirmPassword)
}
