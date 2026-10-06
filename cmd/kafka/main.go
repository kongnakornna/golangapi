// Package main — standalone Kafka service binary
//
// @title         Kafka Order Service
// @version       1.0
// @description   Asynchronous order processing with Kafka + WebSocket
// @host          localhost:5051
// @BasePath      /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

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
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.GetCfg()
	log := logger.NewLogger("kafka")
	log.InitLogger()
	log.Infof("🚀 Starting Kafka Service...")

	// ---------- DB ----------
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	log.Info("✅ PostgreSQL connected")

	// ---------- Kafka producer ----------
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalf("❌ Producer failed: %v", err)
	}
	defer func() { _ = producer.Close() }()
	log.Infof("✅ Kafka producer ready (brokers: %v)", cfg.Kafka.Brokers)

	// ---------- WebSocket hub ----------
	msgQueue := queue.NewQueue()
	hub := websocket.NewHub(websocket.HubConfig{
		Queue: msgQueue,
		Log:   log,
	})
	go hub.Run()
	log.Info("✅ WebSocket hub started")

	// ---------- Usecase ----------
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, producer, cfg.Kafka.Topic, hub)

	// ---------- Kafka consumer ----------
	consumerHandler := func(msg *kafka.OrderMessage) error {
		return orderUsecase.ProcessOrderMessage(msg)
	}
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		consumerHandler,
		cfg.Kafka.Timeout,
	)
	if err != nil {
		log.Fatalf("❌ Consumer failed: %v", err)
	}
	defer func() { _ = consumer.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("❌ Consumer start error: %v", err)
	}
	log.Info("✅ Kafka consumer started")

	// ---------- HTTP ----------
	orderHandler := kafkaHttp.NewOrderHandler(orderUsecase)
	wsHandler := ws.NewWsHandler(hub, orderUsecase, log)

	router := mux.NewRouter()
	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	router.HandleFunc("/ws", wsHandler.ServeWS).Methods("GET")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	port := cfg.Server.Port
	if port == "" {
		port = "5051"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infof("📡 HTTP server listening on port %s", port)
		log.Infof("📚 Swagger: http://localhost:%s/swagger/index.html", port)
		log.Infof("🔌 WebSocket: ws://localhost:%s/ws", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Info("⏳ Shutting down...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Errorf("HTTP shutdown error: %v", err)
	}

	cancel()
	consumer.Wait()
	msgQueue.Close()
	if err := log.Sync(); err != nil {
		log.Errorf("Log sync error: %v", err)
	}
	log.Info("✅ Service stopped gracefully")
}
