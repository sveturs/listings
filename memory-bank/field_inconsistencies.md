# Анализ несоответствий полей между слоями проекта

## Сводная таблица соответствия полей

### 1. Основные идентификаторы

| Поле | База данных | Backend (Go) | Frontend (TS) | OpenSearch | Проблемы |
|------|-------------|--------------|---------------|------------|----------|
| **ID объявления** | `id` (integer) | `ID` / `id` | `id?` (number) | `id` (integer) | ✅ Согласовано |
| **ID пользователя** | `user_id` (integer) | `UserID` / `user_id` | `user_id?` (number) | `user_id` (integer) | ✅ Согласовано |
| **ID категории** | `category_id` (integer) | `CategoryID` / `category_id` | `category_id?` (number) | `category_id` (integer) | ✅ Согласовано |
| **ID витрины** | `storefront_id` (integer) | `StorefrontID` / `storefront_id` | `storefront_id?` (number) | `storefront_id` (integer) | ✅ Согласовано |

### 2. Поля с несоответствиями

#### 2.1. Поля адреса/локации

| Поле | База данных | Backend (Go) | Frontend (TS) | OpenSearch | Проблемы |
|------|-------------|--------------|---------------|------------|----------|
| **Город** | `address_city` | `City` / `city` | `city?` | `city` | ⚠️ Разные названия в БД |
| **Страна** | `address_country` | `Country` / `country` | `country?` | `country` | ⚠️ Разные названия в БД |

**Проблема**: В базе данных используются префиксы `address_`, а в остальных слоях - простые названия.

**Рекомендация**: 
- Либо переименовать поля в БД на `city` и `country`
- Либо добавить маппинг в SQL-запросах: `address_city AS city, address_country AS country`

#### 2.2. Числовые типы

| Поле | База данных | Backend (Go) | Frontend (TS) | OpenSearch |
|------|-------------|--------------|---------------|------------|
| **price** | `numeric` | `float64` | `number` | `float` |
| **latitude** | `numeric` | `*float64` | `number` | `float` |
| **longitude** | `numeric` | `*float64` | `number` | `float` |

**Проблема**: PostgreSQL `numeric` может хранить больше точности, чем `float64`.

**Рекомендация**: Для цен лучше использовать `decimal` тип или хранить в копейках как `bigint`.

### 3. Отсутствующие поля в разных слоях

#### 3.1. Поля, присутствующие только в БД

- `needs_reindex` (boolean) - флаг для переиндексации в OpenSearch
- `external_id` (varchar) - присутствует в Backend/Frontend, но может быть NULL

#### 3.2. Поля, присутствующие только в Backend/Frontend

- `HelpfulVotes` / `helpful_votes` - нет в БД, вероятно вычисляемое поле
- `NotHelpfulVotes` / `not_helpful_votes` - нет в БД, вероятно вычисляемое поле
- `IsFavorite` / `is_favorite` - нет в БД, вероятно из таблицы `marketplace_favorites`
- `AverageRating` / `average_rating` - нет в БД, вычисляемое из таблицы `reviews`
- `ReviewCount` / `review_count` - нет в БД, вычисляемое из таблицы `reviews`
- `OldPrice` / `old_price` - нет в БД
- `HasDiscount` / `has_discount` - нет в БД

**Рекомендация**: Документировать, откуда берутся эти поля (JOIN с другими таблицами или вычисления).

### 4. Типы данных и обязательность

#### 4.1. Обязательные поля в БД vs опциональные в коде

| Поле | БД | Backend | Frontend | Проблема |
|------|----|---------|-----------| --------|
| `title` | NOT NULL | string | string? | Frontend делает все поля опциональными |
| `show_on_map` | NOT NULL (default: true) | bool | boolean? | Frontend делает все поля опциональными |

**Рекомендация**: В TypeScript-типах для создания/обновления объектов сделать обязательные поля required.

### 5. Специфичные проблемы

#### 5.1. Metadata

- **БД**: `jsonb`
- **Backend**: `map[string]interface{}`
- **Frontend**: `{ [key: string]: unknown }`
- **OpenSearch**: отсутствует в маппинге

**Проблема**: Metadata не индексируется в OpenSearch, поиск по metadata невозможен.

#### 5.2. Translations

- **БД**: отсутствует (вероятно, хранится в `metadata` или отдельной таблице)
- **Backend**: `TranslationMap` (map[string]map[string]string)
- **Frontend**: `components['schemas']['models.TranslationMap']`
- **OpenSearch**: есть поля для мультиязычного поиска

### 6. Рекомендации по исправлению

1. **Унифицировать названия полей адреса**:
   ```sql
   ALTER TABLE marketplace_listings 
   RENAME COLUMN address_city TO city;
   
   ALTER TABLE marketplace_listings 
   RENAME COLUMN address_country TO country;
   ```

2. **Добавить отсутствующие поля в БД или документировать их источники**:
   - Создать VIEW с вычисляемыми полями
   - Или документировать JOIN-запросы

3. **Создать строгие TypeScript-типы для операций**:
   ```typescript
   // Для создания
   interface CreateListingInput {
     title: string; // required
     user_id: number; // required
     // ...
   }
   
   // Для ответов API
   interface ListingResponse {
     id: number;
     title: string;
     // все поля обязательные, кроме явно nullable
   }
   ```

4. **Добавить metadata в OpenSearch маппинг**, если требуется поиск по metadata.

5. **Документировать маппинг между слоями** в отдельном файле для разработчиков.

## Критические несоответствия

1. ⚠️ **address_city/address_country vs city/country** - требует изменения БД или SQL-запросов
2. ⚠️ **Вычисляемые поля** - нужна документация источников данных
3. ⚠️ **Обязательность полей** - Frontend типы слишком расслабленные

## Заключение

Основные идентификаторы (id, user_id, category_id, storefront_id) согласованы между всеми слоями. Главные проблемы связаны с:
- Разными названиями полей адреса в БД
- Отсутствием некоторых полей в БД (вычисляемые поля)
- Слишком опциональными типами во Frontend

Рекомендуется начать с исправления названий полей адреса и документирования источников вычисляемых полей.