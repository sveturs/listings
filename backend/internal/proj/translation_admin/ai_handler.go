package translation_admin

import (
	"strings"

	"backend/internal/domain/models"
	"backend/internal/proj/translation_admin/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// AITranslationHandler handles AI translation endpoints
type AITranslationHandler struct {
	logger  zerolog.Logger
	service *Service
}

// NewAITranslationHandler creates new AI translation handler
func NewAITranslationHandler(logger zerolog.Logger, service *Service) *AITranslationHandler {
	return &AITranslationHandler{
		logger:  logger,
		service: service,
	}
}

// TranslateText godoc
// @Summary Translate text using AI
// @Description Translates text using configured AI providers
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body models.AITranslateRequest true "Translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.AITranslateResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/translate [post]
func (h *AITranslationHandler) TranslateText(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.AITranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Валидация запроса
	if req.Text == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.textRequired")
	}

	if req.TargetLang == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.targetLanguageRequired")
	}

	// SourceLang может быть пустым или "auto" для автоопределения

	// Получаем провайдера
	provider := service.GetAIProvider(req.Provider)
	if !provider.IsConfigured() {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.providerNotConfigured")
	}

	// Выполняем перевод
	translation, confidence, err := provider.Translate(ctx, req.Text, req.SourceLang, req.TargetLang)
	if err != nil {
		h.logger.Error().Err(err).
			Str("provider", req.Provider).
			Str("source_lang", req.SourceLang).
			Str("target_lang", req.TargetLang).
			Msg("Failed to translate text")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.translationFailed")
	}

	response := models.AITranslateResponse{
		Translation: translation,
		Confidence:  confidence,
		Provider:    provider.GetName(),
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationSuccess", response)
}

// TranslateBatch godoc
// @Summary Batch translate texts using AI
// @Description Translates multiple texts using configured AI providers
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body models.AITranslateBatchRequest true "Batch translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.AITranslateBatchResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/translate-batch [post]
func (h *AITranslationHandler) TranslateBatch(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.AITranslateBatchRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Валидация
	if len(req.Items) == 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.noItemsToTranslate")
	}

	if req.SourceLang == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.sourceLangRequired")
	}

	if len(req.TargetLangs) == 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.targetLangsRequired")
	}

	// Получаем провайдера
	provider := service.GetAIProvider(req.Provider)
	if !provider.IsConfigured() {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.providerNotConfigured")
	}

	// Переводим каждый элемент
	results := []models.AITranslationResult{}

	for _, item := range req.Items {
		result := models.AITranslationResult{
			Key:          item.Key,
			Module:       item.Module,
			Translations: make(map[string]string),
			Provider:     provider.GetName(),
		}

		// Переводим на каждый целевой язык
		for _, targetLang := range req.TargetLangs {
			if targetLang == req.SourceLang {
				// Пропускаем перевод на исходный язык
				result.Translations[targetLang] = item.Text
				continue
			}

			translation, confidence, err := provider.Translate(ctx, item.Text, req.SourceLang, targetLang)
			if err != nil {
				h.logger.Error().Err(err).
					Str("key", item.Key).
					Str("target_lang", targetLang).
					Msg("Failed to translate item")
				// Записываем ошибку, но продолжаем с другими переводами
				result.Translations[targetLang] = "[TRANSLATION_ERROR] " + item.Text
				result.Error = err.Error()
			} else {
				result.Translations[targetLang] = translation
				result.Confidence = confidence
			}
		}

		results = append(results, result)
	}

	response := models.AITranslateBatchResponse{
		Results:      results,
		TotalItems:   len(req.Items),
		SuccessCount: countSuccessful(results),
		FailedCount:  countFailed(results),
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.batchTranslationComplete", response)
}

// GetAvailableProviders godoc
// @Summary Get available AI translation providers
// @Description Returns list of configured AI translation providers
// @Tags Translation Admin
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]map[string]interface{}}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/providers [get]
func (h *AITranslationHandler) GetAvailableProviders(c *fiber.Ctx) error {
	providers := service.GetAvailableProviders()
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providersRetrieved", providers)
}

// TranslateModule godoc
// @Summary Translate entire module using AI
// @Description Translates all missing translations in a module
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body models.TranslateModuleRequest true "Module translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.TranslateModuleResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/translate-module [post]
func (h *AITranslationHandler) TranslateModule(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.TranslateModuleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Получаем все ключи модуля, которые нужно перевести
	moduleData, err := h.service.GetModuleTranslations(ctx, req.Module)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.moduleLoadError")
	}

	// Получаем провайдера
	provider := service.GetAIProvider(req.Provider)
	if !provider.IsConfigured() {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.providerNotConfigured")
	}

	results := []models.AITranslationResult{}

	// Переводим только те ключи, где есть плейсхолдеры или отсутствуют переводы
	for _, trans := range moduleData {
		translations := trans.Translations
		key := trans.Key
		needsTranslation := false
		sourceText := ""

		// Определяем исходный текст и проверяем необходимость перевода
		if val, ok := translations[req.SourceLang]; ok && !isPlaceholder(val) {
			sourceText = val
		} else {
			// Если нет исходного текста, пропускаем
			continue
		}

		// Проверяем нужен ли перевод для целевых языков
		for _, targetLang := range req.TargetLangs {
			if targetLang == req.SourceLang {
				continue
			}

			if val, ok := translations[targetLang]; !ok || isPlaceholder(val) {
				needsTranslation = true
				break
			}
		}

		if !needsTranslation {
			continue
		}

		// Переводим
		result := models.AITranslationResult{
			Key:          key,
			Module:       req.Module,
			Translations: make(map[string]string),
			Provider:     provider.GetName(),
		}

		for _, targetLang := range req.TargetLangs {
			if targetLang == req.SourceLang {
				result.Translations[targetLang] = sourceText
				continue
			}

			// Проверяем, нужен ли перевод
			if val, ok := translations[targetLang]; ok && !isPlaceholder(val) {
				result.Translations[targetLang] = val
				continue
			}

			// Выполняем перевод
			translation, confidence, err := provider.Translate(ctx, sourceText, req.SourceLang, targetLang)
			if err != nil {
				h.logger.Error().Err(err).
					Str("key", key).
					Str("target_lang", targetLang).
					Msg("Failed to translate key")
				result.Translations[targetLang] = "[ERROR] " + sourceText
				result.Error = err.Error()
			} else {
				result.Translations[targetLang] = translation
				result.Confidence = confidence
			}
		}

		results = append(results, result)
	}

	response := models.TranslateModuleResponse{
		Module:         req.Module,
		Results:        results,
		TotalKeys:      len(moduleData),
		TranslatedKeys: len(results),
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.moduleTranslationComplete", response)
}

// Helper functions
func countSuccessful(results []models.AITranslationResult) int {
	count := 0
	for _, r := range results {
		if r.Error == "" {
			count++
		}
	}
	return count
}

func countFailed(results []models.AITranslationResult) int {
	count := 0
	for _, r := range results {
		if r.Error != "" {
			count++
		}
	}
	return count
}

func isPlaceholder(text string) bool {
	// Проверяем, является ли текст плейсхолдером
	// Например: [RU] Text, [EN] Text, users.auth.error.* и т.д.
	if strings.HasPrefix(text, "[") && strings.Contains(text, "]") {
		return true
	}
	if strings.Contains(text, ".") && !strings.Contains(text, " ") {
		// Вероятно, это ключ типа "users.auth.error.invalid_token"
		return true
	}
	return false
}
