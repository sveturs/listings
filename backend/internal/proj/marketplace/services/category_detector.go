package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"backend/internal/domain/models"
	"backend/internal/proj/marketplace/repository"
	"backend/internal/proj/marketplace/storage/opensearch"
	"backend/internal/storage/postgres"
)

// CategoryDetector предоставляет интеллектуальное определение категории
type CategoryDetector struct {
	logger      *zap.Logger
	osRepo      *opensearch.Repository
	pgRepo      *repository.MarketplaceRepository
	keywordRepo *postgres.CategoryKeywordRepository
	statsRepo   *postgres.CategoryDetectionStatsRepository
}

// NewCategoryDetector создает новый экземпляр CategoryDetector
func NewCategoryDetector(
	logger *zap.Logger,
	osRepo *opensearch.Repository,
	pgRepo *repository.MarketplaceRepository,
	keywordRepo *postgres.CategoryKeywordRepository,
	statsRepo *postgres.CategoryDetectionStatsRepository,
) *CategoryDetector {
	return &CategoryDetector{
		logger:      logger,
		osRepo:      osRepo,
		pgRepo:      pgRepo,
		keywordRepo: keywordRepo,
		statsRepo:   statsRepo,
	}
}

// NewCategoryDetectorFromStorage создает CategoryDetector из Storage
func NewCategoryDetectorFromStorage(db *postgres.Database, osClient interface{}) (*CategoryDetector, error) {
	// Получаем logger
	logger := zap.L()

	// Создаём OpenSearch репозиторий
	osRepo, ok := osClient.(*opensearch.Repository)
	if !ok {
		return nil, errors.New("invalid OpenSearch client type")
	}

	// Создаём PostgreSQL репозитории
	sqlxDB := db.GetSQLXDB()

	pgRepo := repository.NewMarketplaceRepository(sqlxDB)
	keywordRepo := postgres.NewCategoryKeywordRepository(sqlxDB)
	statsRepo := postgres.NewCategoryDetectionStatsRepository(sqlxDB)

	return NewCategoryDetector(logger, osRepo, pgRepo, keywordRepo, statsRepo), nil
}

// DetectionInput входные данные для определения категории
type DetectionInput struct {
	// Из AI анализа
	Keywords    []string               `json:"keywords"`
	Attributes  map[string]interface{} `json:"attributes"`
	Domain      string                 `json:"domain"`
	ProductType string                 `json:"product_type"`
	Language    string                 `json:"language"`

	// Дополнительные данные
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	// Контекст
	UserID    *int32 `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

// DetectionResult результат определения категории
type DetectionResult struct {
	CategoryID            int32                 `json:"category_id"`
	CategoryName          string                `json:"category_name"`
	CategorySlug          string                `json:"category_slug"`
	ConfidenceScore       float64               `json:"confidence_score"`
	Method                string                `json:"method"`
	AlternativeCategories []AlternativeCategory `json:"alternative_categories,omitempty"`

	// Детали для отладки
	KeywordScore     float64  `json:"keyword_score,omitempty"`
	SimilarityScore  float64  `json:"similarity_score,omitempty"`
	MatchedKeywords  []string `json:"matched_keywords,omitempty"`
	ProcessingTimeMs int64    `json:"processing_time_ms"`

	// ID записи статистики для отслеживания подтверждений
	StatsID int32 `json:"stats_id,omitempty"`
}

// AlternativeCategory альтернативная категория
type AlternativeCategory struct {
	CategoryID      int32   `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	CategorySlug    string  `json:"category_slug"`
	ConfidenceScore float64 `json:"confidence_score"`
}

// DetectCategory определяет категорию используя комбинированный подход
func (cd *CategoryDetector) DetectCategory(ctx context.Context, input DetectionInput) (*DetectionResult, error) {
	startTime := time.Now()

	// Проверка на nil
	if cd == nil {
		return nil, errors.New("CategoryDetector is nil")
	}
	if cd.logger == nil {
		return nil, errors.New("logger is nil")
	}
	if cd.keywordRepo == nil {
		return nil, errors.New("keywordRepo is nil")
	}

	cd.logger.Info("начало определения категории",
		zap.Strings("keywords", input.Keywords),
		zap.String("title", input.Title),
		zap.String("language", input.Language))

	// Валидация входных данных
	if len(input.Keywords) == 0 && input.Title == "" && input.Description == "" {
		return nil, errors.New("недостаточно данных для определения категории")
	}

	// Язык по умолчанию
	if input.Language == "" {
		input.Language = "ru"
	}

	// Параллельно выполняем поиск по ключевым словам и similarity search
	keywordCh := make(chan *keywordSearchResult)
	similarityCh := make(chan *similaritySearchResult)

	// Поиск по ключевым словам
	go func() {
		result, err := cd.searchByKeywords(ctx, input)
		if err != nil {
			cd.logger.Error("ошибка поиска по ключевым словам", zap.Error(err))
			keywordCh <- nil
		} else {
			keywordCh <- result
		}
	}()

	// Similarity search
	go func() {
		if input.Title != "" || input.Description != "" {
			result, err := cd.searchBySimilarity(ctx, input)
			if err != nil {
				cd.logger.Error("ошибка similarity search", zap.Error(err))
				similarityCh <- nil
			} else {
				similarityCh <- result
			}
		} else {
			similarityCh <- nil
		}
	}()

	// Получаем результаты
	keywordResult := <-keywordCh
	similarityResult := <-similarityCh

	// Комбинируем результаты
	finalResult := cd.combineResults(keywordResult, similarityResult)

	// Определяем метод
	switch {
	case keywordResult != nil && similarityResult != nil:
		finalResult.Method = "combined"
	case keywordResult != nil:
		finalResult.Method = "keywords"
	case similarityResult != nil:
		finalResult.Method = "similarity"
	default:
		return nil, errors.New("не удалось определить категорию")
	}

	// Обогащаем информацией о категории
	category, err := cd.pgRepo.GetCategoryByID(ctx, finalResult.CategoryID)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения категории")
	}

	finalResult.CategoryName = category.Name
	finalResult.CategorySlug = category.Slug
	finalResult.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	// Обновляем счетчик использования ключевых слов
	if keywordResult != nil && len(keywordResult.matchedKeywords[finalResult.CategoryID]) > 0 {
		// Сохраняем параметры для фоновой задачи
		categoryID := finalResult.CategoryID
		keywords := keywordResult.matchedKeywords[finalResult.CategoryID]
		language := input.Language

		//nolint:contextcheck // Используем новый контекст для фоновой задачи
		go cd.incrementKeywordUsageInBackground(categoryID, keywords, language)
	}

	// Сохраняем статистику для анализа и получаем ID
	statsID := cd.saveDetectionStats(ctx, input, finalResult)
	finalResult.StatsID = statsID

	return finalResult, nil
}

// keywordSearchResult результат поиска по ключевым словам
type keywordSearchResult struct {
	categoryScores  map[int32]float64
	matchedKeywords map[int32][]string
}

// searchByKeywords поиск категории по ключевым словам
func (cd *CategoryDetector) searchByKeywords(ctx context.Context, input DetectionInput) (*keywordSearchResult, error) {
	cd.logger.Info("searchByKeywords начало")

	// Объединяем все ключевые слова для поиска
	allKeywords := make([]string, 0, len(input.Keywords))
	allKeywords = append(allKeywords, input.Keywords...)

	// Извлекаем ключевые слова из title и description
	if input.Title != "" {
		titleKeywords := cd.extractKeywords(input.Title)
		allKeywords = append(allKeywords, titleKeywords...)
	}

	if input.Description != "" {
		descKeywords := cd.extractKeywords(input.Description)
		allKeywords = append(allKeywords, descKeywords...)
	}

	cd.logger.Info("поиск по ключевым словам",
		zap.Strings("keywords", allKeywords),
		zap.String("language", input.Language))

	// Ищем совпадения с ключевыми словами категорий
	matches, err := cd.keywordRepo.FindMatchingCategories(ctx, allKeywords, input.Language)
	if err != nil {
		cd.logger.Error("ошибка FindMatchingCategories", zap.Error(err))
		return nil, errors.Wrap(err, "ошибка поиска категорий по ключевым словам")
	}

	cd.logger.Info("найдено совпадений", zap.Int("count", len(matches)))

	// Подсчитываем score для каждой категории
	categoryScores := make(map[int32]float64)
	matchedKeywords := make(map[int32][]string)

	for _, match := range matches {
		score := match.Weight

		// Учитываем тип ключевого слова
		switch match.KeywordType {
		case "main":
			score *= 2.0
		case "synonym":
			score *= 1.5
		case "brand":
			score *= 1.2
		case "attribute":
			score *= 1.1
		}

		// Учитываем негативные ключевые слова
		if match.IsNegative {
			categoryScores[match.CategoryID] -= score * 3.0
		} else {
			categoryScores[match.CategoryID] += score
			matchedKeywords[match.CategoryID] = append(matchedKeywords[match.CategoryID], match.Keyword)
		}
	}

	// Нормализуем scores
	maxScore := 0.0
	for _, score := range categoryScores {
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore > 0 {
		for catID, score := range categoryScores {
			categoryScores[catID] = score / maxScore
		}
	}

	return &keywordSearchResult{
		categoryScores:  categoryScores,
		matchedKeywords: matchedKeywords,
	}, nil
}

// similaritySearchResult результат similarity search
type similaritySearchResult struct {
	categoryScores map[int32]float64
	topListingID   *int32
	topScore       float64
}

// searchBySimilarity поиск похожих объявлений
func (cd *CategoryDetector) searchBySimilarity(ctx context.Context, input DetectionInput) (*similaritySearchResult, error) {
	// Формируем текст для поиска
	searchText := input.Title
	if input.Description != "" {
		searchText += " " + input.Description
	}

	// Используем OpenSearch для поиска похожих объявлений
	similarListings, err := cd.osRepo.FindSimilarListings(ctx, searchText, 10)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка поиска похожих объявлений")
	}

	if len(similarListings) == 0 {
		return nil, nil //nolint:nilnil // Возвращаем nil для обоих значений, если нет похожих объявлений
	}

	// Подсчитываем score для каждой категории
	categoryScores := make(map[int32]float64)
	var topListingID *int32
	var topScore float64

	for i, listing := range similarListings {
		// Score убывает с позицией
		score := 1.0 / float64(i+1)

		categoryScores[listing.CategoryID] += score

		if i == 0 {
			id := listing.ID
			topListingID = &id
			topScore = listing.Score
		}
	}

	// Нормализуем scores
	maxScore := 0.0
	for _, score := range categoryScores {
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore > 0 {
		for catID, score := range categoryScores {
			categoryScores[catID] = score / maxScore
		}
	}

	return &similaritySearchResult{
		categoryScores: categoryScores,
		topListingID:   topListingID,
		topScore:       topScore,
	}, nil
}

// combineResults комбинирует результаты разных методов
func (cd *CategoryDetector) combineResults(keyword *keywordSearchResult, similarity *similaritySearchResult) *DetectionResult {
	combinedScores := make(map[int32]float64)

	// Веса для комбинирования (60% keywords, 40% similarity)
	keywordWeight := 0.6
	similarityWeight := 0.4

	// Комбинируем scores
	if keyword != nil {
		for catID, score := range keyword.categoryScores {
			combinedScores[catID] += score * keywordWeight
		}
	}

	if similarity != nil {
		for catID, score := range similarity.categoryScores {
			combinedScores[catID] += score * similarityWeight
		}
	}

	// Находим категорию с максимальным score
	var bestCategoryID int32
	var bestScore float64

	for catID, score := range combinedScores {
		if score > bestScore {
			bestCategoryID = catID
			bestScore = score
		}
	}

	// Если не удалось определить категорию с достаточной уверенностью,
	// используем категорию "Elektronika" (ID=1001) как дефолтную
	const defaultCategoryID = 1001 // Electronics - самая общая категория
	const minConfidenceThreshold = 0.3

	if bestScore < minConfidenceThreshold || bestCategoryID == 0 {
		cd.logger.Warn("не удалось определить категорию с достаточной уверенностью, используется дефолтная",
			zap.Float64("bestScore", bestScore),
			zap.Int32("bestCategoryID", bestCategoryID))

		bestCategoryID = defaultCategoryID
		bestScore = 0.1 // Низкая уверенность для дефолтной категории
	}

	// Формируем список альтернативных категорий
	alternatives := make([]AlternativeCategory, 0)
	for catID, score := range combinedScores {
		if catID != bestCategoryID && score > 0.3 { // порог 30%
			alternatives = append(alternatives, AlternativeCategory{
				CategoryID:      catID,
				ConfidenceScore: score,
			})
		}
	}

	// Сортируем альтернативы по score
	sort.Slice(alternatives, func(i, j int) bool {
		return alternatives[i].ConfidenceScore > alternatives[j].ConfidenceScore
	})

	// Ограничиваем количество альтернатив
	if len(alternatives) > 3 {
		alternatives = alternatives[:3]
	}

	result := &DetectionResult{
		CategoryID:            bestCategoryID,
		ConfidenceScore:       bestScore,
		AlternativeCategories: alternatives,
	}

	// Добавляем детали для отладки
	if keyword != nil {
		result.KeywordScore = keyword.categoryScores[bestCategoryID]
		result.MatchedKeywords = keyword.matchedKeywords[bestCategoryID]
	}

	if similarity != nil {
		result.SimilarityScore = similarity.categoryScores[bestCategoryID]
	}

	return result
}

// extractKeywords извлекает ключевые слова из текста
func (cd *CategoryDetector) extractKeywords(text string) []string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	// Заменяем знаки препинания на пробелы, включая дефисы
	// Это позволит "Mercedes-Benz" стать "mercedes benz"
	replacer := strings.NewReplacer(
		"-", " ", // Дефис заменяем на пробел для правильной токенизации брендов
		".", " ",
		",", " ",
		"!", " ",
		"?", " ",
		";", " ",
		":", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"{", " ",
		"}", " ",
		"/", " ",
		"\\", " ",
	)
	text = replacer.Replace(text)

	// Разбиваем на слова
	words := strings.Fields(text)

	// Убираем короткие слова и стоп-слова
	keywords := make([]string, 0)
	stopWords := map[string]bool{
		"и": true, "в": true, "на": true, "с": true, "по": true,
		"для": true, "от": true, "до": true, "из": true,
		"the": true, "a": true, "an": true, "in": true, "on": true,
		"for": true, "to": true, "from": true, "with": true,
	}

	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// saveDetectionStats сохраняет статистику для анализа и возвращает ID записи
func (cd *CategoryDetector) saveDetectionStats(ctx context.Context, input DetectionInput, result *DetectionResult) int32 {
	stats := &models.CategoryDetectionStats{
		UserID:                input.UserID,
		SessionID:             input.SessionID,
		Method:                result.Method,
		AIKeywords:            input.Keywords,
		AIAttributes:          input.Attributes,
		AIDomain:              input.Domain,
		AIProductType:         input.ProductType,
		AISuggestedCategoryID: &result.CategoryID,
		FinalCategoryID:       &result.CategoryID,
		ConfidenceScore:       &result.ConfidenceScore,
		KeywordScore:          &result.KeywordScore,
		SimilarityScore:       &result.SimilarityScore,
		MatchedKeywords:       result.MatchedKeywords,
		ProcessingTimeMs:      &result.ProcessingTimeMs,
	}

	// Сохраняем альтернативные категории
	if len(result.AlternativeCategories) > 0 {
		alternatives := make(map[string]interface{})
		for _, alt := range result.AlternativeCategories {
			alternatives[fmt.Sprintf("%d", alt.CategoryID)] = alt.ConfidenceScore
		}
		alternativesJSON, _ := json.Marshal(alternatives)
		stats.AlternativeCategories = alternativesJSON
	}

	if err := cd.statsRepo.Create(ctx, stats); err != nil {
		cd.logger.Error("ошибка сохранения статистики определения категории", zap.Error(err))
		return 0
	}

	// Возвращаем ID созданной записи
	return stats.ID
}

// UpdateKeywordWeights обновляет веса ключевых слов на основе статистики
func (cd *CategoryDetector) UpdateKeywordWeights(ctx context.Context) error {
	// Получаем статистику за последние 30 дней
	stats, err := cd.statsRepo.GetRecentStats(ctx, 30)
	if err != nil {
		return errors.Wrap(err, "ошибка получения статистики")
	}

	// Анализируем успешные и неуспешные определения
	keywordSuccess := make(map[string]int)
	keywordTotal := make(map[string]int)

	for _, stat := range stats {
		if stat.UserConfirmed != nil && *stat.UserConfirmed {
			// Успешное определение
			for _, keyword := range stat.MatchedKeywords {
				keywordSuccess[keyword]++
				keywordTotal[keyword]++
			}
		} else if stat.UserSelectedCategoryID != nil {
			// Пользователь выбрал другую категорию
			for _, keyword := range stat.MatchedKeywords {
				keywordTotal[keyword]++
			}
		}
	}

	// Обновляем success_rate для ключевых слов
	for keyword, total := range keywordTotal {
		successRate := float64(keywordSuccess[keyword]) / float64(total)

		if err := cd.keywordRepo.UpdateSuccessRate(ctx, keyword, successRate); err != nil {
			cd.logger.Error("ошибка обновления success_rate",
				zap.String("keyword", keyword),
				zap.Error(err))
		}
	}

	return nil
}

// incrementKeywordUsageInBackground увеличивает счетчик использования ключевых слов в фоновом режиме
func (cd *CategoryDetector) incrementKeywordUsageInBackground(categoryID int32, keywords []string, language string) {
	bgCtx := context.Background()
	if err := cd.keywordRepo.IncrementUsageCount(bgCtx, categoryID, keywords, language); err != nil {
		cd.logger.Error("ошибка обновления счетчика использования ключевых слов",
			zap.Error(err),
			zap.Int32("categoryID", categoryID))
	}
}
