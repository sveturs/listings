package search_optimization

import (
	"backend/internal/middleware"
	"backend/internal/proj/search_optimization/handler"
	"backend/internal/proj/search_optimization/service"
	"backend/internal/proj/search_optimization/storage"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
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
	// Группа для админских эндпоинтов поиска (избегаем конфликта с users admin)
	admin := app.Group("/api/v1/search-admin")
	// Временно убираем авторизацию для тестирования
	// admin.Use(middleware.AuthRequiredJWT)
	// admin.Use(middleware.AdminRequired)

	// Группа для поисковой оптимизации
	search := admin.Group("/search")

	// Эндпоинты для управления синонимами
	synonyms := search.Group("/synonyms")
	synonyms.Get("/", m.Handler.GetSynonyms)
	synonyms.Post("/", m.Handler.CreateSynonym)
	synonyms.Put("/:id", m.Handler.UpdateSynonym)
	synonyms.Delete("/:id", m.Handler.DeleteSynonym)

	// Эндпоинты для оптимизации
	optimization := search.Group("/optimization")
	optimization.Post("/start", m.Handler.StartOptimization)
	optimization.Get("/status/:id", m.Handler.GetOptimizationStatus)
	optimization.Post("/cancel/:id", m.Handler.CancelOptimization)
	optimization.Post("/apply/:id", m.Handler.ApplyOptimizedWeights)
	optimization.Get("/analyze", m.Handler.AnalyzeCurrentWeights)
	optimization.Get("/history", m.Handler.GetOptimizationHistory)
	optimization.Get("/config", m.Handler.GetOptimizationConfig)
	optimization.Put("/config", m.Handler.UpdateOptimizationConfig)
	optimization.Post("/backup", m.Handler.CreateWeightBackup)
	optimization.Post("/rollback/:id", m.Handler.RollbackWeights)

	return nil
}

// GetPrefix возвращает префикс для модуля поисковой оптимизации
func (m *Module) GetPrefix() string {
	return "/api/v1/admin/search"
}
