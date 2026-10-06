package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"icmongolang/config"
	apiMiddleware "icmongolang/internal/middleware"
	alarmHttp "icmongolang/internal/modules/alarm/delivery/http"
	alarmRepo "icmongolang/internal/modules/alarm/repository"
	alarmUsecase "icmongolang/internal/modules/alarm/usecase"
	apimanagerHttp "icmongolang/internal/modules/apimanager/delivery/http"
	apimanagerRepo "icmongolang/internal/modules/apimanager/repository"
	apimanagerUsecase "icmongolang/internal/modules/apimanager/usecase"
	auditlogHttp "icmongolang/internal/modules/auditlog/delivery/http"
	auditlogProvider "icmongolang/internal/modules/auditlog/provider"
	auditlogRepo "icmongolang/internal/modules/auditlog/repository"
	auditlogUsecase "icmongolang/internal/modules/auditlog/usecase"
	authHttp "icmongolang/internal/modules/auth/delivery/http"
	authDomain "icmongolang/internal/modules/auth/domain"
	batchHttp "icmongolang/internal/modules/batch/delivery/http"
	batchPgRepo "icmongolang/internal/modules/batch/repository"
	batchUsecase "icmongolang/internal/modules/batch/usecase"
	controlHttp "icmongolang/internal/modules/control/delivery/http"
	controlRepo "icmongolang/internal/modules/control/repository"
	controlUsecase "icmongolang/internal/modules/control/usecase"
	customerHttp "icmongolang/internal/modules/customer/delivery/http"
	customerPgRepo "icmongolang/internal/modules/customer/repository"
	customerUsecase "icmongolang/internal/modules/customer/usecase"
	dashboardHttp "icmongolang/internal/modules/dashboard/delivery/http"
	dashboardRepo "icmongolang/internal/modules/dashboard/repository"
	dashboardUsecase "icmongolang/internal/modules/dashboard/usecase"
	documentHttp "icmongolang/internal/modules/document/delivery/http"
	documentRepo "icmongolang/internal/modules/document/repository"
	documentUsecase "icmongolang/internal/modules/document/usecase"
	elasticsearchHttp "icmongolang/internal/modules/elasticsearch/delivery/http"
	elasticsearchUsecase "icmongolang/internal/modules/elasticsearch/usecase"
	emailHttp "icmongolang/internal/modules/email/delivery/http"
	emailPgRepo "icmongolang/internal/modules/email/repository"
	emailUsecase "icmongolang/internal/modules/email/usecase"
	flowengineHttp "icmongolang/internal/modules/flowengine/delivery/http"
	flowengineRepo "icmongolang/internal/modules/flowengine/repository"
	flowengineUsecase "icmongolang/internal/modules/flowengine/usecase"
	fsHttp "icmongolang/internal/modules/fullschedule/delivery/http"
	fsExecutor "icmongolang/internal/modules/fullschedule/executor"
	fsPgRepo "icmongolang/internal/modules/fullschedule/repository"
	fsScheduler "icmongolang/internal/modules/fullschedule/scheduler"
	fsUsecase "icmongolang/internal/modules/fullschedule/usecase"
	i18nHttp "icmongolang/internal/modules/i18n/delivery/http"
	i18nPgRepo "icmongolang/internal/modules/i18n/repository"
	i18nUsecase "icmongolang/internal/modules/i18n/usecase"
	influxHttp "icmongolang/internal/modules/influxdb/delivery/http"
	influxUsecase "icmongolang/internal/modules/influxdb/usecase"
	iotHttp "icmongolang/internal/modules/iot/delivery/http"
	iotRepo "icmongolang/internal/modules/iot/repository"
	iotUsecase "icmongolang/internal/modules/iot/usecase"
	itemHttp "icmongolang/internal/modules/items/delivery/http"
	itemRepository "icmongolang/internal/modules/items/repository"
	itemUseCase "icmongolang/internal/modules/items/usecase"
	jobHttp "icmongolang/internal/modules/job/delivery/http"
	jobPgRepo "icmongolang/internal/modules/job/repository"
	jobUsecase "icmongolang/internal/modules/job/usecase"
	keyHttp "icmongolang/internal/modules/keys/delivery/http"
	keyRepo "icmongolang/internal/modules/keys/repository"
	keyUsecase "icmongolang/internal/modules/keys/usecase"
	mqttHttp "icmongolang/internal/modules/mqtt/delivery/http"
	mqttUsecase "icmongolang/internal/modules/mqtt/usecase"
	notifierHttp "icmongolang/internal/modules/notifier/delivery/http"
	notifierRepo "icmongolang/internal/modules/notifier/repository"
	notifierUsecase "icmongolang/internal/modules/notifier/usecase"
	paymentHttp "icmongolang/internal/modules/payment/delivery/http"
	paymentPgRepo "icmongolang/internal/modules/payment/repository"
	paymentUsecase "icmongolang/internal/modules/payment/usecase"
	pdpamod "icmongolang/internal/modules/pdpa"
	pdpahttp "icmongolang/internal/modules/pdpa/interfaces/http"
	purchaseorderHttp "icmongolang/internal/modules/purchaseorder/delivery/http"
	purchaseorderRepo "icmongolang/internal/modules/purchaseorder/repository"
	purchaseorderUsecase "icmongolang/internal/modules/purchaseorder/usecase"
	queue "icmongolang/internal/modules/queue"
	quotationHttp "icmongolang/internal/modules/quotation/delivery/http"
	quotationPgRepo "icmongolang/internal/modules/quotation/repository"
	quotationUsecase "icmongolang/internal/modules/quotation/usecase"
	realtimeHttp "icmongolang/internal/modules/realtime/delivery/http"
	realtimeUsecase "icmongolang/internal/modules/realtime/usecase"
	reportHttp "icmongolang/internal/modules/report/delivery/http"
	reportUsecase "icmongolang/internal/modules/report/usecase"
	settingsHttp "icmongolang/internal/modules/settings/delivery/http"
	settingsRepoPkg "icmongolang/internal/modules/settings/repository"
	settingsUCPkg "icmongolang/internal/modules/settings/usecase"
	userHttp "icmongolang/internal/modules/users/delivery/http"
	userDistributor "icmongolang/internal/modules/users/distributor"
	userRepository "icmongolang/internal/modules/users/repository"
	userUseCase "icmongolang/internal/modules/users/usecase"
	vectordataHttp "icmongolang/internal/modules/vectordata/delivery/http"
	vectordataUsecase "icmongolang/internal/modules/vectordata/usecase"
	wsdelivery "icmongolang/internal/modules/websocket/delivery/ws"
	wosHttp "icmongolang/internal/modules/wos/delivery/http"
	wosPgRepo "icmongolang/internal/modules/wos/repository"
	wosUsecase "icmongolang/internal/modules/wos/usecase"
	"icmongolang/internal/repository"
	redisDb "icmongolang/pkg/db/redis"
	"icmongolang/pkg/elasticsearch"
	"icmongolang/pkg/influxdb"
	"icmongolang/pkg/llm"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
	"icmongolang/pkg/sendEmail"
	"icmongolang/pkg/vectordb"
	"icmongolang/pkg/websocket"

	_ "icmongolang/docs"
	httpSwagger "icmongolang/pkg/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ============================================================================
// Middleware composition helpers
// ============================================================================

// composeMiddleware รวม middleware หลายตัวเป็น middleware ตัวเดียว
// ลำดับการเรียก: m1 → m2 → m3 → handler
func composeMiddleware(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			if mws[i] != nil {
				next = mws[i](next)
			}
		}
		return next
	}
}

// ============================================================================
// New — wire router + return Runtime with lifecycle hooks
// ============================================================================

// New wires the main application router and returns a Runtime that carries
// shutdown hooks. parentCtx controls background tasks (PDPA outbox dispatcher,
// FullSchedule scheduler, IoT ingester) — they stop when parentCtx is cancelled.
func New(
	parentCtx context.Context,
	db *gorm.DB,
	redisClient *redis.Client,
	taskRedisClient *asynq.Client,
	cfg *config.Config,
	logger logger.Logger,
	influxClient *influxdb.InfluxClient,
	mqttClient mqtt.Client,
	esClient *elasticsearch.Client,
	llmClient *llm.Client,
) (*Runtime, error) {

	// ========================================================================
	// 1. Repositories and UseCases
	// ========================================================================
	// --- User module: repos ---
	userPgRepo     := userRepository.CreateUserPgRepository(db)
	userRedisRepo  := userRepository.CreateUserRedisRepository(redisClient)
	accessMenuRepo := userRepository.CreateAccessMenuPgRepository(db)
	fileRepo       := userRepository.CreateFilePgRepository(db)
	logRepo        := userRepository.CreateLogPgRepository(db)
	logTypeRepo    := userRepository.CreateLogTypePgRepository(db)
	moduleRepo     := userRepository.CreateModulePgRepository(db)
	menuRepo       := userRepository.CreateMenuPgRepository(db)
	cacheRepo      := userRepository.CreateCacheRepository(redisClient)

	// --- User module: usecases ---
	userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, logger)

	userUC       := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, logger)
	accessMenuUC := userUseCase.NewAccessMenuUseCase(accessMenuRepo, cacheRepo, logger)
	fileUC       := userUseCase.NewFileUseCase(fileRepo, cacheRepo, logger)
	logUC        := userUseCase.NewLogUseCase(logRepo, cacheRepo, logger)
	logTypeUC    := userUseCase.NewLogTypeUseCase(logTypeRepo, cacheRepo, logger)
	moduleUC     := userUseCase.NewModuleUseCase(moduleRepo, cacheRepo, logger)
	menuUC       := userUseCase.NewMenuUseCase(menuRepo, cacheRepo, logger)

	// --- Item module ---
	itemPgRepo := itemRepository.CreateItemPgRepository(db)
	itemUC     := itemUseCase.CreateItemUseCaseI(itemPgRepo, cfg, logger)

	// --- Business modules ---
	poPgRepo := purchaseorderRepo.CreatePurchaseOrderPgRepository(db)
	poUC     := purchaseorderUsecase.CreatePurchaseOrderUseCaseI(poPgRepo, cfg, logger)

	paymentRepoVar := paymentPgRepo.CreatePaymentPgRepository(db)
	paymentUC      := paymentUsecase.CreatePaymentUseCaseI(paymentRepoVar, cfg, logger)

	dashRepoVar := dashboardRepo.CreateDashboardPgRepository(db)
	dashUC      := dashboardUsecase.CreateDashboardUseCaseI(dashRepoVar, cfg, logger)

	docRepoVar := documentRepo.CreateDocumentPgRepository(db)
	docUC      := documentUsecase.CreateDocumentUseCaseI(docRepoVar, cfg, logger)

	emailRepoVar := emailPgRepo.CreateEmailPgRepository(db)
	emailUC      := emailUsecase.CreateEmailUseCaseI(emailRepoVar, cfg, logger)

	batchRepoVar := batchPgRepo.CreateBatchPgRepository(db)
	batchUC      := batchUsecase.CreateBatchUseCaseI(batchRepoVar, cfg, logger)

	i18nRepoVar := i18nPgRepo.CreateI18nPgRepository(db)
	i18nUC      := i18nUsecase.CreateI18nUseCaseI(i18nRepoVar, cfg, logger)

	wosRepoVar := wosPgRepo.CreateWosPgRepository(db)
	wosUC      := wosUsecase.CreateWosUseCaseI(wosRepoVar, cfg, logger)

	jobRepoVar := jobPgRepo.CreateJobPgRepository(db)
	jobUC      := jobUsecase.CreateJobUseCaseI(jobRepoVar, cfg, logger)

	customerPgrepo := customerPgRepo.CreateCustomerPgRepository(db)
	carPgrepo      := customerPgRepo.CreateCarPgRepository(db)
	customerUC     := customerUsecase.CreateCustomerUseCaseI(customerPgrepo, cfg, logger)
	carUC          := customerUsecase.CreateCarUseCaseI(carPgrepo, cfg, logger)

	quotationRepoVar := quotationPgRepo.CreateQuotationPgRepository(db)
	quotationUC      := quotationUsecase.CreateQuotationUseCaseI(quotationRepoVar, cfg, logger)

	// ========================================================================
	// 2. Middleware Manager
	// ========================================================================
	redisRepo := repository.CreateRedisRepository[any](redisClient)
	mw := apiMiddleware.CreateMiddlewareManager(cfg, logger, userUC, redisRepo)
	apiMiddleware.StartMetricsCollector(redisClient)

	// ========================================================================
	// 3. Router + Global middlewares
	// ========================================================================
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(apiMiddleware.MonitoringMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(time.Duration(cfg.Server.ProcessTimeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(mw.Cors()))

	// Static avatar uploads
	uploadRoot := cfg.Upload.Dir
	if uploadRoot == "" {
		uploadRoot = "./uploads"
	}
	avatarFS := http.FileServer(http.Dir(filepath.Join(uploadRoot, "avatar")))
	r.Handle("/uploads/avatar/*", http.StripPrefix("/uploads/avatar/", avatarFS))

	// API sub-router
	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	// ─────────────────────────────────────────────────────────────────────────
	// Runtime — bundle router + shutdown hooks
	// ─────────────────────────────────────────────────────────────────────────
	rt := &Runtime{Router: r}

	// ========================================================================
	// 4. Core handlers
	// ========================================================================
	userHandler := userHttp.CreateUserHandler(
		userUC,
		accessMenuUC,
		fileUC,
		logUC,
		logTypeUC,
		menuUC,
		moduleUC,
		cfg,
		logger,
	)
	authHandler := authHttp.CreateAuthHandler(userUC, logUC, cfg, logger)
	itemHandler := itemHttp.CreateItemHandler(itemUC, cfg, logger)
	poHandler := purchaseorderHttp.CreatePurchaseOrderHandler(poUC, cfg, logger)
	paymentHandler := paymentHttp.CreatePaymentHandler(paymentUC, cfg, logger)
	dashHandler := dashboardHttp.CreateDashboardHandler(dashUC, cfg, logger)
	docHandler := documentHttp.CreateDocumentHandler(docUC, cfg, logger)
	emailHandler := emailHttp.CreateEmailHandler(emailUC, cfg, logger)
	batchHandler := batchHttp.CreateBatchHandler(batchUC, cfg, logger)
	i18nHandler := i18nHttp.CreateI18nHandler(i18nUC, cfg, logger)
	wosHandler := wosHttp.CreateWosHandler(wosUC, cfg, logger)
	jobHandler := jobHttp.CreateJobHandler(jobUC, cfg, logger)
	customerHandler := customerHttp.CreateCustomerHandler(customerUC, carUC, cfg, logger)
	quotationHandler := quotationHttp.CreateQuotationHandler(quotationUC, cfg, logger)
	reportHandler := reportHttp.CreateReportHandler(cfg, logger, customerUC, reportUsecase.NewReportUseCase(customerUC, quotationUC))

	// ========================================================================
	// 5. InfluxDB routes (optional)
	// ========================================================================
	if influxClient != nil {
		influxUC := influxUsecase.NewInfluxUseCase(influxClient, logger)
		influxHandler := influxHttp.CreateInfluxHandler(influxUC, logger)
		influxHttp.MapInfluxRoutes(apiRouter, influxHandler, mw)
		logger.Info("✅ InfluxDB routes registered")
	} else {
		logger.Warn("⚠️ InfluxDB client is nil – Influx routes skipped")
	}

	// ========================================================================
	// 5.1 Elasticsearch + LLM vector search (optional)
	// ========================================================================
	if esClient != nil && llmClient != nil {
		esUC := elasticsearchUsecase.NewElasticsearchUseCase(esClient, llmClient, logger)
		esHandler := elasticsearchHttp.NewElasticsearchHandler(esUC, logger)
		elasticsearchHttp.MapElasticsearchRoutes(apiRouter, esHandler, mw)
		logger.Info("✅ Elasticsearch routes registered (vector search)")
	} else {
		logger.Warn("⚠️ Elasticsearch client is nil – ES routes skipped")
	}

	// ========================================================================
	// 5.2 VectorData routes (abstract vector DB)
	// ========================================================================
	vdb, err := vectordb.New(&cfg.VectorDB, esClient, db, logger)
	if err != nil {
		logger.Errorf("❌ Failed to init vector DB: %v", err)
	} else {
		vdUC := vectordataUsecase.NewVectorDataUseCase(vdb, llmClient, cfg, logger)
		vdHandler := vectordataHttp.NewVectorDataHandler(vdUC, logger)
		vectordataHttp.MapVectorDataRoutes(apiRouter, vdHandler, mw)
		logger.Infof("✅ VectorData routes registered (provider=%s, index=%s)", vdb.Provider(), cfg.VectorDB.Index)
	}

	// ========================================================================
	// 6. WebSocket Hub
	// ========================================================================
	wsHub := websocket.NewHub(websocket.HubConfig{
		Queue: queue.NewQueue(),
		Log:   logger,
		Callbacks: websocket.HubCallbacks{
			OnConnect: func(client *websocket.Client) error {
				logger.Infof("WS client connected: %s (user: %s, role: %s)", client.ID, client.UserID, client.Role)
				return nil
			},
			OnDisconnect: func(client *websocket.Client) {
				logger.Infof("WS client disconnected: %s (user: %s)", client.ID, client.UserID)
			},
			OnSubscribe:   func(client *websocket.Client, topic string) { logger.Debugf("WS sub: %s -> %s", client.ID, topic) },
			OnUnsubscribe: func(client *websocket.Client, topic string) { logger.Debugf("WS unsub: %s -> %s", client.ID, topic) },
			OnJoinRoom:    func(client *websocket.Client, room string) { logger.Debugf("WS join: %s -> %s", client.ID, room) },
			OnLeaveRoom:   func(client *websocket.Client, room string) { logger.Debugf("WS leave: %s -> %s", client.ID, room) },
			OnMessage: func(client *websocket.Client, topic, room string, payload []byte) bool {
				logger.Debugf("WS msg: %s topic=%s room=%s", client.ID, topic, room)
				return true
			},
		},
	})
	go wsHub.Run()

	// ========================================================================
	// 7. Redis Cache
	// ========================================================================
	alarmLogRepo := alarmRepo.NewAlarmLogRepository(db)
	redisCache := redisDb.NewCache(redisClient)

	// ========================================================================
	// 8. IoT MQTT v3 routes
	// ========================================================================
	deviceRepo := iotRepo.NewDeviceRepository(db)
	iotAlarmLogRepo := iotRepo.NewAlarmLogRepository(db)
	deviceStatusRepo := iotRepo.NewDeviceStatusRepository(db)
	deviceConfigRepo := iotRepo.NewDeviceConfigRepository(db)
	iotDataRepo := iotRepo.NewIotDataRepository(db)
	activityLogRepo := iotRepo.NewActivityLogRepository(db)
	commandLogRepo := iotRepo.NewCommandLogRepository(db)
	deviceAlertRepo := iotRepo.NewDeviceAlertRepository(db)

	iotUC := iotUsecase.NewMQTT3UseCase(
		deviceRepo, iotAlarmLogRepo, mqttClient, redisClient, influxClient,
		logger, cfg,
		deviceStatusRepo, deviceConfigRepo, iotDataRepo,
		activityLogRepo, commandLogRepo, deviceAlertRepo,
	)
	if mqttClient != nil && mqttClient.IsConnected() {
		if err := iotUC.StartIngest(parentCtx); err != nil {
			logger.Errorf("❌ Failed to start IoT ingester: %v", err)
		} else {
			logger.Info("✅ IoT MQTT ingester started (topic #)")
		}
	} else {
		logger.Warn("⚠️ MQTT client not connected – IoT MQTT ingester skipped")
	}
	iotHandler := iotHttp.NewMQTT3Handler(iotUC, logger, redisCache)
	if err := iotHttp.MapMQTT3Routes(apiRouter, iotHandler, mw); err != nil {
		return nil, fmt.Errorf("register IoT MQTT3 routes: %w", err)
	}
	logger.Info("✅ MQTT3 routes registered")

	// ========================================================================
	// 9. MQTT (standard) routes
	// ========================================================================
	if mqttClient != nil && mqttClient.IsConnected() {
		var influxDBClient influxdb.Client
		if influxClient != nil {
			influxDBClient = influxClient
		}
		mqttUC, err := mqttUsecase.NewMQTTUseCaseWithClientAndWS(
			mqttClient, cfg, logger, wsHub, db,
			alarmLogRepo, influxDBClient, redisCache,
		)
		if err != nil {
			logger.Errorf("❌ Failed to create MQTT usecase: %v", err)
		} else {
			mqttHandler := mqttHttp.CreateMQTTHandler(mqttUC, logger, redisCache)
			mqttHttp.MapMQTTRoutes(apiRouter, mqttHandler, mw)
			logger.Info("✅ MQTT routes registered with WebSocket, alarm log, cache & InfluxDB")
		}
	} else {
		logger.Warn("⚠️ MQTT client not connected – MQTT routes skipped")
	}

	// ========================================================================
	// 10. Alarm routes
	// ========================================================================
	alarmUC := alarmUsecase.NewAlarmUseCase()
	alarmHandler := alarmHttp.NewAlarmHandler(alarmUC, logger)
	alarmHttp.MapAlarmRoutes(apiRouter, alarmHandler, mw)

	// ========================================================================
	// 11. Core routes (auth, user, item)
	// ========================================================================
	authHttp.MapAuthRoute(apiRouter, authHandler, mw)
	userHttp.MapUserRoute(apiRouter, userHandler, mw)
	itemHttp.MapItemRoute(apiRouter, itemHandler, mw)

	// ========================================================================
	// 12. Business module routes
	// ========================================================================
	purchaseorderHttp.MapPurchaseOrderRoute(apiRouter, poHandler, mw)
	paymentHttp.MapPaymentRoute(apiRouter, paymentHandler, mw)
	paymentHttp.MapReceiptRoute(apiRouter, paymentHandler, mw)
	dashboardHttp.MapDashboardRoute(r, dashHandler, mw)
	documentHttp.MapDocumentRoute(r, docHandler, mw)
	emailHttp.MapEmailRoute(r, emailHandler, mw)
	batchHttp.MapBatchRoute(r, batchHandler, mw)
	i18nHttp.MapI18nRoute(r, i18nHandler, mw)
	wosHttp.MapWosRoute(r, wosHandler, mw)
	jobHttp.MapJobRoute(apiRouter, jobHandler, mw)
	customerHttp.MapCustomerRoute(apiRouter, customerHandler, mw)
	quotationHttp.MapQuotationRoute(apiRouter, quotationHandler, mw)
	reportHttp.MapReportRoute(r, cfg, reportHandler, mw)
	logger.Info("✅ New module routes registered")

	// ========================================================================
	// 13. Settings module
	// ========================================================================
	settingsRepo := settingsRepoPkg.CreateSettingsPgRepository(db)
	settingsUC := settingsUCPkg.CreateSettingsUseCaseI(settingsRepo, mqttClient, redisCache, logger)
	settingsHandler := settingsHttp.CreateSettingsHandler(settingsUC, cfg, logger)
	settingsHttp.MapSettingRoute(apiRouter, settingsHandler, mw)
	logger.Info("✅ Settings routes registered")

	// ========================================================================
	// 14. Platform modules (Phase 0 skeletons)
	// ========================================================================
	notifierRepoVar := notifierRepo.NewPgRepository(db)
	notifierUC := notifierUsecase.NewNotifierUseCase(notifierRepoVar, cfg, mqttClient, logger)
	notifierHandler := notifierHttp.CreateNotifierHandler(notifierUC, logger)
	notifierHttp.MapNotifierRoutes(apiRouter, notifierHandler, mw)
	logger.Info("✅ Notifier routes registered")

	auditlogRepoVar := auditlogRepo.NewPgRepository(db)
	auditlogAnchor := auditlogProvider.NewNoopProvider()
	auditlogUC := auditlogUsecase.NewAuditLogUseCase(auditlogRepoVar, auditlogAnchor, logger)
	auditlogHandler := auditlogHttp.CreateAuditLogHandler(auditlogUC, logger)
	auditlogHttp.MapAuditLogRoutes(apiRouter, auditlogHandler, mw)
	logger.Info("✅ AuditLog routes registered (hash-chain)")

	// --- API Keys module (pkg/apikey backed) ---
	// Configure the API key pepper through the API_KEY_PEPPER environment variable.
	keyRepoVar := keyRepo.NewPostgresRepository(db)
	keyUC := keyUsecase.NewAPIKeyUseCase(keyRepoVar, []byte(os.Getenv("API_KEY_PEPPER")))
	keyHandler := keyHttp.NewHandler(keyUC)
	keyHttp.MapRoutes(apiRouter, keyHandler, mw)
	logger.Info("✅ API keys routes registered")

	// --- APIManager: wire rate-limit cache adapter ---
	apiCache := &apiCacheAdapter{cache: redisCache, redis: redisClient}
	apimanagerRepoVar := apimanagerRepo.NewPgRepository(db)
	apimanagerUC := apimanagerUsecase.NewAPIManagerUseCase(apimanagerRepoVar, apiCache, logger)
	apimanagerHandler := apimanagerHttp.CreateAPIManagerHandler(apimanagerUC, logger)
	apimanagerHttp.MapAPIManagerRoutes(apiRouter, apimanagerHandler, mw)
	logger.Info("✅ APIManager routes registered (key + rate-limit)")

	realtimeUC := realtimeUsecase.NewRealtimeUseCase(wsHub, logger)
	realtimeHandler := realtimeHttp.CreateRealtimeHandler(realtimeUC, logger)
	realtimeHttp.MapRealtimeRoutes(apiRouter, realtimeHandler, mw)
	logger.Info("✅ Realtime routes registered (dashboard real-time push)")

	flowengineRepoVar := flowengineRepo.NewPgRepository(db)
	flowengineUC := flowengineUsecase.NewFlowEngineUseCase(flowengineRepoVar, notifierUC, realtimeUC, logger)
	flowengineHandler := flowengineHttp.CreateFlowEngineHandler(flowengineUC, logger)
	flowengineHttp.MapFlowEngineRoutes(apiRouter, flowengineHandler, mw)
	logger.Info("✅ FlowEngine routes registered")

	// ========================================================================
	// 15. Control module
	// ========================================================================
	controlRepoVar := controlRepo.NewPgRepository(db)
	controlUC := controlUsecase.NewControlUseCase(controlRepoVar, mqttClient, logger)
	controlHandler := controlHttp.CreateControlHandler(controlUC, logger)
	controlHttp.MapControlRoutes(apiRouter, controlHandler, mw)
	logger.Info("✅ Control routes registered (auto-control manager)")

	// ========================================================================
	// 16. FullSchedule module
	// ========================================================================
	fsRepoVar := fsPgRepo.CreateFullSchedulePgRepository(db)
	fsMasterRepo := fsPgRepo.CreateMasterPgRepository(db)
	fsMasterUC := fsUsecase.CreateMasterUseCaseI(fsMasterRepo, cfg, logger)
	fsExec := fsExecutor.New(mqttClient, fsRepoVar, sendEmail.NewEmailSender(cfg), cfg, logger)
	fsUC := fsUsecase.CreateFullScheduleUseCaseI(fsRepoVar, fsExec, fsMasterUC, cfg, logger)
	fsHandler := fsHttp.CreateFullscheduleHandler(fsUC, cfg, logger)
	fsMasterHandler := fsHttp.CreateMasterHandler(fsMasterUC, cfg, logger)
	fsHttp.MapFullscheduleRoutes(apiRouter, fsHandler, mw)
	fsHttp.MapMasterRoutes(apiRouter, fsMasterHandler, mw)
	fsSched := fsScheduler.New(fsUC, redisClient, logger, 0)
	fsSched.Start(parentCtx)
	logger.Info("✅ FullSchedule routes registered + scheduler started")

	// ========================================================================
	// 17. PDPA module (compliance: consent, DSAR, account, policy, audit)
	// ========================================================================
	pdpaMod, err := pdpamod.NewModule(db, redisClient, cfg, logger)
	if err != nil {
		logger.Errorf("❌ Failed to init PDPA module: %v", err)
		// ไม่ Fatal — PDPA optional
	} else {
		// Start outbox dispatcher bound to parent ctx (auto-stop on shutdown)
		pdpaMod.StartBackground(parentCtx)

		// Register cleanup hook — Kafka producer flush + dispatcher stop
		rt.AddShutdown(func(ctx context.Context) error {
			logger.Info("⏳ Shutting down PDPA module...")
			return pdpaMod.Shutdown()
		})

		// HTTP handlers
		pdpaHandlers := &pdpahttp.Handlers{
			Consent: pdpahttp.NewConsentHandler(
				pdpaMod.GrantConsent,
				pdpaMod.RevokeConsent,
				pdpaMod.ListConsentLog,
				logger,
			),
			DSAR: pdpahttp.NewDSARHandler(
				pdpaMod.SubmitDSAR,
				pdpaMod.ProcessDSAR,
				pdpaMod.ListDSAR,
				pdpaMod.GetDSARStatus,
				logger,
			),
			Account: pdpahttp.NewAccountHandler(
				pdpaMod.SuspendAccount,
				pdpaMod.TerminateAccount,
				pdpaMod.ConfirmDeletion,
				logger,
			),
			Policy: pdpahttp.NewPolicyHandler(
				pdpaMod.GetActivePolicy,
				logger,
			),
			Admin: pdpahttp.NewAdminHandler(
				pdpaMod.GetAdminReport,
				logger,
			),
		}

		// Auth middleware: verifier + authenticator + currentUser + activeUser
		pdpaAuthMW := pdpahttp.AuthMiddleware(composeMiddleware(
			mw.Verifier(true),
			mw.Authenticator(),
			mw.CurrentUser(),
			mw.ActiveUser(),
		))

		// Admin middleware: role check only (auth already applied by outer group)
		pdpaAdminMW := pdpahttp.AuthMiddleware(mw.RequireRole(authDomain.RoleAdmin))

		pdpahttp.RegisterRoutes(apiRouter, pdpaHandlers, pdpaAuthMW, pdpaAdminMW)
		logger.Info("✅ PDPA routes registered (consent, DSAR, account, policy, audit)")
	}

	// ========================================================================
	// 18. WebSocket endpoint
	// ========================================================================
	wsHandler := wsdelivery.NewWsHandler(wsHub, nil, logger, mw)
	apiRouter.Get("/ws", mw.WSOriginGuard(cfg.WS.AllowedOrigins)(wsHandler.HandlerFunc()))

	// ========================================================================
	// 19. Metrics + Health
	// ========================================================================
	r.Get("/apimetric", apiMiddleware.MetricsHandler)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// ========================================================================
	// 20. Swagger
	// ========================================================================
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// ========================================================================
	// 21. Detailed health check
	// ========================================================================
	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, map[string]string{"status": "ok", "timestamp": time.Now().UTC().String()})
	})

	apiRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"status": "up",
			"checks": map[string]bool{
				"database":      db != nil,
				"redis":         redisClient != nil,
				"influxdb":      influxClient != nil,
				"elasticsearch": esClient != nil,
				"vectordata":    vdb != nil,
				"llm":           llmClient != nil,
				"mqtt":          mqttClient != nil && mqttClient.IsConnected(),
				"websocket":     true,
				"pdpa":          pdpaMod != nil,
			},
		}
		render.Respond(w, r, status)
	})

	return rt, nil
}