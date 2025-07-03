# СИСТЕМНЫЙ ПАСПОРТ: Analytics Handler

## 📋 Обзор модуля

**Назначение**: Система аналитики и метрик для платформы витрин  
**Расположение**: `/backend/internal/proj/analytics/`  
**Тип**: Backend handler  
**Статус**: ✅ Активный  

### 🎯 Основные функции
- Сбор данных о событиях пользователей в реальном времени
- Агрегация аналитических данных для витрин
- Отслеживание конверсии и поведения пользователей
- Статистика трафика и продаж
- Анализ источников трафика и популярных товаров

## 🏗️ Архитектура модуля

### 📁 Структура файлов
```
backend/internal/proj/analytics/
├── module.go                   # Фабрика модуля и регистрация
├── handler/
│   └── analytics_handler.go    # HTTP обработчики событий
├── routes/
│   └── routes.go              # Регистрация маршрутов
└── service/
    └── analytics_service.go   # Бизнес-логика аналитики
```

### 🔧 Основные компоненты

#### Module (module.go:14-27)
```go
type Module struct {
    handler *handler.AnalyticsHandler
}

func NewModule(db *postgres.Database) *Module {
    storefrontRepo := postgres.NewStorefrontRepository(db)
    analyticsService := service.NewAnalyticsService(storefrontRepo)
    analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
    return &Module{handler: analyticsHandler}
}
```

#### AnalyticsHandler (handler/analytics_handler.go:13-22)
```go
type AnalyticsHandler struct {
    service service.AnalyticsService
}
```

#### AnalyticsService (service/analytics_service.go:9-12)
```go
type AnalyticsService interface {
    RecordEvent(ctx context.Context, event *EventData) error
}
```

## 🛠️ API Endpoints

### 🌐 Публичные маршруты

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| POST | `/api/v1/analytics/event` | RecordEvent | Запись события аналитики |

### 🔐 Защищенные маршруты (в Storefronts Handler)

| Метод | Путь | Функция | Описание |
|-------|------|---------|----------|
| GET | `/api/v1/storefronts/:id/analytics` | GetAnalytics | Получение аналитики витрины |

## 🗄️ Модели данных

### EventRequest (handler/analytics_handler.go:25-31)
```go
type EventRequest struct {
    StorefrontID int             `json:"storefront_id" validate:"required"`
    EventType    string          `json:"event_type" validate:"required,oneof=page_view product_view add_to_cart checkout order"`
    EventData    json.RawMessage `json:"event_data,omitempty"`
    SessionID    string          `json:"session_id" validate:"required"`
    UserID       *int            `json:"user_id,omitempty"`
}
```

### EventData (service/analytics_service.go:15-24)
```go
type EventData struct {
    StorefrontID int             `json:"storefront_id"`
    EventType    string          `json:"event_type"`
    EventData    json.RawMessage `json:"event_data"`
    SessionID    string          `json:"session_id"`
    UserID       *int            `json:"user_id,omitempty"`
    IPAddress    string          `json:"ip_address"`
    UserAgent    string          `json:"user_agent"`
    Referrer     string          `json:"referrer"`
}
```

### StorefrontEvent (storage/postgres)
```go
type StorefrontEvent struct {
    StorefrontID int             `json:"storefront_id"`
    EventType    EventType       `json:"event_type"`
    EventData    json.RawMessage `json:"event_data"`
    UserID       *int            `json:"user_id,omitempty"`
    SessionID    string          `json:"session_id"`
    IPAddress    string          `json:"ip_address,omitempty"`
    UserAgent    string          `json:"user_agent,omitempty"`
    Referrer     string          `json:"referrer,omitempty"`
}
```

### StorefrontAnalytics (domain/models)
```go
type StorefrontAnalytics struct {
    ID           int       `json:"id"`
    StorefrontID int       `json:"storefront_id"`
    Date         time.Time `json:"date"`
    
    // Трафик
    PageViews      int     `json:"page_views"`
    UniqueVisitors int     `json:"unique_visitors"`
    BounceRate     float64 `json:"bounce_rate"`
    AvgSessionTime int     `json:"avg_session_time"` // в секундах
    
    // Продажи
    OrdersCount    int     `json:"orders_count"`
    Revenue        float64 `json:"revenue"`
    AvgOrderValue  float64 `json:"avg_order_value"`
    ConversionRate float64 `json:"conversion_rate"`
    
    // Детальная информация
    PaymentMethodsUsage JSONB `json:"payment_methods_usage"`
    ProductViews        int   `json:"product_views"`
    AddToCartCount      int   `json:"add_to_cart_count"`
    CheckoutCount       int   `json:"checkout_count"`
    TrafficSources      JSONB `json:"traffic_sources"`
    TopProducts         JSONB `json:"top_products"`
    TopCategories       JSONB `json:"top_categories"`
    OrdersByCity        JSONB `json:"orders_by_city"`
    
    CreatedAt time.Time `json:"created_at"`
}
```

## 📊 Типы событий

### EventType константы
```go
const (
    EventPageView    EventType = "page_view"     // Просмотр страницы витрины
    EventProductView EventType = "product_view"  // Просмотр товара
    EventAddToCart   EventType = "add_to_cart"   // Добавление в корзину
    EventCheckout    EventType = "checkout"      // Начало оформления заказа
    EventOrder       EventType = "order"         // Завершение заказа
)
```

### Примеры использования

#### Просмотр витрины
```json
{
    "storefront_id": 123,
    "event_type": "page_view",
    "session_id": "sess_abc123",
    "event_data": {
        "page": "/storefront/123",
        "timestamp": "2024-01-15T10:30:00Z"
    }
}
```

#### Просмотр товара
```json
{
    "storefront_id": 123,
    "event_type": "product_view",
    "session_id": "sess_abc123",
    "user_id": 456,
    "event_data": {
        "product_id": 789,
        "product_name": "iPhone 15 Pro",
        "category": "electronics",
        "price": 999.99
    }
}
```

#### Добавление в корзину
```json
{
    "storefront_id": 123,
    "event_type": "add_to_cart",
    "session_id": "sess_abc123",
    "user_id": 456,
    "event_data": {
        "product_id": 789,
        "quantity": 1,
        "price": 999.99
    }
}
```

## 🔄 Бизнес-процессы

### Запись события (handler/analytics_handler.go:44-85)
1. **Валидация запроса**:
   - Проверка обязательных полей
   - Валидация типа события
   - Проверка storefront_id

2. **Обогащение данных**:
   - Извлечение IP адреса клиента
   - Получение User-Agent и Referrer
   - Добавление user_id из JWT токена (если авторизован)

3. **Сохранение в БД**:
   - Запись в таблицу `storefront_events`
   - Логирование ошибок

### Агрегация данных (Job: storefronts/jobs/analytics_aggregator.go)
Ежедневное выполнение:

1. **Сбор событий за день**:
   - Группировка по витринам и типам событий
   - Подсчет количества событий

2. **Вычисление метрик**:
   - Уникальные посетители по IP + User-Agent
   - Bounce rate (одностраничные сессии)
   - Среднее время сессии
   - Конверсия по воронке

3. **Анализ трафика**:
   - Источники трафика по Referrer
   - Популярные товары
   - География заказов

4. **Сохранение агрегатов**:
   - Запись в таблицу `storefront_analytics`
   - UPSERT на основе (storefront_id, date)

## 🗃️ База данных

### Таблица storefront_events
```sql
CREATE TABLE storefront_events (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES marketplace_storefronts(id),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB DEFAULT '{}',
    user_id INT,
    session_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Таблица storefront_analytics
```sql
CREATE TABLE storefront_analytics (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES marketplace_storefronts(id),
    date DATE NOT NULL,
    
    -- Трафик
    page_views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    bounce_rate DECIMAL(5,2) DEFAULT 0,
    avg_session_time INT DEFAULT 0,
    
    -- Продажи
    orders_count INT DEFAULT 0,
    revenue DECIMAL(10,2) DEFAULT 0,
    avg_order_value DECIMAL(10,2) DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    
    -- JSON агрегаты
    payment_methods_usage JSONB DEFAULT '{}',
    product_views INT DEFAULT 0,
    add_to_cart_count INT DEFAULT 0,
    checkout_count INT DEFAULT 0,
    traffic_sources JSONB DEFAULT '{}',
    top_products JSONB DEFAULT '[]',
    top_categories JSONB DEFAULT '[]',
    orders_by_city JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(storefront_id, date)
);
```

### Индексы для производительности
```sql
-- Для быстрого поиска событий по витрине и дате
CREATE INDEX idx_storefront_events_storefront_date 
ON storefront_events(storefront_id, created_at);

-- Для анализа сессий
CREATE INDEX idx_storefront_events_session 
ON storefront_events(session_id, created_at);

-- Для поиска аналитики по витрине и периоду
CREATE INDEX idx_storefront_analytics_storefront_date 
ON storefront_analytics(storefront_id, date);
```

## 🔒 Безопасность и валидация

### Input Validation (handler/analytics_handler.go:50-53)
```go
if req.StorefrontID <= 0 || req.EventType == "" || req.SessionID == "" {
    return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.validation_failed")
}
```

### Публичный доступ
- Endpoint `/analytics/event` доступен без авторизации
- Позволяет анонимную аналитику для неавторизованных пользователей
- IP и User-Agent используются для защиты от ботов

### Приватность данных
- Не сохраняются персональные данные в event_data
- IP адреса используются только для агрегации
- Возможность анонимного трекинга

## 🔗 Внешние интеграции

### StorefrontRepository
- Метод `RecordEvent()` для сохранения событий
- Метод `GetAnalytics()` для получения агрегированных данных
- Методы для вычисления метрик

### Storefronts Handler
- Endpoint для получения аналитики владельцами витрин
- Авторизация и проверка прав доступа
- Фильтрация по датам

## 📈 Метрики и KPI

### Трафиковые метрики
- **Page Views**: общее количество просмотров
- **Unique Visitors**: уникальные посетители (IP + User-Agent)
- **Bounce Rate**: процент одностраничных сессий
- **Avg Session Time**: среднее время сессии

### Конверсионные метрики
- **Conversion Rate**: процент конверсии из просмотра в заказ
- **Add to Cart Rate**: процент добавлений в корзину
- **Checkout Rate**: процент начала оформления заказа

### Коммерческие метрики
- **Revenue**: общая выручка
- **Orders Count**: количество заказов
- **Average Order Value**: средний чек
- **Payment Methods Usage**: распределение способов оплаты

### Продуктовые метрики
- **Top Products**: популярные товары
- **Top Categories**: популярные категории
- **Product Views**: просмотры товаров

## 🏭 Фабричные методы

### Module Factory (module.go:19-27)
```go
func NewModule(db *postgres.Database) *Module {
    storefrontRepo := postgres.NewStorefrontRepository(db)
    analyticsService := service.NewAnalyticsService(storefrontRepo)
    analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
    return &Module{handler: analyticsHandler}
}
```

### Service Factory (service/analytics_service.go:32-36)
```go
func NewAnalyticsService(storefrontRepo postgres.StorefrontRepository) AnalyticsService {
    return &analyticsServiceImpl{
        storefrontRepo: storefrontRepo,
    }
}
```

## ⚠️ Особенности реализации

### Real-time vs Batch обработка
- **Real-time**: запись событий происходит сразу
- **Batch**: агрегация выполняется ежедневно в background job
- Компромисс между производительностью и актуальностью данных

### Анонимная аналитика
- Поддержка неавторизованных пользователей
- Использование session_id для связывания событий
- IP-based уникальность посетителей

### JSON агрегаты
- Гибкое хранение сложных структур в JSONB
- Возможность добавления новых метрик без изменения схемы
- Эффективные запросы с GIN индексами

## 🔄 Связи с другими модулями

### Входящие зависимости
- `storefronts` handler - основной источник аналитических данных
- `marketplace` handler - данные о товарах и заказах
- `users` handler - информация о пользователях

### Исходящие зависимости
- PostgreSQL storage для сохранения событий и агрегатов
- Background jobs для обработки данных

## 🚀 TODO и улучшения

### Технические улучшения
- [ ] Real-time дашборды через WebSocket
- [ ] Экспорт данных в CSV/Excel
- [ ] API для кастомных отчетов
- [ ] Интеграция с внешними аналитическими системами

### Функциональные улучшения
- [ ] A/B тестирование
- [ ] Когорный анализ пользователей
- [ ] Прогнозирование продаж
- [ ] Alerting при аномалиях

### Производительность
- [ ] Партиционирование таблицы событий
- [ ] Архивирование старых данных
- [ ] Кэширование популярных запросов
- [ ] Асинхронная обработка событий

## 📊 Примеры запросов

### Получение аналитики витрины
```http
GET /api/v1/storefronts/123/analytics?from=2024-01-01&to=2024-01-31
Authorization: Bearer <jwt_token>
```

### Запись события просмотра товара
```http
POST /api/v1/analytics/event
Content-Type: application/json

{
    "storefront_id": 123,
    "event_type": "product_view",
    "session_id": "sess_abc123",
    "event_data": {
        "product_id": 456,
        "category": "electronics"
    }
}
```

---

**Дата создания**: $(date)  
**Версия**: 1.0  
**Статус**: ✅ Активный модуль  
**Последнее обновление**: Система real-time аналитики с batch агрегацией