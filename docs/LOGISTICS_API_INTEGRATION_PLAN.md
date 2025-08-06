# План интеграции с API логистики

## Обзор

Этот документ описывает план интеграции маркетплейса с внешним API логистики для управления доставкой заказов.

## Текущее состояние

### ✅ Готово:
1. **Структура заказов в БД**
   - Таблицы: `storefront_orders`, `storefront_order_items`
   - Поля для логистики: `shipping_method`, `shipping_provider`, `tracking_number`, `shipped_at`, `delivered_at`

2. **API для создания заказов**
   - POST `/api/v1/orders` - создание заказа с уменьшением остатков товара
   - GET `/api/v1/orders` - получение списка заказов пользователя
   - PUT `/api/v1/orders/:id/cancel` - отмена заказа

3. **Статусы заказа**
   - `pending` - новый заказ
   - `confirmed` - подтвержден
   - `shipped` - отправлен
   - `delivered` - доставлен
   - `cancelled` - отменен

## Архитектура интеграции

### 1. Новый модуль логистики

```
backend/internal/proj/logistics/
├── handler/
│   ├── webhook_handler.go    # Обработчик webhook от логистики
│   └── shipping_handler.go   # API для управления доставкой
├── service/
│   ├── logistics_client.go   # Клиент для внешнего API
│   ├── shipping_service.go   # Бизнес-логика доставки
│   └── tracking_service.go   # Отслеживание посылок
├── repository/
│   └── shipping_repository.go
└── module.go

```

### 2. Таблица для логистических данных

```sql
CREATE TABLE logistics_shipments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES storefront_orders(id),
    external_id VARCHAR(100),           -- ID в системе логистики
    provider VARCHAR(50) NOT NULL,      -- Провайдер (DHL, FedEx, etc)
    tracking_number VARCHAR(100),
    status VARCHAR(50),                 -- Статус от провайдера
    pickup_address JSONB,               -- Берется из storefront_orders.pickup_address
    delivery_address JSONB,             -- Берется из storefront_orders.shipping_address
    package_info JSONB,                 -- Размеры, вес и т.д.
    estimated_delivery_date DATE,
    actual_delivery_date DATE,
    cost NUMERIC(12,2),
    currency CHAR(3) DEFAULT 'RSD',
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Примечание**: Адрес забора (pickup_address) автоматически заполняется из данных витрины при создании заказа.

### 3. Интерфейс для провайдеров логистики

```go
type LogisticsProvider interface {
    // Создать заявку на доставку
    CreateShipment(ctx context.Context, order *StorefrontOrder) (*Shipment, error)
    
    // Получить трекинг информацию
    GetTracking(ctx context.Context, trackingNumber string) (*TrackingInfo, error)
    
    // Отменить доставку
    CancelShipment(ctx context.Context, shipmentID string) error
    
    // Рассчитать стоимость доставки
    CalculateShippingCost(ctx context.Context, params ShippingParams) (*ShippingQuote, error)
    
    // Получить доступные методы доставки
    GetAvailableShippingMethods(ctx context.Context, params ShippingParams) ([]ShippingMethod, error)
}
```

## Процесс интеграции

### 1. При создании заказа
1. Сохраняем заказ со статусом `pending`
2. Резервируем товар на складе
3. Очищаем корзину
4. Перенаправляем на страницу успеха

### 2. После подтверждения оплаты
1. Меняем статус заказа на `confirmed`
2. Отправляем запрос в API логистики для создания доставки
3. Сохраняем tracking number и другие данные
4. Уведомляем покупателя по email/SMS

### 3. Webhook от логистики
1. Получаем обновления статуса доставки
2. Обновляем статус заказа
3. Уведомляем покупателя об изменениях

### 4. Отслеживание доставки
1. Страница отслеживания для покупателя
2. API для получения актуального статуса
3. Интеграция с картой для визуализации

## API Endpoints

### Новые endpoints:

```
POST   /api/v1/orders/:id/ship              # Создать доставку
GET    /api/v1/orders/:id/tracking          # Получить трекинг
PUT    /api/v1/orders/:id/delivery-address  # Изменить адрес доставки
GET    /api/v1/shipping/methods             # Доступные методы доставки
POST   /api/v1/shipping/calculate           # Рассчитать стоимость
POST   /api/v1/webhooks/logistics/:provider # Webhook для провайдера
```

## Конфигурация

```env
# Logistics API Configuration
LOGISTICS_PROVIDER=dhl                      # Активный провайдер
LOGISTICS_API_URL=https://api.dhl.com/v1
LOGISTICS_API_KEY=your-api-key
LOGISTICS_API_SECRET=your-api-secret
LOGISTICS_WEBHOOK_SECRET=webhook-secret
LOGISTICS_SANDBOX_MODE=true                 # Тестовый режим
```

## UI изменения

### 1. Checkout страница
- Выбор метода доставки с расчетом стоимости в реальном времени
- Валидация адреса доставки через API логистики

### 2. Страница заказа
- Отображение tracking number
- Кнопка "Отследить посылку"
- Timeline статусов доставки

### 3. Dashboard продавца
- Кнопка "Отправить заказ"
- Печать shipping label
- Массовая отправка заказов

## Безопасность

1. **Webhook валидация**
   - Проверка подписи запроса
   - IP whitelist для провайдеров
   - Rate limiting

2. **Хранение данных**
   - Шифрование чувствительных данных (адреса, телефоны)
   - Логирование всех операций
   - Аудит изменений

## Этапы реализации

### Этап 1: Базовая интеграция (1-2 недели)
- [ ] Создать модуль logistics
- [ ] Реализовать клиент для тестового API
- [ ] Добавить таблицу logistics_shipments
- [ ] Реализовать создание доставки при подтверждении заказа

### Этап 2: Отслеживание и уведомления (1 неделя)
- [ ] Webhook handler для обновлений статуса
- [ ] Email/SMS уведомления
- [ ] UI для отслеживания на frontend

### Этап 3: Расширенные функции (2 недели)
- [ ] Множественные провайдеры логистики
- [ ] Автоматический выбор оптимального провайдера
- [ ] Возвраты и обмены
- [ ] Интеграция с картами для визуализации

### Этап 4: Оптимизация (1 неделя)
- [ ] Кеширование расчетов стоимости
- [ ] Batch операции для массовых отправок
- [ ] Аналитика и отчеты по доставкам

## Тестирование

1. **Unit тесты**
   - Mock клиент для API логистики
   - Тесты бизнес-логики

2. **Интеграционные тесты**
   - Тестовый аккаунт у провайдера
   - E2E сценарии создания и отслеживания доставки

3. **Нагрузочное тестирование**
   - Обработка множественных webhook
   - Массовое создание доставок

## Мониторинг

- Метрики успешности создания доставок
- Время доставки по регионам
- Стоимость доставки
- SLA провайдеров
- Ошибки интеграции

## Риски и митигация

1. **Недоступность API провайдера**
   - Retry механизм с exponential backoff
   - Fallback на альтернативного провайдера
   - Очередь для отложенной обработки

2. **Изменения в API**
   - Версионирование клиентов
   - Мониторинг deprecated endpoints
   - Автоматические тесты совместимости

3. **Безопасность данных**
   - Регулярный аудит
   - Шифрование в покое и при передаче
   - Минимизация хранимых данных