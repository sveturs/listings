package handler

import "backend/internal/domain/models"

// Структуры ответов для Swagger документации

// CartResponse ответ с содержимым корзины
type CartResponse struct {
	Cart *models.ShoppingCart `json:"cart"`
}

// OrderResponse ответ с заказом
type OrderResponse struct {
	Order *models.StorefrontOrder `json:"order"`
}

// OrdersListResponse ответ со списком заказов
type OrdersListResponse struct {
	Orders []models.StorefrontOrder `json:"orders"`
	Total  int                      `json:"total"`
	Page   int                      `json:"page"`
	Limit  int                      `json:"limit"`
}

// InventoryResponse ответ с информацией об остатках
type InventoryResponse struct {
	Available     int                           `json:"available"`
	Reserved      int                           `json:"reserved"`
	LowStockItems []models.LowStockItem         `json:"low_stock_items,omitempty"`
	Movements     []models.InventoryMovementDTO `json:"movements,omitempty"`
}

// ClearCartResponse ответ об очистке корзины
type ClearCartResponse struct {
	Message string `json:"message"`
}
