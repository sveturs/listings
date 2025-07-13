# Отчет: Исправление фильтрации по радиусу в GIS модуле

## Проблема

Пользователь обнаружил, что фильтрация по радиусу на карте не работает - параметр `distance` не применялся при кластеризации объявлений, хотя работал при обычном поиске.

## Проведенный анализ

### 1. Анализ frontend кода (map/page.tsx)
- ✅ **Правильная передача параметров**: Frontend корректно передает `distance` параметр в API запросы
- ✅ **Корректная конвертация**: Радиус правильно конвертируется из метров в километры (`${radiusKm}km`)
- ✅ **Логирование**: Параметры корректно логируются в браузере

### 2. Анализ backend кода

#### 2.1 Handler (spatial_search.go)
- ✅ **SearchListings**: Корректно обрабатывает параметр `distance`
- ❌ **GetClusters**: **НЕ ОБРАБАТЫВАЛ** параметры `distance`, `latitude`, `longitude`

#### 2.2 Repository (postgis_repo.go)
- ✅ **SearchListings**: Корректно применяет SQL фильтр `ST_DWithin`
- ❌ **GetClusters**: **НЕ ПРИМЕНЯЛ** радиусный фильтр в SQL запросе

#### 2.3 Types (geo_types.go)
- ❌ **ClusterParams**: **НЕ СОДЕРЖАЛ** поля `Center` и `RadiusKm`

## Внесенные исправления

### 1. Обновление типов данных
**Файл:** `internal/proj/gis/types/geo_types.go`
```go
// Добавлены поля для радиусной фильтрации
type ClusterParams struct {
    // ... существующие поля ...
    Center      *Point   `json:"center,omitempty"`    // Центр поиска
    RadiusKm    float64  `json:"radius_km,omitempty"` // Радиус в км
    // ... остальные поля ...
}
```

### 2. Обновление обработчика
**Файл:** `internal/proj/gis/handler/spatial_search.go`

Добавлена обработка параметров:
- `center` - центр поиска в формате "lat,lng"
- `latitude` и `longitude` - отдельные параметры координат
- `radius_km` - радиус в километрах
- `distance` - радиус в формате "10km"

```go
// Поддержка latitude/longitude отдельно
if lat := c.QueryFloat("latitude", 0); lat != 0 {
    if lng := c.QueryFloat("longitude", 0); lng != 0 {
        params.Center = &types.Point{Lat: lat, Lng: lng}
    }
}

// Поддержка distance в формате "10km"
if distanceStr := c.Query("distance"); distanceStr != "" {
    if strings.HasSuffix(distanceStr, "km") {
        distanceStr = strings.TrimSuffix(distanceStr, "km")
        if radius, err := strconv.ParseFloat(distanceStr, 64); err == nil {
            params.RadiusKm = radius
        }
    }
}
```

### 3. Обновление репозитория
**Файл:** `internal/proj/gis/repository/postgis_repo.go`

#### 3.1 Добавлен SQL фильтр по радиусу в кластеризацию
```go
// Добавляем фильтр по радиусу
if params.Center != nil && params.RadiusKm > 0 {
    clusterQuery += fmt.Sprintf(" AND ST_DWithin(lg.location::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)",
        argCount, argCount+1, argCount+2)
    args = append(args, params.Center.Lng, params.Center.Lat, params.RadiusKm*1000) // в метрах
    argCount += 3
}
```

#### 3.2 Обновлена передача параметров в SearchListings
```go
searchParams := types.SearchParams{
    // ... существующие параметры ...
    Center:      params.Center,
    RadiusKm:    params.RadiusKm,
    // ... остальные параметры ...
}
```

#### 3.3 Улучшено логирование
```go
log.Printf("[GetClusters] Starting with params: ZoomLevel=%d, Bounds=%+v, Center=%+v, RadiusKm=%f, Categories=%v, CategoryIDs=%v, MinPrice=%v, MaxPrice=%v",
    params.ZoomLevel, params.Bounds, params.Center, params.RadiusKm, params.Categories, params.CategoryIDs, params.MinPrice, params.MaxPrice)
```

## Результаты тестирования

### 1. API тестирование
**Без distance параметра:**
```bash
curl "http://localhost:3000/api/v1/gis/clusters?bounds=44.7,20.3,44.9,20.6&zoom_level=10&min_price=100&max_price=1000"
# Результат: 7 кластеров
```

**С distance параметром:**
```bash
curl "http://localhost:3000/api/v1/gis/clusters?bounds=44.7,20.3,44.9,20.6&zoom_level=10&latitude=44.8176&longitude=20.4649&distance=10km&min_price=100&max_price=1000"
# Результат: 4 кластера (фильтрация применена!)
```

### 2. Логи backend
```
[GetClusters] Starting with params: ZoomLevel=10, Bounds={...}, Center=&{Lat:44.8176 Lng:20.4649}, RadiusKm=10.000000, Categories=[], CategoryIDs=[], MinPrice=0x..., MaxPrice=0x...
[GetClusters] Query: ... AND ST_DWithin(lg.location::geography, ST_SetSRID(ST_MakePoint($6, $7), 4326)::geography, $8) ...
[GetClusters] Args: [20.3 44.7 20.6 44.9 20 20.4649 44.8176 10000 100 1000]
```

### 3. Frontend тестирование
✅ **Реальные запросы из браузера:**
- Параметр `distance=10km` корректно передается
- Количество кластеров динамически меняется в зависимости от радиуса
- Фильтрация работает в реальном времени при изменении радиуса

## Дополнительные улучшения

1. **Форматирование кода**: Код отформатирован с помощью `gofumpt`
2. **Совместимость**: Поддерживаются все варианты передачи параметров:
   - `center=lat,lng` + `radius_km=10`
   - `latitude=lat` + `longitude=lng` + `distance=10km`
   - `latitude=lat` + `longitude=lng` + `radius_km=10`

## Проверка качества

### Backend
- ✅ **Сборка**: `go build ./...` прошла успешно
- ✅ **Форматирование**: `gofumpt` применен
- ⚠️ **Линтинг**: Есть проблема с версией golangci-lint, но код компилируется

### Тесты
- ✅ **API тестирование**: Все endpoints работают корректно
- ✅ **Frontend интеграция**: Карта работает с новой логикой
- ✅ **Логирование**: Детальные логи для отладки

## Коммит
```bash
git commit -m "fix: добавлена поддержка фильтрации по радиусу в кластеризации карты"
```

## Заключение

Проблема с фильтрацией по радиусу **полностью решена**. Теперь фильтрация работает как для обычного поиска объявлений, так и для кластеризации на карте. Пользователи могут использовать ползунок радиуса для динамической фильтрации результатов в реальном времени.

**Файлы изменены:**
- `internal/proj/gis/types/geo_types.go`
- `internal/proj/gis/handler/spatial_search.go`
- `internal/proj/gis/repository/postgis_repo.go`