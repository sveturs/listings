-- Удаляем функцию обновления
DROP FUNCTION IF EXISTS refresh_map_items_cache();

-- Удаляем materialized view
DROP MATERIALIZED VIEW IF EXISTS map_items_cache;