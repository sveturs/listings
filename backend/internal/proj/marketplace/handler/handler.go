package handler
import (
     globalService "backend/internal/proj/global/service"
)

type Handler struct {
	Marketplace *MarketplaceHandler

}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		Marketplace: NewMarketplaceHandler(services),
	}
}
