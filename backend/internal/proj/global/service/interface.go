// backend/internal/proj/global/service/interface.go
package service

import (
	"backend/internal/config"
	balanceService "backend/internal/proj/balance/service"
	chatService "backend/internal/proj/marketplace/service"
	marketplaceService "backend/internal/proj/marketplace/service"
	translationService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	paymentService "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefront/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/storage" 
)

type ServicesInterface interface {
	Auth() userService.AuthServiceInterface
	User() userService.UserServiceInterface
	Config() *config.Config
	Marketplace() marketplaceService.MarketplaceServiceInterface
	Review() reviewService.ReviewServiceInterface
	Chat() chatService.ChatServiceInterface
	Notification() notificationService.NotificationServiceInterface
	Translation() translationService.TranslationServiceInterface
	Balance() balanceService.BalanceServiceInterface
	Payment() paymentService.PaymentServiceInterface
	Storefront() storefrontService.StorefrontServiceInterface
	Storage() storage.Storage 
}
