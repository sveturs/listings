package handler

import (
	"strconv"
	"time"

	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// SavedSearchesHandler обрабатывает запросы, связанные с сохраненными поисками
type SavedSearchesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewSavedSearchesHandler создает новый обработчик сохраненных поисков
func NewSavedSearchesHandler(services globalService.ServicesInterface) *SavedSearchesHandler {
	return &SavedSearchesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// CreateSavedSearchRequest структура запроса для создания сохраненного поиска
type CreateSavedSearchRequest struct {
	Name            string                 `json:"name" validate:"required,min=1,max=100"`
	Filters         map[string]interface{} `json:"filters" validate:"required"`
	SearchType      string                 `json:"search_type"`
	NotifyEnabled   bool                   `json:"notify_enabled"`
	NotifyFrequency string                 `json:"notify_frequency"`
}

// UpdateSavedSearchRequest структура запроса для обновления сохраненного поиска
type UpdateSavedSearchRequest struct {
	Name            string                 `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Filters         map[string]interface{} `json:"filters,omitempty"`
	NotifyEnabled   *bool                  `json:"notify_enabled,omitempty"`
	NotifyFrequency string                 `json:"notify_frequency,omitempty"`
}

// SavedSearchResponse структура ответа сохраненного поиска
type SavedSearchResponse struct {
	ID              int                    `json:"id"`
	UserID          int                    `json:"user_id"`
	Name            string                 `json:"name"`
	Filters         map[string]interface{} `json:"filters"`
	SearchType      string                 `json:"search_type"`
	NotifyEnabled   bool                   `json:"notify_enabled"`
	NotifyFrequency string                 `json:"notify_frequency"`
	ResultsCount    int                    `json:"results_count"`
	LastNotifiedAt  *time.Time             `json:"last_notified_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CreateSavedSearch создает новый сохраненный поиск
// @Summary Create saved search
// @Description Creates a new saved search for the user
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Param body body CreateSavedSearchRequest true "Saved search data"
// @Success 200 {object} utils.SuccessResponseSwag{data=SavedSearchResponse} "Created saved search"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createSavedSearchError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches [post]
func (h *SavedSearchesHandler) CreateSavedSearch(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	var req CreateSavedSearchRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Установка значений по умолчанию
	if req.SearchType == "" {
		req.SearchType = "cars"
	}
	if req.NotifyFrequency == "" {
		req.NotifyFrequency = "daily"
	}

	// Создаем сохраненный поиск
	savedSearch, err := h.marketplaceService.CreateSavedSearch(c.Context(), userID, req.Name, req.Filters, req.SearchType, req.NotifyEnabled, req.NotifyFrequency)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to create saved search")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createSavedSearchError")
	}

	return utils.SuccessResponse(c, savedSearch)
}

// GetSavedSearches получает список сохраненных поисков пользователя
// @Summary Get user's saved searches
// @Description Returns all saved searches for the user
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]SavedSearchResponse} "List of saved searches"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getSavedSearchesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches [get]
func (h *SavedSearchesHandler) GetSavedSearches(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем тип поиска из query параметров (опционально)
	searchType := c.Query("search_type", "")

	savedSearches, err := h.marketplaceService.GetUserSavedSearches(c.Context(), userID, searchType)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get saved searches")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getSavedSearchesError")
	}

	if savedSearches == nil {
		savedSearches = []interface{}{}
	}

	return utils.SuccessResponse(c, savedSearches)
}

// GetSavedSearch получает конкретный сохраненный поиск
// @Summary Get saved search by ID
// @Description Returns a specific saved search
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Param id path int true "Saved search ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=SavedSearchResponse} "Saved search details"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.savedSearchNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getSavedSearchError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches/{id} [get]
func (h *SavedSearchesHandler) GetSavedSearch(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	searchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	savedSearch, err := h.marketplaceService.GetSavedSearchByID(c.Context(), userID, searchID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Int("searchId", searchID).Msg("Failed to get saved search")
		if err.Error() == "saved search not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.savedSearchNotFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getSavedSearchError")
	}

	return utils.SuccessResponse(c, savedSearch)
}

// UpdateSavedSearch обновляет сохраненный поиск
// @Summary Update saved search
// @Description Updates a saved search
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Param id path int true "Saved search ID"
// @Param body body UpdateSavedSearchRequest true "Update data"
// @Success 200 {object} utils.SuccessResponseSwag{data=SavedSearchResponse} "Updated saved search"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.savedSearchNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateSavedSearchError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches/{id} [put]
func (h *SavedSearchesHandler) UpdateSavedSearch(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	searchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	var req UpdateSavedSearchRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Обновляем сохраненный поиск
	savedSearch, err := h.marketplaceService.UpdateSavedSearch(c.Context(), userID, searchID, req.Name, req.Filters, req.NotifyEnabled, req.NotifyFrequency)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Int("searchId", searchID).Msg("Failed to update saved search")
		if err.Error() == "saved search not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.savedSearchNotFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateSavedSearchError")
	}

	return utils.SuccessResponse(c, savedSearch)
}

// DeleteSavedSearch удаляет сохраненный поиск
// @Summary Delete saved search
// @Description Deletes a saved search
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Param id path int true "Saved search ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.savedSearchNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteSavedSearchError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches/{id} [delete]
func (h *SavedSearchesHandler) DeleteSavedSearch(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	searchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	err = h.marketplaceService.DeleteSavedSearch(c.Context(), userID, searchID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Int("searchId", searchID).Msg("Failed to delete saved search")
		if err.Error() == "saved search not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.savedSearchNotFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteSavedSearchError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.savedSearchDeleted",
	})
}

// ExecuteSavedSearch выполняет сохраненный поиск и возвращает результаты
// @Summary Execute saved search
// @Description Executes a saved search and returns the results
// @Tags marketplace-saved-searches
// @Accept json
// @Produce json
// @Param id path int true "Saved search ID"
// @Success 200 {object} utils.SuccessResponseSwag "Search results"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.savedSearchNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.executeSavedSearchError"
// @Security BearerAuth
// @Router /api/v1/marketplace/saved-searches/{id}/execute [get]
func (h *SavedSearchesHandler) ExecuteSavedSearch(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	searchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем сохраненный поиск
	savedSearch, err := h.marketplaceService.GetSavedSearchByID(c.Context(), userID, searchID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Int("searchId", searchID).Msg("Failed to get saved search")
		if err.Error() == "saved search not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.savedSearchNotFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.executeSavedSearchError")
	}

	// Выполняем поиск с сохраненными фильтрами
	results, err := h.marketplaceService.ExecuteSavedSearch(c.Context(), savedSearch)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Int("searchId", searchID).Msg("Failed to execute saved search")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.executeSavedSearchError")
	}

	return utils.SuccessResponse(c, results)
}
