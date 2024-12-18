package service

import (
    "backend/internal/storage"
)

type Service struct {
    Car CarServiceInterface
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Car: NewCarService(storage),
    }
}