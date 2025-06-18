package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage/postgres"
	"backend/internal/storage/filestorage"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Ошибки
var (
	ErrUnauthorized           = errors.New("unauthorized")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrStorefrontLimitReached = errors.New("storefront limit reached for current plan")
	ErrSlugAlreadyExists      = errors.New("slug already exists")
	ErrInvalidLocation        = errors.New("invalid location data")
	ErrFeatureNotAvailable    = errors.New("feature not available in current plan")
)

// Лимиты по тарифным планам
var planLimits = map[models.SubscriptionPlan]struct {
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
	Create(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error)
	GetByID(ctx context.Context, id int) (*models.Storefront, error)
	GetBySlug(ctx context.Context, slug string) (*models.Storefront, error)
	Update(ctx context.Context, userID int, storefrontID int, dto *models.StorefrontUpdateDTO) error
	Delete(ctx context.Context, userID int, storefrontID int) error
	
	// Листинг и поиск
	ListUserStorefronts(ctx context.Context, userID int) ([]*models.Storefront, error)
	Search(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error)
	
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

// storefrontService реализация сервиса
type storefrontService struct {
	repo     postgres.StorefrontRepository
	fileRepo filestorage.FileStorageInterface
}

// NewStorefrontService создает новый сервис витрин
func NewStorefrontService(repo postgres.StorefrontRepository, fileRepo filestorage.FileStorageInterface) StorefrontService {
	return &storefrontService{
		repo:     repo,
		fileRepo: fileRepo,
	}
}

// Create создает новую витрину
func (s *storefrontService) Create(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	// Проверяем можно ли создать еще одну витрину
	canCreate, err := s.CanCreateStorefront(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !canCreate {
		return nil, ErrStorefrontLimitReached
	}

	// Валидация локации
	if err := s.validateLocation(&dto.Location); err != nil {
		return nil, err
	}

	// Генерируем уникальный slug
	slug := s.generateSlug(dto.Name)
	existingCount := 0
	for {
		existing, _ := s.repo.GetBySlug(ctx, slug)
		if existing == nil {
			break
		}
		existingCount++
		slug = fmt.Sprintf("%s-%d", s.generateSlug(dto.Name), existingCount)
	}

	// Создаем витрину через репозиторий
	// TODO: получить userID из контекста
	storefront, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create storefront: %w", err)
	}

	// Загружаем логотип если есть
	if len(dto.Logo) > 0 {
		logoURL, err := s.UploadLogo(ctx, userID, storefront.ID, dto.Logo, "logo.jpg")
		if err != nil {
			// Не критично, продолжаем
			fmt.Printf("Failed to upload logo: %v\n", err)
		} else {
			storefront.LogoURL = logoURL
			// TODO: обновить в БД
		}
	}

	// Загружаем баннер если есть
	if len(dto.Banner) > 0 {
		bannerURL, err := s.UploadBanner(ctx, userID, storefront.ID, dto.Banner, "banner.jpg")
		if err != nil {
			fmt.Printf("Failed to upload banner: %v\n", err)
		} else {
			storefront.BannerURL = bannerURL
			// TODO: обновить в БД
		}
	}

	return storefront, nil
}

// GetByID получает витрину по ID
func (s *storefrontService) GetByID(ctx context.Context, id int) (*models.Storefront, error) {
	return s.repo.GetByID(ctx, id)
}

// GetBySlug получает витрину по slug
func (s *storefrontService) GetBySlug(ctx context.Context, slug string) (*models.Storefront, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// Update обновляет витрину
func (s *storefrontService) Update(ctx context.Context, userID int, storefrontID int, dto *models.StorefrontUpdateDTO) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Валидация локации если обновляется
	if dto.Location != nil {
		if err := s.validateLocation(dto.Location); err != nil {
			return err
		}
	}

	// Проверяем доступность функций если включаются
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		return err
	}

	limits := planLimits[storefront.SubscriptionPlan]
	
	if dto.AIAgentEnabled != nil && *dto.AIAgentEnabled && !limits.CanUseAI {
		return ErrFeatureNotAvailable
	}
	
	if dto.LiveShoppingEnabled != nil && *dto.LiveShoppingEnabled && !limits.CanUseLive {
		return ErrFeatureNotAvailable
	}
	
	if dto.GroupBuyingEnabled != nil && *dto.GroupBuyingEnabled && !limits.CanUseGroup {
		return ErrFeatureNotAvailable
	}

	return s.repo.Update(ctx, storefrontID, dto)
}

// Delete удаляет витрину
func (s *storefrontService) Delete(ctx context.Context, userID int, storefrontID int) error {
	// Только владелец может удалить
	isOwner, err := s.repo.IsOwner(ctx, storefrontID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, storefrontID)
}

// ListUserStorefronts получает витрины пользователя
func (s *storefrontService) ListUserStorefronts(ctx context.Context, userID int) ([]*models.Storefront, error) {
	filter := &models.StorefrontFilter{
		UserID: &userID,
		Limit:  100,
		Offset: 0,
	}
	
	storefronts, _, err := s.repo.List(ctx, filter)
	return storefronts, err
}

// Search ищет витрины
func (s *storefrontService) Search(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error) {
	// Устанавливаем дефолтные значения
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	return s.repo.List(ctx, filter)
}

// GetMapData получает данные для карты
func (s *storefrontService) GetMapData(ctx context.Context, bounds postgres.GeoBounds, filter *models.StorefrontFilter) ([]*models.StorefrontMapData, error) {
	return s.repo.GetMapData(ctx, bounds, filter)
}

// GetNearby получает ближайшие витрины
func (s *storefrontService) GetNearby(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*models.Storefront, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	return s.repo.FindNearby(ctx, lat, lng, radiusKm, limit)
}

// GetBusinessesInBuilding получает все бизнесы в здании
func (s *storefrontService) GetBusinessesInBuilding(ctx context.Context, lat, lng float64) ([]*models.StorefrontMapData, error) {
	return s.repo.GetBusinessesInBuilding(ctx, lat, lng, 30) // 30 метров радиус
}

// UpdateWorkingHours обновляет часы работы
func (s *storefrontService) UpdateWorkingHours(ctx context.Context, userID int, storefrontID int, hours []*models.StorefrontHours) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Проставляем storefrontID всем записям
	for _, h := range hours {
		h.StorefrontID = storefrontID
	}

	return s.repo.SetWorkingHours(ctx, hours)
}

// UpdatePaymentMethods обновляет методы оплаты
func (s *storefrontService) UpdatePaymentMethods(ctx context.Context, userID int, storefrontID int, methods []*models.StorefrontPaymentMethod) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_manage_payments")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Проставляем storefrontID всем записям
	for _, m := range methods {
		m.StorefrontID = storefrontID
	}

	return s.repo.SetPaymentMethods(ctx, methods)
}

// UpdateDeliveryOptions обновляет опции доставки
func (s *storefrontService) UpdateDeliveryOptions(ctx context.Context, userID int, storefrontID int, options []*models.StorefrontDeliveryOption) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Проставляем storefrontID всем записям
	for _, opt := range options {
		opt.StorefrontID = storefrontID
	}

	return s.repo.SetDeliveryOptions(ctx, options)
}

// AddStaff добавляет сотрудника
func (s *storefrontService) AddStaff(ctx context.Context, ownerID int, storefrontID int, userID int, role models.StaffRole) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, ownerID, "can_manage_staff")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	// Проверяем лимит сотрудников
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		return err
	}

	currentStaff, err := s.repo.GetStaff(ctx, storefrontID)
	if err != nil {
		return err
	}

	limits := planLimits[storefront.SubscriptionPlan]
	if limits.MaxStaff != -1 && len(currentStaff) >= limits.MaxStaff {
		return errors.New("staff limit reached for current plan")
	}

	// Определяем права по роли
	permissions := s.getDefaultPermissionsByRole(role)

	staff := &models.StorefrontStaff{
		StorefrontID: storefrontID,
		UserID:       userID,
		Role:         role,
		Permissions:  permissions,
	}

	return s.repo.AddStaff(ctx, staff)
}

// UpdateStaffPermissions обновляет права сотрудника
func (s *storefrontService) UpdateStaffPermissions(ctx context.Context, ownerID int, staffID int, permissions models.JSONB) error {
	// TODO: проверить что ownerID имеет права manage_staff для витрины где работает staffID
	return s.repo.UpdateStaff(ctx, staffID, permissions)
}

// RemoveStaff удаляет сотрудника
func (s *storefrontService) RemoveStaff(ctx context.Context, ownerID int, storefrontID int, userID int) error {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, ownerID, "can_manage_staff")
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInsufficientPermissions
	}

	return s.repo.RemoveStaff(ctx, storefrontID, userID)
}

// GetStaff получает список сотрудников
func (s *storefrontService) GetStaff(ctx context.Context, storefrontID int) ([]*models.StorefrontStaff, error) {
	return s.repo.GetStaff(ctx, storefrontID)
}

// UploadLogo загружает логотип
func (s *storefrontService) UploadLogo(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return "", err
	}
	if !hasPermission {
		return "", ErrInsufficientPermissions
	}

	// Генерируем путь
	path := fmt.Sprintf("storefronts/%d/logo_%d.jpg", storefrontID, time.Now().Unix())
	
	// Загружаем в хранилище
	reader := strings.NewReader(string(data))
	url, err := s.fileRepo.UploadFile(ctx, path, reader, int64(len(data)), "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("failed to upload logo: %w", err)
	}

	// Обновляем URL в БД
	updateDTO := &models.StorefrontUpdateDTO{}
	// TODO: добавить поле LogoURL в DTO
	err = s.repo.Update(ctx, storefrontID, updateDTO)
	if err != nil {
		return "", fmt.Errorf("failed to update logo URL: %w", err)
	}

	return url, nil
}

// UploadBanner загружает баннер
func (s *storefrontService) UploadBanner(ctx context.Context, userID int, storefrontID int, data []byte, filename string) (string, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_edit_storefront")
	if err != nil {
		return "", err
	}
	if !hasPermission {
		return "", ErrInsufficientPermissions
	}

	// Генерируем путь
	path := fmt.Sprintf("storefronts/%d/banner_%d.jpg", storefrontID, time.Now().Unix())
	
	// Загружаем в хранилище
	reader := strings.NewReader(string(data))
	url, err := s.fileRepo.UploadFile(ctx, path, reader, int64(len(data)), "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("failed to upload banner: %w", err)
	}

	// TODO: обновить URL в БД

	return url, nil
}

// RecordView записывает просмотр
func (s *storefrontService) RecordView(ctx context.Context, storefrontID int) error {
	return s.repo.RecordView(ctx, storefrontID)
}

// GetAnalytics получает аналитику
func (s *storefrontService) GetAnalytics(ctx context.Context, userID int, storefrontID int, from, to time.Time) ([]*models.StorefrontAnalytics, error) {
	// Проверяем права
	hasPermission, err := s.repo.HasPermission(ctx, storefrontID, userID, "can_view_analytics")
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, ErrInsufficientPermissions
	}

	return s.repo.GetAnalytics(ctx, storefrontID, from, to)
}

// CanCreateStorefront проверяет можно ли создать витрину
func (s *storefrontService) CanCreateStorefront(ctx context.Context, userID int) (bool, error) {
	// Получаем текущие витрины пользователя
	storefronts, err := s.ListUserStorefronts(ctx, userID)
	if err != nil {
		return false, err
	}

	// Определяем план пользователя (берем максимальный из всех витрин)
	maxPlan := models.SubscriptionPlanStarter
	for _, sf := range storefronts {
		if sf.SubscriptionPlan > maxPlan {
			maxPlan = sf.SubscriptionPlan
		}
	}

	// Проверяем лимит
	limits := planLimits[maxPlan]
	if limits.MaxStorefronts == -1 {
		return true, nil
	}

	return len(storefronts) < limits.MaxStorefronts, nil
}

// CheckFeatureAvailability проверяет доступность функции
func (s *storefrontService) CheckFeatureAvailability(ctx context.Context, storefrontID int, feature string) (bool, error) {
	storefront, err := s.repo.GetByID(ctx, storefrontID)
	if err != nil {
		return false, err
	}

	limits := planLimits[storefront.SubscriptionPlan]

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
		return false, fmt.Errorf("unknown feature: %s", feature)
	}
}

// Вспомогательные функции

// generateSlug генерирует slug из названия
func (s *storefrontService) generateSlug(name string) string {
	// Приводим к нижнему регистру
	slug := strings.ToLower(name)
	
	// Заменяем пробелы на дефисы
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Оставляем только латиницу, цифры и дефисы
	var result []rune
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result = append(result, r)
		}
	}
	
	return string(result)
}

// validateLocation валидирует данные локации
func (s *storefrontService) validateLocation(location *models.Location) error {
	if location == nil {
		return ErrInvalidLocation
	}

	// Проверяем координаты
	if location.BuildingLat < -90 || location.BuildingLat > 90 {
		return fmt.Errorf("invalid latitude: %f", location.BuildingLat)
	}
	
	if location.BuildingLng < -180 || location.BuildingLng > 180 {
		return fmt.Errorf("invalid longitude: %f", location.BuildingLng)
	}

	// Проверяем обязательные поля
	if location.FullAddress == "" {
		return errors.New("address is required")
	}
	
	if location.City == "" {
		return errors.New("city is required")
	}

	return nil
}

// getDefaultPermissionsByRole возвращает права по умолчанию для роли
func (s *storefrontService) getDefaultPermissionsByRole(role models.StaffRole) models.JSONB {
	switch role {
	case models.StaffRoleManager:
		return models.JSONB{
			"can_add_products":    true,
			"can_edit_products":   true,
			"can_delete_products": true,
			"can_view_orders":     true,
			"can_process_orders":  true,
			"can_refund_orders":   true,
			"can_edit_storefront": false,
			"can_manage_staff":    false,
			"can_view_analytics":  true,
			"can_manage_payments": false,
			"can_reply_to_reviews": true,
			"can_send_messages":   true,
		}
	case models.StaffRoleSupport:
		return models.JSONB{
			"can_add_products":    false,
			"can_edit_products":   false,
			"can_delete_products": false,
			"can_view_orders":     true,
			"can_process_orders":  true,
			"can_refund_orders":   false,
			"can_edit_storefront": false,
			"can_manage_staff":    false,
			"can_view_analytics":  false,
			"can_manage_payments": false,
			"can_reply_to_reviews": true,
			"can_send_messages":   true,
		}
	case models.StaffRoleModerator:
		return models.JSONB{
			"can_add_products":    false,
			"can_edit_products":   true,
			"can_delete_products": false,
			"can_view_orders":     false,
			"can_process_orders":  false,
			"can_refund_orders":   false,
			"can_edit_storefront": false,
			"can_manage_staff":    false,
			"can_view_analytics":  false,
			"can_manage_payments": false,
			"can_reply_to_reviews": true,
			"can_send_messages":   false,
		}
	default:
		// Минимальные права по умолчанию
		return models.JSONB{
			"can_add_products":    false,
			"can_edit_products":   false,
			"can_delete_products": false,
			"can_view_orders":     false,
			"can_process_orders":  false,
			"can_refund_orders":   false,
			"can_edit_storefront": false,
			"can_manage_staff":    false,
			"can_view_analytics":  false,
			"can_manage_payments": false,
			"can_reply_to_reviews": false,
			"can_send_messages":   false,
		}
	}
}