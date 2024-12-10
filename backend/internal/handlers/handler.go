//backend/internal/handlers/handler.go
package handlers

import (
    "backend/internal/services"
)

type Handler struct {
    Auth     *AuthHandler
    Rooms    *RoomHandler
    Bookings *BookingHandler
    Users    *UserHandler
    Cars     *CarHandler 
    Marketplace *MarketplaceHandler
}

func NewHandler(services services.ServicesInterface) *Handler {
    return &Handler{
        Auth:     NewAuthHandler(services),
        Rooms:    NewRoomHandler(services),
        Bookings: NewBookingHandler(services),
        Users:    NewUserHandler(services),
        Cars:     NewCarHandler(services),
        Marketplace: NewMarketplaceHandler(services),
    }
}