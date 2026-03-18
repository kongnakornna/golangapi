package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kongnakornna/golangapi/internal/platform/config"
)

// InitDB เริ่มต้นการเชื่อมต่อฐานข้อมูล
func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// ปรับแต่งสำหรับโปรดักชัน: ปรับระดับ log
	logLevel := logger.Warn
	if cfg.Driver == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		PrepareStmt:                              true, // เตรียมคำสั่งไว้ล่วงหน้า เพิ่มประสิทธิภาพ
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("เชื่อมต่อฐานข้อมูลล้มเหลว: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ดึงการเชื่อมต่อฐานข้อมูลล้มเหลว: %w", err)
	}

	// ปรับแต่ง connection pool สำหรับโปรดักชัน
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxLifetime / 2) // หมดเวลาของการเชื่อมต่อที่ไม่ได้ใช้งาน

	// ทดสอบการเชื่อมต่อ
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping ฐานข้อมูลล้มเหลว: %w", err)
	}

	return db, nil
}

// InitRedis เริ่มต้นการเชื่อมต่อ Redis
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10, // ขนาด connection pool
		MinIdleConns: 5,  // จำนวนการเชื่อมต่อที่ว่างขั้นต่ำ
		MaxRetries:   2,  // จำนวนครั้งที่ลองใหม่สูงสุด
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// ทดสอบการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("เชื่อมต่อ Redis ล้มเหลว: %w", err)
	}

	return rdb, nil
}
