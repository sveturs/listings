# ✅ Исправление перевода категорий в хлебных крошках

## Решенная проблема

На страницах товаров маркетплейса категории в хлебных крошках всегда отображались на сербском языке, независимо от выбранного языка интерфейса.

## Реализованное решение

### 1. Backend - добавлена поддержка языка в эндпоинтах

**Файлы:**
- `/backend/internal/proj/marketplace/handler/listings.go`
- `/backend/internal/proj/marketplace/storage/postgres/marketplace.go`

**Изменения:**

1. В handler'ах `GetListing` и `GetListingBySlug` добавлен параметр языка:
```go
// Получаем язык из query параметра
lang := c.Query("lang", "en")

// Создаем контекст с языком
ctx := context.WithValue(c.Context(), "locale", lang)
```

2. В storage изменен SQL запрос для получения переводов категорий:
```sql
WITH RECURSIVE category_path AS (
    SELECT id, name, slug, parent_id, 1 as level
    FROM marketplace_categories
    WHERE id = $1
    UNION ALL
    SELECT c.id, c.name, c.slug, c.parent_id, cp.level + 1
    FROM marketplace_categories c
    JOIN category_path cp ON c.id = cp.parent_id
)
SELECT
    cp.id,
    COALESCE(t.translated_text, cp.name) as name,
    cp.slug
FROM category_path cp
LEFT JOIN translations t ON
    t.entity_type = 'category' AND
    t.entity_id = cp.id AND
    t.field_name = 'name' AND
    t.language = $2
ORDER BY cp.level DESC
```

### 2. Frontend - добавлена передача языка в API запрос

**Файл:** `/frontend/svetu/src/app/[locale]/marketplace/[id]/page.tsx`

**Изменение:**
```typescript
// Было:
let response = await fetch(
  `${config.getApiUrl()}/api/v1/marketplace/listings/${id}`
);

// Стало:
let response = await fetch(
  `${config.getApiUrl()}/api/v1/marketplace/listings/${id}?lang=${locale}`
);
```

## Результаты тестирования

### API возвращает корректные переводы:

**Английский (`?lang=en`):**
```json
{
  "category_path_names": ["Home & Garden", "Furniture"]
}
```

**Русский (`?lang=ru`):**
```json
{
  "category_path_names": ["Дом и сад", "Мебель"]
}
```

**Сербский (`?lang=sr`):**
```json
{
  "category_path_names": ["Дом и башта", "Намештај"]
}
```

## Проверка работоспособности

1. Откройте страницу товара: `http://localhost:3001/ru/marketplace/320`
2. Переключите язык интерфейса
3. Хлебные крошки должны отображаться на выбранном языке

## Технические детали

- Используется `COALESCE` для fallback на оригинальное название, если перевод отсутствует
- Язык передается через контекст для использования в нижних слоях
- Frontend автоматически передает текущий locale в API запросы