package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"icmongolang/config"
	apiMiddleware "icmongolang/internal/middleware"
	authHttp "icmongolang/internal/modules/auth/delivery/http"
	batchHttp "icmongolang/internal/modules/batch/delivery/http"
	batchPgRepo "icmongolang/internal/modules/batch/repository"
	batchUsecase "icmongolang/internal/modules/batch/usecase"
	customerHttp "icmongolang/internal/modules/customer/delivery/http"
	customerPgRepo "icmongolang/internal/modules/customer/repository"
	customerUsecase "icmongolang/internal/modules/customer/usecase"
	dashboardHttp "icmongolang/internal/modules/dashboard/delivery/http"
	dashboardRepo "icmongolang/internal/modules/dashboard/repository"
	dashboardUsecase "icmongolang/internal/modules/dashboard/usecase"
	documentHttp "icmongolang/internal/modules/document/delivery/http"
	documentRepo "icmongolang/internal/modules/document/repository"
	documentUsecase "icmongolang/internal/modules/document/usecase"
	emailHttp "icmongolang/internal/modules/email/delivery/http"
	emailPgRepo "icmongolang/internal/modules/email/repository"
	emailUsecase "icmongolang/internal/modules/email/usecase"
	i18nHttp "icmongolang/internal/modules/i18n/delivery/http"
	i18nPgRepo "icmongolang/internal/modules/i18n/repository"
	i18nUsecase "icmongolang/internal/modules/i18n/usecase"
	itemHttp "icmongolang/internal/modules/items/delivery/http"
	itemRepo "icmongolang/internal/modules/items/repository"
	itemUC "icmongolang/internal/modules/items/usecase"
	jobHttp "icmongolang/internal/modules/job/delivery/http"
	jobPgRepo "icmongolang/internal/modules/job/repository"
	jobUsecase "icmongolang/internal/modules/job/usecase"
	paymentHttp "icmongolang/internal/modules/payment/delivery/http"
	paymentPgRepo "icmongolang/internal/modules/payment/repository"
	paymentUsecase "icmongolang/internal/modules/payment/usecase"
	purchaseorderHttp "icmongolang/internal/modules/purchaseorder/delivery/http"
	purchaseorderRepo "icmongolang/internal/modules/purchaseorder/repository"
	purchaseorderUsecase "icmongolang/internal/modules/purchaseorder/usecase"
	quotationHttp "icmongolang/internal/modules/quotation/delivery/http"
	quotationPgRepo "icmongolang/internal/modules/quotation/repository"
	quotationUsecase "icmongolang/internal/modules/quotation/usecase"
	reportHttp "icmongolang/internal/modules/report/delivery/http"
	reportUsecase "icmongolang/internal/modules/report/usecase"
	userHttp "icmongolang/internal/modules/users/delivery/http"
	userDistributor "icmongolang/internal/modules/users/distributor"
	userRepo "icmongolang/internal/modules/users/repository"
	userUC "icmongolang/internal/modules/users/usecase"
	wosHttp "icmongolang/internal/modules/wos/delivery/http"
	wosPgRepo "icmongolang/internal/modules/wos/repository"
	wosUsecase "icmongolang/internal/modules/wos/usecase"
	"icmongolang/pkg/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(db *gorm.DB, redisClient *redis.Client, taskRedisClient *asynq.Client, cfg *config.Config, log logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	// ---------- 1. Global middleware (order matters) ----------
	r.Use(middleware.Recoverer) // panic recovery
	r.Use(middleware.RequestID) // request ID header
	r.Use(middleware.RealIP)    // real client IP
	r.Use(middleware.Timeout(time.Duration(cfg.Server.ProcessTimeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: func() []string {
			if cfg.Server.Mode == "production" {
				return []string{"https://" + cfg.Server.BaseUrl}
			}
			return []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"}
		}(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Security headers
	secConfig := apiMiddleware.DefaultSecurityConfig
	r.Use(apiMiddleware.SecurityMiddleware(&secConfig))

	// Rate limiting
	rateLimiter := apiMiddleware.NewRateLimitMiddleware(apiMiddleware.DefaultRateLimitConfig)
	r.Use(rateLimiter.Handler)

	// Custom logging and monitoring
	r.Use(apiMiddleware.LoggingMiddleware)
	r.Use(apiMiddleware.MonitoringMiddleware)
	// ---------- 1. Health check Index----------
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ---------- 2. Metrics endpoint (Prometheus) ----------
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	// Also keep the custom JSON metrics endpoint if needed
	r.Get("/apimetric", apiMiddleware.MetricsHandler)

	// ---------- 3. Health check ----------
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ---------- 4. Setup repositories, usecases, handlers ----------
	// User module
	userPgRepo := userRepo.CreateUserPgRepository(db)
	userRedisRepo := userRepo.CreateUserRedisRepository(redisClient)
	userDistrib := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, log)
	userUsecase := userUC.CreateUserUseCaseI(userPgRepo, userRedisRepo, userDistrib, cfg, log)

	// Item module
	itemPgRepo := itemRepo.CreateItemPgRepository(db)
	itemUsecase := itemUC.CreateItemUseCaseI(itemPgRepo, cfg, log)

	// --- New module repos and usecases ---
	poPgRepo := purchaseorderRepo.CreatePurchaseOrderPgRepository(db)
	poUsecase := purchaseorderUsecase.CreatePurchaseOrderUseCaseI(poPgRepo, cfg, log)

	paymentPgRepoVar := paymentPgRepo.CreatePaymentPgRepository(db)
	paymentUsecaseVar := paymentUsecase.CreatePaymentUseCaseI(paymentPgRepoVar, cfg, log)

	dashPgRepoVar := dashboardRepo.CreateDashboardPgRepository(db)
	dashUsecaseVar := dashboardUsecase.CreateDashboardUseCaseI(dashPgRepoVar, cfg, log)

	docPgRepoVar := documentRepo.CreateDocumentPgRepository(db)
	docUsecaseVar := documentUsecase.CreateDocumentUseCaseI(docPgRepoVar, cfg, log)

	emailPgRepoVar := emailPgRepo.CreateEmailPgRepository(db)
	emailUsecaseVar := emailUsecase.CreateEmailUseCaseI(emailPgRepoVar, cfg, log)

	batchPgRepoVar := batchPgRepo.CreateBatchPgRepository(db)
	batchUsecaseVar := batchUsecase.CreateBatchUseCaseI(batchPgRepoVar, cfg, log)

	i18nPgRepoVar := i18nPgRepo.CreateI18nPgRepository(db)
	i18nUsecaseVar := i18nUsecase.CreateI18nUseCaseI(i18nPgRepoVar, cfg, log)

	wosPgRepoVar := wosPgRepo.CreateWosPgRepository(db)
	wosUsecaseVar := wosUsecase.CreateWosUseCaseI(wosPgRepoVar, cfg, log)

	jobPgRepoVar := jobPgRepo.CreateJobPgRepository(db)
	jobUsecaseVar := jobUsecase.CreateJobUseCaseI(jobPgRepoVar, cfg, log)

	custPgRepoVar := customerPgRepo.CreateCustomerPgRepository(db)
	carPgRepoVar := customerPgRepo.CreateCarPgRepository(db)
	customerUsecaseVar := customerUsecase.CreateCustomerUseCaseI(custPgRepoVar, cfg, log)
	carUsecaseVar := customerUsecase.CreateCarUseCaseI(carPgRepoVar, cfg, log)

	quotationPgRepoVar := quotationPgRepo.CreateQuotationPgRepository(db)
	quotationUsecaseVar := quotationUsecase.CreateQuotationUseCaseI(quotationPgRepoVar, cfg, log)

	// Middleware manager
	mw := apiMiddleware.CreateMiddlewareManager(cfg, log, userUsecase)

	// Handlers
	authHandler := authHttp.CreateAuthHandler(userUsecase, cfg, log)
	userHandler := userHttp.CreateUserHandler(userUsecase, cfg, log)
	itemHandler := itemHttp.CreateItemHandler(itemUsecase, cfg, log)
	poHandler := purchaseorderHttp.CreatePurchaseOrderHandler(poUsecase, cfg, log)
	paymentHandler := paymentHttp.CreatePaymentHandler(paymentUsecaseVar, cfg, log)
	dashHandler := dashboardHttp.CreateDashboardHandler(dashUsecaseVar, cfg, log)
	docHandler := documentHttp.CreateDocumentHandler(docUsecaseVar, cfg, log)
	emailHandler := emailHttp.CreateEmailHandler(emailUsecaseVar, cfg, log)
	batchHandler := batchHttp.CreateBatchHandler(batchUsecaseVar, cfg, log)
	i18nHandler := i18nHttp.CreateI18nHandler(i18nUsecaseVar, cfg, log)
	wosHandler := wosHttp.CreateWosHandler(wosUsecaseVar, cfg, log)
	jobHandler := jobHttp.CreateJobHandler(jobUsecaseVar, cfg, log)
	customerHandler := customerHttp.CreateCustomerHandler(customerUsecaseVar, carUsecaseVar, cfg, log)
	quotationHandler := quotationHttp.CreateQuotationHandler(quotationUsecaseVar, cfg, log)
	reportHandler := reportHttp.CreateReportHandler(cfg, log, customerUsecaseVar, reportUsecase.NewReportUseCase(customerUsecaseVar, quotationUsecaseVar))

	// ---------- 5. API routes ----------
	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	// Public routes
	apiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, "pong")
	})

	// Auth routes (public)
	authHttp.MapAuthRoute(apiRouter, authHandler, mw)

	// Protected routes group
	apiRouter.Group(func(r chi.Router) {
		r.Use(mw.Verifier(true))
		r.Use(mw.Authenticator())
		r.Use(mw.CurrentUser())
		r.Use(mw.ActiveUser())

		// User routes
		userHttp.MapUserRoute(apiRouter, userHandler, mw)

		// Item routes
		itemHttp.MapItemRoute(apiRouter, itemHandler, mw)

		// New module routes (compatible with apiRouter)
		purchaseorderHttp.MapPurchaseOrderRoute(apiRouter, poHandler, mw)
		paymentHttp.MapPaymentRoute(apiRouter, paymentHandler, mw)
		paymentHttp.MapReceiptRoute(apiRouter, paymentHandler, mw)
		jobHttp.MapJobRoute(apiRouter, jobHandler, mw)
		customerHttp.MapCustomerRoute(apiRouter, customerHandler, mw)
		quotationHttp.MapQuotationRoute(apiRouter, quotationHandler, mw)
	})

	// New module routes that use /api/xxx prefix (need root router)
	dashboardHttp.MapDashboardRoute(r, dashHandler, mw)
	documentHttp.MapDocumentRoute(r, docHandler, mw)
	emailHttp.MapEmailRoute(r, emailHandler, mw)
	batchHttp.MapBatchRoute(r, batchHandler, mw)
	i18nHttp.MapI18nRoute(r, i18nHandler, mw)
	wosHttp.MapWosRoute(r, wosHandler, mw)
	reportHttp.MapReportRoute(r, cfg, reportHandler, mw)
	log.Info("✅ All module routes registered")

	return r
}
