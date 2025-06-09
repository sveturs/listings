// backend/internal/proj/balance/handler/handler.go

package handler

import (
	globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Balance *BalanceHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Balance: NewBalanceHandler(services.Balance(), services.Payment()),
	}
}
