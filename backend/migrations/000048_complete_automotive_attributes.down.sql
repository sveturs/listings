-- Откат добавления автомобильных атрибутов
-- Migration 000048 DOWN: Remove automotive attributes

-- Удаление связей категорий с атрибутами
DELETE FROM unified_category_attributes 
WHERE attribute_id IN (
    SELECT id FROM unified_attributes 
    WHERE code IN (
        'auto_part_brand', 'auto_part_oem', 'auto_part_condition', 'auto_compatibility',
        'auto_warranty', 'auto_year_from', 'auto_year_to', 'auto_installation',
        'tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 
        'tire_speed_index', 'tire_load_index',
        'rim_diameter', 'rim_width', 'rim_bolt_pattern', 'rim_offset', 'rim_center_bore',
        'engine_volume', 'engine_power', 'engine_type',
        'vehicle_capacity', 'vehicle_seats', 'vehicle_axles'
    )
);

-- Удаление созданных атрибутов
DELETE FROM unified_attributes WHERE code IN (
    'auto_part_brand', 'auto_part_oem', 'auto_part_condition', 'auto_compatibility',
    'auto_warranty', 'auto_year_from', 'auto_year_to', 'auto_installation',
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 
    'tire_speed_index', 'tire_load_index',
    'rim_diameter', 'rim_width', 'rim_bolt_pattern', 'rim_offset', 'rim_center_bore',
    'engine_volume', 'engine_power', 'engine_type',
    'vehicle_capacity', 'vehicle_seats', 'vehicle_axles'
);