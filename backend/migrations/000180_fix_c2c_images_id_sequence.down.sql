-- ============================================================================
-- ОТКАТ: Удаление sequence для c2c_images.id
-- ============================================================================

BEGIN;

-- 1. Удалить DEFAULT
ALTER TABLE c2c_images ALTER COLUMN id DROP DEFAULT;

-- 2. Удалить sequence
DROP SEQUENCE IF EXISTS c2c_images_id_seq CASCADE;

COMMIT;
