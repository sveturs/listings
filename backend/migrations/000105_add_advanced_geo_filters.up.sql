-- Таблица для кеширования изохрон
CREATE TABLE IF NOT EXISTS gis_isochrone_cache (
    id SERIAL PRIMARY KEY,
    center_point GEOGRAPHY(POINT, 4326) NOT NULL,
    transport_mode VARCHAR(20) NOT NULL,
    max_minutes INTEGER NOT NULL,
    polygon GEOGRAPHY(POLYGON, 4326) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Индексы для изохрон
CREATE INDEX idx_isochrone_cache_center ON gis_isochrone_cache USING GIST(center_point);
CREATE INDEX idx_isochrone_cache_expires ON gis_isochrone_cache(expires_at);
CREATE INDEX idx_isochrone_cache_lookup ON gis_isochrone_cache(transport_mode, max_minutes);

-- Таблица для кеширования POI
CREATE TABLE IF NOT EXISTS gis_poi_cache (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    poi_type VARCHAR(50) NOT NULL,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Индексы для POI
CREATE INDEX idx_poi_cache_type ON gis_poi_cache(poi_type);
CREATE INDEX idx_poi_cache_location ON gis_poi_cache USING GIST(location);
CREATE INDEX idx_poi_cache_expires ON gis_poi_cache(expires_at);

-- Таблица для аналитики использования фильтров
CREATE TABLE IF NOT EXISTS gis_filter_analytics (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    session_id VARCHAR(255) NOT NULL,
    filter_type VARCHAR(50) NOT NULL,
    filter_params JSONB NOT NULL,
    result_count INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индексы для аналитики
CREATE INDEX idx_filter_analytics_user ON gis_filter_analytics(user_id);
CREATE INDEX idx_filter_analytics_session ON gis_filter_analytics(session_id);
CREATE INDEX idx_filter_analytics_type ON gis_filter_analytics(filter_type);
CREATE INDEX idx_filter_analytics_created ON gis_filter_analytics(created_at);

-- Материализованное представление для анализа плотности
CREATE MATERIALIZED VIEW IF NOT EXISTS gis_listing_density_grid AS
WITH grid_cells AS (
    SELECT 
        x::numeric as grid_x,
        y::numeric as grid_y,
        ST_MakeEnvelope(
            x::numeric, 
            y::numeric, 
            (x + 0.005)::numeric, 
            (y + 0.005)::numeric,
            4326
        )::geography as cell
    FROM generate_series(
        18.0, 23.0, 0.005  -- Границы Сербии по долготе
    ) as x,
    generate_series(
        42.0, 46.5, 0.005  -- Границы Сербии по широте
    ) as y
)
SELECT 
    ROW_NUMBER() OVER () as id,
    grid_x,
    grid_y,
    cell,
    COUNT(l.id) as listing_count,
    ST_Area(cell) / 1000000.0 as area_km2,
    CASE 
        WHEN ST_Area(cell) > 0 THEN COUNT(l.id) / (ST_Area(cell) / 1000000.0)
        ELSE 0
    END as density
FROM grid_cells g
LEFT JOIN listings_geo lg ON 
    lg.location IS NOT NULL 
    AND ST_Within(lg.location::geometry, g.cell::geometry)
LEFT JOIN marketplace_listings l ON 
    l.id = lg.listing_id 
    AND l.status = 'active'
GROUP BY grid_x, grid_y, cell;

-- Индекс для материализованного представления
CREATE INDEX idx_density_grid_cell ON gis_listing_density_grid USING GIST(cell);
CREATE INDEX idx_density_grid_density ON gis_listing_density_grid(density);

-- Функция для обновления материализованного представления
CREATE OR REPLACE FUNCTION refresh_density_grid()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY gis_listing_density_grid;
END;
$$ LANGUAGE plpgsql;

-- Комментарии к таблицам
COMMENT ON TABLE gis_isochrone_cache IS 'Кеш изохрон для оптимизации запросов к MapBox API';
COMMENT ON TABLE gis_poi_cache IS 'Кеш точек интереса (школы, больницы, метро и т.д.)';
COMMENT ON TABLE gis_filter_analytics IS 'Аналитика использования геофильтров';
COMMENT ON MATERIALIZED VIEW gis_listing_density_grid IS 'Предрассчитанная сетка плотности объявлений';