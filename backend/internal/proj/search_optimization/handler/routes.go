package handler

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterSearchOptimizationRoutes(api fiber.Router, h *SearchOptimizationHandler) {
	searchOptimization := api.Group("/search")

	// Основные операции оптимизации
	searchOptimization.Post("/optimize-weights", h.StartOptimization)
	searchOptimization.Get("/optimization-status/:session_id", h.GetOptimizationStatus)
	searchOptimization.Post("/optimization-cancel/:session_id", h.CancelOptimization)

	// Применение и управление весами
	searchOptimization.Post("/apply-weights", h.ApplyOptimizedWeights)
	searchOptimization.Post("/backup-weights", h.CreateWeightBackup)
	searchOptimization.Post("/rollback-weights", h.RollbackWeights)

	// Анализ и статистика
	searchOptimization.Post("/analyze-weights", h.AnalyzeCurrentWeights)
	searchOptimization.Get("/optimization-history", h.GetOptimizationHistory)

	// Конфигурация
	searchOptimization.Get("/optimization-config", h.GetOptimizationConfig)
	searchOptimization.Put("/optimization-config", h.UpdateOptimizationConfig)

	// Управление синонимами
	searchOptimization.Get("/synonyms", h.GetSynonyms)
	searchOptimization.Post("/synonyms", h.CreateSynonym)
	searchOptimization.Put("/synonyms/:id", h.UpdateSynonym)
	searchOptimization.Delete("/synonyms/:id", h.DeleteSynonym)
}
