-- Исправление путей к изображениям в marketplace_images
-- Изменяем пути с listing_XXX.jpg на XXX/image1.jpg, XXX/image2.jpg и т.д.

-- Обновляем file_path для первых изображений
UPDATE marketplace_images 
SET file_path = listing_id::text || '/image1.jpg',
    public_url = 'http://localhost:9000/listings/' || listing_id::text || '/image1.jpg'
WHERE file_path LIKE 'listing_%' 
  AND file_path NOT LIKE '%_2%' 
  AND file_path NOT LIKE '%_3%'
  AND file_path NOT LIKE '%_4%'
  AND file_path NOT LIKE '%_5%';

-- Обновляем file_path для вторых изображений
UPDATE marketplace_images 
SET file_path = listing_id::text || '/image2.jpg',
    public_url = 'http://localhost:9000/listings/' || listing_id::text || '/image2.jpg'
WHERE file_path LIKE 'listing_%_2.jpg';

-- Обновляем file_path для третьих изображений
UPDATE marketplace_images 
SET file_path = listing_id::text || '/image3.jpg',
    public_url = 'http://localhost:9000/listings/' || listing_id::text || '/image3.jpg'
WHERE file_path LIKE 'listing_%_3.jpg';

-- Обновляем file_path для четвертых изображений
UPDATE marketplace_images 
SET file_path = listing_id::text || '/image4.jpg',
    public_url = 'http://localhost:9000/listings/' || listing_id::text || '/image4.jpg'
WHERE file_path LIKE 'listing_%_4.jpg';

-- Обновляем file_path для пятых изображений
UPDATE marketplace_images 
SET file_path = listing_id::text || '/image5.jpg',
    public_url = 'http://localhost:9000/listings/' || listing_id::text || '/image5.jpg'
WHERE file_path LIKE 'listing_%_5.jpg';

-- Аналогично для storefront_product_images (используя image_url вместо file_path)
UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image1.jpg',
    thumbnail_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image1.jpg'
WHERE image_url LIKE '%product_%' 
  AND image_url NOT LIKE '%_2%' 
  AND image_url NOT LIKE '%_3%'
  AND image_url NOT LIKE '%_4%'
  AND image_url NOT LIKE '%_5%';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image2.jpg',
    thumbnail_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image2.jpg'
WHERE image_url LIKE '%product_%_2.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image3.jpg',
    thumbnail_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image3.jpg'
WHERE image_url LIKE '%product_%_3.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image4.jpg',
    thumbnail_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image4.jpg'
WHERE image_url LIKE '%product_%_4.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image5.jpg',
    thumbnail_url = 'http://localhost:9000/listings/' || storefront_product_id::text || '/image5.jpg'
WHERE image_url LIKE '%product_%_5.jpg';