package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	authservice "github.com/sveturs/auth/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	attributespb "github.com/sveturs/listings/api/proto/attributes/v1"
	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/cache"
	"github.com/sveturs/listings/internal/config"
	"github.com/sveturs/listings/internal/health"
	"github.com/sveturs/listings/internal/metrics"
	"github.com/sveturs/listings/internal/middleware"
	"github.com/sveturs/listings/internal/ratelimit"
	"github.com/sveturs/listings/internal/repository/minio"
	"github.com/sveturs/listings/internal/repository/opensearch"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service"
	"github.com/sveturs/listings/internal/service/listings"
	"github.com/sveturs/listings/internal/timeout"
	grpcTransport "github.com/sveturs/listings/internal/transport/grpc"
	httpTransport "github.com/sveturs/listings/internal/transport/http"
	"github.com/sveturs/listings/internal/worker"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
)

func main() {
	// Handle CLI commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("Listings Service %s (built: %s)\n", Version, BuildTime)
			return
		case "healthcheck":
			fmt.Println("OK")
			return
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize zerolog logger (for HTTP/gRPC transports)
	zerologLogger := initLogger(cfg.App.LogLevel, cfg.App.LogFormat)
	zerologLogger.Info().
		Str("version", Version).
		Str("env", cfg.App.Env).
		Msg("Starting Listings Service")

	// Initialize zerolog logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := zerolog.InfoLevel
	if cfg.App.LogLevel == "debug" {
		logLevel = zerolog.DebugLevel
	}
	logger := zerolog.New(os.Stdout).Level(logLevel).With().
		Timestamp().
		Str("service", "listings").
		Str("version", Version).
		Logger()

	// Initialize metrics
	metricsInstance := metrics.NewMetrics("listings")

	// Start pprof server on separate port (for profiling)
	go func() {
		pprofAddr := ":6060"
		logger.Info().Str("addr", pprofAddr).Msg("Starting pprof server")
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			logger.Error().Err(err).Msg("pprof server failed")
		}
	}()

	// Initialize PostgreSQL
	db, err := postgres.InitDB(
		cfg.DB.DSN(),
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.ConnMaxLifetime,
		cfg.DB.ConnMaxIdleTime,
		zerologLogger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	pgRepo := postgres.NewRepository(db, zerologLogger)

	// Initialize PgxPool for order-related repositories (they require pgx for transactions)
	pgxPool, err := postgres.InitPgxPool(
		context.Background(),
		cfg.DB.DSN(),
		int32(cfg.DB.MaxOpenConns),
		int32(cfg.DB.MaxIdleConns),
		zerologLogger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize pgxpool")
	}
	defer pgxPool.Close()

	// Start DB stats collector
	dbStatsCollector := metrics.NewDBStatsCollector(db, metricsInstance, zerologLogger, 15*time.Second)
	go dbStatsCollector.Start(context.Background())
	defer dbStatsCollector.Stop()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(
		cfg.Redis.Addr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PoolSize,
		cfg.Redis.MinIdleConns,
		cfg.Redis.ListingTTL,
		zerologLogger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize Redis cache")
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close Redis cache")
		}
	}()

	// Initialize OpenSearch (if enabled)
	var searchClient *opensearch.Client
	if cfg.Features.AsyncIndexing {
		searchClient, err = opensearch.NewClient(
			cfg.Search.Addresses,
			cfg.Search.Username,
			cfg.Search.Password,
			cfg.Search.Index,
			zerologLogger,
		)
		if err != nil {
			logger.Warn().Err(err).Msg("OpenSearch not available, continuing without search")
		} else {
			defer func() {
				if err := searchClient.Close(); err != nil {
					logger.Error().Err(err).Msg("failed to close OpenSearch client")
				}
			}()
		}
	}

	// Initialize MinIO (optional)
	var minioClient *minio.Client
	if cfg.Storage.Endpoint != "" {
		minioClient, err = minio.NewClient(
			cfg.Storage.Endpoint,
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			cfg.Storage.Bucket,
			cfg.Storage.UseSSL,
			zerologLogger,
		)
		if err != nil {
			logger.Warn().Err(err).Msg("MinIO not available, continuing without object storage")
		}
	}

	// Initialize Auth Service client (if enabled)
	var authInterceptor *middleware.AuthInterceptor
	if cfg.Auth.Enabled {
		authServiceConfig := &authservice.Config{
			HTTPURL: cfg.Auth.ServiceURL,
			Timeout: cfg.Auth.Timeout,
		}

		authSvc, err := authservice.NewAuthService(authServiceConfig)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to create auth service")
		}

		logger.Info().
			Str("url", cfg.Auth.ServiceURL).
			Msg("Auth service initialized")

		authInterceptor = middleware.NewAuthInterceptor(authSvc, zerologLogger)
	} else {
		logger.Warn().Msg("Auth service DISABLED - all requests will be unauthenticated")
	}

	// Initialize listings service
	listingsService := listings.NewService(pgRepo, redisCache, searchClient, zerologLogger)

	// Initialize storefront service
	storefrontService := listings.NewStorefrontService(pgRepo, &zerologLogger)

	// Initialize attribute repository and service
	attrRepo := postgres.NewAttributeRepository(db, zerologLogger)
	attributeService := service.NewAttributeService(attrRepo, redisCache.GetClient(), zerologLogger)

	// Initialize category service
	categoryService := service.NewCategoryService(pgRepo, redisCache.GetClient(), zerologLogger)

	// Initialize order service dependencies
	cartRepo := postgres.NewCartRepository(pgxPool, zerologLogger)
	orderRepo := postgres.NewOrderRepository(pgxPool, zerologLogger)
	reservationRepo := postgres.NewReservationRepository(pgxPool, zerologLogger)

	// Initialize cart service
	cartService := service.NewCartService(
		cartRepo,
		pgRepo,
		pgRepo,
		pgRepo,
		zerologLogger,
	)

	// Initialize order service
	orderService := service.NewOrderService(
		orderRepo,
		cartRepo,
		reservationRepo,
		pgRepo,
		pgxPool,
		nil, // Use default financial config
		zerologLogger,
	)

	// Initialize health check service
	healthConfig := &health.Config{
		CheckTimeout:     cfg.Health.CheckTimeout,
		CheckInterval:    cfg.Health.CheckInterval,
		StartupTimeout:   cfg.Health.StartupTimeout,
		CacheDuration:    cfg.Health.CacheDuration,
		EnableDeepChecks: cfg.Health.EnableDeepChecks,
	}
	healthChecker := health.NewService(db.DB, redisCache, searchClient, minioClient, healthConfig, zerologLogger)

	// Initialize worker (if enabled and search is available)
	var indexWorker *worker.Worker
	if cfg.Worker.Enabled && searchClient != nil {
		indexWorker = worker.NewWorker(
			pgRepo,
			searchClient,
			metricsInstance,
			cfg.Worker.Concurrency,
			zerologLogger,
		)
		if err := indexWorker.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start indexing worker")
		}
		defer func() {
			if err := indexWorker.Stop(); err != nil {
				logger.Error().Err(err).Msg("failed to stop indexing worker")
			}
		}()
	}

	// Initialize rate limiter (conditionally based on config)
	var rateLimiterInterceptor grpc.UnaryServerInterceptor
	if cfg.Features.RateLimitEnabled {
		rateLimiterConfig := ratelimit.NewDefaultConfig()
		rateLimiterInstance := ratelimit.NewRedisLimiter(redisCache.GetClient(), zerologLogger)
		logger.Info().Msg("Rate limiter initialized with Redis backend")

		// Create rate limiter interceptor with metrics
		rateLimiterInterceptor = ratelimit.UnaryServerInterceptorWithMetrics(
			rateLimiterInstance,
			rateLimiterConfig,
			metricsInstance,
			zerologLogger,
		)
	} else {
		logger.Warn().Msg("Rate limiting DISABLED - not recommended for production")
	}

	// Create timeout interceptor
	timeoutInterceptor := timeout.UnaryServerInterceptor(metricsInstance, zerologLogger)

	// Build interceptor chain (conditionally include rate limiter and auth)
	// Order: timeout (outermost) → auth → rate limiting → metrics (innermost)
	interceptors := []grpc.UnaryServerInterceptor{
		timeoutInterceptor,
	}
	if cfg.Auth.Enabled && authInterceptor != nil {
		interceptors = append(interceptors, authInterceptor.Unary())
	}
	if cfg.Features.RateLimitEnabled {
		interceptors = append(interceptors, rateLimiterInterceptor)
	}
	interceptors = append(interceptors, metricsInstance.UnaryServerInterceptor())

	// Initialize gRPC server with interceptors
	// Order: timeout → auth (if enabled) → rate limiting (if enabled) → metrics
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	grpcHandler := grpcTransport.NewServer(listingsService, storefrontService, attributeService, categoryService, orderService, cartService, metricsInstance, zerologLogger)
	listingspb.RegisterListingsServiceServer(grpcServer, grpcHandler)
	attributespb.RegisterAttributeServiceServer(grpcServer, grpcHandler)
	listingspb.RegisterOrderServiceServer(grpcServer, grpcHandler)

	// Enable gRPC reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start gRPC server in goroutine
	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.GRPCHost, cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create gRPC listener")
	}

	go func() {
		logger.Info().Int("port", cfg.Server.GRPCPort).Msg("Starting gRPC server")
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Error().Err(err).Msg("gRPC server error")
		}
	}()

	// Initialize HTTP server with health checks
	httpHandler := httpTransport.NewMinimalHandler(listingsService, zerologLogger)
	healthHandler := httpTransport.NewHealthHandler(healthChecker, zerologLogger)

	httpApp, err := httpTransport.StartMinimalServer(
		cfg.Server.HTTPHost,
		cfg.Server.HTTPPort,
		httpHandler,
		healthHandler,
		zerologLogger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start HTTP server")
	}

	defer func() {
		if err := httpApp.Shutdown(); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown HTTP server")
		}
	}()

	logger.Info().
		Int("http_port", cfg.Server.HTTPPort).
		Int("grpc_port", cfg.Server.GRPCPort).
		Msg("Listings Service started successfully (HTTP + gRPC)")

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
	logger.Info().Msg("Shutting down gracefully...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop worker first
	if indexWorker != nil {
		if err := indexWorker.Stop(); err != nil {
			logger.Error().Err(err).Msg("error stopping worker")
		}
	}

	// Shutdown gRPC server
	logger.Info().Msg("Stopping gRPC server...")
	grpcServer.GracefulStop()

	// Shutdown HTTP server
	if err := httpApp.ShutdownWithContext(ctx); err != nil {
		logger.Error().Err(err).Msg("error shutting down HTTP server")
	}

	logger.Info().Msg("Listings Service stopped")
}

// initLogger initializes zerolog logger
func initLogger(level, format string) zerolog.Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// Configure output format
	if format == "pretty" || format == "console" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	// Default to JSON
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}
