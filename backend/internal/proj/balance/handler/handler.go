// backend/internal/proj/balance/handler/handler.go

package handler

import (
	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Balance     *BalanceHandler
	jwtParserMW fiber.Handler
}

func NewHandler(services globalService.ServicesInterface, jwtParserMW fiber.Handler) *Handler {
	return &Handler{
		Balance:     NewBalanceHandler(services.Balance(), services.Payment()),
		jwtParserMW: jwtParserMW,
	}
}
