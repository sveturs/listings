// backend/internal/proj/users/service/service.go
package service

import (
	"backend/internal/storage"
)

type Service struct {
	User UserServiceInterface
}

func NewService(store storage.Storage) *Service {
	return &Service{
		User: NewUserService(store),
	}
}
