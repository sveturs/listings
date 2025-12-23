package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
)

// CategoryDetectionRepository - репозиторий для работы с category detection
type CategoryDetectionRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger

	// Кеш брендов для быстрого поиска
	brandCache      map[string]*domain.BrandCategoryMapping // brand_name (lowercase) -> mapping
	brandAliasCache map[string]string                        // alias (lowercase) -> brand_name
}

// NewCategoryDetectionRepository создаёт репозиторий
func NewCategoryDetectionRepository(db *sqlx.DB, logger zerolog.Logger) *CategoryDetectionRepository {
	repo := &CategoryDetectionRepository{
		db:              db,
		logger:          logger.With().Str("component", "category_detection_repo").Logger(),
		brandCache:      make(map[string]*domain.BrandCategoryMapping),
		brandAliasCache: make(map[string]string),
	}

	// Загрузить кеш брендов при инициализации
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := repo.loadBrandCache(ctx); err != nil {
		logger.Warn().Err(err).Msg("Failed to load brand cache on init, will retry on first use")
	}

	return repo
}

// FindByKeywords ищет категории по ключевым словам
func (r *CategoryDetectionRepository) FindByKeywords(
	ctx context.Context,
	keywords []string,
	language string,
) ([]domain.CategoryMatch, error) {
	if len(keywords) == 0 {
		return []domain.CategoryMatch{}, nil
	}

	// Создаём условия для поиска по keywords в meta_keywords JSONB
	// meta_keywords->>'language' содержит строку с keywords через запятую
	var conditions []string
	args := []interface{}{language}
	argIndex := 2

	for _, keyword := range keywords {
		conditions = append(conditions, fmt.Sprintf("meta_keywords->>$1 ILIKE $%d", argIndex))
		args = append(args, "%"+keyword+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " OR ")

	query := fmt.Sprintf(`
		SELECT
			id,
			slug,
			name,
			path,
			level
		FROM categories
		WHERE is_active = true
		  AND (%s)
		ORDER BY level ASC, sort_order ASC
		LIMIT 10
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to find categories by keywords")
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer rows.Close()

	var matches []domain.CategoryMatch

	for rows.Next() {
		var (
			id       uuid.UUID
			slug     string
			nameJSON []byte
			path     string
			level    int
		)

		if err := rows.Scan(&id, &slug, &nameJSON, &path, &level); err != nil {
			r.logger.Warn().Err(err).Msg("failed to scan category row")
			continue
		}

		// Parse name JSONB
		nameMap := make(map[string]string)
		if nameJSON != nil && len(nameJSON) > 0 {
			if err := json.Unmarshal(nameJSON, &nameMap); err != nil {
				r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
			}
		}

		// Get localized name
		categoryName := nameMap[language]
		if categoryName == "" {
			categoryName = nameMap["en"] // fallback to English
		}

		// Calculate confidence based on keyword matches
		matchedKeywords := []string{}
		for _, kw := range keywords {
			// Simple match check
			if strings.Contains(strings.ToLower(categoryName), strings.ToLower(kw)) {
				matchedKeywords = append(matchedKeywords, kw)
			}
		}

		confidence := float64(len(matchedKeywords)) / float64(len(keywords))
		if confidence < 0.3 {
			confidence = 0.3 // minimum confidence для keyword match
		}

		matches = append(matches, domain.CategoryMatch{
			CategoryID:      id,
			CategoryName:    categoryName,
			CategorySlug:    slug,
			CategoryPath:    path,
			ConfidenceScore: confidence,
			DetectionMethod: domain.MethodKeywordMatch,
			MatchedKeywords: matchedKeywords,
		})
	}

	return matches, nil
}

// FindBySimilarity ищет категории по похожести названия (PostgreSQL similarity)
func (r *CategoryDetectionRepository) FindBySimilarity(
	ctx context.Context,
	title string,
	language string,
) ([]domain.CategoryMatch, error) {
	if title == "" {
		return []domain.CategoryMatch{}, nil
	}

	// Используем pg_trgm similarity для поиска похожих категорий
	// Требует установленного расширения pg_trgm в PostgreSQL
	query := `
		SELECT
			id,
			slug,
			name,
			path,
			SIMILARITY(name->>$1, $2) as similarity_score
		FROM categories
		WHERE is_active = true
		  AND SIMILARITY(name->>$1, $2) > 0.2
		ORDER BY similarity_score DESC
		LIMIT 5
	`

	rows, err := r.db.QueryContext(ctx, query, language, title)
	if err != nil {
		r.logger.Warn().Err(err).Msg("similarity search failed (pg_trgm extension may not be installed)")
		return []domain.CategoryMatch{}, nil // не критично, вернём пустой результат
	}
	defer rows.Close()

	var matches []domain.CategoryMatch

	for rows.Next() {
		var (
			id              uuid.UUID
			slug            string
			nameJSON        []byte
			path            string
			similarityScore float64
		)

		if err := rows.Scan(&id, &slug, &nameJSON, &path, &similarityScore); err != nil {
			r.logger.Warn().Err(err).Msg("failed to scan similarity row")
			continue
		}

		// Parse name JSONB
		nameMap := make(map[string]string)
		if nameJSON != nil && len(nameJSON) > 0 {
			if err := json.Unmarshal(nameJSON, &nameMap); err != nil {
				r.logger.Warn().Err(err).Msg("failed to unmarshal name JSONB")
			}
		}

		categoryName := nameMap[language]
		if categoryName == "" {
			categoryName = nameMap["en"]
		}

		matches = append(matches, domain.CategoryMatch{
			CategoryID:      id,
			CategoryName:    categoryName,
			CategorySlug:    slug,
			CategoryPath:    path,
			ConfidenceScore: similarityScore,
			DetectionMethod: domain.MethodSimilarity,
		})
	}

	return matches, nil
}

// SaveDetection сохраняет результаты детекции для tracking и обучения
func (r *CategoryDetectionRepository) SaveDetection(
	ctx context.Context,
	detection *domain.CategoryDetection,
) error {
	// Пока не храним в БД - можно добавить таблицу category_detections
	// для tracking и улучшения алгоритма
	r.logger.Debug().
		Str("detection_id", detection.ID.String()).
		Str("input_title", detection.InputTitle).
		Msg("category detection tracked (not saved to DB)")

	return nil
}

// ConfirmSelection записывает подтверждение выбора пользователя
func (r *CategoryDetectionRepository) ConfirmSelection(
	ctx context.Context,
	detectionID, selectedCategoryID uuid.UUID,
) error {
	// Можно добавить таблицу для tracking подтверждений
	// Это поможет улучшить алгоритм со временем
	r.logger.Info().
		Str("detection_id", detectionID.String()).
		Str("selected_category_id", selectedCategoryID.String()).
		Msg("user selection confirmed (not saved to DB)")

	return nil
}

// loadBrandCache загружает все бренды в память для быстрого поиска
func (r *CategoryDetectionRepository) loadBrandCache(ctx context.Context) error {
	query := `
		SELECT id, brand_name, brand_aliases, category_slug, confidence, is_verified, created_at, updated_at
		FROM brand_category_mapping
		WHERE is_verified = true
		ORDER BY confidence DESC
	`

	var mappings []struct {
		ID           uuid.UUID      `db:"id"`
		BrandName    string         `db:"brand_name"`
		BrandAliases pq.StringArray `db:"brand_aliases"`
		CategorySlug string         `db:"category_slug"`
		Confidence   float64        `db:"confidence"`
		IsVerified   bool           `db:"is_verified"`
		CreatedAt    time.Time      `db:"created_at"`
		UpdatedAt    time.Time      `db:"updated_at"`
	}

	if err := r.db.SelectContext(ctx, &mappings, query); err != nil {
		return fmt.Errorf("failed to load brand mappings: %w", err)
	}

	// Очистить старый кеш
	r.brandCache = make(map[string]*domain.BrandCategoryMapping)
	r.brandAliasCache = make(map[string]string)

	// Заполнить кеш
	for _, m := range mappings {
		brandNameLower := strings.ToLower(m.BrandName)

		mapping := &domain.BrandCategoryMapping{
			ID:           m.ID,
			BrandName:    m.BrandName,
			BrandAliases: []string(m.BrandAliases),
			CategorySlug: m.CategorySlug,
			Confidence:   m.Confidence,
			IsVerified:   m.IsVerified,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}

		r.brandCache[brandNameLower] = mapping

		// Добавить все алиасы в кеш
		for _, alias := range m.BrandAliases {
			aliasLower := strings.ToLower(alias)
			r.brandAliasCache[aliasLower] = brandNameLower
		}
	}

	r.logger.Info().
		Int("brands_loaded", len(r.brandCache)).
		Int("aliases_loaded", len(r.brandAliasCache)).
		Msg("Brand cache loaded successfully")

	return nil
}

// GetBrandCategoryMapping возвращает маппинг бренда на категорию
func (r *CategoryDetectionRepository) GetBrandCategoryMapping(ctx context.Context, brand string) (*domain.BrandCategoryMapping, error) {
	brandLower := strings.ToLower(brand)

	// Проверить кеш
	if len(r.brandCache) == 0 {
		// Кеш пуст, загрузить
		if err := r.loadBrandCache(ctx); err != nil {
			return nil, fmt.Errorf("failed to load brand cache: %w", err)
		}
	}

	// Попробовать найти по прямому названию
	if mapping, ok := r.brandCache[brandLower]; ok {
		return mapping, nil
	}

	// Попробовать найти по алиасу
	if brandName, ok := r.brandAliasCache[brandLower]; ok {
		if mapping, ok := r.brandCache[brandName]; ok {
			return mapping, nil
		}
	}

	return nil, fmt.Errorf("brand not found: %s", brand)
}

// FindBrandInText ищет бренды в тексте (использует кеш)
// Использует word boundary matching для избежания ложных срабатываний
func (r *CategoryDetectionRepository) FindBrandInText(ctx context.Context, text string) ([]*domain.BrandCategoryMapping, error) {
	// Загрузить кеш если пуст
	if len(r.brandCache) == 0 {
		if err := r.loadBrandCache(ctx); err != nil {
			return nil, fmt.Errorf("failed to load brand cache: %w", err)
		}
	}

	textLower := strings.ToLower(text)
	var matches []*domain.BrandCategoryMapping
	seen := make(map[string]bool)

	// Проверить все алиасы с word boundary matching
	for alias, brandName := range r.brandAliasCache {
		if containsWholeWord(textLower, alias) {
			if !seen[brandName] {
				seen[brandName] = true
				if mapping, ok := r.brandCache[brandName]; ok {
					matches = append(matches, mapping)
				}
			}
		}
	}

	// Проверить прямые названия брендов
	for brandName, mapping := range r.brandCache {
		if !seen[brandName] && containsWholeWord(textLower, brandName) {
			seen[brandName] = true
			matches = append(matches, mapping)
		}
	}

	return matches, nil
}

// containsWholeWord проверяет наличие слова/фразы с учётом границ слов
// Избегает ложных срабатываний типа "th" в "with"
func containsWholeWord(text, word string) bool {
	if word == "" {
		return false
	}

	// Для коротких алиасов (2-3 символа) требуем точное совпадение слова
	// Для длинных - можно использовать substring matching
	minLengthForSubstring := 4

	if len(word) < minLengthForSubstring {
		// Ищем как отдельное слово с границами
		// Границы: начало строки, пробел, знаки препинания
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
		return pattern.MatchString(text)
	}

	// Для длинных слов (4+ символов) substring matching безопасен
	return strings.Contains(text, word)
}

// FindByAIMapping ищет категорию по AI-предложенному названию через таблицу маппинга
// Возвращает CategoryMatch с высоким confidence (0.95+) при успехе
func (r *CategoryDetectionRepository) FindByAIMapping(
	ctx context.Context,
	suggestedCategory string,
	language string,
) (*domain.CategoryMatch, error) {
	if suggestedCategory == "" {
		return nil, nil
	}

	// Нормализация: lowercase, trim
	normalized := strings.ToLower(strings.TrimSpace(suggestedCategory))

	query := `
		SELECT
			c.id,
			c.slug,
			c.name,
			c.path,
			c.level,
			cam.confidence_boost,
			cam.priority
		FROM category_ai_mapping cam
		JOIN categories c ON cam.target_category_id = c.id
		WHERE LOWER(cam.ai_category_name) = $1
		  AND c.is_active = true
		ORDER BY cam.priority DESC
		LIMIT 1
	`

	var (
		categoryID      uuid.UUID
		slug            string
		nameJSON        []byte
		path            string
		level           int
		confidenceBoost float64
		priority        int
	)

	err := r.db.QueryRowContext(ctx, query, normalized).Scan(
		&categoryID,
		&slug,
		&nameJSON,
		&path,
		&level,
		&confidenceBoost,
		&priority,
	)

	if err != nil {
		// sql.ErrNoRows не ошибка - просто нет маппинга
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Debug().
				Str("suggested", suggestedCategory).
				Msg("No AI mapping found (not an error)")
			return nil, nil
		}
		return nil, fmt.Errorf("query AI mapping failed: %w", err)
	}

	// Parse name JSONB → map[language]name
	var nameMap map[string]string
	if err := json.Unmarshal(nameJSON, &nameMap); err != nil {
		return nil, fmt.Errorf("failed to parse category name JSONB: %w", err)
	}

	// Локализация: sr → ru → en → first available
	categoryName := nameMap[language]
	if categoryName == "" {
		if alt := nameMap["sr"]; alt != "" {
			categoryName = alt
		} else if alt := nameMap["ru"]; alt != "" {
			categoryName = alt
		} else if alt := nameMap["en"]; alt != "" {
			categoryName = alt
		} else {
			// Берём первое доступное
			for _, name := range nameMap {
				categoryName = name
				break
			}
		}
	}

	match := &domain.CategoryMatch{
		CategoryID:      categoryID,
		CategoryName:    categoryName,
		CategorySlug:    slug,
		CategoryPath:    path,
		ConfidenceScore: 0.80 + confidenceBoost, // Base 0.80 + boost (0.15-0.30) = 0.95-1.10 (cap at 1.0)
		DetectionMethod: "ai_mapping",
		MatchedKeywords: []string{suggestedCategory},
	}

	// Cap confidence at 1.0
	if match.ConfidenceScore > 1.0 {
		match.ConfidenceScore = 1.0
	}

	r.logger.Info().
		Str("suggested", suggestedCategory).
		Str("category", slug).
		Float64("confidence", match.ConfidenceScore).
		Int("priority", priority).
		Msg("AI mapping found")

	return match, nil
}
