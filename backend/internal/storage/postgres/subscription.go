package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
)

// SubscriptionRepository represents subscription repository
type SubscriptionRepository struct {
	db *sqlx.DB
}

// NewSubscriptionRepository creates new subscription repository
func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// GetPlans returns all active subscription plans
func (r *SubscriptionRepository) GetPlans(ctx context.Context) ([]models.SubscriptionPlanDetails, error) {
	query := `
		SELECT id, code, name, COALESCE(price_monthly, 0) as price_monthly, COALESCE(price_yearly, 0) as price_yearly,
			max_storefronts, max_products_per_storefront, max_staff_per_storefront, max_images_total,
			has_ai_assistant, has_live_shopping, has_export_data, has_custom_domain,
			has_analytics, has_priority_support, commission_rate, free_trial_days,
			sort_order, is_active, is_recommended, created_at, updated_at
		FROM subscription_plans
		WHERE is_active = true
		ORDER BY sort_order ASC`

	var plans []models.SubscriptionPlanDetails
	err := r.db.SelectContext(ctx, &plans, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plans: %w", err)
	}

	return plans, nil
}

// GetPlanByCode returns plan by code
func (r *SubscriptionRepository) GetPlanByCode(ctx context.Context, code string) (*models.SubscriptionPlanDetails, error) {
	query := `
		SELECT id, code, name, COALESCE(price_monthly, 0) as price_monthly, COALESCE(price_yearly, 0) as price_yearly,
			max_storefronts, max_products_per_storefront, max_staff_per_storefront, max_images_total,
			has_ai_assistant, has_live_shopping, has_export_data, has_custom_domain,
			has_analytics, has_priority_support, commission_rate, free_trial_days,
			sort_order, is_active, is_recommended, created_at, updated_at
		FROM subscription_plans
		WHERE code = $1 AND is_active = true`

	var plan models.SubscriptionPlanDetails
	err := r.db.GetContext(ctx, &plan, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return &plan, nil
}

// GetUserSubscription returns active user subscription
func (r *SubscriptionRepository) GetUserSubscription(ctx context.Context, userID int) (*models.UserSubscription, error) {
	query := `
		SELECT us.id, us.user_id, us.plan_id, us.status, us.billing_cycle,
			us.started_at, us.trial_ends_at, us.current_period_start, us.current_period_end,
			us.canceled_at, us.expires_at, us.last_payment_id, us.last_payment_at,
			us.next_payment_at, us.payment_method, us.auto_renew, us.used_storefronts,
			us.metadata, us.notes, us.created_at, us.updated_at
		FROM user_subscriptions us
		WHERE us.user_id = $1 AND us.status IN ('active', 'trial')
		LIMIT 1`

	var subscription models.UserSubscription
	err := r.db.GetContext(ctx, &subscription, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return empty subscription, not nil to avoid nilnil error
			return &models.UserSubscription{}, nil
		}
		return nil, fmt.Errorf("failed to get user subscription: %w", err)
	}

	// Load plan details
	plan, err := r.GetPlanByID(ctx, subscription.PlanID)
	if err == nil {
		subscription.Plan = plan
	}

	return &subscription, nil
}

// GetPlanByID returns plan by ID
func (r *SubscriptionRepository) GetPlanByID(ctx context.Context, id int) (*models.SubscriptionPlanDetails, error) {
	query := `
		SELECT id, code, name, COALESCE(price_monthly, 0) as price_monthly, COALESCE(price_yearly, 0) as price_yearly,
			max_storefronts, max_products_per_storefront, max_staff_per_storefront, max_images_total,
			has_ai_assistant, has_live_shopping, has_export_data, has_custom_domain,
			has_analytics, has_priority_support, commission_rate, free_trial_days,
			sort_order, is_active, is_recommended, created_at, updated_at
		FROM subscription_plans
		WHERE id = $1`

	var plan models.SubscriptionPlanDetails
	err := r.db.GetContext(ctx, &plan, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan by ID: %w", err)
	}

	return &plan, nil
}

// CreateSubscription creates new subscription
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.UserSubscription, error) {
	// Get plan details
	plan, err := r.GetPlanByCode(ctx, req.PlanCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Check if user already has subscription
	var existingID int
	checkQuery := `SELECT id FROM user_subscriptions WHERE user_id = $1 AND status IN ('active', 'trial')`
	err = tx.GetContext(ctx, &existingID, checkQuery, req.UserID)
	if err == nil {
		return nil, fmt.Errorf("user already has active subscription")
	}

	// Calculate dates
	now := time.Now()
	var trialEndsAt *time.Time
	var currentPeriodEnd time.Time
	var status models.SubscriptionStatus

	if req.StartTrial && plan.FreeTrialDays > 0 {
		trialEnd := now.AddDate(0, 0, plan.FreeTrialDays)
		trialEndsAt = &trialEnd
		currentPeriodEnd = trialEnd
		status = models.SubscriptionStatusTrial
	} else {
		status = models.SubscriptionStatusActive
		if req.BillingCycle == models.BillingCycleMonthly {
			currentPeriodEnd = now.AddDate(0, 1, 0)
		} else {
			currentPeriodEnd = now.AddDate(1, 0, 0)
		}
	}

	// Insert subscription
	insertQuery := `
		INSERT INTO user_subscriptions (
			user_id, plan_id, status, billing_cycle,
			started_at, trial_ends_at, current_period_start, current_period_end,
			next_payment_at, payment_method, auto_renew, used_storefronts
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 0
		) RETURNING id`

	var subscriptionID int
	err = tx.GetContext(ctx, &subscriptionID, insertQuery,
		req.UserID, plan.ID, status, req.BillingCycle,
		now, trialEndsAt, now, currentPeriodEnd,
		&currentPeriodEnd, req.PaymentMethod, true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Add history record
	historyQuery := `
		INSERT INTO subscription_history (subscription_id, user_id, action, to_plan_id)
		VALUES ($1, $2, $3, $4)`

	_, err = tx.ExecContext(ctx, historyQuery, subscriptionID, req.UserID, models.SubscriptionActionCreated, plan.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to add history: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get created subscription
	return r.GetUserSubscription(ctx, req.UserID)
}

// UpgradeSubscription upgrades existing subscription
func (r *SubscriptionRepository) UpgradeSubscription(ctx context.Context, userID int, req *models.UpgradeSubscriptionRequest) (*models.UserSubscription, error) {
	// Get current subscription
	currentSub, err := r.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current subscription: %w", err)
	}
	if currentSub == nil {
		return nil, fmt.Errorf("no active subscription found")
	}

	// Get new plan
	newPlan, err := r.GetPlanByCode(ctx, req.PlanCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get new plan: %w", err)
	}

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Update subscription
	billingCycle := currentSub.BillingCycle
	if req.BillingCycle != "" {
		billingCycle = req.BillingCycle
	}

	updateQuery := `
		UPDATE user_subscriptions
		SET plan_id = $1, billing_cycle = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	_, err = tx.ExecContext(ctx, updateQuery, newPlan.ID, billingCycle, currentSub.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Add history record
	var action models.SubscriptionAction
	if newPlan.SortOrder > currentSub.Plan.SortOrder {
		action = models.SubscriptionActionUpgraded
	} else {
		action = models.SubscriptionActionDowngraded
	}

	historyQuery := `
		INSERT INTO subscription_history (subscription_id, user_id, action, from_plan_id, to_plan_id)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, historyQuery, currentSub.ID, userID, action, currentSub.PlanID, newPlan.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to add history: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get updated subscription
	return r.GetUserSubscription(ctx, userID)
}

// CancelSubscription cancels subscription
func (r *SubscriptionRepository) CancelSubscription(ctx context.Context, userID int, reason string) error {
	// Get current subscription
	currentSub, err := r.GetUserSubscription(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get current subscription: %w", err)
	}
	if currentSub == nil {
		return fmt.Errorf("no active subscription found")
	}

	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Update subscription
	now := time.Now()
	updateQuery := `
		UPDATE user_subscriptions
		SET status = $1, canceled_at = $2, expires_at = $3, auto_renew = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	_, err = tx.ExecContext(ctx, updateQuery, models.SubscriptionStatusCanceled, now, currentSub.CurrentPeriodEnd, currentSub.ID)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	// Add history record
	historyQuery := `
		INSERT INTO subscription_history (subscription_id, user_id, action, from_plan_id, reason)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, historyQuery, currentSub.ID, userID, models.SubscriptionActionCanceled, currentSub.PlanID, reason)
	if err != nil {
		return fmt.Errorf("failed to add history: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CheckSubscriptionLimits checks if user can use resource
func (r *SubscriptionRepository) CheckSubscriptionLimits(ctx context.Context, userID int, resourceType string, count int) (*models.CheckLimitResponse, error) {
	// Get user subscription or use free plan
	sub, err := r.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	var plan *models.SubscriptionPlanDetails
	var usedStorefronts int

	if sub != nil {
		plan = sub.Plan
		usedStorefronts = sub.UsedStorefronts
	} else {
		// Use free plan
		freePlan, err := r.GetPlanByCode(ctx, "starter")
		if err != nil {
			return nil, fmt.Errorf("failed to get free plan: %w", err)
		}
		plan = freePlan

		// Count current storefronts for user without subscription
		countQuery := `SELECT COUNT(*) FROM storefronts WHERE user_id = $1 AND is_active = true`
		err = r.db.GetContext(ctx, &usedStorefronts, countQuery, userID)
		if err != nil {
			usedStorefronts = 0
		}
	}

	response := &models.CheckLimitResponse{
		ResourceType: resourceType,
	}

	switch resourceType {
	case "storefront":
		response.CurrentUsage = usedStorefronts
		response.Limit = plan.MaxStorefronts
		if plan.MaxStorefronts == -1 {
			response.Allowed = true
			response.Message = "Unlimited storefronts allowed"
		} else {
			response.Allowed = (usedStorefronts + count) <= plan.MaxStorefronts
			if !response.Allowed {
				response.Message = fmt.Sprintf("Storefront limit reached (%d/%d)", usedStorefronts, plan.MaxStorefronts)
				response.RequiredPlan = r.getRequiredPlanForStorefronts(usedStorefronts + count)
			}
		}

	case "product":
		response.Limit = plan.MaxProductsPerStorefront
		if plan.MaxProductsPerStorefront == -1 {
			response.Allowed = true
			response.Message = "Unlimited products allowed"
		} else {
			// TODO: Get current product count for storefront
			response.Allowed = true
		}

	case "staff":
		response.Limit = plan.MaxStaffPerStorefront
		if plan.MaxStaffPerStorefront == -1 {
			response.Allowed = true
			response.Message = "Unlimited staff allowed"
		} else {
			// TODO: Get current staff count for storefront
			response.Allowed = true
		}

	case "image":
		response.Limit = plan.MaxImagesTotal
		if plan.MaxImagesTotal == -1 {
			response.Allowed = true
			response.Message = "Unlimited images allowed"
		} else {
			// TODO: Get current image count
			response.Allowed = true
		}

	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	return response, nil
}

// getRequiredPlanForStorefronts returns minimum plan for storefront count
func (r *SubscriptionRepository) getRequiredPlanForStorefronts(count int) string {
	switch {
	case count <= 1:
		return "starter"
	case count <= 3:
		return "professional"
	case count <= 10:
		return "business"
	default:
		return "enterprise"
	}
}

// GetUserSubscriptionInfo returns detailed subscription info
func (r *SubscriptionRepository) GetUserSubscriptionInfo(ctx context.Context, userID int) (*models.UserSubscriptionInfo, error) {
	query := `SELECT * FROM get_user_subscription($1)`

	var info models.UserSubscriptionInfo
	err := r.db.GetContext(ctx, &info, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No subscription - return free plan info
			freePlan, err := r.GetPlanByCode(ctx, "starter")
			if err != nil {
				return nil, fmt.Errorf("failed to get free plan: %w", err)
			}

			// Count current storefronts
			var usedStorefronts int
			countQuery := `SELECT COUNT(*) FROM storefronts WHERE user_id = $1 AND is_active = true`
			_ = r.db.GetContext(ctx, &usedStorefronts, countQuery, userID)

			return &models.UserSubscriptionInfo{
				PlanCode:        &freePlan.Code,
				PlanName:        &freePlan.Name,
				MaxStorefronts:  freePlan.MaxStorefronts,
				UsedStorefronts: usedStorefronts,
				MaxProducts:     freePlan.MaxProductsPerStorefront,
				MaxStaff:        freePlan.MaxStaffPerStorefront,
				MaxImages:       freePlan.MaxImagesTotal,
				HasAI:           freePlan.HasAIAssistant,
				HasLive:         freePlan.HasLiveShopping,
				HasExport:       freePlan.HasExportData,
				HasCustomDomain: freePlan.HasCustomDomain,
				Price:           decimal.Zero,
			}, nil
		}
		return nil, fmt.Errorf("failed to get subscription info: %w", err)
	}

	// Add price based on billing cycle
	if info.BillingCycle != nil && info.PlanCode != nil {
		plan, err := r.GetPlanByCode(ctx, *info.PlanCode)
		if err == nil {
			if *info.BillingCycle == string(models.BillingCycleMonthly) {
				info.Price = plan.PriceMonthly
			} else {
				info.Price = plan.PriceYearly
			}
		}
	}

	return &info, nil
}

// GetDB returns database connection
func (r *SubscriptionRepository) GetDB() *sqlx.DB {
	return r.db
}

// RecordPayment records subscription payment
func (r *SubscriptionRepository) RecordPayment(ctx context.Context, subscriptionID int, paymentID int, amount decimal.Decimal, status models.SubscriptionPaymentStatus) error {
	// Get subscription
	var sub models.UserSubscription
	query := `SELECT user_id, current_period_start, current_period_end FROM user_subscriptions WHERE id = $1`
	err := r.db.GetContext(ctx, &sub, query, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Insert payment record
	insertQuery := `
		INSERT INTO subscription_payments (
			subscription_id, user_id, payment_id, amount, currency,
			period_start, period_end, status, paid_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var paidAt *time.Time
	if status == models.SubPaymentStatusCompleted {
		now := time.Now()
		paidAt = &now
	}

	_, err = r.db.ExecContext(ctx, insertQuery,
		subscriptionID, sub.UserID, paymentID, amount, "EUR",
		sub.CurrentPeriodStart, sub.CurrentPeriodEnd, status, paidAt,
	)
	if err != nil {
		return fmt.Errorf("failed to record payment: %w", err)
	}

	// Update subscription if payment completed
	if status == models.SubPaymentStatusCompleted {
		updateQuery := `
			UPDATE user_subscriptions
			SET last_payment_id = $1, last_payment_at = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3`

		_, err = r.db.ExecContext(ctx, updateQuery, paymentID, paidAt, subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
	}

	return nil
}
