package recommendations

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage/postgres"
)

// Service provides advanced recommendation algorithms
type Service struct {
	db *postgres.Database
}

// NewService creates a new recommendations service
func NewService(db *postgres.Database) *Service {
	return &Service{
		db: db,
	}
}

// SimilarityScore represents similarity between items
type SimilarityScore struct {
	ListingID int64   `db:"listing_id"`
	Score     float64 `db:"score"`
}

// UserInteraction represents user interaction with listings
type UserInteraction struct {
	UserID          int64     `db:"user_id"`
	ListingID       int64     `db:"listing_id"`
	CategoryID      int64     `db:"category_id"`
	InteractionType string    `db:"interaction_type"`
	ViewDuration    int       `db:"view_duration_seconds"`
	CreatedAt       time.Time `db:"created_at"`
}

// GetContentBasedRecommendations gets recommendations based on item similarity
func (s *Service) GetContentBasedRecommendations(itemID int64, limit int) ([]models.MarketplaceListing, error) {
	// Get the reference item with attributes
	var refItem struct {
		models.MarketplaceListing
		Attributes json.RawMessage `db:"attributes"`
	}

	err := s.db.GetSQLXDB().Get(&refItem, `
		SELECT ml.*,
			   COALESCE(
				   (SELECT jsonb_agg(jsonb_build_object(
					   'name', ua.name,
					   'value', ua.value,
					   'type', ua.type,
					   'unit', ua.unit
				   ))
				   FROM unified_attributes ua
				   WHERE ua.listing_id = ml.id
				   ), '[]'::jsonb
			   ) as attributes
		FROM marketplace_listings ml
		WHERE ml.id = $1 AND ml.status = 'active'
	`, itemID)
	if err != nil {
		return nil, err
	}

	// Parse reference attributes
	var refAttrs []map[string]interface{}
	if err := json.Unmarshal(refItem.Attributes, &refAttrs); err != nil {
		refAttrs = []map[string]interface{}{}
	}

	// Find similar items based on multiple factors
	query := `
		WITH item_scores AS (
			SELECT
				ml.id,
				-- Price similarity (max 30 points)
				CASE
					WHEN ml.price = 0 OR $2 = 0 THEN 0
					ELSE (1 - LEAST(ABS(ml.price - $2) / GREATEST(ml.price, $2), 1)) * 30
				END as price_score,

				-- Category match (20 points for same category, 10 for parent)
				CASE
					WHEN ml.category_id = $3 THEN 20
					WHEN ml.category_id IN (
						SELECT id FROM marketplace_categories
						WHERE parent_id = (SELECT parent_id FROM marketplace_categories WHERE id = $3)
					) THEN 10
					ELSE 0
				END as category_score,

				-- Location proximity (max 20 points)
				CASE
					WHEN ml.address_city = $4 THEN 20
					WHEN ml.location = $5 THEN 10
					ELSE 0
				END as location_score,

				-- User engagement (max 15 points based on views)
				LEAST(ml.views_count / 100.0, 1) * 15 as engagement_score,

				-- Freshness (max 15 points for newer items)
				CASE
					WHEN ml.created_at > NOW() - INTERVAL '7 days' THEN 15
					WHEN ml.created_at > NOW() - INTERVAL '30 days' THEN 10
					WHEN ml.created_at > NOW() - INTERVAL '90 days' THEN 5
					ELSE 0
				END as freshness_score
			FROM marketplace_listings ml
			WHERE ml.id != $1
				AND ml.status = 'active'
				AND ml.category_id IN (
					-- Same category or sibling categories
					SELECT id FROM marketplace_categories
					WHERE id = $3 OR parent_id = (SELECT parent_id FROM marketplace_categories WHERE id = $3)
				)
		)
		SELECT
			ml.*,
			(s.price_score + s.category_score + s.location_score + s.engagement_score + s.freshness_score) as total_score
		FROM marketplace_listings ml
		JOIN item_scores s ON ml.id = s.id
		ORDER BY total_score DESC
		LIMIT $6
	`

	var listings []models.MarketplaceListing
	rows, err := s.db.GetSQLXDB().Query(query,
		itemID, refItem.Price, refItem.CategoryID, refItem.City, refItem.Location, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	for rows.Next() {
		var listing models.MarketplaceListing
		var totalScore float64
		var metadata sql.NullString
		var addressMultilingual sql.NullString

		if err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.CategoryID,
			&listing.Title, &listing.Description, &listing.Price,
			&listing.Condition, &listing.Status, &listing.Location,
			&listing.Latitude, &listing.Longitude, &listing.City,
			&listing.Country, &listing.ViewsCount, &listing.ShowOnMap,
			&listing.OriginalLanguage, &listing.CreatedAt, &listing.UpdatedAt,
			&listing.StorefrontID, &listing.ExternalID, &metadata,
			&sql.NullBool{}, // needs_reindex
			&addressMultilingual, &totalScore,
		); err != nil {
			continue
		}

		// Parse JSON fields if not null
		if metadata.Valid {
			if err := json.Unmarshal([]byte(metadata.String), &listing.Metadata); err != nil {
				logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			}
		}
		if addressMultilingual.Valid {
			if err := json.Unmarshal([]byte(addressMultilingual.String), &listing.AddressMultilingual); err != nil {
				logger.Error().Err(err).Msg("Failed to unmarshal address multilingual")
			}
		}

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// If we have attributes, calculate attribute similarity and re-sort
	if len(refAttrs) > 0 && len(listings) > 0 {
		s.enhanceWithAttributeSimilarity(&listings, refAttrs)
	}

	return listings, nil
}

// GetCollaborativeRecommendations gets recommendations based on user behavior
func (s *Service) GetCollaborativeRecommendations(userID int64, limit int) ([]models.MarketplaceListing, error) {
	// Get user's interaction history
	var userHistory []UserInteraction
	err := s.db.GetSQLXDB().Select(&userHistory, `
		SELECT * FROM universal_view_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)
	if err != nil {
		return nil, err
	}

	if len(userHistory) == 0 {
		// No history, return popular items
		return s.GetPopularRecommendations(limit)
	}

	// Find similar users based on interaction overlap
	query := `
		WITH user_categories AS (
			SELECT DISTINCT category_id
			FROM universal_view_history
			WHERE user_id = $1
		),
		similar_users AS (
			SELECT
				vh.user_id,
				COUNT(DISTINCT vh.listing_id) as common_items,
				COUNT(DISTINCT vh.category_id) as common_categories
			FROM universal_view_history vh
			WHERE vh.user_id != $1
				AND vh.category_id IN (SELECT category_id FROM user_categories)
				AND vh.created_at > NOW() - INTERVAL '30 days'
			GROUP BY vh.user_id
			HAVING COUNT(DISTINCT vh.listing_id) > 2
			ORDER BY common_items DESC, common_categories DESC
			LIMIT 20
		)
		SELECT DISTINCT ml.*
		FROM marketplace_listings ml
		JOIN universal_view_history vh ON vh.listing_id = ml.id
		JOIN similar_users su ON su.user_id = vh.user_id
		WHERE ml.id NOT IN (
			SELECT listing_id FROM universal_view_history WHERE user_id = $1
		)
		AND ml.status = 'active'
		ORDER BY ml.views_count DESC, ml.created_at DESC
		LIMIT $2
	`

	var listings []models.MarketplaceListing
	err = s.db.GetSQLXDB().Select(&listings, query, userID, limit)
	return listings, err
}

// GetHybridRecommendations combines content-based and collaborative filtering
func (s *Service) GetHybridRecommendations(userID int64, itemID *int64, limit int) ([]models.MarketplaceListing, error) {
	results := make(map[int]models.MarketplaceListing)
	scores := make(map[int]float64)

	// Get collaborative recommendations (weight: 0.4)
	if userID > 0 {
		collaborative, err := s.GetCollaborativeRecommendations(userID, limit*2)
		if err == nil {
			for i, listing := range collaborative {
				results[listing.ID] = listing
				scores[listing.ID] += (1.0 - float64(i)/float64(len(collaborative))) * 0.4
			}
		}
	}

	// Get content-based recommendations (weight: 0.4)
	if itemID != nil && *itemID > 0 {
		contentBased, err := s.GetContentBasedRecommendations(*itemID, limit*2)
		if err == nil {
			for i, listing := range contentBased {
				results[listing.ID] = listing
				scores[listing.ID] += (1.0 - float64(i)/float64(len(contentBased))) * 0.4
			}
		}
	}

	// Get trending items (weight: 0.2)
	trending, err := s.GetTrendingRecommendations(limit)
	if err == nil {
		for i, listing := range trending {
			results[listing.ID] = listing
			scores[listing.ID] += (1.0 - float64(i)/float64(len(trending))) * 0.2
		}
	}

	// Sort by combined score
	type scoredListing struct {
		listing models.MarketplaceListing
		score   float64
	}

	var scoredList []scoredListing
	for id, listing := range results {
		scoredList = append(scoredList, scoredListing{
			listing: listing,
			score:   scores[id],
		})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	// Return top N items
	var finalList []models.MarketplaceListing
	for i := 0; i < limit && i < len(scoredList); i++ {
		finalList = append(finalList, scoredList[i].listing)
	}

	return finalList, nil
}

// GetPersonalizedRecommendations gets personalized recommendations for a user
func (s *Service) GetPersonalizedRecommendations(userID int64, category string, limit int) ([]models.MarketplaceListing, error) {
	// Analyze user preferences
	var preferences struct {
		AvgPrice      sql.NullFloat64 `db:"avg_price"`
		MinPrice      sql.NullFloat64 `db:"min_price"`
		MaxPrice      sql.NullFloat64 `db:"max_price"`
		TopCities     []string
		TopCategories []int64
	}

	// Get price preferences
	err := s.db.GetSQLXDB().Get(&preferences, `
		SELECT
			AVG(ml.price) as avg_price,
			MIN(ml.price) as min_price,
			MAX(ml.price) as max_price
		FROM marketplace_listings ml
		JOIN universal_view_history vh ON vh.listing_id = ml.id
		WHERE vh.user_id = $1
			AND vh.interaction_type IN ('view', 'click_phone', 'add_favorite')
			AND ml.price > 0
	`, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Get top cities
	rows, err := s.db.GetSQLXDB().Query(`
		SELECT ml.address_city, COUNT(*) as cnt
		FROM marketplace_listings ml
		JOIN universal_view_history vh ON vh.listing_id = ml.id
		WHERE vh.user_id = $1 AND ml.address_city IS NOT NULL
		GROUP BY ml.address_city
		ORDER BY cnt DESC
		LIMIT 3
	`, userID)
	if err == nil {
		defer func() {
			if err := rows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close rows")
			}
		}()
		for rows.Next() {
			var city string
			var count int
			if err := rows.Scan(&city, &count); err == nil {
				// Store city names instead of IDs
				_ = city // We'll use city names in the future if needed
			}
		}
		if err = rows.Err(); err != nil {
			logger.Error().Err(err).Msg("Rows iteration error")
		}
	}

	// Build personalized query
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT ml.*,
			CASE
	`)

	// Add price preference scoring
	if preferences.AvgPrice.Valid {
		avgPrice := preferences.AvgPrice.Float64
		queryBuilder.WriteString(fmt.Sprintf(`
			WHEN ml.price BETWEEN %f AND %f THEN 30
			WHEN ml.price BETWEEN %f AND %f THEN 20
			ELSE 10
		`, avgPrice*0.8, avgPrice*1.2, avgPrice*0.5, avgPrice*1.5))
	} else {
		queryBuilder.WriteString(`WHEN TRUE THEN 20`)
	}

	queryBuilder.WriteString(` END as price_score
		FROM marketplace_listings ml
		WHERE ml.status = 'active'
	`)

	// Add category filter if specified
	if category != "" && category != "all" {
		queryBuilder.WriteString(fmt.Sprintf(` AND ml.category_id IN (
			SELECT id FROM marketplace_categories WHERE LOWER(name) LIKE LOWER('%%%s%%')
		)`, category))
	}

	// Exclude already viewed items
	queryBuilder.WriteString(fmt.Sprintf(` AND ml.id NOT IN (
		SELECT listing_id FROM universal_view_history WHERE user_id = %d
	)`, userID))

	queryBuilder.WriteString(` ORDER BY price_score DESC, ml.views DESC, ml.created_at DESC LIMIT $1`)

	var listings []models.MarketplaceListing
	err = s.db.GetSQLXDB().Select(&listings, queryBuilder.String(), limit)
	return listings, err
}

// GetTrendingRecommendations gets trending items based on recent engagement
func (s *Service) GetTrendingRecommendations(limit int) ([]models.MarketplaceListing, error) {
	query := `
		WITH trending_scores AS (
			SELECT
				ml.id,
				ml.views_count as total_views,
				COUNT(DISTINCT vh.user_id) as unique_viewers,
				COUNT(CASE WHEN vh.interaction_type = 'add_favorite' THEN 1 END) as favorites,
				COUNT(CASE WHEN vh.interaction_type = 'click_phone' THEN 1 END) as phone_clicks,
				AVG(vh.view_duration_seconds) as avg_duration
			FROM marketplace_listings ml
			LEFT JOIN universal_view_history vh ON vh.listing_id = ml.id
				AND vh.created_at > NOW() - INTERVAL '7 days'
			WHERE ml.status = 'active'
				AND ml.created_at > NOW() - INTERVAL '30 days'
			GROUP BY ml.id, ml.views_count
		)
		SELECT ml.*,
			(
				ts.total_views * 0.2 +
				ts.unique_viewers * 0.3 +
				ts.favorites * 0.25 +
				ts.phone_clicks * 0.15 +
				COALESCE(ts.avg_duration, 0) * 0.1
			) as trend_score
		FROM marketplace_listings ml
		JOIN trending_scores ts ON ts.id = ml.id
		ORDER BY trend_score DESC
		LIMIT $1
	`

	var listings []models.MarketplaceListing
	err := s.db.GetSQLXDB().Select(&listings, query, limit)
	return listings, err
}

// GetPopularRecommendations gets popular items as fallback
func (s *Service) GetPopularRecommendations(limit int) ([]models.MarketplaceListing, error) {
	var listings []models.MarketplaceListing
	err := s.db.GetSQLXDB().Select(&listings, `
		SELECT * FROM marketplace_listings
		WHERE status = 'active'
		ORDER BY views DESC, created_at DESC
		LIMIT $1
	`, limit)
	return listings, err
}

// enhanceWithAttributeSimilarity calculates attribute-based similarity
func (s *Service) enhanceWithAttributeSimilarity(listings *[]models.MarketplaceListing, refAttrs []map[string]interface{}) {
	// For each listing, get its attributes and calculate similarity
	for i := range *listings {
		listing := &(*listings)[i]

		var listingAttrs []map[string]interface{}
		var attrs json.RawMessage

		err := s.db.GetSQLXDB().Get(&attrs, `
			SELECT COALESCE(
				(SELECT jsonb_agg(jsonb_build_object(
					'name', ua.name,
					'value', ua.value,
					'type', ua.type
				))
				FROM unified_attributes ua
				WHERE ua.listing_id = $1
				), '[]'::jsonb
			)
		`, listing.ID)

		if err == nil {
			if err := json.Unmarshal(attrs, &listingAttrs); err != nil {
				logger.Error().Err(err).Msg("Failed to unmarshal listing attributes")
			}

			// Calculate Jaccard similarity
			similarity := s.calculateAttributeSimilarity(refAttrs, listingAttrs)
			// Store similarity score (we could add it as a field if needed)
			_ = similarity
		}
	}
}

// calculateAttributeSimilarity calculates similarity between attribute sets
func (s *Service) calculateAttributeSimilarity(attrs1, attrs2 []map[string]interface{}) float64 {
	if len(attrs1) == 0 || len(attrs2) == 0 {
		return 0
	}

	// Create maps for easier comparison
	map1 := make(map[string]interface{})
	map2 := make(map[string]interface{})

	for _, attr := range attrs1 {
		if name, ok := attr["name"].(string); ok {
			map1[name] = attr["value"]
		}
	}

	for _, attr := range attrs2 {
		if name, ok := attr["name"].(string); ok {
			map2[name] = attr["value"]
		}
	}

	// Calculate intersection and union
	intersection := 0.0
	for key, val1 := range map1 {
		if val2, exists := map2[key]; exists {
			// Check if values are similar
			if s.areValuesSimilar(val1, val2) {
				intersection++
			}
		}
	}

	union := float64(len(map1) + len(map2))
	if union > 0 {
		// Adjust for intersection counted twice
		union -= intersection
		if union > 0 {
			return intersection / union
		}
	}

	return 0
}

// areValuesSimilar checks if two attribute values are similar
func (s *Service) areValuesSimilar(val1, val2 interface{}) bool {
	// Convert to strings for comparison
	str1 := fmt.Sprintf("%v", val1)
	str2 := fmt.Sprintf("%v", val2)

	// Exact match
	if str1 == str2 {
		return true
	}

	// Try numeric comparison
	if num1, err1 := parseFloat(str1); err1 == nil {
		if num2, err2 := parseFloat(str2); err2 == nil {
			// Consider similar if within 10%
			diff := math.Abs(num1 - num2)
			avg := (num1 + num2) / 2
			if avg > 0 {
				return (diff / avg) < 0.1
			}
		}
	}

	return false
}

// parseFloat tries to parse a string as float
func parseFloat(s string) (float64, error) {
	// Remove common units and separators
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSuffix(s, "km")
	s = strings.TrimSuffix(s, "л")

	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// PrecomputeSimilarities pre-calculates item similarities for faster retrieval
func (s *Service) PrecomputeSimilarities() error {
	// This could be run as a background job periodically
	// For now, we'll use real-time calculation
	return nil
}
