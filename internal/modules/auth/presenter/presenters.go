package presenter

// SignInRequest represents the login request using email.
// SignInRequest ใช้สำหรับขอเข้าสู่ระบบด้วยอีเมล
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

// LoginRequest represents the login request using username.
// LoginRequest ใช้สำหรับขอเข้าสู่ระบบด้วยชื่อผู้ใช้
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

// RefreshTokenRequest carries the refresh token for token renewal.
// RefreshTokenRequest ใช้สำหรับขอ refresh token ใหม่
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJSUzI1NiIs..."`
}

// ForgotPasswordRequest contains the email for password reset.
// ForgotPasswordRequest มีอีเมลสำหรับขอรีเซ็ตรหัสผ่าน
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest carries the new password and confirmation.
// ResetPasswordRequest มีรหัสผ่านใหม่และยืนยันรหัสผ่าน
type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"newpass123"`
}

// TokenResponse is the standard token payload returned on successful auth.
// TokenResponse คือ payload token มาตรฐานที่ส่งกลับเมื่อยืนยันตัวตนสำเร็จ
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
