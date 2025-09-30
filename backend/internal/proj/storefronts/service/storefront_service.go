package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/notifications/service"
	"backend/internal/proj/storefronts/common"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/postgres"
)

// Ошибки
var (
	ErrUnauthorized             = errors.New("unauthorized")
	ErrInsufficientPermissions  = errors.New("insufficient permissions")
	ErrStorefrontLimitReached   = errors.New("storefront limit reached for current plan")
	ErrSlugAlreadyExists        = errors.New("slug already exists")
	ErrInvalidLocation          = errors.New("invalid location data")
	ErrFeatureNotAvailable      = errors.New("feature not available in current plan")
	ErrRepositoryNotInitialized = errors.New("storefront repository not initialized")
	ErrStaffLimitReached        = errors.New("staff limit reached for current plan")
	ErrStorefrontNotFound       = errors.New("storefront not found")
)

// Лимиты по тарифным планам
var planLimits = map[models.SubscriptionPlanType]struct {
	MaxStorefronts  int
	MaxProducts     int
	MaxStaff        int
	MaxImages       int
	CanUseAI        bool
	CanUseLive      bool
	CanUseGroup     bool
	CanExportData   bool
	CanCustomDomain bool
}{
	models.SubscriptionPlanStarter: {
		MaxStorefronts: 1,
		MaxProducts:    50,
		MaxStaff:       1,
		MaxImages:      100,
	},
	models.SubscriptionPlanProfessional: {
		MaxStorefronts:  3,
		MaxProducts:     500,
		MaxStaff:        5,
		MaxImages:       1000,
		CanExportData:   true,
		CanCustomDomain: true,
	},
	models.SubscriptionPlanBusiness: {
		MaxStorefronts:  10,
		MaxProducts:     5000,
		MaxStaff:        20,
		MaxImages:       10000,
		CanUseAI:        true,
		CanUseLive:      true,
		CanExportData:   true,
		CanCustomDomain: true,
	},
	models.SubscriptionPlanEnterprise: {
		MaxStorefronts:  -1, // unlimited
		MaxProducts:     -1,
		MaxStaff:        -1,
		MaxImages:       -1,
		CanUseAI:        true,
		CanUseLive:      true,
		CanUseGroup:     true,
		CanExportData:   true,
		CanCustomDomain: true,
	},
}

// StorefrontService интерфейс сервиса витрин
type StorefrontService interface {
	// Основные операции
	CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error)
	GetByID(ctx context.Context, id int) (*models.Storefront, error)
	GetBySlug(ctx context.Context, slug string) (*models.Storefront, error)
	Update(ctx context.Context, userID int, storefrontID int, dto *models.StorefrontUpdateDTO) error
	Delete(ctx context.Context, userID int, storefrontID int) error
	Restore(ctx context.Context, userID int, storefrontID int) error

	// Листинг и поиск
	ListUserStorefronts(ctx context.Context, userID int) ([]*models.Storefront, error)
	Search(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error)
	SearchOpenSearch(ctx context.Context, params *opensearch.StorefrontSearchParams) (*opensearch.StorefrontSearchResult, error)

	// Картографические функции
	GetMapData(ctx context.Context, bounds postgres.GeoBounds, filter *models.StorefrontFilter) ([]*models.StorefrontMapData, error)
	GetNearby(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*models.Storefront, error)
	GetBusinessesInBuilding(ctx context.Context, lat, lng float64) ([]*models.StorefrontMapData, error)

	// Управление настройками
	UpdateWorkingHours(ctx context.Context, userID int, storefrontID int, hours []*models.StorefrontHours) error
	UpdatePaymentMethods(ctx context.Context, userID int, storefrontID int, methods []*models.StorefrontPaymentMethod) error
	UpdateDeliveryOptions(ctx context.Context, userID int, storefrontID int, options []*models.StorefrontDeliveryOption) error

	// Управление персоналом
	AddStaff(ctx context.Context, ownerID int, storefrontID int, userID int, role models.StaffRole) error
	UpdateStaffPermissions(ctx context.Context, ownerID int, staffID int, permissions models.JSONB) error
	RemoveStaff(ctx context.Context, ownerID int, storefrontID int, userID int) error
	GetStaff(ctx context.Context, storefrontID int) ([]*models.StorefrontStaff, error)

	// Загрузка изображений
	UploadLogo(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error)
	UploadBanner(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error)

	// Аналитика
	RecordView(ctx context.Context, storefrontID int) error
	GetAnalytics(ctx context.Context, userID int, storefrontID int, from, to time.Time) ([]*models.StorefrontAnalytics, error)

	// Проверки
	CanCreateStorefront(ctx context.Context, userID int) (bool, error)
	CheckFeatureAvailability(ctx context.Context, storefrontID int, feature string) (bool, error)
}

// NotificationService интерфейс для уведомлений
type NotificationService interface {
	CreateNotification(ctx context.Context, notification *models.Notification) error
}

// ServicesInterface интерфейс для доступа к сервисам
type ServicesInterface interface {
	Storage() storage.Storage
	FileStorage() filestorage.FileStorageInterface
	Notification() service.NotificationServiceInterface
}

// StorefrontServiceImpl реализация сервиса
type StorefrontServiceImpl struct {
	services ServicesInterface
	repo     postgres.StorefrontRepository
}

// NewStorefrontService создает новый сервис витрин
func NewStorefrontService(services ServicesInterface) StorefrontService {
	service := &StorefrontServiceImpl{
		services: services,
	}

	// Инициализируем репозиторий сразу
	if services != nil && services.Storage() != nil {
		repoInterface := services.Storage().Storefront()
		if storefrontRepo, ok := repoInterface.(postgres.StorefrontRepository); ok {
			service.repo = storefrontRepo
			logger.Info().Msg("Storefront repository initialized in NewStorefrontService")
		} else {
			logger.Error().Str("type", fmt.Sprintf("%T", repoInterface)).Msg("Failed to cast storefront repository in NewStorefrontService")
		}
	}

	return service
}

// SetServices устанавливает ссылку на services после инициализации
func (s *StorefrontServiceImpl) SetServices(services ServicesInterface) {
	s.services = services
	if services != nil && services.Storage() != nil {
		// Получаем интерфейс репозитория
		repoInterface := services.Storage().Storefront()
		if storefrontRepo, ok := repoInterface.(postgres.StorefrontRepository); ok {
			s.repo = storefrontRepo
			logger.Info().Msg("Storefront repository initialized successfully")
		} else {
			logger.Error().Str("type", fmt.Sprintf("%T", repoInterface)).Msg("Failed to cast storefront repository")
		}
	}
}

// CreateStorefront создает новую витрину
func (s *StorefrontServiceImpl) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	logger.Info().Int("userID", userID).Str("name", dto.Name).Msg("CreateStorefront: начало создания витрины")

	// Проверяем инициализацию репозитория
	if s.repo == nil {
		logger.Error().Msg("CreateStorefront: репозиторий не инициализирован")
		return nil, ErrRepositoryNotInitialized
	}

	// Проверяем может ли пользователь создать витрину
	canCreate, err := s.canCreateStorefront(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("CreateStorefront: ошибка проверки прав")
		return nil, fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !canCreate {
		logger.Warn().Int("userID", userID).Msg("CreateStorefront: превышен лимит витрин")
		return nil, ErrStorefrontLimitReached
	}

	// Дополняем DTO
	dto.Slug = s.generateSlug(dto.Name)
	dto.UserID = userID
	logger.Info().Str("slug", dto.Slug).Msg("CreateStorefront: сгенерирован slug")

	// Создаем витрину
	storefront, err := s.repo.Create(ctx, dto)
	if err != nil {
		logger.Error().Err(err).Str("slug", dto.Slug).Msg("CreateStorefront: ошибка создания витрины в БД")
		return nil, fmt.Errorf("ошибка создания витрины: %w", err)
	}

	// Индексируем в OpenSearch
	if err := s.services.Storage().IndexStorefront(ctx, storefront); err != nil {
		logger.Error().Err(err).Msg("Ошибка индексации витрины в OpenSearch")
	}

	// Создаем уведомление о создании витрины
	notification := &models.Notification{
		UserID:    userID,
		Type:      "storefront_created",
		Title:     "Витрина создана",
		Message:   fmt.Sprintf("Витрина '%s' успешно создана", storefront.Name),
		ListingID: storefront.ID, // Используем ListingID как общий EntityID
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.services.Notification().CreateNotification(ctx, notification); err != nil {
		// Не прерываем создание витрины из-за ошибки уведомления
		logger.Error().Err(err).Msg("Ошибка создания уведомления о создании витрины")
	}

	return storefront, nil
}

// GetByID получает витрину по ID
func (s *StorefrontServiceImpl) GetByID(ctx context.Context, id int) (*models.Storefront, error) {
	return s.repo.GetByID(ctx, id)
}

// GetBySlug получает витрину по slug
func (s *StorefrontServiceImpl) GetBySlug(ctx context.Context, slug string) (*models.Storefront, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// Update обновляет витрину
func (s *StorefrontServiceImpl) Update(ctx context.Context, userID int, storefrontID int, dto *models.StorefrontUpdateDTO) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Обновляем витрину
	if err := s.repo.Update(ctx, storefrontID, dto); err != nil {
		return fmt.Errorf("ошибка обновления витрины: %w", err)
	}

	// Получаем обновленную витрину для переиндексации
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка получения витрины после обновления")
		return nil
	}

	// Переиндексируем в OpenSearch
	if err := s.services.Storage().IndexStorefront(ctx, storefront); err != nil {
		logger.Error().Err(err).Msg("Ошибка переиндексации витрины в OpenSearch")
	}

	return nil
}

// Delete удаляет витрину
func (s *StorefrontServiceImpl) Delete(ctx context.Context, userID int, storefrontID int) error {
	// Проверяем, является ли пользователь администратором из контекста
	isAdmin, ok := ctx.Value(common.ContextKeyIsAdmin).(bool)
	logger.Info().
		Int("userID", userID).
		Int("storefrontID", storefrontID).
		Bool("isAdmin", isAdmin).
		Bool("isAdminOk", ok).
		Msg("Delete service called")

	if ok && isAdmin {
		// Проверяем флаг жесткого удаления из контекста
		hardDelete, _ := ctx.Value(common.ContextKeyHardDelete).(bool)

		if hardDelete {
			// Администратор выбрал жесткое удаление
			logger.Info().Msgf("Admin user %d hard deleting storefront %d", userID, storefrontID)

			// Жестко удаляем витрину (полное удаление из БД)
			if err := s.repo.HardDeleteDebug(ctx, storefrontID); err != nil {
				logger.Error().Err(err).Int("storefrontID", storefrontID).Msg("Failed to hard delete storefront")
				return fmt.Errorf("ошибка жесткого удаления витрины: %w", err)
			}
		} else {
			// Администратор выбрал мягкое удаление
			logger.Info().Msgf("Admin user %d soft deleting storefront %d", userID, storefrontID)

			// Мягкое удаление (деактивация)
			if err := s.repo.Delete(ctx, storefrontID); err != nil {
				return fmt.Errorf("ошибка мягкого удаления витрины: %w", err)
			}
		}

		// Удаляем из OpenSearch
		if err := s.services.Storage().DeleteStorefrontIndex(ctx, storefrontID); err != nil {
			logger.Error().Err(err).Msg("Ошибка удаления витрины из OpenSearch")
		}

		return nil
	}

	// Для обычных пользователей проверяем является ли пользователь владельцем
	isOwner, err := s.repo.IsOwner(ctx, storefrontID, userID)
	if err != nil {
		return fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	if !isOwner {
		return ErrUnauthorized
	}

	// Удаляем витрину
	if err := s.repo.Delete(ctx, storefrontID); err != nil {
		return fmt.Errorf("ошибка удаления витрины: %w", err)
	}

	// Удаляем из OpenSearch
	if err := s.services.Storage().DeleteStorefrontIndex(ctx, storefrontID); err != nil {
		logger.Error().Err(err).Msg("Ошибка удаления витрины из OpenSearch")
	}

	return nil
}

// Restore восстанавливает деактивированную витрину
func (s *StorefrontServiceImpl) Restore(ctx context.Context, userID int, storefrontID int) error {
	// Проверяем, является ли пользователь администратором из контекста
	isAdmin, ok := ctx.Value(common.ContextKeyIsAdmin).(bool)
	logger.Info().
		Int("userID", userID).
		Int("storefrontID", storefrontID).
		Bool("isAdmin", isAdmin).
		Bool("isAdminOk", ok).
		Msg("Restore service called")

	if !ok || !isAdmin {
		return fmt.Errorf("только администраторы могут восстанавливать витрины")
	}

	// Восстанавливаем витрину (активируем)
	if err := s.repo.Restore(ctx, storefrontID); err != nil {
		logger.Error().Err(err).Int("storefrontID", storefrontID).Msg("Failed to restore storefront")
		return fmt.Errorf("ошибка восстановления витрины: %w", err)
	}

	// Обновляем индекс в OpenSearch
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		logger.Error().Err(err).Int("storefrontID", storefrontID).Msg("Failed to get storefront after restore")
		return fmt.Errorf("ошибка получения витрины после восстановления: %w", err)
	}

	// Индексируем в OpenSearch
	if err := s.services.Storage().IndexStorefront(ctx, storefront); err != nil {
		logger.Error().Err(err).Msg("Ошибка индексации витрины в OpenSearch после восстановления")
		// Не возвращаем ошибку, так как витрина уже восстановлена в БД
	}

	logger.Info().Int("storefrontID", storefrontID).Msg("Storefront restored successfully")
	return nil
}

// ListUserStorefronts получает список витрин пользователя
func (s *StorefrontServiceImpl) ListUserStorefronts(ctx context.Context, userID int) ([]*models.Storefront, error) {
	if s.repo == nil {
		return nil, ErrRepositoryNotInitialized
	}

	filter := &models.StorefrontFilter{
		UserID: &userID,
		Limit:  100,
		Offset: 0,
	}
	storefronts, _, err := s.repo.List(ctx, filter)
	return storefronts, err
}

// Search выполняет поиск витрин
func (s *StorefrontServiceImpl) Search(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error) {
	return s.repo.List(ctx, filter)
}

// SearchOpenSearch выполняет поиск витрин через OpenSearch
func (s *StorefrontServiceImpl) SearchOpenSearch(ctx context.Context, params *opensearch.StorefrontSearchParams) (*opensearch.StorefrontSearchResult, error) {
	return s.services.Storage().SearchStorefrontsOpenSearch(ctx, params)
}

// GetMapData получает данные для карты
func (s *StorefrontServiceImpl) GetMapData(ctx context.Context, bounds postgres.GeoBounds, filter *models.StorefrontFilter) ([]*models.StorefrontMapData, error) {
	return s.repo.GetMapData(ctx, bounds, filter)
}

// GetNearby получает витрины поблизости
func (s *StorefrontServiceImpl) GetNearby(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*models.Storefront, error) {
	return s.repo.FindNearby(ctx, lat, lng, radiusKm, limit)
}

// GetBusinessesInBuilding получает витрины в здании
func (s *StorefrontServiceImpl) GetBusinessesInBuilding(ctx context.Context, lat, lng float64) ([]*models.StorefrontMapData, error) {
	return s.repo.GetBusinessesInBuilding(ctx, lat, lng, 50) // 50 метров радиус
}

// UpdateWorkingHours обновляет часы работы
func (s *StorefrontServiceImpl) UpdateWorkingHours(ctx context.Context, userID int, storefrontID int, hours []*models.StorefrontHours) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	return s.repo.SetWorkingHours(ctx, hours)
}

// UpdatePaymentMethods обновляет методы оплаты
func (s *StorefrontServiceImpl) UpdatePaymentMethods(ctx context.Context, userID int, storefrontID int, methods []*models.StorefrontPaymentMethod) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_manage_payments")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	return s.repo.SetPaymentMethods(ctx, methods)
}

// UpdateDeliveryOptions обновляет опции доставки
func (s *StorefrontServiceImpl) UpdateDeliveryOptions(ctx context.Context, userID int, storefrontID int, options []*models.StorefrontDeliveryOption) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	return s.repo.SetDeliveryOptions(ctx, options)
}

// AddStaff добавляет сотрудника
func (s *StorefrontServiceImpl) AddStaff(ctx context.Context, ownerID int, storefrontID int, userID int, role models.StaffRole) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, ownerID, "can_manage_staff")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	staff := &models.StorefrontStaff{
		StorefrontID: storefrontID,
		UserID:       userID,
		Role:         role,
		Permissions:  s.getDefaultPermissions(role),
	}

	return s.repo.AddStaff(ctx, staff)
}

// UpdateStaffPermissions обновляет права сотрудника
func (s *StorefrontServiceImpl) UpdateStaffPermissions(ctx context.Context, ownerID int, staffID int, permissions models.JSONB) error {
	// Здесь должна быть проверка прав
	return s.repo.UpdateStaff(ctx, staffID, permissions)
}

// RemoveStaff удаляет сотрудника
func (s *StorefrontServiceImpl) RemoveStaff(ctx context.Context, ownerID int, storefrontID int, userID int) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, ownerID, "can_manage_staff")
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	return s.repo.RemoveStaff(ctx, storefrontID, userID)
}

// GetStaff получает список сотрудников
func (s *StorefrontServiceImpl) GetStaff(ctx context.Context, storefrontID int) ([]*models.StorefrontStaff, error) {
	return s.repo.GetStaff(ctx, storefrontID)
}

// UploadLogo загружает логотип
func (s *StorefrontServiceImpl) UploadLogo(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return "", fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return "", ErrInsufficientPermissions
	}

	// Загружаем файл
	path := fmt.Sprintf("storefronts/%d/logo/%s", storefrontID, filename)
	url, err := s.services.FileStorage().UploadFile(ctx, path, bytes.NewReader(data), int64(len(data)), "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки логотипа: %w", err)
	}

	// Обновляем URL в БД
	updates := &models.StorefrontUpdateDTO{
		LogoURL: &url,
	}
	if err := s.repo.Update(ctx, storefrontID, updates); err != nil {
		return "", fmt.Errorf("ошибка обновления URL логотипа: %w", err)
	}

	return url, nil
}

// UploadBanner загружает баннер
func (s *StorefrontServiceImpl) UploadBanner(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return "", fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return "", ErrInsufficientPermissions
	}

	// Загружаем файл
	path := fmt.Sprintf("storefronts/%d/banner/%s", storefrontID, filename)
	url, err := s.services.FileStorage().UploadFile(ctx, path, bytes.NewReader(data), int64(len(data)), "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки баннера: %w", err)
	}

	// Обновляем URL в БД
	updates := &models.StorefrontUpdateDTO{
		BannerURL: &url,
	}
	if err := s.repo.Update(ctx, storefrontID, updates); err != nil {
		return "", fmt.Errorf("ошибка обновления URL баннера: %w", err)
	}

	return url, nil
}

// RecordView записывает просмотр
func (s *StorefrontServiceImpl) RecordView(ctx context.Context, storefrontID int) error {
	return s.repo.RecordView(ctx, storefrontID)
}

// GetAnalytics получает аналитику
func (s *StorefrontServiceImpl) GetAnalytics(ctx context.Context, userID int, storefrontID int, from, to time.Time) ([]*models.StorefrontAnalytics, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_view_analytics")
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !hasPermission {
		return nil, ErrInsufficientPermissions
	}

	return s.repo.GetAnalyticsData(ctx, storefrontID, from, to)
}

// CanCreateStorefront проверяет может ли пользователь создать витрину
func (s *StorefrontServiceImpl) CanCreateStorefront(ctx context.Context, userID int) (bool, error) {
	return s.canCreateStorefront(ctx, userID)
}

// CheckFeatureAvailability проверяет доступность функции
func (s *StorefrontServiceImpl) CheckFeatureAvailability(ctx context.Context, storefrontID int, feature string) (bool, error) {
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		return false, err
	}

	limits, ok := planLimits[storefront.SubscriptionPlan]
	if !ok {
		return false, nil
	}

	switch feature {
	case "ai_agent":
		return limits.CanUseAI, nil
	case "live_shopping":
		return limits.CanUseLive, nil
	case "group_buying":
		return limits.CanUseGroup, nil
	case "export_data":
		return limits.CanExportData, nil
	case "custom_domain":
		return limits.CanCustomDomain, nil
	default:
		return false, nil
	}
}

// Helper методы

func (s *StorefrontServiceImpl) canCreateStorefront(ctx context.Context, userID int) (bool, error) {
	logger.Info().Int("userID", userID).Msg("canCreateStorefront: проверка возможности создания витрины")

	// Получаем текущие витрины пользователя
	currentStorefronts, err := s.ListUserStorefronts(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("canCreateStorefront: ошибка получения списка витрин")
		return false, err
	}
	logger.Info().Int("userID", userID).Int("count", len(currentStorefronts)).Msg("canCreateStorefront: текущее количество витрин")

	// TODO: Интеграция с сервисом подписок для получения актуального плана пользователя
	// Пока используем бесплатный план по умолчанию
	userPlan := models.SubscriptionPlanStarter

	limits, ok := planLimits[userPlan]
	if !ok {
		logger.Error().Str("plan", string(userPlan)).Msg("canCreateStorefront: неизвестный план")
		return false, nil
	}

	// Проверяем лимит
	if limits.MaxStorefronts == -1 {
		logger.Info().Msg("canCreateStorefront: безлимитный план")
		return true, nil // Unlimited
	}

	canCreate := len(currentStorefronts) < limits.MaxStorefronts
	logger.Info().
		Int("current", len(currentStorefronts)).
		Int("limit", limits.MaxStorefronts).
		Bool("canCreate", canCreate).
		Msg("canCreateStorefront: результат проверки")

	return canCreate, nil
}

func (s *StorefrontServiceImpl) generateSlug(name string) string {
	slug := strings.ToLower(name)

	// Заменяем пробелы на дефисы
	slug = strings.ReplaceAll(slug, " ", "-")

	// Удаляем не-ASCII символы
	var result []rune
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			result = append(result, r)
		}
	}

	return string(result)
}

func (s *StorefrontServiceImpl) getDefaultPermissions(role models.StaffRole) models.JSONB {
	switch role {
	case models.StaffRoleOwner:
		return models.JSONB{
			"can_add_products":     true,
			"can_edit_products":    true,
			"can_delete_products":  true,
			"can_view_orders":      true,
			"can_process_orders":   true,
			"can_refund_orders":    true,
			"can_edit_storefront":  true,
			"can_manage_staff":     true,
			"can_view_analytics":   true,
			"can_manage_payments":  true,
			"can_reply_to_reviews": true,
			"can_send_messages":    true,
		}
	case models.StaffRoleManager:
		return models.JSONB{
			"can_add_products":     true,
			"can_edit_products":    true,
			"can_delete_products":  true,
			"can_view_orders":      true,
			"can_process_orders":   true,
			"can_refund_orders":    true,
			"can_edit_storefront":  false,
			"can_manage_staff":     false,
			"can_view_analytics":   true,
			"can_manage_payments":  false,
			"can_reply_to_reviews": true,
			"can_send_messages":    true,
		}
	case models.StaffRoleCashier:
		return models.JSONB{
			"can_add_products":     false,
			"can_edit_products":    false,
			"can_delete_products":  false,
			"can_view_orders":      true,
			"can_process_orders":   true,
			"can_refund_orders":    false,
			"can_edit_storefront":  false,
			"can_manage_staff":     false,
			"can_view_analytics":   false,
			"can_manage_payments":  false,
			"can_reply_to_reviews": false,
			"can_send_messages":    true,
		}
	case models.StaffRoleSupport:
		return models.JSONB{
			"can_add_products":     false,
			"can_edit_products":    false,
			"can_delete_products":  false,
			"can_view_orders":      true,
			"can_process_orders":   false,
			"can_refund_orders":    false,
			"can_edit_storefront":  false,
			"can_manage_staff":     false,
			"can_view_analytics":   false,
			"can_manage_payments":  false,
			"can_reply_to_reviews": true,
			"can_send_messages":    true,
		}
	case models.StaffRoleModerator:
		return models.JSONB{
			"can_add_products":     false,
			"can_edit_products":    false,
			"can_delete_products":  false,
			"can_view_orders":      true,
			"can_process_orders":   false,
			"can_refund_orders":    false,
			"can_edit_storefront":  false,
			"can_manage_staff":     false,
			"can_view_analytics":   false,
			"can_manage_payments":  false,
			"can_reply_to_reviews": true,
			"can_send_messages":    true,
		}
	default:
		return models.JSONB{}
	}
}
