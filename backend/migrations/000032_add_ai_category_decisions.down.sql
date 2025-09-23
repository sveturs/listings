-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS trigger_update_ai_category_decisions_updated_at ON ai_category_decisions;
DROP FUNCTION IF EXISTS update_ai_category_decisions_updated_at();

-- Удаляем таблицу
DROP TABLE IF EXISTS ai_category_decisions;