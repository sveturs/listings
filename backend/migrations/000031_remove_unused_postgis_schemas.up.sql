-- Удаление неиспользуемых PostGIS расширений и схем
-- tiger/tiger_data - для геокодирования адресов США (не используется)
-- topology - для топологических функций PostGIS (не используется)

-- Удаляем расширения (это автоматически удалит связанные схемы)
DROP EXTENSION IF EXISTS postgis_tiger_geocoder CASCADE;
DROP EXTENSION IF EXISTS postgis_topology CASCADE;

-- Удаляем схемы если они остались
DROP SCHEMA IF EXISTS tiger CASCADE;
DROP SCHEMA IF EXISTS tiger_data CASCADE;
DROP SCHEMA IF EXISTS topology CASCADE;
