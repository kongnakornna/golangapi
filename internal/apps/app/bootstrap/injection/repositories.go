package injection

import (
	"os"

	"gorm.io/gorm"

	userrepo "github.com/kongnakornna/golangapi/internal/core/user/repository"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// Repositories คือชุดรวมของ repositories ทั้งหมด
// ประกอบด้วยออบเจ็กต์สำหรับเข้าถึงข้อมูลทุกระดับ ซึ่งทำหน้าที่โต้ตอบกับแหล่งข้อมูล
type Repositories struct {
	// ออบเจ็กต์สำหรับเข้าถึงข้อมูลผู้ใช้
	UserRepo userrepo.UserRepository

	// สามารถเพิ่ม repositories อื่นๆ ได้ที่นี่...
	// ProductRepo repository.ProductRepository
	// OrderRepo repository.OrderRepository
}

// InitRepositories ใช้สำหรับเริ่มต้น repositories ทั้งหมด
// เป็นเลเยอร์แรกของการ dependency injection ที่รับผิดชอบสร้างออบเจ็กต์สำหรับเข้าถึงข้อมูลทั้งหมด
func InitRepositories(db *gorm.DB, log logger.Logger) *Repositories {
	// ตรวจสอบพารามิเตอร์
	if db == nil {
		log.Error("การเชื่อมต่อฐานข้อมูลไม่สามารถเป็นค่าว่างได้")
		os.Exit(1)
	}

	// สร้างอินสแตนซ์ของ repositories ทั้งหมด
	userRepo := userrepo.NewUserRepository(db)

	// ส่งกลับชุดรวมของ repositories
	return &Repositories{
		UserRepo: userRepo,
	}
}
