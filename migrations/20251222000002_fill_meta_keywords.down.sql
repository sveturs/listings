-- Очистка meta_keywords для L1 категорий
UPDATE categories
SET meta_keywords = '{}'::jsonb,
    updated_at = NOW()
WHERE level = 1;
