# Category Detection Service

## Обзор

CategoryDetectionService - интеллектуальный сервис для автоматического определения категорий товаров на основе названия и описания.

## Архитектура

### Многоуровневая стратегия детекции:

1. **Keyword Matching** (быстро, бесплатно)
   - Поиск по ключевым словам в `meta_keywords` JSONB
   - Confidence: 0.3-1.0

2. **Similarity Search** (PostgreSQL pg_trgm)
   - Поиск по похожести названия категории
   - Confidence: 0.2-1.0

3. **Claude AI** (точно, платно)
   - Используется если keyword/similarity дают низкий confidence (<0.5)
   - Model: `claude-3-5-sonnet-20241022`
   - Confidence: обычно 0.7-0.95

4. **Fallback** (всегда)
   - Возвращает пустой результат если все методы провалились

## Использование

### Инициализация

```go
import (
    "github.com/vondi-global/listings/internal/repository/postgres"
    "github.com/vondi-global/listings/internal/service"
)

// В коде инициализации сервиса
detectionRepo := postgres.NewCategoryDetectionRepository(db, logger)
categoryRepo := postgres.NewRepository(db, logger) // existing repo

detectionService := service.NewCategoryDetectionService(
    detectionRepo,
    categoryRepo,
    redisClient,
    os.Getenv("CLAUDE_API_KEY"), // или из config
    logger,
)
```

### DetectFromText - определение категории по тексту

```go
result, err := detectionService.DetectFromText(ctx, domain.DetectFromTextInput{
    Title:       "iPhone 15 Pro 256GB",
    Description: "Новый смартфон Apple iPhone 15 Pro с 256GB памяти",
    Language:    "ru",
})

if err != nil {
    log.Error("Detection failed", err)
    return err
}

// Результат
if result.Primary != nil {
    fmt.Printf("Категория: %s (уверенность: %.2f)\n", 
        result.Primary.CategoryName, 
        result.Primary.ConfidenceScore)
    fmt.Printf("Метод: %s\n", result.Primary.DetectionMethod)
}

// Альтернативы
for _, alt := range result.Alternatives {
    fmt.Printf("Альтернатива: %s (%.2f)\n", 
        alt.CategoryName, 
        alt.ConfidenceScore)
}
```

### DetectFromKeywords - определение по ключевым словам

```go
result, err := detectionService.DetectFromKeywords(ctx, 
    []string{"iphone", "smartphone", "apple"}, 
    "en",
)
```

### DetectBatch - пакетная обработка

```go
batchResult, err := detectionService.DetectBatch(ctx, domain.DetectBatchInput{
    Items: []domain.DetectFromTextInput{
        {Title: "iPhone 15 Pro", Language: "en"},
        {Title: "Samsung Galaxy S24", Language: "en"},
        {Title: "MacBook Pro M3", Language: "en"},
    },
})

for i, detection := range batchResult.Results {
    fmt.Printf("Item %d: %s\n", i+1, detection.Primary.CategoryName)
}
fmt.Printf("Total time: %dms\n", batchResult.TotalProcessingTime)
```

### ConfirmSelection - подтверждение выбора пользователя

```go
// Когда пользователь подтверждает категорию
err := detectionService.ConfirmSelection(ctx, 
    detectionID,        // UUID детекции
    selectedCategoryID, // UUID выбранной категории
)
```

## Структуры данных

### CategoryDetection

```go
type CategoryDetection struct {
    ID               uuid.UUID       // Уникальный ID детекции
    Primary          *CategoryMatch  // Основная категория
    Alternatives     []CategoryMatch // Альтернативы (до 5)
    ProcessingTimeMs int32           // Время обработки
    
    // Входные данные (для tracking)
    InputTitle       string
    InputDescription string
    InputLanguage    string
    
    // Подтверждение пользователя
    UserConfirmed    *bool
    UserSelectedID   *uuid.UUID
    
    CreatedAt time.Time
}
```

### CategoryMatch

```go
type CategoryMatch struct {
    CategoryID      uuid.UUID
    CategoryName    string
    CategorySlug    string
    CategoryPath    string
    ConfidenceScore float64         // 0.0-1.0
    DetectionMethod DetectionMethod // "ai_claude" | "keyword_match" | "similarity"
    MatchedKeywords []string
}
```

## Требования PostgreSQL

### Расширение pg_trgm (опционально)

Для similarity search требуется расширение pg_trgm:

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

Если расширение не установлено, similarity search будет пропущен (fallback на другие методы).

## Производительность

### Типичное время отклика:

- **Keyword matching:** 10-50ms
- **Similarity search:** 20-100ms
- **Claude AI:** 2-10 секунд (используется редко)
- **Fallback:** <1ms

### Оптимизация:

1. Keyword matching выполняется первым (быстрее всего)
2. Claude AI вызывается только при низком confidence (<0.5)
3. Результаты детекции сохраняются асинхронно (не блокируют ответ)

## Настройка confidence thresholds

В коде сервиса можно настроить пороги:

```go
// В методе DetectFromText
if len(keywordMatches) > 0 && keywordMatches[0].ConfidenceScore > 0.7 {
    // Использовать keyword match без AI
}

if keywordMatches[0].ConfidenceScore < 0.5 {
    // Вызвать Claude AI
}
```

## Логирование

Сервис логирует все этапы детекции:

```
DEBUG: Using keyword match result (confidence=0.85, category=electronics)
WARN: Claude AI detection failed, using fallback
INFO: Category detected successfully (method=ai_claude, confidence=0.92)
```

## Будущие улучшения

1. **Кэширование результатов** (Redis) для популярных запросов
2. **ML модель** вместо Claude API для оффлайн детекции
3. **Tracking таблица** в PostgreSQL для анализа точности
4. **Обучение на feedback** пользователей (ConfirmSelection)
5. **Multi-language support** в промптах Claude

## Безопасность

- ✅ CLAUDE_API_KEY хранится в переменных окружения
- ✅ Таймаут HTTP запросов к Claude API (30 секунд)
- ✅ Graceful fallback при ошибках API
- ✅ Rate limiting (на стороне Claude API)

## Стоимость

- **Keyword matching:** Бесплатно (PostgreSQL query)
- **Similarity search:** Бесплатно (PostgreSQL query)
- **Claude AI:** ~$0.003 за запрос (Sonnet 3.5)
  - Используется в <20% случаев
  - Средняя стоимость: ~$0.0006 на детекцию

## Troubleshooting

### "Claude API key not configured"

```bash
# Установить в .env
VONDILISTINGS_CLAUDE_API_KEY=sk-ant-...
```

### "similarity search failed (pg_trgm extension may not be installed)"

```sql
-- Установить расширение
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

Это не критично - сервис будет работать без similarity search.

### Низкий confidence на всех методах

- Проверить что в БД есть категории с заполненными `meta_keywords`
- Проверить что keywords соответствуют языку input.Language
- Добавить больше keywords в `meta_keywords` для категорий
