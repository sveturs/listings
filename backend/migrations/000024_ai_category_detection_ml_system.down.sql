-- Удаление триггеров
DROP TRIGGER IF EXISTS update_category_detection_feedback_updated_at ON category_detection_feedback;
DROP TRIGGER IF EXISTS update_category_ai_mappings_updated_at ON category_ai_mappings;
DROP TRIGGER IF EXISTS update_category_keyword_weights_updated_at ON category_keyword_weights;

-- Удаление функций
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS get_category_by_ai_hints(VARCHAR, VARCHAR);
DROP FUNCTION IF EXISTS update_mapping_stats(VARCHAR, VARCHAR, INTEGER, BOOLEAN);

-- Удаление представлений
DROP VIEW IF EXISTS category_detection_accuracy;
DROP VIEW IF EXISTS category_detection_errors;

-- Удаление индексов
DROP INDEX IF EXISTS idx_feedback_keywords;
DROP INDEX IF EXISTS idx_feedback_ai_hints;
DROP INDEX IF EXISTS idx_feedback_created_at;
DROP INDEX IF EXISTS idx_feedback_user_confirmed;

DROP INDEX IF EXISTS idx_ai_mappings_domain_type;
DROP INDEX IF EXISTS idx_ai_mappings_category;
DROP INDEX IF EXISTS idx_ai_mappings_active;

DROP INDEX IF EXISTS idx_keyword_weights_keyword;
DROP INDEX IF EXISTS idx_keyword_weights_category;
DROP INDEX IF EXISTS idx_keyword_weights_weight;

DROP INDEX IF EXISTS idx_detection_stats_date;
DROP INDEX IF EXISTS idx_detection_stats_algorithm;

DROP INDEX IF EXISTS idx_detection_cache_key;
DROP INDEX IF EXISTS idx_detection_cache_expires;

-- Удаление таблиц
DROP TABLE IF EXISTS category_detection_cache;
DROP TABLE IF EXISTS category_detection_stats;
DROP TABLE IF EXISTS category_detection_experiments;
DROP TABLE IF EXISTS category_keyword_weights;
DROP TABLE IF EXISTS category_ai_mappings;
DROP TABLE IF EXISTS category_detection_feedback;