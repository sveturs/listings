package service

import (
    "backend/internal/storage"
)

type Service struct {
    Review ReviewServiceInterface // Меняем тип на интерфейс
}

func NewService(storage storage.Storage) *Service {
    return &Service{
        Review: NewReviewService(storage),
    }
}