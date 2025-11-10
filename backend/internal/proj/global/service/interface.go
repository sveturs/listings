// backend/internal/proj/global/service/interface.go
package service

import (
	"context"

	authService "github.com/sveturs/auth/pkg/service"

	"backend/internal/config"
	balanceService "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	geocodeService "backend/internal/proj/geocode/service"
	notificationService "backend/internal/proj/notifications/service"
	paymentService "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/services"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/interfaces"
)

// SearchLogsServiceInterface интерфейс для логирования поисковых запросов
type SearchLogsServiceInterface interface {
	LogSearch(ctx context.Context, entry interface{}) error
}

type ServicesInterface interface {
	User() userService.UserServiceInterface
	Config() *config.Config
	Review() reviewService.ReviewServiceInterface
	Notification() notificationService.NotificationServiceInterface
	Balance() balanceService.BalanceServiceInterface
	Payment() paymentService.PaymentServiceInterface
	Storage() storage.Storage
	Geocode() geocodeService.GeocodeServiceInterface

	// FileStorage возвращает сервис для работы с файловым хранилищем
	FileStorage() filestorage.FileStorageInterface

	// UnifiedSearch возвращает сервис для унифицированного поиска
	UnifiedSearch() UnifiedSearchServiceInterface

	// BehaviorTracking возвращает сервис для трекинга поведения пользователей
	BehaviorTracking() behaviorTrackingService.BehaviorTrackingService

	// SearchLogs возвращает сервис для логирования поисковых запросов
	SearchLogs() SearchLogsServiceInterface

	// NewImageService создает новый ImageService
	NewImageService(fileStorage filestorage.FileStorageInterface, repo interfaces.ImageRepositoryInterface, cfg services.ImageServiceConfig) *services.ImageService

	// AuthUserService возвращает UserService из auth библиотеки
	AuthUserService() *authService.UserService
}
