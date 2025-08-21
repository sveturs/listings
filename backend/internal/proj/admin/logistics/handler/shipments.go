package handler

import (
	"strconv"
	"time"

	"backend/internal/domain/logistics"
	"backend/internal/proj/admin/logistics/service"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ShipmentsHandler обработчик для управления отправлениями
type ShipmentsHandler struct {
	monitoringService *service.MonitoringService
	problemService    *service.ProblemService
	logger            *logger.Logger
}

// NewShipmentsHandler создает новый обработчик отправлений
func NewShipmentsHandler(monitoringService *service.MonitoringService, problemService *service.ProblemService) *ShipmentsHandler {
	return &ShipmentsHandler{
		monitoringService: monitoringService,
		problemService:    problemService,
		logger:            logger.GetLogger(),
	}
}

// GetShipments godoc
// @Summary Получить список отправлений с фильтрами
// @Description Возвращает список отправлений с возможностью фильтрации и пагинации
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Фильтр по статусу"
// @Param courier_service query string false "Фильтр по курьерской службе (bex, postexpress)"
// @Param date_from query string false "Дата начала периода (RFC3339)"
// @Param date_to query string false "Дата окончания периода (RFC3339)"
// @Param city query string false "Фильтр по городу"
// @Param tracking_number query string false "Поиск по трек-номеру"
// @Param has_problems query bool false "Только проблемные отправления"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Param sort_by query string false "Поле для сортировки" default(created_at)
// @Param sort_order query string false "Порядок сортировки (asc, desc)" default(desc)
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "List of shipments with total count"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid filters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/shipments [get]
func (h *ShipmentsHandler) GetShipments(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим фильтры
	filter := logistics.ShipmentsFilter{
		Page:      c.QueryInt("page", 1),
		Limit:     c.QueryInt("limit", 20),
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	// Опциональные фильтры
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if courierService := c.Query("courier_service"); courierService != "" {
		filter.CourierService = &courierService
	}

	if city := c.Query("city"); city != "" {
		filter.City = &city
	}

	if trackingNumber := c.Query("tracking_number"); trackingNumber != "" {
		filter.TrackingNumber = &trackingNumber
	}

	if hasProblems := c.Query("has_problems"); hasProblems != "" {
		val := hasProblems == "true"
		filter.HasProblems = &val
	}

	// Парсим даты
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		t, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_from")
		}
		filter.DateFrom = &t
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		t, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_to")
		}
		filter.DateTo = &t
	}

	// Получаем отправления
	shipments, total, err := h.monitoringService.GetShipments(c.Context(), filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.shipments.error")
	}

	// Формируем ответ с метаданными пагинации
	response := map[string]interface{}{
		"shipments": shipments,
		"total":     total,
		"page":      filter.Page,
		"limit":     filter.Limit,
		"pages":     (total + filter.Limit - 1) / filter.Limit,
	}

	return utils.SuccessResponse(c, response)
}

// GetShipmentDetailsByProvider godoc
// @Summary Получить детальную информацию об отправлении по провайдеру
// @Description Возвращает полную информацию об отправлении включая историю статусов
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param provider path string true "Провайдер (BEX, PostExpress)"
// @Param id path int true "ID отправления"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Shipment details"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/shipments/{provider}/{id} [get]
func (h *ShipmentsHandler) GetShipmentDetailsByProvider(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	provider := c.Params("provider")
	if provider == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.provider_required")
	}

	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_shipment_id")
	}

	// Получаем детали отправления
	details, err := h.monitoringService.GetShipmentDetailsByProvider(c.Context(), provider, shipmentID)
	if err != nil {
		if err.Error() == "shipment not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.shipment_not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.error")
	}

	return utils.SuccessResponse(c, details)
}

// GetShipmentDetails godoc
// @Summary Получить детальную информацию об отправлении
// @Description Возвращает полную информацию об отправлении включая историю статусов и проблемы
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID отправления"
// @Param type query string true "Тип отправления (bex, postexpress)"
// @Success 200 {object} utils.SuccessResponseSwag{data=logistics.ShipmentDetails} "Shipment details"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/shipments/{id} [get]
func (h *ShipmentsHandler) GetShipmentDetails(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_shipment_id")
	}

	shipmentType := c.Query("type")
	if shipmentType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.shipment_type_required")
	}

	// Получаем детали отправления
	details, err := h.monitoringService.GetShipmentDetails(c.Context(), shipmentID, shipmentType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "logistics.shipment_not_found")
	}

	return utils.SuccessResponse(c, details)
}

// UpdateShipmentStatus godoc
// @Summary Обновить статус отправления
// @Description Изменяет статус отправления и добавляет запись в лог
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID отправления"
// @Param body body map[string]string true "Новый статус"
// @Success 200 {object} utils.SuccessResponseSwag "Status updated successfully"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden - insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/shipments/{id}/status [put]
func (h *ShipmentsHandler) UpdateShipmentStatus(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_shipment_id")
	}

	// Парсим тело запроса
	var request struct {
		Status       string `json:"status"`
		ShipmentType string `json:"shipment_type"`
		Comment      string `json:"comment"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	// Получаем ID администратора из контекста (мок для тестирования)
	adminID := 1 // TODO: Получать реальный ID администратора из JWT токена

	// Обновляем статус через сервис
	err = h.monitoringService.UpdateShipmentStatus(c.Context(), shipmentID, request.ShipmentType, request.Status, adminID, request.Comment)
	if err != nil {
		h.logger.Error("Failed to update shipment status", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.update_status_error")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"shipment_id": shipmentID,
		"new_status":  request.Status,
		"message":     "logistics.status_updated_successfully",
	})
}

// PerformShipmentAction godoc
// @Summary Выполнить действие над отправлением
// @Description Выполняет специальное действие (связаться с курьером, отправить уведомление и т.д.)
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID отправления"
// @Param body body map[string]interface{} true "Действие и параметры"
// @Success 200 {object} utils.SuccessResponseSwag "Action performed successfully"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden - insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/shipments/{id}/action [post]
func (h *ShipmentsHandler) PerformShipmentAction(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_shipment_id")
	}

	// Парсим тело запроса
	var request struct {
		Action       string                 `json:"action"`
		ShipmentType string                 `json:"shipment_type"`
		Parameters   map[string]interface{} `json:"parameters"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.invalid_request_body")
	}

	// Получаем ID администратора из контекста (мок для тестирования)
	adminID := 1 // TODO: Получать реальный ID администратора из JWT токена

	// Выполняем действие через сервис
	err = h.monitoringService.PerformShipmentAction(c.Context(), shipmentID, request.ShipmentType, request.Action, adminID, request.Parameters)
	if err != nil {
		h.logger.Error("Failed to perform shipment action", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.action_error")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"shipment_id": shipmentID,
		"action":      request.Action,
		"message":     "logistics.action_completed_successfully",
	})
}
