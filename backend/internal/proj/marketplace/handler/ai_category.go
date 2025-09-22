package handler

import (
	"strconv"

	"backend/internal/proj/marketplace/repository"
	"backend/internal/proj/marketplace/services"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ConfirmDetectionRequest represents request for confirming category detection
type ConfirmDetectionRequest struct {
	CorrectCategoryID int32 `json:"correctCategoryId"`
}

type AICategoryHandler struct {
	detector         *services.AICategoryDetector
	validator        *services.AICategoryValidator
	keywordGenerator *services.AIKeywordGenerator
	keywordRepo      *repository.KeywordRepository
	learningSystem   *services.AILearningSystem
	logger           *zap.Logger
}

func NewAICategoryHandler(detector *services.AICategoryDetector, validator *services.AICategoryValidator, keywordGenerator *services.AIKeywordGenerator, keywordRepo *repository.KeywordRepository, learningSystem *services.AILearningSystem, logger *zap.Logger) *AICategoryHandler {
	return &AICategoryHandler{
		detector:         detector,
		validator:        validator,
		keywordGenerator: keywordGenerator,
		keywordRepo:      keywordRepo,
		learningSystem:   learningSystem,
		logger:           logger,
	}
}

// DetectCategory godoc
// @Summary Определение категории товара с помощью AI (с AI Fallback для 99% точности)
// @Description Определяет наиболее подходящую категорию для товара используя многоуровневый AI анализ с Fallback механизмом
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param request body services.AIDetectionInput true "Входные данные для определения категории"
// @Success 200 {object} utils.SuccessResponseSwag{data=services.AIDetectionResult} "Результат определения категории"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/detect-category [post]
func (h *AICategoryHandler) DetectCategory(c *fiber.Ctx) error {
	var input services.AIDetectionInput
	if err := c.BodyParser(&input); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Добавляем UserID из контекста если авторизован
	if user, ok := c.Locals("user").(map[string]interface{}); ok {
		if userID, ok := user["user_id"].(float64); ok {
			input.UserID = int32(userID)
		}
	}

	// Используем новый метод с AI Fallback для максимальной точности
	result, err := h.detector.DetectWithAIFallback(c.Context(), input)
	if err != nil {
		h.logger.Error("Failed to detect category with AI fallback", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.detectionFailed")
	}

	return utils.SuccessResponse(c, result)
}

// ConfirmDetection godoc
// @Summary Подтверждение правильности определения категории
// @Description Пользователь подтверждает или исправляет результат AI определения для улучшения модели
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param feedbackId path int true "ID записи обратной связи"
// @Param request body ConfirmDetectionRequest true "Правильная категория"
// @Success 200 {object} utils.SuccessResponseSwag "Обратная связь сохранена"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/confirm/{feedbackId} [post]
func (h *AICategoryHandler) ConfirmDetection(c *fiber.Ctx) error {
	feedbackIDStr := c.Params("feedbackId")
	feedbackID, err := strconv.ParseInt(feedbackIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidFeedbackId")
	}

	var req ConfirmDetectionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	err = h.detector.ConfirmDetection(c.Context(), feedbackID, req.CorrectCategoryID)
	if err != nil {
		h.logger.Error("Failed to confirm detection", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.confirmFailed")
	}

	// Запускаем обучение в фоне
	go func() {
		if err := h.detector.LearnFromFeedback(c.Context()); err != nil {
			h.logger.Error("Failed to learn from feedback", zap.Error(err))
		}
	}()

	return utils.SuccessResponse(c, nil)
}

// GetAccuracyMetrics godoc
// @Summary Получение метрик точности AI детекции
// @Description Возвращает статистику точности определения категорий за указанный период
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param days query int false "Количество дней для анализа" default(7)
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Метрики точности"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/metrics [get]
func (h *AICategoryHandler) GetAccuracyMetrics(c *fiber.Ctx) error {
	daysStr := c.Query("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		days = 7
	}

	metrics, err := h.detector.GetAccuracyMetrics(c.Context(), days)
	if err != nil {
		h.logger.Error("Failed to get accuracy metrics", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.metricsFailed")
	}

	return utils.SuccessResponse(c, metrics)
}

// TriggerLearning godoc
// @Summary Запуск процесса обучения модели
// @Description Запускает процесс обучения на основе накопленной обратной связи
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Обучение запущено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/learn [post]
func (h *AICategoryHandler) TriggerLearning(c *fiber.Ctx) error {
	// Запускаем обучение в фоне
	go func() {
		if err := h.detector.LearnFromFeedback(c.Context()); err != nil {
			h.logger.Error("Failed to learn from feedback", zap.Error(err))
		} else {
			h.logger.Info("Learning from feedback completed successfully")
		}
	}()

	return utils.SuccessResponse(c, nil)
}

// ValidateCategory godoc
// @Summary Валидация выбора категории через AI
// @Description Проверяет правильность выбора категории для товара используя AI анализ
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param request body services.ValidationRequest true "Данные для валидации категории"
// @Success 200 {object} utils.SuccessResponseSwag{data=services.ValidationResult} "Результат валидации"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/validate-category [post]
func (h *AICategoryHandler) ValidateCategory(c *fiber.Ctx) error {
	var req services.ValidationRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse validation request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Валидируем входные данные
	if req.Title == "" || req.CategoryName == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.missingRequiredFields")
	}

	// Проводим AI валидацию
	result, err := h.validator.ValidateCategory(c.Context(), req)
	if err != nil {
		h.logger.Error("AI validation failed", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.validationFailed")
	}

	// Логируем результат для мониторинга
	h.logger.Info("Category validation completed",
		zap.String("title", req.Title),
		zap.String("category", req.CategoryName),
		zap.Bool("isCorrect", result.IsCorrect),
		zap.Float64("confidence", result.Confidence),
		zap.String("reasoning", result.Reasoning))

	return utils.SuccessResponse(c, result)
}

// GenerateKeywords godoc
// @Summary Генерация ключевых слов для категории
// @Description Генерирует полный набор ключевых слов для указанной категории через AI
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param request body services.KeywordGenerationRequest true "Данные для генерации ключевых слов"
// @Success 200 {object} utils.SuccessResponseSwag{data=services.KeywordGenerationResult} "Результат генерации ключевых слов"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/generate-keywords [post]
func (h *AICategoryHandler) GenerateKeywords(c *fiber.Ctx) error {
	var req services.KeywordGenerationRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse keyword generation request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Валидируем входные данные
	if req.CategoryID == 0 || req.CategoryName == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.missingRequiredFields")
	}

	// Устанавливаем значения по умолчанию
	if req.Language == "" {
		req.Language = "ru"
	}
	if req.MinKeywords == 0 {
		req.MinKeywords = 50
	}

	// Генерируем ключевые слова через AI
	result, err := h.keywordGenerator.GenerateKeywordsForCategory(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to generate keywords", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.keywordGenerationFailed")
	}

	// Сохраняем ключевые слова в базу данных
	if len(result.Keywords) > 0 {
		// Convert services.GeneratedKeyword to repository.GeneratedKeyword
		repoKeywords := make([]repository.GeneratedKeyword, len(result.Keywords))
		for i, kw := range result.Keywords {
			repoKeywords[i] = repository.GeneratedKeyword{
				Keyword:     kw.Keyword,
				Type:        kw.Type,
				Weight:      kw.Weight,
				Confidence:  kw.Confidence,
				Description: kw.Description,
			}
		}
		err = h.keywordRepo.BulkInsertKeywords(c.Context(), req.CategoryID, repoKeywords, "ai_generated")
		if err != nil {
			h.logger.Error("Failed to save generated keywords", zap.Error(err))
			// Не возвращаем ошибку, так как генерация прошла успешно
		}
	}

	h.logger.Info("Keywords generated successfully",
		zap.Int32("categoryId", req.CategoryID),
		zap.String("categoryName", req.CategoryName),
		zap.Int("generatedCount", result.GeneratedCount))

	return utils.SuccessResponse(c, result)
}

// GenerateKeywordsForAllCategories godoc
// @Summary Массовая генерация ключевых слов для всех категорий
// @Description Генерирует ключевые слова для всех категорий, которым не хватает ключевых слов
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param minKeywords query int false "Минимальное количество ключевых слов на категорию" default(50)
// @Success 200 {object} utils.SuccessResponseSwag{data=services.KeywordGenerationResult} "Результат массовой генерации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/generate-keywords-all [post]
func (h *AICategoryHandler) GenerateKeywordsForAllCategories(c *fiber.Ctx) error {
	minKeywordsStr := c.Query("minKeywords", "50")
	minKeywords, err := strconv.Atoi(minKeywordsStr)
	if err != nil || minKeywords < 10 || minKeywords > 200 {
		minKeywords = 50
	}

	// Получаем категории, которым нужны ключевые слова
	categories, err := h.keywordRepo.GetCategoriesNeedingKeywords(c.Context(), minKeywords)
	if err != nil {
		h.logger.Error("Failed to get categories needing keywords", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.getCategoriesFailed")
	}

	if len(categories) == 0 {
		return utils.SuccessResponse(c, map[string]interface{}{
			"message":           "Все категории уже имеют достаточно ключевых слов",
			"categoriesFound":   0,
			"keywordsGenerated": 0,
		})
	}

	h.logger.Info("Starting bulk keyword generation",
		zap.Int("categoriesCount", len(categories)),
		zap.Int("minKeywords", minKeywords))

	// Конвертируем в нужный формат
	var serviceCategories []services.Category
	for _, cat := range categories {
		serviceCategories = append(serviceCategories, services.Category{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}

	// Запускаем массовую генерацию в фоне
	go func() {
		result, err := h.keywordGenerator.GenerateKeywordsForAllCategories(c.Context(), serviceCategories)
		if err != nil {
			h.logger.Error("Bulk keyword generation failed", zap.Error(err))
			return
		}

		// Сохраняем результаты в базу данных
		totalSaved := 0
		for _, categoryMapping := range result.Categories {
			if len(categoryMapping.Keywords) > 0 {
				// Convert services.GeneratedKeyword to repository.GeneratedKeyword
				repoKeywords := make([]repository.GeneratedKeyword, len(categoryMapping.Keywords))
				for i, kw := range categoryMapping.Keywords {
					repoKeywords[i] = repository.GeneratedKeyword{
						Keyword:     kw.Keyword,
						Type:        kw.Type,
						Weight:      kw.Weight,
						Confidence:  kw.Confidence,
						Description: kw.Description,
					}
				}
				err := h.keywordRepo.BulkInsertKeywords(c.Context(), categoryMapping.CategoryID, repoKeywords, "ai_generated_bulk")
				if err != nil {
					h.logger.Error("Failed to save bulk keywords",
						zap.Int32("categoryId", categoryMapping.CategoryID),
						zap.Error(err))
				} else {
					totalSaved += len(categoryMapping.Keywords)
				}
			}
		}

		h.logger.Info("Bulk keyword generation completed",
			zap.Int("categoriesProcessed", len(result.Categories)),
			zap.Int("totalKeywordsSaved", totalSaved),
			zap.Int64("processingTimeMs", result.ProcessingTimeMs))
	}()

	return utils.SuccessResponse(c, map[string]interface{}{
		"message":         "Массовая генерация ключевых слов запущена",
		"categoriesFound": len(categories),
		"status":          "processing",
	})
}

// GetKeywordStats godoc
// @Summary Статистика ключевых слов
// @Description Возвращает статистику по ключевым словам и их эффективности
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param categoryId query int false "ID категории для фильтрации"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Статистика ключевых слов"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/keyword-stats [get]
func (h *AICategoryHandler) GetKeywordStats(c *fiber.Ctx) error {
	categoryIDStr := c.Query("categoryId")

	// Получаем общую статистику по категориям
	counts, err := h.keywordRepo.GetKeywordCountByCategory(c.Context())
	if err != nil {
		h.logger.Error("Failed to get keyword counts", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.getStatsFailed")
	}

	// Получаем топ ключевых слов
	topKeywords, err := h.keywordRepo.GetTopKeywords(c.Context(), 20)
	if err != nil {
		h.logger.Error("Failed to get top keywords", zap.Error(err))
		topKeywords = []repository.KeywordAnalytics{} // Пустой массив при ошибке
	}

	stats := map[string]interface{}{
		"totalCategories":        len(counts),
		"keywordCountByCategory": counts,
		"topKeywords":            topKeywords,
	}

	// Если указана конкретная категория, добавляем детальную информацию
	if categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 32)
		if err == nil {
			keywords, err := h.keywordRepo.GetKeywordsByCategory(c.Context(), int32(categoryID))
			if err == nil {
				keywordsByType, _ := h.keywordRepo.GetKeywordsByTypes(c.Context(), int32(categoryID))
				stats["categoryKeywords"] = keywords
				stats["keywordsByType"] = keywordsByType
				stats["categoryKeywordCount"] = len(keywords)
			}
		}
	}

	return utils.SuccessResponse(c, stats)
}

// LearnFromFeedback godoc
// @Summary Запуск обучения системы на основе обратной связи
// @Description Анализирует обратную связь от AI валидации и улучшает систему категоризации
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=services.LearningMetrics} "Результаты обучения"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/learn-from-feedback [post]
func (h *AICategoryHandler) LearnFromFeedback(c *fiber.Ctx) error {
	if h.learningSystem == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "errors.learningSystemNotAvailable")
	}

	h.logger.Info("Manual learning session triggered")

	// Run learning in background for immediate response
	go func() {
		metrics, err := h.learningSystem.LearnFromValidationFeedback(c.Context())
		if err != nil {
			h.logger.Error("Learning from feedback failed", zap.Error(err))
		} else {
			h.logger.Info("Learning session completed",
				zap.Int("improvements", metrics.ImprovementsApplied),
				zap.Int("keywordsLearned", metrics.KeywordsLearned))
		}
	}()

	return utils.SuccessResponse(c, map[string]interface{}{
		"message": "Learning session started in background",
		"status":  "processing",
	})
}

// AutoImproveKeywords godoc
// @Summary Автоматическое улучшение ключевых слов
// @Description Автоматически улучшает покрытие ключевых слов для плохо работающих категорий
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Процесс улучшения запущен"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/auto-improve [post]
func (h *AICategoryHandler) AutoImproveKeywords(c *fiber.Ctx) error {
	if h.learningSystem == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "errors.learningSystemNotAvailable")
	}

	h.logger.Info("Auto-improvement process triggered")

	// Run improvement in background
	go func() {
		err := h.learningSystem.AutoImproveKeywords(c.Context())
		if err != nil {
			h.logger.Error("Auto-improvement failed", zap.Error(err))
		} else {
			h.logger.Info("Auto-improvement completed successfully")
		}
	}()

	return utils.SuccessResponse(c, map[string]interface{}{
		"message": "Auto-improvement process started",
		"status":  "processing",
	})
}

// GetLearningStats godoc
// @Summary Статистика обучения AI системы
// @Description Возвращает статистику и метрики обучения саморазвивающейся системы
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=services.LearningMetrics} "Статистика обучения"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/learning-stats [get]
func (h *AICategoryHandler) GetLearningStats(c *fiber.Ctx) error {
	if h.learningSystem == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "errors.learningSystemNotAvailable")
	}

	stats, err := h.learningSystem.GetLearningStats(c.Context())
	if err != nil {
		h.logger.Error("Failed to get learning stats", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.getLearningStatsFailed")
	}

	return utils.SuccessResponse(c, stats)
}

// ScheduledLearning godoc
// @Summary Запуск планового обучения системы
// @Description Выполняет все плановые задачи обучения: анализ обратной связи, улучшение ключевых слов
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Плановое обучение запущено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/scheduled-learning [post]
func (h *AICategoryHandler) ScheduledLearning(c *fiber.Ctx) error {
	if h.learningSystem == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "errors.learningSystemNotAvailable")
	}

	h.logger.Info("Scheduled learning session triggered")

	// Run scheduled learning in background
	go func() {
		err := h.learningSystem.ScheduledLearning(c.Context())
		if err != nil {
			h.logger.Error("Scheduled learning failed", zap.Error(err))
		} else {
			h.logger.Info("Scheduled learning completed successfully")
		}
	}()

	return utils.SuccessResponse(c, map[string]interface{}{
		"message": "Scheduled learning started in background",
		"status":  "processing",
	})
}

// DetectCategoryStandard godoc
// @Summary Определение категории только стандартными алгоритмами (без AI Fallback)
// @Description Определяет категорию используя только keyword matching и similarity без AI Fallback для сравнения
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param request body services.AIDetectionInput true "Входные данные для определения категории"
// @Success 200 {object} utils.SuccessResponseSwag{data=services.AIDetectionResult} "Результат определения категории"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/detect-category-standard [post]
func (h *AICategoryHandler) DetectCategoryStandard(c *fiber.Ctx) error {
	var input services.AIDetectionInput
	if err := c.BodyParser(&input); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Добавляем UserID из контекста если авторизован
	if user, ok := c.Locals("user").(map[string]interface{}); ok {
		if userID, ok := user["user_id"].(float64); ok {
			input.UserID = int32(userID)
		}
	}

	// Используем только стандартный метод без AI Fallback
	result, err := h.detector.DetectCategory(c.Context(), input)
	if err != nil {
		h.logger.Error("Failed to detect category with standard method", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.detectionFailed")
	}

	return utils.SuccessResponse(c, result)
}

// SelectCategory godoc
// @Summary Прямой выбор категории через AI из полного списка
// @Description AI анализирует товар и выбирает наиболее подходящую категорию из всех доступных (метод максимальной точности)
// @Tags marketplace-ai
// @Accept json
// @Produce json
// @Param request body services.AIDetectionInput true "Входные данные товара"
// @Success 200 {object} utils.SuccessResponseSwag{data=services.AIDetectionResult} "Результат выбора категории с обоснованием"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /marketplace/ai/select-category [post]
func (h *AICategoryHandler) SelectCategory(c *fiber.Ctx) error {
	var input services.AIDetectionInput
	if err := c.BodyParser(&input); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Проверяем обязательные поля
	if input.Title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.titleRequired")
	}

	// Добавляем UserID из контекста если авторизован
	if user, ok := c.Locals("user").(map[string]interface{}); ok {
		if userID, ok := user["user_id"].(float64); ok {
			input.UserID = int32(userID)
		}
	}

	h.logger.Info("Direct AI category selection requested",
		zap.String("title", input.Title),
		zap.Int("descriptionLength", len(input.Description)))

	// Используем метод прямого выбора категории через AI
	result, err := h.detector.SelectCategoryDirectly(c.Context(), input)
	if err != nil {
		h.logger.Error("Failed to select category via AI",
			zap.Error(err),
			zap.String("errorDetails", err.Error()))
		// Возвращаем подробную ошибку в режиме разработки
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "AI selection failed",
			"details": err.Error(),
		})
	}

	h.logger.Info("AI category selection successful",
		zap.Int32("categoryId", result.CategoryID),
		zap.String("categoryName", result.CategoryName),
		zap.Float64("confidence", result.ConfidenceScore))

	return utils.SuccessResponse(c, result)
}
