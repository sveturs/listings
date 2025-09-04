-- Откат добавления атрибутов для 7 категорий
-- Migration 000047 DOWN: Remove attributes for missing categories

-- Удаление связей категорий с атрибутами
DELETE FROM unified_category_attributes 
WHERE category_id IN (1014, 1013, 1016, 1020, 1019, 1017, 1015);

-- Удаление созданных атрибутов
DELETE FROM unified_attributes WHERE code IN (
    -- Health & Beauty
    'health_product_type', 'health_brand', 'health_condition', 'health_expiry_date', 'health_skin_type',
    -- Kids & Baby
    'kids_age_group', 'kids_product_type', 'kids_brand', 'kids_condition', 'kids_size', 'kids_gender',
    -- Musical Instruments
    'music_instrument_type', 'music_brand', 'music_condition', 'music_skill_level', 'music_acoustic_electric',
    -- Events & Tickets
    'event_type', 'event_date', 'event_location', 'event_quantity', 'event_seat_type',
    -- Education
    'edu_type', 'edu_subject', 'edu_level', 'edu_language', 'edu_format',
    -- Antiques & Art
    'art_type', 'art_period', 'art_condition', 'art_authenticity', 'art_material', 'art_dimensions',
    -- Hobbies & Entertainment
    'hobby_type', 'hobby_brand', 'hobby_condition', 'hobby_age_group', 'hobby_skill_level', 'hobby_players'
);

-- Удаление статистики атрибутов
DELETE FROM unified_attribute_stats 
WHERE category_id IN (1014, 1013, 1016, 1020, 1019, 1017, 1015);