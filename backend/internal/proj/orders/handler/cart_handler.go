package handler

import (
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AddToCart добавляет товар в корзину
// @Summary Add item to cart
// @Description Adds a product to the shopping cart for a specific storefront
// @Tags cart
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param item body models.AddToCartRequest true "Item to add"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.ShoppingCart} "Updated cart"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Product not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/cart/items [post]
func (h *OrdersHandler) AddToCart(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	_, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_storefront_id")
	}

	var req models.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse cart request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_request_body")
	}

	// Получаем user_id или session_id для корзины
	var userID *int
	var sessionID *string

	if userIDRaw := c.Locals("user_id"); userIDRaw != nil {
		userIDVal := userIDRaw.(int)
		userID = &userIDVal
	} else {
		sessionIDVal := c.Get("X-Session-ID")
		if sessionIDVal == "" {
			// Создаем временную сессию для анонимного пользователя
			sessionIDVal = "anon_" + c.IP() + "_" + storefrontIDStr
		}
		sessionID = &sessionIDVal
	}

	storefrontID, _ := strconv.Atoi(storefrontIDStr)

	// Создаем элемент корзины
	cartItem := &models.ShoppingCartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	cart, err := h.orderService.AddToCartWithDetails(c.Context(), cartItem, storefrontID, userID, sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to add item to cart")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.add_failed")
	}

	return utils.SuccessResponse(c, cart)
}

// UpdateCartItem обновляет количество товара в корзине
// @Summary Update cart item quantity
// @Description Updates the quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param item_id path int true "Cart item ID"
// @Param update body models.UpdateCartItemRequest true "Updated quantity"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.ShoppingCart} "Updated cart"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Item not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/cart/items/{item_id} [put]
func (h *OrdersHandler) UpdateCartItem(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_storefront_id")
	}

	itemIDStr := c.Params("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_item_id")
	}

	var req models.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse update cart request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_request_body")
	}

	// Получаем user_id или session_id для корзины
	var userID *int
	var sessionID *string

	if userIDRaw := c.Locals("user_id"); userIDRaw != nil {
		userIDVal := userIDRaw.(int)
		userID = &userIDVal
	} else {
		sessionIDVal := c.Get("X-Session-ID")
		if sessionIDVal == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.missing_session")
		}
		sessionID = &sessionIDVal
	}

	cart, err := h.orderService.UpdateCartItemQuantity(c.Context(), itemID, storefrontID, req.Quantity, userID, sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update cart item")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.update_failed")
	}

	return utils.SuccessResponse(c, cart)
}

// RemoveFromCart удаляет товар из корзины
// @Summary Remove item from cart
// @Description Removes an item from the shopping cart
// @Tags cart
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param item_id path int true "Cart item ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.ShoppingCart} "Updated cart"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Item not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/cart/items/{item_id} [delete]
func (h *OrdersHandler) RemoveFromCart(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_storefront_id")
	}

	itemIDStr := c.Params("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_item_id")
	}

	// Получаем user_id или session_id для корзины
	var userID *int
	var sessionID *string

	if userIDRaw := c.Locals("user_id"); userIDRaw != nil {
		userIDVal := userIDRaw.(int)
		userID = &userIDVal
	} else {
		sessionIDVal := c.Get("X-Session-ID")
		if sessionIDVal == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.missing_session")
		}
		sessionID = &sessionIDVal
	}

	cart, err := h.orderService.RemoveFromCart(c.Context(), itemID, storefrontID, userID, sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to remove item from cart")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.remove_failed")
	}

	return utils.SuccessResponse(c, cart)
}

// GetCart возвращает содержимое корзины
// @Summary Get shopping cart
// @Description Gets the current contents of the shopping cart for a storefront
// @Tags cart
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.ShoppingCart} "Cart contents"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/cart [get]
func (h *OrdersHandler) GetCart(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_storefront_id")
	}

	// Получаем user_id или session_id для корзины
	var userID *int
	var sessionID *string

	if userIDRaw := c.Locals("user_id"); userIDRaw != nil {
		userIDVal := userIDRaw.(int)
		userID = &userIDVal
	} else {
		sessionIDVal := c.Get("X-Session-ID")
		if sessionIDVal == "" {
			// Создаем временную сессию для анонимного пользователя
			sessionIDVal = "anon_" + c.IP() + "_" + storefrontIDStr
		}
		sessionID = &sessionIDVal
	}

	cart, err := h.orderService.GetCart(c.Context(), storefrontID, userID, sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get cart")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.get_failed")
	}

	return utils.SuccessResponse(c, cart)
}

// ClearCart очищает корзину
// @Summary Clear shopping cart
// @Description Removes all items from the shopping cart
// @Tags cart
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=string} "Cart cleared"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/cart [delete]
func (h *OrdersHandler) ClearCart(c *fiber.Ctx) error {
	storefrontIDStr := c.Params("storefront_id")
	storefrontID, err := strconv.Atoi(storefrontIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.invalid_storefront_id")
	}

	// Получаем user_id или session_id для корзины
	var userID *int
	var sessionID *string

	if userIDRaw := c.Locals("user_id"); userIDRaw != nil {
		userIDVal := userIDRaw.(int)
		userID = &userIDVal
	} else {
		sessionIDVal := c.Get("X-Session-ID")
		if sessionIDVal == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "cart.error.missing_session")
		}
		sessionID = &sessionIDVal
	}

	err = h.orderService.ClearCart(c.Context(), storefrontID, userID, sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to clear cart")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.clear_failed")
	}

	return utils.SuccessResponse(c, "cart.cleared")
}

// GetUserCarts возвращает все корзины пользователя
// @Summary Get all user carts
// @Description Gets all shopping carts for the authenticated user across all storefronts
// @Tags cart
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.ShoppingCart} "User's carts"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/carts [get]
func (h *OrdersHandler) GetUserCarts(c *fiber.Ctx) error {
	// Получаем user_id из контекста (только для авторизованных)
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	userID := userIDRaw.(int)

	carts, err := h.orderService.GetUserCarts(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user carts")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "cart.error.get_user_carts_failed")
	}

	return utils.SuccessResponse(c, carts)
}
