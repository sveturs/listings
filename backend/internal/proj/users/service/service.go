// backend/internal/proj/users/service/service.go
package service

import (
	authService "github.com/sveturs/auth/pkg/http/service"
)

type Service struct {
	User UserServiceInterface
}

func NewService(authSvc *authService.AuthService, userSvc *authService.UserService) *Service {
	return &Service{
		User: NewUserService(authSvc, userSvc),
	}
}
