// backend/internal/proj/marketplace/handler/translations.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
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

// UpdateTranslations обновляет переводы для объявления
func (h *TranslationsHandler) UpdateTranslations(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Проверяем существование объявления и права доступа
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", listingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем, является ли пользователь владельцем объявления
	isAdmin, _ := h.services.User().IsUserAdmin(c.Context(), "")
	if listing.UserID != userID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Вы не можете редактировать переводы этого объявления")
	}

	// Парсим данные запроса
	var updateData struct {
		Language     string            `json:"language"`
		Translations map[string]string `json:"translations"`
		IsVerified   bool              `json:"is_verified"`
		Provider     string            `json:"provider"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные запроса")
	}

	// Проверяем корректность языка
	if updateData.Language == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не указан язык перевода")
	}

	// Получаем провайдера из запроса или из параметра запроса
	provider := updateData.Provider
	if provider == "" {
		provider = c.Query("translation_provider", "google")
	}

	// Проверяем корректность провайдера и приводим к типу TranslationProvider
	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == "openai" {
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
		err := h.marketplaceService.UpdateTranslationWithProvider(c.Context(), translation, translationProvider)
		if err != nil {
			log.Printf("Failed to update translation for listing %d, field %s, language %s: %v",
				listingID, fieldName, updateData.Language, err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось обновить перевод")
		}
	}

	// После обновления переводов, переиндексируем объявление
	go func() {
		updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
		if err != nil {
			log.Printf("Failed to get updated listing for reindexing: %v", err)
			return
		}

		err = h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing)
		if err != nil {
			log.Printf("Failed to reindex listing after translation update: %v", err)
		}
	}()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Переводы успешно обновлены",
	})
}

// GetTranslations получает переводы для объявления
func (h *TranslationsHandler) GetTranslations(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
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
		log.Printf("Failed to get translations for listing %d: %v", listingID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить переводы")
	}

	// Фильтруем переводы по языку, если он указан
	var filteredTranslations []models.Translation
	for _, translation := range translations {
		if language == "" || translation.Language == language {
			filteredTranslations = append(filteredTranslations, translation)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    filteredTranslations,
	})
}

// TranslateText переводит текст с использованием выбранного провайдера
func (h *TranslationsHandler) TranslateText(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Text       string `json:"text"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
		Provider   string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные запроса")
	}

	// Проверяем обязательные поля
	if request.Text == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Текст для перевода не указан")
	}

	if request.TargetLang == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Целевой язык не указан")
	}

	// Если исходный язык не указан, определяем его автоматически
	if request.SourceLang == "" {
		detectedLang, _, err := h.services.Translation().DetectLanguage(c.Context(), request.Text)
		if err != nil {
			log.Printf("Failed to detect language: %v", err)
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
		log.Printf("Failed to translate text: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось перевести текст")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"translated_text": translatedText,
			"source_lang":     request.SourceLang,
			"target_lang":     request.TargetLang,
			"provider":        request.Provider,
		},
	})
}

// DetectLanguage определяет язык текста
func (h *TranslationsHandler) DetectLanguage(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Text string `json:"text"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные запроса")
	}

	// Проверяем обязательные поля
	if request.Text == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Текст для определения языка не указан")
	}

	// Определяем язык
	language, confidence, err := h.services.Translation().DetectLanguage(c.Context(), request.Text)
	if err != nil {
		log.Printf("Failed to detect language: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось определить язык")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"language":   language,
			"confidence": confidence,
		},
	})
}

// GetTranslationLimits возвращает лимиты использования сервиса перевода
func (h *TranslationsHandler) GetTranslationLimits(c *fiber.Ctx) error {
	// Пример реализации - обычно это получается от API сервиса перевода
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"daily_limit": 10000,
			"used_today":  3450,
			"remaining":   6550,
			"provider":    "google",
		},
	})
}

// SetTranslationProvider устанавливает провайдера перевода
func (h *TranslationsHandler) SetTranslationProvider(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		Provider string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные запроса")
	}

	// Проверяем, поддерживается ли запрошенный провайдер
	if request.Provider != "google" && request.Provider != "openai" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неподдерживаемый провайдер. Доступные варианты: google, openai")
	}

	// Здесь можно установить провайдер по умолчанию, например, через кеш или настройки пользователя
	// ...

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Провайдер перевода успешно установлен",
		"data": fiber.Map{
			"provider": request.Provider,
		},
	})
}

// BatchTranslateListings переводит несколько объявлений сразу
func (h *TranslationsHandler) BatchTranslateListings(c *fiber.Ctx) error {
	// Парсим данные запроса
	var request struct {
		ListingIDs []int  `json:"listing_ids"`
		TargetLang string `json:"target_lang"`
		Provider   string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные запроса")
	}

	// Проверяем обязательные поля
	if len(request.ListingIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не указаны ID объявлений для перевода")
	}

	if request.TargetLang == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не указан целевой язык")
	}

	// Определяем провайдер перевода
	provider := request.Provider
	if provider == "" {
		provider = "google"
	}

	translationProvider := service.GoogleTranslate
	if strings.ToLower(provider) == "openai" {
		translationProvider = service.OpenAI
	}

	// Запускаем процесс перевода в фоне
	go func() {
		for _, listingID := range request.ListingIDs {
			listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
			if err != nil {
				log.Printf("Failed to get listing %d for translation: %v", listingID, err)
				continue
			}

			if listing.OriginalLanguage == "" {
				// Определяем язык объявления
				detectedLang, _, err := h.services.Translation().DetectLanguage(c.Context(), listing.Title+" "+listing.Description)
				if err != nil {
					log.Printf("Failed to detect language for listing %d: %v", listingID, err)
					listing.OriginalLanguage = "en" // Используем английский по умолчанию
				} else {
					listing.OriginalLanguage = detectedLang
				}

				// Обновляем исходный язык в объявлении
				listing.OriginalLanguage = detectedLang
				err = h.marketplaceService.UpdateListing(c.Context(), listing)
				if err != nil {
					log.Printf("Failed to update listing %d with original language: %v", listingID, err)
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
					c.Context(),
					listing.Title,
					listing.OriginalLanguage,
					request.TargetLang,
					translationProvider,
				)
			} else {
				translatedTitle, err = h.services.Translation().Translate(
					c.Context(),
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
					err = h.marketplaceService.UpdateTranslationWithProvider(c.Context(), titleTranslation, translationProvider)
				} else {
					err = h.marketplaceService.UpdateTranslation(c.Context(), titleTranslation)
				}

				if err != nil {
					log.Printf("Failed to save title translation for listing %d to %s: %v", listing.ID, request.TargetLang, err)
				}
			}

			// Переводим описание
			var translatedDesc string
			if isFactory {
				translatedDesc, err = translationFactory.TranslateWithProvider(
					c.Context(),
					listing.Description,
					listing.OriginalLanguage,
					request.TargetLang,
					translationProvider,
				)
			} else {
				translatedDesc, err = h.services.Translation().Translate(
					c.Context(),
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
					err = h.marketplaceService.UpdateTranslationWithProvider(c.Context(), descTranslation, translationProvider)
				} else {
					err = h.marketplaceService.UpdateTranslation(c.Context(), descTranslation)
				}

				if err != nil {
					log.Printf("Failed to save description translation for listing %d to %s: %v", listing.ID, request.TargetLang, err)
				}
			}

			// Переиндексируем объявление с новыми переводами
			updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listing.ID)
			if err != nil {
				log.Printf("Failed to get updated listing %d for reindexing: %v", listing.ID, err)
				continue
			}

			err = h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing)
			if err != nil {
				log.Printf("Failed to reindex listing %d after translations: %v", listing.ID, err)
			}
		}
	}()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Процесс перевода запущен",
		"data": fiber.Map{
			"listing_count": len(request.ListingIDs),
			"target_lang":   request.TargetLang,
			"provider":      provider,
		},
	})
}
