package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"backend/internal/domain/models"
)

// CategoryDetectionStatsRepository репозиторий для статистики определения категорий
type CategoryDetectionStatsRepository struct {
	db *sqlx.DB
}

// NewCategoryDetectionStatsRepository создает новый репозиторий
func NewCategoryDetectionStatsRepository(db *sqlx.DB) *CategoryDetectionStatsRepository {
	return &CategoryDetectionStatsRepository{db: db}
}

// Create создает новую запись статистики
func (r *CategoryDetectionStatsRepository) Create(ctx context.Context, stats *models.CategoryDetectionStats) error {
	query := `
		INSERT INTO category_detection_stats (
			user_id, session_id, method,
			ai_keywords, ai_attributes, ai_domain, ai_product_type,
			ai_suggested_category_id, final_category_id,
			alternative_categories, confidence_score,
			similarity_score, keyword_score,
			similar_listings_found, top_similar_listing_id, top_similarity_score,
			matched_keywords, matched_negative_keywords,
			processing_time_ms, user_confirmed, user_selected_category_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		) RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		stats.UserID, stats.SessionID, stats.Method,
		pq.Array(stats.AIKeywords), stats.AIAttributes, stats.AIDomain, stats.AIProductType,
		stats.AISuggestedCategoryID, stats.FinalCategoryID,
		stats.AlternativeCategories, stats.ConfidenceScore,
		stats.SimilarityScore, stats.KeywordScore,
		stats.SimilarListingsFound, stats.TopSimilarListingID, stats.TopSimilarityScore,
		pq.Array(stats.MatchedKeywords), pq.Array(stats.MatchedNegativeKeywords),
		stats.ProcessingTimeMs, stats.UserConfirmed, stats.UserSelectedCategoryID,
	).Scan(&stats.ID, &stats.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "ошибка создания записи статистики")
	}

	return nil
}

// UpdateUserFeedback обновляет обратную связь пользователя
func (r *CategoryDetectionStatsRepository) UpdateUserFeedback(ctx context.Context, statsID int32, confirmed bool, selectedCategoryID *int32) error {
	query := `
		UPDATE category_detection_stats
		SET user_confirmed = $2,
		    user_selected_category_id = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, statsID, confirmed, selectedCategoryID)
	if err != nil {
		return errors.Wrap(err, "ошибка обновления обратной связи")
	}

	return nil
}

// GetRecentStats получает статистику за последние N дней
func (r *CategoryDetectionStatsRepository) GetRecentStats(ctx context.Context, days int) ([]*models.CategoryDetectionStats, error) {
	query := `
		SELECT 
			id, user_id, session_id, method,
			ai_domain, ai_product_type,
			ai_suggested_category_id, final_category_id,
			confidence_score, similarity_score, keyword_score,
			similar_listings_found, top_similar_listing_id, top_similarity_score,
			processing_time_ms, user_confirmed, user_selected_category_id,
			created_at
		FROM category_detection_stats
		WHERE created_at > NOW() - INTERVAL '%d days'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(query, days))
	if err != nil {
		return nil, errors.Wrap(err, "ошибка выполнения запроса")
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Failed to close rows: %v", closeErr)
		}
	}()

	var stats []*models.CategoryDetectionStats
	for rows.Next() {
		stat := &models.CategoryDetectionStats{}
		err := rows.Scan(
			&stat.ID,
			&stat.UserID,
			&stat.SessionID,
			&stat.Method,
			&stat.AIDomain,
			&stat.AIProductType,
			&stat.AISuggestedCategoryID,
			&stat.FinalCategoryID,
			&stat.ConfidenceScore,
			&stat.SimilarityScore,
			&stat.KeywordScore,
			&stat.SimilarListingsFound,
			&stat.TopSimilarListingID,
			&stat.TopSimilarityScore,
			&stat.ProcessingTimeMs,
			&stat.UserConfirmed,
			&stat.UserSelectedCategoryID,
			&stat.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка сканирования строки")
		}

		// Отдельно загружаем массивы ключевых слов
		var matchedKeywords, matchedNegativeKeywords []string
		err = r.db.QueryRowContext(ctx,
			`SELECT matched_keywords, matched_negative_keywords 
			 FROM category_detection_stats 
			 WHERE id = $1`, stat.ID).Scan(
			pq.Array(&matchedKeywords),
			pq.Array(&matchedNegativeKeywords),
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "ошибка загрузки ключевых слов")
		}
		stat.MatchedKeywords = matchedKeywords
		stat.MatchedNegativeKeywords = matchedNegativeKeywords

		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "ошибка чтения результатов")
	}

	return stats, nil
}

// MethodEffectiveness эффективность метода определения
type MethodEffectiveness struct {
	Method             string  `db:"method"`
	TotalDetections    int32   `db:"total_detections"`
	Accuracy           float64 `db:"accuracy"`
	AvgConfidence      float64 `db:"avg_confidence"`
	AvgTimeMs          float64 `db:"avg_time_ms"`
	UserConfirmedCount int32   `db:"user_confirmed_count"`
	UserCorrectedCount int32   `db:"user_corrected_count"`
}

// GetMethodEffectiveness получает эффективность методов
func (r *CategoryDetectionStatsRepository) GetMethodEffectiveness(ctx context.Context) ([]*MethodEffectiveness, error) {
	query := `
		SELECT 
			method,
			COUNT(*) as total_detections,
			AVG(CASE WHEN ai_suggested_category_id = final_category_id THEN 1.0 ELSE 0.0 END) as accuracy,
			AVG(confidence_score) as avg_confidence,
			AVG(processing_time_ms) as avg_time_ms,
			COUNT(CASE WHEN user_confirmed = true THEN 1 END) as user_confirmed_count,
			COUNT(CASE WHEN user_selected_category_id IS NOT NULL THEN 1 END) as user_corrected_count
		FROM category_detection_stats
		WHERE created_at > NOW() - INTERVAL '30 days'
		GROUP BY method
		ORDER BY accuracy DESC
	`

	var effectiveness []*MethodEffectiveness
	err := r.db.SelectContext(ctx, &effectiveness, query)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения эффективности методов")
	}

	return effectiveness, nil
}

// ProblematicCategory проблемная категория
type ProblematicCategory struct {
	CategoryID        int32   `db:"id"`
	CategoryName      string  `db:"name"`
	CategorySlug      string  `db:"slug"`
	DetectionAttempts int32   `db:"detection_attempts"`
	SuccessRate       float64 `db:"success_rate"`
	AvgConfidence     float64 `db:"avg_confidence"`
	CorrectionCount   int32   `db:"correction_count"`
}

// GetProblematicCategories получает категории с низким процентом определения
func (r *CategoryDetectionStatsRepository) GetProblematicCategories(ctx context.Context, limit int) ([]*ProblematicCategory, error) {
	query := `
		SELECT 
			c.id,
			c.name,
			c.slug,
			COUNT(DISTINCT s.id) as detection_attempts,
			AVG(CASE WHEN s.ai_suggested_category_id = s.final_category_id THEN 1.0 ELSE 0.0 END) as success_rate,
			AVG(s.confidence_score) as avg_confidence,
			COUNT(CASE WHEN s.user_selected_category_id IS NOT NULL THEN 1 END) as correction_count
		FROM marketplace_categories c
		LEFT JOIN category_detection_stats s ON c.id = s.ai_suggested_category_id
		WHERE s.created_at > NOW() - INTERVAL '30 days'
		GROUP BY c.id, c.name, c.slug
		HAVING COUNT(DISTINCT s.id) > 5
		ORDER BY AVG(CASE WHEN s.ai_suggested_category_id = s.final_category_id THEN 1.0 ELSE 0.0 END) ASC
		LIMIT $1
	`

	var categories []*ProblematicCategory
	err := r.db.SelectContext(ctx, &categories, query, limit)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения проблемных категорий")
	}

	return categories, nil
}

// CategoryStats статистика по категории
type CategoryStats struct {
	CategoryID      int32           `json:"category_id"`
	TotalDetections int32           `json:"total_detections"`
	SuccessRate     float64         `json:"success_rate"`
	AvgConfidence   float64         `json:"avg_confidence"`
	TopKeywords     []string        `json:"top_keywords"`
	CommonMistakes  map[int32]int32 `json:"common_mistakes"` // category_id -> count
}

// GetCategoryStats получает детальную статистику по категории
func (r *CategoryDetectionStatsRepository) GetCategoryStats(ctx context.Context, categoryID int32, days int) (*CategoryStats, error) {
	// Основная статистика
	statsQuery := `
		SELECT 
			COUNT(*) as total_detections,
			AVG(CASE WHEN ai_suggested_category_id = final_category_id THEN 1.0 ELSE 0.0 END) as success_rate,
			AVG(confidence_score) as avg_confidence
		FROM category_detection_stats
		WHERE ai_suggested_category_id = $1
			AND created_at > NOW() - INTERVAL '%d days'
	`

	var stats CategoryStats
	stats.CategoryID = categoryID

	err := r.db.QueryRowContext(ctx, fmt.Sprintf(statsQuery, days), categoryID).Scan(
		&stats.TotalDetections,
		&stats.SuccessRate,
		&stats.AvgConfidence,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "ошибка получения статистики категории")
	}

	// Топ ключевых слов
	keywordsQuery := `
		SELECT unnest(matched_keywords) as keyword, COUNT(*) as count
		FROM category_detection_stats
		WHERE final_category_id = $1
			AND created_at > NOW() - INTERVAL '%d days'
		GROUP BY keyword
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(keywordsQuery, days), categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения топ ключевых слов")
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Failed to close rows: %v", closeErr)
		}
	}()

	stats.TopKeywords = make([]string, 0)
	for rows.Next() {
		var keyword string
		var count int
		if err := rows.Scan(&keyword, &count); err != nil {
			continue
		}
		stats.TopKeywords = append(stats.TopKeywords, keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "ошибка чтения результатов ключевых слов")
	}

	// Частые ошибки (когда AI предложил эту категорию, но пользователь выбрал другую)
	mistakesQuery := `
		SELECT user_selected_category_id, COUNT(*) as count
		FROM category_detection_stats
		WHERE ai_suggested_category_id = $1
			AND user_selected_category_id IS NOT NULL
			AND user_selected_category_id != ai_suggested_category_id
			AND created_at > NOW() - INTERVAL '%d days'
		GROUP BY user_selected_category_id
		ORDER BY count DESC
		LIMIT 5
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(mistakesQuery, days), categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения частых ошибок")
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Failed to close rows: %v", closeErr)
		}
	}()

	stats.CommonMistakes = make(map[int32]int32)
	for rows.Next() {
		var mistakeCategoryID int32
		var count int32
		if err := rows.Scan(&mistakeCategoryID, &count); err != nil {
			continue
		}
		stats.CommonMistakes[mistakeCategoryID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "ошибка чтения результатов ошибок")
	}

	return &stats, nil
}

// CleanupOldStats удаляет старую статистику
func (r *CategoryDetectionStatsRepository) CleanupOldStats(ctx context.Context, olderThanDays int) error {
	query := `
		DELETE FROM category_detection_stats
		WHERE created_at < NOW() - INTERVAL '%d days'
	`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, olderThanDays))
	if err != nil {
		return errors.Wrap(err, "ошибка удаления старой статистики")
	}

	return nil
}
