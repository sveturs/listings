// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// StorefrontEventRepository handles storefront event storage
type StorefrontEventRepository interface {
	// RecordEvent stores a new storefront analytics event
	RecordEvent(ctx context.Context, event *domain.StorefrontEvent) error
}

// storefrontEventRepository implements StorefrontEventRepository
type storefrontEventRepository struct {
	db     *pgxpool.Pool
	logger zerolog.Logger
}

// NewStorefrontEventRepository creates a new storefront event repository
func NewStorefrontEventRepository(pool *pgxpool.Pool, logger zerolog.Logger) StorefrontEventRepository {
	return &storefrontEventRepository{
		db:     pool,
		logger: logger.With().Str("component", "storefront_event_repository").Logger(),
	}
}

// RecordEvent stores a new storefront analytics event
// Target: < 10ms (simple INSERT with indexes)
func (r *storefrontEventRepository) RecordEvent(ctx context.Context, event *domain.StorefrontEvent) error {
	r.logger.Debug().
		Int64("storefront_id", event.StorefrontID).
		Str("event_type", event.EventType).
		Str("session_id", event.SessionID).
		Msg("recording storefront event")

	query := `
		INSERT INTO storefront_events (
			storefront_id,
			event_type,
			event_data,
			user_id,
			session_id,
			ip_address,
			user_agent,
			referrer,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		event.StorefrontID,
		event.EventType,
		event.EventData,
		event.UserID,
		event.SessionID,
		event.IPAddress,
		event.UserAgent,
		event.Referrer,
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Int64("storefront_id", event.StorefrontID).
			Str("event_type", event.EventType).
			Msg("failed to record storefront event")
		return fmt.Errorf("failed to record storefront event: %w", err)
	}

	r.logger.Debug().
		Int64("storefront_id", event.StorefrontID).
		Str("event_type", event.EventType).
		Msg("storefront event recorded successfully")

	return nil
}
