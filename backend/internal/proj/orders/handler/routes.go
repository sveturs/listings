package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// RegisterRoutes регистрирует все маршруты для заказов
func (h *OrdersHandler) RegisterRoutes(app *fiber.App, m *middleware.Middleware) {
	// Корзина покупок (опциональная аутентификация для поддержки анонимных корзин)
	// ВАЖНО: НЕ используем Group с префиксом /api/v1/b2c - это создает middleware leak!
	app.Post("/api/v1/b2c/:storefront_id/cart/items", m.JWTParser(), h.AddToCart)
	app.Put("/api/v1/b2c/:storefront_id/cart/items/:item_id", m.JWTParser(), h.UpdateCartItem)
	app.Delete("/api/v1/b2c/:storefront_id/cart/items/:item_id", m.JWTParser(), h.RemoveFromCart)
	app.Get("/api/v1/b2c/:storefront_id/cart", m.JWTParser(), h.GetCart)
	app.Delete("/api/v1/b2c/:storefront_id/cart", m.JWTParser(), h.ClearCart)

	// Корзины пользователя (требует аутентификации)
	app.Get("/api/v1/user/carts", m.JWTParser(), authMiddleware.RequireAuth(), h.GetUserCarts)

	// Заказы (требует полную аутентификацию)
	app.Post("/api/v1/orders", m.JWTParser(), authMiddleware.RequireAuth(), h.CreateOrder)
	app.Get("/api/v1/orders", m.JWTParser(), authMiddleware.RequireAuth(), h.GetMyOrders)
	app.Get("/api/v1/orders/:id", m.JWTParser(), authMiddleware.RequireAuth(), h.GetOrder)
	app.Put("/api/v1/orders/:id/cancel", m.JWTParser(), authMiddleware.RequireAuth(), h.CancelOrder)

	// Управление заказами витрины (для продавцов)
	app.Get("/api/v1/b2c/:storefront_id/orders", m.JWTParser(), authMiddleware.RequireAuth(), h.GetStorefrontOrders)
	app.Put("/api/v1/b2c/:storefront_id/orders/:order_id/status", m.JWTParser(), authMiddleware.RequireAuth(), h.UpdateOrderStatus)
}
