package handler
import (
     globalService "backend/internal/proj/global/service"
)
type Handler struct {
	Review *ReviewHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	
	return &Handler{
		Review: NewReviewHandler(services),
	}
}
