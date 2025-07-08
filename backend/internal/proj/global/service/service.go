// backend/internal/proj/global/service/service.go
package service

import (
	"log"

	"backend/internal/config"
	balance "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	geocodeService "backend/internal/proj/geocode/service" // Добавить этот импорт
	marketplaceService "backend/internal/proj/marketplace/service"
	translationService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	payment "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefronts/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
)

type Service struct {
	users            *userService.Service
	marketplace      *marketplaceService.Service
	review           *reviewService.Service
	chat             *marketplaceService.Service
	contacts         *marketplaceService.ContactsService
	config           *config.Config
	notification     *notificationService.Service
	translation      translationService.TranslationServiceInterface
	balance          *balance.BalanceService
	payment          payment.PaymentServiceInterface
	storefront       storefrontService.StorefrontService
	storage          storage.Storage
	geocode          geocodeService.GeocodeServiceInterface
	fileStorage      filestorage.FileStorageInterface
	chatAttachment   *marketplaceService.ChatAttachmentService
	unifiedSearch    UnifiedSearchServiceInterface
	behaviorTracking behaviorTrackingService.BehaviorTrackingService
}

func NewService(storage storage.Storage, cfg *config.Config, translationSvc translationService.TranslationServiceInterface) *Service {
	notificationSvc := notificationService.NewService(storage)
	balanceSvc := balance.NewBalanceService(storage)
	geocodeSvc := geocodeService.NewGeocodeService(storage)
	// TODO: behaviorTrackingSvc should be injected from outside since it uses its own repository
	// For now, we'll skip behavior tracking in global service
	var behaviorTrackingSvc behaviorTrackingService.BehaviorTrackingService

	// Создаем сервис витрин (временно без services, передадим позже)
	var storefrontSvc storefrontService.StorefrontService

	// Создаем mock сервис платежей для разработки
	// В продакшене здесь будет AllSecure сервис
	paymentSvc := payment.NewMockPaymentService(cfg.FrontendURL)
	// Create services
	marketplaceSvc := marketplaceService.NewService(storage, notificationSvc.Notification)
	contactsSvc := marketplaceService.NewContactsService(storage)

	// Here we need to set the real translation service to the marketplace service
	if ms, ok := marketplaceSvc.Marketplace.(*marketplaceService.MarketplaceService); ok && translationSvc != nil {
		// Inject the real translation service
		ms.SetTranslationService(translationSvc)
	}

	// Инициализация файлового хранилища
	fileStorageSvc, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		log.Printf("Ошибка инициализации файлового хранилища: %v. Будут использоваться временные файлы.", err)
	}

	// Создаем отдельное хранилище для chat-files
	chatFileStorageConfig := cfg.FileStorage
	chatFileStorageConfig.MinioBucketName = "chat-files"
	chatFileStorageSvc, err := filestorage.NewFileStorage(chatFileStorageConfig)
	if err != nil {
		log.Printf("Ошибка инициализации хранилища чат-файлов: %v", err)
		chatFileStorageSvc = fileStorageSvc // Используем основное хранилище как fallback
	}

	// Инициализация сервиса вложений чата
	chatAttachmentSvc := marketplaceService.NewChatAttachmentService(storage, chatFileStorageSvc, cfg.FileUpload)

	// Установка сервиса вложений в marketplace сервис
	marketplaceSvc.SetChatAttachmentService(chatAttachmentSvc)

	// Создаем экземпляр Service
	s := &Service{
		users:            userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL, cfg.JWTSecret, cfg.JWTExpirationHours),
		marketplace:      marketplaceSvc,
		review:           reviewService.NewService(storage),
		chat:             marketplaceSvc, // Reuse the same service for chat
		contacts:         contactsSvc,
		config:           cfg,
		notification:     notificationSvc,
		translation:      translationSvc,
		balance:          balanceSvc,
		payment:          paymentSvc,
		storefront:       storefrontSvc,
		storage:          storage,
		geocode:          geocodeSvc,
		fileStorage:      fileStorageSvc,
		chatAttachment:   chatAttachmentSvc,
		behaviorTracking: behaviorTrackingSvc,
	}

	// Теперь создаем сервис витрин с правильными зависимостями
	storefrontSvc = storefrontService.NewStorefrontService(s)
	if svc, ok := storefrontSvc.(*storefrontService.StorefrontServiceImpl); ok {
		svc.SetServices(s)
	}
	s.storefront = storefrontSvc

	// Инициализация сервиса унифицированного поиска
	s.unifiedSearch = NewUnifiedSearchService(s)

	return s
}

func (s *Service) Shutdown() {
	log.Println("All services stopped")
}

func (s *Service) Geocode() geocodeService.GeocodeServiceInterface {
	return s.geocode
}

func (s *Service) Storage() storage.Storage {
	return s.storage
}

func (s *Service) Storefront() storefrontService.StorefrontService {
	return s.storefront
}

func (s *Service) Payment() payment.PaymentServiceInterface {
	return s.payment
}

func (s *Service) Balance() balance.BalanceServiceInterface {
	return s.balance
}

// FileStorage возвращает сервис для работы с файловым хранилищем
func (s *Service) FileStorage() filestorage.FileStorageInterface {
	return s.fileStorage
}

// Остальные методы интерфейса ServicesInterface
func (s *Service) Auth() userService.AuthServiceInterface {
	return s.users.Auth
}

func (s *Service) User() userService.UserServiceInterface {
	return s.users.User
}

func (s *Service) Config() *config.Config {
	return s.config
}

func (s *Service) Marketplace() marketplaceService.MarketplaceServiceInterface {
	return s.marketplace.Marketplace
}

func (s *Service) Review() reviewService.ReviewServiceInterface {
	return s.review.Review
}

func (s *Service) Chat() marketplaceService.ChatServiceInterface {
	return s.chat.Chat
}

func (s *Service) Notification() notificationService.NotificationServiceInterface {
	return s.notification.Notification
}

func (s *Service) Translation() translationService.TranslationServiceInterface {
	return s.translation
}

func (s *Service) Contacts() marketplaceService.ContactsServiceInterface {
	return s.contacts
}

func (s *Service) ChatAttachment() marketplaceService.ChatAttachmentServiceInterface {
	return s.chatAttachment
}

func (s *Service) UnifiedSearch() UnifiedSearchServiceInterface {
	return s.unifiedSearch
}

func (s *Service) Orders() marketplaceService.OrderServiceInterface {
	if s.marketplace.Order != nil {
		return s.marketplace.Order
	}
	return nil
}

func (s *Service) BehaviorTracking() behaviorTrackingService.BehaviorTrackingService {
	return s.behaviorTracking
}
