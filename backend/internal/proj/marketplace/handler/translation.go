// backend/internal/proj/marketplace/handler/translation.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/global/service"
	marketplaceService "backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type TranslationHandler struct {
	services service.ServicesInterface
}

func NewTranslationHandler(services service.ServicesInterface) *TranslationHandler {
	return &TranslationHandler{
		services: services,
	}
}

// GetTranslationLimits возвращает информацию о лимитах перевода для разных провайдеров
func (h *TranslationHandler) GetTranslationLimits(c *fiber.Ctx) error {
	// Проверяем авторизацию
	_, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Пытаемся получить доступ к фабрике перевода
	translationFactory, ok := h.services.Translation().(marketplaceService.TranslationFactoryInterface)
	if !ok {
		// Если интерфейс не реализован, возвращаем стандартные значения
		return utils.SuccessResponse(c, fiber.Map{
			"google": fiber.Map{
				"used":  0,
				"limit": 100,
			},
			"openai": fiber.Map{
				"used":  0,
				"limit": 0, // Нет жестких лимитов для OpenAI
			},
		})
	}

	// Получаем информацию о лимитах Google Translate
	googleUsed, googleLimit, err := translationFactory.GetTranslationCount(marketplaceService.GoogleTranslate)
	if err != nil {
		log.Printf("Ошибка получения лимитов Google Translate: %v", err)
		googleUsed = 0
		googleLimit = 100 // Значение по умолчанию
	}

	// Для OpenAI нет встроенных лимитов
	openaiUsed, openaiLimit := 0, 0

	// Получаем список доступных провайдеров
	availableProviders := translationFactory.GetAvailableProviders()
	providersStr := make([]string, len(availableProviders))
	for i, provider := range availableProviders {
		providersStr[i] = string(provider)
	}

	// Формируем ответ
	response := fiber.Map{
		"google": fiber.Map{
			"used":      googleUsed,
			"limit":     googleLimit,
			"available": googleLimit - googleUsed,
		},
		"openai": fiber.Map{
			"used":  openaiUsed,
			"limit": openaiLimit,
		},
		"available_providers": providersStr,
		"default_provider":    string(translationFactory.GetDefaultProvider()),
	}

	return utils.SuccessResponse(c, response)
}

// SetTranslationProvider устанавливает провайдер перевода по умолчанию
func (h *TranslationHandler) SetTranslationProvider(c *fiber.Ctx) error {
	// Проверяем авторизацию
	_, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	var request struct {
		Provider string `json:"provider"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат запроса")
	}

	// Пытаемся получить доступ к фабрике перевода
	translationFactory, ok := h.services.Translation().(marketplaceService.TranslationFactoryInterface)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Управление провайдером перевода недоступно")
	}

	// Преобразуем строку в enum
	var provider marketplaceService.TranslationProvider
	switch request.Provider {
	case "google":
		provider = marketplaceService.GoogleTranslate
	case "openai":
		provider = marketplaceService.OpenAI
	default:
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неизвестный провайдер перевода: "+request.Provider)
	}

	// Устанавливаем провайдер
	if err := translationFactory.SetDefaultProvider(provider); err != nil {
		log.Printf("Ошибка установки провайдера перевода: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка установки провайдера перевода: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"provider": string(provider),
		"message":  "Провайдер перевода успешно установлен",
	})
}

// BatchTranslateListings выполняет групповой перевод объявлений с оплатой
func (h *MarketplaceHandler) BatchTranslateListings(c *fiber.Ctx) error {
	// Проверяем авторизацию
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Парсим запрос
	var request struct {
		ListingIDs      []int    `json:"listing_ids"`
		TargetLanguages []string `json:"target_languages"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
	}

	if len(request.ListingIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No listing IDs provided")
	}

	if len(request.TargetLanguages) == 0 {
		// Если не указаны языки, используем стандартный набор
		request.TargetLanguages = []string{"en", "ru", "sr"}
	}

	// Стоимость перевода одного объявления
	const translationCostPerListing = 25.0

	// Общая стоимость
	totalCost := float64(len(request.ListingIDs)) * translationCostPerListing

	// Проверяем баланс пользователя
	balance, err := h.services.Balance().GetBalance(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting user balance: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user balance")
	}

	if balance.Balance < totalCost {
		return utils.ErrorResponse(c, fiber.StatusPaymentRequired, fmt.Sprintf(
			"Insufficient funds for translation. Required: %.2f RSD, Available: %.2f RSD",
			totalCost,
			balance.Balance,
		))
	}

	// Создаем транзакцию списания
	now := time.Now()
	transaction := &models.BalanceTransaction{
		UserID:        userID,
		Type:          "service_payment",
		Amount:        totalCost,
		Currency:      "RSD",
		Status:        "completed",
		PaymentMethod: "balance",
		Description:   fmt.Sprintf("Перевод %d объявлений", len(request.ListingIDs)),
		CreatedAt:     now,
		CompletedAt:   &now,
	}

	// Начинаем транзакцию в БД
	tx, err := h.services.Storage().BeginTx(c.Context(), nil)
	if err != nil {
		log.Printf("Error starting DB transaction: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	// Создаем транзакцию в БД
	transactionID, err := h.services.Storage().CreateTransaction(c.Context(), transaction)
	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create transaction")
	}

	// Обновляем баланс пользователя
	err = h.services.Storage().UpdateBalance(c.Context(), userID, -totalCost)
	if err != nil {
		log.Printf("Error updating balance: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update balance")
	}

	// Для отслеживания результатов
	results := make(map[int]map[string]string)
	successCount := 0
	failedCount := 0

	// Обрабатываем каждое объявление
	for _, listingID := range request.ListingIDs {
		// Проверяем, принадлежит ли объявление пользователю
		listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
		if err != nil {
			// Пропускаем объявления, которые не найдены или недоступны
			log.Printf("Error getting listing %d: %v", listingID, err)
			results[listingID] = map[string]string{"error": err.Error()}
			failedCount++
			continue
		}

		// Проверяем права доступа
		if listing.UserID != userID && listing.StorefrontID == nil {
			// Пропускаем объявления, которые не принадлежат пользователю
			results[listingID] = map[string]string{"error": "Access denied"}
			failedCount++
			continue
		}

		// Если объявление принадлежит витрине, проверяем права доступа к витрине
		if listing.StorefrontID != nil {
			storefront, err := h.services.Storefront().GetStorefrontByID(c.Context(), *listing.StorefrontID, userID)
			if err != nil || storefront.UserID != userID {
				results[listingID] = map[string]string{"error": "Access denied to storefront"}
				failedCount++
				continue
			}
		}

		// Источниковый язык
		sourceLanguage := listing.OriginalLanguage
		if sourceLanguage == "" {
			// Явное определение языка из текста, если язык не указан
			detectedLang, confidence, err := h.services.Translation().DetectLanguage(c.Context(), listing.Title+"\n"+listing.Description)
			if err == nil && confidence > 0.7 {
				sourceLanguage = detectedLang
				log.Printf("Detected language for listing %d: %s (confidence: %.2f)", listingID, detectedLang, confidence)

				// Обновляем язык в базе данных для будущих операций
				_, err = h.services.Storage().Exec(c.Context(),
					"UPDATE marketplace_listings SET original_language = $1 WHERE id = $2",
					detectedLang, listingID)
				if err != nil {
					log.Printf("Error updating original language: %v", err)
				}
			} else {
				// Получаем язык из заголовка запроса
				if userLang := c.Get("Accept-Language"); ok && userLang != "" {
					sourceLanguage = userLang
					log.Printf("Using user preferred language for listing %d: %s", listingID, userLang)
				} else {
					// Если всё еще нет языка, используем русский по умолчанию, так как большинство пользователей русскоговорящие
					sourceLanguage = "ru"
					log.Printf("Using default language (ru) for listing %d", listingID)
				}
			}

			// Обновляем поле original_language в базе данных
			_, err = h.services.Storage().Exec(c.Context(),
				"UPDATE marketplace_listings SET original_language = $1 WHERE id = $2",
				sourceLanguage, listingID)
			if err != nil {
				log.Printf("Error updating original language: %v", err)
			}
		}

		// Получаем провайдер перевода из запроса или используем Google по умолчанию
		provider := c.Query("translation_provider", "google")
		translationProvider := marketplaceService.GoogleTranslate
		if provider == "openai" {
			translationProvider = marketplaceService.OpenAI
		}

		log.Printf("Using translation provider: %s for batch translations", provider)

		// Переводим на каждый целевой язык
		var listingSuccess bool = false
		for _, targetLang := range request.TargetLanguages {
			// Пропускаем язык оригинала
			if targetLang == sourceLanguage {
				continue
			}

			// Проверяем, используем ли фабрику переводов
			translationFactory, isFactory := h.services.Translation().(marketplaceService.TranslationFactoryInterface)

			var translatedTitle, translatedDesc string
			var err error

			// Если фабрика доступна, используем её для перевода с указанием провайдера
			if isFactory {
				log.Printf("Using translation factory with provider %s for batch translation", provider)
				translatedTitle, err = translationFactory.TranslateWithProvider(c.Context(), listing.Title, sourceLanguage, targetLang, translationProvider)
			} else {
				translatedTitle, err = h.services.Translation().Translate(c.Context(), listing.Title, sourceLanguage, targetLang)
			}

			if err != nil {
				log.Printf("Error translating title for listing %d: %v", listingID, err)
				continue
			}

			// Переводим описание
			if isFactory {
				translatedDesc, err = translationFactory.TranslateWithProvider(c.Context(), listing.Description, sourceLanguage, targetLang, translationProvider)
			} else {
				translatedDesc, err = h.services.Translation().Translate(c.Context(), listing.Description, sourceLanguage, targetLang)
			}

			if err != nil {
				log.Printf("Error translating description for listing %d: %v", listingID, err)
				continue
			}

			// Сохраняем переводы
			titleTranslation := &models.Translation{
				EntityType:          "listing",
				EntityID:            listingID,
				Language:            targetLang,
				FieldName:           "title",
				TranslatedText:      translatedTitle,
				IsMachineTranslated: true,
				IsVerified:          false,
				Metadata:            map[string]interface{}{"provider": provider},
			}

			// Используем UpdateTranslationWithProvider если доступно
			if isFactory {
				err = h.marketplaceService.UpdateTranslationWithProvider(c.Context(), titleTranslation, translationProvider)
			} else {
				err = h.marketplaceService.UpdateTranslation(c.Context(), titleTranslation)
			}

			if err != nil {
				log.Printf("Error saving title translation for listing %d: %v", listingID, err)
			}

			descTranslation := &models.Translation{
				EntityType:          "listing",
				EntityID:            listingID,
				Language:            targetLang,
				FieldName:           "description",
				TranslatedText:      translatedDesc,
				IsMachineTranslated: true,
				IsVerified:          false,
				Metadata:            map[string]interface{}{"provider": provider},
			}

			// Используем UpdateTranslationWithProvider если доступно
			if isFactory {
				err = h.marketplaceService.UpdateTranslationWithProvider(c.Context(), descTranslation, translationProvider)
			} else {
				err = h.marketplaceService.UpdateTranslation(c.Context(), descTranslation)
			}

			if err != nil {
				log.Printf("Error saving description translation for listing %d: %v", listingID, err)
			}

			// Записываем результаты для ответа
			if results[listingID] == nil {
				results[listingID] = make(map[string]string)
			}
			results[listingID][targetLang] = "translated"
			listingSuccess = true
		}

		if listingSuccess {
			successCount++

			// Переиндексируем объявление после успешного перевода
			// Сначала получаем свежую версию объявления с переводами
			updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
			if err != nil {
				log.Printf("Warning: Failed to get listing %d for reindexing after translation: %v", listingID, err)
			} else {
				// Переиндексируем объявление в OpenSearch
				if err := h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing); err != nil {
					log.Printf("Warning: Failed to reindex listing %d after translation: %v", listingID, err)
				} else {
					log.Printf("Successfully reindexed listing %d after translation", listingID)
				}
			}
		} else {
			failedCount++
		}
	}

	// Фиксируем транзакцию в БД
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":        fmt.Sprintf("Successfully translated %d listings. Failed: %d", successCount, failedCount),
		"cost":           totalCost,
		"transaction_id": transactionID,
		"results":        results,
	})
}
