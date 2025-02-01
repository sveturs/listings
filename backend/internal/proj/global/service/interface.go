// backend/internal/proj/global/service/interface.go
package service

import (
    userService "backend/internal/proj/users/service"
    marketplaceService "backend/internal/proj/marketplace/service"
    reviewService "backend/internal/proj/reviews/service"
    chatService "backend/internal/proj/marketplace/service"
    "backend/internal/config"
    notificationService "backend/internal/proj/notifications/service"
    translationService "backend/internal/proj/marketplace/service" 
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
}