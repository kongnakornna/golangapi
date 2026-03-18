package middleware

import (
	"context"
	"net/http"
	"strings"

	httpx "github.com/kongnakornna/golangapi/internal/transport/httpx"
	apperrors "github.com/kongnakornna/golangapi/pkg/errors"
	jwtpkg "github.com/kongnakornna/golangapi/pkg/jwt"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// UserIDKey คีย์สำหรับเก็บ ID ผู้ใช้ในคอนเท็กซ์
type UserIDKey struct{}

// RoleKey คีย์สำหรับเก็บบทบาทผู้ใช้ในคอนเท็กซ์
type RoleKey struct{}

// JWTConfig โครงสร้างกำหนดค่าสำหรับมิดเดิลแวร์ JWT
type JWTConfig struct {
	Secret       string   // คีย์ลับสำหรับ JWT
	ExcludePaths []string // เส้นทางที่ไม่ต้องผ่านการตรวจสอบสิทธิ์
}

// JWTAuth มิดเดิลแวร์สำหรับตรวจสอบสิทธิ์ด้วย JWT
func JWTAuth(config *JWTConfig, log logger.Logger) func(http.Handler) http.Handler {
	if log == nil {
		log = logger.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ข้ามการตรวจสอบสำหรับคำขอ OPTIONS
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// ตรวจสอบว่าเส้นทางปัจจุบันอยู่ในรายการที่ต้องยกเว้นหรือไม่
			path := r.URL.Path
			for _, excludePath := range config.ExcludePaths {
				if strings.HasPrefix(path, excludePath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// ดึงโทเคนจากส่วนหัวคำขอ
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				renderUnauthorized(w, r, "ไม่มีโทเคนยืนยันตัวตน")
				return
			}

			// แยกส่วนประกอบของโทเคน
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				renderUnauthorized(w, r, "รูปแบบโทเคนยืนยันตัวตนไม่ถูกต้อง")
				return
			}
			tokenString := tokenParts[1]

			// แยกวิเคราะห์โทเคน
			claims, err := jwtpkg.ParseToken(tokenString, config.Secret)
			if err != nil {
				ctxLogger := log.WithContext(r.Context())
				ctxLogger.Error("การแยกวิเคราะห์โทเคนล้มเหลว", "error", err, "token", tokenString)
				renderUnauthorized(w, r, "โทเคนยืนยันตัวตนไม่ถูกต้อง")
				return
			}

			// เพิ่ม ID ผู้ใช้และบทบาทลงในคอนเท็กซ์
			ctx := context.WithValue(r.Context(), UserIDKey{}, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey{}, claims.Role)

			// หากมี RequestContext ให้เพิ่มข้อมูลผู้ใช้ลงไปด้วย
			reqCtx := GetRequestContext(ctx)
			if reqCtx != nil {
				reqCtx.UserID = claims.UserID
				reqCtx.UserRole = claims.Role
			}

			// ดำเนินการคำขอต่อไป
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID ดึง ID ผู้ใช้จากคอนเท็กซ์
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(uint)
	return userID, ok
}

// GetRole ดึงบทบาทผู้ใช้จากคอนเท็กซ์
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey{}).(string)
	return role, ok
}

// RequireRole มิดเดิลแวร์สำหรับกำหนดให้ต้องมีบทบาทเฉพาะ
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := GetRole(r.Context())
			if !ok || userRole != role {
				renderForbidden(w, r, "ไม่มีสิทธิ์เข้าถึง")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// renderUnauthorized ส่งการตอบสนองข้อผิดพลาดประเภทไม่ได้รับอนุญาต
func renderUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	err := apperrors.New(apperrors.ErrorTypeUnauthorized, message, nil)
	httpx.Error(w, r, err)
}

// renderForbidden ส่งการตอบสนองข้อผิดพลาดประเภทสิทธิ์ไม่เพียงพอ
func renderForbidden(w http.ResponseWriter, r *http.Request, message string) {
	err := apperrors.New(apperrors.ErrorTypeForbidden, message, nil)
	httpx.Error(w, r, err)
}
