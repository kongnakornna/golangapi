package auth

import "context"

// AuthUseCaseI defines the authentication business logic interface.
// AuthUseCaseI กำหนดอินเทอร์เฟซสำหรับตรรกะทางธุรกิจเกี่ยวกับการยืนยันตัวตน
type AuthUseCaseI interface {
	SignIn(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	Verify(ctx context.Context, verificationCode string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken, newPassword, confirmPassword string) error
}
