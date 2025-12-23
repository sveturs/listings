package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const searchAnalyticsIndex = "search_analytics"

// SearchEvent represents a search event for analytics
type SearchEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Query       string                 `json:"query"`
	QueryKeyword string                `json:"query_keyword"` // For aggregations
	UserID      *int64                 `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id"`
	ResultCount int64                  `json:"result_count"`
	TookMs      int64                  `json:"took_ms"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	Page        int                    `json:"page"`
	HasResults  bool                   `json:"has_results"`
	SearchType  string                 `json:"search_type"` // search, autocomplete, suggest
	Platform    string                 `json:"platform"`    // web, ios, android
	Language    string                 `json:"language"`    // sr, en, ru
	EventType   string                 `json:"event_type"`  // Always "search" for SearchEvent
}

// ClickEvent represents a click on search result
type ClickEvent struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	SearchEventID string    `json:"search_event_id"`
	ListingID     int64     `json:"listing_id"`
	Position      int       `json:"position"`
	UserID        *int64    `json:"user_id,omitempty"`
	SessionID     string    `json:"session_id"`
	EventType     string    `json:"event_type"` // Always "click"
}

// ConversionEvent represents a conversion (cart, purchase, favorite)
type ConversionEvent struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	SearchEventID  string    `json:"search_event_id"`
	ListingID      int64     `json:"listing_id"`
	ConversionType string    `json:"conversion_type"` // cart, purchase, favorite
	UserID         *int64    `json:"user_id,omitempty"`
	EventType      string    `json:"event_type"` // Always "conversion"
}

// SearchAnalyticsReport represents analytics report for a time period
type SearchAnalyticsReport struct {
	TotalSearches        int64                `json:"total_searches"`
	UniqueQueries        int64                `json:"unique_queries"`
	AvgResultCount       float64              `json:"avg_result_count"`
	ZeroResultRate       float64              `json:"zero_result_rate"` // % of searches with 0 results
	AvgLatency           float64              `json:"avg_latency"`
	TopQueries           []QueryStats         `json:"top_queries"`
	TopZeroResultQueries []string             `json:"top_zero_result_queries"`
	SearchesByPlatform   map[string]int64     `json:"searches_by_platform"`
	SearchesByLanguage   map[string]int64     `json:"searches_by_language"`
	CTR                  float64              `json:"ctr"`            // Click-Through Rate
	ConversionRate       float64              `json:"conversion_rate"`
}

// QueryStats represents statistics for a single query
type QueryStats struct {
	Query       string  `json:"query"`
	Count       int64   `json:"count"`
	AvgPosition float64 `json:"avg_position"`
	CTR         float64 `json:"ctr"`
}

// AnalyticsClient handles search analytics operations
type AnalyticsClient struct {
	searchClient *SearchClient
	logger       zerolog.Logger
}

// NewAnalyticsClient creates a new analytics client
func NewAnalyticsClient(searchClient *SearchClient, logger zerolog.Logger) *AnalyticsClient {
	return &AnalyticsClient{
		searchClient: searchClient,
		logger:       logger.With().Str("component", "analytics_client").Logger(),
	}
}

// CreateAnalyticsIndex creates the analytics index with proper mappings
func (ac *AnalyticsClient) CreateAnalyticsIndex(ctx context.Context) error {
	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":               map[string]string{"type": "keyword"},
				"timestamp":        map[string]string{"type": "date"},
				"query":            map[string]string{"type": "text", "analyzer": "standard"},
				"query_keyword":    map[string]string{"type": "keyword"},
				"user_id":          map[string]string{"type": "long"},
				"session_id":       map[string]string{"type": "keyword"},
				"result_count":     map[string]string{"type": "integer"},
				"took_ms":          map[string]string{"type": "integer"},
				"has_results":      map[string]string{"type": "boolean"},
				"search_type":      map[string]string{"type": "keyword"},
				"platform":         map[string]string{"type": "keyword"},
				"language":         map[string]string{"type": "keyword"},
				"filters":          map[string]string{"type": "object"},
				"page":             map[string]string{"type": "integer"},
				"event_type":       map[string]string{"type": "keyword"}, // search, click, conversion
				"search_event_id":  map[string]string{"type": "keyword"},
				"listing_id":       map[string]string{"type": "long"},
				"position":         map[string]string{"type": "integer"},
				"conversion_type":  map[string]string{"type": "keyword"},
			},
		},
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"index": map[string]interface{}{
				"lifecycle": map[string]string{
					"name": "search_analytics_policy", // 30 days retention
				},
			},
		},
	}

	return ac.createIndex(ctx, searchAnalyticsIndex, mappings)
}

// TrackSearch records a search event (async, non-blocking)
func (ac *AnalyticsClient) TrackSearch(ctx context.Context, event *SearchEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	event.EventType = "search"
	event.QueryKeyword = event.Query // For aggregations

	// Index asynchronously to not block search
	go func() {
		asyncCtx := context.Background()
		if err := ac.indexDocument(asyncCtx, searchAnalyticsIndex, event.ID, event); err != nil {
			ac.logger.Warn().
				Err(err).
				Str("event_id", event.ID).
				Str("query", event.Query).
				Msg("failed to track search event")
		}
	}()

	return nil
}

// TrackClick records a click event (async)
func (ac *AnalyticsClient) TrackClick(ctx context.Context, event *ClickEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	event.EventType = "click"

	go func() {
		asyncCtx := context.Background()
		if err := ac.indexDocument(asyncCtx, searchAnalyticsIndex, event.ID, event); err != nil {
			ac.logger.Warn().
				Err(err).
				Str("event_id", event.ID).
				Int64("listing_id", event.ListingID).
				Msg("failed to track click event")
		}
	}()

	return nil
}

// TrackConversion records a conversion event (async)
func (ac *AnalyticsClient) TrackConversion(ctx context.Context, event *ConversionEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	event.EventType = "conversion"

	go func() {
		asyncCtx := context.Background()
		if err := ac.indexDocument(asyncCtx, searchAnalyticsIndex, event.ID, event); err != nil {
			ac.logger.Warn().
				Err(err).
				Str("event_id", event.ID).
				Int64("listing_id", event.ListingID).
				Str("conversion_type", event.ConversionType).
				Msg("failed to track conversion event")
		}
	}()

	return nil
}

// GetSearchAnalytics retrieves analytics report for a time period
func (ac *AnalyticsClient) GetSearchAnalytics(ctx context.Context, from, to time.Time) (*SearchAnalyticsReport, error) {
	// Query all event types within the time range
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": from.Format(time.RFC3339),
					"lte": to.Format(time.RFC3339),
				},
			},
		},
		"aggs": map[string]interface{}{
			// Count by event type (search, click, conversion)
			"by_event_type": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "event_type",
					"size":  10,
				},
			},
			// Search-specific aggregations (filtered)
			"searches": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]string{"event_type": "search"},
				},
				"aggs": map[string]interface{}{
					"unique_queries": map[string]interface{}{
						"cardinality": map[string]string{"field": "query_keyword"},
					},
					"avg_result_count": map[string]interface{}{
						"avg": map[string]string{"field": "result_count"},
					},
					"avg_latency": map[string]interface{}{
						"avg": map[string]string{"field": "took_ms"},
					},
					"zero_results": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]bool{"has_results": false},
						},
					},
					"by_platform": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "platform",
							"size":  10,
						},
					},
					"by_language": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "language",
							"size":  10,
						},
					},
					"top_queries": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "query_keyword",
							"size":  100,
						},
					},
				},
			},
			// Click count
			"clicks": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]string{"event_type": "click"},
				},
			},
			// Conversion count
			"conversions": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]string{"event_type": "conversion"},
				},
			},
		},
	}

	resp, err := ac.searchClient.SearchIndex(ctx, searchAnalyticsIndex, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get search analytics: %w", err)
	}

	return ac.parseAnalyticsReport(resp)
}

// GetTopZeroResultQueries retrieves top queries with zero results
func (ac *AnalyticsClient) GetTopZeroResultQueries(ctx context.Context, from, to time.Time, limit int) ([]string, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"gte": from.Format(time.RFC3339),
								"lte": to.Format(time.RFC3339),
							},
						},
					},
					{
						"term": map[string]bool{"has_results": false},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"zero_queries": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "query_keyword",
					"size":  limit,
				},
			},
		},
	}

	resp, err := ac.searchClient.SearchIndex(ctx, searchAnalyticsIndex, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get zero result queries: %w", err)
	}

	return ac.parseTopQueries(resp, "zero_queries")
}

// parseAnalyticsReport parses SearchResponse into SearchAnalyticsReport
func (ac *AnalyticsClient) parseAnalyticsReport(resp *SearchResponse) (*SearchAnalyticsReport, error) {
	report := &SearchAnalyticsReport{
		SearchesByPlatform: make(map[string]int64),
		SearchesByLanguage: make(map[string]int64),
	}

	if resp.Aggregations == nil {
		return report, nil
	}

	var totalClicks, totalConversions int64

	// Parse clicks count
	if clicks, ok := resp.Aggregations["clicks"].(map[string]interface{}); ok {
		if docCount, ok := clicks["doc_count"].(float64); ok {
			totalClicks = int64(docCount)
		}
	}

	// Parse conversions count
	if conversions, ok := resp.Aggregations["conversions"].(map[string]interface{}); ok {
		if docCount, ok := conversions["doc_count"].(float64); ok {
			totalConversions = int64(docCount)
		}
	}

	// Parse searches aggregation (nested)
	searches, ok := resp.Aggregations["searches"].(map[string]interface{})
	if !ok {
		return report, nil
	}

	// Get total searches from doc_count of searches filter
	if docCount, ok := searches["doc_count"].(float64); ok {
		report.TotalSearches = int64(docCount)
	}

	// Parse unique queries
	if uniqueQueries, ok := searches["unique_queries"].(map[string]interface{}); ok {
		if value, ok := uniqueQueries["value"].(float64); ok {
			report.UniqueQueries = int64(value)
		}
	}

	// Parse avg result count
	if avgResultCount, ok := searches["avg_result_count"].(map[string]interface{}); ok {
		if value, ok := avgResultCount["value"].(float64); ok {
			report.AvgResultCount = value
		}
	}

	// Parse avg latency
	if avgLatency, ok := searches["avg_latency"].(map[string]interface{}); ok {
		if value, ok := avgLatency["value"].(float64); ok {
			report.AvgLatency = value
		}
	}

	// Parse zero results
	if zeroResults, ok := searches["zero_results"].(map[string]interface{}); ok {
		if docCount, ok := zeroResults["doc_count"].(float64); ok && report.TotalSearches > 0 {
			report.ZeroResultRate = (docCount / float64(report.TotalSearches)) * 100
		}
	}

	// Parse by platform
	if byPlatform, ok := searches["by_platform"].(map[string]interface{}); ok {
		if buckets, ok := byPlatform["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					if key, ok := b["key"].(string); ok {
						if docCount, ok := b["doc_count"].(float64); ok {
							report.SearchesByPlatform[key] = int64(docCount)
						}
					}
				}
			}
		}
	}

	// Parse by language
	if byLanguage, ok := searches["by_language"].(map[string]interface{}); ok {
		if buckets, ok := byLanguage["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					if key, ok := b["key"].(string); ok {
						if docCount, ok := b["doc_count"].(float64); ok {
							report.SearchesByLanguage[key] = int64(docCount)
						}
					}
				}
			}
		}
	}

	// Parse top queries from nested aggregation
	report.TopQueries = ac.parseNestedQueryStats(searches, "top_queries")

	// Calculate CTR and Conversion Rate
	if report.TotalSearches > 0 {
		report.CTR = float64(totalClicks) / float64(report.TotalSearches)
		report.ConversionRate = float64(totalConversions) / float64(report.TotalSearches)
	}

	return report, nil
}

// parseQueryStats parses query statistics from aggregations
func (ac *AnalyticsClient) parseQueryStats(resp *SearchResponse, aggName string) []QueryStats {
	var stats []QueryStats

	if resp.Aggregations == nil {
		return stats
	}

	if topQueries, ok := resp.Aggregations[aggName].(map[string]interface{}); ok {
		if buckets, ok := topQueries["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					stat := QueryStats{}
					if key, ok := b["key"].(string); ok {
						stat.Query = key
					}
					if docCount, ok := b["doc_count"].(float64); ok {
						stat.Count = int64(docCount)
					}
					stats = append(stats, stat)
				}
			}
		}
	}

	return stats
}

// parseNestedQueryStats parses query statistics from a nested aggregation bucket
func (ac *AnalyticsClient) parseNestedQueryStats(bucket map[string]interface{}, aggName string) []QueryStats {
	var stats []QueryStats

	if bucket == nil {
		return stats
	}

	if topQueries, ok := bucket[aggName].(map[string]interface{}); ok {
		if buckets, ok := topQueries["buckets"].([]interface{}); ok {
			for _, b := range buckets {
				if bm, ok := b.(map[string]interface{}); ok {
					stat := QueryStats{}
					if key, ok := bm["key"].(string); ok {
						stat.Query = key
					}
					if docCount, ok := bm["doc_count"].(float64); ok {
						stat.Count = int64(docCount)
					}
					stats = append(stats, stat)
				}
			}
		}
	}

	return stats
}

// parseTopQueries parses top queries from aggregations
func (ac *AnalyticsClient) parseTopQueries(resp *SearchResponse, aggName string) ([]string, error) {
	var queries []string

	if resp.Aggregations == nil {
		return queries, nil
	}

	if topQueries, ok := resp.Aggregations[aggName].(map[string]interface{}); ok {
		if buckets, ok := topQueries["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					if key, ok := b["key"].(string); ok {
						queries = append(queries, key)
					}
				}
			}
		}
	}

	return queries, nil
}

// createIndex creates an index with given mappings (helper method)
func (ac *AnalyticsClient) createIndex(ctx context.Context, indexName string, mappings map[string]interface{}) error {
	// Convert mappings to JSON
	body, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("failed to marshal index mappings: %w", err)
	}

	// Check if index exists
	exists, err := ac.searchClient.client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer exists.Body.Close()

	if exists.StatusCode == 200 {
		ac.logger.Info().Str("index", indexName).Msg("index already exists")
		return nil
	}

	// Create index
	res, err := ac.searchClient.client.Indices.Create(
		indexName,
		ac.searchClient.client.Indices.Create.WithBody(bytes.NewReader(body)),
		ac.searchClient.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index [%s]", res.Status())
	}

	ac.logger.Info().Str("index", indexName).Msg("index created successfully")
	return nil
}

// indexDocument indexes a document (helper method)
func (ac *AnalyticsClient) indexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	res, err := ac.searchClient.client.Index(
		indexName,
		bytes.NewReader(body),
		ac.searchClient.client.Index.WithDocumentID(docID),
		ac.searchClient.client.Index.WithContext(ctx),
		ac.searchClient.client.Index.WithRefresh("false"), // Async
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document [%s]", res.Status())
	}

	return nil
}
