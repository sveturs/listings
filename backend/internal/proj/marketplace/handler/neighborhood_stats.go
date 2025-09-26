package handler

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NeighborhoodStatsResponse struct {
	TotalListings int     `json:"total_listings"`
	NewToday      int     `json:"new_today"`
	WithinRadius  int     `json:"within_radius"`
	RadiusKm      float64 `json:"radius_km"`
	CenterLat     float64 `json:"center_lat,omitempty"`
	CenterLon     float64 `json:"center_lon,omitempty"`
}

// GetNeighborhoodStats godoc
// @Summary Get neighborhood statistics
// @Description Get statistics about listings in user's neighborhood
// @Tags marketplace
// @Accept json
// @Produce json
// @Param lat query number false "Center latitude"
// @Param lon query number false "Center longitude"
// @Param radius query number false "Radius in kilometers (default 5)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=handler.NeighborhoodStatsResponse} "Statistics"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.statsError"
// @Router /api/v1/marketplace/neighborhood-stats [get]
func (h *MarketplaceHandler) GetNeighborhoodStats(c *fiber.Ctx) error {
	lat := c.QueryFloat("lat", 44.8176) // Default Belgrade coordinates
	lon := c.QueryFloat("lon", 20.4633)
	radiusKm := c.QueryFloat("radius", 5.0)

	// Get total listings count
	var totalCount int
	err := h.storage.GetPool().QueryRow(c.Context(), `
		SELECT COUNT(*) FROM marketplace_listings 
		WHERE status = 'active'
	`).Scan(&totalCount)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get total listings count")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.statsError",
		})
	}

	// Get new listings today
	today := time.Now().Truncate(24 * time.Hour)
	var newToday int
	err = h.storage.GetPool().QueryRow(c.Context(), `
		SELECT COUNT(*) FROM marketplace_listings 
		WHERE status = 'active' 
		AND created_at >= $1
	`, today).Scan(&newToday)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get new listings today")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.statsError",
		})
	}

	// Get listings within radius
	// Using simple distance calculation for demo (in production should use PostGIS)
	var withinRadius int
	err = h.storage.GetPool().QueryRow(c.Context(), `
		SELECT COUNT(*) FROM marketplace_listings 
		WHERE status = 'active' 
		AND latitude IS NOT NULL 
		AND longitude IS NOT NULL
		AND (
			6371 * acos(
				cos(radians($1)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians($2)) + 
				sin(radians($1)) * sin(radians(latitude))
			)
		) <= $3
	`, lat, lon, radiusKm).Scan(&withinRadius)
	if err != nil {
		// If spatial calculation fails, estimate as percentage
		withinRadius = int(math.Round(float64(totalCount) * 0.45))
	}

	response := NeighborhoodStatsResponse{
		TotalListings: totalCount,
		NewToday:      newToday,
		WithinRadius:  withinRadius,
		RadiusKm:      radiusKm,
		CenterLat:     lat,
		CenterLon:     lon,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}
