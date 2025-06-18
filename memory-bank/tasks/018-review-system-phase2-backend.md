# Улучшение системы отзывов - Фаза 2 (Backend)

## Дата: 17.06.2025
## Ветка: rev_17_06

## Выполненные задачи:

### 1. Автоматическое заполнение entity_origin_type/id

**Файлы изменены:**
- `/backend/internal/proj/reviews/service/review.go` - добавлен метод `setEntityOrigin`

**Логика:**
- При создании отзыва автоматически определяется владелец
- Для товаров: если есть storefront_id → origin = магазин, иначе origin = продавец
- Для магазинов и пользователей: origin = сама сущность

### 2. Материализованные представления для агрегации рейтингов

**Миграция:** `000047_add_rating_aggregation_views.up.sql`

**Созданы представления:**
- `user_ratings` - агрегированные рейтинги пользователей
- `storefront_ratings` - агрегированные рейтинги магазинов

**Функционал:**
- Автоматическое обновление через триггеры
- Распределение оценок (1-5 звезд)
- Тренды за последние 30 дней
- Разбивка по источникам (прямые/через товары/через магазины)

### 3. API endpoints для агрегированных рейтингов

**Новые endpoints:**
```
GET /api/v1/users/{id}/aggregated-rating
GET /api/v1/storefronts/{id}/aggregated-rating
GET /api/v1/reviews/can-review/{type}/{id}
POST /api/v1/reviews/{id}/confirm
POST /api/v1/reviews/{id}/dispute
```

### 4. Новые модели данных

**Файл:** `/backend/internal/domain/models/rating.go`

```go
type AggregatedRating struct {
    Average             float64
    TotalReviews        int
    Distribution        map[int]int
    Breakdown           RatingBreakdown
    VerifiedPercentage  int
    RecentTrend         string // up, down, stable
}

type CanReviewResponse struct {
    CanReview         bool
    Reason            string
    HasExistingReview bool
    ExistingReviewID  *int
}
```

### 5. Реализованные методы Storage

**Файл:** `/backend/internal/storage/postgres/db.go`

- `GetUserAggregatedRating` - получение агрегированного рейтинга пользователя
- `GetStorefrontAggregatedRating` - получение агрегированного рейтинга магазина
- `CreateReviewConfirmation` - создание подтверждения отзыва
- `CreateReviewDispute` - создание спора по отзыву
- `CanUserReviewEntity` - проверка возможности оставить отзыв

### 6. Бизнес-логика в сервисе

**Новые методы ReviewService:**
- `GetUserAggregatedRating` - с расчетом трендов и процентов
- `GetStorefrontAggregatedRating` - аналогично для магазинов
- `ConfirmReview` - подтверждение отзыва продавцом
- `DisputeReview` - создание спора по отзыву
- `CanUserReviewEntity` - проверка прав на отзыв

### 7. Проверки и ограничения

**Реализованы проверки:**
- Один пользователь = один отзыв на объект (уникальный индекс)
- Нельзя оставить отзыв на свой товар
- Нельзя создать несколько споров на один отзыв
- Автоматическое определение тренда рейтинга

## Технические детали:

### Алгоритм определения тренда:
```go
if recentRating - averageRating > 0.2 {
    trend = "up"
} else if recentRating - averageRating < -0.2 {
    trend = "down"
} else {
    trend = "stable"
}
```

### Структура агрегации:
- Отзывы на товары → рейтинг продавца/магазина
- Отзывы на магазин → рейтинг владельца
- Сохранение связи через entity_origin_type/id

## Результаты:

- ✅ Миграции успешно применены
- ✅ Материализованные представления созданы
- ✅ API endpoints работают
- ✅ Swagger документация обновлена
- ✅ Готова инфраструктура для frontend

## Следующие шаги (Frontend):

1. Создать компоненты для отображения рейтингов
2. Интегрировать API в страницу товара
3. Создать страницы продавцов и магазинов
4. Реализовать формы для отзывов
5. Добавить модальные окна для споров и подтверждений