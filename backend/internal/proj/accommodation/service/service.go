package service

import (
    "backend/internal/storage"
)

type Service struct {
    Room    RoomServiceInterface
    Booking BookingServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Room:    NewRoomService(storage),
        Booking: NewBookingService(storage),
    }
}