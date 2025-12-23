package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/opensearch"
)

// AnalyticsHandler handles search analytics HTTP endpoints
type AnalyticsHandler struct {
	analyticsClient *opensearch.AnalyticsClient
	logger          zerolog.Logger
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsClient *opensearch.AnalyticsClient, logger zerolog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsClient: analyticsClient,
		logger:          logger.With().Str("component", "analytics_handler").Logger(),
	}
}

// RegisterRoutes registers analytics routes
func (h *AnalyticsHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	analytics := api.Group("/analytics")

	analytics.Get("/search", h.GetSearchAnalytics)
	analytics.Post("/click", h.TrackClick)
	analytics.Post("/conversion", h.TrackConversion)
}

// GetSearchAnalytics godoc
// @Summary Get search analytics report
// @Description Retrieves search analytics for a time period
// @Tags analytics
// @Accept json
// @Produce json
// @Param from query string false "Start date (RFC3339)" default(30 days ago)
// @Param to query string false "End date (RFC3339)" default(now)
// @Success 200 {object} opensearch.SearchAnalyticsReport
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/analytics/search [get]
func (h *AnalyticsHandler) GetSearchAnalytics(c *fiber.Ctx) error {
	// Parse date range (default: last 30 days)
	fromStr := c.Query("from", time.Now().AddDate(0, 0, -30).Format(time.RFC3339))
	toStr := c.Query("to", time.Now().Format(time.RFC3339))

	fromTime, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid 'from' date format, expected RFC3339",
		})
	}

	toTime, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid 'to' date format, expected RFC3339",
		})
	}

	// Validate date range
	if toTime.Before(fromTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "'to' date must be after 'from' date",
		})
	}

	// Get analytics report
	report, err := h.analyticsClient.GetSearchAnalytics(c.Context(), fromTime, toTime)
	if err != nil {
		h.logger.Error().
			Err(err).
			Time("from", fromTime).
			Time("to", toTime).
			Msg("failed to get search analytics")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve analytics",
		})
	}

	return c.JSON(report)
}

// TrackClickRequest represents click tracking request
type TrackClickRequest struct {
	SearchEventID string `json:"search_event_id" validate:"required"`
	ListingID     int64  `json:"listing_id" validate:"required,min=1"`
	Position      int    `json:"position" validate:"required,min=0"`
	SessionID     string `json:"session_id" validate:"required"`
	UserID        *int64 `json:"user_id,omitempty"`
}

// TrackClick godoc
// @Summary Track click on search result
// @Description Records a click event when user clicks on search result
// @Tags analytics
// @Accept json
// @Produce json
// @Param request body TrackClickRequest true "Click event data"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/analytics/click [post]
func (h *AnalyticsHandler) TrackClick(c *fiber.Ctx) error {
	var req TrackClickRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if req.SearchEventID == "" || req.ListingID == 0 || req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search_event_id, listing_id, and session_id are required",
		})
	}

	// Track click event
	event := &opensearch.ClickEvent{
		SearchEventID: req.SearchEventID,
		ListingID:     req.ListingID,
		Position:      req.Position,
		SessionID:     req.SessionID,
		UserID:        req.UserID,
	}

	if err := h.analyticsClient.TrackClick(c.Context(), event); err != nil {
		h.logger.Error().
			Err(err).
			Int64("listing_id", req.ListingID).
			Str("search_event_id", req.SearchEventID).
			Msg("failed to track click event")

		// Don't fail the request - analytics is non-critical
		// Just log the error and return success
	}

	return c.JSON(fiber.Map{"success": true})
}

// TrackConversionRequest represents conversion tracking request
type TrackConversionRequest struct {
	SearchEventID  string `json:"search_event_id" validate:"required"`
	ListingID      int64  `json:"listing_id" validate:"required,min=1"`
	ConversionType string `json:"conversion_type" validate:"required,oneof=cart purchase favorite"`
	UserID         *int64 `json:"user_id,omitempty"`
}

// TrackConversion godoc
// @Summary Track conversion event
// @Description Records a conversion (cart, purchase, favorite) from search
// @Tags analytics
// @Accept json
// @Produce json
// @Param request body TrackConversionRequest true "Conversion event data"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/analytics/conversion [post]
func (h *AnalyticsHandler) TrackConversion(c *fiber.Ctx) error {
	var req TrackConversionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if req.SearchEventID == "" || req.ListingID == 0 || req.ConversionType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search_event_id, listing_id, and conversion_type are required",
		})
	}

	// Validate conversion type
	validTypes := map[string]bool{"cart": true, "purchase": true, "favorite": true}
	if !validTypes[req.ConversionType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "conversion_type must be one of: cart, purchase, favorite",
		})
	}

	// Track conversion event
	event := &opensearch.ConversionEvent{
		SearchEventID:  req.SearchEventID,
		ListingID:      req.ListingID,
		ConversionType: req.ConversionType,
		UserID:         req.UserID,
	}

	if err := h.analyticsClient.TrackConversion(c.Context(), event); err != nil {
		h.logger.Error().
			Err(err).
			Int64("listing_id", req.ListingID).
			Str("search_event_id", req.SearchEventID).
			Str("conversion_type", req.ConversionType).
			Msg("failed to track conversion event")

		// Don't fail the request - analytics is non-critical
	}

	return c.JSON(fiber.Map{"success": true})
}
