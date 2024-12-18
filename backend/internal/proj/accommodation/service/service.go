// backend/internal/proj/accommodation/service/service.go
package service

import (
    "backend/internal/storage"
)

type Service struct {
    Room    RoomServiceInterface
    Bed     BedServiceInterface
    Booking BookingServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Room:    NewRoomService(storage),
        Bed:     NewBedService(storage),
        Booking: NewBookingService(storage),
    }
}