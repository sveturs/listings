-- Удаление функций
DROP FUNCTION IF EXISTS refresh_density_grid();

-- Удаление материализованного представления
DROP MATERIALIZED VIEW IF EXISTS gis_listing_density_grid;

-- Удаление таблиц
DROP TABLE IF EXISTS gis_filter_analytics;
DROP TABLE IF EXISTS gis_poi_cache;
DROP TABLE IF EXISTS gis_isochrone_cache;