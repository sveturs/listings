// backend/internal/proj/users/handler/admin_balance.go
package handler

import (
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetUserBalance получает баланс пользователя (админский метод)
func (h *UserHandler) GetUserBalance(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	balance, err := h.services.Balance().GetBalance(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения баланса пользователя")
	}

	return utils.SuccessResponse(c, balance)
}

// GetUserTransactions получает транзакции пользователя (админский метод)
func (h *UserHandler) GetUserTransactions(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	limit := utils.StringToInt(c.Query("limit", "20"), 20)
	offset := utils.StringToInt(c.Query("offset", "0"), 0)

	transactions, err := h.services.Balance().GetTransactions(c.Context(), id, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения транзакций пользователя")
	}

	return utils.SuccessResponse(c, transactions)
}
