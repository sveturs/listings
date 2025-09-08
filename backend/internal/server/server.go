// Package server
// backend/internal/server/server.go
package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	globalService "backend/internal/proj/global/service"
	postexpressService "backend/internal/proj/postexpress/service"
	postexpressRepository "backend/internal/proj/postexpress/storage/postgres"
	pkglogger "backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	pkgErrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	_ "backend/docs"
	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/middleware"
	adminLogistics "backend/internal/proj/admin/logistics"
	"backend/internal/proj/analytics"
	balanceHandler "backend/internal/proj/balance/handler"
	"backend/internal/proj/behavior_tracking"
	"backend/internal/proj/bexexpress"
	configHandler "backend/internal/proj/config"
	contactsHandler "backend/internal/proj/contacts/handler"
	docsHandler "backend/internal/proj/docserver/handler"
	geocodeHandler "backend/internal/proj/geocode/handler"
	gisHandler "backend/internal/proj/gis/handler"
	globalHandler "backend/internal/proj/global/handler"
	healthHandler "backend/internal/proj/health"
	marketplaceHandler "backend/internal/proj/marketplace/handler"
	marketplaceService "backend/internal/proj/marketplace/service"
	notificationHandler "backend/internal/proj/notifications/handler"
	"backend/internal/proj/orders"
	paymentHandler "backend/internal/proj/payments/handler"
	postexpressHandler "backend/internal/proj/postexpress/handler"
	reviewHandler "backend/internal/proj/reviews/handler"
	"backend/internal/proj/search_admin"
	"backend/internal/proj/search_optimization"
	"backend/internal/proj/storefronts"
	"backend/internal/proj/subscriptions"
	"backend/internal/proj/translation_admin"
	userHandler "backend/internal/proj/users/handler"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
)

type Server struct {
	app                *fiber.App
	cfg                *config.Config
	configModule       *configHandler.Module
	users              *userHandler.Handler
	middleware         *middleware.Middleware
	review             *reviewHandler.Handler
	marketplace        *marketplaceHandler.Handler
	notifications      *notificationHandler.Handler
	balance            *balanceHandler.Handler
	payments           *paymentHandler.Handler
	postexpress        *postexpressHandler.Handler
	bexexpress         *bexexpress.Module
	adminLogistics     *adminLogistics.Module
	orders             *orders.Module
	storefront         *storefronts.Module
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

	osClient, err := initializeOpenSearch(cfg)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "OpenSearch initialization failed")
	} else {
		logger.Info().Msg("Успешное подключение к OpenSearch")
	}
	db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage, cfg.SearchWeights)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to initialize database")
	}
	if cfg.ReindexOnAPI != "" {
		if err = db.ReindexAllListings(ctx); err != nil {
			logger.Error().Err(err).Msg("reindexAllListings() failded")
		}
	}

	translationService, err := initializeTranslationService(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize translation service: %w", err)
	}

	services := globalService.NewService(ctx, db, cfg, translationService)

	configModule := configHandler.NewModule(cfg)
	usersHandler := userHandler.NewHandler(services)
	reviewHandler := reviewHandler.NewHandler(services)
	notificationsHandler := notificationHandler.NewHandler(services.Notification())
	marketplaceHandlerInstance := marketplaceHandler.NewHandler(services)
	balanceHandler := balanceHandler.NewHandler(services)
	storefrontModule := storefronts.NewModule(services)
	ordersModule, err := orders.NewModule(db, &opensearch.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	})
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to initialize orders module")
	}
	contactsHandler := contactsHandler.NewHandler(services)
	paymentsHandler := paymentHandler.NewHandler(services)

	// Post Express инициализация
	postexpressRepo := postexpressRepository.NewRepository(db.GetSQLXDB())
	postexpressWSPClient := postexpressService.NewWSPClient(&postexpressService.WSPConfig{
		Endpoint:        cfg.PostExpress.BaseURL,
		Username:        cfg.PostExpress.Username,
		Password:        cfg.PostExpress.Password,
		TestMode:        cfg.PostExpress.TestMode,
		Timeout:         30 * time.Second,
		Language:        "sr",
		DeviceType:      2, // 2 для API интеграции согласно документации Post Express
		MaxRetries:      3,
		RetryDelay:      1 * time.Second,
		DeviceName:      "SveTu-Server",
		ApplicationName: "SveTu-Platform",
		Version:         "1.0.0",
	}, *pkglogger.New())
	postexpressServiceInstance := postexpressService.NewService(
		postexpressRepo,
		postexpressWSPClient,
		*pkglogger.New(),
		&postexpressService.ServiceConfig{
			DefaultWarehouseCode: "MAIN",
			PickupExpiryDays:     7,
			EnableAutoTracking:   true,
			TrackingInterval:     15 * time.Minute,
			MaxRetries:           3,
		},
	)
	postexpressHandlerInstance := postexpressHandler.NewHandler(postexpressServiceInstance, *pkglogger.New())

	// BEX Express инициализация
	bexexpressModule, err := bexexpress.NewModule(db.GetSQLXDB().DB, cfg)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize BEX Express module, continuing without it")
		// Не возвращаем ошибку, продолжаем без BEX
	}

	// Admin Logistics инициализация
	adminLogisticsModule, err := adminLogistics.NewModule(db.GetSQLXDB().DB, cfg, pkglogger.New())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize Admin Logistics module, continuing without it")
		// Не возвращаем ошибку, продолжаем без админки логистики
	}

	docsHandlerInstance := docsHandler.NewHandler(cfg.Docs)

	// Health handler
	// Получаем sql.DB из Database структуры
	var sqlDB *sql.DB
	if db != nil {
		sqlDB = db.GetSQLDB()
	}
	healthHandlerInstance := healthHandler.NewHandler(sqlDB, redisClient)
	middleware := middleware.NewMiddleware(cfg, services)
	geocodeHandler := geocodeHandler.NewHandler(services)
	globalHandlerInstance := globalHandler.NewHandler(services, cfg.SearchWeights)
	analyticsModule := analytics.NewModule(db, osClient)
	behaviorTrackingModule := behavior_tracking.NewModule(ctx, db.GetPool())
	translationAdminModule := translation_admin.NewModule(ctx, db.GetSQLXDB(), *logger.Get(), "/data/hostel-booking-system", redisClient, translationService)
	searchAdminModule := search_admin.NewModule(db, osClient, pkglogger.New())
	// TODO: После рефакторинга передать storage или services для переиндексации
	searchOptimizationModule := search_optimization.NewModule(db, *pkglogger.New())
	gisHandlerInstance := gisHandler.NewHandler(db.GetSQLXDB())

	// Инициализация модуля подписок
	// Используем nil для AllSecure - модуль автоматически переключится на mock payments
	subscriptionsModule := subscriptions.NewModule(db.GetSQLXDB(), nil, nil, pkglogger.New())

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
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "::1"},
	})

	server := &Server{
		app:                app,
		cfg:                cfg,
		configModule:       configModule,
		users:              usersHandler,
		middleware:         middleware,
		review:             reviewHandler,
		marketplace:        marketplaceHandlerInstance,
		notifications:      notificationsHandler,
		balance:            balanceHandler,
		payments:           paymentsHandler,
		postexpress:        postexpressHandlerInstance,
		bexexpress:         bexexpressModule,
		adminLogistics:     adminLogisticsModule,
		orders:             ordersModule,
		storefront:         storefrontModule,
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
		fileStorage:        fileStorage,
		health:             healthHandlerInstance,
		redisClient:        redisClient,
	}

	notificationsHandler.ConnectTelegramWebhook()
	server.setupMiddleware() //nolint:contextcheck
	server.setupRoutes()     //nolint:contextcheck

	return server, nil
}

func initializeTranslationService(cfg *config.Config, db *postgres.Database) (marketplaceService.TranslationServiceInterface, error) {
	// Используем новую фабрику V2 с поддержкой 4 провайдеров
	factoryConfig := struct {
		GoogleAPIKey    string
		OpenAIAPIKey    string
		ClaudeAPIKey    string
		DeepLAPIKey     string
		DeepLUseFreeAPI bool
	}{
		GoogleAPIKey:    cfg.GoogleTranslateAPIKey,
		OpenAIAPIKey:    cfg.OpenAIAPIKey,
		ClaudeAPIKey:    cfg.ClaudeAPIKey,
		DeepLAPIKey:     cfg.DeepLAPIKey,
		DeepLUseFreeAPI: cfg.DeepLUseFreeAPI,
	}

	translationFactory, err := marketplaceService.NewTranslationServiceFactoryV2(factoryConfig, db)
	if err == nil {
		availableProviders := translationFactory.GetAvailableProviders()
		logger.Info().
			Interface("providers", availableProviders).
			Int("count", len(availableProviders)).
			Msg("Создана фабрика сервисов перевода V2")
		return translationFactory, nil
	}

	// Fallback на старую версию если V2 не работает
	if cfg.GoogleTranslateAPIKey != "" && cfg.OpenAIAPIKey != "" {
		translationFactory, err := marketplaceService.NewTranslationServiceFactory(cfg.GoogleTranslateAPIKey, cfg.OpenAIAPIKey, db)
		if err == nil {
			logger.Info().Msg("Создана фабрика сервисов перевода (старая версия)")
			return translationFactory, nil
		}
	}

	// Крайний fallback на простой OpenAI сервис
	if cfg.OpenAIAPIKey != "" {
		translationService, err := marketplaceService.NewTranslationService(cfg.OpenAIAPIKey)
		if err != nil {
			return nil, err
		}
		logger.Info().Msg("Создан сервис перевода на базе OpenAI (fallback)")
		return translationService, nil
	}

	return nil, fmt.Errorf("не указан ни один API ключ для перевода")
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

	// Auth Proxy middleware - должен быть рано для перехвата /api/v1/auth/* запросов
	authProxy := middleware.NewAuthProxyMiddleware()
	s.app.Use(authProxy.ProxyToAuthService())

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
		return c.SendString("Svetu API")
	})

	// Health checks и metrics
	s.health.RegisterRoutes(s.app)

	// Swagger документация
	s.app.Get("/swagger/*", swagger.HandlerDefault)
	s.app.Get("/docs/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// WebSocket с проверкой аутентификации и rate limiting
	s.app.Get("/ws/chat", s.middleware.AuthRequiredJWT, s.middleware.RateLimitByUser(30, time.Minute), func(c *fiber.Ctx) error {
		// Проверяем, что это WebSocket запрос
		if websocket.IsWebSocketUpgrade(c) {
			// Сохраняем userID для использования в WebSocket handler
			userID := c.Locals("user_id").(int)

			return websocket.New(func(conn *websocket.Conn) {
				// Передаем userID через контекст соединения
				// В Fiber WebSocket, Locals доступен только для чтения
				// Поэтому создаем обертку с сохраненным userID
				s.marketplace.Chat.HandleWebSocketWithAuth(conn, userID)
			}, websocket.Config{
				HandshakeTimeout:  10 * time.Second,
				ReadBufferSize:    1024,
				WriteBufferSize:   1024,
				EnableCompression: false,
			})(c)
		}
		return fiber.ErrUpgradeRequired
	})

	// CSRF токен - регистрируем ДО проектных роутов чтобы избежать конфликта с AuthRequiredJWT
	s.app.Get("/api/v1/csrf-token", s.middleware.GetCSRFToken())

	// Config роуты - публичные, регистрируем ДО проектных роутов чтобы избежать конфликта с AuthRequiredJWT
	s.app.Get("/api/v1/config/storage", s.configModule.Handler.GetStorageConfig)

	// Регистрируем роуты через новую систему
	s.registerProjectRoutes()

	// Проксирование статических файлов MinIO
	// Эти маршруты должны быть после всех API маршрутов
	s.app.Get("/listings/*", s.ProxyMinIO)
	s.app.Get("/chat-files/*", s.ProxyChatFiles)
	s.app.Get("/storefront-products/*", s.ProxyStorefrontProducts)
}

// registerProjectRoutes регистрирует роуты проектов через новую систему
func (s *Server) registerProjectRoutes() {
	// Создаем слайс всех проектов, которые реализуют RouteRegistrar
	var registrars []RouteRegistrar

	// Добавляем все проекты, которые реализуют RouteRegistrar
	// ВАЖНО: global должен быть первым, чтобы его публичные API не конфликтовали с авторизацией других модулей
	// config регистрируется отдельно до этого метода для публичных роутов
	// searchOptimization должен быть раньше marketplace, чтобы избежать конфликта с глобальным middleware
	// subscriptions должен быть раньше marketplace, чтобы публичные роуты не перехватывались auth middleware
	registrars = append(registrars, s.global, s.notifications, s.users, s.review, s.searchOptimization, s.searchAdmin)

	// Добавляем Subscriptions если он инициализирован - ДО marketplace чтобы избежать конфликтов с auth middleware
	if s.subscriptions != nil {
		registrars = append(registrars, s.subscriptions)
	}

	registrars = append(registrars, s.marketplace, s.balance, s.orders, s.storefront,
		s.geocode, s.gis, s.contacts, s.payments, s.postexpress)

	// Добавляем BEX Express если он инициализирован
	if s.bexexpress != nil {
		registrars = append(registrars, s.bexexpress)
	}

	// Добавляем Admin Logistics если он инициализирован
	if s.adminLogistics != nil {
		registrars = append(registrars, s.adminLogistics)
	}

	registrars = append(registrars, s.docs, s.analytics, s.behaviorTracking, s.translationAdmin)

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
