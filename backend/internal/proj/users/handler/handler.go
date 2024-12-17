//backend/internal/proj/users/handler/handler.go
package handler

import (
     "backend/internal/services"
 )

 type Handler struct {
	Auth *AuthHandler
	User *UserHandler
}

func NewHandler(services services.ServicesInterface) *Handler {
	return &Handler{
		Auth: NewAuthHandler(services),
		User: NewUserHandler(services),
	}
}