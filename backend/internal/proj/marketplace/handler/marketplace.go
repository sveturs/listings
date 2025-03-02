// backend/internal/handlers/marketplace.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"math"
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

	// Устанавливаем ID пользователя из контекста
	listing.UserID = userID

	// Валидация обязательных полей
	if listing.Title == "" || listing.Description == "" || listing.Price <= 0 || listing.CategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Заполните все обязательные поля")
	}

	// Создаем объявление
	listingID, err := h.marketplaceService.CreateListing(c.Context(), &listing)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при создании объявления")
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

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Изображения успешно загружены",
		"images":  uploadedImages,
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

    return utils.SuccessResponse(c, listing)
}

// UpdateListing - обновление объявления
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

	// Обновляем объявление
	err = h.marketplaceService.UpdateListing(c.Context(), &listing)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating listing")
	}

	// Проверяем изменение цены
	if oldListing.Price != listing.Price {
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

// DeleteListing - удаление объявления
func (h *MarketplaceHandler) DeleteListing(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
	}

	err = h.marketplaceService.DeleteListing(c.Context(), listingID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting listing")
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
    
    if query == "" {
        return utils.SuccessResponse(c, fiber.Map{
            "data": []interface{}{},
        })
    }
    
    // Получаем предложения категорий из сервиса
    categories, err := h.marketplaceService.GetCategorySuggestions(c.Context(), query, size)
    if err != nil {
        log.Printf("Ошибка получения предложений категорий: %v", err)
        return utils.SuccessResponse(c, fiber.Map{
            "data": []interface{}{},
        })
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "data": categories,
    })
}
func (h *MarketplaceHandler) SearchListingsAdvanced(c *fiber.Ctx) error {
	// Получаем параметры поиска
	params := &search.ServiceParams{
		Query:         c.Query("q", ""), 
		CategoryID:    c.Query("category_id", ""),
		Condition:     c.Query("condition", ""),
		City:          c.Query("city", ""),
		Country:       c.Query("country", ""),
		StorefrontID:  c.Query("storefront_id", ""),
		Sort:          c.Query("sort_by", ""),
		SortDirection: c.Query("sort_direction", "desc"),
		Distance:      c.Query("distance", "10km"),
		Page:          c.QueryInt("page", 1),
		Size:          c.QueryInt("size", 20),
		Language:      c.Query("language", ""),
	}
    log.Printf("Полученный поисковый запрос: %s", params.Query)
	// Обрабатываем числовые параметры - добавим защиту от ошибок
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

	// Обрабатываем координаты - добавим защиту от ошибок
	if lat := c.Query("lat", ""); lat != "" {
		if val, err := strconv.ParseFloat(lat, 64); err == nil {
			params.Latitude = val
		}
	}

	if lon := c.Query("lon", ""); lon != "" {
		if val, err := strconv.ParseFloat(lon, 64); err == nil {
			params.Longitude = val
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

	// В случае ошибки с OpenSearch, используем обычный поиск
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

		// Пробуем получить обычным методом
		listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, params.Size, (params.Page-1)*params.Size)
		if err != nil {
			log.Printf("Ошибка стандартного поиска: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка выполнения поиска")
		}

		// Формируем такой же ответ, как от OpenSearch
		return utils.SuccessResponse(c, fiber.Map{
			"data": listings,
			"meta": fiber.Map{
				"total":       total,
				"page":        params.Page,
				"size":        params.Size,
				"total_pages": (total + int64(params.Size) - 1) / int64(params.Size),
			},
		})
	}

	// Если OpenSearch ответил успешно
	return utils.SuccessResponse(c, fiber.Map{
		"data": result.Items,
		"meta": fiber.Map{
			"total":       result.Total,
			"page":        result.Page,
			"size":        result.Size,
			"total_pages": result.TotalPages,
			"facets":      result.Facets,
			"suggestions": result.Suggestions,
			"took_ms":     result.Took,
		},
	})
}

// GetSuggestions возвращает предложения автодополнения
func (h *MarketplaceHandler) GetSuggestions(c *fiber.Ctx) error {
    prefix := c.Query("q", "")
    size := c.QueryInt("size", 5)
    
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
        
        return utils.SuccessResponse(c, fiber.Map{
            "data": titles,
        })
    }
    
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
