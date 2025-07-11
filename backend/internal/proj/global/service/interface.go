// backend/internal/proj/global/service/interface.go
package service

import (
	"context"

	"backend/internal/config"
	balanceService "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	geocodeService "backend/internal/proj/geocode/service"
	chatService "backend/internal/proj/marketplace/service"
	contactsService "backend/internal/proj/marketplace/service"
	marketplaceService "backend/internal/proj/marketplace/service"
	translationService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	paymentService "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefronts/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
)

// SearchLogsServiceInterface интерфейс для логирования поисковых запросов
type SearchLogsServiceInterface interface {
	LogSearch(ctx context.Context, entry interface{}) error
}

type ServicesInterface interface {
	Auth() userService.AuthServiceInterface
	User() userService.UserServiceInterface
	Config() *config.Config
	Marketplace() marketplaceService.MarketplaceServiceInterface
	Review() reviewService.ReviewServiceInterface
	Chat() chatService.ChatServiceInterface
	Contacts() contactsService.ContactsServiceInterface
	Notification() notificationService.NotificationServiceInterface
	Translation() translationService.TranslationServiceInterface
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
}
