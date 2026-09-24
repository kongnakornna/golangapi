package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"icmongolang/config"
	kafkaHttp "icmongolang/internal/modules/kafka/delivery/http"
	"icmongolang/internal/modules/kafka/delivery/ws"
	"icmongolang/internal/modules/kafka/repository"
	"icmongolang/internal/modules/kafka/usecase"
	"icmongolang/internal/modules/queue"
	"icmongolang/pkg/db/postgres"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	// Ensure a root command exists for adding subcommands when not defined elsewhere
	if RootCmd == nil {
		RootCmd = &cobra.Command{Use: "icmongolang"}
	}
	RootCmd.AddCommand(kafkaCmd)
}

// RootCmd is the base command for the CLI. If the project defines its own
// RootCmd in another file, this declaration will be ignored due to the
// conditional initialization above.
var RootCmd *cobra.Command

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start Kafka service (producer + consumer + WebSocket)",
	Long:  `Starts HTTP server to accept orders, WebSocket for real-time updates, and Kafka consumer to process orders asynchronously.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()
		log := logger.NewLogger("kafka")
		log.InitLogger()
		log.Infof("✅ Starting Kafka + WebSocket service...")

		db, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			log.Fatalf("❌ Cannot connect to DB: %v", err)
		}
		log.Info("✅ Connected to PostgreSQL")

		producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Timeout)
		if err != nil {
			log.Fatalf("❌ Failed to create Kafka producer: %v", err)
		}
		defer producer.Close()
		log.Infof("✅ Kafka producer connected to %v", cfg.Kafka.Brokers)

		msgQueue := queue.NewQueue()
		hub := websocket.NewHub(msgQueue, log)
		go hub.Run()
		log.Info("✅ WebSocket hub started")

		orderRepo := repository.NewOrderRepository(db)
		orderUsecase := usecase.NewOrderUsecase(orderRepo, producer, cfg.Kafka.Topic, hub)

		consumerHandler := func(msg *kafka.OrderMessage) error {
			return orderUsecase.ProcessOrderMessage(msg)
		}
		consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic, consumerHandler, cfg.Kafka.Timeout)
		if err != nil {
			log.Fatalf("❌ Failed to create Kafka consumer: %v", err)
		}
		defer consumer.Close()

		ctx, cancel := context.WithCancel(context.Background())
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("❌ Failed to start consumer: %v", err)
		}
		log.Info("✅ Kafka consumer started")

		orderHandler := kafkaHttp.NewOrderHandler(orderUsecase)
		wsHandler := ws.NewWsHandler(hub, orderUsecase, log)

		router := mux.NewRouter()
		router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
		router.HandleFunc("/ws", wsHandler.ServeWS).Methods("GET")
		router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}).Methods("GET")
		router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))

		httpPort := cfg.Server.Port
		if httpPort == "" {
			httpPort = "5051"
		}
		srv := &http.Server{
			Addr:    ":" + httpPort,
			Handler: router,
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			log.Infof("🚀 HTTP server listening on port %s", httpPort)
			log.Infof("📚 Swagger UI at http://localhost:%s/swagger/index.html", httpPort)
			log.Infof("🔌 WebSocket endpoint at ws://localhost:%s/ws", httpPort)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("❌ HTTP server error: %v", err)
			}
		}()

		<-stop
		log.Info("⏳ Shutting down...")

		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("❌ HTTP shutdown error: %v", err)
		}

		cancel()
		consumer.Wait()
		msgQueue.Close()
		if err := log.Sync(); err != nil {
			log.Errorf("Log sync error: %v", err)
		}
		log.Info("✅ Service stopped gracefully")
	},
}
