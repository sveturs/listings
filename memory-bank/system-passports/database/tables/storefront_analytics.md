# Паспорт таблицы `storefront_analytics`

## Назначение
Система аналитики для витрин магазинов. Агрегирует ежедневную статистику по трафику, продажам, конверсии и другим ключевым метрикам для владельцев витрин.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS storefront_analytics (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    
    -- Трафик
    page_views INT DEFAULT 0,
    unique_visitors INT DEFAULT 0,
    bounce_rate DECIMAL(5,2) DEFAULT 0,
    avg_session_time INT DEFAULT 0, -- в секундах
    
    -- Продажи
    orders_count INT DEFAULT 0,
    revenue DECIMAL(10,2) DEFAULT 0,
    avg_order_value DECIMAL(10,2) DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    
    -- JSON поля для детальной информации
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
    
    -- Уникальный индекс для предотвращения дубликатов
    UNIQUE(storefront_id, date)
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID записи |
| `storefront_id` | INT | NOT NULL FK | ID витрины |
| `date` | DATE | NOT NULL | Дата аналитики |
| **Трафик** | | | |
| `page_views` | INT | DEFAULT 0 | Количество просмотров страниц |
| `unique_visitors` | INT | DEFAULT 0 | Уникальные посетители |
| `bounce_rate` | DECIMAL(5,2) | DEFAULT 0 | Процент отказов |
| `avg_session_time` | INT | DEFAULT 0 | Среднее время сессии (сек) |
| **Продажи** | | | |
| `orders_count` | INT | DEFAULT 0 | Количество заказов |
| `revenue` | DECIMAL(10,2) | DEFAULT 0 | Выручка |
| `avg_order_value` | DECIMAL(10,2) | DEFAULT 0 | Средний чек |
| `conversion_rate` | DECIMAL(5,2) | DEFAULT 0 | Конверсия в % |
| **Детальная аналитика** | | | |
| `payment_methods_usage` | JSONB | DEFAULT '{}' | Использование способов оплаты |
| `product_views` | INT | DEFAULT 0 | Просмотры товаров |
| `add_to_cart_count` | INT | DEFAULT 0 | Добавления в корзину |
| `checkout_count` | INT | DEFAULT 0 | Оформления заказов |
| `traffic_sources` | JSONB | DEFAULT '{}' | Источники трафика |
| `top_products` | JSONB | DEFAULT '[]' | Топ товаров |
| `top_categories` | JSONB | DEFAULT '[]' | Топ категорий |
| `orders_by_city` | JSONB | DEFAULT '{}' | Заказы по городам |

## Индексы

```sql
-- Основной поиск по витрине и дате
CREATE INDEX idx_storefront_analytics_storefront_date ON storefront_analytics(storefront_id, date DESC);

-- Поиск по датам для общей статистики
CREATE INDEX idx_storefront_analytics_date ON storefront_analytics(date);
```

## Ограничения

- **UNIQUE**: (`storefront_id`, `date`) - одна запись на витрину в день
- **FOREIGN KEY**: `storefront_id` → `storefronts(id)` ON DELETE CASCADE

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `storefront_id` → `storefronts.id` | Many-to-One | Витрина магазина |

## Структура JSONB полей

### payment_methods_usage
```json
{
  "card": 45,
  "paypal": 23,
  "crypto": 12,
  "bank_transfer": 8
}
```

### traffic_sources
```json
{
  "direct": 120,
  "google": 85,
  "facebook": 45,
  "instagram": 32,
  "referral": 18
}
```

### top_products
```json
[
  {"id": 123, "name": "iPhone 15", "views": 234, "sales": 15},
  {"id": 124, "name": "MacBook Pro", "views": 187, "sales": 8},
  {"id": 125, "name": "AirPods", "views": 156, "sales": 23}
]
```

### top_categories
```json
[
  {"id": 1, "name": "Electronics", "views": 450, "sales": 67},
  {"id": 2, "name": "Fashion", "views": 320, "sales": 45}
]
```

### orders_by_city
```json
{
  "Belgrade": 23,
  "Novi Sad": 12,
  "Nis": 8,
  "Kragujevac": 5
}
```

## Бизнес-правила

1. **Ежедневная агрегация**: Одна запись на витрину в день
2. **Расчетные поля**: `avg_order_value` и `conversion_rate` рассчитываются автоматически
3. **Каскадное удаление**: При удалении витрины удаляется вся аналитика
4. **Детализация в JSON**: Подробные данные хранятся в JSONB полях

## Примеры использования

### Добавление дневной статистики
```sql
INSERT INTO storefront_analytics (
    storefront_id, date, page_views, unique_visitors, orders_count, revenue
) VALUES (
    42, '2024-01-15', 250, 180, 8, 1240.50
) ON CONFLICT (storefront_id, date) 
DO UPDATE SET
    page_views = EXCLUDED.page_views,
    unique_visitors = EXCLUDED.unique_visitors,
    orders_count = EXCLUDED.orders_count,
    revenue = EXCLUDED.revenue,
    updated_at = CURRENT_TIMESTAMP;
```

### Получение статистики за месяц
```sql
SELECT 
    date,
    page_views,
    unique_visitors,
    orders_count,
    revenue,
    ROUND(conversion_rate, 2) as conversion_rate
FROM storefront_analytics 
WHERE storefront_id = 42 
  AND date >= '2024-01-01' 
  AND date <= '2024-01-31'
ORDER BY date DESC;
```

### Топ витрин по выручке за период
```sql
SELECT 
    s.name,
    SUM(sa.revenue) as total_revenue,
    SUM(sa.orders_count) as total_orders,
    AVG(sa.conversion_rate) as avg_conversion_rate
FROM storefront_analytics sa
JOIN storefronts s ON sa.storefront_id = s.id
WHERE sa.date >= '2024-01-01'
GROUP BY s.id, s.name
ORDER BY total_revenue DESC
LIMIT 10;
```

### Анализ источников трафика
```sql
SELECT 
    date,
    traffic_sources
FROM storefront_analytics 
WHERE storefront_id = 42 
  AND date >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY date DESC;
```

## Метрики и KPI

### Основные метрики
- **CTR**: Click-through rate на товары
- **Conversion Rate**: Процент посетителей, совершивших покупку
- **AOV**: Average Order Value - средний чек
- **Bounce Rate**: Процент отказов
- **Session Time**: Среднее время на сайте

### Расчетные формулы
```sql
-- Конверсия
conversion_rate = (orders_count * 100.0) / NULLIF(unique_visitors, 0)

-- Средний чек
avg_order_value = revenue / NULLIF(orders_count, 0)

-- CTR на товары
product_ctr = (add_to_cart_count * 100.0) / NULLIF(product_views, 0)
```

## Известные особенности

1. **Агрегированные данные**: Хранит уже обработанную статистику, не raw события
2. **Upsert логика**: Поддержка обновления существующих записей
3. **JSON аналитика**: Детальные данные в структурированном формате
4. **Временные ряды**: Оптимизировано для анализа трендов

## Использование в коде

**Backend**:
- Analytics service: `internal/proj/analytics/`
- Aggregation jobs: `internal/jobs/analytics_aggregator.go`
- API: `/api/v1/storefronts/{id}/analytics`

**Frontend**:
- Dashboard: `src/components/storefront/analytics/`
- Charts: `src/components/charts/`
- Страницы: `src/app/[locale]/storefront/[slug]/analytics/`

## Производительность

- **Партиционирование**: Рекомендуется по дате для больших объемов
- **Индексы**: Оптимизированы для временных диапазонов
- **Агрегация**: Batch-обработка в off-peak часы

## Связанные компоненты

- **Real-time events**: Система событий для сбора данных
- **Analytics dashboard**: Панель аналитики
- **Report generator**: Генерация отчетов
- **Alert system**: Уведомления о важных изменениях метрик