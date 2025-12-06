package domain

import (
	"database/sql/driver"
	"errors"
	"time"
)

// JSONB represents a PostgreSQL JSONB column
type JSONB []byte

// Value implements the driver.Valuer interface for database/sql
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

// Scan implements the sql.Scanner interface for database/sql
func (j *JSONB) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}
	switch v := src.(type) {
	case []byte:
		*j = JSONB(v)
		return nil
	case string:
		*j = JSONB(v)
		return nil
	default:
		return errors.New("incompatible type for JSONB")
	}
}

// MarshalJSON returns j as the JSON encoding of j
func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONB: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Storefront represents B2C store entity (storefronts table)
type Storefront struct {
	// Identification
	ID          int64   `db:"id" json:"id"`
	UserID      int64   `db:"user_id" json:"user_id"`
	Slug        string  `db:"slug" json:"slug"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`

	// Branding
	LogoURL   *string `db:"logo_url" json:"logo_url,omitempty"`
	BannerURL *string `db:"banner_url" json:"banner_url,omitempty"`
	Theme     JSONB   `db:"theme" json:"theme,omitempty"`

	// Contact Information
	Phone   *string `db:"phone" json:"phone,omitempty"`
	Email   *string `db:"email" json:"email,omitempty"`
	Website *string `db:"website" json:"website,omitempty"`

	// Location
	Address             *string  `db:"address" json:"address,omitempty"`
	City                *string  `db:"city" json:"city,omitempty"`
	PostalCode          *string  `db:"postal_code" json:"postal_code,omitempty"`
	Country             *string  `db:"country" json:"country,omitempty"`
	Latitude            *float64 `db:"latitude" json:"latitude,omitempty"`
	Longitude           *float64 `db:"longitude" json:"longitude,omitempty"`
	FormattedAddress    *string  `db:"formatted_address" json:"formatted_address,omitempty"`
	GeoStrategy         string   `db:"geo_strategy" json:"geo_strategy"`
	DefaultPrivacyLevel string   `db:"default_privacy_level" json:"default_privacy_level"`
	AddressVerified     bool     `db:"address_verified" json:"address_verified"`

	// Settings
	Settings JSONB `db:"settings" json:"settings,omitempty"`
	SeoMeta  JSONB `db:"seo_meta" json:"seo_meta,omitempty"`

	// Status and Statistics
	IsActive         bool       `db:"is_active" json:"is_active"`
	IsVerified       bool       `db:"is_verified" json:"is_verified"`
	VerificationDate *time.Time `db:"verification_date" json:"verification_date,omitempty"`
	Rating           float64    `db:"rating" json:"rating"`
	ReviewsCount     int32      `db:"reviews_count" json:"reviews_count"`
	ProductsCount    int32      `db:"products_count" json:"products_count"`
	SalesCount       int32      `db:"sales_count" json:"sales_count"`
	ViewsCount       int32      `db:"views_count" json:"views_count"`

	// Subscription (Monetization)
	SubscriptionPlan      string     `db:"subscription_plan" json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `db:"subscription_expires_at" json:"subscription_expires_at,omitempty"`
	CommissionRate        float64    `db:"commission_rate" json:"commission_rate"`
	SubscriptionID        *int64     `db:"subscription_id" json:"subscription_id,omitempty"`
	IsSubscriptionActive  bool       `db:"is_subscription_active" json:"is_subscription_active"`

	// AI and Killer Features
	AIAgentEnabled      bool  `db:"ai_agent_enabled" json:"ai_agent_enabled"`
	AIAgentConfig       JSONB `db:"ai_agent_config" json:"ai_agent_config,omitempty"`
	LiveShoppingEnabled bool  `db:"live_shopping_enabled" json:"live_shopping_enabled"`
	GroupBuyingEnabled  bool  `db:"group_buying_enabled" json:"group_buying_enabled"`

	// Social Features
	FollowersCount int32 `db:"followers_count" json:"followers_count"`

	// Timestamps
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Soft delete support

	// Related entities (loaded separately)
	Staff           []StorefrontStaff          `db:"-" json:"staff,omitempty"`
	Hours           []StorefrontHours          `db:"-" json:"hours,omitempty"`
	PaymentMethods  []PaymentMethod            `db:"-" json:"payment_methods,omitempty"`
	DeliveryOptions []StorefrontDeliveryOption `db:"-" json:"delivery_options,omitempty"`
}

// IsDeleted returns true if the storefront is soft-deleted
func (s *Storefront) IsDeleted() bool {
	return s.DeletedAt != nil
}

// HasLocation returns true if the storefront has coordinates set
func (s *Storefront) HasLocation() bool {
	return s.Latitude != nil && s.Longitude != nil
}

// GetDisplayName returns the name for display purposes
func (s *Storefront) GetDisplayName() string {
	return s.Name
}

// GetURL returns the storefront URL based on slug
func (s *Storefront) GetURL(baseURL string) string {
	return baseURL + "/store/" + s.Slug
}

// StorefrontStaff represents staff member (storefront_staff table)
type StorefrontStaff struct {
	ID           int64      `db:"id" json:"id"`
	StorefrontID int64      `db:"storefront_id" json:"storefront_id"`
	UserID       int64      `db:"user_id" json:"user_id"`
	Role         string     `db:"role" json:"role"`
	Permissions  JSONB      `db:"permissions" json:"permissions,omitempty"`
	LastActiveAt *time.Time `db:"last_active_at" json:"last_active_at,omitempty"`
	ActionsCount int32      `db:"actions_count" json:"actions_count"`
	InvitationID *int64     `db:"invitation_id" json:"invitation_id,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// StorefrontHours represents working hours (storefront_hours table)
type StorefrontHours struct {
	ID           int64   `db:"id" json:"id"`
	StorefrontID int64   `db:"storefront_id" json:"storefront_id"`
	DayOfWeek    int32   `db:"day_of_week" json:"day_of_week"`
	OpenTime     *string `db:"open_time" json:"open_time,omitempty"`
	CloseTime    *string `db:"close_time" json:"close_time,omitempty"`
	IsClosed     bool    `db:"is_closed" json:"is_closed"`
	SpecialDate  *string `db:"special_date" json:"special_date,omitempty"`
	SpecialNote  *string `db:"special_note" json:"special_note,omitempty"`
}

// PaymentMethod represents payment method (storefront_payment_methods table)
type PaymentMethod struct {
	ID             int64     `db:"id" json:"id"`
	StorefrontID   int64     `db:"storefront_id" json:"storefront_id"`
	MethodType     string    `db:"method_type" json:"method_type"`
	IsEnabled      bool      `db:"is_enabled" json:"is_enabled"`
	Provider       *string   `db:"provider" json:"provider,omitempty"`
	Settings       JSONB     `db:"settings" json:"settings,omitempty"`
	TransactionFee float64   `db:"transaction_fee" json:"transaction_fee"`
	MinAmount      *float64  `db:"min_amount" json:"min_amount,omitempty"`
	MaxAmount      *float64  `db:"max_amount" json:"max_amount,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// StorefrontDeliveryOption represents delivery option (storefront_delivery_options table)
type StorefrontDeliveryOption struct {
	ID               int64     `db:"id" json:"id"`
	StorefrontID     int64     `db:"storefront_id" json:"storefront_id"`
	Name             string    `db:"name" json:"name"`
	Description      *string   `db:"description" json:"description,omitempty"`
	BasePrice        float64   `db:"base_price" json:"base_price"`
	PricePerKm       float64   `db:"price_per_km" json:"price_per_km"`
	PricePerKg       float64   `db:"price_per_kg" json:"price_per_kg"`
	FreeAboveAmount  *float64  `db:"free_above_amount" json:"free_above_amount,omitempty"`
	MinOrderAmount   *float64  `db:"min_order_amount" json:"min_order_amount,omitempty"`
	MaxWeightKg      *float64  `db:"max_weight_kg" json:"max_weight_kg,omitempty"`
	MaxDistanceKm    *float64  `db:"max_distance_km" json:"max_distance_km,omitempty"`
	EstimatedDaysMin int32     `db:"estimated_days_min" json:"estimated_days_min"`
	EstimatedDaysMax int32     `db:"estimated_days_max" json:"estimated_days_max"`
	Zones            JSONB     `db:"zones" json:"zones,omitempty"`
	AvailableDays    JSONB     `db:"available_days" json:"available_days,omitempty"`
	CutoffTime       *string   `db:"cutoff_time" json:"cutoff_time,omitempty"`
	Provider         *string   `db:"provider" json:"provider,omitempty"`
	ProviderConfig   JSONB     `db:"provider_config" json:"provider_config,omitempty"`
	IsActive         bool      `db:"is_active" json:"is_active"`
	DisplayOrder     int32     `db:"display_order" json:"display_order"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// StorefrontUpdate represents fields that can be updated
type StorefrontUpdate struct {
	Name                *string  `json:"name,omitempty"`
	Description         *string  `json:"description,omitempty"`
	IsActive            *bool    `json:"is_active,omitempty"`
	LogoURL             *string  `json:"logo_url,omitempty"`
	BannerURL           *string  `json:"banner_url,omitempty"`
	Theme               JSONB    `json:"theme,omitempty"`
	Phone               *string  `json:"phone,omitempty"`
	Email               *string  `json:"email,omitempty"`
	Website             *string  `json:"website,omitempty"`
	Address             *string  `json:"address,omitempty"`
	City                *string  `json:"city,omitempty"`
	PostalCode          *string  `json:"postal_code,omitempty"`
	Country             *string  `json:"country,omitempty"`
	Latitude            *float64 `json:"latitude,omitempty"`
	Longitude           *float64 `json:"longitude,omitempty"`
	FormattedAddress    *string  `json:"formatted_address,omitempty"`
	Settings            JSONB    `json:"settings,omitempty"`
	SeoMeta             JSONB    `json:"seo_meta,omitempty"`
	AIAgentEnabled      *bool    `json:"ai_agent_enabled,omitempty"`
	LiveShoppingEnabled *bool    `json:"live_shopping_enabled,omitempty"`
	GroupBuyingEnabled  *bool    `json:"group_buying_enabled,omitempty"`
}

// StaffUpdate represents staff fields that can be updated
type StaffUpdate struct {
	Role        *string `json:"role,omitempty"`
	Permissions JSONB   `json:"permissions,omitempty"`
}

// ListStorefrontsFilter represents filter parameters for listing storefronts
type ListStorefrontsFilter struct {
	UserID            *int64   `json:"user_id,omitempty"`
	IsActive          *bool    `json:"is_active,omitempty"`
	IsVerified        *bool    `json:"is_verified,omitempty"`
	City              *string  `json:"city,omitempty"`
	Country           *string  `json:"country,omitempty"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
	RadiusKm          *float64 `json:"radius_km,omitempty"`
	MinRating         *float64 `json:"min_rating,omitempty"`
	SubscriptionPlans []string `json:"subscription_plans,omitempty"`
	PaymentMethods    []string `json:"payment_methods,omitempty"`
	DeliveryProviders []string `json:"delivery_providers,omitempty"`
	SupportsCOD       *bool    `json:"supports_cod,omitempty"`
	HasAIAgent        *bool    `json:"has_ai_agent,omitempty"`
	HasLiveShopping   *bool    `json:"has_live_shopping,omitempty"`
	HasGroupBuying    *bool    `json:"has_group_buying,omitempty"`
	Search            *string  `json:"search,omitempty"`
	SortBy            string   `json:"sort_by"`
	SortOrder         string   `json:"sort_order"`
	Page              int32    `json:"page"`
	Limit             int32    `json:"limit"`
}

// MapBounds represents geographical bounds for map queries
type MapBounds struct {
	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// StorefrontMapData represents minimal data for map display
type StorefrontMapData struct {
	ID            int64   `db:"id" json:"id"`
	Slug          string  `db:"slug" json:"slug"`
	Name          string  `db:"name" json:"name"`
	Latitude      float64 `db:"latitude" json:"latitude"`
	Longitude     float64 `db:"longitude" json:"longitude"`
	Rating        float64 `db:"rating" json:"rating"`
	LogoURL       *string `db:"logo_url" json:"logo_url,omitempty"`
	Address       *string `db:"address" json:"address,omitempty"`
	Phone         *string `db:"phone" json:"phone,omitempty"`
	WorkingNow    bool    `db:"working_now" json:"working_now"`
	ProductsCount int32   `db:"products_count" json:"products_count"`
	SupportsCOD   bool    `db:"supports_cod" json:"supports_cod"`
	HasDelivery   bool    `db:"has_delivery" json:"has_delivery"`
	HasSelfPickup bool    `db:"has_self_pickup" json:"has_self_pickup"`
	AcceptsCards  bool    `db:"accepts_cards" json:"accepts_cards"`
}

// StorefrontDashboardStats represents storefront dashboard statistics
type StorefrontDashboardStats struct {
	TotalProducts    int32   `json:"total_products"`
	ActiveProducts   int32   `json:"active_products"`
	OrdersCount      int32   `json:"orders_count"`
	Revenue          float64 `json:"revenue"`
	AvgOrderValue    float64 `json:"avg_order_value"`
	ViewsCount       int32   `json:"views_count"`
	UniqueVisitors   int32   `json:"unique_visitors"`
	ConversionRate   float64 `json:"conversion_rate"`
	PendingOrders    int32   `json:"pending_orders"`
	LowStockProducts int32   `json:"low_stock_products"`
}

// Includes represents which related entities to include
type Includes struct {
	Staff           bool
	Hours           bool
	PaymentMethods  bool
	DeliveryOptions bool
}
