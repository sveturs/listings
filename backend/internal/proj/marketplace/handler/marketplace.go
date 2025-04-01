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

		// Модерируем текст перед сохранением
		moderatedTitle, err := h.services.Translation().ModerateText(c.Context(), listing.Title, listing.OriginalLanguage)
		if err != nil {
			log.Printf("Failed to moderate title: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка модерации заголовка")
		}
		listing.Title = moderatedTitle

		moderatedDesc, err := h.services.Translation().ModerateText(c.Context(), listing.Description, listing.OriginalLanguage)
		if err != nil {
			log.Printf("Failed to moderate description: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка модерации описания")
		}
		listing.Description = moderatedDesc
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

		// Переводим на другие языки
		targetLanguages := []string{"en", "ru", "sr"}
		for _, targetLang := range targetLanguages {
			if targetLang == listing.OriginalLanguage {
				continue
			}

			// Переводим заголовок
			translatedTitle, err := h.services.Translation().Translate(c.Context(), listing.Title, listing.OriginalLanguage, targetLang)
			if err == nil {
				err = h.marketplaceService.UpdateTranslation(c.Context(), &models.Translation{
					EntityType:          "listing",
					EntityID:            listingID,
					Language:            targetLang,
					FieldName:           "title",
					TranslatedText:      translatedTitle,
					IsMachineTranslated: true,
					IsVerified:          false,
				})
				if err != nil {
					log.Printf("Error saving %s title translation: %v", targetLang, err)
				}
			} else {
				log.Printf("Error translating title to %s: %v", targetLang, err)
			}

			// Переводим описание
			translatedDesc, err := h.services.Translation().Translate(c.Context(), listing.Description, listing.OriginalLanguage, targetLang)
			if err == nil {
				err = h.marketplaceService.UpdateTranslation(c.Context(), &models.Translation{
					EntityType:          "listing",
					EntityID:            listingID,
					Language:            targetLang,
					FieldName:           "description",
					TranslatedText:      translatedDesc,
					IsMachineTranslated: true,
					IsVerified:          false,
				})
				if err != nil {
					log.Printf("Error saving %s description translation: %v", targetLang, err)
				}
			} else {
				log.Printf("Error translating description to %s: %v", targetLang, err)
			}
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

// GetSimilarListings возвращает похожие объявления
func (h *MarketplaceHandler) GetSimilarListings(c *fiber.Ctx) error {
    listingID, err := c.ParamsInt("id")
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
    }

    // Получаем параметры пагинации
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 8) // По умолчанию ограничиваем 8 похожими объявлениями
    
    // Расчет смещения для пагинации
    offset := (page - 1) * limit

    // Получаем исходное объявление
    listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
    if err != nil {
        log.Printf("Ошибка при получении объявления: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить объявление")
    }

    // Формируем запрос для расширенного поиска в OpenSearch
    // Используем multi-match для поиска по заголовку и описанию
    matchQuery := make([]map[string]interface{}, 0)

    // Добавляем поиск по названию объявления
    if listing.Title != "" {
        // Добавляем запрос на похожие названия с небольшой нечеткостью
        matchQuery = append(matchQuery, map[string]interface{}{
            "match": map[string]interface{}{
                "title": map[string]interface{}{
                    "query":     listing.Title,
                    "boost":     3.0, // Высокий вес для названия
                    "fuzziness": "AUTO",
                },
            },
        })
    }

    // Если есть описание, добавляем его в поиск
    if len(listing.Description) > 20 {
        // Берем только первые 200 символов для более точного соответствия
        descriptionExcerpt := listing.Description
        if len(descriptionExcerpt) > 200 {
            descriptionExcerpt = descriptionExcerpt[:200]
        }

        matchQuery = append(matchQuery, map[string]interface{}{
            "match": map[string]interface{}{
                "description": map[string]interface{}{
                    "query":     descriptionExcerpt,
                    "boost":     1.0,
                    "fuzziness": "AUTO",
                },
            },
        })
    }

    // Строим сложный запрос для OpenSearch
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "should": matchQuery,
                "must": []map[string]interface{}{
                    {
                        "term": map[string]interface{}{
                            "category_id": listing.CategoryID,
                        },
                    },
                    {
                        "term": map[string]interface{}{
                            "status": "active",
                        },
                    },
                },
                "must_not": []map[string]interface{}{
                    {
                        "term": map[string]interface{}{
                            "id": listingID, // Исключаем текущее объявление
                        },
                    },
                },
                "minimum_should_match": 1,
            },
        },
        "from": offset,
        "size": limit,
    }

    // Добавляем ценовой диапазон если есть цена
    if listing.Price > 0 {
        // Определяем диапазон цен (например, ±30% от текущей цены)
        minPrice := listing.Price * 0.7
        maxPrice := listing.Price * 1.3

        priceQuery := map[string]interface{}{
            "range": map[string]interface{}{
                "price": map[string]interface{}{
                    "gte":   minPrice,
                    "lte":   maxPrice,
                    "boost": 1.5, // Придаем вес объявлениям с похожей ценой
                },
            },
        }

        shouldClause := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{})
        shouldClause = append(shouldClause, priceQuery)
        query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = shouldClause
    }

    // Учитываем атрибуты для повышения релевантности
    if len(listing.Attributes) > 0 {
        for _, attr := range listing.Attributes {
            // Выбираем только значимые атрибуты
            if isSignificantAttribute(attr.AttributeName) && attr.DisplayValue != "" {
                attrQuery := map[string]interface{}{
                    "nested": map[string]interface{}{
                        "path": "attributes",
                        "query": map[string]interface{}{
                            "bool": map[string]interface{}{
                                "must": []map[string]interface{}{
                                    {
                                        "term": map[string]interface{}{
                                            "attributes.attribute_name": attr.AttributeName,
                                        },
                                    },
                                    {
                                        "match": map[string]interface{}{
                                            "attributes.display_value": map[string]interface{}{
                                                "query": attr.DisplayValue,
                                                "boost": 2.0, // Высокий вес для совпадения по атрибутам
                                            },
                                        },
                                    },
                                },
                            },
                        },
                        "boost": 2.5,
                    },
                }

                shouldClause := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{})
                shouldClause = append(shouldClause, attrQuery)
                query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = shouldClause
            }
        }
    }

    // Формируем параметры поиска для OpenSearch
    searchParams := &search.SearchParams{
        Page:        page,
        Size:        limit,
        CustomQuery: query,
    }

    // Выполняем поиск с использованием OpenSearch
    var similarListings []*models.MarketplaceListing
    result, err := h.marketplaceService.Storage().SearchListings(c.Context(), searchParams)

    if err != nil {
        log.Printf("Ошибка поиска похожих объявлений через OpenSearch: %v", err)

        // Если OpenSearch поиск не удался, используем запасной вариант
        // с простым поиском по категории
        fallbackParams := &search.ServiceParams{
            CategoryID: strconv.Itoa(listing.CategoryID),
            Size:       limit,
            Page:       page,
            Sort:       "date_desc",
        }

        fallbackResult, fallbackErr := h.marketplaceService.SearchListingsAdvanced(c.Context(), fallbackParams)
        if fallbackErr != nil {
            log.Printf("Ошибка запасного поиска: %v", fallbackErr)
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить похожие объявления")
        }

        // Фильтруем результаты, убирая исходное объявление
        for _, item := range fallbackResult.Items {
            if item.ID != listingID {
                similarListings = append(similarListings, item)
            }
        }
    } else {
        // Используем результаты из OpenSearch
        similarListings = result.Listings
    }

    log.Printf("Найдено %d похожих объявлений для объявления ID=%d (страница %d, лимит %d)", 
        len(similarListings), listingID, page, limit)

    // Возвращаем данные с поддержкой пагинации
    return utils.SuccessResponse(c, similarListings)
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

// processAttributesFromRequest обрабатывает атрибуты из запроса и добавляет их в объявление
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

					// Обработка значения в зависимости от типа
					switch attrType {
					case "text", "select":
						if value, ok := attrMap["value"].(string); ok && value != "" {
							attr.TextValue = &value
							log.Printf("DEBUG: Attribute %d (%s) text value: %s", attr.AttributeID, attr.AttributeName, value)
						}
					case "number":
						// Универсальная обработка числовых значений
						var numValue float64
						var isSet bool

						if value, ok := attrMap["value"].(float64); ok {
							numValue = value
							isSet = true
							log.Printf("DEBUG: Получено числовое значение value: %f для атрибута %s", value, attr.AttributeName)
						} else if value, ok := attrMap["numeric_value"].(float64); ok {
							numValue = value
							isSet = true
							log.Printf("DEBUG: Получено числовое значение numeric_value: %f для атрибута %s", value, attr.AttributeName)
						} else if strValue, ok := attrMap["value"].(string); ok && strValue != "" {
							// Используем ParseFloat с 64-битной точностью
							parsedValue, parseErr := strconv.ParseFloat(strValue, 64)
							if parseErr == nil {
								numValue = parsedValue
								isSet = true
								log.Printf("DEBUG: Преобразовано строковое значение '%s' в число %f для атрибута %s", strValue, parsedValue, attr.AttributeName)
							} else {
								// Пробуем удалить запятые/пробелы
								cleanValue := strings.ReplaceAll(strings.ReplaceAll(strValue, ",", ""), " ", "")
								if parsedClean, err := strconv.ParseFloat(cleanValue, 64); err == nil {
									numValue = parsedClean
									isSet = true
									log.Printf("DEBUG: Успешно преобразовано очищенное значение '%s' в число %f для атрибута %s", cleanValue, parsedClean, attr.AttributeName)
								}
							}
						}

						// Если значение установлено, обрабатываем особые случаи
						if isSet {
							// Проверка и корректировка для особых атрибутов
							if attr.AttributeName == "year" {
								currentYear := time.Now().Year()
								if numValue < 1900 || numValue > float64(currentYear+1) {
									if numValue > 0 {
										log.Printf("Keeping year value %f as is", numValue)
									} else {
										numValue = float64(currentYear)
										log.Printf("Setting year to current year: %f", numValue)
									}
								}
								// Округляем год до целого числа
								numValue = math.Floor(numValue)
							} else if attr.AttributeName == "engine_capacity" {
								// Округляем до 1 знака после запятой
								numValue = math.Round(numValue*10) / 10
							}

							attr.NumericValue = &numValue
						}
					case "boolean":
						if value, ok := attrMap["value"].(bool); ok {
							attr.BooleanValue = &value
						} else if strValue, ok := attrMap["value"].(string); ok {
							boolValue := strValue == "true" || strValue == "1"
							attr.BooleanValue = &boolValue
						}
					}
				}

				// Устанавливаем display_value для любого типа атрибута
				if attr.TextValue != nil {
					attr.DisplayValue = *attr.TextValue
				} else if attr.NumericValue != nil {
					attr.DisplayValue = fmt.Sprintf("%g", *attr.NumericValue)
				} else if attr.BooleanValue != nil {
					if *attr.BooleanValue {
						attr.DisplayValue = "Да"
					} else {
						attr.DisplayValue = "Нет"
					}
				}

				// Добавляем атрибут только если есть какое-то значение
				if attr.TextValue != nil || attr.NumericValue != nil || attr.BooleanValue != nil || attr.JSONValue != nil {
					listing.Attributes = append(listing.Attributes, attr)
				}
			}
		}
	}
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
		// Обработка файла
		log.Printf("Processing file %d: Name=%s, Size=%d, ContentType=%s", i, file.Filename, file.Size, file.Header.Get("Content-Type"))
		fileName, err := h.marketplaceService.ProcessImage(file)
		if err != nil {
			log.Printf("Failed to process image: %v", err)
			continue
		}

		// Сохраняем информацию о файле
		image := models.MarketplaceImage{
			ListingID:   listingID,
			FilePath:    fileName,
			FileName:    file.Filename,
			FileSize:    int(file.Size),
			ContentType: file.Header.Get("Content-Type"),
			IsMain:      i == mainImageIndex,
		}

		// Сохраняем файл
		err = c.SaveFile(file, "./uploads/"+fileName)
		if err != nil {
			log.Printf("Failed to save file: %v", err)
			continue
		}
		log.Printf("Image saved: %s", image.FilePath)
		// Сохраняем информацию в базу
		imageID, err := h.marketplaceService.AddListingImage(c.Context(), &image)
		if err != nil {
			log.Printf("Failed to save image info: %v", err)
			continue
		}

		image.ID = imageID
		uploadedImages = append(uploadedImages, image)
	}

	fullListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Failed to get full listing for reindexing: %v", err)
	} else {
		// Переиндексируем объявление с загруженными изображениями
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
// ReindexRatings переиндексирует рейтинги всех объявлений
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
		Type       string      `json:"type"`
		ID         interface{} `json:"id"`
		Title      string      `json:"title"`
		Display    string      `json:"display,omitempty"`
		Priority   int         `json:"priority"`
		CategoryID int         `json:"category_id,omitempty"`
		Path       interface{} `json:"path,omitempty"`
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

// Добавить новый метод
func (h *MarketplaceHandler) GetSubcategories(c *fiber.Ctx) error {
	parentID := c.Query("parent_id")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	categories, err := h.marketplaceService.GetSubcategories(c.Context(), parentID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching subcategories")
	}

	return utils.SuccessResponse(c, categories)
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

	// Создаем контекст с user_id
	ctx := context.WithValue(c.Context(), "user_id", userID)

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
		processAttributesFromRequest(requestBody, &listing)
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
	// Переиндексируем объявление в OpenSearch
	updatedListing, err := h.marketplaceService.GetListingByID(c.Context(), listing.ID)
	if err != nil {
		log.Printf("Failed to get updated listing for reindexing: %v", err)
	} else {
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), updatedListing); err != nil {
			log.Printf("Failed to reindex listing %d: %v", listing.ID, err)
		} else {
			log.Printf("Successfully reindexed listing %d", listing.ID)
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
	}

	if err := c.BodyParser(&updateData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
	}

	// Обновляем каждый переведенный field
	for fieldName, translatedText := range updateData.Translations {
		err := h.marketplaceService.UpdateTranslation(c.Context(), &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            updateData.Language,
			FieldName:           fieldName,
			TranslatedText:      translatedText,
			IsVerified:          updateData.IsVerified,
			IsMachineTranslated: false,
		})
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating translation")
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Translations updated successfully",
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
		sourceLanguage = "sr" // По умолчанию
	}

	// Переводим на каждый целевой язык
	var listingSuccess bool = false
	for _, targetLang := range request.TargetLanguages {
		// Пропускаем язык оригинала
		if targetLang == sourceLanguage {
			continue
		}

		// Переводим заголовок
		translatedTitle, err := h.services.Translation().Translate(c.Context(), listing.Title, sourceLanguage, targetLang)
		if err != nil {
			log.Printf("Error translating title for listing %d: %v", listingID, err)
			continue
		}

		// Переводим описание
		translatedDesc, err := h.services.Translation().Translate(c.Context(), listing.Description, sourceLanguage, targetLang)
		if err != nil {
			log.Printf("Error translating description for listing %d: %v", listingID, err)
			continue
		}

		// Сохраняем переводы
		translationData := &models.Translation{
			EntityType:          "listing",
			EntityID:            listingID,
			Language:            targetLang,
			FieldName:           "title",
			TranslatedText:      translatedTitle,
			IsMachineTranslated: true,
			IsVerified:          false,
		}
		err = h.marketplaceService.UpdateTranslation(c.Context(), translationData)
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
		}
		err = h.marketplaceService.UpdateTranslation(c.Context(), descTranslation)
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

    log.Printf("Извлеченные атрибуты фильтров: %+v", attributeFilters)


	log.Printf("Исходные параметры сортировки из запроса: sort_by=%s", c.Query("sort_by", ""))

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
	// Всегда устанавливаем статус "active" для всех публичных запросов
    params.Status = "active"
	// Добавьте этот код после инициализации params
	log.Printf("Итоговый параметр сортировки в запросе к OpenSearch: %s", params.Sort)
    
    // Проверяем и логируем параметры пагинации
    log.Printf("Параметры пагинации: page=%d, size=%d", params.Page, params.Size)
    log.Printf("Параметр сортировки получен из запроса: %s", params.Sort)
    // Устанавливаем разумные ограничения на параметры пагинации
    if params.Page < 1 {
        params.Page = 1
    }
    if params.Size < 1 {
        params.Size = 20
    } else if params.Size > 1000 {
        // Ограничиваем максимальный размер страницы, но увеличиваем для карты
        // Проверяем, предназначен ли запрос для отображения на карте
        viewMode := c.Query("view_mode", "")
        if viewMode == "map" {
            // Для карты разрешаем больший размер страницы
            params.Size = 1000
        } else {
            // Для обычного списка оставляем прежнее ограничение для производительности
            params.Size = 100
        }
    }
    
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
