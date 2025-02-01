// backend/internal/proj/global/service/service.go
package service

import (
    "backend/internal/config"
    "backend/internal/storage"
    userService "backend/internal/proj/users/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
    notificationService "backend/internal/proj/notifications/service"
    translationService "backend/internal/proj/marketplace/service"  // правильный импорт
    //"log"
)

type Service struct {
    users         *userService.Service
    marketplace   *marketplaceService.Service
    review        *reviewService.Service
    chat          *marketplaceService.Service
    config        *config.Config
    notification  *notificationService.Service
     translation   translationService.TranslationServiceInterface
}

func NewService(storage storage.Storage, cfg *config.Config, translationSvc translationService.TranslationServiceInterface) *Service {
    notificationSvc := notificationService.NewService(storage)
    
    return &Service{
        users:        userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
        marketplace: marketplaceService.NewService(storage, notificationSvc.Notification),
        review:      reviewService.NewService(storage),
        chat:        marketplaceService.NewService(storage, notificationSvc.Notification),
        config:      cfg,
        notification: notificationSvc,
        translation:  translationSvc, // используем переданный сервис
    }
}

// Implement interface methods
func (s *Service) Auth() userService.AuthServiceInterface {
    return s.users.Auth
}

func (s *Service) Notification() notificationService.NotificationServiceInterface {
    return s.notification.Notification
}

func (s *Service) User() userService.UserServiceInterface {
    return s.users.User
}

func (s *Service) Translation() translationService.TranslationServiceInterface {
    return s.translation
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