package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// SubscriptionPlanType тип тарифного плана
type SubscriptionPlanType string

const (
	SubscriptionPlanStarter      SubscriptionPlanType = "starter"
	SubscriptionPlanProfessional SubscriptionPlanType = "professional"
	SubscriptionPlanBusiness     SubscriptionPlanType = "business"
	SubscriptionPlanEnterprise   SubscriptionPlanType = "enterprise"
)

// StaffRole роль сотрудника витрины
type StaffRole string

const (
	StaffRoleOwner     StaffRole = "owner"
	StaffRoleManager   StaffRole = "manager"
	StaffRoleCashier   StaffRole = "cashier"
	StaffRoleSupport   StaffRole = "support"
	StaffRoleModerator StaffRole = "moderator"
)

// PaymentMethodType тип платежного метода
type PaymentMethodType string

const (
	PaymentMethodCash         PaymentMethodType = "cash"          // Наличные в магазине
	PaymentMethodCOD          PaymentMethodType = "cod"           // Cash on Delivery - оплата курьеру
	PaymentMethodCard         PaymentMethodType = "card"          // Банковская карта
	PaymentMethodBankTransfer PaymentMethodType = "bank_transfer" // Банковский перевод
	PaymentMethodPayPal       PaymentMethodType = "paypal"
	PaymentMethodCrypto       PaymentMethodType = "crypto"
	PaymentMethodPostanska    PaymentMethodType = "postanska" // Poštanska štedionica
	PaymentMethodKeks         PaymentMethodType = "keks_pay"  // Keks Pay (популярно в Сербии)
	PaymentMethodIPS          PaymentMethodType = "ips"       // Instant Payment System
)

// DeliveryProvider провайдеры доставки
type DeliveryProvider string

const (
	DeliveryProviderPostaSrbije DeliveryProvider = "posta_srbije"
	DeliveryProviderAKS         DeliveryProvider = "aks"
	DeliveryProviderBEX         DeliveryProvider = "bex"
	DeliveryProviderDExpress    DeliveryProvider = "d_express"
	DeliveryProviderCityExpress DeliveryProvider = "city_express"
	DeliveryProviderSelfPickup  DeliveryProvider = "self_pickup"
	DeliveryProviderOwnDelivery DeliveryProvider = "own_delivery"
)

// StorefrontGeoStrategy стратегия геолокации витрины
type StorefrontGeoStrategy string

const (
	GeoStrategyStorefrontLocation StorefrontGeoStrategy = "storefront_location" // Использовать адрес витрины
	GeoStrategyIndividualLocation StorefrontGeoStrategy = "individual_location" // Использовать индивидуальные адреса товаров
)

// LocationPrivacyLevel уровень приватности адреса
type LocationPrivacyLevel string

const (
	PrivacyLevelExact    LocationPrivacyLevel = "exact"    // Точный адрес
	PrivacyLevelStreet   LocationPrivacyLevel = "street"   // Только улица
	PrivacyLevelDistrict LocationPrivacyLevel = "district" // Только район
	PrivacyLevelCity     LocationPrivacyLevel = "city"     // Только город
)

// JSONB тип для работы с JSONB полями PostgreSQL
type JSONB map[string]interface{}

// Value реализует интерфейс driver.Valuer
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil // nolint:nilnil // Valid pattern for driver.Valuer when nil is expected
	}
	return json.Marshal(j)
}

// Scan реализует интерфейс sql.Scanner
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}
	return json.Unmarshal(bytes, j)
}

// Storefront представляет структуру витрины
type Storefront struct {
	ID          int    `json:"id" db:"id"`
	UserID      int    `json:"user_id" db:"user_id"`
	Slug        string `json:"slug" db:"slug"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`

	// Связанные данные
	User *User `json:"user,omitempty"` // Владелец витрины

	// Брендинг
	LogoURL   string `json:"logo_url,omitempty" db:"logo_url"`
	BannerURL string `json:"banner_url,omitempty" db:"banner_url"`
	Theme     JSONB  `json:"theme" db:"theme"`

	// Контактная информация
	Phone   string `json:"phone,omitempty" db:"phone"`
	Email   string `json:"email,omitempty" db:"email"`
	Website string `json:"website,omitempty" db:"website"`

	// Локация
	Address             string                `json:"address,omitempty" db:"address"`
	City                string                `json:"city,omitempty" db:"city"`
	PostalCode          string                `json:"postal_code,omitempty" db:"postal_code"`
	Country             string                `json:"country" db:"country"`
	Latitude            *float64              `json:"latitude,omitempty" db:"latitude"`
	Longitude           *float64              `json:"longitude,omitempty" db:"longitude"`
	FormattedAddress    *string               `json:"formatted_address,omitempty" db:"formatted_address"`
	GeoStrategy         StorefrontGeoStrategy `json:"geo_strategy" db:"geo_strategy"`
	DefaultPrivacyLevel LocationPrivacyLevel  `json:"default_privacy_level" db:"default_privacy_level"`
	AddressVerified     bool                  `json:"address_verified" db:"address_verified"`

	// Настройки бизнеса
	Settings JSONB `json:"settings" db:"settings"`
	SEOMeta  JSONB `json:"seo_meta" db:"seo_meta"`

	// Статус и статистика
	IsActive         bool       `json:"is_active" db:"is_active"`
	IsVerified       bool       `json:"is_verified" db:"is_verified"`
	VerificationDate *time.Time `json:"verification_date,omitempty" db:"verification_date"`
	Rating           float64    `json:"rating" db:"rating"`
	ReviewsCount     int        `json:"reviews_count" db:"reviews_count"`
	ProductsCount    int        `json:"products_count" db:"products_count"`
	SalesCount       int        `json:"sales_count" db:"sales_count"`
	ViewsCount       int        `json:"views_count" db:"views_count"`

	// Подписка (монетизация)
	SubscriptionPlan      SubscriptionPlanType `json:"subscription_plan" db:"subscription_plan"`
	SubscriptionExpiresAt *time.Time           `json:"subscription_expires_at,omitempty" db:"subscription_expires_at"`
	CommissionRate        float64              `json:"commission_rate" db:"commission_rate"`

	// AI и killer features
	AIAgentEnabled      bool  `json:"ai_agent_enabled" db:"ai_agent_enabled"`
	AIAgentConfig       JSONB `json:"ai_agent_config" db:"ai_agent_config"`
	LiveShoppingEnabled bool  `json:"live_shopping_enabled" db:"live_shopping_enabled"`
	GroupBuyingEnabled  bool  `json:"group_buying_enabled" db:"group_buying_enabled"`

	// Временные метки
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// StorefrontStaff представляет сотрудника витрины
type StorefrontStaff struct {
	ID           int       `json:"id"`
	StorefrontID int       `json:"storefront_id"`
	UserID       int       `json:"user_id"`
	Role         StaffRole `json:"role"`
	Permissions  JSONB     `json:"permissions"`

	// Отслеживание активности
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	ActionsCount int        `json:"actions_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StorefrontHours представляет часы работы витрины
type StorefrontHours struct {
	ID           int    `json:"id"`
	StorefrontID int    `json:"storefront_id"`
	DayOfWeek    int    `json:"day_of_week"` // 0=Sunday, 6=Saturday
	OpenTime     string `json:"open_time,omitempty"`
	CloseTime    string `json:"close_time,omitempty"`
	IsClosed     bool   `json:"is_closed"`

	// Специальные часы
	SpecialDate *time.Time `json:"special_date,omitempty"`
	SpecialNote string     `json:"special_note,omitempty"`
}

// StorefrontPaymentMethod представляет метод оплаты витрины
type StorefrontPaymentMethod struct {
	ID           int               `json:"id"`
	StorefrontID int               `json:"storefront_id"`
	MethodType   PaymentMethodType `json:"method_type"`
	IsEnabled    bool              `json:"is_enabled"`

	// Настройки провайдера
	Provider string `json:"provider,omitempty"`
	Settings JSONB  `json:"settings"`

	// Комиссии и лимиты
	TransactionFee float64  `json:"transaction_fee"`
	MinAmount      *float64 `json:"min_amount,omitempty"`
	MaxAmount      *float64 `json:"max_amount,omitempty"`

	// Специфично для COD (наложенный платеж)
	CODFee *float64 `json:"cod_fee,omitempty"` // Дополнительная плата за наложенный платеж

	CreatedAt time.Time `json:"created_at"`
}

// DeliveryZone представляет зону доставки
type DeliveryZone struct {
	Name          string   `json:"name"`
	PostalCodes   []string `json:"postal_codes"`
	PriceModifier float64  `json:"price_modifier"`
}

// StorefrontDeliveryOption представляет опцию доставки витрины
type StorefrontDeliveryOption struct {
	ID           int              `json:"id"`
	StorefrontID int              `json:"storefront_id"`
	Name         string           `json:"name"`
	Description  string           `json:"description,omitempty"`
	Provider     DeliveryProvider `json:"provider"`

	// Ценообразование
	BasePrice       float64  `json:"base_price"`
	PricePerKm      float64  `json:"price_per_km"`
	PricePerKg      float64  `json:"price_per_kg"`
	FreeAboveAmount *float64 `json:"free_above_amount,omitempty"`

	// Специальные платы
	CODFee          *float64 `json:"cod_fee,omitempty"`          // Плата за наложенный платеж
	InsuranceFee    *float64 `json:"insurance_fee,omitempty"`    // Страховка
	FragileHandling *float64 `json:"fragile_handling,omitempty"` // Обработка хрупких товаров

	// Ограничения доставки
	MinOrderAmount   *float64 `json:"min_order_amount,omitempty"`
	MaxWeightKg      *float64 `json:"max_weight_kg,omitempty"`
	MaxDistanceKm    *float64 `json:"max_distance_km,omitempty"`
	EstimatedDaysMin int      `json:"estimated_days_min"`
	EstimatedDaysMax int      `json:"estimated_days_max"`

	// Зоны и доступность
	Zones         []DeliveryZone `json:"zones"`
	AvailableDays []int          `json:"available_days"`
	CutoffTime    string         `json:"cutoff_time,omitempty"`

	// Поддерживаемые методы оплаты для этой доставки
	SupportedPaymentMethods []PaymentMethodType `json:"supported_payment_methods"`

	// Интеграция с провайдером
	ProviderConfig JSONB `json:"provider_config"`

	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StorefrontAnalytics представляет аналитику витрины
type StorefrontAnalytics struct {
	ID           int       `json:"id"`
	StorefrontID int       `json:"storefront_id"`
	Date         time.Time `json:"date"`

	// Трафик
	PageViews      int     `json:"page_views"`
	UniqueVisitors int     `json:"unique_visitors"`
	BounceRate     float64 `json:"bounce_rate"`
	AvgSessionTime int     `json:"avg_session_time"` // в секундах

	// Продажи
	OrdersCount    int     `json:"orders_count"`
	Revenue        float64 `json:"revenue"`
	AvgOrderValue  float64 `json:"avg_order_value"`
	ConversionRate float64 `json:"conversion_rate"`

	// Методы оплаты
	PaymentMethodsUsage JSONB `json:"payment_methods_usage"` // {cod: 45%, card: 30%, ...}

	// Товары
	ProductViews   int `json:"product_views"`
	AddToCartCount int `json:"add_to_cart_count"`
	CheckoutCount  int `json:"checkout_count"`

	// Источники трафика
	TrafficSources JSONB `json:"traffic_sources"`

	// Топ товары/категории
	TopProducts   JSONB `json:"top_products"`
	TopCategories JSONB `json:"top_categories"`

	// География заказов
	OrdersByCity JSONB `json:"orders_by_city"`

	CreatedAt time.Time `json:"created_at"`
}

// Location представляет расширенную информацию о локации
type Location struct {
	// Координаты клика пользователя
	UserLat float64 `json:"user_lat"`
	UserLng float64 `json:"user_lng"`

	// Координаты здания (после геокодинга)
	BuildingLat float64 `json:"building_lat"`
	BuildingLng float64 `json:"building_lng"`

	// Адресные данные
	FullAddress string `json:"full_address"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	PostalCode  string `json:"postal_code"`
	City        string `json:"city"`
	Country     string `json:"country"`

	// Дополнительная информация
	BuildingInfo JSONB `json:"building_info,omitempty"`
}

// DTO для создания/обновления

// StorefrontCreateDTO данные для создания витрины
type StorefrontCreateDTO struct {
	UserID      int    `json:"-"` // Заполняется из контекста
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Slug        string `json:"slug,omitempty"` // Генерируется автоматически если пустой
	Description string `json:"description,omitempty"`

	// Брендинг
	Logo   []byte `json:"-"`
	Banner []byte `json:"-"`
	Theme  JSONB  `json:"theme,omitempty"`

	// Контактная информация
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty" validate:"omitempty,email"`
	Website string `json:"website,omitempty" validate:"omitempty,url"`

	// Локация
	Location Location `json:"location" validate:"required"`

	// Настройки
	Settings JSONB `json:"settings,omitempty"`
	SEOMeta  JSONB `json:"seo_meta,omitempty"`
}

// StorefrontUpdateDTO данные для обновления витрины
type StorefrontUpdateDTO struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`

	// Брендинг
	LogoURL   *string `json:"logo_url,omitempty"`
	BannerURL *string `json:"banner_url,omitempty"`
	Theme     JSONB   `json:"theme,omitempty"`

	// Контактная информация
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Website *string `json:"website,omitempty" validate:"omitempty,url"`

	// Локация
	Location *Location `json:"location,omitempty"`

	// Настройки
	Settings JSONB `json:"settings,omitempty"`
	SEOMeta  JSONB `json:"seo_meta,omitempty"`

	// Функции
	AIAgentEnabled      *bool `json:"ai_agent_enabled,omitempty"`
	LiveShoppingEnabled *bool `json:"live_shopping_enabled,omitempty"`
	GroupBuyingEnabled  *bool `json:"group_buying_enabled,omitempty"`
}

// StorefrontFilter фильтры для поиска витрин
type StorefrontFilter struct {
	// Базовые фильтры
	UserID     *int    `json:"user_id,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
	City       *string `json:"city,omitempty"`
	Country    *string `json:"country,omitempty"`

	// Геолокация
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	RadiusKm  *float64 `json:"radius_km,omitempty"`

	// Рейтинг
	MinRating *float64 `json:"min_rating,omitempty"`

	// Тарифный план
	SubscriptionPlans []SubscriptionPlanType `json:"subscription_plans,omitempty"`

	// Методы оплаты
	PaymentMethods []PaymentMethodType `json:"payment_methods,omitempty"`

	// Методы доставки
	DeliveryProviders []DeliveryProvider `json:"delivery_providers,omitempty"`

	// Поддержка наложенного платежа
	SupportsCOD *bool `json:"supports_cod,omitempty"`

	// Функции
	HasAIAgent      *bool `json:"has_ai_agent,omitempty"`
	HasLiveShopping *bool `json:"has_live_shopping,omitempty"`
	HasGroupBuying  *bool `json:"has_group_buying,omitempty"`

	// Поиск
	Search *string `json:"search,omitempty"`

	// Сортировка
	SortBy    string `json:"sort_by,omitempty"`    // rating, created_at, products_count, distance
	SortOrder string `json:"sort_order,omitempty"` // asc, desc

	// Пагинация
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// StorefrontMapData данные витрины для отображения на карте
type StorefrontMapData struct {
	ID        int     `json:"id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    float64 `json:"rating"`
	LogoURL   string  `json:"logo_url,omitempty"`

	// Краткая информация для попапа
	Address       string `json:"address"`
	Phone         string `json:"phone,omitempty"`
	WorkingNow    bool   `json:"working_now"`
	ProductsCount int    `json:"products_count"`

	// Ключевые фичи
	SupportsCOD   bool `json:"supports_cod"`
	HasDelivery   bool `json:"has_delivery"`
	HasSelfPickup bool `json:"has_self_pickup"`
	AcceptsCards  bool `json:"accepts_cards"`
}

// StaffPermissions разрешения для сотрудников
type StaffPermissions struct {
	// Управление товарами
	CanAddProducts    bool `json:"can_add_products"`
	CanEditProducts   bool `json:"can_edit_products"`
	CanDeleteProducts bool `json:"can_delete_products"`

	// Управление заказами
	CanViewOrders    bool `json:"can_view_orders"`
	CanProcessOrders bool `json:"can_process_orders"`
	CanRefundOrders  bool `json:"can_refund_orders"`

	// Управление витриной
	CanEditStorefront bool `json:"can_edit_storefront"`
	CanManageStaff    bool `json:"can_manage_staff"`
	CanViewAnalytics  bool `json:"can_view_analytics"`
	CanManagePayments bool `json:"can_manage_payments"`

	// Коммуникации
	CanReplyToReviews bool `json:"can_reply_to_reviews"`
	CanSendMessages   bool `json:"can_send_messages"`
}
