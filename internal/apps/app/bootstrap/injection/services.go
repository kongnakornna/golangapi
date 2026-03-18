package injection

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	authservice "github.com/kongnakornna/golangapi/internal/core/auth/service"
	userservice "github.com/kongnakornna/golangapi/internal/core/user/service"
	"github.com/kongnakornna/golangapi/internal/platform/config"
	"github.com/kongnakornna/golangapi/pkg/cache"
	"github.com/kongnakornna/golangapi/pkg/jwt"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// Services ชุดของบริการทั้งหมด
// ประกอบด้วยออบเจกต์ทั้งหมดของชั้นตรรกะทางธุรกิจ และจัดการกฎหลักทางธุรกิจ
type Services struct {
	// ธุรกิจที่เกี่ยวข้องกับผู้ใช้
	UserService userservice.UserService

	// ธุรกิจที่เกี่ยวข้องกับการยืนยันตัวตน
	AuthService authservice.AuthService

	// สามารถเพิ่มบริการอื่นๆ ได้ที่นี่...
	// ProductService userservice.ProductService
	// OrderService userservice.OrderService
}

// InitServices เริ่มต้นบริการทั้งหมด
// นี่คือชั้นที่สองของการฉีด dependencies ซึ่งขึ้นอยู่กับชั้น repository
func InitServices(
	repos *Repositories,
	validate *validator.Validate,
	db *gorm.DB,
	config *config.AppConfig,
	cacheInstance cache.Cache,
	log logger.Logger,
) *Services {
	// ตรวจสอบพารามิเตอร์
	if repos == nil {
		log.Error("ไม่สามารถให้ repository เป็นค่าว่างได้")
		os.Exit(1)
	}
	if validate == nil {
		log.Error("ไม่สามารถให้ตัวตรวจสอบเป็นค่าว่างได้")
		os.Exit(1)
	}
	if db == nil {
		log.Error("ไม่สามารถให้การเชื่อมต่อฐานข้อมูลเป็นค่าว่างได้")
		os.Exit(1)
	}
	if config == nil {
		log.Error("ไม่สามารถให้การกำหนดค่าเป็นค่าว่างได้")
		os.Exit(1)
	}

	// สร้างการกำหนดค่า JWT
	jwtConfig := createJWTConfig(config, log)

	// สร้างอินสแตนซ์บริการทั้งหมด
	userService := userservice.NewUserService(repos.UserRepo, validate, db, cacheInstance)
	authService := authservice.NewAuthService(repos.UserRepo, validate, db, jwtConfig, cacheInstance)

	// คืนค่าชุดบริการ
	return &Services{
		UserService: userService,
		AuthService: authService,
	}
}

// createJWTConfig สร้างการกำหนดค่า JWT จากการกำหนดค่าแอปพลิเคชัน
// นี่คือฟังก์ชันช่วยเหลือ สำหรับสร้างการกำหนดค่าที่จำเป็นสำหรับบริการ JWT
func createJWTConfig(config *config.AppConfig, log logger.Logger) *jwt.Config {
	if config.JWT.Secret == "" {
		log.Warn("คีย์ลับ JWT ว่างเปล่า ซึ่งอาจทำให้เกิดปัญหาด้านความปลอดภัย")
	}

	return &jwt.Config{
		Secret:          config.JWT.Secret,
		AccessTokenExp:  config.JWT.AccessTokenExp,
		RefreshTokenExp: config.JWT.RefreshTokenExp,
		Issuer:          config.JWT.Issuer,
	}
}