-- Удаление проблемных объявлений с ID из логов ошибок
BEGIN;

-- Для безопасности сохраним информацию об удаляемых объявлениях
CREATE TEMPORARY TABLE broken_listings AS
SELECT * FROM marketplace_listings
WHERE id IN (223, 222, 221, 220, 218, 217, 216, 215, 214, 212, 211, 210, 209, 208, 207, 206, 205);

-- Сохраним информацию о изображениях
CREATE TEMPORARY TABLE broken_images AS
SELECT * FROM marketplace_images
WHERE listing_id IN (SELECT id FROM broken_listings);

-- Удаляем сначала зависимые данные
DELETE FROM marketplace_images WHERE listing_id IN (SELECT id FROM broken_listings);
DELETE FROM listing_attribute_values WHERE listing_id IN (SELECT id FROM broken_listings);
DELETE FROM marketplace_favorites WHERE listing_id IN (SELECT id FROM broken_listings);

-- Затем удаляем сами объявления
DELETE FROM marketplace_listings WHERE id IN (SELECT id FROM broken_listings);

-- Посмотрим, что удалили
SELECT 'Deleted listings:', count(*) FROM broken_listings;
SELECT 'Deleted images:', count(*) FROM broken_images;

-- Если все в порядке, подтверждаем транзакцию
COMMIT;