# Паспорт таблицы: marketplace_chats

## Назначение
Управление чатами между покупателями и продавцами. Поддерживает как чаты по конкретному объявлению, так и прямые чаты между пользователями.

## Структура таблицы

```sql
CREATE TABLE marketplace_chats (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    buyer_id INT REFERENCES users(id),
    seller_id INT REFERENCES users(id),
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN DEFAULT false,
    UNIQUE(listing_id, buyer_id, seller_id)
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор чата (SERIAL)
- `listing_id` - связанное объявление (может быть NULL для прямых чатов)
- `buyer_id` - ID покупателя/инициатора чата (FK к users)
- `seller_id` - ID продавца/получателя (FK к users)

### Статус и время
- `last_message_at` - время последнего сообщения
- `is_archived` - архивирован ли чат (по умолчанию false)

### Системные поля
- `created_at` - дата создания чата
- `updated_at` - дата последнего обновления

### Уникальные ограничения
- `UNIQUE(listing_id, buyer_id, seller_id)` - один чат на комбинацию

## Индексы

### Базовые индексы
1. **idx_marketplace_chats_buyer** - поиск чатов покупателя
2. **idx_marketplace_chats_seller** - поиск чатов продавца
3. **idx_marketplace_chats_updated** - сортировка по активности

### Оптимизационные индексы
4. **idx_marketplace_chats_user_lookup** - быстрый поиск чатов пользователя
5. **idx_marketplace_chats_archived** - фильтр активных чатов
6. **idx_marketplace_chats_listing** - чаты по объявлению
7. **idx_marketplace_chats_participants** - поиск по участникам
8. **idx_marketplace_chats_active_sorted** - активные чаты по времени
9. **idx_marketplace_chats_listing_participants** - составной индекс
10. **idx_unique_direct_chat** - уникальность прямых чатов

## Триггеры

- **update_marketplace_chats_timestamp** - автоматическое обновление updated_at

## Связи с другими таблицами

### Прямые связи
- `listing_id` → `marketplace_listings.id` (CASCADE DELETE)
- `buyer_id` → `users.id`
- `seller_id` → `users.id`

### Обратные связи
- `marketplace_messages.chat_id` - сообщения в чате

## Бизнес-правила

### Типы чатов
1. **Чат по объявлению** - listing_id NOT NULL
   - Один чат на тройку (listing, buyer, seller)
   - Удаляется при удалении объявления

2. **Прямой чат** - listing_id IS NULL
   - Один чат между двумя пользователями
   - Уникальность по паре участников

### Архивация
- `is_archived = true` - чат скрыт из основного списка
- Архивные чаты доступны через отдельный интерфейс
- При новом сообщении чат автоматически разархивируется

### Обновление времени
- `last_message_at` обновляется при каждом новом сообщении
- Используется для сортировки чатов по активности

## Примеры использования

### Создание чата по объявлению
```sql
INSERT INTO marketplace_chats (listing_id, buyer_id, seller_id)
VALUES (123, 456, 789)
ON CONFLICT (listing_id, buyer_id, seller_id) 
DO UPDATE SET is_archived = false;
```

### Создание прямого чата
```sql
INSERT INTO marketplace_chats (buyer_id, seller_id)
VALUES (LEAST(456, 789), GREATEST(456, 789))
ON CONFLICT ON CONSTRAINT idx_unique_direct_chat
DO UPDATE SET is_archived = false;
```

### Получение активных чатов пользователя
```sql
SELECT c.*, 
       l.title as listing_title,
       u1.name as buyer_name,
       u2.name as seller_name,
       (SELECT COUNT(*) FROM marketplace_messages m 
        WHERE m.chat_id = c.id 
        AND m.receiver_id = :user_id 
        AND NOT m.is_read) as unread_count
FROM marketplace_chats c
LEFT JOIN marketplace_listings l ON c.listing_id = l.id
JOIN users u1 ON c.buyer_id = u1.id
JOIN users u2 ON c.seller_id = u2.id
WHERE (c.buyer_id = :user_id OR c.seller_id = :user_id)
  AND NOT c.is_archived
ORDER BY c.last_message_at DESC;
```

### Архивация чата
```sql
UPDATE marketplace_chats 
SET is_archived = true, updated_at = CURRENT_TIMESTAMP
WHERE id = 123;
```

### Обновление времени последнего сообщения
```sql
UPDATE marketplace_chats 
SET last_message_at = CURRENT_TIMESTAMP,
    is_archived = false
WHERE id = 123;
```

## Известные особенности

1. **LEAST/GREATEST для прямых чатов** - обеспечивает уникальность независимо от порядка ID
2. **CASCADE DELETE** - чаты автоматически удаляются с объявлением
3. **Множественные индексы** - оптимизированы для разных сценариев поиска
4. **Автоматическая разархивация** - при новом сообщении
5. **Счетчик непрочитанных** - вычисляется на лету из messages

## API интеграция

### Endpoints
- `GET /api/v1/marketplace/chats` - список чатов пользователя
- `POST /api/v1/marketplace/chats` - создать/открыть чат
- `GET /api/v1/marketplace/chats/{id}` - информация о чате
- `PUT /api/v1/marketplace/chats/{id}/archive` - архивировать
- `DELETE /api/v1/marketplace/chats/{id}` - удалить чат

### WebSocket события
- `chat:new` - новый чат создан
- `chat:updated` - чат обновлен
- `chat:archived` - чат архивирован

## Производительность

1. **Пагинация обязательна** - для списков чатов
2. **Кеширование счетчиков** - unread_count можно кешировать
3. **Индексы покрывают все запросы** - не требуется full scan
4. **Партиционирование** - возможно по created_at для больших объемов

## Миграции

- **000001** - создание таблицы
- **000039** - добавление оптимизационных индексов и уникальных ограничений