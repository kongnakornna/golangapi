package postgres

import (
	"context"
	"fmt"
	"time"

	"icmongolang/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPsqlDB(cfg *config.Config) (*gorm.DB, error) {
	timeout := cfg.Postgres.ConnectionTimeout
	if timeout <= 0 {
		timeout = 10
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok connect_timeout=%d",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname, cfg.Postgres.Port, timeout,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	// Bound the pool so bursts (e.g. the MQTT ingester spawning one DB
	// goroutine per device per message) cannot exhaust Postgres's
	// max_connections. Defaults keep concurrent API requests healthy while
	// leaving headroom for other services sharing the server.
	maxOpen := cfg.Postgres.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25
	}
	maxIdle := cfg.Postgres.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 10
	}
	maxIdle = min(maxIdle, maxOpen)
	connMaxLifetime := time.Duration(cfg.Postgres.ConnMaxLifetimeMin) * time.Minute
	if cfg.Postgres.ConnMaxLifetimeMin <= 0 {
		connMaxLifetime = 5 * time.Minute
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	return db, nil
}
