// backend/internal/proj/global/service/service.go
package service

import (
    "backend/internal/config"
    "backend/internal/storage"
    userService "backend/internal/proj/users/service"
    accommodationService "backend/internal/proj/accommodation/service"
    carService "backend/internal/proj/car/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
)

type Service struct {
    users         *userService.Service
    accommodation *accommodationService.Service
    car          *carService.Service
    marketplace  *marketplaceService.Service
    review       *reviewService.Service
    config       *config.Config
}

func NewService(storage storage.Storage, cfg *config.Config) *Service {
    return &Service{
        users:         userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
        accommodation: accommodationService.NewService(storage),
        car:          carService.NewService(storage),
        marketplace:  marketplaceService.NewService(storage),
        review:       reviewService.NewService(storage),
        config:       cfg,
    }
}

// Реализация методов интерфейса
func (s *Service) Auth() userService.AuthServiceInterface {
    return s.users.Auth
}

func (s *Service) User() userService.UserServiceInterface {
    return s.users.User
}

func (s *Service) Car() carService.CarServiceInterface {
    return s.car.Car
}

func (s *Service) Room() accommodationService.RoomServiceInterface { 
    return s.accommodation.Room
}

func (s *Service) Booking() accommodationService.BookingServiceInterface { 
    return s.accommodation.Booking
}

func (s *Service) Config() *config.Config { 
    return s.config
}

func (s *Service) Marketplace() marketplaceService.MarketplaceServiceInterface {
    return s.marketplace.Marketplace
}

func (s *Service) Review() reviewService.ReviewServiceInterface {
    return s.review.Review
}