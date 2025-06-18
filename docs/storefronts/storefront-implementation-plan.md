# План внедрения витрин в Sve Tu Platform

## Концепция "Лучшие витрины на рынке"

Наша цель - создать систему витрин, которая превзойдет конкурентов по функциональности, удобству и возможностям монетизации.

### Географический фокус
- **Первичный рынок**: Сербия
- **Расширение**: Страны бывшей Югославии (Хорватия, Босния и Герцеговина, Северная Македония, Словения, Черногория)
- **Локализация**: Полная поддержка местных языков, валют и платежных систем

## Архитектура системы

### Преимущества перед конкурентами

1. **Умный импорт товаров**
   - AI-распознавание категорий
   - Автоматическое улучшение описаний
   - Оптимизация изображений

2. **Продвинутая аналитика**
   - Тепловые карты кликов
   - A/B тестирование витрин
   - Прогнозирование продаж

3. **Омниканальность**
   - Интеграция с социальными сетями
   - API для внешних систем
   - Мобильное приложение для управления

## Фаза 1: Базовая инфраструктура (2 недели)

### Backend задачи

#### 1.1 Модели данных
```sql
-- Таблица витрин
CREATE TABLE storefronts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    theme JSONB DEFAULT '{}',
    
    -- Контакты
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    
    -- Адрес и геолокация
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(2),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Настройки
    settings JSONB DEFAULT '{}',
    seo_meta JSONB DEFAULT '{}',
    
    -- Статус и статистика
    is_active BOOLEAN DEFAULT false,
    is_verified BOOLEAN DEFAULT false,
    rating DECIMAL(3, 2) DEFAULT 0,
    reviews_count INT DEFAULT 0,
    views_count INT DEFAULT 0,
    
    -- Подписка
    subscription_plan VARCHAR(50) DEFAULT 'free',
    subscription_expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Связь товаров с витринами
ALTER TABLE marketplace_listings 
ADD COLUMN storefront_id INT REFERENCES storefronts(id),
ADD COLUMN is_storefront_item BOOLEAN DEFAULT false;

-- Таблица сотрудников витрины
CREATE TABLE storefront_staff (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id),
    user_id INT NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL, -- owner, manager, support
    permissions JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(storefront_id, user_id)
);

-- Таблица рабочих часов
CREATE TABLE storefront_hours (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id),
    day_of_week INT NOT NULL, -- 0-6
    open_time TIME,
    close_time TIME,
    is_closed BOOLEAN DEFAULT false
);

-- Таблица способов оплаты
CREATE TABLE storefront_payment_methods (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id),
    method_type VARCHAR(50) NOT NULL, -- cash, card, online, crypto
    is_enabled BOOLEAN DEFAULT true,
    settings JSONB DEFAULT '{}'
);

-- Таблица доставки
CREATE TABLE storefront_delivery_options (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    min_order_amount DECIMAL(10, 2),
    estimated_days INT,
    zones JSONB DEFAULT '[]', -- географические зоны доставки
    is_active BOOLEAN DEFAULT true
);
```

#### 1.2 API endpoints
```go
// Управление витринами
POST   /api/v1/storefronts              // Создать витрину
GET    /api/v1/storefronts              // Список моих витрин
GET    /api/v1/storefronts/:id          // Детали витрины
PUT    /api/v1/storefronts/:id          // Обновить витрину
DELETE /api/v1/storefronts/:id          // Удалить витрину

// Товары витрины
GET    /api/v1/storefronts/:id/listings // Товары витрины
POST   /api/v1/storefronts/:id/listings // Добавить товар
PUT    /api/v1/storefronts/:id/listings/batch // Групповые операции

// Публичные endpoints
GET    /api/v1/public/storefronts/:slug // Публичная витрина
GET    /api/v1/public/storefronts/:slug/listings // Товары
GET    /api/v1/public/storefronts/search // Поиск витрин
```

### Frontend задачи

#### 1.3 Страницы и компоненты
```typescript
// Структура страниц
/storefronts                    // Список моих витрин
/storefronts/create            // Создание витрины
/storefronts/[id]              // Панель управления
/storefronts/[id]/products     // Управление товарами
/storefronts/[id]/orders       // Заказы
/storefronts/[id]/analytics    // Аналитика
/storefronts/[id]/settings     // Настройки

/shop/[slug]                   // Публичная витрина
/shop/[slug]/products          // Каталог товаров
/shop/[slug]/about            // О магазине
/shop/[slug]/reviews          // Отзывы
```

#### 1.4 Redux slices
```typescript
// storefrontSlice.ts
interface StorefrontState {
  myStorefronts: Storefront[];
  currentStorefront: Storefront | null;
  publicStorefront: PublicStorefront | null;
  loading: boolean;
  error: string | null;
}

// storefrontProductsSlice.ts
interface StorefrontProductsState {
  products: Product[];
  filters: ProductFilters;
  pagination: Pagination;
  loading: boolean;
}
```

## Фаза 2: Расширенный функционал (3 недели)

### 2.1 Система импорта товаров

#### Backend
```go
// Модель импорта
type ImportJob struct {
    ID           int
    StorefrontID int
    FileType     string // csv, xml, json, api
    SourceURL    string
    Status       string // pending, processing, completed, failed
    TotalItems   int
    ProcessedItems int
    Errors       []ImportError
    Mapping      map[string]string // соответствие полей
    CreatedAt    time.Time
}

// API endpoints
POST /api/v1/storefronts/:id/import/upload   // Загрузка файла
POST /api/v1/storefronts/:id/import/url      // Импорт по URL
GET  /api/v1/storefronts/:id/import/history  // История импортов
GET  /api/v1/storefronts/:id/import/:jobId   // Статус импорта
```

#### Frontend компоненты
- ImportWizard - мастер импорта с шагами
- FieldMapper - сопоставление полей
- ImportProgress - прогресс импорта
- ImportHistory - история и повтор импортов

### 2.2 Продвинутая настройка витрины

#### Темы оформления
```typescript
interface StorefrontTheme {
  primaryColor: string;
  secondaryColor: string;
  fontFamily: string;
  layout: 'grid' | 'list' | 'masonry';
  headerStyle: 'minimal' | 'full' | 'transparent';
  customCSS?: string;
}
```

#### SEO оптимизация
- Метатеги для каждой страницы
- Структурированные данные (Schema.org)
- Sitemap для витрины
- Open Graph теги

### 2.3 Аналитика и отчеты

#### Метрики
- Просмотры витрины/товаров
- Конверсия в заказы
- Средний чек
- География покупателей
- Источники трафика

#### Дашборд
- Графики продаж
- Топ товаров
- Карта покупателей
- Сравнение периодов

## Фаза 3: Интеграция оплаты и доставки (4 недели)

### 3.1 Система заказов

#### Модель данных
```sql
CREATE TABLE storefront_orders (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES storefronts(id),
    customer_id INT REFERENCES users(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    
    -- Статус
    status VARCHAR(50) NOT NULL, -- pending, confirmed, processing, shipped, delivered, cancelled
    payment_status VARCHAR(50), -- pending, paid, refunded
    
    -- Суммы
    subtotal DECIMAL(10, 2) NOT NULL,
    delivery_fee DECIMAL(10, 2) DEFAULT 0,
    discount DECIMAL(10, 2) DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    
    -- Доставка
    delivery_option_id INT REFERENCES storefront_delivery_options(id),
    delivery_address JSONB,
    tracking_number VARCHAR(255),
    
    -- Оплата
    payment_method VARCHAR(50),
    payment_details JSONB,
    
    -- Контакты
    customer_name VARCHAR(255),
    customer_phone VARCHAR(50),
    customer_email VARCHAR(255),
    
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE storefront_order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES storefront_orders(id),
    listing_id INT NOT NULL REFERENCES marketplace_listings(id),
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    total DECIMAL(10, 2) NOT NULL
);
```

### 3.2 Интеграция платежных систем

#### Поддерживаемые системы (адаптировано для Балкан)
1. **Локальные банки Сербии**
   - **Poštanska štedionica** (основной банк-партнер) ✅
   - Intesa Sanpaolo
   - Raiffeisen Bank
   - UniCredit Bank
   - Платежные карты DinaCard

2. **Региональные системы**
   - **PayTen** - популярно в Сербии
   - **CorvusPay** - Хорватия и регион
   - **MonriPayments** - вся бывшая Югославия
   
3. **Международные**
   - **PayPal** - где доступно
   - **Stripe** - для международных продавцов
   - **Wise** - переводы между странами

4. **Наличные**
   - Оплата при доставке (очень популярно)
   - Оплата в пунктах выдачи

#### Процесс оплаты
1. Покупатель оформляет заказ
2. Выбирает способ оплаты
3. Перенаправление на платежный шлюз
4. Подтверждение оплаты
5. Автоматическое обновление статуса

### 3.3 Система доставки

#### Интеграции для региона Балкан
1. **Почта Сербии (Pošta Srbije)** - основной партнер
   - Доставка по всей Сербии
   - Post Express для быстрой доставки
   - Международные отправления

2. **Курьерские службы**
   - **AKS** - популярно в Сербии
   - **BEX** - быстрая доставка
   - **City Express** - городская доставка
   - **D Express** - региональная сеть

3. **Региональные партнеры**
   - **Hrvatska pošta** - Хорватия
   - **Pošta Slovenije** - Словения
   - **Makedonska pošta** - Северная Македония

4. **Самовывоз**
   - Пункты выдачи в крупных городах
   - Партнерские магазины
   - Локеры (автоматы выдачи)

#### Функционал
- Расчет стоимости доставки
- Отслеживание посылок
- Печать этикеток
- Управление возвратами

## Фаза 4: Продвинутые возможности (2 недели)

### 4.1 AI-функции

#### Умные рекомендации
- Персонализированные товары
- Кросс-продажи
- Похожие витрины

#### Автоматизация
- Генерация описаний товаров
- Оптимизация цен
- Ответы на частые вопросы

### 4.2 Маркетинговые инструменты

#### Промо-кампании
- Скидки и купоны
- Временные акции
- Программа лояльности
- Реферальная система

#### Email-маркетинг
- Автоматические рассылки
- Брошенные корзины
- Новости и акции
- Персонализация

### 4.3 Мобильное приложение

#### Для продавцов
- Управление витриной
- Обработка заказов
- Чат с покупателями
- Push-уведомления

#### Для покупателей
- Избранные витрины
- История покупок
- Отслеживание заказов
- Быстрый заказ

## Технические требования

### Производительность
- Загрузка витрины < 2 сек
- Поддержка 10000+ товаров
- CDN для изображений
- Оптимизация запросов

### Безопасность
- SSL для всех витрин
- PCI DSS для платежей
- Защита от DDoS
- Резервное копирование

### Масштабируемость
- Микросервисная архитектура
- Горизонтальное масштабирование
- Кеширование (Redis)
- Очереди (RabbitMQ)

## Монетизация

### Тарифные планы (адаптировано для рынка Сербии)

#### Starter (Бесплатно)
- До 50 товаров
- Базовая витрина
- 3% комиссия с продаж

#### Professional (2,990 RSD/мес ~ €25)
- До 500 товаров
- Темы оформления
- Импорт товаров
- 2% комиссия

#### Business (9,990 RSD/мес ~ €85)
- Неограниченно товаров
- API доступ
- Приоритетная поддержка
- 1% комиссия

#### Enterprise (Индивидуально)
- White-label решение
- Выделенный менеджер
- Кастомизация
- От 0.5% комиссии

### Дополнительные услуги
- Продвижение витрины - 5,990 RSD/мес (~€50)
- SMS-уведомления - 6 RSD/шт (~€0.05)
- Дополнительные менеджеры - 1,190 RSD/мес (~€10)
- Выделенная поддержка на сербском языке

## График внедрения

### Месяц 1
- Неделя 1-2: Фаза 1 (Базовая инфраструктура)
- Неделя 3-4: Начало Фазы 2 (Импорт и настройки)

### Месяц 2
- Неделя 1: Завершение Фазы 2
- Неделя 2-4: Фаза 3 (Оплата и доставка)

### Месяц 3
- Неделя 1: Завершение Фазы 3
- Неделя 2-3: Фаза 4 (Продвинутые функции)
- Неделя 4: Тестирование и запуск

## KPI успеха

### Через 3 месяца (Сербия)
- 100+ активных витрин
- 10,000+ товаров в витринах
- 1,000+ заказов через витрины
- 1,200,000 RSD (~€10,000) месячный оборот

### Через 6 месяцев (Сербия + 1-2 страны)
- 500+ активных витрин
- 50,000+ товаров
- Запуск в Хорватии и Боснии
- 6,000,000 RSD (~€50,000) месячный оборот

### Через 1 год (Вся бывшая Югославия)
- 2,000+ витрин
- 200,000+ товаров
- Присутствие во всех целевых странах
- 12,000,000 RSD (~€100,000) месячный оборот

## Заключение

Система витрин Sve Tu Platform станет лучшей на рынке Балкан благодаря:

1. **Локальной адаптации**
   - Интеграция с местными банками и платежными системами
   - Поддержка всех языков региона (сербский, хорватский, боснийский, словенский, македонский)
   - Партнерство с локальными службами доставки

2. **Понимании рынка**
   - Популярность оплаты при доставке
   - Важность личных контактов и доверия
   - Предпочтение локальных продавцов

3. **Региональной экспансии**
   - Единая платформа для всей бывшей Югославии
   - Упрощенная трансграничная торговля
   - Общая культура и язык как преимущество

4. **Конкурентных ценах**
   - Адаптированные тарифы для каждой страны
   - Низкие комиссии по сравнению с международными платформами
   - Поддержка малого бизнеса

Это создаст экосистему, объединяющую продавцов и покупателей всего региона, возрождая экономические связи между странами бывшей Югославии через современную цифровую платформу.