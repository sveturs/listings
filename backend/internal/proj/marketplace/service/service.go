// backend/internal/proj/marketplace/service/service.go
package service

import (
    "backend/internal/storage"
)

type Service struct {
    Marketplace MarketplaceServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Marketplace: NewMarketplaceService(storage),
    }
}