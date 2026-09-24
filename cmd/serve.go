package cmd

import (
	"icmongolang/config"
	"icmongolang/internal/distributor"
	"icmongolang/internal/server"
	"icmongolang/pkg/db/postgres"
	"icmongolang/pkg/db/redis"
	"icmongolang/pkg/logger"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()

		appLogger.Infof("--Serve Run--")
		appLogger.Infof(" ✅ AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("❌ เชื่อมต่อ Databse Postgresql ล้มเหลว - Postgresql init: %s", err)
			appLogger.Error("❌ เชื่อมต่อ Databse Postgresql ล้มเหลว", "error", err)
		} else {
			appLogger.Infof("✅ Postgres connected")
			appLogger.Info("✅ เชื่อมต่อ Databse Postgresql สำเร็จ")
		}

		if cfg.Server.MigrateOnStart {
			if err := Migrate(psqlDB, appLogger); err != nil {
				appLogger.Warnf("Migration error: %v", err)
			} else {
				appLogger.Info("✅ Data migrated successfully")
			}
		}

		redisClient := redis.NewRedis(cfg)

		taskRedisClient := distributor.NewRedisClient(cfg)

		server, err := server.NewServer(cfg, psqlDB, redisClient, taskRedisClient, appLogger)
		if err != nil {
			appLogger.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
