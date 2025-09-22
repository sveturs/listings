package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// KeywordRepository handles keyword operations in database
type KeywordRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewKeywordRepository creates a new keyword repository
func NewKeywordRepository(db *sqlx.DB, logger *zap.Logger) *KeywordRepository {
	return &KeywordRepository{
		db:     db,
		logger: logger,
	}
}

// KeywordRecord represents a keyword record in database
type KeywordRecord struct {
	ID           int32     `db:"id"`
	CategoryID   int32     `db:"category_id"`
	Keyword      string    `db:"keyword"`
	Language     string    `db:"language"`
	Weight       float64   `db:"weight"`
	KeywordType  string    `db:"keyword_type"`
	IsNegative   bool      `db:"is_negative"`
	Source       string    `db:"source"`
	UsageCount   int32     `db:"usage_count"`
	SuccessRate  float64   `db:"success_rate"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// BulkInsertKeywords inserts multiple keywords for a category
func (r *KeywordRepository) BulkInsertKeywords(ctx context.Context, categoryID int32, keywords []GeneratedKeyword, source string) error {
	if len(keywords) == 0 {
		return nil
	}

	// Prepare batch insert query
	query := `
		INSERT INTO category_keywords
		(category_id, keyword, language, weight, keyword_type, source, created_at, updated_at)
		VALUES %s
		ON CONFLICT (category_id, keyword, language)
		DO UPDATE SET
			weight = EXCLUDED.weight,
			keyword_type = EXCLUDED.keyword_type,
			source = EXCLUDED.source,
			updated_at = CURRENT_TIMESTAMP
	`

	values := make([]string, 0, len(keywords))
	args := make([]interface{}, 0, len(keywords)*6)
	argIndex := 1

	for _, kw := range keywords {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4, argIndex+5))

		// Determine language - default to 'ru'
		language := "ru"
		if kw.Keyword != "" {
			// Simple heuristic for language detection
			if isEnglishText(kw.Keyword) {
				language = "en"
			}
		}

		args = append(args,
			categoryID,
			strings.ToLower(strings.TrimSpace(kw.Keyword)),
			language,
			kw.Weight,
			kw.Type,
			source,
		)
		argIndex += 6
	}

	finalQuery := fmt.Sprintf(query, strings.Join(values, ","))

	_, err := r.db.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		r.logger.Error("Failed to bulk insert keywords",
			zap.Int32("categoryId", categoryID),
			zap.Int("keywordCount", len(keywords)),
			zap.Error(err))
		return fmt.Errorf("failed to bulk insert keywords: %w", err)
	}

	r.logger.Info("Successfully bulk inserted keywords",
		zap.Int32("categoryId", categoryID),
		zap.Int("keywordCount", len(keywords)),
		zap.String("source", source))

	return nil
}

// GetKeywordsByCategory retrieves all keywords for a category
func (r *KeywordRepository) GetKeywordsByCategory(ctx context.Context, categoryID int32) ([]KeywordRecord, error) {
	query := `
		SELECT id, category_id, keyword, language, weight, keyword_type,
		       is_negative, source, usage_count, success_rate, created_at, updated_at
		FROM category_keywords
		WHERE category_id = $1
		ORDER BY weight DESC, usage_count DESC
	`

	var keywords []KeywordRecord
	err := r.db.SelectContext(ctx, &keywords, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords for category %d: %w", categoryID, err)
	}

	return keywords, nil
}

// GetKeywordCountByCategory returns keyword count for each category
func (r *KeywordRepository) GetKeywordCountByCategory(ctx context.Context) (map[int32]int, error) {
	query := `
		SELECT category_id, COUNT(*) as keyword_count
		FROM category_keywords
		GROUP BY category_id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyword counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[int32]int)
	for rows.Next() {
		var categoryID int32
		var count int
		if err := rows.Scan(&categoryID, &count); err != nil {
			continue
		}
		counts[categoryID] = count
	}

	return counts, nil
}

// SearchKeywords searches for keywords across categories
func (r *KeywordRepository) SearchKeywords(ctx context.Context, searchTerm string, limit int) ([]KeywordRecord, error) {
	query := `
		SELECT id, category_id, keyword, language, weight, keyword_type,
		       is_negative, source, usage_count, success_rate, created_at, updated_at
		FROM category_keywords
		WHERE LOWER(keyword) LIKE LOWER($1)
		ORDER BY weight DESC, usage_count DESC
		LIMIT $2
	`

	var keywords []KeywordRecord
	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	err := r.db.SelectContext(ctx, &keywords, query, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search keywords: %w", err)
	}

	return keywords, nil
}

// UpdateKeywordWeight updates the weight of a keyword
func (r *KeywordRepository) UpdateKeywordWeight(ctx context.Context, keywordID int32, newWeight float64) error {
	query := `
		UPDATE category_keywords
		SET weight = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, keywordID, newWeight)
	if err != nil {
		return fmt.Errorf("failed to update keyword weight: %w", err)
	}

	return nil
}

// UpdateKeywordStats updates usage statistics for a keyword
func (r *KeywordRepository) UpdateKeywordStats(ctx context.Context, keywordID int32, successful bool) error {
	query := `
		UPDATE category_keywords
		SET usage_count = usage_count + 1,
		    success_rate = CASE
		        WHEN usage_count = 0 THEN CASE WHEN $2 THEN 1.0 ELSE 0.0 END
		        ELSE (success_rate * usage_count + CASE WHEN $2 THEN 1.0 ELSE 0.0 END) / (usage_count + 1)
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, keywordID, successful)
	if err != nil {
		return fmt.Errorf("failed to update keyword stats: %w", err)
	}

	return nil
}

// GetCategoriesNeedingKeywords returns categories with less than minKeywords
func (r *KeywordRepository) GetCategoriesNeedingKeywords(ctx context.Context, minKeywords int) ([]CategoryInfo, error) {
	query := `
		SELECT c.id, c.name, c.slug, COALESCE(kw.keyword_count, 0) as current_keywords
		FROM marketplace_categories c
		LEFT JOIN (
			SELECT category_id, COUNT(*) as keyword_count
			FROM category_keywords
			GROUP BY category_id
		) kw ON c.id = kw.category_id
		WHERE COALESCE(kw.keyword_count, 0) < $1
		  AND c.parent_id IS NOT NULL  -- Only leaf categories
		ORDER BY COALESCE(kw.keyword_count, 0) ASC
	`

	var categories []CategoryInfo
	err := r.db.SelectContext(ctx, &categories, query, minKeywords)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories needing keywords: %w", err)
	}

	return categories, nil
}

// DeleteKeyword deletes a keyword by ID
func (r *KeywordRepository) DeleteKeyword(ctx context.Context, keywordID int32) error {
	query := `DELETE FROM category_keywords WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, keywordID)
	if err != nil {
		return fmt.Errorf("failed to delete keyword: %w", err)
	}

	return nil
}

// GetTopKeywords returns most successful keywords for analytics
func (r *KeywordRepository) GetTopKeywords(ctx context.Context, limit int) ([]KeywordAnalytics, error) {
	query := `
		SELECT ck.keyword, ck.keyword_type, ck.weight, ck.usage_count,
		       ck.success_rate, mc.name as category_name, ck.source
		FROM category_keywords ck
		JOIN marketplace_categories mc ON ck.category_id = mc.id
		WHERE ck.usage_count > 0
		ORDER BY ck.success_rate DESC, ck.usage_count DESC
		LIMIT $1
	`

	var analytics []KeywordAnalytics
	err := r.db.SelectContext(ctx, &analytics, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top keywords: %w", err)
	}

	return analytics, nil
}

// GetKeywordsByTypes returns keywords grouped by type for a category
func (r *KeywordRepository) GetKeywordsByTypes(ctx context.Context, categoryID int32) (map[string][]string, error) {
	query := `
		SELECT keyword_type, keyword
		FROM category_keywords
		WHERE category_id = $1
		ORDER BY keyword_type, weight DESC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords by type: %w", err)
	}
	defer rows.Close()

	keywordsByType := make(map[string][]string)
	for rows.Next() {
		var keywordType, keyword string
		if err := rows.Scan(&keywordType, &keyword); err != nil {
			continue
		}

		if _, exists := keywordsByType[keywordType]; !exists {
			keywordsByType[keywordType] = make([]string, 0)
		}
		keywordsByType[keywordType] = append(keywordsByType[keywordType], keyword)
	}

	return keywordsByType, nil
}

// Helper structs
type CategoryInfo struct {
	ID              int32  `db:"id"`
	Name            string `db:"name"`
	Slug            string `db:"slug"`
	CurrentKeywords int    `db:"current_keywords"`
}

type KeywordAnalytics struct {
	Keyword      string  `db:"keyword"`
	KeywordType  string  `db:"keyword_type"`
	Weight       float64 `db:"weight"`
	UsageCount   int32   `db:"usage_count"`
	SuccessRate  float64 `db:"success_rate"`
	CategoryName string  `db:"category_name"`
	Source       string  `db:"source"`
}

// GeneratedKeyword represents a keyword generated by AI (imported from services)
type GeneratedKeyword struct {
	Keyword     string  `json:"keyword"`
	Type        string  `json:"type"`
	Weight      float64 `json:"weight"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// isEnglishText is a simple heuristic to detect English text
func isEnglishText(text string) bool {
	// Count Latin characters vs total characters
	latinCount := 0
	totalChars := 0

	for _, r := range text {
		if r > 32 { // Skip spaces and control characters
			totalChars++
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				latinCount++
			}
		}
	}

	if totalChars == 0 {
		return false
	}

	// If more than 80% Latin characters, consider it English
	return float64(latinCount)/float64(totalChars) > 0.8
}