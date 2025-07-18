# Руководство по хранению геоданных и адресов

## Обзор системы

В системе используется унифицированный подход к хранению геоданных для всех типов объектов (объявления marketplace, товары витрин, сами витрины).

## Структура таблиц

### 1. marketplace_listings
Основная таблица объявлений со следующими геополями:
```sql
- location VARCHAR(255)       -- Текстовый адрес
- latitude NUMERIC(10,8)      -- Широта
- longitude NUMERIC(11,8)     -- Долгота  
- address_city VARCHAR(100)   -- Город
- address_country VARCHAR(100) -- Страна
```

### 2. listings_geo (устаревшая)
Промежуточная таблица для геоданных (миграция 099):
```sql
- listing_id INTEGER          -- ID объявления
- location GEOGRAPHY(Point)   -- PostGIS точка
- geohash VARCHAR(12)         -- Геохеш для быстрого поиска
- formatted_address TEXT      -- Отформатированный адрес
- address_components JSONB    -- Компоненты адреса
```

### 3. unified_geo (основная)
Унифицированная таблица для всех геоданных (миграция 113):
```sql
- source_type geo_source_type -- Тип объекта ('marketplace_listing', 'storefront_product', 'storefront')
- source_id BIGINT           -- ID объекта
- location GEOGRAPHY(Point)   -- PostGIS точка
- geohash VARCHAR(12)        -- Геохеш
- formatted_address TEXT     -- Отформатированный адрес
```

## Процесс сохранения геоданных

### При создании объявления:

1. **Frontend отправляет**:
   ```json
   {
     "location": "Vase Stajica, Novi Sad 21101, Serbia",
     "latitude": 45.250903,
     "longitude": 19.843395,
     "city": "Novi Sad",     // Заполняется автоматически из геокодинга
     "country": "Serbia"      // Заполняется автоматически из геокодинга
   }
   ```

2. **Backend сохраняет в marketplace_listings**:
   - `location` - полный адрес
   - `latitude`, `longitude` - координаты
   - `address_city` - город (НЕ `city` - это ошибка в коде!)
   - `address_country` - страна (НЕ `country` - это ошибка в коде!)

3. **Автоматически создается запись в listings_geo**:
   - Через триггер или в коде создается запись с PostGIS точкой
   - Вычисляется geohash для оптимизации поиска

4. **Должна создаваться запись в unified_geo**:
   - `source_type = 'marketplace_listing'`
   - `source_id` = ID созданного объявления
   - `location` = PostGIS точка из координат
   - `geohash` = вычисленный геохеш

## GIS поиск

GIS API использует таблицу `unified_geo` для поиска:
```sql
SELECT * FROM unified_geo 
WHERE source_type = 'marketplace_listing'
  AND ST_DWithin(location, ST_MakePoint(lng, lat)::geography, radius)
```

## Известные проблемы

1. **Несоответствие имен полей**: Frontend отправляет `city`/`country`, но БД ожидает `address_city`/`address_country`

2. **Отсутствие автоматического переноса в unified_geo**: При создании объявления данные не попадают автоматически в unified_geo

3. **Дублирование данных**: Геоданные хранятся в трех местах (marketplace_listings, listings_geo, unified_geo)

## Рекомендации для исправления

1. **В коде создания объявления** (`marketplace.go:InsertListing`):
   - Исправить маппинг полей city → address_city, country → address_country
   - Добавить создание записи в unified_geo после успешного создания объявления

2. **Для импорта товаров в витрины**:
   - Использовать unified_geo с source_type = 'storefront_product'
   - Сохранять геоданные аналогично объявлениям

3. **Миграция существующих данных**:
   ```sql
   -- Перенос из listings_geo в unified_geo
   INSERT INTO unified_geo (source_type, source_id, location, geohash, formatted_address)
   SELECT 'marketplace_listing', listing_id, location, geohash, formatted_address
   FROM listings_geo
   ON CONFLICT (source_type, source_id) DO NOTHING;
   ```

## Проверка целостности данных

```sql
-- Проверить объявления без геоданных
SELECT ml.id, ml.title 
FROM marketplace_listings ml
LEFT JOIN unified_geo ug ON ug.source_type = 'marketplace_listing' AND ug.source_id = ml.id
WHERE ug.id IS NULL AND ml.status = 'active';

-- Проверить соответствие координат
SELECT ml.id, 
       ml.latitude, ml.longitude,
       ST_Y(ug.location::geometry) as geo_lat,
       ST_X(ug.location::geometry) as geo_lng
FROM marketplace_listings ml
JOIN unified_geo ug ON ug.source_type = 'marketplace_listing' AND ug.source_id = ml.id
WHERE ml.latitude IS NOT NULL;
```