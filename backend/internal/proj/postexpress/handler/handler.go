package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/interfaces"
	"backend/internal/middleware"
	"backend/pkg/utils"
	"backend/internal/proj/postexpress/models"
	"backend/internal/proj/postexpress/service"
	"backend/internal/proj/postexpress/storage"
	"backend/pkg/logger"
)

// Handler представляет HTTP обработчик для Post Express
type Handler struct {
	service service.Service
	logger  logger.Logger
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(service service.Service, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes регистрирует маршруты (реализация интерфейса RouteRegistrar)
func (h *Handler) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	postExpress := app.Group("/api/v1/postexpress")

	// Health check
	postExpress.Get("/health", h.Health)

	// Настройки
	postExpress.Get("/settings", h.GetSettings)
	postExpress.Put("/settings", h.UpdateSettings)

	// Локации и отделения
	postExpress.Get("/locations/search", h.SearchLocations)
	postExpress.Get("/locations/:id", h.GetLocation)
	postExpress.Post("/locations/sync", h.SyncLocations)
	postExpress.Get("/offices", h.GetOffices)
	postExpress.Get("/offices/:code", h.GetOffice)
	postExpress.Post("/offices/sync", h.SyncOffices)

	// Тарифы и расчет стоимости
	postExpress.Post("/calculate-rate", h.CalculateRate)
	postExpress.Get("/rates", h.GetRates)

	// Отправления
	postExpress.Post("/shipments", h.CreateShipment)
	postExpress.Get("/shipments", h.ListShipments)
	postExpress.Get("/shipments/:id", h.GetShipment)
	postExpress.Put("/shipments/:id/status", h.UpdateShipmentStatus)
	postExpress.Post("/shipments/:id/cancel", h.CancelShipment)

	// Документы
	postExpress.Get("/shipments/:id/label", h.GetShipmentLabel)
	postExpress.Get("/shipments/:id/invoice", h.GetShipmentInvoice)

	// Отслеживание
	postExpress.Get("/track/:tracking", h.TrackShipment)
	postExpress.Post("/track/sync", h.SyncAllShipments)

	// Склад и самовывоз
	warehouse := postExpress.Group("/warehouse")
	warehouse.Get("/", h.GetWarehouses)
	warehouse.Get("/:code", h.GetWarehouse)
	warehouse.Post("/pickup-orders", h.CreatePickupOrder)
	warehouse.Get("/pickup-orders/:id", h.GetPickupOrder)
	warehouse.Get("/pickup-orders/code/:code", h.GetPickupOrderByCode)
	warehouse.Post("/pickup-orders/:id/confirm", h.ConfirmPickup)
	warehouse.Post("/pickup-orders/:id/cancel", h.CancelPickupOrder)

	// Статистика
	postExpress.Get("/statistics/shipments", h.GetShipmentStatistics)
	postExpress.Get("/statistics/warehouse/:id", h.GetWarehouseStatistics)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования (реализация интерфейса RouteRegistrar)
func (h *Handler) GetPrefix() string {
	return "postexpress"
}

// =============================================================================
// HEALTH CHECK
// =============================================================================

// Health проверяет состояние Post Express сервиса
// @Summary Health check для Post Express
// @Description Проверяет доступность и состояние сервиса Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=map[string]interface{}}
// @Router /api/v1/postexpress/health [get]
func (h *Handler) Health(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, map[string]interface{}{
		"status":  "ok",
		"service": "postexpress",
		"message": "Post Express service is healthy",
	})
}

// =============================================================================
// НАСТРОЙКИ
// =============================================================================

// GetSettings получает настройки Post Express
// @Summary Get Post Express settings
// @Description Получить текущие настройки интеграции с Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressSettings}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/settings [get]
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetSettings(c.Context())
	if err != nil {
		h.logger.Error("Failed to get Post Express settings: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getSettingsError")
	}

	return utils.SuccessResponse(c, settings)
}

// UpdateSettings обновляет настройки Post Express
// @Summary Update Post Express settings
// @Description Обновить настройки интеграции с Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param settings body backend_internal_proj_postexpress_models.PostExpressSettings true "Настройки Post Express"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressSettings}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/settings [put]
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	var settings models.PostExpressSettings
	if err := c.BodyParser(&settings); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	err := h.service.UpdateSettings(c.Context(), &settings)
	if err != nil {
		h.logger.Error("Failed to update Post Express settings: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.updateSettingsError")
	}

	return utils.SuccessResponse(c, settings)
}

// =============================================================================
// ЛОКАЦИИ И ОТДЕЛЕНИЯ
// =============================================================================

// SearchLocations ищет населенные пункты
// @Summary Search locations
// @Description Поиск населенных пунктов Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_postexpress_models.PostExpressLocation}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/locations/search [get]
func (h *Handler) SearchLocations(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.queryRequired")
	}

	locations, err := h.service.SearchLocations(c.Context(), query)
	if err != nil {
		h.logger.Error("Failed to search locations - query: %s, error: %v", query, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.searchLocationsError")
	}

	return utils.SuccessResponse(c, locations)
}

// GetLocation получает информацию о населенном пункте
// @Summary Get location
// @Description Получить информацию о населенном пункте по ID
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID населенного пункта"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressLocation}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/locations/{id} [get]
func (h *Handler) GetLocation(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	location, err := h.service.GetLocationByID(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get location - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getLocationError")
	}

	if location == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.locationNotFound")
	}

	return utils.SuccessResponse(c, location)
}

// SyncLocations синхронизирует населенные пункты с Post Express
// @Summary Sync locations
// @Description Синхронизировать населенные пункты с API Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/locations/sync [post]
func (h *Handler) SyncLocations(c *fiber.Ctx) error {
	err := h.service.SyncLocations(c.Context())
	if err != nil {
		h.logger.Error("Failed to sync locations: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.syncLocationsError")
	}

	return utils.SuccessResponse(c, "Locations synced successfully")
}

// GetOffices получает список отделений
// @Summary Get offices
// @Description Получить список почтовых отделений для населенного пункта
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param location_id query int true "ID населенного пункта"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_postexpress_models.PostExpressOffice}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/offices [get]
func (h *Handler) GetOffices(c *fiber.Ctx) error {
	locationID, err := strconv.Atoi(c.Query("location_id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.locationIDRequired")
	}

	offices, err := h.service.GetOfficesByLocation(c.Context(), locationID)
	if err != nil {
		h.logger.Error("Failed to get offices - location_id: %d, error: %v", locationID, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getOfficesError")
	}

	return utils.SuccessResponse(c, offices)
}

// GetOffice получает информацию об отделении
// @Summary Get office
// @Description Получить информацию о почтовом отделении по коду
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param code path string true "Код отделения"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressOffice}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/offices/{code} [get]
func (h *Handler) GetOffice(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.officeCodeRequired")
	}

	office, err := h.service.GetOfficeByCode(c.Context(), code)
	if err != nil {
		h.logger.Error("Failed to get office - code: %s, error: %v", code, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getOfficeError")
	}

	if office == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.officeNotFound")
	}

	return utils.SuccessResponse(c, office)
}

// SyncOffices синхронизирует отделения с Post Express
// @Summary Sync offices
// @Description Синхронизировать почтовые отделения с API Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/offices/sync [post]
func (h *Handler) SyncOffices(c *fiber.Ctx) error {
	err := h.service.SyncOffices(c.Context())
	if err != nil {
		h.logger.Error("Failed to sync offices: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.syncOfficesError")
	}

	return utils.SuccessResponse(c, "Offices synced successfully")
}

// =============================================================================
// ТАРИФЫ И РАСЧЕТ СТОИМОСТИ
// =============================================================================

// CalculateRate рассчитывает стоимость доставки
// @Summary Calculate delivery rate
// @Description Рассчитать стоимость доставки Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_postexpress_models.CalculateRateRequest true "Параметры для расчета стоимости"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.CalculateRateResponse}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/calculate-rate [post]
func (h *Handler) CalculateRate(c *fiber.Ctx) error {
	var req models.CalculateRateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	// Простая валидация
	if req.WeightKg <= 0 || req.WeightKg > 20 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.invalidWeight")
	}

	response, err := h.service.CalculateRate(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to calculate rate - weight: %f, error: %v", req.WeightKg, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.calculateRateError")
	}

	return utils.SuccessResponse(c, response)
}

// GetRates получает список тарифов
// @Summary Get delivery rates
// @Description Получить список тарифов доставки Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_postexpress_models.PostExpressRate}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/rates [get]
func (h *Handler) GetRates(c *fiber.Ctx) error {
	rates, err := h.service.GetRates(c.Context())
	if err != nil {
		h.logger.Error("Failed to get rates: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getRatesError")
	}

	return utils.SuccessResponse(c, rates)
}

// =============================================================================
// ОТПРАВЛЕНИЯ
// =============================================================================

// CreateShipment создает новое отправление
// @Summary Create shipment
// @Description Создать новое отправление Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_postexpress_models.CreateShipmentRequest true "Данные для создания отправления"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressShipment}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments [post]
func (h *Handler) CreateShipment(c *fiber.Ctx) error {
	var req models.CreateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	// Валидация
	if req.RecipientName == "" || req.RecipientAddress == "" || req.RecipientPhone == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.recipientInfoRequired")
	}

	if req.WeightKg <= 0 || req.WeightKg > 20 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.invalidWeight")
	}

	shipment, err := h.service.CreateShipment(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create shipment: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.createShipmentError")
	}

	return utils.SuccessResponse(c, shipment)
}

// ListShipments получает список отправлений
// @Summary List shipments
// @Description Получить список отправлений с фильтрацией и пагинацией
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param status query string false "Статус отправления"
// @Param date_from query string false "Дата от (YYYY-MM-DD)"
// @Param date_to query string false "Дата до (YYYY-MM-DD)"
// @Param city query string false "Город получателя"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(20)
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=object}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments [get]
func (h *Handler) ListShipments(c *fiber.Ctx) error {
	// Парсим фильтры
	filters := storage.ShipmentFilters{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if status := c.Query("status"); status != "" {
		shipmentStatus := models.ShipmentStatus(status)
		filters.Status = &shipmentStatus
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters.DateTo = &dateTo
	}

	if city := c.Query("city"); city != "" {
		filters.City = &city
	}

	if orderID := c.QueryInt("order_id", 0); orderID > 0 {
		filters.OrderID = &orderID
	}

	shipments, total, err := h.service.ListShipments(c.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to list shipments: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.listShipmentsError")
	}

	response := map[string]interface{}{
		"shipments": shipments,
		"total":     total,
		"page":      filters.Page,
		"page_size": filters.PageSize,
		"pages":     (total + filters.PageSize - 1) / filters.PageSize,
	}

	return utils.SuccessResponse(c, response)
}

// GetShipment получает информацию об отправлении
// @Summary Get shipment
// @Description Получить информацию об отправлении по ID
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID отправления"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.PostExpressShipment}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments/{id} [get]
func (h *Handler) GetShipment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	shipment, err := h.service.GetShipment(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get shipment - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getShipmentError")
	}

	if shipment == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.shipmentNotFound")
	}

	return utils.SuccessResponse(c, shipment)
}

// UpdateShipmentStatus обновляет статус отправления
// @Summary Update shipment status
// @Description Обновить статус отправления
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID отправления"
// @Param request body object true "Новый статус" example({"status":"delivered"})
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments/{id}/status [put]
func (h *Handler) UpdateShipmentStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	var req struct {
		Status models.ShipmentStatus `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	err = h.service.UpdateShipmentStatus(c.Context(), id, req.Status)
	if err != nil {
		h.logger.Error("Failed to update shipment status - id: %d, status: %s, error: %v", id, req.Status, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.updateStatusError")
	}

	return utils.SuccessResponse(c, "Status updated successfully")
}

// CancelShipment отменяет отправление
// @Summary Cancel shipment
// @Description Отменить отправление
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID отправления"
// @Param request body object true "Причина отмены" example({"reason":"Отмена по требованию клиента"})
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments/{id}/cancel [post]
func (h *Handler) CancelShipment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	err = h.service.CancelShipment(c.Context(), id, req.Reason)
	if err != nil {
		h.logger.Error("Failed to cancel shipment - id: %d, reason: %s, error: %v", id, req.Reason, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.cancelShipmentError")
	}

	return utils.SuccessResponse(c, "Shipment canceled successfully")
}

// =============================================================================
// ДОКУМЕНТЫ
// =============================================================================

// GetShipmentLabel получает этикетку отправления
// @Summary Get shipment label
// @Description Получить этикетку отправления в формате PDF
// @Tags PostExpress
// @Accept json
// @Produce application/pdf
// @Param id path int true "ID отправления"
// @Success 200 {file} file "PDF этикетка"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments/{id}/label [get]
func (h *Handler) GetShipmentLabel(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	labelData, err := h.service.GetShipmentLabel(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get shipment label - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getLabelError")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=label.pdf")
	return c.Send(labelData)
}

// GetShipmentInvoice получает накладную отправления
// @Summary Get shipment invoice
// @Description Получить накладную отправления в формате PDF
// @Tags PostExpress
// @Accept json
// @Produce application/pdf
// @Param id path int true "ID отправления"
// @Success 200 {file} file "PDF накладная"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/shipments/{id}/invoice [get]
func (h *Handler) GetShipmentInvoice(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	invoiceData, err := h.service.GetShipmentInvoice(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get shipment invoice - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getInvoiceError")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=invoice.pdf")
	return c.Send(invoiceData)
}

// =============================================================================
// ОТСЛЕЖИВАНИЕ
// =============================================================================

// TrackShipment отслеживает отправление
// @Summary Track shipment
// @Description Получить события отслеживания отправления по трекинг-номеру
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param tracking path string true "Трекинг-номер"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_postexpress_models.TrackingEvent}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/track/{tracking} [get]
func (h *Handler) TrackShipment(c *fiber.Ctx) error {
	trackingNumber := c.Params("tracking")
	if trackingNumber == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.trackingNumberRequired")
	}

	events, err := h.service.TrackShipment(c.Context(), trackingNumber)
	if err != nil {
		h.logger.Error("Failed to track shipment - tracking: %s, error: %v", trackingNumber, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.trackShipmentError")
	}

	return utils.SuccessResponse(c, events)
}

// SyncAllShipments синхронизирует все активные отправления
// @Summary Sync all active shipments
// @Description Синхронизировать статусы всех активных отправлений с Post Express
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/track/sync [post]
func (h *Handler) SyncAllShipments(c *fiber.Ctx) error {
	err := h.service.SyncAllActiveShipments(c.Context())
	if err != nil {
		h.logger.Error("Failed to sync all shipments: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.syncShipmentsError")
	}

	return utils.SuccessResponse(c, "All shipments synced successfully")
}

// =============================================================================
// СКЛАД И САМОВЫВОЗ
// =============================================================================

// GetWarehouses получает список складов
// @Summary Get warehouses
// @Description Получить список всех активных складов
// @Tags PostExpress
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_postexpress_models.Warehouse}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse [get]
func (h *Handler) GetWarehouses(c *fiber.Ctx) error {
	warehouses, err := h.service.GetWarehouses(c.Context())
	if err != nil {
		h.logger.Error("Failed to get warehouses: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getWarehousesError")
	}

	return utils.SuccessResponse(c, warehouses)
}

// GetWarehouse получает информацию о складе
// @Summary Get warehouse
// @Description Получить информацию о складе по коду
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param code path string true "Код склада"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.Warehouse}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/{code} [get]
func (h *Handler) GetWarehouse(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.warehouseCodeRequired")
	}

	warehouse, err := h.service.GetWarehouseByCode(c.Context(), code)
	if err != nil {
		h.logger.Error("Failed to get warehouse - code: %s, error: %v", code, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getWarehouseError")
	}

	if warehouse == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.warehouseNotFound")
	}

	return utils.SuccessResponse(c, warehouse)
}

// CreatePickupOrder создает заказ на самовывоз
// @Summary Create pickup order
// @Description Создать заказ на самовывоз со склада
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_postexpress_models.CreatePickupOrderRequest true "Данные для создания заказа на самовывоз"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.WarehousePickupOrder}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/pickup-orders [post]
func (h *Handler) CreatePickupOrder(c *fiber.Ctx) error {
	var req models.CreatePickupOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	// Валидация
	if req.CustomerName == "" || req.CustomerPhone == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.customerInfoRequired")
	}

	order, err := h.service.CreatePickupOrder(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create pickup order: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.createPickupOrderError")
	}

	return utils.SuccessResponse(c, order)
}

// GetPickupOrder получает информацию о заказе на самовывоз
// @Summary Get pickup order
// @Description Получить информацию о заказе на самовывоз по ID
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID заказа на самовывоз"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.WarehousePickupOrder}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/pickup-orders/{id} [get]
func (h *Handler) GetPickupOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	order, err := h.service.GetPickupOrder(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get pickup order - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getPickupOrderError")
	}

	if order == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.pickupOrderNotFound")
	}

	return utils.SuccessResponse(c, order)
}

// GetPickupOrderByCode получает информацию о заказе на самовывоз по коду
// @Summary Get pickup order by code
// @Description Получить информацию о заказе на самовывоз по коду
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param code path string true "Код самовывоза"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_postexpress_models.WarehousePickupOrder}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/pickup-orders/code/{code} [get]
func (h *Handler) GetPickupOrderByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "postexpress.pickupCodeRequired")
	}

	order, err := h.service.GetPickupOrderByCode(c.Context(), code)
	if err != nil {
		h.logger.Error("Failed to get pickup order by code - code: %s, error: %v", code, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getPickupOrderError")
	}

	if order == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.pickupOrderNotFound")
	}

	return utils.SuccessResponse(c, order)
}

// ConfirmPickup подтверждает выдачу заказа
// @Summary Confirm pickup
// @Description Подтвердить выдачу заказа со склада
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID заказа на самовывоз"
// @Param request body object true "Данные подтверждения" example({"confirmed_by":"Иван Петрович","document_type":"passport","document_number":"123456789"})
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/pickup-orders/{id}/confirm [post]
func (h *Handler) ConfirmPickup(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	var req struct {
		ConfirmedBy    string `json:"confirmed_by" validate:"required"`
		DocumentType   string `json:"document_type" validate:"required"`
		DocumentNumber string `json:"document_number" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	err = h.service.ConfirmPickup(c.Context(), id, req.ConfirmedBy, req.DocumentType, req.DocumentNumber)
	if err != nil {
		h.logger.Error("Failed to confirm pickup - id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.confirmPickupError")
	}

	return utils.SuccessResponse(c, "Pickup confirmed successfully")
}

// CancelPickupOrder отменяет заказ на самовывоз
// @Summary Cancel pickup order
// @Description Отменить заказ на самовывоз
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID заказа на самовывоз"
// @Param request body object true "Причина отмены" example({"reason":"Отмена по требованию клиента"})
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/warehouse/pickup-orders/{id}/cancel [post]
func (h *Handler) CancelPickupOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidRequestBody")
	}

	err = h.service.CancelPickupOrder(c.Context(), id, req.Reason)
	if err != nil {
		h.logger.Error("Failed to cancel pickup order - id: %d, reason: %s, error: %v", id, req.Reason, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.cancelPickupOrderError")
	}

	return utils.SuccessResponse(c, "Pickup order canceled successfully")
}

// =============================================================================
// СТАТИСТИКА
// =============================================================================

// GetShipmentStatistics получает статистику отправлений
// @Summary Get shipment statistics
// @Description Получить статистику отправлений с фильтрами
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param date_from query string false "Дата от (YYYY-MM-DD)"
// @Param date_to query string false "Дата до (YYYY-MM-DD)"
// @Param group_by query string false "Группировка (day, week, month)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=storage.ShipmentStatistics}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/statistics/shipments [get]
func (h *Handler) GetShipmentStatistics(c *fiber.Ctx) error {
	filters := storage.StatisticsFilters{}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters.DateTo = &dateTo
	}

	if groupBy := c.Query("group_by"); groupBy != "" {
		filters.GroupBy = groupBy
	}

	stats, err := h.service.GetShipmentStatistics(c.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to get shipment statistics: %v", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getStatisticsError")
	}

	return utils.SuccessResponse(c, stats)
}

// GetWarehouseStatistics получает статистику склада
// @Summary Get warehouse statistics
// @Description Получить статистику склада
// @Tags PostExpress
// @Accept json
// @Produce json
// @Param id path int true "ID склада"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=storage.WarehouseStatistics}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/postexpress/statistics/warehouse/{id} [get]
func (h *Handler) GetWarehouseStatistics(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "common.errors.invalidID")
	}

	stats, err := h.service.GetWarehouseStatistics(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get warehouse statistics - warehouse_id: %d, error: %v", id, err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "postexpress.getWarehouseStatisticsError")
	}

	if stats == nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "postexpress.warehouseNotFound")
	}

	return utils.SuccessResponse(c, stats)
}

// Ensure Handler implements RouteRegistrar interface
var _ interfaces.RouteRegistrar = (*Handler)(nil)
