// backend/internal/proj/users/handler/handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Auth        *AuthHandler
	User        *UserHandler
	jwtParserMW fiber.Handler
}

func NewHandler(
	services globalService.ServicesInterface,
	auth *AuthHandler,
	jwtParserMW fiber.Handler,
) *Handler {
	return &Handler{
		Auth:        auth,
		jwtParserMW: jwtParserMW,
		User:        NewUserHandler(services),
	}
}
