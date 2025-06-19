package service

import (
	"backend/internal/storage/postgres"
	"context"
	"encoding/json"
)

// AnalyticsService интерфейс сервиса аналитики
type AnalyticsService interface {
	RecordEvent(ctx context.Context, event *EventData) error
}

// EventData данные события
type EventData struct {
	StorefrontID int             `json:"storefront_id"`
	EventType    string          `json:"event_type"`
	EventData    json.RawMessage `json:"event_data"`
	SessionID    string          `json:"session_id"`
	UserID       *int            `json:"user_id,omitempty"`
	IPAddress    string          `json:"ip_address"`
	UserAgent    string          `json:"user_agent"`
	Referrer     string          `json:"referrer"`
}

// analyticsServiceImpl реализация сервиса
type analyticsServiceImpl struct {
	storefrontRepo postgres.StorefrontRepository
}

// NewAnalyticsService создает новый сервис аналитики
func NewAnalyticsService(storefrontRepo postgres.StorefrontRepository) AnalyticsService {
	return &analyticsServiceImpl{
		storefrontRepo: storefrontRepo,
	}
}

// RecordEvent записывает событие
func (s *analyticsServiceImpl) RecordEvent(ctx context.Context, event *EventData) error {
	// Конвертируем в формат репозитория
	repoEvent := &postgres.StorefrontEvent{
		StorefrontID: event.StorefrontID,
		EventType:    postgres.EventType(event.EventType),
		EventData:    event.EventData,
		SessionID:    event.SessionID,
		UserID:       event.UserID,
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		Referrer:     event.Referrer,
	}

	return s.storefrontRepo.RecordEvent(ctx, repoEvent)
}