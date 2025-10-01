package logistics

import (
	"database/sql"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/proj/admin/logistics/handler"
	"backend/internal/proj/admin/logistics/service"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль админки логистики
type Module struct {
	dashboardHandler  *handler.DashboardHandler
	shipmentsHandler  *handler.ShipmentsHandler
	problemsHandler   *handler.ProblemsHandler
	analyticsHandler  *handler.AnalyticsHandler
	monitoringService *service.MonitoringService
	problemService    *service.ProblemService
	analyticsService  *service.AnalyticsService
}

// NewModule создает новый модуль админки логистики
func NewModule(db *sql.DB, cfg *config.Config, logger *logger.Logger) (*Module, error) {
	// Создаем сервисы
	monitoringService := service.NewMonitoringService(db)
	problemService := service.NewProblemService(db)
	analyticsService := service.NewAnalyticsService(db)

	// Создаем handlers
	dashboardHandler := handler.NewDashboardHandler(monitoringService)
	shipmentsHandler := handler.NewShipmentsHandler(monitoringService, problemService)
	problemsHandler := handler.NewProblemsHandler(problemService, logger)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	return &Module{
		dashboardHandler:  dashboardHandler,
		shipmentsHandler:  shipmentsHandler,
		problemsHandler:   problemsHandler,
		analyticsHandler:  analyticsHandler,
		monitoringService: monitoringService,
		problemService:    problemService,
		analyticsService:  analyticsService,
	}, nil
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Создаем группу для админки логистики
	// ВАЖНО: middleware JWTParser и AdminRequired уже применены от родительской группы /api/v1/admin (marketplace)
	// Поэтому НЕ нужно добавлять их здесь повторно
	adminLogistics := app.Group("/api/v1/admin/logistics")

	// Dashboard роуты
	adminLogistics.Get("/dashboard", m.dashboardHandler.GetDashboardStats)
	adminLogistics.Get("/dashboard/chart", m.dashboardHandler.GetWeeklyChart)

	// Shipments роуты
	adminLogistics.Get("/shipments", m.shipmentsHandler.GetShipments)
	adminLogistics.Get("/shipments/:id", m.shipmentsHandler.GetShipmentDetails)
	adminLogistics.Get("/shipments/:provider/:id", m.shipmentsHandler.GetShipmentDetailsByProvider)
	adminLogistics.Put("/shipments/:id/status", m.shipmentsHandler.UpdateShipmentStatus)
	adminLogistics.Post("/shipments/:id/action", m.shipmentsHandler.PerformShipmentAction)

	// Problems роуты
	adminLogistics.Get("/problems", m.problemsHandler.GetProblems)
	adminLogistics.Post("/problems", m.problemsHandler.CreateProblem)
	adminLogistics.Put("/problems/:id", m.problemsHandler.UpdateProblem)
	adminLogistics.Post("/problems/:id/resolve", m.problemsHandler.ResolveProblem)
	adminLogistics.Post("/problems/:id/assign", m.problemsHandler.AssignProblem)
	adminLogistics.Get("/problems/:id/details", m.problemsHandler.GetProblemDetails)
	adminLogistics.Get("/problems/:id/comments", m.problemsHandler.GetProblemComments)
	adminLogistics.Post("/problems/:id/comments", m.problemsHandler.AddProblemComment)
	adminLogistics.Get("/problems/:id/history", m.problemsHandler.GetProblemHistory)

	// Analytics роуты
	adminLogistics.Get("/analytics/performance", m.analyticsHandler.GetPerformanceMetrics)
	adminLogistics.Get("/analytics/financial", m.analyticsHandler.GetFinancialReport)
	adminLogistics.Get("/analytics/export", m.analyticsHandler.ExportReport)
	adminLogistics.Get("/analytics/couriers", m.analyticsHandler.GetCourierComparison)

	return nil
}

// GetPrefix возвращает префикс маршрутов модуля
func (m *Module) GetPrefix() string {
	return "/api/v1/admin/logistics"
}

// GetMonitoringService возвращает сервис мониторинга для использования в других модулях
func (m *Module) GetMonitoringService() *service.MonitoringService {
	return m.monitoringService
}

// GetProblemService возвращает сервис проблем
func (m *Module) GetProblemService() *service.ProblemService {
	return m.problemService
}

// GetAnalyticsService возвращает сервис аналитики
func (m *Module) GetAnalyticsService() *service.AnalyticsService {
	return m.analyticsService
}

// StartBackgroundTasks запускает фоновые задачи модуля
func (m *Module) StartBackgroundTasks() error {
	// TODO: Запустить cron задачи для:
	// - Автоматического создания проблем для задержанных отправлений
	// - Обновления метрик в таблице logistics_metrics
	// - Очистки устаревшего кеша
	// - Отправки ежедневных/еженедельных отчетов

	return nil
}
