-- Откат добавления сербских переводов для автомобильных категорий
-- Migration 000046 DOWN: Remove Serbian translations for marketplace categories

-- Удаление добавленных переводов
DELETE FROM translations WHERE language = 'sr' AND key_name IN (
    'category.automotive.additional_equipment',
    'category.automotive.electrical_electronic',
    'category.automotive.tires_wheels',
    'category.automotive.body_parts',
    'category.automotive.engine_parts',
    'category.automotive.cooling_system',
    'category.automotive.suspension_steering',
    'category.automotive.exhaust_system',
    'category.automotive.fuel_system',
    'category.automotive.lighting',
    'category.automotive.interior_accessories',
    'category.automotive.exterior_accessories',
    'category.automotive.car_care',
    'category.automotive.tools_equipment',
    'category.automotive.safety_security',
    'category.automotive.navigation_electronics',
    'category.automotive.car_audio',
    'category.automotive.performance_tuning',
    'category.automotive.vintage_classic',
    'category.automotive.motorcycles',
    'category.automotive.bicycles',
    'category.automotive.boats_marine',
    'category.automotive.trucks_commercial',
    'category.automotive.trailers_parts',
    'category.automotive.racing_sports',
    'category.automotive.maintenance_repair',
    'category.automotive.documentation',
    'category.automotive.insurance_financing'
);

-- Удаление записи из статистики
DELETE FROM translation_stats WHERE category = 'automotive_serbian_additions';