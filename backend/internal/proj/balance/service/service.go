// backend/internal/proj/balance/service/service.go

package balance

import (
    "backend/internal/storage"
)

type Service struct {
    Balance BalanceServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Balance: NewBalanceService(storage),
    }
}