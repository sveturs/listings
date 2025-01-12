// backend/internal/proj/marketplace/service/service.go
package service

import (
    "backend/internal/storage"
)

type Service struct {
    Marketplace MarketplaceServiceInterface
    Chat       ChatServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Marketplace: NewMarketplaceService(storage),
        Chat:       NewChatService(storage),
    }
}