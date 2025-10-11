# Анализ: MarketplaceListing vs StorefrontProduct

**Дата:** 2025-10-09
**Статус:** Исследование возможности объединения сущностей
**Автор:** Анализ архитектуры проекта

---

## 📋 Содержание

1. [Executive Summary](#executive-summary)
2. [Текущее состояние](#текущее-состояние)
3. [Количественный анализ полей](#количественный-анализ-полей)
4. [Качественный анализ](#качественный-анализ)
5. [Доказательства исторического расхождения](#доказательства-исторического-расхождения)
6. [Преимущества объединения](#преимущества-объединения)
7. [Недостатки объединения](#недостатки-объединения)
8. [Рекомендации](#рекомендации)
9. [План миграции (если принято решение объединять)](#план-миграции)

---

## Executive Summary

### Ключевые выводы:

- ✅ **Общих полей больше, чем различий**: 22 общих поля (41%) vs 31 уникальное (59%)
- ✅ **Разные названия полей** - результат технического долга, а не архитектурного решения
- ✅ **Функция конвертации уже существует**: `getStorefrontProductAsListing()` доказывает, что система работает с ними как с одним типом
- ⚠️ **Критичные различия**: владение товаром (UserID vs StorefrontID), статус, атрибуты
- 📊 **58% NULL полей** при объединении в одну таблицу

### Гипотеза подтверждена:

> **"Разные типы у одних названий - это следствие того, что мы поменяли в одной сущности и забыли в другой"**

Анализ кода показывает, что это действительно **исторический технический долг**.

---

## Текущее состояние

### MarketplaceListing
**Файл:** `backend/internal/domain/models/models.go:40-94`
**Назначение:** P2P маркетплейс - объявления от частных лиц и компаний
**Модель:** "Доска объявлений" (Avito, OLX, Facebook Marketplace)

**Характеристики:**
- Владелец: `UserID int` (физлицо/юрлицо)
- Статус: `Status string` (draft/active/sold/archived/moderation)
- Социальные функции: голосования, отзывы
- Гибкая категоризация

### StorefrontProduct
**Файл:** `backend/internal/domain/models/storefront_product.go:8-48`
**Назначение:** Товар витрины - продукт в интернет-магазине
**Модель:** "Магазин" (Wildberries, OZON)

**Характеристики:**
- Владелец: `StorefrontID int` (магазин)
- Статус: `IsActive bool` (активен/неактивен)
- Профессиональная торговля: складской учет, варианты товара
- B2C модель

### Интеграция

```go
// backend/internal/proj/marketplace/storage/postgres/marketplace.go:3355
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error)
```

Товары витрин **могут публиковаться** в маркетплейс через конвертацию:
```
StorefrontProduct → (конвертация) → MarketplaceListing
```

---

## Количественный анализ полей

### Статистика

```
MarketplaceListing:
├─ Общие поля:      22 (49%)
├─ Уникальные поля: 23 (51%)
└─ ВСЕГО:           45 полей

StorefrontProduct:
├─ Общие поля:      22 (73%)
├─ Уникальные поля:  8 (27%)
└─ ВСЕГО:           30 полей

Объединенная сущность Product:
├─ Общие поля:      22 (41%)
├─ Уникальные ML:   23 (43%)
├─ Уникальные SP:    8 (15%)
└─ ВСЕГО:           53 поля
```

### Общие поля (22 поля)

| № | Поле | MarketplaceListing | StorefrontProduct | Совместимость |
|---|------|-------------------|-------------------|---------------|
| 1 | **ID** | `ID int` | `ID int` | ✅ 100% |
| 2 | **Название** | `Title string` | `Name string` | ⚠️ 95% (разные имена) |
| 3 | **Описание** | `Description string` | `Description string` | ✅ 100% |
| 4 | **Цена** | `Price float64` | `Price float64` | ✅ 100% |
| 5 | **Категория** | `CategoryID int` | `CategoryID int` | ✅ 100% |
| 6 | **Изображения** | `Images []MarketplaceImage` | `Images []StorefrontProductImage` | ✅ 90% (оба ImageInterface) |
| 7 | **Категория (объект)** | `Category *MarketplaceCategory` | `Category *MarketplaceCategory` | ✅ 100% |
| 8 | **Создано** | `CreatedAt time.Time` | `CreatedAt time.Time` | ✅ 100% |
| 9 | **Обновлено** | `UpdatedAt time.Time` | `UpdatedAt time.Time` | ✅ 100% |
| 10 | **Переводы** | `Translations TranslationMap` | `Translations map[string]map[string]string` | ⚠️ 70% (разные типы) |
| 11 | **Адрес переводы** | `AddressMultilingual map[string]string` | `AddressTranslations map[string]string` | ⚠️ 95% (разные имена) |
| 12 | **Широта** | `Latitude *float64` | `IndividualLatitude *float64` | ⚠️ 100% (разные имена) |
| 13 | **Долгота** | `Longitude *float64` | `IndividualLongitude *float64` | ⚠️ 100% (разные имена) |
| 14 | **Показать на карте** | `ShowOnMap bool` | `ShowOnMap bool` | ✅ 100% |
| 15 | **Приватность локации** | `LocationPrivacy string` | `LocationPrivacy *string` | ⚠️ 95% |
| 16 | **Адрес** | `Location string` + `City` + `Country` | `IndividualAddress *string` | ⚠️ 90% |
| 17 | **Атрибуты** | `Attributes []ListingAttributeValue` | `Attributes JSONB` | ❌ 60% (разная структура) |
| 18 | **Активность** | `Status string` (5 состояний) | `IsActive bool` (2 состояния) | ❌ 70% (разная семантика) |
| 19 | **Просмотры** | `ViewsCount int` | `ViewCount int` | ⚠️ 100% (разные имена) |
| 20 | **Варианты** | `Variants []MarketplaceListingVariant` | `Variants []StorefrontProductVariant` | ❌ 60% (разные структуры) |
| 21 | **Остаток** | `StockQuantity *int` | `StockQuantity *int` | ✅ 100% |
| 22 | **Статус остатка** | `StockStatus *string` | `StockStatus string` | ⚠️ 95% |

**Легенда:**
- ✅ Полностью совместимы
- ⚠️ Совместимы с небольшими изменениями (переименование, nullable)
- ❌ Требуют значительной переработки

### Уникальные поля MarketplaceListing (23 поля)

| № | Поле | Тип | Назначение | NULL? |
|---|------|-----|-----------|-------|
| 1 | `UserID` | `int` | Владелец объявления (P2P) | ❌ NOT NULL |
| 2 | `Condition` | `string` | Состояние (новый/б/у) | ❌ NOT NULL |
| 3 | `Status` | `string` | draft/active/sold/archived | ❌ NOT NULL |
| 4 | `HelpfulVotes` | `int` | Голоса "полезно" | ✅ DEFAULT 0 |
| 5 | `NotHelpfulVotes` | `int` | Голоса "не полезно" | ✅ DEFAULT 0 |
| 6 | `IsFavorite` | `bool` | В избранном текущего пользователя | ✅ DEFAULT false |
| 7 | `OldPrice` | `*float64` | Старая цена (скидки) | ✅ NULL |
| 8 | `HasDiscount` | `bool` | Есть скидка | ✅ DEFAULT false |
| 9 | `DiscountPercentage` | `*int` | Процент скидки | ✅ NULL |
| 10 | `Metadata` | `map[string]interface{}` | Дополнительные данные | ✅ NULL |
| 11 | `AverageRating` | `float64` | Средняя оценка | ✅ DEFAULT 0 |
| 12 | `ReviewCount` | `int` | Количество отзывов | ✅ DEFAULT 0 |
| 13 | `StorefrontID` | `*int` | Связь с витриной | ✅ NULL |
| 14 | `Storefront` | `*Storefront` | Данные витрины | ✅ NULL |
| 15 | `ExternalID` | `string` | ID из внешней системы | ✅ NULL |
| 16 | `IsStorefrontProduct` | `bool` | Флаг товара витрины | ✅ DEFAULT false |
| 17 | `OriginalLanguage` | `string` | Оригинальный язык | ✅ NULL |
| 18 | `RawTranslations` | `interface{}` | Сырые данные переводов | ✅ NULL |
| 19 | `CategoryPathNames` | `[]string` | Путь категорий (названия) | ✅ NULL |
| 20 | `CategoryPathIds` | `[]int` | Путь категорий (ID) | ✅ NULL |
| 21 | `CategoryPathSlugs` | `[]string` | Путь категорий (slugs) | ✅ NULL |
| 22 | `CategoryPath` | `[]string` | Путь категорий | ✅ NULL |
| 23 | `User` | `*User` | Данные пользователя | ✅ NULL (join) |

### Уникальные поля StorefrontProduct (8 полей)

| № | Поле | Тип | Назначение | NULL? |
|---|------|-----|-----------|-------|
| 1 | `StorefrontID` | `int` | Витрина-владелец | ❌ NOT NULL |
| 2 | `Currency` | `string` | Валюта товара | ❌ NOT NULL |
| 3 | `SKU` | `*string` | Артикул продавца | ✅ NULL |
| 4 | `Barcode` | `*string` | Штрихкод (EAN/UPC) | ✅ NULL |
| 5 | `IsActive` | `bool` | Активность товара | ❌ NOT NULL |
| 6 | `SoldCount` | `int` | Количество продаж | ✅ DEFAULT 0 |
| 7 | `HasIndividualLocation` | `bool` | Есть своя локация | ✅ DEFAULT false |
| 8 | `HasVariants` | `bool` | Есть варианты | ✅ DEFAULT false |

---

## Качественный анализ

### Критичные различия (блокируют прямое объединение)

#### 1. Владение товаром - ПРИНЦИПИАЛЬНО разное

```go
// MarketplaceListing - P2P модель
UserID int  // Владелец = физлицо/юрлицо

// StorefrontProduct - B2C модель
StorefrontID int  // Владелец = магазин (у которого свой UserID)
```

**Проблема при объединении:**
```go
type Product struct {
    UserID       *int  // NULL для storefront products
    StorefrontID *int  // NULL для marketplace listings

    // Нужна валидация: ОДИН из двух должен быть заполнен!
    // Невозможно сделать NOT NULL constraint в БД
}
```

#### 2. Состояние товара - разная семантика

```go
// MarketplaceListing
Status string  // "draft" | "active" | "sold" | "archived" | "moderation"

// StorefrontProduct
IsActive bool  // true | false
```

**Жизненный цикл:**
- **MarketplaceListing**: draft → active → sold → archived
- **StorefrontProduct**: draft → active ↔ inactive (+ управление остатками)

#### 3. Атрибуты - несовместимые подходы

```go
// MarketplaceListing - структурированные категорийные атрибуты
Attributes []ListingAttributeValue
// [{id: 1, name: "Цвет", value: "Красный", category_id: 10}, ...]

// StorefrontProduct - произвольные JSON атрибуты продавца
Attributes JSONB
// {"material": "cotton", "size": "XL", "brand": "Nike"}
```

**Вывод:** Это **разные паттерны работы с данными**.

---

## Доказательства исторического расхождения

### 1. Функция конвертации (smoking gun)

**Файл:** `backend/internal/proj/marketplace/storage/postgres/marketplace.go:3355`

```go
func (s *Storage) getStorefrontProductAsListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
    // ...
    err := s.pool.QueryRow(ctx, `
        SELECT
            sp.id, sp.storefront_id, sf.user_id, sp.category_id,
            sp.name,              -- ← StorefrontProduct.Name
            sp.description,
            sp.price,
            'new' as condition,   -- ← ХАРДКОД!
            'active' as status,   -- ← ХАРДКОД!
            '' as location,       -- ← ХАРДКОД пустой строки
            0 as latitude,        -- ← ХАРДКОД 0
            0 as longitude,       -- ← ХАРДКОД 0
            '' as city,           -- ← ХАРДКОД
            '' as country,        -- ← ХАРДКОД
            sp.view_count,        -- ← StorefrontProduct.ViewCount
            sp.created_at, sp.updated_at,
            false as show_on_map, -- ← ХАРДКОД
            'sr' as original_language,
            c.name as category_name, c.slug as category_slug,
            '{}'::jsonb as metadata
        FROM storefront_products sp
        LEFT JOIN storefronts sf ON sp.storefront_id = sf.id
        LEFT JOIN marketplace_categories c ON sp.category_id = c.id
        WHERE sp.id = $1 AND sp.is_active = true
    `, id).Scan(
        &listing.ID, &listing.StorefrontID, &listing.UserID, &listing.CategoryID,
        &listing.Title,        // ← sp.name → listing.Title
        &listing.Description, &listing.Price, &listing.Condition, &listing.Status,
        &listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
        &listing.Country,
        &listing.ViewsCount,   // ← sp.view_count → listing.ViewsCount
        &listing.CreatedAt, &listing.UpdatedAt,
        &listing.ShowOnMap, &listing.OriginalLanguage,
        &categoryName, &categorySlug, &listing.Metadata,
    )

    // Для storefront продуктов нет атрибутов пока что
    listing.Attributes = []models.ListingAttributeValue{}
    listing.Translations = make(map[string]map[string]string)
}
```

### Маппинг полей доказывает гипотезу:

| StorefrontProduct | → | MarketplaceListing | Вывод |
|-------------------|---|-------------------|-------|
| `sp.name` | → | `listing.Title` | **Одно и то же!** Просто разные названия |
| `sp.view_count` | → | `listing.ViewsCount` | **Одно и то же!** `Count` vs `s` |
| `sp.price` | → | `listing.Price` | ✅ Идентичны |
| `sp.category_id` | → | `listing.CategoryID` | ✅ Идентичны |

### 2. Хардкод значений - признак костыля

```sql
'new' as condition,        -- Для storefront всегда "новый"
'active' as status,        -- Всегда активный
'' as location,            -- Локация игнорируется
0 as latitude,             -- Координаты игнорируются (хотя у StorefrontProduct они есть!)
false as show_on_map       -- Не показывать на карте
```

**Вывод:** StorefrontProduct изначально **задумывался как подтип** MarketplaceListing, но со временем структуры разошлись.

### 3. Комментарий "пока что"

```go
// Для storefront продуктов нет атрибутов пока что
listing.Attributes = []models.ListingAttributeValue{}
```

Слово **"пока что"** говорит о намерении унифицировать в будущем!

### 4. Процент использования общих полей

```
StorefrontProduct: 73% общих полей → почти вся структура общая
MarketplaceListing: 49% общих полей → половина специфичных
```

**Вывод:** StorefrontProduct проще и ближе к "базовому" товару, MarketplaceListing - расширенная версия.

---

## Преимущества объединения

### ✅ Технические

1. **Унификация кода**
   - Один репозиторий вместо двух
   - Один набор handlers с общей логикой
   - Меньше дублирования CRUD операций

2. **Упрощение поиска**
   - Единый поисковый индекс в OpenSearch
   - Один API endpoint: `GET /api/v1/products?type=listing`
   - Проще делать общую витрину "все товары"

3. **Гибкость бизнес-логики**
   - Легче "превращать" listing в storefront product
   - Пользователь может "апгрейдить" объявление до магазина
   - Проще миграция данных между типами

4. **Упрощение избранного и корзины**
   - Одна таблица `favorites` для всех типов
   - Одна корзина для marketplace + storefront
   - Не нужно джойнить две таблицы

5. **Аналитика и статистика**
   - Проще считать общую статистику по платформе
   - Один запрос для "топ товаров"
   - Легче строить рекомендации

6. **Устранение технического долга**
   - Унификация названий полей (`Title` vs `Name` → `Title`)
   - Единый формат переводов
   - Последовательная архитектура

### ✅ Бизнесовые

1. **Упрощение UX**
   - Единый поиск по всем товарам
   - Единое избранное
   - Единая корзина

2. **Гибкость монетизации**
   - Легко вводить различные тарифы для типов
   - Простая миграция пользователей между моделями

3. **Развитие продукта**
   - Новые фичи применяются ко всем товарам
   - Не нужно дублировать функционал

---

## Недостатки объединения

### ❌ Технические

#### 1. "Раздутая" модель (God Object)

```go
type Product struct {
    // 22 общих поля
    ID, Title, Description, Price...

    // 23 поля ТОЛЬКО для listing (NULL для storefront)
    Condition string
    HelpfulVotes int
    UserID int
    // ...

    // 8 полей ТОЛЬКО для storefront (NULL для listing)
    StorefrontID int
    SKU string
    Currency string
    // ...
}
```

**53 поля** в одной структуре - сложно понять и поддерживать.

#### 2. Усложнение валидации

```go
func (p *Product) Validate() error {
    if p.Type == "listing" {
        if p.Condition == "" {
            return errors.New("condition required")
        }
        if p.UserID == nil {
            return errors.New("user_id required")
        }
    } else if p.Type == "storefront_product" {
        if p.StorefrontID == nil {
            return errors.New("storefront_id required")
        }
        if p.SKU == nil {
            return errors.New("sku required")
        }
    }
    // ... еще 50 строк условий
}
```

#### 3. Проблемы с базой данных

**Вариант A: Одна таблица**
```sql
CREATE TABLE products (
    id INT PRIMARY KEY,
    type VARCHAR(50),

    -- 22 общих поля
    title VARCHAR(255) NOT NULL,

    -- 23 поля listing (NULL в 50% случаев)
    user_id INT,           -- NULL для storefront
    condition VARCHAR(50), -- NULL для storefront
    helpful_votes INT,     -- NULL для storefront

    -- 8 полей storefront (NULL в 50% случаев)
    storefront_id INT,     -- NULL для listing
    sku VARCHAR(100),      -- NULL для listing
    currency VARCHAR(3)    -- NULL для listing
);
```

**Проблемы:**
- 🔴 **58% NULL значений** (31 поле из 53 пустые)
- 🔴 Невозможно сделать `NOT NULL` constraint для специфичных полей
- 🔴 Индексы раздуваются
- 🔴 Сложные CHECK constraints

**Вариант B: Наследование (Class Table Inheritance)**
```sql
CREATE TABLE products (id, type, title, price, ...);  -- общие
CREATE TABLE listing_specific (product_id, condition, user_id, ...);
CREATE TABLE storefront_specific (product_id, storefront_id, sku, ...);
```

**Проблемы:**
- 🔴 JOIN при каждом запросе (медленнее)
- 🔴 Сложнее код репозитория
- 🟡 По сути это "разделение с общей базой"

#### 4. Усложнение бизнес-логики

```go
func (s *OrderService) CreateOrder(product *Product) error {
    if product.Type == "listing" {
        // P2P логика: эскроу, уведомления продавцу
    } else if product.Type == "storefront_product" {
        // B2C логика: складской учет, Post Express
    }
    // Гигантский if-else во ВСЕХ методах
}
```

Каждый метод превращается в `switch` по типу.

#### 5. Потеря type safety

**Сейчас:**
```go
func CreateListing(listing MarketplaceListing) error {
    // Компилятор знает все поля
}
```

**После объединения:**
```go
func CreateProduct(product Product) error {
    // Нужно проверять Type вручную
    if product.Type == "listing" {
        if product.Condition == "" {
            // runtime error вместо compile error
        }
    }
}
```

#### 6. Сложность миграции

```bash
# Затронет ~50 файлов
grep -r "MarketplaceListing" backend/ | wc -l  # ~500 упоминаний
grep -r "StorefrontProduct" backend/ | wc -l   # ~300 упоминаний
```

**Риски:**
- 🔴 Огромная миграция данных
- 🔴 Высокий риск багов
- 🔴 Невозможно откатить частично

### ❌ Бизнесовые

1. **Разные жизненные циклы**
   - Listing: draft → active → **sold** → archived
   - Storefront: draft → active ↔ inactive (продается постоянно)

2. **Разные права доступа**
   - Listing: владелец (UserID)
   - Storefront: владелец витрины + персонал с permissions

3. **Разные требования к данным**
   - Listing: condition обязательно
   - Storefront: SKU/Barcode критично

---

## Рекомендации

### 🎯 Вариант 1: Hybrid Architecture (РЕКОМЕНДУЕТСЯ)

**Суть:** Объединить только для поиска/отображения, разделить для бизнес-логики.

```go
// 1. Интерфейс для общих операций
type ProductInterface interface {
    GetID() int
    GetTitle() string
    GetDescription() string
    GetPrice() float64
    GetCategoryID() int
    GetImages() []ImageInterface
    GetLocation() (lat, lng float64)
    GetCreatedAt() time.Time
    // ... 14 общих методов
}

// 2. Обе сущности реализуют интерфейс
func (m *MarketplaceListing) GetTitle() string { return m.Title }
func (s *StorefrontProduct) GetTitle() string  { return s.Name }

// 3. Search View для OpenSearch (22 общих поля)
type ProductSearchView struct {
    ID          int       `json:"id"`
    Type        string    `json:"type"`
    Title       string    `json:"title"`
    Price       float64   `json:"price"`
    CategoryID  int       `json:"category_id"`
    Images      []string  `json:"images"`
    CreatedAt   time.Time `json:"created_at"`
    // ... остальные общие поля
}

// 4. БД - раздельные таблицы
// marketplace_listings (45 полей, корректные constraints)
// storefront_products (30 полей, корректные constraints)

// 5. OpenSearch - единый индекс
// products (22 общих поля + type)
```

#### Преимущества:
- ✅ Единый поиск
- ✅ Type safety
- ✅ NO NULL поля в БД
- ✅ Простая валидация
- ✅ Гибкость
- ✅ Performance

#### Недостатки:
- ⚠️ Два репозитория (но это норм для bounded contexts)
- ⚠️ Конвертация в SearchView (но она уже есть)

---

### 🎯 Вариант 2: Полное объединение с Type

**Суть:** Одна таблица `products` с полем `type`.

```go
type Product struct {
    ID   int    `json:"id" db:"id"`
    Type string `json:"type" db:"type"` // "listing" | "storefront_product"

    // Владение (один обязателен)
    UserID       *int `json:"user_id,omitempty" db:"user_id"`
    StorefrontID *int `json:"storefront_id,omitempty" db:"storefront_id"`

    // Общие поля (унифицированные)
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Price       float64   `json:"price" db:"price"`
    Currency    string    `json:"currency" db:"currency"`
    ViewCount   int       `json:"view_count" db:"view_count"`
    // ... остальные 17 общих

    // Специфичные listing
    Condition       *string  `json:"condition,omitempty" db:"condition"`
    HelpfulVotes    *int     `json:"helpful_votes,omitempty" db:"helpful_votes"`
    OldPrice        *float64 `json:"old_price,omitempty" db:"old_price"`
    // ... остальные 20

    // Специфичные storefront
    SKU       *string `json:"sku,omitempty" db:"sku"`
    Barcode   *string `json:"barcode,omitempty" db:"barcode"`
    SoldCount *int    `json:"sold_count,omitempty" db:"sold_count"`
    // ... остальные 5

    // Статус (унифицированный)
    Status   string `json:"status" db:"status"` // "active" | "inactive" | "sold" | "draft"
    IsActive bool   `json:"is_active" db:"is_active"` // computed from Status
}
```

#### Преимущества:
- ✅ Единая структура
- ✅ Один репозиторий
- ✅ Один набор handlers
- ✅ Проще добавлять общие поля

#### Недостатки:
- ❌ 53 поля в структуре
- ❌ 58% NULL в БД
- ❌ Сложная валидация
- ❌ Потеря type safety
- ❌ Гигантская миграция

---

### 🎯 Вариант 3: Не объединять (status quo)

**Суть:** Оставить как есть, но зафиксировать маппинг.

#### Преимущества:
- ✅ Минимум изменений
- ✅ Нет рисков
- ✅ Type safety
- ✅ Простая валидация

#### Недостатки:
- ❌ Технический долг остается
- ❌ Дублирование логики
- ❌ Два репозитория

---

## План миграции

### Если выбран Вариант 1 (Hybrid Architecture)

#### Фаза 1: Подготовка (1-2 дня)

1. ✅ Создать интерфейс `ProductInterface`
2. ✅ Создать `ProductSearchView` структуру
3. ✅ Реализовать интерфейс в обеих сущностях
4. ✅ Написать тесты для интерфейса

#### Фаза 2: OpenSearch (2-3 дня)

1. ✅ Создать единый индекс `products`
2. ✅ Написать функции конвертации в SearchView
3. ✅ Мигрировать данные в новый индекс
4. ✅ Обновить поисковые запросы

#### Фаза 3: Общие сервисы (3-4 дня)

1. ✅ Переписать favorites на интерфейс
2. ✅ Переписать cart на интерфейс
3. ✅ Обновить API endpoints
4. ✅ Тесты

#### Фаза 4: Cleanup (1 день)

1. ✅ Удалить дублирующийся код
2. ✅ Обновить документацию
3. ✅ Code review

**Общее время:** 7-10 дней

---

### Если выбран Вариант 2 (Полное объединение)

#### Фаза 1: Проектирование (2-3 дня)

1. ✅ Финализировать схему Product
2. ✅ Написать валидацию
3. ✅ Спроектировать миграции
4. ✅ План отката

#### Фаза 2: Database Migration (3-5 дней)

1. ✅ Создать таблицу `products`
2. ✅ Мигрировать `marketplace_listings` (type='listing')
3. ✅ Мигрировать `storefront_products` (type='storefront_product')
4. ✅ Обновить foreign keys
5. ✅ Тесты миграции

#### Фаза 3: Code Migration (7-10 дней)

1. ✅ Обновить models
2. ✅ Обновить repositories (~10 файлов)
3. ✅ Обновить services (~15 файлов)
4. ✅ Обновить handlers (~20 файлов)
5. ✅ Обновить тесты (~30 файлов)

#### Фаза 4: OpenSearch (2-3 дня)

1. ✅ Обновить индексацию
2. ✅ Переиндексировать данные
3. ✅ Обновить поисковые запросы

#### Фаза 5: Testing & Rollout (5-7 дней)

1. ✅ E2E тесты
2. ✅ Performance тесты
3. ✅ Staging deployment
4. ✅ Monitoring
5. ✅ Production rollout

**Общее время:** 19-28 дней (3-4 недели)

**Риски:**
- 🔴 Высокий риск регрессии
- 🔴 Сложно откатить
- 🟡 Простой продакшена при миграции БД

---

## Итоговая рекомендация

### 🏆 Рекомендуется: Вариант 1 (Hybrid Architecture)

**Почему:**
1. ✅ Минимальные изменения (7-10 дней vs 3-4 недели)
2. ✅ Низкий риск регрессии
3. ✅ Получаем главную выгоду (единый поиск)
4. ✅ Сохраняем type safety
5. ✅ NO NULL поля в БД
6. ✅ Можно откатить на любом этапе

**Что делать дальше:**
1. Обсудить решение с командой
2. Начать с Фазы 1 (интерфейс)
3. Постепенно мигрировать функционал
4. Измерить эффект

---

## Приложения

### A. Маппинг полей для миграции

| Общее название | MarketplaceListing | StorefrontProduct |
|---------------|-------------------|-------------------|
| `title` | `Title` | `Name` |
| `view_count` | `ViewsCount` | `ViewCount` |
| `latitude` | `Latitude` | `IndividualLatitude` |
| `longitude` | `Longitude` | `IndividualLongitude` |
| `address` | `Location` | `IndividualAddress` |
| `address_translations` | `AddressMultilingual` | `AddressTranslations` |

### B. Примеры кода

#### Интерфейс ProductInterface

```go
// backend/internal/domain/interfaces/product.go
package interfaces

import "time"

type ProductInterface interface {
    // Core fields
    GetID() int
    GetTitle() string
    GetDescription() string
    GetPrice() float64
    GetCategoryID() int

    // Images
    GetImages() []ImageInterface

    // Location
    GetLatitude() *float64
    GetLongitude() *float64
    GetShowOnMap() bool

    // Timestamps
    GetCreatedAt() time.Time
    GetUpdatedAt() time.Time

    // Type
    GetType() string // "listing" | "storefront_product"
}
```

#### Реализация в MarketplaceListing

```go
func (m *MarketplaceListing) GetTitle() string {
    return m.Title
}

func (m *MarketplaceListing) GetType() string {
    return "listing"
}
```

#### Реализация в StorefrontProduct

```go
func (s *StorefrontProduct) GetTitle() string {
    return s.Name
}

func (s *StorefrontProduct) GetType() string {
    return "storefront_product"
}
```

#### ProductSearchView

```go
// backend/internal/domain/models/product_search_view.go
package models

import "time"

type ProductSearchView struct {
    ID          int       `json:"id"`
    Type        string    `json:"type"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    CategoryID  int       `json:"category_id"`
    Images      []string  `json:"images"`
    Latitude    *float64  `json:"latitude,omitempty"`
    Longitude   *float64  `json:"longitude,omitempty"`
    ShowOnMap   bool      `json:"show_on_map"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Converter functions
func ToSearchView(p ProductInterface) *ProductSearchView {
    images := make([]string, 0)
    for _, img := range p.GetImages() {
        images = append(images, img.GetImageURL())
    }

    return &ProductSearchView{
        ID:          p.GetID(),
        Type:        p.GetType(),
        Title:       p.GetTitle(),
        Description: p.GetDescription(),
        Price:       p.GetPrice(),
        CategoryID:  p.GetCategoryID(),
        Images:      images,
        Latitude:    p.GetLatitude(),
        Longitude:   p.GetLongitude(),
        ShowOnMap:   p.GetShowOnMap(),
        CreatedAt:   p.GetCreatedAt(),
        UpdatedAt:   p.GetUpdatedAt(),
    }
}
```

---

**Документ создан:** 2025-10-09
**Следующий шаг:** Обсуждение с командой и выбор варианта
