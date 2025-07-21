-- Скрипт для удаления всех категорий и атрибутов из БД
-- ВНИМАНИЕ: Это удалит ВСЕ данные, связанные с категориями и атрибутами!

BEGIN;

-- 1. Удаляем связанные данные с атрибутами объявлений
DELETE FROM listing_attribute_values;

-- 2. Удаляем переводы опций атрибутов
DELETE FROM attribute_option_translations;

-- 3. Удаляем маппинг категорий и атрибутов
DELETE FROM category_attribute_mapping;

-- 4. Удаляем атрибуты из групп (если есть)
DELETE FROM attribute_group_items;

-- 5. Удаляем группы атрибутов категорий
DELETE FROM category_attribute_groups;

-- 6. Удаляем сами атрибуты категорий
DELETE FROM category_attributes;

-- 7. Удаляем группы атрибутов
DELETE FROM attribute_groups;

-- 8. Удаляем связанные с категориями поля в объявлениях (обнуляем category_id)
UPDATE marketplace_listings SET category_id = NULL;

-- 9. Удаляем все категории
DELETE FROM marketplace_categories;

-- 10. Удаляем атрибуты продуктов (storefront)
DELETE FROM product_variant_attribute_values;
DELETE FROM product_variant_attributes;
DELETE FROM storefront_product_attributes;

-- Выводим результаты
SELECT 'Удаление завершено' as status;
SELECT 'marketplace_categories' as table_name, COUNT(*) as remaining_records FROM marketplace_categories
UNION ALL
SELECT 'category_attributes', COUNT(*) FROM category_attributes
UNION ALL
SELECT 'category_attribute_mapping', COUNT(*) FROM category_attribute_mapping
UNION ALL
SELECT 'listing_attribute_values', COUNT(*) FROM listing_attribute_values
UNION ALL
SELECT 'category_attribute_groups', COUNT(*) FROM category_attribute_groups
UNION ALL
SELECT 'attribute_groups', COUNT(*) FROM attribute_groups
UNION ALL
SELECT 'attribute_option_translations', COUNT(*) FROM attribute_option_translations;

COMMIT;