-- Удаляем атрибут price из дополнительных атрибутов категорий
-- Цена - это базовое поле всех объявлений, а не дополнительный атрибут

-- 1. Удаляем значения атрибута price из существующих объявлений
DELETE FROM listing_attribute_values 
WHERE attribute_id = 2001;

-- 2. Удаляем связи атрибута price с категориями
DELETE FROM category_attribute_mapping 
WHERE attribute_id = 2001;

-- 3. Удаляем атрибут price из групп атрибутов
DELETE FROM attribute_group_items 
WHERE attribute_id = 2001;

-- 4. Удаляем переводы для опций атрибута price (если есть)
DELETE FROM attribute_option_translations 
WHERE attribute_name = 'price';

-- 5. Удаляем сам атрибут price
DELETE FROM category_attributes 
WHERE id = 2001;

-- 6. Удаляем переводы атрибута price
DELETE FROM translations 
WHERE entity_type = 'attribute' AND entity_id = 2001;