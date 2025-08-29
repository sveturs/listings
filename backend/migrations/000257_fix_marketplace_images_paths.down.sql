-- Откат изменений путей к изображениям в marketplace_images
-- Возвращаем пути обратно к формату listing_XXX.jpg

-- Откат для первых изображений
UPDATE marketplace_images 
SET file_path = 'listing_' || listing_id::text || '.jpg',
    public_url = 'http://localhost:9000/listings/listing_' || listing_id::text || '.jpg'
WHERE file_path LIKE '%/image1.jpg';

-- Откат для вторых изображений
UPDATE marketplace_images 
SET file_path = 'listing_' || listing_id::text || '_2.jpg',
    public_url = 'http://localhost:9000/listings/listing_' || listing_id::text || '_2.jpg'
WHERE file_path LIKE '%/image2.jpg';

-- Откат для третьих изображений
UPDATE marketplace_images 
SET file_path = 'listing_' || listing_id::text || '_3.jpg',
    public_url = 'http://localhost:9000/listings/listing_' || listing_id::text || '_3.jpg'
WHERE file_path LIKE '%/image3.jpg';

-- Откат для четвертых изображений
UPDATE marketplace_images 
SET file_path = 'listing_' || listing_id::text || '_4.jpg',
    public_url = 'http://localhost:9000/listings/listing_' || listing_id::text || '_4.jpg'
WHERE file_path LIKE '%/image4.jpg';

-- Откат для пятых изображений
UPDATE marketplace_images 
SET file_path = 'listing_' || listing_id::text || '_5.jpg',
    public_url = 'http://localhost:9000/listings/listing_' || listing_id::text || '_5.jpg'
WHERE file_path LIKE '%/image5.jpg';

-- Аналогично для storefront_product_images (используя image_url вместо file_path)
UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '.jpg',
    thumbnail_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '.jpg'
WHERE image_url LIKE '%/image1.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_2.jpg',
    thumbnail_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_2.jpg'
WHERE image_url LIKE '%/image2.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_3.jpg',
    thumbnail_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_3.jpg'
WHERE image_url LIKE '%/image3.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_4.jpg',
    thumbnail_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_4.jpg'
WHERE image_url LIKE '%/image4.jpg';

UPDATE storefront_product_images 
SET image_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_5.jpg',
    thumbnail_url = 'http://localhost:9000/listings/product_' || storefront_product_id::text || '_5.jpg'
WHERE image_url LIKE '%/image5.jpg';