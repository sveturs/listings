package handler

import (
    globalService "backend/internal/proj/global/service"
    
    
)

type Handler struct {
    Room    *RoomHandler
    Booking *BookingHandler
}

func NewHandler(services globalService.ServicesInterface) *Handler {
    return &Handler{
        Room:    NewRoomHandler(services),
        Booking: NewBookingHandler(services), 
    }
}