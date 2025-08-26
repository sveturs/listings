package models

import (
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
)

// SubscriptionStatus represents subscription status
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusTrial     SubscriptionStatus = "trial"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
)

// BillingCycle represents billing cycle
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

// SubscriptionAction represents action type for history
type SubscriptionAction string

const (
	SubscriptionActionCreated    SubscriptionAction = "created"
	SubscriptionActionUpgraded   SubscriptionAction = "upgraded"
	SubscriptionActionDowngraded SubscriptionAction = "downgraded"
	SubscriptionActionRenewed    SubscriptionAction = "renewed"
	SubscriptionActionCanceled   SubscriptionAction = "canceled"
	SubscriptionActionExpired    SubscriptionAction = "expired"
)

// SubscriptionPaymentStatus for subscription payments
type SubscriptionPaymentStatus string

const (
	SubPaymentStatusPending    SubscriptionPaymentStatus = "pending"
	SubPaymentStatusProcessing SubscriptionPaymentStatus = "processing"
	SubPaymentStatusCompleted  SubscriptionPaymentStatus = "completed"
	SubPaymentStatusFailed     SubscriptionPaymentStatus = "failed"
	SubPaymentStatusRefunded   SubscriptionPaymentStatus = "refunded"
)

// SubscriptionPlanDetails represents a subscription plan with full details
type SubscriptionPlanDetails struct {
	ID                       int             `json:"id" db:"id"`
	Code                     string          `json:"code" db:"code"`
	Name                     string          `json:"name" db:"name"`
	PriceMonthly             decimal.Decimal `json:"price_monthly" db:"price_monthly"`
	PriceYearly              decimal.Decimal `json:"price_yearly" db:"price_yearly"`
	MaxStorefronts           int             `json:"max_storefronts" db:"max_storefronts"`
	MaxProductsPerStorefront int             `json:"max_products_per_storefront" db:"max_products_per_storefront"`
	MaxStaffPerStorefront    int             `json:"max_staff_per_storefront" db:"max_staff_per_storefront"`
	MaxImagesTotal           int             `json:"max_images_total" db:"max_images_total"`
	HasAIAssistant           bool            `json:"has_ai_assistant" db:"has_ai_assistant"`
	HasLiveShopping          bool            `json:"has_live_shopping" db:"has_live_shopping"`
	HasExportData            bool            `json:"has_export_data" db:"has_export_data"`
	HasCustomDomain          bool            `json:"has_custom_domain" db:"has_custom_domain"`
	HasAnalytics             bool            `json:"has_analytics" db:"has_analytics"`
	HasPrioritySupport       bool            `json:"has_priority_support" db:"has_priority_support"`
	CommissionRate           decimal.Decimal `json:"commission_rate" db:"commission_rate"`
	FreeTrialDays            int             `json:"free_trial_days" db:"free_trial_days"`
	SortOrder                int             `json:"sort_order" db:"sort_order"`
	IsActive                 bool            `json:"is_active" db:"is_active"`
	IsRecommended            bool            `json:"is_recommended" db:"is_recommended"`
	CreatedAt                time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at" db:"updated_at"`
}

// UserSubscription represents a user's subscription
type UserSubscription struct {
	ID                 int                      `json:"id" db:"id"`
	UserID             int                      `json:"user_id" db:"user_id"`
	PlanID             int                      `json:"plan_id" db:"plan_id"`
	Status             SubscriptionStatus       `json:"status" db:"status"`
	BillingCycle       BillingCycle             `json:"billing_cycle" db:"billing_cycle"`
	StartedAt          time.Time                `json:"started_at" db:"started_at"`
	TrialEndsAt        *time.Time               `json:"trial_ends_at,omitempty" db:"trial_ends_at"`
	CurrentPeriodStart time.Time                `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time                `json:"current_period_end" db:"current_period_end"`
	CanceledAt         *time.Time               `json:"canceled_at,omitempty" db:"canceled_at"`
	ExpiresAt          *time.Time               `json:"expires_at,omitempty" db:"expires_at"`
	LastPaymentID      *int                     `json:"last_payment_id,omitempty" db:"last_payment_id"`
	LastPaymentAt      *time.Time               `json:"last_payment_at,omitempty" db:"last_payment_at"`
	NextPaymentAt      *time.Time               `json:"next_payment_at,omitempty" db:"next_payment_at"`
	PaymentMethod      *string                  `json:"payment_method,omitempty" db:"payment_method"`
	AutoRenew          bool                     `json:"auto_renew" db:"auto_renew"`
	UsedStorefronts    int                      `json:"used_storefronts" db:"used_storefronts"`
	Metadata           JSONB                    `json:"metadata" db:"metadata"`
	Notes              *string                  `json:"notes,omitempty" db:"notes"`
	CreatedAt          time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at" db:"updated_at"`
	Plan               *SubscriptionPlanDetails `json:"plan,omitempty" db:"-"`
}

// SubscriptionHistory represents subscription history
type SubscriptionHistory struct {
	ID             int                `json:"id" db:"id"`
	SubscriptionID int                `json:"subscription_id" db:"subscription_id"`
	UserID         int                `json:"user_id" db:"user_id"`
	Action         SubscriptionAction `json:"action" db:"action"`
	FromPlanID     *int               `json:"from_plan_id,omitempty" db:"from_plan_id"`
	ToPlanID       *int               `json:"to_plan_id,omitempty" db:"to_plan_id"`
	Reason         *string            `json:"reason,omitempty" db:"reason"`
	Metadata       JSONB              `json:"metadata" db:"metadata"`
	CreatedAt      time.Time          `json:"created_at" db:"created_at"`
	CreatedBy      *int               `json:"created_by,omitempty" db:"created_by"`
}

// SubscriptionPayment represents a subscription payment
type SubscriptionPayment struct {
	ID              int                       `json:"id" db:"id"`
	SubscriptionID  int                       `json:"subscription_id" db:"subscription_id"`
	UserID          int                       `json:"user_id" db:"user_id"`
	PaymentID       *int                      `json:"payment_id,omitempty" db:"payment_id"`
	Amount          decimal.Decimal           `json:"amount" db:"amount"`
	Currency        string                    `json:"currency" db:"currency"`
	PeriodStart     time.Time                 `json:"period_start" db:"period_start"`
	PeriodEnd       time.Time                 `json:"period_end" db:"period_end"`
	Status          SubscriptionPaymentStatus `json:"status" db:"status"`
	PaymentMethod   *string                   `json:"payment_method,omitempty" db:"payment_method"`
	TransactionData JSONB                     `json:"transaction_data" db:"transaction_data"`
	PaidAt          *time.Time                `json:"paid_at,omitempty" db:"paid_at"`
	FailedAt        *time.Time                `json:"failed_at,omitempty" db:"failed_at"`
	RefundedAt      *time.Time                `json:"refunded_at,omitempty" db:"refunded_at"`
	CreatedAt       time.Time                 `json:"created_at" db:"created_at"`
}

// SubscriptionUsage represents resource usage
type SubscriptionUsage struct {
	ID             int       `json:"id" db:"id"`
	SubscriptionID int       `json:"subscription_id" db:"subscription_id"`
	StorefrontID   *int      `json:"storefront_id,omitempty" db:"storefront_id"`
	ResourceType   string    `json:"resource_type" db:"resource_type"`
	ResourceID     *int      `json:"resource_id,omitempty" db:"resource_id"`
	ResourceCount  int       `json:"resource_count" db:"resource_count"`
	Action         string    `json:"action" db:"action"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// UserSubscriptionInfo represents complete subscription info
type UserSubscriptionInfo struct {
	SubscriptionID  *int            `json:"subscription_id,omitempty" db:"subscription_id"`
	PlanCode        *string         `json:"plan_code,omitempty" db:"plan_code"`
	PlanName        *string         `json:"plan_name,omitempty" db:"plan_name"`
	Status          *string         `json:"status,omitempty" db:"status"`
	ExpiresAt       *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	MaxStorefronts  int             `json:"max_storefronts" db:"max_storefronts"`
	UsedStorefronts int             `json:"used_storefronts" db:"used_storefronts"`
	MaxProducts     int             `json:"max_products" db:"max_products"`
	MaxStaff        int             `json:"max_staff" db:"max_staff"`
	MaxImages       int             `json:"max_images" db:"max_images"`
	HasAI           bool            `json:"has_ai" db:"has_ai"`
	HasLive         bool            `json:"has_live" db:"has_live"`
	HasExport       bool            `json:"has_export" db:"has_export"`
	HasCustomDomain bool            `json:"has_custom_domain" db:"has_custom_domain"`
	Price           decimal.Decimal `json:"price,omitempty"`
	BillingCycle    *string         `json:"billing_cycle,omitempty" db:"billing_cycle"`
	NextPaymentAt   *time.Time      `json:"next_payment_at,omitempty" db:"next_payment_at"`
	TrialEndsAt     *time.Time      `json:"trial_ends_at,omitempty" db:"trial_ends_at"`
}

// CreateSubscriptionRequest represents request to create subscription
type CreateSubscriptionRequest struct {
	UserID        int          `json:"user_id" validate:"required"`
	PlanCode      string       `json:"plan_code" validate:"required"`
	BillingCycle  BillingCycle `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	PaymentMethod string       `json:"payment_method" validate:"required"`
	StartTrial    bool         `json:"start_trial"`
}

// UpgradeSubscriptionRequest represents request to upgrade subscription
type UpgradeSubscriptionRequest struct {
	PlanCode     string       `json:"plan_code" validate:"required"`
	BillingCycle BillingCycle `json:"billing_cycle,omitempty" validate:"omitempty,oneof=monthly yearly"`
}

// CheckLimitRequest represents request to check limits
type CheckLimitRequest struct {
	ResourceType string `json:"resource_type" validate:"required,oneof=storefront product staff image"`
	Count        int    `json:"count,omitempty"`
}

// CheckLimitResponse represents response for limit check
type CheckLimitResponse struct {
	Allowed      bool   `json:"allowed"`
	CurrentUsage int    `json:"current_usage"`
	Limit        int    `json:"limit"`
	ResourceType string `json:"resource_type"`
	RequiredPlan string `json:"required_plan,omitempty"`
	Message      string `json:"message,omitempty"`
}

// Value implements driver.Valuer interface
func (s SubscriptionStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Value implements driver.Valuer interface
func (b BillingCycle) Value() (driver.Value, error) {
	return string(b), nil
}

// Value implements driver.Valuer interface
func (a SubscriptionAction) Value() (driver.Value, error) {
	return string(a), nil
}

// Value implements driver.Valuer interface
func (p SubscriptionPaymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}
