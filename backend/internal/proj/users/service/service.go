package service

import (
    "backend/internal/storage"
    
)

type Service struct {
    Auth AuthServiceInterface
    User UserServiceInterface
}

func NewService(storage storage.Storage, googleClientID, googleClientSecret, googleRedirectURL string) *Service {
    return &Service{
        Auth: NewAuthService(googleClientID, googleClientSecret, googleRedirectURL, storage),
        User: NewUserService(storage),
    }
}