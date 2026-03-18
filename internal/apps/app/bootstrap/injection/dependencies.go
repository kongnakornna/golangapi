package injection

import (
	"go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/kongnakornna/golangapiinternal/platform/config"
	"github.com/kongnakornna/golangapipkg/cache"
	"github.com/kongnakornna/golangapipkg/logger"
	"github.com/kongnakornna/golangapipkg/queue"
	"github.com/kongnakornna/golangapipkg/transaction"
)

// Dependencies คอนเทนเนอร์สำหรับ dependencies ของแอปพลิเคชัน
// นี่คือคอนเทนเนอร์หลักสำหรับ dependencies ของแอปพลิเคชัน ซึ่งจัดระเบียบความสัมพันธ์ของ dependencies ในแต่ละชั้นตามรูปแบบ dependency injection
type Dependencies struct {
	// ชั้น data access layer - รับผิดชอบการโต้ตอบกับฐานข้อมูล
	Repositories *Repositories

	// ชั้น business logic layer - ห่อหุ้มกฎทางธุรกิจหลัก
	Services *Services

	// ชั้น presentation layer - จัดการคำขอ HTTP และการตอบสนอง
	Handlers *Handlers

	// คอนฟิกูเรชันของแอปพลิเคชัน - ข้อมูลคอนฟิกทั่วโลก
	Config *config.AppConfig

	// โครงสร้างพื้นฐาน - ให้การสนับสนุนในระดับพื้นฐาน
	Infrastructure struct {
		DB                 *gorm.DB
		Redis              *redis.Client
		Cache              cache.Cache
		Validator          *validator.Validate
		Logger             logger.Logger
		Queue              queue.Queue
		TransactionManager transaction.Manager
	}
}

// NewDependencies เริ่มต้นคอนเทนเนอร์สำหรับ dependencies
// เป็นไปตามหลักการ dependency inversion โดยเริ่มต้นคอมโพเนนต์ในแต่ละชั้นจากล่างขึ้นบน:
// โครงสร้างพื้นฐาน -> ชั้น repository -> ชั้น service -> ชั้น handler
func NewDependencies(
	db *gorm.DB, // การเชื่อมต่อฐานข้อมูล
	rdb *redis.Client, // ไคลเอนต์ Redis
	validate *validator.Validate, // ตัวตรวจสอบความถูกต้อง
	appConfig *config.AppConfig, // คอนฟิกูเรชันของแอปพลิเคชัน
	cacheInstance cache.Cache, // อินสแตนซ์แคช
	appLogger logger.Logger, // ตัวบันทึกข้อมูล
) *Dependencies {
	// สร้างตัวจัดการคิว (รองรับเฉพาะ Redis)
	var queueManager queue.Queue
	if rdb != nil {
		queueManager = queue.NewRedisQueue(rdb, 10)
	} else {
		queueManager = queue.NewNoop()
	}
	// หากไม่มี Redis ฟังก์ชันการทำงานของคิวจะไม่สามารถใช้งานได้

	// สร้างตัวจัดการธุรกรรม
	txManager := transaction.NewGormTransactionManager(db)

	// สร้างคอนเทนเนอร์สำหรับ dependencies
	deps := &Dependencies{
		Config: appConfig,
		Infrastructure: struct {
			DB                 *gorm.DB
			Redis              *redis.Client
			Cache              cache.Cache
			Validator          *validator.Validate
			Logger             logger.Logger
			Queue              queue.Queue
			TransactionManager transaction.Manager
		}{
			DB:                 db,
			Redis:              rdb,
			Cache:              cacheInstance,
			Validator:          validate,
			Logger:             appLogger,
			Queue:              queueManager,
			TransactionManager: txManager,
		},
	}

	// 1. เริ่มต้น dependencies ในชั้น repository - ชั้น data access layer
	deps.Repositories = InitRepositories(db, appLogger)

	// 2. เริ่มต้น dependencies ในชั้น service - ชั้น business logic layer
	deps.Services = InitServices(deps.Repositories, validate, db, appConfig, cacheInstance, appLogger)

	// 3. เริ่มต้น dependencies ในชั้น handler - ชั้น presentation layer
	deps.Handlers = InitHandlers(deps.Services, appLogger, validate, db, rdb)

	// ส่งคืนคอนเทนเนอร์ dependencies ที่ประกอบเสร็จเรียบร้อยแล้ว
	return deps
}