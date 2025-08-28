-- Удаляем технические атрибуты из категории Автомобили (кроме power_hp, который был добавлен отдельно)
DELETE FROM category_attribute_mapping 
WHERE category_id = 1003 
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('body_type', 'drivetrain', 'doors', 'seats')
);

-- Удаляем только атрибут drivetrain, остальные могут использоваться в других категориях
DELETE FROM category_attributes WHERE name = 'drivetrain';