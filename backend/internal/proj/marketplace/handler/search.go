// TEMPORARY: Will be moved to microservice
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"backend/pkg/utils"
)

var osClient *opensearch.Client

// InitOpenSearchClient инициализирует клиент OpenSearch
func InitOpenSearchClient() error {
	cfg := opensearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	osClient = client
	return nil
}

// SearchListings godoc
// @Summary Поиск объявлений
// @Description Публичный поиск объявлений через OpenSearch
// @Tags marketplace
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param category_id query int false "ID категории"
// @Param limit query int false "Лимит результатов" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/search [get]
func (h *Handler) SearchListings(c *fiber.Ctx) error {
	// Ensure OpenSearch client is initialized
	if osClient == nil {
		if err := InitOpenSearchClient(); err != nil {
			h.logger.Error().Err(err).Msg("Failed to initialize OpenSearch client")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.search_unavailable")
		}
	}

	// Parse query parameters
	query := c.Query("query", "")
	categoryID := c.QueryInt("category_id", 0)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Validate parameters
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Build OpenSearch query
	searchQuery := buildSearchQuery(query, categoryID, limit, offset)

	// Convert query to JSON
	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal search query")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.search_query_error")
	}

	// Execute search
	req := opensearchapi.SearchRequest{
		Index: []string{"marketplace_listings"},
		Body:  bytes.NewReader(queryJSON),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := req.Do(ctx, osClient)
	if err != nil {
		h.logger.Error().Err(err).Msg("OpenSearch request failed")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.search_failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		h.logger.Error().
			Str("status", res.Status()).
			Msg("OpenSearch returned error")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.search_error")
	}

	// Parse response
	var searchResp OpenSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse OpenSearch response")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.search_parse_error")
	}

	// Extract listings from hits
	listings := make([]map[string]interface{}, 0, len(searchResp.Hits.Hits))
	for _, hit := range searchResp.Hits.Hits {
		listings = append(listings, hit.Source)
	}

	// Build response
	response := map[string]interface{}{
		"listings": listings,
		"total":    searchResp.Hits.Total.Value,
		"limit":    limit,
		"offset":   offset,
	}

	return utils.SuccessResponse(c, response)
}

// buildSearchQuery создает запрос для OpenSearch
func buildSearchQuery(query string, categoryID, limit, offset int) *map[string]interface{} {
	// Build bool query
	must := []interface{}{}

	// Only active listings
	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{
			"status": "active",
		},
	})

	// Add text search if provided
	if query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^3", "description"},
			},
		})
	}

	// Add category filter if provided
	if categoryID > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": categoryID,
			},
		})
	}

	// Build final query
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"size": limit,
		"from": offset,
		"sort": []interface{}{
			map[string]interface{}{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return &searchQuery
}

// OpenSearch response structures
type OpenSearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
