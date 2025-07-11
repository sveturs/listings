# 🗺️ Advanced GIS: Детальный план реализации

**Версия**: 1.0  
**Дата**: 2025-01-10  
**Длительность**: 12 недель  
**Команда**: 3-4 разработчика

## 📊 Обзор проекта

### Ресурсы:
- **Backend разработчик** (Senior) - 100% 
- **Frontend разработчик** (Senior) - 100%
- **Full-stack разработчик** (Middle+) - 100%
- **DevOps** (частичная занятость) - 30%
- **UI/UX дизайнер** (первые 4 недели) - 50%

### Бюджет технологий:
- Mapbox: $500/месяц (50k загрузок карт)
- TomTom Traffic API: $200/месяц
- Cloudflare: $20/месяц
- Дополнительные серверы: ~$300/месяц

## 📅 ФАЗА 1: Фундамент (3 недели)

### Неделя 1: Инфраструктура и база данных

#### День 1-2: Настройка окружения
```bash
# Backend задачи
- [ ] Установка PostGIS 3.4 на dev/staging
- [ ] Создание Docker образов с PostGIS
- [ ] Настройка репликации для гео-данных
- [ ] Настройка бэкапов с геоданными

# DevOps задачи  
- [ ] Обновление CI/CD для PostGIS миграций
- [ ] Настройка мониторинга PostGIS метрик
- [ ] Конфигурация Cloudflare для тайлов
```

#### День 3-5: Миграция данных
```sql
-- Скрипт миграции существующих координат
BEGIN;

-- Включаем PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_raster;
CREATE EXTENSION IF NOT EXISTS pg_trgm; -- для fuzzy поиска

-- Новая структура для геоданных
CREATE TABLE listings_geo_new (
    id UUID PRIMARY KEY,
    listing_id UUID NOT NULL REFERENCES marketplace_listings(id),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    geohash4 VARCHAR(4) GENERATED ALWAYS AS (ST_GeoHash(location, 4)) STORED,
    geohash6 VARCHAR(6) GENERATED ALWAYS AS (ST_GeoHash(location, 6)) STORED,
    geohash8 VARCHAR(8) GENERATED ALWAYS AS (ST_GeoHash(location, 8)) STORED,
    
    -- Денормализованные данные для скорости
    city VARCHAR(100),
    district VARCHAR(100),
    street VARCHAR(200),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_listings_geo_location ON listings_geo_new USING GIST(location);
CREATE INDEX idx_listings_geo_geohash4 ON listings_geo_new(geohash4);
CREATE INDEX idx_listings_geo_geohash6 ON listings_geo_new(geohash6);
CREATE INDEX idx_listings_geo_city ON listings_geo_new(city);

-- Миграция данных
INSERT INTO listings_geo_new (listing_id, location, city)
SELECT 
    id,
    ST_Point(longitude, latitude)::geography,
    address_city
FROM marketplace_listings
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

COMMIT;
```

### Неделя 2: Базовый API и сервисы

#### Backend структура:
```go
// internal/proj/gis/service/spatial_service.go
package service

type SpatialService struct {
    db        *sql.DB
    cache     *redis.Client
    elastic   *elasticsearch.Client
}

// Базовые методы
func (s *SpatialService) SearchRadius(ctx context.Context, center Point, radiusKm float64, filters Filters) ([]Listing, error)
func (s *SpatialService) SearchBoundingBox(ctx context.Context, bounds BBox, filters Filters) ([]Listing, error)
func (s *SpatialService) GetClusters(ctx context.Context, bounds BBox, zoom int) ([]Cluster, error)
func (s *SpatialService) GetNearestListings(ctx context.Context, point Point, limit int) ([]Listing, error)
```

#### API Endpoints (OpenAPI):
```yaml
/api/v1/gis/search/radius:
  post:
    summary: Поиск в радиусе
    parameters:
      - name: lat
        type: number
        required: true
      - name: lng  
        type: number
        required: true
      - name: radius_km
        type: number
        default: 5
      - name: category_id
        type: integer
      - name: price_min
        type: number
      - name: price_max
        type: number
    responses:
      200:
        schema:
          type: array
          items:
            $ref: '#/definitions/GeoListing'

/api/v1/gis/clusters:
  get:
    summary: Получить кластеры для карты
    parameters:
      - name: bounds
        description: "sw_lat,sw_lng,ne_lat,ne_lng"
      - name: zoom
        type: integer
        minimum: 1
        maximum: 20
```

### Неделя 3: Базовая карта

#### Frontend компоненты:
```typescript
// src/components/GIS/Map/MapContainer.tsx
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

export const MapContainer: React.FC = () => {
  const [map, setMap] = useState<mapboxgl.Map | null>(null);
  
  useEffect(() => {
    const map = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [20.4568, 44.8178], // Белград
      zoom: 12
    });
    
    // Добавляем контролы
    map.addControl(new mapboxgl.NavigationControl());
    map.addControl(new mapboxgl.GeolocateControl());
    
    setMap(map);
  }, []);
  
  return <div ref={mapContainer} className="h-full w-full" />;
};
```

#### Интеграция с API:
```typescript
// src/hooks/useGeoSearch.ts
export const useGeoSearch = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (params: GeoSearchParams) => {
      const response = await api.post('/gis/search/radius', params);
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(['geoListings'], data);
    }
  });
};
```

## 📅 ФАЗА 2: Основной функционал (4 недели)

### Неделя 4-5: Умный поиск и фильтры

#### Natural Language Processing для поиска:
```python
# scripts/train_search_model.py
import spacy
from transformers import pipeline

# Обучаем модель понимать запросы типа:
# "детские товары до 2000 динар рядом с Калемегданом"
# "кафе с wifi в центре"

class GeoSearchNLP:
    def __init__(self):
        self.nlp = spacy.load("sr_core_news_sm")  # Сербская модель
        self.ner = pipeline("ner", model="xlm-roberta-base")
        
    def parse_query(self, query: str) -> dict:
        doc = self.nlp(query)
        
        # Извлекаем сущности
        entities = {
            "categories": [],
            "price_max": None,
            "location": None,
            "amenities": []
        }
        
        for ent in doc.ents:
            if ent.label_ == "LOC":
                entities["location"] = ent.text
            elif ent.label_ == "MONEY":
                entities["price_max"] = self.parse_price(ent.text)
                
        return entities
```

#### Фильтры в реальном времени:
```typescript
// src/components/GIS/Filters/RealtimeFilters.tsx
export const RealtimeFilters = ({ onFiltersChange }) => {
  const [filters, setFilters] = useState<Filters>({
    categories: [],
    priceRange: [0, 10000],
    radius: 5,
    openNow: false,
    hasDelivery: false,
    rating: 0
  });
  
  // Debounced обновление
  const debouncedChange = useMemo(
    () => debounce(onFiltersChange, 300),
    [onFiltersChange]
  );
  
  useEffect(() => {
    debouncedChange(filters);
  }, [filters, debouncedChange]);
  
  return (
    <div className="space-y-4 p-4">
      {/* Категории */}
      <CategoryFilter 
        selected={filters.categories}
        onChange={(cats) => setFilters({...filters, categories: cats})}
      />
      
      {/* Ценовой диапазон */}
      <PriceRangeSlider
        value={filters.priceRange}
        onChange={(range) => setFilters({...filters, priceRange: range})}
      />
      
      {/* Радиус поиска */}
      <RadiusSelector
        value={filters.radius}
        onChange={(radius) => setFilters({...filters, radius})}
      />
    </div>
  );
};
```

### Неделя 6-7: Кластеризация и оптимизация

#### Серверная кластеризация:
```go
// internal/proj/gis/service/clustering.go
func (s *SpatialService) CreateClusters(listings []Listing, zoom int) []Cluster {
    // Размер ячейки в зависимости от зума
    cellSize := getCellSize(zoom)
    
    // Группировка по геохешам
    clusters := make(map[string]*Cluster)
    
    for _, listing := range listings {
        // Вычисляем геохеш для текущего зума
        hash := geohash.Encode(listing.Lat, listing.Lng, getPrecision(zoom))
        
        if cluster, exists := clusters[hash]; exists {
            cluster.Count++
            cluster.Listings = append(cluster.Listings, listing.ID)
            // Пересчитываем центр
            cluster.RecalculateCenter()
        } else {
            clusters[hash] = &Cluster{
                ID:       hash,
                Center:   Point{Lat: listing.Lat, Lng: listing.Lng},
                Count:    1,
                Listings: []string{listing.ID},
            }
        }
    }
    
    return clustersToSlice(clusters)
}
```

#### Клиентская оптимизация:
```typescript
// src/utils/mapOptimization.ts
export class MapOptimizer {
    private renderQueue: Marker[] = [];
    private visibleBounds: LngLatBounds;
    private renderBatchSize = 50;
    
    scheduleRender(markers: Marker[]) {
        // Сортируем по приоритету (ближе к центру - выше приоритет)
        const sorted = this.prioritizeMarkers(markers);
        
        // Рендерим батчами
        this.renderBatch(sorted.slice(0, this.renderBatchSize));
        
        // Остальные в очередь
        this.renderQueue = sorted.slice(this.renderBatchSize);
        this.scheduleNextBatch();
    }
    
    private scheduleNextBatch() {
        requestIdleCallback(() => {
            if (this.renderQueue.length > 0) {
                const batch = this.renderQueue.splice(0, this.renderBatchSize);
                this.renderBatch(batch);
                this.scheduleNextBatch();
            }
        });
    }
}
```

## 📅 ФАЗА 3: Продвинутые функции (3 недели)

### Неделя 8: Зоны доставки

#### Backend модель:
```sql
-- Таблица зон доставки
CREATE TABLE delivery_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    storefront_id UUID NOT NULL REFERENCES user_storefronts(id),
    zone_name VARCHAR(100),
    zone_polygon GEOGRAPHY(POLYGON, 4326) NOT NULL,
    delivery_fee DECIMAL(10,2) DEFAULT 0,
    min_order_amount DECIMAL(10,2) DEFAULT 0,
    max_delivery_time_minutes INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Функция проверки попадания точки в зону
CREATE FUNCTION check_delivery_availability(
    p_storefront_id UUID,
    p_lat DECIMAL,
    p_lng DECIMAL
) RETURNS TABLE (
    zone_id UUID,
    zone_name VARCHAR,
    delivery_fee DECIMAL,
    min_order_amount DECIMAL,
    estimated_time INT
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
        AND ST_Contains(
            dz.zone_polygon::geometry,
            ST_Point(p_lng, p_lat)::geometry
        );
END;
$$ LANGUAGE plpgsql;
```

#### Frontend редактор зон:
```typescript
// src/components/GIS/DeliveryZoneEditor.tsx
import { useMapDraw } from '@/hooks/useMapDraw';

export const DeliveryZoneEditor = ({ storefrontId, onSave }) => {
    const { drawPolygon, editPolygon, deletePolygon } = useMapDraw();
    const [zones, setZones] = useState<DeliveryZone[]>([]);
    
    const handleDrawComplete = (polygon: Polygon) => {
        const newZone: DeliveryZone = {
            id: generateId(),
            polygon,
            deliveryFee: 0,
            minOrderAmount: 0,
            color: generateColor()
        };
        
        setZones([...zones, newZone]);
    };
    
    return (
        <div className="relative h-full">
            <Map
                onDrawComplete={handleDrawComplete}
                drawingMode="polygon"
            >
                {zones.map(zone => (
                    <PolygonLayer
                        key={zone.id}
                        polygon={zone.polygon}
                        color={zone.color}
                        opacity={0.3}
                        onClick={() => editPolygon(zone.id)}
                    />
                ))}
            </Map>
            
            <ZoneSettings
                zones={zones}
                onUpdate={(id, settings) => updateZoneSettings(id, settings)}
                onDelete={(id) => deleteZone(id)}
            />
        </div>
    );
};
```

### Неделя 9: Маршрутизация

#### Интеграция с pgRouting:
```sql
-- Импорт дорожной сети Сербии из OpenStreetMap
-- (выполняется отдельным скриптом osm2pgrouting)

-- Функция построения маршрута
CREATE FUNCTION calculate_route(
    start_lat DECIMAL,
    start_lng DECIMAL,
    end_lat DECIMAL,
    end_lng DECIMAL
) RETURNS TABLE (
    seq INT,
    edge BIGINT,
    cost DOUBLE PRECISION,
    geom GEOMETRY
) AS $$
DECLARE
    start_vertex BIGINT;
    end_vertex BIGINT;
BEGIN
    -- Находим ближайшие вершины графа
    SELECT id INTO start_vertex
    FROM roads_vertices_pgr
    ORDER BY the_geom <-> ST_Point(start_lng, start_lat)::geometry
    LIMIT 1;
    
    SELECT id INTO end_vertex
    FROM roads_vertices_pgr
    ORDER BY the_geom <-> ST_Point(end_lng, end_lat)::geometry
    LIMIT 1;
    
    -- Строим маршрут
    RETURN QUERY
    SELECT 
        d.seq,
        d.edge,
        d.cost,
        r.geom
    FROM pgr_dijkstra(
        'SELECT id, source, target, cost FROM roads',
        start_vertex,
        end_vertex,
        false
    ) d
    JOIN roads r ON d.edge = r.id;
END;
$$ LANGUAGE plpgsql;
```

#### Frontend отображение маршрута:
```typescript
// src/components/GIS/RouteDisplay.tsx
export const RouteDisplay = ({ start, end }) => {
    const { data: route, isLoading } = useQuery({
        queryKey: ['route', start, end],
        queryFn: () => api.getRoute(start, end),
        staleTime: 5 * 60 * 1000 // 5 минут
    });
    
    if (!route) return null;
    
    return (
        <>
            {/* Линия маршрута */}
            <Source
                type="geojson"
                data={{
                    type: 'Feature',
                    geometry: {
                        type: 'LineString',
                        coordinates: route.coordinates
                    }
                }}
            >
                <Layer
                    type="line"
                    paint={{
                        'line-color': '#3b82f6',
                        'line-width': 4,
                        'line-opacity': 0.8
                    }}
                />
            </Source>
            
            {/* Информация о маршруте */}
            <RouteInfo
                distance={route.distance}
                duration={route.duration}
                steps={route.steps}
            />
        </>
    );
};
```

### Неделя 10: Real-time и аналитика

#### WebSocket сервер для отслеживания:
```go
// internal/proj/gis/realtime/tracker.go
type LocationUpdate struct {
    OrderID  string  `json:"order_id"`
    Lat      float64 `json:"lat"`
    Lng      float64 `json:"lng"`
    Speed    float64 `json:"speed"`
    Heading  float64 `json:"heading"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (h *Hub) HandleDeliveryTracking(conn *websocket.Conn, orderID string) {
    // Подписываемся на обновления
    sub := h.Subscribe(fmt.Sprintf("delivery:%s", orderID))
    defer h.Unsubscribe(sub)
    
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case update := <-sub:
            // Отправляем обновление клиенту
            if err := conn.WriteJSON(update); err != nil {
                return
            }
            
        case <-ticker.C:
            // Ping для поддержания соединения
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

#### Дашборд геоаналитики:
```typescript
// src/components/Analytics/GeoDashboard.tsx
export const GeoDashboard = () => {
    const { data: heatmapData } = useGeoAnalytics('heatmap');
    const { data: districtStats } = useGeoAnalytics('districts');
    
    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Тепловая карта продаж */}
            <Card>
                <CardHeader>
                    <CardTitle>Тепловая карта продаж</CardTitle>
                </CardHeader>
                <CardContent>
                    <HeatmapLayer
                        data={heatmapData}
                        gradient={{
                            0.0: 'blue',
                            0.5: 'yellow',
                            1.0: 'red'
                        }}
                    />
                </CardContent>
            </Card>
            
            {/* Статистика по районам */}
            <Card>
                <CardHeader>
                    <CardTitle>Топ районов</CardTitle>
                </CardHeader>
                <CardContent>
                    <DistrictChart data={districtStats} />
                </CardContent>
            </Card>
            
            {/* Оптимальная зона */}
            <Card>
                <CardHeader>
                    <CardTitle>Рекомендуемая зона доставки</CardTitle>
                </CardHeader>
                <CardContent>
                    <OptimalZoneMap />
                </CardContent>
            </Card>
        </div>
    );
};
```

## 📅 ФАЗА 4: Полировка и запуск (2 недели)

### Неделя 11: Оптимизация производительности

#### Чек-лист оптимизаций:
- [ ] Включить HTTP/2 Server Push для тайлов
- [ ] Настроить Service Worker для offline карт
- [ ] Оптимизировать bundle size (tree shaking)
- [ ] Включить Brotli сжатие для GeoJSON
- [ ] Настроить индексы в Elasticsearch
- [ ] Профилирование SQL запросов
- [ ] Оптимизация React рендеринга

#### Метрики производительности:
```typescript
// src/utils/performanceMonitoring.ts
export const trackMapPerformance = () => {
    // Core Web Vitals для карты
    const observer = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
            analytics.track('map_performance', {
                metric: entry.name,
                value: entry.value,
                url: window.location.pathname
            });
        }
    });
    
    observer.observe({ entryTypes: ['largest-contentful-paint', 'first-input-delay'] });
};
```

### Неделя 12: Тестирование и запуск

#### План тестирования:
1. **Unit тесты** (покрытие > 80%)
   - Геометрические функции
   - Сервисы поиска
   - React компоненты

2. **Интеграционные тесты**
   - API endpoints
   - WebSocket соединения
   - Кэширование

3. **E2E тесты** (Cypress)
   - Поиск на карте
   - Фильтрация
   - Создание зон доставки

4. **Нагрузочное тестирование**
   - 1000 одновременных пользователей
   - 10k маркеров на карте
   - 100 запросов/сек на геопоиск

#### Чек-лист запуска:
- [ ] Миграция production данных
- [ ] Настройка CDN и кэширования
- [ ] Мониторинг и алерты
- [ ] Документация API
- [ ] Обучающие материалы
- [ ] A/B тест (10% трафика)
- [ ] Постепенный rollout

## 📊 Метрики успеха (KPI)

### Технические метрики:
| Метрика | Неделя 1 | Неделя 4 | Неделя 12 |
|---------|----------|----------|-----------|
| Uptime | 99% | 99.5% | 99.9% |
| P95 latency | 500ms | 300ms | 200ms |
| Ошибки | < 1% | < 0.5% | < 0.1% |

### Бизнес метрики:
| Метрика | До GIS | После GIS | Рост |
|---------|--------|-----------|------|
| Конверсия | 3.2% | 4.5% | +40% |
| Ср. чек | 2500 RSD | 3000 RSD | +20% |
| Retention | 25% | 35% | +40% |

## 🚀 Риски и митигация

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Задержка с PostGIS | Низкая | Высокое | Подготовить fallback на Haversine |
| Проблемы с Mapbox | Средняя | Среднее | Альтернатива - Leaflet + OSM |
| Низкая производительность | Средняя | Высокое | Профилирование с первого дня |
| Сложность для пользователей | Низкая | Среднее | UX тестирование, обучение |

## ✅ Критерии успеха проекта

1. **Технические:**
   - Все тесты проходят (unit, integration, e2e)
   - Производительность соответствует SLA
   - Нет критических багов

2. **Продуктовые:**
   - 40% пользователей используют карту
   - Конверсия из карты > 8%
   - NPS вырос на 10 пунктов

3. **Бизнес:**
   - ROI положительный через 6 месяцев
   - Рост локальных продаж +30%
   - Снижение support tickets на геолокацию

---

**Advanced GIS - от идеи до production за 12 недель!** 🚀