package search_optimization

import (
	"backend/internal/middleware"
	"backend/internal/proj/search_optimization/handler"
	"backend/internal/proj/search_optimization/service"
	"backend/internal/proj/search_optimization/storage"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

type Module struct {
	Handler *handler.SearchOptimizationHandler
	Service service.SearchOptimizationService
	Storage storage.SearchOptimizationRepository
}

func NewModule(db *postgres.Database, logger logger.Logger) *Module {
	// Создание репозитория с pgxpool
	repo := storage.NewPostgresSearchOptimizationRepository(db.GetPool())

	// Создание сервиса
	svc := service.NewSearchOptimizationService(repo, logger)

	// Создание обработчика
	hdl := handler.NewSearchOptimizationHandler(svc)

	return &Module{
		Handler: hdl,
		Service: svc,
		Storage: repo,
	}
}

// RegisterRoutes регистрирует маршруты модуля поисковой оптимизации
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	// Группа для админских эндпоинтов поиска
	admin := app.Group("/api/v1/admin/search", middleware.JWTParser(), authMiddleware.RequireAuthString("admin"))

	// Эндпоинты для управления синонимами
	admin.Get("/synonyms", m.Handler.GetSynonyms)
	admin.Post("/synonyms", m.Handler.CreateSynonym)
	admin.Put("/synonyms/:id", m.Handler.UpdateSynonym)
	admin.Delete("/synonyms/:id", m.Handler.DeleteSynonym)

	// Эндпоинты для оптимизации
	admin.Post("/optimize-weights", m.Handler.StartOptimization)
	admin.Get("/optimization-status/:session_id", m.Handler.GetOptimizationStatus)
	admin.Post("/optimization-cancel/:session_id", m.Handler.CancelOptimization)
	admin.Post("/apply-weights", m.Handler.ApplyOptimizedWeights)
	admin.Post("/analyze-weights", m.Handler.AnalyzeCurrentWeights)
	admin.Get("/optimization-history", m.Handler.GetOptimizationHistory)
	admin.Get("/optimization-config", m.Handler.GetOptimizationConfig)
	admin.Put("/optimization-config", m.Handler.UpdateOptimizationConfig)
	admin.Post("/backup-weights", m.Handler.CreateWeightBackup)
	admin.Post("/rollback-weights", m.Handler.RollbackWeights)

	return nil
}

// GetPrefix возвращает префикс для модуля поисковой оптимизации
func (m *Module) GetPrefix() string {
	return "/api/v1/admin/search"
}
