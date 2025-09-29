package handler

import (
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/pkg/utils"
	"backend/internal/logger"
	"backend/internal/proj/marketplace/service"

	"github.com/gofiber/fiber/v2"
)

// TestFuzzySearch тестирует нечеткий поиск и расширение синонимами
// @Summary Test fuzzy search functionality
// @Description Tests fuzzy search with synonyms expansion
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param lang query string false "Language" default(ru)
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=handler.FuzzySearchTestResponse} "Test results"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.queryRequired"
// @Router /api/v1/marketplace/test-fuzzy-search [get]
func (h *SearchHandler) TestFuzzySearch(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.queryRequired")
	}

	lang := c.Query("lang", "ru")

	// Расширяем запрос синонимами
	expandedQuery, err := h.marketplaceService.ExpandQueryWithSynonyms(c.Context(), query, lang)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to expand query with synonyms")
		expandedQuery = query
	}

	// Ищем похожие категории
	categoryResults, err := h.marketplaceService.SearchCategoriesFuzzy(c.Context(), query, lang, 0.3)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to search categories")
	}

	response := FuzzySearchTestResponse{
		OriginalQuery: query,
		ExpandedQuery: expandedQuery,
		Language:      lang,
		Categories:    categoryResults,
	}

	return utils.SuccessResponse(c, response)
}

// FuzzySearchTestResponse ответ для теста нечеткого поиска
type FuzzySearchTestResponse struct {
	OriginalQuery string                         `json:"original_query"`
	ExpandedQuery string                         `json:"expanded_query"`
	Language      string                         `json:"language"`
	Categories    []service.CategorySearchResult `json:"categories"`
}

// SearchWithFuzzyParams выполняет поиск с указанными параметрами нечеткости
// @Summary Search with custom fuzzy parameters
// @Description Performs search with customizable fuzzy matching parameters
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param fuzziness query string false "Fuzziness level (AUTO, 0, 1, 2)" default(AUTO)
// @Param minimum_should_match query string false "Minimum should match (e.g., 30%, 2)" default(30%)
// @Param use_synonyms query bool false "Use synonym expansion" default(true)
// @Param category_id query string false "Category ID"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} handler.SearchResponse "Search results"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.queryRequired"
// @Router /api/v1/marketplace/fuzzy-search [get]
func (h *SearchHandler) SearchWithFuzzyParams(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.queryRequired")
	}

	// Получаем параметры нечеткого поиска
	fuzziness := c.Query("fuzziness", "AUTO")
	minimumShouldMatch := c.Query("minimum_should_match", "30%")
	useSynonyms := c.Query("use_synonyms", "true") == "true"

	// Базовые параметры поиска
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	categoryID := c.Query("category_id")

	// Создаем параметры поиска
	params := search.ServiceParams{
		Query:              query,
		Page:               page,
		Size:               limit,
		CategoryID:         categoryID,
		Language:           c.Query("lang", "ru"),
		Fuzziness:          fuzziness,
		MinimumShouldMatch: minimumShouldMatch,
		UseSynonyms:        useSynonyms,
	}

	// Выполняем поиск
	results, err := h.marketplaceService.SearchListingsAdvanced(c.Context(), &params)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to perform fuzzy search")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.searchError")
	}

	// Преобразуем результаты
	items := results.Items
	if items == nil {
		items = []*models.MarketplaceListing{}
	}

	listings := make([]models.MarketplaceListing, 0, len(items))
	for _, item := range items {
		if item != nil {
			listingCopy := *item
			if item.Images != nil {
				listingCopy.Images = make([]models.MarketplaceImage, len(item.Images))
				copy(listingCopy.Images, item.Images)
			}
			listings = append(listings, listingCopy)
		}
	}

	response := SearchResponse{
		Data: listings,
		Meta: SearchMetadata{
			Total:              results.Total,
			Page:               page,
			Size:               limit,
			TotalPages:         (results.Total + limit - 1) / limit,
			HasMore:            page < (results.Total+limit-1)/limit,
			Facets:             results.Facets,
			Suggestions:        results.Suggestions,
			SpellingSuggestion: results.SpellingSuggestion,
			TookMs:             results.Took,
		},
	}

	return c.JSON(response)
}
