// backend/internal/proj/marketplace/handler/auto.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	//"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AutoHandler struct {
	services      globalService.ServicesInterface
	autoService   service.AutoServiceInterface
}

func NewAutoHandler(services globalService.ServicesInterface) *AutoHandler {
	return &AutoHandler{
		services:    services,
		autoService: services.Auto(),
	}
}

// CreateAutoListing создает новое автомобильное объявление
func (h *AutoHandler) CreateAutoListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Структура для приема данных запроса
	var request struct {
		Listing       models.MarketplaceListing `json:"listing"`
		AutoProperties models.AutoProperties    `json:"auto_properties"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные")
	}

	// Устанавливаем ID пользователя из контекста
	request.Listing.UserID = userID

	// Валидация обязательных полей листинга
	if request.Listing.Title == "" || request.Listing.Description == "" || request.Listing.Price <= 0 || request.Listing.CategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Заполните все обязательные поля объявления")
	}

	// Валидация обязательных полей автомобиля
	if request.AutoProperties.Brand == "" || request.AutoProperties.Model == "" || request.AutoProperties.Year == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Заполните все обязательные поля автомобиля (марка, модель, год)")
	}

	// Создаем автомобильное объявление
	listingID, err := h.autoService.CreateAutoListing(c.Context(), &request.Listing, &request.AutoProperties)
	if err != nil {
		log.Printf("Error creating auto listing: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при создании объявления")
	}

	// Обновляем материализованное представление
	if err := h.services.Marketplace().RefreshCategoryListingCounts(c.Context()); err != nil {
		log.Printf("Error refreshing category counts: %v", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      listingID,
		"message": "Объявление успешно создано",
	})
}

// GetAutoListingByID получает информацию об автомобильном объявлении по ID
func (h *AutoHandler) GetAutoListingByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем автомобильное объявление
	listing, err := h.autoService.GetAutoListingByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting auto listing %d: %v", id, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении объявления")
	}

	return utils.SuccessResponse(c, listing)
}

// UpdateAutoListing обновляет автомобильное объявление
func (h *AutoHandler) UpdateAutoListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Структура для приема данных запроса
	var request struct {
		Listing       models.MarketplaceListing `json:"listing"`
		AutoProperties models.AutoProperties    `json:"auto_properties"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректные данные")
	}

	// Устанавливаем ID пользователя и объявления
	request.Listing.UserID = userID
	request.Listing.ID = id
	request.AutoProperties.ListingID = id

	// Получаем текущее объявление для проверки прав доступа
	currentListing, err := h.services.Marketplace().GetListingByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем, что пользователь является владельцем объявления
	if currentListing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "У вас нет прав на редактирование этого объявления")
	}

	// Обновляем автомобильное объявление
	err = h.autoService.UpdateAutoListing(c.Context(), &request.Listing, &request.AutoProperties)
	if err != nil {
		log.Printf("Error updating auto listing: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при обновлении объявления")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Объявление успешно обновлено",
	})
}

// GetAutoListings получает список автомобильных объявлений с фильтрацией
func (h *AutoHandler) GetAutoListings(c *fiber.Ctx) error {
	// Получаем базовые фильтры из запроса
	baseFilters := map[string]string{
		"category_id":   c.Query("category_id"),
		"city":          c.Query("city"),
		"min_price":     c.Query("min_price"),
		"max_price":     c.Query("max_price"),
		"query":         c.Query("query"),
		"condition":     c.Query("condition"),
		"sort_by":       c.Query("sort_by"),
		"storefront_id": c.Query("storefront_id"),
	}

	// Получаем автомобильные фильтры из запроса
	autoFilters := &models.AutoFilter{
		Brand:        c.Query("brand"),
		Model:        c.Query("model"),
		FuelType:     c.Query("fuel_type"),
		Transmission: c.Query("transmission"),
		BodyType:     c.Query("body_type"),
		DriveType:    c.Query("drive_type"),
	}

	// Парсим числовые параметры
	if yearFrom := c.Query("year_from"); yearFrom != "" {
		if year, err := strconv.Atoi(yearFrom); err == nil {
			autoFilters.YearFrom = year
		}
	}

	if yearTo := c.Query("year_to"); yearTo != "" {
		if year, err := strconv.Atoi(yearTo); err == nil {
			autoFilters.YearTo = year
		}
	}

	if mileageFrom := c.Query("mileage_from"); mileageFrom != "" {
		if mileage, err := strconv.Atoi(mileageFrom); err == nil {
			autoFilters.MileageFrom = mileage
		}
	}

	if mileageTo := c.Query("mileage_to"); mileageTo != "" {
		if mileage, err := strconv.Atoi(mileageTo); err == nil {
			autoFilters.MileageTo = mileage
		}
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	// Получаем автомобильные объявления
	listings, total, err := h.autoService.SearchAutoListings(c.Context(), autoFilters, baseFilters, limit, offset)
	if err != nil {
		log.Printf("Error getting auto listings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching auto listings")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": listings,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetAutoConstants возвращает константы для автомобильных свойств
func (h *AutoHandler) GetAutoConstants(c *fiber.Ctx) error {
	constants := h.autoService.GetAutoConstants()
	return utils.SuccessResponse(c, constants)
}

// GetModelsByBrand возвращает список моделей для указанной марки
func (h *AutoHandler) GetModelsByBrand(c *fiber.Ctx) error {
	brand := c.Query("brand")
	if brand == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не указана марка автомобиля")
	}

	models, err := h.autoService.GetModelsByBrand(c.Context(), brand)
	if err != nil {
		log.Printf("Error getting models for brand %s: %v", brand, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении моделей")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"models": models,
	})
}

// GetAvailableBrands возвращает список доступных марок автомобилей
func (h *AutoHandler) GetAvailableBrands(c *fiber.Ctx) error {
	brands, err := h.autoService.GetAvailableBrands(c.Context())
	if err != nil {
		log.Printf("Error getting available brands: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении марок")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"brands": brands,
	})
}

// Обработчик для проверки, является ли категория автомобильной
func (h *AutoHandler) IsAutoCategory(c *fiber.Ctx) error {
    categoryID, err := strconv.Atoi(c.Query("category_id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
    }

    // Логирование запроса
    log.Printf("Проверка категории ID=%d на принадлежность к автомобильным", categoryID)

    // Явная проверка на известные автомобильные категории
    if categoryID == 2000 || categoryID == 2100 {
        log.Printf("Категория %d является автомобильной (по прямому совпадению)", categoryID)
        return utils.SuccessResponse(c, fiber.Map{
            "is_auto": true,
        })
    }

    // Проверка через сервис
    isAuto, err := h.autoService.IsAutoCategory(c.Context(), categoryID)
    if err != nil {
        log.Printf("Ошибка проверки категории %d: %v", categoryID, err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при проверке категории")
    }

    log.Printf("Результат проверки категории %d: is_auto=%v", categoryID, isAuto)
    
    return utils.SuccessResponse(c, fiber.Map{
        "is_auto": isAuto,
    })
}