package postgres

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// CategoryKeywordRepository репозиторий для работы с ключевыми словами категорий
type CategoryKeywordRepository struct {
	db *sqlx.DB
}

// NewCategoryKeywordRepository создает новый репозиторий
func NewCategoryKeywordRepository(db *sqlx.DB) *CategoryKeywordRepository {
	return &CategoryKeywordRepository{db: db}
}

// KeywordMatch совпадение ключевого слова
type KeywordMatch struct {
	CategoryID  int32   `db:"category_id"`
	Keyword     string  `db:"keyword"`
	Weight      float64 `db:"weight"`
	KeywordType string  `db:"keyword_type"`
	IsNegative  bool    `db:"is_negative"`
}

// FindMatchingCategories находит категории по ключевым словам
func (r *CategoryKeywordRepository) FindMatchingCategories(ctx context.Context, keywords []string, language string) ([]*KeywordMatch, error) {
	if len(keywords) == 0 {
		return nil, nil
	}

	// Приводим ключевые слова к нижнему регистру
	lowerKeywords := make([]string, len(keywords))
	for i, kw := range keywords {
		lowerKeywords[i] = strings.ToLower(kw)
	}

	query := `
		SELECT 
			category_id,
			keyword,
			weight,
			keyword_type,
			is_negative
		FROM category_keywords
		WHERE LOWER(keyword) = ANY($1)
			AND (language = $2 OR language = '*')
		ORDER BY weight DESC
	`

	var matches []*KeywordMatch
	err := r.db.SelectContext(ctx, &matches, query, pq.Array(lowerKeywords), language)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка поиска категорий по ключевым словам")
	}

	return matches, nil
}

// UpdateSuccessRate обновляет процент успешного определения для ключевого слова
func (r *CategoryKeywordRepository) UpdateSuccessRate(ctx context.Context, keyword string, successRate float64) error {
	query := `
		UPDATE category_keywords
		SET success_rate = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE LOWER(keyword) = LOWER($1)
	`

	_, err := r.db.ExecContext(ctx, query, keyword, successRate)
	if err != nil {
		return errors.Wrap(err, "ошибка обновления success_rate")
	}

	return nil
}

// IncrementUsageCount увеличивает счетчик использования ключевых слов
func (r *CategoryKeywordRepository) IncrementUsageCount(ctx context.Context, categoryID int32, keywords []string, language string) error {
	if len(keywords) == 0 {
		return nil
	}

	// Приводим ключевые слова к нижнему регистру
	lowerKeywords := make([]string, len(keywords))
	for i, kw := range keywords {
		lowerKeywords[i] = strings.ToLower(kw)
	}

	query := `
		UPDATE category_keywords
		SET usage_count = usage_count + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE category_id = $1
			AND LOWER(keyword) = ANY($2)
			AND (language = $3 OR language = '*')
	`

	_, err := r.db.ExecContext(ctx, query, categoryID, pq.Array(lowerKeywords), language)
	if err != nil {
		return errors.Wrap(err, "ошибка увеличения usage_count")
	}

	return nil
}

// CategoryKeyword модель ключевого слова категории
type CategoryKeyword struct {
	ID          int32     `db:"id"`
	CategoryID  int32     `db:"category_id"`
	Keyword     string    `db:"keyword"`
	Language    string    `db:"language"`
	Weight      float64   `db:"weight"`
	KeywordType string    `db:"keyword_type"`
	IsNegative  bool      `db:"is_negative"`
	Source      string    `db:"source"`
	UsageCount  int32     `db:"usage_count"`
	SuccessRate float64   `db:"success_rate"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// GetKeywordsByCategoryID получает все ключевые слова для категории
func (r *CategoryKeywordRepository) GetKeywordsByCategoryID(ctx context.Context, categoryID int32) ([]*CategoryKeyword, error) {
	query := `
		SELECT * FROM category_keywords
		WHERE category_id = $1
		ORDER BY weight DESC, usage_count DESC
	`

	var keywords []*CategoryKeyword
	err := r.db.SelectContext(ctx, &keywords, query, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения ключевых слов категории")
	}

	return keywords, nil
}

// AddKeyword добавляет новое ключевое слово для категории
func (r *CategoryKeywordRepository) AddKeyword(ctx context.Context, keyword *CategoryKeyword) error {
	query := `
		INSERT INTO category_keywords (
			category_id, keyword, language, weight, 
			keyword_type, is_negative, source
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (category_id, keyword, language) 
		DO UPDATE SET
			weight = EXCLUDED.weight,
			keyword_type = EXCLUDED.keyword_type,
			is_negative = EXCLUDED.is_negative,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	err := r.db.GetContext(ctx, &keyword.ID, query,
		keyword.CategoryID, keyword.Keyword, keyword.Language,
		keyword.Weight, keyword.KeywordType, keyword.IsNegative, keyword.Source)
	if err != nil {
		return errors.Wrap(err, "ошибка добавления ключевого слова")
	}

	return nil
}

// DeleteKeyword удаляет ключевое слово
func (r *CategoryKeywordRepository) DeleteKeyword(ctx context.Context, id int32) error {
	query := `DELETE FROM category_keywords WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "ошибка удаления ключевого слова")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "ошибка получения количества удаленных строк")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
