package recommendations

import (
	"net/http"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/middleware"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler handles recommendations endpoints
type Handler struct {
	db      *postgres.Database
	service *Service
}

// NewHandler creates a new recommendations handler
func NewHandler(db *postgres.Database) *Handler {
	return &Handler{
		db:      db,
		service: NewService(db),
	}
}

// RegisterRoutes registers recommendation routes
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	recommendations := app.Group("/api/v1/recommendations")

	// Просто регистрируем endpoints без auth пока
	recommendations.Post("/", h.GetRecommendations)
	recommendations.Post("/view-history", h.AddViewHistory)
	recommendations.Get("/view-history", h.GetViewHistory)

	recommendations.Get("/view-statistics/:listing_id", h.GetViewStatistics)
	return nil
}

// GetPrefix returns the API prefix for recommendation endpoints
func (h *Handler) GetPrefix() string {
	return "/api/v1/recommendations"
}

// RecommendationRequest represents the request for recommendations
type RecommendationRequest struct {
	Type          string `json:"type" validate:"required,oneof=similar personal trending new discounted recommended"`
	Category      string `json:"category" validate:"required"`
	CurrentItemID int64  `json:"current_item_id,omitempty"`
	UserID        int64  `json:"user_id,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

// GetRecommendations retrieves recommendations based on type
// @Summary Get recommendations
// @Description Get recommendations based on type and category
// @Tags recommendations
// @Accept json
// @Produce json
// @Param request body RecommendationRequest true "Recommendation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceListing} "Recommendations"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Router /api/v1/recommendations [post]
func (h *Handler) GetRecommendations(c *fiber.Ctx) error {
	var req RecommendationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	// Default limit
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Get recommendations based on type using advanced algorithms
	var listings []models.MarketplaceListing
	var err error

	switch req.Type {
	case "similar":
		// Use content-based filtering for similar items
		if req.CurrentItemID > 0 {
			listings, err = h.service.GetContentBasedRecommendations(req.CurrentItemID, req.Limit)
		} else {
			listings, err = h.service.GetPopularRecommendations(req.Limit)
		}
	case "personal":
		// Use collaborative filtering for personalized recommendations
		if req.UserID > 0 {
			listings, err = h.service.GetPersonalizedRecommendations(req.UserID, req.Category, req.Limit)
		} else {
			listings, err = h.service.GetPopularRecommendations(req.Limit)
		}
	case "trending":
		// Get trending items based on recent engagement
		listings, err = h.service.GetTrendingRecommendations(req.Limit)
	case "new":
		// Keep simple new listings
		listings, err = h.getNewListings(req.Category, req.Limit)
	case "recommended":
		// Use hybrid approach for best recommendations
		listings, err = h.service.GetHybridRecommendations(req.UserID, &req.CurrentItemID, req.Limit)
	default:
		// Default to trending
		listings, err = h.service.GetTrendingRecommendations(req.Limit)
	}

	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "error.serverError", nil)
	}

	return utils.SendSuccessResponse(c, listings, "success")
}

func (h *Handler) getSimilarListings(itemID int64, category string, limit int) ([]models.MarketplaceListing, error) {
	// Get the current item details
	var currentItem models.MarketplaceListing
	err := h.db.GetSQLXDB().Get(&currentItem, `
		SELECT * FROM marketplace_listings
		WHERE id = $1 AND status = 'active'
	`, itemID)
	if err != nil {
		return nil, err
	}

	// Find similar items based on price range and category
	var listings []models.MarketplaceListing
	priceMin := currentItem.Price * 0.7
	priceMax := currentItem.Price * 1.3

	query := `
		SELECT * FROM marketplace_listings
		WHERE category_id = $1
		AND id != $2
		AND status = 'active'
		AND price BETWEEN $3 AND $4
		ORDER BY ABS(price - $5) ASC
		LIMIT $6
	`

	err = h.db.GetSQLXDB().Select(&listings, query,
		currentItem.CategoryID, itemID, priceMin, priceMax, currentItem.Price, limit)

	return listings, err
}

func (h *Handler) getTrendingListings(category string, limit int) ([]models.MarketplaceListing, error) {
	var listings []models.MarketplaceListing

	// Get listings with most views in last 7 days
	query := `
		SELECT ml.* FROM marketplace_listings ml
		WHERE ml.status = 'active'
		AND ml.created_at > NOW() - INTERVAL '7 days'
		ORDER BY ml.views DESC, ml.created_at DESC
		LIMIT $1
	`

	err := h.db.GetSQLXDB().Select(&listings, query, limit)
	return listings, err
}

func (h *Handler) getNewListings(category string, limit int) ([]models.MarketplaceListing, error) {
	var listings []models.MarketplaceListing

	query := `
		SELECT * FROM marketplace_listings
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT $1
	`

	err := h.db.GetSQLXDB().Select(&listings, query, limit)
	return listings, err
}

func (h *Handler) getRecommendedListings(category string, userID int64, limit int) ([]models.MarketplaceListing, error) {
	// For now, return featured listings
	var listings []models.MarketplaceListing

	query := `
		SELECT * FROM marketplace_listings
		WHERE status = 'active'
		AND is_featured = true
		ORDER BY created_at DESC
		LIMIT $1
	`

	err := h.db.GetSQLXDB().Select(&listings, query, limit)
	return listings, err
}

// ViewHistoryRequest represents view history entry
type ViewHistoryRequest struct {
	ListingID           int64  `json:"listing_id" validate:"required"`
	CategoryID          int64  `json:"category_id" validate:"required"`
	InteractionType     string `json:"interaction_type" validate:"required,oneof=view click_phone add_favorite"`
	ViewDurationSeconds int    `json:"view_duration_seconds,omitempty"`
}

// AddViewHistory adds a view history entry
// @Summary Add view history
// @Description Add a view history entry for the current user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ViewHistoryRequest true "View history request"
// @Success 200 {object} utils.SuccessResponseSwag "Success"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/recommendations/view-history [post]
func (h *Handler) AddViewHistory(c *fiber.Ctx) error {
	// Пытаемся получить userID из контекста
	userID, ok := c.Locals("userID").(int64)
	if !ok || userID == 0 {
		// Если нет авторизации, пропускаем сохранение истории
		return utils.SendSuccessResponse(c, map[string]string{"status": "ok", "message": "anonymous"}, "success")
	}

	var req ViewHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	// Insert view history
	_, err := h.db.GetSQLXDB().Exec(`
		INSERT INTO universal_view_history
		(user_id, listing_id, category_id, interaction_type, view_duration_seconds, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, userID, req.ListingID, req.CategoryID, req.InteractionType, req.ViewDurationSeconds)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "error.serverError", nil)
	}

	// Update listing views count
	if req.InteractionType == "view" {
		_, _ = h.db.GetSQLXDB().Exec(`
			UPDATE marketplace_listings
			SET views = views + 1
			WHERE id = $1
		`, req.ListingID)
	}

	return utils.SendSuccessResponse(c, map[string]string{"status": "ok"}, "success")
}

// GetViewHistory retrieves view history for the current user
// @Summary Get view history
// @Description Get view history for the current user
// @Tags recommendations
// @Produce json
// @Security Bearer
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceListing} "View history"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/recommendations/view-history [get]
func (h *Handler) GetViewHistory(c *fiber.Ctx) error {
	// Пытаемся получить userID из контекста
	userID, ok := c.Locals("userID").(int64)
	if !ok || userID == 0 {
		// Если нет авторизации, возвращаем пустой список
		return utils.SendSuccessResponse(c, []models.MarketplaceListing{}, "success")
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	var listings []models.MarketplaceListing
	query := `
		SELECT DISTINCT ml.* FROM marketplace_listings ml
		JOIN universal_view_history vh ON vh.listing_id = ml.id
		WHERE vh.user_id = $1
		ORDER BY vh.created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := h.db.GetSQLXDB().Select(&listings, query, userID, limit, offset)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "error.serverError", nil)
	}

	return utils.SendSuccessResponse(c, listings, "success")
}

// GetViewStatistics retrieves view statistics for a listing
// @Summary Get view statistics
// @Description Get view statistics for a specific listing
// @Tags recommendations
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=ViewStatistics} "View statistics"
// @Router /api/v1/recommendations/view-statistics/{listing_id} [get]
func (h *Handler) GetViewStatistics(c *fiber.Ctx) error {
	listingID, err := strconv.ParseInt(c.Params("listing_id"), 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	var stats struct {
		ViewsCount  int     `json:"views_count"`
		UniqueUsers int     `json:"unique_users"`
		AvgDuration float64 `json:"avg_duration"`
	}

	err = h.db.GetSQLXDB().Get(&stats, `
		SELECT
			COUNT(*) as views_count,
			COUNT(DISTINCT user_id) as unique_users,
			COALESCE(AVG(view_duration_seconds), 0) as avg_duration
		FROM universal_view_history
		WHERE listing_id = $1 AND interaction_type = 'view'
	`, listingID)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "error.serverError", nil)
	}

	return utils.SendSuccessResponse(c, stats, "success")
}
