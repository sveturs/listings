-- Восстановление PostGIS расширений (на случай если понадобятся)

-- Создаём расширения (они автоматически создадут нужные схемы)
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
