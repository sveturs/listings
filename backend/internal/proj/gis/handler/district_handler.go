package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/proj/gis/constants"

	"backend/pkg/utils"
	"backend/internal/proj/gis/service"
	"backend/internal/proj/gis/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DistrictHandler handles district and municipality related requests
type DistrictHandler struct {
	service *service.DistrictService
}

// NewDistrictHandler creates a new district handler
func NewDistrictHandler(service *service.DistrictService) *DistrictHandler {
	return &DistrictHandler{service: service}
}

// GetDistricts returns all districts
// @Summary Get districts
// @Description Get all districts with optional filtering
// @Tags gis
// @Accept json
// @Produce json
// @Param country_code query string false "Country code (e.g., RS)"
// @Param city_id query string false "City ID"
// @Param city_ids query []string false "City IDs (comma-separated)"
// @Param name query string false "District name (partial match)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]types.District} "List of districts"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/districts [get]
func (h *DistrictHandler) GetDistricts(c *fiber.Ctx) error {
	params := types.DistrictSearchParams{
		CountryCode: c.Query("country_code", "RS"),
		Name:        c.Query("name"),
	}

	// Parse city_id if provided
	if cityIDStr := c.Query("city_id"); cityIDStr != "" {
		cityID, err := uuid.Parse(cityIDStr)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidCityId")
		}
		params.CityID = &cityID
	}

	// Parse city_ids[] if provided
	cityIDStrs := c.Context().QueryArgs().PeekMulti("city_ids[]")
	if len(cityIDStrs) > 0 {
		var cityIDs []uuid.UUID
		for _, cityIDBytes := range cityIDStrs {
			cityID, err := uuid.Parse(string(cityIDBytes))
			if err != nil {
				return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidCityId")
			}
			cityIDs = append(cityIDs, cityID)
		}
		params.CityIDs = cityIDs
	}

	districts, err := h.service.GetDistricts(c.Context(), params)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("ERROR: Failed to get districts: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToGetDistricts")
	}

	return utils.SuccessResponse(c, districts)
}

// GetDistrictByID returns a district by ID
// @Summary Get district by ID
// @Description Get a single district by its ID
// @Tags gis
// @Accept json
// @Produce json
// @Param id path string true "District ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=types.District} "District details"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "District not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/districts/{id} [get]
func (h *DistrictHandler) GetDistrictByID(c *fiber.Ctx) error {
	districtID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidDistrictId")
	}

	district, err := h.service.GetDistrictByID(c.Context(), districtID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "api.districtNotFound")
	}

	return utils.SuccessResponse(c, district)
}

// GetDistrictBoundary returns district boundary as GeoJSON
// @Summary Get district boundary
// @Description Get district boundary as GeoJSON polygon for map visualization
// @Tags gis
// @Accept json
// @Produce json
// @Param id path string true "District ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=types.DistrictBoundaryResponse} "District boundary"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "District not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/districts/{id}/boundary [get]
func (h *DistrictHandler) GetDistrictBoundary(c *fiber.Ctx) error {
	districtID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidDistrictId")
	}

	boundary, err := h.service.GetDistrictBoundary(c.Context(), districtID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "api.districtNotFound")
	}

	return utils.SuccessResponse(c, boundary)
}

// GetMunicipalities returns all municipalities
// @Summary Get municipalities
// @Description Get all municipalities with optional filtering
// @Tags gis
// @Accept json
// @Produce json
// @Param country_code query string false "Country code (e.g., RS)"
// @Param district_id query string false "District ID"
// @Param name query string false "Municipality name (partial match)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]types.Municipality} "List of municipalities"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/municipalities [get]
func (h *DistrictHandler) GetMunicipalities(c *fiber.Ctx) error {
	params := types.MunicipalitySearchParams{
		CountryCode: c.Query("country_code", "RS"),
		Name:        c.Query("name"),
	}

	// Parse district_id if provided
	if districtIDStr := c.Query("district_id"); districtIDStr != "" {
		districtID, err := uuid.Parse(districtIDStr)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidDistrictId")
		}
		params.DistrictID = &districtID
	}

	municipalities, err := h.service.GetMunicipalities(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToGetMunicipalities")
	}

	return utils.SuccessResponse(c, municipalities)
}

// GetMunicipalityByID returns a municipality by ID
// @Summary Get municipality by ID
// @Description Get a single municipality by its ID
// @Tags gis
// @Accept json
// @Produce json
// @Param id path string true "Municipality ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=types.Municipality} "Municipality details"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Municipality not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/municipalities/{id} [get]
func (h *DistrictHandler) GetMunicipalityByID(c *fiber.Ctx) error {
	municipalityID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidMunicipalityId")
	}

	municipality, err := h.service.GetMunicipalityByID(c.Context(), municipalityID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "api.municipalityNotFound")
	}

	return utils.SuccessResponse(c, municipality)
}

// SearchByDistrict searches listings within a district
// @Summary Search listings by district
// @Description Search for listings within a specific district
// @Tags gis
// @Accept json
// @Produce json
// @Param district_id path string true "District ID"
// @Param category_id query string false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param limit query int false "Limit results (default: 1000, max: 5000)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]types.GeoListing} "Search results"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "District not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/search/by-district/{district_id} [get]
func (h *DistrictHandler) SearchByDistrict(c *fiber.Ctx) error {
	districtID, err := uuid.Parse(c.Params("district_id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidDistrictId")
	}

	params := types.DistrictListingSearchParams{
		Limit:  constants.DEFAULT_LIMIT,
		Offset: 0,
	}

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidCategoryId")
		}
		params.CategoryID = &categoryID
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	results, err := h.service.SearchListingsByDistrict(c.Context(), districtID, params)
	if err != nil {
		// Temporary debug logging
		fmt.Printf("ERROR: SearchListingsByDistrict failed: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToSearchListings")
	}

	return utils.SuccessResponse(c, results)
}

// SearchByMunicipality searches listings within a municipality
// @Summary Search listings by municipality
// @Description Search for listings within a specific municipality
// @Tags gis
// @Accept json
// @Produce json
// @Param municipality_id path string true "Municipality ID"
// @Param category_id query string false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param limit query int false "Limit results (default: 1000, max: 5000)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]types.GeoListing} "Search results"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Municipality not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/search/by-municipality/{municipality_id} [get]
func (h *DistrictHandler) SearchByMunicipality(c *fiber.Ctx) error {
	municipalityID, err := uuid.Parse(c.Params("municipality_id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidMunicipalityId")
	}

	params := types.MunicipalityListingSearchParams{
		Limit:  50,
		Offset: 0,
	}

	// Parse query parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidCategoryId")
		}
		params.CategoryID = &categoryID
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	results, err := h.service.SearchListingsByMunicipality(c.Context(), municipalityID, params)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToSearchListings")
	}

	return utils.SuccessResponse(c, results)
}

// GetCities returns all cities
// @Summary Get cities
// @Description Get all cities with optional filtering by viewport bounds
// @Tags gis
// @Accept json
// @Produce json
// @Param country_code query string false "Country code (e.g., RS)"
// @Param has_districts query bool false "Filter cities with districts"
// @Param bounds_north query number false "Viewport north boundary"
// @Param bounds_south query number false "Viewport south boundary"
// @Param bounds_east query number false "Viewport east boundary"
// @Param bounds_west query number false "Viewport west boundary"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]types.City} "List of cities"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/cities [get]
func (h *DistrictHandler) GetCities(c *fiber.Ctx) error {
	params := types.CitySearchParams{
		CountryCode: c.Query("country_code", "RS"),
		SearchQuery: c.Query("q"),
		Limit:       50,
		Offset:      0,
	}

	// Parse has_districts filter
	if hasDistrictsStr := c.Query("has_districts"); hasDistrictsStr != "" {
		hasDistricts := hasDistrictsStr == "true"
		params.HasDistricts = &hasDistricts
	}

	// Parse bounds if provided
	if northStr := c.Query("bounds_north"); northStr != "" {
		var bounds types.Bounds
		var err error

		if bounds.North, err = strconv.ParseFloat(northStr, 64); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidBounds")
		}
		if bounds.South, err = strconv.ParseFloat(c.Query("bounds_south"), 64); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidBounds")
		}
		if bounds.East, err = strconv.ParseFloat(c.Query("bounds_east"), 64); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidBounds")
		}
		if bounds.West, err = strconv.ParseFloat(c.Query("bounds_west"), 64); err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidBounds")
		}

		params.Bounds = &bounds
	}

	cities, err := h.service.GetCities(c.Context(), params)
	if err != nil {
		fmt.Printf("ERROR: Failed to get cities: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToGetCities")
	}

	return utils.SuccessResponse(c, cities)
}

// GetVisibleCities returns cities visible in viewport with distance calculation
// @Summary Get visible cities
// @Description Get cities visible in viewport bounds with distance to center
// @Tags gis
// @Accept json
// @Produce json
// @Param request body types.VisibleCitiesRequest true "Viewport bounds and center"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=types.VisibleCitiesResponse} "Visible cities with distances"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/gis/cities/visible [post]
func (h *DistrictHandler) GetVisibleCities(c *fiber.Ctx) error {
	var req types.VisibleCitiesRequest

	// Логирование сырого тела запроса
	fmt.Printf("DEBUG: GetVisibleCities raw body: %s\n", string(c.Body()))

	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("ERROR: Failed to parse request body: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidRequestBody")
	}

	// Логирование распарсенного запроса
	fmt.Printf("DEBUG: Parsed request: bounds=%+v, center=%+v\n", req.Bounds, req.Center)

	if err := req.Validate(); err != nil {
		fmt.Printf("ERROR: Request validation failed: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "api.invalidRequestData")
	}

	response, err := h.service.GetVisibleCities(c.Context(), req)
	if err != nil {
		fmt.Printf("ERROR: Failed to get visible cities: %+v\n", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "api.failedToGetVisibleCities")
	}

	return utils.SuccessResponse(c, response)
}
