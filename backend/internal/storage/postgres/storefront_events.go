package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// EventType тип события
type EventType string

const (
	EventPageView    EventType = "page_view"
	EventProductView EventType = "product_view"
	EventAddToCart   EventType = "add_to_cart"
	EventCheckout    EventType = "checkout"
	EventOrder       EventType = "order"
)

// StorefrontEvent событие витрины
type StorefrontEvent struct {
	StorefrontID int             `json:"storefront_id"`
	EventType    EventType       `json:"event_type"`
	EventData    json.RawMessage `json:"event_data"`
	UserID       *int            `json:"user_id,omitempty"`
	SessionID    string          `json:"session_id"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	Referrer     string          `json:"referrer,omitempty"`
}

// RecordEvent записывает событие в базу данных
func (r *storefrontRepo) RecordEvent(ctx context.Context, event *StorefrontEvent) error {
	query := `
		INSERT INTO storefront_events (
			storefront_id, event_type, event_data, user_id, session_id,
			ip_address, user_agent, referrer
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.pool.Exec(ctx, query,
		event.StorefrontID, event.EventType, event.EventData,
		event.UserID, event.SessionID, event.IPAddress,
		event.UserAgent, event.Referrer,
	)
	
	if err != nil {
		return fmt.Errorf("failed to record event: %w", err)
	}
	
	return nil
}

// GetEventStats получает статистику событий за период
func (r *storefrontRepo) GetEventStats(ctx context.Context, storefrontID int, from, to time.Time) (map[EventType]int, error) {
	query := `
		SELECT event_type, COUNT(*) as count
		FROM storefront_events
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
		GROUP BY event_type
	`
	
	rows, err := r.db.pool.Query(ctx, query, storefrontID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get event stats: %w", err)
	}
	defer rows.Close()
	
	stats := make(map[EventType]int)
	for rows.Next() {
		var eventType EventType
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan event stat: %w", err)
		}
		stats[eventType] = count
	}
	
	return stats, nil
}

// GetUniqueVisitors получает количество уникальных посетителей
func (r *storefrontRepo) GetUniqueVisitors(ctx context.Context, storefrontID int, from, to time.Time) (int, error) {
	query := `
		SELECT COUNT(DISTINCT session_id)
		FROM storefront_events
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
			AND event_type = 'page_view'
	`
	
	var count int
	err := r.db.pool.QueryRow(ctx, query, storefrontID, from, to).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unique visitors: %w", err)
	}
	
	return count, nil
}

// GetTrafficSources получает источники трафика
func (r *storefrontRepo) GetTrafficSources(ctx context.Context, storefrontID int, from, to time.Time) (map[string]int, error) {
	query := `
		SELECT 
			CASE 
				WHEN referrer = '' OR referrer IS NULL THEN 'direct'
				WHEN referrer LIKE '%google%' THEN 'google'
				WHEN referrer LIKE '%facebook%' THEN 'facebook'
				WHEN referrer LIKE '%instagram%' THEN 'instagram'
				ELSE 'other'
			END as source,
			COUNT(*) as count
		FROM storefront_events
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
			AND event_type = 'page_view'
		GROUP BY source
		ORDER BY count DESC
	`
	
	rows, err := r.db.pool.Query(ctx, query, storefrontID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get traffic sources: %w", err)
	}
	defer rows.Close()
	
	sources := make(map[string]int)
	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err != nil {
			return nil, fmt.Errorf("failed to scan traffic source: %w", err)
		}
		sources[source] = count
	}
	
	return sources, nil
}