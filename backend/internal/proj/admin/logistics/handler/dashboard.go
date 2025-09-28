package handler

import (
	"backend/internal/proj/admin/logistics/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// DashboardHandler обработчик для dashboard метрик
type DashboardHandler struct {
	monitoringService *service.MonitoringService
}

// NewDashboardHandler создает новый обработчик dashboard
func NewDashboardHandler(monitoringService *service.MonitoringService) *DashboardHandler {
	return &DashboardHandler{
		monitoringService: monitoringService,
	}
}

// GetDashboardStats godoc
// @Summary Получить статистику для dashboard логистики
// @Description Возвращает агрегированные метрики и статистику для отображения на dashboard
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Dashboard statistics"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden - insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/dashboard [get]
func (h *DashboardHandler) GetDashboardStats(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Получаем статистику
	stats, err := h.monitoringService.GetDashboardStats(c.Context())
	if err != nil {
		// Логируем ошибку для отладки
		_ = c.App().Config().ErrorHandler(c, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.dashboard.error")
	}

	return utils.SuccessResponse(c, stats)
}

// GetWeeklyChart godoc
// @Summary Получить данные для недельного графика
// @Description Возвращает статистику доставок за последние 7 дней для графиков
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=[]interface{}} "Weekly chart data"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/dashboard/chart [get]
func (h *DashboardHandler) GetWeeklyChart(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Получаем полную статистику и извлекаем недельные данные
	stats, err := h.monitoringService.GetDashboardStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.chart.error")
	}

	return utils.SuccessResponse(c, stats.WeeklyDeliveries)
}
