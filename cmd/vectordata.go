package cmd

import (
	"context"

	"icmongolang/config"
	"icmongolang/internal/modules/vectordata/usecase"
	"icmongolang/pkg/db/postgres"
	esclient "icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/vectordb"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// vectordataCmd – vector database commands (seed sample data).
var vectordataCmd = &cobra.Command{
	Use:   "vectordata",
	Short: "Vector database commands",
	Long:  "Vector database commands (seed sample data into the configured vector DB)",
}

var vectordataSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed sample vector documents",
	Long: `Seed sample vector documents into the configured vector database
(elasticsearch or pgvector, per vectorDb.provider in config).

Requires the LLM embedding service to be running.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		appLogger.Infof("--vectordata seed Run-- (provider=%s)", cfg.VectorDB.Provider)

		ctx := context.Background()

		var db *gorm.DB
		if cfg.VectorDB.Provider == "pgvector" {
			psqlDB, err := postgres.NewPsqlDB(cfg)
			if err != nil {
				appLogger.Fatalf("เชื่อมต่อไม่สำเร็จ - Postgresql init: %s", err)
			}
			db = psqlDB
		}

		var es *esclient.Client
		if cfg.VectorDB.Provider == "elasticsearch" {
			esClient, err := esclient.NewClient(&cfg.Elasticsearch, appLogger)
			if err != nil {
				appLogger.Fatalf("Elasticsearch init: %s", err)
			}
			es = esClient
		}

		vdb, err := vectordb.New(&cfg.VectorDB, es, db, appLogger)
		if err != nil {
			appLogger.Fatalf("Vector DB init: %s", err)
		}

		if err := vdb.EnsureIndex(ctx, cfg.VectorDB.Dims); err != nil {
			appLogger.Warnf("⚠️ EnsureIndex: %v", err)
		}

		llmClient := llm.NewClient(&cfg.LLM, appLogger)

		vdUC := usecase.NewVectorDataUseCase(vdb, llmClient, cfg, appLogger)
		resp, err := vdUC.Seed(ctx, nil)
		if err != nil {
			appLogger.Fatalf("Seed failed: %s", err)
		}
		appLogger.Infof("🌱 Seeded %d sample documents into %q (provider=%s)", resp.Seeded, resp.Index, resp.Provider)
	},
}

func init() {
	RootCmd.AddCommand(vectordataCmd)
	vectordataCmd.AddCommand(vectordataSeedCmd)
}
