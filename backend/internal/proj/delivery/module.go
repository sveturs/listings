package delivery

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/config"
	"backend/internal/middleware"
	adminLogistics "backend/internal/proj/admin/logistics/service"
	"backend/internal/proj/delivery/grpcclient"
	"backend/internal/proj/delivery/handler"
	"backend/internal/proj/delivery/service"
	notifService "backend/internal/proj/notifications/service"
	"backend/pkg/logger"
)

// Module - модуль универсальной системы доставки
type Module struct {
	handler      *handler.Handler
	adminHandler *handler.AdminHandler
	service      *service.Service
	grpcClient   *grpcclient.Client
	// Сервисы из admin/logistics для консолидации
	monitoringService *adminLogistics.MonitoringService
	problemService    *adminLogistics.ProblemService
	analyticsService  *adminLogistics.AnalyticsService
}

// NewModule создает новый модуль доставки
func NewModule(db *sqlx.DB, cfg *config.Config, logger *logger.Logger) (*Module, error) {
	// Используем дефолтное значение если не задано
	grpcURL := cfg.DeliveryGRPCURL
	if grpcURL == "" {
		grpcURL = "svetu.rs:30051" // Default: production delivery microservice
		log.Info().Str("url", grpcURL).Msg("Using default delivery gRPC URL")
	}

	// Создаем gRPC клиент для delivery микросервиса (ОБЯЗАТЕЛЬНО)
	grpcClient, err := grpcclient.NewClient(grpcURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to delivery gRPC service at %s: %w", grpcURL, err)
	}
	log.Info().Str("url", grpcURL).Msg("Successfully connected to delivery gRPC service")

	// Создаем сервис с обязательным gRPC клиентом
	svc := service.NewService(db, grpcClient)

	// Создаем основной обработчик (он больше не нуждается в providerFactory)
	h := handler.NewHandler(svc)

	// Получаем admin обработчик из основного handler
	adminHandler := h.GetAdminHandler()

	// Создаем сервисы из admin/logistics для полной функциональности
	sqlDB := db.DB // sqlx.DB содержит *sql.DB
	monitoringService := adminLogistics.NewMonitoringService(sqlDB)
	problemService := adminLogistics.NewProblemService(sqlDB)
	analyticsService := adminLogistics.NewAnalyticsService(sqlDB)

	// Устанавливаем сервисы в admin handler для консолидации
	adminHandler.SetLogisticsServices(monitoringService, problemService, analyticsService, logger)

	return &Module{
		handler:           h,
		adminHandler:      adminHandler,
		service:           svc,
		grpcClient:        grpcClient,
		monitoringService: monitoringService,
		problemService:    problemService,
		analyticsService:  analyticsService,
	}, nil
}

// SetNotificationService устанавливает сервис уведомлений
func (m *Module) SetNotificationService(notifSvc notifService.NotificationServiceInterface) {
	if m.service != nil {
		m.service.SetNotificationService(notifSvc)
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Создаем группу для API v1 с авторизацией (для основных роутов)
	api := app.Group("/api/v1", mw.JWTParser(), authMiddleware.RequireAuth())

	// Регистрируем защищенные роуты
	m.handler.RegisterRoutes(api)

	// Регистрируем публичные webhook роуты (без авторизации)
	webhookGroup := app.Group("/api/v1/delivery/webhooks")
	m.handler.RegisterWebhookRoutes(webhookGroup)

	// Регистрируем админские роуты (консолидация из admin/logistics)
	adminGroup := app.Group("/api/v1/admin/delivery")
	adminGroup.Use(mw.JWTParser(), authMiddleware.RequireAuthString("admin"))
	m.adminHandler.RegisterConsolidatedRoutes(adminGroup, mw)

	return nil
}

// GetPrefix возвращает префикс маршрутов модуля
func (m *Module) GetPrefix() string {
	return "/api/v1/delivery"
}
