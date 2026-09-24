package server

import (
	"context"
	"fmt"
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
	ctxTimeout     = 5
)

type Server struct {
	server          *http.Server
	cfg             *config.Config
	logger          logger.Logger
	db              *gorm.DB
	redisClient     *redis.Client
	taskRedisClient *asynq.Client
	mqttClient      mqtt.Client
	influxClient    *influxdb.InfluxClient
	esClient        *elasticsearch.Client
	llmClient       *llm.Client
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

func NewServer(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, taskRedisClient *asynq.Client, log logger.Logger) (*Server, error) {
	log.Info("✅ configuring server...")

	// ---------- สร้าง MQTT client (อาจเชื่อมไม่สำเร็จ) ----------
	mqttClient := mqtt.New(&cfg.MQTT, log)

	connTimeout := time.Duration(cfg.MQTT.ConnectionTimeout) * time.Second
	if connTimeout == 0 {
		connTimeout = 30 * time.Second
		log.Warn("✅ MQTT connection timeout not set in config, using default 30s")
	}

	if strings.HasPrefix(cfg.MQTT.Broker, "mqtt://") {
		log.Warn("✅ Broker URL uses 'mqtt://' scheme; Paho Go client requires 'tcp://' or 'ssl://'. Connection may fail.")
	}

	if err := connectMQTTWithRetry(mqttClient, connTimeout, 3, log); err != nil {
		log.Errorf("❌ Failed to connect MQTT: %v", err)
		// mqttClient จะอยู่ในสถานะ disconnected, routes จะไม่ถูกเพิ่ม (handlers.go จะตรวจสอบ IsConnected)
	} else {
		log.Info("✅ MQTT connected successfully")
	}

	// ---------- สร้าง InfluxDB client (อาจเป็น nil ถ้าเชื่อมไม่สำเร็จ) ----------
	influxClient, err := influxdb.NewInfluxClient(cfg, log)
	if err != nil {
		log.Errorf("❌ Failed to connect InfluxDB: %v", err)
		influxClient = nil
	} else {
		log.Info("✅ InfluxDB connected successfully")
	}

	// ---------- สร้าง Elasticsearch client (อาจเป็น nil ถ้าเชื่อมไม่สำเร็จ) ----------
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

	// ---------- สร้าง LLM embedding client (พร้อมใช้งานเสมอ) ----------
	llmClient := llm.NewClient(&cfg.LLM, log)

	// ---------- สร้าง router โดยส่ง clients ที่สร้างแล้ว ----------
	router, err := New(db, redisClient, taskRedisClient, cfg, log, influxClient, mqttClient, esClient, llmClient)
	if err != nil {
		return nil, fmt.Errorf("⚠️ Failed to create router: %w", err)
	}

	addr := cfg.Server.Port
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	return &Server{
		server: &http.Server{
			Addr:           addr,
			Handler:        router,
			ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
			MaxHeaderBytes: maxHeaderBytes,
		},
		cfg:             cfg,
		logger:          log,
		db:              db,
		redisClient:     redisClient,
		taskRedisClient: taskRedisClient,
		mqttClient:      mqttClient,
		influxClient:    influxClient,
		esClient:        esClient,
		llmClient:       llmClient,
	}, nil
}

func (srv *Server) Start() {
	srv.logger.Info("✅ starting server...")
	go func() {
		srv.logger.Infof("✅ Listening on %s", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
			srv.logger.Errorf("❌ server error: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	srv.logger.Infof("✅ Shutting down server... Reason: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Errorf("❌ server shutdown error: %v", err)
	}

	if srv.mqttClient != nil {
		srv.mqttClient.Disconnect(250)
		srv.logger.Info("✅ MQTT client disconnected")
	}
	if srv.influxClient != nil {
		srv.influxClient.Close()
		srv.logger.Info("✅ InfluxDB client closed")
	}

	srv.logger.Info("✅ Server gracefully stopped")
}
