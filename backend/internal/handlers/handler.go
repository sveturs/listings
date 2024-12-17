//backend/internal/handlers/handler.go
package handlers

import (
    "backend/internal/services"
)

type Handler struct {
    Cars     *CarHandler 
    Marketplace *MarketplaceHandler
    Reviews     *ReviewHandler
}

func NewHandler(services services.ServicesInterface) *Handler {
    return &Handler{
        Cars:     NewCarHandler(services),
        Marketplace: NewMarketplaceHandler(services),
        Reviews:     NewReviewHandler(services),
    }
}