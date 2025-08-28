-- Remove Serbian translations for tire attributes

DELETE FROM translations 
WHERE entity_type = 'attribute' 
AND entity_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 
                   'tire_brand', 'tire_condition', 'tread_depth', 'tire_year', 'tire_quantity')
)
AND language = 'sr';

-- Remove option translations (if they were added by this migration)
DELETE FROM attribute_option_translations 
WHERE attribute_name IN ('tire_season', 'tire_condition', 'tire_quantity');