// backend/internal/proj/global/service/service.go
package service

import (
    "backend/internal/config"
    balance "backend/internal/proj/balance/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    translationService "backend/internal/proj/marketplace/service"
    geocodeService "backend/internal/proj/geocode/service" // Добавить этот импорт
    notificationService "backend/internal/proj/notifications/service"
    payment "backend/internal/proj/payments/service"
    reviewService "backend/internal/proj/reviews/service"
    storefrontService "backend/internal/proj/storefront/service"
    userService "backend/internal/proj/users/service"
    "backend/internal/storage"
)

type Service struct {
    users        *userService.Service
    marketplace  *marketplaceService.Service
    review       *reviewService.Service
    chat         *marketplaceService.Service
    config       *config.Config
    notification *notificationService.Service
    translation  translationService.TranslationServiceInterface
    balance      *balance.BalanceService
    payment      payment.PaymentServiceInterface
    storefront   storefrontService.StorefrontServiceInterface
    storage      storage.Storage
    geocode      geocodeService.GeocodeServiceInterface
}

func NewService(storage storage.Storage, cfg *config.Config, translationSvc translationService.TranslationServiceInterface) *Service {
    notificationSvc := notificationService.NewService(storage)
    balanceSvc := balance.NewBalanceService(storage)
    geocodeSvc := geocodeService.NewGeocodeService(storage)

    // Создаем сервис платежей с передачей сервиса баланса
    stripeService := payment.NewStripeService(
        cfg.StripeAPIKey,
        cfg.StripeWebhookSecret,
        cfg.FrontendURL,
        balanceSvc,
    )

    return &Service{
        users:        userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL),
        marketplace:  marketplaceService.NewService(storage, notificationSvc.Notification),
        review:       reviewService.NewService(storage),
        chat:         marketplaceService.NewService(storage, notificationSvc.Notification),
        config:       cfg,
        notification: notificationSvc,
        translation:  translationSvc,
        balance:      balanceSvc,
        payment:      stripeService,
        storefront:   storefrontService.NewStorefrontService(storage),
        storage:      storage,
        geocode:      geocodeSvc,
    }
}

func (s *Service) Geocode() geocodeService.GeocodeServiceInterface {
    return s.geocode

}


func (s *Service) Storage() storage.Storage {
	return s.storage
}
func (s *Service) Storefront() storefrontService.StorefrontServiceInterface {
	return s.storefront
}
func (s *Service) Payment() payment.PaymentServiceInterface {
	return s.payment
}
func (s *Service) Balance() balance.BalanceServiceInterface {
	return s.balance
}

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
func (s *Service) Auto() marketplaceService.AutoServiceInterface {
    return s.marketplace.Auto
}