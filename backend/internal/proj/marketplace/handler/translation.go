// backend/internal/proj/marketplace/handler/translations.go
package handler

import (
	"context"
	"strconv"
	"strings"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
)

const (
	translationProviderOpenAI = "openai"
)

// TranslationsHandler обрабатывает запросы, связанные с переводами
type TranslationsHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewTranslationsHandler создает новый обработчик переводов
func NewTranslationsHandler(services globalService.ServicesInterface) *TranslationsHandler {
	return &TranslationsHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// UpdateTranslations updates translations for a listing
// @Summary Update listing translations
// @Description Updates translations for a specific listing with support for different providers
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param translation_provider query string false "Translation provider (google, openai)" default(google)
// @Param translations body TranslationUpdateRequest true "Translation data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Translations updated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.updateTranslationError"
// @Security BearerAuth
// @Router /api/v1/marketplace/translations/{id} [put]
func (h *TranslationsHandler) UpdateTranslations(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Проверяем существование объявления и права доступа
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем, является ли пользователь владельцем объявления
	isAdmin, _ := h.services.User().IsUserAdmin(c.Context(), "")
	if listing.UserID != userID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Парсим данные запроса
	var updateData struct {
		Language     string            `json:"language"`
		Translations map[string]string `json:"translations"`
		IsVerified   bool              `json:"is_verified"`
		Provider     string            `json:"provider"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем корректность языка
	if updateData.Language == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.languageRequired")
	}

	// Получаем провайдера из запроса или из параметра запроса
	provider := updateData.Provider
	if provider == "" {
		provider = c.Query("translation_provider", "google")
	}

	// Проверяем корректность провайдера и приводим к типу TranslationProvider
	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == translationProviderOpenAI {
		translationProvider = service.OpenAI
	}

	// Обновляем каждый переведенный field
	for fieldName, translatedText := range updateData.Translations {
		// Проверяем поддерживаемые поля
		if fieldName != "title" && fieldName != "description" {
			continue // Пропускаем неподдерживаемые поля
		}

		translation := &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            updateData.Language,
			FieldName:           fieldName,
			TranslatedText:      translatedText,
			IsVerified:          updateData.IsVerified,
			IsMachineTranslated: false,
			Metadata:            map[string]interface{}{"provider": provider, "updated_by": userID},
		}

		// Обновляем перевод
		err := h.marketplaceService.UpdateTranslationWithProvider(c.Context(), translation, translationProvider, 0)
		if err != nil {
			logger.Error().Err(err).Int("listing_id", listingID).Str("field", fieldName).Str("language", updateData.Language).Msg("Failed to update translation")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateTranslationError")
		}
	}

	// После обновления переводов, переиндексируем объявление
	// Создаем новый контекст для фоновой задачи
	bgCtx := context.Background()
	go func() {
		// Используем bgCtx для предотвращения утечки контекста запроса
		updatedListing, err := h.marketplaceService.GetListingByID(bgCtx, listingID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get updated listing for reindexing")
			return
		}

		err = h.marketplaceService.Storage().IndexListing(bgCtx, updatedListing)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to reindex listing after translation update")
		}
	}()

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.translationsUpdated",
	})
}

// GetTranslations retrieves translations for a listing
// @Summary Get listing translations
// @Description Retrieves all translations for a specific listing, optionally filtered by language
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param language query string false "Language filter"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.Translation} "Translations retrieved successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.getTranslationsError"
// @Router /api/v1/marketplace/translations/{id} [get]
func (h *TranslationsHandler) GetTranslations(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем параметр языка из запроса
	language := c.Query("language")
	if language == "" {
		// Если язык не указан, пытаемся получить его из контекста
		if lang, ok := c.Locals("language").(string); ok {
			language = lang
		} else {
			language = "ru" // Язык по умолчанию
		}
	}

	// Получаем переводы
	translations, err := h.marketplaceService.Storage().GetTranslationsForEntity(c.Context(), "listing", listingID)
	if err != nil {
		logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to get translations")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getTranslationsError")
	}

	// Фильтруем переводы по языку, если он указан
	var filteredTranslations []models.Translation
	for _, translation := range translations {
		if language == "" || translation.Language == language {
			filteredTranslations = append(filteredTranslations, translation)
		}
	}

	return utils.SuccessResponse(c, filteredTranslations)
}

// TranslateText translates text using the selected provider
// @Summary Translate text
// @Description Translates text from source language to target language using specified provider
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param translation body TranslateTextRequest true "Translation request"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=TranslatedTextData} "Text translated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.translateError"
// @Router /api/v1/marketplace/translations/translate [post]
func (h *TranslationsHandler) TranslateText(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Text       string `json:"text"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
		Provider   string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if request.Text == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.textRequired")
	}

	if request.TargetLang == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.targetLanguageRequired")
	}

	// Если исходный язык не указан, определяем его автоматически
	if request.SourceLang == "" {
		detectedLang, _, err := h.services.Translation().DetectLanguage(c.Context(), request.Text)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to detect language")
			request.SourceLang = "en" // Используем английский по умолчанию
		} else {
			request.SourceLang = detectedLang
		}
	}

	// Определяем провайдера перевода
	translationProvider := service.GoogleTranslate
	if strings.ToLower(request.Provider) == "openai" {
		translationProvider = service.OpenAI
	}

	// Переводим текст
	var translatedText string
	var err error

	// Проверяем, поддерживает ли сервис фабрику переводчиков
	translationFactory, isFactory := h.services.Translation().(service.TranslationFactoryInterface)
	if isFactory {
		translatedText, err = translationFactory.TranslateWithProvider(
			c.Context(),
			request.Text,
			request.SourceLang,
			request.TargetLang,
			translationProvider,
		)
	} else {
		// Используем обычный интерфейс перевода
		translatedText, err = h.services.Translation().Translate(
			c.Context(),
			request.Text,
			request.SourceLang,
			request.TargetLang,
		)
	}

	if err != nil {
		logger.Error().Err(err).Msg("Failed to translate text")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.translateError")
	}

	return utils.SuccessResponse(c, TranslatedTextData{
		TranslatedText: translatedText,
		SourceLang:     request.SourceLang,
		TargetLang:     request.TargetLang,
		Provider:       request.Provider,
	})
}

// DetectLanguage detects the language of text
// @Summary Detect text language
// @Description Automatically detects the language of the provided text
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param detection body DetectLanguageRequest true "Text for language detection"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=DetectedLanguageData} "Language detected successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.detectLanguageError"
// @Router /api/v1/marketplace/translations/detect-language [post]
func (h *TranslationsHandler) DetectLanguage(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Text string `json:"text"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if request.Text == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.textForDetectionRequired")
	}

	// Определяем язык
	language, confidence, err := h.services.Translation().DetectLanguage(c.Context(), request.Text)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to detect language")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.detectLanguageError")
	}

	return utils.SuccessResponse(c, DetectedLanguageData{
		Language:   language,
		Confidence: confidence,
	})
}

// GetTranslationLimits returns translation service usage limits
// @Summary Get translation limits
// @Description Returns current usage limits and statistics for the translation service
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=TranslationLimitsData} "Translation limits retrieved successfully"
// @Router /api/v1/translation/limits [get]
func (h *TranslationsHandler) GetTranslationLimits(c *fiber.Ctx) error {
	// Пример реализации - обычно это получается от API сервиса перевода
	return utils.SuccessResponse(c, TranslationLimitsData{
		DailyLimit: 10000,
		UsedToday:  3450,
		Remaining:  6550,
		Provider:   "google",
	})
}

// SetTranslationProvider sets the translation provider
// @Summary Set translation provider
// @Description Sets the default translation provider for the user
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param provider body SetProviderRequest true "Provider configuration (google, openai)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=SetProviderResponse} "Translation provider set successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Router /api/v1/translation/provider [post]
func (h *TranslationsHandler) SetTranslationProvider(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Provider string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем, поддерживается ли запрошенный провайдер
	if request.Provider != "google" && request.Provider != "openai" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.unsupportedProvider")
	}

	// Здесь можно установить провайдер по умолчанию, например, через кеш или настройки пользователя
	// ...

	return utils.SuccessResponse(c, SetProviderResponse{
		Success: true,
		Message: "marketplace.providerSet",
		Data: ProviderData{
			Provider: request.Provider,
		},
	})
}

// BatchTranslateListings translates multiple listings at once
// @Summary Batch translate listings
// @Description Translates multiple listings to the specified target language using the selected provider
// @Tags marketplace-translations
// @Accept json
// @Produce json
// @Param batch body BatchTranslateRequest true "Batch translation request"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=BatchTranslateData} "Batch translation started successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Security BearerAuth
// @Router /api/v1/marketplace/translations/batch-translate [post]
func (h *TranslationsHandler) BatchTranslateListings(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		ListingIDs []int  `json:"listing_ids"`
		TargetLang string `json:"target_lang"`
		Provider   string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if len(request.ListingIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.listingIdsRequired")
	}

	if request.TargetLang == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.targetLanguageRequired")
	}

	// Определяем провайдер перевода
	provider := request.Provider
	if provider == "" {
		provider = "google"
	}

	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == translationProviderOpenAI {
		translationProvider = service.OpenAI
	}

	// Запускаем процесс перевода в фоне
	// Создаем новый контекст для фоновой задачи
	bgCtx := context.Background()
	go func() {
		// Используем bgCtx для предотвращения утечки контекста запроса
		for _, listingID := range request.ListingIDs {
			listing, err := h.marketplaceService.GetListingByID(bgCtx, listingID)
			if err != nil {
				logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to get listing for translation")
				continue
			}

			if listing.OriginalLanguage == "" {
				// Определяем язык объявления
				detectedLang, _, err := h.services.Translation().DetectLanguage(bgCtx, listing.Title+" "+listing.Description)
				if err != nil {
					logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to detect language for listing")
					listing.OriginalLanguage = "en" // Используем английский по умолчанию
				} else {
					listing.OriginalLanguage = detectedLang
				}

				// Обновляем исходный язык в объявлении
				listing.OriginalLanguage = detectedLang
				err = h.marketplaceService.UpdateListing(bgCtx, listing)
				if err != nil {
					logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to update listing with original language")
				}
			}

			// Пропускаем, если исходный язык совпадает с целевым
			if listing.OriginalLanguage == request.TargetLang {
				continue
			}

			// Переводим заголовок и описание
			translationFactory, isFactory := h.services.Translation().(service.TranslationFactoryInterface)

			// Переводим заголовок
			var translatedTitle string
			if isFactory {
				translatedTitle, err = translationFactory.TranslateWithProvider(
					bgCtx,
					listing.Title,
					listing.OriginalLanguage,
					request.TargetLang,
					translationProvider,
				)
			} else {
				translatedTitle, err = h.services.Translation().Translate(
					bgCtx,
					listing.Title,
					listing.OriginalLanguage,
					request.TargetLang,
				)
			}

			if err == nil {
				// Сохраняем перевод заголовка
				titleTranslation := &models.Translation{
					EntityType:          "listing",
					EntityID:            listing.ID,
					Language:            request.TargetLang,
					FieldName:           "title",
					TranslatedText:      translatedTitle,
					IsMachineTranslated: true,
					IsVerified:          false,
					Metadata:            map[string]interface{}{"provider": provider},
				}

				if isFactory {
					err = h.marketplaceService.UpdateTranslationWithProvider(bgCtx, titleTranslation, translationProvider, 0)
				} else {
					err = h.marketplaceService.UpdateTranslation(bgCtx, titleTranslation)
				}

				if err != nil {
					logger.Error().Err(err).Int("listing_id", listing.ID).Str("target_lang", request.TargetLang).Msg("Failed to save title translation")
				}
			}

			// Переводим описание
			var translatedDesc string
			if isFactory {
				translatedDesc, err = translationFactory.TranslateWithProvider(
					bgCtx,
					listing.Description,
					listing.OriginalLanguage,
					request.TargetLang,
					translationProvider,
				)
			} else {
				translatedDesc, err = h.services.Translation().Translate(
					bgCtx,
					listing.Description,
					listing.OriginalLanguage,
					request.TargetLang,
				)
			}

			if err == nil {
				// Сохраняем перевод описания
				descTranslation := &models.Translation{
					EntityType:          "listing",
					EntityID:            listing.ID,
					Language:            request.TargetLang,
					FieldName:           "description",
					TranslatedText:      translatedDesc,
					IsMachineTranslated: true,
					IsVerified:          false,
					Metadata:            map[string]interface{}{"provider": provider},
				}

				if isFactory {
					err = h.marketplaceService.UpdateTranslationWithProvider(bgCtx, descTranslation, translationProvider, 0)
				} else {
					err = h.marketplaceService.UpdateTranslation(bgCtx, descTranslation)
				}

				if err != nil {
					logger.Error().Err(err).Int("listing_id", listing.ID).Str("target_lang", request.TargetLang).Msg("Failed to save description translation")
				}
			}

			// Переиндексируем объявление с новыми переводами
			updatedListing, err := h.marketplaceService.GetListingByID(bgCtx, listing.ID)
			if err != nil {
				logger.Error().Err(err).Int("listing_id", listing.ID).Msg("Failed to get updated listing for reindexing")
				continue
			}

			err = h.marketplaceService.Storage().IndexListing(bgCtx, updatedListing)
			if err != nil {
				logger.Error().Err(err).Int("listing_id", listing.ID).Msg("Failed to reindex listing after translations")
			}
		}
	}()

	return utils.SuccessResponse(c, BatchTranslateData{
		ListingCount: len(request.ListingIDs),
		TargetLang:   request.TargetLang,
		Provider:     provider,
	})
}
