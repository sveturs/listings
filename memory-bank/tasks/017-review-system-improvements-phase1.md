# Улучшение системы отзывов - Фаза 1

## Дата: 17.06.2025
## Ветка: rev_17_06

## Выполненные задачи:

### 1. Реализована проверка через чат для верификации покупок

**Файлы изменены:**
- `/backend/internal/proj/reviews/service/review.go` - обновлен метод `checkVerifiedPurchase`
- `/backend/internal/storage/storage.go` - добавлен новый метод `GetChatActivityStats`
- `/backend/internal/storage/postgres/db.go` - реализован метод получения статистики чата
- `/backend/internal/domain/models/marketplace_chat.go` - добавлена структура `ChatActivityStats`

**Логика верификации:**
- Проверяется наличие чата между покупателем и продавцом
- Минимум 5 сообщений от каждой стороны
- Чат не старше 30 дней
- Общее количество сообщений >= 10

### 2. Добавлен уникальный индекс для защиты от спама

**Миграция:** `000046_add_review_constraints_and_tables.up.sql`
- Создан уникальный индекс `idx_reviews_user_entity_unique`
- Один пользователь может оставить только один отзыв на одну сущность
- Индекс работает только для активных отзывов (status != 'deleted')

### 3. Добавлены ограничения на фото

**Файлы изменены:**
- `/backend/internal/proj/reviews/handler/reviews.go` - метод `UploadPhotos`
- `/backend/internal/proj/reviews/service/review.go` - метод `UpdateReviewPhotos`
- `/frontend/svetu/src/messages/{en,ru}.json` - добавлены переводы для новых ошибок

**Ограничения:**
- Максимум 5 фото на отзыв
- Максимальный размер файла: 5MB
- Разрешенные форматы: JPG, PNG, WebP
- Проверка авторства отзыва при загрузке
- Фото добавляются к существующим, а не заменяют их

### 4. Созданы таблицы для подтверждений и споров

**Новые таблицы:**
- `review_confirmations` - подтверждения отзывов продавцами
- `review_disputes` - споры по отзывам
- `review_dispute_messages` - сообщения в спорах

**Добавлены поля в таблицу reviews:**
- `seller_confirmed` - флаг подтверждения продавцом
- `has_active_dispute` - флаг наличия активного спора

## Технические детали:

### SQL запрос для получения статистики чата:
```sql
WITH chat_info AS (
    SELECT c.id as chat_id, c.created_at as chat_created
    FROM marketplace_chats c
    WHERE c.buyer_id = $1 AND c.seller_id = $2 AND c.listing_id = $3
    LIMIT 1
),
message_stats AS (
    SELECT 
        COUNT(*) as total_messages,
        COUNT(*) FILTER (WHERE m.sender_id = $1) as buyer_messages,
        COUNT(*) FILTER (WHERE m.sender_id = $2) as seller_messages,
        MIN(m.created_at) as first_message_date,
        MAX(m.created_at) as last_message_date
    FROM marketplace_messages m
    INNER JOIN chat_info ci ON m.chat_id = ci.chat_id
)
SELECT 
    CASE WHEN ci.chat_id IS NOT NULL THEN true ELSE false END as chat_exists,
    COALESCE(ms.total_messages, 0) as total_messages,
    -- и т.д.
```

### Миграции:
- Успешно применена миграция 000046
- Создан уникальный индекс и новые таблицы
- Добавлены индексы для быстрого поиска

## Следующие шаги (Фаза 2):

1. Реализовать API endpoints для подтверждения отзывов продавцами
2. Создать систему жалоб на отзывы
3. Разработать интерфейс модерации для администратора
4. Добавить уведомления о новых отзывах и спорах

## Тестирование:

- ✅ Миграции успешно применены
- ✅ Уникальный индекс создан
- ✅ Новые таблицы созданы
- ✅ Swagger документация обновлена
- ✅ Добавлены переводы для новых сообщений об ошибках