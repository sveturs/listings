package handler

import (
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	// Error messages
	errOrderNotFound = "order not found"
	errAccessDenied  = "access denied"
)

// CreateOrder создает новый заказ
// @Summary Create new order
// @Description Creates a new order from cart items or direct items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=models.StorefrontOrder} "Created order"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 409 {object} backend_pkg_utils.ErrorResponseSwag "Insufficient stock"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/orders [post]
func (h *OrdersHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse order request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_request_body")
	}

	// Валидация
	if req.StorefrontID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_storefront_id")
	}

	if len(req.Items) == 0 && req.CartID == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.no_items")
	}

	// Используем новый метод с транзакциями
	order, err := h.orderService.CreateOrderWithTx(c.Context(), h.db, &req, userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create order")

		switch err.Error() {
		case "insufficient stock":
			return utils.ErrorResponse(c, fiber.StatusConflict, "orders.error.insufficient_stock")
		case "storefront not found":
			return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.error.storefront_not_found")
		case "product not found":
			return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.error.product_not_found")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.create_failed")
		}
	}

	return utils.SuccessResponse(c, order)
}

// GetMyOrders возвращает заказы текущего пользователя
// @Summary Get user orders
// @Description Gets a list of orders for the authenticated user
// @Tags orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Param storefront_id query int false "Filter by storefront"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.StorefrontOrder} "User orders"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/orders [get]
func (h *OrdersHandler) GetMyOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	filter := &models.OrderFilter{
		CustomerID: &userID,
	}

	// Парсинг параметров запроса
	if page := c.QueryInt("page", 1); page > 0 {
		filter.Page = page
	}
	if limit := c.QueryInt("limit", 20); limit > 0 && limit <= 100 {
		filter.Limit = limit
	}
	if status := c.Query("status"); status != "" {
		orderStatus := models.OrderStatus(status)
		filter.Status = &orderStatus
	}
	if storefrontID := c.QueryInt("storefront_id"); storefrontID > 0 {
		filter.StorefrontID = &storefrontID
	}

	orders, total, err := h.orderService.GetOrders(c.Context(), filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user orders")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.get_failed")
	}

	response := map[string]interface{}{
		"orders": orders,
		"total":  total,
		"page":   filter.Page,
		"limit":  filter.Limit,
	}

	return utils.SuccessResponse(c, response)
}

// GetOrder возвращает детали заказа
// @Summary Get order details
// @Description Gets detailed information about a specific order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.StorefrontOrder} "Order details"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Order not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/orders/{id} [get]
func (h *OrdersHandler) GetOrder(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_order_id")
	}

	userID := c.Locals("user_id").(int)

	order, err := h.orderService.GetOrderByID(c.Context(), orderID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get order")

		switch err.Error() {
		case errOrderNotFound:
			return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.error.not_found")
		case errAccessDenied:
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.error.access_denied")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.get_failed")
		}
	}

	return utils.SuccessResponse(c, order)
}

// CancelOrder отменяет заказ
// @Summary Cancel order
// @Description Cancels an existing order if it's in a cancellable state
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Param cancel body models.CancelOrderRequest true "Cancellation reason"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.StorefrontOrder} "Canceled order"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Order not found"
// @Failure 409 {object} backend_pkg_utils.ErrorResponseSwag "Order cannot be canceled"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/orders/{id}/cancel [put]
func (h *OrdersHandler) CancelOrder(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_order_id")
	}

	var req models.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse cancel request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_request_body")
	}

	userID := c.Locals("user_id").(int)

	order, err := h.orderService.CancelOrder(c.Context(), orderID, userID, req.Reason)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to cancel order")

		switch err.Error() {
		case errOrderNotFound:
			return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.error.not_found")
		case errAccessDenied:
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.error.access_denied")
		case "order cannot be canceled":
			return utils.ErrorResponse(c, fiber.StatusConflict, "orders.error.cannot_cancel")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.cancel_failed")
		}
	}

	return utils.SuccessResponse(c, order)
}

// GetStorefrontOrders возвращает заказы витрины (для владельца)
// @Summary Get storefront orders
// @Description Gets a list of orders for a specific storefront (owner only)
// @Tags orders
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.StorefrontOrder} "Storefront orders"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/orders [get]
func (h *OrdersHandler) GetStorefrontOrders(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_storefront_id")
	}

	userID := c.Locals("user_id").(int)

	filter := &models.OrderFilter{
		StorefrontID: &storefrontID,
		SellerID:     &userID, // Только заказы витрин пользователя
	}

	// Парсинг параметров запроса
	if page := c.QueryInt("page", 1); page > 0 {
		filter.Page = page
	}
	if limit := c.QueryInt("limit", 20); limit > 0 && limit <= 100 {
		filter.Limit = limit
	}
	if status := c.Query("status"); status != "" {
		orderStatus := models.OrderStatus(status)
		filter.Status = &orderStatus
	}

	orders, total, err := h.orderService.GetOrders(c.Context(), filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get storefront orders")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.get_failed")
	}

	response := map[string]interface{}{
		"orders": orders,
		"total":  total,
		"page":   filter.Page,
		"limit":  filter.Limit,
	}

	return utils.SuccessResponse(c, response)
}

// UpdateOrderStatus обновляет статус заказа (для продавца)
// @Summary Update order status
// @Description Updates the status of an order (seller only)
// @Tags orders
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param order_id path int true "Order ID"
// @Param update body models.UpdateOrderStatusRequest true "Status update"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.StorefrontOrder} "Updated order"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Order not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/orders/{order_id}/status [put]
func (h *OrdersHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_storefront_id")
	}

	orderIDStr := c.Params("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_order_id")
	}

	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse status update request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_request_body")
	}

	userID := c.Locals("user_id").(int)

	order, err := h.orderService.UpdateOrderStatus(c.Context(), orderID, storefrontID, userID, req.Status, req.TrackingNumber, req.SellerNotes)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update order status")

		switch err.Error() {
		case errOrderNotFound:
			return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.error.not_found")
		case errAccessDenied:
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.error.access_denied")
		case "invalid status transition":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.error.invalid_status_transition")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.error.update_failed")
		}
	}

	return utils.SuccessResponse(c, order)
}
