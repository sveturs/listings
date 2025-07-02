# Паспорт таблицы `storefronts`

## Назначение
Основная таблица витрин магазинов - новая полнофункциональная система для создания персональных интернет-магазинов пользователями. Поддерживает брендинг, SEO, аналитику, подписки и AI-функции.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS storefronts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Branding
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    theme JSONB DEFAULT '{"primaryColor": "#1976d2", "layout": "grid"}',
    
    -- Contact information
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    
    -- Location
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'RS',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Business settings
    settings JSONB DEFAULT '{}',
    seo_meta JSONB DEFAULT '{}',
    
    -- Status and stats
    is_active BOOLEAN DEFAULT false,
    is_verified BOOLEAN DEFAULT false,
    verification_date TIMESTAMP,
    rating DECIMAL(3, 2) DEFAULT 0.00,
    reviews_count INT DEFAULT 0,
    products_count INT DEFAULT 0,
    sales_count INT DEFAULT 0,
    views_count INT DEFAULT 0,
    
    -- Subscription (for monetization)
    subscription_plan VARCHAR(50) DEFAULT 'starter',
    subscription_expires_at TIMESTAMP,
    commission_rate DECIMAL(5, 2) DEFAULT 3.00,
    
    -- AI and killer features preparation
    ai_agent_enabled BOOLEAN DEFAULT false,
    ai_agent_config JSONB DEFAULT '{}',
    live_shopping_enabled BOOLEAN DEFAULT false,
    group_buying_enabled BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID витрины |
| `user_id` | INT | NOT NULL FK | ID владельца витрины |
| `slug` | VARCHAR(100) | UNIQUE NOT NULL | URL-slug для витрины |
| `name` | VARCHAR(255) | NOT NULL | Название магазина |
| `description` | TEXT | NULLABLE | Описание магазина |
| **Брендинг** | | | |
| `logo_url` | VARCHAR(500) | NULLABLE | URL логотипа |
| `banner_url` | VARCHAR(500) | NULLABLE | URL баннера |
| `theme` | JSONB | DEFAULT | Тема оформления |
| **Контакты** | | | |
| `phone` | VARCHAR(50) | NULLABLE | Телефон |
| `email` | VARCHAR(255) | NULLABLE | Email |
| `website` | VARCHAR(255) | NULLABLE | Веб-сайт |
| **Местоположение** | | | |
| `address` | TEXT | NULLABLE | Адрес |
| `city` | VARCHAR(100) | NULLABLE | Город |
| `postal_code` | VARCHAR(20) | NULLABLE | Почтовый индекс |
| `country` | VARCHAR(2) | DEFAULT 'RS' | Код страны |
| `latitude` | DECIMAL(10,8) | NULLABLE | Широта |
| `longitude` | DECIMAL(11,8) | NULLABLE | Долгота |
| **Настройки** | | | |
| `settings` | JSONB | DEFAULT '{}' | Настройки магазина |
| `seo_meta` | JSONB | DEFAULT '{}' | SEO метаданные |
| **Статистика** | | | |
| `is_active` | BOOLEAN | DEFAULT false | Активность витрины |
| `is_verified` | BOOLEAN | DEFAULT false | Верификация |
| `verification_date` | TIMESTAMP | NULLABLE | Дата верификации |
| `rating` | DECIMAL(3,2) | DEFAULT 0.00 | Рейтинг (0.00-5.00) |
| `reviews_count` | INT | DEFAULT 0 | Количество отзывов |
| `products_count` | INT | DEFAULT 0 | Количество товаров |
| `sales_count` | INT | DEFAULT 0 | Количество продаж |
| `views_count` | INT | DEFAULT 0 | Количество просмотров |
| **Монетизация** | | | |
| `subscription_plan` | VARCHAR(50) | DEFAULT 'starter' | План подписки |
| `subscription_expires_at` | TIMESTAMP | NULLABLE | Окончание подписки |
| `commission_rate` | DECIMAL(5,2) | DEFAULT 3.00 | Комиссия в % |
| **AI и продвинутые функции** | | | |
| `ai_agent_enabled` | BOOLEAN | DEFAULT false | AI-агент включен |
| `ai_agent_config` | JSONB | DEFAULT '{}' | Конфигурация AI |
| `live_shopping_enabled` | BOOLEAN | DEFAULT false | Лайв-шоппинг |
| `group_buying_enabled` | BOOLEAN | DEFAULT false | Групповые покупки |

## Индексы

```sql
CREATE INDEX idx_storefronts_user_id ON storefronts(user_id);
CREATE INDEX idx_storefronts_slug ON storefronts(slug);
CREATE INDEX idx_storefronts_city ON storefronts(city);
CREATE INDEX idx_storefronts_is_active ON storefronts(is_active);
CREATE INDEX idx_storefronts_rating ON storefronts(rating DESC);
```

## Ограничения

- **UNIQUE**: `slug` - уникальный URL каждой витрины
- **FOREIGN KEY**: `user_id` → `users(id)` ON DELETE CASCADE
- **CHECK**: `rating` должен быть от 0.00 до 5.00
- **CHECK**: `commission_rate` должен быть положительным

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `user_id` → `users.id` | Many-to-One | Владелец витрины |
| `storefront_products` | One-to-Many | Товары витрины |
| `storefront_analytics` | One-to-Many | Аналитика витрины |

## Планы подписки

| План | Описание | Возможности |
|------|----------|-------------|
| `starter` | Бесплатный | Базовые функции |
| `professional` | Платный | Расширенная аналитика |
| `business` | Премиум | AI-функции |
| `enterprise` | Корпоративный | Все возможности |

## Бизнес-правила

1. **Уникальность slug**: Каждая витрина имеет уникальный URL
2. **Верификация**: Влияет на доверие покупателей
3. **Монетизация**: Комиссия с продаж зависит от плана
4. **AI-функции**: Доступны на продвинутых планах
5. **Геолокация**: Поддержка поиска по местоположению

## Примеры использования

### Создание новой витрины
```sql
INSERT INTO storefronts (user_id, slug, name, description, theme)
VALUES (123, 'electronics-pro', 'Electronics Pro', 'Лучшая электроника', 
        '{"primaryColor": "#FF6B35", "layout": "list"}');
```

### Поиск активных витрин в городе
```sql
SELECT * FROM storefronts 
WHERE city = 'Belgrade' AND is_active = true 
ORDER BY rating DESC, views_count DESC;
```

### Обновление статистики
```sql
UPDATE storefronts 
SET products_count = products_count + 1,
    views_count = views_count + 1
WHERE id = 42;
```

## Известные особенности

1. **Новая версия**: Заменяет устаревшую таблицу `user_storefronts`
2. **Богатая функциональность**: Поддержка брендинга, SEO, аналитики
3. **Подготовка к AI**: Поля для будущих AI-функций уже включены
4. **Монетизация**: Встроенная система подписок и комиссий
5. **Геолокация**: Поддержка поиска по местоположению

## Использование в коде

**Backend**:
- Handler: `internal/proj/storefronts/`
- Модели: `internal/domain/storefront.go`
- API: `/api/v1/storefronts/`

**Frontend**:
- Компоненты: `src/components/storefront/`
- Страницы: `src/app/[locale]/storefront/`
- Типы: `@/types/storefront.ts`

## Связанные компоненты

- **Storefront products**: Товары витрины
- **Storefront analytics**: Аналитика продаж
- **Payment system**: Обработка платежей
- **SEO система**: Оптимизация для поисковиков
- **Theme engine**: Система тем и брендинга