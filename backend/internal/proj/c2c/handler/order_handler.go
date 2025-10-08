package handler

import (
	"log"
	"strconv"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/proj/c2c/service"
	"backend/pkg/utils"
)

// OrderHandler обрабатывает запросы связанные с заказами
type OrderHandler struct {
	orderService service.OrderServiceInterface
}

// NewOrderHandler создает новый обработчик заказов
func NewOrderHandler(orderService service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// RegisterRoutes регистрирует маршруты для заказов
func (h *OrderHandler) RegisterRoutes(app fiber.Router) {
	// Маршруты уже находятся под /api/v1/marketplace/orders благодаря группе в handler.go

	// Создание заказа
	app.Post("/create", h.CreateMarketplaceOrder)

	// Списки заказов
	app.Get("/my/purchases", h.GetMyPurchases)
	app.Get("/my/sales", h.GetMySales)

	// Операции с конкретным заказом
	app.Get("/:id", h.GetOrderDetails)
	app.Post("/:id/confirm-payment", h.ConfirmPayment)
	app.Post("/:id/ship", h.MarkAsShipped)
	app.Post("/:id/confirm-delivery", h.ConfirmDelivery)
	app.Post("/:id/dispute", h.OpenDispute)
	app.Post("/:id/message", h.AddMessage)
}

// CreateMarketplaceOrderRequest запрос на создание заказа
type CreateMarketplaceOrderRequest struct {
	ListingID     int64   `json:"listing_id" validate:"required"`
	Message       *string `json:"message,omitempty"`
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=card bank_transfer"`
}

// CreateMarketplaceOrder создает заказ для листинга
// @Summary Create marketplace order
// @Description Creates a new order for marketplace listing with payment preauthorization
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CreateMarketplaceOrderRequest true "Order details"
// @Success 200 {object} utils.SuccessResponseSwag "Order created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/marketplace/orders/create [post]
func (h *OrderHandler) CreateMarketplaceOrder(c *fiber.Ctx) error {
	log.Printf("CreateMarketplaceOrder called")

	userID, _ := authMiddleware.GetUserID(c)
	log.Printf("UserID: %d", userID)

	var req CreateMarketplaceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidRequest")
	}

	log.Printf("Request: ListingID=%d, PaymentMethod=%s", req.ListingID, req.PaymentMethod)

	// TODO: Добавить валидацию структуры
	// if err := utils.ValidateStruct(&req); err != nil {
	//     return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidData")
	// }

	// Создаем заказ
	order, payment, err := h.orderService.CreateOrderFromRequest(c.Context(), service.CreateOrderRequest{
		BuyerID:       int64(userID),
		ListingID:     req.ListingID,
		Message:       req.Message,
		PaymentMethod: req.PaymentMethod,
		ReturnURL:     c.Get("Referer", "/"), // Используем Referer или дефолт
	})
	if err != nil {
		// Логируем ошибку для отладки
		log.Printf("Error creating order: %v", err)

		// Обрабатываем специфичные ошибки
		switch err.Error() {
		case "listing is not active":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.listingNotActive")
		case "cannot buy own listing":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.cannotBuyOwnListing")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.createError")
		}
	}

	// Возвращаем результат
	return utils.SuccessResponse(c, fiber.Map{
		"order_id":    order.ID,
		"payment_url": payment.PaymentURL,
		"message":     "orders.created",
	})
}

// GetMyPurchases получает заказы где пользователь - покупатель
// @Summary Get my purchases
// @Description Get list of orders where current user is buyer
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponseSwag "Orders list"
// @Router /api/v1/marketplace/orders/my/purchases [get]
func (h *OrderHandler) GetMyPurchases(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	orders, total, err := h.orderService.GetBuyerOrders(c.Context(), int64(userID), page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.fetchError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"orders": convertOrdersToResponse(orders),
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetMySales получает заказы где пользователь - продавец
// @Summary Get my sales
// @Description Get list of orders where current user is seller
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponseSwag "Orders list"
// @Router /api/v1/marketplace/orders/my/sales [get]
func (h *OrderHandler) GetMySales(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	orders, total, err := h.orderService.GetSellerOrders(c.Context(), int64(userID), page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.fetchError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"orders": convertOrdersToResponse(orders),
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetOrderDetails получает детали заказа
// @Summary Get order details
// @Description Get detailed information about specific order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceOrder} "Order details"
// @Router /api/v1/marketplace/orders/{id} [get]
func (h *OrderHandler) GetOrderDetails(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	}

	order, err := h.orderService.GetOrderDetails(c.Context(), orderID, int64(userID))
	if err != nil {
		if err.Error() == "unauthorized: not a party of this order" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.accessDenied")
		}
		return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.notFound")
	}

	return utils.SuccessResponse(c, order)
}

// ConfirmPaymentRequest запрос на подтверждение оплаты
type ConfirmPaymentRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

// ConfirmPayment подтверждает оплату заказа
// @Summary Confirm order payment
// @Description Confirm order payment (for mock payment provider)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body ConfirmPaymentRequest true "Payment confirmation"
// @Success 200 {object} utils.SuccessResponseSwag "Payment confirmed"
// @Router /api/v1/marketplace/orders/{id}/confirm-payment [post]
func (h *OrderHandler) ConfirmPayment(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	}

	var req ConfirmPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidRequest")
	}

	// Проверяем что пользователь - покупатель этого заказа
	order, err := h.orderService.GetOrderDetails(c.Context(), orderID, int64(userID))
	if err != nil {
		if err.Error() == "unauthorized: not a party of this order" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.accessDenied")
		}
		return utils.ErrorResponse(c, fiber.StatusNotFound, "orders.notFound")
	}

	if order.BuyerID != int64(userID) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.notBuyer")
	}

	// Подтверждаем оплату
	err = h.orderService.ConfirmPayment(c.Context(), orderID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.confirmPaymentError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "orders.paymentConfirmed"})
}

// MarkAsShippedRequest запрос на отметку отправки
type MarkAsShippedRequest struct {
	ShippingMethod string `json:"shipping_method" validate:"required"`
	TrackingNumber string `json:"tracking_number" validate:"required"`
}

// MarkAsShipped отмечает заказ как отправленный
// @Summary Mark order as shipped
// @Description Mark order as shipped by seller
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body MarkAsShippedRequest true "Shipping details"
// @Success 200 {object} utils.SuccessResponseSwag "Order marked as shipped"
// @Router /api/v1/marketplace/orders/{id}/ship [post]
func (h *OrderHandler) MarkAsShipped(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	}

	var req MarkAsShippedRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidRequest")
	}

	// TODO: Добавить валидацию структуры
	// if err := utils.ValidateStruct(&req); err != nil {
	//     return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidData")
	// }

	err = h.orderService.MarkAsShipped(c.Context(), orderID, int64(userID), req.ShippingMethod, req.TrackingNumber)
	if err != nil {
		switch err.Error() {
		case "unauthorized: not the seller":
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.notSeller")
		case "order must be paid before shipping":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.notPaid")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.updateError")
		}
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "orders.shipped"})
}

// ConfirmDelivery подтверждает получение заказа
// @Summary Confirm order delivery
// @Description Confirm order delivery by buyer
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.SuccessResponseSwag "Delivery confirmed"
// @Router /api/v1/marketplace/orders/{id}/confirm-delivery [post]
func (h *OrderHandler) ConfirmDelivery(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	}

	err = h.orderService.ConfirmDelivery(c.Context(), orderID, int64(userID))
	if err != nil {
		switch err.Error() {
		case "unauthorized: not the buyer":
			return utils.ErrorResponse(c, fiber.StatusForbidden, "orders.notBuyer")
		case "order must be shipped before delivery confirmation":
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.notShipped")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "orders.updateError")
		}
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "orders.delivered"})
}

// OpenDisputeRequest запрос на открытие спора
type OpenDisputeRequest struct {
	Reason string `json:"reason" validate:"required,min=10"`
}

// OpenDispute открывает спор по заказу
// @Summary Open order dispute
// @Description Open dispute for order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body OpenDisputeRequest true "Dispute reason"
// @Success 200 {object} utils.SuccessResponseSwag "Dispute opened"
// @Router /api/v1/marketplace/orders/{id}/dispute [post]
func (h *OrderHandler) OpenDispute(c *fiber.Ctx) error {
	userID, _ := authMiddleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	}

	var req OpenDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidRequest")
	}

	// TODO: Добавить валидацию структуры
	// if err := utils.ValidateStruct(&req); err != nil {
	//     return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidData")
	// }

	err = h.orderService.OpenDispute(c.Context(), orderID, int64(userID), req.Reason)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.disputeError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "orders.disputeOpened"})
}

// AddMessageRequest запрос на добавление сообщения
type AddMessageRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

// AddMessage добавляет сообщение к заказу
// @Summary Add message to order
// @Description Add message to order conversation
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body AddMessageRequest true "Message content"
// @Success 200 {object} utils.SuccessResponseSwag "Message added"
// @Router /api/v1/marketplace/orders/{id}/message [post]
func (h *OrderHandler) AddMessage(c *fiber.Ctx) error {
	// userID, _ := authMiddleware.GetUserID(c)
	// orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	// if err != nil {
	//     return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidID")
	// }

	// var req AddMessageRequest
	// if err := c.BodyParser(&req); err != nil {
	//     return utils.ErrorResponse(c, fiber.StatusBadRequest, "orders.invalidRequest")
	// }

	// TODO: Проверить доступ к заказу и добавить сообщение

	return utils.SuccessResponse(c, fiber.Map{"message": "orders.messageAdded"})
}

// convertOrdersToResponse преобразует заказы в формат для ответа
func convertOrdersToResponse(orders []*models.MarketplaceOrder) []interface{} {
	result := make([]interface{}, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}
	return result
}
