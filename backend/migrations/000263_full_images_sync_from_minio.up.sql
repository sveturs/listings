-- Полная миграция синхронизации изображений с MinIO
-- На основе реального аудита от 29.08.2025
-- Автор: System

-- =====================================================
-- ЧАСТЬ 1: Создание временных таблиц для аудита
-- =====================================================

-- Таблица для хранения результатов аудита MinIO
CREATE TEMP TABLE minio_audit (
    listing_id INTEGER,
    file_name TEXT,
    file_exists BOOLEAN DEFAULT true
);

-- Заполнение таблицы на основе известных данных из аудита
-- Объявление 183: в MinIO есть main.jpg, image2.jpg, image3.jpg
INSERT INTO minio_audit (listing_id, file_name) VALUES
    (183, 'main.jpg'),
    (183, 'image2.jpg'),  
    (183, 'image3.jpg');

-- Объявления 250-267: стандартный набор файлов
DO $$
DECLARE
    lid INTEGER;
BEGIN
    -- Объявления с 4 файлами
    FOR lid IN 250..251 LOOP
        INSERT INTO minio_audit (listing_id, file_name) VALUES
            (lid, 'main.jpg'),
            (lid, 'image2.jpg'),
            (lid, 'image3.jpg'),
            (lid, 'image4.jpg');
    END LOOP;
    
    -- Объявление 252 с 5 файлами (уже исправлено в предыдущей миграции)
    -- Пропускаем
    
    -- Объявление 253 с 3 файлами
    INSERT INTO minio_audit (listing_id, file_name) VALUES
        (253, 'main.jpg'),
        (253, 'image2.jpg'),
        (253, 'image3.jpg');
    
    -- Объявления 254-261 с 4 файлами
    FOR lid IN 254..261 LOOP
        INSERT INTO minio_audit (listing_id, file_name) VALUES
            (lid, 'main.jpg'),
            (lid, 'image2.jpg'),
            (lid, 'image3.jpg'),
            (lid, 'image4.jpg');
    END LOOP;
    
    -- Объявление 262 с 5 файлами
    INSERT INTO minio_audit (listing_id, file_name) VALUES
        (262, 'main.jpg'),
        (262, 'image2.jpg'),
        (262, 'image3.jpg'),
        (262, 'image4.jpg'),
        (262, 'image5.jpg');
    
    -- Объявления 263-265 с 4 файлами
    FOR lid IN 263..265 LOOP
        INSERT INTO minio_audit (listing_id, file_name) VALUES
            (lid, 'main.jpg'),
            (lid, 'image2.jpg'),
            (lid, 'image3.jpg'),
            (lid, 'image4.jpg');
    END LOOP;
    
    -- Объявления 266-267 с 3 файлами
    FOR lid IN 266..267 LOOP
        INSERT INTO minio_audit (listing_id, file_name) VALUES
            (lid, 'main.jpg'),
            (lid, 'image2.jpg'),
            (lid, 'image3.jpg');
    END LOOP;
END $$;

-- =====================================================
-- ЧАСТЬ 2: Удаление старых неправильных записей
-- =====================================================

-- Удаляем записи с неправильными именами файлов
DELETE FROM marketplace_images 
WHERE file_name LIKE 'listing_%'
AND listing_id IN (SELECT DISTINCT listing_id FROM minio_audit);

-- =====================================================
-- ЧАСТЬ 3: Синхронизация изображений из MinIO
-- =====================================================

-- Добавляем недостающие изображения из MinIO
INSERT INTO marketplace_images (
    listing_id, 
    file_path, 
    file_name, 
    file_size,
    content_type, 
    is_main, 
    storage_type, 
    storage_bucket, 
    public_url
)
SELECT 
    ma.listing_id,
    ma.listing_id || '/' || ma.file_name,
    ma.file_name,
    100000, -- Примерный размер
    'image/jpeg',
    (ma.file_name = 'main.jpg'),
    'minio',
    'listings',
    '/listings/' || ma.listing_id || '/' || ma.file_name
FROM minio_audit ma
WHERE NOT EXISTS (
    SELECT 1 
    FROM marketplace_images mi 
    WHERE mi.listing_id = ma.listing_id 
    AND mi.file_name = ma.file_name
)
AND ma.file_exists = true;

-- =====================================================
-- ЧАСТЬ 4: Исправление существующих путей
-- =====================================================

-- Исправляем file_path для всех записей
UPDATE marketplace_images 
SET 
    file_path = listing_id || '/' || file_name,
    public_url = '/listings/' || listing_id || '/' || file_name
WHERE storage_type = 'minio'
AND (
    file_path != listing_id || '/' || file_name
    OR public_url LIKE '%100.88.44.15%'
    OR public_url LIKE '%localhost:9000%'
);

-- =====================================================
-- ЧАСТЬ 5: Обеспечение корректности is_main
-- =====================================================

-- Сбрасываем все is_main
UPDATE marketplace_images 
SET is_main = false 
WHERE listing_id IN (SELECT DISTINCT listing_id FROM minio_audit);

-- Устанавливаем is_main для main.jpg
UPDATE marketplace_images 
SET is_main = true 
WHERE file_name = 'main.jpg'
AND listing_id IN (SELECT DISTINCT listing_id FROM minio_audit);

-- Если main.jpg нет, делаем главным первое изображение
UPDATE marketplace_images mi1
SET is_main = true
WHERE mi1.id = (
    SELECT MIN(id) 
    FROM marketplace_images mi2 
    WHERE mi2.listing_id = mi1.listing_id
)
AND mi1.listing_id IN (SELECT DISTINCT listing_id FROM minio_audit)
AND NOT EXISTS (
    SELECT 1 
    FROM marketplace_images mi3 
    WHERE mi3.listing_id = mi1.listing_id 
    AND mi3.is_main = true
);

-- =====================================================
-- ЧАСТЬ 6: Исправление URL в storefront_product_images
-- =====================================================

-- Убираем IP адреса и localhost:9000 из URL товаров витрин
UPDATE storefront_product_images
SET 
    image_url = REGEXP_REPLACE(
        REGEXP_REPLACE(image_url, 'http://100\.88\.44\.15:9000', '', 'g'),
        'http://localhost:9000', '', 'g'
    ),
    thumbnail_url = REGEXP_REPLACE(
        REGEXP_REPLACE(thumbnail_url, 'http://100\.88\.44\.15:9000', '', 'g'),
        'http://localhost:9000', '', 'g'
    )
WHERE 
    image_url LIKE '%100.88.44.15%' 
    OR image_url LIKE '%localhost:9000%'
    OR thumbnail_url LIKE '%100.88.44.15%'
    OR thumbnail_url LIKE '%localhost:9000%';

-- =====================================================
-- ЧАСТЬ 7: Добавление изображений для объявления 268
-- =====================================================

-- У объявления 268 особый случай с файлом 1756382511472715941.jpg
INSERT INTO marketplace_images (
    listing_id, file_path, file_name, file_size,
    content_type, is_main, storage_type, storage_bucket, public_url
) 
SELECT 268, '268/1756382511472715941.jpg', '1756382511472715941.jpg', 100000,
    'image/jpeg', true, 'minio', 'listings', '/listings/268/1756382511472715941.jpg'
WHERE NOT EXISTS (
    SELECT 1 FROM marketplace_images 
    WHERE listing_id = 268 AND file_name = '1756382511472715941.jpg'
);

-- =====================================================
-- ЧАСТЬ 8: Логирование результатов
-- =====================================================

DO $$
DECLARE
    total_fixed INTEGER;
    total_added INTEGER;
    total_storefront_fixed INTEGER;
BEGIN
    -- Подсчет исправленных записей
    SELECT COUNT(*) INTO total_fixed 
    FROM marketplace_images 
    WHERE listing_id IN (SELECT DISTINCT listing_id FROM minio_audit);
    
    -- Подсчет добавленных записей (приблизительно)
    SELECT COUNT(*) INTO total_added
    FROM marketplace_images
    WHERE created_at >= (CURRENT_TIMESTAMP - INTERVAL '1 minute');
    
    -- Подсчет исправленных товаров витрин
    SELECT COUNT(*) INTO total_storefront_fixed
    FROM storefront_product_images
    WHERE image_url NOT LIKE '%100.88.44.15%'
    AND image_url NOT LIKE '%localhost:9000%';
    
    RAISE NOTICE 'Миграция завершена успешно';
    RAISE NOTICE 'Обработано изображений объявлений: %', total_fixed;
    RAISE NOTICE 'Добавлено новых изображений: %', total_added;
    RAISE NOTICE 'Исправлено изображений товаров витрин: %', total_storefront_fixed;
END $$;