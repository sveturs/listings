-- Удаляем триггер
DROP TRIGGER IF EXISTS problem_status_change_trigger ON problem_shipments;

-- Удаляем функцию триггера
DROP FUNCTION IF EXISTS track_problem_status_changes();

-- Удаляем таблицы
DROP TABLE IF EXISTS problem_status_history;
DROP TABLE IF EXISTS problem_comments;