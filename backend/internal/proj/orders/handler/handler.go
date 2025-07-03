package handler

import (
	"backend/internal/proj/orders/service"
)

// OrdersHandler основной handler для всех операций с заказами
type OrdersHandler struct {
	orderService     *service.OrderService
	inventoryManager *service.InventoryManager
}

// NewOrdersHandler создает новый handler
func NewOrdersHandler(orderService *service.OrderService, inventoryManager *service.InventoryManager) *OrdersHandler {
	return &OrdersHandler{
		orderService:     orderService,
		inventoryManager: inventoryManager,
	}
}
