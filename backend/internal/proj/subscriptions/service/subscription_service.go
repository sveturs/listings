package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	"backend/internal/pkg/allsecure"
	"backend/internal/proj/payments/service"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"
)

// SubscriptionService manages subscriptions
type SubscriptionService struct {
	repo            *postgres.SubscriptionRepository
	paymentService  *service.AllSecureService
	allSecureClient *allsecure.Client
	logger          *logger.Logger
	useMockPayments bool
}

// NewSubscriptionService creates new subscription service
func NewSubscriptionService(
	repo *postgres.SubscriptionRepository,
	paymentService *service.AllSecureService,
	allSecureClient *allsecure.Client,
	logger *logger.Logger,
) *SubscriptionService {
	useMock := paymentService == nil
	if useMock {
		logger.Info("Subscription service initialized with mock payments")
	}

	return &SubscriptionService{
		repo:            repo,
		paymentService:  paymentService,
		allSecureClient: allSecureClient,
		logger:          logger,
		useMockPayments: useMock,
	}
}

// GetPlans returns all active subscription plans
func (s *SubscriptionService) GetPlans(ctx context.Context) ([]models.SubscriptionPlanDetails, error) {
	plans, err := s.repo.GetPlans(ctx)
	if err != nil {
		s.logger.Error("Failed to get plans: %v", err)
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}
	return plans, nil
}

// GetUserSubscription returns user's active subscription
func (s *SubscriptionService) GetUserSubscription(ctx context.Context, userID int) (*models.UserSubscriptionInfo, error) {
	info, err := s.repo.GetUserSubscriptionInfo(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user subscription: %v (user_id: %d)", err, userID)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return info, nil
}

// CreateSubscription creates new subscription
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.UserSubscription, error) {
	// Validate request
	if req.UserID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Check if user already has subscription
	existing, err := s.repo.GetUserSubscription(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user already has active subscription")
	}

	// Get plan details
	plan, err := s.repo.GetPlanByCode(ctx, req.PlanCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// For free plan, create immediately
	if plan.PriceMonthly.IsZero() && plan.PriceYearly.IsZero() {
		subscription, err := s.repo.CreateSubscription(ctx, req)
		if err != nil {
			s.logger.Error("Failed to create free subscription: %v (user_id: %d)", err, req.UserID)
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
		s.logger.Info("Created free subscription for user %d with plan %s", req.UserID, req.PlanCode)
		return subscription, nil
	}

	// For paid plans, initiate payment first
	amount := plan.PriceMonthly
	if req.BillingCycle == models.BillingCycleYearly {
		amount = plan.PriceYearly
	}

	// If trial period, create subscription without payment
	if req.StartTrial && plan.FreeTrialDays > 0 {
		subscription, err := s.repo.CreateSubscription(ctx, req)
		if err != nil {
			s.logger.Error("Failed to create trial subscription: %v (user_id: %d)", err, req.UserID)
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
		s.logger.Info("Created trial subscription for user %d with plan %s (trial days: %d)", req.UserID, req.PlanCode, plan.FreeTrialDays)
		return subscription, nil
	}

	// Create subscription (payment will be handled separately)
	subscription, err := s.repo.CreateSubscription(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create subscription: %v (user_id: %d)", err, req.UserID)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info("Created subscription for user %d with plan %s (amount: %s)", req.UserID, req.PlanCode, amount.String())
	return subscription, nil
}

// UpgradeSubscription upgrades existing subscription
func (s *SubscriptionService) UpgradeSubscription(ctx context.Context, userID int, req *models.UpgradeSubscriptionRequest) (*models.UserSubscription, error) {
	// Get current subscription
	current, err := s.repo.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current subscription: %w", err)
	}
	if current == nil {
		// No subscription - create new one
		createReq := &models.CreateSubscriptionRequest{
			UserID:        userID,
			PlanCode:      req.PlanCode,
			BillingCycle:  req.BillingCycle,
			PaymentMethod: "card", // Default to card
		}
		if req.BillingCycle == "" {
			createReq.BillingCycle = models.BillingCycleMonthly
		}
		return s.CreateSubscription(ctx, createReq)
	}

	// Check if downgrade or upgrade
	newPlan, err := s.repo.GetPlanByCode(ctx, req.PlanCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get new plan: %w", err)
	}

	// Upgrade subscription
	subscription, err := s.repo.UpgradeSubscription(ctx, userID, req)
	if err != nil {
		s.logger.Error("Failed to upgrade subscription: %v (user_id: %d)", err, userID)
		return nil, fmt.Errorf("failed to upgrade subscription: %w", err)
	}

	s.logger.Info("Upgraded subscription for user %d from %s to %s", userID, current.Plan.Code, newPlan.Code)
	return subscription, nil
}

// CancelSubscription cancels user's subscription
func (s *SubscriptionService) CancelSubscription(ctx context.Context, userID int, reason string) error {
	err := s.repo.CancelSubscription(ctx, userID, reason)
	if err != nil {
		s.logger.Error("Failed to cancel subscription: %v (user_id: %d)", err, userID)
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	s.logger.Info("Canceled subscription for user %d (reason: %s)", userID, reason)
	return nil
}

// CheckLimits checks if user can use resource
func (s *SubscriptionService) CheckLimits(ctx context.Context, userID int, req *models.CheckLimitRequest) (*models.CheckLimitResponse, error) {
	count := req.Count
	if count == 0 {
		count = 1
	}

	response, err := s.repo.CheckSubscriptionLimits(ctx, userID, req.ResourceType, count)
	if err != nil {
		s.logger.Error("Failed to check limits: %v (user_id: %d, resource: %s)", err, userID, req.ResourceType)
		return nil, fmt.Errorf("failed to check limits: %w", err)
	}

	return response, nil
}

// ProcessSubscriptionPayment processes payment for subscription
func (s *SubscriptionService) ProcessSubscriptionPayment(ctx context.Context, userID int, subscriptionID int, paymentIntentID string) error {
	// Get subscription details
	sub, err := s.repo.GetUserSubscription(ctx, userID)
	if err != nil || sub == nil {
		return fmt.Errorf("subscription not found")
	}

	// Get plan details
	plan := sub.Plan
	if plan == nil {
		return fmt.Errorf("plan not found")
	}

	// Calculate amount
	amount := plan.PriceMonthly
	if sub.BillingCycle == models.BillingCycleYearly {
		amount = plan.PriceYearly
	}

	// Record payment
	err = s.repo.RecordPayment(ctx, subscriptionID, 0, amount, models.SubPaymentStatusCompleted)
	if err != nil {
		s.logger.Error("Failed to record payment: %v (subscription_id: %d)", err, subscriptionID)
		return fmt.Errorf("failed to record payment: %w", err)
	}

	s.logger.Info("Processed subscription payment for subscription %d (amount: %s)", subscriptionID, amount.String())
	return nil
}

// InitiatePayment initiates payment for subscription
func (s *SubscriptionService) InitiatePayment(ctx context.Context, userID int, planCode string, billingCycle models.BillingCycle, returnURL string) (*PaymentInitiationResponse, error) {
	// Get plan details
	plan, err := s.repo.GetPlanByCode(ctx, planCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Calculate amount
	amount := plan.PriceMonthly
	if billingCycle == models.BillingCycleYearly {
		amount = plan.PriceYearly
	}

	// Skip payment for free plan
	if amount.IsZero() {
		return &PaymentInitiationResponse{
			PaymentRequired: false,
			Message:         "Free plan selected",
		}, nil
	}

	// Use mock payments if AllSecure is not configured
	if s.useMockPayments {
		return InitiatePaymentMock(ctx, userID, planCode, billingCycle, amount, returnURL)
	}

	// Since AllSecure requires a listing, we'll create subscription payments differently
	// For now, return mock payment URL since AllSecure isn't fully integrated for subscriptions
	paymentIntentID := fmt.Sprintf("pi_sub_%d_%s_%d", userID, planCode, time.Now().Unix())
	redirectURL := fmt.Sprintf("%s/subscription/payment-mock?payment_intent=%s&amount=%s&plan=%s&cycle=%s",
		returnURL,
		paymentIntentID,
		amount.String(),
		planCode,
		billingCycle,
	)

	return &PaymentInitiationResponse{
		PaymentRequired: true,
		PaymentIntentID: paymentIntentID,
		RedirectURL:     redirectURL,
		Amount:          amount,
		Currency:        "EUR",
	}, nil
}

// PaymentInitiationResponse represents payment initiation response
type PaymentInitiationResponse struct {
	PaymentRequired bool            `json:"payment_required"`
	PaymentIntentID string          `json:"payment_intent_id,omitempty"`
	RedirectURL     string          `json:"redirect_url,omitempty"`
	Amount          decimal.Decimal `json:"amount,omitempty"`
	Currency        string          `json:"currency,omitempty"`
	Message         string          `json:"message,omitempty"`
}

// RenewSubscription renews expired subscription
func (s *SubscriptionService) RenewSubscription(ctx context.Context, userID int) error {
	// Get current subscription
	sub, err := s.repo.GetUserSubscription(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub == nil {
		return fmt.Errorf("no subscription found")
	}

	// Check if renewal needed
	if sub.Status != models.SubscriptionStatusExpired {
		return fmt.Errorf("subscription is not expired")
	}

	// TODO: Process renewal payment and update subscription
	s.logger.Info("Subscription renewal initiated for user %d", userID)
	return nil
}

// ProcessExpiredSubscriptions processes all expired subscriptions
func (s *SubscriptionService) ProcessExpiredSubscriptions(ctx context.Context) error {
	// This would be called by a cron job
	query := `
		UPDATE user_subscriptions
		SET status = 'expired', updated_at = CURRENT_TIMESTAMP
		WHERE status IN ('active', 'trial')
		AND current_period_end < $1
		AND auto_renew = false`

	result, err := s.repo.GetDB().ExecContext(ctx, query, time.Now())
	if err != nil {
		s.logger.Error("Failed to process expired subscriptions: %v", err)
		return fmt.Errorf("failed to process expired subscriptions: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected > 0 {
		s.logger.Info("Processed %d expired subscriptions", affected)
	}

	return nil
}
