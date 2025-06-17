# Структура системы чатов Backend

## 1. Структура таблиц БД

### marketplace_chats
```sql
CREATE TABLE marketplace_chats (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,  -- может быть NULL для прямых сообщений
    buyer_id INT REFERENCES users(id),
    seller_id INT REFERENCES users(id),
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN DEFAULT false,
    UNIQUE(listing_id, buyer_id, seller_id)
);
```

### marketplace_messages
```sql
CREATE TABLE marketplace_messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES marketplace_chats(id) ON DELETE CASCADE,
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,  -- может быть NULL
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

### chat_attachments
```sql
CREATE TABLE chat_attachments (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES marketplace_messages(id) ON DELETE CASCADE,
    file_type VARCHAR(20) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    storage_type VARCHAR(50) NOT NULL DEFAULT 'minio',
    storage_bucket VARCHAR(100),
    public_url TEXT,
    thumbnail_url TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 2. Связь с объявлениями

- **Чат с объявлением**: `listing_id` содержит ID объявления
- **Прямой чат**: `listing_id IS NULL` - прямое сообщение между пользователями
- **Удаленное объявление**: Чат сохраняется даже после удаления объявления, но `listing_id` остается для истории

## 3. Ключевые файлы системы

### Сервис чатов
- `/backend/internal/proj/marketplace/service/chat.go` - основная бизнес-логика
- `/backend/internal/proj/marketplace/service/chat_interface.go` - интерфейсы
- `/backend/internal/proj/marketplace/service/chat_attachment.go` - работа с вложениями

### Хранилище (PostgreSQL)
- `/backend/internal/proj/marketplace/storage/postgres/chat.go` - работа с БД
- `/backend/internal/storage/postgres/chat_attachments.go` - вложения

### HTTP обработчики
- `/backend/internal/proj/marketplace/handler/chat.go` - REST API endpoints

### Доменные модели
- `/backend/internal/domain/models/marketplace_chat.go` - структуры данных
- `/backend/internal/domain/models/chat_attachment.go` - модели вложений

## 4. Основные методы сервиса

### ChatService
```go
// Сообщения
SendMessage(ctx, msg) error
GetMessages(ctx, listingID, userID, offset, limit) ([]Message, error)
MarkMessagesAsRead(ctx, messageIDs, userID) error
GetMessageByID(ctx, messageID) (*Message, error)

// Чаты
GetChats(ctx, userID) ([]Chat, error)
GetChat(ctx, chatID, userID) (*Chat, error)
ArchiveChat(ctx, chatID, userID) error
GetUnreadMessagesCount(ctx, userID) (int, error)

// WebSocket
BroadcastMessage(msg)
SubscribeToMessages(userID) chan *Message
UnsubscribeFromMessages(userID)

// Статус онлайн
SetUserOnline(userID)
SetUserOffline(userID)
IsUserOnline(userID) bool
GetOnlineUsers() []int
```

## 5. Проверка активности пользователей

### SQL запросы для анализа (см. `/memory-bank/analysis/chat-activity-queries.sql`):

1. **Общая статистика по чатам**
2. **Количество сообщений по пользователю**
3. **Топ активных чатов**
4. **Непрочитанные сообщения**
5. **Активность по дням**
6. **Чаты с удаленными объявлениями**
7. **Детальная статистика пользователя**
8. **Информация по конкретному чату**

## 6. Особенности реализации

### Безопасность
- Санитизация контента сообщений (защита от XSS)
- Проверка прав доступа к чату
- Ограничение длины сообщения (max 10000 символов)

### Уведомления
- Асинхронная отправка через NotificationService
- Поддержка Telegram уведомлений
- Не блокирует основной поток отправки

### WebSocket
- Real-time доставка сообщений
- Отслеживание онлайн статуса
- Подписка на обновления статусов

### Мультиязычность
- Поле `original_language` для исходного языка
- Интеграция с сервисом переводов
- Frontend переводит placeholder'ы

## 7. API Endpoints

- `GET /api/v1/marketplace/chat` - список чатов пользователя
- `GET /api/v1/marketplace/chat/messages` - сообщения чата
- `POST /api/v1/marketplace/messages` - отправка сообщения
- `POST /api/v1/marketplace/messages/read` - отметка как прочитанное
- `GET /api/v1/marketplace/messages/unread` - счетчик непрочитанных
- `POST /api/v1/marketplace/chat/:id/archive` - архивация чата
- `WS /api/v1/marketplace/ws` - WebSocket подключение