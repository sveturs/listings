# 🗺️ ПОЛНЫЙ ПЛАН РАЗВИТИЯ СИСТЕМЫ КАРТ SVETU

## 📊 РЕАЛЬНОЕ СОСТОЯНИЕ СИСТЕМЫ (13.09.2025)

### ✅ Что уже реализовано и работает

#### Инфраструктура БД
- **PostGIS** полностью настроен (миграция 000001)
- **15+ GIST индексов** для оптимизации geo-запросов
- **Таблицы созданы**:
  - `listings_geo` - геоданные объявлений (структура готова, 0 записей)
  - `gis_listing_density_grid` - сетка плотности (901,901 записей)
  - `unified_geo` - унифицированные геоданные (74 записи)
  - `geocoding_cache` - кэш геокодирования (1 запись)
  - `marketplace_listings` - 15 записей с координатами

#### Backend архитектура
```
backend/internal/proj/gis/
├── handler/          ✅ 6 файлов
├── service/          ✅ 6 сервисов
├── repository/       ✅ PostGIS репозиторий
└── types/            ✅ Типы данных
```

#### Frontend компоненты
```
frontend/src/components/GIS/
├── Map/              ✅ 17 компонентов
├── hooks/            ✅ 3 хука (геолокация, поиск, радиус)
├── utils/            ✅ GeoJSON утилиты
└── LocationPicker/   ✅ Выбор локации
```

#### API эндпоинты (работают, но не в Swagger)
- `/api/v1/gis/search/radius` - радиусный поиск
- `/api/v1/gis/geocode/*` - геокодирование
- `/api/v1/gis/nearby` - ближайшие объекты

### ⚠️ Критические проблемы

1. **Нет данных**: `listings_geo` пуста - геоданные не синхронизированы
2. **API не документирован**: GIS эндпоинты отсутствуют в Swagger
3. **Лимиты малы**: API отдает максимум 100 записей
4. **Нет кэширования**: Frontend загружает данные при каждом изменении

---

## 🚀 ФАЗА 0: КРИТИЧЕСКИЕ ИСПРАВЛЕНИЯ (1-2 дня)

### 0.1 Синхронизация геоданных

**Создать миграцию для заполнения `listings_geo`:**

```sql
-- backend/migrations/000XXX_sync_listings_geo.up.sql
INSERT INTO listings_geo (
    listing_id,
    location,
    blurred_location,
    privacy_level,
    address_components
)
SELECT
    id,
    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326),
    -- Размытие на 500м для приватности
    ST_SetSRID(
        ST_MakePoint(
            longitude + (random() - 0.5) * 0.009,
            latitude + (random() - 0.5) * 0.009
        ),
        4326
    ),
    CASE
        WHEN user_id IN (SELECT id FROM users WHERE is_business = true) THEN 1
        ELSE 2
    END,
    jsonb_build_object(
        'city', city,
        'district', district,
        'street', address,
        'formatted', CONCAT(address, ', ', city)
    )
FROM marketplace_listings
WHERE latitude IS NOT NULL
  AND longitude IS NOT NULL
ON CONFLICT (listing_id) DO UPDATE
SET
    location = EXCLUDED.location,
    updated_at = NOW();

-- Обновить unified_geo
INSERT INTO unified_geo (
    entity_type,
    entity_id,
    location,
    metadata
)
SELECT
    'listing',
    id,
    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326),
    jsonb_build_object(
        'title', title,
        'category_id', category_id,
        'price', price
    )
FROM marketplace_listings
WHERE latitude IS NOT NULL
ON CONFLICT (entity_type, entity_id) DO UPDATE
SET
    location = EXCLUDED.location,
    metadata = EXCLUDED.metadata;
```

### 0.2 Увеличение лимитов API

```go
// backend/internal/proj/gis/handler/spatial_search.go
const (
    DEFAULT_LIMIT = 1000  // было 100
    MAX_LIMIT = 5000      // было 500
    DEFAULT_RADIUS = 5000 // 5км по умолчанию
)

func (h *Handler) validateSearchParams(params *SearchParams) error {
    if params.Limit <= 0 {
        params.Limit = DEFAULT_LIMIT
    }
    if params.Limit > MAX_LIMIT {
        params.Limit = MAX_LIMIT
    }
    // Добавить пагинацию для больших объемов
    if params.Offset < 0 {
        params.Offset = 0
    }
    return nil
}
```

### 0.3 Frontend кэширование и дебаунсинг

```typescript
// frontend/svetu/src/hooks/useMapCache.ts
import { useRef, useCallback } from 'react';
import { debounce } from 'lodash';

interface CacheEntry<T> {
  data: T;
  timestamp: number;
  key: string;
}

export function useMapCache<T>(ttl: number = 300000) { // 5 минут
  const cache = useRef<Map<string, CacheEntry<T>>>(new Map());

  const getCached = useCallback((key: string): T | null => {
    const entry = cache.current.get(key);
    if (!entry) return null;

    if (Date.now() - entry.timestamp > ttl) {
      cache.current.delete(key);
      return null;
    }

    return entry.data;
  }, [ttl]);

  const setCached = useCallback((key: string, data: T) => {
    // LRU: удалить старые если > 100 записей
    if (cache.current.size > 100) {
      const oldestKey = Array.from(cache.current.entries())
        .sort((a, b) => a[1].timestamp - b[1].timestamp)[0][0];
      cache.current.delete(oldestKey);
    }

    cache.current.set(key, {
      data,
      timestamp: Date.now(),
      key
    });
  }, []);

  const clearCache = useCallback(() => {
    cache.current.clear();
  }, []);

  return { getCached, setCached, clearCache };
}

// frontend/svetu/src/app/[locale]/map/MapClient.tsx
export default function MapClient() {
  const { getCached, setCached } = useMapCache();

  // Дебаунс для фильтров
  const debouncedLoadListings = useMemo(
    () => debounce(async (filters: any, bounds: any) => {
      const cacheKey = JSON.stringify({ filters, bounds });

      // Проверить кэш
      const cached = getCached(cacheKey);
      if (cached) {
        setListings(cached);
        return;
      }

      // Загрузить новые данные
      const response = await apiClient.get('/api/v1/gis/search/radius', {
        params: { ...filters, ...bounds, limit: 1000 }
      });

      setCached(cacheKey, response.data);
      setListings(response.data);
    }, 300),
    [getCached, setCached]
  );

  // Использовать дебаунс версию
  useEffect(() => {
    debouncedLoadListings(filters, mapBounds);
  }, [filters, mapBounds]);
}
```

### 0.4 Добавление GIS в Swagger

```go
// backend/internal/proj/gis/handler/routes.go
// Добавить Swagger аннотации ко всем эндпоинтам

// SearchRadius godoc
// @Summary Поиск объявлений в радиусе
// @Description Поиск объявлений в заданном радиусе от точки
// @Tags gis
// @Accept json
// @Produce json
// @Param lat query number true "Широта центра поиска"
// @Param lng query number true "Долгота центра поиска"
// @Param radius query number false "Радиус поиска в метрах" default(5000)
// @Param limit query integer false "Лимит результатов" default(1000)
// @Param category_id query integer false "ID категории для фильтрации"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.GeoListing}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Router /api/v1/gis/search/radius [get]
func (h *Handler) SearchRadius(c *fiber.Ctx) error {
    // существующий код
}
```

---

## 🎯 ФАЗА 1: ОПТИМИЗАЦИЯ ПРОИЗВОДИТЕЛЬНОСТИ (3-5 дней)

### 1.1 Активация серверной кластеризации

```go
// backend/internal/proj/gis/handler/cluster.go
package handler

import (
    "encoding/json"
    "github.com/gofiber/fiber/v2"
)

type ClusterPoint struct {
    Lat   float64 `json:"lat"`
    Lng   float64 `json:"lng"`
    Count int     `json:"count"`
    IDs   []int   `json:"ids,omitempty"`
}

// GetClusters godoc
// @Summary Получить кластеры объявлений
// @Tags gis
// @Param zoom query integer true "Уровень зума карты (1-20)"
// @Param bounds query string true "Границы видимой области (west,south,east,north)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]ClusterPoint}
// @Router /api/v1/gis/clusters [get]
func (h *Handler) GetClusters(c *fiber.Ctx) error {
    zoom := c.QueryInt("zoom", 10)
    bounds := c.Query("bounds")

    // Использовать готовую таблицу gis_listing_density_grid
    query := `
        WITH grid_size AS (
            SELECT
                CASE
                    WHEN $1 < 8 THEN 0.5    -- 50км сетка
                    WHEN $1 < 12 THEN 0.1   -- 10км сетка
                    WHEN $1 < 15 THEN 0.01  -- 1км сетка
                    ELSE 0.001              -- 100м сетка
                END as size
        ),
        clusters AS (
            SELECT
                ST_X(ST_Centroid(cell)) as lng,
                ST_Y(ST_Centroid(cell)) as lat,
                SUM(density) as count,
                array_agg(listing_ids) as ids
            FROM gis_listing_density_grid, grid_size
            WHERE cell && ST_MakeEnvelope($2, $3, $4, $5, 4326)
            GROUP BY
                floor(ST_X(ST_Centroid(cell)) / size) * size,
                floor(ST_Y(ST_Centroid(cell)) / size) * size
            HAVING SUM(density) > 0
        )
        SELECT * FROM clusters
        ORDER BY count DESC
        LIMIT 500
    `

    // Парсинг bounds и выполнение запроса
    west, south, east, north := parseBounds(bounds)

    rows, err := h.db.Query(query, zoom, west, south, east, north)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    var clusters []ClusterPoint
    for rows.Next() {
        var cluster ClusterPoint
        var ids []byte
        err := rows.Scan(&cluster.Lng, &cluster.Lat, &cluster.Count, &ids)
        if err != nil {
            continue
        }

        // На больших зумах включать IDs
        if zoom > 15 {
            json.Unmarshal(ids, &cluster.IDs)
        }

        clusters = append(clusters, cluster)
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data": clusters,
        "total": len(clusters),
    })
}
```

### 1.2 Виртуализация маркеров на Frontend

```typescript
// frontend/svetu/src/components/GIS/Map/VirtualizedMarkers.tsx
import { memo, useMemo } from 'react';
import { Marker } from 'react-map-gl';
import type { MapMarkerData } from '@/types/gis';

interface Props {
  markers: MapMarkerData[];
  viewport: {
    latitude: number;
    longitude: number;
    zoom: number;
    bounds: [number, number, number, number]; // [west, south, east, north]
  };
  onMarkerClick: (marker: MapMarkerData) => void;
}

export const VirtualizedMarkers = memo(({
  markers,
  viewport,
  onMarkerClick
}: Props) => {
  // Фильтровать только видимые маркеры + буфер
  const visibleMarkers = useMemo(() => {
    const [west, south, east, north] = viewport.bounds;
    const buffer = 0.01; // ~1км буфер

    return markers.filter(m => {
      const [lng, lat] = m.position;
      return (
        lat >= south - buffer &&
        lat <= north + buffer &&
        lng >= west - buffer &&
        lng <= east + buffer
      );
    });
  }, [markers, viewport.bounds]);

  // Использовать кластеры на малых зумах
  const renderMarkers = useMemo(() => {
    if (viewport.zoom < 12) {
      // Показывать кластеры
      return null; // Кластеры рендерятся отдельным компонентом
    }

    // Лимитировать количество маркеров
    const maxMarkers = viewport.zoom > 15 ? 500 : 200;
    const limited = visibleMarkers.slice(0, maxMarkers);

    return limited.map(marker => (
      <Marker
        key={marker.id}
        longitude={marker.position[0]}
        latitude={marker.position[1]}
        anchor="bottom"
        onClick={() => onMarkerClick(marker)}
      >
        <div className="map-marker">
          {marker.icon || '📍'}
        </div>
      </Marker>
    ));
  }, [visibleMarkers, viewport.zoom, onMarkerClick]);

  return <>{renderMarkers}</>;
});

VirtualizedMarkers.displayName = 'VirtualizedMarkers';
```

### 1.3 Прогрессивная загрузка данных

```typescript
// frontend/svetu/src/hooks/useProgressiveLoading.ts
export function useProgressiveLoading() {
  const [loadingStage, setLoadingStage] = useState<
    'initial' | 'basic' | 'detailed' | 'complete'
  >('initial');

  const loadProgressively = useCallback(async (
    bounds: MapBounds,
    filters: any
  ) => {
    // Этап 1: Загрузить кластеры/основные точки
    setLoadingStage('basic');
    const clusters = await apiClient.get('/api/v1/gis/clusters', {
      params: { bounds, zoom: getZoom() }
    });

    // Отобразить кластеры немедленно
    displayClusters(clusters.data);

    // Этап 2: Загрузить детали видимой области
    setLoadingStage('detailed');
    const details = await apiClient.get('/api/v1/gis/search/radius', {
      params: {
        ...bounds,
        ...filters,
        limit: 200 // Первая порция
      }
    });

    // Обновить маркеры
    updateMarkers(details.data);

    // Этап 3: Догрузить остальное в фоне
    setLoadingStage('complete');
    if (details.data.hasMore) {
      const remaining = await apiClient.get('/api/v1/gis/search/radius', {
        params: {
          ...bounds,
          ...filters,
          offset: 200,
          limit: 800
        }
      });

      // Добавить к существующим
      appendMarkers(remaining.data);
    }
  }, []);

  return { loadProgressively, loadingStage };
}
```

---

## 🗺️ ФАЗА 2: ИНТЕГРАЦИЯ С БИЗНЕС-ЛОГИКОЙ (1 неделя)

### 2.1 Визуализация Post Express на карте

```typescript
// frontend/svetu/src/components/GIS/Map/layers/PostExpressLayer.tsx
import { useEffect, useState } from 'react';
import { Layer, Source } from 'react-map-gl';

export function PostExpressLayer({
  visible = true,
  selectedListing
}: {
  visible?: boolean;
  selectedListing?: any;
}) {
  const [pickupPoints, setPickupPoints] = useState([]);
  const [deliveryZones, setDeliveryZones] = useState(null);

  // Загрузить точки выдачи Post Express
  useEffect(() => {
    if (!visible) return;

    apiClient.get('/api/v1/post-express/pickup-points')
      .then(res => setPickupPoints(res.data));
  }, [visible]);

  // Рассчитать зоны доставки для выбранного товара
  useEffect(() => {
    if (!selectedListing) {
      setDeliveryZones(null);
      return;
    }

    // Создать зоны доставки
    const zones = {
      instant: createCircle(selectedListing.location, 2000), // 2км - 1 час
      sameDay: createCircle(selectedListing.location, 10000), // 10км - в тот же день
      nextDay: createCircle(selectedListing.location, 50000), // 50км - следующий день
    };

    setDeliveryZones(zones);
  }, [selectedListing]);

  if (!visible) return null;

  return (
    <>
      {/* Точки выдачи Post Express */}
      <Source
        id="post-express-points"
        type="geojson"
        data={{
          type: 'FeatureCollection',
          features: pickupPoints.map(point => ({
            type: 'Feature',
            properties: {
              id: point.id,
              name: point.name,
              type: 'pickup',
              workHours: point.work_hours
            },
            geometry: {
              type: 'Point',
              coordinates: [point.lng, point.lat]
            }
          }))
        }}
      >
        <Layer
          id="post-express-icons"
          type="symbol"
          layout={{
            'icon-image': 'post-express-marker',
            'icon-size': 0.8,
            'text-field': ['get', 'name'],
            'text-offset': [0, 1.5],
            'text-anchor': 'top'
          }}
        />
      </Source>

      {/* Зоны доставки */}
      {deliveryZones && (
        <Source
          id="delivery-zones"
          type="geojson"
          data={{
            type: 'FeatureCollection',
            features: [
              {
                type: 'Feature',
                properties: { zone: 'instant', time: '1 час' },
                geometry: {
                  type: 'Polygon',
                  coordinates: [deliveryZones.instant]
                }
              },
              {
                type: 'Feature',
                properties: { zone: 'sameDay', time: 'Сегодня' },
                geometry: {
                  type: 'Polygon',
                  coordinates: [deliveryZones.sameDay]
                }
              },
              {
                type: 'Feature',
                properties: { zone: 'nextDay', time: 'Завтра' },
                geometry: {
                  type: 'Polygon',
                  coordinates: [deliveryZones.nextDay]
                }
              }
            ]
          }}
        >
          <Layer
            id="delivery-zones-fill"
            type="fill"
            paint={{
              'fill-color': [
                'match',
                ['get', 'zone'],
                'instant', '#10B981',
                'sameDay', '#3B82F6',
                'nextDay', '#6366F1',
                '#000000'
              ],
              'fill-opacity': 0.15
            }}
          />
          <Layer
            id="delivery-zones-line"
            type="line"
            paint={{
              'line-color': [
                'match',
                ['get', 'zone'],
                'instant', '#10B981',
                'sameDay', '#3B82F6',
                'nextDay', '#6366F1',
                '#000000'
              ],
              'line-width': 2,
              'line-dasharray': [2, 2]
            }}
          />
        </Source>
      )}
    </>
  );
}
```

### 2.2 Система маркеров C2C/B2C/Services

```typescript
// frontend/svetu/src/components/GIS/Map/utils/markerStyles.ts
export interface MarkerStyle {
  icon: string;
  color: string;
  size: 'small' | 'medium' | 'large';
  shape: 'circle' | 'square' | 'pin';
  priority: number; // Для z-index
}

export function getMarkerStyle(listing: any): MarkerStyle {
  // B2C - Бизнесы
  if (listing.storefront_id) {
    // Недвижимость
    if (listing.category_slug?.includes('real-estate')) {
      return {
        icon: '🏠',
        color: '#DC2626', // red-600
        size: 'large',
        shape: 'square',
        priority: 10
      };
    }

    // Автомобили
    if (listing.category_slug?.includes('auto')) {
      return {
        icon: '🚗',
        color: '#EA580C', // orange-600
        size: 'large',
        shape: 'circle',
        priority: 9
      };
    }

    // Услуги
    if (listing.category_slug?.includes('service')) {
      return {
        icon: '🔧',
        color: '#CA8A04', // yellow-600
        size: 'medium',
        shape: 'pin',
        priority: 8
      };
    }

    // Магазины
    return {
      icon: '🏪',
      color: '#7C3AED', // violet-600
      size: 'medium',
      shape: 'circle',
      priority: 7
    };
  }

  // C2C - Частные объявления
  if (listing.is_urgent) {
    return {
      icon: '⚡',
      color: '#DC2626', // red для срочных
      size: 'medium',
      shape: 'pin',
      priority: 6
    };
  }

  // Обычные C2C
  return {
    icon: '📍',
    color: '#6B7280', // gray-500
    size: 'small',
    shape: 'circle',
    priority: 1
  };
}

// Компонент маркера
export function StyledMarker({ listing, onClick }: any) {
  const style = getMarkerStyle(listing);

  return (
    <div
      className={`
        map-marker
        ${style.shape}
        ${style.size}
        hover:scale-110
        transition-transform
        cursor-pointer
        relative
      `}
      style={{
        backgroundColor: style.color,
        zIndex: style.priority,
        width: style.size === 'large' ? 40 : style.size === 'medium' ? 32 : 24,
        height: style.size === 'large' ? 40 : style.size === 'medium' ? 32 : 24,
      }}
      onClick={() => onClick(listing)}
    >
      <span className="text-white text-center">
        {style.icon}
      </span>

      {/* Бейдж для витрин */}
      {listing.storefront_id && (
        <div className="absolute -top-2 -right-2 w-4 h-4">
          <img
            src={listing.storefront_logo}
            className="rounded-full border border-white"
            alt=""
          />
        </div>
      )}

      {/* Цена для недвижимости */}
      {listing.category_slug?.includes('real-estate') && (
        <div className="absolute -bottom-6 text-xs font-bold whitespace-nowrap">
          €{listing.price}
        </div>
      )}
    </div>
  );
}
```

### 2.3 Режим витрин и связей

```typescript
// frontend/svetu/src/components/GIS/Map/modes/StorefrontMode.tsx
import { useState, useEffect } from 'react';
import { Source, Layer } from 'react-map-gl';

interface StorefrontModeProps {
  storefrontId?: number;
  showInventory?: boolean;
  showDeliveryRadius?: boolean;
  showConnections?: boolean;
}

export function StorefrontMode({
  storefrontId,
  showInventory = true,
  showDeliveryRadius = false,
  showConnections = false
}: StorefrontModeProps) {
  const [storefront, setStorefront] = useState(null);
  const [products, setProducts] = useState([]);

  useEffect(() => {
    if (!storefrontId) return;

    // Загрузить данные витрины
    Promise.all([
      apiClient.get(`/api/v1/storefronts/${storefrontId}`),
      apiClient.get(`/api/v1/storefronts/${storefrontId}/products`)
    ]).then(([storeRes, productsRes]) => {
      setStorefront(storeRes.data);
      setProducts(productsRes.data);
    });
  }, [storefrontId]);

  if (!storefront) return null;

  // Создать линии связи между витриной и товарами
  const connectionLines = products.map(product => ({
    type: 'Feature',
    properties: {
      product_id: product.id,
      storefront_id: storefront.id
    },
    geometry: {
      type: 'LineString',
      coordinates: [
        [storefront.longitude, storefront.latitude],
        [product.longitude, product.latitude]
      ]
    }
  }));

  return (
    <>
      {/* Главный офис/магазин */}
      <Marker
        longitude={storefront.longitude}
        latitude={storefront.latitude}
      >
        <div className="storefront-hq-marker">
          <img
            src={storefront.logo}
            className="w-12 h-12 rounded-full border-2 border-white shadow-lg"
          />
          <div className="badge badge-primary absolute -bottom-2">
            {products.length} товаров
          </div>
        </div>
      </Marker>

      {/* Товары витрины */}
      {showInventory && products.map(product => (
        <Marker
          key={product.id}
          longitude={product.longitude}
          latitude={product.latitude}
        >
          <StyledMarker
            listing={{
              ...product,
              storefront_logo: storefront.logo
            }}
          />
        </Marker>
      ))}

      {/* Линии связи */}
      {showConnections && (
        <Source
          id="storefront-connections"
          type="geojson"
          data={{
            type: 'FeatureCollection',
            features: connectionLines
          }}
        >
          <Layer
            id="connection-lines"
            type="line"
            paint={{
              'line-color': storefront.brand_color || '#8B5CF6',
              'line-width': 1.5,
              'line-opacity': 0.4,
              'line-dasharray': [3, 3]
            }}
          />
        </Source>
      )}

      {/* Радиус доставки */}
      {showDeliveryRadius && (
        <Source
          id="delivery-radius"
          type="geojson"
          data={{
            type: 'Feature',
            geometry: {
              type: 'Polygon',
              coordinates: [
                createCircle(
                  [storefront.longitude, storefront.latitude],
                  storefront.delivery_radius || 10000
                )
              ]
            }
          }}
        >
          <Layer
            id="delivery-radius-fill"
            type="fill"
            paint={{
              'fill-color': storefront.brand_color || '#8B5CF6',
              'fill-opacity': 0.1
            }}
          />
          <Layer
            id="delivery-radius-line"
            type="line"
            paint={{
              'line-color': storefront.brand_color || '#8B5CF6',
              'line-width': 2,
              'line-opacity': 0.5
            }}
          />
        </Source>
      )}
    </>
  );
}
```

---

## 📈 ФАЗА 3: АНАЛИТИКА И МОНИТОРИНГ (3-4 дня)

### 3.1 Сбор метрик производительности

```typescript
// frontend/svetu/src/utils/monitoring/mapPerformance.ts
class MapPerformanceMonitor {
  private metrics: Map<string, number[]> = new Map();

  startMeasure(name: string) {
    performance.mark(`${name}-start`);
  }

  endMeasure(name: string) {
    performance.mark(`${name}-end`);
    performance.measure(name, `${name}-start`, `${name}-end`);

    const measure = performance.getEntriesByName(name)[0];
    if (measure) {
      this.recordMetric(name, measure.duration);

      // Отправить на сервер если критично
      if (measure.duration > 1000) {
        this.reportSlowOperation(name, measure.duration);
      }
    }
  }

  private recordMetric(name: string, value: number) {
    if (!this.metrics.has(name)) {
      this.metrics.set(name, []);
    }

    const values = this.metrics.get(name)!;
    values.push(value);

    // Хранить только последние 100 измерений
    if (values.length > 100) {
      values.shift();
    }
  }

  getStats(name: string) {
    const values = this.metrics.get(name) || [];
    if (values.length === 0) return null;

    const sorted = [...values].sort((a, b) => a - b);
    return {
      avg: values.reduce((a, b) => a + b) / values.length,
      min: sorted[0],
      max: sorted[sorted.length - 1],
      p50: sorted[Math.floor(sorted.length * 0.5)],
      p95: sorted[Math.floor(sorted.length * 0.95)],
      p99: sorted[Math.floor(sorted.length * 0.99)]
    };
  }

  private async reportSlowOperation(name: string, duration: number) {
    await apiClient.post('/api/v1/gis/analytics/performance', {
      operation: name,
      duration,
      viewport: getCurrentViewport(),
      timestamp: Date.now()
    });
  }
}

export const mapMonitor = new MapPerformanceMonitor();

// Использование в компонентах
export function useMonitoredMapLoad() {
  const loadListings = useCallback(async (params: any) => {
    mapMonitor.startMeasure('map-data-load');

    try {
      const data = await apiClient.get('/api/v1/gis/search', params);
      return data;
    } finally {
      mapMonitor.endMeasure('map-data-load');

      // Логировать статистику
      const stats = mapMonitor.getStats('map-data-load');
      console.log('Map load performance:', stats);
    }
  }, []);

  return { loadListings };
}
```

### 3.2 Аналитика использования

```go
// backend/internal/proj/gis/handler/analytics.go
package handler

type MapAnalyticsEvent struct {
    UserID      int                    `json:"user_id"`
    EventType   string                 `json:"event_type"`
    Viewport    map[string]float64     `json:"viewport"`
    Filters     map[string]interface{} `json:"filters"`
    ResultCount int                    `json:"result_count"`
    Duration    int                    `json:"duration"`
    Timestamp   int64                  `json:"timestamp"`
}

// TrackMapEvent godoc
// @Summary Отслеживание событий карты
// @Tags gis-analytics
// @Accept json
// @Param event body MapAnalyticsEvent true "Событие"
// @Success 200 {object} utils.SuccessResponseSwag
// @Router /api/v1/gis/analytics/track [post]
func (h *Handler) TrackMapEvent(c *fiber.Ctx) error {
    var event MapAnalyticsEvent
    if err := c.BodyParser(&event); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    // Сохранить в gis_filter_analytics
    _, err := h.db.Exec(`
        INSERT INTO gis_filter_analytics (
            user_id,
            event_type,
            viewport_bounds,
            filters_used,
            result_count,
            response_time_ms,
            created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, NOW())
    `,
        event.UserID,
        event.EventType,
        event.Viewport,
        event.Filters,
        event.ResultCount,
        event.Duration,
    )

    if err != nil {
        h.logger.Error("Failed to track event", err)
    }

    // Асинхронно обновить популярные фильтры
    go h.updatePopularFilters(event)

    return c.JSON(fiber.Map{"success": true})
}

// GetMapStats возвращает статистику использования карты
func (h *Handler) GetMapStats(c *fiber.Ctx) error {
    stats := struct {
        TotalSearches   int            `json:"total_searches"`
        AvgResponseTime float64        `json:"avg_response_time"`
        PopularFilters  map[string]int `json:"popular_filters"`
        HeatmapData     []HeatmapPoint `json:"heatmap_data"`
    }{}

    // Получить статистику за последние 7 дней
    h.db.QueryRow(`
        SELECT
            COUNT(*),
            AVG(response_time_ms)
        FROM gis_filter_analytics
        WHERE created_at > NOW() - INTERVAL '7 days'
    `).Scan(&stats.TotalSearches, &stats.AvgResponseTime)

    // Популярные фильтры
    rows, _ := h.db.Query(`
        SELECT
            filters_used,
            COUNT(*) as count
        FROM gis_filter_analytics
        WHERE created_at > NOW() - INTERVAL '7 days'
        GROUP BY filters_used
        ORDER BY count DESC
        LIMIT 10
    `)

    // ... обработка результатов

    return c.JSON(stats)
}
```

---

## 🚧 ФАЗА 4: РАСШИРЕННЫЕ ФУНКЦИИ (2 недели)

### 4.1 Изохроны и доступность

```go
// backend/internal/proj/gis/service/isochrone_service.go
func (s *IsochroneService) GetIsochrone(center Coordinates, minutes int, mode string) (*Isochrone, error) {
    // Проверить кэш
    cached, err := s.checkCache(center, minutes, mode)
    if err == nil && cached != nil {
        return cached, nil
    }

    // Рассчитать новый изохрон
    var polygon string

    switch mode {
    case "walking":
        // 5 км/ч скорость ходьбы
        radius := float64(minutes) * 5.0 / 60.0 * 1000 // в метрах
        polygon = s.createCirclePolygon(center, radius)

    case "driving":
        // Использовать дорожную сеть (упрощенно - 40 км/ч средняя)
        radius := float64(minutes) * 40.0 / 60.0 * 1000
        polygon = s.createCirclePolygon(center, radius)

    case "transit":
        // Использовать точки остановок общественного транспорта
        polygon = s.calculateTransitIsochrone(center, minutes)
    }

    // Сохранить в кэш
    s.saveToCache(center, minutes, mode, polygon)

    return &Isochrone{
        Center:   center,
        Minutes:  minutes,
        Mode:     mode,
        Polygon:  polygon,
        CachedAt: time.Now(),
    }, nil
}
```

### 4.2 Тепловые карты

```typescript
// frontend/svetu/src/components/GIS/Map/layers/HeatmapLayer.tsx
import { HeatmapLayer } from '@deck.gl/aggregation-layers';

export function ListingHeatmap({ listings, visible = true }) {
  if (!visible || !listings.length) return null;

  const heatmapData = listings.map(listing => ({
    coordinates: [listing.longitude, listing.latitude],
    weight: listing.views || 1
  }));

  return (
    <DeckGLOverlay
      layers={[
        new HeatmapLayer({
          id: 'listing-heatmap',
          data: heatmapData,
          getPosition: d => d.coordinates,
          getWeight: d => d.weight,
          radiusPixels: 30,
          intensity: 1,
          threshold: 0.05,
          colorRange: [
            [255, 255, 178, 0],
            [254, 217, 118, 127],
            [254, 178, 76, 200],
            [253, 141, 60, 255],
            [240, 59, 32, 255],
            [189, 0, 38, 255]
          ]
        })
      ]}
    />
  );
}
```

### 4.3 Поиск по полигонам

```typescript
// frontend/svetu/src/components/GIS/Map/tools/PolygonSearch.tsx
import { useState } from 'react';
import { DrawControl } from '@mapbox/mapbox-gl-draw';

export function PolygonSearchTool({ onSearch }) {
  const [isDrawing, setIsDrawing] = useState(false);

  const handleCreate = (e: any) => {
    const polygon = e.features[0];
    const coordinates = polygon.geometry.coordinates[0];

    // Поиск в полигоне
    apiClient.post('/api/v1/gis/search/polygon', {
      polygon: coordinates
    }).then(res => {
      onSearch(res.data);
    });
  };

  return (
    <>
      <button
        className={`btn ${isDrawing ? 'btn-error' : 'btn-primary'}`}
        onClick={() => setIsDrawing(!isDrawing)}
      >
        {isDrawing ? 'Отменить' : 'Поиск в области'}
      </button>

      {isDrawing && (
        <DrawControl
          position="top-left"
          displayControlsDefault={false}
          controls={{
            polygon: true,
            trash: true
          }}
          onCreate={handleCreate}
        />
      )}
    </>
  );
}
```

---

## 📊 МЕТРИКИ УСПЕХА

### Технические метрики
- ⏱️ Время загрузки карты < 2 сек
- 📍 Отображение 1000+ маркеров без лагов
- 💾 Кэш hit rate > 60%
- 🔄 Обновление при фильтрации < 500мс

### Бизнес-метрики
- 👥 Увеличение конверсии на 20%
- 🕐 Среднее время на карте +50%
- 🔍 Использование гео-поиска 40% пользователей
- 📦 Выбор доставки через карту 30% заказов

### UX метрики
- 😊 Удовлетворенность картой > 4.5/5
- 🖱️ Среднее количество взаимодействий > 10
- 📱 Mobile использование > 60%
- 🔄 Повторное использование > 70%

---

## 🗓️ TIMELINE

### Неделя 1 (Критические исправления)
- День 1-2: Фаза 0 - синхронизация данных, лимиты, кэш
- День 3-5: Начало Фазы 1 - кластеризация, виртуализация

### Неделя 2 (Оптимизация)
- День 6-8: Завершение Фазы 1 - прогрессивная загрузка
- День 9-10: Начало Фазы 2 - Post Express интеграция

### Неделя 3 (Бизнес-функции)
- День 11-13: Маркеры C2C/B2C, режим витрин
- День 14-15: Фаза 3 - мониторинг и аналитика

### Неделя 4-5 (Расширенные функции)
- День 16-20: Фаза 4 - изохроны, тепловые карты
- День 21-25: Тестирование, оптимизация, документация

---

## ✅ DEFINITION OF DONE

Каждая фаза считается завершенной когда:

1. ✅ Код написан и протестирован
2. ✅ API документирован в Swagger
3. ✅ Frontend компоненты работают на mobile/desktop
4. ✅ Производительность соответствует метрикам
5. ✅ Миграции применены на staging
6. ✅ Документация обновлена
7. ✅ Код прошел review
8. ✅ E2E тесты пройдены

---

## 🔧 ТЕХНИЧЕСКИЙ СТЕК

### Backend
- PostgreSQL 14+ с PostGIS 3.2+
- Go Fiber для API
- OpenSearch для полнотекстового поиска
- Redis для кэширования

### Frontend
- React 19 + Next.js 15
- Mapbox GL JS / react-map-gl
- Deck.gl для визуализаций
- TailwindCSS + DaisyUI

### Инфраструктура
- Docker для контейнеризации
- Nginx для проксирования
- MinIO для хранения тайлов карт
- Grafana для мониторинга

---

## 📚 ДОКУМЕНТАЦИЯ И РЕСУРСЫ

### Внутренние документы
- `/docs/maps/MAP_IMPLEMENTATION_PLAN.md` - этот план
- `/docs/GIS_API.md` - документация API (создать)
- `/docs/MAP_COMPONENTS.md` - компоненты карты (создать)

### Внешние ресурсы
- [PostGIS Documentation](https://postgis.net/docs/)
- [Mapbox GL JS](https://docs.mapbox.com/mapbox-gl-js/)
- [Deck.gl](https://deck.gl/docs)
- [Turf.js](https://turfjs.org/) - гео-утилиты

---

## 🚀 БЫСТРЫЙ СТАРТ ДЛЯ РАЗРАБОТЧИКА

```bash
# 1. Применить миграции для синхронизации geo данных
cd backend
./migrator up

# 2. Запустить backend с GIS модулем
go run cmd/api/main.go

# 3. Запустить frontend
cd ../frontend/svetu
yarn dev -p 3001

# 4. Открыть карту
http://localhost:3001/map

# 5. Проверить API
http://localhost:3000/swagger/index.html#/gis
```

---

*План создан: 13.09.2025*
*Версия: 1.0.0*
*Автор: DevOps Team*