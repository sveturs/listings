# План реализации GIS модуля для маркетплейса Sve Tu

## Этап 1: Базовая инфраструктура (2 недели)

### Неделя 1: Backend инфраструктура

#### День 1-2: Настройка PostGIS и миграции

**Задачи:**
1. Установка и настройка PostGIS расширения
2. Создание миграции для включения PostGIS
3. Создание таблиц для геоданных
4. Настройка пространственных индексов

**Файлы для создания:**
```
backend/migrations/
├── 000050_enable_postgis.up.sql
├── 000050_enable_postgis.down.sql
├── 000051_create_gis_tables.up.sql
├── 000051_create_gis_tables.down.sql
├── 000052_add_spatial_indexes.up.sql
└── 000052_add_spatial_indexes.down.sql
```

**SQL миграции:**
```sql
-- 000050_enable_postgis.up.sql
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- 000051_create_gis_tables.up.sql
-- Таблица административных границ Сербии
CREATE TABLE administrative_boundaries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_cyrillic VARCHAR(100),
    name_latin VARCHAR(100),
    type VARCHAR(50) NOT NULL, -- 'municipality', 'district', 'region', 'country'
    parent_id INT REFERENCES administrative_boundaries(id),
    boundary GEOMETRY(MULTIPOLYGON, 4326),
    center_point GEOMETRY(POINT, 4326),
    population INT,
    area_km2 DECIMAL(10,2),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для кэширования геокодирования
CREATE TABLE geocoding_cache (
    id SERIAL PRIMARY KEY,
    address_hash VARCHAR(64) UNIQUE NOT NULL,
    original_address TEXT NOT NULL,
    normalized_address TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    point GEOMETRY(POINT, 4326),
    formatted_address_sr TEXT,
    formatted_address_en TEXT,
    confidence_score FLOAT,
    provider VARCHAR(50), -- 'nominatim', 'google', 'manual'
    raw_response JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица популярных мест
CREATE TABLE popular_locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_cyrillic VARCHAR(100),
    type VARCHAR(50), -- 'shopping_center', 'market', 'transport_hub'
    point GEOMETRY(POINT, 4326) NOT NULL,
    radius_meters INT DEFAULT 100,
    city VARCHAR(100),
    address TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 000052_add_spatial_indexes.up.sql
-- Преобразование существующих координат в геометрию
ALTER TABLE marketplace_listings 
    ADD COLUMN IF NOT EXISTS location_point GEOMETRY(POINT, 4326);

UPDATE marketplace_listings 
    SET location_point = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

ALTER TABLE user_storefronts 
    ADD COLUMN IF NOT EXISTS location_point GEOMETRY(POINT, 4326);

UPDATE user_storefronts 
    SET location_point = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Создание пространственных индексов
CREATE INDEX idx_listings_location_point ON marketplace_listings USING GIST(location_point);
CREATE INDEX idx_storefronts_location_point ON user_storefronts USING GIST(location_point);
CREATE INDEX idx_boundaries_boundary ON administrative_boundaries USING GIST(boundary);
CREATE INDEX idx_boundaries_center ON administrative_boundaries USING GIST(center_point);
CREATE INDEX idx_geocoding_point ON geocoding_cache USING GIST(point);
CREATE INDEX idx_popular_locations_point ON popular_locations USING GIST(point);

-- Индексы для текстового поиска
CREATE INDEX idx_boundaries_name ON administrative_boundaries(name);
CREATE INDEX idx_boundaries_name_cyrillic ON administrative_boundaries(name_cyrillic);
CREATE INDEX idx_geocoding_hash ON geocoding_cache(address_hash);
```

#### День 3-4: Базовая структура GIS модуля

**Структура модуля:**
```
backend/internal/proj/gis/
├── module.go                    # Инициализация модуля
├── handler/
│   ├── handler.go              # Базовый обработчик
│   ├── spatial_search.go       # Пространственный поиск
│   └── geocoding.go            # Геокодирование
├── service/
│   ├── interface.go            # Интерфейсы сервисов
│   ├── spatial_service.go      # Пространственные операции
│   └── geocoding_service.go    # Сервис геокодирования
├── repository/
│   ├── interface.go            # Интерфейсы репозиториев
│   └── postgres/
│       ├── spatial_repo.go     # Пространственные запросы
│       └── geocoding_repo.go   # Кэш геокодирования
├── types/
│   ├── coordinates.go          # Типы координат
│   ├── geometry.go             # Геометрические типы
│   └── requests.go             # Типы запросов/ответов
└── utils/
    ├── geometry.go             # Геометрические утилиты
    └── validation.go           # Валидация координат
```

**Основные типы данных:**
```go
// types/coordinates.go
package types

type Coordinates struct {
    Latitude  float64 `json:"latitude" validate:"required,latitude"`
    Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type BoundingBox struct {
    MinLat float64 `json:"min_lat"`
    MinLng float64 `json:"min_lng"`
    MaxLat float64 `json:"max_lat"`
    MaxLng float64 `json:"max_lng"`
}

type Point struct {
    Type        string    `json:"type"`
    Coordinates []float64 `json:"coordinates"` // [lng, lat]
}

// types/requests.go
type SpatialSearchRequest struct {
    Center   *Coordinates `json:"center"`
    Radius   int          `json:"radius_km" validate:"min=1,max=100"`
    BBox     *BoundingBox `json:"bbox"`
    Category *int         `json:"category_id"`
    Limit    int          `json:"limit" validate:"min=1,max=100"`
    Offset   int          `json:"offset" validate:"min=0"`
}

type GeocodeRequest struct {
    Address  string `json:"address" validate:"required,min=3"`
    Language string `json:"language" validate:"omitempty,oneof=sr en"`
}
```

#### День 5: Интеграция с основным приложением

**Изменения в существующих файлах:**

1. **backend/internal/server/server.go** - регистрация GIS модуля
2. **backend/internal/proj/global/service/service.go** - добавление GIS сервиса
3. **backend/internal/config/config.go** - конфигурация для GIS

**Конфигурация:**
```go
// Добавить в config.go
type GISConfig struct {
    EnablePostGIS          bool   `env:"GIS_ENABLE_POSTGIS" envDefault:"true"`
    NominatimURL          string `env:"GIS_NOMINATIM_URL" envDefault:"https://nominatim.openstreetmap.org"`
    DefaultSearchRadius   int    `env:"GIS_DEFAULT_RADIUS_KM" envDefault:"10"`
    MaxSearchRadius       int    `env:"GIS_MAX_RADIUS_KM" envDefault:"100"`
    GeocacheTTL          int    `env:"GIS_GEOCACHE_TTL_DAYS" envDefault:"30"`
    EnableClustering     bool   `env:"GIS_ENABLE_CLUSTERING" envDefault:"true"`
    ClusterMinZoom       int    `env:"GIS_CLUSTER_MIN_ZOOM" envDefault:"10"`
}
```

### Неделя 2: Frontend компоненты и API

#### День 6-7: API endpoints

**REST API endpoints:**
```
GET  /api/v1/gis/search/spatial    # Пространственный поиск
POST /api/v1/gis/geocode           # Геокодирование адреса
GET  /api/v1/gis/geocode/reverse   # Обратное геокодирование
GET  /api/v1/gis/boundaries        # Административные границы
GET  /api/v1/gis/popular-locations # Популярные места
```

**Swagger документация:**
```go
// @Summary Spatial search for listings
// @Description Search listings within radius or bounding box
// @Tags GIS
// @Accept json
// @Produce json
// @Param request body types.SpatialSearchRequest true "Search parameters"
// @Success 200 {object} types.SpatialSearchResponse
// @Router /api/v1/gis/search/spatial [get]
```

#### День 8-9: Базовые React компоненты

**Структура frontend компонентов:**
```
frontend/svetu/src/
├── components/
│   └── GIS/
│       ├── Map/
│       │   ├── BaseMap.tsx           # Базовый компонент карты
│       │   ├── MapContainer.tsx      # Контейнер с состоянием
│       │   └── MapControls.tsx       # Элементы управления
│       ├── Markers/
│       │   ├── ListingMarker.tsx     # Маркер объявления
│       │   ├── ClusterMarker.tsx     # Кластер маркеров
│       │   └── UserLocationMarker.tsx # Местоположение пользователя
│       └── Search/
│           ├── RadiusSelector.tsx    # Выбор радиуса поиска
│           └── LocationSearch.tsx    # Поиск по местоположению
├── hooks/
│   └── gis/
│       ├── useGeolocation.ts        # Геолокация пользователя
│       ├── useMapBounds.ts          # Границы видимой области
│       └── useSpatialSearch.ts     # Пространственный поиск
└── lib/
    └── gis/
        ├── constants.ts             # Константы (центр Сербии и т.д.)
        ├── utils.ts                 # Утилиты для работы с координатами
        └── types.ts                 # TypeScript типы
```

**Базовая конфигурация карты:**
```typescript
// lib/gis/constants.ts
export const SERBIA_CENTER = {
  lat: 44.0165,
  lng: 21.0059
};

export const DEFAULT_ZOOM = 7;
export const CITY_ZOOM = 12;
export const STREET_ZOOM = 16;

export const MAJOR_CITIES = {
  BELGRADE: { lat: 44.7866, lng: 20.4489, name: 'Београд' },
  NOVI_SAD: { lat: 45.2671, lng: 19.8335, name: 'Нови Сад' },
  NIS: { lat: 43.3209, lng: 21.8954, name: 'Ниш' },
  KRAGUJEVAC: { lat: 44.0142, lng: 20.9394, name: 'Крагујевац' }
};
```

#### День 10: Интеграция и тестирование

**Тестовые сценарии:**

1. **Unit тесты backend:**
   - Валидация координат
   - Расчет расстояний
   - Преобразование геометрии
   - Кэширование геокодирования

2. **Integration тесты:**
   - API endpoints
   - PostGIS запросы
   - Производительность пространственных индексов

3. **Frontend тесты:**
   - Рендеринг карты
   - Взаимодействие с маркерами
   - Обработка геолокации

**Метрики производительности:**
- Пространственный поиск < 200ms для 10k записей
- Геокодирование < 100ms (с кэшем)
- Загрузка карты < 2 секунды
- Рендеринг 100 маркеров < 50ms

### Оценка рисков и план митигации

**Технические риски:**
1. **PostGIS не установлен на сервере**
   - Митигация: Fallback на обычные SQL запросы с формулой Haversine
   - Подготовить Docker образ с PostGIS

2. **Производительность пространственных запросов**
   - Митигация: Агрессивное кэширование
   - Материализованные представления для популярных запросов

3. **Лимиты Nominatim API**
   - Митигация: Собственный экземпляр Nominatim
   - Альтернативные провайдеры (Mapbox, Here)

**UX риски:**
1. **Медленная загрузка карты на мобильных**
   - Митигация: Lazy loading
   - Упрощенная мобильная версия

2. **Неточное геокодирование сербских адресов**
   - Митигация: Ручная корректировка координат
   - База популярных адресов

### Зависимости

**NPM пакеты:**
```json
{
  "leaflet": "^1.9.4",
  "react-leaflet": "^4.2.1",
  "@types/leaflet": "^1.9.8",
  "leaflet.markercluster": "^1.5.3",
  "@turf/turf": "^6.5.0"
}
```

**Go модули:**
```go
require (
    github.com/paulmach/orb v0.11.1  // Геометрические операции
    github.com/twpayne/go-geom v1.5.3 // PostGIS типы
)
```

### Контрольные точки

**Конец недели 1:**
- ✓ PostGIS настроен и работает
- ✓ Миграции применены
- ✓ Базовая структура GIS модуля создана
- ✓ API endpoints определены

**Конец недели 2:**
- ✓ API endpoints реализованы
- ✓ Базовые React компоненты созданы
- ✓ Карта отображается с тестовыми данными
- ✓ Все тесты проходят

### Определение успеха этапа 1

Этап считается успешно завершенным, если:
1. PostGIS расширение активировано в базе данных
2. Созданы все необходимые таблицы и индексы
3. Backend GIS модуль интегрирован в приложение
4. Реализованы базовые API endpoints (поиск, геокодирование)
5. Создан простейший компонент карты в React
6. Все тесты проходят успешно
7. Документация API обновлена

### Подготовка к следующему этапу

После завершения этапа 1 будет готова база для:
- Реализации полнофункциональной интерактивной карты
- Добавления кластеризации маркеров
- Интеграции с существующим поиском
- Реализации фильтров по местоположению