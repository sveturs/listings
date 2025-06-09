// backend/internal/proj/global/service/service.go
package service

import (
	"backend/internal/config"
	balance "backend/internal/proj/balance/service"
	geocodeService "backend/internal/proj/geocode/service" // Добавить этот импорт
	marketplaceService "backend/internal/proj/marketplace/service"
	translationService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	payment "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefront/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"log"
)

type Service struct {
	users          *userService.Service
	marketplace    *marketplaceService.Service
	review         *reviewService.Service
	chat           *marketplaceService.Service
	contacts       *marketplaceService.ContactsService
	config         *config.Config
	notification   *notificationService.Service
	translation    translationService.TranslationServiceInterface
	balance        *balance.BalanceService
	payment        payment.PaymentServiceInterface
	storefront     storefrontService.StorefrontServiceInterface
	storage        storage.Storage
	geocode        geocodeService.GeocodeServiceInterface
	scheduleImport *storefrontService.ScheduleService
	fileStorage    filestorage.FileStorageInterface
	chatAttachment *marketplaceService.ChatAttachmentService
}

func NewService(storage storage.Storage, cfg *config.Config, translationSvc translationService.TranslationServiceInterface) *Service {
	notificationSvc := notificationService.NewService(storage)
	balanceSvc := balance.NewBalanceService(storage)
	geocodeSvc := geocodeService.NewGeocodeService(storage)
	storefrontSvc := storefrontService.NewStorefrontService(storage)
	scheduleService := storefrontService.NewScheduleService(storage, storefrontSvc)

	// Создаем сервис платежей с передачей сервиса баланса
	stripeService := payment.NewStripeService(
		cfg.StripeAPIKey,
		cfg.StripeWebhookSecret,
		cfg.FrontendURL,
		balanceSvc,
	)
	scheduleService.Start()
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

	return &Service{
		users:          userService.NewService(storage, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL, cfg.JWTSecret, cfg.JWTExpirationHours),
		marketplace:    marketplaceSvc,
		review:         reviewService.NewService(storage),
		chat:           marketplaceSvc, // Reuse the same service for chat
		contacts:       contactsSvc,
		config:         cfg,
		notification:   notificationSvc,
		translation:    translationSvc,
		balance:        balanceSvc,
		payment:        stripeService,
		storefront:     storefrontService.NewStorefrontService(storage),
		storage:        storage,
		geocode:        geocodeSvc,
		scheduleImport: scheduleService,
		fileStorage:    fileStorageSvc,
		chatAttachment: chatAttachmentSvc,
	}
}

func (s *Service) Shutdown() {
	// Останавливаем сервис расписания
	if s.scheduleImport != nil {
		s.scheduleImport.Stop()
	}

	log.Println("All services stopped")
}

func (s *Service) Geocode() geocodeService.GeocodeServiceInterface {
	return s.geocode
}

func (s *Service) Storage() storage.Storage {
	return s.storage
}

func (s *Service) Storefront() storefrontService.StorefrontServiceInterface {
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
