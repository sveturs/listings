package service

import (
	"context"
	"math"
	"strings"

	"backend/internal/domain/models"
)

// SimilarityScore представляет результат расчета похожести
type SimilarityScore struct {
	ListingID      int                    `json:"listing_id"`
	CategoryScore  float64                `json:"category_score"`
	AttributeScore float64                `json:"attribute_score"`
	TextScore      float64                `json:"text_score"`
	PriceScore     float64                `json:"price_score"`
	LocationScore  float64                `json:"location_score"`
	TotalScore     float64                `json:"total_score"`
	MatchReasons   []string               `json:"match_reasons"`
	ScoreBreakdown map[string]interface{} `json:"score_breakdown"`
}

// SimilarityCalculator отвечает за расчет похожести
type SimilarityCalculator struct {
	weightManager *WeightManager
}

func NewSimilarityCalculator() *SimilarityCalculator {
	wm := NewWeightManager()
	wm.InitializeDefaultWeights()

	return &SimilarityCalculator{
		weightManager: wm,
	}
}

func (sc *SimilarityCalculator) CalculateSimilarity(
	ctx context.Context,
	sourceListing *models.MarketplaceListing,
	targetListing *models.MarketplaceListing,
) (*SimilarityScore, error) {
	score := &SimilarityScore{
		ListingID:      targetListing.ID,
		MatchReasons:   make([]string, 0),
		ScoreBreakdown: make(map[string]interface{}),
	}

	// 1. Категория (вес 30%)
	score.CategoryScore = sc.calculateCategoryScore(sourceListing, targetListing)
	if score.CategoryScore > 0.5 {
		score.MatchReasons = append(score.MatchReasons, "Та же категория")
	}

	// 2. Атрибуты (вес 30%)
	score.AttributeScore = sc.calculateAttributeScore(sourceListing, targetListing)
	if score.AttributeScore > 0.6 {
		score.MatchReasons = append(score.MatchReasons, "Похожие характеристики")
	}

	// 3. Текст (вес 20%)
	score.TextScore = sc.calculateTextScore(sourceListing, targetListing)
	if score.TextScore > 0.4 {
		score.MatchReasons = append(score.MatchReasons, "Похожее описание")
	}

	// 4. Цена (вес 15%)
	score.PriceScore = sc.calculatePriceScore(sourceListing, targetListing)
	if score.PriceScore > 0.7 {
		score.MatchReasons = append(score.MatchReasons, "Похожая цена")
	}

	// 5. Местоположение (вес 5%)
	score.LocationScore = sc.calculateLocationScore(sourceListing, targetListing)
	if score.LocationScore > 0.8 {
		score.MatchReasons = append(score.MatchReasons, "То же место")
	}

	// Итоговый скор
	score.TotalScore = score.CategoryScore*0.3 +
		score.AttributeScore*0.3 +
		score.TextScore*0.2 +
		score.PriceScore*0.15 +
		score.LocationScore*0.05

	// Детализированная информация
	score.ScoreBreakdown = map[string]interface{}{
		"category_weight":  0.3,
		"attribute_weight": 0.3,
		"text_weight":      0.2,
		"price_weight":     0.15,
		"location_weight":  0.05,
		"category_raw":     score.CategoryScore,
		"attribute_raw":    score.AttributeScore,
		"text_raw":         score.TextScore,
		"price_raw":        score.PriceScore,
		"location_raw":     score.LocationScore,
	}

	return score, nil
}

func (sc *SimilarityCalculator) calculateCategoryScore(
	source, target *models.MarketplaceListing,
) float64 {
	if source.CategoryID == target.CategoryID {
		return 1.0
	}

	// Проверяем родительские категории (если есть логика)
	// В упрощенной версии возвращаем 0
	return 0.0
}

func (sc *SimilarityCalculator) calculateAttributeScore(
	source, target *models.MarketplaceListing,
) float64 {
	if len(source.Attributes) == 0 || len(target.Attributes) == 0 {
		return 0.0
	}

	weights := sc.weightManager.GetCategoryWeights(source.CategoryID)
	var totalWeight, matchedWeight float64

	// Создаем мапы атрибутов источника и цели
	sourceAttrs := make(map[string]string)
	for _, attr := range source.Attributes {
		sourceAttrs[attr.AttributeName] = strings.ToLower(strings.TrimSpace(attr.DisplayValue))
	}

	targetAttrs := make(map[string]string)
	for _, attr := range target.Attributes {
		targetAttrs[attr.AttributeName] = strings.ToLower(strings.TrimSpace(attr.DisplayValue))
	}

	// Проходим по всем весам категории
	for attrName, weight := range weights {
		sourceValue, sourceExists := sourceAttrs[attrName]
		targetValue, targetExists := targetAttrs[attrName]

		// Если атрибут есть хотя бы в одном объявлении, учитываем его вес
		if sourceExists || targetExists {
			totalWeight += weight

			// Если атрибут есть в обоих объявлениях и значения совпадают
			if sourceExists && targetExists && sc.attributeValuesMatch(sourceValue, targetValue) {
				matchedWeight += weight
			}
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return matchedWeight / totalWeight
}

func (sc *SimilarityCalculator) attributeValuesMatch(value1, value2 string) bool {
	// Приводим к нижнему регистру для сравнения
	value1Lower := strings.ToLower(value1)
	value2Lower := strings.ToLower(value2)

	// Точное совпадение (без учета регистра)
	if value1Lower == value2Lower {
		return true
	}

	// Частичное совпадение для строк
	if strings.Contains(value1Lower, value2Lower) || strings.Contains(value2Lower, value1Lower) {
		return true
	}

	// Для чисел - проверяем диапазон (±10%)
	// Можно добавить более сложную логику

	return false
}

func (sc *SimilarityCalculator) calculateTextScore(
	source, target *models.MarketplaceListing,
) float64 {
	sourceText := strings.ToLower(source.Title + " " + source.Description)
	targetText := strings.ToLower(target.Title + " " + target.Description)

	return sc.calculateTextSimilarity(sourceText, targetText)
}

func (sc *SimilarityCalculator) calculateTextSimilarity(text1, text2 string) float64 {
	// Простой алгоритм на основе общих слов
	words1 := strings.Fields(text1)
	words2 := strings.Fields(text2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Подсчитываем общие слова
	wordSet1 := make(map[string]bool)
	for _, word := range words1 {
		if len(word) > 3 { // Игнорируем короткие слова
			wordSet1[word] = true
		}
	}

	commonWords := 0
	for _, word := range words2 {
		if len(word) > 3 && wordSet1[word] {
			commonWords++
		}
	}

	// Jaccard similarity
	totalUniqueWords := len(wordSet1)
	for _, word := range words2 {
		if len(word) > 3 && !wordSet1[word] {
			totalUniqueWords++
		}
	}

	if totalUniqueWords == 0 {
		return 0.0
	}

	return float64(commonWords) / float64(totalUniqueWords)
}

func (sc *SimilarityCalculator) calculatePriceScore(
	source, target *models.MarketplaceListing,
) float64 {
	if source.Price == 0 || target.Price == 0 {
		return 0.0
	}

	// Чем ближе цены, тем выше скор
	ratio := source.Price / target.Price
	if ratio > 1 {
		ratio = 1 / ratio
	}

	// Логарифмическая шкала для более естественного восприятия
	return math.Max(0, 1-math.Abs(math.Log(ratio))/math.Log(2))
}

func (sc *SimilarityCalculator) calculateLocationScore(
	source, target *models.MarketplaceListing,
) float64 {
	// Город
	if source.City == target.City {
		return 1.0
	}

	// Страна
	if source.Country == target.Country {
		return 0.5
	}

	// Координаты (если есть)
	if source.Latitude != nil && source.Longitude != nil &&
		target.Latitude != nil && target.Longitude != nil {
		distance := sc.calculateDistance(
			*source.Latitude, *source.Longitude,
			*target.Latitude, *target.Longitude,
		)

		// Чем меньше расстояние, тем выше скор
		return math.Max(0, 1-distance/100) // 100 км = 0 скор
	}

	return 0.0
}

func (sc *SimilarityCalculator) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Haversine formula для расчета расстояния между координатами
	const R = 6371 // Радиус Земли в км

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// ScoredListing представляет объявление с рассчитанным скором
type ScoredListing struct {
	Listing *models.MarketplaceListing
	Score   *SimilarityScore
}
