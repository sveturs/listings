// backend/internal/proj/global/service/service.go
package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	authService "github.com/sveturs/auth/pkg/service"

	"backend/internal/config"
	balance "backend/internal/proj/balance/service"
	behaviorTrackingService "backend/internal/proj/behavior_tracking/service"
	behaviorTrackingPostgres "backend/internal/proj/behavior_tracking/storage/postgres"
	geocodeService "backend/internal/proj/geocode/service"
	notificationService "backend/internal/proj/notifications/service"
	payment "backend/internal/proj/payments/service"
	reviewService "backend/internal/proj/reviews/service"
	userService "backend/internal/proj/users/service"
	"backend/internal/services"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/interfaces"
)

type Service struct {
	users            *userService.Service
	review           *reviewService.Service
	config           *config.Config
	notification     *notificationService.Service
	balance          *balance.BalanceService
	payment          payment.PaymentServiceInterface
	storage          storage.Storage
	geocode          geocodeService.GeocodeServiceInterface
	fileStorage      filestorage.FileStorageInterface
	unifiedSearch    UnifiedSearchServiceInterface
	behaviorTracking behaviorTrackingService.BehaviorTrackingService
	authUserService  *authService.UserService // Auth библиотека UserService
}

func NewService(ctx context.Context, storage storage.Storage, cfg *config.Config, authSvc *authService.AuthService, userSvc *authService.UserService) *Service {
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

	// Создаем mock сервис платежей для разработки
	// В продакшене здесь будет AllSecure сервис
	paymentSvc := payment.NewMockPaymentService(cfg.FrontendURL)

	// Инициализация файлового хранилища
	fileStorageSvc, err := filestorage.NewFileStorage(ctx, cfg.FileStorage)
	if err != nil {
		log.Printf("Ошибка инициализации файлового хранилища: %v. Будут использоваться временные файлы.", err)
	}

	// Создаем userService
	usersSvc := userService.NewService(authSvc, userSvc, storage)

	// Создаем экземпляр Service
	s := &Service{
		users:            usersSvc,
		review:           reviewService.NewService(storage),
		config:           cfg,
		notification:     notificationSvc,
		balance:          balanceSvc,
		payment:          paymentSvc,
		storage:          storage,
		geocode:          geocodeSvc,
		fileStorage:      fileStorageSvc,
		behaviorTracking: behaviorTrackingSvc,
		authUserService:  userSvc, // Сохраняем auth UserService
	}

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

func (s *Service) Review() reviewService.ReviewServiceInterface {
	return s.review.Review
}

func (s *Service) Notification() notificationService.NotificationServiceInterface {
	return s.notification.Notification
}

func (s *Service) UnifiedSearch() UnifiedSearchServiceInterface {
	return s.unifiedSearch
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
