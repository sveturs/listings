package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"backend/internal/proj/marketplace/services"
	"backend/pkg/utils"
)

// CategoryDetectorHandler обработчик для умного поиска категорий
type CategoryDetectorHandler struct {
	detector *services.CategoryDetector
	logger   *zap.Logger
}

// NewCategoryDetectorHandler создает новый обработчик
func NewCategoryDetectorHandler(detector *services.CategoryDetector, logger *zap.Logger) *CategoryDetectorHandler {
	return &CategoryDetectorHandler{
		detector: detector,
		logger:   logger,
	}
}

// DetectCategory определяет категорию по семантическим данным
// @Summary Определить категорию по ключевым словам и атрибутам
// @Description Использует комбинированный подход с ключевыми словами и similarity search
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param body body DetectCategoryRequest true "Данные для определения категории"
// @Success 200 {object} utils.SuccessResponseSwag{data=DetectCategoryResponse} "Результат определения категории"
// @Failure 400 {object} utils.ErrorResponseSwag "Недостаточно данных"
// @Failure 500 {object} utils.ErrorResponseSwag "Ошибка сервера"
// @Router /api/v1/marketplace/categories/detect [post]
func (h *CategoryDetectorHandler) DetectCategory(c *fiber.Ctx) error {
	// Используем fmt для гарантированного вывода
	fmt.Println("=== DetectCategory METHOD CALLED ===")
	
	// Проверка на nil в самом начале
	if h == nil {
		fmt.Println("ERROR: CategoryDetectorHandler is nil!")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "handler is nil")
	}
	
	fmt.Println("Handler is not nil, checking detector...")
	
	if h.detector == nil {
		fmt.Println("ERROR: detector is nil!")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "detector is nil")
	}
	
	fmt.Println("Detector is not nil, checking logger...")
	
	// Логируем используя глобальный logger если h.logger nil
	if h.logger == nil {
		fmt.Println("Logger is nil, using global zap logger")
		zap.L().Info(">>>>>> DetectCategory method called! (global logger) <<<<<<")
	} else {
		fmt.Println("Using handler logger")
		// Попробуем вызвать logger в try-catch
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("PANIC when using logger: %v\n", r)
				}
			}()
			h.logger.Info(">>>>>> DetectCategory method called! <<<<<<")
		}()
	}
	
	fmt.Println("Logger call completed, parsing request body...")
	
	var req DetectCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
		// Используем глобальный logger для безопасности
		zap.L().Error("ошибка парсинга запроса", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.validation.invalidRequest")
	}
	
	fmt.Printf("Request parsed successfully. Keywords: %d, Title: %s\n", len(req.Keywords), req.Title)
	
	// Безопасный вызов logger
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC when logging parsed request: %v\n", r)
			}
		}()
		if h.logger != nil {
			h.logger.Info("запрос распарсен", 
				zap.Int("keywords_count", len(req.Keywords)),
				zap.String("title", req.Title))
		} else {
			zap.L().Info("запрос распарсен", 
				zap.Int("keywords_count", len(req.Keywords)),
				zap.String("title", req.Title))
		}
	}()

	fmt.Println("Starting validation...")
	
	// Валидация
	if len(req.Keywords) == 0 && req.Title == "" && req.Description == "" {
		fmt.Println("Validation failed: insufficient data")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.validation.insufficientData")
	}
	
	fmt.Println("Validation passed, getting user ID...")

	// Получаем user_id из контекста
	userID := c.Locals("userID")
	var userIDPtr *int32
	if userID != nil {
		if id, ok := userID.(int); ok {
			id32 := int32(id)
			userIDPtr = &id32
		}
	}

	// Формируем входные данные для детектора
	input := services.DetectionInput{
		Keywords:    req.Keywords,
		Attributes:  req.Attributes,
		Domain:      req.Domain,
		ProductType: req.ProductType,
		Language:    req.Language,
		Title:       req.Title,
		Description: req.Description,
		UserID:      userIDPtr,
		SessionID:   c.Get("X-Session-ID", ""),
	}

	// Определяем категорию
	fmt.Println("About to call detector.DetectCategory...")
	
	// Безопасный вызов logger
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC when logging detector call: %v\n", r)
			}
		}()
		if h.logger != nil {
			h.logger.Info("вызов detector.DetectCategory", zap.Any("input", input))
		} else {
			zap.L().Info("вызов detector.DetectCategory", zap.Any("input", input))
		}
	}()
	
	fmt.Println("Calling detector.DetectCategory method...")
	result, err := h.detector.DetectCategory(c.Context(), input)
	
	if err != nil {
		fmt.Printf("Detector returned error: %v\n", err)
		// Безопасное логирование ошибки
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("PANIC when logging error: %v\n", r)
				}
			}()
			if h.logger != nil {
				h.logger.Error("ошибка определения категории", zap.Error(err))
			} else {
				zap.L().Error("ошибка определения категории", zap.Error(err))
			}
		}()
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.marketplace.categoryDetectionFailed")
	}
	
	fmt.Printf("Detector returned result: categoryID=%d\n", result.CategoryID)
	
	h.logger.Info("категория определена", zap.Int32("categoryID", result.CategoryID))

	// Формируем ответ
	response := DetectCategoryResponse{
		CategoryID:       result.CategoryID,
		CategoryName:     result.CategoryName,
		CategorySlug:     result.CategorySlug,
		ConfidenceScore:  result.ConfidenceScore,
		Method:           result.Method,
		ProcessingTimeMs: result.ProcessingTimeMs,
		StatsID:          result.StatsID,
	}

	// Добавляем предупреждение, если использована категория "Прочее"
	const otherCategoryID = 9999
	if result.CategoryID == otherCategoryID {
		response.Warning = "Не удалось автоматически определить подходящую категорию. Объявление будет размещено в категории 'Прочее'. Пожалуйста, выберите более подходящую категорию вручную."
		h.logger.Warn("использована категория 'Прочее'",
			zap.Strings("keywords", req.Keywords),
			zap.String("title", req.Title))
	}

	// Добавляем альтернативные категории
	if len(result.AlternativeCategories) > 0 {
		response.AlternativeCategories = make([]AlternativeCategoryResponse, len(result.AlternativeCategories))
		for i, alt := range result.AlternativeCategories {
			response.AlternativeCategories[i] = AlternativeCategoryResponse{
				CategoryID:      alt.CategoryID,
				CategoryName:    alt.CategoryName,
				CategorySlug:    alt.CategorySlug,
				ConfidenceScore: alt.ConfidenceScore,
			}
		}
	}

	// Добавляем детали для отладки (только если запрошено)
	if req.IncludeDebugInfo {
		response.DebugInfo = &DebugInfo{
			KeywordScore:    result.KeywordScore,
			SimilarityScore: result.SimilarityScore,
			MatchedKeywords: result.MatchedKeywords,
		}
	}

	return utils.SuccessResponse(c, response)
}

// UpdateCategoryConfirmation обновляет подтверждение пользователя
// @Summary Подтвердить или исправить выбор категории
// @Description Обновляет статистику для улучшения алгоритма
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param stats_id path int true "ID записи статистики"
// @Param body body UpdateConfirmationRequest true "Данные подтверждения"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]bool} "Успешное обновление"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Ошибка сервера"
// @Router /api/v1/marketplace/categories/detect/{stats_id}/confirm [put]
func (h *CategoryDetectorHandler) UpdateCategoryConfirmation(c *fiber.Ctx) error {
	// Получаем ID статистики
	statsIDStr := c.Params("stats_id")
	_, err := strconv.ParseInt(statsIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.validation.invalidStatsID")
	}

	var req UpdateConfirmationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.validation.invalidRequest")
	}

	// TODO: Вызвать метод обновления статистики
	// statsID, _ := strconv.ParseInt(statsIDStr, 10, 32)
	// statsRepo.UpdateUserFeedback(c.Context(), int32(statsID), req.Confirmed, req.SelectedCategoryID)

	return utils.SuccessResponse(c, map[string]bool{"updated": true})
}

// GetCategoryKeywords получает ключевые слова для категории
// @Summary Получить ключевые слова категории
// @Description Возвращает список ключевых слов с весами для указанной категории
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param category_id path int true "ID категории"
// @Param language query string false "Язык (ru, en, sr)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]CategoryKeywordResponse} "Список ключевых слов"
// @Failure 404 {object} utils.ErrorResponseSwag "Категория не найдена"
// @Failure 500 {object} utils.ErrorResponseSwag "Ошибка сервера"
// @Router /api/v1/marketplace/categories/{category_id}/keywords [get]
func (h *CategoryDetectorHandler) GetCategoryKeywords(c *fiber.Ctx) error {
	// Получаем ID категории
	categoryIDStr := c.Params("category_id")
	_, err := strconv.ParseInt(categoryIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.validation.invalidCategoryID")
	}

	// TODO: Получить ключевые слова из репозитория
	// categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 32)
	// keywords, err := keywordRepo.GetKeywordsByCategoryID(c.Context(), int32(categoryID))

	// Заглушка для примера
	keywords := []CategoryKeywordResponse{
		{
			ID:          1,
			Keyword:     "телефон",
			Language:    "ru",
			Weight:      10.0,
			KeywordType: "main",
			UsageCount:  150,
			SuccessRate: 0.85,
		},
	}

	// Фильтруем по языку если указан
	language := c.Query("language")
	if language != "" {
		filtered := make([]CategoryKeywordResponse, 0)
		for _, kw := range keywords {
			if kw.Language == language || kw.Language == "*" {
				filtered = append(filtered, kw)
			}
		}
		keywords = filtered
	}

	return utils.SuccessResponse(c, keywords)
}

// Request/Response структуры

type DetectCategoryRequest struct {
	Keywords         []string               `json:"keywords" example:"телефон,смартфон,айфон"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
	Domain           string                 `json:"domain,omitempty" example:"electronics"`
	ProductType      string                 `json:"product_type,omitempty" example:"smartphone"`
	Language         string                 `json:"language,omitempty" example:"ru"`
	Title            string                 `json:"title,omitempty"`
	Description      string                 `json:"description,omitempty"`
	IncludeDebugInfo bool                   `json:"include_debug_info,omitempty"`
}

type DetectCategoryResponse struct {
	CategoryID            int32                         `json:"category_id"`
	CategoryName          string                        `json:"category_name"`
	CategorySlug          string                        `json:"category_slug"`
	ConfidenceScore       float64                       `json:"confidence_score"`
	Method                string                        `json:"method"`
	AlternativeCategories []AlternativeCategoryResponse `json:"alternative_categories,omitempty"`
	ProcessingTimeMs      int64                         `json:"processing_time_ms"`
	DebugInfo             *DebugInfo                    `json:"debug_info,omitempty"`
	StatsID               int32                         `json:"stats_id,omitempty"`
	Warning               string                        `json:"warning,omitempty"`
}

type AlternativeCategoryResponse struct {
	CategoryID      int32   `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	CategorySlug    string  `json:"category_slug"`
	ConfidenceScore float64 `json:"confidence_score"`
}

type DebugInfo struct {
	KeywordScore    float64  `json:"keyword_score"`
	SimilarityScore float64  `json:"similarity_score"`
	MatchedKeywords []string `json:"matched_keywords,omitempty"`
}

type UpdateConfirmationRequest struct {
	Confirmed          bool   `json:"confirmed"`
	SelectedCategoryID *int32 `json:"selected_category_id,omitempty"`
}

type CategoryKeywordResponse struct {
	ID          int32   `json:"id"`
	Keyword     string  `json:"keyword"`
	Language    string  `json:"language"`
	Weight      float64 `json:"weight"`
	KeywordType string  `json:"keyword_type"`
	IsNegative  bool    `json:"is_negative,omitempty"`
	UsageCount  int32   `json:"usage_count"`
	SuccessRate float64 `json:"success_rate"`
}
