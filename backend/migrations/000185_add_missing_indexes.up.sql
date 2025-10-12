-- Добавление недостающих индексов
-- Эти индексы улучшат производительность JOIN операций

-- 1. Role Permissions - индекс для JOIN с roles
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);

-- 2. Districts - индекс для JOIN с cities
CREATE INDEX IF NOT EXISTS idx_districts_city ON districts(city_id);

-- Примечание: эти индексы помогут в следующих типах запросов:
-- - SELECT * FROM role_permissions WHERE role_id = ?
-- - SELECT * FROM districts WHERE city_id = ?
-- - JOIN операции между roles<->role_permissions и cities<->districts
