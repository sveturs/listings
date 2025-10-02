// backend/internal/proj/global/service/interface.go
package service

import (
	"context"

	"backend/internal/config"
	balanceService "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	geocodeService "backend/internal/proj/geocode/service"
	marketplaceService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	paymentService "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefronts/service"
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
	Marketplace() marketplaceService.MarketplaceServiceInterface
	Review() reviewService.ReviewServiceInterface
	Chat() marketplaceService.ChatServiceInterface
	Contacts() marketplaceService.ContactsServiceInterface
	Notification() notificationService.NotificationServiceInterface
	Translation() marketplaceService.TranslationServiceInterface
	Balance() balanceService.BalanceServiceInterface
	Payment() paymentService.PaymentServiceInterface
	Storefront() storefrontService.StorefrontService
	Storage() storage.Storage
	Geocode() geocodeService.GeocodeServiceInterface

	// FileStorage возвращает сервис для работы с файловым хранилищем
	FileStorage() filestorage.FileStorageInterface

	// ChatAttachment возвращает сервис для работы с вложениями чата
	ChatAttachment() marketplaceService.ChatAttachmentServiceInterface

	// UnifiedSearch возвращает сервис для унифицированного поиска
	UnifiedSearch() UnifiedSearchServiceInterface

	// Orders возвращает сервис для работы с заказами маркетплейса
	Orders() marketplaceService.OrderServiceInterface

	// BehaviorTracking возвращает сервис для трекинга поведения пользователей
	BehaviorTracking() behaviorTrackingService.BehaviorTrackingService

	// SearchLogs возвращает сервис для логирования поисковых запросов
	SearchLogs() SearchLogsServiceInterface

	// UnifiedCar возвращает сервис для работы с автомобилями
	UnifiedCar() *marketplaceService.UnifiedCarService

	// NewImageService создает новый ImageService
	NewImageService(fileStorage filestorage.FileStorageInterface, repo interfaces.ImageRepositoryInterface, cfg services.ImageServiceConfig) *services.ImageService

	// AuthUserService возвращает UserService из auth библиотеки
	AuthUserService() interface{}
}
