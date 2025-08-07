package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes регистрирует все маршруты для заказов
func (h *OrdersHandler) RegisterRoutes(app *fiber.App, m *middleware.Middleware) {
	// Группа для API заказов
	api := app.Group("/api/v1")

	// Корзина покупок (требует аутентификации или session)
	cartGroup := api.Group("/storefronts/:storefront_id/cart")
	cartGroup.Use(m.OptionalAuth()) // Опциональная аутентификация для поддержки анонимных корзин
	{
		cartGroup.Post("/items", h.AddToCart)                 // POST /api/v1/storefronts/:id/cart/items
		cartGroup.Put("/items/:item_id", h.UpdateCartItem)    // PUT /api/v1/storefronts/:id/cart/items/:item_id
		cartGroup.Delete("/items/:item_id", h.RemoveFromCart) // DELETE /api/v1/storefronts/:id/cart/items/:item_id
		cartGroup.Get("/", h.GetCart)                         // GET /api/v1/storefronts/:id/cart
		cartGroup.Delete("/", h.ClearCart)                    // DELETE /api/v1/storefronts/:id/cart
	}

	// Корзины пользователя (требует аутентификации)
	userCartsGroup := api.Group("/user")
	userCartsGroup.Use(m.RequireAuth())
	{
		userCartsGroup.Get("/carts", h.GetUserCarts) // GET /api/v1/user/carts
	}

	// Заказы (требует полную аутентификацию)
	ordersGroup := api.Group("/orders")
	ordersGroup.Use(m.RequireAuth())
	{
		ordersGroup.Post("/", h.CreateOrder)          // POST /api/v1/orders
		ordersGroup.Get("/", h.GetMyOrders)           // GET /api/v1/orders
		ordersGroup.Get("/:id", h.GetOrder)           // GET /api/v1/orders/:id
		ordersGroup.Put("/:id/cancel", h.CancelOrder) // PUT /api/v1/orders/:id/cancel
	}

	// Управление заказами витрины (для продавцов)
	storefrontOrdersGroup := api.Group("/storefronts/:storefront_id/orders")
	storefrontOrdersGroup.Use(m.RequireAuth())
	{
		storefrontOrdersGroup.Get("/", h.GetStorefrontOrders)               // GET /api/v1/storefronts/:id/orders
		storefrontOrdersGroup.Put("/:order_id/status", h.UpdateOrderStatus) // PUT /api/v1/storefronts/:id/orders/:order_id/status
	}
}
