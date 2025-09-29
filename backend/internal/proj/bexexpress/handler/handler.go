package handler

import (
	"fmt"
	"strconv"

	"backend/internal/proj/bexexpress/models"
	"backend/internal/proj/bexexpress/service"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler представляет HTTP handler для BEX Express
type Handler struct {
	service *service.Service
	logger  *logger.Logger
}

// NewHandler создает новый handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
		logger:  logger.GetLogger(),
	}
}

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes(api fiber.Router) {
	bex := api.Group("/bex")

	// Создание отправления
	bex.Post("/shipments", h.CreateShipment)

	// Получение информации об отправлении
	bex.Get("/shipments/:id", h.GetShipment)

	// Получение статуса отправления
	bex.Get("/shipments/:id/status", h.GetShipmentStatus)

	// Получение этикетки для печати
	bex.Get("/shipments/:id/label", h.GetShipmentLabel)

	// Отмена отправления
	bex.Delete("/shipments/:id", h.CancelShipment)

	// Расчет стоимости доставки
	bex.Post("/calculate-rate", h.CalculateRate)

	// Поиск адреса
	bex.Post("/search-address", h.SearchAddress)

	// Получение списка пунктов выдачи
	bex.Get("/parcel-shops", h.GetParcelShops)

	// Отслеживание посылки по номеру
	bex.Get("/track/:tracking", h.TrackShipment)

	// Массовое создание отправлений
	bex.Post("/shipments/bulk", h.CreateBulkShipments)

	// Webhook для обновления статусов
	bex.Post("/webhook/status", h.HandleStatusWebhook)
}

// CreateShipment создает новое отправление
// @Summary Create BEX shipment
// @Description Create a new shipment through BEX Express
// @Tags BEX Express
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_bexexpress_models.CreateShipmentRequest true "Shipment details"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_bexexpress_models.BEXShipment}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments [post]
func (h *Handler) CreateShipment(c *fiber.Ctx) error {
	var req models.CreateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Validate request
	if req.RecipientName == "" || req.RecipientAddress == "" || req.RecipientCity == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.missingRequiredFields")
	}

	// Автоматически определяем категорию по весу если не указана
	if req.ShipmentCategory == 0 {
		req.ShipmentCategory = service.CalculateShipmentCategory(req.WeightKg)
	}

	// Устанавливаем содержимое по умолчанию
	if req.ShipmentContents == 0 {
		req.ShipmentContents = 3 // Mixed content
	}

	shipment, err := h.service.CreateShipment(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.shipmentCreationFailed")
	}

	return utils.SuccessResponse(c, shipment)
}

// GetShipment получает информацию об отправлении
// @Summary Get shipment details
// @Description Get shipment information by ID
// @Tags BEX Express
// @Produce json
// @Param id path int true "Shipment ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_bexexpress_models.BEXShipment}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments/{id} [get]
func (h *Handler) GetShipment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	shipment, err := h.service.GetShipmentByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.shipmentNotFound")
	}

	return utils.SuccessResponse(c, shipment)
}

// GetShipmentStatus получает статус отправления
// @Summary Get shipment status
// @Description Get current status of shipment
// @Tags BEX Express
// @Produce json
// @Param id path int true "Shipment ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_bexexpress_models.BEXShipment}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments/{id}/status [get]
func (h *Handler) GetShipmentStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	shipment, err := h.service.GetShipmentStatus(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.statusNotFound")
	}

	// Возвращаем только статусную информацию
	statusInfo := map[string]interface{}{
		"shipment_id":     shipment.ID,
		"tracking_number": shipment.TrackingNumber,
		"status":          shipment.Status,
		"status_text":     shipment.StatusText,
		"registered_at":   shipment.RegisteredAt,
		"picked_up_at":    shipment.PickedUpAt,
		"delivered_at":    shipment.DeliveredAt,
		"failed_at":       shipment.FailedAt,
		"returned_at":     shipment.ReturnedAt,
		"failed_reason":   shipment.FailedReason,
	}

	return utils.SuccessResponse(c, statusInfo)
}

// GetShipmentLabel получает этикетку для печати
// @Summary Get shipment label
// @Description Get printable label for shipment
// @Tags BEX Express
// @Produce application/pdf
// @Param id path int true "Shipment ID"
// @Param size query int false "Page size (4 for A4, 6 for A6)" default(4)
// @Success 200 {file} binary
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments/{id}/label [get]
func (h *Handler) GetShipmentLabel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	pageSize := 4 // Default A4
	if size := c.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil && (s == 4 || s == 6) {
			pageSize = s
		}
	}

	labelData, err := h.service.GetShipmentLabel(c.Context(), id, pageSize)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.labelNotFound")
	}

	// Return PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=bex-label-%d.pdf", id))
	return c.Send(labelData)
}

// CancelShipment отменяет отправление
// @Summary Cancel shipment
// @Description Cancel a shipment
// @Tags BEX Express
// @Param id path int true "Shipment ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments/{id} [delete]
func (h *Handler) CancelShipment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
	}

	err = h.service.CancelShipment(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.cancelFailed")
	}

	return utils.SuccessResponse(c, nil)
}

// CalculateRate рассчитывает стоимость доставки
// @Summary Calculate shipping rate
// @Description Calculate shipping cost for given parameters
// @Tags BEX Express
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_bexexpress_models.CalculateRateRequest true "Rate calculation parameters"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_bexexpress_models.CalculateRateResponse}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/bex/calculate-rate [post]
func (h *Handler) CalculateRate(c *fiber.Ctx) error {
	var req models.CalculateRateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	// Автоматически определяем категорию по весу если не указана
	if req.ShipmentCategory == 0 {
		req.ShipmentCategory = service.CalculateShipmentCategory(req.WeightKg)
	}

	rate, err := h.service.CalculateRate(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.calculationFailed")
	}

	return utils.SuccessResponse(c, rate)
}

// SearchAddress ищет адреса в справочниках BEX
// @Summary Search addresses
// @Description Search for addresses in BEX database
// @Tags BEX Express
// @Accept json
// @Produce json
// @Param request body backend_internal_proj_bexexpress_models.SearchAddressRequest true "Search parameters"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_bexexpress_models.AddressSuggestion}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/bex/search-address [post]
func (h *Handler) SearchAddress(c *fiber.Ctx) error {
	var req models.SearchAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if len(req.Query) < 2 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.queryTooShort")
	}

	suggestions, err := h.service.SearchAddress(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.searchFailed")
	}

	return utils.SuccessResponse(c, suggestions)
}

// GetParcelShops получает список пунктов выдачи
// @Summary Get parcel shops
// @Description Get list of BEX parcel shops
// @Tags BEX Express
// @Produce json
// @Param city query string false "Filter by city"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_bexexpress_models.BEXParcelShop}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/bex/parcel-shops [get]
func (h *Handler) GetParcelShops(c *fiber.Ctx) error {
	city := c.Query("city")

	shops, err := h.service.GetParcelShops(c.Context(), city)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.fetchFailed")
	}

	return utils.SuccessResponse(c, shops)
}

// TrackShipment отслеживает посылку по номеру
// @Summary Track shipment
// @Description Track shipment by tracking number
// @Tags BEX Express
// @Produce json
// @Param tracking path string true "Tracking number"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_proj_bexexpress_models.BEXShipment}
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag
// @Router /api/v1/bex/track/{tracking} [get]
func (h *Handler) TrackShipment(c *fiber.Ctx) error {
	tracking := c.Params("tracking")
	if tracking == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidTracking")
	}

	shipment, err := h.service.GetShipmentByTracking(c.Context(), tracking)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.shipmentNotFound")
	}

	// Обновляем статус из BEX
	if shipment.BexShipmentID != nil {
		shipment, _ = h.service.GetShipmentStatus(c.Context(), shipment.ID)
	}

	return utils.SuccessResponse(c, shipment)
}

// CreateBulkShipments создает несколько отправлений
// @Summary Create bulk shipments
// @Description Create multiple shipments at once
// @Tags BEX Express
// @Accept json
// @Produce json
// @Param request body []backend_internal_proj_bexexpress_models.CreateShipmentRequest true "Array of shipment details"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_bexexpress_models.BEXShipment}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Security ApiKeyAuth
// @Router /api/v1/bex/shipments/bulk [post]
func (h *Handler) CreateBulkShipments(c *fiber.Ctx) error {
	var requests []models.CreateShipmentRequest
	if err := c.BodyParser(&requests); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidRequest")
	}

	if len(requests) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.emptyRequest")
	}

	if len(requests) > 100 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.tooManyShipments")
	}

	var shipments []*models.BEXShipment
	var errors []string

	for i, req := range requests {
		// Автоматически определяем категорию по весу если не указана
		if req.ShipmentCategory == 0 {
			req.ShipmentCategory = service.CalculateShipmentCategory(req.WeightKg)
		}

		// Устанавливаем содержимое по умолчанию
		if req.ShipmentContents == 0 {
			req.ShipmentContents = 3 // Mixed content
		}

		shipment, err := h.service.CreateShipment(c.Context(), &req)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Shipment %d: %s", i+1, err.Error()))
			continue
		}
		shipments = append(shipments, shipment)
	}

	result := map[string]interface{}{
		"created": shipments,
		"errors":  errors,
		"summary": map[string]int{
			"total":      len(requests),
			"successful": len(shipments),
			"failed":     len(errors),
		},
	}

	return utils.SuccessResponse(c, result)
}

// HandleStatusWebhook обрабатывает webhook от BEX для обновления статусов
// @Summary Handle BEX status webhook
// @Description Receive status updates from BEX
// @Tags BEX Express
// @Accept json
// @Produce json
// @Param webhook body map[string]interface{} true "Webhook data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag
// @Router /api/v1/bex/webhook/status [post]
func (h *Handler) HandleStatusWebhook(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidWebhook")
	}

	// Process webhook data
	shipmentID, ok := data["shipment_id"].(float64)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.missingShipmentID")
	}

	// Update shipment status
	if _, err := h.service.GetShipmentStatus(c.Context(), int(shipmentID)); err != nil {
		// Log error but return success to BEX
		// We don't want them to retry
		h.logger.Error("Failed to update shipment status from webhook: %v", err)
	}

	return utils.SuccessResponse(c, nil)
}
