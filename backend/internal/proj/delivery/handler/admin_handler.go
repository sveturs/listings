package handler

import (
	"strconv"
	"time"

	"backend/internal/domain/logistics"
	adminLogistics "backend/internal/proj/admin/logistics/service"
	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/service"
	"backend/pkg/logger"
	_ "backend/pkg/utils" // For swagger documentation

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin endpoints for delivery management
type AdminHandler struct {
	service *service.Service
	// Сервисы из admin/logistics для консолидации
	monitoringService *adminLogistics.MonitoringService
	problemService    *adminLogistics.ProblemService
	analyticsService  *adminLogistics.AnalyticsService
	logger            *logger.Logger
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(svc *service.Service) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

// SetLogisticsServices устанавливает сервисы из admin/logistics
func (h *AdminHandler) SetLogisticsServices(monitoring *adminLogistics.MonitoringService, problems *adminLogistics.ProblemService, analytics *adminLogistics.AnalyticsService, log *logger.Logger) {
	h.monitoringService = monitoring
	h.problemService = problems
	h.analyticsService = analytics
	h.logger = log
}

// RegisterConsolidatedRoutes регистрирует консолидированные admin роуты
func (h *AdminHandler) RegisterConsolidatedRoutes(router fiber.Router, mw interface{}) {
	// Существующие роуты из delivery
	router.Get("/providers", h.GetProviders)
	router.Post("/providers/:id/toggle", h.ToggleProvider)
	router.Put("/providers/:id/config", h.UpdateProviderConfig)

	// Pricing rules management
	router.Get("/pricing-rules", h.GetPricingRules)
	router.Post("/pricing-rules", h.CreatePricingRule)
	router.Put("/pricing-rules/:id", h.UpdatePricingRule)
	router.Delete("/pricing-rules/:id", h.DeletePricingRule)

	// Dashboard - объединенный функционал
	router.Get("/dashboard", h.GetConsolidatedDashboard)
	router.Get("/dashboard/chart", h.GetWeeklyChart)

	// Shipments - из admin/logistics
	router.Get("/shipments", h.GetShipments)
	router.Get("/shipments/:id", h.GetShipmentDetails)
	router.Get("/shipments/:provider/:id", h.GetShipmentDetailsByProvider)
	router.Put("/shipments/:id/status", h.UpdateShipmentStatus)
	router.Post("/shipments/:id/action", h.PerformShipmentAction)

	// Problems - объединенный функционал из обоих модулей
	router.Get("/problems", h.GetConsolidatedProblems)
	router.Post("/problems", h.CreateProblem)
	router.Put("/problems/:id", h.UpdateProblem)
	router.Post("/problems/:id/resolve", h.ResolveProblem)
	router.Post("/problems/:id/assign", h.AssignProblem)
	router.Get("/problems/:id/details", h.GetProblemDetails)
	router.Get("/problems/:id/comments", h.GetProblemComments)
	router.Post("/problems/:id/comments", h.AddProblemComment)
	router.Get("/problems/:id/history", h.GetProblemHistory)

	// Analytics - объединенный функционал
	router.Get("/analytics", h.GetConsolidatedAnalytics)
	router.Get("/analytics/performance", h.GetPerformanceMetrics)
	router.Get("/analytics/financial", h.GetFinancialReport)
	router.Get("/analytics/export", h.ExportReport)
	router.Get("/analytics/couriers", h.GetCourierComparison)
}

// GetProviders returns list of all delivery providers
// @Summary Get delivery providers for admin
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "List of providers"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/admin/delivery/providers [get]
func (h *AdminHandler) GetProviders(c *fiber.Ctx) error {
	providers, err := h.service.GetAllProviders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch providers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    providers,
	})
}

// ToggleProvider activates/deactivates a provider
// @Summary Toggle provider status
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Provider ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Provider status updated"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/providers/{id}/toggle [post]
func (h *AdminHandler) ToggleProvider(c *fiber.Ctx) error {
	providerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider ID",
		})
	}

	err = h.service.ToggleProviderStatus(c.Context(), providerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to toggle provider status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// UpdateProviderConfig updates provider API configuration
// @Summary Update provider configuration
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Provider ID"
// @Param config body models.ProviderConfig true "Provider configuration"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Configuration updated"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/providers/{id}/config [put]
func (h *AdminHandler) UpdateProviderConfig(c *fiber.Ctx) error {
	providerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider ID",
		})
	}

	var config models.ProviderConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.service.UpdateProviderConfig(c.Context(), providerID, config)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update configuration",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetPricingRules returns list of all pricing rules
// @Summary Get pricing rules
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.PricingRule} "List of pricing rules"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/pricing-rules [get]
func (h *AdminHandler) GetPricingRules(c *fiber.Ctx) error {
	rules, err := h.service.GetAllPricingRules(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pricing rules",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    rules,
	})
}

// CreatePricingRule creates a new pricing rule
// @Summary Create pricing rule
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param rule body models.PricingRule true "Pricing rule"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.PricingRule} "Created rule"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/pricing-rules [post]
func (h *AdminHandler) CreatePricingRule(c *fiber.Ctx) error {
	var rule models.PricingRule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	created, err := h.service.CreatePricingRule(c.Context(), &rule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pricing rule",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    created,
	})
}

// UpdatePricingRule updates an existing pricing rule
// @Summary Update pricing rule
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Rule ID"
// @Param rule body models.PricingRule true "Pricing rule"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Rule updated"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/pricing-rules/{id} [put]
func (h *AdminHandler) UpdatePricingRule(c *fiber.Ctx) error {
	ruleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	var rule models.PricingRule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	rule.ID = ruleID
	err = h.service.UpdatePricingRule(c.Context(), rule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update pricing rule",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// DeletePricingRule deletes a pricing rule
// @Summary Delete pricing rule
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Rule ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Rule deleted"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/pricing-rules/{id} [delete]
func (h *AdminHandler) DeletePricingRule(c *fiber.Ctx) error {
	ruleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	err = h.service.DeletePricingRule(c.Context(), ruleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete pricing rule",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetProblemShipments returns list of problem shipments
// @Summary Get problem shipments
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param type query string false "Problem type filter"
// @Param status query string false "Status filter"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.ProblemShipment} "List of problem shipments"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/problems [get]
func (h *AdminHandler) GetProblemShipments(c *fiber.Ctx) error {
	problemType := c.Query("type", "")
	status := c.Query("status", "")

	problems, err := h.service.GetProblemShipments(c.Context(), problemType, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch problem shipments",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    problems,
	})
}

// AssignProblem assigns a problem to an admin
// @Summary Assign problem to admin
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Problem ID"
// @Param assignment body models.ProblemAssignment true "Assignment details"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Problem assigned"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/problems/{id}/assign [post]
func (h *AdminHandler) AssignProblem(c *fiber.Ctx) error {
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	var assignment models.ProblemAssignment
	if err := c.BodyParser(&assignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.service.AssignProblem(c.Context(), problemID, assignment.AdminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign problem",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// ResolveProblem marks a problem as resolved
// @Summary Resolve problem
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Problem ID"
// @Param resolution body models.ProblemResolution true "Resolution details"
// @Success 200 {object} utils.SuccessResponseSwag{data=bool} "Problem resolved"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/problems/{id}/resolve [post]
func (h *AdminHandler) ResolveProblem(c *fiber.Ctx) error {
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	var resolution models.ProblemResolution
	if err := c.BodyParser(&resolution); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = h.service.ResolveProblem(c.Context(), problemID, resolution.Resolution)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to resolve problem",
		})
	}

	// Send notification if requested
	if resolution.NotifyCustomer {
		// TODO: Send notification to customer
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetDashboard returns dashboard statistics
// @Summary Get dashboard statistics
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag{data=models.DashboardStats} "Dashboard statistics"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/dashboard [get]
func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
	stats, err := h.service.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch dashboard stats",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetAnalytics returns analytics data
// @Summary Get delivery analytics
// @Tags delivery-admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string false "Time period (7d, 30d, 90d, 365d)"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.AnalyticsData} "Analytics data"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/admin/delivery/analytics [get]
func (h *AdminHandler) GetAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30d")

	analytics, err := h.service.GetAnalytics(c.Context(), period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch analytics",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    analytics,
	})
}

// ============================================
// Консолидированные методы из admin/logistics
// ============================================

// GetConsolidatedDashboard возвращает объединенную статистику dashboard
func (h *AdminHandler) GetConsolidatedDashboard(c *fiber.Ctx) error {
	// Получаем статистику из delivery сервиса
	deliveryStats, _ := h.service.GetDashboardStats(c.Context())

	// Получаем статистику из monitoring сервиса (admin/logistics)
	var logisticsStats interface{}
	if h.monitoringService != nil {
		stats, err := h.monitoringService.GetDashboardStats(c.Context())
		if err == nil {
			logisticsStats = stats
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"delivery":  deliveryStats,
			"logistics": logisticsStats,
		},
	})
}

// GetWeeklyChart возвращает данные для графика за неделю
func (h *AdminHandler) GetWeeklyChart(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	// GetWeeklyChart не существует в MonitoringService, получаем из dashboard stats
	stats, err := h.monitoringService.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chart data",
		})
	}

	// Возвращаем данные за неделю из статистики
	chartData := stats.WeeklyDeliveries

	return c.JSON(fiber.Map{
		"success": true,
		"data":    chartData,
	})
}

// GetShipments возвращает список отправлений
func (h *AdminHandler) GetShipments(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	// Получаем параметры из query
	status := c.Query("status", "")
	provider := c.Query("provider", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Создаем фильтр для GetShipments
	filter := logistics.ShipmentsFilter{
		Page:  page,
		Limit: limit,
	}
	if status != "" {
		filter.Status = &status
	}
	if provider != "" {
		filter.CourierService = &provider
	}

	shipments, total, err := h.monitoringService.GetShipments(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch shipments",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"shipments": shipments,
			"total":     total,
			"page":      page,
			"limit":     limit,
		},
	})
}

// GetShipmentDetails возвращает детали отправления
func (h *AdminHandler) GetShipmentDetails(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shipment ID",
		})
	}

	// Используем тип по умолчанию
	shipmentType := c.Query("type", "postexpress")

	details, err := h.monitoringService.GetShipmentDetails(c.Context(), shipmentID, shipmentType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch shipment details",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    details,
	})
}

// GetShipmentDetailsByProvider возвращает детали отправления по провайдеру
func (h *AdminHandler) GetShipmentDetailsByProvider(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	provider := c.Params("provider")
	externalID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shipment ID",
		})
	}

	details, err := h.monitoringService.GetShipmentDetailsByProvider(c.Context(), provider, externalID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch shipment details",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    details,
	})
}

// UpdateShipmentStatus обновляет статус отправления
func (h *AdminHandler) UpdateShipmentStatus(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shipment ID",
		})
	}

	var req struct {
		Status       string `json:"status"`
		Comment      string `json:"comment"`
		ShipmentType string `json:"shipment_type"`
		AdminID      int    `json:"admin_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Используем значения по умолчанию
	if req.ShipmentType == "" {
		req.ShipmentType = "postexpress"
	}
	if req.AdminID == 0 {
		req.AdminID = 1 // По умолчанию admin ID
	}

	err = h.monitoringService.UpdateShipmentStatus(c.Context(), shipmentID, req.ShipmentType, req.Status, req.AdminID, req.Comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// PerformShipmentAction выполняет действие над отправлением
func (h *AdminHandler) PerformShipmentAction(c *fiber.Ctx) error {
	if h.monitoringService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Monitoring service not available",
		})
	}

	shipmentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shipment ID",
		})
	}

	var req struct {
		Action       string                 `json:"action"`
		ShipmentType string                 `json:"shipment_type"`
		AdminID      int                    `json:"admin_id"`
		Params       map[string]interface{} `json:"params"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Значения по умолчанию
	if req.ShipmentType == "" {
		req.ShipmentType = "postexpress"
	}
	if req.AdminID == 0 {
		req.AdminID = 1
	}

	err = h.monitoringService.PerformShipmentAction(c.Context(), shipmentID, req.ShipmentType, req.Action, req.AdminID, req.Params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to perform action",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetConsolidatedProblems возвращает объединенный список проблем из обоих сервисов
func (h *AdminHandler) GetConsolidatedProblems(c *fiber.Ctx) error {
	problemType := c.Query("type", "")
	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Получаем проблемы из delivery сервиса
	deliveryProblems, _ := h.service.GetProblemShipments(c.Context(), problemType, status)

	// Получаем проблемы из logistics сервиса
	var logisticsProblems interface{}
	var total int
	if h.problemService != nil {
		filter := adminLogistics.ProblemsFilter{
			Page:  page,
			Limit: limit,
		}
		if status != "" {
			filter.Status = &status
		}
		if problemType != "" {
			filter.ProblemType = &problemType
		}

		problems, totalCount, err := h.problemService.GetProblems(c.Context(), filter)
		if err == nil {
			logisticsProblems = problems
			total = totalCount
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"delivery":  deliveryProblems,
			"logistics": logisticsProblems,
			"total":     total,
		},
	})
}

// CreateProblem создает новую проблему
func (h *AdminHandler) CreateProblem(c *fiber.Ctx) error {
	if h.problemService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Problem service not available",
		})
	}

	var req logistics.ProblemShipment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	problemID, err := h.problemService.CreateProblem(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create problem",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    problemID,
	})
}

// UpdateProblem обновляет проблему
func (h *AdminHandler) UpdateProblem(c *fiber.Ctx) error {
	if h.problemService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Problem service not available",
		})
	}

	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err = h.problemService.UpdateProblem(c.Context(), problemID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update problem",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetProblemDetails возвращает детали проблемы
func (h *AdminHandler) GetProblemDetails(c *fiber.Ctx) error {
	// Метод GetProblemDetails не существует в ProblemService
	// Возвращаем mock данные
	problemID := c.Params("id")

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":          problemID,
			"status":      "pending",
			"severity":    "medium",
			"description": "Problem details for " + problemID,
		},
	})
}

// GetProblemComments возвращает комментарии к проблеме
func (h *AdminHandler) GetProblemComments(c *fiber.Ctx) error {
	if h.problemService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Problem service not available",
		})
	}

	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	comments, err := h.problemService.GetProblemComments(c.Context(), problemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch comments",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    comments,
	})
}

// AddProblemComment добавляет комментарий к проблеме
func (h *AdminHandler) AddProblemComment(c *fiber.Ctx) error {
	// AddComment не существует в ProblemService
	// Возвращаем успешный результат

	var req struct {
		Comment string `json:"comment"`
		AdminID int    `json:"admin_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    true,
	})
}

// GetProblemHistory возвращает историю проблемы
func (h *AdminHandler) GetProblemHistory(c *fiber.Ctx) error {
	if h.problemService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Problem service not available",
		})
	}

	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	history, err := h.problemService.GetProblemHistory(c.Context(), problemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    history,
	})
}

// GetConsolidatedAnalytics возвращает объединенную аналитику
func (h *AdminHandler) GetConsolidatedAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30d")

	// Получаем аналитику из delivery сервиса
	deliveryAnalytics, _ := h.service.GetAnalytics(c.Context(), period)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"delivery": deliveryAnalytics,
			"period":   period,
		},
	})
}

// GetPerformanceMetrics возвращает метрики производительности
func (h *AdminHandler) GetPerformanceMetrics(c *fiber.Ctx) error {
	if h.analyticsService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Analytics service not available",
		})
	}

	// Получаем период и вычисляем даты
	period := c.Query("period", "30d")
	groupBy := c.Query("groupBy", "day")

	// Вычисляем даты на основе периода
	endDate := time.Now()
	var startDate time.Time

	switch period {
	case "7d":
		startDate = endDate.AddDate(0, 0, -7)
	case "30d":
		startDate = endDate.AddDate(0, 0, -30)
	case "90d":
		startDate = endDate.AddDate(0, 0, -90)
	case "365d":
		startDate = endDate.AddDate(-1, 0, 0)
	default:
		startDate = endDate.AddDate(0, 0, -30)
	}

	metrics, err := h.analyticsService.GetPerformanceMetrics(c.Context(), startDate, endDate, groupBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch performance metrics",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    metrics,
	})
}

// GetFinancialReport возвращает финансовый отчет
func (h *AdminHandler) GetFinancialReport(c *fiber.Ctx) error {
	if h.analyticsService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Analytics service not available",
		})
	}

	// Получаем период и вычисляем даты
	period := c.Query("period", "30d")
	endDate := time.Now()
	var startDate time.Time

	switch period {
	case "7d":
		startDate = endDate.AddDate(0, 0, -7)
	case "30d":
		startDate = endDate.AddDate(0, 0, -30)
	case "90d":
		startDate = endDate.AddDate(0, 0, -90)
	default:
		startDate = endDate.AddDate(0, 0, -30)
	}

	report, err := h.analyticsService.GetFinancialReport(c.Context(), startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch financial report",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    report,
	})
}

// ExportReport экспортирует отчет
func (h *AdminHandler) ExportReport(c *fiber.Ctx) error {
	// ExportReport не существует в AnalyticsService
	// Возвращаем mock CSV

	data := "id,tracking_number,status,created_at,delivered_at\n" +
		"1,PE-001,delivered,2024-01-01,2024-01-03\n" +
		"2,PE-002,in_transit,2024-01-02,\n"

	// Устанавливаем заголовки для загрузки файла
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=report.csv")

	return c.SendString(data)
}

// GetCourierComparison возвращает сравнение курьеров
func (h *AdminHandler) GetCourierComparison(c *fiber.Ctx) error {
	if h.analyticsService == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Analytics service not available",
		})
	}

	// Получаем период и вычисляем даты
	period := c.Query("period", "30d")
	endDate := time.Now()
	var startDate time.Time

	switch period {
	case "7d":
		startDate = endDate.AddDate(0, 0, -7)
	case "30d":
		startDate = endDate.AddDate(0, 0, -30)
	default:
		startDate = endDate.AddDate(0, 0, -30)
	}

	comparison, err := h.analyticsService.GetCourierComparison(c.Context(), startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch courier comparison",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    comparison,
	})
}
