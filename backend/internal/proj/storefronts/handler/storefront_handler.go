package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// StorefrontHandler HTTP handler для витрин
type StorefrontHandler struct {
	service service.StorefrontService
}

// NewStorefrontHandler создает новый handler
func NewStorefrontHandler(service service.StorefrontService) *StorefrontHandler {
	return &StorefrontHandler{
		service: service,
	}
}

// CreateStorefront создает новую витрину
// @Summary Create new storefront
// @Description Creates a new storefront for the authenticated user
// @Tags storefronts
// @Accept json
// @Produce json
// @Param storefront body models.StorefrontCreateDTO true "Storefront data"
// @Success 201 {object} models.Storefront "Created storefront"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Storefront limit reached"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts [post]
func (h *StorefrontHandler) CreateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var dto models.StorefrontCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	storefront, err := h.service.Create(c.Context(), userID, &dto)
	if err != nil {
		switch err {
		case service.ErrStorefrontLimitReached:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Storefront limit reached for your subscription plan",
			})
		case service.ErrInvalidLocation:
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Invalid location data",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to create storefront",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(storefront)
}

// GetStorefront получает витрину по ID
// @Summary Get storefront by ID
// @Description Returns storefront details by ID
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} models.Storefront "Storefront details"
// @Failure 404 {object} ErrorResponse "Storefront not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/{id} [get]
func (h *StorefrontHandler) GetStorefront(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	storefront, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if err == postgres.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Storefront not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get storefront",
		})
	}

	return c.JSON(storefront)
}

// GetStorefrontBySlug получает витрину по slug
// @Summary Get storefront by slug
// @Description Returns storefront details by slug
// @Tags storefronts
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Success 200 {object} models.Storefront "Storefront details"
// @Failure 404 {object} ErrorResponse "Storefront not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/slug/{slug} [get]
func (h *StorefrontHandler) GetStorefrontBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	storefront, err := h.service.GetBySlug(c.Context(), slug)
	if err != nil {
		if err == postgres.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Storefront not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get storefront",
		})
	}

	return c.JSON(storefront)
}

// UpdateStorefront обновляет витрину
// @Summary Update storefront
// @Description Updates storefront details
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param storefront body models.StorefrontUpdateDTO true "Update data"
// @Success 200 {object} SuccessResponse "Storefront updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 404 {object} ErrorResponse "Storefront not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [put]
func (h *StorefrontHandler) UpdateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	var dto models.StorefrontUpdateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err = h.service.Update(c.Context(), userID, id, &dto)
	if err != nil {
		switch err {
		case service.ErrInsufficientPermissions:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		case service.ErrFeatureNotAvailable:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Feature not available in your subscription plan",
			})
		case postgres.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Storefront not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to update storefront",
			})
		}
	}

	return c.JSON(SuccessResponse{
		Message: "Storefront updated successfully",
	})
}

// DeleteStorefront удаляет витрину
// @Summary Delete storefront
// @Description Soft deletes a storefront (marks as inactive)
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} SuccessResponse "Storefront deleted"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Only owner can delete storefront"
// @Failure 404 {object} ErrorResponse "Storefront not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [delete]
func (h *StorefrontHandler) DeleteStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	err = h.service.Delete(c.Context(), userID, id)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Only owner can delete storefront",
			})
		case postgres.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Storefront not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to delete storefront",
			})
		}
	}

	return c.JSON(SuccessResponse{
		Message: "Storefront deleted successfully",
	})
}

// ListStorefronts получает список витрин
// @Summary List storefronts
// @Description Returns paginated list of storefronts with filters
// @Tags storefronts
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Param is_active query bool false "Filter by active status"
// @Param is_verified query bool false "Filter by verification status"
// @Param city query string false "Filter by city"
// @Param min_rating query float32 false "Minimum rating filter"
// @Param search query string false "Search in name and description"
// @Param lat query float64 false "Latitude for geo search"
// @Param lng query float64 false "Longitude for geo search"
// @Param radius_km query float64 false "Radius in km for geo search"
// @Param sort_by query string false "Sort field (rating, products_count, distance)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Results per page (max 100)"
// @Param offset query int false "Results offset"
// @Success 200 {object} StorefrontsListResponse "List of storefronts"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts [get]
func (h *StorefrontHandler) ListStorefronts(c *fiber.Ctx) error {
	filter := &models.StorefrontFilter{
		Limit:  20,
		Offset: 0,
	}

	// Парсим query параметры
	if userID := c.QueryInt("user_id"); userID > 0 {
		filter.UserID = &userID
	}
	
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}
	
	if isVerified := c.Query("is_verified"); isVerified != "" {
		verified := isVerified == "true"
		filter.IsVerified = &verified
	}
	
	if city := c.Query("city"); city != "" {
		filter.City = &city
	}
	
	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		filter.MinRating = &minRating
	}
	
	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}
	
	// Геофильтр
	if lat, err := strconv.ParseFloat(c.Query("lat"), 64); err == nil {
		filter.Latitude = &lat
	}
	if lng, err := strconv.ParseFloat(c.Query("lng"), 64); err == nil {
		filter.Longitude = &lng
	}
	if radius, err := strconv.ParseFloat(c.Query("radius_km"), 64); err == nil {
		filter.RadiusKm = &radius
	}
	
	// Сортировка
	filter.SortBy = c.Query("sort_by", "created_at")
	filter.SortOrder = c.Query("sort_order", "desc")
	
	// Пагинация
	if limit := c.QueryInt("limit", 20); limit > 0 && limit <= 100 {
		filter.Limit = limit
	}
	filter.Offset = c.QueryInt("offset", 0)

	storefronts, total, err := h.service.Search(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to list storefronts",
		})
	}

	return c.JSON(StorefrontsListResponse{
		Storefronts: storefronts,
		Total:       total,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	})
}

// GetMyStorefronts получает витрины текущего пользователя
// @Summary Get my storefronts
// @Description Returns list of storefronts owned by current user
// @Tags storefronts
// @Accept json
// @Produce json
// @Success 200 {object} []models.Storefront "List of user's storefronts"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/my [get]
func (h *StorefrontHandler) GetMyStorefronts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefronts, err := h.service.ListUserStorefronts(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get user storefronts",
		})
	}

	return c.JSON(storefronts)
}

// GetMapData получает данные для карты
// @Summary Get storefronts for map
// @Description Returns storefronts within map bounds with minimal data for performance
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param min_lat query float64 true "Minimum latitude"
// @Param max_lat query float64 true "Maximum latitude"
// @Param min_lng query float64 true "Minimum longitude"
// @Param max_lng query float64 true "Maximum longitude"
// @Param min_rating query float32 false "Minimum rating filter"
// @Success 200 {object} []models.StorefrontMapData "Map markers data"
// @Failure 400 {object} ErrorResponse "Invalid bounds"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/map [get]
func (h *StorefrontHandler) GetMapData(c *fiber.Ctx) error {
	// Парсим границы карты
	minLat, err1 := strconv.ParseFloat(c.Query("min_lat"), 64)
	maxLat, err2 := strconv.ParseFloat(c.Query("max_lat"), 64)
	minLng, err3 := strconv.ParseFloat(c.Query("min_lng"), 64)
	maxLng, err4 := strconv.ParseFloat(c.Query("max_lng"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid map bounds",
		})
	}

	bounds := postgres.GeoBounds{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}

	// Дополнительные фильтры
	filter := &models.StorefrontFilter{}
	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		filter.MinRating = &minRating
	}

	mapData, err := h.service.GetMapData(c.Context(), bounds, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get map data",
		})
	}

	return c.JSON(mapData)
}

// GetNearbyStorefronts получает ближайшие витрины
// @Summary Get nearby storefronts
// @Description Returns storefronts within specified radius from coordinates
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Param radius_km query float64 false "Radius in kilometers (default 5)"
// @Param limit query int false "Maximum results (default 20, max 100)"
// @Success 200 {object} []models.Storefront "Nearby storefronts"
// @Failure 400 {object} ErrorResponse "Invalid coordinates"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/nearby [get]
func (h *StorefrontHandler) GetNearbyStorefronts(c *fiber.Ctx) error {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)

	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid coordinates",
		})
	}

	radiusKm := c.QueryFloat("radius_km", 5.0)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	storefronts, err := h.service.GetNearby(c.Context(), lat, lng, radiusKm, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get nearby storefronts",
		})
	}

	return c.JSON(storefronts)
}

// GetBusinessesInBuilding получает все бизнесы в здании
// @Summary Get businesses in building
// @Description Returns all storefronts located in the same building
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param lat query float64 true "Building latitude"
// @Param lng query float64 true "Building longitude"
// @Success 200 {object} []models.StorefrontMapData "Businesses in building"
// @Failure 400 {object} ErrorResponse "Invalid coordinates"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/building [get]
func (h *StorefrontHandler) GetBusinessesInBuilding(c *fiber.Ctx) error {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)

	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid coordinates",
		})
	}

	businesses, err := h.service.GetBusinessesInBuilding(c.Context(), lat, lng)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get businesses in building",
		})
	}

	return c.JSON(businesses)
}

// UpdateWorkingHours обновляет часы работы
// @Summary Update working hours
// @Description Updates storefront working hours
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param hours body []models.StorefrontHours true "Working hours"
// @Success 200 {object} SuccessResponse "Hours updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/hours [put]
func (h *StorefrontHandler) UpdateWorkingHours(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	var hours []*models.StorefrontHours
	if err := c.BodyParser(&hours); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err = h.service.UpdateWorkingHours(c.Context(), userID, id, hours)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to update working hours",
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Working hours updated successfully",
	})
}

// UpdatePaymentMethods обновляет методы оплаты
// @Summary Update payment methods
// @Description Updates storefront payment methods
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param methods body []models.StorefrontPaymentMethod true "Payment methods"
// @Success 200 {object} SuccessResponse "Payment methods updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/payment-methods [put]
func (h *StorefrontHandler) UpdatePaymentMethods(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	var methods []*models.StorefrontPaymentMethod
	if err := c.BodyParser(&methods); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err = h.service.UpdatePaymentMethods(c.Context(), userID, id, methods)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to update payment methods",
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Payment methods updated successfully",
	})
}

// UpdateDeliveryOptions обновляет опции доставки
// @Summary Update delivery options
// @Description Updates storefront delivery options
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param options body []models.StorefrontDeliveryOption true "Delivery options"
// @Success 200 {object} SuccessResponse "Delivery options updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/delivery-options [put]
func (h *StorefrontHandler) UpdateDeliveryOptions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	var options []*models.StorefrontDeliveryOption
	if err := c.BodyParser(&options); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err = h.service.UpdateDeliveryOptions(c.Context(), userID, id, options)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to update delivery options",
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Delivery options updated successfully",
	})
}

// RecordView записывает просмотр витрины
// @Summary Record storefront view
// @Description Records a view for analytics
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} SuccessResponse "View recorded"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/{id}/view [post]
func (h *StorefrontHandler) RecordView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	err = h.service.RecordView(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to record view",
		})
	}

	return c.JSON(SuccessResponse{
		Message: "View recorded",
	})
}

// GetAnalytics получает аналитику витрины
// @Summary Get storefront analytics
// @Description Returns analytics data for specified period
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param from query string true "Start date (YYYY-MM-DD)"
// @Param to query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} []models.StorefrontAnalytics "Analytics data"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/analytics [get]
func (h *StorefrontHandler) GetAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid from date format",
		})
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid to date format",
		})
	}

	analytics, err := h.service.GetAnalytics(c.Context(), userID, id, from, to)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get analytics",
		})
	}

	return c.JSON(analytics)
}

// Response types

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type StorefrontsListResponse struct {
	Storefronts []*models.Storefront `json:"storefronts"`
	Total       int                     `json:"total"`
	Limit       int                     `json:"limit"`
	Offset      int                     `json:"offset"`
}

// UploadLogo загружает логотип витрины
// @Summary Upload storefront logo
// @Description Uploads a logo image for the storefront
// @Tags storefronts
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Storefront ID"
// @Param logo formData file true "Logo file"
// @Success 200 {object} map[string]string "Logo URL"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/logo [post]
func (h *StorefrontHandler) UploadLogo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Failed to get logo file",
		})
	}

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to open file",
		})
	}
	defer src.Close()

	// Читаем данные
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to read file",
		})
	}

	url, err := h.service.UploadLogo(c.Context(), userID, storefrontID, data, file.Filename)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to upload logo",
		})
	}

	return c.JSON(fiber.Map{
		"logo_url": url,
	})
}

// UploadBanner загружает баннер витрины
// @Summary Upload storefront banner
// @Description Uploads a banner image for the storefront
// @Tags storefronts
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Storefront ID"
// @Param banner formData file true "Banner file"
// @Success 200 {object} map[string]string "Banner URL"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/banner [post]
func (h *StorefrontHandler) UploadBanner(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	file, err := c.FormFile("banner")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Failed to get banner file",
		})
	}

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to open file",
		})
	}
	defer src.Close()

	// Читаем данные
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to read file",
		})
	}

	url, err := h.service.UploadBanner(c.Context(), userID, storefrontID, data, file.Filename)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to upload banner",
		})
	}

	return c.JSON(fiber.Map{
		"banner_url": url,
	})
}