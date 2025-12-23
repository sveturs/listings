# Categories V2 API

## Category Detection Service

Унифицированный сервис автоматического определения категорий для товаров.

### Файлы

- `category_detection.proto` - Proto определение сервиса
- `category_detection.pb.go` - Сгенерированные Go типы
- `category_detection_grpc.pb.go` - Сгенерированный gRPC клиент/сервер

### Методы

#### DetectFromText
Определение категории по тексту (название + описание товара).

```go
import categoriesv2 "github.com/vondi-global/listings/api/proto/categories/v2"

resp, err := client.DetectFromText(ctx, &categoriesv2.DetectFromTextRequest{
    Title:       "iPhone 15 Pro Max 256GB",
    Description: "Latest Apple smartphone with A17 Pro chip",
    Language:    "en",
    Hints: &categoriesv2.CategoryHints{
        Domain:      "electronics",
        ProductType: "smartphone",
        Keywords:    []string{"Apple", "iOS", "5G"},
    },
})
```

**Ответ:**
```go
resp.Primary.CategoryId       // "uuid-категории"
resp.Primary.CategoryName     // "Smartphones"
resp.Primary.ConfidenceScore  // 0.95
resp.Primary.DetectionMethod  // "ai_claude"
resp.Alternatives             // [CategoryMatch, ...]
```

#### DetectFromKeywords
Определение категории по списку ключевых слов.

```go
resp, err := client.DetectFromKeywords(ctx, &categoriesv2.DetectFromKeywordsRequest{
    Keywords: []string{"smartphone", "Apple", "iPhone", "256GB"},
    Language: "en",
})
```

#### DetectBatch
Массовое определение категорий (для импорта товаров).

```go
resp, err := client.DetectBatch(ctx, &categoriesv2.DetectBatchRequest{
    Items: []*categoriesv2.DetectFromTextRequest{
        {Title: "iPhone 15", Description: "...", Language: "en"},
        {Title: "Samsung Galaxy S24", Description: "...", Language: "en"},
        {Title: "MacBook Pro", Description: "...", Language: "en"},
    },
})

// resp.Results[0].Primary.CategoryName -> "Smartphones"
// resp.Results[1].Primary.CategoryName -> "Smartphones"
// resp.Results[2].Primary.CategoryName -> "Laptops"
```

#### ConfirmSelection
Подтверждение выбора пользователя (для машинного обучения).

```go
_, err := client.ConfirmSelection(ctx, &categoriesv2.ConfirmSelectionRequest{
    DetectionId:        resp.DetectionId,
    SelectedCategoryId: "user-selected-category-uuid",
})
```

### Методы определения (detection_method)

- `ai_claude` - Определение через Claude AI (высокая точность)
- `keyword_match` - Определение по ключевым словам
- `similarity` - Определение по семантической близости

### Уровни уверенности (confidence_score)

- `0.9-1.0` - Очень высокая уверенность
- `0.7-0.9` - Высокая уверенность
- `0.5-0.7` - Средняя уверенность (рекомендуется показать альтернативы)
- `< 0.5` - Низкая уверенность (требуется ручной выбор)

### Примеры использования

#### 1. Автоматическое определение при создании объявления

```go
// Продавец вводит название и описание
resp, err := client.DetectFromText(ctx, &categoriesv2.DetectFromTextRequest{
    Title:       userInput.Title,
    Description: userInput.Description,
    Language:    userInput.Language,
})

if resp.Primary.ConfidenceScore >= 0.7 {
    // Автоматически применяем категорию
    listing.CategoryID = resp.Primary.CategoryId
} else {
    // Показываем пользователю варианты для выбора
    suggestedCategories := append([]CategoryMatch{resp.Primary}, resp.Alternatives...)
}
```

#### 2. Импорт товаров из CSV

```go
items := make([]*categoriesv2.DetectFromTextRequest, len(csvRows))
for i, row := range csvRows {
    items[i] = &categoriesv2.DetectFromTextRequest{
        Title:       row.Title,
        Description: row.Description,
        Language:    "en",
    }
}

resp, err := client.DetectBatch(ctx, &categoriesv2.DetectBatchRequest{
    Items: items,
})

for i, result := range resp.Results {
    csvRows[i].CategoryID = result.Primary.CategoryId
}
```

#### 3. Переклассификация существующих товаров

```go
// Получаем все товары без категории
listings := getListingsWithoutCategory()

for _, listing := range listings {
    resp, err := client.DetectFromText(ctx, &categoriesv2.DetectFromTextRequest{
        Title:       listing.Title,
        Description: listing.Description,
        Language:    listing.Language,
    })

    if err != nil {
        continue
    }

    // Применяем только если уверенность > 80%
    if resp.Primary.ConfidenceScore >= 0.8 {
        updateListingCategory(listing.ID, resp.Primary.CategoryId)
    }
}
```

### Интеграция с Frontend

```typescript
// API endpoint в BFF (Next.js)
async function detectCategory(title: string, description: string) {
  const resp = await grpcClient.detectFromText({
    title,
    description,
    language: i18n.language,
  });

  return {
    primary: resp.primary,
    alternatives: resp.alternatives,
  };
}

// React component
const CategorySelector = () => {
  const [suggestions, setSuggestions] = useState([]);

  useEffect(() => {
    if (title && description) {
      detectCategory(title, description).then(setSuggestions);
    }
  }, [title, description]);

  return (
    <div>
      <h3>Suggested Categories</h3>
      {suggestions.map(cat => (
        <CategoryCard
          key={cat.categoryId}
          name={cat.categoryName}
          confidence={cat.confidenceScore}
          onClick={() => selectCategory(cat.categoryId)}
        />
      ))}
    </div>
  );
};
```

### Генерация кода

```bash
cd /p/github.com/vondi-global/listings
make proto
```

### Связанные файлы

- `categories_v2.proto` - Основной CRUD сервис для категорий
- `category_detection.proto` - Сервис определения категорий (этот файл)
