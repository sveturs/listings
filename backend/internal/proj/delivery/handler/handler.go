package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/service"
	"backend/pkg/utils"
)

// Handler - HTTP обработчик для системы доставки
type Handler struct {
	service      *service.Service
	adminHandler *AdminHandler
	logger       zerolog.Logger
}

// NewHandler создает новый обработчик
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		service:      svc,
		adminHandler: NewAdminHandler(svc),
		logger:       log.Logger,
	}
}

// GetAdminHandler возвращает admin handler для консолидации
func (h *Handler) GetAdminHandler() *AdminHandler {
	return h.adminHandler
}

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Админские роуты регистрируются отдельно через module.go

	// Универсальные эндпоинты доставки
	delivery := router.Group("/delivery")
	delivery.Post("/calculate-universal", h.CalculateUniversal)
	delivery.Post("/calculate-cart", h.CalculateCart)
	delivery.Get("/providers", h.GetProviders)

	// Атрибуты доставки товаров
	products := router.Group("/products")
	products.Get("/:id/delivery-attributes", h.GetProductAttributes)
	products.Put("/:id/delivery-attributes", h.UpdateProductAttributes)

	// Дефолтные атрибуты категорий
	categories := router.Group("/categories")
	categories.Get("/:id/delivery-defaults", h.GetCategoryDefaults)
	categories.Put("/:id/delivery-defaults", h.UpdateCategoryDefaults)
	categories.Post("/:id/apply-defaults", h.ApplyCategoryDefaults)

	// Отправления
	shipments := router.Group("/shipments")
	shipments.Post("/", h.CreateShipment)
	shipments.Get("/:id", h.GetShipment)
	shipments.Delete("/:id", h.CancelShipment)
	shipments.Get("/track/:tracking", h.TrackShipment)

	// Тестовые эндпоинты регистрируются отдельно через RegisterTestRoutes (см. module.go)

	// Административные роуты
	admin := router.Group("/admin/delivery")
	admin.Get("/providers", h.GetProvidersAdmin)
	admin.Put("/providers/:id", h.UpdateProvider)
	admin.Post("/pricing-rules", h.CreatePricingRule)
	admin.Get("/analytics", h.GetAnalytics)
}

// RegisterWebhookRoutes регистрирует webhook маршруты (без авторизации)
func (h *Handler) RegisterWebhookRoutes(router fiber.Router) {
	router.Post("/:provider/tracking", h.HandleTrackingWebhook)
}

// RegisterTestRoutes регистрирует тестовые маршруты (без авторизации для удобства тестирования)
func (h *Handler) RegisterTestRoutes(app fiber.Router) {
	// Тестовые эндпоинты (новые, через gRPC микросервис)
	// Используем /api/public/* чтобы НЕ наследовать RequireAuth middleware от /api/v1/*
	test := app.Group("/api/public/delivery/test")
	test.Post("/shipment", h.CreateTestShipment)
	test.Get("/tracking/:tracking_number", h.TrackTestShipment)
	test.Post("/cancel/:id", h.CancelTestShipment)
	test.Post("/calculate", h.CalculateTestRate)
	test.Get("/settlements", h.GetTestSettlements)
	test.Get("/streets", h.GetTestStreets)
	test.Get("/parcel-lockers", h.GetTestParcelLockers)
	test.Get("/delivery-services", h.GetTestDeliveryServices)
	test.Post("/validate-address", h.ValidateTestAddress)
	test.Get("/providers", h.GetTestProviders)
	test.Get("/config", h.GetTestConfig)
	test.Get("/history", h.GetTestHistory)
}

// CalculateUniversal - DEPRECATED: расчет перенесен в delivery microservice
// @Summary Calculate universal delivery rates (DEPRECATED)
// @Description DEPRECATED: This endpoint has been moved to delivery microservice. Use gRPC CalculateRate instead.
// @Tags delivery
// @Accept json
// @Produce json
// @Param request body CalculationRequest true "Calculation request"
// @Success 501 {object} utils.ErrorResponseSwag "Not implemented"
// @Router /api/v1/delivery/calculate-universal [post]
func (h *Handler) CalculateUniversal(c *fiber.Ctx) error {
	return utils.SendErrorResponse(
		c,
		fiber.StatusNotImplemented,
		"delivery.calculation_moved_to_microservice",
		fiber.Map{
			"message": "Calculation functionality has been moved to delivery microservice. Use gRPC CalculateRate method instead.",
		},
	)
}

// CalculateCart - DEPRECATED: расчет перенесен в delivery microservice
// @Summary Calculate delivery for cart (DEPRECATED)
// @Description DEPRECATED: This endpoint has been moved to delivery microservice. Use gRPC CalculateRate instead.
// @Tags delivery
// @Accept json
// @Produce json
// @Param request body CartCalculationRequest true "Cart calculation request"
// @Success 501 {object} utils.ErrorResponseSwag "Not implemented"
// @Router /api/v1/delivery/calculate-cart [post]
func (h *Handler) CalculateCart(c *fiber.Ctx) error {
	return utils.SendErrorResponse(
		c,
		fiber.StatusNotImplemented,
		"delivery.calculation_moved_to_microservice",
		fiber.Map{
			"message": "Cart calculation functionality has been moved to delivery microservice. Use gRPC CalculateRate method instead.",
		},
	)
}

// GetProviders - получает список доступных провайдеров
// @Summary Get delivery providers
// @Description Get list of available delivery providers
// @Tags delivery
// @Produce json
// @Param active query bool false "Only active providers"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.Provider} "List of providers"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/delivery/providers [get]
func (h *Handler) GetProviders(c *fiber.Ctx) error {
	activeOnly := c.QueryBool("active", true)

	providers, err := h.service.GetProviders(c.Context(), activeOnly)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_providers", nil)
	}

	return utils.SendSuccessResponse(c, providers, "Список провайдеров получен")
}

// GetProductAttributes - получает атрибуты доставки товара
// @Summary Get product delivery attributes
// @Description Get delivery attributes for a product
// @Tags delivery
// @Produce json
// @Param id path int true "Product ID"
// @Param type query string false "Product type (listing or storefront_product)" default(listing)
// @Success 200 {object} utils.SuccessResponseSwag{data=models.DeliveryAttributes} "Product attributes"
// @Failure 404 {object} utils.ErrorResponseSwag "Product not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/products/{id}/delivery-attributes [get]
func (h *Handler) GetProductAttributes(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_product_id", nil)
	}

	productType := c.Query("type", "listing")

	attrs, err := h.service.GetProductAttributes(c.Context(), productID, productType)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.product_not_found", nil)
	}

	return utils.SendSuccessResponse(c, attrs, "Атрибуты доставки получены")
}

// UpdateProductAttributes - обновляет атрибуты доставки товара
// @Summary Update product delivery attributes
// @Description Update delivery attributes for a product
// @Tags delivery
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param type query string false "Product type (listing or storefront_product)" default(listing)
// @Param attributes body models.DeliveryAttributes true "Delivery attributes"
// @Success 200 {object} utils.SuccessResponseSwag "Success"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} utils.ErrorResponseSwag "Product not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/products/{id}/delivery-attributes [put]
func (h *Handler) UpdateProductAttributes(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_product_id", nil)
	}

	productType := c.Query("type", "listing")

	var attrs models.DeliveryAttributes
	if err := c.BodyParser(&attrs); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
	}

	if err := h.service.UpdateProductAttributes(c.Context(), productID, productType, &attrs); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update", nil)
	}

	return utils.SendSuccessResponse(c, nil, "Атрибуты доставки обновлены")
}

// GetCategoryDefaults - получает дефолтные атрибуты категории
// @Summary Get category default attributes
// @Description Get default delivery attributes for a category
// @Tags delivery
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.CategoryDefaults} "Category defaults"
// @Failure 404 {object} utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/categories/{id}/delivery-defaults [get]
func (h *Handler) GetCategoryDefaults(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_category_id", nil)
	}

	defaults, err := h.service.GetCategoryDefaults(c.Context(), categoryID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.category_not_found", nil)
	}

	return utils.SendSuccessResponse(c, defaults, "Дефолтные атрибуты получены")
}

// UpdateCategoryDefaults - обновляет дефолтные атрибуты категории
// @Summary Update category default attributes
// @Description Update default delivery attributes for a category
// @Tags delivery
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param defaults body models.CategoryDefaults true "Category defaults"
// @Success 200 {object} utils.SuccessResponseSwag "Success"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/categories/{id}/delivery-defaults [put]
func (h *Handler) UpdateCategoryDefaults(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_category_id", nil)
	}

	var defaults models.CategoryDefaults
	if err := c.BodyParser(&defaults); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
	}

	defaults.CategoryID = categoryID

	if err := h.service.UpdateCategoryDefaults(c.Context(), &defaults); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update", nil)
	}

	return utils.SendSuccessResponse(c, nil, "Дефолтные атрибуты обновлены")
}

// ApplyCategoryDefaults - применяет дефолтные атрибуты к товарам категории
// @Summary Apply category defaults to products
// @Description Apply default delivery attributes to all products in category without attributes
// @Tags delivery
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]int} "Number of updated products"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/categories/{id}/apply-defaults [post]
func (h *Handler) ApplyCategoryDefaults(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_category_id", nil)
	}

	count, err := h.service.ApplyCategoryDefaults(c.Context(), categoryID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_apply", nil)
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"updated_count": count,
	}, "Дефолтные атрибуты применены")
}

// CreateShipment - создает отправление
// @Summary Create shipment
// @Description Create a new shipment with selected provider
// @Tags delivery
// @Accept json
// @Produce json
// @Param shipment body service.CreateShipmentRequest true "Shipment request"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Shipment} "Created shipment"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/shipments [post]
func (h *Handler) CreateShipment(c *fiber.Ctx) error {
	var req service.CreateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
	}

	shipment, err := h.service.CreateShipment(c.Context(), &req)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_shipment", nil)
	}

	return utils.SendSuccessResponse(c, shipment, "Отправление создано")
}

// GetShipment - получает информацию об отправлении
// @Summary Get shipment
// @Description Get shipment information by ID
// @Tags delivery
// @Produce json
// @Param id path int true "Shipment ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Shipment} "Shipment info"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/shipments/{id} [get]
func (h *Handler) GetShipment(c *fiber.Ctx) error {
	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_shipment_id", nil)
	}

	shipment, err := h.service.GetShipment(c.Context(), shipmentID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", nil)
	}

	return utils.SendSuccessResponse(c, shipment, "Информация об отправлении")
}

// CancelShipment - отменяет отправление
// @Summary Cancel shipment
// @Description Cancel an existing shipment
// @Tags delivery
// @Produce json
// @Param id path int true "Shipment ID"
// @Param reason body CancelRequest false "Cancel reason"
// @Success 200 {object} utils.SuccessResponseSwag "Success"
// @Failure 400 {object} utils.ErrorResponseSwag "Cannot cancel"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/shipments/{id} [delete]
func (h *Handler) CancelShipment(c *fiber.Ctx) error {
	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_shipment_id", nil)
	}

	var req CancelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request_body", nil)
	}

	if err := h.service.CancelShipment(c.Context(), shipmentID, req.Reason); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_cancel", nil)
	}

	return utils.SendSuccessResponse(c, nil, "Отправление отменено")
}

// TrackShipment - отслеживает отправление
// @Summary Track shipment
// @Description Track shipment by tracking number
// @Tags delivery
// @Produce json
// @Param tracking path string true "Tracking number"
// @Success 200 {object} utils.SuccessResponseSwag{data=service.TrackingInfo} "Tracking info"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/shipments/track/{tracking} [get]
func (h *Handler) TrackShipment(c *fiber.Ctx) error {
	trackingNumber := c.Params("tracking")
	if trackingNumber == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_tracking_number", nil)
	}

	info, err := h.service.TrackShipment(c.Context(), trackingNumber)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", nil)
	}

	return utils.SendSuccessResponse(c, info, "Информация об отслеживании получена")
}

// GetProvidersAdmin - получает список провайдеров для админки
// Note: Swagger documentation is in AdminHandler.GetProviders
func (h *Handler) GetProvidersAdmin(c *fiber.Ctx) error {
	providers, err := h.service.GetProviders(c.Context(), false)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_providers", nil)
	}

	return utils.SendSuccessResponse(c, providers, "Список провайдеров получен")
}

// UpdateProvider - обновляет провайдера
// Note: Swagger documentation is in AdminHandler.UpdateProviderConfig
func (h *Handler) UpdateProvider(c *fiber.Ctx) error {
	providerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_provider_id", nil)
	}

	var provider models.Provider
	if err := c.BodyParser(&provider); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
	}

	if err := h.service.UpdateProvider(c.Context(), providerID, &provider); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update_provider", nil)
	}

	return utils.SendSuccessResponse(c, nil, "Провайдер обновлен")
}

// CreatePricingRule - создает правило расчета стоимости
// Note: Swagger documentation is in AdminHandler.CreatePricingRule
func (h *Handler) CreatePricingRule(c *fiber.Ctx) error {
	var rule models.PricingRule
	if err := c.BodyParser(&rule); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", nil)
	}

	createdRule, err := h.service.CreatePricingRule(c.Context(), &rule)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_rule", nil)
	}

	return utils.SendSuccessResponse(c, createdRule, "Правило ценообразования создано")
}

// GetAnalytics - получает аналитику по доставкам
// Note: Swagger documentation is in AdminHandler.GetAnalytics
func (h *Handler) GetAnalytics(c *fiber.Ctx) error {
	// Парсим даты
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_from_date", nil)
		}
	} else {
		// По умолчанию - последние 30 дней
		from = time.Now().AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_to_date", nil)
		}
	} else {
		// По умолчанию - сегодня
		to = time.Now()
	}

	// Парсим ID провайдера
	var providerID *int
	if providerIDStr := c.Query("provider_id"); providerIDStr != "" {
		pID, err := strconv.Atoi(providerIDStr)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_provider_id", nil)
		}
		providerID = &pID
	}

	analytics, err := h.service.GetDeliveryAnalytics(c.Context(), from, to, providerID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_analytics", nil)
	}

	return utils.SendSuccessResponse(c, analytics, "Аналитика доставок")
}

// HandleTrackingWebhook - DEPRECATED: webhook обработка перенесена в delivery microservice
// @Summary Handle tracking webhook (DEPRECATED)
// @Description DEPRECATED: Webhook processing moved to delivery microservice. Configure webhooks to point directly to delivery service.
// @Tags delivery
// @Accept json
// @Produce json
// @Param provider path string true "Provider code (post_express, bex_express, etc.)"
// @Param payload body object true "Webhook payload from provider"
// @Success 501 {object} utils.ErrorResponseSwag "Moved to microservice"
// @Router /api/v1/delivery/webhooks/{provider}/tracking [post]
func (h *Handler) HandleTrackingWebhook(c *fiber.Ctx) error {
	return utils.SendErrorResponse(c, fiber.StatusNotImplemented, "delivery.webhook.moved_to_microservice", map[string]interface{}{
		"message": "Webhook processing has been moved to delivery microservice",
		"note":    "Configure webhooks to point directly to the delivery service",
	})
}

// Структуры запросов

// Location представляет локацию для расчета доставки
type Location struct {
	City       string  `json:"city"`
	PostalCode string  `json:"postal_code,omitempty"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// ItemDimensions представляет размеры товара
type ItemDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ItemWithAttrs представляет товар с атрибутами для расчета доставки
type ItemWithAttrs struct {
	ProductID   int            `json:"product_id"`
	Quantity    int            `json:"quantity"`
	Weight      float64        `json:"weight"`
	Dimensions  ItemDimensions `json:"dimensions"`
	Value       float64        `json:"value"`
	Fragile     bool           `json:"fragile,omitempty"`
	Refrigerate bool           `json:"refrigerate,omitempty"`
	Category    string         `json:"category,omitempty"`
}

// CancelRequest - запрос отмены
type CancelRequest struct {
	Reason string `json:"reason"`
}
