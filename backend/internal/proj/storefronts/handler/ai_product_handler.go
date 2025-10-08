package handler

import (
	"backend/internal/proj/ai/handler"
	marketplaceServices "backend/internal/proj/marketplace/services"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AIProductHandler handles AI-related endpoints for storefront products
type AIProductHandler struct {
	categoryDetector *marketplaceServices.AICategoryDetector
	aiHandler        *handler.Handler
	logger           *zap.Logger
}

// NewAIProductHandler creates a new AI product handler
func NewAIProductHandler(
	categoryDetector *marketplaceServices.AICategoryDetector,
	aiHandler *handler.Handler,
	logger *zap.Logger,
) *AIProductHandler {
	return &AIProductHandler{
		categoryDetector: categoryDetector,
		aiHandler:        aiHandler,
		logger:           logger,
	}
}

// AnalyzeProductImageRequest represents the request for analyzing product image
type AnalyzeProductImageRequest struct {
	ImageData string `json:"imageData"` // base64 encoded image
	Language  string `json:"language"`  // ru, en, sr
}

// DetectCategoryRequest represents the request for category detection
type DetectCategoryRequest struct {
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	AIHints     *marketplaceServices.AIHints `json:"aiHints,omitempty"`
	Language    string                       `json:"language"`
}

// ABTestTitlesRequest represents the request for A/B testing titles
type ABTestTitlesRequest struct {
	TitleVariants []string `json:"titleVariants"`
}

// TranslateContentRequest represents the request for content translation
type TranslateContentRequest struct {
	Content struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"content"`
	TargetLanguages []string `json:"targetLanguages"`
	SourceLanguage  string   `json:"sourceLanguage"`
}

// AnalyzeProductImage godoc
// @Summary Analyze product image with AI for storefront products
// @Description Uses Claude AI to analyze product image and extract title, description, category hints, price, etc.
// @Tags AI Storefronts
// @Accept json
// @Produce json
// @Param request body AnalyzeProductImageRequest true "Image data"
// @Success 200 {object} utils.SuccessResponseSwag{data=handler.AnalyzeProductResponse} "AI analysis result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/ai/analyze-product-image [post]
func (h *AIProductHandler) AnalyzeProductImage(c *fiber.Ctx) error {
	var req AnalyzeProductImageRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if req.ImageData == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.imageDataRequired")
	}

	if req.Language == "" {
		req.Language = "ru"
	}

	// Вызываем существующий AI handler (переиспользуем!)
	// Передаем через контекст что это storefront product
	c.Locals("entity_type", "product")

	// Используем существующий метод AnalyzeProduct из AI handler
	return h.aiHandler.AnalyzeProduct(c)
}

// DetectCategory godoc
// @Summary Detect category for storefront product
// @Description Uses AI to detect the most suitable category for a storefront product
// @Tags AI Storefronts
// @Accept json
// @Produce json
// @Param request body DetectCategoryRequest true "Product details"
// @Success 200 {object} utils.SuccessResponseSwag{data=marketplaceServices.AIDetectionResult} "Category detection result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/ai/detect-category [post]
func (h *AIProductHandler) DetectCategory(c *fiber.Ctx) error {
	var req DetectCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if req.Title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.titleRequired")
	}

	// Создаем input для детектора с указанием entity_type = "product"
	input := marketplaceServices.AIDetectionInput{
		Title:       req.Title,
		Description: req.Description,
		AIHints:     req.AIHints,
		EntityType:  "product", // ВАЖНО: указываем что это product!
	}

	// Добавляем UserID из контекста если авторизован
	if user, ok := c.Locals("user").(map[string]interface{}); ok {
		if userID, ok := user["user_id"].(float64); ok {
			input.UserID = int32(userID)
		}
	}

	// Вызываем детектор с AI Fallback
	result, err := h.categoryDetector.DetectWithAIFallback(c.Context(), input)
	if err != nil {
		h.logger.Error("Failed to detect category for product", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.detectionFailed")
	}

	return utils.SuccessResponse(c, result)
}

// ABTestTitles godoc
// @Summary A/B test product titles
// @Description Evaluates multiple title variants and returns the best one
// @Tags AI Storefronts
// @Accept json
// @Produce json
// @Param request body ABTestTitlesRequest true "Title variants"
// @Success 200 {object} utils.SuccessResponseSwag "A/B test result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/ai/ab-test-titles [post]
func (h *AIProductHandler) ABTestTitles(c *fiber.Ctx) error {
	var req ABTestTitlesRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if len(req.TitleVariants) < 2 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.minTwoTitlesRequired")
	}

	// Переиспользуем метод из AI handler
	return h.aiHandler.PerformABTest(c)
}

// TranslateContent godoc
// @Summary Translate product content to multiple languages
// @Description Translates product title and description to specified languages
// @Tags AI Storefronts
// @Accept json
// @Produce json
// @Param request body TranslateContentRequest true "Content to translate"
// @Success 200 {object} utils.SuccessResponseSwag "Translation result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/ai/translate-content [post]
func (h *AIProductHandler) TranslateContent(c *fiber.Ctx) error {
	var req TranslateContentRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request", zap.Error(err))
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if req.Content.Title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.titleRequired")
	}

	if len(req.TargetLanguages) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.targetLanguagesRequired")
	}

	// Переиспользуем метод из AI handler
	return h.aiHandler.TranslateContent(c)
}

// GetMetrics godoc
// @Summary Get AI detection metrics for storefront products
// @Description Returns accuracy and performance metrics for AI category detection
// @Tags AI Storefronts
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Metrics data"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/ai/metrics [get]
func (h *AIProductHandler) GetMetrics(c *fiber.Ctx) error {
	// TODO: Implement metrics collection for storefront products
	// For now, return basic placeholder
	metrics := map[string]interface{}{
		"entity_type":       "product",
		"total_detections":  0,
		"accuracy":          0.0,
		"cache_hit_rate":    0.0,
		"avg_processing_ms": 0,
	}

	return utils.SuccessResponse(c, metrics)
}
