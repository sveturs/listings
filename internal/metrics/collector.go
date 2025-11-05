package metrics

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// DBStatsCollector periodically collects database connection pool statistics
type DBStatsCollector struct {
	db       *sqlx.DB
	metrics  *Metrics
	logger   zerolog.Logger
	interval time.Duration
	stopCh   chan struct{}
}

// NewDBStatsCollector creates a new database stats collector
func NewDBStatsCollector(db *sqlx.DB, metrics *Metrics, logger zerolog.Logger, interval time.Duration) *DBStatsCollector {
	return &DBStatsCollector{
		db:       db,
		metrics:  metrics,
		logger:   logger.With().Str("component", "db_stats_collector").Logger(),
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins collecting database stats at regular intervals
func (c *DBStatsCollector) Start(ctx context.Context) {
	c.logger.Info().Dur("interval", c.interval).Msg("starting DB stats collector")

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// Collect initial stats
	c.collect()

	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-c.stopCh:
			c.logger.Info().Msg("DB stats collector stopped")
			return
		case <-ctx.Done():
			c.logger.Info().Msg("DB stats collector context cancelled")
			return
		}
	}
}

// Stop stops the collector
func (c *DBStatsCollector) Stop() {
	close(c.stopCh)
}

// collect gathers and records current database stats
func (c *DBStatsCollector) collect() {
	stats := c.db.Stats()

	// Update connection pool metrics
	c.metrics.UpdateDBConnectionStats(stats.OpenConnections, stats.Idle)

	// Log detailed stats at debug level
	c.logger.Debug().
		Int("open", stats.OpenConnections).
		Int("idle", stats.Idle).
		Int("in_use", stats.InUse).
		Int("wait_count", int(stats.WaitCount)).
		Dur("wait_duration", stats.WaitDuration).
		Int("max_idle_closed", int(stats.MaxIdleClosed)).
		Int("max_idle_time_closed", int(stats.MaxIdleTimeClosed)).
		Int("max_lifetime_closed", int(stats.MaxLifetimeClosed)).
		Msg("database connection pool stats")

	// Additional metrics that could be tracked (if needed later):
	// - WaitCount: Total number of connections waited for
	// - WaitDuration: Total time blocked waiting for a new connection
	// - MaxIdleClosed: Total number of connections closed due to SetMaxIdleConns
	// - MaxIdleTimeClosed: Total number of connections closed due to SetConnMaxIdleTime
	// - MaxLifetimeClosed: Total number of connections closed due to SetConnMaxLifetime
}
