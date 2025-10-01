package handler

import (
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// DashboardHandler обрабатывает запросы к dashboard витрины
type DashboardHandler struct {
	service        service.StorefrontService
	productService *service.ProductService
	repository     postgres.StorefrontRepository
}

// NewDashboardHandler создает новый dashboard handler
func NewDashboardHandler(storefrontService service.StorefrontService, productService *service.ProductService, repository postgres.StorefrontRepository) *DashboardHandler {
	return &DashboardHandler{
		service:        storefrontService,
		productService: productService,
		repository:     repository,
	}
}

// DashboardNotification уведомление для dashboard
type DashboardNotification struct {
	ID        int    `json:"id"`
	Type      string `json:"type"` // order, message, stock
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	IsRead    bool   `json:"is_read"`
}

// GetDashboardStats получает статистику для dashboard
// @Summary Get dashboard statistics
// @Description Returns statistics for storefront dashboard
// @Tags storefronts,dashboard
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Success 200 {object} postgres.DashboardStats "Dashboard statistics"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Storefront not found"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{slug}/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID, _ := authMiddleware.GetUserID(c)

	// Проверяем доступ к витрине
	storefront, err := h.service.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
	}

	if storefront.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.access_denied")
	}

	// Получаем реальную статистику из БД
	stats, err := h.repository.GetDashboardStats(c.Context(), storefront.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.stats_failed")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetRecentOrders получает последние заказы
// @Summary Get recent orders
// @Description Returns recent orders for storefront dashboard
// @Tags storefronts,dashboard
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param limit query int false "Limit" default(5)
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]postgres.DashboardOrder} "Recent orders"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Storefront not found"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{slug}/dashboard/recent-orders [get]
func (h *DashboardHandler) GetRecentOrders(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID, _ := authMiddleware.GetUserID(c)

	// Проверяем доступ к витрине
	storefront, err := h.service.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
	}

	if storefront.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.access_denied")
	}

	limit := c.QueryInt("limit", 5)
	if limit > 20 {
		limit = 20
	}

	// Получаем реальные заказы из БД
	orders, err := h.repository.GetRecentOrders(c.Context(), storefront.ID, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.orders_failed")
	}

	// Если заказов нет, возвращаем пустой массив
	if orders == nil {
		orders = []*postgres.DashboardOrder{}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    orders,
		"total":   len(orders),
	})
}

// GetLowStockProducts получает товары с низким запасом
// @Summary Get low stock products
// @Description Returns products with low stock for storefront dashboard
// @Tags storefronts,dashboard
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]postgres.LowStockProduct} "Low stock products"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Storefront not found"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{slug}/dashboard/low-stock [get]
func (h *DashboardHandler) GetLowStockProducts(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID, _ := authMiddleware.GetUserID(c)

	// Проверяем доступ к витрине
	storefront, err := h.service.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
	}

	if storefront.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.access_denied")
	}

	// Получаем реальные товары с низким запасом из БД
	products, err := h.repository.GetLowStockProducts(c.Context(), storefront.ID, 10)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.low_stock_failed")
	}

	// Если товаров нет, возвращаем пустой массив
	if products == nil {
		products = []*postgres.LowStockProduct{}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}

// GetDashboardNotifications получает уведомления для dashboard
// @Summary Get dashboard notifications
// @Description Returns notifications for storefront dashboard
// @Tags storefronts,dashboard
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]handler.DashboardNotification} "Dashboard notifications"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Storefront not found"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Access denied"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{slug}/dashboard/notifications [get]
func (h *DashboardHandler) GetDashboardNotifications(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID, _ := authMiddleware.GetUserID(c)

	// Проверяем доступ к витрине
	storefront, err := h.service.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
	}

	if storefront.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.access_denied")
	}

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	// Пока возвращаем базовые уведомления
	// TODO: В будущем создать более продвинутую систему уведомлений для витрин
	allNotifications := []DashboardNotification{
		{
			ID:        1,
			Type:      "order",
			Title:     "notifications.newOrder",
			Message:   "notifications.newOrderMessage",
			CreatedAt: "2024-01-15T10:30:00Z",
			IsRead:    false,
		},
		{
			ID:        2,
			Type:      "stock",
			Title:     "notifications.lowStock",
			Message:   "notifications.lowStockMessage",
			CreatedAt: "2024-01-15T09:00:00Z",
			IsRead:    false,
		},
		{
			ID:        3,
			Type:      "message",
			Title:     "notifications.newMessage",
			Message:   "notifications.newMessageFromCustomer",
			CreatedAt: "2024-01-14T15:45:00Z",
			IsRead:    true,
		},
	}

	// Обрезаем массив согласно лимиту
	notifications := allNotifications
	if len(allNotifications) > limit {
		notifications = allNotifications[:limit]
	}

	// Подсчитываем непрочитанные
	unreadCount := 0
	for _, n := range notifications {
		if !n.IsRead {
			unreadCount++
		}
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"data":         notifications,
		"unread_count": unreadCount,
	})
}
