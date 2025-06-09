// backend/internal/proj/notifications/service/service.go
package service

import (
	"backend/internal/storage"
)

type Service struct {
	Notification NotificationServiceInterface
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		Notification: NewNotificationService(storage),
	}
}
