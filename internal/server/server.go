package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"icmongolang/config"
	"icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 10 // seconds — ให้เวลา PDPA/Kafka ปิดให้ทัน
)

// ============================================================================
// Server
// ============================================================================

type Server struct {
	server          *http.Server
	runtime         *Runtime
	cfg             *config.Config
	logger          logger.Logger
	db              *gorm.DB
	redisClient     *redis.Client
	taskRedisClient *asynq.Client
	mqttClient      mqtt.Client
	influxClient    *influxdb.InfluxClient
	esClient        *elasticsearch.Client
	llmClient       *llm.Client
	cancelBg        context.CancelFunc
}

// normalizeMQTTBroker แปลง URL scheme ให้ Paho MQTT Go client รับได้
//
// Paho Go รองรับเฉพาะ: tcp://, ssl://, ws://, wss://
// แปลง mqtt:// → tcp://, mqtts:// → ssl://
func normalizeMQTTBroker(broker string, log logger.Logger) string {
	switch {
	case strings.HasPrefix(broker, "mqtt://"):
		fixed := "tcp://" + strings.TrimPrefix(broker, "mqtt://")
		log.Warnf("⚠️  Auto-converted MQTT broker %q → %q (Paho Go requires tcp://)", broker, fixed)
		return fixed
	case strings.HasPrefix(broker, "mqtts://"):
		fixed := "ssl://" + strings.TrimPrefix(broker, "mqtts://")
		log.Warnf("⚠️  Auto-converted MQTT broker %q → %q (Paho Go requires ssl://)", broker, fixed)
		return fixed
	}
	return broker
}

func connectMQTTWithRetry(client mqtt.Client, timeout time.Duration, maxRetries int, log logger.Logger) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if client.IsConnected() {
			log.Warn("⚠️  MQTT already connected")
			return nil
		}

		log.Infof("✅ Connecting to MQTT (attempt %d/%d)", attempt, maxRetries)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err := client.Connect(ctx)
		cancel()

		if err == nil {
			log.Info("✅ MQTT connected successfully")
			return nil
		}

		lastErr = err
		log.Errorf("❌ MQTT connection attempt %d failed: %v", attempt, err)

		if attempt < maxRetries {
			backoff := time.Duration(attempt) * time.Second
			log.Infof("✅ Retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("❌ failed to connect MQTT after %d attempts: %w", maxRetries, lastErr)
}

// ============================================================================
// NewServer
// ============================================================================

func NewServer(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	taskRedisClient *asynq.Client,
	log logger.Logger,
) (*Server, error) {
	log.Info("✅ configuring server...")

	// ---------- Auto-fix MQTT broker URL ----------
	cfg.MQTT.Broker = normalizeMQTTBroker(cfg.MQTT.Broker, log)

	mqttClient := mqtt.New(&cfg.MQTT, log)

	connTimeout := time.Duration(cfg.MQTT.ConnectionTimeout) * time.Second
	if connTimeout == 0 {
		connTimeout = 30 * time.Second
		log.Warn("✅ MQTT connection timeout not set in config, using default 30s")
	}

	if err := connectMQTTWithRetry(mqttClient, connTimeout, 3, log); err != nil {
		log.Errorf("❌ Failed to connect MQTT: %v", err)
	} else {
		log.Info("✅ MQTT connected successfully")
	}

	// ---------- InfluxDB ----------
	influxClient, err := influxdb.NewInfluxClient(cfg, log)
	if err != nil {
		log.Errorf("❌ Failed to connect InfluxDB: %v", err)
		influxClient = nil
	} else {
		log.Info("✅ InfluxDB connected successfully")
	}

	// ---------- Elasticsearch ----------
	var esClient *elasticsearch.Client
	if len(cfg.Elasticsearch.Addresses) > 0 {
		esClient, err = elasticsearch.NewClient(&cfg.Elasticsearch, log)
		if err != nil {
			log.Errorf("❌ Failed to connect Elasticsearch: %v", err)
			esClient = nil
		} else {
			log.Info("✅ Elasticsearch connected successfully")
		}
	} else {
		log.Warn("⚠️ No Elasticsearch addresses configured – ES routes skipped")
	}

	// ---------- LLM ----------
	llmClient := llm.NewClient(&cfg.LLM, log)

	// ---------- Background context (bound to server lifecycle) ----------
	bgCtx, cancelBg := context.WithCancel(context.Background())

	// ---------- Router + Runtime ----------
	runtime, err := New(
		bgCtx, // ← parent context for background tasks (PDPA dispatcher, etc.)
		db,
		redisClient,
		taskRedisClient,
		cfg,
		log,
		influxClient,
		mqttClient,
		esClient,
		llmClient,
	)
	if err != nil {
		cancelBg()
		return nil, fmt.Errorf("⚠️ Failed to create router: %w", err)
	}

	addr := cfg.Server.Port
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	return &Server{
		server: &http.Server{
			Addr:           addr,
			Handler:        runtime, // *Runtime implements http.Handler via ServeHTTP
			BaseContext:    func(net.Listener) context.Context { return bgCtx },
			ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
			MaxHeaderBytes: maxHeaderBytes,
		},
		runtime:         runtime, // ใช้ตัวที่ได้จาก New() (มี hooks ครบ)
		cfg:             cfg,
		logger:          log,
		db:              db,
		redisClient:     redisClient,
		taskRedisClient: taskRedisClient,
		mqttClient:      mqttClient,
		influxClient:    influxClient,
		esClient:        esClient,
		llmClient:       llmClient,
		cancelBg:        cancelBg,
	}, nil
}

// ============================================================================
// Start — run server + graceful shutdown
// ============================================================================

func (srv *Server) Start() {
	srv.logger.Info("✅ starting server...")

	go func() {
		srv.logger.Infof("✅ Listening on %s", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srv.logger.Errorf("❌ server error: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	srv.logger.Infof("✅ Shutting down server... Reason: %s", sig)

	// ========================================================================
	// Graceful shutdown sequence:
	//   1. Stop accepting new HTTP requests
	//   2. Cancel background tasks (PDPA dispatcher, FullSchedule scheduler, ingester)
	//   3. Run registered cleanup hooks (PDPA Kafka producer, etc.)
	//   4. Disconnect MQTT + close InfluxDB
	// ========================================================================
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	// 1. HTTP server shutdown (drain in-flight requests)
	if err := srv.server.Shutdown(shutdownCtx); err != nil {
		srv.logger.Errorf("❌ HTTP server shutdown error: %v", err)
	} else {
		srv.logger.Info("✅ HTTP server stopped")
	}

	// 2. Cancel background tasks (uses parent ctx)
	if srv.cancelBg != nil {
		srv.cancelBg()
		srv.logger.Info("✅ Background tasks cancelled")
	}

	// 3. Run registered cleanup hooks (PDPA Kafka producer flush, etc.)
	if srv.runtime != nil {
		if err := srv.runtime.Shutdown(shutdownCtx); err != nil {
			srv.logger.Errorf("❌ Runtime cleanup error: %v", err)
		} else {
			srv.logger.Info("✅ Runtime cleanup completed")
		}
	}

	// 4. MQTT disconnect
	if srv.mqttClient != nil {
		srv.mqttClient.Disconnect(250)
		srv.logger.Info("✅ MQTT client disconnected")
	}

	// 5. InfluxDB close
	if srv.influxClient != nil {
		srv.influxClient.Close()
		srv.logger.Info("✅ InfluxDB client closed")
	}

	srv.logger.Info("✅ Server gracefully stopped")
}