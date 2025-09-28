package handler

import (
	"context"
	"database/sql"
	"net/http"

	"backend/internal/domain/models"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetPublicRecommendations gets recommendations without authentication
// @Summary Get public recommendations
// @Description Get recommendations based on type and category without authentication
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param type query string false "Recommendation type (trending, new, similar, recommended)" default(trending)
// @Param category query string false "Category filter"
// @Param item_id query int false "Current item ID (for similar recommendations)"
// @Param limit query int false "Number of results" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceListing} "Recommendations"
// @Router /api/v1/marketplace/recommendations [get]
func (h *MarketplaceHandler) GetPublicRecommendations(c *fiber.Ctx) error {
	recType := c.Query("type", "trending")
	_ = c.Query("category", "") // Currently not used but kept for future filtering
	itemID := c.QueryInt("item_id", 0)
	limit := c.QueryInt("limit", 10)

	var listings []models.MarketplaceListing
	var err error

	// Simple recommendations based on type
	switch recType {
	case "trending":
		// Get trending items from last 7 days
		query := `
			SELECT * FROM marketplace_listings
			WHERE status = 'active'
			AND created_at > NOW() - INTERVAL '7 days'
			ORDER BY views_count DESC, created_at DESC
			LIMIT $1
		`
		listings, err = h.getListingsFromQuery(query, limit)

	case "new":
		// Get newest items
		query := `
			SELECT * FROM marketplace_listings
			WHERE status = 'active'
			ORDER BY created_at DESC
			LIMIT $1
		`
		listings, err = h.getListingsFromQuery(query, limit)

	case "similar":
		if itemID > 0 {
			// Get similar items by category and price
			query := `
				SELECT ml2.* FROM marketplace_listings ml1
				JOIN marketplace_listings ml2 ON ml2.category_id = ml1.category_id
				WHERE ml1.id = $1 AND ml2.id != $1
				AND ml2.status = 'active'
				AND ml2.price BETWEEN ml1.price * 0.7 AND ml1.price * 1.3
				ORDER BY ABS(ml2.price - ml1.price) ASC
				LIMIT $2
			`
			listings, err = h.getListingsFromQueryWithID(query, itemID, limit)
		} else {
			// Fallback to popular
			query := `
				SELECT * FROM marketplace_listings
				WHERE status = 'active'
				ORDER BY views_count DESC
				LIMIT $1
			`
			listings, err = h.getListingsFromQuery(query, limit)
		}

	default:
		// Default to popular
		query := `
			SELECT * FROM marketplace_listings
			WHERE status = 'active'
			ORDER BY views_count DESC, created_at DESC
			LIMIT $1
		`
		listings, err = h.getListingsFromQuery(query, limit)
	}

	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "error.serverError", nil)
	}

	// Note: Images are not loaded here to keep the response fast
	// Frontend should load images separately if needed

	return utils.SendSuccessResponse(c, listings, "success")
}

// Helper function to get listings from query
func (h *MarketplaceHandler) getListingsFromQuery(query string, limit int) ([]models.MarketplaceListing, error) {
	ctx := context.Background()
	rows, err := h.storage.GetPool().Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []models.MarketplaceListing
	for rows.Next() {
		var listing models.MarketplaceListing
		var metadata sql.NullString
		var addressMultilingual sql.NullString

		err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.CategoryID,
			&listing.Title, &listing.Description, &listing.Price,
			&listing.Condition, &listing.Status, &listing.Location,
			&listing.Latitude, &listing.Longitude, &listing.City,
			&listing.Country, &listing.ViewsCount, &listing.ShowOnMap,
			&listing.OriginalLanguage, &listing.CreatedAt, &listing.UpdatedAt,
			&listing.StorefrontID, &listing.ExternalID, &metadata,
			&sql.NullBool{}, // needs_reindex
			&addressMultilingual,
		)
		if err != nil {
			continue
		}

		listings = append(listings, listing)
	}

	return listings, nil
}

// Helper function to get listings from query with ID parameter
func (h *MarketplaceHandler) getListingsFromQueryWithID(query string, id, limit int) ([]models.MarketplaceListing, error) {
	ctx := context.Background()
	rows, err := h.storage.GetPool().Query(ctx, query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []models.MarketplaceListing
	for rows.Next() {
		var listing models.MarketplaceListing
		var metadata sql.NullString
		var addressMultilingual sql.NullString

		err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.CategoryID,
			&listing.Title, &listing.Description, &listing.Price,
			&listing.Condition, &listing.Status, &listing.Location,
			&listing.Latitude, &listing.Longitude, &listing.City,
			&listing.Country, &listing.ViewsCount, &listing.ShowOnMap,
			&listing.OriginalLanguage, &listing.CreatedAt, &listing.UpdatedAt,
			&listing.StorefrontID, &listing.ExternalID, &metadata,
			&sql.NullBool{}, // needs_reindex
			&addressMultilingual,
		)
		if err != nil {
			continue
		}

		listings = append(listings, listing)
	}

	return listings, nil
}
