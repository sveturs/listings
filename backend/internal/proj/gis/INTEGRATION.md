# Интеграция GIS модуля

## Шаги интеграции

### 1. Добавить миграцию для PostGIS

Создайте миграцию для добавления колонки `location` в таблицу `marketplace_listings`:

```sql
-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Add location column
ALTER TABLE marketplace_listings 
ADD COLUMN IF NOT EXISTS location geography(Point, 4326);

-- Create spatial index
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_location 
ON marketplace_listings USING GIST(location);

-- Add address column if not exists
ALTER TABLE marketplace_listings 
ADD COLUMN IF NOT EXISTS address TEXT;
```

### 2. Обновить OpenSearch маппинг

Добавьте геополя в маппинг OpenSearch:

```json
{
  "mappings": {
    "properties": {
      "location": {
        "type": "geo_point"
      },
      "address": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      }
    }
  }
}
```

### 3. Подключить модуль в main.go

В файле `cmd/api/main.go` добавьте:

```go
import (
    gisHandler "backend/internal/proj/gis/handler"
)

// В функции инициализации маршрутов
gisHandler.RegisterRoutes(app, db, middleware)
```

### 4. Добавить переводы

В файлы переводов `frontend/svetu/src/messages/{en,ru}.json` добавьте:

```json
{
  "gis": {
    "invalidBounds": "Invalid map bounds",
    "invalidCenter": "Invalid center coordinates",
    "invalidRadius": "Invalid search radius",
    "boundsRequired": "Map bounds are required",
    "invalidZoomLevel": "Invalid zoom level",
    "coordinatesRequired": "Coordinates are required",
    "invalidListingId": "Invalid listing ID",
    "listingNotFound": "Listing not found",
    "searchError": "Error searching listings",
    "clusterError": "Error getting clusters",
    "getLocationError": "Error getting location",
    "updateLocationError": "Error updating location",
    "invalidMinPrice": "Invalid minimum price",
    "invalidMaxPrice": "Invalid maximum price",
    "invalidRequest": "Invalid request",
    "validationError": "Validation error"
  }
}
```

### 5. Обновить модель Listing

В существующей модели `marketplace_listings` добавьте поддержку геоданных при создании/обновлении объявлений.

### 6. Обновить Swagger документацию

После добавления маршрутов выполните:

```bash
make generate-types
```

## Использование API

### Поиск объявлений на карте
```bash
curl "http://localhost:3000/api/v1/gis/search?bounds=45.5,44.5,20.5,19.5&categories=electronics"
```

### Получение кластеров
```bash
curl "http://localhost:3000/api/v1/gis/clusters?bounds=45.5,44.5,20.5,19.5&zoom_level=10"
```

### Поиск в радиусе
```bash
curl "http://localhost:3000/api/v1/gis/nearby?lat=45.0&lng=20.0&radius_km=5"
```

### Обновление локации (требует авторизацию)
```bash
curl -X PUT "http://localhost:3000/api/v1/gis/listings/{id}/location" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lat": 45.0,
    "lng": 20.0,
    "address": "Main Street 123, Belgrade"
  }'
```

## Frontend интеграция

Для отображения карты на frontend рекомендуется использовать:
- Leaflet или Mapbox GL JS для отображения карты
- React-Leaflet для интеграции с React
- Marker Clustering плагин для кластеризации

Пример компонента карты будет предоставлен в следующей фазе разработки.