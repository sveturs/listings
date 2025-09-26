package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
)

// CarsHandler handles car makes and models endpoints
type CarsHandler struct {
	service    service.Interface
	carService *service.UnifiedCarService
}

// NewCarsHandler creates new cars handler
func NewCarsHandler(service service.Interface, carService *service.UnifiedCarService) *CarsHandler {
	return &CarsHandler{
		service:    service,
		carService: carService,
	}
}

// GetCarMakes godoc
// @Summary Get all car makes
// @Description Get all available car makes from the database
// @Tags marketplace-cars
// @Accept json
// @Produce json
// @Param country query string false "Filter by country"
// @Param is_domestic query boolean false "Filter domestic brands"
// @Param active_only query boolean false "Show only active brands" default(true)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CarMake} "List of car makes"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/cars/makes [get]
func (h *CarsHandler) GetCarMakes(c *fiber.Ctx) error {
	logger.Info().Msg("GetCarMakes handler called")

	// Parse query parameters
	country := c.Query("country")
	isDomestic := c.QueryBool("is_domestic", false)
	activeOnly := c.QueryBool("active_only", true)

	// Get car makes from service
	makes, err := h.service.GetCarMakes(c.Context(), country, isDomestic, false, activeOnly)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get car makes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.getError")
	}

	return utils.SuccessResponse(c, makes)
}

// GetCarModels godoc
// @Summary Get car models by make
// @Description Get all car models for a specific make
// @Tags marketplace-cars
// @Accept json
// @Produce json
// @Param make_slug path string true "Car make slug"
// @Param active_only query boolean false "Show only active models" default(true)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CarModel} "List of car models"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Make not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/cars/makes/{make_slug}/models [get]
func (h *CarsHandler) GetCarModels(c *fiber.Ctx) error {
	makeSlug := c.Params("make_slug")
	if makeSlug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.required")
	}

	activeOnly := c.QueryBool("active_only", true)

	// Get car models from service
	models, err := h.service.GetCarModelsByMake(c.Context(), makeSlug, activeOnly)
	if err != nil {
		logger.Error().Err(err).Str("make_slug", makeSlug).Msg("Failed to get car models")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.getError")
	}

	if len(models) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "general.notFound")
	}

	return utils.SuccessResponse(c, models)
}

// GetCarGenerations godoc
// @Summary Get car generations by model
// @Description Get all generations for a specific car model
// @Tags marketplace-cars
// @Accept json
// @Produce json
// @Param model_id path int true "Car model ID"
// @Param active_only query boolean false "Show only active generations" default(true)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CarGeneration} "List of car generations"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Model not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/cars/models/{model_id}/generations [get]
func (h *CarsHandler) GetCarGenerations(c *fiber.Ctx) error {
	modelIDStr := c.Params("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidNumber")
	}

	activeOnly := c.QueryBool("active_only", true)

	// Get car generations from service
	generations, err := h.service.GetCarGenerationsByModel(c.Context(), modelID, activeOnly)
	if err != nil {
		logger.Error().Err(err).Int("model_id", modelID).Msg("Failed to get car generations")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.getError")
	}

	if len(generations) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "general.notFound")
	}

	return utils.SuccessResponse(c, generations)
}

// SearchCarMakes godoc
// @Summary Search car makes
// @Description Search car makes by name with fuzzy matching
// @Tags marketplace-cars
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit results" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CarMake} "Search results"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/cars/makes/search [get]
func (h *CarsHandler) SearchCarMakes(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.required")
	}

	limit := c.QueryInt("limit", 10)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Search car makes
	makes, err := h.service.SearchCarMakes(c.Context(), query, limit)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("Failed to search car makes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.searchError")
	}

	return utils.SuccessResponse(c, makes)
}

// DecodeVIN godoc
// @Summary Decode VIN number
// @Description Decode vehicle identification number (VIN) to get vehicle information
// @Tags marketplace-cars
// @Accept json
// @Produce json
// @Param vin path string true "VIN number (17 characters)"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.VINDecodeResult} "VIN decode result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid VIN"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/cars/vin/{vin}/decode [get]
func (h *CarsHandler) DecodeVIN(c *fiber.Ctx) error {
	vin := c.Params("vin")
	if vin == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.required")
	}

	// Проверяем, что carService инициализирован
	if h.carService == nil {
		logger.Error().Msg("UnifiedCarService not initialized")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.serviceError")
	}

	// Декодируем VIN
	result, err := h.carService.DecodeVIN(c.Context(), vin)
	if err != nil {
		logger.Error().Err(err).Str("vin", vin).Msg("Failed to decode VIN")
		// Если VIN некорректный, возвращаем 400
		if err.Error() == "invalid VIN length" || err.Error() == "VIN decoder is disabled" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidVIN")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "general.decodeError")
	}

	return utils.SuccessResponse(c, result)
}

// RegisterRoutes registers car-related routes
func (h *CarsHandler) RegisterRoutes(router fiber.Router) {
	cars := router.Group("/cars")

	// Public routes
	cars.Get("/makes", h.GetCarMakes)
	cars.Get("/makes/search", h.SearchCarMakes)
	cars.Get("/makes/:make_slug/models", h.GetCarModels)
	cars.Get("/models/:model_id/generations", h.GetCarGenerations)

	// VIN декодирование (может потребовать авторизацию в будущем)
	cars.Get("/vin/:vin/decode", h.DecodeVIN)
}
