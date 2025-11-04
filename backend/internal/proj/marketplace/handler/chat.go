// TEMPORARY: Will be moved to microservice
package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/pkg/utils"
)

// GetChats godoc
// @Summary Получить список чатов пользователя
// @Description Получить список чатов текущего пользователя (временная заглушка)
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag{data=[]interface{}}
// @Failure 401 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/chat [get]
func (h *Handler) GetChats(c *fiber.Ctx) error {
	_, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// TODO: Implement when chat microservice is ready
	// Временно возвращаем пустой массив
	return utils.SuccessResponse(c, []interface{}{})
}
