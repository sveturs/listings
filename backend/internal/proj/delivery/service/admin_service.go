package service

import (
	"context"

	"backend/internal/proj/delivery/models"
)

// Admin methods for delivery service

// GetAllProviders returns all delivery providers
func (s *Service) GetAllProviders(ctx context.Context) ([]models.Provider, error) {
	return s.storage.GetAllProviders(ctx)
}

// ToggleProviderStatus toggles provider active status
func (s *Service) ToggleProviderStatus(ctx context.Context, providerID int) error {
	return s.storage.ToggleProviderStatus(ctx, providerID)
}

// UpdateProviderConfig updates provider API configuration
func (s *Service) UpdateProviderConfig(ctx context.Context, providerID int, config models.ProviderConfig) error {
	return s.storage.UpdateProviderConfig(ctx, providerID, config)
}

// GetAllPricingRules returns all pricing rules
func (s *Service) GetAllPricingRules(ctx context.Context) ([]models.PricingRule, error) {
	return s.storage.GetAllPricingRules(ctx)
}

// UpdatePricingRule updates an existing pricing rule
func (s *Service) UpdatePricingRule(ctx context.Context, rule models.PricingRule) error {
	return s.storage.UpdatePricingRule(ctx, rule)
}

// DeletePricingRule deletes a pricing rule
func (s *Service) DeletePricingRule(ctx context.Context, ruleID int) error {
	return s.storage.DeletePricingRule(ctx, ruleID)
}

// GetProblemShipments returns problem shipments based on filters
func (s *Service) GetProblemShipments(ctx context.Context, problemType, status string) ([]models.ProblemShipment, error) {
	return s.storage.GetProblemShipments(ctx, problemType, status)
}

// AssignProblem assigns a problem to an admin
func (s *Service) AssignProblem(ctx context.Context, problemID, adminID int) error {
	return s.storage.AssignProblem(ctx, problemID, adminID)
}

// ResolveProblem resolves a problem shipment
func (s *Service) ResolveProblem(ctx context.Context, problemID int, resolution string) error {
	return s.storage.ResolveProblem(ctx, problemID, resolution)
}

// GetDashboardStats returns dashboard statistics
func (s *Service) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	return s.storage.GetDashboardStats(ctx)
}

// GetAnalytics returns analytics data for a specific period
func (s *Service) GetAnalytics(ctx context.Context, period string) (*models.AnalyticsData, error) {
	return s.storage.GetAnalytics(ctx, period)
}
