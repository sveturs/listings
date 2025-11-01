// backend/internal/proj/marketplace/service/router_adapter.go
package service

import (
	unifiedService "backend/internal/proj/unified/service"
)

// RouterAdapter адаптирует marketplace.TrafficRouter к unified.TrafficRouter interface
type RouterAdapter struct {
	router *TrafficRouter
}

// NewRouterAdapter создаёт адаптер для TrafficRouter
func NewRouterAdapter(router *TrafficRouter) *RouterAdapter {
	return &RouterAdapter{router: router}
}

// ShouldUseMicroservice адаптирует вызов к TrafficRouter
func (a *RouterAdapter) ShouldUseMicroservice(userID string, isAdmin bool) *unifiedService.TrafficRoutingDecision {
	decision := a.router.ShouldUseMicroservice(userID, isAdmin)

	// Конвертируем RoutingDecision -> TrafficRoutingDecision
	return &unifiedService.TrafficRoutingDecision{
		UseМicroservice: decision.UseМicroservice,
		Reason:          decision.Reason,
		UserID:          decision.UserID,
		IsAdmin:         decision.IsAdmin,
		IsCanary:        decision.IsCanary,
		Hash:            decision.Hash,
	}
}

// ValidateConfig делегирует валидацию к TrafficRouter
func (a *RouterAdapter) ValidateConfig() error {
	return a.router.ValidateConfig()
}
