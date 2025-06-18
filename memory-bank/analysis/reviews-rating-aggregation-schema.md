# Схема агрегации рейтингов в системе отзывов

## Принципы агрегации

### 1. Источники оценок

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Отзыв на товар │     │ Отзыв на магазин│     │Отзыв на продавца│
│  entity_type:   │     │  entity_type:   │     │  entity_type:   │
│    listing      │     │   storefront    │     │     user        │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                         │
         │                       │                         │
         ▼                       ▼                         ▼
    ┌────────────────────────────────────────────────────────┐
    │                   АГРЕГАЦИЯ РЕЙТИНГОВ                  │
    └────────────────────────────────────────────────────────┘
```

### 2. Правила агрегации

#### 2.1 Для обычного продавца (без магазина):
```sql
user_rating = AVG(
  -- Прямые отзывы на продавца
  reviews WHERE entity_type = 'user' AND entity_id = user_id
  
  UNION ALL
  
  -- Отзывы на его товары
  reviews WHERE entity_type = 'listing' 
    AND entity_id IN (
      SELECT id FROM marketplace_listings WHERE user_id = user_id
    )
)
```

#### 2.2 Для владельца магазина:
```sql
user_rating = AVG(
  -- Прямые отзывы на продавца
  reviews WHERE entity_type = 'user' AND entity_id = user_id
  
  UNION ALL
  
  -- Отзывы на его магазины
  reviews WHERE entity_type = 'storefront' 
    AND entity_id IN (
      SELECT id FROM storefronts WHERE user_id = user_id
    )
  
  UNION ALL
  
  -- Отзывы на товары его магазинов
  reviews WHERE entity_type = 'listing' 
    AND entity_id IN (
      SELECT l.id FROM marketplace_listings l
      JOIN storefronts s ON l.storefront_id = s.id
      WHERE s.user_id = user_id
    )
)
```

#### 2.3 Для магазина:
```sql
storefront_rating = AVG(
  -- Прямые отзывы на магазин
  reviews WHERE entity_type = 'storefront' AND entity_id = storefront_id
  
  UNION ALL
  
  -- Отзывы на товары магазина
  reviews WHERE entity_type = 'listing' 
    AND entity_id IN (
      SELECT id FROM marketplace_listings 
      WHERE storefront_id = storefront_id
    )
)
```

### 3. Сохранение истории при удалении

#### 3.1 Структура для сохранения связей:
```sql
-- В таблице reviews уже есть поля:
entity_origin_type VARCHAR(50) -- 'user' или 'storefront'  
entity_origin_id INTEGER        -- ID владельца

-- При создании отзыва заполняем:
IF entity_type = 'listing' THEN
  -- Получаем владельца товара
  SELECT user_id, storefront_id INTO @user_id, @storefront_id
  FROM marketplace_listings WHERE id = entity_id;
  
  IF @storefront_id IS NOT NULL THEN
    entity_origin_type = 'storefront';
    entity_origin_id = @storefront_id;
  ELSE
    entity_origin_type = 'user';
    entity_origin_id = @user_id;
  END IF;
  
ELSIF entity_type = 'storefront' THEN
  entity_origin_type = 'storefront';
  entity_origin_id = entity_id;
  
ELSIF entity_type = 'user' THEN
  entity_origin_type = 'user';
  entity_origin_id = entity_id;
END IF;
```

### 4. Материализованные представления для производительности

```sql
-- Обновляется триггерами при изменении отзывов
CREATE MATERIALIZED VIEW user_ratings AS
SELECT 
  u.id as user_id,
  COUNT(DISTINCT r.id) as total_reviews,
  AVG(r.rating) as average_rating,
  COUNT(DISTINCT CASE WHEN r.entity_type = 'user' THEN r.id END) as direct_reviews,
  COUNT(DISTINCT CASE WHEN r.entity_type = 'listing' THEN r.id END) as listing_reviews,
  COUNT(DISTINCT CASE WHEN r.entity_type = 'storefront' THEN r.id END) as storefront_reviews,
  COUNT(DISTINCT CASE WHEN r.is_verified_purchase THEN r.id END) as verified_reviews,
  -- Распределение оценок
  COUNT(CASE WHEN r.rating = 1 THEN 1 END) as rating_1,
  COUNT(CASE WHEN r.rating = 2 THEN 1 END) as rating_2,
  COUNT(CASE WHEN r.rating = 3 THEN 1 END) as rating_3,
  COUNT(CASE WHEN r.rating = 4 THEN 1 END) as rating_4,
  COUNT(CASE WHEN r.rating = 5 THEN 1 END) as rating_5,
  -- Тренд за последние 30 дней
  AVG(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN r.rating END) as recent_rating,
  MAX(r.created_at) as last_review_at
FROM users u
LEFT JOIN reviews r ON (
  -- Прямые отзывы
  (r.entity_type = 'user' AND r.entity_id = u.id) OR
  -- Через origin после удаления
  (r.entity_origin_type = 'user' AND r.entity_origin_id = u.id)
)
WHERE r.status = 'published'
GROUP BY u.id;

-- Аналогично для магазинов
CREATE MATERIALIZED VIEW storefront_ratings AS ...
```

### 5. API endpoints для получения рейтингов

```typescript
// Получить агрегированный рейтинг пользователя
GET /api/v1/users/{id}/rating
Response: {
  average: 4.5,
  total_reviews: 156,
  distribution: {1: 5, 2: 8, 3: 15, 4: 48, 5: 80},
  breakdown: {
    direct: {count: 20, average: 4.6},
    listings: {count: 120, average: 4.4}, 
    storefronts: {count: 16, average: 4.7}
  },
  verified_percentage: 78,
  recent_trend: "stable" // up/down/stable
}

// Получить рейтинг магазина
GET /api/v1/storefronts/{id}/rating

// Получить рейтинг товара
GET /api/v1/listings/{id}/rating
```

### 6. Обновление рейтингов в реальном времени

```go
// При создании/изменении отзыва
func (s *ReviewService) afterReviewChange(ctx context.Context, review *Review) {
  // 1. Обновляем материализованное представление
  s.refreshMaterializedView(review.EntityOriginType, review.EntityOriginID)
  
  // 2. Обновляем кеш Redis
  s.invalidateRatingCache(review.EntityOriginType, review.EntityOriginID)
  
  // 3. Обновляем OpenSearch для товаров
  if review.EntityType == "listing" {
    s.updateListingRatingInSearch(review.EntityID)
  }
  
  // 4. Отправляем событие через WebSocket
  s.publishRatingUpdate(review.EntityOriginType, review.EntityOriginID)
}
```

### 7. Защита целостности данных

```sql
-- Триггер для автоматического заполнения origin
CREATE TRIGGER set_review_origin
BEFORE INSERT ON reviews
FOR EACH ROW
EXECUTE FUNCTION set_review_origin_fields();

-- Триггер для обновления рейтингов
CREATE TRIGGER update_ratings_after_review
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION refresh_rating_views();

-- Периодическая полная пересборка (раз в час)
CREATE OR REPLACE FUNCTION rebuild_all_ratings() ...
```

### 8. Отображение на фронтенде

```typescript
// Компонент для отображения комплексного рейтинга
interface RatingBreakdownProps {
  userId?: number;
  storefrontId?: number;
  showDetails?: boolean;
}

const RatingBreakdown: React.FC<RatingBreakdownProps> = ({...}) => {
  const rating = useAggregatedRating(entityType, entityId);
  
  return (
    <div className="rating-breakdown">
      {/* Общий рейтинг */}
      <div className="overall-rating">
        <span className="rating-value">{rating.average}</span>
        <RatingStars value={rating.average} />
        <span className="review-count">({rating.total_reviews} отзывов)</span>
      </div>
      
      {/* Разбивка по источникам */}
      {showDetails && (
        <div className="rating-sources">
          <div>Как продавец: {rating.breakdown.direct.average}</div>
          <div>За товары: {rating.breakdown.listings.average}</div>
          <div>За магазины: {rating.breakdown.storefronts.average}</div>
        </div>
      )}
      
      {/* График распределения */}
      <RatingDistribution distribution={rating.distribution} />
    </div>
  );
};
```

### 9. Миграция для поддержки агрегации

```sql
-- Добавляем индексы для быстрой агрегации
CREATE INDEX idx_reviews_origin ON reviews(entity_origin_type, entity_origin_id) 
WHERE status = 'published';

CREATE INDEX idx_reviews_aggregation ON reviews(entity_type, entity_id, rating) 
WHERE status = 'published';

-- Создаем таблицу для кеширования рейтингов
CREATE TABLE rating_cache (
  entity_type VARCHAR(50),
  entity_id INTEGER,
  average_rating DECIMAL(3,2),
  total_reviews INTEGER,
  distribution JSONB,
  breakdown JSONB,
  calculated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (entity_type, entity_id)
);
```