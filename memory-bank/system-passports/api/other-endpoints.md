# Паспорт API Endpoints: Остальные группы

## 📋 Обзор
Этот файл содержит паспорта для оставшихся 7 групп API endpoints:
- Chat (💬)
- Search (🔍) 
- Balance (💰)
- Notifications (🔔)
- Analytics (📊)
- Admin (👨‍💼)
- System (🔧)

---

## 💬 Chat Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/marketplace/chat`
- **Handler**: `backend/internal/proj/marketplace/handler/handler.go`
- **Количество endpoints**: 6 + WebSocket
- **Интеграции**: PostgreSQL, WebSocket, Redis

### 🎯 Назначение
Real-time чат система для коммуникации покупателей и продавцов

### 📡 Endpoints
```typescript
GET /api/v1/marketplace/chat/
// Получение списка чатов пользователя
// Handler: h.Chat.GetChats
// Response: Chat[] с последними сообщениями

GET /api/v1/marketplace/chat/messages
// Получение сообщений чата с пагинацией
// Handler: h.Chat.GetMessages
// Query: chat_id, before_id, limit

POST /api/v1/marketplace/chat/messages
// Отправка сообщения
// Handler: h.Chat.SendMessage
// Body: {chat_id, content, attachments?}
// WebSocket: Уведомления в реальном времени

PUT /api/v1/marketplace/chat/messages/read
// Отметить сообщения как прочитанные
// Handler: h.Chat.MarkAsRead
// Body: {chat_id, message_ids}

POST /api/v1/marketplace/chat/:chat_id/archive
// Архивировать чат
// Handler: h.Chat.ArchiveChat

GET /api/v1/marketplace/chat/unread-count
// Количество непрочитанных сообщений
// Handler: h.Chat.GetUnreadCount

// WebSocket
WS /ws/chat
// Реальные уведомления о сообщениях
// Handler: s.marketplace.Chat.HandleWebSocketWithAuth
```

### 🎭 Структуры данных
```typescript
interface Chat {
  id: string;
  listing_id: string;
  buyer_id: string;
  seller_id: string;
  last_message?: Message;
  unread_count: number;
  status: "active" | "archived" | "blocked";
  created_at: string;
}

interface Message {
  id: string;
  chat_id: string;
  sender_id: string;
  content: string;
  attachments: Attachment[];
  type: "text" | "image" | "file" | "system";
  read_at?: string;
  created_at: string;
}
```

---

## 🔍 Search Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/search`
- **Handler**: `backend/internal/proj/global/handler/handler.go`
- **Количество endpoints**: 1
- **Интеграции**: OpenSearch (multiple indices)

### 🎯 Назначение
Унифицированный поиск по всем сущностям платформы

### 📡 Endpoints
```typescript
GET /api/v1/search
// Глобальный поиск по всей платформе
// Handler: h.UnifiedSearch.UnifiedSearch
// Query: q, type, filters, sort
// Searches: listings, storefronts, users, categories
// Response: UnifiedSearchResults с группировкой по типам
```

### 🎭 Структуры данных
```typescript
interface UnifiedSearchResults {
  query: string;
  total_results: number;
  results: {
    listings: SearchResult[];
    storefronts: SearchResult[];
    users: SearchResult[];
    categories: SearchResult[];
  };
  suggestions: string[];
  took: number;
}
```

---

## 💰 Balance Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/balance`
- **Handler**: `backend/internal/proj/balance/handler/routes.go`
- **Количество endpoints**: 4
- **Интеграции**: PostgreSQL, Payment providers

### 🎯 Назначение
Управление балансами пользователей и транзакциями

### 📡 Endpoints
```typescript
GET /api/v1/balance/
// Получение текущего баланса пользователя
// Handler: h.Balance.GetBalance
// Response: UserBalance с мультивалютными счетами

GET /api/v1/balance/transactions
// История транзакций с фильтрацией
// Handler: h.Balance.GetTransactions
// Query: type, date_from, date_to, currency

GET /api/v1/balance/payment-methods
// Привязанные платежные методы
// Handler: h.Balance.GetPaymentMethods
// Response: PaymentMethod[] с токенами

POST /api/v1/balance/deposit
// Пополнение баланса
// Handler: h.Balance.CreateDeposit
// Body: {amount, currency, payment_method_id}
```

### 🎭 Структуры данных
```typescript
interface UserBalance {
  user_id: string;
  balances: {
    RSD: number;
    EUR: number;
  };
  pending_transactions: PendingTransaction[];
  total_earned: number;
  total_spent: number;
}

interface BalanceTransaction {
  id: string;
  user_id: string;
  type: "deposit" | "withdrawal" | "payment" | "refund" | "fee";
  amount: number;
  currency: string;
  status: "pending" | "completed" | "failed";
  description: string;
  created_at: string;
}
```

---

## 🔔 Notifications Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/notifications`
- **Handler**: `backend/internal/proj/notifications/handler/routes.go`
- **Количество endpoints**: 9 (2 публичных, 7 защищенных)
- **Интеграции**: PostgreSQL, Telegram Bot, Email service

### 🎯 Назначение
Система уведомлений через множественные каналы

### 📡 Endpoints

#### Публичные
```typescript
POST /api/v1/notifications/telegram/webhook
// Webhook для Telegram бота
// Handler: h.HandleTelegramWebhook

POST /api/v1/notifications/email/public
// Отправка публичных email (контакты, поддержка)
// Handler: h.SendPublicEmail
```

#### Защищенные
```typescript
GET /api/v1/notifications/
// Список уведомлений пользователя
// Handler: h.GetNotifications
// Query: type, read, limit

GET /api/v1/notifications/settings
// Настройки уведомлений пользователя
// Handler: h.GetSettings

PUT /api/v1/notifications/settings
// Обновление настроек уведомлений
// Handler: h.UpdateSettings

GET /api/v1/notifications/telegram/status
// Статус подключения к Telegram
// Handler: h.GetTelegramStatus

GET /api/v1/notifications/telegram/token
// Получение токена для связи с Telegram
// Handler: h.GetTelegramToken

POST /api/v1/notifications/telegram/connect
// Подключение Telegram аккаунта
// Handler: h.ConnectTelegram

PUT /api/v1/notifications/:id/read
// Отметить уведомление как прочитанное
// Handler: h.MarkAsRead
```

### 🎭 Структуры данных
```typescript
interface Notification {
  id: string;
  user_id: string;
  type: NotificationType;
  title: string;
  content: string;
  data?: Record<string, any>;
  channels: ("in_app" | "email" | "telegram")[];
  read_at?: string;
  created_at: string;
}

type NotificationType = 
  | "new_message" | "listing_sold" | "payment_received"
  | "review_received" | "price_drop" | "system_alert";

interface NotificationSettings {
  user_id: string;
  email_enabled: boolean;
  telegram_enabled: boolean;
  preferences: {
    new_messages: boolean;
    sales_updates: boolean;
    reviews: boolean;
    marketing: boolean;
  };
}
```

---

## 📊 Analytics Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/analytics`
- **Handler**: `backend/internal/proj/analytics/routes/routes.go`
- **Количество endpoints**: 1
- **Интеграции**: PostgreSQL, ClickHouse (опционально)

### 🎯 Назначение
Сбор и анализ пользовательских событий

### 📡 Endpoints
```typescript
POST /api/v1/analytics/event
// Запись аналитического события
// Handler: h.RecordEvent
// Body: {event_type, properties, user_id?, session_id?}
// Events: page_view, listing_click, search, purchase
```

### 🎭 Структуры данных
```typescript
interface AnalyticsEvent {
  event_type: string;
  user_id?: string;
  session_id?: string;
  properties: Record<string, any>;
  timestamp: string;
  ip_address: string;
  user_agent: string;
}
```

---

## 👨‍💼 Admin Endpoints

### 📋 Метаданные
- **Базовый путь**: `/api/v1/admin`
- **Количество endpoints**: 30+
- **Security**: Требует admin роль
- **Интеграции**: Все модули системы

### 🎯 Назначение
Административная панель для управления платформой

### 📡 Основные группы

#### Управление пользователями
```typescript
GET /api/v1/admin/users
PUT /api/v1/admin/users/:id
DELETE /api/v1/admin/users/:id
PUT /api/v1/admin/users/:id/status
```

#### Управление администраторами
```typescript
GET /api/v1/admin/admins
POST /api/v1/admin/admins
DELETE /api/v1/admin/admins/:email
GET /api/v1/admin/admins/check/:email
```

#### Управление категориями
```typescript
POST /api/v1/admin/categories
GET /api/v1/admin/categories
PUT /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id
```

#### Системные операции
```typescript
POST /api/v1/admin/reindex-listings
POST /api/v1/admin/reindex-listings-with-translations
POST /api/v1/admin/sync-discounts
POST /api/v1/admin/reindex-ratings
```

---

## 🔧 System Endpoints

### 📋 Метаданные
- **Назначение**: Служебные и системные функции
- **Security**: Публичные или внутренние

### 📡 Endpoints
```typescript
GET /
// Главная страница API
// Response: API info и версия

GET /api/health
// Health check для мониторинга
// Response: {status: "ok", version, uptime}

GET /swagger/*
// Swagger UI документация
// Static files

GET /docs/*
// Альтернативная документация
// Handler: docserver

GET /api/v1/csrf-token
// Получение CSRF токена
// Response: {token: string}
```

### Дополнительные служебные API

#### Geocoding
```typescript
GET /api/v1/geocode/reverse
// Обратное геокодирование
GET /api/v1/geocode/cities
// Предложения городов
```

#### Contacts & Privacy
```typescript
GET /api/v1/contacts/
POST /api/v1/contacts/
PUT /api/v1/contacts/:contact_user_id/status
GET /api/v1/contacts/privacy
PUT /api/v1/contacts/privacy
```

#### Documentation
```typescript
GET /api/v1/docs/files
// Список файлов документации
GET /api/v1/docs/content
// Содержимое файла документации
```

---

## 🎯 Итоговая статистика API

### По группам:
- **Auth**: 8 endpoints
- **Users**: 3 endpoints  
- **Marketplace**: 41 endpoints
- **Payments**: 6 endpoints
- **Reviews**: 16 endpoints
- **Storefronts**: 43 endpoints
- **Chat**: 6 endpoints + WebSocket
- **Search**: 1 endpoint
- **Balance**: 4 endpoints
- **Notifications**: 9 endpoints
- **Analytics**: 1 endpoint
- **Admin**: 30+ endpoints
- **System & Misc**: 15+ endpoints

### **Общий итог: ~180+ API endpoints**

### Архитектурные особенности:
- REST + WebSocket гибридная архитектура
- OpenAPI/Swagger документация
- Мультиязычность (ru/en)
- Rate limiting и CSRF защита
- Интеграция с PostgreSQL, OpenSearch, MinIO, Redis
- Microservices-ready модульная структура