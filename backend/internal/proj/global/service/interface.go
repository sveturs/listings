package service

import (
    userService "backend/internal/proj/users/service"
    accommodationService "backend/internal/proj/accommodation/service"
    carService "backend/internal/proj/car/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
    "backend/internal/config"
)

type ServicesInterface interface {
    Auth() userService.AuthServiceInterface
    Room() accommodationService.RoomServiceInterface
    Booking() accommodationService.BookingServiceInterface
    User() userService.UserServiceInterface
    Car() carService.CarServiceInterface
    Config() *config.Config
    Marketplace() marketplaceService.MarketplaceServiceInterface
    Review() reviewService.ReviewServiceInterface
}