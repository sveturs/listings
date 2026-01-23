package reindexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchutil"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/indexer"
	osrepo "github.com/vondi-global/listings/internal/repository/opensearch"
)

// ReindexManager manages blue-green index switching for zero-downtime reindexing
type ReindexManager struct {
	client   *opensearch.Client
	osClient *osrepo.Client
	indexer  *indexer.ListingIndexer
	logger   zerolog.Logger
	progress *ReindexProgress
}

// VerificationResult contains reindex verification metrics
type VerificationResult struct {
	Valid           bool
	TotalDocs       int64
	ExpectedDocs    int64
	MismatchedCount int64
	SampleErrors    []string
	FieldCoverage   map[string]float64
}

// ReindexProgress tracks reindexing progress
type ReindexProgress struct {
	Total      int64
	Indexed    int64
	Failed     int64
	StartTime  time.Time
	LastUpdate time.Time
}

// NewReindexManager creates a new ReindexManager instance
func NewReindexManager(osClient *osrepo.Client, indexer *indexer.ListingIndexer, logger zerolog.Logger) *ReindexManager {
	return &ReindexManager{
		client:   osClient.GetClient(),
		osClient: osClient,
		indexer:  indexer,
		logger:   logger.With().Str("component", "reindex_manager").Logger(),
		progress: &ReindexProgress{},
	}
}

// GetProgress returns current reindexing progress
func (m *ReindexManager) GetProgress() *ReindexProgress {
	return m.progress
}

// StartBlueGreenReindex performs blue-green reindexing with zero downtime
// Algorithm:
// 1. Determine current index version (v1 or v2)
// 2. Create new index with opposite version
// 3. Perform full reindexing to new index
// 4. Verify results
// 5. Atomic alias switch from old to new
// 6. Keep old index for 24h rollback window
func (m *ReindexManager) StartBlueGreenReindex(ctx context.Context, batchSize int) error {
	startTime := time.Now()

	// Step 1: Get current index version
	currentVersion, err := m.GetCurrentIndexVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current index version: %w", err)
	}

	m.logger.Info().Str("current_version", currentVersion).Msg("starting blue-green reindex")

	// Step 2: Determine new version
	newVersion := "v2"
	if currentVersion == "v2" {
		newVersion = "v1"
	}
	newIndexName := fmt.Sprintf("marketplace_listings_%s", newVersion)
	oldIndexName := fmt.Sprintf("marketplace_listings_%s", currentVersion)

	m.logger.Info().
		Str("old_index", oldIndexName).
		Str("new_index", newIndexName).
		Msg("determined blue-green versions")

	// Step 3: Delete new index if exists (from previous failed attempt)
	if err := m.DeleteIndexIfExists(ctx, newIndexName); err != nil {
		return fmt.Errorf("failed to delete old new index: %w", err)
	}

	// Step 4: Create new index with latest mappings
	mapping := osrepo.GetListingsIndexMapping()
	if err := m.osClient.CreateIndex(ctx, newIndexName, mapping); err != nil {
		return fmt.Errorf("failed to create new index: %w", err)
	}

	m.logger.Info().Str("index", newIndexName).Msg("created new index with mappings")

	// Step 5: Perform full reindexing to new index
	// We'll reindex directly using OpenSearch bulk API to the new index
	m.logger.Info().Int("batch_size", batchSize).Msg("starting full reindexing")

	if err := m.reindexToNewIndex(ctx, newIndexName, batchSize); err != nil {
		return fmt.Errorf("reindexing failed: %w", err)
	}

	m.logger.Info().Str("index", newIndexName).Msg("reindexing completed")

	// Step 6: Force refresh to ensure all docs are searchable
	if err := m.RefreshIndex(ctx, newIndexName); err != nil {
		return fmt.Errorf("failed to refresh new index: %w", err)
	}

	// Step 7: Verify reindex results
	result, err := m.VerifyReindex(ctx, newIndexName)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	if !result.Valid {
		m.logger.Error().
			Interface("result", result).
			Msg("verification failed - NOT switching alias")
		return fmt.Errorf("verification failed: %+v", result)
	}

	m.logger.Info().
		Int64("total_docs", result.TotalDocs).
		Int64("expected_docs", result.ExpectedDocs).
		Msg("verification passed")

	// Step 8: Atomic alias switch
	if err := m.SwitchAlias(ctx, oldIndexName, newIndexName); err != nil {
		return fmt.Errorf("failed to switch alias: %w", err)
	}

	m.logger.Info().
		Str("from_index", oldIndexName).
		Str("to_index", newIndexName).
		Msg("alias switched successfully")

	// Step 9: Log completion
	elapsed := time.Since(startTime)
	m.logger.Info().
		Str("old_index", oldIndexName).
		Str("new_index", newIndexName).
		Str("duration", elapsed.String()).
		Msg("blue-green reindex completed - old index kept for 24h rollback")

	return nil
}

// GetCurrentIndexVersion determines which index version is currently active
// Returns "v1" or "v2" based on alias resolution
func (m *ReindexManager) GetCurrentIndexVersion(ctx context.Context) (string, error) {
	// Get alias information
	res, err := m.client.Indices.GetAlias(
		m.client.Indices.GetAlias.WithContext(ctx),
		m.client.Indices.GetAlias.WithName("marketplace_listings"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get alias: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// No alias exists - this is initial setup
		// Check which indexes exist
		existsV1, _ := m.IndexExists(ctx, "marketplace_listings_v1")
		existsV2, _ := m.IndexExists(ctx, "marketplace_listings_v2")

		if existsV1 {
			return "v1", nil
		}
		if existsV2 {
			return "v2", nil
		}

		// No indexes exist - default to v1
		return "v1", nil
	}

	if res.IsError() {
		return "", fmt.Errorf("get alias error: %s", res.Status())
	}

	// Parse response
	var aliasResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&aliasResp); err != nil {
		return "", fmt.Errorf("failed to parse alias response: %w", err)
	}

	// Find which index has the alias
	for indexName := range aliasResp {
		if indexName == "marketplace_listings_v1" {
			return "v1", nil
		}
		if indexName == "marketplace_listings_v2" {
			return "v2", nil
		}
	}

	// Default to v1 if nothing found
	return "v1", nil
}

// SwitchAlias atomically switches alias from old index to new index
func (m *ReindexManager) SwitchAlias(ctx context.Context, fromIndex, toIndex string) error {
	// Build atomic alias update request
	body := map[string]interface{}{
		"actions": []map[string]interface{}{
			{"remove": map[string]string{"index": fromIndex, "alias": "marketplace_listings"}},
			{"add": map[string]string{"index": toIndex, "alias": "marketplace_listings"}},
		},
	}

	// Execute atomic alias switch
	res, err := m.client.Indices.UpdateAliases(
		opensearchutil.NewJSONReader(body),
		m.client.Indices.UpdateAliases.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to update aliases: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("update aliases error: %s - %s", res.Status(), string(bodyBytes))
	}

	m.logger.Info().
		Str("from", fromIndex).
		Str("to", toIndex).
		Msg("alias switched atomically")

	return nil
}

// VerifyReindex verifies that reindexing completed successfully
func (m *ReindexManager) VerifyReindex(ctx context.Context, indexName string) (*VerificationResult, error) {
	result := &VerificationResult{
		Valid:         true,
		FieldCoverage: make(map[string]float64),
		SampleErrors:  make([]string, 0),
	}

	// 1. Compare document counts
	indexCount, err := m.osClient.CountDocuments(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to count index documents: %w", err)
	}
	result.TotalDocs = int64(indexCount)

	// Get expected count from database
	dbCount, err := m.getDBCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count database documents: %w", err)
	}
	result.ExpectedDocs = dbCount

	// Allow 5% tolerance for documents in process
	minExpected := float64(dbCount) * 0.95
	if float64(indexCount) < minExpected {
		result.Valid = false
		result.SampleErrors = append(result.SampleErrors,
			fmt.Sprintf("document count too low: expected %d, got %d (%.2f%%)",
				dbCount, indexCount, float64(indexCount)/float64(dbCount)*100))
	}

	// 2. Check field coverage for critical fields
	criticalFields := []string{"title", "price", "category_id", "status"}
	for _, field := range criticalFields {
		coverage, err := m.getFieldCoverage(ctx, indexName, field)
		if err != nil {
			result.SampleErrors = append(result.SampleErrors,
				fmt.Sprintf("failed to check coverage for %s: %v", field, err))
			continue
		}

		result.FieldCoverage[field] = coverage
		if coverage < 99.0 {
			result.Valid = false
			result.SampleErrors = append(result.SampleErrors,
				fmt.Sprintf("field %s coverage too low: %.2f%%", field, coverage))
		}
	}

	// 3. Test sample search queries
	testQueries := []string{"telefon", "patike", "auto"}
	for _, query := range testQueries {
		count, err := m.testSearch(ctx, indexName, query)
		if err != nil {
			result.SampleErrors = append(result.SampleErrors,
				fmt.Sprintf("test query '%s' failed: %v", query, err))
			continue
		}
		if count == 0 {
			// Warning but not failure - query might legitimately return 0
			m.logger.Warn().Str("query", query).Msg("test query returned 0 results")
		}
	}

	return result, nil
}

// RollbackToOldIndex switches alias back to old index version
func (m *ReindexManager) RollbackToOldIndex(ctx context.Context, oldIndex string) error {
	// Get current index
	currentVersion, err := m.GetCurrentIndexVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	currentIndex := fmt.Sprintf("marketplace_listings_%s", currentVersion)

	m.logger.Warn().
		Str("from", currentIndex).
		Str("to", oldIndex).
		Msg("rolling back to old index")

	// Switch alias back
	return m.SwitchAlias(ctx, currentIndex, oldIndex)
}

// DeleteIndexIfExists deletes an index if it exists
func (m *ReindexManager) DeleteIndexIfExists(ctx context.Context, indexName string) error {
	exists, err := m.IndexExists(ctx, indexName)
	if err != nil {
		return err
	}

	if !exists {
		m.logger.Debug().Str("index", indexName).Msg("index does not exist - skipping delete")
		return nil
	}

	return m.osClient.DeleteIndex(ctx, indexName)
}

// IndexExists checks if an index exists
func (m *ReindexManager) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := m.client.Indices.Exists(
		[]string{indexName},
		m.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// RefreshIndex forces index refresh to make all documents searchable
func (m *ReindexManager) RefreshIndex(ctx context.Context, indexName string) error {
	res, err := m.client.Indices.Refresh(
		m.client.Indices.Refresh.WithContext(ctx),
		m.client.Indices.Refresh.WithIndex(indexName),
	)
	if err != nil {
		return fmt.Errorf("failed to refresh index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("refresh index error: %s", res.Status())
	}

	m.logger.Debug().Str("index", indexName).Msg("index refreshed")
	return nil
}

// getDBCount returns the total number of active public listings in database
func (m *ReindexManager) getDBCount(ctx context.Context) (int64, error) {
	return m.indexer.CountActiveListings(ctx)
}

// getFieldCoverage returns percentage of documents that have a specific field
func (m *ReindexManager) getFieldCoverage(ctx context.Context, indexName, field string) (float64, error) {
	// Query for documents that have this field
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"exists": map[string]interface{}{
				"field": field,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := m.client.Count(
		m.client.Count.WithContext(ctx),
		m.client.Count.WithIndex(indexName),
		m.client.Count.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return 0, fmt.Errorf("count query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("count query error: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse count response: %w", err)
	}

	count, ok := result["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid count response format")
	}

	// Get total count
	totalCount, err := m.osClient.CountDocuments(ctx, indexName)
	if err != nil {
		return 0, err
	}

	if totalCount == 0 {
		return 0, nil
	}

	coverage := (count / float64(totalCount)) * 100
	return coverage, nil
}

// testSearch performs a test search query and returns result count
func (m *ReindexManager) testSearch(ctx context.Context, indexName, queryText string) (int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  queryText,
				"fields": []string{"title", "description"},
			},
		},
		"size": 0, // Only count, don't return docs
	}

	body, err := json.Marshal(query)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := m.client.Search(
		m.client.Search.WithContext(ctx),
		m.client.Search.WithIndex(indexName),
		m.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return 0, fmt.Errorf("search query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("search query error: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse search response: %w", err)
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid hits structure")
	}

	total, ok := hits["total"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid total structure")
	}

	value, ok := total["value"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid total value")
	}

	return int(value), nil
}

// reindexToNewIndex performs reindexing using OpenSearch Reindex API
// This is faster than manual bulk indexing as it copies data server-side
func (m *ReindexManager) reindexToNewIndex(ctx context.Context, newIndex string, batchSize int) error {
	// Use OpenSearch Reindex API to copy from current index to new index
	// Get current index name
	currentVersion, err := m.GetCurrentIndexVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	sourceIndex := fmt.Sprintf("marketplace_listings_%s", currentVersion)

	// Note: If this is initial setup and sourceIndex doesn't exist,
	// fall back to using the indexer to populate from DB
	exists, err := m.IndexExists(ctx, sourceIndex)
	if err != nil {
		return err
	}

	if !exists {
		// No source index - this is first-time setup
		// Use indexer to populate from database
		m.logger.Info().Msg("source index does not exist - performing initial indexing from database")
		return m.indexer.ReindexAllWithAttributes(ctx, batchSize)
	}

	// Source index exists - use OpenSearch Reindex API for fast server-side copy
	m.logger.Info().
		Str("source", sourceIndex).
		Str("dest", newIndex).
		Msg("performing server-side reindex")

	body := map[string]interface{}{
		"source": map[string]interface{}{
			"index": sourceIndex,
		},
		"dest": map[string]interface{}{
			"index": newIndex,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal reindex request: %w", err)
	}

	res, err := m.client.Reindex(
		bytes.NewReader(bodyBytes),
		m.client.Reindex.WithContext(ctx),
		m.client.Reindex.WithWaitForCompletion(true),
	)
	if err != nil {
		return fmt.Errorf("reindex request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("reindex error: %s - %s", res.Status(), string(bodyBytes))
	}

	var reindexResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&reindexResp); err != nil {
		return fmt.Errorf("failed to parse reindex response: %w", err)
	}

	m.logger.Info().
		Interface("response", reindexResp).
		Msg("reindex completed")

	return nil
}
