# 📍 GIS Модуль - Детальный план реализации Фазы 1

**Цель фазы**: Полностью работающая интерактивная карта с визуализацией товаров и витрин  
**Результат**: Пользователи могут открыть карту, увидеть товары рядом, кликнуть на маркер и перейти к товару  
**Срок**: 4 недели  
**Команда**: Backend Senior + Frontend Senior + DevOps (30%)

## 📅 Неделя 1: Инфраструктура и данные

### День 1-2: Настройка PostGIS

#### 1. Создание миграций БД

**Файл**: `backend/migrations/047_enable_postgis.up.sql`
```sql
-- Включение расширений PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

**Файл**: `backend/migrations/047_enable_postgis.down.sql`
```sql
DROP EXTENSION IF EXISTS postgis CASCADE;
DROP EXTENSION IF EXISTS postgis_topology CASCADE;
DROP EXTENSION IF EXISTS pg_trgm CASCADE;
```

**Файл**: `backend/migrations/048_create_listings_geo.up.sql`
```sql
-- Таблица геоданных объявлений
CREATE TABLE listings_geo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    -- Геохеши для оптимизации
    geohash4 VARCHAR(4) GENERATED ALWAYS AS (ST_GeoHash(location, 4)) STORED,
    geohash6 VARCHAR(6) GENERATED ALWAYS AS (ST_GeoHash(location, 6)) STORED,
    geohash8 VARCHAR(8) GENERATED ALWAYS AS (ST_GeoHash(location, 8)) STORED,
    -- Денормализованные данные для быстрого доступа
    city VARCHAR(100),
    district VARCHAR(100),
    postal_code VARCHAR(10),
    -- Приватность локации
    is_precise BOOLEAN DEFAULT true,
    blur_radius_meters INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_listings_geo_location ON listings_geo USING GIST(location);
CREATE INDEX idx_listings_geo_geohash4 ON listings_geo(geohash4);
CREATE INDEX idx_listings_geo_geohash6 ON listings_geo(geohash6);
CREATE INDEX idx_listings_geo_city ON listings_geo(city);
CREATE INDEX idx_listings_geo_listing ON listings_geo(listing_id);
```

**Файл**: `backend/migrations/048_create_listings_geo.down.sql`
```sql
DROP TABLE IF EXISTS listings_geo;
```

#### 2. Обновление Docker Compose

**Файл**: `docker-compose.yml` (добавить в postgres секцию)
```yaml
postgres:
  image: postgis/postgis:16-3.4
  environment:
    - POSTGRES_DB=svetu
    - POSTGRES_USER=postgres
    - POSTGRES_PASSWORD=postgres
  volumes:
    - postgres_data:/var/lib/postgresql/data
  ports:
    - "5432:5432"
```

### День 3-4: Миграция координат

#### 3. Скрипт миграции существующих координат

**Файл**: `backend/scripts/migrate_listings_coordinates.go`
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/svetu?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Миграция координат из существующих полей latitude/longitude
    query := `
        INSERT INTO listings_geo (listing_id, location, city, district, postal_code)
        SELECT 
            ml.id,
            ST_Point(ml.longitude, ml.latitude)::geography,
            ml.city,
            ml.district,
            ml.postal_code
        FROM marketplace_listings ml
        WHERE ml.latitude IS NOT NULL 
        AND ml.longitude IS NOT NULL
        AND NOT EXISTS (
            SELECT 1 FROM listings_geo lg WHERE lg.listing_id = ml.id
        )
    `
    
    result, err := db.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
    
    affected, _ := result.RowsAffected()
    fmt.Printf("Migrated %d listings\n", affected)
}
```

### День 5: Базовая структура GIS модуля

#### 4. Создание структуры модуля

**Создать директории**:
```
backend/internal/proj/gis/
├── handler/
├── service/
├── repository/
└── types/
```

**Файл**: `backend/internal/proj/gis/types/geo_types.go`
```go
package types

import (
    "time"
    "github.com/google/uuid"
)

// Point представляет географическую точку
type Point struct {
    Lat float64 `json:"lat" validate:"required,min=-90,max=90"`
    Lng float64 `json:"lng" validate:"required,min=-180,max=180"`
}

// Bounds представляет границы карты
type Bounds struct {
    Southwest Point `json:"southwest" validate:"required"`
    Northeast Point `json:"northeast" validate:"required"`
}

// GeoListing представляет объявление с геоданными
type GeoListing struct {
    ID         uuid.UUID `json:"id"`
    Title      string    `json:"title"`
    Price      float64   `json:"price"`
    Location   Point     `json:"location"`
    Distance   float64   `json:"distance,omitempty"`
    Category   Category  `json:"category"`
    Thumbnail  string    `json:"thumbnail"`
    Seller     Seller    `json:"seller"`
    IsPrecise  bool      `json:"is_precise"`
    BlurRadius int       `json:"blur_radius,omitempty"`
}

// Category информация о категории
type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Icon string `json:"icon,omitempty"`
}

// Seller информация о продавце
type Seller struct {
    ID     uuid.UUID `json:"id"`
    Name   string    `json:"name"`
    Rating float64   `json:"rating"`
}

// Cluster представляет кластер маркеров
type Cluster struct {
    ID     string  `json:"id"`
    Center Point   `json:"center"`
    Count  int     `json:"count"`
    Bounds *Bounds `json:"bounds,omitempty"`
}

// SearchBoundsRequest запрос поиска в границах
type SearchBoundsRequest struct {
    Bounds     Bounds  `json:"bounds" validate:"required"`
    Zoom       int     `json:"zoom" validate:"min=1,max=20"`
    Clustered  bool    `json:"clustered"`
    CategoryID *int    `json:"category_id,omitempty"`
    PriceMin   *float64 `json:"price_min,omitempty"`
    PriceMax   *float64 `json:"price_max,omitempty"`
}

// SearchBoundsResponse ответ поиска в границах
type SearchBoundsResponse struct {
    Listings []GeoListing `json:"listings,omitempty"`
    Clusters []Cluster    `json:"clusters,omitempty"`
    Total    int          `json:"total"`
}
```

## 📅 Неделя 2: Backend API

### День 6-7: Repository слой

**Файл**: `backend/internal/proj/gis/repository/postgis_repo.go`
```go
package repository

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/google/uuid"
    "github.com/lib/pq"
    
    "backend/internal/proj/gis/types"
)

type PostGISRepository struct {
    db *sql.DB
}

func NewPostGISRepository(db *sql.DB) *PostGISRepository {
    return &PostGISRepository{db: db}
}

// SearchInBounds ищет объявления в границах карты
func (r *PostGISRepository) SearchInBounds(
    ctx context.Context,
    bounds types.Bounds,
    zoom int,
    filters map[string]interface{},
) ([]types.GeoListing, error) {
    query := `
        SELECT 
            lg.listing_id,
            ml.title,
            ml.price,
            ST_Y(lg.location::geometry) as lat,
            ST_X(lg.location::geometry) as lng,
            lg.is_precise,
            lg.blur_radius_meters,
            mc.id as category_id,
            mc.name as category_name,
            COALESCE(mi.url, '') as thumbnail,
            us.id as seller_id,
            us.username as seller_name,
            COALESCE(us.rating, 0) as seller_rating
        FROM listings_geo lg
        JOIN marketplace_listings ml ON lg.listing_id = ml.id
        JOIN marketplace_categories mc ON ml.category_id = mc.id
        LEFT JOIN marketplace_images mi ON mi.listing_id = ml.id AND mi.is_primary = true
        JOIN user_storefronts us ON ml.storefront_id = us.id
        WHERE 
            ST_Intersects(
                lg.location,
                ST_MakeEnvelope($1, $2, $3, $4, 4326)::geography
            )
            AND ml.is_active = true
    `
    
    args := []interface{}{
        bounds.Southwest.Lng, bounds.Southwest.Lat,
        bounds.Northeast.Lng, bounds.Northeast.Lat,
    }
    
    // Добавляем фильтры
    argCount := 4
    if categoryID, ok := filters["category_id"].(int); ok && categoryID > 0 {
        argCount++
        query += fmt.Sprintf(" AND ml.category_id = $%d", argCount)
        args = append(args, categoryID)
    }
    
    if priceMin, ok := filters["price_min"].(float64); ok {
        argCount++
        query += fmt.Sprintf(" AND ml.price >= $%d", argCount)
        args = append(args, priceMin)
    }
    
    if priceMax, ok := filters["price_max"].(float64); ok {
        argCount++
        query += fmt.Sprintf(" AND ml.price <= $%d", argCount)
        args = append(args, priceMax)
    }
    
    // Лимит в зависимости от зума
    limit := 100
    if zoom < 10 {
        limit = 50
    } else if zoom > 15 {
        limit = 500
    }
    
    query += fmt.Sprintf(" LIMIT %d", limit)
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("search in bounds: %w", err)
    }
    defer rows.Close()
    
    var listings []types.GeoListing
    for rows.Next() {
        var l types.GeoListing
        var sellerID uuid.UUID
        
        err := rows.Scan(
            &l.ID,
            &l.Title,
            &l.Price,
            &l.Location.Lat,
            &l.Location.Lng,
            &l.IsPrecise,
            &l.BlurRadius,
            &l.Category.ID,
            &l.Category.Name,
            &l.Thumbnail,
            &sellerID,
            &l.Seller.Name,
            &l.Seller.Rating,
        )
        if err != nil {
            return nil, fmt.Errorf("scan listing: %w", err)
        }
        
        l.Seller.ID = sellerID
        
        // Размываем координаты если нужно
        if !l.IsPrecise && l.BlurRadius > 0 {
            l.Location = r.blurLocation(l.Location, l.BlurRadius)
        }
        
        listings = append(listings, l)
    }
    
    return listings, nil
}

// GetListingLocation получает координаты объявления
func (r *PostGISRepository) GetListingLocation(
    ctx context.Context,
    listingID uuid.UUID,
) (*types.Point, error) {
    var point types.Point
    var isPrecise bool
    var blurRadius int
    
    query := `
        SELECT 
            ST_Y(location::geometry) as lat,
            ST_X(location::geometry) as lng,
            is_precise,
            blur_radius_meters
        FROM listings_geo
        WHERE listing_id = $1
    `
    
    err := r.db.QueryRowContext(ctx, query, listingID).Scan(
        &point.Lat, &point.Lng, &isPrecise, &blurRadius,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("get listing location: %w", err)
    }
    
    if !isPrecise && blurRadius > 0 {
        point = r.blurLocation(point, blurRadius)
    }
    
    return &point, nil
}

// UpdateListingLocation обновляет координаты объявления
func (r *PostGISRepository) UpdateListingLocation(
    ctx context.Context,
    listingID uuid.UUID,
    location types.Point,
    isPrecise bool,
) error {
    query := `
        INSERT INTO listings_geo (listing_id, location, is_precise, blur_radius_meters)
        VALUES ($1, ST_Point($2, $3)::geography, $4, $5)
        ON CONFLICT (listing_id) 
        DO UPDATE SET 
            location = EXCLUDED.location,
            is_precise = EXCLUDED.is_precise,
            blur_radius_meters = EXCLUDED.blur_radius_meters,
            updated_at = NOW()
    `
    
    blurRadius := 0
    if !isPrecise {
        blurRadius = 500 // 500 метров по умолчанию
    }
    
    _, err := r.db.ExecContext(ctx, query, 
        listingID, location.Lng, location.Lat, isPrecise, blurRadius,
    )
    if err != nil {
        return fmt.Errorf("update listing location: %w", err)
    }
    
    return nil
}

// CreateClusters создает кластеры для текущего зума
func (r *PostGISRepository) CreateClusters(
    ctx context.Context,
    bounds types.Bounds,
    zoom int,
    minClusterSize int,
) ([]types.Cluster, error) {
    // Размер сетки в зависимости от зума
    gridSize := 0.1
    if zoom >= 10 && zoom < 14 {
        gridSize = 0.01
    } else if zoom >= 14 {
        gridSize = 0.001
    }
    
    query := `
        WITH clustered AS (
            SELECT 
                COUNT(*) as count,
                AVG(ST_Y(lg.location::geometry)) as center_lat,
                AVG(ST_X(lg.location::geometry)) as center_lng,
                MIN(ST_Y(lg.location::geometry)) as min_lat,
                MIN(ST_X(lg.location::geometry)) as min_lng,
                MAX(ST_Y(lg.location::geometry)) as max_lat,
                MAX(ST_X(lg.location::geometry)) as max_lng,
                FLOOR(ST_X(lg.location::geometry) / $5) as grid_x,
                FLOOR(ST_Y(lg.location::geometry) / $5) as grid_y
            FROM listings_geo lg
            JOIN marketplace_listings ml ON lg.listing_id = ml.id
            WHERE 
                ST_Intersects(
                    lg.location,
                    ST_MakeEnvelope($1, $2, $3, $4, 4326)::geography
                )
                AND ml.is_active = true
            GROUP BY grid_x, grid_y
            HAVING COUNT(*) >= $6
        )
        SELECT 
            md5(grid_x::text || grid_y::text) as id,
            count,
            center_lat,
            center_lng,
            min_lat,
            min_lng,
            max_lat,
            max_lng
        FROM clustered
    `
    
    rows, err := r.db.QueryContext(ctx, query,
        bounds.Southwest.Lng, bounds.Southwest.Lat,
        bounds.Northeast.Lng, bounds.Northeast.Lat,
        gridSize, minClusterSize,
    )
    if err != nil {
        return nil, fmt.Errorf("create clusters: %w", err)
    }
    defer rows.Close()
    
    var clusters []types.Cluster
    for rows.Next() {
        var c types.Cluster
        var minLat, minLng, maxLat, maxLng float64
        
        err := rows.Scan(
            &c.ID, &c.Count,
            &c.Center.Lat, &c.Center.Lng,
            &minLat, &minLng, &maxLat, &maxLng,
        )
        if err != nil {
            return nil, fmt.Errorf("scan cluster: %w", err)
        }
        
        c.Bounds = &types.Bounds{
            Southwest: types.Point{Lat: minLat, Lng: minLng},
            Northeast: types.Point{Lat: maxLat, Lng: maxLng},
        }
        
        clusters = append(clusters, c)
    }
    
    return clusters, nil
}

// blurLocation размывает координаты для приватности
func (r *PostGISRepository) blurLocation(p types.Point, radiusMeters int) types.Point {
    // Простое размытие - добавляем случайное смещение
    // В продакшене использовать более сложный алгоритм
    metersPerDegree := 111111.0
    offset := float64(radiusMeters) / metersPerDegree
    
    // Случайное смещение в пределах радиуса
    // TODO: использовать криптографически стойкий рандом
    p.Lat += (0.5 - 0.5) * offset * 2
    p.Lng += (0.5 - 0.5) * offset * 2
    
    return p
}
```

### День 8-9: Service слой

**Файл**: `backend/internal/proj/gis/service/spatial_service.go`
```go
package service

import (
    "context"
    "fmt"
    
    "github.com/google/uuid"
    
    "backend/internal/proj/gis/repository"
    "backend/internal/proj/gis/types"
)

type SpatialService struct {
    repo  *repository.PostGISRepository
    cache CacheService // Redis кэш
}

func NewSpatialService(repo *repository.PostGISRepository, cache CacheService) *SpatialService {
    return &SpatialService{
        repo:  repo,
        cache: cache,
    }
}

// SearchInBounds выполняет поиск в границах с кластеризацией
func (s *SpatialService) SearchInBounds(
    ctx context.Context,
    req types.SearchBoundsRequest,
) (*types.SearchBoundsResponse, error) {
    // Валидация границ
    if err := s.validateBounds(req.Bounds); err != nil {
        return nil, fmt.Errorf("invalid bounds: %w", err)
    }
    
    // Проверяем кэш
    cacheKey := s.buildCacheKey(req)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        return cached.(*types.SearchBoundsResponse), nil
    }
    
    filters := make(map[string]interface{})
    if req.CategoryID != nil {
        filters["category_id"] = *req.CategoryID
    }
    if req.PriceMin != nil {
        filters["price_min"] = *req.PriceMin
    }
    if req.PriceMax != nil {
        filters["price_max"] = *req.PriceMax
    }
    
    response := &types.SearchBoundsResponse{}
    
    if req.Clustered && req.Zoom < 15 {
        // Используем кластеризацию для малых зумов
        clusters, err := s.repo.CreateClusters(ctx, req.Bounds, req.Zoom, 5)
        if err != nil {
            return nil, fmt.Errorf("create clusters: %w", err)
        }
        response.Clusters = clusters
        response.Total = len(clusters)
    } else {
        // Возвращаем отдельные объявления
        listings, err := s.repo.SearchInBounds(ctx, req.Bounds, req.Zoom, filters)
        if err != nil {
            return nil, fmt.Errorf("search listings: %w", err)
        }
        response.Listings = listings
        response.Total = len(listings)
    }
    
    // Кэшируем результат на 5 минут
    s.cache.Set(ctx, cacheKey, response, 300)
    
    return response, nil
}

// GetListingLocation возвращает координаты объявления
func (s *SpatialService) GetListingLocation(
    ctx context.Context,
    listingID uuid.UUID,
) (*types.Point, error) {
    location, err := s.repo.GetListingLocation(ctx, listingID)
    if err != nil {
        return nil, fmt.Errorf("get location: %w", err)
    }
    
    return location, nil
}

// UpdateListingLocation обновляет координаты объявления
func (s *SpatialService) UpdateListingLocation(
    ctx context.Context,
    listingID uuid.UUID,
    location types.Point,
    isPrecise bool,
) error {
    // Валидация координат
    if err := s.validatePoint(location); err != nil {
        return fmt.Errorf("invalid location: %w", err)
    }
    
    err := s.repo.UpdateListingLocation(ctx, listingID, location, isPrecise)
    if err != nil {
        return fmt.Errorf("update location: %w", err)
    }
    
    // Инвалидируем кэш
    s.cache.InvalidatePattern(ctx, "gis:bounds:*")
    
    return nil
}

// validateBounds проверяет корректность границ
func (s *SpatialService) validateBounds(bounds types.Bounds) error {
    if bounds.Southwest.Lat >= bounds.Northeast.Lat {
        return fmt.Errorf("southwest lat must be less than northeast lat")
    }
    if bounds.Southwest.Lng >= bounds.Northeast.Lng {
        return fmt.Errorf("southwest lng must be less than northeast lng")
    }
    
    // Проверяем разумные границы (не слишком большая область)
    latDiff := bounds.Northeast.Lat - bounds.Southwest.Lat
    lngDiff := bounds.Northeast.Lng - bounds.Southwest.Lng
    
    if latDiff > 10 || lngDiff > 10 {
        return fmt.Errorf("bounds area too large")
    }
    
    return nil
}

// validatePoint проверяет корректность точки
func (s *SpatialService) validatePoint(p types.Point) error {
    if p.Lat < -90 || p.Lat > 90 {
        return fmt.Errorf("latitude must be between -90 and 90")
    }
    if p.Lng < -180 || p.Lng > 180 {
        return fmt.Errorf("longitude must be between -180 and 180")
    }
    return nil
}

// buildCacheKey создает ключ кэша
func (s *SpatialService) buildCacheKey(req types.SearchBoundsRequest) string {
    return fmt.Sprintf("gis:bounds:%f:%f:%f:%f:%d:%v:%v:%v:%v",
        req.Bounds.Southwest.Lat, req.Bounds.Southwest.Lng,
        req.Bounds.Northeast.Lat, req.Bounds.Northeast.Lng,
        req.Zoom, req.Clustered, req.CategoryID, req.PriceMin, req.PriceMax,
    )
}
```

### День 10: API Handlers

**Файл**: `backend/internal/proj/gis/handler/spatial_search.go`
```go
package handler

import (
    "net/http"
    "strconv"
    "strings"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    
    "backend/internal/proj/gis/service"
    "backend/internal/proj/gis/types"
    "backend/pkg/utils"
)

type SpatialHandler struct {
    service *service.SpatialService
}

func NewSpatialHandler(service *service.SpatialService) *SpatialHandler {
    return &SpatialHandler{service: service}
}

// SearchBounds godoc
// @Summary Поиск объявлений в границах карты
// @Description Возвращает объявления или кластеры в зависимости от зума
// @Tags GIS
// @Accept json
// @Produce json
// @Param bounds query string true "Границы карты в формате sw_lat,sw_lng,ne_lat,ne_lng"
// @Param zoom query int false "Уровень зума карты" minimum(1) maximum(20)
// @Param clustered query bool false "Использовать кластеризацию" default(true)
// @Param category_id query int false "ID категории"
// @Param price_min query number false "Минимальная цена"
// @Param price_max query number false "Максимальная цена"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.SearchBoundsResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Router /api/v1/gis/search/bounds [get]
func (h *SpatialHandler) SearchBounds(c *fiber.Ctx) error {
    // Парсим границы
    boundsStr := c.Query("bounds")
    if boundsStr == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.gis.boundsRequired")
    }
    
    parts := strings.Split(boundsStr, ",")
    if len(parts) != 4 {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.gis.invalidBounds")
    }
    
    swLat, _ := strconv.ParseFloat(parts[0], 64)
    swLng, _ := strconv.ParseFloat(parts[1], 64)
    neLat, _ := strconv.ParseFloat(parts[2], 64)
    neLng, _ := strconv.ParseFloat(parts[3], 64)
    
    req := types.SearchBoundsRequest{
        Bounds: types.Bounds{
            Southwest: types.Point{Lat: swLat, Lng: swLng},
            Northeast: types.Point{Lat: neLat, Lng: neLng},
        },
        Zoom:      c.QueryInt("zoom", 12),
        Clustered: c.QueryBool("clustered", true),
    }
    
    // Опциональные фильтры
    if categoryID := c.QueryInt("category_id", 0); categoryID > 0 {
        req.CategoryID = &categoryID
    }
    
    if priceMin, err := strconv.ParseFloat(c.Query("price_min"), 64); err == nil {
        req.PriceMin = &priceMin
    }
    
    if priceMax, err := strconv.ParseFloat(c.Query("price_max"), 64); err == nil {
        req.PriceMax = &priceMax
    }
    
    // Выполняем поиск
    result, err := h.service.SearchInBounds(c.Context(), req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.gis.searchFailed")
    }
    
    return utils.SuccessResponse(c, result)
}

// GetListingLocation godoc
// @Summary Получить координаты объявления
// @Description Возвращает географические координаты объявления
// @Tags GIS
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.Point}
// @Failure 404 {object} utils.ErrorResponseSwag
// @Router /api/v1/gis/listings/{id}/location [get]
func (h *SpatialHandler) GetListingLocation(c *fiber.Ctx) error {
    listingID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
    }
    
    location, err := h.service.GetListingLocation(c.Context(), listingID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.gis.getLocationFailed")
    }
    
    if location == nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, "errors.gis.locationNotFound")
    }
    
    return utils.SuccessResponse(c, location)
}

// UpdateListingLocation godoc
// @Summary Обновить координаты объявления
// @Description Обновляет географические координаты объявления
// @Tags GIS
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Param location body types.Point true "Новые координаты"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 403 {object} utils.ErrorResponseSwag
// @Router /api/v1/gis/listings/{id}/location [post]
// @Security BearerAuth
func (h *SpatialHandler) UpdateListingLocation(c *fiber.Ctx) error {
    listingID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidID")
    }
    
    var req struct {
        Location  types.Point `json:"location" validate:"required"`
        IsPrecise bool       `json:"is_precise"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "errors.invalidInput")
    }
    
    // TODO: Проверить права доступа к объявлению
    
    err = h.service.UpdateListingLocation(c.Context(), listingID, req.Location, req.IsPrecise)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.gis.updateLocationFailed")
    }
    
    return utils.SuccessResponse(c, map[string]bool{"success": true})
}
```

**Файл**: `backend/internal/proj/gis/handler/routes.go`
```go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "backend/internal/middleware"
)

// RegisterRoutes регистрирует маршруты GIS модуля
func RegisterRoutes(app *fiber.App, h *SpatialHandler, auth *middleware.AuthMiddleware) {
    gis := app.Group("/api/v1/gis")
    
    // Публичные маршруты
    gis.Get("/search/bounds", h.SearchBounds)
    gis.Get("/listings/:id/location", h.GetListingLocation)
    
    // Защищенные маршруты
    gis.Post("/listings/:id/location", auth.RequireAuth, h.UpdateListingLocation)
}
```

## 📅 Неделя 3: Frontend карта

### День 11-12: Настройка Mapbox

#### 1. Установка зависимостей

```bash
cd frontend/svetu
yarn add mapbox-gl @types/mapbox-gl
yarn add react-map-gl
```

#### 2. Конфигурация

**Файл**: `frontend/svetu/src/config/types.ts` (добавить)
```typescript
export interface EnvVariables {
  // ... существующие поля
  NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN: string;
  NEXT_PUBLIC_MAPBOX_STYLE_URL: string;
}
```

**Файл**: `frontend/svetu/src/config/index.ts` (обновить)
```typescript
export const config = {
  // ... существующие поля
  mapbox: {
    accessToken: process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN || '',
    styleUrl: process.env.NEXT_PUBLIC_MAPBOX_STYLE_URL || 'mapbox://styles/mapbox/streets-v12',
  },
};
```

### День 13-14: Компонент карты

**Файл**: `frontend/svetu/src/components/GIS/Map/InteractiveMap.tsx`
```typescript
import { useEffect, useRef, useState, useCallback } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { config } from '@/config';
import type { GeoListing, Bounds } from '@/types/gis';

// Установка токена Mapbox
mapboxgl.accessToken = config.mapbox.accessToken;

interface InteractiveMapProps {
  center?: [number, number];
  zoom?: number;
  listings?: GeoListing[];
  onBoundsChange?: (bounds: Bounds) => void;
  onMarkerClick?: (listing: GeoListing) => void;
  className?: string;
}

export const InteractiveMap: React.FC<InteractiveMapProps> = ({
  center = [20.4568, 44.8178], // Белград
  zoom = 12,
  listings = [],
  onBoundsChange,
  onMarkerClick,
  className = '',
}) => {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const markers = useRef<mapboxgl.Marker[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  
  // Инициализация карты
  useEffect(() => {
    if (!mapContainer.current || map.current) return;
    
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: config.mapbox.styleUrl,
      center,
      zoom,
      attributionControl: false,
    });
    
    // Добавляем контролы
    map.current.addControl(
      new mapboxgl.NavigationControl(),
      'top-right'
    );
    
    map.current.addControl(
      new mapboxgl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true,
        },
        trackUserLocation: true,
        showUserHeading: true,
      }),
      'top-right'
    );
    
    // События карты
    map.current.on('load', () => {
      setIsLoading(false);
    });
    
    map.current.on('moveend', () => {
      if (onBoundsChange && map.current) {
        const bounds = map.current.getBounds();
        onBoundsChange({
          southwest: {
            lat: bounds.getSouth(),
            lng: bounds.getWest(),
          },
          northeast: {
            lat: bounds.getNorth(),
            lng: bounds.getEast(),
          },
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
    markers.current.forEach((marker) => marker.remove());
    markers.current = [];
    
    // Добавляем новые маркеры
    listings.forEach((listing) => {
      // Создаем элемент маркера
      const el = document.createElement('div');
      el.className = 'custom-marker';
      el.style.width = '40px';
      el.style.height = '40px';
      el.style.backgroundImage = `url(/api/categories/${listing.category.id}/icon.svg)`;
      el.style.backgroundSize = 'cover';
      el.style.borderRadius = '50%';
      el.style.backgroundColor = '#fff';
      el.style.border = '3px solid #fff';
      el.style.boxShadow = '0 2px 8px rgba(0,0,0,0.3)';
      el.style.cursor = 'pointer';
      el.style.transition = 'transform 0.2s';
      
      // Эффект при наведении
      el.addEventListener('mouseenter', () => {
        el.style.transform = 'scale(1.1)';
      });
      
      el.addEventListener('mouseleave', () => {
        el.style.transform = 'scale(1)';
      });
      
      // Размытие для приватных локаций
      if (!listing.is_precise) {
        el.style.opacity = '0.7';
        el.classList.add('blurred-location');
      }
      
      // Создаем маркер
      const marker = new mapboxgl.Marker(el)
        .setLngLat([listing.location.lng, listing.location.lat])
        .addTo(map.current!);
      
      // Клик по маркеру
      el.addEventListener('click', () => {
        if (onMarkerClick) {
          onMarkerClick(listing);
        }
      });
      
      // Попап с информацией
      const popupContent = `
        <div class="p-2 max-w-xs">
          ${listing.thumbnail ? `
            <img 
              src="${listing.thumbnail}" 
              alt="${listing.title}"
              class="w-full h-32 object-cover rounded mb-2"
            />
          ` : ''}
          <h3 class="font-semibold text-base line-clamp-2">${listing.title}</h3>
          <p class="text-lg font-bold text-primary">${listing.price.toLocaleString('sr-RS')} РСД</p>
          <div class="flex items-center gap-2 text-sm text-gray-600 mt-1">
            <span>${listing.seller.name}</span>
            ${listing.seller.rating > 0 ? `
              <span class="flex items-center gap-1">
                <svg class="w-4 h-4 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
                ${listing.seller.rating.toFixed(1)}
              </span>
            ` : ''}
          </div>
          ${listing.distance ? `
            <p class="text-sm text-gray-500 mt-1">
              ${listing.distance < 1000 
                ? `${Math.round(listing.distance)}м` 
                : `${(listing.distance / 1000).toFixed(1)}км`
              }
            </p>
          ` : ''}
        </div>
      `;
      
      const popup = new mapboxgl.Popup({
        offset: 25,
        closeButton: false,
        className: 'custom-popup',
      }).setHTML(popupContent);
      
      marker.setPopup(popup);
      markers.current.push(marker);
    });
  }, [listings, isLoading, onMarkerClick]);
  
  return (
    <div className={`relative w-full h-full ${className}`}>
      <div ref={mapContainer} className="w-full h-full" />
      
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
          <div className="flex flex-col items-center gap-2">
            <div className="loading loading-spinner loading-lg text-primary"></div>
            <span className="text-sm text-gray-600">Загрузка карты...</span>
          </div>
        </div>
      )}
      
      {/* Стили для попапов */}
      <style jsx global>{`
        .custom-popup {
          font-family: inherit;
        }
        
        .custom-popup .mapboxgl-popup-content {
          padding: 0;
          border-radius: 0.5rem;
          box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
        }
        
        .custom-popup .mapboxgl-popup-tip {
          border-top-color: white;
        }
        
        .blurred-location::after {
          content: '';
          position: absolute;
          inset: -10px;
          background: radial-gradient(circle, transparent 30%, rgba(255,255,255,0.8) 70%);
          pointer-events: none;
        }
      `}</style>
    </div>
  );
};
```

### День 15: Контролы и фильтры карты

**Файл**: `frontend/svetu/src/components/GIS/Map/MapControls.tsx`
```typescript
import { useState } from 'react';
import { MapIcon, ListBulletIcon } from '@heroicons/react/24/outline';

interface MapControlsProps {
  viewMode: 'map' | 'list';
  onViewModeChange: (mode: 'map' | 'list') => void;
  mapStyle: string;
  onMapStyleChange: (style: string) => void;
}

export const MapControls: React.FC<MapControlsProps> = ({
  viewMode,
  onViewModeChange,
  mapStyle,
  onMapStyleChange,
}) => {
  const [showStyleMenu, setShowStyleMenu] = useState(false);
  
  const mapStyles = [
    { id: 'streets-v12', name: 'Улицы', icon: '🗺️' },
    { id: 'satellite-streets-v12', name: 'Спутник', icon: '🛰️' },
    { id: 'light-v11', name: 'Светлая', icon: '☀️' },
    { id: 'dark-v11', name: 'Темная', icon: '🌙' },
  ];
  
  return (
    <div className="absolute top-4 left-4 z-10 flex flex-col gap-2">
      {/* Переключатель вида */}
      <div className="btn-group bg-white shadow-lg">
        <button
          className={`btn btn-sm ${viewMode === 'map' ? 'btn-primary' : 'btn-ghost'}`}
          onClick={() => onViewModeChange('map')}
        >
          <MapIcon className="w-4 h-4" />
          Карта
        </button>
        <button
          className={`btn btn-sm ${viewMode === 'list' ? 'btn-primary' : 'btn-ghost'}`}
          onClick={() => onViewModeChange('list')}
        >
          <ListBulletIcon className="w-4 h-4" />
          Список
        </button>
      </div>
      
      {/* Стиль карты */}
      {viewMode === 'map' && (
        <div className="relative">
          <button
            className="btn btn-sm btn-ghost bg-white shadow-lg"
            onClick={() => setShowStyleMenu(!showStyleMenu)}
          >
            Стиль карты
          </button>
          
          {showStyleMenu && (
            <div className="absolute top-full left-0 mt-1 bg-white rounded-lg shadow-xl p-2 min-w-[150px]">
              {mapStyles.map((style) => (
                <button
                  key={style.id}
                  className={`w-full text-left px-3 py-2 rounded hover:bg-gray-100 flex items-center gap-2 ${
                    mapStyle === style.id ? 'bg-primary/10 text-primary' : ''
                  }`}
                  onClick={() => {
                    onMapStyleChange(style.id);
                    setShowStyleMenu(false);
                  }}
                >
                  <span className="text-lg">{style.icon}</span>
                  <span className="text-sm">{style.name}</span>
                </button>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
};
```

**Файл**: `frontend/svetu/src/components/GIS/Map/MapFilters.tsx`
```typescript
import { useState, useEffect } from 'react';
import { FunnelIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { useCategories } from '@/hooks/useCategories';

interface MapFiltersProps {
  onFiltersChange: (filters: MapFilters) => void;
  initialFilters?: MapFilters;
}

export interface MapFilters {
  category_id?: number;
  price_min?: number;
  price_max?: number;
  radius?: number;
  has_delivery?: boolean;
}

export const MapFilters: React.FC<MapFiltersProps> = ({
  onFiltersChange,
  initialFilters = {},
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [filters, setFilters] = useState<MapFilters>(initialFilters);
  const { categories } = useCategories();
  
  const handleFilterChange = (key: keyof MapFilters, value: any) => {
    const newFilters = { ...filters };
    
    if (value === null || value === undefined || value === '') {
      delete newFilters[key];
    } else {
      newFilters[key] = value;
    }
    
    setFilters(newFilters);
  };
  
  const applyFilters = () => {
    onFiltersChange(filters);
    setIsOpen(false);
  };
  
  const clearFilters = () => {
    setFilters({});
    onFiltersChange({});
    setIsOpen(false);
  };
  
  const activeFiltersCount = Object.keys(filters).length;
  
  return (
    <>
      {/* Кнопка фильтров */}
      <button
        className="btn btn-sm bg-white shadow-lg"
        onClick={() => setIsOpen(true)}
      >
        <FunnelIcon className="w-4 h-4" />
        Фильтры
        {activeFiltersCount > 0 && (
          <span className="badge badge-primary badge-sm ml-1">
            {activeFiltersCount}
          </span>
        )}
      </button>
      
      {/* Панель фильтров */}
      {isOpen && (
        <>
          {/* Оверлей */}
          <div
            className="fixed inset-0 bg-black/50 z-40"
            onClick={() => setIsOpen(false)}
          />
          
          {/* Панель */}
          <div className="fixed right-0 top-0 h-full w-80 bg-white shadow-xl z-50 overflow-y-auto">
            <div className="p-4 border-b">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold">Фильтры</h3>
                <button
                  className="btn btn-ghost btn-sm btn-circle"
                  onClick={() => setIsOpen(false)}
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>
            </div>
            
            <div className="p-4 space-y-6">
              {/* Категория */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Категория</span>
                </label>
                <select
                  className="select select-bordered w-full"
                  value={filters.category_id || ''}
                  onChange={(e) => handleFilterChange('category_id', e.target.value ? Number(e.target.value) : null)}
                >
                  <option value="">Все категории</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>
              
              {/* Цена */}
              <div>
                <label className="label">
                  <span className="label-text">Цена (РСД)</span>
                </label>
                <div className="flex gap-2">
                  <input
                    type="number"
                    placeholder="От"
                    className="input input-bordered w-full"
                    value={filters.price_min || ''}
                    onChange={(e) => handleFilterChange('price_min', e.target.value ? Number(e.target.value) : null)}
                  />
                  <input
                    type="number"
                    placeholder="До"
                    className="input input-bordered w-full"
                    value={filters.price_max || ''}
                    onChange={(e) => handleFilterChange('price_max', e.target.value ? Number(e.target.value) : null)}
                  />
                </div>
              </div>
              
              {/* Радиус поиска */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    Радиус поиска: {filters.radius || 5} км
                  </span>
                </label>
                <input
                  type="range"
                  min="1"
                  max="50"
                  value={filters.radius || 5}
                  onChange={(e) => handleFilterChange('radius', Number(e.target.value))}
                  className="range range-primary"
                />
                <div className="w-full flex justify-between text-xs px-2 mt-1">
                  <span>1км</span>
                  <span>25км</span>
                  <span>50км</span>
                </div>
              </div>
              
              {/* Доставка */}
              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">Только с доставкой</span>
                  <input
                    type="checkbox"
                    className="checkbox checkbox-primary"
                    checked={filters.has_delivery || false}
                    onChange={(e) => handleFilterChange('has_delivery', e.target.checked || null)}
                  />
                </label>
              </div>
            </div>
            
            {/* Кнопки действий */}
            <div className="p-4 border-t flex gap-2">
              <button
                className="btn btn-ghost flex-1"
                onClick={clearFilters}
                disabled={activeFiltersCount === 0}
              >
                Очистить
              </button>
              <button
                className="btn btn-primary flex-1"
                onClick={applyFilters}
              >
                Применить
              </button>
            </div>
          </div>
        </>
      )}
    </>
  );
};
```

## 📅 Неделя 4: Интеграция и оптимизация

### День 16-17: API hooks и интеграция

**Файл**: `frontend/svetu/src/hooks/gis/useGeoSearch.ts`
```typescript
import { useQuery, useMutation } from '@tanstack/react-query';
import { api } from '@/services/api';
import type { SearchBoundsRequest, SearchBoundsResponse, Point } from '@/types/gis';

export const useGeoSearch = () => {
  const searchInBounds = useQuery<SearchBoundsResponse, Error, SearchBoundsResponse, [string, SearchBoundsRequest]>({
    queryKey: ['gis', 'search-bounds'],
    queryFn: async ({ queryKey }) => {
      const [_, params] = queryKey;
      const boundsStr = `${params.bounds.southwest.lat},${params.bounds.southwest.lng},${params.bounds.northeast.lat},${params.bounds.northeast.lng}`;
      
      const searchParams = new URLSearchParams({
        bounds: boundsStr,
        zoom: params.zoom.toString(),
        clustered: params.clustered.toString(),
      });
      
      if (params.category_id) {
        searchParams.append('category_id', params.category_id.toString());
      }
      if (params.price_min !== undefined) {
        searchParams.append('price_min', params.price_min.toString());
      }
      if (params.price_max !== undefined) {
        searchParams.append('price_max', params.price_max.toString());
      }
      
      const response = await api.get(`/gis/search/bounds?${searchParams}`);
      return response.data;
    },
    enabled: false, // Запускаем вручную
    staleTime: 5 * 60 * 1000, // 5 минут
  });
  
  return {
    searchInBounds,
    isLoading: searchInBounds.isLoading,
    data: searchInBounds.data,
    error: searchInBounds.error,
  };
};

export const useListingLocation = (listingId?: string) => {
  return useQuery<Point | null>({
    queryKey: ['gis', 'listing-location', listingId],
    queryFn: async () => {
      if (!listingId) return null;
      const response = await api.get(`/gis/listings/${listingId}/location`);
      return response.data;
    },
    enabled: !!listingId,
  });
};

export const useUpdateListingLocation = () => {
  return useMutation({
    mutationFn: async ({ listingId, location, isPrecise }: {
      listingId: string;
      location: Point;
      isPrecise: boolean;
    }) => {
      const response = await api.post(`/gis/listings/${listingId}/location`, {
        location,
        is_precise: isPrecise,
      });
      return response.data;
    },
  });
};
```

**Файл**: `frontend/svetu/src/hooks/gis/useGeolocation.ts`
```typescript
import { useState, useCallback } from 'react';

interface GeolocationState {
  location: GeolocationPosition | null;
  error: GeolocationPositionError | null;
  isLoading: boolean;
}

export const useGeolocation = () => {
  const [state, setState] = useState<GeolocationState>({
    location: null,
    error: null,
    isLoading: false,
  });
  
  const requestLocation = useCallback(() => {
    if (!navigator.geolocation) {
      setState({
        location: null,
        error: {
          code: 0,
          message: 'Геолокация не поддерживается браузером',
          PERMISSION_DENIED: 1,
          POSITION_UNAVAILABLE: 2,
          TIMEOUT: 3,
        } as GeolocationPositionError,
        isLoading: false,
      });
      return;
    }
    
    setState((prev) => ({ ...prev, isLoading: true }));
    
    navigator.geolocation.getCurrentPosition(
      (position) => {
        setState({
          location: position,
          error: null,
          isLoading: false,
        });
      },
      (error) => {
        setState({
          location: null,
          error,
          isLoading: false,
        });
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000, // 5 минут
      }
    );
  }, []);
  
  const watchLocation = useCallback(() => {
    if (!navigator.geolocation) return;
    
    const watchId = navigator.geolocation.watchPosition(
      (position) => {
        setState({
          location: position,
          error: null,
          isLoading: false,
        });
      },
      (error) => {
        setState((prev) => ({
          ...prev,
          error,
          isLoading: false,
        }));
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 5000,
      }
    );
    
    return () => {
      navigator.geolocation.clearWatch(watchId);
    };
  }, []);
  
  return {
    ...state,
    requestLocation,
    watchLocation,
  };
};
```

### День 18-19: Страница с картой

**Файл**: `frontend/svetu/src/app/[locale]/map/page.tsx`
```typescript
'use client';

import { useState, useCallback, useEffect } from 'react';
import { InteractiveMap } from '@/components/GIS/Map/InteractiveMap';
import { MapControls } from '@/components/GIS/Map/MapControls';
import { MapFilters, type MapFilters as MapFiltersType } from '@/components/GIS/Map/MapFilters';
import { useGeoSearch } from '@/hooks/gis/useGeoSearch';
import { useRouter } from 'next/navigation';
import type { Bounds, GeoListing } from '@/types/gis';

export default function MapPage() {
  const router = useRouter();
  const [viewMode, setViewMode] = useState<'map' | 'list'>('map');
  const [mapStyle, setMapStyle] = useState('streets-v12');
  const [bounds, setBounds] = useState<Bounds | null>(null);
  const [filters, setFilters] = useState<MapFiltersType>({});
  const [zoom, setZoom] = useState(12);
  
  const { searchInBounds, data, isLoading } = useGeoSearch();
  
  // Поиск при изменении границ или фильтров
  useEffect(() => {
    if (!bounds) return;
    
    searchInBounds.refetch({
      bounds,
      zoom,
      clustered: zoom < 15,
      ...filters,
    });
  }, [bounds, zoom, filters]);
  
  const handleBoundsChange = useCallback((newBounds: Bounds) => {
    setBounds(newBounds);
  }, []);
  
  const handleMarkerClick = useCallback((listing: GeoListing) => {
    router.push(`/listings/${listing.id}`);
  }, [router]);
  
  const handleFiltersChange = useCallback((newFilters: MapFiltersType) => {
    setFilters(newFilters);
  }, []);
  
  return (
    <div className="h-screen flex flex-col">
      {/* Верхняя панель */}
      <div className="navbar bg-base-100 shadow-md z-20">
        <div className="flex-1">
          <h1 className="text-xl font-bold">Карта товаров</h1>
        </div>
        <div className="flex-none gap-2">
          <MapFilters 
            onFiltersChange={handleFiltersChange}
            initialFilters={filters}
          />
          <button
            className="btn btn-ghost btn-sm"
            onClick={() => router.back()}
          >
            Закрыть
          </button>
        </div>
      </div>
      
      {/* Основной контент */}
      <div className="flex-1 relative">
        {viewMode === 'map' ? (
          <>
            <InteractiveMap
              listings={data?.listings || []}
              onBoundsChange={handleBoundsChange}
              onMarkerClick={handleMarkerClick}
            />
            <MapControls
              viewMode={viewMode}
              onViewModeChange={setViewMode}
              mapStyle={mapStyle}
              onMapStyleChange={setMapStyle}
            />
            
            {/* Счетчик результатов */}
            {data && (
              <div className="absolute bottom-4 left-4 bg-white rounded-lg shadow-lg px-4 py-2">
                <span className="text-sm font-medium">
                  {data.total > 0 ? (
                    <>Найдено: {data.total} {data.clusters ? 'кластеров' : 'товаров'}</>
                  ) : (
                    'Ничего не найдено'
                  )}
                </span>
              </div>
            )}
            
            {/* Индикатор загрузки */}
            {isLoading && (
              <div className="absolute top-20 right-4 bg-white rounded-lg shadow-lg px-4 py-2">
                <div className="flex items-center gap-2">
                  <div className="loading loading-spinner loading-sm"></div>
                  <span className="text-sm">Поиск...</span>
                </div>
              </div>
            )}
          </>
        ) : (
          <div className="p-4">
            {/* Список товаров */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {data?.listings?.map((listing) => (
                <div
                  key={listing.id}
                  className="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
                  onClick={() => handleMarkerClick(listing)}
                >
                  {listing.thumbnail && (
                    <figure>
                      <img
                        src={listing.thumbnail}
                        alt={listing.title}
                        className="h-48 w-full object-cover"
                      />
                    </figure>
                  )}
                  <div className="card-body">
                    <h2 className="card-title line-clamp-2">{listing.title}</h2>
                    <p className="text-xl font-bold text-primary">
                      {listing.price.toLocaleString('sr-RS')} РСД
                    </p>
                    <div className="text-sm text-gray-600">
                      <p>{listing.seller.name}</p>
                      {listing.distance && (
                        <p>
                          {listing.distance < 1000
                            ? `${Math.round(listing.distance)}м`
                            : `${(listing.distance / 1000).toFixed(1)}км`}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
```

### День 20: Оптимизация и тестирование

#### 1. Кэширование на стороне Redis

**Файл**: `backend/internal/proj/gis/repository/cache_repo.go`
```go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type CacheRepository struct {
    client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
    return &CacheRepository{client: client}
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    
    return json.Unmarshal([]byte(val), dest)
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *CacheRepository) InvalidatePattern(ctx context.Context, pattern string) error {
    var cursor uint64
    for {
        keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return err
        }
        
        if len(keys) > 0 {
            if err := r.client.Del(ctx, keys...).Err(); err != nil {
                return err
            }
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    
    return nil
}
```

#### 2. Миграции для тестовых данных

**Файл**: `backend/fixtures/001_gis_test_data.sql`
```sql
-- Тестовые объявления для Белграда с реальными координатами
WITH test_locations AS (
    SELECT * FROM (VALUES
        ('Стари град', 44.8176, 20.4633),
        ('Врачар', 44.7989, 20.4766),
        ('Нови Београд', 44.8096, 20.4049),
        ('Земун', 44.8433, 20.4011),
        ('Палилула', 44.8154, 20.4859),
        ('Звездара', 44.8025, 20.5074),
        ('Чукарица', 44.7456, 20.4202),
        ('Раковица', 44.7589, 20.4569)
    ) AS t(district, base_lat, base_lng)
)
INSERT INTO listings_geo (listing_id, location, city, district, postal_code, is_precise)
SELECT 
    ml.id,
    ST_Point(
        tl.base_lng + (random() - 0.5) * 0.02,
        tl.base_lat + (random() - 0.5) * 0.02
    )::geography,
    'Београд',
    tl.district,
    CASE tl.district
        WHEN 'Стари град' THEN '11000'
        WHEN 'Врачар' THEN '11010'
        WHEN 'Нови Београд' THEN '11070'
        WHEN 'Земун' THEN '11080'
        WHEN 'Палилула' THEN '11050'
        WHEN 'Звездара' THEN '11060'
        WHEN 'Чукарица' THEN '11030'
        WHEN 'Раковица' THEN '11090'
    END,
    CASE WHEN random() > 0.8 THEN false ELSE true END -- 20% приватных локаций
FROM marketplace_listings ml
CROSS JOIN LATERAL (
    SELECT * FROM test_locations ORDER BY random() LIMIT 1
) tl
WHERE ml.city = 'Београд'
AND NOT EXISTS (
    SELECT 1 FROM listings_geo lg WHERE lg.listing_id = ml.id
)
LIMIT 1000;

-- Добавляем blur_radius для приватных локаций
UPDATE listings_geo
SET blur_radius_meters = 300 + (random() * 400)::int
WHERE is_precise = false;
```

#### 3. Мониторинг производительности

**Файл**: `backend/internal/proj/gis/middleware/metrics.go`
```go
package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    gisRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "gis_request_duration_seconds",
        Help: "Duration of GIS requests in seconds",
        Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
    }, []string{"method", "endpoint", "status"})
    
    gisActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "gis_active_requests",
        Help: "Number of active GIS requests",
    })
    
    gisListingsReturned = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "gis_listings_returned",
        Help: "Number of listings returned per request",
        Buckets: []float64{0, 10, 25, 50, 100, 250, 500, 1000},
    }, []string{"endpoint"})
)

func MetricsMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        gisActiveRequests.Inc()
        defer gisActiveRequests.Dec()
        
        err := c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Response().StatusCode()
        
        gisRequestDuration.WithLabelValues(
            c.Method(),
            c.Path(),
            fmt.Sprintf("%d", status),
        ).Observe(duration)
        
        return err
    }
}
```

### День 21: Финальная интеграция и тестирование

#### 1. Интеграционные тесты

**Файл**: `backend/internal/proj/gis/handler/spatial_search_test.go`
```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "backend/internal/proj/gis/types"
)

func TestSearchBounds(t *testing.T) {
    // Тест успешного поиска
    t.Run("successful search", func(t *testing.T) {
        req := httptest.NewRequest(
            "GET",
            "/api/v1/gis/search/bounds?bounds=44.7,20.3,44.9,20.6&zoom=12",
            nil,
        )
        
        resp, err := app.Test(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        
        var result struct {
            Data types.SearchBoundsResponse `json:"data"`
        }
        err = json.NewDecoder(resp.Body).Decode(&result)
        require.NoError(t, err)
        
        assert.NotEmpty(t, result.Data.Listings)
        assert.Greater(t, result.Data.Total, 0)
    })
    
    // Тест кластеризации
    t.Run("clustered search", func(t *testing.T) {
        req := httptest.NewRequest(
            "GET",
            "/api/v1/gis/search/bounds?bounds=44.7,20.3,44.9,20.6&zoom=10&clustered=true",
            nil,
        )
        
        resp, err := app.Test(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        
        var result struct {
            Data types.SearchBoundsResponse `json:"data"`
        }
        err = json.NewDecoder(resp.Body).Decode(&result)
        require.NoError(t, err)
        
        assert.NotEmpty(t, result.Data.Clusters)
        assert.Empty(t, result.Data.Listings)
    })
    
    // Тест с фильтрами
    t.Run("filtered search", func(t *testing.T) {
        req := httptest.NewRequest(
            "GET",
            "/api/v1/gis/search/bounds?bounds=44.8,20.4,44.82,20.42&zoom=15&category_id=1&price_max=5000",
            nil,
        )
        
        resp, err := app.Test(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        
        var result struct {
            Data types.SearchBoundsResponse `json:"data"`
        }
        err = json.NewDecoder(resp.Body).Decode(&result)
        require.NoError(t, err)
        
        // Проверяем что фильтры применились
        for _, listing := range result.Data.Listings {
            assert.Equal(t, 1, listing.Category.ID)
            assert.LessOrEqual(t, listing.Price, 5000.0)
        }
    })
}
```

#### 2. Обновление swagger документации

```bash
cd backend
make generate-types
```

#### 3. Финальная проверка

**Чек-лист перед завершением фазы 1:**

- [ ] PostGIS установлен и работает
- [ ] Миграции применены успешно
- [ ] API endpoints работают корректно
- [ ] Карта отображается и загружается < 2 сек
- [ ] Маркеры кластеризуются при малом зуме
- [ ] Попапы показывают информацию о товаре
- [ ] Фильтры работают корректно
- [ ] Геолокация пользователя работает
- [ ] Кэширование Redis настроено
- [ ] Метрики Prometheus собираются
- [ ] Тесты проходят успешно
- [ ] Документация Swagger обновлена
- [ ] Производительность соответствует требованиям

## 📊 Метрики успеха фазы 1

После завершения фазы 1 должны быть достигнуты следующие показатели:

1. **Производительность**:
   - Загрузка карты < 2 секунды
   - Отображение 1000+ маркеров без лагов
   - P95 latency API < 200ms

2. **Функциональность**:
   - Карта отображает товары с корректными координатами
   - Кластеризация работает при zoom < 15
   - Фильтры по категории и цене применяются
   - Клик по маркеру ведет на страницу товара

3. **Пользовательский опыт**:
   - CTR с маркера на товар > 5%
   - Отсутствие критических ошибок
   - Мобильная версия работает корректно

## 🚀 Запуск и проверка

1. **Применение миграций**:
```bash
cd backend
./migrator migrate
./migrator migrate --only-fixtures
```

2. **Запуск сервисов**:
```bash
docker-compose up -d postgres redis
cd backend && go run ./cmd/api/main.go
cd frontend/svetu && yarn dev -p 3001
```

3. **Проверка функциональности**:
- Открыть http://localhost:3001/map
- Проверить отображение маркеров
- Применить фильтры
- Кликнуть на маркер
- Проверить мобильную версию

## 📝 Документация для команды

После завершения фазы 1 необходимо:

1. Обновить README с инструкциями по настройке PostGIS
2. Добавить примеры использования API в Postman коллекцию
3. Создать руководство по добавлению координат к объявлениям
4. Подготовить метрики дашборд в Grafana

---

**Фаза 1 готова к реализации!** 🎯

При возникновении вопросов обращайтесь к техническому руководителю проекта.