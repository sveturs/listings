-- ============================================================================
-- МИГРАЦИЯ: Исправление sequence для c2c_images.id
-- Дата: 2025-10-11
-- Описание: Добавление SERIAL для автоинкремента ID в таблице c2c_images
-- Проблема: При создании через LIKE не копируются DEFAULT и sequence
-- ============================================================================

BEGIN;

-- 1. Создать sequence для c2c_images.id
CREATE SEQUENCE IF NOT EXISTS c2c_images_id_seq;

-- 2. Установить текущее значение sequence = максимальный ID + 1
SELECT setval('c2c_images_id_seq', COALESCE((SELECT MAX(id) FROM c2c_images), 0) + 1, false);

-- 3. Установить DEFAULT для колонки id
ALTER TABLE c2c_images ALTER COLUMN id SET DEFAULT nextval('c2c_images_id_seq');

-- 4. Установить владельца sequence (важно для правильной работы pg_dump)
ALTER SEQUENCE c2c_images_id_seq OWNED BY c2c_images.id;

COMMIT;

-- ============================================================================
-- РЕЗУЛЬТАТ: Теперь c2c_images.id будет автоматически генерироваться
-- ============================================================================
