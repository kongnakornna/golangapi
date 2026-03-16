package v1

import (
	"github.com/go-chi/chi/v5"
	authmodule "github.com/vadxq/go-rest-starter/internal/core/auth"
	usermodule "github.com/vadxq/go-rest-starter/internal/core/user"
	custommiddleware "github.com/vadxq/go-rest-starter/internal/transport/middleware"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// SetupProtectedRoutes ตั้งค่าเส้นทางที่ได้รับการป้องกัน (ต้องมีการรับรองตัวตน)
func SetupProtectedRoutes(r chi.Router, config RouterConfig, jwtConfig *custommiddleware.JWTConfig, log logger.Logger) {
	// สร้างกลุ่มเส้นทางที่ต้องการการยืนยันตัวตนด้วย JWT
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.JWTAuth(jwtConfig, log))

		// เส้นทางที่เกี่ยวข้องกับการรับรองตัวตน (ต้องการการยืนยันตัวตนแล้ว)
		authmodule.RegisterProtectedRoutes(r, config.AuthHandler)

		// เส้นทางทรัพยากรผู้ใช้
		usermodule.RegisterRoutes(r, config.UserHandler)
	})
}