// Package server
// backend/internal/server/server.go
package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	version "backend/internal/version"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	pkgErrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	authclient "github.com/sveturs/auth/pkg/http/client"
	authService "github.com/sveturs/auth/pkg/service"

	_ "backend/docs"
	"backend/internal/cache"
	"backend/internal/clients/listings"
	"backend/internal/config"
	"backend/internal/interfaces"
	"backend/internal/logger"
	"backend/internal/middleware"
	adminLogistics "backend/internal/proj/admin/logistics"
	testingHandler "backend/internal/proj/admin/testing/handler"
	testingService "backend/internal/proj/admin/testing/service"
	testingStorage "backend/internal/proj/admin/testing/storage/postgres"
	aiHandler "backend/internal/proj/ai/handler"
	"backend/internal/proj/analytics"
	balanceHandler "backend/internal/proj/balance/handler"
	"backend/internal/proj/behavior_tracking"
	chat "backend/internal/proj/chat"
	configHandler "backend/internal/proj/config"
	contactsHandler "backend/internal/proj/contacts/handler"
	creditHandler "backend/internal/proj/credit"
	"backend/internal/proj/delivery"
	delivery_grpcclient "backend/internal/proj/delivery/grpcclient"
	docsHandler "backend/internal/proj/docserver/handler"
	geocodeHandler "backend/internal/proj/geocode/handler"
	gisHandler "backend/internal/proj/gis/handler"
	globalHandler "backend/internal/proj/global/handler"
	globalService "backend/internal/proj/global/service"
	healthHandler "backend/internal/proj/health"
	marketplaceHandler "backend/internal/proj/marketplace/handler"
	notificationHandler "backend/internal/proj/notifications/handler"
	"backend/internal/proj/orders"
	paymentHandler "backend/internal/proj/payments/handler"
	recommendationsHandler "backend/internal/proj/recommendations"
	reviewHandler "backend/internal/proj/reviews/handler"
	"backend/internal/proj/search_admin"
	"backend/internal/proj/search_optimization"
	"backend/internal/proj/subscriptions"
	"backend/internal/proj/tracking"
	"backend/internal/proj/translation_admin"
	userHandler "backend/internal/proj/users/handler"
	"backend/internal/proj/viber"
	vinModule "backend/internal/proj/vin"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	pkglogger "backend/pkg/logger"
)

type Server struct {
	app                *fiber.App
	cfg                *config.Config
	configModule       *configHandler.Module
	ai                 *aiHandler.Handler
	users              *userHandler.Handler
	middleware         *middleware.Middleware
	authService        *authService.AuthService
	jwtParserMW        fiber.Handler
	review             *reviewHandler.Handler
	notifications      *notificationHandler.Handler
	balance            *balanceHandler.Handler
	payments           *paymentHandler.Handler
	adminLogistics     *adminLogistics.Module
	adminTesting       *testingHandler.Handler
	delivery           *delivery.Module
	orders             *orders.Module
	geocode            *geocodeHandler.Handler
	contacts           *contactsHandler.Handler
	docs               *docsHandler.Handler
	analytics          *analytics.Module
	behaviorTracking   *behavior_tracking.Module
	translationAdmin   *translation_admin.Module
	searchAdmin        *search_admin.Module
	searchOptimization *search_optimization.Module
	global             *globalHandler.Handler
	gis                *gisHandler.Handler
	subscriptions      *subscriptions.Module
	tracking           *tracking.Module
	viber              *viber.Module
	vin                *vinModule.Module
	chat               *chat.Module
	marketplace        *marketplaceHandler.Handler
	credit             *creditHandler.Handler
	recommendations    *recommendationsHandler.Handler
	fileStorage        filestorage.FileStorageInterface
	health             *healthHandler.Handler
	redisClient        *redis.Client
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	fileStorage, err := filestorage.NewFileStorage(ctx, cfg.FileStorage)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "Ошибка инициализации файлового хранилища")
	}

	// Инициализируем Redis для кеширования переводов
	var redisClient *redis.Client
	if cfg.Redis.URL != "" {
		logrusLogger := logrus.New()
		redisCache, err := cache.NewRedisCache(ctx, cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.PoolSize, logrusLogger)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to initialize Redis cache, continuing without cache")
		} else {
			redisClient = redisCache.GetClient()
			logger.Info().Msg("Redis cache initialized successfully")
		}
	}

	// Инициализируем универсальный кеш для маркетплейса (будет использован позже)
	// var universalCache *marketplaceCache.UniversalCache
	// if cfg.Redis.URL != "" {
	// 	// Будет интегрирован после полного тестирования
	// 	logger.Info().Msg("Universal marketplace cache will be integrated later")
	// }

	osClient, err := initializeOpenSearch(cfg)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "OpenSearch initialization failed")
	} else {
		logger.Info().Msg("Успешное подключение к OpenSearch")
	}
	db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, cfg.OpenSearch.UnifiedIndex, fileStorage, cfg.SearchWeights, cfg)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to initialize database")
	}
	if cfg.ReindexOnAPI != "" {
		if err = db.ReindexAllListings(ctx); err != nil {
			logger.Error().Err(err).Msg("reindexAllListings() failded")
		}
	}

	// Create auth service client BEFORE creating services
	authClient, err := authclient.NewClientWithResponses(cfg.AuthServiceURL)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to create auth client")
	}

	// Create auth service with local JWT validation (automatically fetches public key from auth service)
	zerologLogger := logger.Get()
	authServiceInstance := authService.NewAuthServiceWithLocalValidation(authClient, zerologLogger)
	logger.Info().
		Str("auth_service_url", cfg.AuthServiceURL).
		Msg("Auth service initialized with local JWT validation (public key will be fetched from auth service)")

	userServiceInstance := authService.NewUserService(authClient, zerologLogger)
	oauthServiceInstance := authService.NewOAuthService(authClient)

	// Now create services with authService and userService
	services := globalService.NewService(ctx, db, cfg, authServiceInstance, userServiceInstance)

	configModule := configHandler.NewModule(cfg)
	aiHandlerInstance := aiHandler.NewHandler(cfg, services)

	authHandler := userHandler.NewAuthHandler(authServiceInstance, oauthServiceInstance, cfg.BackendURL, cfg.FrontendURL, *zerologLogger)
	jwtParserMW := authMiddleware.JWTParser(authServiceInstance)
	usersHandler := userHandler.NewHandler(services, authHandler, jwtParserMW)

	reviewHandler := reviewHandler.NewHandler(services, jwtParserMW)
	notificationsHandler := notificationHandler.NewHandler(services.Notification(), jwtParserMW)
	balanceHandler := balanceHandler.NewHandler(services, jwtParserMW)

	// Delivery система инициализация с консолидацией admin/logistics
	// ВАЖНО: Должна быть инициализирована ДО orders модуля (для передачи deliveryClient)
	deliveryModule, err := delivery.NewModule(db.GetSQLXDB(), cfg, pkglogger.New())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize Delivery module, continuing without it")
		// Не возвращаем ошибку, продолжаем без delivery системы
	} else if deliveryModule != nil && services != nil {
		// Подключаем сервис уведомлений к модулю доставки
		deliveryModule.SetNotificationService(services.Notification())
		logger.Info().Msg("Notification service integrated with delivery module")
	}

	contactsHandler := contactsHandler.NewHandler(services)
	paymentsHandler := paymentHandler.NewHandler(services, jwtParserMW)

	// Admin Logistics инициализация
	adminLogisticsModule, err := adminLogistics.NewModule(db.GetSQLXDB().DB, cfg, pkglogger.New())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize Admin Logistics module, continuing without it")
		// Не возвращаем ошибку, продолжаем без админки логистики
	}

	// Admin Testing module инициализация
	var adminTestingHandler *testingHandler.Handler
	testAdminEmail := os.Getenv("TEST_ADMIN_EMAIL")
	testAdminPassword := os.Getenv("TEST_ADMIN_PASSWORD")
	if testAdminEmail != "" && testAdminPassword != "" {
		testStorage := testingStorage.NewStorage(db.GetSQLXDB(), *zerologLogger)
		testAuthMgr := testingService.NewTestAuthManager(cfg.BackendURL, testAdminEmail, testAdminPassword, *zerologLogger)
		testRunner := testingService.NewTestRunner(testStorage, testAuthMgr, cfg.BackendURL, *zerologLogger)
		adminTestingHandler = testingHandler.NewHandler(testRunner, jwtParserMW, *zerologLogger)
		logger.Info().
			Str("admin_email", testAdminEmail).
			Str("backend_url", cfg.BackendURL).
			Msg("Admin Testing module initialized")
	} else {
		logger.Info().Msg("Admin Testing module disabled (no TEST_ADMIN_EMAIL or TEST_ADMIN_PASSWORD)")
	}

	docsHandlerInstance := docsHandler.NewHandler(cfg.Docs, jwtParserMW)

	// Health handler
	// Получаем sql.DB из Database структуры
	var sqlDB *sql.DB
	if db != nil {
		sqlDB = db.GetSQLDB()
	}
	healthHandlerInstance := healthHandler.NewHandler(sqlDB, redisClient)
	middleware := middleware.NewMiddleware(cfg, services, authServiceInstance, jwtParserMW)
	geocodeHandler := geocodeHandler.NewHandler(services)
	globalHandlerInstance := globalHandler.NewHandler(services, cfg.SearchWeights)
	analyticsModule := analytics.NewModule(db, osClient, jwtParserMW)
	behaviorTrackingModule := behavior_tracking.NewModule(ctx, db.GetPool(), jwtParserMW)
	// Translation service moved to listings microservice - pass nil for now
	translationAdminModule := translation_admin.NewModule(ctx, db.GetSQLXDB(), *logger.Get(), "/data/hostel-booking-system", redisClient, nil, jwtParserMW)
	searchAdminModule := search_admin.NewModule(db, osClient, pkglogger.New())
	// TODO: После рефакторинга передать storage или services для переиндексации
	searchOptimizationModule := search_optimization.NewModule(db, *pkglogger.New())
	gisHandlerInstance := gisHandler.NewHandler(db.GetSQLXDB())

	// Инициализация модуля подписок
	// Используем nil для AllSecure - модуль автоматически переключится на mock payments
	subscriptionsModule := subscriptions.NewModule(db.GetSQLXDB(), nil, nil, pkglogger.New(), jwtParserMW)

	// Инициализация модуля трекинга
	trackingModule := tracking.NewModule(db) //nolint:contextcheck

	// Инициализация модуля Viber
	viberModule := viber.NewModule(services)

	// Инициализация модуля VIN декодера
	vinModule := vinModule.NewModule(db.GetSQLXDB())

	// TEMPORARY: Инициализация модуля чата (будет перенесен в микросервис)
	// Включен с заглушкой - возвращает сообщение об отключении функционала
	chatModule := chat.New(services, cfg, jwtParserMW)

	// Listings gRPC client (Phase 7.4 - Categories Integration)
	var listingsClient *listings.Client
	if cfg.UseListingsMicroservice && cfg.ListingsGRPCURL != "" {
		var err error
		listingsClient, err = listings.NewClient(cfg.ListingsGRPCURL, *zerologLogger)
		if err != nil {
			logger.Error().Err(err).Str("url", cfg.ListingsGRPCURL).Msg("Failed to create listings gRPC client, falling back to monolith")
			listingsClient = nil // Fallback to monolith
		} else {
			logger.Info().Str("url", cfg.ListingsGRPCURL).Msg("Listings gRPC client initialized successfully")
			// Inject listings client into cart repository for product data
			db.SetListingsClientToCart(listingsClient)
		}
	} else {
		logger.Info().Bool("use_microservice", cfg.UseListingsMicroservice).Str("grpc_url", cfg.ListingsGRPCURL).Msg("Listings microservice disabled, using monolith")
	}

	// Orders модуль инициализация (ПОСЛЕ delivery и listings модулей)
	var deliveryClient *delivery_grpcclient.Client
	if deliveryModule != nil {
		deliveryClient = deliveryModule.GetGRPCClient()
	}
	// OpenSearch integration removed (c2c/b2c deprecated)
	ordersModule, err := orders.NewModule(db, deliveryClient, listingsClient, redisClient)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to initialize orders module")
	}

	// TEMPORARY: Marketplace handler (minimal functionality until microservice migration)
	marketplaceHandlerInstance := marketplaceHandler.NewHandler(db.GetSQLXDB(), services, jwtParserMW, *zerologLogger, listingsClient, cfg.UseListingsMicroservice)

	// Инициализация универсальных handlers
	creditHandlerInstance := creditHandler.NewHandler()
	recommendationsHandlerInstance := recommendationsHandler.NewHandler(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Детальное логирование ошибки
			logger.Error().
				Err(err).
				Str("error_type", fmt.Sprintf("%T", err)).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Str("stack", fmt.Sprintf("%+v", err)).
				Msg("Error in handler")

			// Стандартная обработка ошибки
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":  err.Error(),
				"status": code,
				"path":   c.Path(),
			})
		},
		BodyLimit:               50 * 1024 * 1024,
		ReadBufferSize:          16384, // Увеличиваем размер буфера чтения для больших заголовков
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "::1"},
	})

	server := &Server{
		app:                app,
		cfg:                cfg,
		configModule:       configModule,
		ai:                 aiHandlerInstance,
		users:              usersHandler,
		middleware:         middleware,
		authService:        authServiceInstance,
		jwtParserMW:        jwtParserMW,
		review:             reviewHandler,
		notifications:      notificationsHandler,
		balance:            balanceHandler,
		payments:           paymentsHandler,
		adminLogistics:     adminLogisticsModule,
		adminTesting:       adminTestingHandler,
		delivery:           deliveryModule,
		orders:             ordersModule,
		geocode:            geocodeHandler,
		contacts:           contactsHandler,
		docs:               docsHandlerInstance,
		analytics:          analyticsModule,
		behaviorTracking:   behaviorTrackingModule,
		translationAdmin:   translationAdminModule,
		searchAdmin:        searchAdminModule,
		searchOptimization: searchOptimizationModule,
		global:             globalHandlerInstance,
		gis:                gisHandlerInstance,
		subscriptions:      subscriptionsModule,
		tracking:           trackingModule,
		viber:              viberModule,
		vin:                vinModule,
		chat:               chatModule, // Enabled with stub - returns disabled message
		marketplace:        marketplaceHandlerInstance,
		credit:             creditHandlerInstance,
		recommendations:    recommendationsHandlerInstance,
		fileStorage:        fileStorage,
		health:             healthHandlerInstance,
		redisClient:        redisClient,
	}

	notificationsHandler.ConnectTelegramWebhook()

	server.setupMiddleware() //nolint:contextcheck
	server.setupRoutes()     //nolint:contextcheck

	return server, nil
}

func initializeOpenSearch(cfg *config.Config) (*opensearch.OpenSearchClient, error) {
	if cfg.OpenSearch.URL == "" {
		return nil, errors.New("OpenSearch URL не указан, поиск будет отключен")
	}

	osClient, err := opensearch.NewOpenSearchClient(opensearch.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	})
	if err != nil {
		return nil, errors.New("ошибка подключения к OpenSearch")
	}

	return osClient, nil
}

func (s *Server) setupMiddleware() {
	// Общие middleware для observability
	// Security headers должны быть первыми
	s.app.Use(s.middleware.SecurityHeaders())

	// CORS должен быть первым для обработки preflight запросов
	s.app.Use(s.middleware.CORS())
	s.app.Use(s.middleware.Logger())

	// Middleware для определения языка из запроса
	s.app.Use(s.middleware.LocaleMiddleware())

	// Prometheus metrics middleware
	s.app.Use(middleware.PrometheusMiddleware())

	// Инициализируем метрики feature flags
	if s.cfg.FeatureFlags != nil {
		middleware.InitializeFeatureFlagMetrics(map[string]bool{
			"USE_UNIFIED_ATTRIBUTES":      s.cfg.FeatureFlags.UseUnifiedAttributes,
			"UNIFIED_ATTRIBUTES_FALLBACK": s.cfg.FeatureFlags.UnifiedAttributesFallback,
			"DUAL_WRITE_ATTRIBUTES":       s.cfg.FeatureFlags.DualWriteAttributes,
		})
	}

	// Удалено глобальное применение AuthRequiredJWT
}

func (s *Server) setupRoutes() { //nolint:contextcheck // внутренние маршруты
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Svetu API " + version.GetVersion())
	})

	// Health checks и metrics
	s.health.RegisterRoutes(s.app)

	// Prometheus metrics endpoint
	s.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Swagger документация
	s.app.Get("/swagger/*", swagger.HandlerDefault)
	s.app.Get("/docs/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// TEMPORARY: WebSocket /ws/chat (will be moved to chat microservice)
	// Регистрируем WebSocket роут для чата с JWT аутентификацией
	if s.chat != nil {
		_ = s.chat.RegisterRoutes(s.app)
	}

	// WebSocket для трекинга доставок (публичный, по токену)
	s.app.Get("/ws/tracking/:token", func(c *fiber.Ctx) error {
		token := c.Params("token")
		if token == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Missing tracking token")
		}

		// Проверяем, что это WebSocket запрос
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(conn *websocket.Conn) {
				// Проверяем токен и получаем delivery
				if s.tracking != nil && s.tracking.DeliveryService != nil {
					delivery, err := s.tracking.DeliveryService.ValidateTrackingToken(token)
					if err != nil {
						// Отправляем ошибку и закрываем соединение
						_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"Invalid tracking token"}`))
						_ = conn.Close()
						return
					}

					// Используем Hub для обработки WebSocket
					if s.tracking.Hub != nil {
						s.tracking.Hub.HandleWebSocket(conn, delivery.ID)
					} else {
						// Fallback если Hub не инициализирован
						_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"connected","delivery_id":`+strconv.Itoa(delivery.ID)+`}`))
						for {
							_, _, err := conn.ReadMessage()
							if err != nil {
								break
							}
						}
					}
				} else {
					// Если tracking module не инициализирован
					_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"Tracking service unavailable"}`))
					_ = conn.Close()
				}
			}, websocket.Config{
				HandshakeTimeout:  10 * time.Second,
				ReadBufferSize:    1024,
				WriteBufferSize:   1024,
				EnableCompression: false,
			})(c)
		}
		return fiber.ErrUpgradeRequired
	})

	// Config роуты - публичные, регистрируем ДО проектных роутов чтобы избежать конфликта с AuthRequiredJWT
	s.app.Get("/api/v1/config/storage", s.configModule.Handler.GetStorageConfig)

	// Регистрируем роуты через новую систему
	s.registerProjectRoutes()

	// Test route - direct registration without any middleware AFTER project routes
	s.app.Post("/api/v1/auth/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Test route works, registered AFTER project routes"})
	})

	// Проксирование статических файлов MinIO
	// Эти маршруты должны быть после всех API маршрутов
	s.app.Get("/listings/*", s.ProxyMinIO)
	s.app.Get("/chat-files/*", s.ProxyChatFiles)
	s.app.Get("/storefront-products/*", s.ProxyStorefrontProducts)
}

// registerProjectRoutes регистрирует роуты проектов через новую систему
func (s *Server) registerProjectRoutes() {
	// Создаем слайс всех проектов, которые реализуют RouteRegistrar
	var registrars []interfaces.RouteRegistrar

	// Добавляем все проекты, которые реализуют RouteRegistrar
	// ВАЖНО: global должен быть первым, чтобы его публичные API не конфликтовали с авторизацией других модулей
	// config регистрируется отдельно до этого метода для публичных роутов
	// analytics должен быть РАНЬШЕ остальных, чтобы публичный /api/v1/analytics/event зарегистрировался первым
	// searchOptimization должен быть раньше, чтобы избежать конфликта с глобальным middleware
	// subscriptions должен быть раньше, чтобы публичные роуты не перехватывались auth middleware
	// tracking должен быть раньше, чтобы его публичные роуты не перехватывались auth middleware
	// TEMPORARY: marketplace должен быть раньше, чтобы публичные категории не требовали auth
	if s.marketplace != nil {
		registrars = append(registrars, s.marketplace)
	}
	registrars = append(registrars, s.global, s.analytics, s.ai, s.notifications, s.users, s.review, s.searchOptimization, s.searchAdmin, s.tracking)

	// Добавляем Subscriptions если он инициализирован - ДО marketplace чтобы избежать конфликтов с auth middleware
	if s.subscriptions != nil {
		registrars = append(registrars, s.subscriptions)
	}

	registrars = append(registrars, s.balance, s.orders,
		s.geocode, s.gis, s.contacts, s.payments)

	// Добавляем Delivery если он инициализирован
	if s.delivery != nil {
		registrars = append(registrars, s.delivery)
	}

	// Добавляем Admin Logistics если он инициализирован
	if s.adminLogistics != nil {
		registrars = append(registrars, s.adminLogistics)
	}

	// Добавляем Admin Testing если он инициализирован
	if s.adminTesting != nil {
		// Register routes directly since it doesn't implement RouteRegistrar
		s.adminTesting.RegisterRoutes(s.app)
		logger.Info().Str("prefix", s.adminTesting.GetPrefix()).Msg("Admin Testing routes registered")
	}

	registrars = append(registrars, s.docs, s.behaviorTracking, s.translationAdmin, s.viber)

	// Добавляем VIN модуль
	if s.vin != nil {
		registrars = append(registrars, s.vin)
	}

	// TEMPORARY: Marketplace уже зарегистрирован в начале для избежания middleware конфликтов

	// Добавляем универсальные handlers
	if s.credit != nil {
		registrars = append(registrars, s.credit)
	}
	if s.recommendations != nil {
		registrars = append(registrars, s.recommendations)
	}

	// Регистрируем роуты каждого проекта
	for _, registrar := range registrars {
		err := registrar.RegisterRoutes(s.app, s.middleware)
		if err != nil {
			logger.Error().
				Err(err).
				Str("prefix", registrar.GetPrefix()).
				Msg("Ошибка регистрации роутов проекта")
		} else {
			logger.Info().
				Str("prefix", registrar.GetPrefix()).
				Msg("Роуты проекта успешно зарегистрированы")
		}
	}
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
