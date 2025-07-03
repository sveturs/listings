# Сводка изменений от 2025-06-29

## Проблема
Эндпоинт `/api/v1/entity/user/{id}/stats` возвращал ошибку 500:
```
ERROR: relation "user_ratings" does not exist (SQLSTATE 42703)
```

## Исправления

### 1. База данных - Миграция 000062

**Файл:** `000062_create_review_stats_views.up.sql`

Создана миграция для материализованных представлений:

- **user_ratings** - статистика отзывов для пользователей
- **storefront_ratings** - статистика отзывов для витрин (с полем owner_id)
- **listing_ratings** - статистика отзывов для объявлений
- Индексы для производительности
- Функция `refresh_rating_views()` для обновления

**Ключевое исправление:** Добавлено недостающее поле `owner_id` в `storefront_ratings`, которое ожидалось в коде.

### 2. Код сервиса - Обработка ошибок

**Файл:** `/data/hostel-booking-system/backend/internal/proj/reviews/service/review.go`
**Функция:** `GetReviewStats` (строки 292-431)

**Изменения:**
```go
// БЫЛО:
if err == sql.ErrNoRows {

// СТАЛО:
if err == pgx.ErrNoRows {
```

**Места изменений:**
- Строки 314-320: обработка ошибок для пользователей
- Строки 358-363: обработка ошибок для витрин

**Причина:** Проект использует драйвер `pgx/v5`, но код проверял ошибки из пакета `database/sql`.

### 3. Frontend - Структура ответа API

**Файл:** `/data/hostel-booking-system/frontend/svetu/src/app/[locale]/profile/listings/page.tsx`

**Проблема:** `fetchMyListings` логировал "Unexpected response structure"

**Исправления:**
1. Обновлен интерфейс `ListingsResponse`:
```typescript
interface ListingsResponse {
  success: boolean;  // Добавлено поле
  data: UserListing[];
  meta: { total: number; page: number; limit: number; };
}
```

2. Упрощена логика обработки:
```typescript
// БЫЛО:
if (response.data.data && Array.isArray(response.data.data)) {
  setListings(response.data.data);
} else if (Array.isArray(response.data)) {
  setListings(response.data);
} else {
  console.error('Unexpected response structure:', response);

// СТАЛО:
if (!response.error && response.data && response.data.success) {
  if (Array.isArray(response.data.data)) {
    setListings(response.data.data);
  } else {
    console.error('Unexpected response structure:', response);
```

## Результат

✅ **Эндпоинт `/api/v1/entity/user/{id}/stats` работает корректно**
```json
{
  "data": {
    "total_reviews": 0,
    "average_rating": 0,
    "verified_reviews": 0,
    "rating_distribution": {},
    "photo_reviews": 0
  },
  "success": true
}
```

✅ **Ошибки переводов исчезли** (переводы уже были в ru.json)

✅ **fetchMyListings корректно обрабатывает ответы API**

## Команды для применения миграции

```bash
cd /data/hostel-booking-system/backend
# Проверить текущую версию
./migrate version

# Применить новую миграцию
./migrate up

# Или конкретную миграцию
./migrate to 62
```

## Проверка результата

```bash
# Проверить создание представлений
psql -d your_db -c "\d+ user_ratings"
psql -d your_db -c "\d+ storefront_ratings"

# Проверить работу эндпоинта
curl http://localhost:3000/api/v1/entity/user/2/stats
```

---
**Автор:** Claude  
**Дата:** 2025-06-29 19:40  
**Тип:** Исправление критических ошибок