// backend/internal/proj/users/handler/handler.go
package handler

import (
	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Auth *AuthHandler
	User *UserHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Auth: NewAuthHandler(services),
		User: NewUserHandler(services),
	}
}
