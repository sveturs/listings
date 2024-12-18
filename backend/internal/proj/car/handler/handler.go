package handler

import (
    globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Car *CarHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Car: NewCarHandler(services),
	}
}
