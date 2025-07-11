# GIS Module

Модуль для работы с географическими данными и пространственным поиском объявлений.

## Структура модуля

```
gis/
├── types/              # Типы данных
│   ├── geo_types.go   # Point, Bounds, GeoListing, Cluster
│   └── errors.go      # Специфичные ошибки модуля
├── repository/         # Слой работы с БД
│   └── postgis_repo.go # PostGIS репозиторий
├── service/           # Бизнес-логика
│   └── spatial_service.go # Сервис пространственных операций
└── handler/           # HTTP обработчики
    ├── spatial_search.go # API endpoints
    └── routes.go      # Регистрация маршрутов
```

## API Endpoints

### Публичные (без авторизации)

- `GET /api/v1/gis/search` - Пространственный поиск объявлений
- `GET /api/v1/gis/clusters` - Получение кластеров для карты
- `GET /api/v1/gis/nearby` - Ближайшие объявления от точки
- `GET /api/v1/gis/listings/:id/location` - Геоданные объявления

### Защищенные (требуют авторизацию)

- `PUT /api/v1/gis/listings/:id/location` - Обновление геолокации объявления

## Основные типы

### Point
```go
type Point struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}
```

### Bounds
```go
type Bounds struct {
    North float64 `json:"north"`
    South float64 `json:"south"`
    East  float64 `json:"east"`
    West  float64 `json:"west"`
}
```

### SearchParams
Параметры для пространственного поиска с фильтрацией по:
- Географическим границам или радиусу от точки
- Категориям
- Ценовому диапазону
- Текстовому запросу

### ClusterParams
Параметры для кластеризации объявлений на карте с учетом уровня зума.

## Примеры использования

### Поиск в границах карты
```
GET /api/v1/gis/search?bounds=45.5,44.5,20.5,19.5&categories=electronics,clothing
```

### Поиск в радиусе
```
GET /api/v1/gis/search?center=45.0,20.0&radius_km=5&sort_by=distance
```

### Получение кластеров
```
GET /api/v1/gis/clusters?bounds=45.5,44.5,20.5,19.5&zoom_level=10
```

### Ближайшие объявления
```
GET /api/v1/gis/nearby?lat=45.0&lng=20.0&radius_km=2&limit=10
```

## Интеграция

Для подключения модуля в основное приложение:

```go
import "backend/internal/proj/gis/handler"

// В функции инициализации маршрутов
handler.RegisterRoutes(app, db, authMiddleware)
```

## Требования к БД

Модуль требует PostGIS расширение в PostgreSQL:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

Индекс для пространственных запросов создается автоматически при добавлении колонки типа geography/geometry.