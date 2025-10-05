// backend/internal/proj/global/service/service.go
package service

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	authService "github.com/sveturs/auth/pkg/http/service"

	"backend/internal/cache"
	"backend/internal/config"
	balance "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	behaviorTrackingPostgres "backend/internal/proj/behavior_tracking/storage/postgres"
	geocodeService "backend/internal/proj/geocode/service" // Добавить этот импорт
	marketplaceService "backend/internal/proj/marketplace/service"
	notificationService "backend/internal/proj/notifications/service"
	payment "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	storefrontService "backend/internal/proj/storefronts/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/services"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/interfaces"
)

type Service struct {
	users            *userService.Service
	marketplace      *marketplaceService.Service
	review           *reviewService.Service
	chat             *marketplaceService.Service
	contacts         *marketplaceService.ContactsService
	config           *config.Config
	notification     *notificationService.Service
	translation      marketplaceService.TranslationServiceInterface
	balance          *balance.BalanceService
	payment          payment.PaymentServiceInterface
	storefront       storefrontService.StorefrontService
	storage          storage.Storage
	geocode          geocodeService.GeocodeServiceInterface
	fileStorage      filestorage.FileStorageInterface
	chatAttachment   *marketplaceService.ChatAttachmentService
	chatTranslation  *marketplaceService.ChatTranslationService
	unifiedSearch    UnifiedSearchServiceInterface
	behaviorTracking behaviorTrackingService.BehaviorTrackingService
	unifiedCar       *marketplaceService.UnifiedCarService
	authUserService  *authService.UserService // Auth библиотека UserService
}

func NewService(ctx context.Context, storage storage.Storage, cfg *config.Config, translationSvc marketplaceService.TranslationServiceInterface, authSvc *authService.AuthService, userSvc *authService.UserService) *Service {
	notificationSvc := notificationService.NewService(storage)
	balanceSvc := balance.NewBalanceService(storage)
	geocodeSvc := geocodeService.NewGeocodeService(storage)

	// Initialize behavior tracking service
	var behaviorTrackingSvc behaviorTrackingService.BehaviorTrackingService
	// Get pool from storage if it's a postgres.Database
	if poolAccessor, ok := storage.(interface{ GetPool() *pgxpool.Pool }); ok && poolAccessor != nil {
		if pool := poolAccessor.GetPool(); pool != nil {
			// Create behavior tracking repository and service
			behaviorRepo := behaviorTrackingPostgres.NewBehaviorTrackingRepository(pool)
			behaviorTrackingSvc = behaviorTrackingService.NewBehaviorTrackingService(ctx, behaviorRepo)
			log.Println("Behavior tracking service initialized successfully")
		} else {
			log.Println("Warning: PostgreSQL pool not available for behavior tracking")
		}
	} else {
		log.Println("Warning: Storage does not support GetPool() method for behavior tracking")
	}

	// Создаем сервис витрин (временно без services, передадим позже)
	var storefrontSvc storefrontService.StorefrontService

	// Создаем Redis кеш если настроен
	var cacheAdapter marketplaceService.CacheInterface
	var redisClient *redis.Client
	if cfg.Redis.URL != "" {
		logger := logrus.New()
		redisCache, err := cache.NewRedisCache(
			ctx,
			cfg.Redis.URL,
			cfg.Redis.Password,
			cfg.Redis.DB,
			cfg.Redis.PoolSize,
			logger,
		)
		if err != nil {
			log.Printf("Warning: Failed to initialize Redis cache: %v", err)
			// Продолжаем работу без кеша
		} else {
			cacheAdapter = cache.NewAdapter(redisCache)
			redisClient = redisCache.GetClient()
			log.Println("Redis cache initialized successfully")
		}
	} else {
		log.Println("Redis not configured, running without cache")
	}

	// Создаем mock сервис платежей для разработки
	// В продакшене здесь будет AllSecure сервис
	paymentSvc := payment.NewMockPaymentService(cfg.FrontendURL)
	// Create services
	marketplaceSvc := marketplaceService.NewService(storage, notificationSvc.Notification, cfg.SearchWeights, cacheAdapter)
	contactsSvc := marketplaceService.NewContactsService(storage)

	// Here we need to set the real translation service to the marketplace service
	if ms, ok := marketplaceSvc.Marketplace.(*marketplaceService.MarketplaceService); ok && translationSvc != nil {
		// Inject the real translation service
		ms.SetTranslationService(translationSvc)
	}

	// Инициализация файлового хранилища
	fileStorageSvc, err := filestorage.NewFileStorage(ctx, cfg.FileStorage)
	if err != nil {
		log.Printf("Ошибка инициализации файлового хранилища: %v. Будут использоваться временные файлы.", err)
	}

	// Создаем отдельное хранилище для chat-files
	chatFileStorageConfig := cfg.FileStorage
	chatFileStorageConfig.MinioBucketName = cfg.FileStorage.MinioChatBucket
	chatFileStorageSvc, err := filestorage.NewFileStorage(ctx, chatFileStorageConfig)
	if err != nil {
		log.Printf("Ошибка инициализации хранилища чат-файлов: %v", err)
		chatFileStorageSvc = fileStorageSvc // Используем основное хранилище как fallback
	}

	// Инициализация сервиса вложений чата
	chatAttachmentSvc := marketplaceService.NewChatAttachmentService(storage, chatFileStorageSvc, cfg.FileUpload)

	// Установка сервиса вложений в marketplace сервис
	marketplaceSvc.SetChatAttachmentService(chatAttachmentSvc)

	// Создаем userService для chatTranslation (с доступом к storage для chat settings)
	usersSvc := userService.NewService(authSvc, userSvc, storage)

	// Инициализация сервиса переводов чата (теперь с доступом к userService и storage)
	chatTranslationSvc := marketplaceService.NewChatTranslationService(translationSvc, redisClient, usersSvc.User, storage)

	// Установка сервиса переводов в marketplace сервис
	marketplaceSvc.SetChatTranslationService(chatTranslationSvc)

	// Установка зависимостей в ChatService для поддержки персонализированных переводов в WebSocket
	if chatSvc, ok := marketplaceSvc.Chat.(*marketplaceService.ChatService); ok {
		chatSvc.SetChatTranslationService(chatTranslationSvc)
		chatSvc.SetUserService(usersSvc.User)
		log.Println("ChatService dependencies set (translation & user service)")
	}

	// Создаем UnifiedCarService
	carServiceConfig := &marketplaceService.CarServiceConfig{
		CacheTTL:          24 * time.Hour,
		CacheEnabled:      redisClient != nil,
		VINDecoderEnabled: true, // Включаем VIN декодер
	}
	unifiedCarSvc := marketplaceService.NewUnifiedCarService(storage, redisClient, carServiceConfig)

	// Создаем экземпляр Service
	s := &Service{
		users:            usersSvc,
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
		chatTranslation:  chatTranslationSvc,
		behaviorTracking: behaviorTrackingSvc,
		unifiedCar:       unifiedCarSvc,
		authUserService:  userSvc, // Сохраняем auth UserService
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

// UnifiedCar возвращает сервис для работы с автомобилями
func (s *Service) UnifiedCar() *marketplaceService.UnifiedCarService {
	return s.unifiedCar
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

func (s *Service) Translation() marketplaceService.TranslationServiceInterface {
	return s.translation
}

func (s *Service) Contacts() marketplaceService.ContactsServiceInterface {
	return s.contacts
}

func (s *Service) ChatAttachment() marketplaceService.ChatAttachmentServiceInterface {
	return s.chatAttachment
}

func (s *Service) ChatTranslation() *marketplaceService.ChatTranslationService {
	return s.chatTranslation
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

func (s *Service) SearchLogs() SearchLogsServiceInterface {
	// Временно возвращаем nil, пока не реализован сервис логирования поиска
	return nil
}

// NewImageService создает новый ImageService
func (s *Service) NewImageService(fileStorage filestorage.FileStorageInterface, repo interfaces.ImageRepositoryInterface, cfg services.ImageServiceConfig) *services.ImageService {
	return services.NewImageService(fileStorage, repo, cfg)
}

// AuthUserService возвращает UserService из auth библиотеки
func (s *Service) AuthUserService() *authService.UserService {
	return s.authUserService
}
