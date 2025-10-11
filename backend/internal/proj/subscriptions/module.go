package subscriptions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/pkg/allsecure"
	paymentsService "backend/internal/proj/payments/service"
	"backend/internal/proj/subscriptions/handler"
	"backend/internal/proj/subscriptions/service"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"
)

// Module представляет модуль подписок
type Module struct {
	handler *handler.SubscriptionHandler
	service *service.SubscriptionService
	logger  *logger.Logger
}

// NewModule создает новый модуль подписок
func NewModule(db *sqlx.DB, paymentService *paymentsService.AllSecureService, allSecureClient *allsecure.Client, logger *logger.Logger, jwtParserMW fiber.Handler) *Module {
	// Создаем репозиторий
	repo := postgres.NewSubscriptionRepository(db)

	// Создаем сервис
	subscriptionService := service.NewSubscriptionService(
		repo,
		paymentService,
		allSecureClient,
		logger,
	)

	// Создаем handler
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, logger, jwtParserMW)

	return &Module{
		handler: subscriptionHandler,
		service: subscriptionService,
		logger:  logger,
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	m.handler.RegisterRoutes(app, middleware)
	m.logger.Info("Subscription routes registered")
	return nil
}

// GetPrefix возвращает префикс модуля
func (m *Module) GetPrefix() string {
	return "subscriptions"
}

// GetService возвращает сервис подписок (для использования другими модулями)
func (m *Module) GetService() *service.SubscriptionService {
	return m.service
}
