package handler

import (
    "backend/internal/services"
)

type Handler struct {
    Room    *RoomHandler
    Booking *BookingHandler
}

func NewHandler(services services.ServicesInterface) *Handler {
    return &Handler{
        Room:    NewRoomHandler(services),
        Booking: NewBookingHandler(services), 
    }
}