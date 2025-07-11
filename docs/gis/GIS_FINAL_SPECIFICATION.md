# 🗺️ GIS модуль для Sve Tu Marketplace - Финальная спецификация

**Версия**: 1.0 FINAL  
**Дата**: 2025-01-10  
**Статус**: Готово к реализации

## 📋 Оглавление

1. [Резюме проекта](#резюме-проекта)
2. [Бизнес-требования](#бизнес-требования)
3. [Функциональные требования](#функциональные-требования)
4. [Техническая архитектура](#техническая-архитектура)
5. [База данных](#база-данных)
6. [API спецификация](#api-спецификация)
7. [Frontend компоненты](#frontend-компоненты)
8. [Интеграция с доставкой](#интеграция-с-доставкой)
9. [План реализации](#план-реализации)
10. [Метрики успеха](#метрики-успеха)

---

## 1. Резюме проекта

### 🎯 Цель
Создать современную геолокационную платформу для маркетплейса Sve Tu, которая улучшит пользовательский опыт через интерактивные карты, умный геопоиск и интеграцию со службами доставки.

### 📊 Ключевые показатели
- **Срок реализации**: 12 недель
- **Команда**: 3-4 разработчика
- **Бюджет технологий**: ~$1000/месяц
- **ROI**: 6 месяцев

### ✅ Основные возможности
- Интерактивная карта с объявлениями
- Поиск товаров по радиусу и районам
- Зоны доставки для магазинов
- Интеграция с D Express и Поштой Србије
- Геоаналитика для продавцов
- Мобильная оптимизация

---

## 2. Бизнес-требования

### 2.1 Проблемы, которые решаем

1. **Для покупателей**:
   - Сложно найти товары рядом
   - Неясно, доставят ли в мой район
   - Нет информации о пунктах выдачи

2. **Для продавцов**:
   - Не понимают географию своих клиентов
   - Сложно настроить зоны доставки
   - Нет аналитики по районам

3. **Для бизнеса**:
   - Низкая конверсия локального поиска
   - Высокие затраты на логистику
   - Отсутствие геоданных для аналитики

### 2.2 Ожидаемые результаты

| Метрика | Текущее | Целевое | Срок |
|---------|---------|---------|------|
| Конверсия поиска | 3.2% | 4.5% | 3 мес |
| Использование карты | 0% | 40% | 2 мес |
| Локальные продажи | базовый | +30% | 3 мес |
| Время на сайте | 4:30 | 6:00 | 2 мес |

---

## 3. Функциональные требования

### 3.1 Интерактивная карта

#### Основные функции:
- **Отображение объявлений** на карте с маркерами
- **Кластеризация** при большом количестве маркеров
- **Информационные окна** с превью товара
- **Фильтры** по категориям, цене, доставке
- **Режимы карты**: обычная, спутник, упрощенная

#### Пользовательские сценарии:
```
1. Открыть карту → Увидеть товары рядом
2. Применить фильтры → Обновление в реальном времени
3. Кликнуть на маркер → Превью товара
4. Перейти к товару → Страница объявления
```

### 3.2 Геопоиск

#### Типы поиска:
1. **По радиусу** - "в радиусе 5 км от меня"
2. **По району** - "в Новом Белграде"
3. **По маршруту** - "по дороге домой"
4. **Умный поиск** - "детские товары до 5000 динар рядом"

#### Фильтры поиска:
- Категория товара
- Ценовой диапазон
- Наличие доставки
- Рейтинг продавца
- Время работы (для магазинов)

### 3.3 Зоны доставки

#### Для продавцов:
- Рисование зон на карте
- Настройка стоимости по зонам
- Минимальная сумма заказа
- Время доставки

#### Для покупателей:
- Проверка доставки по адресу
- Расчет стоимости
- Выбор пункта самовывоза

### 3.4 Службы доставки

#### D Express:
- Экспресс-доставка (1-2 дня)
- Отслеживание в реальном времени
- API интеграция

#### Пошта Србије:
- Стандартная доставка (2-5 дней)
- Сеть почтовых отделений
- Бюджетные тарифы

### 3.5 Геоаналитика

#### Для продавцов:
- Тепловая карта покупателей
- Статистика по районам
- Анализ конкурентов
- Рекомендации по расширению

#### Для администрации:
- Общая карта активности
- Проблемные зоны доставки
- Тренды по регионам

---

## 4. Техническая архитектура

### 4.1 Технологический стек

#### Backend:
- **PostgreSQL + PostGIS 3.4** - геопространственная БД
- **Go** - высокопроизводительные сервисы
- **Redis** - кэширование геоданных
- **Elasticsearch** - геопоиск

#### Frontend:
- **React 18** + **TypeScript**
- **Mapbox GL JS** - векторные карты
- **TanStack Query** - управление данными
- **Tailwind CSS** - стилизация

#### Инфраструктура:
- **Docker** + **Kubernetes**
- **Cloudflare CDN** - для карт
- **GitHub Actions** - CI/CD
- **Prometheus + Grafana** - мониторинг

### 4.2 Архитектура сервисов

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Web Client    │     │  Mobile Client  │     │   Admin Panel   │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                         │
         └───────────────────────┴─────────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │      API Gateway        │
                    │    (Kong/Nginx)         │
                    └────────────┬────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌────────▼────────┐    ┌────────▼────────┐    ┌────────▼────────┐
│   GIS Service   │    │ Delivery Service │    │ Search Service  │
│      (Go)       │    │      (Go)        │    │ (Elasticsearch) │
└────────┬────────┘    └────────┬────────┘    └────────┬────────┘
         │                       │                       │
         └───────────────────────┴───────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │      PostgreSQL         │
                    │    + PostGIS 3.4        │
                    └─────────────────────────┘
```

### 4.3 Модульная структура

```
backend/
├── internal/
│   └── proj/
│       └── gis/
│           ├── handler/
│           │   ├── spatial_search.go    # Пространственный поиск
│           │   ├── geocoding.go         # Геокодирование
│           │   ├── zones.go             # Зоны доставки
│           │   └── analytics.go         # Геоаналитика
│           ├── service/
│           │   ├── spatial_service.go   # Бизнес-логика
│           │   ├── clustering.go        # Кластеризация
│           │   └── routing.go           # Маршрутизация
│           ├── repository/
│           │   ├── postgis_repo.go      # PostGIS запросы
│           │   └── cache_repo.go        # Redis кэш
│           └── types/
│               └── geo_types.go         # Типы данных
└── migrations/
    ├── 001_enable_postgis.sql
    ├── 002_create_geo_tables.sql
    └── 003_add_spatial_indexes.sql

frontend/svetu/src/
├── components/
│   └── GIS/
│       ├── Map/
│       │   ├── InteractiveMap.tsx       # Основная карта
│       │   ├── MapControls.tsx          # Контролы карты
│       │   └── MapFilters.tsx           # Фильтры
│       ├── Markers/
│       │   ├── ListingMarker.tsx        # Маркер товара
│       │   └── ClusterMarker.tsx        # Кластер
│       ├── Search/
│       │   ├── RadiusSearch.tsx         # Поиск по радиусу
│       │   └── SmartSearch.tsx          # Умный поиск
│       └── Delivery/
│           ├── ZoneEditor.tsx           # Редактор зон
│           └── DeliverySelector.tsx     # Выбор доставки
└── hooks/
    └── gis/
        ├── useGeolocation.ts            # Геолокация
        ├── useGeoSearch.ts              # Поиск
        └── useDeliveryZones.ts          # Зоны доставки
```

---

## 5. База данных

### 5.1 Основные таблицы

```sql
-- Включение PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Геоданные объявлений
CREATE TABLE listings_geo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL REFERENCES marketplace_listings(id),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    -- Геохеши для быстрого поиска
    geohash4 VARCHAR(4) GENERATED ALWAYS AS (ST_GeoHash(location, 4)) STORED,
    geohash6 VARCHAR(6) GENERATED ALWAYS AS (ST_GeoHash(location, 6)) STORED,
    geohash8 VARCHAR(8) GENERATED ALWAYS AS (ST_GeoHash(location, 8)) STORED,
    -- Денормализованные данные
    city VARCHAR(100),
    district VARCHAR(100),
    postal_code VARCHAR(10),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_listings_geo_location ON listings_geo USING GIST(location);
CREATE INDEX idx_listings_geo_geohash4 ON listings_geo(geohash4);
CREATE INDEX idx_listings_geo_geohash6 ON listings_geo(geohash6);
CREATE INDEX idx_listings_geo_city ON listings_geo(city);
CREATE INDEX idx_listings_geo_listing ON listings_geo(listing_id);

-- Зоны доставки
CREATE TABLE delivery_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    storefront_id UUID NOT NULL REFERENCES user_storefronts(id),
    zone_name VARCHAR(100),
    zone_polygon GEOGRAPHY(POLYGON, 4326) NOT NULL,
    delivery_fee DECIMAL(10,2) DEFAULT 0,
    min_order_amount DECIMAL(10,2) DEFAULT 0,
    max_delivery_time_minutes INT,
    color VARCHAR(7), -- HEX цвет для отображения
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_delivery_zones_polygon ON delivery_zones USING GIST(zone_polygon);
CREATE INDEX idx_delivery_zones_storefront ON delivery_zones(storefront_id);

-- Службы доставки
CREATE TABLE delivery_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL, -- 'dexpress', 'posta_srbije'
    name VARCHAR(100) NOT NULL,
    name_cyrillic VARCHAR(100),
    api_endpoint VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    capabilities JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Пункты выдачи
CREATE TABLE pickup_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID REFERENCES delivery_providers(id),
    external_id VARCHAR(100),
    name VARCHAR(200) NOT NULL,
    address VARCHAR(500) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    city VARCHAR(100),
    postal_code VARCHAR(10),
    working_hours JSONB,
    services JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pickup_points_location ON pickup_points USING GIST(location);
CREATE INDEX idx_pickup_points_city ON pickup_points(city);
CREATE INDEX idx_pickup_points_provider ON pickup_points(provider_id);

-- Административные границы Сербии
CREATE TABLE administrative_boundaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    name_cyrillic VARCHAR(100),
    name_latin VARCHAR(100),
    type VARCHAR(50) NOT NULL, -- 'municipality', 'district', 'region'
    parent_id UUID REFERENCES administrative_boundaries(id),
    boundary GEOGRAPHY(MULTIPOLYGON, 4326),
    center_point GEOGRAPHY(POINT, 4326),
    population INT,
    area_km2 DECIMAL(10,2),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_boundaries_type ON administrative_boundaries(type);
CREATE INDEX idx_boundaries_boundary ON administrative_boundaries USING GIST(boundary);

-- Кэш геокодирования
CREATE TABLE geocoding_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address_hash VARCHAR(64) UNIQUE NOT NULL,
    original_address TEXT NOT NULL,
    normalized_address TEXT,
    location GEOGRAPHY(POINT, 4326),
    confidence_score FLOAT,
    provider VARCHAR(50),
    raw_response JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_used TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_geocoding_hash ON geocoding_cache(address_hash);
CREATE INDEX idx_geocoding_location ON geocoding_cache USING GIST(location);
```

### 5.2 Функции и процедуры

```sql
-- Поиск объявлений в радиусе
CREATE OR REPLACE FUNCTION search_listings_radius(
    p_center GEOGRAPHY,
    p_radius_meters INT,
    p_category_id INT DEFAULT NULL,
    p_price_min DECIMAL DEFAULT NULL,
    p_price_max DECIMAL DEFAULT NULL,
    p_limit INT DEFAULT 100
) RETURNS TABLE (
    listing_id UUID,
    distance_meters FLOAT,
    title VARCHAR,
    price DECIMAL,
    location GEOGRAPHY
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        lg.listing_id,
        ST_Distance(lg.location, p_center) as distance_meters,
        ml.title,
        ml.price,
        lg.location
    FROM listings_geo lg
    JOIN marketplace_listings ml ON lg.listing_id = ml.id
    WHERE 
        ST_DWithin(lg.location, p_center, p_radius_meters)
        AND ml.is_active = true
        AND (p_category_id IS NULL OR ml.category_id = p_category_id)
        AND (p_price_min IS NULL OR ml.price >= p_price_min)
        AND (p_price_max IS NULL OR ml.price <= p_price_max)
    ORDER BY distance_meters
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- Проверка доставки в точку
CREATE OR REPLACE FUNCTION check_delivery_availability(
    p_storefront_id UUID,
    p_location GEOGRAPHY
) RETURNS TABLE (
    zone_id UUID,
    zone_name VARCHAR,
    delivery_fee DECIMAL,
    min_order_amount DECIMAL,
    delivery_time INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        dz.id,
        dz.zone_name,
        dz.delivery_fee,
        dz.min_order_amount,
        dz.max_delivery_time_minutes
    FROM delivery_zones dz
    WHERE 
        dz.storefront_id = p_storefront_id
        AND dz.is_active = true
        AND ST_Contains(dz.zone_polygon::geometry, p_location::geometry);
END;
$$ LANGUAGE plpgsql;

-- Кластеризация для карты
CREATE OR REPLACE FUNCTION create_clusters(
    p_bounds GEOGRAPHY,
    p_zoom INT,
    p_category_id INT DEFAULT NULL
) RETURNS TABLE (
    cluster_center GEOGRAPHY,
    cluster_count INT,
    cluster_bounds GEOGRAPHY
) AS $$
DECLARE
    grid_size FLOAT;
BEGIN
    -- Размер сетки в зависимости от зума
    grid_size := CASE 
        WHEN p_zoom < 10 THEN 0.1
        WHEN p_zoom < 14 THEN 0.01
        ELSE 0.001
    END;
    
    RETURN QUERY
    SELECT 
        ST_Centroid(ST_Collect(lg.location::geometry))::geography as cluster_center,
        COUNT(*)::INT as cluster_count,
        ST_Envelope(ST_Collect(lg.location::geometry))::geography as cluster_bounds
    FROM listings_geo lg
    JOIN marketplace_listings ml ON lg.listing_id = ml.id
    WHERE 
        ST_Intersects(lg.location, p_bounds)
        AND ml.is_active = true
        AND (p_category_id IS NULL OR ml.category_id = p_category_id)
    GROUP BY 
        FLOOR(ST_X(lg.location::geometry) / grid_size),
        FLOOR(ST_Y(lg.location::geometry) / grid_size)
    HAVING COUNT(*) > 1;
END;
$$ LANGUAGE plpgsql;
```

---

## 6. API спецификация

### 6.1 Endpoints

#### Пространственный поиск
```yaml
POST /api/v1/gis/search/radius:
  summary: Поиск объявлений в радиусе
  requestBody:
    content:
      application/json:
        schema:
          type: object
          required: [lat, lng]
          properties:
            lat:
              type: number
              description: Широта центра поиска
            lng:
              type: number
              description: Долгота центра поиска
            radius:
              type: integer
              default: 5000
              description: Радиус поиска в метрах
            category_id:
              type: integer
              description: ID категории
            price_min:
              type: number
            price_max:
              type: number
            limit:
              type: integer
              default: 100
              maximum: 500
  responses:
    200:
      content:
        application/json:
          schema:
            type: object
            properties:
              results:
                type: array
                items:
                  $ref: '#/components/schemas/GeoListing'
              total:
                type: integer
              center:
                $ref: '#/components/schemas/Point'

GET /api/v1/gis/search/bounds:
  summary: Поиск в границах карты
  parameters:
    - name: bounds
      in: query
      required: true
      description: "sw_lat,sw_lng,ne_lat,ne_lng"
      schema:
        type: string
    - name: zoom
      in: query
      schema:
        type: integer
        minimum: 1
        maximum: 20
    - name: clustered
      in: query
      schema:
        type: boolean
        default: true

POST /api/v1/gis/geocode:
  summary: Геокодирование адреса
  requestBody:
    content:
      application/json:
        schema:
          type: object
          required: [address]
          properties:
            address:
              type: string
              description: Адрес для геокодирования
            country:
              type: string
              default: "RS"

GET /api/v1/gis/reverse:
  summary: Обратное геокодирование
  parameters:
    - name: lat
      in: query
      required: true
    - name: lng
      in: query
      required: true
```

#### Зоны доставки
```yaml
GET /api/v1/gis/delivery/zones/{storefront_id}:
  summary: Получить зоны доставки магазина
  
POST /api/v1/gis/delivery/zones:
  summary: Создать зону доставки
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/DeliveryZone'

PUT /api/v1/gis/delivery/zones/{zone_id}:
  summary: Обновить зону доставки

DELETE /api/v1/gis/delivery/zones/{zone_id}:
  summary: Удалить зону доставки

POST /api/v1/gis/delivery/check:
  summary: Проверить доступность доставки
  requestBody:
    content:
      application/json:
        schema:
          type: object
          required: [storefront_id, lat, lng]
          properties:
            storefront_id:
              type: string
              format: uuid
            lat:
              type: number
            lng:
              type: number
```

#### Службы доставки
```yaml
GET /api/v1/gis/delivery/providers:
  summary: Список служб доставки

POST /api/v1/gis/delivery/calculate:
  summary: Расчет стоимости доставки
  requestBody:
    content:
      application/json:
        schema:
          type: object
          required: [from_postal, to_postal, weight]
          properties:
            from_postal:
              type: string
            to_postal:
              type: string
            weight:
              type: number
              description: Вес в килограммах
            provider:
              type: string
              enum: [dexpress, posta_srbije]

GET /api/v1/gis/delivery/pickup-points:
  summary: Пункты выдачи
  parameters:
    - name: provider
      in: query
      schema:
        type: string
        enum: [dexpress, posta_srbije]
    - name: city
      in: query
    - name: lat
      in: query
    - name: lng
      in: query
    - name: radius
      in: query

GET /api/v1/gis/delivery/tracking/{tracking_number}:
  summary: Отслеживание посылки
```

#### Аналитика
```yaml
GET /api/v1/gis/analytics/heatmap:
  summary: Тепловая карта активности
  parameters:
    - name: type
      in: query
      schema:
        type: string
        enum: [purchases, views, searches]
    - name: period
      in: query
      schema:
        type: string
        enum: [day, week, month]

GET /api/v1/gis/analytics/districts:
  summary: Статистика по районам
  parameters:
    - name: storefront_id
      in: query
    - name: metric
      in: query
      schema:
        type: string
        enum: [orders, revenue, customers]
```

### 6.2 Модели данных

```typescript
// Основные типы
interface Point {
  lat: number;
  lng: number;
}

interface Bounds {
  southwest: Point;
  northeast: Point;
}

interface GeoListing {
  id: string;
  title: string;
  price: number;
  location: Point;
  distance?: number;
  category: {
    id: number;
    name: string;
  };
  thumbnail: string;
  seller: {
    id: string;
    name: string;
    rating: number;
  };
}

interface Cluster {
  id: string;
  center: Point;
  count: number;
  bounds: Bounds;
}

interface DeliveryZone {
  id?: string;
  storefront_id: string;
  name: string;
  polygon: Point[];
  delivery_fee: number;
  min_order_amount: number;
  max_delivery_time_minutes: number;
  color: string;
  is_active: boolean;
}

interface DeliveryProvider {
  id: string;
  code: 'dexpress' | 'posta_srbije';
  name: string;
  name_cyrillic: string;
  capabilities: string[];
}

interface PickupPoint {
  id: string;
  provider: DeliveryProvider;
  name: string;
  address: string;
  location: Point;
  working_hours: Record<string, string>;
  services: string[];
}

interface TrackingInfo {
  tracking_number: string;
  status: 'pending' | 'picked_up' | 'in_transit' | 'delivered';
  current_location: string;
  events: TrackingEvent[];
  estimated_delivery: string;
}
```

---

## 7. Frontend компоненты

### 7.1 Основная карта

```typescript
// src/components/GIS/Map/InteractiveMap.tsx
import { useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

interface InteractiveMapProps {
  center?: [number, number];
  zoom?: number;
  listings?: GeoListing[];
  onBoundsChange?: (bounds: Bounds) => void;
  onMarkerClick?: (listing: GeoListing) => void;
}

export const InteractiveMap: React.FC<InteractiveMapProps> = ({
  center = [20.4568, 44.8178], // Белград
  zoom = 12,
  listings = [],
  onBoundsChange,
  onMarkerClick
}) => {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    if (!mapContainer.current) return;
    
    // Инициализация карты
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v12',
      center,
      zoom,
      attributionControl: false
    });
    
    // Добавляем контролы
    map.current.addControl(new mapboxgl.NavigationControl(), 'top-right');
    map.current.addControl(new mapboxgl.GeolocateControl({
      positionOptions: { enableHighAccuracy: true },
      trackUserLocation: true
    }), 'top-right');
    
    // События карты
    map.current.on('load', () => {
      setIsLoading(false);
    });
    
    map.current.on('moveend', () => {
      if (onBoundsChange) {
        const bounds = map.current!.getBounds();
        onBoundsChange({
          southwest: { 
            lat: bounds.getSouth(), 
            lng: bounds.getWest() 
          },
          northeast: { 
            lat: bounds.getNorth(), 
            lng: bounds.getEast() 
          }
        });
      }
    });
    
    return () => {
      map.current?.remove();
    };
  }, []);
  
  // Обновление маркеров
  useEffect(() => {
    if (!map.current || isLoading) return;
    
    // Удаляем старые маркеры
    const markers = document.getElementsByClassName('mapboxgl-marker');
    while (markers[0]) {
      markers[0].remove();
    }
    
    // Добавляем новые маркеры
    listings.forEach(listing => {
      const el = document.createElement('div');
      el.className = 'custom-marker';
      el.style.width = '30px';
      el.style.height = '30px';
      el.style.backgroundImage = `url(/api/categories/${listing.category.id}/icon)`;
      el.style.backgroundSize = 'cover';
      el.style.borderRadius = '50%';
      el.style.border = '2px solid white';
      el.style.boxShadow = '0 2px 4px rgba(0,0,0,0.3)';
      el.style.cursor = 'pointer';
      
      // Добавляем эффект размытия для приватных локаций
      if (listing.is_blurred) {
        el.style.opacity = '0.7';
        el.classList.add('blurred-location');
      }
      
      const marker = new mapboxgl.Marker(el)
        .setLngLat([listing.location.lng, listing.location.lat])
        .addTo(map.current!);
        
      el.addEventListener('click', () => {
        if (onMarkerClick) {
          onMarkerClick(listing);
        }
      });
      
      // Попап с информацией
      const popup = new mapboxgl.Popup({ offset: 25 })
        .setHTML(`
          <div class="p-2">
            <img src="${listing.thumbnail}" class="w-full h-24 object-cover rounded mb-2" />
            <h3 class="font-semibold">${listing.title}</h3>
            <p class="text-lg font-bold">${listing.price} РСД</p>
            <p class="text-sm text-gray-600">${listing.distance ? `${listing.distance}м` : ''}</p>
          </div>
        `);
        
      marker.setPopup(popup);
    });
  }, [listings, isLoading]);
  
  return (
    <div className="relative w-full h-full">
      <div ref={mapContainer} className="w-full h-full" />
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      )}
    </div>
  );
};
```

### 7.2 Поиск по радиусу

```typescript
// src/components/GIS/Search/RadiusSearch.tsx
import { useState } from 'react';
import { useGeolocation } from '@/hooks/useGeolocation';
import { useGeoSearch } from '@/hooks/useGeoSearch';

export const RadiusSearch: React.FC = () => {
  const { location, error: geoError, requestLocation } = useGeolocation();
  const { search, isLoading } = useGeoSearch();
  const [radius, setRadius] = useState(5); // км
  const [filters, setFilters] = useState({
    category_id: null,
    price_min: null,
    price_max: null
  });
  
  const handleSearch = () => {
    if (!location) {
      requestLocation();
      return;
    }
    
    search({
      lat: location.coords.latitude,
      lng: location.coords.longitude,
      radius: radius * 1000, // в метры
      ...filters
    });
  };
  
  return (
    <div className="bg-white p-4 rounded-lg shadow-lg">
      <h3 className="text-lg font-semibold mb-4">Поиск рядом с вами</h3>
      
      {/* Радиус поиска */}
      <div className="mb-4">
        <label className="block text-sm font-medium mb-2">
          Радиус поиска: {radius} км
        </label>
        <input
          type="range"
          min="1"
          max="50"
          value={radius}
          onChange={(e) => setRadius(Number(e.target.value))}
          className="w-full"
        />
      </div>
      
      {/* Фильтры */}
      <CategoryFilter 
        value={filters.category_id}
        onChange={(cat) => setFilters({...filters, category_id: cat})}
      />
      
      <PriceRangeFilter
        min={filters.price_min}
        max={filters.price_max}
        onChange={(min, max) => setFilters({...filters, price_min: min, price_max: max})}
      />
      
      {/* Кнопка поиска */}
      <button
        onClick={handleSearch}
        disabled={isLoading}
        className="w-full mt-4 bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50"
      >
        {isLoading ? 'Поиск...' : 'Найти'}
      </button>
      
      {geoError && (
        <p className="text-red-600 text-sm mt-2">
          Не удалось определить местоположение
        </p>
      )}
    </div>
  );
};
```

### 7.3 Редактор зон доставки

```typescript
// src/components/GIS/Delivery/ZoneEditor.tsx
import { useState, useRef } from 'react';
import { MapboxDraw } from '@mapbox/mapbox-gl-draw';
import '@mapbox/mapbox-gl-draw/dist/mapbox-gl-draw.css';

interface DeliveryZoneEditorProps {
  storefrontId: string;
  existingZones?: DeliveryZone[];
  onSave: (zones: DeliveryZone[]) => void;
}

export const DeliveryZoneEditor: React.FC<DeliveryZoneEditorProps> = ({
  storefrontId,
  existingZones = [],
  onSave
}) => {
  const [zones, setZones] = useState<DeliveryZone[]>(existingZones);
  const [selectedZone, setSelectedZone] = useState<string | null>(null);
  const drawRef = useRef<MapboxDraw | null>(null);
  
  useEffect(() => {
    if (!map.current) return;
    
    // Инициализация инструментов рисования
    drawRef.current = new MapboxDraw({
      displayControlsDefault: false,
      controls: {
        polygon: true,
        trash: true
      },
      defaultMode: 'draw_polygon'
    });
    
    map.current.addControl(drawRef.current);
    
    // Загружаем существующие зоны
    existingZones.forEach(zone => {
      drawRef.current!.add({
        type: 'Feature',
        id: zone.id,
        geometry: {
          type: 'Polygon',
          coordinates: [zone.polygon.map(p => [p.lng, p.lat])]
        },
        properties: zone
      });
    });
    
    // События рисования
    map.current.on('draw.create', (e) => {
      const feature = e.features[0];
      const polygon = feature.geometry.coordinates[0].map(coord => ({
        lat: coord[1],
        lng: coord[0]
      }));
      
      const newZone: DeliveryZone = {
        id: feature.id,
        storefront_id: storefrontId,
        name: `Зона ${zones.length + 1}`,
        polygon,
        delivery_fee: 0,
        min_order_amount: 0,
        max_delivery_time_minutes: 60,
        color: generateColor(),
        is_active: true
      };
      
      setZones([...zones, newZone]);
      setSelectedZone(newZone.id);
    });
    
    return () => {
      if (drawRef.current && map.current) {
        map.current.removeControl(drawRef.current);
      }
    };
  }, []);
  
  const updateZone = (zoneId: string, updates: Partial<DeliveryZone>) => {
    setZones(zones.map(z => 
      z.id === zoneId ? { ...z, ...updates } : z
    ));
  };
  
  const deleteZone = (zoneId: string) => {
    drawRef.current?.delete(zoneId);
    setZones(zones.filter(z => z.id !== zoneId));
    setSelectedZone(null);
  };
  
  return (
    <div className="flex h-full">
      {/* Карта */}
      <div className="flex-1">
        <InteractiveMap />
      </div>
      
      {/* Панель настроек */}
      <div className="w-80 bg-white p-4 shadow-lg overflow-y-auto">
        <h3 className="text-lg font-semibold mb-4">Зоны доставки</h3>
        
        {zones.length === 0 ? (
          <p className="text-gray-600 text-sm">
            Нарисуйте зону на карте для начала
          </p>
        ) : (
          <div className="space-y-4">
            {zones.map(zone => (
              <div
                key={zone.id}
                className={`border rounded-lg p-3 cursor-pointer ${
                  selectedZone === zone.id ? 'border-blue-500' : 'border-gray-200'
                }`}
                onClick={() => setSelectedZone(zone.id)}
              >
                <div className="flex items-center justify-between mb-2">
                  <input
                    type="text"
                    value={zone.name}
                    onChange={(e) => updateZone(zone.id!, { name: e.target.value })}
                    className="font-semibold bg-transparent"
                    onClick={(e) => e.stopPropagation()}
                  />
                  <div
                    className="w-6 h-6 rounded"
                    style={{ backgroundColor: zone.color }}
                  />
                </div>
                
                <div className="space-y-2 text-sm">
                  <div>
                    <label className="text-gray-600">Стоимость доставки:</label>
                    <input
                      type="number"
                      value={zone.delivery_fee}
                      onChange={(e) => updateZone(zone.id!, { 
                        delivery_fee: Number(e.target.value) 
                      })}
                      className="ml-2 w-20 border rounded px-1"
                      onClick={(e) => e.stopPropagation()}
                    /> РСД
                  </div>
                  
                  <div>
                    <label className="text-gray-600">Мин. сумма:</label>
                    <input
                      type="number"
                      value={zone.min_order_amount}
                      onChange={(e) => updateZone(zone.id!, { 
                        min_order_amount: Number(e.target.value) 
                      })}
                      className="ml-2 w-20 border rounded px-1"
                      onClick={(e) => e.stopPropagation()}
                    /> РСД
                  </div>
                  
                  <div>
                    <label className="text-gray-600">Время доставки:</label>
                    <input
                      type="number"
                      value={zone.max_delivery_time_minutes}
                      onChange={(e) => updateZone(zone.id!, { 
                        max_delivery_time_minutes: Number(e.target.value) 
                      })}
                      className="ml-2 w-16 border rounded px-1"
                      onClick={(e) => e.stopPropagation()}
                    /> мин
                  </div>
                </div>
                
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteZone(zone.id!);
                  }}
                  className="mt-2 text-red-600 text-sm hover:underline"
                >
                  Удалить зону
                </button>
              </div>
            ))}
          </div>
        )}
        
        <button
          onClick={() => onSave(zones)}
          className="mt-6 w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700"
        >
          Сохранить зоны
        </button>
      </div>
    </div>
  );
};
```

### 7.4 Выбор службы доставки

```typescript
// src/components/GIS/Delivery/ServiceSelector.tsx
export const DeliveryServiceSelector: React.FC<{
  from: Point;
  to: Point;
  weight: number;
  onSelect: (option: DeliveryOption) => void;
}> = ({ from, to, weight, onSelect }) => {
  const { data: options, isLoading } = useQuery({
    queryKey: ['delivery-options', from, to, weight],
    queryFn: () => api.calculateDeliveryOptions({ from, to, weight })
  });
  
  if (isLoading) return <LoadingSpinner />;
  
  return (
    <div className="space-y-4">
      {/* D Express */}
      {options?.dexpress && (
        <div 
          className="border rounded-lg p-4 cursor-pointer hover:border-blue-500 transition-colors"
          onClick={() => onSelect(options.dexpress)}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <img 
                src="/images/dexpress-logo.png" 
                alt="D Express" 
                className="h-10"
              />
              <div>
                <h4 className="font-semibold">D Express</h4>
                <p className="text-sm text-gray-600">
                  Брза достава за 1-2 дана
                </p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold">{options.dexpress.price} РСД</p>
              <p className="text-sm text-gray-600">1-2 дана</p>
            </div>
          </div>
        </div>
      )}
      
      {/* Пошта Србије */}
      {options?.posta_srbije && (
        <div 
          className="border rounded-lg p-4 cursor-pointer hover:border-blue-500 transition-colors"
          onClick={() => onSelect(options.posta_srbije)}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <img 
                src="/images/posta-srbije-logo.png" 
                alt="Пошта Србије" 
                className="h-10"
              />
              <div>
                <h4 className="font-semibold">Пошта Србије</h4>
                <p className="text-sm text-gray-600">
                  Стандардна достава
                </p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold">{options.posta_srbije.price} РСД</p>
              <p className="text-sm text-gray-600">2-5 дана</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
```

---

## 8. Интеграция с доставкой

### 8.1 D Express

```go
// internal/proj/delivery/providers/dexpress/client.go
package dexpress

type Client struct {
    apiURL   string
    apiKey   string
    client   *http.Client
    cache    *redis.Client
}

func (c *Client) CalculatePrice(params PriceParams) (*PriceResponse, error) {
    // Проверка кэша
    cacheKey := fmt.Sprintf("dexpress:price:%s:%s:%.2f", 
        params.FromPostal, params.ToPostal, params.Weight)
    
    if cached, err := c.cache.Get(context.Background(), cacheKey).Result(); err == nil {
        var resp PriceResponse
        json.Unmarshal([]byte(cached), &resp)
        return &resp, nil
    }
    
    // API запрос
    reqBody, _ := json.Marshal(map[string]interface{}{
        "from_postal": params.FromPostal,
        "to_postal":   params.ToPostal,
        "weight":      params.Weight,
        "cod_amount":  params.CODAmount,
    })
    
    req, _ := http.NewRequest("POST", 
        c.apiURL+"/calculate-price", 
        bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var priceResp PriceResponse
    json.NewDecoder(resp.Body).Decode(&priceResp)
    
    // Кэшируем на 1 час
    c.cache.Set(context.Background(), cacheKey, 
        json.Marshal(priceResp), time.Hour)
    
    return &priceResp, nil
}

func (c *Client) CreateShipment(order ShipmentOrder) (*Shipment, error) {
    reqBody, _ := json.Marshal(order)
    
    req, _ := http.NewRequest("POST", 
        c.apiURL+"/shipments", 
        bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    
    var shipment Shipment
    json.NewDecoder(resp.Body).Decode(&shipment)
    
    return &shipment, nil
}

func (c *Client) TrackShipment(trackingNumber string) (*TrackingInfo, error) {
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("%s/tracking/%s", c.apiURL, trackingNumber), nil)
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    
    var tracking TrackingInfo
    json.NewDecoder(resp.Body).Decode(&tracking)
    
    // Маппинг статусов
    tracking.Status = c.mapStatus(tracking.Status)
    
    return &tracking, nil
}
```

### 8.2 Пошта Србије

```go
// internal/proj/delivery/providers/posta/client.go
package posta

type Client struct {
    apiURL   string
    username string
    password string
    client   *http.Client
    db       *sql.DB
}

func (c *Client) CalculatePrice(params PriceParams) (*PriceResponse, error) {
    // Пошта Србије использует зоны
    zone := c.calculateZone(params.FromPostal, params.ToPostal)
    
    // Базовая цена по зоне
    basePrice := map[int]float64{
        1: 250, // внутри города
        2: 350, // между городами
        3: 450, // удаленные районы
    }[zone]
    
    // Доплата за вес
    weightPrice := 0.0
    if params.Weight > 1 {
        weightPrice = (params.Weight - 1) * 50
    }
    
    return &PriceResponse{
        Price:         basePrice + weightPrice,
        EstimatedDays: c.getEstimatedDays(zone),
        Zone:          zone,
    }, nil
}

func (c *Client) calculateZone(from, to string) int {
    // Логика определения зоны по почтовым индексам
    fromRegion := c.getRegionByPostal(from)
    toRegion := c.getRegionByPostal(to)
    
    if fromRegion == toRegion {
        return 1 // внутри региона
    } else if c.isMajorCity(fromRegion) && c.isMajorCity(toRegion) {
        return 2 // между крупными городами
    }
    return 3 // остальное
}

func (c *Client) FindNearestOffices(lat, lng float64, limit int) ([]PostOffice, error) {
    query := `
        SELECT 
            id, name, address, postal_code,
            ST_Distance(location::geography, ST_Point($1, $2)::geography) as distance,
            working_hours, services
        FROM pickup_points
        WHERE 
            provider_id = (SELECT id FROM delivery_providers WHERE code = 'posta_srbije')
            AND is_active = true
        ORDER BY location <-> ST_Point($1, $2)::geography
        LIMIT $3
    `
    
    rows, err := c.db.Query(query, lng, lat, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var offices []PostOffice
    for rows.Next() {
        var office PostOffice
        var servicesJSON []byte
        var hoursJSON []byte
        
        err := rows.Scan(
            &office.ID, &office.Name, &office.Address, 
            &office.PostalCode, &office.Distance,
            &hoursJSON, &servicesJSON,
        )
        if err != nil {
            continue
        }
        
        json.Unmarshal(hoursJSON, &office.WorkingHours)
        json.Unmarshal(servicesJSON, &office.Services)
        
        offices = append(offices, office)
    }
    
    return offices, nil
}
```

---

## 9. План реализации

### 9.1 Команда и ресурсы

| Роль | Занятость | Ответственность |
|------|-----------|-----------------|
| Backend Senior | 100% | PostGIS, API, интеграции |
| Frontend Senior | 100% | Карты, UI компоненты |
| Full-stack Middle+ | 100% | Доставка, аналитика |
| DevOps | 30% | Инфраструктура, мониторинг |
| UI/UX Designer | 50% (4 недели) | Дизайн компонентов |

### 9.2 Фазы разработки

#### Фаза 1: Базовая карта (4 недели)
**Цель**: Полностью работающая интерактивная карта с визуализацией товаров и витрин

**Результат фазы**: Пользователи могут открыть карту, увидеть товары рядом, кликнуть на маркер и перейти к товару

**Неделя 1: Инфраструктура и данные**
- [ ] Настройка PostGIS в Docker
- [ ] Создание таблиц listings_geo
- [ ] Миграция координат существующих объявлений
- [ ] Базовая структура GIS модуля

**Неделя 2: Backend API**
- [ ] GET /api/v1/gis/search/bounds - поиск в границах карты
- [ ] GET /api/v1/gis/listings/{id}/location - координаты товара
- [ ] POST /api/v1/gis/listings/{id}/location - обновление координат
- [ ] Интеграция с Redis для кэширования

**Неделя 3: Frontend карта**  
- [ ] Интеграция Mapbox GL JS
- [ ] Компонент InteractiveMap с маркерами
- [ ] Попапы с превью товара
- [ ] Базовые контролы карты (зум, навигация)

**Неделя 4: Интеграция и оптимизация**
- [ ] Кластеризация маркеров при большом количестве
- [ ] Lazy loading маркеров по мере движения карты
- [ ] Мобильная адаптация карты
- [ ] Интеграционное тестирование

**Метрики успеха фазы 1**:
- Карта загружается < 2 сек
- Отображается 1000+ маркеров без лагов
- CTR с маркера на товар > 5%

#### Фаза 2: Геопоиск и фильтры (4 недели)
**Цель**: Полноценный поиск товаров по местоположению с фильтрами

**Результат фазы**: Пользователи могут искать товары в радиусе от себя, по районам, с фильтрами по категориям и цене

**Неделя 5-6: Пространственный поиск**
- [ ] POST /api/v1/gis/search/radius - поиск в радиусе
- [ ] GET /api/v1/gis/districts - список районов
- [ ] POST /api/v1/gis/search/district - поиск по району
- [ ] Оптимизация PostGIS индексов

**Неделя 7: Геолокация и геокодирование**
- [ ] Hook useGeolocation для определения местоположения
- [ ] POST /api/v1/gis/geocode - адрес в координаты
- [ ] GET /api/v1/gis/reverse - координаты в адрес
- [ ] Кэширование результатов геокодирования

**Неделя 8: UI компоненты поиска**
- [ ] Компонент RadiusSearch с слайдером радиуса
- [ ] Компонент DistrictSelector с выбором района
- [ ] Интеграция фильтров категорий и цены
- [ ] Обновление результатов в реальном времени

**Метрики успеха фазы 2**:
- Поиск по радиусу < 200ms
- Использование геопоиска > 30% пользователей
- Конверсия из геопоиска > 10%

#### Фаза 3: Умный поиск и аналитика (4 недели)  
**Цель**: AI-powered поиск и аналитика для продавцов

**Результат фазы**: Умный поиск с рекомендациями, тепловые карты активности, статистика по районам

**Неделя 9-10: Умный поиск**
- [ ] Интеграция с Elasticsearch для геопоиска
- [ ] POST /api/v1/gis/search/smart - умный поиск
- [ ] Компонент SmartSearch с автокомплитом
- [ ] Персонализация результатов по истории

**Неделя 11: Геоаналитика**  
- [ ] GET /api/v1/gis/analytics/heatmap - тепловая карта
- [ ] GET /api/v1/gis/analytics/districts - статистика районов
- [ ] Компонент HeatmapLayer для карты
- [ ] Дашборд продавца с геостатистикой

**Неделя 12: Доставка и финализация**
- [ ] Базовые зоны доставки для витрин
- [ ] Проверка доставки по адресу
- [ ] Оптимизация производительности
- [ ] Production deploy и мониторинг

**Метрики успеха фазы 3**:
- Умный поиск используют > 50% пользователей
- Настроено > 100 зон доставки
- Uptime > 99.9%

### 9.3 Риски и митигация

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|  
| Задержка PostGIS | Низкая | Высокое | Готовый Docker образ |
| API лимиты карт | Средняя | Среднее | Кэширование, OSM fallback |
| Сложность UX | Низкая | Высокое | Итеративный дизайн |
| Производительность | Средняя | Высокое | Профилирование с день 1 |

---

## 10. Метрики успеха

### 10.1 Технические метрики

| Метрика | Цель | Измерение |
|---------|------|-----------|
| Uptime | 99.9% | Prometheus |
| P95 latency геопоиска | < 200ms | Grafana |
| Загрузка карты | < 2 сек | Google Analytics |
| Ошибки API | < 0.1% | Sentry |

### 10.2 Продуктовые метрики

| Метрика | Базовое | Цель (3 мес) |
|---------|---------|--------------|
| Использование карты | 0% | 40% |
| Конверсия из карты | - | 8% |
| Поиск по геолокации | 0% | 25% |
| Настроенные зоны доставки | 0 | 500 |

### 10.3 Бизнес метрики

| Метрика | Влияние |
|---------|---------|
| Рост локальных продаж | +30% |
| Снижение возвратов | -15% |
| Увеличение среднего чека | +20% |
| NPS продавцов | +15 пунктов |

---

## 📎 Приложения

### A. Конфигурация

```yaml
# .env.example
# Mapbox
MAPBOX_ACCESS_TOKEN=pk.xxx
MAPBOX_STYLE_URL=mapbox://styles/mapbox/streets-v12

# PostGIS
POSTGIS_VERSION=3.4
POSTGIS_ENABLE_RASTER=true

# Геокодирование
NOMINATIM_URL=https://nominatim.openstreetmap.org
GEOCODING_CACHE_TTL=86400

# D Express
DEXPRESS_API_URL=https://api.dexpress.rs/v1
DEXPRESS_API_KEY=xxx

# Пошта Србије
POSTA_API_URL=https://api.posta.rs/v2
POSTA_USERNAME=xxx
POSTA_PASSWORD=xxx

# Кэширование
REDIS_GEO_TTL=300
TILE_CACHE_SIZE=1GB
```

### B. Тестовые данные

```sql
-- Тестовые объявления для Белграда
INSERT INTO listings_geo (listing_id, location, city, district)
SELECT 
    id,
    ST_Point(
        20.4568 + (random() - 0.5) * 0.1,
        44.8178 + (random() - 0.5) * 0.1
    )::geography,
    'Београд',
    CASE (random() * 5)::int
        WHEN 0 THEN 'Стари град'
        WHEN 1 THEN 'Нови Београд'
        WHEN 2 THEN 'Земун'
        WHEN 3 THEN 'Палилула'
        ELSE 'Чукарица'
    END
FROM marketplace_listings
WHERE city = 'Београд'
LIMIT 1000;

-- Тестовые почтовые отделения
INSERT INTO pickup_points (provider_id, name, address, location, city, postal_code)
VALUES
    ((SELECT id FROM delivery_providers WHERE code = 'posta_srbije'),
     'Пошта 11000 Београд 1', 'Таковска 2', ST_Point(20.4589, 44.8098)::geography, 
     'Београд', '11000'),
    ((SELECT id FROM delivery_providers WHERE code = 'posta_srbije'),
     'Пошта 11070 Нови Београд', 'Булевар Михајла Пупина 165', ST_Point(20.4023, 44.8219)::geography, 
     'Београд', '11070');
```

---

## ✅ Контрольный список готовности

### Pre-development
- [ ] Команда сформирована
- [ ] Доступы к API получены
- [ ] Окружение настроено
- [ ] Дизайн-макеты готовы

### Development
- [ ] PostGIS работает
- [ ] API endpoints готовы
- [ ] Карта отображается
- [ ] Поиск функционирует
- [ ] Доставка интегрирована

### Pre-launch
- [ ] Тесты проходят
- [ ] Документация написана
- [ ] Мониторинг настроен
- [ ] Команда обучена

### Launch
- [ ] Production deploy
- [ ] A/B тест запущен
- [ ] Метрики собираются
- [ ] Support готов

---

**Документ готов к передаче в разработку!** 🚀

При возникновении вопросов обращайтесь к техническому руководителю проекта.