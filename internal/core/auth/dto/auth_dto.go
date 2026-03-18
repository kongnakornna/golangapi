package dto

import userdto "github.com/vadxq/go-rest-starter/internal/core/user/dto"

// UserResponse ใช้โครงสร้างการตอบสนองผู้ใช้ซ้ำ
type UserResponse = userdto.UserResponse

// LoginRequest คำขอเข้าสู่ระบบ
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse การตอบสนองการเข้าสู่ระบบ
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest คำขอรีเฟรชโทเค็น
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse การตอบสนองโทเค็น
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}