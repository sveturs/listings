package translation_admin

import (
	"strconv"

	"backend/internal/domain/models"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	translationNotFoundError = "translation not found"
)

// Handler handles translation admin API requests
type Handler struct {
	service *Service
	logger  zerolog.Logger
}

// NewHandler creates a new translation admin handler
func NewHandler(service *Service, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers translation admin routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Frontend translations
	frontend := router.Group("/frontend")
	frontend.Get("/modules", h.GetFrontendModules)
	frontend.Get("/module/:name", h.GetModuleTranslations)
	frontend.Put("/module/:name", h.UpdateModuleTranslations)
	frontend.Post("/validate", h.ValidateTranslations)
	frontend.Post("/sync", h.SyncFrontendToDBOld)

	// Database translations
	database := router.Group("/database")
	database.Get("/", h.GetDatabaseTranslations)
	database.Get("/:id", h.GetTranslation)
	database.Put("/:id", h.UpdateTranslation)
	database.Delete("/:id", h.DeleteTranslation)
	database.Post("/batch", h.BatchOperations)

	// AI translations - основные обработчики в AITranslationHandler (module.go)
	ai := router.Group("/ai")
	// ai.Post("/translate", h.TranslateText) // используется AITranslationHandler.TranslateText
	// ai.Post("/batch", h.BatchTranslate) // используется AITranslationHandler.TranslateBatch
	ai.Post("/apply", h.ApplyAITranslations)
	ai.Get("/providers-old", h.GetProviders) // старый метод
	ai.Put("/providers/:id", h.UpdateProvider)

	// Translation providers - новый эндпоинт для системы провайдеров
	router.Get("/providers", h.GetTranslationProviders)
	router.Put("/providers/:id", h.UpdateTranslationProvider)

	// Export/Import
	router.Get("/export", h.ExportToJSON)
	router.Post("/import", h.ImportFromJSON)
	router.Post("/export/advanced", h.ExportAdvanced)
	router.Post("/import/advanced", h.ImportAdvanced)

	// Bulk operations
	bulk := router.Group("/bulk")
	bulk.Post("/translate", h.BulkTranslate)

	// Synchronization
	sync := router.Group("/sync")
	sync.Post("/frontend-to-db", h.SyncFrontendToDB)
	sync.Post("/db-to-frontend", h.SyncDBToFrontend)
	sync.Post("/db-to-opensearch", h.SyncDBToOpenSearch)
	sync.Get("/status", h.GetSyncStatus)
	sync.Get("/conflicts", h.GetConflicts)
	sync.Post("/conflicts/:id/resolve", h.ResolveConflict)
	sync.Post("/conflicts/resolve", h.ResolveConflictsBatch)

	// Versioning
	versions := router.Group("/versions")
	versions.Get("/:entity/:id", h.GetVersionHistory)
	versions.Get("/translation/:id", h.GetTranslationVersions)
	versions.Post("/rollback", h.RollbackVersion)
	versions.Get("/diff", h.GetVersionDiff)

	// Statistics
	stats := router.Group("/stats")
	stats.Get("/overview", h.GetStatisticsOverview)
	stats.Get("/coverage", h.GetCoverage)
	stats.Get("/quality", h.GetQuality)
	stats.Get("/usage", h.GetUsage)

	// Audit
	audit := router.Group("/audit")
	audit.Get("/logs", h.GetAuditLogs)
	audit.Get("/statistics", h.GetAuditStatistics)
}

// GetFrontendModules godoc
// @Summary Get all frontend translation modules
// @Description Returns list of all frontend modules with statistics
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.FrontendModule}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/frontend/modules [get]
func (h *Handler) GetFrontendModules(c *fiber.Ctx) error {
	ctx := c.Context()

	modules, err := h.service.GetFrontendModules(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get frontend modules")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getModulesError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.modulesRetrieved", modules)
}

// GetModuleTranslations godoc
// @Summary Get translations for a specific module
// @Description Returns all translations for the specified module
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param name path string true "Module name"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.FrontendTranslation}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/frontend/module/{name} [get]
func (h *Handler) GetModuleTranslations(c *fiber.Ctx) error {
	ctx := c.Context()
	moduleName := c.Params("name")

	if moduleName == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidModuleName")
	}

	translations, err := h.service.GetModuleTranslations(ctx, moduleName)
	if err != nil {
		h.logger.Error().Err(err).Str("module", moduleName).Msg("Failed to get module translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getTranslationsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationsRetrieved", translations)
}

// UpdateModuleTranslations godoc
// @Summary Update translations for a module
// @Description Updates multiple translations within a module
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param name path string true "Module name"
// @Param updates body []backend_internal_domain_models.FrontendTranslation true "Translation updates"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/frontend/module/{name} [put]
func (h *Handler) UpdateModuleTranslations(c *fiber.Ctx) error {
	ctx := c.Context()
	moduleName := c.Params("name")

	h.logger.Info().Str("module", moduleName).Msg("UpdateModuleTranslations called")

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		h.logger.Error().Msg("User ID not found in context")
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	h.logger.Info().Int("userID", userID).Msg("User ID found")

	var updates []models.FrontendTranslation
	if err := c.BodyParser(&updates); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	h.logger.Info().Int("updatesCount", len(updates)).Msg("Parsed updates")

	if err := h.service.UpdateModuleTranslations(ctx, moduleName, updates, userID); err != nil {
		h.logger.Error().Err(err).Str("module", moduleName).Msg("Failed to update module translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.updateError")
	}

	h.logger.Info().Msg("Successfully updated translations")
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationsUpdated", nil)
}

// ValidateTranslations godoc
// @Summary Validate translations for consistency
// @Description Validates translations across modules and languages
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.ValidateTranslationsRequest true "Validation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.ValidationResult}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/frontend/validate [post]
func (h *Handler) ValidateTranslations(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.ValidateTranslationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	results, err := h.service.ValidateTranslations(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to validate translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.validationError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.validationComplete", results)
}

// SyncFrontendToDBOld godoc - старый метод
// @Summary Sync frontend translations to database OLD
// @Description Synchronizes all frontend translations to the database (OLD METHOD)
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/sync/frontend-to-db-old [post]
func (h *Handler) SyncFrontendToDBOld(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	if err := h.service.SyncFrontendToDBOld(ctx, userID); err != nil {
		h.logger.Error().Err(err).Msg("Failed to sync frontend to DB")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.syncError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.syncComplete", nil)
}

// GetStatisticsOverview godoc
// @Summary Get translation statistics overview
// @Description Returns comprehensive statistics about translations
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.TranslationStatistics}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/stats/overview [get]
func (h *Handler) GetStatisticsOverview(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.service.GetStatistics(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get statistics")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.statsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.statsRetrieved", stats)
}

// Placeholder handlers for remaining endpoints

// GetDatabaseTranslations godoc
// @Summary Get database translations with filters
// @Description Returns translations from database with optional filters
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type filter"
// @Param entity_id query int false "Entity ID filter"
// @Param language query string false "Language filter"
// @Param field_name query string false "Field name filter"
// @Param is_verified query bool false "Verified status filter"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.Translation}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/database [get]
func (h *Handler) GetDatabaseTranslations(c *fiber.Ctx) error {
	ctx := c.Context()

	// Build filters from query parameters
	filters := make(map[string]interface{})

	if entityType := c.Query("entity_type"); entityType != "" {
		filters["entity_type"] = entityType
	}
	if entityID := c.QueryInt("entity_id", 0); entityID > 0 {
		filters["entity_id"] = entityID
	}
	if language := c.Query("language"); language != "" {
		filters["language"] = language
	}
	if fieldName := c.Query("field_name"); fieldName != "" {
		filters["field_name"] = fieldName
	}
	if isVerified := c.Query("is_verified"); isVerified != "" {
		filters["is_verified"] = isVerified == "true"
	}
	if limit := c.QueryInt("limit", 100); limit > 0 {
		filters["limit"] = limit
	}
	if offset := c.QueryInt("offset", 0); offset >= 0 {
		filters["offset"] = offset
	}

	translations, err := h.service.GetDatabaseTranslations(ctx, filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get database translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getDatabaseError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.databaseTranslationsRetrieved", translations)
}

// GetTranslation godoc
// @Summary Get single translation by ID
// @Description Returns a single translation from database
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param id path int true "Translation ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.Translation}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/database/{id} [get]
func (h *Handler) GetTranslation(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidTranslationID")
	}

	translation, err := h.service.GetTranslationByID(ctx, id)
	if err != nil {
		if err.Error() == translationNotFoundError {
			return utils.SendError(c, fiber.StatusNotFound, "admin.translations.translationNotFound")
		}
		h.logger.Error().Err(err).Int("id", id).Msg("Failed to get translation")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getTranslationError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationRetrieved", translation)
}

// UpdateTranslation godoc
// @Summary Update translation in database
// @Description Updates an existing translation
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param id path int true "Translation ID"
// @Param translation body backend_internal_domain_models.TranslationUpdateRequest true "Translation update"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/database/{id} [put]
func (h *Handler) UpdateTranslation(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidTranslationID")
	}

	var updateReq models.TranslationUpdateRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	if err := h.service.UpdateDatabaseTranslation(ctx, id, &updateReq, userID); err != nil {
		if err.Error() == translationNotFoundError {
			return utils.SendError(c, fiber.StatusNotFound, "admin.translations.translationNotFound")
		}
		h.logger.Error().Err(err).Int("id", id).Msg("Failed to update translation")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.updateDatabaseError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationUpdated", nil)
}

// DeleteTranslation godoc
// @Summary Delete translation from database
// @Description Deletes a translation permanently
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param id path int true "Translation ID"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/database/{id} [delete]
func (h *Handler) DeleteTranslation(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidTranslationID")
	}

	if err := h.service.DeleteDatabaseTranslation(ctx, id, userID); err != nil {
		if err.Error() == translationNotFoundError {
			return utils.SendError(c, fiber.StatusNotFound, "admin.translations.translationNotFound")
		}
		h.logger.Error().Err(err).Int("id", id).Msg("Failed to delete translation")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.deleteDatabaseError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationDeleted", nil)
}

// BatchOperations godoc
// @Summary Perform batch operations on translations
// @Description Create, update or delete multiple translations at once
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param operations body backend_internal_domain_models.BatchOperationsRequest true "Batch operations"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.BatchOperationsResult}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/database/batch [post]
func (h *Handler) BatchOperations(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.BatchOperationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	result, err := h.service.PerformBatchOperations(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to perform batch operations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.batchOperationError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.batchOperationComplete", result)
}

func (h *Handler) TranslateText(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.TranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	result, err := h.service.TranslateText(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to translate text")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.translateError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.textTranslated", result)
}

func (h *Handler) BatchTranslate(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.AIBatchTranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	result, err := h.service.BatchTranslate(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to batch translate")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.batchTranslateError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.batchTranslated", result)
}

func (h *Handler) GetProviders(c *fiber.Ctx) error {
	ctx := c.Context()

	providers, err := h.service.GetAIProviders(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get AI providers")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getProvidersError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providersRetrieved", providers)
}

func (h *Handler) UpdateProvider(c *fiber.Ctx) error {
	ctx := c.Context()
	providerID := c.Params("id")

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var provider models.AIProvider
	if err := c.BodyParser(&provider); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	provider.ID = providerID

	if err := h.service.UpdateAIProvider(ctx, &provider, userID); err != nil {
		h.logger.Error().Err(err).Str("provider_id", providerID).Msg("Failed to update AI provider")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.updateProviderError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providerUpdated", nil)
}

func (h *Handler) SyncDBToOpenSearch(c *fiber.Ctx) error {
	// TODO: Implement DB to OpenSearch sync
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.notImplemented", nil)
}

// ExportToJSON exports database translations to JSON format
// @Summary Export translations to JSON
// @Description Export database translations to JSON format for backup or migration
// @Tags admin-translations
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type to export (frontend, database, all)"
// @Param language query string false "Language code to export (en, ru, sr, all)"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/export [get]
func (h *Handler) ExportToJSON(c *fiber.Ctx) error {
	ctx := c.Context()

	entityType := c.Query("entity_type", "all")
	language := c.Query("language", "all")

	// Экспортируем переводы из БД
	translations, err := h.service.ExportTranslations(ctx, entityType, language)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to export translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.exportError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.exportSuccess", translations)
}

// ImportFromJSON imports translations from JSON to database
// @Summary Import translations from JSON
// @Description Import translations from JSON format to database
// @Tags admin-translations
// @Accept json
// @Produce json
// @Param body body backend_internal_domain_models.ImportTranslationsRequest true "Import request"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/import [post]
func (h *Handler) ImportFromJSON(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.ImportTranslationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Импортируем переводы в БД
	result, err := h.service.ImportTranslations(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to import translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.importError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.importSuccess", result)
}

// SyncFrontendToDB syncs frontend JSON translations to database
// @Summary Sync frontend to database
// @Description Synchronize frontend JSON translations to database
// @Tags admin-translations
// @Accept json
// @Produce json
// @Param module query string false "Module name to sync (all for all modules)"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/sync/frontend-to-db [post]
func (h *Handler) SyncFrontendToDB(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	module := c.Query("module", "all")

	// Синхронизируем Frontend -> DB
	result, err := h.service.SyncFrontendToDB(ctx, module, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to sync frontend to DB")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.syncError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.syncSuccess", result)
}

// SyncDBToFrontend syncs database translations to frontend JSON files
// @Summary Sync database to frontend
// @Description Synchronize database translations to frontend JSON files
// @Tags admin-translations
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type to sync"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/sync/db-to-frontend [post]
func (h *Handler) SyncDBToFrontend(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	entityType := c.Query("entity_type", "all")

	// Синхронизируем DB -> Frontend
	result, err := h.service.SyncDBToFrontend(ctx, entityType, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to sync DB to frontend")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.syncError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.syncSuccess", result)
}

func (h *Handler) GetSyncStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	// Получаем статус синхронизации
	status, err := h.service.GetSyncStatus(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get sync status")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.syncStatusError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.syncStatus", status)
}

func (h *Handler) GetConflicts(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse filter parameters
	filter := &models.ConflictsFilter{}

	if module := c.Query("module"); module != "" {
		filter.Module = &module
	}

	if language := c.Query("language"); language != "" {
		filter.Language = &language
	}

	if conflictType := c.Query("type"); conflictType != "" {
		filter.Type = &conflictType
	}

	if resolved := c.Query("resolved"); resolved != "" {
		resolvedBool := resolved == "true"
		filter.Resolved = &resolvedBool
	}

	// Get conflicts from service
	conflicts, err := h.service.GetConflicts(ctx, filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get conflicts")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getConflictsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.conflictsRetrieved", conflicts.Conflicts)
}

func (h *Handler) ResolveConflict(c *fiber.Ctx) error {
	conflictIDStr := c.Params("id")
	conflictID, err := strconv.Atoi(conflictIDStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidConflictID")
	}

	// TODO: Implement resolve conflict
	_ = conflictID
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.notImplemented", nil)
}

// ResolveConflictsBatch godoc
// @Summary Resolve multiple translation conflicts
// @Description Batch resolve translation sync conflicts
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param resolutions body backend_internal_domain_models.ConflictResolutionBatch true "Conflict resolutions"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.ConflictResolutionResult}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/sync/conflicts/resolve [post]
func (h *Handler) ResolveConflictsBatch(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.ConflictResolutionBatch
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Validate request
	if len(req.Resolutions) == 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.noResolutions")
	}

	// Process resolutions
	result, err := h.service.ResolveConflictsBatch(ctx, req.Resolutions)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to resolve conflicts")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.resolveConflictsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.conflictsResolved", result)
}

// GetVersionHistory godoc
// @Summary Get version history for an entity
// @Description Returns version history for all translations of an entity
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param entity path string true "Entity type (translation, category, listing)"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.VersionHistoryResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/versions/{entity}/{id} [get]
// GetTranslationVersions godoc
// @Summary Get versions for a specific translation
// @Description Returns all versions for a specific translation ID
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param id path int true "Translation ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.TranslationVersion}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/versions/translation/{id} [get]
func (h *Handler) GetTranslationVersions(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	translationID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidTranslationID")
	}

	versions, err := h.service.translationRepo.GetTranslationVersions(ctx, translationID)
	if err != nil {
		h.logger.Error().Err(err).Int("translation_id", translationID).Msg("Failed to get translation versions")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.getVersionsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.versionsRetrieved", versions)
}

func (h *Handler) GetVersionHistory(c *fiber.Ctx) error {
	ctx := c.Context()

	entityType := c.Params("entity")
	idStr := c.Params("id")

	if entityType == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidEntityType")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidEntityID")
	}

	history, err := h.service.GetVersionHistory(ctx, entityType, id)
	if err != nil {
		h.logger.Error().Err(err).Str("entity", entityType).Int("id", id).Msg("Failed to get version history")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.versionHistoryError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.versionHistoryRetrieved", history)
}

// RollbackVersion godoc
// @Summary Rollback translation to previous version
// @Description Rollback a translation to a specific previous version
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.RollbackRequest true "Rollback request"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/versions/rollback [post]
func (h *Handler) RollbackVersion(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.RollbackRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	if req.TranslationID <= 0 || req.VersionID <= 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRollbackRequest")
	}

	err := h.service.RollbackVersion(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Int("translation_id", req.TranslationID).Int("version_id", req.VersionID).Msg("Failed to rollback version")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.rollbackError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.rollbackSuccess", nil)
}

// GetVersionDiff godoc
// @Summary Get differences between two versions
// @Description Compare two translation versions and return differences
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param version1 query int true "First version ID"
// @Param version2 query int true "Second version ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.VersionDiff}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/versions/diff [get]
func (h *Handler) GetVersionDiff(c *fiber.Ctx) error {
	ctx := c.Context()

	version1 := c.QueryInt("version1", 0)
	version2 := c.QueryInt("version2", 0)

	if version1 <= 0 || version2 <= 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidVersionIDs")
	}

	diff, err := h.service.GetVersionDiff(ctx, version1, version2)
	if err != nil {
		h.logger.Error().Err(err).Int("version1", version1).Int("version2", version2).Msg("Failed to get version diff")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.versionDiffError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.versionDiffRetrieved", diff)
}

func (h *Handler) GetCoverage(c *fiber.Ctx) error {
	// TODO: Implement get coverage
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.notImplemented", nil)
}

func (h *Handler) GetQuality(c *fiber.Ctx) error {
	// TODO: Implement get quality
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.notImplemented", nil)
}

func (h *Handler) GetUsage(c *fiber.Ctx) error {
	// TODO: Implement get usage
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.notImplemented", nil)
}

// ApplyAITranslations godoc
// @Summary Apply AI-generated translations
// @Description Apply translations generated by AI providers
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param translations body backend_internal_domain_models.ApplyTranslationsRequest true "Translations to apply"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/apply [post]
func (h *Handler) ApplyAITranslations(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.ApplyTranslationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	if err := h.service.ApplyAITranslations(ctx, &req, userID); err != nil {
		h.logger.Error().Err(err).Msg("Failed to apply AI translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.applyTranslationsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.translationsApplied", nil)
}

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Retrieve audit logs with optional filters
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param user_id query int false "Filter by user ID"
// @Param action query string false "Filter by action type"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.TranslationAuditLog}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/audit/logs [get]
func (h *Handler) GetAuditLogs(c *fiber.Ctx) error {
	ctx := c.Context()

	// Build filters from query parameters
	filters := make(map[string]interface{})

	if limit := c.QueryInt("limit", 100); limit > 0 {
		filters["limit"] = limit
	}
	if userID := c.QueryInt("user_id", 0); userID > 0 {
		filters["user_id"] = userID
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}

	logs, err := h.service.GetAuditLogs(ctx, filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get audit logs")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.auditLogsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.auditLogsRetrieved", logs)
}

// GetAuditStatistics godoc
// @Summary Get audit statistics
// @Description Retrieve statistics about audit logs and user activity
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.AuditStatistics}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/audit/statistics [get]
func (h *Handler) GetAuditStatistics(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.service.GetAuditStatistics(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get audit statistics")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.auditStatisticsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.auditStatisticsRetrieved", stats)
}

// ExportAdvanced godoc
// @Summary Advanced export translations
// @Description Export translations in various formats (JSON, CSV, XLIFF)
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.ExportRequest true "Export request"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/export/advanced [post]
func (h *Handler) ExportAdvanced(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.ExportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Set default format if not specified
	if req.Format == "" {
		req.Format = models.ExportFormatJSON
	}

	data, err := h.service.ExportTranslationsAdvanced(ctx, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("format", string(req.Format)).Msg("Failed to export translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.exportAdvancedError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.exportAdvancedSuccess", data)
}

// ImportAdvanced godoc
// @Summary Advanced import translations
// @Description Import translations from various formats (JSON, CSV, XLIFF)
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.TranslationImportRequest true "Import request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.ImportResult}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/import/advanced [post]
func (h *Handler) ImportAdvanced(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.TranslationImportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Set default format if not specified
	if req.Format == "" {
		req.Format = models.ExportFormatJSON
	}

	result, err := h.service.ImportTranslationsAdvanced(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).Str("format", string(req.Format)).Msg("Failed to import translations")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.importAdvancedError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.importAdvancedSuccess", result)
}

// BulkTranslate godoc
// @Summary Bulk translate multiple entities
// @Description Perform bulk translation of multiple entities across languages
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.BulkTranslateRequest true "Bulk translation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.BatchTranslateResult}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/bulk/translate [post]
func (h *Handler) BulkTranslate(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get user ID from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.SendError(c, fiber.StatusUnauthorized, "admin.translations.unauthorized")
	}

	var req models.BulkTranslateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidRequest")
	}

	// Validate request
	if req.EntityType == "" || req.SourceLanguage == "" || len(req.TargetLanguages) == 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidBulkRequest")
	}

	result, err := h.service.BulkTranslate(ctx, &req, userID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("entity_type", req.EntityType).
			Str("source_lang", req.SourceLanguage).
			Strs("target_langs", req.TargetLanguages).
			Msg("Failed to perform bulk translation")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.bulkTranslateError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.bulkTranslateSuccess", result)
}

// GetTranslationProviders godoc
// @Summary Get all translation providers
// @Description Returns list of available translation providers (Claude, DeepL, etc.)
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.TranslationProvider}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/providers [get]
func (h *Handler) GetTranslationProviders(c *fiber.Ctx) error {
	// Возвращаем пустой массив провайдеров пока они не настроены в системе
	// В будущем здесь можно будет возвращать реальных провайдеров из БД или конфигурации
	providers := []map[string]interface{}{}

	// Можно добавить статические провайдеры для демонстрации
	// providers = append(providers, map[string]interface{}{
	//     "id": 1,
	//     "name": "Claude AI",
	//     "provider_type": "claude",
	//     "is_active": true,
	//     "api_key": "configured",
	// })

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providersRetrieved", providers)
}

// UpdateTranslationProvider godoc
// @Summary Update translation provider configuration
// @Description Updates provider settings and configuration
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Param provider body backend_internal_domain_models.TranslationProvider true "Provider configuration"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/providers/{id} [put]
func (h *Handler) UpdateTranslationProvider(c *fiber.Ctx) error {
	idStr := c.Params("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidProviderID")
	}

	// Пока просто возвращаем успех, так как провайдеры еще не реализованы
	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providerUpdated", nil)
}
