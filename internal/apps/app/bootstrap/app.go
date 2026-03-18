package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	""github.com/redis/go-redis/v9"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/kongnakornna/golangapi/internal/apps/app/bootstrap/injection"
	api "github.com/kongnakornna/golangapi/internal/apps/app/router"
	"github.com/kongnakornna/golangapi/internal/platform/config"
	"github.com/kongnakornna/golangapi/internal/platform/db"
	httpx "github.com/kongnakornna/golangapi/internal/transport/httpx"
	"github.com/kongnakornna/golangapi/pkg/cache"
	"github.com/kongnakornna/golangapi/pkg/logger"
)

// App โครงสร้างแอปพลิเคชัน
type App struct {
	DB        *gorm.DB
	Redis     *redis.Client
	Router    *chi.Mux
	Cache     cache.Cache
	Validator *validator.Validate
	Deps      *injection.Dependencies
	Server    *http.Server
	Config    *config.AppConfig
	logger    logger.Logger
}

// New สร้างอินสแตนซ์แอปพลิเคชันใหม่
func New() (*App, error) {
	// กำหนดค่าเส้นทางไฟล์กำหนดค่า
	configPath := getConfigPath()

	// โหลดการกำหนดค่า
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("โหลดการกำหนดค่าล้มเหลว: %w", err)
	}

	// สร้างตัวบันทึก日志
	appLogger, err := logger.NewLogger(&logger.LogConfig{
		Level:   cfg.Log.Level,
		File:    cfg.Log.File,
		Console: cfg.Log.Console,
	})
	if err != nil {
		return nil, fmt.Errorf("สร้างตัวบันทึก日志ล้มเหลว: %w", err)
	}

	appLogger.Info("โหลดการกำหนดค่าเสร็จสมบูรณ์", "config_path", configPath)
	httpx.SetLogger(appLogger)

	// สร้างอินสแตนซ์แอปพลิเคชัน
	app := &App{
		Config: cfg,
		logger: appLogger,
	}

	// เริ่มต้นแอปพลิเคชัน
	if err := app.initialize(); err != nil {
		return nil, fmt.Errorf("เริ่มต้นแอปพลิเคชันล้มเหลว: %w", err)
	}

	return app, nil
}

// initialize เริ่มต้นส่วนประกอบของแอปพลิเคชัน
func (app *App) initialize() error {
	app.logger.Info("กำลังเริ่มต้นแอปพลิเคชัน...")

	// เริ่มต้นการเชื่อมต่อฐานข้อมูล
	if err := app.initDatabase(); err != nil {
		return fmt.Errorf("เริ่มต้นฐานข้อมูลล้มเหลว: %w", err)
	}

	// เริ่มต้นการเชื่อมต่อ Redis
	if err := app.initRedis(); err != nil {
		return fmt.Errorf("เริ่มต้น Redis ล้มเหลว: %w", err)
	}

	// เริ่มต้นแคช
	if err := app.initCache(); err != nil {
		return fmt.Errorf("เริ่มต้นแคชล้มเหลว: %w", err)
	}

	// เริ่มต้นตัวตรวจสอบ
	app.Validator = validator.New()

	// เริ่มต้นการฉีด dependencies
	if err := app.initDependencies(); err != nil {
		return fmt.Errorf("เริ่มต้นการฉีด dependencies ล้มเหลว: %w", err)
	}

	// เริ่มต้นเส้นทาง
	if err := app.initRouter(); err != nil {
		return fmt.Errorf("เริ่มต้นเส้นทางล้มเหลว: %w", err)
	}

	app.logger.Info("เริ่มต้นแอปพลิเคชันเสร็จสมบูรณ์")
	return nil
}

// initDatabase เริ่มต้นการเชื่อมต่อฐานข้อมูล
func (app *App) initDatabase() error {
	app.logger.Info("กำลังเชื่อมต่อฐานข้อมูล...")

	database, err := db.InitDB(&app.Config.Database)
	if err != nil {
		return err
	}

	app.DB = database
	app.logger.Info("เชื่อมต่อฐานข้อมูลสำเร็จ")
	return nil
}

// initRedis เริ่มต้นการเชื่อมต่อ Redis
func (app *App) initRedis() error {
	app.logger.Info("กำลังเชื่อมต่อ Redis...")

	if !app.Config.Redis.Enabled {
		app.logger.Info("Redis ถูกปิดใช้งาน ข้ามการเริ่มต้น")
		return nil
	}

	redisClient, err := db.InitRedis(&app.Config.Redis)
	if err != nil {
		app.logger.Warn("เชื่อมต่อ Redis ล้มเหลว ลดระดับเป็น Noop", "error", err)
		return nil
	}

	if redisClient == nil {
		app.logger.Info("Redis ไม่ได้เปิดใช้งาน")
		return nil
	}

	app.Redis = redisClient
	app.logger.Info("เชื่อมต่อ Redis สำเร็จ")
	return nil
}

// initCache เริ่มต้นแคช
func (app *App) initCache() error {
	app.logger.Info("กำลังเริ่มต้นแคช...")

	// บริการแคชต้องพึ่งพา Redis
	if app.Redis == nil {
		app.logger.Warn("Redis ไม่พร้อมใช้งาน ลดระดับแคชเป็น Noop")
		app.Cache = cache.NewNoop()
		return nil
	}

	cacheOpts := cache.Options{
		DefaultExpiration: 10 * time.Minute,
		CleanupInterval:   5 * time.Minute,
		RedisAddress:      fmt.Sprintf("%s:%d", app.Config.Redis.Host, app.Config.Redis.Port),
		RedisPassword:     app.Config.Redis.Password,
		RedisDB:           app.Config.Redis.DB,
	}

	app.logger.Info("ใช้ Redis เป็นพื้นที่จัดเก็บแคช")

	cacheInstance, err := cache.NewCache(cacheOpts)
	if err != nil {
		app.logger.Warn("เริ่มต้น Redis แคชล้มเหลว ลดระดับแคชเป็น Noop", "error", err)
		app.Cache = cache.NewNoop()
		return nil
	}

	app.Cache = cacheInstance
	app.logger.Info("เริ่มต้นแคชสำเร็จ")
	return nil
}

// initDependencies เริ่มต้นการฉีด dependencies
func (app *App) initDependencies() error {
	app.logger.Info("กำลังเริ่มต้นระบบการฉีด dependencies...")

	deps := injection.NewDependencies(
		app.DB,
		app.Redis,
		app.Validator,
		app.Config,
		app.Cache,
		app.logger,
	)

	app.Deps = deps
	app.logger.Info("เริ่มต้นระบบการฉีด dependencies เสร็จสมบูรณ์")
	return nil
}

// initRouter เริ่มต้นเส้นทาง
func (app *App) initRouter() error {
	app.logger.Info("กำลังกำหนดค่าเส้นทาง API...")

	router := chi.NewRouter()

	api.Setup(router, api.RouterConfig{
		UserHandler:   app.Deps.Handlers.UserHandler,
		AuthHandler:   app.Deps.Handlers.AuthHandler,
		HealthHandler: app.Deps.Handlers.HealthHandler,
		JWTSecret:     app.Deps.Config.JWT.Secret,
		Logger:        app.logger,
	})

	app.Router = router
	app.logger.Info("กำหนดค่าเส้นทาง API เสร็จสมบูรณ์")
	return nil
}

// StartServer เริ่มต้นเซิร์ฟเวอร์ HTTP
func (app *App) StartServer() <-chan error {
	errCh := make(chan error, 1)

	// สร้างเซิร์ฟเวอร์ HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler:      app.Router,
		ReadTimeout:  app.Config.Server.ReadTimeout,
		WriteTimeout: app.Config.Server.WriteTimeout,
	}

	app.Server = server

	// เริ่มต้นเซิร์ฟเวอร์
	go func() {
		app.logger.Info("เริ่มต้นเซิร์ฟเวอร์ HTTP", "port", app.Config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("เซิร์ฟเวอร์ HTTP ผิดพลาด: %w", err)
		}
	}()

	return errCh
}

// Shutdown ปิดแอปพลิเคชันอย่างราบรื่น
func (app *App) Shutdown(ctx context.Context) error {
	app.logger.Info("กำลังปิดแอปพลิเคชันอย่างราบรื่น...")

	// ใช้ channel เพื่อรวบรวมข้อผิดพลาด
	errChan := make(chan error, 3)

	// ปิดส่วนประกอบต่างๆ พร้อมกัน
	go func() {
		if app.Server != nil {
			app.logger.Info("กำลังปิดเซิร์ฟเวอร์ HTTP...")
			errChan <- app.Server.Shutdown(ctx)
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if app.DB != nil {
			app.logger.Info("กำลังปิดการเชื่อมต่อฐานข้อมูล...")
			if sqlDB, err := app.DB.DB(); err == nil {
				errChan <- sqlDB.Close()
			} else {
				errChan <- err
			}
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if app.Redis != nil {
			app.logger.Info("กำลังปิดการเชื่อมต่อ Redis...")
			errChan <- app.Redis.Close()
		} else {
			errChan <- nil
		}
	}()

	// รอให้การปิดทั้งหมดเสร็จสิ้น
	var hasError bool
	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			app.logger.Error("ปิดส่วนประกอบล้มเหลว", "error", err)
			hasError = true
		}
	}

	if hasError {
		app.logger.Warn("เกิดข้อผิดพลาดขณะปิดแอปพลิเคชัน")
	} else {
		app.logger.Info("ปิดแอปพลิเคชันอย่างราบรื่นเสร็จสมบูรณ์")
	}

	return nil
}

// Logger คืนค่า logger ของแอปพลิเคชัน
func (app *App) Logger() logger.Logger {
	return app.logger
}

// รับเส้นทางไฟล์กำหนดค่า
func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return configPath
}
