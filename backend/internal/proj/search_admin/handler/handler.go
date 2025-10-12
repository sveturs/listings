package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"backend/internal/domain"
	"backend/internal/proj/search_admin/service"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
	logger  *logger.Logger
}

func NewHandler(service *service.Service, logger *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// GetWeights godoc
// @Summary Get all search field weights
// @Description Get weights for all searchable fields
// @Tags search-config
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]domain.SearchWeight} "List of search weights"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/weights [get]
func (h *Handler) GetWeights(c *fiber.Ctx) error {
	weights, err := h.service.GetWeights(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, weights)
}

// GetWeightByField godoc
// @Summary Get weight for specific field
// @Description Get weight configuration for a specific search field
// @Tags search-config
// @Accept json
// @Produce json
// @Param field path string true "Field name"
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.SearchWeight} "Search weight details"
// @Failure 404 {object} utils.ErrorResponseSwag "Weight not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/weights/{field} [get]
func (h *Handler) GetWeightByField(c *fiber.Ctx) error {
	fieldName := c.Params("field")

	weight, err := h.service.GetWeightByField(c.Context(), fieldName)
	if err != nil {
		// h.log.Error("Failed to get weight by field", "field", fieldName, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	if weight == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.weightNotFound")
	}

	return utils.SuccessResponse(c, weight)
}

// CreateWeight godoc
// @Summary Create new search weight
// @Description Create weight configuration for a search field
// @Tags search-config
// @Accept json
// @Produce json
// @Param weight body domain.SearchWeight true "Weight configuration"
// @Success 201 {object} utils.SuccessResponseSwag{data=domain.SearchWeight} "Created weight"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/weights [post]
func (h *Handler) CreateWeight(c *fiber.Ctx) error {
	var weight domain.SearchWeight
	if err := c.BodyParser(&weight); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err := h.service.CreateWeight(c.Context(), &weight)
	if err != nil {
		// h.log.Error("Failed to create weight", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, weight)
}

// UpdateWeight godoc
// @Summary Update search weight
// @Description Update weight configuration for a search field
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Weight ID"
// @Param weight body domain.SearchWeight true "Updated weight configuration"
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.SearchWeight} "Updated weight"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Weight not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/weights/{id} [put]
func (h *Handler) UpdateWeight(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	var weight domain.SearchWeight
	if err := c.BodyParser(&weight); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err = h.service.UpdateWeight(c.Context(), id, &weight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.weightNotFound")
		}
		// h.log.Error("Failed to update weight", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	weight.ID = id
	return utils.SuccessResponse(c, weight)
}

// DeleteWeight godoc
// @Summary Delete search weight
// @Description Delete weight configuration for a search field
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Weight ID"
// @Success 204 {object} nil "Weight deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Weight not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/weights/{id} [delete]
func (h *Handler) DeleteWeight(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	err = h.service.DeleteWeight(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.weightNotFound")
		}
		// h.log.Error("Failed to delete weight", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSynonyms godoc
// @Summary Get search synonyms
// @Description Get all synonyms for search terms with pagination
// @Tags search-config
// @Accept json
// @Produce json
// @Param language query string false "Language code (default: ru)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{} "Paginated list of synonyms"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/synonyms [get]
func (h *Handler) GetSynonyms(c *fiber.Ctx) error {
	language := c.Query("language", "ru")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	search := c.Query("search", "")

	// Защита от некорректных значений
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Получаем все синонимы
	allSynonyms, err := h.service.GetSynonyms(c.Context(), language)
	if err != nil {
		// h.log.Error("Failed to get synonyms", "language", language, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	// Фильтруем по поисковому запросу если есть
	filteredSynonyms := allSynonyms
	if search != "" {
		filteredSynonyms = []domain.SearchSynonym{}
		for _, syn := range allSynonyms {
			// Поиск в term или synonyms
			if contains(syn.Term, search) {
				filteredSynonyms = append(filteredSynonyms, syn)
				continue
			}
			for _, s := range syn.Synonyms {
				if contains(s, search) {
					filteredSynonyms = append(filteredSynonyms, syn)
					break
				}
			}
		}
	}

	// Подсчитываем total
	total := len(filteredSynonyms)
	totalPages := (total + limit - 1) / limit

	// Применяем пагинацию
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if startIndex > total {
		startIndex = total
	}
	if endIndex > total {
		endIndex = total
	}

	paginatedSynonyms := []domain.SearchSynonym{}
	if startIndex < endIndex {
		paginatedSynonyms = filteredSynonyms[startIndex:endIndex]
	}

	// Возвращаем с пагинацией
	response := fiber.Map{
		"success": true,
		"data": fiber.Map{
			"data":        paginatedSynonyms,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	}

	// Логируем ответ для отладки
	log.Printf("GetSynonyms response: success=%v, total=%d, page=%d, synonyms count=%d", true, total, page, len(paginatedSynonyms))

	return c.JSON(response)
}

// Helper function для case-insensitive поиска
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) > 0 && len(s) > 0 &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CreateSynonym godoc
// @Summary Create search synonym
// @Description Create synonym mapping for search terms
// @Tags search-config
// @Accept json
// @Produce json
// @Param synonym body domain.SearchSynonym true "Synonym configuration"
// @Success 201 {object} utils.SuccessResponseSwag{data=domain.SearchSynonym} "Created synonym"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/synonyms [post]
func (h *Handler) CreateSynonym(c *fiber.Ctx) error {
	c.Locals("logger_prefix", "CreateSynonym")
	log.Printf("CreateSynonym: Starting handler")

	var synonym domain.SearchSynonym
	if err := c.BodyParser(&synonym); err != nil {
		log.Printf("CreateSynonym: Failed to parse body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	log.Printf("CreateSynonym: Parsed synonym: %+v", synonym)

	err := h.service.CreateSynonym(c.Context(), &synonym)
	if err != nil {
		log.Printf("CreateSynonym: Service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Service error"})
	}

	log.Printf("CreateSynonym: Success")
	return c.Status(201).JSON(fiber.Map{"success": true, "data": synonym})
}

// UpdateSynonym godoc
// @Summary Update search synonym
// @Description Update synonym mapping for search terms
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Synonym ID"
// @Param synonym body domain.SearchSynonym true "Updated synonym configuration"
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.SearchSynonym} "Updated synonym"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Synonym not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/synonyms/{id} [put]
func (h *Handler) UpdateSynonym(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	var synonym domain.SearchSynonym
	if err := c.BodyParser(&synonym); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err = h.service.UpdateSynonym(c.Context(), id, &synonym)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.synonymNotFound")
		}
		// h.log.Error("Failed to update synonym", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	synonym.ID = id
	return utils.SuccessResponse(c, synonym)
}

// DeleteSynonym godoc
// @Summary Delete search synonym
// @Description Delete synonym mapping
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Synonym ID"
// @Success 204 {object} nil "Synonym deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Synonym not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/synonyms/{id} [delete]
func (h *Handler) DeleteSynonym(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	err = h.service.DeleteSynonym(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.synonymNotFound")
		}
		// h.log.Error("Failed to delete synonym", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetTransliterationRules godoc
// @Summary Get transliteration rules
// @Description Get all transliteration rules for search
// @Tags search-config
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]domain.TransliterationRule} "List of transliteration rules"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/transliteration [get]
func (h *Handler) GetTransliterationRules(c *fiber.Ctx) error {
	rules, err := h.service.GetTransliterationRules(c.Context())
	if err != nil {
		// h.log.Error("Failed to get transliteration rules", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, rules)
}

// CreateTransliterationRule godoc
// @Summary Create transliteration rule
// @Description Create new transliteration rule for search
// @Tags search-config
// @Accept json
// @Produce json
// @Param rule body domain.TransliterationRule true "Transliteration rule"
// @Success 201 {object} utils.SuccessResponseSwag{data=domain.TransliterationRule} "Created rule"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/transliteration [post]
func (h *Handler) CreateTransliterationRule(c *fiber.Ctx) error {
	var rule domain.TransliterationRule
	if err := c.BodyParser(&rule); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err := h.service.CreateTransliterationRule(c.Context(), &rule)
	if err != nil {
		// h.log.Error("Failed to create transliteration rule", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, rule)
}

// UpdateTransliterationRule godoc
// @Summary Update transliteration rule
// @Description Update existing transliteration rule
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param rule body domain.TransliterationRule true "Updated rule"
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.TransliterationRule} "Updated rule"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Rule not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/transliteration/{id} [put]
func (h *Handler) UpdateTransliterationRule(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	var rule domain.TransliterationRule
	if err := c.BodyParser(&rule); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err = h.service.UpdateTransliterationRule(c.Context(), id, &rule)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.ruleNotFound")
		}
		// h.log.Error("Failed to update transliteration rule", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	rule.ID = id
	return utils.SuccessResponse(c, rule)
}

// DeleteTransliterationRule godoc
// @Summary Delete transliteration rule
// @Description Delete transliteration rule
// @Tags search-config
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Success 204 {object} nil "Rule deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 404 {object} utils.ErrorResponseSwag "Rule not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config/transliteration/{id} [delete]
func (h *Handler) DeleteTransliterationRule(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	err = h.service.DeleteTransliterationRule(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.ruleNotFound")
		}
		// h.log.Error("Failed to delete transliteration rule", "id", id, "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSearchStatistics godoc
// @Summary Get search statistics
// @Description Get recent search statistics
// @Tags search-config
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]domain.SearchStatistics} "Search statistics"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/statistics [get]
func (h *Handler) GetSearchStatistics(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 100)

	stats, err := h.service.GetSearchStatistics(c.Context(), limit)
	if err != nil {
		// h.log.Error("Failed to get search statistics", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, stats)
}

// GetPopularSearches godoc
// @Summary Get popular searches
// @Description Get most popular search queries from last 7 days
// @Tags search-config
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 10)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]map[string]interface{}} "Popular searches"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/statistics/popular [get]
func (h *Handler) GetPopularSearches(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)

	searches, err := h.service.GetPopularSearches(c.Context(), limit)
	if err != nil {
		// h.log.Error("Failed to get popular searches", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, searches)
}

// GetConfig godoc
// @Summary Get search configuration
// @Description Get general search configuration settings
// @Tags search-config
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.SearchConfig} "Search configuration"
// @Failure 404 {object} utils.ErrorResponseSwag "Configuration not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config [get]
func (h *Handler) GetConfig(c *fiber.Ctx) error {
	config, err := h.service.GetConfig(c.Context())
	if err != nil {
		// h.log.Error("Failed to get search config", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	if config == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.configNotFound")
	}

	return utils.SuccessResponse(c, config)
}

// UpdateConfig godoc
// @Summary Update search configuration
// @Description Update general search configuration settings
// @Tags search-config
// @Accept json
// @Produce json
// @Param config body domain.SearchConfig true "Search configuration"
// @Success 200 {object} utils.SuccessResponseSwag{data=domain.SearchConfig} "Updated configuration"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/search/config [put]
func (h *Handler) UpdateConfig(c *fiber.Ctx) error {
	var config domain.SearchConfig
	if err := c.BodyParser(&config); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	err := h.service.UpdateConfig(c.Context(), &config)
	if err != nil {
		// h.log.Error("Failed to update search config", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internalServerError")
	}

	return utils.SuccessResponse(c, config)
}

// GetSearchAnalytics - REMOVED (deprecated, use /api/v1/analytics/metrics/search instead)
