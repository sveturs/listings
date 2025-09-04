-- Откат миграции 000042: Удаление атрибутов для категории Pets

-- Удаляем связи категорий с атрибутами
DELETE FROM unified_category_attributes
WHERE attribute_id IN (
    SELECT id FROM unified_attributes WHERE code LIKE 'pet_%'
);

-- Удаляем переводы атрибутов
DELETE FROM translations
WHERE entity_type = 'unified_attribute'
  AND entity_id IN (
    SELECT id FROM unified_attributes WHERE code LIKE 'pet_%'
  );

-- Удаляем сами атрибуты
DELETE FROM unified_attributes WHERE code LIKE 'pet_%';