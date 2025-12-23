package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
)

// Prometheus metrics for OpenSearch monitoring
var (
	opensearchRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "request_duration_seconds",
			Help:      "Duration of OpenSearch requests",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)

	opensearchRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "requests_total",
			Help:      "Total number of OpenSearch requests",
		},
		[]string{"operation", "status"},
	)

	opensearchIndexedDocs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "indexed_documents",
			Help:      "Number of documents in index",
		},
		[]string{"index"},
	)

	opensearchClusterStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "cluster_status",
			Help:      "Cluster health status (0=red, 1=yellow, 2=green)",
		},
		[]string{"cluster"},
	)

	opensearchSearchLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "search_latency_seconds",
			Help:      "Search query latency",
			Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5},
		},
		[]string{"query_type"}, // search, autocomplete, suggest
	)

	opensearchReindexProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "reindex_progress",
			Help:      "Reindexation progress (0-100%)",
		},
		[]string{"index"},
	)

	opensearchIndexSizeMB = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "index_size_mb",
			Help:      "Index size in megabytes",
		},
		[]string{"index"},
	)

	opensearchShardCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "vondi_listings",
			Subsystem: "opensearch",
			Name:      "shard_count",
			Help:      "Number of shards (primary + replica)",
		},
		[]string{"index", "shard_type"}, // primary, replica
	)
)

// HealthStatus represents OpenSearch health status
type HealthStatus struct {
	Status        string        // healthy, degraded, unhealthy
	ClusterHealth string        // green, yellow, red
	IndexExists   bool          // whether index exists
	DocsCount     int64         // document count
	LastCheck     time.Time     // last check timestamp
	Latency       time.Duration // test query latency
	IndexSizeMB   float64       // index size in MB
	Shards        ShardInfo     // shard information
}

// ShardInfo contains shard statistics
type ShardInfo struct {
	Total      int // total shards
	Successful int // successful shards
	Failed     int // failed shards
	Primary    int // primary shards
	Replica    int // replica shards
}

// IndexStats contains index statistics
type IndexStats struct {
	Index        string        // index name
	DocsCount    int64         // document count
	StoreSize    string        // storage size
	StoreSizeMB  float64       // storage size in MB
	SegmentCount int           // segment count
	Shards       ShardInfo     // shard information
	RefreshTime  time.Duration // refresh time
	FlushTime    time.Duration // flush time
}

// ClusterHealthResponse represents cluster health API response
type clusterHealthResponse struct {
	ClusterName         string `json:"cluster_name"`
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}

// indexStatsResponse represents index stats API response
type indexStatsResponse struct {
	Indices map[string]indexStatsDetail `json:"indices"`
}

type indexStatsDetail struct {
	Primaries struct {
		Docs struct {
			Count int64 `json:"count"`
		} `json:"docs"`
		Store struct {
			SizeInBytes int64 `json:"size_in_bytes"`
		} `json:"store"`
		Segments struct {
			Count int `json:"count"`
		} `json:"segments"`
		Refresh struct {
			TotalTimeInMillis int64 `json:"total_time_in_millis"`
		} `json:"refresh"`
		Flush struct {
			TotalTimeInMillis int64 `json:"total_time_in_millis"`
		} `json:"flush"`
	} `json:"primaries"`
	Shards map[string][]struct {
		Routing struct {
			State   string `json:"state"`
			Primary bool   `json:"primary"`
		} `json:"routing"`
	} `json:"shards"`
}

// HealthCheckDetailed performs comprehensive health check on OpenSearch
func (c *Client) HealthCheckDetailed(ctx context.Context) (*HealthStatus, error) {
	start := time.Now()
	status := &HealthStatus{
		LastCheck: time.Now(),
	}

	// 1. Cluster health
	healthResp, err := c.getClusterHealth(ctx)
	if err != nil {
		status.Status = "unhealthy"
		status.ClusterHealth = "unknown"
		c.logger.Error().Err(err).Msg("cluster health check failed")
		return status, err
	}
	status.ClusterHealth = healthResp.Status

	// 2. Index exists check
	status.IndexExists = c.IndexExists(ctx, c.index)
	if !status.IndexExists {
		status.Status = "unhealthy"
		c.logger.Warn().Str("index", c.index).Msg("index does not exist")
		return status, fmt.Errorf("index %s does not exist", c.index)
	}

	// 3. Document count
	count, err := c.CountDocuments(ctx, c.index)
	if err != nil {
		c.logger.Warn().Err(err).Msg("failed to count documents")
	} else {
		status.DocsCount = int64(count)
	}

	// 4. Test query latency
	queryStart := time.Now()
	_, err = c.GetListingByID(ctx, 1) // Try to get a sample document
	if err != nil {
		c.logger.Debug().Err(err).Msg("test query failed (may be normal if no docs)")
	}
	status.Latency = time.Since(queryStart)

	// 5. Get index stats
	indexStats, err := c.GetIndexStats(ctx)
	if err != nil {
		c.logger.Warn().Err(err).Msg("failed to get index stats")
	} else {
		status.IndexSizeMB = indexStats.StoreSizeMB
		status.Shards = indexStats.Shards
	}

	// Determine overall status
	if status.ClusterHealth == "red" || !status.IndexExists {
		status.Status = "unhealthy"
	} else if status.ClusterHealth == "yellow" || status.Latency > 500*time.Millisecond {
		status.Status = "degraded"
	} else {
		status.Status = "healthy"
	}

	// Record total health check duration
	healthCheckDuration := time.Since(start)

	// Update Prometheus metrics
	opensearchClusterStatus.WithLabelValues(c.index).Set(healthStatusToFloat(status.ClusterHealth))
	opensearchIndexedDocs.WithLabelValues(c.index).Set(float64(status.DocsCount))
	opensearchIndexSizeMB.WithLabelValues(c.index).Set(status.IndexSizeMB)
	opensearchShardCount.WithLabelValues(c.index, "primary").Set(float64(status.Shards.Primary))
	opensearchShardCount.WithLabelValues(c.index, "replica").Set(float64(status.Shards.Replica))

	c.logger.Info().
		Str("status", status.Status).
		Str("cluster_health", status.ClusterHealth).
		Int64("docs", status.DocsCount).
		Float64("index_size_mb", status.IndexSizeMB).
		Dur("latency", status.Latency).
		Dur("health_check_duration", healthCheckDuration).
		Msg("OpenSearch health check completed")

	return status, nil
}

// RecordRequest records metrics for an OpenSearch request
func (c *Client) RecordRequest(operation string, start time.Time, err error) {
	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "error"
	}

	opensearchRequestDuration.WithLabelValues(operation, status).Observe(duration)
	opensearchRequestTotal.WithLabelValues(operation, status).Inc()
}

// RecordSearch records metrics for a search query
func (c *Client) RecordSearch(queryType string, start time.Time) {
	duration := time.Since(start).Seconds()
	opensearchSearchLatency.WithLabelValues(queryType).Observe(duration)
}

// UpdateReindexProgress updates reindex progress metric
func (c *Client) UpdateReindexProgress(progress float64) {
	opensearchReindexProgress.WithLabelValues(c.index).Set(progress)
}

// GetIndexStats retrieves detailed index statistics
func (c *Client) GetIndexStats(ctx context.Context) (*IndexStats, error) {
	// GET /{index}/_stats
	res, err := c.client.Indices.Stats(
		c.client.Indices.Stats.WithIndex(c.index),
		c.client.Indices.Stats.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("index stats error: %s - %s", res.Status(), string(bodyBytes))
	}

	var statsResp indexStatsResponse
	if err := json.NewDecoder(res.Body).Decode(&statsResp); err != nil {
		return nil, fmt.Errorf("failed to parse index stats: %w", err)
	}

	indexDetail, exists := statsResp.Indices[c.index]
	if !exists {
		return nil, fmt.Errorf("index %s not found in stats response", c.index)
	}

	// Calculate shard counts
	shards := ShardInfo{}
	for _, shardList := range indexDetail.Shards {
		for _, shard := range shardList {
			shards.Total++
			if shard.Routing.State == "STARTED" {
				shards.Successful++
			} else {
				shards.Failed++
			}
			if shard.Routing.Primary {
				shards.Primary++
			} else {
				shards.Replica++
			}
		}
	}

	// Convert bytes to MB
	sizeInBytes := indexDetail.Primaries.Store.SizeInBytes
	sizeMB := float64(sizeInBytes) / 1024 / 1024

	stats := &IndexStats{
		Index:        c.index,
		DocsCount:    indexDetail.Primaries.Docs.Count,
		StoreSizeMB:  sizeMB,
		StoreSize:    formatBytes(sizeInBytes),
		SegmentCount: indexDetail.Primaries.Segments.Count,
		Shards:       shards,
		RefreshTime:  time.Duration(indexDetail.Primaries.Refresh.TotalTimeInMillis) * time.Millisecond,
		FlushTime:    time.Duration(indexDetail.Primaries.Flush.TotalTimeInMillis) * time.Millisecond,
	}

	c.logger.Debug().
		Str("index", stats.Index).
		Int64("docs", stats.DocsCount).
		Str("size", stats.StoreSize).
		Int("segments", stats.SegmentCount).
		Int("total_shards", stats.Shards.Total).
		Msg("index stats retrieved")

	return stats, nil
}

// IndexExists checks if an index exists
func (c *Client) IndexExists(ctx context.Context, indexName string) bool {
	res, err := c.client.Indices.Exists(
		[]string{indexName},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		c.logger.Error().Err(err).Str("index", indexName).Msg("failed to check index existence")
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}

// getClusterHealth retrieves cluster health information
func (c *Client) getClusterHealth(ctx context.Context) (*clusterHealthResponse, error) {
	res, err := c.client.Cluster.Health(
		c.client.Cluster.Health.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("cluster health request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("cluster health error: %s - %s", res.Status(), string(bodyBytes))
	}

	var healthResp clusterHealthResponse
	if err := json.NewDecoder(res.Body).Decode(&healthResp); err != nil {
		return nil, fmt.Errorf("failed to parse cluster health: %w", err)
	}

	return &healthResp, nil
}

// healthStatusToFloat converts health status string to metric value
func healthStatusToFloat(status string) float64 {
	switch status {
	case "green":
		return 2
	case "yellow":
		return 1
	default: // red, unknown
		return 0
	}
}

// formatBytes formats bytes to human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// StatsCollector periodically collects OpenSearch statistics
type StatsCollector struct {
	client   *Client
	interval time.Duration
	logger   zerolog.Logger
	stopCh   chan struct{}
}

// NewStatsCollector creates a new statistics collector
func NewStatsCollector(client *Client, interval time.Duration, logger zerolog.Logger) *StatsCollector {
	return &StatsCollector{
		client:   client,
		interval: interval,
		logger:   logger.With().Str("component", "opensearch_stats_collector").Logger(),
		stopCh:   make(chan struct{}),
	}
}

// Start begins periodic statistics collection
func (sc *StatsCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(sc.interval)
	defer ticker.Stop()

	sc.logger.Info().Dur("interval", sc.interval).Msg("OpenSearch stats collector started")

	for {
		select {
		case <-ticker.C:
			sc.collectStats(ctx)
		case <-sc.stopCh:
			sc.logger.Info().Msg("OpenSearch stats collector stopped")
			return
		case <-ctx.Done():
			sc.logger.Info().Msg("OpenSearch stats collector stopped (context cancelled)")
			return
		}
	}
}

// Stop stops the statistics collector
func (sc *StatsCollector) Stop() {
	close(sc.stopCh)
}

// collectStats collects and records statistics
func (sc *StatsCollector) collectStats(ctx context.Context) {
	health, err := sc.client.HealthCheckDetailed(ctx)
	if err != nil {
		sc.logger.Error().Err(err).Msg("failed to collect OpenSearch stats")
		return
	}

	sc.logger.Debug().
		Str("status", health.Status).
		Str("cluster", health.ClusterHealth).
		Int64("docs", health.DocsCount).
		Float64("size_mb", health.IndexSizeMB).
		Dur("latency", health.Latency).
		Msg("OpenSearch stats collected")
}
