package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"icmongolang/internal/modules/queue"
	httprest "icmongolang/internal/modules/websocket/delivery/http"
	"icmongolang/internal/modules/websocket/delivery/ws"
	"icmongolang/internal/modules/websocket/repository/postgres"
	"icmongolang/internal/modules/websocket/usecase"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	log := logger.NewLogger("websocket")
	log.InitLogger()
	log.Info("🚀 Starting WebSocket (Redis Queue) Service...")

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()
	log.Info("✅ PostgreSQL connected")

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Info("✅ Redis connected")

	q := queue.NewRedisQueue(redisClient, 10)

	repo := postgres.NewPGRepository(db)
	hub := websocket.NewHub(q, log)
	go hub.Run()
	log.Info("✅ WebSocket hub started")

	uc := usecase.NewWSUsecase(repo, q, hub)

	wsHandler := ws.NewWsHandler(hub, uc, log)
	restHandler := httprest.NewWsRestHandler(hub, uc, log)

	router := mux.NewRouter()
	router.HandleFunc("/ws/health", wsHandler.ServeHTTP).Methods("GET")
	router.HandleFunc("/api/ws/messages", restHandler.SendMessage).Methods("POST")
	router.HandleFunc("/api/ws/messages", restHandler.GetMessages).Methods("GET")
	router.HandleFunc("/api/ws/rooms", restHandler.GetRooms).Methods("GET")
	router.HandleFunc("/api/ws/rooms/{room}/stats", restHandler.GetRoomStats).Methods("GET")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("WEBSOCKET_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infof("📡 WebSocket server listening on port %s", port)
		log.Infof("🔌 WebSocket endpoint: ws://localhost:%s/ws", port)
		log.Infof("📚 REST endpoints at /api/ws/...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Info("⏳ Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("HTTP shutdown error: %v", err)
	}

	hub.Shutdown()
	q.Close()
	if err := log.Sync(); err != nil {
		log.Errorf("Log sync error: %v", err)
	}
	log.Info("✅ Service stopped")
}
