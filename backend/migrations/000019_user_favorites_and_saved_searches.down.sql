-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_saved_searches_updated_at ON saved_searches;

-- Удаляем функцию обновления updated_at (если не используется другими таблицами)
-- DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем таблицы в правильном порядке (из-за foreign keys)
DROP TABLE IF EXISTS user_car_view_history CASCADE;
DROP TABLE IF EXISTS saved_search_notifications CASCADE;
DROP TABLE IF EXISTS saved_searches CASCADE;

-- Примечание: marketplace_favorites остается, так как она существовала до этой миграции