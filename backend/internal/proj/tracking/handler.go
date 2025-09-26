package tracking

import (
	"strconv"

	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// TrackingHandler обрабатывает запросы связанные с трекингом доставок
type TrackingHandler struct {
	deliveryService *DeliveryService
	courierService  *CourierService
	hub             *Hub
}

// NewTrackingHandler создаёт новый обработчик трекинга
func NewTrackingHandler(
	deliveryService *DeliveryService,
	courierService *CourierService,
	hub *Hub,
) *TrackingHandler {
	return &TrackingHandler{
		deliveryService: deliveryService,
		courierService:  courierService,
		hub:             hub,
	}
}

// GetDeliveryByToken получает информацию о доставке по токену
// @Summary Get delivery information by tracking token
// @Description Retrieves delivery details using tracking token
// @Tags tracking
// @Accept json
// @Produce json
// @Param token path string true "Tracking token"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=tracking.Delivery} "Delivery information"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid token"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Delivery not found"
// @Router /api/v1/tracking/{token} [get]
func (h *TrackingHandler) GetDeliveryByToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "tracking.error.invalidToken")
	}

	// For now, return simplified data for Post Express shipments
	shipment, err := h.deliveryService.GetPostExpressShipment(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "tracking.error.notFound")
	}

	return utils.SuccessResponse(c, shipment)
}

// UpdateCourierLocation обновляет местоположение курьера
// @Summary Update courier location
// @Description Updates courier's current location and broadcasts to tracking clients
// @Tags tracking
// @Accept json
// @Produce json
// @Param courier_id path int true "Courier ID"
// @Param location body CourierLocationUpdate true "Location data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Location updated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Courier not found"
// @Router /api/v1/courier/{courier_id}/location [post]
func (h *TrackingHandler) UpdateCourierLocation(c *fiber.Ctx) error {
	courierIDStr := c.Params("courier_id")
	courierID, err := strconv.Atoi(courierIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "courier.error.invalidID")
	}

	var req CourierLocationUpdate
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	// Проверяем существование курьера
	_, err = h.courierService.GetCourier(c.Context(), courierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "courier.error.notFound")
	}

	// Обновляем локацию курьера
	err = h.courierService.UpdateCourierLocation(c.Context(), courierID, req.Latitude, req.Longitude)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "courier.error.updateFailed")
	}

	// Получаем активные доставки курьера
	deliveries, err := h.deliveryService.GetActiveDeliveries(courierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "delivery.error.getActiveFailed")
	}

	// Broadcast обновление локации для каждой доставки
	for _, delivery := range deliveries {
		h.hub.BroadcastLocationUpdate(
			delivery.ID,
			req.Latitude,
			req.Longitude,
			req.Speed,
			req.Heading,
		)
	}

	return utils.SuccessResponse(c, nil)
}

// GetActiveDeliveries получает активные доставки курьера
// @Summary Get active deliveries for courier
// @Description Retrieves list of active deliveries assigned to a courier
// @Tags courier
// @Accept json
// @Produce json
// @Param courier_id path int true "Courier ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]tracking.Delivery} "Active deliveries"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid courier ID"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Courier not found"
// @Router /api/v1/courier/{courier_id}/deliveries [get]
func (h *TrackingHandler) GetActiveDeliveries(c *fiber.Ctx) error {
	courierIDStr := c.Params("courier_id")
	courierID, err := strconv.Atoi(courierIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "courier.error.invalidID")
	}

	deliveries, err := h.deliveryService.GetActiveDeliveries(courierID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "delivery.error.getActiveFailed")
	}

	return utils.SuccessResponse(c, deliveries)
}

// UpdateDeliveryStatus обновляет статус доставки
// @Summary Update delivery status
// @Description Updates the status of a delivery (picked_up, in_transit, delivered, failed)
// @Tags tracking
// @Accept json
// @Produce json
// @Param delivery_id path int true "Delivery ID"
// @Param status body DeliveryStatusUpdate true "Status update"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Status updated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Delivery not found"
// @Router /api/v1/delivery/{delivery_id}/status [put]
func (h *TrackingHandler) UpdateDeliveryStatus(c *fiber.Ctx) error {
	deliveryIDStr := c.Params("delivery_id")
	deliveryID, err := strconv.Atoi(deliveryIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "delivery.error.invalidID")
	}

	var req DeliveryStatusUpdate
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	// Валидация статуса
	validStatuses := []string{"pending", "picked_up", "in_transit", "delivered", "failed"}
	isValid := false
	for _, validStatus := range validStatuses {
		if req.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "delivery.error.invalidStatus")
	}

	// Обновляем статус
	err = h.deliveryService.UpdateDeliveryStatus(c.Context(), deliveryID, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "delivery.error.updateFailed")
	}

	return utils.SuccessResponse(c, nil)
}

// GetWebSocketConnections получает статистику WebSocket соединений
// @Summary Get WebSocket connections statistics
// @Description Returns the number of active WebSocket connections per delivery
// @Tags tracking
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=map[string]int} "Active connections"
// @Router /api/v1/tracking/connections [get]
func (h *TrackingHandler) GetWebSocketConnections(c *fiber.Ctx) error {
	stats := h.hub.GetActiveConnections()
	return utils.SuccessResponse(c, stats)
}

// CreateDelivery создает новую доставку
// @Summary Create new delivery
// @Description Creates a new delivery with courier assignment
// @Tags delivery
// @Accept json
// @Produce json
// @Param delivery body CreateDeliveryRequest true "Delivery data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=tracking.Delivery} "Delivery created"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Router /api/v1/delivery [post]
func (h *TrackingHandler) CreateDelivery(c *fiber.Ctx) error {
	var req CreateDeliveryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	// Создаем доставку
	delivery, err := h.deliveryService.CreateDelivery(
		c.Context(),
		req.OrderID,
		req.CourierID,
		req.PickupAddress,
		req.DeliveryAddress,
		req.PickupLatitude,
		req.PickupLongitude,
		req.DeliveryLatitude,
		req.DeliveryLongitude,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "delivery.error.createFailed")
	}

	return utils.SuccessResponse(c, delivery)
}

// DTOs для запросов

// CourierLocationUpdate запрос на обновление локации курьера
type CourierLocationUpdate struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Speed     float64 `json:"speed" validate:"min=0"`
	Heading   int     `json:"heading" validate:"min=0,max=359"`
}

// DeliveryStatusUpdate запрос на обновление статуса доставки
type DeliveryStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=pending picked_up in_transit delivered failed"`
}

// CreateDeliveryRequest запрос на создание доставки
type CreateDeliveryRequest struct {
	OrderID           int     `json:"order_id" validate:"required,min=1"`
	CourierID         int     `json:"courier_id" validate:"required,min=1"`
	PickupAddress     string  `json:"pickup_address" validate:"required,min=5"`
	DeliveryAddress   string  `json:"delivery_address" validate:"required,min=5"`
	PickupLatitude    float64 `json:"pickup_latitude" validate:"required,min=-90,max=90"`
	PickupLongitude   float64 `json:"pickup_longitude" validate:"required,min=-180,max=180"`
	DeliveryLatitude  float64 `json:"delivery_latitude" validate:"required,min=-90,max=90"`
	DeliveryLongitude float64 `json:"delivery_longitude" validate:"required,min=-180,max=180"`
}
