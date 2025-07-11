# План реализации GIS модуля - Этап 1: Основная инфраструктура

## Обзор первого этапа

**Цель:** Создать базовую инфраструктуру для работы с геопространственными данными в маркетплейсе Sve Tu.

**Длительность:** 2 недели

**Приоритет:** Высокий

## Детальный план задач

### Неделя 1: Backend инфраструктура

#### День 1-2: Настройка PostGIS

**Задачи:**
1. Установка PostGIS расширения
   ```sql
   CREATE EXTENSION IF NOT EXISTS postgis;
   CREATE EXTENSION IF NOT EXISTS postgis_topology;
   ```

2. Создание миграции для пространственных индексов
   ```sql
   -- backend/migrations/000XXX_add_spatial_indexes.up.sql
   -- Создание пространственных колонок
   ALTER TABLE marketplace_listings 
   ADD COLUMN IF NOT EXISTS location GEOGRAPHY(POINT, 4326);
   
   -- Заполнение из существующих координат
   UPDATE marketplace_listings 
   SET location = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
   WHERE longitude IS NOT NULL AND latitude IS NOT NULL;
   
   -- Индексы
   CREATE INDEX idx_listings_location ON marketplace_listings USING GIST(location);
   CREATE INDEX idx_listings_city ON marketplace_listings(address_city);
   ```

3. Создание таблиц для административных единиц
   ```sql
   -- backend/migrations/000XXX_create_administrative_boundaries.up.sql
   CREATE TABLE administrative_boundaries (
       id SERIAL PRIMARY KEY,
       code VARCHAR(20) UNIQUE NOT NULL,
       name VARCHAR(100) NOT NULL,
       name_cyrillic VARCHAR(100),
       name_latin VARCHAR(100),
       type VARCHAR(50) NOT NULL, -- 'municipality', 'district', 'region', 'country'
       parent_id INT REFERENCES administrative_boundaries(id),
       population INT,
       area_km2 DECIMAL(10,2),
       boundary GEOGRAPHY(MULTIPOLYGON, 4326),
       center_point GEOGRAPHY(POINT, 4326),
       metadata JSONB DEFAULT '{}',
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   
   CREATE INDEX idx_boundaries_type ON administrative_boundaries(type);
   CREATE INDEX idx_boundaries_parent ON administrative_boundaries(parent_id);
   CREATE INDEX idx_boundaries_boundary ON administrative_boundaries USING GIST(boundary);
   ```

**Результат:** База данных готова для работы с геоданными

#### День 3-4: Базовая структура GIS модуля

**Создание файловой структуры:**
```
backend/internal/proj/gis/
├── module.go                    # Инициализация модуля
├── handler/
│   └── handler.go               # Базовые HTTP handlers
├── service/
│   ├── interface.go             # Интерфейсы сервисов
│   └── gis_service.go          # Основной GIS сервис
├── repository/
│   ├── interface.go             # Интерфейсы репозиториев
│   └── postgres/
│       └── spatial_repository.go # Пространственные запросы
└── types/
    └── types.go                 # Базовые типы данных
```

**Базовые типы данных:**
```go
// backend/internal/proj/gis/types/types.go
package types

type Point struct {
    Lat float64 `json:"lat" validate:"required,min=-90,max=90"`
    Lng float64 `json:"lng" validate:"required,min=-180,max=180"`
}

type Bounds struct {
    NorthEast Point `json:"northeast"`
    SouthWest Point `json:"southwest"`
}

type SpatialSearchRequest struct {
    Center   *Point   `json:"center"`
    Radius   float64  `json:"radius"`      // в километрах
    Bounds   *Bounds  `json:"bounds"`
    Category *int     `json:"category_id"`
    Limit    int      `json:"limit"`
    Offset   int      `json:"offset"`
}

type GeocodingRequest struct {
    Address string `json:"address" validate:"required"`
    Country string `json:"country" default:"RS"`
}

type GeocodingResponse struct {
    Formatted   string  `json:"formatted"`
    Lat         float64 `json:"lat"`
    Lng         float64 `json:"lng"`
    Confidence  float64 `json:"confidence"`
    Components  map[string]string `json:"components"`
}
```

#### День 5: Интеграция с существующим кодом

**Задачи:**
1. Добавление GIS модуля в глобальный сервис
2. Регистрация маршрутов API
3. Обновление Swagger документации

**API endpoints:**
```
GET  /api/v1/gis/search         # Пространственный поиск
POST /api/v1/gis/geocode        # Геокодирование адреса
GET  /api/v1/gis/reverse        # Обратное геокодирование
GET  /api/v1/gis/boundaries     # Получение административных границ
```

### Неделя 2: Frontend и интеграция

#### День 6-7: Базовые React компоненты

**Установка зависимостей:**
```json
{
  "dependencies": {
    "leaflet": "^1.9.4",
    "react-leaflet": "^4.2.1",
    "@types/leaflet": "^1.9.8",
    "leaflet.markercluster": "^1.5.3"
  }
}
```

**Базовый компонент карты:**
```typescript
// frontend/svetu/src/components/GIS/Map/BaseMap.tsx
import { MapContainer, TileLayer } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';

interface BaseMapProps {
  center?: [number, number];
  zoom?: number;
  height?: string;
  children?: React.ReactNode;
}

export default function BaseMap({ 
  center = [44.8178, 20.4568], // Белград по умолчанию
  zoom = 12,
  height = '400px',
  children 
}: BaseMapProps) {
  return (
    <MapContainer
      center={center}
      zoom={zoom}
      style={{ height, width: '100%' }}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {children}
    </MapContainer>
  );
}
```

#### День 8-9: Хуки и утилиты

**Создание хуков:**
```typescript
// frontend/svetu/src/hooks/useGeolocation.ts
export function useGeolocation() {
  const [location, setLocation] = useState<GeolocationPosition | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  
  // Реализация получения геолокации
}

// frontend/svetu/src/hooks/useGeocode.ts  
export function useGeocode() {
  const [loading, setLoading] = useState(false);
  
  const geocode = async (address: string) => {
    // Вызов API геокодирования
  };
  
  return { geocode, loading };
}
```

#### День 10: Тестирование и документация

**Задачи:**
1. Unit тесты для backend сервисов
2. Интеграционные тесты API
3. Компонентные тесты React
4. Обновление документации API

## Ключевые решения

### 1. Выбор картографической библиотеки

**Leaflet vs Google Maps:**
- ✅ Leaflet: открытый, бесплатный, гибкий
- ❌ Google Maps: платный после лимита, ограничения
- **Решение:** Leaflet с OpenStreetMap

### 2. Геокодирование

**Варианты:**
- Nominatim (OSM) - бесплатный, но ограничения по запросам
- Google Geocoding API - платный, точный
- Mapbox - средний вариант
- **Решение:** Nominatim с кэшированием результатов

### 3. Хранение геоданных

**PostGIS vs отдельные колонки:**
- ✅ PostGIS: мощные пространственные функции
- ✅ Поддержка сложных запросов
- ✅ Стандартизированный подход
- **Решение:** PostGIS с обратной совместимостью

### 4. Административные границы Сербии

**Источники данных:**
- OpenStreetMap export для Сербии
- Официальные данные Республиканского геодетического завода
- **Решение:** Импорт из OSM с возможностью обновления

## Риски и митигация

### Технические риски

1. **Производительность пространственных запросов**
   - Митигация: Правильные индексы, денормализация критичных данных

2. **Точность геокодирования сербских адресов**
   - Митигация: Кэширование, ручная корректировка, альтернативные источники

3. **Нагрузка на сервер от карт**
   - Митигация: CDN для тайлов, ленивая загрузка, кластеризация

### Организационные риски

1. **Недостаток экспертизы в GIS**
   - Митигация: Консультации, обучение, простые решения вначале

2. **Изменение требований**
   - Митигация: Модульная архитектура, итеративная разработка

## Критерии готовности первого этапа

### Backend
- [ ] PostGIS установлен и настроен
- [ ] Миграции выполнены успешно
- [ ] GIS модуль создан и интегрирован
- [ ] API endpoints работают
- [ ] Тесты проходят

### Frontend  
- [ ] Базовый компонент карты отображается
- [ ] Leaflet корректно инициализирован
- [ ] Хуки созданы и протестированы
- [ ] Интеграция с API работает

### Документация
- [ ] API документация обновлена
- [ ] README для GIS модуля создан
- [ ] Примеры использования добавлены

## Следующие шаги

После успешного завершения первого этапа:
1. Демонстрация базовой функциональности
2. Сбор обратной связи
3. Корректировка плана второго этапа
4. Начало работы над интерактивной картой объявлений

## Ресурсы и ссылки

- [PostGIS документация](https://postgis.net/docs/)
- [Leaflet документация](https://leafletjs.com/reference.html)
- [React Leaflet](https://react-leaflet.js.org/)
- [OpenStreetMap данные Сербии](https://download.geofabrik.de/europe/serbia.html)
- [Nominatim API](https://nominatim.org/release-docs/latest/api/Overview/)

---

**Статус:** Готов к обсуждению
**Автор:** GIS команда
**Дата:** 2025-01-10