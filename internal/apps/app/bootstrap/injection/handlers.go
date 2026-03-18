package injection

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	authhandler "github.com/kongnakornna/golangapi/internal/core/auth/handler"
	healthhandler "github.com/kongnakornna/golangapi/internal/core/health/handler"
	userhandler "github.com/kongnakornna/golangapi/internal/core/user/handler"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// Handlers รวมตัวจัดการ HTTP ทั้งหมด
type Handlers struct {
	UserHandler   *userhandler.UserHandler
	AuthHandler   *authhandler.AuthHandler
	HealthHandler *healthhandler.HealthHandler
}

// InitHandlers เริ่มต้นตัวจัดการ HTTP ทั้งหมด
func InitHandlers(
	services *Services,
	log logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	redis *redis.Client,
) *Handlers {
	// เริ่มต้นตัวจัดการผู้ใช้
	userHandler := userhandler.NewUserHandler(
		services.UserService,
		log,
		validator,
	)

	// เริ่มต้นตัวจัดการการยืนยันตัวตน
	authHandler := authhandler.NewAuthHandler(
		services.AuthService,
		log,
		validator,
	)

	// เริ่มต้นตัวจัดการตรวจสอบสถานะ
	healthHandler := healthhandler.NewHealthHandler(
		db,
		redis,
		log,
	)

	return &Handlers{
		UserHandler:   userHandler,
		AuthHandler:   authHandler,
		HealthHandler: healthHandler,
	}
}