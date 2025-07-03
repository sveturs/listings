# СИСТЕМНЫЙ ПАСПОРТ: Notifications Handler

## 📋 Обзор модуля

**Назначение**: Модуль управления системой уведомлений пользователей  
**Расположение**: `/backend/internal/proj/notifications/`  
**Тип**: Backend handler  
**Статус**: ✅ Активный  

### 🎯 Основные функции
- Управление настройками уведомлений пользователей
- Интеграция с Telegram Bot API для мгновенных уведомлений
- Отправка email уведомлений через SMTP
- Обработка webhook от Telegram бота
- Публичная форма обратной связи
- CRUD операции с уведомлениями в БД

## 🏗️ Архитектура модуля

### 📁 Структура файлов
```
backend/internal/proj/notifications/
├── handler/
│   ├── handler.go          # Основные HTTP handlers
│   ├── routes.go           # Регистрация маршрутов
│   └── responses.go        # Структуры ответов
├── service/
│   ├── interface.go        # Интерфейс сервиса
│   ├── service.go          # Фабрика сервисов
│   ├── notification.go     # Бизнес-логика уведомлений
│   └── email.go           # Email сервис
└── storage/
    ├── interface.go        # Интерфейс репозитория
    └── postgres/
        └── notifications.go # PostgreSQL реализация
```

### 🔧 Основные компоненты

#### Handler (handler.go:26-42)
```go
type Handler struct {
    notificationService service.NotificationServiceInterface
    bot                 *tgbotapi.BotAPI
}
```

#### NotificationService (service/notification.go:15-42)
```go
type NotificationService struct {
    storage storage.Storage
    bot     *tgbotapi.BotAPI
    email   *EmailService
}
```

#### EmailService (service/email.go:10-28)
```go
type EmailService struct {
    smtpHost     string
    smtpPort     string
    senderEmail  string
    senderName   string
    smtpUsername string
    smtpPassword string
}
```

## 🛠️ API Endpoints

### 🔐 Защищенные маршруты (JWT Auth)

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| GET | `/api/v1/notifications` | GetNotifications | Получить список уведомлений |
| GET | `/api/v1/notifications/settings` | GetSettings | Получить настройки уведомлений |
| PUT | `/api/v1/notifications/settings` | UpdateSettings | Обновить настройки уведомлений |
| GET | `/api/v1/notifications/telegram/status` | GetTelegramStatus | Статус подключения Telegram |
| GET | `/api/v1/notifications/telegram/token` | GetTelegramToken | Токен для подключения Telegram |
| POST | `/api/v1/notifications/telegram/connect` | ConnectTelegram | Подключить Telegram |
| PUT | `/api/v1/notifications/:id/read` | MarkAsRead | Отметить как прочитанное |

### 🌐 Публичные маршруты

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| POST | `/api/v1/notifications/telegram/webhook` | HandleTelegramWebhook | Webhook Telegram бота |
| POST | `/api/v1/notifications/email/public` | SendPublicEmail | Публичная форма обратной связи |

## 🗄️ Модели данных

### Notification (domain/models/notification.go)
```go
type Notification struct {
    ID          int             `json:"id"`
    UserID      int             `json:"user_id"`
    Type        string          `json:"type"`
    Title       string          `json:"title"`
    Message     string          `json:"message"`
    ListingID   int             `json:"listing_id,omitempty"`
    Data        json.RawMessage `json:"data,omitempty"`
    IsRead      bool            `json:"is_read"`
    DeliveredTo json.RawMessage `json:"delivered_to"`
    CreatedAt   time.Time       `json:"created_at"`
}
```

### NotificationSettings (domain/models/notification.go)
```go
type NotificationSettings struct {
    UserID           int       `json:"user_id"`
    NotificationType string    `json:"notification_type"`
    TelegramEnabled  bool      `json:"telegram_enabled"`
    EmailEnabled     bool      `json:"email_enabled"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

### TelegramConnection (domain/models/notification.go)
```go
type TelegramConnection struct {
    UserID           int       `json:"user_id"`
    TelegramChatID   string    `json:"telegram_chat_id"`
    TelegramUsername string    `json:"telegram_username"`
    ConnectedAt      time.Time `json:"connected_at"`
}
```

## 📊 Типы уведомлений

### Константы (domain/models/models.go)
```go
const (
    NotificationTypeNewMessage     = "new_message"      // Новые сообщения в чатах
    NotificationTypeNewReview      = "new_review"       // Новые отзывы
    NotificationTypeReviewVote     = "review_vote"      // Голосование за отзывы
    NotificationTypeReviewResponse = "review_response"  // Ответы на отзывы
    NotificationTypeListingStatus  = "listing_status"   // Изменение статуса объявлений
    NotificationTypeFavoritePrice  = "favorite_price"   // Изменение цен в избранном
)
```

## 🔒 Безопасность и аутентификация

### JWT Middleware
- Все защищенные маршруты требуют JWT токен
- Извлечение `user_id` из токена через `c.Locals("user_id")`

### Telegram Security
- HMAC SHA256 подпись для токенов подключения (handler.go:226-233)
- Валидация токенов при подключении (handler.go:235-257)
- Секретный ключ: `TELEGRAM_BOT_TOKEN`

### Email Security
- CORS заголовки для публичной формы
- Валидация входных данных
- Логирование всех email операций

## 🗃️ База данных

### Связанные таблицы
- `notifications` - основные уведомления
- `notification_settings` - настройки пользователей
- `user_telegram_connections` - подключения Telegram
- `users` - связь с пользователями

### Основные операции
```sql
-- Получение уведомлений пользователя
SELECT id, user_id, type, title, message, data, is_read, delivered_to, created_at
FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3

-- Обновление настроек (UPSERT)
INSERT INTO notification_settings (user_id, notification_type, telegram_enabled, email_enabled)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, notification_type) 
DO UPDATE SET telegram_enabled = EXCLUDED.telegram_enabled, email_enabled = EXCLUDED.email_enabled

-- Сохранение Telegram подключения
INSERT INTO user_telegram_connections (user_id, telegram_chat_id, telegram_username)
VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET ...
```

## 🔗 Внешние интеграции

### Telegram Bot API
- **Библиотека**: `github.com/go-telegram-bot-api/telegram-bot-api`
- **Webhook URL**: `https://svetu.rs/api/v1/notifications/telegram/webhook`
- **Команды**: `/start` с токеном для подключения
- **Функции**: отправка уведомлений, обработка команд

### SMTP Email Service
- **Сервер**: `mailserver:25` (без TLS)
- **От**: `info@svetu.rs`
- **Поддержка HTML**: да
- **Режимы**: прямая отправка + ручное соединение

## 📈 Бизнес-логика

### Отправка уведомлений (service/notification.go:184-240)
1. Проверка настроек пользователя по типу уведомления
2. Создание записи в БД
3. Отправка в Telegram (если включено)
4. Отправка на email (если включено)
5. Логирование результатов

### Подключение Telegram (handler.go:127-200)
1. Валидация токена из команды `/start`
2. Извлечение `user_id` из подписанного токена
3. Сохранение `chat_id` и `username`
4. Создание базовых настроек уведомлений
5. Подтверждение подключения

### Публичная форма обратной связи (handler.go:338-437)
1. Парсинг данных формы (name, email, message, source)
2. Определение получателя по источнику (`klimagrad` → `klimagrad@svetu.rs`)
3. Ручная отправка через SMTP без TLS
4. Полное логирование процесса

## 🏭 Фабричные методы

### Service Factory (service/service.go:12-16)
```go
func NewService(storage storage.Storage) *Service {
    return &Service{
        Notification: NewNotificationService(storage),
    }
}
```

### Handler Factory (handler.go:31-42)
```go
func NewHandler(service service.NotificationServiceInterface) *Handler {
    // Инициализация Telegram бота
    // Возврат настроенного handler
}
```

## 🔧 Конфигурация

### Переменные окружения
- `TELEGRAM_BOT_TOKEN` - токен Telegram бота
- `EMAIL_PASSWORD` - пароль SMTP (опционально)

### SMTP настройки
- Host: `mailserver`
- Port: `25`
- Auth: Plain (без TLS)
- From: `info@svetu.rs`

## 📝 Структуры ответов

### TelegramTokenResponse (responses.go:13-18)
```go
type TelegramTokenResponse struct {
    Token       string    `json:"token"`
    GeneratedAt time.Time `json:"generated_at"`
}
```

### NotificationSettingsResponse (responses.go:38-41)
```go
type NotificationSettingsResponse struct {
    Data []models.NotificationSettings `json:"data"`
}
```

### PublicEmailSendResponse (responses.go:54-59)
```go
type PublicEmailSendResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

## ⚠️ Особенности реализации

### Управление настройками
- Автоматическое создание базовых настроек при первом запросе
- Partial update - обновление только переданных полей
- Сохранение существующих значений при частичном обновлении

### Email отправка
- Два режима: `smtp.SendMail` и ручное соединение
- Fallback при ошибке первого способа
- HTML шаблоны с ссылками на объявления

### Telegram интеграция
- Автоматическая настройка webhook при инициализации
- Обработка ошибок webhook через логи
- Генерация безопасных токенов для подключения

## 🔄 Связи с другими модулями

### Входящие зависимости
- `users` handler - получение информации о пользователе
- `marketplace` handler - уведомления об объявлениях
- `reviews` handler - уведомления об отзывах
- `payments` handler - уведомления о платежах

### Исходящие зависимости
- PostgreSQL storage для всех операций с БД
- Telegram Bot API для мгновенных уведомлений
- SMTP сервер для email уведомлений

## 🚀 TODO и улучшения

### Технические улучшения
- [ ] Добавить retry механизм для Telegram/Email
- [ ] Реализовать batch отправку уведомлений
- [ ] Добавить метрики отправки/доставки
- [ ] Кэширование настроек пользователей

### Функциональные улучшения
- [ ] Push уведомления для мобильных устройств
- [ ] Шаблоны уведомлений с переменными
- [ ] Группировка уведомлений по типам
- [ ] Отписка от уведомлений по one-click ссылке

### Безопасность
- [ ] Rate limiting для публичных endpoint'ов
- [ ] Encryption для Telegram chat_id
- [ ] Audit log для изменения настроек

## 📊 Метрики и мониторинг

### Логируемые события
- Отправка каждого уведомления
- Ошибки Telegram/Email доставки  
- Подключения/отключения Telegram
- Обновления настроек пользователей

### Рекомендуемые метрики
- Количество отправленных уведомлений по типам
- Success rate доставки Telegram/Email
- Время отклика webhook Telegram
- Активность пользователей по настройкам

---

**Дата создания**: $(date)  
**Версия**: 1.0  
**Статус**: ✅ Активный модуль  
**Последнее обновление**: Обработка email уведомлений и Telegram интеграция