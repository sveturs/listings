package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

const claudeAPIURL = "https://api.anthropic.com/v1/messages"

// CategoryDetectionService - сервис детекции категорий
type CategoryDetectionService struct {
	repo           *postgres.CategoryDetectionRepository
	categoryRepo   CategoryRepositoryInterface
	redisClient    *redis.Client
	claudeAPIKey   string
	logger         zerolog.Logger
	categoryCache  []*domain.CategoryV2 // Cached categories for Claude prompt
	cacheLoadedAt  time.Time            // When cache was last loaded
	cacheTTL       time.Duration        // Cache TTL (default 1 hour)
}

// CategoryRepositoryInterface - интерфейс для получения категорий
type CategoryRepositoryInterface interface {
	GetByUUID(ctx context.Context, id string) (*domain.CategoryV2, error)
	GetBySlugV2(ctx context.Context, slug string) (*domain.CategoryV2, error)
	GetAllActiveV2(ctx context.Context) ([]*domain.CategoryV2, error)
}

// NewCategoryDetectionService создаёт сервис
func NewCategoryDetectionService(
	repo *postgres.CategoryDetectionRepository,
	categoryRepo CategoryRepositoryInterface,
	redisClient *redis.Client,
	claudeAPIKey string,
	logger zerolog.Logger,
) *CategoryDetectionService {
	svc := &CategoryDetectionService{
		repo:         repo,
		categoryRepo: categoryRepo,
		redisClient:  redisClient,
		claudeAPIKey: claudeAPIKey,
		logger:       logger.With().Str("component", "category_detection_service").Logger(),
		cacheTTL:     1 * time.Hour, // Categories don't change often
	}

	// Load categories into cache on startup (non-blocking)
	go svc.refreshCategoryCache()

	return svc
}

// refreshCategoryCache loads all active categories into memory cache
func (s *CategoryDetectionService) refreshCategoryCache() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	categories, err := s.categoryRepo.GetAllActiveV2(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to load categories for cache")
		return
	}

	s.categoryCache = categories
	s.cacheLoadedAt = time.Now()
	s.logger.Info().Int("count", len(categories)).Msg("Category cache loaded successfully")
}

// getCachedCategories returns cached categories, refreshing if needed
func (s *CategoryDetectionService) getCachedCategories() []*domain.CategoryV2 {
	// Check if cache is stale
	if s.categoryCache == nil || time.Since(s.cacheLoadedAt) > s.cacheTTL {
		s.refreshCategoryCache()
	}
	return s.categoryCache
}

// DetectFromText определяет категорию по тексту с каскадной детекцией
func (s *CategoryDetectionService) DetectFromText(
	ctx context.Context,
	input domain.DetectFromTextInput,
) (*domain.CategoryDetection, error) {
	startTime := time.Now()
	fullText := input.Title + " " + input.Description

	// === ЭТАП 0: AI Mapping (НАИВЫСШИЙ ПРИОРИТЕТ, confidence 0.95-1.0) ===
	// Если Claude AI предложил категорию при анализе изображения, используем её ПЕРВОЙ
	if input.SuggestedCategory != "" {
		aiMapping, err := s.repo.FindByAIMapping(ctx, input.SuggestedCategory, input.Language)
		if err != nil {
			s.logger.Warn().Err(err).
				Str("suggested_category", input.SuggestedCategory).
				Msg("AI mapping query failed")
		}

		if aiMapping != nil && aiMapping.ConfidenceScore >= 0.90 {
			s.logger.Info().
				Str("suggested_category", input.SuggestedCategory).
				Str("category", aiMapping.CategorySlug).
				Float64("confidence", aiMapping.ConfidenceScore).
				Str("method", "ai_mapping").
				Msg("Using AI mapping (highest priority)")
			return s.buildDetection([]domain.CategoryMatch{*aiMapping}, startTime, input)
		}

		// Если маппинг нашёлся, но confidence низкий, логируем для анализа
		if aiMapping != nil {
			s.logger.Warn().
				Str("suggested_category", input.SuggestedCategory).
				Float64("confidence", aiMapping.ConfidenceScore).
				Msg("AI mapping found but confidence too low (<0.90)")
		}
	}

	// === ЭТАП 1: Brand Matching (самый точный, confidence 0.95-0.98) ===
	brandMatches, err := s.detectByBrand(ctx, fullText, input.Language)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Brand matching failed")
	}

	if len(brandMatches) > 0 && brandMatches[0].ConfidenceScore >= 0.95 {
		s.logger.Info().
			Float64("confidence", brandMatches[0].ConfidenceScore).
			Str("category", brandMatches[0].CategorySlug).
			Str("method", string(domain.MethodBrandMatch)).
			Msg("Using brand match result")
		return s.buildDetection(brandMatches, startTime, input)
	}

	// === ЭТАП 1.5: ProductType Hint Matching (для часов, телефонов и т.д.) ===
	if input.Hints != nil && input.Hints.ProductType != "" {
		productTypeMatch := s.detectByProductType(ctx, input.Hints.ProductType, input.Language)
		if productTypeMatch != nil && productTypeMatch.ConfidenceScore >= 0.85 {
			s.logger.Info().
				Str("productType", input.Hints.ProductType).
				Str("category", productTypeMatch.CategorySlug).
				Float64("confidence", productTypeMatch.ConfidenceScore).
				Msg("Using productType hint match")
			return s.buildDetection([]domain.CategoryMatch{*productTypeMatch}, startTime, input)
		}
	}

	// === ЭТАП 2: Keyword Matching (быстро и бесплатно) ===
	keywords := extractKeywords(fullText)
	keywordMatches, err := s.repo.FindByKeywords(ctx, keywords, input.Language)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Keyword matching failed")
	}

	// Если есть хороший keyword match (>0.6), используем его
	if len(keywordMatches) > 0 && keywordMatches[0].ConfidenceScore > 0.6 {
		s.logger.Debug().
			Float64("confidence", keywordMatches[0].ConfidenceScore).
			Str("category", keywordMatches[0].CategorySlug).
			Str("method", string(domain.MethodKeywordMatch)).
			Msg("Using keyword match result")

		// Добавить brand matches как альтернативы
		allMatches := mergeMatches(keywordMatches, brandMatches)
		return s.buildDetection(allMatches, startTime, input)
	}

	// === ЭТАП 3: Similarity Matching ===
	similarityMatches, err := s.repo.FindBySimilarity(ctx, input.Title, input.Language)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Similarity matching failed")
	}

	if len(similarityMatches) > 0 && similarityMatches[0].ConfidenceScore > 0.4 {
		s.logger.Debug().
			Float64("confidence", similarityMatches[0].ConfidenceScore).
			Str("category", similarityMatches[0].CategorySlug).
			Str("method", string(domain.MethodSimilarity)).
			Msg("Using similarity match result")

		// Комбинировать все результаты
		allMatches := mergeMatches(similarityMatches, keywordMatches, brandMatches)
		return s.buildDetection(allMatches, startTime, input)
	}

	// === ЭТАП 4: Claude AI (только как fallback, дорого) ===
	if len(keywordMatches) == 0 || keywordMatches[0].ConfidenceScore < 0.4 {
		s.logger.Info().Msg("Using Claude AI as fallback")
		aiMatches, err := s.detectWithClaude(ctx, input)
		if err != nil {
			s.logger.Warn().Err(err).Msg("Claude AI detection failed, using best available")
			// Fallback на лучшие доступные результаты
			allMatches := mergeMatches(keywordMatches, similarityMatches, brandMatches)
			if len(allMatches) > 0 {
				return s.buildDetection(allMatches, startTime, input)
			}

			// ВАЖНО: SuggestedCategory теперь обрабатывается в ЭТАП 0 (highest priority)
			// Если мы дошли сюда - значит маппинг не найден или confidence < 0.90
			return s.fallbackDetection(startTime, input)
		}

		// AI успешен - комбинируем со всеми результатами как альтернативы
		allMatches := mergeMatches(aiMatches, keywordMatches, similarityMatches, brandMatches)
		return s.buildDetection(allMatches, startTime, input)
	}

	// === Используем лучшие доступные результаты ===
	allMatches := mergeMatches(keywordMatches, similarityMatches, brandMatches)
	return s.buildDetection(allMatches, startTime, input)
}

// DetectFromKeywords определяет категорию по ключевым словам
func (s *CategoryDetectionService) DetectFromKeywords(
	ctx context.Context,
	keywords []string,
	language string,
) (*domain.CategoryDetection, error) {
	startTime := time.Now()

	matches, err := s.repo.FindByKeywords(ctx, keywords, language)
	if err != nil {
		return nil, fmt.Errorf("FindByKeywords failed: %w", err)
	}

	input := domain.DetectFromTextInput{
		Title:    strings.Join(keywords, " "),
		Language: language,
	}

	if len(matches) == 0 {
		return s.fallbackDetection(startTime, input)
	}

	return s.buildDetection(matches, startTime, input)
}

// DetectBatch выполняет batch детекцию
func (s *CategoryDetectionService) DetectBatch(
	ctx context.Context,
	input domain.DetectBatchInput,
) (*domain.DetectBatchResult, error) {
	startTime := time.Now()

	results := make([]domain.CategoryDetection, len(input.Items))

	for i, item := range input.Items {
		detection, err := s.DetectFromText(ctx, item)
		if err != nil {
			s.logger.Warn().Err(err).Int("index", i).Msg("Batch item detection failed")
			detection = s.fallbackDetectionSync(item)
		}
		results[i] = *detection
	}

	return &domain.DetectBatchResult{
		Results:             results,
		TotalProcessingTime: int32(time.Since(startTime).Milliseconds()),
	}, nil
}

// ConfirmSelection подтверждает выбор пользователя
func (s *CategoryDetectionService) ConfirmSelection(
	ctx context.Context,
	detectionID, selectedCategoryID uuid.UUID,
) error {
	return s.repo.ConfirmSelection(ctx, detectionID, selectedCategoryID)
}

// detectByBrand определяет категорию по бренду
func (s *CategoryDetectionService) detectByBrand(
	ctx context.Context,
	text string,
	language string,
) ([]domain.CategoryMatch, error) {
	// Найти бренды в тексте через репозиторий (использует кеш)
	brandMappings, err := s.repo.FindBrandInText(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("FindBrandInText failed: %w", err)
	}

	if len(brandMappings) == 0 {
		return []domain.CategoryMatch{}, nil
	}

	var matches []domain.CategoryMatch

	// Конвертировать brand mappings в CategoryMatch
	for _, mapping := range brandMappings {
		// Получить полную информацию о категории
		category, err := s.categoryRepo.GetBySlugV2(ctx, mapping.CategorySlug)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("slug", mapping.CategorySlug).
				Str("brand", mapping.BrandName).
				Msg("Failed to get category for brand")
			continue
		}

		// Локализовать категорию
		localized := category.Localize(language)

		matches = append(matches, domain.CategoryMatch{
			CategoryID:      category.ID,
			CategoryName:    localized.Name,
			CategorySlug:    category.Slug,
			CategoryPath:    category.Path,
			ConfidenceScore: mapping.Confidence,
			DetectionMethod: domain.MethodBrandMatch,
			MatchedKeywords: []string{mapping.BrandName},
		})
	}

	return matches, nil
}

// detectWithClaude использует Claude AI для детекции
func (s *CategoryDetectionService) detectWithClaude(
	ctx context.Context,
	input domain.DetectFromTextInput,
) ([]domain.CategoryMatch, error) {
	if s.claudeAPIKey == "" {
		return nil, fmt.Errorf("Claude API key not configured")
	}

	prompt := s.buildCategoryDetectionPrompt(input)

	requestBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307", // Haiku - fast and cheap for category detection
		"max_tokens": 1024,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.claudeAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	return s.parseCategoryResponse(claudeResp.Content[0].Text, input.Language)
}

// buildCategoryDetectionPrompt строит промпт для Claude с динамическим списком категорий
func (s *CategoryDetectionService) buildCategoryDetectionPrompt(input domain.DetectFromTextInput) string {
	// Get cached categories
	categories := s.getCachedCategories()

	// Build category list for prompt
	var categoryLines strings.Builder
	for _, cat := range categories {
		// Get localized name (prefer English, fallback to Serbian)
		name := cat.Name["en"]
		if name == "" {
			name = cat.Name["sr"]
		}
		if name == "" {
			name = cat.Slug
		}

		// Get description for better context
		desc := cat.Description["en"]
		if desc == "" {
			desc = cat.Description["sr"]
		}

		// Get keywords for better matching
		keywords := cat.MetaKeywords["en"]
		if keywords == "" {
			keywords = cat.MetaKeywords["sr"]
		}

		// Format: UUID | slug | Name | Keywords
		if desc != "" || keywords != "" {
			extra := ""
			if keywords != "" {
				extra = fmt.Sprintf(" (%s)", keywords)
			} else if desc != "" {
				// Truncate description if too long
				if len(desc) > 50 {
					desc = desc[:50] + "..."
				}
				extra = fmt.Sprintf(" - %s", desc)
			}
			categoryLines.WriteString(fmt.Sprintf("- %s | %s | %s%s\n", cat.ID.String(), cat.Slug, name, extra))
		} else {
			categoryLines.WriteString(fmt.Sprintf("- %s | %s | %s\n", cat.ID.String(), cat.Slug, name))
		}
	}

	return fmt.Sprintf(`Analyze this product and determine the best matching category.

Title: %s
Description: %s

AVAILABLE CATEGORIES (format: UUID | slug | Name | Keywords):
%s
IMPORTANT: You MUST return the category_id as the exact UUID from the list above!

Return a JSON response with the following structure:
{
  "category_id": "exact-uuid-from-list",
  "category_slug": "category-slug",
  "category_name": "Category Name",
  "confidence": 0.95,
  "keywords": ["keyword1", "keyword2"],
  "alternatives": [
    {"category_id": "uuid", "category_slug": "slug", "category_name": "Name", "confidence": 0.7}
  ]
}

Return ONLY valid JSON, no explanations.`, input.Title, input.Description, categoryLines.String())
}

// parseCategoryResponse парсит ответ Claude и обогащает результаты из БД
// Поддерживает оба формата: новый с category_id (UUID) и старый с category (slug)
func (s *CategoryDetectionService) parseCategoryResponse(text, language string) ([]domain.CategoryMatch, error) {
	// Очищаем от markdown code blocks
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var response struct {
		// New format with UUID
		CategoryID   string `json:"category_id"`
		CategorySlug string `json:"category_slug"`
		// Legacy format (fallback)
		Category     string   `json:"category"`
		CategoryName string   `json:"category_name"`
		Confidence   float64  `json:"confidence"`
		Keywords     []string `json:"keywords"`
		Alternatives []struct {
			CategoryID   string  `json:"category_id"`
			CategorySlug string  `json:"category_slug"`
			Category     string  `json:"category"`
			CategoryName string  `json:"category_name"`
			Confidence   float64 `json:"confidence"`
		} `json:"alternatives"`
	}

	if err := json.Unmarshal([]byte(text), &response); err != nil {
		return nil, fmt.Errorf("parse Claude response failed: %w", err)
	}

	var matches []domain.CategoryMatch

	// === Parse primary category ===
	// Try UUID first (new format), fallback to slug (legacy)
	categoryID := response.CategoryID
	categorySlug := response.CategorySlug
	if categorySlug == "" {
		categorySlug = response.Category // Legacy field
	}

	primaryMatch := s.enrichCategoryFromResponse(categoryID, categorySlug, response.CategoryName, response.Confidence, response.Keywords, language)
	if primaryMatch != nil {
		matches = append(matches, *primaryMatch)
	}

	// === Parse alternatives ===
	for _, alt := range response.Alternatives {
		altID := alt.CategoryID
		altSlug := alt.CategorySlug
		if altSlug == "" {
			altSlug = alt.Category
		}

		altMatch := s.enrichCategoryFromResponse(altID, altSlug, alt.CategoryName, alt.Confidence, nil, language)
		if altMatch != nil {
			matches = append(matches, *altMatch)
		}
	}

	return matches, nil
}

// enrichCategoryFromResponse обогащает category match данными из БД
// Пытается найти категорию по UUID (если предоставлен), затем по slug
func (s *CategoryDetectionService) enrichCategoryFromResponse(categoryID, slug, name string, confidence float64, keywords []string, language string) *domain.CategoryMatch {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var category *domain.CategoryV2
	var err error

	// Try to find by UUID first (most reliable)
	if categoryID != "" && categoryID != "uuid" && categoryID != "exact-uuid-from-list" {
		// Validate UUID format
		if _, parseErr := uuid.Parse(categoryID); parseErr == nil {
			category, err = s.categoryRepo.GetByUUID(ctx, categoryID)
			if err == nil && category != nil {
				s.logger.Debug().
					Str("category_id", categoryID).
					Str("slug", category.Slug).
					Msg("Category found by UUID from Claude")
			}
		}
	}

	// Fallback to slug-based search if UUID lookup failed
	if category == nil && slug != "" {
		category = s.findCategoryBySlugVariants(ctx, slug)
	}

	if category == nil {
		s.logger.Warn().
			Str("category_id", categoryID).
			Str("slug", slug).
			Str("name", name).
			Msg("Category from Claude not found in DB, using generic match")

		// Возвращаем match без UUID (zero UUID)
		return &domain.CategoryMatch{
			CategorySlug:    slug,
			CategoryName:    name,
			ConfidenceScore: confidence,
			DetectionMethod: domain.MethodAIClaude,
			MatchedKeywords: keywords,
		}
	}

	// Локализуем название
	localized := category.Localize(language)

	return &domain.CategoryMatch{
		CategoryID:      category.ID,
		CategorySlug:    category.Slug,
		CategoryName:    localized.Name,
		CategoryPath:    category.Path,
		ConfidenceScore: confidence,
		DetectionMethod: domain.MethodAIClaude,
		MatchedKeywords: keywords,
	}
}

// findCategoryBySlugVariants tries to find a category by various slug transformations
func (s *CategoryDetectionService) findCategoryBySlugVariants(ctx context.Context, slug string) *domain.CategoryV2 {
	slugVariants := []string{
		slug,                               // as-is
		strings.ReplaceAll(slug, "-", "_"), // home-garden -> home_garden
		strings.ToLower(slug),              // Electronics -> electronics
		s.mapAISlugToDBSlug(slug),          // electronics -> elektronika
	}

	for _, trySlug := range slugVariants {
		if trySlug == "" {
			continue
		}
		category, err := s.categoryRepo.GetBySlugV2(ctx, trySlug)
		if err == nil && category != nil {
			return category
		}
	}

	return nil
}

// enrichCategoryFromDB обогащает category match данными из БД
// Пытается найти категорию по slug и получить реальный UUID
func (s *CategoryDetectionService) enrichCategoryFromDB(slug, name string, confidence float64, keywords []string, language string) *domain.CategoryMatch {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Нормализуем slug для поиска в БД
	// Claude может вернуть "electronics" или "home-garden"
	slugVariants := []string{
		slug,                                    // как есть
		strings.ReplaceAll(slug, "-", "_"),      // electronics -> electronics
		strings.ToLower(slug),                   // Electronics -> electronics
		s.mapAISlugToDBSlug(slug),               // electronics -> elektronika
	}

	var category *domain.CategoryV2
	var err error

	// Пробуем найти по разным вариантам slug
	for _, trySlug := range slugVariants {
		if trySlug == "" {
			continue
		}
		category, err = s.categoryRepo.GetBySlugV2(ctx, trySlug)
		if err == nil && category != nil {
			break
		}
	}

	if category == nil {
		s.logger.Warn().
			Str("slug", slug).
			Str("name", name).
			Msg("Category from Claude not found in DB, using generic match")

		// Возвращаем match без UUID (zero UUID)
		return &domain.CategoryMatch{
			CategorySlug:    slug,
			CategoryName:    name,
			ConfidenceScore: confidence,
			DetectionMethod: domain.MethodAIClaude,
			MatchedKeywords: keywords,
		}
	}

	// Локализуем название
	localized := category.Localize(language)

	return &domain.CategoryMatch{
		CategoryID:      category.ID,
		CategorySlug:    category.Slug,
		CategoryName:    localized.Name,
		CategoryPath:    category.Path,
		ConfidenceScore: confidence,
		DetectionMethod: domain.MethodAIClaude,
		MatchedKeywords: keywords,
	}
}

// resolveSuggestedCategory пытается найти категорию по suggested slug из AI анализа изображения
// Это последний fallback когда все остальные методы не сработали
//nolint:unused // Reserved for future use
func (s *CategoryDetectionService) resolveSuggestedCategory(ctx context.Context, suggestedCategory, language string) *domain.CategoryMatch {
	// Нормализуем slug
	normalizedSlug := strings.ToLower(strings.TrimSpace(suggestedCategory))
	if normalizedSlug == "" {
		return nil
	}

	// Маппим AI slug на DB slug
	dbSlug := s.mapAISlugToDBSlug(normalizedSlug)

	s.logger.Debug().
		Str("suggested", suggestedCategory).
		Str("normalized", normalizedSlug).
		Str("db_slug", dbSlug).
		Msg("Resolving suggested category")

	// Пробуем найти категорию по slug
	category, err := s.categoryRepo.GetBySlugV2(ctx, dbSlug)
	if err != nil || category == nil {
		// Пробуем также оригинальный slug
		if dbSlug != normalizedSlug {
			category, err = s.categoryRepo.GetBySlugV2(ctx, normalizedSlug)
		}
	}

	if category == nil || err != nil {
		s.logger.Warn().
			Str("suggested", suggestedCategory).
			Str("db_slug", dbSlug).
			Err(err).
			Msg("Could not resolve suggested category to DB category")
		return nil
	}

	// Локализуем название
	localized := category.Localize(language)

	s.logger.Info().
		Str("suggested", suggestedCategory).
		Str("resolved_slug", category.Slug).
		Str("resolved_name", localized.Name).
		Msg("Successfully resolved suggested category")

	return &domain.CategoryMatch{
		CategoryID:      category.ID,
		CategorySlug:    category.Slug,
		CategoryName:    localized.Name,
		CategoryPath:    category.Path,
		ConfidenceScore: 0.75, // Средняя уверенность для suggested category
		DetectionMethod: domain.MethodFallback,
		MatchedKeywords: []string{suggestedCategory},
	}
}

// mapAISlugToDBSlug преобразует slug из Claude в slug в БД
// Claude возвращает английские slugs, в БД используются сербские
func (s *CategoryDetectionService) mapAISlugToDBSlug(aiSlug string) string {
	mapping := map[string]string{
		// Electronics
		"electronics": "elektronika",

		// Fashion & Clothing
		"fashion":  "odeca-i-obuca",
		"clothing": "odeca-i-obuca",

		// Home & Garden
		"home-garden": "dom-i-basta",
		"home":        "dom-i-basta",
		"garden":      "dom-i-basta",

		// Furniture (мебель) - подкатегория Home & Garden
		"furniture":           "dom-i-basta",
		"chair":               "namestaj-dnevna-soba",
		"sofa":                "namestaj-dnevna-soba",
		"table":               "namestaj-dnevna-soba",
		"bed":                 "namestaj-spavaca-soba",
		"desk":                "namestaj-kancelarija",
		"kitchen-furniture":   "namestaj-kuhinja",
		"outdoor-furniture":   "bastenska-garnitura",
		"living-room":         "namestaj-dnevna-soba",
		"bedroom":             "namestaj-spavaca-soba",
		"nameštaj":            "dom-i-basta",
		"namestaj":            "dom-i-basta",

		// Sports & Tourism
		"sports": "sport-i-turizam",
		"sport":  "sport-i-turizam",

		// Vehicles & Automotive
		"vehicles":   "automobilizam",
		"automotive": "automobilizam",
		"auto":       "automobilizam",
		"car":        "automobilizam",
		"cars":       "automobilizam",
		"bike":       "automobilizam",
		"bicycle":    "automobilizam",

		// Kids & Baby
		"kids":     "za-bebe-i-decu",
		"baby":     "za-bebe-i-decu",
		"children": "za-bebe-i-decu",
		"toys":     "za-bebe-i-decu",

		// Beauty & Health
		"beauty": "lepota-i-zdravlje",
		"health": "lepota-i-zdravlje",

		// Books & Media
		"books-media": "knjige-i-mediji",
		"books":       "knjige-i-mediji",
		"media":       "knjige-i-mediji",

		// Jewelry & Watches
		"jewelry": "nakit-i-satovi",
		"watches": "nakit-i-satovi",
		"watch":   "nakit-i-satovi",
		// НЕ маппим "accessories" - это широкая категория (сумки, ремни, очки)
		// Для точной детекции используем productType hint

		// Appliances
		"appliances": "kucni-aparati",

		// Pets
		"pets": "kucni-ljubimci",

		// Food & Beverages
		"food":      "hrana-i-pice",
		"beverages": "hrana-i-pice",

		// Office
		"office": "kancelarijski-materijal",

		// Music
		"music": "muzicki-instrumenti",

		// Art & Crafts
		"art":    "umetnost-i-rukotvorine",
		"crafts": "umetnost-i-rukotvorine",

		// Services
		"services": "usluge",

		// Industrial & Tools
		"industrial": "industrija-i-alati",
		"tools":      "industrija-i-alati",

		// Other
		"other": "ostalo",
	}

	if dbSlug, ok := mapping[strings.ToLower(aiSlug)]; ok {
		return dbSlug
	}
	return aiSlug
}

// detectByProductType определяет категорию по productType из AI hint
// Возвращает точное совпадение для известных типов продуктов
func (s *CategoryDetectionService) detectByProductType(ctx context.Context, productType, language string) *domain.CategoryMatch {
	// Маппинг productType → slug категории
	productTypeMapping := map[string]string{
		// Часы → Nakit i satovi
		"watch":       "nakit-i-satovi",
		"watches":     "nakit-i-satovi",
		"wristwatch":  "nakit-i-satovi",
		"smartwatch":  "pametni-satovi", // Умные часы в Electronics
		"jewelry":     "nakit-i-satovi",
		"necklace":    "nakit-i-satovi",
		"ring":        "nakit-i-satovi",
		"bracelet":    "nakit-i-satovi",
		"earrings":    "nakit-i-satovi",

		// Электроника
		"smartphone":  "mobilni-telefoni",
		"phone":       "mobilni-telefoni",
		"mobile":      "mobilni-telefoni",
		"tablet":      "tableti",
		"laptop":      "laptopovi",
		"computer":    "racunari",
		"headphones":  "slusalice",
		"camera":      "fotoaparati",
		"tv":          "televizori",
		"television":  "televizori",
		"console":     "konzole",
		"gaming":      "konzole",

		// Одежда
		"dress":       "zenska-odeca",
		"shirt":       "muska-odeca",
		"shoes":       "obuca",
		"sneakers":    "obuca",
		"bag":         "torbe",
		"handbag":     "torbe",
	}

	slug, ok := productTypeMapping[strings.ToLower(productType)]
	if !ok {
		s.logger.Debug().
			Str("productType", productType).
			Msg("ProductType not found in mapping")
		return nil
	}

	// Ищем категорию в БД
	category, err := s.categoryRepo.GetBySlugV2(ctx, slug)
	if err != nil || category == nil {
		s.logger.Warn().
			Str("productType", productType).
			Str("slug", slug).
			Err(err).
			Msg("Category not found for productType")
		return nil
	}

	localized := category.Localize(language)

	return &domain.CategoryMatch{
		CategoryID:      category.ID,
		CategorySlug:    category.Slug,
		CategoryName:    localized.Name,
		CategoryPath:    category.Path,
		ConfidenceScore: 0.90, // Высокая уверенность для productType match
		DetectionMethod: domain.MethodKeywordMatch, // Используем существующий метод
		MatchedKeywords: []string{productType},
	}
}

// buildDetection создаёт CategoryDetection из matches
func (s *CategoryDetectionService) buildDetection(
	matches []domain.CategoryMatch,
	startTime time.Time,
	input domain.DetectFromTextInput,
) (*domain.CategoryDetection, error) {
	if len(matches) == 0 {
		return s.fallbackDetection(startTime, input)
	}

	primary := &matches[0]
	alternatives := []domain.CategoryMatch{}
	if len(matches) > 1 {
		alternatives = matches[1:]
		if len(alternatives) > 5 {
			alternatives = alternatives[:5]
		}
	}

	detection := domain.NewCategoryDetection(
		primary,
		alternatives,
		int32(time.Since(startTime).Milliseconds()),
		input,
	)

	// Сохраняем для tracking (async, не блокируем)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.repo.SaveDetection(ctx, detection); err != nil {
			s.logger.Error().Err(err).Msg("Failed to save detection")
		}
	}()

	return detection, nil
}

// fallbackDetection возвращает fallback результат
func (s *CategoryDetectionService) fallbackDetection(
	startTime time.Time,
	input domain.DetectFromTextInput,
) (*domain.CategoryDetection, error) {
	detection := &domain.CategoryDetection{
		ID:               uuid.New(),
		Primary:          nil,
		Alternatives:     nil,
		ProcessingTimeMs: int32(time.Since(startTime).Milliseconds()),
		InputTitle:       input.Title,
		InputDescription: input.Description,
		InputLanguage:    input.Language,
		CreatedAt:        time.Now(),
	}
	return detection, nil
}

// fallbackDetectionSync - sync версия для batch
func (s *CategoryDetectionService) fallbackDetectionSync(input domain.DetectFromTextInput) *domain.CategoryDetection {
	return &domain.CategoryDetection{
		ID:               uuid.New(),
		Primary:          nil,
		Alternatives:     nil,
		ProcessingTimeMs: 0,
		InputTitle:       input.Title,
		InputDescription: input.Description,
		InputLanguage:    input.Language,
		CreatedAt:        time.Now(),
	}
}

// claudeResponse - структура ответа Claude API
type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// extractKeywords извлекает ключевые слова из текста
func extractKeywords(text string) []string {
	// Удаляем спецсимволы
	reg := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	text = reg.ReplaceAllString(text, " ")

	// Разбиваем на слова
	words := strings.Fields(strings.ToLower(text))

	// Фильтруем короткие слова и стоп-слова
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "were": true,
		"i": true, "za": true, "na": true, "sa": true, "u": true,
		"и": true, "в": true, "на": true, "с": true, "для": true,
	}

	var keywords []string
	seen := make(map[string]bool)
	for _, word := range words {
		if len(word) < 3 || stopWords[word] || seen[word] {
			continue
		}
		seen[word] = true
		keywords = append(keywords, word)
	}

	return keywords
}

// mergeMatches объединяет matches из нескольких источников, убирая дубликаты
// Сортирует по confidence score (от большего к меньшему)
func mergeMatches(matchesLists ...[]domain.CategoryMatch) []domain.CategoryMatch {
	seen := make(map[string]bool)
	var result []domain.CategoryMatch

	// Объединить все списки
	for _, matches := range matchesLists {
		for _, m := range matches {
			key := m.CategorySlug
			if key == "" {
				key = m.CategoryID.String()
			}
			if !seen[key] {
				seen[key] = true
				result = append(result, m)
			}
		}
	}

	// Сортировать по confidence (desc)
	// Простая bubble sort для малых массивов
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].ConfidenceScore < result[j].ConfidenceScore {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}
