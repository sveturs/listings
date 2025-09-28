package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/proj/subscriptions/service"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

// SubscriptionHandler handles subscription endpoints
type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *logger.Logger
}

// NewSubscriptionHandler creates new subscription handler
func NewSubscriptionHandler(service *service.SubscriptionService, logger *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// GetPlans godoc
// @Summary Get subscription plans
// @Description Get all available subscription plans
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_proj_subscriptions_models.SubscriptionPlanDetails} "List of plans"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/plans [get]
func (h *SubscriptionHandler) GetPlans(c *fiber.Ctx) error {
	plans, err := h.service.GetPlans(c.Context())
	if err != nil {
		h.logger.Error("Failed to get plans: %v", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.getPlansError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.getPlansSuccess", plans)
}

// GetCurrentSubscription godoc
// @Summary Get current subscription
// @Description Get current user's subscription details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserSubscriptionInfo} "Current subscription"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/current [get]
func (h *SubscriptionHandler) GetCurrentSubscription(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	subscription, err := h.service.GetUserSubscription(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get subscription: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.getSubscriptionError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.getSubscriptionSuccess", subscription)
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create new subscription for user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body backend_internal_proj_subscriptions_models.CreateSubscriptionRequest true "Create subscription request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserSubscription} "Created subscription"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 409 {object} utils.ErrorResponseSwag "User already has subscription"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	var req models.CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "common.invalidRequest")
	}

	// Ensure user can only create subscription for themselves
	req.UserID = userID

	subscription, err := h.service.CreateSubscription(c.Context(), &req)
	if err != nil {
		if err.Error() == "user already has active subscription" {
			return utils.SendError(c, fiber.StatusConflict, "subscriptions.alreadyHasSubscription")
		}
		h.logger.Error("Failed to create subscription: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.createError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.createSuccess", subscription)
}

// UpgradeSubscription godoc
// @Summary Upgrade subscription
// @Description Upgrade existing subscription to a new plan
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body backend_internal_proj_subscriptions_models.UpgradeSubscriptionRequest true "Upgrade subscription request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserSubscription} "Upgraded subscription"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Subscription not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/upgrade [post]
func (h *SubscriptionHandler) UpgradeSubscription(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	var req models.UpgradeSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "common.invalidRequest")
	}

	subscription, err := h.service.UpgradeSubscription(c.Context(), userID, &req)
	if err != nil {
		if err.Error() == "no active subscription found" {
			return utils.SendError(c, fiber.StatusNotFound, "subscriptions.notFound")
		}
		h.logger.Error("Failed to upgrade subscription: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.upgradeError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.upgradeSuccess", subscription)
}

// CancelSubscription godoc
// @Summary Cancel subscription
// @Description Cancel current subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param reason body map[string]string false "Cancellation reason"
// @Success 200 {object} utils.SuccessResponseSwag{message=string} "Subscription canceled"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Subscription not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/cancel [post]
func (h *SubscriptionHandler) CancelSubscription(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		body = map[string]string{}
	}

	reason := body["reason"]
	if reason == "" {
		reason = "User requested cancellation"
	}

	err := h.service.CancelSubscription(c.Context(), userID, reason)
	if err != nil {
		if err.Error() == "no active subscription found" {
			return utils.SendError(c, fiber.StatusNotFound, "subscriptions.notFound")
		}
		h.logger.Error("Failed to cancel subscription: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.cancelError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.cancelSuccess", nil)
}

// CheckLimits godoc
// @Summary Check subscription limits
// @Description Check if user can use a resource within subscription limits
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body backend_internal_proj_subscriptions_models.CheckLimitRequest true "Check limit request"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_proj_subscriptions_models.CheckLimitResponse} "Limit check result"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/check-limits [post]
func (h *SubscriptionHandler) CheckLimits(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	var req models.CheckLimitRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "common.invalidRequest")
	}

	response, err := h.service.CheckLimits(c.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to check limits: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.checkLimitsError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.checkLimitsSuccess", response)
}

// InitiatePayment godoc
// @Summary Initiate payment for subscription
// @Description Initiate payment process for subscription plan
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body InitiatePaymentRequest true "Payment initiation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=service.PaymentInitiationResponse} "Payment initiation response"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/initiate-payment [post]
func (h *SubscriptionHandler) InitiatePayment(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64
		userID64, ok := userIDVal.(int64)
		if !ok {
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}

	var req InitiatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "common.invalidRequest")
	}

	// Default billing cycle if not specified
	if req.BillingCycle == "" {
		req.BillingCycle = models.BillingCycleMonthly
	}

	// Generate return URL
	returnURL := req.ReturnURL
	if returnURL == "" {
		returnURL = "https://svetu.rs/subscription/payment-complete"
	}

	response, err := h.service.InitiatePayment(c.Context(), userID, req.PlanCode, req.BillingCycle, returnURL)
	if err != nil {
		h.logger.Error("Failed to initiate payment: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.paymentInitiationError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.paymentInitiationSuccess", response)
}

// InitiatePaymentRequest represents payment initiation request
type InitiatePaymentRequest struct {
	PlanCode     string              `json:"plan_code" validate:"required"`
	BillingCycle models.BillingCycle `json:"billing_cycle,omitempty"`
	ReturnURL    string              `json:"return_url,omitempty"`
}

// CompletePayment godoc
// @Summary Complete payment for subscription
// @Description Complete payment process after returning from payment gateway
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param payment_intent query string true "Payment intent ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserSubscription} "Payment completed"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/subscriptions/complete-payment [post]
func (h *SubscriptionHandler) CompletePayment(c *fiber.Ctx) error {
	// Safe type assertion for userId
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "users.auth.error.unauthorized")
	}

	userID, ok := userIDVal.(int)
	if !ok {
		// Try int64 as fallback
		userID64, ok := userIDVal.(int64)
		if !ok {
			h.logger.Error("Invalid userId type in context: %v", userIDVal)
			return utils.SendError(c, fiber.StatusInternalServerError, "common.internalError")
		}
		userID = int(userID64)
	}
	paymentIntentID := c.Query("payment_intent")

	if paymentIntentID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "subscriptions.paymentIntentRequired")
	}

	// Get current subscription
	sub, err := h.service.GetUserSubscription(c.Context(), userID)
	if err != nil || sub == nil || sub.SubscriptionID == nil {
		return utils.SendError(c, fiber.StatusNotFound, "subscriptions.notFound")
	}

	// Process payment
	err = h.service.ProcessSubscriptionPayment(c.Context(), userID, *sub.SubscriptionID, paymentIntentID)
	if err != nil {
		h.logger.Error("Failed to complete payment: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.paymentError")
	}

	// Get updated subscription
	updatedSub, _ := h.service.GetUserSubscription(c.Context(), userID)

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.paymentSuccess", updatedSub)
}

// AdminGetUserSubscription godoc
// @Summary Get user subscription (Admin)
// @Description Get subscription details for a specific user (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param user_id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserSubscriptionInfo} "User subscription"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "User not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/users/{user_id}/subscription [get]
func (h *SubscriptionHandler) AdminGetUserSubscription(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "common.invalidUserID")
	}

	subscription, err := h.service.GetUserSubscription(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user subscription: %v (user_id: %d)", err, userID)
		return utils.SendError(c, fiber.StatusInternalServerError, "subscriptions.getSubscriptionError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "subscriptions.getSubscriptionSuccess", subscription)
}
