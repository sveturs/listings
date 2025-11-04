// backend/internal/proj/users/handler/admin_balance.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	_ "backend/internal/domain/models" // For Swagger documentation
	"backend/pkg/utils"
)

// GetUserBalance returns user balance information
// @Summary Get user balance (Admin)
// @Description Returns balance information for a specific user
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.UserBalance} "User balance"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.balance.error.invalid_user_id"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.balance.error.fetch_balance_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id}/balance [get]
func (h *UserHandler) GetUserBalance(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.balance.error.invalid_user_id")
	}

	balance, err := h.services.Balance().GetBalance(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.balance.error.fetch_balance_failed")
	}

	return utils.SuccessResponse(c, balance)
}

// GetUserTransactions returns user transaction history
// @Summary Get user transactions (Admin)
// @Description Returns paginated transaction history for a specific user
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.BalanceTransaction} "Transaction history"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.balance.error.invalid_user_id"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.balance.error.fetch_transactions_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id}/transactions [get]
func (h *UserHandler) GetUserTransactions(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.balance.error.invalid_user_id")
	}

	limit := utils.StringToInt(c.Query("limit", "20"), 20)
	offset := utils.StringToInt(c.Query("offset", "0"), 0)

	transactions, err := h.services.Balance().GetTransactions(c.Context(), id, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.balance.error.fetch_transactions_failed")
	}

	return utils.SuccessResponse(c, transactions)
}
