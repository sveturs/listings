# WebSocket Architecture для Chat в Listings Microservice

## Текущая ситуация

- ❌ WebSocket в монолите **отключен** (`chat.temporarily_disabled`)
- ❌ Real-time обновления не работают
- ❌ Нужно обновлять страницу для получения новых сообщений
- ❌ Нет уведомлений о новых сообщениях в header

## Цель

Реализовать полноценный real-time чат с WebSocket в Listings микросервисе.

## Архитектура

### Вариант 1: Прямое подключение к микросервису (финальная цель)

```
Browser (WS) → Listings Microservice:8086/ws/chat → Hub → БД listings_dev_db
```

**Проблема:** Нужно пробросить WebSocket через Nginx/reverse proxy

### Вариант 2: Proxy через монолит (текущая реализация)

```
Browser (WS) → Backend:3000/ws/chat → WebSocket Proxy → Listings Microservice:8086/ws/chat
```

**Преимущество:** Единая точка входа, не нужно менять Nginx конфигурацию

## Компоненты

### 1. Hub (Listings Microservice)

Аналогично `tracking/hub.go`:

```go
type ChatHub struct {
    // Подключения по user_id
    connections map[int64]map[*websocket.Conn]bool

    // Каналы
    broadcast   chan *ChatMessage
    register    chan *ClientSubscription
    unregister  chan *ClientSubscription

    // Защита
    mutex sync.RWMutex
}
```

**Функции:**
- `RegisterConnection(conn, userID)` - регистрация клиента
- `UnregisterConnection(conn, userID)` - отключение клиента
- `BroadcastNewMessage(chatID, message)` - broadcast нового сообщения
- `BroadcastMessageRead(chatID, messageID, userID)` - уведомление о прочтении
- `BroadcastTyping(chatID, userID, isTyping)` - индикатор печати

### 2. HTTP WebSocket Handler (Listings Microservice)

```go
// internal/transport/http/chat_websocket.go
func (h *ChatWebSocketHandler) HandleChatWebSocket(c *fiber.Ctx) error {
    // 1. Получить userID из JWT токена (из query параметра)
    // 2. Upgrade to WebSocket
    // 3. Зарегистрировать в Hub
    // 4. Обработать входящие сообщения (ping/pong, typing)
    // 5. Отписаться при закрытии
}
```

**Endpoint:** `GET /ws/chat?token=<JWT>`

### 3. WebSocket Proxy (Backend Монолит)

```go
// backend/internal/proj/chat/module.go
app.Get("/ws/chat", func(c *fiber.Ctx) error {
    // 1. Проверить JWT токен из query параметра
    // 2. Создать WebSocket proxy к микросервису
    // 3. Пробросить все сообщения туда-обратно
})
```

## Типы сообщений

### От клиента к серверу:

```json
// Ping (keep-alive)
{"type": "ping"}

// Typing indicator
{"type": "typing", "chat_id": 123, "is_typing": true}

// Mark as read
{"type": "mark_read", "message_id": 456}
```

### От сервера к клиенту:

```json
// Подключение установлено
{"type": "connected", "user_id": 6, "timestamp": "2025-11-21T13:00:00Z"}

// Pong (ответ на ping)
{"type": "pong", "timestamp": "2025-11-21T13:00:00Z"}

// Новое сообщение
{
  "type": "new_message",
  "chat_id": 3,
  "message": {
    "id": 789,
    "chat_id": 3,
    "sender_id": 5,
    "content": "Привет!",
    "created_at": "2025-11-21T13:00:00Z"
  }
}

// Сообщение прочитано
{"type": "message_read", "message_id": 456, "chat_id": 3, "read_by": 6}

// Пользователь печатает
{"type": "typing", "chat_id": 3, "user_id": 5, "is_typing": true}
```

## Интеграция с Chat Service

### SendMessage - broadcast после создания

```go
func (s *chatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*Message, error) {
    // ... создание сообщения в БД ...

    // Broadcast через Hub
    if s.hub != nil {
        s.hub.BroadcastNewMessage(message.ChatID, message)
    }

    return message, nil
}
```

### MarkMessageAsRead - broadcast после обновления

```go
func (s *chatService) MarkMessageAsRead(ctx context.Context, messageID, userID int64) error {
    // ... обновление в БД ...

    // Broadcast через Hub
    if s.hub != nil {
        s.hub.BroadcastMessageRead(message.ChatID, messageID, userID)
    }

    return nil
}
```

## Frontend изменения

### WebSocket подключение

```typescript
// src/services/chat-websocket.ts
class ChatWebSocket {
  private ws: WebSocket | null = null;

  connect(token: string) {
    // Подключаемся через монолит (proxy)
    this.ws = new WebSocket(`ws://localhost:3000/ws/chat?token=${token}`);

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      switch(data.type) {
        case 'new_message':
          // Обновить Redux store
          dispatch(addMessage(data.message));
          // Показать уведомление
          showNotification(data.message);
          break;

        case 'message_read':
          // Обновить статус сообщения
          dispatch(markMessageRead(data.message_id));
          break;

        case 'typing':
          // Показать индикатор печати
          dispatch(setTyping(data.chat_id, data.user_id, data.is_typing));
          break;
      }
    };
  }

  sendTyping(chatId: number, isTyping: boolean) {
    this.send({ type: 'typing', chat_id: chatId, is_typing: isTyping });
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}
```

## План реализации

### Шаг 1: Hub в микросервисе ✅
1. Создать `internal/websocket/chat_hub.go`
2. Реализовать регистрацию/отписку клиентов
3. Реализовать broadcast методы

### Шаг 2: HTTP WebSocket Handler ✅
1. Создать `internal/transport/http/chat_websocket.go`
2. Реализовать JWT аутентификацию из query параметра
3. Подключить к Hub

### Шаг 3: Интеграция с Chat Service ✅
1. Передать Hub в `chatService`
2. Добавить broadcast вызовы в `SendMessage`
3. Добавить broadcast вызовы в `MarkMessageAsRead`

### Шаг 4: WebSocket Proxy в монолите ✅
1. Обновить `backend/internal/proj/chat/module.go`
2. Реализовать proxy к микросервису
3. Пробросить JWT токен

### Шаг 5: Frontend обновление ✅
1. Создать `ChatWebSocket` сервис
2. Подключить в `ChatWindow.tsx`
3. Обработать все типы сообщений
4. Добавить уведомления в header

### Шаг 6: Тестирование ✅
**Статус:** ✅ ГОТОВО К ТЕСТИРОВАНИЮ

**Что работает:**
- ✅ WebSocket Hub в микросервисе (порт 8086)
- ✅ WebSocket proxy в монолите (порт 3000)
- ✅ Frontend middleware для обработки WebSocket событий
- ✅ WebSocketManager для автоматического подключения
- ✅ BFF endpoint `/api/v2/ws-token` для получения токена

**Как протестировать:**
1. Запустить backend: `screen -r backend-3000`
2. Запустить microservice: `screen -r listings-microservice-50053`
3. Открыть два браузера/вкладки с разными пользователями
4. Отправить сообщение в одной вкладке
5. Проверить что сообщение появилось в другой вкладке в real-time

**Проверка WebSocket подключения:**
```bash
# Проверить логи backend
tail -f /tmp/backend.log | grep -i websocket

# Проверить логи microservice
tail -f /tmp/listings-microservice.log | grep -i websocket

# В браузерной консоли должны быть сообщения:
# [WebSocket] Connected
# WebSocket message: {type: "connected", ...}
```

## Безопасность

- ✅ JWT аутентификация через query параметр
- ✅ Проверка прав доступа к чату
- ✅ Rate limiting для WebSocket соединений
- ✅ Graceful shutdown с закрытием всех соединений

## Мониторинг

- Количество активных WebSocket соединений
- Метрики broadcast сообщений
- Ошибки соединений/разрывов
- Endpoint: `GET /api/v1/internal/chat/ws/stats`

## Fallback стратегия

Если WebSocket недоступен:
1. Frontend переключается на polling (каждые 5 секунд)
2. Показывается предупреждение пользователю
3. Автоматическое переподключение каждые 30 секунд

---

**Статус:** ✅ ГОТОВО К ТЕСТИРОВАНИЮ
**Приоритет:** 🔥 Высокий
**Дата создания:** 2025-11-21
**Дата завершения:** 2025-11-21
