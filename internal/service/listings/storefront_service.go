package listings

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// StorefrontRepository defines the interface for storefront data access
type StorefrontRepository interface {
	CreateStorefront(ctx context.Context, storefront *domain.Storefront) error
	GetStorefrontByID(ctx context.Context, id int64, includes *domain.Includes) (*domain.Storefront, error)
	GetStorefrontBySlug(ctx context.Context, slug string, includes *domain.Includes) (*domain.Storefront, error)
	UpdateStorefront(ctx context.Context, id int64, update *domain.StorefrontUpdate) error
	DeleteStorefront(ctx context.Context, id int64, hardDelete bool) error
	ListStorefronts(ctx context.Context, filter *domain.ListStorefrontsFilter) ([]domain.Storefront, int, error)
	AddStaff(ctx context.Context, staff *domain.StorefrontStaff) error
	UpdateStaff(ctx context.Context, id int64, update *domain.StaffUpdate) error
	RemoveStaff(ctx context.Context, storefrontID, userID int64) error
	GetStaff(ctx context.Context, storefrontID int64) ([]domain.StorefrontStaff, error)
	SetWorkingHours(ctx context.Context, storefrontID int64, hours []domain.StorefrontHours) error
	GetWorkingHours(ctx context.Context, storefrontID int64) ([]domain.StorefrontHours, error)
	IsOpenNow(ctx context.Context, storefrontID int64) (bool, *time.Time, *time.Time, error)
	SetPaymentMethods(ctx context.Context, storefrontID int64, methods []domain.PaymentMethod) error
	GetPaymentMethods(ctx context.Context, storefrontID int64) ([]domain.PaymentMethod, error)
	SetDeliveryOptions(ctx context.Context, storefrontID int64, options []domain.StorefrontDeliveryOption) error
	GetDeliveryOptions(ctx context.Context, storefrontID int64) ([]domain.StorefrontDeliveryOption, error)
	GetMapData(ctx context.Context, bounds *domain.MapBounds, filter *domain.ListStorefrontsFilter) ([]domain.StorefrontMapData, error)
	GetStorefrontDashboardStats(ctx context.Context, storefrontID int64, from, to *time.Time) (*domain.StorefrontDashboardStats, error)
	IsSlugTaken(ctx context.Context, slug string, excludeID *int64) (bool, error)
	IncrementViewsCount(ctx context.Context, storefrontID int64) error
}

// CreateStorefrontRequest represents storefront creation request
type CreateStorefrontRequest struct {
	UserID      int64
	Name        string
	Slug        string
	Description *string
	Logo        []byte
	Banner      []byte
	Theme       domain.JSONB
	Phone       *string
	Email       *string
	Website     *string
	Location    StorefrontLocation
	Settings    domain.JSONB
	SeoMeta     domain.JSONB
}

// StorefrontLocation represents storefront address information
type StorefrontLocation struct {
	UserLat     float64
	UserLng     float64
	FullAddress string
	City        string
	PostalCode  *string
	Country     string
}

// UpdateStorefrontRequest represents storefront update request
type UpdateStorefrontRequest struct {
	Name                *string
	Description         *string
	IsActive            *bool
	LogoURL             *string
	BannerURL           *string
	Theme               domain.JSONB
	Phone               *string
	Email               *string
	Website             *string
	Location            *StorefrontLocation
	Settings            domain.JSONB
	SeoMeta             domain.JSONB
	AIAgentEnabled      *bool
	LiveShoppingEnabled *bool
	GroupBuyingEnabled  *bool
}

// StorefrontService handles storefront business logic
type StorefrontService struct {
	repo   StorefrontRepository
	logger *zerolog.Logger
}

// NewStorefrontService creates a new storefront service
func NewStorefrontService(repo StorefrontRepository, logger *zerolog.Logger) *StorefrontService {
	return &StorefrontService{
		repo:   repo,
		logger: logger,
	}
}

// CreateStorefront creates a new storefront
func (s *StorefrontService) CreateStorefront(ctx context.Context, req *CreateStorefrontRequest) (*domain.Storefront, error) {
	if err := validateCreateStorefrontRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	} else {
		slug = generateSlug(slug)
	}

	taken, err := s.repo.IsSlugTaken(ctx, slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug: %w", err)
	}
	if taken {
		for i := 1; i <= 100; i++ {
			newSlug := fmt.Sprintf("%s-%d", slug, i)
			taken, err = s.repo.IsSlugTaken(ctx, newSlug, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to check slug: %w", err)
			}
			if !taken {
				slug = newSlug
				break
			}
		}
		if taken {
			return nil, fmt.Errorf("failed to generate unique slug")
		}
	}

	storefront := &domain.Storefront{
		UserID:               req.UserID,
		Slug:                 slug,
		Name:                 req.Name,
		Description:          req.Description,
		Theme:                req.Theme,
		Phone:                req.Phone,
		Email:                req.Email,
		Website:              req.Website,
		Address:              &req.Location.FullAddress,
		City:                 &req.Location.City,
		PostalCode:           req.Location.PostalCode,
		Country:              &req.Location.Country,
		Latitude:             &req.Location.UserLat,
		Longitude:            &req.Location.UserLng,
		FormattedAddress:     &req.Location.FullAddress,
		GeoStrategy:          "storefront_location",
		DefaultPrivacyLevel:  "exact",
		AddressVerified:      false,
		Settings:             req.Settings,
		SeoMeta:              req.SeoMeta,
		IsActive:             true,
		IsVerified:           false,
		Rating:               0.0,
		ReviewsCount:         0,
		ProductsCount:        0,
		SalesCount:           0,
		ViewsCount:           0,
		SubscriptionPlan:     "starter",
		CommissionRate:       3.00,
		IsSubscriptionActive: true,
		AIAgentEnabled:       false,
		LiveShoppingEnabled:  false,
		GroupBuyingEnabled:   false,
		FollowersCount:       0,
	}

	if err := s.repo.CreateStorefront(ctx, storefront); err != nil {
		return nil, fmt.Errorf("failed to create storefront: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefront.ID).
		Str("slug", storefront.Slug).
		Msg("storefront created successfully")

	return storefront, nil
}

// GetStorefront retrieves a storefront by ID or slug
func (s *StorefrontService) GetStorefront(ctx context.Context, id *int64, slug *string, includes *domain.Includes) (*domain.Storefront, error) {
	if id == nil && slug == nil {
		return nil, fmt.Errorf("either id or slug must be provided")
	}

	var storefront *domain.Storefront
	var err error

	if id != nil {
		storefront, err = s.repo.GetStorefrontByID(ctx, *id, includes)
	} else {
		storefront, err = s.repo.GetStorefrontBySlug(ctx, *slug, includes)
	}

	if err != nil {
		return nil, err
	}

	return storefront, nil
}

// UpdateStorefront updates a storefront
func (s *StorefrontService) UpdateStorefront(ctx context.Context, id int64, req *UpdateStorefrontRequest) (*domain.Storefront, error) {
	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return nil, err
		}
	}

	if req.Location != nil {
		if err := validateCoordinates(req.Location.UserLat, req.Location.UserLng); err != nil {
			return nil, err
		}
	}

	update := &domain.StorefrontUpdate{
		Name:                req.Name,
		Description:         req.Description,
		IsActive:            req.IsActive,
		LogoURL:             req.LogoURL,
		BannerURL:           req.BannerURL,
		Theme:               req.Theme,
		Phone:               req.Phone,
		Email:               req.Email,
		Website:             req.Website,
		Settings:            req.Settings,
		SeoMeta:             req.SeoMeta,
		AIAgentEnabled:      req.AIAgentEnabled,
		LiveShoppingEnabled: req.LiveShoppingEnabled,
		GroupBuyingEnabled:  req.GroupBuyingEnabled,
	}

	if req.Location != nil {
		update.Address = &req.Location.FullAddress
		update.City = &req.Location.City
		update.PostalCode = req.Location.PostalCode
		update.Country = &req.Location.Country
		update.Latitude = &req.Location.UserLat
		update.Longitude = &req.Location.UserLng
		update.FormattedAddress = &req.Location.FullAddress
	}

	if err := s.repo.UpdateStorefront(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update storefront: %w", err)
	}

	storefront, err := s.repo.GetStorefrontByID(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated storefront: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", id).
		Msg("storefront updated successfully")

	return storefront, nil
}

// DeleteStorefront deletes a storefront
func (s *StorefrontService) DeleteStorefront(ctx context.Context, id int64, hardDelete bool) error {
	if err := s.repo.DeleteStorefront(ctx, id, hardDelete); err != nil {
		return fmt.Errorf("failed to delete storefront: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", id).
		Bool("hard_delete", hardDelete).
		Msg("storefront deleted successfully")

	return nil
}

// ListStorefronts lists storefronts with filters
func (s *StorefrontService) ListStorefronts(ctx context.Context, filter *domain.ListStorefrontsFilter) ([]domain.Storefront, int, error) {
	storefronts, total, err := s.repo.ListStorefronts(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list storefronts: %w", err)
	}

	return storefronts, total, nil
}

// AddStaff adds a staff member to a storefront
func (s *StorefrontService) AddStaff(ctx context.Context, storefrontID, userID int64, role string, permissions domain.JSONB) (*domain.StorefrontStaff, error) {
	if err := validateStaffRole(role); err != nil {
		return nil, err
	}

	staff := &domain.StorefrontStaff{
		StorefrontID: storefrontID,
		UserID:       userID,
		Role:         role,
		Permissions:  permissions,
		ActionsCount: 0,
	}

	if err := s.repo.AddStaff(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to add staff: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefrontID).
		Int64("user_id", userID).
		Str("role", role).
		Msg("staff member added successfully")

	return staff, nil
}

// UpdateStaff updates a staff member
func (s *StorefrontService) UpdateStaff(ctx context.Context, staffID int64, update *domain.StaffUpdate) (*domain.StorefrontStaff, error) {
	if update.Role != nil {
		if err := validateStaffRole(*update.Role); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateStaff(ctx, staffID, update); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	s.logger.Info().
		Int64("staff_id", staffID).
		Msg("staff member updated successfully")

	return nil, nil
}

// RemoveStaff removes a staff member from a storefront
func (s *StorefrontService) RemoveStaff(ctx context.Context, storefrontID, userID int64) error {
	if err := s.repo.RemoveStaff(ctx, storefrontID, userID); err != nil {
		return fmt.Errorf("failed to remove staff: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefrontID).
		Int64("user_id", userID).
		Msg("staff member removed successfully")

	return nil
}

// GetStaff retrieves all staff members for a storefront
func (s *StorefrontService) GetStaff(ctx context.Context, storefrontID int64) ([]domain.StorefrontStaff, error) {
	staff, err := s.repo.GetStaff(ctx, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return staff, nil
}

// SetWorkingHours sets working hours for a storefront
func (s *StorefrontService) SetWorkingHours(ctx context.Context, storefrontID int64, hours []domain.StorefrontHours) error {
	for _, hour := range hours {
		if err := validateWorkingHour(&hour); err != nil {
			return err
		}
	}

	if err := s.repo.SetWorkingHours(ctx, storefrontID, hours); err != nil {
		return fmt.Errorf("failed to set working hours: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefrontID).
		Int("hours_count", len(hours)).
		Msg("working hours set successfully")

	return nil
}

// GetWorkingHours retrieves working hours for a storefront
func (s *StorefrontService) GetWorkingHours(ctx context.Context, storefrontID int64) ([]domain.StorefrontHours, error) {
	hours, err := s.repo.GetWorkingHours(ctx, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get working hours: %w", err)
	}

	return hours, nil
}

// IsOpenNow checks if a storefront is currently open
func (s *StorefrontService) IsOpenNow(ctx context.Context, storefrontID int64) (bool, *time.Time, *time.Time, error) {
	return s.repo.IsOpenNow(ctx, storefrontID)
}

// SetPaymentMethods sets payment methods for a storefront
func (s *StorefrontService) SetPaymentMethods(ctx context.Context, storefrontID int64, methods []domain.PaymentMethod) error {
	for _, method := range methods {
		if err := validatePaymentMethod(&method); err != nil {
			return err
		}
	}

	if err := s.repo.SetPaymentMethods(ctx, storefrontID, methods); err != nil {
		return fmt.Errorf("failed to set payment methods: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefrontID).
		Int("methods_count", len(methods)).
		Msg("payment methods set successfully")

	return nil
}

// GetPaymentMethods retrieves payment methods for a storefront
func (s *StorefrontService) GetPaymentMethods(ctx context.Context, storefrontID int64) ([]domain.PaymentMethod, error) {
	methods, err := s.repo.GetPaymentMethods(ctx, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	return methods, nil
}

// SetDeliveryOptions sets delivery options for a storefront
func (s *StorefrontService) SetDeliveryOptions(ctx context.Context, storefrontID int64, options []domain.StorefrontDeliveryOption) error {
	for _, option := range options {
		if err := validateDeliveryOption(&option); err != nil {
			return err
		}
	}

	if err := s.repo.SetDeliveryOptions(ctx, storefrontID, options); err != nil {
		return fmt.Errorf("failed to set delivery options: %w", err)
	}

	s.logger.Info().
		Int64("storefront_id", storefrontID).
		Int("options_count", len(options)).
		Msg("delivery options set successfully")

	return nil
}

// GetDeliveryOptions retrieves delivery options for a storefront
func (s *StorefrontService) GetDeliveryOptions(ctx context.Context, storefrontID int64) ([]domain.StorefrontDeliveryOption, error) {
	options, err := s.repo.GetDeliveryOptions(ctx, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery options: %w", err)
	}

	return options, nil
}

// GetMapData retrieves storefronts for map display
func (s *StorefrontService) GetMapData(ctx context.Context, bounds *domain.MapBounds, filter *domain.ListStorefrontsFilter) ([]domain.StorefrontMapData, error) {
	mapData, err := s.repo.GetMapData(ctx, bounds, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get map data: %w", err)
	}

	return mapData, nil
}

// GetDashboardStats retrieves dashboard statistics for a storefront
func (s *StorefrontService) GetDashboardStats(ctx context.Context, storefrontID int64, from, to *time.Time) (*domain.StorefrontDashboardStats, error) {
	stats, err := s.repo.GetStorefrontDashboardStats(ctx, storefrontID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return stats, nil
}

// Validation helper functions

func validateCreateStorefrontRequest(req *CreateStorefrontRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 255 {
		return fmt.Errorf("name too long (max 255 characters)")
	}
	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return err
		}
	}
	if err := validateCoordinates(req.Location.UserLat, req.Location.UserLng); err != nil {
		return err
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return nil
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

func validateStaffRole(role string) error {
	validRoles := map[string]bool{
		"owner":     true,
		"manager":   true,
		"cashier":   true,
		"support":   true,
		"moderator": true,
	}
	if !validRoles[role] {
		return fmt.Errorf("invalid staff role: %s", role)
	}
	return nil
}

func validateWorkingHour(hour *domain.StorefrontHours) error {
	if hour.DayOfWeek < 0 || hour.DayOfWeek > 6 {
		return fmt.Errorf("day_of_week must be between 0 and 6")
	}
	return nil
}

func validatePaymentMethod(method *domain.PaymentMethod) error {
	if method.MethodType == "" {
		return fmt.Errorf("method_type is required")
	}
	validTypes := map[string]bool{
		"cash": true, "cod": true, "card": true, "bank_transfer": true,
		"paypal": true, "crypto": true, "postanska": true, "keks_pay": true, "ips": true,
	}
	if !validTypes[method.MethodType] {
		return fmt.Errorf("invalid payment method type: %s", method.MethodType)
	}
	return nil
}

func validateDeliveryOption(option *domain.StorefrontDeliveryOption) error {
	if option.Name == "" {
		return fmt.Errorf("name is required")
	}
	if option.EstimatedDaysMin < 0 || option.EstimatedDaysMax < 0 {
		return fmt.Errorf("estimated days must be non-negative")
	}
	if option.EstimatedDaysMin > option.EstimatedDaysMax {
		return fmt.Errorf("estimated_days_min cannot be greater than estimated_days_max")
	}
	return nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("store-%d", time.Now().Unix()%10000)
	}
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.TrimRight(slug, "-")
	}
	return slug
}
