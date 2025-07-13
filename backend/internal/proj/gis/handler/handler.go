package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/middleware"
	"backend/internal/proj/gis/repository"
	"backend/internal/proj/gis/service"
)

// Handler основной обработчик GIS модуля
type Handler struct {
	spatialHandler         *SpatialHandler
	geocodingHandler       *GeocodingHandler
	districtHandler        *DistrictHandler
	advancedFiltersHandler *AdvancedFiltersHandler
}

// NewHandler создает новый обработчик GIS модуля
func NewHandler(db *sqlx.DB) *Handler {
	spatialService := service.NewSpatialService(db)
	geocodingService := service.NewGeocodingService(db)

	// District repository and service
	districtRepo := repository.NewDistrictRepository(db)
	districtService := service.NewDistrictService(districtRepo)

	// Advanced filters services
	isochroneService := service.NewIsochroneService(db)
	poiService := service.NewPOIService(db)
	densityService := service.NewDensityService(db)

	spatialHandler := NewSpatialHandler(spatialService)
	geocodingHandler := NewGeocodingHandler(geocodingService)
	districtHandler := NewDistrictHandler(districtService)
	advancedFiltersHandler := NewAdvancedFiltersHandler(isochroneService, poiService, densityService)

	return &Handler{
		spatialHandler:         spatialHandler,
		geocodingHandler:       geocodingHandler,
		districtHandler:        districtHandler,
		advancedFiltersHandler: advancedFiltersHandler,
	}
}

// RegisterRoutes регистрирует все маршруты GIS модуля
func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware *middleware.Middleware) error {
	// Группа маршрутов для GIS
	gis := app.Group("/api/v1/gis")

	// ========== Публичные маршруты геокодирования (Phase 2) ==========
	gis.Post("/geocode/validate", h.geocodingHandler.ValidateGeocode)
	gis.Get("/geocode/suggestions", h.geocodingHandler.SearchAddressSuggestions)
	gis.Get("/geocode/reverse", h.geocodingHandler.ReverseGeocode)
	gis.Get("/geocode/cache/stats", h.geocodingHandler.GetCacheStats)

	// ========== Публичные маршруты поиска ==========
	gis.Get("/search", h.spatialHandler.SearchListings)
	gis.Get("/search/radius", h.spatialHandler.RadiusSearch)
	gis.Get("/nearby", h.spatialHandler.GetNearbyListings)
	gis.Get("/listings/:id/location", h.spatialHandler.GetListingLocation)

	// ========== Публичные маршруты районов и муниципалитетов (Phase 3) ==========
	gis.Get("/districts", h.districtHandler.GetDistricts)
	gis.Get("/districts/:id", h.districtHandler.GetDistrictByID)
	gis.Get("/municipalities", h.districtHandler.GetMunicipalities)
	gis.Get("/municipalities/:id", h.districtHandler.GetMunicipalityByID)
	gis.Get("/search/by-district/:district_id", h.districtHandler.SearchByDistrict)
	gis.Get("/search/by-municipality/:municipality_id", h.districtHandler.SearchByMunicipality)

	// ========== Расширенные геофильтры (Phase 4) ==========
	advanced := gis.Group("/advanced")
	advanced.Post("/isochrone", h.advancedFiltersHandler.GetIsochrone)
	advanced.Get("/poi/search", h.advancedFiltersHandler.SearchPOI)
	advanced.Post("/density/analyze", h.advancedFiltersHandler.AnalyzeDensity)
	advanced.Get("/density/heatmap", h.advancedFiltersHandler.GetDensityHeatmap)
	advanced.Post("/apply-filters", h.advancedFiltersHandler.ApplyAdvancedFilters)

	// ========== Защищенные маршруты (требуют авторизации) ==========
	protected := gis.Group("/", authMiddleware.AuthRequiredJWT)

	// Старые endpoints
	protected.Put("/listings/:id/location", h.spatialHandler.UpdateListingLocation)

	// Новые endpoints Phase 2
	protected.Put("/listings/:id/address", h.spatialHandler.UpdateListingAddress)
	protected.Post("/geocode/cache/cleanup", h.geocodingHandler.CleanupExpiredCache)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (h *Handler) GetPrefix() string {
	return "gis"
}
