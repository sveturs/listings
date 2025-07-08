# Field Naming Conventions - Sve Tu Platform

## Соглашения об именовании полей

В проекте Sve Tu Platform используются следующие соглашения об именовании полей в разных слоях системы:

### 1. SQL (PostgreSQL)
- **Стиль**: snake_case
- **Примеры**: `user_id`, `category_id`, `storefront_id`, `created_at`, `updated_at`
- **Правила**:
  - Все буквы строчные
  - Слова разделяются подчеркиванием
  - Имена таблиц во множественном числе: `users`, `marketplace_listings`
  - Внешние ключи: `<table>_id` (например, `user_id`, `category_id`)

### 2. Go структуры
- **Стиль**: CamelCase (PascalCase)
- **Примеры**: `UserID`, `CategoryID`, `StorefrontID`, `CreatedAt`, `UpdatedAt`
- **Правила**:
  - Первая буква заглавная
  - Аббревиатуры пишутся заглавными: `ID`, `URL`, `API`
  - Приватные поля начинаются со строчной буквы

### 3. Go JSON теги
- **Стиль**: snake_case
- **Примеры**: `json:"user_id"`, `json:"category_id"`, `json:"created_at"`
- **Правила**:
  - Соответствуют именам полей в БД
  - Используются для сериализации/десериализации JSON
  - Опциональные поля помечаются: `json:"field,omitempty"`

### 4. TypeScript/JavaScript
- **Переменные и свойства**: camelCase
- **Интерфейсы API**: snake_case (для совместимости с backend)
- **Примеры**:
  ```typescript
  // Локальные переменные
  const userId = 123;
  const categoryName = "Electronics";
  
  // API интерфейсы
  interface MarketplaceListing {
    id: number;
    user_id: number;      // snake_case из API
    category_id: number;  // snake_case из API
    created_at: string;   // snake_case из API
  }
  ```

### 5. OpenSearch
- **Стиль**: snake_case
- **Примеры**: `user_id`, `category_id`, `storefront_id`
- **Правила**:
  - Соответствуют именам полей в БД
  - Используются для полнотекстового поиска
  - Текстовые поля могут иметь подполя: `title.keyword`, `title.search`

## Специальные случаи и исключения

### 1. Поля адреса
- **Проблема**: В БД используются префиксы `address_city`, `address_country`
- **В других слоях**: `city`, `country`
- **Решение**: Использовать алиасы в SQL-запросах или миграцию БД

### 2. Вычисляемые поля
Поля, которые не хранятся в БД, но вычисляются на лету:
- `average_rating` - среднее значение из таблицы `reviews`
- `review_count` - количество отзывов
- `is_favorite` - из таблицы `marketplace_favorites`
- `helpful_votes`, `not_helpful_votes` - из таблицы голосов

### 3. Метаданные
- **БД**: `metadata` (jsonb)
- **Go**: `Metadata map[string]interface{}`
- **TypeScript**: `metadata: { [key: string]: unknown }`

## Рекомендации

1. **Консистентность**: Всегда следуйте установленным соглашениям в соответствующем слое
2. **Маппинг**: При передаче данных между слоями используйте соответствующие преобразования
3. **Документация**: Документируйте любые отклонения от стандартных соглашений
4. **Валидация**: Используйте линтеры и форматтеры для автоматической проверки

## Инструменты для проверки

- **SQL**: pgFormatter
- **Go**: gofmt, golangci-lint
- **TypeScript**: ESLint, Prettier
- **Общее**: pre-commit hooks для автоматической проверки