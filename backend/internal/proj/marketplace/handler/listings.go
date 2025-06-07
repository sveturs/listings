// backend/internal/proj/marketplace/handler/listings.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
	"time"
)

// ListingsHandler обрабатывает запросы, связанные с объявлениями
type ListingsHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewListingsHandler создает новый обработчик объявлений
func NewListingsHandler(services globalService.ServicesInterface) *ListingsHandler {
	return &ListingsHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// CreateListing создает новое объявление
// @Summary Create new listing
// @Description Creates a new marketplace listing with attributes
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param body body models.MarketplaceListing true "Listing data"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{id=int,message=string}} "Listing created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 409 {object} utils.ErrorResponseSwag "marketplace.duplicateTitle"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings [post]
func (h *ListingsHandler) CreateListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Дополнительная обработка для атрибутов
	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		processAttributesFromRequest(requestBody, &listing)
	}

	listing.UserID = userID
	listing.Status = "active"
	
	// Санитизация полей для защиты от XSS
	listing.Title = utils.SanitizeText(listing.Title)
	listing.Description = utils.SanitizeText(listing.Description)
	if listing.Location != "" {
		listing.Location = utils.SanitizeText(listing.Location)
	}

	// Создаем объявление
	id, err := h.marketplaceService.CreateListing(c.Context(), &listing)
	if err != nil {
		log.Printf("Failed to create listing: %v", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return utils.ErrorResponse(c, fiber.StatusConflict, "marketplace.duplicateTitle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createError")
	}

	// Возвращаем ID созданного объявления
	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"message": "marketplace.createSuccess",
	})
}

// GetListing получает детали объявления
// @Summary Get listing details
// @Description Returns detailed information about a specific listing including attributes and images
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceListing} "Listing details"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getError"
// @Router /api/v1/marketplace/listings/{id} [get]
func (h *ListingsHandler) GetListing(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем детали объявления
	listing, err := h.marketplaceService.GetListingByID(c.Context(), id)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", id, err)
		if err.Error() == "listing not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getError")
	}

	// Делаем запрос на увеличение счетчика просмотров в горутине, чтобы не задерживать ответ
	go func(ctx context.Context, listingID int) {
		err := h.services.Storage().IncrementViewsCount(context.Background(), listingID)
		if err != nil {
			log.Printf("Failed to increment views count for listing %d: %v", listingID, err)
		}
	}(c.Context(), id)

	// Получаем ID пользователя из контекста для проверки избранного
	userID, ok := c.Locals("user_id").(int)
	if ok && userID > 0 {
		// Проверяем, находится ли объявление в избранном у пользователя
		var favorites []models.MarketplaceListing
		favorites, err = h.marketplaceService.GetUserFavorites(c.Context(), userID)
		if err == nil {
			for _, fav := range favorites {
				if fav.ID == listing.ID {
					listing.IsFavorite = true
					break
				}
			}
		}
	}

	// Возвращаем детали объявления
	return utils.SuccessResponse(c, listing)
}

// GetListings получает список объявлений
// @Summary Get listings list
// @Description Returns paginated list of listings with optional filters
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param category_id query int false "Category ID filter"
// @Param condition query string false "Condition filter (new, used, etc.)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param sort_by query string false "Sort order (price_asc, price_desc, date_desc, etc.)"
// @Param user_id query int false "User ID filter"
// @Param storefront_id query int false "Storefront ID filter"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=object{data=[]models.MarketplaceListing,meta=object{total=int,page=int,limit=int}}} "Listings list with pagination"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.listError"
// @Router /api/v1/marketplace/listings [get]
func (h *ListingsHandler) GetListings(c *fiber.Ctx) error {
	// Получаем параметры фильтрации из запроса
	query := c.Query("query")
	category := c.Query("category_id")
	condition := c.Query("condition")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	sortBy := c.Query("sort_by")
	userIDStr := c.Query("user_id")
	storefrontIDStr := c.Query("storefront_id")

	// Значения по умолчанию для пагинации
	limit := 20
	offset := 0

	// Получаем лимит и смещение из запроса
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Формируем фильтры
	filters := make(map[string]string)
	if query != "" {
		filters["query"] = query
	}
	if category != "" {
		filters["category_id"] = category
	}
	if condition != "" {
		filters["condition"] = condition
	}
	if minPrice != "" {
		filters["min_price"] = minPrice
	}
	if maxPrice != "" {
		filters["max_price"] = maxPrice
	}
	if sortBy != "" {
		filters["sort_by"] = sortBy
	}
	if userIDStr != "" {
		filters["user_id"] = userIDStr
	}
	if storefrontIDStr != "" {
		filters["storefront_id"] = storefrontIDStr
	}

	// Получаем список объявлений
	listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, limit, offset)
	if err != nil {
		log.Printf("Failed to get listings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.listError")
	}

	// Проверяем, что listings не nil
	if listings == nil {
		listings = []models.MarketplaceListing{}
	}

	// Возвращаем список объявлений с пагинацией
	return utils.SuccessResponse(c, fiber.Map{
		"data": listings,
		"meta": fiber.Map{
			"total": total,
			"page":  offset/limit + 1,
			"limit": limit,
		},
	})
}

// UpdateListing обновляет существующее объявление
// @Summary Update listing
// @Description Updates an existing marketplace listing. Only the owner can update
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param body body models.MarketplaceListing true "Updated listing data"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Listing updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id} [put]
func (h *ListingsHandler) UpdateListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем текущие данные объявления для проверки владельца
	currentListing, err := h.marketplaceService.GetListingByID(c.Context(), id)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", id, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем, является ли пользователь владельцем объявления
	if currentListing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Парсим данные из запроса
	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Дополнительная обработка для атрибутов
	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		processAttributesFromRequest(requestBody, &listing)
	}

	// Устанавливаем ID объявления и пользователя
	listing.ID = id
	listing.UserID = userID
	
	// Санитизация полей для защиты от XSS
	listing.Title = utils.SanitizeText(listing.Title)
	listing.Description = utils.SanitizeText(listing.Description)
	if listing.Location != "" {
		listing.Location = utils.SanitizeText(listing.Location)
	}

	// Обрабатываем изменение цены - если она отличается, сохраняем в историю
	if currentListing.Price != listing.Price {
		// Создаем запись в истории цен
		priceHistory := models.PriceHistoryEntry{
			ListingID:     id,
			Price:         listing.Price,
			EffectiveFrom: time.Now(),
			ChangeSource:  "manual",
		}

		err = h.services.Storage().ClosePriceHistoryEntry(c.Context(), id)
		if err != nil {
			log.Printf("Failed to close previous price history entry: %v", err)
		}

		err = h.services.Storage().AddPriceHistoryEntry(c.Context(), &priceHistory)
		if err != nil {
			log.Printf("Failed to add price history entry: %v", err)
		}

		// Проверяем, не является ли изменение цены манипуляцией
		isManipulation, err := h.services.Storage().CheckPriceManipulation(c.Context(), id)
		if err != nil {
			log.Printf("Failed to check price manipulation: %v", err)
		}

		if isManipulation {
			log.Printf("Detected price manipulation for listing %d", id)
			// Здесь можно добавить логику для обработки манипуляций с ценой
		}
	}

	// Обновляем объявление
	err = h.marketplaceService.UpdateListing(c.Context(), &listing)
	if err != nil {
		log.Printf("Failed to update listing: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.updateSuccess",
	})
}

// DeleteListing удаляет объявление
// @Summary Delete listing
// @Description Deletes a marketplace listing. Only the owner can delete
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Listing deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id} [delete]
func (h *ListingsHandler) DeleteListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Удаляем объявление
	err = h.marketplaceService.DeleteListing(c.Context(), id, userID)
	if err != nil {
		log.Printf("Failed to delete listing with ID %d: %v", id, err)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "permission") {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteError")
	}

	// Удаляем документ из OpenSearch
	go func() {
		err := h.services.Storage().DeleteListingIndex(context.Background(), fmt.Sprintf("%d", id))
		if err != nil {
			log.Printf("Failed to delete listing index for ID %d: %v", id, err)
		}
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.deleteSuccess",
	})
}

// GetPriceHistory получает историю цен для объявления
// @Summary Get listing price history
// @Description Returns price history for a specific listing
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.PriceHistoryEntry} "Price history entries"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.priceHistoryError"
// @Router /api/v1/marketplace/listings/{id}/price-history [get]
func (h *ListingsHandler) GetPriceHistory(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем историю цен
	priceHistory, err := h.marketplaceService.GetPriceHistory(c.Context(), id)
	if err != nil {
		log.Printf("Failed to get price history for listing %d: %v", id, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.priceHistoryError")
	}

	// Проверяем, что priceHistory не nil
	if priceHistory == nil {
		priceHistory = []models.PriceHistoryEntry{}
	}

	// Возвращаем историю цен
	return utils.SuccessResponse(c, priceHistory)
}

// SynchronizeDiscounts синхронизирует данные о скидках
// @Summary Synchronize discount data
// @Description Synchronizes discount data for all listings (Admin only)
// @Tags marketplace-admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Discounts synchronized successfully"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "admin.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.syncError"
// @Security BearerAuth
// @Router /api/v1/admin/sync-discounts [post]
func (h *ListingsHandler) SynchronizeDiscounts(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.checkError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin.required")
	}

	// Запускаем синхронизацию
	err = h.marketplaceService.SynchronizeDiscountData(c.Context(), 0)
	if err != nil {
		log.Printf("Failed to synchronize discount data: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.syncError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.syncSuccess",
	})
}

// isSignificantAttribute определяет, является ли атрибут значимым для поиска похожих объявлений
func isSignificantAttribute(attrName string) bool {
	// Список значимых атрибутов
	significantAttrs := map[string]bool{
		"brand":           true,
		"model":           true,
		"year":            true,
		"mileage":         true,
		"engine_capacity": true,
		"power":           true,
		"area":            true,
		"land_area":       true,
		"rooms":           true,
		"floor":           true,
		"total_floors":    true,
		"walls":           true,
		"condition":       true,
		"property_type":   true,
		"vehicle_type":    true,
		"screen_size":     true,
		"memory":          true,
	}

	return significantAttrs[attrName]
}

// processAttributesFromRequest обрабатывает атрибуты из запроса
func processAttributesFromRequest(requestBody map[string]interface{}, listing *models.MarketplaceListing) {
	// Проверяем наличие атрибутов в запросе
	if attributesRaw, ok := requestBody["attributes"]; ok {
		if attributesSlice, ok := attributesRaw.([]interface{}); ok {
			var attributes []models.ListingAttributeValue

			for _, attrRaw := range attributesSlice {
				if attrMap, ok := attrRaw.(map[string]interface{}); ok {
					var attr models.ListingAttributeValue

					// Перенос всех полей из JSON-объекта
					if id, ok := attrMap["attribute_id"].(float64); ok {
						attr.AttributeID = int(id)
					}
					if name, ok := attrMap["attribute_name"].(string); ok {
						attr.AttributeName = name
					}
					if displayName, ok := attrMap["display_name"].(string); ok {
						attr.DisplayName = displayName
					}
					if attrType, ok := attrMap["attribute_type"].(string); ok {
						attr.AttributeType = attrType
					}
					if unit, ok := attrMap["unit"].(string); ok {
						attr.Unit = unit
					}
					if displayValue, ok := attrMap["display_value"].(string); ok {
						attr.DisplayValue = displayValue
					}

					// Обрабатываем значение в зависимости от типа атрибута
					switch attr.AttributeType {
					case "text", "select":
						if textValue, ok := attrMap["text_value"].(string); ok && textValue != "" {
							attr.TextValue = &textValue
						} else if textValue, ok := attrMap["value"].(string); ok && textValue != "" {
							attr.TextValue = &textValue
						}
					case "number":
						if numValue, ok := attrMap["numeric_value"].(float64); ok {
							attr.NumericValue = &numValue
						} else if numValue, ok := attrMap["value"].(float64); ok {
							attr.NumericValue = &numValue
						} else if textValue, ok := attrMap["value"].(string); ok && textValue != "" {
							// Иногда числа приходят как строки, преобразуем
							if numVal, err := strconv.ParseFloat(textValue, 64); err == nil {
								attr.NumericValue = &numVal
							}
						}
					case "boolean":
						if boolValue, ok := attrMap["boolean_value"].(bool); ok {
							attr.BooleanValue = &boolValue
						} else if boolValue, ok := attrMap["value"].(bool); ok {
							attr.BooleanValue = &boolValue
						}
					case "multiselect":
						// Для multiselect значение хранится в JSON
						if jsonValues, ok := attrMap["json_value"]; ok {
							jsonBytes, err := json.Marshal(jsonValues)
							if err == nil {
								attr.JSONValue = jsonBytes
							}
						} else if jsonValues, ok := attrMap["value"]; ok {
							jsonBytes, err := json.Marshal(jsonValues)
							if err == nil {
								attr.JSONValue = jsonBytes
							}
						}
					}

					attributes = append(attributes, attr)
				}
			}

			listing.Attributes = attributes
		}
	}
}
