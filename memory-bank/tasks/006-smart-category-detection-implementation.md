# Задание на реализацию умного поиска категорий

## Общее описание

Необходимо реализовать систему умного определения категорий на основе семантического анализа, которая заменит текущий подход с жестко закодированными категориями в AI промптах.

## Архитектура решения

### 1. AI извлекает семантическую информацию

Вместо того чтобы AI выбирал конкретную категорию из списка, он будет извлекать:
```json
{
  "domain": "automotive",           // общая область
  "productType": "tire",           // тип продукта
  "keywords": ["summer", "michelin", "205/55", "R16", "tire"],
  "attributes": {
    "brand": "michelin",
    "size": "205/55 R16",
    "season": "summer"
  }
}
```

### 2. Backend определяет категорию

На основе семантической информации backend найдет наиболее подходящую категорию используя:
- Ключевые слова с весами
- Иерархию категорий
- Атрибуты категории
- Историю выборов пользователей

### 3. Обучение системы

Система будет улучшаться на основе:
- Подтверждений/исправлений пользователей
- Статистики кликов
- Популярности категорий

## Детальный план реализации

### Этап 1: База данных

#### 1.1 Создать таблицу category_keywords

```sql
-- migrations/000166_create_category_keywords.up.sql
CREATE TABLE category_keywords (
  id SERIAL PRIMARY KEY,
  category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
  keyword VARCHAR(100) NOT NULL,
  language VARCHAR(2) DEFAULT 'en',
  weight FLOAT DEFAULT 1.0 CHECK (weight >= 0.0 AND weight <= 10.0),
  source VARCHAR(50) DEFAULT 'manual', -- manual, ai_extracted, user_confirmed
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_keyword_lower (LOWER(keyword)),
  INDEX idx_category_weight (category_id, weight DESC),
  INDEX idx_language (language),
  UNIQUE idx_unique_keyword_category (category_id, keyword, language)
);

-- Триггер для updated_at
CREATE TRIGGER update_category_keywords_updated_at 
  BEFORE UPDATE ON category_keywords 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at();
```

#### 1.2 Создать таблицу category_selection_history

```sql
-- migrations/000167_create_category_selection_history.up.sql
CREATE TABLE category_selection_history (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  session_id VARCHAR(100),
  ai_suggested_category_id INTEGER REFERENCES marketplace_categories(id),
  user_selected_category_id INTEGER REFERENCES marketplace_categories(id),
  keywords TEXT[], -- массив ключевых слов из AI
  attributes JSONB, -- атрибуты из AI
  confidence_score FLOAT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_user_id (user_id),
  INDEX idx_created_at (created_at DESC),
  INDEX idx_categories (ai_suggested_category_id, user_selected_category_id)
);
```

#### 1.3 Заполнить начальные данные

```sql
-- migrations/000168_seed_category_keywords.up.sql

-- Шины и колеса (1304)
INSERT INTO category_keywords (category_id, keyword, language, weight) VALUES
-- Русский
(1304, 'шина', 'ru', 10.0),
(1304, 'резина', 'ru', 9.0),
(1304, 'колесо', 'ru', 8.0),
(1304, 'покрышка', 'ru', 8.0),
(1304, 'диск', 'ru', 7.0),
(1304, 'летняя', 'ru', 5.0),
(1304, 'зимняя', 'ru', 5.0),
(1304, 'всесезонная', 'ru', 5.0),
-- English
(1304, 'tire', 'en', 10.0),
(1304, 'wheel', 'en', 8.0),
(1304, 'rim', 'en', 7.0),
(1304, 'summer', 'en', 5.0),
(1304, 'winter', 'en', 5.0),
(1304, 'all-season', 'en', 5.0),
-- Бренды
(1304, 'michelin', 'en', 3.0),
(1304, 'bridgestone', 'en', 3.0),
(1304, 'continental', 'en', 3.0),
-- Размеры (паттерны)
(1304, '205/55', 'en', 2.0),
(1304, 'R16', 'en', 2.0),
(1304, 'R17', 'en', 2.0);

-- Аналогично для других категорий...
```

### Этап 2: Backend API

#### 2.1 Новые структуры данных

```go
// internal/domain/models/category_search.go
package models

type CategorySearchHints struct {
    Domain      string            `json:"domain"`
    ProductType string            `json:"productType"`
    Keywords    []string          `json:"keywords"`
    Attributes  map[string]string `json:"attributes"`
    Language    string            `json:"language"`
}

type CategoryMatch struct {
    Category   MarketplaceCategory `json:"category"`
    Confidence float64            `json:"confidence"`
    MatchedKeywords []string      `json:"matchedKeywords"`
}

type CategorySearchResult struct {
    BestMatch    *CategoryMatch   `json:"bestMatch"`
    Alternatives []CategoryMatch  `json:"alternatives"`
}
```

#### 2.2 Новый сервис для поиска категорий

```go
// internal/proj/marketplace/service/category_matcher.go
package service

type CategoryMatcher interface {
    FindBestCategory(ctx context.Context, hints CategorySearchHints) (*CategorySearchResult, error)
    RecordUserSelection(ctx context.Context, aiSuggested, userSelected int, hints CategorySearchHints) error
    UpdateKeywordWeights(ctx context.Context) error
}

type categoryMatcherService struct {
    db          *sql.DB
    cache       cache.Cache
    openSearch  opensearch.Client
}

func (s *categoryMatcherService) FindBestCategory(ctx context.Context, hints CategorySearchHints) (*CategorySearchResult, error) {
    // 1. Получить все категории из кэша или БД
    categories := s.getCategories(ctx)
    
    // 2. Для каждой категории подсчитать score
    scores := make(map[int]float64)
    matchedKeywords := make(map[int][]string)
    
    for _, category := range categories {
        score, matched := s.calculateCategoryScore(category, hints)
        scores[category.ID] = score
        matchedKeywords[category.ID] = matched
    }
    
    // 3. Отсортировать по score и вернуть топ результаты
    return s.prepareResults(categories, scores, matchedKeywords), nil
}

func (s *categoryMatcherService) calculateCategoryScore(category Category, hints CategorySearchHints) (float64, []string) {
    score := 0.0
    matched := []string{}
    
    // 1. Проверка keywords
    keywords := s.getCategoryKeywords(category.ID, hints.Language)
    for _, hint := range hints.Keywords {
        for _, kw := range keywords {
            if strings.Contains(strings.ToLower(hint), strings.ToLower(kw.Keyword)) {
                score += kw.Weight
                matched = append(matched, kw.Keyword)
            }
        }
    }
    
    // 2. Проверка domain
    if category.Slug == hints.Domain || strings.Contains(category.Slug, hints.Domain) {
        score += 5.0
    }
    
    // 3. Проверка атрибутов
    categoryAttrs := s.getCategoryAttributes(category.ID)
    for attrName, attrValue := range hints.Attributes {
        if hasAttribute(categoryAttrs, attrName) {
            score += 2.0
        }
    }
    
    // 4. Бонус за популярность (на основе статистики)
    score += s.getPopularityBonus(category.ID)
    
    return score, matched
}
```

#### 2.3 HTTP Handler

```go
// internal/proj/marketplace/handler/category_search.go
package handler

// SmartCategorySearch поиск категории по семантическим подсказкам
// @Summary Smart category search
// @Description Finds best matching category based on semantic hints
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param body body models.CategorySearchHints true "Search hints from AI"
// @Success 200 {object} models.CategorySearchResult
// @Failure 400 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/categories/smart-search [post]
func (h *Handler) SmartCategorySearch(c *fiber.Ctx) error {
    var hints models.CategorySearchHints
    if err := c.BodyParser(&hints); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request")
    }
    
    // Язык по умолчанию
    if hints.Language == "" {
        hints.Language = "en"
    }
    
    result, err := h.categoryMatcher.FindBestCategory(c.Context(), hints)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Category search failed")
    }
    
    return utils.SuccessResponse(c, result)
}

// RecordCategorySelection записывает выбор пользователя для обучения
// @Summary Record user category selection
// @Description Records what category user selected vs what AI suggested
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body RecordSelectionRequest true "Selection data"
// @Success 200 {object} utils.SuccessResponseSwag
// @Router /api/v1/marketplace/categories/record-selection [post]
func (h *Handler) RecordCategorySelection(c *fiber.Ctx) error {
    userID := c.Locals("userID").(int)
    
    var req RecordSelectionRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request")
    }
    
    err := h.categoryMatcher.RecordUserSelection(
        c.Context(),
        req.AISuggestedID,
        req.UserSelectedID,
        req.Hints,
    )
    
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to record selection")
    }
    
    return utils.SuccessResponse(c, nil)
}
```

### Этап 3: Frontend изменения

#### 3.1 Обновить AI промпты

```typescript
// frontend/svetu/src/app/api/ai/analyze/prompts.ts

export function getAnalysisPrompt(userLanguage: string): string {
  const prompts: Record<string, string> = {
    ru: `Ты - эксперт по анализу товаров. Проанализируй изображение и извлеки информацию.
    
ВАЖНО: НЕ выбирай конкретную категорию! Вместо этого опиши семантические характеристики:

1. title: Заголовок товара
2. description: Описание для покупателей
3. categoryHints: {
   domain: Общая область (automotive/electronics/fashion/real-estate/home-garden/etc)
   productType: Конкретный тип (tire/phone/shirt/apartment/furniture/etc)
   keywords: Массив важных ключевых слов, описывающих товар
   attributes: Извлеченные характеристики
}
4. price: Примерная цена
5. attributes: Конкретные атрибуты (размер, бренд, цвет и т.д.)
6. Остальные поля как раньше...

Пример categoryHints для шины:
{
  "domain": "automotive",
  "productType": "tire", 
  "keywords": ["шина", "летняя", "michelin", "205/55", "R16", "автомобильная"],
  "attributes": {
    "tire_brand": "michelin",
    "tire_size": "205/55 R16",
    "tire_season": "summer"
  }
}`,
    // Аналогично для других языков...
  };
  
  return prompts[userLanguage] || prompts.ru;
}
```

#### 3.2 Новый сервис для определения категории

```typescript
// frontend/svetu/src/services/ai/categoryMatcher.ts

interface CategorySearchHints {
  domain: string;
  productType: string;
  keywords: string[];
  attributes: Record<string, string>;
  language: string;
}

interface CategoryMatch {
  category: Category;
  confidence: number;
  matchedKeywords: string[];
}

interface CategorySearchResult {
  bestMatch?: CategoryMatch;
  alternatives: CategoryMatch[];
}

export class CategoryMatcherService {
  async findCategory(hints: CategorySearchHints): Promise<CategorySearchResult> {
    const response = await fetch('/api/v1/marketplace/categories/smart-search', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...hints,
        language: getCurrentLanguage()
      })
    });
    
    if (!response.ok) {
      throw new Error('Category search failed');
    }
    
    return response.json();
  }
  
  async recordSelection(
    aiSuggestedId: number,
    userSelectedId: number,
    hints: CategorySearchHints
  ): Promise<void> {
    await fetch('/api/v1/marketplace/categories/record-selection', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({
        aiSuggestedId,
        userSelectedId,
        hints
      })
    });
  }
}
```

#### 3.3 Обновить компонент create-listing-ai

```typescript
// frontend/svetu/src/app/[locale]/create-listing-ai/page.tsx

const handleAIAnalysis = async (images: File[]) => {
  // 1. AI анализ
  const aiData = await claudeAI.analyzeProduct(imageBase64, locale);
  
  // 2. Поиск категории через новый API
  const categoryResult = await categoryMatcher.findCategory(aiData.categoryHints);
  
  // 3. Если уверенность низкая, показать варианты пользователю
  if (!categoryResult.bestMatch || categoryResult.bestMatch.confidence < 0.7) {
    setShowCategoryConfirmation(true);
    setCategorySuggestions(categoryResult.alternatives);
  } else {
    setSelectedCategory(categoryResult.bestMatch.category);
  }
  
  // 4. Загрузить атрибуты для выбранной категории
  const attributes = await loadCategoryAttributes(selectedCategory.id);
};

// Компонент подтверждения категории
const CategoryConfirmation = ({ suggestions, onSelect }) => (
  <div className="alert alert-warning">
    <h4>Пожалуйста, уточните категорию:</h4>
    <div className="grid grid-cols-2 gap-2 mt-4">
      {suggestions.map(match => (
        <button
          key={match.category.id}
          onClick={() => onSelect(match.category)}
          className="btn btn-outline"
        >
          {match.category.name}
          <span className="text-xs opacity-70">
            ({Math.round(match.confidence * 100)}%)
          </span>
        </button>
      ))}
    </div>
  </div>
);
```

### Этап 4: Обучение системы

#### 4.1 Фоновая задача обновления весов

```go
// internal/proj/marketplace/service/keyword_optimizer.go

func (s *categoryMatcherService) UpdateKeywordWeights(ctx context.Context) error {
    // 1. Получить историю выборов за последние 30 дней
    history := s.getSelectionHistory(ctx, 30)
    
    // 2. Для каждой пары (AI предложил -> пользователь выбрал)
    for _, record := range history {
        if record.AISuggestedID != record.UserSelectedID {
            // Пользователь исправил выбор AI
            
            // Увеличить веса keywords для правильной категории
            for _, keyword := range record.Keywords {
                s.increaseKeywordWeight(ctx, record.UserSelectedID, keyword, 0.1)
            }
            
            // Уменьшить веса для неправильной категории
            for _, keyword := range record.Keywords {
                s.decreaseKeywordWeight(ctx, record.AISuggestedID, keyword, 0.05)
            }
        }
    }
    
    // 3. Добавить новые keywords из популярных выборов
    s.extractNewKeywords(ctx, history)
    
    return nil
}
```

#### 4.2 Cron задача

```go
// cmd/api/main.go

// Запускать каждый день в 3 часа ночи
scheduler.Cron("0 3 * * *", func() {
    if err := categoryMatcher.UpdateKeywordWeights(context.Background()); err != nil {
        logger.Error().Err(err).Msg("Failed to update keyword weights")
    }
})
```

### Этап 5: Мониторинг и метрики

#### 5.1 Метрики для отслеживания

1. **Точность определения категории**
   - % случаев когда пользователь согласился с AI
   - Среднее confidence score

2. **Производительность**
   - Время поиска категории
   - Количество обращений к БД

3. **Обучение**
   - Количество новых keywords
   - Изменение весов

#### 5.2 Dashboard в админке

Создать страницу статистики:
- График точности по дням
- Топ неправильно определяемых категорий
- Список keywords с весами для редактирования

## Преимущества решения

1. **Масштабируемость**: Легко добавлять новые категории
2. **Самообучение**: Система улучшается от использования
3. **Мультиязычность**: Keywords на разных языках
4. **Гибкость**: Веса можно настраивать
5. **Производительность**: Кэширование и индексы

## План внедрения

### Фаза 1 (1 неделя)
- [ ] Создать таблицы БД
- [ ] Реализовать базовый поиск по keywords
- [ ] Обновить AI промпты

### Фаза 2 (1 неделя)
- [ ] Реализовать API endpoints
- [ ] Интегрировать в frontend
- [ ] Добавить начальные keywords

### Фаза 3 (3 дня)
- [ ] Добавить запись истории выборов
- [ ] Реализовать обновление весов
- [ ] Создать метрики

### Фаза 4 (2 дня)
- [ ] Тестирование
- [ ] Оптимизация
- [ ] Документация

## Риски и митигация

1. **Холодный старт**: Мало данных в начале
   - Митигация: Качественные начальные keywords

2. **Производительность**: Поиск по многим категориям
   - Митигация: Кэширование, индексы, ограничение глубины

3. **Неправильное обучение**: Плохие данные от пользователей
   - Митигация: Валидация, ограничения на изменения весов