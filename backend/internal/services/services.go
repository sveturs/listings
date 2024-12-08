package services

import (
    "backend/internal/config"
    "backend/internal/storage"
)

type ServicesInterface interface {
    Auth() AuthServiceInterface
    Room() RoomServiceInterface
    Booking() BookingServiceInterface
    User() UserServiceInterface
    Car() CarServiceInterface
    Config() *config.Config
}

type Services struct {
    auth    AuthServiceInterface
    room    RoomServiceInterface
    booking BookingServiceInterface
    user    UserServiceInterface
    car     CarServiceInterface // Добавляем поле для CarService
    config  *config.Config
}

func NewServices(storage storage.Storage, cfg *config.Config) *Services {
    return &Services{
        auth:    NewAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL, storage),
        room:    NewRoomService(storage),
        booking: NewBookingService(storage),
        user:    NewUserService(storage),
        car:     NewCarService(storage), // Инициализируем CarService
        config:  cfg,
    }
}

func (s *Services) Auth() AuthServiceInterface { 
    return s.auth
}
func (s *Services) Car() CarServiceInterface {
    return s.car
}
func (s *Services) Room() RoomServiceInterface { 
    return s.room
}

func (s *Services) Booking() BookingServiceInterface { 
    return s.booking
}

func (s *Services) User() UserServiceInterface { 
    return s.user
}

func (s *Services) Config() *config.Config { 
    return s.config
}

// Проверяем, что Services реализует ServicesInterface
var _ ServicesInterface = (*Services)(nil)