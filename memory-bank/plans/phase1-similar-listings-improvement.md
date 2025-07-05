# 🚀 **ФАЗА 1: Детальный план реализации улучшенного алгоритма похожих объявлений**

## 📋 **Обзор Фазы 1 (1-2 недели)**

**Цель**: Создать интеллектуальную систему поиска похожих объявлений с многоуровневым скорингом и улучшить общий поиск на сайте.

**Основные компоненты**:
1. ✅ Система весов атрибутов для категорий
2. ✅ Многоуровневый алгоритм скоринга
3. ✅ Улучшенная интеграция с OpenSearch
4. ✅ Улучшение общего поиска и автодополнения
5. ✅ Тестирование и документация

---

## 🔧 **Задача 1: Система весов атрибутов (2-3 дня)**

### **1.1 Создать структуры данных для весов**

**Файл**: `backend/internal/proj/marketplace/service/similarity_weights.go` (новый)

```go
package service

import "context"

// AttributeWeight представляет вес атрибута для конкретной категории
type AttributeWeight struct {
    AttributeName string  `json:"attribute_name"`
    Weight        float64 `json:"weight"`
    Category      string  `json:"category"`
}

// CategoryWeights содержит веса для всех атрибутов категории
type CategoryWeights struct {
    CategoryID   int                        `json:"category_id"`
    CategoryName string                     `json:"category_name"`
    Weights      map[string]float64         `json:"weights"`
    ParentID     *int                       `json:"parent_id,omitempty"`
}

// WeightManager управляет весами атрибутов
type WeightManager struct {
    weights map[int]*CategoryWeights
}

func NewWeightManager() *WeightManager {
    return &WeightManager{
        weights: make(map[int]*CategoryWeights),
    }
}
```

### **1.2 Определить веса для основных категорий**

```go
func (wm *WeightManager) InitializeDefaultWeights() {
    // Недвижимость - Квартиры (ID: 1100)
    wm.weights[1100] = &CategoryWeights{
        CategoryID:   1100,
        CategoryName: "Квартиры",
        Weights: map[string]float64{
            "rooms":         0.9,  // Количество комнат - критично
            "area":          0.85, // Площадь
            "floor":         0.7,  // Этаж
            "property_type": 0.8,  // Тип недвижимости
            "location":      0.75, // Местоположение
            "condition":     0.6,  // Состояние
            "heating":       0.5,  // Отопление
            "parking":       0.4,  // Парковка
            "balcony":       0.3,  // Балкон
            "elevator":      0.25, // Лифт
        },
    }

    // Автомобили (ID: 2000)
    wm.weights[2000] = &CategoryWeights{
        CategoryID:   2000,
        CategoryName: "Автомобили",
        Weights: map[string]float64{
            "make":         0.9,  // Марка
            "model":        0.85, // Модель
            "year":         0.8,  // Год выпуска
            "body_type":    0.75, // Тип кузова
            "fuel_type":    0.7,  // Тип топлива
            "transmission": 0.65, // Коробка передач
            "engine":       0.6,  // Двигатель
            "color":        0.3,  // Цвет
            "mileage":      0.7,  // Пробег
            "condition":    0.6,  // Состояние
        },
    }

    // Электроника (ID: 3000)
    wm.weights[3000] = &CategoryWeights{
        CategoryID:   3000,
        CategoryName: "Электроника",
        Weights: map[string]float64{
            "brand":        0.9,  // Бренд
            "model":        0.85, // Модель
            "type":         0.8,  // Тип устройства
            "condition":    0.7,  // Состояние
            "warranty":     0.5,  // Гарантия
            "color":        0.3,  // Цвет
            "storage":      0.6,  // Объем памяти
            "display_size": 0.5,  // Размер экрана
        },
    }

    // Мебель (ID: 4000)
    wm.weights[4000] = &CategoryWeights{
        CategoryID:   4000,
        CategoryName: "Мебель",
        Weights: map[string]float64{
            "type":      0.9,  // Тип мебели
            "material":  0.8,  // Материал
            "style":     0.7,  // Стиль
            "size":      0.75, // Размер
            "color":     0.6,  // Цвет
            "condition": 0.65, // Состояние
            "brand":     0.5,  // Бренд
        },
    }

    // Одежда (ID: 5000)
    wm.weights[5000] = &CategoryWeights{
        CategoryID:   5000,
        CategoryName: "Одежда",
        Weights: map[string]float64{
            "type":      0.9,  // Тип одежды
            "brand":     0.8,  // Бренд
            "size":      0.85, // Размер
            "color":     0.6,  // Цвет
            "material":  0.7,  // Материал
            "season":    0.65, // Сезон
            "gender":    0.8,  // Пол
            "condition": 0.6,  // Состояние
        },
    }
}
```

### **1.3 Добавить функцию получения весов**

```go
func (wm *WeightManager) GetCategoryWeights(categoryID int) map[string]float64 {
    if weights, exists := wm.weights[categoryID]; exists {
        return weights.Weights
    }
    
    // Попытка найти веса родительской категории
    parentWeights := wm.findParentWeights(categoryID)
    if parentWeights != nil {
        return parentWeights
    }
    
    // Возвращаем базовые веса
    return wm.getDefaultWeights()
}

func (wm *WeightManager) getDefaultWeights() map[string]float64 {
    return map[string]float64{
        "brand":     0.7,
        "model":     0.65,
        "type":      0.8,
        "condition": 0.6,
        "color":     0.4,
        "size":      0.5,
        "material":  0.5,
    }
}

func (wm *WeightManager) findParentWeights(categoryID int) map[string]float64 {
    // Здесь должна быть логика поиска родительской категории
    // Пока возвращаем nil для упрощения
    return nil
}
```

---

## 🎯 **Задача 2: Многоуровневый алгоритм скоринга (3-4 дня)**

### **2.1 Создать структуру для скоринга**

**Файл**: `backend/internal/proj/marketplace/service/similarity_scoring.go` (новый)

```go
package service

import (
    "context"
    "math"
    "strings"
    "fmt"
    "sort"
)

// SimilarityScore представляет результат расчета похожести
type SimilarityScore struct {
    ListingID        int                    `json:"listing_id"`
    CategoryScore    float64                `json:"category_score"`
    AttributeScore   float64                `json:"attribute_score"`
    TextScore        float64                `json:"text_score"`
    PriceScore       float64                `json:"price_score"`
    LocationScore    float64                `json:"location_score"`
    TotalScore       float64                `json:"total_score"`
    MatchReasons     []string               `json:"match_reasons"`
    ScoreBreakdown   map[string]interface{} `json:"score_breakdown"`
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
```

### **2.2 Основной метод расчета похожести**

```go
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
```

### **2.3 Функции расчета отдельных скоров**

```go
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
    
    // Создаем мапу атрибутов источника
    sourceAttrs := make(map[string]string)
    for _, attr := range source.Attributes {
        sourceAttrs[attr.AttributeName] = strings.ToLower(strings.TrimSpace(attr.DisplayValue))
    }
    
    // Сравниваем с атрибутами цели
    for _, targetAttr := range target.Attributes {
        if weight, exists := weights[targetAttr.AttributeName]; exists {
            totalWeight += weight
            
            if sourceValue, found := sourceAttrs[targetAttr.AttributeName]; found {
                targetValue := strings.ToLower(strings.TrimSpace(targetAttr.DisplayValue))
                
                if sc.attributeValuesMatch(sourceValue, targetValue) {
                    matchedWeight += weight
                }
            }
        }
    }
    
    if totalWeight == 0 {
        return 0.0
    }
    
    return matchedWeight / totalWeight
}

func (sc *SimilarityCalculator) attributeValuesMatch(value1, value2 string) bool {
    // Точное совпадение
    if value1 == value2 {
        return true
    }
    
    // Частичное совпадение для строк
    if strings.Contains(value1, value2) || strings.Contains(value2, value1) {
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
```

---

## 🔄 **Задача 3: Обновление GetSimilarListings (2-3 дня)**

### **3.1 Обновить основную функцию**

**Файл**: `backend/internal/proj/marketplace/service/marketplace.go`

```go
func (s *MarketplaceService) GetSimilarListings(
    ctx context.Context, 
    listingID int, 
    limit int,
) ([]*models.MarketplaceListing, error) {
    // Инициализируем калькулятор похожести
    calculator := NewSimilarityCalculator()
    
    // Получаем исходное объявление
    sourceListing, err := s.GetListingByID(ctx, listingID)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения объявления: %w", err)
    }
    
    // Формируем расширенные параметры поиска
    params := s.buildAdvancedSearchParams(sourceListing, limit*3) // Получаем больше для фильтрации
    
    // Выполняем поиск потенциально похожих объявлений
    results, err := s.SearchListingsAdvanced(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("ошибка поиска похожих объявлений: %w", err)
    }
    
    // Рассчитываем скоры похожести для каждого объявления
    var scoredListings []*ScoredListing
    for _, listing := range results.Items {
        if listing.ID == listingID {
            continue // Пропускаем исходное объявление
        }
        
        score, err := calculator.CalculateSimilarity(ctx, sourceListing, listing)
        if err != nil {
            continue // Пропускаем при ошибке расчета
        }
        
        // Добавляем только объявления с высоким скором
        if score.TotalScore > 0.25 {
            scoredListings = append(scoredListings, &ScoredListing{
                Listing: listing,
                Score:   score,
            })
        }
    }
    
    // Сортируем по скору релевантности
    sort.Slice(scoredListings, func(i, j int) bool {
        return scoredListings[i].Score.TotalScore > scoredListings[j].Score.TotalScore
    })
    
    // Ограничиваем результат и добавляем метаданные
    var similarListings []*models.MarketplaceListing
    for i, scored := range scoredListings {
        if i >= limit {
            break
        }
        
        // Добавляем информацию о скоре в метаданные
        if scored.Listing.Metadata == nil {
            scored.Listing.Metadata = make(map[string]interface{})
        }
        
        scored.Listing.Metadata["similarity_score"] = scored.Score.TotalScore
        scored.Listing.Metadata["match_reasons"] = scored.Score.MatchReasons
        scored.Listing.Metadata["score_breakdown"] = scored.Score.ScoreBreakdown
        
        similarListings = append(similarListings, scored.Listing)
    }
    
    return similarListings, nil
}

type ScoredListing struct {
    Listing *models.MarketplaceListing
    Score   *SimilarityScore
}
```

### **3.2 Создать расширенные параметры поиска**

```go
func (s *MarketplaceService) buildAdvancedSearchParams(
    listing *models.MarketplaceListing, 
    limit int,
) *search.ServiceParams {
    params := &search.ServiceParams{
        Size: limit,
        Page: 1,
        Sort: "date_desc",
    }
    
    // Поиск по категории (основной фильтр)
    params.CategoryID = strconv.Itoa(listing.CategoryID)
    
    // Добавляем фильтры по ключевым атрибутам
    if len(listing.Attributes) > 0 {
        attributeFilters := make(map[string]string)
        
        // Используем калькулятор для определения важных атрибутов
        calculator := NewSimilarityCalculator()
        weights := calculator.weightManager.GetCategoryWeights(listing.CategoryID)
        
        for _, attr := range listing.Attributes {
            if weight, exists := weights[attr.AttributeName]; exists && weight > 0.6 {
                attributeFilters[attr.AttributeName] = attr.DisplayValue
            }
        }
        
        if len(attributeFilters) > 0 {
            params.AttributeFilters = attributeFilters
        }
    }
    
    // Добавляем ценовой диапазон (±50%)
    if listing.Price > 0 {
        params.PriceMin = listing.Price * 0.5
        params.PriceMax = listing.Price * 1.5
    }
    
    // Добавляем географический фильтр
    if listing.City != "" {
        params.City = listing.City
    }
    
    return params
}
```

---

## 🔍 **Задача 4: Улучшение общего поиска (2-3 дня)**

### **4.1 Создать endpoint для автодополнения**

**Файл**: `backend/internal/proj/marketplace/handler/search.go`

```go
// GetEnhancedSuggestions возвращает умные подсказки для поиска
// @Summary Получить подсказки для автодополнения
// @Description Возвращает подсказки на основе частичного ввода: популярные запросы, категории, товары
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Поисковый запрос"
// @Param limit query int false "Лимит результатов (по умолчанию 10)"
// @Param types query string false "Типы подсказок: queries,categories,products (по умолчанию все)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]SuggestionItem}
// @Router /api/v1/marketplace/enhanced-suggestions [get]
func (h *SearchHandler) GetEnhancedSuggestions(c *fiber.Ctx) error {
    query := c.Query("query")
    if query == "" {
        return c.JSON(utils.SuccessResponse([]interface{}{}))
    }
    
    limit := c.QueryInt("limit", 10)
    types := c.Query("types", "queries,categories,products")
    
    ctx := c.Context()
    suggestions, err := h.service.GetEnhancedSuggestions(ctx, query, limit, types)
    if err != nil {
        h.logger.WithError(err).Error("Ошибка получения подсказок")
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "suggestions.getError")
    }
    
    return c.JSON(utils.SuccessResponse(suggestions))
}
```

### **4.2 Реализовать сервис автодополнения**

**Файл**: `backend/internal/proj/marketplace/service/search_suggestions.go` (новый)

```go
package service

import (
    "context"
    "strings"
    "sort"
)

type SuggestionType string

const (
    SuggestionTypeQuery    SuggestionType = "query"
    SuggestionTypeCategory SuggestionType = "category"
    SuggestionTypeProduct  SuggestionType = "product"
)

type SuggestionItem struct {
    Type        SuggestionType         `json:"type"`
    Value       string                 `json:"value"`
    Label       string                 `json:"label"`
    Count       int                    `json:"count,omitempty"`
    CategoryID  int                    `json:"category_id,omitempty"`
    ProductID   int                    `json:"product_id,omitempty"`
    Icon        string                 `json:"icon,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (s *MarketplaceService) GetEnhancedSuggestions(
    ctx context.Context,
    query string,
    limit int,
    types string,
) ([]SuggestionItem, error) {
    query = strings.ToLower(strings.TrimSpace(query))
    if query == "" {
        return []SuggestionItem{}, nil
    }
    
    var suggestions []SuggestionItem
    typesList := strings.Split(types, ",")
    
    // 1. Популярные поисковые запросы
    if contains(typesList, "queries") {
        querySuggestions, err := s.getQuerySuggestions(ctx, query)
        if err == nil {
            suggestions = append(suggestions, querySuggestions...)
        }
    }
    
    // 2. Категории
    if contains(typesList, "categories") {
        categorySuggestions, err := s.getCategorySuggestions(ctx, query)
        if err == nil {
            suggestions = append(suggestions, categorySuggestions...)
        }
    }
    
    // 3. Товары/Объявления
    if contains(typesList, "products") {
        productSuggestions, err := s.getProductSuggestions(ctx, query)
        if err == nil {
            suggestions = append(suggestions, productSuggestions...)
        }
    }
    
    // Сортируем по релевантности и ограничиваем
    suggestions = s.rankSuggestions(suggestions, query)
    if len(suggestions) > limit {
        suggestions = suggestions[:limit]
    }
    
    return suggestions, nil
}

func (s *MarketplaceService) getQuerySuggestions(
    ctx context.Context,
    query string,
) ([]SuggestionItem, error) {
    // Получаем популярные запросы из БД или кеша
    popularQueries, err := s.storage.GetPopularSearchQueries(ctx, query, 5)
    if err != nil {
        return nil, err
    }
    
    suggestions := make([]SuggestionItem, 0, len(popularQueries))
    for _, pq := range popularQueries {
        suggestions = append(suggestions, SuggestionItem{
            Type:  SuggestionTypeQuery,
            Value: pq.Query,
            Label: pq.Query,
            Count: pq.SearchCount,
            Metadata: map[string]interface{}{
                "last_searched": pq.LastSearched,
            },
        })
    }
    
    return suggestions, nil
}

func (s *MarketplaceService) getCategorySuggestions(
    ctx context.Context,
    query string,
) ([]SuggestionItem, error) {
    // Ищем категории по имени
    categories, err := s.storage.SearchCategories(ctx, query, 5)
    if err != nil {
        return nil, err
    }
    
    suggestions := make([]SuggestionItem, 0, len(categories))
    for _, cat := range categories {
        suggestions = append(suggestions, SuggestionItem{
            Type:       SuggestionTypeCategory,
            Value:      cat.Slug,
            Label:      cat.Name,
            CategoryID: cat.ID,
            Icon:       cat.Icon,
            Count:      cat.ListingCount,
            Metadata: map[string]interface{}{
                "path": cat.Path,
            },
        })
    }
    
    return suggestions, nil
}

func (s *MarketplaceService) getProductSuggestions(
    ctx context.Context,
    query string,
) ([]SuggestionItem, error) {
    // Быстрый поиск товаров
    searchParams := &search.ServiceParams{
        Query: query,
        Size:  5,
        Page:  1,
    }
    
    results, err := s.SearchListingsAdvanced(ctx, searchParams)
    if err != nil {
        return nil, err
    }
    
    suggestions := make([]SuggestionItem, 0, len(results.Items))
    for _, item := range results.Items {
        suggestions = append(suggestions, SuggestionItem{
            Type:      SuggestionTypeProduct,
            Value:     item.Title,
            Label:     item.Title,
            ProductID: item.ID,
            Metadata: map[string]interface{}{
                "price":    item.Price,
                "image":    getFirstImage(item.Images),
                "category": item.Category.Name,
            },
        })
    }
    
    return suggestions, nil
}

func (s *MarketplaceService) rankSuggestions(
    suggestions []SuggestionItem,
    query string,
) []SuggestionItem {
    // Ранжируем по релевантности
    sort.Slice(suggestions, func(i, j int) bool {
        // Приоритет типов: query > category > product
        if suggestions[i].Type != suggestions[j].Type {
            typeOrder := map[SuggestionType]int{
                SuggestionTypeQuery:    1,
                SuggestionTypeCategory: 2,
                SuggestionTypeProduct:  3,
            }
            return typeOrder[suggestions[i].Type] < typeOrder[suggestions[j].Type]
        }
        
        // По точности совпадения
        iExact := strings.HasPrefix(strings.ToLower(suggestions[i].Label), query)
        jExact := strings.HasPrefix(strings.ToLower(suggestions[j].Label), query)
        if iExact != jExact {
            return iExact
        }
        
        // По популярности (count)
        return suggestions[i].Count > suggestions[j].Count
    })
    
    return suggestions
}
```

### **4.3 Улучшить алгоритм ранжирования результатов поиска**

**Файл**: `backend/internal/proj/global/handler/unified_search.go`

Обновить функцию `UnifiedSearch`:

```go
func (h *UnifiedSearchHandler) UnifiedSearch(c *fiber.Ctx) error {
    // ... существующий код парсинга параметров ...
    
    // Улучшенный алгоритм объединения результатов
    mergedResults := h.mergeAndRankResults(
        marketplaceResults.Items,
        storefrontResults,
        params.Query,
        params.SortBy,
    )
    
    // ... остальной код ...
}

func (h *UnifiedSearchHandler) mergeAndRankResults(
    marketplaceItems []*models.MarketplaceListing,
    storefrontItems []map[string]interface{},
    query string,
    sortBy string,
) []map[string]interface{} {
    // Создаем унифицированную структуру с scoring
    type ScoredItem struct {
        Item          map[string]interface{}
        RelevanceScore float64
    }
    
    var scoredItems []ScoredItem
    
    // Обрабатываем marketplace объявления
    for _, item := range marketplaceItems {
        unifiedItem := h.convertMarketplaceToUnified(item)
        score := h.calculateRelevanceScore(unifiedItem, query)
        
        scoredItems = append(scoredItems, ScoredItem{
            Item:           unifiedItem,
            RelevanceScore: score,
        })
    }
    
    // Обрабатываем storefront товары
    for _, item := range storefrontItems {
        score := h.calculateRelevanceScore(item, query)
        
        // Бонус для верифицированных витрин
        if verified, ok := item["storefront_verified"].(bool); ok && verified {
            score *= 1.2
        }
        
        scoredItems = append(scoredItems, ScoredItem{
            Item:           item,
            RelevanceScore: score,
        })
    }
    
    // Сортируем по выбранному критерию
    h.sortScoredItems(scoredItems, sortBy)
    
    // Конвертируем обратно в массив
    results := make([]map[string]interface{}, len(scoredItems))
    for i, scored := range scoredItems {
        results[i] = scored.Item
        results[i]["relevance_score"] = scored.RelevanceScore
    }
    
    return results
}

func (h *UnifiedSearchHandler) calculateRelevanceScore(
    item map[string]interface{},
    query string,
) float64 {
    score := 0.0
    query = strings.ToLower(query)
    
    // Точное совпадение в заголовке (вес 5.0)
    if title, ok := item["title"].(string); ok {
        titleLower := strings.ToLower(title)
        if titleLower == query {
            score += 5.0
        } else if strings.Contains(titleLower, query) {
            score += 3.0
        }
    }
    
    // Совпадение в описании (вес 2.0)
    if desc, ok := item["description"].(string); ok {
        if strings.Contains(strings.ToLower(desc), query) {
            score += 2.0
        }
    }
    
    // Учитываем популярность (просмотры)
    if views, ok := item["views_count"].(float64); ok {
        score += math.Log10(views+1) * 0.5
    }
    
    // Свежесть объявления
    if createdAt, ok := item["created_at"].(string); ok {
        if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
            daysSinceCreated := time.Since(t).Hours() / 24
            freshnessScore := math.Max(0, 1-daysSinceCreated/30) // Снижаем скор для старых
            score += freshnessScore
        }
    }
    
    return score
}

func (h *UnifiedSearchHandler) sortScoredItems(
    items []ScoredItem,
    sortBy string,
) {
    switch sortBy {
    case "relevance":
        sort.Slice(items, func(i, j int) bool {
            return items[i].RelevanceScore > items[j].RelevanceScore
        })
    case "price_asc":
        sort.Slice(items, func(i, j int) bool {
            priceI := getPrice(items[i].Item)
            priceJ := getPrice(items[j].Item)
            return priceI < priceJ
        })
    case "price_desc":
        sort.Slice(items, func(i, j int) bool {
            priceI := getPrice(items[i].Item)
            priceJ := getPrice(items[j].Item)
            return priceI > priceJ
        })
    case "date_desc":
        sort.Slice(items, func(i, j int) bool {
            dateI := getCreatedAt(items[i].Item)
            dateJ := getCreatedAt(items[j].Item)
            return dateI.After(dateJ)
        })
    default:
        // По умолчанию - по релевантности
        sort.Slice(items, func(i, j int) bool {
            return items[i].RelevanceScore > items[j].RelevanceScore
        })
    }
}
```

### **4.4 Добавить сохранение популярных запросов**

**Файл**: `backend/migrations/XXX_add_search_queries_table.up.sql` (новый)

```sql
CREATE TABLE IF NOT EXISTS search_queries (
    id SERIAL PRIMARY KEY,
    query VARCHAR(255) NOT NULL,
    normalized_query VARCHAR(255) NOT NULL,
    search_count INTEGER DEFAULT 1,
    last_searched TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    language VARCHAR(10) DEFAULT 'ru',
    results_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_search_queries_normalized ON search_queries(normalized_query);
CREATE INDEX idx_search_queries_count ON search_queries(search_count DESC);
CREATE INDEX idx_search_queries_language ON search_queries(language);

-- Функция для обновления счетчика
CREATE OR REPLACE FUNCTION update_search_query_count()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_search_queries_updated_at
    BEFORE UPDATE ON search_queries
    FOR EACH ROW
    EXECUTE FUNCTION update_search_query_count();
```

### **4.5 Обновить фронтенд для работы с новым API**

**Файл**: `frontend/svetu/src/services/unifiedSearch.ts`

Добавить метод для получения подсказок:

```typescript
interface EnhancedSuggestion {
  type: 'query' | 'category' | 'product';
  value: string;
  label: string;
  count?: number;
  category_id?: number;
  product_id?: number;
  icon?: string;
  metadata?: Record<string, any>;
}

class UnifiedSearchService {
  // ... существующий код ...
  
  async getEnhancedSuggestions(
    query: string,
    limit: number = 10,
    types: string = 'queries,categories,products'
  ): Promise<EnhancedSuggestion[]> {
    try {
      const response = await apiClient.get<{
        success: boolean;
        data: EnhancedSuggestion[];
      }>('/api/v1/marketplace/enhanced-suggestions', {
        params: { query, limit, types }
      });
      
      return response.data.data || [];
    } catch (error) {
      console.error('Error fetching suggestions:', error);
      return [];
    }
  }
  
  // Обновить метод search для отслеживания запросов
  async search(params: UnifiedSearchParams): Promise<UnifiedSearchResponse> {
    const response = await apiClient.get<UnifiedSearchResponse>('/api/v1/search', {
      params: this.buildQueryParams(params)
    });
    
    // Сохраняем успешный поиск в истории
    if (params.query && response.data.data.items.length > 0) {
      this.addToSearchHistory(params.query);
    }
    
    return response.data;
  }
}
```

---

## 🧪 **Задача 5: Тестирование (1-2 дня)**

### **5.1 Создать unit тесты**

**Файл**: `backend/internal/proj/marketplace/service/similarity_test.go` (новый)

```go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "your-project/internal/domain/models"
)

func TestSimilarityCalculator_CalculateSimilarity(t *testing.T) {
    calculator := NewSimilarityCalculator()
    
    // Тестовые объявления
    sourceListing := &models.MarketplaceListing{
        ID:          1,
        CategoryID:  1100,
        Title:       "3-комнатная квартира в центре",
        Description: "Просторная квартира с ремонтом",
        Price:       200000,
        City:        "Белград",
        Attributes: []models.ListingAttributeValue{
            {AttributeName: "rooms", DisplayValue: "3"},
            {AttributeName: "area", DisplayValue: "85"},
            {AttributeName: "floor", DisplayValue: "5"},
        },
    }
    
    targetListing := &models.MarketplaceListing{
        ID:          2,
        CategoryID:  1100,
        Title:       "Квартира 3 комнаты центр города",
        Description: "Отличная квартира после ремонта",
        Price:       220000,
        City:        "Белград",
        Attributes: []models.ListingAttributeValue{
            {AttributeName: "rooms", DisplayValue: "3"},
            {AttributeName: "area", DisplayValue: "90"},
            {AttributeName: "floor", DisplayValue: "3"},
        },
    }
    
    // Выполняем расчет
    score, err := calculator.CalculateSimilarity(context.Background(), sourceListing, targetListing)
    
    // Проверяем результат
    assert.NoError(t, err)
    assert.NotNil(t, score)
    assert.Equal(t, 2, score.ListingID)
    assert.Greater(t, score.TotalScore, 0.5) // Ожидаем высокую похожесть
    assert.Equal(t, 1.0, score.CategoryScore) // Та же категория
    assert.Greater(t, score.AttributeScore, 0.7) // Похожие атрибуты
    assert.Contains(t, score.MatchReasons, "Та же категория")
}

func TestWeightManager_GetCategoryWeights(t *testing.T) {
    wm := NewWeightManager()
    wm.InitializeDefaultWeights()
    
    // Тест для квартир
    weights := wm.GetCategoryWeights(1100)
    assert.NotEmpty(t, weights)
    assert.Equal(t, 0.9, weights["rooms"])
    assert.Equal(t, 0.85, weights["area"])
    
    // Тест для несуществующей категории
    defaultWeights := wm.GetCategoryWeights(9999)
    assert.NotEmpty(t, defaultWeights)
    assert.Contains(t, defaultWeights, "brand")
}

func BenchmarkSimilarityCalculation(b *testing.B) {
    calculator := NewSimilarityCalculator()
    // ... настройка тестовых данных ...
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        calculator.CalculateSimilarity(context.Background(), sourceListing, targetListing)
    }
}
```

### **5.2 Создать интеграционные тесты**

**Файл**: `backend/internal/proj/marketplace/service/integration_test.go` (новый)

```go
func TestMarketplaceService_GetSimilarListings_Integration(t *testing.T) {
    // Создаем тестовую БД и сервис
    service := setupTestService(t)
    
    // Создаем тестовые объявления
    listing1 := createTestListing(t, service, "3-комнатная квартира", 1100)
    listing2 := createTestListing(t, service, "3-комнатная квартира люкс", 1100)
    listing3 := createTestListing(t, service, "Автомобиль BMW", 2000)
    
    // Тестируем поиск похожих
    similar, err := service.GetSimilarListings(context.Background(), listing1.ID, 10)
    
    assert.NoError(t, err)
    assert.Len(t, similar, 1) // Должна найтись только одна похожая квартира
    assert.Equal(t, listing2.ID, similar[0].ID)
    assert.NotNil(t, similar[0].Metadata["similarity_score"])
}
```

---

## 📚 **Задача 6: Документация (1 день)**

### **6.1 Обновить API документацию**

**Файл**: `backend/internal/proj/marketplace/handler/search.go`

Добавить Swagger комментарии:

```go
// GetSimilarListings получает похожие объявления с улучшенным алгоритмом релевантности
// @Summary Получить похожие объявления
// @Description Возвращает список похожих объявлений на основе многоуровневого анализа: категория, атрибуты, текст, цена, местоположение
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Param limit query int false "Лимит результатов (по умолчанию 5, максимум 20)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceListing} "Список похожих объявлений с метаданными скора"
// @Failure 404 {object} utils.ErrorResponseSwag "Объявление не найдено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/marketplace/listings/{id}/similar [get]
```

### **6.2 Создать техническую документацию**

**Файл**: `memory-bank/system-passports/backend/improved-similarity-algorithm.md` (новый)

```markdown
# Улучшенный алгоритм поиска похожих объявлений

## Обзор

Новый алгоритм использует многоуровневый подход для определения релевантности объявлений:

### Компоненты скоринга

1. **Категория (30%)** - точное или частичное совпадение категории
2. **Атрибуты (30%)** - взвешенное сравнение характеристик товара
3. **Текст (20%)** - семантическое сходство названия и описания
4. **Цена (15%)** - близость ценового диапазона
5. **Местоположение (5%)** - географическая близость

### Веса атрибутов по категориям

- **Недвижимость**: комнаты (0.9), площадь (0.85), этаж (0.7)
- **Автомобили**: марка (0.9), модель (0.85), год (0.8)
- **Электроника**: бренд (0.9), модель (0.85), тип (0.8)

### API Response

```json
{
  "data": [
    {
      "id": 123,
      "title": "...",
      "metadata": {
        "similarity_score": 0.85,
        "match_reasons": ["Та же категория", "Похожие характеристики"],
        "score_breakdown": {
          "category_raw": 1.0,
          "attribute_raw": 0.8,
          "text_raw": 0.6,
          "price_raw": 0.9,
          "location_raw": 1.0
        }
      }
    }
  ]
}
```
```

---

## ⏱️ **Timeline и приоритеты**

| День | Задача | Приоритет | Ответственный |
|------|--------|-----------|---------------|
| 1-2  | Система весов атрибутов | 🔴 Высокий | Backend dev |
| 3-4  | Алгоритм скоринга похожести | 🔴 Высокий | Backend dev |
| 5-6  | Обновление GetSimilarListings | 🔴 Высокий | Backend dev |
| 7-8  | Улучшение общего поиска | 🔴 Высокий | Backend dev |
| 9    | Unit тесты | 🟡 Средний | Backend dev |
| 10   | Интеграционные тесты | 🟡 Средний | QA |
| 11   | Документация | 🟢 Низкий | Tech writer |
| 12   | Код-ревью и деплой | 🔴 Высокий | Team lead |

---

## 🎯 **Критерии успеха Фазы 1**

1. ✅ **Функциональность**: Новый алгоритм похожести корректно работает
2. ✅ **Поиск**: Улучшенный поиск с автодополнением и ранжированием
3. ✅ **Производительность**: Время ответа < 500ms для поиска и рекомендаций
4. ✅ **Качество**: Релевантность улучшена на 40%+ для обоих алгоритмов
5. ✅ **Покрытие тестами**: >80% покрытие нового кода
6. ✅ **Документация**: Полная техническая документация создана

---

## 🚀 **Готовность к внедрению**

После завершения Фазы 1 будет готова к внедрению основа для более продвинутых алгоритмов (ML, векторный поиск) в Фазах 2 и 3.