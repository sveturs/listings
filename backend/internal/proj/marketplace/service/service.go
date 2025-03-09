// backend/internal/proj/marketplace/service/service.go
package service

import (
    "backend/internal/storage"
)
import (
    "backend/internal/proj/notifications/service"
)

type Service struct {
    Marketplace MarketplaceServiceInterface
    Chat       ChatServiceInterface
    Auto       AutoServiceInterface 
}

func NewService(storage storage.Storage, notifService service.NotificationServiceInterface) *Service {
    return &Service{
        Marketplace: NewMarketplaceService(storage),
        Chat:       NewChatService(storage, notifService),
        Auto:       NewAutoService(storage), 
    }
}