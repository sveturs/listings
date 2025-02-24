// backend/internal/proj/global/service/service.go

package service

import (
    "backend/internal/config"
    "backend/internal/storage"
    userService "backend/internal/proj/users/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
    notificationService "backend/internal/proj/notifications/service"
    translationService "backend/internal/proj/marketplace/service"
    balance "backend/internal/proj/balance/service"
)

type Service struct {
    users         *userService.Service
    marketplace   *marketplaceService.Service
    review        *reviewService.Service
    chat          *marketplaceService.Service
    config        *config.Config
    notification  *notificationService.Service
    translation   translationService.TranslationServiceInterface
    balance       *balance.BalanceService
}

func NewService(storage storage.Storage, cfg *config.Config, translationSvc translationService.TranslationServiceInterface) *Service {
    notificationSvc := notificationService.NewService(storage)

    return &Service{
        users:        userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
        marketplace:  marketplaceService.NewService(storage, notificationSvc.Notification),
        review:       reviewService.NewService(storage),
        chat:         marketplaceService.NewService(storage, notificationSvc.Notification),
        config:       cfg,
        notification: notificationSvc,
        translation:  translationSvc,
        balance:      balance.NewBalanceService(storage),
    }
}
func (s *Service) Balance() balance.BalanceServiceInterface {
    return s.balance
}

// Остальные методы интерфейса ServicesInterface
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

func (s *Service) Notification() notificationService.NotificationServiceInterface {
    return s.notification.Notification
}

func (s *Service) Translation() translationService.TranslationServiceInterface {
    return s.translation
}