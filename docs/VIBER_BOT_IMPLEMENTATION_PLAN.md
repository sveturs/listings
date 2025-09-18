# Viber Bot для Маркетплейса SveTu - План реализации

## 1. Коммерческая модель и стоимость

### Ежемесячные расходы
- **Базовая плата**: €100/месяц за бота
- **Сообщения**:
  - Бесплатно: все ответы в течение 24 часов
  - ~€0.0045 за проактивное уведомление в Сербии

### Оптимизация расходов
- Максимально использовать 24-часовые сессии
- Группировать уведомления
- Использовать кнопки быстрых ответов для инициации диалогов

## 2. Процесс создания бота

### Шаг 1: Регистрация
```
1. Обратиться к официальному партнёру Viber в Сербии:
   - Infobip (https://www.infobip.com/viber-business)
   - BulkGate (https://www.bulkgate.com/en/pricing/viber/rs/serbia/)
   - NTH (https://nth.rs/en/channels/viber/)

2. Или напрямую: https://www.forbusiness.viber.com/
```

### Шаг 2: Подготовка инфраструктуры
```go
// backend/internal/proj/viber/config.go
type ViberBotConfig struct {
    AuthToken    string `env:"VIBER_AUTH_TOKEN"`
    BotName      string `env:"VIBER_BOT_NAME"`
    WebhookURL   string `env:"VIBER_WEBHOOK_URL"`
    AvatarURL    string `env:"VIBER_AVATAR_URL"`
}
```

### Шаг 3: Webhook endpoint
```go
// backend/internal/proj/viber/handler.go
package viber

import (
    "github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
    botService *BotService
}

func (h *WebhookHandler) HandleWebhook(c *fiber.Ctx) error {
    var event ViberEvent
    if err := c.BodyParser(&event); err != nil {
        return err
    }

    switch event.Event {
    case "message":
        return h.handleMessage(c, event)
    case "subscribed":
        return h.handleSubscription(c, event)
    case "conversation_started":
        return h.handleConversationStart(c, event)
    }

    return c.SendStatus(200)
}
```

## 3. Функциональность бота

### Основные команды
```yaml
/start - Приветствие и главное меню
/search - Поиск товаров
/orders - Мои заказы
/cart - Корзина
/help - Помощь
/storefronts - Мои витрины (для продавцов)
```

### Rich Media меню
```json
{
  "Type": "rich_media",
  "ButtonsGroupColumns": 6,
  "ButtonsGroupRows": 2,
  "Buttons": [
    {
      "ActionType": "reply",
      "ActionBody": "search",
      "Text": "🔍 Поиск товаров",
      "TextSize": "medium",
      "Columns": 3,
      "Rows": 1
    },
    {
      "ActionType": "reply",
      "ActionBody": "categories",
      "Text": "📂 Категории",
      "TextSize": "medium",
      "Columns": 3,
      "Rows": 1
    },
    {
      "ActionType": "open-url",
      "ActionBody": "https://svetu.rs/ru/create-listing-choice",
      "Text": "➕ Создать объявление",
      "TextSize": "medium",
      "Columns": 6,
      "Rows": 1
    }
  ]
}
```

## 4. Интеграция с маркетплейсом

### Уведомления (платные, вне сессии)
```go
type NotificationService struct {
    viberBot *ViberBot
}

func (s *NotificationService) SendOrderUpdate(userID, orderID string) error {
    // Проверяем, есть ли активная сессия (24 часа)
    if s.HasActiveSession(userID) {
        // Бесплатное сообщение
        return s.SendMessage(userID, formatOrderUpdate(orderID))
    }

    // Платное уведомление - использовать экономно
    if s.IsHighPriorityNotification(orderID) {
        return s.SendProactiveMessage(userID, formatOrderUpdate(orderID))
    }

    // Откладываем до следующей сессии
    return s.QueueNotification(userID, orderID)
}
```

### Поиск товаров через бота
```go
func (h *BotHandler) handleSearch(query string, senderID string) {
    // Используем существующий поиск
    results, err := h.searchService.Search(context.Background(), &SearchParams{
        Query: query,
        Limit: 5,
    })

    if err != nil {
        h.sendError(senderID, "Ошибка поиска")
        return
    }

    // Формируем карусель товаров
    carousel := h.buildProductCarousel(results)
    h.sendRichMedia(senderID, carousel)
}
```

## 5. База данных

### Миграция для Viber пользователей
```sql
-- migrations/000030_viber_bot_users.up.sql
CREATE TABLE viber_users (
    id SERIAL PRIMARY KEY,
    viber_id VARCHAR(100) UNIQUE NOT NULL,
    user_id INT REFERENCES users(id),
    name VARCHAR(255),
    avatar_url TEXT,
    language VARCHAR(10) DEFAULT 'sr',
    subscribed BOOLEAN DEFAULT true,
    last_session_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE viber_sessions (
    id SERIAL PRIMARY KEY,
    viber_user_id INT REFERENCES viber_users(id),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_message_at TIMESTAMP WITH TIME ZONE,
    context JSONB,
    active BOOLEAN DEFAULT true
);

CREATE INDEX idx_viber_users_viber_id ON viber_users(viber_id);
CREATE INDEX idx_viber_sessions_active ON viber_sessions(active, last_message_at);
```

## 6. Аналитика и оптимизация

### Метрики для отслеживания
```go
type ViberMetrics struct {
    TotalUsers          int
    ActiveSessions      int
    FreeMessages        int
    PaidMessages        int
    MonthlyMessageCost  float64
    ConversionRate      float64
}
```

### Dashboard для мониторинга
- Количество подписчиков
- Активные сессии
- Расходы на сообщения
- Конверсия: просмотр → корзина → покупка

## 7. Примерный код для Go

### Установка SDK
```bash
go get github.com/viber/viber-bot-go
```

### Базовая структура бота
```go
package viber

import (
    "github.com/viber/viber-bot-go"
    "github.com/viber/viber-bot-go/model"
)

type MarketplaceBot struct {
    bot           *viber.Bot
    searchService *search.Service
    orderService  *orders.Service
    cartService   *cart.Service
}

func NewMarketplaceBot(config *ViberBotConfig) (*MarketplaceBot, error) {
    bot := &viber.Bot{
        AppKey: config.AuthToken,
        Sender: viber.Sender{
            Name:   config.BotName,
            Avatar: config.AvatarURL,
        },
        Message: TextMessageHandler,
    }

    return &MarketplaceBot{
        bot: bot,
    }, nil
}

func (b *MarketplaceBot) Start() error {
    // Установка webhook
    _, err := b.bot.SetWebhook(b.config.WebhookURL, nil)
    return err
}
```

## 8. Преимущества для маркетплейса

### Для покупателей:
- Быстрый поиск товаров
- Отслеживание заказов
- Уведомления о скидках
- Чат с продавцами

### Для продавцов:
- Уведомления о новых заказах
- Управление товарами
- Статистика продаж
- Прямая связь с покупателями

## 9. Timeline реализации

### Фаза 1 (2 недели)
- Регистрация бота
- Базовый webhook
- Приветствие и меню

### Фаза 2 (2 недели)
- Поиск товаров
- Показ категорий
- Детали товаров

### Фаза 3 (2 недели)
- Корзина
- Оформление заказов
- Уведомления

### Фаза 4 (1 неделя)
- Тестирование
- Оптимизация
- Запуск

## 10. Бюджет (первый год)

```
Фиксированные расходы:
- €100/месяц × 12 = €1,200

Переменные расходы (прогноз):
- 10,000 проактивных сообщений/месяц × €0.0045 = €45/месяц
- Годовые расходы на сообщения: €540

Итого: ~€1,740/год
```

## 11. Альтернативы и дополнения

### Telegram Bot (бесплатно)
- Без ежемесячной платы
- Популярен в IT-сообществе
- Меньше пользователей в Сербии

### WhatsApp Business API
- Дороже Viber
- Больше международных пользователей

### SMS уведомления
- Дороже (€0.01-0.03 за SMS)
- 100% доставляемость
- Для критичных уведомлений

## 12. Контакты партнёров в Сербии

1. **Infobip** (Глобальный партнёр)
   - Web: https://www.infobip.com
   - Офис в Белграде

2. **NTH** (Локальный партнёр)
   - Web: https://nth.rs
   - Специализация на Viber в Сербии

3. **BulkGate**
   - Web: https://www.bulkgate.com
   - Поддержка DOO компаний