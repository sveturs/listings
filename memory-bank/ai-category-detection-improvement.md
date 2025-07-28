# Улучшение AI определения категорий

## Предлагаемая архитектура

### 1. Изменение структуры AI анализа

```typescript
interface AIAnalysisResult {
  // Базовая информация
  title: string;
  description: string;
  
  // Новое: семантическая информация для поиска категории
  categoryHints: {
    domain: string;           // "automotive", "electronics", etc.
    productType: string;      // "tire", "phone", "furniture"
    keywords: string[];       // ["summer", "michelin", "205/55"]
    attributes: Record<string, any>; // Извлеченные характеристики
  };
  
  // Остальные поля...
}
```

### 2. Создание таблицы category_keywords в БД

```sql
CREATE TABLE category_keywords (
  id SERIAL PRIMARY KEY,
  category_id INTEGER REFERENCES marketplace_categories(id),
  keyword VARCHAR(100) NOT NULL,
  weight FLOAT DEFAULT 1.0,
  language VARCHAR(2) DEFAULT 'en',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_keyword (keyword),
  INDEX idx_category_weight (category_id, weight DESC)
);

-- Примеры данных
INSERT INTO category_keywords (category_id, keyword, weight, language) VALUES
(1304, 'tire', 10.0, 'en'),
(1304, 'шина', 10.0, 'ru'),
(1304, 'guma', 10.0, 'sr'),
(1304, 'wheel', 8.0, 'en'),
(1304, 'michelin', 5.0, 'en'),
(1304, 'summer', 3.0, 'en'),
(1304, '205/55', 2.0, 'en');
```

### 3. API endpoint для умного поиска категории

```go
// backend/internal/proj/marketplace/handler/category_search.go

func (h *Handler) SmartCategorySearch(c *fiber.Ctx) error {
    type Request struct {
        Domain      string              `json:"domain"`
        ProductType string              `json:"productType"`
        Keywords    []string            `json:"keywords"`
        Attributes  map[string]string   `json:"attributes"`
        Language    string              `json:"language"`
    }
    
    var req Request
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    // Поиск категории по весам
    category, confidence := h.service.FindBestCategory(
        req.Domain,
        req.ProductType,
        req.Keywords,
        req.Attributes,
        req.Language,
    )
    
    return c.JSON(fiber.Map{
        "category": category,
        "confidence": confidence,
        "alternatives": h.service.GetAlternativeCategories(req),
    })
}
```

### 4. Упрощенный промпт для AI

```typescript
const SIMPLIFIED_PROMPT = `
Analyze the image and extract:

1. domain: General area (automotive/electronics/fashion/real-estate/etc)
2. productType: Specific type (tire/phone/shirt/apartment/etc)  
3. keywords: Important terms that describe the product
4. attributes: Key characteristics you can identify

Don't worry about exact category - just describe what you see!
`;
```

### 5. Frontend интеграция

```typescript
// services/ai/categoryMatcher.ts
export class CategoryMatcher {
  async findCategory(aiAnalysis: any): Promise<CategoryMatch> {
    // Сначала пробуем найти по ключевым словам
    const response = await fetch('/api/v1/marketplace/smart-category-search', {
      method: 'POST',
      body: JSON.stringify({
        domain: aiAnalysis.categoryHints.domain,
        productType: aiAnalysis.categoryHints.productType,
        keywords: aiAnalysis.categoryHints.keywords,
        attributes: aiAnalysis.attributes,
        language: getCurrentLanguage()
      })
    });
    
    const result = await response.json();
    
    // Если уверенность низкая, показываем пользователю варианты
    if (result.confidence < 0.7) {
      return {
        category: result.category,
        needsConfirmation: true,
        alternatives: result.alternatives
      };
    }
    
    return { category: result.category, needsConfirmation: false };
  }
}
```

## Преимущества подхода:

1. **Масштабируемость**: Легко добавлять новые категории через keywords
2. **Многоязычность**: Ключевые слова на разных языках
3. **Гибкость**: Веса позволяют настроить приоритеты
4. **Производительность**: Быстрый поиск по индексам
5. **Обучаемость**: Можно улучшать веса на основе выборов пользователей

## План внедрения:

1. Создать таблицу category_keywords
2. Заполнить начальными данными для существующих категорий
3. Создать API endpoint для поиска
4. Обновить AI промпты для извлечения hints
5. Интегрировать в frontend
6. Добавить сбор статистики для улучшения весов

## Дополнительные улучшения:

1. **ML модель**: Обучить небольшую модель на истории выборов категорий
2. **Elasticsearch**: Использовать для более умного поиска по keywords
3. **Кэширование**: Кэшировать популярные соответствия
4. **A/B тестирование**: Сравнивать разные алгоритмы matching'а