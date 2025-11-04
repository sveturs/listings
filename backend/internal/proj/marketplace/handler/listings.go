// TEMPORARY: Will be moved to microservice
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/pkg/utils"
)

// CreateListingRequest represents the request body for creating a listing
type CreateListingRequest struct {
	CategoryID   int     `json:"category_id" validate:"required"`
	Title        string  `json:"title" validate:"required,min=3,max=200"`
	Description  *string `json:"description,omitempty"`
	Price        float64 `json:"price" validate:"required,min=0"`
	Currency     string  `json:"currency,omitempty"`
	Quantity     int32   `json:"quantity,omitempty"`
	SKU          *string `json:"sku,omitempty"`
	StorefrontID *int    `json:"storefront_id,omitempty"`
}

// CreateListing godoc
// @Summary Создать новое объявление
// @Description Создать новое объявление в marketplace
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateListingRequest true "Данные объявления"
// @Success 201 {object} utils.SuccessResponseSwag{data=interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings [post]
func (h *Handler) CreateListing(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("CreateListing: user not authenticated")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Parse request body
	var req CreateListingRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Err(err).Msg("CreateListing: failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_request")
	}

	// Basic validation
	if req.CategoryID == 0 {
		h.logger.Error().Msg("CreateListing: category_id is required")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.category_required")
	}
	if req.Title == "" || len(req.Title) < 3 || len(req.Title) > 200 {
		h.logger.Error().Msg("CreateListing: title must be between 3 and 200 characters")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_title")
	}
	if req.Price < 0 {
		h.logger.Error().Msg("CreateListing: price must be non-negative")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_price")
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "RSD"
	}

	// Set default quantity if not provided
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	// Call the storage layer (direct DB insert for now)
	// TEMPORARY: Direct DB insert until microservice fully migrated
	listing, err := h.storage.CreateListing(c.Context(), int(userID), req.CategoryID, req.Title, req.Description, req.Price, req.Currency, req.Quantity, req.SKU, req.StorefrontID)
	if err != nil {
		h.logger.Error().Err(err).
			Int("user_id", int(userID)).
			Int("category_id", req.CategoryID).
			Msg("CreateListing: failed to create listing")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.create_failed")
	}

	h.logger.Info().
		Int("listing_id", listing.ID).
		Int("user_id", int(userID)).
		Str("title", req.Title).
		Msg("Listing created successfully")

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, listing)
}

// GetListing godoc
// @Summary Получить объявление по ID
// @Description Получить детали объявления по его ID
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=interface{}}
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id} [get]
func (h *Handler) GetListing(c *fiber.Ctx) error {
	// Parse listing ID
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("GetListing: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Get listing from storage
	listing, err := h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetListing: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	return utils.SuccessResponse(c, listing)
}
