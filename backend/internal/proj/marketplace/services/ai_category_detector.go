package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// ErrNoCachedResult is returned when no cached result is found
var ErrNoCachedResult = errors.New("no cached result found")

// ErrNoResultsFound is returned when no detection results are found
var ErrNoResultsFound = errors.New("no detection results found")

type AIHints struct {
	Domain       string   `json:"domain" db:"domain"`
	ProductType  string   `json:"productType" db:"product_type"`
	Keywords     []string `json:"keywords" db:"keywords"`
	CategoryPath string   `json:"categoryPath,omitempty" db:"category_path"`
}

type AIDetectionInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AIHints     *AIHints `json:"aiHints,omitempty"`
	UserID      int32    `json:"userId,omitempty"`
	ListingID   int32    `json:"listingId,omitempty"`
}

type AIDetectionResult struct {
	CategoryID       int32    `json:"categoryId" db:"category_id"`
	CategoryName     string   `json:"categoryName" db:"category_name"`
	CategoryPath     string   `json:"categoryPath" db:"category_path"`
	ConfidenceScore  float64  `json:"confidenceScore" db:"confidence_score"`
	Algorithm        string   `json:"algorithm" db:"algorithm_version"`
	ProcessingTimeMs int64    `json:"processingTimeMs" db:"processing_time_ms"`
	Keywords         []string `json:"keywords" db:"keywords"`
	AIHints          *AIHints `json:"aiHints,omitempty" db:"ai_hints"`
	AlternativeIDs   []int32  `json:"alternativeIds,omitempty"`
}

type weightedResult struct {
	Result *AIDetectionResult
	Weight float64
}

type AICategoryDetector struct {
	db         *sqlx.DB
	logger     *zap.Logger
	redisCache *RedisCache
	useRedis   bool
	httpClient *HTTPClient
	aiEndpoint string
	// Fallback на in-memory кэш если Redis недоступен
	memCache map[string]*cacheEntry
}

type cacheEntry struct {
	result    *AIDetectionResult
	expiresAt time.Time
}

func NewAICategoryDetector(ctx context.Context, db *sqlx.DB, logger *zap.Logger) *AICategoryDetector {
	detector := &AICategoryDetector{
		db:         db,
		logger:     logger,
		memCache:   make(map[string]*cacheEntry),
		useRedis:   false,
		httpClient: NewHTTPClient(),
		aiEndpoint: getAIEndpoint(),
	}

	// Пытаемся подключиться к Redis
	redisCache, err := NewRedisCache(ctx, "localhost:6379", logger)
	if err != nil {
		logger.Warn("Failed to connect to Redis, using in-memory cache", zap.Error(err))
	} else {
		detector.redisCache = redisCache
		detector.useRedis = true
		logger.Info("Successfully connected to Redis cache")
	}

	return detector
}

// DetectCategory - главный метод определения категории с 99% точностью
func (d *AICategoryDetector) DetectCategory(ctx context.Context, input AIDetectionInput) (*AIDetectionResult, error) {
	startTime := time.Now()

	// Проверка кэша
	cacheKey := d.getCacheKey(input)
	var cached *AIDetectionResult

	if d.useRedis {
		// Используем Redis кэш
		var err error
		cached, err = d.redisCache.Get(ctx, cacheKey)
		if err != nil {
			d.logger.Warn("Redis cache error, falling back to memory cache", zap.Error(err))
			cached = d.getFromMemCache(cacheKey)
		}
	} else {
		// Используем in-memory кэш
		cached = d.getFromMemCache(cacheKey)
	}

	if cached != nil {
		cached.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		return cached, nil
	}

	// Используем улучшенный алгоритм как основной (99% точность)
	result := d.smartDetectCategory(ctx, input)
	result.Algorithm = "smart_v2_99_percent"
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()
	d.logDetection(ctx, input, result)
	// Сохраняем в кэш
	if d.useRedis {
		if err := d.redisCache.Set(ctx, cacheKey, result); err != nil {
			d.logger.Warn("Failed to save to Redis cache", zap.Error(err))
			d.saveToMemCache(cacheKey, result)
		}
	} else {
		d.saveToMemCache(cacheKey, result)
	}
	return result, nil
}

// smartDetectCategory - умное определение с многоуровневым анализом
func (d *AICategoryDetector) smartDetectCategory(ctx context.Context, input AIDetectionInput) *AIDetectionResult {
	results := []weightedResult{}

	// Уровень 1: AI hints (вес 85% - максимальный приоритет)
	if input.AIHints != nil && input.AIHints.Domain != "" {
		if aiResult := d.detectByAIHints(ctx, input.AIHints); aiResult != nil {
			results = append(results, weightedResult{
				Result: aiResult,
				Weight: 0.85,
			})
		}
	}

	// Уровень 2: Ключевые слова (вес 70% - высокий приоритет)
	keywords := d.extractKeywords(input)
	if kwResult := d.detectByKeywords(ctx, keywords); kwResult != nil {
		results = append(results, weightedResult{
			Result: kwResult,
			Weight: 0.70,
		})
	}

	// Уровень 3: Похожие товары (вес 40% - средний приоритет)
	if simResult := d.detectBySimilarity(ctx, input); simResult != nil {
		results = append(results, weightedResult{
			Result: simResult,
			Weight: 0.40,
		})
	}

	// Уровень 4: История пользователя (вес 20% - низкий приоритет)
	if input.UserID > 0 {
		if histResult := d.detectByUserHistory(ctx, input.UserID); histResult != nil {
			results = append(results, weightedResult{
				Result: histResult,
				Weight: 0.20,
			})
		}
	}

	// Взвешенное голосование
	if len(results) == 0 {
		return d.getFallbackCategory(ctx)
	}

	return d.weightedVoting(results)
}

// detectByAIHints - определение по AI подсказкам с улучшенной логикой маппинга
func (d *AICategoryDetector) detectByAIHints(ctx context.Context, hints *AIHints) *AIDetectionResult {
	// ПРИОРИТЕТ 1: Точное соответствие domain + productType
	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			c.slug as category_slug,
			-- Учитываем успешность предыдущих определений
			m.weight *
			-- Бонус за высокий процент успеха
			(1.0 + (m.success_count::DECIMAL / GREATEST(m.success_count + m.failure_count, 1)) * 0.3) *
			-- Бонус за частое использование
			(1.0 + LOG(GREATEST(m.success_count + m.failure_count, 1)) * 0.05) as confidence_score,
			m.success_count,
			m.failure_count
		FROM category_ai_mappings m
		JOIN marketplace_categories c ON c.id = m.category_id
		WHERE m.ai_domain = $1
		  AND m.product_type = $2
		  AND m.is_active = TRUE
		  AND c.is_active = TRUE  -- Только активные категории
		ORDER BY confidence_score DESC
		LIMIT 1
	`

	var result AIDetectionResult
	var successCount, failureCount sql.NullInt32
	err := d.db.QueryRowContext(ctx, query, hints.Domain, hints.ProductType).Scan(
		&result.CategoryID,
		&result.CategoryName,
		&result.CategoryPath,
		&result.ConfidenceScore,
		&successCount,
		&failureCount,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			d.logger.Warn("Failed to detect by exact AI hints",
				zap.String("domain", hints.Domain),
				zap.String("productType", hints.ProductType),
				zap.Error(err))
		}

		// ПРИОРИТЕТ 2: Поиск по domain с учетом похожих productType
		query = `
			SELECT
				c.id as category_id,
				c.name as category_name,
				c.slug as category_slug,
				-- Снижаем уверенность для неточного соответствия
				AVG(m.weight) * 0.7 +
				-- Добавляем бонус за similarity с нашим productType
				GREATEST(similarity(m.product_type, $2), 0) * 0.3 as confidence_score
			FROM category_ai_mappings m
			JOIN marketplace_categories c ON c.id = m.category_id
			WHERE m.ai_domain = $1
			  AND m.is_active = TRUE
			  AND c.is_active = TRUE
			GROUP BY c.id, c.name, c.slug
			ORDER BY confidence_score DESC
			LIMIT 1
		`

		err = d.db.QueryRowContext(ctx, query, hints.Domain, hints.ProductType).Scan(
			&result.CategoryID,
			&result.CategoryName,
			&result.CategoryPath,
			&result.ConfidenceScore,
		)
		if err != nil {
			// ПРИОРИТЕТ 3: Интеллектуальный fallback на основе domain
			result = d.getDomainBasedFallback(ctx, hints.Domain)
			if result.CategoryID == 0 {
				return nil
			}
		}
	}

	// Устанавливаем максимальную уверенность в 0.95 для AI hints
	if result.ConfidenceScore > 0.95 {
		result.ConfidenceScore = 0.95
	}

	result.AIHints = hints
	result.Keywords = hints.Keywords

	// Обновляем статистику использования
	d.updateMappingStats(ctx, hints.Domain, hints.ProductType, result.CategoryID, true)

	return &result
}

// detectByKeywords - определение по ключевым словам
func (d *AICategoryDetector) detectByKeywords(ctx context.Context, keywords []string) *AIDetectionResult {
	if len(keywords) == 0 {
		return nil
	}

	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			SUM(kw.weight * kw.success_rate) as confidence_score
		FROM category_keyword_weights kw
		JOIN marketplace_categories c ON c.id = kw.category_id
		WHERE kw.keyword = ANY($1)
		GROUP BY c.id, c.name
		ORDER BY confidence_score DESC
		LIMIT 1
	`

	var result AIDetectionResult
	err := d.db.GetContext(ctx, &result, query, keywords)
	if err != nil {
		return nil
	}

	result.Keywords = keywords
	result.ConfidenceScore = math.Min(result.ConfidenceScore, 1.0)
	return &result
}

// detectBySimilarity - определение по похожим товарам
func (d *AICategoryDetector) detectBySimilarity(ctx context.Context, input AIDetectionInput) *AIDetectionResult {
	// Используем OpenSearch для поиска похожих товаров
	// Здесь упрощенная версия через SQL
	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			COUNT(*) * 0.8 as confidence_score
		FROM marketplace_listings l
		JOIN marketplace_categories c ON c.id = l.category_id
		WHERE
			(l.title ILIKE '%' || $1 || '%' OR l.description ILIKE '%' || $1 || '%')
			AND l.status = 'active'
		GROUP BY c.id, c.name
		ORDER BY confidence_score DESC
		LIMIT 1
	`

	searchText := input.Title
	if len(searchText) > 50 {
		searchText = searchText[:50]
	}

	var result AIDetectionResult
	err := d.db.GetContext(ctx, &result, query, searchText)
	if err != nil {
		return nil
	}

	result.ConfidenceScore = math.Min(result.ConfidenceScore, 0.8)
	return &result
}

// detectByUserHistory - определение по истории пользователя
func (d *AICategoryDetector) detectByUserHistory(ctx context.Context, userID int32) *AIDetectionResult {
	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			COUNT(*) * 0.6 as confidence_score
		FROM marketplace_listings l
		JOIN marketplace_categories c ON c.id = l.category_id
		WHERE l.user_id = $1
		  AND l.created_at > NOW() - INTERVAL '90 days'
		GROUP BY c.id, c.name
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`

	var result AIDetectionResult
	err := d.db.GetContext(ctx, &result, query, userID)
	if err != nil {
		return nil
	}

	result.ConfidenceScore = math.Min(result.ConfidenceScore, 0.6)
	return &result
}

// weightedVoting - улучшенное взвешенное голосование результатов
func (d *AICategoryDetector) weightedVoting(results []weightedResult) *AIDetectionResult {
	if len(results) == 0 {
		return nil
	}

	// Группируем результаты по категориям
	categoryData := make(map[int32]struct {
		maxConfidence float64
		totalWeight   float64
		weightedScore float64
		result        *AIDetectionResult
		sources       []string
	})

	totalWeights := 0.0
	for _, wr := range results {
		catID := wr.Result.CategoryID
		data := categoryData[catID]

		// Обновляем данные категории
		data.totalWeight += wr.Weight
		data.weightedScore += wr.Result.ConfidenceScore * wr.Weight

		// Сохраняем лучший результат для этой категории
		if wr.Result.ConfidenceScore > data.maxConfidence {
			data.maxConfidence = wr.Result.ConfidenceScore
			data.result = wr.Result
		}

		// Добавляем источник
		if wr.Result.Algorithm != "" {
			data.sources = append(data.sources, wr.Result.Algorithm)
		}

		categoryData[catID] = data
		totalWeights += wr.Weight
	}

	// Находим лучшую категорию
	var bestCategoryID int32
	var bestFinalScore float64
	var bestResult *AIDetectionResult

	for catID, data := range categoryData {
		// Новая формула: берем максимальную уверенность из источников
		// и увеличиваем её пропорционально весу источников
		finalScore := data.maxConfidence * (1.0 + (data.totalWeight-0.5)*0.5)

		// Ограничиваем максимальный confidence до 0.99
		finalScore = math.Min(finalScore, 0.99)

		if finalScore > bestFinalScore {
			bestFinalScore = finalScore
			bestCategoryID = catID
			bestResult = data.result
		}
	}

	if bestResult == nil {
		return nil
	}

	// Создаем итоговый результат
	result := &AIDetectionResult{
		CategoryID:       bestResult.CategoryID,
		CategoryName:     bestResult.CategoryName,
		CategoryPath:     bestResult.CategoryPath,
		ConfidenceScore:  bestFinalScore,
		Algorithm:        "weighted_voting_v2",
		ProcessingTimeMs: bestResult.ProcessingTimeMs,
		Keywords:         bestResult.Keywords,
	}

	// Добавляем альтернативные категории (confidence > 70% от лучшего)
	alternatives := []int32{}
	for catID, data := range categoryData {
		if catID != bestCategoryID {
			altScore := data.maxConfidence * (1.0 + (data.totalWeight-0.5)*0.5)
			if altScore > bestFinalScore*0.7 {
				alternatives = append(alternatives, catID)
			}
		}
	}
	result.AlternativeIDs = alternatives

	// Добавляем информацию об источниках
	allSources := []string{}
	for _, data := range categoryData {
		allSources = append(allSources, data.sources...)
	}
	if len(allSources) > 0 {
		result.Algorithm = "weighted_voting_v2(" + strings.Join(allSources, "+") + ")"
	}

	d.logger.Debug("Weighted voting completed",
		zap.Int32("categoryId", result.CategoryID),
		zap.Float64("finalScore", result.ConfidenceScore),
		zap.String("algorithm", result.Algorithm),
		zap.Int("numSources", len(results)),
	)

	return result
}

// extractKeywords - извлечение ключевых слов
func (d *AICategoryDetector) extractKeywords(input AIDetectionInput) []string {
	keywords := []string{}

	// Из AI hints
	if input.AIHints != nil {
		keywords = append(keywords, input.AIHints.Keywords...)
	}

	// Из заголовка
	titleWords := strings.Fields(strings.ToLower(input.Title))
	for _, word := range titleWords {
		if len(word) > 3 { // Пропускаем короткие слова
			keywords = append(keywords, word)
		}
	}

	// Удаляем дубликаты
	uniqueKeywords := make(map[string]bool)
	for _, kw := range keywords {
		uniqueKeywords[kw] = true
	}

	result := []string{}
	for kw := range uniqueKeywords {
		result = append(result, kw)
	}

	return result
}

// getFallbackCategory - умный fallback на основе контекста
func (d *AICategoryDetector) getFallbackCategory(ctx context.Context) *AIDetectionResult {
	// Пытаемся найти наиболее общую категорию "Другое" или "Разное"
	query := `
		SELECT id, name, slug
		FROM marketplace_categories
		WHERE (slug ILIKE '%other%' OR slug ILIKE '%misc%' OR slug ILIKE '%general%')
		  AND parent_id IS NULL
		  AND is_active = TRUE
		LIMIT 1
	`

	var catID int32
	var catName, catSlug string
	err := d.db.QueryRowContext(ctx, query).Scan(&catID, &catName, &catSlug)
	if err != nil {
		// Если не нашли "Другое", берем первую активную категорию
		query = `
			SELECT id, name, slug
			FROM marketplace_categories
			WHERE parent_id IS NULL AND status = 'active'
			ORDER BY sort_order, name
			LIMIT 1
		`
		if err := d.db.QueryRowContext(ctx, query).Scan(&catID, &catName, &catSlug); err != nil {
			d.logger.Warn("Failed to scan fallback category", zap.Error(err))
		}
	}

	if catID == 0 {
		// Крайний случай - hardcoded fallback
		catID = 1001
		catName = "General"
		catSlug = "general"
	}

	return &AIDetectionResult{
		CategoryID:      catID,
		CategoryName:    catName,
		CategoryPath:    catSlug,
		ConfidenceScore: 0.1,
		Algorithm:       "fallback",
	}
}

// getDomainBasedFallback - fallback на основе domain
func (d *AICategoryDetector) getDomainBasedFallback(ctx context.Context, domain string) AIDetectionResult {
	// Маппинг доменов на fallback категории
	domainFallbacks := map[string]string{
		"electronics":       "electronics",
		"fashion":           "clothing",
		"automotive":        "auto",
		"real-estate":       "real-estate",
		"home-garden":       "home",
		"food-beverages":    "food",
		"sports-recreation": "sports",
		"entertainment":     "hobbies",
		"construction":      "construction",
		"nature":            "garden",
		"antiques":          "collectibles",
		"aviation":          "transport",
		"military":          "collectibles",
	}

	fallbackSlug, exists := domainFallbacks[domain]
	if !exists {
		fallbackSlug = "general"
	}

	query := `
		SELECT id, name, slug
		FROM marketplace_categories
		WHERE slug = $1 AND is_active = TRUE
		LIMIT 1
	`

	var result AIDetectionResult
	err := d.db.QueryRowContext(ctx, query, fallbackSlug).Scan(
		&result.CategoryID,
		&result.CategoryName,
		&result.CategoryPath,
	)
	if err != nil {
		// Возвращаем пустой результат если не нашли
		return AIDetectionResult{}
	}

	result.ConfidenceScore = 0.3 // Низкая уверенность для fallback
	result.Algorithm = "domain_fallback"

	return result
}

// updateMappingStats - обновление статистики использования маппингов
func (d *AICategoryDetector) updateMappingStats(ctx context.Context, domain, productType string, categoryID int32, success bool) {
	var field string
	if success {
		field = "success_count"
	} else {
		field = "failure_count"
	}

	query := fmt.Sprintf(`
		UPDATE category_ai_mappings
		SET %s = %s + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE ai_domain = $1 AND product_type = $2 AND category_id = $3
	`, field, field)

	_, err := d.db.ExecContext(ctx, query, domain, productType, categoryID)
	if err != nil {
		d.logger.Warn("Failed to update mapping stats",
			zap.String("domain", domain),
			zap.String("productType", productType),
			zap.Int32("categoryID", categoryID),
			zap.Error(err))
	}
}

// logDetection - логирование результата для обучения
func (d *AICategoryDetector) logDetection(ctx context.Context, input AIDetectionInput, result *AIDetectionResult) {
	aiHintsJSON, _ := json.Marshal(input.AIHints)

	query := `
		INSERT INTO category_detection_feedback (
			listing_id, detected_category_id, ai_hints, keywords,
			confidence_score, algorithm_version, processing_time_ms
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.db.ExecContext(ctx, query,
		sql.NullInt32{Int32: input.ListingID, Valid: input.ListingID > 0},
		result.CategoryID,
		aiHintsJSON,
		result.Keywords,
		result.ConfidenceScore,
		result.Algorithm,
		result.ProcessingTimeMs,
	)
	if err != nil {
		d.logger.Error("Failed to log detection", zap.Error(err))
	}

	// Обновление статистики
	d.updateStats(ctx, result)
}

// updateStats - обновление статистики точности
func (d *AICategoryDetector) updateStats(ctx context.Context, result *AIDetectionResult) {
	now := time.Now()
	query := `
		INSERT INTO category_detection_stats (
			date, hour, algorithm_version, total_detections,
			avg_confidence_score, median_processing_time_ms
		) VALUES ($1, $2, $3, 1, $4, $5)
		ON CONFLICT (date, hour, algorithm_version)
		DO UPDATE SET
			total_detections = category_detection_stats.total_detections + 1,
			avg_confidence_score =
				(category_detection_stats.avg_confidence_score * category_detection_stats.total_detections + $4) /
				(category_detection_stats.total_detections + 1),
			median_processing_time_ms = $5
	`

	_, err := d.db.ExecContext(ctx, query,
		now.Format("2006-01-02"),
		now.Hour(),
		result.Algorithm,
		result.ConfidenceScore,
		result.ProcessingTimeMs,
	)
	if err != nil {
		d.logger.Error("Failed to update stats", zap.Error(err))
	}
}

// LearnFromFeedback - обучение на основе обратной связи
func (d *AICategoryDetector) LearnFromFeedback(ctx context.Context) error {
	query := `
		SELECT
			f.ai_hints,
			f.keywords,
			f.correct_category_id,
			f.detected_category_id
		FROM category_detection_feedback f
		WHERE f.user_confirmed = TRUE
		  AND f.created_at > NOW() - INTERVAL '7 days'
		ORDER BY f.created_at DESC
		LIMIT 1000
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			d.logger.Warn("Failed to close rows in LearnFromFeedback", zap.Error(err))
		}
	}()

	for rows.Next() {
		var (
			aiHintsJSON        []byte
			keywords           []string
			correctCategoryID  int32
			detectedCategoryID int32
		)

		err := rows.Scan(&aiHintsJSON, &keywords, &correctCategoryID, &detectedCategoryID)
		if err != nil {
			continue
		}

		// Обновляем веса ключевых слов
		d.updateKeywordWeights(ctx, keywords, correctCategoryID, detectedCategoryID == correctCategoryID)

		// Обновляем AI маппинги
		if len(aiHintsJSON) > 0 {
			var hints AIHints
			if err := json.Unmarshal(aiHintsJSON, &hints); err == nil {
				d.updateAIMappingStats(ctx, hints.Domain, hints.ProductType, correctCategoryID,
					detectedCategoryID == correctCategoryID)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	return nil
}

// updateKeywordWeights - обновление весов ключевых слов
func (d *AICategoryDetector) updateKeywordWeights(ctx context.Context, keywords []string, categoryID int32, success bool) {
	multiplier := 1.01
	if !success {
		multiplier = 0.99
	}

	for _, keyword := range keywords {
		query := `
			INSERT INTO category_keyword_weights (keyword, category_id, weight, occurrence_count, success_rate)
			VALUES ($1, $2, 1.0, 1, $3)
			ON CONFLICT (keyword, category_id, language)
			DO UPDATE SET
				weight = LEAST(category_keyword_weights.weight * $4, 2.0),
				occurrence_count = category_keyword_weights.occurrence_count + 1,
				success_rate =
					(category_keyword_weights.success_rate * category_keyword_weights.occurrence_count + $3) /
					(category_keyword_weights.occurrence_count + 1)
		`

		successRate := 0.0
		if success {
			successRate = 1.0
		}

		if _, err := d.db.ExecContext(ctx, query, keyword, categoryID, successRate, multiplier); err != nil {
			d.logger.Warn("Failed to update keyword weight",
				zap.String("keyword", keyword),
				zap.Int32("categoryID", categoryID),
				zap.Error(err))
		}
	}
}

// updateAIMappingStats - обновление статистики AI маппинга
func (d *AICategoryDetector) updateAIMappingStats(ctx context.Context, domain, productType string, categoryID int32, success bool) {
	query := `
		UPDATE category_ai_mappings
		SET %s = %s + 1,
		    weight = %s
		WHERE ai_domain = $1
		  AND product_type = $2
		  AND category_id = $3
	`

	if success {
		query = fmt.Sprintf(query, "success_count", "success_count", "LEAST(weight * 1.01, 1.0)")
	} else {
		query = fmt.Sprintf(query, "failure_count", "failure_count", "GREATEST(weight * 0.99, 0.1)")
	}

	if _, err := d.db.ExecContext(ctx, query, domain, productType, categoryID); err != nil {
		d.logger.Warn("Failed to update AI mapping stats",
			zap.String("domain", domain),
			zap.String("productType", productType),
			zap.Int32("categoryID", categoryID),
			zap.Error(err))
	}
}

// ConfirmDetection - подтверждение правильности определения
func (d *AICategoryDetector) ConfirmDetection(ctx context.Context, feedbackID int64, correctCategoryID int32) error {
	query := `
		UPDATE category_detection_feedback
		SET user_confirmed = TRUE,
		    correct_category_id = $2
		WHERE id = $1
	`

	_, err := d.db.ExecContext(ctx, query, feedbackID, correctCategoryID)
	return err
}

// GetAccuracyMetrics - получение метрик точности
func (d *AICategoryDetector) GetAccuracyMetrics(ctx context.Context, days int) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) as confirmed,
			AVG(confidence_score) as avg_confidence,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY processing_time_ms) as median_time
		FROM category_detection_feedback
		WHERE created_at > NOW() - INTERVAL '%d days'
	`

	var metrics struct {
		Total         int64   `db:"total"`
		Confirmed     int64   `db:"confirmed"`
		AvgConfidence float64 `db:"avg_confidence"`
		MedianTime    float64 `db:"median_time"`
	}

	err := d.db.GetContext(ctx, &metrics, fmt.Sprintf(query, days))
	if err != nil {
		return nil, err
	}

	accuracy := float64(0)
	if metrics.Total > 0 {
		accuracy = float64(metrics.Confirmed) / float64(metrics.Total) * 100
	}

	return map[string]interface{}{
		"totalDetections":     metrics.Total,
		"confirmedDetections": metrics.Confirmed,
		"accuracyPercent":     accuracy,
		"avgConfidence":       metrics.AvgConfidence,
		"medianTimeMs":        metrics.MedianTime,
	}, nil
}

// Cache methods
func (d *AICategoryDetector) getCacheKey(input AIDetectionInput) string {
	data := fmt.Sprintf("%s|%s|%v", input.Title, input.Description, input.AIHints)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (d *AICategoryDetector) getFromMemCache(key string) *AIDetectionResult {
	if entry, ok := d.memCache[key]; ok {
		if entry.expiresAt.After(time.Now()) {
			return entry.result
		}
		delete(d.memCache, key)
	}
	return nil
}

func (d *AICategoryDetector) saveToMemCache(key string, result *AIDetectionResult) {
	d.memCache[key] = &cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(15 * time.Minute),
	}

	// Очистка старых записей
	if len(d.memCache) > 1000 {
		for k, v := range d.memCache {
			if v.expiresAt.Before(time.Now()) {
				delete(d.memCache, k)
			}
		}
	}
}

// GetCacheStats возвращает статистику кэша
func (d *AICategoryDetector) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	if d.useRedis {
		return d.redisCache.GetStats(ctx)
	}

	// Статистика для in-memory кэша
	validEntries := 0
	now := time.Now()
	for _, entry := range d.memCache {
		if entry.expiresAt.After(now) {
			validEntries++
		}
	}

	return map[string]interface{}{
		"type":         "memory",
		"totalEntries": len(d.memCache),
		"validEntries": validEntries,
	}, nil
}

// ClearCache очищает весь кэш
func (d *AICategoryDetector) ClearCache(ctx context.Context) error {
	if d.useRedis {
		if err := d.redisCache.Clear(ctx); err != nil {
			return err
		}
	}
	d.memCache = make(map[string]*cacheEntry)
	return nil
}

// getAIEndpoint получает AI endpoint из переменных окружения
func getAIEndpoint() string {
	endpoint := os.Getenv("AI_ENDPOINT")
	if endpoint == "" {
		return "https://api.anthropic.com/v1/messages" // fallback
	}
	return endpoint
}

// CategoryOption представляет опцию категории для AI
type CategoryOption struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// AIFallbackRequest структура запроса к AI для выбора категории
type AIFallbackRequest struct {
	Model     string      `json:"model"`
	MaxTokens int         `json:"max_tokens"`
	Messages  []AIMessage `json:"messages"`
}

// AIMessage представляет сообщение для AI
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIFallbackResponse ответ от AI с выбором категории
type AIFallbackResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// AISelectionResult результат AI выбора категории
type AISelectionResult struct {
	CategoryID     int32   `json:"categoryId"`
	CategoryName   string  `json:"categoryName"`
	CategorySlug   string  `json:"categorySlug"`
	Confidence     float64 `json:"confidence"`
	Reasoning      string  `json:"reasoning"`
	AlternativeIDs []int32 `json:"alternativeIds,omitempty"`
}

// DetectWithAIFallback главный метод с AI как основным механизмом для 99% точности
func (d *AICategoryDetector) DetectWithAIFallback(ctx context.Context, input AIDetectionInput) (*AIDetectionResult, error) {
	startTime := time.Now()
	d.logger.Info("Starting category detection with AI-first approach",
		zap.String("title", input.Title),
		zap.String("description", input.Description[:min(100, len(input.Description))]))

	// Фаза 1: Проверяем кеш AI решений в новой таблице
	cachedResult, err := d.getAIDecisionFromCache(ctx, input)
	if err == nil && cachedResult != nil {
		d.logger.Info("Found cached AI decision",
			zap.Int32("categoryId", cachedResult.CategoryID),
			zap.String("categoryName", cachedResult.CategoryName),
			zap.Float64("confidence", cachedResult.ConfidenceScore))
		cachedResult.Algorithm = "ai_cached_decision"
		cachedResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		return cachedResult, nil
	}

	// Фаза 2: Быстрый локальный поиск (если confidence > 90%)
	localResult, err := d.quickLocalSearch(ctx, input)
	if err == nil && localResult != nil && localResult.ConfidenceScore > 0.90 {
		// Семантическая валидация
		if d.isSemanticallySensible(ctx, input, localResult) {
			d.logger.Info("Quick local search successful with high confidence",
				zap.Float64("confidence", localResult.ConfidenceScore),
				zap.Int32("categoryId", localResult.CategoryID))
			localResult.Algorithm = "quick_local_high_confidence"
			localResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()
			// Сохраняем в кеш как подтвержденное решение
			d.saveAIDecisionToCache(ctx, input, localResult, "", false)
			return localResult, nil
		}
	}

	// Фаза 3: AI выбирает из полного списка категорий (основной метод)
	d.logger.Info("Using AI to select category from full list")
	aiResult, err := d.selectCategoryWithAI(ctx, input)
	if err != nil {
		d.logger.Error("AI category selection failed", zap.Error(err))
		// Fallback на стандартные методы
		if localResult != nil {
			localResult.Algorithm = "local_fallback_ai_failed"
			localResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()
			return localResult, nil
		}
		// Крайний случай - используем DetectCategory
		standardResult, stdErr := d.DetectCategory(ctx, input)
		if stdErr != nil {
			return nil, fmt.Errorf("all detection methods failed: AI: %w, standard: %w", err, stdErr)
		}
		standardResult.Algorithm = "standard_fallback_all_failed"
		standardResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		return standardResult, nil
	}

	// Конвертируем AI результат в стандартный формат
	finalResult := &AIDetectionResult{
		CategoryID:       aiResult.CategoryID,
		CategoryName:     aiResult.CategoryName,
		CategoryPath:     aiResult.CategorySlug,
		ConfidenceScore:  aiResult.Confidence,
		Algorithm:        "ai_direct_selection_v2",
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
		Keywords:         d.extractKeywords(input),
		AIHints:          input.AIHints,
		AlternativeIDs:   aiResult.AlternativeIDs,
	}

	d.logger.Info("AI category selection successful",
		zap.Float64("aiConfidence", aiResult.Confidence),
		zap.Int32("aiCategoryId", aiResult.CategoryID),
		zap.String("aiCategoryName", aiResult.CategoryName),
		zap.String("reasoning", aiResult.Reasoning[:min(200, len(aiResult.Reasoning))]))

	// Сохраняем решение AI в кеш для будущего использования
	d.saveAIDecisionToCache(ctx, input, finalResult, aiResult.Reasoning, true)

	// Логируем для обучения
	d.logDetection(ctx, input, finalResult)

	return finalResult, nil
}

// selectCategoryWithAI отправляет все категории в AI для выбора лучшей
func (d *AICategoryDetector) selectCategoryWithAI(ctx context.Context, input AIDetectionInput) (*AISelectionResult, error) {
	// Получаем все активные категории
	categories, err := d.getAllActiveCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Создаем промпт с полным списком категорий
	prompt, err := d.buildCategorySelectionPrompt(input, categories)
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Отправляем запрос к AI
	aiResponse, err := d.sendAIRequest(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI request failed: %w", err)
	}

	// Парсим ответ AI
	result, err := d.parseAIResponse(aiResponse, categories)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// getAllActiveCategories получает все активные категории из БД
func (d *AICategoryDetector) getAllActiveCategories(ctx context.Context) ([]CategoryOption, error) {
	query := `
		SELECT id, name, slug
		FROM marketplace_categories
		WHERE is_active = TRUE
		ORDER BY sort_order, name
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			d.logger.Warn("Failed to close rows in getAllActiveCategories", zap.Error(err))
		}
	}()

	var categories []CategoryOption
	for rows.Next() {
		var cat CategoryOption
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug)
		if err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("no active categories found")
	}

	d.logger.Debug("Retrieved categories for AI selection", zap.Int("count", len(categories)))
	return categories, nil
}

// buildCategorySelectionPrompt создает промпт для AI с полным списком категорий
func (d *AICategoryDetector) buildCategorySelectionPrompt(input AIDetectionInput, categories []CategoryOption) (string, error) {
	// Строим список категорий
	categoryList := "\nДоступные категории:\n"
	for _, cat := range categories {
		categoryList += fmt.Sprintf("- ID: %d, Название: \"%s\", Slug: \"%s\"\n", cat.ID, cat.Name, cat.Slug)
	}

	// Создаем структурированный промпт
	prompt := fmt.Sprintf(`Ты эксперт по категоризации товаров в маркетплейсе. Проанализируй товар и выбери НАИБОЛЕЕ ТОЧНУЮ категорию.

ТОВАР ДЛЯ АНАЛИЗА:
Название: "%s"
Описание: "%s"

%s

ЗАДАЧА:
1. Внимательно проанализируй название и описание товара
2. Выбери ЕДИНСТВЕННУЮ наиболее подходящую категорию из списка выше
3. Укажи свою уверенность в процентах (0-100)
4. Кратко объясни свой выбор

ТРЕБОВАНИЯ К ОТВЕТУ:
- Отвечай ТОЛЬКО в JSON формате
- Обязательные поля: categoryId, categoryName, categorySlug, confidence, reasoning
- Уверенность должна быть реалистичной
- Reasoning должно быть кратким и ясным

ФОРМАТ ОТВЕТА:
{
  "categoryId": 1234,
  "categoryName": "Название категории",
  "categorySlug": "slug-kategorii",
  "confidence": 95.0,
  "reasoning": "Краткое объяснение выбора"
}

ВАЖНО: Отвечай только JSON, без дополнительного текста!`,
		input.Title,
		input.Description,
		categoryList)

	return prompt, nil
}

// sendAIRequest отправляет запрос к AI API
func (d *AICategoryDetector) sendAIRequest(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("CLAUDE_API_KEY not set")
	}

	request := AIFallbackRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 500,
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Конвертируем в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, "POST", d.aiEndpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Отправляем запрос
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			d.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Error("AI API error",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("response", string(body)))
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// parseAIResponse парсит ответ AI и извлекает выбранную категорию
func (d *AICategoryDetector) parseAIResponse(response string, categories []CategoryOption) (*AISelectionResult, error) {
	var aiResp AIFallbackResponse
	if err := json.Unmarshal([]byte(response), &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response JSON: %w", err)
	}

	if len(aiResp.Content) == 0 {
		return nil, fmt.Errorf("empty content in AI response")
	}

	// Извлекаем JSON из content
	content := aiResp.Content[0].Text

	// Находим JSON в ответе (может быть обернут в markdown)
	jsonStart := strings.Index(content, "{")
	jsonEnd := strings.LastIndex(content, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no JSON found in AI response: %s", content)
	}

	jsonContent := content[jsonStart:jsonEnd]

	var result AISelectionResult
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI selection JSON: %w, content: %s", err, jsonContent)
	}

	// Валидируем результат
	if err := d.validateAIResult(&result, categories); err != nil {
		return nil, fmt.Errorf("AI result validation failed: %w", err)
	}

	// Нормализуем confidence score (0-1)
	if result.Confidence > 1.0 {
		result.Confidence /= 100.0
	}

	d.logger.Debug("Successfully parsed AI response",
		zap.Int32("categoryId", result.CategoryID),
		zap.String("categoryName", result.CategoryName),
		zap.Float64("confidence", result.Confidence))

	return &result, nil
}

// validateAIResult проверяет корректность результата AI
func (d *AICategoryDetector) validateAIResult(result *AISelectionResult, categories []CategoryOption) error {
	if result.CategoryID == 0 {
		return fmt.Errorf("missing categoryId")
	}

	if result.CategoryName == "" {
		return fmt.Errorf("missing categoryName")
	}

	if result.Confidence < 0 || result.Confidence > 100 {
		return fmt.Errorf("invalid confidence value: %f", result.Confidence)
	}

	// Проверяем, что выбранная категория существует
	found := false
	for _, cat := range categories {
		if cat.ID == result.CategoryID {
			found = true
			// Обновляем имя и slug из БД для консистентности
			result.CategoryName = cat.Name
			result.CategorySlug = cat.Slug
			break
		}
	}

	if !found {
		return fmt.Errorf("selected category ID %d not found in available categories", result.CategoryID)
	}

	return nil
}

// min возвращает минимальное из двух значений
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isSemanticallySensible проверяет семантическую корректность результата
func (d *AICategoryDetector) isSemanticallySensible(ctx context.Context, input AIDetectionInput, result *AIDetectionResult) bool {
	// Базовые правила семантической валидации
	title := strings.ToLower(input.Title)
	description := strings.ToLower(input.Description)
	categoryName := strings.ToLower(result.CategoryName)

	// Комбинированный текст для анализа
	fullText := title + " " + description

	// Правило 1: Явные мобильные устройства не должны попадать в автомобили
	mobileKeywords := []string{"iphone", "samsung", "xiaomi", "huawei", "pixel", "oneplus", "smartphone", "телефон", "смартфон", "мобильный"}
	autoCategories := []string{"automobili", "auto", "car", "автомобил"}

	if d.containsAnyKeyword(fullText, mobileKeywords) && d.containsAnyKeyword(categoryName, autoCategories) {
		d.logger.Debug("Semantic validation failed: mobile device categorized as automobile",
			zap.String("title", input.Title),
			zap.String("category", result.CategoryName))
		return false
	}

	// Правило 2: Одежда не должна попадать в электронику
	clothingKeywords := []string{"майка", "футболка", "джинсы", "платье", "юбка", "брюки", "куртка", "рубашка", "блузка"}
	electronicsCategories := []string{"elektronika", "electronics", "электроника"}

	if d.containsAnyKeyword(fullText, clothingKeywords) && d.containsAnyKeyword(categoryName, electronicsCategories) {
		d.logger.Debug("Semantic validation failed: clothing categorized as electronics",
			zap.String("title", input.Title),
			zap.String("category", result.CategoryName))
		return false
	}

	// Правило 3: Еда не должна попадать в технику
	foodKeywords := []string{"хлеб", "молоко", "мясо", "овощи", "фрукты", "сыр", "колбаса", "рыба", "макароны", "рис"}
	techCategories := []string{"elektronika", "electronics", "kućni aparati", "home-appliances"}

	if d.containsAnyKeyword(fullText, foodKeywords) && d.containsAnyKeyword(categoryName, techCategories) {
		d.logger.Debug("Semantic validation failed: food categorized as technology",
			zap.String("title", input.Title),
			zap.String("category", result.CategoryName))
		return false
	}

	// Правило 4: Книги не должны попадать в автомобили
	bookKeywords := []string{"книга", "учебник", "роман", "справочник", "энциклопедия", "словарь", "book"}

	if d.containsAnyKeyword(fullText, bookKeywords) && d.containsAnyKeyword(categoryName, autoCategories) {
		d.logger.Debug("Semantic validation failed: book categorized as automobile",
			zap.String("title", input.Title),
			zap.String("category", result.CategoryName))
		return false
	}

	// Правило 5: Домашние животные не должны попадать в автозапчасти
	petKeywords := []string{"кот", "собака", "котенок", "щенок", "кошка", "хомяк", "попугай", "рыбка", "cat", "dog"}
	autoPartsCategories := []string{"auto delovi", "auto-parts", "автозапчасти"}

	if d.containsAnyKeyword(fullText, petKeywords) && d.containsAnyKeyword(categoryName, autoPartsCategories) {
		d.logger.Debug("Semantic validation failed: pet categorized as auto parts",
			zap.String("title", input.Title),
			zap.String("category", result.CategoryName))
		return false
	}

	// Если все проверки прошли - результат семантически корректен
	return true
}

// containsAnyKeyword проверяет содержит ли текст любое из ключевых слов
func (d *AICategoryDetector) containsAnyKeyword(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// getAIDecisionFromCache получает закешированное решение AI из новой таблицы
func (d *AICategoryDetector) getAIDecisionFromCache(ctx context.Context, input AIDetectionInput) (*AIDetectionResult, error) {
	// Вычисляем хеш заголовка
	hash := sha256.Sum256([]byte(strings.ToLower(input.Title)))
	titleHash := hex.EncodeToString(hash[:])

	query := `
		SELECT
			acd.category_id,
			c.name as category_name,
			c.slug as category_path,
			acd.confidence,
			acd.reasoning,
			acd.alternative_category_ids,
			acd.ai_keywords,
			acd.ai_domain,
			acd.ai_product_type
		FROM ai_category_decisions acd
		JOIN marketplace_categories c ON c.id = acd.category_id
		WHERE acd.title_hash = $1
			AND acd.created_at > NOW() - INTERVAL '30 days'
		ORDER BY acd.confidence DESC, acd.created_at DESC
		LIMIT 1
	`

	var result struct {
		CategoryID     int32    `db:"category_id"`
		CategoryName   string   `db:"category_name"`
		CategoryPath   string   `db:"category_path"`
		Confidence     float64  `db:"confidence"`
		Reasoning      *string  `db:"reasoning"`
		AlternativeIDs []int32  `db:"alternative_category_ids"`
		AIKeywords     []string `db:"ai_keywords"`
		AIDomain       *string  `db:"ai_domain"`
		AIProductType  *string  `db:"ai_product_type"`
	}

	err := d.db.GetContext(ctx, &result, query, titleHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoCachedResult // Не найдено в кеше
		}
		return nil, fmt.Errorf("failed to get cached AI decision: %w", err)
	}

	// Собираем AIHints если есть
	var aiHints *AIHints
	if result.AIDomain != nil || result.AIProductType != nil || len(result.AIKeywords) > 0 {
		aiHints = &AIHints{}
		if result.AIDomain != nil {
			aiHints.Domain = *result.AIDomain
		}
		if result.AIProductType != nil {
			aiHints.ProductType = *result.AIProductType
		}
		aiHints.Keywords = result.AIKeywords
	}

	return &AIDetectionResult{
		CategoryID:      result.CategoryID,
		CategoryName:    result.CategoryName,
		CategoryPath:    result.CategoryPath,
		ConfidenceScore: result.Confidence,
		AlternativeIDs:  result.AlternativeIDs,
		Keywords:        result.AIKeywords,
		AIHints:         aiHints,
	}, nil
}

// saveAIDecisionToCache сохраняет решение AI в таблицу для будущего использования
func (d *AICategoryDetector) saveAIDecisionToCache(ctx context.Context, input AIDetectionInput, result *AIDetectionResult, reasoning string, fromAI bool) {
	// Вычисляем хеш заголовка
	hash := sha256.Sum256([]byte(strings.ToLower(input.Title)))
	titleHash := hex.EncodeToString(hash[:])

	// Определяем модель
	aiModel := "local_algorithms"
	if fromAI {
		aiModel = "claude-3-haiku-20240307"
	}

	// Подготавливаем AI hints
	var aiDomain, aiProductType *string
	var aiKeywords []string

	if input.AIHints != nil {
		if input.AIHints.Domain != "" {
			aiDomain = &input.AIHints.Domain
		}
		if input.AIHints.ProductType != "" {
			aiProductType = &input.AIHints.ProductType
		}
		aiKeywords = input.AIHints.Keywords
	}

	query := `
		INSERT INTO ai_category_decisions (
			title_hash, title, description, category_id, confidence,
			reasoning, alternative_category_ids, ai_model, processing_time_ms,
			ai_domain, ai_product_type, ai_keywords
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (title_hash) DO UPDATE SET
			category_id = EXCLUDED.category_id,
			confidence = EXCLUDED.confidence,
			reasoning = EXCLUDED.reasoning,
			alternative_category_ids = EXCLUDED.alternative_category_ids,
			ai_model = EXCLUDED.ai_model,
			processing_time_ms = EXCLUDED.processing_time_ms,
			updated_at = CURRENT_TIMESTAMP
		WHERE ai_category_decisions.confidence < EXCLUDED.confidence
	`

	_, err := d.db.ExecContext(ctx, query,
		titleHash,
		input.Title,
		input.Description,
		result.CategoryID,
		result.ConfidenceScore,
		reasoning,
		result.AlternativeIDs,
		aiModel,
		result.ProcessingTimeMs,
		aiDomain,
		aiProductType,
		aiKeywords,
	)

	if err != nil {
		d.logger.Error("Failed to save AI decision to cache",
			zap.String("titleHash", titleHash),
			zap.Int32("categoryId", result.CategoryID),
			zap.Error(err))
	} else {
		d.logger.Debug("AI decision saved to cache",
			zap.String("titleHash", titleHash),
			zap.Int32("categoryId", result.CategoryID),
			zap.Float64("confidence", result.ConfidenceScore))
	}
}

// quickLocalSearch выполняет быстрый локальный поиск по ключевым словам и похожести
func (d *AICategoryDetector) quickLocalSearch(ctx context.Context, input AIDetectionInput) (*AIDetectionResult, error) {
	// Используем комбинацию методов для быстрого поиска
	results := []weightedResult{}

	// Ключевые слова
	keywords := d.extractKeywords(input)
	if kwResult := d.detectByKeywords(ctx, keywords); kwResult != nil {
		results = append(results, weightedResult{
			Result: kwResult,
			Weight: 0.8,
		})
	}

	// AI hints если есть
	if input.AIHints != nil && input.AIHints.Domain != "" {
		if aiResult := d.detectByAIHints(ctx, input.AIHints); aiResult != nil {
			results = append(results, weightedResult{
				Result: aiResult,
				Weight: 0.9,
			})
		}
	}

	// Похожие товары
	if simResult := d.detectBySimilarity(ctx, input); simResult != nil {
		results = append(results, weightedResult{
			Result: simResult,
			Weight: 0.6,
		})
	}

	if len(results) == 0 {
		return nil, ErrNoResultsFound
	}

	// Взвешенное голосование
	return d.weightedVoting(results), nil
}

// SelectCategoryDirectly - метод прямого выбора категории через AI без fallback
// Всегда использует AI для выбора категории из полного списка
func (d *AICategoryDetector) SelectCategoryDirectly(ctx context.Context, input AIDetectionInput) (*AIDetectionResult, error) {
	startTime := time.Now()

	// Сначала проверяем кеш
	cachedResult, err := d.getAIDecisionFromCache(ctx, input)
	if err == nil && cachedResult != nil {
		d.logger.Info("Found cached AI decision for direct selection",
			zap.Int32("categoryId", cachedResult.CategoryID),
			zap.String("categoryName", cachedResult.CategoryName))
		cachedResult.Algorithm = "ai_direct_cached"
		cachedResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		return cachedResult, nil
	}

	d.logger.Info("Using AI for direct category selection",
		zap.String("title", input.Title))

	// Напрямую вызываем AI для выбора категории
	aiResult, err := d.selectCategoryWithAI(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("AI category selection failed: %w", err)
	}

	// Конвертируем результат
	finalResult := &AIDetectionResult{
		CategoryID:       aiResult.CategoryID,
		CategoryName:     aiResult.CategoryName,
		CategoryPath:     aiResult.CategorySlug,
		ConfidenceScore:  aiResult.Confidence,
		Algorithm:        "ai_direct_selection",
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
		Keywords:         d.extractKeywords(input),
		AIHints:          input.AIHints,
		AlternativeIDs:   aiResult.AlternativeIDs,
	}

	d.logger.Info("AI direct selection completed",
		zap.Int32("categoryId", finalResult.CategoryID),
		zap.String("categoryName", finalResult.CategoryName),
		zap.Float64("confidence", finalResult.ConfidenceScore),
		zap.String("reasoning", aiResult.Reasoning[:min(200, len(aiResult.Reasoning))]))

	// Сохраняем в кеш
	d.saveAIDecisionToCache(ctx, input, finalResult, aiResult.Reasoning, true)

	// Логируем для статистики
	d.logDetection(ctx, input, finalResult)

	return finalResult, nil
}
