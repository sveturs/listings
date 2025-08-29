-- Миграция для синхронизации изображений с MinIO и исправления путей
-- Автор: System
-- Дата: 2025-08-29

-- =====================================================
-- ЧАСТЬ 1: Исправление путей в marketplace_images
-- =====================================================

-- Создаем временную функцию для исправления путей
CREATE OR REPLACE FUNCTION fix_marketplace_image_paths() RETURNS void AS $$
DECLARE
    img_record RECORD;
    correct_path TEXT;
    correct_url TEXT;
BEGIN
    -- Исправляем пути для всех изображений
    FOR img_record IN 
        SELECT id, listing_id, file_path, file_name 
        FROM marketplace_images 
        WHERE storage_type = 'minio'
    LOOP
        -- Формируем правильный путь
        correct_path := img_record.listing_id || '/' || img_record.file_name;
        
        -- Формируем правильный URL (используем относительный путь)
        correct_url := '/listings/' || correct_path;
        
        -- Обновляем запись
        UPDATE marketplace_images 
        SET 
            file_path = correct_path,
            public_url = correct_url
        WHERE id = img_record.id;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Выполняем исправление путей
SELECT fix_marketplace_image_paths();

-- Удаляем временную функцию
DROP FUNCTION fix_marketplace_image_paths();

-- =====================================================
-- ЧАСТЬ 2: Синхронизация с MinIO для marketplace_listings
-- =====================================================

-- Функция для синхронизации изображений конкретного объявления
CREATE OR REPLACE FUNCTION sync_listing_images_from_audit() RETURNS void AS $$
DECLARE
    listing_record RECORD;
BEGIN
    -- Обрабатываем известные проблемные объявления на основе аудита
    
    -- Объявление 252: добавляем недостающий image5.jpg
    IF NOT EXISTS (
        SELECT 1 FROM marketplace_images 
        WHERE listing_id = 252 AND file_name = 'image5.jpg'
    ) THEN
        INSERT INTO marketplace_images (
            listing_id, file_path, file_name, file_size, 
            content_type, is_main, storage_type, storage_bucket, public_url
        ) VALUES (
            252, '252/image5.jpg', 'image5.jpg', 124000,
            'image/jpeg', false, 'minio', 'listings', '/listings/252/image5.jpg'
        );
    END IF;
    
    -- Исправляем is_main для объявления 252
    UPDATE marketplace_images SET is_main = false WHERE listing_id = 252;
    UPDATE marketplace_images SET is_main = true WHERE listing_id = 252 AND file_name = 'main.jpg';
    
    -- Объявление 183: добавляем недостающие файлы
    IF NOT EXISTS (SELECT 1 FROM marketplace_images WHERE listing_id = 183 AND file_name = 'main.jpg') THEN
        INSERT INTO marketplace_images (
            listing_id, file_path, file_name, file_size, 
            content_type, is_main, storage_type, storage_bucket, public_url
        ) VALUES (
            183, '183/main.jpg', 'main.jpg', 100000,
            'image/jpeg', true, 'minio', 'listings', '/listings/183/main.jpg'
        );
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM marketplace_images WHERE listing_id = 183 AND file_name = 'image2.jpg') THEN
        INSERT INTO marketplace_images (
            listing_id, file_path, file_name, file_size, 
            content_type, is_main, storage_type, storage_bucket, public_url
        ) VALUES (
            183, '183/image2.jpg', 'image2.jpg', 100000,
            'image/jpeg', false, 'minio', 'listings', '/listings/183/image2.jpg'
        );
    END IF;
    
    -- Объявления 250-267: обновляем записи для соответствия MinIO
    -- Для этих объявлений удаляем старые записи с неправильными именами и добавляем правильные
    
    FOR listing_record IN 
        SELECT DISTINCT listing_id 
        FROM marketplace_images 
        WHERE listing_id BETWEEN 250 AND 267
    LOOP
        -- Удаляем записи с паттерном listing_XXX.jpg
        DELETE FROM marketplace_images 
        WHERE listing_id = listing_record.listing_id 
        AND file_name LIKE 'listing_%';
        
        -- Добавляем стандартные изображения если их нет
        IF NOT EXISTS (SELECT 1 FROM marketplace_images WHERE listing_id = listing_record.listing_id AND file_name = 'main.jpg') THEN
            INSERT INTO marketplace_images (
                listing_id, file_path, file_name, file_size, 
                content_type, is_main, storage_type, storage_bucket, public_url
            ) VALUES (
                listing_record.listing_id, 
                listing_record.listing_id || '/main.jpg', 
                'main.jpg', 
                100000,
                'image/jpeg', true, 'minio', 'listings', 
                '/listings/' || listing_record.listing_id || '/main.jpg'
            );
        END IF;
        
        -- Добавляем дополнительные изображения на основе типичного паттерна
        IF NOT EXISTS (SELECT 1 FROM marketplace_images WHERE listing_id = listing_record.listing_id AND file_name = 'image2.jpg') THEN
            INSERT INTO marketplace_images (
                listing_id, file_path, file_name, file_size, 
                content_type, is_main, storage_type, storage_bucket, public_url
            ) VALUES (
                listing_record.listing_id, 
                listing_record.listing_id || '/image2.jpg', 
                'image2.jpg', 
                100000,
                'image/jpeg', false, 'minio', 'listings', 
                '/listings/' || listing_record.listing_id || '/image2.jpg'
            );
        END IF;
        
        IF NOT EXISTS (SELECT 1 FROM marketplace_images WHERE listing_id = listing_record.listing_id AND file_name = 'image3.jpg') THEN
            INSERT INTO marketplace_images (
                listing_id, file_path, file_name, file_size, 
                content_type, is_main, storage_type, storage_bucket, public_url
            ) VALUES (
                listing_record.listing_id, 
                listing_record.listing_id || '/image3.jpg', 
                'image3.jpg', 
                100000,
                'image/jpeg', false, 'minio', 'listings', 
                '/listings/' || listing_record.listing_id || '/image3.jpg'
            );
        END IF;
    END LOOP;
    
END;
$$ LANGUAGE plpgsql;

-- Выполняем синхронизацию
SELECT sync_listing_images_from_audit();

-- Удаляем временную функцию
DROP FUNCTION sync_listing_images_from_audit();

-- =====================================================
-- ЧАСТЬ 3: Исправление URL в storefront_product_images
-- =====================================================

-- Исправляем URL для всех изображений товаров витрин
UPDATE storefront_product_images
SET 
    image_url = REGEXP_REPLACE(image_url, 'http://100\.88\.44\.15:9000', '', 'g'),
    thumbnail_url = REGEXP_REPLACE(thumbnail_url, 'http://100\.88\.44\.15:9000', '', 'g')
WHERE 
    image_url LIKE '%100.88.44.15%' OR 
    thumbnail_url LIKE '%100.88.44.15%';

-- Заменяем localhost:9000 на относительные пути
UPDATE storefront_product_images
SET 
    image_url = REGEXP_REPLACE(image_url, 'http://localhost:9000', '', 'g'),
    thumbnail_url = REGEXP_REPLACE(thumbnail_url, 'http://localhost:9000', '', 'g')
WHERE 
    image_url LIKE '%localhost:9000%' OR 
    thumbnail_url LIKE '%localhost:9000%';

-- =====================================================
-- ЧАСТЬ 4: Обеспечение консистентности данных
-- =====================================================

-- Убедимся, что у каждого объявления есть главное изображение
UPDATE marketplace_images mi1
SET is_main = true
WHERE mi1.id = (
    SELECT id 
    FROM marketplace_images mi2 
    WHERE mi2.listing_id = mi1.listing_id 
    AND mi2.file_name = 'main.jpg'
    LIMIT 1
)
AND NOT EXISTS (
    SELECT 1 
    FROM marketplace_images mi3 
    WHERE mi3.listing_id = mi1.listing_id 
    AND mi3.is_main = true
);

-- Если main.jpg нет, делаем главным первое изображение
UPDATE marketplace_images mi1
SET is_main = true
WHERE mi1.id = (
    SELECT id 
    FROM marketplace_images mi2 
    WHERE mi2.listing_id = mi1.listing_id 
    ORDER BY mi2.id 
    LIMIT 1
)
AND NOT EXISTS (
    SELECT 1 
    FROM marketplace_images mi3 
    WHERE mi3.listing_id = mi1.listing_id 
    AND mi3.is_main = true
);

-- Логирование результатов
DO $$
DECLARE
    marketplace_count INTEGER;
    storefront_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO marketplace_count FROM marketplace_images;
    SELECT COUNT(*) INTO storefront_count FROM storefront_product_images;
    
    RAISE NOTICE 'Миграция завершена. Обработано: % изображений объявлений, % изображений товаров', 
        marketplace_count, storefront_count;
END $$;