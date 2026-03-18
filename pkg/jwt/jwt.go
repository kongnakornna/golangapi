package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config การกำหนดค่า JWT
type Config struct {
	Secret          string        // คีย์ลับสำหรับ JWT
	AccessTokenExp  time.Duration // อายุการใช้งานของโทเค็นเข้าถึง
	RefreshTokenExp time.Duration // อายุการใช้งานของโทเค็นรีเฟรช
	Issuer          string        // ผู้ออกโทเค็น
}

// Claims ข้อมูลการอ้างสิทธิ์แบบกำหนดเองสำหรับ JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken สร้างโทเค็นสำหรับการเข้าถึง
func GenerateAccessToken(userID uint, role string, config *Config) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenExp)), // เวลาหมดอายุ
			IssuedAt:  jwt.NewNumericDate(time.Now()),                           // เวลาที่ออก
			NotBefore: jwt.NewNumericDate(time.Now()),                           // เวลาที่เริ่มใช้งานได้
			Issuer:    config.Issuer,                                            // ผู้ออก
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// GenerateRefreshToken สร้างโทเค็นสำหรับรีเฟรช
func GenerateRefreshToken(userID uint, config *Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenExp)), // เวลาหมดอายุ
		IssuedAt:  jwt.NewNumericDate(time.Now()),                           // เวลาที่ออก
		NotBefore: jwt.NewNumericDate(time.Now()),                           // เวลาที่เริ่มใช้งานได้
		Issuer:    config.Issuer,                                            // ผู้ออก
		Subject:   fmt.Sprintf("%d", userID),                                // ระบุ userID ใน subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// ParseToken แปลงค่าและตรวจสอบโทเค็นเข้าถึง
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("วิธีการลงนามที่ไม่คาดคิด: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("โทเค็นไม่ถูกต้อง")
}

// ParseRefreshToken แปลงค่าและตรวจสอบโทเค็นรีเฟรช
func ParseRefreshToken(tokenString string, secret string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("วิธีการลงนามที่ไม่คาดคิด: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// ดึง userID จาก Subject
		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("รหัสผู้ใช้ไม่ถูกต้อง: %w", err)
		}
		return uint(userID), nil
	}

	return 0, fmt.Errorf("โทเค็นไม่ถูกต้อง")
}

// ValidateToken ตรวจสอบว่าโทเค็นถูกต้องหรือไม่
func ValidateToken(tokenString string, secret string) bool {
	_, err := ParseToken(tokenString, secret)
	return err == nil
}
