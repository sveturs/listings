# 📦 План внедрения системы мест хранения остатков товаров

## 📊 Результаты полного аудита текущей системы

### ✅ Что уже реализовано:
1. **Полноценная система управления остатками**
   - Учет остатков на уровне товаров и вариантов
   - Резервирование товаров при заказах
   - История движений товаров
   - Автоматическое освобождение резервов

2. **Интеграции**
   - Синхронизация с OpenSearch
   - ACID транзакции для предотвращения overselling
   - CSV импорт/экспорт остатков

3. **Frontend компоненты**
   - Таблицы управления остатками
   - Групповые операции
   - Аналитика остатков

### ❌ Что отсутствует:
- **Нет системы складов/мест хранения**
- **Все остатки в одном "виртуальном" месте**
- **Нет распределения по локациям**

## 🎯 Цель внедрения

Добавить поддержку множественных мест хранения остатков:
- **Витрина магазина** - товары на витрине для немедленной продажи
- **Склады логистических центров** - региональные склады партнеров
- **Склад маркетплейса** - централизованный склад платформы

## 📐 Архитектура решения

### 1. Структура базы данных

#### Новые таблицы:

```sql
-- Склады и места хранения
CREATE TABLE storage_locations (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL, -- 'storefront', 'warehouse', 'marketplace'
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL, -- Уникальный код склада
    
    -- Владелец
    owner_type VARCHAR(50) NOT NULL, -- 'storefront', 'marketplace', 'partner'
    owner_id INTEGER, -- ID витрины или партнера (NULL для маркетплейса)
    
    -- Адрес и контакты
    address TEXT NOT NULL,
    city VARCHAR(100),
    region VARCHAR(100),
    country VARCHAR(2) DEFAULT 'RS',
    postal_code VARCHAR(20),
    latitude NUMERIC(10,8),
    longitude NUMERIC(11,8),
    
    -- Контактная информация
    contact_name VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    
    -- Настройки
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false, -- Склад по умолчанию для владельца
    priority INTEGER DEFAULT 0, -- Приоритет для автовыбора
    
    -- Возможности склада
    capabilities JSONB DEFAULT '{}', -- {"can_ship": true, "can_pickup": true, "can_return": true}
    working_hours JSONB DEFAULT '{}', -- График работы
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT type_check CHECK (type IN ('storefront', 'warehouse', 'marketplace')),
    CONSTRAINT owner_type_check CHECK (owner_type IN ('storefront', 'marketplace', 'partner'))
);

-- Остатки по местам хранения
CREATE TABLE inventory_location_stock (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES storage_locations(id),
    product_id BIGINT NOT NULL REFERENCES storefront_products(id),
    variant_id BIGINT REFERENCES storefront_product_variants(id),
    
    -- Остатки
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    available_quantity GENERATED ALWAYS AS (quantity - reserved_quantity) STORED,
    
    -- Пороги
    min_stock_level INTEGER DEFAULT 0,
    max_stock_level INTEGER,
    reorder_point INTEGER,
    reorder_quantity INTEGER,
    
    -- Статус
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'blocked'
    
    -- Местоположение на складе
    zone VARCHAR(50), -- Зона склада
    rack VARCHAR(50), -- Стеллаж
    shelf VARCHAR(50), -- Полка
    bin VARCHAR(50), -- Ячейка
    
    -- Метаданные
    last_counted_at TIMESTAMP WITH TIME ZONE,
    last_movement_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(location_id, product_id, variant_id),
    CONSTRAINT quantity_check CHECK (quantity >= 0),
    CONSTRAINT reserved_check CHECK (reserved_quantity >= 0),
    CONSTRAINT reserved_not_greater CHECK (reserved_quantity <= quantity)
);

-- Перемещения между складами
CREATE TABLE inventory_transfers (
    id BIGSERIAL PRIMARY KEY,
    transfer_number VARCHAR(32) UNIQUE NOT NULL,
    
    -- Откуда и куда
    from_location_id BIGINT NOT NULL REFERENCES storage_locations(id),
    to_location_id BIGINT NOT NULL REFERENCES storage_locations(id),
    
    -- Статус перемещения
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending -> approved -> in_transit -> received/cancelled
    
    -- Товары
    product_id BIGINT NOT NULL REFERENCES storefront_products(id),
    variant_id BIGINT REFERENCES storefront_product_variants(id),
    quantity INTEGER NOT NULL,
    
    -- Даты
    requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP WITH TIME ZONE,
    shipped_at TIMESTAMP WITH TIME ZONE,
    received_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    
    -- Пользователи
    requested_by INTEGER NOT NULL REFERENCES users(id),
    approved_by INTEGER REFERENCES users(id),
    received_by INTEGER REFERENCES users(id),
    
    -- Дополнительно
    reason VARCHAR(255),
    notes TEXT,
    tracking_number VARCHAR(100),
    shipping_cost NUMERIC(10,2),
    
    CONSTRAINT quantity_positive CHECK (quantity > 0),
    CONSTRAINT different_locations CHECK (from_location_id != to_location_id)
);

-- Движения товаров по складам
CREATE TABLE inventory_location_movements (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES storage_locations(id),
    product_id BIGINT NOT NULL REFERENCES storefront_products(id),
    variant_id BIGINT REFERENCES storefront_product_variants(id),
    
    -- Тип движения
    type VARCHAR(20) NOT NULL, -- 'in', 'out', 'adjustment', 'transfer_in', 'transfer_out'
    quantity INTEGER NOT NULL,
    
    -- Связанные объекты
    order_id BIGINT REFERENCES storefront_orders(id),
    transfer_id BIGINT REFERENCES inventory_transfers(id),
    
    -- Детали
    reason VARCHAR(100),
    notes TEXT,
    reference_number VARCHAR(100), -- Внешний номер документа
    
    -- Пользователь
    user_id INTEGER NOT NULL REFERENCES users(id),
    
    -- Время
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Правила распределения остатков
CREATE TABLE inventory_allocation_rules (
    id BIGSERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id),
    
    -- Правило
    rule_type VARCHAR(50) NOT NULL, -- 'priority', 'nearest', 'cheapest', 'fastest'
    priority INTEGER DEFAULT 0,
    
    -- Условия
    conditions JSONB DEFAULT '{}', -- {"min_quantity": 10, "max_distance": 50}
    
    -- Локации
    preferred_locations INTEGER[], -- Массив ID предпочтительных складов
    excluded_locations INTEGER[], -- Массив ID исключенных складов
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### Изменения в существующих таблицах:

```sql
-- Добавить в storefront_orders
ALTER TABLE storefront_orders ADD COLUMN fulfillment_location_id BIGINT REFERENCES storage_locations(id);
ALTER TABLE storefront_orders ADD COLUMN pickup_location_id BIGINT REFERENCES storage_locations(id);

-- Добавить в storefront_order_items
ALTER TABLE storefront_order_items ADD COLUMN allocated_location_id BIGINT REFERENCES storage_locations(id);

-- Добавить в inventory_reservations
ALTER TABLE inventory_reservations ADD COLUMN location_id BIGINT REFERENCES storage_locations(id);
```

### 2. Backend реализация

#### Новые сервисы:

**LocationService** - управление складами:
```go
type LocationService interface {
    CreateLocation(ctx context.Context, location *StorageLocation) error
    UpdateLocation(ctx context.Context, id int64, location *StorageLocation) error
    GetLocation(ctx context.Context, id int64) (*StorageLocation, error)
    ListLocations(ctx context.Context, filters LocationFilters) ([]*StorageLocation, error)
    GetDefaultLocation(ctx context.Context, ownerType string, ownerID int) (*StorageLocation, error)
}
```

**MultiLocationInventoryManager** - управление остатками по локациям:
```go
type MultiLocationInventoryManager interface {
    // Получение остатков
    GetStockByLocation(ctx context.Context, productID, variantID, locationID int64) (*LocationStock, error)
    GetTotalStock(ctx context.Context, productID, variantID int64) (*AggregatedStock, error)
    GetStockAllLocations(ctx context.Context, productID, variantID int64) ([]*LocationStock, error)
    
    // Обновление остатков
    UpdateLocationStock(ctx context.Context, locationID int64, updates []StockUpdate) error
    TransferStock(ctx context.Context, transfer *StockTransfer) error
    
    // Резервирование
    ReserveStockMultiLocation(ctx context.Context, orderID int64, items []OrderItem) ([]Reservation, error)
    AllocateStock(ctx context.Context, orderID int64, rules *AllocationRules) ([]Allocation, error)
    
    // Аналитика
    GetLowStockLocations(ctx context.Context, threshold int) ([]*LowStockReport, error)
    GetStockDistribution(ctx context.Context, productID int64) (*DistributionReport, error)
}
```

**TransferService** - перемещения между складами:
```go
type TransferService interface {
    CreateTransfer(ctx context.Context, transfer *Transfer) error
    ApproveTransfer(ctx context.Context, transferID int64, userID int) error
    ShipTransfer(ctx context.Context, transferID int64, trackingNumber string) error
    ReceiveTransfer(ctx context.Context, transferID int64, userID int) error
    CancelTransfer(ctx context.Context, transferID int64, reason string) error
    ListTransfers(ctx context.Context, filters TransferFilters) ([]*Transfer, error)
}
```

#### API Endpoints:

```yaml
# Управление складами
POST   /api/v1/storefronts/{slug}/locations           # Создать склад
GET    /api/v1/storefronts/{slug}/locations           # Список складов
GET    /api/v1/storefronts/{slug}/locations/{id}      # Детали склада
PUT    /api/v1/storefronts/{slug}/locations/{id}      # Обновить склад
DELETE /api/v1/storefronts/{slug}/locations/{id}      # Удалить склад

# Остатки по локациям
GET    /api/v1/storefronts/{slug}/products/{id}/stock/locations      # Остатки по всем складам
GET    /api/v1/storefronts/{slug}/products/{id}/stock/location/{lid} # Остатки на конкретном складе
PUT    /api/v1/storefronts/{slug}/products/{id}/stock/location/{lid} # Обновить остатки на складе

# Перемещения
POST   /api/v1/storefronts/{slug}/transfers           # Создать перемещение
GET    /api/v1/storefronts/{slug}/transfers           # Список перемещений
GET    /api/v1/storefronts/{slug}/transfers/{id}      # Детали перемещения
PUT    /api/v1/storefronts/{slug}/transfers/{id}/approve  # Утвердить
PUT    /api/v1/storefronts/{slug}/transfers/{id}/ship     # Отправить
PUT    /api/v1/storefronts/{slug}/transfers/{id}/receive  # Принять
PUT    /api/v1/storefronts/{slug}/transfers/{id}/cancel   # Отменить

# Правила распределения
GET    /api/v1/storefronts/{slug}/allocation-rules    # Правила распределения
POST   /api/v1/storefronts/{slug}/allocation-rules    # Создать правило
PUT    /api/v1/storefronts/{slug}/allocation-rules/{id} # Обновить правило

# Аналитика
GET    /api/v1/storefronts/{slug}/inventory/distribution  # Распределение остатков
GET    /api/v1/storefronts/{slug}/inventory/low-stock     # Товары с низкими остатками по складам
```

### 3. Frontend компоненты

#### Новые компоненты:

**LocationManager** - управление складами:
```tsx
components/Storefront/Inventory/LocationManager.tsx
- Список складов с фильтрами
- Форма создания/редактирования склада
- Карта с расположением складов
- Настройки склада (график, возможности)
```

**MultiLocationStockTable** - остатки по складам:
```tsx
components/Storefront/Inventory/MultiLocationStockTable.tsx
- Таблица остатков с группировкой по складам
- Фильтры по локациям
- Быстрое редактирование остатков
- Визуализация распределения
```

**TransferManager** - перемещения:
```tsx
components/Storefront/Inventory/TransferManager.tsx
- Создание заявки на перемещение
- Список активных перемещений
- Статусы и трекинг
- История перемещений
```

**StockDistributionChart** - визуализация:
```tsx
components/Storefront/Inventory/StockDistributionChart.tsx
- График распределения по складам
- Тепловая карта остатков
- Аналитика по локациям
```

**AllocationRulesEditor** - правила распределения:
```tsx
components/Storefront/Inventory/AllocationRulesEditor.tsx
- Настройка правил автораспределения
- Приоритеты складов
- Условия выбора склада
```

#### Изменения в существующих компонентах:

1. **UnifiedProductCard** - показывать общие остатки или остатки ближайшего склада
2. **VariantStockTable** - добавить колонку с распределением по складам
3. **AddToCartButton** - учитывать доступность на разных складах
4. **CheckoutPage** - выбор склада для самовывоза или доставки

### 4. Миграция данных

```sql
-- Миграция 001: Создание таблиц
-- (SQL из раздела "Структура базы данных")

-- Миграция 002: Создание склада по умолчанию для каждой витрины
INSERT INTO storage_locations (
    type, name, code, owner_type, owner_id,
    address, city, is_default, is_active
)
SELECT 
    'storefront',
    s.name || ' - Основной склад',
    'STF-' || s.id,
    'storefront',
    s.id,
    s.address,
    COALESCE(s.city, 'Белград'),
    true,
    true
FROM storefronts s;

-- Миграция 003: Перенос существующих остатков
INSERT INTO inventory_location_stock (
    location_id, product_id, variant_id, quantity, reserved_quantity
)
SELECT 
    sl.id,
    sp.id,
    NULL,
    sp.stock_quantity,
    COALESCE(
        (SELECT SUM(ir.quantity) 
         FROM inventory_reservations ir 
         WHERE ir.product_id = sp.id 
         AND ir.variant_id IS NULL 
         AND ir.status = 'active'),
        0
    )
FROM storefront_products sp
JOIN storage_locations sl ON sl.owner_id = sp.storefront_id 
    AND sl.owner_type = 'storefront' 
    AND sl.is_default = true
WHERE sp.has_variants = false;

-- Перенос остатков вариантов
INSERT INTO inventory_location_stock (
    location_id, product_id, variant_id, quantity, reserved_quantity
)
SELECT 
    sl.id,
    spv.product_id,
    spv.id,
    spv.stock_quantity,
    COALESCE(
        (SELECT SUM(ir.quantity) 
         FROM inventory_reservations ir 
         WHERE ir.variant_id = spv.id 
         AND ir.status = 'active'),
        0
    )
FROM storefront_product_variants spv
JOIN storefront_products sp ON sp.id = spv.product_id
JOIN storage_locations sl ON sl.owner_id = sp.storefront_id 
    AND sl.owner_type = 'storefront' 
    AND sl.is_default = true;
```

## 📅 План внедрения по этапам

### Этап 1: Базовая инфраструктура (5 дней)
**День 1-2: База данных**
- [ ] Создать миграции для новых таблиц
- [ ] Выполнить миграцию существующих данных
- [ ] Создать индексы для оптимизации

**День 3-4: Backend сервисы**
- [ ] Реализовать LocationService
- [ ] Создать базовые CRUD операции для складов
- [ ] Добавить API endpoints для управления складами

**День 5: Тестирование**
- [ ] Unit тесты для новых сервисов
- [ ] Интеграционные тесты API

### Этап 2: Управление остатками по локациям (5 дней)
**День 6-7: Backend**
- [ ] Реализовать MultiLocationInventoryManager
- [ ] Модифицировать существующий InventoryManager
- [ ] Добавить поддержку множественных локаций в заказах

**День 8-9: Frontend**
- [ ] Создать LocationManager компонент
- [ ] Реализовать MultiLocationStockTable
- [ ] Интегрировать в панель управления витриной

**День 10: Интеграция**
- [ ] Обновить процесс создания заказа
- [ ] Модифицировать резервирование товаров
- [ ] Тестирование сценариев

### Этап 3: Перемещения между складами (4 дня)
**День 11-12: Backend**
- [ ] Реализовать TransferService
- [ ] Создать API для перемещений
- [ ] Добавить обработку статусов перемещений

**День 13-14: Frontend**
- [ ] Создать TransferManager компонент
- [ ] Добавить формы создания перемещений
- [ ] Реализовать отслеживание статусов

### Этап 4: Автоматическое распределение (3 дня)
**День 15-16: Логика распределения**
- [ ] Реализовать правила распределения
- [ ] Создать алгоритмы выбора оптимального склада
- [ ] Добавить API для настройки правил

**День 17: Frontend**
- [ ] Создать AllocationRulesEditor
- [ ] Интегрировать в настройки витрины

### Этап 5: Аналитика и отчеты (3 дня)
**День 18-19: Backend**
- [ ] Создать сервисы аналитики
- [ ] Реализовать отчеты по складам
- [ ] API для получения статистики

**День 20: Frontend**
- [ ] Создать StockDistributionChart
- [ ] Добавить дашборд с аналитикой
- [ ] Интегрировать отчеты

### Этап 6: Финализация (3 дня)
**День 21: Оптимизация**
- [ ] Оптимизация запросов
- [ ] Кеширование данных
- [ ] Улучшение производительности

**День 22: Документация**
- [ ] API документация
- [ ] Руководство пользователя
- [ ] Обучающие материалы

**День 23: Развертывание**
- [ ] Подготовка production среды
- [ ] Поэтапное развертывание
- [ ] Мониторинг и поддержка

## 🚀 Преимущества внедрения

1. **Гибкость управления**
   - Распределенное хранение товаров
   - Оптимизация логистики
   - Снижение затрат на доставку

2. **Улучшенный контроль**
   - Точный учет по локациям
   - История перемещений
   - Предотвращение потерь

3. **Масштабируемость**
   - Легкое добавление новых складов
   - Интеграция с партнерами
   - Поддержка dropshipping

4. **Автоматизация**
   - Автовыбор оптимального склада
   - Автоматическое пополнение
   - Умное распределение заказов

## 🔧 Технические требования

- PostgreSQL 14+
- Go 1.21+
- React 18+
- Next.js 15+
- Redis для кеширования

## 📝 Риски и митигация

1. **Сложность миграции**
   - Решение: Поэтапная миграция с откатом
   
2. **Производительность**
   - Решение: Индексы, кеширование, оптимизация запросов
   
3. **Обучение пользователей**
   - Решение: Постепенное внедрение, обучающие материалы

## 💰 Оценка трудозатрат

- **Backend разработка**: 12 дней
- **Frontend разработка**: 8 дней
- **Тестирование**: 3 дня
- **Документация**: 2 дня
- **Развертывание**: 1 день

**Итого**: 26 рабочих дней (1 разработчик full-stack)

## 🎯 KPI успешности внедрения

1. Снижение времени обработки заказа на 30%
2. Уменьшение затрат на логистику на 20%
3. Повышение точности учета до 99.9%
4. Сокращение случаев overselling до 0
5. Увеличение скорости инвентаризации в 3 раза