package handler

import (
	"backend/internal/proj/orders/service"

	"github.com/jmoiron/sqlx"
)

// OrdersHandler основной handler для всех операций с заказами
type OrdersHandler struct {
	orderService     *service.OrderService
	inventoryManager *service.InventoryManager
	db               *sqlx.DB
}

// NewOrdersHandler создает новый handler
func NewOrdersHandler(orderService *service.OrderService, inventoryManager *service.InventoryManager, db *sqlx.DB) *OrdersHandler {
	return &OrdersHandler{
		orderService:     orderService,
		inventoryManager: inventoryManager,
		db:               db,
	}
}
