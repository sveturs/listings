package translation_admin

import (
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetCostsSummary godoc
// @Summary Get AI translation costs summary
// @Description Returns summary of AI translation costs across all providers
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/costs [get]
func (h *AITranslationHandler) GetCostsSummary(c *fiber.Ctx) error {
	ctx := c.Context()

	costs, err := h.service.GetAIProviderCosts(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get costs summary")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.costsRetrievalFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.costsRetrieved", costs)
}

// GetCostAlerts godoc
// @Summary Get AI translation cost alerts
// @Description Returns alerts if costs exceed specified limits
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param daily_limit query float64 false "Daily cost limit in USD" default(100)
// @Param monthly_limit query float64 false "Monthly cost limit in USD" default(2000)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]string}
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/costs/alerts [get]
func (h *AITranslationHandler) GetCostAlerts(c *fiber.Ctx) error {
	ctx := c.Context()

	dailyLimit := c.QueryFloat("daily_limit", 100.0)
	monthlyLimit := c.QueryFloat("monthly_limit", 2000.0)

	alerts, err := h.service.GetAIProviderAlerts(ctx, dailyLimit, monthlyLimit)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get cost alerts")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.alertsRetrievalFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.alertsRetrieved", map[string]interface{}{
		"alerts":        alerts,
		"daily_limit":   dailyLimit,
		"monthly_limit": monthlyLimit,
		"has_alerts":    len(alerts) > 0,
	})
}

// ResetProviderCosts godoc
// @Summary Reset AI provider costs
// @Description Resets cost tracking for a specific provider
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param provider path string true "Provider name (openai, google, deepl, claude)"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]string}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/costs/{provider}/reset [post]
func (h *AITranslationHandler) ResetProviderCosts(c *fiber.Ctx) error {
	ctx := c.Context()

	provider := c.Params("provider")
	if provider == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.providerRequired")
	}

	// Validate provider
	validProviders := map[string]bool{
		"openai": true,
		"google": true,
		"deepl":  true,
		"claude": true,
	}

	if !validProviders[provider] {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.invalidProvider")
	}

	err := h.service.ResetAIProviderCosts(ctx, provider)
	if err != nil {
		h.logger.Error().Err(err).Str("provider", provider).Msg("Failed to reset provider costs")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.costsResetFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.costsReset", map[string]string{
		"provider": provider,
		"status":   "reset",
	})
}

// GetProviderCostDetails godoc
// @Summary Get detailed costs for a specific provider
// @Description Returns detailed cost information for a specific AI provider
// @Tags Translation Admin
// @Accept json
// @Produce json
// @Param provider path string true "Provider name (openai, google, deepl, claude)"
// @Success 200 {object} utils.SuccessResponseSwag{data=ProviderCosts}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/admin/translations/ai/costs/{provider} [get]
func (h *AITranslationHandler) GetProviderCostDetails(c *fiber.Ctx) error {
	ctx := c.Context()

	provider := c.Params("provider")
	if provider == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "admin.translations.providerRequired")
	}

	// Get cost tracker from service
	costTracker := h.service.GetCostTracker()
	if costTracker == nil {
		return utils.SendError(c, fiber.StatusServiceUnavailable, "admin.translations.costTrackingUnavailable")
	}

	costs, err := costTracker.GetProviderCosts(ctx, provider)
	if err != nil {
		h.logger.Error().Err(err).Str("provider", provider).Msg("Failed to get provider costs")
		return utils.SendError(c, fiber.StatusInternalServerError, "admin.translations.providerCostsRetrievalFailed")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "admin.translations.providerCostsRetrieved", costs)
}
