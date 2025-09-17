# Интеграция Viber Bot через Infobip

## Обзор

Интеграция Viber Bot для SveTu Marketplace через платформу Infobip, обеспечивающая:
- Отправку уведомлений о доставке
- Real-time трекинг курьеров
- 24-часовые бесплатные окна для общения
- Rich Media сообщения с картами и кнопками

## Учетные данные Infobip

```
API Key: d225d6436a020569e4dee8468919de13-5d9211b6-0f93-4829-bd40-59855773e867
Base URL: d9vgp1.api.infobip.com
Тестовый баланс: 100 бесплатных сообщений
Триал период: 60 дней (до ноября 2025)
```

## Архитектура решения

### 1. Backend компоненты

#### Infobip Client (`internal/proj/viber/infobip/client.go`)
- HTTP клиент для Infobip API
- Поддержка всех типов сообщений (текст, изображения, кнопки, Rich Media)
- Обработка webhook уведомлений
- Отслеживание статусов доставки

#### Viber Bot Service (`internal/proj/viber/service/`)
- `bot_service.go` - прямая интеграция с Viber API
- `infobip_bot_service.go` - интеграция через Infobip
- `session_manager.go` - управление 24-часовыми сессиями

#### Tracking Services (`internal/proj/tracking/`)
- `delivery_service.go` - управление доставками
- `courier_service.go` - управление курьерами
- `websocket_hub.go` - WebSocket сервер для real-time обновлений

### 2. База данных

```sql
-- Таблицы для Viber Bot
viber_users         -- Пользователи Viber
viber_sessions      -- 24-часовые сессии для бесплатных сообщений
viber_messages      -- История сообщений

-- Таблицы для трекинга
couriers            -- Курьеры
deliveries          -- Доставки
delivery_location_history -- История локаций
```

### 3. WebSocket для real-time трекинга

```go
// Подключение к WebSocket
ws://svetu.rs/ws/tracking?token=TRACKING_TOKEN

// Формат сообщений
{
  "type": "location_update",
  "data": {
    "delivery_id": 123,
    "latitude": 44.787197,
    "longitude": 20.457273,
    "speed": 15.5,
    "heading": 180,
    "eta": "2025-09-17T15:30:00Z",
    "distance_meters": 2500
  }
}
```

## Настройка и запуск

### 1. Переменные окружения

Добавьте в `.env`:

```bash
# Infobip API
INFOBIP_API_KEY=d225d6436a020569e4dee8468919de13-5d9211b6-0f93-4829-bd40-59855773e867
INFOBIP_BASE_URL=d9vgp1.api.infobip.com
INFOBIP_SENDER_ID=SveTuBot  # Нужно зарегистрировать в Infobip

# Viber Bot (если используется напрямую)
VIBER_AUTH_TOKEN=your_viber_token
VIBER_BOT_NAME=SveTu Marketplace
VIBER_BOT_AVATAR=https://svetu.rs/logo-720.png
VIBER_WEBHOOK_URL=https://svetu.rs/api/viber/webhook

# Mapbox для карт
MAPBOX_ACCESS_TOKEN=your_mapbox_token
```

### 2. Применение миграций

```bash
cd backend
./migrator up
```

### 3. Регистрация Viber Service в Infobip

1. Войдите в Infobip Dashboard: https://portal.infobip.com
2. Перейдите в Channels → Viber
3. Нажмите "Create Viber Service"
4. Заполните данные:
   - Service Name: SveTu Marketplace
   - Service ID: SveTuBot
   - Description: Маркетплейс с доставкой в Сербии
   - Category: Shopping
   - Webhook URL: https://svetu.rs/api/viber/infobip-webhook

### 4. Запуск сервера

```bash
cd backend
go run cmd/api/main.go
```

## API Endpoints

### Viber Bot

```http
# Webhook для Infobip
POST /api/viber/infobip-webhook

# Отправка сообщения
POST /api/viber/send
{
  "viber_id": "user_viber_id",
  "text": "Ваш заказ в пути!"
}

# Отправка уведомления о трекинге
POST /api/viber/send-tracking
{
  "viber_id": "user_viber_id",
  "delivery_id": 123
}
```

### Courier API

```http
# Обновление локации курьера
POST /api/courier/location
{
  "courier_id": 1,
  "latitude": 44.787197,
  "longitude": 20.457273,
  "speed": 15.5,
  "heading": 180
}

# Получение активных доставок
GET /api/courier/{id}/deliveries
```

### Tracking API

```http
# Получение информации о доставке
GET /api/tracking/{token}

# WebSocket подключение
WS /ws/tracking?token={tracking_token}
```

## Примеры использования

### Отправка текстового сообщения через Infobip

```go
func SendNotification(ctx context.Context, viberID, text string) error {
    cfg := config.LoadViberConfig()
    db := postgres.NewDB()
    
    service := service.NewInfobipBotService(cfg, db)
    return service.SendTextMessage(ctx, viberID, text)
}
```

### Отправка Rich Media с картой

```go
func SendTrackingMap(ctx context.Context, viberID string, delivery *DeliveryInfo) error {
    cfg := config.LoadViberConfig()
    db := postgres.NewDB()
    
    service := service.NewInfobipBotService(cfg, db)
    return service.SendTrackingNotification(ctx, viberID, delivery)
}
```

### Обработка webhook от Infobip

```go
func HandleInfobipWebhook(c *fiber.Ctx) error {
    var webhook infobip.ViberWebhook
    if err := c.BodyParser(&webhook); err != nil {
        return err
    }
    
    service := service.NewInfobipBotService(cfg, db)
    return service.ProcessWebhook(c.Context(), &webhook)
}
```

## Стоимость и биллинг

### Бесплатные сообщения (в рамках 24-часовой сессии)
- Пользователь инициирует диалог
- Все сообщения бесплатны в течение 24 часов
- Автоматическое отслеживание сессий

### Платные сообщения (вне сессии)
- Текстовые: ~€0.015 за сообщение
- Rich Media: ~€0.025 за сообщение
- Промо рассылки: требуют метку PROMOTION

### Оптимизация расходов
1. Используйте 24-часовые окна максимально
2. Группируйте уведомления
3. Используйте webhook для подтверждения доставки
4. Мониторьте статистику через Infobip Dashboard

## Frontend страница трекинга

### Компоненты

1. **TrackingPage** (`frontend/svetu/src/app/[locale]/track/[token]/page.tsx`)
   - Отображение карты с Mapbox GL JS
   - WebSocket подключение для real-time обновлений
   - Показ информации о доставке

2. **MapView** (`frontend/svetu/src/components/tracking/MapView.tsx`)
   - Интеграция с Mapbox
   - Анимация движения курьера
   - Отображение товаров и магазинов поблизости

3. **DeliveryInfo** (`frontend/svetu/src/components/tracking/DeliveryInfo.tsx`)
   - Информация о заказе
   - ETA и расстояние
   - Контакты курьера

## Мониторинг и отладка

### Логи
```bash
# Backend логи
tail -f /var/log/svetu/backend.log

# WebSocket соединения
tail -f /var/log/svetu/websocket.log
```

### Метрики
- Количество активных сессий
- Количество отправленных сообщений
- Стоимость за период
- Процент доставленных сообщений

### Infobip Dashboard
- https://portal.infobip.com/analytics/viber
- Отслеживание расходов
- Статистика доставки
- Детальные отчеты

## Troubleshooting

### Проблема: Сообщения не доставляются
1. Проверьте статус в Infobip Dashboard
2. Убедитесь, что пользователь подписан на бота
3. Проверьте формат Viber ID

### Проблема: WebSocket не подключается
1. Проверьте токен трекинга
2. Убедитесь, что nginx пропускает WebSocket
3. Проверьте CORS настройки

### Проблема: Карта не отображается
1. Проверьте Mapbox токен
2. Убедитесь, что координаты корректны
3. Проверьте CSP заголовки

## Следующие шаги

1. ✅ Интеграция с Infobip API
2. ✅ Создание сервисов для курьеров и доставок
3. ✅ WebSocket сервер для трекинга
4. ⏳ Frontend страница трекинга
5. ⏳ Настройка API endpoints
6. ⏳ Тестирование с реальными Viber аккаунтами
7. ⏳ Оптимизация стоимости сообщений
8. ⏳ Добавление аналитики и метрик

## Контакты и поддержка

- Infobip Support: support@infobip.com
- Документация API: https://www.infobip.com/docs/api
- Viber Bot документация: https://developers.viber.com/docs/api/
- Mapbox документация: https://docs.mapbox.com/