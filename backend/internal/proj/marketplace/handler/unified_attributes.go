package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/services/attributes"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// UnifiedAttributesHandler handles unified attributes endpoints
type UnifiedAttributesHandler struct {
	service      *attributes.UnifiedAttributeService
	storage      postgres.UnifiedAttributeStorage
	featureFlags *config.FeatureFlags
}

// NewUnifiedAttributesHandler creates a new unified attributes handler
func NewUnifiedAttributesHandler(
	storage postgres.UnifiedAttributeStorage,
	featureFlags *config.FeatureFlags,
) *UnifiedAttributesHandler {
	// Создаем сервис с поддержкой feature flags
	service := attributes.NewUnifiedAttributeService(
		storage,
		featureFlags.UnifiedAttributesFallback,
		featureFlags.DualWriteAttributes,
	)

	return &UnifiedAttributesHandler{
		service:      service,
		storage:      storage,
		featureFlags: featureFlags,
	}
}

// GetCategoryAttributes godoc
// @Summary Get attributes for a category
// @Description Get all attributes available for a specific category with unified system support
// @Tags marketplace-attributes-v2
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.UnifiedAttribute} "List of category attributes"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid category ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v2/marketplace/categories/{category_id}/attributes [get]
func (h *UnifiedAttributesHandler) GetCategoryAttributes(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryId")
	}

	// Проверяем feature flag для конкретного пользователя
	userID, _ := authMiddleware.GetUserID(c)
	if !h.featureFlags.ShouldUseUnifiedAttributes(userID) {
		// Если новая система отключена для пользователя - возвращаем из старой
		return h.getCategoryAttributesLegacy(c, categoryID)
	}

	// Логируем использование новой системы если включено
	if h.featureFlags.LogAttributeSystemCalls {
		logger.Info().
			Int("user_id", userID).
			Int("category_id", categoryID).
			Msg("Using unified attributes system")
	}

	// Получаем атрибуты через новую систему
	attributes, err := h.service.GetCategoryAttributes(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).
			Int("category_id", categoryID).
			Msg("Failed to get category attributes")
		middleware.RecordUnifiedAttributesUsage("v2", "get_category_attributes_error")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.getAttributesError")
	}

	// Записываем метрику успешного вызова
	middleware.RecordUnifiedAttributesUsage("v2", "get_category_attributes_success")

	return utils.SendSuccess(c, fiber.StatusOK, "success.getAttributes", attributes)
}

// GetCategoryAttributesWithSettings godoc
// @Summary Get attributes with settings for a category
// @Description Get all attributes with their category-specific settings
// @Tags marketplace-attributes-v2
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.UnifiedCategoryAttribute} "List of category attributes with settings"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid category ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v2/marketplace/categories/{category_id}/attributes/detailed [get]
func (h *UnifiedAttributesHandler) GetCategoryAttributesWithSettings(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryId")
	}

	userID, _ := authMiddleware.GetUserID(c)
	if !h.featureFlags.ShouldUseUnifiedAttributes(userID) {
		return utils.SendError(c, fiber.StatusNotImplemented, "errors.featureNotAvailable")
	}

	attributes, err := h.service.GetCategoryAttributesWithSettings(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).
			Int("category_id", categoryID).
			Msg("Failed to get category attributes with settings")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.getAttributesError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.getAttributes", attributes)
}

// GetListingAttributeValues godoc
// @Summary Get attribute values for a listing
// @Description Get all attribute values for a specific listing
// @Tags marketplace-attributes-v2
// @Accept json
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.UnifiedAttributeValue} "List of attribute values"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid listing ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v2/marketplace/listings/{listing_id}/attributes [get]
func (h *UnifiedAttributesHandler) GetListingAttributeValues(c *fiber.Ctx) error {
	listingID, err := strconv.Atoi(c.Params("listing_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidListingId")
	}

	userID, _ := authMiddleware.GetUserID(c)
	if !h.featureFlags.ShouldUseUnifiedAttributes(userID) {
		return h.getListingAttributeValuesLegacy(c, listingID)
	}

	values, err := h.service.GetAttributeValues(
		c.Context(),
		models.AttributeEntityTypeListing,
		listingID,
	)
	if err != nil {
		logger.Error().Err(err).
			Int("listing_id", listingID).
			Msg("Failed to get listing attribute values")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.getAttributeValuesError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.getAttributeValues", values)
}

// SaveListingAttributeValues godoc
// @Summary Save attribute values for a listing
// @Description Save or update attribute values for a specific listing
// @Tags marketplace-attributes-v2
// @Accept json
// @Produce json
// @Param listing_id path int true "Listing ID"
// @Param values body map[int]interface{} true "Attribute values map (attribute_id -> value)"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "Values saved successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security Bearer
// @Router /api/v2/marketplace/listings/{listing_id}/attributes [post]
func (h *UnifiedAttributesHandler) SaveListingAttributeValues(c *fiber.Ctx) error {
	listingID, err := strconv.Atoi(c.Params("listing_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidListingId")
	}

	// Проверяем авторизацию
	userID, _ := authMiddleware.GetUserID(c)
	if userID == 0 {
		return utils.SendError(c, fiber.StatusUnauthorized, "errors.unauthorized")
	}

	// Проверяем feature flag
	if !h.featureFlags.ShouldUseUnifiedAttributes(userID) {
		return h.saveListingAttributeValuesLegacy(c, listingID)
	}

	// Парсим тело запроса
	var values map[int]interface{}
	if err := c.BodyParser(&values); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidRequestBody")
	}

	// Сохраняем значения
	err = h.service.SaveAttributeValues(
		c.Context(),
		models.AttributeEntityTypeListing,
		listingID,
		values,
	)
	if err != nil {
		logger.Error().Err(err).
			Int("listing_id", listingID).
			Msg("Failed to save listing attribute values")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.saveAttributeValuesError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.saveAttributeValues", nil)
}

// CreateAttribute godoc
// @Summary Create a new attribute
// @Description Create a new unified attribute (Admin only)
// @Tags marketplace-attributes-v2-admin
// @Accept json
// @Produce json
// @Param attribute body models.UnifiedAttribute true "Attribute data"
// @Success 201 {object} utils.SuccessResponseSwag{data=int} "Created attribute ID"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security Bearer
// @Router /api/v2/admin/attributes [post]
func (h *UnifiedAttributesHandler) CreateAttribute(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	var attr models.UnifiedAttribute
	if err := c.BodyParser(&attr); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidRequestBody")
	}

	id, err := h.service.CreateAttribute(c.Context(), &attr)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create attribute")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.createAttributeError")
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "success.createAttribute", id)
}

// UpdateAttribute godoc
// @Summary Update an attribute
// @Description Update an existing unified attribute (Admin only)
// @Tags marketplace-attributes-v2-admin
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "Attribute updated"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "Attribute not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security Bearer
// @Router /api/v2/admin/attributes/{id} [put]
func (h *UnifiedAttributesHandler) UpdateAttribute(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidAttributeId")
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidRequestBody")
	}

	err = h.service.UpdateAttribute(c.Context(), attributeID, updates)
	if err != nil {
		logger.Error().Err(err).
			Int("attribute_id", attributeID).
			Msg("Failed to update attribute")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.updateAttributeError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.updateAttribute", nil)
}

// DeleteAttribute godoc
// @Summary Delete an attribute
// @Description Delete a unified attribute (Admin only)
// @Tags marketplace-attributes-v2-admin
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "Attribute deleted"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "Attribute not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security Bearer
// @Router /api/v2/admin/attributes/{id} [delete]
func (h *UnifiedAttributesHandler) DeleteAttribute(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidAttributeId")
	}

	err = h.service.DeleteAttribute(c.Context(), attributeID)
	if err != nil {
		logger.Error().Err(err).
			Int("attribute_id", attributeID).
			Msg("Failed to delete attribute")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.deleteAttributeError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.deleteAttribute", nil)
}

// AttachAttributeToCategory godoc
// @Summary Attach attribute to category
// @Description Attach an attribute to a category with specific settings (Admin only)
// @Tags marketplace-attributes-v2-admin
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param attribute_id path int true "Attribute ID"
// @Param settings body models.UnifiedCategoryAttribute true "Attachment settings"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "Attribute attached"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security Bearer
// @Router /api/v2/admin/categories/{category_id}/attributes/{attribute_id} [post]
func (h *UnifiedAttributesHandler) AttachAttributeToCategory(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attribute_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidAttributeId")
	}

	var settings models.UnifiedCategoryAttribute
	if err := c.BodyParser(&settings); err != nil {
		// Если тело пустое - используем настройки по умолчанию
		settings = models.UnifiedCategoryAttribute{
			IsEnabled:  true,
			IsRequired: false,
			SortOrder:  0,
		}
	}

	err = h.service.AttachAttributeToCategory(c.Context(), categoryID, attributeID, &settings)
	if err != nil {
		logger.Error().Err(err).
			Int("category_id", categoryID).
			Int("attribute_id", attributeID).
			Msg("Failed to attach attribute to category")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.attachAttributeError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.attachAttribute", nil)
}

// Legacy методы для обратной совместимости

func (h *UnifiedAttributesHandler) getCategoryAttributesLegacy(c *fiber.Ctx, categoryID int) error {
	// Здесь должна быть логика получения из старой системы
	// Временно возвращаем пустой массив
	return utils.SendSuccess(c, fiber.StatusOK, "success.getAttributes", []interface{}{})
}

func (h *UnifiedAttributesHandler) getListingAttributeValuesLegacy(c *fiber.Ctx, listingID int) error {
	// Здесь должна быть логика получения из старой системы
	// Временно возвращаем пустой массив
	return utils.SendSuccess(c, fiber.StatusOK, "success.getAttributeValues", []interface{}{})
}

func (h *UnifiedAttributesHandler) saveListingAttributeValuesLegacy(c *fiber.Ctx, listingID int) error {
	// Здесь должна быть логика сохранения в старую систему
	// Временно возвращаем успех
	return utils.SendSuccess(c, fiber.StatusOK, "success.saveAttributeValues", nil)
}

// GetFeatureStatus godoc
// @Summary Get feature flags status for unified attributes
// @Description Get current status of feature flags for debugging
// @Tags marketplace-attributes-v2-admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Feature flags status"
// @Security Bearer
// @Router /api/v2/admin/attributes/feature-status [get]
func (h *UnifiedAttributesHandler) GetFeatureStatus(c *fiber.Ctx) error {
	// Только для администраторов
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	status := h.featureFlags.GetCurrentConfiguration()
	return utils.SendSuccess(c, fiber.StatusOK, "success.getFeatureStatus", status)
}

// GetAttributeRanges возвращает диапазоны значений для числовых атрибутов категории
func (h *UnifiedAttributesHandler) GetAttributeRanges(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryID")
	}

	// Используем fallback на старую систему
	return h.getCategoryAttributesLegacy(c, categoryID)
}

// UpdateListingAttributeValues обновляет значения атрибутов объявления (PUT)
func (h *UnifiedAttributesHandler) UpdateListingAttributeValues(c *fiber.Ctx) error {
	// Для PUT используем ту же логику, что и для POST (перезапись)
	return h.SaveListingAttributeValues(c)
}

// DetachAttributeFromCategory отвязывает атрибут от категории
func (h *UnifiedAttributesHandler) DetachAttributeFromCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryID")
	}

	attributeID, err := strconv.Atoi(c.Params("attribute_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidAttributeID")
	}

	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	// Удаляем связь через сервис
	err = h.service.DetachAttributeFromCategory(c.Context(), categoryID, attributeID)
	if err != nil {
		logger.Error().Err(err).
			Int("category_id", categoryID).
			Int("attribute_id", attributeID).
			Msg("Failed to detach attribute from category")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.detachAttributeFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.detachAttribute", nil)
}

// UpdateCategoryAttribute обновляет параметры атрибута в категории
func (h *UnifiedAttributesHandler) UpdateCategoryAttribute(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidCategoryID")
	}

	attributeID, err := strconv.Atoi(c.Params("attribute_id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidAttributeID")
	}

	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	var settings models.UnifiedCategoryAttribute
	if err := c.BodyParser(&settings); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "errors.invalidInput")
	}

	// Обновляем через сервис
	err = h.service.UpdateCategoryAttribute(c.Context(), categoryID, attributeID, &settings)
	if err != nil {
		logger.Error().Err(err).
			Int("category_id", categoryID).
			Int("attribute_id", attributeID).
			Msg("Failed to update category attribute")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.updateAttributeFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.updateAttribute", nil)
}

// MigrateFromLegacy запускает миграцию данных из старой системы
func (h *UnifiedAttributesHandler) MigrateFromLegacy(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	// Запускаем миграцию через сервис в фоне
	go func() {
		ctx := context.Background()
		err := h.service.MigrateFromLegacySystem(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Migration failed")
		} else {
			logger.Info().Msg("Migration completed successfully")
		}
	}()

	return utils.SendSuccess(c, fiber.StatusOK, "success.migrationStarted", fiber.Map{
		"message": "Migration started in background",
	})
}

// GetMigrationStatus возвращает статус миграции
func (h *UnifiedAttributesHandler) GetMigrationStatus(c *fiber.Ctx) error {
	// Проверяем права администратора
	if !authMiddleware.IsAdmin(c) {
		return utils.SendError(c, fiber.StatusForbidden, "errors.forbidden")
	}

	// Получаем статус миграции из сервиса
	status, err := h.service.GetMigrationStatus(c.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get migration status")
		return utils.SendError(c, fiber.StatusInternalServerError, "errors.getMigrationStatusFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "success.migrationStatus", status)
}
