// TEMPORARY: Will be moved to microservice
package handler

import (
	"github.com/gofiber/fiber/v2"
)

// NeighborhoodStatsResponse структура ответа для статистики окрестностей
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
// @Success 200 {object} utils.SuccessResponseSwag{data=NeighborhoodStatsResponse} "Statistics"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.statsError"
// @Router /api/v1/marketplace/neighborhood-stats [get]
func (h *Handler) GetNeighborhoodStats(c *fiber.Ctx) error {
	lat := c.QueryFloat("lat", 44.8176) // Default Belgrade coordinates
	lon := c.QueryFloat("lon", 20.4633)
	radiusKm := c.QueryFloat("radius", 5.0)

	// TODO: Реализовать реальную статистику через базу данных или микросервис
	// Пока возвращаем заглушку с нулевыми значениями
	h.logger.Info().
		Float64("lat", lat).
		Float64("lon", lon).
		Float64("radius_km", radiusKm).
		Msg("GetNeighborhoodStats called (stub implementation)")

	response := NeighborhoodStatsResponse{
		TotalListings: 0,
		NewToday:      0,
		WithinRadius:  0,
		RadiusKm:      radiusKm,
		CenterLat:     lat,
		CenterLon:     lon,
	}

	return c.JSON(response)
}
