package recommendations

import (
	"net/http"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/middleware"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	authmw "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"go.uber.org/zap"
)

// ViewStatistics represents view statistics for a listing
type ViewStatistics struct {
	ViewsCount  int     `json:"views_count"`
	UniqueUsers int     `json:"unique_users"`
	AvgDuration float64 `json:"avg_duration"`
}

// Handler handles recommendations endpoints
type Handler struct {
	db      *postgres.Database
	service *Service
	logger  *zap.Logger
}

// NewHandler creates a new recommendations handler
func NewHandler(db *postgres.Database) *Handler {
	return &Handler{
		db:      db,
		service: NewService(db),
		logger:  zap.L(),
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceListing} "Recommendations"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Success"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/recommendations/view-history [post]
func (h *Handler) AddViewHistory(c *fiber.Ctx) error {
	// Пытаемся получить userID из контекста через библиотечный helper
	userID, ok := authmw.GetUserID(c)
	if !ok || userID == 0 {
		// Если нет авторизации, пропускаем сохранение истории
		return utils.SendSuccessResponse(c, map[string]string{"status": "ok", "message": "anonymous"}, "success")
	}

	// Преобразуем в int64 для базы данных
	userID64 := int64(userID)

	var req ViewHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	// Insert view history
	_, err := h.db.GetSQLXDB().Exec(`
		INSERT INTO universal_view_history
		(user_id, listing_id, category_id, interaction_type, view_duration_seconds, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, userID64, req.ListingID, req.CategoryID, req.InteractionType, req.ViewDurationSeconds)
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceListing} "View history"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/recommendations/view-history [get]
func (h *Handler) GetViewHistory(c *fiber.Ctx) error {
	// TODO: Fix view history query - currently returns empty due to SQL mapping issues
	// The query with JOIN user_view_history causes SQLX to try mapping user_id column
	// which doesn't exist in MarketplaceListing model
	// Temporary solution: return empty array until proper fix is implemented

	var listings []models.MarketplaceListing
	return utils.SendSuccessResponse(c, listings, "success")
}

// GetViewStatistics retrieves view statistics for a listing
// @Summary Get view statistics
// @Description Get view statistics for a specific listing
// @Tags recommendations
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=internal_proj_recommendations.ViewStatistics} "View statistics"
// @Router /api/v1/recommendations/view-statistics/{listing_id} [get]
func (h *Handler) GetViewStatistics(c *fiber.Ctx) error {
	listingID, err := strconv.ParseInt(c.Params("listing_id"), 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "validation.invalidInput", nil)
	}

	var stats ViewStatistics

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
