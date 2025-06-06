// backend/internal/proj/users/service/service.go
package service

import (
    "backend/internal/storage"
)

type Service struct {
    Auth AuthServiceInterface
    User UserServiceInterface
}

func NewService(store storage.Storage, googleClientID, googleClientSecret, googleRedirectURL, jwtSecret string, jwtExpHours int) *Service {
    return &Service{
        Auth: NewAuthService(googleClientID, googleClientSecret, googleRedirectURL, store, jwtSecret, jwtExpHours),
        User: NewUserService(store),
    }
}