// backend/internal/handlers/marketplace.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MarketplaceHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

func NewMarketplaceHandler(services globalService.ServicesInterface) *MarketplaceHandler {
	return &MarketplaceHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

func (h *MarketplaceHandler) CreateListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные")
	}

	// Устанавливаем ID пользователя
	listing.UserID = userID

	// Проверяем, указаны ли город и страна в объявлении
	// Если нет, попробуем получить их из профиля пользователя
	if (listing.City == "" || listing.Country == "") && (listing.Latitude == nil || listing.Longitude == nil) {
		userProfile, err := h.services.User().GetUserProfile(c.Context(), userID)
		if err == nil && userProfile != nil {
			// Устанавливаем город из профиля, если не указан в объявлении
			if listing.City == "" && userProfile.City != "" {
				listing.City = userProfile.City
				log.Printf("Using city from user profile: %s", userProfile.City)
			}

			// Устанавливаем страну из профиля, если не указана в объявлении
			if listing.Country == "" && userProfile.Country != "" {
				listing.Country = userProfile.Country
				log.Printf("Using country from user profile: %s", userProfile.Country)
			}
		} else {
			log.Printf("Could not get user profile or profile has no location info: %v", err)
		}
	}

	// Парсим атрибуты из запроса (оставляем как есть)
	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		processAttributesFromRequest(requestBody, &listing)
	}
	if listing.StorefrontID == nil {
		// Получаем провайдер перевода из запроса или используем Google по умолчанию
		provider := c.Query("translation_provider", "google")
		// Здесь не используем translationProvider, так как нам нужно только решить использовать модерацию или нет

		log.Printf("Using translation provider: %s for new listing", provider)

		fullText := listing.Title + "\n" + listing.Description
		detectedLanguage, confidence, err := h.services.Translation().DetectLanguage(c.Context(), fullText)
		if err != nil || confidence < 0.8 {
			// Если не удалось определить язык или уверенность низкая
			if listing.OriginalLanguage != "" {
				detectedLanguage = listing.OriginalLanguage
			} else if userLang, ok := c.Locals("language").(string); ok && userLang != "" {
				detectedLanguage = userLang
			} else {
				detectedLanguage = "en" // Язык по умолчанию
			}
		}
		listing.OriginalLanguage = detectedLanguage

		// Модерируем текст только если используем OpenAI
		if provider == "openai" {
			// Модерируем текст перед сохранением
			log.Printf("Using OpenAI for moderation")
			translationFactory, isFactory := h.services.Translation().(service.TranslationFactoryInterface)

			// Модерация заголовка
			if isFactory {
				openAIService, err := translationFactory.GetTranslationService(service.OpenAI)
				if err == nil {
					moderatedTitle, err := openAIService.ModerateText(c.Context(), listing.Title, listing.OriginalLanguage)
					if err == nil {
						listing.Title = moderatedTitle
					} else {
						log.Printf("Failed to moderate title with OpenAI: %v", err)
					}
				}
			} else {
				// Для совместимости используем обычный сервис
				moderatedTitle, err := h.services.Translation().ModerateText(c.Context(), listing.Title, listing.OriginalLanguage)
				if err != nil {
					log.Printf("Failed to moderate title: %v", err)
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка модерации заголовка")
				}
				listing.Title = moderatedTitle
			}

			// Модерация описания
			if isFactory {
				openAIService, err := translationFactory.GetTranslationService(service.OpenAI)
				if err == nil {
					moderatedDesc, err := openAIService.ModerateText(c.Context(), listing.Description, listing.OriginalLanguage)
					if err == nil {
						listing.Description = moderatedDesc
					} else {
						log.Printf("Failed to moderate description with OpenAI: %v", err)
					}
				}
			} else {
				moderatedDesc, err := h.services.Translation().ModerateText(c.Context(), listing.Description, listing.OriginalLanguage)
				if err != nil {
					log.Printf("Failed to moderate description: %v", err)
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка модерации описания")
				}
				listing.Description = moderatedDesc
			}
		} else {
			log.Printf("Using Google Translate - skipping moderation")
		}
	}
	// Создаем объявление
	listingID, err := h.marketplaceService.CreateListing(c.Context(), &listing)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при создании объявления")
	}

	// Сохраняем переводы только для обычных пользователей (не витрин)
	if listing.StorefrontID == nil {
		// Сохраняем оригинальный текст как перевод для исходного языка
		err = h.marketplaceService.UpdateTranslation(c.Context(), &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            listing.OriginalLanguage,
			FieldName:           "title",
			TranslatedText:      listing.Title,
			IsMachineTranslated: false,
			IsVerified:          true,
		})
		if err != nil {
			log.Printf("Error saving original title translation: %v", err)
		}

		err = h.marketplaceService.UpdateTranslation(c.Context(), &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            listing.OriginalLanguage,
			FieldName:           "description",
			TranslatedText:      listing.Description,
			IsMachineTranslated: false,
			IsVerified:          true,
		})
		if err != nil {
			log.Printf("Error saving original description translation: %v", err)
		}

		// Получаем провайдер перевода из запроса или используем Google по умолчанию
		provider := c.Query("translation_provider", "google")
		translationProvider := service.GoogleTranslate
		if provider == "openai" {
			translationProvider = service.OpenAI
		}

		log.Printf("Using translation provider: %s for translations", provider)

		// Переводим на другие языки
		targetLanguages := []string{"en", "ru", "sr"}
		for _, targetLang := range targetLanguages {
			if targetLang == listing.OriginalLanguage {
				continue
			}

			// Проверяем, используем ли фабрику переводов
			translationFactory, isFactory := h.services.Translation().(service.TranslationFactoryInterface)

			var translatedTitle, translatedDesc string

			// Переводим заголовок
			if isFactory {
				log.Printf("Using translation factory with provider %s", provider)
				translatedTitle, err = translationFactory.TranslateWithProvider(c.Context(), listing.Title, listing.OriginalLanguage, targetLang, translationProvider)
			} else {
				log.Printf("Using default translation service")
				translatedTitle, err = h.services.Translation().Translate(c.Context(), listing.Title, listing.OriginalLanguage, targetLang)
			}

			if err == nil {
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
					log.Printf("Error saving %s title translation: %v", targetLang, err)
				}
			} else {
				log.Printf("Error translating title to %s with provider %s: %v", targetLang, provider, err)
			}

			// Переводим описание
			if isFactory {
				translatedDesc, err = translationFactory.TranslateWithProvider(c.Context(), listing.Description, listing.OriginalLanguage, targetLang, translationProvider)
			} else {
				translatedDesc, err = h.services.Translation().Translate(c.Context(), listing.Description, listing.OriginalLanguage, targetLang)
			}

			if err == nil {
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
					log.Printf("Error saving %s description translation: %v", targetLang, err)
				}
			} else {
				log.Printf("Error translating description to %s with provider %s: %v", targetLang, provider, err)
			}
		}
	}

	// После выполнения всех переводов делаем небольшую задержку
	// чтобы гарантировать, что все транзакции с переводами завершились
	time.Sleep(500 * time.Millisecond)

	// Получаем обновленное объявление со всеми переводами для переиндексации
	updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Warning: Failed to get listing %d for reindexing after translations: %v", listingID, err)
	} else {
		// Дополнительная проверка наличия переводов перед индексацией
		if updatedListing.Translations == nil || len(updatedListing.Translations) == 0 {
			log.Printf("Warning: Listing %d has no translations before reindexing, will try to load them explicitly", listingID)

			// Пытаемся явно загрузить переводы
			translations, err := h.marketplaceService.Storage().GetTranslationsForEntity(c.Context(), "listing", listingID)
			if err != nil {
				log.Printf("Error loading translations for listing %d: %v", listingID, err)
			} else if len(translations) > 0 {
				// Организуем переводы в структуру TranslationMap
				transMap := make(models.TranslationMap)
				for _, t := range translations {
					if _, ok := transMap[t.Language]; !ok {
						transMap[t.Language] = make(map[string]string)
					}
					transMap[t.Language][t.FieldName] = t.TranslatedText
				}
				updatedListing.Translations = transMap
				log.Printf("Explicitly loaded %d translations for listing %d", len(translations), listingID)
			}
		}

		// Переиндексируем объявление в OpenSearch, чтобы поиск работал по всем языкам
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing); err != nil {
			log.Printf("Warning: Failed to reindex listing %d after translations: %v", listingID, err)
		} else {
			log.Printf("Successfully reindexed listing %d with all translations", listingID)
		}
	}

	// Обновляем материализованное представление
	if err := h.marketplaceService.RefreshCategoryListingCounts(c.Context()); err != nil {
		log.Printf("Error refreshing category counts: %v", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      listingID,
		"message": "Объявление успешно создано",
	})
}

var (
	categoryTreeCache      []models.CategoryTreeNode
	categoryTreeLastUpdate time.Time
	categoryTreeMutex      sync.RWMutex
)

// ModerateImage проверяет изображение на запрещенный контент
func (h *MarketplaceHandler) ModerateImage(c *fiber.Ctx) error {
	log.Printf("ModerateImage: начало обработки запроса")

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("ModerateImage: ошибка получения файла: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файла")
	}

	log.Printf("ModerateImage: получен файл: name=%s, size=%d, contentType=%s",
		file.Filename, file.Size, file.Header.Get("Content-Type"))

	// Временно сохраняем файл
	tempDir := "./uploads/temp"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		log.Printf("ModerateImage: ошибка создания директории %s: %v", tempDir, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка создания временной директории")
	}

	tempPath := fmt.Sprintf("%s/%d%s", tempDir, time.Now().UnixNano(), filepath.Ext(file.Filename))
	if err := c.SaveFile(file, tempPath); err != nil {
		log.Printf("ModerateImage: ошибка сохранения файла во временную директорию %s: %v", tempPath, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка сохранения файла")
	}

	log.Printf("ModerateImage: файл временно сохранен в %s", tempPath)

	defer func() {
		if err := os.Remove(tempPath); err != nil {
			log.Printf("ModerateImage: ошибка удаления временного файла %s: %v", tempPath, err)
		} else {
			log.Printf("ModerateImage: временный файл %s успешно удален", tempPath)
		}
	}()

	// Получаем сервис Cloudinary
	log.Printf("ModerateImage: создание сервиса Cloudinary")
	cloudinaryService, err := service.NewCloudinaryService()
	if err != nil {
		log.Printf("ModerateImage: ошибка инициализации сервиса Cloudinary: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка инициализации сервиса обработки изображений")
	}

	log.Printf("ModerateImage: сервис Cloudinary успешно создан, начинаем модерацию изображения")

	// Проверяем изображение на модерацию
	result, err := cloudinaryService.ModerateImage(c.Context(), tempPath)
	if err != nil {
		log.Printf("ModerateImage: ошибка модерации изображения: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка модерации изображения")
	}

	log.Printf("ModerateImage: модерация успешно завершена, результат: %+v", result)

	return utils.SuccessResponse(c, result)
}

// EnhancePreview создает предпросмотр улучшенного изображения
func (h *MarketplaceHandler) EnhancePreview(c *fiber.Ctx) error {

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файла")
	}

	// Временно сохраняем файл
	tempPath := fmt.Sprintf("./uploads/temp/%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
	if err := c.SaveFile(file, tempPath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка сохранения файла")
	}
	defer os.Remove(tempPath) // Удаляем временный файл после обработки

	// Получаем сервис Cloudinary
	cloudinaryService, err := service.NewCloudinaryService()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка инициализации сервиса обработки изображений")
	}

	// Создаем предпросмотр улучшенного изображения
	result, err := cloudinaryService.EnhancePreview(c.Context(), tempPath)
	if err != nil {
		log.Printf("Error creating enhancement preview: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка создания предпросмотра")
	}

	// Добавляем цену за улучшение
	result["price"] = 30.0 // Цена в динарах (регулируйте в зависимости от стоимости)

	return utils.SuccessResponse(c, result)
}

// EnhanceImages улучшает все загруженные изображения
func (h *MarketplaceHandler) EnhanceImages(c *fiber.Ctx) error {
	// Авторизованный пользователь
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем баланс пользователя
	balance, err := h.services.Balance().GetBalance(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка проверки баланса")
	}

	// Получаем форму с файлами
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файлов")
	}

	files := form.File["images"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Нет файлов для обработки")
	}

	// Рассчитываем общую стоимость улучшения
	pricePerImage := 30.0 // RSD
	totalPrice := pricePerImage * float64(len(files))

	// Проверяем достаточность средств
	if balance.Balance < totalPrice {
		return utils.ErrorResponse(c, fiber.StatusPaymentRequired,
			fmt.Sprintf("Недостаточно средств. Требуется: %.2f RSD", totalPrice))
	}

	// Получаем сервис Cloudinary
	cloudinaryService, err := service.NewCloudinaryService()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка инициализации сервиса обработки изображений")
	}

	// Обрабатываем все изображения
	enhancedImages := make([]map[string]interface{}, 0, len(files))

	for _, file := range files {
		// Сохраняем во временную директорию
		tempPath := fmt.Sprintf("./uploads/temp/%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
		if err := c.SaveFile(file, tempPath); err != nil {
			continue // Пропускаем файл при ошибке
		}

		// Улучшаем изображение
		result, err := cloudinaryService.EnhanceImage(c.Context(), tempPath)
		if err != nil {
			os.Remove(tempPath)
			continue
		}

		// Добавляем информацию о файле
		result["original_name"] = file.Filename

		// Добавляем в результат
		enhancedImages = append(enhancedImages, result)

		// Удаляем временный файл
		os.Remove(tempPath)
	}

	// Если улучшено хотя бы одно изображение, списываем средства
	if len(enhancedImages) > 0 {
		// Пересчитываем стоимость на основе фактически улучшенных изображений
		actualPrice := pricePerImage * float64(len(enhancedImages))

		// Списываем средства с баланса
		txDescription := fmt.Sprintf("Улучшение %d изображений", len(enhancedImages))
		now := time.Now()
		transaction := &models.BalanceTransaction{
			UserID:        userID,
			Type:          "service_payment",
			Amount:        -actualPrice, // отрицательное значение для списания
			Currency:      "RSD",
			Status:        "completed",
			PaymentMethod: "balance",
			Description:   txDescription,
			CreatedAt:     now,
			CompletedAt:   &now,
		}

		txID, err := h.services.Storage().CreateTransaction(c.Context(), transaction)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка списания средств")
		}

		// Обновляем баланс напрямую через SQL
		_, err = h.services.Storage().Exec(c.Context(), `
    UPDATE user_balances 
    SET balance = balance - $1, 
        updated_at = NOW()
    WHERE user_id = $2
`, actualPrice, userID)
		if err != nil {
			log.Printf("Error updating balance: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обновления баланса")
		}

		return utils.SuccessResponse(c, fiber.Map{
			"success":         true,
			"enhanced_images": enhancedImages,
			"transaction_id":  txID,
			"total_price":     actualPrice,
		})
	}

	return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось улучшить изображения")
}

// GetSimilarListings возвращает похожие объявления
func (h *MarketplaceHandler) GetSimilarListings(c *fiber.Ctx) error {
	listingID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем параметры пагинации
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 8)
	offset := (page - 1) * limit

	// Получаем исходное объявление
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Ошибка при получении объявления: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить объявление")
	}

	// Формируем запрос для OpenSearch напрямую
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"category_id": listing.CategoryID,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"must_not": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"id": listingID, // Исключаем текущее объявление
						},
					},
				},
				"should":               []interface{}{},
				"minimum_should_match": 1,
			},
		},
		"from": offset,
		"size": limit,
	}

	// Добавляем поиск по заголовку и описанию
	shouldClauses := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]interface{})

	if listing.Title != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query":     listing.Title,
					"boost":     3.0,
					"fuzziness": "AUTO",
				},
			},
		})
	}

	if len(listing.Description) > 0 {
		descriptionExcerpt := listing.Description
		if len(descriptionExcerpt) > 200 {
			descriptionExcerpt = descriptionExcerpt[:200]
		}

		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"description": map[string]interface{}{
					"query":     descriptionExcerpt,
					"boost":     1.0,
					"fuzziness": "AUTO",
				},
			},
		})
	}

	// Добавляем ценовой диапазон
	if listing.Price > 0 {
		minPrice := listing.Price * 0.7
		maxPrice := listing.Price * 1.3

		shouldClauses = append(shouldClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte":   minPrice,
					"lte":   maxPrice,
					"boost": 1.5,
				},
			},
		})
	}

	// Обновляем запрос с новыми should-условиями
	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = shouldClauses

	// Логируем запрос для отладки
	queryJSON, _ := json.MarshalIndent(query, "", "  ")
	log.Printf("Запрос для похожих объявлений: %s", string(queryJSON))

	// Выполняем запрос напрямую через OpenSearch клиент
	osClient, err := h.getOpenSearchClient()

	if err != nil {
		log.Printf("Не удалось получить клиент OpenSearch: %v", err)
		// Используем запасной метод через стандартный поиск
		fallbackParams := &search.ServiceParams{
			CategoryID: strconv.Itoa(listing.CategoryID),
			Size:       limit,
			Page:       page,
			Sort:       "date_desc",
		}

		fallbackResult, fallbackErr := h.marketplaceService.SearchListingsAdvanced(c.Context(), fallbackParams)
		if fallbackErr != nil {
			log.Printf("Ошибка резервного поиска: %v", fallbackErr)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось найти похожие объявления")
		}

		// Фильтруем результаты, убирая исходное объявление
		var similarListings []*models.MarketplaceListing
		for _, item := range fallbackResult.Items {
			if item.ID != listingID {
				similarListings = append(similarListings, item)
				if len(similarListings) >= limit {
					break
				}
			}
		}

		return utils.SuccessResponse(c, similarListings)
	}

	// Выполняем поиск напрямую
	indexName := "marketplace_listings" // или получите название индекса из конфигурации
	queryJSONBytes, _ := json.Marshal(query)
	responseBytes, err := osClient.Execute("POST", "/"+indexName+"/_search", queryJSONBytes)

	if err != nil {
		log.Printf("Ошибка выполнения запроса к OpenSearch: %v", err)
		// Используем запасной метод
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка поиска")
	}

	// Декодируем ответ
	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		log.Printf("Ошибка декодирования ответа: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка декодирования результатов")
	}

	// Извлекаем id похожих объявлений
	var similarIDs []int
	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if id, ok := hitMap["_id"].(string); ok {
						if idInt, err := strconv.Atoi(id); err == nil {
							similarIDs = append(similarIDs, idInt)
						}
					}
				}
			}
		}
	}

	// Загружаем полные данные объявлений по ID
	var similarListings []*models.MarketplaceListing
	for _, id := range similarIDs {
		if listing, err := h.marketplaceService.GetListingByID(c.Context(), id); err == nil {
			similarListings = append(similarListings, listing)
		}
	}

	log.Printf("Найдено %d похожих объявлений", len(similarListings))
	return utils.SuccessResponse(c, similarListings)
}

// Вспомогательный метод для получения клиента OpenSearch
func (h *MarketplaceHandler) getOpenSearchClient() (interface {
	Execute(method, path string, body []byte) ([]byte, error)
}, error) {
	// Получаем доступ к хранилищу через marketplaceService
	storage := h.marketplaceService.Storage()

	// Пробуем получить клиент через метод GetOpenSearchClient, если он есть
	if db, ok := storage.(interface {
		GetOpenSearchClient() (interface {
			Execute(method, path string, body []byte) ([]byte, error)
		}, error)
	}); ok {
		return db.GetOpenSearchClient()
	}

	// Если метод GetOpenSearchClient отсутствует, возвращаем ошибку
	return nil, fmt.Errorf("OpenSearch клиент не доступен")
}

// isSignificantAttribute определяет, является ли атрибут значимым для поиска похожих объявлений
func isSignificantAttribute(attrName string) bool {
	// Список значимых атрибутов, влияющих на определение похожести
	significantAttrs := map[string]bool{
		"make":            true,
		"model":           true,
		"brand":           true,
		"year":            true,
		"manufacturer":    true,
		"type":            true,
		"category":        true,
		"rooms":           true,
		"property_type":   true,
		"body_type":       true,
		"engine_capacity": true,
		"processor":       true,
		"screen_size":     true,
		"memory":          true,
		"ram":             true,
		"os":              true,
		"color":           true,
		"material":        true,
		"size":            true,
	}

	return significantAttrs[attrName]
}

// Обновим функцию processAttributesFromRequest для лучшей обработки атрибутов недвижимости
func processAttributesFromRequest(requestBody map[string]interface{}, listing *models.MarketplaceListing) {
	if attributesRaw, ok := requestBody["attributes"].([]interface{}); ok {
		log.Printf("DEBUG: Found %d attributes in request", len(attributesRaw))

		for _, attrRaw := range attributesRaw {
			if attrMap, ok := attrRaw.(map[string]interface{}); ok {
				var attr models.ListingAttributeValue

				// ID атрибута
				if attrID, ok := attrMap["attribute_id"].(float64); ok {
					attr.AttributeID = int(attrID)
				}

				// Имя атрибута для отладки
				if attrName, ok := attrMap["attribute_name"].(string); ok {
					attr.AttributeName = attrName
				}

				// Имя для отображения
				if displayName, ok := attrMap["display_name"].(string); ok {
					attr.DisplayName = displayName
				}

				// Тип атрибута
				if attrType, ok := attrMap["attribute_type"].(string); ok {
					attr.AttributeType = attrType

					// Список числовых атрибутов
					isNumericAttr := attr.AttributeName == "rooms" ||
						attr.AttributeName == "floor" ||
						attr.AttributeName == "total_floors" ||
						attr.AttributeName == "area" ||
						attr.AttributeName == "land_area" ||
						attr.AttributeName == "mileage" ||
						attr.AttributeName == "year" ||
						attr.AttributeName == "engine_capacity" ||
						attr.AttributeName == "power" ||
						attr.AttributeName == "screen_size"

					if isNumericAttr {
						// Обрабатываем как числовой атрибут, независимо от указанного типа
						var numValue float64
						var isSet bool

						if value, ok := attrMap["value"].(float64); ok {
							numValue = value
							isSet = true
						} else if value, ok := attrMap["numeric_value"].(float64); ok {
							numValue = value
							isSet = true
						} else if strValue, ok := attrMap["value"].(string); ok && strValue != "" {
							// Удаляем все, кроме цифр, точки и минуса
							clean := regexp.MustCompile(`[^\d\.-]`).ReplaceAllString(strValue, "")
							if parsedValue, parseErr := strconv.ParseFloat(clean, 64); parseErr == nil {
								numValue = parsedValue
								isSet = true
							}
						}

						if isSet {
							attr.NumericValue = &numValue
							attr.AttributeType = "number" // Принудительно устанавливаем тип как числовой

							// Устанавливаем единицу измерения
							switch attr.AttributeName {
							case "area":
								attr.Unit = "m²"
								attr.DisplayValue = fmt.Sprintf("%g м²", numValue)
							case "land_area":
								attr.Unit = "ar"
								attr.DisplayValue = fmt.Sprintf("%g сот", numValue)
							case "mileage":
								attr.Unit = "km"
								attr.DisplayValue = fmt.Sprintf("%g км", numValue)
							case "engine_capacity":
								attr.Unit = "l"
								attr.DisplayValue = fmt.Sprintf("%g л", numValue)
							case "power":
								attr.Unit = "ks"
								attr.DisplayValue = fmt.Sprintf("%g л.с.", numValue)
							case "screen_size":
								attr.Unit = "inč"
								attr.DisplayValue = fmt.Sprintf("%g\"", numValue)
							case "rooms":
								attr.Unit = "soba"
								attr.DisplayValue = fmt.Sprintf("%g", numValue)
							case "floor", "total_floors":
								attr.Unit = "sprat"
								attr.DisplayValue = fmt.Sprintf("%g", numValue)
							case "year":
								attr.DisplayValue = fmt.Sprintf("%d", int(numValue))
							default:
								attr.DisplayValue = fmt.Sprintf("%g", numValue)
							}
						}
					} else {
						// Обычная обработка для других атрибутов
						if attr.AttributeType == "text" || attr.AttributeType == "select" {
							if value, ok := attrMap["value"].(string); ok && value != "" {
								strValue := value
								attr.TextValue = &strValue
								attr.DisplayValue = strValue
							} else if value, ok := attrMap["text_value"].(string); ok && value != "" {
								strValue := value
								attr.TextValue = &strValue
								attr.DisplayValue = strValue
							}
						} else if attr.AttributeType == "boolean" {
							if value, ok := attrMap["value"].(bool); ok {
								boolValue := value
								attr.BooleanValue = &boolValue

								if boolValue {
									attr.DisplayValue = "Да"
								} else {
									attr.DisplayValue = "Нет"
								}
							} else if strValue, ok := attrMap["value"].(string); ok {
								boolValue := strValue == "true" || strValue == "да" || strValue == "1"
								attr.BooleanValue = &boolValue

								if boolValue {
									attr.DisplayValue = "Да"
								} else {
									attr.DisplayValue = "Нет"
								}
							}
						}
					}
				}

				// Проверка есть ли единица измерения в запросе
				if unit, ok := attrMap["unit"].(string); ok && unit != "" {
					attr.Unit = unit
				}

				// Добавляем атрибут только если есть какое-то значение
				if attr.TextValue != nil || attr.NumericValue != nil || attr.BooleanValue != nil || attr.JSONValue != nil {
					listing.Attributes = append(listing.Attributes, attr)
				}
			}
		}
	}
}

// GetAttributeRanges возвращает минимальные и максимальные значения для числовых атрибутов категории
func (h *MarketplaceHandler) GetAttributeRanges(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	ranges, err := h.marketplaceService.Storage().GetAttributeRanges(c.Context(), categoryID)
	if err != nil {
		log.Printf("Error fetching attribute ranges: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attribute ranges")
	}

	return utils.SuccessResponse(c, ranges)
}

// GetCategoryAttributes возвращает атрибуты для указанной категории
func (h *MarketplaceHandler) GetCategoryAttributes(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("id")
	log.Printf("Requested attributes for category ID: %d", categoryID)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	attributes, err := h.marketplaceService.GetCategoryAttributes(c.Context(), categoryID)
	if err != nil {
		log.Printf("Error fetching category attributes for category %d: %v", categoryID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch category attributes")
	}

	log.Printf("Found %d attributes for category %d", len(attributes), categoryID)
	for i, attr := range attributes {
		log.Printf("Attribute %d: name=%s, type=%s, options=%v",
			i, attr.Name, attr.AttributeType, string(attr.Options))
	}

	return utils.SuccessResponse(c, attributes)
}

func (h *MarketplaceHandler) GetCategoryTree(c *fiber.Ctx) error {
	categoryTreeMutex.RLock()
	if time.Since(categoryTreeLastUpdate) < 5*time.Minute && categoryTreeCache != nil && len(categoryTreeCache) > 0 {
		categories := categoryTreeCache
		categoryTreeMutex.RUnlock()
		return utils.SuccessResponse(c, categories)
	}
	categoryTreeMutex.RUnlock()

	categoryTreeMutex.Lock()
	defer categoryTreeMutex.Unlock()

	// Повторная проверка после получения блокировки
	if time.Since(categoryTreeLastUpdate) < 5*time.Minute && categoryTreeCache != nil && len(categoryTreeCache) > 0 {
		return utils.SuccessResponse(c, categoryTreeCache)
	}

	categories, err := h.marketplaceService.GetCategoryTree(c.Context())
	if err != nil {
		// В случае ошибки очищаем кеш
		categoryTreeCache = nil
		categoryTreeLastUpdate = time.Time{}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching category tree")
	}

	// Проверяем валидность полученных данных
	if len(categories) == 0 {
		// Если данные пустые, не кешируем их
		categoryTreeCache = nil
		categoryTreeLastUpdate = time.Time{}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Received empty category tree")
	}

	// Сохраняем только валидные данные
	categoryTreeCache = categories
	categoryTreeLastUpdate = time.Now()

	return utils.SuccessResponse(c, categories)
}

func (h *MarketplaceHandler) InvalidateCategoryCache() {
	categoryTreeMutex.Lock()
	defer categoryTreeMutex.Unlock()
	categoryTreeCache = nil
	categoryTreeLastUpdate = time.Time{}
}
func (h *MarketplaceHandler) UploadImages(c *fiber.Ctx) error {
	log.Printf("Starting image upload for listing ID: %v", c.Params("id"))
	listingID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем файлы из формы
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка при получении файлов")
	}

	files := form.File["images"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Нет файлов для загрузки")
	}
	for _, file := range files {
		log.Printf("Received file: name=%s, size=%d, type=%s", file.Filename, file.Size, file.Header.Get("Content-Type"))
	}
	if len(files) > 10 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Максимум 10 фотографий")
	}

	var uploadedImages []models.MarketplaceImage
	mainImageIndex := 0
	if mainIdx := c.FormValue("main_image_index"); mainIdx != "" {
		mainImageIndex, _ = strconv.Atoi(mainIdx)
	}

	for i, file := range files {
		// Загружаем изображение
		image, err := h.marketplaceService.UploadImage(c.Context(), file, listingID, i == mainImageIndex)
		if err != nil {
			log.Printf("Failed to upload image: %v", err)
			continue
		}

		uploadedImages = append(uploadedImages, *image)
	}

	// Переиндексируем объявление с загруженными изображениями
	fullListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Failed to get full listing for reindexing: %v", err)
	} else {
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), fullListing); err != nil {
			log.Printf("Failed to reindex listing after image upload: %v", err)
		} else {
			log.Printf("Successfully reindexed listing %d with %d images", listingID, len(fullListing.Images))
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Изображения успешно загружены",
		"images":  uploadedImages,
	})
}
func (h *MarketplaceHandler) DeleteImage(c *fiber.Ctx) error {
	imageID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID изображения")
	}

	// Проверяем права доступа
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем информацию об изображении и объявлении
	image, err := h.marketplaceService.Storage().GetListingImageByID(c.Context(), imageID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Изображение не найдено")
	}

	// Проверяем, является ли пользователь владельцем объявления
	listing, err := h.marketplaceService.GetListingByID(c.Context(), image.ListingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения информации об объявлении")
	}

	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "У вас нет прав для удаления этого изображения")
	}

	// Удаляем изображение
	err = h.marketplaceService.DeleteImage(c.Context(), imageID)
	if err != nil {
		log.Printf("Failed to delete image: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка удаления изображения")
	}

	// Переиндексируем объявление
	updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listing.ID)
	if err != nil {
		log.Printf("Failed to get updated listing for reindexing: %v", err)
	} else {
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing); err != nil {
			log.Printf("Failed to reindex listing after image deletion: %v", err)
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Изображение успешно удалено",
	})
}

// ReindexRatings переиндексирует рейтинги всех объявлений
func (h *MarketplaceHandler) ReindexRatings(c *fiber.Ctx) error {
	// Проверяем административные права
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Запускаем процесс переиндексации в фоне
	go func() {
		ctx := context.Background()

		// Получаем все активные объявления
		rows, err := h.services.Storage().Query(ctx, `
            SELECT id FROM marketplace_listings 
            WHERE status = 'active'
        `)

		if err != nil {
			log.Printf("Ошибка при получении списка объявлений: %v", err)
			return
		}
		defer rows.Close()

		var listingIDs []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Printf("Ошибка при сканировании ID: %v", err)
				continue
			}
			listingIDs = append(listingIDs, id)
		}

		log.Printf("Начинаем обновление рейтингов для %d объявлений", len(listingIDs))

		// Для каждого объявления получаем рейтинг и обновляем
		for _, id := range listingIDs {
			// Получаем объявление
			listing, err := h.marketplaceService.GetListingByID(ctx, id)
			if err != nil {
				log.Printf("Ошибка получения объявления %d: %v", id, err)
				continue
			}

			// Получаем статистику отзывов напрямую из базы данных
			var reviewCount int
			var averageRating float64

			err = h.services.Storage().QueryRow(ctx, `
                SELECT COUNT(*), COALESCE(AVG(rating), 0)
                FROM reviews
                WHERE entity_type = 'listing' AND entity_id = $1 AND status = 'published'
            `, id).Scan(&reviewCount, &averageRating)

			if err != nil {
				log.Printf("Ошибка получения статистики отзывов для объявления %d: %v", id, err)
				// Если не удалось получить статистику, устанавливаем нулевые значения
				reviewCount = 0
				averageRating = 0
			}

			// Обновляем рейтинг в объекте
			listing.AverageRating = averageRating
			listing.ReviewCount = reviewCount

			// Переиндексируем объявление
			if err := h.marketplaceService.Storage().IndexListing(ctx, listing); err != nil {
				log.Printf("Ошибка индексации объявления %d: %v", id, err)
			} else {
				log.Printf("Обновлен рейтинг для объявления %d: %.2f (%d отзывов)",
					id, averageRating, reviewCount)
			}
		}

		log.Println("Переиндексация рейтингов успешно завершена")
	}()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Запущена переиндексация рейтингов всех объявлений",
	})
}

func (h *MarketplaceHandler) SynchronizeDiscounts(c *fiber.Ctx) error {
	// Проверяем административные права
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Запускаем синхронизацию в фоне
	go func() {
		ctx := context.Background()

		// Изменённый запрос: ищем все объявления с историей цен или упоминанием скидки
		rows, err := h.services.Storage().Query(ctx, `
            SELECT id 
            FROM marketplace_listings 
            WHERE (
                -- Объявления со скидками в описании или метаданных
                description LIKE '%СКИДКА%' OR metadata->>'discount' IS NOT NULL
                OR 
                -- Объявления, у которых есть история цен с изменениями
                id IN (
                    SELECT listing_id 
                    FROM price_history
                    GROUP BY listing_id
                    HAVING COUNT(*) > 1
                )
            )
            AND status = 'active'
        `)

		if err != nil {
			log.Printf("Ошибка при поиске объявлений со скидками: %v", err)
			return
		}

		var listingIDs []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Printf("Ошибка при сканировании ID объявления: %v", err)
				continue
			}
			listingIDs = append(listingIDs, id)
		}
		rows.Close()

		log.Printf("Найдено %d объявлений со скидками для синхронизации", len(listingIDs))

		// Обрабатываем каждое объявление
		for _, id := range listingIDs {
			if err := h.marketplaceService.SynchronizeDiscountData(ctx, id); err != nil {
				log.Printf("Ошибка синхронизации данных о скидке для объявления %d: %v", id, err)
			} else {
				log.Printf("Успешно синхронизированы данные о скидке для объявления %d", id)
			}
		}

		// Запускаем переиндексацию объявлений
		if err := h.marketplaceService.ReindexAllListings(ctx); err != nil {
			log.Printf("Ошибка переиндексации объявлений: %v", err)
		} else {
			log.Println("Синхронизация данных о скидках успешно завершена")
		}
	}()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Запущена синхронизация данных о скидках",
	})
}
func (h *MarketplaceHandler) ReindexAllListings(c *fiber.Ctx) error {
	// Проверяем административные права
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Запускаем процесс переиндексации в фоне
	go func() {
		ctx := context.Background()
		if err := h.marketplaceService.ReindexAllListings(ctx); err != nil {
			log.Printf("Ошибка переиндексации объявлений: %v", err)
		} else {
			log.Println("Переиндексация успешно завершена")
		}
	}()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Запущена переиндексация всех объявлений. Процесс может занять некоторое время.",
	})
}
func (h *MarketplaceHandler) GetEnhancedSuggestions(c *fiber.Ctx) error {
	prefix := c.Query("q", "")
	size := c.QueryInt("size", 8)

	log.Printf("Запрос расширенных подсказок, запрос: '%s', размер: %d", prefix, size)

	if prefix == "" {
		return utils.SuccessResponse(c, fiber.Map{
			"data": []interface{}{},
		})
	}

	// Структура для объединенных результатов
	type SuggestionItem struct {
		Type           string      `json:"type"`
		ID             interface{} `json:"id"`
		Title          string      `json:"title"`
		Display        string      `json:"display,omitempty"`
		Priority       int         `json:"priority"`
		CategoryID     int         `json:"category_id,omitempty"`
		Path           interface{} `json:"path,omitempty"`
		AttributeName  string      `json:"attribute_name,omitempty"`
		AttributeValue string      `json:"attribute_value,omitempty"`
	}

	var suggestions []SuggestionItem

	// 1. Получаем подсказки товаров
	// Сначала пытаемся через OpenSearch или стандартный поиск
	productTitles, err := h.marketplaceService.GetSuggestions(c.Context(), prefix, 3)
	if err != nil || len(productTitles) == 0 {
		// Если OpenSearch не дал результатов, используем прямой поиск в базе данных
		filters := map[string]string{
			"query": prefix,
		}

		products, _, err := h.marketplaceService.GetListings(c.Context(), filters, 3, 0)
		if err == nil && len(products) > 0 {
			for _, product := range products {
				if strings.Contains(strings.ToLower(product.Title), strings.ToLower(prefix)) {
					suggestions = append(suggestions, SuggestionItem{
						Type:       "product",
						ID:         product.ID,
						Title:      product.Title,
						Display:    product.Title,
						Priority:   1,
						CategoryID: product.CategoryID,
					})
				}
			}
		}
	} else {
		// Используем результаты из OpenSearch
		// Но нужно получить ID и категории этих товаров из базы данных
		filters := map[string]string{
			"title_exact": strings.Join(productTitles, "|"), // Специальный параметр для точного поиска по заголовкам
		}

		products, _, err := h.marketplaceService.GetListings(c.Context(), filters, len(productTitles), 0)
		if err == nil {
			// Создаем мапу найденных товаров для быстрого поиска
			productsMap := make(map[string]*models.MarketplaceListing)
			for i := range products {
				productsMap[products[i].Title] = &products[i]
			}

			// Добавляем товары в том же порядке, что и в подсказках
			for _, title := range productTitles {
				if product, ok := productsMap[title]; ok {
					suggestions = append(suggestions, SuggestionItem{
						Type:       "product",
						ID:         product.ID,
						Title:      product.Title,
						Display:    product.Title,
						Priority:   1,
						CategoryID: product.CategoryID,
					})
				} else {
					// Если товар не найден в базе, добавляем только заголовок
					suggestions = append(suggestions, SuggestionItem{
						Type:     "product",
						Title:    title,
						Display:  title,
						Priority: 1,
					})
				}
			}
		} else {
			// Если поиск в базе не удался, используем только заголовки
			for _, title := range productTitles {
				suggestions = append(suggestions, SuggestionItem{
					Type:     "product",
					Title:    title,
					Display:  title,
					Priority: 1,
				})
			}
		}
	}

	// 2. Получаем подсказки категорий
	attributeRows, err := h.marketplaceService.Storage().Query(c.Context(), `
	SELECT DISTINCT 
		ca.name as attr_name, 
		ca.display_name as attr_display, 
		lav.text_value, 
		COUNT(DISTINCT l.id) as listing_count,
		MIN(l.category_id) as category_id
	FROM listing_attribute_values lav
	JOIN category_attributes ca ON lav.attribute_id = ca.id
	JOIN marketplace_listings l ON lav.listing_id = l.id
	WHERE 
		lav.text_value IS NOT NULL 
		AND lav.text_value != ''
		AND LOWER(lav.text_value) LIKE LOWER($1)
		AND l.status = 'active'
		AND ca.name IN ('make', 'model', 'brand', 'property_type', 'body_type')
	GROUP BY ca.name, ca.display_name, lav.text_value
	ORDER BY listing_count DESC
	LIMIT $2
`, "%"+prefix+"%", size)

	if err == nil {
		defer attributeRows.Close()

		for attributeRows.Next() {
			var attrName, attrDisplay, attrValue string
			var listingCount, categoryID int

			if err := attributeRows.Scan(&attrName, &attrDisplay, &attrValue, &listingCount, &categoryID); err != nil {
				log.Printf("Ошибка сканирования атрибута: %v", err)
				continue
			}

			// Определяем приоритет атрибута
			priority := 2
			if attrName == "model" || attrName == "make" {
				priority = 1 // Высокий приоритет для марки и модели
			}

			// Формируем текст для отображения
			display := fmt.Sprintf("%s: %s (%d)", attrDisplay, attrValue, listingCount)

			suggestions = append(suggestions, SuggestionItem{
				Type:           "attribute",
				Title:          attrValue,
				Display:        display,
				Priority:       priority,
				CategoryID:     categoryID,
				AttributeName:  attrName,
				AttributeValue: attrValue,
			})
		}
	}

	categorySuggestions, err := h.marketplaceService.GetCategorySuggestions(c.Context(), prefix, 3)
	if err == nil && len(categorySuggestions) > 0 {
		for _, category := range categorySuggestions {
			suggestions = append(suggestions, SuggestionItem{
				Type:     "category",
				ID:       category.ID,
				Title:    category.Name,
				Display:  fmt.Sprintf("Категория: %s (%d)", category.Name, category.ListingCount),
				Priority: 2,
			})
		}
	}

	// 3. Дополнительно извлекаем категории для найденных товаров
	productCategoryIDs := make(map[int]bool)

	for _, suggestion := range suggestions {
		if suggestion.Type == "product" && suggestion.CategoryID > 0 {
			productCategoryIDs[suggestion.CategoryID] = true
		}
	}

	if len(productCategoryIDs) > 0 {
		// Получаем дерево категорий
		categoryTree, err := h.marketplaceService.GetCategoryTree(c.Context())
		if err == nil {
			// Функция для рекурсивного поиска категории
			var findCategory func(categories []models.CategoryTreeNode, id int, path []map[string]interface{}) (models.CategoryTreeNode, []map[string]interface{}, bool)
			findCategory = func(categories []models.CategoryTreeNode, id int, path []map[string]interface{}) (models.CategoryTreeNode, []map[string]interface{}, bool) {
				for _, category := range categories {
					currentPath := append(path, map[string]interface{}{
						"id":   category.ID,
						"name": category.Name,
						"slug": category.Slug,
					})

					if category.ID == id {
						return category, currentPath, true
					}

					if len(category.Children) > 0 {
						if foundCategory, foundPath, found := findCategory(category.Children, id, currentPath); found {
							return foundCategory, foundPath, true
						}
					}
				}

				return models.CategoryTreeNode{}, nil, false
			}

			// Добавляем категории из найденных товаров
			addedCategories := make(map[int]bool)

			for catID := range productCategoryIDs {
				if _, exists := addedCategories[catID]; exists {
					continue
				}

				category, path, found := findCategory(categoryTree, catID, []map[string]interface{}{})
				if found {
					suggestions = append(suggestions, SuggestionItem{
						Type:     "category",
						ID:       category.ID,
						Title:    category.Name,
						Display:  "Категория: " + category.Name,
						Priority: 2,
						Path:     path,
					})
					addedCategories[catID] = true

					// Добавляем родительскую категорию, если она есть
					if category.ParentID != nil && *category.ParentID > 0 {
						parentCategory, parentPath, found := findCategory(categoryTree, *category.ParentID, []map[string]interface{}{})
						if found && !addedCategories[*category.ParentID] {
							suggestions = append(suggestions, SuggestionItem{
								Type:     "category",
								ID:       parentCategory.ID,
								Title:    parentCategory.Name,
								Display:  "Раздел: " + parentCategory.Name,
								Priority: 3,
								Path:     parentPath,
							})
							addedCategories[*category.ParentID] = true
						}
					}
				}
			}
		}
	}

	// 4. Сортируем результаты по приоритету
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Priority < suggestions[j].Priority
	})

	// 5. Ограничиваем количество результатов

	if len(suggestions) > size {
		suggestions = suggestions[:size]
	}

	log.Printf("Найдено %d расширенных подсказок для запроса '%s'", len(suggestions), prefix)

	return utils.SuccessResponse(c, fiber.Map{
		"data": suggestions,
	})
}

func (h *MarketplaceHandler) GetListings(c *fiber.Ctx) error {
	filters := map[string]string{
		"category_id":   c.Query("category_id"),
		"city":          c.Query("city"),
		"min_price":     c.Query("min_price"),
		"max_price":     c.Query("max_price"),
		"query":         c.Query("query"),
		"condition":     c.Query("condition"),
		"sort_by":       c.Query("sort_by"),
		"storefront_id": c.Query("storefront_id"),
		"user_id":       c.Query("user_id"),
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, limit, offset)
	if err != nil {
		log.Printf("Error getting listings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching listings")
	}

	//log.Printf("Found %d listings", len(listings))

	return utils.SuccessResponse(c, fiber.Map{
		"data": listings,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
func (h *MarketplaceHandler) AddToFavorites(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	err = h.marketplaceService.AddToFavorites(c.Context(), userID, listingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding to favorites")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Added to favorites successfully",
	})
}
func (h *MarketplaceHandler) RemoveFromFavorites(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	err = h.marketplaceService.RemoveFromFavorites(c.Context(), userID, listingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error removing from favorites")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Removed from favorites successfully",
	})
}

func (h *MarketplaceHandler) GetListing(c *fiber.Ctx) error {
	// Получаем user_id из контекста, если пользователь авторизован
	var userID int
	if uid := c.Locals("user_id"); uid != nil {
		var ok bool
		userID, ok = uid.(int)
		if !ok {
			log.Printf("Invalid user_id type in context: %T", uid)
			userID = 0
		}
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	// Добавить проверку на валидность ID
	if id <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	// Создаем контекст с user_id и флагом увеличения счетчика просмотров
	ctx := context.WithValue(c.Context(), "user_id", userID)
	// Устанавливаем флаг увеличения счетчика просмотров для API получения деталей объявления
	ctx = context.WithValue(ctx, "increment_views", true)
	// Добавляем IP-адрес для отслеживания просмотров неавторизованных пользователей
	// Получаем реальный IP-адрес пользователя из заголовков, переданных прокси
	ipAddress := c.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = c.Get("X-Real-IP")
		if ipAddress == "" {
			// Если заголовки прокси не установлены, используем стандартный IP
			ipAddress = c.IP()
		}
	}
	// Если в X-Forwarded-For несколько адресов через запятую, берем первый (клиентский)
	if commaPos := strings.Index(ipAddress, ","); commaPos > 0 {
		ipAddress = strings.TrimSpace(ipAddress[:commaPos])
	}
	log.Printf("Использую IP-адрес для счетчика просмотров: %s", ipAddress)
	ctx = context.WithValue(ctx, "ip_address", ipAddress)

	listing, err := h.marketplaceService.GetListingByID(ctx, id)
	if err != nil {
		log.Printf("Error getting listing %d: %v", id, err)
		if err.Error() == "listing not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching listing")
	}

	// Добавляем отладочную информацию
	log.Printf("DEBUG: Listing %d attributes count: %d", listing.ID, len(listing.Attributes))
	if len(listing.Attributes) > 0 {
		log.Printf("DEBUG: Creating listing with %d attributes", len(listing.Attributes))
		for i, attr := range listing.Attributes {
			log.Printf("DEBUG: Attribute %d: ID=%d, Name=%s, Type=%s",
				i, attr.AttributeID, attr.AttributeName, attr.AttributeType)
		}
	}

	return utils.SuccessResponse(c, listing)
}

func (h *MarketplaceHandler) UpdateListing(c *fiber.Ctx) error {
	// Получаем старую версию объявления
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	oldListing, err := h.marketplaceService.GetListingByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
	}

	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
	}

	listing.ID = id
	listing.UserID = c.Locals("user_id").(int)
	if err := c.BodyParser(&listing); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
	}

	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		if _, hasAttrs := requestBody["attributes"]; !hasAttrs {
			// Атрибуты не переданы, загружаем существующие
			existingAttrs, err := h.marketplaceService.Storage().GetListingAttributes(c.Context(), listing.ID)
			if err == nil && len(existingAttrs) > 0 {
				listing.Attributes = existingAttrs
				log.Printf("Loaded existing %d attributes for listing %d", len(existingAttrs), listing.ID)
			}
		} else {
			// Атрибуты переданы, обрабатываем их как обычно
			processAttributesFromRequest(requestBody, &listing)
		}
	}

	// Обработка изменения цены и метаданных о скидке
	if listing.Price != oldListing.Price {
		// Если у объявления нет метаданных, создаем их
		if listing.Metadata == nil {
			listing.Metadata = make(map[string]interface{})
		}

		// Проверяем, снизилась ли цена
		if listing.Price < oldListing.Price {
			// Определяем исходную цену для расчета скидки
			var originalPrice float64
			var hasExistingDiscount bool

			// Проверяем существующие метаданные о скидке
			if oldListing.Metadata != nil {
				if discount, ok := oldListing.Metadata["discount"].(map[string]interface{}); ok {
					if prevPrice, ok := discount["previous_price"].(float64); ok && prevPrice > 0 {
						// Используем предыдущую цену из существующей скидки
						originalPrice = prevPrice
						hasExistingDiscount = true
						log.Printf("Найдена существующая скидка для объявления %d. Предыдущая цена: %.2f", listing.ID, originalPrice)
					}
				}
			}

			// Если нет существующей скидки, используем текущую цену объявления
			if !hasExistingDiscount {
				originalPrice = oldListing.Price
			}

			// Вычисляем процент скидки от исходной цены
			discountPercent := int((originalPrice - listing.Price) / originalPrice * 100)

			// Добавляем или обновляем информацию о скидке в метаданные
			listing.Metadata["discount"] = map[string]interface{}{
				"discount_percent":  discountPercent,
				"previous_price":    originalPrice,
				"effective_from":    time.Now().Format(time.RFC3339),
				"has_price_history": true,
			}

			log.Printf("Обновлена информация о скидке для объявления %d: %d%% (исходная цена: %.2f, новая цена: %.2f)",
				listing.ID, discountPercent, originalPrice, listing.Price)
		} else if listing.Price > oldListing.Price {
			// Если цена повысилась, проверяем, нужно ли удалить информацию о скидке
			if oldListing.Metadata != nil {
				if _, ok := oldListing.Metadata["discount"]; ok {
					// Удаляем метаданные о скидке при повышении цены
					delete(listing.Metadata, "discount")
					log.Printf("Удалена информация о скидке для объявления %d из-за повышения цены", listing.ID)
				}
			}
		}
		// Если цены равны, оставляем метаданные без изменений
	}

	// Обновляем объявление
	err = h.marketplaceService.UpdateListing(c.Context(), &listing)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating listing")
	}
	// Переиндексируем объявление в OpenSearch с переводами
	// Даём небольшую задержку для завершения транзакций
	time.Sleep(500 * time.Millisecond)

	updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listing.ID)
	if err != nil {
		log.Printf("Failed to get updated listing for reindexing: %v", err)
	} else {
		// Проверяем наличие переводов
		if updatedListing.Translations == nil || len(updatedListing.Translations) == 0 {
			log.Printf("Warning: Listing %d has no translations before reindexing, will try to load them explicitly", listing.ID)

			// Пытаемся явно загрузить переводы
			translations, err := h.marketplaceService.Storage().GetTranslationsForEntity(c.Context(), "listing", listing.ID)
			if err != nil {
				log.Printf("Error loading translations for listing %d: %v", listing.ID, err)
			} else if len(translations) > 0 {
				// Организуем переводы в структуру TranslationMap
				transMap := make(models.TranslationMap)
				for _, t := range translations {
					if _, ok := transMap[t.Language]; !ok {
						transMap[t.Language] = make(map[string]string)
					}
					transMap[t.Language][t.FieldName] = t.TranslatedText
				}
				updatedListing.Translations = transMap
				log.Printf("Explicitly loaded %d translations for listing %d", len(translations), listing.ID)
			}
		}

		if err := h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing); err != nil {
			log.Printf("Failed to reindex listing %d: %v", listing.ID, err)
		} else {
			log.Printf("Successfully reindexed listing %d with all translations", listing.ID)
		}
	}
	// Проверяем изменение цены
	if oldListing.Price != listing.Price {
		// Вызываем явно синхронизацию данных о скидке и истории цен
		err = h.marketplaceService.SynchronizeDiscountData(c.Context(), listing.ID)
		if err != nil {
			log.Printf("Ошибка при синхронизации скидки для объявления %d: %v", listing.ID, err)
			// Не возвращаем ошибку клиенту, чтобы не прерывать обновление
		}

		favoritedUsers, err := h.marketplaceService.GetFavoritedUsers(c.Context(), listing.ID)
		if err != nil {
			log.Printf("Error getting favorited users: %v", err)
		} else {
			priceChange := listing.Price - oldListing.Price
			changeText := "увеличилась"
			if priceChange < 0 {
				changeText = "уменьшилась"
			}

			for _, userID := range favoritedUsers {
				notificationText := fmt.Sprintf(
					"Изменение цены в избранном\nОбъявление: %s\nЦена %s на %.2f руб.\nНовая цена: %.2f руб.",
					listing.Title,
					changeText,
					math.Abs(float64(priceChange)),
					listing.Price,
				)
				if err := h.services.Notification().SendNotification(
					c.Context(),
					userID,
					models.NotificationTypeFavoritePrice,
					notificationText,
					listing.ID,
				); err != nil {
					log.Printf("Error sending notification to user %d: %v", userID, err)
				}
			}
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Listing updated successfully",
	})
}

func (h *MarketplaceHandler) UpdateTranslations(c *fiber.Ctx) error {
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	var updateData struct {
		Language     string            `json:"language"`
		Translations map[string]string `json:"translations"`
		IsVerified   bool              `json:"is_verified"`
		Provider     string            `json:"provider"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
	}

	// Получаем провайдер из запроса или из параметра запроса
	provider := updateData.Provider
	if provider == "" {
		provider = c.Query("translation_provider", "google")
	}

	// Проверяем корректность провайдера и приводим к типу TranslationProvider
	translationProvider := service.GoogleTranslate
	if provider == "openai" {
		translationProvider = service.OpenAI
	}

	// Обновляем каждый переведенный field
	for fieldName, translatedText := range updateData.Translations {
		translation := &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            updateData.Language,
			FieldName:           fieldName,
			TranslatedText:      translatedText,
			IsVerified:          updateData.IsVerified,
			IsMachineTranslated: false,
			Metadata:            map[string]interface{}{"provider": provider},
		}

		// Передаем информацию о провайдере в сервис
		err := h.marketplaceService.UpdateTranslationWithProvider(c.Context(), translation, translationProvider)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating translation")
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Translations updated successfully",
	})
}

func (h *MarketplaceHandler) DeleteListing(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	// Получаем объявление с изображениями до удаления
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		if err.Error() == "listing not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching listing")
	}

	// Проверяем права доступа
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to delete this listing")
	}

	// Сохраняем пути к изображениям
	var imagePaths []string
	if listing.Images != nil {
		for _, img := range listing.Images {
			if img.FilePath != "" {
				imagePaths = append(imagePaths, img.FilePath)
			}
		}
	}

	// Удаляем объявление
	err = h.marketplaceService.DeleteListing(c.Context(), listingID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting listing")
	}

	// Удаляем файлы изображений с диска
	for _, path := range imagePaths {
		filePath := fmt.Sprintf("./uploads/%s", path)
		if err := os.Remove(filePath); err != nil {
			log.Printf("Error removing image file %s: %v", filePath, err)
		} else {
			log.Printf("Successfully removed image file: %s", filePath)
		}
	}

	// Удаляем из индекса OpenSearch
	if err := h.marketplaceService.Storage().DeleteListingIndex(c.Context(), fmt.Sprintf("%d", listingID)); err != nil {
		log.Printf("Error removing listing from search index: %v", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Listing deleted successfully",
	})
}

// GetCategories - получение списка категорий
func (h *MarketplaceHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.marketplaceService.GetCategories(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching categories")
	}

	return utils.SuccessResponse(c, categories)
}
func (h *MarketplaceHandler) GetFavorites(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		//   log.Printf("GetFavorites: no user_id in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// log.Printf("GetFavorites: fetching favorites for userID=%d", userID)

	ctx := context.WithValue(c.Context(), "user_id", userID)

	favorites, err := h.marketplaceService.GetUserFavorites(ctx, userID)
	if err != nil {
		log.Printf("GetFavorites error: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении избранных объявлений")
	}

	//  log.Printf("GetFavorites: found %d favorites for userID=%d", len(favorites), userID)
	return utils.SuccessResponse(c, favorites)
}

// GetCategorySuggestions возвращает предложения категорий на основе поискового запроса
func (h *MarketplaceHandler) GetCategorySuggestions(c *fiber.Ctx) error {
	query := c.Query("q", "")
	size := c.QueryInt("size", 3)

	log.Printf("Запрос на получение предложений категорий, запрос: '%s', размер: %d", query, size)

	if query == "" {
		return utils.SuccessResponse(c, fiber.Map{
			"data": []interface{}{},
		})
	}

	// Выполняем SQL-запрос для поиска категорий, связанных с запросом
	sqlQuery := `
        WITH RECURSIVE category_tree AS (
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            WHERE 1=1
            
            UNION
            
            SELECT c.id, c.name, c.parent_id
            FROM marketplace_categories c
            JOIN category_tree t ON c.parent_id = t.id
        ),
        matching_categories AS (
            SELECT 
                c.id,
                c.name,
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND ml.status = 'active') as listing_count,
                CASE WHEN LOWER(c.name) LIKE LOWER($1) THEN 100 ELSE 0 END +
                (SELECT COUNT(*) FROM marketplace_listings ml 
                 WHERE ml.category_id = c.id 
                 AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1)) 
                 AND ml.status = 'active') as relevance
            FROM marketplace_categories c
            WHERE LOWER(c.name) LIKE LOWER($1)
            OR EXISTS (
                SELECT 1 FROM marketplace_listings ml 
                WHERE ml.category_id = c.id 
                AND (LOWER(ml.title) LIKE LOWER($1) OR LOWER(ml.description) LIKE LOWER($1))
                AND ml.status = 'active'
            )
        )
        SELECT id, name, listing_count
        FROM matching_categories
        WHERE listing_count > 0
        ORDER BY relevance DESC, listing_count DESC
        LIMIT $2
    `

	rows, err := h.marketplaceService.Storage().Query(c.Context(), sqlQuery, "%"+query+"%", size)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса категорий: %v", err)
		return utils.SuccessResponse(c, fiber.Map{
			"data": []interface{}{},
		})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var listingCount int

		if err := rows.Scan(&id, &name, &listingCount); err != nil {
			log.Printf("Ошибка сканирования категории: %v", err)
			continue
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"name":          name,
			"listing_count": listingCount,
		})
	}

	log.Printf("Найдено %d релевантных категорий для запроса '%s'", len(results), query)

	return utils.SuccessResponse(c, fiber.Map{
		"data": results,
	})
}

// Обновите метод в backend/internal/proj/marketplace/handler/marketplace.go

func (h *MarketplaceHandler) SearchListingsAdvanced(c *fiber.Ctx) error {
	// Получаем параметры поиска
	log.Printf("Все параметры запроса: %+v", c.Queries())

	attributeFilters := make(map[string]string)

	// Более надежный способ извлечения атрибутов - проходим по всем параметрам запроса
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if strings.HasPrefix(keyStr, "attr_") {
			attrName := strings.TrimPrefix(keyStr, "attr_")
			valueStr := string(value)
			attributeFilters[attrName] = valueStr
			log.Printf("Извлечен атрибут из запроса: %s = %s", attrName, valueStr)
		}
	})

	log.Printf("Извлеченные атрибуты ФФФильтров: %+v", attributeFilters)

	log.Printf("Исходные параметры сортировки из запроса: sort_by=%s", c.Query("sort_by", ""))

	// ВАЖНОЕ ИЗМЕНЕНИЕ: Получаем параметр view_mode и добавляем его в логи
	viewMode := c.Query("view_mode", "")
	log.Printf("Параметр режима просмотра из URL: view_mode=%s", viewMode)

	// Логирование всех параметров запроса для отладки
	log.Printf("Все параметры запроса: %+v", c.Queries())

	params := &search.ServiceParams{
		Query:            c.Query("q", c.Query("query", "")),
		CategoryID:       c.Query("category_id", ""),
		Condition:        c.Query("condition", ""),
		City:             c.Query("city", ""),
		Country:          c.Query("country", ""),
		StorefrontID:     c.Query("storefront_id", ""),
		Sort:             c.Query("sort_by", ""),
		SortDirection:    c.Query("sort_direction", "desc"),
		Distance:         c.Query("distance", ""),
		Page:             c.QueryInt("page", 1),
		Size:             c.QueryInt("size", 20),
		Language:         c.Query("language", ""),
		AttributeFilters: attributeFilters,
	}

	// Обязательные параметры для публичных запросов
	params.Status = "active"

	// УЛУЧШЕНИЕ: Детальное логирование параметров запроса
	log.Printf("Параметры пагинации: page=%d, size=%d, view_mode=%s",
		params.Page, params.Size, viewMode)

	// Установка размера в зависимости от режима просмотра
	if params.Page < 1 {
		params.Page = 1
	}

	// УЛУЧШЕНИЕ: Принимаем решение о размере выборки на основе viewMode
	if viewMode == "map" {
		// Для режима карты устанавливаем очень большой размер страницы
		// независимо от запрошенного размера
		log.Printf("Обнаружен режим просмотра 'map'. Устанавливаем большой размер страницы.")
		params.Size = 5000 // Устанавливаем максимальное значение для карты
	} else if params.Size < 1 {
		// Для обычного просмотра используем стандартные ограничения
		params.Size = 20
	} else if params.Size > 1000 {
		// Ограничиваем максимальный размер для обычного просмотра
		params.Size = 100
	}

	// Логирование финальных параметров поиска
	log.Printf("Итоговые параметры поиска: size=%d, view_mode=%s", params.Size, viewMode)

	// Добавьте этот код после инициализации params
	log.Printf("Итоговый параметр сортировки в запросе к OpenSearch: %s", params.Sort)

	// Проверяем и логируем параметры пагинации
	log.Printf("Параметры пагинации: page=%d, size=%d", params.Page, params.Size)
	log.Printf("Параметр сортировки получен из запроса: %s", params.Sort)

	// ВАЖНОЕ ИЗМЕНЕНИЕ: Добавляем view_mode в логи
	log.Printf("Параметры пагинации: page=%d, size=%d, view_mode=%s",
		params.Page, params.Size, viewMode)

	// Устанавливаем разумные ограничения на параметры пагинации
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		params.Size = 20
	} else if params.Size > 1000 {
		// Ограничиваем максимальный размер страницы, но увеличиваем для карты
		// ВАЖНОЕ ИЗМЕНЕНИЕ: Проверяем, предназначен ли запрос для отображения на карте
		if viewMode == "map" {
			// Для карты разрешаем больший размер страницы
			params.Size = 5000
			log.Printf("Установлен максимальный размер для режима карты: %d", params.Size)
		} else {
			// Для обычного списка оставляем прежнее ограничение для производительности
			params.Size = 100
			log.Printf("Установлен стандартный максимальный размер для списка: %d", params.Size)
		}
	}
	log.Printf("Параметры пагинации: page=%d, size=%d, view_mode=%s",
		params.Page, params.Size, c.Query("view_mode", ""))
	// Дополнительное логирование
	log.Printf("Полные параметры поиска: %+v", params)
	log.Printf("Атрибуты фильтров: %+v", attributeFilters)

	// ИСПРАВЛЕНИЕ: сначала проверяем наличие координат, потом устанавливаем distance
	latParam := c.Query("latitude", "")
	lonParam := c.Query("longitude", "")

	if latParam != "" && lonParam != "" {
		lat, errLat := strconv.ParseFloat(latParam, 64)
		lon, errLon := strconv.ParseFloat(lonParam, 64)

		if errLat == nil && errLon == nil && (lat != 0 || lon != 0) {
			params.Latitude = lat
			params.Longitude = lon
			log.Printf("Установлены координаты: lat=%.6f, lon=%.6f", lat, lon)
		}
	}

	// Теперь, когда у нас есть координаты, проверяем параметр distance
	if params.Distance != "" && params.Latitude != 0 && params.Longitude != 0 {
		log.Printf("Установлен фильтр по расстоянию: %s от координат (%.6f, %.6f)",
			params.Distance, params.Latitude, params.Longitude)

		// Проверить наличие индекса перед установкой координат
		if err := h.marketplaceService.Storage().PrepareIndex(c.Context()); err != nil {
			log.Printf("Ошибка проверки индекса: %v", err)
			// Но продолжаем выполнение, просто не используем гео-поиск
			params.Distance = ""
		}
	} else if params.Distance != "" {
		log.Printf("Параметр distance указан (%s), но координаты отсутствуют или равны нулю (%.6f, %.6f). "+
			"Параметр distance будет проигнорирован.",
			params.Distance, params.Latitude, params.Longitude)
	}

	log.Printf("Полученный поисковый запрос: %s", params.Query)
	// Обрабатываем числовые параметры
	if priceMin := c.Query("min_price", ""); priceMin != "" {
		if val, err := strconv.ParseFloat(priceMin, 64); err == nil && val >= 0 {
			params.PriceMin = val
		}
	}

	if priceMax := c.Query("max_price", ""); priceMax != "" {
		if val, err := strconv.ParseFloat(priceMax, 64); err == nil && val >= 0 {
			params.PriceMax = val
		}
	}

	// Запрашиваемые агрегации
	if aggs := c.Query("aggs", ""); aggs != "" {
		params.Aggregations = strings.Split(aggs, ",")
	}

	// Если не указан язык, берем из context
	if params.Language == "" {
		if lang, ok := c.Locals("language").(string); ok && lang != "" {
			params.Language = lang
		} else {
			params.Language = "sr"
		}
	}

	// Выполняем поиск
	result, err := h.marketplaceService.SearchListingsAdvanced(c.Context(), params)
	if err != nil {
		log.Printf("Ошибка поиска: %v", err)

		// Используем стандартный поиск
		filters := map[string]string{
			"category_id":   params.CategoryID,
			"condition":     params.Condition,
			"city":          params.City,
			"country":       params.Country,
			"storefront_id": params.StorefrontID,
			"sort_by":       params.Sort,
		}

		// Добавляем числовые фильтры, если они указаны
		if params.PriceMin > 0 {
			filters["min_price"] = fmt.Sprintf("%g", params.PriceMin)
		}
		if params.PriceMax > 0 {
			filters["max_price"] = fmt.Sprintf("%g", params.PriceMax)
		}

		// Добавляем текстовый поиск
		if params.Query != "" {
			filters["query"] = params.Query
		}

		// Определяем смещение для пагинации
		offset := (params.Page - 1) * params.Size

		// Пробуем получить обычным методом
		listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, params.Size, offset)
		if err != nil {
			log.Printf("Ошибка стандартного поиска: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка выполнения поиска")
		}

		// Указываем метаданные пагинации в ответе
		totalPages := int(math.Ceil(float64(total) / float64(params.Size)))

		// Формируем такой же ответ, как от OpenSearch
		return utils.SuccessResponse(c, fiber.Map{
			"data": listings,
			"meta": fiber.Map{
				"total":       total,
				"page":        params.Page,
				"size":        params.Size,
				"total_pages": totalPages,
				"has_more":    params.Page < totalPages,
			},
		})
	}

	// После получения результатов поиска
	log.Printf("Результаты поиска: найдено %d объявлений", len(result.Items))

	// Формируем информацию о пагинации
	hasMore := result.Page < result.TotalPages

	// Если OpenSearch ответил успешно
	return utils.SuccessResponse(c, fiber.Map{
		"data": result.Items,
		"meta": fiber.Map{
			"total":               result.Total,
			"page":                result.Page,
			"size":                result.Size,
			"total_pages":         result.TotalPages,
			"has_more":            hasMore,
			"facets":              result.Facets,
			"suggestions":         result.Suggestions,
			"took_ms":             result.Took,
			"spelling_suggestion": result.SpellingSuggestion,
		},
	})
}

// GetSuggestions возвращает предложения автодополнения
func (h *MarketplaceHandler) GetSuggestions(c *fiber.Ctx) error {
	prefix := c.Query("q", "")
	size := c.QueryInt("size", 5)

	log.Printf("Запрос на получение подсказок, запрос: '%s', размер: %d", prefix, size)

	if prefix == "" {
		return utils.SuccessResponse(c, fiber.Map{
			"data": []string{},
		})
	}

	// Пытаемся получить предложения из OpenSearch
	suggestions, err := h.marketplaceService.GetSuggestions(c.Context(), prefix, size)
	if err != nil {
		log.Printf("Ошибка получения предложений из OpenSearch: %v", err)

		// Используем более простой поиск через базу данных
		// Здесь можно реализовать запасной вариант поиска по префиксу в PostgreSQL
		filters := map[string]string{
			"query": prefix + "%", // Используем префикс для LIKE запроса
		}

		listings, _, err := h.marketplaceService.GetListings(c.Context(), filters, size, 0)
		if err != nil {
			log.Printf("Ошибка запасного поиска: %v", err)
			// В случае полной неудачи возвращаем пустой массив
			return utils.SuccessResponse(c, fiber.Map{
				"data": []string{},
			})
		}

		// Извлекаем названия из найденных объявлений
		titles := make([]string, 0, len(listings))
		for _, listing := range listings {
			titles = append(titles, listing.Title)
		}

		log.Printf("Получены подсказки из базы данных: %v", titles)

		return utils.SuccessResponse(c, fiber.Map{
			"data": titles,
		})
	}

	log.Printf("Получены подсказки из OpenSearch: %v", suggestions)

	return utils.SuccessResponse(c, fiber.Map{
		"data": suggestions,
	})
}

// ReindexAll переиндексирует все объявления
func (h *MarketplaceHandler) ReindexAll(c *fiber.Ctx) error {
	// Проверяем административные права
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Запускаем процесс переиндексации в фоне
	go func() {
		ctx := context.Background()
		if err := h.marketplaceService.ReindexAllListings(ctx); err != nil {
			log.Printf("Ошибка переиндексации: %v", err)
		} else {
			log.Println("Переиндексация успешно завершена")
		}
	}()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Запущена переиндексация всех объявлений",
	})
}

// ReindexAllWithTranslations переиндексирует все объявления с явной загрузкой переводов
func (h *MarketplaceHandler) ReindexAllWithTranslations(c *fiber.Ctx) error {
	// Проверяем административные права
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Запускаем реиндексацию в фоне
	go func() {
		ctx := context.Background()

		// Получаем все ID объявлений
		rows, err := h.services.Storage().Query(ctx, `
            SELECT id FROM marketplace_listings 
            WHERE status = 'active'
            ORDER BY id
        `)

		if err != nil {
			log.Printf("Error getting listing IDs for reindex: %v", err)
			return
		}
		defer rows.Close()

		var listingIDs []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Printf("Error scanning listing ID: %v", err)
				continue
			}
			listingIDs = append(listingIDs, id)
		}

		log.Printf("Starting reindex of %d listings with translations...", len(listingIDs))

		// Реиндексируем каждое объявление с явной загрузкой переводов
		count := 0
		for _, id := range listingIDs {
			// Получаем объявление со всеми данными, включая атрибуты
			listing, err := h.marketplaceService.GetListingByID(ctx, id)
			if err != nil {
				log.Printf("Error getting listing %d: %v", id, err)
				continue
			}

			// Специальная обработка для объявлений без атрибутов
			if listing.Attributes == nil || len(listing.Attributes) == 0 {
				attrs, err := h.marketplaceService.Storage().GetListingAttributes(ctx, id)
				if err == nil && len(attrs) > 0 {
					listing.Attributes = attrs
					log.Printf("Loaded %d attributes for listing %d during reindex", len(attrs), id)
				}
			}

			// Проверяем наличие переводов и явно загружаем их если нужно
			if listing.Translations == nil || len(listing.Translations) == 0 {
				transMap := make(models.TranslationMap)

				// Используем прямой SQL-запрос для надежности
				rows, err := h.marketplaceService.Storage().Query(ctx, `
                    SELECT language, field_name, translated_text 
                    FROM translations 
                    WHERE entity_type = 'listing' AND entity_id = $1
                `, id)

				if err == nil {
					defer rows.Close()
					for rows.Next() {
						var lang, field, text string
						if err := rows.Scan(&lang, &field, &text); err == nil {
							if _, ok := transMap[lang]; !ok {
								transMap[lang] = make(map[string]string)
							}
							transMap[lang][field] = text
						}
					}

					if len(transMap) > 0 {
						listing.Translations = transMap
						log.Printf("Loaded %d languages of translations for listing %d", len(transMap), id)
					}
				}
			}

			// Индексируем объявление со всеми данными
			if err := h.marketplaceService.Storage().IndexListing(ctx, listing); err != nil {
				log.Printf("Error indexing listing %d: %v", id, err)
			} else {
				count++
				if count%20 == 0 {
					log.Printf("Progress: %d/%d listings reindexed", count, len(listingIDs))
				}
			}
		}

		log.Printf("Successfully reindexed %d/%d listings with translations", count, len(listingIDs))
	}()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Переиндексация всех объявлений с переводами запущена...",
	})
}

func (h *MarketplaceHandler) GetPriceHistory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Error parsing listing ID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	// Получаем историю цен
	history, err := h.marketplaceService.GetPriceHistory(c.Context(), id)
	if err != nil {
		log.Printf("Error fetching price history for listing %d: %v", id, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch price history")
	}

	return utils.SuccessResponse(c, history)
}
