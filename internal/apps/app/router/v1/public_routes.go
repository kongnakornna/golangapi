package v1

import (
	"github.com/go-chi/chi/v5"
	authmodule "github.com/kongnakornna/golangapi/internal/core/auth"
	authhandler "github.com/kongnakornna/golangapi/internal/core/auth/handler"
	userhandler "github.com/kongnakornna/golangapi/internal/core/user/handler"
)

// RouterConfig การกำหนดค่าเส้นทาง
type RouterConfig struct {
	UserHandler *userhandler.UserHandler
	AuthHandler *authhandler.AuthHandler
	JWTSecret   string
}

// SetupPublicRoutes ตั้งค่าเส้นทางสาธารณะ (ไม่ต้องรับรองความถูกต้อง)
func SetupPublicRoutes(r chi.Router, config RouterConfig) {
	authmodule.RegisterPublicRoutes(r, config.AuthHandler)
}