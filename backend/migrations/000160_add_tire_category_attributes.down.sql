-- Remove tire category attributes

-- First remove attribute option translations
DELETE FROM attribute_option_translations 
WHERE attribute_name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);

-- Remove translations
DELETE FROM translations 
WHERE entity_type = 'attribute' 
AND entity_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
        'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
        'tire_quantity'
    )
);

-- Remove category mappings
DELETE FROM category_attribute_mapping 
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
        'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
        'tire_quantity'
    )
);

-- Finally remove attributes
DELETE FROM category_attributes 
WHERE name IN (
    'tire_width', 'tire_profile', 'tire_diameter', 'tire_season',
    'tire_brand', 'tire_condition', 'tread_depth', 'tire_year',
    'tire_quantity'
);