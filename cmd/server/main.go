package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/cache"
	"github.com/sveturs/listings/internal/config"
	"github.com/sveturs/listings/internal/metrics"
	"github.com/sveturs/listings/internal/repository/opensearch"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service/listings"
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

	// Initialize logger
	logger := initLogger(cfg.App.LogLevel, cfg.App.LogFormat)
	logger.Info().
		Str("version", Version).
		Str("env", cfg.App.Env).
		Msg("Starting Listings Service")

	// Initialize metrics
	metricsInstance := metrics.NewMetrics("listings")

	// Initialize PostgreSQL
	db, err := postgres.InitDB(
		cfg.DB.DSN(),
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.ConnMaxLifetime,
		cfg.DB.ConnMaxIdleTime,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer db.Close()

	pgRepo := postgres.NewRepository(db, logger)

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(
		cfg.Redis.Addr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PoolSize,
		cfg.Redis.MinIdleConns,
		cfg.Redis.ListingTTL,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize Redis cache")
	}
	defer redisCache.Close()

	// Initialize OpenSearch (if enabled)
	var searchClient *opensearch.Client
	if cfg.Features.AsyncIndexing {
		searchClient, err = opensearch.NewClient(
			cfg.Search.Addresses,
			cfg.Search.Username,
			cfg.Search.Password,
			cfg.Search.Index,
			logger,
		)
		if err != nil {
			logger.Warn().Err(err).Msg("OpenSearch not available, continuing without search")
		} else {
			defer searchClient.Close()
		}
	}

	// Initialize listings service
	listingsService := listings.NewService(pgRepo, redisCache, searchClient, logger)

	// Initialize worker (if enabled and search is available)
	var indexWorker *worker.Worker
	if cfg.Worker.Enabled && searchClient != nil {
		indexWorker = worker.NewWorker(
			pgRepo,
			searchClient,
			metricsInstance,
			cfg.Worker.Concurrency,
			logger,
		)
		if err := indexWorker.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start indexing worker")
		}
		defer indexWorker.Stop()
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	grpcHandler := grpcTransport.NewServer(listingsService, logger)
	pb.RegisterListingsServiceServer(grpcServer, grpcHandler)

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

	// Initialize HTTP server
	httpHandler := httpTransport.NewMinimalHandler(listingsService, logger)
	httpApp, err := httpTransport.StartMinimalServer(
		cfg.Server.HTTPHost,
		cfg.Server.HTTPPort,
		httpHandler,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start HTTP server")
	}
	defer httpApp.Shutdown()

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
