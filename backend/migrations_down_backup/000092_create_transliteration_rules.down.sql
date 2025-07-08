-- Удаление таблицы правил транслитерации
DROP TABLE IF EXISTS transliteration_rules CASCADE;

-- Удаление функции триггера
DROP FUNCTION IF EXISTS update_transliteration_rules_updated_at() CASCADE;