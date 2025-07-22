package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
)

const (
	providerGoogle = "google"
	providerOpenAI = "openai"
	providerManual = "manual"
	entityTypeCategory = "category"
	entityTypeAttribute = "attribute"
)

// AdminTranslationsHandler обрабатывает административные запросы для переводов
type AdminTranslationsHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewAdminTranslationsHandler создает новый обработчик для административных переводов
func NewAdminTranslationsHandler(services globalService.ServicesInterface) *AdminTranslationsHandler {
	return &AdminTranslationsHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// BatchTranslateCategories массово переводит категории
// @Summary Batch translate categories
// @Description Translates multiple categories to specified languages
// @Tags marketplace-admin-translations
// @Accept json
// @Produce json
// @Param body body BatchTranslateCategoriesRequest true "Batch translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=BatchTranslateResponse} "Batch translation started successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.batchTranslateError"
// @Security BearerAuth
// @Router /api/admin/translations/batch-categories [post]
func (h *AdminTranslationsHandler) BatchTranslateCategories(c *fiber.Ctx) error {
	var request struct {
		CategoryIDs     []int    `json:"category_ids"`
		TargetLanguages []string `json:"target_languages"`
		Provider        string   `json:"provider"`
		AutoTranslate   bool     `json:"auto_translate"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if len(request.CategoryIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.categoryIdsRequired")
	}

	if len(request.TargetLanguages) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.targetLanguagesRequired")
	}

	// Определяем провайдер перевода
	provider := request.Provider
	if provider == "" {
		provider = providerGoogle
	}

	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == providerOpenAI {
		translationProvider = service.OpenAI
	}

	// Создаем фоновый контекст для длительной операции
	bgCtx := context.Background()

	// Счетчики для статистики
	totalTranslations := 0
	successfulTranslations := 0

	// Запускаем массовый перевод в фоне
	go func() {
		for _, categoryID := range request.CategoryIDs {
			// Получаем категорию
			categories, err := h.marketplaceService.GetCategories(bgCtx)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get categories")
				continue
			}

			var category *models.MarketplaceCategory
			for i := range categories {
				if categories[i].ID == categoryID {
					category = &categories[i]
					break
				}
			}

			if category == nil {
				logger.Error().Int("category_id", categoryID).Msg("Category not found")
				continue
			}

			// Переводим на каждый целевой язык
			for _, targetLang := range request.TargetLanguages {
				totalTranslations++

				// Пропускаем, если текст уже на целевом языке
				if targetLang == "en" && isLikelyEnglish(category.Name) {
					continue
				}
				if targetLang == "ru" && isLikelyCyrillic(category.Name) {
					continue
				}

				// Переводим название категории
				translatedText, err := h.marketplaceService.TranslateText(bgCtx, category.Name, "auto", targetLang)
				if err != nil {
					logger.Error().Err(err).
						Int("category_id", categoryID).
						Str("target_lang", targetLang).
						Msg("Failed to translate category name")
					continue
				}

				// Сохраняем перевод
				translation := &models.Translation{
					EntityType:          "category",
					EntityID:            categoryID,
					Language:            targetLang,
					FieldName:           "name",
					TranslatedText:      translatedText,
					IsMachineTranslated: true,
					IsVerified:          false,
					Metadata:            map[string]interface{}{"provider": provider},
				}

				if err := h.marketplaceService.UpdateTranslationWithProvider(bgCtx, translation, translationProvider); err != nil {
					logger.Error().Err(err).
						Int("category_id", categoryID).
						Str("target_lang", targetLang).
						Msg("Failed to save category translation")
					continue
				}

				successfulTranslations++
			}
		}

		logger.Info().
			Int("total", totalTranslations).
			Int("successful", successfulTranslations).
			Msg("Batch category translation completed")
	}()

	return utils.SuccessResponse(c, BatchTranslateResponse{
		Success: true,
		Message: "marketplace.batchTranslationStarted",
		Data: BatchTranslateData{
			TotalCount:      len(request.CategoryIDs),
			TargetLanguages: request.TargetLanguages,
			Provider:        provider,
		},
	})
}

// BatchTranslateAttributes массово переводит атрибуты и их опции
// @Summary Batch translate attributes
// @Description Translates multiple attributes and their options to specified languages
// @Tags marketplace-admin-translations
// @Accept json
// @Produce json
// @Param body body BatchTranslateAttributesRequest true "Batch translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=BatchTranslateResponse} "Batch translation started successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.batchTranslateError"
// @Security BearerAuth
// @Router /api/admin/translations/batch-attributes [post]
func (h *AdminTranslationsHandler) BatchTranslateAttributes(c *fiber.Ctx) error {
	var request struct {
		AttributeIDs     []int    `json:"attribute_ids"`
		TargetLanguages  []string `json:"target_languages"`
		Provider         string   `json:"provider"`
		TranslateOptions bool     `json:"translate_options"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if len(request.AttributeIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.attributeIdsRequired")
	}

	if len(request.TargetLanguages) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.targetLanguagesRequired")
	}

	// Определяем провайдер перевода
	provider := request.Provider
	if provider == "" {
		provider = providerGoogle
	}

	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == providerOpenAI {
		translationProvider = service.OpenAI
	}

	// Создаем фоновый контекст
	bgCtx := context.Background()

	// Счетчики для статистики
	totalTranslations := 0
	successfulTranslations := 0

	// Запускаем массовый перевод в фоне
	go func() {
		for _, attributeID := range request.AttributeIDs {
			// Получаем атрибут
			attribute, err := h.marketplaceService.GetAttributeByID(bgCtx, attributeID)
			if err != nil {
				logger.Error().Err(err).Int("attribute_id", attributeID).Msg("Failed to get attribute")
				continue
			}

			if attribute == nil {
				logger.Error().Int("attribute_id", attributeID).Msg("Attribute not found")
				continue
			}

			// Переводим на каждый целевой язык
			for _, targetLang := range request.TargetLanguages {
				// Переводим display_name атрибута
				totalTranslations++

				// Пропускаем, если текст уже на целевом языке
				if targetLang == "en" && isLikelyEnglish(attribute.DisplayName) {
					continue
				}
				if targetLang == "ru" && isLikelyCyrillic(attribute.DisplayName) {
					continue
				}

				translatedName, err := h.marketplaceService.TranslateText(bgCtx, attribute.DisplayName, "auto", targetLang)
				if err != nil {
					logger.Error().Err(err).
						Int("attribute_id", attributeID).
						Str("target_lang", targetLang).
						Msg("Failed to translate attribute display name")
					continue
				}

				// Сохраняем перевод названия
				translation := &models.Translation{
					EntityType:          "attribute",
					EntityID:            attributeID,
					Language:            targetLang,
					FieldName:           "display_name",
					TranslatedText:      translatedName,
					IsMachineTranslated: true,
					IsVerified:          false,
					Metadata:            map[string]interface{}{"provider": provider},
				}

				if err := h.marketplaceService.UpdateTranslationWithProvider(bgCtx, translation, translationProvider); err != nil {
					logger.Error().Err(err).
						Int("attribute_id", attributeID).
						Str("target_lang", targetLang).
						Msg("Failed to save attribute translation")
					continue
				}

				successfulTranslations++

				// Переводим опции, если требуется
				if request.TranslateOptions && attribute.Options != nil {
					// Парсим опции
					var options []map[string]interface{}
					if err := json.Unmarshal(attribute.Options, &options); err == nil {
						for _, option := range options {
							if value, ok := option["value"].(string); ok {
								totalTranslations++

								// Пропускаем, если текст уже на целевом языке
								if targetLang == "en" && isLikelyEnglish(value) {
									continue
								}
								if targetLang == "ru" && isLikelyCyrillic(value) {
									continue
								}

								translatedOption, err := h.marketplaceService.TranslateText(bgCtx, value, "auto", targetLang)
								if err != nil {
									logger.Error().Err(err).
										Int("attribute_id", attributeID).
										Str("option", value).
										Str("target_lang", targetLang).
										Msg("Failed to translate option")
									continue
								}

								// Сохраняем перевод опции
								optionTranslation := &models.Translation{
									EntityType:          "attribute_option",
									EntityID:            attributeID,
									Language:            targetLang,
									FieldName:           value,
									TranslatedText:      translatedOption,
									IsMachineTranslated: true,
									IsVerified:          false,
									Metadata:            map[string]interface{}{"provider": provider},
								}

								if err := h.marketplaceService.UpdateTranslationWithProvider(bgCtx, optionTranslation, translationProvider); err != nil {
									logger.Error().Err(err).
										Int("attribute_id", attributeID).
										Str("option", value).
										Str("target_lang", targetLang).
										Msg("Failed to save option translation")
									continue
								}

								successfulTranslations++
							}
						}
					}
				}
			}
		}

		logger.Info().
			Int("total", totalTranslations).
			Int("successful", successfulTranslations).
			Msg("Batch attribute translation completed")
	}()

	return utils.SuccessResponse(c, BatchTranslateResponse{
		Success: true,
		Message: "marketplace.batchTranslationStarted",
		Data: BatchTranslateData{
			TotalCount:      len(request.AttributeIDs),
			TargetLanguages: request.TargetLanguages,
			Provider:        provider,
		},
	})
}

// GetTranslationStatus получает статус переводов для сущности
// @Summary Get translation status
// @Description Gets translation status for categories or attributes
// @Tags marketplace-admin-translations
// @Accept json
// @Produce json
// @Param entity_type query string true "Entity type (category, attribute)"
// @Param entity_ids query string false "Comma-separated entity IDs"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]TranslationStatusItem} "Translation status retrieved"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidEntityType"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getTranslationStatusError"
// @Security BearerAuth
// @Router /api/admin/translations/status [get]
func (h *AdminTranslationsHandler) GetTranslationStatus(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityIDs := c.Query("entity_ids")

	// Валидация типа сущности
	if entityType != entityTypeCategory && entityType != entityTypeAttribute {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidEntityType")
	}

	// Парсим ID сущностей, если предоставлены
	var ids []int
	if entityIDs != "" {
		idStrings := strings.Split(entityIDs, ",")
		for _, idStr := range idStrings {
			var id int
			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				ids = append(ids, id)
			}
		}
	}

	// Получаем статус переводов
	query := `
		SELECT DISTINCT entity_id, language, field_name, is_machine_translated, is_verified
		FROM translations
		WHERE entity_type = $1
	`
	args := []interface{}{entityType}

	if len(ids) > 0 {
		query += " AND entity_id = ANY($2)"
		args = append(args, ids)
	}

	rows, err := h.marketplaceService.Storage().Query(c.Context(), query, args...)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get translation status")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getTranslationStatusError")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	// Формируем результат
	statusMap := make(map[int]map[string]TranslationFieldStatus)

	for rows.Next() {
		var entityID int
		var language, fieldName string
		var isMachineTranslated, isVerified bool

		if err := rows.Scan(&entityID, &language, &fieldName, &isMachineTranslated, &isVerified); err != nil {
			continue
		}

		if statusMap[entityID] == nil {
			statusMap[entityID] = make(map[string]TranslationFieldStatus)
		}

		statusMap[entityID][language] = TranslationFieldStatus{
			Language:            language,
			IsTranslated:        true,
			IsMachineTranslated: isMachineTranslated,
			IsVerified:          isVerified,
		}
	}

	// Формируем финальный результат
	var result []TranslationStatusItem
	supportedLanguages := []string{"en", "ru", "sr"}

	// Если не указаны конкретные ID, получаем все сущности
	if len(ids) == 0 {
		if entityType == "category" {
			categories, err := h.marketplaceService.GetCategories(c.Context())
			if err == nil {
				for _, cat := range categories {
					ids = append(ids, cat.ID)
				}
			}
		}
		// Для атрибутов нужно получить все атрибуты через SQL
		// TODO: добавить метод GetAllAttributes в сервис
	}

	for _, id := range ids {
		item := TranslationStatusItem{
			EntityID:   id,
			EntityType: entityType,
			Languages:  make(map[string]TranslationFieldStatus),
		}

		for _, lang := range supportedLanguages {
			if status, exists := statusMap[id][lang]; exists {
				item.Languages[lang] = status
			} else {
				item.Languages[lang] = TranslationFieldStatus{
					Language:     lang,
					IsTranslated: false,
				}
			}
		}

		result = append(result, item)
	}

	return utils.SuccessResponse(c, result)
}

// UpdateFieldTranslation обновляет перевод для конкретного поля сущности
// @Summary Update field translation
// @Description Updates translation for a specific field of an entity
// @Tags marketplace-admin-translations
// @Accept json
// @Produce json
// @Param entity_type path string true "Entity type (category, attribute)"
// @Param entity_id path int true "Entity ID"
// @Param field_name path string true "Field name"
// @Param body body UpdateFieldTranslationRequest true "Translation update request"
// @Success 200 {object} utils.SuccessResponseSwag{data=TranslationFieldStatus} "Translation updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateTranslationError"
// @Security BearerAuth
// @Router /api/admin/translations/{entity_type}/{entity_id}/{field_name} [put]
func (h *AdminTranslationsHandler) UpdateFieldTranslation(c *fiber.Ctx) error {
	entityType := c.Params("entity_type")
	entityID, err := c.ParamsInt("entity_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidEntityId")
	}
	fieldName := c.Params("field_name")

	// Валидация типа сущности
	if entityType != entityTypeCategory && entityType != entityTypeAttribute {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidEntityType")
	}

	var request struct {
		Translations map[string]string `json:"translations"`
		Provider     string            `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Определяем провайдер перевода
	provider := request.Provider
	if provider == "" {
		provider = providerManual
	}

	translationProvider := service.Manual
	if strings.ToLower(provider) == "google" {
		translationProvider = service.GoogleTranslate
	} else if strings.ToLower(provider) == "openai" {
		translationProvider = service.OpenAI
	}

	// Обновляем переводы для каждого языка
	for lang, translatedText := range request.Translations {
		translation := &models.Translation{
			EntityType:          entityType,
			EntityID:            entityID,
			Language:            lang,
			FieldName:           fieldName,
			TranslatedText:      translatedText,
			IsMachineTranslated: provider != "manual",
			IsVerified:          provider == "manual",
			Metadata:            map[string]interface{}{"provider": provider},
		}

		if err := h.marketplaceService.UpdateTranslationWithProvider(c.Context(), translation, translationProvider); err != nil {
			logger.Error().Err(err).
				Str("entity_type", entityType).
				Int("entity_id", entityID).
				Str("field_name", fieldName).
				Str("language", lang).
				Msg("Failed to update translation")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateTranslationError")
		}
	}

	// Возвращаем обновленный статус
	updatedStatus := make(map[string]TranslationFieldStatus)
	for lang := range request.Translations {
		updatedStatus[lang] = TranslationFieldStatus{
			Language:            lang,
			IsTranslated:        true,
			IsMachineTranslated: provider != "manual",
			IsVerified:          provider == "manual",
		}
	}

	return utils.SuccessResponse(c, updatedStatus)
}

// Вспомогательные функции
func isLikelyEnglish(text string) bool {
	latinCount := 0
	totalLetters := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalLetters++
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				latinCount++
			}
		}
	}

	if totalLetters == 0 {
		return false
	}

	return float64(latinCount)/float64(totalLetters) > 0.8
}

func isLikelyCyrillic(text string) bool {
	cyrillicCount := 0
	totalLetters := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			totalLetters++
			if unicode.Is(unicode.Cyrillic, r) {
				cyrillicCount++
			}
		}
	}

	if totalLetters == 0 {
		return false
	}

	return float64(cyrillicCount)/float64(totalLetters) > 0.8
}

// Типы для Swagger документации
type BatchTranslateCategoriesRequest struct {
	CategoryIDs     []int    `json:"category_ids" example:"[1,2,3]"`
	TargetLanguages []string `json:"target_languages" example:"[\"en\",\"ru\",\"sr\"]"`
	Provider        string   `json:"provider" example:"google"`
	AutoTranslate   bool     `json:"auto_translate" example:"true"`
}

type BatchTranslateAttributesRequest struct {
	AttributeIDs     []int    `json:"attribute_ids" example:"[1,2,3]"`
	TargetLanguages  []string `json:"target_languages" example:"[\"en\",\"ru\",\"sr\"]"`
	Provider         string   `json:"provider" example:"google"`
	TranslateOptions bool     `json:"translate_options" example:"true"`
}

type TranslationFieldStatus struct {
	Language            string `json:"language"`
	IsTranslated        bool   `json:"is_translated"`
	IsMachineTranslated bool   `json:"is_machine_translated"`
	IsVerified          bool   `json:"is_verified"`
}

type TranslationStatusItem struct {
	EntityID   int                               `json:"entity_id"`
	EntityType string                            `json:"entity_type"`
	Languages  map[string]TranslationFieldStatus `json:"languages"`
}

type UpdateFieldTranslationRequest struct {
	Translations map[string]string `json:"translations" example:"{\"en\":\"Hello\",\"ru\":\"Привет\",\"sr\":\"Здраво\"}"`
	Provider     string            `json:"provider" example:"manual"`
}
