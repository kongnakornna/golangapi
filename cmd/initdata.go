package cmd

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/distributor"
	userDistributor "icmongolang/internal/modules/users/distributor"
	userRepository "icmongolang/internal/modules/users/repository"
	userUseCase "icmongolang/internal/modules/users/usecase"
	"icmongolang/pkg/db/postgres"
	"icmongolang/pkg/db/redis"
	"icmongolang/pkg/logger"

	"github.com/spf13/cobra"
)

var initDataCmd = &cobra.Command{
	Use:   "initdata",
	Short: "Init data",
	Long:  "Init data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		appLogger.Infof("--initdata Run--")
		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("เชื่อมต่อไม่สำเร็จ - Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected เชื่อมต่อสำเร็จ..")
		}

		redisClient := redis.NewRedis(cfg)

		taskRedisClient := distributor.NewRedisClient(cfg)

		// Repository
		userPgRepo := userRepository.CreateUserPgRepository(psqlDB)
		userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)

		// Distributor
		userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, appLogger)

		// UseCase
		userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, appLogger)

		// Create super user if not exists
		isCreated, _ := userUC.CreateSuperUserIfNotExist(context.Background())

		// isCreated, _ := userUC.CreateSuperUserIfNotExist(context.Background(), "admin", "password", "admin")
		if !isCreated {
			appLogger.Info("Super user is exists, skip create")
		} else {
			appLogger.Info("Created super user")
		}
	},
}

func init() {
	RootCmd.AddCommand(initDataCmd)
}
