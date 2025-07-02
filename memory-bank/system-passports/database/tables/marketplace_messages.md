# Паспорт таблицы: marketplace_messages

## Назначение
Хранение сообщений в чатах между пользователями. Поддерживает текстовые сообщения, вложения и отслеживание статуса прочтения.

## Структура таблицы

```sql
CREATE TABLE marketplace_messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES marketplace_chats(id) ON DELETE CASCADE,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    receiver_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    original_language VARCHAR(2) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    has_attachments BOOLEAN DEFAULT false,
    attachments_count INT DEFAULT 0
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор сообщения (SERIAL)
- `chat_id` - привязка к чату (FK с CASCADE DELETE)
- `listing_id` - связанное объявление (FK с CASCADE DELETE)
- `sender_id` - отправитель сообщения (FK к users)
- `receiver_id` - получатель сообщения (FK к users)

### Содержимое
- `content` - текст сообщения (обязательное поле)
- `original_language` - язык оригинала (по умолчанию 'en')

### Статус
- `is_read` - прочитано ли сообщение (по умолчанию false)

### Вложения
- `has_attachments` - есть ли вложения (по умолчанию false)
- `attachments_count` - количество вложений (по умолчанию 0)

### Системные поля
- `created_at` - время отправки сообщения
- `updated_at` - время последнего изменения

## Индексы

### Базовые индексы
1. **idx_marketplace_messages_chat** - сообщения в чате
2. **idx_marketplace_messages_listing** - сообщения по объявлению
3. **idx_marketplace_messages_sender** - отправленные сообщения
4. **idx_marketplace_messages_receiver** - полученные сообщения
5. **idx_marketplace_messages_created** - сортировка по времени

### Оптимизационные индексы
6. **idx_marketplace_messages_unread** - непрочитанные сообщения
7. **idx_marketplace_messages_chat_ordered** - сообщения в чате по времени
8. **idx_marketplace_messages_chat_unread** - непрочитанные в чате
9. **idx_marketplace_messages_receiver_unread_count** - для подсчета непрочитанных
10. **idx_marketplace_messages_chat_last** - последнее сообщение в чате

## Триггеры

- **update_marketplace_messages_timestamp** - автоматическое обновление updated_at

## Связи с другими таблицами

### Прямые связи
- `chat_id` → `marketplace_chats.id` (CASCADE DELETE)
- `listing_id` → `marketplace_listings.id` (CASCADE DELETE)
- `sender_id` → `users.id`
- `receiver_id` → `users.id`

### Обратные связи
- `chat_attachments` - вложения сообщения (если has_attachments = true)

## Бизнес-правила

### Отправка сообщений
1. **Участники чата** - sender_id должен быть buyer_id или seller_id в чате
2. **Получатель** - автоматически определяется как другой участник чата
3. **Обязательный контент** - content не может быть пустым

### Статус прочтения
- `is_read = false` - новое сообщение
- `is_read = true` - прочитанное сообщение
- Обновляется при открытии чата получателем

### Вложения
- Файлы хранятся отдельно в `chat_attachments`
- `has_attachments` флаг для быстрой фильтрации
- `attachments_count` для отображения количества

### Языки
- `original_language` - двухбуквенный код ISO 639-1
- Используется для автоперевода сообщений

## Примеры использования

### Отправка сообщения
```sql
-- Сначала проверяем участников чата
WITH chat_info AS (
    SELECT buyer_id, seller_id 
    FROM marketplace_chats 
    WHERE id = :chat_id
)
INSERT INTO marketplace_messages (
    chat_id, listing_id, sender_id, receiver_id, content, original_language
)
SELECT 
    :chat_id,
    (SELECT listing_id FROM marketplace_chats WHERE id = :chat_id),
    :sender_id,
    CASE 
        WHEN :sender_id = buyer_id THEN seller_id
        ELSE buyer_id
    END,
    :content,
    :language
FROM chat_info;

-- Обновляем время последнего сообщения в чате
UPDATE marketplace_chats 
SET last_message_at = CURRENT_TIMESTAMP
WHERE id = :chat_id;
```

### Получение сообщений чата
```sql
SELECT m.*, 
       s.name as sender_name,
       s.picture_url as sender_avatar
FROM marketplace_messages m
JOIN users s ON m.sender_id = s.id
WHERE m.chat_id = :chat_id
ORDER BY m.created_at DESC
LIMIT 50 OFFSET :offset;
```

### Отметка сообщений как прочитанных
```sql
UPDATE marketplace_messages
SET is_read = true
WHERE chat_id = :chat_id
  AND receiver_id = :user_id
  AND is_read = false;
```

### Получение количества непрочитанных
```sql
-- Для конкретного чата
SELECT COUNT(*) as unread_count
FROM marketplace_messages
WHERE chat_id = :chat_id
  AND receiver_id = :user_id
  AND is_read = false;

-- Общее количество для пользователя
SELECT COUNT(*) as total_unread
FROM marketplace_messages
WHERE receiver_id = :user_id
  AND is_read = false;
```

### Поиск сообщений
```sql
SELECT m.*, c.listing_id
FROM marketplace_messages m
JOIN marketplace_chats c ON m.chat_id = c.id
WHERE (m.sender_id = :user_id OR m.receiver_id = :user_id)
  AND m.content ILIKE '%' || :search_term || '%'
ORDER BY m.created_at DESC;
```

## Известные особенности

1. **CASCADE DELETE** - сообщения удаляются с чатом и объявлением
2. **Дублирование listing_id** - для производительности запросов
3. **Receiver_id вычисляется** - на основе участников чата
4. **Множественные индексы** - для разных сценариев использования
5. **Original_language** - для будущей функции автоперевода

## API интеграция

### Endpoints
- `GET /api/v1/marketplace/chats/{id}/messages` - сообщения чата
- `POST /api/v1/marketplace/chats/{id}/messages` - отправить сообщение
- `PUT /api/v1/marketplace/messages/{id}/read` - отметить прочитанным
- `POST /api/v1/marketplace/messages/{id}/attachments` - добавить вложение

### WebSocket события
- `message:new` - новое сообщение
- `message:read` - сообщение прочитано
- `message:typing` - индикатор набора текста

### Уведомления
- Push уведомление при новом сообщении
- Email дайджест непрочитанных
- Telegram уведомления (если включены)

## Производительность

1. **Пагинация обязательна** - для истории сообщений
2. **Lazy loading вложений** - загружаются отдельно
3. **Batch update для is_read** - обновление пачками
4. **Кеширование счетчиков** - unread count в Redis

## Миграции

- **000001** - создание таблицы
- **000039** - добавление has_attachments, attachments_count и индексов