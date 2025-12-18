-- ═══════════════════════════════════════════════════════════════════
-- 🔴 КРИТИЧЕСКИ ВАЖНО: ПОЛНАЯ ОЧИСТКА ПЕРЕД НОВОЙ СИСТЕМОЙ КАТЕГОРИЙ
-- ═══════════════════════════════════════════════════════════════════
--
-- Миграция: 20251216000002_cleanup.up.sql
-- Дата: 2025-12-16
-- Описание: Удаление всех старых категорий и атрибутов перед внедрением новой системы
--
-- УДАЛЯЕМ:
--   ❌ Все категории
--   ❌ Все атрибуты и их значения
--   ❌ Все связи category_attributes
--
-- СОХРАНЯЕМ:
--   ✅ listings - отвязываем от категорий (category_id = NULL)
--   ✅ products - отвязываем от категорий (category_id = NULL)
--   ✅ Все остальные таблицы без изменений
--
-- ВАЖНО: Эта миграция должна быть выполнена ПОСЛЕ очистки OpenSearch!
--   ./scripts/cleanup_opensearch.sh
--
-- ═══════════════════════════════════════════════════════════════════

BEGIN;

-- Логирование начала
DO $$
BEGIN
    RAISE NOTICE '════════════════════════════════════════════════════════';
    RAISE NOTICE '🔴 НАЧАЛО ОЧИСТКИ БД';
    RAISE NOTICE '════════════════════════════════════════════════════════';
END $$;

-- 1. Убираем NOT NULL constraint с category_id в listings (если есть)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'listings'
        AND column_name = 'category_id'
        AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE listings ALTER COLUMN category_id DROP NOT NULL;
        RAISE NOTICE '1a. Убран NOT NULL constraint с listings.category_id';
    END IF;
END $$;

-- 1b. Отвязываем listings от категорий
DO $$
DECLARE
    listings_affected INTEGER;
BEGIN
    UPDATE listings SET category_id = NULL WHERE category_id IS NOT NULL;
    GET DIAGNOSTICS listings_affected = ROW_COUNT;
    RAISE NOTICE '1b. Отвязано listings от категорий: %', listings_affected;
END $$;

-- 2. Убираем NOT NULL constraint с category_id в products (если есть)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'products') THEN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'products'
            AND column_name = 'category_id'
            AND is_nullable = 'NO'
        ) THEN
            ALTER TABLE products ALTER COLUMN category_id DROP NOT NULL;
            RAISE NOTICE '2a. Убран NOT NULL constraint с products.category_id';
        END IF;
    END IF;
END $$;

-- 2b. Отвязываем products от категорий (если таблица существует)
DO $$
DECLARE
    products_affected INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'products') THEN
        EXECUTE 'UPDATE products SET category_id = NULL WHERE category_id IS NOT NULL';
        GET DIAGNOSTICS products_affected = ROW_COUNT;
        RAISE NOTICE '2b. Отвязано products от категорий: %', products_affected;
    ELSE
        RAISE NOTICE '2b. Таблица products не найдена (пропускаем)';
    END IF;
END $$;

-- 3. Удаляем связи атрибутов с категориями
DO $$
DECLARE
    cat_attrs_count INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'category_attributes') THEN
        SELECT COUNT(*) INTO cat_attrs_count FROM category_attributes;
        TRUNCATE TABLE category_attributes CASCADE;
        RAISE NOTICE '3. Удалено связей category_attributes: %', cat_attrs_count;
    ELSE
        RAISE NOTICE '3. Таблица category_attributes не найдена (создастся позже)';
    END IF;
END $$;

-- 4. Удаляем значения атрибутов
DO $$
DECLARE
    attr_values_count INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'attribute_values') THEN
        SELECT COUNT(*) INTO attr_values_count FROM attribute_values;
        TRUNCATE TABLE attribute_values CASCADE;
        RAISE NOTICE '4. Удалено значений атрибутов: %', attr_values_count;
    ELSE
        RAISE NOTICE '4. Таблица attribute_values не найдена (создастся позже)';
    END IF;
END $$;

-- 5. Удаляем атрибуты
DO $$
DECLARE
    attributes_count INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'attributes') THEN
        SELECT COUNT(*) INTO attributes_count FROM attributes;
        TRUNCATE TABLE attributes CASCADE;
        RAISE NOTICE '5. Удалено атрибутов: %', attributes_count;
    ELSE
        RAISE NOTICE '5. Таблица attributes не найдена (создастся позже)';
    END IF;
END $$;

-- 6. Удаляем ВСЕ категории
DO $$
DECLARE
    categories_count INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'categories') THEN
        SELECT COUNT(*) INTO categories_count FROM categories;
        TRUNCATE TABLE categories CASCADE;
        RAISE NOTICE '6. Удалено категорий: %', categories_count;
    ELSE
        RAISE NOTICE '6. Таблица categories не найдена (создастся позже)';
    END IF;
END $$;

-- 7. Сбрасываем последовательности (если используются SERIAL)
DO $$
BEGIN
    -- Проверяем и сбрасываем sequences если они существуют
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'categories_id_seq' AND relkind = 'S') THEN
        ALTER SEQUENCE categories_id_seq RESTART WITH 1;
        RAISE NOTICE '7. Сброшена последовательность categories_id_seq';
    END IF;

    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'attributes_id_seq' AND relkind = 'S') THEN
        ALTER SEQUENCE attributes_id_seq RESTART WITH 1;
        RAISE NOTICE '7. Сброшена последовательность attributes_id_seq';
    END IF;
END $$;

-- Логирование завершения
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '════════════════════════════════════════════════════════';
    RAISE NOTICE '✅ ОЧИСТКА ЗАВЕРШЕНА УСПЕШНО';
    RAISE NOTICE '════════════════════════════════════════════════════════';
    RAISE NOTICE '';
    RAISE NOTICE 'Удалены:';
    RAISE NOTICE '  ❌ Все категории';
    RAISE NOTICE '  ❌ Все атрибуты и их значения';
    RAISE NOTICE '  ❌ Все связи category_attributes';
    RAISE NOTICE '';
    RAISE NOTICE 'Сохранены:';
    RAISE NOTICE '  ✅ Все таблицы listings, products (category_id = NULL)';
    RAISE NOTICE '  ✅ Все остальные данные без изменений';
    RAISE NOTICE '';
    RAISE NOTICE 'Следующий шаг:';
    RAISE NOTICE '  Выполнить DDL миграции для новой структуры категорий';
    RAISE NOTICE '════════════════════════════════════════════════════════';
END $$;

COMMIT;
