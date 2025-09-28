package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/calculator"
	"backend/internal/proj/delivery/factory"
	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/service"
	"backend/pkg/utils"
)

// Handler - HTTP обработчик для системы доставки
type Handler struct {
	service      *service.Service
	adminHandler *AdminHandler
}

// NewHandler создает новый обработчик
func NewHandler(db *sqlx.DB, providerFactory *factory.ProviderFactory) *Handler {
	svc := service.NewService(db, providerFactory)
	return &Handler{
		service:      svc,
		adminHandler: NewAdminHandler(svc),
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

// CalculateUniversal - универсальный расчет стоимости доставки для всех провайдеров
// @Summary Calculate universal delivery rates
// @Description Calculate delivery rates across all available providers
// @Tags delivery
// @Accept json
// @Produce json
// @Param request body calculator.CalculationRequest true "Calculation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=calculator.CalculationResponse} "Calculation results"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/delivery/calculate-universal [post]
func (h *Handler) CalculateUniversal(c *fiber.Ctx) error {
	var req calculator.CalculationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// Используем мок калькулятор для демонстрации
	mockCalc := calculator.NewMockCalculator()
	resp, err := mockCalc.CalculateMock(c.Context(), &req)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.calculation_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, resp, "Стоимость доставки рассчитана")
}

// CalculateCart - расчет стоимости доставки для корзины
// @Summary Calculate delivery for cart
// @Description Calculate delivery rates for cart items with optimization
// @Tags delivery
// @Accept json
// @Produce json
// @Param request body CartCalculationRequest true "Cart calculation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=calculator.CalculationResponse} "Calculation results"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/delivery/calculate-cart [post]
func (h *Handler) CalculateCart(c *fiber.Ctx) error {
	var req CartCalculationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// Конвертируем в стандартный запрос расчета
	calcReq := &calculator.CalculationRequest{
		FromLocation:   req.FromLocation,
		ToLocation:     req.ToLocation,
		Items:          req.Items,
		InsuranceValue: req.InsuranceValue,
		CODAmount:      req.CODAmount,
		DeliveryType:   req.DeliveryType,
	}

	resp, err := h.service.CalculateDelivery(c.Context(), calcReq)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.calculation_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, resp, "Стоимость доставки для корзины рассчитана")
}

// GetProviders - получает список доступных провайдеров
// @Summary Get delivery providers
// @Description Get list of available delivery providers
// @Tags delivery
// @Produce json
// @Param active query bool false "Only active providers"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_proj_delivery_models.Provider} "List of providers"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/delivery/providers [get]
func (h *Handler) GetProviders(c *fiber.Ctx) error {
	activeOnly := c.QueryBool("active", true)

	providers, err := h.service.GetProviders(c.Context(), activeOnly)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_providers", fiber.Map{
			"error": err.Error(),
		})
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
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_delivery_models.DeliveryAttributes} "Product attributes"
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
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.product_not_found", fiber.Map{
			"error": err.Error(),
		})
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
// @Param attributes body backend_internal_proj_delivery_models.DeliveryAttributes true "Delivery attributes"
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
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.UpdateProductAttributes(c.Context(), productID, productType, &attrs); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, nil, "Атрибуты доставки обновлены")
}

// GetCategoryDefaults - получает дефолтные атрибуты категории
// @Summary Get category default attributes
// @Description Get default delivery attributes for a category
// @Tags delivery
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_delivery_models.CategoryDefaults} "Category defaults"
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
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.category_not_found", fiber.Map{
			"error": err.Error(),
		})
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
// @Param defaults body backend_internal_proj_delivery_models.CategoryDefaults true "Category defaults"
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
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	defaults.CategoryID = categoryID

	if err := h.service.UpdateCategoryDefaults(c.Context(), &defaults); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update", fiber.Map{
			"error": err.Error(),
		})
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
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_apply", fiber.Map{
			"error": err.Error(),
		})
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
// @Param shipment body backend_internal_proj_delivery_service.CreateShipmentRequest true "Shipment request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_delivery_models.Shipment} "Created shipment"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/shipments [post]
func (h *Handler) CreateShipment(c *fiber.Ctx) error {
	var req service.CreateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	shipment, err := h.service.CreateShipment(c.Context(), &req)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_shipment", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, shipment, "Отправление создано")
}

// GetShipment - получает информацию об отправлении
// @Summary Get shipment
// @Description Get shipment information by ID
// @Tags delivery
// @Produce json
// @Param id path int true "Shipment ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_delivery_models.Shipment} "Shipment info"
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
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", fiber.Map{
			"error": err.Error(),
		})
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
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_cancel", fiber.Map{
			"error": err.Error(),
		})
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
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "error.shipment_not_found", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, info, "Информация об отслеживании получена")
}

// GetProvidersAdmin - получает список провайдеров для админки
// @Summary Get providers for admin
// @Description Get detailed list of all providers for admin panel
// @Tags admin-delivery
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_proj_delivery_models.Provider} "List of providers"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/admin/delivery/providers [get]
func (h *Handler) GetProvidersAdmin(c *fiber.Ctx) error {
	providers, err := h.service.GetProviders(c.Context(), false)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_providers", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, providers, "Список провайдеров получен")
}

// UpdateProvider - обновляет провайдера
// @Summary Update provider
// @Description Update provider settings
// @Tags admin-delivery
// @Accept json
// @Produce json
// @Param id path int true "Provider ID"
// @Param provider body backend_internal_proj_delivery_models.Provider true "Provider data"
// @Success 200 {object} utils.SuccessResponseSwag "Success"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/admin/delivery/providers/{id} [put]
func (h *Handler) UpdateProvider(c *fiber.Ctx) error {
	providerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_provider_id", nil)
	}

	var provider models.Provider
	if err := c.BodyParser(&provider); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.UpdateProvider(c.Context(), providerID, &provider); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_update_provider", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, nil, "Провайдер обновлен")
}

// CreatePricingRule - создает правило расчета стоимости
// @Summary Create pricing rule
// @Description Create new pricing rule for provider
// @Tags admin-delivery
// @Accept json
// @Produce json
// @Param rule body backend_internal_proj_delivery_models.PricingRule true "Pricing rule"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_delivery_models.PricingRule} "Created rule"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/admin/delivery/pricing-rules [post]
func (h *Handler) CreatePricingRule(c *fiber.Ctx) error {
	var rule models.PricingRule
	if err := c.BodyParser(&rule); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	createdRule, err := h.service.CreatePricingRule(c.Context(), &rule)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_create_rule", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, createdRule, "Правило ценообразования создано")
}

// GetAnalytics - получает аналитику по доставкам
// @Summary Get delivery analytics
// @Description Get analytics and statistics for deliveries
// @Tags admin-delivery
// @Produce json
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param provider_id query int false "Provider ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=service.DeliveryAnalytics} "Analytics data"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/admin/delivery/analytics [get]
func (h *Handler) GetAnalytics(c *fiber.Ctx) error {
	// Парсим даты
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_from_date", fiber.Map{
				"error": err.Error(),
			})
		}
	} else {
		// По умолчанию - последние 30 дней
		from = time.Now().AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_to_date", fiber.Map{
				"error": err.Error(),
			})
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
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "error.invalid_provider_id", fiber.Map{
				"error": err.Error(),
			})
		}
		providerID = &pID
	}

	analytics, err := h.service.GetDeliveryAnalytics(c.Context(), from, to, providerID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "error.failed_to_get_analytics", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, analytics, "Аналитика доставок")
}

// HandleTrackingWebhook - обработка webhook для отслеживания
// @Summary Handle tracking webhook
// @Description Process tracking status updates from delivery providers
// @Tags delivery
// @Accept json
// @Produce json
// @Param provider path string true "Provider code (post_express, bex_express, etc.)"
// @Param payload body object true "Webhook payload from provider"
// @Success 200 {object} utils.SuccessResponseSwag{data=object} "Webhook processed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} utils.ErrorResponseSwag "Provider not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/delivery/webhooks/{provider}/tracking [post]
func (h *Handler) HandleTrackingWebhook(c *fiber.Ctx) error {
	providerCode := c.Params("provider")
	if providerCode == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.webhook.provider_required", nil)
	}

	// Получаем тело запроса
	payload := c.Body()
	if len(payload) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.webhook.empty_payload", nil)
	}

	// Собираем заголовки
	headers := make(map[string]string)
	c.GetReqHeaders()
	for key, values := range c.GetReqHeaders() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Обрабатываем webhook через сервис
	result, err := h.service.HandleProviderWebhook(c.Context(), providerCode, payload, headers)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "delivery.webhook.processing_error", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, result, "Webhook processed successfully")
}

// Структуры запросов

// CartCalculationRequest - запрос расчета для корзины
type CartCalculationRequest struct {
	FromLocation   calculator.Location        `json:"from_location"`
	ToLocation     calculator.Location        `json:"to_location"`
	Items          []calculator.ItemWithAttrs `json:"items"`
	InsuranceValue float64                    `json:"insurance_value,omitempty"`
	CODAmount      float64                    `json:"cod_amount,omitempty"`
	DeliveryType   string                     `json:"delivery_type,omitempty"`
}

// CancelRequest - запрос отмены
type CancelRequest struct {
	Reason string `json:"reason"`
}
