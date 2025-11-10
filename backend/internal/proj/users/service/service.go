// backend/internal/proj/users/service/service.go
package service

import (
	authService "github.com/sveturs/auth/pkg/service"

	"backend/internal/storage"
)

type Service struct {
	User UserServiceInterface
}

func NewService(authSvc *authService.AuthService, userSvc *authService.UserService, storage storage.Storage) *Service {
	return &Service{
		User: NewUserService(authSvc, userSvc, storage),
	}
}
