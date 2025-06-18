# План реализации картографического сервиса

## Фаза 1: Базовая карта (1 неделя)

### Backend задачи:

1. **Расширение моделей для геоданных**
```go
type Location struct {
    // Координаты клика пользователя
    UserLat     float64 `json:"user_lat"`
    UserLng     float64 `json:"user_lng"`
    
    // Координаты здания (после геокодинга)
    BuildingLat float64 `json:"building_lat"`
    BuildingLng float64 `json:"building_lng"`
    
    // Адресные данные
    FullAddress  string `json:"full_address"`
    Street       string `json:"street"`
    HouseNumber  string `json:"house_number"`
    PostalCode   string `json:"postal_code"`
    City         string `json:"city"`
    Country      string `json:"country"`
}
```

2. **API endpoints для карт**
- `GET /api/v1/map/storefronts` - витрины в области
- `GET /api/v1/map/storefronts/:id/products` - товары витрины с геолокацией
- `POST /api/v1/map/geocode` - геокодирование
- `GET /api/v1/map/building/:lat/:lng` - информация о здании
- `GET /api/v1/map/clusters` - кластеризованные данные

3. **Оптимизация запросов**
```sql
-- Индекс для геозапросов
CREATE INDEX idx_storefronts_location ON storefronts USING GIST (
    point(longitude, latitude)
);

-- Функция поиска в радиусе
CREATE OR REPLACE FUNCTION find_nearby_storefronts(
    lat FLOAT, lng FLOAT, radius_km FLOAT
) RETURNS TABLE(...) AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM storefronts
    WHERE earth_distance(
        ll_to_earth(latitude, longitude),
        ll_to_earth(lat, lng)
    ) <= radius_km * 1000;
END;
$$ LANGUAGE plpgsql;
```

### Frontend задачи:

1. **Установка зависимостей**
```bash
yarn add leaflet react-leaflet
yarn add leaflet.markercluster
yarn add @types/leaflet --dev
```

2. **Базовый компонент карты**
```tsx
// components/maps/BaseMap.tsx
interface BaseMapProps {
  center?: [number, number];
  zoom?: number;
  onLocationSelect?: (location: Location) => void;
}

export const BaseMap: React.FC<BaseMapProps> = ({
  center = [44.8125, 20.4612], // Белград
  zoom = 12
}) => {
  return (
    <MapContainer
      center={center}
      zoom={zoom}
      className="w-full h-full"
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; OpenStreetMap contributors'
      />
      <LocationPicker />
      <MarkerClusterGroup />
    </MapContainer>
  );
};
```

## Фаза 2: Продвинутые функции (2 недели)

### 1. Building Intelligence System

```typescript
// services/maps/BuildingService.ts
class BuildingService {
  async getBusinessesInBuilding(lat: number, lng: number) {
    // Находим все бизнесы в радиусе 30м (одно здание)
    const businesses = await api.get('/map/building', {
      params: { lat, lng, radius: 30 }
    });
    
    // Группируем по этажам если есть данные
    return this.groupByFloor(businesses);
  }
  
  async suggestAddress(lat: number, lng: number) {
    // Умное предложение адреса
    const reverse = await this.reverseGeocode(lat, lng);
    const refined = await this.refineAddress(reverse);
    return refined;
  }
}
```

### 2. Кластеризация с кастомными маркерами

```typescript
// Разные иконки для разных типов
const markerIcons = {
  storefront: L.divIcon({
    html: '<div class="storefront-marker">🏪</div>',
    iconSize: [30, 30],
    className: 'custom-marker'
  }),
  
  realEstate: L.divIcon({
    html: '<div class="realestate-marker">🏠</div>',
    iconSize: [30, 30],
    className: 'custom-marker'
  }),
  
  product: L.divIcon({
    html: '<div class="product-marker">📦</div>',
    iconSize: [25, 25],
    className: 'custom-marker'
  })
};

// Кастомная функция создания кластера
const createClusterIcon = (cluster) => {
  const count = cluster.getChildCount();
  const size = count < 10 ? 'small' : count < 100 ? 'medium' : 'large';
  
  return L.divIcon({
    html: `<div class="cluster-${size}">${count}</div>`,
    className: 'marker-cluster',
    iconSize: L.point(40, 40)
  });
};
```

### 3. Зоны доставки

```typescript
// components/maps/DeliveryZones.tsx
interface DeliveryZone {
  id: string;
  name: string;
  polygon: LatLng[];
  priceModifier: number;
  estimatedTime: string;
}

export const DeliveryZoneEditor: React.FC = () => {
  const [zones, setZones] = useState<DeliveryZone[]>([]);
  
  return (
    <FeatureGroup>
      <EditControl
        position="topright"
        onCreated={(e) => handleZoneCreated(e)}
        draw={{
          polygon: true,
          rectangle: false,
          circle: false,
          marker: false,
          polyline: false
        }}
      />
      {zones.map(zone => (
        <Polygon
          key={zone.id}
          positions={zone.polygon}
          color={getZoneColor(zone.priceModifier)}
          fillOpacity={0.3}
        >
          <Popup>{zone.name} - {zone.estimatedTime}</Popup>
        </Polygon>
      ))}
    </FeatureGroup>
  );
};
```

## Фаза 3: Аналитика и оптимизация (1 неделя)

### 1. Heat Maps

```typescript
// components/maps/HeatMap.tsx
import 'leaflet.heat';

export const PopularityHeatMap: React.FC = () => {
  const heatData = useStorefrontHeatData();
  
  useEffect(() => {
    if (map && heatData) {
      L.heatLayer(heatData, {
        radius: 25,
        blur: 15,
        maxZoom: 17,
        gradient: {
          0.4: 'blue',
          0.6: 'cyan',
          0.7: 'lime',
          0.8: 'yellow',
          1.0: 'red'
        }
      }).addTo(map);
    }
  }, [map, heatData]);
};
```

### 2. Офлайн поддержка

```typescript
// services/maps/OfflineService.ts
class OfflineMapService {
  async cacheArea(bounds: LatLngBounds, zoom: number[]) {
    const tiles = this.calculateTiles(bounds, zoom);
    
    for (const tile of tiles) {
      await this.cacheTile(tile);
    }
  }
  
  private async cacheTile(tile: TileCoords) {
    const url = this.getTileUrl(tile);
    const response = await fetch(url);
    const blob = await response.blob();
    
    await this.indexedDB.tiles.put({
      key: `${tile.z}/${tile.x}/${tile.y}`,
      data: blob,
      timestamp: Date.now()
    });
  }
}
```

## Специальные функции для недвижимости

```typescript
// components/maps/RealEstateMap.tsx
interface PropertyMarker {
  id: string;
  type: 'apartment' | 'house' | 'land' | 'commercial';
  price: number;
  size: number;
  location: Location;
  images: string[];
}

export const RealEstateMap: React.FC = () => {
  const [properties, setProperties] = useState<PropertyMarker[]>([]);
  const [filters, setFilters] = useState<PropertyFilters>({});
  
  return (
    <div className="relative h-screen">
      <PropertyFilters onChange={setFilters} />
      <BaseMap>
        <PropertyMarkers 
          properties={properties}
          filters={filters}
          renderPopup={(property) => (
            <PropertyQuickView property={property} />
          )}
        />
      </BaseMap>
      <PropertyList 
        properties={getVisibleProperties(properties, mapBounds)}
        className="absolute right-0 top-0 w-80"
      />
    </div>
  );
};
```

## Метрики успеха

1. **Производительность**
   - Загрузка карты < 2 сек
   - Отрисовка 10k маркеров < 100ms
   - Плавность 60 FPS при перемещении

2. **UX метрики**
   - Точность геокодинга > 95%
   - Конверсия выбора локации > 80%
   - Использование фильтров > 60%

3. **Бизнес метрики**
   - Увеличение конверсии на 30%
   - Снижение времени поиска на 50%
   - Рост повторных визитов на 40%