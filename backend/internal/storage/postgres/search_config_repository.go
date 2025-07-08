package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SearchConfigRepository struct {
	db *sqlx.DB
}

func NewSearchConfigRepository(db *sqlx.DB) *SearchConfigRepository {
	return &SearchConfigRepository{db: db}
}

// Weights methods
func (r *SearchConfigRepository) GetWeights(ctx context.Context) ([]domain.SearchWeight, error) {
	var weights []domain.SearchWeight
	query := `SELECT id, field_name, weight, description, created_at, updated_at 
	          FROM search_weights ORDER BY field_name`

	err := r.db.SelectContext(ctx, &weights, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get search weights: %w", err)
	}

	return weights, nil
}

func (r *SearchConfigRepository) GetWeightByField(ctx context.Context, fieldName string) (*domain.SearchWeight, error) {
	var weight domain.SearchWeight
	query := `SELECT id, field_name, weight, description, created_at, updated_at 
	          FROM search_weights WHERE field_name = $1`

	err := r.db.GetContext(ctx, &weight, query, fieldName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get weight for field %s: %w", fieldName, err)
	}

	return &weight, nil
}

func (r *SearchConfigRepository) CreateWeight(ctx context.Context, weight *domain.SearchWeight) error {
	query := `INSERT INTO search_weights (field_name, weight, description) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		weight.FieldName, weight.Weight, weight.Description).
		Scan(&weight.ID, &weight.CreatedAt, &weight.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create weight: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) UpdateWeight(ctx context.Context, id int64, weight *domain.SearchWeight) error {
	query := `UPDATE search_weights 
	          SET field_name = $2, weight = $3, description = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $1
	          RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		id, weight.FieldName, weight.Weight, weight.Description).
		Scan(&weight.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update weight: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) DeleteWeight(ctx context.Context, id int64) error {
	query := `DELETE FROM search_weights WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete weight: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Synonyms methods
func (r *SearchConfigRepository) GetSynonyms(ctx context.Context, language string) ([]domain.SearchSynonym, error) {
	var synonyms []domain.SearchSynonym
	query := `SELECT id, term, synonyms, language, created_at, updated_at 
	          FROM search_synonyms_config WHERE language = $1 ORDER BY term`

	rows, err := r.db.QueryContext(ctx, query, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get synonyms: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var synonym domain.SearchSynonym
		err := rows.Scan(&synonym.ID, &synonym.Term,
			pq.Array(&synonym.Synonyms), &synonym.Language,
			&synonym.CreatedAt, &synonym.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan synonym: %w", err)
		}
		synonyms = append(synonyms, synonym)
	}

	return synonyms, nil
}

func (r *SearchConfigRepository) CreateSynonym(ctx context.Context, synonym *domain.SearchSynonym) error {
	query := `INSERT INTO search_synonyms_config (term, synonyms, language) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		synonym.Term, pq.Array(synonym.Synonyms), synonym.Language).
		Scan(&synonym.ID, &synonym.CreatedAt, &synonym.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create synonym: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) UpdateSynonym(ctx context.Context, id int64, synonym *domain.SearchSynonym) error {
	query := `UPDATE search_synonyms_config 
	          SET term = $2, synonyms = $3, language = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $1
	          RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		id, synonym.Term, pq.Array(synonym.Synonyms), synonym.Language).
		Scan(&synonym.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update synonym: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) DeleteSynonym(ctx context.Context, id int64) error {
	query := `DELETE FROM search_synonyms_config WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete synonym: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Transliteration methods
func (r *SearchConfigRepository) GetTransliterationRules(ctx context.Context) ([]domain.TransliterationRule, error) {
	var rules []domain.TransliterationRule
	query := `SELECT id, from_script, to_script, from_pattern, to_pattern, priority, created_at, updated_at 
	          FROM transliteration_rules ORDER BY priority DESC, from_pattern`

	err := r.db.SelectContext(ctx, &rules, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transliteration rules: %w", err)
	}

	return rules, nil
}

func (r *SearchConfigRepository) CreateTransliterationRule(ctx context.Context, rule *domain.TransliterationRule) error {
	query := `INSERT INTO transliteration_rules (from_script, to_script, from_pattern, to_pattern, priority) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		rule.FromScript, rule.ToScript, rule.FromPattern, rule.ToPattern, rule.Priority).
		Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transliteration rule: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) UpdateTransliterationRule(ctx context.Context, id int64, rule *domain.TransliterationRule) error {
	query := `UPDATE transliteration_rules 
	          SET from_script = $2, to_script = $3, from_pattern = $4, to_pattern = $5, priority = $6, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $1
	          RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		id, rule.FromScript, rule.ToScript, rule.FromPattern, rule.ToPattern, rule.Priority).
		Scan(&rule.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update transliteration rule: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) DeleteTransliterationRule(ctx context.Context, id int64) error {
	query := `DELETE FROM transliteration_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transliteration rule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Statistics methods
func (r *SearchConfigRepository) CreateSearchStatistics(ctx context.Context, stats *domain.SearchStatistics) error {
	filtersJSON := ""
	if stats.SearchFilters != "" {
		filtersJSON = stats.SearchFilters
	}

	query := `INSERT INTO search_statistics (query, results_count, search_duration_ms, user_id, search_filters) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		stats.Query, stats.ResultsCount, stats.SearchDuration, stats.UserID, filtersJSON).
		Scan(&stats.ID, &stats.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create search statistics: %w", err)
	}

	return nil
}

func (r *SearchConfigRepository) GetSearchStatistics(ctx context.Context, limit int) ([]domain.SearchStatistics, error) {
	var statistics []domain.SearchStatistics
	query := `SELECT id, query, results_count, search_duration_ms, user_id, search_filters, created_at 
	          FROM search_statistics 
	          ORDER BY created_at DESC 
	          LIMIT $1`

	err := r.db.SelectContext(ctx, &statistics, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search statistics: %w", err)
	}

	return statistics, nil
}

func (r *SearchConfigRepository) GetPopularSearches(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `SELECT query, COUNT(*) as count, AVG(results_count) as avg_results
	          FROM search_statistics 
	          WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '7 days'
	          GROUP BY query 
	          ORDER BY count DESC 
	          LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		var avgResults float64

		err := rows.Scan(&query, &count, &avgResults)
		if err != nil {
			return nil, fmt.Errorf("failed to scan popular search: %w", err)
		}

		results = append(results, map[string]interface{}{
			"query":       query,
			"count":       count,
			"avg_results": avgResults,
		})
	}

	return results, nil
}

// Config methods
func (r *SearchConfigRepository) GetConfig(ctx context.Context) (*domain.SearchConfig, error) {
	var config domain.SearchConfig
	query := `SELECT id, min_search_length, max_suggestions, fuzzy_enabled, fuzzy_max_edits, 
	                 synonyms_enabled, transliteration_enabled, created_at, updated_at 
	          FROM search_config 
	          LIMIT 1`

	err := r.db.GetContext(ctx, &config, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get search config: %w", err)
	}

	return &config, nil
}

func (r *SearchConfigRepository) UpdateConfig(ctx context.Context, config *domain.SearchConfig) error {
	query := `UPDATE search_config 
	          SET min_search_length = $1, max_suggestions = $2, fuzzy_enabled = $3, 
	              fuzzy_max_edits = $4, synonyms_enabled = $5, transliteration_enabled = $6,
	              updated_at = CURRENT_TIMESTAMP
	          WHERE id = (SELECT id FROM search_config LIMIT 1)
	          RETURNING id, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		config.MinSearchLength, config.MaxSuggestions, config.FuzzyEnabled,
		config.FuzzyMaxEdits, config.SynonymsEnabled, config.TransliterationEnabled).
		Scan(&config.ID, &config.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update search config: %w", err)
	}

	return nil
}
