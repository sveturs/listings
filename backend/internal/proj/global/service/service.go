package service

import (
    "backend/internal/config"
    "backend/internal/storage"
    userService "backend/internal/proj/users/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
)

type Service struct {
    users         *userService.Service
    marketplace  *marketplaceService.Service
    review       *reviewService.Service
    chat         *marketplaceService.Service
    config       *config.Config
}

func NewService(storage storage.Storage, cfg *config.Config) *Service {
    return &Service{
        users:         userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
        marketplace:  marketplaceService.NewService(storage),
        review:       reviewService.NewService(storage),
        chat:         marketplaceService.NewService(storage),
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

func (s *Service) Config() *config.Config {
    return s.config
}

func (s *Service) Marketplace() marketplaceService.MarketplaceServiceInterface {
    return s.marketplace.Marketplace
}

func (s *Service) Review() reviewService.ReviewServiceInterface {
    return s.review.Review
}

func (s *Service) Chat() marketplaceService.ChatServiceInterface {
    return s.chat.Chat
}
