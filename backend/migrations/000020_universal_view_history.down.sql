-- Удаляем триггер
DROP TRIGGER IF EXISTS trg_check_return_visit ON user_view_history;

-- Удаляем функции
DROP FUNCTION IF EXISTS check_return_visit();
DROP FUNCTION IF EXISTS update_view_statistics(DATE);

-- Удаляем таблицы
DROP TABLE IF EXISTS view_statistics;
DROP TABLE IF EXISTS user_view_history;

-- Примечание: таблица user_car_view_history не восстанавливается,
-- так как данные уже мигрированы в новую структуру.
-- Если необходимо восстановить старую таблицу, это должно быть сделано отдельно.