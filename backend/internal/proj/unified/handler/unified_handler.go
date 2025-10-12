// backend/internal/proj/unified/handler/unified_handler.go
package handler

import (
	"net/http"

	"backend/internal/domain/models"
	"backend/internal/proj/unified/storage/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type UnifiedHandler struct {
	storage *postgres.UnifiedStorage
	log     *zerolog.Logger
}

func NewUnifiedHandler(storage *postgres.UnifiedStorage, log *zerolog.Logger) *UnifiedHandler {
	return &UnifiedHandler{
		storage: storage,
		log:     log,
	}
}

// GetUnifiedListings godoc
// @Summary Get unified listings (C2C + B2C)
// @Description Возвращает объединенный список C2C объявлений и B2C товаров
// @Tags unified
// @Accept json
// @Produce json
// @Param source_type query string false "Тип источника: all, c2c, b2c" default(all)
// @Param category_id query int false "ID категории"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param condition query string false "Состояние: new, used, refurbished"
// @Param query query string false "Текстовый поиск"
// @Param storefront_id query int false "ID витрины (только для B2C)"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} models.UnifiedListingsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/unified/listings [get]
func (h *UnifiedHandler) GetUnifiedListings(c *fiber.Ctx) error {
	// Парсинг фильтров
	filters := models.UnifiedListingsFilters{
		SourceType:   c.Query("source_type", "all"),
		CategoryID:   c.QueryInt("category_id", 0),
		MinPrice:     float64(c.QueryInt("min_price", 0)),
		MaxPrice:     float64(c.QueryInt("max_price", 0)),
		Condition:    c.Query("condition", ""),
		Query:        c.Query("query", ""),
		StorefrontID: c.QueryInt("storefront_id", 0),
		Limit:        c.QueryInt("limit", 20),
		Offset:       c.QueryInt("offset", 0),
	}

	// Валидация
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}

	if filters.SourceType != "all" &&
		filters.SourceType != "c2c" &&
		filters.SourceType != "b2c" {
		filters.SourceType = "all"
	}

	h.log.Info().
		Str("source_type", filters.SourceType).
		Int("category_id", filters.CategoryID).
		Int("limit", filters.Limit).
		Int("offset", filters.Offset).
		Msg("GetUnifiedListings request")

	// Получить listings из storage
	listings, total, err := h.storage.GetUnifiedListings(c.Context(), filters)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get unified listings")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "unified.failed_to_get_listings",
		})
	}

	response := models.UnifiedListingsResponse{
		Data:       listings,
		Total:      total,
		Limit:      filters.Limit,
		Offset:     filters.Offset,
		SourceType: filters.SourceType,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response.Data,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	})
}

// GetUnifiedListingByID godoc
// @Summary Get unified listing by ID
// @Description Получить unified listing по ID и типу (c2c или b2c)
// @Tags unified
// @Accept json
// @Produce json
// @Param id path int true "ID объявления/товара"
// @Param source_type query string true "Тип источника: c2c, b2c"
// @Success 200 {object} models.UnifiedListing
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/unified/listings/{id} [get]
func (h *UnifiedHandler) GetUnifiedListingByID(c *fiber.Ctx) error {
	// Парсинг параметров
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "unified.invalid_id",
		})
	}

	sourceType := c.Query("source_type", "")
	if sourceType != "c2c" && sourceType != "b2c" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "unified.invalid_source_type",
		})
	}

	h.log.Info().Int("id", id).Str("source_type", sourceType).Msg("GetUnifiedListingByID request")

	// Получить listing
	listing, err := h.storage.GetUnifiedListingByID(c.Context(), id, sourceType)
	if err != nil {
		h.log.Error().Err(err).Int("id", id).Msg("Failed to get unified listing by ID")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "unified.listing_not_found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    listing,
	})
}
