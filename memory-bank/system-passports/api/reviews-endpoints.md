# Паспорт API Endpoints: Reviews (Отзывы)

## 📋 Метаданные
- **Группа API**: Reviews
- **Базовый путь**: `/api/v1/reviews`
- **Handler**: `backend/internal/proj/reviews/handler/handler.go`
- **Количество endpoints**: 16 (7 публичных, 9 защищенных)
- **Интеграции**: PostgreSQL, MinIO (фото), OpenSearch (индексация)

## 🎯 Назначение
Система отзывов и рейтингов для маркетплейса:
- Двухэтапный процесс создания отзывов (draft → publish)
- Фотографии в отзывах с галереей
- Агрегированные рейтинги для пользователей и витрин
- Система голосования за полезность отзывов
- Ответы продавцов на отзывы
- Модерация и система споров

## 📡 Endpoints

### 🌐 Публичные (без авторизации)

#### GET `/api/v1/reviews/`
**Назначение**: Получение списка отзывов с фильтрацией
- **Handler**: `h.Review.GetReviews`
- **Query Parameters**: 
  - `entity_type`: "user" | "listing" | "storefront"
  - `entity_id`: ID сущности
  - `rating_min`, `rating_max`: фильтр по рейтингу
  - `with_photos`: только отзывы с фотографиями
  - `sort`: "newest" | "oldest" | "rating_high" | "rating_low" | "helpful"
- **Response**: Пагинированный список отзывов

#### GET `/api/v1/reviews/:id`
**Назначение**: Получение детального отзыва
- **Handler**: `h.Review.GetReviewByID`
- **Response**: Полная информация об отзыве + фотографии + ответы

#### GET `/api/v1/reviews/stats`
**Назначение**: Общая статистика отзывов системы
- **Handler**: `h.Review.GetStats`
- **Response**: Агрегированные данные для аналитики

#### GET `/api/v1/entity/:type/:id/rating`
**Назначение**: Средний рейтинг конкретной сущности
- **Handler**: `h.Review.GetEntityRating`
- **Params**: type ("user"|"listing"|"storefront"), entity_id
- **Response**: Средний рейтинг + количество отзывов

#### GET `/api/v1/entity/:type/:id/stats`
**Назначение**: Детальная статистика рейтингов сущности
- **Handler**: `h.Review.GetEntityStats`
- **Response**: Распределение по звездам + тренды

#### GET `/api/v1/users/:id/aggregated-rating`
**Назначение**: Агрегированный рейтинг пользователя (продавца)
- **Handler**: `h.Review.GetUserAggregatedRating`
- **Includes**: Рейтинги как продавца + как владельца витрин
- **Response**: Общий рейтинг + детализация по источникам

#### GET `/api/v1/storefronts/:id/aggregated-rating`
**Назначение**: Агрегированный рейтинг витрины
- **Handler**: `h.Review.GetStorefrontAggregatedRating`
- **Includes**: Рейтинги товаров + обслуживания + доставки
- **Response**: Общий рейтинг + детализация по аспектам

### 🔒 Защищенные (требуют авторизации)

#### GET `/api/v1/reviews/can-review/:type/:id`
**Назначение**: Проверка возможности оставить отзыв
- **Handler**: `h.Review.CanReview`
- **Logic**: Проверяет завершенные транзакции, дубликаты
- **Response**: Boolean + причина если нельзя

#### POST `/api/v1/reviews/draft`
**Назначение**: Создание черновика отзыва (этап 1)
- **Handler**: `h.Review.CreateDraftReview`
- **Body**: DraftReviewRequest
- **Response**: Review ID для добавления фотографий
- **Status**: draft

#### POST `/api/v1/reviews/:id/photos`
**Назначение**: Загрузка фотографий к отзыву (этап 2)
- **Handler**: `h.Review.UploadPhotos`
- **Content-Type**: multipart/form-data
- **Limit**: До 5 фотографий, 10MB каждая
- **Integration**: MinIO bucket "reviews"

#### POST `/api/v1/reviews/:id/publish`
**Назначение**: Публикация отзыва (этап 3)
- **Handler**: `h.Review.PublishReview`
- **Effect**: draft → published, индексация в OpenSearch
- **Notifications**: Уведомление получателю отзыва

#### PUT `/api/v1/reviews/:id`
**Назначение**: Редактирование отзыва
- **Handler**: `h.Review.UpdateReview`
- **Security**: Только автор в течение 24 часов
- **Body**: UpdateReviewRequest (частичное)

#### DELETE `/api/v1/reviews/:id`
**Назначение**: Удаление отзыва
- **Handler**: `h.Review.DeleteReview`
- **Security**: Автор или админ
- **Effect**: Soft delete + пересчет рейтингов

#### POST `/api/v1/reviews/:id/vote`
**Назначение**: Голосование за полезность отзыва
- **Handler**: `h.Review.VoteForReview`
- **Body**: {"vote": "helpful" | "not_helpful"}
- **Logic**: Один голос на пользователя
- **Effect**: Влияет на сортировку отзывов

#### POST `/api/v1/reviews/:id/response`
**Назначение**: Ответ на отзыв (для продавцов)
- **Handler**: `h.Review.AddResponse`
- **Security**: Только получатель отзыва
- **Body**: ResponseRequest с текстом ответа
- **Limit**: Один ответ на отзыв

#### POST `/api/v1/reviews/:id/confirm`
**Назначение**: Подтверждение отзыва покупателем
- **Handler**: `h.Review.ConfirmReview`
- **When**: После получения ответа от продавца
- **Effect**: Помечает отзыв как подтвержденный

#### POST `/api/v1/reviews/:id/dispute`
**Назначение**: Оспаривание отзыва
- **Handler**: `h.Review.DisputeReview`
- **Security**: Получатель отзыва
- **Body**: DisputeRequest с причиной
- **Effect**: Отзыв на модерацию

## 🎭 Структуры данных

### Основная модель отзыва
```typescript
interface Review {
  id: string;
  entity_type: "user" | "listing" | "storefront";
  entity_id: string;
  reviewer_id: string;
  rating: number;                    // 1-5 звезд
  title: string;
  content: string;
  photos: ReviewPhoto[];
  aspects?: AspectRatings;           // детализированные оценки
  status: ReviewStatus;
  helpful_votes: number;
  not_helpful_votes: number;
  response?: ReviewResponse;
  verified_purchase: boolean;
  created_at: string;
  published_at?: string;
  updated_at: string;
}

type ReviewStatus = "draft" | "published" | "hidden" | "disputed" | "deleted";

interface AspectRatings {
  quality?: number;                  // качество товара
  communication?: number;            // общение с продавцом
  delivery?: number;                 // скорость доставки
  description?: number;              // соответствие описанию
  packaging?: number;                // качество упаковки
}

interface ReviewPhoto {
  id: string;
  url: string;
  thumbnail_url: string;
  caption?: string;
  order: number;
}

interface ReviewResponse {
  id: string;
  author_id: string;
  content: string;
  created_at: string;
}
```

### Запросы
```typescript
interface DraftReviewRequest {
  entity_type: "user" | "listing" | "storefront";
  entity_id: string;
  transaction_id?: string;           // для верификации покупки
  rating: number;                    // 1-5
  title: string;                     // до 100 символов
  content: string;                   // до 2000 символов
  aspects?: AspectRatings;
  is_anonymous?: boolean;            // анонимный отзыв
}

interface UpdateReviewRequest {
  rating?: number;
  title?: string;
  content?: string;
  aspects?: AspectRatings;
}

interface DisputeRequest {
  reason: "fake" | "inappropriate" | "spam" | "incorrect" | "other";
  explanation: string;
  evidence_urls?: string[];          // ссылки на доказательства
}
```

### Статистика и аналитика
```typescript
interface EntityRatingStats {
  entity_id: string;
  average_rating: number;            // средний рейтинг
  total_reviews: number;
  rating_distribution: {             // распределение по звездам
    1: number;
    2: number;
    3: number;
    4: number;
    5: number;
  };
  aspects_avg?: AspectRatings;       // средние оценки по аспектам
  trends: {                          // тренды за периоды
    last_7_days: number;
    last_30_days: number;
    last_90_days: number;
  };
  verified_percentage: number;       // процент верифицированных отзывов
}

interface ReviewsFilterStats {
  total_reviews: number;
  with_photos: number;
  verified_purchases: number;
  rating_breakdown: Record<number, number>;
  recent_count: number;              // за последние 30 дней
}
```

## 🔄 Интеграции

### Database Schema
```sql
-- Основная таблица отзывов
reviews (
  id, entity_type, entity_id, reviewer_id,
  rating, title, content, aspects_json,
  status, helpful_votes, not_helpful_votes,
  verified_purchase, created_at, published_at, updated_at
);

-- Фотографии отзывов
review_photos (
  id, review_id, file_path, thumbnail_path,
  caption, order_index, created_at
);

-- Ответы на отзывы
review_responses (
  id, review_id, author_id, content, created_at
);

-- Голоса за полезность
review_votes (
  review_id, user_id, vote_type, created_at,
  PRIMARY KEY (review_id, user_id)
);

-- Споры по отзывам
review_disputes (
  id, review_id, disputer_id, reason,
  explanation, status, admin_response,
  created_at, resolved_at
);
```

### MinIO Integration
- **Bucket**: `reviews`
- **Path**: `/reviews/{review_id}/{photo_id}.{ext}`
- **Thumbnails**: 300x300 для галереи
- **Compression**: JPEG качество 85%

### OpenSearch Integration
- **Index**: `reviews`
- **Mapping**: Поддержка поиска по тексту отзывов
- **Aggregations**: Для статистики и фасетов

## 🎛️ Бизнес-логика

### Верификация покупок
```typescript
function canUserReview(
  userId: string, 
  entityType: string, 
  entityId: string
): ReviewEligibility {
  // Проверяем завершенные транзакции
  const transactions = getCompletedTransactions(userId, entityId);
  
  // Проверяем существующие отзывы
  const existingReviews = getUserReviews(userId, entityType, entityId);
  
  return {
    can_review: transactions.length > 0 && existingReviews.length === 0,
    reason: !transactions.length ? "no_completed_purchase" 
          : existingReviews.length > 0 ? "already_reviewed"
          : "eligible",
    verified_purchase: transactions.length > 0
  };
}
```

### Система модерации
```typescript
interface ModerationRules {
  auto_hide_threshold: -5;           // скрыть если helpful_votes < -5
  dispute_review_threshold: 3;       // на модерацию при 3+ спорах
  spam_detection: {
    min_content_length: 20;
    max_duplicate_percentage: 80;
    blocked_words: string[];
  };
  photo_moderation: {
    max_file_size: 10 * 1024 * 1024; // 10MB
    allowed_types: ["image/jpeg", "image/png"];
    ai_content_check: boolean;
  };
}
```

### Влияние на рейтинги
- Рейтинги пересчитываются в реальном времени
- Веса отзывов: верифицированные покупки × 1.5
- Старые отзывы постепенно теряют вес (decay factor)
- Подозрительные отзывы исключаются из расчета

## 🛡️ Безопасность и модерация

### Защита от накрутки
- Лимит: 1 отзыв на транзакцию
- Проверка связанных аккаунтов (IP, устройство)
- AI детекция сгенерированного контента
- Анализ паттернов активности

### Контроль качества
- Минимальная длина отзыва: 20 символов
- Фильтр нецензурной лексики
- Проверка на спам и дубликаты
- Ручная модерация спорных отзывов

## ⚠️ Известные особенности

### Performance
- Кеширование рейтингов сущностей на 15 минут
- Lazy loading фотографий в списках
- Pagination с курсорами для больших объемов
- Индексы БД по entity_type + entity_id

### UX Features
- Draft система для пошагового создания
- Автосохранение черновиков
- Предпросмотр отзыва перед публикацией
- Push уведомления о новых отзывах

### Analytics
- Трекинг конверсии отзывов (draft → published)
- A/B тесты UI компонентов отзывов
- Анализ корреляции отзывов и продаж
- Sentiment analysis текста отзывов

## 🧪 Примеры использования

### Создание отзыва (3 этапа)
```bash
# Этап 1: Создание черновика
curl -X POST /api/v1/reviews/draft \
  -H "Authorization: Bearer <token>" \
  -d '{
    "entity_type": "listing",
    "entity_id": "listing-123",
    "rating": 5,
    "title": "Отличный товар!",
    "content": "Товар полностью соответствует описанию..."
  }'

# Этап 2: Загрузка фотографий
curl -X POST /api/v1/reviews/review-456/photos \
  -H "Authorization: Bearer <token>" \
  -F "photos=@photo1.jpg" \
  -F "photos=@photo2.jpg"

# Этап 3: Публикация
curl -X POST /api/v1/reviews/review-456/publish \
  -H "Authorization: Bearer <token>"
```

### Получение отзывов с фильтрами
```bash
curl "/api/v1/reviews/?entity_type=user&entity_id=user-123&rating_min=4&with_photos=true&sort=helpful"
```

### Голосование за отзыв
```bash
curl -X POST /api/v1/reviews/review-456/vote \
  -H "Authorization: Bearer <token>" \
  -d '{"vote": "helpful"}'
```