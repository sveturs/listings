package search_optimization

import (
	"backend/internal/proj/search_optimization/handler"
	"backend/internal/proj/search_optimization/service"
	"backend/internal/proj/search_optimization/storage"
	"backend/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type Module struct {
	Handler *handler.SearchOptimizationHandler
	Service service.SearchOptimizationService
	Storage storage.SearchOptimizationRepository
}

func NewModule(db *sqlx.DB, logger logger.Logger) *Module {
	// Создание репозитория
	repo := storage.NewPostgresSearchOptimizationRepository(db)

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
