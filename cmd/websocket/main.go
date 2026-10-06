package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"icmongolang/config"
	"icmongolang/internal/middleware"
	"icmongolang/internal/modules/queue"
	httprest "icmongolang/internal/modules/websocket/delivery/http"
	"icmongolang/internal/modules/websocket/delivery/ws"
	"icmongolang/internal/modules/websocket/repository/postgres"
	"icmongolang/internal/modules/websocket/usecase"
	"icmongolang/internal/repository"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/websocket"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	log := logger.NewLogger("websocket")
	log.InitLogger()
	log.Info("🚀 Starting WebSocket (Redis Queue) Service...")

	cfg := config.GetCfg()

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
	hub := websocket.NewHub(websocket.HubConfig{
		Queue: q,
		Log:   log,
		Callbacks: websocket.HubCallbacks{
			OnConnect: func(client *websocket.Client) error {
				log.Infof("WS client connected: %s (user: %s, role: %s)", client.ID, client.UserID, client.Role)
				return nil
			},
			OnDisconnect: func(client *websocket.Client) {
				log.Infof("WS client disconnected: %s (user: %s)", client.ID, client.UserID)
			},
			OnSubscribe: func(client *websocket.Client, topic string) {
				log.Debugf("WS client %s subscribed to topic: %s", client.ID, topic)
			},
			OnUnsubscribe: func(client *websocket.Client, topic string) {
				log.Debugf("WS client %s unsubscribed from topic: %s", client.ID, topic)
			},
			OnJoinRoom: func(client *websocket.Client, room string) {
				log.Debugf("WS client %s joined room: %s", client.ID, room)
			},
			OnLeaveRoom: func(client *websocket.Client, room string) {
				log.Debugf("WS client %s left room: %s", client.ID, room)
			},
			OnMessage: func(client *websocket.Client, topic, room string, payload []byte) bool {
				log.Debugf("WS message from %s to topic=%s room=%s: %s", client.ID, topic, room, string(payload))
				return true // allow broadcast
			},
		},
	})
	go hub.Run()
	log.Info("✅ WebSocket hub started")

	uc := usecase.NewWSUsecase(repo, q, hub)

	// Create stop channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start background retention cleanup scheduler
	cleanupInterval := 1 * time.Hour // Run cleanup every hour
	if cfg.Retention.Enabled && cfg.Retention.CleanupInterval > 0 {
		cleanupInterval = cfg.Retention.CleanupInterval
	}
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		log.Infof("🧹 Retention cleanup scheduler started (interval: %v)", cleanupInterval)
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				stats, err := uc.RunRetentionCleanup(ctx)
				cancel()
				if err != nil {
					log.Errorf("Retention cleanup failed: %v", err)
				} else {
					log.Infof("✅ Retention cleanup completed: deleted=%d messages, freed=%d bytes, policies=%d",
						stats.MessagesDeleted, stats.BytesFreed, stats.ActivePolicies)
				}
			case <-stop:
				log.Info("🛑 Retention cleanup scheduler stopped")
				return
			}
		}
	}()

	// Create middleware manager for WebSocket authentication
	redisRepo := repository.CreateRedisRepository[any](redisClient)
	mw := middleware.CreateMiddlewareManager(cfg, log, nil, redisRepo)

	wsHandler := ws.NewWsHandler(hub, uc, log, mw)
	restHandler := httprest.NewWsRestHandler(hub, uc, log)

	router := mux.NewRouter()
	// WebSocket endpoint with Origin guard
	router.Handle("/ws", mw.WSOriginGuard(cfg.WS.AllowedOrigins)(wsHandler.HandlerFunc())).Methods("GET")
	router.Handle("/ws/health", mw.WSOriginGuard(cfg.WS.AllowedOrigins)(wsHandler.HandlerFunc())).Methods("GET")
	router.HandleFunc("/api/ws/messages", restHandler.SendMessage).Methods("POST")
	router.HandleFunc("/api/ws/messages", restHandler.GetMessages).Methods("GET")
	router.HandleFunc("/api/ws/messages/paginated", restHandler.GetMessagesPaginated).Methods("GET")
	router.HandleFunc("/api/ws/rooms/{room}/messages", restHandler.GetMessagesByRoom).Methods("GET")
	router.HandleFunc("/api/ws/senders/{sender_id}/messages", restHandler.GetMessagesBySender).Methods("GET")
	router.HandleFunc("/api/ws/messages/search", restHandler.SearchMessages).Methods("GET")
	router.HandleFunc("/api/ws/rooms", restHandler.GetRooms).Methods("GET")
	router.HandleFunc("/api/ws/rooms/{room}/stats", restHandler.GetRoomStats).Methods("GET")
	router.HandleFunc("/api/ws/presence/{user_id}", restHandler.GetPresence).Methods("GET")
	router.HandleFunc("/api/ws/presence", restHandler.GetAllPresence).Methods("GET")
	// Room Presence endpoints
	router.HandleFunc("/api/ws/rooms/{room}/presence", restHandler.GetRoomPresence).Methods("GET")
	router.HandleFunc("/api/ws/users/{user_id}/room", restHandler.GetUserCurrentRoom).Methods("GET")
	// Retention Policy endpoints
	router.HandleFunc("/api/ws/retention/policies", restHandler.CreateRetentionPolicy).Methods("POST")
	router.HandleFunc("/api/ws/retention/policies", restHandler.ListRetentionPolicies).Methods("GET")
	router.HandleFunc("/api/ws/retention/policies/{id}", restHandler.GetRetentionPolicy).Methods("GET")
	router.HandleFunc("/api/ws/retention/policies/{id}", restHandler.UpdateRetentionPolicy).Methods("PUT")
	router.HandleFunc("/api/ws/retention/policies/{id}", restHandler.DeleteRetentionPolicy).Methods("DELETE")
	router.HandleFunc("/api/ws/retention/stats", restHandler.GetRetentionStats).Methods("GET")
	router.HandleFunc("/api/ws/retention/cleanup", restHandler.RunRetentionCleanup).Methods("POST")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("WEBSOCKET_PORT")
	if port == "" {
		port = "5000"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

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

	// hub.Shutdown() - method doesn't exist, just close the hub context
	// hub.Shutdown()
	q.Close()
	if err := log.Sync(); err != nil {
		log.Errorf("Log sync error: %v", err)
	}
	log.Info("✅ Service stopped")
}
